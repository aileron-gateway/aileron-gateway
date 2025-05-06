package zxml

import (
	"encoding/xml"
	"fmt"
	"strings"
)

// Encode encodes the given obj into XML document.
// Resulting XML is written into the encoder.
// See the comment on [BadgerFish] for conversion rules.
func (b *BadgerFish) Encode(encoder *xml.Encoder, obj map[string]any) error {
	for k, v := range obj {
		if err := b.encode(encoder, k, v); err != nil {
			return err
		}
	}
	return encoder.Flush()
}

func (b *BadgerFish) encode(encoder *xml.Encoder, key string, elem any) (err error) {
	// Process slice element such as
	// 	"bob": [
	// 		{	"$": "charlie" },
	// 		{ "$": "edgar" }
	// 	]
	if arr, ok := elem.([]any); ok {
		for _, a := range arr {
			if err := b.encode(encoder, key, a); err != nil {
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
		text, attrs, children, err = b.parseItems(e)
		if err != nil {
			return err
		}
	default: // Type is not []any.
		text = elem
	}

	start := xml.StartElement{
		Name: xml.Name{Local: restoreNamespace(b.NamespaceSep, key)},
		Attr: attrs,
	}

	var token xml.Token
	if text != nil {
		if b.XMLValue != nil {
			if token, err = b.XMLValue(text, &start); err != nil {
				return err
			}
		} else {
			if token, err = jsonValueToToken(b.TrimSpace, text); err != nil {
				return err
			}
		}
	}

	if err := encoder.EncodeToken(start); err != nil { // Start element.
		return &XMLError{Err: err, Cause: CauseXMLEncoder}
	}
	if token != nil {
		if err := encoder.EncodeToken(token); err != nil { // Text content.
			return &XMLError{Err: err, Cause: CauseXMLEncoder}
		}
	}
	for k, v := range children {
		if err := b.encode(encoder, k, v); err != nil { // Child element.
			return err
		}
	}
	if err := encoder.EncodeToken(start.End()); err != nil { // End element.
		return &XMLError{Err: err, Cause: CauseXMLEncoder}
	}
	return nil
}

func (b *BadgerFish) parseItems(obj map[string]any) (text any, attrs []xml.Attr, children map[string]any, err error) {
	children = map[string]any{}
	for k, v := range obj {
		switch {
		case k == b.TextKey:
			text = v
		case strings.HasPrefix(k, b.AttrPrefix):
			name := strings.TrimPrefix(k, b.AttrPrefix)
			name = restoreNamespace(b.NamespaceSep, name)
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
			case map[string]any:
				for kk, vv := range value {
					vv, ok := vv.(string)
					if !ok {
						err = &XMLError{
							Cause:  CauseDataType,
							Detail: fmt.Sprintf("Namespace must be string or null. got %T:%+v", v, v),
						}
						return
					}
					if kk == b.TextKey {
						attrs = append(attrs, xml.Attr{
							Name:  xml.Name{Local: name},
							Value: vv,
						})
					} else {
						attrs = append(attrs, xml.Attr{
							Name:  xml.Name{Local: name + ":" + kk},
							Value: vv,
						})
					}
				}
			default:
				err = &XMLError{
					Cause:  CauseDataType,
					Detail: fmt.Sprintf("Attribute must be string,null or map[string]any. got %T:%+v", value, value),
				}
				return
			}
		default:
			children[k] = v
		}
	}
	return text, attrs, children, nil
}
