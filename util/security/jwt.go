// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package security

import (
	"context"
	"encoding/base32"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"sync"

	v1 "github.com/aileron-gateway/aileron-gateway/apis/app/v1"
	"github.com/aileron-gateway/aileron-gateway/app"
	"github.com/aileron-gateway/aileron-gateway/core"
	"github.com/aileron-gateway/aileron-gateway/kernel/hash"
	"github.com/golang-jwt/jwt/v5"
	"github.com/lestrrat-go/jwx/v2/jwk"
)

var (
	ErrNilSpec = errors.New("util/security: nil spec was given")

	ErrInvalidAlg  = errors.New("util/security: invalid algorithm")
	ErrNoKey       = errors.New("util/security: no key specified")
	ErrInvalidType = errors.New("util/security: invalid key type")

	ErrKeyNotFound  = errors.New("util/security: validating key was not found")
	ErrNoKid        = errors.New("util/security: kid is not in the JWT header")
	ErrNoAlg        = errors.New("util/security: alg is not in the JWT header")
	ErrWrongAlg     = errors.New("util/security: wrong algorithm is used for the key")
	ErrNoSigningKey = errors.New("util/security: no keys found for signing a JWT")
	ErrNilToken     = errors.New("util/security: token or token's header is nil")
)

// Algorithm is the type of signing algorithm for JWT.
// Algorithms for JWT are described in RFC 7518 "JSON Web Algorithms (JWA)".
// EdDSA algorithm is described in RFC 8037
// "CFRG Elliptic Curve Diffie-Hellman (ECDH) and Signatures in JSON Object Signing and Encryption (JOSE)".
//   - https://datatracker.ietf.org/doc/rfc7518/
//   - https://datatracker.ietf.org/doc/rfc8037/
type Algorithm string

const (
	NONE  Algorithm = "none"  // No digital signature or MAC
	ES256 Algorithm = "ES256" // ECDSA using P-256 and SHA-256
	ES384 Algorithm = "ES384" // ECDSA using P-384 and SHA-384
	ES512 Algorithm = "ES512" // ECDSA using P-521 and SHA-512
	EdDSA Algorithm = "EdDSA" // EdDSA using Ed25519
	HS256 Algorithm = "HS256" // HMAC using SHA-256
	HS384 Algorithm = "HS384" // HMAC using SHA-384
	HS512 Algorithm = "HS512" // HMAC using SHA-512
	RS256 Algorithm = "RS256" // RSASSA-PKCS1-v1_5 using SHA-256
	RS384 Algorithm = "RS384" // RSASSA-PKCS1-v1_5 using SHA-384
	RS512 Algorithm = "RS512" // RSASSA-PKCS1-v1_5 using SHA-512
	PS256 Algorithm = "PS256" // RSASSA-PSS using SHA-256 and MGF1 with SHA-256
	PS384 Algorithm = "PS384" // RSASSA-PSS using SHA-384 and MGF1 with SHA-384
	PS512 Algorithm = "PS512" // RSASSA-PSS using SHA-512 and MGF1 with SHA-512
)

var (
	// JWTAlgorithm is the JWT signing algorithm.
	JWTAlgorithm = map[v1.SigningKeyAlgorithm]Algorithm{
		v1.SigningKeyAlgorithm_NONE:  NONE,
		v1.SigningKeyAlgorithm_ES256: ES256,
		v1.SigningKeyAlgorithm_ES384: ES384,
		v1.SigningKeyAlgorithm_ES512: ES512,
		v1.SigningKeyAlgorithm_EdDSA: EdDSA,
		v1.SigningKeyAlgorithm_HS256: HS256,
		v1.SigningKeyAlgorithm_HS384: HS384,
		v1.SigningKeyAlgorithm_HS512: HS512,
		v1.SigningKeyAlgorithm_RS256: RS256,
		v1.SigningKeyAlgorithm_RS384: RS384,
		v1.SigningKeyAlgorithm_RS512: RS512,
		v1.SigningKeyAlgorithm_PS256: PS256,
		v1.SigningKeyAlgorithm_PS384: PS384,
		v1.SigningKeyAlgorithm_PS512: PS512,
	}

	// SigningMethods holds the JWT singing methods.
	// Following keys are available.
	// 	- jwt.NONE  for No digital signature or MAC
	// 	- jwt.ES256 for ECDSA using P-256 and SHA-256
	// 	- jwt.ES384 for ECDSA using P-384 and SHA-384
	// 	- jwt.ES512 for ECDSA using P-521 and SHA-512
	// 	- jwt.EdDSA for EdDSA using Ed25519
	// 	- jwt.HS256 for HMAC using SHA-256
	// 	- jwt.HS384 for HMAC using SHA-384
	// 	- jwt.HS512 for HMAC using SHA-512
	// 	- jwt.RS256 for RSASSA-PKCS1-v1_5 using SHA-256
	// 	- jwt.RS384 for RSASSA-PKCS1-v1_5 using SHA-384
	// 	- jwt.RS512 for RSASSA-PKCS1-v1_5 using SHA-512
	// 	- jwt.PS256 for RSASSA-PSS using SHA-256 and MGF1 with SHA-256
	// 	- jwt.PS384 for RSASSA-PSS using SHA-384 and MGF1 with SHA-384
	// 	- jwt.PS512 for RSASSA-PSS using SHA-512 and MGF1 with SHA-512
	SigningMethods = map[Algorithm]jwt.SigningMethod{
		NONE:  jwt.SigningMethodNone,
		ES256: jwt.SigningMethodES256,
		ES384: jwt.SigningMethodES384,
		ES512: jwt.SigningMethodES512,
		EdDSA: jwt.SigningMethodEdDSA,
		HS256: jwt.SigningMethodHS256,
		HS384: jwt.SigningMethodHS384,
		HS512: jwt.SigningMethodHS512,
		RS256: jwt.SigningMethodRS256,
		RS384: jwt.SigningMethodRS384,
		RS512: jwt.SigningMethodRS512,
		PS256: jwt.SigningMethodPS256,
		PS384: jwt.SigningMethodPS384,
		PS512: jwt.SigningMethodPS512,
	}

	// HashAlgorithm holds the hash functions for JWT singing.
	// Following keys are available.
	// 	- jwt.ES256 for ECDSA using P-256 and SHA-256
	// 	- jwt.ES384 for ECDSA using P-384 and SHA-384
	// 	- jwt.ES512 for ECDSA using P-521 and SHA-512
	// 	- jwt.EdDSA for EdDSA using Ed25519
	// 	- jwt.HS256 for HMAC using SHA-256
	// 	- jwt.HS384 for HMAC using SHA-384
	// 	- jwt.HS512 for HMAC using SHA-512
	// 	- jwt.RS256 for RSASSA-PKCS1-v1_5 using SHA-256
	// 	- jwt.RS384 for RSASSA-PKCS1-v1_5 using SHA-384
	// 	- jwt.RS512 for RSASSA-PKCS1-v1_5 using SHA-512
	// 	- jwt.PS256 for RSASSA-PSS using SHA-256 and MGF1 with SHA-256
	// 	- jwt.PS384 for RSASSA-PSS using SHA-384 and MGF1 with SHA-384
	// 	- jwt.PS512 for RSASSA-PSS using SHA-512 and MGF1 with SHA-512
	HashAlgorithm = map[Algorithm]hash.HashFunc{
		ES256: hash.SHA256,
		ES384: hash.SHA384,
		ES512: hash.SHA512,
		EdDSA: hash.SHA512,
		HS256: hash.SHA256,
		HS384: hash.SHA384,
		HS512: hash.SHA512,
		RS256: hash.SHA256,
		RS384: hash.SHA384,
		RS512: hash.SHA512,
		PS256: hash.SHA256,
		PS384: hash.SHA384,
		PS512: hash.SHA512,
	}
)

// SigningKey holds key information for signing JWTs.
type SigningKey struct {
	// kid is the key id in the JWT's header.
	kid string
	// header is the JWT's header values.
	// The header values are associated with the key.
	header map[string]any
	// method is the signing method to use with the key.
	method jwt.SigningMethod
	// key is the key value.
	key any
}

// SigningKeys returns list of JWT signing keys.
// This function will panic when the specs contains nil.
func SigningKeys(specs ...*v1.SigningKeySpec) ([]*SigningKey, error) {
	keys := make([]*SigningKey, 0, len(specs))
	for _, spec := range specs {
		alg, ok := JWTAlgorithm[spec.Algorithm]
		if !ok {
			return nil, app.ErrAppUtilGenerateJWTKey.WithStack(ErrInvalidAlg, nil)
		}

		kid := spec.KeyID
		if kid == "" {
			// KeyID in the spec is optional.
			// So, use the calculated value from the signature of the spec when not set.
			b, _ := json.Marshal(spec.JWTHeader)
			signature := spec.Algorithm.String() + spec.KeyType.String() + spec.KeyFilePath + string(b)
			kid = base32.StdEncoding.EncodeToString(hash.SHA1(hash.SHA3_512([]byte(signature))))
		}

		header := map[string]any{
			"alg": alg,
			"kid": kid,
			"typ": "JWT",
		}
		for k, v := range spec.JWTHeader { // Add user defined JWT headers.
			header[k] = v
		}

		var b []byte // Key bytes.
		if spec.KeyFilePath != "" {
			content, err := os.ReadFile(spec.KeyFilePath)
			if err != nil {
				return nil, app.ErrAppUtilGenerateJWTKey.WithStack(err, nil)
			}
			b = content
		} else {
			content, err := base64.StdEncoding.DecodeString(spec.KeyString)
			if err != nil {
				return nil, app.ErrAppUtilGenerateJWTKey.WithStack(err, nil)
			}
			b = content
		}

		k, err := parseKey(spec.KeyType, alg, b)
		if err != nil {
			return nil, app.ErrAppUtilGenerateJWTKey.WithStack(err, nil)
		}

		keys = append(keys, &SigningKey{
			kid:    kid,
			header: header,
			method: SigningMethods[alg],
			key:    k,
		})
	}

	return keys, nil
}

func parseKey(typ v1.SigningKeyType, alg Algorithm, b []byte) (any, error) {
	if len(b) == 0 && alg != NONE {
		// Key content is required when the alg is not NONE.
		return nil, ErrNoKey
	}

	switch typ {
	case v1.SigningKeyType_PRIVATE:
		return parsePrivateKey(alg, b)
	case v1.SigningKeyType_PUBLIC:
		return parsePublicKey(alg, b)
	case v1.SigningKeyType_COMMON:
		switch alg {
		case HS256, HS384, HS512:
			return b, nil
		}
	default:
		switch alg {
		case NONE:
			return jwt.UnsafeAllowNoneSignatureType, nil
		}
	}

	return nil, ErrInvalidType
}

func parsePrivateKey(alg Algorithm, pem []byte) (any, error) {
	switch alg {
	case RS256, RS384, RS512, PS256, PS384, PS512:
		return jwt.ParseRSAPrivateKeyFromPEM(pem)
	case ES256, ES384, ES512:
		return jwt.ParseECPrivateKeyFromPEM(pem)
	case EdDSA:
		return jwt.ParseEdPrivateKeyFromPEM(pem)
	}
	return nil, ErrInvalidAlg
}

func parsePublicKey(alg Algorithm, pem []byte) (any, error) {
	switch alg {
	case RS256, RS384, RS512, PS256, PS384, PS512:
		return jwt.ParseRSAPublicKeyFromPEM(pem)
	case ES256, ES384, ES512:
		return jwt.ParseECPublicKeyFromPEM(pem)
	case EdDSA:
		return jwt.ParseEdPublicKeyFromPEM(pem)
	}
	return nil, ErrInvalidAlg
}

// NewJWTHandler returns a new JWTHandler.
func NewJWTHandler(spec *v1.JWTHandlerSpec, rt http.RoundTripper) (*JWTHandler, error) {
	if spec == nil {
		return nil, core.ErrCoreGenCreateComponent.WithStack(ErrNilSpec, map[string]any{"reason": "failed to create JWT handler"})
	}

	// Get private keys for signing JWTs.
	sk, err := SigningKeys(spec.PrivateKeys...)
	if err != nil {
		return nil, core.ErrCoreGenCreateComponent.WithStack(err, map[string]any{"reason": "failed to create JWT handler"})
	}
	s := map[string]*SigningKey{}
	for _, k := range sk {
		s[k.kid] = k
	}

	// Get public keys for validating JWTs.
	vk, err := SigningKeys(spec.PublicKeys...)
	if err != nil {
		return nil, core.ErrCoreGenCreateComponent.WithStack(err, map[string]any{"reason": "failed to create JWT handler"})
	}
	v := map[string]*SigningKey{}
	for _, k := range vk {
		v[k.kid] = k
	}

	if spec.JWKs == nil {
		spec.JWKs = map[string]string{}
	}

	if rt == nil {
		rt = http.DefaultTransport
	}

	return &JWTHandler{
		signingKeys:    s,
		validatingKeys: v,
		jkus:           spec.JWKs,
		kids:           map[string][]string{},
		useJKU:         spec.UseJKU,
		rt:             rt,
	}, nil
}

// JWTHandler treat JWTs.
// JWTHandler sign, validate and parse JWTs.
type JWTHandler struct {

	// signingKeys are the private keys for public key algorithms.
	// Because mutex is not used for signingKeys,
	// do not write to this map on runtime to prevent conflict.
	signingKeys map[string]*SigningKey

	// mu protects validatingKeys.
	// validatingKeys will be updated when a key was not found.
	mu sync.RWMutex

	// validationKey is a set of signing keys.
	// Public keys for public key algorithms.
	validatingKeys map[string]*SigningKey

	// jkus is the set of jku, JWKs URI of providers.
	// Keys should be the issuer ID and the value should be the valid JWKs endpoint URI.
	jkus map[string]string

	// kids holds kid of keys obtained from JWKs endpoints.
	// The key is the JWKs URL and the values are kid obtained
	// from that endpoint.
	// All keys obtained from a endpoint should be removed from
	// this JWTHandler when refreshing keys.
	kids map[string][]string

	// useJKU is the flag to use JWKs endpoint set in the "jku" header.
	// If true, this handler tries to get JWK Set from the JWKs endpoint.
	useJKU bool

	// rt is the round tripper used for getting keys from JWKs endpoints.
	rt http.RoundTripper
}

// TokenWithClaims returns a token with signing key.
// This method returns an error when there is no signing key
// registered in this handler.
// Singing key is selected randomly from the registered keys.
func (h *JWTHandler) TokenWithClaims(claims jwt.Claims) (*jwt.Token, error) {
	// Use the key obtained first for sining the claims.
	// Note that the oder of map is uncertain.
	for _, v := range h.signingKeys {
		t := &jwt.Token{
			Method: v.method,
			Header: v.header,
			Claims: claims,
		}
		return t, nil
	}
	return nil, ErrNoSigningKey
}

// SignedString returns a signed string of the given token.
// This method panics when nil token was given.
func (h *JWTHandler) SignedString(token *jwt.Token) (string, error) {
	if token == nil || token.Header == nil {
		return "", ErrNilToken
	}
	kid, ok := token.Header["kid"].(string)
	if !ok {
		return "", ErrNoKid
	}
	key, ok := h.signingKeys[kid]
	if !ok {
		return "", ErrNoSigningKey
	}
	return token.SignedString(key.key)
}

// ParseWithClaims parse claims from token string.
func (h *JWTHandler) ParseWithClaims(token string, claims jwt.Claims, options ...jwt.ParserOption) (*jwt.Token, error) {
	return jwt.ParseWithClaims(token, claims, h.keyFunc, options...)
}

// ValidMapClaims returns jwt.MapClaims if the token was successfully validated.
// This method returns nil claims and non-nil error if the token was invalid.
func (h *JWTHandler) ValidMapClaims(token string, options ...jwt.ParserOption) (jwt.MapClaims, error) {
	c := jwt.MapClaims{}
	_, err := h.ParseWithClaims(token, &c, options...)
	if err != nil {
		return nil, app.ErrAppAuthnParseWithClaims.WithStack(err, map[string]any{"jwt": token})
	}
	return c, nil
}

// keyFunc is the key finding function.
//   - ErrNoKid: if "kid" not found in the header
func (h *JWTHandler) keyFunc(t *jwt.Token) (any, error) {
	kid, ok := t.Header["kid"].(string)
	if !ok {
		return nil, ErrNoKid
	}

	h.mu.RLock()
	key, ok := h.validatingKeys[kid]
	h.mu.RUnlock()
	if !ok {
		// Key not found in this handler.
		// Try to get keys from JWKs endpoints if possible.
		key, ok = h.refreshValidatingKeys(kid, t)
		if !ok {
			return nil, ErrKeyNotFound
		}
	}

	// key.method can be nil when the public key was fetched from a JWK set endpoints.
	// "alg" in the JWK set is optional.
	// RFC 7517 JSON Web Key (JWK) - 4.4.  "alg" (Algorithm) Parameter
	if key.method != nil {
		// "alg" header must be present.
		// RFC 7515 JSON Web Signature (JWS) - 4.1.1. "alg" (Algorithm) Header Parameter
		alg, ok := t.Header["alg"].(string)
		if !ok {
			return nil, ErrNoAlg
		}
		if alg != key.method.Alg() {
			return nil, ErrWrongAlg
		}
	}

	return key.key, nil
}

// refreshValidatingKeys refreshes validating keys.
func (h *JWTHandler) refreshValidatingKeys(kid string, t *jwt.Token) (*SigningKey, bool) {
	// jku is JWK set URI.
	var jku string

	if iss, err := t.Claims.GetIssuer(); err == nil {
		jku = h.jkus[iss]
	}

	if jku == "" && h.useJKU {
		jku, _ = t.Header["jku"].(string)
	}

	if jku == "" {
		// No JWKs URL to get validating keys.
		return nil, false
	}

	keys, err := fetchPublicKeys(h.rt, jku)
	if err != nil {
		return nil, false
	}

	// Obtain lock to update validatingKeys and kids.
	h.mu.Lock()
	defer h.mu.Unlock()

	// Keys are already refreshed while waiting for
	// fetching validating keys of getting a lock.
	if key, ok := h.validatingKeys[kid]; ok {
		return key, true
	}

	// Delete old keys obtained from this JWKs endpoints.
	for _, old := range h.kids[jku] {
		delete(h.validatingKeys, old)
	}

	kids := make([]string, 0, len(keys))
	var key *SigningKey
	for _, k := range keys {
		h.validatingKeys[k.kid] = k
		kids = append(kids, k.kid)
		if k.kid == kid {
			key = k
		}
	}

	h.kids[jku] = kids

	return key, key != nil
}

// fetchPublicKeys get public keys from given JWU, or JWKs URL.
// RoundTripper of the first argument must not nil.
func fetchPublicKeys(rt http.RoundTripper, jku string) ([]*SigningKey, error) {
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, jku, nil)
	if err != nil {
		return nil, app.ErrAppUtilGetJWKSet.WithStack(err, nil)
	}

	res, err := rt.RoundTrip(req)
	if err != nil {
		return nil, app.ErrAppUtilGetJWKSet.WithStack(err, nil)
	}
	defer res.Body.Close()

	set, err := jwk.ParseReader(res.Body)
	if err != nil {
		return nil, app.ErrAppUtilGetJWKSet.WithStack(err, nil)
	}

	var keys []*SigningKey
	for i := 0; i < set.Len(); i++ {
		k, _ := set.Key(i)

		pk, err := k.PublicKey()
		if err != nil {
			// It seems to be better to ignore invalid keys.
			// This things should be logged.
			continue
		}

		var raw any
		if err := pk.Raw(&raw); err != nil {
			// It seems to be better to ignore invalid keys.
			// This things should be logged.
			continue
		}

		key := &SigningKey{
			kid: k.KeyID(),
			header: map[string]any{
				"alg": k.Algorithm().String(), // "alg" can be empty string.
				"kid": k.KeyID(),
				"typ": "JWT",
			},
			// method can be nil because the "alg" is an optional field.
			// RFC 7517 JSON Web Key (JWK) - 4.4.  "alg" (Algorithm) Parameter
			method: SigningMethods[Algorithm(k.Algorithm().String())],
			key:    raw,
		}

		keys = append(keys, key)
	}

	return keys, nil
}
