package httpproxy

import (
	"io"
	"mime"
	"net/http"

	"golang.org/x/net/http/httpguts"
)

// shouldFlushImmediately returns flush interval.
// Negative interval means immediate flush.
// For content-tye "text/event-stream", which means server-sent events
// we have to call Flush() immediately after received response.
// See the https://www.w3.org/TR/eventsource/#text-event-stream
func shouldFlushImmediately(res *http.Response) bool {
	if res.ContentLength < 0 {
		// Any streaming type response.
		// Or almost all HTTP/2, HTTP/3.
		return true
	}
	if mt, _, _ := mime.ParseMediaType(res.Header.Get("Content-Type")); mt == "text/event-stream" {
		return true // Server Sent Events (SSE)
	}
	if httpguts.HeaderValuesContainsToken(res.Header["Transfer-Encoding"], "chunked") {
		return true // Chunked response.
	}
	return false
}

// withImmediateFlush wraps the given http.ResponseWriter
// with immediate flushing writer.
// Immediate flushing won't be applied when the given http.ResponseWriter
// does not implement http.Flusher interface or the flushImmediately is false.
func withImmediateFlush(rw http.ResponseWriter, flushImmediately bool) io.Writer {
	if !flushImmediately {
		return rw
	}
	inner := rw
	for {
		if flusher, ok := inner.(http.Flusher); ok {
			return &immediateFlushWriter{
				inner:   inner,
				flusher: flusher,
			}
		}
		if uw, ok := inner.(interface{ Unwrap() http.ResponseWriter }); ok {
			inner = uw.Unwrap()
			continue
		}
		return rw
	}
}

// immediateFlushWriter flushed internal flusher immediately
// after writing to the writer.
// immediateFlushWriter implements io.Writer interface.
type immediateFlushWriter struct {
	inner io.Writer
	// flusher flushes buffer.
	// flusher must not be nil.
	flusher http.Flusher
}

func (f *immediateFlushWriter) Write(p []byte) (n int, err error) {
	defer f.flusher.Flush()
	return f.inner.Write(p)
}
