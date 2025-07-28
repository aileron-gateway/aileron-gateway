// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package httpclient

import (
	"net/http"
	"slices"

	v1 "github.com/aileron-gateway/aileron-gateway/apis/core/v1"
	"github.com/aileron-gateway/aileron-gateway/apis/kernel"
	"github.com/aileron-gateway/aileron-gateway/core"
	"github.com/aileron-gateway/aileron-gateway/kernel/api"
	"github.com/aileron-gateway/aileron-gateway/kernel/network"
	utilhttp "github.com/aileron-gateway/aileron-gateway/util/http"
	"google.golang.org/protobuf/reflect/protoreflect"
)

const (
	apiVersion = "core/v1"
	kind       = "HTTPClient"
	Key        = apiVersion + "/" + kind
)

var Resource api.Resource = &API{
	BaseResource: &api.BaseResource{
		DefaultProto: &v1.HTTPClient{
			APIVersion: apiVersion,
			Kind:       kind,
			Metadata: &kernel.Metadata{
				Namespace: "default",
				Name:      "default",
			},
			Spec: &v1.HTTPClientSpec{},
		},
	},
}

type API struct {
	*api.BaseResource
}

func (*API) Create(a api.API[*api.Request, *api.Response], msg protoreflect.ProtoMessage) (any, error) {
	c := msg.(*v1.HTTPClient)

	var roundTripper http.RoundTripper = network.DefaultHTTPTransport
	switch c.Spec.Transports.(type) {
	case *v1.HTTPClientSpec_HTTPTransportConfig:
		tp, err := network.HTTPTransport(c.Spec.GetHTTPTransportConfig())
		if err != nil {
			return nil, core.ErrCoreGenCreateObject.WithStack(err, map[string]any{"kind": kind})
		}
		roundTripper = tp
	case *v1.HTTPClientSpec_HTTP2TransportConfig:
		tp, err := network.HTTP2Transport(c.Spec.GetHTTP2TransportConfig())
		if err != nil {
			return nil, core.ErrCoreGenCreateObject.WithStack(err, map[string]any{"kind": kind})
		}
		roundTripper = tp
	case *v1.HTTPClientSpec_HTTP3TransportConfig:
		tp, err := network.HTTP3Transport(c.Spec.GetHTTP3TransportConfig())
		if err != nil {
			return nil, core.ErrCoreGenCreateObject.WithStack(err, map[string]any{"kind": kind})
		}
		roundTripper = tp
	}

	ts, err := api.ReferTypedObjects[core.Tripperware](a, c.Spec.Tripperwares...)
	if err != nil {
		return nil, core.ErrCoreGenCreateObject.WithStack(err, map[string]any{"kind": kind})
	}

	rc := c.Spec.RetryConfig
	if rc != nil && rc.MaxRetry > 0 {
		retryStatus := make([]int, len(rc.RetryStatusCodes))
		for i, s := range rc.RetryStatusCodes {
			retryStatus[i] = int(s)
		}
		slices.Sort(retryStatus) // slices.Compact requires sorted slice.
		ret := &retry{
			maxRetry:         int(rc.MaxRetry),
			maxContentLength: int64(rc.MaxContentLength),
			retryStatus:      slices.Clip(slices.Compact(retryStatus)), // Remove duplicates.
		}
		ts = append(ts, core.Tripperware(ret))
	}

	return utilhttp.TripperwareChain(ts, roundTripper), nil
}
