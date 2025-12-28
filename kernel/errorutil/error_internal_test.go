// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package errorutil

import (
	"errors"
	"testing"

	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
)

func TestError(t *testing.T) {
	type condition struct {
		err *ErrorAttrs
	}

	type action struct {
		value string
	}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"zaro values",
			&condition{
				err: &ErrorAttrs{
					code:  "dummy",
					kind:  "dummy",
					stack: []byte("dummy"),
					name:  "dummy",
					msg:   "",
					err:   errors.New("dummy"),
				},
			},
			&action{
				value: "",
			},
		),
		gen(
			"non-zaro values",
			&condition{
				err: &ErrorAttrs{
					code:  "dummy",
					kind:  "dummy",
					stack: []byte("dummy"),
					name:  "dummy",
					msg:   "test",
					err:   errors.New("dummy"),
				},
			},
			&action{
				value: "test",
			},
		),
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			testutil.Diff(t, tt.A.value, tt.C.err.Error())
		})
	}
}

func TestCode(t *testing.T) {
	type condition struct {
		err *ErrorAttrs
	}

	type action struct {
		value string
	}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"zaro values",
			&condition{
				err: &ErrorAttrs{
					code:  "",
					kind:  "dummy",
					stack: []byte("dummy"),
					name:  "dummy",
					msg:   "dummy",
					err:   errors.New("dummy"),
				},
			},
			&action{
				value: "",
			},
		),
		gen(
			"non-zaro values",
			&condition{
				err: &ErrorAttrs{
					code:  "test",
					kind:  "dummy",
					stack: []byte("dummy"),
					name:  "dummy",
					msg:   "dummy",
					err:   errors.New("dummy"),
				},
			},
			&action{
				value: "test",
			},
		),
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			testutil.Diff(t, tt.A.value, tt.C.err.Code())
		})
	}
}

func TestKind(t *testing.T) {
	type condition struct {
		err *ErrorAttrs
	}

	type action struct {
		value string
	}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"zaro values",
			&condition{
				err: &ErrorAttrs{
					code:  "dummy",
					kind:  "",
					stack: []byte("dummy"),
					name:  "dummy",
					msg:   "dummy",
					err:   errors.New("dummy"),
				},
			},
			&action{
				value: "",
			},
		),
		gen(
			"non-zaro values",
			&condition{
				err: &ErrorAttrs{
					code:  "dummy",
					kind:  "test",
					stack: []byte("dummy"),
					name:  "dummy",
					msg:   "dummy",
					err:   errors.New("dummy"),
				},
			},
			&action{
				value: "test",
			},
		),
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			testutil.Diff(t, tt.A.value, tt.C.err.Kind())
		})
	}
}

func TestStackTrace(t *testing.T) {
	type condition struct {
		err *ErrorAttrs
	}

	type action struct {
		value string
	}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"zaro values",
			&condition{
				err: &ErrorAttrs{
					code:  "dummy",
					kind:  "dummy",
					stack: []byte(""),
					name:  "dummy",
					msg:   "dummy",
					err:   errors.New("dummy"),
				},
			},
			&action{
				value: "",
			},
		),
		gen(
			"non-zaro values",
			&condition{
				err: &ErrorAttrs{
					code:  "dummy",
					kind:  "dummy",
					stack: []byte("test"),
					name:  "dummy",
					msg:   "dummy",
					err:   errors.New("dummy"),
				},
			},
			&action{
				value: "test",
			},
		),
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			testutil.Diff(t, tt.A.value, tt.C.err.StackTrace())
		})
	}
}

func TestName(t *testing.T) {
	type condition struct {
		err *ErrorAttrs
	}

	type action struct {
		value string
	}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"zaro values",
			&condition{
				err: &ErrorAttrs{
					code:  "dummy",
					kind:  "dummy",
					stack: []byte("dummy"),
					name:  "",
					msg:   "dummy",
					err:   errors.New("dummy"),
				},
			},
			&action{
				value: "",
			},
		),
		gen(
			"non-zaro values",
			&condition{
				err: &ErrorAttrs{
					code:  "dummy",
					kind:  "dummy",
					stack: []byte("dummy"),
					name:  "test",
					msg:   "dummy",
					err:   errors.New("dummy"),
				},
			},
			&action{
				value: "test",
			},
		),
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			testutil.Diff(t, tt.A.value, tt.C.err.Name())
		})
	}
}

func TestMap(t *testing.T) {
	type condition struct {
		err *ErrorAttrs
	}

	type action struct {
		value map[string]any
	}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"zaro values",
			&condition{
				err: &ErrorAttrs{
					code:  "",
					kind:  "",
					stack: []byte(""),
					name:  "dummy",
					msg:   "",
					err:   errors.New("dummy"),
				},
			},
			&action{
				value: map[string]any{
					"code":  "",
					"kind":  "",
					"msg":   "",
					"stack": "",
				},
			},
		),
		gen(
			"non-zaro values",
			&condition{
				err: &ErrorAttrs{
					code:  "code",
					kind:  "kind",
					stack: []byte("stack"),
					name:  "dummy",
					msg:   "msg",
					err:   errors.New("dummy"),
				},
			},
			&action{
				value: map[string]any{
					"code":  "code",
					"kind":  "kind",
					"msg":   "msg",
					"stack": "stack",
				},
			},
		),
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			testutil.Diff(t, tt.A.value, tt.C.err.Map())
		})
	}
}

func TestKeyValues(t *testing.T) {
	type condition struct {
		err *ErrorAttrs
	}

	type action struct {
		value []any
	}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"zaro values",
			&condition{
				err: &ErrorAttrs{
					code:  "",
					kind:  "",
					stack: []byte(""),
					name:  "dummy",
					msg:   "",
					err:   errors.New("dummy"),
				},
			},
			&action{
				value: []any{
					"code", "",
					"kind", "",
					"msg", "",
					"stack", "",
				},
			},
		),
		gen(
			"non-zaro values",
			&condition{
				err: &ErrorAttrs{
					code:  "code",
					kind:  "kind",
					stack: []byte("stack"),
					name:  "dummy",
					msg:   "msg",
					err:   errors.New("dummy"),
				},
			},
			&action{
				value: []any{
					"code", "code",
					"kind", "kind",
					"msg", "msg",
					"stack", "stack",
				},
			},
		),
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			testutil.Diff(t, tt.A.value, tt.C.err.KeyValues())
		})
	}
}

func TestUnwrap(t *testing.T) {
	type condition struct {
		err *ErrorAttrs
	}

	type action struct {
		value error
	}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"zaro values",
			&condition{
				err: &ErrorAttrs{
					code:  "dummy",
					kind:  "dummy",
					stack: []byte("dummy"),
					name:  "dummy",
					msg:   "dummy",
					err:   errors.New(""),
				},
			},
			&action{
				value: errors.New(""),
			},
		),
		gen(
			"non-zaro values",
			&condition{
				err: &ErrorAttrs{
					code:  "dummy",
					kind:  "dummy",
					stack: []byte("dummy"),
					name:  "dummy",
					msg:   "dummy",
					err:   errors.New("test"),
				},
			},
			&action{
				value: errors.New("test"),
			},
		),
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			testutil.Diff(t, tt.A.value.Error(), tt.C.err.Unwrap().Error())
		})
	}
}

func TestIs(t *testing.T) {
	type condition struct {
		err *ErrorAttrs
		is  error
	}

	type action struct {
		result bool
	}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"nil target",
			&condition{
				err: &ErrorAttrs{
					code: "test",
				},
				is: nil,
			},
			&action{
				result: false,
			},
		),
		gen(
			"different codes",
			&condition{
				err: &ErrorAttrs{
					code: "test",
				},
				is: &ErrorAttrs{
					code: "test test",
				},
			},
			&action{
				result: false,
			},
		),
		gen(
			"different inner codes",
			&condition{
				err: &ErrorAttrs{
					code: "test",
				},
				is: &ErrorAttrs{
					code: "test test",
					err: &ErrorAttrs{
						code: "test test test",
					},
				},
			},
			&action{
				result: false,
			},
		),
		gen(
			"no coder error",
			&condition{
				err: &ErrorAttrs{
					code: "test",
				},
				is: errors.New("test"),
			},
			&action{
				result: false,
			},
		),
		gen(
			"no coder for inner error",
			&condition{
				err: &ErrorAttrs{
					code: "test test",
				},
				is: &ErrorAttrs{
					code: "test",
					err:  errors.New("test"),
				},
			},
			&action{
				result: false,
			},
		),
		gen(
			"same codes",
			&condition{
				err: &ErrorAttrs{
					code: "test",
				},
				is: &ErrorAttrs{
					code: "test",
				},
			},
			&action{
				result: true,
			},
		),
		gen(
			"same inner codes",
			&condition{
				err: &ErrorAttrs{
					code: "test",
				},
				is: &ErrorAttrs{
					code: "test test",
					err: &ErrorAttrs{
						code: "test",
					},
				},
			},
			&action{
				result: true,
			},
		),
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			testutil.Diff(t, tt.A.result, tt.C.err.Is(tt.C.is))
		})
	}
}
