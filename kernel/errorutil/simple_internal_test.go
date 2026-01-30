// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package errorutil

import (
	"io"
	"testing"

	"github.com/aileron-gateway/aileron-gateway/internal/testutil"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestNewSimple(t *testing.T) {
	e := NewSimple(io.EOF, "message", "foo=%s", "bar")
	want := &SimpleError{
		Cause:   io.EOF,
		Message: "message",
		Detail:  "foo=bar",
	}
	testutil.Diff(t, want, e, cmpopts.EquateErrors())
}

func TestSimpleError_Unwrap(t *testing.T) {
	type condition struct {
		err *SimpleError
	}

	type action struct {
		err error
	}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"nil error",
			&condition{
				err: &SimpleError{
					Cause: nil,
				},
			},
			&action{
				err: nil,
			},
		),
		gen(
			"non nil error",
			&condition{
				err: &SimpleError{
					Cause: io.EOF,
				},
			},
			&action{
				err: io.EOF,
			},
		),
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			err := tt.C.err.Unwrap()
			testutil.Diff(t, tt.A.err, err, cmpopts.EquateErrors())
		})
	}
}

type testWrapErr struct {
	error
}

func (e *testWrapErr) Unwrap() error {
	return e.error
}

func TestSimpleError_Is(t *testing.T) {
	type condition struct {
		err    *SimpleError
		target error
	}

	type action struct {
		is bool
	}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"nil",
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
			&condition{
				err:    &SimpleError{},
				target: &SimpleError{},
			},
			&action{
				is: true,
			},
		),
		gen(
			"match all",
			&condition{
				err: &SimpleError{
					Cause:   io.EOF,
					Message: "message",
					Detail:  "detail",
				},
				target: &SimpleError{
					Cause:   io.EOF,
					Message: "message",
					Detail:  "detail",
				},
			},
			&action{
				is: true,
			},
		),
		gen(
			"message not match",
			&condition{
				err:    &SimpleError{Message: "foo"},
				target: &SimpleError{Message: "bar"},
			},
			&action{
				is: false,
			},
		),
		gen(
			"detail not match",
			&condition{
				err:    &SimpleError{Detail: "foo"},
				target: &SimpleError{Detail: "bar"},
			},
			&action{
				is: true,
			},
		),
		gen(
			"match after unwrap",
			&condition{
				err: &SimpleError{Message: "test"},
				target: &testWrapErr{
					error: &SimpleError{Message: "test"},
				},
			},
			&action{
				is: true,
			},
		),
		gen(
			"not match after unwrap",
			&condition{
				err: &SimpleError{Message: "test"},
				target: &testWrapErr{
					error: io.EOF,
				},
			},
			&action{
				is: false,
			},
		),
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			is := tt.C.err.Is(tt.C.target)
			testutil.Diff(t, tt.A.is, is)
		})
	}
}

func TestSimpleError_Error(t *testing.T) {
	type condition struct {
		err *SimpleError
	}

	type action struct {
		err string
	}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"all",
			&condition{
				err: &SimpleError{
					Cause:   io.EOF,
					Message: "message",
					Detail:  "detail",
				},
			},
			&action{
				err: "message detail [ EOF ]",
			},
		),
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			testutil.Diff(t, tt.A.err, tt.C.err.Error())
		})
	}
}
