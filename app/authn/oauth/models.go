package oauth

import (
	"crypto/sha256"
	"crypto/tls"
	"encoding/base64"
	"errors"

	"github.com/aileron-gateway/aileron-gateway/app"
)

// OAuthTokens holds token data.
type OAuthTokens struct {
	// Context is the context name that are used to obtain IDT or AT.
	Context string `json:"context" msgpack:"context"`

	// IDT is the original ID Token string.
	// IDT can be an empty string.
	// IDT can be a JWT or an Opaque token.
	IDT string `json:"idt" msgpack:"idt"`
	// IDTExp is the expiration time of the ID token.
	// This field will be 0 when the IDT
	// does not contain "exp" claim.
	IDTExp int64 `json:"idt_exp" msgpack:"idt_exp"`
	// IDTClaims is the claims that are extracted from the ID token.
	IDTClaims map[string]any `json:"idt_claims" msgpack:"idt_claims"`

	// AT is the original Access Token string.
	// AT can be an empty string.
	// IDT can be a JWT or an Opaque token.
	AT string `json:"at" msgpack:"at"`
	// ATExp is the expiration time of the Access token.
	// This field will be 0 when the AT
	// does not contain "exp" claim.
	ATExp int64 `json:"at_exp" msgpack:"at_exp"`
	// ATClaims is the claims that are extracted from the access token.
	ATClaims map[string]any `json:"at_claims" msgpack:"at_claims"`

	// RT is the original Refresh Token string.
	// RT can be an empty string.
	// IDT can be a JWT or an Opaque token.
	RT string `json:"rt" msgpack:"rt"`

	// updated is the flag to re-encode session data.
	// If any data in this struct has changed,
	// updated should be set to true for updating session.
	updated bool `json:"-" msgpack:"-"`
}

// Validate the client certificate required by FAPI Part2 - 6.2.1. Protected resources provisions.
// 2. shall adhere to the requirements in MTLS.
// https://openid.net/specs/openid-financial-api-part-2-1_0-final.html
// Based on OAuth 2.0 Mutual-TLS Client Authentication and Certificate-Bound Access Tokens.
// https://datatracker.ietf.org/doc/html/rfc8705
func validateCert(cnf map[string]any, state *tls.ConnectionState) error {
	if cnf == nil {
		err := errors.New("authn/oauth: cnf claim not found")
		err = app.ErrAppAuthnInvalidCredential.WithoutStack(err, map[string]any{"purpose": "mTLS", "reason": "insufficient claims"})
		return err
	}

	if state == nil {
		err := errors.New("authn/oauth: tls state not found")
		err = app.ErrAppAuthnInvalidCredential.WithoutStack(err, map[string]any{"purpose": "mTLS", "reason": "no tls state"})
		return err
	}

	c, ok := cnf["x5t#S256"].(string)
	if !ok {
		err := errors.New("authn/oauth: cnf claims does not contains valid value for x5t#S256")
		err = app.ErrAppAuthnInvalidCredential.WithoutStack(err, map[string]any{"purpose": "mTLS", "reason": "insufficient cnf claims"})
		return err
	}

	if len(state.PeerCertificates) == 0 {
		err := errors.New("authn/oauth: client does not provide certificate")
		err = app.ErrAppAuthnInvalidCredential.WithoutStack(err, map[string]any{"purpose": "mTLS", "reason": "peer certificates not found"})
		return err
	}

	hash := sha256.Sum256(state.PeerCertificates[0].Raw)
	thumbprint := base64.RawURLEncoding.EncodeToString(hash[:])

	if thumbprint != c {
		err := errors.New("authn/oauth: thumbprint mismatched. expect:" + thumbprint + " got:" + c)
		err = app.ErrAppAuthnInvalidCredential.WithoutStack(err, map[string]any{"purpose": "mTLS", "reason": "thumbprint mismatch"})
		return err
	}

	return nil
}
