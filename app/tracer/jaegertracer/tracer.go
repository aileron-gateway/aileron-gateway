package jaegertracer

import (
	"cmp"
	"context"
	"crypto/tls"
	"errors"
	"io"
	"net/http"
	"path"
	"runtime"
	"strconv"
	"strings"

	"github.com/aileron-gateway/aileron-gateway/core"
	"github.com/aileron-gateway/aileron-gateway/kernel/uid"
	utilhttp "github.com/aileron-gateway/aileron-gateway/util/http"
	"github.com/opentracing/opentracing-go"
)

type middlewareContextKey struct{}
type tripperwareContextKey struct{}

var (
	mCtxKey = middlewareContextKey{}
	tCtxKey = tripperwareContextKey{}
)

type jaegerTracer struct {
	// tracer is a jaeger tracer.
	tracer opentracing.Tracer
	// closer closes the tracer.
	// Close method should be called before shutting down
	// the application to flush all tracing data
	// in the internal buffer.
	closer io.Closer

	// mNames is a combination of middleware names and numbers.
	// This can be used to make it possible to label user defined name to the spans.
	mNames map[int]string
	// tNames is a combination of tripperware names and numbers.
	// This can be used to make it possible to label user defined name to the spans.
	tNames map[int]string
	// headers is the list of HTTP header names that will be
	// added to the span tags.
	headers []string
}

func (t *jaegerTracer) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var c int // counter
		if v := r.Context().Value(mCtxKey); v != nil {
			c = v.(int) + 1
		}

		span, ctx := t.spanContext(r.Context(), r.Header, strconv.Itoa(c)+":middleware")
		defer span.Finish()

		if ptr, file, _, ok := runtime.Caller(2); ok {
			f, _ := runtime.CallersFrames([]uintptr{ptr}).Next()
			span.SetTag("caller.file", path.Base(path.Dir(file))+"/"+path.Base(file))
			span.SetTag("caller.func", f.Function)
		}

		// Label operation name.
		// Default value like "1:middleware" will be used when not configured.
		span.SetOperationName(cmp.Or(t.mNames[c], strconv.Itoa(c)+":middleware"))

		// Propagate recorder with context.
		ctx = context.WithValue(ctx, mCtxKey, c)
		r = r.WithContext(ctx)

		// Add information only for the root span.
		if c == 0 {
			// TODO: Implement error handling just in case.
			// According to the opentracing.Tracer.Inject implementation,
			// err should always be nil here.
			// https://pkg.go.dev/github.com/jaegertracing/jaeger-client-go#Tracer.Inject
			_ = t.tracer.Inject(span.Context(), opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(r.Header))

			span.SetTag("http.id", uid.IDFromContext(ctx))
			span.SetTag("http.schema", schema(r.TLS))
			span.SetTag("http.method", r.Method)
			span.SetTag("http.path", r.URL.Path)
			span.SetTag("http.query", r.URL.RawQuery)
			span.SetTag("net.addr", r.RemoteAddr)
			span.SetTag("net.host", r.Host)
			for _, h := range t.headers {
				span.SetTag("http.header."+strings.ToLower(h), r.Header.Values(h))
			}

			ww := utilhttp.WrapWriter(w)
			w = ww
			defer func() {
				span.SetTag("http.status_code", ww.StatusCode())
			}()
		}

		next.ServeHTTP(w, r)
	})
}

func (t *jaegerTracer) Tripperware(next http.RoundTripper) http.RoundTripper {
	return core.RoundTripperFunc(func(r *http.Request) (*http.Response, error) {
		var c int // counter
		if v := r.Context().Value(tCtxKey); v != nil {
			c = v.(int) + 1
		}

		span, ctx := t.spanContext(r.Context(), r.Header, strconv.Itoa(c)+":tripperware")
		defer span.Finish()

		if ptr, file, _, ok := runtime.Caller(2); ok {
			f, _ := runtime.CallersFrames([]uintptr{ptr}).Next()
			span.SetTag("caller.file", path.Base(path.Dir(file))+"/"+path.Base(file))
			span.SetTag("caller.func", f.Function)
		}

		// Label operation name.
		// Default value like "1:tripperware" will be used when not configured.
		span.SetOperationName(cmp.Or(t.tNames[c], strconv.Itoa(c)+":tripperware"))

		// Propagate recorder with context.
		ctx = context.WithValue(ctx, tCtxKey, c)
		r = r.WithContext(ctx)

		if c == 0 {
			// TODO: Implement error handling just in case.
			// According to the opentracing.Tracer.Inject implementation,
			// err should always be nil here.
			// https://pkg.go.dev/github.com/jaegertracing/jaeger-client-go#Tracer.Inject
			_ = t.tracer.Inject(span.Context(), opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(r.Header))
		}

		res, err := next.RoundTrip(r)

		if c == 0 {
			span.SetTag("http.id", uid.IDFromContext(ctx))
			span.SetTag("http.schema", schema(r.TLS))
			span.SetTag("http.method", r.Method)
			span.SetTag("http.path", r.URL.Path)
			span.SetTag("http.query", r.URL.RawQuery)
			span.SetTag("peer.host", r.URL.Host)
			for _, h := range t.headers {
				span.SetTag("http.header."+strings.ToLower(h), r.Header.Values(h))
			}
			if err != nil {
				span.SetTag("http.status_code", 0)
			} else {
				span.SetTag("http.status_code", res.StatusCode)
			}
		}

		return res, err
	})
}

// spanContext returns nre span and context for the request.
func (t *jaegerTracer) spanContext(ctx context.Context, header http.Header, name string) (opentracing.Span, context.Context) {
	var span opentracing.Span
	if parentSpan := opentracing.SpanFromContext(ctx); parentSpan != nil {
		span = t.tracer.StartSpan(name, opentracing.ChildOf(parentSpan.Context()))
	} else {
		carrier := opentracing.HTTPHeadersCarrier(header)
		sc, err := t.tracer.Extract(opentracing.HTTPHeaders, carrier)

		if errors.Is(err, opentracing.ErrSpanContextNotFound) {
			span = t.tracer.StartSpan(name)
		} else {
			span = t.tracer.StartSpan(name, opentracing.ChildOf(sc))
		}
	}

	ctx = opentracing.ContextWithSpan(ctx, span)
	return span, ctx
}

// Trace is the method than can be called from any types of resources.
// Callers must update their context with the returned one.
// The returned function with finishes spans must be called when finishing spans.
func (t *jaegerTracer) Trace(ctx context.Context, name string, tags map[string]string) (context.Context, func()) {
	var span opentracing.Span
	if parentSpan := opentracing.SpanFromContext(ctx); parentSpan != nil {
		span = t.tracer.StartSpan(name, opentracing.ChildOf(parentSpan.Context()))
	} else {
		span = t.tracer.StartSpan(name)
	}

	for k, v := range tags {
		span.SetTag(k, v)
	}

	ctx = opentracing.ContextWithSpan(ctx, span)
	return ctx, span.Finish
}

// Finalize ensures t.closer.Close is called before the application halts.
// Close flushes trace data remaining in the buffer.
// This implements the core.Finalizer interface.
func (t *jaegerTracer) Finalize() error {
	return t.closer.Close()
}

func schema(s *tls.ConnectionState) string {
	if s == nil {
		return "http"
	}
	return "https"
}
