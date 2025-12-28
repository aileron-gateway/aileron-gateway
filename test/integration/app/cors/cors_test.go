// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

//go:build integration

package cors_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aileron-gateway/aileron-gateway/apis/kernel"
	"github.com/aileron-gateway/aileron-gateway/cmd/aileron/app"
	"github.com/aileron-gateway/aileron-gateway/core"
	"github.com/aileron-gateway/aileron-gateway/internal/testutil"
	"github.com/aileron-gateway/aileron-gateway/kernel/api"
	"github.com/aileron-gateway/aileron-gateway/test/integration/common"
)

func TestHeaders(t *testing.T) {
	configs := []string{"./config-headers.yaml"}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "app/v1",
		Kind:       "CORSMiddleware",
	}
	m, err := api.ReferTypedObject[core.Middleware](server, ref)
	testutil.DiffError(t, nil, nil, err)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})
	h := m.Middleware(handler)

	// Preflight request.
	r1 := httptest.NewRequest(http.MethodOptions, "http://cors.com", nil)
	r1.Header.Set("Origin", "http://cors.com")
	r1.Header.Set("Access-Control-Request-Method", http.MethodGet)
	r1.Header.Set("Access-Control-Request-Headers", "Foo,Bar")
	w1 := httptest.NewRecorder()
	h.ServeHTTP(w1, r1)
	testutil.Diff(t, http.StatusOK, w1.Result().StatusCode)
	testutil.Diff(t, "http://cors.com", w1.Header().Get("Access-Control-Allow-Origin"))
	testutil.Diff(t, "true", w1.Header().Get("Access-Control-Allow-Credentials"))
	testutil.Diff(t, "X-Allowed-1,X-Allowed-2", w1.Header().Get("Access-Control-Allow-Headers"))
	testutil.Diff(t, "X-Exposed-1,X-Exposed-2", w1.Header().Get("Access-Control-Expose-Headers"))
	testutil.Diff(t, "3600", w1.Header().Get("Access-Control-Max-Age"))
	testutil.Diff(t, "", w1.Header().Get("Cross-Origin-Embedder-Policy"))
	testutil.Diff(t, "", w1.Header().Get("Cross-Origin-Opener-Policy"))
	testutil.Diff(t, "", w1.Header().Get("Cross-Origin-Resource-Policy"))

	// Actual request.
	r2 := httptest.NewRequest(http.MethodOptions, "http://cors.com", nil)
	r2.Header.Set("Origin", "http://cors.com")
	r2.Header.Set("Access-Control-Request-Headers", "Foo,Bar")
	w2 := httptest.NewRecorder()
	h.ServeHTTP(w2, r2)
	testutil.Diff(t, http.StatusOK, w2.Result().StatusCode)
	testutil.Diff(t, "http://cors.com", w2.Header().Get("Access-Control-Allow-Origin"))
	testutil.Diff(t, "true", w2.Header().Get("Access-Control-Allow-Credentials"))
	testutil.Diff(t, "X-Allowed-1,X-Allowed-2", w2.Header().Get("Access-Control-Allow-Headers"))
	testutil.Diff(t, "X-Exposed-1,X-Exposed-2", w2.Header().Get("Access-Control-Expose-Headers"))
	testutil.Diff(t, "", w2.Header().Get("Access-Control-Max-Age"))
	testutil.Diff(t, "require-corp", w2.Header().Get("Cross-Origin-Embedder-Policy"))
	testutil.Diff(t, "same-origin", w2.Header().Get("Cross-Origin-Opener-Policy"))
	testutil.Diff(t, "same-origin", w2.Header().Get("Cross-Origin-Resource-Policy"))
}

func TestAllowCredentials(t *testing.T) {
	configs := []string{"./config-allow-credentials.yaml"}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "app/v1",
		Kind:       "CORSMiddleware",
	}
	m, err := api.ReferTypedObject[core.Middleware](server, ref)
	testutil.DiffError(t, nil, nil, err)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})
	h := m.Middleware(handler)

	r1 := httptest.NewRequest(http.MethodGet, "http://allowed-origin.com/test", nil)
	r1.Header.Set("Origin", "http://allowed-origin.com")
	w1 := httptest.NewRecorder()
	h.ServeHTTP(w1, r1)
	testutil.Diff(t, http.StatusOK, w1.Result().StatusCode)
	testutil.Diff(t, "http://allowed-origin.com", w1.Header().Get("Access-Control-Allow-Origin"))
	testutil.Diff(t, "true", w1.Header().Get("Access-Control-Allow-Credentials"))

	r2 := httptest.NewRequest(http.MethodGet, "http://disallowed-origin.com/test", nil)
	r2.Header.Set("Origin", "http://disallowed-origin.com")
	w2 := httptest.NewRecorder()
	h.ServeHTTP(w2, r2)
	testutil.Diff(t, http.StatusOK, w2.Result().StatusCode)
	testutil.Diff(t, "*", w2.Header().Get("Access-Control-Allow-Origin"))
	testutil.Diff(t, "true", w2.Header().Get("Access-Control-Allow-Credentials"))
}

func TestAllowedOrigin(t *testing.T) {
	configs := []string{"./config-allowed-origin.yaml"}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "app/v1",
		Kind:       "CORSMiddleware",
	}
	m, err := api.ReferTypedObject[core.Middleware](server, ref)
	testutil.DiffError(t, nil, nil, err)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})
	h := m.Middleware(handler)

	// Preflight request from allowed origin.
	r1 := httptest.NewRequest(http.MethodOptions, "http://allowed-origin.com/test", nil)
	r1.Header.Set("Origin", "http://allowed-origin.com")
	r1.Header.Set("Access-Control-Request-Method", http.MethodGet)
	w1 := httptest.NewRecorder()
	h.ServeHTTP(w1, r1)
	testutil.Diff(t, http.StatusOK, w1.Result().StatusCode)
	testutil.Diff(t, "http://allowed-origin.com", w1.Header().Get("Access-Control-Allow-Origin"))
	b1, _ := io.ReadAll(w1.Result().Body)
	testutil.Diff(t, "", string(b1))

	// Preflight request from disallowed origin.
	r2 := httptest.NewRequest(http.MethodOptions, "http://disallowed-origin.com/test", nil)
	r2.Header.Set("Origin", "http://disallowed-origin.com")
	r2.Header.Set("Access-Control-Request-Method", http.MethodGet)
	w2 := httptest.NewRecorder()
	h.ServeHTTP(w2, r2)
	testutil.Diff(t, http.StatusForbidden, w2.Result().StatusCode)
	testutil.Diff(t, "", w2.Header().Get("Access-Control-Allow-Origin"))
	b2, _ := io.ReadAll(w2.Result().Body)
	testutil.Diff(t, "", string(b2))

	// Actual request from allowed origin with allowed method.
	r3 := httptest.NewRequest(http.MethodGet, "http://allowed-origin.com/test", nil)
	r3.Header.Set("Origin", "http://allowed-origin.com")
	w3 := httptest.NewRecorder()
	h.ServeHTTP(w3, r3)
	testutil.Diff(t, http.StatusOK, w3.Result().StatusCode)
	testutil.Diff(t, "http://allowed-origin.com", w3.Header().Get("Access-Control-Allow-Origin"))
	b3, _ := io.ReadAll(w3.Result().Body)
	testutil.Diff(t, "ok", string(b3))

	// Actual request from allowed origin with disallowed method.
	r4 := httptest.NewRequest(http.MethodDelete, "http://allowed-origin.com/test", nil)
	r4.Header.Set("Origin", "http://allowed-origin.com")
	w4 := httptest.NewRecorder()
	h.ServeHTTP(w4, r4)
	testutil.Diff(t, http.StatusForbidden, w4.Result().StatusCode)
	testutil.Diff(t, "http://allowed-origin.com", w4.Header().Get("Access-Control-Allow-Origin"))
	b4, _ := io.ReadAll(w4.Result().Body)
	testutil.Diff(t, `{"status":403,"statusText":"Forbidden"}`, string(b4))
}

func TestWildcardOrigin(t *testing.T) {
	configs := []string{"./config-wildcard-origin.yaml"}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "app/v1",
		Kind:       "CORSMiddleware",
	}
	m, err := api.ReferTypedObject[core.Middleware](server, ref)
	testutil.DiffError(t, nil, nil, err)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})
	h := m.Middleware(handler)

	// Preflight request from allowed origin.
	r1 := httptest.NewRequest(http.MethodOptions, "http://allowed-origin.com/test", nil)
	r1.Header.Set("Origin", "http://allowed-origin.com")
	r1.Header.Set("Access-Control-Request-Method", http.MethodGet)
	w1 := httptest.NewRecorder()
	h.ServeHTTP(w1, r1)
	testutil.Diff(t, http.StatusOK, w1.Result().StatusCode)
	testutil.Diff(t, "http://allowed-origin.com", w1.Header().Get("Access-Control-Allow-Origin"))
	b1, _ := io.ReadAll(w1.Result().Body)
	testutil.Diff(t, "", string(b1))

	// Preflight request from other origin.
	r2 := httptest.NewRequest(http.MethodOptions, "http://other-origin.com/test", nil)
	r2.Header.Set("Origin", "http://other-origin.com")
	r2.Header.Set("Access-Control-Request-Method", http.MethodGet)
	w2 := httptest.NewRecorder()
	h.ServeHTTP(w2, r2)
	testutil.Diff(t, http.StatusOK, w2.Result().StatusCode)
	testutil.Diff(t, "*", w2.Header().Get("Access-Control-Allow-Origin"))
	b2, _ := io.ReadAll(w2.Result().Body)
	testutil.Diff(t, "", string(b2))

	// Actual request from allowed origin with allowed method.
	r3 := httptest.NewRequest(http.MethodGet, "http://allowed-origin.com/test", nil)
	r3.Header.Set("Origin", "http://allowed-origin.com")
	w3 := httptest.NewRecorder()
	h.ServeHTTP(w3, r3)
	testutil.Diff(t, http.StatusOK, w3.Result().StatusCode)
	testutil.Diff(t, "http://allowed-origin.com", w3.Header().Get("Access-Control-Allow-Origin"))
	b3, _ := io.ReadAll(w3.Result().Body)
	testutil.Diff(t, "ok", string(b3))

	// Actual request from allowed origin with disallowed method.
	r4 := httptest.NewRequest(http.MethodDelete, "http://allowed-origin.com/test", nil)
	r4.Header.Set("Origin", "http://allowed-origin.com")
	w4 := httptest.NewRecorder()
	h.ServeHTTP(w4, r4)
	testutil.Diff(t, http.StatusForbidden, w4.Result().StatusCode)
	testutil.Diff(t, "http://allowed-origin.com", w4.Header().Get("Access-Control-Allow-Origin"))
	b4, _ := io.ReadAll(w4.Result().Body)
	testutil.Diff(t, `{"status":403,"statusText":"Forbidden"}`, string(b4))
}
