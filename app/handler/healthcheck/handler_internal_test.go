// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package healthcheck

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/aileron-gateway/aileron-gateway/app"
	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
	httputil "github.com/aileron-gateway/aileron-gateway/util/http"
)

type healthChecker struct {
	timeout time.Duration
	status  bool
}

func (m *healthChecker) HealthCheck(ctx context.Context) (context.Context, bool) {
	if m.timeout > 0 {
		time.Sleep(m.timeout)
	}
	if m.status {
		return ctx, true
	}
	return ctx, false
}

func TestServeHTTP(t *testing.T) {
	type condition struct {
		timeout  time.Duration
		checkers []app.HealthChecker
	}

	type action struct {
		statusCode int
		body       string
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())

	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"nil checker",
			[]string{},
			[]string{},
			&condition{
				timeout:  time.Second,
				checkers: nil,
			},
			&action{
				statusCode: http.StatusOK,
				body:       "{}",
			},
		),
		gen(
			"single checker succeeded",
			[]string{},
			[]string{},
			&condition{
				timeout: time.Second,
				checkers: []app.HealthChecker{
					&healthChecker{
						status: true,
					},
				},
			},
			&action{
				statusCode: http.StatusOK,
				body:       "{}",
			},
		),
		gen(
			"single checker failed",
			[]string{},
			[]string{},
			&condition{
				timeout: time.Second,
				checkers: []app.HealthChecker{
					&healthChecker{
						status: false,
					},
				},
			},
			&action{
				statusCode: http.StatusInternalServerError,
				body:       `{"status":500,"statusText":"Internal Server Error"}`,
			},
		),
		gen(
			"single checker timeout",
			[]string{},
			[]string{},
			&condition{
				timeout: 1 * time.Millisecond,
				checkers: []app.HealthChecker{
					&healthChecker{
						timeout: 100 * time.Millisecond,
						status:  true,
					},
				},
			},
			&action{
				statusCode: http.StatusGatewayTimeout,
				body:       `{"status":504,"statusText":"Gateway Timeout"}`,
			},
		),
		gen(
			"multiple checkers succeeded",
			[]string{},
			[]string{},
			&condition{
				timeout: time.Second,
				checkers: []app.HealthChecker{
					&healthChecker{
						status: true,
					},
					&healthChecker{
						status: true,
					},
				},
			},
			&action{
				statusCode: http.StatusOK,
				body:       "{}",
			},
		),
		gen(
			"multiple checkers failed",
			[]string{},
			[]string{},
			&condition{
				timeout: time.Second,
				checkers: []app.HealthChecker{
					&healthChecker{
						status: true,
					},
					&healthChecker{
						status: false,
					},
				},
			},
			&action{
				statusCode: http.StatusInternalServerError,
				body:       `{"status":500,"statusText":"Internal Server Error"}`,
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			hc := healthCheck{
				eh:       httputil.GlobalErrorHandler(httputil.DefaultErrorHandlerName),
				timeout:  tt.C().timeout,
				checkers: tt.C().checkers,
			}

			req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
			resp := httptest.NewRecorder()

			hc.ServeHTTP(resp, req)

			testutil.Diff(t, tt.A().statusCode, resp.Code)
			testutil.Diff(t, "application/json; charset=utf-8", resp.Header().Get("Content-Type"))
			testutil.Diff(t, "nosniff", resp.Header().Get("X-Content-Type-Options"))
			testutil.Diff(t, tt.A().body, resp.Body.String())
		})
	}
}
