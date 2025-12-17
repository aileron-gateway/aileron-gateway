// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package oauth

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"regexp"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	v1 "github.com/aileron-gateway/aileron-gateway/apis/app/v1"
	"github.com/aileron-gateway/aileron-gateway/app"
	"github.com/aileron-gateway/aileron-gateway/core"
	"github.com/aileron-gateway/aileron-gateway/kernel/api"
	"github.com/aileron-gateway/aileron-gateway/kernel/log"
	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
	utilhttp "github.com/aileron-gateway/aileron-gateway/util/http"
	"github.com/aileron-gateway/aileron-gateway/util/security"
	"github.com/aileron-gateway/aileron-gateway/util/session"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

type idtClaims struct {
	*jwt.RegisteredClaims
	Azp      string `json:"azp,omitempty"`
	Acr      string `json:"acr,omitempty"`
	AuthTime int    `json:"auth_time,omitempty"`
	Nonce    string `json:"nonce,omitempty"`
}

type invalidAuthTimeIDTClaims struct {
	*jwt.RegisteredClaims
	Azp      string `json:"azp,omitempty"`
	Acr      string `json:"acr,omitempty"`
	AuthTime string `json:"auth_time,omitempty"`
	Nonce    string `json:"nonce,omitempty"`
}

type atClaims struct {
	*jwt.RegisteredClaims
	ClientID  string         `json:"client_id,omitempty"`
	Username  string         `json:"username,omitempty"`
	TokenType string         `json:"token_type,omitempty"`
	Cnf       map[string]any `json:"cnf,omitempty"`
}

// func newJWTHandler() *security.JWTHandler {
// 	spec := &v1.JWTHandlerSpec{
// 		PublicKeys: []*v1.SigningKeySpec{
// 			{
// 				KeyID:       "J1EddGQKpO3dXzBHh8Y28LtiT7uDYuVvJuDLZyU9KSU",
// 				Algorithm:   v1.SigningKeyAlgorithm_RS256,
// 				KeyType:     v1.SigningKeyType_PUBLIC,
// 				KeyFilePath: "../../../test/ut/app/jwt/public.key",
// 			},
// 		},
// 		PrivateKeys: []*v1.SigningKeySpec{
// 			{
// 				KeyID:       "J1EddGQKpO3dXzBHh8Y28LtiT7uDYuVvJuDLZyU9KSU",
// 				Algorithm:   v1.SigningKeyAlgorithm_RS256,
// 				KeyType:     v1.SigningKeyType_PRIVATE,
// 				KeyFilePath: "../../../test/ut/app/jwt/private.key",
// 			},
// 		},
// 	}

// 	jh, _ := security.NewJWTHandler(spec, nil)
// 	return jh
// }

// type testTokenIntrospection struct {
// 	mc     jwt.MapClaims
// 	scope  string
// 	active bool
// 	cnf    map[string]any
// }

// func (t *testTokenIntrospection) tokenIntrospection(ctx context.Context, token string) (*TokenIntrospectionClaims, error) {
// 	if t.mc == nil {
// 		return nil, errors.New("token introspection error")
// 	}

// 	c := &TokenIntrospectionClaims{
// 		RegisteredClaims: &jwt.RegisteredClaims{
// 			ExpiresAt: mapValue[*jwt.NumericDate](t.mc, "exp"),
// 			IssuedAt:  mapValue[*jwt.NumericDate](t.mc, "iat"),
// 			NotBefore: mapValue[*jwt.NumericDate](t.mc, "nbf"),
// 			Subject:   mapValue[string](t.mc, "sub"),
// 			Audience:  mapValue[jwt.ClaimStrings](t.mc, "aud"),
// 			Issuer:    mapValue[string](t.mc, "iss"),
// 			ID:        mapValue[string](t.mc, "jti"),
// 		},
// 		Claims:    t.mc,
// 		Active:    t.active,
// 		Scope:     t.scope,
// 		Scopes:    strings.Split(t.scope, " "),
// 		ClientID:  mapValue[string](t.mc, "client_id"),
// 		Username:  mapValue[string](t.mc, "username"),
// 		TokenType: mapValue[string](t.mc, "token_type"),
// 		Cnf:       t.cnf,
// 	}
// 	return c, nil
// }

type testErrAPI[Q, S any] struct {
	api.API[Q, S]
}

func (a testErrAPI[Q, S]) Serve(ctx context.Context, q Q) (S, error) {
	var s S
	return s, errors.New("get error")
}

func TestNewROPCHandler(t *testing.T) {
	type condition struct {
		bh   *baseHandler
		spec *v1.ROPCHandler
	}

	type action struct {
		h *ropcHandler
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	actNoError := tb.Action("no error", "check that the there is no error")
	table := tb.Build()

	bh := &baseHandler{
		lg: log.GlobalLogger(log.DefaultLoggerName),
		oauthCtxs: map[string]*oauthContext{
			"test": {
				lg:   log.GlobalLogger(log.DefaultLoggerName),
				name: "test-context",
			},
		},
		contextHeaderKey: "test-context-key",
		contextQueryKey:  "test-query-key",
	}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"no error",
			[]string{},
			[]string{actNoError},
			&condition{
				bh: bh,
				spec: &v1.ROPCHandler{
					RedeemTokenPath: "/redeem",
					UsernameKey:     "username",
					PasswordKey:     "password",
				},
			},
			&action{
				h: &ropcHandler{
					baseHandler: bh,
					redeemPath:  "/redeem",
					usernameKey: "username",
					passwordKey: "password",
					queryParams: map[string]string{"grant_type": "password"},
				},
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			h := newROPCHandler(tt.C().bh, tt.C().spec)

			opts := []cmp.Option{
				cmp.Comparer(testutil.ComparePointer[log.Logger]),
				cmp.AllowUnexported(utilhttp.DefaultErrorHandler{}),
				cmp.AllowUnexported(baseHandler{}, oauthContext{}),
				cmp.AllowUnexported(ropcHandler{}, redeemTokenClient{}, provider{}, tokenIntrospectionClient{}, security.JWTHandler{}, security.SigningKey{}),
				cmp.AllowUnexported(testSkipper{}),
				cmpopts.IgnoreUnexported(sync.RWMutex{}, atomic.Bool{}),
			}
			testutil.Diff(t, tt.A().h, h, opts...)
		})
	}
}

type testTokenRedeemer struct {
	status   int
	response *TokenResponse
	err      core.HTTPError

	ctx    context.Context
	params map[string]string
}

func (t *testTokenRedeemer) redeemToken(ctx context.Context, params map[string]string) (int, *TokenResponse, core.HTTPError) {
	t.ctx = ctx
	t.params = params
	return t.status, t.response, t.err
}

func TestROPCHandler_ServeHTTP(t *testing.T) {
	type condition struct {
		r        *http.Request
		redeemer *testTokenRedeemer
	}

	type action struct {
		status     int
		body       string
		query      map[string]string
		err        any // error or errorutil.Kind
		errPattern *regexp.Regexp
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	// cndHeader := tb.Condition("header", "credentials in authorization header")
	// cndEmptyUsername := tb.Condition("empty usernameKey", "usernameKey is empty")
	// cndEmptyPassword := tb.Condition("empty passwordKey", "passwordKey is empty")
	// cndInsufficientCredential := tb.Condition("insufficient credentials", "provided username or password is invalid")
	// cndInvalidAuthHeader := tb.Condition("invalid auth header", "invalid authorization header value")
	// actCheckError := tb.Action("error", "check that the expected error is returned")
	actCheckNoError := tb.Action("no error", "check that the there is no error")
	table := tb.Build()

	defaultCtxReq := httptest.NewRequest(http.MethodPost, "https://test.com/token", nil)
	defaultCtxReq.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte("testUser:testPassword")))
	testCtxReq := httptest.NewRequest(http.MethodPost, "https://test.com/token?oauthContext=testContext", nil)
	testCtxReq.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte("testUser:testPassword")))
	invalidCtxReq := httptest.NewRequest(http.MethodPost, "https://test.com/token?oauthContext=invalidContext", nil)
	invalidCtxReq.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte("testUser:testPassword")))

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"default context",
			[]string{},
			[]string{actCheckNoError},
			&condition{
				r: defaultCtxReq,
				redeemer: &testTokenRedeemer{
					status: http.StatusOK,
					response: &TokenResponse{
						StatusCode:  http.StatusOK,
						AccessToken: "testAccessToken",
						RawBody:     []byte("defaultResponseBody"),
					},
				},
			},
			&action{
				status: http.StatusOK,
				body:   "defaultResponseBody",
				query: map[string]string{
					"username": "testUser",
					"password": "testPassword",
					"scope":    "defaultClientScope",
				},
				err: nil,
			},
		),
		gen(
			"test context",
			[]string{},
			[]string{actCheckNoError},
			&condition{
				r: testCtxReq,
				redeemer: &testTokenRedeemer{
					status: http.StatusOK,
					response: &TokenResponse{
						StatusCode:  http.StatusOK,
						AccessToken: "testAccessToken",
						RawBody:     []byte("testResponseBody"),
					},
				},
			},
			&action{
				status: http.StatusOK,
				body:   "testResponseBody",
				query: map[string]string{
					"username": "testUser",
					"password": "testPassword",
					"scope":    "testClientScope",
				},
				err: nil,
			},
		),
		gen(
			"invalid context",
			[]string{},
			[]string{actCheckNoError},
			&condition{
				r:        invalidCtxReq,
				redeemer: &testTokenRedeemer{},
			},
			&action{
				status:     http.StatusUnauthorized,
				body:       "",
				query:      nil,
				err:        app.ErrAppAuthnAuthentication,
				errPattern: regexp.MustCompile(core.ErrPrefix + `authentication failed`),
			},
		),
		gen(
			"redeemer returns error",
			[]string{},
			[]string{actCheckNoError},
			&condition{
				r: defaultCtxReq,
				redeemer: &testTokenRedeemer{
					err: utilhttp.NewHTTPError(io.EOF, http.StatusUnauthorized),
				},
			},
			&action{
				status: http.StatusUnauthorized,
				body:   "",
				query: map[string]string{
					"username": "testUser",
					"password": "testPassword",
					"scope":    "defaultClientScope",
				},
				err:        io.EOF,
				errPattern: regexp.MustCompile(`EOF`),
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			eh := &testErrorHandler{}
			h := &ropcHandler{
				eh: eh,
				baseHandler: &baseHandler{
					contextQueryKey: "oauthContext",
					oauthCtxs: map[string]*oauthContext{
						"default": {
							lg:   log.GlobalLogger(log.DefaultLoggerName),
							name: "test-context",
							client: &client{
								id:     "defaultClientID",
								secret: "defaultClientSecret",
								scope:  "defaultClientScope",
							},
							provider: &provider{
								issuer:  "https://default.com/",
								tokenEP: "https://default.com/token",
							},
							tokenRedeemer: tt.C().redeemer,
						},
						"testContext": {
							lg:   log.GlobalLogger(log.DefaultLoggerName),
							name: "test-context",
							client: &client{
								id:       "testClientID",
								secret:   "testClientSecret",
								scope:    "testClientScope",
								audience: "test",
							},
							provider: &provider{
								issuer:  "https://test.com/",
								tokenEP: "https://test.com/token",
							},
							tokenRedeemer: tt.C().redeemer,
						},
					},
				},
			}

			w := httptest.NewRecorder()
			h.ServeHTTP(w, tt.C().r)
			b, _ := io.ReadAll(w.Result().Body)
			testutil.Diff(t, tt.A().body, string(b))
			testutil.Diff(t, tt.A().query, tt.C().redeemer.params)

			if tt.A().err != nil {
				testutil.Diff(t, true, eh.called)
				e := eh.err.(*utilhttp.HTTPError)
				testutil.DiffError(t, tt.A().err, tt.A().errPattern, e.Unwrap(), cmpopts.EquateErrors())
				testutil.Diff(t, tt.A().status, e.StatusCode())
			} else {
				testutil.Diff(t, false, eh.called)
				testutil.Diff(t, tt.A().status, w.Result().StatusCode)
			}
		})
	}
}

func TestROPCHandler_ServeAuthn(t *testing.T) {
	jh := newJWTHandler()
	expire := time.Now().Add(1 * time.Hour)
	validAT, _ := jh.TokenWithClaims(
		&atClaims{
			RegisteredClaims: &jwt.RegisteredClaims{
				Issuer:    "https://test.com/",
				ExpiresAt: jwt.NewNumericDate(expire),
			},
		})
	validATString, _ := jh.SignedString(validAT)
	validATMapClaims := jwt.MapClaims{}
	jh.ParseWithClaims(validATString, &validATMapClaims, []jwt.ParserOption{}...)
	// validATMap := map[string]any{"iss": "https://test.com/", "aud": []any{"test"}, "exp": float64(expire.Unix())}
	validATMapWithTrue := map[string]any{"iss": "https://test.com/", "exp": float64(expire.Unix()), "active": true}

	invalidAT, _ := jh.TokenWithClaims(
		&atClaims{
			RegisteredClaims: &jwt.RegisteredClaims{
				Issuer: "https://invalid.com/",
				// Audience:  jwt.ClaimStrings{"testAudience"},
				ExpiresAt: jwt.NewNumericDate(expire),
			},
		})
	invalidATString, _ := jh.SignedString(invalidAT)

	type condition struct {
		redeemer *testTokenRedeemer
		r        *http.Request
		tokens   *OAuthTokens
	}

	type action struct {
		tokens     *OAuthTokens
		body       []byte
		ia         app.AuthResult
		sr         bool
		err        any // error or errorutil.Kind
		errPattern *regexp.Regexp
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	// condRedeemPath := tb.Condition("authentication via redeem token endpoint", "authentication via redeem token endpoint")
	// condTokenRequestFailed := tb.Condition("failed to request token", "failed to request token")
	// condValidClaimsFailed := tb.Condition("failed to valid claims", "failed to valid claims")
	// condExistTokensInSession := tb.Condition("exist tokens in the session", "exist tokens in the session")
	// actError := tb.Action("error", "check that the expected error is returned")
	actNoError := tb.Action("no error", "check that the there is no error")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"valid AT in session",
			[]string{},
			[]string{actNoError},
			&condition{
				redeemer: &testTokenRedeemer{
					status: http.StatusOK,
					response: &TokenResponse{
						AccessToken: validATString,
					},
				},
				r: newReqHeaderCredentials("/token", base64.StdEncoding.EncodeToString([]byte("testUsername:testPassword"))),
				tokens: &OAuthTokens{
					Context:  "default",
					AT:       validATString,
					ATClaims: validATMapClaims,
				},
			},
			&action{
				tokens: &OAuthTokens{
					Context:  "default",
					AT:       validATString,
					ATExp:    expire.Unix(),
					ATClaims: validATMapWithTrue,
				},
				ia:  app.AuthSucceeded,
				sr:  false,
				err: nil,
			},
		),
		gen(
			"valid AT in session / use non-default context",
			[]string{},
			[]string{actNoError},
			&condition{
				redeemer: &testTokenRedeemer{
					status: http.StatusOK,
					response: &TokenResponse{
						AccessToken: validATString,
					},
				},
				r: newReqHeaderCredentials("/token", base64.StdEncoding.EncodeToString([]byte("testUsername:testPassword"))),
				tokens: &OAuthTokens{
					Context:  "default",
					AT:       validATString,
					ATClaims: validATMapClaims,
				},
			},
			&action{
				tokens: &OAuthTokens{
					Context:  "default",
					AT:       validATString,
					ATExp:    expire.Unix(),
					ATClaims: validATMapWithTrue,
				},
				ia:  app.AuthSucceeded,
				sr:  false,
				err: nil,
			},
		),
		gen(
			"no AT in session",
			[]string{},
			[]string{actNoError},
			&condition{
				redeemer: &testTokenRedeemer{
					status: http.StatusOK,
					response: &TokenResponse{
						AccessToken: validATString,
					},
				},
				r: newReqHeaderCredentials("/token", base64.StdEncoding.EncodeToString([]byte("testUsername:testPassword"))),
				tokens: &OAuthTokens{
					Context: "default",
				},
			},
			&action{
				tokens: &OAuthTokens{
					Context:  "default",
					AT:       validATString,
					ATExp:    expire.Unix(),
					ATClaims: validATMapWithTrue,
				},
				ia:  app.AuthSucceeded,
				sr:  false,
				err: nil,
			},
		),
		gen(
			"valid AT in session / redeem new AT",
			[]string{},
			[]string{actNoError},
			&condition{
				redeemer: &testTokenRedeemer{
					status: http.StatusOK,
					response: &TokenResponse{
						StatusCode:  http.StatusOK,
						AccessToken: validATString,
					},
				},
				r: newReqHeaderCredentials("/redeem", base64.StdEncoding.EncodeToString([]byte("testUsername:testPassword"))),
				tokens: &OAuthTokens{
					Context:  "default",
					AT:       validATString,
					ATClaims: validATMapClaims,
				},
			},
			&action{
				tokens: &OAuthTokens{
					Context:  "default",
					AT:       validATString,
					ATExp:    expire.Unix(),
					ATClaims: validATMapWithTrue,
				},
				ia:  app.AuthContinue,
				sr:  true,
				err: nil,
			},
		),
		gen(
			"no AT in session / redeem new AT",
			[]string{},
			[]string{actNoError},
			&condition{
				redeemer: &testTokenRedeemer{
					status: http.StatusOK,
					response: &TokenResponse{
						StatusCode:  http.StatusOK,
						AccessToken: validATString,
					},
				},
				r: newReqHeaderCredentials("/redeem", base64.StdEncoding.EncodeToString([]byte("testUsername:testPassword"))),
				tokens: &OAuthTokens{
					Context: "default",
				},
			},
			&action{
				tokens: &OAuthTokens{
					Context:  "default",
					AT:       validATString,
					ATExp:    expire.Unix(),
					ATClaims: validATMapWithTrue,
				},
				ia:  app.AuthSucceeded,
				sr:  true,
				err: nil,
			},
		),
		gen(
			"no session",
			[]string{},
			[]string{actNoError},
			&condition{
				redeemer: &testTokenRedeemer{
					status: http.StatusOK,
					response: &TokenResponse{
						AccessToken: validATString,
					},
				},
				r:      newReqHeaderCredentials("/token", base64.StdEncoding.EncodeToString([]byte("testUsername:testPassword"))),
				tokens: nil, // No session.
			},
			&action{
				tokens: nil,
				ia:     app.AuthContinue,
				sr:     false,
			},
		),
		gen(
			"invalid AT",
			[]string{},
			[]string{actNoError},
			&condition{
				redeemer: &testTokenRedeemer{
					status: http.StatusOK,
					response: &TokenResponse{
						Error: "invalid access token",
					},
				},
				r: newReqHeaderCredentials("/token", base64.StdEncoding.EncodeToString([]byte("testUsername:testPassword"))),
				tokens: &OAuthTokens{
					Context:  "default",
					AT:       "INVALID_ACCESS_TOKEN",
					ATClaims: validATMapClaims,
				},
			},
			&action{
				tokens: nil,
				ia:     app.AuthFailed,
				sr:     true,
				err:    reAuthenticationRequired,
			},
		),
		gen(
			"expired AT in session",
			[]string{},
			[]string{actNoError},
			&condition{
				redeemer: &testTokenRedeemer{
					response: &TokenResponse{},
				},
				r: newReqHeaderCredentials("/token", base64.StdEncoding.EncodeToString([]byte("testUsername:testPassword"))),
				tokens: &OAuthTokens{
					Context: "default",
					ATExp:   time.Now().Unix() - 10,
					AT:      invalidATString,
				},
			},
			&action{
				tokens: nil,
				ia:     app.AuthFailed,
				sr:     true,
				err:    reAuthenticationRequired,
			},
		),
		gen(
			"no tokens in session / oauth context not found",
			[]string{},
			[]string{actNoError},
			&condition{
				redeemer: &testTokenRedeemer{
					status:   http.StatusOK,
					response: &TokenResponse{},
				},
				r: newReqHeaderCredentials("/token?testContextQueryKey=NotFound", base64.StdEncoding.EncodeToString([]byte("testUsername:testPassword"))),
				tokens: &OAuthTokens{
					Context: "default",
				},
			},
			&action{
				tokens: nil,
				ia:     app.AuthContinue,
				sr:     false,
			},
		),
		gen(
			"no tokens in session / token request error",
			[]string{},
			[]string{actNoError},
			&condition{
				redeemer: &testTokenRedeemer{
					status: http.StatusOK,
					response: &TokenResponse{
						Error: "token request failed",
					},
				},
				r: newReqHeaderCredentials("/token", base64.StdEncoding.EncodeToString([]byte("testUsername:testPassword"))),
				tokens: &OAuthTokens{
					Context: "default",
				},
			},
			&action{
				tokens: nil,
				ia:     app.AuthFailed,
				sr:     true,
				err:    reAuthenticationRequired,
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			eh := &testErrorHandler{}
			h := &ropcHandler{
				eh:         eh,
				redeemPath: "/redeem",
				baseHandler: &baseHandler{
					lg:               log.GlobalLogger(log.DefaultLoggerName),
					contextHeaderKey: "testContextHeaderKey",
					contextQueryKey:  "testContextQueryKey",
					oauthCtxs: map[string]*oauthContext{
						"default": {
							lg:   log.GlobalLogger(log.DefaultLoggerName),
							name: "default",
							client: &client{
								id:       "defaultClientID",
								secret:   "defaultClientSecret",
								scope:    "defaultClientScope",
								audience: "defaultAudience",
							},
							provider: &provider{
								issuer:  "https://test.com/",
								tokenEP: "https://test.com/token",
							},
							tokenRedeemer: tt.C().redeemer,
							jh:            jh,
							atParseOpts:   []jwt.ParserOption{jwt.WithIssuer("https://test.com/")},
							idtParseOpts:  []jwt.ParserOption{jwt.WithIssuer("https://test.com/")},
						},
						"testContext": {
							lg:   log.GlobalLogger(log.DefaultLoggerName),
							name: "testContext",
							client: &client{
								id:       "testClientID",
								secret:   "testClientSecret",
								scope:    "testClientScope",
								audience: "testAudience",
							},
							provider: &provider{
								issuer:  "https://test.com/",
								tokenEP: "https://test.com/token",
							},
							tokenRedeemer: tt.C().redeemer,
							jh:            jh,
							claimsKey:     "testClaimsKey",
							atParseOpts:   []jwt.ParserOption{jwt.WithIssuer("https://test.com/")},
							idtParseOpts:  []jwt.ParserOption{jwt.WithIssuer("https://test.com/")},
						},
					},
				},
			}

			r := tt.C().r
			if tt.C().tokens != nil {
				ss := session.NewDefaultSession(session.SerializeJSON)
				session.MustPersist(ss, ropcSessionKey, tt.C().tokens)
				ctx := session.ContextWithSession(tt.C().r.Context(), ss)
				r = r.WithContext(ctx)
			}

			w := httptest.NewRecorder()
			nr, ia, sr, err := h.ServeAuthn(w, r)
			t.Logf("%#v\n", err)
			t.Logf("%#v, %#v\n", nr, tt.A().tokens)
			if tt.A().tokens != nil && tt.A().tokens != (*OAuthTokens)(nil) {
				newSess := session.SessionFromContext(nr.Context())
				token := &OAuthTokens{}
				err := newSess.Extract(ropcSessionKey, token)
				testutil.Diff(t, tt.A().tokens, token, cmp.AllowUnexported(OAuthTokens{}))
				testutil.Diff(t, nil, err)
			} else {
				testutil.Diff(t, (*http.Request)(nil), nr)
			}
			testutil.Diff(t, tt.A().ia, ia) // authenticated
			testutil.Diff(t, tt.A().sr, sr) // shouldReturn
			testutil.DiffError(t, tt.A().err, tt.A().errPattern, err, cmpopts.EquateErrors())
		})
	}
}

func TestROPCHandler_TokenRequest(t *testing.T) {
	type condition struct {
		usernameKey string
		passwordKey string
		r           *http.Request
		redeemer    *testTokenRedeemer
	}

	type action struct {
		resp       *TokenResponse
		query      map[string]string
		err        any // error or errorutil.Kind
		errPattern *regexp.Regexp
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	cndForm := tb.Condition("form", "credentials in form request body")
	cndHeader := tb.Condition("header", "credentials in authorization header")
	cndEmptyUsername := tb.Condition("empty usernameKey", "usernameKey is empty")
	cndEmptyPassword := tb.Condition("empty passwordKey", "passwordKey is empty")
	cndInsufficientCredential := tb.Condition("insufficient credentials", "provided username or password is invalid")
	cndInvalidAuthHeader := tb.Condition("invalid auth header", "invalid authorization header value")
	actCheckError := tb.Action("error", "check that the expected error is returned")
	actCheckNoError := tb.Action("no error", "check that the there is no error")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"credential in form body",
			[]string{cndForm},
			[]string{actCheckNoError},
			&condition{
				usernameKey: "username",
				passwordKey: "password",
				r:           newReqFormBodyCredentials("/token", bytes.NewReader([]byte("username=testUsername&password=testPassword"))),
				redeemer: &testTokenRedeemer{
					status: http.StatusOK,
					response: &TokenResponse{
						AccessToken: "testAccessToken",
					},
				},
			},
			&action{
				resp: &TokenResponse{
					AccessToken: "testAccessToken",
				},
				query: map[string]string{
					"username":  "testUsername",
					"password":  "testPassword",
					"scope":     "testClientScope",
					"testParam": "testValue",
				},
			},
		),
		gen(
			"credential in authorization header/empty username key",
			[]string{cndHeader, cndEmptyUsername},
			[]string{actCheckNoError},
			&condition{
				usernameKey: "",
				passwordKey: "password",
				r:           newReqHeaderCredentials("/token", base64.StdEncoding.EncodeToString([]byte("testUsername:testPassword"))),
				redeemer: &testTokenRedeemer{
					status: http.StatusOK,
					response: &TokenResponse{
						AccessToken: "testAccessToken",
					},
				},
			},
			&action{
				resp: &TokenResponse{
					AccessToken: "testAccessToken",
				},
				query: map[string]string{
					"username":  "testUsername",
					"password":  "testPassword",
					"scope":     "testClientScope",
					"testParam": "testValue",
				},
			},
		),
		gen(
			"credential in authorization header/empty password key",
			[]string{cndHeader, cndEmptyPassword},
			[]string{actCheckNoError},
			&condition{
				usernameKey: "username",
				passwordKey: "",
				r:           newReqHeaderCredentials("/token", base64.StdEncoding.EncodeToString([]byte("testUsername:testPassword"))),
				redeemer: &testTokenRedeemer{
					status: http.StatusOK,
					response: &TokenResponse{
						AccessToken: "testAccessToken",
					},
				},
			},
			&action{
				resp: &TokenResponse{
					AccessToken: "testAccessToken",
				},
				query: map[string]string{
					"username":  "testUsername",
					"password":  "testPassword",
					"scope":     "testClientScope",
					"testParam": "testValue",
				},
			},
		),
		gen(
			"credential in authorization header/empty username&password keys",
			[]string{cndHeader, cndEmptyUsername, cndEmptyPassword},
			[]string{actCheckNoError},
			&condition{
				usernameKey: "",
				passwordKey: "",
				r:           newReqHeaderCredentials("/token", base64.StdEncoding.EncodeToString([]byte("testUsername:testPassword"))),
				redeemer: &testTokenRedeemer{
					status: http.StatusOK,
					response: &TokenResponse{
						AccessToken: "testAccessToken",
					},
				},
			},
			&action{
				resp: &TokenResponse{
					AccessToken: "testAccessToken",
				},
				query: map[string]string{
					"username":  "testUsername",
					"password":  "testPassword",
					"scope":     "testClientScope",
					"testParam": "testValue",
				},
			},
		),
		gen(
			"body read error",
			[]string{cndForm},
			[]string{actCheckError},
			&condition{
				usernameKey: "username",
				passwordKey: "password",
				r:           newReqFormBodyCredentials("/token", &testutil.ErrorReader{}),
				redeemer: &testTokenRedeemer{
					status: http.StatusOK,
					response: &TokenResponse{
						AccessToken: "testAccessToken",
					},
				},
			},
			&action{
				resp:       nil,
				query:      nil,
				err:        app.ErrAppGenReadHTTPBody,
				errPattern: regexp.MustCompile(core.ErrPrefix + `failed to read request body.`),
			},
		),
		gen(
			"username not found in form body",
			[]string{cndForm, cndInsufficientCredential},
			[]string{actCheckError},
			&condition{
				usernameKey: "username",
				passwordKey: "password",
				r:           newReqFormBodyCredentials("/token", bytes.NewReader([]byte("password=testPassword"))),
				redeemer: &testTokenRedeemer{
					status: http.StatusOK,
					response: &TokenResponse{
						AccessToken: "testAccessToken",
					},
				},
			},
			&action{
				resp:       nil,
				query:      nil,
				err:        app.ErrAppAuthnInvalidCredential,
				errPattern: regexp.MustCompile(core.ErrPrefix + `invalid credential for ROPC. username or password is empty`),
			},
		),
		gen(
			"password not found in form body",
			[]string{cndForm, cndInsufficientCredential},
			[]string{actCheckError},
			&condition{
				usernameKey: "username",
				passwordKey: "password",
				r:           newReqFormBodyCredentials("/token", bytes.NewReader([]byte("username=testUsername"))),
				redeemer: &testTokenRedeemer{
					status: http.StatusOK,
					response: &TokenResponse{
						AccessToken: "testAccessToken",
					},
				},
			},
			&action{
				resp:       nil,
				query:      nil,
				err:        app.ErrAppAuthnInvalidCredential,
				errPattern: regexp.MustCompile(core.ErrPrefix + `invalid credential for ROPC. username or password is empty`),
			},
		),
		gen(
			"authorization header not found",
			[]string{cndHeader, cndInvalidAuthHeader},
			[]string{actCheckError},
			&condition{
				usernameKey: "",
				passwordKey: "",
				r:           newReqHeaderCredentials("/token", ""),
				redeemer: &testTokenRedeemer{
					status: http.StatusOK,
					response: &TokenResponse{
						AccessToken: "testAccessToken",
					},
				},
			},
			&action{
				resp:       nil,
				query:      nil,
				err:        app.ErrAppAuthnInvalidCredential,
				errPattern: regexp.MustCompile(core.ErrPrefix + `invalid credential for ROPC. invalid authorization header`),
			},
		),
		gen(
			"empty username",
			[]string{cndHeader, cndInvalidAuthHeader},
			[]string{actCheckError},
			&condition{
				usernameKey: "",
				passwordKey: "",
				r:           newReqHeaderCredentials("/token", base64.StdEncoding.EncodeToString([]byte(":testPassword"))),
				redeemer: &testTokenRedeemer{
					status: http.StatusOK,
					response: &TokenResponse{
						AccessToken: "testAccessToken",
					},
				},
			},
			&action{
				resp:       nil,
				query:      nil,
				err:        app.ErrAppAuthnInvalidCredential,
				errPattern: regexp.MustCompile(core.ErrPrefix + `invalid credential for ROPC. username or password is empty`),
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			oc := &oauthContext{
				lg:   log.GlobalLogger(log.DefaultLoggerName),
				name: "test-context",
				client: &client{
					id:     "testClientID",
					secret: "testClientSecret",
					scope:  "testClientScope",
				},
				provider: &provider{
					issuer:  "https://test.com/",
					tokenEP: "https://test.com/token",
				},
				tokenRedeemer: tt.C().redeemer,
			}

			h := &ropcHandler{
				usernameKey: tt.C().usernameKey,
				passwordKey: tt.C().passwordKey,
				queryParams: map[string]string{"testParam": "testValue"},
			}

			resp, err := h.tokenRequest(tt.C().r, oc)
			testutil.DiffError(t, tt.A().err, tt.A().errPattern, err)

			testutil.Diff(t, tt.A().resp, resp)
			testutil.Diff(t, tt.A().query, tt.C().redeemer.params)
		})
	}
}

func newReqFormBodyCredentials(path string, body io.Reader) *http.Request {
	req := httptest.NewRequest(http.MethodPost, "https://test.com"+path, body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return req
}

func newReqHeaderCredentials(path string, header string) *http.Request {
	req := httptest.NewRequest(http.MethodPost, "https://test.com"+path, nil)
	if header != "" {
		req.Header.Set("Authorization", "Basic "+header)
	}
	return req
}
