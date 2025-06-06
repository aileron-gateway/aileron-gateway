// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package resilience

import (
	"encoding/binary"
	"hash/fnv"
	"net/http"
	"net/textproto"

	v1 "github.com/aileron-gateway/aileron-gateway/apis/core/v1"
)

// HTTPHasher calculate the hash value from a HTTP request.
type HTTPHasher interface {
	// Hash returns the hash value calculated from the given requests.
	// Hash functions and hash source information are vary
	// depending on the implementers.
	// The returned int is the int converted value of the hash value
	// calculated by binary.BigEndian.Uint32 function.
	// It means the returned int ranges from 0 to 65,535.
	// The returned bool will be false if the hasher
	// failed to calculate a hash from the request.
	Hash(r *http.Request) (int, bool)
}

// NewHTTPHashers returns a new instances of hashers.
func NewHTTPHashers(specs []*v1.HTTPHasherSpec) []HTTPHasher {
	hs := make([]HTTPHasher, 0, len(specs))
	for _, spec := range specs {
		if spec == nil {
			continue
		}
		h := NewHTTPHasher(spec)
		hs = append(hs, h)
	}
	return hs
}

func fnv1a32(b []byte) []byte {
	h := fnv.New32a()
	h.Write(b)
	return h.Sum(nil)
}

// NewHTTPHasher returns a new instance of hasher.
// nil is returned when nil was given as the spec.
func NewHTTPHasher(spec *v1.HTTPHasherSpec) HTTPHasher {
	if spec == nil {
		return nil
	}
	var h HTTPHasher
	switch spec.HasherType {
	case v1.HTTPHasherType_ClientAddr:
		h = &clientAddrHasher{}
	case v1.HTTPHasherType_Header:
		h = &headerHasher{
			name: textproto.CanonicalMIMEHeaderKey(spec.Key),
		}
	case v1.HTTPHasherType_MultiHeader:
		for i := range spec.Keys {
			spec.Keys[i] = textproto.CanonicalMIMEHeaderKey(spec.Keys[i])
		}
		h = &multiHeaderHasher{
			names: spec.Keys,
		}
	case v1.HTTPHasherType_Cookie:
		h = &cookieHasher{
			name: spec.Key,
		}
	case v1.HTTPHasherType_Query:
		h = &queryHasher{
			name: spec.Key,
		}
	case v1.HTTPHasherType_PathParam:
		h = &pathParamHasher{
			name: spec.Key,
		}
	default:
		h = &clientAddrHasher{}
	}
	return h
}

// clientAddrHasher calculate hash from the client's network address.
// Hash is calculated like hashFunc(r.RemoteAddr).
// This hashing always works.
type clientAddrHasher struct {
}

func (h *clientAddrHasher) Hash(r *http.Request) (int, bool) {
	v := fnv1a32([]byte(r.RemoteAddr))
	return int(binary.BigEndian.Uint32(v) >> 1), true // Shift 1 bit to make the value positive.
}

// headerHasher calculate hash from the header value.
// Hash is calculated like hashFunc(r.Header.Get("<name>")).
// When multiple header values were found, only 1 of them is used.
type headerHasher struct {
	// name is the header name.
	// If header value was not found, this hasher will fails.
	name string
}

func (h *headerHasher) Hash(r *http.Request) (int, bool) {
	v := r.Header[h.name]
	if len(v) == 0 {
		return -1, false
	}
	sum := fnv1a32([]byte(v[0]))
	return int(binary.BigEndian.Uint32(sum) >> 1), true // Shift 1 bit to make the value positive.
}

// multiHeaderHasher calculate hash from the multiple header values.
// Hash is calculated like hashFunc(join(r.Header.Get("<name>"))).
// When multiple header values were found for a header, only 1 of them is used.
// If all of the header values were empty, hashing will fail.
type multiHeaderHasher struct {
	// names is the list of header names.
	// If all header values were not found, this hasher will fails.
	names []string
}

func (h *multiHeaderHasher) Hash(r *http.Request) (int, bool) {
	v := ""
	for _, n := range h.names {
		vv := r.Header[n]
		if len(vv) > 0 {
			v += vv[0]
		}
	}
	if v == "" {
		return -1, false
	}
	sum := fnv1a32([]byte(v))
	return int(binary.BigEndian.Uint32(sum) >> 1), true // Shift 1 bit to make the value positive.
}

// cookieHasher calculate hash from the cookie value.
// Hash is calculated like hashFunc(r.Cookie("<name>")).
// When multiple cookie values were found, only 1 of them is used.
type cookieHasher struct {
	// name is the cookie name.
	// If cookie value was not found, this hasher will fails.
	name string
}

func (h *cookieHasher) Hash(r *http.Request) (int, bool) {
	ck, err := r.Cookie(h.name)
	if err != nil {
		return -1, false
	}
	v := ck.Value
	if v == "" {
		return -1, false
	}
	sum := fnv1a32([]byte(v))
	return int(binary.BigEndian.Uint32(sum) >> 1), true // Shift 1 bit to make the value positive.
}

// queryHasher calculate hash from the query parameter.
// Hash is calculated like hashFunc(r.URL.Query().Get("<name>")).
// When multiple query values were found, only 1 of them is used.
type queryHasher struct {
	// name is the query name.
	// If query value was not found, this hasher will fails.
	name string
}

func (h *queryHasher) Hash(r *http.Request) (int, bool) {
	v := r.URL.Query().Get(h.name)
	if v == "" {
		return -1, false
	}
	sum := fnv1a32([]byte(v))
	return int(binary.BigEndian.Uint32(sum) >> 1), true // Shift 1 bit to make the value positive.
}

// pathParamHasher calculate hash from the path parameter.
// Hash is calculated like hashFunc(r.PathValue("<name>")).
type pathParamHasher struct {
	// name is the path param name.
	// If param value was not found, this hasher will fails.
	name string
}

func (h *pathParamHasher) Hash(r *http.Request) (int, bool) {
	v := r.PathValue(h.name)
	if v == "" {
		return -1, false
	}
	sum := fnv1a32([]byte(v))
	return int(binary.BigEndian.Uint32(sum) >> 1), true // Shift 1 bit to make the value positive.
}
