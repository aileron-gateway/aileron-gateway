package compression

import (
	"compress/gzip"
	"io"
	"mime"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"sync"

	"github.com/andybalholm/brotli"
)

// newBrotliWriterPool returns a new instance of
// brotli.Writer pool.
// Panics occurs when an invalid level was given as an argument.
func newBrotliWriterPool(level int) sync.Pool {
	return sync.Pool{
		New: func() any {
			w := brotli.NewWriterLevel(io.Discard, level)
			return w
		},
	}
}

// newGzipWriterPool returns a new instance of
// gzip.Writer pool.
// Given level must be a valid value otherwise
// the pool returns nil writer.
func newGzipWriterPool(level int) sync.Pool {
	return sync.Pool{
		New: func() any {
			// No errors are expected because we know the
			// level is valid value.
			w, _ := gzip.NewWriterLevel(io.Discard, level)
			return w
		},
	}
}

// resettableWriter is a writer that can
// reset the internal io.Writer.
// Both gzip.Writer and brotli.Writer satisfy
// this resettableWriter interface.
type resettableWriter interface {
	io.WriteCloser
	Reset(io.Writer)
}

// compressionWriter compress the response data written to the embedded http.ResponseWriter.
type compressionWriter struct {
	http.ResponseWriter

	// writer is a gzip writer.
	// This writer will be initialized with the embedded http.ResponseWriter
	// when initializing. It is initialized by io.Discard before initialize method is called.
	writer resettableWriter

	// encoding is the compression type string
	// that should be go in the Content-Encoding
	// response header. Currently "gzip" or "br".
	encoding string

	// mimes is the list of MIME types that should be compressed.
	// Compression is skipped when the Content-Type of the response is included in this list.
	// Checkout the link below for all official MIME types.
	// 	- https://www.iana.org/assignments/media-types/media-types.xhtml
	mimes []string

	// minimumSize is the minimum byte size of response body to be compressed.
	// Response bodies that are smaller than this value won't be compressed.
	// Note that compressing small object (~1kB) could be larger than the original size.
	minimumSize int64

	// initialized is the flag to hold if this writer was initialized or not.
	// This flag is set to true when the initialize() method was called.
	initialized bool

	// shouldSkip is the flag to skip compression.
	// This flag will be initialized in initialize().
	shouldSkip bool
}

// initialize initializes this writer and wrapped HTTP response.
func (w *compressionWriter) initialize() {
	w.initialized = true
	wh := w.ResponseWriter.Header()

	length := wh.Get("Content-Length")
	if length == "" {
		w.shouldSkip = true // Unknown response body size.
		return
	}

	if size, _ := strconv.ParseInt(length, 10, 64); size < w.minimumSize {
		w.shouldSkip = true // Response body too small.
		return
	}

	m, _, _ := mime.ParseMediaType(wh.Get("Content-Type"))
	if !slices.Contains(w.mimes, m) {
		w.shouldSkip = true // Not a target mime.
		return
	}

	// Content-Encoding can contain multiple values.
	// So, skip compression if there is at least one "gzip" or "br".
	ce := wh.Get("Content-Encoding")
	if ce != "" {
		// Skip if already compressed.
		if strings.Contains(ce, gzipEncoding) || strings.Contains(ce, brotliEncoding) ||
			strings.Contains(ce, "deflate") || strings.Contains(ce, "compress") || strings.Contains(ce, "zstd") {
			w.shouldSkip = true
			return
		}
		ce = ce + "," + w.encoding
	} else {
		ce = w.encoding
	}

	wh.Add("Vary", "Accept-Encoding") // Don't let client use cached content when required for different types of encoding.
	wh.Set("Content-Encoding", ce)    // Set the encoding by replacing existing one.
	wh.Del("Content-Length")          // Delete Content-Length because the content will be compressed.
	w.writer.Reset(w.ResponseWriter)
}

// WriteHeader writes response status to the response write.
// Response header and compression writer will be initialized
// before writing the status code into the internal response writer.
func (w *compressionWriter) WriteHeader(statusCode int) {
	if !w.initialized {
		w.initialize() // Initialize writer and response headers.
	}
	w.ResponseWriter.WriteHeader(statusCode)
}

// Write write data into this
// Calling this method initialize this compressionWriter.
// Written data will be compressed when skip=false.
func (w *compressionWriter) Write(data []byte) (int, error) {
	// Do not write nil data because compressing empty byte won't be 0 bytes.
	// For example, compressing empty string "" becomes "H4sIAAAAAAAAAwMAAAAAAAAAAAA=".
	if len(data) == 0 {
		return 0, nil
	}

	if !w.initialized {
		w.initialize() // Initialize writer and response headers.
	}
	if w.shouldSkip {
		return w.ResponseWriter.Write(data)
	}
	return w.writer.Write(data)
}
