// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package throttle

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
	httputil "github.com/aileron-gateway/aileron-gateway/util/http"
	"github.com/aileron-projects/go/ztime/zrate"
)

type testThrottler struct {
	// acceptedAt is  the number that the
	// accept method return true.
	// 0, 1, 2, 3, ...
	acceptedAt int
	// releaser is the release function returned when the accept method
	// returns true.
	releaser func()
	// hook is the hook function.
	// This function will be called when the accept method was called.
	// The number of counter will be given at the argument.
	hook func(int)

	// counter is the number that the
	// accepted method was called.
	counter int

	// accepted is the bool that the request
	// has been accepted return true.
	accepted bool

	// released is the bool that indicating
	// whether the release func was called.
	released bool
}

func (t *testThrottler) accept(_ context.Context) (bool, func()) {
	defer func() {
		t.counter += 1
	}()
	if t.hook != nil {
		t.hook(t.counter)
	}
	if t.accepted && !t.released {
		return false, nil
	}

	if t.counter >= t.acceptedAt {
		t.accepted = true
		return true, func() {
			t.released = true
		}
	}
	return false, nil
}

type testMatcher struct {
	match bool
}

func (t testMatcher) Match(s string) bool {
	return t.match
}

func TestMiddleware(t *testing.T) {
	type condition struct {
		throttle          throttle
		sendSecondRequest bool
	}

	type action struct {
		resp1Status int
		resp2Status int
	}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"throttler accepts requests",
			&condition{
				throttle: throttle{
					eh: nil,
					throttlers: []*apiThrottler{
						{
							limiter: zrate.NewConcurrentLimiter(1),
							paths:   testMatcher{match: true},
						},
					},
				},
			},
			&action{
				resp1Status: http.StatusOK,
			},
		),
		gen(
			"methods match",
			&condition{
				throttle: throttle{
					eh: nil,
					throttlers: []*apiThrottler{
						{
							limiter: zrate.NewConcurrentLimiter(1),
							methods: []string{http.MethodGet},
							paths:   testMatcher{match: true},
						},
					},
				},
			},
			&action{
				resp1Status: http.StatusOK,
			},
		),
		gen(
			"methods not match",
			&condition{
				throttle: throttle{
					eh: nil,
					throttlers: []*apiThrottler{
						{
							limiter: zrate.NewConcurrentLimiter(0),
							methods: []string{http.MethodPost},
							paths:   testMatcher{match: true},
						},
					},
				},
			},
			&action{
				resp1Status: http.StatusOK,
			},
		),
		gen(
			"paths not match",
			&condition{
				throttle: throttle{
					eh: nil,
					throttlers: []*apiThrottler{
						{
							limiter: zrate.NewConcurrentLimiter(0),
							methods: []string{http.MethodGet},
							paths:   testMatcher{match: false},
						},
					},
				},
			},
			&action{
				resp1Status: http.StatusOK,
			},
		),
		gen(
			"no throttler",
			&condition{},
			&action{
				resp1Status: http.StatusOK,
			},
		),
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			th := throttle{
				eh:         httputil.GlobalErrorHandler(httputil.DefaultErrorHandlerName),
				throttlers: tt.C.throttle.throttlers,
			}
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			resp := httptest.NewRecorder()
			h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})
			th.Middleware(h).ServeHTTP(resp, req)

			testutil.Diff(t, tt.A.resp1Status, resp.Code)

			if tt.C.sendSecondRequest {
				resp = httptest.NewRecorder()
				th.Middleware(h).ServeHTTP(resp, req)
				testutil.Diff(t, tt.A.resp2Status, resp.Code)
			}
		})
	}
}
