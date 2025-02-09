package cors

import (
	"net/http"
	"slices"
)

type corsPolicy struct {
	// allowedOrigins is the list of origins to allow.
	// Specified origin must exactly match the origin of requests.
	allowedOrigins []string

	// allowedMethods is the list of http methods to allow.
	// All methods must be listed.
	allowedMethods []string

	// embedderPolicy is the value for "Cross-Origin-Embedder-Policy" header.
	// 	- https://docs.w3cub.com/http/headers/cross-origin-embedder-policy
	//  - https://html.spec.whatwg.org/multipage/browsers.html#the-coep-headers
	embedderPolicy string

	// openerPolicy is the value for "Cross-Origin-Opener-Policy" header.
	// 	- https://docs.w3cub.com/http/headers/cross-origin-opener-policy
	//  - https://html.spec.whatwg.org/multipage/browsers.html#cross-origin-opener-policies
	openerPolicy string

	// resourcePolicy is the value for "Cross-Origin-Resource-Policy" header.
	// 	- https://docs.w3cub.com/http/headers/cross-origin-opener-policy
	//  - https://fetch.spec.whatwg.org/#cross-origin-resource-policy-header
	resourcePolicy string

	// allowCredentials allows client javascript to read credentials.
	// Credentials are cookies, authorization headers, or TLS client certificates.
	// When this flag is set to true, "Access-Control-Allow-Credentials: true" is set to the response header.
	// 	- https://docs.w3cub.com/http/headers/access-control-allow-credentials
	allowCredentials bool

	// joinedAllowedMethods is the string of joined string of multiple header values.
	// Values are joined to avoid join strings operation for each request.
	// This value is set to "Access-Control-Allow-Methods" response header.
	// 	- https://docs.w3cub.com/http/headers/access-control-allow-methods
	joinedAllowedMethods string

	// joinedAllowedHeaders is the string of joined string of multiple header values.
	// Values are joined to avoid join strings operation for each request.
	// This value is set to "Access-Control-Allow-Headers" response header.
	// 	- https://docs.w3cub.com/http/headers/access-control-allow-headers
	joinedAllowedHeaders string

	// joinedExposedHeaders is the string of joined string of multiple header values.
	// Values are joined to avoid join strings operation for each request.
	// This value is set to "Access-Control-Expose-Headers" response header.
	// 	- https://docs.w3cub.com/http/headers/access-control-expose-headers
	joinedExposedHeaders string

	// maxAge is the time that clients are allowed to cache.
	// This value is set to "Access-Control-Max-Age" response header.
	// Header is not set when this value is empty.
	// 	- https://docs.w3cub.com/http/headers/access-control-max-age
	// 	- https://www.w3.org/TR/2020/SPSD-cors-20200602/#preflight-result-cache-0
	maxAge string

	// allowPrivateNetwork allows client to share resources with external networks.
	// When this flag is set to true, "Access-Control-Allow-Private-Network: true" is set to the response header.
	// 	- https://wicg.github.io/private-network-access/
	allowPrivateNetwork bool

	// disableWildCardOrigin if true, set the given origin to the
	// "Access-Control-Allow-Origin" header rather than "*".
	// This is, in most cases, insecure than the wildcard origin "*".
	disableWildCardOrigin bool
}

// originAllowed returns if the request is allowed by this CORS policy or not.
// Do not accept the request denied by this policy.
// Make sure not to use requesters origin if it was not matched to a specific origin.
// The concept of origin is defined in the following document.
//   - https://datatracker.ietf.org/doc/rfc6454/
func (p *corsPolicy) originAllowed(r *http.Request, wh http.Header) bool {
	origin := r.Header.Get("Origin")
	// Origin should not be empty.
	// So, deny it even simple requests do not contains origin headers.
	// 	- https://datatracker.ietf.org/doc/rfc6454/
	if origin == "" {
		return false
	}

	if slices.Contains(p.allowedOrigins, origin) {
		wh.Set("Access-Control-Allow-Origin", origin)
		return true
	}

	if slices.Contains(p.allowedOrigins, "*") {
		if p.disableWildCardOrigin {
			wh.Set("Access-Control-Allow-Origin", origin)
		} else {
			wh.Set("Access-Control-Allow-Origin", "*")
		}
		return true
	}

	return false
}

// handlePreflight handles preflight requests, or OPTIONS requests.
// We can return 200, 204, 403 status to the preflight request.
//
// References
//   - https://fetch.spec.whatwg.org/#cors-preflight-fetch-0
//   - https://docs.w3cub.com/http/headers/access-control-allow-credentials
//   - https://docs.w3cub.com/http/headers/access-control-allow-headers
//   - https://docs.w3cub.com/http/headers/access-control-allow-origin
//   - https://docs.w3cub.com/http/headers/access-control-expose-headers
//   - https://docs.w3cub.com/http/headers/access-control-max-age
//   - https://docs.w3cub.com/http/headers/access-control-request-headers
//   - https://docs.w3cub.com/http/headers/access-control-request-method
//   - https://wicg.github.io/private-network-access
func (p *corsPolicy) handlePreflight(w http.ResponseWriter, r *http.Request) int {
	wh := w.Header()

	// First of all, always set Origin in the Vary header.
	// 	- https://github.com/rs/cors/issues/10
	// 	- https://github.com/rs/cors/commit/dbdca4d95feaa7511a46e6f1efb3b3aa505bc43f#commitcomment-12352001
	wh.Add("Vary", "Origin")

	// Always returned all allowed methods.
	// Server must not return wildcard "*" as "Access-Control-Allow-Methods" header value when the clients
	// provide credentials such as cookies, TLS client certificates, and authentication entries.
	// Here, credentials do not need to be checked because p.joinedAllowedMethods does not contain wildcard "*".
	wh.Set("Access-Control-Allow-Methods", p.joinedAllowedMethods)
	wh.Add("Vary", "Access-Control-Request-Method")

	// Always set vary header.
	wh.Add("Vary", "Access-Control-Request-Headers")
	if r.Header.Get("Access-Control-Request-Headers") != "" {
		wh.Set("Access-Control-Allow-Headers", p.joinedAllowedHeaders)
	}

	if p.joinedExposedHeaders != "" {
		wh.Set("Access-Control-Expose-Headers", p.joinedExposedHeaders)
	}

	if p.allowCredentials {
		wh.Set("Access-Control-Allow-Credentials", "true")
	}

	// Max-Age is the duration that browsers can cache.
	if p.maxAge != "" {
		wh.Set("Access-Control-Max-Age", p.maxAge)
	}

	// "Access-Control-Request-Private-Network" does not appear in the fetch standard
	// but it would be useful if it can be set.
	// This header is included in the response to preflight requests and not in actual requests.
	if r.Header.Get("Access-Control-Request-Private-Network") != "" {
		wh.Add("Vary", "Access-Control-Request-Private-Network")
		if p.allowPrivateNetwork {
			wh.Set("Access-Control-Allow-Private-Network", "true")
		}
	}

	// Access-Control-* headers are added. Now check the request origin.
	// Forbidden if the origin is not allowed.
	if !p.originAllowed(r, w.Header()) {
		// Return 403 rather than returning 200 or 204
		// not to reveal server's information from the stand point of security.
		return http.StatusForbidden
	}

	return http.StatusOK
}

// handleActualRequest handles actual requests, or non preflight.
// This method called only when originAllowed is true.
//
// References
//   - https://fetch.spec.whatwg.org/#http-responses
//   - https://docs.w3cub.com/http/headers/access-control-allow-methods
//   - https://html.spec.whatwg.org/multipage/browsers.html#cross-origin-opener-policies
//   - https://html.spec.whatwg.org/multipage/browsers.html#coep
//   - https://docs.w3cub.com/http/headers/cross-origin-embedder-policy
//   - https://docs.w3cub.com/http/headers/cross-origin-opener-policy
//   - https://docs.w3cub.com/http/headers/cross-origin-resource-policy
func (p *corsPolicy) handleActualRequest(w http.ResponseWriter, r *http.Request) int {
	wh := w.Header()

	// First of all, always set Origin in the Vary header.
	// 	- https://github.com/rs/cors/issues/10
	// 	- https://github.com/rs/cors/commit/dbdca4d95feaa7511a46e6f1efb3b3aa505bc43f#commitcomment-12352001
	wh.Add("Vary", "Origin")

	// Always returned all allowed methods.
	// Even this header is not required in the fetch standard.
	wh.Set("Access-Control-Allow-Methods", p.joinedAllowedMethods)

	if p.joinedAllowedHeaders != "" {
		wh.Set("Access-Control-Allow-Headers", p.joinedAllowedHeaders)
	}
	if p.joinedExposedHeaders != "" {
		wh.Set("Access-Control-Expose-Headers", p.joinedExposedHeaders)
	}

	// Setting the "Access-Control-Allow-Credentials" is not allowed
	// 	- when "Access-Control-Allow-Origin" = "*"
	// 	- when "Access-Control-Allow-Headers" contains "*"
	if p.allowCredentials && wh.Get("Access-Control-Allow-Origin") != "*" {
		wh.Set("Access-Control-Allow-Credentials", "true")
	}

	// Access-Control-* headers are added. Now check the request origin.
	// Forbidden if the origin is not allowed.
	if !p.originAllowed(r, wh) {
		return http.StatusForbidden
	}

	// Access-Control-* headers are added. Now check the request method.
	// Forbidden if the method is not allowed.
	if !slices.Contains(p.allowedMethods, r.Method) {
		return http.StatusForbidden
	}

	if p.embedderPolicy != "" {
		wh.Set("Cross-Origin-Embedder-Policy", p.embedderPolicy)
	}
	if p.openerPolicy != "" {
		wh.Set("Cross-Origin-Opener-Policy", p.openerPolicy)
	}
	if p.resourcePolicy != "" {
		wh.Set("Cross-Origin-Resource-Policy", p.resourcePolicy)
	}

	return http.StatusOK
}
