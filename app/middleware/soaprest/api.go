// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package soaprest

import (
	"cmp"
	"encoding/json"
	"encoding/xml"

	v1 "github.com/aileron-gateway/aileron-gateway/apis/app/v1"
	"github.com/aileron-gateway/aileron-gateway/apis/kernel"
	"github.com/aileron-gateway/aileron-gateway/kernel/api"
	"github.com/aileron-projects/go/zencoding/zxml"
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

func (s *API) Default() protoreflect.ProtoMessage {
	return &v1.SOAPRESTMiddleware{
		APIVersion: apiVersion,
		Kind:       kind,
		Metadata: &kernel.Metadata{
			Namespace: "default",
			Name:      "default",
		},
		Spec: &v1.SOAPRESTMiddlewareSpec{},
	}
}

func (*API) Create(a api.API[*api.Request, *api.Response], msg protoreflect.ProtoMessage) (any, error) {
	c := msg.(*v1.SOAPRESTMiddleware)
	eh := utilhttp.GlobalErrorHandler(cmp.Or(c.Metadata.ErrorHandler, utilhttp.DefaultErrorHandlerName))

	var cv *zxml.JSONConverter
	switch t := c.Spec.Rules.(type) {
	case *v1.SOAPRESTMiddlewareSpec_Simple:
		cv = newSimpleConverter(t.Simple)
	case *v1.SOAPRESTMiddlewareSpec_Rayfish:
		cv = newRayfishConverter(t.Rayfish)
	case *v1.SOAPRESTMiddlewareSpec_Badgerfish:
		cv = newBadgerfishConverter(t.Badgerfish)
	default:
		cv = newSimpleConverter(&v1.SimpleSpec{})
	}

	cv.WithJSONDecoderOpts(func(d *json.Decoder) { d.UseNumber() })
	cv.WithJSONEncoderOpts(func(e *json.Encoder) { e.SetEscapeHTML(false) })
	cv.WithXMLEncoderOpts(func(e *xml.Encoder) { e.Indent("", "  ") })

	return &soapREST{
		eh:        eh,
		converter: cv,
	}, nil
}

func newSimpleConverter(spec *v1.SimpleSpec) *zxml.JSONConverter {
	simple := &zxml.Simple{
		TextKey:      cmp.Or(spec.TextKey, "$"),
		AttrPrefix:   cmp.Or(spec.AttrPrefix, "@"),
		NamespaceSep: cmp.Or(spec.NamespaceSep, ":"),
		TrimSpace:    spec.TrimSpace,
		PreferShort:  spec.PreferShort,
	}
	simple.WithEmptyValue(string(""))
	return &zxml.JSONConverter{
		EncodeDecoder: simple,
		Header:        xml.Header,
	}
}

func newRayfishConverter(spec *v1.RayfishSpec) *zxml.JSONConverter {
	rayfish := &zxml.RayFish{
		NameKey:      cmp.Or(spec.NameKey, "#name"),
		TextKey:      cmp.Or(spec.TextKey, "#text"),
		ChildrenKey:  cmp.Or(spec.ChildrenKey, "#children"),
		AttrPrefix:   cmp.Or(spec.AttrPrefix, "@"),
		NamespaceSep: cmp.Or(spec.NamespaceSep, ":"),
		TrimSpace:    spec.TrimSpace,
	}
	rayfish.WithEmptyValue(string(""))
	return &zxml.JSONConverter{
		EncodeDecoder: rayfish,
		Header:        xml.Header,
	}
}

func newBadgerfishConverter(spec *v1.BadgerfishSpec) *zxml.JSONConverter {
	badgerfish := &zxml.BadgerFish{
		TextKey:      cmp.Or(spec.TextKey, "$"),
		AttrPrefix:   cmp.Or(spec.AttrPrefix, "@"),
		NamespaceSep: cmp.Or(spec.NamespaceSep, ":"),
		TrimSpace:    spec.TrimSpace,
	}
	badgerfish.WithEmptyValue(make(map[string]any, 0))
	return &zxml.JSONConverter{
		EncodeDecoder: badgerfish,
		Header:        xml.Header,
	}
}
