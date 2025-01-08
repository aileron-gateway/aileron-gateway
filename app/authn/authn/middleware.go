package authn

import (
	"net/http"

	"github.com/aileron-gateway/aileron-gateway/app"
	"github.com/aileron-gateway/aileron-gateway/core"
	utilhttp "github.com/aileron-gateway/aileron-gateway/util/http"
)

// authn is the middleware that calls authentication handlers
// one by one.
// This implements core.Middleware interface.
type authn struct {
	eh core.ErrorHandler

	handlers []app.AuthenticationHandler
}

func (m *authn) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Yeld authentication for authentication handlers.
		authenticated := false
		for _, h := range m.handlers {
			newReq, result, shouldReturn, err := h.ServeAuthn(w, r)
			if newReq != nil {
				r = newReq // Update request.
			}
			if err != nil {
				m.eh.ServeHTTPError(w, r, err) // Authentication failed.
				return
			}
			if shouldReturn {
				return
			}
			if result == app.AuthSucceeded {
				authenticated = true // Authentication succeeded.
				break
			}
		}

		// Authentication failed.
		if !authenticated {
			m.eh.ServeHTTPError(w, r, utilhttp.ErrUnauthorized)
			return
		}

		// Authentication succeeded.
		next.ServeHTTP(w, r)
	})
}
