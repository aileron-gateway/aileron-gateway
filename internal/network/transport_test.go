// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package network_test

import (
	"crypto/tls"
	"crypto/x509"
	"net"
	"net/http"
	"testing"
	"time"

	k "github.com/aileron-gateway/aileron-gateway/apis/kernel"
	"github.com/aileron-gateway/aileron-gateway/internal/network"
	"github.com/aileron-gateway/aileron-gateway/internal/testutil"
	"github.com/aileron-gateway/aileron-gateway/kernel/er"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/quic-go/quic-go"
	"github.com/quic-go/quic-go/http3"
	"golang.org/x/net/http2"
)

func TestQuickConfig(t *testing.T) {
	type condition struct {
		spec *k.QuicConfig
	}

	type action struct {
		config *quic.Config
		err    error
	}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"input nil",
			&condition{
				spec: nil,
			},
			&action{
				config: nil,
				err:    nil,
			},
		),
		gen(
			"input zero value spec",
			&condition{
				spec: &k.QuicConfig{},
			},
			&action{
				config: &quic.Config{
					Versions: []quic.Version{},
				},
				err: nil,
			},
		),
		gen(
			"valid values for spec",
			&condition{
				spec: &k.QuicConfig{
					Versions:                       []k.QuicConfig_Version{k.QuicConfig_Version1, k.QuicConfig_Version2},
					HandshakeIdleTimeout:           10,
					MaxIdleTimeout:                 11,
					InitialStreamReceiveWindow:     12,
					MaxStreamReceiveWindow:         13,
					InitialConnectionReceiveWindow: 14,
					MaxConnectionReceiveWindow:     15,
					MaxIncomingStreams:             16,
					MaxIncomingUniStreams:          17,
					KeepAlivePeriod:                18,
					DisablePathMTUDiscovery:        true,
					Allow0RTT:                      true,
					EnableDatagrams:                true,
				},
			},
			&action{
				config: &quic.Config{
					Versions:                       []quic.Version{quic.Version1, quic.Version2},
					HandshakeIdleTimeout:           time.Millisecond * 10,
					MaxIdleTimeout:                 time.Millisecond * 11,
					InitialStreamReceiveWindow:     12,
					MaxStreamReceiveWindow:         13,
					InitialConnectionReceiveWindow: 14,
					MaxConnectionReceiveWindow:     15,
					MaxIncomingStreams:             16,
					MaxIncomingUniStreams:          17,
					KeepAlivePeriod:                time.Millisecond * 18,
					DisablePathMTUDiscovery:        true,
					Allow0RTT:                      true,
					EnableDatagrams:                true,
				},
				err: nil,
			},
		),
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			c, err := network.QuicConfig(tt.C.spec)
			testutil.Diff(t, tt.A.err, err, cmpopts.EquateErrors())
			if err != nil {
				return
			}

			testutil.Diff(t, tt.A.config, c)
		})
	}
}

func TestHTTPTransport(t *testing.T) {
	type condition struct {
		spec *k.HTTPTransportConfig
	}

	type action struct {
		tp  *http.Transport
		err error
	}

	newSystemPool := func() *x509.CertPool {
		pool, err := x509.SystemCertPool()
		if err != nil {
			panic(err) // Bad test environment.
		}
		return pool
	}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"input nil",
			&condition{
				spec: nil,
			},
			&action{
				tp: &http.Transport{
					Proxy: http.ProxyFromEnvironment,
					DialContext: (&net.Dialer{
						Timeout:   10 * time.Second,
						KeepAlive: 30 * time.Second,
					}).DialContext,
					DisableCompression:    true,
					MaxIdleConnsPerHost:   1024,
					IdleConnTimeout:       60 * time.Second,
					TLSHandshakeTimeout:   10 * time.Second,
					ExpectContinueTimeout: 1 * time.Second,
				},
				err: nil,
			},
		),
		gen(
			"input zero value spec",
			&condition{
				spec: &k.HTTPTransportConfig{},
			},
			&action{
				tp: &http.Transport{
					Proxy:               http.ProxyFromEnvironment,
					MaxIdleConnsPerHost: 1024,
					TLSNextProto:        make(map[string]func(authority string, c *tls.Conn) http.RoundTripper, 0),
				},
				err: nil,
			},
		),
		gen(
			"specify valid TransportConfig values for spec",
			&condition{
				spec: &k.HTTPTransportConfig{
					TLSHandshakeTimeout:    10,
					DisableKeepAlives:      true,
					DisableCompression:     true,
					MaxIdleConns:           11,
					MaxIdleConnsPerHost:    12,
					MaxConnsPerHost:        13,
					IdleConnTimeout:        14,
					ResponseHeaderTimeout:  15,
					ExpectContinueTimeout:  16,
					MaxResponseHeaderBytes: 17,
					WriteBufferSize:        18,
					ReadBufferSize:         19,
					AllowHTTP2:             true,
				},
			},
			&action{
				tp: &http.Transport{
					Proxy:                  http.ProxyFromEnvironment,
					TLSHandshakeTimeout:    time.Millisecond * 10,
					DisableKeepAlives:      true,
					DisableCompression:     true,
					MaxIdleConns:           11,
					MaxIdleConnsPerHost:    12,
					MaxConnsPerHost:        13,
					IdleConnTimeout:        time.Millisecond * 14,
					ResponseHeaderTimeout:  time.Millisecond * 15,
					ExpectContinueTimeout:  time.Millisecond * 16,
					MaxResponseHeaderBytes: 17,
					WriteBufferSize:        18,
					ReadBufferSize:         19,
					ForceAttemptHTTP2:      false,
				},
				err: nil,
			},
		),
		gen(
			"specify valid TLSConfig values",
			&condition{
				spec: &k.HTTPTransportConfig{
					TLSConfig: &k.TLSConfig{},
				},
			},
			&action{
				tp: &http.Transport{
					Proxy:               http.ProxyFromEnvironment,
					MaxIdleConnsPerHost: 1024,
					TLSClientConfig: &tls.Config{
						RootCAs:      newSystemPool(),
						ClientCAs:    newSystemPool(),
						Certificates: []tls.Certificate{},
					},
					TLSNextProto: make(map[string]func(authority string, c *tls.Conn) http.RoundTripper, 0),
				},
				err: nil,
			},
		),
		gen(
			"specify invalid TLSConfig values",
			&condition{
				spec: &k.HTTPTransportConfig{
					TLSConfig: &k.TLSConfig{
						ClientAuth: 99999,
					},
				},
			},
			&action{
				tp: nil,
				err: &er.Error{
					Package:     network.ErrPkg,
					Type:        network.ErrTypeTransport,
					Description: network.ErrDscNewTransport,
				},
			},
		),
		gen(
			"specify valid DialConfig values/disable http2",
			&condition{
				spec: &k.HTTPTransportConfig{
					DialConfig: &k.DialConfig{},
					AllowHTTP2: false,
				},
			},
			&action{
				tp: &http.Transport{
					Proxy:               http.ProxyFromEnvironment,
					MaxIdleConnsPerHost: 1024,
					TLSNextProto:        make(map[string]func(authority string, c *tls.Conn) http.RoundTripper, 0),
				},
				err: nil,
			},
		),
		gen(
			"specify valid DialConfig values/enable http2",
			&condition{
				spec: &k.HTTPTransportConfig{
					DialConfig: &k.DialConfig{},
					AllowHTTP2: true,
				},
			},
			&action{
				tp: &http.Transport{
					Proxy:               http.ProxyFromEnvironment,
					MaxIdleConnsPerHost: 1024,
					ForceAttemptHTTP2:   true,
				},
				err: nil,
			},
		),
		gen(
			"specify invalid DialConfig values",
			&condition{
				spec: &k.HTTPTransportConfig{
					DialConfig: &k.DialConfig{
						LocalAddress: "tcp://INVALID_ADDRESS",
					},
				},
			},
			&action{
				tp: nil,
				err: &er.Error{
					Package:     network.ErrPkg,
					Type:        network.ErrTypeTransport,
					Description: network.ErrDscNewTransport,
				},
			},
		),
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			tp, err := network.HTTPTransport(tt.C.spec)
			testutil.Diff(t, tt.A.err, err, cmpopts.EquateErrors())
			if err != nil {
				return
			}
			opts := []cmp.Option{
				cmpopts.IgnoreUnexported(http.Transport{}, tls.Config{}),
				cmpopts.IgnoreFields(http.Transport{}, "Proxy", "DialContext", "DialTLSContext"),
			}
			testutil.Diff(t, tt.A.tp, tp, opts...)
		})
	}
}

func TestHTTP2Transport(t *testing.T) {
	type condition struct {
		spec *k.HTTP2TransportConfig
	}

	type action struct {
		tp  *http2.Transport
		err error
	}

	newSystemPool := func() *x509.CertPool {
		pool, err := x509.SystemCertPool()
		if err != nil {
			panic(err) // Bad test environment.
		}
		return pool
	}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"input nil",
			&condition{
				spec: nil,
			},
			&action{
				tp:  &http2.Transport{},
				err: nil,
			},
		),
		gen(
			"input zero value spec",
			&condition{
				spec: &k.HTTP2TransportConfig{},
			},
			&action{
				tp:  &http2.Transport{},
				err: nil,
			},
		),
		gen(
			"specify valid TransportConfig values for spec",
			&condition{
				spec: &k.HTTP2TransportConfig{
					DisableCompression:         true,
					AllowHTTP:                  true,
					MaxHeaderListSize:          10,
					MaxReadFrameSize:           11,
					MaxDecoderHeaderTableSize:  12,
					MaxEncoderHeaderTableSize:  13,
					StrictMaxConcurrentStreams: true,
					IdleConnTimeout:            14,
					ReadIdleTimeout:            15,
					PingTimeout:                16,
					WriteByteTimeout:           17,
				},
			},
			&action{
				tp: &http2.Transport{
					DisableCompression:         true,
					AllowHTTP:                  true,
					MaxHeaderListSize:          10,
					MaxReadFrameSize:           11,
					MaxDecoderHeaderTableSize:  12,
					MaxEncoderHeaderTableSize:  13,
					StrictMaxConcurrentStreams: true,
					IdleConnTimeout:            14 * time.Millisecond,
					ReadIdleTimeout:            15 * time.Millisecond,
					PingTimeout:                16 * time.Millisecond,
					WriteByteTimeout:           17 * time.Millisecond,
				},
				err: nil,
			},
		),
		gen(
			"specify valid TLSConfig values",
			&condition{
				spec: &k.HTTP2TransportConfig{
					TLSConfig: &k.TLSConfig{},
				},
			},
			&action{
				tp: &http2.Transport{
					TLSClientConfig: &tls.Config{
						NextProtos:   []string{http2.NextProtoTLS},
						RootCAs:      newSystemPool(),
						ClientCAs:    newSystemPool(),
						Certificates: []tls.Certificate{},
					},
				},
				err: nil,
			},
		),
		gen(
			"specify invalid TLSConfig values",
			&condition{
				spec: &k.HTTP2TransportConfig{
					TLSConfig: &k.TLSConfig{
						ClientAuth: 99999,
					},
				},
			},
			&action{
				tp: nil,
				err: &er.Error{
					Package:     network.ErrPkg,
					Type:        network.ErrTypeTransport,
					Description: network.ErrDscNewTransport,
				},
			},
		),
		gen(
			"specify valid DialConfig values",
			&condition{
				spec: &k.HTTP2TransportConfig{
					DialConfig: &k.DialConfig{},
				},
			},
			&action{
				tp:  &http2.Transport{},
				err: nil,
			},
		),
		gen(
			"specify valid DialConfig values with TLS",
			&condition{
				spec: &k.HTTP2TransportConfig{
					DialConfig: &k.DialConfig{
						TLSConfig: &k.TLSConfig{},
					},
				},
			},
			&action{
				tp:  &http2.Transport{},
				err: nil,
			},
		),
		gen(
			"specify invalid DialConfig values",
			&condition{
				spec: &k.HTTP2TransportConfig{
					DialConfig: &k.DialConfig{
						LocalAddress: "tcp://INVALID_ADDRESS",
					},
				},
			},
			&action{
				tp: nil,
				err: &er.Error{
					Package:     network.ErrPkg,
					Type:        network.ErrTypeTransport,
					Description: network.ErrDscNewTransport,
				},
			},
		),
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			tp, err := network.HTTP2Transport(tt.C.spec)
			testutil.Diff(t, tt.A.err, err, cmpopts.EquateErrors())
			if err != nil {
				return
			}

			opts := []cmp.Option{
				cmpopts.IgnoreUnexported(http2.Transport{}, tls.Config{}),
				cmpopts.IgnoreFields(http2.Transport{}, "DialTLSContext"),
			}
			testutil.Diff(t, tt.A.tp, tp, opts...)

			// Check that the dialer is applied.
			// testutil.Diff(t, true, tp.DialContext != nil)
		})
	}
}

func TestHTTP3Transport(t *testing.T) {
	type condition struct {
		spec *k.HTTP3TransportConfig
	}

	type action struct {
		rt  *http3.Transport
		err error
	}

	newSystemPool := func() *x509.CertPool {
		pool, err := x509.SystemCertPool()
		if err != nil {
			panic(err) // Bad test environment.
		}
		return pool
	}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"input nil",
			&condition{
				spec: nil,
			},
			&action{
				rt:  &http3.Transport{},
				err: nil,
			},
		),
		gen(
			"input zero value spec",
			&condition{
				spec: &k.HTTP3TransportConfig{},
			},
			&action{
				rt:  &http3.Transport{},
				err: nil,
			},
		),
		gen(
			"valid config",
			&condition{
				spec: &k.HTTP3TransportConfig{
					QuicConfig: &k.QuicConfig{
						Versions:                       []k.QuicConfig_Version{k.QuicConfig_Version1, k.QuicConfig_Version2},
						HandshakeIdleTimeout:           10,
						MaxIdleTimeout:                 11,
						InitialStreamReceiveWindow:     12,
						MaxStreamReceiveWindow:         13,
						InitialConnectionReceiveWindow: 14,
						MaxConnectionReceiveWindow:     15,
						MaxIncomingStreams:             16,
						MaxIncomingUniStreams:          17,
						KeepAlivePeriod:                18,
						DisablePathMTUDiscovery:        true,
						Allow0RTT:                      true,
						EnableDatagrams:                true,
					},
					DisableCompression:     true,
					EnableDatagrams:        true,
					MaxResponseHeaderBytes: 1,
				},
			},
			&action{
				rt: &http3.Transport{
					QUICConfig: &quic.Config{
						Versions:                       []quic.Version{quic.Version1, quic.Version2},
						HandshakeIdleTimeout:           time.Millisecond * 10,
						MaxIdleTimeout:                 time.Millisecond * 11,
						InitialStreamReceiveWindow:     12,
						MaxStreamReceiveWindow:         13,
						InitialConnectionReceiveWindow: 14,
						MaxConnectionReceiveWindow:     15,
						MaxIncomingStreams:             16,
						MaxIncomingUniStreams:          17,
						KeepAlivePeriod:                time.Millisecond * 18,
						DisablePathMTUDiscovery:        true,
						Allow0RTT:                      true,
						EnableDatagrams:                true,
					},
					DisableCompression:     true,
					EnableDatagrams:        true,
					MaxResponseHeaderBytes: 1,
				},
				err: nil,
			},
		),
		gen(
			"valid TLSConfig",
			&condition{
				spec: &k.HTTP3TransportConfig{
					TLSConfig: &k.TLSConfig{},
				},
			},
			&action{
				rt: &http3.Transport{
					TLSClientConfig: &tls.Config{
						RootCAs:      newSystemPool(),
						ClientCAs:    newSystemPool(),
						Certificates: []tls.Certificate{},
					},
				},
				err: nil,
			},
		),
		gen(
			"invalid TLSConfig",
			&condition{
				spec: &k.HTTP3TransportConfig{
					TLSConfig: &k.TLSConfig{
						ClientAuth: 99999,
					},
				},
			},
			&action{
				rt: nil,
				err: &er.Error{
					Package:     network.ErrPkg,
					Type:        network.ErrTypeTransport,
					Description: network.ErrDscNewTransport,
				},
			},
		),
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			rt, err := network.HTTP3Transport(tt.C.spec)
			testutil.Diff(t, tt.A.err, err, cmpopts.EquateErrors())
			if err != nil {
				return
			}

			opts := []cmp.Option{
				cmpopts.IgnoreUnexported(tls.Config{}),
				cmpopts.IgnoreUnexported(http3.Transport{}),
			}
			testutil.Diff(t, tt.A.rt, rt, opts...)
		})
	}
}
