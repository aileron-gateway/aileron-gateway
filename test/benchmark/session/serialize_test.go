package session_test

import (
	"crypto/rand"
	"encoding/json"
	"testing"

	"github.com/vmihailenco/msgpack/v5"
)

var testMap = map[string][]byte{
	"key00": randomBytes(100),
	"key01": randomBytes(100),
	"key02": randomBytes(100),
	"key03": randomBytes(100),
	"key04": randomBytes(100),
	"key05": randomBytes(100),
	"key06": randomBytes(100),
	"key07": randomBytes(100),
	"key08": randomBytes(100),
	"key09": randomBytes(100),
}

func randomBytes(length int) []byte {
	b := make([]byte, length)
	rand.Reader.Read(b)
	return b
}

func BenchmarkMarshalJSON(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		json.Marshal(testMap)
	}
}

func BenchmarkUnmarshalJSON(b *testing.B) {
	in, _ := json.Marshal(testMap)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		into := map[string][]byte{}
		json.Unmarshal(in, &into)
	}
}

func BenchmarkMarshalMsgpack(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		msgpack.Marshal(testMap)
	}
}

func BenchmarkUnmarshalMsgpack(b *testing.B) {
	in, _ := msgpack.Marshal(testMap)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		into := map[string][]byte{}
		msgpack.Unmarshal(in, &into)
	}
}
