// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

//go:build integration

package prommeter_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"

	"github.com/aileron-gateway/aileron-gateway/apis/kernel"
	"github.com/aileron-gateway/aileron-gateway/cmd/aileron/app"
	"github.com/aileron-gateway/aileron-gateway/core"
	"github.com/aileron-gateway/aileron-gateway/internal/testutil"
	"github.com/aileron-gateway/aileron-gateway/kernel/api"
	"github.com/aileron-gateway/aileron-gateway/test/integration/common"
)

type meter interface {
	http.Handler
	core.Middleware
	core.Tripperware
}

// mockRoundTripper is a mock implementation of http.RoundTripper
type mockRoundTripper struct {
	err error
}

func (m *mockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	if m.err != nil {
		return nil, m.err
	}
	return &http.Response{
		StatusCode: http.StatusOK,
		Request:    req,
	}, nil
}

func check(t *testing.T, m meter) {
	t.Helper()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})
	h := m.Middleware(handler)

	r1 := httptest.NewRequest(http.MethodGet, "http://metrics.com/middleware", nil)
	w1 := httptest.NewRecorder()
	h.ServeHTTP(w1, r1)
	b1, _ := io.ReadAll(w1.Result().Body)
	testutil.Diff(t, http.StatusOK, w1.Result().StatusCode)
	testutil.Diff(t, "ok", string(b1))

	mockTransport := &mockRoundTripper{}
	rt := m.Tripperware(mockTransport)
	r2 := httptest.NewRequest(http.MethodGet, "http://metrics.com/tripperware", nil)
	_, err := rt.RoundTrip(r2)
	testutil.DiffError(t, nil, nil, err)

	r3 := httptest.NewRequest(http.MethodGet, "http://metrics.com/metrics", nil)
	w3 := httptest.NewRecorder()
	m.ServeHTTP(w3, r3)
	b3, _ := io.ReadAll(w3.Body)

	metrics := string(b3)
	t.Log(metrics)
	testutil.Diff(t, true, strings.Contains(metrics, `http_requests_total{code="200",host="metrics.com",method="GET",path="/middleware"} 1`))
	testutil.Diff(t, true, strings.Contains(metrics, `http_client_requests_total{code="200",host="metrics.com",method="GET",path="/tripperware"} 1`))
}

func TestEmptyNameNamespace(t *testing.T) {
	configs := []string{"./config-empty-name-namespace.yaml"}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "app/v1",
		Kind:       "PrometheusMeter",
		Name:       "",
		Namespace:  "",
	}
	m, err := api.ReferTypedObject[meter](server, ref)
	testutil.DiffError(t, nil, nil, err)

	check(t, m)
}

func TestEmptyName(t *testing.T) {
	configs := []string{"./config-empty-name.yaml"}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "app/v1",
		Kind:       "PrometheusMeter",
		Name:       "",
		Namespace:  "testNamespace",
	}
	m, err := api.ReferTypedObject[meter](server, ref)
	testutil.DiffError(t, nil, nil, err)

	check(t, m)
}

func TestEmptyNamespace(t *testing.T) {
	configs := []string{"./config-empty-namespace.yaml"}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "app/v1",
		Kind:       "PrometheusMeter",
		Name:       "testName",
		Namespace:  "",
	}
	m, err := api.ReferTypedObject[meter](server, ref)
	testutil.DiffError(t, nil, nil, err)

	check(t, m)
}

func TestEmptySpec(t *testing.T) {
	configs := []string{"./config-empty-spec.yaml"}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "app/v1",
		Kind:       "PrometheusMeter",
		Name:       "testName",
		Namespace:  "testNamespace",
	}
	m, err := api.ReferTypedObject[meter](server, ref)
	testutil.DiffError(t, nil, nil, err)

	check(t, m)
}

func TestInvalidSpec(t *testing.T) {
	configs := []string{"./config-invalid-spec.yaml"}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	errPattern := regexp.MustCompile(core.ErrPrefix + `failed to load configs. config-invalid-spec.yaml`)
	testutil.DiffError(t, app.ErrAppMainLoadConfigs, errPattern, err)
}

func TestMinimalWithoutMetadata(t *testing.T) {
	configs := []string{"./config-minimal-without-metadata.yaml"}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "app/v1",
		Kind:       "PrometheusMeter",
	}
	m, err := api.ReferTypedObject[meter](server, ref)
	testutil.DiffError(t, nil, nil, err)

	check(t, m)
}

func TestMinimalWithMetadata(t *testing.T) {
	configs := []string{"./config-minimal-with-metadata.yaml"}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "app/v1",
		Kind:       "PrometheusMeter",
		Name:       "testName",
		Namespace:  "testNamespace",
	}
	m, err := api.ReferTypedObject[meter](server, ref)
	testutil.DiffError(t, nil, nil, err)

	check(t, m)
}
