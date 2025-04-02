package soaprest

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/aileron-gateway/aileron-gateway/app"
	"github.com/aileron-gateway/aileron-gateway/core"
	"github.com/aileron-gateway/aileron-gateway/kernel/txtutil"
	utilhttp "github.com/aileron-gateway/aileron-gateway/util/http"
)

const (
	soapActionHeaderKey = "SOAPAction"
	soapEnvelopeKey     = "Envelope"
	soapHeaderKey       = "Header"
	soapBodyKey         = "Body"

	separatorChar = "_"

	xsiNamespaceKey = "xsi"
	xmlNamespaceKey = "xmlns"

	soapNamespaceURI = "http://schemas.xmlsoap.org/soap/envelope/"
	xmlNamespaceURI  = "http://www.w3.org/2000/xmlns/"
	xsiNamespaceURI  = "http://www.w3.org/2001/XMLSchema-instance"

	soap11MIMEType = "text/xml"
)

type soapREST struct {
	eh core.ErrorHandler

	// paths is the path matcher to apply SOAP/REST conversion.
	// paths must not be nil.
	paths txtutil.Matcher[string]

	attributeKey string
	textKey      string
	namespaceKey string

	soapNamespacePrefix string

	extractStringElement  bool
	extractBooleanElement bool
	extractIntegerElement bool
	extractFloatElement   bool
}

func (s *soapREST) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// If the request path does not match the configured value,
		// the conversion process will not be executed, and the request will be passed to the next handler.
		if !s.paths.Match(r.URL.Path) {
			next.ServeHTTP(w, r)
			return
		}

		// If it's not a SOAP　1.1 request then return VersionMismatch faultcode.
		if !isSOAPRequest(r) {
			err := app.ErrAppMiddleSOAPRESTVersionMismatch.WithoutStack(nil, nil)
			s.eh.ServeHTTPError(w, r, utilhttp.NewHTTPError(err, http.StatusForbidden))
			return
		}

		// Read the request body
		body, err := io.ReadAll(r.Body)
		if err != nil {
			err = app.ErrAppMiddleSOAPRESTReadRequestBody.WithoutStack(err, map[string]any{"body": body})
			s.eh.ServeHTTPError(w, r, utilhttp.NewHTTPError(err, http.StatusBadRequest))
			return
		}
		r.Body.Close()

		// Parse XML
		var root xmlNode
		if err := xml.Unmarshal(body, &root); err != nil {
			err = app.ErrAppMiddleSOAPRESTUnmarshalRequestBody.WithoutStack(err, map[string]any{"body": body})
			s.eh.ServeHTTPError(w, r, utilhttp.NewHTTPError(err, http.StatusBadRequest))
			return
		}

		// Convert to xmlNode map
		nsCtx := &namespaceContext{
			prefixToURI: map[string]string{},
			uriToPrefix: map[string]string{},
		}
		jsonData := s.xmlToMap(root, nsCtx)

		// Convert the map to JSON bytes
		jsonBody, err := json.Marshal(jsonData)
		if err != nil {
			err = app.ErrAppMiddleSOAPRESTMarshalJSONData.WithoutStack(err, map[string]any{"jsonData": jsonData})
			s.eh.ServeHTTPError(w, r, utilhttp.NewHTTPError(err, http.StatusBadRequest))
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

		// Call the next handler with the modified request
		next.ServeHTTP(ww, newReq)

		// Delete Content-Length because the body will be modified
		ww.ResponseWriter.Header().Del("Content-Length")

		// Convert REST response to SOAP response
		respBody, err := s.convertRESTtoSOAPResponse(ww)
		if err != nil {
			err = app.ErrAppMiddleSOAPRESTDecodeResponseBody.WithoutStack(err, map[string]any{"body": "failed to decode: " + ww.body.String()})
			s.eh.ServeHTTPError(w, r, utilhttp.NewHTTPError(err, http.StatusInternalServerError))
			return
		}

		ww.ResponseWriter.Header().Set("Content-Type", soap11MIMEType+"; charset=utf-8")
		w.WriteHeader(ww.StatusCode())
		_, err = ww.ResponseWriter.Write(respBody)
		if err != nil {
			err = app.ErrAppMiddleSOAPRESTWriteResponseBody.WithoutStack(err, map[string]any{"body": respBody})
			s.eh.ServeHTTPError(w, r, utilhttp.NewHTTPError(err, http.StatusInternalServerError))
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

func removeNewlinesAndTabs(content string) string {
	replacer := strings.NewReplacer(
		"\n", "",
		"\t", "",
		"\r", "",
	)
	return replacer.Replace(content)
}

func (s soapREST) xmlToMap(node xmlNode, nsCtx *namespaceContext) any {
	namespaces := map[string]string{}
	attributes := map[string]string{}

	// Check whether the XML contains `xsi:nil`.
	for _, attr := range node.Attrs {
		if (attr.Name.Space == xsiNamespaceURI || attr.Name.Space == xsiNamespaceKey) &&
			attr.Name.Local == "nil" && strings.ToLower(attr.Value) == "true" {
			return nil
		}
	}

	// Checks whether the namespace URI is included.
	for _, attr := range node.Attrs {
		isXMLNS := (attr.Name.Space == xmlNamespaceKey ||
			attr.Name.Space == xmlNamespaceURI ||
			attr.Name.Local == xmlNamespaceKey)

		if isXMLNS {
			if attr.Value == soapNamespaceURI {
				namespaces[s.soapNamespacePrefix] = attr.Value
				nsCtx.addNamespace(s.soapNamespacePrefix, attr.Value)
			} else {
				namespaces[attr.Name.Local] = attr.Value
				nsCtx.addNamespace(attr.Name.Local, attr.Value)
			}
		} else if !strings.HasPrefix(attr.Name.Space, xmlNamespaceKey) {
			// Obtain the prefix corresponding to the namespace and construct the key
			prefix := nsCtx.getPrefix(attr.Name.Space)
			if prefix != "" {
				attributes[prefix+separatorChar+attr.Name.Local] = attr.Value
			} else {
				attributes[attr.Name.Local] = attr.Value
			}
		}
	}

	content := node.Content
	cleanedContent := removeNewlinesAndTabs(content)
	trimmed := strings.TrimSpace(content)

	resultMap := map[string]any{}
	if len(namespaces) > 0 {
		resultMap[s.namespaceKey] = namespaces
	}
	if len(attributes) > 0 {
		resultMap[s.attributeKey] = attributes
	}

	// Processing child elements
	if len(node.Children) > 0 {
		childrenMap := make(map[string][]any)

		for _, child := range node.Children {
			childValue := s.xmlToMap(child, nsCtx)
			childName := s.getNodeName(child, nsCtx)

			if childMap, ok := childValue.(map[string]any); ok {
				if len(childMap) == 1 {
					var singleKey string
					var singleVal any
					for k, v := range childMap {
						singleKey = k
						singleVal = v
					}

					if singleKey == childName {
						childrenMap[childName] = append(childrenMap[childName], singleVal)
						continue
					}
				}

				childrenMap[childName] = append(childrenMap[childName], childMap)
			} else {
				childrenMap[childName] = append(childrenMap[childName], childValue)
			}
		}

		// Store the processed child elements in `resultMap`
		for k, v := range childrenMap {
			if len(v) == 1 {
				resultMap[k] = v[0]
			} else {
				resultMap[k] = v
			}
		}

		if trimmed != "" {
			resultMap[s.textKey] = s.parseValue(cleanedContent)
		}
	} else {
		if content == "" {
			if len(resultMap) > 0 {
				return resultMap
			}
			return map[string]any{}
		} else {
			if len(resultMap) > 0 {
				resultMap[s.textKey] = s.parseValue(cleanedContent)
				return resultMap
			}
			return s.parseValue(cleanedContent)
		}
	}

	// Only supports conversions for SOAP 1.1
	if node.XMLName.Local == soapEnvelopeKey && node.XMLName.Space == soapNamespaceURI {
		return map[string]any{s.soapNamespacePrefix + separatorChar + soapEnvelopeKey: resultMap}
	}

	// Retrieve element names that include namespace prefixes
	nodeName := s.getNodeName(node, nsCtx)
	return map[string]any{nodeName: resultMap}
}

func (s soapREST) getNodeName(node xmlNode, nsCtx *namespaceContext) string {
	nodeName := node.XMLName.Local

	// Only supports conversions for SOAP 1.1
	if nodeName == soapBodyKey && node.XMLName.Space == soapNamespaceURI {
		// return like "soap_body"
		return s.soapNamespacePrefix + separatorChar + soapBodyKey
	} else if nodeName == soapHeaderKey && node.XMLName.Space == soapNamespaceURI {
		// return like "soap_header"
		return s.soapNamespacePrefix + separatorChar + soapHeaderKey
	}

	if node.XMLName.Space != "" {
		prefix := nsCtx.getPrefix(node.XMLName.Space)
		if prefix != "" {
			return prefix + separatorChar + nodeName
		}
	}

	return nodeName
}

func (s soapREST) convertRESTtoSOAPResponse(wrapper *wrappedWriter) ([]byte, error) {
	decoder := json.NewDecoder(wrapper.body)
	decoder.UseNumber()

	var restData map[string]any
	if err := decoder.Decode(&restData); err != nil {
		return nil, err
	}

	nsManager := &namespaceManager{
		namespaces: make(map[string]string),
	}
	envelope := s.createSOAPEnvelope(restData, nsManager)
	output, _ := xml.MarshalIndent(envelope, "", "  ")

	respBytes := append([]byte(xml.Header), output...)
	return respBytes, nil
}

// soapEnvelope is a struct representing a SOAPEnvelope.
type soapEnvelope struct {
	XMLName xml.Name    `xml:"Envelope"`
	ExtraNS []xml.Attr  `xml:",attr"`
	Attrs   []xml.Attr  `xml:",any,attr,omitempty"`
	Header  *soapHeader `xml:",omitempty"`
	Body    *soapBody

	prefix string
}

// soapHeader is a struct representing a SOAPHeader.
type soapHeader struct {
	XMLName xml.Name   `xml:"Header"`
	Attrs   []xml.Attr `xml:",any,attr,omitempty"`
	Content []xmlElement

	prefix string
}

// soapBody is a struct representing a SOAPBody.
type soapBody struct {
	XMLName xml.Name   `xml:"Body"`
	Attrs   []xml.Attr `xml:",any,attr,omitempty"`
	Content []xmlElement

	prefix string
}

func (e soapEnvelope) MarshalXML(enc *xml.Encoder, start xml.StartElement) error {
	start.Name.Local = fmt.Sprintf("%s:%s", e.prefix, "Envelope")
	start.Attr = append(start.Attr, e.ExtraNS...)

	if len(e.Attrs) > 0 {
		start.Attr = append(start.Attr, e.Attrs...)
	}

	enc.EncodeToken(start)

	if e.Header != nil {
		e.Header.prefix = e.prefix
		enc.Encode(e.Header)
	}

	if e.Body != nil {
		e.Body.prefix = e.prefix
		enc.Encode(e.Body)
	}

	enc.EncodeToken(xml.EndElement{Name: start.Name})
	return nil
}

func (h soapHeader) MarshalXML(enc *xml.Encoder, start xml.StartElement) error {
	start.Name.Local = fmt.Sprintf("%s:%s", h.prefix, "Header")

	if len(h.Attrs) > 0 {
		start.Attr = append(start.Attr, h.Attrs...)
	}

	enc.EncodeToken(start)

	for _, element := range h.Content {
		enc.Encode(element)
	}

	enc.EncodeToken(xml.EndElement{Name: start.Name})
	return nil
}

func (b soapBody) MarshalXML(enc *xml.Encoder, start xml.StartElement) error {
	start.Name.Local = fmt.Sprintf("%s:%s", b.prefix, "Body")

	if len(b.Attrs) > 0 {
		start.Attr = append(start.Attr, b.Attrs...)
	}

	enc.EncodeToken(start)

	for _, element := range b.Content {
		enc.Encode(element)
	}

	enc.EncodeToken(xml.EndElement{Name: start.Name})
	return nil
}

// xmlElement is a struct used for marshaling into XML
type xmlElement struct {
	XMLName     xml.Name
	Attrs       []xml.Attr `xml:",attr"`
	Content     string     `xml:",chardata"`
	children    []xmlElement
	isNil       bool
	hasSiblings bool
}

// xmlElement.MarshalXML is a custom marshaller for encoding an xmlElement struct to XML.
func (e xmlElement) MarshalXML(enc *xml.Encoder, start xml.StartElement) error {
	if e.XMLName.Space != "" {
		start.Name.Local = fmt.Sprintf("%s:%s", e.XMLName.Space, e.XMLName.Local)
	} else {
		start.Name.Local = e.XMLName.Local
	}

	if e.XMLName.Local == "" {
		if e.Content != "" {
			content := e.Content

			// Determine the format based on whether there are sibling elements
			if e.hasSiblings {
				if !strings.Contains(content, "\n") {
					content = "\n    " + content
				}
				// If there are sibling elements, output the text without indentation
				enc.EncodeToken(xml.CharData([]byte(content)))
			} else {
				// If there are no sibling elements, add indentation
				if !strings.Contains(content, "\n") {
					content = "\n    " + content + "\n  "
				}
				enc.EncodeToken(xml.CharData([]byte(content)))
			}
		}
		return nil
	}
	start.Attr = e.Attrs

	if e.isNil {
		start.Attr = append(start.Attr, xml.Attr{
			Name:  xml.Name{Local: "xsi:nil"},
			Value: "true",
		})
	}

	enc.EncodeToken(start)

	// EncodeToken does not perform error handling when the Token is CharData.
	if e.Content != "" {
		enc.EncodeToken(xml.CharData([]byte(e.Content)))
	}

	for _, child := range e.children {
		enc.Encode(child)
	}

	// EncodeToken raises an error if the Token does not match the StartToken or if the Local is an empty string,
	// but no error occurs during the EndToken check in this scenario; therefore, no error handling is performed.
	enc.EncodeToken(xml.EndElement{Name: start.Name})
	return nil
}

func (s soapREST) createSOAPEnvelope(data map[string]any, nsManager *namespaceManager) *soapEnvelope {
	envelope := &soapEnvelope{
		prefix: s.soapNamespacePrefix,
		Header: &soapHeader{
			prefix: s.soapNamespacePrefix,
		},
		Body: &soapBody{
			prefix: s.soapNamespacePrefix,
		},
	}

	if envelopeData, ok := data[s.soapNamespacePrefix+separatorChar+soapEnvelopeKey].(map[string]any); ok {
		if hasNullValue(envelopeData) {
			nsManager.addNamespace(xsiNamespaceKey, xsiNamespaceURI)
		}

		if attrMap, ok := envelopeData[s.attributeKey].(map[string]any); ok {
			envelope.Attrs = mapToXMLAttrs(attrMap, nsManager)
		}

		if nsMap, ok := envelopeData[s.namespaceKey].(map[string]any); ok {
			for prefix, uri := range nsMap {
				if uriStr, isStr := uri.(string); isStr {
					nsManager.addNamespace(prefix, uriStr)
				}
			}
		}

		var nsAttrs []xml.Attr
		for prefix, uri := range nsManager.namespaces {
			var attrName xml.Name
			if prefix == "" {
				attrName = xml.Name{Local: xmlNamespaceKey}
			} else {
				attrName = xml.Name{Local: xmlNamespaceKey + ":" + prefix}
			}
			nsAttrs = append(nsAttrs, xml.Attr{
				Name:  attrName,
				Value: uri,
			})
		}
		envelope.ExtraNS = nsAttrs

		for key, value := range envelopeData {
			valueMap, ok := value.(map[string]any)
			if !ok {
				continue
			}
			parts := strings.SplitN(key, separatorChar, 2)
			elementName := parts[len(parts)-1]

			switch elementName {
			case soapHeaderKey:
				if headerAttrMap, ok := valueMap[s.attributeKey].(map[string]any); ok {
					envelope.Header.Attrs = mapToXMLAttrs(headerAttrMap, nsManager)
				}

				if textContent, ok := valueMap[s.textKey].(string); ok {
					// Check if there are any other child elements
					hasOtherElements := false
					for k := range valueMap {
						if k != s.attributeKey && k != s.textKey && k != s.namespaceKey {
							hasOtherElements = true
							break
						}
					}

					element := xmlElement{
						XMLName:     xml.Name{Local: ""},
						Content:     textContent,
						hasSiblings: hasOtherElements,
					}
					envelope.Header.Content = append(envelope.Header.Content, element)
				}

				childElements := s.mapToXMLElements(valueMap, nsManager)
				envelope.Header.Content = append(envelope.Header.Content, childElements...)

			case soapBodyKey:
				if bodyAttrMap, ok := valueMap[s.attributeKey].(map[string]any); ok {
					envelope.Body.Attrs = mapToXMLAttrs(bodyAttrMap, nsManager)
				}

				if textContent, ok := valueMap[s.textKey].(string); ok {
					// Check if there are any other child elements
					hasOtherElements := false
					for k := range valueMap {
						if k != s.attributeKey && k != s.textKey && k != s.namespaceKey {
							hasOtherElements = true
							break
						}
					}

					element := xmlElement{
						XMLName:     xml.Name{Local: ""},
						Content:     textContent,
						hasSiblings: hasOtherElements,
					}
					envelope.Body.Content = append(envelope.Body.Content, element)
				}

				childElements := s.mapToXMLElements(valueMap, nsManager)
				envelope.Body.Content = append(envelope.Body.Content, childElements...)
			}
		}
	}

	return envelope
}

func (s soapREST) mapToXMLElements(data map[string]any, nsManager *namespaceManager) []xmlElement {
	elements := make([]xmlElement, 0, len(data))
	for key, value := range data {
		// Keys that include attributeKey and textKey have already been processed.
		if key == s.attributeKey || key == s.namespaceKey || key == s.textKey {
			continue
		}

		// If the value is an array, handle each element separately
		if valueArray, isArray := value.([]any); isArray {
			// In the case of arrays, add each element directly under the same name without creating a parent element
			for _, item := range valueArray {
				// Split the key into a namespace prefix and a local name
				var namespace string
				elementName := key

				// Check if the key starts with `separatorChar` at the beginning
				startsWithSeparator := strings.HasPrefix(key, separatorChar)

				parts := strings.SplitN(key, separatorChar, 2)
				if len(parts) == 2 && !startsWithSeparator {
					namespace = parts[0]
					elementName = parts[1]
				}

				element := s.createXMLElementFromValue(elementName, item, namespace, nsManager)
				elements = append(elements, element)
			}
			continue
		}

		// Split the key into a namespace prefix and a local name
		startsWithSeparatorChar := strings.HasPrefix(key, separatorChar)
		parts := strings.SplitN(key, separatorChar, 2)
		var namespace string

		if len(parts) == 2 {
			namespace = parts[0]

			// If the key starts with separatorChar, prepend separatorChar to elementName as well
			if startsWithSeparatorChar {
				parts[1] = separatorChar + parts[1]
			}
		}

		element := s.mapToXMLElement(key, value, namespace, parts, nsManager)
		elements = append(elements, element)
	}
	return elements
}

func (s soapREST) createXMLElementFromValue(elementName string, value any, namespace string, nsManager *namespaceManager) xmlElement {
	element := xmlElement{
		XMLName: xml.Name{
			Space: namespace,
			Local: elementName,
		},
	}

	switch v := value.(type) {
	case nil:
		element.isNil = true
	case map[string]any:
		// Process attributes
		if attrMap, ok := v[s.attributeKey].(map[string]any); ok {
			element.Attrs = mapToXMLAttrs(attrMap, nsManager)
		}

		// Process text content
		if textValue, ok := v[s.textKey]; ok {
			if textValue == nil {
				element.isNil = true
			} else if textContent, isStr := textValue.(string); isStr {
				element.Content = sanitizeControlCharacters(textContent)
			}
		}

		// Process child elements
		for childKey, childValue := range v {
			if childKey == s.attributeKey || childKey == s.textKey || childKey == s.namespaceKey {
				continue
			}

			childParts := strings.SplitN(childKey, separatorChar, 2)
			var childNamespace string
			var childLocalName string

			// If the key starts with separatorChar, prepend separatorChar to childLocalName as well
			if strings.HasPrefix(childKey, separatorChar) {
				childParts[1] = separatorChar + childParts[1]
			}

			if len(childParts) == 2 {
				childNamespace = childParts[0]
				childLocalName = childParts[1]
			} else {
				childLocalName = childKey
			}

			child := s.mapToXMLElement(childLocalName, childValue, childNamespace, childParts, nsManager)
			element.children = append(element.children, child)
		}
	default:
		element.Content = sanitizeControlCharacters(fmt.Sprintf("%v", v))
	}

	return element
}

func (s soapREST) mapToXMLElement(elementName string, value any, namespace string, parts []string, nsManager *namespaceManager) xmlElement {
	// When the JSON data contains an array, the key does not include the separator character
	// so the length of `parts` will not be 2
	if len(parts) == 2 {
		namespace = parts[0]
		elementName = parts[1]
	}

	// Create the basic structure of an xmlElement
	element := xmlElement{
		XMLName: xml.Name{
			Space: namespace,
			Local: elementName,
		},
	}

	// Perform appropriate processing based on the type of the value
	switch v := value.(type) {
	case nil:
		element.isNil = true

	case map[string]any:
		// When a value is not null but is empty
		if len(v) == 0 {
			return element
		}

		// Processing attributes
		if attrMap, ok := v[s.attributeKey].(map[string]any); ok {
			element.Attrs = mapToXMLAttrs(attrMap, nsManager)
		}

		// Processing text content
		if textValue, ok := v[s.textKey]; ok {
			if textValue == nil {
				element.isNil = true
			} else if textContent, isStr := textValue.(string); isStr {
				element.Content = sanitizeControlCharacters(textContent)
			}
		}

		// Processing child elements
		for childKey, childValue := range v {
			if childKey == s.attributeKey || childKey == s.textKey || childKey == s.namespaceKey {
				continue
			}

			// Check whether it is an array
			if childArray, isArray := childValue.([]any); isArray {
				// In the case of an array, add each element as an independent child element
				for _, item := range childArray {
					childParts := strings.SplitN(childKey, separatorChar, 2)
					var childNamespace string
					var childLocalName string

					startsWithSeparator := strings.HasPrefix(childKey, separatorChar)

					if len(childParts) == 2 {
						childNamespace = childParts[0]
						if startsWithSeparator {
							childParts[1] = separatorChar + childParts[1]
						}
						childLocalName = childParts[1]
					} else {
						childLocalName = childKey
					}

					childElement := s.createXMLElementFromValue(childLocalName, item, childNamespace, nsManager)
					element.children = append(element.children, childElement)
				}
			} else {
				childParts := strings.SplitN(childKey, separatorChar, 2)
				var childNamespace string
				var childLocalName string

				startsWithSeparator := strings.HasPrefix(childKey, separatorChar)

				if len(childParts) == 2 {
					childNamespace = childParts[0]
					if startsWithSeparator {
						childParts[1] = separatorChar + childParts[1]
					}
					childLocalName = childParts[1]
				} else {
					childLocalName = childKey
				}

				child := s.mapToXMLElement(childLocalName, childValue, childNamespace, childParts, nsManager)
				element.children = append(element.children, child)
			}
		}

	default:
		element.Content = sanitizeControlCharacters(fmt.Sprintf("%v", v))
	}

	return element
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

func (w *wrappedWriter) Flush() {
	// no-op
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
	if s.extractStringElement {
		if strings.HasPrefix(content, `"`) && strings.HasSuffix(content, `"`) {
			return content[1 : len(content)-1]
		}
	}

	if s.extractBooleanElement {
		trimmedContent := strings.TrimSpace(content)
		if trimmedContent == "true" {
			return true
		}

		if trimmedContent == "false" {
			return false
		}
	}

	trimmedContent := strings.TrimSpace(content)

	// If the conversion to integer fails, it is handled as a string without any error handling.
	if s.extractIntegerElement {
		if i, err := strconv.ParseInt(trimmedContent, 10, 64); err == nil {
			return i
		}
	}

	// If the conversion to float fails, it is handled as a string without any error handling.
	if s.extractFloatElement {
		if f, err := strconv.ParseFloat(trimmedContent, 64); err == nil {
			return f
		}
	}

	return content
}

// sanitizeXMLCharacters removes invalid XML 1.0 characters from the input string.
func sanitizeControlCharacters(input string) string {
	var sanitized strings.Builder
	for _, r := range input {
		if isValidXMLChar(r) {
			sanitized.WriteRune(r)
		}
	}
	return sanitized.String()
}

// isValidXMLChar checks if a rune is valid according to XML 1.0.
func isValidXMLChar(r rune) bool {
	// Valid XML characters (according to XML 1.0 spec)
	return (r == 0x09 || r == 0x0A || r == 0x0D || // Tab, Line Feed, Carriage Return
		(r >= 0x20 && r <= 0xD7FF) || // Basic characters
		(r >= 0xE000 && r <= 0xFFFD) || // Valid non-supplementary characters
		(r >= 0x10000 && r <= 0x10FFFF)) // Supplementary characters
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

// mapToXMLAttrs is a helper function that converts attributes (key-value pairs) in JSON to `xml.Attr`
func mapToXMLAttrs(attrMap map[string]interface{}, nsManager *namespaceManager) []xml.Attr {
	attrs := make([]xml.Attr, 0, len(attrMap))
	for k, v := range attrMap {
		// Handling of attributes that begin with "xmlns_"
		if strings.HasPrefix(k, "xmlns_") {
			localName := strings.TrimPrefix(k, "xmlns_")
			attrs = append(attrs, xml.Attr{
				Name:  xml.Name{Local: "xmlns:" + localName},
				Value: fmt.Sprintf("%v", v),
			})
			continue
		}

		// Handling of attributes that begin with "xsi_"
		if strings.HasPrefix(k, "xsi_") {
			// "xsi_type" → "xsi:type"
			localName := strings.TrimPrefix(k, "xsi_")
			attrs = append(attrs, xml.Attr{
				Name:  xml.Name{Local: "xsi:" + localName},
				Value: fmt.Sprintf("%v", v),
			})
			continue
		}

		// Handling attributes with common prefixes using separatorChar
		parts := strings.SplitN(k, "_", 2)
		if len(parts) == 2 && nsManager.namespaces[parts[0]] != "" {
			attrs = append(attrs, xml.Attr{
				Name:  xml.Name{Local: parts[0] + ":" + parts[1]},
				Value: fmt.Sprintf("%v", v),
			})
		} else {
			attrs = append(attrs, xml.Attr{
				Name:  xml.Name{Local: k},
				Value: fmt.Sprintf("%v", v),
			})
		}
	}
	return attrs
}
