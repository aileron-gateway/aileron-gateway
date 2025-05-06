package zxml

import (
	"bytes"
	"encoding/xml"
	"errors"
	"testing"
)

func TestBadgerFish_Encode(t *testing.T) {
	t.Parallel()
	base := BadgerFish{
		TextKey:      "$",
		AttrPrefix:   "@",
		NamespaceSep: ":",
		emptyVal:     make(map[string]any, 0),
	}
	t.Run("text key", func(t *testing.T) {
		s := base
		s.TextKey = "#text"
		input := map[string]any{"alice": map[string]any{"#text": "bob", "@charlie": "david"}}
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
		input := map[string]any{"alice": map[string]any{"$": "bob", "%charlie": "david"}}
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
		input := map[string]any{"alice": map[string]any{"$": "bob", "@foo_**_charlie": "david"}}
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
		input := map[string]any{"alice": map[string]any{"$": "bob", "@charlie": "david"}}
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
		input := map[string]any{"alice": map[string]any{"$": "bob", "@charlie": "david"}}
		var buf bytes.Buffer
		enc := xml.NewEncoder(&buf)
		err := s.Encode(enc, input)
		AssertEqualErr(t, "error not match", e, err)
	})

	t.Run("slice value error", func(t *testing.T) {
		s := base
		input := map[string]any{"alice": map[string]any{"$": []any{uint(123)}, "@charlie": "david"}}
		var buf bytes.Buffer
		enc := xml.NewEncoder(&buf)
		err := s.Encode(enc, input)
		AssertEqualErr(t, "error not match", &XMLError{Cause: CauseDataType}, err)
	})
	t.Run("invalid text", func(t *testing.T) {
		s := base
		input := map[string]any{"alice": map[string]any{"$": uint(123), "@charlie": "david"}}
		var buf bytes.Buffer
		enc := xml.NewEncoder(&buf)
		err := s.Encode(enc, input)
		AssertEqualErr(t, "error not match", &XMLError{Cause: CauseDataType}, err)
	})
	t.Run("array value invalid", func(t *testing.T) {
		s := base
		input := map[string]any{"alice": []any{uint(123)}}
		var buf bytes.Buffer
		enc := xml.NewEncoder(&buf)
		err := s.Encode(enc, input)
		AssertEqualErr(t, "error not match", &XMLError{Cause: CauseDataType}, err)
	})
	t.Run("invalid attribute", func(t *testing.T) {
		s := base
		input := map[string]any{"alice": map[string]any{"$": "bob", "@charlie": uint(123)}}
		var buf bytes.Buffer
		enc := xml.NewEncoder(&buf)
		err := s.Encode(enc, input)
		AssertEqualErr(t, "error not match", &XMLError{Cause: CauseDataType}, err)
	})
	t.Run("nil attribute", func(t *testing.T) {
		s := base
		input := map[string]any{"alice": map[string]any{"$": "bob", "@charlie": nil}}
		var buf bytes.Buffer
		enc := xml.NewEncoder(&buf)
		err := s.Encode(enc, input)
		AssertEqual(t, "error is not nil", nil, err)
		want := `<alice charlie="">bob</alice>`
		AssertEqual(t, "xml not match", want, buf.String())
	})
	t.Run("empty start elem", func(t *testing.T) {
		s := base
		input := map[string]any{"": map[string]any{"$": "bob"}}
		var buf bytes.Buffer
		enc := xml.NewEncoder(&buf)
		err := s.Encode(enc, input)
		AssertEqualErr(t, "error not match", &XMLError{Cause: CauseXMLEncoder}, err)
	})
	t.Run("token encode error", func(t *testing.T) {
		s := base
		s.XMLValue = func(a any, se *xml.StartElement) (xml.Token, error) { return xml.StartElement{}, nil }
		input := map[string]any{"alice": map[string]any{"$": "bob"}}
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
		input := map[string]any{"alice": map[string]any{"$": "bob"}}
		var buf bytes.Buffer
		enc := xml.NewEncoder(&buf)
		err := s.Encode(enc, input)
		AssertEqualErr(t, "error not match", &XMLError{Cause: CauseXMLEncoder}, err)
	})
	t.Run("child encode error", func(t *testing.T) {
		s := base
		s.XMLValue = func(a any, se *xml.StartElement) (xml.Token, error) { return xml.StartElement{}, nil }
		input := map[string]any{"alice": map[string]any{"bob": map[string]any{"$": "david"}}}
		var buf bytes.Buffer
		enc := xml.NewEncoder(&buf)
		err := s.Encode(enc, input)
		AssertEqualErr(t, "error not match", &XMLError{Cause: CauseXMLEncoder}, err)
	})
}
