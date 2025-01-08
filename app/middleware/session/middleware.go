package session

import (
	"context"
	"net/http"

	"github.com/aileron-gateway/aileron-gateway/core"
	utilhttp "github.com/aileron-gateway/aileron-gateway/util/http"
	"github.com/aileron-gateway/aileron-gateway/util/session"
)

// saveSessionWriter merges http.ResponseWriter and sessionSaver.
// This makes it possible to save session before writing response.
// This is needed when, for example, using cookies as session storage, session data
// must be written to the response header before writing status code.
type saveSessionWriter struct {
	http.ResponseWriter

	ss      session.Session
	store   session.Store
	saved   bool
	saveErr error
}

func (s *saveSessionWriter) saveSession(ctx context.Context) {
	if s.saved {
		return
	}
	s.saved = true
	s.saveErr = s.store.Save(ctx, s.ResponseWriter, s.ss)
}

func (s *saveSessionWriter) Unwrap() http.ResponseWriter {
	return s.ResponseWriter
}

// WriteHeader wraps WriteHeader method of http.ResponseWriter.
// HTTP response headers must be written before setting the status code.
// This is needed especially when using the cookie store as session store.
func (s *saveSessionWriter) WriteHeader(statusCode int) {
	if !s.saved {
		s.saveSession(context.Background())
	}
	s.ResponseWriter.WriteHeader(statusCode)
}

// Write wraps Write method of http.ResponseWriter.
// It can be possible that the response body would written before writing the status code.
// Session data should be saved before writing the body in that situation.
func (s *saveSessionWriter) Write(b []byte) (int, error) {
	if !s.saved {
		s.saveSession(context.Background())
	}
	if s.saveErr != nil {
		return 0, s.saveErr
	}
	return s.ResponseWriter.Write(b)
}

// sessioner manages session by getting, deleting and saving them.
// This implements core.Middleware interface.
// Session management is one of the most important feature from the view of security.
// So, follow the best practices of session management as well as possible.
//
// References
// - https://cheatsheetseries.owasp.org/cheatsheets/Session_Management_Cheat_Sheet.html
type sessioner struct {
	eh core.ErrorHandler
	// store is the session store.
	store session.Store
}

func (m *sessioner) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ss, err := m.store.Get(r)
		if err != nil {
			m.eh.ServeHTTPError(w, r, err)
			return
		}
		ctx := session.ContextWithSession(r.Context(), ss)

		ssw := &saveSessionWriter{
			ResponseWriter: w,
			ss:             ss,
			store:          m.store,
		}
		defer func(ctx context.Context) {
			if !ssw.saved {
				ssw.saveSession(ctx)
			}
			if ssw.saveErr != nil {
				m.eh.ServeHTTPError(w, r, utilhttp.NewHTTPError(ssw.saveErr, http.StatusInternalServerError))
			}
		}(ctx)

		ctx = utilhttp.ContextWithPreProxyHook(ctx, func(r *http.Request) error {
			return m.store.Save(ctx, w, ss)
		})

		r = r.WithContext(ctx)
		next.ServeHTTP(ssw, r)
	})
}
