package er

import (
	"errors"
	"io"
	"testing"

	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestError_Wrap(t *testing.T) {
	type condition struct {
		err  *Error
		wrap error
	}

	type action struct {
		err *Error
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"wrap nil",
			[]string{},
			[]string{},
			&condition{
				err: &Error{
					inner: nil,
				},
				wrap: nil,
			},
			&action{
				err: &Error{
					inner: nil,
				},
			},
		),
		gen(
			"wrap non-nil",
			[]string{},
			[]string{},
			&condition{
				err: &Error{
					inner: nil,
				},
				wrap: io.EOF,
			},
			&action{
				err: &Error{
					inner: io.EOF,
				},
			},
		),
		gen(
			"wrap already existing non-nil",
			[]string{},
			[]string{},
			&condition{
				err: &Error{
					inner: io.ErrUnexpectedEOF,
				},
				wrap: io.EOF,
			},
			&action{
				err: &Error{
					inner: errors.Join(io.EOF, io.ErrUnexpectedEOF),
				},
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			err := tt.C().err.Wrap(tt.C().wrap)
			testutil.Diff(t, tt.A().err, err, cmpopts.EquateErrors())
		})
	}
}

func TestError_Unwrap(t *testing.T) {
	type condition struct {
		err *Error
	}

	type action struct {
		err error
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"nil error",
			[]string{},
			[]string{},
			&condition{
				err: &Error{
					inner: nil,
				},
			},
			&action{
				err: nil,
			},
		),
		gen(
			"non nil error",
			[]string{},
			[]string{},
			&condition{
				err: &Error{
					inner: io.EOF,
				},
			},
			&action{
				err: io.EOF,
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			err := tt.C().err.Unwrap()
			testutil.Diff(t, tt.A().err, err, cmpopts.EquateErrors())
		})
	}
}

type testWrapErr struct {
	error
}

func (e *testWrapErr) Unwrap() error {
	return e.error
}

func TestError_Is(t *testing.T) {
	type condition struct {
		err    *Error
		target error
	}

	type action struct {
		is bool
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"nil",
			[]string{},
			[]string{},
			&condition{
				err:    nil,
				target: nil,
			},
			&action{
				is: true,
			},
		),
		gen(
			"match empty",
			[]string{},
			[]string{},
			&condition{
				err:    &Error{},
				target: &Error{},
			},
			&action{
				is: true,
			},
		),
		gen(
			"match all",
			[]string{},
			[]string{},
			&condition{
				err: &Error{
					Package:     "err",
					Type:        "type",
					Description: "desc",
					inner:       io.EOF,
					Detail:      "detail",
				},
				target: &Error{
					Package:     "err",
					Type:        "type",
					Description: "desc",
				},
			},
			&action{
				is: true,
			},
		),
		gen(
			"package not match",
			[]string{},
			[]string{},
			&condition{
				err:    &Error{Package: "err"},
				target: &Error{},
			},
			&action{
				is: false,
			},
		),
		gen(
			"Type not match",
			[]string{},
			[]string{},
			&condition{
				err:    &Error{Type: "err"},
				target: &Error{},
			},
			&action{
				is: false,
			},
		),
		gen(
			"description not match",
			[]string{},
			[]string{},
			&condition{
				err:    &Error{Description: "err"},
				target: &Error{},
			},
			&action{
				is: false,
			},
		),
		gen(
			"match after unwrap",
			[]string{},
			[]string{},
			&condition{
				err: &Error{Package: "test"},
				target: &testWrapErr{
					error: &Error{Package: "test"},
				},
			},
			&action{
				is: true,
			},
		),
		gen(
			"not match after unwrap",
			[]string{},
			[]string{},
			&condition{
				err: &Error{Package: "test"},
				target: &testWrapErr{
					error: io.EOF,
				},
			},
			&action{
				is: false,
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			is := tt.C().err.Is(tt.C().target)
			testutil.Diff(t, tt.A().is, is)
		})
	}
}

func TestError_Error(t *testing.T) {
	type condition struct {
		err *Error
	}

	type action struct {
		err string
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"all",
			[]string{},
			[]string{},
			&condition{
				err: &Error{
					inner:       io.EOF,
					Package:     "pkg",
					Type:        "type",
					Description: "desc",
					Detail:      "detail",
				},
			},
			&action{
				err: "pkg: type: desc detail [ EOF ]",
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			testutil.Diff(t, tt.A().err, tt.C().err.Error())
		})
	}
}
