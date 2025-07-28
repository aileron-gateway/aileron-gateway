// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package throttle

import (
	"net/http"
	"slices"

	"github.com/aileron-gateway/aileron-gateway/app"
	"github.com/aileron-gateway/aileron-gateway/core"
	utilhttp "github.com/aileron-gateway/aileron-gateway/util/http"
)

// throttle throttles requests.
// This implements core.Middleware interface.
type throttle struct {
	eh         core.ErrorHandler
	throttlers []*apiThrottler
}

func (m *throttle) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for _, t := range m.throttlers {
			if len(t.methods) > 0 && !slices.Contains(t.methods, r.Method) {
				continue
			}
			if !t.paths.Match(r.URL.Path) {
				continue
			}

			accepted, release := t.accept(r.Context())
			if accepted {
				defer release()
				break
			}

			err := app.ErrAppMiddleThrottle.WithoutStack(nil, nil)
			m.eh.ServeHTTPError(w, r, utilhttp.NewHTTPError(err, http.StatusTooManyRequests))
			return
		}

		next.ServeHTTP(w, r)
	})
}
