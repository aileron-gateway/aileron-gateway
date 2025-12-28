// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package oauth

import (
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"regexp"
	"testing"

	"github.com/aileron-gateway/aileron-gateway/app"
	"github.com/aileron-gateway/aileron-gateway/core"
	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
)

func x5tS256(in []byte) string {
	b := sha256.Sum256(in)
	return base64.RawURLEncoding.EncodeToString(b[:])
}

func TestValidateCert(t *testing.T) {
	type condition struct {
		cnf   map[string]any
		state *tls.ConnectionState
	}

	type action struct {
		err        any // error or errorutil.Kind
		errPattern *regexp.Regexp
	}
	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"success",
			&condition{
				cnf: map[string]any{
					"x5t#S256": x5tS256([]byte("test")),
				},
				state: &tls.ConnectionState{
					PeerCertificates: []*x509.Certificate{
						{
							Raw: []byte("test"),
						},
					},
				},
			},
			&action{
				err: nil,
			},
		),
		gen(
			"cnf is nil",
			&condition{
				cnf: nil,
				state: &tls.ConnectionState{
					PeerCertificates: []*x509.Certificate{
						{
							Raw: []byte("test"),
						},
					},
				},
			},
			&action{
				err:        app.ErrAppAuthnInvalidCredential,
				errPattern: regexp.MustCompile(core.ErrPrefix + `invalid credential for mTLS. insufficient claims`),
			},
		),
		gen(
			"tls state is nil",
			&condition{
				cnf: map[string]any{
					"x5t#S256": x5tS256([]byte("test")),
				},
				state: nil,
			},
			&action{
				err:        app.ErrAppAuthnInvalidCredential,
				errPattern: regexp.MustCompile(core.ErrPrefix + `invalid credential for mTLS. no tls state`),
			},
		),
		gen(
			"x5t#S256 is not string",
			&condition{
				cnf: map[string]any{
					"x5t#S256": 12345,
				},
				state: &tls.ConnectionState{
					PeerCertificates: []*x509.Certificate{
						{
							Raw: []byte("test"),
						},
					},
				},
			},
			&action{
				err:        app.ErrAppAuthnInvalidCredential,
				errPattern: regexp.MustCompile(core.ErrPrefix + `invalid credential for mTLS. insufficient cnf claims`),
			},
		),
		gen(
			"no peer certs",
			&condition{
				cnf: map[string]any{
					"x5t#S256": x5tS256([]byte("test")),
				},
				state: &tls.ConnectionState{
					PeerCertificates: []*x509.Certificate{},
				},
			},
			&action{
				err:        app.ErrAppAuthnInvalidCredential,
				errPattern: regexp.MustCompile(core.ErrPrefix + `invalid credential for mTLS. peer certificates not found`),
			},
		),
		gen(
			"thumbprint mismatched",
			&condition{
				cnf: map[string]any{
					"x5t#S256": x5tS256([]byte("test")),
				},
				state: &tls.ConnectionState{
					PeerCertificates: []*x509.Certificate{
						{
							Raw: []byte("test-test"),
						},
					},
				},
			},
			&action{
				err:        app.ErrAppAuthnInvalidCredential,
				errPattern: regexp.MustCompile(core.ErrPrefix + `invalid credential for mTLS. thumbprint mismatch`),
			},
		),
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			err := validateCert(tt.C.cnf, tt.C.state)
			testutil.DiffError(t, tt.A.err, tt.A.errPattern, err)
		})
	}
}
