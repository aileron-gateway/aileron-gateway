// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package xmlconv_test

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"io"
	"reflect"
	"slices"
	"testing"
)

func xmlTokens(decoder *xml.Decoder, end xml.EndElement) (map[string]any, error) {
	key := ""
	m := map[string]any{}
	var tokens []xml.Token
	for {
		token, err := decoder.Token()
		if err != nil {
			if err == io.EOF {
				return m, nil
			}
			return m, err
		}
		switch t := token.(type) {
		case xml.Comment, xml.ProcInst, xml.Directive:
			continue
		case xml.StartElement:
			slices.SortFunc(t.Attr, func(a, b xml.Attr) int {
				if a.Name.Space+a.Name.Local > b.Name.Space+b.Name.Local {
					return 1
				}
				return -1
			})
			key = t.Name.Space + ":" + t.Name.Local
			children, err := xmlTokens(decoder, t.End())
			if err != nil {
				return nil, err
			}
			m[key] = children
		case xml.CharData:
			t = bytes.TrimSpace([]byte(t))
			if len(t) == 0 {
				continue
			}
			token = t
		case xml.EndElement:
			if t == end {
				return m, nil
			}
			v, ok := m[key]
			if ok {
				vv := v.([]xml.Token)
				m[key] = append(vv, tokens...)
			} else {
				m[key] = tokens
			}
			tokens = nil
		}
	}
}

func equalXML(t *testing.T, a, b []byte) bool {
	tokens1, err := xmlTokens(xml.NewDecoder(bytes.NewReader(a)), xml.EndElement{})
	if err != nil {
		panic(err)
	}
	tokens2, err := xmlTokens(xml.NewDecoder(bytes.NewReader(b)), xml.EndElement{})
	if err != nil {
		panic(err)
	}
	if equal := reflect.DeepEqual(tokens1, tokens2); equal {
		return true
	}
	t.Logf("XML-1: %#v\n", tokens1)
	t.Logf("XML-2: %#v\n", tokens2)
	return false
}

func equalJSON(t *testing.T, a, b []byte) bool {
	var obj1, obj2 interface{}
	if err := json.Unmarshal(a, &obj1); err != nil {
		panic(err)
	}
	if err := json.Unmarshal(b, &obj2); err != nil {
		panic(err)
	}
	if equal := reflect.DeepEqual(obj1, obj2); equal {
		return true
	}
	t.Logf("JSON-1: %#v\n", obj1)
	t.Logf("JSON-2: %#v\n", obj2)
	return false
}
