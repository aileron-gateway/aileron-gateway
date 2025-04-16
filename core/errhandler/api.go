// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package errhandler

import (
	v1 "github.com/aileron-gateway/aileron-gateway/apis/core/v1"
	"github.com/aileron-gateway/aileron-gateway/apis/kernel"
	"github.com/aileron-gateway/aileron-gateway/core"
	"github.com/aileron-gateway/aileron-gateway/kernel/api"
	"github.com/aileron-gateway/aileron-gateway/kernel/log"
	utilhttp "github.com/aileron-gateway/aileron-gateway/util/http"
	"google.golang.org/protobuf/reflect/protoreflect"
)

const (
	apiVersion = "core/v1"
	kind       = "ErrorHandler"
	Key        = apiVersion + "/" + kind
)

var Resource api.Resource = &API{
	BaseResource: &api.BaseResource{
		DefaultProto: &v1.ErrorHandler{
			APIVersion: apiVersion,
			Kind:       kind,
			Metadata: &kernel.Metadata{
				Namespace: "default",
				Name:      "default",
			},
			Spec: &v1.ErrorHandlerSpec{},
		},
	},
}

type API struct {
	*api.BaseResource
}

func (*API) Create(a api.API[*api.Request, *api.Response], msg protoreflect.ProtoMessage) (any, error) {
	c := msg.(*v1.ErrorHandler)

	messages := make([]*utilhttp.ErrorMessage, 0, len(c.Spec.ErrorMessages))
	for _, em := range c.Spec.ErrorMessages {
		if len(em.MIMEContents) == 0 {
			continue
		}
		msg, err := utilhttp.NewErrorMessage(em)
		if err != nil {
			return nil, core.ErrCoreGenCreateObject.WithStack(err, map[string]any{"kind": kind})
		}
		messages = append(messages, msg)
	}

	return &utilhttp.DefaultErrorHandler{
		LG:          log.DefaultOr(c.Metadata.Logger),
		StackAlways: c.Spec.StackAlways,
		Msgs:        messages,
	}, nil
}
