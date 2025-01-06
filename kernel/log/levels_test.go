package log_test

import (
	"testing"

	k "github.com/aileron-gateway/aileron-gateway/apis/kernel"
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

func TestLevel(t *testing.T) {
	type condition struct {
		level k.LogLevel
	}

	type action struct {
		expect log.LogLevel
	}

	tb := testutil.NewTableBuilder[*condition, *action]().Name(t.Name())
	cndLvTrace := tb.Condition("Trace", "input Trace as log.LogLevel")
	cndLvDebug := tb.Condition("Debug", "input Debug as log.LogLevel")
	cndLvInfo := tb.Condition("Info", "input Info as log.LogLevel")
	cndLvWarn := tb.Condition("Warn", "input Warn as log.LogLevel")
	cndLvError := tb.Condition("Error", "input Error as log.LogLevel")
	cndLvFatal := tb.Condition("Fatal", "input Fatal as log.LogLevel")
	cndLvOther := tb.Condition("Other", "input other log level")
	actLvTrace := tb.Action("Debug", "check that the returned log.LogLevel is Trace")
	actLvDebug := tb.Action("Debug", "check that the returned log.LogLevel is Debug")
	actLvInfo := tb.Action("Info", "check that the returned log.LogLevel is Info")
	actLvWarn := tb.Action("Warn", "check that the returned log.LogLevel is Warn")
	actLvError := tb.Action("Error", "check that the returned log.LogLevel is Error")
	actLvFatal := tb.Action("Fatal", "check that the returned log.LogLevel is Fatal")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"Trace-1",
			[]string{cndLvOther},
			[]string{actLvTrace},
			&condition{
				level: k.LogLevel_Trace - 1,
			},
			&action{
				expect: log.LvTrace,
			},
		),
		gen(
			"Trace",
			[]string{cndLvTrace},
			[]string{actLvTrace},
			&condition{
				level: k.LogLevel_Trace,
			},
			&action{
				expect: log.LvTrace,
			},
		),
		gen(
			"Debug",
			[]string{cndLvDebug},
			[]string{actLvDebug},
			&condition{
				level: k.LogLevel_Debug,
			},
			&action{
				expect: log.LvDebug,
			},
		),
		gen(
			"Info",
			[]string{cndLvInfo},
			[]string{actLvInfo},
			&condition{
				level: k.LogLevel_Info,
			},
			&action{
				expect: log.LvInfo,
			},
		),
		gen(
			"Warn",
			[]string{cndLvWarn},
			[]string{actLvWarn},
			&condition{
				level: k.LogLevel_Warn,
			},
			&action{
				expect: log.LvWarn,
			},
		),
		gen(
			"Error",
			[]string{cndLvError},
			[]string{actLvError},
			&condition{
				level: k.LogLevel_Error,
			},
			&action{
				expect: log.LvError,
			},
		),
		gen(
			"Fatal",
			[]string{cndLvFatal},
			[]string{actLvFatal},
			&condition{
				level: k.LogLevel_Fatal,
			},
			&action{
				expect: log.LvFatal,
			},
		),
		gen(
			"Fatal+1",
			[]string{cndLvOther},
			[]string{actLvFatal},
			&condition{
				level: k.LogLevel_Fatal + 1,
			},
			&action{
				expect: log.LvFatal,
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			got := log.Level(tt.C().level)
			testutil.Diff(t, tt.A().expect, got)
		})
	}
}
