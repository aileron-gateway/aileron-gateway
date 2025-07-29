// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package throttle

import (
	"context"
	"testing"

	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
	"github.com/aileron-projects/go/ztime/zrate"
)

func TestRetryThrottler(t *testing.T) {
	type condition struct {
		throttler *apiThrottler
	}
	type action struct {
		accepted []bool
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"not exceeded/no retry",
			[]string{}, []string{},
			&condition{
				throttler: &apiThrottler{
					limiter:  zrate.NewFixedWindowLimiter(1),
					maxRetry: 0,
				},
			},
			&action{
				accepted: []bool{true},
			},
		),
		gen(
			"exceeded/no retry",
			[]string{}, []string{},
			&condition{
				throttler: &apiThrottler{
					limiter:  zrate.NewFixedWindowLimiter(2),
					maxRetry: 0,
				},
			},
			&action{
				accepted: []bool{true, true, false},
			},
		),
		gen(
			"exceeded/fail for retry",
			[]string{}, []string{},
			&condition{
				throttler: &apiThrottler{
					limiter:  zrate.NewFixedWindowLimiter(2),
					maxRetry: 1,
				},
			},
			&action{
				accepted: []bool{true, true, false},
			},
		),
		gen(
			"exceeded/with retry",
			[]string{}, []string{},
			&condition{
				throttler: &apiThrottler{
					limiter:  zrate.NewFixedWindowLimiter(2),
					maxRetry: 100,
				},
			},
			&action{
				accepted: []bool{true, true, true},
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			rt := tt.C().throttler
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			for _, a := range tt.A().accepted {
				ok, release := rt.accept(ctx)
				if release != nil {
					defer release()
				}
				testutil.Diff(t, a, ok)
			}
		})
	}
}
