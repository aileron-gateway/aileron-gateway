// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package httpserver

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aileron-gateway/aileron-gateway/kernel/log"
	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
	"github.com/aileron-gateway/aileron-gateway/kernel/uid"
	httputil "github.com/aileron-gateway/aileron-gateway/util/http"
)

func TestMiddleware(t *testing.T) {
	type condition struct {
		panics    any
		requestID string
	}

	type action struct {
		statusCode int
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())

	CndSetNoPanic := tb.Condition("set no panic", "set no panic")
	CndSetPanic := tb.Condition("set panic", "set panic")
	CndSetErrAbortHandler := tb.Condition("set ErrAbortHandler err", "set ErrAbortHandler err")
	CndSetOtherErrHandler := tb.Condition("set other error handler", "set other error handler")
	CndSetReqID := tb.Condition("set request id", "set request id")
	ActCheckExpected := tb.Action("expected value returned", "check that an expected value returned")

	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"no panic",
			[]string{CndSetNoPanic},
			[]string{ActCheckExpected},
			&condition{},
			&action{
				statusCode: http.StatusOK,
			},
		),
		gen(
			"panic occur",
			[]string{CndSetPanic},
			[]string{ActCheckExpected},
			&condition{
				panics: "panic occurred",
			},
			&action{
				statusCode: http.StatusInternalServerError,
			},
		),
		gen(
			"ErrAbortHandler pattern",
			[]string{CndSetErrAbortHandler},
			[]string{ActCheckExpected},
			&condition{
				panics: http.ErrAbortHandler,
			},
			&action{
				statusCode: http.StatusOK,
			},
		),
		gen(
			"other error handler pattern",
			[]string{CndSetOtherErrHandler},
			[]string{ActCheckExpected},
			&condition{
				panics: http.ErrContentLength,
			},
			&action{
				statusCode: http.StatusInternalServerError,
			},
		),
		gen(
			"set request id",
			[]string{CndSetReqID},
			[]string{ActCheckExpected},
			&condition{
				panics:    "panic occurred",
				requestID: "test-request-id",
			},
			&action{
				statusCode: http.StatusInternalServerError,
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			rc := recoverer{
				lg: log.GlobalLogger(log.DefaultLoggerName),
				eh: httputil.GlobalErrorHandler(httputil.DefaultErrorHandlerName),
			}

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			ctx := uid.ContextWithID(req.Context(), tt.C().requestID)
			req = req.WithContext(ctx)

			resp := httptest.NewRecorder()

			h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if tt.C().panics != nil {
					panic(tt.C().panics)
				}
			})
			rc.Middleware(h).ServeHTTP(resp, req)

			testutil.Diff(t, tt.A().statusCode, resp.Code)
		})
	}
}
