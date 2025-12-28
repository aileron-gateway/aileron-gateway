// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package app_test

import (
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
	"github.com/aileron-gateway/aileron-gateway/internal/testutil"
	"github.com/aileron-gateway/aileron-gateway/kernel/api"
	"github.com/google/go-cmp/cmp/cmpopts"
)

// testDir is the path to the test data.
var testDir = "../../../test/"

func TestLoadEnvFiles(t *testing.T) {
	type condition struct {
		paths []string
	}

	type action struct {
		envs       map[string]string
		err        any // error or errorutil.Kind
		errPattern *regexp.Regexp
	}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"input nil",
			&condition{
				paths: []string{},
			},
			&action{},
		),
		gen(
			"read env",
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

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			err := app.LoadEnvFiles(tt.C.paths)
			testutil.DiffError(t, tt.A.err, tt.A.errPattern, err)
			for k, v := range tt.A.envs {
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

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"input nil",
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

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			err := app.LoadConfigFiles(tt.C.server, tt.C.paths)
			t.Logf("%#v\n", err)
			testutil.DiffError(t, tt.A.err, tt.A.errPattern, err)
			testutil.Diff(t, tt.A.reqs, tt.C.server.reqs)
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

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"input nil",
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

	// Set default exit function at the end of this test.
	defer func() {
		app.Exit = os.Exit
	}()

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

			app.ShowTemplate(tt.C.server, tt.C.tpl, tt.C.out)
			testutil.Diff(t, tt.A.reqs, tt.C.server.reqs, cmpopts.IgnoreUnexported(k.Reference{}))

			w.Close()

			out, err := io.ReadAll(r)
			testutil.Diff(t, nil, err)
			testutil.Diff(t, true, strings.Contains(string(out), tt.A.contains))
			t.Log(string(out))
		})
	}
}

func TestSplitMultiDoc(t *testing.T) {
	type condition struct {
		docs []byte
		sep  string
	}
	type action struct {
		contents [][]byte
	}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"read docs with 1 block",
			&condition{
				docs: []byte("test"),
				sep:  "",
			},
			&action{
				contents: [][]byte{
					[]byte("test"),
				},
			},
		),
		gen(
			"read docs with 2 blocks",
			&condition{
				docs: []byte("test1\n---\ntest2"),
				sep:  "",
			},
			&action{
				contents: [][]byte{
					[]byte("test1"),
					[]byte("test2"),
				},
			},
		),
		gen(
			"read empty docs",
			&condition{
				docs: []byte(""),
				sep:  "",
			},
			&action{
				contents: nil,
			},
		),
		gen(
			"read docs double separator",
			&condition{
				docs: []byte("test1\n---\n---\ntest2"),
				sep:  "",
			},
			&action{
				contents: [][]byte{
					[]byte("test1"),
					[]byte("test2"),
				},
			},
		),
		gen(
			"read docs with only separator",
			&condition{
				docs: []byte("---\n"),
				sep:  "",
			},
			&action{
				contents: nil,
			},
		),
		gen(
			"read docs with non default separator",
			&condition{
				docs: []byte("test1\n***\ntest2"),
				sep:  "***\n",
			},
			&action{
				contents: [][]byte{
					[]byte("test1"),
					[]byte("test2"),
				},
			},
		),
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			contents := app.SplitMultiDoc(tt.C.docs, tt.C.sep)
			testutil.Diff(t, tt.A.contents, contents)
		})
	}
}
