// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package hash

import (
	k "github.com/aileron-gateway/aileron-gateway/apis/kernel"
	"github.com/aileron-projects/go/zcrypto/zblake2b"
	"github.com/aileron-projects/go/zcrypto/zblake2s"
	"github.com/aileron-projects/go/zcrypto/zsha1"
	"github.com/aileron-projects/go/zcrypto/zsha256"
	"github.com/aileron-projects/go/zcrypto/zsha3"
	"github.com/aileron-projects/go/zcrypto/zsha512"
)

// HashFunc is the function that returns hash value of the given data.
type HashFunc func(data []byte) (hash []byte)

// HMACFunc is the function of Hash-based Message Authentication Code.
type HMACFunc func(message []byte, key []byte) (hash []byte)

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
		AlgSHA1:        zsha1.Sum,
		AlgSHA224:      zsha256.Sum224,
		AlgSHA256:      zsha256.Sum256,
		AlgSHA384:      zsha512.Sum384,
		AlgSHA512:      zsha512.Sum512,
		AlgSHA512_224:  zsha512.Sum224,
		AlgSHA512_256:  zsha512.Sum256,
		AlgSHA3_224:    zsha3.Sum224,
		AlgSHA3_256:    zsha3.Sum256,
		AlgSHA3_384:    zsha3.Sum384,
		AlgSHA3_512:    zsha3.Sum512,
		AlgSHAKE128:    zsha3.SumShake128,
		AlgSHAKE256:    zsha3.SumShake256,
		AlgBLAKE2s_256: zblake2s.Sum256,
		AlgBLAKE2b_256: zblake2b.Sum256,
		AlgBLAKE2b_384: zblake2b.Sum384,
		AlgBLAKE2b_512: zblake2b.Sum512,
	}
	return algToFunc[a]
}

// FromHashAlg returns hash function
// corresponding to the given hash algorithm.
// This function returns nil when hash function not found.
func FromHashAlg(t k.HashAlg) HashFunc {
	typeToFunc := map[k.HashAlg]HashFunc{
		k.HashAlg_SHA1:        zsha1.Sum,
		k.HashAlg_SHA224:      zsha256.Sum224,
		k.HashAlg_SHA256:      zsha256.Sum256,
		k.HashAlg_SHA384:      zsha512.Sum384,
		k.HashAlg_SHA512:      zsha512.Sum512,
		k.HashAlg_SHA512_224:  zsha512.Sum224,
		k.HashAlg_SHA512_256:  zsha512.Sum256,
		k.HashAlg_SHA3_224:    zsha3.Sum224,
		k.HashAlg_SHA3_256:    zsha3.Sum256,
		k.HashAlg_SHA3_384:    zsha3.Sum384,
		k.HashAlg_SHA3_512:    zsha3.Sum512,
		k.HashAlg_SHAKE128:    zsha3.SumShake128,
		k.HashAlg_SHAKE256:    zsha3.SumShake256,
		k.HashAlg_BLAKE2s_256: zblake2s.Sum256,
		k.HashAlg_BLAKE2b_256: zblake2b.Sum256,
		k.HashAlg_BLAKE2b_384: zblake2b.Sum384,
		k.HashAlg_BLAKE2b_512: zblake2b.Sum512,
	}
	return typeToFunc[t]
}

// FromAlgorithm returns HMAC function by searching with the given hash algorithm.
// This function returns nil when no HMAC function found.
func HMACFromAlgorithm(a Algorithm) HMACFunc {
	algToFunc := map[Algorithm]HMACFunc{
		AlgSHA1:        zsha1.HMACSum,
		AlgSHA224:      zsha256.HMACSum224,
		AlgSHA256:      zsha256.HMACSum256,
		AlgSHA384:      zsha512.HMACSum384,
		AlgSHA512:      zsha512.HMACSum512,
		AlgSHA512_224:  zsha512.HMACSum224,
		AlgSHA512_256:  zsha512.HMACSum256,
		AlgSHA3_224:    zsha3.HMACSum224,
		AlgSHA3_256:    zsha3.HMACSum256,
		AlgSHA3_384:    zsha3.HMACSum384,
		AlgSHA3_512:    zsha3.HMACSum512,
		AlgSHAKE128:    nil,
		AlgSHAKE256:    nil,
		AlgBLAKE2s_256: zblake2s.HMACSum256,
		AlgBLAKE2b_256: zblake2b.HMACSum256,
		AlgBLAKE2b_384: zblake2b.HMACSum384,
		AlgBLAKE2b_512: zblake2b.HMACSum512,
	}
	return algToFunc[a]
}

// FromHashAlg returns HMAC function by searching with the given hash algorithm.
// This function returns nil when no HMAC function found.
func HMACFromHashAlg(t k.HashAlg) HMACFunc {
	typeToFunc := map[k.HashAlg]HMACFunc{
		k.HashAlg_SHA1:        zsha1.HMACSum,
		k.HashAlg_SHA224:      zsha256.HMACSum224,
		k.HashAlg_SHA256:      zsha256.HMACSum256,
		k.HashAlg_SHA384:      zsha512.HMACSum384,
		k.HashAlg_SHA512:      zsha512.HMACSum512,
		k.HashAlg_SHA512_224:  zsha512.HMACSum224,
		k.HashAlg_SHA512_256:  zsha512.HMACSum256,
		k.HashAlg_SHA3_224:    zsha3.HMACSum224,
		k.HashAlg_SHA3_256:    zsha3.HMACSum256,
		k.HashAlg_SHA3_384:    zsha3.HMACSum384,
		k.HashAlg_SHA3_512:    zsha3.HMACSum512,
		k.HashAlg_SHAKE128:    nil,
		k.HashAlg_SHAKE256:    nil,
		k.HashAlg_BLAKE2s_256: zblake2s.HMACSum256,
		k.HashAlg_BLAKE2b_256: zblake2b.HMACSum256,
		k.HashAlg_BLAKE2b_384: zblake2b.HMACSum384,
		k.HashAlg_BLAKE2b_512: zblake2b.HMACSum512,
	}
	return typeToFunc[t]
}
