package errorutil_test

import (
	"io"
	"testing"

	"github.com/aileron-gateway/aileron-gateway/core"
)

func BenchmarkWithoutStack(b *testing.B) {
	err := core.ErrPrimitive.WithStack(io.EOF, nil)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		core.ErrPrimitive.WithoutStack(err, nil)
	}
}

func BenchmarkWithStack(b *testing.B) {
	err := core.ErrPrimitive.WithStack(io.EOF, nil)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		core.ErrPrimitive.WithStack(err, nil)
	}
}
