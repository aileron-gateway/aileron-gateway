// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package errhandler

import (
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
	actCheckError := tb.Action("check the returned error", "check that the returned error is the one expected")
	actCheckNoError := tb.Action("check no error", "check that there is no error returned")
	table := tb.Build()

	lg := log.GlobalLogger(log.DefaultLoggerName)
	em, _ := utilhttp.NewErrorMessage(&v1.ErrorMessageSpec{Codes: []string{"E999"},
		MIMEContents: []*v1.MIMEContentSpec{{MIMEType: "text/plain", Template: "hello"}},
	})

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
				expect: &utilhttp.DefaultErrorHandler{LG: lg, Msgs: []*utilhttp.ErrorMessage{}},
				err:    nil,
			},
		),
		gen(
			"create nil content",
			[]string{cndDefaultManifest},
			[]string{actCheckNoError},
			&condition{
				manifest: &v1.ErrorHandler{
					Metadata: &k.Metadata{},
					Spec: &v1.ErrorHandlerSpec{
						ErrorMessages: []*v1.ErrorMessageSpec{
							{
								Codes:        []string{"E999"},
								MIMEContents: nil,
							},
						},
					},
				},
				server: api.NewContainerAPI(),
			},
			&action{
				expect: &utilhttp.DefaultErrorHandler{LG: lg, Msgs: []*utilhttp.ErrorMessage{}},
				err:    nil,
			},
		),
		gen(
			"create with a content",
			[]string{cndDefaultManifest},
			[]string{actCheckNoError},
			&condition{
				manifest: &v1.ErrorHandler{
					Metadata: &k.Metadata{},
					Spec: &v1.ErrorHandlerSpec{
						ErrorMessages: []*v1.ErrorMessageSpec{
							{
								Codes: []string{"E999"},
								MIMEContents: []*v1.MIMEContentSpec{
									{MIMEType: "text/plain", Template: "hello"},
								},
							},
						},
					},
				},
				server: api.NewContainerAPI(),
			},
			&action{
				expect: &utilhttp.DefaultErrorHandler{LG: lg, Msgs: []*utilhttp.ErrorMessage{em}},
				err:    nil,
			},
		),
		gen(
			"error message create error",
			[]string{},
			[]string{actCheckError},
			&condition{
				manifest: &v1.ErrorHandler{
					Metadata: &k.Metadata{},
					Spec: &v1.ErrorHandlerSpec{
						ErrorMessages: []*v1.ErrorMessageSpec{
							{
								Codes:        []string{"E999"},
								Messages:     []string{"[0-9a-"},
								MIMEContents: []*v1.MIMEContentSpec{{}},
							},
						},
					},
				},
				server: api.NewContainerAPI(),
			},
			&action{
				expect:     nil,
				err:        core.ErrCoreGenCreateObject,
				errPattern: regexp.MustCompile(core.ErrPrefix + `failed to create ErrorHandler`),
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
				cmp.AllowUnexported(utilhttp.ErrorMessage{}),
				cmpopts.IgnoreInterfaces(struct{ txtutil.Template }{}),
			}
			testutil.Diff(t, tt.A().expect, got, opts...)
		})
	}
}
