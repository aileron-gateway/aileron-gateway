// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package http_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aileron-gateway/aileron-gateway/internal/testutil"
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

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"new with ResponseWriter without writing status code", &condition{
				w: httptest.NewRecorder(),
			},
			&action{
				written:    false,
				statusCode: 0,
			},
		),
		gen(
			"nil ResponseWriter", &condition{
				w: nil,
			},
			&action{
				written:    false,
				statusCode: 0,
			},
		),
		gen(
			"new with ResponseWriter with writing status code", &condition{
				w:          httptest.NewRecorder(),
				statusCode: 999,
			},
			&action{
				written:    true,
				statusCode: 999,
			},
		),
		gen(
			"new with ResponseWriter with writing status code", &condition{
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

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			ww := utilhttp.WrapWriter(tt.C.w)
			if tt.C.statusCode > 0 {
				ww.WriteHeader(tt.C.statusCode)
			}

			testutil.Diff(t, tt.A.written, ww.Written())
			testutil.Diff(t, tt.A.statusCode, ww.StatusCode())

			wwPtr := ww.(*utilhttp.WrappedWriter)
			target := tt.C.w
			if ww, ok := tt.C.w.(*utilhttp.WrappedWriter); ok {
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

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"unwrap nil", &condition{
				ww: &utilhttp.WrappedWriter{
					ResponseWriter: nil,
				},
			},
			&action{
				w: nil,
			},
		),
		gen(
			"unwrap non-nil", &condition{
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

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			w := tt.C.ww.Unwrap()
			testutil.Diff(t, tt.A.w, w, cmp.AllowUnexported(testResponseWriter{}))
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

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"status code 100", &condition{
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
			"status code 999", &condition{
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

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			w := httptest.NewRecorder()
			ww := utilhttp.WrapWriter(w)
			ww.WriteHeader(tt.C.statusCode)

			testutil.Diff(t, true, ww.Written())
			testutil.Diff(t, tt.A.statusCode, ww.StatusCode())
			testutil.Diff(t, tt.A.statusCode, w.Result().StatusCode)
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

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"status code 100", &condition{
				statusCode: 100,
				body:       "test",
			},
			&action{
				statusCode: 100,
				body:       "test",
			},
		),
		gen(
			"status code 999", &condition{
				statusCode: 999,
				body:       "test",
			},
			&action{
				statusCode: 999,
				body:       "test",
			},
		),
		gen(
			"status code 0 (don't write the code)", &condition{
				statusCode: 0,
				body:       "test",
			},
			&action{
				statusCode: 200,
				body:       "test",
			},
		),
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			w := httptest.NewRecorder()
			ww := utilhttp.WrapWriter(w)
			if tt.C.statusCode > 0 {
				ww.WriteHeader(tt.C.statusCode)
			}
			ww.Write([]byte(tt.C.body))

			testutil.Diff(t, true, ww.Written())
			testutil.Diff(t, tt.A.statusCode, ww.StatusCode())
			testutil.Diff(t, w.Result().StatusCode, ww.StatusCode())

			body, _ := io.ReadAll(w.Body)
			testutil.Diff(t, tt.A.body, string(body))
			testutil.Diff(t, len(tt.A.body), int(ww.ContentLength()))
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

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"don't write status code", &condition{
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
			"write status code", &condition{
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

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			if tt.C.write {
				tt.C.ww.WriteHeader(999)
			}

			testutil.Diff(t, tt.A.written, tt.C.ww.Written())
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

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"status code 100", &condition{
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
			"status code 999", &condition{
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

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			tt.C.ww.WriteHeader(tt.C.statusCode)
			testutil.Diff(t, tt.A.statusCode, tt.C.ww.StatusCode())
		})
	}
}
