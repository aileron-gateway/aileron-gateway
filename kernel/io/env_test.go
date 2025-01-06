package io_test

import (
	"errors"
	"os"
	"testing"

	"github.com/aileron-gateway/aileron-gateway/kernel/er"
	"github.com/aileron-gateway/aileron-gateway/kernel/io"
	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestResolveEnv(t *testing.T) {
	type condition struct {
		inputNil bool
		input    string
		extraEnv map[string]string
	}

	type action struct {
		expect string
	}

	CndInputNil := "input nil"
	CndInputEmpty := "input empty"
	CndInnerEnv := "define env"
	ActCheckResolved := "returned env or resolved value"
	ActCheckEmpty := "empty string returned"

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	tb.Condition(CndInputNil, "input nil value")
	tb.Condition(CndInputEmpty, "input empty string value")
	tb.Condition(CndInnerEnv, "environmental variable is defined in the other environmental variable")
	tb.Action(ActCheckResolved, "check that the environmental variable is resolved")
	tb.Action(ActCheckEmpty, "check that the empty string is returned")
	table := tb.Build()

	testVarName1 := "TEST_ENV_FOO"
	testVarValue1 := "foo_value"
	testVarName2 := "TEST_ENV_BAR"
	testVarValue2 := "bar_value"

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"input nil",
			[]string{CndInputNil},
			[]string{ActCheckEmpty},
			&condition{
				inputNil: true,
			},
			&action{
				expect: "",
			},
		),
		gen(
			"input empty",
			[]string{CndInputEmpty},
			[]string{ActCheckEmpty},
			&condition{
				inputNil: false,
				input:    "",
			},
			&action{
				expect: "",
			},
		),
		gen(
			"input text without env",
			[]string{},
			[]string{},
			&condition{
				inputNil: false,
				input: `
				test text input
				test text input
				test text input`,
			},
			&action{
				expect: `
				test text input
				test text input
				test text input`,
			},
		),
		gen(
			"input ${ENV}",
			[]string{},
			[]string{ActCheckResolved},
			&condition{
				inputNil: false,
				input:    "test env ${" + testVarName1 + "}.",
			},
			&action{
				expect: "test env " + testVarValue1 + ".",
			},
		),
		gen(
			"input ${ENV} with line break",
			[]string{},
			[]string{ActCheckResolved},
			&condition{
				inputNil: false,
				input:    "test env ${" + testVarName1 + "\n}.",
			},
			&action{
				expect: "test env ${" + testVarName1 + "\n}.",
			},
		),
		gen(
			"input ${${INNER_ENV}}",
			[]string{CndInnerEnv},
			[]string{ActCheckResolved},
			&condition{
				inputNil: false,
				input:    "test env ${${INNER_ENV}}.",
				extraEnv: map[string]string{"INNER_ENV": testVarName1},
			},
			&action{
				expect: "test env " + testVarValue1 + ".",
			},
		),
		gen(
			"input ${ENV:-${INNER_ENV}}",
			[]string{CndInnerEnv},
			[]string{ActCheckResolved},
			&condition{
				inputNil: false,
				input:    "test env ${OUTER_ENV:-${INNER_ENV}}.",
				extraEnv: map[string]string{
					"OUTER_ENV": "",
					"INNER_ENV": "default",
				},
			},
			&action{
				expect: "test env default.",
			},
		),
		gen(
			"input ${${INNER_ENV}:-default}",
			[]string{CndInnerEnv},
			[]string{ActCheckResolved},
			&condition{
				inputNil: false,
				input:    "test env ${${INNER_ENV}:-default}.",
				extraEnv: map[string]string{
					"INNER_ENV": "OUTER_ENV",
				},
			},
			&action{
				expect: "test env default.",
			},
		),
		gen(
			"input ${ENV-${INNER_ENV}}",
			[]string{CndInnerEnv},
			[]string{ActCheckResolved},
			&condition{
				inputNil: false,
				input:    "test env ${OUTER_ENV:-${INNER_ENV}}.",
				extraEnv: map[string]string{
					"INNER_ENV": "default",
				},
			},
			&action{
				expect: "test env default.",
			},
		),
		gen(
			"input ${${INNER_ENV}-default}",
			[]string{CndInnerEnv},
			[]string{ActCheckResolved},
			&condition{
				inputNil: false,
				input:    "test env ${${INNER_ENV}:-default}.",
				extraEnv: map[string]string{
					"OUTER_ENV": "",
					"INNER_ENV": "OUTER_ENV",
				},
			},
			&action{
				expect: "test env default.",
			},
		),
		gen(
			"input ${ENV+${INNER_ENV}}",
			[]string{CndInnerEnv},
			[]string{ActCheckResolved},
			&condition{
				inputNil: false,
				input:    "test env ${OUTER_ENV+${INNER_ENV}}.",
				extraEnv: map[string]string{
					"OUTER_ENV": "",
					"INNER_ENV": "default",
				},
			},
			&action{
				expect: "test env default.",
			},
		),
		gen(
			"input ${${INNER_ENV}+default}",
			[]string{CndInnerEnv},
			[]string{ActCheckResolved},
			&condition{
				inputNil: false,
				input:    "test env ${${INNER_ENV}+default}.",
				extraEnv: map[string]string{
					"OUTER_ENV": "",
					"INNER_ENV": "OUTER_ENV",
				},
			},
			&action{
				expect: "test env default.",
			},
		),
		gen(
			"input ${#${INNER_ENV}}",
			[]string{CndInnerEnv},
			[]string{ActCheckResolved},
			&condition{
				inputNil: false,
				input:    "test env ${#${INNER_ENV}}.",
				extraEnv: map[string]string{
					"OUTER_ENV": "1234567890",
					"INNER_ENV": "OUTER_ENV",
				},
			},
			&action{
				expect: "test env 10.",
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			// Initialize
			func() {
				os.Setenv(testVarName1, testVarValue1)
				os.Setenv(testVarName2, testVarValue2)
				for k, v := range tt.C().extraEnv {
					os.Setenv(k, v)
				}
			}()

			// Unset env on exit not to affect other tests.
			defer func() {
				os.Unsetenv(testVarName1)
				os.Unsetenv(testVarName2)
				for k := range tt.C().extraEnv {
					os.Unsetenv(k)
				}
			}()

			var got []byte
			if tt.C().inputNil {
				got = io.ResolveEnv(nil)
			} else {
				got = io.ResolveEnv([]byte(tt.C().input))
			}

			testutil.Diff(t, tt.A().expect, string(got))
		})
	}
}

func TestLoadEnv(t *testing.T) {
	type condition struct {
		overwrite bool
		input     [][]byte
		checkVar  string
	}

	type action struct {
		err         error
		expectedVal string
	}

	CndInputNil := "input nil"
	CndOverwrite := "overwrite"
	ActCheckErrorLoadEnv := "returned env or resolved value"
	ActCheckErrorSetEnv := "returned env or resolved value"

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	tb.Condition(CndInputNil, "input nil value")
	tb.Condition(CndOverwrite, "enable overwriting env")
	tb.Action(ActCheckErrorLoadEnv, "enable overwriting env")
	tb.Action(ActCheckErrorSetEnv, "enable overwriting env")
	table := tb.Build()

	testVarName := "TEST_ENV_FOO"
	testVarValue := "default"

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"input nil with overwrite",
			[]string{CndInputNil, CndOverwrite},
			[]string{},
			&condition{
				overwrite: true,
				checkVar:  testVarName,
				input:     nil,
			},
			&action{
				expectedVal: "default",
			},
		),
		gen(
			"input nil with overwrite",
			[]string{CndInputNil},
			[]string{},
			&condition{
				overwrite: false,
				checkVar:  testVarName,
				input:     nil,
			},
			&action{
				expectedVal: "default",
			},
		),
		gen(
			"input env with overwrite",
			[]string{CndOverwrite},
			[]string{},
			&condition{
				overwrite: true,
				checkVar:  testVarName,
				input: [][]byte{
					[]byte(`
					TEST_ENV_FOO=overwritten
					`),
				},
			},
			&action{
				expectedVal: "overwritten",
			},
		),
		gen(
			"input env without overwrite",
			[]string{},
			[]string{},
			&condition{
				overwrite: false,
				checkVar:  testVarName,
				input: [][]byte{
					[]byte(`
					TEST_ENV_FOO=overwritten
					`),
				},
			},
			&action{
				expectedVal: "default",
			},
		),
		gen(
			"failed to load env",
			[]string{},
			[]string{ActCheckErrorLoadEnv},
			&condition{
				input: [][]byte{
					[]byte(`
					this is invalid env file
					`),
				},
			},
			&action{
				err: &er.Error{
					Package:     io.ErrPkg,
					Type:        io.ErrTypeEnv,
					Description: io.ErrDscLoadEnv,
				},
			},
		),
		gen(
			"set env error",
			[]string{},
			[]string{},
			&condition{
				overwrite: true,
				checkVar:  testVarName,
				input: [][]byte{
					[]byte(`
					TEST_ENV_FOO=overwritten
					`),
				},
			},
			&action{
				expectedVal: "default",
				err: &er.Error{
					Package:     io.ErrPkg,
					Type:        io.ErrTypeEnv,
					Description: io.ErrDscSetEnv,
				},
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			// Initialize
			os.Setenv(testVarName, testVarValue)
			if tt.A().err != nil {
				io.Setenv = func(key, value string) error { return errors.New("env set error") }
				defer func() { io.Setenv = os.Setenv }()
			}

			// Unset env on exit not to affect other tests.
			defer os.Unsetenv(testVarName)

			err := io.LoadEnv(tt.C().overwrite, tt.C().input...)
			testutil.Diff(t, tt.A().err, err, cmpopts.EquateErrors())
			testutil.Diff(t, tt.A().expectedVal, os.Getenv(tt.C().checkVar))
		})
	}
}
