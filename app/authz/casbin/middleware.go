// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package casbin

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/aileron-gateway/aileron-gateway/app"
	"github.com/aileron-gateway/aileron-gateway/core"
	"github.com/aileron-gateway/aileron-gateway/kernel/log"
	utilhttp "github.com/aileron-gateway/aileron-gateway/util/http"
	"github.com/casbin/casbin/v3"
)

// authzClaims holds request information
// that can be used for authorization.
type authzClaims struct {
	Auth   any         `json:"auth" yaml:"auth" msgpack:"auth"`
	Host   string      `json:"host" yaml:"host" msgpack:"host"`
	Remote string      `json:"remote" yaml:"remote" msgpack:"remote"`
	Method string      `json:"method" yaml:"method" msgpack:"method"`
	API    string      `json:"api" yaml:"api" msgpack:"api"`
	Query  url.Values  `json:"query" yaml:"query" msgpack:"query"`
	Header http.Header `json:"header" yaml:"header" msgpack:"header"`
}

// authz do authorization using casbin authorization framework.
// This implements app.AuthorizationHandler interface.
//
// References:
//   - https://casbin.org/
//   - https://casbin.org/docs/supported-models
type authz struct {
	lg log.Logger
	w  io.Writer

	eh core.ErrorHandler

	// enforcers are casbin authorization enforcers.
	enforcers []casbin.IEnforcer

	// key is the context key to extract subject's information.
	key string
	// extraKeys is the list of context key name
	// which value is added to the request_definition.
	// This will be used if this gateway were extended as a library.
	extraKeys []string

	// explain is the flag to log authorization information.
	// This should be false in production environment
	// because of the performance consideration.
	explain bool
}

func (m *authz) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims := &authzClaims{
			Auth:   r.Context().Value(m.key),
			Host:   r.Host,
			Remote: r.RemoteAddr,
			Method: r.Method,
			API:    r.URL.Path,
			Query:  r.URL.Query(),
			Header: r.Header,
		}

		var vals = []any{claims, r.URL.Path, r.Method}
		for _, k := range m.extraKeys {
			vals = append(vals, r.Context().Value(k))
		}

		for _, enf := range m.enforcers {
			var ok bool
			var err error

			if m.explain {
				var msgs [][]string
				ok, msgs, err = enf.EnforceEx(vals...)
				for _, msg := range msgs {
					m.w.Write([]byte(strings.Join(msg, ", ")))
				}
			} else {
				ok, err = enf.Enforce(vals...)
			}

			// Unauthorized
			if err != nil {
				if m.lg.Enabled(log.LvDebug) {
					msg := fmt.Sprintf("forbidden by Casbin for method=%s path=%s", r.Method, r.URL.Path)
					err := app.ErrAppAuthzAuthorization.WithoutStack(err, nil)
					m.lg.Debug(r.Context(), msg, err.Name(), err.Map())
				}
				continue
			}

			// Unauthorized
			if !ok {
				if m.lg.Enabled(log.LvDebug) {
					msg := fmt.Sprintf("forbidden by Casbin for method=%s path=%s", r.Method, r.URL.Path)
					err := app.ErrAppAuthzAuthorization.WithoutStack(nil, nil)
					m.lg.Debug(r.Context(), msg, err.Name(), err.Map())
				}
				continue
			}

			// Authorized.
			next.ServeHTTP(w, r)
			return
		}

		// Forbidden.
		m.eh.ServeHTTPError(w, r, utilhttp.ErrForbidden)
	})
}
