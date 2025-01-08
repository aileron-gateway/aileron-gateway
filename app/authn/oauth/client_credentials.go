package oauth

import (
	"maps"
	"net/http"
	"sync"

	v1 "github.com/aileron-gateway/aileron-gateway/apis/app/v1"
	"github.com/aileron-gateway/aileron-gateway/app"
	utilhttp "github.com/aileron-gateway/aileron-gateway/util/http"
)

// newClientCredentialsHandler returns a new instance of ClientCredentials handler.
func newClientCredentialsHandler(bh *baseHandler, _ *v1.ClientCredentialsHandler) *clientCredentialsHandler {
	p := map[string]string{
		"grant_type": "client_credentials",
	}
	tokens := make(map[string]*OAuthTokens, len(bh.oauthCtxs))
	for name := range bh.oauthCtxs {
		tokens[name] = &OAuthTokens{}
	}

	return &clientCredentialsHandler{
		baseHandler: bh,
		queryParams: p,
		oauthTokens: tokens,
	}
}

type clientCredentialsHandler struct {
	*baseHandler

	// queryParams is the set of query parameters
	// that will be sent to authorization server
	// in token requests.
	queryParams map[string]string

	// mu protects oauthTokens.
	mu sync.RWMutex
	// oauthTokens is the cached token for this client.
	oauthTokens map[string]*OAuthTokens
}

// ServeAuthn handle authentication request.
// return newRequest, authenticated, shouldReturn, error.
func (h *clientCredentialsHandler) ServeAuthn(w http.ResponseWriter, r *http.Request) (*http.Request, app.AuthResult, bool, error) {
	oc, err := h.oauthContext(r)
	if err != nil {
		return nil, app.AuthContinue, false, utilhttp.NewHTTPError(err, http.StatusUnauthorized)
	}

	// Deep copy for goroutine safety.
	h.mu.RLock()
	oldTokens := h.oauthTokens[oc.name]
	tokens := &OAuthTokens{
		Context: oc.name, // Context name handling this authentication flow.
		AT:      oldTokens.AT,
		RT:      oldTokens.RT,
	}
	h.mu.RUnlock()

	// Get oauth claims with basic validation.
	err = oc.validOauthClaims(r.Context(), tokens)
	if err != nil {
		if err != reAuthenticationRequired {
			return nil, false, true, err
		}
	} else {
		r = r.WithContext(oc.contextWithToken(r.Context(), tokens))
		// Update cached tokens before return.
		h.mu.Lock()
		defer h.mu.Unlock()
		oldTokens.AT = tokens.AT
		oldTokens.RT = tokens.RT
		return r, app.AuthSucceeded, false, nil // Authentication succeeded !!
	}

	q := make(map[string]string, len(h.queryParams))
	maps.Copy(q, h.queryParams)
	q["scope"] = oc.client.scope
	_, resp, err := oc.redeemToken(r.Context(), q)
	if err != nil {
		return nil, false, true, err
	}
	tokens.AT = resp.AccessToken
	tokens.RT = resp.RefreshToken

	err = oc.validOauthClaims(r.Context(), tokens)
	if err != nil {
		return nil, app.AuthFailed, true, err
	}

	// Save tokens in the context so the succeeding middleware
	// such as authorization middleware can use them.
	ctx := oc.contextWithToken(r.Context(), tokens)
	r = r.WithContext(ctx)

	// Update cached tokens before return.
	h.mu.Lock()
	defer h.mu.Unlock()
	oldTokens.AT = tokens.AT
	oldTokens.RT = tokens.RT

	return r, true, false, nil // Authentication succeeded !!
}
