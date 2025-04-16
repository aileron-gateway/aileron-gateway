// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package httpproxy

import (
	"context"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	v1 "github.com/aileron-gateway/aileron-gateway/apis/core/v1"
	k "github.com/aileron-gateway/aileron-gateway/apis/kernel"
	"github.com/aileron-gateway/aileron-gateway/core"
	"github.com/aileron-gateway/aileron-gateway/kernel/api"
	"github.com/aileron-gateway/aileron-gateway/kernel/log"
	"github.com/aileron-gateway/aileron-gateway/kernel/network"
	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
	"github.com/aileron-gateway/aileron-gateway/kernel/txtutil"
	utilhttp "github.com/aileron-gateway/aileron-gateway/util/http"
	"github.com/aileron-gateway/aileron-gateway/util/resilience"
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

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	cndDefault := tb.Condition("default", "input default manifest")
	cndNoLB := tb.Condition("no LB", "default values are applied to LB upstream")
	actCheckMutated := tb.Action("check mutated", "check that the intended fields are mutated")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"mutate default",
			[]string{cndDefault},
			[]string{actCheckMutated},
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
			[]string{},
			[]string{},
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
			[]string{cndNoLB},
			[]string{actCheckMutated},
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
								CheckInterval: 1, // In seconds.
								NetworkType:   k.NetworkType_HTTP,
								Protocol:      "icmp",
								Address:       "",
								// CircuitBreaker: &v1.CircuitBreaker{
								// 	FailureThreshold:        3, // Count for ConsecutiveCounter, not percentage.
								// 	SuccessThreshold:        3, // Count for ConsecutiveCounter, not percentage.
								// 	EffectiveFailureSamples: 10,
								// 	EffectiveSuccessSamples: 10,
								// 	WaitDuration:            180, // In seconds.
								// 	CircuitBreakerCounter: &v1.CircuitBreaker_ConsecutiveCounter{
								// 		ConsecutiveCounter: &v1.ConsecutiveCounterSpec{},
								// 	},
								// },
							}},
						}},
					},
				},
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			manifest := Resource.Mutate(tt.C().manifest)

			opts := []cmp.Option{
				cmpopts.IgnoreUnexported(v1.ReverseProxyHandler{}, v1.ReverseProxyHandlerSpec{}),
				cmpopts.IgnoreUnexported(v1.LoadBalancerSpec{}, v1.UpstreamSpec{}, v1.PathMatcherSpec{}),
				cmpopts.IgnoreUnexported(v1.CircuitBreaker{}, v1.CircuitBreaker_ConsecutiveCounter{}, v1.ConsecutiveCounterSpec{}),
				cmpopts.IgnoreUnexported(k.Metadata{}, k.Status{}, k.Reference{}),
			}
			testutil.Diff(t, tt.A().manifest, manifest, opts...)
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

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	cndDefaultManifest := tb.Condition("default manifest", "input default manifest")
	cndInvalidReference := tb.Condition("invalid reference", "input an invalid reference")
	actCheckNoError := tb.Action("no error", "check that there is no error")
	actCheckError := tb.Action("error", "check that there is an error")
	table := tb.Build()

	server := api.NewContainerAPI()
	postTestResource(server, "roundTripper", http.DefaultTransport)

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"create with default manifest",
			[]string{cndDefaultManifest},
			[]string{actCheckNoError},
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
			[]string{},
			[]string{actCheckNoError},
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
			"fail to get errorhandler",
			[]string{cndInvalidReference},
			[]string{actCheckError},
			&condition{
				manifest: &v1.ReverseProxyHandler{
					Metadata: &k.Metadata{},
					Spec: &v1.ReverseProxyHandlerSpec{
						ErrorHandler: &k.Reference{
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
			"fail to get round tripper",
			[]string{cndInvalidReference},
			[]string{actCheckError},
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
			[]string{cndInvalidReference},
			[]string{actCheckError},
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
			[]string{},
			[]string{actCheckError},
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

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			a := &API{}
			rp, err := a.Create(server, tt.C().manifest)
			testutil.DiffError(t, tt.A().err, tt.A().errPattern, err)

			opts := []cmp.Option{
				cmp.Comparer(testutil.ComparePointer[log.Logger]),
				cmp.Comparer(testutil.ComparePointer[*http.Transport]),
				cmp.AllowUnexported(utilhttp.DefaultErrorHandler{}),
				cmp.AllowUnexported(reverseProxy{}),
			}
			testutil.Diff(t, tt.A().rp, rp, opts...)
		})
	}
}

func TestNewLoadBalancers(t *testing.T) {
	type condition struct {
		specs []*v1.LoadBalancerSpec
	}

	type action struct {
		lbs        []loadBalancer
		upstreams  []upstream
		err        any // error or errorutil.Kind
		errPattern *regexp.Regexp
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	cndInputNil := tb.Condition("input nil", "input nil specs")
	cndNoUpstream := tb.Condition("no upstream", "no upstream")
	cndWithHosts := tb.Condition("with hosts", "specs contains at least 1 host")
	cndWithMethods := tb.Condition("with methods", "specs contains at least 1 method")
	cndRoundRobin := tb.Condition("round robin", "round-robin load balancer")
	cndRandom := tb.Condition("random", "random load balancer")
	cndInvalidSpec := tb.Condition("invalid spec", "input an invalid spec which should result in an error")
	actCheckError := tb.Action("error", "check that there is an error")
	actCheckNoError := tb.Action("no error", "check that there is no error")
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
			"input nil",
			[]string{cndInputNil},
			[]string{actCheckNoError},
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
			[]string{cndNoUpstream, cndWithHosts, cndWithMethods},
			[]string{actCheckNoError},
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
					&nonHashLB{
						lbMatcher: &lbMatcher{
							pathMatchers:  []matcherFunc{(&matcher{pattern: "/"}).prefix},
							paramMatchers: []txtutil.Matcher[*http.Request]{},
						},
						LoadBalancer: &resilience.RoundRobinLB[upstream]{
							// Content not checked.
						},
					},
				},
				err: nil,
			},
		),
		gen(
			"multiple path matchers",
			[]string{cndNoUpstream, cndWithHosts, cndWithMethods},
			[]string{actCheckNoError},
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
					&nonHashLB{
						lbMatcher: &lbMatcher{
							pathMatchers: []matcherFunc{
								(&matcher{pattern: "/"}).prefix,
								(&matcher{pattern: "/foo"}).suffix,
								(&matcher{pattern: "/bar"}).contains,
							},
							paramMatchers: []txtutil.Matcher[*http.Request]{},
						},
						LoadBalancer: &resilience.RoundRobinLB[upstream]{
							// Content not checked.
						},
					},
				},
				err: nil,
			},
		),
		gen(
			"roundrobin/1 host/1 method",
			[]string{cndNoUpstream, cndWithHosts, cndWithMethods},
			[]string{actCheckNoError},
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
					&nonHashLB{
						lbMatcher: &lbMatcher{
							pathMatchers:  []matcherFunc{(&matcher{pattern: "/"}).prefix},
							hosts:         []string{"test1.com"},
							methods:       []string{http.MethodGet},
							paramMatchers: []txtutil.Matcher[*http.Request]{},
						},
						LoadBalancer: &resilience.RoundRobinLB[upstream]{
							// Content not checked.
						},
					},
				},
				err: nil,
			},
		),
		gen(
			"roundrobin/2 host2/2 methods",
			[]string{cndNoUpstream, cndWithHosts, cndWithMethods},
			[]string{actCheckNoError},
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
					&nonHashLB{
						lbMatcher: &lbMatcher{
							pathMatchers:  []matcherFunc{(&matcher{pattern: "/"}).prefix},
							hosts:         []string{"test1.com", "test2.com"},
							methods:       []string{http.MethodGet, http.MethodHead},
							paramMatchers: []txtutil.Matcher[*http.Request]{},
						},
						LoadBalancer: &resilience.RoundRobinLB[upstream]{
							// Content not checked.
						},
					},
				},
				err: nil,
			},
		),
		gen(
			"roundrobin/param matcher",
			[]string{cndNoUpstream},
			[]string{actCheckNoError},
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
					&nonHashLB{
						lbMatcher: &lbMatcher{
							pathMatchers: []matcherFunc{(&matcher{pattern: "/"}).prefix},
							paramMatchers: []txtutil.Matcher[*http.Request]{
								&headerMatcher{key: "Headerparam", f: mustMatcher(txtutil.MatchTypePrefix, "/hp")},
								&queryMatcher{key: "queryParam", f: mustMatcher(txtutil.MatchTypePrefix, "/qp")},
								&pathParamMatcher{key: "pathParam", f: mustMatcher(txtutil.MatchTypePrefix, "/pp")},
							},
						},
						LoadBalancer: &resilience.RoundRobinLB[upstream]{
							// Content not checked.
						},
					},
				},
				err: nil,
			},
		),
		gen(
			"roundrobin/0 upstream",
			[]string{cndNoUpstream},
			[]string{actCheckNoError},
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
					&nonHashLB{
						lbMatcher: &lbMatcher{
							pathMatchers:  []matcherFunc{(&matcher{pattern: "/"}).prefix},
							paramMatchers: []txtutil.Matcher[*http.Request]{},
						},
						LoadBalancer: &resilience.RoundRobinLB[upstream]{
							// Content not checked.
						},
					},
				},
				err: nil,
			},
		),
		gen(
			"roundrobin/1 upstream/weight 0",
			[]string{cndRoundRobin},
			[]string{actCheckNoError},
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
					&nonHashLB{
						lbMatcher: &lbMatcher{
							pathMatchers:  []matcherFunc{(&matcher{pattern: "/"}).prefix},
							paramMatchers: []txtutil.Matcher[*http.Request]{},
						},
						LoadBalancer: &resilience.RoundRobinLB[upstream]{
							// Content not checked.
						},
					},
				},
				err: nil,
			},
		),
		gen(
			"roundrobin/1 upstream/weight 1",
			[]string{cndRoundRobin},
			[]string{actCheckNoError},
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
					&nonHashLB{
						lbMatcher: &lbMatcher{
							pathMatchers:  []matcherFunc{(&matcher{pattern: "/"}).prefix},
							paramMatchers: []txtutil.Matcher[*http.Request]{},
						},
						LoadBalancer: &resilience.RoundRobinLB[upstream]{
							// Content not checked.
						},
					},
				},
				err: nil,
			},
		),
		gen(
			"roundrobin/2 upstream/weight contain 0",
			[]string{cndRoundRobin},
			[]string{actCheckNoError},
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
					&nonHashLB{
						lbMatcher: &lbMatcher{
							pathMatchers:  []matcherFunc{(&matcher{pattern: "/"}).prefix},
							paramMatchers: []txtutil.Matcher[*http.Request]{},
						},
						LoadBalancer: &resilience.RoundRobinLB[upstream]{
							// Content not checked.
						},
					},
				},
				err: nil,
			},
		),
		gen(
			"roundrobin/2 upstream/same weights",
			[]string{cndRoundRobin},
			[]string{actCheckNoError},
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
					&nonHashLB{
						lbMatcher: &lbMatcher{
							pathMatchers:  []matcherFunc{(&matcher{pattern: "/"}).prefix},
							paramMatchers: []txtutil.Matcher[*http.Request]{},
						},
						LoadBalancer: &resilience.RoundRobinLB[upstream]{
							// Content not checked.
						},
					},
				},
				err: nil,
			},
		),
		gen(
			"roundrobin/2 upstream/different weights",
			[]string{cndRoundRobin},
			[]string{actCheckNoError},
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
					&nonHashLB{
						lbMatcher: &lbMatcher{
							pathMatchers:  []matcherFunc{(&matcher{pattern: "/"}).prefix},
							paramMatchers: []txtutil.Matcher[*http.Request]{},
						},
						LoadBalancer: &resilience.RoundRobinLB[upstream]{
							// Content not checked.
						},
					},
				},
				err: nil,
			},
		),
		gen(
			"random",
			[]string{cndRandom},
			[]string{actCheckNoError},
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
					&nonHashLB{
						lbMatcher: &lbMatcher{
							pathMatchers:  []matcherFunc{(&matcher{pattern: "/"}).prefix},
							hosts:         []string{"test1.com"},
							methods:       []string{http.MethodGet},
							paramMatchers: []txtutil.Matcher[*http.Request]{},
						},
						LoadBalancer: &resilience.RandomLB[upstream]{
							// Content not checked.
						},
					},
				},
				err: nil,
			},
		),
		gen(
			"maglev",
			[]string{},
			[]string{actCheckNoError},
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
					&hashBasedLB{
						lbMatcher: &lbMatcher{
							pathMatchers:  []matcherFunc{(&matcher{pattern: "/"}).prefix},
							hosts:         []string{"test1.com"},
							methods:       []string{http.MethodGet},
							paramMatchers: []txtutil.Matcher[*http.Request]{},
						},
						LoadBalancer: &resilience.MaglevLB[upstream]{
							// Content not checked.
						},
						hashers: []resilience.HTTPHasher{},
					},
				},
				err: nil,
			},
		),
		gen(
			"ring hash",
			[]string{},
			[]string{actCheckNoError},
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
					&hashBasedLB{
						lbMatcher: &lbMatcher{
							pathMatchers:  []matcherFunc{(&matcher{pattern: "/"}).prefix},
							hosts:         []string{"test1.com"},
							methods:       []string{http.MethodGet},
							paramMatchers: []txtutil.Matcher[*http.Request]{},
						},
						LoadBalancer: &resilience.RingHashLB[upstream]{
							// Content not checked.
						},
						hashers: []resilience.HTTPHasher{},
					},
				},
				err: nil,
			},
		),
		gen(
			"direct hash",
			[]string{},
			[]string{actCheckNoError},
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
					&directHashLB{
						lbMatcher: &lbMatcher{
							pathMatchers:  []matcherFunc{(&matcher{pattern: "/"}).prefix},
							hosts:         []string{"test1.com"},
							methods:       []string{http.MethodGet},
							paramMatchers: []txtutil.Matcher[*http.Request]{},
						},
						LoadBalancer: &resilience.DirectHashLB[upstream]{
							// Content not checked.
						},
						hashers: []resilience.HTTPHasher{},
					},
				},
				err: nil,
			},
		),
		gen(
			"path matcher create error",
			[]string{cndInvalidSpec},
			[]string{actCheckError},
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
			[]string{cndInvalidSpec},
			[]string{actCheckError},
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
			[]string{cndInvalidSpec},
			[]string{actCheckError},
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
			[]string{cndInvalidSpec},
			[]string{actCheckError},
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
			[]string{cndInvalidSpec},
			[]string{actCheckError},
			&condition{
				specs: []*v1.LoadBalancerSpec{
					{
						PathMatcher: &v1.PathMatcherSpec{
							Match:     "/",
							MatchType: k.MatchType_Prefix,
						},
						Upstreams: []*v1.UpstreamSpec{
							{
								URL:          "http://test.com",
								Weight:       1,
								EnableActive: true,
								NetworkType:  k.NetworkType_HTTP,
								Address:      "INVALID \n",
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

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			lbs, err := newLoadBalancers(http.DefaultTransport, tt.C().specs)
			testutil.DiffError(t, tt.A().err, tt.A().errPattern, err)

			opts := []cmp.Option{
				cmp.AllowUnexported(lbMatcher{}, nonHashLB{}, hashBasedLB{}, directHashLB{}),
				cmp.AllowUnexported(noopUpstream{}),
				cmp.AllowUnexported(headerMatcher{}, queryMatcher{}, pathParamMatcher{}),
				cmp.AllowUnexported(sync.Mutex{}, sync.RWMutex{}, atomic.Int32{}),
				cmp.Comparer(testutil.ComparePointer[matcherFunc]),
				cmp.Comparer(testutil.ComparePointer[txtutil.MatchFunc[string]]),
				cmpopts.IgnoreTypes(resilience.RoundRobinLB[upstream]{}),
				cmpopts.IgnoreTypes(resilience.RandomLB[upstream]{}),
				cmpopts.IgnoreTypes(resilience.MaglevLB[upstream]{}),
				cmpopts.IgnoreTypes(resilience.DirectHashLB[upstream]{}),
				cmpopts.IgnoreTypes(resilience.RingHashLB[upstream]{}),
			}
			testutil.Diff(t, tt.A().lbs, lbs, opts...)
			// testutil.Diff(t, tt.A().upstreams, lbs., opts...)
		})
	}
}

func TestNewLBUpstreams(t *testing.T) {
	type condition struct {
		specs []*v1.UpstreamSpec
	}

	type action struct {
		ups        []upstream
		err        any // error or errorutil.Kind
		errPattern *regexp.Regexp
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	cndInputNil := tb.Condition("input nil", "input nil")
	cndValidSpecs := tb.Condition("valid specs", "input valid spec")
	actCheckNoError := tb.Action("no error", "check that no error was returned")
	actCheckError := tb.Action("error", "check that an error was returned")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"input nil",
			[]string{cndInputNil},
			[]string{actCheckNoError},
			&condition{
				specs: nil,
			},
			&action{
				ups: nil,
			},
		),
		gen(
			"1 valid spec",
			[]string{cndValidSpecs},
			[]string{actCheckNoError},
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
						weight:    1,
						rawURL:    "http://test.com",
						parsedURL: &url.URL{Scheme: "http", Host: "test.com"},
					},
				},
			},
		),
		gen(
			"multiple valid specs",
			[]string{cndValidSpecs},
			[]string{actCheckNoError},
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
						weight:    1,
						rawURL:    "http://test.com/foo",
						parsedURL: &url.URL{Scheme: "http", Host: "test.com", Path: "/foo"},
					},
					&noopUpstream{
						weight:    2,
						rawURL:    "http://test.com/bar",
						parsedURL: &url.URL{Scheme: "http", Host: "test.com", Path: "/bar"},
					},
				},
			},
		),
		gen(
			"invalid spec",
			[]string{},
			[]string{actCheckError},
			&condition{
				specs: []*v1.UpstreamSpec{
					{
						URL:          "http://test.com",
						EnableActive: true,
						NetworkType:  k.NetworkType_HTTP,
						Address:      "INVALID \n",
						Weight:       1,
					},
				},
			},
			&action{
				ups:        nil,
				err:        core.ErrCoreGenCreateComponent,
				errPattern: regexp.MustCompile(core.ErrPrefix + `failed to create component. .*net/url: invalid control character in URL`),
			},
		),
		gen(
			"weight 0",
			[]string{cndValidSpecs},
			[]string{actCheckNoError},
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
						weight:    1,
						rawURL:    "http://test.com",
						parsedURL: &url.URL{Scheme: "http", Host: "test.com"},
					},
				},
			},
		),
		gen(
			"weight -1",
			[]string{cndValidSpecs},
			[]string{actCheckNoError},
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

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			ups, err := newLBUpstreams(http.DefaultTransport, tt.C().specs)
			testutil.DiffError(t, tt.A().err, tt.A().errPattern, err)

			opts := []cmp.Option{
				cmp.AllowUnexported(lbUpstream{}, noopUpstream{}),
				cmpopts.IgnoreFields(lbUpstream{}, "closer"),
			}
			testutil.Diff(t, tt.A().ups, ups, opts...)

			if len(ups) == 0 {
				return
			}

			time.Sleep(20 * time.Millisecond) // Wait activeCheck goroutine to start
			for _, up := range ups {
				if u, ok := up.(*lbUpstream); ok {
					u.close() // Close goroutine
				}
			}
		})
	}
}

func TestNewLBUpstream(t *testing.T) {
	type condition struct {
		spec *v1.UpstreamSpec
	}

	type action struct {
		ups        upstream
		err        any // error or errorutil.Kind
		errPattern *regexp.Regexp
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	cndHTTP := tb.Condition("http", "specify network type HTTP")
	cndTCP := tb.Condition("tcp", "specify network type TCP")
	cndUDP := tb.Condition("udp", "specify network type UDP")
	cndIP := tb.Condition("ip", "specify network type IP")
	actCheckUpstream := tb.Action("check upstream", "check the returned upstream config")
	actCheckNoError := tb.Action("no error", "check that no error was returned")
	actCheckError := tb.Action("error", "check that there is an error")
	table := tb.Build()

	// Get available port for testing to avoid dialing unknown service.
	ln4, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		panic(err)
	}
	ln4.Close()
	testPort := strconv.Itoa(ln4.Addr().(*net.TCPAddr).Port)

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"noop upstream",
			[]string{},
			[]string{actCheckUpstream, actCheckNoError},
			&condition{
				spec: &v1.UpstreamSpec{
					URL:          "http://test.com",
					EnableActive: false,
				},
			},
			&action{
				ups: &noopUpstream{
					rawURL:    "http://test.com",
					parsedURL: &url.URL{Scheme: "http", Host: "test.com"},
				},
			},
		),
		gen(
			"url has trailing slash",
			[]string{},
			[]string{actCheckUpstream, actCheckNoError},
			&condition{
				spec: &v1.UpstreamSpec{
					URL:          "http://test.com/",
					EnableActive: false,
				},
			},
			&action{
				ups: &noopUpstream{
					rawURL:    "http://test.com", // Suffix "/" will be trimmed.
					parsedURL: &url.URL{Scheme: "http", Host: "test.com"},
				},
			},
		),
		gen(
			"invalid url",
			[]string{},
			[]string{actCheckUpstream, actCheckNoError},
			&condition{
				spec: &v1.UpstreamSpec{
					URL:          "http://test com",
					EnableActive: false,
				},
			},
			&action{
				ups:        nil,
				err:        core.ErrCoreGenCreateComponent,
				errPattern: regexp.MustCompile(core.ErrPrefix + `failed to create component.`),
			},
		),
		gen(
			"HTTP",
			[]string{cndHTTP},
			[]string{actCheckUpstream, actCheckNoError},
			&condition{
				spec: &v1.UpstreamSpec{
					URL:           "http://test.com/",
					EnablePassive: true,
					EnableActive:  true,
					Weight:        9,
					InitialDelay:  10,
					CheckInterval: 20,
					NetworkType:   k.NetworkType_HTTP,
					Address:       "http://127.0.0.1:" + testPort,
				},
			},
			&action{
				ups: &lbUpstream{
					weight:         9,
					rawURL:         "http://test.com",
					parsedURL:      &url.URL{Scheme: "http", Host: "test.com"},
					passiveEnabled: true,
					initialDelay:   time.Second * 10,
					interval:       time.Second * 20,
				},
			},
		),
		gen(
			"HTTP error",
			[]string{cndHTTP},
			[]string{actCheckUpstream, actCheckError},
			&condition{
				spec: &v1.UpstreamSpec{
					URL:           "http://test.com/",
					EnablePassive: true,
					EnableActive:  true,
					Weight:        9,
					InitialDelay:  10,
					CheckInterval: 20,
					NetworkType:   k.NetworkType_HTTP,
					Address:       "INVALID \n",
				},
			},
			&action{
				ups:        nil,
				err:        core.ErrCoreGenCreateComponent,
				errPattern: regexp.MustCompile(core.ErrPrefix + `failed to create component. .*net/url: invalid control character in URL`),
			},
		),
		gen(
			"TCP",
			[]string{cndTCP},
			[]string{actCheckUpstream, actCheckNoError},
			&condition{
				spec: &v1.UpstreamSpec{
					URL:           "http://test.com/",
					EnablePassive: true,
					EnableActive:  true,
					Weight:        9,
					InitialDelay:  10,
					CheckInterval: 20,
					NetworkType:   k.NetworkType_TCP,
					Address:       "127.0.0.1:" + testPort,
				},
			},
			&action{
				ups: &lbUpstream{
					weight:         9,
					rawURL:         "http://test.com",
					parsedURL:      &url.URL{Scheme: "http", Host: "test.com"},
					passiveEnabled: true,
					initialDelay:   time.Second * 10,
					interval:       time.Second * 20,
				},
			},
		),
		gen(
			"TCP4",
			[]string{cndTCP},
			[]string{actCheckUpstream, actCheckNoError},
			&condition{
				spec: &v1.UpstreamSpec{
					URL:           "http://test.com/",
					EnablePassive: true,
					EnableActive:  true,
					Weight:        9,
					InitialDelay:  10,
					CheckInterval: 20,
					NetworkType:   k.NetworkType_TCP4,
					Address:       "127.0.0.1:" + testPort,
				},
			},
			&action{
				ups: &lbUpstream{
					weight:         9,
					rawURL:         "http://test.com",
					parsedURL:      &url.URL{Scheme: "http", Host: "test.com"},
					passiveEnabled: true,
					initialDelay:   time.Second * 10,
					interval:       time.Second * 20,
				},
			},
		),
		gen(
			"TCP6",
			[]string{cndTCP},
			[]string{actCheckUpstream, actCheckNoError},
			&condition{
				spec: &v1.UpstreamSpec{
					URL:           "http://test.com/",
					EnablePassive: true,
					EnableActive:  true,
					Weight:        9,
					InitialDelay:  10,
					CheckInterval: 20,
					NetworkType:   k.NetworkType_TCP6,
					Address:       "[::1]:" + testPort,
				},
			},
			&action{
				ups: &lbUpstream{
					weight:         9,
					rawURL:         "http://test.com",
					parsedURL:      &url.URL{Scheme: "http", Host: "test.com"},
					passiveEnabled: true,
					initialDelay:   time.Second * 10,
					interval:       time.Second * 20,
				},
			},
		),
		gen(
			"TCP error",
			[]string{cndTCP},
			[]string{actCheckUpstream, actCheckError},
			&condition{
				spec: &v1.UpstreamSpec{
					URL:           "http://test.com/",
					EnablePassive: true,
					EnableActive:  true,
					Weight:        9,
					InitialDelay:  10,
					CheckInterval: 20,
					NetworkType:   k.NetworkType_TCP,
					Address:       "INVALID",
				},
			},
			&action{
				ups:        nil,
				err:        core.ErrCoreGenCreateComponent,
				errPattern: regexp.MustCompile(core.ErrPrefix + `failed to create component. .*INVALID: missing port in address`),
			},
		),
		gen(
			"UDP",
			[]string{cndUDP},
			[]string{actCheckUpstream, actCheckNoError},
			&condition{
				spec: &v1.UpstreamSpec{
					URL:           "http://test.com/",
					EnablePassive: true,
					EnableActive:  true,
					Weight:        9,
					InitialDelay:  10,
					CheckInterval: 20,
					NetworkType:   k.NetworkType_UDP,
					Address:       "127.0.0.1:" + testPort,
				},
			},
			&action{
				ups: &lbUpstream{
					weight:         9,
					rawURL:         "http://test.com",
					parsedURL:      &url.URL{Scheme: "http", Host: "test.com"},
					passiveEnabled: true,
					initialDelay:   time.Second * 10,
					interval:       time.Second * 20,
				},
			},
		),
		gen(
			"UDP4",
			[]string{cndUDP},
			[]string{actCheckUpstream, actCheckNoError},
			&condition{
				spec: &v1.UpstreamSpec{
					URL:           "http://test.com/",
					EnablePassive: true,
					EnableActive:  true,
					Weight:        9,
					InitialDelay:  10,
					CheckInterval: 20,
					NetworkType:   k.NetworkType_UDP4,
					Address:       "127.0.0.1:" + testPort,
				},
			},
			&action{
				ups: &lbUpstream{
					weight:         9,
					rawURL:         "http://test.com",
					parsedURL:      &url.URL{Scheme: "http", Host: "test.com"},
					passiveEnabled: true,
					initialDelay:   time.Second * 10,
					interval:       time.Second * 20,
				},
			},
		),
		gen(
			"UDP6",
			[]string{cndUDP},
			[]string{actCheckUpstream, actCheckNoError},
			&condition{
				spec: &v1.UpstreamSpec{
					URL:           "http://test.com/",
					EnablePassive: true,
					EnableActive:  true,
					Weight:        9,
					InitialDelay:  10,
					CheckInterval: 20,
					NetworkType:   k.NetworkType_UDP6,
					Address:       "[::1]:" + testPort,
				},
			},
			&action{
				ups: &lbUpstream{
					weight:         9,
					rawURL:         "http://test.com",
					parsedURL:      &url.URL{Scheme: "http", Host: "test.com"},
					passiveEnabled: true,
					initialDelay:   time.Second * 10,
					interval:       time.Second * 20,
				},
			},
		),
		gen(
			"UDP error",
			[]string{cndUDP},
			[]string{actCheckUpstream, actCheckError},
			&condition{
				spec: &v1.UpstreamSpec{
					URL:           "http://test.com/",
					EnablePassive: true,
					EnableActive:  true,
					Weight:        9,
					InitialDelay:  10,
					CheckInterval: 20,
					NetworkType:   k.NetworkType_UDP,
					Address:       "INVALID",
				},
			},
			&action{
				ups:        nil,
				err:        core.ErrCoreGenCreateComponent,
				errPattern: regexp.MustCompile(core.ErrPrefix + `failed to create component. .*INVALID: missing port in address`),
			},
		),
		gen(
			"IP",
			[]string{cndIP},
			[]string{actCheckUpstream, actCheckNoError},
			&condition{
				spec: &v1.UpstreamSpec{
					URL:           "http://test.com/",
					EnablePassive: true,
					EnableActive:  true,
					Weight:        9,
					InitialDelay:  10,
					CheckInterval: 20,
					NetworkType:   k.NetworkType_IP,
					Address:       "127.0.0.1",
					Protocol:      "icmp",
				},
			},
			&action{
				ups: &lbUpstream{
					weight:         9,
					rawURL:         "http://test.com",
					parsedURL:      &url.URL{Scheme: "http", Host: "test.com"},
					passiveEnabled: true,
					initialDelay:   time.Second * 10,
					interval:       time.Second * 20,
				},
			},
		),
		gen(
			"IP4",
			[]string{cndIP},
			[]string{actCheckUpstream, actCheckNoError},
			&condition{
				spec: &v1.UpstreamSpec{
					URL:           "http://test.com/",
					EnablePassive: true,
					EnableActive:  true,
					Weight:        9,
					InitialDelay:  10,
					CheckInterval: 20,
					NetworkType:   k.NetworkType_IP4,
					Address:       "127.0.0.1",
					Protocol:      "icmp",
				},
			},
			&action{
				ups: &lbUpstream{
					weight:         9,
					rawURL:         "http://test.com",
					parsedURL:      &url.URL{Scheme: "http", Host: "test.com"},
					passiveEnabled: true,
					initialDelay:   time.Second * 10,
					interval:       time.Second * 20,
				},
			},
		),
		gen(
			"IP6",
			[]string{cndIP},
			[]string{actCheckUpstream, actCheckNoError},
			&condition{
				spec: &v1.UpstreamSpec{
					URL:           "http://test.com/",
					EnablePassive: true,
					EnableActive:  true,
					Weight:        9,
					InitialDelay:  10,
					CheckInterval: 20,
					NetworkType:   k.NetworkType_IP6,
					Address:       "::1",
					Protocol:      "icmp",
				},
			},
			&action{
				ups: &lbUpstream{
					weight:         9,
					rawURL:         "http://test.com",
					parsedURL:      &url.URL{Scheme: "http", Host: "test.com"},
					passiveEnabled: true,
					initialDelay:   time.Second * 10,
					interval:       time.Second * 20,
				},
			},
		),
		gen(
			"IP error",
			[]string{cndIP},
			[]string{actCheckUpstream, actCheckError},
			&condition{
				spec: &v1.UpstreamSpec{
					URL:           "http://test.com/",
					EnablePassive: true,
					EnableActive:  true,
					Weight:        9,
					InitialDelay:  10,
					CheckInterval: 20,
					NetworkType:   k.NetworkType_IP,
					Address:       "INVALID \n",
					Protocol:      "icmp",
				},
			},
			&action{
				ups:        nil,
				err:        core.ErrCoreGenCreateComponent,
				errPattern: regexp.MustCompile(core.ErrPrefix + `failed to create component. .*\[lookup INVALID`),
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			ups, err := newLBUpstream(http.DefaultTransport, tt.C().spec)
			testutil.DiffError(t, tt.A().err, tt.A().errPattern, err)

			if ups == nil {
				return
			}

			opts := []cmp.Option{
				cmp.AllowUnexported(lbUpstream{}, noopUpstream{}),
				cmpopts.IgnoreFields(lbUpstream{}, "closer"),
			}
			testutil.Diff(t, tt.A().ups, ups, opts...)

			if u, ok := ups.(*lbUpstream); ok {
				time.Sleep(20 * time.Millisecond) // Wait activeCheck goroutine to start
				u.close()                         // Close goroutine
			}
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
