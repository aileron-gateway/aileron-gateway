// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package encrypt

import (
	"bytes"
	"crypto/cipher"
	"crypto/des"
	"crypto/rand"
	"errors"
	"io"
	"reflect"
	"testing"

	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
)

var (
	_ = desFunc(newDES)       // Check that the function satisfies the signature.
	_ = desFunc(newTripleDES) // Check that the function satisfies the signature.
)

func TestNewDES(t *testing.T) {
	type condition struct {
		key    []byte
		reader io.Reader
	}

	type action struct {
		cipher cipher.Block
		iv     []byte
		err    error
	}

	cndValidKeyLength := "valid key length"
	cndErrReader := "use error reader"
	actCheckNoError := "check no error"
	actCheckError := "check error"

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	tb.Condition(cndValidKeyLength, "input non-zero or non-nil password")
	tb.Condition(cndErrReader, "error reading random bytes")
	tb.Action(actCheckNoError, "check that there is no error")
	tb.Action(actCheckError, "check that there the expected error")
	table := tb.Build()

	validCipher, _ := des.NewCipher([]byte("test_key"))
	reader := bytes.NewReader([]byte("123456789012345678901234"))

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"success to generate a new cipher",
			[]string{cndValidKeyLength},
			[]string{actCheckNoError},
			&condition{
				key:    []byte("test_key"),
				reader: reader,
			},
			&action{
				cipher: validCipher,
				iv:     []byte("12345678"),
				err:    nil,
			},
		),
		gen(
			"invalid key length",
			[]string{},
			[]string{actCheckError},
			&condition{
				key:    []byte("invalid_length_key"),
				reader: reader,
			},
			&action{
				cipher: nil,
				iv:     nil,
				err:    errors.New("crypto/des: invalid key size 18"),
			},
		),
		gen(
			"err generate iv",
			[]string{cndValidKeyLength, cndErrReader},
			[]string{actCheckError},
			&condition{
				key:    []byte("test_key"),
				reader: &testutil.ErrorReader{},
			},
			&action{
				cipher: nil,
				iv:     nil,
				err:    errors.New("rand read error"),
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			tmp := rand.Reader
			rand.Reader = tt.C().reader
			defer func() {
				rand.Reader = tmp
			}()

			cipher, iv, err := newDES(tt.C().key)
			if tt.A().err != nil {
				testutil.Diff(t, tt.A().err.Error(), err.Error())
			} else {
				testutil.Diff(t, nil, err)
			}
			testutil.Diff(t, true, reflect.DeepEqual(tt.A().cipher, cipher))
			testutil.Diff(t, tt.A().iv, iv)
		})
	}
}

func TestNewTripleDES(t *testing.T) {
	type condition struct {
		key    []byte
		reader io.Reader
	}

	type action struct {
		cipher cipher.Block
		iv     []byte
		err    error
	}

	cndValidKeyLength := "valid key length"
	cndErrReader := "use error reader"
	actCheckNoError := "check no error"
	actCheckError := "check error"

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	tb.Condition(cndValidKeyLength, "input non-zero or non-nil password")
	tb.Condition(cndErrReader, "error reading random bytes")
	tb.Action(actCheckNoError, "check that there is no error")
	tb.Action(actCheckError, "check that there the expected error")
	table := tb.Build()

	validCipher, _ := des.NewTripleDESCipher([]byte("testkey_testkey_testkey_"))
	reader := bytes.NewReader([]byte("123456789012345678901234"))

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"success to generate a new cipher",
			[]string{cndValidKeyLength},
			[]string{actCheckNoError},
			&condition{
				key:    []byte("testkey_testkey_testkey_"),
				reader: reader,
			},
			&action{
				cipher: validCipher,
				iv:     []byte("12345678"),
				err:    nil,
			},
		),
		gen(
			"invalid key length",
			[]string{},
			[]string{actCheckError},
			&condition{
				key:    []byte("invalid_length_key"),
				reader: reader,
			},
			&action{
				cipher: nil,
				iv:     nil,
				err:    errors.New("crypto/des: invalid key size 18"),
			},
		),
		gen(
			"err generate iv",
			[]string{cndValidKeyLength, cndErrReader},
			[]string{actCheckError},
			&condition{
				key:    []byte("testkey_testkey_testkey_"),
				reader: &testutil.ErrorReader{},
			},
			&action{
				cipher: nil,
				iv:     nil,
				err:    errors.New("rand read error"),
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			tmp := rand.Reader
			rand.Reader = tt.C().reader
			defer func() {
				rand.Reader = tmp
			}()

			cipher, iv, err := newTripleDES(tt.C().key)
			if tt.A().err != nil {
				testutil.Diff(t, tt.A().err.Error(), err.Error())
			} else {
				testutil.Diff(t, nil, err)
			}
			testutil.Diff(t, true, reflect.DeepEqual(tt.A().cipher, cipher))
			testutil.Diff(t, tt.A().iv, iv)
		})
	}
}
