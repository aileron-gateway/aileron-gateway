package zxml_test

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"io"
	"maps"
	"reflect"
	"slices"
	"testing"

	"github.com/davecgh/go-spew/spew"
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
	var obj1, obj2 any
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

// AssertEqual checks if the given two values are the same.
// See also https://go.dev/wiki/TestComments
func AssertEqual[T comparable](t *testing.T, errReason string, want, got T) {
	t.Helper()
	if want == got {
		return
	}
	errReason += "\n"
	errReason += "-want: " + spew.Sdump(want)
	errReason += "+got: " + spew.Sdump(got)
	t.Error(errReason)
}

// AssertNotEqual checks if the given two values are not same.
// See also https://go.dev/wiki/TestComments
func AssertNotEqual[T comparable](t *testing.T, errReason string, check, got T) {
	t.Helper()
	if check != got {
		return
	}
	errReason += "\n"
	errReason += "-check: " + spew.Sdump(check)
	errReason += "+got: " + spew.Sdump(got)
	t.Error(errReason)
}

// AssertEqualSlice checks if the given two slices are the same using [slices.Equal].
// See also https://go.dev/wiki/TestComments
func AssertEqualSlice[S ~[]E, E comparable](t *testing.T, errReason string, want, got S) {
	t.Helper()
	if slices.Equal(want, got) {
		return
	}
	errReason += "\n"
	errReason += "-want: " + spew.Sdump(want)
	errReason += "+got: " + spew.Sdump(got)
	t.Error(errReason)
}

// AssertEqualMap checks if the given two maps are the same using [maps.Equal].
// errReason can be a expression for fmt.Sprintf(errReason, want, got).
// See also https://go.dev/wiki/TestComments
func AssertEqualMap[M1, M2 ~map[K]V, K, V comparable](t *testing.T, errReason string, want M1, got M2) {
	t.Helper()
	if maps.Equal(want, got) {
		return
	}
	errReason += "\n"
	errReason += "-want: " + spew.Sdump(want)
	errReason += "+got: " + spew.Sdump(got)
	t.Error(errReason)
}

// AssertEqualErr checks if the given two errors are the same.
// Errors are checked by following order and considered the same
// when one of them returned true.
//
//   - Compare pointer: want == got
//   - Compare error: errors.Is(got, want)
//   - Compare message: want.Error() == got.Error()
//
// See also https://go.dev/wiki/TestComments
func AssertEqualErr(t *testing.T, errReason string, want, got error) {
	t.Helper()
	if want == got {
		return // nil == nil is also here.
	}
	if errors.Is(got, want) {
		return
	}
	if want != nil && got != nil {
		if want.Error() == got.Error() {
			return
		}
	}
	errReason += "\n"
	errReason += "-want: " + spew.Sdump(want)
	errReason += "+got: " + spew.Sdump(got)
	t.Error(errReason)
}
