package hash_test

import (
	"testing"

	"github.com/aileron-gateway/aileron-gateway/kernel/hash"
)

var (
	testInputBytes = []byte("1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
)

func BenchmarkSHA1(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		hash.SHA1(testInputBytes)
	}
}

func BenchmarkSHA224(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		hash.SHA224(testInputBytes)
	}
}

func BenchmarkSHA256(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		hash.SHA256(testInputBytes)
	}
}

func BenchmarkSHA384(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		hash.SHA384(testInputBytes)
	}
}

func BenchmarkSHA512(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		hash.SHA512(testInputBytes)
	}
}

func BenchmarkSHA512_224(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		hash.SHA512_224(testInputBytes)
	}
}

func BenchmarkSHA512_256(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		hash.SHA512_256(testInputBytes)
	}
}

func BenchmarkSHA3_224(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		hash.SHA3_224(testInputBytes)
	}
}

func BenchmarkSHA3_256(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		hash.SHA3_256(testInputBytes)
	}
}

func BenchmarkSHA3_384(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		hash.SHA3_384(testInputBytes)
	}
}

func BenchmarkSHA3_512(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		hash.SHA3_512(testInputBytes)
	}
}

func BenchmarkSHAKE128(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		hash.SHAKE128(testInputBytes)
	}
}

func BenchmarkSHAKE256(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		hash.SHAKE256(testInputBytes)
	}
}

func BenchmarkMD5(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		hash.MD5(testInputBytes)
	}
}

func BenchmarkFNV1_32(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		hash.FNV1_32(testInputBytes)
	}
}

func BenchmarkFNV1a_32(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		hash.FNV1a_32(testInputBytes)
	}
}

func BenchmarkFNV1_64(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		hash.FNV1_64(testInputBytes)
	}
}

func BenchmarkFNV1a_64(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		hash.FNV1a_64(testInputBytes)
	}
}

func BenchmarkFNV1_128(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		hash.FNV1_128(testInputBytes)
	}
}

func BenchmarkFNV1a_128(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		hash.FNV1a_128(testInputBytes)
	}
}

func BenchmarkCRC32(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		hash.CRC32(testInputBytes)
	}
}

func BenchmarkCRC64ISO(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		hash.CRC64ISO(testInputBytes)
	}
}

func BenchmarkCRC64ECMA(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		hash.CRC64ECMA(testInputBytes)
	}
}

func BenchmarkBLAKE2s_256(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		hash.BLAKE2s_256(testInputBytes)
	}
}

func BenchmarkBLAKE2b_256(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		hash.BLAKE2b_256(testInputBytes)
	}
}

func BenchmarkBLAKE2b_384(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		hash.BLAKE2b_384(testInputBytes)
	}
}

func BenchmarkBLAKE2b_512(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		hash.BLAKE2b_512(testInputBytes)
	}
}
