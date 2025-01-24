// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package xmlconv

import (
	"encoding/xml"
	"fmt"
)

// NewSimple returns a new instance of [Simple]
// with default configuration.
func NewSimple() *Simple {
	return &Simple{
		TextKey:      "$",
		AttrPrefix:   "@",
		NamespaceSep: ":",
		TrimSpace:    true,
		PreferShort:  true,
		emptyVal:     "",
	}
}

// Simple is the simple XML JSON converter.
// It is not an implementation of convention such as [RayFish] or [BadgerFish].
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
//	│                                     ││   "alice": {                        │
//	│                                     ││     "$": "bob"                      │
//	│                                     ││   }                                 │
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
//	│   <david>edgar</david>              ││     "bob": {                        │
//	│ </alice>                            ││       "$": "charlie"                │
//	│                                     ││     },                              │
//	│                                     ││     "david": {                      │
//	│                                     ││       "$": "edgar"                  │
//	│                                     ││     }                               │
//	│                                     ││   }                                 │
//	│                                     ││ }                                   │
//	└─────────────────────────────────────┘└─────────────────────────────────────┘
//
//	| An element with children. Children have the same name.
//	┌─────────────────────────────────────┐┌─────────────────────────────────────┐
//	│ <alice>                             ││ {                                   │
//	│   <bob>charlie</bob>                ││   "alice": {                        │
//	│   <bob>edgar</bob>                  ││     "bob": [                        │
//	│ </alice>                            ││       {                             │
//	│                                     ││         "$": "charlie"              │
//	│                                     ││       },                            │
//	│                                     ││       {                             │
//	│                                     ││         "$": "edgar"                │
//	│                                     ││       }                             │
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
//	│   edgar                             ││     "charlie": {                    │
//	│ </alice>                            ││       "$": "david"                  │
//	│                                     ││     }                               │
//	│                                     ││   }                                 │
//	│                                     ││ }                                   │
//	└─────────────────────────────────────┘└─────────────────────────────────────┘
//
//	| An element with namespaces.
//	┌─────────────────────────────────────┐┌─────────────────────────────────────┐
//	│ <alice xmlns:ns="http://abc.com/">  ││ {                                   │
//	│   <ns:bob>charlie</ns:bob>          ││   "alice": {                        │
//	│ </alice>                            ││     "@xmlns:ns": "http://abc.com/", │
//	│                                     ││     "ns:bob": {                     │
//	│                                     ││       "$": "charlie"                │
//	│                                     ││     }                               │
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
//	| Option [Simple.PreferShort] directory associates text value to the key.
//	| This option works for the element without any child elements and attributes.
//	┌─────────────────────────────────────┐┌─────────────────────────────────────┐
//	│ <alice>bob</alice>                  ││ {                                   │
//	│                                     ││   "alice": "bob"                    │
//	│                                     ││ }                                   │
//	└─────────────────────────────────────┘└─────────────────────────────────────┘
//
//	| [Simple.WithEmptyValue] replaces the JSON value for empty elements.
//	| Default is the converter created with [NewSimple].
//	┌─────────────────────────────────────┐┌─────────────────────────────────────┐
//	│ <!-- Default -->                    ││ {                                   │
//	│ <alice></alice>                     ││   "alice": ""                       │
//	│                                     ││ }                                   │
//	└─────────────────────────────────────┘└─────────────────────────────────────┘
//	┌─────────────────────────────────────┐┌─────────────────────────────────────┐
//	│ <!-- WithEmptyValue(nil) -->        ││ {                                   │
//	│ <alice></alice>                     ││   "alice": null                     │
//	│                                     ││ }                                   │
//	└─────────────────────────────────────┘└─────────────────────────────────────┘
//	┌─────────────────────────────────────┐┌─────────────────────────────────────┐
//	│ <!-- m := map[string]any{} -->      ││ {                                   │
//	│ <!-- WithEmptyValue(m) -->          ││   "alice": {}                       │
//	│ <alice></alice>                     ││ }                                   │
//	└─────────────────────────────────────┘└─────────────────────────────────────┘
type Simple struct {
	// TextKey is the json key name to store
	// content of XML elements.
	// Typically "$" is used.
	// TextKey should not be empty.
	TextKey string
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
	// PreferShort if true, use short format.
	// For XML to JSON conversion, if content has no attribute and no child elements,
	// JSON will be {"key": "value"} rather than {"key": {"$": "value"}}.
	// For JSON to XML conversion, JSON can always use {"key": "value"}
	// and {"key": {"$": "value"}} expression without this configuration.
	// This option is used in XML to JSON conversion.
	PreferShort bool
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
func (s *Simple) WithEmptyValue(v any) {
	switch v.(type) {
	case nil, string, map[string]any:
	default:
		panic(fmt.Errorf("xmlconv: %#v is not allowed for empty value", v))
	}
	s.emptyVal = v
}
