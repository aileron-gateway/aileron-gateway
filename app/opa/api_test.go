// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package opa

import (
	"regexp"
	"testing"

	v1 "github.com/aileron-gateway/aileron-gateway/apis/app/v1"
	k "github.com/aileron-gateway/aileron-gateway/apis/kernel"
	"github.com/aileron-gateway/aileron-gateway/core"
	"github.com/aileron-gateway/aileron-gateway/kernel/api"
	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func TestCreate(t *testing.T) {
	type condition struct {
		manifest protoreflect.ProtoMessage
	}

	type action struct {
		err        any // error or errorutil.Kind
		errPattern *regexp.Regexp
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	CndDefaultManifest := tb.Condition("input default manifest", "input default manifest")
	CndErrorReferenceLogSet := tb.Condition("input error reference to logger or log creator", "input error reference to logger or log creator")
	CndErrorErrorHandlerSet := tb.Condition("input error reference to errorhandler", "input error reference to errorhandler")
	CndSetRegoConfig := tb.Condition("input rego config", "input rego config")
	CndInputInvalid := tb.Condition("input manifest with invalid value", "input wrong type of manifest as an argument")
	CndInvalidRegoPolicyContent := tb.Condition("create with invalid rego policy content", "create with invalid rego policy content")
	ActCheckNoError := tb.Action("check no error was returned", "check no error was returned")
	ActCheckErrorMsg := tb.Action("check error message", "check the error messages that was returned")

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"create with default manifest",
			[]string{CndDefaultManifest},
			[]string{ActCheckNoError},
			&condition{
				manifest: &v1.OPAAuthzMiddleware{
					Metadata: &k.Metadata{},
					Spec: &v1.OPAAuthzMiddlewareSpec{
						ClaimsKey: "AuthnClaims",
					},
				},
			},
			&action{
				err: nil,
			},
		),
		gen(
			"fail to get logger",
			[]string{CndErrorReferenceLogSet},
			[]string{ActCheckErrorMsg},
			&condition{
				manifest: &v1.OPAAuthzMiddleware{
					Metadata: &k.Metadata{},
					Spec: &v1.OPAAuthzMiddlewareSpec{
						ClaimsKey: "AuthnClaims",
						Logger: &k.Reference{
							APIVersion: "wrong",
						},
					},
				},
			},
			&action{
				err:        core.ErrCoreGenCreateObject,
				errPattern: regexp.MustCompile(core.ErrPrefix + `failed to create OPAAuthzMiddleware`),
			},
		),
		gen(
			"fail to get errorhandler",
			[]string{CndErrorErrorHandlerSet},
			[]string{ActCheckErrorMsg},
			&condition{
				manifest: &v1.OPAAuthzMiddleware{
					Metadata: &k.Metadata{},
					Spec: &v1.OPAAuthzMiddlewareSpec{
						ClaimsKey: "AuthnClaims",
						ErrorHandler: &k.Reference{
							APIVersion: "wrong",
						},
					},
				},
			},
			&action{
				err:        core.ErrCoreGenCreateObject,
				errPattern: regexp.MustCompile(core.ErrPrefix + `failed to create OPAAuthzMiddleware`),
			},
		),
		gen(
			"create with rego config",
			[]string{CndSetRegoConfig},
			[]string{ActCheckNoError},
			&condition{
				manifest: &v1.OPAAuthzMiddleware{
					Metadata: &k.Metadata{},
					Spec: &v1.OPAAuthzMiddlewareSpec{
						ClaimsKey: "AuthnClaims",
						Regos: []*v1.RegoSpec{
							{
								QueryParameter: "data.example.authz.allow",
								PolicyFiles: []string{
									"../../../test/ut/app/opa/policy.rego",
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
			"create with invalid rego config",
			[]string{CndInputInvalid},
			[]string{ActCheckErrorMsg},
			&condition{
				manifest: &v1.OPAAuthzMiddleware{
					Metadata: &k.Metadata{},
					Spec: &v1.OPAAuthzMiddlewareSpec{
						ClaimsKey: "AuthnClaims",
						Regos: []*v1.RegoSpec{
							{
								QueryParameter: "data.example.authz.allow",
								PolicyFiles: []string{
									"../../../test/ut/app/opa/notexist_policy.rego",
								},
							},
						},
					},
				},
			},
			&action{
				err:        core.ErrCoreGenCreateObject,
				errPattern: regexp.MustCompile(core.ErrPrefix + `failed to create OPAAuthzMiddleware`),
			},
		),
		gen(
			"create with invalid rego policy content",
			[]string{CndInvalidRegoPolicyContent},
			[]string{ActCheckErrorMsg},
			&condition{
				manifest: &v1.OPAAuthzMiddleware{
					Metadata: &k.Metadata{},
					Spec: &v1.OPAAuthzMiddlewareSpec{
						ClaimsKey: "AuthnClaims",
						Regos: []*v1.RegoSpec{
							{
								QueryParameter: "data.example.authz.allow",
								PolicyFiles: []string{
									"../../../test/ut/app/opa/invalid_policy.rego",
								},
							},
						},
					},
				},
			},
			&action{
				err:        core.ErrCoreGenCreateObject,
				errPattern: regexp.MustCompile(core.ErrPrefix + `failed to create OPAAuthzMiddleware`),
			},
		),
	}
	table := tb.Build()

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			server := api.NewContainerAPI()
			a := &API{BaseResource: &api.BaseResource{}}
			_, err := a.Create(server, tt.C().manifest)
			testutil.DiffError(t, tt.A().err, tt.A().errPattern, err)
		})
	}
}

// func TestRegoQueries(t *testing.T) {

// 	type condition struct {
// 		spec  *v1.RegoSpec
// 		input rego.EvalOption
// 	}

// 	type action struct {
// 		err        any // error or errorutil.Kind
// 		errPattern *regexp.Regexp
// 	}

// 	tb := testutil.NewTableBuilder[*condition, *action]()
// 	tb.Name(t.Name())

// 	gen := testutil.NewCase[*condition, *action]
// 	testCases := []*testutil.Case[*condition, *action]{
// 		gen(
// 			"",
// 			[]string{},
// 			[]string{},
// 			&condition{
// 				spec: &v1.RegoSpec{
// 					// PolicyFiles: "",
// 				},
// 			},
// 			&action{
// 				err: nil,
// 			},
// 		),
// 	}
// 	table := tb.Build()

// 	testutil.Register(table, testCases...)
//
// 	for _, tt := range table.Entries() {
// 		tt := tt
// 		t.Run(tt.Name(), func(t *testing.T) {

// 			_, err := regoQueries(tt.C().spec)
// 			testutil.DiffError(t, tt.A().err, tt.A().errPattern, err)

// 		})
// 	}

// }

func TestEnvData(t *testing.T) {
	type condition struct {
		spec *v1.EnvDataSpec
	}

	type action struct {
		expect map[string]any
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"",
			[]string{},
			[]string{},
			&condition{
				spec: &v1.EnvDataSpec{
					// PolicyFiles: "",
				},
			},
			&action{
				expect: map[string]any{},
			},
		),
	}
	table := tb.Build()

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			m := envData(tt.C().spec)
			testutil.Diff(t, tt.A().expect, m)
		})
	}
}
