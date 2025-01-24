// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package xmlconv

import (
	"encoding/xml"
	"fmt"
)

// NewRayFish returns a new instance of [RayFish]
// with default configuration.
func NewRayFish() *RayFish {
	return &RayFish{
		NameKey:      "#name",
		TextKey:      "#text",
		ChildrenKey:  "#children",
		AttrPrefix:   "@",
		NamespaceSep: ":",
		TrimSpace:    true,
		emptyVal:     "",
	}
}

// RayFish is the XML JSON converter following RayFish convention.
// Basic converting rules are shown in the tables below.
//
// References:
//
//   - https://www.onperl.org/blog/onperl/page/rayfish
//   - https://github.com/bramstein/xsltjson/
//   - https://wiki.open311.org/JSON_and_XML_Conversion/
//   - https://pypi.org/project/xmljson/
//
// Conversion Rules:
//
//	┌─────────────────────────────────────┐┌─────────────────────────────────────┐
//	│============= XML Input =============││============ JSON Output ============│
//	└─────────────────────────────────────┘└─────────────────────────────────────┘
//
//	| A simple element without attributes or children.
//	┌─────────────────────────────────────┐┌─────────────────────────────────────┐
//	│ <alice>bob</alice>                  ││ {                                   │
//	│                                     ││   "#name": "alice",                 │
//	│                                     ││   "#text": "bob",                   │
//	│                                     ││   "#children": []                   │
//	│                                     ││ }                                   │
//	└─────────────────────────────────────┘└─────────────────────────────────────┘
//
//	| An element with attribute. Prefix "@" is added to the name.
//	┌─────────────────────────────────────┐┌─────────────────────────────────────┐
//	│ <alice charlie="david">bob</alice>  ││ {                                   │
//	│                                     ││   "#name": "alice",                 │
//	│                                     ││   "#text": "bob",                   │
//	│                                     ││   "#children": [                    │
//	│                                     ││     {                               │
//	│                                     ││       "#name": "@charlie",          │
//	│                                     ││       "#text": "david",             │
//	│                                     ││       "#children": []               │
//	│                                     ││     }                               │
//	│                                     ││   ]                                 │
//	│                                     ││ }                                   │
//	└─────────────────────────────────────┘└─────────────────────────────────────┘
//
//	| An element with children. Children have different name.
//	┌─────────────────────────────────────┐┌─────────────────────────────────────┐
//	│ <alice>                             ││ {                                   │
//	│   <bob>charlie</bob>                ││   "#name": "alice",                 │
//	│   <david>edgar</david>              ││   "#text": "",                      │
//	│ </alice>                            ││   "#children": [                    │
//	│                                     ││     {                               │
//	│                                     ││       "#name": "bob",               │
//	│                                     ││       "#text": "charlie",           │
//	│                                     ││       "#children": []               │
//	│                                     ││     },                              │
//	│                                     ││     {                               │
//	│                                     ││       "#name": "david",             │
//	│                                     ││       "#text": "edgar",             │
//	│                                     ││       "#children": []               │
//	│                                     ││     }                               │
//	│                                     ││   ]                                 │
//	│                                     ││ }                                   │
//	└─────────────────────────────────────┘└─────────────────────────────────────┘
//
//	| An element with children. Children have the same name.
//	┌─────────────────────────────────────┐┌─────────────────────────────────────┐
//	│ <alice>                             ││ {                                   │
//	│   <bob>charlie</bob>                ││   "#name": "alice",                 │
//	│   <bob>edgar</bob>                  ││   "#text": "",                      │
//	│ </alice>                            ││   "#children": [                    │
//	│                                     ││     {                               │
//	│                                     ││       "#name": "bob",               │
//	│                                     ││       "#text": "charlie",           │
//	│                                     ││       "#children": []               │
//	│                                     ││     },                              │
//	│                                     ││     {                               │
//	│                                     ││       "#name": "bob",               │
//	│                                     ││       "#text": "edgar",             │
//	│                                     ││       "#children": []               │
//	│                                     ││     }                               │
//	│                                     ││   ]                                 │
//	│                                     ││ }                                   │
//	└─────────────────────────────────────┘└─────────────────────────────────────┘
//
//	| An element with mixed content. Text contents are joined.
//	┌─────────────────────────────────────┐┌─────────────────────────────────────┐
//	│ <alice>                             ││ {                                   │
//	│   bob                               ││   "#name": "alice",                 │
//	│   <charlie>david</charlie>          ││   "#text": "bobedgar",              │
//	│   edgar                             ││   "#children": [                    │
//	│ </alice>                            ││     {                               │
//	│                                     ││       "#name": "charlie",           │
//	│                                     ││       "#text": "david",             │
//	│                                     ││       "#children": []               │
//	│                                     ││     }                               │
//	│                                     ││   ]                                 │
//	│                                     ││ }                                   │
//	└─────────────────────────────────────┘└─────────────────────────────────────┘
//
//	| An element with namespaces.
//	┌─────────────────────────────────────┐┌─────────────────────────────────────┐
//	│ <alice xmlns:ns="http://abc.com/">  ││ {                                   │
//	│   <ns:bob>charlie</ns:bob>          ││   "#name": "alice",                 │
//	│ </alice>                            ││   "#text": "",                      │
//	│                                     ││   "#children": [                    │
//	│                                     ││     {                               │
//	│                                     ││       "#name": "@xmlns:ns",         │
//	│                                     ││       "#text": "http://abc.com/",   │
//	│                                     ││       "#children": []               │
//	│                                     ││     },                              │
//	│                                     ││     {                               │
//	│                                     ││       "#name": "ns:bob",            │
//	│                                     ││       "#text": "charlie",           │
//	│                                     ││       "#children": []               │
//	│                                     ││     }                               │
//	│                                     ││   ]                                 │
//	│                                     ││ }                                   │
//	└─────────────────────────────────────┘└─────────────────────────────────────┘
//
//	| An empty element with empty attribute.
//	┌─────────────────────────────────────┐┌─────────────────────────────────────┐
//	│ <alice charlie=""></alice>          ││ {                                   │
//	│                                     ││   "#name": "alice",                 │
//	│                                     ││   "#text": "",                      │
//	│                                     ││   "#children": [                    │
//	│                                     ││     {                               │
//	│                                     ││       "#name": "@charlie",          │
//	│                                     ││       "#text": "",                  │
//	│                                     ││       "#children": []               │
//	│                                     ││     }                               │
//	│                                     ││   ]                                 │
//	│                                     ││ }                                   │
//	└─────────────────────────────────────┘└─────────────────────────────────────┘
//
//	| [RayFish.WithEmptyValue] replaces the JSON value for empty elements.
//	| Default is the converter created with [NewRayFish].
//	┌─────────────────────────────────────┐┌─────────────────────────────────────┐
//	│ <!-- Default -->                    ││ {                                   │
//	│ <alice></alice>                     ││   "#name": "alice",                 │
//	│                                     ││   "#text": "",                      │
//	│                                     ││   "#children": []                   │
//	│                                     ││ }                                   │
//	└─────────────────────────────────────┘└─────────────────────────────────────┘
//	┌─────────────────────────────────────┐┌─────────────────────────────────────┐
//	│ <!-- WithEmptyValue(nil) -->        ││ {                                   │
//	│ <alice></alice>                     ││   "#name": "alice",                 │
//	│                                     ││   "#text": null,                    │
//	│                                     ││   "#children": []                   │
//	│                                     ││ }                                   │
//	└─────────────────────────────────────┘└─────────────────────────────────────┘
//	┌─────────────────────────────────────┐┌─────────────────────────────────────┐
//	│ <!-- m := map[string]any{} -->      ││ {                                   │
//	│ <!-- WithEmptyValue(m) -->          ││   "#name": "alice",                 │
//	│ <alice></alice>                     ││   "#text": {},                      │
//	│                                     ││   "#children": []                   │
//	│                                     ││ }                                   │
//	└─────────────────────────────────────┘└─────────────────────────────────────┘
type RayFish struct {
	// NameKeys is the json key name to store
	// XML element names.
	// Typically "#name" is used.
	// NameKey should not be empty.
	NameKey string
	// TextKey is the json key name to store
	// content of XML elements.
	// Typically "#text" is used.
	// TextKey should not be empty.
	TextKey string
	// ChildrenKey is the json key name to store
	// attributes of a element and its child elements.
	// Typically "#children" is used.
	// ChildrenKey should not be empty.
	ChildrenKey string
	// AttrPrefix is the json key name prefix for XML attributes.
	// Attribute names are stored in a json with this prefix.
	// For example, XML attribute foo="" is converted into {"@foo": "bar"}
	// Typically "@" is used.
	// AttrPrefix should not be empty.
	AttrPrefix string
	// NamespaceSep is the name space separator.
	// Namespace separator ":" in XML element names are converted
	// into the specified string.
	// Note that, general RayFish convention discard namespace information
	// but this converter keep it.
	// Use ":" if there is no reason to change.
	// NamespaceSep should not be empty.
	NamespaceSep string

	// XMLValue convert JSON value into XML value.
	// Input value is the any type value decoded by [json.Decoder].
	// XML attributes can be used but should not be modified.
	// Returned values are recognized by [xml.Encoder].
	// The given [xml.StartElement] can be modified in the function.
	// This option is used in JSON to XML conversion.
	XMLValue func(any, *xml.StartElement) (xml.Token, error)
	// JSONValue convert XML value into JSON value.
	// Input value is the text part of a XML element.
	// Returned values are recognized by [json.Encoder].
	// This function is used in XML to JSON conversion.
	JSONValue func(string, xml.StartElement) (any, error)

	// TrimSpace if true, trims unicode space from xml text.
	// See the [unicode.IsSpace] for space definition.
	// This option is used in XML to JSON conversion.
	TrimSpace bool
	// IgnoreUnknownKey if true, ignores undefined or unused
	// JSON attributes rather than error.
	// This option is used in JSON to XML conversion.
	IgnoreUnusedKey bool

	// emptyVal is the value for empty XML element
	// that have no attributes and no child elements.
	emptyVal any
}

// WithEmptyValue replaces the value corresponding to empty element
// that do not have any attributes, text content and child elements.
// The given value is used in XML to JSON conversion.
// Allowed values are listed below and other will result in panic.
//
// Allowed values are:
//
//   - nil
//   - string("")
//   - make(map[string]any,0)
func (r *RayFish) WithEmptyValue(v any) {
	switch v.(type) {
	case nil, string, map[string]any:
	default:
		panic(fmt.Errorf("xmlconv: %#v is not allowed for empty value", v))
	}
	r.emptyVal = v
}
