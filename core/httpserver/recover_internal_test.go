// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package httpserver

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aileron-gateway/aileron-gateway/kernel/log"
	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
	httputil "github.com/aileron-gateway/aileron-gateway/util/http"
	"github.com/aileron-projects/go/zx/zuid"
)

func TestMiddleware(t *testing.T) {
	type condition struct {
		panics    any
		requestID string
	}

	type action struct {
		statusCode int
	}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"no panic",
			&condition{},
			&action{
				statusCode: http.StatusOK,
			},
		),
		gen(
			"panic occur",
			&condition{
				panics: "panic occurred",
			},
			&action{
				statusCode: http.StatusInternalServerError,
			},
		),
		gen(
			"ErrAbortHandler pattern",
			&condition{
				panics: http.ErrAbortHandler,
			},
			&action{
				statusCode: http.StatusOK,
			},
		),
		gen(
			"other error handler pattern",
			&condition{
				panics: http.ErrContentLength,
			},
			&action{
				statusCode: http.StatusInternalServerError,
			},
		),
		gen(
			"set request id",
			&condition{
				panics:    "panic occurred",
				requestID: "test-request-id",
			},
			&action{
				statusCode: http.StatusInternalServerError,
			},
		),
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			rc := recoverer{
				lg: log.GlobalLogger(log.DefaultLoggerName),
				eh: httputil.GlobalErrorHandler(httputil.DefaultErrorHandlerName),
			}

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			ctx := zuid.ContextWithID(req.Context(), "context", tt.C.requestID)
			req = req.WithContext(ctx)

			resp := httptest.NewRecorder()

			h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if tt.C.panics != nil {
					panic(tt.C.panics)
				}
			})
			rc.Middleware(h).ServeHTTP(resp, req)

			testutil.Diff(t, tt.A.statusCode, resp.Code)
		})
	}
}
