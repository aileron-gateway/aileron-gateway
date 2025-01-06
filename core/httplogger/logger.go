package httplogger

import (
	"cmp"
	"context"
	"io"
	"mime"
	"net/http"
	"strconv"
	"time"

	"github.com/aileron-gateway/aileron-gateway/core"
	"github.com/aileron-gateway/aileron-gateway/kernel/log"
	utilhttp "github.com/aileron-gateway/aileron-gateway/util/http"
)

var (
	emptyHeader = make(map[string]string, 0)
	noopFunc    = func() {}
	nopByteFunc = func() []byte { return nil }
)

// WrappedWriter wraps http.ResponseWriter and holds http status code.
// This implements io.Writer interface.
type wrappedWriter struct {
	http.ResponseWriter
	status      int
	written     bool
	writtenFunc func()

	dump bool      // dump is set in writtenFunc.
	w    io.Writer // w is set in writtenFunc.

	flushChecked bool
	flushFunc    func()
}

func (w *wrappedWriter) Unwrap() http.ResponseWriter {
	return w.ResponseWriter
}

// Flush flushes the buffer of internal ResponseWriter.
// This is the implementation of [http.Flusher] interface.
// Implementing [http.Flusher.]Flush() make the reverse proxy handler
// use this wrappedWriter when proxy streaming bodies.
func (w *wrappedWriter) Flush() {
	if w.flushChecked {
		w.flushFunc()
		return
	}
	w.flushChecked = true
	w.flushFunc = noopFunc
	rw := w.ResponseWriter
	for {
		if flusher, ok := rw.(http.Flusher); ok {
			w.flushFunc = flusher.Flush
			flusher.Flush()
			return
		}
		if uw, ok := rw.(interface{ Unwrap() http.ResponseWriter }); ok {
			rw = uw.Unwrap()
			continue
		}
		return
	}
}

func (w *wrappedWriter) WriteHeader(statusCode int) {
	w.status = statusCode
	if !w.written {
		w.written = true
		w.writtenFunc()
		w.writtenFunc = nil // Release for GC.
	}
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *wrappedWriter) Write(b []byte) (int, error) {
	if !w.written {
		w.written = true
		w.writtenFunc()
		w.writtenFunc = nil // Release for GC.
	}
	if w.dump {
		_, _ = w.w.Write(b) // Explicitly ignore the returned  error.
	}
	return w.ResponseWriter.Write(b)
}

// httpLogger records log attributes.
// This implements core.Middleware interface.
type httpLogger struct {
	req *baseLogger
	res *baseLogger

	zone    *time.Location
	timeFmt string
}

func (lg *httpLogger) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var id string // Log ID in the case for there is no request IDs.
		if v := r.Context().Value(idContextKey); v != nil {
			id = v.(string)
		} else {
			id = newLogID()
			r = r.WithContext(context.WithValue(r.Context(), idContextKey, id))
		}

		startTime := time.Now()

		attr := &requestAttrs{
			typ:    "server",
			id:     id,
			time:   time.Now().In(lg.zone).Format(lg.timeFmt),
			method: r.Method,
			path:   r.URL.Path,
			query:  lg.req.logQuery(r.URL.RawQuery),
			host:   r.Host,
			remote: r.RemoteAddr,
			proto:  r.Proto,
			size:   r.ContentLength,
			header: lg.req.logHeaders(r.Header),
		}
		lg.req.logOutput(r.Context(), "server", attr.accessKeyValues(), attr.TagFunc)

		ww := &wrappedWriter{
			ResponseWriter: w,
		}
		w = ww

		// writtenFunc is fired only once when response status or body was written to the ResponseWriter.
		ww.writtenFunc = func() { //nolint:contextcheck // Function `Middleware$1$1` should pass the context parameter
			size := int64(-1) // Use -1 when there is no Content-Length header. Streaming, HTTP/2, etc...
			if cl := w.Header()["Content-Length"]; len(cl) > 0 {
				size, _ = strconv.ParseInt(cl[0], 10, 64)
			}
			attr := &responseAttrs{
				typ:      "server",
				id:       id,
				time:     time.Now().In(lg.zone).Format(lg.timeFmt),
				duration: time.Since(startTime).Microseconds(),
				status:   ww.status,
				size:     size,
				header:   lg.res.logHeaders(w.Header()),
			}
			lg.res.logOutput(r.Context(), "server", attr.accessKeyValues(), attr.TagFunc)
		}

		next.ServeHTTP(w, r)
	})
}

func (lg *httpLogger) Tripperware(next http.RoundTripper) http.RoundTripper {
	return core.RoundTripperFunc(func(r *http.Request) (w *http.Response, err error) {
		var id string // Log ID in the case for there is no request IDs.
		if v := r.Context().Value(idContextKey); v != nil {
			id = v.(string)
		} else {
			id = newLogID()
			r = r.WithContext(context.WithValue(r.Context(), idContextKey, id))
		}

		startTime := time.Now()

		attr := &requestAttrs{
			typ:    "client",
			id:     id,
			time:   time.Now().In(lg.zone).Format(lg.timeFmt),
			method: r.Method,
			path:   r.URL.Path,
			query:  lg.req.logQuery(r.URL.RawQuery),
			host:   cmp.Or(r.Host, r.URL.Host),
			remote: r.RemoteAddr,
			proto:  r.Proto,
			size:   r.ContentLength,
			header: lg.req.logHeaders(r.Header),
		}
		lg.req.logOutput(r.Context(), "client", attr.accessKeyValues(), attr.TagFunc)

		defer func(ctx context.Context) {
			attr := &responseAttrs{
				typ:      "client",
				id:       id,
				time:     time.Now().In(lg.zone).Format(lg.timeFmt),
				duration: time.Since(startTime).Microseconds(),
				header:   emptyHeader,
			}
			if w != nil {
				attr.status = w.StatusCode
				attr.size = w.ContentLength
				attr.header = lg.res.logHeaders(w.Header)
			}
			lg.res.logOutput(ctx, "client", attr.accessKeyValues(), attr.TagFunc)
		}(r.Context())

		return next.RoundTrip(r)
	})
}

// httpLogger records log attributes.
// This implements core.Middleware interface.
type journalLogger struct {
	lg  log.Logger
	eh  core.ErrorHandler
	req *baseLogger
	res *baseLogger

	zone    *time.Location
	timeFmt string
}

func (lg *journalLogger) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var id string // Log ID in the case for there is no request IDs.
		if v := r.Context().Value(idContextKey); v != nil {
			id = v.(string)
		} else {
			id = newLogID()
			r = r.WithContext(context.WithValue(r.Context(), idContextKey, id))
		}

		startTime := time.Now()

		var body []byte
		if r.Body != nil && r.Body != http.NoBody {
			mt, _, _ := mime.ParseMediaType(r.Header.Get("Content-Type"))
			b, rc, err := lg.req.bodyReadCloser(id+".svr.req.bin", mt, r.ContentLength, r.Body)
			if err != nil {
				err = core.ErrCoreLogger.WithStack(err, nil)
				lg.eh.ServeHTTPError(w, r, utilhttp.NewHTTPError(err, http.StatusInternalServerError))
				return
			}
			body = b
			r.Body = rc
			defer rc.Close()
		}

		attr := &requestAttrs{
			typ:    "server",
			id:     id,
			time:   time.Now().In(lg.zone).Format(lg.timeFmt),
			method: r.Method,
			path:   r.URL.Path,
			query:  lg.req.logQuery(r.URL.RawQuery),
			host:   r.Host,
			remote: r.RemoteAddr,
			proto:  r.Proto,
			size:   r.ContentLength,
			header: lg.req.logHeaders(r.Header),
			body:   string(body),
		}
		lg.req.logOutput(r.Context(), "server", attr.journalKeyValues(), attr.TagFunc)

		ww := &wrappedWriter{
			ResponseWriter: w,
		}
		w = ww

		size := int64(-1)       // Use -1 when there is no Content-Length header. Streaming, HTTP/2, etc...
		byteFunc := nopByteFunc // byteFunc returns response body bytes.
		// writtenFunc is fired only once when response status or body was written to the ResponseWriter.
		ww.writtenFunc = func() { //nolint:contextcheck // Function `Middleware$1$1` should pass the context parameter
			if cl := w.Header()["Content-Length"]; len(cl) > 0 {
				size, _ = strconv.ParseInt(cl[0], 10, 64)
			}
			mt, _, _ := mime.ParseMediaType(w.Header().Get("Content-Type"))
			bf, bw, err := lg.res.bodyWriter(id+".svr.res.bin", mt, size)
			if err != nil {
				err := core.ErrCoreLogger.WithStack(err, nil)
				lg.lg.Error(r.Context(), err.Name(), err.Map()) // Logging only.
			}
			if bw != nil {
				ww.dump = true
				ww.w = bw
				byteFunc = bf
			}
		}

		defer func(ctx context.Context) {
			attr := &responseAttrs{
				typ:      "server",
				id:       id,
				time:     time.Now().In(lg.zone).Format(lg.timeFmt),
				duration: time.Since(startTime).Microseconds(),
				status:   ww.status,
				size:     size,
				header:   lg.res.logHeaders(w.Header()),
				body:     string(byteFunc()),
			}
			lg.res.logOutput(ctx, "server", attr.journalKeyValues(), attr.TagFunc)
		}(r.Context())

		next.ServeHTTP(w, r)
	})
}

func (lg *journalLogger) Tripperware(next http.RoundTripper) http.RoundTripper {
	return core.RoundTripperFunc(func(r *http.Request) (w *http.Response, err error) {
		var id string // Log ID in the case for there is no request IDs.
		if v := r.Context().Value(idContextKey); v != nil {
			id = v.(string)
		} else {
			id = newLogID()
			r = r.WithContext(context.WithValue(r.Context(), idContextKey, id))
		}

		startTime := time.Now()

		var body []byte
		if r.Body != nil && r.Body != http.NoBody {
			mt, _, _ := mime.ParseMediaType(r.Header.Get("Content-Type"))
			b, rc, err := lg.req.bodyReadCloser(id+".cli.req.bin", mt, r.ContentLength, r.Body)
			if err != nil {
				return nil, core.ErrCoreLogger.WithStack(err, nil)
			}
			body = b
			r.Body = rc
			defer rc.Close()
		}

		attr := &requestAttrs{
			typ:    "client",
			id:     id,
			time:   time.Now().In(lg.zone).Format(lg.timeFmt),
			method: r.Method,
			path:   r.URL.Path,
			query:  lg.req.logQuery(r.URL.RawQuery),
			host:   cmp.Or(r.Host, r.URL.Host),
			remote: r.RemoteAddr,
			proto:  r.Proto,
			size:   r.ContentLength,
			header: lg.req.logHeaders(r.Header),
			body:   string(body),
		}
		lg.req.logOutput(r.Context(), "client", attr.journalKeyValues(), attr.TagFunc)

		defer func(ctx context.Context) {
			attr := &responseAttrs{
				typ:      "client",
				id:       id,
				time:     time.Now().In(lg.zone).Format(lg.timeFmt),
				duration: time.Since(startTime).Microseconds(),
				header:   emptyHeader,
			}
			if w != nil {
				attr.status = w.StatusCode
				attr.size = w.ContentLength
				attr.header = lg.res.logHeaders(w.Header)
				mt, _, _ := mime.ParseMediaType(w.Header.Get("Content-Type"))
				body, rc, err := lg.res.bodyReadCloser(id+".cli.res.bin", mt, w.ContentLength, w.Body)
				if err != nil {
					err := core.ErrCoreLogger.WithStack(err, nil)
					lg.lg.Error(ctx, err.Name(), err.Map()) // Logging only because the upstream already returned response.
				}
				w.Body = rc
				attr.body = string(body)
			}
			lg.res.logOutput(ctx, "client", attr.journalKeyValues(), attr.TagFunc)
		}(r.Context())

		return next.RoundTrip(r)
	})
}
