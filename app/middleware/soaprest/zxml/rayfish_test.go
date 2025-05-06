// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package zxml_test

import (
	"encoding/json"
	"encoding/xml"
	"os"
	"testing"

	"github.com/aileron-gateway/aileron-gateway/app/middleware/soaprest/zxml"
)

var rayFishDecodeTest = []struct {
	file string
	err  error
}{
	{"00_basic", nil},
	{"00_datatype", nil},
	{"00_multiple", nil},
	{"00_wiki", nil},
	{"ng_case01", &zxml.XMLError{Cause: zxml.CauseXMLDecoder}},
	{"ng_case02", &zxml.XMLError{Cause: zxml.CauseXMLDecoder}},
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
	{"soap_example", nil},
	{"soap11_basic_fault", nil},
	{"soap11_basic", nil},
	{"soap12_basic_fault", nil},
	{"soap12_basic", nil},
}

func TestRayFish_Decode(t *testing.T) {
	t.Parallel()

	c := zxml.JSONConverter{
		EncodeDecoder: zxml.NewRayFish(),
		Header:        xml.Header,
	}
	c.WithJSONDecoderOpts(func(d *json.Decoder) { d.UseNumber() })
	c.WithJSONEncoderOpts(func(e *json.Encoder) { e.SetEscapeHTML(false) })

	for _, tc := range rayFishDecodeTest {
		t.Run(tc.file, func(t *testing.T) {
			xmlBytes, _ := os.ReadFile("./testdata/xml/" + tc.file + ".xml")
			jsonBytes, _ := os.ReadFile("./testdata/rayfish/" + tc.file + ".json")
			b, err := c.XMLtoJSON(xmlBytes)
			AssertEqualErr(t, "error not match", tc.err, err)
			if tc.err != nil {
				return
			}
			if !equalJSON(t, b, jsonBytes) {
				t.Error("decode result not match (xml to json)")
			}
		})
	}
}

var rayFishEncodeTest = []struct {
	file string
	err  error
}{
	{"00_basic", nil},
	{"00_datatype", nil},
	{"00_multiple", nil},
	{"00_wiki", nil},
	{"ng_case01", &zxml.XMLError{Cause: zxml.CauseJSONDecoder}},
	{"ng_case02", &zxml.XMLError{Cause: zxml.CauseJSONStruct}},
	{"ng_case03", &zxml.XMLError{Cause: zxml.CauseJSONStruct}},
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
	{"soap_example", nil},
	{"soap11_basic_fault", nil},
	{"soap11_basic", nil},
	{"soap12_basic_fault", nil},
	{"soap12_basic", nil},
}

func TestRayFish_Encode(t *testing.T) {
	t.Parallel()

	c := zxml.JSONConverter{
		EncodeDecoder: zxml.NewRayFish(),
		Header:        xml.Header,
	}
	c.WithJSONDecoderOpts(func(d *json.Decoder) { d.UseNumber() })
	c.WithJSONEncoderOpts(func(e *json.Encoder) { e.SetEscapeHTML(false) })

	for _, tc := range rayFishEncodeTest {
		t.Run(tc.file, func(t *testing.T) {
			xmlBytes, _ := os.ReadFile("./testdata/rayfish/" + tc.file + ".xml")
			jsonBytes, _ := os.ReadFile("./testdata/rayfish/" + tc.file + ".json")
			b, err := c.JSONtoXML(jsonBytes)
			AssertEqualErr(t, "error not match", tc.err, err)
			if tc.err != nil {
				return
			}
			if !equalXML(t, b, xmlBytes) {
				t.Error("encode result not match (json to xml)")
			}
		})
	}
}
