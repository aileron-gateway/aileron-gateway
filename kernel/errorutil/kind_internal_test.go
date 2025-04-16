// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package errorutil

import (
	"bytes"
	"errors"
	"io"
	"testing"

	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

var (
	_ = ErrorKind(&Kind{}) // Check that the struct satisfies interface.
	_ = Creator(&Kind{})   // Check that the struct satisfies interface.
)

func TestNewKind(t *testing.T) {
	type condition struct {
		code   string
		kind   string
		tpl    string
		panics bool
	}

	type action struct {
		errMsg string
		kind   *Kind
	}

	CndValidTemplate := "input valid template"
	actCheckNil := "check that the returned kind is nil"
	actCheckPanic := "check that an panic happen"

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	tb.Condition(CndValidTemplate, "input valid fasttemplate as an argument")
	tb.Action(actCheckNil, "check that the returned kind is nil")
	tb.Action(actCheckPanic, "check that a panic will happen")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"valid template",
			[]string{CndValidTemplate},
			[]string{},
			&condition{
				code: "test-code",
				kind: "test-kind",
				tpl:  "this is {{test}}",
			},
			&action{
				kind: &Kind{
					code: "test-code",
					kind: "test-kind",
					tpl:  newTemplate("this is {{test}}"),
				},
			},
		),
		gen(
			"template is empty",
			[]string{CndValidTemplate},
			[]string{},
			&condition{
				code: "test-code",
				kind: "test-kind",
				tpl:  "",
			},
			&action{
				kind: &Kind{
					code: "test-code",
					kind: "test-kind",
					tpl:  newTemplate(""),
				},
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			if tt.C().panics {
				defer func() {
					e := recover().(error)
					testutil.Diff(t, tt.A().errMsg, e.Error())
				}()
			}

			k := NewKind(tt.C().code, tt.C().kind, tt.C().tpl)
			testutil.Diff(t, tt.A().kind, k, cmp.AllowUnexported(Kind{}, template{}))
		})
	}
}

func TestKind_Code(t *testing.T) {
	type condition struct {
		kind *Kind
	}

	type action struct {
		code string
	}

	cndNonEmptyCode := "non-empty code"
	actCheckReturned := "check returned code"

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	tb.Condition(cndNonEmptyCode, "input non-empty string as code")
	tb.Action(actCheckReturned, "check that the returned code is the one expected")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"non-empty code",
			[]string{cndNonEmptyCode},
			[]string{actCheckReturned},
			&condition{
				kind: &Kind{
					code: "test-code",
				},
			},
			&action{
				code: "test-code",
			},
		),
		gen(
			"empty string as code",
			[]string{},
			[]string{actCheckReturned},
			&condition{
				kind: &Kind{
					code: "",
				},
			},
			&action{
				code: "",
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			code := tt.C().kind.Code()
			testutil.Diff(t, tt.A().code, code)
		})
	}
}

func TestKind_Kind(t *testing.T) {
	type condition struct {
		kind *Kind
	}

	type action struct {
		kind string
	}

	cndNonEmptyKind := "non-empty kind"
	actCheckReturned := "check returned kind"

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	tb.Condition(cndNonEmptyKind, "input non-empty string as kind")
	tb.Action(actCheckReturned, "check that the returned kind is the one expected")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"non-empty kind",
			[]string{cndNonEmptyKind},
			[]string{actCheckReturned},
			&condition{
				kind: &Kind{
					kind: "test-kind",
				},
			},
			&action{
				kind: "test-kind",
			},
		),
		gen(
			"empty string as kind",
			[]string{},
			[]string{actCheckReturned},
			&condition{
				kind: &Kind{
					kind: "",
				},
			},
			&action{
				kind: "",
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			kind := tt.C().kind.Kind()
			testutil.Diff(t, tt.A().kind, kind)
		})
	}
}

func TestKind_WithoutStack(t *testing.T) {
	type condition struct {
		kind *Kind
		err  error
	}

	type action struct {
		attrs *ErrorAttrs
	}

	cndInputError := "pass an error"
	cndInputErrorAttrs := "pass an ErrorAttrs"
	actCheckReturned := "check returned ErrorAttrs"

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	tb.Condition(cndInputError, "pass an error which is not instance or ErrorAttrs as argument")
	tb.Condition(cndInputErrorAttrs, "pass an error of ErrorAttrs as argument")
	tb.Action(actCheckReturned, "check that the returned ErrorAttrs is the one with expected values")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"with no error",
			[]string{},
			[]string{actCheckReturned},
			&condition{
				kind: &Kind{
					kind: "test-kind",
					code: "test-code",
					tpl:  newTemplate("test-message"),
				},
			},
			&action{
				attrs: &ErrorAttrs{
					kind: "test-kind",
					code: "test-code",
					name: "error",
					msg:  "test-message",
				},
			},
		),
		gen(
			"with an error",
			[]string{cndInputError},
			[]string{actCheckReturned},
			&condition{
				kind: &Kind{
					kind: "test-kind",
					code: "test-code",
					tpl:  newTemplate("test-message"),
				},
				err: errors.New("test-error"),
			},
			&action{
				attrs: &ErrorAttrs{
					kind: "test-kind",
					code: "test-code",
					name: "error",
					msg:  "test-message [test-error]",
					err:  errors.New("test-error"),
				},
			},
		),
		gen(
			"with an error of ErrorAttrs",
			[]string{cndInputErrorAttrs},
			[]string{actCheckReturned},
			&condition{
				kind: &Kind{
					kind: "test-kind",
					code: "test-code",
					tpl:  newTemplate("test-message"),
				},
				err: &ErrorAttrs{
					kind: "error-kind",
					code: "error-code",
					name: "error",
					msg:  "test-error",
				},
			},
			&action{
				attrs: &ErrorAttrs{
					kind: "test-kind",
					code: "test-code",
					name: "error",
					msg:  "test-message [test-error]",
				},
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			attrs := tt.C().kind.WithoutStack(tt.C().err, nil)
			testutil.Diff(t, tt.A().attrs, attrs, cmp.AllowUnexported(ErrorAttrs{}), cmpopts.EquateErrors())
		})
	}
}

func TestKind_WithStack(t *testing.T) {
	type condition struct {
		kind *Kind
		err  error
	}

	type action struct {
		attrs *ErrorAttrs
	}

	cndInputError := "pass an error"
	cndInputErrorAttrs := "pass an ErrorAttrs"
	actCheckReturned := "check returned ErrorAttrs"

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	tb.Condition(cndInputError, "pass an error which is not instance or ErrorAttrs as argument")
	tb.Condition(cndInputErrorAttrs, "pass an error of ErrorAttrs as argument")
	tb.Action(actCheckReturned, "check that the returned ErrorAttrs is the one with expected values")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"with no error",
			[]string{},
			[]string{actCheckReturned},
			&condition{
				kind: &Kind{
					kind: "test-kind",
					code: "test-code",
					tpl:  newTemplate("test-message"),
				},
			},
			&action{
				attrs: &ErrorAttrs{
					kind: "test-kind",
					code: "test-code",
					name: "error",
					msg:  "test-message",
				},
			},
		),
		gen(
			"with an error",
			[]string{cndInputError},
			[]string{actCheckReturned},
			&condition{
				kind: &Kind{
					kind: "test-kind",
					code: "test-code",
					tpl:  newTemplate("test-message"),
				},
				err: errors.New("test-error"),
			},
			&action{
				attrs: &ErrorAttrs{
					kind: "test-kind",
					code: "test-code",
					name: "error",
					msg:  "test-message [test-error]",
					err:  errors.New("test-error"),
				},
			},
		),
		gen(
			"with an error of ErrorAttrs",
			[]string{cndInputErrorAttrs},
			[]string{actCheckReturned},
			&condition{
				kind: &Kind{
					kind: "test-kind",
					code: "test-code",
					tpl:  newTemplate("test-message"),
				},
				err: &ErrorAttrs{
					kind: "error-kind",
					code: "error-code",
					name: "error",
					msg:  "test-error",
				},
			},
			&action{
				attrs: &ErrorAttrs{
					kind: "test-kind",
					code: "test-code",
					name: "error",
					msg:  "test-message [test-error]",
				},
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			errAttrs := tt.C().kind.WithStack(tt.C().err, nil)
			attrs := errAttrs.(*ErrorAttrs)
			testutil.Diff(t, tt.A().attrs, attrs, cmp.AllowUnexported(ErrorAttrs{}), cmpopts.IgnoreFields(ErrorAttrs{}, "stack"), cmpopts.EquateErrors())
			testutil.Diff(t, true, bytes.Contains(attrs.stack, []byte("goroutine")))
		})
	}
}

func TestKind_Is(t *testing.T) {
	type condition struct {
		kind *Kind
		err  error
	}

	type action struct {
		ok bool
	}

	cndNil := "nil"
	cndMatchCode := "match code"
	actCheckMatch := "check match"
	actCheckUnMatched := "check unmatched"

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	tb.Condition(cndNil, "input nil as the target")
	tb.Condition(cndMatchCode, "input error which matches as the target")
	tb.Action(actCheckMatch, "check that the errors matches")
	tb.Action(actCheckUnMatched, "check that the errors does not match")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"nil",
			[]string{cndNil},
			[]string{actCheckUnMatched},
			&condition{
				kind: &Kind{
					code: "test-code",
				},
				err: nil,
			},
			&action{
				ok: false,
			},
		),
		gen(
			"an error of errors.New",
			[]string{},
			[]string{actCheckUnMatched},
			&condition{
				kind: &Kind{
					code: "test-code",
				},
				err: errors.New("test-error"),
			},
			&action{
				ok: false,
			},
		),
		gen(
			"an ErrorAttrs which matches",
			[]string{cndMatchCode},
			[]string{actCheckMatch},
			&condition{
				kind: &Kind{
					code: "test-code",
				},
				err: &ErrorAttrs{
					code: "test-code",
				},
			},
			&action{
				ok: true,
			},
		),
		gen(
			"an ErrorAttrs which does not match",
			[]string{},
			[]string{actCheckUnMatched},
			&condition{
				kind: &Kind{
					code: "test-code",
				},
				err: &ErrorAttrs{
					code: "error-code",
				},
			},
			&action{
				ok: false,
			},
		),
		gen(
			"an ErrorAttrs containing another ErrorAttrs which matches",
			[]string{cndMatchCode},
			[]string{actCheckMatch},
			&condition{
				kind: &Kind{
					code: "test-code",
				},
				err: &ErrorAttrs{
					code: "error-code",
					err: &ErrorAttrs{
						code: "test-code",
					},
				},
			},
			&action{
				ok: true,
			},
		),
		gen(
			"an ErrorAttrs which mates containing another ErrorAttrs",
			[]string{cndMatchCode},
			[]string{actCheckMatch},
			&condition{
				kind: &Kind{
					code: "test-code",
				},
				err: &ErrorAttrs{
					code: "test-code",
					err: &ErrorAttrs{
						code: "error-code",
					},
				},
			},
			&action{
				ok: true,
			},
		),
		gen(
			"an ErrorAttrs containing another ErrorAttrs which does not match",
			[]string{},
			[]string{actCheckUnMatched},
			&condition{
				kind: &Kind{
					code: "test-code",
				},
				err: &ErrorAttrs{
					code: "error-code",
					err: &ErrorAttrs{
						code: "error-code",
					},
				},
			},
			&action{
				ok: false,
			},
		),
		gen(
			"an ErrorAttrs containing io.EOF which does not match",
			[]string{},
			[]string{actCheckUnMatched},
			&condition{
				kind: &Kind{
					code: "test-code",
				},
				err: &ErrorAttrs{
					code: "error-code",
					err:  io.EOF,
				},
			},
			&action{
				ok: false,
			},
		),
		gen(
			"an ErrorAttrs containing errors.New error which does not match",
			[]string{},
			[]string{actCheckUnMatched},
			&condition{
				kind: &Kind{
					kind: "test-kind",
					code: "test-code",
				},
				err: &ErrorAttrs{
					kind: "error-kind",
					code: "error-code",
					err:  errors.New("test-error"),
				},
			},
			&action{
				ok: false,
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			ok := tt.C().kind.Is(tt.C().err)
			testutil.Diff(t, tt.A().ok, ok)
		})
	}
}
