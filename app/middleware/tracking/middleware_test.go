// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package tracking

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aileron-gateway/aileron-gateway/app"
	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
	"github.com/aileron-gateway/aileron-gateway/kernel/uid"
	utilhttp "github.com/aileron-gateway/aileron-gateway/util/http"
)

// Mock functions to simulate NewRequestID and NewTraceID behavior.
var (
	mockNewRequestID = func() (string, error) {
		return "mock-request-id", nil
	}
	mockNewTraceID = func(_ string) (string, error) {
		return "mock-trace-id", nil
	}
)

func TestMiddleware(t *testing.T) {
	type condition struct {
		trcExtractHeader string
		requestID        string
		traceID          string
		newReqIDFunc     func() (string, error)
		newTrcIDFunc     func(string) (string, error)
	}

	type action struct {
		statusCode int
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"successful request with no errors",
			[]string{},
			[]string{},
			&condition{
				requestID:    "mock-request-id",
				traceID:      "mock-trace-id",
				newReqIDFunc: mockNewRequestID,
				newTrcIDFunc: mockNewTraceID,
			},
			&action{
				statusCode: http.StatusOK,
			},
		),
		gen(
			"trace header provided in request",
			[]string{},
			[]string{},
			&condition{
				requestID:        "mock-request-id",
				traceID:          "existing-trace-id",
				trcExtractHeader: "X-Trace-ID",
				newReqIDFunc:     mockNewRequestID,
				newTrcIDFunc:     mockNewTraceID,
			},
			&action{
				statusCode: http.StatusOK,
			},
		),
		gen(
			"error generating request ID",
			[]string{},
			[]string{},
			&condition{
				requestID: "",
				traceID:   "",
				newReqIDFunc: func() (string, error) {
					return "", app.ErrAppMiddleGenID.WithStack(nil, map[string]any{"type": "request"})
				},
				newTrcIDFunc: mockNewTraceID,
			},
			&action{
				statusCode: http.StatusInternalServerError,
			},
		),
		gen(
			"error generating trace ID",
			[]string{},
			[]string{},
			&condition{
				requestID:    "",
				traceID:      "",
				newReqIDFunc: mockNewRequestID,
				newTrcIDFunc: func(s string) (string, error) {
					return "", app.ErrAppMiddleGenID.WithStack(nil, map[string]any{"type": "trace"})
				},
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
			tr := tracker{
				eh:               utilhttp.GlobalErrorHandler(utilhttp.DefaultErrorHandlerName),
				reqProxyHeader:   "X-Request-ID",
				trcProxyHeader:   "X-Trace-ID",
				trcExtractHeader: tt.C().trcExtractHeader,
				newReqID:         tt.C().newReqIDFunc,
				newTrcID:         tt.C().newTrcIDFunc,
			}

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			if tt.C().trcExtractHeader != "" && tt.C().traceID != "" {
				req.Header.Set(tt.C().trcExtractHeader, tt.C().traceID)
			}

			resp := httptest.NewRecorder()

			h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				testutil.Diff(t, tt.C().requestID, uid.IDFromContext(r.Context()))
				if tt.C().traceID != "" {
					h := utilhttp.ProxyHeaderFromContext(r.Context())
					testutil.Diff(t, tt.C().traceID, h.Get("X-Trace-ID"))
				}
				w.WriteHeader(http.StatusOK)
			})
			tr.Middleware(h).ServeHTTP(resp, req)

			testutil.Diff(t, tt.A().statusCode, resp.Code)
		})
	}
}
