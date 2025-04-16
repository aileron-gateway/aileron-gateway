// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package cron

import (
	"testing"
	"time"

	"github.com/aileron-gateway/aileron-gateway/kernel/er"
	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestNormalize(t *testing.T) {
	type condition struct {
		exp string
	}

	type action struct {
		out string
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"@yearly",
			[]string{},
			[]string{},
			&condition{
				exp: "@yearly",
			},
			&action{
				out: "0 0 1 1 *",
			},
		),
		gen(
			"CRON_TZ=UTC @yearly",
			[]string{},
			[]string{},
			&condition{
				exp: "CRON_TZ=UTC @yearly",
			},
			&action{
				out: "TZ=UTC 0 0 1 1 *",
			},
		),
		gen(
			"@YEARLY",
			[]string{},
			[]string{},
			&condition{
				exp: "@YEARLY",
			},
			&action{
				out: "@YEARLY",
			},
		),
		gen(
			"@annually",
			[]string{},
			[]string{},
			&condition{
				exp: "@annually",
			},
			&action{
				out: "0 0 1 1 *",
			},
		),
		gen(
			"CRON_TZ=UTC @annually",
			[]string{},
			[]string{},
			&condition{
				exp: "CRON_TZ=UTC @annually",
			},
			&action{
				out: "TZ=UTC 0 0 1 1 *",
			},
		),
		gen(
			"@ANNUALLY",
			[]string{},
			[]string{},
			&condition{
				exp: "@ANNUALLY",
			},
			&action{
				out: "@ANNUALLY",
			},
		),
		gen(
			"@monthly",
			[]string{},
			[]string{},
			&condition{
				exp: "@monthly",
			},
			&action{
				out: "0 0 1 * *",
			},
		),
		gen(
			"CRON_TZ=UTC @monthly",
			[]string{},
			[]string{},
			&condition{
				exp: "CRON_TZ=UTC @monthly",
			},
			&action{
				out: "TZ=UTC 0 0 1 * *",
			},
		),
		gen(
			"@MONTHLY",
			[]string{},
			[]string{},
			&condition{
				exp: "@MONTHLY",
			},
			&action{
				out: "@MONTHLY",
			},
		),
		gen(
			"@weekly",
			[]string{},
			[]string{},
			&condition{
				exp: "@weekly",
			},
			&action{
				out: "0 0 * * 0",
			},
		),
		gen(
			"CRON_TZ=UTC @weekly",
			[]string{},
			[]string{},
			&condition{
				exp: "CRON_TZ=UTC @weekly",
			},
			&action{
				out: "TZ=UTC 0 0 * * 0",
			},
		),
		gen(
			"@WEEKLY",
			[]string{},
			[]string{},
			&condition{
				exp: "@WEEKLY",
			},
			&action{
				out: "@WEEKLY",
			},
		),
		gen(
			"@daily",
			[]string{},
			[]string{},
			&condition{
				exp: "@daily",
			},
			&action{
				out: "0 0 * * *",
			},
		),
		gen(
			"CRON_TZ=UTC @weekly",
			[]string{},
			[]string{},
			&condition{
				exp: "CRON_TZ=UTC @daily",
			},
			&action{
				out: "TZ=UTC 0 0 * * *",
			},
		),
		gen(
			"@DAILY",
			[]string{},
			[]string{},
			&condition{
				exp: "@DAILY",
			},
			&action{
				out: "@DAILY",
			},
		),
		gen(
			"@hourly",
			[]string{},
			[]string{},
			&condition{
				exp: "@hourly",
			},
			&action{
				out: "0 * * * *",
			},
		),
		gen(
			"CRON_TZ=UTC @hourly",
			[]string{},
			[]string{},
			&condition{
				exp: "CRON_TZ=UTC @hourly",
			},
			&action{
				out: "TZ=UTC 0 * * * *",
			},
		),
		gen(
			"@HOURLY",
			[]string{},
			[]string{},
			&condition{
				exp: "@HOURLY",
			},
			&action{
				out: "@HOURLY",
			},
		),
		gen(
			"@sunday",
			[]string{},
			[]string{},
			&condition{
				exp: "@sunday",
			},
			&action{
				out: "0 0 * * 0",
			},
		),
		gen(
			"@monday",
			[]string{},
			[]string{},
			&condition{
				exp: "@monday",
			},
			&action{
				out: "0 0 * * 1",
			},
		),
		gen(
			"@tuesday",
			[]string{},
			[]string{},
			&condition{
				exp: "@tuesday",
			},
			&action{
				out: "0 0 * * 2",
			},
		),
		gen(
			"@wednesday",
			[]string{},
			[]string{},
			&condition{
				exp: "@wednesday",
			},
			&action{
				out: "0 0 * * 3",
			},
		),
		gen(
			"@thursday",
			[]string{},
			[]string{},
			&condition{
				exp: "@thursday",
			},
			&action{
				out: "0 0 * * 4",
			},
		),
		gen(
			"@friday",
			[]string{},
			[]string{},
			&condition{
				exp: "@friday",
			},
			&action{
				out: "0 0 * * 5",
			},
		),
		gen(
			"@saturday",
			[]string{},
			[]string{},
			&condition{
				exp: "@saturday",
			},
			&action{
				out: "0 0 * * 6",
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			out := normalize(tt.C().exp)
			testutil.Diff(t, tt.A().out, out)
		})
	}
}

func TestReplaceMonth(t *testing.T) {
	type condition struct {
		exp string
	}

	type action struct {
		out string
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"comma joined JAN to DEC",
			[]string{},
			[]string{},
			&condition{
				exp: "JAN,FEB,MAR,APR,MAY,JUN,JUL,AUG,SEP,OCT,NOV,DEC",
			},
			&action{
				out: "1,2,3,4,5,6,7,8,9,10,11,12",
			},
		),
		gen(
			"comma joined jan to dec",
			[]string{},
			[]string{},
			&condition{
				exp: "jan,feb,mar,apr,may,jun,jul,aug,sep,oct,nov,dec",
			},
			&action{
				out: "jan,feb,mar,apr,may,jun,jul,aug,sep,oct,nov,dec",
			},
		),
		gen(
			"hyphen joined JAN to DEC",
			[]string{},
			[]string{},
			&condition{
				exp: "JAN-DEC",
			},
			&action{
				out: "1-12",
			},
		),
		gen(
			"hyphen joined jan to dec",
			[]string{},
			[]string{},
			&condition{
				exp: "jan-dec",
			},
			&action{
				out: "jan-dec",
			},
		),
		gen(
			"invalid join",
			[]string{},
			[]string{},
			&condition{
				exp: "JANFEB", // Should not be "12"
			},
			&action{
				out: "JANFEB",
			},
		),
		gen(
			"not exist month",
			[]string{},
			[]string{},
			&condition{
				exp: "FOO",
			},
			&action{
				out: "FOO",
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			out := replaceMonth(tt.C().exp)
			testutil.Diff(t, tt.A().out, out)
		})
	}
}

func TestReplaceWeek(t *testing.T) {
	type condition struct {
		exp string
	}

	type action struct {
		out string
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"comma joined SUN to STA",
			[]string{},
			[]string{},
			&condition{
				exp: "SUN,MON,TUE,WED,THU,FRI,SAT",
			},
			&action{
				out: "0,1,2,3,4,5,6",
			},
		),
		gen(
			"comma joined sun to sat",
			[]string{},
			[]string{},
			&condition{
				exp: "sun,mon,tue,wed,thu,fri,sat",
			},
			&action{
				out: "sun,mon,tue,wed,thu,fri,sat",
			},
		),
		gen(
			"hyphen joined sun to sat",
			[]string{},
			[]string{},
			&condition{
				exp: "SUN-SAT",
			},
			&action{
				out: "0-6",
			},
		),
		gen(
			"hyphen joined sun to sat",
			[]string{},
			[]string{},
			&condition{
				exp: "sun-sat",
			},
			&action{
				out: "sun-sat",
			},
		),
		gen(
			"invalid join",
			[]string{},
			[]string{},
			&condition{
				exp: "SUNMON", // Should not be "01"
			},
			&action{
				out: "SUNMON",
			},
		),
		gen(
			"not exist week",
			[]string{},
			[]string{},
			&condition{
				exp: "FOO",
			},
			&action{
				out: "FOO",
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			out := replaceWeek(tt.C().exp)
			testutil.Diff(t, tt.A().out, out)
		})
	}
}

func TestParseValue(t *testing.T) {
	type condition struct {
		exp string
		min int
		max int
	}

	type action struct {
		val uint64
		ok  bool
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"* 0to9",
			[]string{},
			[]string{},
			&condition{
				exp: "*",
				min: 0,
				max: 9,
			},
			&action{
				val: 0b_11111_11111,
				ok:  true,
			},
		),
		gen(
			"* 3to9",
			[]string{},
			[]string{},
			&condition{
				exp: "*",
				min: 3,
				max: 9,
			},
			&action{
				val: 0b_11111_11000,
				ok:  true,
			},
		),
		gen(
			"5 0to9",
			[]string{},
			[]string{},
			&condition{
				exp: "5",
				min: 0,
				max: 9,
			},
			&action{
				val: 0b_1_00000,
				ok:  true,
			},
		),
		gen(
			"5 3to9",
			[]string{},
			[]string{},
			&condition{
				exp: "5",
				min: 3,
				max: 9,
			},
			&action{
				val: 0b_1_00000,
				ok:  true,
			},
		),
		gen(
			"5 6to9",
			[]string{},
			[]string{},
			&condition{
				exp: "5",
				min: 6,
				max: 9,
			},
			&action{
				val: 0,
				ok:  false,
			},
		),
		gen(
			"5-7 0to9",
			[]string{},
			[]string{},
			&condition{
				exp: "5-7",
				min: 0,
				max: 9,
			},
			&action{
				val: 0b_111_00000,
				ok:  true,
			},
		),
		gen(
			"5-7 5to9",
			[]string{},
			[]string{},
			&condition{
				exp: "5-7",
				min: 5,
				max: 9,
			},
			&action{
				val: 0b_111_00000,
				ok:  true,
			},
		),
		gen(
			"5-7 6to9",
			[]string{},
			[]string{},
			&condition{
				exp: "5-7",
				min: 6,
				max: 9,
			},
			&action{
				val: 0,
				ok:  false,
			},
		),
		gen(
			"5-7 4to6",
			[]string{},
			[]string{},
			&condition{
				exp: "5-7",
				min: 4,
				max: 6,
			},
			&action{
				val: 0,
				ok:  false,
			},
		),
		gen(
			"*/2 0to9",
			[]string{},
			[]string{},
			&condition{
				exp: "*/2",
				min: 0,
				max: 9,
			},
			&action{
				val: 0b_01010_10101,
				ok:  true,
			},
		),
		gen(
			"1/2 0to9",
			[]string{},
			[]string{},
			&condition{
				exp: "1/2",
				min: 0,
				max: 9,
			},
			&action{
				val: 0b_10101_01010,
				ok:  true,
			},
		),
		gen(
			"1/3 0to9",
			[]string{},
			[]string{},
			&condition{
				exp: "1/3",
				min: 0,
				max: 9,
			},
			&action{
				val: 0b_00100_10010,
				ok:  true,
			},
		),
		gen(
			"1-6/2 0to9",
			[]string{},
			[]string{},
			&condition{
				exp: "1-6/2",
				min: 0,
				max: 9,
			},
			&action{
				val: 0b_00001_01010,
				ok:  true,
			},
		),
		gen(
			"x-5 0to9",
			[]string{},
			[]string{},
			&condition{
				exp: "x-5",
				min: 0,
				max: 9,
			},
			&action{
				val: 0,
				ok:  false,
			},
		),
		gen(
			"1-x 0to9",
			[]string{},
			[]string{},
			&condition{
				exp: "1-x",
				min: 0,
				max: 9,
			},
			&action{
				val: 0,
				ok:  false,
			},
		),
		gen(
			"1-5/x 0to9",
			[]string{},
			[]string{},
			&condition{
				exp: "1-5/x",
				min: 0,
				max: 9,
			},
			&action{
				val: 0,
				ok:  false,
			},
		),
		gen(
			"x-5/2 0to9",
			[]string{},
			[]string{},
			&condition{
				exp: "x-5/2",
				min: 0,
				max: 9,
			},
			&action{
				val: 0,
				ok:  false,
			},
		),
		gen(
			"1-x/2 0to9",
			[]string{},
			[]string{},
			&condition{
				exp: "1-x/2",
				min: 0,
				max: 9,
			},
			&action{
				val: 0,
				ok:  false,
			},
		),
		gen(
			"5/0 0to9",
			[]string{},
			[]string{},
			&condition{
				exp: "5/0",
				min: 0,
				max: 9,
			},
			&action{
				val: 0,
				ok:  false,
			},
		),
		gen(
			"5/1/2 0to9",
			[]string{},
			[]string{},
			&condition{
				exp: "5/1/2",
				min: 0,
				max: 9,
			},
			&action{
				val: 0,
				ok:  false,
			},
		),
		gen(
			"min>max",
			[]string{},
			[]string{},
			&condition{
				exp: "*",
				min: 2,
				max: 1,
			},
			&action{
				val: 0,
				ok:  false,
			},
		),
		gen(
			"1,2,5 0to9",
			[]string{},
			[]string{},
			&condition{
				exp: "1,2,5",
				min: 0,
				max: 9,
			},
			&action{
				val: 0b_00001_00110,
				ok:  true,
			},
		),
		gen(
			"1,2,5,* 0to9",
			[]string{},
			[]string{},
			&condition{
				exp: "1,2,5,*",
				min: 0,
				max: 9,
			},
			&action{
				val: 0b_11111_11111,
				ok:  true,
			},
		),
		gen(
			"1,2,5-7 0to9",
			[]string{},
			[]string{},
			&condition{
				exp: "1,2,5-7",
				min: 0,
				max: 9,
			},
			&action{
				val: 0b_00111_00110,
				ok:  true,
			},
		),
		gen(
			"*/2,5-7 0to9",
			[]string{},
			[]string{},
			&condition{
				exp: "*/2,5-7",
				min: 0,
				max: 9,
			},
			&action{
				val: 0b_01111_10101,
				ok:  true,
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			val, ok := parseValue(tt.C().exp, tt.C().min, tt.C().max)
			testutil.Diff(t, tt.A().val, val)
			testutil.Diff(t, tt.A().ok, ok)
		})
	}
}

func TestParseRange(t *testing.T) {
	type condition struct {
		exp string
		min int
		max int
	}

	type action struct {
		min int
		max int
		ok  bool
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"wildcard",
			[]string{},
			[]string{},
			&condition{
				exp: "*",
				min: 1,
				max: 10,
			},
			&action{
				min: 1,
				max: 10,
				ok:  true,
			},
		),
		gen(
			"invalid wildcard",
			[]string{},
			[]string{},
			&condition{
				exp: "**",
				min: 1,
				max: 10,
			},
			&action{
				min: 0,
				max: 0,
				ok:  false,
			},
		),
		gen(
			"valid number/in min,max",
			[]string{},
			[]string{},
			&condition{
				exp: "5",
				min: 1,
				max: 10,
			},
			&action{
				min: 5,
				max: 5,
				ok:  true,
			},
		),
		gen(
			"valid number/out of min,max",
			[]string{},
			[]string{},
			&condition{
				exp: "5",
				min: 10,
				max: 20,
			},
			&action{
				min: 0,
				max: 0,
				ok:  false,
			},
		),
		gen(
			"invalid number",
			[]string{},
			[]string{},
			&condition{
				exp: "x5",
				min: 1,
				max: 10,
			},
			&action{
				min: 0,
				max: 0,
				ok:  false,
			},
		),
		gen(
			"valid number/in min,max",
			[]string{},
			[]string{},
			&condition{
				exp: "5",
				min: 1,
				max: 10,
			},
			&action{
				min: 5,
				max: 5,
				ok:  true,
			},
		),
		gen(
			"valid range/in min,max",
			[]string{},
			[]string{},
			&condition{
				exp: "2-8",
				min: 1,
				max: 10,
			},
			&action{
				min: 2,
				max: 8,
				ok:  true,
			},
		),
		gen(
			"valid range/ini out of min,max",
			[]string{},
			[]string{},
			&condition{
				exp: "2-8",
				min: 5,
				max: 10,
			},
			&action{
				min: 0,
				max: 0,
				ok:  false,
			},
		),
		gen(
			"valid range/end out of min,max",
			[]string{},
			[]string{},
			&condition{
				exp: "2-8",
				min: 1,
				max: 5,
			},
			&action{
				min: 0,
				max: 0,
				ok:  false,
			},
		),
		gen(
			"invalid range/ini>end",
			[]string{},
			[]string{},
			&condition{
				exp: "8-2",
				min: 1,
				max: 10,
			},
			&action{
				min: 0,
				max: 0,
				ok:  false,
			},
		),
		gen(
			"invalid range",
			[]string{},
			[]string{},
			&condition{
				exp: "x2-8",
				min: 1,
				max: 10,
			},
			&action{
				min: 0,
				max: 0,
				ok:  false,
			},
		),
		gen(
			"invalid range expression/many hyphen",
			[]string{},
			[]string{},
			&condition{
				exp: "2-5-8",
				min: 1,
				max: 10,
			},
			&action{
				min: 0,
				max: 0,
				ok:  false,
			},
		),
		gen(
			"invalid range expression/no end",
			[]string{},
			[]string{},
			&condition{
				exp: "2-",
				min: 1,
				max: 10,
			},
			&action{
				min: 0,
				max: 0,
				ok:  false,
			},
		),
		gen(
			"invalid min,max",
			[]string{},
			[]string{},
			&condition{
				exp: "*",
				min: 10,
				max: 1,
			},
			&action{
				min: 0,
				max: 0,
				ok:  false,
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			min, max, ok := parseRange(tt.C().exp, tt.C().min, tt.C().max)
			testutil.Diff(t, tt.A().min, min)
			testutil.Diff(t, tt.A().max, max)
			testutil.Diff(t, tt.A().ok, ok)
		})
	}
}

func TestParse(t *testing.T) {
	type condition struct {
		exp string
	}

	type action struct {
		c   *Crontab
		err error
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"* * * * *",
			[]string{},
			[]string{},
			&condition{
				exp: "* * * * *",
			},
			&action{
				c: &Crontab{
					second: 0b_00001, // 0
					minute: 0b_11111_11111_11111_11111_11111_11111_11111_11111_11111_11111_11111_11111,
					hour:   0b_01111_11111_11111_11111_11111,
					day:    0b_00011_11111_11111_11111_11111_11111_11110,
					month:  0b_00111_11111_11110,
					week:   0b_00011_11111,
					loc:    time.Local,
					timer:  time.Now,
				},
			},
		),
		gen(
			"* * * * * *",
			[]string{},
			[]string{},
			&condition{
				exp: "* * * * * *",
			},
			&action{
				c: &Crontab{
					second: 0b_11111_11111_11111_11111_11111_11111_11111_11111_11111_11111_11111_11111,
					minute: 0b_11111_11111_11111_11111_11111_11111_11111_11111_11111_11111_11111_11111,
					hour:   0b_01111_11111_11111_11111_11111,
					day:    0b_00011_11111_11111_11111_11111_11111_11110,
					month:  0b_00111_11111_11110,
					week:   0b_00011_11111,
					loc:    time.Local,
					timer:  time.Now,
				},
			},
		),
		gen(
			"TZ=UTC * * * * *",
			[]string{},
			[]string{},
			&condition{
				exp: "TZ=UTC * * * * *",
			},
			&action{
				c: &Crontab{
					second: 0b_00001, // 0
					minute: 0b_11111_11111_11111_11111_11111_11111_11111_11111_11111_11111_11111_11111,
					hour:   0b_01111_11111_11111_11111_11111,
					day:    0b_00011_11111_11111_11111_11111_11111_11110,
					month:  0b_00111_11111_11110,
					week:   0b_00011_11111,
					loc:    time.UTC,
					timer:  time.Now,
				},
			},
		),
		gen(
			"TZ=UTC * * * * * *",
			[]string{},
			[]string{},
			&condition{
				exp: "TZ=UTC * * * * * *",
			},
			&action{
				c: &Crontab{
					second: 0b_11111_11111_11111_11111_11111_11111_11111_11111_11111_11111_11111_11111,
					minute: 0b_11111_11111_11111_11111_11111_11111_11111_11111_11111_11111_11111_11111,
					hour:   0b_01111_11111_11111_11111_11111,
					day:    0b_00011_11111_11111_11111_11111_11111_11110,
					month:  0b_00111_11111_11110,
					week:   0b_00011_11111,
					loc:    time.UTC,
					timer:  time.Now,
				},
			},
		),
		gen(
			"TZ=Local * * * * *",
			[]string{},
			[]string{},
			&condition{
				exp: "TZ=Local * * * * *",
			},
			&action{
				c: &Crontab{
					second: 0b_00001, // 0
					minute: 0b_11111_11111_11111_11111_11111_11111_11111_11111_11111_11111_11111_11111,
					hour:   0b_01111_11111_11111_11111_11111,
					day:    0b_00011_11111_11111_11111_11111_11111_11110,
					month:  0b_00111_11111_11110,
					week:   0b_00011_11111,
					loc:    time.Local,
					timer:  time.Now,
				},
			},
		),
		gen(
			"TZ=Local * * * * * *",
			[]string{},
			[]string{},
			&condition{
				exp: "TZ=Local * * * * * *",
			},
			&action{
				c: &Crontab{
					second: 0b_11111_11111_11111_11111_11111_11111_11111_11111_11111_11111_11111_11111,
					minute: 0b_11111_11111_11111_11111_11111_11111_11111_11111_11111_11111_11111_11111,
					hour:   0b_01111_11111_11111_11111_11111,
					day:    0b_00011_11111_11111_11111_11111_11111_11110,
					month:  0b_00111_11111_11110,
					week:   0b_00011_11111,
					loc:    time.Local,
					timer:  time.Now,
				},
			},
		),
		gen(
			"CRON_TZ=UTC * * * * *",
			[]string{},
			[]string{},
			&condition{
				exp: "CRON_TZ=UTC * * * * *",
			},
			&action{
				c: &Crontab{
					second: 0b_00001, // 0
					minute: 0b_11111_11111_11111_11111_11111_11111_11111_11111_11111_11111_11111_11111,
					hour:   0b_01111_11111_11111_11111_11111,
					day:    0b_00011_11111_11111_11111_11111_11111_11110,
					month:  0b_00111_11111_11110,
					week:   0b_00011_11111,
					loc:    time.UTC,
					timer:  time.Now,
				},
			},
		),
		gen(
			"CRON_TZ=UTC * * * * * *",
			[]string{},
			[]string{},
			&condition{
				exp: "CRON_TZ=UTC * * * * * *",
			},
			&action{
				c: &Crontab{
					second: 0b_11111_11111_11111_11111_11111_11111_11111_11111_11111_11111_11111_11111,
					minute: 0b_11111_11111_11111_11111_11111_11111_11111_11111_11111_11111_11111_11111,
					hour:   0b_01111_11111_11111_11111_11111,
					day:    0b_00011_11111_11111_11111_11111_11111_11110,
					month:  0b_00111_11111_11110,
					week:   0b_00011_11111,
					loc:    time.UTC,
					timer:  time.Now,
				},
			},
		),
		gen(
			"0 0 1 1 *",
			[]string{},
			[]string{},
			&condition{
				exp: "0 0 1 1 *",
			},
			&action{
				c: &Crontab{
					second: 0b_00001, // 0
					minute: 0b_00001,
					hour:   0b_00001,
					day:    0b_00010,
					month:  0b_00010,
					week:   0b_00011_11111,
					loc:    time.Local,
					timer:  time.Now,
				},
			},
		),
		gen(
			"invalid sec",
			[]string{},
			[]string{},
			&condition{
				exp: "x * * * * *",
			},
			&action{
				c: nil,
				err: &er.Error{
					Package:     ErrPkg,
					Type:        ErrTypeParse,
					Description: ErrDscParse,
				},
			},
		),
		gen(
			"invalid min",
			[]string{},
			[]string{},
			&condition{
				exp: "* x * * * *",
			},
			&action{
				c: nil,
				err: &er.Error{
					Package:     ErrPkg,
					Type:        ErrTypeParse,
					Description: ErrDscParse,
				},
			},
		),
		gen(
			"invalid hour",
			[]string{},
			[]string{},
			&condition{
				exp: "* * x * * *",
			},
			&action{
				c: nil,
				err: &er.Error{
					Package:     ErrPkg,
					Type:        ErrTypeParse,
					Description: ErrDscParse,
				},
			},
		),
		gen(
			"invalid day of month",
			[]string{},
			[]string{},
			&condition{
				exp: "* * * x * *",
			},
			&action{
				c: nil,
				err: &er.Error{
					Package:     ErrPkg,
					Type:        ErrTypeParse,
					Description: ErrDscParse,
				},
			},
		),
		gen(
			"invalid month",
			[]string{},
			[]string{},
			&condition{
				exp: "* * * * x *",
			},
			&action{
				c: nil,
				err: &er.Error{
					Package:     ErrPkg,
					Type:        ErrTypeParse,
					Description: ErrDscParse,
				},
			},
		),
		gen(
			"invalid day of week",
			[]string{},
			[]string{},
			&condition{
				exp: "* * * * * x",
			},
			&action{
				c: nil,
				err: &er.Error{
					Package:     ErrPkg,
					Type:        ErrTypeParse,
					Description: ErrDscParse,
				},
			},
		),
		gen(
			"too many fields",
			[]string{},
			[]string{},
			&condition{
				exp: "* * * * * * *",
			},
			&action{
				c: nil,
				err: &er.Error{
					Package:     ErrPkg,
					Type:        ErrTypeParse,
					Description: ErrDscParse,
				},
			},
		),
		gen(
			"invalid timezone",
			[]string{},
			[]string{},
			&condition{
				exp: "TZ=FooBar * * * * *",
			},
			&action{
				c: nil,
				err: &er.Error{
					Package:     ErrPkg,
					Type:        ErrTypeParse,
					Description: ErrDscParse,
				},
			},
		),
		gen(
			"unschedulable",
			[]string{},
			[]string{},
			&condition{
				exp: "* * 30 2 *", // Feb 30th
			},
			&action{
				c: nil,
				err: &er.Error{
					Package:     ErrPkg,
					Type:        ErrTypeParse,
					Description: ErrDscParse,
				},
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			c, err := Parse(tt.C().exp)
			testutil.Diff(t, tt.A().err, err, cmpopts.EquateErrors())

			opts := []cmp.Option{
				cmp.AllowUnexported(Crontab{}),
				cmp.Comparer(testutil.ComparePointer[*time.Location]),
				cmp.Comparer(testutil.ComparePointer[func() time.Time]),
			}
			testutil.Diff(t, tt.A().c, c, opts...)
		})
	}
}
