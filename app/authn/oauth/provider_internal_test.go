// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package oauth

import (
	"bytes"
	"errors"
	"net/http"
	"regexp"
	"testing"
	"time"

	v1 "github.com/aileron-gateway/aileron-gateway/apis/app/v1"
	"github.com/aileron-gateway/aileron-gateway/core"
	"github.com/aileron-gateway/aileron-gateway/internal/testutil"
	"github.com/aileron-gateway/aileron-gateway/kernel/log"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestNewClient(t *testing.T) {
	type condition struct {
		spec *v1.OAuthClient
	}

	type action struct {
		client *client
	}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"base and paths",
			&condition{
				spec: &v1.OAuthClient{
					ID:       "test-id",
					Secret:   "test-secret",
					Audience: "test-audience",
					Scopes:   []string{"s1", "s2"},
				},
			},
			&action{
				client: &client{
					id:       "test-id",
					secret:   "test-secret",
					audience: "test-audience",
					scope:    "s1 s2",
				},
			},
		),
		gen(
			"fill audience",
			&condition{
				spec: &v1.OAuthClient{
					ID:       "test-id",
					Audience: "",
				},
			},
			&action{
				client: &client{
					id:       "test-id",
					audience: "test-id",
				},
			},
		),
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			c, _ := newClient(tt.C.spec)
			testutil.Diff(t, tt.A.client, c, cmp.AllowUnexported(client{}))
		})
	}
}

func TestNewProvider(t *testing.T) {
	type condition struct {
		spec *v1.OAuthProvider
		rt   http.RoundTripper
	}

	type action struct {
		provider   *provider
		err        any // error or errorutil.Kind
		errPattern *regexp.Regexp
	}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"base and paths",
			&condition{
				spec: &v1.OAuthProvider{
					Issuer:  "http://test.com/issuer",
					BaseURL: "http://test.com",
					Endpoints: &v1.ProviderEndpoints{
						Authorization: "/authorization_endpoint",
						Token:         "/token_endpoint",
						Userinfo:      "/userinfo_endpoint",
						Introspection: "/introspection_endpoint",
						Revocation:    "/revocation_endpoint",
						JWKs:          "/jwks_uri",
						Discovery:     "/discovery",
					},
				},
			},
			&action{
				provider: &provider{
					issuer:          "http://test.com/issuer",
					authorizationEP: "http://test.com/authorization_endpoint",
					tokenEP:         "http://test.com/token_endpoint",
					userinfoEP:      "http://test.com/userinfo_endpoint",
					jwksEP:          "http://test.com/jwks_uri",
					introspectEP:    "http://test.com/introspection_endpoint",
					revocationEP:    "http://test.com/revocation_endpoint",
					discoveryEP:     "http://test.com/discovery",
				},
			},
		),
		gen(
			"empty base url",
			&condition{
				spec: &v1.OAuthProvider{
					Issuer:  "http://test.com/issuer",
					BaseURL: "",
					Endpoints: &v1.ProviderEndpoints{
						Authorization: "http://test.com/authorization_endpoint",
						Token:         "http://test.com/token_endpoint",
						Userinfo:      "http://test.com/userinfo_endpoint",
						Introspection: "http://test.com/introspection_endpoint",
						Revocation:    "http://test.com/revocation_endpoint",
						JWKs:          "http://test.com/jwks_uri",
						Discovery:     "http://test.com/discovery",
					},
				},
			},
			&action{
				provider: &provider{
					issuer:          "http://test.com/issuer",
					authorizationEP: "http://test.com/authorization_endpoint",
					tokenEP:         "http://test.com/token_endpoint",
					userinfoEP:      "http://test.com/userinfo_endpoint",
					jwksEP:          "http://test.com/jwks_uri",
					introspectEP:    "http://test.com/introspection_endpoint",
					revocationEP:    "http://test.com/revocation_endpoint",
					discoveryEP:     "http://test.com/discovery",
				},
			},
		),
		gen(
			"empty path",
			&condition{
				spec: &v1.OAuthProvider{
					BaseURL: "http://test.com",
					Endpoints: &v1.ProviderEndpoints{
						Authorization: "",
					},
				},
			},
			&action{
				provider: &provider{
					authorizationEP: "",
				},
			},
		),
		gen(
			"path join error",
			&condition{
				spec: &v1.OAuthProvider{
					BaseURL: "http://test.com\n",
					Endpoints: &v1.ProviderEndpoints{
						Authorization: "/authorization_endpoint",
					},
				},
			},
			&action{
				provider:   nil,
				err:        core.ErrCoreGenCreateComponent,
				errPattern: regexp.MustCompile(core.ErrPrefix + `failed to create component. failed to create provider`),
			},
		),
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			p, err := newProvider(tt.C.spec, tt.C.rt)
			testutil.DiffError(t, tt.A.err, tt.A.errPattern, err)

			if p != nil && p.discoveryEP != "" {
				p.close = make(chan struct{})
				close(p.close)
			}

			opts := []cmp.Option{
				cmp.AllowUnexported(provider{}),
				cmpopts.IgnoreFields(provider{}, "lg", "rt", "ticker", "close"),
			}
			testutil.Diff(t, tt.A.provider, p, opts...)
		})
	}
}

func TestProvider_discover(t *testing.T) {
	type condition struct {
		client   *testClient
		provider *provider
	}

	type action struct {
		provider *provider
	}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"update with new value",
			&condition{
				provider: &provider{
					lg: log.GlobalLogger(log.DefaultLoggerName),
					rt: &testClient{
						status: http.StatusOK,
						body: bytes.NewReader([]byte(`{
							"issuer":"http://test.com/new/issuer",
							"authorization_endpoint":"http://test.com/new/authorization_endpoint",
							"token_endpoint":"http://test.com/new/token_endpoint",
							"userinfo_endpoint":"http://test.com/new/userinfo_endpoint",
							"jwks_uri":"http://test.com/new/jwks_uri",
							"introspection_endpoint":"http://test.com/new/introspection_endpoint",
							"revocation_endpoint":"http://test.com/new/revocation_endpoint"
							}`)),
					},
					ticker:      time.NewTicker(time.Microsecond),
					close:       make(chan struct{}),
					discoveryEP: "http://test.com/discovery",
				},
			},
			&action{
				provider: &provider{
					issuer:          "http://test.com/new/issuer",
					authorizationEP: "http://test.com/new/authorization_endpoint",
					tokenEP:         "http://test.com/new/token_endpoint",
					userinfoEP:      "http://test.com/new/userinfo_endpoint",
					jwksEP:          "http://test.com/new/jwks_uri",
					introspectEP:    "http://test.com/new/introspection_endpoint",
					revocationEP:    "http://test.com/new/revocation_endpoint",
					discoveryEP:     "http://test.com/discovery",
				},
			},
		),
		gen(
			"respect old value",
			&condition{
				provider: &provider{
					lg: log.GlobalLogger(log.DefaultLoggerName),
					rt: &testClient{
						status: http.StatusOK,
						body: bytes.NewReader([]byte(`{
							"issuer":"http://test.com/new/issuer",
							"authorization_endpoint":"http://test.com/new/authorization_endpoint",
							"token_endpoint":"http://test.com/new/token_endpoint",
							"userinfo_endpoint":"http://test.com/new/userinfo_endpoint",
							"jwks_uri":"http://test.com/new/jwks_uri",
							"introspection_endpoint":"http://test.com/new/introspection_endpoint",
							"revocation_endpoint":"http://test.com/new/revocation_endpoint"
							}`)),
					},
					ticker:          time.NewTicker(time.Microsecond),
					close:           make(chan struct{}),
					issuer:          "http://test.com/old/issuer",
					authorizationEP: "http://test.com/old/authorization_endpoint",
					tokenEP:         "http://test.com/old/token_endpoint",
					userinfoEP:      "http://test.com/old/userinfo_endpoint",
					jwksEP:          "http://test.com/old/jwks_uri",
					introspectEP:    "http://test.com/old/introspection_endpoint",
					revocationEP:    "http://test.com/old/revocation_endpoint",
					discoveryEP:     "http://test.com/discovery",
				},
			},
			&action{
				provider: &provider{
					issuer:          "http://test.com/old/issuer",
					authorizationEP: "http://test.com/old/authorization_endpoint",
					tokenEP:         "http://test.com/old/token_endpoint",
					userinfoEP:      "http://test.com/old/userinfo_endpoint",
					jwksEP:          "http://test.com/old/jwks_uri",
					introspectEP:    "http://test.com/old/introspection_endpoint",
					revocationEP:    "http://test.com/old/revocation_endpoint",
					discoveryEP:     "http://test.com/discovery",
				},
			},
		),
		gen(
			"invalid url returned",
			&condition{
				provider: &provider{
					lg: log.GlobalLogger(log.DefaultLoggerName),
					rt: &testClient{
						status: http.StatusOK,
						body:   bytes.NewReader([]byte(`{"issuer":"http://test.com/bad\n"}`)),
					},
					ticker:      time.NewTicker(time.Microsecond),
					close:       make(chan struct{}),
					issuer:      "http://test.com/old",
					discoveryEP: "http://test.com/discovery",
				},
			},
			&action{
				provider: &provider{
					issuer:      "http://test.com/old",
					discoveryEP: "http://test.com/discovery",
				},
			},
		),
		gen(
			"request generate error",
			&condition{
				provider: &provider{
					lg: log.GlobalLogger(log.DefaultLoggerName),
					rt: &testClient{
						status: http.StatusOK,
						body:   bytes.NewReader(nil),
					},
					ticker:      time.NewTicker(time.Microsecond),
					close:       make(chan struct{}),
					issuer:      "http://test.com/old",
					discoveryEP: "http://test.com/discovery\n", // Bad character ad the end.
				},
			},
			&action{
				provider: &provider{
					issuer:      "http://test.com/old",
					discoveryEP: "http://test.com/discovery\n",
				},
			},
		),
		gen(
			"parse invalid url",
			&condition{
				provider: &provider{
					lg: log.GlobalLogger(log.DefaultLoggerName),
					rt: &testClient{
						status: http.StatusOK,
						body:   bytes.NewReader([]byte(`{"issuer":"http://test.com/bad\n"}`)),
					},
					ticker:      time.NewTicker(time.Microsecond),
					close:       make(chan struct{}),
					issuer:      "",
					discoveryEP: "http://test.com/discovery",
				},
			},
			&action{
				provider: &provider{
					issuer:      "",
					discoveryEP: "http://test.com/discovery",
				},
			},
		),
		gen(
			"round trip error",
			&condition{
				provider: &provider{
					lg: log.GlobalLogger(log.DefaultLoggerName),
					rt: &testClient{
						err: errors.New("round trip error"),
					},
					ticker:      time.NewTicker(time.Microsecond),
					close:       make(chan struct{}),
					issuer:      "http://test.com/old",
					discoveryEP: "http://test.com/discovery",
				},
			},
			&action{
				provider: &provider{
					issuer:      "http://test.com/old",
					discoveryEP: "http://test.com/discovery",
				},
			},
		),
		gen(
			"response body read error",
			&condition{
				provider: &provider{
					lg: log.GlobalLogger(log.DefaultLoggerName),
					rt: &testClient{
						status: http.StatusOK,
						body:   &testReader{err: errors.New("body read error")},
					},
					ticker:      time.NewTicker(time.Microsecond),
					close:       make(chan struct{}),
					issuer:      "http://test.com/old",
					discoveryEP: "http://test.com/discovery",
				},
			},
			&action{
				provider: &provider{
					issuer:      "http://test.com/old",
					discoveryEP: "http://test.com/discovery",
				},
			},
		),
		gen(
			"500 server error",
			&condition{
				provider: &provider{
					lg: log.GlobalLogger(log.DefaultLoggerName),
					rt: &testClient{
						status: http.StatusInternalServerError,
						body:   bytes.NewReader(nil),
					},
					ticker:      time.NewTicker(time.Microsecond),
					close:       make(chan struct{}),
					discoveryEP: "http://test.com/discovery",
				},
			},
			&action{
				provider: &provider{
					discoveryEP: "http://test.com/discovery",
				},
			},
		),
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			go func() {
				time.Sleep(10 * time.Millisecond)
				close(tt.C.provider.close)
			}()
			tt.C.provider.discover()

			opts := []cmp.Option{
				cmp.AllowUnexported(provider{}),
				cmpopts.IgnoreFields(provider{}, "lg", "rt", "ticker", "close"),
			}
			testutil.Diff(t, tt.A.provider, tt.C.provider, opts...)
		})
	}
}
