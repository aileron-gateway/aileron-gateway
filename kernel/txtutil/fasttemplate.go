// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package txtutil

import (
	"bytes"
	"fmt"
	"io"
	"regexp"
	"slices"
	"strconv"
)

const (
	fastTplValue = iota
	fastTplTag
)

// NewFastTemplate returns a new instance of fast template.
// Allowed tag pattern is `[0-9a-zA-Z_\-\.]+`.
// Note that the fasttemplate is inspired by
// the https://github.com/valyala/fasttemplate.
func NewFastTemplate(tpl string, start, end string) *FastTemplate {
	exp := []byte(tpl)
	reg := regexp.MustCompile(regexp.QuoteMeta(start) + ` *[0-9a-zA-Z_\-\.]+ *` + regexp.QuoteMeta(end))
	indexes := reg.FindAllIndex(exp, -1)

	types := []int{}
	values := [][]byte{}

	pos := 0
	for _, ids := range indexes {
		if ids[0] > pos {
			types = append(types, fastTplValue)
			values = append(values, exp[pos:ids[0]])
		}
		types = append(types, fastTplTag)
		values = append(values, bytes.Trim(exp[ids[0]+len(start):ids[1]-len(end)], " "))
		pos = ids[1]
	}
	if pos < len(exp) {
		types = append(types, 0)
		values = append(values, exp[pos:])
	}

	return &FastTemplate{
		valTypes: slices.Clip(types),
		values:   slices.Clip(values),
	}
}

// FastTemplate is the fast template object.
// Use NewFastTemplate to create a new instance of FastTemplate.
// Note that the fasttemplate is inspired by
// the https://github.com/valyala/fasttemplate.
type FastTemplate struct {
	valTypes []int
	values   [][]byte
}

func (t *FastTemplate) Execute(m map[string]any) []byte {
	var buf bytes.Buffer
	t.ExecuteWriter(&buf, m)
	return buf.Bytes()
}

func (t *FastTemplate) ExecuteFunc(f func(string) []byte) []byte {
	var buf bytes.Buffer
	t.ExecuteFuncWriter(&buf, f)
	return buf.Bytes()
}

func (t *FastTemplate) ExecuteWriter(w io.Writer, m map[string]any) {
	mv := mapVal(m)
	for i := range t.valTypes {
		switch t.valTypes[i] {
		case fastTplValue:
			_, _ = w.Write(t.values[i])
		case fastTplTag:
			_, _ = w.Write(mv.Value(string(t.values[i])))
		}
	}
}

func (t *FastTemplate) ExecuteFuncWriter(w io.Writer, f func(string) []byte) {
	for i := range t.valTypes {
		switch t.valTypes[i] {
		case fastTplValue:
			_, _ = w.Write(t.values[i])
		case fastTplTag:
			_, _ = w.Write(f(string(t.values[i])))
		}
	}
}

type mapVal map[string]any

func (m mapVal) Value(tag string) []byte {
	if m == nil {
		return nil
	}
	val, ok := m[tag]
	if !ok {
		return nil
	}
	switch v := val.(type) {
	case nil:
		return []byte("<nil>")
	case string:
		return []byte(v)
	case []byte:
		return v
	case int:
		return []byte(strconv.FormatInt(int64(v), 10))
	case int32:
		return []byte(strconv.FormatInt(int64(v), 10))
	case int64:
		return []byte(strconv.FormatInt(v, 10))
	case float32:
		return []byte(strconv.FormatFloat(float64(v), 'f', -1, 32))
	case float64:
		return []byte(strconv.FormatFloat(float64(v), 'f', -1, 64))
	default:
		return []byte(fmt.Sprint(v)) // "%+v"
	}
}
