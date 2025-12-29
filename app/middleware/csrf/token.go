// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package csrf

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"io"
	"mime"
	"net/http"

	"github.com/aileron-gateway/aileron-gateway/internal/hash"
	"github.com/tidwall/gjson"
)

// extractor extract a csrf token bounded to the given requests.
// extractor returns a token even if it is invalid.
type extractor interface {
	// extract returns csrf token by extracting from the given request.
	// It is caller's responsible to verify the returned token.
	extract(*http.Request) (string, error)
}

type csrfToken struct {
	secret   []byte
	seedSize int
	hashSize int
	hmac     hash.HMACFunc
}

func (h *csrfToken) new() (string, error) {
	seed := make([]byte, h.seedSize)
	if _, err := io.ReadFull(rand.Reader, seed); err != nil {
		return "", err
	}
	digest := h.hmac(seed, h.secret)
	b := append(seed, digest...)
	return hex.EncodeToString(b), nil
}

func (h *csrfToken) verify(token string) bool {
	b, err := hex.DecodeString(token)
	if err != nil {
		return false
	}
	if len(b) != h.seedSize+h.hashSize {
		return false
	}
	seed, digest := b[:h.seedSize], b[h.seedSize:]
	truth := h.hmac(seed, h.secret)
	return bytes.Equal(truth, digest)
}

// formExtractor returns CSRF token extracted from HTTP header.
// Empty string "" is returned when token not found.
type headerExtractor struct {
	headerName string
}

func (e *headerExtractor) extract(r *http.Request) (string, error) {
	return r.Header.Get(e.headerName), nil
}

// formExtractor returns CSRF token extracted from form body.
// Content-Type must be "application/x-www-form-urlencoded" .
// Empty string "" is returned when token not found.
type formExtractor struct {
	paramName string
}

func (e *formExtractor) extract(r *http.Request) (string, error) {
	mt, _, _ := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if mt != "application/x-www-form-urlencoded" {
		return "", errInvalidMIME
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return "", err
	}
	if r.GetBody == nil {
		r.GetBody = func() (io.ReadCloser, error) {
			return io.NopCloser(bytes.NewReader(body)), nil
		}
	}
	r.Body, _ = r.GetBody()
	val := r.PostFormValue(e.paramName)
	r.Body, _ = r.GetBody()
	return val, nil
}

// jsonExtractor returns CSRF token extracted from json body.
// Content-Type must be "application/json" .
// Empty string "" is returned when token not found.
type jsonExtractor struct {
	jsonPath string
}

func (e *jsonExtractor) extract(r *http.Request) (string, error) {
	mt, _, _ := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if mt != "application/json" {
		return "", errInvalidMIME
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return "", err
	}
	if r.GetBody == nil {
		r.GetBody = func() (io.ReadCloser, error) {
			return io.NopCloser(bytes.NewReader(body)), nil
		}
	}
	r.Body, _ = r.GetBody()
	return gjson.GetBytes(body, e.jsonPath).String(), nil
}
