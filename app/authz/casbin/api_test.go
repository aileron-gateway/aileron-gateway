// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package casbin

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
	"time"

	v1 "github.com/aileron-gateway/aileron-gateway/apis/app/v1"
	k "github.com/aileron-gateway/aileron-gateway/apis/kernel"
	"github.com/aileron-gateway/aileron-gateway/core"
	"github.com/aileron-gateway/aileron-gateway/kernel/api"
	"github.com/aileron-gateway/aileron-gateway/kernel/cron"
	"github.com/aileron-gateway/aileron-gateway/kernel/er"
	"github.com/aileron-gateway/aileron-gateway/kernel/log"
	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
	"github.com/casbin/casbin/v3/model"
	"github.com/casbin/casbin/v3/persist"
	fileadapter "github.com/casbin/casbin/v3/persist/file-adapter"
	"github.com/google/go-cmp/cmp/cmpopts"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func TestCreate(t *testing.T) {
	type condition struct {
		manifest protoreflect.ProtoMessage
	}

	type action struct {
		err        any // error or errorutil.Kind
		errPattern *regexp.Regexp
		expect     any
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	CndDefaultManifest := tb.Condition("input default manifest", "input default manifest")
	CndErrorReferenceLogSet := tb.Condition("input error reference to logger or log creator", "input error reference to logger or log creator")
	CndErrorErrorHandlerSet := tb.Condition("input error reference to errorhandler", "input error reference to errorhandler")
	ActCheckNoError := tb.Action("check no error was returned", "check no error was returned")
	ActCheckErrorMsg := tb.Action("check error message", "check the error messages that was returned")
	table := tb.Build()

	rootPath := "../../../test/ut/app/casbin/"

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"create with default manifest",
			[]string{CndDefaultManifest},
			[]string{ActCheckNoError},
			&condition{
				manifest: Resource.Default(),
			},
			&action{
				err:    nil,
				expect: &authz{},
			},
		),
		gen(
			"use auth logger",
			[]string{CndErrorReferenceLogSet},
			[]string{ActCheckErrorMsg},
			&condition{
				manifest: &v1.CasbinAuthzMiddleware{
					Metadata: &k.Metadata{},
					Spec: &v1.CasbinAuthzMiddlewareSpec{
						Logger: &k.Reference{
							APIVersion: "container/v1",
							Kind:       "Container",
							Namespace:  "default",
							Name:       "noopLogger",
						},
					},
				},
			},
			&action{
				err: nil,
				expect: &authz{
					lg: log.NoopLogger,
				},
			},
		),
		gen(
			"fail to get logger",
			[]string{CndErrorReferenceLogSet},
			[]string{ActCheckErrorMsg},
			&condition{
				manifest: &v1.CasbinAuthzMiddleware{
					Metadata: &k.Metadata{},
					Spec: &v1.CasbinAuthzMiddlewareSpec{
						Logger: &k.Reference{APIVersion: "wrong"},
					},
				},
			},
			&action{
				err:        core.ErrCoreGenCreateObject,
				errPattern: regexp.MustCompile(core.ErrPrefix + `failed to create CasbinAuthzMiddleware`),
			},
		),
		gen(
			"fail to get errorhandler",
			[]string{CndErrorErrorHandlerSet},
			[]string{ActCheckErrorMsg},
			&condition{
				manifest: &v1.CasbinAuthzMiddleware{
					Metadata: &k.Metadata{},
					Spec: &v1.CasbinAuthzMiddlewareSpec{
						ErrorHandler: &k.Reference{APIVersion: "wrong"},
					},
				},
			},
			&action{
				err:        core.ErrCoreGenCreateObject,
				errPattern: regexp.MustCompile(core.ErrPrefix + `failed to create CasbinAuthzMiddleware`),
			},
		),
		gen(
			"create success CasbinAuthorization",
			[]string{CndDefaultManifest},
			[]string{ActCheckNoError},
			&condition{
				manifest: &v1.CasbinAuthzMiddleware{
					Metadata: &k.Metadata{},
					Spec: &v1.CasbinAuthzMiddlewareSpec{
						Enforcers: []*v1.EnforcerSpec{
							{
								ModelPath: rootPath + "abac_model.conf",
								Policies: &v1.EnforcerSpec_PolicyPath{
									PolicyPath: rootPath + "abac_policy.csv",
								},
							},
						},
					},
				},
			},
			&action{
				err: nil,
			},
		),
		gen(
			"failed to create CasbinAuthorization",
			[]string{CndDefaultManifest},
			[]string{ActCheckNoError},
			&condition{
				manifest: &v1.CasbinAuthzMiddleware{
					Metadata: &k.Metadata{},
					Spec: &v1.CasbinAuthzMiddlewareSpec{
						Enforcers: []*v1.EnforcerSpec{
							{
								ModelPath: "notExist.conf",
								Policies:  &v1.EnforcerSpec_PolicyPath{PolicyPath: "notExist.csv"},
							},
						},
					},
				},
			},
			&action{
				err:        core.ErrCoreGenCreateObject,
				errPattern: regexp.MustCompile(core.ErrPrefix + `failed to create CasbinAuthzMiddleware`),
			},
		),
		gen(
			"enforcer create error",
			[]string{},
			[]string{},
			&condition{
				manifest: &v1.CasbinAuthzMiddleware{
					Metadata: &k.Metadata{},
					Spec: &v1.CasbinAuthzMiddlewareSpec{
						Enforcers: []*v1.EnforcerSpec{
							{
								ModelPath: rootPath + "abac_model.conf",
								Policies:  &v1.EnforcerSpec_PolicyPath{PolicyPath: rootPath + "abac_model.conf"},
								Cron:      "**********",
							},
						},
					},
				},
			},
			&action{
				err:        core.ErrCoreGenCreateObject,
				errPattern: regexp.MustCompile(core.ErrPrefix + `failed to create CasbinAuthzMiddleware`),
			},
		),
	}
	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			server := api.NewContainerAPI()
			postTestResource(server, "noopLogger", log.NoopLogger)

			a := &API{}
			_, err := a.Create(server, tt.C().manifest)
			testutil.DiffError(t, tt.A().err, tt.A().errPattern, err)
		})
	}
}

func TestNewEnforcers(t *testing.T) {
	type condition struct {
		spec *v1.EnforcerSpec
	}

	type action struct {
		err error
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	lg := log.GlobalLogger(log.DefaultLoggerName)
	rootDir := "../../../test/ut/app/casbin/"

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"model from endpoint",
			[]string{},
			[]string{},
			&condition{
				spec: &v1.EnforcerSpec{
					RoundTripper: &k.Reference{
						APIVersion: "container/v1",
						Kind:       "Container",
						Namespace:  "default",
						Name:       "testRoundTripper",
					},
					ModelPath: "http://localhost/abac_model.conf",
					Policies: &v1.EnforcerSpec_PolicyPath{
						PolicyPath: rootDir + "abac_policy.csv",
					},
				},
			},
			&action{},
		),
		gen(
			"invalid model from endpoint",
			[]string{},
			[]string{},
			&condition{
				spec: &v1.EnforcerSpec{
					RoundTripper: &k.Reference{
						APIVersion: "container/v1",
						Kind:       "Container",
						Namespace:  "default",
						Name:       "testRoundTripper",
					},
					ModelPath: "http://localhost/not-exist-model.conf", // Not found model.
					Policies: &v1.EnforcerSpec_PolicyPath{
						PolicyPath: rootDir + "abac_policy.csv",
					},
				},
			},
			&action{
				err: &er.Error{
					Package:     "authz/casbin",
					Type:        "load model",
					Description: "model endpoint returned non 200 OK status.",
				},
			},
		),
		gen(
			"file adapter",
			[]string{},
			[]string{},
			&condition{
				spec: &v1.EnforcerSpec{
					ModelPath: rootDir + "abac_model.conf",
					Policies: &v1.EnforcerSpec_PolicyPath{
						PolicyPath: rootDir + "abac_policy.csv",
					},
				},
			},
			&action{},
		),
		gen(
			"http adapter",
			[]string{},
			[]string{},
			&condition{
				spec: &v1.EnforcerSpec{
					RoundTripper: &k.Reference{
						APIVersion: "container/v1",
						Kind:       "Container",
						Namespace:  "default",
						Name:       "testRoundTripper",
					},
					ModelPath: rootDir + "abac_model.conf",
					Policies: &v1.EnforcerSpec_PolicyURL{
						PolicyURL: "http://localhost/abac_policy.csv",
					},
				},
			},
			&action{},
		),
		gen(
			"external adapter",
			[]string{},
			[]string{},
			&condition{
				spec: &v1.EnforcerSpec{
					ModelPath: rootDir + "abac_model.conf",
					Policies: &v1.EnforcerSpec_ExternalAdapter{
						ExternalAdapter: &k.Reference{
							APIVersion: "container/v1",
							Kind:       "Container",
							Namespace:  "default",
							Name:       "testAdapter",
						},
					},
				},
			},
			&action{},
		),
		gen(
			"ref error/round tripper",
			[]string{},
			[]string{},
			&condition{
				spec: &v1.EnforcerSpec{
					RoundTripper: &k.Reference{
						APIVersion: "container/v1",
						Kind:       "Container",
						Namespace:  "default",
						Name:       "not_exist_roundTripper",
					},
				},
			},
			&action{
				err: &er.Error{
					Package:     api.ErrPkg,
					Type:        api.ErrTypeUtil,
					Description: api.ErrDscAssert,
				},
			},
		),
		gen(
			"ref error/external adapter",
			[]string{},
			[]string{},
			&condition{
				spec: &v1.EnforcerSpec{
					ModelPath: rootDir + "abac_model.conf",
					Policies: &v1.EnforcerSpec_ExternalAdapter{
						ExternalAdapter: &k.Reference{
							APIVersion: "container/v1",
							Kind:       "Container",
							Namespace:  "default",
							Name:       "not_exist_adapter",
						},
					},
				},
			},
			&action{
				err: &er.Error{
					Package:     api.ErrPkg,
					Type:        api.ErrTypeUtil,
					Description: api.ErrDscAssert,
				},
			},
		),
		gen(
			"invalid cron",
			[]string{},
			[]string{},
			&condition{
				spec: &v1.EnforcerSpec{
					Cron:      "*********",
					ModelPath: rootDir + "abac_model.conf",
					Policies: &v1.EnforcerSpec_PolicyPath{
						PolicyPath: rootDir + "abac_policy.csv",
					},
				},
			},
			&action{
				err: &er.Error{
					Package:     cron.ErrPkg,
					Type:        cron.ErrTypeParse,
					Description: cron.ErrDscParse,
				},
			},
		),
	}
	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			// Prepare test resources
			server := api.NewContainerAPI()
			rt := &policyRoundTripper{Handler: http.FileServer(http.Dir(rootDir)), mime: "text/csv"}
			adp := &testAdapter{Adapter: fileadapter.NewAdapter(rootDir + "abac_policy.csv")}
			postTestResource(server, "testRoundTripper", rt)
			postTestResource(server, "testAdapter", adp)

			enfs, err := newEnforcers(server, lg, []*v1.EnforcerSpec{tt.C().spec})
			testutil.Diff(t, tt.A().err, err, cmpopts.EquateErrors())
			if err != nil {
				testutil.Diff(t, 0, len(enfs))
				return
			}

			enf := enfs[0]
			ok, err := enf.Enforce(map[string]any{"age": 21}, "/test", http.MethodGet)
			testutil.Diff(t, true, ok)
			testutil.Diff(t, nil, err)
			ok, err = enf.Enforce(map[string]any{"age": 5}, "/test", http.MethodGet)
			testutil.Diff(t, false, ok)
			testutil.Diff(t, nil, err)
		})
	}
}

type testAdapter struct {
	persist.Adapter
	called    int
	reloadErr error
}

func (a *testAdapter) LoadPolicy(model model.Model) error {
	a.called += 1
	if a.called >= 2 && a.reloadErr != nil {
		return a.reloadErr
	}
	return a.Adapter.LoadPolicy(model)
}

func TestRunReloadCronJob(t *testing.T) {
	type condition struct {
		cron      string
		reloadErr error
	}

	type action struct {
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	lg := log.GlobalLogger(log.DefaultLoggerName)
	rootDir := "../../../test/ut/app/casbin/"

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"reload success",
			[]string{},
			[]string{},
			&condition{
				cron: "* * * * * *", // Every seconds.
			},
			&action{},
		),
		gen(
			"reload failure",
			[]string{},
			[]string{},
			&condition{
				cron:      "* * * * * *", // Every seconds.
				reloadErr: errors.New("reload error"),
			},
			&action{},
		),
	}
	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			// Prepare test resources
			server := api.NewContainerAPI()
			adp := &testAdapter{
				Adapter:   fileadapter.NewAdapter(rootDir + "abac_policy.csv"),
				reloadErr: tt.C().reloadErr,
			}
			postTestResource(server, "testAdapter", adp)

			spec := &v1.EnforcerSpec{
				ModelPath: rootDir + "abac_model.conf",
				Policies: &v1.EnforcerSpec_ExternalAdapter{
					ExternalAdapter: &k.Reference{
						APIVersion: "container/v1",
						Kind:       "Container",
						Namespace:  "default",
						Name:       "testAdapter",
					},
				},
				Cron: tt.C().cron,
			}
			enfs, err := newEnforcers(server, lg, []*v1.EnforcerSpec{spec})
			testutil.Diff(t, nil, err)

			enf := enfs[0]
			ok, err := enf.Enforce(map[string]any{"age": 21}, "/test", http.MethodGet)
			testutil.Diff(t, true, ok)
			testutil.Diff(t, nil, err)
			ok, err = enf.Enforce(map[string]any{"age": 5}, "/test", http.MethodGet)
			testutil.Diff(t, false, ok)
			testutil.Diff(t, nil, err)

			time.Sleep(time.Second)
			testutil.Diff(t, 2, adp.called)
			time.Sleep(time.Second)
			testutil.Diff(t, 3, adp.called)

			// Works correct even after the policy reload.
			ok, err = enf.Enforce(map[string]any{"age": 21}, "/test", http.MethodGet)
			testutil.Diff(t, true, ok)
			testutil.Diff(t, nil, err)
			ok, err = enf.Enforce(map[string]any{"age": 5}, "/test", http.MethodGet)
			testutil.Diff(t, false, ok)
			testutil.Diff(t, nil, err)
		})
	}
}

func postTestResource(server api.API[*api.Request, *api.Response], name string, res any) {
	ref := &k.Reference{
		APIVersion: "container/v1",
		Kind:       "Container",
		Namespace:  "default",
		Name:       name,
	}
	req := &api.Request{
		Method:  api.MethodPost,
		Key:     ref.APIVersion + "/" + ref.Kind + "/" + ref.Namespace + "/" + ref.Name,
		Content: res,
	}
	if _, err := server.Serve(context.Background(), req); err != nil {
		panic(err)
	}
}

type modelRoundTripper struct {
	http.Handler
	err  error
	body io.ReadCloser

	called  int
	rHeader http.Header
}

func (rt *modelRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	rt.called += 1
	rt.rHeader = r.Header
	w := httptest.NewRecorder()
	rt.Handler.ServeHTTP(w, r)
	resp := &http.Response{
		StatusCode: w.Result().StatusCode,
		Header:     w.Result().Header,
		Body:       io.NopCloser(bytes.NewReader(w.Body.Bytes())),
	}
	if rt.body != nil {
		resp.Body = rt.body
	}
	return resp, rt.err
}

func TestLoadModel(t *testing.T) {
	type condition struct {
		path   string
		header map[string]string
		rt     *modelRoundTripper
	}

	type action struct {
		err error
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	rootDir := "../../../test/ut/app/casbin/"

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"model from file",
			[]string{},
			[]string{},
			&condition{
				path: rootDir + "abac_policy.csv",
			},
			&action{},
		),
		gen(
			"model from endpoint",
			[]string{},
			[]string{},
			&condition{
				path: "http://localhost/abac_model.conf",
				rt:   &modelRoundTripper{Handler: http.FileServer(http.Dir(rootDir))},
			},
			&action{},
		),
		gen(
			"invalid path",
			[]string{},
			[]string{},
			&condition{
				path: "http://localhost/\n\n",
				rt:   &modelRoundTripper{Handler: http.FileServer(http.Dir(rootDir))},
			},
			&action{
				err: &er.Error{
					Package:     "authz/casbin",
					Type:        "load model",
					Description: "failed create request.",
				},
			},
		),
		gen(
			"round trip error",
			[]string{},
			[]string{},
			&condition{
				path: "http://localhost/abac_model.conf",
				rt: &modelRoundTripper{
					Handler: http.FileServer(http.Dir(rootDir)),
					err:     io.ErrUnexpectedEOF,
				},
			},
			&action{
				err: &er.Error{
					Package:     "authz/casbin",
					Type:        "load model",
					Description: "failed to get model from endpoint.",
				},
			},
		),
		gen(
			"non 200 OK",
			[]string{},
			[]string{},
			&condition{
				path: "http://localhost/not-exist-model.conf",
				rt:   &modelRoundTripper{Handler: http.FileServer(http.Dir(rootDir))},
			},
			&action{
				err: &er.Error{
					Package:     "authz/casbin",
					Type:        "load model",
					Description: "model endpoint returned non 200 OK status.",
				},
			},
		),
		gen(
			"body read error",
			[]string{},
			[]string{},
			&condition{
				path: "http://localhost/abac_model.conf",
				rt: &modelRoundTripper{
					Handler: http.FileServer(http.Dir(rootDir)),
					body:    io.NopCloser(&testutil.ErrorReader{}),
				},
			},
			&action{
				err: &er.Error{
					Package:     "authz/casbin",
					Type:        "load model",
					Description: "failed to read model from response body.",
				},
			},
		),
		gen(
			"model parse error",
			[]string{},
			[]string{},
			&condition{
				path: "http://localhost/abac_policy.csv",
				rt:   &modelRoundTripper{Handler: http.FileServer(http.Dir(rootDir))},
			},
			&action{
				err: &er.Error{
					Package:     "authz/casbin",
					Type:        "load model",
					Description: "failed to load model.",
				},
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			hdr := map[string]string{"foo": "bar"}
			rt := addHeader(hdr).Tripperware(tt.C().rt)

			_, err := loadModel(tt.C().path, rt)
			testutil.Diff(t, tt.A().err, err, cmpopts.EquateErrors())
			if tt.C().rt != nil && tt.C().rt.called > 0 {
				testutil.Diff(t, "bar", tt.C().rt.rHeader.Get("foo"))
			}
		})
	}
}
