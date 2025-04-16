// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

//go:build integration

package jaegertracer_test

import (
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/aileron-gateway/aileron-gateway/apis/kernel"
	"github.com/aileron-gateway/aileron-gateway/app/tracer/jaegertracer"
	"github.com/aileron-gateway/aileron-gateway/cmd/aileron/app"
	"github.com/aileron-gateway/aileron-gateway/core"
	"github.com/aileron-gateway/aileron-gateway/kernel/api"
	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
	"github.com/aileron-gateway/aileron-gateway/test/integration/common"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
)

func check(t *testing.T, m core.Middleware) {
	t.Helper()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Log(r.Header)
		testutil.Diff(t, true, len(r.Header["Uber-Trace-Id"]) > 0)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	h := m.Middleware(handler)

	r1 := httptest.NewRequest(http.MethodGet, "http://jaeger-test.com/test", nil)
	w1 := httptest.NewRecorder()
	h.ServeHTTP(w1, r1)
	testutil.Diff(t, http.StatusOK, w1.Result().StatusCode)
}

func TestMinimalWithMetadata(t *testing.T) {
	reporter := jaeger.NewInMemoryReporter()
	jaegertracer.TestOptions = []config.Option{config.Reporter(reporter)}
	defer func() {
		jaegertracer.TestOptions = nil
	}()

	configs := []string{"./config-minimal-with-metadata.yaml"}
	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "app/v1",
		Kind:       "JaegerTracer",
		Name:       "testName",
		Namespace:  "testNamespace",
	}
	m, err := api.ReferTypedObject[core.Middleware](server, ref)
	testutil.DiffError(t, nil, nil, err)

	check(t, m)
	t.Log(reporter.GetSpans())
	testutil.Diff(t, 1, reporter.SpansSubmitted())
}

func TestMinimalWithoutMetadata(t *testing.T) {
	reporter := jaeger.NewInMemoryReporter()
	jaegertracer.TestOptions = []config.Option{config.Reporter(reporter)}
	defer func() {
		jaegertracer.TestOptions = nil
	}()

	configs := []string{"./config-minimal-without-metadata.yaml"}
	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "app/v1",
		Kind:       "JaegerTracer",
	}
	m, err := api.ReferTypedObject[core.Middleware](server, ref)
	testutil.DiffError(t, nil, nil, err)

	check(t, m)
	t.Log(reporter.GetSpans())
	testutil.Diff(t, 1, reporter.SpansSubmitted())
}

func TestEmptyName(t *testing.T) {
	reporter := jaeger.NewInMemoryReporter()
	jaegertracer.TestOptions = []config.Option{config.Reporter(reporter)}
	defer func() {
		jaegertracer.TestOptions = nil
	}()

	configs := []string{"./config-empty-name.yaml"}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "app/v1",
		Kind:       "JaegerTracer",
		Name:       "",
		Namespace:  "testNamespace",
	}
	m, err := api.ReferTypedObject[core.Middleware](server, ref)
	testutil.DiffError(t, nil, nil, err)

	check(t, m)
	t.Log(reporter.GetSpans())
	testutil.Diff(t, 1, reporter.SpansSubmitted())
}

func TestEmptyNamespace(t *testing.T) {
	reporter := jaeger.NewInMemoryReporter()
	jaegertracer.TestOptions = []config.Option{config.Reporter(reporter)}
	defer func() {
		jaegertracer.TestOptions = nil
	}()

	configs := []string{"./config-empty-namespace.yaml"}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "app/v1",
		Kind:       "JaegerTracer",
		Name:       "testName",
		Namespace:  "",
	}
	m, err := api.ReferTypedObject[core.Middleware](server, ref)
	testutil.DiffError(t, nil, nil, err)

	check(t, m)
	t.Log(reporter.GetSpans())
	testutil.Diff(t, 1, reporter.SpansSubmitted())
}

func TestEmptyNameNamespace(t *testing.T) {
	reporter := jaeger.NewInMemoryReporter()
	jaegertracer.TestOptions = []config.Option{config.Reporter(reporter)}
	defer func() {
		jaegertracer.TestOptions = nil
	}()

	configs := []string{"./config-empty-name-namespace.yaml"}
	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "app/v1",
		Kind:       "JaegerTracer",
		Name:       "",
		Namespace:  "",
	}
	m, err := api.ReferTypedObject[core.Middleware](server, ref)
	testutil.DiffError(t, nil, nil, err)

	check(t, m)
	t.Log(reporter.GetSpans())
	testutil.Diff(t, 1, reporter.SpansSubmitted())
}

func TestEmptySpec(t *testing.T) {
	reporter := jaeger.NewInMemoryReporter()
	jaegertracer.TestOptions = []config.Option{config.Reporter(reporter)}
	defer func() {
		jaegertracer.TestOptions = nil
	}()

	configs := []string{"./config-empty-spec.yaml"}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "app/v1",
		Kind:       "JaegerTracer",
	}
	m, err := api.ReferTypedObject[core.Middleware](server, ref)
	testutil.DiffError(t, nil, nil, err)

	check(t, m)
	t.Log(reporter.GetSpans())
	testutil.Diff(t, 1, reporter.SpansSubmitted())
}

func TestInvalidSpec(t *testing.T) {
	configs := []string{"./config-invalid-spec.yaml"}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	errPattern := regexp.MustCompile(core.ErrPrefix + `failed to load configs. config-invalid-spec.yaml`)
	testutil.DiffError(t, app.ErrAppMainLoadConfigs, errPattern, err)
}
