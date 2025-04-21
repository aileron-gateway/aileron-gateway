// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package encrypt

import (
	"bytes"
	"crypto/rand"
	"crypto/sha1"
	"crypto/sha256"
	"hash"
	"io"
	"testing"

	k "github.com/aileron-gateway/aileron-gateway/apis/kernel"
	"github.com/aileron-gateway/aileron-gateway/kernel/er"
	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"golang.org/x/crypto/bcrypt"
)

func TestNewPasswordCrypt(t *testing.T) {
	type condition struct {
		spec *k.PasswordCryptSpec
	}

	type action struct {
		c any
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"nil",
			[]string{},
			[]string{},
			&condition{
				spec: nil,
			},
			&action{
				c: nil,
			},
		),
		gen(
			"unknown",
			[]string{},
			[]string{},
			&condition{
				spec: &k.PasswordCryptSpec{
					PasswordCrypts: nil,
				},
			},
			&action{
				c: nil,
			},
		),
		gen(
			"BCrypt",
			[]string{},
			[]string{},
			&condition{
				spec: &k.PasswordCryptSpec{
					PasswordCrypts: &k.PasswordCryptSpec_BCrypt{
						BCrypt: &k.BCryptSpec{},
					},
				},
			},
			&action{
				c: &BCrypt{cost: 10},
			},
		),
		gen(
			"SCrypt",
			[]string{},
			[]string{},
			&condition{
				spec: &k.PasswordCryptSpec{
					PasswordCrypts: &k.PasswordCryptSpec_SCrypt{
						SCrypt: &k.SCryptSpec{},
					},
				},
			},
			&action{
				c: &SCrypt{
					saltLen: 32,
					n:       32768,
					r:       8,
					p:       1,
					keyLen:  32,
				},
			},
		),
		gen(
			"PBKDF2",
			[]string{},
			[]string{},
			&condition{
				spec: &k.PasswordCryptSpec{
					PasswordCrypts: &k.PasswordCryptSpec_PBKDF2{
						PBKDF2: &k.PBKDF2Spec{},
					},
				},
			},
			&action{
				c: &PBKDF2{
					saltLen:  32,
					iter:     4096,
					keyLen:   32,
					hashFunc: sha256.New,
				},
			},
		),
		gen(
			"Argon2i",
			[]string{},
			[]string{},
			&condition{
				spec: &k.PasswordCryptSpec{
					PasswordCrypts: &k.PasswordCryptSpec_Argon2I{
						Argon2I: &k.Argon2Spec{},
					},
				},
			},
			&action{
				c: &Argon2i{
					saltLen: 32,
					time:    3,
					memory:  32768,
					threads: 4,
					keyLen:  32,
				},
			},
		),
		gen(
			"Argon2id",
			[]string{},
			[]string{},
			&condition{
				spec: &k.PasswordCryptSpec{
					PasswordCrypts: &k.PasswordCryptSpec_Argon2Id{
						Argon2Id: &k.Argon2Spec{},
					},
				},
			},
			&action{
				c: &Argon2id{
					saltLen: 32,
					time:    1,
					memory:  65536,
					threads: 4,
					keyLen:  32,
				},
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			c, _ := NewPasswordCrypt(tt.C().spec)

			opts := []cmp.Option{
				cmp.AllowUnexported(BCrypt{}, SCrypt{}),
				cmp.AllowUnexported(PBKDF2{}),
				cmp.AllowUnexported(Argon2i{}, Argon2id{}),
				cmp.Comparer(testutil.ComparePointer[func() hash.Hash]),
			}
			testutil.Diff(t, tt.A().c, c, opts...)
		})
	}
}

func TestNewBCrypt(t *testing.T) {
	type condition struct {
		spec *k.BCryptSpec
	}

	type action struct {
		c   *BCrypt
		err error
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"empty spec",
			[]string{},
			[]string{},
			&condition{
				spec: &k.BCryptSpec{},
			},
			&action{
				c: &BCrypt{cost: 10},
			},
		),
		gen(
			"min cost",
			[]string{},
			[]string{},
			&condition{
				spec: &k.BCryptSpec{Cost: int32(bcrypt.MinCost)},
			},
			&action{
				c: &BCrypt{cost: bcrypt.MinCost},
			},
		),
		gen(
			"max cost",
			[]string{},
			[]string{},
			&condition{
				spec: &k.BCryptSpec{Cost: int32(bcrypt.MaxCost)},
			},
			&action{
				c: &BCrypt{cost: bcrypt.MaxCost},
			},
		),
		gen(
			"invalid cost",
			[]string{},
			[]string{},
			&condition{
				spec: &k.BCryptSpec{Cost: int32(bcrypt.MinCost - 1)},
			},
			&action{
				err: &er.Error{
					Package:     ErrPkg,
					Type:        ErrTypeBCrypt,
					Description: ErrDscHashValid,
				},
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			c, err := NewBCrypt(tt.C().spec)
			testutil.Diff(t, tt.A().err, err, cmpopts.EquateErrors())
			testutil.Diff(t, tt.A().c, c, cmp.AllowUnexported(BCrypt{}))
		})
	}
}

func TestBCrypt(t *testing.T) {
	type condition struct {
		c          *BCrypt
		passwd     []byte
		passwdRepl []byte
	}

	type action struct {
		prefix []byte
		err    error
		match  bool
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	cndInputNonNilPassword := tb.Condition("non-nil password", "input non-zero or non-nil password")
	cndInputNilPassword := tb.Condition("nil password", "input nil password")
	cndInputTooLongPassword := tb.Condition("too long password", "input too long password. length grater than 72")
	actCheckNoError := tb.Action("check no error", "check that there is no error")
	actCheckError := tb.Action("check error", "check that there is an error")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"non-nil password",
			[]string{cndInputNonNilPassword},
			[]string{actCheckNoError},
			&condition{
				c:      &BCrypt{},
				passwd: []byte("password"),
			},
			&action{
				prefix: []byte("$2a$10$"),
				match:  true,
			},
		),
		gen(
			"nil password",
			[]string{cndInputNilPassword},
			[]string{actCheckNoError},
			&condition{
				c:      &BCrypt{},
				passwd: nil,
			},
			&action{
				prefix: []byte("$2a$10$"),
				match:  true,
			},
		),
		gen(
			"hash not match",
			[]string{cndInputNilPassword},
			[]string{actCheckNoError},
			&condition{
				c:          &BCrypt{},
				passwd:     []byte("test input"),
				passwdRepl: []byte("replaced password"),
			},
			&action{
				prefix: []byte("$2a$10$"),
				match:  false,
			},
		),
		gen(
			"too long password",
			[]string{cndInputTooLongPassword},
			[]string{actCheckError},
			&condition{
				c:      &BCrypt{},
				passwd: []byte("too long password. length is grater than 72. too long password. password length is grater than 72."),
			},
			&action{
				prefix: []byte(""),
				match:  false,
				err: &er.Error{
					Package:     ErrPkg,
					Type:        ErrTypeBCrypt,
					Description: ErrDscHash,
				},
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			got, err := tt.C().c.Hash(tt.C().passwd)
			testutil.Diff(t, tt.A().err, err, cmpopts.EquateErrors())
			testutil.Diff(t, true, bytes.HasPrefix(got, tt.A().prefix))

			if len(tt.C().passwdRepl) > 0 {
				tt.C().passwd = tt.C().passwdRepl
			}
			err = tt.C().c.Compare(got, tt.C().passwd)
			testutil.Diff(t, tt.A().match, err == nil)
		})
	}
}

func TestNewSCrypt(t *testing.T) {
	type condition struct {
		spec *k.SCryptSpec
	}

	type action struct {
		c   *SCrypt
		err error
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"empty spec",
			[]string{},
			[]string{},
			&condition{
				spec: &k.SCryptSpec{},
			},
			&action{
				c: &SCrypt{
					saltLen: 32,
					n:       32768,
					r:       8,
					p:       1,
					keyLen:  32,
				},
			},
		),
		gen(
			"full spec",
			[]string{},
			[]string{},
			&condition{
				spec: &k.SCryptSpec{
					SaltLen: 16,
					N:       1024 * 4,
					R:       2,
					P:       3,
					KeyLen:  8,
				},
			},
			&action{
				c: &SCrypt{
					saltLen: 16,
					n:       1024 * 4,
					r:       2,
					p:       3,
					keyLen:  8,
				},
			},
		),
		gen(
			"invalid N",
			[]string{},
			[]string{},
			&condition{
				spec: &k.SCryptSpec{N: -1},
			},
			&action{
				c: nil,
				err: &er.Error{
					Package:     ErrPkg,
					Type:        ErrTypeSCrypt,
					Description: ErrDscHashValid,
				},
			},
		),
		gen(
			"invalid r",
			[]string{},
			[]string{},
			&condition{
				spec: &k.SCryptSpec{R: -1},
			},
			&action{
				c: nil,
				err: &er.Error{
					Package:     ErrPkg,
					Type:        ErrTypeSCrypt,
					Description: ErrDscHashValid,
				},
			},
		),
		gen(
			"invalid p",
			[]string{},
			[]string{},
			&condition{
				spec: &k.SCryptSpec{P: -1},
			},
			&action{
				c: nil,
				err: &er.Error{
					Package:     ErrPkg,
					Type:        ErrTypeSCrypt,
					Description: ErrDscHashValid,
				},
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			c, err := NewSCrypt(tt.C().spec)
			testutil.Diff(t, tt.A().err, err, cmpopts.EquateErrors())
			testutil.Diff(t, tt.A().c, c, cmp.AllowUnexported(SCrypt{}))
		})
	}
}

func TestSCrypt(t *testing.T) {
	type condition struct {
		c          *SCrypt
		saltReader io.Reader
		passwd     []byte
		passwdRepl []byte  // Replaced before compare
		cryptRepl  *SCrypt // Replaced before compare
	}

	type action struct {
		err   error
		match bool
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"nil password",
			[]string{},
			[]string{},
			&condition{
				c: &SCrypt{
					saltLen: 32,
					n:       32768,
					r:       8,
					p:       1,
					keyLen:  32,
				},
				passwd:     nil,
				saltReader: bytes.NewReader([]byte("12345678901234567890123456789012")),
			},
			&action{
				match: true,
			},
		),
		gen(
			"non-nil password",
			[]string{},
			[]string{},
			&condition{
				c: &SCrypt{
					saltLen: 32,
					n:       32768,
					r:       8,
					p:       1,
					keyLen:  32,
				},
				passwd:     []byte("password"),
				saltReader: bytes.NewReader([]byte("12345678901234567890123456789012")),
			},
			&action{
				match: true,
			},
		),
		gen(
			"password not match",
			[]string{},
			[]string{},
			&condition{
				c: &SCrypt{
					saltLen: 32,
					n:       32768,
					r:       8,
					p:       1,
					keyLen:  32,
				},
				passwd:     []byte("password"),
				passwdRepl: []byte("replaced password"),
				saltReader: bytes.NewReader([]byte("12345678901234567890123456789012")),
			},
			&action{
				match: false,
			},
		),
		gen(
			"read error",
			[]string{},
			[]string{},
			&condition{
				c: &SCrypt{
					saltLen: 32,
					n:       32768,
					r:       8,
					p:       1,
					keyLen:  32,
				},
				passwd:     []byte("password"),
				saltReader: &testutil.ErrorReader{},
			},
			&action{
				match: false,
				err: &er.Error{
					Package:     ErrPkg,
					Type:        ErrTypeSCrypt,
					Description: ErrDscHash,
				},
			},
		),
		gen(
			"invalid param",
			[]string{},
			[]string{},
			&condition{
				c: &SCrypt{
					saltLen: 32,
					n:       12345, // Invalid
					r:       8,
					p:       1,
					keyLen:  32,
				},
				passwd:     []byte("password"),
				saltReader: bytes.NewReader([]byte("12345678901234567890123456789012")),
			},
			&action{
				match: false,
				err: &er.Error{
					Package:     ErrPkg,
					Type:        ErrTypeSCrypt,
					Description: ErrDscHash,
				},
			},
		),
		gen(
			"replace",
			[]string{},
			[]string{},
			&condition{
				c: &SCrypt{
					saltLen: 32,
					n:       32768,
					r:       8,
					p:       1,
					keyLen:  32,
				},
				passwd:     []byte("password"),
				saltReader: bytes.NewReader([]byte("12345678901234567890123456789012")),
				cryptRepl:  &SCrypt{},
			},
			&action{
				match: false,
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			tmp := rand.Reader
			rand.Reader = tt.C().saltReader
			defer func() {
				rand.Reader = tmp
			}()

			got, err := tt.C().c.Hash(tt.C().passwd)
			testutil.Diff(t, tt.A().err, err, cmpopts.EquateErrors())

			if len(tt.C().passwdRepl) > 0 {
				tt.C().passwd = tt.C().passwdRepl
			}
			if tt.C().cryptRepl != nil {
				tt.C().c = tt.C().cryptRepl
			}
			err = tt.C().c.Compare(got, tt.C().passwd)
			testutil.Diff(t, tt.A().match, err == nil)
		})
	}
}

func TestNewPBKDF2(t *testing.T) {
	type condition struct {
		spec *k.PBKDF2Spec
	}

	type action struct {
		c   *PBKDF2
		err error
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"empty spec",
			[]string{},
			[]string{},
			&condition{
				spec: &k.PBKDF2Spec{},
			},
			&action{
				c: &PBKDF2{
					saltLen:  32,
					iter:     4096,
					keyLen:   32,
					hashFunc: sha256.New,
				},
			},
		),
		gen(
			"full spec",
			[]string{},
			[]string{},
			&condition{
				spec: &k.PBKDF2Spec{
					SaltLen: 16,
					Iter:    1024,
					KeyLen:  8,
					HashAlg: k.HashAlg_SHA1,
				},
			},
			&action{
				c: &PBKDF2{
					saltLen:  16,
					iter:     1024,
					keyLen:   8,
					hashFunc: sha1.New,
				},
			},
		),
		gen(
			"unsupported hash",
			[]string{},
			[]string{},
			&condition{
				spec: &k.PBKDF2Spec{HashAlg: k.HashAlg_FNV1_32},
			},
			&action{
				c: nil,
				err: &er.Error{
					Package:     ErrPkg,
					Type:        ErrTypePBKDF2,
					Description: ErrDscHashValid,
				},
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			c, err := NewPBKDF2(tt.C().spec)
			testutil.Diff(t, tt.A().err, err, cmpopts.EquateErrors())
			testutil.Diff(t, tt.A().c, c, cmp.AllowUnexported(PBKDF2{}), cmp.Comparer(testutil.ComparePointer[func() hash.Hash]))
		})
	}
}

func TestPBKDF2(t *testing.T) {
	type condition struct {
		c          *PBKDF2
		saltReader io.Reader
		passwd     []byte
		passwdRepl []byte  // Replaced before compare
		cryptRepl  *PBKDF2 // Replaced before compare
	}

	type action struct {
		err   error
		match bool
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"nil password",
			[]string{},
			[]string{},
			&condition{
				c: &PBKDF2{
					saltLen:  32,
					iter:     10,
					keyLen:   32,
					hashFunc: sha256.New,
				},
				passwd:     nil,
				saltReader: bytes.NewReader([]byte("12345678901234567890123456789012")),
			},
			&action{
				match: true,
			},
		),
		gen(
			"non-nil password",
			[]string{},
			[]string{},
			&condition{
				c: &PBKDF2{
					saltLen:  32,
					iter:     10,
					keyLen:   32,
					hashFunc: sha256.New,
				},
				passwd:     []byte("password"),
				saltReader: bytes.NewReader([]byte("12345678901234567890123456789012")),
			},
			&action{
				match: true,
			},
		),
		gen(
			"password not match",
			[]string{},
			[]string{},
			&condition{
				c: &PBKDF2{
					saltLen:  32,
					iter:     10,
					keyLen:   32,
					hashFunc: sha256.New,
				},
				passwd:     []byte("password"),
				passwdRepl: []byte("replaced password"),
				saltReader: bytes.NewReader([]byte("12345678901234567890123456789012")),
			},
			&action{
				match: false,
			},
		),
		gen(
			"read error",
			[]string{},
			[]string{},
			&condition{
				c: &PBKDF2{
					saltLen:  32,
					iter:     10,
					keyLen:   32,
					hashFunc: sha256.New,
				},
				passwd:     []byte("password"),
				saltReader: &testutil.ErrorReader{},
			},
			&action{
				match: false,
				err: &er.Error{
					Package:     ErrPkg,
					Type:        ErrTypePBKDF2,
					Description: ErrDscHash,
				},
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			tmp := rand.Reader
			rand.Reader = tt.C().saltReader
			defer func() {
				rand.Reader = tmp
			}()

			got, err := tt.C().c.Hash(tt.C().passwd)
			testutil.Diff(t, tt.A().err, err, cmpopts.EquateErrors())

			if len(tt.C().passwdRepl) > 0 {
				tt.C().passwd = tt.C().passwdRepl
			}
			if tt.C().cryptRepl != nil {
				tt.C().c = tt.C().cryptRepl
			}
			err = tt.C().c.Compare(got, tt.C().passwd)
			testutil.Diff(t, tt.A().match, err == nil)
		})
	}
}

func TestNewArgon2i(t *testing.T) {
	type condition struct {
		spec *k.Argon2Spec
	}

	type action struct {
		c   *Argon2i
		err error
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"empty spec",
			[]string{},
			[]string{},
			&condition{
				spec: &k.Argon2Spec{},
			},
			&action{
				c: &Argon2i{
					saltLen: 32,
					time:    3,
					memory:  32 * 1024,
					threads: 4,
					keyLen:  32,
				},
			},
		),
		gen(
			"full spec",
			[]string{},
			[]string{},
			&condition{
				spec: &k.Argon2Spec{
					SaltLen: 16,
					Time:    5,
					Memory:  16 * 1024,
					Threads: 2,
					KeyLen:  8,
				},
			},
			&action{
				c: &Argon2i{
					saltLen: 16,
					time:    5,
					memory:  16 * 1024,
					threads: 2,
					keyLen:  8,
				},
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			c, err := NewArgon2i(tt.C().spec)
			testutil.Diff(t, tt.A().err, err, cmpopts.EquateErrors())
			testutil.Diff(t, tt.A().c, c, cmp.AllowUnexported(Argon2i{}))
		})
	}
}

func TestArgon2i(t *testing.T) {
	type condition struct {
		c          *Argon2i
		saltReader io.Reader
		passwd     []byte
		passwdRepl []byte // Replaced before compare
	}

	type action struct {
		err   error
		match bool
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"nil password",
			[]string{},
			[]string{},
			&condition{
				c: &Argon2i{
					saltLen: 32,
					time:    1,
					memory:  1024,
					threads: 1,
					keyLen:  32,
				},
				passwd:     nil,
				saltReader: bytes.NewReader([]byte("12345678901234567890123456789012")),
			},
			&action{
				match: true,
			},
		),
		gen(
			"non-nil password",
			[]string{},
			[]string{},
			&condition{
				c: &Argon2i{
					saltLen: 32,
					time:    1,
					memory:  1024,
					threads: 1,
					keyLen:  32,
				},
				passwd:     []byte("password"),
				saltReader: bytes.NewReader([]byte("12345678901234567890123456789012")),
			},
			&action{
				match: true,
			},
		),
		gen(
			"password not match",
			[]string{},
			[]string{},
			&condition{
				c: &Argon2i{
					saltLen: 32,
					time:    1,
					memory:  1024,
					threads: 1,
					keyLen:  32,
				},
				passwd:     []byte("password"),
				passwdRepl: []byte("replaced password"),
				saltReader: bytes.NewReader([]byte("12345678901234567890123456789012")),
			},
			&action{
				match: false,
			},
		),
		gen(
			"read error",
			[]string{},
			[]string{},
			&condition{
				c: &Argon2i{
					saltLen: 32,
					time:    1,
					memory:  1024,
					threads: 1,
					keyLen:  32,
				},
				passwd:     []byte("password"),
				saltReader: &testutil.ErrorReader{},
			},
			&action{
				match: false,
				err: &er.Error{
					Package:     ErrPkg,
					Type:        ErrTypeArgon2i,
					Description: ErrDscHash,
				},
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			tmp := rand.Reader
			rand.Reader = tt.C().saltReader
			defer func() {
				rand.Reader = tmp
			}()

			got, err := tt.C().c.Hash(tt.C().passwd)
			testutil.Diff(t, tt.A().err, err, cmpopts.EquateErrors())

			if len(tt.C().passwdRepl) > 0 {
				tt.C().passwd = tt.C().passwdRepl
			}
			err = tt.C().c.Compare(got, tt.C().passwd)
			testutil.Diff(t, tt.A().match, err == nil)
		})
	}
}

func TestNewArgon2id(t *testing.T) {
	type condition struct {
		spec *k.Argon2Spec
	}

	type action struct {
		c   *Argon2id
		err error
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"empty spec",
			[]string{},
			[]string{},
			&condition{
				spec: &k.Argon2Spec{},
			},
			&action{
				c: &Argon2id{
					saltLen: 32,
					time:    1,
					memory:  64 * 1024,
					threads: 4,
					keyLen:  32,
				},
			},
		),
		gen(
			"full spec",
			[]string{},
			[]string{},
			&condition{
				spec: &k.Argon2Spec{
					SaltLen: 16,
					Time:    5,
					Memory:  16 * 1024,
					Threads: 2,
					KeyLen:  8,
				},
			},
			&action{
				c: &Argon2id{
					saltLen: 16,
					time:    5,
					memory:  16 * 1024,
					threads: 2,
					keyLen:  8,
				},
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			c, err := NewArgon2id(tt.C().spec)
			testutil.Diff(t, tt.A().err, err, cmpopts.EquateErrors())
			testutil.Diff(t, tt.A().c, c, cmp.AllowUnexported(Argon2id{}))
		})
	}
}

func TestArgon2id(t *testing.T) {
	type condition struct {
		c          *Argon2id
		saltReader io.Reader
		passwd     []byte
		passwdRepl []byte // Replaced before compare
	}

	type action struct {
		err   error
		match bool
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"nil password",
			[]string{},
			[]string{},
			&condition{
				c: &Argon2id{
					saltLen: 32,
					time:    1,
					memory:  1024,
					threads: 1,
					keyLen:  32,
				},
				passwd:     nil,
				saltReader: bytes.NewReader([]byte("12345678901234567890123456789012")),
			},
			&action{
				match: true,
			},
		),
		gen(
			"non-nil password",
			[]string{},
			[]string{},
			&condition{
				c: &Argon2id{
					saltLen: 32,
					time:    1,
					memory:  1024,
					threads: 1,
					keyLen:  32,
				},
				passwd:     []byte("password"),
				saltReader: bytes.NewReader([]byte("12345678901234567890123456789012")),
			},
			&action{
				match: true,
			},
		),
		gen(
			"password not match",
			[]string{},
			[]string{},
			&condition{
				c: &Argon2id{
					saltLen: 32,
					time:    1,
					memory:  1024,
					threads: 1,
					keyLen:  32,
				},
				passwd:     []byte("password"),
				passwdRepl: []byte("replaced password"),
				saltReader: bytes.NewReader([]byte("12345678901234567890123456789012")),
			},
			&action{
				match: false,
			},
		),
		gen(
			"read error",
			[]string{},
			[]string{},
			&condition{
				c: &Argon2id{
					saltLen: 32,
					time:    1,
					memory:  1024,
					threads: 1,
					keyLen:  32,
				},
				passwd:     []byte("password"),
				saltReader: &testutil.ErrorReader{},
			},
			&action{
				match: false,
				err: &er.Error{
					Package:     ErrPkg,
					Type:        ErrTypeArgon2id,
					Description: ErrDscHash,
				},
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			tmp := rand.Reader
			rand.Reader = tt.C().saltReader
			defer func() {
				rand.Reader = tmp
			}()

			got, err := tt.C().c.Hash(tt.C().passwd)
			testutil.Diff(t, tt.A().err, err, cmpopts.EquateErrors())

			if len(tt.C().passwdRepl) > 0 {
				tt.C().passwd = tt.C().passwdRepl
			}
			err = tt.C().c.Compare(got, tt.C().passwd)
			testutil.Diff(t, tt.A().match, err == nil)
		})
	}
}
