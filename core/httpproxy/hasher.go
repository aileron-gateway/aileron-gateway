// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package httpproxy

import (
	"net/http"

	v1 "github.com/aileron-gateway/aileron-gateway/apis/core/v1"
	"github.com/cespare/xxhash/v2"
)

// HTTPHasher calculates hash value from HTTP requests.
type HTTPHasher interface {
	// Hash returns the hash value calculated from the given requests.
	// The returned uint64 is the hashed value.
	Hash(r *http.Request) uint64
}

// newHTTPHasher returns a new instance of hasher.
func newHTTPHasher(spec *v1.HTTPHasherSpec) HTTPHasher {
	if spec == nil {
		return clientAddrHasher("")
	}
	switch spec.HashSource {
	case v1.HTTPHasherSpec_ClientAddr:
		return clientAddrHasher("")
	case v1.HTTPHasherSpec_Header:
		return headerHasher(http.CanonicalHeaderKey(spec.Key))
	case v1.HTTPHasherSpec_Cookie:
		return cookieHasher(spec.Key)
	case v1.HTTPHasherSpec_Query:
		return queryHasher(spec.Key)
	case v1.HTTPHasherSpec_PathParam:
		return pathParamHasher(spec.Key)
	default:
		return clientAddrHasher("")
	}
}

// clientAddrHasher calculates hash value from the client's network address.
type clientAddrHasher string

func (h clientAddrHasher) Hash(r *http.Request) uint64 {
	return xxhash.Sum64String(r.RemoteAddr)
}

// headerHasher calculates hash value from the header value.
type headerHasher string

func (h headerHasher) Hash(r *http.Request) uint64 {
	name := string(h)
	return xxhash.Sum64String(r.Header.Get(name))
}

// cookieHasher calculates hash value from the cookie value.
type cookieHasher string

func (h cookieHasher) Hash(r *http.Request) uint64 {
	name := string(h)
	ck, err := r.Cookie(name)
	if err != nil {
		return xxhash.Sum64String("")
	}
	return xxhash.Sum64String(ck.Value)
}

// queryHasher calculates hash value from the query parameter.
type queryHasher string

func (h queryHasher) Hash(r *http.Request) uint64 {
	name := string(h)
	return xxhash.Sum64String(r.URL.Query().Get(name))
}

// pathParamHasher calculates hash value from path parameter.
type pathParamHasher string

func (h pathParamHasher) Hash(r *http.Request) uint64 {
	name := string(h)
	return xxhash.Sum64String(r.PathValue(name))
}
