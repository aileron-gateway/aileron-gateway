// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package csrf

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/aileron-gateway/aileron-gateway/apis/kernel"
	"github.com/aileron-gateway/aileron-gateway/kernel/hash"
	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
	"github.com/aileron-gateway/aileron-gateway/util/session"
	"github.com/aileron-projects/go/zcrypto/zsha256"
)

type mockExtractor struct {
	token string
	err   error
}

func (m *mockExtractor) extract(r *http.Request) (string, error) {
	return m.token, m.err
}

func TestCustomRequestHeaders(t *testing.T) {
	type condition struct {
		headerValue string
		pattern     string
		validToken  bool
	}

	type action struct {
		expectError bool
		status      int
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"empty header",
			[]string{},
			[]string{},
			&condition{
				headerValue: "",
				pattern:     ".*",
				validToken:  false,
			},
			&action{
				expectError: true,
				status:      http.StatusForbidden,
			},
		),
		gen(
			"invalid pattern",
			[]string{},
			[]string{},
			&condition{
				headerValue: "invalid-token",
				pattern:     "^valid-.*",
				validToken:  false,
			},
			&action{
				expectError: true,
				status:      http.StatusForbidden,
			},
		),
		gen(
			"valid token",
			[]string{},
			[]string{},
			&condition{
				headerValue: "valid-token",
				pattern:     ".*",
				validToken:  true,
			},
			&action{
				expectError: false,
				status:      http.StatusOK,
			},
		),
		gen(
			"valid token and pattern",
			[]string{},
			[]string{},
			&condition{
				headerValue: "valid-token",
				pattern:     "^valid-.*",
				validToken:  true,
			},
			&action{
				expectError: false,
				status:      http.StatusOK,
			},
		),
		gen(
			"set token",
			[]string{},
			[]string{},
			&condition{
				headerValue: "",
				pattern:     "",
				validToken:  true,
			},
			&action{
				expectError: false,
				status:      http.StatusOK,
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {

			csrfToken := &csrfToken{
				secret:   []byte("some-secret"),
				seedSize: 32,
				hashSize: 32,
				hmac:     hash.HMACFromHashAlg(kernel.HashAlg_SHA256),
			}

			// Generate if a valid token is required.
			validToken, _ := csrfToken.new()
			headerValue := tt.C().headerValue
			if tt.C().validToken {
				headerValue = validToken
			}

			var pattern *regexp.Regexp
			if tt.C().pattern != "" {
				pattern = regexp.MustCompile(tt.C().pattern)
			}

			strategy := &customRequestHeaders{
				headerName: "X-CSRF-Token",
				pattern:    pattern,
				token:      csrfToken,
			}

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			w := httptest.NewRecorder()

			req.Header.Set("X-CSRF-Token", headerValue)

			token, err := strategy.get(req)
			if tt.A().expectError {
				testutil.Diff(t, err != nil, true)
				testutil.Diff(t, token, "")
			} else {
				testutil.Diff(t, true, err == nil)
				testutil.Diff(t, token, req.Header.Get("X-CSRF-Token"))
			}

			err = strategy.set(validToken, w, req)
			testutil.Diff(t, true, err == nil)
			testutil.Diff(t, "", w.Header().Get("X-CSRF-Token"))
		})
	}
}

func TestDoubleSubmitCookies(t *testing.T) {
	type condition struct {
		cookieToken  string
		requestToken string
		validToken   bool
		noCookie     bool
		extractError bool
	}

	type action struct {
		expectError bool
		status      int
		expectSet   bool
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"no cookie",
			[]string{},
			[]string{},
			&condition{
				cookieToken:  "",
				requestToken: "",
				validToken:   false,
				noCookie:     true,
			},
			&action{
				expectError: true,
				status:      http.StatusForbidden,
			},
		),
		gen(
			"token mismatch",
			[]string{},
			[]string{},
			&condition{
				cookieToken:  "cookie-token",
				requestToken: "request-token",
				validToken:   false,
			},
			&action{
				expectError: true,
				status:      http.StatusForbidden,
			},
		),
		gen(
			"valid token",
			[]string{},
			[]string{},
			&condition{
				validToken: true,
			},
			&action{
				expectError: false,
				status:      http.StatusOK,
			},
		),
		gen(
			"extract token error",
			[]string{},
			[]string{},
			&condition{
				validToken:   true,
				extractError: true,
			},
			&action{
				expectError: true,
				status:      http.StatusForbidden,
			},
		),
		gen(
			"set token successfully",
			[]string{},
			[]string{},
			&condition{
				cookieToken: "new-valid-token",
				validToken:  true,
			},
			&action{
				expectError: false,
				status:      http.StatusOK,
				expectSet:   true,
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {

			csrfToken := &csrfToken{
				secret:   []byte("some-secret"),
				seedSize: 32,
				hashSize: 32,
				hmac:     zsha256.HMACSum256,
			}

			cookieToken, _ := csrfToken.new()
			requestToken := cookieToken
			if !tt.C().validToken {
				requestToken = "invalid-token"
			}

			var mockError error
			if tt.C().extractError {
				mockError = errors.New("extract tokken failed.")
			}
			strategy := &doubleSubmitCookies{
				token:      csrfToken,
				ext:        &mockExtractor{token: requestToken, err: mockError},
				cookieName: "csrf-token",
				cookie:     &mockCookieCreateor{},
			}

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			resp := httptest.NewRecorder()

			if !tt.C().noCookie {
				ck := &http.Cookie{
					Name:  "csrf-token",
					Value: cookieToken,
				}
				req.AddCookie(ck)
			}
			token, err := strategy.get(req)

			if tt.A().expectError {
				testutil.Diff(t, true, err != nil)
				testutil.Diff(t, "", token)
			} else {
				testutil.Diff(t, true, err == nil)
				testutil.Diff(t, token, requestToken)
			}

			if tt.A().expectSet {
				err := strategy.set(tt.C().cookieToken, resp, req)
				testutil.Diff(t, err == nil, true)

				cookieHeader := resp.Header().Get("Set-Cookie")
				testutil.Diff(t, "csrf-token=new-valid-token; Path=/; Max-Age=100", cookieHeader)
			}
		})
	}
}

type mockCookieCreateor struct{}

func (m *mockCookieCreateor) NewCookie() *http.Cookie {
	return &http.Cookie{
		Path:   "/",
		MaxAge: 100,
	}
}

// mockSession is a mock that implements the session.Session interface.
type mockSession struct {
	token      string
	err        error
	persistErr error
}

func (m *mockSession) Extract(key string, value any) error {
	if key == "__csrf_token__" {
		if m.err != nil {
			return m.err
		}
		*(value.(*[]byte)) = []byte(m.token)
		return nil
	}
	return errors.New("invalid key")
}

func (m *mockSession) Persist(key string, value any) error {
	if m.persistErr != nil {
		return m.persistErr
	}

	if key == "__csrf_token__" {
		if token, ok := value.([]byte); ok {
			m.token = string(token)
			return nil
		}
		return errors.New("invalid value type.")
	}
	return errors.New("invalid key.")
}

func (m *mockSession) Delete(key string) {
}

func (m *mockSession) SetFlag(flag uint) uint {
	return flag
}

func (m *mockSession) MarshalBinary() ([]byte, error) {
	return []byte(m.token), nil
}

func (m *mockSession) UnmarshalBinary(data []byte) error {
	m.token = string(data)
	return nil
}

func (m *mockSession) Attributes() map[string]any {
	return nil
}

func TestSynchronizerToken(t *testing.T) {
	type condition struct {
		sessionToken  string
		requestToken  string
		validToken    bool
		sessionExists bool
		extractError  bool
		tokenError    bool
		persistError  bool
	}

	type action struct {
		expectError  bool
		persistError bool
		status       int
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"no session",
			[]string{},
			[]string{},
			&condition{
				sessionToken:  "",
				requestToken:  "",
				validToken:    false,
				sessionExists: false,
			},
			&action{
				expectError: true,
				status:      http.StatusInternalServerError,
			},
		),
		gen(
			"session token extraction failure",
			[]string{},
			[]string{},
			&condition{
				sessionToken:  "",
				requestToken:  "",
				validToken:    false,
				sessionExists: true,
				extractError:  true,
			},
			&action{
				expectError: true,
				status:      http.StatusInternalServerError,
			},
		),
		gen(
			"request token extraction failure",
			[]string{},
			[]string{},
			&condition{
				sessionToken:  "valid-token",
				requestToken:  "",
				validToken:    false,
				sessionExists: true,
				tokenError:    true,
			},
			&action{
				expectError: true,
				status:      http.StatusForbidden,
			},
		),
		gen(
			"token mismatch",
			[]string{},
			[]string{},
			&condition{
				sessionToken:  "session-token",
				requestToken:  "request-token",
				validToken:    false,
				sessionExists: true,
			},
			&action{
				expectError: true,
				status:      http.StatusForbidden,
			},
		),
		gen(
			"valid token",
			[]string{},
			[]string{},
			&condition{
				validToken:    true,
				sessionExists: true,
			},
			&action{
				expectError: false,
				status:      http.StatusOK,
			},
		),
		gen(
			"persist token failure",
			[]string{},
			[]string{},
			&condition{
				sessionToken:  "valid-token",
				requestToken:  "valid-token",
				validToken:    true,
				sessionExists: true,
				persistError:  true,
			},
			&action{
				persistError: true,
				status:       http.StatusInternalServerError,
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {

			csrfToken := &csrfToken{
				secret:   []byte("some-secret"),
				seedSize: 32,
				hashSize: 32,
				hmac:     zsha256.HMACSum256,
			}

			sessionToken, _ := csrfToken.new()
			requestToken := sessionToken
			if !tt.C().validToken {
				requestToken = "invalid-token"
			}

			var mockExtErr error
			if tt.C().tokenError {
				mockExtErr = errors.New("request token extraction failed")
			}

			strategy := &synchronizerToken{
				ext:   &mockExtractor{token: requestToken, err: mockExtErr},
				token: csrfToken,
			}

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			ctx := req.Context()

			if tt.C().sessionExists {
				var mockSessErr error
				if tt.C().extractError {
					mockSessErr = errors.New("session token extraction failed")
				}
				ctx = session.ContextWithSession(
					req.Context(),
					&mockSession{
						token: sessionToken,
						err:   mockSessErr,
						persistErr: func() error {
							if tt.C().persistError {
								return errors.New("persist failed")
							}
							return nil
						}(),
					},
				)
			}

			req = req.WithContext(ctx)
			token, err := strategy.get(req)
			if tt.A().expectError {
				testutil.Diff(t, true, err != nil)
				testutil.Diff(t, "", token)
			} else {
				testutil.Diff(t, true, err == nil)
				testutil.Diff(t, requestToken, token)
			}

			err = strategy.set(requestToken, nil, req)
			if tt.C().persistError {
				testutil.Diff(t, true, err != nil)
			}
		})
	}
}
