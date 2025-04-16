// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package timeout

import (
	"context"
	"net/http"
	"slices"
	"time"

	"github.com/aileron-gateway/aileron-gateway/app"
	"github.com/aileron-gateway/aileron-gateway/core"
	"github.com/aileron-gateway/aileron-gateway/kernel/txtutil"
	utilhttp "github.com/aileron-gateway/aileron-gateway/util/http"
)

// apiTimeout applies timeout for requests
// based on the HTTP method and path.
type apiTimeout struct {
	// methods is the list of http methods
	// to apply this timeout.
	// Empty is equal to all methods.
	methods []string
	// paths is the path matcher to apply timeout.
	// paths must not be nil.
	paths txtutil.Matcher[string]
	// timeout is the timeout duration applied for
	// the request matched to the methods and paths.
	// No timeout is applied when this value is set to 0.
	timeout time.Duration
}

func (t *apiTimeout) duration(r *http.Request) (time.Duration, bool) {
	if len(t.methods) > 0 && !slices.Contains(t.methods, r.Method) {
		return 0, false
	}
	if !t.paths.Match(r.URL.Path) {
		return 0, false
	}
	return t.timeout, true
}

// timeout applies request timeout for all requests.
// This implements core.Middleware interface.
type timeout struct {
	eh core.ErrorHandler

	// defaultTimeout is the timeout duration applied by default.
	// This timeout is applied to all of the requests which were not
	// handled by the apiTimeout listed in the apiTimeouts variable.
	// No timeout is applied when this value is set to 0.
	defaultTimeout time.Duration

	// apiTimeouts is the list of apiTimeout to apply timeout
	// for requests with different methods and paths.
	apiTimeouts []*apiTimeout
}

func (m *timeout) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Determine timeout duration.
		// Apply default timeout if not matched to any apiTimeouts.
		to := m.defaultTimeout
		for _, at := range m.apiTimeouts {
			if d, ok := at.duration(r); ok {
				to = d // Update timeout duration
				break
			}
		}

		ctx := r.Context()

		// Do not apply timeout when the timeout duration is 0.
		if to > 0 {
			c, cancel := context.WithTimeout(ctx, to)
			ctx = c
			defer cancel()

			ww := utilhttp.WrapWriter(w)
			w = ww
			defer func() {
				if ctx.Err() == context.DeadlineExceeded && !ww.Written() {
					err := app.ErrAppMiddleAPITimeout.WithoutStack(ctx.Err(), nil)
					m.eh.ServeHTTPError(w, r, utilhttp.NewHTTPError(err, http.StatusGatewayTimeout))
				}
			}()
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
