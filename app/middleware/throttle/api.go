// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package throttle

import (
	"cmp"
	"time"

	v1 "github.com/aileron-gateway/aileron-gateway/apis/app/v1"
	"github.com/aileron-gateway/aileron-gateway/apis/kernel"
	"github.com/aileron-gateway/aileron-gateway/core"
	"github.com/aileron-gateway/aileron-gateway/kernel/api"
	"github.com/aileron-gateway/aileron-gateway/kernel/log"
	"github.com/aileron-gateway/aileron-gateway/kernel/txtutil"
	utilhttp "github.com/aileron-gateway/aileron-gateway/util/http"
	"github.com/aileron-projects/go/ztime/zrate"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

const (
	apiVersion = "app/v1"
	kind       = "ThrottleMiddleware"
	Key        = apiVersion + "/" + kind
)

var Resource api.Resource = &API{
	BaseResource: &api.BaseResource{
		DefaultProto: &v1.ThrottleMiddleware{
			APIVersion: apiVersion,
			Kind:       kind,
			Metadata: &kernel.Metadata{
				Namespace: "default",
				Name:      "default",
			},
			Spec: &v1.ThrottleMiddlewareSpec{},
		},
	},
}

type API struct {
	*api.BaseResource
}

// Mutate changes configured values.
// The values of the msg which is given as the argument is the merged message of default values and user defined values.
// Changes for the fields of msg in this function make the final values which will be the input for validate and create function.
// Default values for "repeated" or "oneof" fields can also be applied in this function if necessary.
// Please check msg!=nil and asserting the mgs does not panic even they won't from the view of overall architecture of the gateway.
func (*API) Mutate(msg protoreflect.ProtoMessage) protoreflect.ProtoMessage {
	c := msg.(*v1.ThrottleMiddleware)
	for _, t := range c.Spec.APIThrottlers {
		switch t := t.Throttlers.(type) {
		case *v1.APIThrottlerSpec_MaxConnections:
			baseSpec := &v1.MaxConnectionsSpec{
				MaxConns: 128,
			}
			proto.Merge(baseSpec, t.MaxConnections)
			t.MaxConnections = baseSpec
		case *v1.APIThrottlerSpec_FixedWindow:
			baseSpec := &v1.FixedWindowSpec{
				WindowSize: 1000, // 1000 ms
				Limit:      1000, // 1000 req / 1000ms
			}
			proto.Merge(baseSpec, t.FixedWindow)
			t.FixedWindow = baseSpec
		case *v1.APIThrottlerSpec_TokenBucket:
			baseSpec := &v1.TokenBucketSpec{
				BucketSize:   1000,
				FillInterval: 1000, // 1000 ms
				FillRate:     1000, // 1000 req / 1000 ms
			}
			proto.Merge(baseSpec, t.TokenBucket)
			t.TokenBucket = baseSpec
		case *v1.APIThrottlerSpec_LeakyBucket:
			baseSpec := &v1.LeakyBucketSpec{
				BucketSize:   1000,
				LeakInterval: 1, // 1 ms
			}
			proto.Merge(baseSpec, t.LeakyBucket)
			t.LeakyBucket = baseSpec
		}
	}
	return c
}

func (*API) Create(a api.API[*api.Request, *api.Response], msg protoreflect.ProtoMessage) (any, error) {
	c := msg.(*v1.ThrottleMiddleware)
	_ = log.DefaultOr(c.Metadata.Logger)
	eh := utilhttp.GlobalErrorHandler(cmp.Or(c.Metadata.ErrorHandler, utilhttp.DefaultErrorHandlerName))

	throttlers, err := apiThrottlers(c.Spec.APIThrottlers...)
	if err != nil {
		return nil, core.ErrCoreGenCreateObject.WithStack(err, map[string]any{"kind": kind})
	}

	return &throttle{
		eh:         eh,
		throttlers: throttlers,
	}, nil
}

// apiThrottlers returns new apiThrottlers.
func apiThrottlers(specs ...*v1.APIThrottlerSpec) ([]*apiThrottler, error) {
	ths := make([]*apiThrottler, 0, len(specs))
	for _, spec := range specs {
		if spec == nil || spec.Matcher == nil {
			continue
		}
		m, err := txtutil.NewStringMatcher(txtutil.MatchTypes[spec.Matcher.MatchType], spec.Matcher.Patterns...)
		if err != nil {
			return nil, err // Return err as-is.
		}

		var limiter zrate.Limiter
		switch spec.Throttlers.(type) {
		case *v1.APIThrottlerSpec_MaxConnections:
			s := spec.GetMaxConnections()
			limiter = zrate.NewConcurrentLimiter(int(s.MaxConns))
		case *v1.APIThrottlerSpec_FixedWindow:
			s := spec.GetFixedWindow()
			limiter = zrate.NewFixedWindowLimiterWidth(int(s.Limit), time.Millisecond*time.Duration(s.WindowSize))
		case *v1.APIThrottlerSpec_LeakyBucket:
			s := spec.GetLeakyBucket()
			limiter = zrate.NewLeakyBucketLimiter(int(s.BucketSize), time.Millisecond*time.Duration(s.LeakInterval))
		case *v1.APIThrottlerSpec_TokenBucket:
			s := spec.GetTokenBucket()
			limiter = zrate.NewTokenBucketInterval(int(s.BucketSize), int(s.FillRate), time.Millisecond*time.Duration(s.FillInterval))
		}

		ths = append(ths, &apiThrottler{
			limiter:  limiter,
			methods:  utilhttp.Methods(spec.Methods),
			paths:    m,
			maxRetry: int(spec.MaxRetry),
		})
	}
	return ths, nil
}
