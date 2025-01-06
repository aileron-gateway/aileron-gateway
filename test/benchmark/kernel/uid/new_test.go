package uid_test

import (
	"testing"

	"github.com/aileron-gateway/aileron-gateway/kernel/uid"
)

func BenchmarkNewID(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		uid.NewID()
	}
}

func BenchmarkNewHostedID(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		uid.NewHostedID()
	}
}
