package httpproxy

import (
	"io"
	"net/http"
	"testing"

	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
	utilhttp "github.com/aileron-gateway/aileron-gateway/util/http"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

type testWriter struct {
	err      error
	written  []byte
	readOnly int
}

func (w *testWriter) Write(p []byte) (n int, err error) {
	if w.readOnly > 0 {
		w.written = append(w.written, p[:w.readOnly]...)
		return w.readOnly, w.err
	}
	w.written = append(w.written, p...)
	return len(p), w.err
}

type testWriteFlusher struct {
	*testWriter
	flushed bool
}

func (wf *testWriteFlusher) Flush() {
	wf.flushed = true
}

type testResponse struct {
	id     string
	status int
	writer io.Writer
	header http.Header
}

func (r *testResponse) Header() http.Header {
	return r.header
}

func (r *testResponse) WriteHeader(statusCode int) {
	r.status = statusCode
}

func (r *testResponse) Write(p []byte) (int, error) {
	return r.writer.Write(p)
}

type testFlushResponse struct {
	*testResponse
	flushed bool
}

func (r *testFlushResponse) Flush() {
	r.flushed = true
}

func TestWithImmediateFlush(t *testing.T) {
	type condition struct {
		w     http.ResponseWriter
		flush bool
	}

	type action struct {
		w io.Writer
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"no flush",
			[]string{},
			[]string{},
			&condition{
				w: &testResponse{
					id: "test",
				},
				flush: false,
			},
			&action{
				w: &testResponse{
					id: "test",
				},
			},
		),
		gen(
			"flush with no-flusher",
			[]string{},
			[]string{},
			&condition{
				w: &testResponse{
					id: "test",
				},
				flush: true,
			},
			&action{
				w: &testResponse{
					id: "test",
				},
			},
		),
		gen(
			"flush with flusher",
			[]string{},
			[]string{},
			&condition{
				w: &testFlushResponse{
					testResponse: &testResponse{
						id: "test",
					},
				},
				flush: true,
			},
			&action{
				w: &immediateFlushWriter{
					inner: &testFlushResponse{
						testResponse: &testResponse{
							id: "test",
						},
					},
					flusher: &testFlushResponse{
						testResponse: &testResponse{
							id: "test",
						},
					},
				},
			},
		),
		gen(
			"flush with wrapped flusher",
			[]string{},
			[]string{},
			&condition{
				w: utilhttp.WrapWriter(
					&testFlushResponse{
						testResponse: &testResponse{
							id: "test",
						},
					},
				),
				flush: true,
			},
			&action{
				w: &immediateFlushWriter{
					inner: &utilhttp.WrappedWriter{
						ResponseWriter: &testFlushResponse{
							testResponse: &testResponse{
								id: "test",
							},
						},
					},
					flusher: &utilhttp.WrappedWriter{
						ResponseWriter: &testFlushResponse{
							testResponse: &testResponse{
								id: "test",
							},
						},
					},
				},
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			w := withImmediateFlush(tt.C().w, tt.C().flush)

			opts := []cmp.Option{
				cmp.AllowUnexported(immediateFlushWriter{}),
				cmp.AllowUnexported(testWriteFlusher{}, testWriter{}),
				cmp.AllowUnexported(testResponse{}, testFlushResponse{}),
				cmpopts.IgnoreUnexported(utilhttp.WrappedWriter{}),
			}
			testutil.Diff(t, tt.A().w, w, opts...)
		})
	}
}

func TestShouldFlushImmediately(t *testing.T) {
	type condition struct {
		res *http.Response
	}

	type action struct {
		flush bool
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"no immediate flush",
			[]string{},
			[]string{},
			&condition{
				res: &http.Response{},
			},
			&action{
				flush: false,
			},
		),
		gen(
			"no flush content-type",
			[]string{},
			[]string{},
			&condition{
				res: &http.Response{
					Header: http.Header{
						"Content-Type": {"application/octet-stream"},
					},
				},
			},
			&action{
				flush: false,
			},
		),
		gen(
			"flush event-stream",
			[]string{},
			[]string{},
			&condition{
				res: &http.Response{
					Header: http.Header{
						"Content-Type": {"text/event-stream"},
					},
				},
			},
			&action{
				flush: true,
			},
		),
		gen(
			"positive content length",
			[]string{},
			[]string{},
			&condition{
				res: &http.Response{
					ContentLength: 100,
				},
			},
			&action{
				flush: false,
			},
		),
		gen(
			"flush event-stream",
			[]string{},
			[]string{},
			&condition{
				res: &http.Response{
					ContentLength: -1,
				},
			},
			&action{
				flush: true,
			},
		),
		gen(
			"chunked",
			[]string{},
			[]string{},
			&condition{
				res: &http.Response{
					ContentLength: 0,
					Header: http.Header{
						"Content-Type":      {"text/plain"},
						"Transfer-Encoding": {"gzip, chunked"},
					},
				},
			},
			&action{
				flush: true,
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			shouldFlush := shouldFlushImmediately(tt.C().res)
			testutil.Diff(t, tt.A().flush, shouldFlush)
		})
	}
}

func TestImmediateFlushWriter(t *testing.T) {
	type condition struct {
		wf    *testWriteFlusher
		write [][]byte
	}

	type action struct {
		written []byte
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"immediate flush for nil",
			[]string{},
			[]string{},
			&condition{
				wf: &testWriteFlusher{
					testWriter: &testWriter{},
				},
				write: [][]byte{nil},
			},
			&action{
				written: nil,
			},
		),
		gen(
			"immediate flush",
			[]string{},
			[]string{},
			&condition{
				wf: &testWriteFlusher{
					testWriter: &testWriter{},
				},
				write: [][]byte{
					[]byte("foo"),
					[]byte("bar"),
				},
			},
			&action{
				written: []byte("foobar"),
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			wf := &immediateFlushWriter{
				inner:   tt.C().wf,
				flusher: tt.C().wf,
			}

			for _, b := range tt.C().write {
				tt.C().wf.flushed = false
				wf.Write(b)
				testutil.Diff(t, true, tt.C().wf.flushed)
			}
			testutil.Diff(t, tt.A().written, tt.C().wf.written)
		})
	}
}
