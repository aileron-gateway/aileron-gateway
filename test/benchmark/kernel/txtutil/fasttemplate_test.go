// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package txtutil_test

import (
	"bytes"
	"fmt"
	"sync"
	"testing"

	"github.com/aileron-gateway/aileron-gateway/kernel/txtutil"
)

var fastTplIn = map[string]any{
	"name":     "Jon Doe",
	"age":      "20",
	"gender":   "male",
	"birthday": "Dec 31, 2020",
	"parents": map[string]any{
		"mother": "alice",
		"father": "bob",
	},
}

var fastTplFormat = `
He is {{name}} and {{age}} years old.
His birthday is {{birthday}}.
His parents are {{parents}}.
`

func ExampleFastTemplate_Execute() {

	format := `
He is {{name}} and {{age}} years old.
His parents are {{parents}}.
`
	tpl := txtutil.NewFastTemplate(format, "{{", "}}")

	fmt.Println(string(tpl.Execute(fastTplIn)))
	// Output:
	// He is Jon Doe and 20 years old.
	// His parents are map[father:bob mother:alice].

}

func BenchmarkFastTemplate_Execute(b *testing.B) {

	tpl := txtutil.NewFastTemplate(fastTplFormat, "{{", "}}")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tpl.Execute(fastTplIn)
	}

}

func BenchmarkFastTemplate_ExecuteFunc(b *testing.B) {

	tpl := txtutil.NewFastTemplate(fastTplFormat, "{{", "}}")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tpl.ExecuteFunc(func(tag string) []byte {
			return []byte(fmt.Sprint(fastTplIn[tag]))
		})
	}

}

var pool = &sync.Pool{
	New: func() any {
		var buf bytes.Buffer
		return &buf
	},
}

func BenchmarkFastTemplate_ExecuteBuf(b *testing.B) {

	tpl := txtutil.NewFastTemplate(fastTplFormat, "{{", "}}")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf := pool.Get().(*bytes.Buffer)
		buf.Reset()
		tpl.ExecuteWriter(buf, fastTplIn)
		_ = buf.Bytes()
		pool.Put(buf)
	}

}

func BenchmarkFastTemplate_ExecuteFuncBuf(b *testing.B) {

	tpl := txtutil.NewFastTemplate(fastTplFormat, "{{", "}}")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf := pool.Get().(*bytes.Buffer)
		buf.Reset()
		tpl.ExecuteFuncWriter(buf, func(tag string) []byte {
			return []byte(fmt.Sprint(fastTplIn[tag]))
		})
		_ = buf.Bytes()
		pool.Put(buf)
	}

}
