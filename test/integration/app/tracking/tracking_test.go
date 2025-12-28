// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

//go:build integration

package tracking_test

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
	utilhttp "github.com/aileron-gateway/aileron-gateway/util/http"
	"github.com/aileron-projects/go/zx/zuid"
)

func TestProxyID(t *testing.T) {
	configs := []string{"./config-proxy-id.yaml"}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "app/v1",
		Kind:       "TrackingMiddleware",
	}
	m, err := api.ReferTypedObject[core.Middleware](server, ref)
	testutil.DiffError(t, nil, nil, err)

	r1 := httptest.NewRequest(http.MethodGet, "/", nil)
	w1 := httptest.NewRecorder()

	m.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqID := zuid.FromContext(r.Context(), "context")
		h := utilhttp.ProxyHeaderFromContext(r.Context())
		testutil.Diff(t, reqID, h.Get("X-Request-ID"))
		testutil.Diff(t, true, h.Get("X-Request-ID") != "")
		testutil.Diff(t, true, h.Get("X-Trace-ID") != "")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})).ServeHTTP(w1, r1)

	resp := w1.Result()
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	testutil.Diff(t, http.StatusOK, resp.StatusCode)
	testutil.Diff(t, "ok", string(body))
}

func TestExtractID(t *testing.T) {
	configs := []string{"./config-extract-id.yaml"}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "app/v1",
		Kind:       "TrackingMiddleware",
	}
	m, err := api.ReferTypedObject[core.Middleware](server, ref)
	testutil.DiffError(t, nil, nil, err)

	r1 := httptest.NewRequest(http.MethodGet, "/", nil)
	r1.Header.Set("X-Trace-Extract", "test-trace-id")
	w1 := httptest.NewRecorder()

	m.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h := utilhttp.ProxyHeaderFromContext(r.Context())
		testutil.Diff(t, "test-trace-id", h.Get("X-Trace-ID"))
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})).ServeHTTP(w1, r1)

	resp := w1.Result()
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	testutil.Diff(t, http.StatusOK, resp.StatusCode)
	testutil.Diff(t, "ok", string(body))
}
