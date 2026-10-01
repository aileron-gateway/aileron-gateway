// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package key

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/aileron-gateway/aileron-gateway/app"
	"github.com/aileron-gateway/aileron-gateway/core"
	"github.com/aileron-gateway/aileron-gateway/internal/encoder"
	"github.com/aileron-gateway/aileron-gateway/internal/hash"
	"github.com/aileron-gateway/aileron-gateway/internal/kvs"
	utilhttp "github.com/aileron-gateway/aileron-gateway/util/http"
)

var (
	errUnregisteredKey = utilhttp.NewHTTPError(errors.New("non registered api key provided"), http.StatusForbidden)
	errInvalidKey      = utilhttp.NewHTTPError(errors.New("invalid api key provided"), http.StatusForbidden)
)

type handler struct {
	eh core.ErrorHandler

	keyHeaderName string

	// claimsKey is the key of the claims
	// to save authn info in the context.
	claimsKey string

	// keep is the flag to keep API key header.
	// If true, the Authorization header won't be
	// removed and be sent upstream services.
	keep bool

	// sore is the API key store.
	// API key and bounded info are obtained.
	store kvs.Commander[string, credential]

	encodeFunc encoder.EncodeToStringFunc

	hashFunc hash.HashFunc
	hmacFunc hash.HMACFunc
	hmacKey  []byte
}

func (h *handler) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		newReq, status, err := h.ServeAuthn(w, r)
		r = newReq
		if err != nil {
			h.eh.ServeHTTPError(w, r, err)
			return
		}
		if status&app.AuthSuccess > 0 {
			next.ServeHTTP(w, r)
			return
		}

		h.eh.ServeHTTPError(w, r, utilhttp.ErrUnauthorized)
	})
}

func (h *handler) ServeAuthn(w http.ResponseWriter, r *http.Request) (*http.Request, app.AuthStatus, error) {
	key := r.Header.Get(h.keyHeaderName)
	if key == "" {
		return r, app.AuthFail, utilhttp.ErrForbidden
	}

	givenSecret := []byte(key)
	if len(h.hmacKey) > 0 {
		givenSecret = h.hmacFunc([]byte(key), h.hmacKey)
	} else if h.hashFunc != nil {
		givenSecret = h.hashFunc([]byte(key))
	}

	var id string
	if h.encodeFunc != nil {
		id = h.encodeFunc(givenSecret)
	} else {
		id = string(givenSecret)
	}

	cred, err := h.store.Get(r.Context(), id)
	if err != nil {
		if err == kvs.Nil {
			return r, app.AuthFail, errUnregisteredKey
		}
		return r, app.AuthFail, utilhttp.NewHTTPError(err, http.StatusInternalServerError)
	}

	if string(cred.Secret()) != string(givenSecret) {
		return r, app.AuthFail, errInvalidKey
	}

	claims := &Claims{
		Method:   "APIKey",
		AuthTime: time.Now().Unix(),
		Key:      key,
		Attrs:    cred.Attributes(),
	}

	// Save claims in the context so it can be used for authorization.
	//nolint:staticcheck // SA1029: should not use built-in type string as key for value; define your own type to avoid collisions
	ctx := context.WithValue(r.Context(), h.claimsKey, claims)
	r = r.WithContext(ctx)

	if !h.keep {
		delete(r.Header, h.keyHeaderName)
	}

	return r, app.AuthSuccess, nil
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
	Key      string `json:"key" msgpack:"key"`
	Attrs    any    `json:"attrs" msgpack:"attrs"`
}
