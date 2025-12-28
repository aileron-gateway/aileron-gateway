// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package errorutil

import (
	"bytes"
	"errors"
	"io"
	"testing"

	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
	"github.com/aileron-projects/go/ztext"
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

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"valid template",
			&condition{
				code: "test-code",
				kind: "test-kind",
				tpl:  "this is {{test}}",
			},
			&action{
				kind: &Kind{
					code: "test-code",
					kind: "test-kind",
					tpl:  ztext.NewTemplate("this is {{test}}", "{{", "}}"),
				},
			},
		),
		gen(
			"template is empty",
			&condition{
				code: "test-code",
				kind: "test-kind",
				tpl:  "",
			},
			&action{
				kind: &Kind{
					code: "test-code",
					kind: "test-kind",
					tpl:  ztext.NewTemplate("", "{{", "}}"),
				},
			},
		),
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			if tt.C.panics {
				defer func() {
					e := recover().(error)
					testutil.Diff(t, tt.A.errMsg, e.Error())
				}()
			}

			k := NewKind(tt.C.code, tt.C.kind, tt.C.tpl)
			testutil.Diff(t, tt.A.kind, k, cmp.AllowUnexported(Kind{}, ztext.Template{}))
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

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"non-empty code",
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

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			code := tt.C.kind.Code()
			testutil.Diff(t, tt.A.code, code)
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

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"non-empty kind",
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

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			kind := tt.C.kind.Kind()
			testutil.Diff(t, tt.A.kind, kind)
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

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"with no error",
			&condition{
				kind: &Kind{
					kind: "test-kind",
					code: "test-code",
					tpl:  ztext.NewTemplate("test-message", "{{", "}}"),
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
			&condition{
				kind: &Kind{
					kind: "test-kind",
					code: "test-code",
					tpl:  ztext.NewTemplate("test-message", "{{", "}}"),
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
			&condition{
				kind: &Kind{
					kind: "test-kind",
					code: "test-code",
					tpl:  ztext.NewTemplate("test-message", "{{", "}}"),
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

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			attrs := tt.C.kind.WithoutStack(tt.C.err, nil)
			testutil.Diff(t, tt.A.attrs, attrs, cmp.AllowUnexported(ErrorAttrs{}), cmpopts.EquateErrors())
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

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"with no error",
			&condition{
				kind: &Kind{
					kind: "test-kind",
					code: "test-code",
					tpl:  ztext.NewTemplate("test-message", "{{", "}}"),
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
			&condition{
				kind: &Kind{
					kind: "test-kind",
					code: "test-code",
					tpl:  ztext.NewTemplate("test-message", "{{", "}}"),
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
			&condition{
				kind: &Kind{
					kind: "test-kind",
					code: "test-code",
					tpl:  ztext.NewTemplate("test-message", "{{", "}}"),
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

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			errAttrs := tt.C.kind.WithStack(tt.C.err, nil)
			attrs := errAttrs.(*ErrorAttrs)
			testutil.Diff(t, tt.A.attrs, attrs, cmp.AllowUnexported(ErrorAttrs{}), cmpopts.IgnoreFields(ErrorAttrs{}, "stack"), cmpopts.EquateErrors())
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

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"nil",
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

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			ok := tt.C.kind.Is(tt.C.err)
			testutil.Diff(t, tt.A.ok, ok)
		})
	}
}
