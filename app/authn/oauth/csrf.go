// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package oauth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"io"
	"net/url"

	v1 "github.com/aileron-gateway/aileron-gateway/apis/app/v1"
	"github.com/aileron-gateway/aileron-gateway/core"
)

// csrfSessionKey is the key string to save CSRF parameters
// in the session.
const csrfSessionKey = "_csrf"

const (
	// stateLength is the byte length of "state" parameter
	// described in RFC 6749 The OAuth 2.0 Authorization Framework.
	// 84 bytes are encoded into 112 (=84*8/6) ascii characters
	// by Base64 URL encoding.
	//	- https://datatracker.ietf.org/doc/html/rfc6749
	stateLength = 84

	// nonceLength is the byte length of "nonce" parameter
	// described in OpenID Connect Core 1.0 incorporating errata set 2.
	// 84 bytes are encoded into 112 (=84*8/6) ascii characters
	// by Base64 URL encoding.
	//	- https://openid.net/specs/openid-connect-core-1_0.html
	nonceLength = 84

	// codeVerifierLength is the byte length of code verifier
	// described in RFC 7636 Proof Key for Code Exchange by OAuth Public Clients.
	// 84 bytes are encoded into 112 (=84*8/6) ascii characters
	// by Base64 URL encoding.
	//	- https://datatracker.ietf.org/doc/rfc7636/
	codeVerifierLength = 84
)

// codeChallengeMethod is the mapping of
// v1.PKCEMethod and its string representation.
//   - https://datatracker.ietf.org/doc/html/rfc7636
var codeChallengeMethod = map[v1.PKCEMethod]string{
	v1.PKCEMethod_S256:  "S256",
	v1.PKCEMethod_Plain: "plain",
}

// csrfStateGenerator generates CSRF parameters.
// state and nonce are described in RFC6749
// and PKCE is described in RFC7636.
//   - https://datatracker.ietf.org/doc/rfc6749/
//   - https://datatracker.ietf.org/doc/rfc7636/
type csrfStateGenerator struct {
	stateDisabled bool
	nonceDisabled bool
	pkceDisabled  bool
	method        string // PKCE Code Challenge Method
}

// new generates a new CSRF state parameters.
func (s *csrfStateGenerator) new() (*csrfStates, error) {
	csrf := &csrfStates{}

	if !s.stateDisabled {
		b := make([]byte, stateLength)
		if _, err := io.ReadFull(rand.Reader, b); err != nil {
			return nil, core.ErrCoreGenCreateComponent.WithStack(err, map[string]any{"reason": "random bytes read error"})
		}
		csrf.State = base64.RawURLEncoding.EncodeToString(b)
	}

	if !s.nonceDisabled {
		b := make([]byte, nonceLength)
		if _, err := io.ReadFull(rand.Reader, b); err != nil {
			return nil, core.ErrCoreGenCreateComponent.WithStack(err, map[string]any{"reason": "random bytes read error"})
		}
		csrf.Nonce = base64.RawURLEncoding.EncodeToString(b)
	}

	if !s.pkceDisabled {
		csrf.method = s.method
		b := make([]byte, codeVerifierLength)
		if _, err := io.ReadFull(rand.Reader, b); err != nil {
			return nil, core.ErrCoreGenCreateComponent.WithStack(err, map[string]any{"reason": "random bytes read error"})
		}
		csrf.Verifier = base64.RawURLEncoding.EncodeToString(b)
		if s.method == "plain" {
			// When "plain", code_challenge = code_verifier
			csrf.challenge = csrf.Verifier
		} else {
			// When "S256", code_challenge = BASE64URL-ENCODE(SHA256(ASCII(code_verifier)))
			d := sha256.Sum256([]byte(csrf.Verifier))
			csrf.challenge = base64.RawURLEncoding.EncodeToString(d[:])
			csrf.method = "S256"
		}
	}

	return csrf, nil
}

// csrfStates is the CSRF parameters for oauth.
//   - https://datatracker.ietf.org/doc/rfc6749/
//   - https://datatracker.ietf.org/doc/rfc7636/
type csrfStates struct {
	State     string `json:"s,omitempty" msgpack:"s,omitempty"`
	Nonce     string `json:"n,omitempty" msgpack:"n,omitempty"`
	Verifier  string `json:"v,omitempty" msgpack:"v,omitempty"`
	method    string // method is not need to be saved in the session.
	challenge string // challenge is not need to be saved in the session.
}

// set sets the CSRF parameters to the given url.Values.
// This method do nothing if nil was given by the argument.
func (s *csrfStates) set(v url.Values) {
	if v == nil {
		return
	}
	if s.State != "" {
		v.Set("state", s.State)
	}
	if s.Nonce != "" {
		v.Set("nonce", s.Nonce)
	}
	if s.challenge != "" {
		v.Set("code_challenge_method", s.method)
		v.Set("code_challenge", s.challenge)
	}
}
