// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

//go:build integration

package compression_test

import (
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"regexp"
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

const (
	flagNoCompression = 1 << iota
	flagGzip
	flagBrotli
)

func check(t *testing.T, m core.Middleware, flag uint) {
	t.Helper()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body := []byte(`{"message": "Test Body"}`)
		w.Header().Set("Content-Length", strconv.Itoa(len(body)))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(body)
	})

	h := m.Middleware(handler)

	if flagNoCompression&flag > 0 {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Accept-Encoding", "br, gzip")
		resp := httptest.NewRecorder()
		h.ServeHTTP(resp, req)
		testutil.Diff(t, "", resp.Result().Header.Get("Content-Encoding"))
		body, err := io.ReadAll(resp.Result().Body)
		testutil.DiffError(t, nil, nil, err)
		testutil.Diff(t, `{"message": "Test Body"}`, string(body))
	}
	if flagGzip&flag > 0 {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Accept-Encoding", "gzip")
		resp := httptest.NewRecorder()
		h.ServeHTTP(resp, req)
		testutil.Diff(t, "gzip", resp.Result().Header.Get("Content-Encoding"))
		reader, _ := gzip.NewReader(resp.Result().Body)
		defer reader.Close()
		body, err := io.ReadAll(reader)
		testutil.DiffError(t, nil, nil, err)
		testutil.Diff(t, `{"message": "Test Body"}`, string(body))
	}
	if flagBrotli&flag > 0 {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Accept-Encoding", "br")
		resp := httptest.NewRecorder()
		h.ServeHTTP(resp, req)
		testutil.Diff(t, "br", resp.Result().Header.Get("Content-Encoding"))
		reader := brotli.NewReader(resp.Result().Body)
		body, err := io.ReadAll(reader)
		testutil.DiffError(t, nil, nil, err)
		testutil.Diff(t, `{"message": "Test Body"}`, string(body))
	}
}

func TestMinimalWithoutMetadata(t *testing.T) {
	configs := []string{"./config-minimal-without-metadata.yaml"}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "app/v1",
		Kind:       "CompressionMiddleware",
		Name:       "default",
		Namespace:  "",
	}
	m, err := api.ReferTypedObject[core.Middleware](server, ref)
	testutil.DiffError(t, nil, nil, err)

	check(t, m, flagNoCompression)
}

func TestMinimalWithMetadata(t *testing.T) {
	configs := []string{"./config-minimal-with-metadata.yaml"}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "app/v1",
		Kind:       "CompressionMiddleware",
		Name:       "testName",
		Namespace:  "testNamespace",
	}
	m, err := api.ReferTypedObject[core.Middleware](server, ref)
	testutil.DiffError(t, nil, nil, err)

	check(t, m, flagNoCompression)
}

func TestEmptyNameNamespace(t *testing.T) {
	configs := []string{"./config-empty-name-namespace.yaml"}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "app/v1",
		Kind:       "CompressionMiddleware",
		Name:       "default",
		Namespace:  "",
	}
	m, err := api.ReferTypedObject[core.Middleware](server, ref)
	testutil.DiffError(t, nil, nil, err)

	check(t, m, flagBrotli)
}

func TestEmptyName(t *testing.T) {
	configs := []string{"./config-empty-name.yaml"}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "app/v1",
		Kind:       "CompressionMiddleware",
		Name:       "default",
		Namespace:  "testNamespace",
	}
	m, err := api.ReferTypedObject[core.Middleware](server, ref)
	testutil.DiffError(t, nil, nil, err)

	check(t, m, flagBrotli)
}

func TestEmptyNamespace(t *testing.T) {
	configs := []string{"./config-empty-namespace.yaml"}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "app/v1",
		Kind:       "CompressionMiddleware",
		Name:       "testName",
		Namespace:  "",
	}
	m, err := api.ReferTypedObject[core.Middleware](server, ref)
	testutil.DiffError(t, nil, nil, err)

	check(t, m, flagBrotli)
}

func TestEmptySpec(t *testing.T) {
	configs := []string{"./config-empty-spec.yaml"}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "app/v1",
		Kind:       "CompressionMiddleware",
		Name:       "default",
		Namespace:  "",
	}
	m, err := api.ReferTypedObject[core.Middleware](server, ref)
	testutil.DiffError(t, nil, nil, err)

	check(t, m, flagNoCompression)
}

func TestInvalidSpec(t *testing.T) {
	configs := []string{"./config-invalid-spec.yaml"}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	errPattern := regexp.MustCompile(core.ErrPrefix + `failed to load configs. config-invalid-spec.yaml`)
	testutil.DiffError(t, app.ErrAppMainLoadConfigs, errPattern, err)
}
