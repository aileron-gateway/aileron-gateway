package oauth

import (
	"cmp"
	"errors"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	v1 "github.com/aileron-gateway/aileron-gateway/apis/app/v1"
	"github.com/aileron-gateway/aileron-gateway/app"
	"github.com/aileron-gateway/aileron-gateway/core"
	utilhttp "github.com/aileron-gateway/aileron-gateway/util/http"
	"github.com/aileron-gateway/aileron-gateway/util/security"
	"github.com/aileron-gateway/aileron-gateway/util/session"
)

const (
	// authzCodeSessionKey is the authentication method identifier string.
	// This string is embedded in claims to make it possible
	// to identify the authentication method in the later middleware or handlers.
	authzCodeSessionKey = "__authzCode"

	// authzCodeRedirectSessionKey is the key string to save
	// redirect destination.
	authzCodeRedirectSessionKey = "_authzCodeRed"

	// authzCodeCallbackSessionKey is the key string to save
	// callback path.
	authzCodeCallbackSessionKey = "_authzCodeCall"
)

func newAuthorizationCodeHandler(bh *baseHandler, spec *v1.AuthorizationCodeHandler) (*authorizationCodeHandler, error) {
	csrf := &csrfStateGenerator{
		stateDisabled: spec.DisableState,
		nonceDisabled: spec.DisableNonce,
		pkceDisabled:  spec.DisablePKCE,
		method:        codeChallengeMethod[spec.PKCEMethod],
	}

	u, err := url.Parse(spec.CallbackURL)
	if err != nil {
		return nil, core.ErrCoreGenCreateObject.WithStack(err, map[string]any{"kind": kind})
	}

	var ro *requestObjectGenerator
	if spec.RequestObject != nil {
		requestObjectJH, err := security.NewJWTHandler(spec.RequestObject.JWTHandler, nil)
		if err != nil {
			return nil, core.ErrCoreGenCreateObject.WithStack(err, map[string]any{"kind": kind})
		}

		requestURIPath, err := url.Parse(spec.RequestObject.RequestURI)
		if err != nil {
			return nil, core.ErrCoreGenCreateObject.WithStack(err, map[string]any{"kind": kind})
		}

		ro = &requestObjectGenerator{
			requestURI:     spec.RequestObject.RequestURI,
			requestURIPath: requestURIPath.Path,
			jh:             requestObjectJH,
			exp:            spec.RequestObject.Exp,
			nbf:            spec.RequestObject.Nbf,
			cacheDisabled:  spec.RequestObject.DisableCache,
		}
	}

	var jarm *jarmValidator
	if spec.JARM != nil {
		jarmJH, err := security.NewJWTHandler(spec.JARM.JWTHandler, nil)
		if err != nil {
			return nil, core.ErrCoreGenCreateObject.WithStack(err, map[string]any{"kind": kind})
		}
		jarm = &jarmValidator{
			responseMode: spec.JARM.ResponseMode,
			jh:           jarmJH,
		}
	}

	return &authorizationCodeHandler{
		baseHandler:         bh,
		csrf:                csrf,
		loginPath:           spec.LoginPath,
		callbackPath:        u.Path,
		callbackURL:         spec.CallbackURL,
		redirectPath:        spec.RedirectPath,
		redirectPathPattern: regexp.MustCompile(cmp.Or(spec.RedirectPathPattern, "^$")),
		redirectToLogin:     spec.RedirectToLogin,
		unauthorizeAny:      spec.UnauthorizeAny,
		restoreRequest:      spec.RestoreRequest,
		urlParams:           url.PathEscape(strings.Join(spec.URLParams, "&")),
		fapiEnabled:         spec.EnabledFAPI,
		requestObject:       ro,
		jarm:                jarm,
	}, nil
}

type authorizationCodeHandler struct {
	*baseHandler

	csrf *csrfStateGenerator

	loginPath    string
	callbackPath string
	callbackURL  string

	redirectPath        string
	redirectPathPattern *regexp.Regexp

	redirectToLogin bool
	unauthorizeAny  bool
	restoreRequest  bool

	// urlParams is the additional url query parameters
	// defined by the users.
	// This parameter will be set for authentication request.
	urlParams string

	maxAge int64

	fapiEnabled bool

	requestObject *requestObjectGenerator

	jarm *jarmValidator
}

// ServeAuthn handle authentication request.
// return newRequest, authenticated, shouldReturn, error.
func (h *authorizationCodeHandler) ServeAuthn(w http.ResponseWriter, r *http.Request) (*http.Request, app.AuthResult, bool, error) {
	// authorizationCodeHandler requires session object.
	// If not found, yield authentication to other authentication handlers.
	ss := session.SessionFromContext(r.Context())
	if ss == nil {
		return nil, app.AuthContinue, false, nil
	}

	oauthTokens := &OAuthTokens{}
	if err := ss.Extract(authzCodeSessionKey, oauthTokens); err == nil { // OAuthTokens was found in the session.
		oc := h.oauthContextFromName(oauthTokens.Context)
		err := oc.validOauthClaims(r.Context(), oauthTokens, maxAgeValidation(h.maxAge))
		if err != nil {
			ss.Delete(ropcSessionKey)
			if err != reAuthenticationRequired {
				return nil, app.AuthContinue, true, err
			}
		} else {
			// Valid OAuth tokens are found in the session.
			session.MustPersist(ss, authzCodeSessionKey, oauthTokens)
			r = r.WithContext(oc.contextWithToken(r.Context(), oauthTokens))
			if r.URL.Path == h.loginPath || r.URL.Path == h.callbackPath {
				w.Header().Set("Content-Type", "text/plain; charset=utf-8")
				w.Header().Set("X-Content-Type-Options", "nosniff")
				w.Header().Set("Location", h.redirectPath)
				w.WriteHeader(http.StatusFound)
				w.Write(nil)
				return r, app.AuthContinue, true, nil // Redirect to redirectPath.
			}
			return r, app.AuthSucceeded, false, nil // Authentication succeeded !!
		}
	}

	// -----------------------------------------------------------
	// OAuthTokens was not found in the session.
	// From here, process users who are not authenticated
	// or authentication expired.
	// -----------------------------------------------------------

	// Check if there is an oauthContext that can handle this request.
	// If not found, yield authentication to other authentication handlers.
	oc, err := h.oauthContext(r)
	if err != nil {
		return nil, app.AuthContinue, false, nil
	}

	if h.requestObject != nil && r.URL.Path == h.requestObject.requestURIPath {
		if err := h.handleRequestObjectURI(w, oc); err != nil {
			return nil, false, false, utilhttp.NewHTTPError(err, http.StatusUnauthorized)
		}
		return nil, false, true, nil
	}

	if r.URL.Path != h.callbackPath {
		err := h.handleLogin(w, r, oc)
		return nil, app.AuthContinue, true, err
	}

	r, _, err = h.handleCallback(r, oc)
	if err != nil {
		return nil, app.AuthFailed, true, err
	}

	if h.restoreRequest {
		r, err = session.ExtractRequest(ss, r)
		if err != nil {
			err = app.ErrAppGenSessionOperation.WithStack(err, map[string]any{"operation": "extract original request from session."})
			return nil, app.AuthFailed, true, utilhttp.NewHTTPError(err, http.StatusInternalServerError)
		}
		return r, app.AuthSucceeded, false, nil // Authentication succeeded !!
	} else {
		rd := []byte{}
		if err := ss.Extract(authzCodeRedirectSessionKey, &rd); err != nil {
			return r, app.AuthFailed, true, utilhttp.NewHTTPError(err, http.StatusForbidden)
		}
		ss.Delete(authzCodeRedirectSessionKey)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("Location", string(rd))
		w.WriteHeader(http.StatusFound)
		w.Write(nil)
		return nil, app.AuthSucceeded, true, nil // Authentication succeeded !!
	}
}

// authenticationRequest executes an authentication request.
func (h *authorizationCodeHandler) handleLogin(w http.ResponseWriter, r *http.Request, oc *oauthContext) core.HTTPError {
	ss := session.SessionFromContext(r.Context())
	if ss == nil {
		err := app.ErrAppGenSessionOperation.WithStack(nil, map[string]any{"operation": "get session from context.", "reason": "session is nil"})
		return utilhttp.NewHTTPError(err, http.StatusInternalServerError)
	}

	// Store the redirect URI value by FAPI Part1 - 5.2.3. Public client.
	// 4. shall store the redirect URI value in the resource owner's user-agents (such as browser) session and compare it
	// with the redirect URI that the authorization response was received at, where, if the URIs do not match,
	// the client shall terminate the process with error;
	// https://openid.net/specs/openid-financial-api-part-1-1_0-final.html
	if h.fapiEnabled {
		if err := ss.Persist(authzCodeCallbackSessionKey, &h.callbackPath); err != nil {
			err := app.ErrAppGenSessionOperation.WithStack(err, map[string]any{"operation": "persist callback path in the session."})
			return utilhttp.NewHTTPError(err, http.StatusInternalServerError)
		}
	}

	if r.URL.Path != h.loginPath {
		if h.redirectToLogin {
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			w.Header().Set("X-Content-Type-Options", "nosniff")
			w.Header().Set("Location", h.loginPath)
			w.WriteHeader(http.StatusFound)
			w.Write(nil)
			return nil // Redirect to login path.
		}
		if h.unauthorizeAny {
			return utilhttp.ErrUnauthorized.(*utilhttp.HTTPError)
		}
		if h.restoreRequest {
			if err := session.PersistRequest(ss, r); err != nil {
				err := app.ErrAppGenSessionOperation.WithStack(err, map[string]any{"operation": "persist original request in session."})
				return utilhttp.NewHTTPError(err, http.StatusInternalServerError)
			}
		}
	}

	rd := []byte(cmp.Or(r.URL.Query().Get("rd"), h.redirectPath, r.URL.RequestURI(), "/"))
	if !h.redirectPathPattern.Match(rd) {
		return utilhttp.ErrUnauthorized.(*utilhttp.HTTPError)
	}
	if err := ss.Persist(authzCodeRedirectSessionKey, rd); err != nil {
		err := app.ErrAppGenSessionOperation.WithStack(err, map[string]any{"operation": "persist redirect path in session."})
		return utilhttp.NewHTTPError(err, http.StatusInternalServerError)
	}

	baseParams := url.Values{
		"response_type": []string{"code"},
		"client_id":     []string{oc.client.id},
	}

	optParams := url.Values{
		"redirect_uri": []string{h.callbackURL},
	}

	s, err := h.csrf.new()
	if err != nil {
		err := app.ErrAppAuthnGenerateCSRFParams.WithStack(err, nil)
		return utilhttp.NewHTTPError(err, http.StatusInternalServerError)
	}

	csrfParams := url.Values{}
	s.set(csrfParams)
	if err := ss.Persist(csrfSessionKey, s); err != nil {
		err := app.ErrAppGenSessionOperation.WithStack(err, map[string]any{"operation": "persist csrf params in the session."})
		return utilhttp.NewHTTPError(err, http.StatusInternalServerError)
	}

	// Use response_mode parameter as specified in FAPI Part 2 - 5.2.3.2 JARM
	// In addition, if the response_type value code is used in conjunction with the response_mode value jwt, the client
	// 1. shall verify the authorization responses as specified in JARM, Section 4.4.
	// https://openid.net/specs/openid-financial-api-part-2-1_0-final.html#jarm-1
	if h.jarm != nil {
		optParams.Add("response_mode", responseModeMethods[h.jarm.responseMode])
	}

	urlParamsValues, err := url.ParseQuery(h.urlParams)
	if err != nil {
		return utilhttp.NewHTTPError(err, http.StatusInternalServerError)
	}
	for k, v := range urlParamsValues {
		optParams[k] = v
	}

	// Use Request Object as specified in FAPI Part 2 - 5.2.3 Confidential client
	// 2. shall include the request or request_uri parameter as defined in Section 6
	// of OIDC in the authentication request;
	// https://openid.net/specs/openid-financial-api-part-2-1_0-final.html
	if h.requestObject != nil {
		ro, err := h.requestObject.new(oc, baseParams, optParams)
		if err != nil {
			err := app.ErrAppAuthnGenerateRequestObject.WithStack(err, nil)
			return utilhttp.NewHTTPError(err, http.StatusInternalServerError)
		}

		if h.requestObject.requestURI == "" {
			optParams = url.Values{"request": []string{ro}}
		} else {
			optParams = url.Values{"request_uri": []string{h.requestObject.getRequestURI(ro)}}
		}
	}

	// Authentication request by redirecting.
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("Location", oc.provider.authorizationEP+"?"+baseParams.Encode()+"&"+csrfParams.Encode()+"&scope="+url.PathEscape(oc.client.scope)+"&"+optParams.Encode())
	w.WriteHeader(http.StatusFound)
	w.Write(nil)

	return nil
}

// handleCallback handles callback response from authorization server.
func (h *authorizationCodeHandler) handleCallback(r *http.Request, oc *oauthContext) (*http.Request, *OAuthTokens, core.HTTPError) {
	ss := session.SessionFromContext(r.Context())
	if ss == nil {
		err := app.ErrAppGenSessionOperation.WithStack(nil, map[string]any{"operation": "session not found in context."})
		return r, nil, utilhttp.NewHTTPError(err, http.StatusUnauthorized)
	}

	// Compare the redirect URI value in the session with the redirect URI as specified in FAPI Part 1 - 5.2.3 for Public Client.
	// 4. shall store the redirect URI value in the resource owner's user-agents (such as browser) session and compare it
	// with the redirect URI that the authorization response was received at, where, if the URIs do not match,
	// the client shall terminate the process with error;
	// https://openid.net/specs/openid-financial-api-part-1-1_0-final.html
	if h.fapiEnabled {
		cp := ""
		if err := ss.Extract(authzCodeCallbackSessionKey, &cp); err != nil {
			return r, nil, utilhttp.NewHTTPError(err, http.StatusUnauthorized)
		}
		if cp != h.callbackPath {
			return r, nil, utilhttp.NewHTTPError(errors.New("redirect URI mismatch"), http.StatusUnauthorized)
		}
	}

	state := r.URL.Query().Get("state")
	code := r.URL.Query().Get("code")
	if h.jarm != nil {
		var err error
		code, state, err = h.jarm.valid(oc, r, h.csrf.stateDisabled)
		if err != nil {
			return r, nil, utilhttp.NewHTTPError(err, http.StatusUnauthorized)
		}
	}

	// Extract CSRF parameters and remove them from the session.
	s := &csrfStates{}
	if err := ss.Extract(csrfSessionKey, s); err != nil {
		err := app.ErrAppGenSessionOperation.WithStack(err, map[string]any{"operation": "extract csrf params from session."})
		return r, nil, utilhttp.NewHTTPError(err, http.StatusUnauthorized)
	}
	ss.Delete(csrfSessionKey)

	// =================================================
	// TODO: support alternative response modes.
	// https://openid.net/specs/oauth-v2-multiple-response-types-1_0.html
	// https://openid.net/specs/oauth-v2-form-post-response-mode-1_0.html
	// https://openid.net/specs/openid-financial-api-jarm.html

	if !h.csrf.stateDisabled && s.State != state {
		err := app.ErrAppAuthnInvalidParameters.WithStack(nil, map[string]any{"name": "CSRF", "reason": "state mismatch. " + s.State + " != " + r.URL.Query().Get("state")})
		return r, nil, utilhttp.NewHTTPError(err, http.StatusUnauthorized)
	}

	if code == "" {
		err := app.ErrAppAuthnInvalidParameters.WithStack(nil, map[string]any{"name": "callback", "reason": "authorization code not found in query"})
		return r, nil, utilhttp.NewHTTPError(err, http.StatusUnauthorized)
	}

	params := map[string]string{
		"grant_type":   "authorization_code",
		"code":         code,
		"redirect_uri": h.callbackURL,
	}

	if s.Verifier != "" {
		params["code_verifier"] = s.Verifier
	}
	// =================================================

	_, resp, err := oc.redeemToken(r.Context(), params)
	if err != nil {
		return r, nil, err
	}

	tokens := &OAuthTokens{
		Context: oc.name, // Context name handling this authentication flow.
		IDT:     resp.IDToken,
		AT:      resp.AccessToken,
		RT:      resp.RefreshToken,
	}

	err = oc.validOauthClaims(r.Context(), tokens, maxAgeValidation(h.maxAge), nonceValidation(s.Nonce))
	if err != nil {
		return r, tokens, err
	}

	session.MustPersist(ss, authzCodeSessionKey, tokens)

	// Save tokens in the context so the succeeding middleware
	// such as authorization middleware can use them.
	ctx := oc.contextWithToken(r.Context(), tokens)
	r = r.WithContext(ctx)

	return r, tokens, nil
}

func (h *authorizationCodeHandler) handleRequestObjectURI(w http.ResponseWriter, oc *oauthContext) error {
	baseParams := url.Values{
		"response_type": []string{"code"},
		"client_id":     []string{oc.client.id},
	}

	optParams := url.Values{
		"redirect_uri": []string{h.callbackURL},
	}

	// Use response_mode parameter as specified in FAPI Part 2 - 5.2.3.2 JARM
	// In addition, if the response_type value code is used in conjunction with the response_mode value jwt, the client
	// 1. shall verify the authorization responses as specified in JARM, Section 4.4.
	// https://openid.net/specs/openid-financial-api-part-2-1_0-final.html#jarm-1
	if h.jarm != nil {
		optParams.Add("response_mode", responseModeMethods[h.jarm.responseMode])
	}

	urlParamsValues, err := url.ParseQuery(h.urlParams)
	if err != nil {
		return err
	}
	for k, v := range urlParamsValues {
		optParams[k] = v
	}

	ro, err := h.requestObject.new(oc, baseParams, optParams)
	if err != nil {
		err := app.ErrAppAuthnGenerateRequestObject.WithStack(err, nil)
		return err
	}

	w.Header().Set("Content-Type", "application/jwt; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(ro))

	return nil
}
