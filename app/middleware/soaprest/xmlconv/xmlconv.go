// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package xmlconv

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
)

// EncodeDecoder provides interface for encoding
// and decoding XML document from and to JSON object.
type EncodeDecoder interface {
	// Encode encodes the given obj into XML document.
	// Structure of the obj must follow the implementer's instruction.
	Encode(encoder *xml.Encoder, obj map[string]any) error
	// Decode decodes XML document read from the decoder
	// and convert it into an object (map[string]any) or []object ([]map[string]any)
	// which is intended to be converted into JSON.
	Decode(decoder *xml.Decoder) (any, error)
}

// Converter converts XML to JSON document and
// JSON to XML document.
type Converter struct {
	EncodeDecoder
	// Header is the XML header line string.
	// If not empty, the value is written into the
	// header part of XML document.
	// Typically [encoding/xml.Header] should be used.
	Header string

	xmlEncoderOpts  []func(*xml.Encoder)
	xmlDecoderOpts  []func(*xml.Decoder)
	jsonEncoderOpts []func(*json.Encoder)
	jsonDecoderOpts []func(*json.Decoder)
}

// WithXMLEncoderOpts registers XML encoder options.
// Following example apply indent to the output.
//
// Example:
//
//	func(e *xml.Encoder){
//		e.Indent("", "  ")
//	}
func (c *Converter) WithXMLEncoderOpts(opts ...func(*xml.Encoder)) {
	for _, o := range opts {
		if o == nil {
			continue
		}
		c.xmlEncoderOpts = append(c.xmlEncoderOpts, o)
	}
}

// WithXMLDecoderOpts registers XML decoder options.
// Following example adds ability to read non UTF-8 XML document.
//
// Example:
//
//	func(d *xml.Decoder) {
//		d.CharsetReader = func(charset string, input io.Reader) (io.Reader, error) {
//			switch {
//			case strings.EqualFold(charset, "Shift_JIS"):
//				return transform.NewReader(input, japanese.ShiftJIS.NewDecoder()), nil
//			}
//			return input, nil
//		}
//	}
func (c *Converter) WithXMLDecoderOpts(opts ...func(*xml.Decoder)) {
	for _, o := range opts {
		if o == nil {
			continue
		}
		c.xmlDecoderOpts = append(c.xmlDecoderOpts, o)
	}
}

// WithJSONEncoderOpts registers JSON encoder options.
// Following example make the encoder not to apply HTML escape and
// to use indent.
//
// Example:
//
//	func(e *json.Encoder) {
//		e.SetEscapeHTML(false)
//		e.SetIndent("", "  ")
//	}
func (c *Converter) WithJSONEncoderOpts(opts ...func(*json.Encoder)) {
	for _, o := range opts {
		if o == nil {
			continue
		}
		c.jsonEncoderOpts = append(c.jsonEncoderOpts, o)
	}
}

// WithJSONDecoderOpts registers JSON decoder options.
// Following example applies decoder to parse numbers as [encoding/json.Number].
//
// Example:
//
//	func(d *json.Decoder){
//		d.UseNumber()
//	}
func (c *Converter) WithJSONDecoderOpts(opts ...func(*json.Decoder)) {
	for _, o := range opts {
		if o == nil {
			continue
		}
		c.jsonDecoderOpts = append(c.jsonDecoderOpts, o)
	}
}

// XMLtoJSON converts XML document into JSON document.
func (c *Converter) XMLtoJSON(b []byte) ([]byte, error) {
	decoder := xml.NewDecoder(bytes.NewReader(b))
	for _, opt := range c.xmlDecoderOpts {
		opt(decoder)
	}
	obj, err := c.Decode(decoder)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	for _, opt := range c.jsonEncoderOpts {
		opt(encoder)
	}
	if err := encoder.Encode(obj); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// JSONtoXML converts JSON document into XML document.
// Given json structure must follow the instruction of
// the [Converter.EncodeDecoder].
func (c *Converter) JSONtoXML(b []byte) ([]byte, error) {
	decoder := json.NewDecoder(bytes.NewReader(b))
	for _, opt := range c.jsonDecoderOpts {
		opt(decoder)
	}
	var a any
	if err := decoder.Decode(&a); err != nil {
		return nil, err
	}

	var objs []map[string]any
	switch t := a.(type) {
	case map[string]any:
		objs = []map[string]any{t}
	case []any:
		for _, tt := range t {
			m, ok := tt.(map[string]any)
			if !ok {
				return nil, ErrJSON.New("got %T:%+v", tt, tt)
			}
			objs = append(objs, m)
		}
	default:
		return nil, ErrJSON.New("got %T:%+v", a, a)
	}

	var buf bytes.Buffer
	if c.Header != "" {
		_, _ = buf.Write([]byte(c.Header))
	}
	encoder := xml.NewEncoder(&buf)
	for _, opt := range c.xmlEncoderOpts {
		opt(encoder)
	}
	for _, obj := range objs {
		if err := c.Encode(encoder, obj); err != nil {
			fmt.Println(string(b))
			panic(err)
			// return nil, err
		}
	}
	err := encoder.Flush() // Just in case.
	return buf.Bytes(), err
}
