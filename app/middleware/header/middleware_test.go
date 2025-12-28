// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package header

import (
	"net/http"
	"net/http/httptest"
	"testing"

	k "github.com/aileron-gateway/aileron-gateway/apis/kernel"
	"github.com/aileron-gateway/aileron-gateway/internal/testutil"
	"github.com/aileron-gateway/aileron-gateway/kernel/txtutil"
	httputil "github.com/aileron-gateway/aileron-gateway/util/http"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestWrappedWriter_Unwrap(t *testing.T) {
	type condition struct {
		ww *wrappedWriter
	}

	type action struct {
		w http.ResponseWriter
	}

	testWriter := &httptest.ResponseRecorder{}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"nil writer",
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
			"non nil writer",
			&condition{
				ww: &wrappedWriter{
					ResponseWriter: testWriter,
				},
			},
			&action{
				w: testWriter,
			},
		),
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			w := tt.C.ww.Unwrap()
			testutil.Diff(t, tt.A.w, w, cmp.Comparer(testutil.ComparePointer[http.ResponseWriter]))
		})
	}
}

func TestWrappedWriter_WriteHeader(t *testing.T) {
	type condition struct {
		applied bool
		p       *policy
	}

	type action struct {
		header map[string][]string
	}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"apply policy",
			&condition{
				applied: false,
				p:       &policy{add: map[string]string{"Test": "value"}},
			},
			&action{
				header: map[string][]string{"Test": {"value"}},
			},
		),
		gen(
			"policy already applied",
			&condition{
				applied: true,
				p:       &policy{add: map[string]string{"Test": "value"}},
			},
			&action{
				header: map[string][]string{"Test": nil},
			},
		),
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			resp := httptest.NewRecorder()
			w := &wrappedWriter{
				ResponseWriter: resp,
				applied:        tt.C.applied,
				p:              tt.C.p,
			}

			w.WriteHeader(http.StatusOK)
			testutil.Diff(t, true, w.applied)
			for k, v := range tt.A.header {
				testutil.Diff(t, v, resp.Header()[k])
			}
		})
	}
}

func TestWrappedWriter_Write(t *testing.T) {
	type condition struct {
		applied bool
		p       *policy
	}

	type action struct {
		header map[string][]string
	}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"apply policy",
			&condition{
				applied: false,
				p:       &policy{add: map[string]string{"Test": "value"}},
			},
			&action{
				header: map[string][]string{"Test": {"value"}},
			},
		),
		gen(
			"policy already applied",
			&condition{
				applied: true,
				p:       &policy{add: map[string]string{"Test": "value"}},
			},
			&action{
				header: map[string][]string{"Test": nil},
			},
		),
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			resp := httptest.NewRecorder()
			w := &wrappedWriter{
				ResponseWriter: resp,
				applied:        tt.C.applied,
				p:              tt.C.p,
			}

			w.Write([]byte("test"))
			testutil.Diff(t, true, w.applied)
			for k, v := range tt.A.header {
				testutil.Diff(t, v, resp.Header()[k])
			}
		})
	}
}

func TestPolicy(t *testing.T) {
	type condition struct {
		p      *policy
		header http.Header
	}

	type action struct {
		header http.Header
	}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"allows",
			&condition{
				p: &policy{
					allows: []string{"Foo"},
				},
				header: http.Header{"Foo": {"foo"}, "Bar": {"bar"}},
			},
			&action{
				header: http.Header{"Foo": {"foo"}},
			},
		),
		gen(
			"removes",
			&condition{
				p: &policy{
					removes: []string{"Foo"},
				},
				header: http.Header{"Foo": {"foo"}, "Bar": {"bar"}},
			},
			&action{
				header: http.Header{"Foo": nil, "Bar": {"bar"}},
			},
		),
		gen(
			"add to empty",
			&condition{
				p:      &policy{add: map[string]string{"Test": "value"}},
				header: http.Header{},
			},
			&action{
				header: http.Header{"Test": {"value"}},
			},
		),
		gen(
			"add to existing",
			&condition{
				p:      &policy{add: map[string]string{"Test": "value"}},
				header: http.Header{"Test": {"exist"}},
			},
			&action{
				header: http.Header{"Test": {"exist", "value"}},
			},
		),
		gen(
			"set to empty",
			&condition{
				p:      &policy{set: map[string]string{"Test": "value"}},
				header: http.Header{},
			},
			&action{
				header: http.Header{"Test": {"value"}},
			},
		),
		gen(
			"set to existing",
			&condition{
				p:      &policy{set: map[string]string{"Test": "value"}},
				header: http.Header{"Test": {"exist"}},
			},
			&action{
				header: http.Header{"Test": {"value"}},
			},
		),
		gen(
			"replace",
			&condition{
				p: &policy{
					repls: map[string]txtutil.ReplaceFunc[string]{"Test": func(s string) string { return "***" }},
				},
				header: http.Header{"Test": {"value1", "value2"}},
			},
			&action{
				header: http.Header{"Test": {"***", "***"}},
			},
		),
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			tt.C.p.apply(tt.C.header)
			testutil.Diff(t, tt.A.header, tt.C.header, cmpopts.SortMaps(func(a, b string) bool { return a > b }))
		})
	}
}

func TestHeaderPolicyMiddleware(t *testing.T) {
	type condition struct {
		headerPolicy      headerPolicy
		header            http.Header
		length            int64
		writeResponseBody bool
	}

	type action struct {
		statusCode int
		header     http.Header
	}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"allowed MIME",
			&condition{
				headerPolicy: headerPolicy{
					allowedMIMEs: []string{"application/json"},
					eh:           httputil.GlobalErrorHandler(httputil.DefaultErrorHandlerName),
				},
				header: http.Header{"Content-Type": {"application/json"}},
			},
			&action{
				statusCode: http.StatusOK,
				header:     http.Header{},
			},
		),
		gen(
			"not allowed MIME",
			&condition{
				headerPolicy: headerPolicy{
					allowedMIMEs: []string{"application/xml"},
					eh:           httputil.GlobalErrorHandler(httputil.DefaultErrorHandlerName),
				},
				header: http.Header{"Content-Type": {"application/json"}},
			},
			&action{
				statusCode: http.StatusUnsupportedMediaType,
				header:     http.Header{},
			},
		),
		gen(
			"allowed content length",
			&condition{
				headerPolicy: headerPolicy{
					maxContentLength: 1024,
					eh:               httputil.GlobalErrorHandler(httputil.DefaultErrorHandlerName),
				},
				header: http.Header{},
				length: 512,
			},
			&action{
				statusCode: http.StatusOK,
				header:     http.Header{},
			},
		),
		gen(
			"content length exceeded",
			&condition{
				headerPolicy: headerPolicy{
					maxContentLength: 1024,
					eh:               httputil.GlobalErrorHandler(httputil.DefaultErrorHandlerName),
				},
				header: http.Header{},
				length: 2048,
			},
			&action{
				statusCode: http.StatusRequestEntityTooLarge,
				header:     http.Header{},
			},
		),
		gen(
			"content length required",
			&condition{
				headerPolicy: headerPolicy{
					maxContentLength: 1024, // 1KB
					eh:               httputil.GlobalErrorHandler(httputil.DefaultErrorHandlerName),
				},
				header: http.Header{},
				length: -1,
			},
			&action{
				statusCode: http.StatusLengthRequired,
			},
		),
		gen(
			"apply request header policy",
			&condition{
				headerPolicy: headerPolicy{
					eh: httputil.GlobalErrorHandler(httputil.DefaultErrorHandlerName),
					reqPolicy: &policy{
						add: map[string]string{"Test": "value"},
					},
				},
				header: http.Header{},
			},
			&action{
				statusCode: http.StatusOK,
				header:     http.Header{"Test": {"value"}},
			},
		),
		gen(
			"apply response header policy",
			&condition{
				headerPolicy: headerPolicy{
					eh: httputil.GlobalErrorHandler(httputil.DefaultErrorHandlerName),
					resPolicy: &policy{
						add: map[string]string{"Test": "value"},
					},
				},
				header: http.Header{},
			},
			&action{
				statusCode: http.StatusOK,
				header:     http.Header{"Test": {"value"}},
			},
		),
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.Header = tt.C.header
			req.ContentLength = tt.C.length

			rec := httptest.NewRecorder()
			next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})
			tt.C.headerPolicy.Middleware(next).ServeHTTP(rec, req)

			testutil.Diff(t, tt.A.statusCode, rec.Code)
			if tt.C.headerPolicy.reqPolicy != nil {
				testutil.Diff(t, tt.A.header, tt.C.header, cmpopts.SortMaps(func(a, b string) bool { return a > b }))
			}
			if tt.C.headerPolicy.resPolicy != nil {
				testutil.Diff(t, tt.A.header, rec.Header(), cmpopts.SortMaps(func(a, b string) bool { return a > b }))
			}
		})
	}
}

func createStringReplacer(from string, to string) txtutil.Replacer[string] {
	spec := &k.ReplacerSpec{
		Replacers: &k.ReplacerSpec_Value{
			Value: &k.ValueReplacer{
				FromTo: map[string]string{
					from: to,
				},
			},
		},
	}
	replacer, _ := txtutil.NewStringReplacer(spec)
	return replacer
}
