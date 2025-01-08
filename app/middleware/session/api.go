package session

import (
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"os"

	v1 "github.com/aileron-gateway/aileron-gateway/apis/app/v1"
	"github.com/aileron-gateway/aileron-gateway/apis/kernel"
	"github.com/aileron-gateway/aileron-gateway/app"
	"github.com/aileron-gateway/aileron-gateway/core"
	"github.com/aileron-gateway/aileron-gateway/kernel/api"
	"github.com/aileron-gateway/aileron-gateway/kernel/log"
	utilhttp "github.com/aileron-gateway/aileron-gateway/util/http"
	"github.com/aileron-gateway/aileron-gateway/util/security"
	"github.com/aileron-gateway/aileron-gateway/util/session"
	"google.golang.org/protobuf/reflect/protoreflect"
)

const (
	apiVersion = "app/v1"
	kind       = "SessionMiddleware"
	Key        = apiVersion + "/" + kind
)

var (
	host, _     = os.Hostname()
	hmacSecret  = sha512.Sum512([]byte(host)) // Create default secret.
	cryptSecret = sha256.Sum256([]byte(host)) // Create default secret.
)

var Resource api.Resource = &API{
	BaseResource: &api.BaseResource{
		DefaultProto: &v1.SessionMiddleware{
			APIVersion: apiVersion,
			Kind:       kind,
			Metadata: &kernel.Metadata{
				Namespace: "default",
				Name:      "default",
			},
			Spec: &v1.SessionMiddlewareSpec{
				CookieName: "_SESSION",
				SecureEncoder: &v1.SecureEncoderSpec{
					HashAlg:            kernel.HashAlg_SHA256,
					HMACSecret:         base64.StdEncoding.EncodeToString(hmacSecret[:]),
					CommonKeyCryptType: kernel.CommonKeyCryptType_AESGCM,
					CryptSecret:        base64.StdEncoding.EncodeToString(cryptSecret[:]),
				},
			},
		},
	},
}

type API struct {
	*api.BaseResource
}

func (*API) Create(a api.API[*api.Request, *api.Response], msg protoreflect.ProtoMessage) (any, error) {
	c := msg.(*v1.SessionMiddleware)

	// TODO: Output debug logs in the session middleware.
	_ = log.DefaultOr(c.Metadata.Logger)

	eh, err := utilhttp.ErrorHandler(a, c.Spec.ErrorHandler)
	if err != nil {
		return nil, core.ErrCoreGenCreateObject.WithStack(err, map[string]any{"kind": kind})
	}

	var kvs sessionKVS
	if c.Spec.Storage != nil {
		kvs, err = api.ReferTypedObject[sessionKVS](a, c.Spec.Storage)
		if err != nil {
			return nil, core.ErrCoreGenCreateObject.WithStack(err, map[string]any{"kind": kind})
		}
	}

	enc, err := security.NewSecureEncoder(c.Spec.SecureEncoder)
	if err != nil {
		return nil, core.ErrCoreGenCreateObject.WithStack(err, map[string]any{"kind": kind})
	}

	var tracer app.Tracer
	if c.Spec.Tracer != nil {
		t, err := api.ReferTypedObject[app.Tracer](a, c.Spec.Tracer)
		if err != nil {
			return nil, core.ErrCoreGenCreateObject.WithStack(err, map[string]any{"kind": kind})
		}
		tracer = t
	}

	var store session.Store
	if kvs != nil {
		// Use external session store like redis or memcached, etc..
		store = &kvsSessionStore{
			cookieName: c.Spec.CookieName,
			cc:         utilhttp.NewCookieCreator(c.Spec.Cookie),
			tracer:     tracer,
			prefix:     c.Spec.Prefix,
			store:      kvs,
			enc:        enc,
		}
	} else {
		// Use cookie as session store.
		store = &cookieSessionStore{
			cookiePrefix: c.Spec.CookieName,
			cc:           utilhttp.NewCookieCreator(c.Spec.Cookie),
			enc:          enc,
		}
	}

	return &sessioner{
		eh:    eh,
		store: store,
	}, nil
}
