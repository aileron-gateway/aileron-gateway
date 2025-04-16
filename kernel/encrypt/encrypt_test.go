// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package encrypt_test

import (
	"testing"

	k "github.com/aileron-gateway/aileron-gateway/apis/kernel"
	"github.com/aileron-gateway/aileron-gateway/kernel/encrypt"
	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
	"github.com/google/go-cmp/cmp"
)

func TestEncrypterFromType(t *testing.T) {
	type condition struct {
		typ k.CommonKeyCryptType
	}

	type action struct {
		enc encrypt.EncryptFunc
	}

	cndInvalidType := "invalid encryption type"
	actCheckNil := "check nil"

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	tb.Condition(cndInvalidType, "input invalid encryption type")
	tb.Action(actCheckNil, "check that nil was returned")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"AESGCM",
			[]string{},
			[]string{},
			&condition{
				typ: k.CommonKeyCryptType_AESGCM,
			},
			&action{
				enc: encrypt.EncryptAESGCM,
			},
		),
		gen(
			"AESCBC",
			[]string{},
			[]string{},
			&condition{
				typ: k.CommonKeyCryptType_AESCBC,
			},
			&action{
				enc: encrypt.EncryptAESCBC,
			},
		),
		gen(
			"AESCFB",
			[]string{},
			[]string{},
			&condition{
				typ: k.CommonKeyCryptType_AESCFB,
			},
			&action{
				enc: encrypt.EncryptAESCFB,
			},
		),
		gen(
			"AESCTR",
			[]string{},
			[]string{},
			&condition{
				typ: k.CommonKeyCryptType_AESCTR,
			},
			&action{
				enc: encrypt.EncryptAESCTR,
			},
		),
		gen(
			"AESOFB",
			[]string{},
			[]string{},
			&condition{
				typ: k.CommonKeyCryptType_AESOFB,
			},
			&action{
				enc: encrypt.EncryptAESOFB,
			},
		),
		gen(
			"DESCBC",
			[]string{},
			[]string{},
			&condition{
				typ: k.CommonKeyCryptType_DESCBC,
			},
			&action{
				enc: encrypt.EncryptDESCBC,
			},
		),
		gen(
			"DESCFB",
			[]string{},
			[]string{},
			&condition{
				typ: k.CommonKeyCryptType_DESCFB,
			},
			&action{
				enc: encrypt.EncryptDESCFB,
			},
		),
		gen(
			"DESCTR",
			[]string{},
			[]string{},
			&condition{
				typ: k.CommonKeyCryptType_DESCTR,
			},
			&action{
				enc: encrypt.EncryptDESCTR,
			},
		),
		gen(
			"DESOFB",
			[]string{},
			[]string{},
			&condition{
				typ: k.CommonKeyCryptType_DESOFB,
			},
			&action{
				enc: encrypt.EncryptDESOFB,
			},
		),
		gen(
			"TripleDESCBC",
			[]string{},
			[]string{},
			&condition{
				typ: k.CommonKeyCryptType_TripleDESCBC,
			},
			&action{
				enc: encrypt.EncryptTripleDESCBC,
			},
		),
		gen(
			"TripleDESCFB",
			[]string{},
			[]string{},
			&condition{
				typ: k.CommonKeyCryptType_TripleDESCFB,
			},
			&action{
				enc: encrypt.EncryptTripleDESCFB,
			},
		),
		gen(
			"TripleDESCTR",
			[]string{},
			[]string{},
			&condition{
				typ: k.CommonKeyCryptType_TripleDESCTR,
			},
			&action{
				enc: encrypt.EncryptTripleDESCTR,
			},
		),
		gen(
			"TripleDESOFB",
			[]string{},
			[]string{},
			&condition{
				typ: k.CommonKeyCryptType_TripleDESOFB,
			},
			&action{
				enc: encrypt.EncryptTripleDESOFB,
			},
		),
		gen(
			"RC4",
			[]string{},
			[]string{},
			&condition{
				typ: k.CommonKeyCryptType_RC4,
			},
			&action{
				enc: encrypt.EncryptRC4,
			},
		),
		gen(
			"INVALID",
			[]string{cndInvalidType},
			[]string{actCheckNil},
			&condition{
				typ: 99999999,
			},
			&action{
				enc: nil,
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			enc := encrypt.EncrypterFromType(tt.C().typ)
			testutil.Diff(t, tt.A().enc, enc, cmp.Comparer(testutil.ComparePointer[encrypt.EncryptFunc]))
		})
	}
}

func TestDecrypterFromType(t *testing.T) {
	type condition struct {
		typ k.CommonKeyCryptType
	}

	type action struct {
		dec encrypt.DecryptFunc
	}

	cndInvalidType := "invalid encryption type"
	actCheckNil := "check nil"

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	tb.Condition(cndInvalidType, "input invalid encryption type")
	tb.Action(actCheckNil, "check that nil was returned")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"AESGCM",
			[]string{},
			[]string{},
			&condition{
				typ: k.CommonKeyCryptType_AESGCM,
			},
			&action{
				dec: encrypt.DecryptAESGCM,
			},
		),
		gen(
			"AESCBC",
			[]string{},
			[]string{},
			&condition{
				typ: k.CommonKeyCryptType_AESCBC,
			},
			&action{
				dec: encrypt.DecryptAESCBC,
			},
		),
		gen(
			"AESCFB",
			[]string{},
			[]string{},
			&condition{
				typ: k.CommonKeyCryptType_AESCFB,
			},
			&action{
				dec: encrypt.DecryptAESCFB,
			},
		),
		gen(
			"AESCTR",
			[]string{},
			[]string{},
			&condition{
				typ: k.CommonKeyCryptType_AESCTR,
			},
			&action{
				dec: encrypt.DecryptAESCTR,
			},
		),
		gen(
			"AESOFB",
			[]string{},
			[]string{},
			&condition{
				typ: k.CommonKeyCryptType_AESOFB,
			},
			&action{
				dec: encrypt.DecryptAESOFB,
			},
		),
		gen(
			"DESCBC",
			[]string{},
			[]string{},
			&condition{
				typ: k.CommonKeyCryptType_DESCBC,
			},
			&action{
				dec: encrypt.DecryptDESCBC,
			},
		),
		gen(
			"DESCFB",
			[]string{},
			[]string{},
			&condition{
				typ: k.CommonKeyCryptType_DESCFB,
			},
			&action{
				dec: encrypt.DecryptDESCFB,
			},
		),
		gen(
			"DESCTR",
			[]string{},
			[]string{},
			&condition{
				typ: k.CommonKeyCryptType_DESCTR,
			},
			&action{
				dec: encrypt.DecryptDESCTR,
			},
		),
		gen(
			"DESOFB",
			[]string{},
			[]string{},
			&condition{
				typ: k.CommonKeyCryptType_DESOFB,
			},
			&action{
				dec: encrypt.DecryptDESOFB,
			},
		),
		gen(
			"TripleDESCBC",
			[]string{},
			[]string{},
			&condition{
				typ: k.CommonKeyCryptType_TripleDESCBC,
			},
			&action{
				dec: encrypt.DecryptTripleDESCBC,
			},
		),
		gen(
			"TripleDESCFB",
			[]string{},
			[]string{},
			&condition{
				typ: k.CommonKeyCryptType_TripleDESCFB,
			},
			&action{
				dec: encrypt.DecryptTripleDESCFB,
			},
		),
		gen(
			"TripleDESCTR",
			[]string{},
			[]string{},
			&condition{
				typ: k.CommonKeyCryptType_TripleDESCTR,
			},
			&action{
				dec: encrypt.DecryptTripleDESCTR,
			},
		),
		gen(
			"TripleDESOFB",
			[]string{},
			[]string{},
			&condition{
				typ: k.CommonKeyCryptType_TripleDESOFB,
			},
			&action{
				dec: encrypt.DecryptTripleDESOFB,
			},
		),
		gen(
			"RC4",
			[]string{},
			[]string{},
			&condition{
				typ: k.CommonKeyCryptType_RC4,
			},
			&action{
				dec: encrypt.DecryptRC4,
			},
		),
		gen(
			"INVALID",
			[]string{cndInvalidType},
			[]string{actCheckNil},
			&condition{
				typ: 99999999,
			},
			&action{
				dec: nil,
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			dec := encrypt.DecrypterFromType(tt.C().typ)
			testutil.Diff(t, tt.A().dec, dec, cmp.Comparer(testutil.ComparePointer[encrypt.DecryptFunc]))
		})
	}
}
