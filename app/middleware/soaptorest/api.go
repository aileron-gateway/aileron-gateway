package soaptorest

import (
	v1 "github.com/aileron-gateway/aileron-gateway/apis/app/v1"
	"github.com/aileron-gateway/aileron-gateway/apis/kernel"
	"github.com/aileron-gateway/aileron-gateway/core"
	"github.com/aileron-gateway/aileron-gateway/kernel/api"
	"github.com/aileron-gateway/aileron-gateway/kernel/log"
	"google.golang.org/protobuf/reflect/protoreflect"

	utilhttp "github.com/aileron-gateway/aileron-gateway/util/http"
)

const (
	apiVersion = "app/v1"
	kind       = "SoapToRestMiddleware"
	Key        = apiVersion + "/" + kind
)

var Resource api.Resource = &API{
	BaseResource: &api.BaseResource{
		DefaultProto: &v1.SoapToRestMiddleware{
			APIVersion: apiVersion,
			Kind:       kind,
			Metadata: &kernel.Metadata{
				Namespace: "default",
				Name:      "default",
			},
			Spec: &v1.SoapToRestMiddlewareSpec{},
		},
	},
}

type API struct {
	*api.BaseResource
}

func (*API) Create(a api.API[*api.Request, *api.Response], msg protoreflect.ProtoMessage) (any, error) {
	c := msg.(*v1.SoapToRestMiddleware)

	utilhttp.SetGlobalErrorHandler(utilhttp.DefaultErrorHandlerName, &soapErrorHandler{
		lg: log.GlobalLogger(log.DefaultLoggerName),
	})
	eh, err := utilhttp.ErrorHandler(a, c.Spec.ErrorHandler)
	if err != nil {
		return nil, core.ErrCoreGenCreateObject.WithStack(err, map[string]any{"kind": kind})
	}

	return &soapToRest{
		eh:            eh,
		attributeKey:  c.Spec.AttributeKey,
		namespaceKey:  c.Spec.NameSpaceKey,
		arrayKey:      c.Spec.ArrayKey,
		textKey:       c.Spec.TextKey,
		separatorChar: c.Spec.SeparatorChar,
	}, nil
}
