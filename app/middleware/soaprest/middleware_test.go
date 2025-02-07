package soaprest

import (
	"bytes"
	"encoding/xml"
	"testing"

	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
)

// type mockErrorHandler struct {
// 	err  error
// 	code int
// }

// func (m *mockErrorHandler) ServeHTTPError(w http.ResponseWriter, r *http.Request, err error) {
// 	m.err = err
// 	m.code = err.(core.HTTPError).StatusCode()
// 	w.WriteHeader(m.code)
// }

// type mockReader struct{}

// func (m *mockReader) Read(p []byte) (int, error) {
// 	return 0, errors.New("mock read error")
// }

// type errorResponseRecorder struct {
// 	header     http.Header
// 	code       int
// 	writeError error
// }

// func (rec *errorResponseRecorder) Header() http.Header {
// 	if rec.header == nil {
// 		rec.header = make(http.Header)
// 	}
// 	return rec.header
// }

// func (rec *errorResponseRecorder) WriteHeader(code int) {
// 	rec.code = code
// }

// func (rec *errorResponseRecorder) Write(b []byte) (int, error) {
// 	rec.writeError = errors.New("mock write error")
// 	return 0, rec.writeError
// }

// func TestMiddleware(t *testing.T) {
// 	type condition struct {
// 		method      string
// 		contentType string
// 		body        string

// 		readBodyError    bool
// 		nextHandlerError bool
// 	}

// 	type action struct {
// 		body     string
// 		err      error
// 		respCode int
// 	}

// 	tb := testutil.NewTableBuilder[*condition, *action]()
// 	tb.Name(t.Name())
// 	table := tb.Build()

// 	gen := testutil.NewCase[*condition, *action]
// 	testCases := []*testutil.Case[*condition, *action]{
// 		gen(
// 			"NonSOAPRequest",
// 			[]string{},
// 			[]string{},
// 			&condition{
// 				method:      http.MethodGet,
// 				contentType: "application/json",
// 				body:        ``,
// 			},
// 			&action{
// 				err:      errInvalidSOAP11Request,
// 				respCode: 403,
// 			},
// 		),
// 		gen(
// 			"ReadBodyError",
// 			[]string{},
// 			[]string{},
// 			&condition{
// 				method:      http.MethodPost,
// 				contentType: "text/xml",
// 				body:        `<invalid_xml>`,

// 				readBodyError: true,
// 			},
// 			&action{
// 				err:      utilhttp.ErrBadRequest,
// 				respCode: 400,
// 			},
// 		),
// 		gen(
// 			"UnmarshalError",
// 			[]string{},
// 			[]string{},
// 			&condition{
// 				method:      http.MethodPost,
// 				contentType: "text/xml",
// 				body:        `<invalid_xml>`,

// 				readBodyError: false,
// 			},
// 			&action{
// 				err:      utilhttp.ErrBadRequest,
// 				respCode: 400,
// 			},
// 		),
// 		gen(
// 			"MarshalErrorNaN",
// 			[]string{},
// 			[]string{},
// 			&condition{
// 				method:      http.MethodPost,
// 				contentType: "text/xml",
// 				body: `<?xml version="1.0" encoding="UTF-8"?>
// 			<Envelope xmlns="http://schemas.xmlsoap.org/soap/envelope/">
// 			<Body>
// 			<Value>NaN</Value>
// 			</Body>
// 			</Envelope>`,
// 			},
// 			&action{
// 				err:      utilhttp.ErrBadRequest,
// 				respCode: 400,
// 			},
// 		),
// 		gen(
// 			"NextHandlerError",
// 			[]string{},
// 			[]string{},
// 			&condition{
// 				method:      http.MethodPost,
// 				contentType: "text/xml",
// 				body: `<?xml version="1.0" encoding="utf-8"?>
// 			<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/" xmlns:ns="http://example.com/">
// 			<soap:Header/>
// 			<soap:Body>
// 				<ns:Test>
// 				</ns:Test>
// 			</soap:Body>
// 			</soap:Envelope>`,

// 				nextHandlerError: true,
// 			},
// 			&action{
// 				err:      utilhttp.ErrInternalServerError,
// 				respCode: 500,
// 			},
// 		),
// 		gen(
// 			"ResponseWriterWriteError",
// 			[]string{},
// 			[]string{},
// 			&condition{
// 				method:      http.MethodPost,
// 				contentType: "text/xml",
// 				body: `<?xml version="1.0" encoding="utf-8"?>
// 		<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/" xmlns:ns="http://example.com/">
// 		<soap:Header/>
// 		<soap:Body>
// 			<ns:Test>
// 			</ns:Test>
// 		</soap:Body>
// 		</soap:Envelope>`,
// 			},
// 			&action{
// 				err:      utilhttp.ErrInternalServerError,
// 				respCode: 500,
// 			},
// 		),
// 	}

// 	testutil.Register(table, testCases...)

// 	for _, tt := range table.Entries() {
// 		tt := tt
// 		t.Run(tt.Name(), func(t *testing.T) {
// 			meh := &mockErrorHandler{}
// 			m := &soapREST{
// 				eh:                    meh,
// 				attributeKey:          "@attribute",
// 				textKey:               "#text",
// 				namespaceKey:          "_namespace",
// 				arrayKey:              "item",
// 				separatorChar:         ":",
// 				extractStringElement:  true,
// 				extractBooleanElement: true,
// 				extractIntegerElement: true,
// 				extractFloatElement:   true,
// 			}

// 			var next http.Handler
// 			if tt.C().nextHandlerError {
// 				next = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 					respJSON := `brokenJSON`
// 					w.Header().Set("Content-Type", "application/json; charset=utf-8")
// 					w.WriteHeader(http.StatusInternalServerError)
// 					w.Write([]byte(respJSON))
// 				})
// 			} else {
// 				next = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 					respJSON := `{"message": "success"}`
// 					w.Header().Set("Content-Type", "application/json; charset=utf-8")
// 					w.WriteHeader(http.StatusOK)
// 					w.Write([]byte(respJSON))
// 				})
// 			}

// 			h := m.Middleware(next)
// 			req := httptest.NewRequest(tt.C().method, "http://test.com/test", strings.NewReader(tt.C().body))
// 			if tt.C().readBodyError {
// 				req.Body = io.NopCloser(&mockReader{})
// 			}
// 			req.Header.Set("Content-Type", tt.C().contentType)

// 			var resp http.ResponseWriter
// 			if tt.Name() == "ResponseWriterWriteError" {
// 				resp = &errorResponseRecorder{}
// 			} else {
// 				resp = httptest.NewRecorder()
// 			}

// 			h.ServeHTTP(resp, req)

// 			if rec, ok := resp.(*httptest.ResponseRecorder); ok {
// 				testutil.Diff(t, tt.A().respCode, rec.Code)
// 				testutil.Diff(t, strings.TrimSpace(tt.A().body), strings.TrimSpace(rec.Body.String()))
// 			} else if rec, ok := resp.(*errorResponseRecorder); ok {
// 				testutil.Diff(t, tt.A().respCode, rec.code)
// 			}

// 			opts := []cmp.Option{
// 				cmpopts.EquateErrors(),
// 			}
// 			testutil.DiffError(t, tt.A().err, nil, meh.err, opts...)
// 		})
// 	}
// }

// func TestXmlToMap(t *testing.T) {
// 	type condition struct {
// 		xmlInput xmlNode
// 	}

// 	type action struct {
// 		expected any
// 	}

// 	tb := testutil.NewTableBuilder[*condition, *action]()
// 	tb.Name(t.Name())
// 	table := tb.Build()
// 	gen := testutil.NewCase[*condition, *action]

// 	testCases := []*testutil.Case[*condition, *action]{
// 		gen(
// 			"SimpleTextNode",
// 			[]string{},
// 			[]string{},
// 			&condition{
// 				xmlInput: xmlNode{
// 					XMLName: xml.Name{Local: "Test"},
// 					Content: "Test Content",
// 				},
// 			},
// 			&action{
// 				expected: "Test Content",
// 			},
// 		),
// 		gen(
// 			"MapChildNodeWithMapTypeChildren",
// 			[]string{},
// 			[]string{},
// 			&condition{
// 				xmlInput: xmlNode{
// 					XMLName: xml.Name{Local: "Test"},
// 					Children: []xmlNode{
// 						{
// 							XMLName: xml.Name{
// 								Local: "ChildLocal",
// 							},
// 						},
// 					},
// 				},
// 			},
// 			&action{
// 				expected: map[string]any{
// 					"Test": map[string]any{
// 						"ChildLocal": map[string]any{},
// 					},
// 				},
// 			},
// 		),
// 		gen(
// 			"NodeWithAttributes",
// 			[]string{},
// 			[]string{},
// 			&condition{
// 				xmlInput: xmlNode{
// 					XMLName: xml.Name{Local: "Person"},
// 					Attrs: []xml.Attr{
// 						{Name: xml.Name{Local: "id"}, Value: "123"},
// 						{Name: xml.Name{Local: "role"}, Value: "admin"},
// 					},
// 					Children: []xmlNode{
// 						{
// 							XMLName: xml.Name{Local: "Name"},
// 							Content: "John Doe",
// 						},
// 						{
// 							XMLName: xml.Name{Local: "Email"},
// 							Content: "john@example.com",
// 						},
// 					},
// 				},
// 			},
// 			&action{
// 				expected: map[string]any{
// 					"Person": map[string]any{
// 						"@attribute": map[string]string{
// 							"id":   "123",
// 							"role": "admin",
// 						},
// 						"Name":  "John Doe",
// 						"Email": "john@example.com",
// 					},
// 				},
// 			},
// 		),
// 		gen(
// 			"NodeWithNamespaces",
// 			[]string{},
// 			[]string{},
// 			&condition{
// 				xmlInput: xmlNode{
// 					XMLName: xml.Name{Space: "http://example.com/ns", Local: "Test"},
// 					Attrs: []xml.Attr{
// 						{Name: xml.Name{Space: "xmlns", Local: "ns"}, Value: "http://example.com/ns"},
// 					},
// 					Children: []xmlNode{
// 						{
// 							XMLName: xml.Name{Local: "Value"},
// 							Content: "42",
// 						},
// 					},
// 				},
// 			},
// 			&action{
// 				expected: map[string]any{
// 					"ns:Test": map[string]any{
// 						"_namespace": map[string]string{
// 							"ns": "http://example.com/ns",
// 						},
// 						"Value": int64(42),
// 					},
// 				},
// 			},
// 		),
// 		gen(
// 			"NodeWithXsiNilAttribute",
// 			[]string{},
// 			[]string{},
// 			&condition{
// 				xmlInput: xmlNode{
// 					XMLName: xml.Name{Local: "OptionalValue"},
// 					Attrs: []xml.Attr{
// 						{
// 							Name:  xml.Name{Space: "xsi", Local: "nil"},
// 							Value: "true",
// 						},
// 					},
// 				},
// 			},
// 			&action{
// 				expected: nil,
// 			},
// 		),
// 		gen(
// 			"NestedNodes",
// 			[]string{},
// 			[]string{},
// 			&condition{
// 				xmlInput: xmlNode{
// 					XMLName: xml.Name{Local: "Order"},
// 					Children: []xmlNode{
// 						{
// 							XMLName: xml.Name{Local: "Items"},
// 							Children: []xmlNode{
// 								{
// 									XMLName: xml.Name{Local: "Item"},
// 									Children: []xmlNode{
// 										{
// 											XMLName: xml.Name{Local: "Name"},
// 											Content: "Item1",
// 										},
// 										{
// 											XMLName: xml.Name{Local: "Quantity"},
// 											Content: "2",
// 										},
// 									},
// 								},
// 								{
// 									XMLName: xml.Name{Local: "Item"},
// 									Children: []xmlNode{
// 										{
// 											XMLName: xml.Name{Local: "Name"},
// 											Content: "Item2",
// 										},
// 										{
// 											XMLName: xml.Name{Local: "Quantity"},
// 											Content: "3",
// 										},
// 									},
// 								},
// 							},
// 						},
// 						{
// 							XMLName: xml.Name{Local: "Total"},
// 							Content: "100.50",
// 						},
// 					},
// 				},
// 			},
// 			&action{
// 				expected: map[string]any{
// 					"Order": map[string]any{
// 						"Items": map[string]any{
// 							"Item": []any{
// 								map[string]any{
// 									"Name":     "Item1",
// 									"Quantity": int64(2),
// 								},
// 								map[string]any{
// 									"Name":     "Item2",
// 									"Quantity": int64(3),
// 								},
// 							},
// 						},
// 						"Total": 100.50,
// 					},
// 				},
// 			},
// 		),
// 		gen(
// 			"SOAPEnvelopeWithTextContentAndAttributes",
// 			[]string{},
// 			[]string{},
// 			&condition{
// 				xmlInput: xmlNode{
// 					XMLName: xml.Name{
// 						Space: "http://schemas.xmlsoap.org/soap/envelope/",
// 						Local: "Envelope",
// 					},
// 					Attrs: []xml.Attr{
// 						{
// 							Name:  xml.Name{Local: "xmlns:soap"},
// 							Value: "http://schemas.xmlsoap.org/soap/envelope/",
// 						},
// 					},
// 					Content: "Text Content",
// 					Children: []xmlNode{
// 						{
// 							XMLName: xml.Name{Local: "Body"},
// 							Children: []xmlNode{
// 								{
// 									XMLName: xml.Name{Local: "someData"},
// 									Content: "Some Value",
// 								},
// 							},
// 						},
// 					},
// 				},
// 			},
// 			&action{
// 				expected: map[string]any{
// 					"soap:Envelope": map[string]any{
// 						"@attribute": map[string]string{
// 							"xmlns:soap": "http://schemas.xmlsoap.org/soap/envelope/",
// 						},
// 						"#text": "Text Content",
// 						"Body": map[string]any{
// 							"someData": "Some Value",
// 						},
// 					},
// 				},
// 			},
// 		),
// 	}

// 	testutil.Register(table, testCases...)

// 	for _, tt := range table.Entries() {
// 		tt := tt
// 		t.Run(tt.Name(), func(t *testing.T) {
// 			s := soapREST{
// 				attributeKey:          "@attribute",
// 				textKey:               "#text",
// 				namespaceKey:          "_namespace",
// 				arrayKey:              "item",
// 				separatorChar:         ":",
// 				extractStringElement:  false,
// 				extractBooleanElement: true,
// 				extractIntegerElement: true,
// 				extractFloatElement:   true,
// 			}

// 			nsCtx := &namespaceContext{
// 				prefixToURI: map[string]string{},
// 				uriToPrefix: map[string]string{},
// 			}

// 			result := s.xmlToMap(tt.C().xmlInput, nsCtx)
// 			testutil.Diff(t, tt.A().expected, result)
// 		})
// 	}
// }

// func TestConvertRESTtoSOAPResponse(t *testing.T) {
// 	type condition struct {
// 		restData []byte
// 	}

// 	type action struct {
// 		xml        []byte
// 		err        error
// 		errPattern *regexp.Regexp
// 	}

// 	tb := testutil.NewTableBuilder[*condition, *action]()
// 	tb.Name(t.Name())
// 	table := tb.Build()
// 	gen := testutil.NewCase[*condition, *action]

// 	testCases := []*testutil.Case[*condition, *action]{
// 		// <TODO> Consider whether it is possible to compare XML while ignoring whitespace.
// 		gen(
// 			"ValidSOAPResponse",
// 			[]string{},
// 			[]string{},
// 			&condition{
// 				restData: []byte(`{"soap:Envelope": {"soap:Body": {"Response": {"Result": "Success"}}}}`),
// 			},
// 			&action{
// 				xml: []byte(`<?xml version="1.0" encoding="UTF-8"?>
// <soap:Envelope>
//   <soap:Header></soap:Header>
//   <soap:Body>
//     <Response>
//       <Result>Success</Result>
//     </Response>
//   </soap:Body>
// </soap:Envelope>`),
// 				err:        nil,
// 				errPattern: nil,
// 			},
// 		),
// 		gen(
// 			"DecodeError",
// 			[]string{},
// 			[]string{},
// 			&condition{
// 				restData: nil,
// 			},
// 			&action{
// 				xml:        nil,
// 				err:        nil,
// 				errPattern: regexp.MustCompile(`xml: encoding error: invalid UTF-8`),
// 			},
// 		),
// 	}

// 	testutil.Register(table, testCases...)

// 	for _, tt := range table.Entries() {
// 		tt := tt
// 		t.Run(tt.Name(), func(t *testing.T) {
// 			s := &soapREST{
// 				attributeKey:  "@attribute",
// 				textKey:       "#text",
// 				namespaceKey:  "_namespace",
// 				arrayKey:      "item",
// 				separatorChar: ":",
// 			}

// 			wrapper := &wrappedWriter{
// 				body: bytes.NewBuffer(tt.C().restData),
// 			}
// 			result, err := s.convertRESTtoSOAPResponse(wrapper)

// 			opts := []cmp.Option{
// 				cmpopts.IgnoreFields(soapEnvelope{}, "XMLName"),
// 				cmpopts.IgnoreFields(soapBody{}, "XMLName"),
// 				cmpopts.IgnoreFields(soapHeader{}, "XMLName"),
// 			}

// 			if tt.A().err != nil {
// 				testutil.DiffError(t, tt.A().err, tt.A().errPattern, err)
// 			} else {
// 				testutil.Diff(t, tt.A().xml, result, opts...)
// 			}
// 		})
// 	}
// }

func TestXmlElement_MarshalXML(t *testing.T) {
	type condition struct {
		element xmlElement
	}

	type action struct {
		xmlOutput string
		err       error
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()
	gen := testutil.NewCase[*condition, *action]

	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"SimpleElement",
			nil,
			nil,
			&condition{
				element: xmlElement{
					XMLName: xml.Name{Local: "Test"},
					Content: "TestContent",
				},
			},
			&action{
				xmlOutput: `<Test>TestContent</Test>`,
				err:       nil,
			},
		),
		gen(
			"NotEmptySpace",
			nil,
			nil,
			&condition{
				element: xmlElement{
					XMLName: xml.Name{Local: "Test", Space: "ns"},
					Content: "TestContent",
				},
			},
			&action{
				xmlOutput: `<ns:Test xmlns="ns">TestContent</ns:Test>`,
				err:       nil,
			},
		),
		gen(
			"NilElement",
			nil,
			nil,
			&condition{
				element: xmlElement{
					XMLName: xml.Name{Local: "Test"},
					Content: "NilContent",
					isNil:   true,
				},
			},
			&action{
				xmlOutput: `<Test xsi:nil="true">NilContent</Test>`,
				err:       nil,
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			var buf bytes.Buffer
			var err error

			enc := xml.NewEncoder(&buf)
			err = tt.C().element.MarshalXML(enc, xml.StartElement{Name: tt.C().element.XMLName})
			if tt.A().err != nil {
				testutil.DiffError(t, tt.A().err, nil, err)
				return
			}

			err = enc.Flush()
			testutil.DiffError(t, tt.A().err, nil, err)

			testutil.Diff(t, string([]byte(tt.A().xmlOutput)), string(bytes.TrimSpace(buf.Bytes())))
		})
	}
}
