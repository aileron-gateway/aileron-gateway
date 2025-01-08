package csrf

import (
	"errors"
	"net/http"
	"regexp"

	"github.com/aileron-gateway/aileron-gateway/core"
	utilhttp "github.com/aileron-gateway/aileron-gateway/util/http"
	"github.com/aileron-gateway/aileron-gateway/util/session"
)

var (
	errInvalidMIME     = errors.New("csrf: invalid mime type")
	errInvalidToken    = errors.New("csrf: invalid csrf token")
	errSessionNotFound = errors.New("csrf: session not found")
)

// strategy is the csrf strategy.
type strategy interface {
	// get returns csrf token for given request.
	// The returned token is the validated one.
	// An empty string and non-nil error should be
	// returned when the request has no valid csrf token.
	get(*http.Request) (string, core.HTTPError)
	// set sets a csrf token to the request.
	set(string, http.ResponseWriter, *http.Request) core.HTTPError
}

// customRequestHeaders protects CSRF using custom header strategy.
// This strategy does not require collaboration with upstream.
// This implements csrf.setupVerifier interface.
// Learn more at
// https://cheatsheetseries.owasp.org/cheatsheets/Cross-Site_Request_Forgery_Prevention_Cheat_Sheet.html#custom-request-headers
//
// How to protect:
//   - Access to a new CSRF token with setup by accessing the token issue endpoint.
//   - Put the CSRF to a API request.
//   - The verifier check if the custom header exists.
//   - The verifier accept the request if the custom header exists and match to the pattern.
type customRequestHeaders struct {
	headerName string
	pattern    *regexp.Regexp
	token      *csrfToken
}

func (s *customRequestHeaders) get(r *http.Request) (string, core.HTTPError) {
	token := r.Header.Get(s.headerName)
	if token == "" {
		return "", utilhttp.NewHTTPError(errInvalidToken, http.StatusForbidden)
	}
	if s.pattern != nil && s.pattern.MatchString(token) {
		return token, nil
	}
	if !s.token.verify(token) {
		return "", utilhttp.NewHTTPError(errInvalidToken, http.StatusForbidden)
	}
	return token, nil
}

func (s *customRequestHeaders) set(token string, _ http.ResponseWriter, _ *http.Request) core.HTTPError {
	return nil
}

// synchronizerTokenVerifier protects CSRF using synchronizer token strategy.
// Session + custom header. Session + hidden parameter.
// This implements csrf.setupVerifier interface.
// Learn more at
// https://cheatsheetseries.owasp.org/cheatsheets/Cross-Site_Request_Forgery_Prevention_Cheat_Sheet.html#custom-request-headers
//
// How to protect:
//   - Access to a new CSRF token with setup by accessing the token issue endpoint. The token is protected by HMAC or encryption.
//   - Created token is also stored in the session.
//   - Put the CSRF to a API request.
//   - The verify method verifies the given token and accept or reject request.
type synchronizerToken struct {
	ext   extractor
	token *csrfToken
}

func (s *synchronizerToken) get(r *http.Request) (string, core.HTTPError) {
	ss := session.SessionFromContext(r.Context())
	if ss == nil {
		return "", utilhttp.NewHTTPError(errSessionNotFound, http.StatusInternalServerError)
	}
	var b []byte
	if err := ss.Extract("__csrf_token__", &b); err != nil {
		return "", utilhttp.NewHTTPError(err, http.StatusInternalServerError)
	}
	token1 := string(b)
	token2, err := s.ext.extract(r)
	if err != nil {
		return "", utilhttp.NewHTTPError(err, http.StatusForbidden)
	}
	if token1 != token2 || !s.token.verify(token1) {
		return "", utilhttp.NewHTTPError(errInvalidToken, http.StatusForbidden)
	}
	return token1, nil
}

func (s *synchronizerToken) set(token string, _ http.ResponseWriter, r *http.Request) core.HTTPError {
	ss := session.SessionFromContext(r.Context())
	if ss == nil {
		return utilhttp.NewHTTPError(errSessionNotFound, http.StatusInternalServerError)
	}
	if err := ss.Persist("__csrf_token__", []byte(token)); err != nil {
		return utilhttp.NewHTTPError(err, http.StatusInternalServerError)
	}
	return nil
}

// doubleSubmitCookies protects CSRF using double submit cookie strategy.
// CSRF token set to the cookie and the form hidden params or json bodies.
// Compare tokens between cookie and form params or json body.
// This implements csrf.setupVerifier interface.
// Learn more at
// https://cheatsheetseries.owasp.org/cheatsheets/Cross-Site_Request_Forgery_Prevention_Cheat_Sheet.html#custom-request-headers
//
// How to protect:
//   - Access to a new CSRF token with setup by accessing the token issue endpoint. The token is protected by HMAC or encryption.
//   - Created token is also stored in the cookie.
//   - Put the CSRF to a API request.
//   - The verify method verifies the given token and accept or reject request.
type doubleSubmitCookies struct {
	token      *csrfToken
	ext        extractor
	cookieName string
	cookie     core.CookieCreator
}

func (s *doubleSubmitCookies) get(r *http.Request) (string, core.HTTPError) {
	ck, err := r.Cookie(s.cookieName)
	if err != nil {
		return "", utilhttp.NewHTTPError(err, http.StatusForbidden)
	}
	token1 := ck.Value
	token2, err := s.ext.extract(r)
	if err != nil {
		return "", utilhttp.NewHTTPError(err, http.StatusForbidden)
	}
	if token1 != token2 || !s.token.verify(token1) {
		return "", utilhttp.NewHTTPError(errInvalidToken, http.StatusForbidden)
	}
	return token1, nil
}

func (s *doubleSubmitCookies) set(token string, w http.ResponseWriter, _ *http.Request) core.HTTPError {
	ck := s.cookie.NewCookie()
	ck.Name = s.cookieName
	ck.Value = token
	w.Header().Add("Set-Cookie", ck.String())
	return nil
}
