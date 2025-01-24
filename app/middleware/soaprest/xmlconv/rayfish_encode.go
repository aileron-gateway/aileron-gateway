// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package xmlconv

import (
	"encoding/xml"
	"strings"
)

// Encode encodes the given obj into XML document.
// Resulting XML is written into the encoder.
// See the comment on [RayFish] for conversion rules.
func (r *RayFish) Encode(encoder *xml.Encoder, obj map[string]any) error {
	if len(obj) == 0 { // Ignore empty object.
		return nil
	}

	name, text, allChildren, err := r.parseItems(obj)
	if err != nil {
		return err
	}
	attrs, children, err := r.separateAttrs(allChildren)
	if err != nil {
		return err
	}

	start := xml.StartElement{
		Name: xml.Name{Local: restoreNamespace(r.NamespaceSep, name)},
		Attr: attrs,
	}

	var token xml.Token
	if r.JSONValue != nil {
		if token, err = r.XMLValue(text, &start); err != nil {
			return err
		}
	} else {
		if token, err = jsonValueToToken(r.TrimSpace, text); err != nil {
			return err
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
	for _, child := range children {
		if err := r.Encode(encoder, child); err != nil { // Child element.
			return err
		}
	}
	if err := encoder.EncodeToken(start.End()); err != nil { // End element.
		return err
	}
	return encoder.Flush()
}

// parseItems parses name, text and children.
// By default, it raises an error if undefined JSON key found.
// Use [RayFish.IgnoreUnusedKey] to ignore unknown, or unused JSON key.
//
// The input obj structure must follows:
//
//	map[string]any{
//		"#name": "alice", -----| Cannot omit. Must be string. [RayFish.NameKey] is used to extract.
//		"#text": "bob",   -----| Can omit. Any types. Nil if not exists. [RayFish.TextKey] is used to extract.
//		"#children": [    -----| Can omit. Must be []any. Nil if not exists. [RayFish.ChildrenKey] is used to extract.
//			......
//		],
//		"#unknown": "mallory", ---| Other keys are error by default. Use [RayFish.IgnoreUnusedKey] to ignore.
//	}
func (r *RayFish) parseItems(obj map[string]any) (name string, text any, children []any, err error) {
	var ok bool
	for k, v := range obj {
		switch k {
		case r.NameKey:
			name, ok = v.(string)
			if !ok {
				err = ErrDataType.New("NAME must be string. got %T:%+v", v, v)
				return
			}
		case r.TextKey:
			text = v
		case r.ChildrenKey:
			children, ok = v.([]any)
			if !ok {
				err = ErrDataType.New("CHILDREN must be array. got %T:%+v", v, v)
				return
			}
		default:
			if r.IgnoreUnusedKey {
				continue
			}
			err = ErrJSON.New("key not allowed. got %+v", k)
			return
		}
	}
	return name, text, children, nil
}

// separateAttrs separates children into XML attributes and other elements.
// Attributes should have [RayFish.AttrPrefix] in their names.
// On the other hand, non attributes must not have [RayFish.AttrPrefix] as the name prefix.
// All elements of the given children must be map[string]any type. Otherwise, results in an error.
// All elements must have an JSON attributes named [RayFish.NameKey].
func (r *RayFish) separateAttrs(children []any) (attrs []xml.Attr, elems []map[string]any, err error) {
	for _, child := range children {
		childMap, ok := child.(map[string]any)
		if !ok {
			err = ErrDataType.New("each CHILDREN must be object. got %T:%+v", child, child)
			return
		}
		nameAny, ok := childMap[r.NameKey]
		if !ok {
			err = ErrJSON.New("NAME not found. got %+v", childMap)
			return
		}
		nameStr, ok := nameAny.(string)
		if !ok {
			err = ErrDataType.New("NAME must be string. got %T:%v", nameAny, nameAny)
			return
		}
		if !strings.HasPrefix(nameStr, r.AttrPrefix) {
			elems = append(elems, childMap)
			continue
		}
		switch value := childMap[r.TextKey].(type) {
		case string:
			attrs = append(attrs, xml.Attr{
				Name:  xml.Name{Local: strings.TrimPrefix(nameStr, r.AttrPrefix)},
				Value: value,
			})
		case nil:
			attrs = append(attrs, xml.Attr{
				Name: xml.Name{Local: strings.TrimPrefix(nameStr, r.AttrPrefix)},
			})
		default:
			err = ErrDataType.New("Attribute value must be string or null. got %T:%+v", value, value)
			return
		}
	}
	return attrs, elems, nil
}
