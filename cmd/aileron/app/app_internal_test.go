// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package app

import (
	"context"
	"errors"
	"os"
	"regexp"
	"testing"

	"github.com/aileron-gateway/aileron-gateway/core"
	"github.com/aileron-gateway/aileron-gateway/kernel/api"
	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
	"github.com/google/go-cmp/cmp"
	"github.com/spf13/pflag"
)

// testDir is the path to the test data.
var testDir = "../../../test/"

func TestNewApp(t *testing.T) {
	type condition struct {
	}

	type action struct {
		app *App
	}

	cndNew := "new App"
	actCheckApp := "check API request"

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	tb.Condition(cndNew, "create a new App instance with NewApp() function")
	tb.Action(actCheckApp, "check the values in returned App instance")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"new app",
			[]string{cndNew},
			[]string{actCheckApp},
			&condition{},
			&action{
				app: &App{},
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			app := New()
			testutil.Diff(t, tt.A().app, app, cmp.AllowUnexported(App{}))
		})
	}
}

func TestApp_ParseArgs(t *testing.T) {
	type condition struct {
		args        []string
		customFlags []*pflag.FlagSet
	}

	type action struct {
		args []string
	}

	cndWithArgs := "input non-zero args"
	actCheckArgs := "check args"

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	tb.Condition(cndWithArgs, "input at least 1 arguments")
	tb.Action(actCheckArgs, "check parsed arguments")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"no args",
			[]string{},
			[]string{actCheckArgs},
			&condition{
				args: []string{},
			},
			&action{
				args: []string{"-h"},
			},
		),
		gen(
			"with args",
			[]string{cndWithArgs},
			[]string{actCheckArgs},
			&condition{
				args: []string{"-f", "config.yaml"},
			},
			&action{
				args: []string{"-f", "config.yaml"},
			},
		),
	}

	testutil.Register(table, testCases...)

	Exit = func(_ int) {} // Avoid os.Exit during the tests.
	defer func() {
		Exit = os.Exit
	}()

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			a := &App{}
			a.ParseArgs(tt.C().args, tt.C().customFlags...)
			testutil.Diff(t, tt.A().args, a.args)
		})
	}
}

type testEntrypoint struct {
	err    error
	called bool
}

func (e *testEntrypoint) Run(ctx context.Context) error {
	e.called = true
	return e.err
}

type runTestServer struct {
	res *api.Response
	err error
}

func (s *runTestServer) Serve(ctx context.Context, req *api.Request) (*api.Response, error) {
	return s.res, s.err
}

func TestApp_Run(t *testing.T) {
	type condition struct {
		app    *App
		server *runTestServer
	}

	type action struct {
		err        any // error or errorutil.Kind
		errPattern *regexp.Regexp
	}

	cndConsole := "run as console"
	cndInvalidArgs := "invalid arg"
	cndFileLoadError := "file load error"
	cndServerError := "server returns error"
	cndEntrypointError := "entrypoint error"
	actCheckError := "non-nil error"
	actCheckNoError := "no error"

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	tb.Condition(cndConsole, "run as console app")
	tb.Condition(cndInvalidArgs, "args contains any invalid args which result in an error")
	tb.Condition(cndFileLoadError, "file paths are specified in the args which cannot read successfully")
	tb.Condition(cndServerError, "server returns an error")
	tb.Condition(cndEntrypointError, "entry point has invalid interface or returns an error")
	tb.Action(actCheckError, "check that a non-nil error was returned")
	tb.Action(actCheckNoError, "check that there was no error")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"successfully run console",
			[]string{cndConsole},
			[]string{actCheckNoError},
			&condition{
				app: &App{
					opts: &Options{
						Metadata: &MetadataOptions{},
						Basic:    &BasicOptions{},
					},
				},
				server: &runTestServer{
					res: &api.Response{Content: &testEntrypoint{}},
				},
			},
			&action{},
		),
		gen(
			"load env fails",
			[]string{cndConsole, cndFileLoadError},
			[]string{actCheckError},
			&condition{
				app: &App{
					args: []string{"-e", testDir + "ut/cmd/aileron/app/not-exist.txt"},
					opts: ParseArgs([]string{"-e", testDir + "ut/cmd/aileron/app/not-exist.txt"}),
				},
				server: &runTestServer{
					res: &api.Response{Content: &testEntrypoint{}},
				},
			},
			&action{
				err:        ErrAppMainLoadEnv,
				errPattern: regexp.MustCompile(core.ErrPrefix + `failed to load environmental variables.`),
			},
		),
		gen(
			"load config fails",
			[]string{cndConsole, cndFileLoadError},
			[]string{actCheckError},
			&condition{
				app: &App{
					args: []string{"-f", testDir + "ut/cmd/aileron/app/not-exist.txt"},
					opts: ParseArgs([]string{"-f", testDir + "ut/cmd/aileron/app/not-exist.txt"}),
				},
				server: &runTestServer{
					res: &api.Response{Content: &testEntrypoint{}},
				},
			},
			&action{
				err:        ErrAppMainLoadConfigs,
				errPattern: regexp.MustCompile(core.ErrPrefix + `failed to load configs.`),
			},
		),
		gen(
			"server error",
			[]string{cndConsole, cndServerError},
			[]string{actCheckError},
			&condition{
				app: &App{
					opts: &Options{
						Metadata: &MetadataOptions{},
						Basic:    &BasicOptions{},
					},
				},
				server: &runTestServer{
					res: &api.Response{Content: &testEntrypoint{}},
					err: errors.New("test server error"),
				},
			},
			&action{
				err:        ErrAppMainGetEntrypoint,
				errPattern: regexp.MustCompile(core.ErrPrefix + `failed to get entrypoint resource .*\[test server error\]`),
			},
		),
		gen(
			"entrypoint interface error",
			[]string{cndConsole, cndEntrypointError},
			[]string{actCheckError},
			&condition{
				app: &App{
					opts: &Options{
						Metadata: &MetadataOptions{},
						Basic:    &BasicOptions{},
					},
				},
				server: &runTestServer{
					res: &api.Response{Content: "conte must implement Runner interface"},
				},
			},
			&action{
				err:        ErrAppMainGetEntrypoint,
				errPattern: regexp.MustCompile(core.ErrPrefix + `failed to get entrypoint resource`),
			},
		),
		gen(
			"entrypoint run error",
			[]string{cndConsole, cndEntrypointError},
			[]string{actCheckError},
			&condition{
				app: &App{
					opts: &Options{
						Metadata: &MetadataOptions{},
						Basic:    &BasicOptions{},
					},
				},
				server: &runTestServer{
					res: &api.Response{
						Content: &testEntrypoint{
							err: errors.New("test entrypoint error"),
						},
					},
				},
			},
			&action{
				err:        ErrAppMainRun,
				errPattern: regexp.MustCompile(core.ErrPrefix + `running service failed .*\[test entrypoint error\]`),
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			err := tt.C().app.Run(tt.C().server)
			testutil.DiffError(t, tt.A().err, tt.A().errPattern, err)
		})
	}
}
