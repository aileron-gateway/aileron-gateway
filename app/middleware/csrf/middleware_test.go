// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package csrf

import (
	"crypto/rand"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aileron-gateway/aileron-gateway/apis/kernel"
	"github.com/aileron-gateway/aileron-gateway/core"
	"github.com/aileron-gateway/aileron-gateway/internal/hash"
	"github.com/aileron-gateway/aileron-gateway/internal/testutil"
	utilhttp "github.com/aileron-gateway/aileron-gateway/util/http"
	"github.com/aileron-projects/go/zcrypto/zsha256"
)

type mockStrategy struct {
	getToken string
	getError core.HTTPError
	setError core.HTTPError
}

func (m *mockStrategy) get(_ *http.Request) (string, core.HTTPError) {
	return m.getToken, m.getError
}

func (m *mockStrategy) set(_ string, _ http.ResponseWriter, _ *http.Request) core.HTTPError {
	return m.setError
}

func TestMiddleware(t *testing.T) {

	type condition struct {
		proxyHeaderName string
		strategy        *mockStrategy
		initialHeader   bool
	}

	type action struct {
		status         int
		expectedHeader string
		expectError    bool
	}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"valid token",
			&condition{
				proxyHeaderName: "X-CSRF-Token",
				strategy: &mockStrategy{
					getToken: "valid-token",
					getError: nil,
				},
				initialHeader: true,
			},
			&action{
				status:         http.StatusOK,
				expectedHeader: "valid-token",
				expectError:    false,
			},
		),
		gen(
			"invalid token",
			&condition{
				proxyHeaderName: "X-CSRF-Token",
				strategy: &mockStrategy{
					getToken: "",
					getError: utilhttp.NewHTTPError(nil, http.StatusForbidden),
				},
				initialHeader: true,
			},
			&action{
				status:      http.StatusForbidden,
				expectError: true,
			},
		),
		gen(
			"no token",
			&condition{
				proxyHeaderName: "X-CSRF-Token",
				strategy: &mockStrategy{
					getToken: "",
					getError: nil,
				},
				initialHeader: true,
			},
			&action{
				status:      http.StatusInternalServerError,
				expectError: true,
			},
		),
		gen(
			"valid token with no initial header",
			&condition{
				proxyHeaderName: "X-CSRF-Token",
				strategy: &mockStrategy{
					getToken: "valid-token",
					getError: nil,
				},
				initialHeader: false,
			},
			&action{
				status:         http.StatusOK,
				expectedHeader: "valid-token",
				expectError:    false,
			},
		),
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {

			csrfMiddleware := &csrf{
				proxyHeaderName: tt.C.proxyHeaderName,
				eh:              utilhttp.GlobalErrorHandler(utilhttp.DefaultErrorHandlerName),
				st:              tt.C.strategy,
			}
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			ctx := req.Context()
			if tt.C.initialHeader {
				header := make(http.Header)
				ctx = utilhttp.ContextWithProxyHeader(ctx, header)
			}
			req = req.WithContext(ctx)
			resp := httptest.NewRecorder()

			nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Check if the header is set when the next handler is called
				if tt.A.expectedHeader != "" {
					header := utilhttp.ProxyHeaderFromContext(r.Context())
					testutil.Diff(t, tt.A.expectedHeader, header.Get(tt.C.proxyHeaderName))
				}
			})

			csrfMiddleware.Middleware(nextHandler).ServeHTTP(resp, req)

			testutil.Diff(t, tt.A.status, resp.Code)

		})
	}
}

func TestServeHTTP(t *testing.T) {

	type condition struct {
		issueNew bool
		strategy *mockStrategy
		token    *csrfToken
		accept   string
	}

	type action struct {
		status         int
		expectNewToken bool
		contentType    string
		expectError    bool
		body           string
	}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"issue new token",
			&condition{
				issueNew: true,
				token: &csrfToken{
					secret:   []byte("some-secret-key"),
					seedSize: 32,
					hashSize: 32,
					hmac:     hash.HMACFromHashAlg(kernel.HashAlg_SHA256),
				},
				strategy: &mockStrategy{},
			},
			&action{
				status:         http.StatusOK,
				expectNewToken: true,
				contentType:    "text/plain; charset=utf-8",
			},
		),
		gen(
			"return existing token",
			&condition{
				issueNew: false,
				token: &csrfToken{
					secret:   []byte("some-secret-key"),
					seedSize: 32,
					hashSize: 32,
					hmac:     hash.HMACFromHashAlg(kernel.HashAlg_SHA256),
				},
				strategy: &mockStrategy{
					getToken: "existing-token",
				},
			},
			&action{
				status:         http.StatusOK,
				expectNewToken: false,
				contentType:    "text/plain; charset=utf-8",
			},
		),
		gen(
			"issue new token with text/plain",
			&condition{
				issueNew: false,
				token: &csrfToken{
					secret:   []byte("some-secret-key"),
					seedSize: 32,
					hashSize: 32,
					hmac:     hash.HMACFromHashAlg(kernel.HashAlg_SHA256),
				},
				accept:   "text/plain",
				strategy: &mockStrategy{},
			},
			&action{
				status:         http.StatusOK,
				contentType:    "text/plain; charset=utf-8",
				expectNewToken: true,
				expectError:    false,
			},
		),
		gen(
			"issue new token with application/json",
			&condition{
				issueNew: false,
				token: &csrfToken{
					secret:   []byte("some-secret-key"),
					seedSize: 32,
					hashSize: 32,
					hmac:     hash.HMACFromHashAlg(kernel.HashAlg_SHA256),
				},
				accept:   "application/json",
				strategy: &mockStrategy{},
			},
			&action{
				status:         http.StatusOK,
				contentType:    "application/json; charset=utf-8",
				expectNewToken: true,
				expectError:    false,
			},
		),
		gen(
			"issue new token with application/xml",
			&condition{
				issueNew: false,
				token: &csrfToken{
					secret:   []byte("some-secret-key"),
					seedSize: 32,
					hashSize: 32,
					hmac:     hash.HMACFromHashAlg(kernel.HashAlg_SHA256),
				},
				accept:   "application/xml",
				strategy: &mockStrategy{},
			},
			&action{
				status:         http.StatusOK,
				contentType:    "application/xml; charset=utf-8",
				expectNewToken: true,
				expectError:    false,
			},
		),
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {

			csrfToken := tt.C.token
			csrfMiddleware := &csrf{
				issueNew: tt.C.issueNew,
				token:    csrfToken,
				st:       tt.C.strategy,
				eh:       utilhttp.GlobalErrorHandler(utilhttp.DefaultErrorHandlerName),
			}

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.Header.Set("Accept", tt.C.accept)
			resp := httptest.NewRecorder()

			csrfMiddleware.ServeHTTP(resp, req)

			testutil.Diff(t, tt.A.status, resp.Code)
			testutil.Diff(t, tt.A.contentType, resp.Header().Get("Content-Type"))

			if tt.A.expectNewToken {
				// Check if a new token has been issued
				token := resp.Body.String()
				testutil.Diff(t, true, len(token) > 0)
			}

			if tt.A.expectError {
				testutil.Diff(t, tt.A.status, http.StatusInternalServerError)
				testutil.Diff(t, tt.A.body, resp.Body.String())
			}
		})
	}
}

// Variable to save the original rand.Reader
var originalRandReader = rand.Reader

type mockRandReader struct{}

func (m *mockRandReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("mock random read error")
}

func TestServeHTTPWithRandError(t *testing.T) {
	rand.Reader = &mockRandReader{}
	defer func() {
		rand.Reader = originalRandReader
	}()

	csrfToken := &csrfToken{
		secret:   []byte("some-secret"),
		seedSize: 32,
		hashSize: 32,
		hmac:     zsha256.HMACSum256,
	}

	strategy := &mockStrategy{
		getToken: "",
		getError: nil,
		setError: nil,
	}

	csrfMiddleware := &csrf{
		issueNew: true,
		token:    csrfToken,
		st:       strategy,
		eh:       utilhttp.GlobalErrorHandler(utilhttp.DefaultErrorHandlerName),
	}

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Accept", "text/plain")
	resp := httptest.NewRecorder()

	csrfMiddleware.ServeHTTP(resp, req)

	testutil.Diff(t, http.StatusInternalServerError, resp.Code)

	body := resp.Body.String()
	testutil.Diff(t, true, len(body) > 0)
}

func TestServeHTTPWithSetError(t *testing.T) {
	csrfToken := &csrfToken{
		secret:   []byte("some-secret"),
		seedSize: 32,
		hashSize: 32,
		hmac:     zsha256.HMACSum256,
	}

	strategy := &mockStrategy{
		getToken: "valid-token",
		getError: nil,
		setError: utilhttp.NewHTTPError(errors.New("set token error"), http.StatusInternalServerError),
	}

	csrfMiddleware := &csrf{
		issueNew: true,
		token:    csrfToken,
		st:       strategy,
		eh:       utilhttp.GlobalErrorHandler(utilhttp.DefaultErrorHandlerName),
	}

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Accept", "text/plain")
	resp := httptest.NewRecorder()

	csrfMiddleware.ServeHTTP(resp, req)

	testutil.Diff(t, http.StatusInternalServerError, resp.Code)

	body := resp.Body.String()
	testutil.Diff(t, true, len(body) > 0)
}
