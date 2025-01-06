package http

import (
	"net/http"
	"testing"

	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
	"github.com/google/go-cmp/cmp"
)

type testResponseWriter struct {
	http.ResponseWriter
	id string
}

func (w *testResponseWriter) Unwrap() http.ResponseWriter {
	return w.ResponseWriter
}

type testFlushResponseWriter struct {
	http.ResponseWriter
	id      string
	flushed bool
}

func (w *testFlushResponseWriter) Flush() {
	w.flushed = true
}

func TestWrappedWriter_Flush(t *testing.T) {
	type condition struct {
		w http.ResponseWriter
	}

	type action struct {
		w http.ResponseWriter
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"non flusher",
			[]string{},
			[]string{},
			&condition{
				w: &testResponseWriter{id: "test"},
			},
			&action{
				w: &testResponseWriter{id: "test"},
			},
		),
		gen(
			"flusher",
			[]string{},
			[]string{},
			&condition{
				w: &testFlushResponseWriter{id: "test"},
			},
			&action{
				w: &testFlushResponseWriter{id: "test", flushed: true},
			},
		),
		gen(
			"non inner flusher",
			[]string{},
			[]string{},
			&condition{
				w: &testResponseWriter{
					id:             "out",
					ResponseWriter: &testResponseWriter{id: "inner"},
				},
			},
			&action{
				w: &testResponseWriter{
					id:             "out",
					ResponseWriter: &testResponseWriter{id: "inner"},
				},
			},
		),
		gen(
			"non inner flusher",
			[]string{},
			[]string{},
			&condition{
				w: &testResponseWriter{
					id:             "out",
					ResponseWriter: &testFlushResponseWriter{id: "inner"},
				},
			},
			&action{
				w: &testResponseWriter{
					id:             "out",
					ResponseWriter: &testFlushResponseWriter{id: "inner", flushed: true},
				},
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			w := &WrappedWriter{
				ResponseWriter: tt.C().w,
			}
			w.Flush()

			testutil.Diff(t, true, w.flushChecked)
			testutil.Diff(t, true, w.flushFunc != nil)

			opts := []cmp.Option{
				cmp.AllowUnexported(testResponseWriter{}),
				cmp.AllowUnexported(testFlushResponseWriter{}),
				// cmp.Comparer(testutil.ComparePointer[foo.Bar])
				// testutil.Po
			}
			testutil.Diff(t, tt.A().w, tt.C().w, opts...)

			w.Flush()
			testutil.Diff(t, true, w.flushChecked)
			testutil.Diff(t, true, w.flushFunc != nil)
			testutil.Diff(t, tt.A().w, tt.C().w, opts...)
		})
	}
}
