// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package httpproxy

import (
	"net/http"
	"net/url"
	"slices"

	"github.com/aileron-gateway/aileron-gateway/kernel/txtutil"
	"github.com/aileron-gateway/aileron-gateway/util/resilience"
)

// loadBalancer is the interface of load balancers.
type loadBalancer interface {
	// upstream returns upstream url to send request.
	// The returned bool is true when the request should be
	// load balanced with this load balancer.
	// The returned upstream will be nil when no upstream found.
	upstream(*http.Request) (upstream, *url.URL, bool)
}

func toProxyURL(u *url.URL, newPath string) *url.URL {
	return &url.URL{
		Scheme:      u.Scheme,
		Host:        u.Host,
		Path:        newPath, // This field is replaced by resulting proxy path.
		RawPath:     newPath, // This field is replaced by resulting proxy path.
		RawQuery:    u.RawQuery,
		Fragment:    u.Fragment,
		RawFragment: u.RawFragment,
		// Other fields are not used for proxy.
	}
}

// lbMatcher is the request matcher that will be used
// from load balancers.
type lbMatcher struct {
	// pathMatcher is the url path matcher function.
	// The function should return true when this
	// load balancer can accept the target request.
	// This must not be nil, otherwise panics.
	pathMatchers []matcherFunc
	// methods is the HTTP method name list
	// that this load balancer can accept.
	// If empty, all methods are accepted.
	methods []string
	// hosts is the hostname list
	// that this load balancer can accept.
	// If empty, all hosts are accepted.
	hosts []string
	// paramMatchers is the matcher function for header, query
	// and path parameter.
	paramMatchers []txtutil.Matcher[*http.Request]
}

func (lb *lbMatcher) match(r *http.Request) (string, bool) {
	if len(lb.methods) > 0 && !slices.Contains(lb.methods, r.Method) {
		return "", false
	}
	if len(lb.hosts) > 0 && !slices.Contains(lb.hosts, r.Host) {
		return "", false
	}
	for _, f := range lb.paramMatchers {
		if !f.Match(r) {
			return "", false
		}
	}
	for _, m := range lb.pathMatchers {
		if path, ok := m(r.URL.Path); ok {
			return path, true
		}
	}
	return "", false
}

type nonHashLB struct {
	// lbMatcher is the underlying matcher object.
	// Requests that this load balancer can accept are
	// determined by the match function of this object.
	// This must not be nil, otherwise panics.
	*lbMatcher

	//  LoadBalancer is the internal load balancer.
	// This LoadBalancer should not be a hash-based.
	// Use one of following LBs.
	//  - resilience.RandomLB
	//  - resilience.RoundRobinLB
	resilience.LoadBalancer[upstream]
}

func (lb *nonHashLB) upstream(r *http.Request) (upstream, *url.URL, bool) {
	path, ok := lb.match(r)
	if !ok {
		return nil, nil, false
	}

	ups := lb.Get(-1)
	if ups == nil || !ups.Active() {
		return nil, nil, true // Upstream not available.
	}

	u := toProxyURL(ups.url(), path)
	return ups, u, true // Upstream available.
}

type directHashLB struct {
	// lbMatcher is the underlying matcher object.
	// Requests that this load balancer can accept are
	// determined by the match function of this object.
	// This must not be nil, otherwise panics.
	*lbMatcher

	//  LoadBalancer is the internal load balancer.
	// This LoadBalancer should be a hash-based.
	// Use following LB.
	//  - resilience.DirectHashLB
	resilience.LoadBalancer[upstream]

	// hashers is the list of hasher that will be used for
	// calculating hash values of each requests.
	hashers []resilience.HTTPHasher
}

func (lb *directHashLB) upstream(r *http.Request) (upstream, *url.URL, bool) {
	path, ok := lb.match(r)
	if !ok {
		return nil, nil, false
	}

	for _, h := range lb.hashers {
		val, ok := h.Hash(r) // hash value will be 0-65,535 when ok.
		if !ok {
			continue // Try next hash.
		}
		ups := lb.Get(val)
		if ups == nil || !ups.Active() {
			continue // Try next hash.
		}

		u := toProxyURL(ups.url(), path)
		return ups, u, true // Upstream available.
	}

	return nil, nil, true // Upstream not available
}

type hashBasedLB struct {
	// lbMatcher is the underlying matcher object.
	// Requests that this load balancer can accept are
	// determined by the match function of this object.
	// This must not be nil, otherwise panics.
	*lbMatcher

	//  LoadBalancer is the internal load balancer.
	// This LoadBalancer should be a hash-based.
	// Use one of following LBs.
	//  - resilience.MaglevLB
	//  - resilience.RingHashLB
	resilience.LoadBalancer[upstream]

	// hashers is the list of hasher that will be used for
	// calculating hash values of each requests.
	hashers []resilience.HTTPHasher
}

func (lb *hashBasedLB) upstream(r *http.Request) (upstream, *url.URL, bool) {
	path, ok := lb.match(r)
	if !ok {
		return nil, nil, false
	}

	ok = false
	var val int // hash value will be 0-65,535 when valid.
	for _, h := range lb.hashers {
		if val, ok = h.Hash(r); ok {
			break
		}
	}
	if !ok {
		return nil, nil, true // Hash not available.
	}

	ups := lb.Get(val)
	if ups == nil || !ups.Active() {
		return nil, nil, true // Upstream not available.
	}

	u := toProxyURL(ups.url(), path)
	return ups, u, true // Upstream available.
}
