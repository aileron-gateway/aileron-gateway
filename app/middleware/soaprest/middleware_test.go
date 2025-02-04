package soaprest

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/aileron-gateway/aileron-gateway/core"
	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
	utilhttp "github.com/aileron-gateway/aileron-gateway/util/http"
	"github.com/google/go-cmp/cmp"
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

func TestMiddleware(t *testing.T) {
	type condition struct {
		method      string
		contentType string
		body        string

		readBodyError    bool
		nextHandlerError bool
		jsonMarshalError bool
	}

	type action struct {
		body     string
		err      error
		respCode int
		ehCode   int //
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"NonSOAPRequest",
			[]string{},
			[]string{},
			&condition{
				method:      http.MethodGet,
				contentType: "application/json",
				body:        ``,
			},
			&action{
				err:      errInvalidSOAP11Request,
				respCode: 403,
				ehCode:   403,
			},
		),
		gen(
			"ReadBodyError",
			[]string{},
			[]string{},
			&condition{
				method:      http.MethodPost,
				contentType: "text/xml",
				body:        `<invalid_xml>`,

				readBodyError: true,
			},
			&action{
				err:      utilhttp.ErrBadRequest,
				respCode: 400,
				ehCode:   400,
			},
		),
		gen(
			"UnmarshalError",
			[]string{},
			[]string{},
			&condition{
				method:      http.MethodPost,
				contentType: "text/xml",
				body:        `<invalid_xml>`,

				readBodyError: false,
			},
			&action{
				err:      utilhttp.ErrBadRequest,
				respCode: 400,
				ehCode:   400,
			},
		),
		// <TODO>: need to consider how to generate json.Marshal error
		// gen(
		// 	"MarshalError",
		// 	[]string{},
		// 	[]string{},
		// 	&condition{
		// 		method:      http.MethodPost,
		// 		contentType: "text/xml",
		// 		body: `<?xml version="1.0" encoding="utf-8"?>
		// 			<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/" xmlns:ns="http://example.com/">
		// 			<soap:Header/>
		// 			<soap:Body>
		// 				<ns:Test>
		// 				</ns:Test>
		// 			</soap:Body>
		// 			</soap:Envelope>`,

		// 		jsonMarshalError: true,
		// 	},
		// 	&action{
		// 		err:  utilhttp.ErrBadRequest,
		// 		code: 400,
		// 	},
		// ),
		gen(
			"NextHandlerError",
			[]string{},
			[]string{},
			&condition{
				method:      http.MethodPost,
				contentType: "text/xml",
				body: `<?xml version="1.0" encoding="utf-8"?>
					<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/" xmlns:ns="http://example.com/">
					<soap:Header/>
					<soap:Body>
						<ns:Test>
						</ns:Test>
					</soap:Body>
					</soap:Envelope>`,

				nextHandlerError: true,
			},
			&action{
				err:      utilhttp.ErrInternalServerError,
				respCode: 500,
				ehCode:   500,
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			meh := &mockErrorHandler{}
			m := &soapREST{
				eh:                    meh,
				attributeKey:          "@attribute",
				textKey:               "#text",
				namespaceKey:          "_namespace",
				arrayKey:              "item",
				separatorChar:         ":",
				extractStringElement:  true,
				extractBooleanElement: true,
				extractIntegerElement: true,
				extractFloatElement:   true,
			}

			var next http.Handler
			if tt.C().nextHandlerError {
				next = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					respJSON := `brokenJSON`
					w.Header().Set("Content-Type", "application/json; charset=utf-8")
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte(respJSON))
				})
			} else {
				next = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					respJSON := ``
					w.Header().Set("Content-Type", "application/json; charset=utf-8")
					w.WriteHeader(http.StatusOK)
					w.Write([]byte(respJSON))
				})
			}

			h := m.Middleware(next)
			req := httptest.NewRequest(tt.C().method, "http://test.com/test", strings.NewReader(tt.C().body))
			if tt.C().readBodyError {
				req.Body = io.NopCloser(&mockReader{})
			}
			req.Header.Set("Content-Type", tt.C().contentType)
			resp := httptest.NewRecorder()

			h.ServeHTTP(resp, req)
			testutil.Diff(t, tt.A().respCode, resp.Code)
			testutil.Diff(t, tt.A().ehCode, meh.code)
			testutil.Diff(t, strings.TrimSpace(tt.A().body), strings.TrimSpace(resp.Body.String()))

			opts := []cmp.Option{
				cmpopts.EquateErrors(),
			}
			testutil.DiffError(t, tt.A().err, nil, meh.err, opts...)
		})
	}
}
