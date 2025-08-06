// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package httpserver

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"expvar"
	"net"
	"net/http"
	"net/http/pprof"
	"regexp"
	"strconv"
	"sync/atomic"
	"testing"
	"time"

	v1 "github.com/aileron-gateway/aileron-gateway/apis/core/v1"
	k "github.com/aileron-gateway/aileron-gateway/apis/kernel"
	"github.com/aileron-gateway/aileron-gateway/core"
	"github.com/aileron-gateway/aileron-gateway/kernel/api"
	"github.com/aileron-gateway/aileron-gateway/kernel/log"
	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/quic-go/quic-go"
	"github.com/quic-go/quic-go/http3"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func TestMutate(t *testing.T) {
	type condition struct {
		manifest protoreflect.ProtoMessage
	}

	type action struct {
		manifest protoreflect.ProtoMessage
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	cndDefault := tb.Condition("default", "input default manifest")
	actCheckMutated := tb.Action("check mutated", "check that the intended fields are mutated")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"mutate default",
			[]string{cndDefault},
			[]string{actCheckMutated},
			&condition{
				manifest: Resource.Default(),
			},
			&action{
				manifest: &v1.HTTPServer{
					APIVersion: apiVersion,
					Kind:       kind,
					Metadata: &k.Metadata{
						Namespace: "default",
						Name:      "default",
					},
					Spec: &v1.HTTPServerSpec{
						Addr:            ":8080",
						ShutdownTimeout: 30,
						HTTPConfig: &v1.HTTPConfig{
							ReadTimeout:       30,
							ReadHeaderTimeout: 30,
							WriteTimeout:      30,
							IdleTimeout:       10,
							MaxHeaderBytes:    1 << 13,
						},
					},
				},
			},
		),

		gen(
			"mutate NextProto",
			[]string{},
			[]string{actCheckMutated},
			&condition{
				manifest: &v1.HTTPServer{
					Spec: &v1.HTTPServerSpec{
						HTTPConfig: &v1.HTTPConfig{
							ListenConfig: &k.ListenConfig{
								TLSConfig: &k.TLSConfig{},
							},
						},
						HTTP2Config: &v1.HTTP2Config{},
					},
				},
			},
			&action{
				manifest: &v1.HTTPServer{
					Spec: &v1.HTTPServerSpec{
						HTTPConfig: &v1.HTTPConfig{
							ListenConfig: &k.ListenConfig{
								TLSConfig: &k.TLSConfig{
									NextProtos: []string{http2.NextProtoTLS},
								},
							},
						},
						HTTP2Config: &v1.HTTP2Config{},
					},
				},
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			manifest := Resource.Mutate(tt.C().manifest)

			opts := []cmp.Option{
				cmpopts.IgnoreUnexported(k.Metadata{}, k.Status{}, k.Reference{}),
				cmpopts.IgnoreUnexported(v1.HTTPServer{}, v1.HTTPServerSpec{}),
				cmpopts.IgnoreUnexported(v1.HTTPConfig{}, v1.HTTP2Config{}, v1.HTTP3Config{}),
				cmpopts.IgnoreUnexported(k.ListenConfig{}, k.TLSConfig{}),
				cmpopts.IgnoreUnexported(v1.VirtualHostSpec{}),
			}
			testutil.Diff(t, tt.A().manifest, manifest, opts...)
		})
	}
}

func TestCreate(t *testing.T) {
	type condition struct {
		manifest protoreflect.ProtoMessage
	}

	type action struct {
		expect     any
		err        any // error or errorutil.Kind
		errPattern *regexp.Regexp
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	// Get available port for testing.
	ln, err := net.Listen("tcp", ":0")
	if err != nil {
		panic(err)
	}
	ln.Close()

	// Create an API mux and register apis to the mux.
	testAPI := api.NewContainerAPI()
	postTestResource(testAPI, "mux", &testServeMux{
		ServeMux: &http.ServeMux{},
	})
	postTestResource(testAPI, "handler", &testHandler{
		headers:  map[string]string{"test": "handler"},
		body:     "test",
		patterns: []string{"/test"},
	})
	postTestResource(testAPI, "middleware", &testMiddleware{
		headers: map[string]string{"test": "middleware"},
	})

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"create with default manifest",
			[]string{}, []string{},
			&condition{
				manifest: Resource.Default(),
			},
			&action{
				expect: &runner{
					lg:      log.GlobalLogger(log.DefaultLoggerName),
					timeout: 30 * time.Second,
					svr: &http2Server{
						svr: &http.Server{
							Addr:                         ln.Addr().String(),
							DisableGeneralOptionsHandler: true,
							ReadTimeout:                  30 * time.Second,
							ReadHeaderTimeout:            30 * time.Second,
							WriteTimeout:                 30 * time.Second,
							IdleTimeout:                  10 * time.Second,
							MaxHeaderBytes:               1 << 13,
						},
						listener: &net.TCPListener{},
					},
				},
				err: nil,
			},
		),
		gen(
			"create http1 server only",
			[]string{}, []string{},
			&condition{
				manifest: &v1.HTTPServer{
					Metadata: &k.Metadata{},
					Spec: &v1.HTTPServerSpec{
						Addr:       ln.Addr().String(),
						HTTPConfig: &v1.HTTPConfig{},
						VirtualHosts: []*v1.VirtualHostSpec{
							{
								Pattern: "/test",
								Hosts:   []string{"test.com"},
							},
						},
					},
				},
			},
			&action{
				expect: &runner{
					lg: log.GlobalLogger(log.DefaultLoggerName),
					svr: &http2Server{
						listener: &net.TCPListener{},
						svr: &http.Server{
							Addr:                         ln.Addr().String(),
							Handler:                      &http.ServeMux{},
							DisableGeneralOptionsHandler: true,
						},
					},
				},
				err: nil,
			},
		),
		gen(
			"create http2 server only",
			[]string{}, []string{},
			&condition{
				manifest: &v1.HTTPServer{
					Metadata: &k.Metadata{},
					Spec: &v1.HTTPServerSpec{
						Addr:        ln.Addr().String(),
						HTTPConfig:  &v1.HTTPConfig{},
						HTTP2Config: &v1.HTTP2Config{},
						VirtualHosts: []*v1.VirtualHostSpec{
							{
								Pattern: "/test",
								Hosts:   []string{"test.com"},
							},
						},
					},
				},
			},
			&action{
				expect: &runner{
					lg: log.GlobalLogger(log.DefaultLoggerName),
					svr: &http2Server{
						listener: &net.TCPListener{},
						svr: &http.Server{
							Addr:                         ln.Addr().String(),
							Handler:                      &http.ServeMux{},
							DisableGeneralOptionsHandler: true,
							TLSConfig: &tls.Config{
								NextProtos:               []string{"h2", "http/1.1"},
								PreferServerCipherSuites: true,
							},
						},
					},
				},
				err: nil,
			},
		),
		gen(
			"create http3 server only",
			[]string{}, []string{},
			&condition{
				manifest: &v1.HTTPServer{
					Metadata: &k.Metadata{},
					Spec: &v1.HTTPServerSpec{
						Addr:        ln.Addr().String(),
						HTTP3Config: &v1.HTTP3Config{},
						VirtualHosts: []*v1.VirtualHostSpec{
							{
								Pattern: "/test",
								Hosts:   []string{"test.com"},
							},
						},
					},
				},
			},
			&action{
				expect: &runner{
					lg: log.GlobalLogger(log.DefaultLoggerName),
					svr: &http3Server{
						svr: &http3.Server{
							Addr:    ln.Addr().String(),
							Handler: &http.ServeMux{},
						},
					},
				},
				err: nil,
			},
		),
		gen(
			"specify virtual hosts",
			[]string{}, []string{},
			&condition{
				manifest: &v1.HTTPServer{
					Metadata: &k.Metadata{},
					Spec: &v1.HTTPServerSpec{
						Addr:       ln.Addr().String(),
						HTTPConfig: &v1.HTTPConfig{},
						VirtualHosts: []*v1.VirtualHostSpec{
							{
								Pattern: "/test",
								Hosts:   []string{"test.com"},
							},
						},
					},
				},
			},
			&action{
				expect: &runner{
					lg: log.GlobalLogger(log.DefaultLoggerName),
					svr: &http2Server{
						listener: &net.TCPListener{},
						svr: &http.Server{
							Addr:                         ln.Addr().String(),
							Handler:                      &http.ServeMux{},
							DisableGeneralOptionsHandler: true,
						},
					},
				},
				err: nil,
			},
		),
		gen(
			"register not found",
			[]string{}, []string{},
			&condition{
				manifest: &v1.HTTPServer{
					Metadata: &k.Metadata{},
					Spec: &v1.HTTPServerSpec{
						Addr:       ln.Addr().String(),
						HTTPConfig: &v1.HTTPConfig{},
						VirtualHosts: []*v1.VirtualHostSpec{
							{
								Pattern: "/test",
								Handlers: []*v1.HTTPHandlerSpec{{
									Handler: testResourceRef("handler"),
								}},
							},
						},
					},
				},
			},
			&action{
				expect: &runner{
					lg: log.GlobalLogger(log.DefaultLoggerName),
					svr: &http2Server{
						listener: &net.TCPListener{},
						svr: &http.Server{
							Addr:                         ln.Addr().String(),
							Handler:                      &http.ServeMux{},
							DisableGeneralOptionsHandler: true,
						},
					},
				},
				err: nil,
			},
		),
		gen(
			"error register virtual host mux",
			[]string{}, []string{},
			&condition{
				manifest: &v1.HTTPServer{
					Metadata: &k.Metadata{},
					Spec: &v1.HTTPServerSpec{
						Addr: ln.Addr().String(),
						VirtualHosts: []*v1.VirtualHostSpec{
							{
								Pattern: "/test",
								Hosts:   []string{"test.com", "test.com"},
								Middleware: []*k.Reference{
									testResourceRef("not exist middleware"),
								},
							},
						},
					},
				},
			},
			&action{
				expect:     nil,
				err:        core.ErrCoreGenCreateObject,
				errPattern: regexp.MustCompile(core.ErrPrefix + `failed to create HTTPServer`),
			},
		),
		gen(
			"fail to get middleware",
			[]string{}, []string{},
			&condition{
				manifest: &v1.HTTPServer{
					Metadata: &k.Metadata{},
					Spec: &v1.HTTPServerSpec{
						Middleware: []*k.Reference{
							{
								APIVersion: "wrong",
							},
						},
						Addr: ln.Addr().String(),
					},
				},
			},
			&action{
				expect:     nil,
				err:        core.ErrCoreGenCreateObject,
				errPattern: regexp.MustCompile(core.ErrPrefix + `failed to create HTTPServer`),
			},
		),
		gen(
			"invalid http2 server",
			[]string{}, []string{},
			&condition{
				manifest: &v1.HTTPServer{
					Metadata: &k.Metadata{},
					Spec: &v1.HTTPServerSpec{
						Addr: ln.Addr().String(),
						HTTPConfig: &v1.HTTPConfig{
							ListenConfig: &k.ListenConfig{
								TLSConfig: &k.TLSConfig{
									ClientAuth: 999,
								},
							},
						},
					},
				},
			},
			&action{
				expect:     nil,
				err:        core.ErrCoreGenCreateObject,
				errPattern: regexp.MustCompile(core.ErrPrefix + `failed to create HTTPServer`),
			},
		),
		gen(
			"invalid http3 server",
			[]string{}, []string{},
			&condition{
				manifest: &v1.HTTPServer{
					Metadata: &k.Metadata{},
					Spec: &v1.HTTPServerSpec{
						Addr: ln.Addr().String(),
						HTTP3Config: &v1.HTTP3Config{
							TLSConfig: &k.TLSConfig{
								ClientAuth: 999,
							},
						},
					},
				},
			},
			&action{
				expect:     nil,
				err:        core.ErrCoreGenCreateObject,
				errPattern: regexp.MustCompile(core.ErrPrefix + `failed to create HTTPServer`),
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			got, err := Resource.Create(testAPI, tt.C().manifest)
			testutil.DiffError(t, tt.A().err, tt.A().errPattern, err)

			defer func() {
				server, ok := got.(*runner)
				if ok && server.svr != nil {
					ctx, cancel := context.WithCancel(context.Background())
					go func() {
						time.Sleep(100 * time.Millisecond) // Wait server to  start.
						cancel()
					}()
					server.Run(ctx)
					time.Sleep(200 * time.Millisecond) // Wait server to stop.
				}
			}()

			opts := []cmp.Option{
				cmp.Comparer(testutil.ComparePointer[log.Logger]),
				cmp.AllowUnexported(runner{}),
				cmp.AllowUnexported(http2Server{}, http3Server{}),
				cmpopts.IgnoreUnexported(http.Server{}, http2.Server{}, http.ServeMux{}, http3.Server{}, net.TCPListener{}, tls.Config{}),
				cmpopts.IgnoreUnexported(atomic.Int32{}),
				cmpopts.IgnoreTypes(http.HandlerFunc(nil)),
				cmpopts.IgnoreInterfaces(struct{ http.Handler }{}),
				cmpopts.IgnoreFields(http.Server{}, "TLSNextProto"),
				cmpopts.IgnoreFields(http.Server{}, "Addr"),
			}
			testutil.Diff(t, tt.A().expect, got, opts...)
		})
	}
}

func TestNewHTTP2Server(t *testing.T) {
	type condition struct {
		addr string
		h    http.Handler
		c1   *v1.HTTPConfig
		c2   *v1.HTTP2Config
	}

	type action struct {
		svr        *http2Server
		err        any // error or errorutil.Kind
		errPattern *regexp.Regexp
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	cndInvalidTLSConfig := tb.Condition("invalid TLS config", "input an invalid TLS configuration")
	actCheckError := tb.Action("check the returned error", "check that the returned error is the one expected")
	actCheckNoError := tb.Action("check no error", "check that there is no error returned")
	table := tb.Build()

	// Get available port for testing.
	ln, err := net.Listen("tcp", ":0")
	if err != nil {
		panic(err)
	}
	ln.Close()
	testPort := ":" + strconv.Itoa(ln.Addr().(*net.TCPAddr).Port)

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"nil",
			[]string{},
			[]string{actCheckNoError},
			&condition{
				addr: testPort,
				h:    &testHandler{id: "test"},
			},
			&action{
				svr: nil,
				err: nil,
			},
		),
		gen(
			"zero HTTP config",
			[]string{},
			[]string{actCheckNoError},
			&condition{
				addr: testPort,
				h:    &testHandler{id: "test"},
				c1:   &v1.HTTPConfig{},
			},
			&action{
				svr: &http2Server{
					svr: &http.Server{
						Addr:                         "",
						Handler:                      &testHandler{id: "test"},
						DisableGeneralOptionsHandler: true,
					},
					listener: &net.TCPListener{},
				},
				err: nil,
			},
		),
		gen(
			"full HTTP config",
			[]string{},
			[]string{actCheckNoError},
			&condition{
				addr: testPort,
				h:    &testHandler{id: "test"},
				c1: &v1.HTTPConfig{
					EnableGeneralOptionsHandler: true,
					ReadTimeout:                 10,
					ReadHeaderTimeout:           11,
					WriteTimeout:                12,
					IdleTimeout:                 13,
					MaxHeaderBytes:              14,
					DisableKeepAlive:            true,
					AltSvc:                      "foo",
				},
			},
			&action{
				svr: &http2Server{
					svr: &http.Server{
						Addr:                         "",
						Handler:                      altSvcMiddleware("foo").Middleware(&testHandler{id: "test"}),
						DisableGeneralOptionsHandler: false,
						ReadTimeout:                  10 * time.Second,
						ReadHeaderTimeout:            11 * time.Second,
						WriteTimeout:                 12 * time.Second,
						IdleTimeout:                  13 * time.Second,
						MaxHeaderBytes:               14,
					},
					listener: &net.TCPListener{},
				},
				err: nil,
			},
		),
		gen(
			"zero HTTP2 config",
			[]string{},
			[]string{actCheckNoError},
			&condition{
				addr: testPort,
				h:    &testHandler{id: "test"},
				c1:   &v1.HTTPConfig{},
				c2:   &v1.HTTP2Config{},
			},
			&action{
				svr: &http2Server{
					svr: &http.Server{
						Addr:                         "",
						Handler:                      &testHandler{id: "test"},
						DisableGeneralOptionsHandler: true,
						TLSConfig: &tls.Config{
							NextProtos:               []string{"h2", "http/1.1"},
							PreferServerCipherSuites: true,
						},
					},
					listener: &net.TCPListener{},
				},
				err: nil,
			},
		),
		gen(
			"full HTTP2 config",
			[]string{},
			[]string{actCheckNoError},
			&condition{
				addr: testPort,
				h:    &testHandler{id: "test"},
				c1:   &v1.HTTPConfig{},
				c2: &v1.HTTP2Config{
					MaxConcurrentStreams:         10,
					MaxDecoderHeaderTableSize:    11,
					MaxEncoderHeaderTableSize:    12,
					MaxReadFrameSize:             13,
					PermitProhibitedCipherSuites: true,
					IdleTimeout:                  14,
					MaxUploadBufferPerConnection: 15,
					MaxUploadBufferPerStream:     16,
					EnableH2C:                    true,
				},
			},
			&action{
				svr: &http2Server{
					svr: &http.Server{
						Addr: "",
						Handler: h2c.NewHandler(&testHandler{id: "test"},
							&http2.Server{
								MaxConcurrentStreams:         10,
								MaxDecoderHeaderTableSize:    11,
								MaxEncoderHeaderTableSize:    12,
								MaxReadFrameSize:             13,
								PermitProhibitedCipherSuites: true,
								IdleTimeout:                  14 * time.Second,
								MaxUploadBufferPerConnection: 15,
								MaxUploadBufferPerStream:     16,
							}),
						DisableGeneralOptionsHandler: true,
						TLSConfig: &tls.Config{
							NextProtos:               []string{"h2", "http/1.1"},
							PreferServerCipherSuites: true,
						},
					},
					listener: &net.TCPListener{},
				},
				err: nil,
			},
		),
		gen(
			"invalid TLS Config",
			[]string{cndInvalidTLSConfig},
			[]string{actCheckError},
			&condition{
				addr: testPort,
				h:    &testHandler{id: "test"},
				c1: &v1.HTTPConfig{
					ListenConfig: &k.ListenConfig{
						TLSConfig: &k.TLSConfig{
							ClientAuth: 999,
						},
					},
				},
			},
			&action{
				svr:        nil,
				err:        core.ErrCoreGenCreateObject,
				errPattern: regexp.MustCompile(core.ErrPrefix + `failed to create HTTPServer`),
			},
		),
		gen(
			"invalid HTTP2 config",
			[]string{},
			[]string{actCheckNoError},
			&condition{
				addr: testPort,
				h:    &testHandler{id: "test"},
				c1: &v1.HTTPConfig{
					ListenConfig: &k.ListenConfig{
						TLSConfig: &k.TLSConfig{
							MinVersion: tls.VersionTLS10,
							TLSCiphers: []k.TLSCipher{
								k.TLSCipher_TLS_CHACHA20_POLY1305_SHA256,
							},
						},
					},
				},
				c2: &v1.HTTP2Config{},
			},
			&action{
				svr:        nil,
				err:        core.ErrCoreGenCreateObject,
				errPattern: regexp.MustCompile(core.ErrPrefix + `failed to create HTTPServer`),
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			srv, err := newHTTP2Server(tt.C().addr, tt.C().h, tt.C().c1, tt.C().c2)
			testutil.DiffError(t, tt.A().err, tt.A().errPattern, err)
			if srv != nil {
				srv.Shutdown(context.Background())
				srv.listener.Close()
			}

			opts := []cmp.Option{
				cmp.AllowUnexported(http2Server{}, testHandler{}),
				cmpopts.IgnoreUnexported(http.Server{}, http2.Server{}),
				cmpopts.IgnoreUnexported(net.TCPListener{}, tls.Config{}),
				cmpopts.IgnoreTypes(http.HandlerFunc(nil)), // Skip alt-svc middlewarte.
				cmpopts.IgnoreFields(http.Server{}, "TLSNextProto"),
				testutil.DeepAllowUnexported(h2c.NewHandler(nil, nil)),
			}
			testutil.Diff(t, tt.A().svr, srv, opts...)
		})
	}
}

func TestNewHTTP3Server(t *testing.T) {
	type condition struct {
		addr string
		h    http.Handler
		c    *v1.HTTP3Config
	}

	type action struct {
		svr        *http3Server
		err        any // error or errorutil.Kind
		errPattern *regexp.Regexp
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	cndInvalidTLSConfig := tb.Condition("invalid TLS config", "input an invalid TLS configuration")
	actCheckError := tb.Action("check the returned error", "check that the returned error is the one expected")
	actCheckNoError := tb.Action("check no error", "check that there is no error returned")
	table := tb.Build()

	// Get available port for testing.
	ln, err := net.Listen("tcp", ":0")
	if err != nil {
		panic(err)
	}
	ln.Close()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"nil",
			[]string{},
			[]string{actCheckNoError},
			&condition{
				addr: ln.Addr().String(),
				h:    &testHandler{id: "test"},
			},
			&action{
				svr: nil,
				err: nil,
			},
		),
		gen(
			"zero config",
			[]string{},
			[]string{actCheckNoError},
			&condition{
				addr: ln.Addr().String(),
				h:    &testHandler{id: "test"},
				c:    &v1.HTTP3Config{},
			},
			&action{
				svr: &http3Server{
					svr: &http3.Server{
						Addr:    ln.Addr().String(),
						Handler: &testHandler{id: "test"},
					},
				},
				err: nil,
			},
		),
		gen(
			"full config",
			[]string{},
			[]string{actCheckNoError},
			&condition{
				addr: ln.Addr().String(),
				h:    &testHandler{id: "test"},
				c: &v1.HTTP3Config{
					QuicConfig:     &k.QuicConfig{},
					TLSConfig:      &k.TLSConfig{},
					MaxHeaderBytes: 10,
					AltSvc:         "foo",
				},
			},
			&action{
				svr: &http3Server{
					svr: &http3.Server{
						Addr:           ln.Addr().String(),
						Handler:        altSvcMiddleware("foo").Middleware(&testHandler{id: "test"}),
						MaxHeaderBytes: 10,
						QUICConfig: &quic.Config{
							Versions: []quic.Version{},
						},
						TLSConfig: &tls.Config{
							RootCAs:      x509.NewCertPool(),
							ClientCAs:    x509.NewCertPool(),
							Certificates: []tls.Certificate{},
						},
					},
				},
				err: nil,
			},
		),
		gen(
			"invalid TLS Config",
			[]string{cndInvalidTLSConfig},
			[]string{actCheckError},
			&condition{
				addr: ln.Addr().String(),
				h:    &testHandler{id: "test"},
				c: &v1.HTTP3Config{
					TLSConfig: &k.TLSConfig{
						ClientAuth: 999,
					},
				},
			},
			&action{
				svr:        nil,
				err:        core.ErrCoreGenCreateObject,
				errPattern: regexp.MustCompile(core.ErrPrefix + `failed to create HTTPServer`),
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			srv, err := newHTTP3Server(tt.C().addr, tt.C().h, tt.C().c)
			testutil.DiffError(t, tt.A().err, tt.A().errPattern, err)
			if srv != nil {
				srv.Shutdown(context.Background())
			}

			opts := []cmp.Option{
				cmp.AllowUnexported(http3Server{}, testHandler{}, quic.Config{}),
				cmpopts.IgnoreUnexported(http3.Server{}),
				cmpopts.IgnoreUnexported(tls.Config{}),
				cmpopts.IgnoreFields(tls.Config{}, "RootCAs", "ClientCAs"),
				cmpopts.IgnoreTypes(http.HandlerFunc(nil)), // Skip alt-svc middlewarte.
			}
			testutil.Diff(t, tt.A().svr, srv, opts...)
		})
	}
}

type testMiddleware struct {
	id      string
	headers map[string]string
}

func (m *testMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for k, v := range m.headers {
			w.Header().Add(k, v)
		}

		next.ServeHTTP(w, r)
	})
}

type testHandler struct {
	id       string
	headers  map[string]string
	status   int
	body     string
	patterns []string
	methods  []string
}

func (h *testHandler) Patterns() []string {
	return h.patterns
}

func (h *testHandler) Methods() []string {
	return h.methods
}

func (h *testHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for k, v := range h.headers {
		w.Header().Add(k, v)
	}
	if h.status > 0 {
		w.WriteHeader(h.status)
	}
	w.Write([]byte(h.body))
}

type testServeMux struct {
	*http.ServeMux
}

func (m *testServeMux) Method(method string, path string, h http.Handler) {
	m.ServeMux.Handle(path, h)
}

func (m *testServeMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.ServeMux.ServeHTTP(w, r)
}

type testMux struct {
	Mux
	hs map[string]http.Handler
}

func (m *testMux) Handle(path string, h http.Handler) {
	m.hs[path] = h
	if m.Mux != nil {
		m.Mux.Handle(path, h)
	}
}

func (m *testMux) Method(string, string, http.Handler) {}

func (m *testMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if m.Mux != nil {
		m.Mux.ServeHTTP(w, r)
	}
}

func TestRegisterProfile(t *testing.T) {
	type condition struct {
		enabled bool
	}

	type action struct {
		path    string
		handler http.Handler
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"index handler",
			[]string{},
			[]string{},
			&condition{
				enabled: true,
			},
			&action{
				path:    "GET /debug/pprof/",
				handler: http.HandlerFunc(pprof.Index),
			},
		),
		gen(
			"cmdline handler",
			[]string{},
			[]string{},
			&condition{
				enabled: true,
			},
			&action{
				path:    "GET /debug/pprof/cmdline",
				handler: http.HandlerFunc(pprof.Cmdline),
			},
		),
		gen(
			"profile handler",
			[]string{},
			[]string{},
			&condition{
				enabled: true,
			},
			&action{
				path:    "GET /debug/pprof/profile",
				handler: http.HandlerFunc(pprof.Profile),
			},
		),
		gen(
			"symbol handler",
			[]string{},
			[]string{},
			&condition{
				enabled: true,
			},
			&action{
				path:    "GET /debug/pprof/symbol",
				handler: http.HandlerFunc(pprof.Symbol),
			},
		),
		gen(
			"trace handler",
			[]string{},
			[]string{},
			&condition{
				enabled: true,
			},
			&action{
				path:    "GET /debug/pprof/trace",
				handler: http.HandlerFunc(pprof.Trace),
			},
		),
		gen(
			"profile disabled",
			[]string{},
			[]string{},
			&condition{
				enabled: false,
			},
			&action{
				path:    "GET /debug/pprof/",
				handler: nil,
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			mux := &testMux{
				hs: map[string]http.Handler{},
			}
			registerProfile(mux, tt.C().enabled)
			h := mux.hs[tt.A().path]
			testutil.Diff(t, tt.A().handler, h, cmp.Comparer(testutil.ComparePointer[http.Handler]))
		})
	}
}

func TestRegisterExpvar(t *testing.T) {
	type condition struct {
		enabled bool
	}
	type action struct {
		path    string
		handler http.Handler
	}
	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"register",
			[]string{},
			[]string{},
			&condition{
				enabled: true,
			},
			&action{
				path:    "GET /debug/vars",
				handler: expvar.Handler(),
			},
		),
		gen(
			"disabled",
			[]string{},
			[]string{},
			&condition{
				enabled: false,
			},
			&action{
				path:    "GET /debug/vars",
				handler: nil,
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			mux := &testMux{
				hs: map[string]http.Handler{},
			}
			registerExpvar(mux, tt.C().enabled)
			h := mux.hs[tt.A().path]
			testutil.Diff(t, tt.A().handler, h, cmp.Comparer(testutil.ComparePointer[http.Handler]))
		})
	}
}

func postTestResource(server api.API[*api.Request, *api.Response], name string, res any) {
	ref := testResourceRef(name)
	req := &api.Request{
		Method:  api.MethodPost,
		Key:     ref.APIVersion + "/" + ref.Kind + "/" + ref.Namespace + "/" + ref.Name,
		Content: res,
	}
	if _, err := server.Serve(context.Background(), req); err != nil {
		panic(err)
	}
}

func testResourceRef(name string) *k.Reference {
	return &k.Reference{
		APIVersion: "core/v1",
		Kind:       "Container",
		Namespace:  "test",
		Name:       name,
	}
}
