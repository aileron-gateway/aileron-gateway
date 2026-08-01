// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package httpserver

import (
	"net/http"
	"path"
	"slices"
	"strings"

	v1 "github.com/aileron-gateway/aileron-gateway/apis/core/v1"
	"github.com/aileron-gateway/aileron-gateway/core"
	"github.com/aileron-gateway/aileron-gateway/kernel/api"
	utilhttp "github.com/aileron-gateway/aileron-gateway/util/http"
)

// notFoundHandler returns a not found handler.
// The returned handler will panic if the nil error handler was given.
func notFoundHandler(eh core.ErrorHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pattern := r.Method + " " + r.Host + r.URL.Path
		err := core.ErrCoreServerNotFound.WithoutStack(nil, map[string]any{"pattern": pattern})
		eh.ServeHTTPError(w, r, utilhttp.NewHTTPError(err, http.StatusNotFound))
	})
}

type Mux interface {
	http.Handler
	Handle(pattern string, h http.Handler)
}

// registerHandlers register virtual host handlers to the given mux..
// The function panics if the mux is nil.
func registerHandlers(a api.API[*api.Request, *api.Response], mux Mux, specs []*v1.VirtualHostSpec, notFound http.Handler) (err error) {
	defer func() {
		if err != nil {
			return
		}
		// Registering invalid path pattern to the serve mux can panic.
		if e, ok := recover().(error); ok {
			reason := "failed to register handler to HTTP router"
			err = core.ErrCoreGenCreateComponent.WithStack(e, map[string]any{"reason": reason})
		}
	}()

	for _, vhSpec := range specs {
		middleware, err := api.ReferTypedObjects[core.Middleware](a, vhSpec.Middleware...)
		if err != nil {
			return core.ErrCoreGenCreateComponent.WithStack(err, map[string]any{"reason": "failed to get middleware"})
		}

		for _, hSpec := range vhSpec.Handlers {
			methods, paths, handler, err := utilhttp.Handler(a, hSpec)
			if err != nil {
				return core.ErrCoreGenCreateComponent.WithStack(err, map[string]any{"reason": "failed to create handler"})
			}
			if len(paths) == 0 {
				paths = append(paths, hSpec.Pattern)
			}
			for i, p := range paths {
				paths[i] = path.Clean("/" + vhSpec.Pattern + p) // path.Clean removes trailing slash.
				if strings.HasSuffix(p, "/") && paths[i] != "/" {
					paths[i] = paths[i] + "/" // Append trailing slash if necessary.
				}
			}
			if len(vhSpec.Methods) > 0 {
				if len(methods) == 0 {
					methods = utilhttp.Methods(vhSpec.Methods)
				} else {
					methods = intersectionString(methods, utilhttp.Methods(vhSpec.Methods))
				}
			}
			for _, pattern := range generatePatterns(vhSpec.Hosts, paths) {
				h := utilhttp.MiddlewareChain(middleware, handler)
				mux.Handle(pattern, &methodCheckHandler{
					handler:          h,
					methodNotAllowed: notFound,
					allowMethods:     methods,
					allowAll:         len(methods) == 0,
				})
			}
		}
		for _, pattern := range generatePatterns(vhSpec.Hosts, nil) {
			registerNotFound(mux, pattern, notFound)
		}
	}

	registerNotFound(mux, "/", notFound)
	return nil
}

func registerNotFound(mux Mux, pattern string, h http.Handler) {
	defer func() { recover() }()
	mux.Handle(pattern, h)
}

// IntersectionString returns intersection of the two given set.
//
// example:
//
//	s1 := []string{"A", "B"}
//	s2 := []string{"B", "C"}
//	s := util.IntersectionString(s1,s2) // []string{"B"}
func intersectionString(set1, set2 []string) []string {
	if len(set1) == 0 || len(set2) == 0 {
		return nil
	}
	s := []string{}
	for _, v := range set1 {
		if slices.Contains(set2, v) {
			s = append(s, v)
		}
	}

	slices.Sort(s)                     // slices.Compact require sorts.
	s = slices.Clip(slices.Compact(s)) // Remove duplicates.
	return s
}

// generatePatterns returns pattern for serve mux.
func generatePatterns(hosts []string, paths []string) []string {
	if len(hosts) == 0 {
		hosts = []string{""}
	}
	if len(paths) == 0 {
		paths = []string{"/"}
	}
	patterns := make([]string, 0, len(hosts)*len(paths))
	for _, h := range hosts {
		for _, p := range paths {
			p = "/" + strings.TrimPrefix(p, "/")
			patterns = append(patterns, h+p) // "[HOST]/[PATH]" https://pkg.go.dev/net/http#ServeMux
		}
	}
	slices.Sort(patterns)           // slices.Compact require sorts.
	return slices.Compact(patterns) // Remove duplicates.
}

type methodCheckHandler struct {
	handler          http.Handler
	methodNotAllowed http.Handler
	allowMethods     []string // Allowed methods list.
	allowAll         bool     // Allow all HTTP methods.
}

func (h methodCheckHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if h.allowAll || slices.Contains(h.allowMethods, r.Method) {
		h.handler.ServeHTTP(w, r)
		return
	}
	h.methodNotAllowed.ServeHTTP(w, r)
}
