// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package template

import (
	"net/http"
	"regexp"
	"testing"

	v1 "github.com/aileron-gateway/aileron-gateway/apis/core/v1"
	k "github.com/aileron-gateway/aileron-gateway/apis/kernel"
	"github.com/aileron-gateway/aileron-gateway/core"
	"github.com/aileron-gateway/aileron-gateway/kernel/api"
	"github.com/aileron-gateway/aileron-gateway/kernel/log"
	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
	"github.com/aileron-gateway/aileron-gateway/kernel/txtutil"
	utilhttp "github.com/aileron-gateway/aileron-gateway/util/http"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func TestMutate(t *testing.T) {
	type condition struct {
		manifest protoreflect.ProtoMessage
	}

	type action struct {
		manifest protoreflect.ProtoMessage
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	cndDefault := tb.Condition("default", "input default manifest")
	cndEmptyMIME := tb.Condition("empty mime", "input empty string as mime type")
	cndWrongStatus := tb.Condition("wrong status code", "status code is invalid")
	actCheckMutated := tb.Action("check mutated", "check that the intended fields are mutated")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"mutate default",
			[]string{cndDefault},
			[]string{actCheckMutated},
			&condition{
				manifest: Resource.Default(),
			},
			&action{
				manifest: &v1.TemplateHandler{
					APIVersion: apiVersion,
					Kind:       kind,
					Metadata: &k.Metadata{
						Namespace: "default",
						Name:      "default",
					},
					Spec: &v1.TemplateHandlerSpec{},
				},
			},
		),
		gen(
			"mutate mime type",
			[]string{cndEmptyMIME},
			[]string{actCheckMutated},
			&condition{
				manifest: &v1.TemplateHandler{
					Spec: &v1.TemplateHandlerSpec{
						MIMEContents: []*v1.MIMEContentSpec{
							{
								MIMEType:   "",
								StatusCode: 500,
							},
						},
					},
				},
			},
			&action{
				manifest: &v1.TemplateHandler{
					Spec: &v1.TemplateHandlerSpec{
						MIMEContents: []*v1.MIMEContentSpec{
							{
								MIMEType:   "text/plain",
								StatusCode: 500,
							},
						},
					},
				},
			},
		),
		gen(
			"mutate status code",
			[]string{cndWrongStatus},
			[]string{actCheckMutated},
			&condition{
				manifest: &v1.TemplateHandler{
					Spec: &v1.TemplateHandlerSpec{
						MIMEContents: []*v1.MIMEContentSpec{
							{
								MIMEType:   "application/json",
								StatusCode: 0,
							},
						},
					},
				},
			},
			&action{
				manifest: &v1.TemplateHandler{
					Spec: &v1.TemplateHandlerSpec{
						MIMEContents: []*v1.MIMEContentSpec{
							{
								MIMEType:   "application/json",
								StatusCode: 200,
							},
						},
					},
				},
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			manifest := Resource.Mutate(tt.C().manifest)

			opts := []cmp.Option{
				cmpopts.IgnoreUnexported(k.Metadata{}, k.Reference{}),
				cmpopts.IgnoreUnexported(v1.TemplateHandler{}, v1.TemplateHandlerSpec{}, v1.MIMEContentSpec{}),
			}
			testutil.Diff(t, tt.A().manifest, manifest, opts...)
		})
	}
}

func TestCreate(t *testing.T) {
	type condition struct {
		manifest protoreflect.ProtoMessage
		server   api.API[*api.Request, *api.Response]
	}

	type action struct {
		expect     any
		err        any // error or errorutil.Kind
		errPattern *regexp.Regexp
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	cndDefaultManifest := tb.Condition("input default manifest", "input default manifest")
	cndErrorReference := tb.Condition("input error reference", "input an error reference to an object")
	actCheckError := tb.Action("check the returned error", "check that the returned error is the one expected")
	actCheckNoError := tb.Action("check no error", "check that there is no error returned")
	table := tb.Build()

	tpl, _ := txtutil.NewTemplate(txtutil.TplGoText, "{{.test}}")

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"create with default manifest",
			[]string{cndDefaultManifest},
			[]string{actCheckNoError},
			&condition{
				manifest: Resource.Default(),
				server:   api.NewContainerAPI(),
			},
			&action{
				expect: &templateHandler{
					HandlerBase: &utilhttp.HandlerBase{},
					eh:          utilhttp.GlobalErrorHandler(utilhttp.DefaultErrorHandlerName),
					contents:    []*utilhttp.MIMEContent{},
				},
				err: nil,
			},
		),
		gen(
			"create with mime content",
			[]string{},
			[]string{actCheckNoError},
			&condition{
				manifest: &v1.TemplateHandler{
					Metadata: &k.Metadata{},
					Spec: &v1.TemplateHandlerSpec{
						MIMEContents: []*v1.MIMEContentSpec{
							{
								MIMEType:     "text/plain",
								StatusCode:   http.StatusOK,
								Header:       map[string]string{"foo": "bar"},
								TemplateType: v1.TemplateType_GoText,
								Template:     `{{.test}}`,
							},
						},
					},
				},
				server: api.NewContainerAPI(),
			},
			&action{
				expect: &templateHandler{
					HandlerBase: &utilhttp.HandlerBase{},
					eh:          utilhttp.GlobalErrorHandler(utilhttp.DefaultErrorHandlerName),
					contents: []*utilhttp.MIMEContent{
						{
							Template:   tpl,
							StatusCode: http.StatusOK,
							Header:     http.Header{"Foo": []string{"bar"}},
							MIMEType:   "text/plain",
						},
					},
				},
				err: nil,
			},
		),
		gen(
			"fail to get mime content",
			[]string{cndErrorReference},
			[]string{actCheckError},
			&condition{
				manifest: &v1.TemplateHandler{
					Metadata: &k.Metadata{},
					Spec: &v1.TemplateHandlerSpec{
						MIMEContents: []*v1.MIMEContentSpec{
							{
								TemplateType: v1.TemplateType_GoText,
								Template:     `[0-9a-z`,
							},
						},
					},
				},
				server: api.NewContainerAPI(),
			},
			&action{
				expect:     nil,
				err:        core.ErrCoreGenCreateObject,
				errPattern: regexp.MustCompile(core.ErrPrefix + `failed to create TemplateHandler`),
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			got, err := Resource.Create(tt.C().server, tt.C().manifest)
			testutil.DiffError(t, tt.A().err, tt.A().errPattern, err)

			opts := []cmp.Option{
				cmp.Comparer(testutil.ComparePointer[log.Logger]),
				cmp.AllowUnexported(utilhttp.DefaultErrorHandler{}),
				cmp.AllowUnexported(templateHandler{}, utilhttp.MIMEContent{}),
				cmpopts.IgnoreInterfaces(struct{ txtutil.Template }{}),
			}

			testutil.Diff(t, tt.A().expect, got, opts...)
		})
	}
}
