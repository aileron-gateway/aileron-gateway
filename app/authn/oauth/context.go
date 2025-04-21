// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package oauth

import (
	"cmp"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/aileron-gateway/aileron-gateway/app"
	"github.com/aileron-gateway/aileron-gateway/core"
	"github.com/aileron-gateway/aileron-gateway/kernel/log"
	utilhttp "github.com/aileron-gateway/aileron-gateway/util/http"
	"github.com/aileron-gateway/aileron-gateway/util/security"
	"github.com/golang-jwt/jwt/v5"
)

var (
	reAuthenticationRequired = utilhttp.NewHTTPError(errors.New("authn/oauth: authentication required"), http.StatusUnauthorized)
	errContextNotFound       = utilhttp.NewHTTPError(errors.New("authn/oauth: oauth context not found"), http.StatusUnauthorized)
)

type baseHandler struct {
	lg log.Logger
	eh core.ErrorHandler

	// oauthCtxs is the map of context name and
	// the corresponding context.
	// OAuth providers' information are contained
	// in the context and should provide the right operation
	// against the provider.
	// "default" is the special name of the context that will be
	// used by default when no context key were provided.
	oauthCtxs map[string]*oauthContext

	// contextQueryKey is the URL query key name
	// to get the name of oauthContext to use for authentication.
	// HTTP requests should contain the URL query like below.
	// Ex. "https://example.com/foo?<contextQueryKey>=provider1"
	contextQueryKey string

	// contextHeaderKey is the HTTP header key name
	// to get the name of oauthContext to use for authentication.
	// HTTP requests should contain the header like below.
	// Ex. "<contextHeaderKey>: provider1"
	contextHeaderKey string
}

// oauthContextFromName returns a oauth context from the given name.
// oauthContextFromName does not care if a context with the given name
// actually exists or not.
// oauthContextFromName returns nil when no context found.
func (h *baseHandler) oauthContextFromName(name string) *oauthContext {
	return h.oauthCtxs[name]
}

// oauthContext returns oauth context by the key.
// The name of the context can be specified by URL query or HTTP header.
// "default" is the special value that is used when no keys were found.
// oauthContext returns an error and nil context when no context found.
// Context names are looked up by the following order.
//   - URL Query: search context with query key if the key is configured.
//   - HTTP header: search context with header key if the key is configured.
//   - default: a special name "default".
func (h *baseHandler) oauthContext(r *http.Request) (*oauthContext, error) {
	var name string
	if h.contextQueryKey != "" {
		name = r.URL.Query().Get(h.contextQueryKey)
	} else if h.contextHeaderKey != "" {
		name = r.Header.Get(h.contextHeaderKey)
	}

	if name == "" {
		// Use "default" if no context obtained.
		name = "default"
	}

	if oc, ok := h.oauthCtxs[name]; ok {
		return oc, nil
	}

	return nil, app.ErrAppAuthnAuthentication.WithoutStack(errContextNotFound, nil)
}

// baseHandler is the base struct for authentication handlers.
// This struct is used by embedding to authentication handlers.
type oauthContext struct {
	tokenRedeemer
	tokenIntrospector
	userInfoRequestor

	lg log.Logger

	// name is the name of this context.
	name string

	// provider is the information of the OAuth provider.
	// Only use issuer value for JWT validation
	// in this oauthContext.
	provider *provider
	// client is the information of this OAuth client.
	client *client

	// atParseOpts is the list of access token parse options.
	atParseOpts []jwt.ParserOption
	// idtParseOpts is the list of ID token parse options.
	idtParseOpts []jwt.ParserOption
	// skipUnexpiredAT skips validation of unexpired access token.
	// This option works when the token was restored from session.
	skipUnexpiredAT bool
	// skipUnexpiredIDT skips validation of unexpired ID token.
	// This option works when the token was restored from session.
	skipUnexpiredIDT bool

	// jh is a JWT handler.
	// JWT handler should be configured when
	// JWT validation is needed during authentication flow.
	jh *security.JWTHandler

	// introspectionEnabled is the flag to enable token introspection.
	// JWT attributes are not validated when the introspection is enabled.
	// Introspection should be enabled when the token is not JWT.
	// If introspection enabled, local validation will be disabled.
	// It is not possible to use both introspection and local validation.
	introspectionEnabled bool

	// claimsKey is the key string to save claims in the session.
	// Claims are not saved in the context when this value is empty.
	claimsKey string
	// atProxyHeader is the header name for proxy an AccessToken.
	// Token will not be sent upstream when this value is empty.
	atProxyHeader string
	// idtProxyHeader is the header name for proxy an ID Token.
	// Token will not be sent upstream when this value is empty.
	idtProxyHeader string
}

// validOauthClaims returns a validated oauth claims.
// The method only do basic validation for the given tokens.
// Additional validation should be taken by callers if necessary.
func (c *oauthContext) validOauthClaims(ctx context.Context, tokens *OAuthTokens, opts ...validateOption) core.HTTPError {
	if (tokens == nil) || (tokens.IDT == "" && tokens.AT == "") {
		return reAuthenticationRequired
	}

	now := time.Now().Unix()

	// Validate the access token if exists.
	// Skip validation when the access token has not expired yet.
	if tokens.AT != "" && !(c.skipUnexpiredAT && (now < tokens.ATExp)) && !skipATValidation {
		claims, err := c.validateAT(ctx, tokens.AT)
		if err != nil && err != reAuthenticationRequired {
			return err
		}

		if !mapValue[bool](claims, "active") {
			if tokens.RT == "" {
				// Access token is invalid and no refresh token was given for refreshing the access token.
				// Caller should do authentication and obtain a new access token.
				return reAuthenticationRequired
			}

			// Try token refresh.
			params := map[string]string{
				"grant_type":    "refresh_token",
				"refresh_token": tokens.RT,
			}
			status, resp, err := c.redeemToken(ctx, params)
			if err != nil {
				if status >= 400 && status < 500 {
					return reAuthenticationRequired
				}
				return err
			}

			tokens.updated = true
			tokens.AT = resp.AccessToken
			tokens.RT = resp.RefreshToken
			tokens.IDT = cmp.Or(resp.IDToken, tokens.IDT) // Keep old ID token if not returned.

			claims, err = c.validateAT(ctx, tokens.AT)
			if err != nil {
				// Refresh token is exists but it might be expired.
				// Then require re-authentication.
				// Not to mention, this error can be a server-side error.
				return err
			}

			if !mapValue[bool](claims, "active") {
				// Token refresh has been succeeded.
				// But the token could be revoked for some reasons.
				// Then re-authentication should be required.
				return reAuthenticationRequired
			}
		}

		tokens.updated = true
		tokens.ATClaims = claims
		tokens.ATExp = int64(mapValue[float64](claims, "exp"))
		if c.lg.Enabled(log.LvDebug) { // Debug logging.
			attr := log.NewCustomAttrs("access_token", claims)
			c.lg.Debug(ctx, "access token validation succeeded", attr.Name(), attr.Map())
		}
	}

	// Validate the ID token if exists.
	// Skip validation when the ID token has not expired yet.
	if tokens.IDT != "" && !(c.skipUnexpiredIDT && (now < tokens.IDTExp)) && !skipIDTValidation {
		claims, err := c.validateIDT(tokens.IDT, opts)
		if err != nil {
			// Unauthorize for invalid ID token or require re-authentication if expired.
			return err
		}

		tokens.updated = true
		tokens.IDTClaims = claims
		tokens.IDTExp = int64(mapValue[float64](claims, "exp"))
		if c.lg.Enabled(log.LvDebug) { // Debug logging.
			attr := log.NewCustomAttrs("id_token", claims)
			c.lg.Debug(ctx, "ID token validation succeeded", attr.Name(), attr.Map())
		}
	}

	return nil
}

// validateAT validates access tokens and returns parsed map claims or HTTP error.
// Caller should respect the HTTP response returned from this method.
// We prepare validateOptions argument for future use.
func (c *oauthContext) validateAT(ctx context.Context, at string) (jwt.MapClaims, core.HTTPError) {
	if c.introspectionEnabled {
		_, claims, err := c.tokenIntrospection(ctx, map[string]string{"token": at})
		if err == nil {
			return claims, nil
		}

		if c.lg.Enabled(log.LvDebug) {
			err := app.ErrAppAuthnIntrospection.WithoutStack(err, map[string]any{"token": at})
			c.lg.Debug(ctx, "token introspection failed", err.Name(), err.Map())
		}
		return nil, err
	} else {
		claims, err := c.jh.ValidMapClaims(at, c.atParseOpts...)
		if err == nil {
			claims["active"] = true
			return claims, nil
		}

		sud, _ := claims.GetSubject()
		if sud == "" {
			err := app.ErrAppAuthnInvalidToken.WithoutStack(err, map[string]any{"name": "access token", "reason": "sub in access token does not exist.", "token": at})
			return nil, utilhttp.NewHTTPError(err, http.StatusUnauthorized)
		}

		if c.lg.Enabled(log.LvDebug) {
			err := app.ErrAppAuthnParseWithClaims.WithoutStack(err, map[string]any{"jwt": at})
			c.lg.Debug(ctx, "failed to parse claims", err.Name(), err.Map())
		}
		return nil, reAuthenticationRequired
	}
}

// validateIDT validates ID tokens and returns parsed map claims or HTTP error.
// Caller should respect the HTTP response returned from this method.
// Additional validation can be given by callers with validateOptions argument.
func (c *oauthContext) validateIDT(idt string, opts []validateOption) (jwt.MapClaims, core.HTTPError) {
	// Notes:
	// "aud" claim should be validated with audience.
	// client id should be validated with "azp" claim.
	claims, err := c.jh.ValidMapClaims(idt, c.idtParseOpts...)
	if err != nil {
		err := app.ErrAppAuthnInvalidToken.WithoutStack(err, map[string]any{"name": "id token", "reason": "token validation failed.", "token": idt})
		return nil, utilhttp.NewHTTPError(err, http.StatusUnauthorized)
	}

	sud, _ := claims.GetSubject()
	if sud == "" {
		err := app.ErrAppAuthnInvalidToken.WithoutStack(err, map[string]any{"name": "id token", "reason": "sub in ID token does not exist.", "token": idt})
		return nil, utilhttp.NewHTTPError(err, http.StatusUnauthorized)
	}

	iat, _ := claims.GetIssuedAt()
	if iat == nil {
		err := app.ErrAppAuthnInvalidToken.WithoutStack(err, map[string]any{"name": "id token", "reason": "iat in ID token does not exist.", "token": idt})
		return nil, utilhttp.NewHTTPError(err, http.StatusUnauthorized)
	}

	aud, _ := claims.GetAudience()
	azp, _ := (claims["azp"]).(string)
	if (len(aud) > 1 || azp != "") && azp != c.client.id {
		err := app.ErrAppAuthnInvalidToken.WithoutStack(err, map[string]any{"name": "id token", "reason": "azp in ID token does not match to client id.", "token": idt})
		return nil, utilhttp.NewHTTPError(err, http.StatusUnauthorized)
	}

	// Do optional validations.
	for _, opt := range opts {
		if err := opt.validate(claims); err != nil {
			return nil, err
		}
	}

	return claims, nil
}

// contextWithToken saves the oauth token in the context.
// contextWithToken set the ID Token and Access token in the proxy header if configured.
func (h *oauthContext) contextWithToken(ctx context.Context, tokens *OAuthTokens) context.Context {
	ph := utilhttp.ProxyHeaderFromContext(ctx)
	if ph == nil {
		ph = make(http.Header)
		ctx = utilhttp.ContextWithProxyHeader(ctx, ph)
	}
	if h.atProxyHeader != "" && tokens.AT != "" {
		ph.Set(h.atProxyHeader, tokens.AT)
	}
	if h.idtProxyHeader != "" && tokens.IDT != "" {
		ph.Set(h.idtProxyHeader, tokens.IDT)
	}

	// Claims are saved in the context with string typed key to
	// make it possible to use them from authorization middleware.
	if h.claimsKey != "" {
		//nolint:staticcheck // SA1029: should not use built-in type string as key for value; define your own type to avoid collisions
		ctx = context.WithValue(ctx, h.claimsKey, tokens)
	}

	return ctx
}

func (c *oauthContext) validateUserInfo(ui []byte, idt string) core.HTTPError {
	info := map[string]any{}
	err := json.Unmarshal(ui, &info)
	if err != nil {
		err := app.ErrAppGenUnmarshal.WithoutStack(err, map[string]any{"from": "json", "to": "map[string]any{}", "content": string(ui)})
		return utilhttp.NewHTTPError(app.ErrAppAuthnInvalidUserInfo.WithoutStack(err, nil), http.StatusUnauthorized)
	}

	uiSub, ok := info["sub"].(string)
	if !ok {
		err := app.ErrAppAuthnInvalidUserInfo.WithoutStack(err, map[string]any{"reason": "sub in UserInfo response isn't string type."})
		return utilhttp.NewHTTPError(err, http.StatusUnauthorized)
	}

	claims, err := c.jh.ValidMapClaims(idt, c.idtParseOpts...)
	if err != nil {
		err := app.ErrAppAuthnInvalidToken.WithoutStack(err, map[string]any{"name": "id token", "reason": "token validation failed.", "token": idt})
		return utilhttp.NewHTTPError(app.ErrAppAuthnInvalidUserInfo.WithoutStack(err, nil), http.StatusUnauthorized)
	}

	idtSud, _ := claims.GetSubject()
	if idtSud == "" {
		err := app.ErrAppAuthnInvalidToken.WithoutStack(err, map[string]any{"name": "id token", "reason": "sub in ID token does not exist.", "token": idt})
		return utilhttp.NewHTTPError(app.ErrAppAuthnInvalidUserInfo.WithoutStack(err, nil), http.StatusUnauthorized)
	}

	if uiSub != idtSud {
		err := app.ErrAppAuthnInvalidUserInfo.WithoutStack(err, map[string]any{"reason": "sub in UserInfo does not match sub in ID token."})
		return utilhttp.NewHTTPError(err, http.StatusUnauthorized)
	}

	return nil
}

// mapValue returns a value by searching with the given key.
// This function returns zero value if key not found.
func mapValue[T any](c jwt.MapClaims, key string) T {
	var t T
	v, ok := c[key]
	if !ok {
		return t
	}
	vv, ok := v.(T)
	if !ok {
		return t
	}
	return vv
}

type validateOption interface {
	validate(jwt.MapClaims) core.HTTPError
}

type maxAgeValidation int64
type nonceValidation string

func (v maxAgeValidation) validate(claims jwt.MapClaims) core.HTTPError {
	if int64(v) <= 0 {
		return nil
	}
	at := int64(mapValue[float64](claims, "auth_time"))
	if at == 0 {
		err := app.ErrAppAuthnInvalidToken.WithoutStack(nil, map[string]any{"name": "ID token", "reason": "auth_time not found."})
		return utilhttp.NewHTTPError(err, http.StatusUnauthorized)
	}

	if int64(v) < (time.Now().Unix() - at) {
		return reAuthenticationRequired
	}

	return nil
}

func (v nonceValidation) validate(claims jwt.MapClaims) core.HTTPError {
	if string(v) == "" {
		return nil
	}
	if string(v) != claims["nonce"] {
		err := app.ErrAppAuthnInvalidToken.WithoutStack(nil, map[string]any{"name": "ID token", "reason": "invalid nonce."})
		return utilhttp.NewHTTPError(err, http.StatusUnauthorized)
	}
	return nil
}
