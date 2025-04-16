// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package casbin

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
	"github.com/casbin/casbin/v3"
	"github.com/casbin/casbin/v3/persist"
	"github.com/google/go-cmp/cmp/cmpopts"
)

var (
	_ persist.Adapter = &noopAdapter{}
)

func TestNoopAdapter(t *testing.T) {
	type condition struct {
	}

	type action struct {
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"error",
			[]string{},
			[]string{},
			&condition{},
			&action{},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			policy := &noopAdapter{}

			var err error

			err = policy.AddPolicies("", "", nil)
			testutil.Diff(t, errNotImplemented, err, cmpopts.EquateErrors())

			err = policy.AddPolicy("", "", nil)
			testutil.Diff(t, errNotImplemented, err, cmpopts.EquateErrors())

			err = policy.LoadPolicy(nil)
			testutil.Diff(t, errNotImplemented, err, cmpopts.EquateErrors())

			err = policy.RemoveFilteredPolicy("", "", 0)
			testutil.Diff(t, errNotImplemented, err, cmpopts.EquateErrors())

			err = policy.RemovePolicies("", "", nil)
			testutil.Diff(t, errNotImplemented, err, cmpopts.EquateErrors())

			err = policy.RemovePolicy("", "", nil)
			testutil.Diff(t, errNotImplemented, err, cmpopts.EquateErrors())

			err = policy.SavePolicy(nil)
			testutil.Diff(t, errNotImplemented, err, cmpopts.EquateErrors())
		})
	}
}

func TestCsvAdapter(t *testing.T) {
	rootDir := "../../../test/ut/app/casbin/csv/"

	type condition struct {
		model  string // With extension ".conf"
		policy string // Without extension ".csv"
		sub    any
		path   string
		method string
	}

	type action struct {
		authorized       bool // If authorized with ${policy}.csv
		authorizedReload bool // If authorized with ${policy}_reload.csv
		loadErr          *regexp.Regexp
		enforceErr       *regexp.Regexp
		reLoadErr        *regexp.Regexp
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"abac/success",
			[]string{},
			[]string{},
			&condition{
				model:  rootDir + "abac_model.conf",
				policy: rootDir + "abac_policy", // abac_policy.csv & abac_policy_reload.csv
				sub:    map[string]any{"age": 21},
				path:   "/",
				method: http.MethodPost, // Only age>20
			},
			&action{
				authorized:       true,
				authorizedReload: false,
			},
		),
		gen(
			"abac/failure",
			[]string{},
			[]string{},
			&condition{
				model:  rootDir + "abac_model.conf",
				policy: rootDir + "abac_policy",
				sub:    map[string]any{"age": 5},
				path:   "/",
				method: http.MethodGet,
			},
			&action{
				authorized:       false,
				authorizedReload: false,
			},
		),
		gen(
			"rbac/success",
			[]string{},
			[]string{},
			&condition{
				model:  rootDir + "rbac_model.conf",
				policy: rootDir + "rbac_policy",
				sub:    map[string]any{"name": "alice"},
				path:   "/foo/123",
				method: http.MethodPost,
			},
			&action{
				authorized:       true,  // alice is admin
				authorizedReload: false, // alice is no longer admin
			},
		),
		gen(
			"rbac/failure",
			[]string{},
			[]string{},
			&condition{
				model:  rootDir + "rbac_model.conf",
				policy: rootDir + "rbac_policy",
				sub:    map[string]any{"name": "bob"},
				path:   "/bar/123",
				method: http.MethodPost,
			},
			&action{
				authorized:       false, // bob is not admin
				authorizedReload: true,  // now bob is admin
			},
		),
		gen(
			"v9",
			[]string{},
			[]string{},
			&condition{
				model:  rootDir + "model_v9.conf",
				policy: rootDir + "policy_v9",
				method: http.MethodGet,
			},
			&action{
				authorized:       true,
				authorizedReload: false,
			},
		),
		gen(
			"no policy",
			[]string{},
			[]string{},
			&condition{
				model:  rootDir + "abac_model.conf",
				policy: rootDir + "no_policy",
			},
			&action{
				authorized:       false,
				authorizedReload: false,
				// Casbin panics when no policy line found.
				enforceErr: regexp.MustCompile(`invalid memory address or nil pointer dereference`),
				reLoadErr:  regexp.MustCompile(`(cannot find the file|no such file)`),
			},
		),
		gen(
			"empty file",
			[]string{},
			[]string{},
			&condition{
				model:  rootDir + "abac_model.conf",
				policy: rootDir + "empty",
			},
			&action{
				authorized:       false,
				authorizedReload: false,
				// Casbin panics when no policy line found.
				enforceErr: regexp.MustCompile(`invalid memory address or nil pointer dereference`),
				reLoadErr:  regexp.MustCompile(`(cannot find the file|no such file)`),
			},
		),
		gen(
			"file not found",
			[]string{},
			[]string{},
			&condition{
				model:  rootDir + "abac_model.conf",
				policy: rootDir + "not_found",
			},
			&action{
				loadErr: regexp.MustCompile(`(cannot find the file|no such file)`),
			},
		),
		gen(
			"invalid format",
			[]string{},
			[]string{},
			&condition{
				model:  rootDir + "abac_model.conf",
				policy: rootDir + "invalid",
			},
			&action{
				loadErr: regexp.MustCompile(`invalid memory address or nil pointer dereference`),
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			policy := &csvAdapter{
				Adapter:  &noopAdapter{},
				filePath: tt.C().policy + ".csv",
			}

			enf, err := casbin.NewEnforcer(tt.C().model, policy)
			t.Logf("%s\n", err)
			if tt.A().loadErr != nil {
				testutil.Diff(t, true, tt.A().loadErr.MatchString(err.Error()))
				testutil.Diff(t, (*casbin.Enforcer)(nil), enf)
				return
			} else {
				testutil.Diff(t, nil, err)
			}
			enf.AddFunction("mapValue", mapValue)

			ok, err := enf.Enforce(tt.C().sub, tt.C().path, tt.C().method)
			testutil.Diff(t, tt.A().authorized, ok)
			t.Logf("%s\n", err)
			if tt.A().enforceErr != nil {
				testutil.Diff(t, true, tt.A().enforceErr.MatchString(err.Error()))
			} else {
				testutil.Diff(t, nil, err)
			}

			policy.filePath = tt.C().policy + "_reload.csv"
			err = enf.LoadPolicy()
			t.Logf("%s\n", err)
			if tt.A().reLoadErr != nil {
				testutil.Diff(t, true, tt.A().reLoadErr.MatchString(err.Error()))
			} else {
				testutil.Diff(t, nil, err)
			}
			ok, _ = enf.Enforce(tt.C().sub, tt.C().path, tt.C().method)
			testutil.Diff(t, tt.A().authorizedReload, ok)
		})
	}
}

func TestJSONAdapter(t *testing.T) {
	rootDir := "../../../test/ut/app/casbin/json/"

	type condition struct {
		model  string // With extension ".conf"
		policy string // Without extension ".json"
		sub    any
		path   string
		method string
	}

	type action struct {
		authorized       bool // If authorized with ${policy}.json
		authorizedReload bool // If authorized with ${policy}_reload.json
		loadErr          *regexp.Regexp
		enforceErr       *regexp.Regexp
		reLoadErr        *regexp.Regexp
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"abac/success",
			[]string{},
			[]string{},
			&condition{
				model:  rootDir + "abac_model.conf",
				policy: rootDir + "abac_policy", // abac_policy.json & abac_policy_reload.json
				sub:    map[string]any{"age": 21},
				path:   "/",
				method: http.MethodPost, // Only age>20
			},
			&action{
				authorized:       true,
				authorizedReload: false,
			},
		),
		gen(
			"abac/failure",
			[]string{},
			[]string{},
			&condition{
				model:  rootDir + "abac_model.conf",
				policy: rootDir + "abac_policy",
				sub:    map[string]any{"age": 5},
				path:   "/",
				method: http.MethodGet,
			},
			&action{
				authorized:       false,
				authorizedReload: false,
			},
		),
		gen(
			"rbac/success",
			[]string{},
			[]string{},
			&condition{
				model:  rootDir + "rbac_model.conf",
				policy: rootDir + "rbac_policy",
				sub:    map[string]any{"name": "alice"},
				path:   "/foo/123",
				method: http.MethodPost,
			},
			&action{
				authorized:       true,  // alice is admin
				authorizedReload: false, // alice is no longer admin
			},
		),
		gen(
			"rbac/failure",
			[]string{},
			[]string{},
			&condition{
				model:  rootDir + "rbac_model.conf",
				policy: rootDir + "rbac_policy",
				sub:    map[string]any{"name": "bob"},
				path:   "/bar/123",
				method: http.MethodPost,
			},
			&action{
				authorized:       false, // bob is not admin
				authorizedReload: true,  // now bob is admin
			},
		),
		gen(
			"v9",
			[]string{},
			[]string{},
			&condition{
				model:  rootDir + "model_v9.conf",
				policy: rootDir + "policy_v9",
				method: http.MethodGet,
			},
			&action{
				authorized:       true,
				authorizedReload: false,
			},
		),
		gen(
			"no policy",
			[]string{},
			[]string{},
			&condition{
				model:  rootDir + "abac_model.conf",
				policy: rootDir + "no_policy",
			},
			&action{
				loadErr: regexp.MustCompile(`cannot unmarshal object into Go value`),
			},
		),
		gen(
			"empty file",
			[]string{},
			[]string{},
			&condition{
				model:  rootDir + "abac_model.conf",
				policy: rootDir + "empty",
			},
			&action{
				loadErr: regexp.MustCompile(`unexpected end of JSON input`),
			},
		),
		gen(
			"file not found",
			[]string{},
			[]string{},
			&condition{
				model:  rootDir + "abac_model.conf",
				policy: rootDir + "not_found",
			},
			&action{
				loadErr: regexp.MustCompile(`(cannot find the file|no such file)`),
			},
		),
		gen(
			"invalid format",
			[]string{},
			[]string{},
			&condition{
				model:  rootDir + "abac_model.conf",
				policy: rootDir + "invalid",
			},
			&action{
				loadErr: regexp.MustCompile(`invalid character`),
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			policy := &jsonAdapter{
				Adapter:  &noopAdapter{},
				filePath: tt.C().policy + ".json",
			}

			enf, err := casbin.NewEnforcer(tt.C().model, policy)
			t.Logf("%s\n", err)
			if tt.A().loadErr != nil {
				testutil.Diff(t, true, tt.A().loadErr.MatchString(err.Error()))
				testutil.Diff(t, (*casbin.Enforcer)(nil), enf)
				return
			} else {
				testutil.Diff(t, nil, err)
			}
			enf.AddFunction("mapValue", mapValue)

			ok, err := enf.Enforce(tt.C().sub, tt.C().path, tt.C().method)
			testutil.Diff(t, tt.A().authorized, ok)
			t.Logf("%s\n", err)
			if tt.A().enforceErr != nil {
				testutil.Diff(t, true, tt.A().enforceErr.MatchString(err.Error()))
			} else {
				testutil.Diff(t, nil, err)
			}

			policy.filePath = tt.C().policy + "_reload.json"
			err = enf.LoadPolicy()
			t.Logf("%s\n", err)
			if tt.A().reLoadErr != nil {
				testutil.Diff(t, true, tt.A().reLoadErr.MatchString(err.Error()))
			} else {
				testutil.Diff(t, nil, err)
			}
			ok, _ = enf.Enforce(tt.C().sub, tt.C().path, tt.C().method)
			testutil.Diff(t, tt.A().authorizedReload, ok)
		})
	}
}

func TestYamlAdapter(t *testing.T) {
	rootDir := "../../../test/ut/app/casbin/yaml/"

	type condition struct {
		model  string // With extension ".conf"
		policy string // Without extension ".yaml"
		sub    any
		path   string
		method string
	}

	type action struct {
		authorized       bool // If authorized with ${policy}.yaml
		authorizedReload bool // If authorized with ${policy}_reload.yaml
		loadErr          *regexp.Regexp
		enforceErr       *regexp.Regexp
		reLoadErr        *regexp.Regexp
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"abac/success",
			[]string{},
			[]string{},
			&condition{
				model:  rootDir + "abac_model.conf",
				policy: rootDir + "abac_policy", // abac_policy.csv & abac_policy_reload.csv
				sub:    map[string]any{"age": 21},
				path:   "/",
				method: http.MethodPost, // Only age>20
			},
			&action{
				authorized:       true,
				authorizedReload: false,
			},
		),
		gen(
			"abac/failure",
			[]string{},
			[]string{},
			&condition{
				model:  rootDir + "abac_model.conf",
				policy: rootDir + "abac_policy",
				sub:    map[string]any{"age": 5},
				path:   "/",
				method: http.MethodGet,
			},
			&action{
				authorized:       false,
				authorizedReload: false,
			},
		),
		gen(
			"rbac/success",
			[]string{},
			[]string{},
			&condition{
				model:  rootDir + "rbac_model.conf",
				policy: rootDir + "rbac_policy",
				sub:    map[string]any{"name": "alice"},
				path:   "/foo/123",
				method: http.MethodPost,
			},
			&action{
				authorized:       true,  // alice is admin
				authorizedReload: false, // alice is no longer admin
			},
		),
		gen(
			"rbac/failure",
			[]string{},
			[]string{},
			&condition{
				model:  rootDir + "rbac_model.conf",
				policy: rootDir + "rbac_policy",
				sub:    map[string]any{"name": "bob"},
				path:   "/bar/123",
				method: http.MethodPost,
			},
			&action{
				authorized:       false, // bob is not admin
				authorizedReload: true,  // now bob is admin
			},
		),
		gen(
			"v9",
			[]string{},
			[]string{},
			&condition{
				model:  rootDir + "model_v9.conf",
				policy: rootDir + "policy_v9",
				method: http.MethodGet,
			},
			&action{
				authorized:       true,
				authorizedReload: false,
			},
		),
		gen(
			"no policy",
			[]string{},
			[]string{},
			&condition{
				model:  rootDir + "abac_model.conf",
				policy: rootDir + "no_policy",
			},
			&action{
				enforceErr: regexp.MustCompile(`invalid memory address or nil pointer dereference`),
				reLoadErr:  regexp.MustCompile(`(cannot find the file|no such file)`),
			},
		),
		gen(
			"empty file",
			[]string{},
			[]string{},
			&condition{
				model:  rootDir + "abac_model.conf",
				policy: rootDir + "empty",
			},
			&action{
				enforceErr: regexp.MustCompile(`invalid memory address or nil pointer dereference`),
				reLoadErr:  regexp.MustCompile(`(cannot find the file|no such file)`),
			},
		),
		gen(
			"file not found",
			[]string{},
			[]string{},
			&condition{
				model:  rootDir + "abac_model.conf",
				policy: rootDir + "not_found",
			},
			&action{
				loadErr: regexp.MustCompile(`(cannot find the file|no such file)`),
			},
		),
		gen(
			"invalid format",
			[]string{},
			[]string{},
			&condition{
				model:  rootDir + "abac_model.conf",
				policy: rootDir + "invalid",
			},
			&action{
				loadErr: regexp.MustCompile(`cannot unmarshal`),
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			policy := &yamlAdapter{
				Adapter:  &noopAdapter{},
				filePath: tt.C().policy + ".yaml",
			}

			enf, err := casbin.NewEnforcer(tt.C().model, policy)
			t.Logf("%s\n", err)
			if tt.A().loadErr != nil {
				testutil.Diff(t, true, tt.A().loadErr.MatchString(err.Error()))
				testutil.Diff(t, (*casbin.Enforcer)(nil), enf)
				return
			} else {
				testutil.Diff(t, nil, err)
			}
			enf.AddFunction("mapValue", mapValue)

			ok, err := enf.Enforce(tt.C().sub, tt.C().path, tt.C().method)
			testutil.Diff(t, tt.A().authorized, ok)
			t.Logf("%s\n", err)
			if tt.A().enforceErr != nil {
				testutil.Diff(t, true, tt.A().enforceErr.MatchString(err.Error()))
			} else {
				testutil.Diff(t, nil, err)
			}

			policy.filePath = tt.C().policy + "_reload.yaml"
			err = enf.LoadPolicy()
			t.Logf("%s\n", err)
			if tt.A().reLoadErr != nil {
				testutil.Diff(t, true, tt.A().reLoadErr.MatchString(err.Error()))
			} else {
				testutil.Diff(t, nil, err)
			}
			ok, _ = enf.Enforce(tt.C().sub, tt.C().path, tt.C().method)
			testutil.Diff(t, tt.A().authorizedReload, ok)
		})
	}
}

func TestXmlAdapter(t *testing.T) {
	rootDir := "../../../test/ut/app/casbin/xml/"

	type condition struct {
		model  string // With extension ".conf"
		policy string // Without extension ".xml"
		sub    any
		path   string
		method string
	}

	type action struct {
		authorized       bool // If authorized with ${policy}.xml
		authorizedReload bool // If authorized with ${policy}_reload.xml
		loadErr          *regexp.Regexp
		enforceErr       *regexp.Regexp
		reLoadErr        *regexp.Regexp
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"abac/success",
			[]string{},
			[]string{},
			&condition{
				model:  rootDir + "abac_model.conf",
				policy: rootDir + "abac_policy", // abac_policy.xml & abac_policy_reload.xml
				sub:    map[string]any{"age": 21},
				path:   "/",
				method: http.MethodPost, // Only age>20
			},
			&action{
				authorized:       true,
				authorizedReload: false,
			},
		),
		gen(
			"abac/failure",
			[]string{},
			[]string{},
			&condition{
				model:  rootDir + "abac_model.conf",
				policy: rootDir + "abac_policy",
				sub:    map[string]any{"age": 5},
				path:   "/",
				method: http.MethodGet,
			},
			&action{
				authorized:       false,
				authorizedReload: false,
			},
		),
		gen(
			"rbac/success",
			[]string{},
			[]string{},
			&condition{
				model:  rootDir + "rbac_model.conf",
				policy: rootDir + "rbac_policy",
				sub:    map[string]any{"name": "alice"},
				path:   "/foo/123",
				method: http.MethodPost,
			},
			&action{
				authorized:       true,  // alice is admin
				authorizedReload: false, // alice is no longer admin
			},
		),
		gen(
			"rbac/failure",
			[]string{},
			[]string{},
			&condition{
				model:  rootDir + "rbac_model.conf",
				policy: rootDir + "rbac_policy",
				sub:    map[string]any{"name": "bob"},
				path:   "/bar/123",
				method: http.MethodPost,
			},
			&action{
				authorized:       false, // bob is not admin
				authorizedReload: true,  // now bob is admin
			},
		),
		gen(
			"v9",
			[]string{},
			[]string{},
			&condition{
				model:  rootDir + "model_v9.conf",
				policy: rootDir + "policy_v9",
				method: http.MethodGet,
			},
			&action{
				authorized:       true,
				authorizedReload: false,
			},
		),
		gen(
			"no policy",
			[]string{},
			[]string{},
			&condition{
				model:  rootDir + "abac_model.conf",
				policy: rootDir + "no_policy",
			},
			&action{
				enforceErr: regexp.MustCompile(`invalid memory address or nil pointer dereference`),
				reLoadErr:  regexp.MustCompile(`(cannot find the file|no such file)`),
			},
		),
		gen(
			"empty file",
			[]string{},
			[]string{},
			&condition{
				model:  rootDir + "abac_model.conf",
				policy: rootDir + "empty",
			},
			&action{
				loadErr: regexp.MustCompile(`EOF`),
			},
		),
		gen(
			"file not found",
			[]string{},
			[]string{},
			&condition{
				model:  rootDir + "abac_model.conf",
				policy: rootDir + "not_found",
			},
			&action{
				loadErr: regexp.MustCompile(`(cannot find the file|no such file)`),
			},
		),
		gen(
			"invalid format",
			[]string{},
			[]string{},
			&condition{
				model:  rootDir + "abac_model.conf",
				policy: rootDir + "invalid",
			},
			&action{
				loadErr: regexp.MustCompile(`EOF`),
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			policy := &xmlAdapter{
				Adapter:  &noopAdapter{},
				filePath: tt.C().policy + ".xml",
			}

			enf, err := casbin.NewEnforcer(tt.C().model, policy)
			t.Logf("%s\n", err)
			if tt.A().loadErr != nil {
				testutil.Diff(t, true, tt.A().loadErr.MatchString(err.Error()))
				testutil.Diff(t, (*casbin.Enforcer)(nil), enf)
				return
			} else {
				testutil.Diff(t, nil, err)
			}
			enf.AddFunction("mapValue", mapValue)

			ok, err := enf.Enforce(tt.C().sub, tt.C().path, tt.C().method)
			testutil.Diff(t, tt.A().authorized, ok)
			t.Logf("%s\n", err)
			if tt.A().enforceErr != nil {
				testutil.Diff(t, true, tt.A().enforceErr.MatchString(err.Error()))
			} else {
				testutil.Diff(t, nil, err)
			}

			policy.filePath = tt.C().policy + "_reload.xml"
			err = enf.LoadPolicy()
			t.Logf("%s\n", err)
			if tt.A().reLoadErr != nil {
				testutil.Diff(t, true, tt.A().reLoadErr.MatchString(err.Error()))
			} else {
				testutil.Diff(t, nil, err)
			}
			ok, _ = enf.Enforce(tt.C().sub, tt.C().path, tt.C().method)
			testutil.Diff(t, tt.A().authorizedReload, ok)
		})
	}
}

type policyRoundTripper struct {
	http.Handler
	mime string
	body io.ReadCloser
	err  error

	called      int
	notModified bool
}

func (rt *policyRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	rt.called += 1
	if rt.err != nil {
		return nil, rt.err
	}

	w := httptest.NewRecorder()
	rt.Handler.ServeHTTP(w, r)

	resp := &http.Response{
		StatusCode: w.Result().StatusCode,
		Header:     w.Result().Header,
		Body:       io.NopCloser(bytes.NewReader(w.Body.Bytes())),
	}
	if rt.called >= 2 && rt.notModified {
		resp.StatusCode = http.StatusNotModified
		resp.Body = io.NopCloser(bytes.NewReader(nil))
		return resp, nil
	}
	if rt.mime != "" {
		resp.Header.Set("Content-Type", rt.mime)
	}
	if rt.body != nil {
		resp.Body = rt.body
	}
	return resp, nil
}

type testErrReader struct {
	io.Reader
	err error
}

func (r *testErrReader) Read(p []byte) (n int, err error) {
	return 0, r.err
}

func TestNewFileAdapter(t *testing.T) {
	type condition struct {
		policy string
	}

	type action struct {
		loadErr *regexp.Regexp
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	rootDir := "../../../test/ut/app/casbin/http/"

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"csv",
			[]string{},
			[]string{},
			&condition{
				policy: rootDir + "abac_policy.csv",
			},
			&action{},
		),
		gen(
			"json",
			[]string{},
			[]string{},
			&condition{
				policy: rootDir + "abac_policy.json",
			},
			&action{},
		),
		gen(
			"yaml",
			[]string{},
			[]string{},
			&condition{
				policy: rootDir + "abac_policy.yaml",
			},
			&action{},
		),
		gen(
			"xml",
			[]string{},
			[]string{},
			&condition{
				policy: rootDir + "abac_policy.xml",
			},
			&action{},
		),
		gen(
			"extension",
			[]string{},
			[]string{},
			&condition{
				policy: rootDir + "abac_policy.txt",
			},
			&action{
				loadErr: regexp.MustCompile(`unsupported file extension`),
			},
		),
		gen(
			"no extension",
			[]string{},
			[]string{},
			&condition{
				policy: rootDir + "abac_policy",
			},
			&action{
				loadErr: regexp.MustCompile(`unsupported file extension`),
			},
		),
		gen(
			"empty file name",
			[]string{},
			[]string{},
			&condition{
				policy: "",
			},
			&action{
				loadErr: regexp.MustCompile(`unsupported file extension`),
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			policy, err := newFileAdapter(tt.C().policy)
			t.Logf("%s\n", err)
			if tt.A().loadErr != nil {
				testutil.Diff(t, true, tt.A().loadErr.MatchString(err.Error()))
				return
			} else {
				testutil.Diff(t, nil, err)
			}

			enf, err := casbin.NewEnforcer(rootDir+"abac_model.conf", policy)
			t.Logf("%s\n", err)
			if tt.A().loadErr != nil {
				testutil.Diff(t, true, tt.A().loadErr.MatchString(err.Error()))
				testutil.Diff(t, (*casbin.Enforcer)(nil), enf)
				return
			} else {
				testutil.Diff(t, nil, err)
			}
			enf.AddFunction("mapValue", mapValue)

			ok, err := enf.Enforce(map[string]any{"age": 21}, "/foo", http.MethodPost)
			testutil.Diff(t, true, ok)
			testutil.Diff(t, nil, err)

			ok, err = enf.Enforce(map[string]any{"age": 5}, "/bar", http.MethodGet)
			testutil.Diff(t, false, ok)
			testutil.Diff(t, nil, err)
		})
	}
}

func TestHttpAdapter(t *testing.T) {
	type condition struct {
		endpoint string
		rt       http.RoundTripper
	}

	type action struct {
		loadErr *regexp.Regexp
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	rootDir := "../../../test/ut/app/casbin/http/"
	handler := http.FileServer(http.Dir("../../../test/ut/app/casbin/http/"))

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"application/csv",
			[]string{},
			[]string{},
			&condition{
				endpoint: "http://localhost/abac_policy.csv",
				rt: &policyRoundTripper{
					Handler: handler,
					mime:    "application/csv",
				},
			},
			&action{},
		),
		gen(
			"text/csv",
			[]string{},
			[]string{},
			&condition{
				endpoint: "http://localhost/abac_policy.csv",
				rt: &policyRoundTripper{
					Handler: handler,
					mime:    "text/csv",
				},
			},
			&action{},
		),
		gen(
			"application/json",
			[]string{},
			[]string{},
			&condition{
				endpoint: "http://localhost/abac_policy.json",
				rt: &policyRoundTripper{
					Handler: handler,
					mime:    "application/json",
				},
			},
			&action{},
		),
		gen(
			"text/json",
			[]string{},
			[]string{},
			&condition{
				endpoint: "http://localhost/abac_policy.json",
				rt: &policyRoundTripper{
					Handler: handler,
					mime:    "text/json",
				},
			},
			&action{},
		),
		gen(
			"application/yaml",
			[]string{},
			[]string{},
			&condition{
				endpoint: "http://localhost/abac_policy.yaml",
				rt: &policyRoundTripper{
					Handler: handler,
					mime:    "application/yaml",
				},
			},
			&action{},
		),
		gen(
			"text/yaml",
			[]string{},
			[]string{},
			&condition{
				endpoint: "http://localhost/abac_policy.yaml",
				rt: &policyRoundTripper{
					Handler: handler,
					mime:    "text/yaml",
				},
			},
			&action{},
		),
		gen(
			"application/yml",
			[]string{},
			[]string{},
			&condition{
				endpoint: "http://localhost/abac_policy.yaml",
				rt: &policyRoundTripper{
					Handler: handler,
					mime:    "application/yml",
				},
			},
			&action{},
		),
		gen(
			"text/yml",
			[]string{},
			[]string{},
			&condition{
				endpoint: "http://localhost/abac_policy.yaml",
				rt: &policyRoundTripper{
					Handler: handler,
					mime:    "text/yml",
				},
			},
			&action{},
		),
		gen(
			"application/xml",
			[]string{},
			[]string{},
			&condition{
				endpoint: "http://localhost/abac_policy.xml",
				rt: &policyRoundTripper{
					Handler: handler,
					mime:    "application/xml",
				},
			},
			&action{},
		),
		gen(
			"text/xml",
			[]string{},
			[]string{},
			&condition{
				endpoint: "http://localhost/abac_policy.xml",
				rt: &policyRoundTripper{
					Handler: handler,
					mime:    "text/xml",
				},
			},
			&action{},
		),
		gen(
			"invalid type",
			[]string{},
			[]string{},
			&condition{
				endpoint: "http://localhost/abac_policy.csv",
				rt: &policyRoundTripper{
					Handler: handler,
					mime:    "text/plain",
				},
			},
			&action{
				loadErr: regexp.MustCompile(`unsupported media type text/plain`),
			},
		),
		gen(
			"not modified",
			[]string{},
			[]string{},
			&condition{
				endpoint: "http://localhost/abac_policy.csv",
				rt: &policyRoundTripper{
					Handler:     handler,
					mime:        "text/csv",
					notModified: true,
				},
			},
			&action{},
		),
		gen(
			"not found",
			[]string{},
			[]string{},
			&condition{
				endpoint: "http://localhost/not_found.csv",
				rt: &policyRoundTripper{
					Handler: handler,
				},
			},
			&action{
				loadErr: regexp.MustCompile(`failed to get policy from`),
			},
		),
		gen(
			"round trip error",
			[]string{},
			[]string{},
			&condition{
				endpoint: "http://localhost/not_found.csv",
				rt: &policyRoundTripper{
					Handler: handler,
					err:     http.ErrAbortHandler, // Dummy error
				},
			},
			&action{
				loadErr: regexp.MustCompile(`abort Handler`),
			},
		),
		gen(
			"body read error",
			[]string{},
			[]string{},
			&condition{
				endpoint: "http://localhost/not_found.csv",
				rt: &policyRoundTripper{
					Handler: handler,
					body:    io.NopCloser(&testErrReader{err: io.ErrUnexpectedEOF}),
				},
			},
			&action{
				loadErr: regexp.MustCompile(`unexpected EOF`),
			},
		),
		gen(
			"invalid endpoint",
			[]string{},
			[]string{},
			&condition{
				endpoint: "invalid*://foobar.com/",
				rt: &policyRoundTripper{
					Handler: handler,
				},
			},
			&action{
				loadErr: regexp.MustCompile(`first path segment in URL cannot contain colon`),
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			policy := &httpAdapter{
				Adapter:  &noopAdapter{},
				endpoint: tt.C().endpoint,
				rt:       tt.C().rt,
			}

			enf, err := casbin.NewEnforcer(rootDir+"abac_model.conf", policy)
			t.Logf("%s\n", err)
			if tt.A().loadErr != nil {
				testutil.Diff(t, true, tt.A().loadErr.MatchString(err.Error()))
				testutil.Diff(t, (*casbin.Enforcer)(nil), enf)
				return
			} else {
				testutil.Diff(t, nil, err)
			}
			enf.AddFunction("mapValue", mapValue)

			// Force reload policies to test 302 NotModified.
			policy.LoadPolicy(enf.GetModel())

			ok, err := enf.Enforce(map[string]any{"age": 21}, "/foo", http.MethodPost)
			testutil.Diff(t, true, ok)
			testutil.Diff(t, nil, err)

			ok, err = enf.Enforce(map[string]any{"age": 5}, "/bar", http.MethodGet)
			testutil.Diff(t, false, ok)
			testutil.Diff(t, nil, err)
		})
	}
}
