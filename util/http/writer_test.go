package http_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
	utilhttp "github.com/aileron-gateway/aileron-gateway/util/http"
	"github.com/google/go-cmp/cmp"
)

func TestWrapWriter(t *testing.T) {
	type condition struct {
		w          http.ResponseWriter
		statusCode int
	}

	type action struct {
		written    bool
		statusCode int
	}

	CndNilResponseWriter := "nil response writer"
	CndWriteStatusCode := "write status code"

	ActCheckWritten := "check if the code was written"
	ActCheckStatus := "check the status code"

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	tb.Condition(CndNilResponseWriter, "give nil http.ResponseWriter as an argument")
	tb.Condition(CndWriteStatusCode, "writer status code to the created io.Writer")

	tb.Action(ActCheckWritten, "check that a status code was written")
	tb.Action(ActCheckStatus, "check the status code")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"new with ResponseWriter without writing status code",
			[]string{},
			[]string{ActCheckWritten, ActCheckStatus},
			&condition{
				w: httptest.NewRecorder(),
			},
			&action{
				written:    false,
				statusCode: 0,
			},
		),
		gen(
			"nil ResponseWriter",
			[]string{CndNilResponseWriter},
			[]string{ActCheckWritten, ActCheckStatus},
			&condition{
				w: nil,
			},
			&action{
				written:    false,
				statusCode: 0,
			},
		),
		gen(
			"new with ResponseWriter with writing status code",
			[]string{CndWriteStatusCode},
			[]string{ActCheckWritten, ActCheckStatus},
			&condition{
				w:          httptest.NewRecorder(),
				statusCode: 999,
			},
			&action{
				written:    true,
				statusCode: 999,
			},
		),
		gen(
			"new with ResponseWriter with writing status code",
			[]string{CndWriteStatusCode},
			[]string{ActCheckWritten, ActCheckStatus},
			&condition{
				w: &utilhttp.WrappedWriter{
					ResponseWriter: httptest.NewRecorder(),
				},
				statusCode: 999,
			},
			&action{
				written:    true,
				statusCode: 999,
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			ww := utilhttp.WrapWriter(tt.C().w)
			if tt.C().statusCode > 0 {
				ww.WriteHeader(tt.C().statusCode)
			}

			testutil.Diff(t, tt.A().written, ww.Written())
			testutil.Diff(t, tt.A().statusCode, ww.StatusCode())

			wwPtr := ww.(*utilhttp.WrappedWriter)
			target := tt.C().w
			if ww, ok := tt.C().w.(*utilhttp.WrappedWriter); ok {
				target = ww.ResponseWriter
			}
			testutil.Diff(t, target, wwPtr.ResponseWriter, cmp.Comparer(testutil.ComparePointer[http.ResponseWriter]))
		})
	}
}

type testResponseWriter struct {
	http.ResponseWriter
	id string
}

func TestWrappedWriter_Unwrap(t *testing.T) {
	type condition struct {
		ww *utilhttp.WrappedWriter
	}

	type action struct {
		w http.ResponseWriter
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	cndNonNil := tb.Condition("non-nil", "inner writer is non-nil")
	actCheckWriter := tb.Action("check writer", "check the inner writer")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"unwrap nil",
			[]string{},
			[]string{actCheckWriter},
			&condition{
				ww: &utilhttp.WrappedWriter{
					ResponseWriter: nil,
				},
			},
			&action{
				w: nil,
			},
		),
		gen(
			"unwrap non-nil",
			[]string{cndNonNil},
			[]string{actCheckWriter},
			&condition{
				ww: &utilhttp.WrappedWriter{
					ResponseWriter: &testResponseWriter{
						id: "inner",
					},
				},
			},
			&action{
				w: &testResponseWriter{
					id: "inner",
				},
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			w := tt.C().ww.Unwrap()
			testutil.Diff(t, tt.A().w, w, cmp.AllowUnexported(testResponseWriter{}))
		})
	}
}

func TestWrappedWriter_WriteHeader(t *testing.T) {
	type condition struct {
		ww         *utilhttp.WrappedWriter
		statusCode int
	}

	type action struct {
		statusCode int
	}

	CndWriteStatusCode := "write status code"
	ActCheckStatus := "check the status code"

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	tb.Condition(CndWriteStatusCode, "write status code")
	tb.Action(ActCheckStatus, "check the status code")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"status code 100",
			[]string{CndWriteStatusCode},
			[]string{ActCheckStatus},
			&condition{
				ww: &utilhttp.WrappedWriter{
					ResponseWriter: httptest.NewRecorder(),
				},
				statusCode: 100,
			},
			&action{
				statusCode: 100,
			},
		),
		gen(
			"status code 999",
			[]string{CndWriteStatusCode},
			[]string{ActCheckStatus},
			&condition{
				ww: &utilhttp.WrappedWriter{
					ResponseWriter: httptest.NewRecorder(),
				},
				statusCode: 999,
			},
			&action{
				statusCode: 999,
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			w := httptest.NewRecorder()
			ww := utilhttp.WrapWriter(w)
			ww.WriteHeader(tt.C().statusCode)

			testutil.Diff(t, true, ww.Written())
			testutil.Diff(t, tt.A().statusCode, ww.StatusCode())
			testutil.Diff(t, tt.A().statusCode, w.Result().StatusCode)
		})
	}
}

func TestWrappedWriter_Write(t *testing.T) {
	type condition struct {
		statusCode int
		body       string
	}

	type action struct {
		statusCode int
		body       string
	}

	CndWriteStatusCode := "write status code"
	ActCheckStatus := "check the status code"
	ActCheckBody := "check the written body"

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	tb.Condition(CndWriteStatusCode, "write status code")
	tb.Action(ActCheckStatus, "check the status code")
	tb.Action(ActCheckBody, "check the written body was the one expected")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"status code 100",
			[]string{CndWriteStatusCode},
			[]string{ActCheckStatus, ActCheckBody},
			&condition{
				statusCode: 100,
				body:       "test",
			},
			&action{
				statusCode: 100,
				body:       "test",
			},
		),
		gen(
			"status code 999",
			[]string{CndWriteStatusCode},
			[]string{ActCheckStatus, ActCheckBody},
			&condition{
				statusCode: 999,
				body:       "test",
			},
			&action{
				statusCode: 999,
				body:       "test",
			},
		),
		gen(
			"status code 0 (don't write the code)",
			[]string{},
			[]string{ActCheckStatus, ActCheckBody},
			&condition{
				statusCode: 0,
				body:       "test",
			},
			&action{
				statusCode: 200,
				body:       "test",
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			w := httptest.NewRecorder()
			ww := utilhttp.WrapWriter(w)
			if tt.C().statusCode > 0 {
				ww.WriteHeader(tt.C().statusCode)
			}
			ww.Write([]byte(tt.C().body))

			testutil.Diff(t, true, ww.Written())
			testutil.Diff(t, tt.A().statusCode, ww.StatusCode())
			testutil.Diff(t, w.Result().StatusCode, ww.StatusCode())

			body, _ := io.ReadAll(w.Body)
			testutil.Diff(t, tt.A().body, string(body))
			testutil.Diff(t, len(tt.A().body), int(ww.ContentLength()))
		})
	}
}

func TestWrappedWriter_Written(t *testing.T) {
	type condition struct {
		ww    *utilhttp.WrappedWriter
		write bool
	}

	type action struct {
		written bool
	}

	CndWriteStatusCode := "write status code"
	ActCheckWritten := "check if written"

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	tb.Condition(CndWriteStatusCode, "writer status code")
	tb.Action(ActCheckWritten, "check that the status code was written or not")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"don't write status code",
			[]string{CndWriteStatusCode},
			[]string{ActCheckWritten},
			&condition{
				ww: &utilhttp.WrappedWriter{
					ResponseWriter: httptest.NewRecorder(),
				},
				write: false,
			},
			&action{
				written: false,
			},
		),
		gen(
			"write status code",
			[]string{CndWriteStatusCode},
			[]string{ActCheckWritten},
			&condition{
				ww: &utilhttp.WrappedWriter{
					ResponseWriter: httptest.NewRecorder(),
				},
				write: true,
			},
			&action{
				written: true,
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			if tt.C().write {
				tt.C().ww.WriteHeader(999)
			}

			testutil.Diff(t, tt.A().written, tt.C().ww.Written())
		})
	}
}

func TestWrappedWriter_StatusCode(t *testing.T) {
	type condition struct {
		ww         *utilhttp.WrappedWriter
		statusCode int
	}

	type action struct {
		statusCode int
	}

	CndWriteStatusCode := "write status code"
	ActCheckStatus := "check the status code"

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	tb.Condition(CndWriteStatusCode, "writer status code")
	tb.Action(ActCheckStatus, "check the status code")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"status code 100",
			[]string{CndWriteStatusCode},
			[]string{ActCheckStatus},
			&condition{
				ww: &utilhttp.WrappedWriter{
					ResponseWriter: httptest.NewRecorder(),
				},
				statusCode: 100,
			},
			&action{
				statusCode: 100,
			},
		),
		gen(
			"status code 999",
			[]string{CndWriteStatusCode},
			[]string{ActCheckStatus},
			&condition{
				ww: &utilhttp.WrappedWriter{
					ResponseWriter: httptest.NewRecorder(),
				},
				statusCode: 999,
			},
			&action{
				statusCode: 999,
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			tt.C().ww.WriteHeader(tt.C().statusCode)
			testutil.Diff(t, tt.A().statusCode, tt.C().ww.StatusCode())
		})
	}
}
