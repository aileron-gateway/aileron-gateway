package csrf

import (
	"net/http"
	"strings"

	"github.com/aileron-gateway/aileron-gateway/app"
	"github.com/aileron-gateway/aileron-gateway/core"
	utilhttp "github.com/aileron-gateway/aileron-gateway/util/http"
)

// csrf is the middleware which applies CSRF(Cross-Site Request Forgery).
// This implements core.Middleware interface.
//
// References
//   - https://owasp.org/www-community/attacks/csrf
type csrf struct {
	*utilhttp.HandlerBase

	eh core.ErrorHandler

	// proxyHeaderName is the header name to proxy
	// verified CSRF token to upstream.
	// This header is removed for non verified requests.
	proxyHeaderName string
	// issueNew is the flag to always generate new csrf token
	// in the token issue handler.
	issueNew bool

	// token is the csrf token create and verifier.
	token *csrfToken

	// st is the csrf strategy.
	st strategy
}

func (m *csrf) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Delete the http header named m.proxyHeaderName if exists
		// because it is not verified in this middleware.
		if m.proxyHeaderName != "" {
			delete(r.Header, m.proxyHeaderName)
		}

		token, err := m.st.get(r)
		if err != nil || token == "" {
			m.eh.ServeHTTPError(w, r, err)
			return
		}

		// Send the verified token value to upstream.
		if m.proxyHeaderName != "" {
			ctx := r.Context()
			h := utilhttp.ProxyHeaderFromContext(ctx)
			if h == nil {
				h = make(http.Header)
			}
			h.Set(m.proxyHeaderName, token)
			ctx = utilhttp.ContextWithProxyHeader(ctx, h)
			r = r.WithContext(ctx)
		}

		next.ServeHTTP(w, r)
	})
}

// ServeHTTP is the implementation of the http.Handler interface.
// This handler returns a new or an already exist csrf token to the clients.
// This handler MUST be protected with CORS policy, or CSRF protection does not work correctly.
// Get a new CSRF token from this handler and put it for API requests.
func (m *csrf) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// verifier should be return a HTTP error if any.
	var token string

	if !m.issueNew {
		token, _ = m.st.get(r)
	}

	if token == "" {
		newToken, err := m.token.new()
		if err != nil {
			err = app.ErrAppMiddleCSRFNewToken.WithoutStack(err, nil)
			m.eh.ServeHTTPError(w, r, utilhttp.NewHTTPError(err, http.StatusInternalServerError))
			return
		}
		token = newToken
	}

	if err := m.st.set(token, w, r); err != nil {
		m.eh.ServeHTTPError(w, r, err)
		return
	}

	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("Cache-Control", "no-store")

	accept := r.Header.Get("Accept")
	switch {
	case strings.Contains(accept, "text/plain"):
		w.Header()["Content-Type"] = []string{"text/plain; charset=utf-8"}
		w.Write([]byte(token))
	case strings.Contains(accept, "application/json") || strings.Contains(accept, "text/json"):
		w.Header()["Content-Type"] = []string{"application/json; charset=utf-8"}
		w.Write([]byte(`{"token":"` + token + `"}`))
	case strings.Contains(accept, "application/xml") || strings.Contains(accept, "text/xml"):
		w.Header()["Content-Type"] = []string{"application/xml; charset=utf-8"}
		w.Write([]byte(`<token>` + token + `</token>`))
	default:
		w.Header()["Content-Type"] = []string{"text/plain; charset=utf-8"}
		w.Write([]byte(token))
	}
}
