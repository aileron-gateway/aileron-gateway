package main

import (
	stdcmp "cmp"
	"context"
	"errors"
	"os"
	"testing"

	"github.com/aileron-gateway/aileron-gateway/apis/kernel"
	"github.com/aileron-gateway/aileron-gateway/cmd/aileron/app"
	"github.com/aileron-gateway/aileron-gateway/kernel/api"
	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
)

// testDir is the path to the test data.
// This path can be changed by the environmental variable.
var testDir = stdcmp.Or(os.Getenv("TEST_DIR"), "../../test/")

// testRunner is a runner that will be called in the main function.
// This implements daemon.Runner and core.Runner interfaces.
type testRunner struct {
	err    error
	called bool
	got    context.Context
}

func (r *testRunner) Run(ctx context.Context) error {
	r.got = ctx
	r.called = true
	return r.err
}

func TestMainFunc(t *testing.T) {

	type condition struct {
		args   []string
		runner *testRunner
	}

	type action struct {
		shouldRunnerCalled bool
		shouldExit         bool
		exitCode           int
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	cndError := tb.Condition("runner error", "test runner returns a non-nil error")
	cndNoError := tb.Condition("runner no error", "test runner do not returns any error")
	actCheckExist0 := tb.Action("exit 0", "check that the main function exited with code 0")
	actCheckExist1 := tb.Action("exit 1", "check that the main function exited with code 1")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"exit with 0",
			[]string{cndNoError},
			[]string{actCheckExist0},
			&condition{
				args:   []string{"aileron", "-f", testDir + "ut/cmd/aileron/main-test.yaml"},
				runner: &testRunner{},
			},
			&action{
				shouldRunnerCalled: true,
				shouldExit:         false,
				exitCode:           0,
			},
		),
		gen(
			"exit with 1",
			[]string{cndError},
			[]string{actCheckExist1},
			&condition{
				args: []string{"aileron", "-f", testDir + "ut/cmd/aileron/main-test.yaml"},
				runner: &testRunner{
					err: errors.New("test"),
				},
			},
			&action{
				shouldRunnerCalled: true,
				shouldExit:         true,
				exitCode:           1,
			},
		),
	}

	testutil.Register(table, testCases...)

	defer func() {
		// This test change the server status.
		// So, set new fresh server on exit of this test.
		svr = api.NewDefaultServeMux()
	}()

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {

			tmpArgs := os.Args
			os.Args = tt.C().args
			defer func() {
				os.Args = tmpArgs
			}()

			app.Exit = func(code int) {
				testutil.Diff(t, true, tt.A().shouldExit)
				testutil.Diff(t, tt.A().exitCode, code)
			}
			defer func() { app.Exit = os.Exit }()

			svr = api.NewDefaultServeMux()
			svr.Handle("container/", api.NewContainerAPI())
			postTestResource(svr, "testRunner", tt.C().runner)

			main()

			testutil.Diff(t, tt.A().shouldRunnerCalled, tt.C().runner.called)

		})
	}

}

func postTestResource(server api.API[*api.Request, *api.Response], name string, res any) {
	ref := testResourceRef(name)
	req := &api.Request{
		Method:  api.MethodPost,
		Key:     ref.APIVersion + "/" + ref.Kind + "/" + ref.Namespace + "/" + ref.Name,
		Content: res,
	}
	if _, err := server.Serve(context.Background(), req); err != nil {
		panic(err)
	}
}

func testResourceRef(name string) *kernel.Reference {
	return &kernel.Reference{
		APIVersion: "container/v1",
		Kind:       "Container",
		Namespace:  "test",
		Name:       name,
	}
}
