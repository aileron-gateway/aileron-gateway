// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package csrf

import (
	"crypto/sha512"
	"encoding/base64"
	"os"
	"regexp"
	"testing"

	v1 "github.com/aileron-gateway/aileron-gateway/apis/app/v1"
	"github.com/aileron-gateway/aileron-gateway/apis/kernel"
	"github.com/aileron-gateway/aileron-gateway/core"
	"github.com/aileron-gateway/aileron-gateway/internal/hash"
	"github.com/aileron-gateway/aileron-gateway/internal/testutil"
	"github.com/aileron-gateway/aileron-gateway/kernel/api"
	utilhttp "github.com/aileron-gateway/aileron-gateway/util/http"
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

	host, _ := os.Hostname()
	hmacSecret := sha512.Sum512([]byte(host)) // Create default secret.

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"default values",
			&condition{
				manifest: Resource.Default(),
			},
			&action{
				manifest: &v1.CSRFMiddleware{
					APIVersion: apiVersion,
					Kind:       kind,
					Metadata: &kernel.Metadata{
						Namespace: "default",
						Name:      "default",
					},
					Spec: &v1.CSRFMiddlewareSpec{
						Patterns: []string{"/token"},
						SeedSize: 20,
						Secret:   base64.StdEncoding.EncodeToString(hmacSecret[:]),
						HashAlg:  kernel.HashAlg_SHA256,
						CSRFPatterns: &v1.CSRFMiddlewareSpec_CustomRequestHeader{
							CustomRequestHeader: &v1.CustomRequestHeaderSpec{
								HeaderName: "X-Requested-With",
							},
						},
					},
				},
			},
		),
		gen(
			"custom header/overwrite",
			&condition{
				manifest: &v1.CSRFMiddleware{
					APIVersion: apiVersion,
					Kind:       kind,
					Metadata: &kernel.Metadata{
						Namespace: "default",
						Name:      "default",
					},
					Spec: &v1.CSRFMiddlewareSpec{
						SeedSize: 20,
						Secret:   base64.StdEncoding.EncodeToString(hmacSecret[:]),
						HashAlg:  kernel.HashAlg_SHA256,
						CSRFPatterns: &v1.CSRFMiddlewareSpec_CustomRequestHeader{
							CustomRequestHeader: &v1.CustomRequestHeaderSpec{
								HeaderName:     "__foo",
								AllowedPattern: ".+",
							},
						},
					},
				},
			},
			&action{
				manifest: &v1.CSRFMiddleware{
					APIVersion: apiVersion,
					Kind:       kind,
					Metadata: &kernel.Metadata{
						Namespace: "default",
						Name:      "default",
					},
					Spec: &v1.CSRFMiddlewareSpec{
						Patterns: []string{"/token"},
						SeedSize: 20,
						Secret:   base64.StdEncoding.EncodeToString(hmacSecret[:]),
						HashAlg:  kernel.HashAlg_SHA256,
						CSRFPatterns: &v1.CSRFMiddlewareSpec_CustomRequestHeader{
							CustomRequestHeader: &v1.CustomRequestHeaderSpec{
								HeaderName:     "__foo",
								AllowedPattern: ".+",
							},
						},
					},
				},
			},
		),
		gen(
			"double submit cookie",
			&condition{
				manifest: &v1.CSRFMiddleware{
					APIVersion: apiVersion,
					Kind:       kind,
					Metadata: &kernel.Metadata{
						Namespace: "default",
						Name:      "default",
					},
					Spec: &v1.CSRFMiddlewareSpec{
						SeedSize: 20,
						Secret:   base64.StdEncoding.EncodeToString(hmacSecret[:]),
						HashAlg:  kernel.HashAlg_SHA256,
						CSRFPatterns: &v1.CSRFMiddlewareSpec_DoubleSubmitCookies{
							DoubleSubmitCookies: &v1.DoubleSubmitCookiesSpec{},
						},
					},
				},
			},
			&action{
				manifest: &v1.CSRFMiddleware{
					APIVersion: apiVersion,
					Kind:       kind,
					Metadata: &kernel.Metadata{
						Namespace: "default",
						Name:      "default",
					},
					Spec: &v1.CSRFMiddlewareSpec{
						Patterns: []string{"/token"},
						SeedSize: 20,
						Secret:   base64.StdEncoding.EncodeToString(hmacSecret[:]),
						HashAlg:  kernel.HashAlg_SHA256,
						CSRFPatterns: &v1.CSRFMiddlewareSpec_DoubleSubmitCookies{
							DoubleSubmitCookies: &v1.DoubleSubmitCookiesSpec{
								CookieName:  "__csrfToken",
								TokenSource: v1.TokenSource_Header,
								SourceKey:   "__csrfToken",
							},
						},
					},
				},
			},
		),
		gen(
			"double submit cookie/overwrite",
			&condition{
				manifest: &v1.CSRFMiddleware{
					APIVersion: apiVersion,
					Kind:       kind,
					Metadata: &kernel.Metadata{
						Namespace: "default",
						Name:      "default",
					},
					Spec: &v1.CSRFMiddlewareSpec{
						SeedSize: 20,
						Secret:   base64.StdEncoding.EncodeToString(hmacSecret[:]),
						HashAlg:  kernel.HashAlg_SHA256,
						CSRFPatterns: &v1.CSRFMiddlewareSpec_DoubleSubmitCookies{
							DoubleSubmitCookies: &v1.DoubleSubmitCookiesSpec{
								CookieName:  "__foo",
								TokenSource: v1.TokenSource_Form,
								SourceKey:   "__bar",
							},
						},
					},
				},
			},
			&action{
				manifest: &v1.CSRFMiddleware{
					APIVersion: apiVersion,
					Kind:       kind,
					Metadata: &kernel.Metadata{
						Namespace: "default",
						Name:      "default",
					},
					Spec: &v1.CSRFMiddlewareSpec{
						Patterns: []string{"/token"},
						SeedSize: 20,
						Secret:   base64.StdEncoding.EncodeToString(hmacSecret[:]),
						HashAlg:  kernel.HashAlg_SHA256,
						CSRFPatterns: &v1.CSRFMiddlewareSpec_DoubleSubmitCookies{
							DoubleSubmitCookies: &v1.DoubleSubmitCookiesSpec{
								CookieName:  "__foo",
								TokenSource: v1.TokenSource_Form,
								SourceKey:   "__bar",
							},
						},
					},
				},
			},
		),
		gen(
			"synchronizer token",
			&condition{
				manifest: &v1.CSRFMiddleware{
					APIVersion: apiVersion,
					Kind:       kind,
					Metadata: &kernel.Metadata{
						Namespace: "default",
						Name:      "default",
					},
					Spec: &v1.CSRFMiddlewareSpec{
						SeedSize: 20,
						Secret:   base64.StdEncoding.EncodeToString(hmacSecret[:]),
						HashAlg:  kernel.HashAlg_SHA256,
						CSRFPatterns: &v1.CSRFMiddlewareSpec_SynchronizerToken{
							SynchronizerToken: &v1.SynchronizerTokenSpec{},
						},
					},
				},
			},
			&action{
				manifest: &v1.CSRFMiddleware{
					APIVersion: apiVersion,
					Kind:       kind,
					Metadata: &kernel.Metadata{
						Namespace: "default",
						Name:      "default",
					},
					Spec: &v1.CSRFMiddlewareSpec{
						Patterns: []string{"/token"},
						SeedSize: 20,
						Secret:   base64.StdEncoding.EncodeToString(hmacSecret[:]),
						HashAlg:  kernel.HashAlg_SHA256,
						CSRFPatterns: &v1.CSRFMiddlewareSpec_SynchronizerToken{
							SynchronizerToken: &v1.SynchronizerTokenSpec{
								TokenSource: v1.TokenSource_Header,
								SourceKey:   "__csrfToken",
							},
						},
					},
				},
			},
		),
		gen(
			"synchronizer token/overwrite",
			&condition{
				manifest: &v1.CSRFMiddleware{
					APIVersion: apiVersion,
					Kind:       kind,
					Metadata: &kernel.Metadata{
						Namespace: "default",
						Name:      "default",
					},
					Spec: &v1.CSRFMiddlewareSpec{
						SeedSize: 20,
						Secret:   base64.StdEncoding.EncodeToString(hmacSecret[:]),
						HashAlg:  kernel.HashAlg_SHA256,
						CSRFPatterns: &v1.CSRFMiddlewareSpec_SynchronizerToken{
							SynchronizerToken: &v1.SynchronizerTokenSpec{
								TokenSource: v1.TokenSource_Form,
								SourceKey:   "__bar",
							},
						},
					},
				},
			},
			&action{
				manifest: &v1.CSRFMiddleware{
					APIVersion: apiVersion,
					Kind:       kind,
					Metadata: &kernel.Metadata{
						Namespace: "default",
						Name:      "default",
					},
					Spec: &v1.CSRFMiddlewareSpec{
						Patterns: []string{"/token"},
						SeedSize: 20,
						Secret:   base64.StdEncoding.EncodeToString(hmacSecret[:]),
						HashAlg:  kernel.HashAlg_SHA256,
						CSRFPatterns: &v1.CSRFMiddlewareSpec_SynchronizerToken{
							SynchronizerToken: &v1.SynchronizerTokenSpec{
								TokenSource: v1.TokenSource_Form,
								SourceKey:   "__bar",
							},
						},
					},
				},
			},
		),
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			msg := Resource.Mutate(tt.C.manifest)

			opts := []cmp.Option{
				cmpopts.IgnoreUnexported(v1.CSRFMiddleware{}, v1.CSRFMiddlewareSpec{}),
				cmpopts.IgnoreUnexported(v1.CSRFMiddlewareSpec_CustomRequestHeader{}, v1.CustomRequestHeaderSpec{}),
				cmpopts.IgnoreUnexported(v1.CSRFMiddlewareSpec_DoubleSubmitCookies{}, v1.DoubleSubmitCookiesSpec{}),
				cmpopts.IgnoreUnexported(v1.CSRFMiddlewareSpec_SynchronizerToken{}, v1.SynchronizerTokenSpec{}),
				cmpopts.IgnoreUnexported(kernel.Metadata{}, kernel.Status{}),
			}
			testutil.Diff(t, tt.A.manifest, msg, opts...)
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
		expect     any
	}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"create with default manifest",
			&condition{
				manifest: Resource.Default(),
			},
			&action{
				err: nil,
				expect: &csrf{
					HandlerBase:     &utilhttp.HandlerBase{},
					eh:              utilhttp.GlobalErrorHandler(utilhttp.DefaultErrorHandlerName),
					proxyHeaderName: "",
					issueNew:        false,
					token: &csrfToken{
						secret:   defaultSecret[:],
						seedSize: 32,
						hashSize: 32,
						hmac:     hash.HMACFromHashAlg(kernel.HashAlg_SHA256),
					},
					st: &customRequestHeaders{
						headerName: "X-Requested-With",
						token: &csrfToken{
							secret:   defaultSecret[:],
							seedSize: 32,
							hashSize: 32,
							hmac:     hash.HMACFromHashAlg(kernel.HashAlg_SHA256),
						},
					},
				},
			},
		),
		gen(
			"invalid secret",
			&condition{
				manifest: &v1.CSRFMiddleware{
					APIVersion: apiVersion,
					Kind:       kind,
					Metadata: &kernel.Metadata{
						Namespace: "default",
						Name:      "default",
					},
					Spec: &v1.CSRFMiddlewareSpec{
						Secret: "invalid base64",
					},
				},
			},
			&action{
				err:        core.ErrCoreGenCreateObject,
				errPattern: regexp.MustCompile(core.ErrPrefix + `failed to create CSRFMiddleware`),
			},
		),
		gen(
			"invalid hash algorithm",
			&condition{
				manifest: &v1.CSRFMiddleware{
					APIVersion: apiVersion,
					Kind:       kind,
					Metadata: &kernel.Metadata{
						Namespace: "default",
						Name:      "default",
					},
					Spec: &v1.CSRFMiddlewareSpec{
						HashAlg: 999,
					},
				},
			},
			&action{
				err:        core.ErrCoreGenCreateObject,
				errPattern: regexp.MustCompile(core.ErrPrefix + `failed to create CSRFMiddleware`),
			},
		),
		gen(
			"custom request header",
			&condition{
				manifest: &v1.CSRFMiddleware{
					APIVersion: apiVersion,
					Kind:       kind,
					Metadata: &kernel.Metadata{
						Namespace: "default",
						Name:      "default",
					},
					Spec: &v1.CSRFMiddlewareSpec{
						CSRFPatterns: &v1.CSRFMiddlewareSpec_CustomRequestHeader{
							CustomRequestHeader: &v1.CustomRequestHeaderSpec{
								HeaderName: "X-Requested-With",
							},
						},
						HashAlg: kernel.HashAlg_SHA256,
					},
				},
			},
			&action{
				err: nil,
				expect: &csrf{
					HandlerBase: &utilhttp.HandlerBase{},
					eh:          utilhttp.GlobalErrorHandler(utilhttp.DefaultErrorHandlerName),
					token: &csrfToken{
						secret:   []byte{},
						seedSize: 32,
						hashSize: 32,
						hmac:     hash.HMACFromHashAlg(kernel.HashAlg_SHA256),
					},
					st: &customRequestHeaders{
						headerName: "X-Requested-With",
						token: &csrfToken{
							secret:   []byte{},
							seedSize: 32,
							hashSize: 32,
							hmac:     hash.HMACFromHashAlg(kernel.HashAlg_SHA256),
						},
					},
				},
			},
		),
		gen(
			"custom request header failure",
			&condition{
				manifest: &v1.CSRFMiddleware{
					APIVersion: apiVersion,
					Kind:       kind,
					Metadata: &kernel.Metadata{
						Namespace: "default",
						Name:      "default",
					},
					Spec: &v1.CSRFMiddlewareSpec{
						CSRFPatterns: &v1.CSRFMiddlewareSpec_CustomRequestHeader{
							CustomRequestHeader: &v1.CustomRequestHeaderSpec{
								HeaderName:     "X-Requested-With",
								AllowedPattern: "[invalid-regex",
							},
						},
						HashAlg: kernel.HashAlg_SHA256,
					},
				},
			},
			&action{
				err:        core.ErrCoreGenCreateObject,
				errPattern: regexp.MustCompile(core.ErrPrefix + `failed to create CSRFMiddleware`),
			},
		),
		gen(
			"double submit cookies",
			&condition{
				manifest: &v1.CSRFMiddleware{
					APIVersion: apiVersion,
					Kind:       kind,
					Metadata: &kernel.Metadata{
						Namespace: "default",
						Name:      "default",
					},
					Spec: &v1.CSRFMiddlewareSpec{
						CSRFPatterns: &v1.CSRFMiddlewareSpec_DoubleSubmitCookies{
							DoubleSubmitCookies: &v1.DoubleSubmitCookiesSpec{
								CookieName: "__csrfToken",
							},
						},
						HashAlg: kernel.HashAlg_SHA256,
					},
				},
			},
			&action{
				err: nil,
				expect: &csrf{
					HandlerBase: &utilhttp.HandlerBase{},
					eh:          utilhttp.GlobalErrorHandler(utilhttp.DefaultErrorHandlerName),
					token: &csrfToken{
						secret:   []byte{},
						seedSize: 32,
						hashSize: 32,
						hmac:     hash.HMACFromHashAlg(kernel.HashAlg_SHA256),
					},
					st: &doubleSubmitCookies{
						cookieName: "__csrfToken",
						ext:        &headerExtractor{},
						cookie:     &utilhttp.CookieCreator{Path: "/", Secure: true, HTTPOnly: true, SameSite: 1},
						token: &csrfToken{
							secret:   []byte{},
							seedSize: 32,
							hashSize: 32,
							hmac:     hash.HMACFromHashAlg(kernel.HashAlg_SHA256),
						},
					},
				},
			},
		),
		gen(
			"double submit cookies with form",
			&condition{
				manifest: &v1.CSRFMiddleware{
					APIVersion: apiVersion,
					Kind:       kind,
					Metadata: &kernel.Metadata{
						Namespace: "default",
						Name:      "default",
					},
					Spec: &v1.CSRFMiddlewareSpec{
						CSRFPatterns: &v1.CSRFMiddlewareSpec_DoubleSubmitCookies{
							DoubleSubmitCookies: &v1.DoubleSubmitCookiesSpec{
								CookieName:  "__csrfToken",
								TokenSource: v1.TokenSource_Form,
								SourceKey:   "__csrfToken",
							},
						},
						HashAlg: kernel.HashAlg_SHA256,
					},
				},
			},
			&action{
				err: nil,
				expect: &csrf{
					HandlerBase: &utilhttp.HandlerBase{},
					eh:          utilhttp.GlobalErrorHandler(utilhttp.DefaultErrorHandlerName),
					token: &csrfToken{
						secret:   []byte{},
						seedSize: 32,
						hashSize: 32,
						hmac:     hash.HMACFromHashAlg(kernel.HashAlg_SHA256),
					},
					st: &doubleSubmitCookies{
						cookieName: "__csrfToken",
						ext: &formExtractor{
							paramName: "__csrfToken",
						},
						cookie: &utilhttp.CookieCreator{Path: "/", Secure: true, HTTPOnly: true, SameSite: 1},
						token: &csrfToken{
							secret:   []byte{},
							seedSize: 32,
							hashSize: 32,
							hmac:     hash.HMACFromHashAlg(kernel.HashAlg_SHA256),
						},
					},
				},
			},
		),
		gen(
			"double submit cookies with JSON",
			&condition{
				manifest: &v1.CSRFMiddleware{
					APIVersion: apiVersion,
					Kind:       kind,
					Metadata: &kernel.Metadata{
						Namespace: "default",
						Name:      "default",
					},
					Spec: &v1.CSRFMiddlewareSpec{
						CSRFPatterns: &v1.CSRFMiddlewareSpec_DoubleSubmitCookies{
							DoubleSubmitCookies: &v1.DoubleSubmitCookiesSpec{
								CookieName:  "__csrfToken",
								TokenSource: v1.TokenSource_JSON,
								SourceKey:   "csrf.token",
							},
						},
						HashAlg: kernel.HashAlg_SHA256,
					},
				},
			},
			&action{
				err: nil,
				expect: &csrf{
					HandlerBase: &utilhttp.HandlerBase{},
					eh:          utilhttp.GlobalErrorHandler(utilhttp.DefaultErrorHandlerName),
					token: &csrfToken{
						secret:   []byte{},
						seedSize: 32,
						hashSize: 32,
						hmac:     hash.HMACFromHashAlg(kernel.HashAlg_SHA256),
					},
					st: &doubleSubmitCookies{
						cookieName: "__csrfToken",
						ext: &jsonExtractor{
							jsonPath: "csrf.token",
						},
						cookie: &utilhttp.CookieCreator{Path: "/", Secure: true, HTTPOnly: true, SameSite: 1},
						token: &csrfToken{
							secret:   []byte{},
							seedSize: 32,
							hashSize: 32,
							hmac:     hash.HMACFromHashAlg(kernel.HashAlg_SHA256),
						},
					},
				},
			},
		),
		gen(
			"synchronizer token",
			&condition{
				manifest: &v1.CSRFMiddleware{
					APIVersion: apiVersion,
					Kind:       kind,
					Metadata: &kernel.Metadata{
						Namespace: "default",
						Name:      "default",
					},
					Spec: &v1.CSRFMiddlewareSpec{
						CSRFPatterns: &v1.CSRFMiddlewareSpec_SynchronizerToken{
							SynchronizerToken: &v1.SynchronizerTokenSpec{
								SourceKey: "__csrfToken",
							},
						},
						HashAlg: kernel.HashAlg_SHA256,
					},
				},
			},
			&action{
				err: nil,
				expect: &csrf{
					HandlerBase: &utilhttp.HandlerBase{},
					eh:          utilhttp.GlobalErrorHandler(utilhttp.DefaultErrorHandlerName),
					token: &csrfToken{
						secret:   []byte{},
						seedSize: 32,
						hashSize: 32,
						hmac:     hash.HMACFromHashAlg(kernel.HashAlg_SHA256),
					},
					st: &synchronizerToken{
						ext: &headerExtractor{headerName: "__csrftoken"},
						token: &csrfToken{
							secret:   []byte{},
							seedSize: 32,
							hashSize: 32,
							hmac:     hash.HMACFromHashAlg(kernel.HashAlg_SHA256),
						},
					},
				},
			},
		),
		gen(
			"synchronizer token with form",
			&condition{
				manifest: &v1.CSRFMiddleware{
					APIVersion: apiVersion,
					Kind:       kind,
					Metadata: &kernel.Metadata{
						Namespace: "default",
						Name:      "default",
					},
					Spec: &v1.CSRFMiddlewareSpec{
						CSRFPatterns: &v1.CSRFMiddlewareSpec_SynchronizerToken{
							SynchronizerToken: &v1.SynchronizerTokenSpec{
								SourceKey:   "__csrfToken",
								TokenSource: v1.TokenSource_Form,
							},
						},
						HashAlg: kernel.HashAlg_SHA256,
					},
				},
			},
			&action{
				err: nil,
				expect: &csrf{
					HandlerBase: &utilhttp.HandlerBase{},
					eh:          utilhttp.GlobalErrorHandler(utilhttp.DefaultErrorHandlerName),
					token: &csrfToken{
						secret:   []byte{},
						seedSize: 32,
						hashSize: 32,
						hmac:     hash.HMACFromHashAlg(kernel.HashAlg_SHA256),
					},
					st: &synchronizerToken{
						ext: &formExtractor{
							paramName: "__csrfToken",
						},
						token: &csrfToken{
							secret:   []byte{},
							seedSize: 32,
							hashSize: 32,
							hmac:     hash.HMACFromHashAlg(kernel.HashAlg_SHA256),
						},
					},
				},
			},
		),
		gen(
			"synchronizer token with JSON",
			&condition{
				manifest: &v1.CSRFMiddleware{
					APIVersion: apiVersion,
					Kind:       kind,
					Metadata: &kernel.Metadata{
						Namespace: "default",
						Name:      "default",
					},
					Spec: &v1.CSRFMiddlewareSpec{
						CSRFPatterns: &v1.CSRFMiddlewareSpec_SynchronizerToken{
							SynchronizerToken: &v1.SynchronizerTokenSpec{
								SourceKey:   "csrf.token",
								TokenSource: v1.TokenSource_JSON,
							},
						},
						HashAlg: kernel.HashAlg_SHA256,
					},
				},
			},
			&action{
				err: nil,
				expect: &csrf{
					HandlerBase: &utilhttp.HandlerBase{},
					eh:          utilhttp.GlobalErrorHandler(utilhttp.DefaultErrorHandlerName),
					token: &csrfToken{
						secret:   []byte{},
						seedSize: 32,
						hashSize: 32,
						hmac:     hash.HMACFromHashAlg(kernel.HashAlg_SHA256),
					},
					st: &synchronizerToken{
						ext: &jsonExtractor{
							jsonPath: "csrf.token",
						},
						token: &csrfToken{
							secret:   []byte{},
							seedSize: 32,
							hashSize: 32,
							hmac:     hash.HMACFromHashAlg(kernel.HashAlg_SHA256),
						},
					},
				},
			},
		),
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {

			server := api.NewContainerAPI()
			a := &API{}
			got, err := a.Create(server, tt.C.manifest)

			opts := []cmp.Option{
				protocmp.Transform(),
				cmp.AllowUnexported(csrf{}, csrfToken{}, customRequestHeaders{}, doubleSubmitCookies{}, synchronizerToken{}, headerExtractor{}, formExtractor{}, jsonExtractor{}),
				cmp.Comparer(testutil.ComparePointer[core.ErrorHandler]),
				cmp.Comparer(func(x, y *csrfToken) bool {
					return cmp.Equal(x.secret, y.secret) &&
						x.seedSize == y.seedSize &&
						x.hashSize == y.hashSize &&
						x.hmac != nil && y.hmac != nil // Avoid function comparison failures
				}),
			}
			testutil.DiffError(t, tt.A.err, tt.A.errPattern, err)
			testutil.Diff(t, tt.A.expect, got, opts...)

		})
	}
}
