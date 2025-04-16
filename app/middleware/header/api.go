// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package header

import (
	"net/textproto"

	v1 "github.com/aileron-gateway/aileron-gateway/apis/app/v1"
	"github.com/aileron-gateway/aileron-gateway/apis/kernel"
	"github.com/aileron-gateway/aileron-gateway/core"
	"github.com/aileron-gateway/aileron-gateway/kernel/api"
	"github.com/aileron-gateway/aileron-gateway/kernel/txtutil"
	httputil "github.com/aileron-gateway/aileron-gateway/util/http"
	"google.golang.org/protobuf/reflect/protoreflect"
)

const (
	apiVersion = "app/v1"
	kind       = "HeaderPolicyMiddleware"
	Key        = apiVersion + "/" + kind
)

var Resource api.Resource = &API{
	BaseResource: &api.BaseResource{
		DefaultProto: &v1.HeaderPolicyMiddleware{
			APIVersion: apiVersion,
			Kind:       kind,
			Metadata: &kernel.Metadata{
				Namespace: "default",
				Name:      "default",
			},
			Spec: &v1.HeaderPolicyMiddlewareSpec{},
		},
	},
}

type API struct {
	*api.BaseResource
}

func (*API) Create(a api.API[*api.Request, *api.Response], msg protoreflect.ProtoMessage) (any, error) {
	c := msg.(*v1.HeaderPolicyMiddleware)

	eh, err := httputil.ErrorHandler(a, c.Spec.ErrorHandler)
	if err != nil {
		return nil, core.ErrCoreGenCreateObject.WithStack(err, map[string]any{"kind": kind})
	}

	middleware := &headerPolicy{
		eh: eh,

		allowedMIMEs:     c.Spec.AllowMIMEs,
		maxContentLength: c.Spec.MaxContentLength,
	}

	if c.Spec.RequestPolicy != nil {
		p := c.Spec.RequestPolicy
		repls, err := newRewriters(p.Rewrites)
		if err != nil {
			return nil, core.ErrCoreGenCreateObject.WithStack(err, map[string]any{"kind": kind})
		}
		middleware.reqPolicy = &policy{
			allows:  canonicalSlice(p.Allows),
			removes: canonicalSlice(p.Removes),
			add:     canonicalMapKey(p.Add),
			set:     canonicalMapKey(p.Set),
			repls:   repls,
		}
	}
	if c.Spec.ResponsePolicy != nil {
		p := c.Spec.ResponsePolicy
		repls, err := newRewriters(p.Rewrites)
		if err != nil {
			return nil, core.ErrCoreGenCreateObject.WithStack(err, map[string]any{"kind": kind})
		}
		middleware.resPolicy = &policy{
			allows:  canonicalSlice(p.Allows),
			removes: canonicalSlice(p.Removes),
			add:     canonicalMapKey(p.Add),
			set:     canonicalMapKey(p.Set),
			repls:   repls,
		}
	}

	return middleware, nil
}

func newRewriters(specs []*v1.HeaderRewriteSpec) (map[string]txtutil.ReplaceFunc[string], error) {
	result := make(map[string]txtutil.ReplaceFunc[string], len(specs))
	for _, spec := range specs {
		if spec == nil || spec.Name == "" || spec.Replacer == nil {
			continue
		}
		replacer, err := txtutil.NewStringReplacer(spec.Replacer)
		if err != nil {
			return nil, err
		}
		result[textproto.CanonicalMIMEHeaderKey(spec.Name)] = replacer.Replace
	}
	return result, nil
}

func canonicalSlice(headers []string) []string {
	h := make([]string, len(headers))
	for i, key := range headers {
		h[i] = textproto.CanonicalMIMEHeaderKey(key)
	}
	return h
}

func canonicalMapKey(headers map[string]string) map[string]string {
	h := make(map[string]string, len(headers))
	for key, value := range headers {
		h[textproto.CanonicalMIMEHeaderKey(key)] = value
	}
	return h
}
