package zxml

import (
	"encoding/xml"
	"errors"
	"reflect"
	"strings"
	"testing"
)

func TestRayFish_Decode(t *testing.T) {
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
		d, err := s.Decode(xml.NewDecoder(strings.NewReader(`<alice charlie="david">bob</alice>`)))
		want := map[string]any{
			"$":         "alice",
			"#text":     "bob",
			"#children": []map[string]any{{"$": "@charlie", "#text": "david", "#children": []map[string]any{}}},
		}
		AssertEqual(t, "map not match", true, reflect.DeepEqual(want, d))
		AssertEqual(t, "error is not nil", nil, err)
	})
	t.Run("text key", func(t *testing.T) {
		s := base
		s.TextKey = "$"
		d, err := s.Decode(xml.NewDecoder(strings.NewReader(`<alice charlie="david">bob</alice>`)))
		want := map[string]any{
			"#name":     "alice",
			"$":         "bob",
			"#children": []map[string]any{{"#name": "@charlie", "$": "david", "#children": []map[string]any{}}},
		}
		AssertEqual(t, "map not match", true, reflect.DeepEqual(want, d))
		AssertEqual(t, "error is not nil", nil, err)
	})
	t.Run("attr prefix", func(t *testing.T) {
		s := base
		s.AttrPrefix = "%"
		d, err := s.Decode(xml.NewDecoder(strings.NewReader(`<alice charlie="david">bob</alice>`)))
		want := map[string]any{
			"#name":     "alice",
			"#text":     "bob",
			"#children": []map[string]any{{"#name": "%charlie", "#text": "david", "#children": []map[string]any{}}},
		}
		AssertEqual(t, "map not match", true, reflect.DeepEqual(want, d))
		AssertEqual(t, "error is not nil", nil, err)
	})
	t.Run("namespace sep", func(t *testing.T) {
		s := base
		s.NamespaceSep = "_**_"
		d, err := s.Decode(xml.NewDecoder(strings.NewReader(`<alice foo:charlie="david">bob</alice>`)))
		want := map[string]any{
			"#name":     "alice",
			"#text":     "bob",
			"#children": []map[string]any{{"#name": "@foo_**_charlie", "#text": "david", "#children": []map[string]any{}}},
		}
		AssertEqual(t, "map not match", true, reflect.DeepEqual(want, d))
		AssertEqual(t, "error is not nil", nil, err)
	})
	t.Run("trim space", func(t *testing.T) {
		s := base
		s.TrimSpace = true
		d, err := s.Decode(xml.NewDecoder(strings.NewReader(`<alice charlie="david"> bob </alice>`)))
		want := map[string]any{
			"#name":     "alice",
			"#text":     "bob",
			"#children": []map[string]any{{"#name": "@charlie", "#text": "david", "#children": []map[string]any{}}},
		}
		AssertEqual(t, "map not match", true, reflect.DeepEqual(want, d))
		AssertEqual(t, "error is not nil", nil, err)
	})
	t.Run("json value", func(t *testing.T) {
		s := base
		s.JSONValue = func(s string, se xml.StartElement) (any, error) { return "value", nil }
		d, err := s.Decode(xml.NewDecoder(strings.NewReader(`<alice charlie="david">bob</alice>`)))
		want := map[string]any{
			"#name":     "alice",
			"#text":     "value",
			"#children": []map[string]any{{"#name": "@charlie", "#text": "david", "#children": []map[string]any{}}},
		}
		AssertEqual(t, "map not match", true, reflect.DeepEqual(want, d))
		AssertEqual(t, "error is not nil", nil, err)
	})
	t.Run("json value error", func(t *testing.T) {
		s := base
		e := errors.New("parse error")
		s.JSONValue = func(s string, se xml.StartElement) (any, error) { return nil, e }
		d, err := s.Decode(xml.NewDecoder(strings.NewReader(`<alice charlie="david">bob</alice>`)))
		AssertEqualErr(t, "error not match", e, err)
		AssertEqual(t, "map is not nil", true, reflect.DeepEqual(nil, d))
	})
	t.Run("inner json value error", func(t *testing.T) {
		s := base
		e := errors.New("parse error")
		s.JSONValue = func(s string, se xml.StartElement) (any, error) { return nil, e }
		d, err := s.Decode(xml.NewDecoder(strings.NewReader(`<alice><bob>charlie</bob></alice>`)))
		AssertEqualErr(t, "error not match", e, err)
		AssertEqual(t, "map is not nil", true, reflect.DeepEqual(nil, d))
	})
}

func TestRayFish_WithEmptyValue(t *testing.T) {
	t.Parallel()
	base := RayFish{
		NameKey:      "#name",
		TextKey:      "#text",
		ChildrenKey:  "#children",
		AttrPrefix:   "@",
		NamespaceSep: ":",
		emptyVal:     "",
	}
	t.Run("nil", func(t *testing.T) {
		s := base
		s.WithEmptyValue(nil)
		d, err := s.Decode(xml.NewDecoder(strings.NewReader(`<alice></alice>`)))
		want := map[string]any{
			"#name":     "alice",
			"#text":     nil,
			"#children": []map[string]any{},
		}
		AssertEqual(t, "map not match", true, reflect.DeepEqual(want, d))
		AssertEqual(t, "error is not nil", nil, err)
	})
	t.Run("string", func(t *testing.T) {
		s := base
		s.WithEmptyValue("")
		d, err := s.Decode(xml.NewDecoder(strings.NewReader(`<alice></alice>`)))
		want := map[string]any{
			"#name":     "alice",
			"#text":     "",
			"#children": []map[string]any{},
		}
		AssertEqual(t, "map not match", true, reflect.DeepEqual(want, d))
		AssertEqual(t, "error is not nil", nil, err)
	})
	t.Run("panic", func(t *testing.T) {
		s := base
		defer func() {
			r := recover()
			AssertEqualErr(t, "error not match", &XMLError{Cause: CauseEmptyVal}, r.(error))
		}()
		s.WithEmptyValue(123)
	})
}
