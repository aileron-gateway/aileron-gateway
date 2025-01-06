package httpproxy

import (
	"net/http"
	"net/textproto"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	v1 "github.com/aileron-gateway/aileron-gateway/apis/core/v1"
	"github.com/aileron-gateway/aileron-gateway/apis/kernel"
	"github.com/aileron-gateway/aileron-gateway/core"
	"github.com/aileron-gateway/aileron-gateway/kernel/txtutil"
)

// newMatcher returns a new matcherFunc.
// The given spec must not be nil, otherwise panics.
func newMatcher(spec *v1.PathMatcherSpec) (matcherFunc, error) {
	matcher := &matcher{
		pattern:      spec.Match,
		rewrite:      spec.Rewrite,
		trimPrefix:   spec.TrimPrefix,
		appendPrefix: spec.AppendPrefix,
	}
	switch spec.MatchType {
	case kernel.MatchType_Exact:
		return matcher.exact, nil
	case kernel.MatchType_Prefix:
		return matcher.prefix, nil
	case kernel.MatchType_Suffix:
		return matcher.suffix, nil
	case kernel.MatchType_Contains:
		return matcher.contains, nil
	case kernel.MatchType_Path:
		if _, err := path.Match(spec.Match, "syntax test"); err != nil {
			reason := "invalid path pattern for 'Path' type. See https://pkg.go.dev/path#Match"
			return nil, core.ErrCoreGenCreateComponent.WithStack(err, map[string]any{"reason": reason})
		}
		return matcher.path, nil
	case kernel.MatchType_FilePath:
		if _, err := filepath.Match(spec.Match, "syntax test"); err != nil {
			reason := "invalid path pattern for 'FilePath' type. See https://pkg.go.dev/filepath#Match"
			return nil, core.ErrCoreGenCreateComponent.WithStack(err, map[string]any{"reason": reason})
		}
		return matcher.filePath, nil
	case kernel.MatchType_Regex:
		regexp, err := regexp.Compile(spec.Match)
		if err != nil {
			reason := "invalid regular expression for 'Regex' type. See https://pkg.go.dev/regexp/syntax"
			return nil, core.ErrCoreGenCreateComponent.WithStack(err, map[string]any{"reason": reason})
		}
		matcher.regexp = regexp
		return matcher.regEx, nil
	case kernel.MatchType_RegexPOSIX:
		regexp, err := regexp.CompilePOSIX(spec.Match)
		if err != nil {
			reason := "invalid regular expression for 'RegexPOSIX' type. See https://pkg.go.dev/regexp/syntax"
			return nil, core.ErrCoreGenCreateComponent.WithStack(err, map[string]any{"reason": reason})
		}
		matcher.regexp = regexp
		return matcher.regEx, nil
	default:
		return matcher.prefix, nil
	}
}

// matcherFunc is the function of path matcher.
// matcherFunc gets requested URL path and return proxy path.
// The returned value of bool represents is the path matched of not.
// If false, the proxy have to reject the request or try other matcherFuncs.
type matcherFunc func(target string) (string, bool)

// matcher is a path matcher object.
// This struct provides matcherFunc
type matcher struct {
	// pattern is the URL path pattern.
	pattern string
	// trimPrefix is the prefix string
	// which is removed from the target.
	// This prefix is trimmed before checking match.
	trimPrefix string
	// appendPrefix is the prefix string
	// which is appended to the target.
	// This prefix is appended after checking match.
	appendPrefix string
	// regexp is the pre compiled regular expression pattern.
	// This field is used when matching the path with regular expression.
	regexp *regexp.Regexp
	// rewrite is the path rewrite pattern.
	// This field is used for modifying the requested path for proxy.
	rewrite string
}

// exact is the matherFunc of exact matching.
func (m *matcher) exact(target string) (string, bool) {
	target = strings.TrimPrefix(target, m.trimPrefix)
	return m.appendPrefix + target, (m.pattern == target)
}

// prefix is the matherFunc of prefix matching.
func (m *matcher) prefix(target string) (string, bool) {
	target = strings.TrimPrefix(target, m.trimPrefix)
	return m.appendPrefix + target, strings.HasPrefix(target, m.pattern)
}

// suffix is the matherFunc of suffix matching.
func (m *matcher) suffix(target string) (string, bool) {
	target = strings.TrimPrefix(target, m.trimPrefix)
	return m.appendPrefix + target, strings.HasSuffix(target, m.pattern)
}

// contains is the matherFunc of containing matching.
func (m *matcher) contains(target string) (string, bool) {
	target = strings.TrimPrefix(target, m.trimPrefix)
	return m.appendPrefix + target, strings.Contains(target, m.pattern)
}

// path is the matherFunc of path matching.
// Checkout https://pkg.go.dev/path#Match
func (m *matcher) path(target string) (string, bool) {
	target = strings.TrimPrefix(target, m.trimPrefix)
	matched, _ := path.Match(m.pattern, target)
	return m.appendPrefix + target, matched
}

// filePath is the matherFunc of filePath matching.
// Checkout https://pkg.go.dev/filepath#Match
func (m *matcher) filePath(target string) (string, bool) {
	target = strings.TrimPrefix(target, m.trimPrefix)
	matched, _ := filepath.Match(m.pattern, target)
	return m.appendPrefix + target, matched
}

// regEx is the matherFunc of regular expression matching.
// Checkout https://pkg.go.dev/regexp#Regexp.ExpandString
func (m *matcher) regEx(target string) (string, bool) {
	target = strings.TrimPrefix(target, m.trimPrefix)

	if !m.regexp.MatchString(target) {
		return m.appendPrefix + target, false
	}

	// No need to modify requested path.
	if m.rewrite == "" {
		return m.appendPrefix + target, true
	}

	// Modify requested path for proxy.
	result := []byte{}
	for _, submatches := range m.regexp.FindAllStringSubmatchIndex(target, -1) {
		result = m.regexp.ExpandString(result, m.rewrite, target, submatches)
	}

	return m.appendPrefix + string(result), true
}

// pathParamMatchers returns pathParam matchers.
// nil spec and specs with empty key string are ignored.
func pathParamMatchers(specs ...*v1.ParamMatcherSpec) ([]txtutil.Matcher[*http.Request], error) {
	matchers := make([]txtutil.Matcher[*http.Request], 0, len(specs))
	for _, s := range specs {
		if s == nil || s.Key == "" {
			continue
		}
		matchFunc, err := txtutil.NewStringMatcher(txtutil.MatchType(s.MatchType), s.Patterns...)
		if err != nil {
			return nil, err
		}
		matchers = append(matchers, &pathParamMatcher{
			key: s.Key,
			f:   matchFunc.Match,
		})
	}
	return matchers, nil
}

// pathParamMatcher is a matcher for a path parameter.
// pathParamMatcher implements core.Matcher[*http.Request] interface.
type pathParamMatcher struct {
	key string
	f   txtutil.MatchFunc[string]
}

func (m *pathParamMatcher) Match(r *http.Request) bool {
	p := r.PathValue(m.key)
	if p == "" {
		return false
	}
	return m.f(p)
}

// headerMatchers returns header matchers.
// nil spec and specs with empty key string are ignored.
func headerMatchers(specs ...*v1.ParamMatcherSpec) ([]txtutil.Matcher[*http.Request], error) {
	matchers := make([]txtutil.Matcher[*http.Request], 0, len(specs))
	for _, s := range specs {
		if s == nil || s.Key == "" {
			continue
		}
		matchFunc, err := txtutil.NewStringMatcher(txtutil.MatchType(s.MatchType), s.Patterns...)
		if err != nil {
			return nil, err
		}
		matchers = append(matchers, &headerMatcher{
			key: textproto.CanonicalMIMEHeaderKey(s.Key),
			f:   matchFunc.Match,
		})
	}
	return matchers, nil
}

// headerMatcher is a matcher for a HTTP header.
// If multiple header values were found, joined string by commas
// "," is input for the match function.
// headerMatcher implements core.Matcher[*http.Request] interface.
type headerMatcher struct {
	// key is the header name.
	// This must be canonical format.
	// Use textproto.CanonicalMIMEHeaderKey.
	key string
	f   txtutil.MatchFunc[string]
}

func (m *headerMatcher) Match(r *http.Request) bool {
	v := r.Header[m.key]
	switch len(v) {
	case 0:
		return false
	case 1:
		return m.f(v[0])
	default:
		return m.f(strings.Join(v, ","))
	}
}

// queryMatchers returns query matchers.
// nil spec and specs with empty key string are ignored.
func queryMatchers(specs ...*v1.ParamMatcherSpec) ([]txtutil.Matcher[*http.Request], error) {
	matchers := make([]txtutil.Matcher[*http.Request], 0, len(specs))
	for _, s := range specs {
		if s == nil || s.Key == "" {
			continue
		}
		matchFunc, err := txtutil.NewStringMatcher(txtutil.MatchType(s.MatchType), s.Patterns...)
		if err != nil {
			return nil, err
		}
		matchers = append(matchers, &queryMatcher{
			key: s.Key,
			f:   matchFunc.Match,
		})
	}
	return matchers, nil
}

// queryMatcher is a matcher for a URL query.
// If multiple query values were found, joined string by commas
// "," is input for the match function.
// queryMatcher implements core.Matcher[*http.Request] interface.
type queryMatcher struct {
	// key is the query key name.
	key string
	f   txtutil.MatchFunc[string]
}

func (m *queryMatcher) Match(r *http.Request) bool {
	v := r.URL.Query()[m.key]
	switch len(v) {
	case 0:
		return false
	case 1:
		return m.f(v[0])
	default:
		return m.f(strings.Join(v, ","))
	}
}
