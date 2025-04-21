// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package log

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"os"
	"strings"
	"testing"

	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestNewTextSLogger(t *testing.T) {
	type condition struct {
		opts *slog.HandlerOptions
		w    io.Writer
	}

	type action struct {
		w  io.Writer
		lv LogLevel
		lg *slog.Logger
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	cndLvDebug := tb.Condition("debug level", "set logger to debug level")
	cndLvInfo := tb.Condition("info level", "set logger to info level")
	cndLvWarn := tb.Condition("warn level", "set logger to warn level")
	cndLvError := tb.Condition("error level", "set logger to error level")
	cndStdout := tb.Condition("Stdout", "set stdout as writer")
	cndStderr := tb.Condition("Stderr", "set stderr as writer")
	cndNilWriter := tb.Condition("nil writer", "set nil as writer")
	cndNilOption := tb.Condition("nil option", "set nil as option")
	actCheckLogs := tb.Action("check logs", "check that the expected log was output")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"debug/stdout",
			[]string{cndLvDebug, cndStdout},
			[]string{actCheckLogs},
			&condition{
				w:    os.Stdout,
				opts: &slog.HandlerOptions{Level: slog.LevelDebug},
			},
			&action{
				w:  os.Stdout,
				lv: LvDebug,
				lg: slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})),
			},
		),
		gen(
			"info/stdout",
			[]string{cndLvInfo, cndStdout},
			[]string{actCheckLogs},
			&condition{
				w:    os.Stdout,
				opts: &slog.HandlerOptions{Level: slog.LevelInfo},
			},
			&action{
				w:  os.Stdout,
				lv: LvInfo,
				lg: slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})),
			},
		),
		gen(
			"warn/stdout",
			[]string{cndLvWarn, cndStdout},
			[]string{actCheckLogs},
			&condition{
				w:    os.Stdout,
				opts: &slog.HandlerOptions{Level: slog.LevelWarn},
			},
			&action{
				w:  os.Stdout,
				lv: LvWarn,
				lg: slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelWarn})),
			},
		),
		gen(
			"error/stdout",
			[]string{cndLvError, cndStdout},
			[]string{actCheckLogs},
			&condition{
				w:    os.Stdout,
				opts: &slog.HandlerOptions{Level: slog.LevelError},
			},
			&action{
				w:  os.Stdout,
				lv: LvError,
				lg: slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError})),
			},
		),
		gen(
			"stderr",
			[]string{cndLvError, cndStderr},
			[]string{actCheckLogs},
			&condition{
				w:    os.Stderr,
				opts: &slog.HandlerOptions{Level: slog.LevelError},
			},
			&action{
				w:  os.Stderr,
				lv: LvError,
				lg: slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError})),
			},
		),
		gen(
			"nil writer",
			[]string{cndLvInfo, cndNilWriter},
			[]string{actCheckLogs},
			&condition{
				w:    nil,
				opts: &slog.HandlerOptions{Level: slog.LevelError},
			},
			&action{
				w:  os.Stdout,
				lv: LvError,
				lg: slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError})),
			},
		),
		gen(
			"nil option",
			[]string{cndLvInfo, cndNilOption},
			[]string{actCheckLogs},
			&condition{
				w:    os.Stdout,
				opts: nil,
			},
			&action{
				w:  os.Stdout,
				lv: LvInfo,
				lg: slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})),
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			got := NewTextSLogger(tt.C().w, tt.C().opts)

			testutil.Diff(t, tt.A().w, got.w, cmp.Comparer(testutil.ComparePointer[io.Writer]))
			testutil.Diff(t, tt.A().lg, got.lg, cmp.AllowUnexported(slog.Logger{}), cmpopts.IgnoreUnexported(slog.TextHandler{}))
			testutil.Diff(t, tt.A().lv, got.lv)
		})
	}
}

func TestNewJSONSLogger(t *testing.T) {
	type condition struct {
		opts *slog.HandlerOptions
		w    io.Writer
	}

	type action struct {
		w  io.Writer
		lv LogLevel
		lg *slog.Logger
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	cndLvDebug := tb.Condition("debug level", "set logger to debug level")
	cndLvInfo := tb.Condition("info level", "set logger to info level")
	cndLvWarn := tb.Condition("warn level", "set logger to warn level")
	cndLvError := tb.Condition("error level", "set logger to error level")
	cndStdout := tb.Condition("Stdout", "set stdout as writer")
	cndStderr := tb.Condition("Stderr", "set stderr as writer")
	cndNilWriter := tb.Condition("nil writer", "set nil as writer")
	cndNilOption := tb.Condition("nil option", "set nil as option")
	actCheckLogs := tb.Action("check logs", "check that the expected log was output")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"debug/stdout",
			[]string{cndLvDebug, cndStdout},
			[]string{actCheckLogs},
			&condition{
				w:    os.Stdout,
				opts: &slog.HandlerOptions{Level: slog.LevelDebug},
			},
			&action{
				w:  os.Stdout,
				lv: LvDebug,
				lg: slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})),
			},
		),
		gen(
			"info/stdout",
			[]string{cndLvInfo, cndStdout},
			[]string{actCheckLogs},
			&condition{
				w:    os.Stdout,
				opts: &slog.HandlerOptions{Level: slog.LevelInfo},
			},
			&action{
				w:  os.Stdout,
				lv: LvInfo,
				lg: slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})),
			},
		),
		gen(
			"warn/stdout",
			[]string{cndLvWarn, cndStdout},
			[]string{actCheckLogs},
			&condition{
				w:    os.Stdout,
				opts: &slog.HandlerOptions{Level: slog.LevelWarn},
			},
			&action{
				w:  os.Stdout,
				lv: LvWarn,
				lg: slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelWarn})),
			},
		),
		gen(
			"error/stdout",
			[]string{cndLvError, cndStdout},
			[]string{actCheckLogs},
			&condition{
				w:    os.Stdout,
				opts: &slog.HandlerOptions{Level: slog.LevelError},
			},
			&action{
				w:  os.Stdout,
				lv: LvError,
				lg: slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError})),
			},
		),
		gen(
			"stderr",
			[]string{cndLvError, cndStderr},
			[]string{actCheckLogs},
			&condition{
				w:    os.Stderr,
				opts: &slog.HandlerOptions{Level: slog.LevelError},
			},
			&action{
				w:  os.Stderr,
				lv: LvError,
				lg: slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError})),
			},
		),
		gen(
			"nil writer",
			[]string{cndLvInfo, cndNilWriter},
			[]string{actCheckLogs},
			&condition{
				w:    nil,
				opts: &slog.HandlerOptions{Level: slog.LevelError},
			},
			&action{
				w:  os.Stdout,
				lv: LvError,
				lg: slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError})),
			},
		),
		gen(
			"nil option",
			[]string{cndLvInfo, cndNilOption},
			[]string{actCheckLogs},
			&condition{
				w:    os.Stdout,
				opts: nil,
			},
			&action{
				w:  os.Stdout,
				lv: LvInfo,
				lg: slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})),
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			got := NewJSONSLogger(tt.C().w, tt.C().opts)

			testutil.Diff(t, tt.A().w, got.w, cmp.Comparer(testutil.ComparePointer[io.Writer]))
			testutil.Diff(t, tt.A().lg, got.lg, cmp.AllowUnexported(slog.Logger{}), cmpopts.IgnoreUnexported(slog.JSONHandler{}))
			testutil.Diff(t, tt.A().lv, got.lv)
		})
	}
}

func TestEnabled(t *testing.T) {
	type condition struct {
		level slog.Level
	}

	type action struct {
		enabled  []LogLevel
		disabled []LogLevel
	}

	cndLvDebug := "debug level"
	cndLvInfo := "info level"
	cndLvWarn := "warn level"
	cndLvError := "error level"
	actCheckEnabled := "check enabled"

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	tb.Condition(cndLvDebug, "set logger to debug level")
	tb.Condition(cndLvInfo, "set logger to info level")
	tb.Condition(cndLvWarn, "set logger to warn level")
	tb.Condition(cndLvError, "set logger to error level")
	tb.Action(actCheckEnabled, "check if the right value was returned for Enabled")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"debug",
			[]string{cndLvDebug},
			[]string{actCheckEnabled},
			&condition{
				level: slog.LevelDebug,
			},
			&action{
				enabled:  []LogLevel{LvDebug, LvInfo, LvWarn, LvError, LvFatal},
				disabled: []LogLevel{LvTrace},
			},
		),
		gen(
			"info",
			[]string{cndLvDebug},
			[]string{actCheckEnabled},
			&condition{
				level: slog.LevelInfo,
			},
			&action{
				enabled:  []LogLevel{LvInfo, LvWarn, LvError, LvFatal},
				disabled: []LogLevel{LvTrace, LvDebug},
			},
		),
		gen(
			"warn",
			[]string{cndLvWarn},
			[]string{actCheckEnabled},
			&condition{
				level: slog.LevelWarn,
			},
			&action{
				enabled:  []LogLevel{LvWarn, LvError, LvFatal},
				disabled: []LogLevel{LvTrace, LvDebug, LvInfo},
			},
		),
		gen(
			"error",
			[]string{cndLvInfo},
			[]string{actCheckEnabled},
			&condition{
				level: slog.LevelError,
			},
			&action{
				enabled:  []LogLevel{LvError, LvFatal},
				disabled: []LogLevel{LvTrace, LvDebug, LvInfo, LvWarn}},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			var buf bytes.Buffer
			lg := &SLogger{
				lg: slog.New(slog.NewJSONHandler(
					&buf,
					&slog.HandlerOptions{
						Level: tt.C().level,
					},
				)),
				lv: LvFromSLogLevel(tt.C().level),
			}

			for _, lv := range tt.A().enabled {
				testutil.Diff(t, true, lg.Enabled(lv))
			}
			for _, lv := range tt.A().disabled {
				testutil.Diff(t, false, lg.Enabled(lv))
			}
		})
	}
}

func TestDebug(t *testing.T) {
	type condition struct {
		level LogLevel
		kvs   []any
	}

	type action struct {
		noOutput bool
	}

	cndLvTrace := "trace level"
	cndLvDebug := "debug level"
	cndLvInfo := "info level"
	cndLvWarn := "warn level"
	cndLvError := "error level"
	cndLvFatal := "fatal level"
	actCheckLogs := "check logs"

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	tb.Condition(cndLvTrace, "set logger to trace level")
	tb.Condition(cndLvDebug, "set logger to debug level")
	tb.Condition(cndLvInfo, "set logger to info level")
	tb.Condition(cndLvWarn, "set logger to warn level")
	tb.Condition(cndLvError, "set logger to error level")
	tb.Condition(cndLvFatal, "set logger to fatal level")
	tb.Action(actCheckLogs, "check that the expected log was output")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"trace",
			[]string{cndLvTrace},
			[]string{actCheckLogs},
			&condition{
				level: LvTrace,
				kvs:   []any{"A", "B"},
			},
			&action{
				noOutput: false,
			},
		),
		gen(
			"debug",
			[]string{cndLvDebug},
			[]string{actCheckLogs},
			&condition{
				level: LvDebug,
				kvs:   []any{"A", "B"},
			},
			&action{
				noOutput: false,
			},
		),
		gen(
			"info",
			[]string{cndLvDebug},
			[]string{actCheckLogs},
			&condition{
				level: LvInfo,
				kvs:   []any{"A", "B"},
			},
			&action{
				noOutput: true,
			},
		),
		gen(
			"warn",
			[]string{cndLvWarn},
			[]string{actCheckLogs},
			&condition{
				level: LvWarn,
				kvs:   []any{"A", "B"},
			},
			&action{
				noOutput: true,
			},
		),
		gen(
			"error",
			[]string{cndLvError},
			[]string{actCheckLogs},
			&condition{
				level: LvError,
				kvs:   []any{"A", "B"},
			},
			&action{
				noOutput: true,
			},
		),
		gen(
			"fatal",
			[]string{cndLvFatal},
			[]string{actCheckLogs},
			&condition{
				level: LvFatal,
				kvs:   []any{"A", "B"},
			},
			&action{
				noOutput: true,
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			var buf bytes.Buffer
			lg := &SLogger{
				w: &buf,
				lg: slog.New(slog.NewJSONHandler(
					&buf,
					&slog.HandlerOptions{
						Level: LvToSLogLevel(tt.C().level),
					},
				)),
				lv: tt.C().level,
			}

			msg := "test message"
			ctx := context.Background()
			ctx = ContextWithAttrs(ctx, NewCustomAttrs("test", map[string]any{"foo": "bar"}))
			lg.Debug(ctx, msg, tt.C().kvs...)

			if tt.A().noOutput {
				testutil.Diff(t, "", buf.String())
				return
			}

			// Get log attributes from structured json
			attrs := map[string]any{}
			testutil.Diff(t, nil, json.Unmarshal(buf.Bytes(), &attrs))
			testutil.Diff(t, Debug, attrs["level"])
			testutil.Diff(t, msg, attrs["msg"])
			testutil.Diff(t, map[string]any{"foo": "bar"}, attrs["test"])
			for i := 0; i < len(tt.C().kvs); i += 2 {
				testutil.Diff(t, tt.C().kvs[i+1], attrs[tt.C().kvs[i].(string)])
			}

			lg.Write([]byte("always written"))
			t.Log(buf.String())
			testutil.Diff(t, true, strings.HasSuffix(buf.String(), "always written"))
		})
	}
}

func TestInfo(t *testing.T) {
	type condition struct {
		level LogLevel
		kvs   []any
	}

	type action struct {
		noOutput bool
	}

	cndLvTrace := "trace level"
	cndLvDebug := "debug level"
	cndLvInfo := "info level"
	cndLvWarn := "warn level"
	cndLvError := "error level"
	cndLvFatal := "fatal level"
	actCheckLogs := "check logs"

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	tb.Condition(cndLvTrace, "set logger to trace level")
	tb.Condition(cndLvDebug, "set logger to debug level")
	tb.Condition(cndLvInfo, "set logger to info level")
	tb.Condition(cndLvWarn, "set logger to warn level")
	tb.Condition(cndLvError, "set logger to error level")
	tb.Condition(cndLvFatal, "set logger to fatal level")
	tb.Action(actCheckLogs, "check that the expected log was output")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"trace",
			[]string{cndLvTrace},
			[]string{actCheckLogs},
			&condition{
				level: LvTrace,
				kvs:   []any{"A", "B"},
			},
			&action{
				noOutput: false,
			},
		),
		gen(
			"debug",
			[]string{cndLvDebug},
			[]string{actCheckLogs},
			&condition{
				level: LvDebug,
				kvs:   []any{"A", "B"},
			},
			&action{
				noOutput: false,
			},
		),
		gen(
			"info",
			[]string{cndLvDebug},
			[]string{actCheckLogs},
			&condition{
				level: LvInfo,
				kvs:   []any{"A", "B"},
			},
			&action{
				noOutput: false,
			},
		),
		gen(
			"warn",
			[]string{cndLvWarn},
			[]string{actCheckLogs},
			&condition{
				level: LvWarn,
				kvs:   []any{"A", "B"},
			},
			&action{
				noOutput: true,
			},
		),
		gen(
			"error",
			[]string{cndLvError},
			[]string{actCheckLogs},
			&condition{
				level: LvError,
				kvs:   []any{"A", "B"},
			},
			&action{
				noOutput: true,
			},
		),
		gen(
			"fatal",
			[]string{cndLvFatal},
			[]string{actCheckLogs},
			&condition{
				level: LvFatal,
				kvs:   []any{"A", "B"},
			},
			&action{
				noOutput: true,
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			var buf bytes.Buffer
			lg := &SLogger{
				w: &buf,
				lg: slog.New(slog.NewJSONHandler(
					&buf,
					&slog.HandlerOptions{
						Level: LvToSLogLevel(tt.C().level),
					},
				)),
				lv: tt.C().level,
			}

			msg := "test message"
			ctx := context.Background()
			ctx = ContextWithAttrs(ctx, NewCustomAttrs("test", map[string]any{"foo": "bar"}))
			lg.Info(ctx, msg, tt.C().kvs...)

			if tt.A().noOutput {
				testutil.Diff(t, "", buf.String())
				return
			}

			// Get log attributes from structured json
			attrs := map[string]any{}
			testutil.Diff(t, nil, json.Unmarshal(buf.Bytes(), &attrs))

			testutil.Diff(t, Info, attrs["level"])
			testutil.Diff(t, msg, attrs["msg"])
			testutil.Diff(t, map[string]any{"foo": "bar"}, attrs["test"])
			for i := 0; i < len(tt.C().kvs); i += 2 {
				testutil.Diff(t, tt.C().kvs[i+1], attrs[tt.C().kvs[i].(string)])
			}

			lg.Write([]byte("always written"))
			t.Log(buf.String())
			testutil.Diff(t, true, strings.HasSuffix(buf.String(), "always written"))
		})
	}
}

func TestWarn(t *testing.T) {
	type condition struct {
		level LogLevel
		kvs   []any
	}

	type action struct {
		noOutput bool
	}

	cndLvTrace := "trace level"
	cndLvDebug := "debug level"
	cndLvInfo := "info level"
	cndLvWarn := "warn level"
	cndLvError := "error level"
	cndLvFatal := "fatal level"
	actCheckLogs := "check logs"

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	tb.Condition(cndLvTrace, "set logger to trace level")
	tb.Condition(cndLvDebug, "set logger to debug level")
	tb.Condition(cndLvInfo, "set logger to info level")
	tb.Condition(cndLvWarn, "set logger to warn level")
	tb.Condition(cndLvError, "set logger to error level")
	tb.Condition(cndLvFatal, "set logger to fatal level")
	tb.Action(actCheckLogs, "check that the expected log was output")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"trace",
			[]string{cndLvTrace},
			[]string{actCheckLogs},
			&condition{
				level: LvTrace,
				kvs:   []any{"A", "B"},
			},
			&action{
				noOutput: false,
			},
		),
		gen(
			"debug",
			[]string{cndLvDebug},
			[]string{actCheckLogs},
			&condition{
				level: LvDebug,
				kvs:   []any{"A", "B"},
			},
			&action{
				noOutput: false,
			},
		),
		gen(
			"info",
			[]string{cndLvDebug},
			[]string{actCheckLogs},
			&condition{
				level: LvInfo,
				kvs:   []any{"A", "B"},
			},
			&action{
				noOutput: false,
			},
		),
		gen(
			"warn",
			[]string{cndLvWarn},
			[]string{actCheckLogs},
			&condition{
				level: LvWarn,
				kvs:   []any{"A", "B"},
			},
			&action{
				noOutput: false,
			},
		),
		gen(
			"error",
			[]string{cndLvError},
			[]string{actCheckLogs},
			&condition{
				level: LvError,
				kvs:   []any{"A", "B"},
			},
			&action{
				noOutput: true,
			},
		),
		gen(
			"fatal",
			[]string{cndLvFatal},
			[]string{actCheckLogs},
			&condition{
				level: LvFatal,
				kvs:   []any{"A", "B"},
			},
			&action{
				noOutput: true,
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			var buf bytes.Buffer
			lg := &SLogger{
				w: &buf,
				lg: slog.New(slog.NewJSONHandler(
					&buf,
					&slog.HandlerOptions{
						Level: LvToSLogLevel(tt.C().level),
					},
				)),
				lv: tt.C().level,
			}

			msg := "test message"
			ctx := context.Background()
			ctx = ContextWithAttrs(ctx, NewCustomAttrs("test", map[string]any{"foo": "bar"}))
			lg.Warn(ctx, msg, tt.C().kvs...)

			if tt.A().noOutput {
				testutil.Diff(t, "", buf.String())
				return
			}

			// Get log attributes from structured json
			attrs := map[string]any{}
			testutil.Diff(t, nil, json.Unmarshal(buf.Bytes(), &attrs))

			testutil.Diff(t, Warn, attrs["level"])
			testutil.Diff(t, msg, attrs["msg"])
			testutil.Diff(t, map[string]any{"foo": "bar"}, attrs["test"])
			for i := 0; i < len(tt.C().kvs); i += 2 {
				testutil.Diff(t, tt.C().kvs[i+1], attrs[tt.C().kvs[i].(string)])
			}

			lg.Write([]byte("always written"))
			t.Log(buf.String())
			testutil.Diff(t, true, strings.HasSuffix(buf.String(), "always written"))
		})
	}
}

func TestError(t *testing.T) {
	type condition struct {
		level LogLevel
		kvs   []any
	}

	type action struct {
		noOutput bool
	}

	cndLvTrace := "trace level"
	cndLvDebug := "debug level"
	cndLvInfo := "info level"
	cndLvWarn := "warn level"
	cndLvError := "error level"
	cndLvFatal := "fatal level"
	actCheckLogs := "check logs"

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	tb.Condition(cndLvTrace, "set logger to trace level")
	tb.Condition(cndLvDebug, "set logger to debug level")
	tb.Condition(cndLvInfo, "set logger to info level")
	tb.Condition(cndLvWarn, "set logger to warn level")
	tb.Condition(cndLvError, "set logger to error level")
	tb.Condition(cndLvFatal, "set logger to fatal level")
	tb.Action(actCheckLogs, "check that the expected log was output")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"trace",
			[]string{cndLvTrace},
			[]string{actCheckLogs},
			&condition{
				level: LvTrace,
				kvs:   []any{"A", "B"},
			},
			&action{
				noOutput: false,
			},
		),
		gen(
			"debug",
			[]string{cndLvDebug},
			[]string{actCheckLogs},
			&condition{
				level: LvDebug,
				kvs:   []any{"A", "B"},
			},
			&action{
				noOutput: false,
			},
		),
		gen(
			"info",
			[]string{cndLvDebug},
			[]string{actCheckLogs},
			&condition{
				level: LvInfo,
				kvs:   []any{"A", "B"},
			},
			&action{
				noOutput: false,
			},
		),
		gen(
			"warn",
			[]string{cndLvWarn},
			[]string{actCheckLogs},
			&condition{
				level: LvWarn,
				kvs:   []any{"A", "B"},
			},
			&action{
				noOutput: false,
			},
		),
		gen(
			"error",
			[]string{cndLvError},
			[]string{actCheckLogs},
			&condition{
				level: LvError,
				kvs:   []any{"A", "B"},
			},
			&action{
				noOutput: false,
			},
		),
		gen(
			"fatal",
			[]string{cndLvFatal},
			[]string{actCheckLogs},
			&condition{
				level: LvFatal,
				kvs:   []any{"A", "B"},
			},
			&action{
				noOutput: true,
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			var buf bytes.Buffer
			lg := &SLogger{
				w: &buf,
				lg: slog.New(slog.NewJSONHandler(
					&buf,
					&slog.HandlerOptions{
						Level: LvToSLogLevel(tt.C().level),
					},
				)),
				lv: tt.C().level,
			}

			msg := "test message"
			ctx := context.Background()
			ctx = ContextWithAttrs(ctx, NewCustomAttrs("test", map[string]any{"foo": "bar"}))
			lg.Error(ctx, msg, tt.C().kvs...)

			if tt.A().noOutput {
				testutil.Diff(t, "", buf.String())
				return
			}
			// t.Error(buf.String())
			// Get log attributes from structured json
			attrs := map[string]any{}
			testutil.Diff(t, nil, json.Unmarshal(buf.Bytes(), &attrs))

			testutil.Diff(t, Error, attrs["level"])
			testutil.Diff(t, msg, attrs["msg"])
			testutil.Diff(t, map[string]any{"foo": "bar"}, attrs["test"])
			for i := 0; i < len(tt.C().kvs); i += 2 {
				testutil.Diff(t, tt.C().kvs[i+1], attrs[tt.C().kvs[i].(string)])
			}

			lg.Write([]byte("always written"))
			t.Log(buf.String())
			testutil.Diff(t, true, strings.HasSuffix(buf.String(), "always written"))
		})
	}
}
