// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package hash_test

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"testing"

	k "github.com/aileron-gateway/aileron-gateway/apis/kernel"
	"github.com/aileron-gateway/aileron-gateway/kernel/hash"
	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
)

var (
	_ = hash.HashFunc(hash.SHA1)        // Check that the function satisfies the signature.
	_ = hash.HashFunc(hash.SHA224)      // Check that the function satisfies the signature.
	_ = hash.HashFunc(hash.SHA256)      // Check that the function satisfies the signature.
	_ = hash.HashFunc(hash.SHA384)      // Check that the function satisfies the signature.
	_ = hash.HashFunc(hash.SHA512)      // Check that the function satisfies the signature.
	_ = hash.HashFunc(hash.SHA512_224)  // Check that the function satisfies the signature.
	_ = hash.HashFunc(hash.SHA512_256)  // Check that the function satisfies the signature.
	_ = hash.HashFunc(hash.SHA3_224)    // Check that the function satisfies the signature.
	_ = hash.HashFunc(hash.SHA3_256)    // Check that the function satisfies the signature.
	_ = hash.HashFunc(hash.SHA3_384)    // Check that the function satisfies the signature.
	_ = hash.HashFunc(hash.SHA3_512)    // Check that the function satisfies the signature.
	_ = hash.HashFunc(hash.SHAKE128)    // Check that the function satisfies the signature.
	_ = hash.HashFunc(hash.SHAKE256)    // Check that the function satisfies the signature.
	_ = hash.HashFunc(hash.MD5)         // Check that the function satisfies the signature.
	_ = hash.HashFunc(hash.FNV1_32)     // Check that the function satisfies the signature.
	_ = hash.HashFunc(hash.FNV1a_32)    // Check that the function satisfies the signature.
	_ = hash.HashFunc(hash.FNV1_64)     // Check that the function satisfies the signature.
	_ = hash.HashFunc(hash.FNV1a_64)    // Check that the function satisfies the signature.
	_ = hash.HashFunc(hash.FNV1_128)    // Check that the function satisfies the signature.
	_ = hash.HashFunc(hash.FNV1a_128)   // Check that the function satisfies the signature.
	_ = hash.HashFunc(hash.CRC32)       // Check that the function satisfies the signature.
	_ = hash.HashFunc(hash.CRC64ISO)    // Check that the function satisfies the signature.
	_ = hash.HashFunc(hash.CRC64ECMA)   // Check that the function satisfies the signature.
	_ = hash.HashFunc(hash.BLAKE2s_256) // Check that the function satisfies the signature.
	_ = hash.HashFunc(hash.BLAKE2b_256) // Check that the function satisfies the signature.
	_ = hash.HashFunc(hash.BLAKE2b_384) // Check that the function satisfies the signature.
	_ = hash.HashFunc(hash.BLAKE2b_512) // Check that the function satisfies the signature.
)

func ExampleSHA1() {
	msg := []byte("example message")
	digest := hash.SHA1(msg)

	encoded := base64.StdEncoding.EncodeToString(digest)
	fmt.Println(len(digest), encoded)
	// Output:
	// 20 KR61XDC/qdsHp3J0r6vvz9j0Wls=
}

func ExampleSHA224() {
	msg := []byte("example message")
	digest := hash.SHA224(msg)

	encoded := base64.StdEncoding.EncodeToString(digest)
	fmt.Println(len(digest), encoded)
	// Output:
	// 28 XJfEMzdcEZ9s4fC1EZmIIIsiI8i1FesHY5wClw==
}

func ExampleSHA256() {
	msg := []byte("example message")
	digest := hash.SHA256(msg)

	encoded := base64.StdEncoding.EncodeToString(digest)
	fmt.Println(len(digest), encoded)
	// Output:
	// 32 rYTNCxD8Aoc4lxsHgSSuwqDnxtmGo4G+CzhvMr7oh68=
}

func ExampleSHA384() {
	msg := []byte("example message")
	digest := hash.SHA384(msg)

	encoded := base64.StdEncoding.EncodeToString(digest)
	fmt.Println(len(digest), encoded)
	// Output:
	// 48 5ay2kbj4b33/3/5o9tDD58hAvNVh5s6BcxC3F0UQww9KwJOwpVvMZxlshKpvwfiR
}

func ExampleSHA512() {
	msg := []byte("example message")
	digest := hash.SHA512(msg)

	encoded := base64.StdEncoding.EncodeToString(digest)
	fmt.Println(len(digest), encoded)
	// Output:
	// 64 lhhV4GZKBv3ASbCUoXZKp2vhyLuz+VBvZKZpfwxSw/G2WpE0mfSq0Qk3AoAjRnoe0Vf1ZGOvASBR6yjHT2PlvA==
}

// TODO: Remove underscore.
// func ExampleSHA512_224() {
// 	msg := []byte("example message")
// 	digest := hash.SHA512_224(msg)

// 	encoded := base64.StdEncoding.EncodeToString(digest)
// 	fmt.Println(len(digest), encoded)
// 	// Output:
// 	// 28 5s1q7ygDA5ldyiKf/BKtbvV1FB9q7SHBbf3xsQ==
// }

// TODO: Remove underscore.
// func ExampleSHA512_256() {
// 	msg := []byte("example message")
// 	digest := hash.SHA512_256(msg)

// 	encoded := base64.StdEncoding.EncodeToString(digest)
// 	fmt.Println(len(digest), encoded)
// 	// Output:
// 	// 32 xjrf+BvbxNIeOTXUlS094J4CZx0hl8Fkg1CJrm7ta9s=
// }

// TODO: Remove underscore.
// func ExampleSHA3_224() {
// 	msg := []byte("example message")
// 	digest := hash.SHA3_224(msg)

// 	encoded := base64.StdEncoding.EncodeToString(digest)
// 	fmt.Println(len(digest), encoded)
// 	// Output:
// 	// 28 wRtKqBTSRNjkFX0MJOaHgBWZOksgC+b+wkEHfg==
// }

// TODO: Remove underscore.
// func ExampleSHA3_256() {
// 	msg := []byte("example message")
// 	digest := hash.SHA3_256(msg)

// 	encoded := base64.StdEncoding.EncodeToString(digest)
// 	fmt.Println(len(digest), encoded)
// 	// Output:
// 	// 32 d1+c8o4B0hBF2VvWSJbWD8NmroU2r4ybKZip9CRMVbM=
// }

// TODO: Remove underscore.
// func ExampleSHA3_384() {
// 	msg := []byte("example message")
// 	digest := hash.SHA3_384(msg)

// 	encoded := base64.StdEncoding.EncodeToString(digest)
// 	fmt.Println(len(digest), encoded)
// 	// Output:
// 	// 48 gsp6zrQeI9FkBFMiPGrkLYXKNeVYALNL+NF1u1j0f32eCYC5ncVVgTBS478l6IVO
// }

// TODO: Remove underscore.
// func ExampleSHA3_512() {
// 	msg := []byte("example message")
// 	digest := hash.SHA3_512(msg)

// 	encoded := base64.StdEncoding.EncodeToString(digest)
// 	fmt.Println(len(digest), encoded)
// 	// Output:
// 	// 64 zwA1lRotJWRBMf1X4I8TD2qu2n91r+Bl5WENP6OoIsLgQZwamQZ9XnxF2q3aCFBCJrPX3q7qoSrwNGEQvOI1xA==
// }

func ExampleSHAKE128() {
	msg := []byte("example message")
	digest := hash.SHAKE128(msg)

	encoded := base64.StdEncoding.EncodeToString(digest)
	fmt.Println(len(digest), encoded)
	// Output:
	// 32 jr+6If2cdlgjcLSyO9okbFKSs1LTpoAS4zvbP31hmcI=
}

func ExampleSHAKE256() {
	msg := []byte("example message")
	digest := hash.SHAKE256(msg)

	encoded := base64.StdEncoding.EncodeToString(digest)
	fmt.Println(len(digest), encoded)
	// Output:
	// 64 Q3EQDoBvE7j+LkW1yr68C/ZIkgfYStRyFUwTA1+UXGZW0VZkn6yDX9ebpZUJNlK1whk39MkGjJzczWHyiMWtMQ==
}

func ExampleMD5() {
	msg := []byte("example message")
	digest := hash.MD5(msg)

	encoded := base64.StdEncoding.EncodeToString(digest)
	fmt.Println(len(digest), encoded)
	// Output:
	// 16 XH2d+PYe8Ow9sFrCwFIlvg==
}

// TODO: Remove underscore.
// func ExampleFNV1_32() {
// 	msg := []byte("example message")
// 	digest := hash.FNV1_32(msg)

// 	encoded := base64.StdEncoding.EncodeToString(digest)
// 	fmt.Println(len(digest), encoded)
// 	// Output:
// 	// 4 ITiJjA==
// }

// TODO: Remove underscore.
// func ExampleFNV1a_32() {
// 	msg := []byte("example message")
// 	digest := hash.FNV1a_32(msg)

// 	encoded := base64.StdEncoding.EncodeToString(digest)
// 	fmt.Println(len(digest), encoded)
// 	// Output:
// 	// 4 YlvGGg==
// }

// TODO: Remove underscore.
// func ExampleFNV1_64() {
// 	msg := []byte("example message")
// 	digest := hash.FNV1_64(msg)

// 	encoded := base64.StdEncoding.EncodeToString(digest)
// 	fmt.Println(len(digest), encoded)
// 	// Output:
// 	// 8 KqBkJxMgXIw=
// }

// TODO: Remove underscore.
// func ExampleFNV1a_64() {
// 	msg := []byte("example message")
// 	digest := hash.FNV1a_64(msg)

// 	encoded := base64.StdEncoding.EncodeToString(digest)
// 	fmt.Println(len(digest), encoded)
// 	// Output:
// 	// 8 kMjZ9d7ErPo=
// }

// TODO: Remove underscore.
// func ExampleFNV1_128() {
// 	msg := []byte("example message")
// 	digest := hash.FNV1_128(msg)

// 	encoded := base64.StdEncoding.EncodeToString(digest)
// 	fmt.Println(len(digest), encoded)
// 	// Output:
// 	// 16 c00vazjIQklF3y7clzcujA==
// }

// TODO: Remove underscore.
// func ExampleFNV1a_128() {
// 	msg := []byte("example message")
// 	digest := hash.FNV1a_128(msg)

// 	encoded := base64.StdEncoding.EncodeToString(digest)
// 	fmt.Println(len(digest), encoded)
// 	// Output:
// 	// 16 h6ApEo0Rh1k2OrSz39ZUAg==
// }

func ExampleCRC32() {
	msg := []byte("example message")
	digest := hash.CRC32(msg)

	encoded := base64.StdEncoding.EncodeToString(digest)
	fmt.Println(len(digest), encoded)
	// Output:
	// 4 627EcA==
}

func ExampleCRC64ISO() {
	msg := []byte("example message")
	digest := hash.CRC64ISO(msg)

	encoded := base64.StdEncoding.EncodeToString(digest)
	fmt.Println(len(digest), encoded)
	// Output:
	// 8 3xYFy42wOts=
}

func ExampleCRC64ECMA() {
	msg := []byte("example message")
	digest := hash.CRC64ECMA(msg)

	encoded := base64.StdEncoding.EncodeToString(digest)
	fmt.Println(len(digest), encoded)
	// Output:
	// 8 UGgb3N/wvgk=
}

// TODO: Remove underscore.
// func ExampleBLAKE2s_256() {
// 	msg := []byte("example message")
// 	digest := hash.BLAKE2s_256(msg)

// 	encoded := base64.StdEncoding.EncodeToString(digest)
// 	fmt.Println(len(digest), encoded)
// 	// Output:
// 	// 32 ue0Z/vperOCkrXA3AP8dlVj6i2g2xN674wFsnmopTY4=
// }

// TODO: Remove underscore.
// func ExampleBLAKE2b_256() {
// 	msg := []byte("example message")
// 	digest := hash.BLAKE2b_256(msg)

// 	encoded := base64.StdEncoding.EncodeToString(digest)
// 	fmt.Println(len(digest), encoded)
// 	// Output:
// 	// 32 M21hHgLK+KHPi3FlQlNOBOjw/R2i5+Q1lwogbDj5oJ0=
// }

// TODO: Remove underscore.
// func ExampleBLAKE2b_384() {
// 	msg := []byte("example message")
// 	digest := hash.BLAKE2b_384(msg)

// 	encoded := base64.StdEncoding.EncodeToString(digest)
// 	fmt.Println(len(digest), encoded)
// 	// Output:
// 	// 48 d6/ampcIiOpnkVb/r9LBRnWcPrzy/N0utbpZkcj6tKSRTbmZCtrgNcjr8arepU28
// }

// TODO: Remove underscore.
// func ExampleBLAKE2b_512() {
// 	msg := []byte("example message")
// 	digest := hash.BLAKE2b_512(msg)

// 	encoded := base64.StdEncoding.EncodeToString(digest)
// 	fmt.Println(len(digest), encoded)
// 	// Output:
// 	// 64 VadbqHfocIQf4UKzBATW6VF7DBvo0gfz7xGaVYkjyqdcwphJPNPa6DtiGYlpoPS2QevIsaEzvkT0pccqNgcagA==
// }

func TestFromAlgorithm(t *testing.T) {
	type condition struct {
		alg int
	}

	type action struct {
		f hash.HashFunc
	}

	CndHashExists := "Hash exists"
	actCheckNil := "nil"

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	tb.Condition(CndHashExists, "give an existing algorithm id")
	tb.Action(actCheckNil, "check that nil was returned")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"SHA1",
			[]string{CndHashExists},
			[]string{},
			&condition{
				alg: int(hash.AlgSHA1),
			},
			&action{
				f: hash.SHA1,
			},
		),
		gen(
			"SHA224",
			[]string{CndHashExists},
			[]string{},
			&condition{
				alg: int(hash.AlgSHA224),
			},
			&action{
				f: hash.SHA224,
			},
		),
		gen(
			"SHA256",
			[]string{CndHashExists},
			[]string{},
			&condition{
				alg: int(hash.AlgSHA256),
			},
			&action{
				f: hash.SHA256,
			},
		),
		gen(
			"SHA384",
			[]string{CndHashExists},
			[]string{},
			&condition{
				alg: int(hash.AlgSHA384),
			},
			&action{
				f: hash.SHA384,
			},
		),
		gen(
			"SHA512",
			[]string{CndHashExists},
			[]string{},
			&condition{
				alg: int(hash.AlgSHA512),
			},
			&action{
				f: hash.SHA512,
			},
		),
		gen(
			"SHA512_224",
			[]string{CndHashExists},
			[]string{},
			&condition{
				alg: int(hash.AlgSHA512_224),
			},
			&action{
				f: hash.SHA512_224,
			},
		),
		gen(
			"SHA512_256",
			[]string{CndHashExists},
			[]string{},
			&condition{
				alg: int(hash.AlgSHA512_256),
			},
			&action{
				f: hash.SHA512_256,
			},
		),
		gen(
			"SHA3_224",
			[]string{CndHashExists},
			[]string{},
			&condition{
				alg: int(hash.AlgSHA3_224),
			},
			&action{
				f: hash.SHA3_224,
			},
		),
		gen(
			"SHA3_256",
			[]string{CndHashExists},
			[]string{},
			&condition{
				alg: int(hash.AlgSHA3_256),
			},
			&action{
				f: hash.SHA3_256,
			},
		),
		gen(
			"SHA3_384",
			[]string{CndHashExists},
			[]string{},
			&condition{
				alg: int(hash.AlgSHA3_384),
			},
			&action{
				f: hash.SHA3_384,
			},
		),
		gen(
			"SHA3_512",
			[]string{CndHashExists},
			[]string{},
			&condition{
				alg: int(hash.AlgSHA3_512),
			},
			&action{
				f: hash.SHA3_512,
			},
		),
		gen(
			"SHAKE128",
			[]string{CndHashExists},
			[]string{},
			&condition{
				alg: int(hash.AlgSHAKE128),
			},
			&action{
				f: hash.SHAKE128,
			},
		),
		gen(
			"SHAKE256",
			[]string{CndHashExists},
			[]string{},
			&condition{
				alg: int(hash.AlgSHAKE256),
			},
			&action{
				f: hash.SHAKE256,
			},
		),
		gen(
			"MD5",
			[]string{CndHashExists},
			[]string{},
			&condition{
				alg: int(hash.AlgMD5),
			},
			&action{
				f: hash.MD5,
			},
		),
		gen(
			"FNV1_32",
			[]string{CndHashExists},
			[]string{},
			&condition{
				alg: int(hash.AlgFNV1_32),
			},
			&action{
				f: hash.FNV1_32,
			},
		),
		gen(
			"FNV1a_32",
			[]string{CndHashExists},
			[]string{},
			&condition{
				alg: int(hash.AlgFNV1a_32),
			},
			&action{
				f: hash.FNV1a_32,
			},
		),
		gen(
			"FNV1_64",
			[]string{CndHashExists},
			[]string{},
			&condition{
				alg: int(hash.AlgFNV1_64),
			},
			&action{
				f: hash.FNV1_64,
			},
		),
		gen(
			"FNV1a_64",
			[]string{CndHashExists},
			[]string{},
			&condition{
				alg: int(hash.AlgFNV1a_64),
			},
			&action{
				f: hash.FNV1a_64,
			},
		),
		gen(
			"FNV1_128",
			[]string{CndHashExists},
			[]string{},
			&condition{
				alg: int(hash.AlgFNV1_128),
			},
			&action{
				f: hash.FNV1_128,
			},
		),
		gen(
			"FNV1a_128",
			[]string{CndHashExists},
			[]string{},
			&condition{
				alg: int(hash.AlgFNV1a_128),
			},
			&action{
				f: hash.FNV1a_128,
			},
		),
		gen(
			"BLAKE2s_256",
			[]string{CndHashExists},
			[]string{},
			&condition{
				alg: int(hash.AlgBLAKE2s_256),
			},
			&action{
				f: hash.BLAKE2s_256,
			},
		),
		gen(
			"BLAKE2b_256",
			[]string{CndHashExists},
			[]string{},
			&condition{
				alg: int(hash.AlgBLAKE2b_256),
			},
			&action{
				f: hash.BLAKE2b_256,
			},
		),
		gen(
			"BLAKE2b_384",
			[]string{CndHashExists},
			[]string{},
			&condition{
				alg: int(hash.AlgBLAKE2b_384),
			},
			&action{
				f: hash.BLAKE2b_384,
			},
		),
		gen(
			"BLAKE2b_512",
			[]string{CndHashExists},
			[]string{},
			&condition{
				alg: int(hash.AlgBLAKE2b_512),
			},
			&action{
				f: hash.BLAKE2b_512,
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
			h := hash.FromAlgorithm(hash.Algorithm(tt.C().alg))
			if tt.A().f == nil {
				testutil.Diff(t, tt.A().f, h)
				return
			}

			msg := []byte("example message")
			testutil.Diff(t, tt.A().f(msg), h(msg))
		})
	}
}

func TestFromHashAlg(t *testing.T) {
	type condition struct {
		typ int32
	}

	type action struct {
		f hash.HashFunc
	}

	CndHashExists := "Hash exists"
	actCheckNil := "error"

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	tb.Condition(CndHashExists, "give an existing algorithm id")
	tb.Action(actCheckNil, "check that nil was returned")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"SHA1",
			[]string{CndHashExists},
			[]string{},
			&condition{
				typ: int32(k.HashAlg_SHA1),
			},
			&action{
				f: hash.SHA1,
			},
		),
		gen(
			"SHA224",
			[]string{CndHashExists},
			[]string{},
			&condition{
				typ: int32(k.HashAlg_SHA224),
			},
			&action{
				f: hash.SHA224,
			},
		),
		gen(
			"SHA256",
			[]string{CndHashExists},
			[]string{},
			&condition{
				typ: int32(k.HashAlg_SHA256),
			},
			&action{
				f: hash.SHA256,
			},
		),
		gen(
			"SHA384",
			[]string{CndHashExists},
			[]string{},
			&condition{
				typ: int32(k.HashAlg_SHA384),
			},
			&action{
				f: hash.SHA384,
			},
		),
		gen(
			"SHA512",
			[]string{CndHashExists},
			[]string{},
			&condition{
				typ: int32(k.HashAlg_SHA512),
			},
			&action{
				f: hash.SHA512,
			},
		),
		gen(
			"SHA512_224",
			[]string{CndHashExists},
			[]string{},
			&condition{
				typ: int32(k.HashAlg_SHA512_224),
			},
			&action{
				f: hash.SHA512_224,
			},
		),
		gen(
			"SHA512_256",
			[]string{CndHashExists},
			[]string{},
			&condition{
				typ: int32(k.HashAlg_SHA512_256),
			},
			&action{
				f: hash.SHA512_256,
			},
		),
		gen(
			"SHA3_224",
			[]string{CndHashExists},
			[]string{},
			&condition{
				typ: int32(k.HashAlg_SHA3_224),
			},
			&action{
				f: hash.SHA3_224,
			},
		),
		gen(
			"SHA3_256",
			[]string{CndHashExists},
			[]string{},
			&condition{
				typ: int32(k.HashAlg_SHA3_256),
			},
			&action{
				f: hash.SHA3_256,
			},
		),
		gen(
			"SHA3_384",
			[]string{CndHashExists},
			[]string{},
			&condition{
				typ: int32(k.HashAlg_SHA3_384),
			},
			&action{
				f: hash.SHA3_384,
			},
		),
		gen(
			"SHA3_512",
			[]string{CndHashExists},
			[]string{},
			&condition{
				typ: int32(k.HashAlg_SHA3_512),
			},
			&action{
				f: hash.SHA3_512,
			},
		),
		gen(
			"SHAKE128",
			[]string{CndHashExists},
			[]string{},
			&condition{
				typ: int32(k.HashAlg_SHAKE128),
			},
			&action{
				f: hash.SHAKE128,
			},
		),
		gen(
			"SHAKE256",
			[]string{CndHashExists},
			[]string{},
			&condition{
				typ: int32(k.HashAlg_SHAKE256),
			},
			&action{
				f: hash.SHAKE256,
			},
		),
		gen(
			"MD5",
			[]string{CndHashExists},
			[]string{},
			&condition{
				typ: int32(k.HashAlg_MD5),
			},
			&action{
				f: hash.MD5,
			},
		),
		gen(
			"FNV1_32",
			[]string{CndHashExists},
			[]string{},
			&condition{
				typ: int32(k.HashAlg_FNV1_32),
			},
			&action{
				f: hash.FNV1_32,
			},
		),
		gen(
			"FNV1a_32",
			[]string{CndHashExists},
			[]string{},
			&condition{
				typ: int32(k.HashAlg_FNV1a_32),
			},
			&action{
				f: hash.FNV1a_32,
			},
		),
		gen(
			"FNV1_64",
			[]string{CndHashExists},
			[]string{},
			&condition{
				typ: int32(k.HashAlg_FNV1_64),
			},
			&action{
				f: hash.FNV1_64,
			},
		),
		gen(
			"FNV1a_64",
			[]string{CndHashExists},
			[]string{},
			&condition{
				typ: int32(k.HashAlg_FNV1a_64),
			},
			&action{
				f: hash.FNV1a_64,
			},
		),
		gen(
			"FNV1_128",
			[]string{CndHashExists},
			[]string{},
			&condition{
				typ: int32(k.HashAlg_FNV1_128),
			},
			&action{
				f: hash.FNV1_128,
			},
		),
		gen(
			"FNV1a_128",
			[]string{CndHashExists},
			[]string{},
			&condition{
				typ: int32(k.HashAlg_FNV1a_128),
			},
			&action{
				f: hash.FNV1a_128,
			},
		),
		gen(
			"BLAKE2s_256",
			[]string{CndHashExists},
			[]string{},
			&condition{
				typ: int32(k.HashAlg_BLAKE2s_256),
			},
			&action{
				f: hash.BLAKE2s_256,
			},
		),
		gen(
			"BLAKE2b_256",
			[]string{CndHashExists},
			[]string{},
			&condition{
				typ: int32(k.HashAlg_BLAKE2b_256),
			},
			&action{
				f: hash.BLAKE2b_256,
			},
		),
		gen(
			"BLAKE2b_384",
			[]string{CndHashExists},
			[]string{},
			&condition{
				typ: int32(k.HashAlg_BLAKE2b_384),
			},
			&action{
				f: hash.BLAKE2b_384,
			},
		),
		gen(
			"BLAKE2b_512",
			[]string{CndHashExists},
			[]string{},
			&condition{
				typ: int32(k.HashAlg_BLAKE2b_512),
			},
			&action{
				f: hash.BLAKE2b_512,
			},
		),
		gen(
			"unknown",
			[]string{},
			[]string{actCheckNil},
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
			[]string{actCheckNil},
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
			h := hash.FromHashAlg(k.HashAlg(tt.C().typ))
			if tt.A().f == nil {
				testutil.Diff(t, tt.A().f, h)
				return
			}

			msg := []byte("example message")
			testutil.Diff(t, tt.A().f(msg), h(msg))
		})
	}
}

func hexMustDecode(s string) []byte {
	b, _ := hex.DecodeString(s)
	return b
}

func TestSHA1(t *testing.T) {
	type condition struct {
		input    []byte
		notMatch bool
	}

	type action struct {
		expect []byte
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	cndInputMessage := tb.Condition("non-nil message", "input non-zero or non-nil message")
	cndInputNil := tb.Condition("nil message", "input nil as message")
	actCheckHash := tb.Action("check hash", "check that the expected hash is returned")
	actCheckMatch := tb.Action("check match", "check that the expected hash matched")
	actCheckNotMatch := tb.Action("check not match", "check that the expected hash does not matched")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"nil",
			[]string{cndInputNil},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				input: nil,
			},
			&action{
				// The expected value is generated by openssl.
				// echo -n "" | openssl dgst -sha1
				expect: hexMustDecode("da39a3ee5e6b4b0d3255bfef95601890afd80709"),
			},
		),
		gen(
			"zero length",
			[]string{cndInputMessage},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				input: []byte{},
			},
			&action{
				// The expected value is generated by openssl.
				// echo -n "" | openssl dgst -sha1
				expect: hexMustDecode("da39a3ee5e6b4b0d3255bfef95601890afd80709"),
			},
		),
		gen(
			"non-zero length 1",
			[]string{cndInputMessage},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				input: []byte("abc"),
			},
			&action{
				// The expected value is generated by openssl.
				// echo -n "abc" | openssl dgst -sha1
				expect: hexMustDecode("a9993e364706816aba3e25717850c26c9cd0d89d"),
			},
		),
		gen(
			"non-zero length 2",
			[]string{cndInputMessage},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				input: []byte("test"),
			},
			&action{
				// The expected value is generated by openssl.
				// echo -n "test" | openssl dgst -sha1
				expect: hexMustDecode("a94a8fe5ccb19ba61c4c0873d391e987982fbbd3"),
			},
		),
		gen(
			"not match",
			[]string{cndInputMessage},
			[]string{actCheckHash, actCheckNotMatch},
			&condition{
				input:    []byte("test"),
				notMatch: true,
			},
			&action{
				// The first character is wrong.
				expect: hexMustDecode("094a8fe5ccb19ba61c4c0873d391e987982fbbd3"),
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			out := hash.SHA1(tt.C().input)
			testutil.Diff(t, hash.SizeSHA1, len(out))
			testutil.Diff(t, hash.HashSize[k.HashAlg_SHA1], len(out))

			if tt.C().notMatch {
				if bytes.Equal(tt.A().expect, out) {
					t.Error("error unexpectedly matched", tt.A().expect, out)
				}
			} else {
				testutil.Diff(t, tt.A().expect, out)
			}
		})
	}
}

func TestSHA224(t *testing.T) {
	type condition struct {
		input    []byte
		notMatch bool
	}

	type action struct {
		expect []byte
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	cndInputMessage := tb.Condition("non-nil message", "input non-zero or non-nil message")
	cndInputNil := tb.Condition("nil message", "input nil as message")
	actCheckHash := tb.Action("check hash", "check that the expected hash is returned")
	actCheckMatch := tb.Action("check match", "check that the expected hash matched")
	actCheckNotMatch := tb.Action("check not match", "check that the expected hash does not matched")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"nil",
			[]string{cndInputNil},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				input: nil,
			},
			&action{
				// The expected value is generated by openssl.
				// echo -n "" | openssl dgst -sha224
				expect: hexMustDecode("d14a028c2a3a2bc9476102bb288234c415a2b01f828ea62ac5b3e42f"),
			},
		),
		gen(
			"zero length",
			[]string{cndInputMessage},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				input: []byte{},
			},
			&action{
				// The expected value is generated by openssl.
				// echo -n "" | openssl dgst -sha224
				expect: hexMustDecode("d14a028c2a3a2bc9476102bb288234c415a2b01f828ea62ac5b3e42f"),
			},
		),
		gen(
			"non-zero length 1",
			[]string{cndInputMessage},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				input: []byte("abc"),
			},
			&action{
				// The expected value is generated by openssl.
				// echo -n "abc" | openssl dgst -sha224
				expect: hexMustDecode("23097d223405d8228642a477bda255b32aadbce4bda0b3f7e36c9da7"),
			},
		),
		gen(
			"non-zero length 2",
			[]string{cndInputMessage},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				input: []byte("test"),
			},
			&action{
				// The expected value is generated by openssl.
				// echo -n "test" | openssl dgst -sha224
				expect: hexMustDecode("90a3ed9e32b2aaf4c61c410eb925426119e1a9dc53d4286ade99a809"),
			},
		),
		gen(
			"not match",
			[]string{cndInputMessage},
			[]string{actCheckHash, actCheckNotMatch},
			&condition{
				input:    []byte("test"),
				notMatch: true,
			},
			&action{
				// The first character is wrong.
				expect: hexMustDecode("00a3ed9e32b2aaf4c61c410eb925426119e1a9dc53d4286ade99a809"),
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			out := hash.SHA224(tt.C().input)
			testutil.Diff(t, hash.SizeSHA224, len(out))
			testutil.Diff(t, hash.HashSize[k.HashAlg_SHA224], len(out))

			if tt.C().notMatch {
				if bytes.Equal(tt.A().expect, out) {
					t.Error("error unexpectedly matched", tt.A().expect, out)
				}
			} else {
				testutil.Diff(t, tt.A().expect, out)
			}
		})
	}
}

func TestSHA256(t *testing.T) {
	type condition struct {
		input    []byte
		notMatch bool
	}

	type action struct {
		expect []byte
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	cndInputMessage := tb.Condition("non-nil message", "input non-zero or non-nil message")
	cndInputNil := tb.Condition("nil message", "input nil as message")
	actCheckHash := tb.Action("check hash", "check that the expected hash is returned")
	actCheckMatch := tb.Action("check match", "check that the expected hash matched")
	actCheckNotMatch := tb.Action("check not match", "check that the expected hash does not matched")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"nil",
			[]string{cndInputNil},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				input: nil,
			},
			&action{
				// The expected value is generated by openssl.
				// echo -n "" | openssl dgst -sha256
				expect: hexMustDecode("e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"),
			},
		),
		gen(
			"zero length",
			[]string{cndInputMessage},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				input: []byte{},
			},
			&action{
				// The expected value is generated by openssl.
				// echo -n "" | openssl dgst -sha256
				expect: hexMustDecode("e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"),
			},
		),
		gen(
			"non-zero length 1",
			[]string{cndInputMessage},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				input: []byte("abc"),
			},
			&action{
				// The expected value is generated by openssl.
				// echo -n "abc" | openssl dgst -sha256
				expect: hexMustDecode("ba7816bf8f01cfea414140de5dae2223b00361a396177a9cb410ff61f20015ad"),
			},
		),
		gen(
			"non-zero length 2",
			[]string{cndInputMessage},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				input: []byte("test"),
			},
			&action{
				// The expected value is generated by openssl.
				// echo -n "test" | openssl dgst -sha256
				expect: hexMustDecode("9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08"),
			},
		),
		gen(
			"not match",
			[]string{cndInputMessage},
			[]string{actCheckHash, actCheckNotMatch},
			&condition{
				input:    []byte("test"),
				notMatch: true,
			},
			&action{
				// The first character is wrong
				expect: hexMustDecode("0f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08"),
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			out := hash.SHA256(tt.C().input)
			testutil.Diff(t, hash.SizeSHA256, len(out))
			testutil.Diff(t, hash.HashSize[k.HashAlg_SHA256], len(out))

			if tt.C().notMatch {
				if bytes.Equal(tt.A().expect, out) {
					t.Error("error unexpectedly matched", tt.A().expect, out)
				}
			} else {
				testutil.Diff(t, tt.A().expect, out)
			}
		})
	}
}

func TestSHA384(t *testing.T) {
	type condition struct {
		input    []byte
		notMatch bool
	}

	type action struct {
		expect []byte
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	cndInputMessage := tb.Condition("non-nil message", "input non-zero or non-nil message")
	cndInputNil := tb.Condition("nil message", "input nil as message")
	actCheckHash := tb.Action("check hash", "check that the expected hash is returned")
	actCheckMatch := tb.Action("check match", "check that the expected hash matched")
	actCheckNotMatch := tb.Action("check not match", "check that the expected hash does not matched")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"nil",
			[]string{cndInputNil},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				input: nil,
			},
			&action{
				// The expected value is generated by openssl.
				// echo -n "" | openssl dgst -sha384
				expect: hexMustDecode("38b060a751ac96384cd9327eb1b1e36a21fdb71114be07434c0cc7bf63f6e1da274edebfe76f65fbd51ad2f14898b95b"),
			},
		),
		gen(
			"zero length",
			[]string{cndInputMessage},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				input: []byte{},
			},
			&action{
				// The expected value is generated by openssl.
				// echo -n "" | openssl dgst -sha384
				expect: hexMustDecode("38b060a751ac96384cd9327eb1b1e36a21fdb71114be07434c0cc7bf63f6e1da274edebfe76f65fbd51ad2f14898b95b"),
			},
		),
		gen(
			"non-zero length 1",
			[]string{cndInputMessage},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				input: []byte("abc"),
			},
			&action{
				// The expected value is generated by openssl.
				// echo -n "abc" | openssl dgst -sha384
				expect: hexMustDecode("cb00753f45a35e8bb5a03d699ac65007272c32ab0eded1631a8b605a43ff5bed8086072ba1e7cc2358baeca134c825a7"),
			},
		),
		gen(
			"non-zero length 2",
			[]string{cndInputMessage},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				input: []byte("test"),
			},
			&action{
				// The expected value is generated by openssl.
				// echo -n "test" | openssl dgst -sha384
				expect: hexMustDecode("768412320f7b0aa5812fce428dc4706b3cae50e02a64caa16a782249bfe8efc4b7ef1ccb126255d196047dfedf17a0a9"),
			},
		),
		gen(
			"not match",
			[]string{cndInputMessage},
			[]string{actCheckHash, actCheckNotMatch},
			&condition{
				input:    []byte("test"),
				notMatch: true,
			},
			&action{
				// The first character is wrong.
				expect: hexMustDecode("068412320f7b0aa5812fce428dc4706b3cae50e02a64caa16a782249bfe8efc4b7ef1ccb126255d196047dfedf17a0a9"),
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			out := hash.SHA384(tt.C().input)
			testutil.Diff(t, hash.SizeSHA384, len(out))
			testutil.Diff(t, hash.HashSize[k.HashAlg_SHA384], len(out))

			if tt.C().notMatch {
				if bytes.Equal(tt.A().expect, out) {
					t.Error("error unexpectedly matched", tt.A().expect, out)
				}
			} else {
				testutil.Diff(t, tt.A().expect, out)
			}
		})
	}
}

func TestSHA512(t *testing.T) {
	type condition struct {
		input    []byte
		notMatch bool
	}

	type action struct {
		expect []byte
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	cndInputMessage := tb.Condition("non-nil message", "input non-zero or non-nil message")
	cndInputNil := tb.Condition("nil message", "input nil as message")
	actCheckHash := tb.Action("check hash", "check that the expected hash is returned")
	actCheckMatch := tb.Action("check match", "check that the expected hash matched")
	actCheckNotMatch := tb.Action("check not match", "check that the expected hash does not matched")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"nil",
			[]string{cndInputNil},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				input: nil,
			},
			&action{
				// The expected value is generated by openssl.
				// echo -n "" | openssl dgst -sha512
				expect: hexMustDecode("cf83e1357eefb8bdf1542850d66d8007d620e4050b5715dc83f4a921d36ce9ce47d0d13c5d85f2b0ff8318d2877eec2f63b931bd47417a81a538327af927da3e"),
			},
		),
		gen(
			"zero length",
			[]string{cndInputMessage},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				input: []byte{},
			},
			&action{
				// The expected value is generated by openssl.
				// echo -n "" | openssl dgst -sha512
				expect: hexMustDecode("cf83e1357eefb8bdf1542850d66d8007d620e4050b5715dc83f4a921d36ce9ce47d0d13c5d85f2b0ff8318d2877eec2f63b931bd47417a81a538327af927da3e"),
			},
		),
		gen(
			"non-zero length 1",
			[]string{cndInputMessage},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				input: []byte("abc"),
			},
			&action{
				// The expected value is generated by openssl.
				// echo -n "abc" | openssl dgst -sha512
				expect: hexMustDecode("ddaf35a193617abacc417349ae20413112e6fa4e89a97ea20a9eeee64b55d39a2192992a274fc1a836ba3c23a3feebbd454d4423643ce80e2a9ac94fa54ca49f"),
			},
		),
		gen(
			"non-zero length 2",
			[]string{cndInputMessage},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				input: []byte("test"),
			},
			&action{
				// The expected value is generated by openssl.
				// echo -n "test" | openssl dgst -sha512
				expect: hexMustDecode("ee26b0dd4af7e749aa1a8ee3c10ae9923f618980772e473f8819a5d4940e0db27ac185f8a0e1d5f84f88bc887fd67b143732c304cc5fa9ad8e6f57f50028a8ff"),
			},
		),
		gen(
			"not match",
			[]string{cndInputMessage},
			[]string{actCheckHash, actCheckNotMatch},
			&condition{
				input:    []byte("test"),
				notMatch: true,
			},
			&action{
				// The first character is wrong.
				expect: hexMustDecode("0e26b0dd4af7e749aa1a8ee3c10ae9923f618980772e473f8819a5d4940e0db27ac185f8a0e1d5f84f88bc887fd67b143732c304cc5fa9ad8e6f57f50028a8ff"),
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			out := hash.SHA512(tt.C().input)
			testutil.Diff(t, hash.SizeSHA512, len(out))
			testutil.Diff(t, hash.HashSize[k.HashAlg_SHA512], len(out))

			if tt.C().notMatch {
				if bytes.Equal(tt.A().expect, out) {
					t.Error("error unexpectedly matched", tt.A().expect, out)
				}
			} else {
				testutil.Diff(t, tt.A().expect, out)
			}
		})
	}
}

func TestSHA512_224(t *testing.T) {
	type condition struct {
		input    []byte
		notMatch bool
	}

	type action struct {
		expect []byte
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	cndInputMessage := tb.Condition("non-nil message", "input non-zero or non-nil message")
	cndInputNil := tb.Condition("nil message", "input nil as message")
	actCheckHash := tb.Action("check hash", "check that the expected hash is returned")
	actCheckMatch := tb.Action("check match", "check that the expected hash matched")
	actCheckNotMatch := tb.Action("check not match", "check that the expected hash does not matched")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"nil",
			[]string{cndInputNil},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				input: nil,
			},
			&action{
				// The expected value is generated by openssl.
				// echo -n "" | openssl dgst -sha512-224
				expect: hexMustDecode("6ed0dd02806fa89e25de060c19d3ac86cabb87d6a0ddd05c333b84f4"),
			},
		),
		gen(
			"zero length",
			[]string{cndInputMessage},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				input: []byte{},
			},
			&action{
				// The expected value is generated by openssl.
				// echo -n "" | openssl dgst -sha512-224
				expect: hexMustDecode("6ed0dd02806fa89e25de060c19d3ac86cabb87d6a0ddd05c333b84f4"),
			},
		),
		gen(
			"non-zero length 1",
			[]string{cndInputMessage},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				input: []byte("abc"),
			},
			&action{
				// The expected value is generated by openssl.
				// echo -n "abc" | openssl dgst -sha512-224
				expect: hexMustDecode("4634270f707b6a54daae7530460842e20e37ed265ceee9a43e8924aa"),
			},
		),
		gen(
			"non-zero length 2",
			[]string{cndInputMessage},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				input: []byte("test"),
			},
			&action{
				// The expected value is generated by openssl.
				// echo -n "test" | openssl dgst -sha512-224
				expect: hexMustDecode("06001bf08dfb17d2b54925116823be230e98b5c6c278303bc4909a8c"),
			},
		),
		gen(
			"not match",
			[]string{cndInputMessage},
			[]string{actCheckHash, actCheckNotMatch},
			&condition{
				input:    []byte("test"),
				notMatch: true,
			},
			&action{
				// The first character is wrong.
				expect: hexMustDecode("16001bf08dfb17d2b54925116823be230e98b5c6c278303bc4909a8c"),
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			out := hash.SHA512_224(tt.C().input)
			testutil.Diff(t, hash.SizeSHA512_224, len(out))
			testutil.Diff(t, hash.HashSize[k.HashAlg_SHA512_224], len(out))

			if tt.C().notMatch {
				if bytes.Equal(tt.A().expect, out) {
					t.Error("error unexpectedly matched", tt.A().expect, out)
				}
			} else {
				testutil.Diff(t, tt.A().expect, out)
			}
		})
	}
}

func TestSHA512_256(t *testing.T) {
	type condition struct {
		input    []byte
		notMatch bool
	}

	type action struct {
		expect []byte
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	cndInputMessage := tb.Condition("non-nil message", "input non-zero or non-nil message")
	cndInputNil := tb.Condition("nil message", "input nil as message")
	actCheckHash := tb.Action("check hash", "check that the expected hash is returned")
	actCheckMatch := tb.Action("check match", "check that the expected hash matched")
	actCheckNotMatch := tb.Action("check not match", "check that the expected hash does not matched")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"nil",
			[]string{cndInputNil},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				input: nil,
			},
			&action{
				// The expected value is generated by openssl.
				// echo -n "" | openssl dgst -sha512-256
				expect: hexMustDecode("c672b8d1ef56ed28ab87c3622c5114069bdd3ad7b8f9737498d0c01ecef0967a"),
			},
		),
		gen(
			"zero length",
			[]string{cndInputMessage},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				input: []byte{},
			},
			&action{
				// The expected value is generated by openssl.
				// echo -n "" | openssl dgst -sha512-256
				expect: hexMustDecode("c672b8d1ef56ed28ab87c3622c5114069bdd3ad7b8f9737498d0c01ecef0967a"),
			},
		),
		gen(
			"non-zero length 1",
			[]string{cndInputMessage},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				input: []byte("abc"),
			},
			&action{
				// The expected value is generated by openssl.
				// echo -n "abc" | openssl dgst -sha512-256
				expect: hexMustDecode("53048e2681941ef99b2e29b76b4c7dabe4c2d0c634fc6d46e0e2f13107e7af23"),
			},
		),
		gen(
			"non-zero length 2",
			[]string{cndInputMessage},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				input: []byte("test"),
			},
			&action{
				// The expected value is generated by openssl.
				// echo -n "test" | openssl dgst -sha512-256
				expect: hexMustDecode("3d37fe58435e0d87323dee4a2c1b339ef954de63716ee79f5747f94d974f913f"),
			},
		),
		gen(
			"not match",
			[]string{cndInputMessage},
			[]string{actCheckHash, actCheckNotMatch},
			&condition{
				input:    []byte("test"),
				notMatch: true,
			},
			&action{
				// The first character is wrong.
				expect: hexMustDecode("0d37fe58435e0d87323dee4a2c1b339ef954de63716ee79f5747f94d974f913f"),
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			out := hash.SHA512_256(tt.C().input)
			testutil.Diff(t, hash.SizeSHA512_256, len(out))
			testutil.Diff(t, hash.HashSize[k.HashAlg_SHA512_256], len(out))

			if tt.C().notMatch {
				if bytes.Equal(tt.A().expect, out) {
					t.Error("error unexpectedly matched", tt.A().expect, out)
				}
			} else {
				testutil.Diff(t, tt.A().expect, out)
			}
		})
	}
}

func TestSHA3_224(t *testing.T) {
	type condition struct {
		input    []byte
		notMatch bool
	}

	type action struct {
		expect []byte
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	cndInputMessage := tb.Condition("non-nil message", "input non-zero or non-nil message")
	cndInputNil := tb.Condition("nil message", "input nil as message")
	actCheckHash := tb.Action("check hash", "check that the expected hash is returned")
	actCheckMatch := tb.Action("check match", "check that the expected hash matched")
	actCheckNotMatch := tb.Action("check not match", "check that the expected hash does not matched")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"nil",
			[]string{cndInputNil},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				input: nil,
			},
			&action{
				// The expected value is generated by openssl.
				// echo -n "" | openssl dgst -sha3-224
				expect: hexMustDecode("6b4e03423667dbb73b6e15454f0eb1abd4597f9a1b078e3f5b5a6bc7"),
			},
		),
		gen(
			"zero length",
			[]string{cndInputMessage},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				input: []byte{},
			},
			&action{
				// The expected value is generated by openssl.
				// echo -n "" | openssl dgst -sha3-224
				expect: hexMustDecode("6b4e03423667dbb73b6e15454f0eb1abd4597f9a1b078e3f5b5a6bc7"),
			},
		),
		gen(
			"non-zero length 1",
			[]string{cndInputMessage},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				input: []byte("abc"),
			},
			&action{
				// The expected value is generated by openssl.
				// echo -n "abc" | openssl dgst -sha3-224
				expect: hexMustDecode("e642824c3f8cf24ad09234ee7d3c766fc9a3a5168d0c94ad73b46fdf"),
			},
		),
		gen(
			"non-zero length 2",
			[]string{cndInputMessage},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				input: []byte("test"),
			},
			&action{
				// Generated by https://emn178.github.io/online-tools/sha3_224.html
				expect: hexMustDecode("3797bf0afbbfca4a7bbba7602a2b552746876517a7f9b7ce2db0ae7b"),
			},
		),
		gen(
			"not match",
			[]string{cndInputMessage},
			[]string{actCheckHash, actCheckNotMatch},
			&condition{
				input:    []byte("test"),
				notMatch: true,
			},
			&action{
				// The first character is wrong.
				expect: hexMustDecode("0797bf0afbbfca4a7bbba7602a2b552746876517a7f9b7ce2db0ae7b"),
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			out := hash.SHA3_224(tt.C().input)
			testutil.Diff(t, hash.SizeSHA3_224, len(out))
			testutil.Diff(t, hash.HashSize[k.HashAlg_SHA3_224], len(out))

			if tt.C().notMatch {
				if bytes.Equal(tt.A().expect, out) {
					t.Error("error unexpectedly matched", tt.A().expect, out)
				}
			} else {
				testutil.Diff(t, tt.A().expect, out)
			}
		})
	}
}

func TestSHA3_256(t *testing.T) {
	type condition struct {
		input    []byte
		notMatch bool
	}

	type action struct {
		expect []byte
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	cndInputMessage := tb.Condition("non-nil message", "input non-zero or non-nil message")
	cndInputNil := tb.Condition("nil message", "input nil as message")
	actCheckHash := tb.Action("check hash", "check that the expected hash is returned")
	actCheckMatch := tb.Action("check match", "check that the expected hash matched")
	actCheckNotMatch := tb.Action("check not match", "check that the expected hash does not matched")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"nil",
			[]string{cndInputNil},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				input: nil,
			},
			&action{
				// The expected value is generated by openssl.
				// echo -n "" | openssl dgst -sha3-256
				expect: hexMustDecode("a7ffc6f8bf1ed76651c14756a061d662f580ff4de43b49fa82d80a4b80f8434a"),
			},
		),
		gen(
			"zero length",
			[]string{cndInputMessage},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				input: []byte{},
			},
			&action{
				// The expected value is generated by openssl.
				// echo -n "" | openssl dgst -sha3-256
				expect: hexMustDecode("a7ffc6f8bf1ed76651c14756a061d662f580ff4de43b49fa82d80a4b80f8434a"),
			},
		),
		gen(
			"non-zero length 1",
			[]string{cndInputMessage},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				input: []byte("abc"),
			},
			&action{
				// The expected value is generated by openssl.
				// echo -n "abc" | openssl dgst -sha3-256
				expect: hexMustDecode("3a985da74fe225b2045c172d6bd390bd855f086e3e9d525b46bfe24511431532"),
			},
		),
		gen(
			"non-zero length 2",
			[]string{cndInputMessage},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				input: []byte("test"),
			},
			&action{
				// The expected value is generated by openssl.
				// echo -n "test" | openssl dgst -sha3-256
				expect: hexMustDecode("36f028580bb02cc8272a9a020f4200e346e276ae664e45ee80745574e2f5ab80"),
			},
		),
		gen(
			"not match",
			[]string{cndInputMessage},
			[]string{actCheckHash, actCheckNotMatch},
			&condition{
				input:    []byte("test"),
				notMatch: true,
			},
			&action{
				// The first character is wrong.
				expect: hexMustDecode("06f028580bb02cc8272a9a020f4200e346e276ae664e45ee80745574e2f5ab80"),
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			out := hash.SHA3_256(tt.C().input)
			testutil.Diff(t, hash.SizeSHA3_256, len(out))
			testutil.Diff(t, hash.HashSize[k.HashAlg_SHA3_256], len(out))

			if tt.C().notMatch {
				if bytes.Equal(tt.A().expect, out) {
					t.Error("error unexpectedly matched", tt.A().expect, out)
				}
			} else {
				testutil.Diff(t, tt.A().expect, out)
			}
		})
	}
}

func TestSHA3_384(t *testing.T) {
	type condition struct {
		input    []byte
		notMatch bool
	}

	type action struct {
		expect []byte
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	cndInputMessage := tb.Condition("non-nil message", "input non-zero or non-nil message")
	cndInputNil := tb.Condition("nil message", "input nil as message")
	actCheckHash := tb.Action("check hash", "check that the expected hash is returned")
	actCheckMatch := tb.Action("check match", "check that the expected hash matched")
	actCheckNotMatch := tb.Action("check not match", "check that the expected hash does not matched")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"nil",
			[]string{cndInputNil},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				input: nil,
			},
			&action{
				// The expected value is generated by openssl.
				// echo -n "" | openssl dgst -sha3-384
				expect: hexMustDecode("0c63a75b845e4f7d01107d852e4c2485c51a50aaaa94fc61995e71bbee983a2ac3713831264adb47fb6bd1e058d5f004"),
			},
		),
		gen(
			"zero length",
			[]string{cndInputMessage},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				input: []byte{},
			},
			&action{
				// The expected value is generated by openssl.
				// echo -n "" | openssl dgst -sha3-384
				expect: hexMustDecode("0c63a75b845e4f7d01107d852e4c2485c51a50aaaa94fc61995e71bbee983a2ac3713831264adb47fb6bd1e058d5f004"),
			},
		),
		gen(
			"non-zero length 1",
			[]string{cndInputMessage},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				input: []byte("abc"),
			},
			&action{
				// The expected value is generated by openssl.
				// echo -n "abc" | openssl dgst -sha3-384
				expect: hexMustDecode("ec01498288516fc926459f58e2c6ad8df9b473cb0fc08c2596da7cf0e49be4b298d88cea927ac7f539f1edf228376d25"),
			},
		),
		gen(
			"non-zero length 2",
			[]string{cndInputMessage},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				input: []byte("test"),
			},
			&action{
				// The expected value is generated by openssl.
				// echo -n "test" | openssl dgst -sha3-384
				expect: hexMustDecode("e516dabb23b6e30026863543282780a3ae0dccf05551cf0295178d7ff0f1b41eecb9db3ff219007c4e097260d58621bd"),
			},
		),
		gen(
			"not match",
			[]string{cndInputMessage},
			[]string{actCheckHash, actCheckNotMatch},
			&condition{
				input:    []byte("test"),
				notMatch: true,
			},
			&action{
				// The first character is wrong.
				expect: hexMustDecode("0516dabb23b6e30026863543282780a3ae0dccf05551cf0295178d7ff0f1b41eecb9db3ff219007c4e097260d58621bd"),
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			out := hash.SHA3_384(tt.C().input)
			testutil.Diff(t, hash.SizeSHA3_384, len(out))
			testutil.Diff(t, hash.HashSize[k.HashAlg_SHA3_384], len(out))

			if tt.C().notMatch {
				if bytes.Equal(tt.A().expect, out) {
					t.Error("error unexpectedly matched", tt.A().expect, out)
				}
			} else {
				testutil.Diff(t, tt.A().expect, out)
			}
		})
	}
}

func TestSHA3_512(t *testing.T) {
	type condition struct {
		input    []byte
		notMatch bool
	}

	type action struct {
		expect []byte
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	cndInputMessage := tb.Condition("non-nil message", "input non-zero or non-nil message")
	cndInputNil := tb.Condition("nil message", "input nil as message")
	actCheckHash := tb.Action("check hash", "check that the expected hash is returned")
	actCheckMatch := tb.Action("check match", "check that the expected hash matched")
	actCheckNotMatch := tb.Action("check not match", "check that the expected hash does not matched")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"nil",
			[]string{cndInputNil},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				input: nil,
			},
			&action{
				// The expected value is generated by openssl.
				// echo -n "" | openssl dgst -sha3-512
				expect: hexMustDecode("a69f73cca23a9ac5c8b567dc185a756e97c982164fe25859e0d1dcc1475c80a615b2123af1f5f94c11e3e9402c3ac558f500199d95b6d3e301758586281dcd26"),
			},
		),
		gen(
			"zero length",
			[]string{cndInputMessage},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				input: []byte{},
			},
			&action{
				// The expected value is generated by openssl.
				// echo -n "" | openssl dgst -sha3-512
				expect: hexMustDecode("a69f73cca23a9ac5c8b567dc185a756e97c982164fe25859e0d1dcc1475c80a615b2123af1f5f94c11e3e9402c3ac558f500199d95b6d3e301758586281dcd26"),
			},
		),
		gen(
			"non-zero length 1",
			[]string{cndInputMessage},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				input: []byte("abc"),
			},
			&action{
				// The expected value is generated by openssl.
				// echo -n "abc" | openssl dgst -sha3-512
				expect: hexMustDecode("b751850b1a57168a5693cd924b6b096e08f621827444f70d884f5d0240d2712e10e116e9192af3c91a7ec57647e3934057340b4cf408d5a56592f8274eec53f0"),
			},
		),
		gen(
			"non-zero length 2",
			[]string{cndInputMessage},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				input: []byte("test"),
			},
			&action{
				// The expected value is generated by openssl.
				// echo -n "test" | openssl dgst -sha3-512
				expect: hexMustDecode("9ece086e9bac491fac5c1d1046ca11d737b92a2b2ebd93f005d7b710110c0a678288166e7fbe796883a4f2e9b3ca9f484f521d0ce464345cc1aec96779149c14"),
			},
		),
		gen(
			"not match",
			[]string{cndInputMessage},
			[]string{actCheckHash, actCheckNotMatch},
			&condition{
				input:    []byte("test"),
				notMatch: true,
			},
			&action{
				// The first character is wrong.
				expect: hexMustDecode("0ece086e9bac491fac5c1d1046ca11d737b92a2b2ebd93f005d7b710110c0a678288166e7fbe796883a4f2e9b3ca9f484f521d0ce464345cc1aec96779149c14"),
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			out := hash.SHA3_512(tt.C().input)
			testutil.Diff(t, hash.SizeSHA3_512, len(out))
			testutil.Diff(t, hash.HashSize[k.HashAlg_SHA3_512], len(out))

			if tt.C().notMatch {
				if bytes.Equal(tt.A().expect, out) {
					t.Error("error unexpectedly matched", tt.A().expect, out)
				}
			} else {
				testutil.Diff(t, tt.A().expect, out)
			}
		})
	}
}

func TestSHAKE128(t *testing.T) {
	type condition struct {
		input    []byte
		notMatch bool
	}

	type action struct {
		expect []byte
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	cndInputMessage := tb.Condition("non-nil message", "input non-zero or non-nil message")
	cndInputNil := tb.Condition("nil message", "input nil as message")
	actCheckHash := tb.Action("check hash", "check that the expected hash is returned")
	actCheckMatch := tb.Action("check match", "check that the expected hash matched")
	actCheckNotMatch := tb.Action("check not match", "check that the expected hash does not matched")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"nil",
			[]string{cndInputNil},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				input: nil,
			},
			&action{
				// The expected value can be checked by openssl (openssl only shows first 16 bytes).
				// echo -n "" | openssl dgst -shake128
				expect: hexMustDecode("7f9c2ba4e88f827d616045507605853ed73b8093f6efbc88eb1a6eacfa66ef26"),
			},
		),
		gen(
			"zero length",
			[]string{cndInputMessage},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				input: []byte{},
			},
			&action{
				// The expected value can be checked by openssl (openssl only shows first 16 bytes).
				// echo -n "" | openssl dgst -shake128
				expect: hexMustDecode("7f9c2ba4e88f827d616045507605853ed73b8093f6efbc88eb1a6eacfa66ef26"),
			},
		),
		gen(
			"non-zero length 1",
			[]string{cndInputMessage},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				input: []byte("abc"),
			},
			&action{
				// The expected value can be checked by openssl (openssl only shows first 16 bytes).
				// echo -n "abc" | openssl dgst -shake128
				// Online tools like https://emn178.github.io/online-tools/shake_128.html
				// can be used for generating a value.
				expect: hexMustDecode("5881092dd818bf5cf8a3ddb793fbcba74097d5c526a6d35f97b83351940f2cc8"),
			},
		),
		gen(
			"non-zero length 2",
			[]string{cndInputMessage},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				input: []byte("test"),
			},
			&action{
				// The expected value can be checked by openssl (openssl only shows first 16 bytes).
				// echo -n "test" | openssl dgst -shake128
				// Online tools like https://emn178.github.io/online-tools/shake_128.html
				// can be used for generating a value.
				expect: hexMustDecode("d3b0aa9cd8b7255622cebc631e867d4093d6f6010191a53973c45fec9b07c774"),
			},
		),
		gen(
			"not match",
			[]string{cndInputMessage},
			[]string{actCheckHash, actCheckNotMatch},
			&condition{
				input:    []byte("test"),
				notMatch: true,
			},
			&action{
				// The first character is wrong.
				expect: hexMustDecode("03b0aa9cd8b7255622cebc631e867d4093d6f6010191a53973c45fec9b07c774"),
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			out := hash.SHAKE128(tt.C().input)
			testutil.Diff(t, hash.SizeSHAKE128, len(out))
			testutil.Diff(t, hash.HashSize[k.HashAlg_SHAKE128], len(out))

			if tt.C().notMatch {
				if bytes.Equal(tt.A().expect, out) {
					t.Error("error unexpectedly matched", tt.A().expect, out)
				}
			} else {
				testutil.Diff(t, tt.A().expect, out)
			}
		})
	}
}

func TestSHAKE256(t *testing.T) {
	type condition struct {
		input    []byte
		notMatch bool
	}

	type action struct {
		expect []byte
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	cndInputMessage := tb.Condition("non-nil message", "input non-zero or non-nil message")
	cndInputNil := tb.Condition("nil message", "input nil as message")
	actCheckHash := tb.Action("check hash", "check that the expected hash is returned")
	actCheckMatch := tb.Action("check match", "check that the expected hash matched")
	actCheckNotMatch := tb.Action("check not match", "check that the expected hash does not matched")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"nil",
			[]string{cndInputNil},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				input: nil,
			},
			&action{
				// The expected value can be checked by openssl (openssl only shows first 32 bytes).
				// echo -n "" | openssl dgst -shake256
				expect: hexMustDecode("46b9dd2b0ba88d13233b3feb743eeb243fcd52ea62b81b82b50c27646ed5762fd75dc4ddd8c0f200cb05019d67b592f6fc821c49479ab48640292eacb3b7c4be"),
			},
		),
		gen(
			"zero length",
			[]string{cndInputMessage},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				input: []byte{},
			},
			&action{
				// The expected value can be checked by openssl (openssl only shows first 32 bytes).
				// echo -n "" | openssl dgst -shake256
				expect: hexMustDecode("46b9dd2b0ba88d13233b3feb743eeb243fcd52ea62b81b82b50c27646ed5762fd75dc4ddd8c0f200cb05019d67b592f6fc821c49479ab48640292eacb3b7c4be"),
			},
		),
		gen(
			"non-zero length 1",
			[]string{cndInputMessage},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				input: []byte("abc"),
			},
			&action{
				// The expected value can be checked by openssl (openssl only shows first 32 bytes).
				// echo -n "abc" | openssl dgst -shake256
				// Online tools like https://emn178.github.io/online-tools/shake_256.html
				// can be used for generating a value.
				expect: hexMustDecode("483366601360a8771c6863080cc4114d8db44530f8f1e1ee4f94ea37e78b5739d5a15bef186a5386c75744c0527e1faa9f8726e462a12a4feb06bd8801e751e4"),
			},
		),
		gen(
			"non-zero length 2",
			[]string{cndInputMessage},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				input: []byte("test"),
			},
			&action{
				// The expected value can be checked by openssl (openssl only shows first 32 bytes).
				// echo -n "test" | openssl dgst -shake256
				// Online tools like https://emn178.github.io/online-tools/shake_256.html
				// can be used for generating a value.
				expect: hexMustDecode("b54ff7255705a71ee2925e4a3e30e41aed489a579d5595e0df13e32e1e4dd202a7c7f68b31d6418d9845eb4d757adda6ab189e1bb340db818e5b3bc725d992fa"),
			},
		),
		gen(
			"not match",
			[]string{cndInputMessage},
			[]string{actCheckHash, actCheckNotMatch},
			&condition{
				input:    []byte("test"),
				notMatch: true,
			},
			&action{
				// The first character is wrong.
				expect: hexMustDecode("054ff7255705a71ee2925e4a3e30e41aed489a579d5595e0df13e32e1e4dd202a7c7f68b31d6418d9845eb4d757adda6ab189e1bb340db818e5b3bc725d992fa"),
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			out := hash.SHAKE256(tt.C().input)
			testutil.Diff(t, hash.SizeSHAKE256, len(out))
			testutil.Diff(t, hash.HashSize[k.HashAlg_SHAKE256], len(out))

			if tt.C().notMatch {
				if bytes.Equal(tt.A().expect, out) {
					t.Error("error unexpectedly matched", tt.A().expect, out)
				}
			} else {
				testutil.Diff(t, tt.A().expect, out)
			}
		})
	}
}

func TestMD5(t *testing.T) {
	type condition struct {
		input    []byte
		notMatch bool
	}

	type action struct {
		expect []byte
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	cndInputMessage := tb.Condition("non-nil message", "input non-zero or non-nil message")
	cndInputNil := tb.Condition("nil message", "input nil as message")
	actCheckHash := tb.Action("check hash", "check that the expected hash is returned")
	actCheckMatch := tb.Action("check match", "check that the expected hash matched")
	actCheckNotMatch := tb.Action("check not match", "check that the expected hash does not matched")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"nil",
			[]string{cndInputNil},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				input: nil,
			},
			&action{
				// The expected value is generated by openssl.
				// echo -n "" | openssl dgst -md5
				expect: hexMustDecode("d41d8cd98f00b204e9800998ecf8427e"),
			},
		),
		gen(
			"zero length",
			[]string{cndInputMessage},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				input: []byte{},
			},
			&action{
				// The expected value is generated by openssl.
				// echo -n "" | openssl dgst -md5
				expect: hexMustDecode("d41d8cd98f00b204e9800998ecf8427e"),
			},
		),
		gen(
			"non-zero length 1",
			[]string{cndInputMessage},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				input: []byte("abc"),
			},
			&action{
				// The expected value is generated by openssl.
				// echo -n "abc" | openssl dgst -md5
				expect: hexMustDecode("900150983cd24fb0d6963f7d28e17f72"),
			},
		),
		gen(
			"non-zero length 2",
			[]string{cndInputMessage},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				input: []byte("test"),
			},
			&action{
				// The expected value is generated by openssl.
				// echo -n "test" | openssl dgst -md5
				expect: hexMustDecode("098f6bcd4621d373cade4e832627b4f6"),
			},
		),
		gen(
			"not match",
			[]string{cndInputMessage},
			[]string{actCheckHash, actCheckNotMatch},
			&condition{
				input:    []byte("test"),
				notMatch: true,
			},
			&action{
				// The first character is wrong.
				expect: hexMustDecode("198f6bcd4621d373cade4e832627b4f6"),
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			out := hash.MD5(tt.C().input)
			testutil.Diff(t, hash.SizeMD5, len(out))
			testutil.Diff(t, hash.HashSize[k.HashAlg_MD5], len(out))

			if tt.C().notMatch {
				if bytes.Equal(tt.A().expect, out) {
					t.Error("error unexpectedly matched", tt.A().expect, out)
				}
			} else {
				testutil.Diff(t, tt.A().expect, out)
			}
		})
	}
}

func TestFNV1_32(t *testing.T) {
	type condition struct {
		input    []byte
		notMatch bool
	}

	type action struct {
		expect []byte
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	cndInputMessage := tb.Condition("non-nil message", "input non-zero or non-nil message")
	cndInputNil := tb.Condition("nil message", "input nil as message")
	actCheckHash := tb.Action("check hash", "check that the expected hash is returned")
	actCheckMatch := tb.Action("check match", "check that the expected hash matched")
	actCheckNotMatch := tb.Action("check not match", "check that the expected hash does not matched")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"nil",
			[]string{cndInputNil},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				input: nil,
			},
			&action{
				// Generated by https://fnvhash.github.io/fnv-calculator-online/
				expect: hexMustDecode("811c9dc5"),
			},
		),
		gen(
			"zero length",
			[]string{cndInputMessage},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				input: []byte{},
			},
			&action{
				// Generated by https://fnvhash.github.io/fnv-calculator-online/
				expect: hexMustDecode("811c9dc5"),
			},
		),
		gen(
			"non-zero length 1",
			[]string{cndInputMessage},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				input: []byte("abc"),
			},
			&action{
				// Generated by https://fnvhash.github.io/fnv-calculator-online/
				expect: hexMustDecode("439c2f4b"),
			},
		),
		gen(
			"non-zero length 2",
			[]string{cndInputMessage},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				input: []byte("test"),
			},
			&action{
				// Generated by https://fnvhash.github.io/fnv-calculator-online/
				expect: hexMustDecode("bc2c0be9"),
			},
		),
		gen(
			"not match",
			[]string{cndInputMessage},
			[]string{actCheckHash, actCheckNotMatch},
			&condition{
				input:    []byte("test"),
				notMatch: true,
			},
			&action{
				// The first character is wrong.
				expect: hexMustDecode("0c2c0be9"),
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			out := hash.FNV1_32(tt.C().input)
			testutil.Diff(t, hash.SizeFNV1_32, len(out))
			testutil.Diff(t, hash.HashSize[k.HashAlg_FNV1_32], len(out))

			if tt.C().notMatch {
				if bytes.Equal(tt.A().expect, out) {
					t.Error("error unexpectedly matched", tt.A().expect, out)
				}
			} else {
				testutil.Diff(t, tt.A().expect, out)
			}
		})
	}
}

func TestFNV1a_32(t *testing.T) {
	type condition struct {
		input    []byte
		notMatch bool
	}

	type action struct {
		expect []byte
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	cndInputMessage := tb.Condition("non-nil message", "input non-zero or non-nil message")
	cndInputNil := tb.Condition("nil message", "input nil as message")
	actCheckHash := tb.Action("check hash", "check that the expected hash is returned")
	actCheckMatch := tb.Action("check match", "check that the expected hash matched")
	actCheckNotMatch := tb.Action("check not match", "check that the expected hash does not matched")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"nil",
			[]string{cndInputNil},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				input: nil,
			},
			&action{
				// Generated by https://fnvhash.github.io/fnv-calculator-online/
				expect: hexMustDecode("811c9dc5"),
			},
		),
		gen(
			"zero length",
			[]string{cndInputMessage},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				input: []byte{},
			},
			&action{
				// Generated by https://fnvhash.github.io/fnv-calculator-online/
				expect: hexMustDecode("811c9dc5"),
			},
		),
		gen(
			"non-zero length 1",
			[]string{cndInputMessage},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				input: []byte("abc"),
			},
			&action{
				// Generated by https://fnvhash.github.io/fnv-calculator-online/
				expect: hexMustDecode("1a47e90b"),
			},
		),
		gen(
			"non-zero length 2",
			[]string{cndInputMessage},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				input: []byte("test"),
			},
			&action{
				// Generated by https://fnvhash.github.io/fnv-calculator-online/
				expect: hexMustDecode("afd071e5"),
			},
		),
		gen(
			"not match",
			[]string{cndInputMessage},
			[]string{actCheckHash, actCheckNotMatch},
			&condition{
				input:    []byte("test"),
				notMatch: true,
			},
			&action{
				// The first character is wrong.
				expect: hexMustDecode("0fd071e5"),
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			out := hash.FNV1a_32(tt.C().input)
			testutil.Diff(t, hash.SizeFNV1a_32, len(out))
			testutil.Diff(t, hash.HashSize[k.HashAlg_FNV1a_32], len(out))

			if tt.C().notMatch {
				if bytes.Equal(tt.A().expect, out) {
					t.Error("error unexpectedly matched", tt.A().expect, out)
				}
			} else {
				testutil.Diff(t, tt.A().expect, out)
			}
		})
	}
}

func TestFNV1_64(t *testing.T) {
	type condition struct {
		input    []byte
		notMatch bool
	}

	type action struct {
		expect []byte
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	cndInputMessage := tb.Condition("non-nil message", "input non-zero or non-nil message")
	cndInputNil := tb.Condition("nil message", "input nil as message")
	actCheckHash := tb.Action("check hash", "check that the expected hash is returned")
	actCheckMatch := tb.Action("check match", "check that the expected hash matched")
	actCheckNotMatch := tb.Action("check not match", "check that the expected hash does not matched")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"nil",
			[]string{cndInputNil},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				input: nil,
			},
			&action{
				// Generated by https://fnvhash.github.io/fnv-calculator-online/
				expect: hexMustDecode("cbf29ce484222325"),
			},
		),
		gen(
			"zero length",
			[]string{cndInputMessage},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				input: []byte{},
			},
			&action{
				// Generated by https://fnvhash.github.io/fnv-calculator-online/
				expect: hexMustDecode("cbf29ce484222325"),
			},
		),
		gen(
			"non-zero length 1",
			[]string{cndInputMessage},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				input: []byte("abc"),
			},
			&action{
				// Generated by https://fnvhash.github.io/fnv-calculator-online/
				expect: hexMustDecode("d8dcca186bafadcb"),
			},
		),
		gen(
			"non-zero length 2",
			[]string{cndInputMessage},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				input: []byte("test"),
			},
			&action{
				// Generated by https://fnvhash.github.io/fnv-calculator-online/
				expect: hexMustDecode("8c093f7e9fccbf69"),
			},
		),
		gen(
			"not match",
			[]string{cndInputMessage},
			[]string{actCheckHash, actCheckNotMatch},
			&condition{
				input:    []byte("test"),
				notMatch: true,
			},
			&action{
				// The first character is wrong.
				expect: hexMustDecode("0c093f7e9fccbf69"),
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			out := hash.FNV1_64(tt.C().input)
			testutil.Diff(t, hash.SizeFNV1_64, len(out))
			testutil.Diff(t, hash.HashSize[k.HashAlg_FNV1_64], len(out))

			if tt.C().notMatch {
				if bytes.Equal(tt.A().expect, out) {
					t.Error("error unexpectedly matched", tt.A().expect, out)
				}
			} else {
				testutil.Diff(t, tt.A().expect, out)
			}
		})
	}
}

func TestFNV1a_64(t *testing.T) {
	type condition struct {
		input    []byte
		notMatch bool
	}

	type action struct {
		expect []byte
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	cndInputMessage := tb.Condition("non-nil message", "input non-zero or non-nil message")
	cndInputNil := tb.Condition("nil message", "input nil as message")
	actCheckHash := tb.Action("check hash", "check that the expected hash is returned")
	actCheckMatch := tb.Action("check match", "check that the expected hash matched")
	actCheckNotMatch := tb.Action("check not match", "check that the expected hash does not matched")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"nil",
			[]string{cndInputNil},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				input: nil,
			},
			&action{
				// Generated by https://fnvhash.github.io/fnv-calculator-online/
				expect: hexMustDecode("cbf29ce484222325"),
			},
		),
		gen(
			"zero length",
			[]string{cndInputMessage},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				input: []byte{},
			},
			&action{
				// Generated by https://fnvhash.github.io/fnv-calculator-online/
				expect: hexMustDecode("cbf29ce484222325"),
			},
		),
		gen(
			"non-zero length 1",
			[]string{cndInputMessage},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				input: []byte("abc"),
			},
			&action{
				// Generated by https://fnvhash.github.io/fnv-calculator-online/
				expect: hexMustDecode("e71fa2190541574b"),
			},
		),
		gen(
			"non-zero length 2",
			[]string{cndInputMessage},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				input: []byte("test"),
			},
			&action{
				// Generated by https://fnvhash.github.io/fnv-calculator-online/
				expect: hexMustDecode("f9e6e6ef197c2b25"),
			},
		),
		gen(
			"not match",
			[]string{cndInputMessage},
			[]string{actCheckHash, actCheckNotMatch},
			&condition{
				input:    []byte("test"),
				notMatch: true,
			},
			&action{
				// The first character is wrong.
				expect: hexMustDecode("09e6e6ef197c2b25"),
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			out := hash.FNV1a_64(tt.C().input)
			testutil.Diff(t, hash.SizeFNV1a_64, len(out))
			testutil.Diff(t, hash.HashSize[k.HashAlg_FNV1a_64], len(out))

			if tt.C().notMatch {
				if bytes.Equal(tt.A().expect, out) {
					t.Error("error unexpectedly matched", tt.A().expect, out)
				}
			} else {
				testutil.Diff(t, tt.A().expect, out)
			}
		})
	}
}

func TestFNV1_128(t *testing.T) {
	type condition struct {
		input    []byte
		notMatch bool
	}

	type action struct {
		expect []byte
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	cndInputMessage := tb.Condition("non-nil message", "input non-zero or non-nil message")
	cndInputNil := tb.Condition("nil message", "input nil as message")
	actCheckHash := tb.Action("check hash", "check that the expected hash is returned")
	actCheckMatch := tb.Action("check match", "check that the expected hash matched")
	actCheckNotMatch := tb.Action("check not match", "check that the expected hash does not matched")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"nil",
			[]string{cndInputNil},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				input: nil,
			},
			&action{
				// Generated by https://fnvhash.github.io/fnv-calculator-online/
				expect: hexMustDecode("6c62272e07bb014262b821756295c58d"),
			},
		),
		gen(
			"zero length",
			[]string{cndInputMessage},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				input: []byte{},
			},
			&action{
				// Generated by https://fnvhash.github.io/fnv-calculator-online/
				expect: hexMustDecode("6c62272e07bb014262b821756295c58d"),
			},
		),
		gen(
			"non-zero length 1",
			[]string{cndInputMessage},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				input: []byte("abc"),
			},
			&action{
				// Generated by https://fnvhash.github.io/fnv-calculator-online/
				expect: hexMustDecode("a68bb2a4348b5822836dbc78c6aee73b"),
			},
		),
		gen(
			"non-zero length 2",
			[]string{cndInputMessage},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				input: []byte("test"),
			},
			&action{
				// Generated by https://fnvhash.github.io/fnv-calculator-online/
				expect: hexMustDecode("66ab2a8b6f757277b806e89c56faf339"),
			},
		),
		gen(
			"not match",
			[]string{cndInputMessage},
			[]string{actCheckHash, actCheckNotMatch},
			&condition{
				input:    []byte("test"),
				notMatch: true,
			},
			&action{
				// The first character is wrong.
				expect: hexMustDecode("06ab2a8b6f757277b806e89c56faf339"),
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			out := hash.FNV1_128(tt.C().input)
			testutil.Diff(t, hash.SizeFNV1_128, len(out))
			testutil.Diff(t, hash.HashSize[k.HashAlg_FNV1_128], len(out))

			if tt.C().notMatch {
				if bytes.Equal(tt.A().expect, out) {
					t.Error("error unexpectedly matched", tt.A().expect, out)
				}
			} else {
				testutil.Diff(t, tt.A().expect, out)
			}
		})
	}
}

func TestFNV1a_128(t *testing.T) {
	type condition struct {
		input    []byte
		notMatch bool
	}

	type action struct {
		expect []byte
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	cndInputMessage := tb.Condition("non-nil message", "input non-zero or non-nil message")
	cndInputNil := tb.Condition("nil message", "input nil as message")
	actCheckHash := tb.Action("check hash", "check that the expected hash is returned")
	actCheckMatch := tb.Action("check match", "check that the expected hash matched")
	actCheckNotMatch := tb.Action("check not match", "check that the expected hash does not matched")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"nil",
			[]string{cndInputNil},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				input: nil,
			},
			&action{
				// Generated by https://fnvhash.github.io/fnv-calculator-online/
				expect: hexMustDecode("6c62272e07bb014262b821756295c58d"),
			},
		),
		gen(
			"zero length",
			[]string{cndInputMessage},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				input: []byte{},
			},
			&action{
				// Generated by https://fnvhash.github.io/fnv-calculator-online/
				expect: hexMustDecode("6c62272e07bb014262b821756295c58d"),
			},
		),
		gen(
			"non-zero length 1",
			[]string{cndInputMessage},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				input: []byte("abc"),
			},
			&action{
				// Generated by https://fnvhash.github.io/fnv-calculator-online/
				expect: hexMustDecode("a68d622cec8b5822836dbc7977af7f3b"),
			},
		),
		gen(
			"non-zero length 2",
			[]string{cndInputMessage},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				input: []byte("test"),
			},
			&action{
				// Generated by https://fnvhash.github.io/fnv-calculator-online/
				expect: hexMustDecode("69d061a9c5757277b806e99413dd99a5"),
			},
		),
		gen(
			"not match",
			[]string{cndInputMessage},
			[]string{actCheckHash, actCheckNotMatch},
			&condition{
				input:    []byte("test"),
				notMatch: true,
			},
			&action{
				// The first character is wrong.
				expect: hexMustDecode("09d061a9c5757277b806e99413dd99a5"),
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			out := hash.FNV1a_128(tt.C().input)
			testutil.Diff(t, hash.SizeFNV1a_128, len(out))
			testutil.Diff(t, hash.HashSize[k.HashAlg_FNV1a_128], len(out))

			if tt.C().notMatch {
				if bytes.Equal(tt.A().expect, out) {
					t.Error("error unexpectedly matched", tt.A().expect, out)
				}
			} else {
				testutil.Diff(t, tt.A().expect, out)
			}
		})
	}
}

func TestCRC32(t *testing.T) {
	type condition struct {
		input    []byte
		notMatch bool
	}

	type action struct {
		expect []byte
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	cndInputMessage := tb.Condition("non-nil message", "input non-zero or non-nil message")
	cndInputNil := tb.Condition("nil message", "input nil as message")
	actCheckHash := tb.Action("check hash", "check that the expected hash is returned")
	actCheckMatch := tb.Action("check match", "check that the expected hash matched")
	actCheckNotMatch := tb.Action("check not match", "check that the expected hash does not matched")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"nil",
			[]string{cndInputNil},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				input: nil,
			},
			&action{
				expect: hexMustDecode("00000000"),
			},
		),
		gen(
			"zero length",
			[]string{cndInputMessage},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				input: []byte{},
			},
			&action{
				expect: hexMustDecode("00000000"),
			},
		),
		gen(
			"non-zero length 1",
			[]string{cndInputMessage},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				input: []byte("abc"),
			},
			&action{
				expect: hexMustDecode("352441c2"),
			},
		),
		gen(
			"non-zero length 2",
			[]string{cndInputMessage},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				input: []byte("test"),
			},
			&action{
				expect: hexMustDecode("d87f7e0c"),
			},
		),
		gen(
			"not match",
			[]string{cndInputMessage},
			[]string{actCheckHash, actCheckNotMatch},
			&condition{
				input:    []byte("test"),
				notMatch: true,
			},
			&action{
				// The first character is wrong.
				expect: hexMustDecode("087f7e0c"),
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			out := hash.CRC32(tt.C().input)
			testutil.Diff(t, hash.SizeCRC32, len(out))
			testutil.Diff(t, hash.HashSize[k.HashAlg_CRC32], len(out))

			if tt.C().notMatch {
				if bytes.Equal(tt.A().expect, out) {
					t.Error("error unexpectedly matched", tt.A().expect, out)
				}
			} else {
				testutil.Diff(t, tt.A().expect, out)
			}
		})
	}
}

func TestCRC64ISO(t *testing.T) {
	type condition struct {
		input    []byte
		notMatch bool
	}

	type action struct {
		expect []byte
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	cndInputMessage := tb.Condition("non-nil message", "input non-zero or non-nil message")
	cndInputNil := tb.Condition("nil message", "input nil as message")
	actCheckHash := tb.Action("check hash", "check that the expected hash is returned")
	actCheckMatch := tb.Action("check match", "check that the expected hash matched")
	actCheckNotMatch := tb.Action("check not match", "check that the expected hash does not matched")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"nil",
			[]string{cndInputNil},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				input: nil,
			},
			&action{
				expect: hexMustDecode("0000000000000000"),
			},
		),
		gen(
			"zero length",
			[]string{cndInputMessage},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				input: []byte{},
			},
			&action{
				expect: hexMustDecode("0000000000000000"),
			},
		),
		gen(
			"non-zero length 1",
			[]string{cndInputMessage},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				input: []byte("abc"),
			},
			&action{
				expect: hexMustDecode("3776c42000000000"),
			},
		),
		gen(
			"non-zero length 2",
			[]string{cndInputMessage},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				input: []byte("test"),
			},
			&action{
				expect: hexMustDecode("287c72c850000000"),
			},
		),
		gen(
			"not match",
			[]string{cndInputMessage},
			[]string{actCheckHash, actCheckNotMatch},
			&condition{
				input:    []byte("test"),
				notMatch: true,
			},
			&action{
				// The first character is wrong.
				expect: hexMustDecode("087c72c850000000"),
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			out := hash.CRC64ISO(tt.C().input)
			testutil.Diff(t, hash.SizeCRC64ISO, len(out))
			testutil.Diff(t, hash.HashSize[k.HashAlg_CRC64ISO], len(out))

			if tt.C().notMatch {
				if bytes.Equal(tt.A().expect, out) {
					t.Error("error unexpectedly matched", tt.A().expect, out)
				}
			} else {
				testutil.Diff(t, tt.A().expect, out)
			}
		})
	}
}

func TestCRC64ECMA(t *testing.T) {
	type condition struct {
		input    []byte
		notMatch bool
	}

	type action struct {
		expect []byte
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	cndInputMessage := tb.Condition("non-nil message", "input non-zero or non-nil message")
	cndInputNil := tb.Condition("nil message", "input nil as message")
	actCheckHash := tb.Action("check hash", "check that the expected hash is returned")
	actCheckMatch := tb.Action("check match", "check that the expected hash matched")
	actCheckNotMatch := tb.Action("check not match", "check that the expected hash does not matched")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"nil",
			[]string{cndInputNil},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				input: nil,
			},
			&action{
				expect: hexMustDecode("0000000000000000"),
			},
		),
		gen(
			"zero length",
			[]string{cndInputMessage},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				input: []byte{},
			},
			&action{
				expect: hexMustDecode("0000000000000000"),
			},
		),
		gen(
			"non-zero length 1",
			[]string{cndInputMessage},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				input: []byte("abc"),
			},
			&action{
				expect: hexMustDecode("2cd8094a1a277627"),
			},
		),
		gen(
			"non-zero length 2",
			[]string{cndInputMessage},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				input: []byte("test"),
			},
			&action{
				expect: hexMustDecode("fa15fda7c10c75a5"),
			},
		),
		gen(
			"not match",
			[]string{cndInputMessage},
			[]string{actCheckHash, actCheckNotMatch},
			&condition{
				input:    []byte("test"),
				notMatch: true,
			},
			&action{
				// The first character is wrong.
				expect: hexMustDecode("0a15fda7c10c75a5"),
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			out := hash.CRC64ECMA(tt.C().input)
			testutil.Diff(t, hash.SizeCRC64ECMA, len(out))
			testutil.Diff(t, hash.HashSize[k.HashAlg_CRC64ECMA], len(out))

			if tt.C().notMatch {
				if bytes.Equal(tt.A().expect, out) {
					t.Error("error unexpectedly matched", tt.A().expect, out)
				}
			} else {
				testutil.Diff(t, tt.A().expect, out)
			}
		})
	}
}

func TestBLAKE2s_256(t *testing.T) {
	type condition struct {
		input    []byte
		notMatch bool
	}

	type action struct {
		expect []byte
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	cndInputMessage := tb.Condition("non-nil message", "input non-zero or non-nil message")
	cndInputNil := tb.Condition("nil message", "input nil as message")
	actCheckHash := tb.Action("check hash", "check that the expected hash is returned")
	actCheckMatch := tb.Action("check match", "check that the expected hash matched")
	actCheckNotMatch := tb.Action("check not match", "check that the expected hash does not matched")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"nil",
			[]string{cndInputNil},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				input: nil,
			},
			&action{
				// The expected value is generated by openssl.
				// echo -n "" | openssl dgst -blake2s256
				expect: hexMustDecode("69217a3079908094e11121d042354a7c1f55b6482ca1a51e1b250dfd1ed0eef9"),
			},
		),
		gen(
			"zero length",
			[]string{cndInputMessage},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				input: []byte{},
			},
			&action{
				// The expected value is generated by openssl.
				// echo -n "" | openssl dgst -blake2s256
				expect: hexMustDecode("69217a3079908094e11121d042354a7c1f55b6482ca1a51e1b250dfd1ed0eef9"),
			},
		),
		gen(
			"non-zero length",
			[]string{cndInputMessage},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				input: []byte("abc"),
			},
			&action{
				// See https://www.rfc-editor.org/rfc/rfc7693#appendix-B
				// The expected value is generated by openssl.
				// echo -n "abc" | openssl dgst -blake2s256
				expect: hexMustDecode("508c5e8c327c14e2e1a72ba34eeb452f37458b209ed63a294d999b4c86675982"),
			},
		),
		gen(
			"not match",
			[]string{cndInputMessage},
			[]string{actCheckHash, actCheckNotMatch},
			&condition{
				input:    []byte("abc"),
				notMatch: true,
			},
			&action{
				// The first character is wrong.
				expect: hexMustDecode("008c5e8c327c14e2e1a72ba34eeb452f37458b209ed63a294d999b4c86675982"),
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			out := hash.BLAKE2s_256(tt.C().input)
			testutil.Diff(t, hash.SizeBLAKE2s_256, len(out))
			testutil.Diff(t, hash.HashSize[k.HashAlg_BLAKE2s_256], len(out))

			if tt.C().notMatch {
				if bytes.Equal(tt.A().expect, out) {
					t.Error("error unexpectedly matched", tt.A().expect, out)
				}
			} else {
				testutil.Diff(t, tt.A().expect, out)
			}
		})
	}
}

func TestBLAKE2b_256(t *testing.T) {
	type condition struct {
		input    []byte
		notMatch bool
	}

	type action struct {
		expect []byte
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	cndInputMessage := tb.Condition("non-nil message", "input non-zero or non-nil message")
	cndInputNil := tb.Condition("nil message", "input nil as message")
	actCheckHash := tb.Action("check hash", "check that the expected hash is returned")
	actCheckMatch := tb.Action("check match", "check that the expected hash matched")
	actCheckNotMatch := tb.Action("check not match", "check that the expected hash does not matched")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"nil",
			[]string{cndInputNil},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				input: nil,
			},
			&action{
				// Generated by https://asecuritysite.com/hash/blake?m=
				expect: hexMustDecode("0e5751c026e543b2e8ab2eb06099daa1d1e5df47778f7787faab45cdf12fe3a8"),
			},
		),
		gen(
			"zero length",
			[]string{cndInputMessage},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				input: []byte{},
			},
			&action{
				// Generated by https://asecuritysite.com/hash/blake?m=
				expect: hexMustDecode("0e5751c026e543b2e8ab2eb06099daa1d1e5df47778f7787faab45cdf12fe3a8"),
			},
		),
		gen(
			"non-zero length",
			[]string{cndInputMessage},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				input: []byte("abc"),
			},
			&action{
				// Generated by https://8gwifi.org/MessageDigest.jsp
				expect: hexMustDecode("bddd813c634239723171ef3fee98579b94964e3bb1cb3e427262c8c068d52319"),
			},
		),
		gen(
			"not match",
			[]string{cndInputMessage},
			[]string{actCheckHash, actCheckNotMatch},
			&condition{
				input:    []byte("abc"),
				notMatch: true,
			},
			&action{
				// The first character is wrong.
				expect: hexMustDecode("0ddd813c634239723171ef3fee98579b94964e3bb1cb3e427262c8c068d52319"),
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			out := hash.BLAKE2b_256(tt.C().input)
			testutil.Diff(t, hash.SizeBLAKE2b_256, len(out))
			testutil.Diff(t, hash.HashSize[k.HashAlg_BLAKE2b_256], len(out))

			if tt.C().notMatch {
				if bytes.Equal(tt.A().expect, out) {
					t.Error("error unexpectedly matched", tt.A().expect, out)
				}
			} else {
				testutil.Diff(t, tt.A().expect, out)
			}
		})
	}
}

func TestBLAKE2b_384(t *testing.T) {
	type condition struct {
		input    []byte
		notMatch bool
	}

	type action struct {
		expect []byte
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	cndInputMessage := tb.Condition("non-nil message", "input non-zero or non-nil message")
	cndInputNil := tb.Condition("nil message", "input nil as message")
	actCheckHash := tb.Action("check hash", "check that the expected hash is returned")
	actCheckMatch := tb.Action("check match", "check that the expected hash matched")
	actCheckNotMatch := tb.Action("check not match", "check that the expected hash does not matched")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"nil",
			[]string{cndInputNil},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				input: nil,
			},
			&action{
				// See https://en.wikipedia.org/wiki/BLAKE_(hash_function)
				expect: hexMustDecode("b32811423377f52d7862286ee1a72ee540524380fda1724a6f25d7978c6fd3244a6caf0498812673c5e05ef583825100"),
			},
		),
		gen(
			"zero length",
			[]string{cndInputMessage},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				input: []byte{},
			},
			&action{
				// See https://en.wikipedia.org/wiki/BLAKE_(hash_function)
				expect: hexMustDecode("b32811423377f52d7862286ee1a72ee540524380fda1724a6f25d7978c6fd3244a6caf0498812673c5e05ef583825100"),
			},
		),
		gen(
			"non-zero length",
			[]string{cndInputMessage},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				input: []byte("abc"),
			},
			&action{
				// Generated by https://8gwifi.org/MessageDigest.jsp
				expect: hexMustDecode("6f56a82c8e7ef526dfe182eb5212f7db9df1317e57815dbda46083fc30f54ee6c66ba83be64b302d7cba6ce15bb556f4"),
			},
		),
		gen(
			"not match",
			[]string{cndInputMessage},
			[]string{actCheckHash, actCheckNotMatch},
			&condition{
				input:    []byte("abc"),
				notMatch: true,
			},
			&action{
				// The first character is wrong.
				expect: hexMustDecode("0f56a82c8e7ef526dfe182eb5212f7db9df1317e57815dbda46083fc30f54ee6c66ba83be64b302d7cba6ce15bb556f4"),
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			out := hash.BLAKE2b_384(tt.C().input)
			testutil.Diff(t, hash.SizeBLAKE2b_384, len(out))
			testutil.Diff(t, hash.HashSize[k.HashAlg_BLAKE2b_384], len(out))

			if tt.C().notMatch {
				if bytes.Equal(tt.A().expect, out) {
					t.Error("error unexpectedly matched", tt.A().expect, out)
				}
			} else {
				testutil.Diff(t, tt.A().expect, out)
			}
		})
	}
}

func TestBLAKE2b_512(t *testing.T) {
	type condition struct {
		input    []byte
		notMatch bool
	}

	type action struct {
		expect []byte
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	cndInputMessage := tb.Condition("non-nil message", "input non-zero or non-nil message")
	cndInputNil := tb.Condition("nil message", "input nil as message")
	actCheckHash := tb.Action("check hash", "check that the expected hash is returned")
	actCheckMatch := tb.Action("check match", "check that the expected hash matched")
	actCheckNotMatch := tb.Action("check not match", "check that the expected hash does not matched")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"nil",
			[]string{cndInputNil},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				input: nil,
			},
			&action{
				// The expected value is generated by openssl.
				// echo -n "" | openssl dgst -blake2b512
				expect: hexMustDecode("786a02f742015903c6c6fd852552d272912f4740e15847618a86e217f71f5419d25e1031afee585313896444934eb04b903a685b1448b755d56f701afe9be2ce"),
			},
		),
		gen(
			"zero length",
			[]string{cndInputMessage},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				input: []byte{},
			},
			&action{
				// The expected value is generated by openssl.
				// echo -n "" | openssl dgst -blake2b512
				expect: hexMustDecode("786a02f742015903c6c6fd852552d272912f4740e15847618a86e217f71f5419d25e1031afee585313896444934eb04b903a685b1448b755d56f701afe9be2ce"),
			},
		),
		gen(
			"non-zero length",
			[]string{cndInputMessage},
			[]string{actCheckHash, actCheckMatch},
			&condition{
				input: []byte("abc"),
			},
			&action{
				// See https://www.rfc-editor.org/rfc/rfc7693#appendix-A
				// The expected value is generated by openssl.
				// echo -n "abc" | openssl dgst -blake2b512
				expect: hexMustDecode("ba80a53f981c4d0d6a2797b69f12f6e94c212f14685ac4b74b12bb6fdbffa2d17d87c5392aab792dc252d5de4533cc9518d38aa8dbf1925ab92386edd4009923"),
			},
		),
		gen(
			"not match",
			[]string{cndInputMessage},
			[]string{actCheckHash, actCheckNotMatch},
			&condition{
				input:    []byte("abc"),
				notMatch: true,
			},
			&action{
				// The first character is wrong.
				expect: hexMustDecode("0a80a53f981c4d0d6a2797b69f12f6e94c212f14685ac4b74b12bb6fdbffa2d17d87c5392aab792dc252d5de4533cc9518d38aa8dbf1925ab92386edd4009923"),
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			out := hash.BLAKE2b_512(tt.C().input)
			testutil.Diff(t, hash.SizeBLAKE2b_512, len(out))
			testutil.Diff(t, hash.HashSize[k.HashAlg_BLAKE2b_512], len(out))

			if tt.C().notMatch {
				if bytes.Equal(tt.A().expect, out) {
					t.Error("error unexpectedly matched", tt.A().expect, out)
				}
			} else {
				testutil.Diff(t, tt.A().expect, out)
			}
		})
	}
}
