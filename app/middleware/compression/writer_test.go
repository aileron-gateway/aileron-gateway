package compression

import (
	"maps"
	"net/http"
	"net/http/httptest"
	"testing"

	"io"

	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
)

// mockResettableWriter is a mock resettableWriter for testing.
type mockResettableWriter struct {
	writer io.Writer
	data   []byte
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

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"response body too small/skip compression",
			[]string{},
			[]string{},
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
			[]string{},
			[]string{},
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
			[]string{},
			[]string{},
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
			[]string{},
			[]string{},
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
			[]string{},
			[]string{},
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
			[]string{},
			[]string{},
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
			[]string{},
			[]string{},
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
			"compress response body",
			[]string{},
			[]string{},
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
			[]string{},
			[]string{},
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

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			rec := httptest.NewRecorder()
			maps.Copy(rec.Header(), tt.C().header)
			cw := &compressionWriter{
				ResponseWriter: rec,
				writer:         &mockResettableWriter{},
				encoding:       tt.C().encoding,
				mimes:          tt.C().mimes,
				minimumSize:    tt.C().minimumSize,
			}

			if tt.C().status != 0 {
				cw.WriteHeader(tt.C().status)
			}
			n, err := cw.Write(tt.C().data)
			testutil.DiffError(t, nil, nil, err)
			testutil.Diff(t, len(tt.C().data), n)
			b, _ := io.ReadAll(rec.Result().Body)
			testutil.Diff(t, string(tt.C().data), string(b))

			testutil.Diff(t, tt.A().initialized, cw.initialized)
			testutil.Diff(t, tt.A().shouldSkip, cw.shouldSkip)
			testutil.Diff(t, tt.A().encoding, rec.Header().Get("Content-Encoding"))
		})
	}
}
