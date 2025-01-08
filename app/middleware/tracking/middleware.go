package tracking

import (
	"net/http"

	"github.com/aileron-gateway/aileron-gateway/app"
	"github.com/aileron-gateway/aileron-gateway/core"
	"github.com/aileron-gateway/aileron-gateway/kernel/log"
	"github.com/aileron-gateway/aileron-gateway/kernel/uid"
	utilhttp "github.com/aileron-gateway/aileron-gateway/util/http"
)

// NewRequestID returns a new request ID.
// Returned ID should not be empty when the error is nil.
// Non nil error will result in an internal server error.
var NewRequestID func() (string, error)

// NewTraceID returns a new trace ID.
// Returned ID should not be empty when the error is nil.
// Non nil error will result in an internal server error.
// This function can use the given request ID if necessary.
var NewTraceID func(reqID string) (string, error)

// tracker is the tracking middleware.
// This implements core.Middleware interface.
type tracker struct {
	eh core.ErrorHandler

	// reqProxyHeader is the HTTP header name to proxy a request ID.
	// Value must be formatted with textproto.CanonicalMIMEHeaderKey.
	// If empty, request ID is not proxied.
	reqProxyHeader string
	// trcProxyHeader is the HTTP header name to proxy a trace ID.
	// Value must be formatted with textproto.CanonicalMIMEHeaderKey.
	// If empty, trace ID is not set in proxy header.
	trcProxyHeader string
	// trcExtractHeader is the HTTP header name to extract trace ID from header.
	// Value must be formatted with textproto.CanonicalMIMEHeaderKey.
	// If empty, newly generated trace ID is always used.
	trcExtractHeader string

	// newReqID returns a new request ID.
	newReqID func() (string, error)
	// newTrcID returns a new trace ID.
	// newTrcID can use the given request ID
	// to generate a trace ID.
	newTrcID func(string) (string, error)
}

func (m *tracker) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var reqID string // Request ID
		var trcID string // Trace ID
		var err error

		ctx := r.Context()
		reqID = uid.IDFromContext(ctx)
		if reqID == "" {
			reqID, err = m.newReqID()
			if err != nil {
				err = app.ErrAppMiddleGenID.WithStack(err, map[string]any{"type": "request"})
				m.eh.ServeHTTPError(w, r, utilhttp.NewHTTPError(err, http.StatusInternalServerError))
				return
			}
			ctx = uid.ContextWithID(ctx, reqID)
		}

		if m.trcExtractHeader != "" {
			trcID = r.Header.Get(m.trcExtractHeader) // Use given trace ID.
		}
		if trcID == "" {
			trcID, err = m.newTrcID(reqID) // Generate a new trace ID.
			if err != nil {
				err = app.ErrAppMiddleGenID.WithStack(err, map[string]any{"type": "trace"})
				m.eh.ServeHTTPError(w, r, utilhttp.NewHTTPError(err, http.StatusInternalServerError))
				return
			}
		}

		h := utilhttp.ProxyHeaderFromContext(ctx)
		if h == nil {
			h = make(http.Header)
			ctx = utilhttp.ContextWithProxyHeader(ctx, h)
		}
		if m.reqProxyHeader != "" {
			h.Set(m.reqProxyHeader, reqID)
		}
		if m.trcProxyHeader != "" {
			h.Set(m.trcProxyHeader, trcID)
		}

		// Save request id and trace id in the context so the
		// access loggers can output it to the access logs.
		ids := log.NewCustomAttrs("ids", map[string]any{"request": reqID, "trace": trcID})
		ctx = log.ContextWithAttrs(ctx, ids)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
