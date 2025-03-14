package throttle

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
	httputil "github.com/aileron-gateway/aileron-gateway/util/http"
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

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"throttler accepts requests",
			[]string{},
			[]string{},
			&condition{
				throttle: throttle{
					eh: nil,
					throttlers: []*apiThrottler{
						{
							throttler: &testThrottler{
								acceptedAt: 0,
								releaser:   noopReleaser,
							},
							paths: testMatcher{match: true},
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
			[]string{},
			[]string{},
			&condition{
				throttle: throttle{
					eh: nil,
					throttlers: []*apiThrottler{
						{
							throttler: &testThrottler{
								acceptedAt: 0,
								releaser:   noopReleaser,
							},
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
			[]string{},
			[]string{},
			&condition{
				throttle: throttle{
					eh: nil,
					throttlers: []*apiThrottler{
						{
							throttler: &testThrottler{
								acceptedAt: 999999999, // Never used.
							},
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
			[]string{},
			[]string{},
			&condition{
				throttle: throttle{
					eh: nil,
					throttlers: []*apiThrottler{
						{
							throttler: &testThrottler{
								acceptedAt: 999999999, // Never used.
							},
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
			[]string{},
			[]string{},
			&condition{},
			&action{
				resp1Status: http.StatusOK,
			},
		),
		gen(
			"first request is fail, second request is success",
			[]string{},
			[]string{},
			&condition{
				throttle: throttle{
					throttlers: []*apiThrottler{
						{
							throttler: &testThrottler{
								acceptedAt: 1,
								releaser:   noopReleaser,
							},
							paths: testMatcher{match: true},
						},
					},
				},
				sendSecondRequest: true,
			},
			&action{
				resp1Status: http.StatusTooManyRequests,
				resp2Status: http.StatusOK,
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			th := throttle{
				eh:         httputil.GlobalErrorHandler(httputil.DefaultErrorHandlerName),
				throttlers: tt.C().throttle.throttlers,
			}
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			resp := httptest.NewRecorder()
			h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})
			th.Middleware(h).ServeHTTP(resp, req)

			testutil.Diff(t, tt.A().resp1Status, resp.Code)

			if tt.C().sendSecondRequest {
				resp = httptest.NewRecorder()
				th.Middleware(h).ServeHTTP(resp, req)
				testutil.Diff(t, tt.A().resp2Status, resp.Code)
			}
		})
	}
}
