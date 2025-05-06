// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package zxml

import (
	"encoding/xml"
	"io"
	"strings"
)

// Decode decodes XML document read by the decoder into Go values.
//
// Returned value can be:
//   - map[string]any
//   - []map[string]any
//
// Decode:
//   - Ignores [encoding/xml.Comment]
//   - Ignores [encoding/xml.ProcInst]
//   - Ignores [encoding/xml.Directive]
//   - Does not identify CDATA (Limitation of [encoding/xml])
func (r *RayFish) Decode(decoder *xml.Decoder) (any, error) {
	objs := make([]map[string]any, 0)
	var token xml.Token
	var err error
	for {
		if token, err = decoder.Token(); err != nil {
			if err != io.EOF {
				return nil, &XMLError{Err: err, Cause: CauseXMLDecoder}
			}
			break
		}
		switch t := token.(type) {
		case xml.StartElement:
			content, err := r.decode(decoder, t, t.End(), nil)
			if err != nil {
				return nil, err
			}
			objs = append(objs, content)
		}
	}
	if err == io.EOF {
		err = nil
	}
	switch len(objs) {
	case 0:
		return map[string]any{}, err
	case 1:
		return objs[0], err
	default:
		return objs, err
	}
}

func (r *RayFish) decode(decoder *xml.Decoder, start xml.StartElement, end xml.EndElement, ns [][2]string) (map[string]any, error) {
	// Append namespace of this element.
	ns = append(ns, parseNamespace(start.Attr, ns)...)

	// Register attributes as children.
	children := make([]map[string]any, 0, len(start.Attr))
	for _, attr := range start.Attr {
		children = append(children, map[string]any{
			r.NameKey:     attrName(attr.Name, r.AttrPrefix, r.NamespaceSep, ns),
			r.TextKey:     attr.Value,
			r.ChildrenKey: make([]map[string]any, 0), // Attribute has no child.
		})
	}

	var text string
Loop:
	for {
		token, err := decoder.Token()
		if err != nil {
			return nil, &XMLError{Err: err, Cause: CauseXMLDecoder}
		}
		switch t := token.(type) {
		case xml.StartElement:
			content, err := r.decode(decoder, t, t.End(), ns)
			if err != nil {
				return nil, err
			}
			children = append(children, content)
		case xml.CharData:
			if r.TrimSpace {
				text += strings.TrimSpace(string(t))
			} else {
				text += string(t)
			}
		case xml.EndElement:
			if t == end {
				break Loop
			}
		}
	}

	// Convert XML text value to JSON value.
	var val any = text
	if r.JSONValue != nil {
		v, err := r.JSONValue(text, start)
		if err != nil {
			return nil, err
		}
		val = v
	} else {
		if text == "" {
			val = r.emptyVal
		}
	}

	return map[string]any{
		r.NameKey:     tokenName(start.Name, r.NamespaceSep, ns),
		r.TextKey:     val,
		r.ChildrenKey: children,
	}, nil
}
