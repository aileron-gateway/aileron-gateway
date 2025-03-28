package httpproxy

import (
	"sync"
	"time"

	v1 "github.com/aileron-gateway/aileron-gateway/apis/core/v1"
)

var (
	// unixMilli returns the current unix milliseconds.
	// This function will be replaced when testing.
	unixMilli func() int64 = func() int64 {
		return time.Now().UnixMilli()
	}
)

// newCircuitBreaker returns a new circuit breaker.
// This function returns nil if a nil spec was given by the argument.
func newCircuitBreaker(spec *v1.CircuitBreaker) circuitBreaker {
	if spec == nil {
		return nil
	}

	cb := &circuitBreakerController{
		status:           closed,
		successThreshold: int(spec.SuccessThreshold),
		failureThreshold: int(spec.FailureThreshold),
		wait:             time.Second * time.Duration(spec.WaitDuration),
	}

	switch v := spec.CircuitBreakerCounter.(type) {
	case *v1.CircuitBreaker_ConsecutiveCounter:
		cb.counter = &consecutiveCounter{}
	case *v1.CircuitBreaker_CountBasedFixedWindowCounter:
		s := v.CountBasedFixedWindowCounter
		cb.counter = &countBasedFixedWindowCounter{
			maxTotal:                int(s.Samples),
			effectiveFailureSamples: int(spec.EffectiveFailureSamples),
			effectiveSuccessSamples: int(spec.EffectiveSuccessSamples),
		}
	case *v1.CircuitBreaker_TimeBasedFixedWindowCounter:
		s := v.TimeBasedFixedWindowCounter
		cb.counter = &timeBasedFixedWindowCounter{
			lastResetTime:           unixMilli(),
			windowWidth:             int64(1000 * s.WindowWidth), // Convert second to millisecond.
			effectiveFailureSamples: int(spec.EffectiveFailureSamples),
			effectiveSuccessSamples: int(spec.EffectiveSuccessSamples),
		}
	case *v1.CircuitBreaker_CountBasedSlidingWindowCounter:
		s := v.CountBasedSlidingWindowCounter
		cb.counter = &countBasedSlidingWindowCounter{
			maxTotal:                int(s.Samples),
			effectiveFailureSamples: int(spec.EffectiveFailureSamples),
			effectiveSuccessSamples: int(spec.EffectiveSuccessSamples),
			totalHistory:            make([]int, s.HistoryLimit),
			successCountHistory:     make([]int, s.HistoryLimit),
			failureCountHistory:     make([]int, s.HistoryLimit),
		}
	case *v1.CircuitBreaker_TimeBasedSlidingWindowCounter:
		s := v.TimeBasedSlidingWindowCounter
		cb.counter = &timeBasedSlidingWindowCounter{
			lastResetTime:           unixMilli(),
			windowWidth:             int64(1000 * s.WindowWidth), // Convert second to millisecond.
			effectiveFailureSamples: int(spec.EffectiveFailureSamples),
			effectiveSuccessSamples: int(spec.EffectiveSuccessSamples),
			totalHistory:            make([]int, s.HistoryLimit),
			successCountHistory:     make([]int, s.HistoryLimit),
			failureCountHistory:     make([]int, s.HistoryLimit),
		}
	}

	return cb
}

// circuitBreakerStatus is the type of status of circuit breaker.
// Following three statuses are defined for circuit breaker.
//   - closed
//   - opened
//   - halfOpened
type circuitBreakerStatus int

const (
	closed     circuitBreakerStatus = iota // Circuit connectivity is ready.
	opened                                 // Circuit connectivity is non-ready.
	halfOpened                             // Circuit connectivity is ready but close to non-ready.
)

// circuitBreaker is the interface of circuit breaker.
// Circuit breaker monitors target servers.
type circuitBreaker interface {
	// active returns true when this circuit breaker is closed or half-opened.
	// It's implementers responsible to make this method safe
	// for concurrency call.
	Active() bool

	// countSuccess increment count of successes.
	// It's implementers responsible to make this method safe
	// for concurrency call.
	countSuccess()

	// countFailure increment count of failures.
	// It's implementers responsible to make this method safe
	// for concurrency call.
	countFailure()
}

// circuitBreakerController controls status of circuit breaker.
// This implements proxy.circuitBreaker interface.
type circuitBreakerController struct {

	// mu protects status and counter.
	mu sync.RWMutex

	// status is the status of this circuit breaker.
	// This status must be changed by calling "changeStatus" method.
	// Do not change this field by directory accessing this field.
	status circuitBreakerStatus

	// counter is a success and failure counter
	// which holds calculates failure rate and success rate.
	// Methods are not safe for concurrent call.
	counter circuitBreakerCounter

	successThreshold int
	failureThreshold int

	// wait is the duration to wait until
	wait time.Duration
}

// Active returns the status, or availability, of this circuit.
// This method is safe for concurrent call.
func (cb *circuitBreakerController) Active() bool {
	cb.mu.RLock()
	defer cb.mu.RUnlock()
	// circuit breaker is available when not opened.
	return cb.status != opened
}

// changeStatus changes the status of this circuit breaker.
// This method do nothing if the given status is the same as current.
// this method is not safe for concurrent call.
func (cb *circuitBreakerController) changeStatus(s circuitBreakerStatus) {
	if cb.status == s {
		return // No need to be changed.
	}

	cb.status = s
	cb.counter.reset()

	if s == opened {
		// Now status is changed to opened and this circuit breaker is not active.
		// Change from opened to halfOpened after wait duration.
		go func() {
			time.Sleep(cb.wait)
			cb.mu.Lock()
			defer cb.mu.Unlock()
			cb.status = halfOpened // Change status opened -> halfOpened.
		}()
	}
}

// countSuccess counts success.
// It depends on the circuit breaker status if the internal
// counter will be actually incremented or not.
// This method is safe for concurrent call.
func (cb *circuitBreakerController) countSuccess() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	if cb.status != halfOpened {
		return // Do nothing for closed and opened status.
	}

	// success() have to be called only when the status is halfOpened.
	successRate := cb.counter.success()
	if successRate >= cb.successThreshold {
		cb.changeStatus(closed)
	}
}

// countFailure counts failure.
// It depends on the circuit breaker status if the internal
// counter will be actually incremented or not.
// This method is safe for concurrent call.
func (cb *circuitBreakerController) countFailure() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	switch cb.status {
	case closed:
		// failure() have to be called only when the status is closed.
		failureRate := cb.counter.failure()
		if failureRate >= cb.failureThreshold {
			cb.changeStatus(opened) // closed -> opened
		}
	case halfOpened:
		cb.changeStatus(opened) // halfOpened -> opened
	}

	// No operation needed for opened status.
}

// circuitBreakerCounter is the interface of success and failure counter
// for circuit breaker.
type circuitBreakerCounter interface {
	// success increments success counter and return success rate in percentage.
	// success should be called when the current status of circuit breaker is halfOpened.
	// This method is not safe for concurrent call.
	// It is caller's responsibility to protect calling this method for concurrent call.
	success() int
	// failure increments failure counter and return failure rate in percentage.
	// failure should be called when the current status of circuit breaker is closed.
	// This method is not safe for concurrent call.
	// It is caller's responsibility to protect calling this method for concurrent call.
	failure() int
	// reset resets the success and failure counters.
	// failure should be called when the circuit breaker status was changed.
	// This method is not safe for concurrent call.
	// It is caller's responsibility to protect calling this method for concurrent call.
	reset()
}

// consecutiveCounter is the circuit breaker counter
// with consecutive counting algorithm.
// This implements proxy.circuitBreakerCounter interface.
type consecutiveCounter struct {
	// successCount is the consecutive successful count.
	// The upper bound of successCount is successThreshold of circuitBreakerController.
	successCount int
	// failureCount is the consecutive failure count.
	// The upper bound of failureCount is successThreshold of circuitBreakerController.
	failureCount int
}

func (c *consecutiveCounter) success() int {
	c.successCount += 1
	c.failureCount = 0
	return c.successCount
}

func (c *consecutiveCounter) failure() int {
	c.successCount = 0
	c.failureCount += 1
	return c.failureCount
}

func (c *consecutiveCounter) reset() {
	c.successCount = 0
	c.failureCount = 0
}

// countBasedFixedWindowCounter is the circuit breaker counter
// with count-based fixed window algorithm.
// This implements proxy.circuitBreakerCounter interface.
type countBasedFixedWindowCounter struct {

	// maxTotal is the window size.
	// This is the maximum count of total success and failures.
	maxTotal int

	effectiveSuccessSamples int
	effectiveFailureSamples int

	// totalCount is the sum of successful and failure results in the current window.
	// totalCount = successCount + failureCount <= maxTotal.
	totalCount int
	// successCount is the successful count in the current window.
	// The upper bound of successCount is maxTotal.
	// For more specific, successCount <= maxTotal * successThreshold/100.
	successCount int
	// failureCount is the failure count in the current window.
	// The upper bound of failureCount is maxTotal.
	// For more specific, failureCount <= maxTotal * failureThreshold/100.
	failureCount int
}

func (c *countBasedFixedWindowCounter) success() int {
	if c.totalCount >= c.maxTotal {
		c.reset()
	}
	c.totalCount += 1
	c.successCount += 1
	if c.totalCount < c.effectiveSuccessSamples {
		// Success rate should be 0 until the
		// number of samples exceeds the effectiveSuccessSamples.
		return 0
	}
	// The result of percentage is rounded down.
	return (100 * c.successCount) / c.totalCount
}

func (c *countBasedFixedWindowCounter) failure() int {
	if c.totalCount >= c.maxTotal {
		c.reset()
	}
	c.totalCount += 1
	c.failureCount += 1
	if c.totalCount < c.effectiveFailureSamples {
		// Failure rate should be 0 until the
		// number of samples exceeds the effectiveSuccessSamples.
		return 0
	}
	// The result of percentage is rounded down.
	return (100 * c.failureCount) / c.totalCount
}

func (c *countBasedFixedWindowCounter) reset() {
	c.totalCount = 0
	c.successCount = 0
	c.failureCount = 0
}

// timeBasedFixedWindowCounter is the circuit breaker counter
// with time-based fixed window algorithm.
// This implements proxy.circuitBreakerCounter interface.
type timeBasedFixedWindowCounter struct {
	// lastResetTime is the unix milliseconds of the time when
	// the reset() method was called last time.
	lastResetTime int64
	// windowWidth is the window interval to reset counter
	// in milliseconds.
	windowWidth int64

	effectiveSuccessSamples int
	effectiveFailureSamples int

	// totalCount is the sum of successful and failure in the current window.
	// The windowWidth must be selected as the totalCount not to
	// exceed the math.MaxInt.
	// totalCount = successCount + failureCount <= math.MaxInt.
	totalCount int
	// successCount is the successful count in the current window.
	// The windowWidth must be selected as the successCount not to
	// exceed the math.MaxInt.
	successCount int
	// failureCount is the failure count in the current window.
	// The windowWidth must be selected as the failureCount not to
	// exceed the math.MaxInt.
	failureCount int
}

func (c *timeBasedFixedWindowCounter) success() int {
	if unixMilli()-c.lastResetTime > c.windowWidth {
		c.reset()
	}
	c.totalCount += 1
	c.successCount += 1
	if c.totalCount < c.effectiveSuccessSamples {
		// Success rate should be 0 until the
		// number of samples exceeds the effectiveSuccessSamples.
		return 0
	}
	// The result of percentage is rounded down.
	return (100 * c.successCount) / c.totalCount
}

func (c *timeBasedFixedWindowCounter) failure() int {
	if unixMilli()-c.lastResetTime > c.windowWidth {
		c.reset()
	}
	c.totalCount += 1
	c.failureCount += 1
	if c.totalCount < c.effectiveFailureSamples {
		// Failure rate should be 0 until the
		// number of samples exceeds the effectiveSuccessSamples.
		return 0
	}
	// The result of percentage is rounded down.
	return (100 * c.failureCount) / c.totalCount
}

func (c *timeBasedFixedWindowCounter) reset() {
	c.totalCount = 0
	c.successCount = 0
	c.failureCount = 0
	c.lastResetTime = unixMilli()
}

// countBasedSlidingWindowCounter is the circuit breaker counter
// with count-based sliding window algorithm.
// This implements proxy.circuitBreakerCounter interface.
type countBasedSlidingWindowCounter struct {
	maxTotal int

	effectiveSuccessSamples int
	effectiveFailureSamples int

	// totalCount is the sum of successful and failure results in the current window.
	// totalCount = successCount + failureCount <= maxTotal.
	totalCount int
	// successCount is the successful count in the current window.
	// The upper bound of successCount is maxTotal.
	// For more specific, successCount <= maxTotal * successThreshold/100.
	successCount int
	// failureCount is the failure count in the current window.
	// The upper bound of failureCount is maxTotal.
	// For more specific, failureCount <= maxTotal * failureThreshold/100.
	failureCount int

	cumulativeTotal        int
	cumulativeSuccessCount int
	cumulativeFailureCount int

	oldestPosition      int
	totalHistory        []int
	successCountHistory []int
	failureCountHistory []int
}

func (c *countBasedSlidingWindowCounter) success() int {
	if c.totalCount >= c.maxTotal {
		c.reset()
	}
	c.totalCount += 1
	c.successCount += 1
	c.cumulativeTotal += 1
	c.cumulativeSuccessCount += 1
	if c.totalCount < c.effectiveSuccessSamples {
		// Success rate should be 0 until the
		// number of samples exceeds the effectiveSuccessSamples.
		return 0
	}
	// The result of percentage is rounded down.
	return (100 * c.cumulativeSuccessCount) / c.cumulativeTotal
}

func (c *countBasedSlidingWindowCounter) failure() int {
	if c.totalCount >= c.maxTotal {
		c.reset()
	}
	c.totalCount += 1
	c.failureCount += 1
	c.cumulativeTotal += 1
	c.cumulativeFailureCount += 1
	if c.totalCount < c.effectiveFailureSamples {
		// Failure rate should be 0 until the
		// number of samples exceeds the effectiveSuccessSamples.
		return 0
	}
	// The result of percentage is rounded down.
	return (100 * c.cumulativeFailureCount) / c.cumulativeTotal
}

func (c *countBasedSlidingWindowCounter) reset() {
	// Subtract the oldest values from the cumulative values.
	c.cumulativeTotal -= c.totalHistory[c.oldestPosition]
	c.cumulativeSuccessCount -= c.successCountHistory[c.oldestPosition]
	c.cumulativeFailureCount -= c.failureCountHistory[c.oldestPosition]

	// Overwrite the oldest values with the current values.
	c.totalHistory[c.oldestPosition] = c.totalCount
	c.successCountHistory[c.oldestPosition] = c.successCount
	c.failureCountHistory[c.oldestPosition] = c.failureCount

	// Reset current values.
	c.totalCount = 0
	c.successCount = 0
	c.failureCount = 0

	c.oldestPosition += 1
	if c.oldestPosition >= len(c.totalHistory) {
		c.oldestPosition = 0
	}
}

// timeBasedSlidingWindowCounter is the circuit breaker counter
// with time-based sliding window algorithm.
// This implements proxy.circuitBreakerCounter interface.
type timeBasedSlidingWindowCounter struct {
	// lastResetTime is the unix milliseconds of the time when
	// the reset() method was called last time.
	lastResetTime int64
	// windowWidth is the window interval to reset counter
	// in milliseconds.
	windowWidth int64

	effectiveSuccessSamples int
	effectiveFailureSamples int

	// totalCount is the sum of successful and failure in the current window.
	// The windowWidth must be selected as the totalCount not to
	// exceed the math.MaxInt.
	// totalCount = successCount + failureCount <= math.MaxInt.
	totalCount int
	// successCount is the successful count in the current window.
	// The windowWidth must be selected as the successCount not to
	// exceed the math.MaxInt.
	successCount int
	// failureCount is the failure count in the current window.
	// The windowWidth must be selected as the failureCount not to
	// exceed the math.MaxInt.
	failureCount int

	cumulativeTotal        int
	cumulativeSuccessCount int
	cumulativeFailureCount int

	oldestPosition      int
	totalHistory        []int
	successCountHistory []int
	failureCountHistory []int
}

func (c *timeBasedSlidingWindowCounter) success() int {
	if unixMilli()-c.lastResetTime > c.windowWidth {
		c.reset()
	}
	c.totalCount += 1
	c.successCount += 1
	c.cumulativeTotal += 1
	c.cumulativeSuccessCount += 1
	if c.totalCount < c.effectiveSuccessSamples {
		// Success rate should be 0 until the
		// number of samples exceeds the effectiveSuccessSamples.
		return 0
	}
	// The result of percentage is rounded down.
	return (100 * c.cumulativeSuccessCount) / c.cumulativeTotal
}

func (c *timeBasedSlidingWindowCounter) failure() int {
	if unixMilli()-c.lastResetTime > c.windowWidth {
		c.reset()
	}
	c.totalCount += 1
	c.failureCount += 1
	c.cumulativeTotal += 1
	c.cumulativeFailureCount += 1
	if c.totalCount < c.effectiveFailureSamples {
		// Failure rate should be 0 until the
		// number of samples exceeds the effectiveSuccessSamples.
		return 0
	}
	// The result of percentage is rounded down.
	return (100 * c.cumulativeFailureCount) / c.cumulativeTotal
}

func (c *timeBasedSlidingWindowCounter) reset() {
	// Subtract the oldest values from the cumulative values.
	c.cumulativeTotal -= c.totalHistory[c.oldestPosition]
	c.cumulativeSuccessCount -= c.successCountHistory[c.oldestPosition]
	c.cumulativeFailureCount -= c.failureCountHistory[c.oldestPosition]

	// Overwrite the oldest values with the current values.
	c.totalHistory[c.oldestPosition] = c.totalCount
	c.successCountHistory[c.oldestPosition] = c.successCount
	c.failureCountHistory[c.oldestPosition] = c.failureCount

	// Reset current values.
	c.totalCount = 0
	c.successCount = 0
	c.failureCount = 0

	c.oldestPosition += 1
	if c.oldestPosition >= len(c.totalHistory) {
		c.oldestPosition = 0
	}

	c.lastResetTime = unixMilli()
}
