package skipper

import (
	"net/http"
	"slices"

	"github.com/aileron-gateway/aileron-gateway/core"
	"github.com/aileron-gateway/aileron-gateway/kernel/log"
	"github.com/aileron-gateway/aileron-gateway/kernel/txtutil"
	utilhttp "github.com/aileron-gateway/aileron-gateway/util/http"
)

type skipper struct {
	// methods is the list of http methods to skip.
	// Empty is equal to all methods.
	methods []string
	// paths is the path matcher to skip.
	// paths must not be nil.
	paths txtutil.Matcher[string]
}

func (t *skipper) shouldSkip(r *http.Request) bool {
	if len(t.methods) > 0 && !slices.Contains(t.methods, r.Method) {
		return false
	}
	return t.paths.Match(r.URL.Path)
}

type skippable struct {
	skippers []*skipper

	ms []core.Middleware
	ts []core.Tripperware

	// name is the name of this resource
	// for debug logging.
	name string
	// lg is the logger for debug logging.
	lg log.Logger
}

func (s *skippable) Middleware(next http.Handler) http.Handler {
	whenSkip := next
	whenNotSkip := utilhttp.MiddlewareChain(s.ms, next)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for _, skipper := range s.skippers {
			if skipper.shouldSkip(r) {
				if s.lg.Enabled(log.LvDebug) {
					s.lg.Debug(r.Context(), s.name+" is skipped")
				}
				whenSkip.ServeHTTP(w, r)
				return
			}
		}

		whenNotSkip.ServeHTTP(w, r)
	})
}

func (s *skippable) Tripperware(next http.RoundTripper) http.RoundTripper {
	whenSkip := next
	whenNotSkip := utilhttp.TripperwareChain(s.ts, next)

	return core.RoundTripperFunc(func(r *http.Request) (*http.Response, error) {
		for _, skipper := range s.skippers {
			if skipper.shouldSkip(r) {
				if s.lg.Enabled(log.LvDebug) {
					s.lg.Debug(r.Context(), s.name+" is skipped")
				}
				return whenSkip.RoundTrip(r)
			}
		}

		return whenNotSkip.RoundTrip(r)
	})
}
