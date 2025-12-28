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

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"create with default manifest",
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
			&condition{
				manifest: &v1.OPAAuthzMiddleware{
					Metadata: &k.Metadata{},
					Spec: &v1.OPAAuthzMiddlewareSpec{
						ClaimsKey: "AuthnClaims",
						Regos: []*v1.RegoSpec{
							{
								QueryParameter: "data.example.authz.allow",
								PolicyFiles: []string{
									"../../test/ut/app/opa/policy.rego",
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

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			server := api.NewContainerAPI()
			a := &API{BaseResource: &api.BaseResource{}}
			_, err := a.Create(server, tt.C.manifest)
			testutil.DiffError(t, tt.A.err, tt.A.errPattern, err)
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

//
//

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
//

//
//
// 	for _, tt := range testCases  {
// 		tt := tt
// 		t.Run(tt.Name, func(t *testing.T) {

// 			_, err := regoQueries(tt.C.spec)
// 			testutil.DiffError(t, tt.A.err, tt.A.errPattern, err)

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

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"",
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

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			m := envData(tt.C.spec)
			testutil.Diff(t, tt.A.expect, m)
		})
	}
}
