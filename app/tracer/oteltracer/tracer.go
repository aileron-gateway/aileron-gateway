package oteltracer

import (
	"context"
	"crypto/tls"
	"net/http"
	"path"
	"runtime"
	"strconv"
	"strings"

	"github.com/aileron-gateway/aileron-gateway/core"
	"github.com/aileron-gateway/aileron-gateway/kernel/uid"
	utilhttp "github.com/aileron-gateway/aileron-gateway/util/http"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

type middlewareContextKey struct{}
type tripperwareContextKey struct{}

var (
	mCtxKey = middlewareContextKey{}
	tCtxKey = tripperwareContextKey{}
)

type otelTracer struct {
	tracer trace.Tracer
	tp     *sdktrace.TracerProvider
	pg     propagation.TextMapPropagator
	// headers is the list of HTTP header names that will be
	// added to the span attributes.
	headers []string
}

func (t *otelTracer) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var c int // counter
		if v := r.Context().Value(mCtxKey); v != nil {
			c = v.(int) + 1
		}

		span, ctx := t.spanContext(r.Context(), r.Header, strconv.Itoa(c)+":middleware")
		defer span.End()

		ctx = context.WithValue(ctx, mCtxKey, c)
		r = r.WithContext(ctx)

		if ptr, file, _, ok := runtime.Caller(2); ok {
			f, _ := runtime.CallersFrames([]uintptr{ptr}).Next()
			span.SetAttributes(attribute.String("caller.file", path.Base(path.Dir(file))+"/"+path.Base(file)))
			span.SetAttributes(attribute.String("caller.func", f.Function))
		}

		// Label operation name.
		span.SetName(strconv.Itoa(c) + ":middleware")

		// Add information only for the root span.
		if c == 0 {
			// Inject will be performed only for the root span.
			// If you set more than one OpenTelemetryTracer middleware, context propagation may not function properly.
			t.pg.Inject(ctx, propagation.HeaderCarrier(r.Header))

			span.SetAttributes(attribute.String("http.id", uid.IDFromContext(ctx)))
			span.SetAttributes(attribute.String("http.schema", schema(r.TLS)))
			span.SetAttributes(attribute.String("http.method", r.Method))
			span.SetAttributes(attribute.String("http.path", r.URL.Path))
			span.SetAttributes(attribute.String("http.query", r.URL.RawQuery))
			span.SetAttributes(attribute.String("net.addr", r.RemoteAddr))
			span.SetAttributes(attribute.String("net.host", r.Host))
			for _, h := range t.headers {
				span.SetAttributes(attribute.StringSlice("http.header."+strings.ToLower(h), r.Header.Values(h)))
			}

			ww := utilhttp.WrapWriter(w)
			w = ww
			defer func() {
				span.SetAttributes(attribute.Int("http.status_code", ww.StatusCode()))
			}()
		}

		next.ServeHTTP(w, r)
	})
}

func (t *otelTracer) Tripperware(next http.RoundTripper) http.RoundTripper {
	return core.RoundTripperFunc(func(r *http.Request) (*http.Response, error) {
		var c int // counter
		if v := r.Context().Value(tCtxKey); v != nil {
			c = v.(int) + 1
		}

		span, ctx := t.spanContext(r.Context(), r.Header, strconv.Itoa(c)+":tripperware")
		defer span.End()

		ctx = context.WithValue(ctx, tCtxKey, c)
		r = r.WithContext(ctx)

		if ptr, file, _, ok := runtime.Caller(2); ok {
			f, _ := runtime.CallersFrames([]uintptr{ptr}).Next()
			span.SetAttributes(attribute.String("caller.file", path.Base(path.Dir(file))+"/"+path.Base(file)))
			span.SetAttributes(attribute.String("caller.func", f.Function))
		}

		// Label operation name.
		span.SetName(strconv.Itoa(c) + ":tripperware")

		// Inject will be performed only for the root span.
		// If you set more than one OpenTelemetryTracer tripperware, context propagation may not function properly.
		if c == 0 {
			t.pg.Inject(ctx, propagation.HeaderCarrier(r.Header))
		}

		res, err := next.RoundTrip(r)

		if c == 0 {
			span.SetAttributes(attribute.String("http.id", uid.IDFromContext(ctx)))
			span.SetAttributes(attribute.String("http.schema", schema(r.TLS)))
			span.SetAttributes(attribute.String("http.method", r.Method))
			span.SetAttributes(attribute.String("http.path", r.URL.Path))
			span.SetAttributes(attribute.String("http.query", r.URL.RawQuery))
			span.SetAttributes(attribute.String("peer.host", r.URL.Host))
			for _, h := range t.headers {
				span.SetAttributes(attribute.StringSlice("http.header."+strings.ToLower(h), r.Header.Values(h)))
			}
			if err != nil {
				span.SetAttributes(attribute.Int("http.status_code", 0))
			} else {
				span.SetAttributes(attribute.Int("http.status_code", res.StatusCode))
			}
		}

		return res, err
	})
}

// spanContext returns nre span and context for the request.
func (t *otelTracer) spanContext(ctx context.Context, header http.Header, name string) (trace.Span, context.Context) {
	var span trace.Span

	if parentSpan := trace.SpanFromContext(ctx); parentSpan.SpanContext().IsValid() {
		ctx, span = t.tracer.Start(
			ctx,
			name,
			trace.WithSpanKind(trace.SpanKindInternal),
			trace.WithLinks(trace.Link{SpanContext: parentSpan.SpanContext()}),
		)
	} else {
		remoteCtx := t.pg.Extract(ctx, propagation.HeaderCarrier(header))
		sc := trace.SpanFromContext(remoteCtx).SpanContext()

		if sc.IsValid() {
			ctx, span = t.tracer.Start(
				remoteCtx,
				name,
				trace.WithSpanKind(trace.SpanKindInternal),
				trace.WithLinks(trace.Link{SpanContext: sc}),
			)
		} else {
			ctx, span = t.tracer.Start(
				ctx,
				name,
				trace.WithSpanKind(trace.SpanKindServer),
			)
		}
	}

	return span, ctx
}

// Trace is the method than can be called from any types of resources.
// Callers must update their context with the returned one.
// The returned function with finishes spans must be called when finishing spans.
func (t *otelTracer) Trace(ctx context.Context, name string, tags map[string]string) (context.Context, func()) {
	var span trace.Span
	var spanCtx context.Context

	if parentSpan := trace.SpanFromContext(ctx); parentSpan.SpanContext().IsValid() {
		spanCtx, span = t.tracer.Start(
			ctx,
			name,
			trace.WithSpanKind(trace.SpanKindInternal),
			trace.WithLinks(trace.Link{SpanContext: parentSpan.SpanContext()}),
		)
	} else {
		spanCtx, span = t.tracer.Start(
			ctx,
			name,
			trace.WithSpanKind(trace.SpanKindServer),
		)
	}

	attributes := convertTagsToAttributes(tags)
	span.SetAttributes(attributes...)

	return spanCtx, func() {
		span.End()
	}
}

// Finalize ensures t.tp.Shutdown is called before the application halts.
// Shutdown flushes trace data remaining in the buffer.
// This implements the core.Finalizer interface.
func (t *otelTracer) Finalize() error {
	return t.tp.Shutdown(context.Background())
}

func schema(s *tls.ConnectionState) string {
	if s == nil {
		return "http"
	}
	return "https"
}

func convertTagsToAttributes(tags map[string]string) []attribute.KeyValue {
	attributes := make([]attribute.KeyValue, 0, len(tags))
	for k, v := range tags {
		attributes = append(attributes, attribute.String(k, v))
	}
	return attributes
}
