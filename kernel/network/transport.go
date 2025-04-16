// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package network

import (
	"cmp"
	"context"
	"crypto/tls"
	"net"
	"net/http"
	"slices"
	"time"

	k "github.com/aileron-gateway/aileron-gateway/apis/kernel"
	"github.com/aileron-gateway/aileron-gateway/kernel/er"
	"github.com/quic-go/quic-go"
	"github.com/quic-go/quic-go/http3"
	"golang.org/x/net/http2"
)

var (
	DefaultHTTPTransport = &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   10 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		DisableCompression:    true,
		ForceAttemptHTTP2:     false, // Save memory.
		MaxIdleConns:          0,     // No limit.
		MaxIdleConnsPerHost:   1024,  // If 0, 2 is used.
		MaxConnsPerHost:       0,     // No limit.
		IdleConnTimeout:       60 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
	DefaultHTTP2Transport = &http2.Transport{}
	DefaultHTTP3Transport = &http3.Transport{}
)

// HTTPTransport returns a new http.Transport from the given HTTP2TransportConfig.
// Default transport will be returned when nil spec is given.
func HTTPTransport(spec *k.HTTPTransportConfig) (*http.Transport, error) {
	if spec == nil {
		return DefaultHTTPTransport, nil
	}

	tlsConfig, err := TLSConfig(spec.TLSConfig)
	if err != nil {
		return nil, (&er.Error{
			Package:     ErrPkg,
			Type:        ErrTypeTransport,
			Description: ErrDscNewTransport,
			Detail:      "target transport is HTTP.",
		}).Wrap(err)
	}

	h := make(http.Header, len(spec.ProxyConnectHeaders))
	for k, v := range spec.ProxyConnectHeaders {
		h.Set(k, v)
	}
	transport := &http.Transport{
		Proxy:                  http.ProxyFromEnvironment,
		TLSClientConfig:        tlsConfig,
		TLSHandshakeTimeout:    time.Millisecond * time.Duration(spec.TLSHandshakeTimeout),
		DisableKeepAlives:      spec.DisableKeepAlives,
		DisableCompression:     spec.DisableCompression,
		MaxIdleConns:           int(spec.MaxIdleConns),
		MaxIdleConnsPerHost:    cmp.Or(int(spec.MaxIdleConnsPerHost), 1024),
		MaxConnsPerHost:        int(spec.MaxConnsPerHost),
		IdleConnTimeout:        time.Millisecond * time.Duration(spec.IdleConnTimeout),
		ResponseHeaderTimeout:  time.Millisecond * time.Duration(spec.ResponseHeaderTimeout),
		ExpectContinueTimeout:  time.Millisecond * time.Duration(spec.ExpectContinueTimeout),
		ProxyConnectHeader:     h,
		MaxResponseHeaderBytes: spec.MaxResponseHeaderBytes,
		WriteBufferSize:        int(spec.WriteBufferSize),
		ReadBufferSize:         int(spec.ReadBufferSize),
		ForceAttemptHTTP2:      false,
	}

	if !spec.AllowHTTP2 {
		// Disable http2.
		transport.TLSNextProto = make(map[string]func(authority string, c *tls.Conn) http.RoundTripper, 0)
	}

	if spec.DialConfig != nil {
		if spec.AllowHTTP2 {
			transport.ForceAttemptHTTP2 = true
		}
		spec.DialConfig.TLSConfig = cmp.Or(spec.DialConfig.TLSConfig, spec.TLSConfig)
		dialer, err := NewDialerFromSpec(spec.DialConfig)
		if err != nil {
			return nil, (&er.Error{
				Package:     ErrPkg,
				Type:        ErrTypeTransport,
				Description: ErrDscNewTransport,
				Detail:      "target transport is HTTP.",
			}).Wrap(err)
		}
		transport.DialContext = dialer.DialContext
		transport.DialTLSContext = dialer.DialContext
	}

	return transport, nil
}

// HTTP2Transport returns a new http2.Transport from the given HTTP2TransportConfig.
// Default transport will be returned when nil spec is given.
func HTTP2Transport(spec *k.HTTP2TransportConfig) (*http2.Transport, error) {
	if spec == nil {
		return DefaultHTTP2Transport, nil
	}

	tlsConfig, err := TLSConfig(spec.TLSConfig)
	if err != nil {
		return nil, (&er.Error{
			Package:     ErrPkg,
			Type:        ErrTypeTransport,
			Description: ErrDscNewTransport,
			Detail:      "target transport is HTTP2.",
		}).Wrap(err)
	}
	if tlsConfig != nil && !slices.Contains(tlsConfig.NextProtos, http2.NextProtoTLS) {
		tlsConfig.NextProtos = append([]string{http2.NextProtoTLS}, tlsConfig.NextProtos...)
	}

	transport := &http2.Transport{
		TLSClientConfig:            tlsConfig,
		DisableCompression:         spec.DisableCompression,
		AllowHTTP:                  spec.AllowHTTP, // Must set http2.NextProtoTLS in the TLSConfig of Dialer.
		MaxHeaderListSize:          spec.MaxHeaderListSize,
		MaxReadFrameSize:           spec.MaxReadFrameSize,
		MaxDecoderHeaderTableSize:  spec.MaxDecoderHeaderTableSize,
		MaxEncoderHeaderTableSize:  spec.MaxEncoderHeaderTableSize,
		StrictMaxConcurrentStreams: spec.StrictMaxConcurrentStreams,
		IdleConnTimeout:            time.Millisecond * time.Duration(spec.IdleConnTimeout),
		ReadIdleTimeout:            time.Millisecond * time.Duration(spec.ReadIdleTimeout),
		PingTimeout:                time.Millisecond * time.Duration(spec.PingTimeout),
		WriteByteTimeout:           time.Millisecond * time.Duration(spec.WriteByteTimeout),
	}

	if spec.MultiIPConnPool {
		pool := &http2ConnPool{
			t:               transport,
			conns:           map[string]*hostConns{},
			connMap:         map[*http2.ClientConn]*hostConns{},
			resolveInterval: time.Millisecond * time.Duration(spec.MinLookupInterval),
		}
		transport.ConnPool = pool
	}

	if (spec.AllowHTTP || spec.MultiIPConnPool) && spec.DialConfig == nil {
		spec.DialConfig = &k.DialConfig{}
	}
	if spec.DialConfig != nil {
		tc := spec.DialConfig.TLSConfig
		tc = cmp.Or(tc, spec.TLSConfig)
		if tc != nil && !slices.Contains(tc.NextProtos, http2.NextProtoTLS) {
			tc.NextProtos = append([]string{http2.NextProtoTLS}, tc.NextProtos...)
		}
		dialer, err := NewDialerFromSpec(spec.DialConfig)
		if err != nil {
			return nil, (&er.Error{
				Package:     ErrPkg,
				Type:        ErrTypeTransport,
				Description: ErrDscNewTransport,
				Detail:      "target transport is HTTP2.",
			}).Wrap(err)
		}
		transport.DialTLSContext = func(ctx context.Context, network, addr string, _ *tls.Config) (net.Conn, error) {
			return dialer.DialContext(ctx, network, addr)
		}
	}

	return transport, nil
}

// HTTP3Transport returns a new http3.Transport from the given HTTP3TransportConfig.
// Zero valued http3.Transport will be returned when nil spec is given.
func HTTP3Transport(spec *k.HTTP3TransportConfig) (*http3.Transport, error) {
	if spec == nil {
		return DefaultHTTP3Transport, nil
	}

	tlsConfig, err := TLSConfig(spec.TLSConfig)
	if err != nil {
		return nil, (&er.Error{
			Package:     ErrPkg,
			Type:        ErrTypeTransport,
			Description: ErrDscNewTransport,
			Detail:      "target transport is HTTP3.",
		}).Wrap(err)
	}

	// Function QuicConfig does not return any error for now.
	quicConfig, _ := QuicConfig(spec.QuicConfig)

	return &http3.Transport{
		TLSClientConfig:        tlsConfig,
		QUICConfig:             quicConfig,
		DisableCompression:     spec.DisableCompression,
		EnableDatagrams:        spec.EnableDatagrams,
		MaxResponseHeaderBytes: spec.MaxResponseHeaderBytes,
	}, nil
}

// QuicConfig returns a new quic.Config from the given core.k.QuicConfig.
// The second returned value error is always nil for now.
// It is placed for future extension.
func QuicConfig(spec *k.QuicConfig) (*quic.Config, error) {
	if spec == nil {
		return nil, nil
	}

	vs := make([]quic.Version, len(spec.Versions))
	for i, v := range spec.Versions {
		switch v {
		case k.QuickVersion_Version1:
			vs[i] = quic.Version1
		case k.QuickVersion_Version2:
			vs[i] = quic.Version2
		}
	}

	return &quic.Config{
		Versions:                       vs,
		HandshakeIdleTimeout:           time.Millisecond * time.Duration(spec.HandshakeIdleTimeout),
		MaxIdleTimeout:                 time.Millisecond * time.Duration(spec.MaxIdleTimeout),
		InitialStreamReceiveWindow:     spec.InitialStreamReceiveWindow,
		MaxStreamReceiveWindow:         spec.MaxStreamReceiveWindow,
		InitialConnectionReceiveWindow: spec.InitialConnectionReceiveWindow,
		MaxConnectionReceiveWindow:     spec.MaxConnectionReceiveWindow,
		MaxIncomingStreams:             spec.MaxIncomingStreams,
		MaxIncomingUniStreams:          spec.MaxIncomingUniStreams,
		KeepAlivePeriod:                time.Millisecond * time.Duration(spec.KeepAlivePeriod),
		DisablePathMTUDiscovery:        spec.DisablePathMTUDiscovery,
		Allow0RTT:                      spec.Allow0RTT,
		EnableDatagrams:                spec.EnableDatagrams,
	}, nil
}
