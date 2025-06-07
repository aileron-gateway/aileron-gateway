// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package log

import (
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

type testLogger struct {
	Logger
	id string
}

func TestReplaceTime(t *testing.T) {
	type condition struct {
		key    string
		format string
	}

	type action struct {
		format string
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"time",
			[]string{},
			[]string{},
			&condition{
				key:    slog.TimeKey,
				format: time.DateOnly,
			},
			&action{
				format: time.DateTime,
			},
		),
		gen(
			"non time",
			[]string{},
			[]string{},
			&condition{
				key:    "non time",
				format: time.DateOnly,
			},
			&action{
				format: time.DateOnly,
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			now := time.Now()
			attr := slog.Attr{
				Key:   tt.C().key,
				Value: slog.StringValue(now.Format(tt.C().format)),
			}
			attr = replaceTime(nil, attr)
			testutil.Diff(t, now.Format(tt.A().format), attr.Value.String())
		})
	}
}

func TestDefaultOr(t *testing.T) {
	type condition struct {
		name string
	}

	type action struct {
		expect Logger
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"default logger",
			[]string{},
			[]string{},
			&condition{
				name: DefaultLoggerName,
			},
			&action{
				expect: GlobalLogger(DefaultLoggerName),
			},
		),
		gen(
			"test logger",
			[]string{},
			[]string{},
			&condition{
				name: "test",
			},
			&action{
				expect: GlobalLogger(DefaultLoggerName),
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			lg := DefaultOr(tt.C().name)
			testutil.Diff(t, tt.A().expect, lg, cmp.Comparer(testutil.ComparePointer[*SLogger]))
		})
	}
}

func TestSetGlobalLogger(t *testing.T) {
	type condition struct {
		setLogger bool
		logger    Logger
		name      string
	}

	type action struct {
		expect Logger
	}

	CndSetNil := "set nil"
	CndSetNonNil := "set non-nil"
	CndDefaultName := "default name"
	ActCheckReplaced := "check replaced"
	ActCheckStored := "check stored"
	ActCheckNil := "check nil"

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	tb.Condition(CndSetNil, "set nil as logger")
	tb.Condition(CndSetNonNil, "set non-nil logger as input")
	tb.Condition(CndDefaultName, "set logger by default name")
	tb.Action(ActCheckReplaced, "check that logger is replaced")
	tb.Action(ActCheckStored, "check that the logger is stored")
	tb.Action(ActCheckNil, "check that the returned value is nil")
	table := tb.Build()

	testLg := &testLogger{
		Logger: NewJSONSLogger(os.Stdout, nil),
		id:     "test",
	}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"nil by default name",
			[]string{CndSetNil, CndDefaultName},
			[]string{ActCheckReplaced},
			&condition{
				setLogger: true,
				logger:    nil,
				name:      DefaultLoggerName,
			},
			&action{
				expect: NewJSONSLogger(os.Stdout, nil),
			},
		),
		gen(
			"nil by not default name",
			[]string{CndSetNil},
			[]string{ActCheckNil},
			&condition{
				setLogger: true,
				logger:    nil,
				name:      "test",
			},
			&action{
				expect: nil,
			},
		),
		gen(
			"nil by empty name",
			[]string{CndSetNil},
			[]string{ActCheckNil},
			&condition{
				setLogger: true,
				logger:    nil,
				name:      "",
			},
			&action{
				expect: nil,
			},
		),
		gen(
			"non-nil by default name",
			[]string{CndSetNonNil, CndDefaultName},
			[]string{ActCheckReplaced},
			&condition{
				setLogger: true,
				logger:    testLg,
				name:      DefaultLoggerName,
			},
			&action{
				expect: testLg,
			},
		),
		gen(
			"non-nil by not default name",
			[]string{CndSetNonNil},
			[]string{},
			&condition{
				setLogger: true,
				logger:    testLg,
				name:      "test",
			},
			&action{
				expect: testLg,
			},
		),
		gen(
			"non-nil by empty name",
			[]string{CndSetNonNil},
			[]string{},
			&condition{
				setLogger: true,
				logger:    testLg,
				name:      "",
			},
			&action{
				expect: testLg,
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			tmp := GlobalLogger(DefaultLoggerName)
			defer func() {
				SetGlobalLogger(tt.C().name, nil)
				SetGlobalLogger(DefaultLoggerName, tmp)
			}()

			if tt.C().setLogger {
				SetGlobalLogger(tt.C().name, tt.C().logger)
			}

			lg := GlobalLogger(tt.C().name)

			opts := []cmp.Option{
				cmp.AllowUnexported(SLogger{}, slog.Logger{}),
				cmpopts.IgnoreUnexported(slog.TextHandler{}, slog.JSONHandler{}),
				cmp.Comparer(testutil.ComparePointer[*os.File]),
				cmp.Comparer(testutil.ComparePointer[*time.Location]),
			}
			if v, ok := tt.A().expect.(*testLogger); ok {
				testutil.Diff(t, v.id, lg.(*testLogger).id)
			} else {
				testutil.Diff(t, tt.A().expect, lg, opts...)
			}
		})
	}
}

func TestGlobalLogger(t *testing.T) {
	type condition struct {
		name string
	}
	type action struct {
		expect Logger
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	testLg := &testLogger{
		Logger: NewJSONSLogger(os.Stdout, nil),
		id:     "test",
	}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"default name",
			[]string{}, []string{},
			&condition{
				name: DefaultLoggerName,
			},
			&action{
				expect: NewJSONSLogger(os.Stdout, nil),
			},
		),
		gen(
			"not default name",
			[]string{}, []string{},
			&condition{
				name: "test_logger",
			},
			&action{
				expect: testLg,
			},
		),
		gen(
			"not-nil logger",
			[]string{}, []string{},
			&condition{
				name: "not_exist_logger_name",
			},
			&action{
				expect: nil,
			},
		),
		gen(
			"not-nil logger by empty name",
			[]string{}, []string{},
			&condition{
				name: "",
			},
			&action{
				expect: nil,
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			SetGlobalLogger("test_logger", testLg)

			lg := GlobalLogger(tt.C().name)

			opts := []cmp.Option{
				cmp.AllowUnexported(SLogger{}, slog.Logger{}),
				cmpopts.IgnoreUnexported(slog.TextHandler{}, slog.JSONHandler{}),
				cmp.Comparer(testutil.ComparePointer[*os.File]),
				cmp.Comparer(testutil.ComparePointer[*time.Location]),
			}
			if v, ok := tt.A().expect.(*testLogger); ok {
				testutil.Diff(t, v.id, lg.(*testLogger).id)
			} else {
				testutil.Diff(t, tt.A().expect, lg, opts...)
			}
		})
	}
}
