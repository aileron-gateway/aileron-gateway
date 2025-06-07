// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package timeout

import (
	"cmp"
	"time"

	v1 "github.com/aileron-gateway/aileron-gateway/apis/app/v1"
	"github.com/aileron-gateway/aileron-gateway/apis/kernel"
	"github.com/aileron-gateway/aileron-gateway/core"
	"github.com/aileron-gateway/aileron-gateway/kernel/api"
	"github.com/aileron-gateway/aileron-gateway/kernel/txtutil"
	utilhttp "github.com/aileron-gateway/aileron-gateway/util/http"
	"google.golang.org/protobuf/reflect/protoreflect"
)

const (
	apiVersion = "app/v1"
	kind       = "TimeoutMiddleware"
	Key        = apiVersion + "/" + kind
)

var Resource api.Resource = &API{
	BaseResource: &api.BaseResource{
		DefaultProto: &v1.TimeoutMiddleware{
			APIVersion: apiVersion,
			Kind:       kind,
			Metadata: &kernel.Metadata{
				Namespace: "default",
				Name:      "default",
			},
			Spec: &v1.TimeoutMiddlewareSpec{},
		},
	},
}

type API struct {
	*api.BaseResource
}

func (*API) Create(a api.API[*api.Request, *api.Response], msg protoreflect.ProtoMessage) (any, error) {
	c := msg.(*v1.TimeoutMiddleware)
	eh := utilhttp.GlobalErrorHandler(cmp.Or(c.Metadata.ErrorHandler, utilhttp.DefaultErrorHandlerName))

	timeouts, err := apiTimeouts(c.Spec.APITimeouts...)
	if err != nil {
		return nil, core.ErrCoreGenCreateObject.WithStack(err, map[string]any{"kind": kind})
	}

	return &timeout{
		eh:             eh,
		defaultTimeout: time.Millisecond * time.Duration(c.Spec.DefaultTimeout),
		apiTimeouts:    timeouts,
	}, nil
}

func apiTimeouts(specs ...*v1.APITimeoutSpec) ([]*apiTimeout, error) {
	timeouts := make([]*apiTimeout, 0, len(specs))
	for _, spec := range specs {
		if spec == nil || spec.Timeout <= 0 || spec.Matcher == nil {
			continue
		}

		m, err := txtutil.NewStringMatcher(txtutil.MatchTypes[spec.Matcher.MatchType], spec.Matcher.Patterns...)
		if err != nil {
			return nil, err // Return err as-is.
		}

		to := &apiTimeout{
			methods: utilhttp.Methods(spec.Methods),
			paths:   m,
			timeout: time.Millisecond * time.Duration(spec.Timeout),
		}

		timeouts = append(timeouts, to)
	}

	return timeouts, nil
}
