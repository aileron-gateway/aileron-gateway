// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package xmlconv

import (
	"encoding/xml"
	"strings"
)

// Encode encodes the given obj into XML document.
// Resulting XML is written into the encoder.
// See the comment on [Simple] for conversion rules.
func (s *Simple) Encode(encoder *xml.Encoder, obj map[string]any) error {
	for k, v := range obj {
		if err := s.encode(encoder, k, v); err != nil {
			return err
		}
	}
	return encoder.Flush()
}

func (s *Simple) encode(encoder *xml.Encoder, key string, elem any) (err error) {
	// Process slice element such as
	// 	"bob": [
	// 		{	"$": "charlie" },
	// 		{ "$": "edgar" }
	// 	]
	// or
	// 	"bob": [ "charlie", "edgar" ]
	if arr, ok := elem.([]any); ok {
		for _, a := range arr {
			if err := s.encode(encoder, key, a); err != nil {
				return err
			}
		}
		return nil
	}

	var text any
	var attrs []xml.Attr
	var children map[string]any

	switch e := elem.(type) {
	case map[string]any:
		text, attrs, children, err = s.parseItems(e)
		if err != nil {
			return err
		}
	default: // Type is not []any.
		text = elem
	}

	start := xml.StartElement{
		Name: xml.Name{Local: restoreNamespace(s.NamespaceSep, key)},
		Attr: attrs,
	}

	var token xml.Token
	if text != nil {
		if s.JSONValue != nil {
			if token, err = s.XMLValue(text, &start); err != nil {
				return err
			}
		} else {
			if token, err = jsonValueToToken(s.TrimSpace, text); err != nil {
				return err
			}
		}
	}

	if err := encoder.EncodeToken(start); err != nil { // Start element.
		return err
	}
	if token != nil {
		if err := encoder.EncodeToken(token); err != nil { // Text content.
			return err
		}
	}
	for k, v := range children {
		if err := s.encode(encoder, k, v); err != nil { // Child element.
			return err
		}
	}
	if err := encoder.EncodeToken(start.End()); err != nil { // End element.
		return err
	}
	return nil
}

func (r *Simple) parseItems(obj map[string]any) (text any, attrs []xml.Attr, children map[string]any, err error) {
	children = map[string]any{}
	for k, v := range obj {
		switch {
		case k == r.TextKey:
			text = v
		case strings.HasPrefix(k, r.AttrPrefix):
			name := strings.TrimPrefix(k, r.AttrPrefix)
			switch value := v.(type) {
			case string:
				attrs = append(attrs, xml.Attr{
					Name:  xml.Name{Local: name},
					Value: value,
				})
			case nil:
				attrs = append(attrs, xml.Attr{
					Name: xml.Name{Local: name},
				})
			default:
				err = ErrDataType.New("Attribute value must be string or null. got %T:%+v", value, value)
				return
			}
		default:
			children[k] = v
		}
	}
	return text, attrs, children, nil
}
