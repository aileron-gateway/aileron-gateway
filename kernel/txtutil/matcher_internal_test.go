// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package txtutil

import (
	"regexp"
	"testing"

	k "github.com/aileron-gateway/aileron-gateway/apis/kernel"
	"github.com/aileron-gateway/aileron-gateway/kernel/er"
	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
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

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	cndValidPatterns := tb.Condition("valid patterns", "input valid patterns")
	cndInvalidPattern := tb.Condition("invalid pattern", "input an invalid pattern")
	cndUnsupportedType := tb.Condition("unsupported type", "input unsupported MatchType")
	actCheckMatched := tb.Action("matched", "check matched")
	actCheckUnmatched := tb.Action("un-match", "check not matched")
	actCheckNoError := tb.Action("no error", "check not matched target")
	actCheckError := tb.Action("non-nil error", "check not matched target")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"exact",
			[]string{cndValidPatterns},
			[]string{actCheckMatched, actCheckUnmatched, actCheckNoError},
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
			[]string{cndValidPatterns},
			[]string{actCheckMatched, actCheckUnmatched, actCheckNoError},
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
			[]string{cndValidPatterns},
			[]string{actCheckMatched, actCheckUnmatched, actCheckNoError},
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
			[]string{cndValidPatterns},
			[]string{actCheckMatched, actCheckUnmatched, actCheckNoError},
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
			[]string{cndValidPatterns},
			[]string{actCheckMatched, actCheckUnmatched, actCheckNoError},
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
			[]string{cndValidPatterns},
			[]string{actCheckMatched, actCheckUnmatched, actCheckNoError},
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
			[]string{cndValidPatterns},
			[]string{actCheckMatched, actCheckUnmatched, actCheckNoError},
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
			[]string{cndValidPatterns},
			[]string{actCheckMatched, actCheckUnmatched, actCheckNoError},
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
			[]string{cndUnsupportedType},
			[]string{actCheckError},
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
			[]string{cndInvalidPattern},
			[]string{actCheckError},
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
			[]string{cndInvalidPattern},
			[]string{actCheckError},
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
			[]string{cndInvalidPattern},
			[]string{actCheckError},
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
			[]string{cndInvalidPattern},
			[]string{actCheckError},
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

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			m, err := NewStringMatcher(tt.C().typ, tt.C().patterns...)
			testutil.Diff(t, tt.A().err, err, cmpopts.EquateErrors())
			if err != nil {
				return
			}

			for _, v := range tt.A().match {
				t.Log("Check match:", v)
				testutil.Diff(t, true, m.Match(v))
			}

			for _, v := range tt.A().unmatch {
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

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	cndNilSpec := tb.Condition("nil spec", "input nil spec")
	cndValidSpec := tb.Condition("valid spec", "input valid specs")
	actCheckError := tb.Action("error", "check that the expected error is returned")
	actCheckNoError := tb.Action("no error", "check that the there is no error")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"no spec",
			[]string{},
			[]string{actCheckNoError},
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
			[]string{cndValidSpec},
			[]string{actCheckNoError},
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
			[]string{cndValidSpec},
			[]string{actCheckNoError},
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
			[]string{cndNilSpec},
			[]string{actCheckNoError},
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
			[]string{},
			[]string{actCheckError},
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
			[]string{},
			[]string{actCheckError},
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

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			ms, err := NewStringMatchers(tt.C().specs...)
			testutil.Diff(t, tt.A().err, err, cmpopts.EquateErrors())

			for _, v := range tt.A().match {
				matched := false
				for _, m := range ms {
					if matched = m.Match(v); matched {
						break
					}
				}
				t.Log("Check match:", v)
				testutil.Diff(t, true, matched)
			}

			for _, v := range tt.A().unmatch {
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

	cndNoPatterns := "no patterns"
	cndOnePattern := "1 pattern"
	cndMultiplePattern := "multiple patterns"
	actCheckMatched := "matched"
	actCheckUnmatched := "un-match"

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	tb.Condition(cndNoPatterns, "no patterns are set")
	tb.Condition(cndOnePattern, "1 patten is set")
	tb.Condition(cndMultiplePattern, "multiple patterns are set")
	tb.Action(actCheckMatched, "check matched")
	tb.Action(actCheckUnmatched, "check not matched")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"match / 1 pattern",
			[]string{cndOnePattern},
			[]string{actCheckMatched},
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
			[]string{cndMultiplePattern},
			[]string{actCheckUnmatched},
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
			[]string{cndOnePattern},
			[]string{actCheckUnmatched},
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
			[]string{cndMultiplePattern},
			[]string{actCheckUnmatched},
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
			[]string{cndNoPatterns},
			[]string{actCheckUnmatched},
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

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			tt.C().matcher.match = tt.C().matcher.exact
			matched := tt.C().matcher.Match(tt.C().target)
			testutil.Diff(t, tt.A().matched, matched)
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

	cndNoPatterns := "no patterns"
	cndOnePattern := "1 pattern"
	cndMultiplePattern := "multiple patterns"
	actCheckMatched := "matched"
	actCheckUnmatched := "un-match"

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	tb.Condition(cndNoPatterns, "no patterns are set")
	tb.Condition(cndOnePattern, "1 patten is set")
	tb.Condition(cndMultiplePattern, "multiple patterns are set")
	tb.Action(actCheckMatched, "check matched")
	tb.Action(actCheckUnmatched, "check not matched")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"match / 1 pattern",
			[]string{cndOnePattern},
			[]string{actCheckMatched},
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
			[]string{cndMultiplePattern},
			[]string{actCheckUnmatched},
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
			[]string{cndOnePattern},
			[]string{actCheckUnmatched},
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
			[]string{cndMultiplePattern},
			[]string{actCheckUnmatched},
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
			[]string{cndNoPatterns},
			[]string{actCheckUnmatched},
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

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			tt.C().matcher.match = tt.C().matcher.prefix
			matched := tt.C().matcher.Match(tt.C().target)
			testutil.Diff(t, tt.A().matched, matched)
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

	cndNoPatterns := "no patterns"
	cndOnePattern := "1 pattern"
	cndMultiplePattern := "multiple patterns"
	actCheckMatched := "matched"
	actCheckUnmatched := "un-match"

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	tb.Condition(cndNoPatterns, "no patterns are set")
	tb.Condition(cndOnePattern, "1 patten is set")
	tb.Condition(cndMultiplePattern, "multiple patterns are set")
	tb.Action(actCheckMatched, "check matched")
	tb.Action(actCheckUnmatched, "check not matched")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"match / 1 pattern",
			[]string{cndOnePattern},
			[]string{actCheckMatched},
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
			[]string{cndMultiplePattern},
			[]string{actCheckUnmatched},
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
			[]string{cndOnePattern},
			[]string{actCheckUnmatched},
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
			[]string{cndMultiplePattern},
			[]string{actCheckUnmatched},
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
			[]string{cndNoPatterns},
			[]string{actCheckUnmatched},
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

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			tt.C().matcher.match = tt.C().matcher.suffix
			matched := tt.C().matcher.Match(tt.C().target)
			testutil.Diff(t, tt.A().matched, matched)
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

	cndNoPatterns := "no patterns"
	cndOnePattern := "1 pattern"
	cndMultiplePattern := "multiple patterns"
	actCheckMatched := "matched"
	actCheckUnmatched := "un-match"

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	tb.Condition(cndNoPatterns, "no patterns are set")
	tb.Condition(cndOnePattern, "1 patten is set")
	tb.Condition(cndMultiplePattern, "multiple patterns are set")
	tb.Action(actCheckMatched, "check matched")
	tb.Action(actCheckUnmatched, "check not matched")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"match / 1 pattern",
			[]string{cndOnePattern},
			[]string{actCheckMatched},
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
			[]string{cndMultiplePattern},
			[]string{actCheckUnmatched},
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
			[]string{cndOnePattern},
			[]string{actCheckUnmatched},
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
			[]string{cndMultiplePattern},
			[]string{actCheckUnmatched},
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
			[]string{cndNoPatterns},
			[]string{actCheckUnmatched},
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

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			tt.C().matcher.match = tt.C().matcher.contain
			matched := tt.C().matcher.Match(tt.C().target)
			testutil.Diff(t, tt.A().matched, matched)
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

	cndNoPatterns := "no patterns"
	cndOnePattern := "1 pattern"
	cndMultiplePattern := "multiple patterns"
	actCheckMatched := "matched"
	actCheckUnmatched := "un-match"

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	tb.Condition(cndNoPatterns, "no patterns are set")
	tb.Condition(cndOnePattern, "1 patten is set")
	tb.Condition(cndMultiplePattern, "multiple patterns are set")
	tb.Action(actCheckMatched, "check matched")
	tb.Action(actCheckUnmatched, "check not matched")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"match / 1 pattern",
			[]string{cndOnePattern},
			[]string{actCheckMatched},
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
			[]string{cndMultiplePattern},
			[]string{actCheckUnmatched},
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
			[]string{cndOnePattern},
			[]string{actCheckUnmatched},
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
			[]string{cndMultiplePattern},
			[]string{actCheckUnmatched},
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
			[]string{cndNoPatterns},
			[]string{actCheckUnmatched},
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

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			tt.C().matcher.match = tt.C().matcher.path
			matched := tt.C().matcher.Match(tt.C().target)
			testutil.Diff(t, tt.A().matched, matched)
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

	cndNoPatterns := "no patterns"
	cndOnePattern := "1 pattern"
	cndMultiplePattern := "multiple patterns"
	actCheckMatched := "matched"
	actCheckUnmatched := "un-match"

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	tb.Condition(cndNoPatterns, "no patterns are set")
	tb.Condition(cndOnePattern, "1 patten is set")
	tb.Condition(cndMultiplePattern, "multiple patterns are set")
	tb.Action(actCheckMatched, "check matched")
	tb.Action(actCheckUnmatched, "check not matched")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"match / 1 pattern",
			[]string{cndOnePattern},
			[]string{actCheckMatched},
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
			[]string{cndMultiplePattern},
			[]string{actCheckUnmatched},
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
			[]string{cndOnePattern},
			[]string{actCheckUnmatched},
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
			[]string{cndMultiplePattern},
			[]string{actCheckUnmatched},
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
			[]string{cndNoPatterns},
			[]string{actCheckUnmatched},
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

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			tt.C().matcher.match = tt.C().matcher.regex
			matched := tt.C().matcher.Match(tt.C().target)
			testutil.Diff(t, tt.A().matched, matched)
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

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	cndNilSpec := tb.Condition("nil spec", "input nil spec")
	cndValidSpec := tb.Condition("valid spec", "input valid specs")
	actCheckError := tb.Action("error", "check that the expected error is returned")
	actCheckNoError := tb.Action("no error", "check that the there is no error")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"no spec",
			[]string{},
			[]string{actCheckNoError},
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
			[]string{cndValidSpec},
			[]string{actCheckNoError},
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
			[]string{cndValidSpec},
			[]string{actCheckNoError},
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
			[]string{cndNilSpec},
			[]string{actCheckNoError},
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
			[]string{},
			[]string{actCheckError},
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
			[]string{},
			[]string{actCheckError},
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

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			ms, err := NewBytesMatchers(tt.C().specs...)
			testutil.Diff(t, tt.A().err, err, cmpopts.EquateErrors())

			for _, v := range tt.A().match {
				matched := false
				for _, m := range ms {
					if matched = m.Match([]byte(v)); matched {
						break
					}
				}
				t.Log("Check match:", v)
				testutil.Diff(t, true, matched)
			}

			for _, v := range tt.A().unmatch {
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

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	cndValidPatterns := tb.Condition("valid patterns", "input valid patterns")
	cndInvalidPattern := tb.Condition("invalid pattern", "input an invalid pattern")
	cndUnsupportedType := tb.Condition("unsupported type", "input unsupported MatchType")
	actCheckMatched := tb.Action("matched", "check matched")
	actCheckUnmatched := tb.Action("un-match", "check not matched")
	actCheckNoError := tb.Action("no error", "check not matched target")
	actCheckError := tb.Action("non-nil error", "check not matched target")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"exact",
			[]string{cndValidPatterns},
			[]string{actCheckMatched, actCheckUnmatched, actCheckNoError},
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
			[]string{cndValidPatterns},
			[]string{actCheckMatched, actCheckUnmatched, actCheckNoError},
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
			[]string{cndValidPatterns},
			[]string{actCheckMatched, actCheckUnmatched, actCheckNoError},
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
			[]string{cndValidPatterns},
			[]string{actCheckMatched, actCheckUnmatched, actCheckNoError},
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
			[]string{cndValidPatterns},
			[]string{actCheckMatched, actCheckUnmatched, actCheckNoError},
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
			[]string{cndValidPatterns},
			[]string{actCheckMatched, actCheckUnmatched, actCheckNoError},
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
			[]string{cndValidPatterns},
			[]string{actCheckMatched, actCheckUnmatched, actCheckNoError},
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
			[]string{cndValidPatterns},
			[]string{actCheckMatched, actCheckUnmatched, actCheckNoError},
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
			[]string{cndUnsupportedType},
			[]string{actCheckError},
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
			[]string{cndInvalidPattern},
			[]string{actCheckError},
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
			[]string{cndInvalidPattern},
			[]string{actCheckError},
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
			[]string{cndInvalidPattern},
			[]string{actCheckError},
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
			[]string{cndInvalidPattern},
			[]string{actCheckError},
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

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			m, err := NewBytesMatcher(tt.C().typ, tt.C().patterns...)
			testutil.Diff(t, tt.A().err, err, cmpopts.EquateErrors())
			if err != nil {
				return
			}

			for _, v := range tt.A().match {
				t.Log("Check match:", v)
				testutil.Diff(t, true, m.Match([]byte(v)))
			}

			for _, v := range tt.A().unmatch {
				t.Log("Check unmatch:", v)
				testutil.Diff(t, false, m.Match([]byte(v)))
			}
		})
	}
}
