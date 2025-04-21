// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package httplogger

import (
	"time"

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
	kind       = "HTTPLogger"
	Key        = apiVersion + "/" + kind
)

var Resource api.Resource = &API{
	BaseResource: &api.BaseResource{
		DefaultProto: &v1.HTTPLogger{
			APIVersion: apiVersion,
			Kind:       kind,
			Metadata: &kernel.Metadata{
				Namespace: "default",
				Name:      "default",
			},
			Spec: &v1.HTTPLoggerSpec{
				Journal:    false,
				Timezone:   "Local",
				TimeFormat: "2006-01-02 15:04:05.000",
				Request: &v1.LoggingSpec{
					MaxContentLength: 1 << 12, // 4,096 = 1MB
				},
				Response: &v1.LoggingSpec{
					MaxContentLength: 1 << 12, // 4,096 = 1MB
				},
			},
		},
	},
}

type API struct {
	*api.BaseResource
}

func (*API) Mutate(msg protoreflect.ProtoMessage) protoreflect.ProtoMessage {
	c := msg.(*v1.HTTPLogger)

	// Set default journal MIMEs.
	// All media types are listed here.
	// https://www.iana.org/assignments/media-types/media-types.xhtml
	if len(c.Spec.Request.MIMEs) == 0 {
		c.Spec.Request.MIMEs = []string{
			"application/json",
			"application/x-www-form-urlencoded",
			"application/xml",
			"application/soap+xml",
			"application/graphql+json",
			"text/plain",
			"text/html",
			"text/xml",
		}
	}
	if len(c.Spec.Response.MIMEs) == 0 {
		c.Spec.Response.MIMEs = []string{
			"application/json",
			"application/x-www-form-urlencoded",
			"application/xml",
			"application/soap+xml",
			"application/graphql+json",
			"text/plain",
			"text/html",
			"text/xml",
		}
	}

	return c
}

func (*API) Create(a api.API[*api.Request, *api.Response], msg protoreflect.ProtoMessage) (any, error) {
	c := msg.(*v1.HTTPLogger)

	lg := log.DefaultOr(c.Metadata.Logger)

	if c.Spec.Logger != nil {
		x, err := api.ReferTypedObject[log.Logger](a, c.Spec.Logger)
		if err != nil {
			return nil, core.ErrCoreGenCreateObject.WithStack(err, map[string]any{"kind": kind})
		}
		lg = x
	}

	eh, err := utilhttp.ErrorHandler(a, c.Spec.ErrorHandler)
	if err != nil {
		return nil, core.ErrCoreGenCreateObject.WithStack(err, map[string]any{"kind": kind})
	}

	zone, err := time.LoadLocation(c.Spec.Timezone)
	if err != nil {
		return nil, core.ErrCoreGenCreateObject.WithStack(err, map[string]any{"kind": kind})
	}

	reqBL, err := newBaseLogger(c.Spec.Request, lg)
	if err != nil {
		return nil, core.ErrCoreGenCreateObject.WithStack(err, map[string]any{"kind": kind})
	}
	resBL, err := newBaseLogger(c.Spec.Response, lg)
	if err != nil {
		return nil, core.ErrCoreGenCreateObject.WithStack(err, map[string]any{"kind": kind})
	}

	if c.Spec.Journal { // journal logging which output body.
		return &journalLogger{
			lg:      lg,
			eh:      eh,
			req:     reqBL,
			res:     resBL,
			zone:    zone,
			timeFmt: c.Spec.TimeFormat,
		}, nil
	}

	return &httpLogger{
		req:     reqBL,
		res:     resBL,
		zone:    zone,
		timeFmt: c.Spec.TimeFormat,
	}, nil
}
