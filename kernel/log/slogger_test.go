// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package log_test

import (
	"log/slog"
	"testing"

	"github.com/aileron-gateway/aileron-gateway/internal/testutil"
	"github.com/aileron-gateway/aileron-gateway/kernel/log"
)

func TestLvToSLogLevel(t *testing.T) {
	type condition struct {
		level log.LogLevel
	}

	type action struct {
		expect slog.Level
	}
	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"Trace-1",
			&condition{
				level: log.LvTrace - 1,
			},
			&action{
				expect: slog.LevelDebug,
			},
		),
		gen(
			"Trace",
			&condition{
				level: log.LvTrace,
			},
			&action{
				expect: slog.LevelDebug,
			},
		),
		gen(
			"Trace+1",
			&condition{
				level: log.LvTrace + 1,
			},
			&action{
				expect: slog.LevelDebug,
			},
		),
		gen(
			"Debug-1",
			&condition{
				level: log.LvDebug,
			},
			&action{
				expect: slog.LevelDebug,
			},
		),
		gen(
			"Debug",
			&condition{
				level: log.LvDebug,
			},
			&action{
				expect: slog.LevelDebug,
			},
		),
		gen(
			"Debug+1",
			&condition{
				level: log.LvDebug + 1,
			},
			&action{
				expect: slog.LevelInfo,
			},
		),
		gen(
			"Info-1",
			&condition{
				level: log.LvInfo - 1,
			},
			&action{
				expect: slog.LevelInfo,
			},
		),
		gen(
			"Info",
			&condition{
				level: log.LvInfo,
			},
			&action{
				expect: slog.LevelInfo,
			},
		),
		gen(
			"Info+1",
			&condition{
				level: log.LvInfo + 1,
			},
			&action{
				expect: slog.LevelWarn,
			},
		),
		gen(
			"Warn-1",
			&condition{
				level: log.LvWarn - 1,
			},
			&action{
				expect: slog.LevelWarn,
			},
		),
		gen(
			"Warn",
			&condition{
				level: log.LvWarn,
			},
			&action{
				expect: slog.LevelWarn,
			},
		),
		gen(
			"Warn+1",
			&condition{
				level: log.LvWarn + 1,
			},
			&action{
				expect: slog.LevelError,
			},
		),
		gen(
			"Error-1",
			&condition{
				level: log.LvError - 1,
			},
			&action{
				expect: slog.LevelError,
			},
		),
		gen(
			"Error",
			&condition{
				level: log.LvError,
			},
			&action{
				expect: slog.LevelError,
			},
		),
		gen(
			"Error+1",
			&condition{
				level: log.LvError + 1,
			},
			&action{
				expect: slog.LevelError,
			},
		),
		gen(
			"Fatal-1",
			&condition{
				level: log.LvFatal - 1,
			},
			&action{
				expect: slog.LevelError,
			},
		),
		gen(
			"Fatal",
			&condition{
				level: log.LvFatal,
			},
			&action{
				expect: slog.LevelError,
			},
		),
		gen(
			"Fatal+1",
			&condition{
				level: log.LvFatal + 1,
			},
			&action{
				expect: slog.LevelError,
			},
		),
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			got := log.LvToSLogLevel(tt.C.level)
			testutil.Diff(t, tt.A.expect, got)
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

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"Debug-1",
			&condition{
				level: slog.LevelDebug - 1,
			},
			&action{
				expect: log.LvDebug,
			},
		),
		gen(
			"Debug",
			&condition{
				level: slog.LevelDebug,
			},
			&action{
				expect: log.LvDebug,
			},
		),
		gen(
			"Debug+1",
			&condition{
				level: slog.LevelDebug + 1,
			},
			&action{
				expect: log.LvInfo,
			},
		),
		gen(
			"Info-1",
			&condition{
				level: slog.LevelInfo - 1,
			},
			&action{
				expect: log.LvInfo,
			},
		),
		gen(
			"Info",
			&condition{
				level: slog.LevelInfo,
			},
			&action{
				expect: log.LvInfo,
			},
		),
		gen(
			"Info+1",
			&condition{
				level: slog.LevelInfo + 1,
			},
			&action{
				expect: log.LvWarn,
			},
		),
		gen(
			"Warn-1",
			&condition{
				level: slog.LevelWarn - 1,
			},
			&action{
				expect: log.LvWarn,
			},
		),
		gen(
			"Warn",
			&condition{
				level: slog.LevelWarn,
			},
			&action{
				expect: log.LvWarn,
			},
		),
		gen(
			"Warn+1",
			&condition{
				level: slog.LevelWarn + 1,
			},
			&action{
				expect: log.LvError,
			},
		),
		gen(
			"Error-1",
			&condition{
				level: slog.LevelError - 1,
			},
			&action{
				expect: log.LvError,
			},
		),
		gen(
			"Error",
			&condition{
				level: slog.LevelError,
			},
			&action{
				expect: log.LvError,
			},
		),
		gen(
			"Error+1",
			&condition{
				level: slog.LevelError + slog.Level(1),
			},
			&action{
				expect: log.LvError,
			},
		),
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			got := log.LvFromSLogLevel(tt.C.level)
			testutil.Diff(t, tt.A.expect, got)
		})
	}
}
