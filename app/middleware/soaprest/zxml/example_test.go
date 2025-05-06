// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package zxml_test

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"strings"

	"github.com/aileron-gateway/aileron-gateway/app/middleware/soaprest/zxml"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

func ExampleSimple_xmljson() {
	c := zxml.JSONConverter{
		EncodeDecoder: zxml.NewSimple(),
	}

	b, err := c.XMLtoJSON([]byte(`<alice charlie="david">bob</alice>`))
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))
	// Output:
	// {"alice":{"$":"bob","@charlie":"david"}}
}

func ExampleSimple_jsonxml() {
	c := zxml.JSONConverter{
		EncodeDecoder: zxml.NewSimple(),
	}

	b, err := c.JSONtoXML([]byte(`{"alice":{"$":"bob","@charlie":"david"}}`))
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))
	// Output:
	// <alice charlie="david">bob</alice>
}

func ExampleRayFish_xmljson() {
	c := zxml.JSONConverter{
		EncodeDecoder: zxml.NewRayFish(),
	}

	b, err := c.XMLtoJSON([]byte(`<alice charlie="david">bob</alice>`))
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))
	// Output:
	// {"#children":[{"#children":[],"#name":"@charlie","#text":"david"}],"#name":"alice","#text":"bob"}
}

func ExampleRayFish_jsonxml() {
	c := zxml.JSONConverter{
		EncodeDecoder: zxml.NewRayFish(),
	}

	b, err := c.JSONtoXML([]byte(`
{
  "#name": "alice",
  "#text": "bob",
  "#children": [
    {
      "#name": "@charlie",
      "#text": "david",
      "#children": []
    }
  ]
}
	`))
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))
	// Output:
	// <alice charlie="david">bob</alice>
}

func ExampleBadgerFish_xmljson() {
	c := zxml.JSONConverter{
		EncodeDecoder: zxml.NewBadgerFish(),
	}

	b, err := c.XMLtoJSON([]byte(`<alice charlie="david">bob</alice>`))
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))
	// Output:
	// {"alice":{"$":"bob","@charlie":"david"}}
}

func ExampleBadgerFish_jsonxml() {
	c := zxml.JSONConverter{
		EncodeDecoder: zxml.NewBadgerFish(),
	}

	b, err := c.JSONtoXML([]byte(`{"alice":{"$":"bob","@charlie":"david"}}`))
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))
	// Output:
	// <alice charlie="david">bob</alice>
}

func ExampleJSONConverter_WithXMLDecoderOpts() {
	c := zxml.JSONConverter{
		EncodeDecoder: zxml.NewSimple(),
	}
	c.WithXMLDecoderOpts(func(d *xml.Decoder) {
		d.CharsetReader = func(charset string, input io.Reader) (io.Reader, error) {
			switch {
			case strings.EqualFold(charset, "Shift_JIS"):
				return transform.NewReader(input, japanese.ShiftJIS.NewDecoder()), nil
			}
			return input, nil
		}
	})

	header := `<?xml version="1.0" encoding="SHIFT_JIS"?>`
	b, err := c.XMLtoJSON([]byte(header + "<alice>\xea\x9f</alice>"))
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))
	// Output:
	// {"alice":"堯"}
}

func ExampleJSONConverter_WithXMLEncoderOpts() {
	c := zxml.JSONConverter{
		EncodeDecoder: zxml.NewSimple(),
	}
	c.WithXMLEncoderOpts(func(e *xml.Encoder) {
		e.Indent("", "  ")
	})

	b, err := c.JSONtoXML([]byte(`{"alice":{"bob":{"$":"david"}}}`))
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))
	// Output:
	// <alice>
	//   <bob>david</bob>
	// </alice>
}

func ExampleJSONConverter_WithJSONEncoderOpts() {
	c := zxml.JSONConverter{
		EncodeDecoder: zxml.NewSimple(),
	}
	c.WithJSONEncoderOpts(func(e *json.Encoder) {
		e.SetIndent("", "  ")
	})

	b, err := c.XMLtoJSON([]byte(`<alice>bob</alice>`))
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))
	// Output:
	// {
	//   "alice": "bob"
	// }
}

func ExampleJSONConverter_WithJSONDecoderOpts() {
	c := zxml.JSONConverter{
		EncodeDecoder: zxml.NewSimple(),
	}
	c.WithJSONDecoderOpts(func(d *json.Decoder) {
		d.UseNumber()
	})

	b, err := c.JSONtoXML([]byte(`{"alice":{"$":1.2345678901234567890}}`))
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))
	// Output:
	// <alice>1.2345678901234567890</alice>
}
