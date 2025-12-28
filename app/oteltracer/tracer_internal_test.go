// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package oteltracer

import (
	"context"
	"crypto/tls"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aileron-gateway/aileron-gateway/core"
	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"go.opentelemetry.io/contrib/propagators/autoprop"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
)

func TestMiddleware(t *testing.T) {
	type condition struct {
		mCtxKey    any
		headers    []string
		childSpan  bool
		parentSpan bool
		httpsFlag  bool
	}

	type action struct {
		statusCode int
		attributes []attribute.KeyValue
		body       string
		name       string
	}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"nil tracer",
			&condition{},
			&action{
				statusCode: http.StatusOK,
				attributes: []attribute.KeyValue{
					attribute.String("caller.file", "oteltracer/tracer_internal_test.go"),
					attribute.String("caller.func", "github.com/aileron-gateway/aileron-gateway/app/oteltracer.(*otelTracer).Middleware"),
					attribute.String("http.id", ""),
					attribute.String("http.schema", "http"),
					attribute.String("http.method", "GET"),
					attribute.String("http.path", "/"),
					attribute.String("http.query", ""),
					attribute.String("net.addr", "192.0.2.1:1234"),
					attribute.String("net.host", "example.com"),
					attribute.Int("http.status_code", 0),
				},
				name: "0:middleware",
			},
		),
		gen(
			"set mCtxKey",
			&condition{
				mCtxKey: 1,
			},
			&action{
				statusCode: http.StatusOK,
				attributes: []attribute.KeyValue{
					attribute.String("caller.file", "oteltracer/tracer_internal_test.go"),
					attribute.String("caller.func", "github.com/aileron-gateway/aileron-gateway/app/oteltracer.(*otelTracer).Middleware"),
				},
				name: "2:middleware",
			},
		),
		gen(
			"start ChildSpan",
			&condition{
				childSpan:  true,
				parentSpan: false,
			},
			&action{
				statusCode: http.StatusOK,
				attributes: []attribute.KeyValue{
					attribute.String("caller.file", "oteltracer/tracer_internal_test.go"),
					attribute.String("caller.func", "github.com/aileron-gateway/aileron-gateway/app/oteltracer.(*otelTracer).Middleware"),
					attribute.String("http.id", ""),
					attribute.String("http.schema", "http"),
					attribute.String("http.method", "GET"),
					attribute.String("http.path", "/"),
					attribute.String("http.query", ""),
					attribute.String("net.addr", "192.0.2.1:1234"),
					attribute.String("net.host", "example.com"),
					attribute.Int("http.status_code", 0),
				},
				name: "0:middleware",
			},
		),
		gen(
			"start ParentSpan",
			&condition{
				childSpan:  false,
				parentSpan: true,
			},
			&action{
				statusCode: http.StatusOK,
				attributes: []attribute.KeyValue{
					attribute.String("caller.file", "oteltracer/tracer_internal_test.go"),
					attribute.String("caller.func", "github.com/aileron-gateway/aileron-gateway/app/oteltracer.(*otelTracer).Middleware"),
					attribute.String("http.id", ""),
					attribute.String("http.schema", "http"),
					attribute.String("http.method", "GET"),
					attribute.String("http.path", "/"),
					attribute.String("http.query", ""),
					attribute.String("net.addr", "192.0.2.1:1234"),
					attribute.String("net.host", "example.com"),
					attribute.Int("http.status_code", 0),
				},
				name: "0:middleware",
			},
		),
		gen(
			"start Headers",
			&condition{
				headers: []string{"testHeader"},
			},
			&action{
				statusCode: http.StatusOK,
				attributes: []attribute.KeyValue{
					attribute.String("caller.file", "oteltracer/tracer_internal_test.go"),
					attribute.String("caller.func", "github.com/aileron-gateway/aileron-gateway/app/oteltracer.(*otelTracer).Middleware"),
					attribute.String("http.id", ""),
					attribute.String("http.schema", "http"),
					attribute.String("http.method", "GET"),
					attribute.String("http.path", "/"),
					attribute.String("http.query", ""),
					attribute.String("net.addr", "192.0.2.1:1234"),
					attribute.String("net.host", "example.com"),
					attribute.StringSlice("http.header.testheader", []string{}),
					attribute.Int("http.status_code", 0),
				},
				name: "0:middleware",
			},
		),
		gen(
			"set HTTPS schema",
			&condition{
				httpsFlag: true,
			},
			&action{
				statusCode: http.StatusOK,
				attributes: []attribute.KeyValue{
					attribute.String("caller.file", "oteltracer/tracer_internal_test.go"),
					attribute.String("caller.func", "github.com/aileron-gateway/aileron-gateway/app/oteltracer.(*otelTracer).Middleware"),
					attribute.String("http.id", ""),
					attribute.String("http.schema", "https"),
					attribute.String("http.method", "GET"),
					attribute.String("http.path", "/"),
					attribute.String("http.query", ""),
					attribute.String("net.addr", "192.0.2.1:1234"),
					attribute.String("net.host", "example.com"),
					attribute.Int("http.status_code", 0),
				},
				name: "0:middleware",
			},
		),
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			exporter := tracetest.NewInMemoryExporter()
			tracerProvider := sdktrace.NewTracerProvider(
				sdktrace.WithSyncer(exporter),
			)

			tracer := tracerProvider.Tracer("oteltracer")

			ctx := context.Background()

			if tt.C.parentSpan {
				parentCtx, parentSpan := tracer.Start(ctx, "parentSpan")
				defer parentSpan.End()

				ctx = context.WithValue(parentCtx, mCtxKey, tt.C.mCtxKey)
			} else {
				ctx = context.WithValue(ctx, mCtxKey, tt.C.mCtxKey)
			}

			// In unit test, only set TraceContext and Baggage, and verify the behavior in subsequent tests.
			props := []propagation.TextMapPropagator{
				propagation.TraceContext{},
				propagation.Baggage{},
			}
			pg := autoprop.NewTextMapPropagator(props...)

			ot := &otelTracer{
				tracer:  tracer,
				tp:      tracerProvider,
				pg:      pg,
				headers: tt.C.headers,
			}

			exporter.ExportSpans(ctx, exporter.GetSpans().Snapshots())

			req := httptest.NewRequest(http.MethodGet, "/", nil)

			if tt.C.childSpan {
				childCtx, childSpan := tracer.Start(ctx, "childSpan")
				defer childSpan.End()

				ot.pg.Inject(childCtx, propagation.HeaderCarrier(req.Header))
			}

			if tt.C.httpsFlag {
				req.TLS = &tls.ConnectionState{}
			}

			req = req.WithContext(ctx)
			resp := httptest.NewRecorder()
			h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

			ot.Middleware(h).ServeHTTP(resp, req)

			for _, span := range exporter.GetSpans() {
				testutil.Diff(t, tt.A.attributes, span.Attributes, cmp.AllowUnexported(attribute.Value{}))
				testutil.Diff(t, tt.A.name, span.Name)
			}
			testutil.Diff(t, tt.A.statusCode, resp.Code)
			testutil.Diff(t, tt.A.body, resp.Body.String())
		})
	}
}

// Custom error to simulate round trip error
var errRoundTrip = errors.New("round trip error")

func TestTripperware(t *testing.T) {
	type condition struct {
		tCtxKey      any
		headers      []string
		roundTripErr bool
	}

	type action struct {
		err        error
		statusCode int
		attributes []attribute.KeyValue
		name       string
	}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"nil tracer",
			&condition{},
			&action{
				statusCode: http.StatusOK,
				attributes: []attribute.KeyValue{
					attribute.String("caller.file", "oteltracer/tracer_internal_test.go"),
					attribute.String("caller.func", "github.com/aileron-gateway/aileron-gateway/app/oteltracer.(*otelTracer).Tripperware"),
					attribute.String("http.id", ""),
					attribute.String("http.schema", "http"),
					attribute.String("http.method", "GET"),
					attribute.String("http.path", "/"),
					attribute.String("http.query", ""),
					attribute.String("peer.host", ""),
					attribute.Int("http.status_code", 200),
				},
				name: "0:tripperware",
			},
		),
		gen(
			"set tCtxKey",
			&condition{
				tCtxKey: 1,
			},
			&action{
				statusCode: http.StatusOK,
				attributes: []attribute.KeyValue{
					attribute.String("caller.file", "oteltracer/tracer_internal_test.go"),
					attribute.String("caller.func", "github.com/aileron-gateway/aileron-gateway/app/oteltracer.(*otelTracer).Tripperware"),
				},
				name: "2:tripperware",
			},
		),
		gen(
			"set Headers",
			&condition{
				headers: []string{"testHeader"},
			},
			&action{
				statusCode: http.StatusOK,
				attributes: []attribute.KeyValue{
					attribute.String("caller.file", "oteltracer/tracer_internal_test.go"),
					attribute.String("caller.func", "github.com/aileron-gateway/aileron-gateway/app/oteltracer.(*otelTracer).Tripperware"),
					attribute.String("http.id", ""),
					attribute.String("http.schema", "http"),
					attribute.String("http.method", "GET"),
					attribute.String("http.path", "/"),
					attribute.String("http.query", ""),
					attribute.String("peer.host", ""),
					attribute.StringSlice("http.header.testheader", []string{}),
					attribute.Int("http.status_code", 200),
				},
				name: "0:tripperware",
			},
		),
		gen(
			"cause RoundTripError",
			&condition{
				roundTripErr: true,
			},
			&action{
				statusCode: http.StatusOK,
				attributes: []attribute.KeyValue{
					attribute.String("caller.file", "oteltracer/tracer_internal_test.go"),
					attribute.String("caller.func", "github.com/aileron-gateway/aileron-gateway/app/oteltracer.(*otelTracer).Tripperware"),
					attribute.String("http.id", ""),
					attribute.String("http.schema", "http"),
					attribute.String("http.method", "GET"),
					attribute.String("http.path", "/"),
					attribute.String("http.query", ""),
					attribute.String("peer.host", ""),
					attribute.Int("http.status_code", 0),
				},
				name: "0:tripperware",
				err:  errRoundTrip,
			},
		),
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			exporter := tracetest.NewInMemoryExporter()
			tracerProvider := sdktrace.NewTracerProvider(
				sdktrace.WithSyncer(exporter),
			)
			tracer := tracerProvider.Tracer("oteltracer")

			ctx := context.Background()
			ctx = context.WithValue(ctx, tCtxKey, tt.C.tCtxKey)

			// In unit test, only set TraceContext and Baggage, and verify the behavior in subsequent tests.
			props := []propagation.TextMapPropagator{
				propagation.TraceContext{},
				propagation.Baggage{},
			}
			pg := autoprop.NewTextMapPropagator(props...)

			ot := &otelTracer{
				tracer: tracer,
				// TODO: Implement here so that the propagation settings can be modified.
				pg:      pg,
				tp:      tracerProvider,
				headers: tt.C.headers,
			}

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req = req.WithContext(ctx)

			r := core.RoundTripperFunc(func(r *http.Request) (*http.Response, error) {
				resp := &http.Response{
					StatusCode: http.StatusOK,
					Status:     "StatusOK",
					Proto:      "HTTP/1.0",
					ProtoMajor: 1,
					ProtoMinor: 0,
					Body:       nil,
				}
				if tt.C.roundTripErr {
					return resp, errRoundTrip
				}
				return resp, nil
			})

			opts := []cmp.Option{
				cmpopts.EquateErrors(),
			}

			resp, err := ot.Tripperware(r).RoundTrip(req)
			testutil.Diff(t, tt.A.err, err, opts...)

			for _, span := range exporter.GetSpans() {
				testutil.Diff(t, tt.A.attributes, span.Attributes, cmp.AllowUnexported(attribute.Value{}))
				testutil.Diff(t, tt.A.name, span.Name)
			}
			testutil.Diff(t, tt.A.statusCode, resp.StatusCode)
		})
	}
}

func TestTrace(t *testing.T) {
	type condition struct {
		name       string
		tags       map[string]string
		parentSpan bool
	}

	type action struct {
		name       string
		attributes []attribute.KeyValue
	}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"empty name and attributes",
			&condition{},
			&action{
				name:       "",
				attributes: []attribute.KeyValue(nil),
			},
		),
		gen(
			"set name",
			&condition{
				name: "testName",
			},
			&action{
				name:       "testName",
				attributes: []attribute.KeyValue(nil),
			},
		),
		gen(
			"set single attribute",
			&condition{
				name: "testName",
				tags: map[string]string{
					"testTagKey": "testTagValue",
				},
			},
			&action{
				name: "testName",
				attributes: []attribute.KeyValue{
					attribute.String("testTagKey", "testTagValue"),
				},
			},
		),
		gen(
			"set multiple attributes",
			&condition{
				name: "testName",
				tags: map[string]string{
					"testFirstTagKey":  "testFirstTagValue",
					"testSecondTagKey": "testSecondTagValue",
				},
			},
			&action{
				name: "testName",
				attributes: []attribute.KeyValue{
					attribute.String("testFirstTagKey", "testFirstTagValue"),
					attribute.String("testSecondTagKey", "testSecondTagValue"),
				},
			},
		),
		gen(
			"set parent span",
			&condition{
				parentSpan: true,
			},
			&action{
				name:       "",
				attributes: []attribute.KeyValue(nil),
			},
		),
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			exporter := tracetest.NewInMemoryExporter()
			tracerProvider := sdktrace.NewTracerProvider(
				sdktrace.WithSyncer(exporter),
			)

			tracer := tracerProvider.Tracer("oteltracer")

			ctx := context.Background()

			ot := &otelTracer{
				tracer: tracer,
				tp:     tracerProvider,
				// TODO: Implement here so that the propagation settings can be modified.
				pg: propagation.TraceContext{},
			}

			if tt.C.parentSpan {
				parentCtx, parentSpan := tracer.Start(ctx, "parentSpan")
				defer parentSpan.End()

				ctx = context.WithValue(parentCtx, mCtxKey, 0)
			}

			_, finish := ot.Trace(ctx, tt.C.name, tt.C.tags)
			finish()

			opts := []cmp.Option{
				cmp.AllowUnexported(attribute.Value{}),
				cmpopts.SortSlices(func(i, j attribute.KeyValue) bool {
					return i.Key < j.Key
				}),
			}

			for _, span := range exporter.GetSpans() {
				testutil.Diff(t, tt.A.attributes, span.Attributes, opts...)
				testutil.Diff(t, tt.A.name, span.Name)
			}
		})
	}
}

type mockSpanProcessor struct {
	isShutdown bool
}

func (m *mockSpanProcessor) OnStart(parent context.Context, s sdktrace.ReadWriteSpan) {
}

func (m *mockSpanProcessor) OnEnd(s sdktrace.ReadOnlySpan) {
}

func (m *mockSpanProcessor) Shutdown(ctx context.Context) error {
	m.isShutdown = true
	return nil
}

func (m *mockSpanProcessor) ForceFlush(ctx context.Context) error {
	return nil
}

func TestFinalize(t *testing.T) {
	type condition struct {
		mockSpanProcessor *mockSpanProcessor
	}

	type action struct {
		isShutdown bool
		err        error
	}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"default TracerProvider",
			&condition{
				mockSpanProcessor: &mockSpanProcessor{},
			},
			&action{
				isShutdown: true,
				err:        nil,
			},
		),
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			exporter := tracetest.NewInMemoryExporter()
			mockSpanProcessor := tt.C.mockSpanProcessor

			tracerProvider := sdktrace.NewTracerProvider(
				sdktrace.WithSpanProcessor(mockSpanProcessor),
				sdktrace.WithSyncer(exporter),
			)

			ot := &otelTracer{
				tp: tracerProvider,
			}
			err := ot.Finalize()

			testutil.Diff(t, tt.A.isShutdown, mockSpanProcessor.isShutdown)
			testutil.Diff(t, tt.A.err, err)
		})
	}
}
