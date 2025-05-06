// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package zxml

import (
	"bytes"
	"encoding/xml"
	"errors"
	"testing"
)

func TestRayfish_Encode(t *testing.T) {
	t.Parallel()
	base := RayFish{
		NameKey:      "#name",
		TextKey:      "#text",
		ChildrenKey:  "#children",
		AttrPrefix:   "@",
		NamespaceSep: ":",
		emptyVal:     "",
	}
	t.Run("name key", func(t *testing.T) {
		s := base
		s.NameKey = "$"
		input := map[string]any{
			"$":         "alice",
			"#text":     "bob",
			"#children": []any{map[string]any{"$": "@charlie", "#text": "david", "#children": []any{}}},
		}
		var buf bytes.Buffer
		enc := xml.NewEncoder(&buf)
		err := s.Encode(enc, input)
		AssertEqual(t, "error is not nil", nil, err)
		want := `<alice charlie="david">bob</alice>`
		AssertEqual(t, "xml not match", want, buf.String())
	})
	t.Run("text key", func(t *testing.T) {
		s := base
		s.TextKey = "$"
		input := map[string]any{
			"#name":     "alice",
			"$":         "bob",
			"#children": []any{map[string]any{"#name": "@charlie", "$": "david", "#children": []any{}}},
		}
		var buf bytes.Buffer
		enc := xml.NewEncoder(&buf)
		err := s.Encode(enc, input)
		AssertEqual(t, "error is not nil", nil, err)
		want := `<alice charlie="david">bob</alice>`
		AssertEqual(t, "xml not match", want, buf.String())
	})
	t.Run("attr prefix", func(t *testing.T) {
		s := base
		s.AttrPrefix = "%"
		input := map[string]any{
			"#name":     "alice",
			"#text":     "bob",
			"#children": []any{map[string]any{"#name": "%charlie", "#text": "david", "#children": []any{}}},
		}
		var buf bytes.Buffer
		enc := xml.NewEncoder(&buf)
		err := s.Encode(enc, input)
		AssertEqual(t, "error is not nil", nil, err)
		want := `<alice charlie="david">bob</alice>`
		AssertEqual(t, "xml not match", want, buf.String())
	})
	t.Run("namespace sep", func(t *testing.T) {
		s := base
		s.NamespaceSep = "_**_"
		input := map[string]any{
			"#name":     "alice",
			"#text":     "bob",
			"#children": []any{map[string]any{"#name": "@foo_**_charlie", "#text": "david", "#children": []any{}}},
		}
		var buf bytes.Buffer
		enc := xml.NewEncoder(&buf)
		err := s.Encode(enc, input)
		AssertEqual(t, "error is not nil", nil, err)
		want := `<alice foo:charlie="david">bob</alice>`
		AssertEqual(t, "xml not match", want, buf.String())
	})
	t.Run("xml value", func(t *testing.T) {
		s := base
		s.XMLValue = func(a any, se *xml.StartElement) (xml.Token, error) {
			return xml.CharData([]byte("value")), nil
		}
		input := map[string]any{
			"#name":     "alice",
			"#text":     "bob",
			"#children": []any{map[string]any{"#name": "@charlie", "#text": "david", "#children": []any{}}},
		}
		var buf bytes.Buffer
		enc := xml.NewEncoder(&buf)
		err := s.Encode(enc, input)
		AssertEqual(t, "error is not nil", nil, err)
		want := `<alice charlie="david">value</alice>`
		AssertEqual(t, "xml not match", want, buf.String())
	})
	t.Run("xml value error", func(t *testing.T) {
		s := base
		e := errors.New("parse error")
		s.XMLValue = func(a any, se *xml.StartElement) (xml.Token, error) {
			return xml.CharData([]byte("value")), e
		}
		input := map[string]any{
			"#name":     "alice",
			"#text":     "bob",
			"#children": []any{map[string]any{"#name": "@charlie", "#text": "david", "#children": []any{}}},
		}
		var buf bytes.Buffer
		enc := xml.NewEncoder(&buf)
		err := s.Encode(enc, input)
		AssertEqualErr(t, "error not match", e, err)
	})

	t.Run("name not found", func(t *testing.T) {
		s := base
		input := map[string]any{
			"#name":     "alice",
			"#text":     "bob",
			"#children": []any{map[string]any{"#text": "david", "#children": []any{}}},
		}
		var buf bytes.Buffer
		enc := xml.NewEncoder(&buf)
		err := s.Encode(enc, input)
		AssertEqualErr(t, "error not match", &XMLError{Cause: CauseJSONStruct}, err)
	})
	t.Run("name not string", func(t *testing.T) {
		s := base
		input := map[string]any{
			"#name":     "alice",
			"#text":     "bob",
			"#children": []any{map[string]any{"#name": 123, "#text": "david", "#children": []any{}}},
		}
		var buf bytes.Buffer
		enc := xml.NewEncoder(&buf)
		err := s.Encode(enc, input)
		AssertEqualErr(t, "error not match", &XMLError{Cause: CauseDataType}, err)
	})
	t.Run("name type error", func(t *testing.T) {
		s := base
		input := map[string]any{
			"#name":     uint(123),
			"#text":     "bob",
			"#children": []any{map[string]any{"#name": "@charlie", "#text": "david", "#children": []any{}}},
		}
		var buf bytes.Buffer
		enc := xml.NewEncoder(&buf)
		err := s.Encode(enc, input)
		AssertEqualErr(t, "error not match", &XMLError{Cause: CauseDataType}, err)
	})
	t.Run("nil name", func(t *testing.T) {
		s := base
		input := map[string]any{
			"#name":     nil,
			"#text":     "bob",
			"#children": []any{map[string]any{"#name": "@charlie", "#text": "david", "#children": []any{}}},
		}
		var buf bytes.Buffer
		enc := xml.NewEncoder(&buf)
		err := s.Encode(enc, input)
		AssertEqualErr(t, "error not match", &XMLError{Cause: CauseDataType}, err)
	})
	t.Run("text type error", func(t *testing.T) {
		s := base
		input := map[string]any{
			"#name":     "alice",
			"#text":     uint(123),
			"#children": []any{map[string]any{"#name": "@charlie", "#text": "david", "#children": []any{}}},
		}
		var buf bytes.Buffer
		enc := xml.NewEncoder(&buf)
		err := s.Encode(enc, input)
		AssertEqualErr(t, "error not match", &XMLError{Cause: CauseDataType}, err)
	})
	t.Run("nil text", func(t *testing.T) {
		s := base
		input := map[string]any{
			"#name":     "alice",
			"#text":     nil,
			"#children": []any{map[string]any{"#name": "@charlie", "#text": "david", "#children": []any{}}},
		}
		var buf bytes.Buffer
		enc := xml.NewEncoder(&buf)
		err := s.Encode(enc, input)
		AssertEqual(t, "error is not nil", nil, err)
		want := `<alice charlie="david"></alice>`
		AssertEqual(t, "xml not match", want, buf.String())
	})
	t.Run("nil attribute", func(t *testing.T) {
		s := base
		input := map[string]any{
			"#name":     "alice",
			"#text":     "bob",
			"#children": []any{map[string]any{"#name": "@charlie", "#text": nil, "#children": []any{}}},
		}
		var buf bytes.Buffer
		enc := xml.NewEncoder(&buf)
		err := s.Encode(enc, input)
		AssertEqual(t, "error is not nil", nil, err)
		want := `<alice charlie="">bob</alice>`
		AssertEqual(t, "xml not match", want, buf.String())
	})
	t.Run("attribute type error", func(t *testing.T) {
		s := base
		input := map[string]any{
			"#name":     "alice",
			"#text":     "bob",
			"#children": []any{map[string]any{"#name": "@charlie", "#text": uint(123), "#children": []any{}}},
		}
		var buf bytes.Buffer
		enc := xml.NewEncoder(&buf)
		err := s.Encode(enc, input)
		AssertEqualErr(t, "error not match", &XMLError{Cause: CauseDataType}, err)
	})
	t.Run("children type invalid", func(t *testing.T) {
		s := base
		input := map[string]any{
			"#name":     "alice",
			"#text":     "bob",
			"#children": map[string]any{},
		}
		var buf bytes.Buffer
		enc := xml.NewEncoder(&buf)
		err := s.Encode(enc, input)
		AssertEqualErr(t, "error not match", &XMLError{Cause: CauseDataType}, err)
	})
	t.Run("children type invalid", func(t *testing.T) {
		s := base
		input := map[string]any{
			"#name":     "alice",
			"#text":     nil,
			"#children": []any{"invalid"},
		}
		var buf bytes.Buffer
		enc := xml.NewEncoder(&buf)
		err := s.Encode(enc, input)
		AssertEqualErr(t, "error not match", &XMLError{Cause: CauseDataType}, err)
	})
	t.Run("invalid key", func(t *testing.T) {
		s := base
		input := map[string]any{
			"#name":    "alice",
			"#text":    "bob",
			"#invalid": []any{},
		}
		var buf bytes.Buffer
		enc := xml.NewEncoder(&buf)
		err := s.Encode(enc, input)
		AssertEqualErr(t, "error not match", &XMLError{Cause: CauseJSONStruct}, err)
	})
	t.Run("empty start elem", func(t *testing.T) {
		s := base
		input := map[string]any{
			"#name":     "",
			"#text":     "bob",
			"#children": []any{map[string]any{"#name": "@charlie", "#text": "david", "#children": []any{}}},
		}
		var buf bytes.Buffer
		enc := xml.NewEncoder(&buf)
		err := s.Encode(enc, input)
		AssertEqualErr(t, "error not match", &XMLError{Cause: CauseXMLEncoder}, err)
	})
	t.Run("token encode error", func(t *testing.T) {
		s := base
		s.XMLValue = func(a any, se *xml.StartElement) (xml.Token, error) { return xml.StartElement{}, nil }
		input := map[string]any{
			"#name":     "alice",
			"#text":     "bob",
			"#children": []any{map[string]any{"#name": "@charlie", "#text": "david", "#children": []any{}}},
		}
		var buf bytes.Buffer
		enc := xml.NewEncoder(&buf)
		err := s.Encode(enc, input)
		AssertEqualErr(t, "error not match", &XMLError{Cause: CauseXMLEncoder}, err)
	})
	t.Run("end token encode error", func(t *testing.T) {
		s := base
		s.XMLValue = func(a any, se *xml.StartElement) (xml.Token, error) {
			return xml.EndElement{Name: xml.Name{Local: "alice"}}, nil
		}
		input := map[string]any{
			"#name":     "alice",
			"#text":     "bob",
			"#children": []any{map[string]any{"#name": "@charlie", "#text": "david", "#children": []any{}}},
		}
		var buf bytes.Buffer
		enc := xml.NewEncoder(&buf)
		err := s.Encode(enc, input)
		AssertEqualErr(t, "error not match", &XMLError{Cause: CauseXMLEncoder}, err)
	})
	t.Run("child encode error", func(t *testing.T) {
		s := base
		s.XMLValue = func(a any, se *xml.StartElement) (xml.Token, error) {
			if s := a.(string); s == "david" {
				return xml.StartElement{}, nil
			}
			return nil, nil
		}
		input := map[string]any{
			"#name":     "alice",
			"#text":     "bob",
			"#children": []any{map[string]any{"#name": "charlie", "#text": "david", "#children": []any{}}},
		}
		var buf bytes.Buffer
		enc := xml.NewEncoder(&buf)
		err := s.Encode(enc, input)
		AssertEqualErr(t, "error not match", &XMLError{Cause: CauseXMLEncoder}, err)
	})
}
