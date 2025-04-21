// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package oauth

import (
	"crypto/sha256"
	"encoding/base64"
	"net/url"
	"strings"
	"time"

	"github.com/aileron-gateway/aileron-gateway/app"
	"github.com/aileron-gateway/aileron-gateway/util/security"
	"github.com/golang-jwt/jwt/v5"
)

type requestObjectGenerator struct {
	requestURI     string
	requestURIPath string
	jh             *security.JWTHandler // JWTHandler is a required parameter.
	exp            int64
	nbf            int64
	cacheDisabled  bool
}

func (g *requestObjectGenerator) new(oc *oauthContext, base url.Values, opt url.Values) (string, error) {
	nbf := time.Now().Unix() + g.nbf
	requestObjectClaims := jwt.MapClaims{
		"iss":   oc.client.id,
		"aud":   jwt.ClaimStrings{oc.provider.issuer}, // 10. shall send the aud claim in the request object as the OP's Issuer Identifier URL;
		"exp":   nbf + g.exp,                          // 11. shall send an exp claim in the request object that has a lifetime of no longer than 60 minutes;
		"nbf":   nbf,                                  // 14. shall send a nbf claim in the request object;
		"scope": oc.client.scope,
	}

	for k, v := range base {
		requestObjectClaims[k] = strings.Join(v, " ")
	}
	for k, v := range opt {
		requestObjectClaims[k] = strings.Join(v, " ")
	}

	token, err := g.jh.TokenWithClaims(requestObjectClaims)
	if err != nil {
		return "", app.ErrAppAuthnGenerateTokenWithClaims.WithStack(err, nil)
	}

	ro, err := g.jh.SignedString(token)
	if err != nil {
		return "", app.ErrAppAuthnSignToken.WithStack(err, nil)
	}

	return ro, nil
}

func (g *requestObjectGenerator) getRequestURI(ro string) string {
	if !g.cacheDisabled {
		return g.requestURI
	}

	hash := sha256.Sum256([]byte(ro))
	hashEncoded := base64.RawURLEncoding.EncodeToString(hash[:])

	return g.requestURI + "#" + hashEncoded
}
