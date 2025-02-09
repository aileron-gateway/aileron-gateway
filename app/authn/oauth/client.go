package oauth

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"io"
	"maps"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/aileron-gateway/aileron-gateway/app"
	"github.com/aileron-gateway/aileron-gateway/core"
	"github.com/aileron-gateway/aileron-gateway/kernel/log"
	utilhttp "github.com/aileron-gateway/aileron-gateway/util/http"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// tokenRedeemer is an interface of token introspection request client.
type tokenIntrospector interface {
	// tokenIntrospection send token introspection request.
	// The returned status will be -1 when the request was not successfully created or sent.
	// The returned jwt.MapClaims is always nil when a non-nil error is returned.
	// Caller should respect the http error returned from this method.
	tokenIntrospection(context.Context, map[string]string) (statusCode int, claims jwt.MapClaims, err core.HTTPError)
}

// tokenRedeemer is an interface of token request client.
type tokenRedeemer interface {
	// redeemToken send token request.
	// This method returns status code, token response and error.
	// The returned status will be -1 when the request was not successfully created or sent.
	// The returned TokenResponse is always nil when a non-nil error is returned.
	// Caller should respect the http error returned from this method.
	redeemToken(context.Context, map[string]string) (statusCode int, resp *TokenResponse, err core.HTTPError)
}

// requester is an interface of HTTP request client.
type requester interface {
	// doRequest send a HTTP request.
	// This method returns status code, response body and error.
	// The returned status will be -1 when the request was not successfully created or sent.
	// The returned body is always nil when a non-nil error is returned.
	doRequest(context.Context, string, map[string]string) (int, []byte, error)
}

// TokenResponse is the response body model for token requests.
type TokenResponse struct {
	IDToken          string `json:"id_token,omitempty"`
	AccessToken      string `json:"access_token,omitempty"`
	RefreshToken     string `json:"refresh_token,omitempty"`
	TokenType        string `json:"token_type,omitempty"`
	ExpiresIn        int64  `json:"expires_in,omitempty"`
	RefreshExpiresIn int64  `json:"refresh_expires_in,omitempty"`
	Scope            string `json:"scope,omitempty"`

	Error       string `json:"error,omitempty"`
	Description string `json:"error_description,omitempty"`
	URI         string `json:"error_uri,omitempty"`

	StatusCode int    `json:"-"`
	RawBody    []byte `json:"-"`
}

// tokenIntrospectionClient is a http client for token introspection.
// This implements oauth.tokenIntrospector interface.
type tokenIntrospectionClient struct {
	requester
	lg       log.Logger
	provider *provider
}

// tokenIntrospection send a token introspection request for authorization server.
// This method returns status code, token response and error.
// The returned status will be -1 when the request was not successfully created or sent.
// The returned TokenResponse is always nil when an error is returned.
// This method returns non-nil error when the response status was grater than or equal to 500.
// It should always be a server-side error when the error is not nil.
func (c *tokenIntrospectionClient) tokenIntrospection(ctx context.Context, queryParams map[string]string) (int, jwt.MapClaims, core.HTTPError) {
	status, b, err := c.doRequest(ctx, c.provider.introspectEP, queryParams)
	if err != nil {
		err := app.ErrAppAuthnIntrospection.WithStack(err, nil)
		if c.lg.Enabled(log.LvDebug) {
			c.lg.Debug(ctx, "token introspection failed", err.Name(), err.Map())
		}
		return status, nil, utilhttp.NewHTTPError(err, http.StatusInternalServerError)
	}

	if status != http.StatusOK {
		err := app.ErrAppAuthnIntrospection.WithStack(nil, map[string]any{"info": "status:" + strconv.Itoa(status) + " body:" + string(b)})
		if status >= 500 {
			return status, nil, utilhttp.NewHTTPError(err, http.StatusInternalServerError)
		}
		return status, nil, utilhttp.NewHTTPError(err, http.StatusUnauthorized)
	}

	resp := jwt.MapClaims{}
	if err := json.Unmarshal(b, &resp); err != nil {
		err := app.ErrAppGenUnmarshal.WithStack(err, map[string]any{"from": "json", "to": "MapClaims", "content": string(b)})
		return status, nil, utilhttp.NewHTTPError(err, http.StatusInternalServerError)
	}

	return status, resp, nil
}

// redeemTokenClient is a http client for token request.
// This implements oauth.tokenRedeemer interface.
type redeemTokenClient struct {
	requester
	lg       log.Logger
	provider *provider
}

// redeemToken send a token request for authorization server.
// This method returns status code, token response and error.
// The returned status will be -1 when the request was not successfully created or sent.
// The returned TokenResponse is always nil when an error is returned.
// This method returns non-nil error when the response status was grater than or equal to 500.
// It should always be a server-side error when the error is not nil.
func (c *redeemTokenClient) redeemToken(ctx context.Context, queryParams map[string]string) (int, *TokenResponse, core.HTTPError) {
	status, b, err := c.doRequest(ctx, c.provider.tokenEP, queryParams)
	if err != nil {
		err := app.ErrAppAuthnRedeemToken.WithStack(err, nil)
		if c.lg.Enabled(log.LvDebug) {
			c.lg.Debug(ctx, "redeem token failed", err.Name(), err.Map())
		}
		return status, nil, utilhttp.NewHTTPError(err, http.StatusInternalServerError)
	}

	if status != http.StatusOK {
		err := app.ErrAppAuthnRedeemToken.WithStack(nil, map[string]any{"info": "status:" + strconv.Itoa(status) + " body:" + string(b)})
		if status >= 500 {
			return status, nil, utilhttp.NewHTTPError(err, http.StatusInternalServerError)
		}
		return status, nil, utilhttp.NewHTTPError(err, http.StatusUnauthorized)
	}

	resp := &TokenResponse{
		StatusCode: status,
		RawBody:    b,
	}
	if err := json.Unmarshal(b, resp); err != nil {
		err := app.ErrAppGenUnmarshal.WithStack(err, map[string]any{"from": "json", "to": "TokenResponseModel", "content": string(b)})
		return status, nil, utilhttp.NewHTTPError(err, http.StatusInternalServerError)
	}

	return status, resp, nil
}

type clientRequester struct {
	lg log.Logger

	client           *client
	rt               http.RoundTripper
	clientAuthMethod clientAuthMethod
	extraHeader      url.Values
	extraQuery       url.Values
}

// doRequest send an auth request with client authentication.
// This method returns status code, raw body and error.
// The returned status will be -1 when the request was not successfully created or sent.
// The body is always nil when an error is returned.
// It should always be a server-side error when the error is not nil.
func (c *clientRequester) doRequest(ctx context.Context, endpoint string, queryParams map[string]string) (int, []byte, error) {
	q := make(url.Values, len(c.extraQuery)+len(queryParams)+2)
	maps.Copy(q, c.extraQuery)
	for k, v := range queryParams {
		q.Set(k, v)
	}

	if c.clientAuthMethod == clientAuthForm {
		q.Add("client_id", c.client.id)
		q.Add("client_secret", c.client.secret)
	}

	if c.clientAuthMethod == clientAuthJWT || c.clientAuthMethod == clientAuthPrivateKeyJWT {
		id, err := uuid.NewRandom()
		if err != nil {
			err := app.ErrAppAuthnGenerateClientAssertion.WithStack(err, map[string]any{"reason": "failed to generate ID claim"})
			return -1, nil, err
		}

		token, err := c.client.jh.TokenWithClaims(
			&jwt.RegisteredClaims{
				Issuer:    c.client.id,
				Subject:   c.client.id,
				Audience:  jwt.ClaimStrings{endpoint},
				ID:        id.String(),
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute)),
			})
		if err != nil {
			err := app.ErrAppAuthnGenerateClientAssertion.WithStack(err, map[string]any{"reason": "failed to generate token with claim"})
			return -1, nil, err
		}

		ca, err := c.client.jh.SignedString(token)
		if err != nil {
			err := app.ErrAppAuthnGenerateClientAssertion.WithStack(err, map[string]any{"reason": "failed to sign token"})
			return -1, nil, err
		}

		q.Add("client_assertion_type", "urn:ietf:params:oauth:client-assertion-type:jwt-bearer")
		q.Add("client_assertion", ca)
	}

	if c.clientAuthMethod == clientAuthTLSClientAuth || c.clientAuthMethod == clientAuthSelfSignedTLSClientAuth {
		q.Add("client_id", c.client.id)
	}

	body := []byte(q.Encode())
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		err := app.ErrAppGenCreateRequest.WithStack(err, map[string]any{"method": http.MethodPost, "url": endpoint, "body": string(body)})
		return -1, nil, err
	}

	maps.Copy(req.Header, c.extraHeader)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if c.clientAuthMethod == clientAuthBasic {
		req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(c.client.id+":"+c.client.secret)))
	}

	res, err := c.rt.RoundTrip(req)
	if err != nil {
		err := app.ErrAppGenRoundTrip.WithStack(err, map[string]any{"method": http.MethodPost, "url": endpoint, "body": string(body)})
		return -1, nil, err
	}
	defer res.Body.Close()

	b, err := io.ReadAll(res.Body)
	if err != nil {
		err := app.ErrAppGenReadHTTPBody.WithStack(err, map[string]any{"direction": "response", "body": string(b)})
		return res.StatusCode, nil, err
	}

	if !strings.Contains(res.Header.Get("Content-Type"), "application/json") {
		info := map[string]any{
			"method": http.MethodPost,
			"url":    endpoint,
			"status": strconv.Itoa(res.StatusCode),
			"type":   res.Header.Get("Content-Type"),
			"body":   string(b),
		}
		err := app.ErrAppGenInvalidResponse.WithStack(err, info)
		return res.StatusCode, nil, err
	}

	if c.lg.Enabled(log.LvDebug) {
		info := map[string]any{
			"method": http.MethodPost,
			"url":    endpoint,
			"status": strconv.Itoa(res.StatusCode),
			"type":   res.Header.Get("Content-Type"),
			"body":   string(b),
		}
		attr := log.NewCustomAttrs("authz", info)
		c.lg.Debug(ctx, "http response obtained from authorization server", attr.Name(), attr.Map())
	}

	return res.StatusCode, b, nil
}
