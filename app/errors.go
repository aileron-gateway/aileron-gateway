// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package app

import (
	"github.com/aileron-gateway/aileron-gateway/kernel/errorutil"
)

var (
	// ---------------------------------------------------------
	// general: E3000 - E3049
	ErrAppGenCreateRequest    = errorutil.NewKind("E3000", "AppGenCreateRequest", "failed to create http request. method={{method}} url={{url}} body={{body}}")
	ErrAppGenRoundTrip        = errorutil.NewKind("E3001", "AppGenRoundTrip", "failed to round trip. method={{method}} url={{url}} body={{body}}")
	ErrAppGenInvalidResponse  = errorutil.NewKind("E3002", "AppGenInvalidResponse", "invalid response. method={{method}} url={{url}} Content-Type={{type}} Status={{status}} body={{body}}")
	ErrAppGenReadHTTPBody     = errorutil.NewKind("E3003", "AppGenReadHTTPBody", "failed to read {{direction}} body. read={{body}}")
	ErrAppGenSessionOperation = errorutil.NewKind("E3004", "AppGenSessionOperation", "session operation failed. {{operation}} {{reason}}")
	ErrAppGenUnmarshal        = errorutil.NewKind("E3005", "AppGenUnmarshal", "failed to unmarshal from={{from}} to={{to}} {{content}}")
	// ---------------------------------------------------------

	// ---------------------------------------------------------
	// app/authn: E3050 - E3099
	ErrAppAuthnAuthentication          = errorutil.NewKind("E3050", "AppAuthnAuthentication", "authentication failed")
	ErrAppAuthnParseWithClaims         = errorutil.NewKind("E3051", "AppAuthnParseWithClaims", "failed to parse JWT claims. {{jwt}}")
	ErrAppAuthnRedeemToken             = errorutil.NewKind("E3052", "AppAuthnRedeemToken", "failed to redeem token. {{info}}")
	ErrAppAuthnRefreshToken            = errorutil.NewKind("E3053", "AppAuthnRefreshToken", "failed to refresh token. {{info}}")
	ErrAppAuthnIntrospection           = errorutil.NewKind("E3054", "AppAuthnIntrospection", "failed to token introspection. {{info}}")
	ErrAppAuthnInvalidToken            = errorutil.NewKind("E3055", "AppAuthnInvalidToken", "invalid {{name}}. {{reason}} {{token}}")
	ErrAppAuthnInvalidCredential       = errorutil.NewKind("E3056", "AppAuthnInvalidCredential", "invalid credential for {{purpose}}. {{reason}}")
	ErrAppAuthnInvalidParameters       = errorutil.NewKind("E3057", "AppAuthnInvalidParameters", "invalid {{name}} parameters. {{reason}}")
	ErrAppAuthnGenerateCSRFParams      = errorutil.NewKind("E3058", "AppAuthnGenerateCSRFParams", "failed to generate csrf states")
	ErrAppAuthnNoSession               = errorutil.NewKind("E3059", "AppAuthnNoSession", "session not found in the context")
	ErrAppAuthnGenerateRequestObject   = errorutil.NewKind("E3060", "AppAuthnGenerateRequestObject", "failed to generate request object")
	ErrAppAuthnGenerateClientAssertion = errorutil.NewKind("E3061", "AppAuthnGenerateClientAssertion", "failed to generate client assertion. {{reason}}")
	ErrAppAuthnGenerateTokenWithClaims = errorutil.NewKind("E3062", "AppAuthnGenerateTokenWithClaims", "failed to generate JWT claims.")
	ErrAppAuthnSignToken               = errorutil.NewKind("E3063", "AppAuthnSignToken", "failed to sign JWT claims.")
	ErrAppAuthnUserInfo                = errorutil.NewKind("E3064", "AppAuthnUserInfo", "failed to userinfo request. {{info}}")
	ErrAppAuthnInvalidUserInfo         = errorutil.NewKind("E3065", "AppAuthnInvalidUserInfo", "invalid userinfo. {{reason}}")
	// ---------------------------------------------------------

	// ---------------------------------------------------------
	// app/authz: E3100 - E3149
	ErrAppAuthzAuthorization = errorutil.NewKind("E3100", "AppAuthzAuthorization", "authorization failed")
	ErrAppAuthzForbidden     = errorutil.NewKind("E3101", "ErrAppAuthzForbidden", "forbidden on authorization")
	// ---------------------------------------------------------

	// ---------------------------------------------------------
	// app/handler: E3150 - E3199
	// ---------------------------------------------------------

	// ---------------------------------------------------------
	// app/middleware: E3200 - E3249
	ErrAppMiddleGenID              = errorutil.NewKind("E3201", "AppMiddleGenID", "failed to generate {{type}} ID")
	ErrAppMiddleAPITimeout         = errorutil.NewKind("E3202", "AppMiddleAPITimeout", "api timeout occurred")
	ErrAppMiddlePanicRecovered     = errorutil.NewKind("E3203", "AppMiddlePanicRecovered", "panic recovered")
	ErrAppMiddleCORSForbidden      = errorutil.NewKind("E3204", "AppMiddleCORSForbidden", "forbidden by cors policy")
	ErrAppMiddleBodyTooLarge       = errorutil.NewKind("E3205", "AppMiddleBodyTooLarge", "request body size exceeded the body limit.")
	ErrAppMiddleInvalidLength      = errorutil.NewKind("E3206", "AppMiddleInvalidLength", "Content-Length header invalid or not found.")
	ErrAppMiddleBodyLimit          = errorutil.NewKind("E3207", "AppMiddleBodyLimit", "error limiting body size")
	ErrAppMiddleCSRFNewToken       = errorutil.NewKind("E3208", "AppMiddleCSRFNewToken", "failed to create new CSRF token.")
	ErrAppMiddleCSRFToken          = errorutil.NewKind("E3209", "AppMiddleCSRFToken", "error on checking CSRF token..")
	ErrAppMiddleCSRFSession        = errorutil.NewKind("E3210", "AppMiddleCSRFSession", "csrf session operation error.")
	ErrAppMiddleHeaderPolicy       = errorutil.NewKind("E3211", "AppMiddleHeaderPolicy", "header policy error. {{reason}}")
	ErrAppMiddleSession            = errorutil.NewKind("E3212", "AppMiddleSession", "session operation failed.")
	ErrAppMiddleThrottle           = errorutil.NewKind("E3213", "AppMiddleThrottle", "too many requests.")
	ErrAppMiddleInvalidCert        = errorutil.NewKind("E3214", "AppMiddleInvalidCert", "client certificate invalid or not found. {{reason}}")
	ErrAppMiddleInvalidFingerprint = errorutil.NewKind("E3214", "AppMiddleInvalidFingerprint", "fingerprint invalid.")
	// ---------------------------------------------------------

	// ---------------------------------------------------------
	// app/storage: E3250 - E3299
	ErrAppStorageKVS = errorutil.NewKind("E3250", "AppStorageKVS", "error occurred while operating the key-value store")
	// ---------------------------------------------------------

	// ---------------------------------------------------------
	// app/tracer: E3300 - E3349
	// ---------------------------------------------------------

	// ---------------------------------------------------------
	// app/util: E3400 - E3449
	ErrAppUtilGenerateJWTKey = errorutil.NewKind("E3400", "AppUtilGenerateJWTKey", "failed to generate JWT key")
	ErrAppUtilGetJWKSet      = errorutil.NewKind("E3401", "AppUtilGetJWKSet", "failed to get JWK sets from JWKs endpoint")
	// ---------------------------------------------------------
)
