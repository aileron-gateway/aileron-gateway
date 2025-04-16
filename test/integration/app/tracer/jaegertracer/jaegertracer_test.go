// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package jaegertracer_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aileron-gateway/aileron-gateway/apis/kernel"
	"github.com/aileron-gateway/aileron-gateway/app/tracer/jaegertracer"
	"github.com/aileron-gateway/aileron-gateway/cmd/aileron/app"
	"github.com/aileron-gateway/aileron-gateway/core"
	"github.com/aileron-gateway/aileron-gateway/kernel/api"
	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
	"github.com/aileron-gateway/aileron-gateway/test/integration/common"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
)

func TestHeader(t *testing.T) {
	reporter := jaeger.NewInMemoryReporter()
	jaegertracer.TestOptions = []config.Option{config.Reporter(reporter)}
	defer func() {
		jaegertracer.TestOptions = nil
	}()

	configs := []string{"./config-header.yaml"}
	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "app/v1",
		Kind:       "JaegerTracer",
	}
	m, err := api.ReferTypedObject[core.Middleware](server, ref)
	testutil.DiffError(t, nil, nil, err)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Log(r.Header)
		testutil.Diff(t, 1, len(r.Header["Uber-Trace-Id"]))
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})
	h := m.Middleware(handler)

	r1 := httptest.NewRequest(http.MethodGet, "http://jaeger-test.com/test", nil)
	r1.Header.Set("X-Custom-Header", "FooBar")
	w1 := httptest.NewRecorder()
	h.ServeHTTP(w1, r1)
	testutil.Diff(t, http.StatusOK, w1.Result().StatusCode)
	t.Log(reporter.GetSpans())
	testutil.Diff(t, 1, reporter.SpansSubmitted())
	span := reporter.GetSpans()[0].(*jaeger.Span)
	testutil.Diff(t, []string{"FooBar"}, span.Tags()["http.header.x-custom-header"])
}

func TestDisabled(t *testing.T) {
	reporter := jaeger.NewInMemoryReporter()
	jaegertracer.TestOptions = []config.Option{config.Reporter(reporter)}
	defer func() {
		jaegertracer.TestOptions = nil
	}()

	configs := []string{"./config-disabled.yaml"}
	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "app/v1",
		Kind:       "JaegerTracer",
	}
	m, err := api.ReferTypedObject[core.Middleware](server, ref)
	testutil.DiffError(t, nil, nil, err)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Log(r.Header)
		testutil.Diff(t, 0, len(r.Header["Uber-Trace-Id"]))
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})
	h := m.Middleware(handler)

	r1 := httptest.NewRequest(http.MethodGet, "http://jaeger-test.com/test", nil)
	w1 := httptest.NewRecorder()
	h.ServeHTTP(w1, r1)
	testutil.Diff(t, http.StatusOK, w1.Result().StatusCode)
	t.Log(reporter.GetSpans())
	testutil.Diff(t, 0, reporter.SpansSubmitted())
}

func TestMinimal_middleware(t *testing.T) {
	reporter := jaeger.NewInMemoryReporter()
	jaegertracer.TestOptions = []config.Option{config.Reporter(reporter)}
	defer func() {
		jaegertracer.TestOptions = nil
	}()

	configs := []string{"./config-minimal.yaml"}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "app/v1",
		Kind:       "JaegerTracer",
	}
	m, err := api.ReferTypedObject[core.Middleware](server, ref)
	testutil.DiffError(t, nil, nil, err)

	h := m.Middleware(m.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Log(r.Header)
		testutil.Diff(t, 1, len(r.Header["Uber-Trace-Id"]))
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})))

	r1 := httptest.NewRequest(http.MethodGet, "http://jaeger-test.com/test", nil)
	w1 := httptest.NewRecorder()
	h.ServeHTTP(w1, r1)
	testutil.Diff(t, http.StatusOK, w1.Result().StatusCode)
	t.Log(reporter.GetSpans())
	testutil.Diff(t, 2, reporter.SpansSubmitted())
	span1 := reporter.GetSpans()[0].(*jaeger.Span)
	span2 := reporter.GetSpans()[1].(*jaeger.Span)
	testutil.Diff(t, jaeger.SpanID(0), span2.SpanContext().ParentID())
	testutil.Diff(t, span1.SpanContext().ParentID(), span2.SpanContext().SpanID())
}

// mockRoundTripper is a mock implementation of http.RoundTripper
type mockRoundTripper struct {
	err error
}

func (m *mockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	if m.err != nil {
		return nil, m.err
	}
	// Simulate a successful HTTP response
	return &http.Response{
		StatusCode: http.StatusOK,
		Request:    req,
	}, nil
}

func TestMinimal_tripperware(t *testing.T) {
	reporter := jaeger.NewInMemoryReporter()
	jaegertracer.TestOptions = []config.Option{config.Reporter(reporter)}
	defer func() {
		jaegertracer.TestOptions = nil
	}()

	configs := []string{"./config-minimal.yaml"}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "app/v1",
		Kind:       "JaegerTracer",
	}
	m, err := api.ReferTypedObject[core.Tripperware](server, ref)
	testutil.DiffError(t, nil, nil, err)

	mockTransport := &mockRoundTripper{}
	tripperware := m.Tripperware(m.Tripperware(mockTransport))

	r1, _ := http.NewRequest(http.MethodGet, "http://tracer.com/integration", nil)
	w1, err := tripperware.RoundTrip(r1)
	testutil.DiffError(t, nil, nil, err)

	testutil.Diff(t, http.StatusOK, w1.StatusCode)
	t.Log(reporter.GetSpans())
	testutil.Diff(t, 2, reporter.SpansSubmitted())
	span1 := reporter.GetSpans()[0].(*jaeger.Span)
	span2 := reporter.GetSpans()[1].(*jaeger.Span)
	testutil.Diff(t, jaeger.SpanID(0), span2.SpanContext().ParentID())
	testutil.Diff(t, span1.SpanContext().ParentID(), span2.SpanContext().SpanID())
}

func TestMinimal_tripperwareError(t *testing.T) {
	reporter := jaeger.NewInMemoryReporter()
	jaegertracer.TestOptions = []config.Option{config.Reporter(reporter)}
	defer func() {
		jaegertracer.TestOptions = nil
	}()

	configs := []string{"./config-minimal.yaml"}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "app/v1",
		Kind:       "JaegerTracer",
	}
	m, err := api.ReferTypedObject[core.Tripperware](server, ref)
	testutil.DiffError(t, nil, nil, err)

	mockTransport := &mockRoundTripper{
		err: io.ErrUnexpectedEOF,
	}
	tripperware := m.Tripperware(m.Tripperware(mockTransport))

	r1, _ := http.NewRequest(http.MethodGet, "http://tracer.com/integration", nil)
	w1, err := tripperware.RoundTrip(r1)
	testutil.Diff(t, io.ErrUnexpectedEOF, err, cmpopts.EquateErrors())
	testutil.Diff(t, (*http.Response)(nil), w1)

	t.Log(reporter.GetSpans())
	testutil.Diff(t, 2, reporter.SpansSubmitted())
	span1 := reporter.GetSpans()[0].(*jaeger.Span)
	span2 := reporter.GetSpans()[1].(*jaeger.Span)
	testutil.Diff(t, jaeger.SpanID(0), span2.SpanContext().ParentID())
	testutil.Diff(t, span1.SpanContext().ParentID(), span2.SpanContext().SpanID())
}
