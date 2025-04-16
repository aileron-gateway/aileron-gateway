// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package httpproxy

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"sync/atomic"
	"testing"

	v1 "github.com/aileron-gateway/aileron-gateway/apis/core/v1"
	"github.com/aileron-gateway/aileron-gateway/apis/kernel"
	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
	"github.com/aileron-gateway/aileron-gateway/kernel/txtutil"
	"github.com/aileron-gateway/aileron-gateway/util/resilience"
	"github.com/google/go-cmp/cmp"
)

func TestLBMatcher(t *testing.T) {
	type condition struct {
		matcher *lbMatcher
		pattern string
		url     string
		method  string
		header  http.Header
	}

	type action struct {
		path    string
		matched bool
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	mustMatcher := func(typ txtutil.MatchType, patterns ...string) txtutil.MatchFunc[string] {
		mf, err := txtutil.NewStringMatcher(typ, patterns...)
		if err != nil {
			panic(err)
		}
		return mf.Match
	}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"minimum matchers/match",
			[]string{},
			[]string{},
			&condition{
				pattern: "/test/{pp}",
				matcher: &lbMatcher{
					pathMatchers: []matcherFunc{(&matcher{pattern: "/test"}).prefix},
				},
				url:    "http://test.com/test/path-param?qp=query-param",
				method: http.MethodGet,
				header: http.Header{"hp": []string{"header-param"}},
			},
			&action{
				path:    "/test/path-param",
				matched: true,
			},
		),
		gen(
			"all matchers/match",
			[]string{},
			[]string{},
			&condition{
				pattern: "/test/{pp}",
				matcher: &lbMatcher{
					hosts:        []string{"test.com"},
					methods:      []string{http.MethodGet},
					pathMatchers: []matcherFunc{(&matcher{pattern: "/test"}).prefix},
					paramMatchers: []txtutil.Matcher[*http.Request]{
						&headerMatcher{key: "hp", f: mustMatcher(txtutil.MatchTypeExact, "header-param")},
						&queryMatcher{key: "qp", f: mustMatcher(txtutil.MatchTypeExact, "query-param")},
						&pathParamMatcher{key: "pp", f: mustMatcher(txtutil.MatchTypeExact, "path-param")},
					},
				},
				url:    "http://test.com/test/path-param?qp=query-param",
				method: http.MethodGet,
				header: http.Header{"hp": []string{"header-param"}},
			},
			&action{
				path:    "/test/path-param",
				matched: true,
			},
		),
		gen(
			"all matchers/host mismatch",
			[]string{},
			[]string{},
			&condition{
				pattern: "/test/{pp}",
				matcher: &lbMatcher{
					hosts:        []string{"test.com"},
					methods:      []string{http.MethodGet},
					pathMatchers: []matcherFunc{(&matcher{pattern: "/test"}).prefix},
					paramMatchers: []txtutil.Matcher[*http.Request]{
						&headerMatcher{key: "hp", f: mustMatcher(txtutil.MatchTypeExact, "header-param")},
						&queryMatcher{key: "qp", f: mustMatcher(txtutil.MatchTypeExact, "query-param")},
						&pathParamMatcher{key: "pp", f: mustMatcher(txtutil.MatchTypeExact, "path-param")},
					},
				},
				url:    "http://mismatch.com/test/path-param?qp=query-param",
				method: http.MethodGet,
				header: http.Header{"hp": []string{"header-param"}},
			},
			&action{
				path:    "",
				matched: false,
			},
		),
		gen(
			"all matchers/method mismatch",
			[]string{},
			[]string{},
			&condition{
				pattern: "/test/{pp}",
				matcher: &lbMatcher{
					hosts:        []string{"test.com"},
					methods:      []string{http.MethodGet},
					pathMatchers: []matcherFunc{(&matcher{pattern: "/test"}).prefix},
					paramMatchers: []txtutil.Matcher[*http.Request]{
						&headerMatcher{key: "hp", f: mustMatcher(txtutil.MatchTypeExact, "header-param")},
						&queryMatcher{key: "qp", f: mustMatcher(txtutil.MatchTypeExact, "query-param")},
						&pathParamMatcher{key: "pp", f: mustMatcher(txtutil.MatchTypeExact, "path-param")},
					},
				},
				url:    "http://test.com/test/path-param?qp=query-param",
				method: http.MethodPost,
				header: http.Header{"hp": []string{"header-param"}},
			},
			&action{
				path:    "",
				matched: false,
			},
		),
		gen(
			"all matchers/path param mismatch",
			[]string{},
			[]string{},
			&condition{
				pattern: "/test/{pp}",
				matcher: &lbMatcher{
					hosts:        []string{"test.com"},
					methods:      []string{http.MethodGet},
					pathMatchers: []matcherFunc{(&matcher{pattern: "/test"}).prefix},
					paramMatchers: []txtutil.Matcher[*http.Request]{
						&headerMatcher{key: "hp", f: mustMatcher(txtutil.MatchTypeExact, "header-param")},
						&queryMatcher{key: "qp", f: mustMatcher(txtutil.MatchTypeExact, "query-param")},
						&pathParamMatcher{key: "pp", f: mustMatcher(txtutil.MatchTypeExact, "path-param")},
					},
				},
				url:    "http://test.com/test/mismatch-param?qp=query-param",
				method: http.MethodGet,
				header: http.Header{"hp": []string{"header-param"}},
			},
			&action{
				path:    "",
				matched: false,
			},
		),
		gen(
			"all matchers/query mismatch",
			[]string{},
			[]string{},
			&condition{
				pattern: "/test/{pp}",
				matcher: &lbMatcher{
					hosts:        []string{"test.com"},
					methods:      []string{http.MethodGet},
					pathMatchers: []matcherFunc{(&matcher{pattern: "/test"}).prefix},
					paramMatchers: []txtutil.Matcher[*http.Request]{
						&headerMatcher{key: "hp", f: mustMatcher(txtutil.MatchTypeExact, "header-param")},
						&queryMatcher{key: "qp", f: mustMatcher(txtutil.MatchTypeExact, "query-param")},
						&pathParamMatcher{key: "pp", f: mustMatcher(txtutil.MatchTypeExact, "path-param")},
					},
				},
				url:    "http://test.com/test/path-param?qp=mismatch-param",
				method: http.MethodGet,
				header: http.Header{"hp": []string{"header-param"}},
			},
			&action{
				path:    "",
				matched: false,
			},
		),
		gen(
			"all matchers/header mismatch",
			[]string{},
			[]string{},
			&condition{
				pattern: "/test/{pp}",
				matcher: &lbMatcher{
					hosts:        []string{"test.com"},
					methods:      []string{http.MethodGet},
					pathMatchers: []matcherFunc{(&matcher{pattern: "/test"}).prefix},
					paramMatchers: []txtutil.Matcher[*http.Request]{
						&headerMatcher{key: "hp", f: mustMatcher(txtutil.MatchTypeExact, "header-param")},
						&queryMatcher{key: "qp", f: mustMatcher(txtutil.MatchTypeExact, "query-param")},
						&pathParamMatcher{key: "pp", f: mustMatcher(txtutil.MatchTypeExact, "path-param")},
					},
				},
				url:    "http://test.com/test/path-param?qp=query-param",
				method: http.MethodGet,
				header: http.Header{"hp": []string{"mismatch-param"}},
			},
			&action{
				path:    "",
				matched: false,
			},
		),
		gen(
			"no hosts/match",
			[]string{},
			[]string{},
			&condition{
				pattern: "/test/{pp}",
				matcher: &lbMatcher{
					hosts:        []string{},
					methods:      []string{http.MethodGet},
					pathMatchers: []matcherFunc{(&matcher{pattern: "/test"}).prefix},
					paramMatchers: []txtutil.Matcher[*http.Request]{
						&headerMatcher{key: "hp", f: mustMatcher(txtutil.MatchTypeExact, "header-param")},
						&queryMatcher{key: "qp", f: mustMatcher(txtutil.MatchTypeExact, "query-param")},
						&pathParamMatcher{key: "pp", f: mustMatcher(txtutil.MatchTypeExact, "path-param")},
					},
				},
				url:    "http://dummy.com/test/path-param?qp=query-param",
				method: http.MethodGet,
				header: http.Header{"hp": []string{"header-param"}},
			},
			&action{
				path:    "/test/path-param",
				matched: true,
			},
		),
		gen(
			"no methods/match",
			[]string{},
			[]string{},
			&condition{
				pattern: "/test/{pp}",
				matcher: &lbMatcher{
					hosts:        []string{"test.com"},
					methods:      []string{},
					pathMatchers: []matcherFunc{(&matcher{pattern: "/test"}).prefix},
					paramMatchers: []txtutil.Matcher[*http.Request]{
						&headerMatcher{key: "hp", f: mustMatcher(txtutil.MatchTypeExact, "header-param")},
						&queryMatcher{key: "qp", f: mustMatcher(txtutil.MatchTypeExact, "query-param")},
						&pathParamMatcher{key: "pp", f: mustMatcher(txtutil.MatchTypeExact, "path-param")},
					},
				},
				url:    "http://test.com/test/path-param?qp=query-param",
				method: http.MethodHead,
				header: http.Header{"hp": []string{"header-param"}},
			},
			&action{
				path:    "/test/path-param",
				matched: true,
			},
		),
		gen(
			"no path param matcher/match",
			[]string{},
			[]string{},
			&condition{
				pattern: "/test/{pp}",
				matcher: &lbMatcher{
					hosts:        []string{"test.com"},
					methods:      []string{http.MethodGet},
					pathMatchers: []matcherFunc{(&matcher{pattern: "/test"}).prefix},
					paramMatchers: []txtutil.Matcher[*http.Request]{
						&headerMatcher{key: "hp", f: mustMatcher(txtutil.MatchTypeExact, "header-param")},
						&queryMatcher{key: "qp", f: mustMatcher(txtutil.MatchTypeExact, "query-param")},
					},
				},
				url:    "http://test.com/test/dummy?qp=query-param",
				method: http.MethodGet,
				header: http.Header{"hp": []string{"header-param"}},
			},
			&action{
				path:    "/test/dummy",
				matched: true,
			},
		),
		gen(
			"no header matcher/match",
			[]string{},
			[]string{},
			&condition{
				pattern: "/test/{pp}",
				matcher: &lbMatcher{
					hosts:        []string{"test.com"},
					methods:      []string{http.MethodGet},
					pathMatchers: []matcherFunc{(&matcher{pattern: "/test"}).prefix},
					paramMatchers: []txtutil.Matcher[*http.Request]{
						&queryMatcher{key: "qp", f: mustMatcher(txtutil.MatchTypeExact, "query-param")},
						&pathParamMatcher{key: "pp", f: mustMatcher(txtutil.MatchTypeExact, "path-param")},
					},
				},
				url:    "http://test.com/test/path-param?qp=query-param",
				method: http.MethodGet,
				header: http.Header{"hp": []string{"dummy"}},
			},
			&action{
				path:    "/test/path-param",
				matched: true,
			},
		),
		gen(
			"no query matcher/match",
			[]string{},
			[]string{},
			&condition{
				pattern: "/test/{pp}",
				matcher: &lbMatcher{
					hosts:        []string{"test.com"},
					methods:      []string{http.MethodGet},
					pathMatchers: []matcherFunc{(&matcher{pattern: "/test"}).prefix},
					paramMatchers: []txtutil.Matcher[*http.Request]{
						&headerMatcher{key: "hp", f: mustMatcher(txtutil.MatchTypeExact, "header-param")},
						&pathParamMatcher{key: "pp", f: mustMatcher(txtutil.MatchTypeExact, "path-param")},
					},
				},
				url:    "http://test.com/test/path-param?qp=dummy",
				method: http.MethodGet,
				header: http.Header{"hp": []string{"header-param"}},
			},
			&action{
				path:    "/test/path-param",
				matched: true,
			},
		),
		gen(
			"multiple header",
			[]string{},
			[]string{},
			&condition{
				pattern: "/",
				matcher: &lbMatcher{
					pathMatchers: []matcherFunc{(&matcher{pattern: "/test"}).prefix},
					paramMatchers: []txtutil.Matcher[*http.Request]{
						&headerMatcher{key: "hp", f: mustMatcher(txtutil.MatchTypeExact, "value1,value2")},
					},
				},
				url:    "http://test.com/test",
				method: http.MethodGet,
				header: http.Header{"hp": []string{"value1", "value2"}},
			},
			&action{
				path:    "/test",
				matched: true,
			},
		),
		gen(
			"multiple query",
			[]string{},
			[]string{},
			&condition{
				pattern: "/",
				matcher: &lbMatcher{
					pathMatchers: []matcherFunc{(&matcher{pattern: "/"}).prefix},
					paramMatchers: []txtutil.Matcher[*http.Request]{
						&queryMatcher{key: "qp", f: mustMatcher(txtutil.MatchTypeExact, "value1,value2")},
					},
				},
				url:    "http://test.com/test?qp=value1&qp=value2",
				method: http.MethodGet,
				header: http.Header{"hp": []string{"value1", "value2"}},
			},
			&action{
				path:    "/test",
				matched: true,
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			mux := &http.ServeMux{}
			mux.HandleFunc(tt.C().pattern, func(w http.ResponseWriter, r *http.Request) {
				path, matched := tt.C().matcher.match(r)
				testutil.Diff(t, tt.A().matched, matched)
				testutil.Diff(t, tt.A().path, path)
			})

			r, _ := http.NewRequest(tt.C().method, tt.C().url, nil)
			r.Header = tt.C().header
			w := httptest.NewRecorder()
			mux.ServeHTTP(w, r)
		})
	}
}

func TestNonHashLB(t *testing.T) {
	type condition struct {
		matcher   func(string) (string, bool)
		upstreams []upstream
	}

	type action struct {
		upstream []upstream
		url      []*url.URL
		matched  []bool
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	cndPathNotMatch := tb.Condition("path not match", "path matches")
	cndNoUpstream := tb.Condition("no upstream", "load balancer has no upstream")
	cndSingleUpstream := tb.Condition("single upstream", "set single upstream")
	cndMultiUpstream := tb.Condition("multi upstreams", "set multiple upstreams")
	cndInactiveUpstream := tb.Condition("inactive upstreams", "upstream is inactive")
	actNoUpstream := tb.Action("no upstream", "check that no upstream returned")
	table := tb.Build()

	ups1 := &noopUpstream{rawURL: "http://test1.com", weight: 1, parsedURL: &url.URL{Scheme: "http", Host: "test.com", RawQuery: "bar1=baz1"}}
	ups2 := &noopUpstream{rawURL: "http://test2.com", weight: 1, parsedURL: &url.URL{Scheme: "http", Host: "test.com", RawQuery: "bar2=baz2"}}
	url1 := &url.URL{Scheme: "http", Host: "test.com", Path: "/foo", RawPath: "/foo", RawQuery: "bar1=baz1"}
	url2 := &url.URL{Scheme: "http", Host: "test.com", Path: "/foo", RawPath: "/foo", RawQuery: "bar2=baz2"}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"no upstream",
			[]string{cndNoUpstream},
			[]string{actNoUpstream},
			&condition{
				upstreams: []upstream{},
				matcher:   func(s string) (string, bool) { return s, true },
			},
			&action{
				upstream: []upstream{nil, nil, nil}, // Check 3 times
				url:      []*url.URL{nil, nil, nil},
				matched:  []bool{true, true, true},
			},
		),
		gen(
			"single upstream",
			[]string{cndSingleUpstream},
			[]string{},
			&condition{
				upstreams: []upstream{ups1},
				matcher:   func(s string) (string, bool) { return s, true },
			},
			&action{
				upstream: []upstream{ups1, ups1, ups1}, // Check 3 times
				url:      []*url.URL{url1, url1, url1},
				matched:  []bool{true, true, true},
			},
		),
		gen(
			"multiple upstream",
			[]string{cndMultiUpstream},
			[]string{},
			&condition{
				upstreams: []upstream{ups1, ups2},
				matcher:   func(s string) (string, bool) { return s, true },
			},
			&action{
				upstream: []upstream{ups1, ups2, ups1, ups2}, // Check 4 times
				url:      []*url.URL{url1, url2, url1, url2},
				matched:  []bool{true, true, true, true},
			},
		),
		gen(
			"single upstream/inactive",
			[]string{cndInactiveUpstream},
			[]string{actNoUpstream},
			&condition{
				upstreams: []upstream{
					&lbUpstream{
						circuitBreaker: &circuitBreakerController{status: opened},
						rawURL:         "http://inactive.com",
						parsedURL:      &url.URL{Scheme: "http", Host: "inactive.com"},
					},
				},
				matcher: func(s string) (string, bool) { return s, true },
			},
			&action{
				upstream: []upstream{nil, nil, nil}, // Check 3 times
				url:      []*url.URL{nil, nil, nil},
				matched:  []bool{true, true, true},
			},
		),
		gen(
			"multiple upstream/inactive",
			[]string{cndInactiveUpstream},
			[]string{actNoUpstream},
			&condition{
				upstreams: []upstream{
					&lbUpstream{
						circuitBreaker: &circuitBreakerController{status: opened},
						rawURL:         "http://inactive1.com",
						parsedURL:      &url.URL{Scheme: "http", Host: "inactive1.com"},
					},
					&lbUpstream{
						circuitBreaker: &circuitBreakerController{status: opened},
						rawURL:         "http://inactive2.com",
						parsedURL:      &url.URL{Scheme: "http", Host: "inactive2.com"},
					},
				},
				matcher: func(s string) (string, bool) { return s, true },
			},
			&action{
				upstream: []upstream{nil, nil, nil}, // Check 3 times
				url:      []*url.URL{nil, nil, nil},
				matched:  []bool{true, true, true},
			},
		),
		gen(
			"path not match",
			[]string{cndPathNotMatch},
			[]string{actNoUpstream},
			&condition{
				upstreams: []upstream{ups1, ups2},
				matcher:   func(string) (string, bool) { return "", false },
			},
			&action{
				upstream: []upstream{nil, nil, nil}, // Check 3 times
				url:      []*url.URL{nil, nil, nil},
				matched:  []bool{false, false, false},
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			rlb := &resilience.RoundRobinLB[upstream]{}
			rlb.Add(tt.C().upstreams...)
			lb := &nonHashLB{
				lbMatcher: &lbMatcher{
					pathMatchers: []matcherFunc{tt.C().matcher},
				},
				LoadBalancer: rlb,
			}

			r := httptest.NewRequest(http.MethodGet, "http://test.com/foo", nil)

			for i := 0; i < len(tt.A().upstream); i++ {
				upstream, url, matched := lb.upstream(r)
				testutil.Diff(t, tt.A().upstream[i], upstream, cmp.AllowUnexported(lbUpstream{}, noopUpstream{}, atomic.Int32{}))
				testutil.Diff(t, tt.A().url[i], url)
				testutil.Diff(t, tt.A().matched[i], matched)
			}
		})
	}
}

func TestDirectHashLB(t *testing.T) {
	type condition struct {
		matcher   func(string) (string, bool)
		upstreams []upstream
		hashers   []resilience.HTTPHasher
	}

	type action struct {
		upstream []upstream
		url      []*url.URL
		matched  []bool
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	cndPathNotMatch := tb.Condition("path not match", "path matches")
	cndNoUpstream := tb.Condition("no upstream", "load balancer has no upstream")
	cndSingleUpstream := tb.Condition("single upstream", "set single upstream")
	cndMultiUpstream := tb.Condition("multi upstreams", "set multiple upstreams")
	cndInactiveUpstream := tb.Condition("inactive upstreams", "upstream is inactive")
	actNoUpstream := tb.Action("no upstream", "check that no upstream returned")
	table := tb.Build()

	ups1 := &noopUpstream{rawURL: "http://test1.com", weight: 1, parsedURL: &url.URL{Scheme: "http", Host: "test.com", RawQuery: "bar1=baz1"}}
	ups2 := &noopUpstream{rawURL: "http://test2.com", weight: 1, parsedURL: &url.URL{Scheme: "http", Host: "test.com", RawQuery: "bar2=baz2"}}
	ups3 := &noopUpstream{rawURL: "http://test3.com", weight: 1, parsedURL: &url.URL{Scheme: "http", Host: "test.com", RawQuery: "bar3=baz3"}}
	url1 := &url.URL{Scheme: "http", Host: "test.com", Path: "/foo", RawPath: "/foo", RawQuery: "bar1=baz1"}
	// url2 := &url.URL{Scheme: "http", Host: "test.com", Path: "/foo", RawPath: "/foo", RawQuery: "bar2=baz2"}
	url3 := &url.URL{Scheme: "http", Host: "test.com", Path: "/foo", RawPath: "/foo", RawQuery: "bar3=baz3"}
	hasher1 := resilience.NewHTTPHasher(&v1.HTTPHasherSpec{HasherType: v1.HTTPHasherType_Header, Key: "test", HashAlg: kernel.HashAlg_FNV1a_32})
	hasher2 := resilience.NewHTTPHasher(&v1.HTTPHasherSpec{HasherType: v1.HTTPHasherType_Header, Key: "foo", HashAlg: kernel.HashAlg_FNV1a_32})
	hasher3 := resilience.NewHTTPHasher(&v1.HTTPHasherSpec{HasherType: v1.HTTPHasherType_Header, Key: "test", HashAlg: kernel.HashAlg_FNV1a_128})

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"no upstream/no hasher",
			[]string{cndNoUpstream},
			[]string{actNoUpstream},
			&condition{
				upstreams: []upstream{},
				matcher:   func(s string) (string, bool) { return s, true },
			},
			&action{
				upstream: []upstream{nil},
				url:      []*url.URL{nil},
				matched:  []bool{true},
			},
		),
		gen(
			"single upstream/no hasher",
			[]string{cndSingleUpstream},
			[]string{},
			&condition{
				upstreams: []upstream{ups1},
				matcher:   func(s string) (string, bool) { return s, true },
			},
			&action{
				upstream: []upstream{nil},
				url:      []*url.URL{nil},
				matched:  []bool{true},
			},
		),
		gen(
			"multiple upstream/no hasher",
			[]string{cndMultiUpstream},
			[]string{},
			&condition{
				upstreams: []upstream{ups1, ups2, ups3},
				matcher:   func(s string) (string, bool) { return s, true },
			},
			&action{
				upstream: []upstream{nil},
				url:      []*url.URL{nil},
				matched:  []bool{true},
			},
		),
		gen(
			"no upstream/single hasher",
			[]string{cndNoUpstream},
			[]string{actNoUpstream},
			&condition{
				upstreams: []upstream{},
				matcher:   func(s string) (string, bool) { return s, true },
				hashers:   []resilience.HTTPHasher{hasher1},
			},
			&action{
				upstream: []upstream{nil},
				url:      []*url.URL{nil},
				matched:  []bool{true},
			},
		),
		gen(
			"single upstream/single hasher",
			[]string{cndSingleUpstream},
			[]string{},
			&condition{
				upstreams: []upstream{ups1},
				matcher:   func(s string) (string, bool) { return s, true },
				hashers:   []resilience.HTTPHasher{hasher1},
			},
			&action{
				upstream: []upstream{ups1},
				url:      []*url.URL{url1},
				matched:  []bool{true},
			},
		),
		gen(
			"multiple upstream/single hasher",
			[]string{cndMultiUpstream},
			[]string{},
			&condition{
				upstreams: []upstream{ups1, ups2, ups3},
				matcher:   func(s string) (string, bool) { return s, true },
				hashers:   []resilience.HTTPHasher{hasher1},
			},
			&action{
				upstream: []upstream{ups3},
				url:      []*url.URL{url3},
				matched:  []bool{true},
			},
		),
		gen(
			"no upstream/multiple hasher",
			[]string{cndNoUpstream},
			[]string{actNoUpstream},
			&condition{
				upstreams: []upstream{},
				matcher:   func(s string) (string, bool) { return s, true },
				hashers:   []resilience.HTTPHasher{hasher2, hasher3, hasher1},
			},
			&action{
				upstream: []upstream{nil},
				url:      []*url.URL{nil},
				matched:  []bool{true},
			},
		),
		gen(
			"single upstream/multiple hasher",
			[]string{cndSingleUpstream},
			[]string{},
			&condition{
				upstreams: []upstream{ups1},
				matcher:   func(s string) (string, bool) { return s, true },
				hashers:   []resilience.HTTPHasher{hasher2, hasher3, hasher1},
			},
			&action{
				upstream: []upstream{ups1},
				url:      []*url.URL{url1},
				matched:  []bool{true},
			},
		),
		gen(
			"multiple upstream/multiple hasher",
			[]string{cndMultiUpstream},
			[]string{},
			&condition{
				upstreams: []upstream{ups1, ups2, ups3},
				matcher:   func(s string) (string, bool) { return s, true },
				hashers:   []resilience.HTTPHasher{hasher2, hasher3, hasher1},
			},
			&action{
				upstream: []upstream{ups1},
				url:      []*url.URL{url1},
				matched:  []bool{true},
			},
		),
		gen(
			"single upstream/inactive upstream",
			[]string{cndSingleUpstream, cndInactiveUpstream},
			[]string{},
			&condition{
				upstreams: []upstream{
					&lbUpstream{
						weight:         1,
						circuitBreaker: &circuitBreakerController{status: opened},
						rawURL:         "http://inactive.com",
						parsedURL:      &url.URL{Scheme: "http", Host: "inactive"},
					},
				},
				matcher: func(s string) (string, bool) { return s, true },
				hashers: []resilience.HTTPHasher{hasher1},
			},
			&action{
				upstream: []upstream{nil},
				url:      []*url.URL{nil},
				matched:  []bool{true},
			},
		),
		gen(
			"multiple upstream/1 inactive upstream",
			[]string{cndMultiUpstream, cndInactiveUpstream},
			[]string{},
			&condition{
				upstreams: []upstream{
					&lbUpstream{
						weight:         1,
						circuitBreaker: &circuitBreakerController{status: opened},
						rawURL:         "http://inactive.com",
						parsedURL:      &url.URL{Scheme: "http", Host: "inactive"},
					},
					ups2, ups3},
				matcher: func(s string) (string, bool) { return s, true },
				hashers: []resilience.HTTPHasher{hasher2, hasher3, hasher1},
			},
			&action{
				upstream: []upstream{ups3},
				url:      []*url.URL{url3},
				matched:  []bool{true},
			},
		),
		gen(
			"multiple upstream/2 inactive upstream",
			[]string{cndMultiUpstream, cndInactiveUpstream},
			[]string{},
			&condition{
				upstreams: []upstream{
					&lbUpstream{
						weight:         1,
						circuitBreaker: &circuitBreakerController{status: opened},
						rawURL:         "http://inactive1.com",
						parsedURL:      &url.URL{Scheme: "http", Host: "inactive1"},
					},
					&lbUpstream{
						weight:         1,
						circuitBreaker: &circuitBreakerController{status: opened},
						rawURL:         "http://inactive2.com",
						parsedURL:      &url.URL{Scheme: "http", Host: "inactive2"},
					},
					ups3},
				matcher: func(s string) (string, bool) { return s, true },
				hashers: []resilience.HTTPHasher{hasher2, hasher3, hasher1},
			},
			&action{
				upstream: []upstream{ups3},
				url:      []*url.URL{url3},
				matched:  []bool{true},
			},
		),
		gen(
			"path not match",
			[]string{cndPathNotMatch},
			[]string{actNoUpstream},
			&condition{
				upstreams: []upstream{ups1},
				matcher:   func(string) (string, bool) { return "", false },
			},
			&action{
				upstream: []upstream{nil},
				url:      []*url.URL{nil},
				matched:  []bool{false},
			},
		),
		gen(
			"hasher failed",
			[]string{cndMultiUpstream},
			[]string{},
			&condition{
				upstreams: []upstream{ups1, ups2, ups3},
				matcher:   func(s string) (string, bool) { return s, true },
				hashers:   []resilience.HTTPHasher{hasher2},
			},
			&action{
				upstream: []upstream{nil},
				url:      []*url.URL{nil},
				matched:  []bool{true},
			},
		),
		gen(
			"all hashers failed",
			[]string{cndMultiUpstream},
			[]string{},
			&condition{
				upstreams: []upstream{ups1, ups2, ups3},
				matcher:   func(s string) (string, bool) { return s, true },
				hashers:   []resilience.HTTPHasher{hasher2, hasher2, hasher2},
			},
			&action{
				upstream: []upstream{nil},
				url:      []*url.URL{nil},
				matched:  []bool{true},
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			rlb := &resilience.DirectHashLB[upstream]{}
			rlb.Add(tt.C().upstreams...)
			lb := &directHashLB{
				lbMatcher: &lbMatcher{
					pathMatchers: []matcherFunc{tt.C().matcher},
				},
				LoadBalancer: rlb,
				hashers:      tt.C().hashers,
			}

			r := httptest.NewRequest(http.MethodGet, "http://test.com/foo", nil)
			r.Header.Set("test", "hash input")

			for i := 0; i < len(tt.A().upstream); i++ {
				upstream, url, matched := lb.upstream(r)
				testutil.Diff(t, tt.A().upstream[i], upstream, cmp.AllowUnexported(lbUpstream{}, noopUpstream{}, atomic.Int32{}))
				testutil.Diff(t, tt.A().url[i], url)
				testutil.Diff(t, tt.A().matched[i], matched)
			}
		})
	}
}

func TestMaglevLB(t *testing.T) {
	type condition struct {
		matcher   func(string) (string, bool)
		upstreams []upstream
		hashers   []resilience.HTTPHasher
	}

	type action struct {
		upstream []upstream
		url      []*url.URL
		matched  []bool
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	cndPathNotMatch := tb.Condition("path not match", "path matches")
	cndNoUpstream := tb.Condition("no upstream", "load balancer has no upstream")
	cndSingleUpstream := tb.Condition("single upstream", "set single upstream")
	cndMultiUpstream := tb.Condition("multi upstreams", "set multiple upstreams")
	cndInactiveUpstream := tb.Condition("inactive upstreams", "upstream is inactive")
	actNoUpstream := tb.Action("no upstream", "check that no upstream returned")
	table := tb.Build()

	ups1 := &noopUpstream{rawURL: "http://test1.com", weight: 1, parsedURL: &url.URL{Scheme: "http", Host: "test.com", RawQuery: "bar1=baz1"}}
	ups2 := &noopUpstream{rawURL: "http://test2.com", weight: 1, parsedURL: &url.URL{Scheme: "http", Host: "test.com", RawQuery: "bar2=baz2"}}
	ups3 := &noopUpstream{rawURL: "http://test3.com", weight: 1, parsedURL: &url.URL{Scheme: "http", Host: "test.com", RawQuery: "bar3=baz3"}}
	url1 := &url.URL{Scheme: "http", Host: "test.com", Path: "/foo", RawPath: "/foo", RawQuery: "bar1=baz1"}
	url2 := &url.URL{Scheme: "http", Host: "test.com", Path: "/foo", RawPath: "/foo", RawQuery: "bar2=baz2"}
	url3 := &url.URL{Scheme: "http", Host: "test.com", Path: "/foo", RawPath: "/foo", RawQuery: "bar3=baz3"}
	hasher1 := resilience.NewHTTPHasher(&v1.HTTPHasherSpec{HasherType: v1.HTTPHasherType_Header, Key: "test", HashAlg: kernel.HashAlg_FNV1a_32})
	hasher2 := resilience.NewHTTPHasher(&v1.HTTPHasherSpec{HasherType: v1.HTTPHasherType_Header, Key: "foo", HashAlg: kernel.HashAlg_FNV1a_32})
	hasher3 := resilience.NewHTTPHasher(&v1.HTTPHasherSpec{HasherType: v1.HTTPHasherType_Header, Key: "test", HashAlg: kernel.HashAlg_FNV1a_128})

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"no upstream/no hasher",
			[]string{cndNoUpstream},
			[]string{actNoUpstream},
			&condition{
				upstreams: []upstream{},
				matcher:   func(s string) (string, bool) { return s, true },
			},
			&action{
				upstream: []upstream{nil},
				url:      []*url.URL{nil},
				matched:  []bool{true},
			},
		),
		gen(
			"single upstream/no hasher",
			[]string{cndSingleUpstream},
			[]string{},
			&condition{
				upstreams: []upstream{ups1},
				matcher:   func(s string) (string, bool) { return s, true },
			},
			&action{
				upstream: []upstream{nil},
				url:      []*url.URL{nil},
				matched:  []bool{true},
			},
		),
		gen(
			"multiple upstream/no hasher",
			[]string{cndMultiUpstream},
			[]string{},
			&condition{
				upstreams: []upstream{ups1, ups2, ups3},
				matcher:   func(s string) (string, bool) { return s, true },
			},
			&action{
				upstream: []upstream{nil},
				url:      []*url.URL{nil},
				matched:  []bool{true},
			},
		),
		gen(
			"no upstream/single hasher",
			[]string{cndNoUpstream},
			[]string{actNoUpstream},
			&condition{
				upstreams: []upstream{},
				matcher:   func(s string) (string, bool) { return s, true },
				hashers:   []resilience.HTTPHasher{hasher1},
			},
			&action{
				upstream: []upstream{nil},
				url:      []*url.URL{nil},
				matched:  []bool{true},
			},
		),
		gen(
			"single upstream/single hasher",
			[]string{cndSingleUpstream},
			[]string{},
			&condition{
				upstreams: []upstream{ups1},
				matcher:   func(s string) (string, bool) { return s, true },
				hashers:   []resilience.HTTPHasher{hasher1},
			},
			&action{
				upstream: []upstream{ups1},
				url:      []*url.URL{url1},
				matched:  []bool{true},
			},
		),
		gen(
			"multiple upstream/single hasher",
			[]string{cndMultiUpstream},
			[]string{},
			&condition{
				upstreams: []upstream{ups1, ups2, ups3},
				matcher:   func(s string) (string, bool) { return s, true },
				hashers:   []resilience.HTTPHasher{hasher1},
			},
			&action{
				upstream: []upstream{ups3},
				url:      []*url.URL{url3},
				matched:  []bool{true},
			},
		),
		gen(
			"no upstream/multiple hasher",
			[]string{cndNoUpstream},
			[]string{actNoUpstream},
			&condition{
				upstreams: []upstream{},
				matcher:   func(s string) (string, bool) { return s, true },
				hashers:   []resilience.HTTPHasher{hasher2, hasher3, hasher1},
			},
			&action{
				upstream: []upstream{nil},
				url:      []*url.URL{nil},
				matched:  []bool{true},
			},
		),
		gen(
			"single upstream/multiple hasher",
			[]string{cndSingleUpstream},
			[]string{},
			&condition{
				upstreams: []upstream{ups1},
				matcher:   func(s string) (string, bool) { return s, true },
				hashers:   []resilience.HTTPHasher{hasher2, hasher3, hasher1},
			},
			&action{
				upstream: []upstream{ups1},
				url:      []*url.URL{url1},
				matched:  []bool{true},
			},
		),
		gen(
			"multiple upstream/multiple hasher",
			[]string{cndMultiUpstream},
			[]string{},
			&condition{
				upstreams: []upstream{ups1, ups2, ups3},
				matcher:   func(s string) (string, bool) { return s, true },
				hashers:   []resilience.HTTPHasher{hasher2, hasher3, hasher1},
			},
			&action{
				upstream: []upstream{ups2},
				url:      []*url.URL{url2},
				matched:  []bool{true},
			},
		),
		gen(
			"single upstream/inactive upstream",
			[]string{cndSingleUpstream, cndInactiveUpstream},
			[]string{},
			&condition{
				upstreams: []upstream{
					&lbUpstream{
						weight:         1,
						circuitBreaker: &circuitBreakerController{status: opened},
						rawURL:         "http://inactive.com",
						parsedURL:      &url.URL{Scheme: "http", Host: "inactive.com"},
					},
				},
				matcher: func(s string) (string, bool) { return s, true },
				hashers: []resilience.HTTPHasher{hasher1},
			},
			&action{
				upstream: []upstream{nil},
				url:      []*url.URL{nil},
				matched:  []bool{true},
			},
		),
		gen(
			"multiple upstream/1 inactive upstream",
			[]string{cndMultiUpstream, cndInactiveUpstream},
			[]string{},
			&condition{
				upstreams: []upstream{
					&lbUpstream{
						weight:         1,
						circuitBreaker: &circuitBreakerController{status: opened},
						rawURL:         "http://inactive.com",
						parsedURL:      &url.URL{Scheme: "http", Host: "inactive.com"},
					},
					ups2, ups3},
				matcher: func(s string) (string, bool) { return s, true },
				hashers: []resilience.HTTPHasher{hasher2, hasher3, hasher1},
			},
			&action{
				upstream: []upstream{ups2},
				url:      []*url.URL{url2},
				matched:  []bool{true},
			},
		),
		gen(
			"multiple upstream/2 inactive upstream",
			[]string{cndMultiUpstream, cndInactiveUpstream},
			[]string{},
			&condition{
				upstreams: []upstream{
					&lbUpstream{
						weight:         1,
						circuitBreaker: &circuitBreakerController{status: opened},
						rawURL:         "http://inactive1.com",
						parsedURL:      &url.URL{Scheme: "http", Host: "inactive1.com"},
					},
					&lbUpstream{
						weight:         1,
						circuitBreaker: &circuitBreakerController{status: opened},
						rawURL:         "http://inactive2.com",
						parsedURL:      &url.URL{Scheme: "http", Host: "inactive2.com"},
					},
					ups3},
				matcher: func(s string) (string, bool) { return s, true },
				hashers: []resilience.HTTPHasher{hasher2, hasher3, hasher1},
			},
			&action{
				upstream: []upstream{ups3},
				url:      []*url.URL{url3},
				matched:  []bool{true},
			},
		),
		gen(
			"path not match",
			[]string{cndPathNotMatch},
			[]string{actNoUpstream},
			&condition{
				upstreams: []upstream{ups1},
				matcher:   func(string) (string, bool) { return "", false },
			},
			&action{
				upstream: []upstream{nil},
				url:      []*url.URL{nil},
				matched:  []bool{false},
			},
		),
		gen(
			"hasher failed",
			[]string{cndMultiUpstream},
			[]string{},
			&condition{
				upstreams: []upstream{ups1, ups2, ups3},
				matcher:   func(s string) (string, bool) { return s, true },
				hashers:   []resilience.HTTPHasher{hasher2},
			},
			&action{
				upstream: []upstream{nil},
				url:      []*url.URL{nil},
				matched:  []bool{true},
			},
		),
		gen(
			"all hashers failed",
			[]string{cndMultiUpstream},
			[]string{},
			&condition{
				upstreams: []upstream{ups1, ups2, ups3},
				matcher:   func(s string) (string, bool) { return s, true },
				hashers:   []resilience.HTTPHasher{hasher2, hasher2, hasher2},
			},
			&action{
				upstream: []upstream{nil},
				url:      []*url.URL{nil},
				matched:  []bool{true},
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			rlb := &resilience.MaglevLB[upstream]{}
			rlb.Add(tt.C().upstreams...)
			lb := &hashBasedLB{
				lbMatcher: &lbMatcher{
					pathMatchers: []matcherFunc{tt.C().matcher},
				},
				LoadBalancer: rlb,
				hashers:      tt.C().hashers,
			}

			r := httptest.NewRequest(http.MethodGet, "http://test.com/foo", nil)
			r.Header.Set("test", "hash input")

			for i := 0; i < len(tt.A().upstream); i++ {
				upstream, url, matched := lb.upstream(r)
				testutil.Diff(t, tt.A().upstream[i], upstream, cmp.AllowUnexported(lbUpstream{}, noopUpstream{}, atomic.Int32{}))
				testutil.Diff(t, tt.A().url[i], url)
				testutil.Diff(t, tt.A().matched[i], matched)
			}
		})
	}
}
