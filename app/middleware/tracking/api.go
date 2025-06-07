// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package tracking

import (
	"cmp"
	"net/textproto"

	v1 "github.com/aileron-gateway/aileron-gateway/apis/app/v1"
	"github.com/aileron-gateway/aileron-gateway/apis/kernel"
	"github.com/aileron-gateway/aileron-gateway/kernel/api"
	utilhttp "github.com/aileron-gateway/aileron-gateway/util/http"
	"google.golang.org/protobuf/reflect/protoreflect"
)

const (
	apiVersion = "app/v1"
	kind       = "TrackingMiddleware"
	Key        = apiVersion + "/" + kind
)

var Resource api.Resource = &API{
	BaseResource: &api.BaseResource{
		DefaultProto: &v1.TrackingMiddleware{
			APIVersion: apiVersion,
			Kind:       kind,
			Metadata: &kernel.Metadata{
				Namespace: "default",
				Name:      "default",
			},
			Spec: &v1.TrackingMiddlewareSpec{},
		},
	},
}

type API struct {
	*api.BaseResource
}

func (*API) Create(a api.API[*api.Request, *api.Response], msg protoreflect.ProtoMessage) (any, error) {
	c := msg.(*v1.TrackingMiddleware)
	eh := utilhttp.GlobalErrorHandler(cmp.Or(c.Metadata.ErrorHandler, utilhttp.DefaultErrorHandlerName))
	return &tracker{
		eh:               eh,
		reqProxyHeader:   textproto.CanonicalMIMEHeaderKey(c.Spec.RequestIDProxyName),
		trcProxyHeader:   textproto.CanonicalMIMEHeaderKey(c.Spec.TraceIDProxyName),
		trcExtractHeader: textproto.CanonicalMIMEHeaderKey(c.Spec.TraceIDExtractName),
	}, nil
}
