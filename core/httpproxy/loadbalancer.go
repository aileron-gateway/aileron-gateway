// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package httpproxy

import (
	"net/http"
	"net/url"
	"slices"

	"github.com/aileron-gateway/aileron-gateway/internal/txtutil"
	"github.com/aileron-projects/go/zx/zlb"
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

type loadbalancer struct {
	// lbMatcher is the underlying matcher object.
	// Requests that this load balancer can accept are
	// determined by the match function of this object.
	// This must not be nil, otherwise panics.
	*lbMatcher
	//  LoadBalancer is the internal load balancer.
	// This LoadBalancer should not be a hash-based.
	zlb.LoadBalancer[upstream]
	// hasher is the hasher to calculate hash values of a request.
	hasher HTTPHasher
}

func (lb *loadbalancer) upstream(r *http.Request) (upstream, *url.URL, bool) {
	path, ok := lb.match(r)
	if !ok {
		return nil, nil, false
	}
	digest := uint64(0)
	if lb.hasher != nil {
		digest = lb.hasher.Hash(r) // hash value will be 0-65,535 when ok.
	}
	ups, found := lb.Get(digest)
	if !found || ups == nil || !ups.Active() {
		return nil, nil, true // Upstream not available.
	}
	u := toProxyURL(ups.url(), path)
	return ups, u, true // Upstream available.
}
