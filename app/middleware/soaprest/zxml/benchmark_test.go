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

var (
	xmlFile  = "00_wiki.xml"
	jsonFile = "00_wiki.json"
)

func BenchmarkSimple_XMLtoJSON(b *testing.B) {
	c := zxml.JSONConverter{
		EncodeDecoder: zxml.NewSimple(),
		Header:        xml.Header,
	}
	c.WithJSONDecoderOpts(func(d *json.Decoder) { d.UseNumber() })
	c.WithJSONEncoderOpts(func(e *json.Encoder) { e.SetEscapeHTML(false) })

	data, err := os.ReadFile("./testdata/xml/" + xmlFile)
	if err != nil {
		panic(err)
	}

	b.ResetTimer()
	for b.Loop() {
		_, err := c.XMLtoJSON(data)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkSimple_JSONtoXML(b *testing.B) {
	c := zxml.JSONConverter{
		EncodeDecoder: zxml.NewSimple(),
		Header:        xml.Header,
	}
	c.WithJSONDecoderOpts(func(d *json.Decoder) { d.UseNumber() })
	c.WithJSONEncoderOpts(func(e *json.Encoder) { e.SetEscapeHTML(false) })

	data, err := os.ReadFile("./testdata/simple/" + jsonFile)
	if err != nil {
		panic(err)
	}

	b.ResetTimer()
	for b.Loop() {
		_, err := c.JSONtoXML(data)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkRayFish_XMLtoJSON(b *testing.B) {
	c := zxml.JSONConverter{
		EncodeDecoder: zxml.NewRayFish(),
		Header:        xml.Header,
	}
	c.WithJSONDecoderOpts(func(d *json.Decoder) { d.UseNumber() })
	c.WithJSONEncoderOpts(func(e *json.Encoder) { e.SetEscapeHTML(false) })

	data, err := os.ReadFile("./testdata/xml/" + xmlFile)
	if err != nil {
		panic(err)
	}

	b.ResetTimer()
	for b.Loop() {
		_, err := c.XMLtoJSON(data)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkRayFish_JSONtoXML(b *testing.B) {
	c := zxml.JSONConverter{
		EncodeDecoder: zxml.NewRayFish(),
		Header:        xml.Header,
	}
	c.WithJSONDecoderOpts(func(d *json.Decoder) { d.UseNumber() })
	c.WithJSONEncoderOpts(func(e *json.Encoder) { e.SetEscapeHTML(false) })

	data, err := os.ReadFile("./testdata/rayfish/" + jsonFile)
	if err != nil {
		panic(err)
	}

	b.ResetTimer()
	for b.Loop() {
		_, err := c.JSONtoXML(data)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkBadgerFish_XMLtoJSON(b *testing.B) {
	c := zxml.JSONConverter{
		EncodeDecoder: zxml.NewBadgerFish(),
		Header:        xml.Header,
	}
	c.WithJSONDecoderOpts(func(d *json.Decoder) { d.UseNumber() })
	c.WithJSONEncoderOpts(func(e *json.Encoder) { e.SetEscapeHTML(false) })

	data, err := os.ReadFile("./testdata/xml/" + xmlFile)
	if err != nil {
		panic(err)
	}

	b.ResetTimer()
	for b.Loop() {
		_, err := c.XMLtoJSON(data)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkBadgerFish_JSONtoXML(b *testing.B) {
	c := zxml.JSONConverter{
		EncodeDecoder: zxml.NewBadgerFish(),
		Header:        xml.Header,
	}
	c.WithJSONDecoderOpts(func(d *json.Decoder) { d.UseNumber() })
	c.WithJSONEncoderOpts(func(e *json.Encoder) { e.SetEscapeHTML(false) })

	data, err := os.ReadFile("./testdata/badgerfish/" + jsonFile)
	if err != nil {
		panic(err)
	}

	b.ResetTimer()
	for b.Loop() {
		_, err := c.JSONtoXML(data)
		if err != nil {
			b.Fatal(err)
		}
	}
}
