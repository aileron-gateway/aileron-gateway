// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package template

import (
	"cmp"
	"net/http"

	v1 "github.com/aileron-gateway/aileron-gateway/apis/core/v1"
	"github.com/aileron-gateway/aileron-gateway/apis/kernel"
	"github.com/aileron-gateway/aileron-gateway/core"
	"github.com/aileron-gateway/aileron-gateway/kernel/api"
	utilhttp "github.com/aileron-gateway/aileron-gateway/util/http"
	"google.golang.org/protobuf/reflect/protoreflect"
)

const (
	apiVersion = "core/v1"
	kind       = "TemplateHandler"
	Key        = apiVersion + "/" + kind
)

var Resource api.Resource = &API{
	BaseResource: &api.BaseResource{
		DefaultProto: &v1.TemplateHandler{
			APIVersion: apiVersion,
			Kind:       kind,
			Metadata: &kernel.Metadata{
				Namespace: "default",
				Name:      "default",
			},
			Spec: &v1.TemplateHandlerSpec{},
		},
	},
}

type API struct {
	*api.BaseResource
}

func (*API) Mutate(msg protoreflect.ProtoMessage) protoreflect.ProtoMessage {
	c := msg.(*v1.TemplateHandler)
	for _, mc := range c.Spec.MIMEContents {
		if mc.MIMEType == "" {
			mc.MIMEType = "text/plain"
		}
		if mc.StatusCode <= 0 {
			mc.StatusCode = http.StatusOK
		}
	}
	return c
}

func (*API) Create(a api.API[*api.Request, *api.Response], msg protoreflect.ProtoMessage) (any, error) {
	c := msg.(*v1.TemplateHandler)

	contents := make([]*utilhttp.MIMEContent, 0, len(c.Spec.MIMEContents))
	for _, cs := range c.Spec.MIMEContents {
		c, err := utilhttp.NewMIMEContent(cs)
		if err != nil {
			return nil, core.ErrCoreGenCreateObject.WithStack(err, map[string]any{"kind": kind})
		}
		contents = append(contents, c)
	}

	return &templateHandler{
		HandlerBase: &utilhttp.HandlerBase{
			AcceptPatterns: c.Spec.Patterns,
			AcceptMethods:  utilhttp.Methods(c.Spec.Methods),
		},
		eh:       utilhttp.GlobalErrorHandler(cmp.Or(c.Metadata.ErrorHandler, utilhttp.DefaultErrorHandlerName)),
		contents: contents,
	}, nil
}
