package zxml

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"reflect"
	"testing"
)

func TestXMLError_Error(t *testing.T) {
	t.Parallel()
	t.Run("without inner error", func(t *testing.T) {
		err := &XMLError{
			Cause:  "cause",
			Detail: "detail",
		}
		msg := err.Error()
		AssertEqual(t, "error message not match", "zxml: cause detail", msg)
	})
	t.Run("with inner error", func(t *testing.T) {
		err := &XMLError{
			Err:    io.EOF,
			Cause:  "cause",
			Detail: "detail",
		}
		msg := err.Error()
		AssertEqual(t, "error message not match", "zxml: cause detail [EOF]", msg)
	})
}

func TestXMLError_Unwrap(t *testing.T) {
	t.Parallel()
	err := &XMLError{Err: io.EOF, Cause: "cause"}
	inner := err.Unwrap()
	AssertEqualErr(t, "errors not match", io.EOF, inner)
}

func TestXMLError_Is(t *testing.T) {
	t.Parallel()
	t.Run("same", func(t *testing.T) {
		err1 := &XMLError{Cause: "cause"}
		err2 := &XMLError{Cause: "cause"}
		same := err1.Is(err2)
		AssertEqual(t, "errors should be the same", true, same)
	})
	t.Run("same after unwrap", func(t *testing.T) {
		err1 := &XMLError{Cause: "cause"}
		err2 := fmt.Errorf("outer error [%w]", &XMLError{Cause: "cause"})
		same := err1.Is(err2)
		AssertEqual(t, "errors should be the same", true, same)
	})
	t.Run("not same", func(t *testing.T) {
		err1 := &XMLError{Cause: "cause"}
		err2 := io.EOF
		same := err1.Is(err2)
		AssertEqual(t, "errors should not be the same", false, same)
	})
	t.Run("different cause", func(t *testing.T) {
		err1 := &XMLError{Cause: "cause1"}
		err2 := &XMLError{Cause: "cause2"}
		same := err1.Is(err2)
		AssertEqual(t, "errors should not be the same", false, same)
	})
	t.Run("nil", func(t *testing.T) {
		err1 := &XMLError{Cause: "cause"}
		err2 := error(nil)
		same := err1.Is(err2)
		AssertEqual(t, "errors should not be the same", false, same)
	})
}

func TestParseNamespace(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		attrs []xml.Attr
		ns    [][2]string
	}{
		"case01": {
			attrs: []xml.Attr{{Name: xml.Name{Space: "", Local: "foo"}, Value: "test"}},
			ns:    nil,
		},
		"case02": {
			attrs: []xml.Attr{{Name: xml.Name{Space: "", Local: "xmlns"}, Value: "test"}},
			ns:    [][2]string{{"xmlns", "test"}},
		},
		"case03": {
			attrs: []xml.Attr{{Name: xml.Name{Space: "xmlns", Local: ""}, Value: "test"}},
			ns:    [][2]string{{"", "test"}},
		},
		"case04": {
			attrs: []xml.Attr{{Name: xml.Name{Space: "xmlns", Local: "ns"}, Value: "test"}},
			ns:    [][2]string{{"ns", "test"}},
		},
		"case05": {
			attrs: []xml.Attr{{Name: xml.Name{Space: "ns", Local: "foo"}, Value: "test"}},
			ns:    [][2]string{{"foo", "test"}},
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			ns := parseNamespace(tc.attrs, nil)
			AssertEqualSlice(t, "namespace not match", tc.ns, ns)
		})
	}
}

func TestAttrName(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		name        xml.Name
		prefix, sep string
		ns          [][2]string
		out         string
	}{
		"case01": {xml.Name{Space: "", Local: "foo"}, "@", ":", nil, "@foo"},
		"case02": {xml.Name{Space: "ns", Local: "foo"}, "@", ":", nil, "@ns:foo"},
		"case03": {xml.Name{Space: "ns", Local: "foo"}, "@", "_", nil, "@ns_foo"},
		"case04": {xml.Name{Space: "alias", Local: "foo"}, "@", ":", [][2]string{{"xmlns", "alias"}}, "@xmlns:foo"},
		"case05": {xml.Name{Space: "alias", Local: "foo"}, "@", ":", [][2]string{{"ns", "alias"}}, "@ns:foo"},
		"case06": {xml.Name{Space: "alias", Local: "foo"}, "@", ":", [][2]string{{"ns1", "alias"}, {"ns2", "alias"}}, "@ns2:foo"},
		"case07": {xml.Name{Space: "http://www.w3.org/XML/1998/namespace", Local: "foo"}, "@", ":", nil, "@xml:foo"},
		"case08": {xml.Name{Space: "http://www.w3.org/XML/1998/namespace", Local: "foo"}, "@", "_", nil, "@xml_foo"},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			out := attrName(tc.name, tc.prefix, tc.sep, tc.ns)
			AssertEqual(t, "result not match", tc.out, out)
		})
	}
}

func TestTokenName(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		name xml.Name
		sep  string
		ns   [][2]string
		out  string
	}{
		"case01": {xml.Name{Space: "", Local: "foo"}, ":", nil, "foo"},
		"case02": {xml.Name{Space: "ns", Local: "foo"}, ":", nil, "ns:foo"},
		"case03": {xml.Name{Space: "ns", Local: "foo"}, "_", nil, "ns_foo"},
		"case04": {xml.Name{Space: "alias", Local: "foo"}, ":", [][2]string{{"xmlns", "alias"}}, "foo"},
		"case05": {xml.Name{Space: "alias", Local: "foo"}, ":", [][2]string{{"ns", "alias"}}, "ns:foo"},
		"case06": {xml.Name{Space: "alias", Local: "foo"}, ":", [][2]string{{"ns1", "alias"}, {"ns2", "alias"}}, "ns2:foo"},
		"case07": {xml.Name{Space: "http://www.w3.org/XML/1998/namespace", Local: "foo"}, ":", nil, "xml:foo"},
		"case08": {xml.Name{Space: "http://www.w3.org/XML/1998/namespace", Local: "foo"}, "_", nil, "xml_foo"},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			out := tokenName(tc.name, tc.sep, tc.ns)
			AssertEqual(t, "result not match", tc.out, out)
		})
	}
}

func TestRestoreNamespace(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		sep string
		key string
		out string
	}{
		"case01": {"", "", ""},
		"case02": {"", "foo", "foo"},
		"case03": {"", "ns:foo", "ns:foo"},
		"case04": {"", "ns_foo", "ns_foo"},
		"case05": {":", "", ""},
		"case06": {":", "foo", "foo"},
		"case07": {":", ":foo", ":foo"},
		"case08": {":", "ns:foo", "ns:foo"},
		"case09": {"_", "", ""},
		"case10": {"_", "foo", "foo"},
		"case11": {"_", "_foo", "_foo"},
		"case12": {"_", "ns_foo", "ns:foo"},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			out := restoreNamespace(tc.sep, tc.key)
			AssertEqual(t, "output not match", tc.out, out)
		})
	}
}

func TestJsonValueToToken(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		trim  bool
		value any
		token xml.Token
		err   error
	}{
		"case01": {false, nil, xml.CharData([]byte(nil)), nil},
		"case02": {false, "test", xml.CharData([]byte("test")), nil},
		"case03": {true, " test ", xml.CharData([]byte("test")), nil},
		"case04": {false, json.Number("123"), xml.CharData([]byte("123")), nil},
		"case05": {false, float64(1.2345), xml.CharData([]byte("1.2345")), nil},
		"case06": {false, true, xml.CharData([]byte("true")), nil},
		"case07": {false, false, xml.CharData([]byte("false")), nil},
		"case08": {false, 123, nil, &XMLError{Cause: CauseDataType}},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			token, err := jsonValueToToken(tc.trim, tc.value)
			AssertEqual(t, "token not match", true, reflect.DeepEqual(tc.token, token))
			AssertEqualErr(t, "error should be nil", tc.err, err)
		})
	}
}
