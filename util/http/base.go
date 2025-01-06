package http

import (
	"net/http"
	"slices"

	v1 "github.com/aileron-gateway/aileron-gateway/apis/core/v1"
)

// HTTPMethods is the all mapping of HTTP methods.
// Currently following methods are contained.
//   - Get
//   - Head
//   - Post
//   - Put
//   - Patch
//   - Delete
//   - Connect
//   - Options
//   - Trace
var HTTPMethods = map[v1.HTTPMethod]string{
	v1.HTTPMethod_GET:     http.MethodGet,
	v1.HTTPMethod_HEAD:    http.MethodHead,
	v1.HTTPMethod_POST:    http.MethodPost,
	v1.HTTPMethod_PUT:     http.MethodPut,
	v1.HTTPMethod_PATCH:   http.MethodPatch,
	v1.HTTPMethod_DELETE:  http.MethodDelete,
	v1.HTTPMethod_CONNECT: http.MethodConnect,
	v1.HTTPMethod_OPTIONS: http.MethodOptions,
	v1.HTTPMethod_TRACE:   http.MethodTrace,
}

// Methods returns a list of http methods in string.
// Unknown HTTP methods are ignored.
// Currently following methods are supported.
//   - Get
//   - Head
//   - Post
//   - Put
//   - Patch
//   - Delete
//   - Connect
//   - Options
//   - Trace
func Methods(methods []v1.HTTPMethod) []string {
	if len(methods) == 0 {
		return nil
	}
	ms := make([]string, 0, len(methods))
	for _, method := range methods {
		if m, ok := HTTPMethods[method]; ok {
			ms = append(ms, m)
		}
	}
	slices.Sort(ms)
	return slices.Compact(slices.Clip(ms))
}

// HandlerBase is the base struct for a http handler.
type HandlerBase struct {
	// AcceptPattern is the url path patterns this handler can accepts.
	// Path patterns that http.ServeMux can handle are allowed.
	// Do not contain hostname or methods.
	// 	- https://pkg.go.dev/net/http#ServeMux
	AcceptPatterns []string

	// AcceptMethods is the http  methods this handler can accept.
	// Methods defined in the http package are allowed.
	// 	- https://pkg.go.dev/net/http#pkg-constants
	AcceptMethods []string
}

// Patterns returns the http path pattern
// that this handler can handle.
func (h *HandlerBase) Patterns() []string {
	return h.AcceptPatterns
}

// Methods returns the http methods
// this handler can handle.
func (h *HandlerBase) Methods() []string {
	return h.AcceptMethods
}
