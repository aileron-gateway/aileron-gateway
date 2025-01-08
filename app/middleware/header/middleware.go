package header

import (
	"mime"
	"net/http"
	"slices"

	"github.com/aileron-gateway/aileron-gateway/core"
	"github.com/aileron-gateway/aileron-gateway/kernel/txtutil"
	httputil "github.com/aileron-gateway/aileron-gateway/util/http"
)

// wrappedWriter wraps ResponseWriter and apply header policies to
// response headers before response body written.
type wrappedWriter struct {
	http.ResponseWriter
	// applied is set to true when the policy have been applied
	// to the response header.
	applied bool
	// p is the header policy that should be applied to the
	// response headers.
	p *policy
}

// Unwrap returns internal ResponseWriter.
// Unwrap method is conventionally required when wrapping
// the other interface or struct.
func (w *wrappedWriter) Unwrap() http.ResponseWriter {
	return w.ResponseWriter
}

// WriteHeader is the WriteHeader method of http.ResponseWriter interface.
func (w *wrappedWriter) WriteHeader(statusCode int) {
	// Response headers must be fixed before writing status code.
	// Apply policies if it have not been applied yet.
	if !w.applied {
		w.applied = true
		w.p.apply(w.Header())
	}
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *wrappedWriter) Write(b []byte) (int, error) {
	// Apply policies if it have not been applied yet.
	if !w.applied {
		w.applied = true
		w.p.apply(w.Header())
	}
	return w.ResponseWriter.Write(b)
}

// policy is the request/response header policy.
type policy struct {
	allows  []string
	removes []string
	add     map[string]string
	set     map[string]string
	repls   map[string]txtutil.ReplaceFunc[string]
}

func (p *policy) apply(h http.Header) {
	if len(p.allows) > 0 {
		for k := range h {
			if !slices.Contains(p.allows, k) {
				h.Del(k) // Remove not allowed headers.
			}
		}
	}
	for _, k := range p.removes {
		// Set nil instead of removing it in the
		// case for "Date" response header.
		// See https://pkg.go.dev/net/http#ResponseWriter.Header
		h[k] = nil
	}
	for k, v := range p.add {
		h.Add(k, v)
	}
	for k, v := range p.set {
		h.Set(k, v)
	}
	for name, replFunc := range p.repls {
		values := h.Values(name)
		for i, v := range values {
			values[i] = replFunc(v)
		}
	}
}

// headerPolicy is the middleware that applies header policies
// to requests and response headers.
// This implements core.Middleware interface.
type headerPolicy struct {
	eh core.ErrorHandler

	// allowedMIMEs is the allowed list of media types.
	// Values are evaluated by exactly matching.
	// This is effective only for requests
	// and not for response.
	// If empty, mimes are not checked.
	allowedMIMEs []string

	// maxContentLength is the max body size to allow.
	// This is effective only for requests
	// and not for response.
	// If zero, content length are not checked.
	maxContentLength int64

	// reqPolicy is the header policy applied to requests.
	// reqPolicy must not be nil.
	reqPolicy *policy

	// resPolicy is the header policy applied to responses.
	// resPolicy must not be nil.
	resPolicy *policy
}

func (m *headerPolicy) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Apply MIME whitelist.
		if len(m.allowedMIMEs) > 0 {
			mt, _, _ := mime.ParseMediaType(r.Header.Get("Content-Type"))
			if !slices.Contains(m.allowedMIMEs, mt) {
				m.eh.ServeHTTPError(w, r, httputil.ErrUnsupportedMediaType)
				return
			}
		}

		// Restrict maximum content length.
		if m.maxContentLength > 0 {
			if r.ContentLength == -1 { // Unknown sized body or streaming body.
				m.eh.ServeHTTPError(w, r, httputil.ErrLengthRequired)
				return
			} else if r.ContentLength > m.maxContentLength {
				m.eh.ServeHTTPError(w, r, httputil.ErrRequestEntityTooLarge)
				return
			}
		}

		// Apply policy to request headers.
		if m.reqPolicy != nil {
			m.reqPolicy.apply(r.Header)
		}

		// Apply policy to response headers.
		// Response header must be rewritten before
		// the StatusCode or the Body is written.
		if m.resPolicy != nil {
			w = &wrappedWriter{
				ResponseWriter: w,
				p:              m.resPolicy,
			}
		}

		next.ServeHTTP(w, r)
	})
}
