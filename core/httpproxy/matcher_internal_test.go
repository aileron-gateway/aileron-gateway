// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package httpproxy

import (
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"

	v1 "github.com/aileron-gateway/aileron-gateway/apis/core/v1"
	k "github.com/aileron-gateway/aileron-gateway/apis/kernel"
	"github.com/aileron-gateway/aileron-gateway/core"
	"github.com/aileron-gateway/aileron-gateway/kernel/er"
	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
	"github.com/aileron-gateway/aileron-gateway/kernel/txtutil"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestNewMatcher(t *testing.T) {
	type condition struct {
		spec *v1.PathMatcherSpec
	}

	type action struct {
		matchPaths    []string
		notMatchPaths []string
		err           any // error or errorutil.Kind
		errPattern    *regexp.Regexp
	}

	cndExactMatcher := "exact matcher"
	cndPrefixMatcher := "prefix matcher"
	cndSuffixMatcher := "suffix matcher"
	cndContainMatcher := "contain matcher"
	cndPathMatcher := "path matcher"
	cndFilePathMatcher := "path matcher"
	cndRegexMatcher := "regex matcher"
	cndRegexPOSIXMatcher := "regex POSIX matcher"
	cndWithTrimPrefix := "with trim prefix"
	cndWithAppendPrefix := "with append prefix"
	actCheckError := "check error"
	actCheckNoError := "check no error"

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	tb.Condition(cndExactMatcher, "create a exact matcher")
	tb.Condition(cndPrefixMatcher, "create a new prefix matcher")
	tb.Condition(cndSuffixMatcher, "create a new suffix matcher")
	tb.Condition(cndContainMatcher, "create a new contain matcher")
	tb.Condition(cndPathMatcher, "create a new path matcher")
	tb.Condition(cndFilePathMatcher, "create a new filePath matcher")
	tb.Condition(cndRegexMatcher, "create a new regex matcher")
	tb.Condition(cndRegexPOSIXMatcher, "create a new regex POSIX matcher")
	tb.Condition(cndWithTrimPrefix, "specify trim prefix")
	tb.Condition(cndWithAppendPrefix, "specify append prefix")
	tb.Action(actCheckError, "check that there is an error returned")
	tb.Action(actCheckNoError, "check that there is no error returned")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"default will be prefix matcher",
			[]string{},
			[]string{actCheckNoError},
			&condition{
				spec: &v1.PathMatcherSpec{
					Match:     "/test",
					MatchType: k.MatchType(-1),
				},
			},
			&action{
				matchPaths:    []string{"/test", "/test/path"},
				notMatchPaths: []string{"", "/", "/te"},
			},
		),
		gen(
			"exact matcher",
			[]string{cndExactMatcher},
			[]string{actCheckNoError},
			&condition{
				spec: &v1.PathMatcherSpec{
					Match:     "/test",
					MatchType: k.MatchType_Exact,
				},
			},
			&action{
				matchPaths:    []string{"/test"},
				notMatchPaths: []string{"", "/", "/te", "/test/path"},
			},
		),
		gen(
			"exact matcher with trim",
			[]string{cndExactMatcher, cndWithTrimPrefix},
			[]string{actCheckNoError},
			&condition{
				spec: &v1.PathMatcherSpec{
					Match:      "/test",
					MatchType:  k.MatchType_Exact,
					TrimPrefix: "/trim/me",
				},
			},
			&action{
				matchPaths: []string{"/test", "/trim/me/test"},
				notMatchPaths: []string{"", "/", "/te", "/test/path",
					"/trim/me", "/trim/me/", "/trim/me/te", "/trim/me/test/path"},
			},
		),
		gen(
			"exact matcher with append",
			[]string{cndExactMatcher, cndWithTrimPrefix},
			[]string{actCheckNoError},
			&condition{
				spec: &v1.PathMatcherSpec{
					Match:        "/test",
					MatchType:    k.MatchType_Exact,
					AppendPrefix: "/append/me",
				},
			},
			&action{
				matchPaths:    []string{"/test"},
				notMatchPaths: []string{"", "/", "/te", "/test/path"},
			},
		),
		gen(
			"prefix matcher",
			[]string{cndPrefixMatcher},
			[]string{actCheckNoError},
			&condition{
				spec: &v1.PathMatcherSpec{
					Match:     "/test",
					MatchType: k.MatchType_Prefix,
				},
			},
			&action{
				matchPaths:    []string{"/test", "/test/path"},
				notMatchPaths: []string{"", "/", "/te"},
			},
		),
		gen(
			"prefix matcher with trim",
			[]string{cndPrefixMatcher, cndWithTrimPrefix},
			[]string{actCheckNoError},
			&condition{
				spec: &v1.PathMatcherSpec{
					Match:      "/test",
					MatchType:  k.MatchType_Prefix,
					TrimPrefix: "/trim/me",
				},
			},
			&action{
				matchPaths: []string{"/test", "/test/path", "/trim/me/test", "/trim/me/test/path"},
				notMatchPaths: []string{"", "/", "/te",
					"/trim/me", "/trim/me/", "/trim/me/te"},
			},
		),
		gen(
			"prefix matcher with append",
			[]string{cndPrefixMatcher},
			[]string{actCheckNoError},
			&condition{
				spec: &v1.PathMatcherSpec{
					Match:        "/test",
					MatchType:    k.MatchType_Prefix,
					AppendPrefix: "/append/me",
				},
			},
			&action{
				matchPaths:    []string{"/test", "/test/path"},
				notMatchPaths: []string{"", "/", "/te"},
			},
		),
		gen(
			"suffix matcher",
			[]string{cndSuffixMatcher},
			[]string{actCheckNoError},
			&condition{
				spec: &v1.PathMatcherSpec{
					Match:     "/test",
					MatchType: k.MatchType_Suffix,
				},
			},
			&action{
				matchPaths:    []string{"/test", "/path/test"},
				notMatchPaths: []string{"", "/", "/test/path"},
			},
		),
		gen(
			"suffix matcher with trim",
			[]string{cndSuffixMatcher, cndWithTrimPrefix},
			[]string{actCheckNoError},
			&condition{
				spec: &v1.PathMatcherSpec{
					Match:      "/test",
					MatchType:  k.MatchType_Suffix,
					TrimPrefix: "/trim/me",
				},
			},
			&action{
				matchPaths: []string{"/test", "/path/test", "/trim/me/test", "/trim/me/path/test"},
				notMatchPaths: []string{"", "/", "/test/path",
					"/trim/me", "/trim/me/", "/trim/me/test/path"},
			},
		),
		gen(
			"suffix matcher with append",
			[]string{cndSuffixMatcher},
			[]string{actCheckNoError},
			&condition{
				spec: &v1.PathMatcherSpec{
					Match:        "/test",
					MatchType:    k.MatchType_Suffix,
					AppendPrefix: "/append/me",
				},
			},
			&action{
				matchPaths:    []string{"/test", "/path/test"},
				notMatchPaths: []string{"", "/", "/test/path"},
			},
		),
		gen(
			"contain matcher",
			[]string{cndContainMatcher},
			[]string{actCheckNoError},
			&condition{
				spec: &v1.PathMatcherSpec{
					Match:     "/test",
					MatchType: k.MatchType_Contains,
				},
			},
			&action{
				matchPaths:    []string{"/test", "/test/path", "/foo/test/bar"},
				notMatchPaths: []string{"", "/", "test", "/path"},
			},
		),
		gen(
			"contain matcher with trim",
			[]string{cndContainMatcher, cndWithTrimPrefix},
			[]string{actCheckNoError},
			&condition{
				spec: &v1.PathMatcherSpec{
					Match:      "/test",
					MatchType:  k.MatchType_Contains,
					TrimPrefix: "/trim/me",
				},
			},
			&action{
				matchPaths: []string{"/test", "/test/path", "/foo/test/bar",
					"/trim/me/test", "/trim/me/test/path", "/trim/me/foo/test/bar"},
				notMatchPaths: []string{"", "/", "test", "/path",
					"/trim/me", "/trim/me/", "/trim/me/path"},
			},
		),
		gen(
			"contain matcher with append",
			[]string{cndContainMatcher},
			[]string{actCheckNoError},
			&condition{
				spec: &v1.PathMatcherSpec{
					Match:        "/test",
					MatchType:    k.MatchType_Contains,
					AppendPrefix: "/append/me",
				},
			},
			&action{
				matchPaths:    []string{"/test", "/test/path", "/foo/test/bar"},
				notMatchPaths: []string{"", "/", "test", "/path"},
			},
		),
		gen(
			"path matcher",
			[]string{cndPathMatcher},
			[]string{actCheckNoError},
			&condition{
				spec: &v1.PathMatcherSpec{
					Match:     "/*/test",
					MatchType: k.MatchType_Path,
				},
			},
			&action{
				matchPaths:    []string{"/foo/test", "/bar/test"},
				notMatchPaths: []string{"", "/", "/test", "/test/baz"},
			},
		),
		gen(
			"path matcher with trim",
			[]string{cndPathMatcher, cndWithTrimPrefix},
			[]string{actCheckNoError},
			&condition{
				spec: &v1.PathMatcherSpec{
					Match:      "/*/test",
					MatchType:  k.MatchType_Path,
					TrimPrefix: "/trim/me",
				},
			},
			&action{
				matchPaths: []string{"/foo/test", "/bar/test"},
				notMatchPaths: []string{"", "/", "/test", "/test/baz", "/foo/test/",
					"/trim/me", "/trim/me/", "/trim/me/test", "/trim/me/test/baz", "/trim/me/foo/test/"},
			},
		),
		gen(
			"path matcher with append",
			[]string{cndPathMatcher},
			[]string{actCheckNoError},
			&condition{
				spec: &v1.PathMatcherSpec{
					Match:        "/*/test",
					MatchType:    k.MatchType_Path,
					AppendPrefix: "/append/me",
				},
			},
			&action{
				matchPaths:    []string{"/foo/test", "/bar/test"},
				notMatchPaths: []string{"", "/", "/test", "/test/baz"},
			},
		),
		gen(
			"filepath matcher",
			[]string{cndFilePathMatcher},
			[]string{actCheckNoError},
			&condition{
				spec: &v1.PathMatcherSpec{
					Match:     "/*/test",
					MatchType: k.MatchType_FilePath,
				},
			},
			&action{
				matchPaths:    []string{"/foo/test", "/bar/test"},
				notMatchPaths: []string{"", "/", "/test", "/test/baz"},
			},
		),
		gen(
			"filepath matcher with trim",
			[]string{cndFilePathMatcher, cndWithTrimPrefix},
			[]string{actCheckNoError},
			&condition{
				spec: &v1.PathMatcherSpec{
					Match:      "/*/test",
					MatchType:  k.MatchType_FilePath,
					TrimPrefix: "/trim/me",
				},
			},
			&action{
				matchPaths: []string{"/foo/test", "/bar/test"},
				notMatchPaths: []string{"", "/", "/test", "/test/baz", "/foo/test/",
					"/trim/me", "/trim/me/", "/trim/me/test", "/trim/me/test/baz", "/trim/me/foo/test/"},
			},
		),
		gen(
			"filepath matcher with append",
			[]string{cndFilePathMatcher},
			[]string{actCheckNoError},
			&condition{
				spec: &v1.PathMatcherSpec{
					Match:        "/*/test",
					MatchType:    k.MatchType_FilePath,
					AppendPrefix: "/append/me",
				},
			},
			&action{
				matchPaths:    []string{"/foo/test", "/bar/test"},
				notMatchPaths: []string{"", "/", "/test", "/test/baz"},
			},
		),
		gen(
			"regex matcher",
			[]string{cndRegexMatcher},
			[]string{actCheckNoError},
			&condition{
				spec: &v1.PathMatcherSpec{
					Match:     "^/test/(foo|bar)$",
					MatchType: k.MatchType_Regex,
				},
			},
			&action{
				matchPaths:    []string{"/test/foo", "/test/bar"},
				notMatchPaths: []string{"", "/", "/test", "/test/baz"},
			},
		),
		gen(
			"regex matcher with trim",
			[]string{cndRegexMatcher, cndWithTrimPrefix},
			[]string{actCheckNoError},
			&condition{
				spec: &v1.PathMatcherSpec{
					Match:      "^/test/(foo|bar)$",
					MatchType:  k.MatchType_Regex,
					TrimPrefix: "/trim/me",
				},
			},
			&action{
				matchPaths: []string{"/test/foo", "/test/bar"},
				notMatchPaths: []string{"", "/", "/test", "/test/baz",
					"/trim/me", "/trim/me/", "/trim/me/test", "/trim/me/test/baz"},
			},
		),
		gen(
			"regex matcher with append",
			[]string{cndRegexMatcher},
			[]string{actCheckNoError},
			&condition{
				spec: &v1.PathMatcherSpec{
					Match:        "^/test/(foo|bar)$",
					MatchType:    k.MatchType_Regex,
					AppendPrefix: "/append/me",
				},
			},
			&action{
				matchPaths:    []string{"/test/foo", "/test/bar"},
				notMatchPaths: []string{"", "/", "/test", "/test/baz"},
			},
		),
		gen(
			"POSIX matcher",
			[]string{cndRegexPOSIXMatcher},
			[]string{actCheckNoError},
			&condition{
				spec: &v1.PathMatcherSpec{
					Match:     "^/test/(foo|bar)$",
					MatchType: k.MatchType_RegexPOSIX,
				},
			},
			&action{
				matchPaths:    []string{"/test/foo", "/test/bar"},
				notMatchPaths: []string{"", "/", "/test", "/test/baz"},
			},
		),
		gen(
			"POSIX matcher with trim",
			[]string{cndRegexPOSIXMatcher, cndWithTrimPrefix},
			[]string{actCheckNoError},
			&condition{
				spec: &v1.PathMatcherSpec{
					Match:      "^/test/(foo|bar)$",
					MatchType:  k.MatchType_RegexPOSIX,
					TrimPrefix: "/trim/me",
				},
			},
			&action{
				matchPaths: []string{"/test/foo", "/test/bar"},
				notMatchPaths: []string{"", "/", "/test", "/test/baz",
					"/trim/me", "/trim/me/", "/trim/me/test", "/trim/me/test/baz"},
			},
		),
		gen(
			"POSIX matcher with append",
			[]string{cndRegexPOSIXMatcher},
			[]string{actCheckNoError},
			&condition{
				spec: &v1.PathMatcherSpec{
					Match:        "^/test/(foo|bar)$",
					MatchType:    k.MatchType_RegexPOSIX,
					AppendPrefix: "/append/me",
				},
			},
			&action{
				matchPaths:    []string{"/test/foo", "/test/bar"},
				notMatchPaths: []string{"", "/", "/test", "/test/baz"},
			},
		),
		gen(
			"error validating path pattern",
			[]string{cndPathMatcher},
			[]string{actCheckError},
			&condition{
				spec: &v1.PathMatcherSpec{
					Match:     "[0-9a-",
					MatchType: k.MatchType_Path,
				},
			},
			&action{
				err:        core.ErrCoreGenCreateComponent,
				errPattern: regexp.MustCompile(core.ErrPrefix + `failed to create component. invalid path pattern for 'Path' type`),
			},
		),
		gen(
			"error validating filePath pattern",
			[]string{cndFilePathMatcher},
			[]string{actCheckError},
			&condition{
				spec: &v1.PathMatcherSpec{
					Match:     "[0-9a-",
					MatchType: k.MatchType_FilePath,
				},
			},
			&action{
				err:        core.ErrCoreGenCreateComponent,
				errPattern: regexp.MustCompile(core.ErrPrefix + `failed to create component. invalid path pattern for 'FilePath' type`),
			},
		),
		gen(
			"error compiling regular expression",
			[]string{cndRegexMatcher},
			[]string{actCheckError},
			&condition{
				spec: &v1.PathMatcherSpec{
					Match:     "(?<",
					MatchType: k.MatchType_Regex,
				},
			},
			&action{
				err:        core.ErrCoreGenCreateComponent,
				errPattern: regexp.MustCompile(core.ErrPrefix + `failed to create component. invalid regular expression for 'Regex' type`),
			},
		),
		gen(
			"error compiling POSIX expression",
			[]string{cndRegexPOSIXMatcher},
			[]string{actCheckError},
			&condition{
				spec: &v1.PathMatcherSpec{
					Match:     "(?<",
					MatchType: k.MatchType_RegexPOSIX,
				},
			},
			&action{
				err:        core.ErrCoreGenCreateComponent,
				errPattern: regexp.MustCompile(core.ErrPrefix + `failed to create component. invalid regular expression for 'RegexPOSIX' type`),
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			mf, err := newMatcher(tt.C().spec)
			testutil.DiffError(t, tt.A().err, tt.A().errPattern, err)

			trimPrefix := tt.C().spec.TrimPrefix
			if trimPrefix == "" {
				trimPrefix = "MATCH_TO_NOTHING"
			}
			appendPrefix := tt.C().spec.AppendPrefix

			for _, p := range tt.A().matchPaths {
				path, matched := mf(p)
				testutil.Diff(t, true, matched)
				testutil.Diff(t, false, strings.HasPrefix(path, trimPrefix))
				testutil.Diff(t, true, strings.HasPrefix(path, appendPrefix))
			}

			for _, p := range tt.A().notMatchPaths {
				path, matched := mf(p)
				testutil.Diff(t, false, matched)
				testutil.Diff(t, false, strings.HasPrefix(path, trimPrefix))
				testutil.Diff(t, true, strings.HasPrefix(path, appendPrefix))
			}
		})
	}
}

func TestMatcher_exact(t *testing.T) {
	type condition struct {
		pattern string
		target  string
	}

	type action struct {
		path    string
		matched bool
	}

	cndEmptyPattern := "empty pattern"
	actCheckPath := "check returned path"
	actCheckMatched := "check error"
	actCheckNotMatched := "check no error"

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	tb.Condition(cndEmptyPattern, "set empty pattern")
	tb.Action(actCheckPath, "checked the returned path")
	tb.Action(actCheckMatched, "check that the pattern was matched ot the requested path")
	tb.Action(actCheckNotMatched, "check that the pattern was not matched to the requested path")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"empty string matched",
			[]string{cndEmptyPattern},
			[]string{actCheckPath, actCheckMatched},
			&condition{
				pattern: "",
				target:  "",
			},
			&action{
				path:    "",
				matched: true,
			},
		),
		gen(
			"empty string not matched",
			[]string{cndEmptyPattern},
			[]string{actCheckPath, actCheckNotMatched},
			&condition{
				pattern: "",
				target:  "/test",
			},
			&action{
				path:    "/test",
				matched: false,
			},
		),
		gen(
			"path matched",
			[]string{},
			[]string{actCheckPath, actCheckMatched},
			&condition{
				pattern: "/test",
				target:  "/test",
			},
			&action{
				path:    "/test",
				matched: true,
			},
		),
		gen(
			"path not matched",
			[]string{},
			[]string{actCheckPath, actCheckNotMatched},
			&condition{
				pattern: "/test",
				target:  "/",
			},
			&action{
				path:    "/",
				matched: false,
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			m := &matcher{
				pattern: tt.C().pattern,
			}

			path, matched := m.exact(tt.C().target)

			testutil.Diff(t, tt.A().matched, matched)
			testutil.Diff(t, tt.A().path, path)
		})
	}
}

func TestMatcher_prefix(t *testing.T) {
	type condition struct {
		pattern string
		target  string
	}

	type action struct {
		path    string
		matched bool
	}

	cndEmptyPattern := "empty pattern"
	actCheckPath := "check returned path"
	actCheckMatched := "check error"
	actCheckNotMatched := "check no error"

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	tb.Condition(cndEmptyPattern, "set empty pattern")
	tb.Action(actCheckPath, "checked the returned path")
	tb.Action(actCheckMatched, "check that the pattern was matched ot the requested path")
	tb.Action(actCheckNotMatched, "check that the pattern was not matched to the requested path")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"empty string matched",
			[]string{cndEmptyPattern},
			[]string{actCheckPath, actCheckMatched},
			&condition{
				pattern: "",
				target:  "",
			},
			&action{
				path:    "",
				matched: true,
			},
		),
		gen(
			"empty string is a refix for all path",
			[]string{cndEmptyPattern},
			[]string{actCheckPath, actCheckMatched},
			&condition{
				pattern: "",
				target:  "/test",
			},
			&action{
				path:    "/test",
				matched: true,
			},
		),
		gen(
			"exactly the same path",
			[]string{},
			[]string{actCheckPath, actCheckMatched},
			&condition{
				pattern: "/test",
				target:  "/test",
			},
			&action{
				path:    "/test",
				matched: true,
			},
		),
		gen(
			"prefix matched",
			[]string{},
			[]string{actCheckPath, actCheckMatched},
			&condition{
				pattern: "/test",
				target:  "/test/path",
			},
			&action{
				path:    "/test/path",
				matched: true,
			},
		),
		gen(
			"prefix not matched",
			[]string{},
			[]string{actCheckPath, actCheckNotMatched},
			&condition{
				pattern: "/test/path",
				target:  "/test",
			},
			&action{
				path:    "/test",
				matched: false,
			},
		),
		gen(
			"path not matched",
			[]string{},
			[]string{actCheckPath, actCheckNotMatched},
			&condition{
				pattern: "/test",
				target:  "/",
			},
			&action{
				path:    "/",
				matched: false,
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			m := &matcher{
				pattern: tt.C().pattern,
			}

			path, matched := m.prefix(tt.C().target)

			testutil.Diff(t, tt.A().matched, matched)
			testutil.Diff(t, tt.A().path, path)
		})
	}
}

func TestMatcher_suffix(t *testing.T) {
	type condition struct {
		pattern string
		target  string
	}

	type action struct {
		path    string
		matched bool
	}

	cndEmptyPattern := "empty pattern"
	actCheckPath := "check returned path"
	actCheckMatched := "check error"
	actCheckNotMatched := "check no error"

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	tb.Condition(cndEmptyPattern, "set empty pattern")
	tb.Action(actCheckPath, "checked the returned path")
	tb.Action(actCheckMatched, "check that the pattern was matched ot the requested path")
	tb.Action(actCheckNotMatched, "check that the pattern was not matched to the requested path")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"empty string matched",
			[]string{cndEmptyPattern},
			[]string{actCheckPath, actCheckMatched},
			&condition{
				pattern: "",
				target:  "",
			},
			&action{
				path:    "",
				matched: true,
			},
		),
		gen(
			"empty string is a suffix for all path",
			[]string{cndEmptyPattern},
			[]string{actCheckPath, actCheckMatched},
			&condition{
				pattern: "",
				target:  "/test",
			},
			&action{
				path:    "/test",
				matched: true,
			},
		),
		gen(
			"exactly the same path",
			[]string{},
			[]string{actCheckPath, actCheckMatched},
			&condition{
				pattern: "/test",
				target:  "/test",
			},
			&action{
				path:    "/test",
				matched: true,
			},
		),
		gen(
			"suffix matched",
			[]string{},
			[]string{actCheckPath, actCheckMatched},
			&condition{
				pattern: "/test",
				target:  "/path/test",
			},
			&action{
				path:    "/path/test",
				matched: true,
			},
		),
		gen(
			"suffix not matched",
			[]string{},
			[]string{actCheckPath, actCheckNotMatched},
			&condition{
				pattern: "/path/test",
				target:  "/test",
			},
			&action{
				path:    "/test",
				matched: false,
			},
		),
		gen(
			"path not matched",
			[]string{},
			[]string{actCheckPath, actCheckNotMatched},
			&condition{
				pattern: "/test",
				target:  "/",
			},
			&action{
				path:    "/",
				matched: false,
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			m := &matcher{
				pattern: tt.C().pattern,
			}

			path, matched := m.suffix(tt.C().target)

			testutil.Diff(t, tt.A().matched, matched)
			testutil.Diff(t, tt.A().path, path)
		})
	}
}

func TestMatcher_contains(t *testing.T) {
	type condition struct {
		pattern string
		target  string
	}

	type action struct {
		path    string
		matched bool
	}

	cndEmptyPattern := "empty pattern"
	actCheckPath := "check returned path"
	actCheckMatched := "check error"
	actCheckNotMatched := "check no error"

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	tb.Condition(cndEmptyPattern, "set empty pattern")
	tb.Action(actCheckPath, "checked the returned path")
	tb.Action(actCheckMatched, "check that the pattern was matched ot the requested path")
	tb.Action(actCheckNotMatched, "check that the pattern was not matched to the requested path")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"empty string matched",
			[]string{cndEmptyPattern},
			[]string{actCheckPath, actCheckMatched},
			&condition{
				pattern: "",
				target:  "",
			},
			&action{
				path:    "",
				matched: true,
			},
		),
		gen(
			"all path contains empty string",
			[]string{cndEmptyPattern},
			[]string{actCheckPath, actCheckMatched},
			&condition{
				pattern: "",
				target:  "/test",
			},
			&action{
				path:    "/test",
				matched: true,
			},
		),
		gen(
			"exactly the same path",
			[]string{},
			[]string{actCheckPath, actCheckMatched},
			&condition{
				pattern: "/test",
				target:  "/test",
			},
			&action{
				path:    "/test",
				matched: true,
			},
		),
		gen(
			"path contains pattern",
			[]string{},
			[]string{actCheckPath, actCheckMatched},
			&condition{
				pattern: "/test",
				target:  "/prefix/test/suffix",
			},
			&action{
				path:    "/prefix/test/suffix",
				matched: true,
			},
		),
		gen(
			"path not contains pattern",
			[]string{},
			[]string{actCheckPath, actCheckNotMatched},
			&condition{
				pattern: "/foo",
				target:  "/prefix/test/suffix",
			},
			&action{
				path:    "/prefix/test/suffix",
				matched: false,
			},
		),
		gen(
			"path not matched",
			[]string{},
			[]string{actCheckPath, actCheckNotMatched},
			&condition{
				pattern: "/test",
				target:  "/",
			},
			&action{
				path:    "/",
				matched: false,
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			m := &matcher{
				pattern: tt.C().pattern,
			}

			path, matched := m.contains(tt.C().target)

			testutil.Diff(t, tt.A().matched, matched)
			testutil.Diff(t, tt.A().path, path)
		})
	}
}

func TestMatcher_regEx(t *testing.T) {
	type condition struct {
		pattern string
		rewrite string
		target  string
	}

	type action struct {
		path    string
		matched bool
	}

	cndEmptyPattern := "empty pattern"
	cndContainsRegex := "regular expression"
	cndRewrite := "rewrite path"
	actCheckPath := "check returned path"
	actCheckMatched := "check error"
	actCheckNotMatched := "check no error"

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	tb.Condition(cndEmptyPattern, "set empty pattern")
	tb.Condition(cndContainsRegex, "pattern contains regular expression")
	tb.Condition(cndRewrite, "rewrite pattern")
	tb.Action(actCheckPath, "checked the returned path")
	tb.Action(actCheckMatched, "check that the pattern was matched ot the requested path")
	tb.Action(actCheckNotMatched, "check that the pattern was not matched to the requested path")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"empty string matched",
			[]string{cndEmptyPattern},
			[]string{actCheckPath, actCheckMatched},
			&condition{
				pattern: "",
				target:  "",
			},
			&action{
				path:    "",
				matched: true,
			},
		),
		gen(
			"empty pattern matches to all path",
			[]string{cndEmptyPattern},
			[]string{actCheckPath, actCheckMatched},
			&condition{
				pattern: "",
				target:  "/test",
			},
			&action{
				path:    "/test",
				matched: true,
			},
		),
		gen(
			"exactly the same path",
			[]string{},
			[]string{actCheckPath, actCheckMatched},
			&condition{
				pattern: "/test",
				target:  "/test",
			},
			&action{
				path:    "/test",
				matched: true,
			},
		),
		gen(
			"pattern not matched",
			[]string{},
			[]string{actCheckPath, actCheckMatched},
			&condition{
				pattern: "/test",
				target:  "/",
			},
			&action{
				path:    "/",
				matched: false,
			},
		),
		gen(
			"pattern matched to a regular expression",
			[]string{cndContainsRegex},
			[]string{actCheckPath, actCheckMatched},
			&condition{
				pattern: "^/test/(foo|bar)$",
				target:  "/test/foo",
			},
			&action{
				path:    "/test/foo",
				matched: true,
			},
		),
		gen(
			"pattern matched to a regular expression",
			[]string{cndContainsRegex},
			[]string{actCheckPath, actCheckMatched},
			&condition{
				pattern: "^/test/(foo|bar)$",
				target:  "/test/bar",
			},
			&action{
				path:    "/test/bar",
				matched: true,
			},
		),
		gen(
			"pattern matched to a regular expression",
			[]string{cndContainsRegex},
			[]string{actCheckPath, actCheckMatched},
			&condition{
				pattern: "^/..../(foo|bar)$",
				target:  "/test/bar",
			},
			&action{
				path:    "/test/bar",
				matched: true,
			},
		),
		gen(
			"pattern not matched to a regular expression",
			[]string{cndContainsRegex},
			[]string{actCheckPath, actCheckMatched},
			&condition{
				pattern: "^/..../(foo|bar)$",
				target:  "/test/baz",
			},
			&action{
				path:    "/test/baz",
				matched: false,
			},
		),
		gen(
			"pattern matched to a regular expression",
			[]string{cndContainsRegex},
			[]string{actCheckPath, actCheckMatched},
			&condition{
				pattern: "^/..../(foo|bar)$",
				target:  "/test/bar",
			},
			&action{
				path:    "/test/bar",
				matched: true,
			},
		),
		gen(
			"pattern matched to a regular expression",
			[]string{cndContainsRegex, cndRewrite},
			[]string{actCheckPath, actCheckMatched},
			&condition{
				pattern: "^/(?P<key>\\w+)/(?P<value>\\w+)$",
				rewrite: "/$value/$key",
				target:  "/foo/bar",
			},
			&action{
				path:    "/bar/foo",
				matched: true,
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			m := &matcher{
				pattern: tt.C().pattern,
				regexp:  regexp.MustCompile(tt.C().pattern),
				rewrite: tt.C().rewrite,
			}

			path, matched := m.regEx(tt.C().target)

			testutil.Diff(t, tt.A().matched, matched)
			testutil.Diff(t, tt.A().path, path)
		})
	}
}

func TestPathParamMatchers(t *testing.T) {
	type condition struct {
		specs   []*v1.ParamMatcherSpec
		pattern string // Pattern to register to a serve mux.
		url     string
	}

	type action struct {
		numMatcher int
		matchIndex int // -1 is not match.
		err        error
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"no spec",
			[]string{},
			[]string{},
			&condition{
				specs:   nil,
				pattern: "/",
				url:     "http://test.com/",
			},
			&action{
				numMatcher: 0,
				matchIndex: -1,
			},
		),
		gen(
			"empty key",
			[]string{},
			[]string{},
			&condition{
				specs: []*v1.ParamMatcherSpec{
					{Key: ""},
				},
				pattern: "/",
				url:     "http://test.com/",
			},
			&action{
				numMatcher: 0,
				matchIndex: -1,
			},
		),
		gen(
			"param not found",
			[]string{},
			[]string{},
			&condition{
				specs: []*v1.ParamMatcherSpec{
					{
						Key:       "foo",
						Patterns:  []string{"bar"},
						MatchType: k.MatchType_Exact,
					},
				},
				pattern: "/test/{alice}",
				url:     "http://test.com/test/bar",
			},
			&action{
				numMatcher: 1,
				matchIndex: -1,
			},
		),
		gen(
			"single matcher",
			[]string{},
			[]string{},
			&condition{
				specs: []*v1.ParamMatcherSpec{
					{
						Key:       "foo",
						Patterns:  []string{"bar"},
						MatchType: k.MatchType_Exact,
					},
				},
				pattern: "/test/{foo}",
				url:     "http://test.com/test/bar",
			},
			&action{
				numMatcher: 1,
				matchIndex: 0,
			},
		),
		gen(
			"multiple matcher",
			[]string{},
			[]string{},
			&condition{
				specs: []*v1.ParamMatcherSpec{
					{
						Key:       "alice",
						Patterns:  []string{"bob"},
						MatchType: k.MatchType_Exact,
					},
					{
						Key:       "foo",
						Patterns:  []string{"bar"},
						MatchType: k.MatchType_Exact,
					},
				},
				pattern: "/test/{foo}",
				url:     "http://test.com/test/bar",
			},
			&action{
				numMatcher: 2,
				matchIndex: 1,
			},
		),
		gen(
			"contains nil spec",
			[]string{},
			[]string{},
			&condition{
				specs: []*v1.ParamMatcherSpec{
					{
						Key:       "alice",
						Patterns:  []string{"bob"},
						MatchType: k.MatchType_Exact,
					},
					nil, nil, nil,
					{
						Key:       "foo",
						Patterns:  []string{"bar"},
						MatchType: k.MatchType_Exact,
					},
				},
				pattern: "/test/{foo}",
				url:     "http://test.com/test/bar",
			},
			&action{
				numMatcher: 2,
				matchIndex: 1,
			},
		),
		gen(
			"contains empty keys",
			[]string{},
			[]string{},
			&condition{
				specs: []*v1.ParamMatcherSpec{
					{
						Key:       "alice",
						Patterns:  []string{"bob"},
						MatchType: k.MatchType_Exact,
					},
					{Key: ""}, {Key: ""}, {Key: ""},
					{
						Key:       "foo",
						Patterns:  []string{"bar"},
						MatchType: k.MatchType_Exact,
					},
				},
				pattern: "/test/{foo}",
				url:     "http://test.com/test/bar",
			},
			&action{
				numMatcher: 2,
				matchIndex: 1,
			},
		),
		gen(
			"matcher create error",
			[]string{},
			[]string{},
			&condition{
				specs: []*v1.ParamMatcherSpec{
					{
						Key:       "foo",
						Patterns:  []string{"[0-9a-"},
						MatchType: k.MatchType_Regex,
					},
				},
				pattern: "/test/{foo}",
				url:     "http://test.com/test/bar",
			},
			&action{
				numMatcher: 0,
				matchIndex: -1,
				err: &er.Error{
					Package:     txtutil.ErrPkg,
					Type:        txtutil.ErrTypeMatcher,
					Description: txtutil.ErrDscPattern,
				},
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			matchers, err := pathParamMatchers(tt.C().specs...)
			testutil.Diff(t, tt.A().err, err, cmpopts.EquateErrors())
			testutil.Diff(t, tt.A().numMatcher, len(matchers))

			var called bool
			mux := &http.ServeMux{}
			mux.Handle(tt.C().pattern, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				called = true
				matchIndex := -1
				for i, matcher := range matchers {
					ok := matcher.Match(r)
					if ok {
						matchIndex = i
						break
					}
				}
				testutil.Diff(t, tt.A().matchIndex, matchIndex)
			}))

			r, err := http.NewRequest(http.MethodGet, tt.C().url, nil)
			testutil.Diff(t, nil, err)
			w := httptest.NewRecorder()
			mux.ServeHTTP(w, r)
			testutil.Diff(t, true, called) // Check the handler was called.
		})
	}
}

func TestQueryMatchers(t *testing.T) {
	type condition struct {
		specs []*v1.ParamMatcherSpec
		url   string
	}

	type action struct {
		numMatcher int
		matchIndex int // -1 is not match.
		err        error
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"no spec",
			[]string{},
			[]string{},
			&condition{
				specs: nil,
				url:   "http://test.com/",
			},
			&action{
				numMatcher: 0,
				matchIndex: -1,
			},
		),
		gen(
			"empty key",
			[]string{},
			[]string{},
			&condition{
				specs: []*v1.ParamMatcherSpec{
					{Key: ""},
				},
				url: "http://test.com/",
			},
			&action{
				numMatcher: 0,
				matchIndex: -1,
			},
		),
		gen(
			"query not found",
			[]string{},
			[]string{},
			&condition{
				specs: []*v1.ParamMatcherSpec{
					{
						Key:       "foo",
						Patterns:  []string{"bar"},
						MatchType: k.MatchType_Exact,
					},
				},
				url: "http://test.com/?alice=bob",
			},
			&action{
				numMatcher: 1,
				matchIndex: -1,
			},
		),
		gen(
			"single matcher",
			[]string{},
			[]string{},
			&condition{
				specs: []*v1.ParamMatcherSpec{
					{
						Key:       "foo",
						Patterns:  []string{"bar"},
						MatchType: k.MatchType_Exact,
					},
				},
				url: "http://test.com/?foo=bar",
			},
			&action{
				numMatcher: 1,
				matchIndex: 0,
			},
		),
		gen(
			"single matcher/multiple value found",
			[]string{},
			[]string{},
			&condition{
				specs: []*v1.ParamMatcherSpec{
					{
						Key:       "foo",
						Patterns:  []string{"bar,baz"},
						MatchType: k.MatchType_Exact,
					},
				},
				url: "http://test.com/?foo=bar&foo=baz",
			},
			&action{
				numMatcher: 1,
				matchIndex: 0,
			},
		),
		gen(
			"multiple matcher",
			[]string{},
			[]string{},
			&condition{
				specs: []*v1.ParamMatcherSpec{
					{
						Key:       "alice",
						Patterns:  []string{"bob"},
						MatchType: k.MatchType_Exact,
					},
					{
						Key:       "foo",
						Patterns:  []string{"bar"},
						MatchType: k.MatchType_Exact,
					},
				},
				url: "http://test.com/?foo=bar",
			},
			&action{
				numMatcher: 2,
				matchIndex: 1,
			},
		),
		gen(
			"contains nil spec",
			[]string{},
			[]string{},
			&condition{
				specs: []*v1.ParamMatcherSpec{
					{
						Key:       "alice",
						Patterns:  []string{"bob"},
						MatchType: k.MatchType_Exact,
					},
					nil, nil, nil,
					{
						Key:       "foo",
						Patterns:  []string{"bar"},
						MatchType: k.MatchType_Exact,
					},
				},
				url: "http://test.com/?foo=bar",
			},
			&action{
				numMatcher: 2,
				matchIndex: 1,
			},
		),
		gen(
			"contains empty keys",
			[]string{},
			[]string{},
			&condition{
				specs: []*v1.ParamMatcherSpec{
					{
						Key:       "alice",
						Patterns:  []string{"bob"},
						MatchType: k.MatchType_Exact,
					},
					{Key: ""}, {Key: ""}, {Key: ""},
					{
						Key:       "foo",
						Patterns:  []string{"bar"},
						MatchType: k.MatchType_Exact,
					},
				},
				url: "http://test.com/?foo=bar",
			},
			&action{
				numMatcher: 2,
				matchIndex: 1,
			},
		),
		gen(
			"matcher create error",
			[]string{},
			[]string{},
			&condition{
				specs: []*v1.ParamMatcherSpec{
					{
						Key:       "foo",
						Patterns:  []string{"[0-9a-"},
						MatchType: k.MatchType_Regex,
					},
				},
				url: "http://test.com/test/bar",
			},
			&action{
				numMatcher: 0,
				matchIndex: -1,
				err: &er.Error{
					Package:     txtutil.ErrPkg,
					Type:        txtutil.ErrTypeMatcher,
					Description: txtutil.ErrDscPattern,
				},
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			matchers, err := queryMatchers(tt.C().specs...)
			testutil.Diff(t, tt.A().err, err, cmpopts.EquateErrors())
			testutil.Diff(t, tt.A().numMatcher, len(matchers))

			r, err := http.NewRequest(http.MethodGet, tt.C().url, nil)
			testutil.Diff(t, nil, err)

			matchIndex := -1
			for i, matcher := range matchers {
				ok := matcher.Match(r)
				if ok {
					matchIndex = i
					break
				}
			}
			testutil.Diff(t, tt.A().matchIndex, matchIndex)
		})
	}
}

func TestHeaderMatchers(t *testing.T) {
	type condition struct {
		specs  []*v1.ParamMatcherSpec
		header http.Header
	}

	type action struct {
		numMatcher int
		matchIndex int // -1 is not match.
		err        error
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"no spec",
			[]string{},
			[]string{},
			&condition{
				specs:  nil,
				header: http.Header{"Foo": {"bar"}},
			},
			&action{
				numMatcher: 0,
				matchIndex: -1,
			},
		),
		gen(
			"empty key",
			[]string{},
			[]string{},
			&condition{
				specs: []*v1.ParamMatcherSpec{
					{Key: ""},
				},
				header: http.Header{"Foo": {"bar"}},
			},
			&action{
				numMatcher: 0,
				matchIndex: -1,
			},
		),
		gen(
			"header not found",
			[]string{},
			[]string{},
			&condition{
				specs: []*v1.ParamMatcherSpec{
					{
						Key:       "foo",
						Patterns:  []string{"bar"},
						MatchType: k.MatchType_Exact,
					},
				},
				header: http.Header{"Alice": {"bob"}},
			},
			&action{
				numMatcher: 1,
				matchIndex: -1,
			},
		),
		gen(
			"single matcher",
			[]string{},
			[]string{},
			&condition{
				specs: []*v1.ParamMatcherSpec{
					{
						Key:       "foo",
						Patterns:  []string{"bar"},
						MatchType: k.MatchType_Exact,
					},
				},
				header: http.Header{"Foo": {"bar"}},
			},
			&action{
				numMatcher: 1,
				matchIndex: 0,
			},
		),
		gen(
			"single matcher/multiple value found",
			[]string{},
			[]string{},
			&condition{
				specs: []*v1.ParamMatcherSpec{
					{
						Key:       "foo",
						Patterns:  []string{"bar,baz"},
						MatchType: k.MatchType_Exact,
					},
				},
				header: http.Header{"Foo": {"bar", "baz"}},
			},
			&action{
				numMatcher: 1,
				matchIndex: 0,
			},
		),
		gen(
			"multiple matcher",
			[]string{},
			[]string{},
			&condition{
				specs: []*v1.ParamMatcherSpec{
					{
						Key:       "alice",
						Patterns:  []string{"bob"},
						MatchType: k.MatchType_Exact,
					},
					{
						Key:       "foo",
						Patterns:  []string{"bar"},
						MatchType: k.MatchType_Exact,
					},
				},
				header: http.Header{"Foo": {"bar"}},
			},
			&action{
				numMatcher: 2,
				matchIndex: 1,
			},
		),
		gen(
			"contains nil spec",
			[]string{},
			[]string{},
			&condition{
				specs: []*v1.ParamMatcherSpec{
					{
						Key:       "alice",
						Patterns:  []string{"bob"},
						MatchType: k.MatchType_Exact,
					},
					nil, nil, nil,
					{
						Key:       "foo",
						Patterns:  []string{"bar"},
						MatchType: k.MatchType_Exact,
					},
				},
				header: http.Header{"Foo": {"bar"}},
			},
			&action{
				numMatcher: 2,
				matchIndex: 1,
			},
		),
		gen(
			"contains empty keys",
			[]string{},
			[]string{},
			&condition{
				specs: []*v1.ParamMatcherSpec{
					{
						Key:       "alice",
						Patterns:  []string{"bob"},
						MatchType: k.MatchType_Exact,
					},
					{Key: ""}, {Key: ""}, {Key: ""},
					{
						Key:       "foo",
						Patterns:  []string{"bar"},
						MatchType: k.MatchType_Exact,
					},
				},
				header: http.Header{"Foo": {"bar"}},
			},
			&action{
				numMatcher: 2,
				matchIndex: 1,
			},
		),
		gen(
			"matcher create error",
			[]string{},
			[]string{},
			&condition{
				specs: []*v1.ParamMatcherSpec{
					{
						Key:       "foo",
						Patterns:  []string{"[0-9a-"},
						MatchType: k.MatchType_Regex,
					},
				},
				header: http.Header{"Foo": {"bar"}},
			},
			&action{
				numMatcher: 0,
				matchIndex: -1,
				err: &er.Error{
					Package:     txtutil.ErrPkg,
					Type:        txtutil.ErrTypeMatcher,
					Description: txtutil.ErrDscPattern,
				},
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			matchers, err := headerMatchers(tt.C().specs...)
			testutil.Diff(t, tt.A().err, err, cmpopts.EquateErrors())
			testutil.Diff(t, tt.A().numMatcher, len(matchers))

			r, err := http.NewRequest(http.MethodGet, "http://test.com/", nil)
			testutil.Diff(t, nil, err)
			r.Header = tt.C().header

			matchIndex := -1
			for i, matcher := range matchers {
				ok := matcher.Match(r)
				if ok {
					matchIndex = i
					break
				}
			}
			testutil.Diff(t, tt.A().matchIndex, matchIndex)
		})
	}
}
