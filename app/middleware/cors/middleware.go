// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package cors

import (
	"net/http"

	"github.com/aileron-gateway/aileron-gateway/app"
	"github.com/aileron-gateway/aileron-gateway/core"
	utilhttp "github.com/aileron-gateway/aileron-gateway/util/http"
)

// cors is the middleware which applies CORS(Cross-Origin Resource Sharing).
// This implements core.Middleware interface.
//
// References on CORS specification.
//   - https://datatracker.ietf.org/doc/rfc6454/ [The Web Origin Concept]
//   - https://fetch.spec.whatwg.org/ [Fetch]
//   - https://docs.w3cub.com/http/headers/cross-origin-embedder-policy
//   - https://docs.w3cub.com/http/headers/cross-origin-opener-policy
//   - https://docs.w3cub.com/http/headers/cross-origin-resource-policy
//   - https://docs.w3cub.com/http/headers/access-control-allow-credentials
//   - https://docs.w3cub.com/http/headers/access-control-allow-headers
//   - https://docs.w3cub.com/http/headers/access-control-allow-methods
//   - https://docs.w3cub.com/http/headers/access-control-allow-origin
//   - https://docs.w3cub.com/http/headers/access-control-expose-headers
//   - https://docs.w3cub.com/http/headers/access-control-max-age
//   - https://docs.w3cub.com/http/headers/access-control-request-headers
//   - https://docs.w3cub.com/http/headers/access-control-request-method
type cors struct {
	eh core.ErrorHandler
	// policies is the list of CORS policy
	// which applied to the all incoming requests.
	// Requests are forbidden when they are
	// not allowed by at least one policy in this list.
	policy *corsPolicy
}

func (m *cors) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// "Access-Control-Request-Method" is only included in the preflight requests
		// and not in the actual requests.
		// We check this header to identify if this is a preflight request or an actual request.
		// https://fetch.spec.whatwg.org/#http-requests
		if r.Method == http.MethodOptions && r.Header.Get("Access-Control-Request-Method") != "" {
			status := m.policy.handlePreflight(w, r)
			w.WriteHeader(status)
			w.Write(nil)
			return
		}

		status := m.policy.handleActualRequest(w, r)
		if status != http.StatusOK {
			err := app.ErrAppMiddleCORSForbidden.WithoutStack(nil, nil)
			m.eh.ServeHTTPError(w, r, utilhttp.NewHTTPError(err, http.StatusForbidden))
			return
		}

		next.ServeHTTP(w, r)
	})
}
