// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package log

import (
	"strings"
	"testing"
	"time"

	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
)

var (
	_ Attributes = &LocationAttrs{}
	_ Attributes = &DatetimeAttrs{}
	_ Attributes = &CustomAttrs{}
)

func TestNewLocationAttrs(t *testing.T) {
	type condition struct {
		skip int
	}

	type action struct {
		noData bool
		file   string
		fn     string
	}

	ActCheckFile := "check file"
	ActCheckFunc := "check func"
	ActCheckLine := "check line"

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	tb.Action(ActCheckFile, "check that the returned file is the expected one")
	tb.Action(ActCheckFunc, "check that the returned func is the expected one")
	tb.Action(ActCheckLine, "check that the returned line is the expected one")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"skip -10000",
			[]string{},
			[]string{ActCheckFile, ActCheckFunc, ActCheckLine},
			&condition{
				skip: -10000,
			},
			&action{
				file: "runtime/extern.go",
				fn:   "",
			},
		),
		gen(
			"skip -1",
			[]string{},
			[]string{ActCheckFile, ActCheckFunc, ActCheckLine},
			&condition{
				skip: -1,
			},
			&action{
				file: "runtime/extern.go",
				fn:   "",
			},
		),
		gen(
			"skip 0",
			[]string{},
			[]string{ActCheckFile, ActCheckFunc, ActCheckLine},
			&condition{
				skip: 0,
			},
			&action{
				file: "log/attributes.go",
				fn:   "log.NewLocationAttrs",
			},
		),
		gen(
			"skip 1",
			[]string{},
			[]string{ActCheckFile, ActCheckFunc, ActCheckLine},
			&condition{
				skip: 1,
			},
			&action{
				file: "log/attributes_internal_test.go",
				fn:   "log.TestNewLocationAttrs",
			},
		),
		gen(
			"skip 10000",
			[]string{},
			[]string{ActCheckFile, ActCheckFunc, ActCheckLine},
			&condition{
				skip: 10000,
			},
			&action{
				noData: true,
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			loc := NewLocationAttrs(tt.C().skip)
			if tt.A().noData {
				testutil.Diff(t, "", loc.file)
				testutil.Diff(t, "", loc.fn)
				testutil.Diff(t, 0, loc.line)
				return
			}

			testutil.Diff(t, true, strings.Contains(loc.file, tt.A().file))
			testutil.Diff(t, true, strings.Contains(loc.fn, tt.A().fn))
			testutil.Diff(t, true, loc.line > 0)
		})
	}
}

func TestLocationAttrs_Name(t *testing.T) {
	type condition struct {
		loc *LocationAttrs
	}

	type action struct {
		name string
	}

	ActCheckName := "check name"

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	tb.Action(ActCheckName, "check that the returned name is the one expected")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"non empty name",
			[]string{},
			[]string{},
			&condition{
				loc: &LocationAttrs{
					name: "test",
				},
			},
			&action{
				name: "test",
			},
		),
		gen(
			"empty name",
			[]string{},
			[]string{},
			&condition{
				loc: &LocationAttrs{
					name: "",
				},
			},
			&action{
				name: "",
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			testutil.Diff(t, tt.A().name, tt.C().loc.Name())
		})
	}
}

func TestLocationAttrs_Map(t *testing.T) {
	type condition struct {
		loc *LocationAttrs
	}

	type action struct {
		expect map[string]any
	}

	CndWithFile := "with file"
	CndWithFunc := "with func"
	CndWithLine := "with line"
	ActCheckMap := "check map"

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	tb.Condition(CndWithFile, "file has non empty value")
	tb.Condition(CndWithFunc, "func has non empty value")
	tb.Condition(CndWithLine, "line is not 0")
	tb.Action(ActCheckMap, "check that the returned map is the one expected")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"without values",
			[]string{},
			[]string{ActCheckMap},
			&condition{
				loc: &LocationAttrs{
					name: "test",
				},
			},
			&action{
				expect: map[string]any{
					"file": "",
					"line": 0,
					"func": "",
				},
			},
		),
		gen(
			"with file",
			[]string{CndWithFile},
			[]string{ActCheckMap},
			&condition{
				loc: &LocationAttrs{
					name: "test",
					file: "test_file",
				},
			},
			&action{
				expect: map[string]any{
					"file": "test_file",
					"line": 0,
					"func": "",
				},
			},
		),
		gen(
			"with func",
			[]string{CndWithFunc},
			[]string{ActCheckMap},
			&condition{
				loc: &LocationAttrs{
					name: "test",
					fn:   "test_func",
				},
			},
			&action{
				expect: map[string]any{
					"file": "",
					"line": 0,
					"func": "test_func",
				},
			},
		),
		gen(
			"with line",
			[]string{CndWithLine},
			[]string{ActCheckMap},
			&condition{
				loc: &LocationAttrs{
					name: "test",
					line: 10,
				},
			},
			&action{
				expect: map[string]any{
					"file": "",
					"line": 10,
					"func": "",
				},
			},
		),
		gen(
			"with file/func/line",
			[]string{CndWithFile, CndWithFunc, CndWithLine},
			[]string{ActCheckMap},
			&condition{
				loc: &LocationAttrs{
					name: "test",
					file: "test_file",
					fn:   "test_func",
					line: 10,
				},
			},
			&action{
				expect: map[string]any{
					"file": "test_file",
					"line": 10,
					"func": "test_func",
				},
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			testutil.Diff(t, tt.A().expect, tt.C().loc.Map())
		})
	}
}

func TestLocationAttrs_KeyValues(t *testing.T) {
	type condition struct {
		loc *LocationAttrs
	}

	type action struct {
		expect []any
	}

	CndWithFile := "with file"
	CndWithFunc := "with func"
	CndWithLine := "with line"
	ActCheckSlice := "check slice"

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	tb.Condition(CndWithFile, "file has non empty value")
	tb.Condition(CndWithFunc, "func has non empty value")
	tb.Condition(CndWithLine, "line is not 0")
	tb.Action(ActCheckSlice, "check that the returned slice is the one expected")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"without values",
			[]string{},
			[]string{ActCheckSlice},
			&condition{
				loc: &LocationAttrs{
					name: "test",
				},
			},
			&action{
				expect: []any{
					"file", "",
					"line", 0,
					"func", "",
				},
			},
		),
		gen(
			"with file",
			[]string{CndWithFile},
			[]string{ActCheckSlice},
			&condition{
				loc: &LocationAttrs{
					name: "test",
					file: "test_file",
				},
			},
			&action{
				expect: []any{
					"file", "test_file",
					"line", 0,
					"func", "",
				},
			},
		),
		gen(
			"with func",
			[]string{CndWithFunc},
			[]string{ActCheckSlice},
			&condition{
				loc: &LocationAttrs{
					name: "test",
					fn:   "test_func",
				},
			},
			&action{
				expect: []any{
					"file", "",
					"line", 0,
					"func", "test_func",
				},
			},
		),
		gen(
			"with line",
			[]string{CndWithLine},
			[]string{ActCheckSlice},
			&condition{
				loc: &LocationAttrs{
					name: "test",
					line: 10,
				},
			},
			&action{
				expect: []any{
					"file", "",
					"line", 10,
					"func", "",
				},
			},
		),
		gen(
			"with file/func/line",
			[]string{CndWithFile, CndWithFunc, CndWithLine},
			[]string{ActCheckSlice},
			&condition{
				loc: &LocationAttrs{
					name: "test",
					file: "test_file",
					fn:   "test_func",
					line: 10,
				},
			},
			&action{
				expect: []any{
					"file", "test_file",
					"line", 10,
					"func", "test_func",
				},
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			testutil.Diff(t, tt.A().expect, tt.C().loc.KeyValues())
		})
	}
}

func TestNewDatetimeAttrs(t *testing.T) {
	type condition struct {
		dfmt string
		tfmt string
		loc  *time.Location
	}

	type action struct {
		date string
		time string
		zone string
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"valid zone",
			[]string{},
			[]string{},
			&condition{
				dfmt: "2006-01-02",
				tfmt: "15:04:05",
				loc:  time.UTC,
			},
			&action{
				date: "2006-01-02",
				time: "15:04:05",
				zone: "UTC",
			},
		),
		gen(
			"nil zone",
			[]string{},
			[]string{},
			&condition{
				dfmt: "2006-01-02",
				tfmt: "15:04:05",
				loc:  nil,
			},
			&action{
				date: "2006-01-02",
				time: "15:04:05",
				zone: "Local",
			},
		),
		gen(
			"empty dfmt",
			[]string{},
			[]string{},
			&condition{
				dfmt: "",
				tfmt: "15:04:05",
				loc:  time.UTC,
			},
			&action{
				date: "",
				time: "15:04:05",
				zone: "UTC",
			},
		),
		gen(
			"empty tfmt",
			[]string{},
			[]string{},
			&condition{
				dfmt: "2006-01-02",
				tfmt: "",
				loc:  time.UTC,
			},
			&action{
				date: "2006-01-02",
				time: "",
				zone: "UTC",
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			dt := NewDatetimeAttrs(tt.C().dfmt, tt.C().tfmt, tt.C().loc)
			now := time.Now()
			if tt.C().loc != nil {
				now = now.In(tt.C().loc)
			}

			testutil.Diff(t, now.Format(tt.A().date), dt.date)
			testutil.Diff(t, now.Format(tt.A().time), dt.time)
			testutil.Diff(t, now.Location().String(), tt.A().zone)
		})
	}
}

func TestDatetimeAttrs_Name(t *testing.T) {
	type condition struct {
		dt *DatetimeAttrs
	}

	type action struct {
		name string
	}

	ActCheckName := "check name"

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	tb.Action(ActCheckName, "check that the returned name is the one expected")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"non empty name",
			[]string{},
			[]string{},
			&condition{
				dt: &DatetimeAttrs{
					name: "test",
				},
			},
			&action{
				name: "test",
			},
		),
		gen(
			"empty name",
			[]string{},
			[]string{},
			&condition{
				dt: &DatetimeAttrs{
					name: "",
				},
			},
			&action{
				name: "",
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			testutil.Diff(t, tt.A().name, tt.C().dt.Name())
		})
	}
}

func TestDatetimeAttrs_Map(t *testing.T) {
	type condition struct {
		dt *DatetimeAttrs
	}

	type action struct {
		expect map[string]any
	}

	CndWithDate := "with date"
	CndWithTime := "with time"
	CndWithZone := "with zone"
	ActCheckMap := "check map"

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	tb.Condition(CndWithDate, "date has non empty value")
	tb.Condition(CndWithTime, "time has non empty value")
	tb.Condition(CndWithZone, "zone has non empty value")
	tb.Action(ActCheckMap, "check that the returned map is the one expected")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"without values",
			[]string{},
			[]string{ActCheckMap},
			&condition{
				dt: &DatetimeAttrs{
					name: "test",
				},
			},
			&action{
				expect: map[string]any{
					"date": "",
					"time": "",
					"zone": "",
				},
			},
		),
		gen(
			"with date",
			[]string{CndWithDate},
			[]string{ActCheckMap},
			&condition{
				dt: &DatetimeAttrs{
					name: "test",
					date: "test_date",
				},
			},
			&action{
				expect: map[string]any{
					"date": "test_date",
					"time": "",
					"zone": "",
				},
			},
		),
		gen(
			"with time",
			[]string{CndWithTime},
			[]string{ActCheckMap},
			&condition{
				dt: &DatetimeAttrs{
					name: "test",
					time: "test_time",
				},
			},
			&action{
				expect: map[string]any{
					"date": "",
					"time": "test_time",
					"zone": "",
				},
			},
		),
		gen(
			"with zone",
			[]string{CndWithZone},
			[]string{ActCheckMap},
			&condition{
				dt: &DatetimeAttrs{
					name: "test",
					zone: "test_zone",
				},
			},
			&action{
				expect: map[string]any{
					"date": "",
					"time": "",
					"zone": "test_zone",
				},
			},
		),
		gen(
			"with date/time/zone",
			[]string{CndWithDate, CndWithTime, CndWithZone},
			[]string{ActCheckMap},
			&condition{
				dt: &DatetimeAttrs{
					name: "test",
					date: "test_date",
					time: "test_time",
					zone: "test_zone",
				},
			},
			&action{
				expect: map[string]any{
					"date": "test_date",
					"time": "test_time",
					"zone": "test_zone",
				},
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			testutil.Diff(t, tt.A().expect, tt.C().dt.Map())
		})
	}
}

func TestDatetimeAttrs_KeyValues(t *testing.T) {
	type condition struct {
		dt *DatetimeAttrs
	}

	type action struct {
		expect []any
	}

	CndWithDate := "with date"
	CndWithTime := "with time"
	CndWithZone := "with zone"
	ActCheckSlice := "check slice"

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	tb.Condition(CndWithDate, "date has non empty value")
	tb.Condition(CndWithTime, "time has non empty value")
	tb.Condition(CndWithZone, "zone has non empty value")
	tb.Action(ActCheckSlice, "check that the returned slice is the one expected")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"without values",
			[]string{},
			[]string{ActCheckSlice},
			&condition{
				dt: &DatetimeAttrs{
					name: "test",
				},
			},
			&action{
				expect: []any{
					"date", "",
					"time", "",
					"zone", "",
				},
			},
		),
		gen(
			"with date",
			[]string{CndWithDate},
			[]string{ActCheckSlice},
			&condition{
				dt: &DatetimeAttrs{
					name: "test",
					date: "test_date",
				},
			},
			&action{
				expect: []any{
					"date", "test_date",
					"time", "",
					"zone", "",
				},
			},
		),
		gen(
			"with time",
			[]string{CndWithTime},
			[]string{ActCheckSlice},
			&condition{
				dt: &DatetimeAttrs{
					name: "test",
					time: "test_time",
				},
			},
			&action{
				expect: []any{
					"date", "",
					"time", "test_time",
					"zone", "",
				},
			},
		),
		gen(
			"with zone",
			[]string{CndWithZone},
			[]string{ActCheckSlice},
			&condition{
				dt: &DatetimeAttrs{
					name: "test",
					zone: "test_zone",
				},
			},
			&action{
				expect: []any{
					"date", "",
					"time", "",
					"zone", "test_zone",
				},
			},
		),
		gen(
			"with date/time/zone",
			[]string{CndWithDate, CndWithTime, CndWithZone},
			[]string{ActCheckSlice},
			&condition{
				dt: &DatetimeAttrs{
					name: "test",
					date: "test_date",
					time: "test_time",
					zone: "test_zone",
				},
			},
			&action{
				expect: []any{
					"date", "test_date",
					"time", "test_time",
					"zone", "test_zone",
				},
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			testutil.Diff(t, tt.A().expect, tt.C().dt.KeyValues())
		})
	}
}

func TestNewCustomAttrs(t *testing.T) {
	type condition struct {
		name  string
		attrs map[string]any
	}

	type action struct {
		attrs map[string]any
	}

	CndInputNilMap := "input nil map"
	ActCheckName := "check name"
	ActCheckMap := "check map"

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	tb.Condition(CndInputNilMap, "")
	tb.Action(ActCheckName, "check that the returned file is the expected one")
	tb.Action(ActCheckMap, "check that the returned func is the expected one")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"non nil map",
			[]string{},
			[]string{ActCheckName, ActCheckMap},
			&condition{
				name:  "test",
				attrs: map[string]any{"foo": "bar", "hoge": "fuga"},
			},
			&action{
				attrs: map[string]any{"foo": "bar", "hoge": "fuga"},
			},
		),
		gen(
			"nil map",
			[]string{CndInputNilMap},
			[]string{ActCheckName, ActCheckMap},
			&condition{
				name:  "test",
				attrs: nil,
			},
			&action{
				attrs: map[string]any{},
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			cus := NewCustomAttrs(tt.C().name, tt.C().attrs)

			testutil.Diff(t, tt.C().name, cus.name)
			testutil.Diff(t, tt.A().attrs, cus.m)
		})
	}
}

func TestCustomAttrs_Name(t *testing.T) {
	type condition struct {
		ct *CustomAttrs
	}

	type action struct {
		name string
	}

	ActCheckName := "check name"

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	tb.Action(ActCheckName, "check that the returned name is the one expected")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"non empty name",
			[]string{},
			[]string{},
			&condition{
				ct: &CustomAttrs{
					name: "test",
				},
			},
			&action{
				name: "test",
			},
		),
		gen(
			"empty name",
			[]string{},
			[]string{},
			&condition{
				ct: &CustomAttrs{
					name: "",
				},
			},
			&action{
				name: "",
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			testutil.Diff(t, tt.A().name, tt.C().ct.Name())
		})
	}
}

func TestCustomAttrs_Map(t *testing.T) {
	type condition struct {
		ct *CustomAttrs
	}

	type action struct {
		expect map[string]any
	}

	CndWithAttr := "with date"
	ActCheckMap := "check map"

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	tb.Condition(CndWithAttr, "attribute has non empty value")
	tb.Action(ActCheckMap, "check that the returned map is the one expected")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"without attribute",
			[]string{},
			[]string{ActCheckMap},
			&condition{
				ct: &CustomAttrs{
					name: "test",
				},
			},
			&action{
				expect: nil,
			},
		),
		gen(
			"with attribute",
			[]string{CndWithAttr},
			[]string{ActCheckMap},
			&condition{
				ct: &CustomAttrs{
					name: "test",
					m:    map[string]any{"foo": "bar"},
				},
			},
			&action{
				expect: map[string]any{"foo": "bar"},
			},
		),
		gen(
			"with 2 attributes",
			[]string{CndWithAttr},
			[]string{ActCheckMap},
			&condition{
				ct: &CustomAttrs{
					name: "test",
					m:    map[string]any{"foo": "bar", "hoge": "fuga"},
				},
			},
			&action{
				expect: map[string]any{"foo": "bar", "hoge": "fuga"},
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			testutil.Diff(t, tt.A().expect, tt.C().ct.Map())
		})
	}
}

func TestCustomAttrs_KeyValues(t *testing.T) {
	type condition struct {
		ct *CustomAttrs
	}

	type action struct {
		expect []any
	}

	CndWithAttr := "with date"
	ActCheckSlice := "check slice"

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	tb.Condition(CndWithAttr, "attribute has non empty value")
	tb.Action(ActCheckSlice, "check that the returned slice is the one expected")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"without attribute",
			[]string{},
			[]string{ActCheckSlice},
			&condition{
				ct: &CustomAttrs{
					name: "test",
				},
			},
			&action{
				expect: []any{},
			},
		),
		gen(
			"with attribute",
			[]string{CndWithAttr},
			[]string{ActCheckSlice},
			&condition{
				ct: &CustomAttrs{
					name: "test",
					m: map[string]any{
						"foo": "bar",
					},
				},
			},
			&action{
				expect: []any{
					"foo", "bar",
				},
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			testutil.Diff(t, tt.A().expect, tt.C().ct.KeyValues())
		})
	}
}
