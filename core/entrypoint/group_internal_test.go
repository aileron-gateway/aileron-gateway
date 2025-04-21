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
var (
	_ = core.Runner(&channelGroup{})
	_ = core.Runner(&waitGroup{})
)

type testRunner struct {
	name   string
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

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	cndNoRunners := tb.Condition("no runners", "no runners")
	cndErrorRunner := tb.Condition("error runner", "at least one runners are contained which returns an error in the runners")
	actCheckNoError := tb.Action("no error", "check that there is no error returned")
	actCheckError := tb.Action("error", "check that an expected error was returned")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"no runners",
			[]string{cndNoRunners},
			[]string{actCheckNoError},
			&condition{},
			&action{
				err: nil,
			},
		),
		gen(
			"one non-error runner",
			[]string{},
			[]string{actCheckNoError},
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
			[]string{},
			[]string{actCheckNoError},
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
			[]string{cndErrorRunner},
			[]string{actCheckError},
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
			[]string{cndErrorRunner},
			[]string{actCheckError},
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
			[]string{cndErrorRunner},
			[]string{actCheckError},
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

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt

		t.Run(tt.Name(), func(t *testing.T) {
			cg := &channelGroup{
				lg:      log.GlobalLogger(log.DefaultLoggerName),
				runners: tt.C().runners,
			}

			err := cg.Run(context.Background())
			testutil.DiffError(t, tt.A().err, tt.A().errPattern, err)
		})
	}
}

func TestChannelGroup_Finalize(t *testing.T) {
	type condition struct {
		runner *channelGroup
	}

	type action struct {
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	cndEmpty := tb.Condition("no finalizer", "no finalizer was specified")
	cndContainNil := tb.Condition("nil finalizer", "nil finalizer was contained")
	cndError := tb.Condition("error finalizer", "finalizer returns an error")
	actCheckCalled := tb.Action("called", "check that all finalizer were called")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"empty",
			[]string{cndEmpty},
			[]string{actCheckCalled},
			&condition{
				runner: &channelGroup{
					finalizers: nil,
				},
			},
			&action{},
		),
		gen(
			"1 finalizer",
			[]string{},
			[]string{actCheckCalled},
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
			[]string{cndError},
			[]string{actCheckCalled},
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
			[]string{},
			[]string{actCheckCalled},
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
			[]string{cndError},
			[]string{actCheckCalled},
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
			[]string{cndContainNil},
			[]string{actCheckCalled},
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

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt

		t.Run(tt.Name(), func(t *testing.T) {
			tt.C().runner.finalize()
			for _, f := range tt.C().runner.finalizers {
				if f == nil {
					continue
				}
				ff := f.(*testFinalizer)
				testutil.Diff(t, true, ff.called)
			}
		})
	}
}

func TestWaitGroup_Run(t *testing.T) {
	type condition struct {
		runners []core.Runner
	}

	type action struct {
		err        any // error or errorutil.Kind
		errPattern *regexp.Regexp
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	cndNoRunners := tb.Condition("no runners", "no runners")
	cndErrorRunner := tb.Condition("error runner", "at least one runners are contained which returns an error in the runners")
	actCheckNoError := tb.Action("no error", "check that there is no error returned")
	actCheckError := tb.Action("error", "check that an expected error was returned")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"no runners",
			[]string{cndNoRunners},
			[]string{actCheckNoError},
			&condition{},
			&action{
				err: nil,
			},
		),
		gen(
			"one non-error runner",
			[]string{},
			[]string{actCheckNoError},
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
			[]string{},
			[]string{actCheckNoError},
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
			[]string{cndErrorRunner},
			[]string{actCheckError},
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
			"two error and non-error runners",
			[]string{cndErrorRunner},
			[]string{actCheckError},
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
			[]string{cndErrorRunner},
			[]string{actCheckError},
			&condition{
				runners: []core.Runner{
					&testRunner{err: errors.New("test error")},
					&testRunner{err: errors.New("test error")},
				},
			},
			&action{
				err:        core.ErrCoreEntrypointRun,
				errPattern: regexp.MustCompile(core.ErrPrefix + `error on running entrypoint .*\[test error\ntest error\]`),
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			wg := &waitGroup{
				lg:      log.GlobalLogger(log.DefaultLoggerName),
				runners: tt.C().runners,
			}

			err := wg.Run(context.Background())
			testutil.DiffError(t, tt.A().err, tt.A().errPattern, err)
		})
	}
}

func TestWaitGroup_Finalize(t *testing.T) {
	type condition struct {
		runner *waitGroup
	}

	type action struct {
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	cndEmpty := tb.Condition("no finalizer", "no finalizer was specified")
	cndContainNil := tb.Condition("nil finalizer", "nil finalizer was contained")
	cndError := tb.Condition("error finalizer", "finalizer returns an error")
	actCheckCalled := tb.Action("called", "check that all finalizer were called")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"empty",
			[]string{cndEmpty},
			[]string{actCheckCalled},
			&condition{
				runner: &waitGroup{
					finalizers: nil,
				},
			},
			&action{},
		),
		gen(
			"1 finalizer",
			[]string{},
			[]string{actCheckCalled},
			&condition{
				runner: &waitGroup{
					finalizers: []core.Finalizer{
						&testFinalizer{},
					},
				},
			},
			&action{},
		),
		gen(
			"1 finalizer/error",
			[]string{cndError},
			[]string{actCheckCalled},
			&condition{
				runner: &waitGroup{
					finalizers: []core.Finalizer{
						&testFinalizer{err: errors.New("test error")},
					},
				},
			},
			&action{},
		),
		gen(
			"2 finalizers",
			[]string{},
			[]string{actCheckCalled},
			&condition{
				runner: &waitGroup{
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
			[]string{cndError},
			[]string{actCheckCalled},
			&condition{
				runner: &waitGroup{
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
			[]string{cndContainNil},
			[]string{actCheckCalled},
			&condition{
				runner: &waitGroup{
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

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt

		t.Run(tt.Name(), func(t *testing.T) {
			tt.C().runner.finalize()
			for _, f := range tt.C().runner.finalizers {
				if f == nil {
					continue
				}
				ff := f.(*testFinalizer)
				testutil.Diff(t, true, ff.called)
			}
		})
	}
}
