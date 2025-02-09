package oauth

import (
	"context"
	"net/http"
	"regexp"
	"testing"
	"time"

	"github.com/aileron-gateway/aileron-gateway/app"
	"github.com/aileron-gateway/aileron-gateway/core"
	"github.com/aileron-gateway/aileron-gateway/kernel/log"
	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
	utilhttp "github.com/aileron-gateway/aileron-gateway/util/http"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/go-cmp/cmp/cmpopts"
)

// newMapWithKeys returns a new map with
// values obtained by the given keys from the given source map.
// This function ignore keys if they cannot be found in the source map.
// This function returns an zero-length map, or non-nil, even no keys
// found in the source map.
func newMapWithKeys(keys []string, src map[string]any) map[string]any {
	result := make(map[string]any, len(keys))
	for _, k := range keys {
		if v, ok := src[k]; ok {
			result[k] = v
		}
	}
	return result
}

// testSimpleJWT is a simple JWT for test.
// JWT is signed by HS256 with "password".
//
//	{
//		"alg": "HS256",
//		"typ": "JWT",
//		"kid": "test"
//	}
//	{
//		"iss": "test-iss",
//		"aud": "test-aud",
//		"azp": "test-id"
//	}
const testSimpleJWT = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6InRlc3QifQ.eyJpc3MiOiJ0ZXN0LWlzcyIsImF1ZCI6InRlc3QtYXVkIiwiYXpwIjoidGVzdC1pZCJ9.pUrve9DuJIsKngajNrMAJy02VL46_I4pQ1soAAjtoCQ"

// testWringSigJWT is the same as testSimpleJWT
// except for that this JWT is signed with "wrong-password"
const testWringSigJWT = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6InRlc3QifQ.eyJpc3MiOiJ0ZXN0LWlzcyIsImF1ZCI6InRlc3QtYXVkIiwiYXpwIjoidGVzdC1pZCJ9.EB-mX-kz4eQfakeqOacglDg1fsy63rXMFrsk5L9Zn-E"

// testSimpleJWTCnf is the almost same as testSimpleJWT
// except for that this JWT has cnf claim.
// The value is x5tS256([]byte("client certificate")).
//
//	"cnf":{
//	    "x5t#S256": "hcP0Kggtjr3QGCLtW9ZJHba2IuxGMEyCn4H59Qz7A9g"
//	 }
const testSimpleJWTCnf = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6InRlc3QifQ.eyJpc3MiOiJ0ZXN0LWlzcyIsImF1ZCI6InRlc3QtYXVkIiwiYXpwIjoidGVzdC1pZCIsImNuZiI6eyJ4NXQjUzI1NiI6ImhjUDBLZ2d0anIzUUdDTHRXOVpKSGJhMkl1eEdNRXlDbjRINTlRejdBOWcifX0.kz-zILroUwvmHLIBkru31938jHxmzijYE-NKkReZ5Es"

// testBasicOauthContext is the oauth context that can be used
// for verifying the testSimpleJWT with local validation.
var testDefaultOauthContext = &oauthContext{
	name: "default",
	lg:   log.GlobalLogger(log.DefaultLoggerName),
	provider: &provider{
		issuer: "test-iss",
	},
	client: &client{
		id:       "test-id",
		audience: "test-aud",
	},
	jh: newJWTHandler(),
}

func TestValidOauthClaims(t *testing.T) {
	type condition struct {
		oc     *oauthContext
		tokens *OAuthTokens
	}

	type action struct {
		tokens     *OAuthTokens
		errStatus  int
		err        any
		errPattern *regexp.Regexp
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	cndIDT := tb.Condition("IDT", "input non empty ID token")
	cndAT := tb.Condition("AT signature", "input non empty access token")
	cndTokenInvalid := tb.Condition("validate iss", "input token is invalid")
	actCheckReAuth := tb.Action("re-auth", "check that re-authentication required")
	actCheckToken := tb.Action("error", "check that a non-nil token is returned")
	actCheckError := tb.Action("error", "check that an error was returned")
	actCheckNoError := tb.Action("no error", "check that no error was returned")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"nil token",
			[]string{},
			[]string{actCheckError, actCheckReAuth},
			&condition{
				oc: &oauthContext{
					atParseOpts:  []jwt.ParserOption{jwt.WithIssuer("test-iss"), jwt.WithAudience("test-aud")},
					idtParseOpts: []jwt.ParserOption{jwt.WithIssuer("test-iss"), jwt.WithAudience("test-aud")},
				},
				tokens: nil,
			},
			&action{
				tokens:    nil,
				errStatus: http.StatusUnauthorized,
				err:       reAuthenticationRequired,
			},
		),
		gen(
			"ID and Access token empty",
			[]string{},
			[]string{actCheckError, actCheckReAuth},
			&condition{
				oc: &oauthContext{
					atParseOpts:  []jwt.ParserOption{jwt.WithIssuer("test-iss"), jwt.WithAudience("test-aud")},
					idtParseOpts: []jwt.ParserOption{jwt.WithIssuer("test-iss"), jwt.WithAudience("test-aud")},
				},
				tokens: &OAuthTokens{
					IDT: "",
					AT:  "",
				},
			},
			&action{
				tokens:    &OAuthTokens{},
				errStatus: http.StatusUnauthorized,
				err:       reAuthenticationRequired,
			},
		),
		gen(
			"valid AT",
			[]string{cndAT},
			[]string{actCheckNoError, actCheckToken},
			&condition{
				oc: &oauthContext{
					lg: log.GlobalLogger(log.DefaultLoggerName),
					provider: &provider{
						issuer: "test-iss",
					},
					client: &client{
						id:       "test-id",
						audience: "test-aud",
					},
					jh:           newJWTHandler(),
					atParseOpts:  []jwt.ParserOption{jwt.WithIssuer("test-iss"), jwt.WithAudience("test-aud")},
					idtParseOpts: []jwt.ParserOption{jwt.WithIssuer("test-iss"), jwt.WithAudience("test-aud")},
				},
				tokens: &OAuthTokens{
					AT: testSimpleJWT,
				},
			},
			&action{
				tokens: &OAuthTokens{
					AT:       testSimpleJWT,
					ATClaims: map[string]any{"iss": "test-iss", "aud": "test-aud", "azp": "test-id", "active": true},
				},
				err: nil,
			},
		),
		gen(
			"validate AT failed with introspection",
			[]string{cndAT, cndTokenInvalid},
			[]string{actCheckError},
			&condition{
				oc: &oauthContext{
					lg: log.GlobalLogger(log.DefaultLoggerName),
					provider: &provider{
						issuer: "test-iss-xxx",
					},
					client: &client{
						id:       "test-id",
						audience: "test-aud",
					},
					jh:                   newJWTHandler(),
					introspectionEnabled: true,
					tokenIntrospector: &testIntrospector{
						status: http.StatusUnauthorized,
						claims: nil,
						err:    utilhttp.NewHTTPError(app.ErrAppAuthnIntrospection.WithoutStack(nil, nil), http.StatusUnauthorized),
					},
				},
				tokens: &OAuthTokens{
					AT: testSimpleJWT,
				},
			},
			&action{
				tokens: &OAuthTokens{
					AT: testSimpleJWT,
				},
				errStatus:  http.StatusUnauthorized,
				err:        app.ErrAppAuthnIntrospection,
				errPattern: regexp.MustCompile(core.ErrPrefix + `failed to token introspection.`),
			},
		),
		gen(
			"validate AT failed with local validation",
			[]string{cndAT, cndTokenInvalid},
			[]string{actCheckError, actCheckReAuth},
			&condition{
				oc: &oauthContext{
					lg: log.GlobalLogger(log.DefaultLoggerName),
					provider: &provider{
						issuer: "test-iss-xxx",
					},
					client: &client{
						id:       "test-id",
						audience: "test-aud",
					},
					jh:           newJWTHandler(),
					atParseOpts:  []jwt.ParserOption{jwt.WithIssuer("test-iss-xxx"), jwt.WithAudience("test-aud")},
					idtParseOpts: []jwt.ParserOption{jwt.WithIssuer("test-iss-xxx"), jwt.WithAudience("test-aud")},
				},
				tokens: &OAuthTokens{
					AT: testSimpleJWT,
				},
			},
			&action{
				tokens: &OAuthTokens{
					AT: testSimpleJWT,
				},
				errStatus: http.StatusUnauthorized,
				err:       reAuthenticationRequired,
			},
		),
		gen(
			"validate AT failed and Refresh failed with 400",
			[]string{cndAT, cndTokenInvalid},
			[]string{actCheckError, actCheckReAuth},
			&condition{
				oc: &oauthContext{
					lg: log.GlobalLogger(log.DefaultLoggerName),
					provider: &provider{
						issuer: "test-iss-xxx",
					},
					client: &client{
						id:       "test-id",
						audience: "test-aud",
					},
					jh: newJWTHandler(),
					tokenRedeemer: &testRedeemer{
						status: http.StatusUnauthorized,
						err:    utilhttp.NewHTTPError(app.ErrAppAuthnRedeemToken.WithoutStack(nil, nil), http.StatusUnauthorized),
					},
					atParseOpts:  []jwt.ParserOption{jwt.WithIssuer("test-iss-xxx"), jwt.WithAudience("test-aud")},
					idtParseOpts: []jwt.ParserOption{jwt.WithIssuer("test-iss-xxx"), jwt.WithAudience("test-aud")},
				},
				tokens: &OAuthTokens{
					AT: testSimpleJWT,
					RT: testSimpleJWT,
				},
			},
			&action{
				tokens: &OAuthTokens{
					AT: testSimpleJWT,
					RT: testSimpleJWT,
				},
				errStatus: http.StatusUnauthorized,
				err:       reAuthenticationRequired,
			},
		),
		gen(
			"validate AT failed and Refresh failed with 500",
			[]string{cndAT},
			[]string{actCheckError},
			&condition{
				oc: &oauthContext{
					lg: log.GlobalLogger(log.DefaultLoggerName),
					provider: &provider{
						issuer: "test-iss-xxx",
					},
					client: &client{
						id:       "test-id",
						audience: "test-aud",
					},
					jh: newJWTHandler(),
					tokenRedeemer: &testRedeemer{
						status: http.StatusInternalServerError,
						err:    utilhttp.NewHTTPError(app.ErrAppAuthnRedeemToken.WithoutStack(nil, nil), http.StatusInternalServerError),
					},
					atParseOpts:  []jwt.ParserOption{jwt.WithIssuer("test-iss-xxx"), jwt.WithAudience("test-aud")},
					idtParseOpts: []jwt.ParserOption{jwt.WithIssuer("test-iss-xxx"), jwt.WithAudience("test-aud")},
				},
				tokens: &OAuthTokens{
					AT: testSimpleJWT,
					RT: testSimpleJWT,
				},
			},
			&action{
				tokens: &OAuthTokens{
					AT: testSimpleJWT,
					RT: testSimpleJWT,
				},
				errStatus:  http.StatusInternalServerError,
				err:        app.ErrAppAuthnRedeemToken,
				errPattern: regexp.MustCompile(core.ErrPrefix + `failed to redeem token.`),
			},
		),
		gen(
			"validate AT failed and Refreshed AT invalid",
			[]string{cndAT, cndTokenInvalid},
			[]string{actCheckError, actCheckReAuth},
			&condition{
				oc: &oauthContext{
					lg: log.GlobalLogger(log.DefaultLoggerName),
					provider: &provider{
						issuer: "test-iss-xxx",
					},
					client: &client{
						id:       "test-id",
						audience: "test-aud",
					},
					jh: newJWTHandler(),
					tokenRedeemer: &testRedeemer{
						status: http.StatusOK,
						resp: &TokenResponse{
							AccessToken: testWringSigJWT,
						},
					},
					atParseOpts:  []jwt.ParserOption{jwt.WithIssuer("test-iss-xxx"), jwt.WithAudience("test-aud")},
					idtParseOpts: []jwt.ParserOption{jwt.WithIssuer("test-iss-xxx"), jwt.WithAudience("test-aud")},
				},
				tokens: &OAuthTokens{
					AT: testSimpleJWT,
					RT: testSimpleJWT,
				},
			},
			&action{
				tokens: &OAuthTokens{
					AT: testWringSigJWT,
					RT: "",
				},
				errStatus: http.StatusUnauthorized,
				err:       reAuthenticationRequired,
			},
		),
		gen(
			"valid IDT",
			[]string{cndIDT},
			[]string{actCheckNoError, actCheckToken},
			&condition{
				oc: &oauthContext{
					lg: log.GlobalLogger(log.DefaultLoggerName),
					provider: &provider{
						issuer: "test-iss",
					},
					client: &client{
						id:       "test-id",
						audience: "test-aud",
					},
					jh:           newJWTHandler(),
					atParseOpts:  []jwt.ParserOption{jwt.WithIssuer("test-iss"), jwt.WithAudience("test-aud")},
					idtParseOpts: []jwt.ParserOption{jwt.WithIssuer("test-iss"), jwt.WithAudience("test-aud")},
				},
				tokens: &OAuthTokens{
					IDT: testSimpleJWT,
				},
			},
			&action{
				tokens: &OAuthTokens{
					IDT:       testSimpleJWT,
					IDTClaims: map[string]any{"iss": "test-iss", "aud": "test-aud", "azp": "test-id"},
				},
				err: nil,
			},
		),
		gen(
			"validate IDT failed",
			[]string{cndIDT},
			[]string{actCheckError},
			&condition{
				oc: &oauthContext{
					lg: log.GlobalLogger(log.DefaultLoggerName),
					provider: &provider{
						issuer: "test-iss-xxx",
					},
					client: &client{
						id:       "test-id",
						audience: "test-aud",
					},
					jh:           newJWTHandler(),
					atParseOpts:  []jwt.ParserOption{jwt.WithIssuer("test-iss-xxx"), jwt.WithAudience("test-aud")},
					idtParseOpts: []jwt.ParserOption{jwt.WithIssuer("test-iss-xxx"), jwt.WithAudience("test-aud")},
				},
				tokens: &OAuthTokens{
					IDT: testSimpleJWT,
				},
			},
			&action{
				tokens: &OAuthTokens{
					IDT: testSimpleJWT,
				},
				errStatus:  http.StatusUnauthorized,
				err:        app.ErrAppAuthnInvalidToken,
				errPattern: regexp.MustCompile(core.ErrPrefix + `invalid id token. token validation failed.`),
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			err := tt.C().oc.validOauthClaims(context.Background(), tt.C().tokens)

			if tt.C().tokens != nil {
				testutil.Diff(t, tt.A().tokens, tt.C().tokens, cmpopts.IgnoreUnexported(OAuthTokens{}))
			}

			if tt.A().err == nil {
				testutil.Diff(t, nil, err)
				return
			}

			testutil.Diff(t, tt.A().errStatus, err.StatusCode())
			if tt.A().err == reAuthenticationRequired {
				testutil.Diff(t, reAuthenticationRequired.Error(), err.Error())
			} else {
				e := err.(*utilhttp.HTTPError)
				testutil.DiffError(t, tt.A().err, tt.A().errPattern, e.Unwrap())
			}
		})
	}
}

func TestOauthContext_validateAT(t *testing.T) {
	type condition struct {
		oc  *oauthContext
		ctx context.Context
		at  string
	}

	type action struct {
		claims     jwt.MapClaims
		errStatus  int
		err        error
		errPattern *regexp.Regexp
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	cndIntrospection := tb.Condition("introspection", "validate access token with introspection")
	cndValidateSign := tb.Condition("validate signature", "validate jwt signature")
	cndValidateIss := tb.Condition("validate iss", "validate iss claim in the JWT")
	cndValidateAud := tb.Condition("validate aud", "validate aud claim in the JWT")
	actCheckError := tb.Action("error", "check that an error was returned")
	actCheckNoError := tb.Action("no error", "check that no error was returned")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"introspect validation / valid minimal",
			[]string{cndIntrospection},
			[]string{actCheckNoError},
			&condition{
				oc: &oauthContext{
					lg: log.GlobalLogger(log.DefaultLoggerName),
					provider: &provider{
						issuer: "test-iss",
					},
					client: &client{
						id:       "test-id",
						audience: "test-aud",
					},
					jh:                   newJWTHandler(),
					introspectionEnabled: true,
					tokenIntrospector: &testIntrospector{
						status: http.StatusOK,
						claims: jwt.MapClaims{
							"active": true,
							"iss":    "test-iss",
							"aud":    "test-aud",
							"azp":    "test-id",
						},
					},
				},
				ctx: context.Background(),
				at:  "test-access-token",
			},
			&action{
				claims: jwt.MapClaims{
					"active": true,
					"iss":    "test-iss",
					"aud":    "test-aud",
					"azp":    "test-id",
				},
				err: nil,
			},
		),
		gen(
			"introspect validation / failed",
			[]string{cndIntrospection},
			[]string{actCheckError},
			&condition{
				oc: &oauthContext{
					lg: log.GlobalLogger(log.DefaultLoggerName),
					provider: &provider{
						issuer: "test-iss",
					},
					client: &client{
						id:       "test-id",
						audience: "test-aud",
					},
					jh:                   newJWTHandler(),
					introspectionEnabled: true,
					tokenIntrospector: &testIntrospector{
						status: http.StatusUnauthorized,
						claims: nil,
						err:    reAuthenticationRequired,
					},
				},
				ctx: context.Background(),
				at:  "test-access-token",
			},
			&action{
				claims:    nil,
				errStatus: http.StatusUnauthorized,
				err:       reAuthenticationRequired,
			},
		),
		gen(
			"local validation / valid minimal",
			[]string{cndValidateSign, cndValidateIss, cndValidateAud},
			[]string{actCheckNoError},
			&condition{
				oc: &oauthContext{
					lg: log.GlobalLogger(log.DefaultLoggerName),
					provider: &provider{
						issuer: "test-iss",
					},
					client: &client{
						id:       "test-id",
						audience: "test-aud",
					},
					jh: newJWTHandler(),
				},
				ctx: context.Background(),
				at:  testSimpleJWT,
			},
			&action{
				claims: jwt.MapClaims{
					"active": true, // This claim will be added by the method.
					"iss":    "test-iss",
					"aud":    "test-aud",
					"azp":    "test-id",
				},
				err: nil,
			},
		),
		gen(
			"local validation / invalid JWT signature",
			[]string{cndValidateSign, cndValidateIss, cndValidateAud},
			[]string{actCheckError},
			&condition{
				oc: &oauthContext{
					lg: log.GlobalLogger(log.DefaultLoggerName),
					provider: &provider{
						issuer: "test-iss",
					},
					client: &client{
						id:       "test-id",
						audience: "test-aud",
					},
					jh: newJWTHandler(),
				},
				ctx: context.Background(),
				// token is signed by HS256 with password "wrong-password".
				at: testWringSigJWT,
			},
			&action{
				claims:    nil,
				errStatus: http.StatusUnauthorized,
				err:       reAuthenticationRequired,
			},
		),
		gen(
			"local validation / invalid iss claim",
			[]string{cndValidateSign, cndValidateIss, cndValidateAud},
			[]string{actCheckError},
			&condition{
				oc: &oauthContext{
					lg: log.GlobalLogger(log.DefaultLoggerName),
					provider: &provider{
						issuer: "test-iss-xxx",
					},
					client: &client{
						id:       "test-id",
						audience: "test-aud",
					},
					jh:           newJWTHandler(),
					atParseOpts:  []jwt.ParserOption{jwt.WithIssuer("test-iss-xxx"), jwt.WithAudience("test-aud")},
					idtParseOpts: []jwt.ParserOption{jwt.WithIssuer("test-iss-xxx"), jwt.WithAudience("test-aud")},
				},
				ctx: context.Background(),
				at:  testSimpleJWT,
			},
			&action{
				claims:    nil,
				errStatus: http.StatusUnauthorized,
				err:       reAuthenticationRequired,
			},
		),
		gen(
			"local validation / invalid aud claim",
			[]string{cndValidateSign, cndValidateIss, cndValidateAud},
			[]string{actCheckError},
			&condition{
				oc: &oauthContext{
					lg: log.GlobalLogger(log.DefaultLoggerName),
					provider: &provider{
						issuer: "test-iss",
					},
					client: &client{
						id:       "test-id",
						audience: "test-aud-xxx",
					},
					jh:           newJWTHandler(),
					atParseOpts:  []jwt.ParserOption{jwt.WithIssuer("test-iss"), jwt.WithAudience("test-aud-xxx")},
					idtParseOpts: []jwt.ParserOption{jwt.WithIssuer("test-iss"), jwt.WithAudience("test-aud-xxx")},
				},
				ctx: context.Background(),
				at:  testSimpleJWT,
			},
			&action{
				claims:    nil,
				errStatus: http.StatusUnauthorized,
				err:       reAuthenticationRequired,
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			claims, err := tt.C().oc.validateAT(tt.C().ctx, tt.C().at)

			testutil.Diff(t, tt.A().claims, claims)

			if tt.A().err == nil {
				testutil.Diff(t, nil, err)
				return
			}

			testutil.Diff(t, tt.A().errStatus, err.StatusCode())
			if tt.A().err == reAuthenticationRequired {
				testutil.Diff(t, tt.A().err.Error(), err.Error())
			} else {
				e := err.(*utilhttp.HTTPError)
				testutil.DiffError(t, tt.A().err, tt.A().errPattern, e.Unwrap())
			}
		})
	}
}

func TestOauthContext_validateIDT(t *testing.T) {
	type condition struct {
		oc   *oauthContext
		idt  string
		opts []validateOption
	}

	type action struct {
		claims     jwt.MapClaims
		errStatus  int
		err        any
		errPattern *regexp.Regexp
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	cndValidateSign := tb.Condition("validate signature", "validate jwt signature")
	cndValidateIss := tb.Condition("validate iss", "validate iss claim in the JWT")
	cndValidateAud := tb.Condition("validate aud", "validate aud claim in the JWT")
	cndValidateAzp := tb.Condition("validate azp", "validate azp claim in the JWT")
	cndValidateOption := tb.Condition("validate options", "do optional validation")
	actCheckError := tb.Action("error", "check that an error was returned")
	actCheckNoError := tb.Action("no error", "check that no error was returned")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"valid minimal",
			[]string{cndValidateSign, cndValidateIss, cndValidateAud, cndValidateAzp},
			[]string{actCheckNoError},
			&condition{
				oc: &oauthContext{
					provider: &provider{
						issuer: "test-iss",
					},
					client: &client{
						id:       "test-id",
						audience: "test-aud",
					},
					jh: newJWTHandler(),
				},
				idt: testSimpleJWT,
			},
			&action{
				claims: jwt.MapClaims{
					"iss": "test-iss",
					"aud": "test-aud",
					"azp": "test-id",
				},
				err: nil,
			},
		),
		gen(
			"valid with option",
			[]string{cndValidateSign, cndValidateIss, cndValidateAud, cndValidateAzp, cndValidateOption},
			[]string{actCheckNoError},
			&condition{
				oc: &oauthContext{
					provider: &provider{
						issuer: "test-iss",
					},
					client: &client{
						id:       "test-id",
						audience: "test-aud",
					},
					jh: newJWTHandler(),
				},
				idt: testSimpleJWT,
			},
			&action{
				claims: jwt.MapClaims{
					"iss": "test-iss",
					"aud": "test-aud",
					"azp": "test-id",
				},
				err: nil,
			},
		),
		gen(
			"invalid JWT signature",
			[]string{cndValidateSign, cndValidateIss, cndValidateAud},
			[]string{actCheckError},
			&condition{
				oc: &oauthContext{
					provider: &provider{
						issuer: "test-iss-xxx",
					},
					client: &client{
						id:       "test-id",
						audience: "test-aud",
					},
					jh: newJWTHandler(),
				},
				// token is signed by HS256 with password "wrong-password".
				idt: testWringSigJWT,
			},
			&action{
				claims:     nil,
				errStatus:  http.StatusUnauthorized,
				err:        app.ErrAppAuthnInvalidToken,
				errPattern: regexp.MustCompile(core.ErrPrefix + `invalid id token. token validation failed. .* signature is invalid`),
			},
		),
		gen(
			"invalid iss claim",
			[]string{cndValidateSign, cndValidateIss, cndValidateAud},
			[]string{actCheckError},
			&condition{
				oc: &oauthContext{
					provider: &provider{
						issuer: "test-iss-xxx",
					},
					client: &client{
						id:       "test-id",
						audience: "test-aud",
					},
					jh:           newJWTHandler(),
					atParseOpts:  []jwt.ParserOption{jwt.WithIssuer("test-iss-xxx"), jwt.WithAudience("test-aud")},
					idtParseOpts: []jwt.ParserOption{jwt.WithIssuer("test-iss-xxx"), jwt.WithAudience("test-aud")},
				},
				idt: testSimpleJWT,
			},
			&action{
				claims:     nil,
				errStatus:  http.StatusUnauthorized,
				err:        app.ErrAppAuthnInvalidToken,
				errPattern: regexp.MustCompile(core.ErrPrefix + `invalid id token. token validation failed. .* token has invalid issuer`),
			},
		),
		gen(
			"invalid aud claim",
			[]string{cndValidateSign, cndValidateIss, cndValidateAud},
			[]string{actCheckError},
			&condition{
				oc: &oauthContext{
					provider: &provider{
						issuer: "test-iss",
					},
					client: &client{
						id:       "test-id",
						audience: "test-aud-xxx",
					},
					jh:           newJWTHandler(),
					atParseOpts:  []jwt.ParserOption{jwt.WithIssuer("test-iss"), jwt.WithAudience("test-aud-xxx")},
					idtParseOpts: []jwt.ParserOption{jwt.WithIssuer("test-iss"), jwt.WithAudience("test-aud-xxx")},
				},
				idt: testSimpleJWT,
			},
			&action{
				claims:     nil,
				errStatus:  http.StatusUnauthorized,
				err:        app.ErrAppAuthnInvalidToken,
				errPattern: regexp.MustCompile(core.ErrPrefix + `invalid id token. token validation failed. .* token has invalid audience`),
			},
		),
		gen(
			"invalid azp claim",
			[]string{cndValidateSign, cndValidateIss, cndValidateAud, cndValidateAzp},
			[]string{actCheckError},
			&condition{
				oc: &oauthContext{
					provider: &provider{
						issuer: "test-iss",
					},
					client: &client{
						id:       "test-id-xxx",
						audience: "test-aud",
					},
					jh: newJWTHandler(),
				},
				idt: testSimpleJWT,
			},
			&action{
				claims:     nil,
				errStatus:  http.StatusUnauthorized,
				err:        app.ErrAppAuthnInvalidToken,
				errPattern: regexp.MustCompile(core.ErrPrefix + `invalid id token. azp in ID token does not match to client id.`),
			},
		),
		gen(
			"option validation failed",
			[]string{cndValidateSign, cndValidateIss, cndValidateAud, cndValidateAzp, cndValidateOption},
			[]string{actCheckError},
			&condition{
				oc: &oauthContext{
					provider: &provider{
						issuer: "test-iss",
					},
					client: &client{
						id:       "test-id",
						audience: "test-aud",
					},
					jh:           newJWTHandler(),
					atParseOpts:  []jwt.ParserOption{jwt.WithIssuer("test-iss"), jwt.WithAudience("test-aud")},
					idtParseOpts: []jwt.ParserOption{jwt.WithIssuer("test-iss"), jwt.WithAudience("test-aud")},
				},
				idt:  testSimpleJWT,
				opts: []validateOption{maxAgeValidation(300)},
			},
			&action{
				claims:     nil,
				errStatus:  http.StatusUnauthorized,
				err:        app.ErrAppAuthnInvalidToken,
				errPattern: regexp.MustCompile(core.ErrPrefix + `invalid ID token. auth_time not found.`),
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			claims, err := tt.C().oc.validateIDT(tt.C().idt, tt.C().opts)

			testutil.Diff(t, tt.A().claims, claims)

			if tt.A().err == nil {
				testutil.Diff(t, nil, err)
				return
			}

			testutil.Diff(t, tt.A().errStatus, err.StatusCode())
			if tt.A().err == reAuthenticationRequired {
				testutil.Diff(t, reAuthenticationRequired.Error(), err.Error())
			} else {
				e := err.(*utilhttp.HTTPError)
				testutil.DiffError(t, tt.A().err, tt.A().errPattern, e.Unwrap())
			}
		})
	}
}

func TestMaxAgeValidation(t *testing.T) {
	type condition struct {
		opt    validateOption
		claims jwt.MapClaims
	}

	type action struct {
		errStatus  int
		err        any
		errPattern *regexp.Regexp
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	cndSkip := tb.Condition("skip", "set max age 0 to skip validation")
	cndMaxAge := tb.Condition("auth_time exists", "set positive maxAge to verify auth_time")
	cndAuthTime := tb.Condition("auth_time exists", "auth_time claim is exist in the claims")
	cndAuthTimeInvalid := tb.Condition("auth_time invalid", "invalid data type of auth_time claim")
	actCheckError := tb.Action("error", "check that an error was returned")
	actCheckNoError := tb.Action("no error", "check that no error was returned")
	actCheckReAuth := tb.Action("re-auth", "check that re-authentication was required")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"no max age",
			[]string{cndSkip, cndAuthTime},
			[]string{actCheckNoError},
			&condition{
				opt: maxAgeValidation(0),
				claims: jwt.MapClaims{
					"auth_time": float64(1704034800), // Too old.
				},
			},
			&action{
				errStatus: 0,
				err:       nil,
			},
		),
		gen(
			"fresh auth time",
			[]string{cndMaxAge, cndAuthTime},
			[]string{actCheckNoError},
			&condition{
				opt: maxAgeValidation(300),
				claims: jwt.MapClaims{
					"auth_time": float64(time.Now().Unix() - 100), // Fresh.
				},
			},
			&action{
				errStatus: 0,
				err:       nil,
			},
		),
		gen(
			"no auth time claim",
			[]string{cndMaxAge},
			[]string{actCheckError},
			&condition{
				opt:    maxAgeValidation(300),
				claims: jwt.MapClaims{},
			},
			&action{
				errStatus:  http.StatusUnauthorized,
				err:        app.ErrAppAuthnInvalidToken,
				errPattern: regexp.MustCompile(core.ErrPrefix + `invalid ID token. auth_time not found.`),
			},
		),
		gen(
			"invalid auth time claim",
			[]string{cndMaxAge, cndAuthTime, cndAuthTimeInvalid},
			[]string{actCheckError},
			&condition{
				opt: maxAgeValidation(300),
				claims: jwt.MapClaims{
					"auth_time": "1704034800", // String will not accepted.
				},
			},
			&action{
				errStatus:  http.StatusUnauthorized,
				err:        app.ErrAppAuthnInvalidToken,
				errPattern: regexp.MustCompile(core.ErrPrefix + `invalid ID token. auth_time not found.`),
			},
		),
		gen(
			"old auth time",
			[]string{cndMaxAge, cndAuthTime},
			[]string{actCheckError, actCheckReAuth},
			&condition{
				opt: maxAgeValidation(300),
				claims: jwt.MapClaims{
					"auth_time": float64(1704034800), // Too old.
				},
			},
			&action{
				errStatus: http.StatusUnauthorized,
				err:       reAuthenticationRequired,
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			err := tt.C().opt.validate(tt.C().claims)
			if tt.A().err == nil {
				testutil.Diff(t, nil, err)
				return
			}

			testutil.Diff(t, tt.A().errStatus, err.StatusCode())
			if tt.A().err == reAuthenticationRequired {
				testutil.Diff(t, reAuthenticationRequired.Error(), err.Error())
			} else {
				e := err.(*utilhttp.HTTPError)
				testutil.DiffError(t, tt.A().err, tt.A().errPattern, e.Unwrap())
			}
		})
	}
}

func TestNonceValidation(t *testing.T) {
	type condition struct {
		opt    validateOption
		claims jwt.MapClaims
	}

	type action struct {
		errStatus  int
		err        any
		errPattern *regexp.Regexp
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	cndSkip := tb.Condition("skip", "set empty string for nonce to skip verification")
	cndValidateNonce := tb.Condition("validate nonce", "set non empty nonce to validate nonce claim")
	cndNonceClaims := tb.Condition("nonce claim", "set non empty nonce claim in the claims")
	cndNonceInvalid := tb.Condition("invalid nonce", "set invalid nonce claim in the claims")
	actCheckError := tb.Action("error", "check that an error was returned")
	actCheckNoError := tb.Action("no error", "check that no error was returned")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"no nonce",
			[]string{cndSkip, cndNonceClaims},
			[]string{actCheckNoError},
			&condition{
				opt: nonceValidation(""),
				claims: jwt.MapClaims{
					"nonce": "test_nonce_value",
				},
			},
			&action{
				errStatus: 0,
				err:       nil,
			},
		),
		gen(
			"valid nonce",
			[]string{cndValidateNonce, cndNonceClaims},
			[]string{actCheckNoError},
			&condition{
				opt: nonceValidation("test_nonce_value"),
				claims: jwt.MapClaims{
					"nonce": "test_nonce_value",
				},
			},
			&action{
				errStatus: 0,
				err:       nil,
			},
		),
		gen(
			"no nonce claim",
			[]string{cndValidateNonce},
			[]string{actCheckError},
			&condition{
				opt:    nonceValidation("test_nonce_value"),
				claims: jwt.MapClaims{},
			},
			&action{
				errStatus:  http.StatusUnauthorized,
				err:        app.ErrAppAuthnInvalidToken,
				errPattern: regexp.MustCompile(core.ErrPrefix + `invalid ID token. invalid nonce.`),
			},
		),
		gen(
			"invalid nonce claim",
			[]string{cndValidateNonce, cndNonceClaims, cndNonceInvalid},
			[]string{actCheckError},
			&condition{
				opt: nonceValidation("test_nonce_value"),
				claims: jwt.MapClaims{
					"nonce": "invalid_nonce_value",
				},
			},
			&action{
				errStatus:  http.StatusUnauthorized,
				err:        app.ErrAppAuthnInvalidToken,
				errPattern: regexp.MustCompile(core.ErrPrefix + `invalid ID token. invalid nonce.`),
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			err := tt.C().opt.validate(tt.C().claims)
			if tt.A().err == nil {
				testutil.Diff(t, nil, err)
				return
			}

			testutil.Diff(t, tt.A().errStatus, err.StatusCode())
			if tt.A().err == reAuthenticationRequired {
				testutil.Diff(t, reAuthenticationRequired.Error(), err.Error())
			} else {
				e := err.(*utilhttp.HTTPError)
				testutil.DiffError(t, tt.A().err, tt.A().errPattern, e.Unwrap())
			}
		})
	}
}
func TestNewMapWithKeys(t *testing.T) {
	type condition struct {
		keys []string
		src  map[string]any
	}

	type action struct {
		result map[string]any
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	cndKeyExists := tb.Condition("key exists", "key exists in the given source map")
	cndKeyProvided := tb.Condition("key provided", "non empty key list is provided")
	cndMapProvided := tb.Condition("map provided", "non-nil source map is provided")
	actCheckValue := tb.Action("check value", "check that the matched values are returned by the map")
	actCheckEmpty := tb.Action("check empty", "check that empty map was returned")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"all keys found",
			[]string{cndKeyProvided, cndMapProvided, cndKeyExists},
			[]string{actCheckValue},
			&condition{
				src: map[string]any{
					"key1": "string",
					"key2": []byte("byte"),
					"key3": 123,
					"key4": 123.45,
					"key5": map[string]any{"foo": "bar"},
				},
				keys: []string{"key1", "key2", "key3", "key4", "key5"},
			},
			&action{
				result: map[string]any{
					"key1": "string",
					"key2": []byte("byte"),
					"key3": 123,
					"key4": 123.45,
					"key5": map[string]any{"foo": "bar"},
				},
			},
		),
		gen(
			"keys found",
			[]string{cndKeyProvided, cndMapProvided, cndKeyExists},
			[]string{actCheckValue},
			&condition{
				src: map[string]any{
					"key1": "string",
					"key2": []byte("byte"),
					"key3": 123,
					"key4": 123.45,
					"key5": map[string]any{"foo": "bar"},
				},
				keys: []string{"key5"},
			},
			&action{
				result: map[string]any{
					"key5": map[string]any{"foo": "bar"},
				},
			},
		),
		gen(
			"keys not found",
			[]string{cndKeyProvided, cndMapProvided, cndKeyExists},
			[]string{actCheckEmpty},
			&condition{
				src: map[string]any{
					"key1": "string",
					"key2": []byte("byte"),
					"key3": 123,
					"key4": 123.45,
					"key5": map[string]any{"foo": "bar"},
				},
				keys: []string{"no-key1", "no-key2", "no-key3", "no-key4", "no-key5"},
			},
			&action{
				result: map[string]any{},
			},
		),
		gen(
			"keys not provided",
			[]string{cndMapProvided},
			[]string{actCheckEmpty},
			&condition{
				src: map[string]any{
					"key1": "string",
					"key2": []byte("byte"),
					"key3": 123,
					"key4": 123.45,
					"key5": map[string]any{"foo": "bar"},
				},
				keys: nil,
			},
			&action{
				result: map[string]any{},
			},
		),
		gen(
			"source map not provided",
			[]string{cndKeyProvided, cndKeyExists},
			[]string{actCheckEmpty},
			&condition{
				src:  nil,
				keys: []string{"key1", "key2", "key3", "key4", "key5"},
			},
			&action{
				result: map[string]any{},
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			result := newMapWithKeys(tt.C().keys, tt.C().src)
			testutil.Diff(t, tt.A().result, result)
		})
	}
}
