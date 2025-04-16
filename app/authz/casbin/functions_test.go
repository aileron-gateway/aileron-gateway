// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package casbin

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestMapValue(t *testing.T) {
	type condition struct {
		a []any
	}

	type action struct {
		v   any
		err any // error or errorutil.Kind
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	condInvalidArgLength := tb.Condition("the argument of authValue is less than 2", "the argument of authValue is less than 2")
	condKeyNotString := tb.Condition("key is not string", "key is not string")
	condKeyNotFound := tb.Condition("key does not exists", "key does not exists")
	condMultipleKeys := tb.Condition("multiple keys", "use multiple keys")
	condInt := tb.Condition("value type is int", "value type is int")
	condInt8 := tb.Condition("value type is int8", "value type is int8")
	condInt16 := tb.Condition("value type is int16", "value type is int16")
	condInt32 := tb.Condition("value type is int32", "value type is int32")
	condInt64 := tb.Condition("value type is int64", "value type is int64")
	condFloat32 := tb.Condition("value type is float32", "value type is float32")
	condFloat64 := tb.Condition("value type is float64", "value type is float64")
	condJWTMapClaims := tb.Condition("first argument type is jwt.MapClaims", "first argument type is jwt.MapClaims")
	actError := tb.Action("error", "check that the expected error is returned")
	actNoError := tb.Action("no error", "check that the there is no error")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"int",
			[]string{condInt},
			[]string{actNoError},
			&condition{
				a: []any{
					map[string]any{"int": 1},
					"int",
				},
			},
			&action{
				v:   float64(1),
				err: nil,
			},
		),
		gen(
			"int8",
			[]string{condInt8},
			[]string{actNoError},
			&condition{
				a: []any{
					map[string]any{"int8": int8(1)},
					"int8",
				},
			},
			&action{
				v:   float64(1),
				err: nil,
			},
		),
		gen(
			"int16",
			[]string{condInt16},
			[]string{actNoError},
			&condition{
				a: []any{
					map[string]any{"int16": int16(1)},
					"int16",
				},
			},
			&action{
				v:   float64(1),
				err: nil,
			},
		),
		gen(
			"int32",
			[]string{condInt32},
			[]string{actNoError},
			&condition{
				a: []any{
					map[string]any{"int32": int32(1)},
					"int32",
				},
			},
			&action{
				v:   float64(1),
				err: nil,
			},
		),
		gen(
			"int64",
			[]string{condInt64},
			[]string{actNoError},
			&condition{
				a: []any{
					map[string]any{"int64": int64(1)},
					"int64",
				},
			},
			&action{
				v:   float64(1),
				err: nil,
			},
		),
		gen(
			"float32",
			[]string{condFloat32},
			[]string{actNoError},
			&condition{
				a: []any{
					map[string]any{"float32": float32(1)},
					"float32",
				},
			},
			&action{
				v:   float64(1),
				err: nil,
			},
		),
		gen(
			"float64",
			[]string{condFloat64},
			[]string{actNoError},
			&condition{
				a: []any{
					map[string]any{"float64": float64(1)},
					"float64",
				},
			},
			&action{
				v:   float64(1),
				err: nil,
			},
		),
		gen(
			"use multiple keys",
			[]string{condMultipleKeys},
			[]string{actNoError},
			&condition{
				a: []any{
					map[string]any{
						"map": map[string]any{"xx": map[string]any{"yy": "YY"}},
					},
					"map", "xx", "yy",
				},
			},
			&action{
				v:   "YY",
				err: nil,
			},
		),
		gen(
			"the argument of authValue is less than 2",
			[]string{condInvalidArgLength},
			[]string{actError},
			&condition{
				a: []any{},
			},
			&action{
				v:   nil,
				err: errAuthValueInvalidArgLength,
			},
		),
		gen(
			"first argument jwt.MapClaims",
			[]string{condInt, condJWTMapClaims},
			[]string{actNoError},
			&condition{
				a: []any{
					jwt.MapClaims(map[string]any{"int": 1, "int8": int8(1), "int16": int16(1)}),
					"int",
				},
			},
			&action{
				v:   float64(1),
				err: nil,
			},
		),
		gen(
			"first argument map[string]any",
			[]string{},
			[]string{actNoError},
			&condition{
				a: []any{
					map[string]any{"int": 1, "int8": int8(1), "int16": int16(1)},
					"int",
				},
			},
			&action{
				v:   float64(1),
				err: nil,
			},
		),
		gen(
			"first argument nil",
			[]string{},
			[]string{actError},
			&condition{
				a: []any{
					nil,
					"int",
				},
			},
			&action{
				v:   nil,
				err: errAuthValueInvalidType,
			},
		),
		gen(
			"invalid first argument",
			[]string{},
			[]string{actError},
			&condition{
				a: []any{
					"invalid",
					"int",
				},
			},
			&action{
				v:   nil,
				err: errAuthValueInvalidType,
			},
		),
		gen(
			"key is not string",
			[]string{condKeyNotString},
			[]string{actError},
			&condition{
				a: []any{
					map[string]any{},
					1,
				},
			},
			&action{
				v:   nil,
				err: errKeyNotString,
			},
		),
		gen(
			"key does not exists",
			[]string{condKeyNotFound},
			[]string{actError},
			&condition{
				a: []any{
					map[string]any{
						"int":     1,
						"int8":    int8(1),
						"int16":   int16(1),
						"int32":   int32(1),
						"int64":   int64(1),
						"float32": float32(1),
						"float64": float64(1),
						"map": map[string]any{
							"xx": map[string]any{
								"yy": "YY",
							},
						},
					},
					"error",
				},
			},
			&action{
				v:   nil,
				err: errKeyNotFound,
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			v, err := mapValue(tt.C().a...)
			testutil.Diff(t, tt.A().err, err, cmpopts.EquateErrors())
			testutil.Diff(t, tt.A().v, v)
		})
	}
}

func TestQueryValue(t *testing.T) {
	type condition struct {
		args []any
	}

	type action struct {
		v   any
		err any // error or errorutil.Kind
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"valid url args/1 value found",
			[]string{},
			[]string{},
			&condition{
				args: []any{
					url.Values{"foo": []string{"bar"}},
					"foo",
				},
			},
			&action{
				v:   "bar",
				err: nil,
			},
		),
		gen(
			"valid url args/2 values found",
			[]string{},
			[]string{},
			&condition{
				args: []any{
					url.Values{"foo": []string{"baz", "bar"}},
					"foo",
				},
			},
			&action{
				v:   "baz",
				err: nil,
			},
		),
		gen(
			"valid map args/1 value found",
			[]string{},
			[]string{},
			&condition{
				args: []any{
					map[string][]string{"foo": {"bar"}},
					"foo",
				},
			},
			&action{
				v:   "bar",
				err: nil,
			},
		),
		gen(
			"valid map args/2 value found",
			[]string{},
			[]string{},
			&condition{
				args: []any{
					map[string][]string{"foo": {"baz", "bar"}},
					"foo",
				},
			},
			&action{
				v:   "baz",
				err: nil,
			},
		),
		gen(
			"1st arg nil",
			[]string{},
			[]string{},
			&condition{
				args: []any{
					nil, // Invalid type.
					"foo",
				},
			},
			&action{
				v:   "",
				err: errQueryValueInvalidType,
			},
		),
		gen(
			"invalid valid 1st arg",
			[]string{},
			[]string{},
			&condition{
				args: []any{
					map[string]int{"foo": 123}, // Invalid type.
					"foo",
				},
			},
			&action{
				v:   "",
				err: errQueryValueInvalidType,
			},
		),
		gen(
			"invalid valid 2nd arg",
			[]string{},
			[]string{},
			&condition{
				args: []any{
					map[string][]string{"foo": {"bar"}},
					123, // Invalid type.
				},
			},
			&action{
				v:   "",
				err: errKeyNotString,
			},
		),
		gen(
			"invalid arg length",
			[]string{},
			[]string{},
			&condition{
				args: []any{
					map[string][]string{"foo": {"bar"}},
					"foo",
					123, // Invalid type.
				},
			},
			&action{
				v:   "",
				err: errQueryValueInvalidArgLength,
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			v, err := queryValue(tt.C().args...)
			testutil.Diff(t, tt.A().err, err, cmpopts.EquateErrors())
			testutil.Diff(t, tt.A().v, v)
		})
	}
}

func TestHeaderValue(t *testing.T) {
	type condition struct {
		args []any
	}

	type action struct {
		v   any
		err any // error or errorutil.Kind
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"valid header args/1 value found",
			[]string{},
			[]string{},
			&condition{
				args: []any{
					http.Header{"Foo": []string{"bar"}},
					"foo", // CanonicalMIMEHeaderKey is used
				},
			},
			&action{
				v:   "bar",
				err: nil,
			},
		),
		gen(
			"valid header args/2 values found",
			[]string{},
			[]string{},
			&condition{
				args: []any{
					http.Header{"Foo": []string{"baz", "bar"}},
					"foo", // CanonicalMIMEHeaderKey is used
				},
			},
			&action{
				v:   "baz",
				err: nil,
			},
		),
		gen(
			"valid map args/1 value found",
			[]string{},
			[]string{},
			&condition{
				args: []any{
					map[string][]string{"foo": {"bar"}},
					"foo",
				},
			},
			&action{
				v:   "bar",
				err: nil,
			},
		),
		gen(
			"valid map args/2 value found",
			[]string{},
			[]string{},
			&condition{
				args: []any{
					map[string][]string{"foo": {"baz", "bar"}},
					"foo",
				},
			},
			&action{
				v:   "baz",
				err: nil,
			},
		),
		gen(
			"first arg nil",
			[]string{},
			[]string{},
			&condition{
				args: []any{
					nil, // Invalid type.
					"foo",
				},
			},
			&action{
				v:   "",
				err: errHeaderValueInvalidType,
			},
		),
		gen(
			"invalid valid 1st arg",
			[]string{},
			[]string{},
			&condition{
				args: []any{
					map[string]int{"foo": 123}, // Invalid type.
					"foo",
				},
			},
			&action{
				v:   "",
				err: errHeaderValueInvalidType,
			},
		),
		gen(
			"invalid valid 2nd arg",
			[]string{},
			[]string{},
			&condition{
				args: []any{
					map[string][]string{"foo": {"bar"}},
					123, // Invalid type.
				},
			},
			&action{
				v:   "",
				err: errKeyNotString,
			},
		),
		gen(
			"invalid arg length",
			[]string{},
			[]string{},
			&condition{
				args: []any{
					map[string][]string{"foo": {"bar"}},
					"foo",
					123, // Invalid type.
				},
			},
			&action{
				v:   "",
				err: errHeaderValueInvalidArgLength,
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			v, err := headerValue(tt.C().args...)
			testutil.Diff(t, tt.A().err, err, cmpopts.EquateErrors())
			testutil.Diff(t, tt.A().v, v)
		})
	}
}

func TestContains(t *testing.T) {
	type condition struct {
		f    func(...any) (any, error)
		args []any
	}

	type action struct {
		v   any
		err any // error or errorutil.Kind
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"contain string value",
			[]string{},
			[]string{},
			&condition{
				f: contains[string],
				args: []any{
					[]string{"a", "b", "c"},
					"a",
				},
			},
			&action{
				v:   true,
				err: nil,
			},
		),
		gen(
			"contain int value",
			[]string{},
			[]string{},
			&condition{
				f: contains[int],
				args: []any{
					[]int{123, 456, 789},
					789,
				},
			},
			&action{
				v:   true,
				err: nil,
			},
		),
		gen(
			"less than 2 args",
			[]string{},
			[]string{},
			&condition{
				f:    contains[string],
				args: []any{},
			},
			&action{
				v:   nil,
				err: errContainsInvalidArgLength,
			},
		),
		gen(
			"greater than 2 args",
			[]string{},
			[]string{},
			&condition{
				f: contains[string],
				args: []any{
					[]string{"a", "b", "c"},
					"a", "b", "c",
				},
			},
			&action{
				v:   nil,
				err: errContainsInvalidArgLength,
			},
		),
		gen(
			"first argument is invalid slice",
			[]string{},
			[]string{},
			&condition{
				f: contains[string],
				args: []any{
					"error", // requires []string
					"a",
				},
			},
			&action{
				v:   nil,
				err: errContainsInvalidSlice,
			},
		),
		gen(
			"first argument is invalid slice",
			[]string{},
			[]string{},
			&condition{
				f: contains[string],
				args: []any{
					[]int{123, 456}, // requires []string
					"a",
				},
			},
			&action{
				v:   nil,
				err: errContainsInvalidSlice,
			},
		),
		gen(
			"second argument is invalid type",
			[]string{},
			[]string{},
			&condition{
				f: contains[string],
				args: []any{
					[]string{"a", "b", "c"}, // requires []string
					123,                     // requires string
				},
			},
			&action{
				v:   nil,
				err: errContainsInvalidValue,
			},
		),
		gen(
			"not found",
			[]string{},
			[]string{},
			&condition{
				f: contains[string],
				args: []any{
					[]string{"a", "b", "c"},
					"x",
				},
			},
			&action{
				v:   false,
				err: nil,
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			v, err := tt.C().f(tt.C().args...)
			testutil.Diff(t, tt.A().err, err, cmpopts.EquateErrors())
			testutil.Diff(t, tt.A().v, v)
		})
	}
}

func TestContainsNumber(t *testing.T) {
	type condition struct {
		f    func(...any) (any, error)
		args []any
	}

	type action struct {
		v   any
		err any // error or errorutil.Kind
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"contain int value",
			[]string{},
			[]string{},
			&condition{
				f: containsNumber[int],
				args: []any{
					[]int{123, 456, 789},
					456,
				},
			},
			&action{
				v:   true,
				err: nil,
			},
		),
		gen(
			"contain float64 value",
			[]string{},
			[]string{},
			&condition{
				f: containsNumber[float64],
				args: []any{
					[]float64{123, 456, 789},
					float64(789),
				},
			},
			&action{
				v:   true,
				err: nil,
			},
		),
		gen(
			"contain float64 value in []int",
			[]string{},
			[]string{},
			&condition{
				f: containsNumber[int],
				args: []any{
					[]int{123, 456, 789},
					float64(789),
				},
			},
			&action{
				v:   true,
				err: nil,
			},
		),
		gen(
			"less than 2 args",
			[]string{},
			[]string{},
			&condition{
				f:    containsNumber[int],
				args: []any{},
			},
			&action{
				v:   nil,
				err: errContainsInvalidArgLength,
			},
		),
		gen(
			"greater than 2 args",
			[]string{},
			[]string{},
			&condition{
				f: containsNumber[int],
				args: []any{
					[]int{123, 456, 789},
					1, 2, 3,
				},
			},
			&action{
				v:   nil,
				err: errContainsInvalidArgLength,
			},
		),
		gen(
			"first argument is invalid slice",
			[]string{},
			[]string{},
			&condition{
				f: containsNumber[int],
				args: []any{
					"error", // requires []int
					"a",
				},
			},
			&action{
				v:   nil,
				err: errContainsInvalidSlice,
			},
		),
		gen(
			"first argument is invalid slice",
			[]string{},
			[]string{},
			&condition{
				f: containsNumber[int],
				args: []any{
					[]string{"a", "b"}, // requires []int
					123,
				},
			},
			&action{
				v:   nil,
				err: errContainsInvalidSlice,
			},
		),
		gen(
			"second argument is invalid type",
			[]string{},
			[]string{},
			&condition{
				f: containsNumber[int],
				args: []any{
					[]int{123, 456, 789},
					"abc", // requires int
				},
			},
			&action{
				v:   nil,
				err: errContainsInvalidValue,
			},
		),
		gen(
			"not found",
			[]string{},
			[]string{},
			&condition{
				f: containsNumber[int],
				args: []any{
					[]int{123, 456, 789},
					999,
				},
			},
			&action{
				v:   false,
				err: nil,
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			v, err := tt.C().f(tt.C().args...)
			testutil.Diff(t, tt.A().err, err, cmpopts.EquateErrors())
			testutil.Diff(t, tt.A().v, v)
		})
	}
}

func TestAsSlice(t *testing.T) {
	type condition struct {
		f    func(...any) (any, error)
		args []any
	}

	type action struct {
		v   any
		err any // error or errorutil.Kind
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"convert zero value",
			[]string{},
			[]string{},
			&condition{
				f:    asSlice[string],
				args: []any{},
			},
			&action{
				v:   []string{},
				err: nil,
			},
		),
		gen(
			"convert single value",
			[]string{},
			[]string{},
			&condition{
				f:    asSlice[string],
				args: []any{"test"},
			},
			&action{
				v:   []string{"test"},
				err: nil,
			},
		),
		gen(
			"convert multiple value",
			[]string{},
			[]string{},
			&condition{
				f:    asSlice[string],
				args: []any{"foo", "bar"},
			},
			&action{
				v:   []string{"foo", "bar"},
				err: nil,
			},
		),
		gen(
			"all incompatible data type",
			[]string{},
			[]string{},
			&condition{
				f:    asSlice[string],
				args: []any{123, 456},
			},
			&action{
				v:   nil,
				err: errAsSliceInvalidType,
			},
		),
		gen(
			"contains incompatible data type",
			[]string{},
			[]string{},
			&condition{
				f:    asSlice[string],
				args: []any{"abc", 456},
			},
			&action{
				v:   nil,
				err: errAsSliceInvalidType,
			},
		),
		gen(
			"int",
			[]string{},
			[]string{},
			&condition{
				f:    asSlice[int],
				args: []any{123, 456},
			},
			&action{
				v:   []int{123, 456},
				err: nil,
			},
		),
		gen(
			"float32",
			[]string{},
			[]string{},
			&condition{
				f:    asSlice[float32],
				args: []any{float32(123), float32(456)},
			},
			&action{
				v:   []float32{123, 456},
				err: nil,
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			v, err := tt.C().f(tt.C().args...)
			testutil.Diff(t, tt.A().err, err, cmpopts.EquateErrors())
			testutil.Diff(t, tt.A().v, v)
		})
	}
}

func TestAsSliceNumber(t *testing.T) {
	type condition struct {
		f    func(...any) (any, error)
		args []any
	}

	type action struct {
		v   any
		err any // error or errorutil.Kind
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"convert zero value",
			[]string{},
			[]string{},
			&condition{
				f:    asSliceNumber[int],
				args: []any{},
			},
			&action{
				v:   []int{},
				err: nil,
			},
		),
		gen(
			"convert single value",
			[]string{},
			[]string{},
			&condition{
				f:    asSliceNumber[int],
				args: []any{123},
			},
			&action{
				v:   []int{123},
				err: nil,
			},
		),
		gen(
			"convert multiple value",
			[]string{},
			[]string{},
			&condition{
				f:    asSliceNumber[int],
				args: []any{123, 456},
			},
			&action{
				v:   []int{123, 456},
				err: nil,
			},
		),
		gen(
			"all incompatible data type",
			[]string{},
			[]string{},
			&condition{
				f:    asSliceNumber[int],
				args: []any{"abc", "def"},
			},
			&action{
				v:   nil,
				err: errAsSliceInvalidType,
			},
		),
		gen(
			"contains incompatible data type",
			[]string{},
			[]string{},
			&condition{
				f:    asSliceNumber[int],
				args: []any{123, "abc"},
			},
			&action{
				v:   nil,
				err: errAsSliceInvalidType,
			},
		),
		gen(
			"contains float64",
			[]string{},
			[]string{},
			&condition{
				f:    asSliceNumber[int],
				args: []any{123, float64(456)},
			},
			&action{
				v:   []int{123, 456},
				err: nil,
			},
		),
		gen(
			"float32",
			[]string{},
			[]string{},
			&condition{
				f:    asSliceNumber[float32],
				args: []any{float32(123), float32(456)},
			},
			&action{
				v:   []float32{123, 456},
				err: nil,
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			v, err := tt.C().f(tt.C().args...)
			testutil.Diff(t, tt.A().err, err, cmpopts.EquateErrors())
			testutil.Diff(t, tt.A().v, v)
		})
	}
}
