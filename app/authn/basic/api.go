// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package basic

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"

	v1 "github.com/aileron-gateway/aileron-gateway/apis/app/v1"
	"github.com/aileron-gateway/aileron-gateway/apis/kernel"
	"github.com/aileron-gateway/aileron-gateway/core"
	"github.com/aileron-gateway/aileron-gateway/internal/encrypt"
	"github.com/aileron-gateway/aileron-gateway/internal/kvs"
	"github.com/aileron-gateway/aileron-gateway/kernel/api"
	"github.com/aileron-gateway/aileron-gateway/kernel/log"
	utilhttp "github.com/aileron-gateway/aileron-gateway/util/http"
	"google.golang.org/protobuf/reflect/protoreflect"
)

const (
	apiVersion = "app/v1"
	kind       = "BasicAuthnMiddleware"
	Key        = apiVersion + "/" + kind
)

var Resource api.Resource = &API{
	BaseResource: &api.BaseResource{
		DefaultProto: &v1.BasicAuthnMiddleware{
			APIVersion: apiVersion,
			Kind:       kind,
			Metadata: &kernel.Metadata{
				Namespace: "default",
				Name:      "default",
			},
			Spec: &v1.BasicAuthnMiddlewareSpec{
				ClaimsKey: "AuthnClaims",
				Providers: &v1.BasicAuthnMiddlewareSpec_EnvProvider{
					EnvProvider: &v1.BasicAuthnEnvProvider{
						UsernamePrefix: "GATEWAY_BASIC_USERNAME_",
						PasswordPrefix: "GATEWAY_BASIC_PASSWORD_",
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
	c := msg.(*v1.BasicAuthnMiddleware)

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

	var store kvs.Commander[string, credential]
	switch v := c.Spec.Providers.(type) {
	case *v1.BasicAuthnMiddlewareSpec_EnvProvider:
		store, err = newEnvProvider(v.EnvProvider)
	case *v1.BasicAuthnMiddlewareSpec_FileProvider:
		store, err = newFileProvider(v.FileProvider)
	default:
		s := &kvs.MapKVS[string, credential]{}
		s.Open(context.Background())
		store = s
	}
	if err != nil {
		return nil, core.ErrCoreGenCreateObject.WithStack(err, map[string]any{"kind": kind})
	}

	var passwd []byte
	var decryptFunc encrypt.DecryptFunc
	if c.Spec.CryptSecret != "" {
		passwd, err = base64.StdEncoding.DecodeString(c.Spec.CryptSecret)
		if err != nil {
			return nil, core.ErrCoreGenCreateObject.WithStack(err, map[string]any{"kind": kind})
		}
		decryptFunc = encrypt.DecrypterFromType(c.Spec.CommonKeyCryptType)
		if decryptFunc == nil {
			err = errors.New("unsupported common key crypt type")
			return nil, core.ErrCoreGenCreateObject.WithStack(err, map[string]any{"kind": kind})
		}
	}

	compareFunc := equal
	if c.Spec.PasswordCrypt != nil {
		crypt, err := encrypt.NewPasswordCrypt(c.Spec.PasswordCrypt)
		if err != nil {
			return nil, core.ErrCoreGenCreateObject.WithStack(err, map[string]any{"kind": kind})
		}
		println(c.Spec.PasswordCrypt.String())
		compareFunc = crypt.Compare
	}

	return &handler{
		eh: eh,

		claimsKey: c.Spec.ClaimsKey,
		keep:      c.Spec.KeepCredentials,
		challenge: fmt.Sprintf(`Basic realm="%s" charset=UTF-8"`, c.Spec.Realm),

		store: store,

		passwd:      passwd,
		decryptFunc: decryptFunc,
		compareFunc: compareFunc,

		preferError: c.Spec.PreferError,
	}, nil
}

func equal(a, b []byte) error {
	if string(a) == string(b) {
		return nil
	}
	return ErrNotMatch
}

var ErrNotMatch = errors.New("password hash not matched")
