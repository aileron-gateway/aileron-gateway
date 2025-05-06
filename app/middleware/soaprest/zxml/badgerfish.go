// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package zxml

import (
	"encoding/xml"
	"fmt"
)

// NewBadgerFish returns a new instance of [BadgerFish]
// with default configuration.
func NewBadgerFish() *BadgerFish {
	return &BadgerFish{
		TextKey:      "$",
		AttrPrefix:   "@",
		NamespaceSep: ":",
		TrimSpace:    true,
		emptyVal:     make(map[string]any, 0),
	}
}

// BadgerFish is the XML JSON converter following BadgerFish encode-decoder.
// Basic conversion rules are shown in the table below.
//
// References:
//
//   - http://www.sklar.com/badgerfish/
//   - https://wiki.open311.org/JSON_and_XML_Conversion/
//   - http://dropbox.ashlock.us/open311/json-xml/
//   - https://github.com/bramstein/xsltjson/
//   - https://pypi.org/project/xmljson/
//   - https://cloud.google.com/apigee/docs/api-platform/reference/policies/xml-json-policy
//   - https://cloud.google.com/apigee/docs/api-platform/reference/policies/json-xml-policy
//
// Conversion Rules:
//
//	┌=====================================┐┌=====================================┐
//	│              XML Input              ││            JSON Output              |
//	└=====================================┘└=====================================┘
//
//	| A simple element without attributes or children.
//	┌─────────────────────────────────────┐┌─────────────────────────────────────┐
//	│ <alice>bob</alice>                  ││ {                                   │
//	│                                     ││   "alice": { "$": "bob" }           │
//	│                                     ││ }                                   │
//	└─────────────────────────────────────┘└─────────────────────────────────────┘
//
//	| An element with attribute. Prefix "@" is added to the name.
//	┌─────────────────────────────────────┐┌─────────────────────────────────────┐
//	│ <alice charlie="david">bob</alice>  ││ {                                   │
//	│                                     ││   "alice": {                        │
//	│                                     ││     "$": "bob",                     │
//	│                                     ││     "@charlie": "david"             │
//	│                                     ││   }                                 │
//	│                                     ││ }                                   │
//	└─────────────────────────────────────┘└─────────────────────────────────────┘
//
//	| An element with children. Children have different name.
//	┌─────────────────────────────────────┐┌─────────────────────────────────────┐
//	│ <alice>                             ││ {                                   │
//	│   <bob>charlie</bob>                ││   "alice": {                        │
//	│   <david>edgar</david>              ││     "bob": { "$": "charlie" },      │
//	│ </alice>                            ││     "david": { "$": "edgar" }       │
//	│                                     ││   }                                 │
//	│                                     ││ }                                   │
//	└─────────────────────────────────────┘└─────────────────────────────────────┘
//
//	| An element with children. Children have the same name.
//	┌─────────────────────────────────────┐┌─────────────────────────────────────┐
//	│ <alice>                             ││ {                                   │
//	│   <bob>charlie</bob>                ││   "alice": {                        │
//	│   <bob>edgar</bob>                  ││     "bob": [                        │
//	│ </alice>                            ││       { "$": "charlie" },           │
//	│                                     ││       { "$": "edgar" }              │
//	│                                     ││     ]                               │
//	│                                     ││   }                                 │
//	│                                     ││ }                                   │
//	└─────────────────────────────────────┘└─────────────────────────────────────┘
//
//	| An element with mixed content. Text contents are joined.
//	┌─────────────────────────────────────┐┌─────────────────────────────────────┐
//	│ <alice>                             ││ {                                   │
//	│   bob                               ││   "alice": {                        │
//	│   <charlie>david</charlie>          ││     "$": "bobedgar",                │
//	│   edgar                             ││     "charlie": { "$": "david" }     │
//	│ </alice>                            ││   }                                 │
//	│                                     ││ }                                   │
//	└─────────────────────────────────────┘└─────────────────────────────────────┘
//
//	| An element with namespaces.
//	┌─────────────────────────────────────┐┌─────────────────────────────────────┐
//	│ <alice xmlns:ns="http://abc.com/">  ││ {                                   │
//	│   <ns:bob>charlie</ns:bob>          ││   "alice": {                        │
//	│ </alice>                            ││     "@xmlns": {                     │
//	│                                     ││       "ns": "http://abc.com/"       │
//	│                                     ││     },                              │
//	│                                     ││     "ns:bob": { "$": "charlie" }    │
//	│                                     ││   }                                 │
//	│                                     ││ }                                   │
//	└─────────────────────────────────────┘└─────────────────────────────────────┘
//
//	| An empty element with empty attribute.
//	┌─────────────────────────────────────┐┌─────────────────────────────────────┐
//	│ <alice charlie=""></alice>          ││ {                                   │
//	│                                     ││   "alice": {                        │
//	│                                     ││     "@charlie": ""                  │
//	│                                     ││   }                                 │
//	│                                     ││ }                                   │
//	└─────────────────────────────────────┘└─────────────────────────────────────┘
//
//	| [BadgerFish.WithEmptyValue] replaces the JSON value for empty elements.
//	┌─────────────────────────────────────┐┌─────────────────────────────────────┐
//	│ <!-- m := map[string]any{} -->      ││ {                                   │
//	│ <!-- WithEmptyValue(m) -->          ││   "alice": {}                       │
//	│ <alice></alice>                     ││ }                                   │
//	└─────────────────────────────────────┘└─────────────────────────────────────┘
//	┌─────────────────────────────────────┐┌─────────────────────────────────────┐
//	│ <!-- WithEmptyValue("") -->         ││ {                                   │
//	│ <alice></alice>                     ││   "alice": ""                       │
//	│                                     ││ }                                   │
//	└─────────────────────────────────────┘└─────────────────────────────────────┘
//	┌─────────────────────────────────────┐┌─────────────────────────────────────┐
//	│ <!-- WithEmptyValue(nil) -->        ││ {                                   │
//	│ <alice></alice>                     ││   "alice": null                     │
//	│                                     ││ }                                   │
//	└─────────────────────────────────────┘└─────────────────────────────────────┘
type BadgerFish struct {
	// TextKey is the json key name to store text content of elements.
	// Typically "$" is used.
	// TextKey should not be empty.
	TextKey string
	// AttrPrefix is the json key name prefix for XML attributes.
	// Attribute names are stored in json with this prefix.
	// For example, XML attribute foo="bar" is converted into {"@foo": "bar"}.
	// Typically "@" is used.
	// AttrPrefix should not be empty.
	AttrPrefix string
	// NamespaceSep is the name space separator.
	// Namespace separator ":" in XML element names are converted
	// into the specified string.
	// Use ":" if there is no reason to change.
	// NamespaceSep should not be empty.
	NamespaceSep string
	// TrimSpace if true, trims unicode space from xml text.
	// See the [unicode.IsSpace] for space definition.
	// This option is used in XML to JSON conversion.
	TrimSpace bool

	// XMLValue convert JSON value into XML value.
	// Input value is the any type value decoded by [json.Decoder].
	// Returned value is recognized by [xml.Encoder].
	// The given [xml.StartElement] can be modified in the function.
	// This option is used in JSON to XML conversion.
	XMLValue func(any, *xml.StartElement) (xml.Token, error)
	// JSONValue convert XML value into JSON value.
	// Input value is the text part of a XML element.
	// Returned value is recognized by [json.Encoder].
	// This function is used in XML to JSON conversion.
	JSONValue func(string, xml.StartElement) (any, error)

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
func (b *BadgerFish) WithEmptyValue(v any) {
	switch v.(type) {
	case nil, string, map[string]any:
	default:
		panic(&XMLError{Cause: CauseEmptyVal, Detail: fmt.Sprintf("%T:%+v", v, v)})
	}
	b.emptyVal = v
}
