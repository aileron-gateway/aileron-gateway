// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package encrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"testing"

	"github.com/aileron-gateway/aileron-gateway/kernel/er"
	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestBlockEncrypt(t *testing.T) {
	type condition struct {
		blockSize int
		plaintext []byte
	}

	type action struct {
		ciphertext string
		err        error
	}

	cndInvalidBlockSize := "invalid block size"
	actCheckNoError := "check no error"
	actCheckError := "check error"

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	tb.Condition(cndInvalidBlockSize, "input invalid block size")
	tb.Action(actCheckNoError, "check that there is no error")
	tb.Action(actCheckError, "check that there the expected error")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"success to encrypt",
			[]string{},
			[]string{actCheckNoError},
			&condition{
				blockSize: 16,
				plaintext: []byte("test"),
			},
			&action{
				ciphertext: "5a468e86a48ad37675a52b8f39d734df",
			},
		),
		gen(
			"invalid block size",
			[]string{cndInvalidBlockSize},
			[]string{actCheckError},
			&condition{
				blockSize: 999,
				plaintext: []byte("test"),
			},
			&action{
				ciphertext: "",
				err: &er.Error{
					Package:     ErrPkg,
					Type:        ErrTypeBlock,
					Description: ErrDscEncrypt,
				},
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			iv := []byte("1234567890123456")
			c, _ := aes.NewCipher([]byte("1234567890123456"))
			ciphertext, err := blockEncrypt(tt.C().blockSize, cipher.NewCBCEncrypter(c, iv), tt.C().plaintext)

			testutil.Diff(t, tt.A().err, err, cmpopts.EquateErrors())
			testutil.Diff(t, tt.A().ciphertext, hex.EncodeToString(ciphertext))
		})
	}
}

func TestBlockDecrypt(t *testing.T) {
	type condition struct {
		blockSize  int
		ciphertext string
	}

	type action struct {
		plaintext []byte
		err       error
	}

	cndInvalidBlockSize := "invalid block size"
	actCheckNoError := "check no error"
	actCheckError := "check error"

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	tb.Condition(cndInvalidBlockSize, "input invalid block size")
	tb.Action(actCheckNoError, "check that there is no error")
	tb.Action(actCheckError, "check that there the expected error")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"success to decrypt",
			[]string{},
			[]string{actCheckNoError},
			&condition{
				blockSize:  16,
				ciphertext: "5a468e86a48ad37675a52b8f39d734df",
			},
			&action{
				plaintext: []byte("test"),
			},
		),
		gen(
			"invalid block size",
			[]string{cndInvalidBlockSize},
			[]string{actCheckError},
			&condition{
				blockSize:  999,
				ciphertext: "5a468e86a48ad37675a52b8f39d734df",
			},
			&action{
				plaintext: nil,
				err: &er.Error{
					Package:     ErrPkg,
					Type:        ErrTypeBlock,
					Description: ErrDscDecrypt,
				},
			},
		),
		gen(
			"failed to unpadding",
			[]string{},
			[]string{actCheckError},
			&condition{
				blockSize:  16,
				ciphertext: "5a468e86a48ad37675a52b8f39d734ff",
			},
			&action{
				plaintext: nil,
				err: &er.Error{
					Package:     ErrPkg,
					Type:        ErrTypeBlock,
					Description: ErrDscDecrypt,
				},
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			iv := []byte("1234567890123456")
			c, _ := aes.NewCipher([]byte("1234567890123456"))
			ciphertext, _ := hex.DecodeString(tt.C().ciphertext)
			plaintext, err := blockDecrypt(tt.C().blockSize, cipher.NewCBCDecrypter(c, iv), ciphertext)

			testutil.Diff(t, tt.A().err, err, cmpopts.EquateErrors())
			testutil.Diff(t, tt.A().plaintext, plaintext)
		})
	}
}
