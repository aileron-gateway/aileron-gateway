// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package digest

import (
	"context"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/aileron-gateway/aileron-gateway/app"
	"github.com/aileron-gateway/aileron-gateway/core"
	"github.com/aileron-gateway/aileron-gateway/internal/encrypt"
	"github.com/aileron-gateway/aileron-gateway/internal/kvs"
	utilhttp "github.com/aileron-gateway/aileron-gateway/util/http"
)

type handler struct {
	eh core.ErrorHandler

	// claimsKey is the key of the claims
	// to save authn info in the context.
	claimsKey string

	// keep is the flag to keep Authorization header.
	// If true, the Authorization header won't be
	// removed and be sent upstream services.
	keep bool

	realm     string
	algorithm string
	hashFunc  func(string) string

	store kvs.Commander[string, credential]

	passwd      []byte
	decryptFunc encrypt.DecryptFunc
	compareFunc encrypt.PasswordCompareFunc
}

func (h *handler) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		newReq, status, err := h.ServeAuthn(w, r)
		r = newReq
		if err != nil {
			h.eh.ServeHTTPError(w, r, utilhttp.NewHTTPError(err, http.StatusUnauthorized))
			return
		}
		if status&app.AuthReturn > 0 {
			return
		}
		if status&app.AuthSuccess > 0 || status&app.AuthSkip > 0 {
			next.ServeHTTP(w, r)
			return
		}

		h.eh.ServeHTTPError(w, r, utilhttp.ErrUnauthorized)
	})
}

func (h *handler) ServeAuthn(w http.ResponseWriter, r *http.Request) (*http.Request, app.AuthStatus, error) {
	un, cred, err := h.authenticate(r.Context(), r.Header.Get("Authorization"), r.Method)
	if err != nil {
		b := make([]byte, 30)
		io.ReadFull(rand.Reader, b)
		nonce := base64.StdEncoding.EncodeToString(b)
		challenge := `Digest algorithm=` + h.algorithm + `,qop="auth",realm="` + h.realm + `",nonce="` + nonce + `",charset=UTF-8`
		w.Header().Set("WWW-Authenticate", challenge)
		w.WriteHeader(http.StatusUnauthorized)
		w.Write(nil)
		return r, app.AuthReturn, nil
	}

	claims := &Claims{
		Method:   "Digest",
		AuthTime: time.Now().Unix(),
		Name:     un,
		Attrs:    cred.Attributes(),
	}

	// Save claims in the context so it can be used for authorization.
	//nolint:staticcheck // SA1029: should not use built-in type string as key for value; define your own type to avoid collisions
	ctx := context.WithValue(r.Context(), h.claimsKey, claims)
	r = r.WithContext(ctx)

	if !h.keep {
		delete(r.Header, "Authorization")
	}

	return r, app.AuthSuccess, nil
}

func (h *handler) authenticate(ctx context.Context, auth, method string) (username string, cred credential, err error) {
	if auth == "" {
		return "", nil, errors.New("")
	}

	const prefix = "Digest "
	// Case insensitive prefix match. See Issue 22736.
	// https://github.com/golang/go/issues/22736
	if len(auth) < len(prefix) || !strings.EqualFold(auth[:len(prefix)], prefix) {
		return "", nil, errors.New("")
	}
	auth = auth[len(prefix):]

	var secret []byte
	var uri, nonce, response, nc, cnonce, qop string

	tokens := strings.Split(auth, ",")
	for i := 0; i < len(tokens); i++ {
		key, value, ok := strings.Cut(strings.TrimSpace(tokens[i]), "=")
		if !ok {
			continue
		}
		value = strings.Trim(value, `"`)
		switch key {
		case "uri":
			uri = value
		case "nonce":
			nonce = value
		case "response":
			response = value
		case "nc":
			nc = value
		case "cnonce":
			cnonce = value
		case "qop":
			qop = value
		case "username":
			username = value
			cred, err = h.store.Get(ctx, username)
			if err != nil {
				return "", nil, err
			}
			secret = cred.Secret()
			if h.decryptFunc != nil {
				secret, err = h.decryptFunc(h.passwd, secret)
				if err != nil {
					return "", nil, err
				}
			}
		}
	}

	// a1 is the parameter A1.
	// "<algorithm>-sess" type is not supported here.
	// RFC 7616: HTTP Digest Access Authentication 3.4.2 A1
	// https://www.rfc-editor.org/rfc/rfc7616.html#section-3.4.2
	// Should not use the client provided realm. Use the correct realm we have.
	a1 := h.hashFunc(username + ":" + h.realm + ":" + string(secret))

	// a2 is the parameter A2.
	// "auth-int" type of qop is not supported here.
	// RFC 7616: HTTP Digest Access Authentication 3.4.3 A2
	// https://www.rfc-editor.org/rfc/rfc7616.html#section-3.4.3
	a2 := h.hashFunc(method + ":" + uri)

	// kd is the keyed digest.
	// RFC 7616: HTTP Digest Access Authentication 3.4.1 Response
	// https://www.rfc-editor.org/rfc/rfc7616.html#section-3.4.1
	kd := h.hashFunc(a1 + ":" + nonce + ":" + nc + ":" + cnonce + ":" + qop + ":" + a2)

	if response != kd {
		return username, nil, ErrNotMatch
	}

	return username, cred, nil
}

// hashFuncMD5 is the MD5 hash calculation function.
// MD5 is one of the supported algorithm for digest authentication.
// algorithm:MD5.
// Hash algorithms are defined at here.
//   - https://www.rfc-editor.org/rfc/rfc7616.html#section-6.1
func hashFuncMD5(s string) string {
	h := md5.Sum([]byte(s))
	return hex.EncodeToString(h[:])
}

// hashFuncSHA256 is the SHA256 hash calculation function.
// SHA256 is one of the supported algorithm for digest authentication.
// algorithm:SHA-256.
// Hash algorithms are defined at here.
//   - https://www.rfc-editor.org/rfc/rfc7616.html#section-6.1
func hashFuncSHA256(s string) string {
	h := sha256.Sum256([]byte(s))
	return hex.EncodeToString(h[:])
}

// hashFuncSHA512 is the SHA-512-256 hash calculation function.
// SHA-512-256 is one of the supported algorithm for digest authentication.
// algorithm:SHA-512-256
// Hash algorithms are defined at here.
//   - https://www.rfc-editor.org/rfc/rfc7616.html#section-6.1
func hashFuncSHA512(s string) string {
	h := sha512.Sum512_256([]byte(s))
	return hex.EncodeToString(h[:])
}

type credential interface {
	Secret() []byte
	Attributes() any
}

type defaultCredential struct {
	secret []byte
	attrs  any
}

func (c *defaultCredential) Secret() []byte {
	return c.secret
}

func (c *defaultCredential) Attributes() any {
	return c.attrs
}

// Claims is the claims for
// successfully authenticated entity.
// Claims are saved in the request context
// and can be used for authorization.
type Claims struct {
	Method   string `json:"method" msgpack:"method"`
	AuthTime int64  `json:"auth_time" msgpack:"auth_time"`
	Name     string `json:"name" msgpack:"name"`
	Attrs    any    `json:"attrs" msgpack:"attrs"`
}
