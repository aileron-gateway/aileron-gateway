// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package oauth

import (
	"crypto/rand"
	"fmt"
	"net/http"
	"strings"
	"time"

	v1 "github.com/aileron-gateway/aileron-gateway/apis/app/v1"
	"github.com/aileron-gateway/aileron-gateway/app"
	utilhttp "github.com/aileron-gateway/aileron-gateway/util/http"
	"github.com/google/uuid"
)

// newResourceServerHandler returns a new instance of resource server authentication handler.
func newResourceServerHandler(bh *baseHandler, spec *v1.ResourceServerHandler) *resourceServerHandler {
	if spec.HeaderKey == "" {
		spec.HeaderKey = "Authorization"
	}
	return &resourceServerHandler{
		baseHandler: bh,
		headerKey:   spec.HeaderKey,
		fapiEnabled: spec.EnabledFAPI,
	}
}

type resourceServerHandler struct {
	*baseHandler

	// header key is the HTTP request header name
	// to extract a bearer token from it.
	headerKey string

	// fapiEnabled is the flag to enable
	// validation described in FAPI specifications.
	fapiEnabled bool
}

// ServeAuthn handle authentication.
// return newRequest, isAuthenticated, shouldReturn, error.
func (h *resourceServerHandler) ServeAuthn(w http.ResponseWriter, r *http.Request) (*http.Request, app.AuthResult, bool, error) {
	// Find oauth context that can handle this authentication flow.
	// If not found, yield authentication to other authentication handlers.
	oc, err := h.oauthContext(r)
	if err != nil {
		return nil, app.AuthContinue, false, nil
	}

	// If token not found in the header,
	// yield authentication to other authentication handlers.
	token := r.Header.Get(h.headerKey)
	if token == "" {
		return nil, app.AuthContinue, false, nil
	}
	token = strings.TrimPrefix(token, "Bearer ")

	tokens := &OAuthTokens{
		Context:  oc.name, // Context name handling this authentication flow.
		AT:       token,   // Requested access token as-is.
		ATExp:    0,       // Expected to be filled in validOauthClaims method.
		ATClaims: nil,     // Expected to be filled in validOauthClaims method.
	}

	// Check if the token is valid or not.
	if err := oc.validOauthClaims(r.Context(), tokens); err != nil {
		fmt.Println("invalid", tokens, err)
		if e, ok := err.(interface{ Header() http.Header }); ok {
			e.Header().Add("WWW-Authenticate", `Bearer error="invalid_token"`)
		}
		return nil, app.AuthFailed, true, err
	}
	fmt.Println("valid", tokens, err)

	if h.fapiEnabled {
		// Validate the client certificate required by FAPI Part2 - 6.2.1. Protected resources provisions.
		// 2. shall adhere to the requirements in MTLS.
		// https://openid.net/specs/openid-financial-api-part-2-1_0-final.html
		// Based on OAuth 2.0 Mutual-TLS Client Authentication and Certificate-Bound Access Tokens.
		// https://datatracker.ietf.org/doc/html/rfc8705
		cnf := mapValue[map[string]any](tokens.ATClaims, "cnf")
		if err := validateCert(cnf, r.TLS); err != nil {
			err := app.ErrAppAuthnAuthentication.WithoutStack(err, nil)
			httpErr := utilhttp.NewHTTPError(err, http.StatusUnauthorized)
			httpErr.Header().Add("WWW-Authenticate", `Bearer error="invalid_token"`)
			return nil, app.AuthFailed, true, httpErr
		}

		// Set the FAPI interaction ID required by FAPI Part1 - 6.2.1. Protected resources provisions.
		// 11. shall set the response header x-fapi-interaction-id to the value received from the corresponding FAPI client request header or to a RFC4122 UUID value
		// if the request header was not provided to track the interaction, e.g., x-fapi-interaction-id: c770aef3-6784-41f7-8e0e-ff5f97bddb3a;
		// 12. shall log the value of x-fapi-interaction-id in the log entry;
		// https://openid.net/specs/openid-financial-api-part-1-1_0.html
		id := r.Header.Get("x-fapi-interaction-id")
		if id == "" {
			uid, err := uuid.NewRandomFromReader(rand.Reader)
			if err != nil {
				err := app.ErrAppAuthnAuthentication.WithStack(err, nil)
				return nil, app.AuthFailed, true, utilhttp.NewHTTPError(err, http.StatusInternalServerError)
			}
			id = uid.String()
		}
		w.Header().Set("x-fapi-interaction-id", id)
		h.lg.Debug(r.Context(), "x-fapi-interaction-id:"+id)

		// Set the date header required by FAPI Part1 - 6.2.1. Protected resources provisions.
		// 10. shall send the server date in HTTP Date header as in Section 7.1.1.2 of RFC7231;
		// https://openid.net/specs/openid-financial-api-part-1-1_0.html
		w.Header().Set("Date", time.Now().Format(http.TimeFormat))

		// Set the date header required by FAPI Part1 - 6.2.1. Protected resources provisions.
		// 6. shall identify the associated entity to the access token;
		// https://openid.net/specs/openid-financial-api-part-1-1_0.html
		h.lg.Debug(r.Context(), "Access token verified.", "claims", tokens.ATClaims)
	}

	// Save tokens in the context so the succeeding middleware
	// such as authorization middleware can use them.
	ctx := oc.contextWithToken(r.Context(), tokens)
	r = r.WithContext(ctx)

	return r, app.AuthSucceeded, false, nil // Authentication succeeded !!
}
