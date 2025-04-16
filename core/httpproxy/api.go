// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package httpproxy

import (
	"errors"
	"net"
	"net/http"
	"net/url"
	"slices"
	"strings"
	"time"

	v1 "github.com/aileron-gateway/aileron-gateway/apis/core/v1"
	"github.com/aileron-gateway/aileron-gateway/apis/kernel"
	"github.com/aileron-gateway/aileron-gateway/core"
	"github.com/aileron-gateway/aileron-gateway/kernel/api"
	"github.com/aileron-gateway/aileron-gateway/kernel/log"
	"github.com/aileron-gateway/aileron-gateway/kernel/network"
	utilhttp "github.com/aileron-gateway/aileron-gateway/util/http"
	"github.com/aileron-gateway/aileron-gateway/util/resilience"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

const (
	apiVersion = "core/v1"
	kind       = "ReverseProxyHandler"
	Key        = apiVersion + "/" + kind
)

var Resource api.Resource = &API{
	BaseResource: &api.BaseResource{
		DefaultProto: &v1.ReverseProxyHandler{
			APIVersion: apiVersion,
			Kind:       kind,
			Metadata: &kernel.Metadata{
				Namespace: "default",
				Name:      "default",
			},
			Spec: &v1.ReverseProxyHandlerSpec{},
		},
	},
}

type API struct {
	*api.BaseResource
}

func (*API) Mutate(msg protoreflect.ProtoMessage) protoreflect.ProtoMessage {
	c := msg.(*v1.ReverseProxyHandler)

	for _, spec := range c.Spec.LoadBalancers {
		for j, t := range spec.Upstreams {
			baseSpec := &v1.UpstreamSpec{
				Weight:        1,
				EnablePassive: false, // Passive health check.
				EnableActive:  false, // Active health check.
				InitialDelay:  0,     // Delay until stating health check in seconds.
				CheckInterval: 1,     // Health check interval in seconds.
				NetworkType:   kernel.NetworkType_HTTP,
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
			}
			proto.Merge(baseSpec, t)
			spec.Upstreams[j] = baseSpec
		}
	}

	return c
}

func (*API) Create(a api.API[*api.Request, *api.Response], msg protoreflect.ProtoMessage) (any, error) {
	c := msg.(*v1.ReverseProxyHandler)

	eh, err := utilhttp.ErrorHandler(a, c.Spec.ErrorHandler)
	if err != nil {
		return nil, core.ErrCoreGenCreateObject.WithStack(err, map[string]any{"kind": kind})
	}

	// Use http.DefaultTransport as the default round tripper.
	// Replace it if the c.Spec.RoundTripper is set.
	var roundTripper http.RoundTripper = network.DefaultHTTPTransport
	if c.Spec.RoundTripper != nil {
		rt, err := api.ReferTypedObject[http.RoundTripper](a, c.Spec.RoundTripper)
		if err != nil {
			return nil, core.ErrCoreGenCreateObject.WithStack(err, map[string]any{"kind": kind})
		}
		roundTripper = rt
	}

	ts, err := api.ReferTypedObjects[core.Tripperware](a, c.Spec.Tripperwares...)
	if err != nil {
		return nil, core.ErrCoreGenCreateObject.WithStack(err, map[string]any{"kind": kind})
	}

	lbs, err := newLoadBalancers(roundTripper, c.Spec.LoadBalancers)
	if err != nil {
		return nil, core.ErrCoreGenCreateObject.WithStack(err, map[string]any{"kind": kind})
	}

	return &reverseProxy{
		HandlerBase: &utilhttp.HandlerBase{
			AcceptPatterns: c.Spec.Patterns,
			AcceptMethods:  utilhttp.Methods(c.Spec.Methods),
		},
		lg:  log.DefaultOr(c.Metadata.Logger),
		eh:  eh,
		rt:  utilhttp.TripperwareChain(ts, roundTripper),
		lbs: lbs,
	}, nil
}

// newLoadBalancers returns load balancers.
// The given round tripper will be used for active health checking if enabled.
func newLoadBalancers(rt http.RoundTripper, specs []*v1.LoadBalancerSpec) ([]loadBalancer, error) {
	lbs := make([]loadBalancer, 0, len(specs))
	for _, spec := range specs {
		if spec.PathMatcher != nil {
			spec.PathMatchers = slices.Insert(spec.PathMatchers, 0, spec.PathMatcher)
		}
		var pathMatchers []matcherFunc
		for _, s := range spec.PathMatchers {
			mf, err := newMatcher(s)
			if err != nil {
				return nil, core.ErrCoreGenCreateComponent.WithStack(err, map[string]any{"reason": "loadBalancer creation failed"})
			}
			pathMatchers = append(pathMatchers, mf)
		}

		upstreams, err := newLBUpstreams(rt, spec.Upstreams)
		if err != nil {
			return nil, core.ErrCoreGenCreateComponent.WithStack(err, map[string]any{"reason": "loadBalancer creation failed"})
		}

		hMatchers, hErr := headerMatchers(spec.HeaderMatchers...)
		qMatchers, qErr := queryMatchers(spec.QueryMatchers...)
		pMatchers, pErr := pathParamMatchers(spec.PathParamMatchers...)
		if err = errors.Join(hErr, qErr, pErr); err != nil {
			return nil, core.ErrCoreGenCreateComponent.WithStack(err, map[string]any{"reason": "invalid parameter matcher config"})
		}
		matchers := slices.Clip(append(append(hMatchers, qMatchers...), pMatchers...))

		m := &lbMatcher{
			pathMatchers:  pathMatchers,
			methods:       utilhttp.Methods(spec.Methods),
			hosts:         slices.Clip(spec.Hosts),
			paramMatchers: matchers,
		}

		var lb loadBalancer
		switch spec.LBAlgorithm {
		case v1.LBAlgorithm_Maglev:
			rlb := &resilience.MaglevLB[upstream]{
				Size: int(spec.HashTableSize),
			}
			rlb.Add(upstreams...)
			lb = &hashBasedLB{
				lbMatcher:    m,
				LoadBalancer: rlb,
				hashers:      resilience.NewHTTPHashers(spec.Hashers),
			}
		case v1.LBAlgorithm_RingHash:
			rlb := &resilience.RingHashLB[upstream]{
				Size: int(spec.HashTableSize),
			}
			rlb.Add(upstreams...)
			lb = &hashBasedLB{
				lbMatcher:    m,
				LoadBalancer: rlb,
				hashers:      resilience.NewHTTPHashers(spec.Hashers),
			}
		case v1.LBAlgorithm_DirectHash:
			rlb := &resilience.DirectHashLB[upstream]{}
			rlb.Add(upstreams...)
			lb = &directHashLB{
				lbMatcher:    m,
				LoadBalancer: rlb,
				hashers:      resilience.NewHTTPHashers(spec.Hashers),
			}
		case v1.LBAlgorithm_Random:
			rlb := &resilience.RandomLB[upstream]{}
			rlb.Add(upstreams...)
			lb = &nonHashLB{
				lbMatcher:    m,
				LoadBalancer: rlb,
			}
		case v1.LBAlgorithm_RoundRobin:
			fallthrough // Use default.
		default:
			rlb := &resilience.RoundRobinLB[upstream]{}
			rlb.Add(upstreams...)
			lb = &nonHashLB{
				lbMatcher:    m,
				LoadBalancer: rlb,
			}
		}

		lbs = append(lbs, lb)
	}

	return slices.Clip(lbs), nil
}

// newLBUpstreams returns upstreams.
// newLBUpstreams ignore specs which weight is 0 or negative.
// The upstreams with weights less than or equal to 0 are ignored.
func newLBUpstreams(rt http.RoundTripper, specs []*v1.UpstreamSpec) ([]upstream, error) {
	if len(specs) == 0 {
		return nil, nil
	}
	ts := []upstream{}
	for _, spec := range specs {
		if spec.Weight < 0 {
			continue
		}
		if spec.Weight == 0 {
			spec.Weight = 1
		}
		t, err := newLBUpstream(rt, spec)
		if err != nil {
			return nil, err
		}
		ts = append(ts, t)
	}
	return slices.Clip(ts), nil
}

// newLBUpstream returns a new upstream object from given spec.
// The given argument rt and spec must not be nil.
// This function panics if a nil value was given.
func newLBUpstream(rt http.RoundTripper, spec *v1.UpstreamSpec) (upstream, error) {
	rawURL := strings.TrimSuffix(spec.URL, "/")
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return nil, core.ErrCoreGenCreateComponent.WithStack(err, nil)
	}

	// If active health check is not enabled,
	// there is no need to use circuit breaker.
	// Just return noop upstream.
	if !spec.EnableActive {
		return &noopUpstream{
			weight:    int(spec.Weight),
			rawURL:    rawURL,
			parsedURL: parsedURL,
		}, nil
	}

	tt := &lbUpstream{
		circuitBreaker: newCircuitBreaker(spec.CircuitBreaker),
		weight:         int(spec.Weight),
		rawURL:         rawURL,
		parsedURL:      parsedURL,
		passiveEnabled: spec.EnablePassive,
		interval:       time.Second * time.Duration(spec.CheckInterval),
		initialDelay:   time.Second * time.Duration(spec.InitialDelay),
	}

	var t upstream = tt
	nw := network.NetworkType(spec.NetworkType)

	switch spec.NetworkType {
	case kernel.NetworkType_HTTP:
		r, err := http.NewRequest(http.MethodGet, spec.Address, nil)
		if err != nil {
			return nil, core.ErrCoreGenCreateComponent.WithStack(err, nil)
		}
		go tt.activeCheckHTTP(rt, r)

	case kernel.NetworkType_TCP, kernel.NetworkType_TCP4, kernel.NetworkType_TCP6:
		addr, err := net.ResolveTCPAddr(nw, spec.Address)
		if err != nil {
			return nil, core.ErrCoreGenCreateComponent.WithStack(err, nil)
		}
		go tt.activeCheck(nw, addr.String())

	case kernel.NetworkType_UDP, kernel.NetworkType_UDP4, kernel.NetworkType_UDP6:
		addr, err := net.ResolveUDPAddr(nw, spec.Address)
		if err != nil {
			return nil, core.ErrCoreGenCreateComponent.WithStack(err, nil)
		}
		go tt.activeCheck(nw, addr.String())

	case kernel.NetworkType_IP, kernel.NetworkType_IP4, kernel.NetworkType_IP6:
		if spec.Protocol != "" {
			nw = nw + ":" + spec.Protocol
		}
		addr, err := net.ResolveIPAddr(nw, spec.Address)
		if err != nil {
			return nil, core.ErrCoreGenCreateComponent.WithStack(err, nil)
		}
		go tt.activeCheck(nw, addr.String())
	}

	return t, nil
}
