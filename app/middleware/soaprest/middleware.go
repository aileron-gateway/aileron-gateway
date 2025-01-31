package soaprest

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/aileron-gateway/aileron-gateway/core"
	utilhttp "github.com/aileron-gateway/aileron-gateway/util/http"
)

const (
	soapActionHeaderKey = "SOAPAction"
	soapEnvelopeKey     = "Envelope"
	soapHeaderKey       = "Header"
	soapBodyKey         = "Body"

	soapNameSpaceKey = "soap"
	xsiNameSpaceKey  = "xsi"
	xmlNameSpaceKey  = "xmlns"

	soapNameSpaceURI = "http://schemas.xmlsoap.org/soap/envelope/"
	xsiNameSpaceURI  = "http://www.w3.org/2001/XMLSchema-instance"

	soap11MIMEType = "text/xml"
)

var (
	errInvalidSOAP11Request = utilhttp.NewHTTPError(errors.New("expected a SOAP 1.1 request, but received a request in a different format"), http.StatusForbidden)
)

type soapREST struct {
	eh core.ErrorHandler

	attributeKey  string
	textKey       string
	namespaceKey  string
	arrayKey      string
	separatorChar string

	extractStringElement  bool
	extractBooleanElement bool
	extractIntegerElement bool
	extractFloatElement   bool
}

func (m *soapREST) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// If it's not a SOAP　1.1 request then return VersionMismatch faultcode.
		if !isSOAPRequest(r) {
			m.eh.ServeHTTPError(w, r, errInvalidSOAP11Request)
			return
		}

		// Read the request body
		body, err := io.ReadAll(r.Body)
		if err != nil {
			m.eh.ServeHTTPError(w, r, utilhttp.ErrBadRequest)
			return
		}
		r.Body.Close()

		// Parse XML
		var root xmlNode
		if err := xml.Unmarshal(body, &root); err != nil {
			m.eh.ServeHTTPError(w, r, utilhttp.ErrBadRequest)
			return
		}

		// Convert to xmlNode map
		nsCtx := &namespaceContext{
			prefixToURI: map[string]string{},
			uriToPrefix: map[string]string{},
		}
		jsonData := m.xmlToMap(root, nsCtx)

		// Convert the map to JSON bytes
		jsonBody, err := json.Marshal(jsonData)
		if err != nil {
			m.eh.ServeHTTPError(w, r, utilhttp.ErrBadRequest)
			return
		}

		// Create new request with JSON body
		newReq := r.Clone(r.Context())
		newReq.Body = io.NopCloser(bytes.NewReader(jsonBody))
		newReq.ContentLength = int64(len(jsonBody))
		newReq.Header.Set("Content-Type", "application/json")
		newReq.Method = "POST"

		ww := &wrappedWriter{
			ResponseWriter: w,
			body:           &bytes.Buffer{},
		}
		w = ww

		// Call the next handler with the modified request
		next.ServeHTTP(w, newReq)

		// Convert REST response to SOAP response
		respBody, err := m.convertRestToSoapResponse(ww)
		if err != nil {
			m.eh.ServeHTTPError(w, r, utilhttp.ErrInternalServerError)
			return
		}

		ww.ResponseWriter.Header().Set("Content-Type", soap11MIMEType+"; charset=utf-8")
		ww.ResponseWriter.Header().Set("Content-Length", strconv.Itoa(len(respBody)))

		w.WriteHeader(ww.StatusCode())
		_, err = ww.ResponseWriter.Write(respBody)
		if err != nil {
			m.eh.ServeHTTPError(w, r, utilhttp.ErrInternalServerError)
		}
	})
}

// xmlNode is a struct designed to understand the structure of XML
// and is used for converting SOAP to REST.
type xmlNode struct {
	XMLName  xml.Name
	Content  string     `xml:",chardata"` // text contents
	Attrs    []xml.Attr `xml:",any,attr"` // attributes
	Children []xmlNode  `xml:",any"`      // recursive contents
}

func (s soapREST) xmlToMap(node xmlNode, nsCtx *namespaceContext) any {
	// Check whether the XML contains `xsi:nil`.
	for _, attr := range node.Attrs {
		if (attr.Name.Space == xsiNameSpaceURI || attr.Name.Space == xsiNameSpaceKey) &&
			attr.Name.Local == "nil" && strings.ToLower(attr.Value) == "true" {
			return nil
		}
	}

	namespaces := map[string]string{}
	attributes := map[string]string{}

	// Checks whether the namespace URI is included.
	for _, attr := range node.Attrs {
		if attr.Name.Space == xmlNameSpaceKey || attr.Name.Local == xmlNameSpaceKey {
			prefix := attr.Name.Local
			if attr.Name.Space == xmlNameSpaceKey {
				prefix = attr.Name.Local
			}
			namespaces[prefix] = attr.Value
			nsCtx.addNamespace(prefix, attr.Value)
		} else if !strings.HasPrefix(attr.Name.Space, xmlNameSpaceKey) {
			attributes[attr.Name.Local] = attr.Value
		}
	}

	resultMap := map[string]any{}
	if len(namespaces) > 0 {
		resultMap[s.namespaceKey] = namespaces
	}
	if len(attributes) > 0 {
		resultMap[s.attributeKey] = attributes
	}

	// Processing XML content
	content := strings.TrimSpace(node.Content)
	if len(node.Children) > 0 {
		childrenMap := make(map[string][]any)

		// Determine whether the content treated as an array.
		isArray := false
		if len(node.Children) > 0 {
			firstChild := node.Children[0]
			allSameName := true
			for _, child := range node.Children {
				if child.XMLName.Local != firstChild.XMLName.Local ||
					child.XMLName.Space != firstChild.XMLName.Space {
					allSameName = false
					break
				}
			}
			// Check if the array key matches the config value.
			isArrayElement := firstChild.XMLName.Local == s.arrayKey
			isArray = allSameName && isArrayElement
		}

		// Handling array content
		if isArray {
			values := make([]any, 0, len(node.Children))
			for _, child := range node.Children {
				childValue := s.xmlToMap(child, nsCtx)
				if childMap, ok := childValue.(map[string]any); ok {
					for _, v := range childMap {
						values = append(values, v)
					}
				} else {
					values = append(values, childValue)
				}
			}
			nodeName := s.getNodeName(node, nsCtx)
			return map[string]any{nodeName: values}
		} else {
			// Handling not array content
			for _, child := range node.Children {
				childName := s.getNodeName(child, nsCtx)
				childValue := s.xmlToMap(child, nsCtx)
				if childMap, ok := childValue.(map[string]any); ok {
					if len(childMap) == 1 {
						// Single entry child node
						//
						// <singleEntry>
						//   <value>OnlyValue</value>
						// </singleEntry>
						//
						for _, v := range childMap {
							childrenMap[childName] = append(childrenMap[childName], v)
						}
					} else {
						// Map child node
						//
						// <mapNode>
						//   <entry1>Value1</entry1>
						//   <entry2>Value2</entry2>
						// </mapNode>
						//
						childrenMap[childName] = append(childrenMap[childName], childMap)
					}
				} else {
					// Non map child node
					//
					// <nonMapNode>JustText</nonMapNode>
					//
					childrenMap[childName] = append(childrenMap[childName], childValue)
				}
			}

			for k, v := range childrenMap {
				if len(v) == 1 {
					resultMap[k] = v[0]
				} else {
					resultMap[k] = v
				}
			}
		}
	}

	// text node with children and attributes
	if content != "" && (len(node.Children) > 0 || len(node.Attrs) > 0) {
		resultMap[s.textKey] = s.parseValue(content)
	} else if content != "" {
		// text node without children or attributes
		return s.parseValue(content)
	} else if len(resultMap) == 0 {
		// doesn't contain text node and children
		return map[string]any{}
	}

	// Only supports conversions for SOAP 1.1
	if node.XMLName.Local == soapEnvelopeKey &&
		node.XMLName.Space == soapNameSpaceURI {
		return map[string]any{soapNameSpaceKey + s.separatorChar + soapEnvelopeKey: resultMap}
	}

	nodeName := s.getNodeName(node, nsCtx)
	return map[string]any{nodeName: resultMap}
}

func (s soapREST) getNodeName(node xmlNode, nsCtx *namespaceContext) string {
	nodeName := node.XMLName.Local

	// Only supports conversions for SOAP 1.1
	if nodeName == soapBodyKey && node.XMLName.Space == soapNameSpaceURI {
		// return like "soap:body"
		return soapNameSpaceKey + s.separatorChar + soapBodyKey
	} else if nodeName == soapHeaderKey && node.XMLName.Space == soapNameSpaceURI {
		// return like "soap:header"
		return soapNameSpaceKey + s.separatorChar + soapHeaderKey
	}

	if node.XMLName.Space != "" {
		prefix := nsCtx.getPrefix(node.XMLName.Space)
		if prefix != "" {
			return prefix + s.separatorChar + nodeName
		}
	}

	return nodeName
}

func (m *soapREST) convertRestToSoapResponse(wrapper *wrappedWriter) ([]byte, error) {
	var restData map[string]any
	if err := json.NewDecoder(wrapper.body).Decode(&restData); err != nil {
		return nil, err
	}

	nsManager := &namespaceManager{
		namespaces: make(map[string]string),
	}
	envelope := m.createSOAPEnvelope(restData, nsManager)

	output, err := xml.MarshalIndent(envelope, "", "  ")
	if err != nil {
		return nil, err
	}

	responseBytes := append([]byte(xml.Header), output...)
	return responseBytes, err
}

// soapEnvelope is a struct representing a SOAPEnvelope.
type soapEnvelope struct {
	XMLName xml.Name    `xml:"soap:Envelope"`
	ExtraNS []xml.Attr  `xml:",attr"`
	Header  *soapHeader `xml:"soap:Header,omitempty"`
	Body    *soapBody   `xml:"soap:Body"`
}

// soapHeader is a struct representing a SOAPHeader.
type soapHeader struct {
	XMLName xml.Name `xml:"soap:Header"`
	Content []xmlElement
}

// soapBody is a struct representing a SOAPBody.
type soapBody struct {
	XMLName xml.Name `xml:"soap:Body"`
	Content []xmlElement
}

// xmlElement is a struct used for marshaling into XML
type xmlElement struct {
	XMLName  xml.Name
	Attrs    []xml.Attr `xml:",attr"`
	Content  string     `xml:",chardata"`
	children []xmlElement
	isNil    bool
}

// xmlElement.MarshalXML is a custom marshaller for encoding an xmlElement struct to XML.
func (e xmlElement) MarshalXML(enc *xml.Encoder, start xml.StartElement) error {
	if e.XMLName.Space != "" {
		start.Name.Local = fmt.Sprintf("%s:%s", e.XMLName.Space, e.XMLName.Local)
	} else {
		start.Name.Local = e.XMLName.Local
	}

	start.Attr = e.Attrs

	if e.isNil {
		start.Attr = append(start.Attr, xml.Attr{
			Name:  xml.Name{Local: "xsi:nil"},
			Value: "true",
		})
	}

	if err := enc.EncodeToken(start); err != nil {
		return err
	}

	if e.Content != "" {
		if err := enc.EncodeToken(xml.CharData([]byte(e.Content))); err != nil {
			return err
		}
	}

	for _, child := range e.children {
		if err := enc.Encode(child); err != nil {
			return err
		}
	}

	return enc.EncodeToken(xml.EndElement{Name: start.Name})
}

func (s soapREST) createSOAPEnvelope(data map[string]any, nsManager *namespaceManager) *soapEnvelope {
	envelope := &soapEnvelope{
		Header: &soapHeader{},
		Body:   &soapBody{},
	}

	if envelopeData, ok := data[soapNameSpaceKey+s.separatorChar+soapEnvelopeKey].(map[string]any); ok {
		if hasNullValue(envelopeData) {
			nsManager.addNamespace(xsiNameSpaceKey, xsiNameSpaceURI)
		}

		if nsMap, nsOk := envelopeData[s.namespaceKey].(map[string]any); nsOk {
			for prefix, uri := range nsMap {
				if uriStr, isStr := uri.(string); isStr {
					nsManager.addNamespace(prefix, uriStr)
				}
			}
		}

		var nsAttrs []xml.Attr
		for prefix, uri := range nsManager.namespaces {
			nsAttrs = append(nsAttrs, xml.Attr{
				Name:  xml.Name{Local: xmlNameSpaceKey + s.separatorChar + prefix},
				Value: uri,
			})
		}
		envelope.ExtraNS = nsAttrs

		for key, value := range envelopeData {
			if valueMap, ok := value.(map[string]any); ok {
				parts := strings.SplitN(key, s.separatorChar, 2)
				elementName := parts[len(parts)-1]

				switch elementName {
				case soapHeaderKey:
					envelope.Header.Content = s.mapToXMLElements(valueMap, nsManager)
				case soapBodyKey:
					envelope.Body.Content = s.mapToXMLElements(valueMap, nsManager)
				}
			}
		}
	}

	return envelope
}

func (s soapREST) mapToXMLElements(data map[string]any, nsManager *namespaceManager) []xmlElement {
	elements := make([]xmlElement, 0, len(data))
	for key, value := range data {
		// Keys that include attributeKey and textKey have already been processed.
		if key == s.attributeKey || key == s.namespaceKey {
			continue
		}

		parts := strings.SplitN(key, s.separatorChar, 2)
		var namespace string
		if len(parts) == 2 {
			namespace = parts[0]
		}

		if valueMap, ok := value.(map[string]any); ok {
			if nsMap, nsOk := valueMap[s.namespaceKey].(map[string]any); nsOk {
				for prefix, uri := range nsMap {
					if uriStr, isStr := uri.(string); isStr {
						nsManager.addNamespace(prefix, uriStr)
					}
				}
			}
		}

		element := s.mapToXMLElement(key, value, namespace)
		elements = append(elements, element)
	}
	return elements
}

// mapToXMLElement is a function that converts JSON data into xmlElement.
func (s soapREST) mapToXMLElement(key string, value any, namespace string) xmlElement {
	// Separate the namespace and local name from a key.
	parts := strings.SplitN(key, s.separatorChar, 2)
	elementName := key
	// When the JSON data contains an array, the key does not include the separator character,
	// so the length of `parts` will not be 2.
	if len(parts) == 2 {
		namespace = parts[0]
		elementName = parts[1]
	}

	// Create the basic structure of an xmlElement.
	element := xmlElement{
		XMLName: xml.Name{
			Space: namespace,
			Local: elementName,
		},
	}

	// Perform appropriate processing based on the type of the value.
	switch v := value.(type) {
	case nil:
		element.isNil = true

	case map[string]any:
		// When a value is not null but is empty.
		if len(v) == 0 {
			return element
		}

		element.Attrs = s.extractAttributes(v)

		if textContent, ok := v[s.textKey].(string); ok {
			element.Content = sanitizeControlCharacters(textContent)
		}

		for k, childValue := range v {
			// Keys that include attributeKey and textKey have already been processed.
			if k == s.attributeKey || k == s.textKey {
				continue
			}
			child := s.mapToXMLElement(k, childValue, namespace)
			element.children = append(element.children, child)
		}

	case []any:
		// Process each element of the array as a child element.
		for _, item := range v {
			child := s.mapToXMLElement(s.arrayKey, item, namespace)
			element.children = append(element.children, child)
		}

	case float64:
		element.Content = formatFloat(v)

	default:
		element.Content = sanitizeControlCharacters(fmt.Sprintf("%v", v))
	}

	return element
}

func (s soapREST) extractAttributes(data map[string]any) []xml.Attr {
	var attrs []xml.Attr
	if attrMap, ok := data[s.attributeKey].(map[string]any); ok {
		for key, value := range attrMap {
			attrs = append(attrs, xml.Attr{
				Name:  xml.Name{Local: key},
				Value: fmt.Sprintf("%v", value),
			})
		}
	}
	return attrs
}

// wrappedWriter wraps http.ResponseWriter.
// This implements io.Writer interface and utilhttp.Writer interface.
type wrappedWriter struct {
	http.ResponseWriter
	code    int
	written bool
	length  int64
	body    *bytes.Buffer
}

func (w *wrappedWriter) Unwrap() http.ResponseWriter {
	return w.ResponseWriter
}

func (w *wrappedWriter) WriteHeader(statusCode int) {
	if w.written {
		return
	}
	w.code = statusCode
	w.written = true
}

func (w *wrappedWriter) Write(b []byte) (int, error) {
	w.written = true
	w.length += int64(len(b))
	return w.body.Write(b)
}

func (w *wrappedWriter) Body() *bytes.Buffer {
	return w.body
}

func (w *wrappedWriter) Written() bool {
	return w.written
}

func (w *wrappedWriter) StatusCode() int {
	if w.written && w.code == 0 {
		return http.StatusOK
	}
	return w.code
}

func (w *wrappedWriter) ContentLength() int64 {
	return w.length
}

// namespaceContext is a struct used for managing namespaces.
type namespaceContext struct {
	prefixToURI map[string]string
	uriToPrefix map[string]string
}

func (nc *namespaceContext) addNamespace(prefix, uri string) {
	nc.prefixToURI[prefix] = uri
	nc.uriToPrefix[uri] = prefix
}

func (nc *namespaceContext) getPrefix(uri string) string {
	if prefix, ok := nc.uriToPrefix[uri]; ok {
		return prefix
	}
	return ""
}

type namespaceManager struct {
	namespaces map[string]string
}

func (nm *namespaceManager) addNamespace(prefix, uri string) {
	if _, exists := nm.namespaces[prefix]; !exists {
		nm.namespaces[prefix] = uri
	}
}

// Only supports conversions for SOAP 1.1
// The request is determined to be SOAP 1.1 if the Content-Type includes "text/xml" or
// if the value of the SOAPAction header is not empty.
func isSOAPRequest(r *http.Request) bool {
	return strings.Contains(r.Header.Get("Content-Type"), soap11MIMEType) ||
		r.Header.Get(soapActionHeaderKey) != ""
}

// parseValue extracts XML elements according to their type.
func (s soapREST) parseValue(content string) any {
	content = strings.TrimSpace(content)

	if s.extractStringElement {
		if strings.HasPrefix(content, `"`) && strings.HasSuffix(content, `"`) {
			return content[1 : len(content)-1]
		}
	}

	if s.extractBooleanElement {
		if content == "true" {
			return true
		}

		if content == "false" {
			return false
		}
	}

	if s.extractIntegerElement {
		if i, err := strconv.ParseInt(content, 10, 64); err == nil {
			return i
		}
	}

	if s.extractFloatElement {
		if f, err := strconv.ParseFloat(content, 64); err == nil {
			return f
		}
	}

	return content
}

// formatFloat converts float to string, removing unnecessary trailing zeros and decimal point.
func formatFloat(v float64) string {
	str := strconv.FormatFloat(v, 'f', -1, 64)
	str = strings.TrimRight(strings.TrimRight(str, "0"), ".")
	return str
}

// sanitizeControlCharacters sanitizes characters that are not allowed in XML.
func sanitizeControlCharacters(input string) string {
	replacer := strings.NewReplacer(
		"\b", "", // backspace
		"\f", "", // form feed
		"\\b", "", // escaped backspace
		"\\f", "", // escaped form feed
	)
	return replacer.Replace(input)
}

// Recursively check if the JSON contains null.
// If a JSON element contains null, the xsi namespace will be added to the definition.
func hasNullValue(data any) bool {
	switch v := data.(type) {
	case nil:
		return true
	case map[string]any:
		for _, value := range v {
			if hasNullValue(value) {
				return true
			}
		}
	case []any:
		for _, item := range v {
			if hasNullValue(item) {
				return true
			}
		}
	}
	return false
}
