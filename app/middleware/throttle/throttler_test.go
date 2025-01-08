package throttle

import (
	"context"
	"testing"
	"time"

	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
	"github.com/aileron-gateway/aileron-gateway/util/resilience"
)

func TestRetryThrottler(t *testing.T) {
	type condition struct {
		retryThrottler      *retryThrottler
		endContextAt        int
		firstRequestRelease bool
		secondRequest       bool
	}

	type action struct {
		accepted              bool
		secondRequestAccepted bool
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())

	cndNormalRequest := tb.Condition("normal request", "normal request")
	cndRequestWithoutRelease := tb.Condition("request without release", "request without release")
	cndRequestExecuteAfterRelease := tb.Condition("request execute after release", "request execute after release")
	cndResourceReleaseDuringRetryOfRequest := tb.Condition("resource release during retry of a request", "resource release during retry of a request")
	actError := tb.Action("error", "check that an expected value returned")
	actNoError := tb.Action("no error", "check that the there is no error")

	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"normal request",
			[]string{cndNormalRequest},
			[]string{actNoError},
			&condition{
				retryThrottler: &retryThrottler{
					throttler: &testThrottler{
						acceptedAt: 0,
						releaser:   noopReleaser,
					},
					maxRetry: 3,
					waiter:   resilience.NewWaiter(nil),
				},
			},
			&action{
				accepted: true,
			},
		),
		gen(
			"request without release",
			[]string{cndRequestWithoutRelease},
			[]string{actNoError},
			&condition{
				retryThrottler: &retryThrottler{
					throttler: &testThrottler{
						acceptedAt: 0,
						releaser:   noopReleaser,
					},
					maxRetry: 3,
					waiter:   resilience.NewWaiter(nil),
				},
				secondRequest: true,
			},
			&action{
				accepted:              true,
				secondRequestAccepted: false,
			},
		),
		gen(
			"request execute after release",
			[]string{cndRequestExecuteAfterRelease},
			[]string{actNoError},
			&condition{
				retryThrottler: &retryThrottler{
					throttler: &testThrottler{
						acceptedAt: 2,
						releaser:   noopReleaser,
					},
					maxRetry: 3,
					waiter:   resilience.NewWaiter(nil),
				},
				firstRequestRelease: true,
				secondRequest:       true,
			},
			&action{
				accepted:              true,
				secondRequestAccepted: true,
			},
		),
		gen(
			"resource release during retry of a request",
			[]string{cndResourceReleaseDuringRetryOfRequest},
			[]string{actError},
			&condition{
				retryThrottler: &retryThrottler{
					throttler: &testThrottler{
						acceptedAt: 3,
						releaser:   noopReleaser,
					},
					maxRetry: 3,
					waiter:   resilience.NewWaiter(nil),
				},
				endContextAt: 2,
			},
			&action{
				accepted: false,
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			rt := tt.C().retryThrottler
			ctx, cancel := context.WithCancel(context.Background())
			rt.throttler.(*testThrottler).hook = func(count int) {
				if tt.C().endContextAt > 0 && count >= tt.C().endContextAt {
					cancel()
				}
			}

			ok, release := rt.accept(ctx)
			if release != nil {
				defer release()
			}
			if tt.C().firstRequestRelease {
				if release != nil {
					release()
				}
			}
			testutil.Diff(t, tt.A().accepted, ok)

			if tt.C().secondRequest {
				ok, release = rt.accept(ctx)
				if release != nil {
					defer release()
				}
				testutil.Diff(t, tt.A().secondRequestAccepted, ok)
			}
		})
	}
}

func TestMaxConnections(t *testing.T) {
	type condition struct {
		max int
	}

	type action struct {
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"max 0",
			[]string{},
			[]string{},
			&condition{
				max: 0,
			},
			&action{},
		),
		gen(
			"max 1",
			[]string{},
			[]string{},
			&condition{
				max: 1,
			},
			&action{},
		),
		gen(
			"max 5",
			[]string{},
			[]string{},
			&condition{
				max: 5,
			},
			&action{},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			mc := &maxConnections{
				sem: make(chan struct{}, tt.C().max),
			}
			ctx := context.Background()

			for iter := 0; iter < 5; iter++ {
				var releasers []func()
				for i := 0; i < tt.C().max; i++ {
					ok, release := mc.accept(ctx)
					testutil.Diff(t, true, ok)
					releasers = append(releasers, release)
				}
				for i := 0; i < 10; i++ {
					ok, _ := mc.accept(ctx)
					testutil.Diff(t, false, ok)
				}
				for _, release := range releasers {
					release()
				}
			}
		})
	}
}

func TestFixedWindow(t *testing.T) {
	type condition struct {
		size int
	}

	type action struct {
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"size 0",
			[]string{},
			[]string{},
			&condition{
				size: 0,
			},
			&action{},
		),
		gen(
			"size 1",
			[]string{},
			[]string{},
			&condition{
				size: 1,
			},
			&action{},
		),
		gen(
			"size 5",
			[]string{},
			[]string{},
			&condition{
				size: 5,
			},
			&action{},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			fw := &fixedWindow{
				bucket: make(chan struct{}, tt.C().size),
				window: 100 * time.Millisecond,
			}
			ctx := context.Background()

			go fw.fill()
			time.Sleep(150 * time.Millisecond)

			for i := 0; i < tt.C().size; i++ {
				ok, _ := fw.accept(ctx)
				testutil.Diff(t, true, ok)
			}
			for i := 0; i < 10; i++ {
				ok, _ := fw.accept(ctx)
				testutil.Diff(t, false, ok)
			}
		})
	}
}

func TestTokenBucket(t *testing.T) {
	type condition struct {
		size int
		rate int
	}

	type action struct {
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"size=rate=0",
			[]string{},
			[]string{},
			&condition{
				size: 0,
				rate: 0,
			},
			&action{},
		),
		gen(
			"size=rate=1",
			[]string{},
			[]string{},
			&condition{
				size: 1,
				rate: 1,
			},
			&action{},
		),
		gen(
			"size=rate",
			[]string{},
			[]string{},
			&condition{
				size: 5,
				rate: 5,
			},
			&action{},
		),
		gen(
			"size>rate",
			[]string{},
			[]string{},
			&condition{
				size: 10,
				rate: 5,
			},
			&action{},
		),
		gen(
			"size<rate",
			[]string{},
			[]string{},
			&condition{
				size: 5,
				rate: 10,
			},
			&action{},
		),
	}

	testutil.Register(table, testCases...)

	min := func(x, y int) int {
		if x < y {
			return x
		}
		return y
	}

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			tb := &tokenBucket{
				bucket:   make(chan struct{}, tt.C().size),
				rate:     tt.C().rate,
				interval: 100 * time.Millisecond,
			}
			ctx := context.Background()

			go tb.fill()
			time.Sleep(150 * time.Millisecond)

			for i := 0; i < min(tt.C().rate, tt.C().size); i++ {
				ok, _ := tb.accept(ctx)
				testutil.Diff(t, true, ok)
			}
			for i := 0; i < 10; i++ {
				ok, _ := tb.accept(ctx)
				testutil.Diff(t, false, ok)
			}
		})
	}
}

func TestLeakyBucket(t *testing.T) {
	type condition struct {
		size int
		rate int
	}

	type action struct {
		accepted bool
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"size=rate=0",
			[]string{},
			[]string{},
			&condition{
				size: 0,
				rate: 0,
			},
			&action{},
		),
		gen(
			"size=rate=1",
			[]string{},
			[]string{},
			&condition{
				size: 1,
				rate: 1,
			},
			&action{},
		),
		gen(
			"size=rate",
			[]string{},
			[]string{},
			&condition{
				size: 5,
				rate: 5,
			},
			&action{},
		),
		gen(
			"size>rate",
			[]string{},
			[]string{},
			&condition{
				size: 10,
				rate: 5,
			},
			&action{},
		),
		gen(
			"size<rate",
			[]string{},
			[]string{},
			&condition{
				size: 5,
				rate: 10,
			},
			&action{},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			lb := &leakyBucket{
				bucket:   make(chan chan struct{}, tt.C().size),
				rate:     tt.C().rate,
				interval: 200 * time.Millisecond,
			}
			ctx := context.Background()

			go lb.leak()
			time.Sleep(250 * time.Millisecond)

			for i := 0; i < tt.C().rate; i++ { // Immediately accepted.
				ok, _ := lb.accept(ctx)
				testutil.Diff(t, true, ok)
			}
			for i := 0; i < tt.C().size; i++ { // Wait in the bucket.
				go func() {
					ok, _ := lb.accept(ctx)
					testutil.Diff(t, true, ok)
				}()
			}
			time.Sleep(50 * time.Millisecond) // Wait all goroutines done.

			for i := 0; i < 10; i++ { // Overflowed from the bucket.
				ok, _ := lb.accept(ctx)
				testutil.Diff(t, false, ok)
			}

			// Return to initial state.
			time.Sleep(250 * time.Millisecond)

			for i := 0; i < min(tt.C().rate, tt.C().size); i++ { // Immediately accepted.
				ok, _ := lb.accept(ctx)
				testutil.Diff(t, true, ok)
			}
		})
	}
}
