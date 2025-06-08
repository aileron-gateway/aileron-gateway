// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package key

import (
	"context"
	"encoding/base64"
	"errors"
	"net/textproto"

	v1 "github.com/aileron-gateway/aileron-gateway/apis/app/v1"
	"github.com/aileron-gateway/aileron-gateway/apis/kernel"
	"github.com/aileron-gateway/aileron-gateway/core"
	"github.com/aileron-gateway/aileron-gateway/kernel/api"
	"github.com/aileron-gateway/aileron-gateway/kernel/encoder"
	"github.com/aileron-gateway/aileron-gateway/kernel/hash"
	"github.com/aileron-gateway/aileron-gateway/kernel/kvs"
	"github.com/aileron-gateway/aileron-gateway/kernel/log"
	utilhttp "github.com/aileron-gateway/aileron-gateway/util/http"
	"google.golang.org/protobuf/reflect/protoreflect"
)

const (
	apiVersion = "app/v1"
	kind       = "KeyAuthnMiddleware"
	Key        = apiVersion + "/" + kind
)

var Resource api.Resource = &API{
	BaseResource: &api.BaseResource{
		DefaultProto: &v1.KeyAuthnMiddleware{
			APIVersion: apiVersion,
			Kind:       kind,
			Metadata: &kernel.Metadata{
				Namespace: "default",
				Name:      "default",
			},
			Spec: &v1.KeyAuthnMiddlewareSpec{
				ClaimsKey:     "AuthnClaims",
				KeyHeaderName: "X-Api-Key",
				Providers: &v1.KeyAuthnMiddlewareSpec_EnvProvider{
					EnvProvider: &v1.KeyAuthnEnvProvider{
						KeyPrefix: "GATEWAY_APIKEY_",
					},
				},
			},
		},
	},
}

type API struct {
	*api.BaseResource
}

func (*API) Create(a api.API[*api.Request, *api.Response], msg protoreflect.ProtoMessage) (any, error) {
	c := msg.(*v1.KeyAuthnMiddleware)

	lg := log.DefaultOr(c.Metadata.Logger)

	authLg := lg
	if c.Spec.Logger != nil {
		alg, err := api.ReferTypedObject[log.Logger](a, c.Spec.Logger)
		if err != nil {
			return nil, core.ErrCoreGenCreateObject.WithStack(err, map[string]any{"kind": kind})
		}
		authLg = alg
	}
	_ = authLg // Will be used for auth logs.

	eh, err := utilhttp.ErrorHandler(a, c.Spec.ErrorHandler)
	if err != nil {
		return nil, core.ErrCoreGenCreateObject.WithStack(err, map[string]any{"kind": kind})
	}

	var encodeFunc encoder.EncodeToStringFunc
	var store kvs.Commander[string, credential]
	switch v := c.Spec.Providers.(type) {
	case *v1.KeyAuthnMiddlewareSpec_EnvProvider:
		store, encodeFunc, err = newEnvProvider(v.EnvProvider)
	case *v1.KeyAuthnMiddlewareSpec_FileProvider:
		store, encodeFunc, err = newFileProvider(v.FileProvider)
	default:
		s := &kvs.MapKVS[string, credential]{}
		s.Open(context.Background())
		store = s
	}
	if err != nil {
		return nil, core.ErrCoreGenCreateObject.WithStack(err, map[string]any{"kind": kind})
	}

	var hmacKey []byte
	if c.Spec.HMACSecret != "" {
		key, err := base64.StdEncoding.DecodeString(c.Spec.HMACSecret)
		if err != nil {
			return nil, core.ErrCoreGenCreateObject.WithStack(err, map[string]any{"kind": kind})
		}
		hmacKey = key
		if c.Spec.HashAlg == kernel.HashAlg_HashAlgUnknown {
			err := errors.New("hmac key provided but hash algorithm is not set")
			return nil, core.ErrCoreGenCreateObject.WithStack(err, map[string]any{"kind": kind})
		}
	}

	return &handler{
		eh: eh,

		claimsKey:     c.Spec.ClaimsKey,
		keep:          c.Spec.KeepCredentials,
		keyHeaderName: textproto.CanonicalMIMEHeaderKey(c.Spec.KeyHeaderName),

		store: store,

		encodeFunc: encodeFunc,

		hmacKey:  hmacKey, // If hmacKey is not empty, do hmac.
		hmacFunc: hash.HMACFromHashAlg(c.Spec.HashAlg),
		hashFunc: hash.FromHashAlg(c.Spec.HashAlg),
	}, nil
}
