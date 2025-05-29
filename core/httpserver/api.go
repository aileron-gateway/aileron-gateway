// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package httpserver

import (
	"cmp"
	"crypto/tls"
	"expvar"
	"net/http"
	"net/http/pprof"
	"runtime/debug"
	"slices"
	"strings"
	"time"

	v1 "github.com/aileron-gateway/aileron-gateway/apis/core/v1"
	"github.com/aileron-gateway/aileron-gateway/apis/kernel"
	"github.com/aileron-gateway/aileron-gateway/core"
	"github.com/aileron-gateway/aileron-gateway/kernel/api"
	"github.com/aileron-gateway/aileron-gateway/kernel/log"
	"github.com/aileron-gateway/aileron-gateway/kernel/network"
	utilhttp "github.com/aileron-gateway/aileron-gateway/util/http"
	"github.com/quic-go/quic-go/http3"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/protobuf/reflect/protoreflect"
)

const (
	apiVersion = "core/v1"
	kind       = "HTTPServer"
	Key        = apiVersion + "/" + kind
)

var Resource api.Resource = &API{
	BaseResource: &api.BaseResource{
		DefaultProto: &v1.HTTPServer{
			APIVersion: apiVersion,
			Kind:       kind,
			Metadata: &kernel.Metadata{
				Namespace: "default",
				Name:      "default",
			},
			Spec: &v1.HTTPServerSpec{
				Addr:            ":8080",
				ShutdownTimeout: 30, // In second.
				HTTPConfig: &v1.HTTPConfig{
					ReadTimeout:       30, // In second.
					ReadHeaderTimeout: 30, // In second.
					WriteTimeout:      30, // In second.
					IdleTimeout:       10, // In second.
					MaxHeaderBytes:    1 << 13,
				},
				Profile: &v1.ProfileSpec{
					Enable:     false,
					PathPrefix: "/debug/pprof",
				},
				Expvar: &v1.ExpvarSpec{
					Enable: false,
					Path:   "/debug/vars",
				},
			},
		},
	},
}

type API struct {
	*api.BaseResource
}

func (*API) Mutate(msg protoreflect.ProtoMessage) protoreflect.ProtoMessage {
	c := msg.(*v1.HTTPServer)

	// Add "h2" to the NextProtos if HTTP2 server is enabled.
	h1c := c.Spec.HTTPConfig
	if h1c != nil && h1c.ListenConfig != nil && h1c.ListenConfig.TLSConfig != nil {
		tc := h1c.ListenConfig.TLSConfig
		if c.Spec.HTTP2Config != nil && !slices.Contains(tc.NextProtos, http2.NextProtoTLS) {
			tc.NextProtos = append([]string{http2.NextProtoTLS}, tc.NextProtos...)
		}
	}

	return c
}

func (*API) Create(a api.API[*api.Request, *api.Response], msg protoreflect.ProtoMessage) (any, error) {
	c := msg.(*v1.HTTPServer)

	lg := log.DefaultOr(c.Metadata.Logger)

	eh, err := utilhttp.ErrorHandler(a, c.Spec.ErrorHandler)
	if err != nil {
		return nil, core.ErrCoreGenCreateObject.WithStack(err, map[string]any{"kind": kind})
	}

	mux := &http.ServeMux{}
	registerProfile(mux, c.Spec.Profile)
	registerExpvar(mux, c.Spec.Expvar)

	nfh := notFoundHandler(eh)
	handlers, err := registerHandlers(a, mux, c.Spec.VirtualHosts, nfh)
	if err != nil {
		return nil, core.ErrCoreGenCreateObject.WithStack(err, map[string]any{"kind": kind})
	}

	// Register not found handler if possible.
	skipNotFound := false
	for k := range handlers {
		skipNotFound = skipNotFound || wildcardPath.MatchString(k)
	}
	if !skipNotFound {
		mux.Handle("/", notFoundHandler(eh))
	}

	middleware, err := api.ReferTypedObjects[core.Middleware](a, c.Spec.Middleware...)
	if err != nil {
		return nil, core.ErrCoreGenCreateObject.WithStack(err, map[string]any{"kind": kind})
	}
	middleware = append([]core.Middleware{&recoverer{lg: lg, eh: eh}}, middleware...)
	handler := utilhttp.MiddlewareChain(middleware, mux)

	runner := &runner{
		svr:     nil,
		lg:      lg,
		timeout: time.Duration(c.Spec.ShutdownTimeout) * time.Second,
	}

	if c.Spec.HTTP3Config != nil {
		svr, err := newHTTP3Server(c.Spec.Addr, handler, c.Spec.HTTP3Config)
		if err != nil {
			return nil, err
		}
		runner.svr = svr
	} else {
		println("===================== ", c.Spec.Addr)
		debug.PrintStack()
		svr, err := newHTTP2Server(c.Spec.Addr, handler, c.Spec.HTTPConfig, c.Spec.HTTP2Config)
		if err != nil {
			return nil, err
		}
		runner.svr = svr
	}
	return runner, nil
}

// newHTTP2Server returns a new http2 server.
// This function returns nil if the given HTTPConfig was nil.
// The listen address addr must not be an empty string.
// The http.Handler of h should not be nil.
func newHTTP2Server(addr string, h http.Handler, c *v1.HTTPConfig, c2 *v1.HTTP2Config) (*http2Server, error) {
	if c == nil {
		return nil, nil
	}

	c.ListenConfig = cmp.Or(c.ListenConfig, &kernel.ListenConfig{})
	c.ListenConfig.Network = cmp.Or(c.ListenConfig.Network, "tcp")
	c.ListenConfig.Addr = cmp.Or(c.ListenConfig.Addr, addr)
	listener, err := network.NewListenerFromSpec(c.ListenConfig)
	if err != nil {
		return nil, core.ErrCoreGenCreateObject.WithStack(err, map[string]any{"kind": kind})
	}

	// Even the TLS has already been configured for the listener,
	// tlsConfig is required when configuring http2 server.
	// No errors here.
	tlsConfig, _ := network.TLSConfig(c.ListenConfig.TLSConfig)

	if c.AltSvc != "" {
		h = altSvcMiddleware(c.AltSvc).Middleware(h)
	}

	svr := &http.Server{
		Addr:                         "", // Listener has already been configured.
		Handler:                      h,
		DisableGeneralOptionsHandler: !c.EnableGeneralOptionsHandler,
		TLSConfig:                    tlsConfig,
		ReadTimeout:                  time.Second * time.Duration(c.ReadTimeout),
		ReadHeaderTimeout:            time.Second * time.Duration(c.ReadHeaderTimeout),
		WriteTimeout:                 time.Second * time.Duration(c.WriteTimeout),
		IdleTimeout:                  time.Second * time.Duration(c.IdleTimeout),
		MaxHeaderBytes:               int(c.MaxHeaderBytes),
	}
	svr.SetKeepAlivesEnabled(!c.DisableKeepAlive)
	if !c.AllowHTTP2 {
		svr.TLSNextProto = make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0)
	}

	if c2 != nil {
		svr.TLSNextProto = nil
		conf := &http2.Server{
			MaxConcurrentStreams:         c2.MaxConcurrentStreams,
			MaxDecoderHeaderTableSize:    c2.MaxDecoderHeaderTableSize,
			MaxEncoderHeaderTableSize:    c2.MaxEncoderHeaderTableSize,
			MaxReadFrameSize:             c2.MaxReadFrameSize,
			PermitProhibitedCipherSuites: c2.PermitProhibitedCipherSuites,
			IdleTimeout:                  time.Second * time.Duration(c2.IdleTimeout),
			MaxUploadBufferPerConnection: c2.MaxUploadBufferPerConnection,
			MaxUploadBufferPerStream:     c2.MaxUploadBufferPerStream,
		}
		if err := http2.ConfigureServer(svr, conf); err != nil {
			listener.Close()
			return nil, core.ErrCoreGenCreateObject.WithStack(err, map[string]any{"kind": kind})
		}
		if c2.EnableH2C {
			svr.Handler = h2c.NewHandler(h, conf)
		}
	}

	return &http2Server{
		svr:      svr,
		listener: listener,
	}, nil
}

// newHTTP3Server returns a new http3 server.
// This function returns nil if the given HTTP3Config was nil.
// The listen address addr must not be an empty string.
// The http.Handler of h should not be nil.
func newHTTP3Server(addr string, h http.Handler, c *v1.HTTP3Config) (*http3Server, error) {
	if c == nil {
		return nil, nil
	}

	tlsConfig, err := network.TLSConfig(c.TLSConfig)
	if err != nil {
		return nil, core.ErrCoreGenCreateObject.WithStack(err, map[string]any{"kind": kind})
	}

	qc, _ := network.QuicConfig(c.QuicConfig) // No error here.

	if c.AltSvc != "" {
		h = altSvcMiddleware(c.AltSvc).Middleware(h)
	}

	return &http3Server{
		svr: &http3.Server{
			Addr:            addr,
			Port:            0, // This field is not used.
			TLSConfig:       tlsConfig,
			QUICConfig:      qc,
			Handler:         h,
			EnableDatagrams: false, // Should be set in QuicConfig.
			MaxHeaderBytes:  int(c.MaxHeaderBytes),
		},
	}, nil
}

// registerProfile registers profile handlers to the server if specified.
// Given argument must not be nil. This function panics when a nil was given as any argument.
// See https://pkg.go.dev/net/http/pprof.
// This function registers profile handlers to the mux with these paths.
//   - pprof.Index for PathPrefix+"/"
//   - pprof.Cmdline for PathPrefix+"/cmdline"
//   - pprof.Profile for PathPrefix+"/profile"
//   - pprof.Symbol for PathPrefix+"/symbol"
//   - pprof.Trace for PathPrefix+"/trace"
func registerProfile(mux Mux, spec *v1.ProfileSpec) {
	http.DefaultServeMux = &http.ServeMux{}
	if spec == nil || !spec.Enable {
		return // Profiling disabled.
	}
	mux.Handle("GET "+strings.TrimSuffix(spec.PathPrefix, "/")+"/", http.HandlerFunc(pprof.Index))
	mux.Handle("GET "+strings.TrimSuffix(spec.PathPrefix, "/")+"/cmdline", http.HandlerFunc(pprof.Cmdline))
	mux.Handle("GET "+strings.TrimSuffix(spec.PathPrefix, "/")+"/profile", http.HandlerFunc(pprof.Profile))
	mux.Handle("GET "+strings.TrimSuffix(spec.PathPrefix, "/")+"/symbol", http.HandlerFunc(pprof.Symbol))
	mux.Handle("GET "+strings.TrimSuffix(spec.PathPrefix, "/")+"/trace", http.HandlerFunc(pprof.Trace))
}

// registerExpvar registers expvar handlers to the server if specified.
// Given argument must not be nil. This function panics when a nil was given as any argument.
// See https://pkg.go.dev/expvar
func registerExpvar(mux Mux, spec *v1.ExpvarSpec) {
	http.DefaultServeMux = &http.ServeMux{}
	if spec == nil || !spec.Enable {
		return // Expvar disabled.
	}
	mux.Handle("GET "+spec.Path, expvar.Handler())
}
