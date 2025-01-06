package main

import (
	"net/http"

	"github.com/aileron-gateway/aileron-gateway/core"
	"github.com/aileron-gateway/aileron-gateway/kernel/log"
)

var Plugin = plugin{}

type plugin struct {
	lg log.Logger
	eh core.ErrorHandler
}

// Init is the implements of goplugin.Initializer interface.
// Init is called only when it is implemented.
func (p *plugin) Init(lg log.Logger, eh core.ErrorHandler) error {
	p.lg = lg
	p.eh = eh
	return nil
}

// Middleware is the implements of core.Middleware interface.
// The plugin can be used as middleware by implementing this method.
func (p *plugin) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		p.lg.Info(r.Context(), "Plugin middleware called")

		w.Header().Add("foo", "bar")
		next.ServeHTTP(w, r)

	})
}

// Tripperware is the implements of core.Tripperware interface.
// The plugin can be used as tripperware by implementing this method.
func (p *plugin) Tripperware(next http.RoundTripper) http.RoundTripper {
	return core.RoundTripperFunc(func(r *http.Request) (w *http.Response, err error) {

		p.lg.Info(r.Context(), "Plugin tripperware called")

		r.Header.Add("alice", "bob")
		return next.RoundTrip(r)

	})
}

// ServeHTTP is the implements of http.Handler interface.
// The plugin can be used as a http handler by implementing this method.
func (p *plugin) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	p.lg.Info(r.Context(), "Plugin handler called")

	w.Header().Set("hello", "plugin")
	w.Write([]byte("GoPlugin Handler !!"))

}

// Patterns returns path pattern for registering this
// plugin as a handler for multiplexers.
// If method not implemented or empty slice were returned,
// the handler is registered with a path "/".
func (p *plugin) Patterns() []string {
	return []string{"/goplugin"}
}

// Methods returns method names for registering this
// plugin as a handler for multiplexers.
// If method not implemented or empty slice were returned,
// all methods are accepted.
func (p *plugin) Methods() []string {
	return []string{http.MethodGet}
}
