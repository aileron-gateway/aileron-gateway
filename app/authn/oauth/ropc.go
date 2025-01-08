package oauth

import (
	"bytes"
	"encoding/json"
	"io"
	"maps"
	"net/http"

	v1 "github.com/aileron-gateway/aileron-gateway/apis/app/v1"
	"github.com/aileron-gateway/aileron-gateway/app"
	"github.com/aileron-gateway/aileron-gateway/core"
	utilhttp "github.com/aileron-gateway/aileron-gateway/util/http"
	"github.com/aileron-gateway/aileron-gateway/util/session"
)

// ropcSessionKey is the authentication method identifier string.
// This string is embedded in claims to make it possible
// to identify the authentication method in the later middleware or handlers.
const ropcSessionKey = "_ropc"

// newROPCHandler returns a new instance of ROPC handler.
func newROPCHandler(bh *baseHandler, spec *v1.ROPCHandler) *ropcHandler {
	return &ropcHandler{
		baseHandler: bh,
		redeemPath:  spec.RedeemTokenPath,
		usernameKey: spec.UsernameKey,
		passwordKey: spec.PasswordKey,
		queryParams: map[string]string{"grant_type": "password"},
	}
}

type ropcHandler struct {
	*utilhttp.HandlerBase
	*baseHandler

	eh core.ErrorHandler

	// redeemPath is the path to get access tokens with ROPC flow.
	// This path will be checked when this handler is used as middleware.
	// Instead of using this field, uses can use this ropcHandler as a http.Handler.
	redeemPath              string
	saveSessionAtRedeemPath bool

	// usernameKey is the parameter key of username in POST form body.
	// Set BOTH usernameKey and passwordKey to parse username and password
	// from form request body.
	// Otherwise, ropcHandler tries to parse username and password from the
	// authorization header.
	usernameKey string
	// passwordKey is the parameter key of password in POST form body.
	// Set BOTH usernameKey and passwordKey to parse username and password
	// from form request body.
	// Otherwise, ropcHandler tries to parse username and password from the
	// authorization header.
	passwordKey string

	// queryParams is the set of query parameters
	// that are be sent to authorization server in token requests.
	queryParams map[string]string
}

// ServeHTTP expose redeem token endpoint as a HTTP handler.
// This is the implementation of http.Handler interface.
func (h *ropcHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	oc, err := h.oauthContext(r)
	if err != nil {
		h.eh.ServeHTTPError(w, r, utilhttp.NewHTTPError(err, http.StatusUnauthorized))
		return
	}
	resp, err := h.tokenRequest(r, oc)
	if err != nil {
		h.eh.ServeHTTPError(w, r, err)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(resp.StatusCode)
	w.Write(resp.RawBody)
}

// ServeAuthn handle authentication request.
// returns newRequest, authenticated, shouldReturn, error.
func (h *ropcHandler) ServeAuthn(w http.ResponseWriter, r *http.Request) (*http.Request, app.AuthResult, bool, error) {
	// ropcHandler requires session object.
	// If not found, yield authentication to other authentication handlers.
	ss := session.SessionFromContext(r.Context())
	if ss == nil {
		return nil, app.AuthContinue, false, nil
	}

	// Check if tokens exist in the session.
	tokens := &OAuthTokens{}
	if err := ss.Extract(ropcSessionKey, tokens); err == nil {
		oc := h.oauthContextFromName(tokens.Context)
		err := oc.validOauthClaims(r.Context(), tokens)
		if err != nil {
			ss.Delete(ropcSessionKey)
			if err != reAuthenticationRequired {
				return nil, app.AuthContinue, true, err
			}
		} else {
			// Valid OAuth tokens are found in the session.
			session.MustPersist(ss, ropcSessionKey, tokens)
			r = r.WithContext(oc.contextWithToken(r.Context(), tokens))
			if r.URL.Path == h.redeemPath {
				b, err := json.Marshal(tokens)
				if err != nil {
					return nil, app.AuthFailed, true, utilhttp.NewHTTPError(err, http.StatusInternalServerError)
				}
				w.Header().Set("Content-Type", "application/json; charset=utf-8")
				w.Header().Set("X-Content-Type-Options", "nosniff")
				w.WriteHeader(http.StatusOK)
				w.Write(b)
				return r, app.AuthContinue, true, nil // Return token info to the client.
			}
			return r, app.AuthSucceeded, false, nil // Authentication succeeded !!
		}
	}

	// -----------------------------------------------------------
	// From here, process users who are not authenticated
	// or authentication expired.
	// -----------------------------------------------------------

	// Check if there is an oauthContext that can handle this request.
	// If not found, yield authentication to other authentication handlers.
	oc, err := h.oauthContext(r)
	if err != nil {
		return nil, app.AuthContinue, false, nil
	}

	resp, err := h.tokenRequest(r, oc)
	if err != nil {
		return nil, app.AuthFailed, true, err // Return the err as it is.
	}

	tokens = &OAuthTokens{
		Context: oc.name,
		IDT:     resp.IDToken,
		AT:      resp.AccessToken,
		RT:      resp.RefreshToken,
	}

	if err = oc.validOauthClaims(r.Context(), tokens); err != nil {
		return nil, app.AuthFailed, true, err // Return the err as it is.
	}

	session.MustPersist(ss, ropcSessionKey, tokens)

	// Save tokens in the context so the succeeding middleware
	// such as authorization middleware can use them.
	ctx := oc.contextWithToken(r.Context(), tokens)
	r = r.WithContext(ctx)

	if r.URL.Path == h.redeemPath {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.WriteHeader(resp.StatusCode)
		w.Write(resp.RawBody)
		return r, app.AuthSucceeded, true, nil
	}

	return r, app.AuthSucceeded, false, nil
}

func (h *ropcHandler) tokenRequest(r *http.Request, oc *oauthContext) (*TokenResponse, core.HTTPError) {
	var un string // Provided username
	var pw string // Provided password

	// Parse username and password from the form body
	// when BOTH h.usernameKey and h.passwordKey are set.
	// Otherwise, try to parse them from authorization header.
	if h.usernameKey != "" && h.passwordKey != "" {
		b, err := io.ReadAll(r.Body)
		if err != nil {
			err := app.ErrAppGenReadHTTPBody.WithStack(err, map[string]any{"direction": "request", "body": string(b)})
			return nil, utilhttp.NewHTTPError(err, http.StatusInternalServerError)
		}
		r.Body = io.NopCloser(bytes.NewReader(b))

		un = r.PostFormValue(h.usernameKey)
		pw = r.PostFormValue(h.passwordKey)
	} else {
		var ok bool
		un, pw, ok = r.BasicAuth()
		if !ok {
			err := app.ErrAppAuthnInvalidCredential.WithStack(nil, map[string]any{"purpose": "ROPC", "reason": "invalid authorization header"})
			return nil, utilhttp.NewHTTPError(err, http.StatusUnauthorized)
		}
	}

	if un == "" || pw == "" {
		err := app.ErrAppAuthnInvalidCredential.WithStack(nil, map[string]any{"purpose": "ROPC", "reason": "username or password is empty"})
		return nil, utilhttp.NewHTTPError(err, http.StatusUnauthorized)
	}

	q := map[string]string{
		"username": un,
		"password": pw,
		"scope":    oc.client.scope,
	}
	maps.Copy(q, h.queryParams)

	_, resp, err := oc.redeemToken(r.Context(), q)

	return resp, err
}
