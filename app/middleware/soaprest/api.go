package soaprest

import (
	v1 "github.com/aileron-gateway/aileron-gateway/apis/app/v1"
	"github.com/aileron-gateway/aileron-gateway/apis/kernel"
	"github.com/aileron-gateway/aileron-gateway/core"
	"github.com/aileron-gateway/aileron-gateway/kernel/api"
	"github.com/aileron-gateway/aileron-gateway/kernel/txtutil"
	"google.golang.org/protobuf/reflect/protoreflect"

	utilhttp "github.com/aileron-gateway/aileron-gateway/util/http"
)

const (
	apiVersion = "app/v1"
	kind       = "SOAPRESTMiddleware"
	Key        = apiVersion + "/" + kind
)

var Resource api.Resource = &API{
	BaseResource: &api.BaseResource{},
}

type API struct {
	*api.BaseResource
}

func (o *API) Default() protoreflect.ProtoMessage {
	return &v1.SOAPRESTMiddleware{
		APIVersion: apiVersion,
		Kind:       kind,
		Metadata: &kernel.Metadata{
			Namespace: "default",
			Name:      "default",
		},
		Spec: &v1.SOAPRESTMiddlewareSpec{
			Matcher: &kernel.MatcherSpec{Patterns: nil, MatchType: kernel.MatchType_Exact},

			AttributeKey: "attrKey",
			NamespaceKey: "nsKey",
			TextKey:      "textKey",

			SOAPNamespacePrefix: "soap",
		},
	}
}

func (*API) Create(a api.API[*api.Request, *api.Response], msg protoreflect.ProtoMessage) (any, error) {
	c := msg.(*v1.SOAPRESTMiddleware)

	eh, err := utilhttp.ErrorHandler(a, c.Spec.ErrorHandler)
	if err != nil {
		return nil, core.ErrCoreGenCreateObject.WithStack(err, map[string]any{"kind": kind})
	}

	m, err := txtutil.NewStringMatcher(txtutil.MatchTypes[c.Spec.Matcher.MatchType], c.Spec.Matcher.Patterns...)
	if err != nil {
		return nil, core.ErrCoreGenCreateObject.WithStack(err, map[string]any{"kind": kind})
	}

	return &soapREST{
		eh:    eh,
		paths: m,

		attributeKey: c.Spec.AttributeKey,
		namespaceKey: c.Spec.NamespaceKey,
		textKey:      c.Spec.TextKey,

		soapNamespacePrefix: c.Spec.SOAPNamespacePrefix,

		extractStringElement:  c.Spec.ExtractStringElement,
		extractBooleanElement: c.Spec.ExtractBooleanElement,
		extractIntegerElement: c.Spec.ExtractIntegerElement,
		extractFloatElement:   c.Spec.ExtractFloatElement,
	}, nil
}
