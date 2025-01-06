package mac_test

import (
	"testing"

	"github.com/aileron-gateway/aileron-gateway/kernel/mac"
)

var (
	testInputMsg = []byte("1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	testInputKey = []byte("test key test key test key")
)

func BenchmarkSHA1(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mac.SHA1(testInputMsg, testInputKey)
	}
}

func BenchmarkSHA224(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mac.SHA224(testInputMsg, testInputKey)
	}
}

func BenchmarkSHA256(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mac.SHA256(testInputMsg, testInputKey)
	}
}

func BenchmarkSHA384(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mac.SHA384(testInputMsg, testInputKey)
	}
}

func BenchmarkSHA512(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mac.SHA512(testInputMsg, testInputKey)
	}
}

func BenchmarkSHA512_224(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mac.SHA512_224(testInputMsg, testInputKey)
	}
}

func BenchmarkSHA512_256(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mac.SHA512_256(testInputMsg, testInputKey)
	}
}

func BenchmarkSHA3_224(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mac.SHA3_224(testInputMsg, testInputKey)
	}
}

func BenchmarkSHA3_256(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mac.SHA3_256(testInputMsg, testInputKey)
	}
}

func BenchmarkSHA3_384(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mac.SHA3_384(testInputMsg, testInputKey)
	}
}

func BenchmarkSHA3_512(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mac.SHA3_512(testInputMsg, testInputKey)
	}
}

func BenchmarkSHAKE128(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mac.SHAKE128(testInputMsg, testInputKey)
	}
}

func BenchmarkSHAKE256(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mac.SHAKE256(testInputMsg, testInputKey)
	}
}

func BenchmarkMD5(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mac.MD5(testInputMsg, testInputKey)
	}
}

func BenchmarkFNV1_32(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mac.FNV1_32(testInputMsg, testInputKey)
	}
}

func BenchmarkFNV1a_32(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mac.FNV1a_32(testInputMsg, testInputKey)
	}
}

func BenchmarkFNV1_64(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mac.FNV1_64(testInputMsg, testInputKey)
	}
}

func BenchmarkFNV1a_64(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mac.FNV1a_64(testInputMsg, testInputKey)
	}
}

func BenchmarkFNV1_128(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mac.FNV1_128(testInputMsg, testInputKey)
	}
}

func BenchmarkFNV1a_128(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mac.FNV1a_128(testInputMsg, testInputKey)
	}
}

func BenchmarkBLAKE2s_256(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mac.BLAKE2s_256(testInputMsg, testInputKey)
	}
}

func BenchmarkBLAKE2b_256(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mac.BLAKE2b_256(testInputMsg, testInputKey)
	}
}

func BenchmarkBLAKE2b_384(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mac.BLAKE2b_384(testInputMsg, testInputKey)
	}
}

func BenchmarkBLAKE2b_512(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mac.BLAKE2b_512(testInputMsg, testInputKey)
	}
}
