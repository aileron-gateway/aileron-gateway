// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package throttle

import (
	"context"
	"sync"
	"time"

	"github.com/aileron-gateway/aileron-gateway/util/resilience"
)

// throttler apply rate limiting to the requests.
type throttler interface {
	// accept returns the result if the requests
	// should be accepted or not.
	// ok is true when requests are acceptable.
	// ok is false when requests should be denied.
	// releaser must be called after responses are sent to client.
	accept(context.Context) (ok bool, releaser func())
}

// releaser call release function.
// This releaser is intended to be used
// by maxConnection throttler.
type releaser struct {
	once        sync.Once
	releaseFunc func()
}

func (r *releaser) release() {
	r.once.Do(r.releaseFunc)
}

// noopReleaser do nothing.
// This releaser is intended to be used by
// throttlers except for maxConnectionThrottler.
var noopReleaser = func() {}

// retryThrottler wraps throttlers
// and apply retrying to requests.
type retryThrottler struct {
	throttler
	// maxRetry is the maximum retry count.
	// Requests will be rejected when all retries are failed for maxRetry times.
	maxRetry int
	// waiter determines wait duration until next retry.
	waiter resilience.Waiter
}

func (t *retryThrottler) accept(ctx context.Context) (bool, func()) {
	if ok, f := t.throttler.accept(ctx); ok {
		return true, f
	}

	ticker := time.NewTicker(time.Millisecond)
	defer ticker.Stop()

loop:
	for i := 1; i <= t.maxRetry; i++ {
		ticker.Reset(t.waiter.Wait(i))
		select {
		case <-ticker.C:
		case <-ctx.Done():
			break loop
		}

		// Check if the request is now acceptable or not.
		if ok, f := t.throttler.accept(ctx); ok {
			return true, f
		}
	}

	// Retry failed over maxRetry counts.
	// So, return false and decline the request.
	return false, noopReleaser
}

// maxConnections is a throttler with max connections algorithm.
// This implements throttler interface.
type maxConnections struct {
	// sem is the semaphore variable.
	// The capacity of the sem is the same as the number of max connections.
	sem chan struct{}
}

func (t *maxConnections) accept(ctx context.Context) (bool, func()) {
	select {
	case t.sem <- struct{}{}:
		return true, (&releaser{releaseFunc: func() { <-t.sem }}).release // Accept
	default:
		// Request not acceptable.
		return false, nil // Decline.
	}
}

// fixedWindow is a throttler with fixed window algorithm.
// This implements throttler interface.
type fixedWindow struct {
	// bucket holds tokens.
	// 1 token will be consumed for accepting 1 request.
	// This bucket has size N.
	bucket chan struct{}
	// window is the time duration that tokens
	// should be refilled in the bucket.
	// The bucket will be fulfilled every window intervals.
	window time.Duration
}

func (t *fixedWindow) fill() {
	ticker := time.NewTicker(t.window)
	defer ticker.Stop()
	for {
		<-ticker.C
		n := cap(t.bucket) - len(t.bucket)
		for i := 0; i < n; i++ {
			t.bucket <- struct{}{}
		}
	}
}

func (t *fixedWindow) accept(ctx context.Context) (bool, func()) {
	select {
	case <-t.bucket:
		return true, noopReleaser // Accept.
	default:
		return false, nil // Decline.
	}
}

// tokenBucket is a throttler with token bucket algorithm.
// This implements throttler interface.
type tokenBucket struct {
	// bucket holds tokens.
	// 1 token will be consumed for accepting 1 request.
	// This bucket has size N.
	bucket chan struct{}
	// rate is the number of tokens refilled in
	// the bucket every "interval" time.
	// This number should be less than or equal to the bucket size N.
	rate int
	// interval is the time duration that tokens
	// should be refilled in the bucket.
	interval time.Duration
}

func (t *tokenBucket) fill() {
	ticker := time.NewTicker(t.interval)
	defer ticker.Stop()
	for {
		<-ticker.C
		n := cap(t.bucket) - len(t.bucket)
		n = minInt(n, t.rate)
		for i := 0; i < n; i++ {
			t.bucket <- struct{}{}
		}
	}
}

func (t *tokenBucket) accept(ctx context.Context) (bool, func()) {
	select {
	case <-t.bucket:
		return true, noopReleaser // Accept.
	default:
		return false, nil // Decline.
	}
}

func minInt(x, y int) int {
	if x > y {
		return y
	}
	return x
}

// leakyBucket is a throttler with leaky bucket algorithm.
// This implements throttler interface.
type leakyBucket struct {
	// bucket holds tokens.
	// 1 token will be consumed for accepting 1 request.
	// This bucket has size N.
	bucket chan chan struct{}
	// rate is the number of tokens refilled in
	// the bucket every "interval" time.
	// This number is the same as the bucket size N.
	rate int
	// interval is the time duration that tokens
	// should be refilled in the bucket.
	interval time.Duration
}

func (t *leakyBucket) leak() {
	ticker := time.NewTicker(t.interval)
	defer ticker.Stop()
	for {
		<-ticker.C
		for i := 0; i < t.rate; i++ {
			select {
			case c := <-t.bucket:
				close(c) // Accept.
			case <-ticker.C:
				i = -1 // Reset the counter.
			}
		}
	}
}

func (t *leakyBucket) accept(ctx context.Context) (bool, func()) {
	c := make(chan struct{})
	select {
	case t.bucket <- c:
		break // Acceptable. Wait in the bucket.
	default:
		close(c)
		return false, nil // Decline. The bucket is full.
	}
	select {
	case <-c:
		return true, noopReleaser // Accepted.
	case <-ctx.Done():
		close(c)
		return false, nil // Request canceled.
	}
}
