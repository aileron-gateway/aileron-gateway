// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package xmlconv_test

import (
	"encoding/json"
	"encoding/xml"
	"os"
	"testing"

	"github.com/aileron-gateway/aileron-gateway/app/middleware/soaprest/xmlconv"
)

var simpleDecodeTest = []struct {
	file string
	err  error
}{
	{"00_basic", nil},
	{"00_datatype", nil},
	{"00_multiple", nil},
	{"00_wiki", nil},
	{"ng_case01", &xml.SyntaxError{Msg: "unexpected EOF", Line: 2}},
	{"ng_case02", &xml.SyntaxError{Msg: "expected element name after <", Line: 2}},
	{"ok_case01", nil},
	{"ok_case02", nil},
	{"ok_case03", nil},
	{"ok_case04", nil},
	{"ok_case05", nil},
	{"ok_case06", nil},
	{"ok_case07", nil},
	{"ok_case08", nil},
	{"ok_case09", nil},
	{"ok_case10", nil},
	{"ok_case11", nil},
	{"ok_case12", nil},
	{"ok_case13", nil},
	{"ok_case14", nil},
	{"ok_case15", nil},
	{"ok_case16", nil},
	{"ok_case17", nil},
	{"ok_case18", nil},
	{"ok_case19", nil},
	{"ok_case20", nil},
	{"ok_case21", nil},
}

func TestSimple_Decode(t *testing.T) {
	t.Parallel()

	ed := xmlconv.NewSimple()
	ed.PreferShort = false
	c := xmlconv.Converter{
		EncodeDecoder: ed,
		Header:        xml.Header,
	}
	c.WithJSONDecoderOpts(func(d *json.Decoder) { d.UseNumber() })
	c.WithJSONEncoderOpts(func(e *json.Encoder) { e.SetEscapeHTML(false) })

	for _, tc := range simpleDecodeTest {
		t.Run(tc.file, func(t *testing.T) {
			xmlBytes, _ := os.ReadFile("./testdata/xml/" + tc.file + ".xml")
			jsonBytes, _ := os.ReadFile("./testdata/simple/" + tc.file + ".json")
			b, err := c.XMLtoJSON(xmlBytes)
			if tc.err != nil {
				if tc.err.Error() != err.Error() {
					t.Errorf("error not match. want:%v got:%v", tc.err, err)
				}
				return
			} else {
				if err != nil {
					t.Errorf("error not match. want:%v got:%v", tc.err, err)
				}
			}
			if !equalJSON(t, b, jsonBytes) {
				t.Error("decode result not match (xml to json)")
			}
		})
	}
}

func TestSimple_Decode_short(t *testing.T) {
	t.Parallel()

	c := xmlconv.Converter{
		EncodeDecoder: xmlconv.NewSimple(),
		Header:        xml.Header,
	}
	c.WithJSONDecoderOpts(func(d *json.Decoder) { d.UseNumber() })
	c.WithJSONEncoderOpts(func(e *json.Encoder) { e.SetEscapeHTML(false) })

	for _, tc := range simpleDecodeTest {
		t.Run(tc.file, func(t *testing.T) {
			xmlBytes, _ := os.ReadFile("./testdata/xml/" + tc.file + ".xml")
			jsonBytes, _ := os.ReadFile("./testdata/simple/" + tc.file + "_short.json")
			b, err := c.XMLtoJSON(xmlBytes)
			if tc.err != nil {
				if tc.err.Error() != err.Error() {
					t.Errorf("error not match. want:%v got:%v", tc.err, err)
				}
				return
			} else {
				if err != nil {
					t.Errorf("error not match. want:%v got:%v", tc.err, err)
				}
			}
			if !equalJSON(t, b, jsonBytes) {
				t.Error("decode result not match (xml to json)")
			}
		})
	}
}

var simpleEncodeTest = []struct {
	file string
	err  error
}{
	{"00_basic", nil},
	{"00_datatype", nil},
	{"00_multiple", nil},
	{"00_wiki", nil},
	{"ok_case01", nil},
	{"ok_case02", nil},
	{"ok_case03", nil},
	{"ok_case04", nil},
	{"ok_case05", nil},
	{"ok_case06", nil},
	{"ok_case07", nil},
	{"ok_case08", nil},
	{"ok_case09", nil},
	{"ok_case10", nil},
	{"ok_case11", nil},
	{"ok_case12", nil},
	{"ok_case13", nil},
	{"ok_case14", nil},
	{"ok_case15", nil},
	{"ok_case16", nil},
	{"ok_case17", nil},
	{"ok_case18", nil},
	{"ok_case19", nil},
	{"ok_case20", nil},
	{"ok_case21", nil},
}

func TestSimple_Encode(t *testing.T) {
	t.Parallel()

	ed := xmlconv.NewSimple()
	ed.PreferShort = false
	c := xmlconv.Converter{
		EncodeDecoder: ed,
		Header:        xml.Header,
	}
	c.WithJSONDecoderOpts(func(d *json.Decoder) { d.UseNumber() })
	c.WithJSONEncoderOpts(func(e *json.Encoder) { e.SetEscapeHTML(false) })

	for _, tc := range simpleEncodeTest {
		t.Run(tc.file, func(t *testing.T) {
			xmlBytes, _ := os.ReadFile("./testdata/simple/" + tc.file + ".xml")
			jsonBytes, _ := os.ReadFile("./testdata/simple/" + tc.file + ".json")
			b, err := c.JSONtoXML(jsonBytes)
			if tc.err != nil {
				if tc.err.Error() != err.Error() {
					t.Errorf("error not match. want:%v got:%v", tc.err, err)
				}
				return
			} else {
				if err != nil {
					t.Errorf("error not match. want:%v got:%v", tc.err, err)
				}
			}
			if !equalXML(t, b, xmlBytes) {
				t.Error("encode result not match (json to xml)")
			}
		})
	}
}

func TestSimple_Encode_short(t *testing.T) {
	t.Parallel()

	c := xmlconv.Converter{
		EncodeDecoder: xmlconv.NewSimple(),
		Header:        xml.Header,
	}
	c.WithJSONDecoderOpts(func(d *json.Decoder) { d.UseNumber() })
	c.WithJSONEncoderOpts(func(e *json.Encoder) { e.SetEscapeHTML(false) })

	for _, tc := range simpleEncodeTest {
		t.Run(tc.file, func(t *testing.T) {
			xmlBytes, _ := os.ReadFile("./testdata/simple/" + tc.file + ".xml")
			jsonBytes, _ := os.ReadFile("./testdata/simple/" + tc.file + "_short.json")
			b, err := c.JSONtoXML(jsonBytes)
			if tc.err != nil {
				if tc.err.Error() != err.Error() {
					t.Errorf("error not match. want:%v got:%v", tc.err, err)
				}
				return
			} else {
				if err != nil {
					t.Errorf("error not match. want:%v got:%v", tc.err, err)
				}
			}
			if !equalXML(t, b, xmlBytes) {
				t.Error("encode result not match (json to xml)")
			}
		})
	}
}

// func Example_simple_dec() {
// 	ed := xmlconv.NewSimple()
// 	ed.PreferShort = false
// 	c := xmlconv.Converter{
// 		EncodeDecoder: ed,
// 		Header:        xml.Header,
// 	}
// 	cs := xmlconv.Converter{
// 		EncodeDecoder: xmlconv.NewSimple(),
// 		Header:        xml.Header,
// 	}
// 	b, _ := os.ReadFile("./testdata/xml/00_wiki.xml")
// 	bb, err := c.XMLtoJSON(b)
// 	if err != nil {
// 		panic(err)
// 	}
// 	fmt.Println(string(bb))

// 	bbb, err := cs.XMLtoJSON(b)
// 	if err != nil {
// 		panic(err)
// 	}
// 	fmt.Println(string(bbb))
// 	// Output:
// }

// func Example_simple_enc() {
// 	ed := xmlconv.NewSimple()
// 	ed.PreferShort = false
// 	c := xmlconv.Converter{
// 		EncodeDecoder: ed,
// 		Header:        xml.Header,
// 	}
// 	b, _ := os.ReadFile("./testdata/simple/00_wiki.json")
// 	bb, err := c.JSONtoXML(b)
// 	if err != nil {
// 		panic(err)
// 	}
// 	fmt.Println(string(bb))
// 	// Output:
// }

// func ExampleSimple_sample() {
// 	conv := xmlconv.NewSimple()
// 	conv.WithEmptyValue(map[string]any{})
// 	c := xmlconv.Converter{
// 		EncodeDecoder: conv,
// 		Header:        xml.Header,
// 	}
// 	b := `
// <alice></alice>
// 	`
// 	bb, err := c.XMLtoJSON([]byte(b))
// 	if err != nil {
// 		panic(err)
// 	}
// 	fmt.Println(string(bb))
// 	// Output:
// }
