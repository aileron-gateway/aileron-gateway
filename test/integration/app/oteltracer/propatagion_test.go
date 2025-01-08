package oteltracer_test

import (
	"encoding/hex"
	"io"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/aileron-gateway/aileron-gateway/apis/kernel"
	"github.com/aileron-gateway/aileron-gateway/cmd/aileron/app"
	"github.com/aileron-gateway/aileron-gateway/core"
	"github.com/aileron-gateway/aileron-gateway/kernel/api"
	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
	"github.com/aileron-gateway/aileron-gateway/test/integration/common"
	"go.opencensus.io/trace"
	"go.opencensus.io/trace/propagation"
)

func TestPropagationTraceContext(t *testing.T) {
	server := common.NewAPI()
	err := app.LoadConfigFiles(server, []string{"./config-propagation-tracecontext.yaml"})
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "app/v1",
		Kind:       "OpenTelemetryTracer",
	}
	m, err := api.ReferTypedObject[core.Middleware](server, ref)
	testutil.DiffError(t, nil, nil, err)

	parentID := "00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-01"
	idPattern := regexp.MustCompile(`^00-[0-9a-f]{32}-[0-9a-f]{16}-01$`)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Log(r.Header)
		testutil.Diff(t, true, parentID != r.Header.Get("Traceparent"))
		testutil.Diff(t, true, idPattern.MatchString(r.Header.Get("Traceparent")))
		w.Write([]byte("ok"))
	})
	h := m.Middleware(handler)

	r1 := httptest.NewRequest(http.MethodGet, "http://localhost:8080/", nil)
	r1.Header.Add("Traceparent", parentID) // With parent span.
	w1 := httptest.NewRecorder()
	h.ServeHTTP(w1, r1)
	b1, _ := io.ReadAll(w1.Result().Body)
	testutil.Diff(t, http.StatusOK, w1.Result().StatusCode)
	testutil.Diff(t, "ok", string(b1))

	r2 := httptest.NewRequest(http.MethodGet, "http://localhost:8080/", nil)
	w2 := httptest.NewRecorder()
	h.ServeHTTP(w2, r2)
	b2, _ := io.ReadAll(w2.Result().Body)
	testutil.Diff(t, http.StatusOK, w2.Result().StatusCode)
	testutil.Diff(t, "ok", string(b2))
}

func TestPropagationB3(t *testing.T) {
	server := common.NewAPI()
	err := app.LoadConfigFiles(server, []string{"./config-propagation-b3.yaml"})
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "app/v1",
		Kind:       "OpenTelemetryTracer",
	}
	m, err := api.ReferTypedObject[core.Middleware](server, ref)
	testutil.DiffError(t, nil, nil, err)

	b3ID := "80f198ee56343ba864fe8b2a57d3eff7-e457b5a2e4d86bd1-1"
	idPattern := regexp.MustCompile(`^[0-9a-f]{32}-[0-9a-f]{16}-1$`)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Log(r.Header)
		testutil.Diff(t, true, b3ID != r.Header.Get("B3"))
		testutil.Diff(t, true, idPattern.MatchString(r.Header.Get("B3")))
		w.Write([]byte("ok"))
	})
	h := m.Middleware(handler)

	r1 := httptest.NewRequest(http.MethodGet, "http://localhost:8080/", nil)
	r1.Header.Add("B3", b3ID)
	w1 := httptest.NewRecorder()
	h.ServeHTTP(w1, r1)
	b1, _ := io.ReadAll(w1.Result().Body)
	testutil.Diff(t, http.StatusOK, w1.Result().StatusCode)
	testutil.Diff(t, "ok", string(b1))

	r2 := httptest.NewRequest(http.MethodGet, "http://localhost:8080/", nil)
	w2 := httptest.NewRecorder()
	h.ServeHTTP(w2, r2)
	b2, _ := io.ReadAll(w2.Result().Body)
	testutil.Diff(t, http.StatusOK, w2.Result().StatusCode)
	testutil.Diff(t, "ok", string(b2))
}

func TestPropagationJaeger(t *testing.T) {
	server := common.NewAPI()
	err := app.LoadConfigFiles(server, []string{"./config-propagation-jaeger.yaml"})
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "app/v1",
		Kind:       "OpenTelemetryTracer",
	}
	m, err := api.ReferTypedObject[core.Middleware](server, ref)
	testutil.DiffError(t, nil, nil, err)

	uberID := "4bf92f3577b34da6a3ce929d0e0e4736:00f067aa0ba902b7:0:1"
	uberIDPattern := regexp.MustCompile(`^[0-9a-f]{32}:[0-9a-f]{16}:0:1$`)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Log(r.Header)
		testutil.Diff(t, true, uberID != r.Header.Get("Uber-Trace-Id"))
		testutil.Diff(t, true, uberIDPattern.MatchString(r.Header.Get("Uber-Trace-Id")))
		w.Write([]byte("ok"))
	})
	h := m.Middleware(handler)

	r1 := httptest.NewRequest(http.MethodGet, "http://localhost:8080/", nil)
	r1.Header.Add("Uber-Trace-Id", uberID)
	w1 := httptest.NewRecorder()
	h.ServeHTTP(w1, r1)
	b1, _ := io.ReadAll(w1.Result().Body)
	testutil.Diff(t, http.StatusOK, w1.Result().StatusCode)
	testutil.Diff(t, "ok", string(b1))

	r2 := httptest.NewRequest(http.MethodGet, "http://localhost:8080/", nil)
	w2 := httptest.NewRecorder()
	h.ServeHTTP(w2, r2)
	b2, _ := io.ReadAll(w2.Result().Body)
	testutil.Diff(t, http.StatusOK, w2.Result().StatusCode)
	testutil.Diff(t, "ok", string(b2))
}

func TestPropagationXRay(t *testing.T) {
	server := common.NewAPI()
	err := app.LoadConfigFiles(server, []string{"./config-propagation-xray.yaml"})
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "app/v1",
		Kind:       "OpenTelemetryTracer",
	}
	m, err := api.ReferTypedObject[core.Middleware](server, ref)
	testutil.DiffError(t, nil, nil, err)

	xrayID := "Root=1-12345678-abcdef012345678912345678;Parent=53995c3f42cd8ad8;Sampled=1"
	xrayIDPattern := regexp.MustCompile(`^Root=1-12345678-abcdef012345678912345678;Parent=[0-9a-f]{16};Sampled=1`)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Log(r.Header)
		testutil.Diff(t, true, xrayID != r.Header.Get("X-Amzn-Trace-Id"))
		testutil.Diff(t, true, xrayIDPattern.MatchString(r.Header.Get("X-Amzn-Trace-Id")))
		w.Write([]byte("ok"))
	})
	h := m.Middleware(handler)

	r1 := httptest.NewRequest(http.MethodGet, "http://localhost:8080/", nil)
	r1.Header.Add("X-Amzn-Trace-Id", xrayID)
	w1 := httptest.NewRecorder()
	h.ServeHTTP(w1, r1)
	b1, _ := io.ReadAll(w1.Result().Body)
	testutil.Diff(t, http.StatusOK, w1.Result().StatusCode)
	testutil.Diff(t, "ok", string(b1))
}

func TestPropagationOpenCensus(t *testing.T) {
	server := common.NewAPI()
	err := app.LoadConfigFiles(server, []string{"./config-propagation-opencensus.yaml"})
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "app/v1",
		Kind:       "OpenTelemetryTracer",
	}
	m, err := api.ReferTypedObject[core.Middleware](server, ref)
	testutil.DiffError(t, nil, nil, err)

	opencensusTraceID := "4bf92f3577b34da6a3ce929d0e0e4736"
	opencensusSpanID := "00f067aa0ba902b7"
	opencensusHexTraceID, _ := hex.DecodeString(opencensusTraceID)
	opencensusHexSpanID, _ := hex.DecodeString(opencensusSpanID)
	spanContext := trace.SpanContext{
		TraceID: trace.TraceID(opencensusHexTraceID),
		SpanID:  trace.SpanID(opencensusHexSpanID),
	}
	bin := propagation.Binary(spanContext)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Log(r.Header)
		testutil.Diff(t, true, string(bin) != r.Header.Get("Grpc-Trace-Bin"))
		w.Write([]byte("ok"))
	})
	h := m.Middleware(handler)

	r1 := httptest.NewRequest(http.MethodGet, "http://localhost:8080/", nil)
	r1.Header.Add("Grpc-Trace-Bin", string(bin))
	w1 := httptest.NewRecorder()
	h.ServeHTTP(w1, r1)
	b1, _ := io.ReadAll(w1.Result().Body)
	testutil.Diff(t, http.StatusOK, w1.Result().StatusCode)
	testutil.Diff(t, "ok", string(b1))
}

func TestPropagationOpenTracing(t *testing.T) {
	server := common.NewAPI()
	err := app.LoadConfigFiles(server, []string{"./config-propagation-opentracing.yaml"})
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "app/v1",
		Kind:       "OpenTelemetryTracer",
	}
	m, err := api.ReferTypedObject[core.Middleware](server, ref)
	testutil.DiffError(t, nil, nil, err)

	// In the OpenTracing format, a 64-bit TraceID can be used.
	opentracingTraceID := "a3ce929d0e0e4736"
	opentracingSpanID := "00f067aa0ba902b7"
	idPattern := regexp.MustCompile(`[0-9a-f]{16}`)
	opentracingSampled := "1"

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Log(r.Header)
		testutil.Diff(t, true, idPattern.MatchString(r.Header.Get("ot-tracer-traceid")))
		testutil.Diff(t, true, idPattern.MatchString(r.Header.Get("ot-tracer-spanid")))
		testutil.Diff(t, true, opentracingTraceID == r.Header.Get("ot-tracer-traceid"))
		testutil.Diff(t, true, opentracingSpanID != r.Header.Get("ot-tracer-spanid"))
		testutil.Diff(t, "true", r.Header.Get("ot-tracer-sampled"))
		w.Write([]byte("ok"))
	})
	h := m.Middleware(handler)

	r1 := httptest.NewRequest(http.MethodGet, "http://localhost:8080/", nil)
	r1.Header.Add("ot-tracer-traceid", opentracingTraceID)
	r1.Header.Add("ot-tracer-spanid", opentracingSpanID)
	r1.Header.Add("ot-tracer-sampled", opentracingSampled)
	w1 := httptest.NewRecorder()
	h.ServeHTTP(w1, r1)
	b1, _ := io.ReadAll(w1.Result().Body)
	testutil.Diff(t, http.StatusOK, w1.Result().StatusCode)
	testutil.Diff(t, "ok", string(b1))
}
