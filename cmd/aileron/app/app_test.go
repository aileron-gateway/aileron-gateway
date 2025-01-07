package app_test

import (
	stdcmp "cmp"
	"context"
	"errors"
	"io"
	"os"
	"regexp"
	"strings"
	"testing"

	k "github.com/aileron-gateway/aileron-gateway/apis/kernel"
	"github.com/aileron-gateway/aileron-gateway/cmd/aileron/app"
	"github.com/aileron-gateway/aileron-gateway/core"
	"github.com/aileron-gateway/aileron-gateway/kernel/api"
	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
	"github.com/google/go-cmp/cmp/cmpopts"
)

// testDir is the path to the test data.
// This path can be changed by the environmental variable.
var testDir = stdcmp.Or(os.Getenv("TEST_DIR"), "../../../test/")

func TestLoadEnvFiles(t *testing.T) {
	type condition struct {
		paths []string
	}

	type action struct {
		envs       map[string]string
		err        any // error or errorutil.Kind
		errPattern *regexp.Regexp
	}

	cndInputNil := "input nil"
	cndFileExists := "file exists"
	cndFileNotExists := "file not exists"
	actCheckValues := "check values"
	actCheckNoError := "no error"
	actCheckError := "non-nil error"

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	tb.Condition(cndInputNil, "input nil slice as the argument")
	tb.Condition(cndFileExists, "input file path which exists")
	tb.Condition(cndFileNotExists, "input file path which is not exists")
	tb.Action(actCheckValues, "check the environmental variable read from files")
	tb.Action(actCheckNoError, "check that there is no error")
	tb.Action(actCheckError, "check that non-nil error was returned")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"input nil",
			[]string{cndInputNil},
			[]string{actCheckNoError},
			&condition{
				paths: []string{},
			},
			&action{},
		),
		gen(
			"read env",
			[]string{cndFileExists},
			[]string{actCheckValues, actCheckNoError},
			&condition{
				paths: []string{
					testDir + "ut/cmd/aileron/app/env1.txt",
				},
			},
			&action{
				envs: map[string]string{"UNIT_TEST_ENV": "TEST_VALUE"},
			},
		),
		gen(
			"read invalid env",
			[]string{cndFileExists},
			[]string{actCheckError},
			&condition{
				paths: []string{
					testDir + "ut/cmd/aileron/app/env2.txt",
				},
			},
			&action{
				err:        app.ErrAppMainLoadEnv,
				errPattern: regexp.MustCompile(core.ErrPrefix + `failed to load environmental variables.`),
			},
		),
		gen(
			"no args",
			[]string{cndFileNotExists},
			[]string{actCheckError},
			&condition{
				paths: []string{
					testDir + "ut/cmd/aileron/app/not-exist.txt",
				},
			},
			&action{
				err:        app.ErrAppMainLoadEnv,
				errPattern: regexp.MustCompile(core.ErrPrefix + `failed to load environmental variables.`),
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			err := app.LoadEnvFiles(tt.C().paths)
			testutil.DiffError(t, tt.A().err, tt.A().errPattern, err)
			for k, v := range tt.A().envs {
				testutil.Diff(t, v, os.Getenv(k))
			}
		})
	}
}

type testServer struct {
	reqs []*api.Request
	res  *api.Response
	err  error
}

func (s *testServer) Serve(ctx context.Context, r *api.Request) (*api.Response, error) {
	s.reqs = append(s.reqs, r)
	return s.res, s.err
}

func TestLoadConfigFiles(t *testing.T) {
	type condition struct {
		server *testServer
		paths  []string
	}

	type action struct {
		reqs       []*api.Request
		err        any // error or errorutil.Kind
		errPattern *regexp.Regexp
	}

	cndInputNil := "input nil"
	cndFileExists := "file exists"
	cndJSON := "json file"
	cndYAML := "yaml file"
	cndUnsupported := "unsupported file extension"
	cndFieldInsufficient := "insufficient fields"
	cndFileInvalid := "file content is invalid"
	actCheckRequests := "check requests"
	actCheckNoError := "no error"
	actCheckError := "non-nil error"

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	tb.Condition(cndInputNil, "input nil for paths")
	tb.Condition(cndFileExists, "given files are exist")
	tb.Condition(cndJSON, "input file paths with .json extension")
	tb.Condition(cndYAML, "input file paths with .yaml extension")
	tb.Condition(cndUnsupported, "input file paths with unsupported file extension")
	tb.Condition(cndFieldInsufficient, "both apiVersion and kind fields are not defined")
	tb.Condition(cndFileInvalid, "invalid yaml or json format")
	tb.Action(actCheckRequests, "check the API request sent for the server")
	tb.Action(actCheckNoError, "check that there is no error")
	tb.Action(actCheckError, "check that non-nil error was returned")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"input nil",
			[]string{cndInputNil},
			[]string{actCheckNoError},
			&condition{
				server: &testServer{},
				paths:  []string{},
			},
			&action{
				reqs: nil,
			},
		),
		gen(
			"read yaml",
			[]string{cndYAML, cndFileExists},
			[]string{actCheckRequests, actCheckNoError},
			&condition{
				server: &testServer{},
				paths: []string{
					testDir + "ut/cmd/aileron/app/config1.yaml",
				},
			},
			&action{
				reqs: []*api.Request{
					{
						Method:  api.MethodPost,
						Key:     "test1/test2", // "apiVersion/kind"
						Format:  api.FormatYAML,
						Content: []byte("apiVersion: test1\nkind: test2\nmetadata:\n  namespace: test3\n  name: test4"),
					},
				},
			},
		),
		gen(
			"read yml",
			[]string{cndYAML, cndFileExists},
			[]string{actCheckRequests, actCheckNoError},
			&condition{
				server: &testServer{},
				paths: []string{
					testDir + "ut/cmd/aileron/app/config2.yml",
				},
			},
			&action{
				reqs: []*api.Request{
					{
						Method:  api.MethodPost,
						Key:     "test1/test2", // "apiVersion/kind"
						Format:  api.FormatYAML,
						Content: []byte("apiVersion: test1\nkind: test2\nmetadata:\n  namespace: test3\n  name: test4"),
					},
				},
			},
		),
		gen(
			"read json",
			[]string{cndJSON, cndFileExists},
			[]string{actCheckRequests, actCheckNoError},
			&condition{
				server: &testServer{},
				paths: []string{
					testDir + "ut/cmd/aileron/app/config3.json",
				},
			},
			&action{
				reqs: []*api.Request{
					{
						Method:  api.MethodPost,
						Key:     "test1/test2", // "apiVersion/kind"
						Format:  api.FormatJSON,
						Content: []byte(`{"apiVersion":"test1","kind":"test2","metadata":{"namespace":"test3","name":"test4"}}`),
					},
				},
			},
		),
		gen(
			"file not exists",
			[]string{cndYAML},
			[]string{actCheckError},
			&condition{
				server: &testServer{},
				paths: []string{
					testDir + "ut/cmd/aileron/app/not-exist.yaml",
				},
			},
			&action{
				reqs:       nil,
				err:        app.ErrAppMainLoadConfigs,
				errPattern: regexp.MustCompile(core.ErrPrefix + `failed to load configs.`),
			},
		),
		gen(
			"not supported extension",
			[]string{cndUnsupported, cndFileExists},
			[]string{actCheckNoError},
			&condition{
				server: &testServer{},
				paths: []string{
					testDir + "ut/cmd/aileron/app/config4.txt",
				},
			},
			&action{
				reqs: nil,
			},
		),
		gen(
			"invalid yaml format",
			[]string{cndYAML, cndFileExists, cndFileInvalid},
			[]string{actCheckError},
			&condition{
				server: &testServer{},
				paths: []string{
					testDir + "ut/cmd/aileron/app/config5.yaml",
				},
			},
			&action{
				reqs:       nil,
				err:        app.ErrAppMainLoadConfigs,
				errPattern: regexp.MustCompile(core.ErrPrefix + `failed to load configs.`),
			},
		),
		gen(
			"apiVersion and kind are empty",
			[]string{cndYAML, cndFileExists, cndFieldInsufficient},
			[]string{actCheckNoError},
			&condition{
				server: &testServer{},
				paths: []string{
					testDir + "ut/cmd/aileron/app/config6.yaml",
				},
			},
			&action{
				reqs: nil,
			},
		),
		gen(
			"server error",
			[]string{cndYAML, cndFileExists},
			[]string{actCheckRequests, actCheckError},
			&condition{
				server: &testServer{
					err: errors.New("test server error"),
				},
				paths: []string{
					testDir + "ut/cmd/aileron/app/config1.yaml",
				},
			},
			&action{
				reqs: []*api.Request{
					{
						Method:  api.MethodPost,
						Key:     "test1/test2", // "apiVersion/kind"
						Format:  api.FormatYAML,
						Content: []byte("apiVersion: test1\nkind: test2\nmetadata:\n  namespace: test3\n  name: test4"),
					},
				},
				err:        app.ErrAppMainLoadConfigs,
				errPattern: regexp.MustCompile(core.ErrPrefix + `failed to load configs.`),
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			err := app.LoadConfigFiles(tt.C().server, tt.C().paths)
			t.Logf("%#v\n", err)
			testutil.DiffError(t, tt.A().err, tt.A().errPattern, err)
			testutil.Diff(t, tt.A().reqs, tt.C().server.reqs)
		})
	}
}

func TestShowTemplate(t *testing.T) {
	type condition struct {
		server *testServer
		tpl    string
		out    string
	}

	type action struct {
		shouldExit bool
		exitCode   int
		contains   string
		reqs       []*api.Request
	}

	cndInputEmpty := "empty string template"
	cndInvalidTemplate := "invalid template"
	cndOutJSON := "expect json"
	cndOutYAML := "expect yaml"
	cndOutUnsupported := "unsupported template"
	cndServerError := "server returns error"
	actCheckRequest := "check API request"
	actCheckOutput := "check output"
	actCheckErrorExit := "exit with error"
	actCheckSuccessExit := "exit successfully"

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	tb.Condition(cndInputEmpty, "input an empty string as template")
	tb.Condition(cndInvalidTemplate, "input invalid format as template")
	tb.Condition(cndOutJSON, "set out format to expect json format")
	tb.Condition(cndOutYAML, "set out format to expect yaml format")
	tb.Condition(cndOutUnsupported, "input unsupported format for out")
	tb.Condition(cndServerError, "server returns an error for API request")
	tb.Action(actCheckRequest, "check the API request sent for the server")
	tb.Action(actCheckOutput, "check the string in the standard output")
	tb.Action(actCheckErrorExit, "check the function exist with an error")
	tb.Action(actCheckSuccessExit, "check the function exist without any errors")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"input nil",
			[]string{cndInputEmpty},
			[]string{},
			&condition{
				server: &testServer{
					res: &api.Response{Content: []byte("test response")},
				},
				tpl: "",
				out: "",
			},
			&action{
				shouldExit: false,
				reqs:       nil,
			},
		),
		gen(
			"apiGroup/apiVersion/kind",
			[]string{cndOutYAML},
			[]string{actCheckRequest, actCheckOutput, actCheckSuccessExit},
			&condition{
				server: &testServer{
					res: &api.Response{Content: []byte("test response")},
				},
				tpl: "apiGroup/apiVersion/kind",
				out: "",
			},
			&action{
				shouldExit: true,
				exitCode:   0,
				contains:   "test response",
				reqs: []*api.Request{
					{
						Method: api.MethodGet,
						Key:    "apiGroup/apiVersion/kind",
						Params: map[string]string{api.KeyAccept: string(api.FormatYAML)},
						Format: api.FormatProtoReference,
						Content: &k.Reference{
							APIVersion: "apiGroup/apiVersion",
							Kind:       "kind",
							Namespace:  "template",
							Name:       "template",
						},
					},
				},
			},
		),
		gen(
			"apiGroup/apiVersion/kind/namespace/name",
			[]string{cndOutYAML},
			[]string{actCheckRequest, actCheckOutput, actCheckSuccessExit},
			&condition{
				server: &testServer{
					res: &api.Response{Content: []byte("test response")},
				},
				tpl: "apiGroup/apiVersion/kind/namespace/name",
				out: "",
			},
			&action{
				shouldExit: true,
				exitCode:   0,
				contains:   "test response",
				reqs: []*api.Request{
					{
						Method: api.MethodGet,
						Key:    "apiGroup/apiVersion/kind",
						Params: map[string]string{api.KeyAccept: string(api.FormatYAML)},
						Format: api.FormatProtoReference,
						Content: &k.Reference{
							APIVersion: "apiGroup/apiVersion",
							Kind:       "kind",
							Namespace:  "namespace",
							Name:       "name",
						},
					},
				},
			},
		),
		gen(
			"output in json",
			[]string{cndOutJSON},
			[]string{actCheckRequest, actCheckOutput, actCheckSuccessExit},
			&condition{
				server: &testServer{
					res: &api.Response{Content: []byte("test response")},
				},
				tpl: "apiGroup/apiVersion/kind",
				out: "json",
			},
			&action{
				shouldExit: true,
				exitCode:   0,
				contains:   "test response",
				reqs: []*api.Request{
					{
						Method: api.MethodGet,
						Key:    "apiGroup/apiVersion/kind",
						Params: map[string]string{api.KeyAccept: string(api.FormatJSON)},
						Format: api.FormatProtoReference,
						Content: &k.Reference{
							APIVersion: "apiGroup/apiVersion",
							Kind:       "kind",
							Namespace:  "template",
							Name:       "template",
						},
					},
				},
			},
		),
		gen(
			"output in unsupported format",
			[]string{cndOutUnsupported},
			[]string{actCheckRequest, actCheckOutput, actCheckSuccessExit},
			&condition{
				server: &testServer{
					res: &api.Response{Content: []byte("test response")},
				},
				tpl: "apiGroup/apiVersion/kind",
				out: "unsupported",
			},
			&action{
				shouldExit: true,
				exitCode:   0,
				contains:   "test response",
				reqs: []*api.Request{
					{
						Method: api.MethodGet,
						Key:    "apiGroup/apiVersion/kind",
						Params: map[string]string{api.KeyAccept: string(api.FormatYAML)},
						Format: api.FormatProtoReference,
						Content: &k.Reference{
							APIVersion: "apiGroup/apiVersion",
							Kind:       "kind",
							Namespace:  "template",
							Name:       "template",
						},
					},
				},
			},
		),
		gen(
			"invalid template",
			[]string{cndInvalidTemplate},
			[]string{actCheckOutput, actCheckErrorExit},
			&condition{
				server: &testServer{
					res: &api.Response{Content: []byte("test response")},
				},
				tpl: "test1", // Invalid template format
				out: "",
			},
			&action{
				shouldExit: true,
				exitCode:   2,
				contains:   "invalid template",
				reqs:       nil,
			},
		),
		gen(
			"server error",
			[]string{cndServerError},
			[]string{actCheckRequest, actCheckOutput, actCheckErrorExit},
			&condition{
				server: &testServer{
					res: &api.Response{Content: []byte("test response")},
					err: errors.New("test server error"),
				},
				tpl: "apiGroup/apiVersion/kind",
				out: "",
			},
			&action{
				shouldExit: true,
				exitCode:   2,
				contains:   "test server error",
				reqs: []*api.Request{
					{
						Method: api.MethodGet,
						Key:    "apiGroup/apiVersion/kind",
						Params: map[string]string{api.KeyAccept: string(api.FormatYAML)},
						Format: api.FormatProtoReference,
						Content: &k.Reference{
							APIVersion: "apiGroup/apiVersion",
							Kind:       "kind",
							Namespace:  "template",
							Name:       "template",
						},
					},
				},
			},
		),
	}

	testutil.Register(table, testCases...)

	// Set default exit function at the end of this test.
	defer func() {
		app.Exit = os.Exit
	}()

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

			app.ShowTemplate(tt.C().server, tt.C().tpl, tt.C().out)
			testutil.Diff(t, tt.A().reqs, tt.C().server.reqs, cmpopts.IgnoreUnexported(k.Reference{}))

			w.Close()

			out, err := io.ReadAll(r)
			testutil.Diff(t, nil, err)
			testutil.Diff(t, true, strings.Contains(string(out), tt.A().contains))
			t.Log(string(out))
		})
	}
}
