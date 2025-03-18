package soaprest

import (
	"bytes"
	"cmp"
	"encoding/json"
	"encoding/xml"
	"errors"
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

		attributeKey string
		textKey      string
		namespaceKey string
		arrayKey     string

		soapNamespacePrefix string

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
				body: `{"soap_Envelope": {
							"_namespace": {
								"ns": "http://example.com/",
								"soap": "http://schemas.xmlsoap.org/soap/envelope/"
							},
							"soap_Body": {
								"ns_Test": {
									"@attribute": {
										"testAttributeKey": "someValue"
									},
									"ns_Value": 123
								},
								"ns_Array": {
									"item": [100, 3.14, true, "someText"]
								}
							},
							"soap_Header": {
								"#text": "double quoted text",
								"ns_Auth": {
									"ns_Username": "TestUser",
									"ns_Password": "password"
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
				body: `{"soap_Envelope": {
							"_namespace": {
								"ns": "http://example.com/",
								"soap": "http://schemas.xmlsoap.org/soap/envelope/"
							},
							"soap_Body": {
								"ns_Test": {
									"@attribute": {
										"testAttributeKey": "someValue"
									},
									"ns_Value": 123
								},
								"ns_Array": {
									"item": [100, 3.14, true, "someText"]
								}
							},
							"soap_Header": {
								"#text": "double quoted text",
								"ns_Auth": {
									"ns_Username": "TestUser",
									"ns_Password": "password"
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
				err:  app.ErrAppMiddleSOAPRESTUnmarshalRequestBody,
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

				attributeKey: "@attr!",
				textKey:      "#textKey",
				namespaceKey: "_ns",
				arrayKey:     "element",

				extractStringElement:  true,
				extractBooleanElement: true,
				extractIntegerElement: true,
				extractFloatElement:   true,
			},
			&action{
				body: `{"soap_Envelope": {
							"_ns": {
								"ns": "http://example.com/",
								"soap": "http://schemas.xmlsoap.org/soap/envelope/"
							},
							"soap_Body": {
								"ns_Test": {
									"@attr!": {
										"testAttributeKey": "someValue"
									},
									"ns_Value": 123
								},
								"ns_Array": {
									"item": [100, 3.14, true, "someText"]
								}
							},
							"soap_Header": {
								"#textKey": "double quoted text",
								"ns_Auth": {
									"ns_Username": "TestUser",
									"ns_Password": "password"
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
				body: `{"soap_Envelope": {
							"_namespace": {
								"ns": "http://example.com/",
								"soap": "http://schemas.xmlsoap.org/soap/envelope/"
							},
							"soap_Body": {
								"ns_Test": {
									"@attribute": {
										"testAttributeKey": "someValue"
									},
									"ns_Value": "123"
								},
								"ns_Array": {
									"item": ["100", "3.14", "true", "someText"]
								}
							},
							"soap_Header": {
								"#text": "\"double quoted text\"",
								"ns_Auth": {
									"ns_Username": "TestUser",
									"ns_Password": "password"
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
			"WrongURLPath",
			nil,
			nil,
			&condition{
				method:      http.MethodPost,
				contentType: "text/xml",
				body:        `{}`,

				paths: &testMatcher{match: false},

				extractStringElement:  false,
				extractBooleanElement: false,
				extractIntegerElement: false,
				extractFloatElement:   false,
			},
			&action{
				body: `{}`,
				code: 200,
			},
		),
		gen(
			"DefaultNamespaceDeclaration",
			nil,
			nil,
			&condition{
				method:      http.MethodPost,
				contentType: "text/xml",
				body: `<?xml version="1.0" encoding="UTF-8"?>
						<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/" xmlns:ns="http://testNamespace.com/">
							<soap:Header></soap:Header>
							<soap:Body>
								<DefaultNamespace xmlns="http://example.com/">
									<Item>default</Item>
									<ns:Item>defaultWithAnotherNamespace</ns:Item>
								</DefaultNamespace>
							</soap:Body>
						</soap:Envelope>`,

				paths: &testMatcher{match: true},

				extractStringElement:  false,
				extractBooleanElement: false,
				extractIntegerElement: false,
				extractFloatElement:   false,
			},
			&action{
				body: `{"soap_Envelope": {
							"_namespace": {
								"ns": "http://testNamespace.com/",
								"soap": "http://schemas.xmlsoap.org/soap/envelope/"
							},
							"soap_Body": {
								"xmlns_DefaultNamespace": {
									"_namespace": {
										"xmlns":"http://example.com/"
									},
									"ns_Item": "defaultWithAnotherNamespace",
									"xmlns_Item": "default"
								}
							},
							"soap_Header": {}
				}}`,
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
				eh:           meh,
				paths:        tt.C().paths,
				attributeKey: cmp.Or(tt.C().attributeKey, "@attribute"),
				textKey:      cmp.Or(tt.C().textKey, "#text"),
				namespaceKey: cmp.Or(tt.C().namespaceKey, "_namespace"),
				arrayKey:     cmp.Or(tt.C().arrayKey, "item"),

				soapNamespacePrefix: cmp.Or(tt.C().soapNamespacePrefix, "soap"),

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
			testutil.Diff(t, tt.A().code, resp.Code)
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

	attrMapA := make(map[string]string)
	attrMapB := make(map[string]string)

	for _, attr := range a.Attr {
		key := attr.Name.Space + ":" + attr.Name.Local
		attrMapA[key] = attr.Value
	}

	for _, attr := range b.Attr {
		key := attr.Name.Space + ":" + attr.Name.Local
		attrMapB[key] = attr.Value
	}

	for key, valueA := range attrMapA {
		if valueB, ok := attrMapB[key]; !ok || valueA != valueB {
			return false
		}
	}

	if a.Text != b.Text {
		return false
	}

	if len(a.ChildNodes) != len(b.ChildNodes) {
		return false
	}

	childMapA := make(map[string]*testNode)
	childMapB := make(map[string]*testNode)

	for i, child := range a.ChildNodes {
		key := child.Name.Space + ":" + child.Name.Local
		childMapA[key] = a.ChildNodes[i]
	}

	for i, child := range b.ChildNodes {
		key := child.Name.Space + ":" + child.Name.Local
		childMapB[key] = b.ChildNodes[i]
	}

	for key, nodeA := range childMapA {
		if nodeB, ok := childMapB[key]; !ok || !compareNodes(nodeA, nodeB) {
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
		paths       *testMatcher

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
                    "soap_Envelope": {
                        "nsKey": {
                            "soap": "http://schemas.xmlsoap.org/soap/envelope/",
							"ns": "http://example.com/"
                        },
                        "soap_Body": {
                            "ns_Test": {
                                "attrKey": {
                                    "testAttributeKey": "someValue"
                                },
                                "ns_Value": 123
                            },
                            "ns_Array": {
                                "item": [100, 3.14, true, "someText"]
                            }
                        },
                        "soap_Header": {
                            "textKey": "double quoted text",
                            "ns_Auth": {
                                "ns_Username": "TestUser",
                                "ns_Password": "password"
                            }
                        }
                    }
                }`,
				paths: &testMatcher{match: true},
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
				paths:       &testMatcher{match: true},
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
				paths:       &testMatcher{match: true},
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
				eh:           meh,
				attributeKey: "attrKey",
				textKey:      "textKey",
				namespaceKey: "nsKey",
				arrayKey:     "arrayKey",

				paths: tt.C().paths,

				soapNamespacePrefix: "soap",

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
						"attrKey": map[string]string{
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
					"ns_Test": map[string]any{
						"nsKey": map[string]string{
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
			"ContentWithSpace",
			nil,
			nil,
			&condition{
				xmlInput: xmlNode{
					XMLName: xml.Name{Local: "Test"},
					Content: " ",
				},
			},
			&action{
				expected: map[string]any{},
			},
		),
		gen(
			"ContentWithSpaceAndChildContent",
			nil,
			nil,
			&condition{
				xmlInput: xmlNode{
					XMLName: xml.Name{Local: "Test"},
					Content: " test ",
					Children: []xmlNode{
						{
							XMLName: xml.Name{Local: "Outer"},
							Attrs: []xml.Attr{
								{
									Name:  xml.Name{Local: "Inner"},
									Value: "testValue",
								},
							},
							Content: "testContent",
						},
					},
				},
			},
			&action{
				expected: map[string]any{
					"Test": map[string]any{
						"textKey": " test ",
						"Outer": map[string]any{
							"textKey": "testContent",
							"attrKey": map[string]string{
								"Inner": "testValue",
							},
						},
					},
				},
			},
		),
		gen(
			"DefaultNamespaceDeclaration",
			nil,
			nil,
			&condition{
				xmlInput: xmlNode{
					XMLName: xml.Name{
						Space: "http://schemas.xmlsoap.org/soap/envelope/",
						Local: "Envelope",
					},
					Children: []xmlNode{
						{
							XMLName: xml.Name{
								Space: "http://schemas.xmlsoap.org/soap/envelope/",
								Local: "Header",
							},
						},
						{
							XMLName: xml.Name{
								Space: "http://schemas.xmlsoap.org/soap/envelope/",
								Local: "Body",
							},
							Children: []xmlNode{
								{
									XMLName: xml.Name{
										Space: "",
										Local: "DefaultNamespace",
									},
									Attrs: []xml.Attr{
										{
											Name: xml.Name{
												Local: "xmlns",
											},
											Value: "http://example.com/",
										},
									},
									Children: []xmlNode{
										{
											XMLName: xml.Name{
												Local: "Item",
											},
											Content: "default",
										},
										{
											XMLName: xml.Name{
												Space: "http://testNamespace.com/",
												Local: "Item",
											},
											Content: "anotherItemWithNamespace",
										},
									},
								},
							},
						},
					},
				},
			},
			&action{
				expected: map[string]any{
					"soap_Envelope": map[string]any{
						"soap_Body": map[string]any{
							"DefaultNamespace": map[string]any{
								"ns_Item": "anotherItemWithNamespace",
								"Item":    "default",
								"nsKey": map[string]string{
									"xmlns": "http://example.com/",
								},
							},
						},
						"soap_Header": map[string]any{},
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
				attributeKey: "attrKey",
				textKey:      "textKey",
				namespaceKey: "nsKey",
				arrayKey:     "arrayKey",

				soapNamespacePrefix: "soap",

				extractStringElement:  false,
				extractBooleanElement: true,
				extractIntegerElement: true,
				extractFloatElement:   true,
			}

			nsCtx := &namespaceContext{
				prefixToURI: map[string]string{
					"soap": "http://schemas.xmlsoap.org/soap/envelope/",
					"ns":   "http://testNamespace.com/",
				},
				uriToPrefix: map[string]string{
					"http://schemas.xmlsoap.org/soap/envelope/": "soap",
					"http://testNamespace.com/":                 "ns",
				},
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
		errPattern string
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
				restData: []byte(`{
									"soap_Envelope": {
										"soap_Body": {
											"Response": {
												"Result": "Success"
											}
										}
									}}`),
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
			},
		),
		gen(
			"MultipleAttributes",
			nil,
			nil,
			&condition{
				restData: []byte(`{
									"soap_Envelope": {
										"nsKey": {
											"soap": "http://schemas.xmlsoap.org/soap/envelope/"
										},
										"attrKey": {
											"testKey": "testValue"
										},
										"soap_Body": {
											"attrKey": {
												"testBodyKey": "testBodyValue"
											}
										},
										"soap_Header": {
											"attrKey": {
												"testHeaderKey": "testHeaderValue"
											}
										}
									}}`),
			},
			&action{
				xml: []byte(`<?xml version="1.0" encoding="UTF-8"?>
								<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/" testKey="testValue">
									<soap:Header testHeaderKey="testHeaderValue"></soap:Header>
									<soap:Body testBodyKey="testBodyValue"></soap:Body>
								</soap:Envelope>`),
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
				errPattern: "EOF",
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			s := &soapREST{
				attributeKey: "attrKey",
				textKey:      "textKey",
				namespaceKey: "nsKey",
				arrayKey:     "arrayKey",

				soapNamespacePrefix: "soap",
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
			if actualErr != nil {
				testutil.Diff(t, tt.A().errPattern, actualErr.Error())
			}
		})
	}
}

func TestXmlElement_MarshalXML(t *testing.T) {
	type condition struct {
		element xmlElement
	}

	type action struct {
		xmlOutput string
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
			},
		),
		gen(
			"TextNodeDirectlyUnderHeaderOrBody",
			nil,
			nil,
			&condition{
				element: xmlElement{
					XMLName: xml.Name{Local: ""},
					Content: "TextNode",
					isNil:   true,
				},
			},
			&action{
				xmlOutput: "\n    TextNode\n  ",
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			var buf bytes.Buffer

			enc := xml.NewEncoder(&buf)
			tt.C().element.MarshalXML(enc, xml.StartElement{Name: tt.C().element.XMLName})

			enc.Flush()
			testutil.Diff(t, string([]byte(tt.A().xmlOutput)), buf.String())
		})
	}
}

func TestSOAPREST_CreateSOAPEnvelope(t *testing.T) {
	type condition struct {
		data       map[string]any
		namespaces map[string]string
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
					"soap_Envelope": map[string]any{
						"soap_Header": map[string]any{},
						"soap_Body": map[string]any{
							"outerKey": map[string]any{
								"textKey": "textNode",
							},
						},
					},
				},
				namespaces: map[string]string{},
			},
			&action{
				expected: &soapEnvelope{
					prefix: "soap",
					Header: &soapHeader{
						prefix: "soap",
					},
					Body: &soapBody{
						Content: []xmlElement{
							{
								XMLName: xml.Name{Local: "outerKey"},
								Content: "textNode",
							},
						},
						prefix: "soap",
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
					"soap_Envelope": map[string]any{},
				},
				namespaces: map[string]string{},
			},
			&action{
				expected: &soapEnvelope{
					prefix: "soap",
					Header: &soapHeader{
						prefix: "soap",
					},
					Body: &soapBody{
						prefix: "soap",
					},
				},
			},
		),
		gen(
			"EnvelopeWithNamespaces",
			nil,
			nil,
			&condition{
				data: map[string]any{
					"soap_Envelope": map[string]any{
						"nsKey": map[string]any{
							"soap": "http://schemas.xmlsoap.org/soap/envelope/",
							"xsi":  "http://www.w3.org/2001/XMLSchema-instance",
						},
						"soap_Body": map[string]any{},
					},
				},
				namespaces: map[string]string{},
			},
			&action{
				expected: &soapEnvelope{
					prefix: "soap",
					ExtraNS: []xml.Attr{
						{Name: xml.Name{Local: "xmlns:soap"}, Value: "http://schemas.xmlsoap.org/soap/envelope/"},
						{Name: xml.Name{Local: "xmlns:xsi"}, Value: "http://www.w3.org/2001/XMLSchema-instance"},
					},
					Header: &soapHeader{
						prefix: "soap",
					},
					Body: &soapBody{
						prefix: "soap",
					},
				},
			},
		),
		gen(
			"EnvelopeWithHeader",
			nil,
			nil,
			&condition{
				data: map[string]any{
					"soap_Envelope": map[string]any{
						"soap_Header": map[string]any{
							"TestHeader": map[string]any{
								"Key1": "Value1",
								"Key2": "Value2",
							},
						},
						"soap_Body": map[string]any{},
					},
				},
				namespaces: map[string]string{},
			},
			&action{
				expected: &soapEnvelope{
					prefix: "soap",
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
						prefix: "soap",
					},
					Body: &soapBody{
						prefix: "soap",
					},
				},
			},
		),
		gen(
			"EnvelopeWithNamespacesAndNullValue",
			nil,
			nil,
			&condition{
				data: map[string]any{
					"soap_Envelope": map[string]any{
						"nsKey": map[string]any{
							"soap": "http://schemas.xmlsoap.org/soap/envelope/",
							"xsi":  "http://www.w3.org/2001/XMLSchema-instance",
						},
						"soap_Body": map[string]any{
							"PartialResponse": map[string]any{
								"Value": nil,
							},
						},
					},
				},
				namespaces: map[string]string{},
			},
			&action{
				expected: &soapEnvelope{
					prefix: "soap",
					ExtraNS: []xml.Attr{
						{Name: xml.Name{Local: "xmlns:xsi"}, Value: "http://www.w3.org/2001/XMLSchema-instance"},
						{Name: xml.Name{Local: "xmlns:soap"}, Value: "http://schemas.xmlsoap.org/soap/envelope/"},
					},
					Header: &soapHeader{
						prefix: "soap",
					},
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
						prefix: "soap",
					},
				},
			},
		),
		gen(
			"AttributesDirectlyUnderTheSOAPEnvelope",
			nil,
			nil,
			&condition{
				data: map[string]any{
					"soap_Envelope": map[string]any{
						"nsKey": map[string]any{
							"soap": "http://schemas.xmlsoap.org/soap/envelope/",
						},
						"attrKey": map[string]any{
							"testAttr": "exampleAttribute",
						},
					},
				},
				namespaces: map[string]string{},
			},
			&action{
				expected: &soapEnvelope{
					prefix: "soap",
					ExtraNS: []xml.Attr{
						{Name: xml.Name{Local: "xmlns:soap"}, Value: "http://schemas.xmlsoap.org/soap/envelope/"},
					},
					Attrs: []xml.Attr{
						{
							Name:  xml.Name{Local: "testAttr"},
							Value: "exampleAttribute",
						},
					},
					Header: &soapHeader{
						prefix: "soap",
					},
					Body: &soapBody{
						prefix: "soap",
					},
				},
			},
		),
		gen(
			"ChildElementsUnderTheSOAPEnvelopeThatDoNotHave`map[string]any`",
			nil,
			nil,
			&condition{
				data: map[string]any{
					"soap_Envelope": map[string]any{
						"testKey": "testValue",
					},
				},
				namespaces: map[string]string{},
			},
			&action{
				expected: &soapEnvelope{
					prefix: "soap",
					Header: &soapHeader{
						prefix: "soap",
					},
					Body: &soapBody{
						prefix: "soap",
					},
				},
			},
		),
		gen(
			"AttributesDirectlyUnderTheSOAPHeader",
			nil,
			nil,
			&condition{
				data: map[string]any{
					"soap_Envelope": map[string]any{
						"soap_Header": map[string]any{
							"attrKey": map[string]any{
								"testAttr": "exampleAttribute",
							},
						},
					},
				},
				namespaces: map[string]string{},
			},
			&action{
				expected: &soapEnvelope{
					prefix: "soap",
					Header: &soapHeader{
						prefix: "soap",
						Attrs: []xml.Attr{
							{
								Name:  xml.Name{Local: "testAttr"},
								Value: "exampleAttribute",
							},
						},
					},
					Body: &soapBody{
						prefix: "soap",
					},
				},
			},
		),
		gen(
			"AttributesDirectlyUnderTheSOAPBody",
			nil,
			nil,
			&condition{
				data: map[string]any{
					"soap_Envelope": map[string]any{
						"soap_Body": map[string]any{
							"attrKey": map[string]any{
								"testAttr": "exampleAttribute",
							},
							"textKey": "testText",
						},
					},
				},
				namespaces: map[string]string{},
			},
			&action{
				expected: &soapEnvelope{
					prefix: "soap",
					Header: &soapHeader{
						prefix: "soap",
					},
					Body: &soapBody{
						Attrs: []xml.Attr{
							{
								Name:  xml.Name{Local: "testAttr"},
								Value: "exampleAttribute",
							},
						},
						Content: []xmlElement{
							{
								Content: "testText",
							},
						},
						prefix: "soap",
					},
				},
			},
		),
		gen(
			"DefaultNamespaceDeclaration(EmptyPrefix)",
			nil,
			nil,
			&condition{
				data: map[string]any{
					"soap_Envelope": map[string]any{
						"soap_Body": map[string]any{
							"attrKey": map[string]any{
								"": "http://example.com/",
							},
						},
					},
				},
				namespaces: map[string]string{
					"": "http://example.com/",
				},
			},
			&action{
				expected: &soapEnvelope{
					prefix: "soap",
					Header: &soapHeader{
						prefix: "soap",
					},
					Body: &soapBody{
						Attrs: []xml.Attr{
							{
								Name:  xml.Name{Local: ""},
								Value: "http://example.com/",
							},
						},
						prefix: "soap",
					},
					ExtraNS: []xml.Attr{
						{
							Name: xml.Name{
								Local: "xmlns",
							},
							Value: "http://example.com/",
						},
					},
				},
			},
		),
		gen(
			"TextContentWithOtherElement",
			nil,
			nil,
			&condition{
				data: map[string]any{
					"soap_Envelope": map[string]any{
						"soap_Body": map[string]any{
							"attrKey": map[string]any{
								"testAttr": "exampleAttr",
							},
							"textKey": "testText",
							"siblingElement": map[string]any{
								"childElement": "childValue",
							},
						},
					},
				},
				namespaces: map[string]string{},
			},
			&action{
				expected: &soapEnvelope{
					prefix: "soap",
					Header: &soapHeader{
						prefix: "soap",
					},
					Body: &soapBody{
						Attrs: []xml.Attr{
							{
								Name:  xml.Name{Local: "testAttr"},
								Value: "exampleAttr",
							},
						},
						Content: []xmlElement{
							{
								XMLName:     xml.Name{Local: ""},
								Content:     "testText",
								hasSiblings: true,
							},
							{
								XMLName: xml.Name{Local: "siblingElement"},
								children: []xmlElement{
									{
										XMLName: xml.Name{Local: "childElement"},
										Content: "childValue",
									},
								},
							},
						},
						prefix: "soap",
					},
				},
			},
		),
		gen(
			"KeyStartsWithSeparatorChar",
			nil,
			nil,
			&condition{
				data: map[string]any{
					"soap_Envelope": map[string]any{
						"soap_Body": map[string]any{
							"_text_key": "testText",
						},
					},
				},
				namespaces: map[string]string{},
			},
			&action{
				expected: &soapEnvelope{
					prefix: "soap",
					Header: &soapHeader{
						prefix: "soap",
					},
					Body: &soapBody{
						Content: []xmlElement{
							{
								XMLName: xml.Name{Local: "_text_key"},
								Content: "testText",
							},
						},
						prefix: "soap",
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
				attributeKey: "attrKey",
				textKey:      "textKey",
				namespaceKey: "nsKey",
				arrayKey:     "arrayKey",

				soapNamespacePrefix: "soap",

				extractStringElement:  true,
				extractBooleanElement: true,
				extractIntegerElement: true,
				extractFloatElement:   true,
			}

			nsManager := &namespaceManager{
				namespaces: tt.C().namespaces,
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
					namespaceManager{},
					soapEnvelope{},
					soapHeader{},
					soapBody{},
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

		attributeKey string
		namespaceKey string
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
					"ns_key": "value",
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
			"MultipleNamespaceElements",
			nil,
			nil,
			&condition{
				data: map[string]any{
					"data": map[string]any{
						"nsKey": map[string]any{
							"test": "http://example.com/",
						},
					},
				},
			},
			&action{
				expected: []xmlElement{
					{
						XMLName: xml.Name{Local: "data"},
					},
				},
			},
		),
		gen(
			"KeyStartsWithSeparatorChar",
			nil,
			nil,
			&condition{
				data: map[string]any{
					"soap_Body": map[string]any{
						"_test_key": "testValue",
					},
				},
			},
			&action{
				expected: []xmlElement{
					{
						XMLName: xml.Name{Space: "soap", Local: "Body"},
						children: []xmlElement{
							{
								XMLName: xml.Name{Local: "_test_key"},
								Content: "testValue",
							},
						},
					},
				},
			},
		),
		gen(
			"ArrayData",
			nil,
			nil,
			&condition{
				data: map[string]any{
					"soap_Body": []any{
						"value1",
						"value2",
						"value3",
					},
				},
			},
			&action{
				expected: []xmlElement{
					{
						XMLName: xml.Name{Space: "soap", Local: "Body"},
						Content: "value1",
					},
					{
						XMLName: xml.Name{Space: "soap", Local: "Body"},
						Content: "value2",
					},
					{
						XMLName: xml.Name{Space: "soap", Local: "Body"},
						Content: "value3",
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
				attributeKey: cmp.Or(tt.C().attributeKey, "attrKey"),
				namespaceKey: cmp.Or(tt.C().namespaceKey, "nsKey"),
				textKey:      "textKey",
				arrayKey:     "arrayKey",

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

func TestSOAPREST_CreateXMLElementFromValue(t *testing.T) {
	type condition struct {
		elementName string
		value       any
		namespace   string

		attributeKey string
		textKey      string
		namespaceKey string
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
			"NilElement",
			nil,
			nil,
			&condition{},
			&action{
				expected: xmlElement{
					isNil: true,
				},
			},
		),
		gen(
			"ElementWithAttribute",
			nil,
			nil,
			&condition{
				value: map[string]any{
					"attrKey": map[string]any{
						"testAttr": "exampleAttribute",
					},
				},
			},
			&action{
				expected: xmlElement{
					Attrs: []xml.Attr{
						{
							Name:  xml.Name{Local: "testAttr"},
							Value: "exampleAttribute",
						},
					},
				},
			},
		),
		gen(
			"ElementWithNilText",
			nil,
			nil,
			&condition{
				value: map[string]any{
					"textKey": nil,
				},
			},
			&action{
				expected: xmlElement{
					isNil: true,
				},
			},
		),
		gen(
			"ElementWithText",
			nil,
			nil,
			&condition{
				value: map[string]any{
					"textKey": "textElement",
				},
			},
			&action{
				expected: xmlElement{
					Content: "textElement",
				},
			},
		),
		gen(
			"ElementWithChildContent",
			nil,
			nil,
			&condition{
				value: map[string]any{
					"elementNode": map[string]any{
						"childElementNode": map[string]any{
							"childKey": "childValue",
						},
					},
				},
			},
			&action{
				expected: xmlElement{
					children: []xmlElement{
						{
							XMLName: xml.Name{Local: "elementNode"},
							children: []xmlElement{
								{
									XMLName: xml.Name{Local: "childElementNode"},
									children: []xmlElement{
										{
											XMLName: xml.Name{Local: "childKey"},
											Content: "childValue",
										},
									},
								},
							},
						},
					},
				},
			},
		),
		gen(
			"ElementWithSeparatocrChildKey",
			nil,
			nil,
			&condition{
				value: map[string]any{
					"_test_Key": map[string]any{
						"testKey": "testValue",
					},
				},
			},
			&action{
				expected: xmlElement{
					children: []xmlElement{
						{
							XMLName: xml.Name{Local: "_test_Key"},
							children: []xmlElement{
								{
									XMLName: xml.Name{Local: "testKey"},
									Content: "testValue",
								},
							},
						},
					},
				},
			},
		),
		gen(
			"DefaultContent",
			nil,
			nil,
			&condition{
				value: "default",
			},
			&action{
				expected: xmlElement{
					Content: "default",
				},
			},
		),
		gen(
			"SanitizeContent",
			nil,
			nil,
			&condition{
				value: "\b\f\x01\x00",
			},
			&action{
				expected: xmlElement{
					Content: "",
				},
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			s := soapREST{
				attributeKey: cmp.Or(tt.C().attributeKey, "attrKey"),
				namespaceKey: cmp.Or(tt.C().namespaceKey, "nsKey"),
				textKey:      cmp.Or(tt.C().textKey, "textKey"),
				arrayKey:     "arrayKey",

				extractStringElement:  true,
				extractBooleanElement: true,
				extractIntegerElement: true,
				extractFloatElement:   true,
			}

			result := s.createXMLElementFromValue(tt.C().elementName, tt.C().value, tt.C().namespace)
			testutil.Diff(t, tt.A().expected, result, testutil.DeepAllowUnexported(xmlElement{}))
		})
	}
}

func TestSOAPREST_MapToXMLElement(t *testing.T) {
	type condition struct {
		elementName string
		value       any
		namespace   string
		parts       []string

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
			"StringValue",
			nil,
			nil,
			&condition{
				elementName: "TestKey",
				value:       "TestValue",
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
				elementName: "Integer",
				value:       42,
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
				elementName: "Float",
				value:       3.14,
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
				elementName: "Float",
				value:       100.0,
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
				elementName: "OptionalElement",
				value:       nil,
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
				elementName: "EmptyElement",
				value:       map[string]any{},
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
				elementName: "test",
				value: map[string]any{
					"attrKey": map[string]any{
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
				elementName: "test",
				value: map[string]any{
					"textKey": "textContent",
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
				elementName: "test",
				value: map[string]any{
					"textKey": "\b\f",
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
			"NilTextContent",
			nil,
			nil,
			&condition{
				elementName: "test",
				value: map[string]any{
					"textKey": nil,
				},
			},
			&action{
				expected: xmlElement{
					XMLName: xml.Name{Local: "test"},
					isNil:   true,
				},
			},
		),
		gen(
			"MapValueWithChildElements",
			nil,
			nil,
			&condition{
				elementName: "test",
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
			"SeparatorCharKey",
			nil,
			nil,
			&condition{
				elementName: "test",
				value: map[string]any{
					"childElements": map[string]any{
						"_test_Key": []any{
							1,
							2,
							3,
						},
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
									XMLName: xml.Name{Local: "_test_Key"},
									Content: "1",
								},
								{
									XMLName: xml.Name{Local: "_test_Key"},
									Content: "2",
								},
								{
									XMLName: xml.Name{Local: "_test_Key"},
									Content: "3",
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
				elementName: "item",
				value:       []any{"item1", "item2", "item3"},
			},
			&action{
				expected: xmlElement{
					XMLName: xml.Name{Local: "item"},
					Content: "[item1 item2 item3]",
				},
			},
		),
		gen(
			"EmptyArrayValue",
			nil,
			nil,
			&condition{
				elementName: "item",
				value:       []any{},
			},
			&action{
				expected: xmlElement{
					XMLName: xml.Name{Local: "item"},
					Content: "[]",
				},
			},
		),
		gen(
			"StartsWithSeparatorChar",
			nil,
			nil,
			&condition{
				elementName: "_elementName",
				value:       []any{},
				parts:       []string{"_", "item"},
			},
			&action{
				expected: xmlElement{
					XMLName: xml.Name{Space: "_", Local: "item"},
					Content: "[]",
				},
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			s := soapREST{
				attributeKey: cmp.Or(tt.C().attributeKey, "attrKey"),
				textKey:      cmp.Or(tt.C().textKey, "textKey"),
				namespaceKey: "nsKey",
				arrayKey:     cmp.Or(tt.C().arrayKey, "arrayKey"),

				extractStringElement:  true,
				extractBooleanElement: true,
				extractIntegerElement: true,
				extractFloatElement:   true,
			}

			got := s.mapToXMLElement(tt.C().elementName, tt.C().value, tt.C().namespace, tt.C().parts)

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

func TestWrappedWriter_Flush(t *testing.T) {
	type condition struct{}
	type action struct{}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"no-op test",
			nil,
			nil,
			&condition{},
			&action{},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			ww := wrappedWriter{}
			ww.Flush()
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

func TestMapToXMLAttrs(t *testing.T) {
	type condition struct {
		attrMap map[string]any
	}
	type action struct {
		expect []xml.Attr
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"empty map",
			nil,
			nil,
			&condition{
				attrMap: map[string]any{},
			},
			&action{
				expect: []xml.Attr{},
			},
		),
		gen(
			"map with string value",
			nil,
			nil,
			&condition{
				attrMap: map[string]any{
					"testKey": "testValue",
				},
			},
			&action{
				expect: []xml.Attr{
					{
						Name:  xml.Name{Local: "testKey"},
						Value: "testValue",
					},
				},
			},
		),
		gen(
			"map with boolean value",
			nil,
			nil,
			&condition{
				attrMap: map[string]any{
					"enabled": true,
				},
			},
			&action{
				expect: []xml.Attr{
					{
						Name:  xml.Name{Local: "enabled"},
						Value: "true",
					},
				},
			},
		),
		gen(
			"map with integer value",
			nil,
			nil,
			&condition{
				attrMap: map[string]any{
					"count": 42,
				},
			},
			&action{
				expect: []xml.Attr{
					{
						Name:  xml.Name{Local: "count"},
						Value: "42",
					},
				},
			},
		),
		gen(
			"map with float value",
			nil,
			nil,
			&condition{
				attrMap: map[string]any{
					"price": 19.99,
				},
			},
			&action{
				expect: []xml.Attr{
					{
						Name:  xml.Name{Local: "price"},
						Value: "19.99",
					},
				},
			},
		),
		gen(
			"map with nil value",
			nil,
			nil,
			&condition{
				attrMap: map[string]any{
					"data": nil,
				},
			},
			&action{
				expect: []xml.Attr{
					{
						Name:  xml.Name{Local: "data"},
						Value: "<nil>",
					},
				},
			},
		),
		gen(
			"map with multiple values of different types",
			nil,
			nil,
			&condition{
				attrMap: map[string]any{
					"id":      "abc123",
					"count":   100,
					"enabled": false,
					"empty":   nil,
				},
			},
			&action{
				expect: []xml.Attr{
					{
						Name:  xml.Name{Local: "id"},
						Value: "abc123",
					},
					{
						Name:  xml.Name{Local: "count"},
						Value: "100",
					},
					{
						Name:  xml.Name{Local: "enabled"},
						Value: "false",
					},
					{
						Name:  xml.Name{Local: "empty"},
						Value: "<nil>",
					},
				},
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			result := mapToXMLAttrs(tt.C().attrMap)

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
			}

			testutil.Diff(t, tt.A().expect, result, opts...)
		})
	}
}
