// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

//go:build integration

package compression_test

import (
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/aileron-gateway/aileron-gateway/apis/kernel"
	"github.com/aileron-gateway/aileron-gateway/cmd/aileron/app"
	"github.com/aileron-gateway/aileron-gateway/core"
	"github.com/aileron-gateway/aileron-gateway/kernel/api"
	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
	"github.com/aileron-gateway/aileron-gateway/test/integration/common"
	"github.com/andybalholm/brotli"
)

func TestTargetMIMEs_gzip(t *testing.T) {
	configs := []string{"./config-target-mimes.yaml"}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "app/v1",
		Kind:       "CompressionMiddleware",
	}
	m, err := api.ReferTypedObject[core.Middleware](server, ref)
	testutil.DiffError(t, nil, nil, err)

	bodyContentType := "application/json"
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body := []byte(`{"message": "Test Body"}`)
		w.Header().Set("Content-Length", strconv.Itoa(len(body)))
		w.Header().Set("Content-Type", bodyContentType)
		w.WriteHeader(http.StatusOK)
		w.Write(body)
	})
	h := m.Middleware(handler)

	r1 := httptest.NewRequest(http.MethodGet, "/", nil)
	r1.Header.Set("Accept-Encoding", "deflate, gzip, zstd")

	// Expect compressed body.
	bodyContentType = "application/json"
	w1 := httptest.NewRecorder()
	h.ServeHTTP(w1, r1)
	testutil.Diff(t, w1.Code, http.StatusOK)
	testutil.Diff(t, "gzip", w1.Result().Header.Get("Content-Encoding"))
	reader, err := gzip.NewReader(w1.Result().Body)
	testutil.DiffError(t, nil, nil, err)
	defer reader.Close()
	body1, _ := io.ReadAll(reader)
	testutil.Diff(t, `{"message": "Test Body"}`, string(body1))

	// Expect uncompressed body.
	bodyContentType = "text/plain"
	w2 := httptest.NewRecorder()
	h.ServeHTTP(w2, r1)
	testutil.Diff(t, w2.Code, http.StatusOK)
	testutil.Diff(t, "", w2.Result().Header.Get("Content-Encoding"))
	body2, _ := io.ReadAll(w2.Result().Body)
	testutil.Diff(t, `{"message": "Test Body"}`, string(body2))
}

func TestTargetMIMEs_brotli(t *testing.T) {
	configs := []string{"./config-target-mimes.yaml"}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "app/v1",
		Kind:       "CompressionMiddleware",
	}
	m, err := api.ReferTypedObject[core.Middleware](server, ref)
	testutil.DiffError(t, nil, nil, err)

	bodyContentType := "application/json"
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body := []byte(`{"message": "Test Body"}`)
		w.Header().Set("Content-Length", strconv.Itoa(len(body)))
		w.Header().Set("Content-Type", bodyContentType)
		w.WriteHeader(http.StatusOK)
		w.Write(body)
	})
	h := m.Middleware(handler)

	r1 := httptest.NewRequest(http.MethodGet, "/", nil)
	r1.Header.Set("Accept-Encoding", "deflate, br, zstd")

	// Expect compressed body.
	bodyContentType = "application/json"
	w1 := httptest.NewRecorder()
	h.ServeHTTP(w1, r1)
	testutil.Diff(t, w1.Code, http.StatusOK)
	testutil.Diff(t, "br", w1.Result().Header.Get("Content-Encoding"))
	reader := brotli.NewReader(w1.Result().Body)
	body1, _ := io.ReadAll(reader)
	testutil.Diff(t, `{"message": "Test Body"}`, string(body1))

	// Expect uncompressed body.
	bodyContentType = "text/plain"
	w2 := httptest.NewRecorder()
	h.ServeHTTP(w2, r1)
	testutil.Diff(t, w2.Code, http.StatusOK)
	testutil.Diff(t, "", w2.Result().Header.Get("Content-Encoding"))
	body2, _ := io.ReadAll(w2.Result().Body)
	testutil.Diff(t, `{"message": "Test Body"}`, string(body2))
}

func TestMinimumSize(t *testing.T) {
	configs := []string{"./config-minimum-size.yaml"}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "app/v1",
		Kind:       "CompressionMiddleware",
	}
	m, err := api.ReferTypedObject[core.Middleware](server, ref)
	testutil.DiffError(t, nil, nil, err)

	body := []byte("")
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Length", strconv.Itoa(len(body)))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(body)
	})
	h := m.Middleware(handler)

	r1 := httptest.NewRequest(http.MethodGet, "/", nil)
	r1.Header.Set("Accept-Encoding", "gzip")

	// Expect compressed body.
	body = []byte(`{"message": "Test Body"}`)
	w1 := httptest.NewRecorder()
	h.ServeHTTP(w1, r1)
	testutil.Diff(t, w1.Code, http.StatusOK)
	testutil.Diff(t, "gzip", w1.Result().Header.Get("Content-Encoding"))
	reader, err := gzip.NewReader(w1.Result().Body)
	testutil.DiffError(t, nil, nil, err)
	defer reader.Close()
	body1, _ := io.ReadAll(reader)
	testutil.Diff(t, `{"message": "Test Body"}`, string(body1))

	// Expect uncompressed body.
	body = []byte(`{"message": "Test"}`)
	w2 := httptest.NewRecorder()
	h.ServeHTTP(w2, r1)
	testutil.Diff(t, w2.Code, http.StatusOK)
	testutil.Diff(t, "", w2.Result().Header.Get("Content-Encoding"))
	body2, _ := io.ReadAll(w2.Result().Body)
	testutil.Diff(t, `{"message": "Test"}`, string(body2))
}
