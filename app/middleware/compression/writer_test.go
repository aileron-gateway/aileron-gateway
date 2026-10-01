// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package compression

import (
	"bytes"
	"compress/gzip"
	"maps"
	"net/http"
	"net/http/httptest"
	"testing"

	"io"

	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
	"github.com/google/go-cmp/cmp/cmpopts"
)

// mockResettableWriter is a mock resettableWriter for testing.
type mockResettableWriter struct {
	writer   io.Writer
	data     []byte
	flushed  bool
	flushErr error
}

func (m *mockResettableWriter) Write(p []byte) (n int, err error) {
	m.data = append(m.data, p...)
	if m.writer != nil {
		return m.writer.Write(p)
	}
	return len(p), nil
}

func (m *mockResettableWriter) Close() error {
	return nil
}

func (m *mockResettableWriter) Reset(w io.Writer) {
	m.writer = w
}

func (m *mockResettableWriter) Flush() error {
	m.flushed = true
	return m.flushErr
}

func TestCompressionWriter(t *testing.T) {
	type condition struct {
		// contentType      string
		// contentLength    string
		// acceptEncoding   string
		// existingEncoding string
		encoding    string
		minimumSize int64
		mimes       []string

		header http.Header
		status int
		data   []byte
	}

	type action struct {
		initialized bool
		shouldSkip  bool
		encoding    string
	}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"response body too small/skip compression",
			&condition{
				encoding:    "gzip",
				minimumSize: 1024,
				mimes:       []string{"text/html", "application/json"},
				header:      http.Header{"Content-Type": {"text/html"}, "Content-Length": {"512"}},
				data:        []byte("response body"),
			},
			&action{
				initialized: true,
				shouldSkip:  true,
				encoding:    "",
			},
		),
		gen(
			"response body large enough, apply gzip compression",
			&condition{
				encoding:    "gzip",
				minimumSize: 1024,
				mimes:       []string{"text/html", "application/json"},
				header:      http.Header{"Content-Type": {"text/html"}, "Content-Length": {"2048"}},
				data:        []byte("response body"),
			},
			&action{
				initialized: true,
				shouldSkip:  false,
				encoding:    "gzip",
			},
		),
		gen(
			"target MIME type/compress",
			&condition{
				encoding:    "gzip",
				minimumSize: 1024,
				mimes:       []string{"text/html", "application/json"},
				header:      http.Header{"Content-Type": {"text/html"}, "Content-Length": {"2048"}},
				data:        []byte("response body"),
			},
			&action{
				initialized: true,
				shouldSkip:  false,
				encoding:    "gzip",
			},
		),
		gen(
			"non target MIME type, skip compression",
			&condition{
				encoding:    "gzip",
				minimumSize: 1024,
				mimes:       []string{"text/html", "application/json"},
				header:      http.Header{"Content-Type": {"image/png"}, "Content-Length": {"2048"}},
				data:        []byte("response body"),
			},
			&action{
				initialized: true,
				shouldSkip:  true,
				encoding:    "",
			},
		),
		gen(
			"already compressed with gzip",
			&condition{
				minimumSize: 1024,
				mimes:       []string{"text/html", "application/json"},
				header: http.Header{
					"Content-Type":     {"text/html"},
					"Content-Length":   {"2048"},
					"Content-Encoding": {"gzip"},
				},
				data: []byte("response body"),
			},
			&action{
				initialized: true,
				shouldSkip:  true,
				encoding:    "gzip",
			},
		),
		gen(
			"already compressed with brotli",
			&condition{
				minimumSize: 1024,
				mimes:       []string{"text/html", "application/json"},
				header: http.Header{
					"Content-Type":     {"text/html"},
					"Content-Length":   {"2048"},
					"Content-Encoding": {"br"},
				},
				data: []byte("response body"),
			},
			&action{
				initialized: true,
				shouldSkip:  true,
				encoding:    "br",
			},
		),
		gen(
			"already compressed with deflate",
			&condition{
				minimumSize: 1024,
				mimes:       []string{"text/html", "application/json"},
				header: http.Header{
					"Content-Type":     {"text/html"},
					"Content-Length":   {"2048"},
					"Content-Encoding": {"deflate"},
				},
				data: []byte("response body"),
			},
			&action{
				initialized: true,
				shouldSkip:  true,
				encoding:    "deflate",
			},
		),
		gen(
			"status code with no body",
			&condition{
				minimumSize: 1024,
				status:      204,
				mimes:       []string{"text/html"},
				header: http.Header{
					"Content-Type":     {},
					"Content-Length":   {},
					"Content-Encoding": {"dummy"},
				},
				encoding: "gzip",
				data:     nil,
			},
			&action{
				initialized: true,
				shouldSkip:  true,
				encoding:    "", // should be deleted
			},
		),
		gen(
			"compress response body",
			&condition{
				minimumSize: 1024,
				mimes:       []string{"text/html", "application/json"},
				header: http.Header{
					"Content-Type":     {"text/html"},
					"Content-Length":   {"2048"},
					"Content-Encoding": {"unknown"},
				},
				encoding: "gzip",
				data:     []byte("response body"),
			},
			&action{
				initialized: true,
				shouldSkip:  false,
				encoding:    "unknown,gzip",
			},
		),
		gen(
			"write empty body",
			&condition{
				minimumSize: 1024,
				mimes:       []string{"text/html", "application/json"},
				header:      http.Header{"Content-Type": {"text/html"}, "Content-Length": {"2048"}},
				encoding:    "gzip",
				data:        []byte(nil),
			},
			&action{
				initialized: false,
				shouldSkip:  false,
				encoding:    "",
			},
		),
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			maps.Copy(rec.Header(), tt.C.header)
			cw := &compressionWriter{
				ResponseWriter: rec,
				flush:          flushFunc(rec),
				writer:         &mockResettableWriter{},
				encoding:       tt.C.encoding,
				mimes:          tt.C.mimes,
				minimumSize:    tt.C.minimumSize,
			}

			if tt.C.status != 0 {
				cw.WriteHeader(tt.C.status)
			}
			n, err := cw.Write(tt.C.data)
			testutil.DiffError(t, nil, nil, err)
			testutil.Diff(t, len(tt.C.data), n)
			b, _ := io.ReadAll(rec.Result().Body)
			testutil.Diff(t, string(tt.C.data), string(b))

			testutil.Diff(t, tt.A.initialized, cw.initialized)
			testutil.Diff(t, tt.A.shouldSkip, cw.shouldSkip)
			testutil.Diff(t, tt.A.encoding, rec.Header().Get("Content-Encoding"))
		})
	}
}

func TestCompressionWriter_FlushError(t *testing.T) {
	t.Parallel()
	t.Run("non-nil flush", func(t *testing.T) {
		var buf bytes.Buffer
		gw := gzip.NewWriter(&buf)
		flushed := 0
		cw := &compressionWriter{
			ResponseWriter: httptest.NewRecorder(),
			flush:          func() error { flushed += 1; return nil },
			writer:         gw,
			initialized:    true,
		}
		cw.Write([]byte("foo"))
		err := cw.FlushError()
		testutil.Diff(t, nil, err)
		testutil.Diff(t, 1, flushed)
		cw.Write([]byte("bar"))
		err = cw.FlushError()
		testutil.Diff(t, nil, err)
		testutil.Diff(t, 2, flushed)
		gw.Close()
		r, _ := gzip.NewReader(&buf)
		b, _ := io.ReadAll(r)
		testutil.Diff(t, "foobar", string(b))
	})
	t.Run("nil flush", func(t *testing.T) {
		var buf bytes.Buffer
		gw := gzip.NewWriter(&buf)
		cw := &compressionWriter{
			ResponseWriter: httptest.NewRecorder(),
			flush:          nil,
			writer:         gw,
			initialized:    true,
		}
		cw.Write([]byte("foo"))
		err := cw.FlushError()
		testutil.Diff(t, nil, err)
		cw.Write([]byte("bar"))
		err = cw.FlushError()
		testutil.Diff(t, nil, err)
		gw.Close()
		r, _ := gzip.NewReader(&buf)
		b, _ := io.ReadAll(r)
		testutil.Diff(t, "foobar", string(b))
	})
	t.Run("flush error", func(t *testing.T) {
		var buf bytes.Buffer
		gw := gzip.NewWriter(&buf)
		cw := &compressionWriter{
			ResponseWriter: httptest.NewRecorder(),
			flush:          nil,
			writer:         &mockResettableWriter{flushErr: io.EOF},
			initialized:    true,
		}
		cw.Write([]byte("foo"))
		gw.Close()
		err := cw.FlushError()
		testutil.Diff(t, io.EOF, err, cmpopts.EquateErrors())
	})
}

func TestCompressionWriter_Flush(t *testing.T) {
	t.Parallel()
	t.Run("non-nil flush", func(t *testing.T) {
		var buf bytes.Buffer
		gw := gzip.NewWriter(&buf)
		flushed := 0
		cw := &compressionWriter{
			ResponseWriter: httptest.NewRecorder(),
			flush:          func() error { flushed += 1; return nil },
			writer:         gw,
			initialized:    true,
		}
		cw.Write([]byte("foo"))
		cw.Flush()
		testutil.Diff(t, 1, flushed)
		cw.Write([]byte("bar"))
		cw.Flush()
		testutil.Diff(t, 2, flushed)
		gw.Close()
		r, _ := gzip.NewReader(&buf)
		b, _ := io.ReadAll(r)
		testutil.Diff(t, "foobar", string(b))
	})
	t.Run("nil flush", func(t *testing.T) {
		var buf bytes.Buffer
		gw := gzip.NewWriter(&buf)
		cw := &compressionWriter{
			ResponseWriter: httptest.NewRecorder(),
			flush:          nil,
			writer:         gw,
			initialized:    true,
		}
		cw.Write([]byte("foo"))
		cw.Flush()
		cw.Write([]byte("bar"))
		cw.Flush()
		gw.Close()
		r, _ := gzip.NewReader(&buf)
		b, _ := io.ReadAll(r)
		testutil.Diff(t, "foobar", string(b))
	})
}

type mockFlushError struct {
	http.ResponseWriter
	flushed  bool
	flushErr error
}

func (m *mockFlushError) FlushError() error {
	m.flushed = true
	return m.flushErr
}

type mockFlush struct {
	http.ResponseWriter
	flushed  bool
	flushErr error
}

func (m *mockFlush) Flush() error {
	m.flushed = true
	return m.flushErr
}

type mockHTTPFlush struct {
	http.ResponseWriter
	flushed bool
}

func (m *mockHTTPFlush) Flush() {
	m.flushed = true
}

type mockUnwrapWriter struct {
	http.ResponseWriter
}

func (m *mockUnwrapWriter) Unwrap() http.ResponseWriter {
	return m.ResponseWriter
}

func TestFlushFunc(t *testing.T) {
	t.Parallel()
	t.Run("FlushError", func(t *testing.T) {
		rw := &mockFlushError{}
		f := flushFunc(rw)
		f()
		testutil.Diff(t, true, rw.flushed)
	})
	t.Run("Flush", func(t *testing.T) {
		rw := &mockFlush{}
		f := flushFunc(rw)
		f()
		testutil.Diff(t, true, rw.flushed)
	})
	t.Run("HTTPFlush", func(t *testing.T) {
		rw := &mockHTTPFlush{}
		f := flushFunc(rw)
		f()
		testutil.Diff(t, true, rw.flushed)
	})
	t.Run("unwrap", func(t *testing.T) {
		rw := &mockHTTPFlush{}
		f := flushFunc(&mockUnwrapWriter{rw})
		f()
		testutil.Diff(t, true, rw.flushed)
	})
	t.Run("unwrap nil", func(t *testing.T) {
		rw := &mockUnwrapWriter{&mockUnwrapWriter{}}
		f := flushFunc(rw)
		testutil.Diff(t, true, f == nil)
	})
}
