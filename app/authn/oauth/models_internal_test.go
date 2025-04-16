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

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	cndValidCnf := tb.Condition("valid cnf", "cnf has a value")
	cndTLSState := tb.Condition("has tls state", "client tls state has non-nil value")
	cndThumbDiff := tb.Condition("thumb diff", "cnf and client cert have different thumbprint")
	actCheckError := tb.Action("error", "check that an error was returned")
	actCheckNoError := tb.Action("no error", "check that there is no error")

	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"success",
			[]string{cndValidCnf, cndTLSState},
			[]string{actCheckNoError},
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
			[]string{cndTLSState},
			[]string{actCheckError},
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
			[]string{cndValidCnf},
			[]string{actCheckError},
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
			[]string{cndTLSState},
			[]string{actCheckError},
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
			[]string{cndValidCnf},
			[]string{actCheckError},
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
			[]string{cndThumbDiff},
			[]string{actCheckError},
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

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			err := validateCert(tt.C().cnf, tt.C().state)
			testutil.DiffError(t, tt.A().err, tt.A().errPattern, err)
		})
	}
}
