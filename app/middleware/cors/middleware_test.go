// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package cors

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
	utilhttp "github.com/aileron-gateway/aileron-gateway/util/http"
)

func TestCORS_Middleware(t *testing.T) {
	type condition struct {
		method  string
		headers map[string]string
		policy  *corsPolicy
	}

	type action struct {
		status int
		header http.Header
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"Preflight/method allowed",
			[]string{},
			[]string{},
			&condition{
				method: http.MethodOptions,
				headers: map[string]string{
					"Access-Control-Request-Method": http.MethodGet,
					"Origin":                        "http://test.com",
				},
				policy: &corsPolicy{
					allowedOrigins:       []string{"http://test.com"},
					allowedMethods:       []string{http.MethodGet},
					joinedAllowedMethods: http.MethodGet,
				},
			},
			&action{
				status: http.StatusOK,
				header: http.Header{
					"Access-Control-Allow-Origin":  {"http://test.com"},
					"Access-Control-Allow-Methods": {"GET"},
					"Vary":                         {"Origin", "Access-Control-Request-Method", "Access-Control-Request-Headers"},
				},
			},
		),
		gen(
			"Preflight/method disallowed",
			[]string{},
			[]string{},
			&condition{
				method: http.MethodOptions,
				headers: map[string]string{
					"Access-Control-Request-Method": http.MethodPost,
					"Origin":                        "http://test.com",
				},
				policy: &corsPolicy{
					allowedOrigins:       []string{"http://test.com"},
					allowedMethods:       []string{http.MethodGet},
					joinedAllowedMethods: http.MethodGet,
				},
			},
			&action{
				status: http.StatusOK,
				header: http.Header{
					"Access-Control-Allow-Origin":  {"http://test.com"},
					"Access-Control-Allow-Methods": {"GET"},
					"Vary":                         {"Origin", "Access-Control-Request-Method", "Access-Control-Request-Headers"},
				},
			},
		),
		gen(
			"Preflight/origin disallowed",
			[]string{},
			[]string{},
			&condition{
				method: http.MethodOptions,
				headers: map[string]string{
					"Access-Control-Request-Method": http.MethodGet,
					"Origin":                        "http://disallowed.com",
				},
				policy: &corsPolicy{
					allowedOrigins:       []string{"http://test.com"},
					allowedMethods:       []string{http.MethodGet},
					joinedAllowedMethods: http.MethodGet,
				},
			},
			&action{
				status: http.StatusForbidden,
				header: http.Header{
					"Access-Control-Allow-Methods": {"GET"},
					"Vary":                         {"Origin", "Access-Control-Request-Method", "Access-Control-Request-Headers"},
				},
			},
		),
		gen(
			"Preflight/wildcard origin",
			[]string{},
			[]string{},
			&condition{
				method: http.MethodOptions,
				headers: map[string]string{
					"Access-Control-Request-Method": http.MethodGet,
					"Origin":                        "http://test.com",
				},
				policy: &corsPolicy{
					allowedOrigins:       []string{"*"},
					allowedMethods:       []string{http.MethodGet},
					joinedAllowedMethods: http.MethodGet,
				},
			},
			&action{
				status: http.StatusOK,
				header: http.Header{
					"Access-Control-Allow-Origin":  {"*"},
					"Access-Control-Allow-Methods": {"GET"},
					"Vary":                         {"Origin", "Access-Control-Request-Method", "Access-Control-Request-Headers"},
				},
			},
		),
		gen(
			"Preflight/wildcard origin with disabled wildcard origin",
			[]string{},
			[]string{},
			&condition{
				method: http.MethodOptions,
				headers: map[string]string{
					"Access-Control-Request-Method": http.MethodGet,
					"Origin":                        "http://test.com",
				},
				policy: &corsPolicy{
					allowedOrigins:        []string{"*"},
					disableWildCardOrigin: true,
					allowedMethods:        []string{http.MethodGet},
					joinedAllowedMethods:  http.MethodGet,
				},
			},
			&action{
				status: http.StatusOK,
				header: http.Header{
					"Access-Control-Allow-Origin":  {"http://test.com"},
					"Access-Control-Allow-Methods": {"GET"},
					"Vary":                         {"Origin", "Access-Control-Request-Method", "Access-Control-Request-Headers"},
				},
			},
		),
		gen(
			"Preflight/with Access-Control-Request-Headers",
			[]string{},
			[]string{},
			&condition{
				method: http.MethodOptions,
				headers: map[string]string{
					"Access-Control-Request-Method":  http.MethodGet,
					"Access-Control-Request-Headers": "X-Test-Header",
					"Origin":                         "http://test.com",
				},
				policy: &corsPolicy{
					allowedOrigins:       []string{"http://test.com"},
					allowedMethods:       []string{http.MethodGet},
					joinedAllowedMethods: http.MethodGet,
					joinedAllowedHeaders: "X-Custom-Header",
				},
			},
			&action{
				status: http.StatusOK,
				header: http.Header{
					"Access-Control-Allow-Origin":  {"http://test.com"},
					"Access-Control-Allow-Methods": {"GET"},
					"Access-Control-Allow-Headers": {"X-Custom-Header"},
					"Vary":                         {"Origin", "Access-Control-Request-Method", "Access-Control-Request-Headers"},
				},
			},
		),
		gen(
			"Preflight/with exposed headers",
			[]string{},
			[]string{},
			&condition{
				method: http.MethodOptions,
				headers: map[string]string{
					"Access-Control-Request-Method": http.MethodGet,
					"Origin":                        "http://test.com",
				},
				policy: &corsPolicy{
					allowedOrigins:       []string{"http://test.com"},
					allowedMethods:       []string{http.MethodGet},
					joinedAllowedMethods: http.MethodGet,
					joinedExposedHeaders: "X-Expose-Header",
				},
			},
			&action{
				status: http.StatusOK,
				header: http.Header{
					"Access-Control-Allow-Origin":   {"http://test.com"},
					"Access-Control-Allow-Methods":  {"GET"},
					"Access-Control-Expose-Headers": {"X-Expose-Header"},
					"Vary":                          {"Origin", "Access-Control-Request-Method", "Access-Control-Request-Headers"},
				},
			},
		),
		gen(
			"Preflight/with allowed credentials",
			[]string{},
			[]string{},
			&condition{
				method: http.MethodOptions,
				headers: map[string]string{
					"Access-Control-Request-Method": http.MethodGet,
					"Origin":                        "http://test.com",
				},
				policy: &corsPolicy{
					allowedOrigins:       []string{"http://test.com"},
					allowedMethods:       []string{http.MethodGet},
					allowCredentials:     true,
					joinedAllowedMethods: http.MethodGet,
				},
			},
			&action{
				status: http.StatusOK,
				header: http.Header{
					"Access-Control-Allow-Origin":      {"http://test.com"},
					"Access-Control-Allow-Methods":     {"GET"},
					"Access-Control-Allow-Credentials": {"true"},
					"Vary":                             {"Origin", "Access-Control-Request-Method", "Access-Control-Request-Headers"},
				},
			},
		),
		gen(
			"Preflight/with private network allowed",
			[]string{},
			[]string{},
			&condition{
				method: http.MethodOptions,
				headers: map[string]string{
					"Access-Control-Request-Method":          http.MethodGet,
					"Access-Control-Request-Private-Network": "true",
					"Origin":                                 "http://test.com",
				},
				policy: &corsPolicy{
					allowedOrigins:       []string{"http://test.com"},
					allowedMethods:       []string{http.MethodGet},
					joinedAllowedMethods: http.MethodGet,
					allowPrivateNetwork:  true,
				},
			},
			&action{
				status: http.StatusOK,
				header: http.Header{
					"Access-Control-Allow-Origin":          {"http://test.com"},
					"Access-Control-Allow-Methods":         {"GET"},
					"Access-Control-Allow-Private-Network": {"true"},
					"Vary": {"Origin", "Access-Control-Request-Method",
						"Access-Control-Request-Headers", "Access-Control-Request-Private-Network"},
				},
			},
		),
		gen(
			"Preflight/with max age",
			[]string{},
			[]string{},
			&condition{
				method: http.MethodOptions,
				headers: map[string]string{
					"Access-Control-Request-Method": http.MethodGet,
					"Origin":                        "http://test.com",
				},
				policy: &corsPolicy{
					allowedOrigins:       []string{"http://test.com"},
					allowedMethods:       []string{http.MethodGet},
					maxAge:               "100",
					joinedAllowedMethods: http.MethodGet,
				},
			},
			&action{
				status: http.StatusOK,
				header: http.Header{
					"Access-Control-Allow-Origin":  {"http://test.com"},
					"Access-Control-Allow-Methods": {"GET"},
					"Access-Control-Max-Age":       {"100"},
					"Vary":                         {"Origin", "Access-Control-Request-Method", "Access-Control-Request-Headers"},
				},
			},
		),
		gen(
			"Actual request with allowed headers",
			[]string{},
			[]string{},
			&condition{
				method: http.MethodGet,
				headers: map[string]string{
					"Origin": "http://test.com",
				},
				policy: &corsPolicy{
					allowedOrigins:       []string{"http://test.com"},
					allowedMethods:       []string{http.MethodGet},
					joinedAllowedMethods: http.MethodGet,
					joinedAllowedHeaders: "X-Custom-Header",
				},
			},
			&action{
				status: http.StatusOK,
				header: http.Header{
					"Access-Control-Allow-Origin":  {"http://test.com"},
					"Access-Control-Allow-Methods": {"GET"},
					"Access-Control-Allow-Headers": {"X-Custom-Header"},
					"Vary":                         {"Origin"},
				},
			},
		),
		gen(
			"Actual request with exposed headers",
			[]string{},
			[]string{},
			&condition{
				method: http.MethodGet,
				headers: map[string]string{
					"Origin": "http://test.com",
				},
				policy: &corsPolicy{
					allowedOrigins:       []string{"http://test.com"},
					allowedMethods:       []string{http.MethodGet},
					joinedAllowedMethods: http.MethodGet,
					joinedExposedHeaders: "X-Expose-Header",
				},
			},
			&action{
				status: http.StatusOK,
				header: http.Header{
					"Access-Control-Allow-Origin":   {"http://test.com"},
					"Access-Control-Allow-Methods":  {"GET"},
					"Access-Control-Expose-Headers": {"X-Expose-Header"},
					"Vary":                          {"Origin"},
				},
			},
		),
		gen(
			"Actual request with allowed method",
			[]string{},
			[]string{},
			&condition{
				method: http.MethodGet,
				headers: map[string]string{
					"Origin": "http://test.com",
				},
				policy: &corsPolicy{
					allowedOrigins:       []string{"http://test.com"},
					allowedMethods:       []string{http.MethodGet},
					joinedAllowedMethods: http.MethodGet,
				},
			},
			&action{
				status: http.StatusOK,
				header: http.Header{
					"Access-Control-Allow-Origin":  {"http://test.com"},
					"Access-Control-Allow-Methods": {"GET"},
					"Vary":                         {"Origin"},
				},
			},
		),
		gen(
			"Actual request with credentials allowed",
			[]string{},
			[]string{},
			&condition{
				method: http.MethodGet,
				headers: map[string]string{
					"Origin": "http://test.com",
				},
				policy: &corsPolicy{
					allowedOrigins:       []string{"http://test.com"},
					allowedMethods:       []string{http.MethodGet},
					joinedAllowedMethods: http.MethodGet,
					allowCredentials:     true,
				},
			},
			&action{
				status: http.StatusOK,
				header: http.Header{
					"Access-Control-Allow-Origin":      {"http://test.com"},
					"Access-Control-Allow-Methods":     {"GET"},
					"Access-Control-Allow-Credentials": {"true"},
					"Vary":                             {"Origin"},
				},
			},
		),
		gen(
			"Actual request with disallowed origin",
			[]string{},
			[]string{},
			&condition{
				method: http.MethodGet,
				headers: map[string]string{
					"Origin": "http://disallowed.com",
				},
				policy: &corsPolicy{
					allowedOrigins:       []string{"http://test.com"},
					allowedMethods:       []string{http.MethodGet},
					joinedAllowedMethods: http.MethodGet,
				},
			},
			&action{
				status: http.StatusForbidden,
				header: http.Header{
					"Access-Control-Allow-Origin":  nil,
					"Access-Control-Allow-Methods": {"GET"},
				},
			},
		),
		gen(
			"Actual request with embedder policy",
			[]string{},
			[]string{},
			&condition{
				method: http.MethodGet,
				headers: map[string]string{
					"Origin": "http://test.com",
				},
				policy: &corsPolicy{
					allowedOrigins:       []string{"http://test.com"},
					allowedMethods:       []string{http.MethodGet},
					joinedAllowedMethods: http.MethodGet,
					embedderPolicy:       "require-corp",
				},
			},
			&action{
				status: http.StatusOK,
			},
		),
		gen(
			"Actual request with opener policy",
			[]string{},
			[]string{},
			&condition{
				method: http.MethodGet,
				headers: map[string]string{
					"Origin": "http://test.com",
				},
				policy: &corsPolicy{
					allowedOrigins:       []string{"http://test.com"},
					allowedMethods:       []string{http.MethodGet},
					joinedAllowedMethods: http.MethodGet,
					openerPolicy:         "same-origin",
				},
			},
			&action{
				status: http.StatusOK,
			},
		),
		gen(
			"Actual request with opener policy",
			[]string{},
			[]string{},
			&condition{
				method: http.MethodGet,
				headers: map[string]string{
					"Origin": "http://test.com",
				},
				policy: &corsPolicy{
					allowedOrigins:       []string{"http://test.com"},
					allowedMethods:       []string{http.MethodGet},
					joinedAllowedMethods: http.MethodGet,
					resourcePolicy:       "same-origin",
				},
			},
			&action{
				status: http.StatusOK,
			},
		),
		gen(
			"Actual request with empty origin",
			[]string{},
			[]string{},
			&condition{
				method:  http.MethodGet,
				headers: map[string]string{},
				policy: &corsPolicy{
					allowedOrigins:       []string{"http://test.com"},
					allowedMethods:       []string{http.MethodPost},
					joinedAllowedMethods: http.MethodGet,
				},
			},
			&action{
				status: http.StatusForbidden,
			},
		),
		gen(
			"Forbidden actual request",
			[]string{},
			[]string{},
			&condition{
				method: http.MethodGet,
				headers: map[string]string{
					"Origin": "http://test.com",
				},
				policy: &corsPolicy{
					allowedOrigins:       []string{"http://test.com"},
					allowedMethods:       []string{http.MethodPost},
					joinedAllowedMethods: http.MethodGet,
				},
			},
			&action{
				status: http.StatusForbidden,
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {

			// Prepare the CORS middleware
			corsMiddleware := &cors{
				eh:     utilhttp.GlobalErrorHandler(utilhttp.DefaultErrorHandlerName),
				policy: tt.C().policy,
			}

			// Create a test request
			req := httptest.NewRequest(tt.C().method, "http://test.com", nil)
			for k, v := range tt.C().headers {
				req.Header.Set(k, v)
			}

			// Create a test response recorder
			resp := httptest.NewRecorder()
			// Call the middleware
			corsMiddleware.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})).ServeHTTP(resp, req)

			// Verify the status code
			testutil.Diff(t, tt.A().status, resp.Code)
			t.Log(resp.Header())
			for k, v := range tt.A().header {
				testutil.Diff(t, v, resp.Header()[k])
			}

		})
	}
}
