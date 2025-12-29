// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package encrypt

import (
	"cmp"
	"crypto"

	"github.com/aileron-gateway/aileron-gateway/apis/kernel"
	"github.com/aileron-projects/go/zcrypto/zargon2"
	"github.com/aileron-projects/go/zcrypto/zbcrypt"
	"github.com/aileron-projects/go/zcrypto/zpbkdf2"
	"github.com/aileron-projects/go/zcrypto/zscrypt"
)

type PasswordCrypt interface {
	Sum(password []byte) ([]byte, error)
	Compare(hashed []byte, password []byte) error
}

// PasswordHashFunc is the type of function that returns
// hash of the given password.
type PasswordHashFunc func(password []byte) (hashed []byte, err error)

// PasswordCompareFunc is the type of function that compare two passwords.
// This function returns nil when the passwords are correct
// and non-nil error when in-correct.
type PasswordCompareFunc func(hashed []byte, password []byte) error

// NewPasswordCrypt returns an password crypt object.
// NewPasswordCrypt return nil object and nil error
// when nil spec were given or unknown crypt type was given.
func NewPasswordCrypt(spec *kernel.PasswordCryptSpec) (PasswordCrypt, error) {
	if spec == nil {
		return nil, nil
	}
	switch v := spec.PasswordCrypts.(type) {
	case *kernel.PasswordCryptSpec_BCrypt:
		return NewBCrypt(v.BCrypt)
	case *kernel.PasswordCryptSpec_SCrypt:
		return NewSCrypt(v.SCrypt)
	case *kernel.PasswordCryptSpec_PBKDF2:
		return NewPBKDF2(v.PBKDF2)
	case *kernel.PasswordCryptSpec_Argon2I:
		return NewArgon2i(v.Argon2I)
	case *kernel.PasswordCryptSpec_Argon2Id:
		return NewArgon2id(v.Argon2Id)
	default:
		return nil, nil
	}
}

func NewBCrypt(spec *kernel.BCryptSpec) (*zbcrypt.BCrypt, error) {
	cost := cmp.Or(int(spec.Cost), 10)
	return zbcrypt.New(cost)
}

func NewSCrypt(spec *kernel.SCryptSpec) (*zscrypt.SCrypt, error) {
	saltLen := cmp.Or(int(spec.SaltLen), 32)
	n := cmp.Or(int(spec.N), 32768)
	r := cmp.Or(int(spec.R), 8)
	p := cmp.Or(int(spec.P), 1)
	keyLen := cmp.Or(int(spec.KeyLen), 32)
	return zscrypt.New(saltLen, n, r, p, keyLen)
}

func NewPBKDF2(spec *kernel.PBKDF2Spec) (*zpbkdf2.PBKDF2, error) {
	h := map[kernel.HashAlg]crypto.Hash{
		kernel.HashAlg_HashAlgUnknown: crypto.SHA256, // Default SHA256
		kernel.HashAlg_SHA1:           crypto.SHA1,
		kernel.HashAlg_SHA224:         crypto.SHA224,
		kernel.HashAlg_SHA256:         crypto.SHA256,
		kernel.HashAlg_SHA384:         crypto.SHA384,
		kernel.HashAlg_SHA512_224:     crypto.SHA512_224,
		kernel.HashAlg_SHA512_256:     crypto.SHA512_256,
		kernel.HashAlg_SHA3_224:       crypto.SHA3_224,
		kernel.HashAlg_SHA3_256:       crypto.SHA3_256,
		kernel.HashAlg_SHA3_384:       crypto.SHA3_384,
		kernel.HashAlg_SHA3_512:       crypto.SHA3_512,
	}
	saltLen := cmp.Or(int(spec.SaltLen), 32)
	iter := cmp.Or(int(spec.Iter), 4096)
	keyLen := cmp.Or(int(spec.KeyLen), 32)
	return zpbkdf2.New(saltLen, iter, keyLen, h[spec.HashAlg])
}

func NewArgon2i(spec *kernel.Argon2Spec) (*zargon2.Argon2i, error) {
	saltLen := cmp.Or(int(spec.SaltLen), 32)
	time := cmp.Or(spec.Time, 3)
	memory := cmp.Or(spec.Memory, 32*1024)
	threads := cmp.Or(uint8(spec.Threads), 4) //nolint:gosec // G115: integer overflow conversion int32 -> uint8
	keyLen := cmp.Or(spec.KeyLen, 32)
	return zargon2.NewArgon2i(saltLen, time, memory, threads, keyLen)
}

func NewArgon2id(spec *kernel.Argon2Spec) (*zargon2.Argon2id, error) {
	saltLen := cmp.Or(int(spec.SaltLen), 32)
	time := cmp.Or(spec.Time, 1)
	memory := cmp.Or(spec.Memory, 64*1024)
	threads := cmp.Or(uint8(spec.Threads), 4) //nolint:gosec // G115: integer overflow conversion int32 -> uint8
	keyLen := cmp.Or(spec.KeyLen, 32)
	return zargon2.NewArgon2id(saltLen, time, memory, threads, keyLen)
}
