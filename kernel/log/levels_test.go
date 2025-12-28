// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package log_test

import (
	"testing"

	"github.com/aileron-gateway/aileron-gateway/internal/testutil"
	"github.com/aileron-gateway/aileron-gateway/kernel/log"
)

func TestLevelFromText(t *testing.T) {
	type condition struct {
		level string
	}

	type action struct {
		expect log.LogLevel
	}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"Trace",
			&condition{
				level: log.Trace,
			},
			&action{
				expect: log.LvTrace,
			},
		),
		gen(
			"Debug",
			&condition{
				level: log.Debug,
			},
			&action{
				expect: log.LvDebug,
			},
		),
		gen(
			"INFO",
			&condition{
				level: log.Info,
			},
			&action{
				expect: log.LvInfo,
			},
		),
		gen(
			"WARN",
			&condition{
				level: log.Warn,
			},
			&action{
				expect: log.LvWarn,
			},
		),
		gen(
			"ERROR",
			&condition{
				level: log.Error,
			},
			&action{
				expect: log.LvError,
			},
		),
		gen(
			"FATAL",
			&condition{
				level: log.Fatal,
			},
			&action{
				expect: log.LvFatal,
			},
		),
		gen(
			"UNKNOWN",
			&condition{
				level: "UNKNOWN",
			},
			&action{
				expect: log.LvInfo,
			},
		),
		gen(
			"DeBug",
			&condition{
				level: "DeBug",
			},
			&action{
				expect: log.LvDebug,
			},
		),
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			got := log.LevelFromText(tt.C.level)
			testutil.Diff(t, tt.A.expect, got)
		})
	}
}
