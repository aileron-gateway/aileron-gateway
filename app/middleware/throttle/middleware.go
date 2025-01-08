package throttle

import (
	"net/http"
	"slices"

	"github.com/aileron-gateway/aileron-gateway/core"
	"github.com/aileron-gateway/aileron-gateway/kernel/txtutil"
	utilhttp "github.com/aileron-gateway/aileron-gateway/util/http"
)

// apiThrottler applies throttling for requests
// that matched to the method and paths.
type apiThrottler struct {
	throttler
	methods []string
	paths   txtutil.Matcher[string]
}

// throttle throttles requests.
// This implements core.Middleware interface.
type throttle struct {
	eh core.ErrorHandler

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

			m.eh.ServeHTTPError(w, r, utilhttp.ErrTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}
