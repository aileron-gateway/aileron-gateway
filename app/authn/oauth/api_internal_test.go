// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package oauth

import (
	"encoding/base64"
	"regexp"
	"testing"

	v1 "github.com/aileron-gateway/aileron-gateway/apis/app/v1"
	k "github.com/aileron-gateway/aileron-gateway/apis/kernel"
	"github.com/aileron-gateway/aileron-gateway/core"
	"github.com/aileron-gateway/aileron-gateway/kernel/api"
	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/testing/protocmp"
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
	CndDefaultManifest := tb.Condition("input default manifest", "input default manifest")
	CndMutateBaseContext := tb.Condition("mutate base context", "mutate base context")
	CndMutateContext := tb.Condition("mutate context", "mutate context")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"apply default values",
			[]string{CndDefaultManifest},
			[]string{},
			&condition{
				manifest: Resource.Default(),
			},
			&action{
				manifest: &v1.OAuthAuthenticationHandler{
					APIVersion: apiVersion,
					Kind:       kind,
					Metadata: &k.Metadata{
						Namespace: "default",
						Name:      "default",
					},
					Spec: &v1.OAuthAuthenticationHandlerSpec{},
				},
			},
		),
		gen(
			"mutate base context",
			[]string{CndMutateBaseContext},
			[]string{},
			&condition{
				manifest: &v1.OAuthAuthenticationHandler{
					APIVersion: apiVersion,
					Kind:       kind,
					Metadata: &k.Metadata{
						Name: "default",
					},
					Spec: &v1.OAuthAuthenticationHandlerSpec{
						Contexts: []*v1.Context{{}},
					},
				},
			},
			&action{
				manifest: &v1.OAuthAuthenticationHandler{
					APIVersion: apiVersion,
					Kind:       kind,
					Metadata: &k.Metadata{
						Name: "default",
					},
					Spec: &v1.OAuthAuthenticationHandlerSpec{
						Contexts: []*v1.Context{
							{
								Name: "default",
								Provider: &v1.OAuthProvider{
									Endpoints: &v1.ProviderEndpoints{},
								},
								Client:            &v1.OAuthClient{},
								TokenRedeemer:     &v1.ClientRequester{},
								TokenIntrospector: &v1.ClientRequester{},
								JWTHandler:        &v1.JWTHandlerSpec{},
								ClaimsKey:         "AuthnClaims",
								IDTProxyHeader:    "",
								ATProxyHeader:     "",
								ATValidation: &v1.TokenValidation{
									Leeway: 5,
								},
								IDTValidation: &v1.TokenValidation{
									Leeway: 5,
								},
							},
						},
					},
				},
			},
		),
		gen(
			"mutate context",
			[]string{CndMutateContext},
			[]string{},
			&condition{
				manifest: &v1.OAuthAuthenticationHandler{
					APIVersion: apiVersion,
					Kind:       kind,
					Metadata: &k.Metadata{
						Name: "default",
					},
					Spec: &v1.OAuthAuthenticationHandlerSpec{
						Contexts: []*v1.Context{
							{
								Name:           "test-name",
								ClaimsKey:      "test-claims-key",
								ATProxyHeader:  "test-at-proxy-header",
								IDTProxyHeader: "test-idt-proxy-header",
								ATValidation: &v1.TokenValidation{
									Iss: "at-iss",
									Aud: "at-aud",
								},
								IDTValidation: &v1.TokenValidation{
									Iss: "idt-iss",
									Aud: "idt-aud",
								},
							},
						},
					},
				},
			},
			&action{
				manifest: &v1.OAuthAuthenticationHandler{
					APIVersion: apiVersion,
					Kind:       kind,
					Metadata: &k.Metadata{
						Name: "default",
					},
					Spec: &v1.OAuthAuthenticationHandlerSpec{
						Contexts: []*v1.Context{
							{
								Name: "test-name",
								Provider: &v1.OAuthProvider{
									Endpoints: &v1.ProviderEndpoints{},
								},
								Client:            &v1.OAuthClient{},
								TokenRedeemer:     &v1.ClientRequester{},
								TokenIntrospector: &v1.ClientRequester{},
								JWTHandler:        &v1.JWTHandlerSpec{},
								ClaimsKey:         "test-claims-key",
								IDTProxyHeader:    "test-idt-proxy-header",
								ATProxyHeader:     "test-at-proxy-header",
								ATValidation: &v1.TokenValidation{
									Iss:    "at-iss",
									Aud:    "at-aud",
									Leeway: 5,
								},
								IDTValidation: &v1.TokenValidation{
									Iss:    "idt-iss",
									Aud:    "idt-aud",
									Leeway: 5,
								},
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
			msg := Resource.Mutate(tt.C().manifest)

			opts := []cmp.Option{
				protocmp.Transform(),
				cmpopts.IgnoreUnexported(k.Metadata{}, k.Status{}),
			}
			testutil.Diff(t, tt.A().manifest, msg, opts...)
		})
	}
}

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
	CndErrorGetOAuthContext := tb.Condition("input error get OAuthContext", "input error get OAuthContext")
	CndSetHeaderKey := tb.Condition("input HeaderKey", "input HeaderKey")
	CndSuccessToGetOAuthContext := tb.Condition("input OAuthContext", "input OAuthContext")
	ActCheckNoError := tb.Action("check no error was returned", "check no error was returned")
	ActCheckErrorMsg := tb.Action("check error message", "check the error messages that was returned")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"create with default manifest",
			[]string{CndDefaultManifest},
			[]string{ActCheckNoError},
			&condition{
				manifest: &v1.OAuthAuthenticationHandler{
					Metadata: &k.Metadata{},
					Spec: &v1.OAuthAuthenticationHandlerSpec{
						Handlers: &v1.OAuthAuthenticationHandlerSpec_ResourceServerHandler{
							ResourceServerHandler: &v1.ResourceServerHandler{},
						},
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
				manifest: &v1.OAuthAuthenticationHandler{
					Metadata: &k.Metadata{},
					Spec: &v1.OAuthAuthenticationHandlerSpec{
						Handlers: &v1.OAuthAuthenticationHandlerSpec_ResourceServerHandler{
							ResourceServerHandler: &v1.ResourceServerHandler{},
						},
						Logger: &k.Reference{
							APIVersion: "wrong",
						},
					},
				},
			},
			&action{
				err:        core.ErrCoreGenCreateObject,
				errPattern: regexp.MustCompile(core.ErrPrefix + `failed to create OAuthAuthenticationHandler`),
			},
		),
		gen(
			"fail to get OAuthContext",
			[]string{CndErrorGetOAuthContext},
			[]string{ActCheckErrorMsg},
			&condition{
				manifest: &v1.OAuthAuthenticationHandler{
					Metadata: &k.Metadata{},
					Spec: &v1.OAuthAuthenticationHandlerSpec{
						Handlers: &v1.OAuthAuthenticationHandlerSpec_ResourceServerHandler{
							ResourceServerHandler: &v1.ResourceServerHandler{},
						},
						Contexts: []*v1.Context{
							{
								Client: &v1.OAuthClient{
									ID:       "test-id",
									Secret:   "test-secret",
									Audience: "test-audience",
									Scopes:   []string{"test-scope"},
								},
								Provider: &v1.OAuthProvider{
									BaseURL: "wrong-base-URL",
									Endpoints: &v1.ProviderEndpoints{
										Authorization: "wrong-authorization",
									},
								},
								ATValidation:  &v1.TokenValidation{},
								IDTValidation: &v1.TokenValidation{},
							},
						},
					},
				},
			},
			&action{
				err:        core.ErrCoreGenCreateObject,
				errPattern: regexp.MustCompile(core.ErrPrefix + `failed to create OAuthAuthenticationHandler`),
			},
		),
		gen(
			"input HeaderKey",
			[]string{CndSetHeaderKey},
			[]string{ActCheckNoError},
			&condition{
				manifest: &v1.OAuthAuthenticationHandler{
					Metadata: &k.Metadata{},
					Spec: &v1.OAuthAuthenticationHandlerSpec{
						Handlers: &v1.OAuthAuthenticationHandlerSpec_ResourceServerHandler{
							ResourceServerHandler: &v1.ResourceServerHandler{
								HeaderKey: "test-HeaderKey",
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
			"success to get OAuthContext",
			[]string{CndSuccessToGetOAuthContext},
			[]string{ActCheckErrorMsg},
			&condition{
				manifest: &v1.OAuthAuthenticationHandler{
					Metadata: &k.Metadata{},
					Spec: &v1.OAuthAuthenticationHandlerSpec{
						Handlers: &v1.OAuthAuthenticationHandlerSpec_ResourceServerHandler{
							ResourceServerHandler: &v1.ResourceServerHandler{},
						},
						Contexts: []*v1.Context{
							{
								Client: &v1.OAuthClient{
									ID:       "test-id",
									Secret:   "test-secret",
									Audience: "test-audience",
									Scopes:   []string{"test-scope"},
								},
								JWTHandler: &v1.JWTHandlerSpec{
									PrivateKeys: []*v1.SigningKeySpec{
										{
											KeyID:     "test",
											Algorithm: v1.SigningKeyAlgorithm_HS256,
											KeyType:   v1.SigningKeyType_COMMON,
											KeyString: base64.StdEncoding.EncodeToString([]byte("password")),
										},
									},
									PublicKeys: []*v1.SigningKeySpec{
										{
											KeyID:     "test",
											Algorithm: v1.SigningKeyAlgorithm_HS256,
											KeyType:   v1.SigningKeyType_COMMON,
											KeyString: base64.StdEncoding.EncodeToString([]byte("password")),
										},
									},
								},
								Provider: &v1.OAuthProvider{
									BaseURL: "wrong-base-URL",
									Endpoints: &v1.ProviderEndpoints{
										Authorization: "wrong-authorization",
									},
								},
								TokenRedeemer: &v1.ClientRequester{
									RoundTripper: nil,
								},
								TokenIntrospector: &v1.ClientRequester{
									RoundTripper: nil,
								},
								ATValidation:  &v1.TokenValidation{},
								IDTValidation: &v1.TokenValidation{},
							},
						},
					},
				},
			},
			&action{
				err: nil,
			},
		),
	}
	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			server := api.NewContainerAPI()
			a := &API{&api.BaseResource{}}
			_, err := a.Create(server, tt.C().manifest)
			testutil.DiffError(t, tt.A().err, tt.A().errPattern, err)
		})
	}
}
