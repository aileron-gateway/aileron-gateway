package oauth

import (
	"cmp"
	"net/http"
	"os"
	"strconv"
	"time"

	v1 "github.com/aileron-gateway/aileron-gateway/apis/app/v1"
	"github.com/aileron-gateway/aileron-gateway/apis/kernel"
	"github.com/aileron-gateway/aileron-gateway/app"
	"github.com/aileron-gateway/aileron-gateway/core"
	"github.com/aileron-gateway/aileron-gateway/kernel/api"
	"github.com/aileron-gateway/aileron-gateway/kernel/log"
	utilhttp "github.com/aileron-gateway/aileron-gateway/util/http"
	"github.com/aileron-gateway/aileron-gateway/util/security"
	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

const (
	apiVersion = "app/v1"
	kind       = "OAuthAuthenticationHandler"
	Key        = apiVersion + "/" + kind
)

var (
	skipATValidation  = false
	skipIDTValidation = false
)

var Resource api.Resource = &API{
	BaseResource: &api.BaseResource{
		DefaultProto: &v1.OAuthAuthenticationHandler{
			APIVersion: apiVersion,
			Kind:       kind,
			Metadata: &kernel.Metadata{
				Namespace: "default",
				Name:      "default",
			},
			Spec: &v1.OAuthAuthenticationHandlerSpec{},
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
	c := msg.(*v1.OAuthAuthenticationHandler)

	for name, spec := range c.Spec.Contexts {
		baseSpec := &v1.Context{
			Name: "default", // "default" is the special name that will be used by default.
			Provider: &v1.OAuthProvider{
				Endpoints: &v1.ProviderEndpoints{},
			},
			Client:            &v1.OAuthClient{},
			TokenRedeemer:     &v1.ClientRequester{},
			TokenIntrospector: &v1.ClientRequester{},
			JWTHandler:        &v1.JWTHandlerSpec{},
			ClaimsKey:         "AuthnClaims",
			ATProxyHeader:     "", // HTTP header name when proxy Access Token to upstream services.
			IDTProxyHeader:    "", // HTTP header name when proxy ID Token to upstream services.
			ATValidation: &v1.TokenValidation{
				Leeway: 5,
			},
			IDTValidation: &v1.TokenValidation{
				Leeway: 5,
			},
		}
		proto.Merge(baseSpec, spec)
		baseSpec.ATValidation.Iss = cmp.Or(baseSpec.ATValidation.Iss, baseSpec.Provider.Issuer)
		baseSpec.ATValidation.Aud = cmp.Or(baseSpec.ATValidation.Aud, baseSpec.Client.Audience, baseSpec.Client.ID)
		baseSpec.IDTValidation.Iss = cmp.Or(baseSpec.IDTValidation.Iss, baseSpec.Provider.Issuer)
		baseSpec.IDTValidation.Aud = cmp.Or(baseSpec.IDTValidation.Aud, baseSpec.Client.Audience, baseSpec.Client.ID)
		c.Spec.Contexts[name] = baseSpec
	}

	return c
}

func (*API) Create(a api.API[*api.Request, *api.Response], msg protoreflect.ProtoMessage) (any, error) {
	var err error

	skipATValidation, err = getEnvAsBool("AILERON_SKIP_AT_VALIDATION", false)
	if err != nil {
		return nil, core.ErrCoreGenCreateObject.WithStack(err, map[string]any{"kind": kind})
	}

	skipIDTValidation, err = getEnvAsBool("AILERON_SKIP_IDT_VALIDATION", false)
	if err != nil {
		return nil, core.ErrCoreGenCreateObject.WithStack(err, map[string]any{"kind": kind})
	}

	c := msg.(*v1.OAuthAuthenticationHandler)

	lg := log.DefaultOr(c.Metadata.Logger)

	authLg := lg
	if c.Spec.Logger != nil {
		alg, err := api.ReferTypedObject[log.Logger](a, c.Spec.Logger)
		if err != nil {
			return nil, core.ErrCoreGenCreateObject.WithStack(err, map[string]any{"kind": kind})
		}
		authLg = alg
	}

	eh, err := utilhttp.ErrorHandler(a, c.Spec.ErrorHandler)
	if err != nil {
		return nil, core.ErrCoreGenCreateObject.WithStack(err, map[string]any{"kind": kind})
	}

	ctxs := make(map[string]*oauthContext, len(c.Spec.Contexts))
	for _, spec := range c.Spec.Contexts {
		oc, err := newOAuthContext(a, spec, lg)
		if err != nil {
			return nil, core.ErrCoreGenCreateObject.WithStack(err, map[string]any{"kind": kind})
		}
		ctxs[spec.Name] = oc
	}

	bh := &baseHandler{
		lg:               authLg,
		eh:               eh,
		oauthCtxs:        ctxs,
		contextQueryKey:  c.Spec.ContextQueryKey,
		contextHeaderKey: c.Spec.ContextHeaderKey,
	}

	var handler app.AuthenticationHandler
	switch spec := c.Spec.Handlers.(type) {
	case *v1.OAuthAuthenticationHandlerSpec_AuthorizationCodeHandler:
		handler, err = newAuthorizationCodeHandler(bh, spec.AuthorizationCodeHandler)
	case *v1.OAuthAuthenticationHandlerSpec_ClientCredentialsHandler:
		handler = newClientCredentialsHandler(bh, spec.ClientCredentialsHandler)
	case *v1.OAuthAuthenticationHandlerSpec_ROPCHandler:
		handler = newROPCHandler(bh, spec.ROPCHandler)
	case *v1.OAuthAuthenticationHandlerSpec_ResourceServerHandler:
		handler = newResourceServerHandler(bh, spec.ResourceServerHandler)
	}
	if err != nil {
		return nil, core.ErrCoreGenCreateObject.WithStack(err, map[string]any{"kind": kind})
	}

	return handler, nil
}

func newOAuthContext(a api.API[*api.Request, *api.Response], spec *v1.Context, lg log.Logger) (*oauthContext, error) {
	client, err := newClient(spec.Client)
	if err != nil {
		return nil, core.ErrCoreGenCreateObject.WithStack(err, map[string]any{"kind": kind})
	}

	providerRT := http.DefaultTransport
	if spec.Provider.RoundTripper != nil {
		providerRT, err = api.ReferTypedObject[http.RoundTripper](a, spec.Provider.RoundTripper)
		if err != nil {
			return nil, core.ErrCoreGenCreateObject.WithStack(err, map[string]any{"kind": kind})
		}
	}

	provider, err := newProvider(spec.Provider, providerRT)
	if err != nil {
		return nil, core.ErrCoreGenCreateObject.WithStack(err, map[string]any{"kind": kind})
	}

	jh, err := security.NewJWTHandler(spec.JWTHandler, http.DefaultTransport)
	if err != nil {
		return nil, core.ErrCoreGenCreateObject.WithStack(err, map[string]any{"kind": kind})
	}

	redeemTokenRT := http.DefaultTransport
	if spec.TokenRedeemer.RoundTripper != nil {
		redeemTokenRT, err = api.ReferTypedObject[http.RoundTripper](a, spec.TokenRedeemer.RoundTripper)
		if err != nil {
			return nil, core.ErrCoreGenCreateObject.WithStack(err, map[string]any{"kind": kind})
		}
	}
	tokenRedeemer := &redeemTokenClient{
		requester: &clientRequester{
			lg:               lg,
			client:           client,
			rt:               redeemTokenRT,
			clientAuthMethod: clientAuthMethods[spec.TokenRedeemer.ClientAuthMethod],
		},
		lg:       lg,
		provider: provider,
	}

	introspectionRT := http.DefaultTransport
	if spec.TokenIntrospector.RoundTripper != nil {
		introspectionRT, err = api.ReferTypedObject[http.RoundTripper](a, spec.TokenIntrospector.RoundTripper)
		if err != nil {
			return nil, core.ErrCoreGenCreateObject.WithStack(err, map[string]any{"kind": kind})
		}
	}
	tokenIntrospector := &tokenIntrospectionClient{
		requester: &clientRequester{
			lg:               lg,
			client:           client,
			rt:               introspectionRT,
			clientAuthMethod: clientAuthMethods[spec.TokenIntrospector.ClientAuthMethod],
		},
		lg:       lg,
		provider: provider,
	}

	var atParseOpts []jwt.ParserOption
	atParseOpts = appendWithIssuer(atParseOpts, cmp.Or(spec.ATValidation.Iss, provider.issuer))
	atParseOpts = appendWithAudience(atParseOpts, cmp.Or(spec.ATValidation.Aud, client.audience, client.id))
	atParseOpts = appendWithExpirationRequired(atParseOpts, spec.ATValidation.ExpOptional)
	atParseOpts = appendWithIssuedAt(atParseOpts, spec.ATValidation.IatDisabled)
	atParseOpts = appendWithLeeway(atParseOpts, time.Duration(spec.ATValidation.Leeway)*time.Second)
	atParseOpts = appendWithValidMethods(atParseOpts, spec.ATValidation.ValidMethods)
	var idtParseOpts []jwt.ParserOption
	idtParseOpts = appendWithIssuer(idtParseOpts, cmp.Or(spec.IDTValidation.Iss, provider.issuer))
	idtParseOpts = appendWithAudience(idtParseOpts, cmp.Or(spec.IDTValidation.Aud, client.audience, client.id))
	idtParseOpts = appendWithExpirationRequired(idtParseOpts, spec.IDTValidation.ExpOptional)
	idtParseOpts = appendWithIssuedAt(idtParseOpts, spec.IDTValidation.IatDisabled)
	idtParseOpts = appendWithLeeway(idtParseOpts, time.Duration(spec.IDTValidation.Leeway)*time.Second)
	idtParseOpts = appendWithValidMethods(idtParseOpts, spec.IDTValidation.ValidMethods)

	return &oauthContext{
		tokenRedeemer:     tokenRedeemer,
		tokenIntrospector: tokenIntrospector,

		lg: lg,

		name: spec.Name,

		atParseOpts:      atParseOpts,
		idtParseOpts:     idtParseOpts,
		skipUnexpiredAT:  spec.ATValidation.SkipUnexpired,
		skipUnexpiredIDT: spec.IDTValidation.SkipUnexpired,

		provider:             provider,
		client:               client,
		jh:                   jh,
		introspectionEnabled: spec.EnableIntrospection,
		claimsKey:            spec.ClaimsKey,
		atProxyHeader:        spec.ATProxyHeader,
		idtProxyHeader:       spec.IDTProxyHeader,
	}, nil
}

func appendWithIssuer(opts []jwt.ParserOption, iss string) []jwt.ParserOption {
	if iss == "" || iss == "-" { // "-" disables validation.
		return opts
	}
	return append(opts, jwt.WithIssuer(iss))
}

func appendWithAudience(opts []jwt.ParserOption, aud string) []jwt.ParserOption {
	if aud == "" || aud == "-" { // "-" disables validation.
		return opts
	}
	return append(opts, jwt.WithAudience(aud))
}

func appendWithExpirationRequired(opts []jwt.ParserOption, optional bool) []jwt.ParserOption {
	// By default jwt.Validator validate exp if exists.
	// If optional is true, use the default behavior.
	// Otherwise, mandate to exp validation with WithExpirationRequired.
	// https://github.com/golang-jwt/jwt/blob/v5.2.1/validator.go#L105
	// https://pkg.go.dev/github.com/golang-jwt/jwt/v5#Validator.Validate
	if optional {
		return opts
	}
	return append(opts, jwt.WithExpirationRequired())
}

func appendWithIssuedAt(opts []jwt.ParserOption, disabled bool) []jwt.ParserOption {
	// By default jwt.Validator does not validate iat.
	// If disabled is true, use the default behavior.
	// Otherwise, validate the iat if exists with WithIssuedAt.
	// https://github.com/golang-jwt/jwt/blob/v5.2.1/validator.go#L116
	// https://pkg.go.dev/github.com/golang-jwt/jwt/v5#Validator.Validate
	if disabled {
		return opts
	}
	return append(opts, jwt.WithIssuedAt())
}

func appendWithLeeway(opts []jwt.ParserOption, d time.Duration) []jwt.ParserOption {
	if d > 0 {
		return opts
	}
	return append(opts, jwt.WithLeeway(d))
}

func appendWithValidMethods(opts []jwt.ParserOption, methods []string) []jwt.ParserOption {
	if len(methods) == 0 {
		return opts
	}
	return append(opts, jwt.WithValidMethods(methods))
}

func getEnvAsBool(envName string, defaultValue bool) (bool, error) {
	value := os.Getenv(envName)
	if value == "" {
		return defaultValue, nil
	}
	parsedValue, err := strconv.ParseBool(value)
	if err != nil {
		return defaultValue, err
	}
	return parsedValue, nil
}
