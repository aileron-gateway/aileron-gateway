// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package cors

import (
	"strconv"
	"strings"

	v1 "github.com/aileron-gateway/aileron-gateway/apis/app/v1"
	corev1 "github.com/aileron-gateway/aileron-gateway/apis/core/v1"
	"github.com/aileron-gateway/aileron-gateway/apis/kernel"
	"github.com/aileron-gateway/aileron-gateway/core"
	"github.com/aileron-gateway/aileron-gateway/kernel/api"
	"github.com/aileron-gateway/aileron-gateway/kernel/log"
	utilhttp "github.com/aileron-gateway/aileron-gateway/util/http"
	"google.golang.org/protobuf/reflect/protoreflect"
)

const (
	apiVersion = "app/v1"
	kind       = "CORSMiddleware"
	Key        = apiVersion + "/" + kind
)

var Resource api.Resource = &API{
	BaseResource: &api.BaseResource{
		DefaultProto: &v1.CORSMiddleware{
			APIVersion: apiVersion,
			Kind:       kind,
			Metadata: &kernel.Metadata{
				Namespace: "default",
				Name:      "default",
			},
			Spec: &v1.CORSMiddlewareSpec{
				CORSPolicy: &v1.CORSPolicySpec{},
			},
		},
	},
}

type API struct {
	*api.BaseResource
}

// Mutate changes configured values.
// The values of the msg which is given as the argument is the merged message of default values and user defined values.
// Changes for the fields of msg in this function make the final values which will be the input for validate and create function.
// Default values for "repeated" or "oneof" fields can also be applied in this function if necessary.
// Please check msg!=nil and asserting the mgs does not panic even they won't from the view of overall architecture of the gateway.
func (*API) Mutate(msg protoreflect.ProtoMessage) protoreflect.ProtoMessage {
	c := msg.(*v1.CORSMiddleware)

	// Set default values for CORS policy.
	p := c.Spec.CORSPolicy
	if len(p.AllowedOrigins) == 0 {
		p.AllowedOrigins = []string{"*"}
	}
	if len(p.AllowedMethods) == 0 {
		// Use simple methods by default defined in the following document.
		// https://www.w3.org/TR/2020/SPSD-cors-20200602/#simple-method
		p.AllowedMethods = []corev1.HTTPMethod{corev1.HTTPMethod_POST, corev1.HTTPMethod_GET, corev1.HTTPMethod_OPTIONS}
	}
	if len(p.AllowedHeaders) == 0 {
		p.AllowedHeaders = []string{"Content-Type", "X-Requested-With"}
	}

	return c
}

func (*API) Create(a api.API[*api.Request, *api.Response], msg protoreflect.ProtoMessage) (any, error) {
	c := msg.(*v1.CORSMiddleware)

	// TODO: Output debug logs in the CORS middleware.
	_ = log.DefaultOr(c.Metadata.Logger)

	eh, err := utilhttp.ErrorHandler(a, c.Spec.ErrorHandler)
	if err != nil {
		return nil, core.ErrCoreGenCreateObject.WithStack(err, map[string]any{"kind": kind})
	}

	cp := c.Spec.CORSPolicy

	return &cors{
		eh: eh,
		policy: &corsPolicy{
			allowedOrigins:        cp.AllowedOrigins,
			allowedMethods:        utilhttp.Methods(cp.AllowedMethods),
			joinedAllowedMethods:  strings.Join(utilhttp.Methods(cp.AllowedMethods), ","),
			joinedAllowedHeaders:  strings.Join(cp.AllowedHeaders, ","),
			joinedExposedHeaders:  strings.Join(cp.ExposedHeaders, ","),
			allowCredentials:      cp.AllowCredentials,
			maxAge:                strconv.Itoa(int(cp.MaxAge)),
			embedderPolicy:        embedderPolicy(cp.CORSEmbedderPolicy),
			openerPolicy:          openerPolicy(cp.CORSOpenerPolicy),
			resourcePolicy:        resourcePolicy(cp.CORSResourcePolicy),
			allowPrivateNetwork:   cp.AllowPrivateNetwork,
			disableWildCardOrigin: cp.DisableWildCardOrigin,
		},
	}, nil
}

// embedderPolicy returns value for the "Cross-Origin-Embedder-Policy" header.
//   - https://docs.w3cub.com/http/headers/cross-origin-embedder-policy
//   - https://html.spec.whatwg.org/multipage/browsers.html#the-coep-headers
func embedderPolicy(p v1.CORSEmbedderPolicy) string {
	switch p {
	case v1.CORSEmbedderPolicy_EmbedderUnsafeNone:
		return "unsafe-none"
	case v1.CORSEmbedderPolicy_EmbedderRequireCorp:
		return "require-corp"
	case v1.CORSEmbedderPolicy_EmbedderCredentialless:
		return "credentialless"
	default:
		return "" // Response header won't be set.
	}
}

// openerPolicy returns value for the "Cross-Origin-Opener-Policy" header.
//   - https://docs.w3cub.com/http/headers/cross-origin-opener-policy
//   - https://html.spec.whatwg.org/multipage/browsers.html#cross-origin-opener-policies
func openerPolicy(p v1.CORSOpenerPolicy) string {
	switch p {
	case v1.CORSOpenerPolicy_OpenerUnsafeNone:
		return "unsafe-none"
	case v1.CORSOpenerPolicy_OpenerSameOriginAllowPopups:
		return "same-origin-allow-popups"
	case v1.CORSOpenerPolicy_OpenerSameOrigin:
		return "same-origin"
	default:
		return "" // Response header won't be set.
	}
}

// resourcePolicy returns value for the "Cross-Origin-Resource-Policy" header.
//   - https://docs.w3cub.com/http/headers/cross-origin-resource-policy
//   - https://fetch.spec.whatwg.org/#cross-origin-resource-policy-header
func resourcePolicy(p v1.CORSResourcePolicy) string {
	switch p {
	case v1.CORSResourcePolicy_ResourceSameSite:
		return "same-site"
	case v1.CORSResourcePolicy_ResourceSameOrigin:
		return "same-origin"
	case v1.CORSResourcePolicy_ResourceCrossOrigin:
		return "cross-origin"
	default:
		return "" // Response header won't be set.
	}
}
