// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package oauth

import (
	"crypto/rand"
	"crypto/x509"
	"maps"
	"net/http"
	"net/http/httptest"
	"regexp"
	"slices"
	"testing"

	v1 "github.com/aileron-gateway/aileron-gateway/apis/app/v1"
	"github.com/aileron-gateway/aileron-gateway/app"
	"github.com/aileron-gateway/aileron-gateway/core"
	"github.com/aileron-gateway/aileron-gateway/kernel/log"
	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
	utilhttp "github.com/aileron-gateway/aileron-gateway/util/http"
	"github.com/google/go-cmp/cmp"
)

type testSkipper struct {
	methods []string
	paths   []string
}

func (s *testSkipper) Match(r *http.Request) bool {
	if !slices.Contains(s.methods, r.Method) {
		return false
	}
	if !slices.Contains(s.paths, r.URL.Path) {
		return false
	}
	return true
}

type testErrorHandler struct {
	err    any // error or errorutil.Kind
	called bool
	hook   func(error)
}

func (h *testErrorHandler) ServeHTTPError(w http.ResponseWriter, r *http.Request, err error) {
	h.called = true
	h.err = err
	if h.hook != nil {
		h.hook(err)
	}
}

func TestNewResourceServerHandler(t *testing.T) {
	type condition struct {
		bh     *baseHandler
		spec   *v1.ResourceServerHandler
		header http.Header
	}

	type action struct {
		h *resourceServerHandler
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	condRedeemPath := tb.Condition("authentication via redeem token endpoint", "authentication via redeem token endpoint")
	condHeaderKeyIsEmpty := tb.Condition("input data HeaderKey is Empty", "input data HeaderKey is Empty")
	actNoError := tb.Action("no error", "check that the there is no error")

	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"no error",
			[]string{condRedeemPath},
			[]string{actNoError},
			&condition{
				spec: &v1.ResourceServerHandler{
					HeaderKey:   "test-key",
					EnabledFAPI: true,
				},
			},
			&action{
				h: &resourceServerHandler{
					headerKey:   "test-key",
					fapiEnabled: true,
				},
			},
		),
		gen(
			"empty header key",
			[]string{condHeaderKeyIsEmpty},
			[]string{actNoError},
			&condition{
				spec: &v1.ResourceServerHandler{
					HeaderKey:   "",
					EnabledFAPI: true,
				},
			},
			&action{
				h: &resourceServerHandler{
					headerKey:   "Authorization",
					fapiEnabled: true,
				},
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			handler := newResourceServerHandler(tt.C().bh, tt.C().spec)
			opts := []cmp.Option{
				cmp.AllowUnexported(resourceServerHandler{}, baseHandler{}),
				// cmpopts.IgnoreInterfaces(struct{ http.RoundTripper }{}),
			}
			testutil.Diff(t, tt.A().h, handler, opts...)
		})
	}
}

func TestResourceServer_ServeAuthn(t *testing.T) {
	type condition struct {
		h        *resourceServerHandler
		r        *http.Request
		header   http.Header
		uidError bool
	}

	type action struct {
		authenticated    app.AuthResult
		shouldReturn     bool
		fapiHeaderExists bool

		err        any // error or errorutil.Kind
		errPattern *regexp.Regexp
		errStatus  int
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	condEmptyAT := tb.Condition("access token has empty string", "access token has empty string")
	condTokenIntrospection := tb.Condition("token introspection enabled", "token introspection enabled")
	condValidateTokenFailed := tb.Condition("failed to validate token", "failed to validate token")
	condFAPIEnabled := tb.Condition("fapi enabled", "fapi enabled")
	condValidateCertFailed := tb.Condition("failed to validate client certificate", "failed to validate client certificate")
	condGenerateUUIDFailed := tb.Condition("fapi-interaction-id header doesn't exist in the request and failed to generate uuid", "fapi-interaction-id header doesn't exist in the request and failed to generate uuid")
	condOauthContextNotFound := tb.Condition("failed to validate OAuthContext", "failed to validate OAuthContext")
	condAuthenticationSucceeded := tb.Condition("Authentication Succeeded", "Authentication Succeeded")
	condOauthContextFromHeader := tb.Condition("OAuthContext validation succeeds from header", "OAuthContext validation succeeds from header")
	condOauthContextFromQuery := tb.Condition("OAuthContext validation succeeds from query", "OAuthContext validation succeeds from query")
	actError := tb.Action("error", "check that the expected error is returned")
	actNoError := tb.Action("no error", "check that the there is no error")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"context not found",
			[]string{condOauthContextNotFound},
			[]string{actError},
			&condition{
				h: &resourceServerHandler{
					baseHandler: &baseHandler{
						lg:        log.GlobalLogger(log.DefaultLoggerName),
						oauthCtxs: map[string]*oauthContext{},
					},
				},
				r:      httptest.NewRequest(http.MethodGet, "https://test.com/", nil),
				header: http.Header{},
			},
			&action{
				authenticated: app.AuthContinue,
				shouldReturn:  false,
				err:           nil,
			},
		),
		gen(
			"access token not found",
			[]string{condEmptyAT},
			[]string{actError},
			&condition{
				h: &resourceServerHandler{
					baseHandler: &baseHandler{
						lg:        log.GlobalLogger(log.DefaultLoggerName),
						oauthCtxs: map[string]*oauthContext{"default": {}},
					},
					headerKey: "Authorization",
				},
				r: httptest.NewRequest(http.MethodGet, "https://test.com/", nil),
				header: http.Header{
					"Authorization": []string{""},
				},
			},
			&action{
				authenticated: app.AuthContinue,
				shouldReturn:  false,
				err:           nil,
			},
		),
		gen(
			"AT validation error",
			[]string{condTokenIntrospection, condValidateTokenFailed},
			[]string{actError},
			&condition{
				h: &resourceServerHandler{
					baseHandler: &baseHandler{
						lg:        log.GlobalLogger(log.DefaultLoggerName),
						oauthCtxs: map[string]*oauthContext{"default": testDefaultOauthContext},
					},
					headerKey: "Authorization",
				},
				r: httptest.NewRequest(http.MethodGet, "https://test.com/", nil),
				header: http.Header{
					"Authorization": []string{"Bearer TestAccessToken"},
				},
			},
			&action{
				authenticated: false,
				shouldReturn:  true,
				errStatus:     http.StatusUnauthorized,
				err:           reAuthenticationRequired,
			},
		),
		gen(
			"authentication succeeded",
			[]string{condAuthenticationSucceeded},
			[]string{actNoError},
			&condition{
				h: &resourceServerHandler{
					baseHandler: &baseHandler{
						lg:        log.GlobalLogger(log.DefaultLoggerName),
						oauthCtxs: map[string]*oauthContext{"default": testDefaultOauthContext},
					},
					headerKey: "Authorization",
				},
				r: httptest.NewRequest(http.MethodGet, "https://test.com/", nil),
				header: http.Header{
					"Authorization": []string{"Bearer " + testSimpleJWT},
				},
			},
			&action{
				authenticated: true,
				shouldReturn:  false,
				err:           nil,
			},
		),
		gen(
			"OAuthContext validation succeeds from query",
			[]string{condOauthContextFromQuery},
			[]string{actNoError},
			&condition{
				h: &resourceServerHandler{
					baseHandler: &baseHandler{
						lg:              log.GlobalLogger(log.DefaultLoggerName),
						oauthCtxs:       map[string]*oauthContext{"default": testDefaultOauthContext},
						contextQueryKey: "contextQuery",
					},
					headerKey: "Authorization",
				},
				r: httptest.NewRequest(http.MethodGet, "https://test.com?contextQuery=default", nil),
				header: http.Header{
					"Authorization": []string{"Bearer " + testSimpleJWT},
				},
			},
			&action{
				authenticated: true,
				shouldReturn:  false,
				err:           nil,
			},
		),
		gen(
			"OAuthContext validation succeeds from header",
			[]string{condOauthContextFromHeader},
			[]string{actNoError},
			&condition{
				h: &resourceServerHandler{
					baseHandler: &baseHandler{
						lg:               log.GlobalLogger(log.DefaultLoggerName),
						oauthCtxs:        map[string]*oauthContext{"default": testDefaultOauthContext},
						contextHeaderKey: "contextHeader",
					},
					headerKey: "Authorization",
				},
				r: httptest.NewRequest(http.MethodGet, "https://test.com/", nil),
				header: http.Header{
					"Authorization": []string{"Bearer " + testSimpleJWT},
					"contextHeader": []string{"default"},
				},
			},
			&action{
				authenticated: true,
				shouldReturn:  false,
				err:           nil,
			},
		),
		gen(
			"mtls success",
			[]string{condFAPIEnabled},
			[]string{actNoError},
			&condition{
				h: &resourceServerHandler{
					baseHandler: &baseHandler{
						lg:        log.GlobalLogger(log.DefaultLoggerName),
						oauthCtxs: map[string]*oauthContext{"default": testDefaultOauthContext},
					},
					headerKey:   "Authorization",
					fapiEnabled: true,
				},
				r: httptest.NewRequest(http.MethodGet, "https://test.com/", nil),
				header: http.Header{
					"Authorization": []string{"Bearer " + testSimpleJWTCnf},
				},
			},
			&action{
				authenticated:    true,
				shouldReturn:     false,
				fapiHeaderExists: true,
				err:              nil,
			},
		),
		gen(
			"client cert validation error",
			[]string{condFAPIEnabled, condValidateCertFailed},
			[]string{actError},
			&condition{
				h: &resourceServerHandler{
					baseHandler: &baseHandler{
						lg:        log.GlobalLogger(log.DefaultLoggerName),
						oauthCtxs: map[string]*oauthContext{"default": testDefaultOauthContext},
					},
					headerKey:   "Authorization",
					fapiEnabled: true,
				},
				r: httptest.NewRequest(http.MethodGet, "https://test.com/", nil),
				header: http.Header{
					"Authorization": []string{"Bearer " + testSimpleJWT},
				},
			},
			&action{
				authenticated:    false,
				shouldReturn:     true,
				fapiHeaderExists: false,
				errStatus:        http.StatusUnauthorized,
				err:              app.ErrAppAuthnAuthentication,
				errPattern:       regexp.MustCompile(core.ErrPrefix + `authentication failed`),
			},
		),
		gen(
			"uid generate error",
			[]string{condFAPIEnabled, condValidateCertFailed, condGenerateUUIDFailed},
			[]string{actError},
			&condition{
				h: &resourceServerHandler{
					baseHandler: &baseHandler{
						lg:        log.GlobalLogger(log.DefaultLoggerName),
						oauthCtxs: map[string]*oauthContext{"default": testDefaultOauthContext},
					},
					headerKey:   "Authorization",
					fapiEnabled: true,
				},
				r: httptest.NewRequest(http.MethodGet, "https://test.com/", nil),
				header: http.Header{
					"Authorization": []string{"Bearer " + testSimpleJWTCnf},
				},
				uidError: true,
			},
			&action{
				authenticated:    false,
				shouldReturn:     true,
				fapiHeaderExists: false,
				errStatus:        http.StatusInternalServerError,
				err:              app.ErrAppAuthnAuthentication,
				errPattern:       regexp.MustCompile(core.ErrPrefix + `authentication failed`),
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			if tt.C().uidError {
				tmp := rand.Reader
				rand.Reader = &testutil.ErrorReader{}
				defer func() {
					rand.Reader = tmp
				}()
			}

			w := httptest.NewRecorder()
			r := tt.C().r
			maps.Copy(r.Header, tt.C().header)
			r.TLS.PeerCertificates = []*x509.Certificate{
				{
					Raw: []byte("client certificate"),
				},
			}

			r, authenticated, shouldReturn, err := tt.C().h.ServeAuthn(w, r)

			testutil.Diff(t, tt.A().authenticated, authenticated)
			testutil.Diff(t, tt.A().shouldReturn, shouldReturn)

			if tt.A().fapiHeaderExists {
				// Only check the value exists.
				testutil.Diff(t, true, w.Header().Get("Date") != "")
				testutil.Diff(t, true, w.Header().Get("x-fapi-interaction-id") != "")
			} else {
				testutil.Diff(t, "", w.Header().Get("Date"))
				testutil.Diff(t, "", w.Header().Get("x-fapi-interaction-id"))
			}

			if tt.A().err == nil {
				testutil.Diff(t, nil, err)
				return
			}

			e := err.(core.HTTPError)
			testutil.Diff(t, tt.A().errStatus, e.StatusCode())
			if tt.A().err == reAuthenticationRequired {
				testutil.Diff(t, reAuthenticationRequired.Error(), err.Error())
			} else {
				e := err.(*utilhttp.HTTPError)
				testutil.DiffError(t, tt.A().err, tt.A().errPattern, e.Unwrap())
			}
		})
	}
}
