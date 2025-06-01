// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package log_test

import (
	"testing"

	"github.com/aileron-gateway/aileron-gateway/kernel/log"
	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
)

func TestLevelFromText(t *testing.T) {
	type condition struct {
		level string
	}

	type action struct {
		expect log.LogLevel
	}

	tb := testutil.NewTableBuilder[*condition, *action]().Name(t.Name())
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"Trace",
			[]string{},
			[]string{},
			&condition{
				level: log.Trace,
			},
			&action{
				expect: log.LvTrace,
			},
		),
		gen(
			"Debug",
			[]string{},
			[]string{},
			&condition{
				level: log.Debug,
			},
			&action{
				expect: log.LvDebug,
			},
		),
		gen(
			"INFO",
			[]string{},
			[]string{},
			&condition{
				level: log.Info,
			},
			&action{
				expect: log.LvInfo,
			},
		),
		gen(
			"WARN",
			[]string{},
			[]string{},
			&condition{
				level: log.Warn,
			},
			&action{
				expect: log.LvWarn,
			},
		),
		gen(
			"ERROR",
			[]string{},
			[]string{},
			&condition{
				level: log.Error,
			},
			&action{
				expect: log.LvError,
			},
		),
		gen(
			"FATAL",
			[]string{},
			[]string{},
			&condition{
				level: log.Fatal,
			},
			&action{
				expect: log.LvFatal,
			},
		),
		gen(
			"UNKNOWN",
			[]string{},
			[]string{},
			&condition{
				level: "UNKNOWN",
			},
			&action{
				expect: log.LvInfo,
			},
		),
		gen(
			"DeBug",
			[]string{},
			[]string{},
			&condition{
				level: "DeBug",
			},
			&action{
				expect: log.LvDebug,
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			got := log.LevelFromText(tt.C().level)
			testutil.Diff(t, tt.A().expect, got)
		})
	}
}
