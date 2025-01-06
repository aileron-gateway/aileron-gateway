package app_test

import (
	"io"
	"os"
	"strings"
	"testing"

	"github.com/aileron-gateway/aileron-gateway/app"
	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
	"github.com/spf13/pflag"
)

func TestParseArgs(t *testing.T) {
	type condition struct {
		args        []string
		customFlags []*pflag.FlagSet
	}

	type action struct {
		shouldExit  bool
		exitCode    int
		checkOutput []string
	}

	cndBuildInfo := "build info"
	cndHelp := "help"
	cndInvalidFlag := "invalid flag"
	cndInvalidArgs := "invalid args"
	cndWithCustomFlag := "custom flag"
	actCheckExit := "check exit"
	actCheckNoExit := "check no-exit"
	actCheckOutput := "check output message"

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	tb.Condition(cndBuildInfo, "input --info flag to show build information")
	tb.Condition(cndHelp, "input --help flag to show help message")
	tb.Condition(cndInvalidFlag, "input invalid flags")
	tb.Condition(cndInvalidArgs, "input invalid argument")
	tb.Condition(cndWithCustomFlag, "input with extra custom flag")
	tb.Action(actCheckExit, "check that exit function was called")
	tb.Action(actCheckNoExit, "check that exit function was not called")
	tb.Action(actCheckOutput, "check that the message output was as expected")
	table := tb.Build()

	// Create a test flag set
	var testFlag bool
	testFS := pflag.NewFlagSet("test", pflag.ContinueOnError)
	testFS.BoolVarP(&testFlag, "test-flag", "z", false, "test flag")

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"no exit",
			[]string{},
			[]string{actCheckNoExit},
			&condition{
				args: []string{"-f", "config.yaml"},
			},
			&action{
				shouldExit:  false,
				checkOutput: []string{},
			},
		),
		gen(
			"version flag",
			[]string{cndBuildInfo},
			[]string{actCheckExit, actCheckOutput},
			&condition{
				args: []string{"--version"},
			},
			&action{
				shouldExit: true,
				exitCode:   0,
				checkOutput: []string{
					"UNSET",
				},
			},
		),
		gen(
			"info flag",
			[]string{cndBuildInfo},
			[]string{actCheckExit, actCheckOutput},
			&condition{
				args: []string{"--info"},
			},
			&action{
				shouldExit: true,
				exitCode:   0,
				checkOutput: []string{
					"go",
				},
			},
		),
		gen(
			"help flag",
			[]string{cndHelp},
			[]string{actCheckExit, actCheckOutput},
			&condition{
				args: []string{"--help"},
			},
			&action{
				shouldExit: true,
				exitCode:   0,
				checkOutput: []string{
					"Options :",
				},
			},
		),
		gen(
			"invalid flag",
			[]string{cndInvalidFlag},
			[]string{actCheckExit, actCheckOutput},
			&condition{
				args: []string{"--invalid-flag"},
			},
			&action{
				shouldExit: true,
				exitCode:   2,
				checkOutput: []string{
					"unknown flag: --invalid-flag",
					"Options :",
				},
			},
		),
		gen(
			"invalid arg",
			[]string{cndInvalidArgs},
			[]string{actCheckExit, actCheckOutput},
			&condition{
				args: []string{"invalid-arg"},
			},
			&action{
				shouldExit: true,
				exitCode:   2,
				checkOutput: []string{
					"invalid arguments:  [invalid-arg]",
					"Options :",
				},
			},
		),
		gen(
			"with custom flag",
			[]string{cndWithCustomFlag, cndHelp},
			[]string{actCheckNoExit, actCheckOutput},
			&condition{
				args:        []string{"--test-flag", "--help"},
				customFlags: []*pflag.FlagSet{testFS},
			},
			&action{
				shouldExit: true,
				checkOutput: []string{
					"-z, --test-flag",
					"Options :",
				},
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			tmp := os.Stdout
			defer func() {
				os.Stdout = tmp
			}()
			r, w, _ := os.Pipe()
			os.Stdout = w

			app.Exit = func(code int) {
				testutil.Diff(t, true, tt.A().shouldExit)
				testutil.Diff(t, tt.A().exitCode, code)
			}
			defer func() { app.Exit = os.Exit }()

			app.ParseArgs(tt.C().args, tt.C().customFlags...)
			w.Close()

			out, err := io.ReadAll(r)
			testutil.Diff(t, nil, err)

			for _, subStr := range tt.A().checkOutput {
				t.Log("expect contains", subStr)
				t.Log("but got", string(out))
				testutil.Diff(t, true, strings.Contains(string(out), subStr))
			}
		})
	}
}

func TestMetadataOptions(t *testing.T) {
	type condition struct {
	}

	type action struct {
		flags []string
	}

	cndNewFlag := "generate new flag"
	actCheckInfoFlag := "check file flag"
	actCheckHelpFlag := "check env flag"

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	tb.Condition(cndNewFlag, "generate a new flagset from options")
	tb.Action(actCheckInfoFlag, "check that the info option was registered to the flag")
	tb.Action(actCheckHelpFlag, "check that the help option was registered to the flag")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"check registered flags",
			[]string{cndNewFlag},
			[]string{actCheckInfoFlag, actCheckHelpFlag},
			&condition{},
			&action{
				flags: []string{
					"-i, --info",
					"-h, --help",
				},
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			opt := &app.MetadataOptions{}
			flg := opt.FlagSet()
			usage := flg.FlagUsages()

			for _, s := range tt.A().flags {
				t.Log("expect contains", s)
				t.Log("but got", usage)
				testutil.Diff(t, true, strings.Contains(usage, s))
			}
		})
	}
}

func TestBasicOptions(t *testing.T) {
	type condition struct {
	}

	type action struct {
		flags []string
	}

	cndNewFlag := "generate new flag"
	actCheckFileFlag := "check file flag"
	actCheckEnvFlag := "check env flag"

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	tb.Condition(cndNewFlag, "generate a new flagset from options")
	tb.Action(actCheckFileFlag, "check that the file option was registered to the flag")
	tb.Action(actCheckEnvFlag, "check that the env option was registered to the flag")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"check registered flags",
			[]string{cndNewFlag},
			[]string{actCheckFileFlag, actCheckEnvFlag},
			&condition{},
			&action{
				flags: []string{
					"-f, --file stringArray",
					"-e, --env stringArray",
				},
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			opt := &app.BasicOptions{}
			flg := opt.FlagSet()
			usage := flg.FlagUsages()

			for _, s := range tt.A().flags {
				t.Log("expect contains", s)
				t.Log("but got", usage)
				testutil.Diff(t, true, strings.Contains(usage, s))
			}
		})
	}
}
