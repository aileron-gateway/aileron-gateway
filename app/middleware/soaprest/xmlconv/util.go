// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package xmlconv

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type emptyValue any

var (
	EmptyNil emptyValue = nil
	EmptyMap emptyValue = make(map[string]any, 0)
	EmptyStr emptyValue = ""
)

// ErrCause is the error cause description
// occurred while conversions.
type ErrCause string

var (
	ErrDataType ErrCause = "xmlconv: invalid data type. "
	ErrJSON     ErrCause = "xmlconv: invalid json structure. "
	ErrJSONKey  ErrCause = "xmlconv: found invalid json key. "
)

func (c ErrCause) New(format string, args ...any) error {
	return errors.New(string(c) + fmt.Sprintf(format, args...))
}

// attrNameString convert XML attributes token name to string.
// Returned string is mainly used as JSON key.
//
// Examples (prefix="@", sep=":"):
//
//	┌────────────────────────┬─────────────────┐
//	│      XML element       │ Returned string │
//	├────────────────────────┼─────────────────┤
//	│ <elem foo="bar">       │ @foo            │
//	│------------------------│-----------------│
//	│ <elem ns:foo="bar">    │ @ns:foo         │
//	│------------------------│-----------------│
//	│ <elem xmlns:foo="bar"> │ @xmlns:foo      │
//	│------------------------│-----------------│
//	│ <elem xmlns="bar">     │ @xmlns          │
//	└────────────────────────┴─────────────────┘
func AttrNameString(name xml.Name, prefix, sep string, ns [][2]string) string {
	if name.Space == "" {
		return prefix + name.Local
	}
	for i := len(ns) - 1; i >= 0; i-- {
		if name.Space == ns[i][1] {
			return prefix + ns[i][0] + sep + name.Local
		}
	}
	return prefix + name.Space + sep + name.Local
}

// TokenName converts token name given as [encoding/xml.Name] to string
// with namespace consideration.
// It returns string with <name> or <namespace><sep><name> format.
// The <namespace> can be URI or alias name.
// Namespace URI and corresponding alias name can be given by the
// second argument ns.
// For example [2]string{"xsd","http://www.w3.org/2001/XMLSchema"}
// can be given in the ns.
//
// Why the ns should be given? The answer is the standard
// [encoding/xml.Name.Space] always contains namespace value
// which mostly URI. It does not have alias name information.
func TokenName(name xml.Name, sep string, ns [][2]string) string {
	if name.Space == "" {
		return name.Local
	}
	for i := len(ns) - 1; i >= 0; i-- {
		if name.Space == ns[i][1] {
			if ns[i][0] == "xmlns" {
				return name.Local
			}
			return ns[i][0] + sep + name.Local
		}
	}
	return name.Space + sep + name.Local
}

// ParseNamespace parses namespace from the given [encoding/xml.Attr].
func ParseNamespace(attrs []xml.Attr, ns [][2]string) [][2]string {
	for _, attr := range attrs {
		name := attr.Name
		switch name.Space {
		case "":
			if name.Local == "xmlns" {
				ns = append(ns, [2]string{"xmlns", attr.Value})
			}
		case "xmlns":
			ns = append(ns, [2]string{name.Local, attr.Value})
		default:
			ns = append(ns, [2]string{name.Local, attr.Value})
		}
	}
	return ns
}

// jsonValueToToken returns a [encoding/xml.Token] from the value
// parsed from JSON document.
// Supported types follow the [encoding/json.Token] which also listed below.
// Unsupported types results in an error.
//
// Supported Types:
//   - nil
//   - string
//   - bool
//   - float64
//   - json.Number
func jsonValueToToken(trimSpace bool, value any) (xml.Token, error) {
	switch v := value.(type) {
	case string:
		if trimSpace {
			v = strings.TrimSpace(v)
		}
		return xml.CharData(v), nil
	case bool:
		return xml.CharData(strconv.FormatBool(v)), nil
	case float64:
		return xml.CharData(strconv.FormatFloat(v, 'g', -1, 64)), nil
	case json.Number:
		return xml.CharData(v), nil
	case nil: // Keep empty
		return xml.CharData(nil), nil
	default:
		return nil, ErrDataType.New("must be null, string, bool or number. got %T:%+v", value, value)
	}
}

func restoreNamespace(nsSep string, key string) string {
	if nsSep == ":" || nsSep == "" {
		return key
	}
	before, after, found := strings.Cut(key, nsSep)
	if found {
		return before + ":" + after
	}
	return key
}
