package http

import (
	"net/http"
)

// noopFunc is the function that do nothing.
// This is used as a dummy function like io.NooCloser func.
var noopFunc = func() {}

// Writer is the response writer that
// can obtain written status code and written bytes.
type Writer interface {
	http.ResponseWriter

	// Written returns if status code or body were written or not.
	// Calling Write(nil) is considered to be written even it does not
	// write any bytes.
	// If ContentLength() > 0, Written() always returns true.
	// If ContentLength() == 0, Written() can be true or false.
	Written() bool

	// StatusCode returns the written status code.
	// If both status code and body were not written at all,
	// StatusCode() returns zero.
	// If body was written without writing status code,
	// StatusCode() returns 200 which is the same behavior with
	// [http.ResponseWriter.WriteHeader] method.
	// See the comment of https://pkg.go.dev/net/http#ResponseWriter.WriteHeader
	StatusCode() int

	// ContentLength returns the actual written bytes.
	// If ContentLength() > 0, Written() always returns true.
	// If ContentLength() == 0, Written() can be true or false.
	ContentLength() int64
}

// WrapWriter wraps http.ResponseWriter with a writer that
// can record written status code and written bytes.
// When nil is given as the argument, this returns a new io.Writer with wrapping nil http.ResponseWriter
// and the returned writer will panic on write.
// Actual implementation of returned Writer is *WrappedWriter.
func WrapWriter(w http.ResponseWriter) Writer {
	if ww, ok := w.(Writer); ok {
		return ww
	}
	return &WrappedWriter{
		ResponseWriter: w,
	}
}

// WrappedWriter wraps http.ResponseWriter and holds http status code.
// This implements io.Writer interface.
type WrappedWriter struct {
	http.ResponseWriter
	code    int
	written bool
	length  int64

	flushChecked bool
	flushFunc    func()
}

func (w *WrappedWriter) Unwrap() http.ResponseWriter {
	return w.ResponseWriter
}

func (w *WrappedWriter) Flush() {
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

func (w *WrappedWriter) WriteHeader(statusCode int) {
	w.written = true
	w.code = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *WrappedWriter) Write(b []byte) (int, error) {
	w.written = true
	n, err := w.ResponseWriter.Write(b)
	w.length += int64(n)
	return n, err
}

func (w *WrappedWriter) Written() bool {
	return w.written
}

func (w *WrappedWriter) StatusCode() int {
	if w.written && w.code == 0 {
		return http.StatusOK
	}
	return w.code
}

func (w *WrappedWriter) ContentLength() int64 {
	return w.length
}
