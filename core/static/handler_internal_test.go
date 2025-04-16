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

	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
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

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	cndStatus200 := tb.Condition("200", "write status 200")
	cndStatus399 := tb.Condition("399", "write status 399")
	cndStatus400 := tb.Condition("400", "write status 400")
	cndStatus499 := tb.Condition("499", "write status 499")
	cndStatus500 := tb.Condition("500", "write status 500")
	actCheckStatus := tb.Action("check status", "check the responded status code")
	actCheckBody := tb.Action("check body", "check the written body")
	actCheckError := tb.Action("error", "check that there is an error")
	actCheckNoError := tb.Action("no error", "check that there is no error")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"200 OK",
			[]string{cndStatus200},
			[]string{actCheckStatus, actCheckBody, actCheckNoError},
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
			[]string{cndStatus200},
			[]string{actCheckStatus, actCheckBody, actCheckNoError},
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
			[]string{cndStatus200},
			[]string{actCheckStatus, actCheckBody, actCheckNoError},
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
			[]string{cndStatus399},
			[]string{actCheckStatus, actCheckBody, actCheckNoError},
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
			[]string{cndStatus400},
			[]string{actCheckStatus, actCheckBody, actCheckError},
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
			[]string{cndStatus499},
			[]string{actCheckStatus, actCheckBody, actCheckError},
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
			[]string{cndStatus500},
			[]string{actCheckStatus, actCheckBody, actCheckError},
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

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			w := httptest.NewRecorder()
			r, _ := http.NewRequest(http.MethodGet, "http://test.com"+tt.C().path, nil)

			tt.C().h.ServeHTTP(w, r)

			res := w.Result()
			defer res.Body.Close()
			b, _ := io.ReadAll(res.Body)
			t.Log(string(b))

			testutil.Diff(t, tt.A().path, r.URL.Path)
			testutil.Diff(t, tt.A().status, res.StatusCode)
			testutil.Diff(t, true, tt.A().body.Match(b))
			for k, v := range tt.A().header {
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

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	cndFileExist := tb.Condition("file exists", "access to a file exists")
	cndStatError := tb.Condition("file not exists", "getting file info returns an error ")
	cndDirectory := tb.Condition("directory", "access to a directory exists")
	actCheckError := tb.Action("discard", "check that the written stats are discarded")
	actCheckNoError := tb.Action("written", "check that the written status are written to the underlying writer")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"file exists",
			[]string{cndFileExist},
			[]string{actCheckNoError},
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
			[]string{},
			[]string{actCheckError},
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
			[]string{cndStatError},
			[]string{actCheckError},
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
			[]string{cndDirectory},
			[]string{actCheckError},
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

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			fod := &fileOnlyDir{
				fs: tt.C().fs,
			}

			f, err := fod.Open(tt.C().name)
			testutil.DiffError(t, tt.A().err, nil, err, cmpopts.EquateErrors())

			if err != nil {
				testutil.Diff(t, nil, f)
			} else {
				var b bytes.Buffer
				b.ReadFrom(f)
				testutil.Diff(t, tt.A().content, b.String())
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

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	cndStatus0 := tb.Condition("0", "write to the 0 status writer")
	cndStatus399 := tb.Condition("399", "write to the 399 status writer")
	cndStatus400 := tb.Condition("400", "write to the 400 status writer")
	cndStatus500 := tb.Condition("500", "write to the 500 status writer")
	actCheckDiscard := tb.Action("discard", "check that the written bytes are discarded")
	actCheckWritten := tb.Action("written", "check that the written bytes are written to the underlying writer")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"status 0",
			[]string{cndStatus0},
			[]string{actCheckWritten},
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
			[]string{cndStatus399},
			[]string{actCheckWritten},
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
			[]string{cndStatus400},
			[]string{actCheckDiscard},
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
			[]string{cndStatus500},
			[]string{actCheckDiscard},
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

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			tt.C().dw.Write(tt.C().write)

			w := tt.C().dw.ResponseWriter.(*testResponseWriter)
			testutil.Diff(t, tt.A().expect, w.b.Bytes())
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

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	cndStatus0 := tb.Condition("0", "write status 0")
	cndStatus399 := tb.Condition("399", "write status 399")
	cndStatus400 := tb.Condition("400", "write status 400")
	cndStatus500 := tb.Condition("500", "write status 500")
	actCheckDiscard := tb.Action("discard", "check that the written stats are discarded")
	actCheckWritten := tb.Action("written", "check that the written status are written to the underlying writer")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"status 0",
			[]string{cndStatus0},
			[]string{actCheckWritten},
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
			[]string{cndStatus399},
			[]string{actCheckWritten},
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
			[]string{cndStatus400},
			[]string{actCheckDiscard},
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
			[]string{cndStatus500},
			[]string{actCheckDiscard},
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

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			tt.C().dw.WriteHeader(tt.C().status)

			w := tt.C().dw.ResponseWriter.(*testResponseWriter)
			testutil.Diff(t, tt.A().status, w.status)
		})
	}
}
