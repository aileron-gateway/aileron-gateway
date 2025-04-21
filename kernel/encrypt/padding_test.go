// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package encrypt_test

import (
	"bytes"
	"crypto/rand"
	"io"
	"testing"

	"github.com/aileron-gateway/aileron-gateway/kernel/encrypt"
	"github.com/aileron-gateway/aileron-gateway/kernel/er"
	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
	"github.com/google/go-cmp/cmp/cmpopts"
)

// PaddingFunc is the type of function that add padding
// to the given data.
// The length of the padded data will be a multiple of blockSize.
type PaddingFunc func(blockSize int, data []byte) (padded []byte, err error)

// UnPaddingFunc is the type of function that remove padding
// from to the given data.
// The length of the given data should be a multiple of blockSize.
type UnPaddingFunc func(blockSize int, data []byte) (unpadded []byte, err error)

var (
	_ = PaddingFunc(encrypt.PKCS7Pad)        // Check that the function satisfies the signature.
	_ = PaddingFunc(encrypt.ISO7816Pad)      // Check that the function satisfies the signature.
	_ = PaddingFunc(encrypt.ISO10126Pad)     // Check that the function satisfies the signature.
	_ = UnPaddingFunc(encrypt.PKCS7UnPad)    // Check that the function satisfies the signature.
	_ = UnPaddingFunc(encrypt.ISO7816UnPad)  // Check that the function satisfies the signature.
	_ = UnPaddingFunc(encrypt.ISO10126UnPad) // Check that the function satisfies the signature.
)

func TestPKCS7Pad(t *testing.T) {
	type condition struct {
		inputNil  bool
		input     []byte
		blockSize int
	}

	type action struct {
		expect []byte
		err    error
	}

	CndInputNonNil := "non-nil input"
	CndInputNil := "nil input"
	CndBlockSize0 := "blockSize 0"
	CndBlockSize1 := "blockSize 1"
	CndBlockSize255 := "blockSize 255"
	CndBlockSize256 := "blockSize 256"
	ActCheckNoError := "check no error"
	ActCheckError := "check error"
	ActCheckData := "check data"

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	tb.Condition(CndInputNonNil, "input non-zero or non-nil bytes")
	tb.Condition(CndInputNil, "input nil as bytes")
	tb.Condition(CndBlockSize0, "input 0 as block size")
	tb.Condition(CndBlockSize1, "input 1 as block size")
	tb.Condition(CndBlockSize255, "input 255 as block size")
	tb.Condition(CndBlockSize256, "input 256 as block size")
	tb.Action(ActCheckNoError, "check that there is no error occurred")
	tb.Action(ActCheckError, "check that expected error occurred")
	tb.Action(ActCheckData, "check that the expected result is returned")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"nil and blockSize 0",
			[]string{CndInputNil, CndBlockSize0},
			[]string{ActCheckError, ActCheckData},
			&condition{
				inputNil:  true,
				input:     nil,
				blockSize: 0,
			},
			&action{
				expect: nil,
				err: &er.Error{
					Package:     encrypt.ErrPkg,
					Type:        encrypt.ErrTypePadding,
					Description: encrypt.ErrDscPadding,
				},
			},
		),
		gen(
			"nil and blockSize 1",
			[]string{CndInputNil, CndBlockSize1},
			[]string{ActCheckNoError, ActCheckData},
			&condition{
				inputNil:  true,
				input:     nil,
				blockSize: 1,
			},
			&action{
				expect: []byte{0x01},
			},
		),
		gen(
			"nil and blockSize 255",
			[]string{CndInputNil, CndBlockSize255},
			[]string{ActCheckNoError, ActCheckData},
			&condition{
				inputNil:  true,
				input:     nil,
				blockSize: 255,
			},
			&action{
				expect: []byte{
					0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
					0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
					0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
					0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
					0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
					0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
					0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
					0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
					0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
				},
			},
		),
		gen(
			"nil and blockSize 256",
			[]string{CndInputNil, CndBlockSize256},
			[]string{ActCheckError, ActCheckData},
			&condition{
				inputNil:  true,
				input:     nil,
				blockSize: 256,
			},
			&action{
				expect: nil,
				err: &er.Error{
					Package:     encrypt.ErrPkg,
					Type:        encrypt.ErrTypePadding,
					Description: encrypt.ErrDscPadding,
				},
			},
		),
		gen(
			"non-nil and blockSize 0",
			[]string{CndInputNonNil, CndBlockSize0},
			[]string{ActCheckError, ActCheckData},
			&condition{
				input:     []byte("abc"),
				blockSize: 0,
			},
			&action{
				expect: nil,
				err: &er.Error{
					Package:     encrypt.ErrPkg,
					Type:        encrypt.ErrTypePadding,
					Description: encrypt.ErrDscPadding,
				},
			},
		),
		gen(
			"non-nil and blockSize 1",
			[]string{CndInputNonNil, CndBlockSize1},
			[]string{ActCheckNoError, ActCheckData},
			&condition{
				input:     []byte("abc"),
				blockSize: 1,
			},
			&action{
				expect: []byte{0x61, 0x62, 0x63, 0x01},
			},
		),
		gen(
			"non-nil and blockSize 255",
			[]string{CndInputNonNil, CndBlockSize255},
			[]string{ActCheckNoError, ActCheckData},
			&condition{
				input:     []byte("abc"),
				blockSize: 255,
			},
			&action{
				expect: []byte{
					0x61, 0x62, 0x63, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc,
					0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc,
					0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc,
					0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc,
					0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc,
					0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc,
					0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc,
					0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc,
					0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc,
				},
			},
		),
		gen(
			"non-nil and blockSize 256",
			[]string{CndInputNonNil, CndBlockSize256},
			[]string{ActCheckError, ActCheckData},
			&condition{
				input:     []byte("abc"),
				blockSize: 256,
			},
			&action{
				expect: nil,
				err: &er.Error{
					Package:     encrypt.ErrPkg,
					Type:        encrypt.ErrTypePadding,
					Description: encrypt.ErrDscPadding,
				},
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			var out []byte
			var err error
			if tt.C().inputNil {
				out, err = encrypt.PKCS7Pad(tt.C().blockSize, nil)
			} else {
				out, err = encrypt.PKCS7Pad(tt.C().blockSize, tt.C().input)
			}

			testutil.Diff(t, tt.A().expect, out)
			testutil.Diff(t, tt.A().err, err, cmpopts.EquateErrors())
		})
	}
}

func TestPKCS7UnPad(t *testing.T) {
	type condition struct {
		inputNil  bool
		input     []byte
		blockSize int
	}

	type action struct {
		expect []byte
		err    error
	}

	CndInputNonNil := "non-nil input"
	CndInputNil := "nil input"
	CndBlockSize0 := "blockSize 0"
	CndBlockSize1 := "blockSize 1"
	CndBlockSize5 := "blockSize 5"
	CndBlockSize255 := "blockSize 255"
	CndBlockSize256 := "blockSize 256"
	CndDataSizeTooShort := "data size shorter than blockSize"
	CndDataSizeInvalid := "data size is invalid"
	CndInvalidPaddingSize := "padding size is invalid"
	ActCheckNoError := "check no error"
	ActCheckError := "check error"
	ActCheckData := "check data"

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	tb.Condition(CndInputNonNil, "input non-zero or non-nil bytes")
	tb.Condition(CndInputNil, "input nil as bytes")
	tb.Condition(CndBlockSize0, "input 0 as block size")
	tb.Condition(CndBlockSize1, "input 1 as block size")
	tb.Condition(CndBlockSize5, "input 5 as block size")
	tb.Condition(CndBlockSize255, "input 255 as block size")
	tb.Condition(CndBlockSize256, "input 256 as block size")
	tb.Condition(CndDataSizeTooShort, "data length is shorter than block size")
	tb.Condition(CndDataSizeInvalid, "data length is not multiple of block size")
	tb.Condition(CndInvalidPaddingSize, "padding size is invalid")
	tb.Action(ActCheckNoError, "check that there is no error occurred")
	tb.Action(ActCheckError, "check that expected error occurred")
	tb.Action(ActCheckData, "check that the expected result is returned")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"nil and blockSize 0",
			[]string{CndInputNil, CndBlockSize0},
			[]string{ActCheckError, ActCheckData},
			&condition{
				inputNil:  true,
				input:     nil,
				blockSize: 0,
			},
			&action{
				expect: nil,
				err: &er.Error{
					Package:     encrypt.ErrPkg,
					Type:        encrypt.ErrTypePadding,
					Description: encrypt.ErrDscUnpadding,
				},
			},
		),
		gen(
			"nil and blockSize 1",
			[]string{CndInputNil, CndBlockSize1, CndDataSizeTooShort},
			[]string{ActCheckError, ActCheckData},
			&condition{
				inputNil:  true,
				input:     nil,
				blockSize: 1,
			},
			&action{
				expect: nil,
				err: &er.Error{
					Package:     encrypt.ErrPkg,
					Type:        encrypt.ErrTypePadding,
					Description: encrypt.ErrDscUnpadding,
				},
			},
		),
		gen(
			"nil and blockSize 255",
			[]string{CndInputNil, CndBlockSize255, CndDataSizeTooShort},
			[]string{ActCheckError, ActCheckData},
			&condition{
				inputNil:  true,
				input:     nil,
				blockSize: 255,
			},
			&action{
				expect: nil,
				err: &er.Error{
					Package:     encrypt.ErrPkg,
					Type:        encrypt.ErrTypePadding,
					Description: encrypt.ErrDscUnpadding,
				},
			},
		),
		gen(
			"nil and blockSize 256",
			[]string{CndInputNil, CndBlockSize256},
			[]string{ActCheckError, ActCheckData},
			&condition{
				inputNil:  true,
				input:     nil,
				blockSize: 256,
			},
			&action{
				expect: nil,
				err: &er.Error{
					Package:     encrypt.ErrPkg,
					Type:        encrypt.ErrTypePadding,
					Description: encrypt.ErrDscUnpadding,
				},
			},
		),
		gen(
			"non-nil and blockSize 0",
			[]string{CndInputNonNil, CndBlockSize0},
			[]string{ActCheckError, ActCheckData},
			&condition{
				input:     []byte("abc"),
				blockSize: 0,
			},
			&action{
				expect: nil,
				err: &er.Error{
					Package:     encrypt.ErrPkg,
					Type:        encrypt.ErrTypePadding,
					Description: encrypt.ErrDscUnpadding,
				},
			},
		),
		gen(
			"non-nil and blockSize 1 with valid padding 1",
			[]string{CndInputNonNil, CndBlockSize1},
			[]string{ActCheckNoError, ActCheckData},
			&condition{
				input:     []byte{0x61, 0x62, 0x63, 0x01},
				blockSize: 1,
			},
			&action{
				expect: []byte{0x61, 0x62, 0x63},
			},
		),
		gen(
			"non-nil and blockSize 1 with valid padding 3",
			[]string{CndInputNonNil, CndBlockSize1},
			[]string{ActCheckNoError, ActCheckData},
			&condition{
				input:     []byte{0x61, 0x62, 0x63, 0x03},
				blockSize: 1,
			},
			&action{
				expect: []byte{0x61},
			},
		),
		gen(
			"non-nil and blockSize 1 with invalid data length",
			[]string{CndInputNonNil, CndBlockSize1, CndDataSizeTooShort},
			[]string{ActCheckError, ActCheckData},
			&condition{
				input:     []byte{},
				blockSize: 1,
			},
			&action{
				expect: nil,
				err: &er.Error{
					Package:     encrypt.ErrPkg,
					Type:        encrypt.ErrTypePadding,
					Description: encrypt.ErrDscUnpadding,
				},
			},
		),
		gen(
			"non-nil and blockSize 1 with invalid padding size",
			[]string{CndInputNonNil, CndBlockSize1, CndInvalidPaddingSize},
			[]string{ActCheckError, ActCheckData},
			&condition{
				input:     []byte{0x61, 0x62, 0x63, 0x05},
				blockSize: 1,
			},
			&action{
				expect: nil,
				err: &er.Error{
					Package:     encrypt.ErrPkg,
					Type:        encrypt.ErrTypePadding,
					Description: encrypt.ErrDscUnpadding,
				},
			},
		),
		gen(
			"non-nil and blockSize 5 with invalid data length",
			[]string{CndInputNonNil, CndBlockSize5, CndDataSizeInvalid},
			[]string{ActCheckError, ActCheckData},
			&condition{
				input:     []byte{0x61, 0x62, 0x63, 0x61, 0x62, 0x63, 0x02, 0x02},
				blockSize: 5,
			},
			&action{
				expect: nil,
				err: &er.Error{
					Package:     encrypt.ErrPkg,
					Type:        encrypt.ErrTypePadding,
					Description: encrypt.ErrDscUnpadding,
				},
			},
		),
		gen(
			"non-nil and blockSize 255",
			[]string{CndInputNonNil, CndBlockSize255},
			[]string{ActCheckNoError, ActCheckData},
			&condition{
				input: []byte{
					0x61, 0x62, 0x63, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc,
					0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc,
					0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc,
					0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc,
					0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc,
					0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc,
					0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc,
					0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc,
					0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc,
				},
				blockSize: 255,
			},
			&action{
				expect: []byte("abc"),
			},
		),
		gen(
			"non-nil and blockSize 255 with invalid data size",
			[]string{CndInputNonNil, CndBlockSize255, CndDataSizeTooShort},
			[]string{ActCheckNoError, ActCheckData},
			&condition{
				input: []byte{
					0x61, 0x62, 0x63,
				},
				blockSize: 255,
			},
			&action{
				expect: nil,
				err: &er.Error{
					Package:     encrypt.ErrPkg,
					Type:        encrypt.ErrTypePadding,
					Description: encrypt.ErrDscUnpadding,
				},
			},
		),
		gen(
			"non-nil and blockSize 256",
			[]string{CndInputNonNil, CndBlockSize256},
			[]string{ActCheckError, ActCheckData},
			&condition{
				input:     []byte("abc"),
				blockSize: 256,
			},
			&action{
				expect: nil,
				err: &er.Error{
					Package:     encrypt.ErrPkg,
					Type:        encrypt.ErrTypePadding,
					Description: encrypt.ErrDscUnpadding,
				},
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			var out []byte
			var err error
			if tt.C().inputNil {
				out, err = encrypt.PKCS7UnPad(tt.C().blockSize, nil)
			} else {
				out, err = encrypt.PKCS7UnPad(tt.C().blockSize, tt.C().input)
			}

			testutil.Diff(t, tt.A().expect, out)
			testutil.Diff(t, tt.A().err, err, cmpopts.EquateErrors())
		})
	}
}

func TestISO7816Pad(t *testing.T) {
	type condition struct {
		inputNil  bool
		input     []byte
		blockSize int
	}

	type action struct {
		expect []byte
		err    error
	}

	CndInputNonNil := "non-nil input"
	CndInputNil := "nil input"
	CndBlockSize0 := "blockSize 0"
	CndBlockSize1 := "blockSize 1"
	CndBlockSize255 := "blockSize 255"
	CndBlockSize256 := "blockSize 256"
	ActCheckNoError := "check no error"
	ActCheckError := "check error"
	ActCheckData := "check data"

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	tb.Condition(CndInputNonNil, "input non-zero or non-nil bytes")
	tb.Condition(CndInputNil, "input nil as bytes")
	tb.Condition(CndBlockSize0, "input 0 as block size")
	tb.Condition(CndBlockSize1, "input 1 as block size")
	tb.Condition(CndBlockSize255, "input 255 as block size")
	tb.Condition(CndBlockSize256, "input 256 as block size")
	tb.Action(ActCheckNoError, "check that there is no error occurred")
	tb.Action(ActCheckError, "check that expected error occurred")
	tb.Action(ActCheckData, "check that the expected result is returned")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"nil and blockSize 0",
			[]string{CndInputNil, CndBlockSize0},
			[]string{ActCheckError, ActCheckData},
			&condition{
				inputNil:  true,
				input:     nil,
				blockSize: 0,
			},
			&action{
				expect: nil,
				err: &er.Error{
					Package:     encrypt.ErrPkg,
					Type:        encrypt.ErrTypePadding,
					Description: encrypt.ErrDscPadding,
				},
			},
		),
		gen(
			"nil and blockSize 1",
			[]string{CndInputNil, CndBlockSize1},
			[]string{ActCheckNoError, ActCheckData},
			&condition{
				inputNil:  true,
				input:     nil,
				blockSize: 1,
			},
			&action{
				expect: []byte{0x80},
			},
		),
		gen(
			"nil and blockSize 255",
			[]string{CndInputNil, CndBlockSize255},
			[]string{ActCheckNoError, ActCheckData},
			&condition{
				inputNil:  true,
				input:     nil,
				blockSize: 255,
			},
			&action{
				expect: []byte{
					0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				},
			},
		),
		gen(
			"nil and blockSize 256",
			[]string{CndInputNil, CndBlockSize256},
			[]string{ActCheckError, ActCheckData},
			&condition{
				inputNil:  true,
				input:     nil,
				blockSize: 256,
			},
			&action{
				expect: nil,
				err: &er.Error{
					Package:     encrypt.ErrPkg,
					Type:        encrypt.ErrTypePadding,
					Description: encrypt.ErrDscPadding,
				},
			},
		),
		gen(
			"non-nil and blockSize 0",
			[]string{CndInputNonNil, CndBlockSize0},
			[]string{ActCheckError, ActCheckData},
			&condition{
				input:     []byte("abc"),
				blockSize: 0,
			},
			&action{
				expect: nil,
				err: &er.Error{
					Package:     encrypt.ErrPkg,
					Type:        encrypt.ErrTypePadding,
					Description: encrypt.ErrDscPadding,
				},
			},
		),
		gen(
			"non-nil and blockSize 1",
			[]string{CndInputNonNil, CndBlockSize1},
			[]string{ActCheckNoError, ActCheckData},
			&condition{
				input:     []byte("abc"),
				blockSize: 1,
			},
			&action{
				expect: []byte{0x61, 0x62, 0x63, 0x80},
			},
		),
		gen(
			"non-nil and blockSize 255",
			[]string{CndInputNonNil, CndBlockSize255},
			[]string{ActCheckNoError, ActCheckData},
			&condition{
				input:     []byte("abc"),
				blockSize: 255,
			},
			&action{
				expect: []byte{
					0x61, 0x62, 0x63, 0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				},
			},
		),
		gen(
			"non-nil and blockSize 256",
			[]string{CndInputNonNil, CndBlockSize256},
			[]string{ActCheckError, ActCheckData},
			&condition{
				input:     []byte("abc"),
				blockSize: 256,
			},
			&action{
				expect: nil,
				err: &er.Error{
					Package:     encrypt.ErrPkg,
					Type:        encrypt.ErrTypePadding,
					Description: encrypt.ErrDscPadding,
				},
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			var out []byte
			var err error
			if tt.C().inputNil {
				out, err = encrypt.ISO7816Pad(tt.C().blockSize, nil)
			} else {
				out, err = encrypt.ISO7816Pad(tt.C().blockSize, tt.C().input)
			}

			testutil.Diff(t, tt.A().expect, out)
			testutil.Diff(t, tt.A().err, err, cmpopts.EquateErrors())
		})
	}
}

func TestISO7816UnPad(t *testing.T) {
	type condition struct {
		inputNil  bool
		input     []byte
		blockSize int
	}

	type action struct {
		expect []byte
		err    error
	}

	CndInputNonNil := "non-nil input"
	CndInputNil := "nil input"
	CndBlockSize0 := "blockSize 0"
	CndBlockSize1 := "blockSize 1"
	CndBlockSize5 := "blockSize 5"
	CndBlockSize255 := "blockSize 255"
	CndBlockSize256 := "blockSize 256"
	CndDataSizeTooShort := "data size shorter than blockSize"
	CndDataSizeInvalid := "data size is invalid"
	CndInvalidPaddingSize := "padding size is invalid"
	ActCheckNoError := "check no error"
	ActCheckError := "check error"
	ActCheckData := "check data"

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	tb.Condition(CndInputNonNil, "input non-zero or non-nil bytes")
	tb.Condition(CndInputNil, "input nil as bytes")
	tb.Condition(CndBlockSize0, "input 0 as block size")
	tb.Condition(CndBlockSize1, "input 1 as block size")
	tb.Condition(CndBlockSize5, "input 5 as block size")
	tb.Condition(CndBlockSize255, "input 255 as block size")
	tb.Condition(CndBlockSize256, "input 256 as block size")
	tb.Condition(CndDataSizeTooShort, "data length is shorter than block size")
	tb.Condition(CndDataSizeInvalid, "data length is not multiple of block size")
	tb.Condition(CndInvalidPaddingSize, "padding size is invalid")
	tb.Action(ActCheckNoError, "check that there is no error occurred")
	tb.Action(ActCheckError, "check that expected error occurred")
	tb.Action(ActCheckData, "check that the expected result is returned")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"nil and blockSize 0",
			[]string{CndInputNil, CndBlockSize0},
			[]string{ActCheckError, ActCheckData},
			&condition{
				inputNil:  true,
				input:     nil,
				blockSize: 0,
			},
			&action{
				expect: nil,
				err: &er.Error{
					Package:     encrypt.ErrPkg,
					Type:        encrypt.ErrTypePadding,
					Description: encrypt.ErrDscUnpadding,
				},
			},
		),
		gen(
			"nil and blockSize 1",
			[]string{CndInputNil, CndBlockSize1, CndDataSizeTooShort},
			[]string{ActCheckError, ActCheckData},
			&condition{
				inputNil:  true,
				input:     nil,
				blockSize: 1,
			},
			&action{
				expect: nil,
				err: &er.Error{
					Package:     encrypt.ErrPkg,
					Type:        encrypt.ErrTypePadding,
					Description: encrypt.ErrDscUnpadding,
				},
			},
		),
		gen(
			"nil and blockSize 255",
			[]string{CndInputNil, CndBlockSize255, CndDataSizeTooShort},
			[]string{ActCheckError, ActCheckData},
			&condition{
				inputNil:  true,
				input:     nil,
				blockSize: 255,
			},
			&action{
				expect: nil,
				err: &er.Error{
					Package:     encrypt.ErrPkg,
					Type:        encrypt.ErrTypePadding,
					Description: encrypt.ErrDscUnpadding,
				},
			},
		),
		gen(
			"nil and blockSize 256",
			[]string{CndInputNil, CndBlockSize256},
			[]string{ActCheckError, ActCheckData},
			&condition{
				inputNil:  true,
				input:     nil,
				blockSize: 256,
			},
			&action{
				expect: nil,
				err: &er.Error{
					Package:     encrypt.ErrPkg,
					Type:        encrypt.ErrTypePadding,
					Description: encrypt.ErrDscUnpadding,
				},
			},
		),
		gen(
			"non-nil and blockSize 0",
			[]string{CndInputNonNil, CndBlockSize0},
			[]string{ActCheckError, ActCheckData},
			&condition{
				input:     []byte("abc"),
				blockSize: 0,
			},
			&action{
				expect: nil,
				err: &er.Error{
					Package:     encrypt.ErrPkg,
					Type:        encrypt.ErrTypePadding,
					Description: encrypt.ErrDscUnpadding,
				},
			},
		),
		gen(
			"non-nil and blockSize 1 with valid padding 1",
			[]string{CndInputNonNil, CndBlockSize1},
			[]string{ActCheckNoError, ActCheckData},
			&condition{
				input:     []byte{0x61, 0x62, 0x63, 0x80},
				blockSize: 1,
			},
			&action{
				expect: []byte{0x61, 0x62, 0x63},
			},
		),
		gen(
			"non-nil and blockSize 1 with invalid data length",
			[]string{CndInputNonNil, CndBlockSize1, CndDataSizeTooShort},
			[]string{ActCheckError, ActCheckData},
			&condition{
				input:     []byte{},
				blockSize: 1,
			},
			&action{
				expect: nil,
				err: &er.Error{
					Package:     encrypt.ErrPkg,
					Type:        encrypt.ErrTypePadding,
					Description: encrypt.ErrDscUnpadding,
				},
			},
		),
		gen(
			"non-nil and blockSize 1 with invalid padding size",
			[]string{CndInputNonNil, CndBlockSize1, CndInvalidPaddingSize},
			[]string{ActCheckError, ActCheckData},
			&condition{
				input:     []byte{0x61, 0x62, 0x63, 0x80, 0x00},
				blockSize: 1,
			},
			&action{
				expect: nil,
				err: &er.Error{
					Package:     encrypt.ErrPkg,
					Type:        encrypt.ErrTypePadding,
					Description: encrypt.ErrDscUnpadding,
				},
			},
		),
		gen(
			"non-nil and blockSize 5 with valid data length",
			[]string{CndInputNonNil, CndBlockSize5, CndDataSizeInvalid},
			[]string{ActCheckNoError, ActCheckData},
			&condition{
				input:     []byte{0x61, 0x62, 0x63, 0x61, 0x62, 0x63, 0x80, 0x00, 0x00, 0x00},
				blockSize: 5,
			},
			&action{
				expect: []byte{0x61, 0x62, 0x63, 0x61, 0x62, 0x63},
			},
		),
		gen(
			"non-nil and blockSize 5 with invalid data length",
			[]string{CndInputNonNil, CndBlockSize5, CndDataSizeInvalid},
			[]string{ActCheckError, ActCheckData},
			&condition{
				input:     []byte{0x61, 0x62, 0x63, 0x61, 0x62, 0x63, 0x80, 0x00},
				blockSize: 5,
			},
			&action{
				expect: nil,
				err: &er.Error{
					Package:     encrypt.ErrPkg,
					Type:        encrypt.ErrTypePadding,
					Description: encrypt.ErrDscUnpadding,
				},
			},
		),
		gen(
			"non-nil and blockSize 5 with invalid padding size",
			[]string{CndInputNonNil, CndBlockSize5, CndInvalidPaddingSize},
			[]string{ActCheckError, ActCheckData},
			&condition{
				input:     []byte{0x61, 0x62, 0x63, 0x00, 0x00},
				blockSize: 1,
			},
			&action{
				expect: nil,
				err: &er.Error{
					Package:     encrypt.ErrPkg,
					Type:        encrypt.ErrTypePadding,
					Description: encrypt.ErrDscUnpadding,
				},
			},
		),
		gen(
			"non-nil and blockSize 255",
			[]string{CndInputNonNil, CndBlockSize255},
			[]string{ActCheckNoError, ActCheckData},
			&condition{
				input: []byte{
					0x61, 0x62, 0x63, 0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				},
				blockSize: 255,
			},
			&action{
				expect: []byte("abc"),
			},
		),
		gen(
			"non-nil and blockSize 255 with invalid data size",
			[]string{CndInputNonNil, CndBlockSize255, CndDataSizeTooShort},
			[]string{ActCheckNoError, ActCheckData},
			&condition{
				input: []byte{
					0x61, 0x62, 0x63,
				},
				blockSize: 255,
			},
			&action{
				expect: nil,
				err: &er.Error{
					Package:     encrypt.ErrPkg,
					Type:        encrypt.ErrTypePadding,
					Description: encrypt.ErrDscUnpadding,
				},
			},
		),
		gen(
			"non-nil and blockSize 256",
			[]string{CndInputNonNil, CndBlockSize256},
			[]string{ActCheckError, ActCheckData},
			&condition{
				input:     []byte("abc"),
				blockSize: 256,
			},
			&action{
				expect: nil,
				err: &er.Error{
					Package:     encrypt.ErrPkg,
					Type:        encrypt.ErrTypePadding,
					Description: encrypt.ErrDscUnpadding,
				},
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			var out []byte
			var err error
			if tt.C().inputNil {
				out, err = encrypt.ISO7816UnPad(tt.C().blockSize, nil)
			} else {
				out, err = encrypt.ISO7816UnPad(tt.C().blockSize, tt.C().input)
			}

			testutil.Diff(t, tt.A().expect, out)
			testutil.Diff(t, tt.A().err, err, cmpopts.EquateErrors())
		})
	}
}

func TestISO10126Pad(t *testing.T) {
	type condition struct {
		inputNil  bool
		input     []byte
		blockSize int
		reader    io.Reader
	}

	type action struct {
		padSize int
		err     error
	}

	CndInputNonNil := "non-nil input"
	CndInputNil := "nil input"
	CndBlockSize0 := "blockSize 0"
	CndBlockSize1 := "blockSize 1"
	CndBlockSize255 := "blockSize 255"
	CndBlockSize256 := "blockSize 256"
	ActCheckNoError := "check no error"
	ActCheckError := "check error"
	ActCheckData := "check data"

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	tb.Condition(CndInputNonNil, "input non-zero or non-nil bytes")
	tb.Condition(CndInputNil, "input nil as bytes")
	tb.Condition(CndBlockSize0, "input 0 as block size")
	tb.Condition(CndBlockSize1, "input 1 as block size")
	tb.Condition(CndBlockSize255, "input 255 as block size")
	tb.Condition(CndBlockSize256, "input 256 as block size")
	tb.Action(ActCheckNoError, "check that there is no error occurred")
	tb.Action(ActCheckError, "check that expected error occurred")
	tb.Action(ActCheckData, "check that the expected result is returned")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"nil and blockSize 0",
			[]string{CndInputNil, CndBlockSize0},
			[]string{ActCheckError, ActCheckData},
			&condition{
				inputNil:  true,
				input:     nil,
				blockSize: 0,
			},
			&action{
				err: &er.Error{
					Package:     encrypt.ErrPkg,
					Type:        encrypt.ErrTypePadding,
					Description: encrypt.ErrDscPadding,
				},
			},
		),
		gen(
			"nil and blockSize 1",
			[]string{CndInputNil, CndBlockSize1},
			[]string{ActCheckNoError, ActCheckData},
			&condition{
				inputNil:  true,
				input:     nil,
				blockSize: 1,
			},
			&action{
				padSize: 1,
			},
		),
		gen(
			"nil and blockSize 255",
			[]string{CndInputNil, CndBlockSize255},
			[]string{ActCheckNoError, ActCheckData},
			&condition{
				inputNil:  true,
				input:     nil,
				blockSize: 255,
			},
			&action{
				padSize: 255,
			},
		),
		gen(
			"nil and blockSize 256",
			[]string{CndInputNil, CndBlockSize256},
			[]string{ActCheckError, ActCheckData},
			&condition{
				inputNil:  true,
				input:     nil,
				blockSize: 256,
			},
			&action{
				err: &er.Error{
					Package:     encrypt.ErrPkg,
					Type:        encrypt.ErrTypePadding,
					Description: encrypt.ErrDscPadding,
				},
			},
		),
		gen(
			"non-nil and blockSize 0",
			[]string{CndInputNonNil, CndBlockSize0},
			[]string{ActCheckError, ActCheckData},
			&condition{
				input:     []byte("abc"),
				blockSize: 0,
			},
			&action{
				err: &er.Error{
					Package:     encrypt.ErrPkg,
					Type:        encrypt.ErrTypePadding,
					Description: encrypt.ErrDscPadding,
				},
			},
		),
		gen(
			"non-nil and blockSize 1",
			[]string{CndInputNonNil, CndBlockSize1},
			[]string{ActCheckNoError, ActCheckData},
			&condition{
				input:     []byte("abc"),
				blockSize: 1,
			},
			&action{
				padSize: 1,
			},
		),
		gen(
			"non-nil and blockSize 255",
			[]string{CndInputNonNil, CndBlockSize255},
			[]string{ActCheckNoError, ActCheckData},
			&condition{
				input:     []byte("abc"),
				blockSize: 255,
			},
			&action{
				padSize: 252,
			},
		),
		gen(
			"non-nil and blockSize 256",
			[]string{CndInputNonNil, CndBlockSize256},
			[]string{ActCheckError, ActCheckData},
			&condition{
				input:     []byte("abc"),
				blockSize: 256,
			},
			&action{
				err: &er.Error{
					Package:     encrypt.ErrPkg,
					Type:        encrypt.ErrTypePadding,
					Description: encrypt.ErrDscPadding,
				},
			},
		),
		gen(
			"error reading random values",
			[]string{CndInputNonNil, CndBlockSize256},
			[]string{ActCheckError, ActCheckData},
			&condition{
				input:     []byte("abc"),
				blockSize: 255,
				reader:    &testutil.ErrorReader{},
			},
			&action{
				err: &er.Error{
					Package:     encrypt.ErrPkg,
					Type:        encrypt.ErrTypePadding,
					Description: encrypt.ErrDscPadding,
				},
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			if tt.C().reader != nil {
				tmp := rand.Reader
				rand.Reader = tt.C().reader
				defer func() {
					rand.Reader = tmp
				}()
			}

			var out []byte
			var err error
			if tt.C().inputNil {
				out, err = encrypt.ISO10126Pad(tt.C().blockSize, nil)
			} else {
				out, err = encrypt.ISO10126Pad(tt.C().blockSize, tt.C().input)
			}

			if tt.A().err != nil {
				testutil.Diff(t, []byte(nil), out)
			} else if !bytes.HasPrefix(out, tt.C().input) {
				t.Error("invalid prefix", out)
			}

			if tt.A().padSize != 0 && int(out[len(out)-1]) != tt.A().padSize {
				t.Error("invalid padding size", out)
			}

			testutil.Diff(t, tt.A().err, err, cmpopts.EquateErrors())
		})
	}
}

func TestISO10126UnPad(t *testing.T) {
	type condition struct {
		inputNil  bool
		input     []byte
		blockSize int
	}

	type action struct {
		expect []byte
		err    error
	}

	CndInputNonNil := "non-nil input"
	CndInputNil := "nil input"
	CndBlockSize0 := "blockSize 0"
	CndBlockSize1 := "blockSize 1"
	CndBlockSize5 := "blockSize 5"
	CndBlockSize255 := "blockSize 255"
	CndBlockSize256 := "blockSize 256"
	CndDataSizeTooShort := "data size shorter than blockSize"
	CndDataSizeInvalid := "data size is invalid"
	CndInvalidPaddingSize := "padding size is invalid"
	ActCheckNoError := "check no error"
	ActCheckError := "check error"
	ActCheckData := "check data"

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	tb.Condition(CndInputNonNil, "input non-zero or non-nil bytes")
	tb.Condition(CndInputNil, "input nil as bytes")
	tb.Condition(CndBlockSize0, "input 0 as block size")
	tb.Condition(CndBlockSize1, "input 1 as block size")
	tb.Condition(CndBlockSize5, "input 5 as block size")
	tb.Condition(CndBlockSize255, "input 255 as block size")
	tb.Condition(CndBlockSize256, "input 256 as block size")
	tb.Condition(CndDataSizeTooShort, "data length is shorter than block size")
	tb.Condition(CndDataSizeInvalid, "data length is not multiple of block size")
	tb.Condition(CndInvalidPaddingSize, "padding size is invalid")
	tb.Action(ActCheckNoError, "check that there is no error occurred")
	tb.Action(ActCheckError, "check that expected error occurred")
	tb.Action(ActCheckData, "check that the expected result is returned")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"nil and blockSize 0",
			[]string{CndInputNil, CndBlockSize0},
			[]string{ActCheckError, ActCheckData},
			&condition{
				inputNil:  true,
				input:     nil,
				blockSize: 0,
			},
			&action{
				expect: nil,
				err: &er.Error{
					Package:     encrypt.ErrPkg,
					Type:        encrypt.ErrTypePadding,
					Description: encrypt.ErrDscUnpadding,
				},
			},
		),
		gen(
			"nil and blockSize 1",
			[]string{CndInputNil, CndBlockSize1, CndDataSizeTooShort},
			[]string{ActCheckError, ActCheckData},
			&condition{
				inputNil:  true,
				input:     nil,
				blockSize: 1,
			},
			&action{
				expect: nil,
				err: &er.Error{
					Package:     encrypt.ErrPkg,
					Type:        encrypt.ErrTypePadding,
					Description: encrypt.ErrDscUnpadding,
				},
			},
		),
		gen(
			"nil and blockSize 255",
			[]string{CndInputNil, CndBlockSize255, CndDataSizeTooShort},
			[]string{ActCheckError, ActCheckData},
			&condition{
				inputNil:  true,
				input:     nil,
				blockSize: 255,
			},
			&action{
				expect: nil,
				err: &er.Error{
					Package:     encrypt.ErrPkg,
					Type:        encrypt.ErrTypePadding,
					Description: encrypt.ErrDscUnpadding,
				},
			},
		),
		gen(
			"nil and blockSize 256",
			[]string{CndInputNil, CndBlockSize256},
			[]string{ActCheckError, ActCheckData},
			&condition{
				inputNil:  true,
				input:     nil,
				blockSize: 256,
			},
			&action{
				expect: nil,
				err: &er.Error{
					Package:     encrypt.ErrPkg,
					Type:        encrypt.ErrTypePadding,
					Description: encrypt.ErrDscUnpadding,
				},
			},
		),
		gen(
			"non-nil and blockSize 0",
			[]string{CndInputNonNil, CndBlockSize0},
			[]string{ActCheckError, ActCheckData},
			&condition{
				input:     []byte("abc"),
				blockSize: 0,
			},
			&action{
				expect: nil,
				err: &er.Error{
					Package:     encrypt.ErrPkg,
					Type:        encrypt.ErrTypePadding,
					Description: encrypt.ErrDscUnpadding,
				},
			},
		),
		gen(
			"non-nil and blockSize 1 with valid padding 1",
			[]string{CndInputNonNil, CndBlockSize1},
			[]string{ActCheckNoError, ActCheckData},
			&condition{
				input:     []byte{0x61, 0x62, 0x63, 0x01},
				blockSize: 1,
			},
			&action{
				expect: []byte{0x61, 0x62, 0x63},
			},
		),
		gen(
			"non-nil and blockSize 1 with valid padding 3",
			[]string{CndInputNonNil, CndBlockSize1},
			[]string{ActCheckNoError, ActCheckData},
			&condition{
				input:     []byte{0x61, 0x62, 0x63, 0x03},
				blockSize: 1,
			},
			&action{
				expect: []byte{0x61},
			},
		),
		gen(
			"non-nil and blockSize 1 with invalid data length",
			[]string{CndInputNonNil, CndBlockSize1, CndDataSizeTooShort},
			[]string{ActCheckError, ActCheckData},
			&condition{
				input:     []byte{},
				blockSize: 1,
			},
			&action{
				expect: nil,
				err: &er.Error{
					Package:     encrypt.ErrPkg,
					Type:        encrypt.ErrTypePadding,
					Description: encrypt.ErrDscUnpadding,
				},
			},
		),
		gen(
			"non-nil and blockSize 1 with invalid padding size",
			[]string{CndInputNonNil, CndBlockSize1, CndInvalidPaddingSize},
			[]string{ActCheckError, ActCheckData},
			&condition{
				input:     []byte{0x61, 0x62, 0x63, 0x05},
				blockSize: 1,
			},
			&action{
				expect: nil,
				err: &er.Error{
					Package:     encrypt.ErrPkg,
					Type:        encrypt.ErrTypePadding,
					Description: encrypt.ErrDscUnpadding,
				},
			},
		),
		gen(
			"non-nil and blockSize 5 with invalid data length",
			[]string{CndInputNonNil, CndBlockSize5, CndDataSizeInvalid},
			[]string{ActCheckError, ActCheckData},
			&condition{
				input:     []byte{0x61, 0x62, 0x63, 0x61, 0x62, 0x63, 0x02, 0x02},
				blockSize: 5,
			},
			&action{
				expect: nil,
				err: &er.Error{
					Package:     encrypt.ErrPkg,
					Type:        encrypt.ErrTypePadding,
					Description: encrypt.ErrDscUnpadding,
				},
			},
		),
		gen(
			"non-nil and blockSize 255",
			[]string{CndInputNonNil, CndBlockSize255},
			[]string{ActCheckNoError, ActCheckData},
			&condition{
				input: []byte{
					0x61, 0x62, 0x63, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc,
					0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc,
					0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc,
					0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc,
					0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc,
					0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc,
					0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc,
					0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc,
					0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc, 0xfc,
				},
				blockSize: 255,
			},
			&action{
				expect: []byte("abc"),
			},
		),
		gen(
			"non-nil and blockSize 255 with invalid data size",
			[]string{CndInputNonNil, CndBlockSize255, CndDataSizeTooShort},
			[]string{ActCheckNoError, ActCheckData},
			&condition{
				input: []byte{
					0x61, 0x62, 0x63,
				},
				blockSize: 255,
			},
			&action{
				expect: nil,
				err: &er.Error{
					Package:     encrypt.ErrPkg,
					Type:        encrypt.ErrTypePadding,
					Description: encrypt.ErrDscUnpadding,
				},
			},
		),
		gen(
			"non-nil and blockSize 256",
			[]string{CndInputNonNil, CndBlockSize256},
			[]string{ActCheckError, ActCheckData},
			&condition{
				input:     []byte("abc"),
				blockSize: 256,
			},
			&action{
				expect: nil,
				err: &er.Error{
					Package:     encrypt.ErrPkg,
					Type:        encrypt.ErrTypePadding,
					Description: encrypt.ErrDscUnpadding,
				},
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			var out []byte
			var err error
			if tt.C().inputNil {
				out, err = encrypt.ISO10126UnPad(tt.C().blockSize, nil)
			} else {
				out, err = encrypt.ISO10126UnPad(tt.C().blockSize, tt.C().input)
			}

			testutil.Diff(t, tt.A().expect, out)
			testutil.Diff(t, tt.A().err, err, cmpopts.EquateErrors())
		})
	}
}
