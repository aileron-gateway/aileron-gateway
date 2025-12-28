// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package soaprest

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"io"
	"mime"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"regexp"
	"slices"
	"strings"
	"testing"

	"github.com/aileron-gateway/aileron-gateway/app"
	"github.com/aileron-gateway/aileron-gateway/core"
	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
	"github.com/aileron-projects/go/zencoding/zxml"
	gocmp "github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

type mockErrorHandler struct {
	err  error
	code int
}

func (m *mockErrorHandler) ServeHTTPError(w http.ResponseWriter, r *http.Request, err error) {
	m.err = err
	m.code = err.(core.HTTPError).StatusCode()
	w.WriteHeader(m.code)
}

type mockReader struct{}

func (m *mockReader) Read(p []byte) (int, error) {
	return 0, errors.New("mock read error")
}

func equalJSON(t *testing.T, a, b []byte) bool {
	var obj1, obj2 interface{}
	if err := json.Unmarshal(a, &obj1); err != nil {
		panic(err)
	}
	if err := json.Unmarshal(b, &obj2); err != nil {
		panic(err)
	}
	if equal := reflect.DeepEqual(obj1, obj2); equal {
		return true
	}
	t.Logf("JSON-1: %#v\n", obj1)
	t.Logf("JSON-2: %#v\n", obj2)
	return false
}

func TestSOAPREST_Middleware_RequestConversion(t *testing.T) {
	type condition struct {
		file          string
		method        string
		contentType   string
		converter     *zxml.JSONConverter
		readBodyError bool
		pathNotMatch  bool
		setSOAPAction bool
	}

	type action struct {
		file       string
		err        any // error or errorKind
		errPattern *regexp.Regexp
		code       int
	}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"SOAP 1.1 request with SOAPAction",
			&condition{
				file:        "ok_case01",
				method:      http.MethodPost,
				contentType: "text/xml",
				converter: &zxml.JSONConverter{
					EncodeDecoder: zxml.NewSimple(),
				},
				setSOAPAction: true,
			},
			&action{
				file: "ok_case01",
				// This case: no errors on the request side and the upstream handler writes no body,
				// so the response Content-Type remains the default (not application/json).
				// As a result, the middleware should abort with an InvalidContentType error.
				err:  app.ErrAppMiddleSOAPRESTInvalidContentType,
				code: 500,
			},
		),
		gen(
			"SOAP 1.1 request with SOAPAction and Get method",
			&condition{
				file:        "ok_case01",
				method:      http.MethodGet,
				contentType: "text/xml",
				converter: &zxml.JSONConverter{
					EncodeDecoder: zxml.NewSimple(),
				},
				setSOAPAction: true,
			},
			&action{
				file: "ok_case01",

				// This case: no errors on the request side and the upstream handler writes no body,
				// so the response Content-Type remains the default (not application/json).
				// As a result, the middleware should abort with an InvalidContentType error.
				err:  app.ErrAppMiddleSOAPRESTInvalidContentType,
				code: 500,
			},
		),
		gen(
			"SOAP 1.1 request without SOAPAction",
			&condition{
				file:        "ok_case01",
				method:      http.MethodPost,
				contentType: "text/xml",
				converter: &zxml.JSONConverter{
					EncodeDecoder: zxml.NewSimple(),
				},
				setSOAPAction: false,
			},
			&action{
				file: "ok_case01",

				err:  app.ErrAppMiddleSOAPRESTVersionMismatch,
				code: 403,
			},
		),
		gen(
			"SOAP1.2 request",
			&condition{
				file:        "ok_case02",
				method:      http.MethodPost,
				contentType: "application/soap+xml",
				converter: &zxml.JSONConverter{
					EncodeDecoder: zxml.NewSimple(),
				},
			},
			&action{
				file: "ok_case02",
				// This case: no errors on the request side and the upstream handler writes no body,
				// so the response Content-Type remains the default (not application/json).
				// As a result, the middleware should abort with an InvalidContentType error.
				err:  app.ErrAppMiddleSOAPRESTInvalidContentType,
				code: 500,
			},
		),
		gen(
			"SOAP1.2 request with action",
			&condition{
				file:        "ok_case02",
				method:      http.MethodPost,
				contentType: "application/soap+xml",
				converter: &zxml.JSONConverter{
					EncodeDecoder: zxml.NewSimple(),
				},
				setSOAPAction: true,
			},
			&action{
				file: "ok_case02",
				// This case: no errors on the request side and the upstream handler writes no body,
				// so the response Content-Type remains the default (not application/json).
				// As a result, the middleware should abort with an InvalidContentType error.
				err:  app.ErrAppMiddleSOAPRESTInvalidContentType,
				code: 500,
			},
		),
		gen(
			"non SOAP request",
			&condition{
				file:        "ng_case01",
				method:      http.MethodPost,
				contentType: "application/json",
			},
			&action{
				file: "ng_case01",
				err:  app.ErrAppMiddleSOAPRESTVersionMismatch,
				code: 403,
			},
		),
		gen(
			"read body error",
			&condition{
				file:          "ng_case01",
				method:        http.MethodPost,
				contentType:   "text/xml",
				readBodyError: true,
				setSOAPAction: true,
			},
			&action{
				file: "ng_case01",
				err:  app.ErrAppMiddleSOAPRESTReadXMLBody,
				code: 400,
			},
		),
		gen(
			"request convert error",
			&condition{
				file:        "ng_case02",
				method:      http.MethodPost,
				contentType: "text/xml",
				converter: &zxml.JSONConverter{
					EncodeDecoder: zxml.NewSimple(),
					Header:        xml.Header,
				},

				setSOAPAction: true,
			},
			&action{
				file: "ng_case02",
				err:  app.ErrAppMiddleSOAPRESTConvertXMLtoJSON,
				code: 400,
			},
		),
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			meh := &mockErrorHandler{}
			sr := &soapREST{
				eh:        meh,
				converter: tt.C.converter,
			}

			nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				bodyBytes, err := io.ReadAll(r.Body)
				if err != nil {
					t.Fatalf("Failed to read request body: %v", err)
				}
				r.Body.Close()

				// Check whether the SOAP/XML request is being converted
				// into a REST/JSON request by the SOAPRESTMiddleware.
				jsonBytes, _ := os.ReadFile("./testdata/Simple/" + tt.A.file + ".json")
				if !equalJSON(t, bodyBytes, jsonBytes) {
					t.Error("decode result not match (xml to json)")
				}

				// Verify that the Content-Type is correctly modified.
				mt, _, _ := mime.ParseMediaType(r.Header.Get("Content-Type"))
				testutil.Diff(t, "application/json", mt)
				// Verify that the Accept header is set to application/json.
				testutil.Diff(t, "application/json", r.Header.Get("Accept"))
				// Verify that the request method is preserved.
				testutil.Diff(t, tt.C.method, r.Method)

				// Verify that the original Content-Type is preserved in X-Content-Type
				if tt.C.contentType == "text/xml" || tt.C.contentType == "application/soap+xml" {
					testutil.Diff(t, tt.C.contentType, r.Header.Get("X-Content-Type"))
				}

				w.WriteHeader(http.StatusOK)
			})

			h := sr.Middleware(nextHandler)

			// If the path does not match, the XML to JSON conversion will not be performed,
			// so the JSON will be stored in the request body.
			xmlBytes, _ := os.ReadFile("./testdata/xml/" + tt.C.file + ".xml")
			if tt.C.pathNotMatch {
				xmlBytes, _ = os.ReadFile("./testdata/Simple/" + tt.C.file + ".json")
			}

			req := httptest.NewRequest(tt.C.method, "http://test.com/", strings.NewReader(string(xmlBytes)))
			if tt.C.readBodyError {
				req.Body = io.NopCloser(&mockReader{})
			}

			if tt.C.contentType == "application/soap+xml" && tt.C.setSOAPAction {
				req.Header.Set("Content-Type", tt.C.contentType+"; charset=utf-8; action=\"http://example.com/\"")
			} else {
				req.Header.Set("Content-Type", tt.C.contentType+"; charset=utf-8")
			}
			if tt.C.contentType == "text/xml" && tt.C.setSOAPAction {
				req.Header.Set("SOAPAction", "http://example.com/")
			}

			resp := httptest.NewRecorder()
			h.ServeHTTP(resp, req)

			opts := []gocmp.Option{
				cmpopts.EquateErrors(),
			}
			testutil.DiffError(t, tt.A.err, tt.A.errPattern, meh.err, opts...)
			testutil.Diff(t, tt.A.code, resp.Code)
		})
	}
}

type errorResponseRecorder struct {
	header     http.Header
	code       int
	writeError error
}

func (rec *errorResponseRecorder) Header() http.Header {
	if rec.header == nil {
		rec.header = make(http.Header)
	}
	return rec.header
}

func (rec *errorResponseRecorder) WriteHeader(code int) {
	rec.code = code
}

func (rec *errorResponseRecorder) Write(b []byte) (int, error) {
	rec.writeError = errors.New("mock write error")
	return 0, rec.writeError
}

func xmlTokens(decoder *xml.Decoder, end xml.EndElement) (map[string]any, error) {
	key := ""
	m := map[string]any{}
	var tokens []xml.Token
	for {
		token, err := decoder.Token()
		if err != nil {
			if err == io.EOF {
				return m, nil
			}
			return m, err
		}
		switch t := token.(type) {
		case xml.Comment, xml.ProcInst, xml.Directive:
			continue
		case xml.StartElement:
			slices.SortFunc(t.Attr, func(a, b xml.Attr) int {
				if a.Name.Space+a.Name.Local > b.Name.Space+b.Name.Local {
					return 1
				}
				return -1
			})
			key = t.Name.Space + ":" + t.Name.Local
			children, err := xmlTokens(decoder, t.End())
			if err != nil {
				return nil, err
			}
			m[key] = children
		case xml.CharData:
			t = bytes.TrimSpace([]byte(t))
			if len(t) == 0 {
				continue
			}
			token = t
		case xml.EndElement:
			if t == end {
				return m, nil
			}
			v, ok := m[key]
			if ok {
				vv := v.([]xml.Token)
				m[key] = append(vv, tokens...)
			} else {
				m[key] = tokens
			}
			tokens = nil
		}
	}
}

func equalXML(t *testing.T, a, b []byte) bool {
	tokens1, err := xmlTokens(xml.NewDecoder(bytes.NewReader(a)), xml.EndElement{})
	if err != nil {
		panic(err)
	}
	tokens2, err := xmlTokens(xml.NewDecoder(bytes.NewReader(b)), xml.EndElement{})
	if err != nil {
		panic(err)
	}
	if equal := reflect.DeepEqual(tokens1, tokens2); equal {
		return true
	}
	t.Logf("XML-1: %#v\n", tokens1)
	t.Logf("XML-2: %#v\n", tokens2)
	return false
}

func TestSOAPREST_Middleware_ResponseConversion(t *testing.T) {
	type condition struct {
		file        string
		method      string
		contentType string
		charset     string

		converter *zxml.JSONConverter

		invalidContentTypeError bool
		responseConvertError    bool
		responseWriteError      bool
		setSOAPAction           bool
	}

	type action struct {
		file        string
		contentType string
		actionURI   string

		err        any // error or errorKind
		errPattern *regexp.Regexp
		code       int
	}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"SOAP1.1 response",
			&condition{
				file:        "ok_case01",
				method:      http.MethodPost,
				contentType: "text/xml",
				charset:     "utf-8",
				converter: &zxml.JSONConverter{
					EncodeDecoder: zxml.NewSimple(),
				},
				setSOAPAction: true,
			},
			&action{
				file:        "ok_case01",
				contentType: "text/xml; charset=utf-8",
				actionURI:   "http://example.com/",
				err:         nil,
				code:        200,
			},
		),
		gen(
			"SOAP1.2 response, request with SOAPAction",
			&condition{
				file:        "ok_case02",
				method:      http.MethodPost,
				contentType: "application/soap+xml",
				charset:     "utf-8",
				converter: &zxml.JSONConverter{
					EncodeDecoder: zxml.NewSimple(),
				},
				setSOAPAction: true,
			},
			&action{
				file:        "ok_case02",
				contentType: "application/soap+xml; charset=utf-8",
				actionURI:   "http://example.com/",
				err:         nil,
				code:        200,
			},
		),
		gen(
			"SOAP1.2 response, request without SOAPAction",
			&condition{
				file:        "ok_case02",
				method:      http.MethodPost,
				contentType: "application/soap+xml",
				charset:     "utf-8",
				converter: &zxml.JSONConverter{
					EncodeDecoder: zxml.NewSimple(),
				},
				setSOAPAction: false,
			},
			&action{
				file:        "ok_case02",
				contentType: "application/soap+xml; charset=utf-8",
				actionURI:   "",
				err:         nil,
				code:        200,
			},
		),
		gen(
			"invalid Content-Type error",
			&condition{
				file:        "ng_case03",
				method:      http.MethodPost,
				contentType: "text/xml",
				charset:     "utf-8",
				converter: &zxml.JSONConverter{
					EncodeDecoder: zxml.NewSimple(),
				},
				invalidContentTypeError: true,
				setSOAPAction:           true,
			},
			&action{
				// Don't specify a comparison file when conversion fails, as the Body will be empty.
				file:        "notExist",
				contentType: "invalid",
				actionURI:   "http://example.com/",
				err:         app.ErrAppMiddleSOAPRESTInvalidContentType,
				code:        500,
			},
		),
		gen(
			"response convert error",
			&condition{
				file:        "ng_case03",
				method:      http.MethodPost,
				contentType: "text/xml",
				charset:     "utf-8",
				converter: &zxml.JSONConverter{
					EncodeDecoder: zxml.NewSimple(),
				},
				responseConvertError: true,
				setSOAPAction:        true,
			},
			&action{
				// Don't specify a comparison file when conversion fails, as the Body will be empty.
				file:        "notExist",
				contentType: "application/json; charset=utf-8",
				actionURI:   "http://example.com/",
				err:         app.ErrAppMiddleSOAPRESTConvertJSONtoXML,
				code:        500,
			},
		),
		gen(
			"responseWriter write error",
			&condition{
				file:        "ng_case01",
				method:      http.MethodPost,
				contentType: "text/xml",
				charset:     "utf-8",
				converter: &zxml.JSONConverter{
					EncodeDecoder: zxml.NewSimple(),
				},
				responseWriteError: true,
				setSOAPAction:      true,
			},
			&action{
				// Don't specify a comparison file when conversion fails, as the Body will be empty.
				file:        "notExist",
				contentType: "text/xml; charset=utf-8",
				actionURI:   "http://example.com/",

				err:  app.ErrAppMiddleSOAPRESTWriteResponseBody,
				code: 500,
			},
		),
		gen(
			"empty charset with SOAP1.1 request",
			&condition{
				file:        "ok_case01",
				method:      http.MethodPost,
				contentType: "text/xml",
				converter: &zxml.JSONConverter{
					EncodeDecoder: zxml.NewSimple(),
				},
				setSOAPAction: true,
			},
			&action{
				file:        "ok_case01",
				contentType: "text/xml; charset=utf-8",
				actionURI:   "http://example.com/",

				err:  nil,
				code: 200,
			},
		),
		gen(
			"empty charset with SOAP1.2 request",
			&condition{
				file:        "ok_case01",
				method:      http.MethodPost,
				contentType: "application/soap+xml",
				converter: &zxml.JSONConverter{
					EncodeDecoder: zxml.NewSimple(),
				},
				setSOAPAction: true,
			},
			&action{
				file:        "ok_case01",
				contentType: "application/soap+xml; charset=utf-8",
				actionURI:   "http://example.com/",

				err:  nil,
				code: 200,
			},
		),
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			meh := &mockErrorHandler{}
			m := &soapREST{
				eh:        meh,
				converter: tt.C.converter,
			}

			nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json; charset="+tt.C.charset)
				if tt.C.invalidContentTypeError {
					w.Header().Set("Content-Type", "invalid")
				}

				if !tt.C.responseConvertError {
					jsonBytes, _ := os.ReadFile("./testdata/Simple/" + tt.C.file + ".json")
					w.Write([]byte(jsonBytes))
				}

				// Verify that the action value is stored in the custom header.
				testutil.Diff(t, tt.A.actionURI, r.Header.Get("X-SOAP-Action"))
				w.WriteHeader(http.StatusOK)
			})

			h := m.Middleware(nextHandler)
			req := httptest.NewRequest(tt.C.method, "http://test.com/", nil)

			switch tt.C.contentType {
			case "application/soap+xml":
				if tt.C.setSOAPAction {
					req.Header.Set("Content-Type", tt.C.contentType+"; charset=utf-8; action=\"http://example.com/\"")
				} else {
					req.Header.Set("Content-Type", tt.C.contentType+"; charset=utf-8")
				}
			case "text/xml":
				if tt.C.setSOAPAction {
					req.Header.Set("SOAPAction", "http://example.com/")
				}
				req.Header.Set("Content-Type", tt.C.contentType+"; charset=utf-8")
			}

			var resp http.ResponseWriter
			if tt.C.responseWriteError {
				resp = &errorResponseRecorder{}
			} else {
				resp = httptest.NewRecorder()
			}

			h.ServeHTTP(resp, req)

			if rec, ok := resp.(*httptest.ResponseRecorder); ok {
				testutil.Diff(t, tt.A.code, rec.Code)
				xmlBytes, _ := os.ReadFile("./testdata/xml/" + tt.C.file + ".xml")
				if !equalXML(t, xmlBytes, rec.Body.Bytes()) {
					t.Error("encode result not match (json to xml)")
				}
			} else if rec, ok := resp.(*errorResponseRecorder); ok {
				testutil.Diff(t, tt.A.code, rec.code)
			}

			testutil.Diff(t, tt.A.contentType, resp.Header().Get("Content-Type"))

			opts := []gocmp.Option{
				cmpopts.EquateErrors(),
			}
			testutil.DiffError(t, tt.A.err, tt.A.errPattern, meh.err, opts...)
		})
	}
}

type mockResponseWriter struct {
	http.ResponseWriter
	id string
}

func TestWrappedWriter_Unwrap(t *testing.T) {
	type condition struct {
		ww *wrappedWriter
	}

	type action struct {
		w http.ResponseWriter
	}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"unwrap nil",
			&condition{
				ww: &wrappedWriter{
					ResponseWriter: nil,
				},
			},
			&action{
				w: nil,
			},
		),
		gen(
			"unwrap non-nil",
			&condition{
				ww: &wrappedWriter{
					ResponseWriter: &mockResponseWriter{
						id: "inner",
					},
				},
			},
			&action{
				w: &mockResponseWriter{
					id: "inner",
				},
			},
		),
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			w := tt.C.ww.Unwrap()
			testutil.Diff(t, tt.A.w, w, gocmp.AllowUnexported(mockResponseWriter{}))
		})
	}
}

// In the wrappedWriter, writing directly to the response is not performed;
// instead, it simply holds values in the wrappedWriter structure.
// Therefore, it does not check whether the statusCode stored in http.ResponseWriter matches the statusCode of the condition.
func TestWrappedWriter_WriteHeader(t *testing.T) {
	type condition struct {
		ww      *wrappedWriter
		code    int
		written bool
	}

	type action struct {
		code    int
		written bool
	}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"status code 100",
			&condition{
				ww: &wrappedWriter{
					ResponseWriter: httptest.NewRecorder(),
				},
				code:    100,
				written: false,
			},
			&action{
				code:    100,
				written: true,
			},
		),
		gen(
			"status code 999",
			&condition{
				ww: &wrappedWriter{
					ResponseWriter: httptest.NewRecorder(),
				},
				code:    999,
				written: false,
			},
			&action{
				code:    999,
				written: true,
			},
		),
		gen(
			"written wrappedwriter",
			&condition{
				ww: &wrappedWriter{
					ResponseWriter: httptest.NewRecorder(),
				},
				written: true,
			},
			&action{
				code:    0,
				written: true,
			},
		),
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			w := httptest.NewRecorder()
			ww := &wrappedWriter{
				ResponseWriter: w,
				written:        tt.C.written,
			}
			ww.WriteHeader(tt.C.code)

			testutil.Diff(t, tt.A.code, ww.code)
			testutil.Diff(t, tt.A.written, ww.written)
		})
	}
}

func TestWrappedWriter_Write(t *testing.T) {
	type condition struct {
		code int
		body string
	}

	type action struct {
		code int
		body string
	}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"status code 100",
			&condition{
				code: 100,
				body: "test",
			},
			&action{
				code: 100,
				body: "test",
			},
		),
		gen(
			"status code 999",
			&condition{
				code: 999,
				body: "test",
			},
			&action{
				code: 999,
				body: "test",
			},
		),
		gen(
			"status code 0 (don't write the code)",
			&condition{
				code: 0,
				body: "test",
			},
			&action{
				code: 0,
				body: "test",
			},
		),
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			w := httptest.NewRecorder()
			ww := &wrappedWriter{
				ResponseWriter: w,
				body:           &bytes.Buffer{},
			}
			if tt.C.code > 0 {
				ww.WriteHeader(tt.C.code)
			}
			ww.Write([]byte(tt.C.body))

			testutil.Diff(t, tt.A.code, ww.code)
			body, _ := io.ReadAll(ww.body)
			testutil.Diff(t, tt.A.body, string(body))
		})
	}
}

func TestWrappedWriter_Written(t *testing.T) {
	type condition struct {
		ww    *wrappedWriter
		write bool
	}

	type action struct {
		written bool
	}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"don't write status code",
			&condition{
				ww: &wrappedWriter{
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
			&condition{
				ww: &wrappedWriter{
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
			testutil.Diff(t, tt.A.written, tt.C.ww.written)
		})
	}
}

func TestWrappedWriter_StatusCode(t *testing.T) {
	type condition struct {
		ww   *wrappedWriter
		code int
	}

	type action struct {
		code int
	}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"status code 100",
			&condition{
				ww: &wrappedWriter{
					ResponseWriter: httptest.NewRecorder(),
				},
				code: 100,
			},
			&action{
				code: 100,
			},
		),
		gen(
			"status code 999",
			&condition{
				ww: &wrappedWriter{
					ResponseWriter: httptest.NewRecorder(),
				},
				code: 999,
			},
			&action{
				code: 999,
			},
		),
		gen(
			"written is false and code is 0",
			&condition{
				ww: &wrappedWriter{
					ResponseWriter: httptest.NewRecorder(),
				},
				code: 0,
			},
			&action{
				code: 200,
			},
		),
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			tt.C.ww.WriteHeader(tt.C.code)
			testutil.Diff(t, tt.A.code, tt.C.ww.StatusCode())
		})
	}
}

func TestWrappedWriter_Flush(t *testing.T) {
	type condition struct{}
	type action struct{}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"no-op test",
			&condition{},
			&action{},
		),
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			ww := wrappedWriter{}
			ww.Flush()
		})
	}
}
