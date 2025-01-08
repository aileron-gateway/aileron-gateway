package healthcheck

import (
	"context"
	"net/http"
	"time"

	"github.com/aileron-gateway/aileron-gateway/app"
	"github.com/aileron-gateway/aileron-gateway/core"
	utilhttp "github.com/aileron-gateway/aileron-gateway/util/http"
)

// healthCheck is the HTTP handler for health check.
// This implements http.Handler interface.
type healthCheck struct {
	*utilhttp.HandlerBase

	eh core.ErrorHandler

	// timeout is the timeout duration in second.
	// The handler return failure status when health
	// checking time exceeded the duration.
	timeout time.Duration

	// checkers is the list of health checkers.
	// Health check status will be fail when at least
	// 1 checker returned unhealthy status.
	checkers []app.HealthChecker
}

func (h *healthCheck) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Apply timeout fo the context.
	ctx, cancel := context.WithTimeout(r.Context(), h.timeout)
	defer cancel()

	w = utilhttp.WrapWriter(w)

	// Call external health checkers.
	// Return an error if one of them returned failure.
	for _, c := range h.checkers {
		newCtx, ok := c.HealthCheck(ctx)
		if !ok {
			h.eh.ServeHTTPError(w, r, utilhttp.ErrInternalServerError)
			return
		}

		if ctx.Err() == context.DeadlineExceeded {
			h.eh.ServeHTTPError(w, r, utilhttp.ErrGatewayTimeout)
			return
		}

		ctx = newCtx // Update context.
	}

	// TODO: make response Content-Type and body configurable.
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("{}"))
}
