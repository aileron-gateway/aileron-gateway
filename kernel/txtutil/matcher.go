// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package txtutil

import (
	"bytes"
	"path"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"

	k "github.com/aileron-gateway/aileron-gateway/apis/kernel"
	"github.com/aileron-gateway/aileron-gateway/kernel/er"
)

// MatchType is the definition of match types.
// Following types are supported.
//
//   - Exact
//   - Prefix
//   - Suffix
//   - Contains
//   - Path
//   - FilePath
//   - Regex
//   - RegexPOSIX
type MatchType int

const (
	MatchTypeExact      MatchType = iota // Exact matching.
	MatchTypePrefix                      // Prefix matching.
	MatchTypeSuffix                      // Suffix matching.
	MatchTypeContains                    // Contains matching.
	MatchTypePath                        // Path matching.
	MatchTypeFilePath                    // FilePath matching.
	MatchTypeRegex                       // Regular expression matching.
	MatchTypeRegexPOSIX                  // POSIX regular expression matching.
)

// MatchTypes is the match types mapping
// from kernel.MatchType to MatchType.
//
// Following types are contained.
//   - Exact
//   - Prefix
//   - Suffix
//   - Contains
//   - Path
//   - FilePath
//   - Regex
//   - RegexPOSIX
var MatchTypes = map[k.MatchType]MatchType{
	k.MatchType_Exact:      MatchTypeExact,
	k.MatchType_Prefix:     MatchTypePrefix,
	k.MatchType_Suffix:     MatchTypeSuffix,
	k.MatchType_Contains:   MatchTypeContains,
	k.MatchType_Path:       MatchTypePath,
	k.MatchType_FilePath:   MatchTypeFilePath,
	k.MatchType_Regex:      MatchTypeRegex,
	k.MatchType_RegexPOSIX: MatchTypeRegexPOSIX,
}

// MatchFunc is a function that accept an object
// and returns if it matches to a certain condition.
//
// Example to define a string match function:
//
//	// func(string) bool
//	type StringMatchFunc MatchFunc[string]
type MatchFunc[T any] func(T) bool

// Matcher is an interface that accept an object
// and returns if it matches to a certain condition.
//
// Example to define a string match function:
//
//	// interface{ Match(string) bool }
//	type StringMatcher Matcher[string]
type Matcher[T any] interface {
	Match(T) bool
}

// NewStringMatchers returns string matchers from given matcher specs.
// nil specs are ignored.
func NewStringMatchers(specs ...*k.MatcherSpec) ([]Matcher[string], error) {
	ms := make([]Matcher[string], 0, len(specs))
	for _, s := range specs {
		if s == nil {
			continue // Ignore nil spec.
		}
		typ, ok := MatchTypes[s.MatchType]
		if !ok {
			return nil, &er.Error{Package: ErrPkg, Type: ErrTypeMatcher, Description: ErrDscUnsupported, Detail: "NewStringMatchers"}
		}
		m, err := NewStringMatcher(typ, s.Patterns...)
		if err != nil {
			return nil, err // Return err as-is.
		}
		ms = append(ms, m)
	}
	return ms, nil
}

// NewStringMatcher returns a new instance of string matcher.
func NewStringMatcher(typ MatchType, patterns ...string) (Matcher[string], error) {
	pats := slices.Clone(patterns)
	slices.Sort(pats)                        // slices.Compact require sorts.
	pats = slices.Clip(slices.Compact(pats)) // Remove duplicates.

	m := &stringMatcher{
		patterns: pats,
	}

	f, ok := map[MatchType](func(int, string) bool){
		MatchTypeExact:      m.exact,
		MatchTypePrefix:     m.prefix,
		MatchTypeSuffix:     m.suffix,
		MatchTypeContains:   m.contain,
		MatchTypePath:       m.path,
		MatchTypeFilePath:   m.filePath,
		MatchTypeRegex:      m.regex,
		MatchTypeRegexPOSIX: m.regex,
	}[typ]
	if !ok {
		return nil, &er.Error{Package: ErrPkg, Type: ErrTypeMatcher, Description: ErrDscUnsupported, Detail: "match type " + strconv.Itoa(int(typ))}
	}
	m.match = f

	switch typ {
	case MatchTypePath: // Format check.
		for _, p := range pats {
			_, err := path.Match(p, "")
			if err != nil {
				return nil, (&er.Error{Package: ErrPkg, Type: ErrTypeMatcher, Description: ErrDscPattern, Detail: "Path `" + p + "`"}).Wrap(err)
			}
		}
	case MatchTypeFilePath: // Format check.
		for _, p := range pats {
			_, err := filepath.Match(p, "")
			if err != nil {
				return nil, (&er.Error{Package: ErrPkg, Type: ErrTypeMatcher, Description: ErrDscPattern, Detail: "FilePath `" + p + "`"}).Wrap(err)
			}
		}
	case MatchTypeRegex:
		m.regexps = make([]*regexp.Regexp, len(pats))
		for i, p := range pats {
			exp, err := regexp.Compile(p)
			if err != nil {
				return nil, (&er.Error{Package: ErrPkg, Type: ErrTypeMatcher, Description: ErrDscPattern, Detail: "Regex `" + p + "`"}).Wrap(err)
			}
			m.regexps[i] = exp
		}
	case MatchTypeRegexPOSIX:
		m.regexps = make([]*regexp.Regexp, len(pats))
		for i, p := range pats {
			exp, err := regexp.CompilePOSIX(p)
			if err != nil {
				return nil, (&er.Error{Package: ErrPkg, Type: ErrTypeMatcher, Description: ErrDscPattern, Detail: "Regex POSIX `" + p + "`"}).Wrap(err)
			}
			m.regexps[i] = exp
		}
	}

	return m, nil
}

// stringMatcher provides string matching functions.
// This implements Matcher[string] interface.
type stringMatcher struct {
	patterns []string
	regexps  []*regexp.Regexp
	match    func(int, string) bool
}

func (m *stringMatcher) Match(target string) bool {
	for i := 0; i < len(m.patterns); i++ {
		if m.match(i, target) {
			return true
		}
	}
	return false
}

func (m *stringMatcher) exact(i int, target string) bool {
	return m.patterns[i] == target
}

func (m *stringMatcher) prefix(i int, target string) bool {
	return strings.HasPrefix(target, m.patterns[i])
}

func (m *stringMatcher) suffix(i int, target string) bool {
	return strings.HasSuffix(target, m.patterns[i])
}

func (m *stringMatcher) contain(i int, target string) bool {
	return strings.Contains(target, m.patterns[i])
}

func (m *stringMatcher) path(i int, target string) bool {
	ok, _ := path.Match(m.patterns[i], target)
	return ok
}

func (m *stringMatcher) filePath(i int, target string) bool {
	ok, _ := filepath.Match(m.patterns[i], target)
	return ok
}

func (m *stringMatcher) regex(i int, target string) bool {
	return m.regexps[i].MatchString(target)
}

// NewBytesMatchers returns bytes matchers from given matcher specs.
// nil specs are ignored.
func NewBytesMatchers(specs ...*k.MatcherSpec) ([]Matcher[[]byte], error) {
	ms := make([]Matcher[[]byte], 0, len(specs))
	for _, s := range specs {
		if s == nil {
			continue // Ignore nil spec.
		}
		typ, ok := MatchTypes[s.MatchType]
		if !ok {
			return nil, &er.Error{Package: ErrPkg, Type: ErrTypeMatcher, Description: ErrDscUnsupported, Detail: s.MatchType.String()}
		}
		m, err := NewBytesMatcher(typ, s.Patterns...)
		if err != nil {
			return nil, err // Return err as-is.
		}
		ms = append(ms, m)
	}
	return ms, nil
}

// NewBytesMatcher returns a new instance of bytes matcher.
func NewBytesMatcher(typ MatchType, patterns ...string) (Matcher[[]byte], error) {
	pats := make([][]byte, 0, len(patterns))
	for _, p := range patterns {
		pats = append(pats, []byte(p))
	}
	m := &bytesMatcher{
		patternsStr: slices.Clip(slices.Clone(patterns)),
		patterns:    pats,
	}

	f, ok := map[MatchType](func(int, []byte) bool){
		MatchTypeExact:      m.exact,
		MatchTypePrefix:     m.prefix,
		MatchTypeSuffix:     m.suffix,
		MatchTypeContains:   m.contain,
		MatchTypePath:       m.path,
		MatchTypeFilePath:   m.filePath,
		MatchTypeRegex:      m.regex,
		MatchTypeRegexPOSIX: m.regex,
	}[typ]
	if !ok {
		return nil, &er.Error{Package: ErrPkg, Type: ErrTypeMatcher, Description: ErrDscUnsupported, Detail: "match type " + strconv.Itoa(int(typ))}
	}
	m.match = f

	switch typ {
	case MatchTypePath: // Format check.
		for _, p := range patterns {
			_, err := path.Match(p, "")
			if err != nil {
				return nil, &er.Error{Package: ErrPkg, Type: ErrTypeMatcher, Description: ErrDscPattern, Detail: "Path `" + p + "`"}
			}
		}
	case MatchTypeFilePath: // Format check.
		for _, p := range patterns {
			_, err := filepath.Match(p, "")
			if err != nil {
				return nil, &er.Error{Package: ErrPkg, Type: ErrTypeMatcher, Description: ErrDscPattern, Detail: "FilePath `" + p + "`"}
			}
		}
	case MatchTypeRegex:
		m.regexps = make([]*regexp.Regexp, len(patterns))
		for i, p := range patterns {
			exp, err := regexp.Compile(p)
			if err != nil {
				return nil, &er.Error{Package: ErrPkg, Type: ErrTypeMatcher, Description: ErrDscPattern, Detail: "Regex `" + p + "`"}
			}
			m.regexps[i] = exp
		}
	case MatchTypeRegexPOSIX:
		m.regexps = make([]*regexp.Regexp, len(patterns))
		for i, p := range patterns {
			exp, err := regexp.CompilePOSIX(p)
			if err != nil {
				return nil, &er.Error{Package: ErrPkg, Type: ErrTypeMatcher, Description: ErrDscPattern, Detail: "Regex POSIX `" + p + "`"}
			}
			m.regexps[i] = exp
		}
	}

	return m, nil
}

// bytesMatcher provides bytes matching functions.
// This implements Matcher[[]byte] interface.
type bytesMatcher struct {
	patternsStr []string
	patterns    [][]byte
	regexps     []*regexp.Regexp
	match       func(int, []byte) bool
}

func (m *bytesMatcher) Match(target []byte) bool {
	for i := 0; i < len(m.patterns); i++ {
		if m.match(i, target) {
			return true
		}
	}
	return false
}

func (m *bytesMatcher) exact(i int, target []byte) bool {
	return bytes.Equal(m.patterns[i], target)
}

func (m *bytesMatcher) prefix(i int, target []byte) bool {
	return bytes.HasPrefix(target, m.patterns[i])
}

func (m *bytesMatcher) suffix(i int, target []byte) bool {
	return bytes.HasSuffix(target, m.patterns[i])
}

func (m *bytesMatcher) contain(i int, target []byte) bool {
	return bytes.Contains(target, m.patterns[i])
}

func (m *bytesMatcher) path(i int, target []byte) bool {
	ok, _ := path.Match(m.patternsStr[i], string(target))
	return ok
}

func (m *bytesMatcher) filePath(i int, target []byte) bool {
	ok, _ := filepath.Match(m.patternsStr[i], string(target))
	return ok
}

func (m *bytesMatcher) regex(i int, target []byte) bool {
	return m.regexps[i].Match(target)
}
