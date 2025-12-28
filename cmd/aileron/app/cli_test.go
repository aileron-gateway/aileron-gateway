// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package app_test

import (
	"io"
	"os"
	"strings"
	"testing"

	"github.com/aileron-gateway/aileron-gateway/cmd/aileron/app"
	"github.com/aileron-gateway/aileron-gateway/internal/testutil"
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

	// Create a test flag set
	var testFlag bool
	testFS := pflag.NewFlagSet("test", pflag.ContinueOnError)
	testFS.BoolVarP(&testFlag, "test-flag", "z", false, "test flag")

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"no exit",
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

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			tmp := os.Stdout
			defer func() {
				os.Stdout = tmp
			}()
			r, w, _ := os.Pipe()
			os.Stdout = w

			app.Exit = func(code int) {
				testutil.Diff(t, true, tt.A.shouldExit)
				testutil.Diff(t, tt.A.exitCode, code)
			}
			defer func() { app.Exit = os.Exit }()

			app.ParseArgs(tt.C.args, tt.C.customFlags...)
			w.Close()

			out, err := io.ReadAll(r)
			testutil.Diff(t, nil, err)

			for _, subStr := range tt.A.checkOutput {
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

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"check registered flags",
			&condition{},
			&action{
				flags: []string{
					"-i, --info",
					"-h, --help",
				},
			},
		),
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			opt := &app.MetadataOptions{}
			flg := opt.FlagSet()
			usage := flg.FlagUsages()

			for _, s := range tt.A.flags {
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

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"check registered flags",
			&condition{},
			&action{
				flags: []string{
					"-f, --file stringArray",
					"-e, --env stringArray",
				},
			},
		),
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			opt := &app.BasicOptions{}
			flg := opt.FlagSet()
			usage := flg.FlagUsages()

			for _, s := range tt.A.flags {
				t.Log("expect contains", s)
				t.Log("but got", usage)
				testutil.Diff(t, true, strings.Contains(usage, s))
			}
		})
	}
}
