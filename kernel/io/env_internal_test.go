// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package io

import (
	"os"
	"strconv"
	"testing"

	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
)

func TestLoadEnv(t *testing.T) {
	type condition struct {
		overwrite bool
		envs      [][]byte
		preset    map[string]string
	}

	type action struct {
		checkVal map[string]string
	}

	CndInputNil := "input nil"
	CndInputEmpty := "input empty"
	CndDefineEnv := "define env"
	CndEmptyValue := "empty value"
	ActCheckResolved := "returned env or resolved value"
	ActCheckAsItIs := "returned as it is"
	ActCheckEmpty := "empty string returned"

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	tb.Condition(CndInputNil, "input nil value")
	tb.Condition(CndInputEmpty, "input empty string value")
	tb.Condition(CndDefineEnv, "define environmental variable")
	tb.Condition(CndEmptyValue, "the value of the environmental variable is an empty string")
	tb.Action(ActCheckResolved, "check that the environmental variable is resolved")
	tb.Action(ActCheckAsItIs, "check that the string is returned as it was passed by an argument")
	tb.Action(ActCheckEmpty, "check that the empty string is returned")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"input nil",
			[]string{CndInputNil},
			[]string{ActCheckEmpty},
			&condition{
				overwrite: false,
				envs: [][]byte{
					[]byte(""),
				},
				preset: map[string]string{},
			},
			&action{
				checkVal: map[string]string{},
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			// Initialize
			for k, v := range tt.C().preset {
				os.Setenv(k, v)
			}

			// Unset env on exit not to affect other tests.
			defer func() {
				for k := range tt.C().preset {
					os.Unsetenv(k)
				}
				for k := range tt.A().checkVal {
					os.Unsetenv(k)
				}
			}()

			for k, v := range tt.A().checkVal {
				val, exists := os.LookupEnv(k)
				testutil.Diff(t, true, exists)
				testutil.Diff(t, v, val)
			}
		})
	}
}

func TestResolve(t *testing.T) {
	type condition struct {
		inputNil   bool
		define     bool
		env        string
		expression string
	}

	type action struct {
		expect string
	}

	CndInputNil := "input nil"
	CndInputEmpty := "input empty"
	CndDefineEnv := "define env"
	CndEmptyValue := "empty value"
	ActCheckResolved := "returned env or resolved value"
	ActCheckAsItIs := "returned as it is"
	ActCheckEmpty := "empty string returned"

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	tb.Condition(CndInputNil, "input nil value")
	tb.Condition(CndInputEmpty, "input empty string value")
	tb.Condition(CndDefineEnv, "define environmental variable")
	tb.Condition(CndEmptyValue, "the value of the environmental variable is an empty string")
	tb.Action(ActCheckResolved, "check that the environmental variable is resolved")
	tb.Action(ActCheckAsItIs, "check that the string is returned as it was passed by an argument")
	tb.Action(ActCheckEmpty, "check that the empty string is returned")
	table := tb.Build()

	testVarName := "TEST_ENV_FOO"

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"input nil",
			[]string{CndInputNil},
			[]string{ActCheckEmpty},
			&condition{
				inputNil: true,
				define:   false,
			},
			&action{
				expect: "",
			},
		),
		gen(
			"input empty",
			[]string{CndInputEmpty},
			[]string{ActCheckAsItIs, ActCheckEmpty},
			&condition{
				inputNil:   false,
				define:     false,
				expression: "",
			},
			&action{
				expect: "",
			},
		),
		gen(
			"invalid expression $[ENV]",
			[]string{CndDefineEnv},
			[]string{ActCheckAsItIs},
			&condition{
				inputNil:   false,
				define:     true,
				env:        "dummy",
				expression: "$[" + testVarName + "]",
			},
			&action{
				expect: "",
			},
		),
		gen(
			"invalid expression ${ENV",
			[]string{CndDefineEnv},
			[]string{ActCheckAsItIs},
			&condition{
				inputNil:   false,
				define:     true,
				env:        "dummy",
				expression: "${" + testVarName + "",
			},
			&action{
				expect: "dummy",
			},
		),
		gen(
			"invalid expression $ENV}",
			[]string{CndDefineEnv},
			[]string{ActCheckAsItIs},
			&condition{
				inputNil:   false,
				define:     true,
				env:        "dummy",
				expression: "$" + testVarName + "}",
			},
			&action{
				expect: "dummy",
			},
		),
		gen(
			"invalid ${}",
			[]string{},
			[]string{ActCheckEmpty},
			&condition{
				expression: "${}",
			},
			&action{
				expect: "",
			},
		),
		gen(
			"invalid ${:-}",
			[]string{CndDefineEnv},
			[]string{ActCheckEmpty},
			&condition{
				expression: "${:-}",
			},
			&action{
				expect: "",
			},
		),
		gen(
			"invalid ${-}",
			[]string{CndDefineEnv},
			[]string{ActCheckEmpty},
			&condition{
				expression: "${-}",
			},
			&action{
				expect: "",
			},
		),
		gen(
			"invalid ${+}",
			[]string{CndDefineEnv},
			[]string{ActCheckEmpty},
			&condition{
				expression: "${+}",
			},
			&action{
				expect: "",
			},
		),
		gen(
			"invalid ${#}",
			[]string{},
			[]string{ActCheckResolved},
			&condition{
				expression: "${#}",
			},
			&action{
				expect: "0",
			},
		),
		gen(
			"resolve ${ENV} #1",
			[]string{CndDefineEnv},
			[]string{ActCheckResolved},
			&condition{
				define:     true,
				env:        "test_env_value",
				expression: "${" + testVarName + "}",
			},
			&action{
				expect: "test_env_value",
			},
		),
		gen(
			"resolve ${ENV} #2",
			[]string{CndDefineEnv, CndEmptyValue},
			[]string{ActCheckResolved, ActCheckEmpty},
			&condition{
				define:     true,
				env:        "",
				expression: "${" + testVarName + "}",
			},
			&action{
				expect: "",
			},
		),
		gen(
			"resolve ${ENV} #3",
			[]string{CndEmptyValue},
			[]string{ActCheckResolved, ActCheckEmpty},
			&condition{
				define:     false,
				env:        "",
				expression: "${" + testVarName + "}",
			},
			&action{
				expect: "",
			},
		),
		gen(
			"resolve ${ENV} #4",
			[]string{},
			[]string{ActCheckResolved, ActCheckEmpty},
			&condition{
				define:     false,
				env:        "test_env_value",
				expression: "${" + testVarName + "}",
			},
			&action{
				expect: "",
			},
		),
		gen(
			"resolve ${ENV:-default} #1",
			[]string{CndDefineEnv},
			[]string{ActCheckResolved},
			&condition{
				define:     true,
				env:        "test_env_value",
				expression: "${" + testVarName + ":-default}",
			},
			&action{
				expect: "test_env_value",
			},
		),
		gen(
			"resolve ${ENV:-default} #2",
			[]string{CndDefineEnv, CndEmptyValue},
			[]string{ActCheckResolved},
			&condition{
				define:     true,
				env:        "",
				expression: "${" + testVarName + ":-default}",
			},
			&action{
				expect: "default",
			},
		),
		gen(
			"resolve ${ENV:-default} #3",
			[]string{CndEmptyValue},
			[]string{ActCheckResolved},
			&condition{
				define:     false,
				env:        "",
				expression: "${" + testVarName + ":-default}",
			},
			&action{
				expect: "default",
			},
		),
		gen(
			"resolve ${ENV:-default} #4",
			[]string{},
			[]string{ActCheckResolved},
			&condition{
				define:     false,
				env:        "test_env_value",
				expression: "${" + testVarName + ":-default}",
			},
			&action{
				expect: "default",
			},
		),
		gen(
			"resolve ${ENV-default} #1",
			[]string{CndDefineEnv},
			[]string{ActCheckResolved},
			&condition{
				define:     true,
				env:        "test_env_value",
				expression: "${" + testVarName + "-default}",
			},
			&action{
				expect: "test_env_value",
			},
		),
		gen(
			"resolve ${ENV-default} #2",
			[]string{CndDefineEnv, CndEmptyValue},
			[]string{ActCheckResolved, ActCheckEmpty},
			&condition{
				define:     true,
				env:        "",
				expression: "${" + testVarName + "-default}",
			},
			&action{
				expect: "",
			},
		),
		gen(
			"resolve ${ENV-default} #3",
			[]string{CndEmptyValue},
			[]string{ActCheckResolved},
			&condition{
				define:     false,
				env:        "",
				expression: "${" + testVarName + "-default}",
			},
			&action{
				expect: "default",
			},
		),
		gen(
			"resolve ${ENV-default} #4",
			[]string{},
			[]string{ActCheckResolved},
			&condition{
				define:     false,
				env:        "test_env_value",
				expression: "${" + testVarName + "-default}",
			},
			&action{
				expect: "default",
			},
		),
		gen(
			"resolve ${ENV+default} #1",
			[]string{CndDefineEnv},
			[]string{ActCheckResolved},
			&condition{
				define:     true,
				env:        "test_env_value",
				expression: "${" + testVarName + "+default}",
			},
			&action{
				expect: "default",
			},
		),
		gen(
			"resolve ${ENV+default} #2",
			[]string{CndDefineEnv, CndEmptyValue},
			[]string{ActCheckResolved},
			&condition{
				define:     true,
				env:        "",
				expression: "${" + testVarName + "+default}",
			},
			&action{
				expect: "default",
			},
		),
		gen(
			"resolve ${ENV+default} #3",
			[]string{CndEmptyValue},
			[]string{ActCheckResolved, ActCheckEmpty},
			&condition{
				define:     false,
				env:        "",
				expression: "${" + testVarName + "+default}",
			},
			&action{
				expect: "",
			},
		),
		gen(
			"resolve ${ENV+default} #4",
			[]string{},
			[]string{ActCheckResolved, ActCheckEmpty},
			&condition{
				define:     false,
				env:        "test_env_value",
				expression: "${" + testVarName + "+default}",
			},
			&action{
				expect: "",
			},
		),
		gen(
			"resolve ${#ENV} #1",
			[]string{CndDefineEnv},
			[]string{ActCheckResolved},
			&condition{
				define:     true,
				env:        "test_env_value",
				expression: "${#" + testVarName + "}",
			},
			&action{
				expect: strconv.Itoa(len("test_env_value")),
			},
		),
		gen(
			"resolve ${#ENV} #2",
			[]string{CndDefineEnv, CndEmptyValue},
			[]string{ActCheckResolved},
			&condition{
				define:     true,
				env:        "",
				expression: "${#" + testVarName + "}",
			},
			&action{
				expect: strconv.Itoa(len("")),
			},
		),
		gen(
			"resolve ${#ENV} #3",
			[]string{CndEmptyValue},
			[]string{ActCheckResolved},
			&condition{
				define:     false,
				env:        "",
				expression: "${#" + testVarName + "}",
			},
			&action{
				expect: strconv.Itoa(len("")),
			},
		),
		gen(
			"resolve ${#ENV} #4",
			[]string{},
			[]string{ActCheckResolved},
			&condition{
				define:     false,
				env:        "test_env_value",
				expression: "${#" + testVarName + "}",
			},
			&action{
				expect: strconv.Itoa(len("")),
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			// Initialize
			func() {
				if tt.C().define {
					os.Setenv(testVarName, tt.C().env)
				}
			}()

			// Unset env on exit not to affect other tests.
			defer func() {
				os.Unsetenv(testVarName)
			}()

			var got []byte
			if tt.C().inputNil {
				got = resolve(nil)
			} else {
				got = resolve([]byte(tt.C().expression))
			}

			testutil.Diff(t, tt.A().expect, string(got))
		})
	}
}
