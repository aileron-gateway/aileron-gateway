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

	CndInputNonZeroValues := "input non-zero values"
	CndInputZeroValues := "input zero values"
	ActCheckExpected := "expected value"

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	tb.Condition(CndInputNonZeroValues, "input non-zero value")
	tb.Condition(CndInputZeroValues, "input zero value")
	tb.Action(ActCheckExpected, "check that an expected values returned")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"zaro values",
			[]string{},
			[]string{ActCheckExpected},
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
			[]string{CndInputNonZeroValues},
			[]string{ActCheckExpected},
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

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			testutil.Diff(t, tt.A().value, tt.C().err.Error())
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

	CndInputNonZeroValues := "input non-zero values"
	CndInputZeroValues := "input zero values"
	ActCheckExpected := "expected value"

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	tb.Condition(CndInputNonZeroValues, "input non-zero value")
	tb.Condition(CndInputZeroValues, "input zero value")
	tb.Action(ActCheckExpected, "check that an expected values returned")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"zaro values",
			[]string{},
			[]string{ActCheckExpected},
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
			[]string{CndInputNonZeroValues},
			[]string{ActCheckExpected},
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

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			testutil.Diff(t, tt.A().value, tt.C().err.Code())
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

	CndInputNonZeroValues := "input non-zero values"
	CndInputZeroValues := "input zero values"
	ActCheckExpected := "expected value"

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	tb.Condition(CndInputNonZeroValues, "input non-zero value")
	tb.Condition(CndInputZeroValues, "input zero value")
	tb.Action(ActCheckExpected, "check that an expected values returned")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"zaro values",
			[]string{},
			[]string{ActCheckExpected},
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
			[]string{CndInputNonZeroValues},
			[]string{ActCheckExpected},
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

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			testutil.Diff(t, tt.A().value, tt.C().err.Kind())
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

	CndInputNonZeroValues := "input non-zero values"
	CndInputZeroValues := "input zero values"
	ActCheckExpected := "expected value"

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	tb.Condition(CndInputNonZeroValues, "input non-zero value")
	tb.Condition(CndInputZeroValues, "input zero value")
	tb.Action(ActCheckExpected, "check that an expected values returned")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"zaro values",
			[]string{},
			[]string{ActCheckExpected},
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
			[]string{CndInputNonZeroValues},
			[]string{ActCheckExpected},
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

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			testutil.Diff(t, tt.A().value, tt.C().err.StackTrace())
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

	CndInputNonZeroValues := "input non-zero values"
	CndInputZeroValues := "input zero values"
	ActCheckExpected := "expected value"

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	tb.Condition(CndInputNonZeroValues, "input non-zero value")
	tb.Condition(CndInputZeroValues, "input zero value")
	tb.Action(ActCheckExpected, "check that an expected values returned")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"zaro values",
			[]string{},
			[]string{ActCheckExpected},
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
			[]string{CndInputNonZeroValues},
			[]string{ActCheckExpected},
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

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			testutil.Diff(t, tt.A().value, tt.C().err.Name())
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

	CndInputNonZeroValues := "input non-zero values"
	CndInputZeroValues := "input zero values"
	ActCheckExpected := "expected value"

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	tb.Condition(CndInputNonZeroValues, "input non-zero value")
	tb.Condition(CndInputZeroValues, "input zero value")
	tb.Action(ActCheckExpected, "check that an expected values returned")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"zaro values",
			[]string{},
			[]string{ActCheckExpected},
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
			[]string{CndInputNonZeroValues},
			[]string{ActCheckExpected},
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

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			testutil.Diff(t, tt.A().value, tt.C().err.Map())
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

	CndInputNonZeroValues := "input non-zero values"
	CndInputZeroValues := "input zero values"
	ActCheckExpected := "expected value"

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	tb.Condition(CndInputNonZeroValues, "input non-zero value")
	tb.Condition(CndInputZeroValues, "input zero value")
	tb.Action(ActCheckExpected, "check that an expected values returned")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"zaro values",
			[]string{},
			[]string{ActCheckExpected},
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
			[]string{CndInputNonZeroValues},
			[]string{ActCheckExpected},
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

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			testutil.Diff(t, tt.A().value, tt.C().err.KeyValues())
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

	CndInputNonZeroValues := "input non-zero values"
	CndInputZeroValues := "input zero values"
	ActCheckExpected := "expected value"

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	tb.Condition(CndInputNonZeroValues, "input non-zero value")
	tb.Condition(CndInputZeroValues, "input zero value")
	tb.Action(ActCheckExpected, "check that an expected values returned")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"zaro values",
			[]string{},
			[]string{ActCheckExpected},
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
			[]string{CndInputNonZeroValues},
			[]string{ActCheckExpected},
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

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			testutil.Diff(t, tt.A().value.Error(), tt.C().err.Unwrap().Error())
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

	CndInputNil := "input nil"
	CndDifferentCodes := "different codes"
	CndNoCoderInterface := "no coder interface"
	ActCheckTrue := "true returned"

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	tb.Condition(CndInputNil, "input nil")
	tb.Condition(CndDifferentCodes, "two codes are different")
	tb.Condition(CndNoCoderInterface, "no coder interface is implemented")
	tb.Action(ActCheckTrue, "check that true is returned")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"nil target",
			[]string{},
			[]string{},
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
			[]string{CndDifferentCodes},
			[]string{},
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
			[]string{CndDifferentCodes},
			[]string{},
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
			[]string{CndDifferentCodes, CndNoCoderInterface},
			[]string{},
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
			[]string{CndDifferentCodes, CndNoCoderInterface},
			[]string{},
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
			[]string{},
			[]string{ActCheckTrue},
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
			[]string{},
			[]string{ActCheckTrue},
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

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			testutil.Diff(t, tt.A().result, tt.C().err.Is(tt.C().is))
		})
	}
}
