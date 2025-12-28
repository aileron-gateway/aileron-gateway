// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package static

import (
	"bytes"
	"io"
	"io/fs"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/aileron-gateway/aileron-gateway/internal/testutil"
	utilhttp "github.com/aileron-gateway/aileron-gateway/util/http"
	"github.com/google/go-cmp/cmp/cmpopts"
)

type testResponseWriter struct {
	status int
	header http.Header
	b      bytes.Buffer
}

func (w *testResponseWriter) Header() http.Header {
	return w.header
}

func (w *testResponseWriter) Write(b []byte) (int, error) {
	return w.b.Write(b)
}

func (w *testResponseWriter) WriteHeader(status int) {
	w.status = status
}

type errorFS struct {
	http.FileSystem
}

func (f *errorFS) Open(_ string) (http.File, error) {
	return &statErrorFile{}, nil
}

type statErrorFile struct {
	http.File
}

func (f *statErrorFile) Stat() (fs.FileInfo, error) {
	return nil, fs.ErrInvalid
}

type testStatusHandler int

func (h testStatusHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(int(h))
	w.Write([]byte("test"))
}

func TestHandler_ServeHTTP(t *testing.T) {
	type condition struct {
		h    *handler
		path string
	}

	type action struct {
		path   string
		status int
		header map[string]string
		body   *regexp.Regexp
	}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"200 OK",
			&condition{
				h: &handler{
					Handler: testStatusHandler(http.StatusOK),
					eh:      utilhttp.GlobalErrorHandler(utilhttp.DefaultErrorHandlerName),
				},
			},
			&action{
				path:   "/",
				status: http.StatusOK,
				body:   regexp.MustCompile(`test`),
			},
		),
		gen(
			"200 OK with header",
			&condition{
				h: &handler{
					Handler: testStatusHandler(http.StatusOK),
					eh:      utilhttp.GlobalErrorHandler(utilhttp.DefaultErrorHandlerName),
					header:  map[string]string{"Cache-Control": "max-age=604800"},
				},
			},
			&action{
				path:   "/",
				status: http.StatusOK,
				body:   regexp.MustCompile(`test`),
				header: map[string]string{"Cache-Control": "max-age=604800"},
			},
		),
		gen(
			"200 OK with path",
			&condition{
				h: &handler{
					Handler: testStatusHandler(http.StatusOK),
					eh:      utilhttp.GlobalErrorHandler(utilhttp.DefaultErrorHandlerName),
				},
				path: "/test",
			},
			&action{
				path:   "/test",
				status: http.StatusOK,
				body:   regexp.MustCompile(`test`),
			},
		),
		gen(
			"399",
			&condition{
				h: &handler{
					Handler: testStatusHandler(300),
					eh:      utilhttp.GlobalErrorHandler(utilhttp.DefaultErrorHandlerName),
				},
			},
			&action{
				path:   "/",
				status: 300,
				body:   regexp.MustCompile(`test`),
			},
		),
		gen(
			"400 Bad Request",
			&condition{
				h: &handler{
					Handler: testStatusHandler(http.StatusBadRequest),
					eh:      utilhttp.GlobalErrorHandler(utilhttp.DefaultErrorHandlerName),
				},
			},
			&action{
				path:   "/",
				status: http.StatusBadRequest,
				body:   regexp.MustCompile(`Bad Request`),
			},
		),
		gen(
			"499",
			&condition{
				h: &handler{
					Handler: testStatusHandler(499),
					eh:      utilhttp.GlobalErrorHandler(utilhttp.DefaultErrorHandlerName),
				},
			},
			&action{
				path:   "/",
				status: 499,
				body:   regexp.MustCompile(`499`),
			},
		),
		gen(
			"500 Internal Server Error",
			&condition{
				h: &handler{
					Handler: testStatusHandler(http.StatusInternalServerError),
					eh:      utilhttp.GlobalErrorHandler(utilhttp.DefaultErrorHandlerName),
				},
			},
			&action{
				path:   "/",
				status: http.StatusInternalServerError,
				body:   regexp.MustCompile(`Internal Server Error`),
			},
		),
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r, _ := http.NewRequest(http.MethodGet, "http://test.com"+tt.C.path, nil)

			tt.C.h.ServeHTTP(w, r)

			res := w.Result()
			defer res.Body.Close()
			b, _ := io.ReadAll(res.Body)
			t.Log(string(b))

			testutil.Diff(t, tt.A.path, r.URL.Path)
			testutil.Diff(t, tt.A.status, res.StatusCode)
			testutil.Diff(t, true, tt.A.body.Match(b))
			for k, v := range tt.A.header {
				testutil.Diff(t, v, res.Header.Get(k))
			}
		})
	}
}

func TestFileOnlyDir_Open(t *testing.T) {
	type condition struct {
		fs   http.FileSystem
		name string
	}

	type action struct {
		content string
		err     any // error or errorutil.Kind
	}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"file exists",
			&condition{
				fs:   http.Dir(testDir + "ut/core/static/"),
				name: "/testdir/test.txt",
			},
			&action{
				content: "test",
			},
		),
		gen(
			"file not exists",
			&condition{
				fs:   http.Dir(testDir + "ut/core/static/"),
				name: "/testdir/not-exist.txt",
			},
			&action{
				content: "",
				err:     fs.ErrNotExist,
			},
		),
		gen(
			"stat error",
			&condition{
				fs:   &errorFS{},
				name: "/testdir/test.txt",
			},
			&action{
				content: "",
				err:     fs.ErrInvalid,
			},
		),
		gen(
			"directory",
			&condition{
				fs:   http.Dir(testDir + "ut/core/static/"),
				name: "/testdir",
			},
			&action{
				content: "",
				err:     fs.ErrNotExist,
			},
		),
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			fod := &fileOnlyDir{
				fs: tt.C.fs,
			}

			f, err := fod.Open(tt.C.name)
			testutil.DiffError(t, tt.A.err, nil, err, cmpopts.EquateErrors())

			if err != nil {
				testutil.Diff(t, nil, f)
			} else {
				var b bytes.Buffer
				b.ReadFrom(f)
				testutil.Diff(t, tt.A.content, b.String())
			}
		})
	}
}

func TestDiscardWriter_Write(t *testing.T) {
	type condition struct {
		dw    *discardWriter
		write []byte
	}

	type action struct {
		expect []byte
	}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"status 0",
			&condition{
				dw: &discardWriter{
					ResponseWriter: &testResponseWriter{},
					status:         0,
				},
				write: []byte("test"),
			},
			&action{
				expect: []byte("test"),
			},
		),
		gen(
			"status 399",
			&condition{
				dw: &discardWriter{
					ResponseWriter: &testResponseWriter{},
					status:         399,
				},
				write: []byte("test"),
			},
			&action{
				expect: []byte("test"),
			},
		),
		gen(
			"status 400",
			&condition{
				dw: &discardWriter{
					ResponseWriter: &testResponseWriter{},
					status:         400,
				},
				write: []byte("test"),
			},
			&action{
				expect: nil,
			},
		),
		gen(
			"status 500",
			&condition{
				dw: &discardWriter{
					ResponseWriter: &testResponseWriter{},
					status:         500,
				},
				write: []byte("test"),
			},
			&action{
				expect: nil,
			},
		),
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			tt.C.dw.Write(tt.C.write)

			w := tt.C.dw.ResponseWriter.(*testResponseWriter)
			testutil.Diff(t, tt.A.expect, w.b.Bytes())
		})
	}
}

func TestDiscardWriter_WriteHeader(t *testing.T) {
	type condition struct {
		dw     *discardWriter
		status int
	}

	type action struct {
		status int
	}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"status 0",
			&condition{
				dw: &discardWriter{
					ResponseWriter: &testResponseWriter{},
					status:         0,
				},
				status: 0,
			},
			&action{
				status: 0,
			},
		),
		gen(
			"status 399",
			&condition{
				dw: &discardWriter{
					ResponseWriter: &testResponseWriter{},
					status:         399,
				},
				status: 399,
			},
			&action{
				status: 399,
			},
		),
		gen(
			"status 400",
			&condition{
				dw: &discardWriter{
					ResponseWriter: &testResponseWriter{},
					status:         400,
				},
				status: 400,
			},
			&action{
				status: 0,
			},
		),
		gen(
			"status 500",
			&condition{
				dw: &discardWriter{
					ResponseWriter: &testResponseWriter{},
					status:         500,
				},
				status: 500,
			},
			&action{
				status: 0,
			},
		),
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			tt.C.dw.WriteHeader(tt.C.status)

			w := tt.C.dw.ResponseWriter.(*testResponseWriter)
			testutil.Diff(t, tt.A.status, w.status)
		})
	}
}
