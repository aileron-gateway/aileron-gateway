// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package httplogger

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/aileron-gateway/aileron-gateway/core"
	"github.com/aileron-gateway/aileron-gateway/kernel/log"
	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
	utilhttp "github.com/aileron-gateway/aileron-gateway/util/http"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

type testResponseWriter struct {
	http.ResponseWriter
	id      string
	h       http.Header
	written []byte
	status  int
	err     error
}

func (w *testResponseWriter) Unwrap() http.ResponseWriter {
	return w.ResponseWriter
}

func (w *testResponseWriter) Header() http.Header {
	return w.h
}

func (w *testResponseWriter) Write(b []byte) (int, error) {
	w.written = append(w.written, b...)
	return len(b), w.err
}

func (w *testResponseWriter) WriteHeader(statusCode int) {
	w.status = statusCode
}

func TestLoggingWriter_Unwrap(t *testing.T) {
	type condition struct {
		w *wrappedWriter
	}

	type action struct {
		w http.ResponseWriter
	}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"nil",
			&condition{
				w: &wrappedWriter{
					ResponseWriter: nil,
				},
			},
			&action{
				w: nil,
			},
		),
		gen(
			"non-nil",
			&condition{
				w: &wrappedWriter{
					ResponseWriter: &testResponseWriter{
						id: "test",
					},
				},
			},
			&action{
				w: &testResponseWriter{
					id: "test",
				},
			},
		),
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			w := tt.C.w.Unwrap()
			testutil.Diff(t, tt.A.w, w, cmp.AllowUnexported(testResponseWriter{}))
		})
	}
}

type testFlushResponseWriter struct {
	http.ResponseWriter
	id      string
	flushed bool
}

func (w *testFlushResponseWriter) Flush() {
	w.flushed = true
}

func TestLoggingWriter_Flush(t *testing.T) {
	type condition struct {
		w http.ResponseWriter
	}

	type action struct {
		w http.ResponseWriter
	}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"non flusher",
			&condition{
				w: &testResponseWriter{id: "test"},
			},
			&action{
				w: &testResponseWriter{id: "test"},
			},
		),
		gen(
			"flusher",
			&condition{
				w: &testFlushResponseWriter{id: "test"},
			},
			&action{
				w: &testFlushResponseWriter{id: "test", flushed: true},
			},
		),
		gen(
			"non inner flusher",
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

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			w := &wrappedWriter{
				ResponseWriter: tt.C.w,
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
			testutil.Diff(t, tt.A.w, tt.C.w, opts...)

			w.Flush()
			testutil.Diff(t, true, w.flushChecked)
			testutil.Diff(t, true, w.flushFunc != nil)
			testutil.Diff(t, tt.A.w, tt.C.w, opts...)
		})
	}
}

func TestWrappedWriter_WriteHeader(t *testing.T) {
	type condition struct {
		w      *wrappedWriter
		status int
	}

	type action struct {
	}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"write 0",
			&condition{
				w: &wrappedWriter{
					ResponseWriter: &testResponseWriter{},
				},
				status: http.StatusOK,
			},
			&action{},
		),
		gen(
			"write 200",
			&condition{
				w: &wrappedWriter{
					ResponseWriter: &testResponseWriter{},
				},
				status: http.StatusOK,
			},
			&action{},
		),
		gen(
			"write 500",
			&condition{
				w: &wrappedWriter{
					ResponseWriter: &testResponseWriter{},
				},
				status: http.StatusInternalServerError,
			},
			&action{},
		),
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			called := 0
			tt.C.w.writtenFunc = func() {
				called += 1
			}

			tt.C.w.WriteHeader(tt.C.status) // Call once
			tt.C.w.WriteHeader(tt.C.status) // Call twice

			testutil.Diff(t, 1, called)
			testutil.Diff(t, true, tt.C.w.written)
			testutil.Diff(t, tt.C.status, tt.C.w.status)
		})
	}
}

func TestWrappedWriter_Write(t *testing.T) {
	type condition struct {
		w     *wrappedWriter
		write [][]byte
	}

	type action struct {
		written string
	}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"dump",
			&condition{
				w: &wrappedWriter{
					ResponseWriter: &testResponseWriter{},
					dump:           true,
				},
				write: [][]byte{[]byte("foo"), []byte("bar")},
			},
			&action{
				written: "foobar",
			},
		),
		gen(
			"not dump",
			&condition{
				w: &wrappedWriter{
					ResponseWriter: &testResponseWriter{},
					dump:           false,
				},
				write: [][]byte{[]byte("foo"), []byte("bar")},
			},
			&action{
				written: "",
			},
		),
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			var buf bytes.Buffer
			tt.C.w.w = &buf

			called := 0
			tt.C.w.writtenFunc = func() {
				called += 1
			}

			for _, v := range tt.C.write {
				tt.C.w.Write(v)
			}

			testutil.Diff(t, 1, called)
			testutil.Diff(t, true, tt.C.w.written)
			testutil.Diff(t, tt.A.written, buf.String())
		})
	}
}

type testLogger struct {
	log.Logger
	written int
	b       []byte
	msg     string
	kvs     []any
}

func (l *testLogger) Write(p []byte) (n int, err error) {
	l.b = append(l.b, p...)
	l.written += len(p)
	return l.Logger.(io.Writer).Write(p)
}

func (l *testLogger) Info(ctx context.Context, msg string, keyValues ...any) {
	l.msg = msg
	l.kvs = keyValues
}

func TestHTTPLogger_Middleware(t *testing.T) {
	type condition struct {
		logger    *httpLogger
		req       *http.Request
		resHeader http.Header
		resStatus int
		resBody   [][]byte
	}

	type action struct {
		status    int
		body      string
		reqFmtLog string
		resFmtLog string
		reqKVs    []any
		resKVs    []any
	}

	testGetReq, _ := http.NewRequest(http.MethodGet, "http://test.com/get?foo=bar&alice=bob", nil)
	testGetReq = testGetReq.WithContext(context.WithValue(context.Background(), idContextKey, "test-id"))
	testPostReq, _ := http.NewRequest(http.MethodPost, "http://test.com/post?foo=bar&alice=bob", bytes.NewReader([]byte("testRequestBody")))
	testPostReq.Header.Set("Content-Type", "text/plain")

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"record GET",
			&condition{
				logger: &httpLogger{
					req:  &baseLogger{allHeaders: true},
					res:  &baseLogger{allHeaders: true},
					zone: time.Local,
				},
				req:       testGetReq,
				resHeader: http.Header{"test": {"value"}, "Content-Length": {"16"}},
				resStatus: http.StatusOK,
				resBody:   [][]byte{[]byte("testResponseBody")},
			},
			&action{
				status: http.StatusOK,
				body:   "testResponseBody",
				reqKVs: []any{
					"request",
					map[string]any{
						keyID:     "",
						keyTime:   "",
						keyMethod: "GET",
						keyPath:   "/get",
						keyQuery:  "foo=bar&alice=bob",
						keyHost:   "test.com",
						keyRemote: "",
						keyProto:  "HTTP/1.1",
						keyHeader: map[string]string{},
						keySize:   int64(0),
					},
				},
				resKVs: []any{
					"response",
					map[string]any{
						keyID:       "",
						keyTime:     "",
						keyDuration: int64(0),
						keyStatus:   200,
						keyHeader:   map[string]string{"Test": "value", "Content-Length": "16"},
						keySize:     int64(16), // = len("testResponseBody")
					},
				},
			},
		),
		gen(
			"record POST",
			&condition{
				logger: &httpLogger{
					req:  &baseLogger{allHeaders: true},
					res:  &baseLogger{allHeaders: true},
					zone: time.Local,
				},
				req:       testPostReq,
				resHeader: http.Header{"test": {"value"}, "Content-Length": {"16"}},
				resStatus: http.StatusOK,
				resBody:   [][]byte{[]byte("testResponseBody")},
			},
			&action{
				status: http.StatusOK,
				body:   "testResponseBody",
				reqKVs: []any{
					"request",
					map[string]any{
						keyID:     "",
						keyTime:   "",
						keyMethod: "POST",
						keyPath:   "/post",
						keyQuery:  "foo=bar&alice=bob",
						keyHost:   "test.com",
						keyRemote: "",
						keyProto:  "HTTP/1.1",
						keyHeader: map[string]string{"Content-Type": "text/plain"},
						keySize:   int64(15),
					},
				},
				resKVs: []any{
					"response",
					map[string]any{
						keyID:       "",
						keyTime:     "",
						keyDuration: int64(0),
						keyStatus:   200,
						keyHeader:   map[string]string{"Test": "value", "Content-Length": "16"},
						keySize:     int64(16), // = len("testResponseBody")
					},
				},
			},
		),
		gen(
			"Status was not written",
			&condition{
				logger: &httpLogger{
					req:  &baseLogger{allHeaders: true},
					res:  &baseLogger{allHeaders: true},
					zone: time.Local,
				},
				req:       testGetReq,
				resHeader: http.Header{},
				resStatus: 0,
				resBody:   [][]byte{[]byte("testResponseBody")},
			},
			&action{
				status: http.StatusOK,
				body:   "testResponseBody",
				reqKVs: []any{
					"request",
					map[string]any{
						keyID:     "",
						keyTime:   "",
						keyMethod: "GET",
						keyPath:   "/get",
						keyQuery:  "foo=bar&alice=bob",
						keyHost:   "test.com",
						keyRemote: "",
						keyProto:  "HTTP/1.1",
						keyHeader: map[string]string{},
						keySize:   int64(0),
					},
				},
				resKVs: []any{
					"response",
					map[string]any{
						keyID:       "",
						keyTime:     "",
						keyDuration: int64(0),
						keyStatus:   0,
						keyHeader:   map[string]string{},
						keySize:     int64(-1), // = No Content-Length header.
					},
				},
			},
		),
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				for k, vv := range tt.C.resHeader {
					for _, v := range vv {
						w.Header().Add(k, v)
					}
				}
				if tt.C.resStatus > 0 {
					w.WriteHeader(tt.C.resStatus)
				}
				for _, b := range tt.C.resBody {
					w.Write(b)
				}
			})

			reqLg := &testLogger{Logger: log.GlobalLogger(log.DefaultLoggerName)}
			resLg := &testLogger{Logger: log.GlobalLogger(log.DefaultLoggerName)}
			tt.C.logger.req.lg, tt.C.logger.req.w = reqLg, reqLg
			tt.C.logger.res.lg, tt.C.logger.res.w = resLg, resLg

			handler := tt.C.logger.Middleware(h)
			w := httptest.NewRecorder()
			handler.ServeHTTP(w, tt.C.req)
			w.Result()
			body, _ := io.ReadAll(w.Result().Body)
			testutil.Diff(t, tt.A.status, w.Result().StatusCode)
			testutil.Diff(t, tt.A.body, string(body))

			opts := []cmp.Option{
				cmpopts.EquateEmpty(),
				cmpopts.IgnoreMapEntries(func(k string, _ any) bool {
					return k == keyDuration || k == keyID || k == keyTime
				}),
			}
			testutil.Diff(t, tt.A.reqFmtLog, string(reqLg.b))
			testutil.Diff(t, tt.A.resFmtLog, string(resLg.b))
			testutil.Diff(t, tt.A.reqKVs, reqLg.kvs, opts...)
			testutil.Diff(t, tt.A.resKVs, resLg.kvs, opts...)
		})
	}
}

func TestJournalLogger_Middleware(t *testing.T) {
	type condition struct {
		logger    *journalLogger
		req       *http.Request
		resHeader http.Header
		resStatus int
		resBody   [][]byte
	}

	type action struct {
		status    int
		body      string
		reqFmtLog string
		resFmtLog string
		reqKVs    []any
		resKVs    []any
	}

	testGetReq, _ := http.NewRequest(http.MethodGet, "http://test.com/get?foo=bar&alice=bob", nil)
	testGetReq = testGetReq.WithContext(context.WithValue(context.Background(), idContextKey, "test-id"))
	testPostReq := func() *http.Request {
		r, _ := http.NewRequest(http.MethodPost, "http://test.com/post?foo=bar&alice=bob", bytes.NewReader([]byte("testRequestBody")))
		r.Header.Set("Content-Type", "text/plain")
		return r
	}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"record GET",
			&condition{
				logger: &journalLogger{
					req:  &baseLogger{allHeaders: true, maxBody: 100, mimes: []string{"text/plain"}},
					res:  &baseLogger{allHeaders: true, maxBody: 100, mimes: []string{"text/plain"}},
					zone: time.Local,
				},
				req: testGetReq,
				resHeader: http.Header{
					"test":           {"value"},
					"Content-Length": {"16"},
					"Content-Type":   {"text/plain"},
				},
				resStatus: http.StatusOK,
				resBody:   [][]byte{[]byte("testResponseBody")},
			},
			&action{
				status: http.StatusOK,
				body:   "testResponseBody",
				reqKVs: []any{
					"request",
					map[string]any{
						keyID:     "",
						keyTime:   "",
						keyMethod: "GET",
						keyPath:   "/get",
						keyQuery:  "foo=bar&alice=bob",
						keyHost:   "test.com",
						keyRemote: "",
						keyProto:  "HTTP/1.1",
						keyHeader: map[string]string{},
						keySize:   int64(0),
						keyBody:   "",
					},
				},
				resKVs: []any{
					"response",
					map[string]any{
						keyID:       "",
						keyTime:     "",
						keyDuration: int64(0),
						keyStatus:   200,
						keyHeader: map[string]string{
							"Test":           "value",
							"Content-Length": "16",
							"Content-Type":   "text/plain",
						},
						keySize: int64(16), // = len("testResponseBody")
						keyBody: "testResponseBody",
					},
				},
			},
		),
		gen(
			"record POST",
			&condition{
				logger: &journalLogger{
					req:  &baseLogger{allHeaders: true, maxBody: 100, mimes: []string{"text/plain"}},
					res:  &baseLogger{allHeaders: true, maxBody: 100, mimes: []string{"text/plain"}},
					zone: time.Local,
				},
				req: testPostReq(),
				resHeader: http.Header{
					"test":           {"value"},
					"Content-Length": {"16"},
					"Content-Type":   {"text/plain"},
				},
				resStatus: http.StatusOK,
				resBody:   [][]byte{[]byte("testResponseBody")},
			},
			&action{
				status: http.StatusOK,
				body:   "testResponseBody",
				reqKVs: []any{
					"request",
					map[string]any{
						keyID:     "",
						keyTime:   "",
						keyMethod: "POST",
						keyPath:   "/post",
						keyQuery:  "foo=bar&alice=bob",
						keyHost:   "test.com",
						keyRemote: "",
						keyProto:  "HTTP/1.1",
						keyHeader: map[string]string{"Content-Type": "text/plain"},
						keySize:   int64(15),
						keyBody:   "testRequestBody",
					},
				},
				resKVs: []any{
					"response",
					map[string]any{
						keyID:       "",
						keyTime:     "",
						keyDuration: int64(0),
						keyStatus:   200,
						keyHeader: map[string]string{
							"Test":           "value",
							"Content-Length": "16",
							"Content-Type":   "text/plain",
						},
						keySize: int64(16), // = len("testResponseBody")
						keyBody: "testResponseBody",
					},
				},
			},
		),
		gen(
			"Status was not written",
			&condition{
				logger: &journalLogger{
					req:  &baseLogger{allHeaders: true, maxBody: 100, mimes: []string{"text/plain"}},
					res:  &baseLogger{allHeaders: true, maxBody: 100, mimes: []string{"text/plain"}},
					zone: time.Local,
				},
				req:       testGetReq,
				resHeader: http.Header{},
				resStatus: 0,
				resBody:   [][]byte{[]byte("testResponseBody")},
			},
			&action{
				status: http.StatusOK,
				body:   "testResponseBody",
				reqKVs: []any{
					"request",
					map[string]any{
						keyID:     "",
						keyTime:   "",
						keyMethod: "GET",
						keyPath:   "/get",
						keyQuery:  "foo=bar&alice=bob",
						keyHost:   "test.com",
						keyRemote: "",
						keyProto:  "HTTP/1.1",
						keyHeader: map[string]string{},
						keySize:   int64(0),
						keyBody:   "",
					},
				},
				resKVs: []any{
					"response",
					map[string]any{
						keyID:       "",
						keyTime:     "",
						keyDuration: int64(0),
						keyStatus:   0,
						keyHeader: map[string]string{
							"Content-Type": "text/plain; charset=utf-8",
						},
						keySize: int64(-1), // = No Content-Length header.
						keyBody: "",
					},
				},
			},
		),
		gen(
			"reader create error",
			&condition{
				logger: &journalLogger{
					lg: log.GlobalLogger(log.DefaultLoggerName),
					eh: utilhttp.GlobalErrorHandler(utilhttp.DefaultErrorHandlerName),
					req: &baseLogger{
						lg:         log.GlobalLogger(log.DefaultLoggerName),
						allHeaders: true,
						maxBody:    1,
						mimes:      []string{"text/plain"},
						bodyPath:   "./\r\n\t\x00",
					},
					res: &baseLogger{
						allHeaders: true,
						maxBody:    100,
						mimes:      []string{"text/plain"},
					},
					zone: time.Local,
				},
				req: testPostReq(),
				resHeader: http.Header{
					"test":           {"value"},
					"Content-Length": {"16"},
					"Content-Type":   {"text/plain"},
				},
				resStatus: http.StatusOK,
				resBody:   [][]byte{[]byte("testResponseBody")},
			},
			&action{
				status: http.StatusInternalServerError,
				body:   `{"status":500,"statusText":"Internal Server Error"}`,
				reqKVs: nil,
				resKVs: nil,
			},
		),
		gen(
			"writer create error",
			&condition{
				logger: &journalLogger{
					lg: log.GlobalLogger(log.DefaultLoggerName),
					eh: utilhttp.GlobalErrorHandler(utilhttp.DefaultErrorHandlerName),
					req: &baseLogger{
						allHeaders: true,
						maxBody:    100,
						mimes:      []string{"text/plain"},
					},
					res: &baseLogger{
						allHeaders: true,
						maxBody:    1,
						mimes:      []string{"text/plain"},
						bodyPath:   "./\r\n\t\x00",
					},
					zone: time.Local,
				},
				req: testPostReq(),
				resHeader: http.Header{
					"test":           {"value"},
					"Content-Length": {"16"},
					"Content-Type":   {"text/plain"},
				},
				resStatus: http.StatusOK,
				resBody:   [][]byte{[]byte("testResponseBody")},
			},
			&action{
				status: http.StatusOK,
				body:   "testResponseBody",
				reqKVs: []any{
					"request",
					map[string]any{
						keyID:     "",
						keyTime:   "",
						keyMethod: "POST",
						keyPath:   "/post",
						keyQuery:  "foo=bar&alice=bob",
						keyHost:   "test.com",
						keyRemote: "",
						keyProto:  "HTTP/1.1",
						keyHeader: map[string]string{"Content-Type": "text/plain"},
						keySize:   int64(15),
						keyBody:   "testRequestBody",
					},
				},
				resKVs: []any{
					"response",
					map[string]any{
						keyID:       "",
						keyTime:     "",
						keyDuration: int64(0),
						keyStatus:   200,
						keyHeader: map[string]string{
							"Test":           "value",
							"Content-Length": "16",
							"Content-Type":   "text/plain",
						},
						keySize: int64(16), // = len("testResponseBody")
						keyBody: "",        // Not written by error
					},
				},
			},
		),
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				for k, vv := range tt.C.resHeader {
					for _, v := range vv {
						w.Header().Add(k, v)
					}
				}
				if tt.C.resStatus > 0 {
					w.WriteHeader(tt.C.resStatus)
				}
				for _, b := range tt.C.resBody {
					w.Write(b)
				}
			})

			reqLg := &testLogger{Logger: log.GlobalLogger(log.DefaultLoggerName)}
			resLg := &testLogger{Logger: log.GlobalLogger(log.DefaultLoggerName)}
			tt.C.logger.req.lg, tt.C.logger.req.w = reqLg, reqLg
			tt.C.logger.res.lg, tt.C.logger.res.w = resLg, resLg

			handler := tt.C.logger.Middleware(h)
			w := httptest.NewRecorder()
			handler.ServeHTTP(w, tt.C.req)
			w.Result()
			body, _ := io.ReadAll(w.Result().Body)
			testutil.Diff(t, tt.A.status, w.Result().StatusCode)
			testutil.Diff(t, tt.A.body, string(body))

			opts := []cmp.Option{
				cmpopts.EquateEmpty(),
				cmpopts.IgnoreMapEntries(func(k string, _ any) bool {
					return k == keyDuration || k == keyID || k == keyTime
				}),
			}
			testutil.Diff(t, tt.A.reqFmtLog, string(reqLg.b))
			testutil.Diff(t, tt.A.resFmtLog, string(resLg.b))
			testutil.Diff(t, tt.A.reqKVs, reqLg.kvs, opts...)
			testutil.Diff(t, tt.A.resKVs, resLg.kvs, opts...)
		})
	}
}

func TestHTTPLogger_Tripperware(t *testing.T) {
	type condition struct {
		logger *httpLogger
		req    *http.Request
		res    *http.Response
		err    error
	}

	type action struct {
		status    int
		body      string
		reqFmtLog string
		resFmtLog string
		reqKVs    []any
		resKVs    []any
	}

	testGetReq, _ := http.NewRequest(http.MethodGet, "http://test.com/get?foo=bar&alice=bob", nil)
	testGetReq = testGetReq.WithContext(context.WithValue(context.Background(), idContextKey, "test-id"))
	testPostReq, _ := http.NewRequest(http.MethodPost, "http://test.com/post?foo=bar&alice=bob", bytes.NewReader([]byte("testRequestBody")))
	testPostReq.Header.Set("Content-Type", "text/plain")

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"record GET",
			&condition{
				logger: &httpLogger{
					req:  &baseLogger{allHeaders: true},
					res:  &baseLogger{allHeaders: true},
					zone: time.Local,
				},
				req: testGetReq,
				res: &http.Response{
					StatusCode:    http.StatusOK,
					Header:        http.Header{"Test": {"value"}},
					Body:          io.NopCloser(bytes.NewReader([]byte("testResponseBody"))),
					ContentLength: int64(len("testResponseBody")),
				},
			},
			&action{
				status: http.StatusOK,
				body:   "testResponseBody",
				reqKVs: []any{
					"request",
					map[string]any{
						keyID:     "",
						keyTime:   "",
						keyMethod: "GET",
						keyPath:   "/get",
						keyQuery:  "foo=bar&alice=bob",
						keyHost:   "test.com",
						keyRemote: "",
						keyProto:  "HTTP/1.1",
						keyHeader: map[string]string{},
						keySize:   int64(0),
					},
				},
				resKVs: []any{
					"response",
					map[string]any{
						keyID:       "",
						keyTime:     "",
						keyDuration: int64(0),
						keyStatus:   200,
						keyHeader:   map[string]string{"Test": "value"},
						keySize:     int64(16), // = len("testResponseBody")
					},
				},
			},
		),
		gen(
			"record POST",
			&condition{
				logger: &httpLogger{
					req:  &baseLogger{allHeaders: true},
					res:  &baseLogger{allHeaders: true},
					zone: time.Local,
				},
				req: testPostReq,
				res: &http.Response{
					StatusCode:    http.StatusOK,
					Header:        http.Header{"Test": {"value"}},
					Body:          io.NopCloser(bytes.NewReader([]byte("testResponseBody"))),
					ContentLength: int64(len("testResponseBody")),
				},
			},
			&action{
				status: http.StatusOK,
				body:   "testResponseBody",
				reqKVs: []any{
					"request",
					map[string]any{
						keyID:     "",
						keyTime:   "",
						keyMethod: "POST",
						keyPath:   "/post",
						keyQuery:  "foo=bar&alice=bob",
						keyHost:   "test.com",
						keyRemote: "",
						keyProto:  "HTTP/1.1",
						keyHeader: map[string]string{"Content-Type": "text/plain"},
						keySize:   int64(15),
					},
				},
				resKVs: []any{
					"response",
					map[string]any{
						keyID:       "",
						keyTime:     "",
						keyDuration: int64(0),
						keyStatus:   200,
						keyHeader:   map[string]string{"Test": "value"},
						keySize:     int64(16), // = len("testResponseBody")
					},
				},
			},
		),
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			r := core.RoundTripperFunc(func(r *http.Request) (*http.Response, error) {
				return tt.C.res, tt.C.err
			})

			reqLg := &testLogger{Logger: log.GlobalLogger(log.DefaultLoggerName)}
			resLg := &testLogger{Logger: log.GlobalLogger(log.DefaultLoggerName)}
			tt.C.logger.req.lg, tt.C.logger.req.w = reqLg, reqLg
			tt.C.logger.res.lg, tt.C.logger.res.w = resLg, resLg

			roundTripper := tt.C.logger.Tripperware(r)
			w, _ := roundTripper.RoundTrip(tt.C.req)

			body, _ := io.ReadAll(w.Body)
			testutil.Diff(t, tt.A.status, w.StatusCode)
			testutil.Diff(t, tt.A.body, string(body))

			opts := []cmp.Option{
				cmpopts.EquateEmpty(),
				cmpopts.IgnoreMapEntries(func(k string, _ any) bool {
					return k == keyDuration || k == keyID || k == keyTime
				}),
			}
			testutil.Diff(t, tt.A.reqFmtLog, string(reqLg.b))
			testutil.Diff(t, tt.A.resFmtLog, string(resLg.b))
			testutil.Diff(t, tt.A.reqKVs, reqLg.kvs, opts...)
			testutil.Diff(t, tt.A.resKVs, resLg.kvs, opts...)
		})
	}
}

func TestJournalLogger_Tripperware(t *testing.T) {
	type condition struct {
		logger *journalLogger
		req    *http.Request
		res    *http.Response
		err    error
	}

	type action struct {
		status    int
		body      string
		reqFmtLog string
		resFmtLog string
		reqKVs    []any
		resKVs    []any
	}

	testGetReq, _ := http.NewRequest(http.MethodGet, "http://test.com/get?foo=bar&alice=bob", nil)
	testGetReq = testGetReq.WithContext(context.WithValue(context.Background(), idContextKey, "test-id"))
	testPostReq := func() *http.Request {
		r, _ := http.NewRequest(http.MethodPost, "http://test.com/post?foo=bar&alice=bob", bytes.NewReader([]byte("testRequestBody")))
		r.Header.Set("Content-Type", "text/plain")
		return r
	}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"record GET",
			&condition{
				logger: &journalLogger{
					req:  &baseLogger{allHeaders: true, maxBody: 100, mimes: []string{"text/plain"}},
					res:  &baseLogger{allHeaders: true, maxBody: 100, mimes: []string{"text/plain"}},
					zone: time.Local,
				},
				req: testGetReq,
				res: &http.Response{
					StatusCode:    http.StatusOK,
					Header:        http.Header{"Test": {"value"}, "Content-Type": {"text/plain"}},
					Body:          io.NopCloser(bytes.NewReader([]byte("testResponseBody"))),
					ContentLength: int64(len("testResponseBody")),
				},
			},
			&action{
				status: http.StatusOK,
				body:   "testResponseBody",
				reqKVs: []any{
					"request",
					map[string]any{
						keyID:     "",
						keyTime:   "",
						keyMethod: "GET",
						keyPath:   "/get",
						keyQuery:  "foo=bar&alice=bob",
						keyHost:   "test.com",
						keyRemote: "",
						keyProto:  "HTTP/1.1",
						keyHeader: map[string]string{},
						keySize:   int64(0),
						keyBody:   "",
					},
				},
				resKVs: []any{
					"response",
					map[string]any{
						keyID:       "",
						keyTime:     "",
						keyDuration: int64(0),
						keyStatus:   200,
						keyHeader:   map[string]string{"Test": "value", "Content-Type": "text/plain"},
						keySize:     int64(16), // = len("testResponseBody")
						keyBody:     "testResponseBody",
					},
				},
			},
		),
		gen(
			"record POST",
			&condition{
				logger: &journalLogger{
					req:  &baseLogger{allHeaders: true, maxBody: 100, mimes: []string{"text/plain"}},
					res:  &baseLogger{allHeaders: true, maxBody: 100, mimes: []string{"text/plain"}},
					zone: time.Local,
				},
				req: testPostReq(),
				res: &http.Response{
					StatusCode:    http.StatusOK,
					Header:        http.Header{"Test": {"value"}, "Content-Type": {"text/plain"}},
					Body:          io.NopCloser(bytes.NewReader([]byte("testResponseBody"))),
					ContentLength: int64(len("testResponseBody")),
				},
			},
			&action{
				status: http.StatusOK,
				body:   "testResponseBody",
				reqKVs: []any{
					"request",
					map[string]any{
						keyID:     "",
						keyTime:   "",
						keyMethod: "POST",
						keyPath:   "/post",
						keyQuery:  "foo=bar&alice=bob",
						keyHost:   "test.com",
						keyRemote: "",
						keyProto:  "HTTP/1.1",
						keyHeader: map[string]string{"Content-Type": "text/plain"},
						keySize:   int64(15),
						keyBody:   "testRequestBody",
					},
				},
				resKVs: []any{
					"response",
					map[string]any{
						keyID:       "",
						keyTime:     "",
						keyDuration: int64(0),
						keyStatus:   200,
						keyHeader:   map[string]string{"Test": "value", "Content-Type": "text/plain"},
						keySize:     int64(16), // = len("testResponseBody")
						keyBody:     "testResponseBody",
					},
				},
			},
		),
		gen(
			"reader create error",
			&condition{
				logger: &journalLogger{
					lg: log.GlobalLogger(log.DefaultLoggerName),
					eh: utilhttp.GlobalErrorHandler(utilhttp.DefaultErrorHandlerName),
					req: &baseLogger{
						lg:         log.GlobalLogger(log.DefaultLoggerName),
						allHeaders: true,
						maxBody:    1,
						mimes:      []string{"text/plain"},
						bodyPath:   "./\r\n\t\x00",
					},
					res: &baseLogger{
						allHeaders: true,
						maxBody:    100,
						mimes:      []string{"text/plain"},
					},
					zone: time.Local,
				},
				req: testPostReq(),
				res: &http.Response{
					StatusCode:    http.StatusOK,
					Header:        http.Header{"Test": {"value"}, "Content-Type": {"text/plain"}},
					Body:          io.NopCloser(bytes.NewReader([]byte("testResponseBody"))),
					ContentLength: int64(len("testResponseBody")),
				},
			},
			&action{
				status: 0,  // Not checked.
				body:   ``, // Not checked.
				reqKVs: nil,
				resKVs: nil,
			},
		),
		gen(
			"writer create error",
			&condition{
				logger: &journalLogger{
					lg: log.GlobalLogger(log.DefaultLoggerName),
					eh: utilhttp.GlobalErrorHandler(utilhttp.DefaultErrorHandlerName),
					req: &baseLogger{
						lg:         log.GlobalLogger(log.DefaultLoggerName),
						allHeaders: true,
						maxBody:    100,
						mimes:      []string{"text/plain"},
					},
					res: &baseLogger{
						allHeaders: true,
						maxBody:    1,
						mimes:      []string{"text/plain"},
						bodyPath:   "./\r\n\t\x00",
					},
					zone: time.Local,
				},
				req: testPostReq(),
				res: &http.Response{
					StatusCode:    http.StatusOK,
					Header:        http.Header{"Test": {"value"}, "Content-Type": {"text/plain"}},
					Body:          io.NopCloser(bytes.NewReader([]byte("testResponseBody"))),
					ContentLength: int64(len("testResponseBody")),
				},
			},
			&action{
				status: http.StatusOK,
				body:   "testResponseBody",
				reqKVs: []any{
					"request",
					map[string]any{
						keyID:     "",
						keyTime:   "",
						keyMethod: "POST",
						keyPath:   "/post",
						keyQuery:  "foo=bar&alice=bob",
						keyHost:   "test.com",
						keyRemote: "",
						keyProto:  "HTTP/1.1",
						keyHeader: map[string]string{"Content-Type": "text/plain"},
						keySize:   int64(15),
						keyBody:   "testRequestBody",
					},
				},
				resKVs: []any{
					"response",
					map[string]any{
						keyID:       "",
						keyTime:     "",
						keyDuration: int64(0),
						keyStatus:   200,
						keyHeader: map[string]string{
							"Test":         "value",
							"Content-Type": "text/plain",
						},
						keySize: int64(16), // = len("testResponseBody")
						keyBody: "",        // Not written by error
					},
				},
			},
		),
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			r := core.RoundTripperFunc(func(r *http.Request) (*http.Response, error) {
				return tt.C.res, tt.C.err
			})

			reqLg := &testLogger{Logger: log.GlobalLogger(log.DefaultLoggerName)}
			resLg := &testLogger{Logger: log.GlobalLogger(log.DefaultLoggerName)}
			tt.C.logger.req.lg, tt.C.logger.req.w = reqLg, reqLg
			tt.C.logger.res.lg, tt.C.logger.res.w = resLg, resLg

			roundTripper := tt.C.logger.Tripperware(r)
			w, err := roundTripper.RoundTrip(tt.C.req)
			if err != nil {
				testutil.Diff(t, (*http.Response)(nil), w)
				return
			}

			body, _ := io.ReadAll(w.Body)
			testutil.Diff(t, tt.A.status, w.StatusCode)
			testutil.Diff(t, tt.A.body, string(body))

			opts := []cmp.Option{
				cmpopts.EquateEmpty(),
				cmpopts.IgnoreMapEntries(func(k string, _ any) bool {
					return k == keyDuration || k == keyID || k == keyTime
				}),
			}
			testutil.Diff(t, tt.A.reqFmtLog, string(reqLg.b))
			testutil.Diff(t, tt.A.resFmtLog, string(resLg.b))
			testutil.Diff(t, tt.A.reqKVs, reqLg.kvs, opts...)
			testutil.Diff(t, tt.A.resKVs, resLg.kvs, opts...)
		})
	}
}

func TestLogger_isCompressed(t *testing.T) {
	tests := []struct {
		name    string
		headers http.Header
		expect  bool
	}{
		{
			name:    "No Content-Encoding header",
			headers: http.Header{},
			expect:  false,
		},
		{
			name:    "Gzip encoding",
			headers: http.Header{"Content-Encoding": []string{"gzip"}},
			expect:  true,
		},
		{
			name:    "Brotli encoding",
			headers: http.Header{"Content-Encoding": []string{"br"}},
			expect:  true,
		},
		{
			name:    "Multiple encodings including gzip",
			headers: http.Header{"Content-Encoding": []string{"gzip, deflate"}},
			expect:  true,
		},
		{
			name:    "Unrelated encoding",
			headers: http.Header{"Content-Encoding": []string{"identity"}},
			expect:  false,
		},
		{
			name:    "Zstd encoding",
			headers: http.Header{"Content-Encoding": []string{"zstd"}},
			expect:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isCompressed(tt.headers)
			testutil.Diff(t, tt.expect, result)
		})
	}
}
