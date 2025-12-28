// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package encoder

import (
	"testing"

	"github.com/aileron-gateway/aileron-gateway/kernel/er"
	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestMarshalYaml(t *testing.T) {
	type condition struct {
		in any
	}

	type action struct {
		out string
		err error
	}

	type testStruct struct {
		Name string `yaml:"name"`
		Age  int    `yaml:"age"`
	}

	type testStructInvalidTag struct {
		Name string `yaml:"name"`
		Age  int    `yaml:"age,invalid"`
	}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"encode struct",
			&condition{
				in: &testStruct{Name: "John Doe", Age: 20},
			},
			&action{
				out: "name: John Doe\nage: 20\n",
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
			"panic on marshal",
			&condition{
				in: func() {},
			},
			&action{
				out: "",
				err: &er.Error{
					Package:     ErrPkg,
					Type:        ErrTypeYaml,
					Description: ErrDscMarshal,
				},
			},
		),
		gen(
			"failed to marshal",
			&condition{
				in: testStructInvalidTag{
					Name: "foo",
					Age:  20,
				},
			},
			&action{
				out: "",
				err: &er.Error{
					Package:     ErrPkg,
					Type:        ErrTypeYaml,
					Description: ErrDscMarshal,
				},
			},
		),
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			out, err := MarshalYAML(tt.C.in)
			testutil.Diff(t, tt.A.err, err, cmpopts.EquateErrors())
			testutil.Diff(t, tt.A.out, string(out))
		})
	}
}

func TestUnmarshalYaml(t *testing.T) {
	type condition struct {
		in   string
		into any
	}

	type action struct {
		result any
		err    error
	}

	type testStruct struct {
		Name string `yaml:"name"`
		Age  int    `yaml:"age"`
	}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"decode yaml string",
			&condition{
				in:   "name: John Doe\nage: 20\n",
				into: &testStruct{},
			},
			&action{
				result: &testStruct{Name: "John Doe", Age: 20},
				err:    nil,
			},
		),
		gen(
			"decode yaml string into valued struct",
			&condition{
				in:   "name: John Doe\nage: 20\n",
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
				in:   "Invalid:Yaml",
				into: &testStruct{},
			},
			&action{
				result: &testStruct{},
				err: &er.Error{
					Package:     ErrPkg,
					Type:        ErrTypeYaml,
					Description: ErrDscUnmarshal,
				},
			},
		),
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			err := UnmarshalYAML([]byte(tt.C.in), tt.C.into)
			testutil.Diff(t, tt.A.err, err, cmpopts.EquateErrors())
			testutil.Diff(t, tt.A.result, tt.C.into)
		})
	}
}
