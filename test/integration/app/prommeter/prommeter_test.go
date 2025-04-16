// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

//go:build integration

package prommeter_test

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/aileron-gateway/aileron-gateway/apis/kernel"
	"github.com/aileron-gateway/aileron-gateway/cmd/aileron/app"
	"github.com/aileron-gateway/aileron-gateway/core"
	"github.com/aileron-gateway/aileron-gateway/kernel/api"
	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
	"github.com/aileron-gateway/aileron-gateway/test/integration/common"
)

func TestMinimal_middleware(t *testing.T) {
	configs := []string{"./config-minimal.yaml"}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "app/v1",
		Kind:       "PrometheusMeter",
	}
	m, err := api.ReferTypedObject[core.Middleware](server, ref)
	testutil.DiffError(t, nil, nil, err)

	// Use PrometheusMeter Middleware on an actual handler
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})
	h := m.Middleware(handler)

	// Perform GET requests.
	for i := 0; i < 5; i++ {
		r := httptest.NewRequest(http.MethodGet, "http://metrics.com/get", nil)
		w := httptest.NewRecorder()
		h.ServeHTTP(w, r)
		b, _ := io.ReadAll(w.Result().Body)
		testutil.Diff(t, http.StatusOK, w.Result().StatusCode)
		testutil.Diff(t, "ok", string(b))
	}
	// Perform POST requests.
	for i := 0; i < 5; i++ {
		r := httptest.NewRequest(http.MethodPost, "http://metrics.com/post", bytes.NewBuffer([]byte("test body")))
		w := httptest.NewRecorder()
		h.ServeHTTP(w, r)
		b, _ := io.ReadAll(w.Result().Body)
		testutil.Diff(t, http.StatusOK, w.Result().StatusCode)
		testutil.Diff(t, "ok", string(b))
	}

	mh, err := api.ReferTypedObject[http.Handler](server, ref)
	testutil.DiffError(t, nil, nil, err)
	r := httptest.NewRequest(http.MethodGet, "http://metrics.com/metrics", nil)
	w := httptest.NewRecorder()
	mh.ServeHTTP(w, r)
	b, _ := io.ReadAll(w.Body)
	metrics := string(b)
	t.Log(metrics)
	testutil.Diff(t, true, strings.Contains(metrics, `http_requests_total{code="200",host="metrics.com",method="GET",path="/get"} 5`))
	testutil.Diff(t, true, strings.Contains(metrics, `http_requests_total{code="200",host="metrics.com",method="POST",path="/post"} 5`))
}

func TestMinimal_tripperware(t *testing.T) {
	configs := []string{"./config-minimal.yaml"}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "app/v1",
		Kind:       "PrometheusMeter",
	}
	m, err := api.ReferTypedObject[core.Tripperware](server, ref)
	testutil.DiffError(t, nil, nil, err)

	mockTransport := &mockRoundTripper{}
	rt := m.Tripperware(mockTransport)

	// Perform GET requests.
	for i := 0; i < 5; i++ {
		r := httptest.NewRequest(http.MethodGet, "http://metrics.com/get", nil)
		_, err := rt.RoundTrip(r)
		testutil.DiffError(t, nil, nil, err)
	}
	// Perform POST requests.
	for i := 0; i < 5; i++ {
		r := httptest.NewRequest(http.MethodPost, "http://metrics.com/post", bytes.NewBuffer([]byte("test body")))
		_, err := rt.RoundTrip(r)
		testutil.DiffError(t, nil, nil, err)
	}

	mh, err := api.ReferTypedObject[http.Handler](server, ref)
	testutil.DiffError(t, nil, nil, err)
	r := httptest.NewRequest(http.MethodGet, "http://metrics.com/metrics", nil)
	w := httptest.NewRecorder()
	mh.ServeHTTP(w, r)
	b, _ := io.ReadAll(w.Body)
	metrics := string(b)
	t.Log(metrics)
	testutil.Diff(t, true, strings.Contains(metrics, `http_client_requests_total{code="200",host="metrics.com",method="GET",path="/get"} 5`))
	testutil.Diff(t, true, strings.Contains(metrics, `http_client_requests_total{code="200",host="metrics.com",method="POST",path="/post"} 5`))
}
