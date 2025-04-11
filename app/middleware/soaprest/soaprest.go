package soaprest

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

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

// removeNewlinesAndTabs removes unnecessary CR, LF, and Tab characters.
func removeNewlinesAndTabs(content string) string {
	replacer := strings.NewReplacer(
		"\n", "",
		"\t", "",
		"\r", "",
	)
	return replacer.Replace(content)
}

// getNodeName returns the name of the given xmlNode as a string.
func getNodeName(node xmlNode, nsCtx *namespaceContext, nsPrefix string) string {
	nodeName := node.XMLName.Local

	// Only supports conversions for SOAP 1.1
	if nodeName == soapBodyKey && node.XMLName.Space == soapNamespaceURI {
		// return like "soap_body"
		return nsPrefix + separatorChar + soapBodyKey
	} else if nodeName == soapHeaderKey && node.XMLName.Space == soapNamespaceURI {
		// return like "soap_header"
		return nsPrefix + separatorChar + soapHeaderKey
	}

	if node.XMLName.Space != "" {
		prefix := nsCtx.getPrefix(node.XMLName.Space)
		if prefix != "" {
			return prefix + separatorChar + nodeName
		}
	}

	return nodeName
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
func mapToXMLAttrs(attrMap map[string]any, nsManager *namespaceManager) []xml.Attr {
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

// toString is a heloper function that converts the given value into a string and returns it.
func toString(val any) string {
	switch v := val.(type) {
	case nil:
		return "<nil>"
	case string:
		return v
	case []byte:
		return string(v)
	case bool:
		return strconv.FormatBool(v)
	case int:
		return strconv.FormatInt(int64(v), 10)
	case int8:
		return strconv.FormatInt(int64(v), 10)
	case int16:
		return strconv.FormatInt(int64(v), 10)
	case int32:
		return strconv.FormatInt(int64(v), 10)
	case int64:
		return strconv.FormatInt(v, 10)
	case float32:
		return strconv.FormatFloat(float64(v), 'f', -1, 32)
	case float64:
		return strconv.FormatFloat(float64(v), 'f', -1, 64)
	case uint:
		return strconv.FormatUint(uint64(v), 10)
	case uint8:
		return strconv.FormatUint(uint64(v), 10)
	case uint16:
		return strconv.FormatUint(uint64(v), 10)
	case uint32:
		return strconv.FormatUint(uint64(v), 10)
	case uint64:
		return strconv.FormatUint(uint64(v), 10)
	case complex64:
		return strconv.FormatComplex(complex128(v), 'f', -1, 64)
	case complex128:
		return strconv.FormatComplex(complex128(v), 'f', -1, 128)
	default:
		return fmt.Sprint(v) // Fallback to "%+v"
	}
}

func makeCharsetReader(charset string, input io.Reader) (io.Reader, error) {
	if charset == "Shift-JIS" {
		return transform.NewReader(input, japanese.ShiftJIS.NewDecoder()), nil
	}
	return input, nil
}
