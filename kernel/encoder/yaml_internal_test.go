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

	cndNil := "input nil"
	cndInvalidVal := "input invalid value"
	actCheckExpected := "expected value returned"
	actCheckNoError := "no error"
	actCheckError := "expected error returned"

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	tb.Condition(cndNil, "give a valid encoded string")
	tb.Condition(cndInvalidVal, "give an invalid value which will result in an error")
	tb.Action(actCheckExpected, "Check that an expected value returned")
	tb.Action(actCheckNoError, "Check that returned error is nil")
	tb.Action(actCheckError, "Check that an expected error was returned")
	table := tb.Build()

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
			[]string{},
			[]string{actCheckExpected, actCheckNoError},
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
			[]string{cndNil},
			[]string{actCheckExpected, actCheckNoError},
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
			[]string{cndInvalidVal},
			[]string{actCheckError},
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
			[]string{cndInvalidVal},
			[]string{actCheckError},
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

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			out, err := MarshalYAML(tt.C().in)
			testutil.Diff(t, tt.A().err, err, cmpopts.EquateErrors())
			testutil.Diff(t, tt.A().out, string(out))
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

	cndNil := "input nil"
	cndInvalidVal := "input invalid value"
	actCheckExpected := "expected value returned"
	actCheckNoError := "no error"
	actCheckError := "expected error returned"

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	tb.Condition(cndNil, "give a nil value as an input")
	tb.Condition(cndInvalidVal, "give an invalid value which will result in an error")
	tb.Action(actCheckExpected, "Check that an expected value returned")
	tb.Action(actCheckNoError, "Check that returned error is nil")
	tb.Action(actCheckError, "Check that an expected error was returned")
	table := tb.Build()

	type testStruct struct {
		Name string `yaml:"name"`
		Age  int    `yaml:"age"`
	}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"decode yaml string",
			[]string{},
			[]string{actCheckExpected, actCheckNoError},
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
			[]string{},
			[]string{actCheckExpected, actCheckNoError},
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
			[]string{cndNil},
			[]string{actCheckExpected, actCheckNoError},
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
			[]string{cndInvalidVal},
			[]string{actCheckError},
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

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			err := UnmarshalYAML([]byte(tt.C().in), tt.C().into)
			testutil.Diff(t, tt.A().err, err, cmpopts.EquateErrors())
			testutil.Diff(t, tt.A().result, tt.C().into)
		})
	}
}

func TestUnmarshalYAMLFile(t *testing.T) {
	type condition struct {
		path string
		into any
	}

	type action struct {
		result any
		err    error
	}

	cndNil := "input nil"
	cndInvalidVal := "input invalid value"
	actCheckExpected := "expected value returned"
	actCheckNoError := "no error"
	actCheckError := "expected error returned"

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	tb.Condition(cndNil, "give a nil value as an argument")
	tb.Condition(cndInvalidVal, "give an invalid value which will result in an error")
	tb.Action(actCheckExpected, "Check that an expected value returned")
	tb.Action(actCheckNoError, "Check that returned error is nil")
	tb.Action(actCheckError, "Check that an expected error was returned")
	table := tb.Build()

	type testStruct struct {
		Name string `yaml:"name"`
		Age  int    `yaml:"age"`
	}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"decode yaml string",
			[]string{},
			[]string{actCheckExpected, actCheckNoError},
			&condition{
				path: testDir + "ut/kernel/encoder/yaml1.txt",
				into: &testStruct{},
			},
			&action{
				result: &testStruct{Name: "John Doe", Age: 20},
				err:    nil,
			},
		),
		gen(
			"decode into nil",
			[]string{},
			[]string{actCheckExpected, actCheckNoError},
			&condition{
				path: testDir + "ut/kernel/encoder/yaml1.txt",
				into: nil,
			},
			&action{
				result: nil,
				err:    nil,
			},
		),
		gen(
			"decode into func",
			[]string{cndInvalidVal},
			[]string{actCheckError},
			&condition{
				path: testDir + "ut/kernel/encoder/yaml1.txt",
				into: func() {},
			},
			&action{
				result: nil,
				err: &er.Error{
					Package:     ErrPkg,
					Type:        ErrTypeYaml,
					Description: ErrDscUnmarshal,
				},
			},
		),
		gen(
			"empty file",
			[]string{},
			[]string{actCheckExpected, actCheckError},
			&condition{
				path: testDir + "ut/kernel/encoder/empty.txt",
				into: &testStruct{},
			},
			&action{
				result: &testStruct{},
			},
		),
		gen(
			"not-exist file",
			[]string{},
			[]string{actCheckExpected, actCheckError},
			&condition{
				path: testDir + "ut/kernel/encoder/not-exist.txt",
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

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			err := UnmarshalYAMLFile(tt.C().path, tt.C().into)
			testutil.Diff(t, tt.A().err, err, cmpopts.EquateErrors())
			if err == nil {
				testutil.Diff(t, tt.A().result, tt.C().into)
			}
		})
	}
}
