package httpproxy

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"

	v1 "github.com/aileron-gateway/aileron-gateway/apis/core/v1"
	"github.com/aileron-gateway/aileron-gateway/apis/kernel"
	"github.com/aileron-gateway/aileron-gateway/util/resilience"
)

var (
	hasher1 = resilience.NewHTTPHasher(&v1.HTTPHasherSpec{Key: "test", HashAlg: kernel.HashAlg_FNV1_32})
	hasher2 = resilience.NewHTTPHasher(&v1.HTTPHasherSpec{Key: "foo", HashAlg: kernel.HashAlg_FNV1a_32})
	hasher3 = resilience.NewHTTPHasher(&v1.HTTPHasherSpec{Key: "test", HashAlg: kernel.HashAlg_FNV1_64})
)

func BenchmarkRoundRobin(b *testing.B) {
	ups1 := &noopUpstream{rawURL: "ups1", weight: 1, parsedURL: &url.URL{}}
	ups2 := &noopUpstream{rawURL: "ups2", weight: 2, parsedURL: &url.URL{}}
	ups3 := &noopUpstream{rawURL: "ups3", weight: 3, parsedURL: &url.URL{}}
	ups4 := &noopUpstream{rawURL: "ups4", weight: 4, parsedURL: &url.URL{}}
	ups5 := &noopUpstream{rawURL: "ups5", weight: 5, parsedURL: &url.URL{}}
	ups0 := &lbUpstream{
		circuitBreaker: &circuitBreakerController{status: opened},
		weight:         1,
		rawURL:         "inactive",
		parsedURL:      &url.URL{},
	}

	upstreams := []upstream{
		ups1, ups2, ups3, ups4, ups5, ups1, ups2, ups3, ups4, ups0,
		ups1, ups2, ups3, ups4, ups5, ups1, ups2, ups3, ups4, ups0,
		ups1, ups2, ups3, ups4, ups5, ups1, ups2, ups3, ups4, ups0,
		ups1, ups2, ups3, ups4, ups5, ups1, ups2, ups3, ups4, ups0,
		ups1, ups2, ups3, ups4, ups5, ups1, ups2, ups3, ups4, ups0,

		ups1, ups2, ups3, ups4, ups5, ups1, ups2, ups3, ups4, ups0,
		ups1, ups2, ups3, ups4, ups5, ups1, ups2, ups3, ups4, ups0,
		ups1, ups2, ups3, ups4, ups5, ups1, ups2, ups3, ups4, ups0,
		ups1, ups2, ups3, ups4, ups5, ups1, ups2, ups3, ups4, ups0,
		ups1, ups2, ups3, ups4, ups5, ups1, ups2, ups3, ups4, ups0,
	}

	rlb := &resilience.RoundRobinLB[upstream]{}
	rlb.Add(upstreams...)
	lb := &nonHashLB{
		lbMatcher: &lbMatcher{
			pathMatchers: []matcherFunc{func(string) (string, bool) { return "", true }},
		},
		LoadBalancer: rlb,
	}

	r := httptest.NewRequest(http.MethodGet, "http://test.com/", nil)

	// all := map[string]int{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r.Header.Set("test", strconv.Itoa(i))
		lb.upstream(r)
		// u, _, _ := lb.upstream(r)
		// if u != nil {
		// 	all[u.ID()] += 1
		// } else {
		// 	all["inactive"] += 1
		// }
	}

	// b.Logf("%#v\n", all)
}

func BenchmarkRandom(b *testing.B) {
	ups1 := &noopUpstream{rawURL: "ups1", weight: 1, parsedURL: &url.URL{}}
	ups2 := &noopUpstream{rawURL: "ups2", weight: 2, parsedURL: &url.URL{}}
	ups3 := &noopUpstream{rawURL: "ups3", weight: 3, parsedURL: &url.URL{}}
	ups4 := &noopUpstream{rawURL: "ups4", weight: 4, parsedURL: &url.URL{}}
	ups5 := &noopUpstream{rawURL: "ups5", weight: 5, parsedURL: &url.URL{}}
	ups0 := &lbUpstream{
		circuitBreaker: &circuitBreakerController{status: opened},
		weight:         1,
		rawURL:         "inactive",
		parsedURL:      &url.URL{},
	}

	upstreams := []upstream{
		ups1, ups2, ups3, ups4, ups5, ups1, ups2, ups3, ups4, ups0,
		ups1, ups2, ups3, ups4, ups5, ups1, ups2, ups3, ups4, ups0,
		ups1, ups2, ups3, ups4, ups5, ups1, ups2, ups3, ups4, ups0,
		ups1, ups2, ups3, ups4, ups5, ups1, ups2, ups3, ups4, ups0,
		ups1, ups2, ups3, ups4, ups5, ups1, ups2, ups3, ups4, ups0,

		ups1, ups2, ups3, ups4, ups5, ups1, ups2, ups3, ups4, ups0,
		ups1, ups2, ups3, ups4, ups5, ups1, ups2, ups3, ups4, ups0,
		ups1, ups2, ups3, ups4, ups5, ups1, ups2, ups3, ups4, ups0,
		ups1, ups2, ups3, ups4, ups5, ups1, ups2, ups3, ups4, ups0,
		ups1, ups2, ups3, ups4, ups5, ups1, ups2, ups3, ups4, ups0,
	}

	rlb := &resilience.RandomLB[upstream]{}
	rlb.Add(upstreams...)

	lb := &nonHashLB{
		lbMatcher: &lbMatcher{
			pathMatchers: []matcherFunc{func(string) (string, bool) { return "", true }},
		},
		LoadBalancer: rlb,
	}

	r := httptest.NewRequest(http.MethodGet, "http://test.com/", nil)

	// all := map[string]int{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r.Header.Set("test", strconv.Itoa(i))
		lb.upstream(r)
		// u, _, _ := lb.upstream(r)
		// if u != nil {
		// 	all[u.ID()] += 1
		// } else {
		// 	all["inactive"] += 1
		// }
	}

	// b.Logf("%#v\n", all)
}

func BenchmarkDirectHash(b *testing.B) {
	ups1 := &noopUpstream{rawURL: "ups1", weight: 1, parsedURL: &url.URL{}}
	ups2 := &noopUpstream{rawURL: "ups2", weight: 2, parsedURL: &url.URL{}}
	ups3 := &noopUpstream{rawURL: "ups3", weight: 3, parsedURL: &url.URL{}}
	ups4 := &noopUpstream{rawURL: "ups4", weight: 4, parsedURL: &url.URL{}}
	ups5 := &noopUpstream{rawURL: "ups5", weight: 5, parsedURL: &url.URL{}}
	ups0 := &lbUpstream{
		circuitBreaker: &circuitBreakerController{status: opened},
		weight:         1,
		rawURL:         "inactive",
		parsedURL:      &url.URL{},
	}

	upstreams := []upstream{
		ups1, ups2, ups3, ups4, ups5, ups1, ups2, ups3, ups4, ups0,
		ups1, ups2, ups3, ups4, ups5, ups1, ups2, ups3, ups4, ups0,
		ups1, ups2, ups3, ups4, ups5, ups1, ups2, ups3, ups4, ups0,
		ups1, ups2, ups3, ups4, ups5, ups1, ups2, ups3, ups4, ups0,
		ups1, ups2, ups3, ups4, ups5, ups1, ups2, ups3, ups4, ups0,

		ups1, ups2, ups3, ups4, ups5, ups1, ups2, ups3, ups4, ups0,
		ups1, ups2, ups3, ups4, ups5, ups1, ups2, ups3, ups4, ups0,
		ups1, ups2, ups3, ups4, ups5, ups1, ups2, ups3, ups4, ups0,
		ups1, ups2, ups3, ups4, ups5, ups1, ups2, ups3, ups4, ups0,
		ups1, ups2, ups3, ups4, ups5, ups1, ups2, ups3, ups4, ups0,
	}

	rlb := &resilience.DirectHashLB[upstream]{}
	rlb.Add(upstreams...)

	lb := &directHashLB{
		lbMatcher: &lbMatcher{
			pathMatchers: []matcherFunc{func(string) (string, bool) { return "", true }},
		},
		LoadBalancer: rlb,
		hashers:      []resilience.HTTPHasher{hasher1, hasher2, hasher3},
	}

	r := httptest.NewRequest(http.MethodGet, "http://test.com/", nil)

	// all := map[string]int{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r.Header.Set("test", strconv.Itoa(i))
		lb.upstream(r)
		// u, _, _ := lb.upstream(r)
		// if u != nil {
		// 	all[u.ID()] += 1
		// } else {
		// 	all["inactive"] += 1
		// }
	}

	// b.Logf("%#v\n", all)
}

func BenchmarkRingHash(b *testing.B) {
	ups1 := &noopUpstream{rawURL: "ups1", weight: 1, parsedURL: &url.URL{}}
	ups2 := &noopUpstream{rawURL: "ups2", weight: 2, parsedURL: &url.URL{}}
	ups3 := &noopUpstream{rawURL: "ups3", weight: 3, parsedURL: &url.URL{}}
	ups4 := &noopUpstream{rawURL: "ups4", weight: 4, parsedURL: &url.URL{}}
	ups5 := &noopUpstream{rawURL: "ups5", weight: 5, parsedURL: &url.URL{}}
	ups0 := &lbUpstream{
		circuitBreaker: &circuitBreakerController{status: opened},
		weight:         1,
		rawURL:         "inactive",
		parsedURL:      &url.URL{},
	}

	upstreams := []upstream{
		ups1, ups2, ups3, ups4, ups5, ups1, ups2, ups3, ups4, ups0,
		ups1, ups2, ups3, ups4, ups5, ups1, ups2, ups3, ups4, ups0,
		ups1, ups2, ups3, ups4, ups5, ups1, ups2, ups3, ups4, ups0,
		ups1, ups2, ups3, ups4, ups5, ups1, ups2, ups3, ups4, ups0,
		ups1, ups2, ups3, ups4, ups5, ups1, ups2, ups3, ups4, ups0,

		ups1, ups2, ups3, ups4, ups5, ups1, ups2, ups3, ups4, ups0,
		ups1, ups2, ups3, ups4, ups5, ups1, ups2, ups3, ups4, ups0,
		ups1, ups2, ups3, ups4, ups5, ups1, ups2, ups3, ups4, ups0,
		ups1, ups2, ups3, ups4, ups5, ups1, ups2, ups3, ups4, ups0,
		ups1, ups2, ups3, ups4, ups5, ups1, ups2, ups3, ups4, ups0,
	}

	rlb := &resilience.RingHashLB[upstream]{
		Size: 1_000_000,
	}
	rlb.Add(upstreams...)

	lb := &hashBasedLB{
		lbMatcher: &lbMatcher{
			pathMatchers: []matcherFunc{func(string) (string, bool) { return "", true }},
		},
		hashers:      []resilience.HTTPHasher{hasher1, hasher2, hasher3},
		LoadBalancer: rlb,
	}

	r := httptest.NewRequest(http.MethodGet, "http://test.com/", nil)

	// all := map[string]int{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r.Header.Set("test", strconv.Itoa(i))
		lb.upstream(r)
		// u, _, _ := lb.upstream(r)
		// if u != nil {
		// 	all[u.ID()] += 1
		// } else {
		// 	all["inactive"] += 1
		// }
	}

	// b.Logf("%#v\n", all)
}

func BenchmarkMaglev(b *testing.B) {
	ups1 := &noopUpstream{rawURL: "ups1", weight: 1, parsedURL: &url.URL{}}
	ups2 := &noopUpstream{rawURL: "ups2", weight: 2, parsedURL: &url.URL{}}
	ups3 := &noopUpstream{rawURL: "ups3", weight: 3, parsedURL: &url.URL{}}
	ups4 := &noopUpstream{rawURL: "ups4", weight: 4, parsedURL: &url.URL{}}
	ups5 := &noopUpstream{rawURL: "ups5", weight: 5, parsedURL: &url.URL{}}
	ups0 := &lbUpstream{
		circuitBreaker: &circuitBreakerController{status: opened},
		weight:         1,
		rawURL:         "inactive",
		parsedURL:      &url.URL{},
	}

	upstreams := []upstream{
		ups1, ups2, ups3, ups4, ups5, ups1, ups2, ups3, ups4, ups0,
		ups1, ups2, ups3, ups4, ups5, ups1, ups2, ups3, ups4, ups0,
		ups1, ups2, ups3, ups4, ups5, ups1, ups2, ups3, ups4, ups0,
		ups1, ups2, ups3, ups4, ups5, ups1, ups2, ups3, ups4, ups0,
		ups1, ups2, ups3, ups4, ups5, ups1, ups2, ups3, ups4, ups0,

		ups1, ups2, ups3, ups4, ups5, ups1, ups2, ups3, ups4, ups0,
		ups1, ups2, ups3, ups4, ups5, ups1, ups2, ups3, ups4, ups0,
		ups1, ups2, ups3, ups4, ups5, ups1, ups2, ups3, ups4, ups0,
		ups1, ups2, ups3, ups4, ups5, ups1, ups2, ups3, ups4, ups0,
		ups1, ups2, ups3, ups4, ups5, ups1, ups2, ups3, ups4, ups0,
	}

	rlb := &resilience.MaglevLB[upstream]{}
	rlb.Add(upstreams...)

	lb := &hashBasedLB{
		lbMatcher: &lbMatcher{
			pathMatchers: []matcherFunc{func(string) (string, bool) { return "", true }},
		},
		LoadBalancer: rlb,
		hashers:      []resilience.HTTPHasher{hasher1, hasher2, hasher3},
	}

	r := httptest.NewRequest(http.MethodGet, "http://test.com/", nil)

	// all := map[string]int{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r.Header.Set("test", strconv.Itoa(i))
		lb.upstream(r)
		// u, _, _ := lb.upstream(r)
		// if u != nil {
		// 	all[u.ID()] += 1
		// } else {
		// 	all["inactive"] += 1
		// }
	}

	// b.Logf("%#v\n", all)
}
