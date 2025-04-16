// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

//go:build integration

package header_test

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aileron-gateway/aileron-gateway/apis/kernel"
	"github.com/aileron-gateway/aileron-gateway/cmd/aileron/app"
	"github.com/aileron-gateway/aileron-gateway/core"
	"github.com/aileron-gateway/aileron-gateway/kernel/api"
	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
	"github.com/aileron-gateway/aileron-gateway/test/integration/common"
)

func TestAllowMIMEs(t *testing.T) {
	configs := []string{"./config-allow-mimes.yaml"}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "app/v1",
		Kind:       "HeaderPolicyMiddleware",
	}
	m, err := api.ReferTypedObject[core.Middleware](server, ref)
	testutil.DiffError(t, nil, nil, err)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`ok`))
	})
	h := m.Middleware(handler)

	r1 := httptest.NewRequest(http.MethodGet, "http://test.com", nil)
	r1.Header.Set("Content-Type", "application/json")
	w1 := httptest.NewRecorder()
	h.ServeHTTP(w1, r1)
	b1, _ := io.ReadAll(w1.Result().Body)
	testutil.Diff(t, `ok`, string(b1))
	testutil.Diff(t, http.StatusOK, w1.Result().StatusCode)

	r2 := httptest.NewRequest(http.MethodGet, "http://test.com", nil)
	r2.Header.Set("Content-Type", "application/xml")
	w2 := httptest.NewRecorder()
	h.ServeHTTP(w2, r2)
	b2, _ := io.ReadAll(w2.Result().Body)
	testutil.Diff(t, `{"status":415,"statusText":"Unsupported Media Type"}`, string(b2))
	testutil.Diff(t, http.StatusUnsupportedMediaType, w2.Result().StatusCode)
}

func TestMaxContentLength(t *testing.T) {
	configs := []string{"./config-max-content-length.yaml"}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "app/v1",
		Kind:       "HeaderPolicyMiddleware",
	}
	m, err := api.ReferTypedObject[core.Middleware](server, ref)
	testutil.DiffError(t, nil, nil, err)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`ok`))
	})
	h := m.Middleware(handler)

	r1 := httptest.NewRequest(http.MethodGet, "http://test.com", bytes.NewReader(bytes.Repeat([]byte("*"), 15)))
	w1 := httptest.NewRecorder()
	h.ServeHTTP(w1, r1)
	b1, _ := io.ReadAll(w1.Result().Body)
	testutil.Diff(t, `ok`, string(b1))
	testutil.Diff(t, http.StatusOK, w1.Result().StatusCode)

	r2 := httptest.NewRequest(http.MethodGet, "http://test.com", bytes.NewReader(bytes.Repeat([]byte("*"), 30)))
	w2 := httptest.NewRecorder()
	h.ServeHTTP(w2, r2)
	b2, _ := io.ReadAll(w2.Result().Body)
	testutil.Diff(t, `{"status":413,"statusText":"Request Entity Too Large"}`, string(b2))
	testutil.Diff(t, http.StatusRequestEntityTooLarge, w2.Result().StatusCode)

	r3 := httptest.NewRequest(http.MethodGet, "http://test.com", nil)
	r3.ContentLength = -1
	w3 := httptest.NewRecorder()
	h.ServeHTTP(w3, r3)
	b3, _ := io.ReadAll(w3.Result().Body)
	testutil.Diff(t, `{"status":411,"statusText":"Length Required"}`, string(b3))
	testutil.Diff(t, http.StatusLengthRequired, w3.Result().StatusCode)
}

func TestRequestAllows(t *testing.T) {
	configs := []string{"./config-request-allows.yaml"}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "app/v1",
		Kind:       "HeaderPolicyMiddleware",
	}
	m, err := api.ReferTypedObject[core.Middleware](server, ref)
	testutil.DiffError(t, nil, nil, err)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		testutil.Diff(t, "foo", r.Header.Get("X-Allows-Foo")) // Check if allowed.
		testutil.Diff(t, "bar", r.Header.Get("X-Allows-Bar")) // Check if allowed.
		testutil.Diff(t, "", r.Header.Get("X-Allows-Baz"))    // Check if allowed.
		w.Write([]byte(`ok`))
	})
	h := m.Middleware(handler)

	r1 := httptest.NewRequest(http.MethodGet, "http://test.com", nil)
	r1.Header.Set("X-Allows-Foo", "foo")
	r1.Header.Set("X-Allows-Bar", "bar")
	r1.Header.Set("X-Allows-Baz", "baz")

	w1 := httptest.NewRecorder()
	h.ServeHTTP(w1, r1)
	b1, _ := io.ReadAll(w1.Result().Body)
	testutil.Diff(t, `ok`, string(b1))
}

func TestRequest_removes(t *testing.T) {
	configs := []string{"./config-request.yaml"}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "app/v1",
		Kind:       "HeaderPolicyMiddleware",
	}
	m, err := api.ReferTypedObject[core.Middleware](server, ref)
	testutil.DiffError(t, nil, nil, err)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		testutil.Diff(t, "", r.Header.Get("X-Removes-Foo"))    // Check if removed.
		testutil.Diff(t, "", r.Header.Get("X-Removes-Bar"))    // Check if removed.
		testutil.Diff(t, "baz", r.Header.Get("X-Removes-Baz")) // Check if removed.
		w.Write([]byte(`ok`))
	})
	h := m.Middleware(handler)

	r1 := httptest.NewRequest(http.MethodGet, "http://test.com", nil)
	r1.Header.Set("X-Removes-Foo", "foo")
	r1.Header.Set("X-Removes-Bar", "bar")
	r1.Header.Set("X-Removes-Baz", "baz")

	w1 := httptest.NewRecorder()
	h.ServeHTTP(w1, r1)
	b1, _ := io.ReadAll(w1.Result().Body)
	testutil.Diff(t, `ok`, string(b1))
}

func TestRequest_add(t *testing.T) {
	configs := []string{"./config-request.yaml"}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "app/v1",
		Kind:       "HeaderPolicyMiddleware",
	}
	m, err := api.ReferTypedObject[core.Middleware](server, ref)
	testutil.DiffError(t, nil, nil, err)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		testutil.Diff(t, []string{"FOO", "foo"}, r.Header["X-Add-Foo"]) // Check if added.
		testutil.Diff(t, []string{"bar"}, r.Header["X-Add-Bar"])        // Check if added.
		w.Write([]byte(`ok`))
	})
	h := m.Middleware(handler)

	r1 := httptest.NewRequest(http.MethodGet, "http://test.com", nil)
	r1.Header.Add("X-Add-Foo", "FOO")

	w1 := httptest.NewRecorder()
	h.ServeHTTP(w1, r1)
	b1, _ := io.ReadAll(w1.Result().Body)
	testutil.Diff(t, `ok`, string(b1))
}

func TestRequest_set(t *testing.T) {
	configs := []string{"./config-request.yaml"}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "app/v1",
		Kind:       "HeaderPolicyMiddleware",
	}
	m, err := api.ReferTypedObject[core.Middleware](server, ref)
	testutil.DiffError(t, nil, nil, err)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		testutil.Diff(t, []string{"foo"}, r.Header["X-Set-Foo"]) // Check if set.
		testutil.Diff(t, []string{"bar"}, r.Header["X-Set-Bar"]) // Check if set.
		w.Write([]byte(`ok`))
	})
	h := m.Middleware(handler)

	r1 := httptest.NewRequest(http.MethodGet, "http://test.com", nil)
	r1.Header.Add("X-Set-Foo", "FOO")

	w1 := httptest.NewRecorder()
	h.ServeHTTP(w1, r1)
	b1, _ := io.ReadAll(w1.Result().Body)
	testutil.Diff(t, `ok`, string(b1))
}

func TestRequest_rewrites(t *testing.T) {
	configs := []string{"./config-request.yaml"}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "app/v1",
		Kind:       "HeaderPolicyMiddleware",
	}
	m, err := api.ReferTypedObject[core.Middleware](server, ref)
	testutil.DiffError(t, nil, nil, err)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		testutil.Diff(t, []string{"fOo", "fOo"}, r.Header["X-Rewrites-Foo"]) // Check if rewritten.
		testutil.Diff(t, []string{"bAr"}, r.Header["X-Rewrites-Bar"])        // Check if rewritten.
		w.Write([]byte(`ok`))
	})
	h := m.Middleware(handler)

	r1 := httptest.NewRequest(http.MethodGet, "http://test.com", nil)
	r1.Header.Add("X-Rewrites-Foo", "FOO")
	r1.Header.Add("X-Rewrites-Foo", "foo")
	r1.Header.Add("X-Rewrites-Bar", "bar")

	w1 := httptest.NewRecorder()
	h.ServeHTTP(w1, r1)
	b1, _ := io.ReadAll(w1.Result().Body)
	testutil.Diff(t, `ok`, string(b1))
}

func TestResponseAllows(t *testing.T) {
	configs := []string{"./config-response-allows.yaml"}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "app/v1",
		Kind:       "HeaderPolicyMiddleware",
	}
	m, err := api.ReferTypedObject[core.Middleware](server, ref)
	testutil.DiffError(t, nil, nil, err)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Allows-Foo", "foo")
		w.Header().Set("X-Allows-Bar", "bar")
		w.Header().Set("X-Allows-Baz", "baz")
		w.Write([]byte(`ok`))
	})
	h := m.Middleware(handler)

	r1 := httptest.NewRequest(http.MethodGet, "http://test.com", nil)
	w1 := httptest.NewRecorder()
	h.ServeHTTP(w1, r1)
	testutil.Diff(t, "foo", w1.Header().Get("X-Allows-Foo")) // Check if allowed.
	testutil.Diff(t, "bar", w1.Header().Get("X-Allows-Bar")) // Check if allowed.
	testutil.Diff(t, "", w1.Header().Get("X-Allows-Baz"))    // Check if allowed.
	b1, _ := io.ReadAll(w1.Result().Body)
	testutil.Diff(t, `ok`, string(b1))
}

func TestResponse_removes(t *testing.T) {
	configs := []string{"./config-response.yaml"}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "app/v1",
		Kind:       "HeaderPolicyMiddleware",
	}
	m, err := api.ReferTypedObject[core.Middleware](server, ref)
	testutil.DiffError(t, nil, nil, err)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Removes-Foo", "foo")
		w.Header().Set("X-Removes-Bar", "bar")
		w.Header().Set("X-Removes-Baz", "baz")
		w.Write([]byte(`ok`))
	})
	h := m.Middleware(handler)

	r1 := httptest.NewRequest(http.MethodGet, "http://test.com", nil)
	w1 := httptest.NewRecorder()
	h.ServeHTTP(w1, r1)
	testutil.Diff(t, "", w1.Header().Get("X-Removes-Foo"))    // Check if removed.
	testutil.Diff(t, "", w1.Header().Get("X-Removes-Bar"))    // Check if removed.
	testutil.Diff(t, "baz", w1.Header().Get("X-Removes-Baz")) // Check if removed.
	b1, _ := io.ReadAll(w1.Result().Body)
	testutil.Diff(t, `ok`, string(b1))
}

func TestResponse_add(t *testing.T) {
	configs := []string{"./config-response.yaml"}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "app/v1",
		Kind:       "HeaderPolicyMiddleware",
	}
	m, err := api.ReferTypedObject[core.Middleware](server, ref)
	testutil.DiffError(t, nil, nil, err)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("X-Add-Foo", "FOO")
		w.Write([]byte(`ok`))
	})
	h := m.Middleware(handler)

	r1 := httptest.NewRequest(http.MethodGet, "http://test.com", nil)
	w1 := httptest.NewRecorder()
	h.ServeHTTP(w1, r1)
	testutil.Diff(t, []string{"FOO", "foo"}, w1.Header()["X-Add-Foo"]) // Check if added.
	testutil.Diff(t, []string{"bar"}, w1.Header()["X-Add-Bar"])        // Check if added.
	b1, _ := io.ReadAll(w1.Result().Body)
	testutil.Diff(t, `ok`, string(b1))
}

func TestResponse_set(t *testing.T) {
	configs := []string{"./config-response.yaml"}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "app/v1",
		Kind:       "HeaderPolicyMiddleware",
	}
	m, err := api.ReferTypedObject[core.Middleware](server, ref)
	testutil.DiffError(t, nil, nil, err)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("X-Set-Foo", "FOO")
		w.Write([]byte(`ok`))
	})
	h := m.Middleware(handler)

	r1 := httptest.NewRequest(http.MethodGet, "http://test.com", nil)
	w1 := httptest.NewRecorder()
	h.ServeHTTP(w1, r1)
	testutil.Diff(t, []string{"foo"}, w1.Header()["X-Set-Foo"]) // Check if set.
	testutil.Diff(t, []string{"bar"}, w1.Header()["X-Set-Bar"]) // Check if set.
	b1, _ := io.ReadAll(w1.Result().Body)
	testutil.Diff(t, `ok`, string(b1))
}

func TestResponse_rewrites(t *testing.T) {
	configs := []string{"./config-response.yaml"}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "app/v1",
		Kind:       "HeaderPolicyMiddleware",
	}
	m, err := api.ReferTypedObject[core.Middleware](server, ref)
	testutil.DiffError(t, nil, nil, err)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("X-Rewrites-Foo", "FOO")
		w.Header().Add("X-Rewrites-Foo", "foo")
		w.Header().Add("X-Rewrites-Bar", "bar")
		w.Write([]byte(`ok`))
	})
	h := m.Middleware(handler)

	r1 := httptest.NewRequest(http.MethodGet, "http://test.com", nil)
	w1 := httptest.NewRecorder()
	h.ServeHTTP(w1, r1)
	testutil.Diff(t, []string{"fOo", "fOo"}, w1.Header()["X-Rewrites-Foo"]) // Check if rewritten.
	testutil.Diff(t, []string{"bAr"}, w1.Header()["X-Rewrites-Bar"])        // Check if rewritten.
	b1, _ := io.ReadAll(w1.Result().Body)
	testutil.Diff(t, `ok`, string(b1))
}
