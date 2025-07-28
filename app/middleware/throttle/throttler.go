// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package throttle

import (
	"context"
	"time"

	"github.com/aileron-gateway/aileron-gateway/kernel/txtutil"
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
	limiter  zrate.Limiter
	allowNow bool
	methods  []string
	paths    txtutil.Matcher[string]
	// maxRetry is the maximum retry count.
	// Requests will be rejected when all retries are failed for maxRetry times.
	maxRetry int
}

func (t *apiThrottler) token(ctx context.Context) zrate.Token {
	if t.allowNow {
		return t.limiter.AllowNow()
	}
	return t.limiter.WaitNow(ctx)
}

func (t *apiThrottler) accept(ctx context.Context) (bool, func()) {
	if tk := t.token(ctx); tk.OK() {
		return true, tk.Release
	}
	if t.maxRetry <= 0 {
		return false, noopReleaser
	}

	ticker := time.NewTicker(time.Millisecond)
	defer ticker.Stop()

loop:
	for i := 1; i <= t.maxRetry; i++ {
		ticker.Reset(backoff.Attempt(i))
		select {
		case <-ticker.C:
		case <-ctx.Done():
			break loop
		}
		// Check if the request is now acceptable or not.
		if tk := t.token(ctx); tk.OK() {
			return true, tk.Release
		}
	}
	// Retry failed over maxRetry counts.
	// So, return false and decline the request.
	return false, noopReleaser
}
