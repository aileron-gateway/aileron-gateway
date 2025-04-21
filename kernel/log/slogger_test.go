// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package log_test

import (
	"log/slog"
	"testing"

	"github.com/aileron-gateway/aileron-gateway/kernel/log"
	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
)

func TestLvToSLogLevel(t *testing.T) {
	type condition struct {
		level log.LogLevel
	}

	type action struct {
		expect slog.Level
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	cndLvTrace := tb.Condition("Trace", "input Trace as  log.LogLevel")
	cndLvDebug := tb.Condition("Debug", "input Debug as  log.LogLevel")
	cndLvInfo := tb.Condition("Info", "input Info as  log.LogLevel")
	cndLvWarn := tb.Condition("Warn", "input Warn as  log.LogLevel")
	cndLvError := tb.Condition("Error", "input Error as  log.LogLevel")
	cndLvFatal := tb.Condition("Fatal", "input Fatal as  log.LogLevel")
	cndLvOther := tb.Condition("Other", "input other log level")
	actLvDebug := tb.Action("Debug", "check that the returned slog.Level is Debug")
	actLvInfo := tb.Action("Info", "check that the returned slog.Level is Info")
	actLvWarn := tb.Action("Warn", "check that the returned slog.Level is Warn")
	actLvError := tb.Action("Error", "check that the returned slog.Level is Error")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"Trace-1",
			[]string{cndLvOther},
			[]string{actLvDebug},
			&condition{
				level: log.LvTrace - 1,
			},
			&action{
				expect: slog.LevelDebug,
			},
		),
		gen(
			"Trace",
			[]string{cndLvTrace},
			[]string{actLvDebug},
			&condition{
				level: log.LvTrace,
			},
			&action{
				expect: slog.LevelDebug,
			},
		),
		gen(
			"Trace+1",
			[]string{cndLvOther},
			[]string{actLvDebug},
			&condition{
				level: log.LvTrace + 1,
			},
			&action{
				expect: slog.LevelDebug,
			},
		),
		gen(
			"Debug-1",
			[]string{cndLvOther},
			[]string{actLvDebug},
			&condition{
				level: log.LvDebug,
			},
			&action{
				expect: slog.LevelDebug,
			},
		),
		gen(
			"Debug",
			[]string{cndLvDebug},
			[]string{actLvDebug},
			&condition{
				level: log.LvDebug,
			},
			&action{
				expect: slog.LevelDebug,
			},
		),
		gen(
			"Debug+1",
			[]string{cndLvOther},
			[]string{actLvInfo},
			&condition{
				level: log.LvDebug + 1,
			},
			&action{
				expect: slog.LevelInfo,
			},
		),
		gen(
			"Info-1",
			[]string{cndLvOther},
			[]string{actLvInfo},
			&condition{
				level: log.LvInfo - 1,
			},
			&action{
				expect: slog.LevelInfo,
			},
		),
		gen(
			"Info",
			[]string{cndLvInfo},
			[]string{actLvInfo},
			&condition{
				level: log.LvInfo,
			},
			&action{
				expect: slog.LevelInfo,
			},
		),
		gen(
			"Info+1",
			[]string{cndLvOther},
			[]string{actLvWarn},
			&condition{
				level: log.LvInfo + 1,
			},
			&action{
				expect: slog.LevelWarn,
			},
		),
		gen(
			"Warn-1",
			[]string{cndLvOther},
			[]string{actLvWarn},
			&condition{
				level: log.LvWarn - 1,
			},
			&action{
				expect: slog.LevelWarn,
			},
		),
		gen(
			"Warn",
			[]string{cndLvWarn},
			[]string{actLvWarn},
			&condition{
				level: log.LvWarn,
			},
			&action{
				expect: slog.LevelWarn,
			},
		),
		gen(
			"Warn+1",
			[]string{cndLvOther},
			[]string{actLvError},
			&condition{
				level: log.LvWarn + 1,
			},
			&action{
				expect: slog.LevelError,
			},
		),
		gen(
			"Error-1",
			[]string{cndLvOther},
			[]string{actLvError},
			&condition{
				level: log.LvError - 1,
			},
			&action{
				expect: slog.LevelError,
			},
		),
		gen(
			"Error",
			[]string{cndLvError},
			[]string{actLvError},
			&condition{
				level: log.LvError,
			},
			&action{
				expect: slog.LevelError,
			},
		),
		gen(
			"Error+1",
			[]string{cndLvOther},
			[]string{actLvError},
			&condition{
				level: log.LvError + 1,
			},
			&action{
				expect: slog.LevelError,
			},
		),
		gen(
			"Fatal-1",
			[]string{cndLvOther},
			[]string{actLvError},
			&condition{
				level: log.LvFatal - 1,
			},
			&action{
				expect: slog.LevelError,
			},
		),
		gen(
			"Fatal",
			[]string{cndLvFatal},
			[]string{actLvError},
			&condition{
				level: log.LvFatal,
			},
			&action{
				expect: slog.LevelError,
			},
		),
		gen(
			"Fatal+1",
			[]string{cndLvOther},
			[]string{actLvError},
			&condition{
				level: log.LvFatal + 1,
			},
			&action{
				expect: slog.LevelError,
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			got := log.LvToSLogLevel(tt.C().level)
			testutil.Diff(t, tt.A().expect, got)
		})
	}
}

func TestLvFromSLogLevel(t *testing.T) {
	type condition struct {
		level slog.Level
	}

	type action struct {
		expect log.LogLevel
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	cndLvDebug := tb.Condition("Debug", "input Debug as slog.Level")
	cndLvInfo := tb.Condition("Info", "input Info as slog.Level")
	cndLvWarn := tb.Condition("Warn", "input Warn as slog.Level")
	cndLvError := tb.Condition("Error", "input Error as slog.Level")
	cndLvOther := tb.Condition("Other", "input other log level")
	actLvDebug := tb.Action("Debug", "check that the returned  log.LogLevel is Debug")
	actLvInfo := tb.Action("Info", "check that the returned  log.LogLevel is Info")
	actLvWarn := tb.Action("Warn", "check that the returned  log.LogLevel is Warn")
	actLvError := tb.Action("Error", "check that the returned  log.LogLevel is Error")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"Debug-1",
			[]string{cndLvOther},
			[]string{actLvDebug},
			&condition{
				level: slog.LevelDebug - 1,
			},
			&action{
				expect: log.LvDebug,
			},
		),
		gen(
			"Debug",
			[]string{cndLvDebug},
			[]string{actLvDebug},
			&condition{
				level: slog.LevelDebug,
			},
			&action{
				expect: log.LvDebug,
			},
		),
		gen(
			"Debug+1",
			[]string{cndLvOther},
			[]string{actLvInfo},
			&condition{
				level: slog.LevelDebug + 1,
			},
			&action{
				expect: log.LvInfo,
			},
		),
		gen(
			"Info-1",
			[]string{cndLvOther},
			[]string{actLvInfo},
			&condition{
				level: slog.LevelInfo - 1,
			},
			&action{
				expect: log.LvInfo,
			},
		),
		gen(
			"Info",
			[]string{cndLvInfo},
			[]string{actLvInfo},
			&condition{
				level: slog.LevelInfo,
			},
			&action{
				expect: log.LvInfo,
			},
		),
		gen(
			"Info+1",
			[]string{cndLvOther},
			[]string{actLvWarn},
			&condition{
				level: slog.LevelInfo + 1,
			},
			&action{
				expect: log.LvWarn,
			},
		),
		gen(
			"Warn-1",
			[]string{cndLvOther},
			[]string{actLvWarn},
			&condition{
				level: slog.LevelWarn - 1,
			},
			&action{
				expect: log.LvWarn,
			},
		),
		gen(
			"Warn",
			[]string{cndLvWarn},
			[]string{actLvWarn},
			&condition{
				level: slog.LevelWarn,
			},
			&action{
				expect: log.LvWarn,
			},
		),
		gen(
			"Warn+1",
			[]string{cndLvOther},
			[]string{actLvError},
			&condition{
				level: slog.LevelWarn + 1,
			},
			&action{
				expect: log.LvError,
			},
		),
		gen(
			"Error-1",
			[]string{cndLvOther},
			[]string{actLvError},
			&condition{
				level: slog.LevelError - 1,
			},
			&action{
				expect: log.LvError,
			},
		),
		gen(
			"Error",
			[]string{cndLvError},
			[]string{actLvError},
			&condition{
				level: slog.LevelError,
			},
			&action{
				expect: log.LvError,
			},
		),
		gen(
			"Error+1",
			[]string{cndLvOther},
			[]string{actLvError},
			&condition{
				level: slog.LevelError + slog.Level(1),
			},
			&action{
				expect: log.LvError,
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			got := log.LvFromSLogLevel(tt.C().level)
			testutil.Diff(t, tt.A().expect, got)
		})
	}
}
