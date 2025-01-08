package csrf

import (
	"crypto/sha512"
	"encoding/base64"
	"errors"
	"net/textproto"
	"os"
	"regexp"

	v1 "github.com/aileron-gateway/aileron-gateway/apis/app/v1"
	"github.com/aileron-gateway/aileron-gateway/apis/kernel"
	"github.com/aileron-gateway/aileron-gateway/core"
	"github.com/aileron-gateway/aileron-gateway/kernel/api"
	"github.com/aileron-gateway/aileron-gateway/kernel/hash"
	"github.com/aileron-gateway/aileron-gateway/kernel/log"
	"github.com/aileron-gateway/aileron-gateway/kernel/mac"
	utilhttp "github.com/aileron-gateway/aileron-gateway/util/http"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

const (
	apiVersion = "app/v1"
	kind       = "CSRFMiddleware"
	Key        = apiVersion + "/" + kind
)

var (
	host, _       = os.Hostname()
	defaultSecret = sha512.Sum512([]byte(host))
)

var Resource api.Resource = &API{
	BaseResource: &api.BaseResource{
		DefaultProto: &v1.CSRFMiddleware{
			APIVersion: apiVersion,
			Kind:       kind,
			Metadata: &kernel.Metadata{
				Namespace: "default",
				Name:      "default",
			},
			Spec: &v1.CSRFMiddlewareSpec{
				SeedSize: 20,
				Secret:   base64.StdEncoding.EncodeToString(defaultSecret[:]),
				HashAlg:  kernel.HashAlg_SHA256,
				CSRFPatterns: &v1.CSRFMiddlewareSpec_CustomRequestHeader{
					CustomRequestHeader: &v1.CustomRequestHeaderSpec{
						HeaderName: "X-Requested-With",
					},
				},
			},
		},
	},
}

type API struct {
	*api.BaseResource
}

func (*API) Mutate(msg protoreflect.ProtoMessage) protoreflect.ProtoMessage {
	c := msg.(*v1.CSRFMiddleware)

	if len(c.Spec.Patterns) == 0 {
		c.Spec.Patterns = append(c.Spec.Patterns, "/token")
	}

	switch c.Spec.CSRFPatterns.(type) {
	case *v1.CSRFMiddlewareSpec_DoubleSubmitCookies:
		baseSpec := &v1.DoubleSubmitCookiesSpec{
			CookieName:  "__csrfToken",
			TokenSource: v1.TokenSource_Header,
			SourceKey:   "__csrfToken",
		}
		proto.Merge(baseSpec, c.Spec.GetDoubleSubmitCookies())
		c.Spec.CSRFPatterns = &v1.CSRFMiddlewareSpec_DoubleSubmitCookies{
			DoubleSubmitCookies: baseSpec,
		}
	case *v1.CSRFMiddlewareSpec_SynchronizerToken:
		baseSpec := &v1.SynchronizerTokenSpec{
			TokenSource: v1.TokenSource_Header,
			SourceKey:   "__csrfToken",
		}
		proto.Merge(baseSpec, c.Spec.GetSynchronizerToken())
		c.Spec.CSRFPatterns = &v1.CSRFMiddlewareSpec_SynchronizerToken{
			SynchronizerToken: baseSpec,
		}
	}

	return c
}

func (*API) Create(a api.API[*api.Request, *api.Response], msg protoreflect.ProtoMessage) (any, error) {
	c := msg.(*v1.CSRFMiddleware)

	// TODO: Output debug logs in the CSRF middleware.
	_ = log.DefaultOr(c.Metadata.Logger)

	eh, err := utilhttp.ErrorHandler(a, c.Spec.ErrorHandler)
	if err != nil {
		return nil, core.ErrCoreGenCreateObject.WithStack(err, map[string]any{"kind": kind})
	}

	s, err := base64.StdEncoding.DecodeString(c.Spec.Secret)
	if err != nil {
		return nil, core.ErrCoreGenCreateObject.WithStack(err, map[string]any{"kind": kind})
	}
	mac := mac.FromHashAlg(c.Spec.HashAlg)
	if mac == nil {
		err := errors.New("invalid hash algorithm " + c.Spec.HashAlg.String())
		return nil, core.ErrCoreGenCreateObject.WithStack(err, map[string]any{"kind": kind})
	}
	hashSize := hash.HashSize[c.Spec.HashAlg]
	token := &csrfToken{
		secret:   s,
		seedSize: hash.HashSize[c.Spec.HashAlg],
		hashSize: hashSize,
		hmac:     mac,
	}

	var st strategy
	switch c.Spec.CSRFPatterns.(type) {
	case *v1.CSRFMiddlewareSpec_CustomRequestHeader:
		st, err = newCustomRequestHeaders(token, c.Spec.GetCustomRequestHeader())
		if err != nil {
			return nil, core.ErrCoreGenCreateObject.WithStack(err, map[string]any{"kind": kind})
		}
	case *v1.CSRFMiddlewareSpec_DoubleSubmitCookies:
		st = newDoubleSubmitCookies(token, c.Spec.GetDoubleSubmitCookies())
	case *v1.CSRFMiddlewareSpec_SynchronizerToken:
		st = newSynchronizerToken(token, c.Spec.GetSynchronizerToken())
	}

	return &csrf{
		HandlerBase: &utilhttp.HandlerBase{
			AcceptPatterns: c.Spec.Patterns,
			AcceptMethods:  utilhttp.Methods(c.Spec.Methods),
		},
		eh:              eh,
		proxyHeaderName: textproto.CanonicalMIMEHeaderKey(c.Spec.ProxyHeaderName),
		issueNew:        c.Spec.IssueNew,
		token:           token,
		st:              st,
	}, nil
}

func newCustomRequestHeaders(t *csrfToken, spec *v1.CustomRequestHeaderSpec) (strategy, error) {
	p, err := regexp.Compile(spec.AllowedPattern)
	if err != nil {
		return nil, err
	}
	if spec.AllowedPattern == "" {
		p = nil
	}
	return &customRequestHeaders{
		headerName: textproto.CanonicalMIMEHeaderKey(spec.HeaderName),
		pattern:    p,
		token:      t,
	}, nil
}

func newDoubleSubmitCookies(t *csrfToken, spec *v1.DoubleSubmitCookiesSpec) strategy {
	var ext extractor
	switch spec.TokenSource {
	case v1.TokenSource_Header:
		ext = &headerExtractor{
			headerName: textproto.CanonicalMIMEHeaderKey(spec.SourceKey),
		}
	case v1.TokenSource_Form:
		ext = &formExtractor{
			paramName: spec.SourceKey,
		}
	case v1.TokenSource_JSON:
		ext = &jsonExtractor{
			jsonPath: spec.SourceKey,
		}
	}

	return &doubleSubmitCookies{
		token:      t,
		ext:        ext,
		cookieName: spec.CookieName,
		cookie:     utilhttp.NewCookieCreator(spec.Cookie),
	}
}

func newSynchronizerToken(t *csrfToken, spec *v1.SynchronizerTokenSpec) strategy {
	var ext extractor
	switch spec.TokenSource {
	case v1.TokenSource_Header:
		ext = &headerExtractor{
			headerName: textproto.CanonicalMIMEHeaderKey(spec.SourceKey),
		}
	case v1.TokenSource_Form:
		ext = &formExtractor{
			paramName: spec.SourceKey,
		}
	case v1.TokenSource_JSON:
		ext = &jsonExtractor{
			jsonPath: spec.SourceKey,
		}
	}
	return &synchronizerToken{
		token: t,
		ext:   ext,
	}
}
