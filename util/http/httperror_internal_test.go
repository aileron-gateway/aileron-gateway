// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package http

import (
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestNewHTTPError(t *testing.T) {
	type condition struct {
		err    error
		status int
	}

	type action struct {
		err *HTTPError
	}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"with non-nil error", &condition{
				err:    io.EOF,
				status: http.StatusOK,
			},
			&action{
				err: &HTTPError{
					inner:  io.EOF,
					status: http.StatusOK,
					header: make(http.Header),
				},
			},
		),
		gen(
			"with nil error", &condition{
				err:    nil,
				status: http.StatusOK,
			},
			&action{
				err: &HTTPError{
					inner:  nil,
					status: http.StatusOK,
					header: make(http.Header),
				},
			},
		),
		gen(
			"status 0", &condition{
				err:    io.EOF,
				status: 0,
			},
			&action{
				err: &HTTPError{
					inner:  io.EOF,
					status: 0,
					header: make(http.Header),
				},
			},
		),
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			e := NewHTTPError(tt.C.err, tt.C.status)

			opts := []cmp.Option{
				cmp.AllowUnexported(HTTPError{}),
				cmpopts.IgnoreFields(HTTPError{}, "inner"),
			}
			testutil.Diff(t, tt.A.err, e, opts...)
			testutil.Diff(t, tt.A.err.inner, tt.C.err, cmpopts.EquateErrors())
		})
	}
}

func TestHTTPError(t *testing.T) {
	type condition struct {
		err    *HTTPError
		errs   []*ErrorElem
		accept string
	}

	type action struct {
		err         any // error or errorutil.Kind
		errStr      string
		status      int
		contentType string
		header      http.Header
		body        string
	}

	testErr := errors.New("test error")

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"status 0", &condition{
				err: &HTTPError{
					status: 0,
				},
			},
			&action{
				errStr:      "http status 0 ",
				status:      0,
				contentType: "application/json",
				body:        `{"status":0,"statusText":""}`,
			},
		),
		gen(
			"no internal error", &condition{
				err: &HTTPError{
					status: http.StatusOK,
				},
			},
			&action{
				errStr:      "http status 200 OK",
				status:      http.StatusOK,
				contentType: "application/json",
				body:        `{"status":200,"statusText":"OK"}`,
			},
		),
		gen(
			"internal error", &condition{
				err: &HTTPError{
					inner:  testErr,
					status: http.StatusOK,
				},
			},
			&action{
				err:         testErr,
				errStr:      "test error",
				status:      http.StatusOK,
				contentType: "application/json",
				body:        `{"status":200,"statusText":"OK"}`,
			},
		),
		gen(
			"unsupported mime", &condition{
				err: &HTTPError{
					status: http.StatusOK,
				},
				accept: "text/html",
			},
			&action{
				errStr:      "http status 200 OK",
				status:      http.StatusOK,
				contentType: "application/json",
				body:        `{"status":200,"statusText":"OK"}`,
			},
		),
		gen(
			"non empty header mime", &condition{
				err: &HTTPError{
					status: http.StatusOK,
					header: http.Header{"Foo": []string{"bar"}},
				},
			},
			&action{
				errStr:      "http status 200 OK",
				status:      http.StatusOK,
				contentType: "application/json",
				header:      http.Header{"Foo": []string{"bar"}},
				body:        `{"status":200,"statusText":"OK"}`,
			},
		),
		gen(
			"application/json", &condition{
				err: &HTTPError{
					status: http.StatusOK,
				},
				accept: "application/json",
			},
			&action{
				errStr:      "http status 200 OK",
				status:      http.StatusOK,
				contentType: "application/json",
				body:        `{"status":200,"statusText":"OK"}`,
			},
		),
		gen(
			"application/json with elem", &condition{
				err: &HTTPError{
					status: http.StatusOK,
				},
				errs:   []*ErrorElem{{Code: "c", Message: "m", Detail: "d"}},
				accept: "application/json",
			},
			&action{
				errStr:      "http status 200 OK",
				status:      http.StatusOK,
				contentType: "application/json",
				body:        `{"status":200,"statusText":"OK","errors":[{"code":"c","message":"m","detail":"d"}]}`,
			},
		),
		gen(
			"application/xml", &condition{
				err: &HTTPError{
					status: http.StatusOK,
				},
				accept: "application/xml",
			},
			&action{
				errStr:      "http status 200 OK",
				status:      http.StatusOK,
				contentType: "application/xml",
				body:        `<result><status>200</status><statusText>OK</statusText><errors></errors></result>`,
			},
		),
		gen(
			"application/xml with elem", &condition{
				err: &HTTPError{
					status: http.StatusOK,
				},
				errs:   []*ErrorElem{{Code: "c", Message: "m", Detail: "d"}},
				accept: "application/xml",
			},
			&action{
				errStr:      "http status 200 OK",
				status:      http.StatusOK,
				contentType: "application/xml",
				body:        `<result><status>200</status><statusText>OK</statusText><errors><error><code>c</code><message>m</message><detail>d</detail></error></errors></result>`,
			},
		),
		gen(
			"text/plain", &condition{
				err: &HTTPError{
					status: http.StatusOK,
				},
				accept: "text/plain",
			},
			&action{
				errStr:      "http status 200 OK",
				status:      http.StatusOK,
				contentType: "text/plain",
				body:        "status: 200\nstatusText: OK\n",
			},
		),
		gen(
			"text/plain with elem", &condition{
				err: &HTTPError{
					status: http.StatusOK,
				},
				errs:   []*ErrorElem{{Code: "c", Message: "m", Detail: "d"}},
				accept: "text/plain",
			},
			&action{
				errStr:      "http status 200 OK",
				status:      http.StatusOK,
				contentType: "text/plain",
				body:        "status: 200\nstatusText: OK\nerrors:\n    - code: c\n      message: m\n      detail: d\n",
			},
		),
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			e := tt.C.err
			for _, elem := range tt.C.errs {
				e.AddError(elem)
			}
			contentType, body := e.Content(tt.C.accept)
			testutil.Diff(t, tt.A.err, e.Unwrap(), cmpopts.EquateErrors())
			testutil.Diff(t, tt.A.errStr, e.Error())
			testutil.Diff(t, tt.A.status, e.StatusCode())
			testutil.Diff(t, tt.A.contentType, contentType)
			testutil.Diff(t, tt.A.header, e.Header())
			testutil.Diff(t, tt.A.body, string(body))
		})
	}
}
