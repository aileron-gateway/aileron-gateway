package httpserver

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/aileron-gateway/aileron-gateway/core"
	"github.com/aileron-gateway/aileron-gateway/kernel/errorutil"
	"github.com/aileron-gateway/aileron-gateway/kernel/log"
	utilhttp "github.com/aileron-gateway/aileron-gateway/util/http"
)

// recoverer handle panics and recover from panics.
// This implements core.Middleware interface.
// See https://go.dev/blog/defer-panic-and-recover
type recoverer struct {
	lg log.Logger
	eh core.ErrorHandler
}

func (m *recoverer) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ww := utilhttp.WrapWriter(w)

		defer func(ctx context.Context) {
			// Get panicked content.
			// rec is not always error but can be string or other types.
			rec := recover()
			if rec == nil {
				return
			}

			var err errorutil.Attributes
			if e, ok := rec.(error); ok {
				// panic(http.ErrAbortHandler) does not output stacktrace on the server's log.
				// This happens, for example, when client canceled the connection.
				// So, do not output logs and stacks to the application log as well.
				// Check out the comment of http.ErrAbortHandler.
				if errors.Is(e, http.ErrAbortHandler) {
					m.eh.ServeHTTPError(w, r, utilhttp.NewHTTPError(e, -1)) // Logging only.
					return
				}
				err = core.ErrCoreServerRecover.WithStack(e, nil)
			} else {
				err = core.ErrCoreServerRecover.WithStack(fmt.Errorf("%v", rec), nil)
			}

			m.lg.Error(ctx, "panic recovered", err.Name(), err.Map())

			if !ww.Written() {
				m.eh.ServeHTTPError(w, r, utilhttp.NewHTTPError(err, http.StatusInternalServerError))
			}
		}(r.Context())

		next.ServeHTTP(ww, r)
	})
}
