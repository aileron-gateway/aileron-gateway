// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package tracking

import (
	"net/textproto"

	v1 "github.com/aileron-gateway/aileron-gateway/apis/app/v1"
	"github.com/aileron-gateway/aileron-gateway/apis/kernel"
	"github.com/aileron-gateway/aileron-gateway/core"
	"github.com/aileron-gateway/aileron-gateway/kernel/api"
	"github.com/aileron-gateway/aileron-gateway/kernel/encoder"
	"github.com/aileron-gateway/aileron-gateway/kernel/uid"
	httputil "github.com/aileron-gateway/aileron-gateway/util/http"
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
			Spec: &v1.TrackingMiddlewareSpec{
				Encoding: kernel.EncodingType_Base32HexEscaped,
			},
		},
	},
}

type API struct {
	*api.BaseResource
}

func (*API) Create(a api.API[*api.Request, *api.Response], msg protoreflect.ProtoMessage) (any, error) {
	c := msg.(*v1.TrackingMiddleware)

	// Obtain an error handler.
	// Default error handler is returned when not configured.
	eh, err := httputil.ErrorHandler(a, c.Spec.ErrorHandler)
	if err != nil {
		return nil, core.ErrCoreGenCreateObject.WithStack(err, map[string]any{"kind": kind})
	}

	t := &tracker{
		eh: eh,

		reqProxyHeader:   textproto.CanonicalMIMEHeaderKey(c.Spec.RequestIDProxyName),
		trcProxyHeader:   textproto.CanonicalMIMEHeaderKey(c.Spec.TraceIDProxyName),
		trcExtractHeader: textproto.CanonicalMIMEHeaderKey(c.Spec.TraceIDExtractName),

		newReqID: NewRequestID,
		newTrcID: NewTraceID,
	}

	if t.newReqID == nil {
		t.newReqID = newReqIDFunc(c.Spec.Encoding)
	}
	if t.newTrcID == nil {
		t.newTrcID = newTraceID
	}

	return t, nil
}

func newReqIDFunc(e kernel.EncodingType) func() (string, error) {
	enc, _ := encoder.EncoderDecoder(e)
	return func() (string, error) {
		id, err := uid.NewHostedID()
		if err != nil {
			return "", err
		}
		return enc(id), nil
	}
}

func newTraceID(reqID string) (string, error) {
	return reqID, nil
}
