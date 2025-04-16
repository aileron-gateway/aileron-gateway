// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package oauth

import (
	"cmp"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	v1 "github.com/aileron-gateway/aileron-gateway/apis/app/v1"
	"github.com/aileron-gateway/aileron-gateway/app"
	"github.com/aileron-gateway/aileron-gateway/core"
	"github.com/aileron-gateway/aileron-gateway/kernel/log"
	"github.com/aileron-gateway/aileron-gateway/util/security"
)

// clientAuthMethod is the type of method
// used for OAuth client authentication.
//   - https://datatracker.ietf.org/doc/rfc6749/
type clientAuthMethod int

const (
	clientAuthBasic clientAuthMethod = iota // Basic authentication. default
	clientAuthForm                          // Form authentication
	clientAuthJWT
	clientAuthPrivateKeyJWT
	clientAuthTLSClientAuth
	clientAuthSelfSignedTLSClientAuth
)

var clientAuthMethods = map[v1.ClientAuthMethod]clientAuthMethod{
	v1.ClientAuthMethod_BasicAuth:               clientAuthBasic,
	v1.ClientAuthMethod_FormAuth:                clientAuthForm,
	v1.ClientAuthMethod_ClientSecretJWT:         clientAuthJWT,
	v1.ClientAuthMethod_PrivateKeyJWT:           clientAuthPrivateKeyJWT,
	v1.ClientAuthMethod_TLSClientAuth:           clientAuthTLSClientAuth,
	v1.ClientAuthMethod_SelfSignedTLSClientAuth: clientAuthSelfSignedTLSClientAuth,
}

// newClient returns new OAuth client.
func newClient(spec *v1.OAuthClient) (*client, error) {
	var jh *security.JWTHandler
	if spec.JWTHandler != nil {
		j, err := security.NewJWTHandler(spec.JWTHandler, nil)
		if err != nil {
			return nil, err
		}
		jh = j
	}

	return &client{
		id:       spec.ID,
		secret:   spec.Secret,
		audience: cmp.Or(spec.Audience, spec.ID), // Use client ID by default.Audience is prior to the ID.
		scope:    strings.Join(spec.Scopes, " "),
		jh:       jh,
	}, nil
}

// client is the OAuth client.
type client struct {
	id       string // id is the OAuth client ID.
	secret   string // secret is the OAuth client secret.
	audience string // audience is the audience defined in OAuth.
	scope    string // scope is space delimitated scopes.
	jh       *security.JWTHandler
}

// newProvider returns new OAuth provider.
func newProvider(spec *v1.OAuthProvider, rt http.RoundTripper) (*provider, error) {
	errs := []error{}
	validURL := func(base string, path string) string {
		if path == "" {
			return ""
		}
		if base == "" {
			base = path // base cannot be empty for url.JoinPath.
			path = ""   // path can be empty.
		}
		rawURL, err := url.JoinPath(base, path)
		if err != nil {
			errs = append(errs, err)
			return ""
		}
		return rawURL
	}

	if rt == nil {
		rt = http.DefaultTransport
	}

	p := &provider{
		lg:              log.GlobalLogger(log.DefaultLoggerName),
		rt:              rt,
		issuer:          spec.Issuer,
		authorizationEP: validURL(spec.BaseURL, spec.Endpoints.Authorization),
		tokenEP:         validURL(spec.BaseURL, spec.Endpoints.Token),
		introspectEP:    validURL(spec.BaseURL, spec.Endpoints.Introspection),
		userinfoEP:      validURL(spec.BaseURL, spec.Endpoints.Userinfo),
		revocationEP:    validURL(spec.BaseURL, spec.Endpoints.Revocation),
		jwksEP:          validURL(spec.BaseURL, spec.Endpoints.JWKs),
		discoveryEP:     validURL(spec.BaseURL, spec.Endpoints.Discovery),
	}

	if len(errs) > 0 {
		err := errors.Join(errs...)
		err = core.ErrCoreGenCreateComponent.WithStack(err, map[string]any{"reason": "failed to create provider"})
		return nil, err
	}

	go p.discover()
	return p, nil
}

// ProviderMetadata holds metadata of an authorization server.
// Other parameters can be added when they are necessary.
//
//   - https://datatracker.ietf.org/doc/html/rfc8414#section-7.1.2
//   - https://openid.net/specs/openid-connect-discovery-1_0.html#ProviderConfigurationResponse
type providerMetadata struct {
	Issuer          string `json:"issuer,omitempty"`
	AuthzEP         string `json:"authorization_endpoint,omitempty"`
	TokenEP         string `json:"token_endpoint,omitempty"`
	UserinfoEP      string `json:"userinfo_endpoint,omitempty"`
	JWKsEP          string `json:"jwks_uri,omitempty"`
	IntrospectionEP string `json:"introspection_endpoint,omitempty"`
	RevocationEP    string `json:"revocation_endpoint,omitempty"`
}

// provider is the OAuth provider.
type provider struct {
	issuer          string
	authorizationEP string
	tokenEP         string
	userinfoEP      string
	jwksEP          string
	introspectEP    string
	revocationEP    string
	discoveryEP     string

	lg log.Logger

	// rt is the round tripper that is used
	// for OpenID discovering.
	rt http.RoundTripper

	ticker *time.Ticker

	// close ends up the discovery loop.
	close chan struct{}
}

// discover discovers provider metadata from discovery endpoints.
// Do not set the discovery endpoint to prevent this method to work.
func (p *provider) discover() {
	if p.discoveryEP == "" {
		return
	}

	ctx := context.Background()

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, p.discoveryEP, nil)
	if err != nil {
		err := app.ErrAppGenCreateRequest.WithoutStack(err, map[string]any{"method": http.MethodGet, "url": p.discoveryEP})
		p.lg.Error(ctx, "failed to create discovery request", err.Name(), err.Map())
		return
	}

	validURL := func(target string, src string) string {
		if target != "" {
			return target
		}
		if src == "" {
			return ""
		}
		if _, err := url.Parse(src); err != nil {
			return ""
		}
		return src
	}

	if p.ticker == nil {
		p.ticker = time.NewTicker(5 * time.Second)
	}
	defer p.ticker.Stop()

loop:
	for {
		res, err := p.rt.RoundTrip(req)
		if err != nil {
			err := app.ErrAppGenRoundTrip.WithoutStack(err, map[string]any{"method": http.MethodGet, "url": p.discoveryEP})
			p.lg.Info(ctx, "failed to send discovery request", err.Name(), err.Map())
		} else {
			body, err := io.ReadAll(res.Body)
			if err != nil {
				err := app.ErrAppGenReadHTTPBody.WithoutStack(err, map[string]any{"direction": "request", "body": string(body)})
				p.lg.Info(ctx, "failed to read response body", err.Name(), err.Map())
			}
			res.Body.Close()

			metadata := &providerMetadata{}
			if err := json.Unmarshal(body, metadata); err != nil {
				err := app.ErrAppGenUnmarshal.WithoutStack(err, map[string]any{"from": "byte", "to": "providerMetadata", "content": string(body)})
				p.lg.Info(ctx, "failed to unmarshal response json", err.Name(), err.Map())
			}

			if res.StatusCode == http.StatusOK {
				p.lg.Info(ctx, "provider discovery succeeded: ", "metadata", metadata)
				p.issuer = validURL(p.issuer, metadata.Issuer)
				p.authorizationEP = validURL(p.authorizationEP, metadata.AuthzEP)
				p.tokenEP = validURL(p.tokenEP, metadata.TokenEP)
				p.introspectEP = validURL(p.introspectEP, metadata.IntrospectionEP)
				p.userinfoEP = validURL(p.userinfoEP, metadata.UserinfoEP)
				p.revocationEP = validURL(p.revocationEP, metadata.RevocationEP)
				p.jwksEP = validURL(p.jwksEP, metadata.JWKsEP)
				break
			}
		}

		// Wait until closed or timer finished.
		select {
		case <-p.close:
			break loop
		case <-p.ticker.C:
		}
	}
}
