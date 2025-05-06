// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package soaprest

import (
	"cmp"
	"encoding/json"
	"encoding/xml"

	v1 "github.com/aileron-gateway/aileron-gateway/apis/app/v1"
	"github.com/aileron-gateway/aileron-gateway/apis/kernel"
	"github.com/aileron-gateway/aileron-gateway/app/middleware/soaprest/zxml"
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

func (s *API) Default() protoreflect.ProtoMessage {
	return &v1.SOAPRESTMiddleware{
		APIVersion: apiVersion,
		Kind:       kind,
		Metadata: &kernel.Metadata{
			Namespace: "default",
			Name:      "default",
		},
		Spec: &v1.SOAPRESTMiddlewareSpec{
			Matcher: &kernel.MatcherSpec{Patterns: nil, MatchType: kernel.MatchType_Exact},
			Method: &v1.SOAPRESTMiddlewareSpec_SimpleMethodSpec{
				SimpleMethodSpec: &v1.SimpleMethodSpec{
					TextKey:         "$",
					AttrPrefix:      "@",
					NamespaceSep:    ":",
					TrimSpace:       false,
					PreferShort:     false,
					IgnoreUnusedKey: false,
				},
			},
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

	var cv *zxml.JSONConverter
	switch c.Spec.Method.(type) {
	case *v1.SOAPRESTMiddlewareSpec_SimpleMethodSpec:
		cv = newSimpleConverter(c.Spec.GetSimpleMethodSpec())
	case *v1.SOAPRESTMiddlewareSpec_RayfishMethodSpec:
		cv = newRayfishConverter(c.Spec.GetRayfishMethodSpec())
	case *v1.SOAPRESTMiddlewareSpec_BadgerfishMethodSpec:
		cv = newBadgerfishConverter(c.Spec.GetBadgerfishMethodSpec())
	}

	// Use json.Number instead of float64 for JSON numbers.
	cv.WithJSONDecoderOpts(func(d *json.Decoder) { d.UseNumber() })
	// Prevent escaping of HTML special characters.
	cv.WithJSONEncoderOpts(func(e *json.Encoder) { e.SetEscapeHTML(false) })

	return &soapREST{
		eh:        eh,
		paths:     m,
		converter: cv,
	}, nil
}

func newSimpleConverter(spec *v1.SimpleMethodSpec) *zxml.JSONConverter {
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

func newRayfishConverter(spec *v1.RayfishMethodSpec) *zxml.JSONConverter {
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

func newBadgerfishConverter(spec *v1.BadgerfishMethodSpec) *zxml.JSONConverter {
	badgerfish := &zxml.BadgerFish{
		TextKey:      cmp.Or(spec.TextKey, "$"),
		AttrPrefix:   cmp.Or(spec.AttrPrefix, "@"),
		NamespaceSep: cmp.Or(spec.NamespaceSep, ":"),
		TrimSpace:    spec.TrimSpace,
	}
	badgerfish.WithEmptyValue(string(""))

	return &zxml.JSONConverter{
		EncodeDecoder: badgerfish,
		Header:        xml.Header,
	}
}
