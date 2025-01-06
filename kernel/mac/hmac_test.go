package mac_test

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"regexp"
	"testing"

	k "github.com/aileron-gateway/aileron-gateway/apis/kernel"
	"github.com/aileron-gateway/aileron-gateway/kernel/hash"
	"github.com/aileron-gateway/aileron-gateway/kernel/mac"
	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
)

func ExampleSHA1() {
	msg := []byte("example message")
	key := []byte("secret for hmac")

	digest := mac.SHA1(msg, key)

	encoded := base64.StdEncoding.EncodeToString(digest)
	fmt.Println(len(digest), encoded)
	// Output:
	// 20 4LOZkt8mWyQf5NxMEcEDX5WhN8c=
}

func ExampleSHA224() {
	msg := []byte("example message")
	key := []byte("secret for hmac")

	digest := mac.SHA224(msg, key)

	encoded := base64.StdEncoding.EncodeToString(digest)
	fmt.Println(len(digest), encoded)
	// Output:
	// 28 JAPTLg1JivkWrK7aKBbTycRe64jLOPc87eSGyQ==
}

func ExampleSHA256() {
	msg := []byte("example message")
	key := []byte("secret for hmac")

	digest := mac.SHA256(msg, key)

	encoded := base64.StdEncoding.EncodeToString(digest)
	fmt.Println(len(digest), encoded)
	// Output:
	// 32 WFZrK6ngcbdD/W2azEZzBMLg40nPVeqabygE3DMFtEY=
}

func ExampleSHA384() {
	msg := []byte("example message")
	key := []byte("secret for hmac")

	digest := mac.SHA384(msg, key)

	encoded := base64.StdEncoding.EncodeToString(digest)
	fmt.Println(len(digest), encoded)
	// Output:
	// 48 dyD5h1Q7D41gdp+yBAUDe1zzKwD1VlOSyyAASoMOZR6AJVHAaE8fkNqiSsv1yORO
}

func ExampleSHA512() {
	msg := []byte("example message")
	key := []byte("secret for hmac")

	digest := mac.SHA512(msg, key)

	encoded := base64.StdEncoding.EncodeToString(digest)
	fmt.Println(len(digest), encoded)
	// Output:
	// 64 PcMMi83k/XsvUrVGc2fcU2e+b2GszxMdu2v2Pk2r3TXlwS4jdCohdTPkf2eVAs4r1zW77lWtMPppF9NyLrR2hQ==
}

func ExampleSHA512_224() {
	msg := []byte("example message")
	key := []byte("secret for hmac")

	digest := mac.SHA512_224(msg, key)

	encoded := base64.StdEncoding.EncodeToString(digest)
	fmt.Println(len(digest), encoded)
	// Output:
	// 28 tTsSLZ/bxSp2Ifofm2Va/IBliDJhCNTY15z8fQ==
}

func ExampleSHA512_256() {
	msg := []byte("example message")
	key := []byte("secret for hmac")

	digest := mac.SHA512_256(msg, key)

	encoded := base64.StdEncoding.EncodeToString(digest)
	fmt.Println(len(digest), encoded)
	// Output:
	// 32 NHOVfFLktbYpX2rNRz2wgX55EMfebYUZFgBacLUquU4=
}

func ExampleSHA3_224() {
	msg := []byte("example message")
	key := []byte("secret for hmac")

	digest := mac.SHA3_224(msg, key)

	encoded := base64.StdEncoding.EncodeToString(digest)
	fmt.Println(len(digest), encoded)
	// Output:
	// 28 Li652UKyUyc2PVWVjqJWQni737KlTCA1KQt/cQ==
}

func ExampleSHA3_256() {
	msg := []byte("example message")
	key := []byte("secret for hmac")

	digest := mac.SHA3_256(msg, key)

	encoded := base64.StdEncoding.EncodeToString(digest)
	fmt.Println(len(digest), encoded)
	// Output:
	// 32 boiMtQ4CXU8noZieiIYPifiQQkRu2pX1nPWHovdPYUg=
}

func ExampleSHA3_384() {
	msg := []byte("example message")
	key := []byte("secret for hmac")

	digest := mac.SHA3_384(msg, key)

	encoded := base64.StdEncoding.EncodeToString(digest)
	fmt.Println(len(digest), encoded)
	// Output:
	// 48 xj73sf9wrFfxS6CV8G5jZgJMmIOV/eTuiZCz0MplSSbF8bysVnAHXORCd1G2Mw/g
}

func ExampleSHA3_512() {
	msg := []byte("example message")
	key := []byte("secret for hmac")

	digest := mac.SHA3_512(msg, key)

	encoded := base64.StdEncoding.EncodeToString(digest)
	fmt.Println(len(digest), encoded)
	// Output:
	// 64 hCasKS6FBZ/cEJvArXJ7NUIqCZKxHfQJUCsEKYDU0xzRnMTrkFZNb3NlSbzVt6Oq8PRJxlcSs97AXhWQrceb6A==
}

func ExampleSHAKE128() {
	msg := []byte("example message")
	key := []byte("secret for hmac")

	digest := mac.SHAKE128(msg, key)

	encoded := base64.StdEncoding.EncodeToString(digest)
	fmt.Println(len(digest), encoded)
	// Output:
	// 32 Ae8X1jG2bd43G0bYG8fd4nsC7I8kI6anYy/lquEiBNU=
}

func ExampleSHAKE256() {
	msg := []byte("example message")
	key := []byte("secret for hmac")

	digest := mac.SHAKE256(msg, key)

	encoded := base64.StdEncoding.EncodeToString(digest)
	fmt.Println(len(digest), encoded)
	// Output:
	// 64 d6z2IHJm4LC1uiHL0YznRyJDl4gaGu85dpECxV4BxOR2G6pJ9MgQWKLST62jn84Nb6fG5Tt/K4R9SnRDUFElbQ==
}

func ExampleMD5() {
	msg := []byte("example message")
	key := []byte("secret for hmac")

	digest := mac.MD5(msg, key)

	encoded := base64.StdEncoding.EncodeToString(digest)
	fmt.Println(len(digest), encoded)
	// Output:
	// 16 eErgLGC6EFa7UAQjVwuG/g==
}

func ExampleFNV1_32() {
	msg := []byte("example message")
	key := []byte("secret for hmac")

	digest := mac.FNV1_32(msg, key)

	encoded := base64.StdEncoding.EncodeToString(digest)
	fmt.Println(len(digest), encoded)
	// Output:
	// 4 SYWIwg==
}

func ExampleFNV1a_32() {
	msg := []byte("example message")
	key := []byte("secret for hmac")

	digest := mac.FNV1a_32(msg, key)

	encoded := base64.StdEncoding.EncodeToString(digest)
	fmt.Println(len(digest), encoded)
	// Output:
	// 4 dLgb2w==
}

func ExampleFNV1_64() {
	msg := []byte("example message")
	key := []byte("secret for hmac")

	digest := mac.FNV1_64(msg, key)

	encoded := base64.StdEncoding.EncodeToString(digest)
	fmt.Println(len(digest), encoded)
	// Output:
	// 8 OjabFqF/jaI=
}

func ExampleFNV1a_64() {
	msg := []byte("example message")
	key := []byte("secret for hmac")

	digest := mac.FNV1a_64(msg, key)

	encoded := base64.StdEncoding.EncodeToString(digest)
	fmt.Println(len(digest), encoded)
	// Output:
	// 8 TSoLdIf1zVg=
}

func ExampleFNV1_128() {
	msg := []byte("example message")
	key := []byte("secret for hmac")

	digest := mac.FNV1_128(msg, key)

	encoded := base64.StdEncoding.EncodeToString(digest)
	fmt.Println(len(digest), encoded)
	// Output:
	// 16 Q3XsQtHcR5+wEgOs1T6/sg==
}

func ExampleFNV1a_128() {
	msg := []byte("example message")
	key := []byte("secret for hmac")

	digest := mac.FNV1a_128(msg, key)

	encoded := base64.StdEncoding.EncodeToString(digest)
	fmt.Println(len(digest), encoded)
	// Output:
	// 16 TP9tuAX/L2bWrH+mdkueoA==
}

func ExampleCRC32() {
	msg := []byte("example message")
	key := []byte("secret for hmac")

	digest := mac.CRC32(msg, key)

	encoded := base64.StdEncoding.EncodeToString(digest)
	fmt.Println(len(digest), encoded)
	// Output:
	// 4 wl7DIg==
}

func ExampleCRC64ISO() {
	msg := []byte("example message")
	key := []byte("secret for hmac")

	digest := mac.CRC64ISO(msg, key)

	encoded := base64.StdEncoding.EncodeToString(digest)
	fmt.Println(len(digest), encoded)
	// Output:
	// 8 g1bqLbSoSUU=
}

func ExampleCRC64ECMA() {
	msg := []byte("example message")
	key := []byte("secret for hmac")

	digest := mac.CRC64ECMA(msg, key)

	encoded := base64.StdEncoding.EncodeToString(digest)
	fmt.Println(len(digest), encoded)
	// Output:
	// 8 sF7ENfkRxx8=
}

func ExampleBLAKE2s_256() {
	msg := []byte("example message")
	key := []byte("secret for hmac")

	digest := mac.BLAKE2s_256(msg, key)

	encoded := base64.StdEncoding.EncodeToString(digest)
	fmt.Println(len(digest), encoded)
	// Output:
	// 32 nj4oa8YOPZ/U/GDQ6ToPJOKp8jgxDgkfy3f6StxS+z8=
}

func ExampleBLAKE2b_256() {
	msg := []byte("example message")
	key := []byte("secret for hmac")

	digest := mac.BLAKE2b_256(msg, key)

	encoded := base64.StdEncoding.EncodeToString(digest)
	fmt.Println(len(digest), encoded)
	// Output:
	// 32 /p4fX5gCU5/AoI8H0XzThAxwxAXP+qK49a64chXqhZg=
}

func ExampleBLAKE2b_384() {
	msg := []byte("example message")
	key := []byte("secret for hmac")

	digest := mac.BLAKE2b_384(msg, key)

	encoded := base64.StdEncoding.EncodeToString(digest)
	fmt.Println(len(digest), encoded)
	// Output:
	// 48 0iSbUI+jjC0zOWwAXDzq1GGKcy8YApcMljgvQOC8dnGIPMpXXild7pWcXn3CC9sf
}

func ExampleBLAKE2b_512() {
	msg := []byte("example message")
	key := []byte("secret for hmac")

	digest := mac.BLAKE2b_512(msg, key)

	encoded := base64.StdEncoding.EncodeToString(digest)
	fmt.Println(len(digest), encoded)
	// Output:
	// 64 FrNHiVgaojnHhWT3bN3lQJk97eg5CFMdQJMXg3RknNe0VMzS0JdC9iaMfrIP+UfGpQFXVhX88ghkhxEhIjmZ7g==
}

func hexMustDecode(s string) []byte {
	b, _ := hex.DecodeString(s)
	return b
}

func TestFromAlgorithm(t *testing.T) {
	type condition struct {
		alg int
	}

	type action struct {
		f mac.HMACFunc
	}

	CndHMACExists := "HMAC exists"
	actCheckNil := "nil"

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	tb.Condition(CndHMACExists, "give an existing algorithm id")
	tb.Action(actCheckNil, "check that nil was returned")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"SHA1",
			[]string{CndHMACExists},
			[]string{},
			&condition{
				alg: int(mac.AlgSHA1),
			},
			&action{
				f: mac.SHA1,
			},
		),
		gen(
			"SHA224",
			[]string{CndHMACExists},
			[]string{},
			&condition{
				alg: int(mac.AlgSHA224),
			},
			&action{
				f: mac.SHA224,
			},
		),
		gen(
			"SHA256",
			[]string{CndHMACExists},
			[]string{},
			&condition{
				alg: int(mac.AlgSHA256),
			},
			&action{
				f: mac.SHA256,
			},
		),
		gen(
			"SHA384",
			[]string{CndHMACExists},
			[]string{},
			&condition{
				alg: int(mac.AlgSHA384),
			},
			&action{
				f: mac.SHA384,
			},
		),
		gen(
			"SHA512",
			[]string{CndHMACExists},
			[]string{},
			&condition{
				alg: int(mac.AlgSHA512),
			},
			&action{
				f: mac.SHA512,
			},
		),
		gen(
			"SHA512_224",
			[]string{CndHMACExists},
			[]string{},
			&condition{
				alg: int(mac.AlgSHA512_224),
			},
			&action{
				f: mac.SHA512_224,
			},
		),
		gen(
			"SHA512_256",
			[]string{CndHMACExists},
			[]string{},
			&condition{
				alg: int(mac.AlgSHA512_256),
			},
			&action{
				f: mac.SHA512_256,
			},
		),
		gen(
			"SHA3_224",
			[]string{CndHMACExists},
			[]string{},
			&condition{
				alg: int(mac.AlgSHA3_224),
			},
			&action{
				f: mac.SHA3_224,
			},
		),
		gen(
			"SHA3_256",
			[]string{CndHMACExists},
			[]string{},
			&condition{
				alg: int(mac.AlgSHA3_256),
			},
			&action{
				f: mac.SHA3_256,
			},
		),
		gen(
			"SHA3_384",
			[]string{CndHMACExists},
			[]string{},
			&condition{
				alg: int(mac.AlgSHA3_384),
			},
			&action{
				f: mac.SHA3_384,
			},
		),
		gen(
			"SHA3_512",
			[]string{CndHMACExists},
			[]string{},
			&condition{
				alg: int(mac.AlgSHA3_512),
			},
			&action{
				f: mac.SHA3_512,
			},
		),
		gen(
			"SHAKE128",
			[]string{CndHMACExists},
			[]string{},
			&condition{
				alg: int(mac.AlgSHAKE128),
			},
			&action{
				f: mac.SHAKE128,
			},
		),
		gen(
			"SHAKE256",
			[]string{CndHMACExists},
			[]string{},
			&condition{
				alg: int(mac.AlgSHAKE256),
			},
			&action{
				f: mac.SHAKE256,
			},
		),
		gen(
			"MD5",
			[]string{CndHMACExists},
			[]string{},
			&condition{
				alg: int(mac.AlgMD5),
			},
			&action{
				f: mac.MD5,
			},
		),
		gen(
			"FNV1_32",
			[]string{CndHMACExists},
			[]string{},
			&condition{
				alg: int(mac.AlgFNV1_32),
			},
			&action{
				f: mac.FNV1_32,
			},
		),
		gen(
			"FNV1a_32",
			[]string{CndHMACExists},
			[]string{},
			&condition{
				alg: int(mac.AlgFNV1a_32),
			},
			&action{
				f: mac.FNV1a_32,
			},
		),
		gen(
			"FNV1_64",
			[]string{CndHMACExists},
			[]string{},
			&condition{
				alg: int(mac.AlgFNV1_64),
			},
			&action{
				f: mac.FNV1_64,
			},
		),
		gen(
			"FNV1a_64",
			[]string{CndHMACExists},
			[]string{},
			&condition{
				alg: int(mac.AlgFNV1a_64),
			},
			&action{
				f: mac.FNV1a_64,
			},
		),
		gen(
			"FNV1_128",
			[]string{CndHMACExists},
			[]string{},
			&condition{
				alg: int(mac.AlgFNV1_128),
			},
			&action{
				f: mac.FNV1_128,
			},
		),
		gen(
			"FNV1a_128",
			[]string{CndHMACExists},
			[]string{},
			&condition{
				alg: int(mac.AlgFNV1a_128),
			},
			&action{
				f: mac.FNV1a_128,
			},
		),
		gen(
			"BLAKE2s_256",
			[]string{CndHMACExists},
			[]string{},
			&condition{
				alg: int(mac.AlgBLAKE2s_256),
			},
			&action{
				f: mac.BLAKE2s_256,
			},
		),
		gen(
			"BLAKE2b_256",
			[]string{CndHMACExists},
			[]string{},
			&condition{
				alg: int(mac.AlgBLAKE2b_256),
			},
			&action{
				f: mac.BLAKE2b_256,
			},
		),
		gen(
			"BLAKE2b_384",
			[]string{CndHMACExists},
			[]string{},
			&condition{
				alg: int(mac.AlgBLAKE2b_384),
			},
			&action{
				f: mac.BLAKE2b_384,
			},
		),
		gen(
			"BLAKE2b_512",
			[]string{CndHMACExists},
			[]string{},
			&condition{
				alg: int(mac.AlgBLAKE2b_512),
			},
			&action{
				f: mac.BLAKE2b_512,
			},
		),
		gen(
			"not exist",
			[]string{},
			[]string{actCheckNil},
			&condition{
				alg: 9999,
			},
			&action{
				f: nil,
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			h := mac.FromAlgorithm(mac.Algorithm(tt.C().alg))
			if tt.A().f == nil {
				testutil.Diff(t, tt.A().f, h)
				return
			}

			msg := []byte("example message")
			key := []byte("secret for hmac")
			testutil.Diff(t, tt.A().f(msg, key), h(msg, key))
		})
	}
}

func TestFromHashAlg(t *testing.T) {
	type condition struct {
		typ int32
	}

	type action struct {
		f mac.HMACFunc
	}

	CndHMACExists := "HMAC exists"
	actCheckNil := "nil"

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	tb.Condition(CndHMACExists, "give an existing algorithm type")
	tb.Action(actCheckNil, "check that nil was returned")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"SHA1",
			[]string{CndHMACExists},
			[]string{},
			&condition{
				typ: int32(k.HashAlg_SHA1),
			},
			&action{
				f: mac.SHA1,
			},
		),
		gen(
			"SHA224",
			[]string{CndHMACExists},
			[]string{},
			&condition{
				typ: int32(k.HashAlg_SHA224),
			},
			&action{
				f: mac.SHA224,
			},
		),
		gen(
			"SHA256",
			[]string{CndHMACExists},
			[]string{},
			&condition{
				typ: int32(k.HashAlg_SHA256),
			},
			&action{
				f: mac.SHA256,
			},
		),
		gen(
			"SHA384",
			[]string{CndHMACExists},
			[]string{},
			&condition{
				typ: int32(k.HashAlg_SHA384),
			},
			&action{
				f: mac.SHA384,
			},
		),
		gen(
			"SHA512",
			[]string{CndHMACExists},
			[]string{},
			&condition{
				typ: int32(k.HashAlg_SHA512),
			},
			&action{
				f: mac.SHA512,
			},
		),
		gen(
			"SHA512_224",
			[]string{CndHMACExists},
			[]string{},
			&condition{
				typ: int32(k.HashAlg_SHA512_224),
			},
			&action{
				f: mac.SHA512_224,
			},
		),
		gen(
			"SHA512_256",
			[]string{CndHMACExists},
			[]string{},
			&condition{
				typ: int32(k.HashAlg_SHA512_256),
			},
			&action{
				f: mac.SHA512_256,
			},
		),
		gen(
			"SHA3_224",
			[]string{CndHMACExists},
			[]string{},
			&condition{
				typ: int32(k.HashAlg_SHA3_224),
			},
			&action{
				f: mac.SHA3_224,
			},
		),
		gen(
			"SHA3_256",
			[]string{CndHMACExists},
			[]string{},
			&condition{
				typ: int32(k.HashAlg_SHA3_256),
			},
			&action{
				f: mac.SHA3_256,
			},
		),
		gen(
			"SHA3_384",
			[]string{CndHMACExists},
			[]string{},
			&condition{
				typ: int32(k.HashAlg_SHA3_384),
			},
			&action{
				f: mac.SHA3_384,
			},
		),
		gen(
			"SHA3_512",
			[]string{CndHMACExists},
			[]string{},
			&condition{
				typ: int32(k.HashAlg_SHA3_512),
			},
			&action{
				f: mac.SHA3_512,
			},
		),
		gen(
			"SHAKE128",
			[]string{CndHMACExists},
			[]string{},
			&condition{
				typ: int32(k.HashAlg_SHAKE128),
			},
			&action{
				f: mac.SHAKE128,
			},
		),
		gen(
			"SHAKE256",
			[]string{CndHMACExists},
			[]string{},
			&condition{
				typ: int32(k.HashAlg_SHAKE256),
			},
			&action{
				f: mac.SHAKE256,
			},
		),
		gen(
			"MD5",
			[]string{CndHMACExists},
			[]string{},
			&condition{
				typ: int32(k.HashAlg_MD5),
			},
			&action{
				f: mac.MD5,
			},
		),
		gen(
			"FNV1_32",
			[]string{CndHMACExists},
			[]string{},
			&condition{
				typ: int32(k.HashAlg_FNV1_32),
			},
			&action{
				f: mac.FNV1_32,
			},
		),
		gen(
			"FNV1a_32",
			[]string{CndHMACExists},
			[]string{},
			&condition{
				typ: int32(k.HashAlg_FNV1a_32),
			},
			&action{
				f: mac.FNV1a_32,
			},
		),
		gen(
			"FNV1_64",
			[]string{CndHMACExists},
			[]string{},
			&condition{
				typ: int32(k.HashAlg_FNV1_64),
			},
			&action{
				f: mac.FNV1_64,
			},
		),
		gen(
			"FNV1a_64",
			[]string{CndHMACExists},
			[]string{},
			&condition{
				typ: int32(k.HashAlg_FNV1a_64),
			},
			&action{
				f: mac.FNV1a_64,
			},
		),
		gen(
			"FNV1_128",
			[]string{CndHMACExists},
			[]string{},
			&condition{
				typ: int32(k.HashAlg_FNV1_128),
			},
			&action{
				f: mac.FNV1_128,
			},
		),
		gen(
			"FNV1a_128",
			[]string{CndHMACExists},
			[]string{},
			&condition{
				typ: int32(k.HashAlg_FNV1a_128),
			},
			&action{
				f: mac.FNV1a_128,
			},
		),
		gen(
			"BLAKE2s_256",
			[]string{CndHMACExists},
			[]string{},
			&condition{
				typ: int32(k.HashAlg_BLAKE2s_256),
			},
			&action{
				f: mac.BLAKE2s_256,
			},
		),
		gen(
			"BLAKE2b_256",
			[]string{CndHMACExists},
			[]string{},
			&condition{
				typ: int32(k.HashAlg_BLAKE2b_256),
			},
			&action{
				f: mac.BLAKE2b_256,
			},
		),
		gen(
			"BLAKE2b_384",
			[]string{CndHMACExists},
			[]string{},
			&condition{
				typ: int32(k.HashAlg_BLAKE2b_384),
			},
			&action{
				f: mac.BLAKE2b_384,
			},
		),
		gen(
			"BLAKE2b_512",
			[]string{CndHMACExists},
			[]string{},
			&condition{
				typ: int32(k.HashAlg_BLAKE2b_512),
			},
			&action{
				f: mac.BLAKE2b_512,
			},
		),
		gen(
			"unknown",
			[]string{},
			[]string{},
			&condition{
				typ: int32(k.HashAlg_HashAlgUnknown),
			},
			&action{
				f: nil,
			},
		),
		gen(
			"not exist",
			[]string{},
			[]string{},
			&condition{
				typ: 9999,
			},
			&action{
				f: nil,
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			h := mac.FromHashAlg(k.HashAlg(tt.C().typ))
			if tt.A().f == nil {
				testutil.Diff(t, tt.A().f, h)
				return
			}

			msg := []byte("example message")
			key := []byte("secret for hmac")
			testutil.Diff(t, tt.A().f(msg, key), h(msg, key))
		})
	}
}

func TestSHA1(t *testing.T) {
	type condition struct {
		msg []byte
		key []byte
	}

	type action struct {
		nonMatch bool
		expect   []byte
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	cndInputMessage := tb.Condition("non-nil message", "input non-zero or non-nil message")
	cndInputNilMessage := tb.Condition("nil message", "input nil message")
	cndInputKey := tb.Condition("non-nil key", "input non-nil key")
	cndInputNilKey := tb.Condition("nil key", "input nil key")
	actCheckHash := tb.Action("check hash", "check that the expected hash is returned")
	actCheckMatch := tb.Action("check match", "check that the expected hash matched")
	actCheckNotMatch := tb.Action("check not match", "check that the expected hash does not matched")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"non-nil msg & non-nil key",
			[]string{cndInputMessage, cndInputKey},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				msg: []byte("abc"),
				key: []byte("test"),
			},
			&action{
				// The expected value is generated by openssl.
				// echo -n "abc" | openssl dgst -hmac "test" -sha1
				expect: hexMustDecode("890e4312dd7b0dfa4d7c49f3850d188299fef454"),
			},
		),
		gen(
			"non-nil msg & nil key",
			[]string{cndInputMessage, cndInputNilKey},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				msg: []byte("abc"),
				key: nil,
			},
			&action{
				// The expected value is generated by openssl.
				// echo -n "abc" | openssl dgst -hmac "" -sha1
				expect: hexMustDecode("9b4a918f398d74d3e367970aba3cbe54e4d2b5d9"),
			},
		),
		gen(
			"nil msg & non-nil key",
			[]string{cndInputNilMessage, cndInputKey},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				msg: nil,
				key: []byte("test"),
			},
			&action{
				// The expected value is generated by openssl.
				// echo -n "" | openssl dgst -hmac "test" -sha1
				expect: hexMustDecode("fc85087452696e5bcbe3b7a71fde00e320af2cca"),
			},
		),
		gen(
			"nil msg & nil key",
			[]string{cndInputNilMessage, cndInputNilKey},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				msg: nil,
				key: nil,
			},
			&action{
				// The expected value is generated by openssl.
				// echo -n "" | openssl dgst -hmac "" -sha1
				expect: hexMustDecode("fbdb1d1b18aa6c08324b7d64b71fb76370690e1d"),
			},
		),
		gen(
			"wrong hash",
			[]string{cndInputMessage, cndInputKey},
			[]string{actCheckHash, actCheckNotMatch},
			&condition{
				msg: []byte("abc"),
				key: []byte("test"),
			},
			&action{
				// The first character is wrong.
				expect:   hexMustDecode("090e4312dd7b0dfa4d7c49f3850d188299fef454"),
				nonMatch: true,
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			out := mac.SHA1(tt.C().msg, tt.C().key)
			testutil.Diff(t, mac.SizeSHA1, len(out))
			testutil.Diff(t, mac.HashSize[k.HashAlg_SHA1], len(out))

			if tt.A().nonMatch {
				if bytes.Equal(tt.A().expect, out) {
					t.Error("error unexpectedly matched", tt.A().expect, out)
				}
			} else {
				testutil.Diff(t, tt.A().expect, out)
			}

			// Check that the digest generated by HMAC is different from the plain mac.
			plainHash := hash.SHA1(tt.C().msg)
			testutil.Diff(t, false, bytes.Equal(plainHash, out))
		})
	}
}

func TestSHA224(t *testing.T) {
	type condition struct {
		msg []byte
		key []byte
	}

	type action struct {
		nonMatch bool
		expect   []byte
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	cndInputMessage := tb.Condition("non-nil message", "input non-zero or non-nil message")
	cndInputNilMessage := tb.Condition("nil message", "input nil message")
	cndInputKey := tb.Condition("non-nil key", "input non-nil key")
	cndInputNilKey := tb.Condition("nil key", "input nil key")
	actCheckHash := tb.Action("check hash", "check that the expected hash is returned")
	actCheckMatch := tb.Action("check match", "check that the expected hash matched")
	actCheckNotMatch := tb.Action("check not match", "check that the expected hash does not matched")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"non-nil msg & non-nil key",
			[]string{cndInputMessage, cndInputKey},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				msg: []byte("abc"),
				key: []byte("test"),
			},
			&action{
				// The expected value is generated by openssl.
				// echo -n "abc" | openssl dgst -hmac "test" -sha224
				expect: hexMustDecode("cfa85a0c1fa92b8f1cdc520b92098ed7ffa5e255c89c64e79efd01a1"),
			},
		),
		gen(
			"non-nil msg & nil key",
			[]string{cndInputMessage, cndInputNilKey},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				msg: []byte("abc"),
				key: nil,
			},
			&action{
				// The expected value is generated by openssl.
				// echo -n "abc" | openssl dgst -hmac "" -sha224
				expect: hexMustDecode("d473c456fa6aad72bbec9c6ad63ca92d8675caa0b7f451fa4b692081"),
			},
		),
		gen(
			"nil msg & non-nil key",
			[]string{cndInputNilMessage, cndInputKey},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				msg: nil,
				key: []byte("test"),
			},
			&action{
				// The expected value is generated by openssl.
				// echo -n "" | openssl dgst -hmac "test" -sha224
				expect: hexMustDecode("2d015ac5c9d652e088a81131172d982ba765ce38d2486899dceedb7c"),
			},
		),
		gen(
			"nil msg & nil key",
			[]string{cndInputNilMessage, cndInputNilKey},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				msg: nil,
				key: nil,
			},
			&action{
				// The expected value is generated by openssl.
				// echo -n "" | openssl dgst -hmac "" -sha224
				expect: hexMustDecode("5ce14f72894662213e2748d2a6ba234b74263910cedde2f5a9271524"),
			},
		),
		gen(
			"wrong hash",
			[]string{cndInputMessage, cndInputKey},
			[]string{actCheckHash, actCheckNotMatch},
			&condition{
				msg: []byte("abc"),
				key: []byte("test"),
			},
			&action{
				// The first character is wrong.
				expect:   hexMustDecode("0ce14f72894662213e2748d2a6ba234b74263910cedde2f5a9271524"),
				nonMatch: true,
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			out := mac.SHA224(tt.C().msg, tt.C().key)
			testutil.Diff(t, mac.SizeSHA224, len(out))
			testutil.Diff(t, mac.HashSize[k.HashAlg_SHA224], len(out))

			if tt.A().nonMatch {
				if bytes.Equal(tt.A().expect, out) {
					t.Error("error unexpectedly matched", tt.A().expect, out)
				}
			} else {
				testutil.Diff(t, tt.A().expect, out)
			}

			// Check that the digest generated by HMAC is different from the plain mac.
			plainHash := hash.SHA224(tt.C().msg)
			testutil.Diff(t, false, bytes.Equal(plainHash, out))
		})
	}
}

func TestSHA256(t *testing.T) {
	type condition struct {
		msg []byte
		key []byte
	}

	type action struct {
		nonMatch bool
		expect   []byte

		err        any // error or errorutil.Kind
		errPattern *regexp.Regexp
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	cndInputMessage := tb.Condition("non-nil message", "input non-zero or non-nil message")
	cndInputNilMessage := tb.Condition("nil message", "input nil message")
	cndInputKey := tb.Condition("non-nil key", "input non-nil key")
	cndInputNilKey := tb.Condition("nil key", "input nil key")
	actCheckHash := tb.Action("check hash", "check that the expected hash is returned")
	actCheckMatch := tb.Action("check match", "check that the expected hash matched")
	actCheckNotMatch := tb.Action("check not match", "check that the expected hash does not matched")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"non-nil msg & non-nil key",
			[]string{cndInputMessage, cndInputKey},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				msg: []byte("abc"),
				key: []byte("test"),
			},
			&action{
				// The expected value is generated by openssl.
				// echo -n "abc" | openssl dgst -hmac "test" -sha256
				expect: hexMustDecode("d796579aed123e7b743ccaf5b150affa1223e31ecba8b88c9da9ccf7ad5e0594"),
			},
		),
		gen(
			"non-nil msg & nil key",
			[]string{cndInputMessage, cndInputNilKey},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				msg: []byte("abc"),
				key: nil,
			},
			&action{
				// The expected value is generated by openssl.
				// echo -n "abc" | openssl dgst -hmac "" -sha256
				expect: hexMustDecode("fd7adb152c05ef80dccf50a1fa4c05d5a3ec6da95575fc312ae7c5d091836351"),
			},
		),
		gen(
			"nil msg & non-nil key",
			[]string{cndInputNilMessage, cndInputKey},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				msg: nil,
				key: []byte("test"),
			},
			&action{
				// The expected value is generated by openssl.
				// echo -n "" | openssl dgst -hmac "test" -sha256
				expect: hexMustDecode("ad71148c79f21ab9eec51ea5c7dd2b668792f7c0d3534ae66b22f71c61523fb3"),
			},
		),
		gen(
			"nil msg & nil key",
			[]string{cndInputNilMessage, cndInputNilKey},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				msg: nil,
				key: nil,
			},
			&action{
				// The expected value is generated by openssl.
				// echo -n "" | openssl dgst -hmac "" -sha256
				expect: hexMustDecode("b613679a0814d9ec772f95d778c35fc5ff1697c493715653c6c712144292c5ad"),
			},
		),
		gen(
			"wrong hash",
			[]string{cndInputMessage, cndInputKey},
			[]string{actCheckHash, actCheckNotMatch},
			&condition{
				msg: []byte("abc"),
				key: []byte("test"),
			},
			&action{
				// The first character is wrong.
				expect:   hexMustDecode("0796579aed123e7b743ccaf5b150affa1223e31ecba8b88c9da9ccf7ad5e0594"),
				nonMatch: true,
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			out := mac.SHA256(tt.C().msg, tt.C().key)
			testutil.Diff(t, mac.SizeSHA256, len(out))
			testutil.Diff(t, mac.HashSize[k.HashAlg_SHA256], len(out))

			if tt.A().nonMatch {
				if bytes.Equal(tt.A().expect, out) {
					t.Error("error unexpectedly matched", tt.A().expect, out)
				}
			} else {
				testutil.Diff(t, tt.A().expect, out)
			}

			// Check that the digest generated by HMAC is different from the plain mac.
			plainHash := hash.SHA256(tt.C().msg)
			testutil.Diff(t, false, bytes.Equal(plainHash, out))
		})
	}
}

func TestSHA384(t *testing.T) {
	type condition struct {
		msg []byte
		key []byte
	}

	type action struct {
		nonMatch bool
		expect   []byte
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	cndInputMessage := tb.Condition("non-nil message", "input non-zero or non-nil message")
	cndInputNilMessage := tb.Condition("nil message", "input nil message")
	cndInputKey := tb.Condition("non-nil key", "input non-nil key")
	cndInputNilKey := tb.Condition("nil key", "input nil key")
	actCheckHash := tb.Action("check hash", "check that the expected hash is returned")
	actCheckMatch := tb.Action("check match", "check that the expected hash matched")
	actCheckNotMatch := tb.Action("check not match", "check that the expected hash does not matched")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"non-nil msg & non-nil key",
			[]string{cndInputMessage, cndInputKey},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				msg: []byte("abc"),
				key: []byte("test"),
			},
			&action{
				// The expected value is generated by openssl.
				// echo -n "abc" | openssl dgst -hmac "test" -sha384
				expect: hexMustDecode("edc15e84defef9af6eee3b0f878d4b1d043a0dd7d75682f6f5f927848762e715f1e7db68715ceb5af3b32e05f7807a9d"),
			},
		),
		gen(
			"non-nil msg & nil key",
			[]string{cndInputMessage, cndInputNilKey},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				msg: []byte("abc"),
				key: nil,
			},
			&action{
				// The expected value is generated by openssl.
				// echo -n "abc" | openssl dgst -hmac "" -sha384
				expect: hexMustDecode("948f7c5caa500c31d7d4a0f52f3e3da7e33c8a9fe6ef528b8a9ac3e4adc4e24d908e6f40b737510e82354759dc5e9f06"),
			},
		),
		gen(
			"nil msg & non-nil key",
			[]string{cndInputNilMessage, cndInputKey},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				msg: nil,
				key: []byte("test"),
			},
			&action{
				// The expected value is generated by openssl.
				// echo -n "" | openssl dgst -hmac "test" -sha384
				expect: hexMustDecode("bda08a334994873233c844d24f0e7cf8c76c6e9feeb9c25ce97b9446e8efe3e06c261741ca21580360f20f1fd2190e0a"),
			},
		),
		gen(
			"nil msg & nil key",
			[]string{cndInputNilMessage, cndInputNilKey},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				msg: nil,
				key: nil,
			},
			&action{
				// The expected value is generated by openssl.
				// echo -n "" | openssl dgst -hmac "" -sha384
				expect: hexMustDecode("6c1f2ee938fad2e24bd91298474382ca218c75db3d83e114b3d4367776d14d3551289e75e8209cd4b792302840234adc"),
			},
		),
		gen(
			"wrong hash",
			[]string{cndInputMessage, cndInputKey},
			[]string{actCheckHash, actCheckNotMatch},
			&condition{
				msg: []byte("abc"),
				key: []byte("test"),
			},
			&action{
				// The first character is wrong.
				expect:   hexMustDecode("0dc15e84defef9af6eee3b0f878d4b1d043a0dd7d75682f6f5f927848762e715f1e7db68715ceb5af3b32e05f7807a9d"),
				nonMatch: true,
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			out := mac.SHA384(tt.C().msg, tt.C().key)
			testutil.Diff(t, mac.SizeSHA384, len(out))
			testutil.Diff(t, mac.HashSize[k.HashAlg_SHA384], len(out))

			if tt.A().nonMatch {
				if bytes.Equal(tt.A().expect, out) {
					t.Error("error unexpectedly matched", tt.A().expect, out)
				}
			} else {
				testutil.Diff(t, tt.A().expect, out)
			}

			// Check that the digest generated by HMAC is different from the plain mac.
			plainHash := hash.SHA384(tt.C().msg)
			testutil.Diff(t, false, bytes.Equal(plainHash, out))
		})
	}
}

func TestSHA512(t *testing.T) {
	type condition struct {
		msg []byte
		key []byte
	}

	type action struct {
		nonMatch bool
		expect   []byte

		err        any // error or errorutil.Kind
		errPattern *regexp.Regexp
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	cndInputMessage := tb.Condition("non-nil message", "input non-zero or non-nil message")
	cndInputNilMessage := tb.Condition("nil message", "input nil message")
	cndInputKey := tb.Condition("non-nil key", "input non-nil key")
	cndInputNilKey := tb.Condition("nil key", "input nil key")
	actCheckHash := tb.Action("check hash", "check that the expected hash is returned")
	actCheckMatch := tb.Action("check match", "check that the expected hash matched")
	actCheckNotMatch := tb.Action("check not match", "check that the expected hash does not matched")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"non-nil msg & non-nil key",
			[]string{cndInputMessage, cndInputKey},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				msg: []byte("abc"),
				key: []byte("test"),
			},
			&action{
				// The expected value is generated by openssl.
				// echo -n "abc" | openssl dgst -hmac "test" -sha512
				expect: hexMustDecode("9e60c79cefe413ddb61a0435ffc62546dc82b9172a78f4a3f49c76d5ad218fcbe699b801ee3693139d1ef132009a831a1c97fcd062c673087c5e0f74786f89c7"),
			},
		),
		gen(
			"non-nil msg & nil key",
			[]string{cndInputMessage, cndInputNilKey},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				msg: []byte("abc"),
				key: nil,
			},
			&action{
				// The expected value is generated by openssl.
				// echo -n "abc" | openssl dgst -hmac "" -sha512
				expect: hexMustDecode("29689f6b79a8dd686068c2eeae97fd8769ad3ba65cb5381f838358a8045a358ee3ba1739c689c7805e31734fb6072f87261d1256995370d55725cba00d10bdd0"),
			},
		),
		gen(
			"nil msg & non-nil key",
			[]string{cndInputNilMessage, cndInputKey},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				msg: nil,
				key: []byte("test"),
			},
			&action{
				// The expected value is generated by openssl.
				// echo -n "" | openssl dgst -hmac "test" -sha512
				expect: hexMustDecode("01917bf85be0c998598a2332f75c2fe6f662c0900d4391123ca2bc61f073ede360af8f3afd6e5d3f28dff4b57cc22890aa7b7498cf441f32a6f6e78aca3cafe8"),
			},
		),
		gen(
			"nil msg & nil key",
			[]string{cndInputNilMessage, cndInputNilKey},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				msg: nil,
				key: nil,
			},
			&action{
				// The expected value is generated by openssl.
				// echo -n "" | openssl dgst -hmac "" -sha512
				expect: hexMustDecode("b936cee86c9f87aa5d3c6f2e84cb5a4239a5fe50480a6ec66b70ab5b1f4ac6730c6c515421b327ec1d69402e53dfb49ad7381eb067b338fd7b0cb22247225d47"),
			},
		),
		gen(
			"wrong hash",
			[]string{cndInputMessage, cndInputKey},
			[]string{actCheckHash, actCheckNotMatch},
			&condition{
				msg: []byte("abc"),
				key: []byte("test"),
			},
			&action{
				// The first character is wrong.
				expect:   hexMustDecode("0e60c79cefe413ddb61a0435ffc62546dc82b9172a78f4a3f49c76d5ad218fcbe699b801ee3693139d1ef132009a831a1c97fcd062c673087c5e0f74786f89c7"),
				nonMatch: true,
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			out := mac.SHA512(tt.C().msg, tt.C().key)
			testutil.Diff(t, mac.SizeSHA512, len(out))
			testutil.Diff(t, mac.HashSize[k.HashAlg_SHA512], len(out))

			if tt.A().nonMatch {
				if bytes.Equal(tt.A().expect, out) {
					t.Error("error unexpectedly matched", tt.A().expect, out)
				}
			} else {
				testutil.Diff(t, tt.A().expect, out)
			}

			// Check that the digest generated by HMAC is different from the plain mac.
			plainHash := hash.SHA512(tt.C().msg)
			testutil.Diff(t, false, bytes.Equal(plainHash, out))
		})
	}
}

func TestSHA512_224(t *testing.T) {
	type condition struct {
		msg []byte
		key []byte
	}

	type action struct {
		nonMatch bool
		expect   []byte

		err        any // error or errorutil.Kind
		errPattern *regexp.Regexp
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	cndInputMessage := tb.Condition("non-nil message", "input non-zero or non-nil message")
	cndInputNilMessage := tb.Condition("nil message", "input nil message")
	cndInputKey := tb.Condition("non-nil key", "input non-nil key")
	cndInputNilKey := tb.Condition("nil key", "input nil key")
	actCheckHash := tb.Action("check hash", "check that the expected hash is returned")
	actCheckMatch := tb.Action("check match", "check that the expected hash matched")
	actCheckNotMatch := tb.Action("check not match", "check that the expected hash does not matched")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"non-nil msg & non-nil key",
			[]string{cndInputMessage, cndInputKey},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				msg: []byte("abc"),
				key: []byte("test"),
			},
			&action{
				// The expected value is generated by openssl.
				// echo -n "abc" | openssl dgst -hmac "test" -sha512-224
				expect: hexMustDecode("2de0ca71ec29a2d3443d861123eee62eb658d1d598bf6ff4463850d8"),
			},
		),
		gen(
			"non-nil msg & nil key",
			[]string{cndInputMessage, cndInputNilKey},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				msg: []byte("abc"),
				key: nil,
			},
			&action{
				// The expected value is generated by openssl.
				// echo -n "abc" | openssl dgst -hmac "" -sha512-224
				expect: hexMustDecode("bbff170efa4e26333b19388fe35029e8d7c49c6443f5b2acb6e110dc"),
			},
		),
		gen(
			"nil msg & non-nil key",
			[]string{cndInputNilMessage, cndInputKey},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				msg: nil,
				key: []byte("test"),
			},
			&action{
				// The expected value is generated by openssl.
				// echo -n "" | openssl dgst -hmac "test" -sha512-224
				expect: hexMustDecode("5c7dc567bc41657076f3c966dfa319648f47102eb76ad28efa5034d8"),
			},
		),
		gen(
			"nil msg & nil key",
			[]string{cndInputNilMessage, cndInputNilKey},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				msg: nil,
				key: nil,
			},
			&action{
				// The expected value is generated by openssl.
				// echo -n "" | openssl dgst -hmac "" -sha512-224
				expect: hexMustDecode("de43f6b96f2d08cebe1ee9c02c53d96b68c1e55b6c15d6843b410d4c"),
			},
		),
		gen(
			"wrong hash",
			[]string{cndInputMessage, cndInputKey},
			[]string{actCheckHash, actCheckNotMatch},
			&condition{
				msg: []byte("abc"),
				key: []byte("test"),
			},
			&action{
				// The first character is wrong.
				expect:   hexMustDecode("0de0ca71ec29a2d3443d861123eee62eb658d1d598bf6ff4463850d8"),
				nonMatch: true,
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			out := mac.SHA512_224(tt.C().msg, tt.C().key)
			testutil.Diff(t, mac.SizeSHA512_224, len(out))
			testutil.Diff(t, mac.HashSize[k.HashAlg_SHA512_224], len(out))

			if tt.A().nonMatch {
				if bytes.Equal(tt.A().expect, out) {
					t.Error("error unexpectedly matched", tt.A().expect, out)
				}
			} else {
				testutil.Diff(t, tt.A().expect, out)
			}

			// Check that the digest generated by HMAC is different from the plain mac.
			plainHash := hash.SHA512_224(tt.C().msg)
			testutil.Diff(t, false, bytes.Equal(plainHash, out))
		})
	}
}

func TestSHA512_256(t *testing.T) {
	type condition struct {
		msg []byte
		key []byte
	}

	type action struct {
		nonMatch bool
		expect   []byte

		err        any // error or errorutil.Kind
		errPattern *regexp.Regexp
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	cndInputMessage := tb.Condition("non-nil message", "input non-zero or non-nil message")
	cndInputNilMessage := tb.Condition("nil message", "input nil message")
	cndInputKey := tb.Condition("non-nil key", "input non-nil key")
	cndInputNilKey := tb.Condition("nil key", "input nil key")
	actCheckHash := tb.Action("check hash", "check that the expected hash is returned")
	actCheckMatch := tb.Action("check match", "check that the expected hash matched")
	actCheckNotMatch := tb.Action("check not match", "check that the expected hash does not matched")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"non-nil msg & non-nil key",
			[]string{cndInputMessage, cndInputKey},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				msg: []byte("abc"),
				key: []byte("test"),
			},
			&action{
				// The expected value is generated by openssl.
				// echo -n "abc" | openssl dgst -hmac "test" -sha512-256
				expect: hexMustDecode("f0b2935b1664431097be514ebb2c462d98a6892b37ea73c7c959a72e61bf1645"),
			},
		),
		gen(
			"non-nil msg & nil key",
			[]string{cndInputMessage, cndInputNilKey},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				msg: []byte("abc"),
				key: nil,
			},
			&action{
				// The expected value is generated by openssl.
				// echo -n "abc" | openssl dgst -hmac "" -sha512-256
				expect: hexMustDecode("784cac6aafefd5517029bae0cd223d58111dc37f390d982fae2a0548b5aa67ea"),
			},
		),
		gen(
			"nil msg & non-nil key",
			[]string{cndInputNilMessage, cndInputKey},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				msg: nil,
				key: []byte("test"),
			},
			&action{
				// The expected value is generated by openssl.
				// echo -n "" | openssl dgst -hmac "test" -sha512-256
				expect: hexMustDecode("030d89c3c7fb441f2d7da37c033c53a827bb81c0236957d4566643117a24ec1c"),
			},
		),
		gen(
			"nil msg & nil key",
			[]string{cndInputNilMessage, cndInputNilKey},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				msg: nil,
				key: nil,
			},
			&action{
				// The expected value is generated by openssl.
				// echo -n "" | openssl dgst -hmac "" -sha512-256
				expect: hexMustDecode("b79c9951df595274582dc094a1ba46c33e4a36878b2d83cb8553f0fe467dcdcf"),
			},
		),
		gen(
			"wrong hash",
			[]string{cndInputMessage, cndInputKey},
			[]string{actCheckHash, actCheckNotMatch},
			&condition{
				msg: []byte("abc"),
				key: []byte("test"),
			},
			&action{
				// The first character is wrong.
				expect:   hexMustDecode("00b2935b1664431097be514ebb2c462d98a6892b37ea73c7c959a72e61bf1645"),
				nonMatch: true,
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			out := mac.SHA512_256(tt.C().msg, tt.C().key)
			testutil.Diff(t, mac.SizeSHA512_256, len(out))
			testutil.Diff(t, mac.HashSize[k.HashAlg_SHA512_256], len(out))

			if tt.A().nonMatch {
				if bytes.Equal(tt.A().expect, out) {
					t.Error("error unexpectedly matched", tt.A().expect, out)
				}
			} else {
				testutil.Diff(t, tt.A().expect, out)
			}

			// Check that the digest generated by HMAC is different from the plain mac.
			plainHash := hash.SHA512_256(tt.C().msg)
			testutil.Diff(t, false, bytes.Equal(plainHash, out))
		})
	}
}

func TestSHA3_224(t *testing.T) {
	type condition struct {
		msg []byte
		key []byte
	}

	type action struct {
		nonMatch bool
		expect   []byte
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	cndInputMessage := tb.Condition("non-nil message", "input non-zero or non-nil message")
	cndInputNilMessage := tb.Condition("nil message", "input nil message")
	cndInputKey := tb.Condition("non-nil key", "input non-nil key")
	cndInputNilKey := tb.Condition("nil key", "input nil key")
	actCheckHash := tb.Action("check hash", "check that the expected hash is returned")
	actCheckMatch := tb.Action("check match", "check that the expected hash matched")
	actCheckNotMatch := tb.Action("check not match", "check that the expected hash does not matched")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"non-nil msg & non-nil key",
			[]string{cndInputMessage, cndInputKey},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				msg: []byte("abc"),
				key: []byte("test"),
			},
			&action{
				// The expected value is generated by openssl.
				// echo -n "abc" | openssl dgst -hmac "test" -sha3-224
				expect: hexMustDecode("44e52c93252fb5ebbb33df3e81ff432d3cf70f6a33e19379f806ae92"),
			},
		),
		gen(
			"non-nil msg & nil key",
			[]string{cndInputMessage, cndInputNilKey},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				msg: []byte("abc"),
				key: nil,
			},
			&action{
				// The expected value is generated by openssl.
				// echo -n "abc" | openssl dgst -hmac "" -sha3-224
				expect: hexMustDecode("45590e15e0bcb331f9de5c3bf8ef8e7d018a45397735e5a8a1ce0b9b"),
			},
		),
		gen(
			"nil msg & non-nil key",
			[]string{cndInputNilMessage, cndInputKey},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				msg: nil,
				key: []byte("test"),
			},
			&action{
				// The expected value is generated by openssl.
				// echo -n "" | openssl dgst -hmac "test" -sha3-224
				expect: hexMustDecode("d30278220969497275016b0287903d08d0274e1afe57cc9729204b31"),
			},
		),
		gen(
			"nil msg & nil key",
			[]string{cndInputNilMessage, cndInputNilKey},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				msg: nil,
				key: nil,
			},
			&action{
				// The expected value is generated by openssl.
				// echo -n "" | openssl dgst -hmac "" -sha3-224
				expect: hexMustDecode("1b9044e0d5bb4ef944bc00f1b26c483ac3e222f4640935d089a49083"),
			},
		),
		gen(
			"wrong hash",
			[]string{cndInputMessage, cndInputKey},
			[]string{actCheckHash, actCheckNotMatch},
			&condition{
				msg: []byte("abc"),
				key: []byte("test"),
			},
			&action{
				// The first character is wrong.
				expect:   hexMustDecode("04e52c93252fb5ebbb33df3e81ff432d3cf70f6a33e19379f806ae92"),
				nonMatch: true,
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			out := mac.SHA3_224(tt.C().msg, tt.C().key)
			testutil.Diff(t, mac.SizeSHA3_224, len(out))
			testutil.Diff(t, mac.HashSize[k.HashAlg_SHA3_224], len(out))

			if tt.A().nonMatch {
				if bytes.Equal(tt.A().expect, out) {
					t.Error("error unexpectedly matched", tt.A().expect, out)
				}
			} else {
				testutil.Diff(t, tt.A().expect, out)
			}

			// Check that the digest generated by HMAC is different from the plain mac.
			plainHash := hash.SHA3_224(tt.C().msg)
			testutil.Diff(t, false, bytes.Equal(plainHash, out))
		})
	}
}

func TestSHA3_256(t *testing.T) {
	type condition struct {
		msg []byte
		key []byte
	}

	type action struct {
		nonMatch bool
		expect   []byte
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	cndInputMessage := tb.Condition("non-nil message", "input non-zero or non-nil message")
	cndInputNilMessage := tb.Condition("nil message", "input nil message")
	cndInputKey := tb.Condition("non-nil key", "input non-nil key")
	cndInputNilKey := tb.Condition("nil key", "input nil key")
	actCheckHash := tb.Action("check hash", "check that the expected hash is returned")
	actCheckMatch := tb.Action("check match", "check that the expected hash matched")
	actCheckNotMatch := tb.Action("check not match", "check that the expected hash does not matched")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"non-nil msg & non-nil key",
			[]string{cndInputMessage, cndInputKey},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				msg: []byte("abc"),
				key: []byte("test"),
			},
			&action{
				// The expected value is generated by openssl.
				// echo -n "abc" | openssl dgst -hmac "test" -sha3-256
				expect: hexMustDecode("99a994571f3133ba8358086c137d01cef2b4372decee2d0362e9ad53f5348b0c"),
			},
		),
		gen(
			"non-nil msg & nil key",
			[]string{cndInputMessage, cndInputNilKey},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				msg: []byte("abc"),
				key: nil,
			},
			&action{
				// The expected value is generated by openssl.
				// echo -n "abc" | openssl dgst -hmac "" -sha3-256
				expect: hexMustDecode("776bdf4f598121a2ac38c408d375731a5681f10998e77dcb92fb474adfef8f90"),
			},
		),
		gen(
			"nil msg & non-nil key",
			[]string{cndInputNilMessage, cndInputKey},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				msg: nil,
				key: []byte("test"),
			},
			&action{
				// The expected value is generated by openssl.
				// echo -n "" | openssl dgst -hmac "test" -sha3-256
				expect: hexMustDecode("d1177a2cb9cb5ba5bc74891e3f12764656a16c0f872f317255c59737cae921d4"),
			},
		),
		gen(
			"nil msg & nil key",
			[]string{cndInputNilMessage, cndInputNilKey},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				msg: nil,
				key: nil,
			},
			&action{
				// The expected value is generated by openssl.
				// echo -n "" | openssl dgst -hmac "" -sha3-256
				expect: hexMustDecode("e841c164e5b4f10c9f3985587962af72fd607a951196fc92fb3a5251941784ea"),
			},
		),
		gen(
			"wrong hash",
			[]string{cndInputMessage, cndInputKey},
			[]string{actCheckHash, actCheckNotMatch},
			&condition{
				msg: []byte("abc"),
				key: []byte("test"),
			},
			&action{
				// The first character is wrong.
				expect:   hexMustDecode("09a994571f3133ba8358086c137d01cef2b4372decee2d0362e9ad53f5348b0c"),
				nonMatch: true,
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			out := mac.SHA3_256(tt.C().msg, tt.C().key)
			testutil.Diff(t, mac.SizeSHA3_256, len(out))
			testutil.Diff(t, mac.HashSize[k.HashAlg_SHA3_256], len(out))

			if tt.A().nonMatch {
				if bytes.Equal(tt.A().expect, out) {
					t.Error("error unexpectedly matched", tt.A().expect, out)
				}
			} else {
				testutil.Diff(t, tt.A().expect, out)
			}

			// Check that the digest generated by HMAC is different from the plain mac.
			plainHash := hash.SHA3_256(tt.C().msg)
			testutil.Diff(t, false, bytes.Equal(plainHash, out))
		})
	}
}

func TestSHA3_384(t *testing.T) {
	type condition struct {
		msg []byte
		key []byte
	}

	type action struct {
		nonMatch bool
		expect   []byte
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	cndInputMessage := tb.Condition("non-nil message", "input non-zero or non-nil message")
	cndInputNilMessage := tb.Condition("nil message", "input nil message")
	cndInputKey := tb.Condition("non-nil key", "input non-nil key")
	cndInputNilKey := tb.Condition("nil key", "input nil key")
	actCheckHash := tb.Action("check hash", "check that the expected hash is returned")
	actCheckMatch := tb.Action("check match", "check that the expected hash matched")
	actCheckNotMatch := tb.Action("check not match", "check that the expected hash does not matched")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"non-nil msg & non-nil key",
			[]string{cndInputMessage, cndInputKey},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				msg: []byte("abc"),
				key: []byte("test"),
			},
			&action{
				// The expected value is generated by openssl.
				// echo -n "abc" | openssl dgst -hmac "test" -sha3-384
				expect: hexMustDecode("e6dfd4e9a8dfb2c387f0f072cab7920f10ffbc999490fc603a92b99cf976040f3c24cca78c98a4d65ea0c29f02c44877"),
			},
		),
		gen(
			"non-nil msg & nil key",
			[]string{cndInputMessage, cndInputNilKey},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				msg: []byte("abc"),
				key: nil,
			},
			&action{
				// The expected value is generated by openssl.
				// echo -n "abc" | openssl dgst -hmac "" -sha3-384
				expect: hexMustDecode("175064a7663c17b5fa7abcd6fb9d0969b7be431dc4d310d2196d06b23cd8a9db5b13013038437b636de8dfa38edc3452"),
			},
		),
		gen(
			"nil msg & non-nil key",
			[]string{cndInputNilMessage, cndInputKey},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				msg: nil,
				key: []byte("test"),
			},
			&action{
				// The expected value is generated by openssl.
				// echo -n "" | openssl dgst -hmac "test" -sha3-384
				expect: hexMustDecode("52cc09670eb96a16a894559378b122c877bb7aa3eb7a79c4672f3a5efbe2c0f874526a1bcabfee094febc82cc385bb49"),
			},
		),
		gen(
			"nil msg & nil key",
			[]string{cndInputNilMessage, cndInputNilKey},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				msg: nil,
				key: nil,
			},
			&action{
				// The expected value is generated by openssl.
				// echo -n "" | openssl dgst -hmac "" -sha3-384
				expect: hexMustDecode("adca89f07bbfbeaf58880c1572379ea2416568fd3b66542bd42599c57c4567e6ae086299ea216c6f3e7aef90b6191d24"),
			},
		),
		gen(
			"wrong hash",
			[]string{cndInputMessage, cndInputKey},
			[]string{actCheckHash, actCheckNotMatch},
			&condition{
				msg: []byte("abc"),
				key: []byte("test"),
			},
			&action{
				// The first character is wrong.
				expect:   hexMustDecode("06dfd4e9a8dfb2c387f0f072cab7920f10ffbc999490fc603a92b99cf976040f3c24cca78c98a4d65ea0c29f02c44877"),
				nonMatch: true,
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			out := mac.SHA3_384(tt.C().msg, tt.C().key)
			testutil.Diff(t, mac.SizeSHA3_384, len(out))
			testutil.Diff(t, mac.HashSize[k.HashAlg_SHA3_384], len(out))

			if tt.A().nonMatch {
				if bytes.Equal(tt.A().expect, out) {
					t.Error("error unexpectedly matched", tt.A().expect, out)
				}
			} else {
				testutil.Diff(t, tt.A().expect, out)
			}

			// Check that the digest generated by HMAC is different from the plain mac.
			plainHash := hash.SHA3_384(tt.C().msg)
			testutil.Diff(t, false, bytes.Equal(plainHash, out))
		})
	}
}

func TestSHA3_512(t *testing.T) {
	type condition struct {
		msg []byte
		key []byte
	}

	type action struct {
		nonMatch bool
		expect   []byte
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	cndInputMessage := tb.Condition("non-nil message", "input non-zero or non-nil message")
	cndInputNilMessage := tb.Condition("nil message", "input nil message")
	cndInputKey := tb.Condition("non-nil key", "input non-nil key")
	cndInputNilKey := tb.Condition("nil key", "input nil key")
	actCheckHash := tb.Action("check hash", "check that the expected hash is returned")
	actCheckMatch := tb.Action("check match", "check that the expected hash matched")
	actCheckNotMatch := tb.Action("check not match", "check that the expected hash does not matched")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"non-nil msg & non-nil key",
			[]string{cndInputMessage, cndInputKey},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				msg: []byte("abc"),
				key: []byte("test"),
			},
			&action{
				// The expected value is generated by openssl.
				// echo -n "abc" | openssl dgst -hmac "test" -sha3-512
				expect: hexMustDecode("3bd7e107d43bb12a825b0aa11ed14d89cf90607dd0d0c52868f72411d67a0f11dfad0afc830ab41d43b2bccaf929d7fc8608f86ec8c199c65e9f57c31bdd99cf"),
			},
		),
		gen(
			"non-nil msg & nil key",
			[]string{cndInputMessage, cndInputNilKey},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				msg: []byte("abc"),
				key: nil,
			},
			&action{
				// The expected value is generated by openssl.
				// echo -n "abc" | openssl dgst -hmac "" -sha3-512
				expect: hexMustDecode("aedf400b73941c7a3bfb6f72c9e6f4530216fe82d382b04da2dd5eb099f8ca4f9319c0158e5ca9368cb9fb497b0863e4b5fc62701c75eef6a8b30fdf7836bd59"),
			},
		),
		gen(
			"nil msg & non-nil key",
			[]string{cndInputNilMessage, cndInputKey},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				msg: nil,
				key: []byte("test"),
			},
			&action{
				// The expected value is generated by openssl.
				// echo -n "" | openssl dgst -hmac "test" -sha3-512
				expect: hexMustDecode("e8b82f5f0f82c8af5c79c81c9b7b5a702c687f826b420cfecc43bf2ea41e7b763f7c21a55da89e21e4234b819eb01e844a7b9e9e329e31b1cf457d5b415ca688"),
			},
		),
		gen(
			"nil msg & nil key",
			[]string{cndInputNilMessage, cndInputNilKey},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				msg: nil,
				key: nil,
			},
			&action{
				// The expected value is generated by openssl.
				// echo -n "" | openssl dgst -hmac "" -sha3-512
				expect: hexMustDecode("cbcf45540782d4bc7387fbbf7d30b3681d6d66cc435cafd82546b0fce96b367ea79662918436fba442e81a01d0f9592dfcd30f7a7a8f1475693d30be4150ca84"),
			},
		),
		gen(
			"wrong hash",
			[]string{cndInputMessage, cndInputKey},
			[]string{actCheckHash, actCheckNotMatch},
			&condition{
				msg: []byte("abc"),
				key: []byte("test"),
			},
			&action{
				// The first character is wrong.
				expect:   hexMustDecode("0bd7e107d43bb12a825b0aa11ed14d89cf90607dd0d0c52868f72411d67a0f11dfad0afc830ab41d43b2bccaf929d7fc8608f86ec8c199c65e9f57c31bdd99cf"),
				nonMatch: true,
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			out := mac.SHA3_512(tt.C().msg, tt.C().key)
			testutil.Diff(t, mac.SizeSHA3_512, len(out))
			testutil.Diff(t, mac.HashSize[k.HashAlg_SHA3_512], len(out))

			if tt.A().nonMatch {
				if bytes.Equal(tt.A().expect, out) {
					t.Error("error unexpectedly matched", tt.A().expect, out)
				}
			} else {
				testutil.Diff(t, tt.A().expect, out)
			}

			// Check that the digest generated by HMAC is different from the plain mac.
			plainHash := hash.SHA3_512(tt.C().msg)
			testutil.Diff(t, false, bytes.Equal(plainHash, out))
		})
	}
}

func TestSHAKE128(t *testing.T) {
	type condition struct {
		msg []byte
		key []byte
	}

	type action struct {
		nonMatch bool
		expect   []byte
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	cndInputMessage := tb.Condition("non-nil message", "input non-zero or non-nil message")
	cndInputNilMessage := tb.Condition("nil message", "input nil message")
	cndInputKey := tb.Condition("non-nil key", "input non-nil key")
	cndInputNilKey := tb.Condition("nil key", "input nil key")
	actCheckHash := tb.Action("check hash", "check that the expected hash is returned")
	actCheckMatch := tb.Action("check match", "check that the expected hash matched")
	actCheckNotMatch := tb.Action("check not match", "check that the expected hash does not matched")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"non-nil msg & non-nil key",
			[]string{cndInputMessage, cndInputKey},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				msg: []byte("abc"),
				key: []byte("test"),
			},
			&action{
				expect: hexMustDecode("517a4d18ca1edd6f1af1b5374f6023095b8a439373e8ec03c0ab520f658566f0"),
			},
		),
		gen(
			"non-nil msg & nil key",
			[]string{cndInputMessage, cndInputNilKey},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				msg: []byte("abc"),
				key: nil,
			},
			&action{
				expect: hexMustDecode("71e22c30fc7b1d92dfc6c97d604831179ecb327a213aa5ba776a24523f3c7386"),
			},
		),
		gen(
			"nil msg & non-nil key",
			[]string{cndInputNilMessage, cndInputKey},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				msg: nil,
				key: []byte("test"),
			},
			&action{
				expect: hexMustDecode("44031fd2078bbb2a9d654b4aaa5fa9a8d7e7a222681312752098e8decc4a3f16"),
			},
		),
		gen(
			"nil msg & nil key",
			[]string{cndInputNilMessage, cndInputNilKey},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				msg: nil,
				key: nil,
			},
			&action{
				expect: hexMustDecode("cb16e7aad79e0f44efc343d51c42ac752eeead4ec64a56f73737285f4847c293"),
			},
		),
		gen(
			"wrong hash",
			[]string{cndInputMessage, cndInputKey},
			[]string{actCheckHash, actCheckNotMatch},
			&condition{
				msg: []byte("abc"),
				key: []byte("test"),
			},
			&action{
				// The first character is wrong.
				expect:   hexMustDecode("017a4d18ca1edd6f1af1b5374f6023095b8a439373e8ec03c0ab520f658566f0"),
				nonMatch: true,
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			out := mac.SHAKE128(tt.C().msg, tt.C().key)
			testutil.Diff(t, mac.SizeSHAKE128, len(out))
			testutil.Diff(t, mac.HashSize[k.HashAlg_SHAKE128], len(out))

			if tt.A().nonMatch {
				if bytes.Equal(tt.A().expect, out) {
					t.Error("error unexpectedly matched", tt.A().expect, out)
				}
			} else {
				testutil.Diff(t, tt.A().expect, out)
			}

			// Check that the digest generated by HMAC is different from the plain mac.
			plainHash := hash.SHAKE128(tt.C().msg)
			testutil.Diff(t, false, bytes.Equal(plainHash, out))
		})
	}
}

func TestSHAKE256(t *testing.T) {
	type condition struct {
		msg []byte
		key []byte
	}

	type action struct {
		nonMatch bool
		expect   []byte
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	cndInputMessage := tb.Condition("non-nil message", "input non-zero or non-nil message")
	cndInputNilMessage := tb.Condition("nil message", "input nil message")
	cndInputKey := tb.Condition("non-nil key", "input non-nil key")
	cndInputNilKey := tb.Condition("nil key", "input nil key")
	actCheckHash := tb.Action("check hash", "check that the expected hash is returned")
	actCheckMatch := tb.Action("check match", "check that the expected hash matched")
	actCheckNotMatch := tb.Action("check not match", "check that the expected hash does not matched")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"non-nil msg & non-nil key",
			[]string{cndInputMessage, cndInputKey},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				msg: []byte("abc"),
				key: []byte("test"),
			},
			&action{
				expect: hexMustDecode("bb43958e647ee6fefdddd3a519f023067e369d8f492e31fcc444ecf34d13fa079656bbd61da8f3de7b4820f4698eaa6b06998dfad84aff1ac7a36924926cdf29"),
			},
		),
		gen(
			"non-nil msg & nil key",
			[]string{cndInputMessage, cndInputNilKey},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				msg: []byte("abc"),
				key: nil,
			},
			&action{
				expect: hexMustDecode("ad3586181bb4f105f209eb3c109bf9982286cccc2ebe224271c3c8eb9a27632c2d9144f1fa5a2a8907dc2f096d822ef743f5902e671acd0ccf14213342122202"),
			},
		),
		gen(
			"nil msg & non-nil key",
			[]string{cndInputNilMessage, cndInputKey},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				msg: nil,
				key: []byte("test"),
			},
			&action{
				expect: hexMustDecode("20f83210311cb06fe34ca04d02cbd71e31c86322ce135f1961f7acf3e5da8e1866036120bae0e73096a82682f630762b96b48cfb91901adc3e1807566a7687e9"),
			},
		),
		gen(
			"nil msg & nil key",
			[]string{cndInputNilMessage, cndInputNilKey},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				msg: nil,
				key: nil,
			},
			&action{
				expect: hexMustDecode("990afebd7c22034f814822c92918eee01ea81856634afc061f7cda24097ec9ae86b4179749c5ad115404c29589b0c74d33c00ae576930bc08e7d532ea711ef29"),
			},
		),
		gen(
			"wrong hash",
			[]string{cndInputMessage, cndInputKey},
			[]string{actCheckHash, actCheckNotMatch},
			&condition{
				msg: []byte("abc"),
				key: []byte("test"),
			},
			&action{
				// The first character is wrong.
				expect:   hexMustDecode("0b43958e647ee6fefdddd3a519f023067e369d8f492e31fcc444ecf34d13fa079656bbd61da8f3de7b4820f4698eaa6b06998dfad84aff1ac7a36924926cdf29"),
				nonMatch: true,
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			out := mac.SHAKE256(tt.C().msg, tt.C().key)
			testutil.Diff(t, mac.SizeSHAKE256, len(out))
			testutil.Diff(t, mac.HashSize[k.HashAlg_SHAKE256], len(out))

			if tt.A().nonMatch {
				if bytes.Equal(tt.A().expect, out) {
					t.Error("error unexpectedly matched", tt.A().expect, out)
				}
			} else {
				testutil.Diff(t, tt.A().expect, out)
			}

			// Check that the digest generated by HMAC is different from the plain mac.
			plainHash := hash.SHAKE256(tt.C().msg)
			testutil.Diff(t, false, bytes.Equal(plainHash, out))
		})
	}
}

func TestMD5(t *testing.T) {
	type condition struct {
		msg []byte
		key []byte
	}

	type action struct {
		nonMatch bool
		expect   []byte
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	cndInputMessage := tb.Condition("non-nil message", "input non-zero or non-nil message")
	cndInputNilMessage := tb.Condition("nil message", "input nil message")
	cndInputKey := tb.Condition("non-nil key", "input non-nil key")
	cndInputNilKey := tb.Condition("nil key", "input nil key")
	actCheckHash := tb.Action("check hash", "check that the expected hash is returned")
	actCheckMatch := tb.Action("check match", "check that the expected hash matched")
	actCheckNotMatch := tb.Action("check not match", "check that the expected hash does not matched")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"non-nil msg & non-nil key",
			[]string{cndInputMessage, cndInputKey},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				msg: []byte("abc"),
				key: []byte("test"),
			},
			&action{
				// The expected value is generated by openssl.
				// echo -n "abc" | openssl dgst -hmac "test" -md5
				expect: hexMustDecode("ea273a03b6d975926e973ffca956bded"),
			},
		),
		gen(
			"non-nil msg & nil key",
			[]string{cndInputMessage, cndInputNilKey},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				msg: []byte("abc"),
				key: nil,
			},
			&action{
				// The expected value is generated by openssl.
				// echo -n "abc" | openssl dgst -hmac "" -md5
				expect: hexMustDecode("dd2701993d29fdd0b032c233cec63403"),
			},
		),
		gen(
			"nil msg & non-nil key",
			[]string{cndInputNilMessage, cndInputKey},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				msg: nil,
				key: []byte("test"),
			},
			&action{
				// The expected value is generated by openssl.
				// echo -n "" | openssl dgst -hmac "test" -md5
				expect: hexMustDecode("6aebb41d4366c2794687e1d9e7aee307"),
			},
		),
		gen(
			"nil msg & nil key",
			[]string{cndInputNilMessage, cndInputNilKey},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				msg: nil,
				key: nil,
			},
			&action{
				// The expected value is generated by openssl.
				// echo -n "" | openssl dgst -hmac "" -md5
				expect: hexMustDecode("74e6f7298a9c2d168935f58c001bad88"),
			},
		),
		gen(
			"wrong hash",
			[]string{cndInputMessage, cndInputKey},
			[]string{actCheckHash, actCheckNotMatch},
			&condition{
				msg: []byte("abc"),
				key: []byte("test"),
			},
			&action{
				// The first character is wrong.
				expect:   hexMustDecode("0a273a03b6d975926e973ffca956bded"),
				nonMatch: true,
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			out := mac.MD5(tt.C().msg, tt.C().key)
			testutil.Diff(t, mac.SizeMD5, len(out))
			testutil.Diff(t, mac.HashSize[k.HashAlg_MD5], len(out))

			if tt.A().nonMatch {
				if bytes.Equal(tt.A().expect, out) {
					t.Error("error unexpectedly matched", tt.A().expect, out)
				}
			} else {
				testutil.Diff(t, tt.A().expect, out)
			}

			// Check that the digest generated by HMAC is different from the plain mac.
			plainHash := hash.MD5(tt.C().msg)
			testutil.Diff(t, false, bytes.Equal(plainHash, out))
		})
	}
}

func TestFNV1_32(t *testing.T) {
	type condition struct {
		msg []byte
		key []byte
	}

	type action struct {
		nonMatch bool
		expect   []byte
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	cndInputMessage := tb.Condition("non-nil message", "input non-zero or non-nil message")
	cndInputNilMessage := tb.Condition("nil message", "input nil message")
	cndInputKey := tb.Condition("non-nil key", "input non-nil key")
	cndInputNilKey := tb.Condition("nil key", "input nil key")
	actCheckHash := tb.Action("check hash", "check that the expected hash is returned")
	actCheckMatch := tb.Action("check match", "check that the expected hash matched")
	actCheckNotMatch := tb.Action("check not match", "check that the expected hash does not matched")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"non-nil msg & non-nil key",
			[]string{cndInputMessage, cndInputKey},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				msg: []byte("abc"),
				key: []byte("test"),
			},
			&action{
				expect: hexMustDecode("e68f0cc5"),
			},
		),
		gen(
			"non-nil msg & nil key",
			[]string{cndInputMessage, cndInputNilKey},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				msg: []byte("abc"),
				key: nil,
			},
			&action{
				expect: hexMustDecode("788ca764"),
			},
		),
		gen(
			"nil msg & non-nil key",
			[]string{cndInputNilMessage, cndInputKey},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				msg: nil,
				key: []byte("test"),
			},
			&action{
				expect: hexMustDecode("0b3c0506"),
			},
		),
		gen(
			"nil msg & nil key",
			[]string{cndInputNilMessage, cndInputNilKey},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				msg: nil,
				key: nil,
			},
			&action{
				expect: hexMustDecode("fa082b06"),
			},
		),
		gen(
			"wrong hash",
			[]string{cndInputMessage, cndInputKey},
			[]string{actCheckHash, actCheckNotMatch},
			&condition{
				msg: []byte("abc"),
				key: []byte("test"),
			},
			&action{
				// The first character is wrong.
				expect:   hexMustDecode("068f0cc5"),
				nonMatch: true,
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			out := mac.FNV1_32(tt.C().msg, tt.C().key)
			testutil.Diff(t, mac.SizeFNV1_32, len(out))
			testutil.Diff(t, mac.HashSize[k.HashAlg_FNV1_32], len(out))

			if tt.A().nonMatch {
				if bytes.Equal(tt.A().expect, out) {
					t.Error("error unexpectedly matched", tt.A().expect, out)
				}
			} else {
				testutil.Diff(t, tt.A().expect, out)
			}

			// Check that the digest generated by HMAC is different from the plain mac.
			plainHash := hash.FNV1_32(tt.C().msg)
			testutil.Diff(t, false, bytes.Equal(plainHash, out))
		})
	}
}

func TestFNV1a_32(t *testing.T) {
	type condition struct {
		msg []byte
		key []byte
	}

	type action struct {
		nonMatch bool
		expect   []byte
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	cndInputMessage := tb.Condition("non-nil message", "input non-zero or non-nil message")
	cndInputNilMessage := tb.Condition("nil message", "input nil message")
	cndInputKey := tb.Condition("non-nil key", "input non-nil key")
	cndInputNilKey := tb.Condition("nil key", "input nil key")
	actCheckHash := tb.Action("check hash", "check that the expected hash is returned")
	actCheckMatch := tb.Action("check match", "check that the expected hash matched")
	actCheckNotMatch := tb.Action("check not match", "check that the expected hash does not matched")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"non-nil msg & non-nil key",
			[]string{cndInputMessage, cndInputKey},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				msg: []byte("abc"),
				key: []byte("test"),
			},
			&action{
				expect: hexMustDecode("a6735281"),
			},
		),
		gen(
			"non-nil msg & nil key",
			[]string{cndInputMessage, cndInputNilKey},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				msg: []byte("abc"),
				key: nil,
			},
			&action{
				expect: hexMustDecode("ff7ea2f3"),
			},
		),
		gen(
			"nil msg & non-nil key",
			[]string{cndInputNilMessage, cndInputKey},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				msg: nil,
				key: []byte("test"),
			},
			&action{
				expect: hexMustDecode("3017bf78"),
			},
		),
		gen(
			"nil msg & nil key",
			[]string{cndInputNilMessage, cndInputNilKey},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				msg: nil,
				key: nil,
			},
			&action{
				expect: hexMustDecode("c80045f6"),
			},
		),
		gen(
			"wrong hash",
			[]string{cndInputMessage, cndInputKey},
			[]string{actCheckHash, actCheckNotMatch},
			&condition{
				msg: []byte("abc"),
				key: []byte("test"),
			},
			&action{
				// The first character is wrong.
				expect:   hexMustDecode("06735281"),
				nonMatch: true,
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			out := mac.FNV1a_32(tt.C().msg, tt.C().key)
			testutil.Diff(t, mac.SizeFNV1a_32, len(out))
			testutil.Diff(t, mac.HashSize[k.HashAlg_FNV1a_32], len(out))

			if tt.A().nonMatch {
				if bytes.Equal(tt.A().expect, out) {
					t.Error("error unexpectedly matched", tt.A().expect, out)
				}
			} else {
				testutil.Diff(t, tt.A().expect, out)
			}

			// Check that the digest generated by HMAC is different from the plain mac.
			plainHash := hash.FNV1a_32(tt.C().msg)
			testutil.Diff(t, false, bytes.Equal(plainHash, out))
		})
	}
}

func TestFNV1_64(t *testing.T) {
	type condition struct {
		msg []byte
		key []byte
	}

	type action struct {
		nonMatch bool
		expect   []byte
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	cndInputMessage := tb.Condition("non-nil message", "input non-zero or non-nil message")
	cndInputNilMessage := tb.Condition("nil message", "input nil message")
	cndInputKey := tb.Condition("non-nil key", "input non-nil key")
	cndInputNilKey := tb.Condition("nil key", "input nil key")
	actCheckHash := tb.Action("check hash", "check that the expected hash is returned")
	actCheckMatch := tb.Action("check match", "check that the expected hash matched")
	actCheckNotMatch := tb.Action("check not match", "check that the expected hash does not matched")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"non-nil msg & non-nil key",
			[]string{cndInputMessage, cndInputKey},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				msg: []byte("abc"),
				key: []byte("test"),
			},
			&action{
				expect: hexMustDecode("8419f182f5814b95"),
			},
		),
		gen(
			"non-nil msg & nil key",
			[]string{cndInputMessage, cndInputNilKey},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				msg: []byte("abc"),
				key: nil,
			},
			&action{
				expect: hexMustDecode("2667725d9f1aea2e"),
			},
		),
		gen(
			"nil msg & non-nil key",
			[]string{cndInputNilMessage, cndInputKey},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				msg: nil,
				key: []byte("test"),
			},
			&action{
				expect: hexMustDecode("d5c261dc5d6f56c1"),
			},
		),
		gen(
			"nil msg & nil key",
			[]string{cndInputNilMessage, cndInputNilKey},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				msg: nil,
				key: nil,
			},
			&action{
				expect: hexMustDecode("e970411f2790dcd9"),
			},
		),
		gen(
			"wrong hash",
			[]string{cndInputMessage, cndInputKey},
			[]string{actCheckHash, actCheckNotMatch},
			&condition{
				msg: []byte("abc"),
				key: []byte("test"),
			},
			&action{
				// The first character is wrong.
				expect:   hexMustDecode("0419f182f5814b95"),
				nonMatch: true,
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			out := mac.FNV1_64(tt.C().msg, tt.C().key)
			testutil.Diff(t, mac.SizeFNV1_64, len(out))
			testutil.Diff(t, mac.HashSize[k.HashAlg_FNV1_64], len(out))

			if tt.A().nonMatch {
				if bytes.Equal(tt.A().expect, out) {
					t.Error("error unexpectedly matched", tt.A().expect, out)
				}
			} else {
				testutil.Diff(t, tt.A().expect, out)
			}

			// Check that the digest generated by HMAC is different from the plain mac.
			plainHash := hash.FNV1_64(tt.C().msg)
			testutil.Diff(t, false, bytes.Equal(plainHash, out))
		})
	}
}

func TestFNV1a_64(t *testing.T) {
	type condition struct {
		msg []byte
		key []byte
	}

	type action struct {
		nonMatch bool
		expect   []byte
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	cndInputMessage := tb.Condition("non-nil message", "input non-zero or non-nil message")
	cndInputNilMessage := tb.Condition("nil message", "input nil message")
	cndInputKey := tb.Condition("non-nil key", "input non-nil key")
	cndInputNilKey := tb.Condition("nil key", "input nil key")
	actCheckHash := tb.Action("check hash", "check that the expected hash is returned")
	actCheckMatch := tb.Action("check match", "check that the expected hash matched")
	actCheckNotMatch := tb.Action("check not match", "check that the expected hash does not matched")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"non-nil msg & non-nil key",
			[]string{cndInputMessage, cndInputKey},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				msg: []byte("abc"),
				key: []byte("test"),
			},
			&action{
				expect: hexMustDecode("5bddc2a6b3d16d26"),
			},
		),
		gen(
			"non-nil msg & nil key",
			[]string{cndInputMessage, cndInputNilKey},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				msg: []byte("abc"),
				key: nil,
			},
			&action{
				expect: hexMustDecode("fdffe0d798355841"),
			},
		),
		gen(
			"nil msg & non-nil key",
			[]string{cndInputNilMessage, cndInputKey},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				msg: nil,
				key: []byte("test"),
			},
			&action{
				expect: hexMustDecode("162b58f8b94fe1e2"),
			},
		),
		gen(
			"nil msg & nil key",
			[]string{cndInputNilMessage, cndInputNilKey},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				msg: nil,
				key: nil,
			},
			&action{
				expect: hexMustDecode("5a8ba6933d31b551"),
			},
		),
		gen(
			"wrong hash",
			[]string{cndInputMessage, cndInputKey},
			[]string{actCheckHash, actCheckNotMatch},
			&condition{
				msg: []byte("abc"),
				key: []byte("test"),
			},
			&action{
				// The first character is wrong.
				expect:   hexMustDecode("0bddc2a6b3d16d26"),
				nonMatch: true,
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			out := mac.FNV1a_64(tt.C().msg, tt.C().key)
			testutil.Diff(t, mac.SizeFNV1a_64, len(out))
			testutil.Diff(t, mac.HashSize[k.HashAlg_FNV1a_64], len(out))

			if tt.A().nonMatch {
				if bytes.Equal(tt.A().expect, out) {
					t.Error("error unexpectedly matched", tt.A().expect, out)
				}
			} else {
				testutil.Diff(t, tt.A().expect, out)
			}

			// Check that the digest generated by HMAC is different from the plain mac.
			plainHash := hash.FNV1a_64(tt.C().msg)
			testutil.Diff(t, false, bytes.Equal(plainHash, out))
		})
	}
}

func TestFNV1_128(t *testing.T) {
	type condition struct {
		msg []byte
		key []byte
	}

	type action struct {
		nonMatch bool
		expect   []byte
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	cndInputMessage := tb.Condition("non-nil message", "input non-zero or non-nil message")
	cndInputNilMessage := tb.Condition("nil message", "input nil message")
	cndInputKey := tb.Condition("non-nil key", "input non-nil key")
	cndInputNilKey := tb.Condition("nil key", "input nil key")
	actCheckHash := tb.Action("check hash", "check that the expected hash is returned")
	actCheckMatch := tb.Action("check match", "check that the expected hash matched")
	actCheckNotMatch := tb.Action("check not match", "check that the expected hash does not matched")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"non-nil msg & non-nil key",
			[]string{cndInputMessage, cndInputKey},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				msg: []byte("abc"),
				key: []byte("test"),
			},
			&action{
				expect: hexMustDecode("1e2de37279cd025684fdf3208dff9bc5"),
			},
		),
		gen(
			"non-nil msg & nil key",
			[]string{cndInputMessage, cndInputNilKey},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				msg: []byte("abc"),
				key: nil,
			},
			&action{
				expect: hexMustDecode("c04fa8dafde97d01a209a0993e294d33"),
			},
		),
		gen(
			"nil msg & non-nil key",
			[]string{cndInputNilMessage, cndInputKey},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				msg: nil,
				key: []byte("test"),
			},
			&action{
				expect: hexMustDecode("4903aaf686b8438a7ae6ab503cb792c5"),
			},
		),
		gen(
			"nil msg & nil key",
			[]string{cndInputNilMessage, cndInputNilKey},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				msg: nil,
				key: nil,
			},
			&action{
				expect: hexMustDecode("a56420334df413c579b19d048d773b1d"),
			},
		),
		gen(
			"wrong hash",
			[]string{cndInputMessage, cndInputKey},
			[]string{actCheckHash, actCheckNotMatch},
			&condition{
				msg: []byte("abc"),
				key: []byte("test"),
			},
			&action{
				// The first character is wrong.
				expect:   hexMustDecode("0e2de37279cd025684fdf3208dff9bc5"),
				nonMatch: true,
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			out := mac.FNV1_128(tt.C().msg, tt.C().key)
			testutil.Diff(t, mac.SizeFNV1_128, len(out))
			testutil.Diff(t, mac.HashSize[k.HashAlg_FNV1_128], len(out))

			if tt.A().nonMatch {
				if bytes.Equal(tt.A().expect, out) {
					t.Error("error unexpectedly matched", tt.A().expect, out)
				}
			} else {
				testutil.Diff(t, tt.A().expect, out)
			}

			// Check that the digest generated by HMAC is different from the plain mac.
			plainHash := hash.FNV1_128(tt.C().msg)
			testutil.Diff(t, false, bytes.Equal(plainHash, out))
		})
	}
}

func TestFNV1a_128(t *testing.T) {
	type condition struct {
		msg []byte
		key []byte
	}

	type action struct {
		nonMatch bool
		expect   []byte
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	cndInputMessage := tb.Condition("non-nil message", "input non-zero or non-nil message")
	cndInputNilMessage := tb.Condition("nil message", "input nil message")
	cndInputKey := tb.Condition("non-nil key", "input non-nil key")
	cndInputNilKey := tb.Condition("nil key", "input nil key")
	actCheckHash := tb.Action("check hash", "check that the expected hash is returned")
	actCheckMatch := tb.Action("check match", "check that the expected hash matched")
	actCheckNotMatch := tb.Action("check not match", "check that the expected hash does not matched")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"non-nil msg & non-nil key",
			[]string{cndInputMessage, cndInputKey},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				msg: []byte("abc"),
				key: []byte("test"),
			},
			&action{
				expect: hexMustDecode("b7df8c86319e86f078f9c37ac8411e04"),
			},
		),
		gen(
			"non-nil msg & nil key",
			[]string{cndInputMessage, cndInputNilKey},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				msg: []byte("abc"),
				key: nil,
			},
			&action{
				expect: hexMustDecode("5d1b169d259720fb2f2f6f7b4fe6cbc5"),
			},
		),
		gen(
			"nil msg & non-nil key",
			[]string{cndInputNilMessage, cndInputKey},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				msg: nil,
				key: []byte("test"),
			},
			&action{
				expect: hexMustDecode("1ed197b5dd9cf43eee0f118b6c9b142f"),
			},
		),
		gen(
			"nil msg & nil key",
			[]string{cndInputNilMessage, cndInputNilKey},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				msg: nil,
				key: nil,
			},
			&action{
				expect: hexMustDecode("6e23bed80cbbd704ae919d0481cb7faa"),
			},
		),
		gen(
			"wrong hash",
			[]string{cndInputMessage, cndInputKey},
			[]string{actCheckHash, actCheckNotMatch},
			&condition{
				msg: []byte("abc"),
				key: []byte("test"),
			},
			&action{
				// The first character is wrong.
				expect:   hexMustDecode("07df8c86319e86f078f9c37ac8411e04"),
				nonMatch: true,
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			out := mac.FNV1a_128(tt.C().msg, tt.C().key)
			testutil.Diff(t, mac.SizeFNV1a_128, len(out))
			testutil.Diff(t, mac.HashSize[k.HashAlg_FNV1a_128], len(out))

			if tt.A().nonMatch {
				if bytes.Equal(tt.A().expect, out) {
					t.Error("error unexpectedly matched", tt.A().expect, out)
				}
			} else {
				testutil.Diff(t, tt.A().expect, out)
			}

			// Check that the digest generated by HMAC is different from the plain mac.
			plainHash := hash.FNV1a_128(tt.C().msg)
			testutil.Diff(t, false, bytes.Equal(plainHash, out))
		})
	}
}

func TestCRC32(t *testing.T) {
	type condition struct {
		msg []byte
		key []byte
	}

	type action struct {
		nonMatch bool
		expect   []byte
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	cndInputMessage := tb.Condition("non-nil message", "input non-zero or non-nil message")
	cndInputNilMessage := tb.Condition("nil message", "input nil message")
	cndInputKey := tb.Condition("non-nil key", "input non-nil key")
	cndInputNilKey := tb.Condition("nil key", "input nil key")
	actCheckHash := tb.Action("check hash", "check that the expected hash is returned")
	actCheckMatch := tb.Action("check match", "check that the expected hash matched")
	actCheckNotMatch := tb.Action("check not match", "check that the expected hash does not matched")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"non-nil msg & non-nil key",
			[]string{cndInputMessage, cndInputKey},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				msg: []byte("abc"),
				key: []byte("test"),
			},
			&action{
				expect: hexMustDecode("85a9ee18"),
			},
		),
		gen(
			"non-nil msg & nil key",
			[]string{cndInputMessage, cndInputNilKey},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				msg: []byte("abc"),
				key: nil,
			},
			&action{
				expect: hexMustDecode("80388ba4"),
			},
		),
		gen(
			"nil msg & non-nil key",
			[]string{cndInputNilMessage, cndInputKey},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				msg: nil,
				key: []byte("test"),
			},
			&action{
				expect: hexMustDecode("0bb3b95b"),
			},
		),
		gen(
			"nil msg & nil key",
			[]string{cndInputNilMessage, cndInputNilKey},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				msg: nil,
				key: nil,
			},
			&action{
				expect: hexMustDecode("9df5b03c"),
			},
		),
		gen(
			"wrong hash",
			[]string{cndInputMessage, cndInputKey},
			[]string{actCheckHash, actCheckNotMatch},
			&condition{
				msg: []byte("abc"),
				key: []byte("test"),
			},
			&action{
				// The first character is wrong.
				expect:   hexMustDecode("05a9ee18"),
				nonMatch: true,
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			out := mac.CRC32(tt.C().msg, tt.C().key)
			testutil.Diff(t, mac.SizeCRC32, len(out))
			testutil.Diff(t, mac.HashSize[k.HashAlg_CRC32], len(out))

			if tt.A().nonMatch {
				if bytes.Equal(tt.A().expect, out) {
					t.Error("error unexpectedly matched", tt.A().expect, out)
				}
			} else {
				testutil.Diff(t, tt.A().expect, out)
			}

			// Check that the digest generated by HMAC is different from the plain mac.
			plainHash := hash.CRC32(tt.C().msg)
			testutil.Diff(t, false, bytes.Equal(plainHash, out))
		})
	}
}

func TestCRC64ISO(t *testing.T) {
	type condition struct {
		msg []byte
		key []byte
	}

	type action struct {
		nonMatch bool
		expect   []byte
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	cndInputMessage := tb.Condition("non-nil message", "input non-zero or non-nil message")
	cndInputNilMessage := tb.Condition("nil message", "input nil message")
	cndInputKey := tb.Condition("non-nil key", "input non-nil key")
	cndInputNilKey := tb.Condition("nil key", "input nil key")
	actCheckHash := tb.Action("check hash", "check that the expected hash is returned")
	actCheckMatch := tb.Action("check match", "check that the expected hash matched")
	actCheckNotMatch := tb.Action("check not match", "check that the expected hash does not matched")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"non-nil msg & non-nil key",
			[]string{cndInputMessage, cndInputKey},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				msg: []byte("abc"),
				key: []byte("test"),
			},
			&action{
				expect: hexMustDecode("8c08ff0a830d7b76"),
			},
		),
		gen(
			"non-nil msg & nil key",
			[]string{cndInputMessage, cndInputNilKey},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				msg: []byte("abc"),
				key: nil,
			},
			&action{
				expect: hexMustDecode("ae80ffd2a15d7b76"),
			},
		),
		gen(
			"nil msg & non-nil key",
			[]string{cndInputNilMessage, cndInputKey},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				msg: nil,
				key: []byte("test"),
			},
			&action{
				expect: hexMustDecode("fb08ffffffff0ab5"),
			},
		),
		gen(
			"nil msg & nil key",
			[]string{cndInputNilMessage, cndInputNilKey},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				msg: nil,
				key: nil,
			},
			&action{
				expect: hexMustDecode("ae80ffffffffd297"),
			},
		),
		gen(
			"wrong hash",
			[]string{cndInputMessage, cndInputKey},
			[]string{actCheckHash, actCheckNotMatch},
			&condition{
				msg: []byte("abc"),
				key: []byte("test"),
			},
			&action{
				// The first character is wrong.
				expect:   hexMustDecode("0c08ff0a830d7b76"),
				nonMatch: true,
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			out := mac.CRC64ISO(tt.C().msg, tt.C().key)
			testutil.Diff(t, mac.SizeCRC64ISO, len(out))
			testutil.Diff(t, mac.HashSize[k.HashAlg_CRC64ISO], len(out))

			if tt.A().nonMatch {
				if bytes.Equal(tt.A().expect, out) {
					t.Error("error unexpectedly matched", tt.A().expect, out)
				}
			} else {
				testutil.Diff(t, tt.A().expect, out)
			}

			// Check that the digest generated by HMAC is different from the plain mac.
			plainHash := hash.CRC64ISO(tt.C().msg)
			testutil.Diff(t, false, bytes.Equal(plainHash, out))
		})
	}
}

func TestCRC64ECMA(t *testing.T) {
	type condition struct {
		msg []byte
		key []byte
	}

	type action struct {
		nonMatch bool
		expect   []byte
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	cndInputMessage := tb.Condition("non-nil message", "input non-zero or non-nil message")
	cndInputNilMessage := tb.Condition("nil message", "input nil message")
	cndInputKey := tb.Condition("non-nil key", "input non-nil key")
	cndInputNilKey := tb.Condition("nil key", "input nil key")
	actCheckHash := tb.Action("check hash", "check that the expected hash is returned")
	actCheckMatch := tb.Action("check match", "check that the expected hash matched")
	actCheckNotMatch := tb.Action("check not match", "check that the expected hash does not matched")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"non-nil msg & non-nil key",
			[]string{cndInputMessage, cndInputKey},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				msg: []byte("abc"),
				key: []byte("test"),
			},
			&action{
				expect: hexMustDecode("ee1a69bfd2740cf5"),
			},
		),
		gen(
			"non-nil msg & nil key",
			[]string{cndInputMessage, cndInputNilKey},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				msg: []byte("abc"),
				key: nil,
			},
			&action{
				expect: hexMustDecode("f0d1bd36d3be4912"),
			},
		),
		gen(
			"nil msg & non-nil key",
			[]string{cndInputNilMessage, cndInputKey},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				msg: nil,
				key: []byte("test"),
			},
			&action{
				expect: hexMustDecode("20e991b2933807a1"),
			},
		),
		gen(
			"nil msg & nil key",
			[]string{cndInputNilMessage, cndInputNilKey},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				msg: nil,
				key: nil,
			},
			&action{
				expect: hexMustDecode("899f897d1b9ed353"),
			},
		),
		gen(
			"wrong hash",
			[]string{cndInputMessage, cndInputKey},
			[]string{actCheckHash, actCheckNotMatch},
			&condition{
				msg: []byte("abc"),
				key: []byte("test"),
			},
			&action{
				// The first character is wrong.
				expect:   hexMustDecode("0e1a69bfd2740cf5"),
				nonMatch: true,
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			out := mac.CRC64ECMA(tt.C().msg, tt.C().key)
			testutil.Diff(t, mac.SizeCRC64ECMA, len(out))
			testutil.Diff(t, mac.HashSize[k.HashAlg_CRC64ECMA], len(out))

			if tt.A().nonMatch {
				if bytes.Equal(tt.A().expect, out) {
					t.Error("error unexpectedly matched", tt.A().expect, out)
				}
			} else {
				testutil.Diff(t, tt.A().expect, out)
			}

			// Check that the digest generated by HMAC is different from the plain mac.
			plainHash := hash.CRC64ECMA(tt.C().msg)
			testutil.Diff(t, false, bytes.Equal(plainHash, out))
		})
	}
}

func TestBLAKE2s_256(t *testing.T) {
	type condition struct {
		msg []byte
		key []byte
	}

	type action struct {
		nonMatch bool
		expect   []byte
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	cndInputMessage := tb.Condition("non-nil message", "input non-zero or non-nil message")
	cndInputNilMessage := tb.Condition("nil message", "input nil message")
	cndInputKey := tb.Condition("non-nil key", "input non-nil key")
	cndInputNilKey := tb.Condition("nil key", "input nil key")
	actCheckHash := tb.Action("check hash", "check that the expected hash is returned")
	actCheckMatch := tb.Action("check match", "check that the expected hash matched")
	actCheckNotMatch := tb.Action("check not match", "check that the expected hash does not matched")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"non-nil msg & non-nil key",
			[]string{cndInputMessage, cndInputKey},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				msg: []byte("abc"),
				key: []byte("test"),
			},
			&action{
				// The expected value is generated by openssl.
				// echo -n "abc" | openssl dgst -hmac "test" -blake2s256
				expect: hexMustDecode("66fc46e757a5a0dda945ae11ad69a49ca17c8ada7472b363bb1e803d275aea51"),
			},
		),
		gen(
			"non-nil msg & nil key",
			[]string{cndInputMessage, cndInputNilKey},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				msg: []byte("abc"),
				key: nil,
			},
			&action{
				// The expected value is generated by openssl.
				// echo -n "abc" | openssl dgst -hmac "" -blake2s256
				expect: hexMustDecode("4a79a31b10a8c34b7c89735b8088e6b810b4bad4f36d671b65ca2ea0914952bf"),
			},
		),
		gen(
			"nil msg & non-nil key",
			[]string{cndInputNilMessage, cndInputKey},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				msg: nil,
				key: []byte("test"),
			},
			&action{
				// The expected value is generated by openssl.
				// echo -n "" | openssl dgst -hmac "test" -blake2s256
				expect: hexMustDecode("b7bf9b2ed48e763db328f7ce5e3639cfe239475ca8949371eb20cfe118615b32"),
			},
		),
		gen(
			"nil msg & nil key",
			[]string{cndInputNilMessage, cndInputNilKey},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				msg: nil,
				key: nil,
			},
			&action{
				// The expected value is generated by openssl.
				// echo -n "" | openssl dgst -hmac "" -blake2s256
				expect: hexMustDecode("eaf4bb25938f4d20e72656bbbc7a9bf63c0c18537333c35bdb67db1402661acd"),
			},
		),
		gen(
			"wrong hash",
			[]string{cndInputMessage, cndInputKey},
			[]string{actCheckHash, actCheckNotMatch},
			&condition{
				msg: []byte("abc"),
				key: []byte("test"),
			},
			&action{
				// The first character is wrong.
				expect:   hexMustDecode("06fc46e757a5a0dda945ae11ad69a49ca17c8ada7472b363bb1e803d275aea51"),
				nonMatch: true,
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			out := mac.BLAKE2s_256(tt.C().msg, tt.C().key)
			testutil.Diff(t, mac.SizeBLAKE2s_256, len(out))
			testutil.Diff(t, mac.HashSize[k.HashAlg_BLAKE2s_256], len(out))

			if tt.A().nonMatch {
				if bytes.Equal(tt.A().expect, out) {
					t.Error("error unexpectedly matched", tt.A().expect, out)
				}
			} else {
				testutil.Diff(t, tt.A().expect, out)
			}

			// Check that the digest generated by HMAC is different from the plain mac.
			plainHash := hash.BLAKE2s_256(tt.C().msg)
			testutil.Diff(t, false, bytes.Equal(plainHash, out))
		})
	}
}

func TestBLAKE2b_256(t *testing.T) {
	type condition struct {
		msg []byte
		key []byte
	}

	type action struct {
		nonMatch bool
		expect   []byte
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	cndInputMessage := tb.Condition("non-nil message", "input non-zero or non-nil message")
	cndInputNilMessage := tb.Condition("nil message", "input nil message")
	cndInputKey := tb.Condition("non-nil key", "input non-nil key")
	cndInputNilKey := tb.Condition("nil key", "input nil key")
	actCheckHash := tb.Action("check hash", "check that the expected hash is returned")
	actCheckMatch := tb.Action("check match", "check that the expected hash matched")
	actCheckNotMatch := tb.Action("check not match", "check that the expected hash does not matched")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"non-nil msg & non-nil key",
			[]string{cndInputMessage, cndInputKey},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				msg: []byte("abc"),
				key: []byte("test"),
			},
			&action{
				// The expected value is generated by python.
				// >> import hmac, hashlib
				// >> def blake2b_256():
				// >>     return hashlib.blake2b(digest_size=32)
				// >> hmac.new(b"test",b"abc", blake2b_256).hexdigest()
				expect: hexMustDecode("fd3cb97f34e3b743fe1d3f46ae2c5bb11d339a445f5157e0874b2c567f315cd2"),
			},
		),
		gen(
			"non-nil msg & nil key",
			[]string{cndInputMessage, cndInputNilKey},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				msg: []byte("abc"),
				key: nil,
			},
			&action{
				// The expected value is generated by python.
				// >> import hmac, hashlib
				// >> def blake2b_256():
				// >>     return hashlib.blake2b(digest_size=32)
				// >> hmac.new(b"",b"abc", blake2b_256).hexdigest()
				expect: hexMustDecode("81e452d5d85b8c170e068cd17109f0a648147ad66b830d09fad5456bc80c5f91"),
			},
		),
		gen(
			"nil msg & non-nil key",
			[]string{cndInputNilMessage, cndInputKey},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				msg: nil,
				key: []byte("test"),
			},
			&action{
				// The expected value is generated by python.
				// >> import hmac, hashlib
				// >> def blake2b_256():
				// >>     return hashlib.blake2b(digest_size=32)
				// >> hmac.new(b"test",b"", blake2b_256).hexdigest()
				expect: hexMustDecode("31daadffed4dfda67bda7580d32bb1c3cb917b0eac020230dce21f7c5acf3409"),
			},
		),
		gen(
			"nil msg & nil key",
			[]string{cndInputNilMessage, cndInputNilKey},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				msg: nil,
				key: nil,
			},
			&action{
				// The expected value is generated by python.
				// >> import hmac, hashlib
				// >> def blake2b_256():
				// >>     return hashlib.blake2b(digest_size=32)
				// >> hmac.new(b"",b"", blake2b_256).hexdigest()
				expect: hexMustDecode("486b62b89b06365cf96f77c388e093b92aa774ba9eb7530cae6e68a3acbab9e8"),
			},
		),
		gen(
			"wrong hash",
			[]string{cndInputMessage, cndInputKey},
			[]string{actCheckHash, actCheckNotMatch},
			&condition{
				msg: []byte("abc"),
				key: []byte("test"),
			},
			&action{
				// The first character is wrong.
				expect:   hexMustDecode("0d3cb97f34e3b743fe1d3f46ae2c5bb11d339a445f5157e0874b2c567f315cd2"),
				nonMatch: true,
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			out := mac.BLAKE2b_256(tt.C().msg, tt.C().key)
			testutil.Diff(t, mac.SizeBLAKE2b_256, len(out))
			testutil.Diff(t, mac.HashSize[k.HashAlg_BLAKE2b_256], len(out))

			if tt.A().nonMatch {
				if bytes.Equal(tt.A().expect, out) {
					t.Error("error unexpectedly matched", tt.A().expect, out)
				}
			} else {
				testutil.Diff(t, tt.A().expect, out)
			}

			// Check that the digest generated by HMAC is different from the plain mac.
			plainHash := hash.BLAKE2b_256(tt.C().msg)
			testutil.Diff(t, false, bytes.Equal(plainHash, out))
		})
	}
}

func TestBLAKE2b_384(t *testing.T) {
	type condition struct {
		msg []byte
		key []byte
	}

	type action struct {
		nonMatch bool
		expect   []byte
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	cndInputMessage := tb.Condition("non-nil message", "input non-zero or non-nil message")
	cndInputNilMessage := tb.Condition("nil message", "input nil message")
	cndInputKey := tb.Condition("non-nil key", "input non-nil key")
	cndInputNilKey := tb.Condition("nil key", "input nil key")
	actCheckHash := tb.Action("check hash", "check that the expected hash is returned")
	actCheckMatch := tb.Action("check match", "check that the expected hash matched")
	actCheckNotMatch := tb.Action("check not match", "check that the expected hash does not matched")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"non-nil msg & non-nil key",
			[]string{cndInputMessage, cndInputKey},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				msg: []byte("abc"),
				key: []byte("test"),
			},
			&action{
				// The expected value is generated by python.
				// >> import hmac, hashlib
				// >> def blake2b_384():
				// >>     return hashlib.blake2b(digest_size=48)
				// >> hmac.new(b"test",b"abc", blake2b_384).hexdigest()
				expect: hexMustDecode("13b26e291775bcc09a5246703c608cd85c0b1e8cb061b8e05d24f324358d4e5fbd26722430b51ef7c4e889d0d47af912"),
			},
		),
		gen(
			"non-nil msg & nil key",
			[]string{cndInputMessage, cndInputNilKey},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				msg: []byte("abc"),
				key: nil,
			},
			&action{
				// The expected value is generated by python.
				// >> import hmac, hashlib
				// >> def blake2b_384():
				// >>     return hashlib.blake2b(digest_size=48)
				// >> hmac.new(b"",b"abc", blake2b_384).hexdigest()
				expect: hexMustDecode("4d44ea1084f2d41d5edd27241cb9d4c3227392357894b8607712166dad1e24b65bd0ad772e9b8243ab6caaaf3b987df5"),
			},
		),
		gen(
			"nil msg & non-nil key",
			[]string{cndInputNilMessage, cndInputKey},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				msg: nil,
				key: []byte("test"),
			},
			&action{
				// The expected value is generated by python.
				// >> import hmac, hashlib
				// >> def blake2b_384():
				// >>     return hashlib.blake2b(digest_size=48)
				// >> hmac.new(b"test",b"", blake2b_384).hexdigest()
				expect: hexMustDecode("725e60ec477427492897fd0c4b41a7132fe721af3bc076ecf19189c3a120ce22310dedbd6a1505e8948087a8784e6218"),
			},
		),
		gen(
			"nil msg & nil key",
			[]string{cndInputNilMessage, cndInputNilKey},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				msg: nil,
				key: nil,
			},
			&action{
				// The expected value is generated by python.
				// >> import hmac, hashlib
				// >> def blake2b_384():
				// >>     return hashlib.blake2b(digest_size=48)
				// >> hmac.new(b"",b"", blake2b_384).hexdigest()
				expect: hexMustDecode("6fd408e56ff442b97ee0fd03baed766940d1d305fe79d8ca821337ed43cdff27b0979b69a9bb4e3ed31ef56266f9ebf4"),
			},
		),
		gen(
			"wrong hash",
			[]string{cndInputMessage, cndInputKey},
			[]string{actCheckHash, actCheckNotMatch},
			&condition{
				msg: []byte("abc"),
				key: []byte("test"),
			},
			&action{
				// The first character is wrong.
				expect:   hexMustDecode("03b26e291775bcc09a5246703c608cd85c0b1e8cb061b8e05d24f324358d4e5fbd26722430b51ef7c4e889d0d47af912"),
				nonMatch: true,
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			out := mac.BLAKE2b_384(tt.C().msg, tt.C().key)
			testutil.Diff(t, mac.SizeBLAKE2b_384, len(out))
			testutil.Diff(t, mac.HashSize[k.HashAlg_BLAKE2b_384], len(out))

			if tt.A().nonMatch {
				if bytes.Equal(tt.A().expect, out) {
					t.Error("error unexpectedly matched", tt.A().expect, out)
				}
			} else {
				testutil.Diff(t, tt.A().expect, out)
			}

			// Check that the digest generated by HMAC is different from the plain mac.
			plainHash := hash.BLAKE2b_384(tt.C().msg)
			testutil.Diff(t, false, bytes.Equal(plainHash, out))
		})
	}
}

func TestBLAKE2b_512(t *testing.T) {
	type condition struct {
		msg []byte
		key []byte
	}

	type action struct {
		nonMatch bool
		expect   []byte
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	cndInputMessage := tb.Condition("non-nil message", "input non-zero or non-nil message")
	cndInputNilMessage := tb.Condition("nil message", "input nil message")
	cndInputKey := tb.Condition("non-nil key", "input non-nil key")
	cndInputNilKey := tb.Condition("nil key", "input nil key")
	actCheckHash := tb.Action("check hash", "check that the expected hash is returned")
	actCheckMatch := tb.Action("check match", "check that the expected hash matched")
	actCheckNotMatch := tb.Action("check not match", "check that the expected hash does not matched")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"non-nil msg & non-nil key",
			[]string{cndInputMessage, cndInputKey},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				msg: []byte("abc"),
				key: []byte("test"),
			},
			&action{
				// The expected value is generated by openssl.
				// echo -n "abc" | openssl dgst -hmac "test" -blake2b512
				expect: hexMustDecode("81b84c0e4e60ca6e65749b2ff9d8b4ab6eeda53d096cd3052aa7aaa9cdafe7db9b9595c05fd728390e337cfeabdac48e271d386bdb9858052e168747fb6ff06f"),
			},
		),
		gen(
			"non-nil msg & nil key",
			[]string{cndInputMessage, cndInputNilKey},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				msg: []byte("abc"),
				key: nil,
			},
			&action{
				// The expected value is generated by openssl.
				// echo -n "abc" | openssl dgst -hmac "" -blake2b512
				expect: hexMustDecode("c23d7eb81509be346bab43f60e2b889e86301b0c029b3843ebbccac591a047e704566cddda2f068c3ca504182a25e11ca9933911cd696eebe0f1df7ca8d98158"),
			},
		),
		gen(
			"nil msg & non-nil key",
			[]string{cndInputNilMessage, cndInputKey},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				msg: nil,
				key: []byte("test"),
			},
			&action{
				// The expected value is generated by openssl.
				// echo -n "" | openssl dgst -hmac "test" -blake2b512
				expect: hexMustDecode("114b9e0ac5ec8eaf901284427592afd65ae1a99cc0f639d79c1f9f7c939328596f4f56caa0321b42046dcf557730ed1ecbc6dbe546972d81c6a853dfd1eea99d"),
			},
		),
		gen(
			"nil msg & nil key",
			[]string{cndInputNilMessage, cndInputNilKey},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				msg: nil,
				key: nil,
			},
			&action{
				// The expected value is generated by openssl.
				// echo -n "" | openssl dgst -hmac "" -blake2b512
				expect: hexMustDecode("198cd2006f66ff83fbbd913f78aca2251caf4f19fe9475aade8cf2091b99a68466775177424f58286886cbae8229644cec747237d4b721735485e17372fdf59c"),
			},
		),
		gen(
			"wrong hash",
			[]string{cndInputMessage, cndInputKey},
			[]string{actCheckHash, actCheckNotMatch},
			&condition{
				msg: []byte("abc"),
				key: []byte("test"),
			},
			&action{
				// The first character is wrong.
				expect:   hexMustDecode("01b84c0e4e60ca6e65749b2ff9d8b4ab6eeda53d096cd3052aa7aaa9cdafe7db9b9595c05fd728390e337cfeabdac48e271d386bdb9858052e168747fb6ff06f"),
				nonMatch: true,
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			out := mac.BLAKE2b_512(tt.C().msg, tt.C().key)
			testutil.Diff(t, mac.SizeBLAKE2b_512, len(out))
			testutil.Diff(t, mac.HashSize[k.HashAlg_BLAKE2b_512], len(out))

			if tt.A().nonMatch {
				if bytes.Equal(tt.A().expect, out) {
					t.Error("error unexpectedly matched", tt.A().expect, out)
				}
			} else {
				testutil.Diff(t, tt.A().expect, out)
			}

			// Check that the digest generated by HMAC is different from the plain mac.
			plainHash := hash.BLAKE2b_512(tt.C().msg)
			testutil.Diff(t, false, bytes.Equal(plainHash, out))
		})
	}
}
