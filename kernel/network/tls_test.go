// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package network

import (
	"crypto/tls"
	"crypto/x509"
	"testing"

	k "github.com/aileron-gateway/aileron-gateway/apis/kernel"
	"github.com/aileron-gateway/aileron-gateway/kernel/er"
	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestTLSConfig(t *testing.T) {
	type condition struct {
		spec *k.TLSConfig
	}

	type action struct {
		config *tls.Config
		err    error
	}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"nil",
			&condition{
				spec: nil,
			},
			&action{
				err: nil,
			},
		),
		gen(
			"root CA error",
			&condition{
				spec: &k.TLSConfig{
					RootCAs: []string{testDir + "ut/core/utilhttp/not-exists.pem"},
				},
			},
			&action{
				err: &er.Error{
					Package:     ErrPkg,
					Type:        ErrTypeTLS,
					Description: ErrDscTLS,
				},
			},
		),
		gen(
			"client CA error",
			&condition{
				spec: &k.TLSConfig{
					ClientCAs: []string{testDir + "ut/core/utilhttp/not-exists.pem"},
				},
			},
			&action{
				err: &er.Error{
					Package:     ErrPkg,
					Type:        ErrTypeTLS,
					Description: ErrDscTLS,
				},
			},
		),
		gen(
			"client auth error",
			&condition{
				spec: &k.TLSConfig{
					ClientAuth: k.ClientAuthType(999), // Invalid value.
				},
			},
			&action{
				err: &er.Error{
					Package:     ErrPkg,
					Type:        ErrTypeTLS,
					Description: ErrDscTLS,
				},
			},
		),
		gen(
			"renegotiation error",
			&condition{
				spec: &k.TLSConfig{
					Renegotiation: k.RenegotiationSupport(999), // Invalid value.
				},
			},
			&action{
				err: &er.Error{
					Package:     ErrPkg,
					Type:        ErrTypeTLS,
					Description: ErrDscTLS,
				},
			},
		),
		gen(
			"certificate error",
			&condition{
				spec: &k.TLSConfig{
					CertKeyPairs: []*k.CertKeyPair{
						{
							CertFile: testDir + "ut/core/utilhttp/not-exists.crt",
							KeyFile:  testDir + "ut/core/utilhttp/not-exists.key",
						},
					},
				},
			},
			&action{
				err: &er.Error{
					Package:     ErrPkg,
					Type:        ErrTypeTLS,
					Description: ErrDscTLS,
				},
			},
		),
		gen(
			"valid config",
			&condition{
				spec: &k.TLSConfig{
					CertKeyPairs: []*k.CertKeyPair{
						{
							CertFile: testDir + "ut/core/utilhttp/test.crt",
							KeyFile:  testDir + "ut/core/utilhttp/test.key",
						},
					},
					RootCAs:                     nil,
					ServerName:                  "test-server",
					ClientAuth:                  1,
					ClientCAs:                   nil,
					NextProtos:                  []string{"foo", "bar"},
					InsecureSkipVerify:          true,
					TLSCiphers:                  nil,
					SessionTicketsDisabled:      true,
					MinVersion:                  1,
					MaxVersion:                  2,
					CurvePreferences:            nil,
					DynamicRecordSizingDisabled: true,
					Renegotiation:               1,
				},
			},
			&action{
				config: &tls.Config{
					Certificates:                nil,
					RootCAs:                     x509.NewCertPool(),
					NextProtos:                  []string{"foo", "bar"},
					ServerName:                  "test-server",
					ClientAuth:                  1,
					ClientCAs:                   x509.NewCertPool(),
					InsecureSkipVerify:          true,
					CipherSuites:                nil,
					SessionTicketsDisabled:      true,
					MinVersion:                  1,
					MaxVersion:                  2,
					CurvePreferences:            nil,
					DynamicRecordSizingDisabled: true,
					Renegotiation:               1,
				},
			},
		),
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			config, err := TLSConfig(tt.C.spec)
			testutil.Diff(t, tt.A.err, err, cmpopts.EquateErrors())
			if config == nil {
				return
			}

			opts := []cmp.Option{
				cmpopts.IgnoreUnexported(tls.Config{}),
				cmpopts.IgnoreFields(tls.Config{}, "Certificates"),
				cmpopts.IgnoreFields(tls.Config{}, "RootCAs", "ClientCAs"),
			}
			testutil.Diff(t, tt.A.config, config, opts...)
			testutil.Diff(t, len(tt.C.spec.CertKeyPairs), len(config.Certificates)) // TODO: Check better way.
		})
	}
}

func TestCiphers(t *testing.T) {
	type condition struct {
		cs []k.TLSCipher
	}

	type action struct {
		cs []uint16
	}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"nil",
			&condition{
				cs: []k.TLSCipher{},
			},
			&action{
				cs: nil,
			},
		),
		gen(
			"invalid",
			&condition{
				cs: []k.TLSCipher{k.TLSCipher(-999)}, // Invalid cipher.
			},
			&action{
				cs: []uint16{},
			},
		),
		gen(
			"TLS_RSA_WITH_RC4_128_SHA",
			&condition{
				cs: []k.TLSCipher{k.TLSCipher_TLS_RSA_WITH_RC4_128_SHA},
			},
			&action{
				cs: []uint16{tls.TLS_RSA_WITH_RC4_128_SHA},
			},
		),
		gen(
			"TLS_RSA_WITH_3DES_EDE_CBC_SHA",
			&condition{
				cs: []k.TLSCipher{k.TLSCipher_TLS_RSA_WITH_3DES_EDE_CBC_SHA},
			},
			&action{
				cs: []uint16{tls.TLS_RSA_WITH_3DES_EDE_CBC_SHA},
			},
		),
		gen(
			"TLS_RSA_WITH_AES_128_CBC_SHA",
			&condition{
				cs: []k.TLSCipher{k.TLSCipher_TLS_RSA_WITH_AES_128_CBC_SHA},
			},
			&action{
				cs: []uint16{tls.TLS_RSA_WITH_AES_128_CBC_SHA},
			},
		),
		gen(
			"TLS_RSA_WITH_AES_256_CBC_SHA",
			&condition{
				cs: []k.TLSCipher{k.TLSCipher_TLS_RSA_WITH_AES_256_CBC_SHA},
			},
			&action{
				cs: []uint16{tls.TLS_RSA_WITH_AES_256_CBC_SHA},
			},
		),
		gen(
			"TLS_RSA_WITH_AES_128_CBC_SHA256",
			&condition{
				cs: []k.TLSCipher{k.TLSCipher_TLS_RSA_WITH_AES_128_CBC_SHA256},
			},
			&action{
				cs: []uint16{tls.TLS_RSA_WITH_AES_128_CBC_SHA256},
			},
		),

		gen(
			"TLS_RSA_WITH_AES_128_GCM_SHA256",
			&condition{
				cs: []k.TLSCipher{k.TLSCipher_TLS_RSA_WITH_AES_128_GCM_SHA256},
			},
			&action{
				cs: []uint16{tls.TLS_RSA_WITH_AES_128_GCM_SHA256},
			},
		),
		gen(
			"TLS_RSA_WITH_AES_256_GCM_SHA384",
			&condition{
				cs: []k.TLSCipher{k.TLSCipher_TLS_RSA_WITH_AES_256_GCM_SHA384},
			},
			&action{
				cs: []uint16{tls.TLS_RSA_WITH_AES_256_GCM_SHA384},
			},
		),
		gen(
			"TLS_ECDHE_ECDSA_WITH_RC4_128_SHA",
			&condition{
				cs: []k.TLSCipher{k.TLSCipher_TLS_ECDHE_ECDSA_WITH_RC4_128_SHA},
			},
			&action{
				cs: []uint16{tls.TLS_ECDHE_ECDSA_WITH_RC4_128_SHA},
			},
		),
		gen(
			"TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA",
			&condition{
				cs: []k.TLSCipher{k.TLSCipher_TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA},
			},
			&action{
				cs: []uint16{tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA},
			},
		),
		gen(
			"TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA",
			&condition{
				cs: []k.TLSCipher{k.TLSCipher_TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA},
			},
			&action{
				cs: []uint16{tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA},
			},
		),
		gen(
			"TLS_ECDHE_RSA_WITH_RC4_128_SHA",
			&condition{
				cs: []k.TLSCipher{k.TLSCipher_TLS_ECDHE_RSA_WITH_RC4_128_SHA},
			},
			&action{
				cs: []uint16{tls.TLS_ECDHE_RSA_WITH_RC4_128_SHA},
			},
		),
		gen(
			"TLS_ECDHE_RSA_WITH_3DES_EDE_CBC_SHA",
			&condition{
				cs: []k.TLSCipher{k.TLSCipher_TLS_ECDHE_RSA_WITH_3DES_EDE_CBC_SHA},
			},
			&action{
				cs: []uint16{tls.TLS_ECDHE_RSA_WITH_3DES_EDE_CBC_SHA},
			},
		),
		gen(
			"TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA",
			&condition{
				cs: []k.TLSCipher{k.TLSCipher_TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA},
			},
			&action{
				cs: []uint16{tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA},
			},
		),
		gen(
			"TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA",
			&condition{
				cs: []k.TLSCipher{k.TLSCipher_TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA},
			},
			&action{
				cs: []uint16{tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA},
			},
		),
		gen(
			"TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256",
			&condition{
				cs: []k.TLSCipher{k.TLSCipher_TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256},
			},
			&action{
				cs: []uint16{tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256},
			},
		),
		gen(
			"TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256",
			&condition{
				cs: []k.TLSCipher{k.TLSCipher_TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256},
			},
			&action{
				cs: []uint16{tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256},
			},
		),
		gen(
			"TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256",
			&condition{
				cs: []k.TLSCipher{k.TLSCipher_TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256},
			},
			&action{
				cs: []uint16{tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256},
			},
		),
		gen(
			"TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256",
			&condition{
				cs: []k.TLSCipher{k.TLSCipher_TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256},
			},
			&action{
				cs: []uint16{tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256},
			},
		),
		gen(
			"TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384",
			&condition{
				cs: []k.TLSCipher{k.TLSCipher_TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384},
			},
			&action{
				cs: []uint16{tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384},
			},
		),
		gen(
			"TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384",
			&condition{
				cs: []k.TLSCipher{k.TLSCipher_TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384},
			},
			&action{
				cs: []uint16{tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384},
			},
		),
		gen(
			"TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305_SHA256",
			&condition{
				cs: []k.TLSCipher{k.TLSCipher_TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305_SHA256},
			},
			&action{
				cs: []uint16{tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305_SHA256},
			},
		),
		gen(
			"TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305_SHA256",
			&condition{
				cs: []k.TLSCipher{k.TLSCipher_TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305_SHA256},
			},
			&action{
				cs: []uint16{tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305_SHA256},
			},
		),
		gen(
			"TLS_AES_128_GCM_SHA256",
			&condition{
				cs: []k.TLSCipher{k.TLSCipher_TLS_AES_128_GCM_SHA256},
			},
			&action{
				cs: []uint16{tls.TLS_AES_128_GCM_SHA256},
			},
		),
		gen(
			"TLS_AES_256_GCM_SHA384",
			&condition{
				cs: []k.TLSCipher{k.TLSCipher_TLS_AES_256_GCM_SHA384},
			},
			&action{
				cs: []uint16{tls.TLS_AES_256_GCM_SHA384},
			},
		),
		gen(
			"TLS_CHACHA20_POLY1305_SHA256",
			&condition{
				cs: []k.TLSCipher{k.TLSCipher_TLS_CHACHA20_POLY1305_SHA256},
			},
			&action{
				cs: []uint16{tls.TLS_CHACHA20_POLY1305_SHA256},
			},
		),
		gen(
			"TLS_FALLBACK_SCSV",
			&condition{
				cs: []k.TLSCipher{k.TLSCipher_TLS_FALLBACK_SCSV},
			},
			&action{
				cs: []uint16{tls.TLS_FALLBACK_SCSV},
			},
		),
		gen(
			"TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305",
			&condition{
				cs: []k.TLSCipher{k.TLSCipher_TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305},
			},
			&action{
				cs: []uint16{tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305},
			},
		),
		gen(
			"TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305",
			&condition{
				cs: []k.TLSCipher{k.TLSCipher_TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305},
			},
			&action{
				cs: []uint16{tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305},
			},
		),
		gen(
			"input 2 ciphers",
			&condition{
				cs: []k.TLSCipher{
					k.TLSCipher_TLS_RSA_WITH_RC4_128_SHA,
					k.TLSCipher_TLS_RSA_WITH_3DES_EDE_CBC_SHA,
				},
			},
			&action{
				cs: []uint16{
					tls.TLS_RSA_WITH_RC4_128_SHA,
					tls.TLS_RSA_WITH_3DES_EDE_CBC_SHA,
				},
			},
		),
		gen(
			"input 3 ciphers",
			&condition{
				cs: []k.TLSCipher{
					k.TLSCipher_TLS_RSA_WITH_RC4_128_SHA,
					k.TLSCipher_TLS_RSA_WITH_3DES_EDE_CBC_SHA,
					k.TLSCipher_TLS_RSA_WITH_AES_128_CBC_SHA,
				},
			},
			&action{
				cs: []uint16{
					tls.TLS_RSA_WITH_RC4_128_SHA,
					tls.TLS_RSA_WITH_3DES_EDE_CBC_SHA,
					tls.TLS_RSA_WITH_AES_128_CBC_SHA,
				},
			},
		),
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			cs := tlsCiphers(tt.C.cs)
			testutil.Diff(t, tt.A.cs, cs)
		})
	}
}

func TestCurveIDs(t *testing.T) {
	type condition struct {
		ids []k.CurveID
	}

	type action struct {
		ids []tls.CurveID
	}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"nil",
			&condition{
				ids: []k.CurveID{},
			},
			&action{
				ids: nil,
			},
		),
		gen(
			"invalid",
			&condition{
				ids: []k.CurveID{
					k.CurveID(999), // Invalid curve id.
				},
			},
			&action{
				ids: []tls.CurveID{},
			},
		),
		gen(
			"P256",
			&condition{
				ids: []k.CurveID{
					k.CurveID_CurveP256,
				},
			},
			&action{
				ids: []tls.CurveID{
					tls.CurveP256,
				},
			},
		),
		gen(
			"P384",
			&condition{
				ids: []k.CurveID{
					k.CurveID_CurveP384,
				},
			},
			&action{
				ids: []tls.CurveID{
					tls.CurveP384,
				},
			},
		),
		gen(
			"P521",
			&condition{
				ids: []k.CurveID{
					k.CurveID_CurveP521,
				},
			},
			&action{
				ids: []tls.CurveID{
					tls.CurveP521,
				},
			},
		),
		gen(
			"X25519",
			&condition{
				ids: []k.CurveID{
					k.CurveID_X25519,
				},
			},
			&action{
				ids: []tls.CurveID{
					tls.X25519,
				},
			},
		),
		gen(
			"all",
			&condition{
				ids: []k.CurveID{
					k.CurveID_CurveP256,
					k.CurveID_CurveP384,
					k.CurveID_CurveP521,
					k.CurveID_X25519,
				},
			},
			&action{
				ids: []tls.CurveID{
					tls.CurveP256,
					tls.CurveP384,
					tls.CurveP521,
					tls.X25519,
				},
			},
		),
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			ids := curveIDs(tt.C.ids)
			testutil.Diff(t, tt.A.ids, ids)
		})
	}
}
