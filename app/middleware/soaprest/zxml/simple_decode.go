package zxml

import (
	"encoding/xml"
	"io"
	"slices"
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
func (s *Simple) Decode(decoder *xml.Decoder) (any, error) {
	objs := make([]any, 0, 1)
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
			content, err := s.parseContent(decoder, t, t.End(), nil)
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

func (s *Simple) parseContent(decoder *xml.Decoder, start xml.StartElement, end xml.EndElement, ns [][2]string) (map[string]any, error) {
	// Append namespace of this element.
	ns = append(ns, parseNamespace(start.Attr, ns)...)

	// Put attributes in the map.
	attrs := make(map[string]any, len(start.Attr))
	for _, attr := range start.Attr {
		attrs[attrName(attr.Name, s.AttrPrefix, s.NamespaceSep, ns)] = attr.Value
	}

	var text string
	var keys []string
	var children []any

Loop:
	for {
		token, err := decoder.Token()
		if err != nil {
			if err != io.EOF {
				err = &XMLError{Err: err, Cause: CauseXMLDecoder}
			}
			return nil, err
		}

		switch t := token.(type) {
		case xml.CharData:
			if s.TrimSpace {
				text += strings.TrimSpace(string(t))
			} else {
				text += string(t)
			}
		case xml.StartElement:
			content, err := s.parseContent(decoder, t, t.End(), ns)
			if err != nil {
				return nil, err
			}
			for k, v := range content {
				keys = append(keys, k)
				children = append(children, v)
			}
		case xml.EndElement:
			if t == end {
				break Loop
			}
		}
	}

	// Merge attributes and child elements.
	// Key-value pairs of keys-children are stored
	// in the attrs variable.
	mergeChildren(attrs, keys, children)

	// Convert XML text value to JSON value.
	var val any
	if s.JSONValue != nil {
		v, err := s.JSONValue(text, start)
		if err != nil {
			return nil, err
		}
		val = v
	} else {
		if text != "" {
			val = text
		}
	}

	if len(attrs) == 0 {
		if val == nil { // Replace empty value.
			val = s.emptyVal
		}
		if s.PreferShort {
			return map[string]any{
				tokenName(start.Name, s.NamespaceSep, ns): val,
			}, nil
		}
	}

	// Add XML content if it is not nil or empty string.
	if val != nil {
		attrs[s.TextKey] = val
	}

	return map[string]any{
		tokenName(start.Name, s.NamespaceSep, ns): attrs,
	}, nil
}

// mergeChildren merges child elements into the target map.
// The given target must not be nil.
// mergeChildren does not check if a key is already exists in the
// target map but overwrite values.
// Children with same keys are converted in to a set of values as []any.
// Length of keys and children must be same and the key[i]:children[i]
// must be a pair as key-value.
// As shown in the table below, values with the same keys are packed
// into a slice. Values with different keys are put in
// the target map object with their key-value pairs.
func mergeChildren(target map[string]any, keys []string, children []any) {
	if len(keys) == 1 {
		target[keys[0]] = children[0]
		return
	}
	numKey := map[string]int{}
	for _, k := range keys {
		numKey[k] += 1 // Check key duplication by counting on it.
	}
	if len(numKey) == 1 { // Special case. Shortcut for performance.
		target[keys[0]] = children
		return
	}
	for key, num := range numKey {
		if num == 1 { // The key is unique.
			target[key] = children[slices.Index(keys, key)]
			continue
		}
		arr := make([]any, 0, num)
		for i, k := range keys {
			if k == key {
				arr = append(arr, children[i])
			}
		}
		target[key] = arr
	}
}
