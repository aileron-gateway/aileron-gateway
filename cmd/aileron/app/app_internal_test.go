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
	"github.com/aileron-gateway/aileron-gateway/internal/testutil"
	"github.com/aileron-gateway/aileron-gateway/kernel/api"
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

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"new app",
			&condition{},
			&action{
				app: &App{},
			},
		),
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			app := New()
			testutil.Diff(t, tt.A.app, app, cmp.AllowUnexported(App{}))
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

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"no args",
			&condition{
				args: []string{},
			},
			&action{
				args: []string{"-h"},
			},
		),
		gen(
			"with args",
			&condition{
				args: []string{"-f", "config.yaml"},
			},
			&action{
				args: []string{"-f", "config.yaml"},
			},
		),
	}

	Exit = func(_ int) {} // Avoid os.Exit during the tests.
	defer func() {
		Exit = os.Exit
	}()

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			a := &App{}
			a.ParseArgs(tt.C.args, tt.C.customFlags...)
			testutil.Diff(t, tt.A.args, a.args)
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

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"successfully run console",
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

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			err := tt.C.app.Run(tt.C.server)
			testutil.DiffError(t, tt.A.err, tt.A.errPattern, err)
		})
	}
}
