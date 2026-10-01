// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package timeout

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
	httputil "github.com/aileron-gateway/aileron-gateway/util/http"
)

type testMatcher struct {
	match bool
}

func (t testMatcher) Match(s string) bool {
	return t.match
}

func TestAPITimeout(t *testing.T) {
	type condition struct {
		timeouts *apiTimeout
		method   string
	}

	type action struct {
		timeout time.Duration
		match   bool
	}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"no methods",
			&condition{
				method: http.MethodGet,
				timeouts: &apiTimeout{
					methods: nil,
					paths:   &testMatcher{match: true},
					timeout: 10,
				},
			},
			&action{
				match:   true,
				timeout: 10,
			},
		),
		gen(
			"match method",
			&condition{
				method: http.MethodGet,
				timeouts: &apiTimeout{
					methods: []string{http.MethodGet},
					paths:   &testMatcher{match: true},
					timeout: 10,
				},
			},
			&action{
				match:   true,
				timeout: 10,
			},
		),
		gen(
			"not match method",
			&condition{
				method: http.MethodGet,
				timeouts: &apiTimeout{
					methods: []string{http.MethodPost},
					paths:   &testMatcher{match: true},
					timeout: 10,
				},
			},
			&action{
				match:   false,
				timeout: 0,
			},
		),
		gen(
			"not match path",
			&condition{
				method: http.MethodGet,
				timeouts: &apiTimeout{
					methods: nil,
					paths:   &testMatcher{match: false},
					timeout: 10,
				},
			},
			&action{
				match:   false,
				timeout: 0,
			},
		),
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			r := httptest.NewRequest(tt.C.method, "http://test.com", nil)
			to, match := tt.C.timeouts.duration(r)

			testutil.Diff(t, tt.A.timeout, to)
			testutil.Diff(t, tt.A.match, match)
		})
	}
}

func TestMiddleware(t *testing.T) {
	type condition struct {
		defaultTimeout time.Duration
		apiTimeouts    []*apiTimeout
		makeTimeout    bool
	}

	type action struct {
		status  int
		timeout int64
	}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"default timeout",
			&condition{
				defaultTimeout: 10 * time.Second,
				makeTimeout:    false,
			},
			&action{
				status:  http.StatusOK,
				timeout: 10,
			},
		),
		gen(
			"single apiTimeout/match",
			&condition{
				defaultTimeout: 10 * time.Second,
				apiTimeouts: []*apiTimeout{
					{
						methods: []string{http.MethodGet},
						paths:   &testMatcher{match: true},
						timeout: 20 * time.Second,
					},
				},
				makeTimeout: false,
			},
			&action{
				status:  http.StatusOK,
				timeout: 20,
			},
		),
		gen(
			"single apiTimeout/not match",
			&condition{
				defaultTimeout: 10 * time.Second,
				apiTimeouts: []*apiTimeout{
					{
						methods: []string{http.MethodGet},
						paths:   &testMatcher{match: false},
						timeout: 20 * time.Second,
					},
				},
				makeTimeout: false,
			},
			&action{
				status:  http.StatusOK,
				timeout: 10,
			},
		),
		gen(
			"multiple *apiTimeout/match none",
			&condition{
				defaultTimeout: 10 * time.Second,
				apiTimeouts: []*apiTimeout{
					{
						methods: []string{http.MethodGet},
						paths:   &testMatcher{match: false},
						timeout: 20 * time.Second,
					},
					{
						methods: []string{http.MethodGet},
						paths:   &testMatcher{match: false},
						timeout: 30 * time.Second,
					},
				},
				makeTimeout: false,
			},
			&action{
				status:  http.StatusOK,
				timeout: 10,
			},
		),
		gen(
			"multiple *apiTimeout/match first",
			&condition{
				defaultTimeout: 10 * time.Second,
				apiTimeouts: []*apiTimeout{
					{
						methods: []string{http.MethodGet},
						paths:   &testMatcher{match: true},
						timeout: 20 * time.Second,
					},
					{
						methods: []string{http.MethodGet},
						paths:   &testMatcher{match: false},
						timeout: 30 * time.Second,
					},
				},
				makeTimeout: false,
			},
			&action{
				status:  http.StatusOK,
				timeout: 20,
			},
		),
		gen(
			"multiple *apiTimeout/match second",
			&condition{
				defaultTimeout: 10 * time.Second,
				apiTimeouts: []*apiTimeout{
					{
						methods: []string{http.MethodGet},
						paths:   &testMatcher{match: false},
						timeout: 20 * time.Second,
					},
					{
						methods: []string{http.MethodGet},
						paths:   &testMatcher{match: true},
						timeout: 30 * time.Second,
					},
				},
				makeTimeout: false,
			},
			&action{
				status:  http.StatusOK,
				timeout: 30,
			},
		),
		gen(
			"timeout occur",
			&condition{
				defaultTimeout: 10 * time.Second,
				apiTimeouts: []*apiTimeout{
					{
						paths:   testMatcher{match: true},
						timeout: time.Second,
					},
				},
				makeTimeout: true,
			},
			&action{
				status:  http.StatusGatewayTimeout,
				timeout: 1,
			},
		),
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			to := timeout{
				defaultTimeout: tt.C.defaultTimeout,
				apiTimeouts:    tt.C.apiTimeouts,
				eh:             httputil.GlobalErrorHandler(httputil.DefaultErrorHandlerName),
			}

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			resp := httptest.NewRecorder()

			h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				deadline, _ := r.Context().Deadline()
				testutil.Diff(t, tt.A.timeout, deadline.Unix()-time.Now().Unix())
				if tt.C.makeTimeout {
					<-r.Context().Done()
				}
			})
			to.Middleware(h).ServeHTTP(resp, req)

			testutil.Diff(t, tt.A.status, resp.Code)
		})
	}
}
