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

	soapNameSpaceKey = "soap"
	xsiNameSpaceKey  = "xsi"
	xmlNameSpaceKey  = "xmlns"

	soapNameSpaceURI = "http://schemas.xmlsoap.org/soap/envelope/"
	xsiNameSpaceURI  = "http://www.w3.org/2001/XMLSchema-instance"

	soap11MIMEType = "text/xml"
)

type soapREST struct {
	eh core.ErrorHandler

	// paths is the path matcher to apply SOAP/REST conversion.
	// paths must not be nil.
	paths txtutil.Matcher[string]

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

		// Convert REST response to SOAP response
		respBody, err := s.convertRESTtoSOAPResponse(ww)
		if err != nil {
			s.eh.ServeHTTPError(w, r, err)
			return
		}

		ww.ResponseWriter.Header().Set("Content-Type", soap11MIMEType+"; charset=utf-8")
		ww.ResponseWriter.Header().Set("Content-Length", strconv.Itoa(len(respBody)))

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

	// Processing child elements.
	if len(node.Children) > 0 {
		childrenMap := make(map[string][]any)

		for _, child := range node.Children {
			childName := s.getNodeName(child, nsCtx)
			childValue := s.xmlToMap(child, nsCtx)

			if childMap, ok := childValue.(map[string]any); ok {
				if len(childMap) == 1 {
					for k, v := range childMap {
						// Preserve namespace prefixes if they are included
						if strings.Contains(k, s.separatorChar) {
							childrenMap[k] = append(childrenMap[k], v)
						} else {
							childrenMap[childName] = append(childrenMap[childName], v)
						}
					}
				} else {
					childrenMap[childName] = append(childrenMap[childName], childMap)
				}
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
		if content != "" {
			if trimmed == "" {
				// Even if the content consists only of whitespace, retain it if attributes are present
				if len(resultMap) > 0 {
					return resultMap
				}
				return map[string]any{}
			}

			// If attribute information is already stored, retain it
			if len(resultMap) > 0 {
				resultMap[s.textKey] = s.parseValue(cleanedContent)
				return resultMap
			}

			return s.parseValue(cleanedContent)
		} else {
			if len(resultMap) == 0 {
				return map[string]any{}
			}
		}
	}

	// Only supports conversions for SOAP 1.1
	if node.XMLName.Local == soapEnvelopeKey && node.XMLName.Space == soapNameSpaceURI {
		return map[string]any{soapNameSpaceKey + s.separatorChar + soapEnvelopeKey: resultMap}
	}

	// Retrieve element names that include namespace prefixes
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

func (s soapREST) convertRESTtoSOAPResponse(wrapper *wrappedWriter) ([]byte, error) {
	var restData map[string]any
	if err := json.NewDecoder(wrapper.body).Decode(&restData); err != nil {
		err = app.ErrAppMiddleSOAPRESTDecodeResponseBody.WithoutStack(err, map[string]any{"body": "failed to decode: " + wrapper.body.String()})
		return nil, utilhttp.NewHTTPError(err, http.StatusInternalServerError)
	}

	nsManager := &namespaceManager{
		namespaces: make(map[string]string),
	}
	envelope := s.createSOAPEnvelope(restData, nsManager)

	output, err := xml.MarshalIndent(envelope, "", "  ")
	if err != nil {
		return nil, err
	}

	respBytes := append([]byte(xml.Header), output...)
	return respBytes, nil
}

// soapEnvelope is a struct representing a SOAPEnvelope.
type soapEnvelope struct {
	XMLName xml.Name    `xml:"soap:Envelope"`
	ExtraNS []xml.Attr  `xml:",attr"`
	Attrs   []xml.Attr  `xml:",any,attr,omitempty"`
	Header  *soapHeader `xml:"soap:Header,omitempty"`
	Body    *soapBody   `xml:"soap:Body"`
}

// soapHeader is a struct representing a SOAPHeader.
type soapHeader struct {
	XMLName xml.Name   `xml:"soap:Header"`
	Attrs   []xml.Attr `xml:",any,attr,omitempty"`
	Content []xmlElement
}

// soapBody is a struct representing a SOAPBody.
type soapBody struct {
	XMLName xml.Name   `xml:"soap:Body"`
	Attrs   []xml.Attr `xml:",any,attr,omitempty"`
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

	if e.XMLName.Local == "" {
		if e.Content != "" {
			content := e.Content
			if !strings.Contains(content, "\n") {
				content = "\n    " + content
			}
			enc.EncodeToken(xml.CharData([]byte(content)))
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
	return enc.EncodeToken(xml.EndElement{Name: start.Name})
}

// mapToXMLAttrs is a helper function that converts attributes (key-value pairs) in JSON to `xml.Attr`
func mapToXMLAttrs(attrMap map[string]any) []xml.Attr {
	attrs := make([]xml.Attr, 0, len(attrMap))
	for k, v := range attrMap {
		attrs = append(attrs, xml.Attr{
			Name:  xml.Name{Local: k},
			Value: fmt.Sprintf("%v", v),
		})
	}
	return attrs
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

		if attrMap, ok := envelopeData[s.attributeKey].(map[string]any); ok {
			envelope.Attrs = mapToXMLAttrs(attrMap)
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
			nsAttrs = append(nsAttrs, xml.Attr{
				Name:  xml.Name{Local: xmlNameSpaceKey + ":" + prefix},
				Value: uri,
			})
		}
		envelope.ExtraNS = nsAttrs

		for key, value := range envelopeData {
			valueMap, ok := value.(map[string]any)
			if !ok {
				continue
			}
			parts := strings.SplitN(key, s.separatorChar, 2)
			elementName := parts[len(parts)-1]

			switch elementName {
			case soapHeaderKey:
				if headerAttrMap, ok := valueMap[s.attributeKey].(map[string]any); ok {
					envelope.Header.Attrs = mapToXMLAttrs(headerAttrMap)
				}

				if textContent, ok := valueMap[s.textKey].(string); ok {
					element := xmlElement{
						XMLName: xml.Name{Local: ""},
						Content: textContent,
					}
					envelope.Header.Content = append(envelope.Header.Content, element)
				}

				childElements := s.mapToXMLElements(valueMap, nsManager)
				envelope.Header.Content = append(envelope.Header.Content, childElements...)

			case soapBodyKey:
				if bodyAttrMap, ok := valueMap[s.attributeKey].(map[string]any); ok {
					envelope.Body.Attrs = mapToXMLAttrs(bodyAttrMap)
				}

				if textContent, ok := valueMap[s.textKey].(string); ok {
					element := xmlElement{
						XMLName: xml.Name{Local: ""},
						Content: textContent,
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

		// Split the key into a namespace prefix and a local name
		parts := strings.SplitN(key, s.separatorChar, 2)
		var namespace string

		if len(parts) == 2 {
			namespace = parts[0]
		}

		// If the value is a map, process the namespace information within it.
		if valueMap, ok := value.(map[string]any); ok {
			if nsMap, nsOk := valueMap[s.namespaceKey].(map[string]any); nsOk {
				for prefix, uri := range nsMap {
					if uriStr, isStr := uri.(string); isStr {
						nsManager.addNamespace(prefix, uriStr)
					}
				}
			}
		}

		element := s.mapToXMLElement(key, value, namespace, parts)
		elements = append(elements, element)
	}
	return elements
}

func (s soapREST) mapToXMLElement(elementName string, value any, namespace string, parts []string) xmlElement {
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

		// Processing attributes
		if attrMap, ok := v[s.attributeKey].(map[string]any); ok {
			for attrKey, attrValue := range attrMap {
				element.Attrs = append(element.Attrs, xml.Attr{
					Name:  xml.Name{Local: attrKey},
					Value: fmt.Sprintf("%v", attrValue),
				})
			}
		}

		// Processing text content
		if textContent, ok := v[s.textKey].(string); ok {
			element.Content = sanitizeControlCharacters(textContent)
		}

		// Processing child elements
		for childKey, childValue := range v {
			if childKey == s.attributeKey || childKey == s.textKey || childKey == s.namespaceKey {
				continue
			}

			// When processing child elements,
			// parse the namespace information from the child element's own key
			child := s.mapToXMLElement(childKey, childValue, "", nil)
			element.children = append(element.children, child)
		}

	case []any:
		// Process each element of the array as a child element
		for _, item := range v {
			child := s.mapToXMLElement(s.arrayKey, item, "", nil)
			element.children = append(element.children, child)
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
