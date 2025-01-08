package bodylimit

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
	utilhttp "github.com/aileron-gateway/aileron-gateway/util/http"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestLimitReadCloser_Unwrap(t *testing.T) {
	type condition struct {
		rc *limitReadCloser
	}

	type action struct {
		inner io.ReadCloser
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"nil ReadCloser",
			[]string{},
			[]string{},
			&condition{
				rc: &limitReadCloser{
					ReadCloser: nil,
				},
			},
			&action{
				inner: nil,
			},
		),
		gen(
			"non-nil ReadCloser",
			[]string{},
			[]string{},
			&condition{
				rc: &limitReadCloser{
					ReadCloser: os.Stdin,
				},
			},
			&action{
				inner: os.Stdin,
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			inner := tt.C().rc.Unwrap()
			testutil.Diff(t, tt.A().inner, inner, cmp.Comparer(testutil.ComparePointer[io.ReadCloser]))
		})
	}
}

func TestLimitReadCloser_Read(t *testing.T) {
	type condition struct {
		rc *limitReadCloser
	}

	type action struct {
		exceed bool
		read   string
		rec    any
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"shorter than max size",
			[]string{},
			[]string{},
			&condition{
				rc: &limitReadCloser{
					ReadCloser: io.NopCloser(bytes.NewReader([]byte("1234567890"))),
					maxSize:    100,
				},
			},
			&action{
				exceed: false,
				read:   "1234567890",
			},
		),
		gen(
			"same as max size",
			[]string{},
			[]string{},
			&condition{
				rc: &limitReadCloser{
					ReadCloser: io.NopCloser(bytes.NewReader([]byte("1234567890"))),
					maxSize:    10,
				},
			},
			&action{
				exceed: false,
				read:   "1234567890",
			},
		),
		gen(
			"longer than max size",
			[]string{},
			[]string{},
			&condition{
				rc: &limitReadCloser{
					ReadCloser: io.NopCloser(bytes.NewReader([]byte("1234567890"))),
					maxSize:    5,
				},
			},
			&action{
				exceed: true,
				read:   "1234567890",
				rec:    http.ErrAbortHandler,
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			defer func() {
				testutil.Diff(t, tt.A().rec, recover(), cmpopts.EquateErrors())
			}()

			p := make([]byte, 100)
			exceed := false
			tt.C().rc.exceedFunc = func() {
				exceed = true
			}
			n, err := tt.C().rc.Read(p)
			testutil.Diff(t, nil, err)
			testutil.Diff(t, tt.A().exceed, exceed)
			testutil.Diff(t, tt.A().read, string(p[:n]))
		})
	}
}

type errorReader struct {
	err error
}

func (er *errorReader) Read(p []byte) (n int, err error) {
	return len(p), er.err
}

func TestBodyLimit(t *testing.T) {
	type condition struct {
		bl     *bodyLimit
		body   io.Reader
		length int64
	}

	type action struct {
		status int
		rec    any
		body   string
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"empty body",
			[]string{},
			[]string{},
			&condition{
				bl: &bodyLimit{
					eh:       utilhttp.GlobalErrorHandler(utilhttp.DefaultErrorHandlerName),
					maxSize:  10,
					memLimit: 10,
					tempPath: "./",
				},
				body:   nil,
				length: 0,
			},
			&action{
				status: http.StatusOK,
				rec:    nil,
				body:   "",
			},
		),
		gen(
			"skip checking",
			[]string{},
			[]string{},
			&condition{
				bl: &bodyLimit{
					eh:       utilhttp.GlobalErrorHandler(utilhttp.DefaultErrorHandlerName),
					maxSize:  -1,
					memLimit: 100,
					tempPath: "./",
				},
				body:   bytes.NewReader([]byte("1234567890")),
				length: 10,
			},
			&action{
				status: http.StatusOK,
				rec:    nil,
				body:   "1234567890",
			},
		),
		gen(
			"length required",
			[]string{},
			[]string{},
			&condition{
				bl: &bodyLimit{
					eh:       utilhttp.GlobalErrorHandler(utilhttp.DefaultErrorHandlerName),
					maxSize:  100,
					memLimit: 100,
					tempPath: "./",
				},
				body:   bytes.NewReader([]byte("1234567890")),
				length: -1,
			},
			&action{
				status: http.StatusLengthRequired,
				rec:    nil,
				body:   "1234567890",
			},
		),
		gen(
			"body shorter than limit/load on memory",
			[]string{},
			[]string{},
			&condition{
				bl: &bodyLimit{
					eh:       utilhttp.GlobalErrorHandler(utilhttp.DefaultErrorHandlerName),
					maxSize:  10,
					memLimit: 10,
					tempPath: "./",
				},
				body:   bytes.NewReader([]byte("12345")),
				length: 5,
			},
			&action{
				status: http.StatusOK,
				rec:    nil,
				body:   "12345",
			},
		),
		gen(
			"body same as limit/load on memory",
			[]string{},
			[]string{},
			&condition{
				bl: &bodyLimit{
					eh:       utilhttp.GlobalErrorHandler(utilhttp.DefaultErrorHandlerName),
					maxSize:  10,
					memLimit: 10,
					tempPath: "./",
				},
				body:   bytes.NewReader([]byte("1234567890")),
				length: 10,
			},
			&action{
				status: http.StatusOK,
				rec:    nil,
				body:   "1234567890",
			},
		),
		gen(
			"body longer than limit/load on memory",
			[]string{},
			[]string{},
			&condition{
				bl: &bodyLimit{
					eh:       utilhttp.GlobalErrorHandler(utilhttp.DefaultErrorHandlerName),
					maxSize:  10,
					memLimit: 10,
					tempPath: "./",
				},
				body:   bytes.NewReader([]byte("12345678901")),
				length: 11,
			},
			&action{
				status: http.StatusRequestEntityTooLarge,
				rec:    nil,
				body:   "12345678901",
			},
		),
		gen(
			"content length mismatch/load on memory",
			[]string{},
			[]string{},
			&condition{
				bl: &bodyLimit{
					eh:       utilhttp.GlobalErrorHandler(utilhttp.DefaultErrorHandlerName),
					maxSize:  10,
					memLimit: 10,
					tempPath: "./",
				},
				body:   bytes.NewReader([]byte("12345")),
				length: 6,
			},
			&action{
				status: http.StatusBadRequest,
				rec:    nil,
				body:   "12345",
			},
		),
		gen(
			"body shorter than limit/load on file",
			[]string{},
			[]string{},
			&condition{
				bl: &bodyLimit{
					eh:       utilhttp.GlobalErrorHandler(utilhttp.DefaultErrorHandlerName),
					maxSize:  100,
					memLimit: 5,
					tempPath: "./",
				},
				body:   bytes.NewReader([]byte("1234567890")),
				length: 10,
			},
			&action{
				status: http.StatusOK,
				rec:    nil,
				body:   "1234567890",
			},
		),
		gen(
			"body same as limit/load on file",
			[]string{},
			[]string{},
			&condition{
				bl: &bodyLimit{
					eh:       utilhttp.GlobalErrorHandler(utilhttp.DefaultErrorHandlerName),
					maxSize:  10,
					memLimit: 5,
					tempPath: "./",
				},
				body:   bytes.NewReader([]byte("1234567890")),
				length: 10,
			},
			&action{
				status: http.StatusOK,
				rec:    nil,
				body:   "1234567890",
			},
		),
		gen(
			"body longer than limit/load on file",
			[]string{},
			[]string{},
			&condition{
				bl: &bodyLimit{
					eh:       utilhttp.GlobalErrorHandler(utilhttp.DefaultErrorHandlerName),
					maxSize:  10,
					memLimit: 5,
					tempPath: "./",
				},
				body:   bytes.NewReader([]byte("12345678901")),
				length: 11,
			},
			&action{
				status: http.StatusRequestEntityTooLarge,
				rec:    nil,
				body:   "12345678901",
			},
		),
		gen(
			"content length mismatch/load on file",
			[]string{},
			[]string{},
			&condition{
				bl: &bodyLimit{
					eh:       utilhttp.GlobalErrorHandler(utilhttp.DefaultErrorHandlerName),
					maxSize:  10,
					memLimit: 5,
					tempPath: "./",
				},
				body:   bytes.NewReader([]byte("12345")),
				length: 6,
			},
			&action{
				status: http.StatusBadRequest,
				rec:    nil,
				body:   "12345",
			},
		),
		gen(
			"body shorter than limit/check on read",
			[]string{},
			[]string{},
			&condition{
				bl: &bodyLimit{
					eh:       utilhttp.GlobalErrorHandler(utilhttp.DefaultErrorHandlerName),
					maxSize:  100,
					memLimit: -1,
					tempPath: "./",
				},
				body:   bytes.NewReader([]byte("1234567890")),
				length: 10,
			},
			&action{
				status: http.StatusOK,
				rec:    nil,
				body:   "1234567890",
			},
		),
		gen(
			"body same as limit/check on read",
			[]string{},
			[]string{},
			&condition{
				bl: &bodyLimit{
					eh:       utilhttp.GlobalErrorHandler(utilhttp.DefaultErrorHandlerName),
					maxSize:  10,
					memLimit: -1,
					tempPath: "./",
				},
				body:   bytes.NewReader([]byte("1234567890")),
				length: 10,
			},
			&action{
				status: http.StatusOK,
				rec:    nil,
				body:   "1234567890",
			},
		),
		gen(
			"body longer than limit/check on read",
			[]string{},
			[]string{},
			&condition{
				bl: &bodyLimit{
					eh:       utilhttp.GlobalErrorHandler(utilhttp.DefaultErrorHandlerName),
					maxSize:  5,
					memLimit: -1,
					tempPath: "./",
				},
				body:   bytes.NewReader([]byte("1234567890")),
				length: 0,
			},
			&action{
				status: http.StatusRequestEntityTooLarge,
				rec:    http.ErrAbortHandler,
				body:   "1234567890",
			},
		),
		gen(
			"temp file create error",
			[]string{},
			[]string{},
			&condition{
				bl: &bodyLimit{
					eh:       utilhttp.GlobalErrorHandler(utilhttp.DefaultErrorHandlerName),
					maxSize:  100,
					memLimit: 5,
					tempPath: "./foo/bar/baz/",
				},
				body:   bytes.NewReader([]byte("1234567890")),
				length: 10,
			},
			&action{
				status: http.StatusInternalServerError,
				rec:    nil,
				body:   "1234567890",
			},
		),
		gen(
			"body shorter than limit/body read error",
			[]string{},
			[]string{},
			&condition{
				bl: &bodyLimit{
					eh:       utilhttp.GlobalErrorHandler(utilhttp.DefaultErrorHandlerName),
					maxSize:  100,
					memLimit: 5,
					tempPath: "./",
				},
				body:   &errorReader{err: io.ErrUnexpectedEOF},
				length: 10,
			},
			&action{
				status: http.StatusInternalServerError,
				rec:    nil,
				body:   "1234567890",
			},
		),
		gen(
			"body longer than limit/body read error",
			[]string{},
			[]string{},
			&condition{
				bl: &bodyLimit{
					eh:       utilhttp.GlobalErrorHandler(utilhttp.DefaultErrorHandlerName),
					maxSize:  10,
					memLimit: 10,
					tempPath: "./",
				},
				body:   &errorReader{err: io.ErrUnexpectedEOF},
				length: 5,
			},
			&action{
				status: http.StatusRequestEntityTooLarge,
				rec:    nil,
				body:   "1234567890",
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/test", tt.C().body)
			req.ContentLength = tt.C().length
			rec := httptest.NewRecorder()
			next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				defer func() {
					testutil.Diff(t, tt.A().rec, recover(), cmpopts.EquateErrors())
				}()
				b, err := io.ReadAll(r.Body)
				testutil.Diff(t, nil, err)
				testutil.Diff(t, tt.A().body, string(b))
				w.WriteHeader(http.StatusOK)
			})
			tt.C().bl.Middleware(next).ServeHTTP(rec, req)
			testutil.Diff(t, tt.A().status, rec.Result().StatusCode)
		})
	}
}
