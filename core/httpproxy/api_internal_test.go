// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package httpproxy

import (
	"context"
	"net/http"
	"net/url"
	"regexp"
	"sync"
	"sync/atomic"
	"testing"

	v1 "github.com/aileron-gateway/aileron-gateway/apis/core/v1"
	k "github.com/aileron-gateway/aileron-gateway/apis/kernel"
	"github.com/aileron-gateway/aileron-gateway/core"
	"github.com/aileron-gateway/aileron-gateway/internal/testutil"
	"github.com/aileron-gateway/aileron-gateway/kernel/api"
	"github.com/aileron-gateway/aileron-gateway/kernel/log"
	"github.com/aileron-gateway/aileron-gateway/kernel/network"
	"github.com/aileron-gateway/aileron-gateway/kernel/txtutil"
	utilhttp "github.com/aileron-gateway/aileron-gateway/util/http"
	"github.com/aileron-projects/go/zx/zlb"
	"github.com/cespare/xxhash/v2"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func TestMutate(t *testing.T) {
	type condition struct {
		manifest protoreflect.ProtoMessage
	}

	type action struct {
		manifest protoreflect.ProtoMessage
	}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"mutate default",
			&condition{
				manifest: Resource.Default(),
			},
			&action{
				manifest: &v1.ReverseProxyHandler{
					APIVersion: apiVersion,
					Kind:       kind,
					Metadata: &k.Metadata{
						Namespace: "default",
						Name:      "default",
					},
					Spec: &v1.ReverseProxyHandlerSpec{},
				},
			},
		),
		gen(
			"mutate default",
			&condition{
				manifest: Resource.Default(),
			},
			&action{
				manifest: &v1.ReverseProxyHandler{
					APIVersion: apiVersion,
					Kind:       kind,
					Metadata: &k.Metadata{
						Namespace: "default",
						Name:      "default",
					},
					Spec: &v1.ReverseProxyHandlerSpec{},
				},
			},
		),
		gen(
			"mutate lb spec",
			&condition{
				manifest: &v1.ReverseProxyHandler{
					APIVersion: apiVersion,
					Kind:       kind,
					Metadata: &k.Metadata{
						Namespace: "default",
						Name:      "default",
					},
					Spec: &v1.ReverseProxyHandlerSpec{
						LoadBalancers: []*v1.LoadBalancerSpec{{
							Upstreams: []*v1.UpstreamSpec{{}},
						}},
					},
				},
			},
			&action{
				manifest: &v1.ReverseProxyHandler{
					APIVersion: apiVersion,
					Kind:       kind,
					Metadata: &k.Metadata{
						Namespace: "default",
						Name:      "default",
					},
					Spec: &v1.ReverseProxyHandlerSpec{
						LoadBalancers: []*v1.LoadBalancerSpec{{
							Upstreams: []*v1.UpstreamSpec{{
								Weight:        1,
								EnablePassive: false,
								EnableActive:  false,
								InitialDelay:  0, // In seconds.
							}},
						}},
					},
				},
			},
		),
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			manifest := Resource.Mutate(tt.C.manifest)

			opts := []cmp.Option{
				cmpopts.IgnoreUnexported(v1.ReverseProxyHandler{}, v1.ReverseProxyHandlerSpec{}),
				cmpopts.IgnoreUnexported(v1.LoadBalancerSpec{}, v1.UpstreamSpec{}, v1.PathMatcherSpec{}),
				cmpopts.IgnoreUnexported(k.Metadata{}, k.Status{}, k.Reference{}),
			}
			testutil.Diff(t, tt.A.manifest, manifest, opts...)
		})
	}
}

func TestCreate(t *testing.T) {
	type condition struct {
		manifest protoreflect.ProtoMessage
	}

	type action struct {
		rp         http.Handler
		err        any // error or errorutil.Kind
		errPattern *regexp.Regexp
	}

	server := api.NewContainerAPI()
	postTestResource(server, "roundTripper", http.DefaultTransport)

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"create with default manifest",
			&condition{
				manifest: Resource.Default(),
			},
			&action{
				rp: &reverseProxy{
					HandlerBase: &utilhttp.HandlerBase{},
					lg:          log.GlobalLogger(log.DefaultLoggerName),
					eh:          utilhttp.GlobalErrorHandler(utilhttp.DefaultErrorHandlerName),
					rt:          network.DefaultHTTPTransport,
					lbs:         []loadBalancer{},
				},
				err: nil,
			},
		),
		gen(
			"create with round tripper",
			&condition{
				manifest: &v1.ReverseProxyHandler{
					Metadata: &k.Metadata{},
					Spec: &v1.ReverseProxyHandlerSpec{
						RoundTripper: testResourceRef("roundTripper"),
					},
				},
			},
			&action{
				rp: &reverseProxy{
					HandlerBase: &utilhttp.HandlerBase{},
					lg:          log.GlobalLogger(log.DefaultLoggerName),
					eh:          utilhttp.GlobalErrorHandler(utilhttp.DefaultErrorHandlerName),
					rt:          http.DefaultTransport,
					lbs:         []loadBalancer{},
				},
				err: nil,
			},
		),
		gen(
			"fail to get round tripper",
			&condition{
				manifest: &v1.ReverseProxyHandler{
					Metadata: &k.Metadata{},
					Spec: &v1.ReverseProxyHandlerSpec{
						RoundTripper: &k.Reference{
							APIVersion: "wrong",
						},
					},
				},
			},
			&action{
				rp:         nil,
				err:        core.ErrCoreGenCreateObject,
				errPattern: regexp.MustCompile(core.ErrPrefix + `failed to create ReverseProxyHandler`),
			},
		),
		gen(
			"fail to refer tripperware",
			&condition{
				manifest: &v1.ReverseProxyHandler{
					Metadata: &k.Metadata{},
					Spec: &v1.ReverseProxyHandlerSpec{
						Tripperwares: []*k.Reference{
							{
								APIVersion: "wrong",
							},
						},
					},
				},
			},
			&action{
				rp:         nil,
				err:        core.ErrCoreGenCreateObject,
				errPattern: regexp.MustCompile(core.ErrPrefix + `failed to create ReverseProxyHandler`),
			},
		),
		gen(
			"fail to create load balancer",
			&condition{
				manifest: &v1.ReverseProxyHandler{
					Metadata: &k.Metadata{},
					Spec: &v1.ReverseProxyHandlerSpec{
						LoadBalancers: []*v1.LoadBalancerSpec{
							{
								PathMatcher: &v1.PathMatcherSpec{
									Match:     `[0-9a-`,
									MatchType: k.MatchType_Regex,
								},
								Upstreams: []*v1.UpstreamSpec{{URL: "http://test.com"}},
							},
						},
					},
				},
			},
			&action{
				rp:         nil,
				err:        core.ErrCoreGenCreateObject,
				errPattern: regexp.MustCompile(core.ErrPrefix + `failed to create ReverseProxyHandler`),
			},
		),
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			a := &API{}
			rp, err := a.Create(server, tt.C.manifest)
			testutil.DiffError(t, tt.A.err, tt.A.errPattern, err)

			opts := []cmp.Option{
				cmp.Comparer(testutil.ComparePointer[log.Logger]),
				cmp.Comparer(testutil.ComparePointer[*http.Transport]),
				cmp.AllowUnexported(utilhttp.DefaultErrorHandler{}),
				cmp.AllowUnexported(reverseProxy{}),
			}
			testutil.Diff(t, tt.A.rp, rp, opts...)
		})
	}
}

func TestNewLoadBalancers(t *testing.T) {
	type condition struct {
		specs []*v1.LoadBalancerSpec
	}
	type action struct {
		lbs        []loadBalancer
		err        any // error or errorutil.Kind
		errPattern *regexp.Regexp
	}

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
			"input nil",
			&condition{
				specs: nil,
			},
			&action{
				lbs: []loadBalancer{},
				err: nil,
			},
		),
		gen(
			"single path matcher",
			&condition{
				specs: []*v1.LoadBalancerSpec{
					{
						PathMatcher: &v1.PathMatcherSpec{
							Match:     "/",
							MatchType: k.MatchType_Prefix,
						},
					},
				},
			},
			&action{
				lbs: []loadBalancer{
					&loadbalancer{
						lbMatcher: &lbMatcher{
							pathMatchers:  []matcherFunc{(&matcher{pattern: "/"}).prefix},
							paramMatchers: []txtutil.Matcher[*http.Request]{},
						},
						LoadBalancer: zlb.NewBasicRoundRobin[upstream](),
					},
				},
				err: nil,
			},
		),
		gen(
			"multiple path matchers",
			&condition{
				specs: []*v1.LoadBalancerSpec{
					{
						PathMatcher: &v1.PathMatcherSpec{Match: "/", MatchType: k.MatchType_Prefix},
						PathMatchers: []*v1.PathMatcherSpec{
							{Match: "/foo", MatchType: k.MatchType_Suffix},
							{Match: "/bar", MatchType: k.MatchType_Contains},
						},
					},
				},
			},
			&action{
				lbs: []loadBalancer{
					&loadbalancer{
						lbMatcher: &lbMatcher{
							pathMatchers: []matcherFunc{
								(&matcher{pattern: "/"}).prefix,
								(&matcher{pattern: "/foo"}).suffix,
								(&matcher{pattern: "/bar"}).contains,
							},
							paramMatchers: []txtutil.Matcher[*http.Request]{},
						},
						LoadBalancer: zlb.NewBasicRoundRobin[upstream](),
					},
				},
				err: nil,
			},
		),
		gen(
			"roundrobin/1 host/1 method",
			&condition{
				specs: []*v1.LoadBalancerSpec{
					{
						PathMatcher: &v1.PathMatcherSpec{
							Match:     "/",
							MatchType: k.MatchType_Prefix,
						},
						Hosts:   []string{"test1.com"},
						Methods: []v1.HTTPMethod{v1.HTTPMethod_GET},
					},
				},
			},
			&action{
				lbs: []loadBalancer{
					&loadbalancer{
						lbMatcher: &lbMatcher{
							pathMatchers:  []matcherFunc{(&matcher{pattern: "/"}).prefix},
							hosts:         []string{"test1.com"},
							methods:       []string{http.MethodGet},
							paramMatchers: []txtutil.Matcher[*http.Request]{},
						},
						LoadBalancer: zlb.NewBasicRoundRobin[upstream](),
					},
				},
				err: nil,
			},
		),
		gen(
			"roundrobin/2 host2/2 methods",
			&condition{
				specs: []*v1.LoadBalancerSpec{
					{
						PathMatcher: &v1.PathMatcherSpec{
							Match:     "/",
							MatchType: k.MatchType_Prefix,
						},
						Hosts:   []string{"test1.com", "test2.com"},
						Methods: []v1.HTTPMethod{v1.HTTPMethod_GET, v1.HTTPMethod_HEAD},
					},
				},
			},
			&action{
				lbs: []loadBalancer{
					&loadbalancer{
						lbMatcher: &lbMatcher{
							pathMatchers:  []matcherFunc{(&matcher{pattern: "/"}).prefix},
							hosts:         []string{"test1.com", "test2.com"},
							methods:       []string{http.MethodGet, http.MethodHead},
							paramMatchers: []txtutil.Matcher[*http.Request]{},
						},
						LoadBalancer: zlb.NewBasicRoundRobin[upstream](),
					},
				},
				err: nil,
			},
		),
		gen(
			"roundrobin/param matcher",
			&condition{
				specs: []*v1.LoadBalancerSpec{
					{
						PathMatcher: &v1.PathMatcherSpec{
							Match:     "/",
							MatchType: k.MatchType_Prefix,
						},
						PathParamMatchers: []*v1.ParamMatcherSpec{{Key: "pathParam", Patterns: []string{"/pp"}, MatchType: k.MatchType_Exact}},
						HeaderMatchers:    []*v1.ParamMatcherSpec{{Key: "headerParam", Patterns: []string{"/hp"}, MatchType: k.MatchType_Prefix}},
						QueryMatchers:     []*v1.ParamMatcherSpec{{Key: "queryParam", Patterns: []string{"/qp"}, MatchType: k.MatchType_Suffix}},
					},
				},
			},
			&action{
				lbs: []loadBalancer{
					&loadbalancer{
						lbMatcher: &lbMatcher{
							pathMatchers: []matcherFunc{(&matcher{pattern: "/"}).prefix},
							paramMatchers: []txtutil.Matcher[*http.Request]{
								&headerMatcher{key: "Headerparam", f: mustMatcher(txtutil.MatchTypePrefix, "/hp")},
								&queryMatcher{key: "queryParam", f: mustMatcher(txtutil.MatchTypePrefix, "/qp")},
								&pathParamMatcher{key: "pathParam", f: mustMatcher(txtutil.MatchTypePrefix, "/pp")},
							},
						},
						LoadBalancer: zlb.NewBasicRoundRobin[upstream](),
					},
				},
				err: nil,
			},
		),
		gen(
			"roundrobin/0 upstream",
			&condition{
				specs: []*v1.LoadBalancerSpec{
					{
						PathMatcher: &v1.PathMatcherSpec{
							Match:     "/",
							MatchType: k.MatchType_Prefix,
						},
						Upstreams: nil,
					},
				},
			},
			&action{
				lbs: []loadBalancer{
					&loadbalancer{
						lbMatcher: &lbMatcher{
							pathMatchers:  []matcherFunc{(&matcher{pattern: "/"}).prefix},
							paramMatchers: []txtutil.Matcher[*http.Request]{},
						},
						LoadBalancer: zlb.NewBasicRoundRobin[upstream](),
					},
				},
				err: nil,
			},
		),
		gen(
			"roundrobin/1 upstream/weight 0",
			&condition{
				specs: []*v1.LoadBalancerSpec{
					{
						PathMatcher: &v1.PathMatcherSpec{
							Match:     "/",
							MatchType: k.MatchType_Prefix,
						},
						Upstreams: []*v1.UpstreamSpec{
							{URL: "http://test.com", Weight: 0},
						},
					},
				},
			},
			&action{
				lbs: []loadBalancer{
					&loadbalancer{
						lbMatcher: &lbMatcher{
							pathMatchers:  []matcherFunc{(&matcher{pattern: "/"}).prefix},
							paramMatchers: []txtutil.Matcher[*http.Request]{},
						},
						LoadBalancer: zlb.NewBasicRoundRobin[upstream](),
					},
				},
				err: nil,
			},
		),
		gen(
			"roundrobin/1 upstream/weight 1",
			&condition{
				specs: []*v1.LoadBalancerSpec{
					{
						PathMatcher: &v1.PathMatcherSpec{
							Match:     "/",
							MatchType: k.MatchType_Prefix,
						},
						Upstreams: []*v1.UpstreamSpec{
							{URL: "http://test.com", Weight: 1},
						},
					},
				},
			},
			&action{
				lbs: []loadBalancer{
					&loadbalancer{
						lbMatcher: &lbMatcher{
							pathMatchers:  []matcherFunc{(&matcher{pattern: "/"}).prefix},
							paramMatchers: []txtutil.Matcher[*http.Request]{},
						},
						LoadBalancer: zlb.NewBasicRoundRobin[upstream](),
					},
				},
				err: nil,
			},
		),
		gen(
			"roundrobin/2 upstream/weight contain 0",
			&condition{
				specs: []*v1.LoadBalancerSpec{
					{
						PathMatcher: &v1.PathMatcherSpec{
							Match:     "/",
							MatchType: k.MatchType_Prefix,
						},
						Upstreams: []*v1.UpstreamSpec{
							{URL: "http://test1.com", Weight: 1},
							{URL: "http://test2.com", Weight: 0},
						},
					},
				},
			},
			&action{
				lbs: []loadBalancer{
					&loadbalancer{
						lbMatcher: &lbMatcher{
							pathMatchers:  []matcherFunc{(&matcher{pattern: "/"}).prefix},
							paramMatchers: []txtutil.Matcher[*http.Request]{},
						},
						LoadBalancer: zlb.NewBasicRoundRobin[upstream](),
					},
				},
				err: nil,
			},
		),
		gen(
			"roundrobin/2 upstream/same weights",
			&condition{
				specs: []*v1.LoadBalancerSpec{
					{
						PathMatcher: &v1.PathMatcherSpec{
							Match:     "/",
							MatchType: k.MatchType_Prefix,
						},
						Upstreams: []*v1.UpstreamSpec{
							{URL: "http://test1.com", Weight: 1},
							{URL: "http://test2.com", Weight: 1},
						},
					},
				},
			},
			&action{
				lbs: []loadBalancer{
					&loadbalancer{
						lbMatcher: &lbMatcher{
							pathMatchers:  []matcherFunc{(&matcher{pattern: "/"}).prefix},
							paramMatchers: []txtutil.Matcher[*http.Request]{},
						},
						LoadBalancer: zlb.NewBasicRoundRobin[upstream](),
					},
				},
				err: nil,
			},
		),
		gen(
			"roundrobin/2 upstream/different weights",
			&condition{
				specs: []*v1.LoadBalancerSpec{
					{
						PathMatcher: &v1.PathMatcherSpec{
							Match:     "/",
							MatchType: k.MatchType_Prefix,
						},
						Upstreams: []*v1.UpstreamSpec{
							{URL: "http://test1.com", Weight: 1},
							{URL: "http://test2.com", Weight: 2},
						},
					},
				},
			},
			&action{
				lbs: []loadBalancer{
					&loadbalancer{
						lbMatcher: &lbMatcher{
							pathMatchers:  []matcherFunc{(&matcher{pattern: "/"}).prefix},
							paramMatchers: []txtutil.Matcher[*http.Request]{},
						},
						LoadBalancer: zlb.NewBasicRoundRobin[upstream](),
					},
				},
				err: nil,
			},
		),
		gen(
			"random",
			&condition{
				specs: []*v1.LoadBalancerSpec{
					{
						LBAlgorithm: v1.LBAlgorithm_Random,
						PathMatcher: &v1.PathMatcherSpec{
							Match:     "/",
							MatchType: k.MatchType_Prefix,
						},
						Hosts:   []string{"test1.com"},
						Methods: []v1.HTTPMethod{v1.HTTPMethod_GET},
					},
				},
			},
			&action{
				lbs: []loadBalancer{
					&loadbalancer{
						lbMatcher: &lbMatcher{
							pathMatchers:  []matcherFunc{(&matcher{pattern: "/"}).prefix},
							hosts:         []string{"test1.com"},
							methods:       []string{http.MethodGet},
							paramMatchers: []txtutil.Matcher[*http.Request]{},
						},
						LoadBalancer: zlb.NewRandomW[upstream](),
					},
				},
				err: nil,
			},
		),
		gen(
			"maglev",
			&condition{
				specs: []*v1.LoadBalancerSpec{
					{
						LBAlgorithm: v1.LBAlgorithm_Maglev,
						PathMatcher: &v1.PathMatcherSpec{
							Match:     "/",
							MatchType: k.MatchType_Prefix,
						},
						Hosts:   []string{"test1.com"},
						Methods: []v1.HTTPMethod{v1.HTTPMethod_GET},
					},
				},
			},
			&action{
				lbs: []loadBalancer{
					&loadbalancer{
						lbMatcher: &lbMatcher{
							pathMatchers:  []matcherFunc{(&matcher{pattern: "/"}).prefix},
							hosts:         []string{"test1.com"},
							methods:       []string{http.MethodGet},
							paramMatchers: []txtutil.Matcher[*http.Request]{},
						},
						LoadBalancer: zlb.NewMaglev[upstream](),
						hasher:       clientAddrHasher(""),
					},
				},
				err: nil,
			},
		),
		gen(
			"ring hash",
			&condition{
				specs: []*v1.LoadBalancerSpec{
					{
						LBAlgorithm: v1.LBAlgorithm_RingHash,
						PathMatcher: &v1.PathMatcherSpec{
							Match:     "/",
							MatchType: k.MatchType_Prefix,
						},
						Hosts:   []string{"test1.com"},
						Methods: []v1.HTTPMethod{v1.HTTPMethod_GET},
					},
				},
			},
			&action{
				lbs: []loadBalancer{
					&loadbalancer{
						lbMatcher: &lbMatcher{
							pathMatchers:  []matcherFunc{(&matcher{pattern: "/"}).prefix},
							hosts:         []string{"test1.com"},
							methods:       []string{http.MethodGet},
							paramMatchers: []txtutil.Matcher[*http.Request]{},
						},
						LoadBalancer: zlb.NewRingHash[upstream](),
						hasher:       clientAddrHasher(""),
					},
				},
				err: nil,
			},
		),
		gen(
			"direct hash",
			&condition{
				specs: []*v1.LoadBalancerSpec{
					{
						LBAlgorithm: v1.LBAlgorithm_DirectHash,
						PathMatcher: &v1.PathMatcherSpec{
							Match:     "/",
							MatchType: k.MatchType_Prefix,
						},
						Hosts:   []string{"test1.com"},
						Methods: []v1.HTTPMethod{v1.HTTPMethod_GET},
					},
				},
			},
			&action{
				lbs: []loadBalancer{
					&loadbalancer{
						lbMatcher: &lbMatcher{
							pathMatchers:  []matcherFunc{(&matcher{pattern: "/"}).prefix},
							hosts:         []string{"test1.com"},
							methods:       []string{http.MethodGet},
							paramMatchers: []txtutil.Matcher[*http.Request]{},
						},
						LoadBalancer: zlb.NewDirectHashW[upstream](),
						hasher:       clientAddrHasher(""),
					},
				},
				err: nil,
			},
		),
		gen(
			"path matcher create error",
			&condition{
				specs: []*v1.LoadBalancerSpec{
					{
						PathMatcher: &v1.PathMatcherSpec{
							Match:     "[0-9a-",
							MatchType: k.MatchType_Regex,
						},
					},
				},
			},
			&action{
				lbs:        nil,
				err:        core.ErrCoreGenCreateComponent,
				errPattern: regexp.MustCompile(core.ErrPrefix + `failed to create component. loadBalancer creation failed`),
			},
		),
		gen(
			"path param matcher error",
			&condition{
				specs: []*v1.LoadBalancerSpec{
					{
						PathMatcher: &v1.PathMatcherSpec{
							Match:     "/",
							MatchType: k.MatchType_Prefix,
						},
						PathParamMatchers: []*v1.ParamMatcherSpec{
							{Key: "foo", Patterns: []string{"[0-9a-"}, MatchType: k.MatchType_Regex},
						},
					},
				},
			},
			&action{
				lbs:        nil,
				err:        core.ErrCoreGenCreateComponent,
				errPattern: regexp.MustCompile(core.ErrPrefix + `failed to create component. invalid parameter matcher config`),
			},
		),
		gen(
			"header param matcher error",
			&condition{
				specs: []*v1.LoadBalancerSpec{
					{
						PathMatcher: &v1.PathMatcherSpec{
							Match:     "/",
							MatchType: k.MatchType_Prefix,
						},
						HeaderMatchers: []*v1.ParamMatcherSpec{
							{Key: "foo", Patterns: []string{"[0-9a-"}, MatchType: k.MatchType_Regex},
						},
					},
				},
			},
			&action{
				lbs:        nil,
				err:        core.ErrCoreGenCreateComponent,
				errPattern: regexp.MustCompile(core.ErrPrefix + `failed to create component. invalid parameter matcher config`),
			},
		),
		gen(
			"query param matcher error",
			&condition{
				specs: []*v1.LoadBalancerSpec{
					{
						PathMatcher: &v1.PathMatcherSpec{
							Match:     "/",
							MatchType: k.MatchType_Prefix,
						},
						QueryMatchers: []*v1.ParamMatcherSpec{
							{Key: "foo", Patterns: []string{"[0-9a-"}, MatchType: k.MatchType_Regex},
						},
					},
				},
			},
			&action{
				lbs:        nil,
				err:        core.ErrCoreGenCreateComponent,
				errPattern: regexp.MustCompile(core.ErrPrefix + `failed to create component. invalid parameter matcher config`),
			},
		),
		gen(
			"upstream create error",
			&condition{
				specs: []*v1.LoadBalancerSpec{
					{
						PathMatcher: &v1.PathMatcherSpec{
							Match:     "/",
							MatchType: k.MatchType_Prefix,
						},
						Upstreams: []*v1.UpstreamSpec{
							{
								URL:          "\n://test.com",
								Weight:       1,
								EnableActive: true,
							},
						},
					},
				},
			},
			&action{
				lbs:        nil,
				err:        core.ErrCoreGenCreateComponent,
				errPattern: regexp.MustCompile(core.ErrPrefix + `failed to create component. loadBalancer creation failed`),
			},
		),
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			lbs, err := newLoadBalancers(http.DefaultTransport, tt.C.specs)
			testutil.DiffError(t, tt.A.err, tt.A.errPattern, err)

			opts := []cmp.Option{
				cmp.AllowUnexported(lbMatcher{}, loadbalancer{}),
				cmp.AllowUnexported(noopUpstream{}),
				cmp.AllowUnexported(headerMatcher{}, queryMatcher{}, pathParamMatcher{}),
				cmp.AllowUnexported(sync.Mutex{}, sync.RWMutex{}, atomic.Int32{}),
				cmp.Comparer(testutil.ComparePointer[matcherFunc]),
				cmp.Comparer(testutil.ComparePointer[txtutil.MatchFunc[string]]),
				cmpopts.IgnoreTypes(zlb.BasicRoundRobin[upstream]{}),
				cmpopts.IgnoreTypes(zlb.RoundRobin[upstream]{}),
				cmpopts.IgnoreTypes(zlb.RandomW[upstream]{}),
				cmpopts.IgnoreTypes(zlb.Maglev[upstream]{}),
				cmpopts.IgnoreTypes(zlb.DirectHashW[upstream]{}),
				cmpopts.IgnoreTypes(zlb.RingHash[upstream]{}),
			}
			testutil.Diff(t, tt.A.lbs, lbs, opts...)
			// testutil.Diff(t, tt.A.upstreams, lbs., opts...)
		})
	}
}

func TestNewUpstreams(t *testing.T) {
	type condition struct {
		specs []*v1.UpstreamSpec
	}

	type action struct {
		ups       []upstream
		shouldErr bool
	}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"input nil",
			&condition{
				specs: nil,
			},
			&action{
				ups: nil,
			},
		),
		gen(
			"1 valid spec",
			&condition{
				specs: []*v1.UpstreamSpec{
					{
						URL:          "http://test.com/",
						EnableActive: false,
						Weight:       1,
					},
				},
			},
			&action{
				ups: []upstream{
					&noopUpstream{
						id:        xxhash.Sum64String("http://test.com"),
						weight:    1,
						rawURL:    "http://test.com",
						parsedURL: &url.URL{Scheme: "http", Host: "test.com"},
					},
				},
			},
		),
		gen(
			"multiple valid specs",
			&condition{
				specs: []*v1.UpstreamSpec{
					{
						URL:          "http://test.com/foo",
						EnableActive: false,
						Weight:       1,
					},
					{
						URL:          "http://test.com/bar",
						EnableActive: false,
						Weight:       2,
					},
				},
			},
			&action{
				ups: []upstream{
					&noopUpstream{
						id:        xxhash.Sum64String("http://test.com/foo"),
						weight:    1,
						rawURL:    "http://test.com/foo",
						parsedURL: &url.URL{Scheme: "http", Host: "test.com", Path: "/foo"},
					},
					&noopUpstream{
						id:        xxhash.Sum64String("http://test.com/bar"),
						weight:    2,
						rawURL:    "http://test.com/bar",
						parsedURL: &url.URL{Scheme: "http", Host: "test.com", Path: "/bar"},
					},
				},
			},
		),
		gen(
			"invalid spec",
			&condition{
				specs: []*v1.UpstreamSpec{
					{
						URL:    "\nhttp://test.com",
						Weight: 1,
					},
				},
			},
			&action{
				ups:       nil,
				shouldErr: true,
			},
		),
		gen(
			"weight 0",
			&condition{
				specs: []*v1.UpstreamSpec{
					{
						URL:          "http://test.com/",
						EnableActive: false,
						Weight:       0,
					},
				},
			},
			&action{
				ups: []upstream{
					&noopUpstream{
						id:        xxhash.Sum64String("http://test.com"),
						weight:    1,
						rawURL:    "http://test.com",
						parsedURL: &url.URL{Scheme: "http", Host: "test.com"},
					},
				},
			},
		),
		gen(
			"weight -1",
			&condition{
				specs: []*v1.UpstreamSpec{
					{
						URL:          "http://test.com/",
						EnableActive: false,
						Weight:       -1,
					},
				},
			},
			&action{
				ups: []upstream{},
			},
		),
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			ups, err := newUpstreams(http.DefaultTransport, tt.C.specs)
			if tt.A.shouldErr {
				testutil.Diff(t, true, err != nil)
			}
			opts := []cmp.Option{
				cmp.AllowUnexported(lbUpstream{}, noopUpstream{}),
				cmpopts.IgnoreFields(lbUpstream{}, "closer"),
			}
			testutil.Diff(t, tt.A.ups, ups, opts...)
		})
	}
}

func postTestResource(server api.API[*api.Request, *api.Response], name string, res any) {
	ref := testResourceRef(name)
	req := &api.Request{
		Method:  api.MethodPost,
		Key:     ref.APIVersion + "/" + ref.Kind + "/" + ref.Namespace + "/" + ref.Name,
		Content: res,
	}
	if _, err := server.Serve(context.Background(), req); err != nil {
		panic(err)
	}
}

func testResourceRef(name string) *k.Reference {
	return &k.Reference{
		APIVersion: "container/v1",
		Kind:       "Container",
		Namespace:  "test",
		Name:       name,
	}
}
