// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package tracking

import (
	"encoding/base32"
	"net/http"

	"github.com/aileron-gateway/aileron-gateway/core"
	"github.com/aileron-gateway/aileron-gateway/kernel/log"
	"github.com/aileron-gateway/aileron-gateway/kernel/uid"
	utilhttp "github.com/aileron-gateway/aileron-gateway/util/http"
	"github.com/aileron-projects/go/zx/zuid"
)

var base32HexEscaped = base32.NewEncoding("0123456789BCDFGHJKLMNPQRSTUVWXYZ")

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
}

func (m *tracker) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var reqID string // Request ID
		var trcID string // Trace ID

		ctx := r.Context()
		reqID = uid.IDFromContext(ctx)
		if reqID == "" {
			reqID = base32HexEscaped.EncodeToString(zuid.NewHost())
			ctx = uid.ContextWithID(ctx, reqID)
		}

		if m.trcExtractHeader != "" {
			trcID = r.Header.Get(m.trcExtractHeader) // Use given trace ID.
		}
		if trcID == "" {
			trcID = reqID
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
