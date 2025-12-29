// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package oauth

import (
	"net/http"

	v1 "github.com/aileron-gateway/aileron-gateway/apis/app/v1"
	"github.com/aileron-gateway/aileron-gateway/app"
	"github.com/aileron-gateway/aileron-gateway/internal/security"
	"github.com/aileron-gateway/aileron-gateway/kernel/log"
	"github.com/golang-jwt/jwt/v5"
)

var responseModeMethods = map[v1.ResponseModeMethod]string{
	v1.ResponseModeMethod_ResponseModeJWT:         "jwt",
	v1.ResponseModeMethod_ResponseModeQueryJWT:    "query.jwt",
	v1.ResponseModeMethod_ResponseModeFragmentJWT: "fragment.jwt",
	v1.ResponseModeMethod_ResponseModeFormPostJWT: "form_post.jwt",
}

type jarmValidator struct {
	responseMode v1.ResponseModeMethod
	jh           *security.JWTHandler
}

func (j *jarmValidator) valid(oc *oauthContext, r *http.Request, stateDisabled bool) (string, string, error) {
	var jarmString string
	if j.responseMode == v1.ResponseModeMethod_ResponseModeQueryJWT || j.responseMode == v1.ResponseModeMethod_ResponseModeJWT {
		jarmString = r.URL.Query().Get("response")
	}
	if j.responseMode == v1.ResponseModeMethod_ResponseModeFormPostJWT {
		if err := r.ParseForm(); err != nil {
			return "", "", app.ErrAppAuthnInvalidParameters.WithStack(err, map[string]any{"name": "jarm response", "reason": "jarm response not found"})
		}
		jarmString = r.FormValue("response")
	}

	claims, err := j.jh.ValidMapClaims(jarmString, jwt.WithIssuer(oc.provider.issuer), jwt.WithAudience(oc.client.audience))
	if err != nil {
		if oc.lg.Enabled(log.LvDebug) {
			err := app.ErrAppAuthnParseWithClaims.WithStack(err, map[string]any{"jwt": jarmString})
			oc.lg.Debug(r.Context(), "failed to parse claims", err.KeyValues()...)
		}
		return "", "", app.ErrAppAuthnInvalidParameters.WithStack(err, map[string]any{"name": "jarm response", "reason": "jarm response validation failed"})
	}

	code, ok := claims["code"].(string)
	if !ok {
		return "", "", app.ErrAppAuthnInvalidParameters.WithStack(err, map[string]any{"name": "jarm response", "reason": "invalid code"})
	}

	if stateDisabled {
		return code, "", nil
	}

	state, ok := claims["state"].(string)
	if !ok {
		return "", "", app.ErrAppAuthnInvalidParameters.WithStack(err, map[string]any{"name": "jarm response", "reason": "invalid state"})
	}

	return code, state, nil
}
