package http

import (
	"math"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"time"

	v1 "github.com/aileron-gateway/aileron-gateway/apis/core/v1"
	"github.com/aileron-gateway/aileron-gateway/core"
)

// GetCookie returns single cookie value bounded to the given key.
// For example, if the cookies have "session0=value1", "session1=value2", "session2=value3"
// then the joined values of them "value1value2value3" will be returned.
// Use SetCookie to save a large value into the cookie.
// An empty string "" will be returned when the given key was empty.
func GetCookie(key string, cks []*http.Cookie) string {
	if key == "" {
		return ""
	}

	val := make([]string, len(cks))

	for _, ck := range cks {
		if !strings.HasPrefix(ck.Name, key) {
			continue
		}

		i, err := strconv.Atoi(ck.Name[len(key):])
		if err != nil {
			// Cookie name might be changed by the user.
			continue
		}

		if i > len(cks)-1 {
			// Cookie name might be changed by the user.
			continue
		}

		val[i] = ck.Value
	}

	return strings.Join(val, "")
}

// SetCookie sets a large value to the cookie by splitting it into smaller values.
// Maximum size of the value is restricted to 1<<12 - 1<<7 = 3968 bytes.
// Use GetCookie to extract the original value from cookies.
// This function do nothing and return immediately when the given prefix was empty..
func SetCookie(w http.ResponseWriter, cookieNames []string, c core.CookieCreator, prefix string, val string) {
	if prefix == "" {
		return
	}

	// let maximum length of cookie be 4096 - 128 bytes.
	// many browsers restrict the size of cookie to 4096 bytes at maximum.
	// Cookie name and other attributes are considered to be contained in the 128 bytes.
	max := 1<<12 - 1<<7
	n := len(val)
	m := int(math.Ceil(float64(len(val)) / float64(max)))
	l := int(math.Ceil(float64(n) / float64(m)))

	// Set cookies.
	for i := 0; i < m; i++ {
		name := prefix + strconv.Itoa(i)
		v := val[i*l : minInt((i+1)*l, n)] // Get the part of the value.
		ck := c.NewCookie()
		ck.Name = name
		ck.Value = v
		w.Header().Add("Set-Cookie", ck.String())
	}

	// Delete cookies.
	for i := m; i < len(cookieNames); i++ {
		name := prefix + strconv.Itoa(i)
		if slices.Contains(cookieNames, name) {
			ck := &http.Cookie{}
			ck.Name = name
			ck.MaxAge = -1
			w.Header().Add("Set-Cookie", ck.String())
		}
	}
}

// DeleteCookie deletes the cookies which have the given prefix.
// If the prefix is "session", then the cookies with name
// "session0", "session1", "session2", ... will be deleted.
func DeleteCookie(w http.ResponseWriter, cookieNames []string, prefix string) {
	if prefix == "" {
		return
	}
	for i := 0; i < len(cookieNames); i++ {
		name := prefix + strconv.Itoa(i)
		if slices.Contains(cookieNames, name) {
			ck := &http.Cookie{}
			ck.Name = name
			ck.MaxAge = -1
			w.Header().Add("Set-Cookie", ck.String())
		}
	}
}

// CookieNames returns a list of cookie names.
func CookieNames(cks []*http.Cookie) []string {
	names := make([]string, len(cks))
	for i, ck := range cks {
		names[i] = ck.Name
	}
	return names
}

// minInt returns the smaller value.
func minInt(x, y int) int {
	if x > y {
		return y
	}
	return x
}

// DefaultCookieCreator return a new instance of cookie creator
// of the default implementation.
// Attributes are as follows. Cookie name and value should be changed at least.
//   - Name: ""
//   - Value: ""
//   - Path: "/"
//   - Domain: ""
//   - ExpiresIn: 0
//   - MaxAge: 0
//   - Secure: true
//   - HTTPOnly: true
//   - SameSite: Default
//
// References:
//   - https://datatracker.ietf.org/doc/html/rfc6265
//   - https://developer.mozilla.org/en-US/docs/Web/HTTP/Cookies
//   - https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Set-Cookie
func DefaultCookieCreator() *CookieCreator {
	return &CookieCreator{
		Name:      "", // Name must be set by user.
		Value:     "", // Value must be set by user.
		Path:      "/",
		Domain:    "",
		ExpiresIn: 0,
		MaxAge:    0,
		Secure:    true,
		HTTPOnly:  true,
		SameSite:  http.SameSiteDefaultMode,
	}
}

// NewCookieCreator return a new cookie creator object.
// Default cookie creator is returned when nil was given.
//
// References:
//   - https://datatracker.ietf.org/doc/html/rfc6265
//   - https://developer.mozilla.org/en-US/docs/Web/HTTP/Cookies
//   - https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Set-Cookie
func NewCookieCreator(spec *v1.CookieSpec) *CookieCreator {
	if spec == nil {
		return DefaultCookieCreator()
	}
	return &CookieCreator{
		Name:      spec.Name,
		Value:     spec.Value,
		Path:      spec.Path,
		Domain:    spec.Domain,
		ExpiresIn: time.Second * time.Duration(spec.ExpiresIn),
		MaxAge:    int(spec.MaxAge),
		Secure:    spec.Secure,
		HTTPOnly:  spec.HTTPOnly,
		SameSite:  sameSite(spec.SameSite),
	}
}

// sameSite returns a http.SameSite corresponding to the given v1.SameSite.
// The returned value will be
//   - v1.SameSite_Default > http.SameSiteDefaultMode
//   - v1.SameSite_Lax     > http.SameSiteLaxMode
//   - v1.SameSite_Strict  > http.SameSiteStrictMode
//   - v1.SameSite_None    > http.SameSiteNoneMode
//   - Others              > http.SameSite(0)
func sameSite(val v1.SameSite) http.SameSite {
	switch val {
	case v1.SameSite_Default:
		return http.SameSiteDefaultMode
	case v1.SameSite_Lax:
		return http.SameSiteLaxMode
	case v1.SameSite_Strict:
		return http.SameSiteStrictMode
	case v1.SameSite_None:
		return http.SameSiteNoneMode
	}
	return http.SameSite(0)
}

// CookieCreator creates a new cookie with predefined cookie attributes.
// Be careful to set all attribute as secure as possible.
//
// References:
//   - https://datatracker.ietf.org/doc/html/rfc6265
//   - https://developer.mozilla.org/en-US/docs/Web/HTTP/Cookies
//   - https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Set-Cookie
type CookieCreator struct {
	Name      string
	Value     string
	Path      string
	Domain    string
	ExpiresIn time.Duration
	// MaxAge=0 means no 'Max-Age' attribute specified.
	// MaxAge<0 means delete cookie now, equivalently 'Max-Age: 0'
	// MaxAge>0 means Max-Age attribute present and given in seconds
	MaxAge   int
	Secure   bool
	HTTPOnly bool
	SameSite http.SameSite
}

// NewCookie returns a new instance of cookie.
func (c *CookieCreator) NewCookie() *http.Cookie {
	ck := &http.Cookie{
		Name:     c.Name,
		Value:    c.Value,
		Path:     c.Path,
		Domain:   c.Domain,
		MaxAge:   c.MaxAge,
		Secure:   c.Secure,
		HttpOnly: c.HTTPOnly,
		SameSite: c.SameSite,
	}
	if c.ExpiresIn > 0 {
		ck.Expires = time.Now().Add(c.ExpiresIn)
	}
	return ck
}
