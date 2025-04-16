// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package hash

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"hash/crc32"
	"hash/crc64"
	"hash/fnv"

	k "github.com/aileron-gateway/aileron-gateway/apis/kernel"
	"golang.org/x/crypto/blake2b"
	"golang.org/x/crypto/blake2s"
	"golang.org/x/crypto/sha3"
)

var (
	crc64ISOTable  = crc64.MakeTable(crc64.ISO)
	crc64ECMATable = crc64.MakeTable(crc64.ECMA)
)

// HashFunc is the function that returns hash value of the given data.
type HashFunc func(data []byte) (hash []byte)

// Algorithm is the hash algorithm.
type Algorithm int

const (
	AlgUnknown Algorithm = iota
	AlgSHA1
	AlgSHA224
	AlgSHA256
	AlgSHA384
	AlgSHA512
	AlgSHA512_224
	AlgSHA512_256
	AlgSHA3_224
	AlgSHA3_256
	AlgSHA3_384
	AlgSHA3_512
	AlgSHAKE128
	AlgSHAKE256
	AlgMD5
	AlgFNV1_32
	AlgFNV1a_32 //lint:ignore ST1003 should not use underscores in Go names; const AlgFNV1a_32 should be AlgFNV1a32
	AlgFNV1_64
	AlgFNV1a_64 //lint:ignore ST1003 should not use underscores in Go names; const AlgFNV1a_64 should be AlgFNV1a64
	AlgFNV1_128
	AlgFNV1a_128 //lint:ignore ST1003 should not use underscores in Go names; const AlgFNV1a_128 should be AlgFNV1a128
	AlgCRC32
	AlgCRC64ISO
	AlgCRC64ECMA
	AlgBLAKE2s_256 //lint:ignore ST1003 should not use underscores in Go names; const AlgBLAKE2s_256 should be AlgBLAKE2s256
	AlgBLAKE2b_256 //lint:ignore ST1003 should not use underscores in Go names; const AlgBLAKE2b_256 should be AlgBLAKE2b256
	AlgBLAKE2b_384 //lint:ignore ST1003 should not use underscores in Go names; const AlgBLAKE2b_384 should be AlgBLAKE2b384
	AlgBLAKE2b_512 //lint:ignore ST1003 should not use underscores in Go names; const AlgBLAKE2b_512 should be AlgBLAKE2b512
)

const (
	SizeSHA1        = 20
	SizeSHA224      = 28
	SizeSHA256      = 32
	SizeSHA384      = 48
	SizeSHA512      = 64
	SizeSHA512_224  = 28
	SizeSHA512_256  = 32
	SizeSHA3_224    = 28
	SizeSHA3_256    = 32
	SizeSHA3_384    = 48
	SizeSHA3_512    = 64
	SizeSHAKE128    = 32
	SizeSHAKE256    = 64
	SizeMD5         = 16
	SizeFNV1_32     = 4
	SizeFNV1a_32    = 4 //lint:ignore ST1003 should not use underscores in Go names; const SizeFNV1a_32 should be SizeFNV1a32
	SizeFNV1_64     = 8
	SizeFNV1a_64    = 8 //lint:ignore ST1003 should not use underscores in Go names; const SizeFNV1a_64 should be SizeFNV1a64
	SizeFNV1_128    = 16
	SizeFNV1a_128   = 16 //lint:ignore ST1003 should not use underscores in Go names; const SizeFNV1a_128 should be SizeFNV1a128
	SizeCRC32       = 4
	SizeCRC64ISO    = 8
	SizeCRC64ECMA   = 8
	SizeBLAKE2s_256 = 32 //lint:ignore ST1003 should not use underscores in Go names; const SizeBLAKE2s_256 should be SizeBLAKE2s256
	SizeBLAKE2b_256 = 32 //lint:ignore ST1003 should not use underscores in Go names; const SizeBLAKE2b_256 should be SizeBLAKE2b256
	SizeBLAKE2b_384 = 48 //lint:ignore ST1003 should not use underscores in Go names; const SizeBLAKE2b_384 should be SizeBLAKE2b384
	SizeBLAKE2b_512 = 64 //lint:ignore ST1003 should not use underscores in Go names; const SizeBLAKE2b_512 should be SizeBLAKE2b512
)

var HashSize = map[k.HashAlg]int{
	k.HashAlg_SHA1:        SizeSHA1,
	k.HashAlg_SHA224:      SizeSHA224,
	k.HashAlg_SHA256:      SizeSHA256,
	k.HashAlg_SHA384:      SizeSHA384,
	k.HashAlg_SHA512:      SizeSHA512,
	k.HashAlg_SHA512_224:  SizeSHA512_224,
	k.HashAlg_SHA512_256:  SizeSHA512_256,
	k.HashAlg_SHA3_224:    SizeSHA3_224,
	k.HashAlg_SHA3_256:    SizeSHA3_256,
	k.HashAlg_SHA3_384:    SizeSHA3_384,
	k.HashAlg_SHA3_512:    SizeSHA3_512,
	k.HashAlg_SHAKE128:    SizeSHAKE128,
	k.HashAlg_SHAKE256:    SizeSHAKE256,
	k.HashAlg_MD5:         SizeMD5,
	k.HashAlg_FNV1_32:     SizeFNV1_32,
	k.HashAlg_FNV1a_32:    SizeFNV1a_32,
	k.HashAlg_FNV1_64:     SizeFNV1_64,
	k.HashAlg_FNV1a_64:    SizeFNV1a_64,
	k.HashAlg_FNV1_128:    SizeFNV1_128,
	k.HashAlg_FNV1a_128:   SizeFNV1a_128,
	k.HashAlg_CRC32:       SizeCRC32,
	k.HashAlg_CRC64ISO:    SizeCRC64ISO,
	k.HashAlg_CRC64ECMA:   SizeCRC64ECMA,
	k.HashAlg_BLAKE2s_256: SizeBLAKE2s_256,
	k.HashAlg_BLAKE2b_256: SizeBLAKE2b_256,
	k.HashAlg_BLAKE2b_384: SizeBLAKE2b_384,
	k.HashAlg_BLAKE2b_512: SizeBLAKE2b_512,
}

// FromAlgorithm returns hash function
// corresponding to the given hash algorithm.
// This function returns nil when hash function not found.
func FromAlgorithm(a Algorithm) HashFunc {
	algToFunc := map[Algorithm]HashFunc{
		AlgSHA1:        SHA1,
		AlgSHA224:      SHA224,
		AlgSHA256:      SHA256,
		AlgSHA384:      SHA384,
		AlgSHA512:      SHA512,
		AlgSHA512_224:  SHA512_224,
		AlgSHA512_256:  SHA512_256,
		AlgSHA3_224:    SHA3_224,
		AlgSHA3_256:    SHA3_256,
		AlgSHA3_384:    SHA3_384,
		AlgSHA3_512:    SHA3_512,
		AlgSHAKE128:    SHAKE128,
		AlgSHAKE256:    SHAKE256,
		AlgMD5:         MD5,
		AlgFNV1_32:     FNV1_32,
		AlgFNV1a_32:    FNV1a_32,
		AlgFNV1_64:     FNV1_64,
		AlgFNV1a_64:    FNV1a_64,
		AlgFNV1_128:    FNV1_128,
		AlgFNV1a_128:   FNV1a_128,
		AlgCRC32:       CRC32,
		AlgCRC64ISO:    CRC64ISO,
		AlgCRC64ECMA:   CRC64ECMA,
		AlgBLAKE2s_256: BLAKE2s_256,
		AlgBLAKE2b_256: BLAKE2b_256,
		AlgBLAKE2b_384: BLAKE2b_384,
		AlgBLAKE2b_512: BLAKE2b_512,
	}
	return algToFunc[a]
}

// FromHashAlg returns hash function
// corresponding to the given hash algorithm.
// This function returns nil when hash function not found.
func FromHashAlg(t k.HashAlg) HashFunc {
	typeToFunc := map[k.HashAlg]HashFunc{
		k.HashAlg_SHA1:        SHA1,
		k.HashAlg_SHA224:      SHA224,
		k.HashAlg_SHA256:      SHA256,
		k.HashAlg_SHA384:      SHA384,
		k.HashAlg_SHA512:      SHA512,
		k.HashAlg_SHA512_224:  SHA512_224,
		k.HashAlg_SHA512_256:  SHA512_256,
		k.HashAlg_SHA3_224:    SHA3_224,
		k.HashAlg_SHA3_256:    SHA3_256,
		k.HashAlg_SHA3_384:    SHA3_384,
		k.HashAlg_SHA3_512:    SHA3_512,
		k.HashAlg_SHAKE128:    SHAKE128,
		k.HashAlg_SHAKE256:    SHAKE256,
		k.HashAlg_MD5:         MD5,
		k.HashAlg_FNV1_32:     FNV1_32,
		k.HashAlg_FNV1a_32:    FNV1a_32,
		k.HashAlg_FNV1_64:     FNV1_64,
		k.HashAlg_FNV1a_64:    FNV1a_64,
		k.HashAlg_FNV1_128:    FNV1_128,
		k.HashAlg_FNV1a_128:   FNV1a_128,
		k.HashAlg_CRC32:       CRC32,
		k.HashAlg_CRC64ISO:    CRC64ISO,
		k.HashAlg_CRC64ECMA:   CRC64ECMA,
		k.HashAlg_BLAKE2s_256: BLAKE2s_256,
		k.HashAlg_BLAKE2b_256: BLAKE2b_256,
		k.HashAlg_BLAKE2b_384: BLAKE2b_384,
		k.HashAlg_BLAKE2b_512: BLAKE2b_512,
	}
	return typeToFunc[t]
}

// SHA1 returns SHA1 checksum of the given bytes.
// [20]byte slice is returned.
// Technically, sha1.Sum(b) is used.
//   - https://pkg.go.dev/crypto/sha1
func SHA1(b []byte) []byte {
	x := sha1.Sum(b)
	return x[:]
}

// SHA224 returns SHA224 checksum of the given data.
// [28]byte slice is returned.
// Technically, sha256.Sum224(b) is used.
//   - https://pkg.go.dev/crypto/sha256
func SHA224(b []byte) []byte {
	x := sha256.Sum224(b)
	return x[:]
}

// SHA256 returns SHA256 checksum of the given data.
// [32]byte slice is returned.
// Technically, sha256.Sum256(b) is used.
//   - https://pkg.go.dev/crypto/sha256
func SHA256(b []byte) []byte {
	x := sha256.Sum256(b)
	return x[:]
}

// SHA384 returns SHA384 checksum of the given data.
// [48]byte slice is returned.
// Technically, sha512.Sum384(b) is used.
//   - https://pkg.go.dev/crypto/sha512
func SHA384(b []byte) []byte {
	x := sha512.Sum384(b)
	return x[:]
}

// SHA512 returns SHA512 checksum of the given data.
// [64]byte slice is returned.
// Technically, sha512.Sum512(b) is used.
//   - https://pkg.go.dev/crypto/sha512
func SHA512(b []byte) []byte {
	x := sha512.Sum512(b)
	return x[:]
}

// SHA512_224 returns SHA512/224 checksum of the given data.
// [28]byte slice is returned.
// Technically, sha512.Sum512_224(b) is used.
//   - https://pkg.go.dev/crypto/sha512
//
//lint:ignore ST1003 should not use ALL_CAPS in Go names; use CamelCase instead
func SHA512_224(b []byte) []byte {
	x := sha512.Sum512_224(b)
	return x[:]
}

// SHA512_256 returns SHA512/256 checksum of the given data.
// [32]byte slice is returned.
// Technically, sha512.Sum512_256(b) is used.
//   - https://pkg.go.dev/crypto/sha512
//
//lint:ignore ST1003 should not use ALL_CAPS in Go names; use CamelCase instead
func SHA512_256(b []byte) []byte {
	x := sha512.Sum512_256(b)
	return x[:]
}

// SHA3_224 returns SHA3/224 checksum of the given data.
// [28]byte slice is returned.
// Technically, sha3.Sum224(b) is used.
//   - https://pkg.go.dev/golang.org/x/crypto/sha3
//
//lint:ignore ST1003 should not use ALL_CAPS in Go names; use CamelCase instead
func SHA3_224(b []byte) []byte {
	x := sha3.Sum224(b)
	return x[:]
}

// SHA3_256 returns SHA3/256 checksum of the given data.
// [32]byte slice is returned.
// Technically, sha3.Sum256(b) is used.
//   - https://pkg.go.dev/golang.org/x/crypto/sha3
//
//lint:ignore ST1003 should not use ALL_CAPS in Go names; use CamelCase instead
func SHA3_256(b []byte) []byte {
	x := sha3.Sum256(b)
	return x[:]
}

// SHA3_384 returns SHA3/384 checksum of the given data.
// [48]byte slice is returned.
// Technically, sha3.Sum384(b) is used.
//   - https://pkg.go.dev/golang.org/x/crypto/sha3
//
//lint:ignore ST1003 should not use ALL_CAPS in Go names; use CamelCase instead
func SHA3_384(b []byte) []byte {
	x := sha3.Sum384(b)
	return x[:]
}

// SHA3_512 returns SHA3/512 checksum of the given data.
// [64]byte slice is returned.
// Technically, sha3.Sum512(b) is used.
//   - https://pkg.go.dev/golang.org/x/crypto/sha3
//
//lint:ignore ST1003 should not use ALL_CAPS in Go names; use CamelCase instead
func SHA3_512(b []byte) []byte {
	x := sha3.Sum512(b)
	return x[:]
}

// SHAKE128 returns SHAKE128 checksum of the given data.
// [32]byte slice is returned.
// Technically, sha3.NewShake128() is used.
//   - https://pkg.go.dev/golang.org/x/crypto/sha3
func SHAKE128(b []byte) []byte {
	h := sha3.NewShake128()
	h.Write(b)
	return h.Sum(make([]byte, 0, SizeSHAKE128))
}

// SHAKE256 returns SHAKE256 checksum of the given data.
// [32]byte slice is returned.
// Technically, sha3.NewShake256() is used.
//   - https://pkg.go.dev/golang.org/x/crypto/sha3
func SHAKE256(b []byte) []byte {
	h := sha3.NewShake256()
	h.Write(b)
	return h.Sum(make([]byte, 0, SizeSHAKE256))
}

// MD5 returns MD5 checksum of the given data.
// [16]byte slice is returned.
// Technically, md5.Sum(b) is used.
//   - https://pkg.go.dev/golang.org/x/crypto/sha3
func MD5(b []byte) []byte {
	x := md5.Sum(b)
	return x[:]
}

// FNV1_32 returns FNV1/32 checksum of the given data.
// [4]byte slice is returned.
// Technically, fnv.New32() is used.
//   - https://pkg.go.dev/hash/fnv
//
//lint:ignore ST1003 should not use ALL_CAPS in Go names; use CamelCase instead
func FNV1_32(b []byte) []byte {
	h := fnv.New32()
	h.Write(b)
	return h.Sum(make([]byte, 0, SizeFNV1_32))
}

// FNV1a_32 returns FNV1a/32 checksum of the given data.
// [4]byte slice is returned.
// Technically, fnv.New32a() is used.
//   - https://pkg.go.dev/hash/fnv
//
//lint:ignore ST1003 should not use underscores in Go names; func FNV1a_32 should be FNV1a32
func FNV1a_32(b []byte) []byte {
	h := fnv.New32a()
	h.Write(b)
	return h.Sum(make([]byte, 0, SizeFNV1a_32))
}

// FNV1_64 returns FNV1/64 checksum of the given data.
// [8]byte slice is returned.
// Technically, fnv.New64() is used.
//   - https://pkg.go.dev/hash/fnv
//
//lint:ignore ST1003 should not use ALL_CAPS in Go names; use CamelCase instead
func FNV1_64(b []byte) []byte {
	h := fnv.New64()
	h.Write(b)
	return h.Sum(make([]byte, 0, SizeFNV1_64))
}

// FNV1a_64 returns FNV1a/64 checksum of the given data.
// [8]byte slice is returned.
// Technically, fnv.New64a() is used.
//   - https://pkg.go.dev/hash/fnv
//
//lint:ignore ST1003 should not use underscores in Go names; func FNV1a_64 should be FNV1a64
func FNV1a_64(b []byte) []byte {
	h := fnv.New64a()
	h.Write(b)
	return h.Sum(make([]byte, 0, SizeFNV1a_64))
}

// FNV1_128 returns FNV1/128 checksum of the given data.
// [16]byte slice is returned.
// Technically, fnv.New128() is used.
//   - https://pkg.go.dev/hash/fnv
//
//lint:ignore ST1003 should not use ALL_CAPS in Go names; use CamelCase instead
func FNV1_128(b []byte) []byte {
	h := fnv.New128()
	h.Write(b)
	return h.Sum(make([]byte, 0, SizeFNV1_128))
}

// FNV1a_128 returns FNV1a/128 checksum of the given data.
// [16]byte slice is returned.
// Technically, fnv.New128a() is used.
//   - https://pkg.go.dev/hash/fnv
//
//lint:ignore ST1003 should not use underscores in Go names; func FNV1a_128 should be FNV1a128
func FNV1a_128(b []byte) []byte {
	h := fnv.New128a()
	h.Write(b)
	return h.Sum(make([]byte, 0, SizeFNV1a_128))
}

// CRC32 returns CRC32 checksum of the given data.
// [4]byte slice is returned.
// Technically, crc32.New(crc32.IEEETable) is used.
//   - https://pkg.go.dev/hash/crc32
func CRC32(b []byte) []byte {
	h := crc32.New(crc32.IEEETable)
	h.Write(b)
	return h.Sum(make([]byte, 0, SizeCRC32))
}

// CRC64ISO returns CRC64 (using ISO table) checksum of the given data.
// [8]byte slice is returned.
// Technically, crc64.New(crc64ISOTable)is used.
//   - https://pkg.go.dev/hash/crc64
func CRC64ISO(b []byte) []byte {
	h := crc64.New(crc64ISOTable)
	h.Write(b)
	return h.Sum(make([]byte, 0, SizeCRC64ISO))
}

// CRC64ECMA returns CRC64 (using ECMA table) checksum of the given data.
// [8]byte slice is returned.
// Technically, crc64.New(crc64ECMATable)is used.
//   - https://pkg.go.dev/hash/crc64
func CRC64ECMA(b []byte) []byte {
	h := crc64.New(crc64ECMATable)
	h.Write(b)
	return h.Sum(make([]byte, 0, SizeCRC64ECMA))
}

// BLAKE2s_256 returns BLAKE2s/256 checksum of the given data.
// [32]byte slice is returned.
// Technically, blake2s.New256(nil) is used.
//   - https://pkg.go.dev/golang.org/x/crypto/blake2s
//
//lint:ignore ST1003 should not use underscores in Go names; func BLAKE2s_256 should be BLAKE2s256
func BLAKE2s_256(b []byte) []byte {
	h, _ := blake2s.New256(nil)
	h.Write(b)
	return h.Sum(make([]byte, 0, SizeBLAKE2s_256))
}

// BLAKE2b_256 returns BLAKE2b/256 checksum of the given data.
// [32]byte slice is returned.
// Technically, blake2b.New256(nil) is used.
//   - https://pkg.go.dev/golang.org/x/crypto/blake2b
//
//lint:ignore ST1003 should not use underscores in Go names; func BLAKE2b_256 should be BLAKE2b256
func BLAKE2b_256(b []byte) []byte {
	h, _ := blake2b.New256(nil)
	h.Write(b)
	return h.Sum(make([]byte, 0, SizeBLAKE2b_256))
}

// BLAKE2b_384 returns BLAKE2b/384 checksum of the given data.
// [48]byte slice is returned.
// Technically, blake2b.New384(nil) is used.
//   - https://pkg.go.dev/golang.org/x/crypto/blake2b
//
//lint:ignore ST1003 should not use underscores in Go names; func BLAKE2b_384 should be BLAKE2b384
func BLAKE2b_384(b []byte) []byte {
	h, _ := blake2b.New384(nil)
	h.Write(b)
	return h.Sum(make([]byte, 0, SizeBLAKE2b_384))
}

// BLAKE2b_512 returns BLAKE2b/512 checksum of the given data.
// [64]byte slice is returned.
// Technically, blake2b.New512(nil) is used.
//   - https://pkg.go.dev/golang.org/x/crypto/blake2b
//
//lint:ignore ST1003 should not use underscores in Go names; func BLAKE2b_512 should be BLAKE2b512
func BLAKE2b_512(b []byte) []byte {
	h, _ := blake2b.New512(nil)
	h.Write(b)
	return h.Sum(make([]byte, 0, SizeBLAKE2b_512))
}
