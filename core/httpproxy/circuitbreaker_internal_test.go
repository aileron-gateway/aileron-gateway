// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package httpproxy

import (
	"fmt"
	"sync"
	"testing"
	"time"

	v1 "github.com/aileron-gateway/aileron-gateway/apis/core/v1"
	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

type counterMethod int

const (
	success counterMethod = iota
	failure
	reset
)

type testTimer struct {
	current int
	times   []int64
}

func (t *testTimer) unixMilli() int64 {
	tt := t.times[t.current]
	t.current += 1
	return tt
}

func TestNewCircuitBreaker(t *testing.T) {
	type condition struct {
		spec *v1.CircuitBreaker
	}

	type action struct {
		cb circuitBreaker
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	cndNil := tb.Condition("input nil", "input nil spec")
	cndConsecutive := tb.Condition("consecutive", "configure consecutive")
	cndCountFixed := tb.Condition("countBasedFixedWindow", "configure countBasedFixedWindow")
	cndTimeFixed := tb.Condition("timeBasedFixedWindow", "configure timeBasedFixedWindow")
	cndCountSlid := tb.Condition("countBasedSlidingWindow", "configure countBasedSlidingWindow")
	cndTimeSlid := tb.Condition("timeBasedSlidingWindow", "configure timeBasedSlidingWindow")
	actCheckCounter := tb.Action("check counter", "check the returned counter was configured with input values")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"nil",
			[]string{cndNil},
			[]string{actCheckCounter},
			&condition{
				spec: nil,
			},
			&action{
				cb: nil,
			},
		),
		gen(
			"consecutive",
			[]string{cndConsecutive},
			[]string{actCheckCounter},
			&condition{
				spec: &v1.CircuitBreaker{
					FailureThreshold:        10,
					SuccessThreshold:        20,
					EffectiveFailureSamples: 30,
					EffectiveSuccessSamples: 40,
					WaitDuration:            50,
					CircuitBreakerCounter:   &v1.CircuitBreaker_ConsecutiveCounter{},
				},
			},
			&action{
				cb: &circuitBreakerController{
					status:           closed,
					counter:          &consecutiveCounter{},
					failureThreshold: 10,
					successThreshold: 20,
					wait:             50 * time.Second,
				},
			},
		),
		gen(
			"countBasedFixedWindow",
			[]string{cndCountFixed},
			[]string{actCheckCounter},
			&condition{
				spec: &v1.CircuitBreaker{
					FailureThreshold:        10,
					SuccessThreshold:        20,
					EffectiveFailureSamples: 30,
					EffectiveSuccessSamples: 40,
					WaitDuration:            50,
					CircuitBreakerCounter: &v1.CircuitBreaker_CountBasedFixedWindowCounter{
						CountBasedFixedWindowCounter: &v1.CountBasedFixedWindowCounterSpec{
							Samples: 60,
						},
					},
				},
			},
			&action{
				cb: &circuitBreakerController{
					status: closed,
					counter: &countBasedFixedWindowCounter{
						maxTotal:                60,
						effectiveFailureSamples: 30,
						effectiveSuccessSamples: 40,
					},
					failureThreshold: 10,
					successThreshold: 20,
					wait:             50 * time.Second,
				},
			},
		),
		gen(
			"timeBasedFixedWindow",
			[]string{cndTimeFixed},
			[]string{actCheckCounter},
			&condition{
				spec: &v1.CircuitBreaker{
					FailureThreshold:        10,
					SuccessThreshold:        20,
					EffectiveFailureSamples: 30,
					EffectiveSuccessSamples: 40,
					WaitDuration:            50,
					CircuitBreakerCounter: &v1.CircuitBreaker_TimeBasedFixedWindowCounter{
						TimeBasedFixedWindowCounter: &v1.TimeBasedFixedWindowCounterSpec{
							WindowWidth: 60,
						},
					},
				},
			},
			&action{
				cb: &circuitBreakerController{
					status: closed,
					counter: &timeBasedFixedWindowCounter{
						windowWidth:             60 * 1000, // second to millisecond
						effectiveFailureSamples: 30,
						effectiveSuccessSamples: 40,
						lastResetTime:           time.Unix(0, 0).UnixMilli(), // timer is set to epoch for this test.
					},
					failureThreshold: 10,
					successThreshold: 20,
					wait:             50 * time.Second,
				},
			},
		),
		gen(
			"countBasedSlidingWindow",
			[]string{cndCountSlid},
			[]string{actCheckCounter},
			&condition{
				spec: &v1.CircuitBreaker{
					FailureThreshold:        10,
					SuccessThreshold:        20,
					EffectiveFailureSamples: 30,
					EffectiveSuccessSamples: 40,
					WaitDuration:            50,
					CircuitBreakerCounter: &v1.CircuitBreaker_CountBasedSlidingWindowCounter{
						CountBasedSlidingWindowCounter: &v1.CountBasedSlidingWindowCounterSpec{
							Samples:      60,
							HistoryLimit: 70,
						},
					},
				},
			},
			&action{
				cb: &circuitBreakerController{
					status: closed,
					counter: &countBasedSlidingWindowCounter{
						maxTotal:                60,
						effectiveFailureSamples: 30,
						effectiveSuccessSamples: 40,
						totalHistory:            make([]int, 70),
						failureCountHistory:     make([]int, 70),
						successCountHistory:     make([]int, 70),
					},
					failureThreshold: 10,
					successThreshold: 20,
					wait:             50 * time.Second,
				},
			},
		),
		gen(
			"timeBasedSlidingWindow",
			[]string{cndTimeSlid},
			[]string{actCheckCounter},
			&condition{
				spec: &v1.CircuitBreaker{
					FailureThreshold:        10,
					SuccessThreshold:        20,
					EffectiveFailureSamples: 30,
					EffectiveSuccessSamples: 40,
					WaitDuration:            50,
					CircuitBreakerCounter: &v1.CircuitBreaker_TimeBasedSlidingWindowCounter{
						TimeBasedSlidingWindowCounter: &v1.TimeBasedSlidingWindowCounterSpec{
							WindowWidth:  60,
							HistoryLimit: 70,
						},
					},
				},
			},
			&action{
				cb: &circuitBreakerController{
					status: closed,
					counter: &timeBasedSlidingWindowCounter{
						windowWidth:             60 * 1000, // second to millisecond
						effectiveFailureSamples: 30,
						effectiveSuccessSamples: 40,
						totalHistory:            make([]int, 70),
						failureCountHistory:     make([]int, 70),
						successCountHistory:     make([]int, 70),
						lastResetTime:           time.Unix(0, 0).UnixMilli(), // timer is set to epoch for this test.
					},
					failureThreshold: 10,
					successThreshold: 20,
					wait:             50 * time.Second,
				},
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			// Replace timer for testing.
			tmp := unixMilli
			unixMilli = func() int64 { return time.Unix(0, 0).UnixMilli() }
			defer func() {
				unixMilli = tmp
			}()

			cb := newCircuitBreaker(tt.C().spec)

			opts := []cmp.Option{
				cmp.AllowUnexported(circuitBreakerController{}),
				cmp.AllowUnexported(consecutiveCounter{}),
				cmp.AllowUnexported(countBasedFixedWindowCounter{}),
				cmp.AllowUnexported(timeBasedFixedWindowCounter{}),
				cmp.AllowUnexported(countBasedSlidingWindowCounter{}),
				cmp.AllowUnexported(timeBasedSlidingWindowCounter{}),
				cmpopts.IgnoreUnexported(sync.RWMutex{}),
			}
			testutil.Diff(t, tt.A().cb, cb, opts...)
		})
	}
}

func TestCircuitBreakerController(t *testing.T) {
	type condition struct {
		cb          *circuitBreakerController
		callMethods []counterMethod
	}

	type action struct {
		status circuitBreakerStatus
		active bool
	}

	cndInitialClosed := "initial closed"
	cndInitialHalfOpened := "initial halfOpened"
	cndInitialOpened := "initial opened"
	cndChangeToClosed := "change to closed"
	cndChangeToHalfOpened := "change to halfOpened"
	cndChangeToOpened := "change to opened"
	actCheckStatus := "check status"
	actCheckActive := "check active"

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	tb.Condition(cndInitialClosed, "success method was called")
	tb.Condition(cndInitialHalfOpened, "failure method was called")
	tb.Condition(cndInitialOpened, "reset method was called")
	tb.Condition(cndChangeToClosed, "success method was called")
	tb.Condition(cndChangeToHalfOpened, "failure method was called")
	tb.Condition(cndChangeToOpened, "reset method was called")
	tb.Action(actCheckStatus, "check the final count of success")
	tb.Action(actCheckActive, "check the final count of failure")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"all success (closed > closed)",
			[]string{cndInitialClosed, cndChangeToClosed},
			[]string{actCheckStatus, actCheckActive},
			&condition{
				cb: &circuitBreakerController{
					status:           closed,
					counter:          &consecutiveCounter{},
					successThreshold: 2,
					failureThreshold: 2,
				},
				callMethods: []counterMethod{
					success, success, success, success, success,
				},
			},
			&action{
				status: closed,
				active: true,
			},
		),
		gen(
			"failed but not over threshold (closed > closed)",
			[]string{cndInitialClosed, cndChangeToClosed},
			[]string{actCheckStatus, actCheckActive},
			&condition{
				cb: &circuitBreakerController{
					status:           closed,
					counter:          &consecutiveCounter{},
					successThreshold: 2,
					failureThreshold: 2,
				},
				callMethods: []counterMethod{
					success, success, success, success, failure,
				},
			},
			&action{
				status: closed,
				active: true,
			},
		),
		gen(
			"failed over threshold (closed > opened)",
			[]string{cndInitialClosed, cndChangeToOpened},
			[]string{actCheckStatus, actCheckActive},
			&condition{
				cb: &circuitBreakerController{
					status:           closed,
					counter:          &consecutiveCounter{},
					successThreshold: 2,
					failureThreshold: 2,
					wait:             time.Second, // avoid to become halfOpened
				},
				callMethods: []counterMethod{
					success, success, success, failure, failure,
				},
			},
			&action{
				status: opened,
				active: false,
			},
		),
		gen(
			"recover from opened to half-opened (opened > halfOpened)",
			[]string{cndInitialOpened, cndChangeToHalfOpened},
			[]string{actCheckStatus, actCheckActive},
			&condition{
				cb: &circuitBreakerController{
					status:           closed,
					counter:          &consecutiveCounter{},
					successThreshold: 2,
					failureThreshold: 2,
					wait:             time.Millisecond,
				},
				callMethods: []counterMethod{
					success, success, success, failure, failure,
				},
			},
			&action{
				status: halfOpened,
				active: true,
			},
		),
		gen(
			"fail on halfOpened (halfOpened > opened)",
			[]string{cndInitialHalfOpened, cndChangeToOpened},
			[]string{actCheckStatus, actCheckActive},
			&condition{
				cb: &circuitBreakerController{
					status:           halfOpened,
					counter:          &consecutiveCounter{},
					successThreshold: 2,
					failureThreshold: 2,
					wait:             time.Second, // avoid to become halfOpened
				},
				callMethods: []counterMethod{
					failure,
				},
			},
			&action{
				status: opened,
				active: false,
			},
		),
		gen(
			"success on halfOpened (halfOpened > closed)",
			[]string{cndInitialHalfOpened, cndChangeToClosed},
			[]string{actCheckStatus, actCheckActive},
			&condition{
				cb: &circuitBreakerController{
					status:           halfOpened,
					counter:          &consecutiveCounter{},
					successThreshold: 2,
					failureThreshold: 2,
				},
				callMethods: []counterMethod{
					success, success,
				},
			},
			&action{
				status: closed,
				active: true,
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			cb := tt.C().cb
			for _, m := range tt.C().callMethods {
				switch m {
				case success:
					cb.countSuccess()
				case failure:
					cb.countFailure()
				}
			}

			time.Sleep(100 * time.Millisecond)

			testutil.Diff(t, tt.A().status, cb.status)
			testutil.Diff(t, tt.A().active, cb.Active())
		})
	}
}

func TestCircuitBreakerController_changeStatus(t *testing.T) {
	// This test is to take coverage.

	type condition struct {
		cb     *circuitBreakerController
		status circuitBreakerStatus
	}

	type action struct {
		status circuitBreakerStatus
		active bool
	}

	cndClosed := "closed to closed"
	cndOpened := "opened to opened"
	cndHalfOpened := "halfOpened to halfOpened"
	actCheckStatus := "check final status"

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	tb.Condition(cndClosed, "success method was called")
	tb.Condition(cndOpened, "failure method was called")
	tb.Condition(cndHalfOpened, "reset method was called")
	tb.Action(actCheckStatus, "check the final count of success")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"closed",
			[]string{cndClosed},
			[]string{actCheckStatus},
			&condition{
				cb: &circuitBreakerController{
					status: closed,
				},
				status: closed, // closed > closed
			},
			&action{
				status: closed,
				active: true,
			},
		),
		gen(
			"opened",
			[]string{cndOpened},
			[]string{actCheckStatus},
			&condition{
				cb: &circuitBreakerController{
					status: opened,
				},
				status: opened, // opened > opened
			},
			&action{
				status: opened,
				active: false,
			},
		),
		gen(
			"halfOpened",
			[]string{cndHalfOpened},
			[]string{actCheckStatus},
			&condition{
				cb: &circuitBreakerController{
					status: halfOpened,
				},
				status: halfOpened, // halfOpened > halfOpened
			},
			&action{
				status: halfOpened,
				active: true,
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			cb := tt.C().cb
			cb.changeStatus(tt.C().status)

			testutil.Diff(t, tt.A().status, cb.status)
			testutil.Diff(t, tt.A().active, cb.Active())
		})
	}
}

func TestConsecutiveCounter(t *testing.T) {
	type condition struct {
		counter     *consecutiveCounter
		callMethods []counterMethod
	}

	type action struct {
		successCount int
		failureCount int
		result       int
	}

	cndSuccess := "success"
	cndFailure := "failure"
	cndReset := "reset"
	actCheckSuccessCount := "check success count"
	actCheckFailureCount := "check failure count"

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	tb.Condition(cndSuccess, "success method was called")
	tb.Condition(cndFailure, "failure method was called")
	tb.Condition(cndReset, "reset method was called")
	tb.Action(actCheckSuccessCount, "check the final count of success")
	tb.Action(actCheckFailureCount, "check the final count of failure")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"no call",
			[]string{},
			[]string{actCheckSuccessCount, actCheckFailureCount},
			&condition{
				counter:     &consecutiveCounter{},
				callMethods: []counterMethod{},
			},
			&action{
				successCount: 0,
				failureCount: 0,
				result:       0,
			},
		),
		gen(
			"success only",
			[]string{cndSuccess},
			[]string{actCheckSuccessCount, actCheckFailureCount},
			&condition{
				counter:     &consecutiveCounter{},
				callMethods: []counterMethod{success, success, success},
			},
			&action{
				successCount: 3,
				failureCount: 0,
				result:       3,
			},
		),
		gen(
			"failure only",
			[]string{cndFailure},
			[]string{actCheckSuccessCount, actCheckFailureCount},
			&condition{
				counter:     &consecutiveCounter{},
				callMethods: []counterMethod{failure, failure, failure},
			},
			&action{
				successCount: 0,
				failureCount: 3,
				result:       3,
			},
		),
		gen(
			"finally successful",
			[]string{cndSuccess, cndFailure},
			[]string{actCheckSuccessCount, actCheckFailureCount},
			&condition{
				counter:     &consecutiveCounter{},
				callMethods: []counterMethod{failure, failure, success, success},
			},
			&action{
				successCount: 2,
				failureCount: 0,
				result:       2,
			},
		),
		gen(
			"finally failure",
			[]string{cndSuccess, cndFailure},
			[]string{actCheckSuccessCount, actCheckFailureCount},
			&condition{
				counter:     &consecutiveCounter{},
				callMethods: []counterMethod{success, success, failure, failure},
			},
			&action{
				successCount: 0,
				failureCount: 2,
				result:       2,
			},
		),
		gen(
			"reset",
			[]string{cndSuccess, cndFailure, cndReset},
			[]string{actCheckSuccessCount, actCheckFailureCount},
			&condition{
				counter:     &consecutiveCounter{},
				callMethods: []counterMethod{success, failure, reset},
			},
			&action{
				successCount: 0,
				failureCount: 0,
				result:       0,
			},
		),
		gen(
			"success after reset",
			[]string{cndSuccess, cndReset},
			[]string{actCheckSuccessCount, actCheckFailureCount},
			&condition{
				counter:     &consecutiveCounter{},
				callMethods: []counterMethod{reset, success, success},
			},
			&action{
				successCount: 2,
				failureCount: 0,
				result:       2,
			},
		),
		gen(
			"failure after reset",
			[]string{cndFailure, cndReset},
			[]string{actCheckSuccessCount, actCheckFailureCount},
			&condition{
				counter:     &consecutiveCounter{},
				callMethods: []counterMethod{reset, failure, failure},
			},
			&action{
				successCount: 0,
				failureCount: 2,
				result:       2,
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			counter := tt.C().counter

			result := 0
			for _, m := range tt.C().callMethods {
				switch m {
				case success:
					result = counter.success()
				case failure:
					result = counter.failure()
				case reset:
					counter.reset()
					result = 0
				}
			}

			testutil.Diff(t, tt.A().successCount, counter.successCount)
			testutil.Diff(t, tt.A().failureCount, counter.failureCount)
			testutil.Diff(t, tt.A().result, result)
		})
	}
}

func TestCountBasedFixedWindowCounter(t *testing.T) {
	type condition struct {
		counter     *countBasedFixedWindowCounter
		callMethods []counterMethod
	}

	type action struct {
		successCount int
		failureCount int
		totalCount   int
		result       int
	}

	cndWindowSizeZero := "window size 0"
	cndNonZeroEffectiveSuccess := "non-zero effective success"
	cndNonZeroEffectiveFailure := "non-zero effective failure"
	cndSuccess := "success"
	cndFailure := "failure"
	cndReset := "reset"
	actCheckSuccessCount := "check success count"
	actCheckFailureCount := "check failure count"

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	tb.Condition(cndWindowSizeZero, "window size is zero")
	tb.Condition(cndNonZeroEffectiveSuccess, "effective success count is not zero")
	tb.Condition(cndNonZeroEffectiveFailure, "effective failure count is not zero")
	tb.Condition(cndSuccess, "success method was called")
	tb.Condition(cndFailure, "failure method was called")
	tb.Condition(cndReset, "reset method was called")
	tb.Action(actCheckSuccessCount, "check the final count of success")
	tb.Action(actCheckFailureCount, "check the final count of failure")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"window size is 0, failure for the last",
			[]string{cndWindowSizeZero, cndSuccess, cndFailure},
			[]string{actCheckSuccessCount, actCheckFailureCount},
			&condition{
				counter: &countBasedFixedWindowCounter{
					maxTotal: 0,
				},
				callMethods: []counterMethod{success, success, failure, failure},
			},
			&action{
				totalCount:   1,
				successCount: 0,
				failureCount: 1,
				result:       100,
			},
		),
		gen(
			"window size is 0, success for the last",
			[]string{cndWindowSizeZero, cndSuccess, cndFailure},
			[]string{actCheckSuccessCount, actCheckFailureCount},
			&condition{
				counter: &countBasedFixedWindowCounter{
					maxTotal: 0,
				},
				callMethods: []counterMethod{failure, failure, success, success},
			},
			&action{
				totalCount:   1,
				successCount: 1,
				failureCount: 0,
				result:       100,
			},
		),
		gen(
			"call count is less than window size, failure for the last",
			[]string{cndSuccess, cndFailure},
			[]string{actCheckSuccessCount, actCheckFailureCount},
			&condition{
				counter: &countBasedFixedWindowCounter{
					maxTotal: 4,
				},
				callMethods: []counterMethod{success, success, failure, failure},
			},
			&action{
				totalCount:   4,
				successCount: 2,
				failureCount: 2,
				result:       50,
			},
		),
		gen(
			"call count is less than window size, success for the last",
			[]string{cndSuccess, cndFailure},
			[]string{actCheckSuccessCount, actCheckFailureCount},
			&condition{
				counter: &countBasedFixedWindowCounter{
					maxTotal: 4,
				},
				callMethods: []counterMethod{failure, failure, success, success},
			},
			&action{
				totalCount:   4,
				successCount: 2,
				failureCount: 2,
				result:       50,
			},
		),
		gen(
			"call count is more than window size, failure for the last",
			[]string{cndSuccess, cndFailure},
			[]string{actCheckSuccessCount, actCheckFailureCount},
			&condition{
				counter: &countBasedFixedWindowCounter{
					maxTotal: 4,
				},
				callMethods: []counterMethod{success, success, failure, failure, success, failure},
			},
			&action{
				totalCount:   2,
				successCount: 1,
				failureCount: 1,
				result:       50,
			},
		),
		gen(
			"call count is less than window size, success for the last",
			[]string{cndSuccess, cndFailure},
			[]string{actCheckSuccessCount, actCheckFailureCount},
			&condition{
				counter: &countBasedFixedWindowCounter{
					maxTotal: 4,
				},
				callMethods: []counterMethod{failure, failure, success, success},
			},
			&action{
				totalCount:   4,
				successCount: 2,
				failureCount: 2,
				result:       50,
			},
		),
		gen(
			"call count is more than window size, success for the last",
			[]string{cndSuccess, cndFailure},
			[]string{actCheckSuccessCount, actCheckFailureCount},
			&condition{
				counter: &countBasedFixedWindowCounter{
					maxTotal: 4,
				},
				callMethods: []counterMethod{failure, failure, success, success, failure, success},
			},
			&action{
				totalCount:   2,
				successCount: 1,
				failureCount: 1,
				result:       50,
			},
		),
		gen(
			"reset after count",
			[]string{cndSuccess, cndFailure, cndReset},
			[]string{actCheckSuccessCount, actCheckFailureCount},
			&condition{
				counter: &countBasedFixedWindowCounter{
					maxTotal: 4,
				},
				callMethods: []counterMethod{failure, failure, success, success, reset},
			},
			&action{
				totalCount:   0,
				successCount: 0,
				failureCount: 0,
				result:       0,
			},
		),
		gen(
			"0 for less than the effectiveSuccessSamples",
			[]string{cndNonZeroEffectiveSuccess},
			[]string{actCheckSuccessCount, actCheckFailureCount},
			&condition{
				counter: &countBasedFixedWindowCounter{
					maxTotal:                4,
					effectiveSuccessSamples: 4,
				},
				callMethods: []counterMethod{failure, failure, success},
			},
			&action{
				totalCount:   3,
				successCount: 1,
				failureCount: 2,
				result:       0,
			},
		),
		gen(
			"actual success rate for more than the effectiveSuccessSamples",
			[]string{cndNonZeroEffectiveSuccess},
			[]string{actCheckSuccessCount, actCheckFailureCount},
			&condition{
				counter: &countBasedFixedWindowCounter{
					maxTotal:                4,
					effectiveSuccessSamples: 4,
				},
				callMethods: []counterMethod{failure, failure, success, success},
			},
			&action{
				totalCount:   4,
				successCount: 2,
				failureCount: 2,
				result:       50,
			},
		),
		gen(
			"0 for less than the effectiveFailureSamples",
			[]string{cndNonZeroEffectiveFailure},
			[]string{actCheckSuccessCount, actCheckFailureCount},
			&condition{
				counter: &countBasedFixedWindowCounter{
					maxTotal:                4,
					effectiveFailureSamples: 4,
				},
				callMethods: []counterMethod{success, success, failure},
			},
			&action{
				totalCount:   3,
				successCount: 2,
				failureCount: 1,
				result:       0,
			},
		),
		gen(
			"actual failure rate for more than the effectiveFailureSamples",
			[]string{cndNonZeroEffectiveFailure},
			[]string{actCheckSuccessCount, actCheckFailureCount},
			&condition{
				counter: &countBasedFixedWindowCounter{
					maxTotal:                4,
					effectiveFailureSamples: 4,
				},
				callMethods: []counterMethod{success, success, failure, failure},
			},
			&action{
				totalCount:   4,
				successCount: 2,
				failureCount: 2,
				result:       50,
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			counter := tt.C().counter

			result := 0
			for _, m := range tt.C().callMethods {
				switch m {
				case success:
					result = counter.success()
				case failure:
					result = counter.failure()
				case reset:
					counter.reset()
					result = 0
				}
			}

			testutil.Diff(t, tt.A().totalCount, counter.totalCount)
			testutil.Diff(t, tt.A().successCount, counter.successCount)
			testutil.Diff(t, tt.A().failureCount, counter.failureCount)
			testutil.Diff(t, tt.A().result, result)
		})
	}
}

func TestTimeBasedFixedWindowCounter(t *testing.T) {
	type condition struct {
		counter     *timeBasedFixedWindowCounter
		callMethods []counterMethod
		timer       *testTimer
	}

	type action struct {
		successCount int
		failureCount int
		totalCount   int
		result       int
	}

	cndWindowSizeZero := "window size 0"
	cndNonZeroEffectiveSuccess := "non-zero effective success"
	cndNonZeroEffectiveFailure := "non-zero effective failure"
	cndSuccess := "success"
	cndFailure := "failure"
	cndReset := "reset"
	actCheckSuccessCount := "check success count"
	actCheckFailureCount := "check failure count"

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	tb.Condition(cndWindowSizeZero, "window size is zero")
	tb.Condition(cndNonZeroEffectiveSuccess, "effective success count is not zero")
	tb.Condition(cndNonZeroEffectiveFailure, "effective failure count is not zero")
	tb.Condition(cndSuccess, "success method was called")
	tb.Condition(cndFailure, "failure method was called")
	tb.Condition(cndReset, "reset method was called")
	tb.Action(actCheckSuccessCount, "check the final count of success")
	tb.Action(actCheckFailureCount, "check the final count of failure")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"window size is 0, failure for the last",
			[]string{cndWindowSizeZero, cndSuccess, cndFailure},
			[]string{actCheckSuccessCount, actCheckFailureCount},
			&condition{
				counter: &timeBasedFixedWindowCounter{
					lastResetTime: 0,
					windowWidth:   0,
				},
				callMethods: []counterMethod{success, success, failure, failure},
				timer:       &testTimer{times: []int64{10, 20, 30, 40, 50, 60, 70, 80, 90, 100}},
			},
			&action{
				totalCount:   1,
				successCount: 0,
				failureCount: 1,
				result:       100,
			},
		),
		gen(
			"window size is 0, success for the last",
			[]string{cndWindowSizeZero, cndSuccess, cndFailure},
			[]string{actCheckSuccessCount, actCheckFailureCount},
			&condition{
				counter: &timeBasedFixedWindowCounter{
					lastResetTime: 0,
					windowWidth:   0,
				},
				callMethods: []counterMethod{failure, failure, success, success},
				timer:       &testTimer{times: []int64{10, 20, 30, 40, 50, 60, 70, 80, 90, 100}},
			},
			&action{
				totalCount:   1,
				successCount: 1,
				failureCount: 0,
				result:       100,
			},
		),
		gen(
			"window size is 0, failure for the last",
			[]string{cndWindowSizeZero, cndSuccess, cndFailure},
			[]string{actCheckSuccessCount, actCheckFailureCount},
			&condition{
				counter: &timeBasedFixedWindowCounter{
					lastResetTime: 0,
					windowWidth:   100,
				},
				callMethods: []counterMethod{success, success, failure, failure},
				timer:       &testTimer{times: []int64{10, 20, 30, 40, 50, 60, 70, 80, 90, 100}},
			},
			&action{
				totalCount:   4,
				successCount: 2,
				failureCount: 2,
				result:       50,
			},
		),
		gen(
			"window size is 0, success for the last",
			[]string{cndWindowSizeZero, cndSuccess, cndFailure},
			[]string{actCheckSuccessCount, actCheckFailureCount},
			&condition{
				counter: &timeBasedFixedWindowCounter{
					lastResetTime: 0,
					windowWidth:   100,
				},
				callMethods: []counterMethod{failure, failure, success, success},
				timer:       &testTimer{times: []int64{10, 20, 30, 40, 50, 60, 70, 80, 90, 100}},
			},
			&action{
				totalCount:   4,
				successCount: 2,
				failureCount: 2,
				result:       50,
			},
		),
		gen(
			"call count is less than window size, failure for the last",
			[]string{cndSuccess, cndFailure},
			[]string{actCheckSuccessCount, actCheckFailureCount},
			&condition{
				counter: &timeBasedFixedWindowCounter{
					lastResetTime: 0,
					windowWidth:   20,
				},
				callMethods: []counterMethod{success, success, failure, failure},
				timer:       &testTimer{times: []int64{10, 20, 30, 40, 50, 60, 70, 80, 90, 100}},
			},
			&action{
				totalCount:   2,
				successCount: 0,
				failureCount: 2,
				result:       100,
			},
		),
		gen(
			"call count is less than window size, success for the last",
			[]string{cndSuccess, cndFailure},
			[]string{actCheckSuccessCount, actCheckFailureCount},
			&condition{
				counter: &timeBasedFixedWindowCounter{
					lastResetTime: 0,
					windowWidth:   20,
				},
				callMethods: []counterMethod{failure, failure, success, success},
				timer:       &testTimer{times: []int64{10, 20, 30, 40, 50, 60, 70, 80, 90, 100}},
			},
			&action{
				totalCount:   2,
				successCount: 2,
				failureCount: 0,
				result:       100,
			},
		),
		gen(
			"result is 0 for less than the effectiveSuccessSamples",
			[]string{cndNonZeroEffectiveSuccess},
			[]string{actCheckSuccessCount, actCheckFailureCount},
			&condition{
				counter: &timeBasedFixedWindowCounter{
					lastResetTime:           0,
					windowWidth:             100,
					effectiveSuccessSamples: 4,
				},
				callMethods: []counterMethod{failure, failure, success},
				timer:       &testTimer{times: []int64{10, 20, 30, 40, 50, 60, 70, 80, 90, 100}},
			},
			&action{
				totalCount:   3,
				successCount: 1,
				failureCount: 2,
				result:       0,
			},
		),
		gen(
			"actual success rate for more than the effectiveSuccessSamples",
			[]string{cndNonZeroEffectiveSuccess},
			[]string{actCheckSuccessCount, actCheckFailureCount},
			&condition{
				counter: &timeBasedFixedWindowCounter{
					lastResetTime:           0,
					windowWidth:             100,
					effectiveSuccessSamples: 4,
				},
				callMethods: []counterMethod{failure, failure, success, success},
				timer:       &testTimer{times: []int64{10, 20, 30, 40, 50, 60, 70, 80, 90, 100}},
			},
			&action{
				totalCount:   4,
				successCount: 2,
				failureCount: 2,
				result:       50,
			},
		),
		gen(
			"result is 0 for less than the effectiveFailureSamples",
			[]string{cndNonZeroEffectiveFailure},
			[]string{actCheckSuccessCount, actCheckFailureCount},
			&condition{
				counter: &timeBasedFixedWindowCounter{
					lastResetTime:           0,
					windowWidth:             100,
					effectiveFailureSamples: 4,
				},
				callMethods: []counterMethod{success, success, failure},
				timer:       &testTimer{times: []int64{10, 20, 30, 40, 50, 60, 70, 80, 90, 100}},
			},
			&action{
				totalCount:   3,
				successCount: 2,
				failureCount: 1,
				result:       0,
			},
		),
		gen(
			"actual failure rate for more than the effectiveFailureSamples",
			[]string{cndNonZeroEffectiveFailure},
			[]string{actCheckSuccessCount, actCheckFailureCount},
			&condition{
				counter: &timeBasedFixedWindowCounter{
					lastResetTime:           0,
					windowWidth:             100,
					effectiveFailureSamples: 4,
				},
				callMethods: []counterMethod{success, success, failure, failure},
				timer:       &testTimer{times: []int64{10, 20, 30, 40, 50, 60, 70, 80, 90, 100}},
			},
			&action{
				totalCount:   4,
				successCount: 2,
				failureCount: 2,
				result:       50,
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			// Replace timer for testing.
			tmp := unixMilli
			unixMilli = tt.C().timer.unixMilli
			defer func() {
				unixMilli = tmp
			}()

			counter := tt.C().counter

			result := 0
			for _, m := range tt.C().callMethods {
				switch m {
				case success:
					result = counter.success()
				case failure:
					result = counter.failure()
				case reset:
					counter.reset()
					result = 0
				}
			}

			testutil.Diff(t, tt.A().totalCount, counter.totalCount)
			testutil.Diff(t, tt.A().successCount, counter.successCount)
			testutil.Diff(t, tt.A().failureCount, counter.failureCount)
			testutil.Diff(t, tt.A().result, result)
		})
	}
}

func TestCountBasedSlidingWindowCounter(t *testing.T) {
	type condition struct {
		counter     *countBasedSlidingWindowCounter
		callMethods []counterMethod
	}

	type action struct {
		successCount int
		failureCount int
		totalCount   int
		result       int
	}

	cndWindowSizeZero := "window size 0"
	cndNonZeroEffectiveSuccess := "non-zero effective success"
	cndNonZeroEffectiveFailure := "non-zero effective failure"
	cndSuccess := "success"
	cndFailure := "failure"
	cndReset := "reset"
	actCheckSuccessCount := "check success count"
	actCheckFailureCount := "check failure count"

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	tb.Condition(cndWindowSizeZero, "window size is zero")
	tb.Condition(cndNonZeroEffectiveSuccess, "effective success count is not zero")
	tb.Condition(cndNonZeroEffectiveFailure, "effective failure count is not zero")
	tb.Condition(cndSuccess, "success method was called")
	tb.Condition(cndFailure, "failure method was called")
	tb.Condition(cndReset, "reset method was called")
	tb.Action(actCheckSuccessCount, "check the final count of success")
	tb.Action(actCheckFailureCount, "check the final count of failure")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"window size is 0, failure for the last",
			[]string{cndWindowSizeZero, cndSuccess, cndFailure},
			[]string{actCheckSuccessCount, actCheckFailureCount},
			&condition{
				counter: &countBasedSlidingWindowCounter{
					maxTotal:            0,
					totalHistory:        make([]int, 2),
					successCountHistory: make([]int, 2),
					failureCountHistory: make([]int, 2),
				},
				callMethods: []counterMethod{success, success, failure, failure},
			},
			&action{
				totalCount:   1,
				successCount: 0,
				failureCount: 1,
				result:       66, // The last 2/3 failed.
			},
		),
		gen(
			"window size is 0, success for the last",
			[]string{cndWindowSizeZero, cndSuccess, cndFailure},
			[]string{actCheckSuccessCount, actCheckFailureCount},
			&condition{
				counter: &countBasedSlidingWindowCounter{
					maxTotal:            0,
					totalHistory:        make([]int, 2),
					successCountHistory: make([]int, 2),
					failureCountHistory: make([]int, 2),
				},
				callMethods: []counterMethod{failure, failure, success, success},
			},
			&action{
				totalCount:   1,
				successCount: 1,
				failureCount: 0,
				result:       66, // The last 2/3 succeeded.
			},
		),
		gen(
			"call count is less than window size, failure for the last",
			[]string{cndSuccess, cndFailure},
			[]string{actCheckSuccessCount, actCheckFailureCount},
			&condition{
				counter: &countBasedSlidingWindowCounter{
					maxTotal:            4,
					totalHistory:        make([]int, 2),
					successCountHistory: make([]int, 2),
					failureCountHistory: make([]int, 2),
				},
				callMethods: []counterMethod{success, success, failure, failure},
			},
			&action{
				totalCount:   4,
				successCount: 2,
				failureCount: 2,
				result:       50,
			},
		),
		gen(
			"call count is less than window size, success for the last",
			[]string{cndSuccess, cndFailure},
			[]string{actCheckSuccessCount, actCheckFailureCount},
			&condition{
				counter: &countBasedSlidingWindowCounter{
					maxTotal:            4,
					totalHistory:        make([]int, 2),
					successCountHistory: make([]int, 2),
					failureCountHistory: make([]int, 2),
				},
				callMethods: []counterMethod{failure, failure, success, success},
			},
			&action{
				totalCount:   4,
				successCount: 2,
				failureCount: 2,
				result:       50,
			},
		),
		gen(
			"call count is more than window size, failure for the last",
			[]string{cndSuccess, cndFailure},
			[]string{actCheckSuccessCount, actCheckFailureCount},
			&condition{
				counter: &countBasedSlidingWindowCounter{
					maxTotal:            4,
					totalHistory:        make([]int, 2),
					successCountHistory: make([]int, 2),
					failureCountHistory: make([]int, 2),
				},
				callMethods: []counterMethod{success, success, failure, failure, success, failure},
			},
			&action{
				totalCount:   2,
				successCount: 1,
				failureCount: 1,
				result:       50,
			},
		),
		gen(
			"call count is less than window size, success for the last",
			[]string{cndSuccess, cndFailure},
			[]string{actCheckSuccessCount, actCheckFailureCount},
			&condition{
				counter: &countBasedSlidingWindowCounter{
					maxTotal:            4,
					totalHistory:        make([]int, 2),
					successCountHistory: make([]int, 2),
					failureCountHistory: make([]int, 2),
				},
				callMethods: []counterMethod{failure, failure, success, success},
			},
			&action{
				totalCount:   4,
				successCount: 2,
				failureCount: 2,
				result:       50,
			},
		),
		gen(
			"call count is more than window size, success for the last",
			[]string{cndSuccess, cndFailure},
			[]string{actCheckSuccessCount, actCheckFailureCount},
			&condition{
				counter: &countBasedSlidingWindowCounter{
					maxTotal:            4,
					totalHistory:        make([]int, 2),
					successCountHistory: make([]int, 2),
					failureCountHistory: make([]int, 2),
				},
				callMethods: []counterMethod{failure, failure, success, success, failure, success},
			},
			&action{
				totalCount:   2,
				successCount: 1,
				failureCount: 1,
				result:       50,
			},
		),
		gen(
			"reset after count",
			[]string{cndSuccess, cndFailure, cndReset},
			[]string{actCheckSuccessCount, actCheckFailureCount},
			&condition{
				counter: &countBasedSlidingWindowCounter{
					maxTotal:            4,
					totalHistory:        make([]int, 2),
					successCountHistory: make([]int, 2),
					failureCountHistory: make([]int, 2),
				},
				callMethods: []counterMethod{failure, failure, success, success, reset},
			},
			&action{
				totalCount:   0,
				successCount: 0,
				failureCount: 0,
				result:       0,
			},
		),
		gen(
			"result is 0 for less than the effectiveSuccessSamples",
			[]string{cndNonZeroEffectiveSuccess},
			[]string{actCheckSuccessCount, actCheckFailureCount},
			&condition{
				counter: &countBasedSlidingWindowCounter{
					maxTotal:                4,
					effectiveSuccessSamples: 4,
					totalHistory:            make([]int, 2),
					successCountHistory:     make([]int, 2),
					failureCountHistory:     make([]int, 2),
				},
				callMethods: []counterMethod{failure, failure, success},
			},
			&action{
				totalCount:   3,
				successCount: 1,
				failureCount: 2,
				result:       0,
			},
		),
		gen(
			"actual success rate for more than the effectiveSuccessSamples",
			[]string{cndNonZeroEffectiveSuccess},
			[]string{actCheckSuccessCount, actCheckFailureCount},
			&condition{
				counter: &countBasedSlidingWindowCounter{
					maxTotal:                4,
					effectiveSuccessSamples: 4,
					totalHistory:            make([]int, 2),
					successCountHistory:     make([]int, 2),
					failureCountHistory:     make([]int, 2),
				},
				callMethods: []counterMethod{failure, failure, success, success},
			},
			&action{
				totalCount:   4,
				successCount: 2,
				failureCount: 2,
				result:       50,
			},
		),
		gen(
			"result is 0 for less than the effectiveFailureSamples",
			[]string{cndNonZeroEffectiveFailure},
			[]string{actCheckSuccessCount, actCheckFailureCount},
			&condition{
				counter: &countBasedSlidingWindowCounter{
					maxTotal:                4,
					effectiveFailureSamples: 4,
					totalHistory:            make([]int, 2),
					successCountHistory:     make([]int, 2),
					failureCountHistory:     make([]int, 2),
				},
				callMethods: []counterMethod{success, success, failure},
			},
			&action{
				totalCount:   3,
				successCount: 2,
				failureCount: 1,
				result:       0,
			},
		),
		gen(
			"actual failure rate for more than the effectiveFailureSamples",
			[]string{cndNonZeroEffectiveFailure},
			[]string{actCheckSuccessCount, actCheckFailureCount},
			&condition{
				counter: &countBasedSlidingWindowCounter{
					maxTotal:                4,
					effectiveFailureSamples: 4,
					totalHistory:            make([]int, 2),
					successCountHistory:     make([]int, 2),
					failureCountHistory:     make([]int, 2),
				},
				callMethods: []counterMethod{success, success, failure, failure},
			},
			&action{
				totalCount:   4,
				successCount: 2,
				failureCount: 2,
				result:       50,
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			counter := tt.C().counter

			result := 0
			for _, m := range tt.C().callMethods {
				switch m {
				case success:
					result = counter.success()
				case failure:
					result = counter.failure()
				case reset:
					counter.reset()
					result = 0
				}
			}

			testutil.Diff(t, tt.A().totalCount, counter.totalCount)
			testutil.Diff(t, tt.A().successCount, counter.successCount)
			testutil.Diff(t, tt.A().failureCount, counter.failureCount)
			testutil.Diff(t, tt.A().result, result)
		})
	}
}

func TestTimeBasedSlidingWindowCounter(t *testing.T) {
	type condition struct {
		counter     *timeBasedSlidingWindowCounter
		callMethods []counterMethod
		timer       *testTimer
	}

	type action struct {
		successCount int
		failureCount int
		totalCount   int
		result       int
	}

	cndWindowSizeZero := "window size 0"
	cndNonZeroEffectiveSuccess := "non-zero effective success"
	cndNonZeroEffectiveFailure := "non-zero effective failure"
	cndSuccess := "success"
	cndFailure := "failure"
	cndReset := "reset"
	actCheckSuccessCount := "check success count"
	actCheckFailureCount := "check failure count"

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	tb.Condition(cndWindowSizeZero, "window size is zero")
	tb.Condition(cndNonZeroEffectiveSuccess, "effective success count is not zero")
	tb.Condition(cndNonZeroEffectiveFailure, "effective failure count is not zero")
	tb.Condition(cndSuccess, "success method was called")
	tb.Condition(cndFailure, "failure method was called")
	tb.Condition(cndReset, "reset method was called")
	tb.Action(actCheckSuccessCount, "check the final count of success")
	tb.Action(actCheckFailureCount, "check the final count of failure")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"window size is 0, failure for the last",
			[]string{cndWindowSizeZero, cndSuccess, cndFailure},
			[]string{actCheckSuccessCount, actCheckFailureCount},
			&condition{
				counter: &timeBasedSlidingWindowCounter{
					lastResetTime:       0,
					windowWidth:         0,
					totalHistory:        make([]int, 2),
					successCountHistory: make([]int, 2),
					failureCountHistory: make([]int, 2),
				},
				callMethods: []counterMethod{success, success, failure, failure},
				timer:       &testTimer{times: []int64{10, 20, 30, 40, 50, 60, 70, 80, 90, 100}},
			},
			&action{
				totalCount:   1,
				successCount: 0,
				failureCount: 1,
				result:       66, // The last 2/3 failed.
			},
		),
		gen(
			"window size is 0, success for the last",
			[]string{cndWindowSizeZero, cndSuccess, cndFailure},
			[]string{actCheckSuccessCount, actCheckFailureCount},
			&condition{
				counter: &timeBasedSlidingWindowCounter{
					lastResetTime:       0,
					windowWidth:         0,
					totalHistory:        make([]int, 2),
					successCountHistory: make([]int, 2),
					failureCountHistory: make([]int, 2),
				},
				callMethods: []counterMethod{failure, failure, success, success},
				timer:       &testTimer{times: []int64{10, 20, 30, 40, 50, 60, 70, 80, 90, 100}},
			},
			&action{
				totalCount:   1,
				successCount: 1,
				failureCount: 0,
				result:       66, // The last 2/3 succeeded.
			},
		),
		gen(
			"all results within the window, failure for the last",
			[]string{cndSuccess, cndFailure},
			[]string{actCheckSuccessCount, actCheckFailureCount},
			&condition{
				counter: &timeBasedSlidingWindowCounter{
					lastResetTime:       0,
					windowWidth:         40,
					totalHistory:        make([]int, 2),
					successCountHistory: make([]int, 2),
					failureCountHistory: make([]int, 2),
				},
				callMethods: []counterMethod{success, success, failure, failure},
				timer:       &testTimer{times: []int64{10, 20, 30, 40, 50, 60, 70, 80, 90, 100}},
			},
			&action{
				totalCount:   4,
				successCount: 2,
				failureCount: 2,
				result:       50,
			},
		),
		gen(
			"all results within the window, success for the last",
			[]string{cndSuccess, cndFailure},
			[]string{actCheckSuccessCount, actCheckFailureCount},
			&condition{
				counter: &timeBasedSlidingWindowCounter{
					lastResetTime:       0,
					windowWidth:         40,
					totalHistory:        make([]int, 2),
					successCountHistory: make([]int, 2),
					failureCountHistory: make([]int, 2),
				},
				callMethods: []counterMethod{failure, failure, success, success},
				timer:       &testTimer{times: []int64{10, 20, 30, 40, 50, 60, 70, 80, 90, 100}},
			},
			&action{
				totalCount:   4,
				successCount: 2,
				failureCount: 2,
				result:       50,
			},
		),
		gen(
			"more than the window size, failure for the last",
			[]string{cndSuccess, cndFailure},
			[]string{actCheckSuccessCount, actCheckFailureCount},
			&condition{
				counter: &timeBasedSlidingWindowCounter{
					lastResetTime:       0,
					windowWidth:         20,
					totalHistory:        make([]int, 2),
					successCountHistory: make([]int, 2),
					failureCountHistory: make([]int, 2),
				},
				callMethods: []counterMethod{success, success, failure, failure},
				timer:       &testTimer{times: []int64{10, 20, 30, 40, 50, 60, 70, 80, 90, 100}},
			},
			&action{
				totalCount:   2,
				successCount: 0,
				failureCount: 2,
				result:       50,
			},
		),
		gen(
			"more than the window size,  success for the last",
			[]string{cndSuccess, cndFailure},
			[]string{actCheckSuccessCount, actCheckFailureCount},
			&condition{
				counter: &timeBasedSlidingWindowCounter{
					lastResetTime:       0,
					windowWidth:         20,
					totalHistory:        make([]int, 2),
					successCountHistory: make([]int, 2),
					failureCountHistory: make([]int, 2),
				},
				callMethods: []counterMethod{failure, failure, success, success},
				timer:       &testTimer{times: []int64{10, 20, 30, 40, 50, 60, 70, 80, 90, 100}},
			},
			&action{
				totalCount:   2,
				successCount: 2,
				failureCount: 0,
				result:       50,
			},
		),
		gen(
			"result is 0 for less than the effectiveSuccessSamples",
			[]string{cndNonZeroEffectiveSuccess},
			[]string{actCheckSuccessCount, actCheckFailureCount},
			&condition{
				counter: &timeBasedSlidingWindowCounter{
					lastResetTime:           0,
					windowWidth:             40,
					effectiveSuccessSamples: 4,
					totalHistory:            make([]int, 2),
					successCountHistory:     make([]int, 2),
					failureCountHistory:     make([]int, 2),
				},
				callMethods: []counterMethod{failure, failure, success},
				timer:       &testTimer{times: []int64{10, 20, 30, 40, 50, 60, 70, 80, 90, 100}},
			},
			&action{
				totalCount:   3,
				successCount: 1,
				failureCount: 2,
				result:       0,
			},
		),
		gen(
			"actual success rate for more than the effectiveSuccessSamples",
			[]string{cndNonZeroEffectiveSuccess},
			[]string{actCheckSuccessCount, actCheckFailureCount},
			&condition{
				counter: &timeBasedSlidingWindowCounter{
					lastResetTime:           0,
					windowWidth:             40,
					effectiveSuccessSamples: 4,
					totalHistory:            make([]int, 2),
					successCountHistory:     make([]int, 2),
					failureCountHistory:     make([]int, 2),
				},
				callMethods: []counterMethod{failure, failure, success, success},
				timer:       &testTimer{times: []int64{10, 20, 30, 40, 50, 60, 70, 80, 90, 100}},
			},
			&action{
				totalCount:   4,
				successCount: 2,
				failureCount: 2,
				result:       50,
			},
		),
		gen(
			"result is 0 for less than the effectiveFailureSamples",
			[]string{cndNonZeroEffectiveFailure},
			[]string{actCheckSuccessCount, actCheckFailureCount},
			&condition{
				counter: &timeBasedSlidingWindowCounter{
					lastResetTime:           0,
					windowWidth:             40,
					effectiveFailureSamples: 4,
					totalHistory:            make([]int, 2),
					successCountHistory:     make([]int, 2),
					failureCountHistory:     make([]int, 2),
				},
				callMethods: []counterMethod{success, success, failure},
				timer:       &testTimer{times: []int64{10, 20, 30, 40, 50, 60, 70, 80, 90, 100}},
			},
			&action{
				totalCount:   3,
				successCount: 2,
				failureCount: 1,
				result:       0,
			},
		),
		gen(
			"actual failure rate for more than the effectiveFailureSamples",
			[]string{cndNonZeroEffectiveFailure},
			[]string{actCheckSuccessCount, actCheckFailureCount},
			&condition{
				counter: &timeBasedSlidingWindowCounter{
					lastResetTime:           0,
					windowWidth:             40,
					effectiveFailureSamples: 4,
					totalHistory:            make([]int, 2),
					successCountHistory:     make([]int, 2),
					failureCountHistory:     make([]int, 2),
				},
				callMethods: []counterMethod{success, success, failure, failure},
				timer:       &testTimer{times: []int64{10, 20, 30, 40, 50, 60, 70, 80, 90, 100}},
			},
			&action{
				totalCount:   4,
				successCount: 2,
				failureCount: 2,
				result:       50,
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			// Replace timer for testing.
			tmp := unixMilli
			unixMilli = tt.C().timer.unixMilli
			defer func() {
				unixMilli = tmp
			}()

			counter := tt.C().counter

			result := 0
			for _, m := range tt.C().callMethods {
				switch m {
				case success:
					result = counter.success()
				case failure:
					result = counter.failure()
				case reset:
					counter.reset()
					result = 0
				}
			}

			testutil.Diff(t, tt.A().totalCount, counter.totalCount)
			testutil.Diff(t, tt.A().successCount, counter.successCount)
			testutil.Diff(t, tt.A().failureCount, counter.failureCount)
			testutil.Diff(t, tt.A().result, result)
			fmt.Printf("%#v\n", counter)
		})
	}
}
