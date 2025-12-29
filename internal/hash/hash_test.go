// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package hash_test

import (
	"testing"

	k "github.com/aileron-gateway/aileron-gateway/apis/kernel"
	"github.com/aileron-gateway/aileron-gateway/internal/hash"
	"github.com/aileron-gateway/aileron-gateway/internal/testutil"
	"github.com/aileron-projects/go/zcrypto/zblake2b"
	"github.com/aileron-projects/go/zcrypto/zblake2s"
	"github.com/aileron-projects/go/zcrypto/zsha1"
	"github.com/aileron-projects/go/zcrypto/zsha256"
	"github.com/aileron-projects/go/zcrypto/zsha3"
	"github.com/aileron-projects/go/zcrypto/zsha512"
	"github.com/google/go-cmp/cmp"
)

func TestFromAlgorithm(t *testing.T) {
	testCases := map[string]struct {
		alg hash.Algorithm
		f   hash.HashFunc
	}{
		"SHA1":        {alg: hash.AlgSHA1, f: zsha1.Sum},
		"SHA224":      {alg: hash.AlgSHA224, f: zsha256.Sum224},
		"SHA256":      {alg: hash.AlgSHA256, f: zsha256.Sum256},
		"SHA384":      {alg: hash.AlgSHA384, f: zsha512.Sum384},
		"SHA512":      {alg: hash.AlgSHA512, f: zsha512.Sum512},
		"SHA512_224":  {alg: hash.AlgSHA512_224, f: zsha512.Sum224},
		"SHA3_224":    {alg: hash.AlgSHA3_224, f: zsha3.Sum224},
		"SHA3_256":    {alg: hash.AlgSHA3_256, f: zsha3.Sum256},
		"SHA3_384":    {alg: hash.AlgSHA3_384, f: zsha3.Sum384},
		"SHA3_512":    {alg: hash.AlgSHA3_512, f: zsha3.Sum512},
		"SHAKE128":    {alg: hash.AlgSHAKE128, f: zsha3.SumShake128},
		"SHAKE256":    {alg: hash.AlgSHAKE256, f: zsha3.SumShake256},
		"BLAKE2s_256": {alg: hash.AlgBLAKE2s_256, f: zblake2s.Sum256},
		"BLAKE2b_256": {alg: hash.AlgBLAKE2b_256, f: zblake2b.Sum256},
		"BLAKE2b_384": {alg: hash.AlgBLAKE2b_384, f: zblake2b.Sum384},
		"BLAKE2b_512": {alg: hash.AlgBLAKE2b_512, f: zblake2b.Sum512},
		"not exist":   {alg: 9999, f: nil},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			h := hash.FromAlgorithm(tc.alg)
			testutil.Diff(t, tc.f, h, cmp.Comparer(testutil.ComparePointer[hash.HashFunc]))
			if h == nil {
				return
			}
			msg := []byte("example message")
			testutil.Diff(t, tc.f(msg), h(msg))
		})
	}
}

func TestFromHashAlg(t *testing.T) {
	testCases := map[string]struct {
		typ k.HashAlg
		f   hash.HashFunc
	}{
		"SHA1":        {typ: k.HashAlg_SHA1, f: zsha1.Sum},
		"SHA224":      {typ: k.HashAlg_SHA224, f: zsha256.Sum224},
		"SHA256":      {typ: k.HashAlg_SHA256, f: zsha256.Sum256},
		"SHA384":      {typ: k.HashAlg_SHA384, f: zsha512.Sum384},
		"SHA512":      {typ: k.HashAlg_SHA512, f: zsha512.Sum512},
		"SHA512_224":  {typ: k.HashAlg_SHA512_224, f: zsha512.Sum224},
		"SHA3_224":    {typ: k.HashAlg_SHA3_224, f: zsha3.Sum224},
		"SHA3_256":    {typ: k.HashAlg_SHA3_256, f: zsha3.Sum256},
		"SHA3_384":    {typ: k.HashAlg_SHA3_384, f: zsha3.Sum384},
		"SHA3_512":    {typ: k.HashAlg_SHA3_512, f: zsha3.Sum512},
		"SHAKE128":    {typ: k.HashAlg_SHAKE128, f: zsha3.SumShake128},
		"SHAKE256":    {typ: k.HashAlg_SHAKE256, f: zsha3.SumShake256},
		"BLAKE2s_256": {typ: k.HashAlg_BLAKE2s_256, f: zblake2s.Sum256},
		"BLAKE2b_256": {typ: k.HashAlg_BLAKE2b_256, f: zblake2b.Sum256},
		"BLAKE2b_384": {typ: k.HashAlg_BLAKE2b_384, f: zblake2b.Sum384},
		"BLAKE2b_512": {typ: k.HashAlg_BLAKE2b_512, f: zblake2b.Sum512},
		"not exist":   {typ: 9999, f: nil},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			h := hash.FromHashAlg(k.HashAlg(tc.typ))
			testutil.Diff(t, tc.f, h, cmp.Comparer(testutil.ComparePointer[hash.HashFunc]))
			if h == nil {
				return
			}
			msg := []byte("example message")
			testutil.Diff(t, tc.f(msg), h(msg))
		})
	}
}

func TestHMACFromAlgorithm(t *testing.T) {
	testCases := map[string]struct {
		alg hash.Algorithm
		f   hash.HMACFunc
	}{
		"SHA1":        {alg: hash.AlgSHA1, f: zsha1.HMACSum},
		"SHA224":      {alg: hash.AlgSHA224, f: zsha256.HMACSum224},
		"SHA256":      {alg: hash.AlgSHA256, f: zsha256.HMACSum256},
		"SHA384":      {alg: hash.AlgSHA384, f: zsha512.HMACSum384},
		"SHA512":      {alg: hash.AlgSHA512, f: zsha512.HMACSum512},
		"SHA512_224":  {alg: hash.AlgSHA512_224, f: zsha512.HMACSum224},
		"SHA3_224":    {alg: hash.AlgSHA3_224, f: zsha3.HMACSum224},
		"SHA3_256":    {alg: hash.AlgSHA3_256, f: zsha3.HMACSum256},
		"SHA3_384":    {alg: hash.AlgSHA3_384, f: zsha3.HMACSum384},
		"SHA3_512":    {alg: hash.AlgSHA3_512, f: zsha3.HMACSum512},
		"SHAKE128":    {alg: hash.AlgSHAKE128, f: nil},
		"SHAKE256":    {alg: hash.AlgSHAKE256, f: nil},
		"BLAKE2s_256": {alg: hash.AlgBLAKE2s_256, f: zblake2s.HMACSum256},
		"BLAKE2b_256": {alg: hash.AlgBLAKE2b_256, f: zblake2b.HMACSum256},
		"BLAKE2b_384": {alg: hash.AlgBLAKE2b_384, f: zblake2b.HMACSum384},
		"BLAKE2b_512": {alg: hash.AlgBLAKE2b_512, f: zblake2b.HMACSum512},
		"not exist":   {alg: 9999, f: nil},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			h := hash.HMACFromAlgorithm(tc.alg)
			testutil.Diff(t, tc.f, h, cmp.Comparer(testutil.ComparePointer[hash.HMACFunc]))
			if h == nil {
				return
			}
			msg := []byte("example message")
			key := []byte("secret for hmac")
			testutil.Diff(t, tc.f(msg, key), h(msg, key))
		})
	}
}

func TestHMACFromHashAlg(t *testing.T) {
	testCases := map[string]struct {
		typ k.HashAlg
		f   hash.HMACFunc
	}{
		"SHA1":        {typ: k.HashAlg_SHA1, f: zsha1.HMACSum},
		"SHA224":      {typ: k.HashAlg_SHA224, f: zsha256.HMACSum224},
		"SHA256":      {typ: k.HashAlg_SHA256, f: zsha256.HMACSum256},
		"SHA384":      {typ: k.HashAlg_SHA384, f: zsha512.HMACSum384},
		"SHA512":      {typ: k.HashAlg_SHA512, f: zsha512.HMACSum512},
		"SHA512_224":  {typ: k.HashAlg_SHA512_224, f: zsha512.HMACSum224},
		"SHA3_224":    {typ: k.HashAlg_SHA3_224, f: zsha3.HMACSum224},
		"SHA3_256":    {typ: k.HashAlg_SHA3_256, f: zsha3.HMACSum256},
		"SHA3_384":    {typ: k.HashAlg_SHA3_384, f: zsha3.HMACSum384},
		"SHA3_512":    {typ: k.HashAlg_SHA3_512, f: zsha3.HMACSum512},
		"SHAKE128":    {typ: k.HashAlg_SHAKE128, f: nil},
		"SHAKE256":    {typ: k.HashAlg_SHAKE256, f: nil},
		"BLAKE2s_256": {typ: k.HashAlg_BLAKE2s_256, f: zblake2s.HMACSum256},
		"BLAKE2b_256": {typ: k.HashAlg_BLAKE2b_256, f: zblake2b.HMACSum256},
		"BLAKE2b_384": {typ: k.HashAlg_BLAKE2b_384, f: zblake2b.HMACSum384},
		"BLAKE2b_512": {typ: k.HashAlg_BLAKE2b_512, f: zblake2b.HMACSum512},
		"not exist":   {typ: 9999, f: nil},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			h := hash.HMACFromHashAlg(tc.typ)
			testutil.Diff(t, tc.f, h, cmp.Comparer(testutil.ComparePointer[hash.HMACFunc]))
			if h == nil {
				return
			}
			msg := []byte("example message")
			key := []byte("secret for hmac")
			testutil.Diff(t, tc.f(msg, key), h(msg, key))
		})
	}
}
