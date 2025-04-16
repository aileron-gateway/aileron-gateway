// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package oauth

/*
func TestNewAuthorizationCodeHandler(t *testing.T) {

	type condition struct {
		bh   *baseHandler
		spec *v1.AuthorizationCodeHandler
	}

	type action struct {
		h          *authorizationCodeHandler
		err        any
		errPattern *regexp.Regexp
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"empty spec",
			[]string{},
			[]string{},
			&condition{
				bh:   nil,
				spec: &v1.AuthorizationCodeHandler{},
			},
			&action{
				h: &authorizationCodeHandler{
					csrf: &csrfStateGenerator{
						method: "S256",
					},
					redirectPathPattern: regexp.MustCompile("^$"),
				},
				err: nil,
			},
		),
		gen(
			"no error",
			[]string{},
			[]string{},
			&condition{
				bh: nil,
				spec: &v1.AuthorizationCodeHandler{
					DisableState:        true,
					DisableNonce:        true,
					DisablePKCE:         true,
					PKCEMethod:          v1.PKCEMethod_Plain,
					LoginPath:           "/login",
					CallbackURL:         "https://test.com/callback",
					RedirectPath:        "/redirect",
					RedirectKey:         "redirectKey",
					RedirectPathPattern: `^[a-z]+\[[0-9]+\]$`,
					RedirectToLogin:     true,
					UnauthorizeAny:      true,
					RestoreRequest:      true,
					URLParams:           []string{"key1=value1", "key2=value2"},
				},
			},
			&action{
				h: &authorizationCodeHandler{
					loginPath:           "/login",
					redirectPath:        "/redirect",
					redirectPathPattern: regexp.MustCompile(`^[a-z]+\[[0-9]+\]$`),
					callbackPath:        "/callback",
					callbackURL:         "https://test.com/callback",
					redirectToLogin:     true,
					unauthorizeAny:      true,
					restoreRequest:      true,
					csrf: &csrfStateGenerator{
						stateDisabled: true,
						nonceDisabled: true,
						pkceDisabled:  true,
						method:        "plain",
					},
					urlParams: "key1=value1&key2=value2",
				},
				err: nil,
			},
		),
		gen(
			"invalid callback url",
			[]string{},
			[]string{},
			&condition{
				bh: nil,
				spec: &v1.AuthorizationCodeHandler{
					CallbackURL: "https://test.com/\r\n", // Invalid URL.
				},
			},
			&action{
				h:          nil,
				err:        core.ErrCoreGenCreateObject,
				errPattern: regexp.MustCompile(core.ErrPrefix + `failed to create OAuthAuthenticationHandler`),
			},
		),
	}

	testutil.Register(table, testCases...)


	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {

			h, err := newAuthorizationCodeHandler(tt.C().bh, tt.C().spec)
			testutil.DiffError(t, tt.A().err, tt.A().errPattern, err)

			opts := []cmp.Option{
				cmpopts.IgnoreInterfaces(struct{ http.RoundTripper }{}),
				cmp.AllowUnexported(authorizationCodeHandler{}, redeemTokenClient{}, provider{}, tokenIntrospectionClient{}),
				cmp.AllowUnexported(security.JWTHandler{}, security.SigningKey{}, csrfStateGenerator{}),
				cmp.AllowUnexported(atomic.Bool{}, regexp.Regexp{}),
				cmpopts.IgnoreUnexported(sync.RWMutex{}, sync.Mutex{}, http.Transport{}),
			}

			testutil.Diff(t, tt.A().h, h, opts...)

		})
	}
}

func TestAuthorizationCodeHandler_ServeAuthn(t *testing.T) {

	testClient := &client{
		id:     "testClientID",
		secret: "testClientSecret",
	}
	testProvider := &provider{
		issuer:          "https://test.com/",
		authorizationEP: "https://test.com/auth",
		tokenEP:         "https://test.com/token",
	}

	// JWT handler.
	jh := newJWTHandler()

	// AT.
	at, _ := jh.TokenWithClaims(
		&atClaims{
			RegisteredClaims: &jwt.RegisteredClaims{
				Issuer:    "https://test.com/",
				Audience:  jwt.ClaimStrings{"resourceServer"},
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
			},
		})
	atStr, _ := jh.SignedString(at)
	atmc := jwt.MapClaims{}
	jh.ParseWithClaims(atStr, &atmc, []jwt.ParserOption{}...)

	type condition struct {
		h                *authorizationCodeHandler
		rt               tokenRedeemer
		session          *OAuthTokens
		loginPath        string
		sessionNotExists bool
		restoreRequest   bool
		persistRequest   bool
		redirectPath     any
		r                *http.Request
	}

	type action struct {
		ctx     *OAuthClaims
		session *OAuthTokens
		ia      bool
		sr      bool
		err     core.HTTPError
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"url path is login path or callback path when claims exist in session",
			[]string{},
			[]string{},
			&condition{
				h: &authorizationCodeHandler{
					baseHandler: &baseHandler{
						lg: log.GlobalLogger(log.DefaultLoggerName),
						eh: utilhttp.GlobalErrorHandler(utilhttp.DefaultErrorHandlerName),
						oauthCtxs: map[string]*oauthContext{"default": {
							name:     "default",
							provider: testProvider,
							client:   testClient,
						}},
					},
				},
				loginPath: "/login",
				session:   &OAuthTokens{AT: atStr},
				r:         httptest.NewRequest(http.MethodPost, "https://test.com/login", nil),
			},
			&action{
				ctx:     nil,
				session: nil,
				ia:      false,
				sr:      true,
				err:     utilhttp.NewHTTPError(nil, http.StatusForbidden, "application/json", []byte(`{"status":"Forbidden"}`)),
			},
		),
		// gen(
		// 	"invalid claims when claims exist in session",
		// 	[]string{condSessionExists, condInvalidClaims},
		// 	[]string{actError},
		// 	&condition{
		// 		session: &OAuthTokens{AT: ""},
		// 		r:       httptest.NewRequest(http.MethodPost, "https://test.com/", nil),
		// 	},
		// 	&action{
		// 		ctx:     nil,
		// 		session: nil,
		// 		ia:      false,
		// 		sr:      true,
		// 		err:     utilhttp.NewHTTPError(base.ErrAppAuthnInvalidToken.WithStack(nil, map[string]any{"name": "access token"}), http.StatusUnauthorized, "application/json", []byte(`{"status":"Unauthorized"}`)),
		// 	},
		// ),
		// gen(
		// 	"claims exist in session",
		// 	[]string{condSessionExists},
		// 	[]string{actNoError},
		// 	&condition{
		// 		session: &OAuthTokens{AT: atStr},
		// 		r:       httptest.NewRequest(http.MethodPost, "https://test.com/", nil),
		// 	},
		// 	&action{
		// 		ctx: &OAuthClaims{
		// 			Method: acMethod,
		// 			AT:     atmc,
		// 		},
		// 		session: &OAuthTokens{AT: atStr},
		// 		ia:      true,
		// 		sr:      false,
		// 		err:     nil,
		// 	},
		// ),
		// gen(
		// 	"url path is not login path or callback path when login path exists",
		// 	[]string{condNotCallbackPath, condNotLoginPath},
		// 	[]string{actError},
		// 	&condition{
		// 		loginPath: "/login",
		// 		r:         httptest.NewRequest(http.MethodPost, "https://test.com/", nil),
		// 	},
		// 	&action{
		// 		ctx:     nil,
		// 		session: nil,
		// 		ia:      false,
		// 		sr:      true,
		// 		err:     utilhttp.NewHTTPError(nil, http.StatusUnauthorized, "application/json", []byte(`{"status":"Unauthorized"}`)),
		// 	},
		// ),
		// gen(
		// 	"failed to authorization request when url path is not callback path",
		// 	[]string{condNotCallbackPath, condAuthorizationRequestFailed},
		// 	[]string{actError},
		// 	&condition{
		// 		sessionNotExists: true,
		// 		r:                httptest.NewRequest(http.MethodPost, "https://test.com/", nil),
		// 	},
		// 	&action{
		// 		ctx:     nil,
		// 		session: nil,
		// 		ia:      false,
		// 		sr:      true,
		// 		err:     utilhttp.NewHTTPError(base.ErrSessionNotExists.WithStack(nil, nil), http.StatusUnauthorized, "application/json", []byte(`{"status":"Unauthorized"}`)),
		// 	},
		// ),
		// gen(
		// 	"url path is not callback path",
		// 	[]string{condNotCallbackPath},
		// 	[]string{actNoError},
		// 	&condition{
		// 		r: httptest.NewRequest(http.MethodPost, "https://test.com/", nil),
		// 	},
		// 	&action{
		// 		ctx:     nil,
		// 		session: nil,
		// 		ia:      false,
		// 		sr:      true,
		// 		err:     nil,
		// 	},
		// ),
		// gen(
		// 	"failed to handle callback",
		// 	[]string{condHandleCallbackFailed},
		// 	[]string{actError},
		// 	&condition{
		// 		sessionNotExists: true,
		// 		r:                httptest.NewRequest(http.MethodPost, "https://test.com/callback", nil),
		// 	},
		// 	&action{
		// 		ctx:     nil,
		// 		session: nil,
		// 		ia:      false,
		// 		sr:      true,
		// 		err:     utilhttp.NewHTTPError(base.ErrSessionNotExists.WithStack(nil, nil), http.StatusUnauthorized, "application/json", []byte(`{"status":"Unauthorized"}`)),
		// 	},
		// ),
		// gen(
		// 	"failed to handle callback",
		// 	[]string{condInvalidClaims},
		// 	[]string{actError},
		// 	&condition{
		// 		rt: &testRedeemToken{
		// 			tr: &TokenResponseModel{
		// 				AccessToken: "",
		// 			},
		// 		},
		// 		r: httptest.NewRequest(http.MethodPost, "https://test.com/callback", nil),
		// 	},
		// 	&action{
		// 		ctx:     nil,
		// 		session: nil,
		// 		ia:      false,
		// 		sr:      true,
		// 		err:     utilhttp.NewHTTPError(base.ErrAppAuthnInvalidToken.WithStack(nil, map[string]any{"name": "access token"}), http.StatusUnauthorized, "application/json", []byte(`{"status":"Unauthorized"}`)),
		// 	},
		// ),
		// gen(
		// 	"failed to extract request if restore request",
		// 	[]string{condRestoreRequest, condExtractRequestFailed},
		// 	[]string{actError},
		// 	&condition{
		// 		rt: &testRedeemToken{
		// 			tr: &TokenResponseModel{
		// 				AccessToken: atStr,
		// 			},
		// 		},
		// 		restoreRequest: true,
		// 		r:              httptest.NewRequest(http.MethodPost, "https://test.com/callback", nil),
		// 	},
		// 	&action{
		// 		ctx:     nil,
		// 		session: nil,
		// 		ia:      false,
		// 		sr:      true,
		// 		err:     utilhttp.NewHTTPError(errors.New("request not exist"), http.StatusInternalServerError, "", nil),
		// 	},
		// ),
		// gen(
		// 	"extract request if restore request",
		// 	[]string{condRestoreRequest},
		// 	[]string{actNoError},
		// 	&condition{
		// 		rt: &testRedeemToken{
		// 			tr: &TokenResponseModel{
		// 				AccessToken: atStr,
		// 			},
		// 		},
		// 		restoreRequest: true,
		// 		persistRequest: true,
		// 		r:              httptest.NewRequest(http.MethodPost, "https://test.com/callback", nil),
		// 	},
		// 	&action{
		// 		ctx: &OAuthClaims{
		// 			Method: acMethod,
		// 			AT:     atmc,
		// 		},
		// 		session: &OAuthTokens{AT: atStr},
		// 		ia:      true,
		// 		sr:      false,
		// 		err:     nil,
		// 	},
		// ),
		// gen(
		// 	"failed to extract redirect path in session",
		// 	[]string{condRedirectSessionExists, condExtractRedirectFailed},
		// 	[]string{actNoError},
		// 	&condition{
		// 		rt: &testRedeemToken{
		// 			tr: &TokenResponseModel{
		// 				AccessToken: atStr,
		// 			},
		// 		},
		// 		redirectPath: &OAuthClaims{},
		// 		r:            httptest.NewRequest(http.MethodPost, "https://test.com/callback", nil),
		// 	},
		// 	&action{
		// 		ctx: &OAuthClaims{
		// 			Method: acMethod,
		// 			AT:     atmc,
		// 		},
		// 		session: &OAuthTokens{AT: atStr},
		// 		ia:      false,
		// 		sr:      true,
		// 		err:     utilhttp.NewHTTPError(errors.New("(E1133.ErrUnmarshal)failed to unmarshal from=msgpack to=any [msgpack: invalid code=8b decoding string/bytes length]"), http.StatusForbidden, "application/json", []byte(`{"status":"Forbidden"}`)),
		// 	},
		// ),
		// gen(
		// 	"extract redirect session",
		// 	[]string{condRedirectSessionExists},
		// 	[]string{actNoError},
		// 	&condition{
		// 		rt: &testRedeemToken{
		// 			tr: &TokenResponseModel{
		// 				AccessToken: atStr,
		// 			},
		// 		},
		// 		r: httptest.NewRequest(http.MethodPost, "https://test.com/callback", nil),
		// 	},
		// 	&action{
		// 		ctx:     nil,
		// 		session: nil,
		// 		ia:      true,
		// 		sr:      true,
		// 		err:     nil,
		// 	},
		// ),
	}

	testutil.Register(table, testCases...)


	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			h := &authorizationCodeHandler{
				skipper:        skipper,
				tokenRedeemer:  tt.C().rt,
				client:         cl,
				provider:       pv,
				claimsKey:      "oauthClaims",
				resourceServer: "resourceServer",
				jh:             jh,
				csrf:           &csrfStateGenerator{},
				loginPath:      tt.C().loginPath,
				callbackPath:   "/callback",
				restoreRequest: tt.C().restoreRequest,
			}
			w := httptest.NewRecorder()
			r := tt.C().r
			// Create session.
			if !tt.C().sessionNotExists {
				ss, _ := session.NewSession()
				ctx := session.ContextWithSession(r.Context(), ss)
				r = r.WithContext(ctx)

				if tt.C().session != nil {
					// Persist tokens in the session.
					authn.PersistClaims(r.Context(), h.claimsKey, acMethod, tt.C().session)
				}
				if tt.C().persistRequest {
					persistRequest(r)
				}
				if tt.C().redirectPath != nil {
					ss.Persist(redirectSessionKey, tt.C().redirectPath, session.SerializeMsgPack)
				}
			}

			nr, ia, sr, err := h.ServeAuthn(w, r)

			// Check for error.
			if tt.A().err != nil {
				e := err.(core.ErrorResponse)
				// testutil.Diff(t, tt.A().err.Error(), e.Error())
				testutil.Diff(t, tt.A().err.StatusCode(), e.StatusCode())
				testutil.Diff(t, tt.A().err.Body(), e.Body())
			} else {
				testutil.Diff(t, nil, err)
			}

			// Check isAuthenticated.
			testutil.Diff(t, tt.A().ia, ia)
			// Check shouldReturn.
			testutil.Diff(t, tt.A().sr, sr)

			// Check if newRequest is nil.
			if tt.A().ctx == nil && tt.A().session == nil {
				testutil.Diff(t, (*http.Request)(nil), nr)
				return
			}

			// Check the context of the newRequest.
			opts := []cmp.Option{
				cmpopts.IgnoreFields(OAuthClaims{}, "AuthTime"),
			}
			testutil.Diff(t, tt.A().ctx, nr.Context().Value(h.claimsKey), opts...)

			// Check the session.
			tokens, _ := authn.ExtractClaims[OAuthTokens](nr.Context(), h.claimsKey, acMethod)
			testutil.Diff(t, tt.A().session, tokens)
		})
	}
}
*/

// func TestAuthorizationRequest(t *testing.T) {
// 	authEP, _ := url.Parse("https://test.com/auth")

// 	// Provider.
// 	pv := &provider{
// 		Issuer:          "https://test.com/",
// 		AuthorizationEP: authEP,
// 	}
// 	// Client.
// 	cl := &client{
// 		ID:     "testClientID",
// 		Secret: "testClientSecret",
// 		Scope:  "scope1 scope2 scope3",
// 	}

// 	// Response writer.
// 	w := httptest.NewRecorder()
// 	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
// 	w.WriteHeader(http.StatusFound)
// 	w.Write(nil)

// 	type condition struct {
// 		sessionNotExists bool
// 		r                *http.Request
// 		restoreRequest   bool
// 		w                *httptest.ResponseRecorder
// 	}

// 	type action struct {
// 		err        error
// 		errPattern *regexp.Regexp
// 	}

// 	condSessionNotExists := "session does not exists"
// 	condRestoreRequest := "restore request"
// 	condPersistRequestFailed := "failed to persist request in the session"
// 	condGetPathFromQuery := "get the path to redirect from the query"
// 	condGetPathFromRequestURI := "get the path to redirect from the request uri"

// 	actError := "error"
// 	actNoError := "no error"

// 	tb := testutil.NewTableBuilder[*condition, *action]()
// 	tb.Name(t.Name())
// 	tb.Condition(condSessionNotExists, "session does not exists")
// 	tb.Condition(condRestoreRequest, "restore request")
// 	tb.Condition(condPersistRequestFailed, "failed to persist request in the session")
// 	tb.Condition(condGetPathFromQuery, "get the path to redirect from the query")
// 	tb.Condition(condGetPathFromRequestURI, "get the path to redirect from the request uri")
// 	tb.Action(actError, "check that the expected error is returned")
// 	tb.Action(actNoError, "check that the there is no error")

// 	table := tb.Build()

// 	gen := testutil.NewCase[*condition, *action]
// 	testCases := []*testutil.Case[*condition, *action]{
// 		gen(
// 			"session dose not exists",
// 			[]string{condSessionNotExists},
// 			[]string{actError},
// 			&condition{
// 				sessionNotExists: true,
// 				r:                httptest.NewRequest(http.MethodPost, "https://test.com/", nil),
// 				w:                httptest.NewRecorder(),
// 			},
// 			&action{
// 				err:        base.ErrSessionNotExists,
// 				errPattern: regexp.MustCompile(core.ErrPrefix + `session does not exists`),
// 			},
// 		),
// 		gen(
// 			"failed to persist request in the session",
// 			[]string{condRestoreRequest, condPersistRequestFailed},
// 			[]string{actError},
// 			&condition{
// 				r:              newRequestWithCredential(httptest.NewRequest(http.MethodPost, "https://test.com/", &testutil.ErrorReader{}), ""),
// 				restoreRequest: true,
// 				w:              httptest.NewRecorder(),
// 			},
// 			&action{
// 				err:        base.ErrPersistRequestFailed,
// 				errPattern: regexp.MustCompile(core.ErrPrefix + `failed to persist request in the session`),
// 			},
// 		),
// 		gen(
// 			"get the path to redirect from the query",
// 			[]string{condGetPathFromQuery},
// 			[]string{actNoError},
// 			&condition{
// 				r: httptest.NewRequest(http.MethodPost, "https://test.com/?rd=https://test.com/redirect", nil),
// 				w: w,
// 			},
// 			&action{
// 				err: nil,
// 			},
// 		),
// 		gen(
// 			"get the path to redirect from the request uri",
// 			[]string{condGetPathFromRequestURI},
// 			[]string{actNoError},
// 			&condition{
// 				r: httptest.NewRequest(http.MethodPost, "https://test.com/", nil),
// 				w: w,
// 			},
// 			&action{
// 				err: nil,
// 			},
// 		),
// 	}

// 	testutil.Register(table, testCases...)
//

// 	for _, tt := range table.Entries() {
// 		tt := tt
// 		t.Run(tt.Name(), func(t *testing.T) {
// 			h := &authorizationCodeHandler{
// 				provider:       pv,
// 				client:         cl,
// 				callbackURL:    "https://test.com/callback",
// 				csrf:           &csrfStateGenerator{method: "S256"},
// 				restoreRequest: tt.C().restoreRequest,
// 				loginPath:      "",
// 				extraParams:    map[string]string{"key1": "value1", "key2": "value2"},
// 			}

// 			w := httptest.NewRecorder()

// 			// Create request with session.
// 			r := tt.C().r
// 			if !tt.C().sessionNotExists {
// 				ss, _ := session.NewSession()
// 				ctx := session.ContextWithSession(r.Context(), ss)
// 				r = r.WithContext(ctx)

// 				// Persist csrf states in the session.
// 				ss.Persist(csrfSessionKey, &CSRFStates{}, session.SerializeMsgPack)
// 			}

// 			err := h.authorizationRequest(w, r)

// 			// Check for error.
// 			testutil.DiffError(t, tt.A().err, tt.A().errPattern, err)

// 			// Check response writer.
// 			testutil.Diff(t, w.Header().Get("Content-Type"), tt.C().w.Header().Get("Content-Type"))
// 			testutil.Diff(t, w.Result().StatusCode, tt.C().w.Result().StatusCode)
// 			testutil.Diff(t, w.Body.Bytes(), tt.C().w.Body.Bytes())

// 		})
// 	}
// }

// func TestHandleCallback(t *testing.T) {
// 	type condition struct {
// 		sessionNotExists bool
// 		csrfStates       any
// 		tr               tokenRedeemer
// 	}

// 	type action struct {
// 		t          *OAuthTokens
// 		nonce      string
// 		err        error
// 		errPattern *regexp.Regexp
// 	}

// 	condSessionNotExists := "session does not exists"
// 	condExtractCSRFParamsFailed := "failed to extract CSRF parameters"
// 	condInvalidState := "invalid state"
// 	condVerifierExists := "code verifier exists"

// 	actError := "error"
// 	actNoError := "no error"

// 	tb := testutil.NewTableBuilder[*condition, *action]()
// 	tb.Name(t.Name())
// 	tb.Condition(condSessionNotExists, "session does not exists")
// 	tb.Condition(condExtractCSRFParamsFailed, "failed to extract CSRF parameters")
// 	tb.Condition(condInvalidState, "invalid state")
// 	tb.Condition(condVerifierExists, "code verifier exists")
// 	tb.Action(actError, "check that the expected error is returned")
// 	tb.Action(actNoError, "check that the there is no error")

// 	table := tb.Build()

// 	gen := testutil.NewCase[*condition, *action]
// 	testCases := []*testutil.Case[*condition, *action]{
// 		gen(
// 			"session dose not exists",
// 			[]string{condSessionNotExists},
// 			[]string{actError},
// 			&condition{
// 				sessionNotExists: true,
// 			},
// 			&action{
// 				t:          nil,
// 				nonce:      "",
// 				err:        base.ErrSessionNotExists,
// 				errPattern: regexp.MustCompile(core.ErrPrefix + `session does not exists`),
// 			},
// 		),
// 		gen(
// 			"failed to extract CSRF parameters",
// 			[]string{condExtractCSRFParamsFailed},
// 			[]string{actError},
// 			&condition{
// 				csrfStates: "invalidCSRFStates",
// 			},
// 			&action{
// 				t:          nil,
// 				nonce:      "",
// 				err:        base.ErrSessionNotExists,
// 				errPattern: regexp.MustCompile(core.ErrPrefix + `failed to extract CSRF parameters`),
// 			},
// 		),
// 		gen(
// 			"invalid state",
// 			[]string{condInvalidState},
// 			[]string{actError},
// 			&condition{
// 				csrfStates: &CSRFStates{
// 					State: "invalidState",
// 				},
// 			},
// 			&action{
// 				t:          nil,
// 				nonce:      "",
// 				err:        base.ErrSessionNotExists,
// 				errPattern: regexp.MustCompile(core.ErrPrefix + `invalid state`),
// 			},
// 		),
// 		gen(
// 			"code verifier exists",
// 			[]string{condVerifierExists},
// 			[]string{actNoError},
// 			&condition{
// 				tr: &testRedeemToken{
// 					&TokenResponseModel{
// 						IDToken:      "IDT",
// 						AccessToken:  "AT",
// 						RefreshToken: "RT",
// 					},
// 				},
// 				csrfStates: &CSRFStates{
// 					State:    "testState",
// 					Verifier: "testVerifier",
// 					Nonce:    "testNonce",
// 				},
// 			},
// 			&action{
// 				t: &OAuthTokens{
// 					IDT: "IDT",
// 					AT:  "AT",
// 					RT:  "RT",
// 				},
// 				nonce: "testNonce",
// 				err:   nil,
// 			},
// 		),
// 	}

// 	testutil.Register(table, testCases...)
//

// 	for _, tt := range table.Entries() {
// 		tt := tt
// 		t.Run(tt.Name(), func(t *testing.T) {
// 			h := &authorizationCodeHandler{
// 				tokenRedeemer: tt.C().tr,
// 				callbackURL:   "https://test.com/callback",
// 				csrf:          &csrfStateGenerator{},
// 			}

// 			// Create request with session.
// 			r := httptest.NewRequest(http.MethodPost, "https://test.com/", nil)
// 			q := r.URL.Query()
// 			q.Set("state", "testState")
// 			q.Set("code", "testCode")
// 			r.URL.RawQuery = q.Encode()

// 			if !tt.C().sessionNotExists {
// 				ss, _ := session.NewSession()
// 				ctx := session.ContextWithSession(r.Context(), ss)
// 				r = r.WithContext(ctx)

// 				if tt.C().csrfStates != nil {
// 					// Persist csrf states in the session.
// 					ss.Persist(csrfSessionKey, tt.C().csrfStates, session.SerializeMsgPack)
// 				}
// 			}

// 			tokens, nonce, err := h.handleCallback(r)

// 			// Check for error.
// 			testutil.DiffError(t, tt.A().err, tt.A().errPattern, err)

// 			// Check tokens.
// 			testutil.Diff(t, tt.A().t, tokens)

// 			// Check nonce.
// 			testutil.Diff(t, tt.A().nonce, nonce)

// 		})
// 	}
// }

// func TestAuthorizationCodeValidClaims(t *testing.T) {
// 	tokenEP, _ := url.Parse("https://test.com/token")
// 	introspectEP, _ := url.Parse("https://test.com/tokenIntrospection")

// 	// Provider.
// 	pv := &provider{
// 		Issuer:       "https://test.com/",
// 		TokenEP:      tokenEP,
// 		IntrospectEP: introspectEP,
// 	}
// 	// Client.
// 	cl := &client{
// 		ID:     "testClientID",
// 		Secret: "testClientSecret",
// 	}

// 	// JWT handler.
// 	jh := newJWTHandler()

// 	// AT.
// 	at, _ := jh.TokenWithClaims(
// 		&atClaims{
// 			RegisteredClaims: &jwt.RegisteredClaims{
// 				Issuer:    "https://test.com/",
// 				Audience:  jwt.ClaimStrings{"resourceServer"},
// 				ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
// 			},
// 		})
// 	atStr, _ := jh.SignedString(at)
// 	atmc := &jwt.MapClaims{}
// 	jh.ParseWithClaims(atStr, atmc, []jwt.ParserOption{}...)

// 	// Invalid AT.
// 	invalidAT, _ := jh.TokenWithClaims(
// 		&atClaims{
// 			RegisteredClaims: &jwt.RegisteredClaims{
// 				Issuer:    "https://test.com/error",
// 				Audience:  jwt.ClaimStrings{"resourceServer"},
// 				ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
// 			},
// 		})
// 	invalidATStr, _ := jh.SignedString(invalidAT)

// 	// IDT.
// 	idt, _ := jh.TokenWithClaims(
// 		&idtClaims{
// 			RegisteredClaims: &jwt.RegisteredClaims{
// 				Issuer:    "https://test.com/",
// 				Audience:  jwt.ClaimStrings{"testClientID", "testSubClientID"},
// 				ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
// 			},
// 			Azp:      "testClientID",
// 			Acr:      "1",
// 			AuthTime: int(time.Now().Unix() - 10),
// 			Nonce:    "testNonce",
// 		})
// 	idtStr, _ := jh.SignedString(idt)
// 	idtmc := &jwt.MapClaims{}
// 	jh.ParseWithClaims(idtStr, idtmc, []jwt.ParserOption{}...)

// 	// IDT with invalid iss.
// 	invalidIssIDT, _ := jh.TokenWithClaims(
// 		&idtClaims{
// 			RegisteredClaims: &jwt.RegisteredClaims{
// 				Issuer:    "https://test.com/error",
// 				Audience:  jwt.ClaimStrings{"testClientID"},
// 				ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
// 			},
// 			Azp:      "testClientID",
// 			Acr:      "1",
// 			AuthTime: 0,
// 		})
// 	invalidIssIDTStr, _ := jh.SignedString(invalidIssIDT)

// 	// IDT with invalid azp.
// 	invalidAzpIDT, _ := jh.TokenWithClaims(&idtClaims{
// 		RegisteredClaims: &jwt.RegisteredClaims{
// 			Issuer:    "https://test.com/",
// 			Audience:  jwt.ClaimStrings{"testClientID", "testSubClientID"},
// 			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
// 		},
// 		Azp:      "testSubClientID2",
// 		Acr:      "1",
// 		AuthTime: 0,
// 	})
// 	invalidAzpIDTStr, _ := jh.SignedString(invalidAzpIDT)

// 	// IDT without auth_time.
// 	authTimeNotExistsIDT, _ := jh.TokenWithClaims(
// 		&idtClaims{
// 			RegisteredClaims: &jwt.RegisteredClaims{
// 				Issuer:    "https://test.com/",
// 				Audience:  jwt.ClaimStrings{"testClientID", "testSubClientID"},
// 				ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
// 			},
// 			Azp: "testClientID",
// 			Acr: "1",
// 		})
// 	authTimeNotExistsIDTStr, _ := jh.SignedString(authTimeNotExistsIDT)

// 	// IDT with invalid auth_time.
// 	invalidAuthTimeIDT, _ := jh.TokenWithClaims(
// 		&invalidAuthTimeIDTClaims{
// 			RegisteredClaims: &jwt.RegisteredClaims{
// 				Issuer:    "https://test.com/",
// 				Audience:  jwt.ClaimStrings{"testClientID", "testSubClientID"},
// 				ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
// 			},
// 			Azp:      "testClientID",
// 			Acr:      "1",
// 			AuthTime: "invalidAuthTime",
// 		})
// 	invalidAuthTimeIDTStr, _ := jh.SignedString(invalidAuthTimeIDT)

// 	type condition struct {
// 		redeemer             tokenRedeemer
// 		introspector         tokenIntrospector
// 		introspectionEnabled bool
// 		tokens               *OAuthTokens
// 		extraParams          map[string]string
// 		nonce                string
// 	}

// 	type action struct {
// 		t          *OAuthTokens
// 		c          *OAuthClaims
// 		err        error
// 		errPattern *regexp.Regexp
// 	}

// 	condEmptyAT := "access token has empty string"
// 	condTokenIntrospection := "token introspection enabled"
// 	condTokenRefresh := "token refresh"
// 	condValidateIDT := "validate id token"
// 	condTokenIntrospectionFailed := "token introspection failed"
// 	condTokenRefreshFailed := "failed to token refresh"
// 	condValidateRefreshedATFailed := "failed to refreshed access token"

// 	condValidateIDTFailed := "failed to validate id token"
// 	condValidateAzp := "validate azp with id token"
// 	condInvalidAzp := "invalid azp of id token"
// 	condValidateAuthTime := "validate auth_time with id token"
// 	condInvalidMaxAge := "invalid max_age for request parameter"
// 	condAuthTimeNotExists := "auth_time dose not exists in id token"
// 	condInvalidAuthTime := "invalid auth_time in id token"
// 	condExceededAuthTime := "token authorization time exceeded"
// 	condInvalidNonce := "invalid nonce"

// 	actError := "error"
// 	actNoError := "no error"

// 	tb := testutil.NewTableBuilder[*condition, *action]()
// 	tb.Name(t.Name())
// 	tb.Condition(condEmptyAT, "access token has empty string")
// 	tb.Condition(condTokenIntrospection, "token introspection failed")
// 	tb.Condition(condTokenRefresh, "token refresh")
// 	tb.Condition(condValidateIDT, "validate id token")
// 	tb.Condition(condTokenIntrospectionFailed, "token introspection enabled")
// 	tb.Condition(condTokenRefreshFailed, "failed to token refresh")
// 	tb.Condition(condValidateRefreshedATFailed, "failed to refreshed access token")
// 	tb.Condition(condValidateIDTFailed, "failed to validate id token")
// 	tb.Condition(condInvalidAzp, "failed to validate azp of id token")
// 	tb.Condition(condValidateAzp, "validate azp with id token")
// 	tb.Condition(condValidateAuthTime, "validate auth_time with id token")
// 	tb.Condition(condInvalidMaxAge, "invalid max_age for request parameter")
// 	tb.Condition(condAuthTimeNotExists, "auth_time dose not exists in id token")
// 	tb.Condition(condInvalidAuthTime, "invalid auth_time in id token")
// 	tb.Condition(condExceededAuthTime, "token authorization time exceeded")
// 	tb.Condition(condInvalidNonce, "invalid nonce")
// 	tb.Action(actError, "check that the expected error is returned")
// 	tb.Action(actNoError, "check that the there is no error")

// 	table := tb.Build()

// 	gen := testutil.NewCase[*condition, *action]
// 	testCases := []*testutil.Case[*condition, *action]{
// 		gen(
// 			"access token is empty",
// 			[]string{condEmptyAT},
// 			[]string{actError},
// 			&condition{
// 				tokens: &OAuthTokens{
// 					AT: "",
// 				},
// 			},
// 			&action{
// 				t:          nil,
// 				c:          nil,
// 				err:        base.ErrAppAuthnInvalidToken,
// 				errPattern: regexp.MustCompile(core.ErrPrefix + `invalid access token`),
// 			},
// 		),
// 		gen(
// 			"validation by token introspection",
// 			[]string{condTokenIntrospection},
// 			[]string{actNoError},
// 			&condition{
// 				tokens: &OAuthTokens{
// 					AT: atStr,
// 				},
// 				introspector: &testTokenIntrospection{
// 					mc:     *atmc,
// 					active: true,
// 				},
// 				introspectionEnabled: true,
// 			},
// 			&action{
// 				t: &OAuthTokens{
// 					AT: atStr,
// 				},
// 				c: &OAuthClaims{
// 					Method: acMethod,
// 					AT:     *atmc,
// 				},
// 				err: nil,
// 			},
// 		),
// 		gen(
// 			"failed to do token introspection",
// 			[]string{condTokenIntrospection, condTokenIntrospectionFailed},
// 			[]string{actError},
// 			&condition{
// 				tokens: &OAuthTokens{
// 					AT: atStr,
// 				},
// 				introspector: &testTokenIntrospection{
// 					mc: nil,
// 				},
// 				introspectionEnabled: true,
// 			},
// 			&action{
// 				t:          nil,
// 				err:        base.ErrTokenIntrospection,
// 				errPattern: regexp.MustCompile(core.ErrPrefix + `failed to token introspection .* \[token introspection error\]`),
// 			},
// 		),
// 		gen(
// 			"validation by local validation",
// 			[]string{},
// 			[]string{actNoError},
// 			&condition{
// 				tokens: &OAuthTokens{
// 					AT: atStr,
// 				},
// 			},
// 			&action{
// 				t: &OAuthTokens{
// 					AT: atStr,
// 				},
// 				c: &OAuthClaims{
// 					Method: acMethod,
// 					AT:     *atmc,
// 				},
// 				err: nil,
// 			},
// 		),
// 		gen(
// 			"token refresh",
// 			[]string{condTokenRefreshFailed},
// 			[]string{actNoError},
// 			&condition{
// 				tokens: &OAuthTokens{
// 					AT: invalidATStr,
// 					RT: "RT",
// 				},
// 				redeemer: &testRedeemToken{
// 					&TokenResponseModel{
// 						AccessToken:  atStr,
// 						RefreshToken: "Response of RT",
// 					},
// 				},
// 			},
// 			&action{
// 				t: &OAuthTokens{
// 					AT: atStr,
// 					RT: "Response of RT",
// 				},
// 				c: &OAuthClaims{
// 					Method: acMethod,
// 					AT:     *atmc,
// 				},
// 				err: nil,
// 			},
// 		),
// 		gen(
// 			"token refresh failure",
// 			[]string{condTokenRefresh, condTokenRefreshFailed},
// 			[]string{actError},
// 			&condition{
// 				tokens: &OAuthTokens{
// 					AT: invalidATStr,
// 					RT: "RT",
// 				},
// 				redeemer: &testRedeemToken{nil},
// 			},
// 			&action{
// 				t:          nil,
// 				c:          nil,
// 				err:        base.ErrAppAuthnRefreshToken,
// 				errPattern: regexp.MustCompile(core.ErrPrefix + `failed to refresh token`),
// 			},
// 		),
// 		gen(
// 			"failed to validate refreshed access token",
// 			[]string{condTokenRefresh, condValidateRefreshedATFailed},
// 			[]string{actError},
// 			&condition{
// 				tokens: &OAuthTokens{
// 					AT: invalidATStr,
// 					RT: "RT",
// 				},
// 				redeemer: &testRedeemToken{
// 					&TokenResponseModel{
// 						AccessToken:  invalidATStr,
// 						RefreshToken: "RT of response",
// 					},
// 				},
// 			},
// 			&action{
// 				t:          nil,
// 				c:          nil,
// 				err:        base.ErrTokenValidation,
// 				errPattern: regexp.MustCompile(core.ErrPrefix + `failed to validate refreshed access token`),
// 			},
// 		),
// 		gen(
// 			"validate id token",
// 			[]string{condValidateIDT, condValidateAzp, condValidateAuthTime},
// 			[]string{actNoError},
// 			&condition{
// 				tokens: &OAuthTokens{
// 					AT:  atStr,
// 					IDT: idtStr,
// 				},
// 				extraParams: map[string]string{"max_age": "60"},
// 				nonce:       "testNonce",
// 			},
// 			&action{
// 				t: &OAuthTokens{
// 					AT:  atStr,
// 					IDT: idtStr,
// 				},
// 				c: &OAuthClaims{
// 					Method: acMethod,
// 					AT:     *atmc,
// 					IDT:    *idtmc,
// 				},
// 				err: nil,
// 			},
// 		),
// 		gen(
// 			"failed to validate id token",
// 			[]string{condValidateIDT, condValidateIDTFailed},
// 			[]string{actError},
// 			&condition{
// 				tokens: &OAuthTokens{
// 					AT:  atStr,
// 					IDT: invalidIssIDTStr,
// 				},
// 			},
// 			&action{
// 				t:          nil,
// 				c:          nil,
// 				err:        base.ErrTokenValidation,
// 				errPattern: regexp.MustCompile(core.ErrPrefix + `failed to validate id token`),
// 			},
// 		),
// 		gen(
// 			"failed to validate azp of idt",
// 			[]string{condValidateIDT, condValidateAzp, condInvalidAzp},
// 			[]string{actError},
// 			&condition{
// 				tokens: &OAuthTokens{
// 					AT:  atStr,
// 					IDT: invalidAzpIDTStr,
// 				},
// 			},
// 			&action{
// 				t:          nil,
// 				c:          nil,
// 				err:        base.ErrInvalidAzp,
// 				errPattern: regexp.MustCompile(core.ErrPrefix + `failed to validate id token. token has invalid claims: token has invalid authorized party`),
// 			},
// 		),
// 		gen(
// 			"invalid max_age for request parameter",
// 			[]string{condValidateIDT, condValidateAuthTime, condInvalidMaxAge},
// 			[]string{actError},
// 			&condition{
// 				tokens: &OAuthTokens{
// 					AT:  atStr,
// 					IDT: idtStr,
// 				},
// 				extraParams: map[string]string{"max_age": "invalidMaxAge"},
// 			},
// 			&action{
// 				t:          nil,
// 				c:          nil,
// 				err:        base.ErrInvalidMaxAge,
// 				errPattern: regexp.MustCompile(core.ErrPrefix + `invalid max_age for request parameter`),
// 			},
// 		),
// 		gen(
// 			"auth_time dose not exists in id token",
// 			[]string{condValidateIDT, condValidateAuthTime, condAuthTimeNotExists},
// 			[]string{actError},
// 			&condition{
// 				tokens: &OAuthTokens{
// 					AT:  atStr,
// 					IDT: authTimeNotExistsIDTStr,
// 				},
// 				extraParams: map[string]string{"max_age": "60"},
// 			},
// 			&action{
// 				t:          nil,
// 				c:          nil,
// 				err:        base.ErrAuthTimeNotExists,
// 				errPattern: regexp.MustCompile(core.ErrPrefix + `auth_time dose not exists in id token`),
// 			},
// 		),
// 		gen(
// 			"invalid auth_time in id token",
// 			[]string{condValidateIDT, condValidateAuthTime, condInvalidAuthTime},
// 			[]string{actError},
// 			&condition{
// 				tokens: &OAuthTokens{
// 					AT:  atStr,
// 					IDT: invalidAuthTimeIDTStr,
// 				},
// 				extraParams: map[string]string{"max_age": "60"},
// 			},
// 			&action{
// 				t:          nil,
// 				c:          nil,
// 				err:        base.ErrInvalidAuthTime,
// 				errPattern: regexp.MustCompile(`.*`),
// 			},
// 		),
// 		gen(
// 			"token authorization time exceeded",
// 			[]string{condValidateIDT, condValidateAuthTime, condExceededAuthTime},
// 			[]string{actError},
// 			&condition{
// 				tokens: &OAuthTokens{
// 					AT:  atStr,
// 					IDT: idtStr,
// 				},
// 				extraParams: map[string]string{"max_age": "0"},
// 			},
// 			&action{
// 				t:          nil,
// 				c:          nil,
// 				err:        base.ErrExceededAuthTime,
// 				errPattern: regexp.MustCompile(core.ErrPrefix + `failed to validate id token. token has invalid claims: token authorization time exceeded`),
// 			},
// 		),
// 		gen(
// 			"invalid nonce",
// 			[]string{condValidateIDT, condValidateAuthTime, condInvalidNonce},
// 			[]string{actError},
// 			&condition{
// 				tokens: &OAuthTokens{
// 					AT:  atStr,
// 					IDT: idtStr,
// 				},
// 				nonce: "invalidNonce",
// 			},
// 			&action{
// 				t:          nil,
// 				c:          nil,
// 				err:        base.ErrInvalidNonce,
// 				errPattern: regexp.MustCompile(core.ErrPrefix + `failed to validate id token. token has invalid claims: token has invalid nonce`),
// 			},
// 		),
// 	}

// 	testutil.Register(table, testCases...)
//

// 	for _, tt := range table.Entries() {
// 		tt := tt
// 		t.Run(tt.Name(), func(t *testing.T) {
// 			h := &authorizationCodeHandler{
// 				tokenRedeemer:             tt.C().redeemer,
// 				tokenIntrospector:         tt.C().introspector,
// 				provider:                  pv,
// 				client:                    cl,
// 				tokenIntrospectionEnabled: tt.C().introspectionEnabled,
// 				resourceServer:            "resourceServer",
// 				jh:                        jh,
// 				extraParams:               tt.C().extraParams,
// 			}

// 			r := httptest.NewRequest(http.MethodPost, "https://test.com/token", nil)

// 			tokens, claims, err := h.validClaims(r, tt.C().tokens, tt.C().nonce)

// 			// Check for error.
// 			testutil.DiffError(t, tt.A().err, tt.A().errPattern, err)

// 			// Check tokens.
// 			testutil.Diff(t, tt.A().t, tokens)

// 			// Check claims.
// 			opts := []cmp.Option{
// 				cmpopts.IgnoreFields(OAuthClaims{}, "AuthTime"),
// 			}
// 			testutil.Diff(t, tt.A().c, claims, opts...)

// 		})
// 	}
// }

// func TestAuthorizationCodeRefreshToken(t *testing.T) {

// 	type condition struct {
// 		tr tokenRedeemer
// 		t  *OAuthTokens
// 	}

// 	type action struct {
// 		t          *OAuthTokens
// 		err        error
// 		errPattern *regexp.Regexp
// 	}

// 	condEmptyRT := "refresh token has empty string"
// 	condRedeemTokenFailed := "redeem token failed"
// 	condEmptyResIDT := "IDT in response has empty string"
// 	condEmptyResRT := "RT in response has empty string"
// 	actError := "error"
// 	actNoError := "no error"

// 	tb := testutil.NewTableBuilder[*condition, *action]()
// 	tb.Name(t.Name())
// 	tb.Condition(condEmptyRT, "refresh token is empty")
// 	tb.Condition(condRedeemTokenFailed, "redeem token failed")
// 	tb.Condition(condEmptyResIDT, "IDT in response is empty")
// 	tb.Condition(condEmptyResRT, "RT in response has empty string")
// 	tb.Action(actError, "check that the expected error is returned")
// 	tb.Action(actNoError, "check that the there is no error")

// 	table := tb.Build()

// 	gen := testutil.NewCase[*condition, *action]
// 	testCases := []*testutil.Case[*condition, *action]{
// 		gen(
// 			"no error",
// 			[]string{},
// 			[]string{actNoError},
// 			&condition{
// 				tr: &testRedeemToken{
// 					&TokenResponseModel{
// 						AccessToken:  "AT of response",
// 						RefreshToken: "RT of response",
// 						IDToken:      "IDT of response",
// 					},
// 				},
// 				t: &OAuthTokens{
// 					AT:  "AT",
// 					RT:  "RT",
// 					IDT: "IDT",
// 				},
// 			},
// 			&action{
// 				t: &OAuthTokens{
// 					AT:  "AT of response",
// 					RT:  "RT of response",
// 					IDT: "IDT of response",
// 				},
// 				err: nil,
// 			},
// 		),
// 		gen(
// 			"IDT and RT are not returned in response",
// 			[]string{condEmptyResRT, condEmptyResIDT},
// 			[]string{actNoError},
// 			&condition{
// 				tr: &testRedeemToken{
// 					&TokenResponseModel{
// 						AccessToken: "AT of response",
// 					},
// 				},
// 				t: &OAuthTokens{
// 					AT:  "AT",
// 					RT:  "RT",
// 					IDT: "IDT",
// 				},
// 			},
// 			&action{
// 				t: &OAuthTokens{
// 					AT:  "AT of response",
// 					RT:  "RT",
// 					IDT: "IDT",
// 				},
// 				err: nil,
// 			},
// 		),
// 		gen(
// 			"RT is empty",
// 			[]string{condEmptyRT},
// 			[]string{actError},
// 			&condition{
// 				tr: &testRedeemToken{
// 					&TokenResponseModel{
// 						AccessToken:  "AT of response",
// 						RefreshToken: "RT of response",
// 						IDToken:      "IDT of response",
// 					},
// 				},
// 				t: &OAuthTokens{
// 					AT:  "AT",
// 					RT:  "",
// 					IDT: "IDT",
// 				},
// 			},
// 			&action{
// 				t:          nil,
// 				err:        base.ErrAppAuthnInvalidToken,
// 				errPattern: regexp.MustCompile(core.ErrPrefix + `invalid access token`),
// 			},
// 		),
// 		gen(
// 			"redeem token failure",
// 			[]string{condRedeemTokenFailed},
// 			[]string{actError},
// 			&condition{
// 				tr: &testRedeemToken{nil},
// 				t: &OAuthTokens{
// 					AT:  "AT",
// 					RT:  "RT",
// 					IDT: "IDT",
// 				},
// 			},
// 			&action{
// 				t:          nil,
// 				err:        base.ErrAppAuthnRedeemToken,
// 				errPattern: regexp.MustCompile(core.ErrPrefix + `failed to redeem token`),
// 			},
// 		),
// 	}

// 	testutil.Register(table, testCases...)
//

// 	for _, tt := range table.Entries() {
// 		tt := tt
// 		t.Run(tt.Name(), func(t *testing.T) {
// 			r := httptest.NewRequest(http.MethodPost, "https://test.com/token", nil)

// 			h := &authorizationCodeHandler{
// 				tokenRedeemer: tt.C().tr,
// 			}

// 			tokens, err := h.refreshToken(r, tt.C().t)

// 			// Check for error.
// 			testutil.DiffError(t, tt.A().err, tt.A().errPattern, err)

// 			// Check token introspection claims.
// 			testutil.Diff(t, tt.A().t, tokens)

// 		})
// 	}

// }
