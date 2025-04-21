// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package authn

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aileron-gateway/aileron-gateway/app"
	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
	utilhttp "github.com/aileron-gateway/aileron-gateway/util/http"
)

type testHandler struct {
	name         string
	result       app.AuthResult
	shouldReturn bool
	err          error
}

func (h *testHandler) ServeAuthn(w http.ResponseWriter, r *http.Request) (*http.Request, app.AuthResult, bool, error) {
	r.Header.Add("Handlers", h.name)
	return r, h.result, h.shouldReturn, h.err
}

func TestMiddleware(t *testing.T) {
	type condition struct {
		handlers []app.AuthenticationHandler
	}

	type action struct {
		authenticated bool
		statusCode    int
		handlers      []string
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"0 handler",
			[]string{},
			[]string{},
			&condition{
				handlers: []app.AuthenticationHandler{},
			},
			&action{
				authenticated: false,
				statusCode:    http.StatusUnauthorized,
			},
		),
		gen(
			"1 handler/unauthorized",
			[]string{},
			[]string{},
			&condition{
				handlers: []app.AuthenticationHandler{
					&testHandler{
						name:   "handler-1",
						result: app.AuthFailed,
					},
				},
			},
			&action{
				authenticated: false,
				statusCode:    http.StatusUnauthorized,
				handlers:      []string{"handler-1"},
			},
		),
		gen(
			"1 handler/authorized",
			[]string{},
			[]string{},
			&condition{
				handlers: []app.AuthenticationHandler{
					&testHandler{
						name:   "handler-1",
						result: app.AuthSucceeded,
					},
				},
			},
			&action{
				authenticated: true,
				statusCode:    http.StatusOK,
				handlers:      []string{"handler-1"},
			},
		),
		gen(
			"2 handler/unauthorized",
			[]string{},
			[]string{},
			&condition{
				handlers: []app.AuthenticationHandler{
					&testHandler{
						name:   "handler-1",
						result: app.AuthFailed,
					},
					&testHandler{
						name:   "handler-2",
						result: app.AuthFailed,
					},
				},
			},
			&action{
				authenticated: false,
				statusCode:    http.StatusUnauthorized,
				handlers:      []string{"handler-1", "handler-2"},
			},
		),
		gen(
			"2 handler/authorized 1st",
			[]string{},
			[]string{},
			&condition{
				handlers: []app.AuthenticationHandler{
					&testHandler{
						name:   "handler-1",
						result: app.AuthSucceeded,
					},
					&testHandler{
						name:   "handler-2",
						result: app.AuthFailed,
					},
				},
			},
			&action{
				authenticated: true,
				statusCode:    http.StatusOK,
				handlers:      []string{"handler-1"},
			},
		),
		gen(
			"2 handler/authorized 2nd",
			[]string{},
			[]string{},
			&condition{
				handlers: []app.AuthenticationHandler{
					&testHandler{
						name:   "handler-1",
						result: app.AuthFailed,
					},
					&testHandler{
						name:   "handler-2",
						result: app.AuthSucceeded,
					},
				},
			},
			&action{
				authenticated: true,
				statusCode:    http.StatusOK,
				handlers:      []string{"handler-1", "handler-2"},
			},
		),
		gen(
			"should return",
			[]string{},
			[]string{},
			&condition{
				handlers: []app.AuthenticationHandler{
					&testHandler{
						name:         "handler-1",
						result:       app.AuthFailed,
						shouldReturn: true,
					},
					&testHandler{
						name:   "handler-2",
						result: app.AuthFailed,
					},
				},
			},
			&action{
				authenticated: false,
				statusCode:    http.StatusOK,
				handlers:      []string{"handler-1"},
			},
		),
		gen(
			"error",
			[]string{},
			[]string{},
			&condition{
				handlers: []app.AuthenticationHandler{
					&testHandler{
						name:   "handler-1",
						result: app.AuthFailed,
						err:    errors.New("test error"),
					},
					&testHandler{
						name:   "handler-2",
						result: app.AuthFailed,
					},
				},
			},
			&action{
				authenticated: false,
				statusCode:    http.StatusInternalServerError,
				handlers:      []string{"handler-1"},
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			resp := httptest.NewRecorder()
			authenticated := false
			h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				authenticated = true
				w.WriteHeader(http.StatusOK)
			})
			m := &authn{
				eh:       utilhttp.GlobalErrorHandler(utilhttp.DefaultErrorHandlerName),
				handlers: tt.C().handlers,
			}
			m.Middleware(h).ServeHTTP(resp, req)

			testutil.Diff(t, tt.A().authenticated, authenticated)
			testutil.Diff(t, tt.A().statusCode, resp.Result().StatusCode)
			testutil.Diff(t, tt.A().handlers, req.Header["Handlers"])
		})
	}
}
