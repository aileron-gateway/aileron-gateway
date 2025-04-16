// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package network_test

import (
	stdcmp "cmp"
	"crypto/tls"
	"crypto/x509"
	"os"
	"testing"

	k "github.com/aileron-gateway/aileron-gateway/apis/kernel"
	"github.com/aileron-gateway/aileron-gateway/kernel/er"
	"github.com/aileron-gateway/aileron-gateway/kernel/network"
	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

// testDir is the path to the test data.
// This path can be changed by the environmental variable.
var testDir = stdcmp.Or(os.Getenv("TEST_DIR"), "../../test/")

func TestTLSConfig(t *testing.T) {
	type condition struct {
		spec *k.TLSConfig
	}

	type action struct {
		config *tls.Config
		err    error
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	cndInvalidRootCA := tb.Condition("invalid root CA", "input invalid root ca file path")
	cndInvalidClientCA := tb.Condition("invalid client CA", "input invalid client ca file path")
	cndInvalidClientAuth := tb.Condition("invalid client auth", "input invalid value for client auth")
	cndInvalidRenegotiation := tb.Condition("invalid renegotiation", "input invalid value for renegotiation")
	cndInvalidCertKey := tb.Condition("invalid cert key pair", "input invalid value for cert key pairs")
	cndValidConfig := tb.Condition("valid config", "input valid tls config")
	actCheckNoError := tb.Action("no error", "check that there is no error")
	actCheckError := tb.Action("error", "check that there is an error")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"nil",
			[]string{},
			[]string{actCheckNoError},
			&condition{
				spec: nil,
			},
			&action{
				err: nil,
			},
		),
		gen(
			"root CA error",
			[]string{cndInvalidRootCA},
			[]string{actCheckError},
			&condition{
				spec: &k.TLSConfig{
					RootCAs: []string{testDir + "ut/core/utilhttp/not-exists.pem"},
				},
			},
			&action{
				err: &er.Error{
					Package:     network.ErrPkg,
					Type:        network.ErrTypeTLS,
					Description: network.ErrDscTLS,
				},
			},
		),
		gen(
			"client CA error",
			[]string{cndInvalidClientCA},
			[]string{actCheckError},
			&condition{
				spec: &k.TLSConfig{
					ClientCAs: []string{testDir + "ut/core/utilhttp/not-exists.pem"},
				},
			},
			&action{
				err: &er.Error{
					Package:     network.ErrPkg,
					Type:        network.ErrTypeTLS,
					Description: network.ErrDscTLS,
				},
			},
		),
		gen(
			"client auth error",
			[]string{cndInvalidClientAuth},
			[]string{actCheckError},
			&condition{
				spec: &k.TLSConfig{
					ClientAuth: k.ClientAuthType(999), // Invalid value.
				},
			},
			&action{
				err: &er.Error{
					Package:     network.ErrPkg,
					Type:        network.ErrTypeTLS,
					Description: network.ErrDscTLS,
				},
			},
		),
		gen(
			"renegotiation error",
			[]string{cndInvalidRenegotiation},
			[]string{actCheckError},
			&condition{
				spec: &k.TLSConfig{
					Renegotiation: k.RenegotiationSupport(999), // Invalid value.
				},
			},
			&action{
				err: &er.Error{
					Package:     network.ErrPkg,
					Type:        network.ErrTypeTLS,
					Description: network.ErrDscTLS,
				},
			},
		),
		gen(
			"certificate error",
			[]string{cndInvalidCertKey},
			[]string{actCheckError},
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
					Package:     network.ErrPkg,
					Type:        network.ErrTypeTLS,
					Description: network.ErrDscTLS,
				},
			},
		),
		gen(
			"valid config",
			[]string{cndValidConfig},
			[]string{actCheckNoError},
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
					RootCAsIgnoreSystemCerts:    true,
					ClientCAsIgnoreSystemCerts:  true,
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
					CurvePreferences:            []tls.CurveID{},
					DynamicRecordSizingDisabled: true,
					Renegotiation:               1,
				},
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			config, err := network.TLSConfig(tt.C().spec)
			testutil.Diff(t, tt.A().err, err, cmpopts.EquateErrors())
			if config == nil {
				return
			}

			opts := []cmp.Option{
				cmpopts.IgnoreUnexported(tls.Config{}),
				cmpopts.IgnoreFields(tls.Config{}, "Certificates"),
			}
			testutil.Diff(t, tt.A().config, config, opts...)
			testutil.Diff(t, len(tt.C().spec.CertKeyPairs), len(config.Certificates)) // TODO: Check better way.
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

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	cndInputNil := tb.Condition("nil", "")
	cndInputInvalid := tb.Condition("invalid", "")
	actCheckIncluded := tb.Action("included", "check that input cipher is included in the returned list")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"nil",
			[]string{cndInputNil},
			[]string{},
			&condition{
				cs: []k.TLSCipher{},
			},
			&action{
				cs: nil,
			},
		),
		gen(
			"invalid",
			[]string{cndInputInvalid},
			[]string{},
			&condition{
				cs: []k.TLSCipher{k.TLSCipher(-999)}, // Invalid cipher.
			},
			&action{
				cs: []uint16{},
			},
		),
		gen(
			"TLS_RSA_WITH_RC4_128_SHA",
			[]string{},
			[]string{actCheckIncluded},
			&condition{
				cs: []k.TLSCipher{k.TLSCipher_TLS_RSA_WITH_RC4_128_SHA},
			},
			&action{
				cs: []uint16{tls.TLS_RSA_WITH_RC4_128_SHA},
			},
		),
		gen(
			"TLS_RSA_WITH_3DES_EDE_CBC_SHA",
			[]string{},
			[]string{actCheckIncluded},
			&condition{
				cs: []k.TLSCipher{k.TLSCipher_TLS_RSA_WITH_3DES_EDE_CBC_SHA},
			},
			&action{
				cs: []uint16{tls.TLS_RSA_WITH_3DES_EDE_CBC_SHA},
			},
		),
		gen(
			"TLS_RSA_WITH_AES_128_CBC_SHA",
			[]string{},
			[]string{actCheckIncluded},
			&condition{
				cs: []k.TLSCipher{k.TLSCipher_TLS_RSA_WITH_AES_128_CBC_SHA},
			},
			&action{
				cs: []uint16{tls.TLS_RSA_WITH_AES_128_CBC_SHA},
			},
		),
		gen(
			"TLS_RSA_WITH_AES_256_CBC_SHA",
			[]string{},
			[]string{actCheckIncluded},
			&condition{
				cs: []k.TLSCipher{k.TLSCipher_TLS_RSA_WITH_AES_256_CBC_SHA},
			},
			&action{
				cs: []uint16{tls.TLS_RSA_WITH_AES_256_CBC_SHA},
			},
		),
		gen(
			"TLS_RSA_WITH_AES_128_CBC_SHA256",
			[]string{},
			[]string{actCheckIncluded},
			&condition{
				cs: []k.TLSCipher{k.TLSCipher_TLS_RSA_WITH_AES_128_CBC_SHA256},
			},
			&action{
				cs: []uint16{tls.TLS_RSA_WITH_AES_128_CBC_SHA256},
			},
		),

		gen(
			"TLS_RSA_WITH_AES_128_GCM_SHA256",
			[]string{},
			[]string{actCheckIncluded},
			&condition{
				cs: []k.TLSCipher{k.TLSCipher_TLS_RSA_WITH_AES_128_GCM_SHA256},
			},
			&action{
				cs: []uint16{tls.TLS_RSA_WITH_AES_128_GCM_SHA256},
			},
		),
		gen(
			"TLS_RSA_WITH_AES_256_GCM_SHA384",
			[]string{},
			[]string{actCheckIncluded},
			&condition{
				cs: []k.TLSCipher{k.TLSCipher_TLS_RSA_WITH_AES_256_GCM_SHA384},
			},
			&action{
				cs: []uint16{tls.TLS_RSA_WITH_AES_256_GCM_SHA384},
			},
		),
		gen(
			"TLS_ECDHE_ECDSA_WITH_RC4_128_SHA",
			[]string{},
			[]string{actCheckIncluded},
			&condition{
				cs: []k.TLSCipher{k.TLSCipher_TLS_ECDHE_ECDSA_WITH_RC4_128_SHA},
			},
			&action{
				cs: []uint16{tls.TLS_ECDHE_ECDSA_WITH_RC4_128_SHA},
			},
		),
		gen(
			"TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA",
			[]string{},
			[]string{actCheckIncluded},
			&condition{
				cs: []k.TLSCipher{k.TLSCipher_TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA},
			},
			&action{
				cs: []uint16{tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA},
			},
		),
		gen(
			"TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA",
			[]string{},
			[]string{actCheckIncluded},
			&condition{
				cs: []k.TLSCipher{k.TLSCipher_TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA},
			},
			&action{
				cs: []uint16{tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA},
			},
		),
		gen(
			"TLS_ECDHE_RSA_WITH_RC4_128_SHA",
			[]string{},
			[]string{actCheckIncluded},
			&condition{
				cs: []k.TLSCipher{k.TLSCipher_TLS_ECDHE_RSA_WITH_RC4_128_SHA},
			},
			&action{
				cs: []uint16{tls.TLS_ECDHE_RSA_WITH_RC4_128_SHA},
			},
		),
		gen(
			"TLS_ECDHE_RSA_WITH_3DES_EDE_CBC_SHA",
			[]string{},
			[]string{actCheckIncluded},
			&condition{
				cs: []k.TLSCipher{k.TLSCipher_TLS_ECDHE_RSA_WITH_3DES_EDE_CBC_SHA},
			},
			&action{
				cs: []uint16{tls.TLS_ECDHE_RSA_WITH_3DES_EDE_CBC_SHA},
			},
		),
		gen(
			"TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA",
			[]string{},
			[]string{actCheckIncluded},
			&condition{
				cs: []k.TLSCipher{k.TLSCipher_TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA},
			},
			&action{
				cs: []uint16{tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA},
			},
		),
		gen(
			"TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA",
			[]string{},
			[]string{actCheckIncluded},
			&condition{
				cs: []k.TLSCipher{k.TLSCipher_TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA},
			},
			&action{
				cs: []uint16{tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA},
			},
		),
		gen(
			"TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256",
			[]string{},
			[]string{actCheckIncluded},
			&condition{
				cs: []k.TLSCipher{k.TLSCipher_TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256},
			},
			&action{
				cs: []uint16{tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256},
			},
		),
		gen(
			"TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256",
			[]string{},
			[]string{actCheckIncluded},
			&condition{
				cs: []k.TLSCipher{k.TLSCipher_TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256},
			},
			&action{
				cs: []uint16{tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256},
			},
		),
		gen(
			"TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256",
			[]string{},
			[]string{actCheckIncluded},
			&condition{
				cs: []k.TLSCipher{k.TLSCipher_TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256},
			},
			&action{
				cs: []uint16{tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256},
			},
		),
		gen(
			"TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256",
			[]string{},
			[]string{actCheckIncluded},
			&condition{
				cs: []k.TLSCipher{k.TLSCipher_TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256},
			},
			&action{
				cs: []uint16{tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256},
			},
		),
		gen(
			"TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384",
			[]string{},
			[]string{actCheckIncluded},
			&condition{
				cs: []k.TLSCipher{k.TLSCipher_TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384},
			},
			&action{
				cs: []uint16{tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384},
			},
		),
		gen(
			"TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384",
			[]string{},
			[]string{actCheckIncluded},
			&condition{
				cs: []k.TLSCipher{k.TLSCipher_TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384},
			},
			&action{
				cs: []uint16{tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384},
			},
		),
		gen(
			"TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305_SHA256",
			[]string{},
			[]string{actCheckIncluded},
			&condition{
				cs: []k.TLSCipher{k.TLSCipher_TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305_SHA256},
			},
			&action{
				cs: []uint16{tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305_SHA256},
			},
		),
		gen(
			"TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305_SHA256",
			[]string{},
			[]string{actCheckIncluded},
			&condition{
				cs: []k.TLSCipher{k.TLSCipher_TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305_SHA256},
			},
			&action{
				cs: []uint16{tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305_SHA256},
			},
		),
		gen(
			"TLS_AES_128_GCM_SHA256",
			[]string{},
			[]string{actCheckIncluded},
			&condition{
				cs: []k.TLSCipher{k.TLSCipher_TLS_AES_128_GCM_SHA256},
			},
			&action{
				cs: []uint16{tls.TLS_AES_128_GCM_SHA256},
			},
		),
		gen(
			"TLS_AES_256_GCM_SHA384",
			[]string{},
			[]string{actCheckIncluded},
			&condition{
				cs: []k.TLSCipher{k.TLSCipher_TLS_AES_256_GCM_SHA384},
			},
			&action{
				cs: []uint16{tls.TLS_AES_256_GCM_SHA384},
			},
		),
		gen(
			"TLS_CHACHA20_POLY1305_SHA256",
			[]string{},
			[]string{actCheckIncluded},
			&condition{
				cs: []k.TLSCipher{k.TLSCipher_TLS_CHACHA20_POLY1305_SHA256},
			},
			&action{
				cs: []uint16{tls.TLS_CHACHA20_POLY1305_SHA256},
			},
		),
		gen(
			"TLS_FALLBACK_SCSV",
			[]string{},
			[]string{actCheckIncluded},
			&condition{
				cs: []k.TLSCipher{k.TLSCipher_TLS_FALLBACK_SCSV},
			},
			&action{
				cs: []uint16{tls.TLS_FALLBACK_SCSV},
			},
		),
		gen(
			"TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305",
			[]string{},
			[]string{actCheckIncluded},
			&condition{
				cs: []k.TLSCipher{k.TLSCipher_TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305},
			},
			&action{
				cs: []uint16{tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305},
			},
		),
		gen(
			"TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305",
			[]string{},
			[]string{actCheckIncluded},
			&condition{
				cs: []k.TLSCipher{k.TLSCipher_TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305},
			},
			&action{
				cs: []uint16{tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305},
			},
		),
		gen(
			"input 2 ciphers",
			[]string{},
			[]string{actCheckIncluded},
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
			[]string{},
			[]string{actCheckIncluded},
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

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			cs := network.TLSCiphers(tt.C().cs)
			testutil.Diff(t, tt.A().cs, cs)
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

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	cndNil := tb.Condition("nil", "input nil as curve ids")
	cndInvalid := tb.Condition("invalid", "input invalid curve id")
	cndInputP256 := tb.Condition("P256", "input P256 curve id")
	cndInputP384 := tb.Condition("P384", "input P384 curve id")
	cndInputP512 := tb.Condition("P512", "input P512 curve id")
	cndInputX25519 := tb.Condition("X25519", "input X25519 curve id")
	actCheckP256 := tb.Action("P256", "check that P256 is included")
	actCheckP384 := tb.Action("P384", "check that P384 is included")
	actCheckP512 := tb.Action("P512", "check that P512 is included")
	actCheckX25519 := tb.Action("X25519", "check that X25519 is included")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"nil",
			[]string{cndNil},
			[]string{},
			&condition{
				ids: []k.CurveID{},
			},
			&action{
				ids: []tls.CurveID{},
			},
		),
		gen(
			"invalid",
			[]string{cndInvalid},
			[]string{},
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
			[]string{cndInputP256},
			[]string{actCheckP256},
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
			[]string{cndInputP384},
			[]string{actCheckP384},
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
			[]string{cndInputP512},
			[]string{actCheckP512},
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
			[]string{cndInputX25519},
			[]string{actCheckX25519},
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
			[]string{cndInputP256, cndInputP384, cndInputP512, cndInputX25519},
			[]string{actCheckP256, actCheckP384, actCheckP512, actCheckX25519},
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

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			ids := network.CurveIDs(tt.C().ids)
			testutil.Diff(t, tt.A().ids, ids)
		})
	}
}

func TestCAs(t *testing.T) {
	type condition struct {
		ignore bool
		files  []string
	}

	type action struct {
		pool *x509.CertPool
		err  error
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	cndSystemPool := tb.Condition("system pool", "")
	cndFreshPool := tb.Condition("fresh pool", "")
	cndWithValidCert := tb.Condition("valid cert", "")
	cndWithInvalidCert := tb.Condition("invalid cert", "")
	actCheckNoError := tb.Action("no error", "no error and a valid cert pool returned")
	actCheckError := tb.Action("error", "an error and no cert pool returned")
	table := tb.Build()

	newSystemPool := func(files ...string) *x509.CertPool {
		pool, err := x509.SystemCertPool()
		if err != nil {
			panic(err) // Bad test environment.
		}
		for _, f := range files {
			b, err := os.ReadFile(f)
			if err != nil {
				panic(err) // Bad test condition.
			}
			if !pool.AppendCertsFromPEM(b) {
				panic("failed to read pem for test")
			}
		}
		return pool
	}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"system cert pool",
			[]string{cndSystemPool},
			[]string{actCheckNoError},
			&condition{
				ignore: false,
				files:  nil,
			},
			&action{
				pool: newSystemPool(),
			},
		),
		gen(
			"fresh cert pool",
			[]string{cndFreshPool},
			[]string{actCheckNoError},
			&condition{
				ignore: true,
				files:  nil,
			},
			&action{
				pool: x509.NewCertPool(),
			},
		),
		gen(
			"system cert pool with pem",
			[]string{cndSystemPool, cndWithValidCert},
			[]string{actCheckNoError},
			&condition{
				ignore: false,
				files:  []string{testDir + "ut/core/utilhttp/test.crt"},
			},
			&action{
				pool: newSystemPool(testDir + "ut/core/utilhttp/test.crt"),
			},
		),
		gen(
			"pem not exists",
			[]string{cndFreshPool, cndWithInvalidCert},
			[]string{actCheckError},
			&condition{
				ignore: true,
				files:  []string{testDir + "ut/core/utilhttp/not-exists.pem"},
			},
			&action{
				err: &er.Error{
					Package:     network.ErrPkg,
					Type:        network.ErrTypeTLSCert,
					Description: network.ErrDscTLSCert,
				},
			},
		),
		gen(
			"invalid pem",
			[]string{cndFreshPool, cndWithInvalidCert},
			[]string{actCheckError},
			&condition{
				ignore: true,
				files:  []string{testDir + "ut/core/utilhttp/invalid.pem"},
			},
			&action{
				err: &er.Error{
					Package:     network.ErrPkg,
					Type:        network.ErrTypeTLSCert,
					Description: network.ErrDscTLSCert,
				},
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			pool, err := network.CAs(tt.C().ignore, tt.C().files)
			testutil.Diff(t, tt.A().err, err, cmpopts.EquateErrors())
			if err != nil {
				testutil.Diff(t, (*x509.CertPool)(nil), pool)
				return
			}

			testutil.Diff(t, true, tt.A().pool.Equal(pool))
		})
	}
}
