// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package throttle

import (
	"context"
	"time"

	"github.com/aileron-gateway/aileron-gateway/internal/txtutil"
	"github.com/aileron-projects/go/ztime/zbackoff"
	"github.com/aileron-projects/go/ztime/zrate"
)

var (
	// backoff is the backoff used for retrying.
	// Backoff algorithm is currently fixed and not configurable.
	backoff = zbackoff.NewExponentialBackoffFullJitter(time.Millisecond, 10*time.Second, time.Millisecond)
)

// noopReleaser do nothing.
// This releaser is intended to be used by
// throttlers except for maxConnectionThrottler.
var noopReleaser = func() {}

// apiThrottler applies throttling for requests
// that matched to the method and paths.
type apiThrottler struct {
	limiter zrate.Limiter
	methods []string
	paths   txtutil.Matcher[string]
	// maxRetry is the maximum retry count.
	// Requests will be rejected when all retries are failed for maxRetry times.
	maxRetry int
}

func (t *apiThrottler) accept(ctx context.Context) (bool, func()) {
	if tk := t.limiter.Accept(ctx); tk.OK() {
		return true, tk.Release
	}
	if t.maxRetry <= 0 {
		return false, noopReleaser
	}
loop:
	for i := 0; i <= t.maxRetry; i++ {
		println(i, ":", backoff.Attempt(i))
		select {
		case <-time.After(backoff.Attempt(i)):
		case <-ctx.Done():
			break loop
		}
		if tk := t.limiter.Accept(ctx); tk.OK() {
			return true, tk.Release
		}
	}
	// Retry failed over maxRetry counts.
	// So, return false and decline the request.
	return false, noopReleaser
}
