package mac

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"hash"
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

// HMACFunc is the function of Hash-based Message Authentication Code.
type HMACFunc func(message []byte, key []byte) (hash []byte)

// Algorithm is the HMAC underlying hash algorithm.
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
	AlgCRC32
	AlgCRC64ISO
	AlgCRC64ECMA
	AlgMD5
	AlgFNV1_32
	AlgFNV1a_32 //lint:ignore ST1003 should not use underscores in Go names; const AlgFNV1a_32 should be AlgFNV1a32
	AlgFNV1_64
	AlgFNV1a_64 //lint:ignore ST1003 should not use underscores in Go names; const AlgFNV1a_64 should be AlgFNV1a64
	AlgFNV1_128
	AlgFNV1a_128   //lint:ignore ST1003 should not use underscores in Go names; const AlgFNV1a_128 should be AlgFNV1a128
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

// FromAlgorithm returns HMAC function by searching with the given hash algorithm.
// This function returns nil when no HMAC function found.
func FromAlgorithm(a Algorithm) HMACFunc {
	algToFunc := map[Algorithm]HMACFunc{
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
		AlgBLAKE2s_256: BLAKE2s_256,
		AlgBLAKE2b_256: BLAKE2b_256,
		AlgBLAKE2b_384: BLAKE2b_384,
		AlgBLAKE2b_512: BLAKE2b_512,
	}
	return algToFunc[a]
}

// FromHashAlg returns HMAC function by searching with the given hash algorithm.
// This function returns nil when no HMAC function found.
func FromHashAlg(t k.HashAlg) HMACFunc {
	typeToFunc := map[k.HashAlg]HMACFunc{
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
		k.HashAlg_BLAKE2s_256: BLAKE2s_256,
		k.HashAlg_BLAKE2b_256: BLAKE2b_256,
		k.HashAlg_BLAKE2b_384: BLAKE2b_384,
		k.HashAlg_BLAKE2b_512: BLAKE2b_512,
	}
	return typeToFunc[t]
}

// SHA1 returns the HMAC (keyed-hash message authentication code) using sha1.
// [20]byte slice is returned.
func SHA1(msg, key []byte) []byte {
	mac := hmac.New(sha1.New, key)
	mac.Write(msg)
	return mac.Sum(make([]byte, 0, SizeSHA1))
}

// SHA224 returns the HMAC (keyed-hash message authentication code) using sha224.
// [28]byte slice is returned.
func SHA224(msg, key []byte) []byte {
	mac := hmac.New(sha256.New224, key)
	mac.Write(msg)
	return mac.Sum(make([]byte, 0, SizeSHA224))
}

// SHA256 returns the HMAC (keyed-hash message authentication code) using sha256.
// [32]byte slice is returned.
func SHA256(msg, key []byte) []byte {
	mac := hmac.New(sha256.New, key)
	mac.Write(msg)
	return mac.Sum(make([]byte, 0, SizeSHA256))
}

// SHA384 returns the HMAC (keyed-hash message authentication code) using sha384.
// [48]byte slice is returned.
func SHA384(msg, key []byte) []byte {
	mac := hmac.New(sha512.New384, key)
	mac.Write(msg)
	return mac.Sum(make([]byte, 0, SizeSHA384))
}

// SHA512 returns the HMAC (keyed-hash message authentication code) using sha512.
// [64]byte slice is returned.
func SHA512(msg, key []byte) []byte {
	mac := hmac.New(sha512.New, key)
	mac.Write(msg)
	return mac.Sum(make([]byte, 0, SizeSHA512))
}

// SHA512_224 returns the HMAC (keyed-hash message authentication code) using sha512/224.
// [28]byte slice is returned.
//
//lint:ignore ST1003 should not use ALL_CAPS in Go names; use CamelCase instead
func SHA512_224(msg, key []byte) []byte {
	mac := hmac.New(sha512.New512_224, key)
	mac.Write(msg)
	return mac.Sum(make([]byte, 0, SizeSHA512_224))
}

// SHA512_256 returns the HMAC (keyed-hash message authentication code) using sha512/256.
// [32]byte slice is returned.
//
//lint:ignore ST1003 should not use ALL_CAPS in Go names; use CamelCase instead
func SHA512_256(msg, key []byte) []byte {
	mac := hmac.New(sha512.New512_256, key)
	mac.Write(msg)
	return mac.Sum(make([]byte, 0, SizeSHA512_256))
}

// SHA3_224 returns the HMAC (keyed-hash message authentication code) using sha3/224.
// [28]byte slice is returned.
//
//lint:ignore ST1003 should not use ALL_CAPS in Go names; use CamelCase instead
func SHA3_224(msg, key []byte) []byte {
	mac := hmac.New(sha3.New224, key)
	mac.Write(msg)
	return mac.Sum(make([]byte, 0, SizeSHA3_224))
}

// SHA3_256 returns the HMAC (keyed-hash message authentication code) using sha3/256.
// [32]byte slice is returned.
//
//lint:ignore ST1003 should not use ALL_CAPS in Go names; use CamelCase instead
func SHA3_256(msg, key []byte) []byte {
	mac := hmac.New(sha3.New256, key)
	mac.Write(msg)
	return mac.Sum(make([]byte, 0, SizeSHA3_256))
}

// SHA3_384 returns the HMAC (keyed-hash message authentication code) using sha3/384.
// [48]byte slice is returned.
//
//lint:ignore ST1003 should not use ALL_CAPS in Go names; use CamelCase instead
func SHA3_384(msg, key []byte) []byte {
	mac := hmac.New(sha3.New384, key)
	mac.Write(msg)
	return mac.Sum(make([]byte, 0, SizeSHA3_384))
}

// SHA3_512 returns the HMAC (keyed-hash message authentication code) using sha3/512.
// [64]byte slice is returned.
//
//lint:ignore ST1003 should not use ALL_CAPS in Go names; use CamelCase instead
func SHA3_512(msg, key []byte) []byte {
	mac := hmac.New(sha3.New512, key)
	mac.Write(msg)
	return mac.Sum(make([]byte, 0, SizeSHA3_512))
}

// SHAKE128 returns the HMAC (keyed-hash message authentication code) using SHAKE128.
// [32]byte slice is returned.
func SHAKE128(msg, key []byte) []byte {
	mac := hmac.New(func() hash.Hash { return sha3.NewShake128() }, key)
	mac.Write(msg)
	return mac.Sum(make([]byte, 0, SizeSHAKE128))
}

// SHAKE256 returns the HMAC (keyed-hash message authentication code) using SHAKE256.
// [64]byte slice is returned.
func SHAKE256(msg, key []byte) []byte {
	mac := hmac.New(func() hash.Hash { return sha3.NewShake256() }, key)
	mac.Write(msg)
	return mac.Sum(make([]byte, 0, SizeSHAKE256))
}

// MD5 returns the HMAC (keyed-hash message authentication code) using md5.
// [16]byte slice is returned.
func MD5(msg, key []byte) []byte {
	mac := hmac.New(md5.New, key)
	mac.Write(msg)
	return mac.Sum(make([]byte, 0, SizeMD5))
}

// FNV1_32 returns the HMAC (keyed-hash message authentication code) using FNV1/32.
// Because FNV is not cryptography safe, do not use this for protecting sensitive data.
// [4]byte slice is returned.
//
//lint:ignore ST1003 should not use ALL_CAPS in Go names; use CamelCase instead
func FNV1_32(msg, key []byte) []byte {
	mac := hmac.New(func() hash.Hash { return fnv.New32() }, key)
	mac.Write(msg)
	return mac.Sum(make([]byte, 0, SizeFNV1_32))
}

// FNV1a_32 returns the HMAC (keyed-hash message authentication code) using FNV1a/32.
// Because FNV is not cryptography safe, do not use this for protecting sensitive data.
// [4]byte slice is returned.
//
//lint:ignore ST1003 should not use underscores in Go names; func FNV1a_32 should be FNV1a32
func FNV1a_32(msg, key []byte) []byte {
	mac := hmac.New(func() hash.Hash { return fnv.New32a() }, key)
	mac.Write(msg)
	return mac.Sum(make([]byte, 0, SizeFNV1a_32))
}

// FNV1_64 returns the HMAC (keyed-hash message authentication code) using FNV1/64.
// Because FNV is not cryptography safe, do not use this for protecting sensitive data.
// [8]byte slice is returned.
//
//lint:ignore ST1003 should not use ALL_CAPS in Go names; use CamelCase instead
func FNV1_64(msg, key []byte) []byte {
	mac := hmac.New(func() hash.Hash { return fnv.New64() }, key)
	mac.Write(msg)
	return mac.Sum(make([]byte, 0, SizeFNV1_64))
}

// FNV1a_64 returns the HMAC (keyed-hash message authentication code) using FNV1a/64.
// Because FNV is not cryptography safe, do not use this for protecting sensitive data.
// [8]byte slice is returned.
//
//lint:ignore ST1003 should not use underscores in Go names; func FNV1a_64 should be FNV1a64
func FNV1a_64(msg, key []byte) []byte {
	mac := hmac.New(func() hash.Hash { return fnv.New64a() }, key)
	mac.Write(msg)
	return mac.Sum(make([]byte, 0, SizeFNV1a_64))
}

// FNV1_128 returns the HMAC (keyed-hash message authentication code) using FNV1/128.
// Because FNV is not cryptography safe, do not use this for protecting sensitive data.
// [16]byte slice is returned.
//
//lint:ignore ST1003 should not use ALL_CAPS in Go names; use CamelCase instead
func FNV1_128(msg, key []byte) []byte {
	mac := hmac.New(fnv.New128, key)
	mac.Write(msg)
	return mac.Sum(make([]byte, 0, SizeFNV1_128))
}

// FNV1a_128 returns the HMAC (keyed-hash message authentication code) using FNV1a/128.
// Because FNV is not cryptography safe, do not use this for protecting sensitive data.
// [16]byte slice is returned.
//
//lint:ignore ST1003 should not use ALL_CAPS in Go names; use CamelCase instead
func FNV1a_128(msg, key []byte) []byte {
	mac := hmac.New(fnv.New128a, key)
	mac.Write(msg)
	return mac.Sum(make([]byte, 0, SizeFNV1a_128))
}

// CRC32 returns the HMAC (keyed-hash message authentication code) using CRC32.
// Because CRC32 is not cryptography safe, do not use this for protecting sensitive data.
// [4]byte slice is returned.
func CRC32(msg, key []byte) []byte {
	mac := hmac.New(func() hash.Hash {
		return crc32.New(crc32.IEEETable)
	}, key)
	mac.Write(msg)
	return mac.Sum(make([]byte, 0, SizeCRC32))
}

// CRC64ISO returns the HMAC (keyed-hash message authentication code) using CRC64ISO.
// Because CRC64 is not cryptography safe, do not use this for protecting sensitive data.
// [8]byte slice is returned.
func CRC64ISO(msg, key []byte) []byte {
	mac := hmac.New(func() hash.Hash {
		return crc64.New(crc64ISOTable)
	}, key)
	mac.Write(msg)
	return mac.Sum(make([]byte, 0, SizeCRC64ISO))
}

// CRC64ECMA returns the HMAC (keyed-hash message authentication code) using CRC64ECMA.
// Because CRC64 is not cryptography safe, do not use this for protecting sensitive data.
// [8]byte slice is returned.
func CRC64ECMA(msg, key []byte) []byte {
	mac := hmac.New(func() hash.Hash {
		return crc64.New(crc64ECMATable)
	}, key)
	mac.Write(msg)
	return mac.Sum(make([]byte, 0, SizeCRC64ECMA))
}

// BLAKE2s_256 returns the HMAC (keyed-hash message authentication code) using BLAKE2s/256.
// [32]byte slice is returned.
//
//lint:ignore ST1003 should not use ALL_CAPS in Go names; use CamelCase instead
func BLAKE2s_256(msg, key []byte) []byte {
	mac := hmac.New(func() hash.Hash {
		// Not use the built-in MAC mechanism.
		// Note: openssl also does not use built-in hmac for blake.
		h, _ := blake2s.New256(nil)
		return h
	}, key)
	mac.Write(msg)
	return mac.Sum(make([]byte, 0, SizeBLAKE2s_256))
}

// BLAKE2b_256 returns the HMAC (keyed-hash message authentication code) using BLAKE2b/256.
// [32]byte slice is returned.
//
//lint:ignore ST1003 should not use ALL_CAPS in Go names; use CamelCase instead
func BLAKE2b_256(msg, key []byte) []byte {
	mac := hmac.New(func() hash.Hash {
		// Not use the built-in MAC mechanism.
		// Note: openssl also does not use built-in hmac for blake.
		h, _ := blake2b.New256(nil)
		return h
	}, key)
	mac.Write(msg)
	return mac.Sum(make([]byte, 0, SizeBLAKE2b_256))
}

// BLAKE2b_384 returns the HMAC (keyed-hash message authentication code) using BLAKE2b/384.
// [48]byte slice is returned.
//
//lint:ignore ST1003 should not use ALL_CAPS in Go names; use CamelCase instead
func BLAKE2b_384(msg, key []byte) []byte {
	mac := hmac.New(func() hash.Hash {
		// Not use the built-in MAC mechanism.
		// Note: openssl also does not use built-in hmac for blake.
		h, _ := blake2b.New384(nil)
		return h
	}, key)
	mac.Write(msg)
	return mac.Sum(make([]byte, 0, SizeBLAKE2b_384))
}

// BLAKE2b_512 returns the HMAC (keyed-hash message authentication code) using BLAKE2b/512.
// [64]byte slice is returned.
//
//lint:ignore ST1003 should not use ALL_CAPS in Go names; use CamelCase instead
func BLAKE2b_512(msg, key []byte) []byte {
	mac := hmac.New(func() hash.Hash {
		// Not use the built-in MAC mechanism.
		// Note: openssl also does not use built-in hmac for blake.
		h, _ := blake2b.New512(nil)
		return h
	}, key)
	mac.Write(msg)
	return mac.Sum(make([]byte, 0, SizeBLAKE2b_512))
}
