// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package basic

import (
	"context"
	"net/http"
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

	// challenge is the basic authentication challenge
	// header value which means "WWW-Authenticate: <challenge>".
	challenge string

	store kvs.Commander[string, credential]

	passwd      []byte
	decryptFunc encrypt.DecryptFunc
	compareFunc encrypt.PasswordCompareFunc

	// preferError if true, returns an error
	// when authentication failed rather than asking
	// a new username and password.
	preferError bool
}

func (h *handler) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		newReq, status, err := h.ServeAuthn(w, r)
		r = newReq
		if err != nil {
			err = app.ErrAppAuthnAuthentication.WithoutStack(err, nil)
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

		err = app.ErrAppAuthnAuthentication.WithoutStack(nil, nil)
		h.eh.ServeHTTPError(w, r, utilhttp.NewHTTPError(err, http.StatusUnauthorized))
	})
}

func (h *handler) ServeAuthn(w http.ResponseWriter, r *http.Request) (*http.Request, app.AuthStatus, error) {
	un, pw, ok := r.BasicAuth()
	if !ok {
		w.Header().Set("WWW-Authenticate", h.challenge)
		w.WriteHeader(http.StatusUnauthorized)
		w.Write(nil)
		return r, app.AuthReturn, nil
	}

	cred, err := h.authenticate(r.Context(), un, pw)
	if err != nil {
		if h.preferError { // Fail first. Do not ask re-authentication.
			return r, app.AuthFail, app.ErrAppAuthnInvalidCredential.WithoutStack(err, map[string]any{"purpose": "Basic authentication"})
		}
		// Invalid credentials passed and re-require authentication.
		w.Header().Set("WWW-Authenticate", h.challenge)
		w.WriteHeader(http.StatusUnauthorized)
		w.Write(nil)
		return r, app.AuthReturn, nil
	}

	claims := &Claims{
		Method:   "Basic",
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

func (h *handler) authenticate(ctx context.Context, un, pw string) (credential, error) {
	cred, err := h.store.Get(ctx, un)
	if err != nil {
		return nil, err
	}

	secret := cred.Secret()
	if h.decryptFunc != nil {
		secret, err = h.decryptFunc(h.passwd, secret)
		if err != nil {
			return nil, err
		}
	}
	if err := h.compareFunc(secret, []byte(pw)); err != nil {
		return nil, err
	}

	return cred, nil
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
