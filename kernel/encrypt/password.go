// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package encrypt

import (
	"cmp"
	"crypto/rand"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"errors"
	"hash"
	"io"

	k "github.com/aileron-gateway/aileron-gateway/apis/kernel"
	"github.com/aileron-gateway/aileron-gateway/kernel/er"
	"golang.org/x/crypto/argon2"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/pbkdf2"
	"golang.org/x/crypto/scrypt"
	"golang.org/x/crypto/sha3"
)

var ErrNotMatch = errors.New("password hash not matched")

type PasswordCrypt interface {
	Hash(password []byte) ([]byte, error)
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
// when nil spec were given.
func NewPasswordCrypt(spec *k.PasswordCryptSpec) (PasswordCrypt, error) {
	if spec == nil {
		return nil, nil
	}
	switch v := spec.PasswordCrypts.(type) {
	case *k.PasswordCryptSpec_BCrypt:
		return NewBCrypt(v.BCrypt)
	case *k.PasswordCryptSpec_SCrypt:
		return NewSCrypt(v.SCrypt)
	case *k.PasswordCryptSpec_PBKDF2:
		return NewPBKDF2(v.PBKDF2)
	case *k.PasswordCryptSpec_Argon2I:
		return NewArgon2i(v.Argon2I)
	case *k.PasswordCryptSpec_Argon2Id:
		return NewArgon2id(v.Argon2Id)
	default:
		return nil, nil
	}
}

func NewBCrypt(spec *k.BCryptSpec) (*BCrypt, error) {
	crypt := &BCrypt{
		cost: cmp.Or(int(spec.Cost), bcrypt.DefaultCost),
	}
	if crypt.cost < bcrypt.MinCost || crypt.cost > bcrypt.MaxCost {
		return nil, &er.Error{
			Package:     ErrPkg,
			Type:        ErrTypeBCrypt,
			Description: ErrDscHashValid,
			Detail:      "cost out of range.",
		}
	}
	return crypt, nil
}

type BCrypt struct {
	cost int
}

func (c *BCrypt) Hash(password []byte) ([]byte, error) {
	hashed, err := bcrypt.GenerateFromPassword(password, c.cost)
	if err != nil {
		return nil, (&er.Error{
			Package:     ErrPkg,
			Type:        ErrTypeBCrypt,
			Description: ErrDscHash,
		}).Wrap(err)
	}
	return hashed, nil
}

func (c *BCrypt) Compare(hashed []byte, password []byte) error {
	err := bcrypt.CompareHashAndPassword(hashed, password)
	if err != nil {
		return ErrNotMatch
	}
	return nil
}

func NewSCrypt(spec *k.SCryptSpec) (*SCrypt, error) {
	crypt := &SCrypt{
		saltLen: cmp.Or(int(spec.SaltLen), 32),
		n:       cmp.Or(int(spec.N), 32768),
		r:       cmp.Or(int(spec.R), 8),
		p:       cmp.Or(int(spec.P), 1),
		keyLen:  cmp.Or(int(spec.KeyLen), 32),
	}
	if _, err := crypt.Hash([]byte("test")); err != nil {
		return nil, (&er.Error{
			Package:     ErrPkg,
			Type:        ErrTypeSCrypt,
			Description: ErrDscHashValid,
		}).Wrap(err.(*er.Error).Unwrap())
	}
	return crypt, nil
}

type SCrypt struct {
	saltLen int
	n       int
	r       int
	p       int
	keyLen  int
}

func (c *SCrypt) Hash(password []byte) ([]byte, error) {
	salt := make([]byte, c.saltLen)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return nil, (&er.Error{
			Package:     ErrPkg,
			Type:        ErrTypeSCrypt,
			Description: ErrDscHash,
		}).Wrap(err)
	}
	hashed, err := scrypt.Key(password, salt, c.n, c.r, c.p, c.keyLen)
	if err != nil {
		return nil, (&er.Error{
			Package:     ErrPkg,
			Type:        ErrTypeSCrypt,
			Description: ErrDscHash,
		}).Wrap(err)
	}
	return append(salt, hashed...), nil
}

func (c *SCrypt) Compare(hashed []byte, password []byte) error {
	n := len(hashed) - c.saltLen
	if n < 0 {
		return ErrNotMatch
	}
	salt, hashed := hashed[:c.saltLen], hashed[c.saltLen:]
	x, err := scrypt.Key(password, salt, c.n, c.r, c.p, c.keyLen)
	if err != nil {
		return ErrNotMatch
	}
	if string(x) != string(hashed) {
		return ErrNotMatch
	}
	return nil
}

func NewPBKDF2(spec *k.PBKDF2Spec) (*PBKDF2, error) {
	funcs := map[k.HashAlg]func() hash.Hash{
		k.HashAlg_HashAlgUnknown: sha256.New, // Default SHA256
		k.HashAlg_SHA1:           sha1.New,
		k.HashAlg_SHA224:         sha256.New224,
		k.HashAlg_SHA256:         sha256.New,
		k.HashAlg_SHA384:         sha512.New384,
		k.HashAlg_SHA512_224:     sha512.New512_224,
		k.HashAlg_SHA512_256:     sha512.New512_256,
		k.HashAlg_SHA3_224:       sha3.New224,
		k.HashAlg_SHA3_256:       sha3.New256,
		k.HashAlg_SHA3_384:       sha3.New384,
		k.HashAlg_SHA3_512:       sha3.New512,
	}
	hashFunc, ok := funcs[spec.HashAlg]
	if !ok {
		return nil, &er.Error{
			Package:     ErrPkg,
			Type:        ErrTypePBKDF2,
			Description: ErrDscHashValid,
			Detail:      "unsupported hash function. given " + spec.HashAlg.String(),
		}
	}
	crypt := &PBKDF2{
		saltLen:  cmp.Or(int(spec.SaltLen), 32),
		iter:     cmp.Or(int(spec.Iter), 4096),
		keyLen:   cmp.Or(int(spec.KeyLen), 32),
		hashFunc: hashFunc,
	}
	return crypt, nil
}

type PBKDF2 struct {
	saltLen  int
	iter     int
	keyLen   int
	hashFunc func() hash.Hash
}

func (c *PBKDF2) Hash(password []byte) ([]byte, error) {
	salt := make([]byte, c.saltLen)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return nil, (&er.Error{
			Package:     ErrPkg,
			Type:        ErrTypePBKDF2,
			Description: ErrDscHash,
		}).Wrap(err)
	}
	x := pbkdf2.Key(password, salt, c.iter, c.keyLen, c.hashFunc)
	return append(salt, x...), nil
}

func (c *PBKDF2) Compare(hashed []byte, password []byte) error {
	n := len(hashed) - c.saltLen
	if n < 0 {
		return ErrNotMatch
	}
	salt, hashed := hashed[:c.saltLen], hashed[c.saltLen:]
	x := pbkdf2.Key(password, salt, c.iter, c.keyLen, c.hashFunc)
	if string(x) != string(hashed) {
		return ErrNotMatch
	}
	return nil
}

func NewArgon2i(spec *k.Argon2Spec) (*Argon2i, error) {
	crypt := &Argon2i{
		saltLen: cmp.Or(int(spec.SaltLen), 32),
		time:    cmp.Or(spec.Time, 3),
		memory:  cmp.Or(spec.Memory, 32*1024),
		threads: cmp.Or(uint8(spec.Threads), 4), //nolint:gosec // G115: integer overflow conversion int32 -> uint8
		keyLen:  cmp.Or(spec.KeyLen, 32),
	}
	return crypt, nil
}

type Argon2i struct {
	saltLen int
	time    uint32
	memory  uint32
	threads uint8
	keyLen  uint32
}

func (c *Argon2i) Hash(password []byte) ([]byte, error) {
	salt := make([]byte, c.saltLen)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return nil, (&er.Error{
			Package:     ErrPkg,
			Type:        ErrTypeArgon2i,
			Description: ErrDscHash,
		}).Wrap(err)
	}
	x := argon2.Key(password, salt, c.time, c.memory, c.threads, c.keyLen)
	return append(salt, x...), nil
}

func (c *Argon2i) Compare(hashed []byte, password []byte) error {
	n := len(hashed) - c.saltLen
	if n < 0 {
		return ErrNotMatch
	}
	salt, hashed := hashed[:c.saltLen], hashed[c.saltLen:]
	x := argon2.Key(password, salt, c.time, c.memory, c.threads, c.keyLen)
	if string(x) != string(hashed) {
		return ErrNotMatch
	}
	return nil
}

func NewArgon2id(spec *k.Argon2Spec) (*Argon2id, error) {
	crypt := &Argon2id{
		saltLen: cmp.Or(int(spec.SaltLen), 32),
		time:    cmp.Or(spec.Time, 1),
		memory:  cmp.Or(spec.Memory, 64*1024),
		threads: cmp.Or(uint8(spec.Threads), 4), //nolint:gosec // G115: integer overflow conversion int32 -> uint8
		keyLen:  cmp.Or(spec.KeyLen, 32),
	}
	return crypt, nil
}

type Argon2id struct {
	saltLen int
	time    uint32
	memory  uint32
	threads uint8
	keyLen  uint32
}

func (c *Argon2id) Hash(password []byte) ([]byte, error) {
	salt := make([]byte, c.saltLen)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return nil, (&er.Error{
			Package:     ErrPkg,
			Type:        ErrTypeArgon2id,
			Description: ErrDscHash,
		}).Wrap(err)
	}
	x := argon2.IDKey(password, salt, c.time, c.memory, c.threads, c.keyLen)
	return append(salt, x...), nil
}

func (c *Argon2id) Compare(hashed []byte, password []byte) error {
	n := len(hashed) - c.saltLen
	if n < 0 {
		return ErrNotMatch
	}
	salt, hashed := hashed[:c.saltLen], hashed[c.saltLen:]
	x := argon2.IDKey(password, salt, c.time, c.memory, c.threads, c.keyLen)
	if string(x) != string(hashed) {
		return ErrNotMatch
	}
	return nil
}
