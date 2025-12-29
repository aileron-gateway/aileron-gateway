// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package txtutil

import (
	"regexp"
	"testing"

	k "github.com/aileron-gateway/aileron-gateway/apis/kernel"
	"github.com/aileron-gateway/aileron-gateway/internal/testutil"
	"github.com/aileron-gateway/aileron-gateway/kernel/er"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestNewStringMatcher(t *testing.T) {
	type condition struct {
		typ      MatchType
		patterns []string
	}

	type action struct {
		match   []string
		unmatch []string
		err     error
	}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"exact",
			&condition{
				typ:      MatchTypeExact,
				patterns: []string{"test1", "test2", "test3"},
			},
			&action{
				match:   []string{"test1", "test2", "test3"},
				unmatch: []string{"", "test", "test4"},
			},
		),
		gen(
			"prefix", &condition{
				typ:      MatchTypePrefix,
				patterns: []string{"testX", "testY", "testZ"},
			},
			&action{
				match:   []string{"testXXX", "testYYY", "testZZZ"},
				unmatch: []string{"", "test"},
			},
		),
		gen(
			"suffix",
			&condition{
				typ:      MatchTypeSuffix,
				patterns: []string{"XXX", "YYY", "ZZZ"},
			},
			&action{
				match:   []string{"testXXX", "testYYY", "testZZZ"},
				unmatch: []string{"", "test", "testX", "testY", "testZ"},
			},
		),
		gen(
			"contains",
			&condition{
				typ:      MatchTypeContains,
				patterns: []string{"XXX", "YYY", "ZZZ"},
			},
			&action{
				match:   []string{"testXXXtest", "testYYYtest", "testZZZtest"},
				unmatch: []string{"", "test", "X", "Y", "Z"},
			},
		),
		gen(
			"path",
			&condition{
				typ:      MatchTypePath,
				patterns: []string{"/foo/*", "/bar/*"},
			},
			&action{
				match:   []string{"/foo/", "/foo/test", "/bar/", "/bar/test"},
				unmatch: []string{"", "/foo", "/foo*", "/foo/test1/test2"},
			},
		),
		gen(
			"filepath",
			&condition{
				typ:      MatchTypeFilePath,
				patterns: []string{"/foo/*", "/bar/*"},
			},
			&action{
				match:   []string{"/foo/", "/foo/test", "/bar/", "/bar/test"},
				unmatch: []string{"", "/foo", "/foo*"},
			},
		),
		gen(
			"regex",
			&condition{
				typ:      MatchTypeRegex,
				patterns: []string{"test", "^foo[0-9a-zA-Z]*$"},
			},
			&action{
				match:   []string{"test", "foo", "foobar"},
				unmatch: []string{"", "t", "foo#"},
			},
		),
		gen(
			"regex POSIX",
			&condition{
				typ:      MatchTypeRegexPOSIX,
				patterns: []string{"test", "^foo[0-9a-zA-Z]*$"},
			},
			&action{
				match:   []string{"test", "foo", "foobar"},
				unmatch: []string{"", "t", "foo#"},
			},
		),
		gen(
			"unsupported type",
			&condition{
				typ:      MatchType(999), // Does not exist.
				patterns: []string{""},
			},
			&action{
				err: &er.Error{
					Package:     ErrPkg,
					Type:        ErrTypeMatcher,
					Description: ErrDscUnsupported,
				},
			},
		),
		gen(
			"path error",
			&condition{
				typ:      MatchTypePath,
				patterns: []string{"[0-9a-"},
			},
			&action{
				err: &er.Error{
					Package:     ErrPkg,
					Type:        ErrTypeMatcher,
					Description: ErrDscPattern,
				},
			},
		),
		gen(
			"filepath error",
			&condition{
				typ:      MatchTypeFilePath,
				patterns: []string{"[0-9a-"},
			},
			&action{
				err: &er.Error{
					Package:     ErrPkg,
					Type:        ErrTypeMatcher,
					Description: ErrDscPattern,
				},
			},
		),
		gen(
			"regex error",
			&condition{
				typ:      MatchTypeRegex,
				patterns: []string{"[0-9a-"},
			},
			&action{
				err: &er.Error{
					Package:     ErrPkg,
					Type:        ErrTypeMatcher,
					Description: ErrDscPattern,
				},
			},
		),
		gen(
			"regex POSIX error",
			&condition{
				typ:      MatchTypeRegexPOSIX,
				patterns: []string{"[0-9a-"},
			},
			&action{
				err: &er.Error{
					Package:     ErrPkg,
					Type:        ErrTypeMatcher,
					Description: ErrDscPattern,
				},
			},
		),
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			m, err := NewStringMatcher(tt.C.typ, tt.C.patterns...)
			testutil.Diff(t, tt.A.err, err, cmpopts.EquateErrors())
			if err != nil {
				return
			}

			for _, v := range tt.A.match {
				t.Log("Check match:", v)
				testutil.Diff(t, true, m.Match(v))
			}

			for _, v := range tt.A.unmatch {
				t.Log("Check unmatch:", v)
				testutil.Diff(t, false, m.Match(v))
			}
		})
	}
}

func TestNewStringMatchers(t *testing.T) {
	type condition struct {
		specs []*k.MatcherSpec
	}

	type action struct {
		match   []string
		unmatch []string
		err     error
	}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"no spec",
			&condition{
				specs: []*k.MatcherSpec{},
			},
			&action{
				match:   []string{},
				unmatch: []string{"", "test"},
				err:     nil,
			},
		),
		gen(
			"one valid spec",
			&condition{
				specs: []*k.MatcherSpec{
					{
						MatchType: k.MatchType_Exact,
						Patterns:  []string{"test1", "test2", "test3"},
					},
				},
			},
			&action{
				match:   []string{"test1", "test2", "test3"},
				unmatch: []string{"", "test", "test4"},
				err:     nil,
			},
		),
		gen(
			"multiple valid specs",
			&condition{
				specs: []*k.MatcherSpec{
					{
						MatchType: k.MatchType_Exact,
						Patterns:  []string{"test1", "test2", "test3"},
					},
					{
						MatchType: k.MatchType_Prefix,
						Patterns:  []string{"testX", "testY", "testZ"},
					},
				},
			},
			&action{
				match:   []string{"testXXX", "testYYY", "testZZZ"},
				unmatch: []string{"", "test"},
				err:     nil,
			},
		),
		gen(
			"input nil spec",
			&condition{
				specs: []*k.MatcherSpec{
					nil, // This should be ignored.
					{
						MatchType: k.MatchType_Exact,
						Patterns:  []string{"test1", "test2", "test3"},
					},
					{
						MatchType: k.MatchType_Prefix,
						Patterns:  []string{"testX", "testY", "testZ"},
					},
				},
			},
			&action{
				match:   []string{"testXXX", "testYYY", "testZZZ"},
				unmatch: []string{"", "test"},
				err:     nil,
			},
		),
		gen(
			"unsupported type",
			&condition{
				specs: []*k.MatcherSpec{
					{
						MatchType: k.MatchType(999), // 999 does not exist.
						Patterns:  []string{},
					},
				},
			},
			&action{
				err: &er.Error{
					Package:     ErrPkg,
					Type:        ErrTypeMatcher,
					Description: ErrDscUnsupported,
				},
			},
		),
		gen(
			"regex error",
			&condition{
				specs: []*k.MatcherSpec{
					{
						MatchType: k.MatchType_Regex,
						Patterns:  []string{"[0-9a-"},
					},
				},
			},
			&action{
				err: &er.Error{
					Package:     ErrPkg,
					Type:        ErrTypeMatcher,
					Description: ErrDscPattern,
				},
			},
		),
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			ms, err := NewStringMatchers(tt.C.specs...)
			testutil.Diff(t, tt.A.err, err, cmpopts.EquateErrors())

			for _, v := range tt.A.match {
				matched := false
				for _, m := range ms {
					if matched = m.Match(v); matched {
						break
					}
				}
				t.Log("Check match:", v)
				testutil.Diff(t, true, matched)
			}

			for _, v := range tt.A.unmatch {
				matched := false
				for _, m := range ms {
					if matched = m.Match(v); matched {
						break
					}
				}
				t.Log("Check unmatch:", v)
				testutil.Diff(t, false, matched)
			}
		})
	}
}

func TestStringMatcher_Exact(t *testing.T) {
	type condition struct {
		matcher *stringMatcher
		target  string
	}

	type action struct {
		matched bool
	}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"match / 1 pattern",
			&condition{
				matcher: &stringMatcher{
					patterns: []string{"test"},
				},
				target: "test",
			},
			&action{
				matched: true,
			},
		),
		gen(
			"match / multiple pattern",
			&condition{
				matcher: &stringMatcher{
					patterns: []string{"test1", "test2", "test3"},
				},
				target: "test3",
			},
			&action{
				matched: true,
			},
		),
		gen(
			"not match / 1 pattern",
			&condition{
				matcher: &stringMatcher{
					patterns: []string{"test"},
				},
				target: "un-match",
			},
			&action{
				matched: false,
			},
		),
		gen(
			"not match / multiple patterns",
			&condition{
				matcher: &stringMatcher{
					patterns: []string{"test1", "test3", "test3"},
				},
				target: "un-match",
			},
			&action{
				matched: false,
			},
		),
		gen(
			"not match / no patterns",
			&condition{
				matcher: &stringMatcher{
					patterns: []string{},
				},
				target: "", // Empty string.
			},
			&action{
				matched: false,
			},
		),
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			tt.C.matcher.match = tt.C.matcher.exact
			matched := tt.C.matcher.Match(tt.C.target)
			testutil.Diff(t, tt.A.matched, matched)
		})
	}
}

func TestStringMatcher_Prefix(t *testing.T) {
	type condition struct {
		matcher *stringMatcher
		target  string
	}

	type action struct {
		matched bool
	}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"match / 1 pattern",
			&condition{
				matcher: &stringMatcher{
					patterns: []string{"testX"},
				},
				target: "testXXX",
			},
			&action{
				matched: true,
			},
		),
		gen(
			"match / multiple pattern",
			&condition{
				matcher: &stringMatcher{
					patterns: []string{"testX", "testY", "testZ"},
				},
				target: "testZZZ",
			},
			&action{
				matched: true,
			},
		),
		gen(
			"not match / 1 pattern",
			&condition{
				matcher: &stringMatcher{
					patterns: []string{"testX"},
				},
				target: "un-match",
			},
			&action{
				matched: false,
			},
		),
		gen(
			"not match / multiple patterns",
			&condition{
				matcher: &stringMatcher{
					patterns: []string{"testX", "testY", "testZ"},
				},
				target: "un-match",
			},
			&action{
				matched: false,
			},
		),
		gen(
			"not match / no patterns",
			&condition{
				matcher: &stringMatcher{
					patterns: []string{},
				},
				target: "", // Empty string.
			},
			&action{
				matched: false,
			},
		),
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			tt.C.matcher.match = tt.C.matcher.prefix
			matched := tt.C.matcher.Match(tt.C.target)
			testutil.Diff(t, tt.A.matched, matched)
		})
	}
}

func TestStringMatcher_Suffix(t *testing.T) {
	type condition struct {
		matcher *stringMatcher
		target  string
	}

	type action struct {
		matched bool
	}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"match / 1 pattern",
			&condition{
				matcher: &stringMatcher{
					patterns: []string{"XXX"},
				},
				target: "testXXX",
			},
			&action{
				matched: true,
			},
		),
		gen(
			"match / multiple pattern",
			&condition{
				matcher: &stringMatcher{
					patterns: []string{"XXX", "YYY", "ZZZ"},
				},
				target: "testZZZ",
			},
			&action{
				matched: true,
			},
		),
		gen(
			"not match / 1 pattern",
			&condition{
				matcher: &stringMatcher{
					patterns: []string{"XXX"},
				},
				target: "un-match",
			},
			&action{
				matched: false,
			},
		),
		gen(
			"not match / multiple patterns",
			&condition{
				matcher: &stringMatcher{
					patterns: []string{"XXX", "YYY", "ZZZ"},
				},
				target: "un-match",
			},
			&action{
				matched: false,
			},
		),
		gen(
			"not match / no patterns",
			&condition{
				matcher: &stringMatcher{
					patterns: []string{},
				},
				target: "", // Empty string.
			},
			&action{
				matched: false,
			},
		),
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			tt.C.matcher.match = tt.C.matcher.suffix
			matched := tt.C.matcher.Match(tt.C.target)
			testutil.Diff(t, tt.A.matched, matched)
		})
	}
}

func TestStringMatcher_Contain(t *testing.T) {
	type condition struct {
		matcher *stringMatcher
		target  string
	}

	type action struct {
		matched bool
	}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"match / 1 pattern",
			&condition{
				matcher: &stringMatcher{
					patterns: []string{"XXX"},
				},
				target: "testXXXtest",
			},
			&action{
				matched: true,
			},
		),
		gen(
			"match / multiple pattern",
			&condition{
				matcher: &stringMatcher{
					patterns: []string{"XXX", "YYY", "ZZZ"},
				},
				target: "testZZZtest",
			},
			&action{
				matched: true,
			},
		),
		gen(
			"not match / 1 pattern",
			&condition{
				matcher: &stringMatcher{
					patterns: []string{"XXX"},
				},
				target: "un-match",
			},
			&action{
				matched: false,
			},
		),
		gen(
			"not match / multiple patterns",
			&condition{
				matcher: &stringMatcher{
					patterns: []string{"XXX", "YYY", "ZZZ"},
				},
				target: "un-match",
			},
			&action{
				matched: false,
			},
		),
		gen(
			"not match / no patterns",
			&condition{
				matcher: &stringMatcher{
					patterns: []string{},
				},
				target: "", // Empty string.
			},
			&action{
				matched: false,
			},
		),
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			tt.C.matcher.match = tt.C.matcher.contain
			matched := tt.C.matcher.Match(tt.C.target)
			testutil.Diff(t, tt.A.matched, matched)
		})
	}
}

func TestStringMatcher_Path(t *testing.T) {
	type condition struct {
		matcher *stringMatcher
		target  string
	}

	type action struct {
		matched bool
	}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"match / 1 pattern",
			&condition{
				matcher: &stringMatcher{
					patterns: []string{"/test/foo/*"},
				},
				target: "/test/foo/bar",
			},
			&action{
				matched: true,
			},
		),
		gen(
			"match / multiple pattern",
			&condition{
				matcher: &stringMatcher{
					patterns: []string{"/", "/test", "/test/*/*"},
				},
				target: "/test/foo/bar",
			},
			&action{
				matched: true,
			},
		),
		gen(
			"not match / 1 pattern",
			&condition{
				matcher: &stringMatcher{
					patterns: []string{"XXX"},
				},
				target: "un-match",
			},
			&action{
				matched: false,
			},
		),
		gen(
			"not match / multiple patterns",
			&condition{
				matcher: &stringMatcher{
					patterns: []string{"XXX", "YYY", "ZZZ"},
				},
				target: "un-match",
			},
			&action{
				matched: false,
			},
		),
		gen(
			"not match / no patterns",
			&condition{
				matcher: &stringMatcher{
					patterns: []string{},
				},
				target: "", // Empty string.
			},
			&action{
				matched: false,
			},
		),
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			tt.C.matcher.match = tt.C.matcher.path
			matched := tt.C.matcher.Match(tt.C.target)
			testutil.Diff(t, tt.A.matched, matched)
		})
	}
}

func TestStringMatcher_Regex(t *testing.T) {
	type condition struct {
		matcher *stringMatcher
		target  string
	}

	type action struct {
		matched bool
	}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"match / 1 pattern",
			&condition{
				matcher: &stringMatcher{
					patterns: []string{"te.*"},
					regexps: []*regexp.Regexp{
						regexp.MustCompile(`te.*`),
					},
				},
				target: "test",
			},
			&action{
				matched: true,
			},
		),
		gen(
			"match / multiple pattern",
			&condition{
				matcher: &stringMatcher{
					patterns: []string{"foo", "bar", "te.*"},
					regexps: []*regexp.Regexp{
						regexp.MustCompile(`foo`),
						regexp.MustCompile(`bar`),
						regexp.MustCompile(`te.*`),
					},
				},
				target: "test",
			},
			&action{
				matched: true,
			},
		),
		gen(
			"not match / 1 pattern",
			&condition{
				matcher: &stringMatcher{
					patterns: []string{"XXX"},
					regexps: []*regexp.Regexp{
						regexp.MustCompile(`XXX`),
					},
				},
				target: "un-match",
			},
			&action{
				matched: false,
			},
		),
		gen(
			"not match / multiple patterns",
			&condition{
				matcher: &stringMatcher{
					patterns: []string{"XXX", "YYY", "ZZZ"},
					regexps: []*regexp.Regexp{
						regexp.MustCompile(`XXX`),
						regexp.MustCompile(`YYY`),
						regexp.MustCompile(`ZZZ`),
					},
				},
				target: "un-match",
			},
			&action{
				matched: false,
			},
		),
		gen(
			"not match / no patterns",
			&condition{
				matcher: &stringMatcher{
					patterns: []string{},
				},
				target: "", // Empty string.
			},
			&action{
				matched: false,
			},
		),
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			tt.C.matcher.match = tt.C.matcher.regex
			matched := tt.C.matcher.Match(tt.C.target)
			testutil.Diff(t, tt.A.matched, matched)
		})
	}
}

func TestNewBytesMatchers(t *testing.T) {
	type condition struct {
		specs []*k.MatcherSpec
	}

	type action struct {
		match   []string
		unmatch []string
		err     error
	}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"no spec",
			&condition{
				specs: []*k.MatcherSpec{},
			},
			&action{
				match:   []string{},
				unmatch: []string{"", "test"},
				err:     nil,
			},
		),
		gen(
			"one valid spec",
			&condition{
				specs: []*k.MatcherSpec{
					{
						MatchType: k.MatchType_Exact,
						Patterns:  []string{"test1", "test2", "test3"},
					},
				},
			},
			&action{
				match:   []string{"test1", "test2", "test3"},
				unmatch: []string{"", "test", "test4"},
				err:     nil,
			},
		),
		gen(
			"multiple valid specs",
			&condition{
				specs: []*k.MatcherSpec{
					{
						MatchType: k.MatchType_Exact,
						Patterns:  []string{"test1", "test2", "test3"},
					},
					{
						MatchType: k.MatchType_Prefix,
						Patterns:  []string{"testX", "testY", "testZ"},
					},
				},
			},
			&action{
				match:   []string{"testXXX", "testYYY", "testZZZ"},
				unmatch: []string{"", "test"},
				err:     nil,
			},
		),
		gen(
			"input nil spec",
			&condition{
				specs: []*k.MatcherSpec{
					nil, // This should be ignored.
					{
						MatchType: k.MatchType_Exact,
						Patterns:  []string{"test1", "test2", "test3"},
					},
					{
						MatchType: k.MatchType_Prefix,
						Patterns:  []string{"testX", "testY", "testZ"},
					},
				},
			},
			&action{
				match:   []string{"testXXX", "testYYY", "testZZZ"},
				unmatch: []string{"", "test"},
				err:     nil,
			},
		),
		gen(
			"unsupported type",
			&condition{
				specs: []*k.MatcherSpec{
					{
						MatchType: k.MatchType(999), // 999 does not exist.
						Patterns:  []string{},
					},
				},
			},
			&action{
				err: &er.Error{
					Package:     ErrPkg,
					Type:        ErrTypeMatcher,
					Description: ErrDscUnsupported,
				},
			},
		),
		gen(
			"regex error",
			&condition{
				specs: []*k.MatcherSpec{
					{
						MatchType: k.MatchType_Regex,
						Patterns:  []string{"[0-9a-"},
					},
				},
			},
			&action{
				err: &er.Error{
					Package:     ErrPkg,
					Type:        ErrTypeMatcher,
					Description: ErrDscPattern,
				},
			},
		),
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			ms, err := NewBytesMatchers(tt.C.specs...)
			testutil.Diff(t, tt.A.err, err, cmpopts.EquateErrors())

			for _, v := range tt.A.match {
				matched := false
				for _, m := range ms {
					if matched = m.Match([]byte(v)); matched {
						break
					}
				}
				t.Log("Check match:", v)
				testutil.Diff(t, true, matched)
			}

			for _, v := range tt.A.unmatch {
				matched := false
				for _, m := range ms {
					if matched = m.Match([]byte(v)); matched {
						break
					}
				}
				t.Log("Check unmatch:", v)
				testutil.Diff(t, false, matched)
			}
		})
	}
}

func TestNewBytesMatcher(t *testing.T) {
	type condition struct {
		typ      MatchType
		patterns []string
	}

	type action struct {
		match   []string
		unmatch []string
		err     error
	}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"exact",
			&condition{
				typ:      MatchTypeExact,
				patterns: []string{"test1", "test2", "test3"},
			},
			&action{
				match:   []string{"test1", "test2", "test3"},
				unmatch: []string{"", "test", "test4"},
			},
		),
		gen(
			"prefix",
			&condition{
				typ:      MatchTypePrefix,
				patterns: []string{"testX", "testY", "testZ"},
			},
			&action{
				match:   []string{"testXXX", "testYYY", "testZZZ"},
				unmatch: []string{"", "test"},
			},
		),
		gen(
			"suffix",
			&condition{
				typ:      MatchTypeSuffix,
				patterns: []string{"XXX", "YYY", "ZZZ"},
			},
			&action{
				match:   []string{"testXXX", "testYYY", "testZZZ"},
				unmatch: []string{"", "test", "testX", "testY", "testZ"},
			},
		),
		gen(
			"contains",
			&condition{
				typ:      MatchTypeContains,
				patterns: []string{"XXX", "YYY", "ZZZ"},
			},
			&action{
				match:   []string{"testXXXtest", "testYYYtest", "testZZZtest"},
				unmatch: []string{"", "test", "X", "Y", "Z"},
			},
		),
		gen(
			"path",
			&condition{
				typ:      MatchTypePath,
				patterns: []string{"/foo/*", "/bar/*"},
			},
			&action{
				match:   []string{"/foo/", "/foo/test", "/bar/", "/bar/test"},
				unmatch: []string{"", "/foo", "/foo*", "/foo/test1/test2"},
			},
		),
		gen(
			"filepath",
			&condition{
				typ:      MatchTypeFilePath,
				patterns: []string{"/foo/*", "/bar/*"},
			},
			&action{
				match:   []string{"/foo/", "/foo/test", "/bar/", "/bar/test"},
				unmatch: []string{"", "/foo", "/foo*"},
			},
		),
		gen(
			"regex",
			&condition{
				typ:      MatchTypeRegex,
				patterns: []string{"test", "^foo[0-9a-zA-Z]*$"},
			},
			&action{
				match:   []string{"test", "foo", "foobar"},
				unmatch: []string{"", "t", "foo#"},
			},
		),
		gen(
			"regex POSIX",
			&condition{
				typ:      MatchTypeRegexPOSIX,
				patterns: []string{"test", "^foo[0-9a-zA-Z]*$"},
			},
			&action{
				match:   []string{"test", "foo", "foobar"},
				unmatch: []string{"", "t", "foo#"},
			},
		),
		gen(
			"unsupported type",
			&condition{
				typ:      MatchType(999), // Does not exist.
				patterns: []string{""},
			},
			&action{
				err: &er.Error{
					Package:     ErrPkg,
					Type:        ErrTypeMatcher,
					Description: ErrDscUnsupported,
				},
			},
		),
		gen(
			"path error",
			&condition{
				typ:      MatchTypePath,
				patterns: []string{"[0-9a-"},
			},
			&action{
				err: &er.Error{
					Package:     ErrPkg,
					Type:        ErrTypeMatcher,
					Description: ErrDscPattern,
				},
			},
		),
		gen(
			"filepath error",
			&condition{
				typ:      MatchTypeFilePath,
				patterns: []string{"[0-9a-"},
			},
			&action{
				err: &er.Error{
					Package:     ErrPkg,
					Type:        ErrTypeMatcher,
					Description: ErrDscPattern,
				},
			},
		),
		gen(
			"regex error",
			&condition{
				typ:      MatchTypeRegex,
				patterns: []string{"[0-9a-"},
			},
			&action{
				err: &er.Error{
					Package:     ErrPkg,
					Type:        ErrTypeMatcher,
					Description: ErrDscPattern,
				},
			},
		),
		gen(
			"regex POSIX error",
			&condition{
				typ:      MatchTypeRegexPOSIX,
				patterns: []string{"[0-9a-"},
			},
			&action{
				err: &er.Error{
					Package:     ErrPkg,
					Type:        ErrTypeMatcher,
					Description: ErrDscPattern,
				},
			},
		),
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			m, err := NewBytesMatcher(tt.C.typ, tt.C.patterns...)
			testutil.Diff(t, tt.A.err, err, cmpopts.EquateErrors())
			if err != nil {
				return
			}

			for _, v := range tt.A.match {
				t.Log("Check match:", v)
				testutil.Diff(t, true, m.Match([]byte(v)))
			}

			for _, v := range tt.A.unmatch {
				t.Log("Check unmatch:", v)
				testutil.Diff(t, false, m.Match([]byte(v)))
			}
		})
	}
}
