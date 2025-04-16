// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package httpproxy

import (
	"net"
	"net/http"
	"net/textproto"
	"net/url"
	"strings"

	"golang.org/x/net/http/httpguts"
)

// upgradeType returns the string of protocol to upgrade.
// "Upgrade" header must be present at most 1 in the given h.
// The header key in the given h must be canonicalized.
// i.e. "Upgrade" not "upgrade", "Connection" not "connection".
//
// Reference
//   - https://developer.mozilla.org/en-US/docs/Web/HTTP/Protocol_upgrade_mechanism
//   - https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Upgrade
func upgradeType(h http.Header) string {
	if !httpguts.HeaderValuesContainsToken(h["Connection"], "Upgrade") {
		return ""
	}
	if v := h["Upgrade"]; len(v) > 0 {
		return v[0] // "Get" a value.
	}
	return ""
}

// copyHeader copies headers from src to dst.
// The destination header, dst, must not be nil.
// Copied keys are modified with textproto.CanonicalMIMEHeaderKey if necessary.
func copyHeader(dst, src http.Header) {
	for k, v := range src {
		k = textproto.CanonicalMIMEHeaderKey(k)
		dst[k] = append(dst[k], v...) // "Add" values.
	}
}

// copyTrailer copies trailers from src to dst.
// Header keys are copied with the prefixed defined as http.TrailerPrefix.
// Copied keys are modified with textproto.CanonicalMIMEHeaderKey if necessary.
func copyTrailer(dst, src http.Header) {
	for k, v := range src {
		if !strings.HasPrefix(k, http.TrailerPrefix) {
			k = http.TrailerPrefix + k
		}
		dst[k] = append(dst[k], v...) // "Add" values.
	}
}

// removeHopByHopHeaders removes hop-by-hop headers.
// All header keys removed in removeHopByHopHeaders
// must be canonical format in the given header.
//
// Reference
//   - https://go.dev/src/net/http/httputil/reverseproxy.go
//   - https://datatracker.ietf.org/doc/rfc7230/
//   - https://datatracker.ietf.org/doc/rfc2616/
func removeHopByHopHeaders(h http.Header) {
	// RFC 7230, section 6.1: Remove headers listed in the "Connection" header.
	for _, conn := range h["Connection"] {
		for _, c := range strings.Split(conn, ",") {
			delete(h, textproto.CanonicalMIMEHeaderKey(textproto.TrimString(c))) // "Del" a header.
		}
	}
	// RFC 2616, section 13.5.1: Remove a set of known hop-by-hop headers.
	// This behavior is superseded by the RFC 7230 Connection header, but
	// preserve it for backwards compatibility.
	delete(h, "Connection")          // "Del" a header.
	delete(h, "Keep-Alive")          // "Del" a header.
	delete(h, "Proxy-Connection")    // "Del" a header.
	delete(h, "Proxy-Authenticate")  // "Del" a header.
	delete(h, "Proxy-Authorization") // "Del" a header.
	delete(h, "Te")                  // "Del" a header.
	delete(h, "Trailer")             // "Del" a header.
	delete(h, "Transfer-Encoding")   // "Del" a header.
	delete(h, "Upgrade")             // "Del" a header.
}

// handleHopByHopHeaders removes or adds hop-by-hop headers.
// For incoming requests, call removeHopByHopHeaders before calling this function.
// Reference
//   - https://go.dev/src/net/http/httputil/reverseproxy.go
//   - https://datatracker.ietf.org/doc/rfc7230/
//   - https://datatracker.ietf.org/doc/rfc2616/
func handleHopByHopHeaders(inHeader, outHeader http.Header) {
	// The TE request header specifies the transfer encodings the user agent is willing to accept.
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/TE
	if v, ok := inHeader["Te"]; ok && httpguts.HeaderValuesContainsToken(v, "trailers") {
		outHeader["Te"] = []string{"trailers"} // "Set" a value.
	}
	if v, ok := inHeader["Connection"]; ok && httpguts.HeaderValuesContainsToken(v, "Upgrade") {
		outHeader["Connection"] = []string{"Upgrade"} // "Set" a value.
		up := inHeader["Upgrade"]
		outHeader["Upgrade"] = append(make([]string, 0, len(up)), up...) // "Set" values.
	}

	if _, ok := outHeader["User-Agent"]; !ok {
		// If the outbound request doesn't have a User-Agent header set,
		// don't send the default Go HTTP client User-Agent.
		outHeader["User-Agent"] = []string{""} // "Set" a header.
	}
}

// rewriteRequestURL rewrite URL to target one.
// Reference
//   - https://go.dev/src/net/http/httputil/reverseproxy.go
func rewriteRequestURL(req *url.URL, target *url.URL) {
	req.Scheme = target.Scheme
	req.Host = target.Host
	req.Path = target.Path
	req.RawPath = target.RawPath
	if target.RawQuery == "" || req.RawQuery == "" {
		req.RawQuery = target.RawQuery + req.RawQuery
	} else {
		req.RawQuery = target.RawQuery + "&" + req.RawQuery
	}
	if target.Fragment != "" {
		req.Fragment = target.Fragment
		req.RawFragment = target.RawFragment
	}
}

// setXForwardedHeaders sets the Forwarded, X-Forwarded-For, X-Forwarded-Host, and
// X-Forwarded-Proto headers of the outbound request.
// The given req and header must not be nil.
// Forwarded header is defined in RFC7239.
// Reference
//   - https://go.dev/src/net/http/httputil/reverseproxy.go
//   - https://datatracker.ietf.org/doc/rfc7239/
//   - https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Forwarded
//   - https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/X-Forwarded-For
//   - https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/X-Forwarded-Host
//   - https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/X-Forwarded-Proto
func setXForwardedHeaders(req *http.Request, header http.Header) {
	ip, port, err := net.SplitHostPort(req.RemoteAddr)
	if err == nil {
		if prior, ok := req.Header["X-Forwarded-For"]; ok {
			// We can't know the right order of each header values.
			// So the order of joined ip may not be correct
			// when multiple "X-Forwarded-For" headers were exist.
			ip = strings.Join(prior, ", ") + ", " + ip
		}
		header["X-Forwarded-For"] = []string{ip}    // "Set" a value.
		header["X-Forwarded-Port"] = []string{port} // "Set" a value.
	} else {
		delete(header, "X-Forwarded-For")  // "Del" header.
		delete(header, "X-Forwarded-Port") // "Del" header.
	}

	header["X-Forwarded-Host"] = []string{req.Host} // "Set" a value.

	if req.TLS == nil {
		header["X-Forwarded-Proto"] = []string{"http"} // "Set" a value.
	} else {
		header["X-Forwarded-Proto"] = []string{"https"} // "Set" a value.
	}
}
