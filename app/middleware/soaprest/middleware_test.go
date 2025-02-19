package soaprest

import (
	"bytes"
	"cmp"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/http/httptest"
	"regexp"
	"sort"
	"strings"
	"testing"

	"github.com/aileron-gateway/aileron-gateway/app"
	"github.com/aileron-gateway/aileron-gateway/core"
	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
	utilhttp "github.com/aileron-gateway/aileron-gateway/util/http"
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

type testMatcher struct {
	match bool
}

func (t testMatcher) Match(s string) bool {
	return t.match
}

func TestSOAPREST_Middleware_RequestConversion(t *testing.T) {
	type condition struct {
		body        string
		method      string
		contentType string

		paths *testMatcher

		attributeKey  string
		textKey       string
		namespaceKey  string
		arrayKey      string
		separatorChar string

		extractStringElement  bool
		extractBooleanElement bool
		extractIntegerElement bool
		extractFloatElement   bool

		readBodyError bool
	}

	type action struct {
		body       string
		err        any // error or errorKind
		errPattern *regexp.Regexp
		code       int
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"GetSOAPRequest",
			nil,
			nil,
			&condition{
				method:      http.MethodGet,
				contentType: "text/xml",
				body: `<?xml version="1.0" encoding="utf-8"?>
						<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/" xmlns:ns="http://example.com/">
							<soap:Header>
								<ns:Auth>
									<ns:Username>TestUser</ns:Username>
									<ns:Password>password</ns:Password>
								</ns:Auth>
								"double quoted text"
							</soap:Header>
							<soap:Body>
								<ns:Test testAttributeKey="someValue">
									<ns:Value>123</ns:Value>
								</ns:Test>
								<ns:Array>
									<item>100</item>
									<item>3.14</item>
									<item>true</item>
									<item>someText</item>
								</ns:Array>
							</soap:Body>
						</soap:Envelope>`,

				paths: &testMatcher{match: true},

				extractStringElement:  true,
				extractBooleanElement: true,
				extractIntegerElement: true,
				extractFloatElement:   true,
			},
			&action{
				body: `{"soap:Envelope": {
							"_namespace": {
								"ns": "http://example.com/",
								"soap": "http://schemas.xmlsoap.org/soap/envelope/"
							},
							"soap:Body": {
								"ns:Test": {
									"@attribute": {
										"testAttributeKey": "someValue"
									},
									"ns:Value": 123
								},
								"ns:Array": {
									"item": [100, 3.14, true, "someText"]
								}
							},
							"soap:Header": {
								"#text": "double quoted text",
								"ns:Auth": {
									"ns:Username": "TestUser",
									"ns:Password": "password"
								}
							}
						}}`,

				// This is a case where no errors occur on the request side, and no body is set for the response.
				// As a result, a decode error occurs on the response side.
				err:  app.ErrAppMiddleSOAPRESTDecodeResponseBody,
				code: 500,
			},
		),
		gen(
			"PostSOAPRequest",
			nil,
			nil,
			&condition{
				method:      http.MethodPost,
				contentType: "text/xml",
				body: `<?xml version="1.0" encoding="utf-8"?>
						<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/" xmlns:ns="http://example.com/">
							<soap:Header>
								<ns:Auth>
									<ns:Username>TestUser</ns:Username>
									<ns:Password>password</ns:Password>
								</ns:Auth>
								"double quoted text"
							</soap:Header>
							<soap:Body>
								<ns:Test testAttributeKey="someValue">
									<ns:Value>123</ns:Value>
								</ns:Test>
								<ns:Array>
									<item>100</item>
									<item>3.14</item>
									<item>true</item>
									<item>someText</item>
								</ns:Array>
							</soap:Body>
						</soap:Envelope>`,

				paths: &testMatcher{match: true},

				extractStringElement:  true,
				extractBooleanElement: true,
				extractIntegerElement: true,
				extractFloatElement:   true,
			},
			&action{
				body: `{"soap:Envelope": {
							"_namespace": {
								"ns": "http://example.com/",
								"soap": "http://schemas.xmlsoap.org/soap/envelope/"
							},
							"soap:Body": {
								"ns:Test": {
									"@attribute": {
										"testAttributeKey": "someValue"
									},
									"ns:Value": 123
								},
								"ns:Array": {
									"item": [100, 3.14, true, "someText"]
								}
							},
							"soap:Header": {
								"#text": "double quoted text",
								"ns:Auth": {
									"ns:Username": "TestUser",
									"ns:Password": "password"
								}
							}
						}}`,
				// This is a case where no errors occur on the request side, and no body is set for the response.
				// As a result, a decode error occurs on the response side.
				err:  app.ErrAppMiddleSOAPRESTDecodeResponseBody,
				code: 500,
			},
		),
		gen(
			"NonSOAPRequest",
			nil,
			nil,
			&condition{
				method:      http.MethodPost,
				contentType: "application/json",
				body:        ``,
				paths:       &testMatcher{match: true},
			},
			&action{
				err:  app.ErrAppMiddleSOAPRESTVersionMismatch,
				code: 403,
			},
		),
		gen(
			"ReadBodyError",
			nil,
			nil,
			&condition{
				method:      http.MethodPost,
				contentType: "text/xml",
				body:        `<invalid_xml>`,
				paths:       &testMatcher{match: true},

				readBodyError: true,
			},
			&action{
				err:  app.ErrAppMiddleSOAPRESTReadRequestBody,
				code: 400,
			},
		),
		gen(
			"UnmarshalError",
			nil,
			nil,
			&condition{
				method:      http.MethodPost,
				contentType: "text/xml",
				body:        `<invalid_xml>`,
				paths:       &testMatcher{match: true},

				readBodyError: false,
			},
			&action{
				err:  app.ErrAppGenUnmarshal,
				code: 400,
			},
		),
		gen(
			"MarshalErrorNaN",
			nil,
			nil,
			&condition{
				method:      http.MethodPost,
				contentType: "text/xml",
				body: `<?xml version="1.0" encoding="UTF-8"?>
						<Envelope xmlns="http://schemas.xmlsoap.org/soap/envelope/">
							<Body>
								<Value>NaN</Value>
							</Body>
						</Envelope>`,
				paths: &testMatcher{match: true},

				extractFloatElement: true,
			},
			&action{
				err:  app.ErrAppMiddleSOAPRESTMarshalJSONData,
				code: 400,
			},
		),
		gen(
			"PostSOAPRequestWithModifiedKeys",
			nil,
			nil,
			&condition{
				method:      http.MethodPost,
				contentType: "text/xml",
				body: `<?xml version="1.0" encoding="utf-8"?>
						<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/" xmlns:ns="http://example.com/">
							<soap:Header>
								<ns:Auth>
									<ns:Username>TestUser</ns:Username>
									<ns:Password>password</ns:Password>
								</ns:Auth>
								"double quoted text"
							</soap:Header>
							<soap:Body>
								<ns:Test testAttributeKey="someValue">
									<ns:Value>123</ns:Value>
								</ns:Test>
								<ns:Array>
									<item>100</item>
									<item>3.14</item>
									<item>true</item>
									<item>someText</item>
								</ns:Array>
							</soap:Body>
						</soap:Envelope>`,

				paths: &testMatcher{match: true},

				attributeKey:  "@attr!",
				textKey:       "#textKey",
				namespaceKey:  "_ns",
				arrayKey:      "element",
				separatorChar: "*",

				extractStringElement:  true,
				extractBooleanElement: true,
				extractIntegerElement: true,
				extractFloatElement:   true,
			},
			&action{
				body: `{"soap*Envelope": {
							"_ns": {
								"ns": "http://example.com/",
								"soap": "http://schemas.xmlsoap.org/soap/envelope/"
							},
							"soap*Body": {
								"ns*Test": {
									"@attr!": {
										"testAttributeKey": "someValue"
									},
									"ns*Value": 123
								},
								"ns*Array": {
									"item": [100, 3.14, true, "someText"]
								}
							},
							"soap*Header": {
								"#textKey": "double quoted text",
								"ns*Auth": {
									"ns*Username": "TestUser",
									"ns*Password": "password"
								}
							}
						}}`,

				// This is a case where no errors occur on the request side, and no body is set for the response.
				// As a result, a decode error occurs on the response side.
				err:  app.ErrAppMiddleSOAPRESTDecodeResponseBody,
				code: 500,
			},
		),
		gen(
			"PostSOAPRequestWithModifiedKeys",
			nil,
			nil,
			&condition{
				method:      http.MethodPost,
				contentType: "text/xml",
				body: `<?xml version="1.0" encoding="utf-8"?>
						<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/" xmlns:ns="http://example.com/">
							<soap:Header>
								<ns:Auth>
									<ns:Username>TestUser</ns:Username>
									<ns:Password>password</ns:Password>
								</ns:Auth>
								"double quoted text"
							</soap:Header>
							<soap:Body>
								<ns:Test testAttributeKey="someValue">
									<ns:Value>123</ns:Value>
								</ns:Test>
								<ns:Array>
									<item>100</item>
									<item>3.14</item>
									<item>true</item>
									<item>someText</item>
								</ns:Array>
							</soap:Body>
						</soap:Envelope>`,

				paths: &testMatcher{match: true},

				extractStringElement:  false,
				extractBooleanElement: false,
				extractIntegerElement: false,
				extractFloatElement:   false,
			},
			&action{
				body: `{"soap:Envelope": {
							"_namespace": {
								"ns": "http://example.com/",
								"soap": "http://schemas.xmlsoap.org/soap/envelope/"
							},
							"soap:Body": {
								"ns:Test": {
									"@attribute": {
										"testAttributeKey": "someValue"
									},
									"ns:Value": "123"
								},
								"ns:Array": {
									"item": ["100", "3.14", "true", "someText"]
								}
							},
							"soap:Header": {
								"#text": "\"double quoted text\"",
								"ns:Auth": {
									"ns:Username": "TestUser",
									"ns:Password": "password"
								}
							}
						}}`,

				// This is a case where no errors occur on the request side, and no body is set for the response.
				// As a result, a decode error occurs on the response side.
				err:  app.ErrAppMiddleSOAPRESTDecodeResponseBody,
				code: 500,
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			meh := &mockErrorHandler{}
			sr := &soapREST{
				eh:                    meh,
				paths:                 tt.C().paths,
				attributeKey:          cmp.Or(tt.C().attributeKey, "@attribute"),
				textKey:               cmp.Or(tt.C().textKey, "#text"),
				namespaceKey:          cmp.Or(tt.C().namespaceKey, "_namespace"),
				arrayKey:              cmp.Or(tt.C().arrayKey, "item"),
				separatorChar:         cmp.Or(tt.C().separatorChar, ":"),
				extractStringElement:  tt.C().extractStringElement,
				extractBooleanElement: tt.C().extractBooleanElement,
				extractIntegerElement: tt.C().extractIntegerElement,
				extractFloatElement:   tt.C().extractFloatElement,
			}

			nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				bodyBytes, err := io.ReadAll(r.Body)
				if err != nil {
					t.Fatalf("Failed to read request body: %v", err)
				}
				r.Body.Close()

				var actualJSON, expectedJSON any
				if err := json.Unmarshal([]byte(bodyBytes), &actualJSON); err != nil {
					t.Fatalf("Failed to unmarshal transformed JSON: %v", err)
				}

				if err := json.Unmarshal([]byte(tt.A().body), &expectedJSON); err != nil {
					t.Fatalf("Failed to unmarshal expected JSON: %v", err)
				}

				// Check whether the SOAP/XML request is being converted
				// into a REST/JSON request by the SOAPRESTMiddleware.
				testutil.Diff(t, expectedJSON, actualJSON)
				w.WriteHeader(http.StatusOK)
			})

			h := sr.Middleware(nextHandler)
			req := httptest.NewRequest(tt.C().method, "http://test.com/test", strings.NewReader(tt.C().body))
			if tt.C().readBodyError {
				req.Body = io.NopCloser(&mockReader{})
			}
			req.Header.Set("Content-Type", tt.C().contentType)

			resp := httptest.NewRecorder()
			h.ServeHTTP(resp, req)

			opts := []gocmp.Option{
				cmpopts.EquateErrors(),
			}
			testutil.DiffError(t, tt.A().err, tt.A().errPattern, meh.err, opts...)
		})
	}
}

// testNode is a struct that represents an XML node in middleware test
type testNode struct {
	Name       xml.Name
	Attr       []xml.Attr
	Text       string
	ChildNodes []*testNode
}

// parseXML parses an XML string and builds a Node tree.
func parseXML(r io.Reader) (*testNode, error) {
	decoder := xml.NewDecoder(r)
	var root *testNode
	var stack []*testNode

	for {
		tok, err := decoder.Token()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		switch tok := tok.(type) {
		case xml.StartElement:
			node := &testNode{
				Name: tok.Name,
				Attr: tok.Attr,
			}
			sort.Slice(node.Attr, func(i, j int) bool {
				return node.Attr[i].Name.Local < node.Attr[j].Name.Local
			})
			if len(stack) == 0 {
				root = node
			} else {
				parent := stack[len(stack)-1]
				parent.ChildNodes = append(parent.ChildNodes, node)
			}
			stack = append(stack, node)

		case xml.EndElement:
			if len(stack) > 0 {
				stack = stack[:len(stack)-1]
			}

		case xml.CharData:
			if len(stack) > 0 {
				text := strings.TrimSpace(string(tok))
				if text != "" {
					currentNode := stack[len(stack)-1]
					currentNode.Text += text
				}
			}
		}
	}

	return root, nil
}

func compareNodes(a, b *testNode) bool {
	// If the result of parseXML is nil, it is considered a match.
	if a == nil && b == nil {
		return true
	}

	if a.Name.Local != b.Name.Local || a.Name.Space != b.Name.Space {
		return false
	}

	if len(a.Attr) != len(b.Attr) {
		return false
	}
	for i := range a.Attr {
		if a.Attr[i].Name.Local != b.Attr[i].Name.Local ||
			a.Attr[i].Name.Space != b.Attr[i].Name.Space ||
			a.Attr[i].Value != b.Attr[i].Value {
			return false
		}
	}

	if a.Text != b.Text {
		return false
	}

	if len(a.ChildNodes) != len(b.ChildNodes) {
		return false
	}
	for i := range a.ChildNodes {
		if !compareNodes(a.ChildNodes[i], b.ChildNodes[i]) {
			return false
		}
	}

	return true
}

func TestSOAPREST_Middleware_ResponseConversion(t *testing.T) {
	type condition struct {
		body        string
		method      string
		contentType string

		responseWriteError bool
	}

	type action struct {
		body       string
		err        any // error or errorKind
		errPattern *regexp.Regexp
		code       int
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"SimpleSOAPResponse",
			nil,
			nil,
			&condition{
				method:      http.MethodPost,
				contentType: "application/json",
				body: `{
                    "soap:Envelope": {
                        "_namespace": {
                            "ns": "http://example.com/",
                            "soap": "http://schemas.xmlsoap.org/soap/envelope/"
                        },
                        "soap:Body": {
                            "ns:Test": {
                                "@attribute": {
                                    "testAttributeKey": "someValue"
                                },
                                "ns:Value": 123
                            },
                            "ns:Array": {
                                "item": [100, 3.14, true, "someText"]
                            }
                        },
                        "soap:Header": {
                            "#text": "double quoted text",
                            "ns:Auth": {
                                "ns:Username": "TestUser",
                                "ns:Password": "password"
                            }
                        }
                    }
                }`,
			},
			&action{
				body: `<?xml version="1.0" encoding="utf-8"?>
                        <soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/" xmlns:ns="http://example.com/">
                            <soap:Header>
                                double quoted text
                                <ns:Auth>
                                    <ns:Username>TestUser</ns:Username>
                                    <ns:Password>password</ns:Password>
                                </ns:Auth>
                            </soap:Header>
                            <soap:Body>
                                <ns:Test testAttributeKey="someValue">
                                    <ns:Value>123</ns:Value>
                                </ns:Test>
                                <ns:Array>
                                    <item>100</item>
                                    <item>3.14</item>
                                    <item>true</item>
                                    <item>someText</item>
                                </ns:Array>
                            </soap:Body>
                        </soap:Envelope>`,
				err:  nil,
				code: 0,
			},
		),
		gen(
			"NextHandlerError",
			nil,
			nil,
			&condition{
				method:      http.MethodPost,
				contentType: "text/xml",
				body:        `{brokenJSON}`,
			},
			&action{
				body:       ``,
				err:        app.ErrAppMiddleSOAPRESTDecodeResponseBody,
				errPattern: regexp.MustCompile("failed to decode:"),
				code:       500,
			},
		),
		gen(
			"ResponseWriterWriteError",
			nil,
			nil,
			&condition{
				responseWriteError: true,

				method:      http.MethodPost,
				contentType: "text/xml",
				body:        `{"soap:Envelope":{"_namespace":{"ns":"http://example.com/","soap":"http://schemas.xmlsoap.org/soap/envelope/"},"soap:Body":{"ns:Test":{}},"soap:Header":{}}}`,
			},
			&action{
				err:  app.ErrAppMiddleSOAPRESTWriteResponseBody,
				code: 500,
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

			var actualXML []byte

			nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", tt.C().contentType+"; charset=utf-8")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(tt.C().body))
			})

			h := m.Middleware(nextHandler)

			req := httptest.NewRequest(http.MethodPost, "http://example.com",
				strings.NewReader(`<?xml version="1.0" encoding="utf-8"?> <soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/"> </soap:Envelope>`))
			req.Header.Set("Content-Type", "text/xml; charset=utf-8")

			var resp http.ResponseWriter
			if tt.C().responseWriteError {
				resp = &errorResponseRecorder{}
			} else {
				resp = httptest.NewRecorder()
			}

			h.ServeHTTP(resp, req)

			if rec, ok := resp.(*httptest.ResponseRecorder); ok {
				actualXML = rec.Body.Bytes()
				fmt.Printf("Generated XML in Test: \n%s\n", string(actualXML))

				expectedNode, err := parseXML(strings.NewReader(tt.A().body))
				if err != nil {
					t.Fatalf("Failed to parse expected XML: %v", err)
				}

				actualNode, err := parseXML(bytes.NewReader(actualXML))
				if err != nil {
					t.Fatalf("Failed to parse actual XML: %v", err)
				}
				testutil.Diff(t, true, compareNodes(expectedNode, actualNode))
			} else if rec, ok := resp.(*errorResponseRecorder); ok {
				testutil.Diff(t, tt.A().code, rec.code)
			}

			opts := []gocmp.Option{
				cmpopts.EquateErrors(),
			}
			testutil.DiffError(t, tt.A().err, tt.A().errPattern, meh.err, opts...)
		})
	}
}

func TestSOAPREST_XmlToMap(t *testing.T) {
	type condition struct {
		xmlInput xmlNode
	}

	type action struct {
		expected any
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()
	gen := testutil.NewCase[*condition, *action]

	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"SimpleTextNode",
			nil,
			nil,
			&condition{
				xmlInput: xmlNode{
					XMLName: xml.Name{Local: "Test"},
					Content: "Test Content",
				},
			},
			&action{
				expected: "Test Content",
			},
		),
		gen(
			"MapChildNodeWithMapTypeChildren",
			nil,
			nil,
			&condition{
				xmlInput: xmlNode{
					XMLName: xml.Name{Local: "Test"},
					Children: []xmlNode{
						{
							XMLName: xml.Name{
								Local: "ChildLocal",
							},
						},
					},
				},
			},
			&action{
				expected: map[string]any{
					"Test": map[string]any{
						"ChildLocal": map[string]any{},
					},
				},
			},
		),
		gen(
			"NodeWithAttributes",
			nil,
			nil,
			&condition{
				xmlInput: xmlNode{
					XMLName: xml.Name{Local: "Person"},
					Attrs: []xml.Attr{
						{Name: xml.Name{Local: "id"}, Value: "123"},
						{Name: xml.Name{Local: "role"}, Value: "admin"},
					},
					Children: []xmlNode{
						{
							XMLName: xml.Name{Local: "Name"},
							Content: "John Doe",
						},
						{
							XMLName: xml.Name{Local: "Email"},
							Content: "john@example.com",
						},
					},
				},
			},
			&action{
				expected: map[string]any{
					"Person": map[string]any{
						"@attribute": map[string]string{
							"id":   "123",
							"role": "admin",
						},
						"Name":  "John Doe",
						"Email": "john@example.com",
					},
				},
			},
		),
		gen(
			"NodeWithNamespaces",
			nil,
			nil,
			&condition{
				xmlInput: xmlNode{
					XMLName: xml.Name{Space: "http://example.com/ns", Local: "Test"},
					Attrs: []xml.Attr{
						{Name: xml.Name{Space: "xmlns", Local: "ns"}, Value: "http://example.com/ns"},
					},
					Children: []xmlNode{
						{
							XMLName: xml.Name{Local: "Value"},
							Content: "42",
						},
					},
				},
			},
			&action{
				expected: map[string]any{
					"ns:Test": map[string]any{
						"_namespace": map[string]string{
							"ns": "http://example.com/ns",
						},
						"Value": int64(42),
					},
				},
			},
		),
		gen(
			"NodeWithXsiNilAttribute",
			nil,
			nil,
			&condition{
				xmlInput: xmlNode{
					XMLName: xml.Name{Local: "OptionalValue"},
					Attrs: []xml.Attr{
						{
							Name:  xml.Name{Space: "xsi", Local: "nil"},
							Value: "true",
						},
					},
				},
			},
			&action{
				expected: nil,
			},
		),
		gen(
			"NestedNodes",
			nil,
			nil,
			&condition{
				xmlInput: xmlNode{
					XMLName: xml.Name{Local: "Order"},
					Children: []xmlNode{
						{
							XMLName: xml.Name{Local: "Items"},
							Children: []xmlNode{
								{
									XMLName: xml.Name{Local: "Item"},
									Children: []xmlNode{
										{
											XMLName: xml.Name{Local: "Name"},
											Content: "Item1",
										},
										{
											XMLName: xml.Name{Local: "Quantity"},
											Content: "2",
										},
									},
								},
								{
									XMLName: xml.Name{Local: "Item"},
									Children: []xmlNode{
										{
											XMLName: xml.Name{Local: "Name"},
											Content: "Item2",
										},
										{
											XMLName: xml.Name{Local: "Quantity"},
											Content: "3",
										},
									},
								},
							},
						},
						{
							XMLName: xml.Name{Local: "Total"},
							Content: "100.50",
						},
					},
				},
			},
			&action{
				expected: map[string]any{
					"Order": map[string]any{
						"Items": map[string]any{
							"Item": []any{
								map[string]any{
									"Name":     "Item1",
									"Quantity": int64(2),
								},
								map[string]any{
									"Name":     "Item2",
									"Quantity": int64(3),
								},
							},
						},
						"Total": 100.50,
					},
				},
			},
		),
		gen(
			"SOAPEnvelopeWithTextContentAndAttributes",
			nil,
			nil,
			&condition{
				xmlInput: xmlNode{
					XMLName: xml.Name{
						Space: "http://schemas.xmlsoap.org/soap/envelope/",
						Local: "Envelope",
					},
					Attrs: []xml.Attr{
						{
							Name:  xml.Name{Local: "xmlns:soap"},
							Value: "http://schemas.xmlsoap.org/soap/envelope/",
						},
					},
					Content: "Text Content",
					Children: []xmlNode{
						{
							XMLName: xml.Name{Local: "Body"},
							Children: []xmlNode{
								{
									XMLName: xml.Name{Local: "someData"},
									Content: "Some Value",
								},
							},
						},
					},
				},
			},
			&action{
				expected: map[string]any{
					"soap:Envelope": map[string]any{
						"@attribute": map[string]string{
							"xmlns:soap": "http://schemas.xmlsoap.org/soap/envelope/",
						},
						"#text": "Text Content",
						"Body": map[string]any{
							"someData": "Some Value",
						},
					},
				},
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			s := soapREST{
				attributeKey:          "@attribute",
				textKey:               "#text",
				namespaceKey:          "_namespace",
				arrayKey:              "item",
				separatorChar:         ":",
				extractStringElement:  false,
				extractBooleanElement: true,
				extractIntegerElement: true,
				extractFloatElement:   true,
			}

			nsCtx := &namespaceContext{
				prefixToURI: map[string]string{},
				uriToPrefix: map[string]string{},
			}

			result := s.xmlToMap(tt.C().xmlInput, nsCtx)
			testutil.Diff(t, tt.A().expected, result)
		})
	}
}

func TestSOAPREST_ConvertRESTtoSOAPResponse(t *testing.T) {
	type condition struct {
		restData []byte
	}

	type action struct {
		xml        []byte
		err        any // error or errorutil.Kind
		errPattern *regexp.Regexp
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()
	gen := testutil.NewCase[*condition, *action]

	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"ValidSOAPResponse",
			nil,
			nil,
			&condition{
				restData: []byte(`{"soap:Envelope": {"soap:Body": {"Response": {"Result": "Success"}}}}`),
			},
			&action{
				xml: []byte(`<?xml version="1.0" encoding="UTF-8"?>
								<soap:Envelope>
									<soap:Header></soap:Header>
									<soap:Body>
										<Response>
											<Result>Success</Result>
										</Response>
									</soap:Body>
								</soap:Envelope>`),
				err:        nil,
				errPattern: nil,
			},
		),
		gen(
			"DecodeError",
			nil,
			nil,
			&condition{
				restData: nil,
			},
			&action{
				xml:        nil,
				err:        app.ErrAppMiddleSOAPRESTDecodeResponseBody,
				errPattern: regexp.MustCompile(core.ErrPrefix + `failed to decode response body.`),
			},
		),
		gen(
			"EmptyJSONKey",
			nil,
			nil,
			&condition{
				restData: []byte(`{"soap:Envelope": {"soap:Body": {"": "EmptyKeyValue"}}}`),
			},
			&action{
				xml: nil,
				err: utilhttp.NewHTTPError(
					app.ErrAppMiddleSOAPRESTMarshalResponseEnvelope.WithoutStack(nil, nil),
					http.StatusInternalServerError,
				),
				errPattern: regexp.MustCompile(core.ErrPrefix + `failed to marshal response envelope.`),
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			s := &soapREST{
				attributeKey:  "@attribute",
				textKey:       "#text",
				namespaceKey:  "_namespace",
				arrayKey:      "item",
				separatorChar: ":",
			}

			wrapper := &wrappedWriter{
				body: bytes.NewBuffer(tt.C().restData),
			}
			result, actualErr := s.convertRESTtoSOAPResponse(wrapper)

			expectedNode, err := parseXML(strings.NewReader(string(tt.A().xml)))
			if err != nil {
				t.Fatalf("Failed to parse expected XML: %v", err)
			}

			actualNode, err := parseXML(strings.NewReader(string(result)))
			if err != nil {
				t.Fatalf("Failed to parse actual XML: %v", err)
			}
			testutil.Diff(t, true, compareNodes(expectedNode, actualNode))

			testutil.DiffError(t, tt.A().err, tt.A().errPattern, actualErr, cmpopts.EquateErrors())
		})
	}
}

func TestXmlElement_MarshalXML(t *testing.T) {
	type condition struct {
		element xmlElement
	}

	type action struct {
		xmlOutput  string
		err        any // error or errorutil.Kind
		errPattern *regexp.Regexp
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()
	gen := testutil.NewCase[*condition, *action]

	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"DebugCase",
			nil,
			nil,
			&condition{
				element: xmlElement{
					XMLName: xml.Name{Local: "#textKey"},
					Content: "textNode",
				},
			},
			&action{
				xmlOutput: `<#textKey>textNode</#textKey>`,
				err:       nil,
			},
		),
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
		gen(
			"InvalidStartElement",
			nil,
			nil,
			&condition{
				element: xmlElement{
					XMLName: xml.Name{Local: ""},
				},
			},
			&action{
				xmlOutput: ``,
				err: utilhttp.NewHTTPError(
					app.ErrAppMiddleSOAPRESTMarshalResponseEnvelope.WithoutStack(nil, nil),
					http.StatusInternalServerError,
				),
				errPattern: regexp.MustCompile(core.ErrPrefix + `failed to marshal response envelope.`),
			},
		),
		gen(
			"InvalidStartElementInChildren",
			nil,
			nil,
			&condition{
				element: xmlElement{
					XMLName: xml.Name{Local: "Test"},
					children: []xmlElement{
						{
							XMLName: xml.Name{Local: ""},
						},
					},
				},
			},
			&action{
				xmlOutput: `<Test>`,
				err: utilhttp.NewHTTPError(
					app.ErrAppMiddleSOAPRESTMarshalResponseEnvelope.WithoutStack(nil, nil),
					http.StatusInternalServerError,
				),
				errPattern: regexp.MustCompile(core.ErrPrefix + `failed to marshal response envelope`),
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
			testutil.DiffError(t, tt.A().err, tt.A().errPattern, err, cmpopts.EquateErrors())

			enc.Flush()
			testutil.Diff(t, string([]byte(tt.A().xmlOutput)), string(bytes.TrimSpace(buf.Bytes())))
		})
	}
}

func TestSOAPREST_CreateSOAPEnvelope(t *testing.T) {
	type condition struct {
		data map[string]any
	}

	type action struct {
		expected *soapEnvelope
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()
	gen := testutil.NewCase[*condition, *action]

	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"TextNodeInBody",
			nil,
			nil,
			&condition{
				data: map[string]any{
					"soap:Envelope": map[string]any{
						"soap:Header": map[string]any{},
						"soap:Body": map[string]any{
							"outerKey": map[string]any{
								"#text": "textNode",
							},
						},
					},
				},
			},
			&action{
				expected: &soapEnvelope{
					Header: &soapHeader{},
					Body: &soapBody{
						Content: []xmlElement{
							{
								XMLName: xml.Name{Local: "outerKey"},
								Content: "textNode",
							},
						},
					},
				},
			},
		),
		gen(
			"EmptyEnvelope",
			nil,
			nil,
			&condition{
				data: map[string]any{
					"soap:Envelope": map[string]any{},
				},
			},
			&action{
				expected: &soapEnvelope{
					Header: &soapHeader{},
					Body:   &soapBody{},
				},
			},
		),
		gen(
			"EnvelopeWithNamespaces",
			nil,
			nil,
			&condition{
				data: map[string]any{
					"soap:Envelope": map[string]any{
						"_namespace": map[string]any{
							"soap": "http://schemas.xmlsoap.org/soap/envelope/",
							"xsi":  "http://www.w3.org/2001/XMLSchema-instance",
						},
						"soap:Body": map[string]any{},
					},
				},
			},
			&action{
				expected: &soapEnvelope{
					ExtraNS: []xml.Attr{
						{Name: xml.Name{Local: "xmlns:soap"}, Value: "http://schemas.xmlsoap.org/soap/envelope/"},
						{Name: xml.Name{Local: "xmlns:xsi"}, Value: "http://www.w3.org/2001/XMLSchema-instance"},
					},
					Header: &soapHeader{},
					Body:   &soapBody{},
				},
			},
		),
		gen(
			"EnvelopeWithHeader",
			nil,
			nil,
			&condition{
				data: map[string]any{
					"soap:Envelope": map[string]any{
						"soap:Header": map[string]any{
							"TestHeader": map[string]any{
								"Key1": "Value1",
								"Key2": "Value2",
							},
						},
						"soap:Body": map[string]any{},
					},
				},
			},
			&action{
				expected: &soapEnvelope{
					Header: &soapHeader{
						Content: []xmlElement{
							{
								XMLName: xml.Name{Local: "TestHeader"},
								children: []xmlElement{
									{XMLName: xml.Name{Local: "Key1"}, Content: "Value1"},
									{XMLName: xml.Name{Local: "Key2"}, Content: "Value2"},
								},
							},
						},
					},
					Body: &soapBody{},
				},
			},
		),
		gen(
			"EnvelopeWithNamespacesAndNullValue",
			nil,
			nil,
			&condition{
				data: map[string]any{
					"soap:Envelope": map[string]any{
						"_namespace": map[string]any{
							"soap": "http://schemas.xmlsoap.org/soap/envelope/",
							"xsi":  "http://www.w3.org/2001/XMLSchema-instance",
						},
						"soap:Body": map[string]any{
							"PartialResponse": map[string]any{
								"Value": nil,
							},
						},
					},
				},
			},
			&action{
				expected: &soapEnvelope{
					ExtraNS: []xml.Attr{
						{Name: xml.Name{Local: "xmlns:xsi"}, Value: "http://www.w3.org/2001/XMLSchema-instance"},
						{Name: xml.Name{Local: "xmlns:soap"}, Value: "http://schemas.xmlsoap.org/soap/envelope/"},
					},
					Header: &soapHeader{},
					Body: &soapBody{
						Content: []xmlElement{
							{
								XMLName: xml.Name{Local: "PartialResponse"},
								children: []xmlElement{
									{
										XMLName: xml.Name{Local: "Value"},
										isNil:   true,
									},
								},
							},
						},
					},
				},
			},
		),
		// <TODO> Implement test cases with []any in hasNullValue
	}

	testutil.Register(table, testCases...)
	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			// <TODO> Implement custom key specific tests
			s := soapREST{
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

			nsManager := &namespaceManager{
				namespaces: make(map[string]string),
			}

			result := s.createSOAPEnvelope(tt.C().data, nsManager)

			opts := []gocmp.Option{
				cmpopts.SortSlices(func(a, b xml.Attr) bool {
					if a.Name.Local < b.Name.Local {
						return true
					}
					if a.Name.Local > b.Name.Local {
						return false
					}
					return a.Value < b.Value
				}),
				cmpopts.SortSlices(func(a, b xmlElement) bool {
					if a.XMLName.Space < b.XMLName.Space {
						return true
					}
					if a.XMLName.Space > b.XMLName.Space {
						return false
					}
					return a.XMLName.Local < b.XMLName.Local
				}),
				testutil.DeepAllowUnexported(
					xmlElement{},
					soapEnvelope{},
					soapHeader{},
					soapBody{},
					namespaceManager{},
				),
				cmpopts.EquateEmpty(),
			}

			testutil.Diff(t, tt.A().expected, result, opts...)
		})
	}
}
func TestSOAPREST_MapToXMLElements(t *testing.T) {
	type condition struct {
		data map[string]any

		attributeKey  string
		namespaceKey  string
		separatorChar string
	}

	type action struct {
		expected []xmlElement
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()
	gen := testutil.NewCase[*condition, *action]

	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"EmptyMap",
			nil,
			nil,
			&condition{
				data: map[string]any{},
			},
			&action{
				expected: nil,
			},
		),
		gen(
			"SingleElement",
			nil,
			nil,
			&condition{
				data: map[string]any{
					"Key": "Value",
				},
			},
			&action{
				expected: []xmlElement{
					{
						XMLName: xml.Name{Local: "Key"},
						Content: "Value",
					},
				},
			},
		),
		gen(
			"MultipleElements",
			nil,
			nil,
			&condition{
				data: map[string]any{
					"Key1": "Value1",
					"Key2": 0,
				},
			},
			&action{
				expected: []xmlElement{
					{
						XMLName: xml.Name{Local: "Key1"},
						Content: "Value1",
					},
					{
						XMLName: xml.Name{Local: "Key2"},
						Content: "0",
					},
				},
			},
		),
		gen(
			"NestedElements",
			nil,
			nil,
			&condition{
				data: map[string]any{
					"Outer": map[string]any{
						"OuterKey": "Value",
						"Inner": map[string]any{
							"InnerKey1": "InnerValue1",
							"InnerKey2": "InnerValue2",
						},
					},
				},
			},
			&action{
				expected: []xmlElement{
					{
						XMLName: xml.Name{Local: "Outer"},
						children: []xmlElement{
							{
								XMLName: xml.Name{Local: "OuterKey"},
								Content: "Value",
							},
							{
								XMLName: xml.Name{Local: "Inner"},
								children: []xmlElement{
									{
										XMLName: xml.Name{Local: "InnerKey1"},
										Content: "InnerValue1",
									},
									{
										XMLName: xml.Name{Local: "InnerKey2"},
										Content: "InnerValue2",
									},
								},
							},
						},
					},
				},
			},
		),
		gen(
			"ElementsWithNilValue",
			nil,
			nil,
			&condition{
				data: map[string]any{
					"OptionalField": nil,
				},
			},
			&action{
				expected: []xmlElement{
					{
						XMLName: xml.Name{Local: "OptionalField"},
						isNil:   true,
					},
				},
			},
		),
		gen(
			"SpecificAttributeKey",
			nil,
			nil,
			&condition{
				attributeKey: "@attr",
				data: map[string]any{
					"@attr": "attribute",
				},
			},
			&action{
				expected: nil,
			},
		),
		gen(
			"SpecificNamespaceKey",
			nil,
			nil,
			&condition{
				namespaceKey: "_ns",
				data: map[string]any{
					"_ns": "namespace",
				},
			},
			&action{
				expected: nil,
			},
		),
		gen(
			"ContainsSeparatorChar",
			nil,
			nil,
			&condition{
				data: map[string]any{
					"ns:key": "value",
				},
			},
			&action{
				expected: []xmlElement{
					{
						XMLName: xml.Name{Space: "ns", Local: "key"},
						Content: "value",
					},
				},
			},
		),
		gen(
			"ValueContainsNamespaceKeyMap",
			nil,
			nil,
			&condition{
				data: map[string]any{
					"soap:Envelope": map[string]any{
						"_namespace": map[string]any{
							"soap": "http://schemas.xmlsoap.org/soap/envelope/",
						},
					},
				},
			},
			&action{
				expected: []xmlElement{
					{
						XMLName: xml.Name{Space: "soap", Local: "Envelope"},
						children: []xmlElement{
							{
								XMLName: xml.Name{Space: "soap", Local: "_namespace"},
								children: []xmlElement{
									{
										XMLName: xml.Name{Space: "soap", Local: "soap"},
										Content: "http://schemas.xmlsoap.org/soap/envelope/",
									},
								},
							},
						},
					},
				},
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			nsManager := &namespaceManager{
				namespaces: make(map[string]string),
			}

			s := soapREST{
				attributeKey:          cmp.Or(tt.C().attributeKey, "@attribute"),
				namespaceKey:          cmp.Or(tt.C().namespaceKey, "_namespace"),
				textKey:               "#text",
				arrayKey:              "item",
				separatorChar:         cmp.Or(tt.C().separatorChar, ":"),
				extractStringElement:  true,
				extractBooleanElement: true,
				extractIntegerElement: true,
				extractFloatElement:   true,
			}

			result := s.mapToXMLElements(tt.C().data, nsManager)

			opts := []gocmp.Option{
				cmpopts.SortSlices(func(a, b xmlElement) bool {
					if a.XMLName.Space < b.XMLName.Space {
						return true
					}
					if a.XMLName.Space > b.XMLName.Space {
						return false
					}
					return a.XMLName.Local < b.XMLName.Local
				}),
				testutil.DeepAllowUnexported(
					xmlElement{},
					soapEnvelope{},
					soapHeader{},
					soapBody{},
					namespaceManager{}),
				cmpopts.EquateEmpty(),
			}

			testutil.Diff(t, tt.A().expected, result, opts...)
		})
	}
}

func TestSOAPREST_MapToXMLElement(t *testing.T) {
	type condition struct {
		key       string
		value     any
		namespace string

		attributeKey string
		textKey      string
		arrayKey     string
	}

	type action struct {
		expected xmlElement
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()
	gen := testutil.NewCase[*condition, *action]

	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"",
			nil,
			nil,
			&condition{
				key:   "ns:key",
				value: "value",
			},
			&action{
				expected: xmlElement{
					XMLName: xml.Name{Space: "ns", Local: "key"},
					Content: "value",
				},
			},
		),
		gen(
			"StringValue",
			nil,
			nil,
			&condition{
				key:   "TestKey",
				value: "TestValue",
			},
			&action{
				expected: xmlElement{
					XMLName: xml.Name{Local: "TestKey"},
					Content: "TestValue",
				},
			},
		),
		gen(
			"IntegerValue",
			nil,
			nil,
			&condition{
				key:   "Integer",
				value: 42,
			},
			&action{
				expected: xmlElement{
					XMLName: xml.Name{Local: "Integer"},
					Content: "42",
				},
			},
		),
		gen(
			"FloatValue",
			nil,
			nil,
			&condition{
				key:   "Float",
				value: 3.14,
			},
			&action{
				expected: xmlElement{
					XMLName: xml.Name{Local: "Float"},
					Content: "3.14",
				},
			},
		),
		gen(
			"FloatValueWithTrailingZero",
			nil,
			nil,
			&condition{
				key:   "Float",
				value: 100.0,
			},
			&action{
				expected: xmlElement{
					XMLName: xml.Name{Local: "Float"},
					Content: "100",
				},
			},
		),
		gen(
			"NilValue",
			nil,
			nil,
			&condition{
				key:   "OptionalElement",
				value: nil,
			},
			&action{
				expected: xmlElement{
					XMLName: xml.Name{Local: "OptionalElement"},
					isNil:   true,
				},
			},
		),
		gen(
			"EmptyValue",
			nil,
			nil,
			&condition{
				key:   "EmptyElement",
				value: map[string]any{},
			},
			&action{
				expected: xmlElement{
					XMLName: xml.Name{Local: "EmptyElement"},
				},
			},
		),
		gen(
			"MapValueWithAttributes",
			nil,
			nil,
			&condition{
				key: "test",
				value: map[string]any{
					"@attribute": map[string]any{
						"localName": map[string]any{
							"key": "value",
						},
					},
				},
			},
			&action{
				expected: xmlElement{
					XMLName: xml.Name{Local: "test"},
					Attrs: []xml.Attr{
						{Name: xml.Name{Local: "localName"}, Value: "map[key:value]"},
					},
				},
			},
		),
		gen(
			"ContainsTextContent",
			nil,
			nil,
			&condition{
				key: "test",
				value: map[string]any{
					"#text": "textContent",
				},
			},
			&action{
				expected: xmlElement{
					XMLName: xml.Name{Local: "test"},
					Content: "textContent",
				},
			},
		),
		gen(
			"ContainsBackspaceAndFormFeedInTextContent",
			nil,
			nil,
			&condition{
				key: "test",
				value: map[string]any{
					"#text": "\b\f\\b\\f",
				},
			},
			&action{
				expected: xmlElement{
					XMLName: xml.Name{Local: "test"},
					Content: "",
				},
			},
		),
		gen(
			"MapValueWithChildElements",
			nil,
			nil,
			&condition{
				key: "test",
				value: map[string]any{
					"childElements": map[string]any{
						"key1": "value1",
						"key2": "value2",
					},
				},
			},
			&action{
				expected: xmlElement{
					XMLName: xml.Name{Local: "test"},
					children: []xmlElement{
						{
							XMLName: xml.Name{Local: "childElements"},
							children: []xmlElement{
								{
									XMLName: xml.Name{Local: "key1"},
									Content: "value1",
								},
								{
									XMLName: xml.Name{Local: "key2"},
									Content: "value2",
								},
							},
						},
					},
				},
			},
		),
		gen(
			"ArrayValue",
			nil,
			nil,
			&condition{
				key:   "item",
				value: []any{"item1", "item2", "item3"},
			},
			&action{
				expected: xmlElement{
					XMLName: xml.Name{Local: "item"},
					children: []xmlElement{
						{
							XMLName: xml.Name{Local: "item"},
							Content: "item1",
						},
						{
							XMLName: xml.Name{Local: "item"},
							Content: "item2",
						},
						{
							XMLName: xml.Name{Local: "item"},
							Content: "item3",
						},
					},
				},
			},
		),
		gen(
			"EmptyArrayValue",
			nil,
			nil,
			&condition{
				key:   "item",
				value: []any{},
			},
			&action{
				expected: xmlElement{
					XMLName: xml.Name{Local: "item"},
				},
			},
		),
		//<TODO> Implement test cases with different config values.
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			s := soapREST{
				attributeKey:          cmp.Or(tt.C().attributeKey, "@attribute"),
				textKey:               cmp.Or(tt.C().textKey, "#text"),
				namespaceKey:          "_namespace",
				arrayKey:              cmp.Or(tt.C().arrayKey, "item"),
				separatorChar:         ":",
				extractStringElement:  true,
				extractBooleanElement: true,
				extractIntegerElement: true,
				extractFloatElement:   true,
			}

			got := s.mapToXMLElement(tt.C().key, tt.C().value, tt.C().namespace)

			opts := []gocmp.Option{
				testutil.DeepAllowUnexported(xmlElement{}),
				cmpopts.EquateEmpty(),
				cmpopts.SortSlices(func(a, b xmlElement) bool {
					return a.XMLName.Local < b.XMLName.Local
				}),
			}

			testutil.Diff(t, tt.A().expected, got, opts...)
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

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"unwrap nil",
			nil,
			nil,
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
			nil,
			nil,
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

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			w := tt.C().ww.Unwrap()
			testutil.Diff(t, tt.A().w, w, gocmp.AllowUnexported(mockResponseWriter{}))
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

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"status code 100",
			nil,
			nil,
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
			nil,
			nil,
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
			nil,
			nil,
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

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			w := httptest.NewRecorder()
			ww := &wrappedWriter{
				ResponseWriter: w,
				written:        tt.C().written,
			}
			ww.WriteHeader(tt.C().code)

			testutil.Diff(t, true, ww.Written())
			testutil.Diff(t, tt.A().code, ww.code)
			testutil.Diff(t, tt.A().written, ww.written)
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

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"status code 100",
			nil,
			nil,
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
			nil,
			nil,
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
			nil,
			nil,
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

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			w := httptest.NewRecorder()
			ww := &wrappedWriter{
				ResponseWriter: w,
				body:           &bytes.Buffer{},
			}
			if tt.C().code > 0 {
				ww.WriteHeader(tt.C().code)
			}
			ww.Write([]byte(tt.C().body))

			testutil.Diff(t, true, ww.Written())
			testutil.Diff(t, tt.A().code, ww.code)

			body, _ := io.ReadAll(ww.body)
			testutil.Diff(t, tt.A().body, string(body))
			testutil.Diff(t, len(tt.A().body), int(ww.ContentLength()))
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

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"don't write status code",
			nil,
			nil,
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
			nil,
			nil,
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

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			if tt.C().write {
				tt.C().ww.WriteHeader(999)
			}
			testutil.Diff(t, tt.A().written, tt.C().ww.written)
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

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"status code 100",
			nil,
			nil,
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
			nil,
			nil,
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
			nil,
			nil,
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

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			tt.C().ww.WriteHeader(tt.C().code)
			testutil.Diff(t, tt.A().code, tt.C().ww.StatusCode())
		})
	}
}

func TestNamespaceContext_AddNamespace(t *testing.T) {
	type condition struct {
		prefix string
		uri    string
	}

	type action struct {
		prefix string
		uri    string
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"Prefix exists in namespaceContext",
			nil,
			nil,
			&condition{
				prefix: "test",
				uri:    "http://test.com/",
			},
			&action{
				prefix: "test",
				uri:    "http://test.com/",
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			nc := &namespaceContext{
				prefixToURI: map[string]string{},
				uriToPrefix: map[string]string{},
			}

			nc.addNamespace(tt.C().prefix, tt.C().uri)
			testutil.Diff(t, tt.A().prefix, nc.uriToPrefix[tt.C().uri])
			testutil.Diff(t, tt.A().uri, nc.prefixToURI[tt.C().prefix])
		})
	}
}

func TestNamespaceContext_GetPrefix(t *testing.T) {
	type condition struct {
		prefix string
		uri    string

		notExist bool
	}

	type action struct {
		prefix string
		uri    string
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"Prefix exists in namespaceContext",
			nil,
			nil,
			&condition{
				prefix:   "test",
				uri:      "http://test.com/",
				notExist: false,
			},
			&action{
				prefix: "test",
				uri:    "http://test.com/",
			},
		),
		gen(
			"Prefix does not exist in namespaceContext",
			nil,
			nil,
			&condition{
				notExist: true,
			},
			&action{
				prefix: "",
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			nc := &namespaceContext{
				prefixToURI: map[string]string{
					tt.C().prefix: tt.C().uri,
				},
				uriToPrefix: map[string]string{
					tt.C().uri: tt.C().prefix,
				},
			}

			if tt.C().notExist {
				testutil.Diff(t, tt.A().prefix, nc.getPrefix("notExists"))
			} else {
				testutil.Diff(t, tt.A().prefix, nc.getPrefix(tt.C().uri))
			}
		})
	}
}

func TestNamespaceManager_AddNamespace(t *testing.T) {
	type condition struct {
		prefix      string
		originalUri string
		anotherUri  string
	}

	type action struct {
		uri string
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"add namespace once",
			nil,
			nil,
			&condition{
				prefix:      "test",
				originalUri: "http://test.com/",
			},
			&action{
				uri: "http://test.com/",
			},
		),
		gen(
			"add namespace multiple times",
			nil,
			nil,
			&condition{
				prefix:      "test",
				originalUri: "http://original.com/",
				anotherUri:  "http://another.com/",
			},
			&action{
				uri: "http://original.com/",
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			nm := &namespaceManager{
				namespaces: map[string]string{},
			}
			nm.addNamespace(tt.C().prefix, tt.C().originalUri)
			testutil.Diff(t, tt.A().uri, nm.namespaces[tt.C().prefix])

		})
	}
}

func TestSOAPREST_ParseValue(t *testing.T) {
	type condition struct {
		content               string
		extractStringElement  bool
		extractBooleanElement bool
		extractIntegerElement bool
		extractFloatElement   bool
	}

	type action struct {
		expect any
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"string element",
			nil,
			nil,
			&condition{
				content:              "test",
				extractStringElement: true,
			},
			&action{
				expect: "test",
			},
		),
		gen(
			"escaped string element",
			nil,
			nil,
			&condition{
				content:              "\"test\"",
				extractStringElement: true,
			},
			&action{
				expect: "test",
			},
		),
		gen(
			"escaped string element but not extracted",
			nil,
			nil,
			&condition{
				content:              "\"test\"",
				extractStringElement: false,
			},
			&action{
				expect: "\"test\"",
			},
		),
		gen(
			"boolean element (true)",
			nil,
			nil,
			&condition{
				content:               "true",
				extractBooleanElement: true,
			},
			&action{
				expect: true,
			},
		),
		gen(
			"boolean element (false)",
			nil,
			nil,
			&condition{
				content:               "false",
				extractBooleanElement: true,
			},
			&action{
				expect: false,
			},
		),
		gen(
			"boolean element (true) but not extracted",
			nil,
			nil,
			&condition{
				content:               "true",
				extractBooleanElement: false,
			},
			&action{
				expect: "true",
			},
		),
		gen(
			"boolean element (false) but not extracted",
			nil,
			nil,
			&condition{
				content:               "false",
				extractBooleanElement: false,
			},
			&action{
				expect: "false",
			},
		),
		gen(
			"integer element (zero)",
			nil,
			nil,
			&condition{
				content:               "0",
				extractIntegerElement: true,
			},
			&action{
				expect: int64(0),
			},
		),
		gen(
			"integer element (positive)",
			nil,
			nil,
			&condition{
				content:               "100",
				extractIntegerElement: true,
			},
			&action{
				expect: int64(100),
			},
		),
		gen(
			"integer element (negative)",
			nil,
			nil,
			&condition{
				content:               "-100",
				extractIntegerElement: true,
			},
			&action{
				expect: int64(-100),
			},
		),
		gen(
			"integer element (upper bound)",
			nil,
			nil,
			&condition{
				content:               "9223372036854775807",
				extractIntegerElement: true,
			},
			&action{
				expect: int64(math.MaxInt64),
			},
		),
		gen(
			"integer element (lower bound)",
			nil,
			nil,
			&condition{
				content:               "-9223372036854775808",
				extractIntegerElement: true,
			},
			&action{
				expect: int64(math.MinInt64),
			},
		),
		gen(
			"integer element (over upper bound)",
			nil,
			nil,
			&condition{
				content:               "9223372036854775808",
				extractIntegerElement: true,
			},
			&action{
				expect: "9223372036854775808",
			},
		),
		gen(
			"integer element (over lower bound)",
			nil,
			nil,
			&condition{
				content:               "-9223372036854775809",
				extractIntegerElement: true,
			},
			&action{
				expect: "-9223372036854775809",
			},
		),
		gen(
			"integer element (zero) but not extracted",
			nil,
			nil,
			&condition{
				content:               "0",
				extractIntegerElement: false,
			},
			&action{
				expect: "0",
			},
		),
		gen(
			"integer element (positive) but not extracted",
			nil,
			nil,
			&condition{
				content:               "100",
				extractIntegerElement: false,
			},
			&action{
				expect: "100",
			},
		),
		gen(
			"integer element (negative) but not extracted",
			nil,
			nil,
			&condition{
				content:               "-100",
				extractIntegerElement: false,
			},
			&action{
				expect: "-100",
			},
		),
		gen(
			"integer element (upper bound) but not extracted",
			nil,
			nil,
			&condition{
				content:               "9223372036854775807",
				extractIntegerElement: false,
			},
			&action{
				expect: "9223372036854775807",
			},
		),
		gen(
			"integer element (lower bound) but not extracted",
			nil,
			nil,
			&condition{
				content:               "-9223372036854775808",
				extractIntegerElement: false,
			},
			&action{
				expect: "-9223372036854775808",
			},
		),
		gen(
			"integer element (over upper bound) but not extracted",
			nil,
			nil,
			&condition{
				content:               "9223372036854775808",
				extractIntegerElement: false,
			},
			&action{
				expect: "9223372036854775808",
			},
		),
		gen(
			"integer element (over lower bound) but not extracted",
			nil,
			nil,
			&condition{
				content:               "-9223372036854775809",
				extractIntegerElement: false,
			},
			&action{
				expect: "-9223372036854775809",
			},
		),
		gen(
			"float element (zero)",
			nil,
			nil,
			&condition{
				content:             "0.00",
				extractFloatElement: true,
			},
			&action{
				expect: float64(0),
			},
		),
		gen(
			"float element (positive)",
			nil,
			nil,
			&condition{
				content:             "100.123000",
				extractFloatElement: true,
			},
			&action{
				expect: float64(100.123),
			},
		),
		gen(
			"float element (negative)",
			nil,
			nil,
			&condition{
				content:             "-100.123000",
				extractFloatElement: true,
			},
			&action{
				expect: float64(-100.123),
			},
		),
		gen(
			"float element (upper bound)",
			nil,
			nil,
			&condition{
				content:             "1.79769313486231570814527423731704356798070e+308",
				extractFloatElement: true,
			},
			&action{
				expect: math.MaxFloat64,
			},
		),
		gen(
			"float element (lower bound)",
			nil,
			nil,
			&condition{
				content:             "4.9406564584124654417656879286822137236505980e-324",
				extractFloatElement: true,
			},
			&action{
				expect: math.SmallestNonzeroFloat64,
			},
		),
		gen(
			"float element (positive false precision)",
			nil,
			nil,
			&condition{
				content:             "0.1234567890123456789",
				extractFloatElement: true,
			},
			&action{
				expect: float64(0.12345678901234568),
			},
		),
		gen(
			"float element (negative false precision)",
			nil,
			nil,
			&condition{
				content:             "-0.1234567890123456789",
				extractFloatElement: true,
			},
			&action{
				expect: float64(-0.12345678901234568),
			},
		),
		gen(
			"float element (zero) but not extracted",
			nil,
			nil,
			&condition{
				content:             "0.00",
				extractFloatElement: false,
			},
			&action{
				expect: "0.00",
			},
		),
		gen(
			"float element (positive) but not extracted",
			nil,
			nil,
			&condition{
				content:             "100.123000",
				extractFloatElement: false,
			},
			&action{
				expect: "100.123000",
			},
		),
		gen(
			"float element (negative) but not extracted",
			nil,
			nil,
			&condition{
				content:             "-100.123000",
				extractFloatElement: false,
			},
			&action{
				expect: "-100.123000",
			},
		),
		gen(
			"float element (upper bound) but not extracted",
			nil,
			nil,
			&condition{
				content:             "1.79769313486231570814527423731704356798070e+308",
				extractFloatElement: false,
			},
			&action{
				expect: "1.79769313486231570814527423731704356798070e+308",
			},
		),
		gen(
			"float element (lower bound) but not extracted",
			nil,
			nil,
			&condition{
				content:             "4.9406564584124654417656879286822137236505980e-324",
				extractFloatElement: false,
			},
			&action{
				expect: "4.9406564584124654417656879286822137236505980e-324",
			},
		),
		gen(
			"float element (positive false precision) but not extracted",
			nil,
			nil,
			&condition{
				content:             "0.1234567890123456789",
				extractFloatElement: false,
			},
			&action{
				expect: "0.1234567890123456789",
			},
		),
		gen(
			"float element (negative false precision) but not extracted",
			nil,
			nil,
			&condition{
				content:             "-0.1234567890123456789",
				extractFloatElement: false,
			},
			&action{
				expect: "-0.1234567890123456789",
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			sr := soapREST{
				extractStringElement:  tt.C().extractStringElement,
				extractBooleanElement: tt.C().extractBooleanElement,
				extractIntegerElement: tt.C().extractIntegerElement,
				extractFloatElement:   tt.C().extractFloatElement,
			}

			testutil.Diff(t, tt.A().expect, sr.parseValue(tt.C().content))
		})
	}
}

func TestHasNullValue(t *testing.T) {
	type condition struct {
		data any
	}
	type action struct {
		expect bool
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"contains nil data",
			nil,
			nil,
			&condition{
				data: nil,
			},
			&action{
				expect: true,
			},
		),
		gen(
			"contains map[string]any data without nil",
			nil,
			nil,
			&condition{
				data: map[string]any{
					"testKey": "nonNil",
				},
			},
			&action{
				expect: false,
			},
		),
		gen(
			"contains map[string]any data with nil",
			nil,
			nil,
			&condition{
				data: map[string]any{
					"nilKey": nil,
				},
			},
			&action{
				expect: true,
			},
		),
		gen(
			"contains map[string]any data with multiple nil",
			nil,
			nil,
			&condition{
				data: map[string]any{
					"nonNilKey": "nonNil",
					"otherKey":  "other",
					"nilKey1":   nil,
					"nilKey2":   nil,
				},
			},
			&action{
				expect: true,
			},
		),
		gen(
			"contains []any data without nil",
			nil,
			nil,
			&condition{
				data: []any{
					"nonNil",
				},
			},
			&action{
				expect: false,
			},
		),
		gen(
			"contains []any data with nil",
			nil,
			nil,
			&condition{
				data: []any{
					nil,
				},
			},
			&action{
				expect: true,
			},
		),
		gen(
			"contains []any data with multiple nil",
			nil,
			nil,
			&condition{
				data: []any{
					"nonNil",
					"other",
					nil,
					nil,
				},
			},
			&action{
				expect: true,
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			testutil.Diff(t, tt.A().expect, hasNullValue(tt.C().data))
		})
	}
}
