// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package entrypoint

import (
	"context"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/aileron-gateway/aileron-gateway/core"
	"github.com/aileron-gateway/aileron-gateway/kernel/log"
	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
)

// Assert interfaces, or interface test.
var _ core.Runner = &channelGroup{}

type testRunner struct {
	called bool
	err    error
}

func (r *testRunner) Run(_ context.Context) error {
	r.called = true
	if r.err != nil {
		return r.err
	}
	time.Sleep(10 * time.Millisecond)
	return nil
}

func TestChannelGroup_Run(t *testing.T) {
	type condition struct {
		runners []core.Runner
	}

	type action struct {
		err        any // error or errorutil.Kind
		errPattern *regexp.Regexp
	}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"no runners",
			&condition{},
			&action{
				err: nil,
			},
		),
		gen(
			"one non-error runner",
			&condition{
				runners: []core.Runner{
					&testRunner{},
				},
			},
			&action{
				err: nil,
			},
		),
		gen(
			"two non-error runners",
			&condition{
				runners: []core.Runner{
					&testRunner{},
					&testRunner{},
				},
			},
			&action{
				err: nil,
			},
		),
		gen(
			"one error runner",
			&condition{
				runners: []core.Runner{
					&testRunner{err: errors.New("test error")},
				},
			},
			&action{
				err:        core.ErrCoreEntrypointRun,
				errPattern: regexp.MustCompile(core.ErrPrefix + `error on running entrypoint .*\[test error\]`),
			},
		),
		gen(
			"error and non-error runners",
			&condition{
				runners: []core.Runner{
					&testRunner{},
					&testRunner{err: errors.New("test error")},
				},
			},
			&action{
				err:        core.ErrCoreEntrypointRun,
				errPattern: regexp.MustCompile(core.ErrPrefix + `error on running entrypoint .*\[test error\]`),
			},
		),
		gen(
			"two error runners",
			&condition{
				runners: []core.Runner{
					&testRunner{err: errors.New("test error")},
					&testRunner{err: errors.New("test error")},
				},
			},
			&action{
				err:        core.ErrCoreEntrypointRun,
				errPattern: regexp.MustCompile(core.ErrPrefix + `error on running entrypoint .*\[test error\]`),
			},
		),
	}

	for _, tt := range testCases {
		tt := tt

		t.Run(tt.Name, func(t *testing.T) {
			cg := &channelGroup{
				lg:      log.GlobalLogger(log.DefaultLoggerName),
				runners: tt.C.runners,
			}

			err := cg.Run(context.Background())
			testutil.DiffError(t, tt.A.err, tt.A.errPattern, err)
		})
	}
}

func TestChannelGroup_Finalize(t *testing.T) {
	type condition struct {
		runner *channelGroup
	}

	type action struct {
	}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"empty",
			&condition{
				runner: &channelGroup{
					finalizers: nil,
				},
			},
			&action{},
		),
		gen(
			"1 finalizer",
			&condition{
				runner: &channelGroup{
					finalizers: []core.Finalizer{
						&testFinalizer{},
					},
				},
			},
			&action{},
		),
		gen(
			"1 finalizer/error",
			&condition{
				runner: &channelGroup{
					finalizers: []core.Finalizer{
						&testFinalizer{err: errors.New("test error")},
					},
				},
			},
			&action{},
		),
		gen(
			"2 finalizers",
			&condition{
				runner: &channelGroup{
					finalizers: []core.Finalizer{
						&testFinalizer{},
						&testFinalizer{},
					},
				},
			},
			&action{},
		),
		gen(
			"2 finalizers/error",
			&condition{
				runner: &channelGroup{
					finalizers: []core.Finalizer{
						&testFinalizer{},
						&testFinalizer{err: errors.New("test error")},
					},
				},
			},
			&action{},
		),
		gen(
			"nil finalizer",
			&condition{
				runner: &channelGroup{
					finalizers: []core.Finalizer{
						&testFinalizer{},
						nil,
						&testFinalizer{},
					},
				},
			},
			&action{},
		),
	}

	for _, tt := range testCases {
		tt := tt

		t.Run(tt.Name, func(t *testing.T) {
			tt.C.runner.finalize()
			for _, f := range tt.C.runner.finalizers {
				if f == nil {
					continue
				}
				ff := f.(*testFinalizer)
				testutil.Diff(t, true, ff.called)
			}
		})
	}
}
