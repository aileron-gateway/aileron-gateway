// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package encoder

import (
	"testing"

	"github.com/aileron-gateway/aileron-gateway/kernel/er"
	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
	"github.com/google/go-cmp/cmp/cmpopts"
)

// testDir is the path to the test data.
var testDir = "../../test/"

func TestMarshalJSON(t *testing.T) {
	type condition struct {
		in any
	}

	type action struct {
		out string
		err error
	}
	type testStruct struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"encode struct",
			&condition{
				in: &testStruct{Name: "John Doe", Age: 20},
			},
			&action{
				out: "{\n  \"name\": \"John Doe\",\n  \"age\": 20\n}\n",
				err: nil,
			},
		),
		gen(
			"nil",
			&condition{
				in: nil,
			},
			&action{
				out: "",
				err: nil,
			},
		),
		gen(
			"failed to marshal",
			&condition{
				in: func() {},
			},
			&action{
				out: "",
				err: &er.Error{
					Package:     ErrPkg,
					Type:        ErrTypeJSON,
					Description: ErrDscMarshal,
				},
			},
		),
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			out, err := MarshalJSON(tt.C.in)
			testutil.Diff(t, tt.A.err, err, cmpopts.EquateErrors())
			testutil.Diff(t, tt.A.out, string(out))
		})
	}
}

func TestUnmarshalJSON(t *testing.T) {
	type condition struct {
		in   string
		into any
	}

	type action struct {
		result any
		err    error
	}

	type testStruct struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"decode json string",
			&condition{
				in:   `{"name":"John Doe", "age":20}`,
				into: &testStruct{},
			},
			&action{
				result: &testStruct{Name: "John Doe", Age: 20},
				err:    nil,
			},
		),
		gen(
			"decode json string into valued struct",
			&condition{
				in:   `{"name":"John Doe"}`,
				into: &testStruct{Age: 20},
			},
			&action{
				result: &testStruct{Name: "John Doe", Age: 20},
				err:    nil,
			},
		),
		gen(
			"nil",
			&condition{
				in:   "",
				into: nil,
			},
			&action{
				result: nil,
				err:    nil,
			},
		),
		gen(
			"failed to marshal",
			&condition{
				in:   `{Invalid JSON}`,
				into: &testStruct{},
			},
			&action{
				result: &testStruct{},
				err: &er.Error{
					Package:     ErrPkg,
					Type:        ErrTypeJSON,
					Description: ErrDscUnmarshal,
				},
			},
		),
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			err := UnmarshalJSON([]byte(tt.C.in), tt.C.into)
			testutil.Diff(t, tt.A.err, err, cmpopts.EquateErrors())
			testutil.Diff(t, tt.A.result, tt.C.into)
		})
	}
}
