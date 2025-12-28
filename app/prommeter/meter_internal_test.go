// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package prommeter

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aileron-gateway/aileron-gateway/core"
	"github.com/aileron-gateway/aileron-gateway/internal/testutil"
	"github.com/google/go-cmp/cmp"
	"github.com/prometheus/client_golang/prometheus"
	promtestutil "github.com/prometheus/client_golang/prometheus/testutil"
)

func TestRegistry(t *testing.T) {
	type condition struct {
		reg *prometheus.Registry
	}

	type action struct {
	}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"nil registry",
			&condition{
				reg: nil,
			},
			&action{},
		),
		gen(
			"non nil registry",
			&condition{
				reg: prometheus.NewRegistry(),
			},
			&action{},
		),
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			metrics := &metrics{
				reg: tt.C.reg,
			}

			reg := metrics.Registry()

			opts := []cmp.Option{
				cmp.Comparer(testutil.ComparePointer[*prometheus.Registry]),
			}
			testutil.Diff(t, tt.C.reg, reg, opts...)
		})
	}
}

func TestMiddleware(t *testing.T) {
	type condition struct {
		requestCount int
		statusCode   int
	}

	type action struct {
		value float64
	}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"single request",
			&condition{
				requestCount: 1,
				statusCode:   200,
			},
			&action{
				value: 1,
			},
		),
		gen(
			"multiple requests",
			&condition{
				requestCount: 5,
				statusCode:   200,
			},
			&action{
				value: 5,
			},
		),
		gen(
			"single request with 500 error",
			&condition{
				requestCount: 1,
				statusCode:   500,
			},
			&action{
				value: 1,
			},
		),
		gen(
			"single request with 404 error",
			&condition{
				requestCount: 1,
				statusCode:   404,
			},
			&action{
				value: 1,
			},
		),
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			metrics := &metrics{
				mAPICalls: prometheus.NewCounterVec(
					prometheus.CounterOpts{
						Name: "http_requests_total",
						Help: "Total number of received http requests",
					},
					[]string{"host", "path", "code", "method"},
				),
			}

			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			resp := httptest.NewRecorder()
			h := metrics.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.C.statusCode)
			}))

			for i := 0; i < tt.C.requestCount; i++ {
				h.ServeHTTP(resp, req)
			}

			v := promtestutil.ToFloat64(metrics.mAPICalls)
			testutil.Diff(t, tt.A.value, v)
		})
	}
}

func TestTripperware(t *testing.T) {
	type condition struct {
		requestCount int
		responseCode int
		returnError  bool
	}

	type action struct {
		value float64
	}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"single request",
			&condition{
				requestCount: 1,
				responseCode: 200,
			},
			&action{
				value: 1,
			},
		),
		gen(
			"multiple requests",
			&condition{
				requestCount: 5,
				responseCode: 200,
			},
			&action{
				value: 5,
			},
		),
		gen(
			"single request with 404 error",
			&condition{
				requestCount: 1,
				responseCode: 404,
			},
			&action{
				value: 1,
			},
		),
		gen(
			"single request with network error",
			&condition{
				requestCount: 1,
				returnError:  true,
			},
			&action{
				value: 1,
			},
		),
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			metrics := &metrics{
				tAPICalls: prometheus.NewCounterVec(
					prometheus.CounterOpts{
						Name: "http_client_requests_total",
						Help: "Total number of sent http requests",
					},
					[]string{"host", "path", "code", "method"},
				),
			}

			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			h := metrics.Tripperware(core.RoundTripperFunc(func(r *http.Request) (*http.Response, error) {
				if tt.C.returnError {
					return nil, http.ErrHandlerTimeout
				}
				return &http.Response{
					StatusCode: tt.C.responseCode,
					Request:    r,
				}, nil
			}))

			for i := 0; i < tt.C.requestCount; i++ {
				h.RoundTrip(req)
			}

			v := promtestutil.ToFloat64(metrics.tAPICalls)
			testutil.Diff(t, tt.A.value, v)
		})
	}
}
