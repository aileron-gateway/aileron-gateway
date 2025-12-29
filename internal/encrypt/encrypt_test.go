// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package encrypt_test

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"testing"

	k "github.com/aileron-gateway/aileron-gateway/apis/kernel"
	"github.com/aileron-gateway/aileron-gateway/internal/encrypt"
	"github.com/aileron-gateway/aileron-gateway/internal/testutil"
	"github.com/aileron-projects/go/zcrypto/zaes"
	"github.com/google/go-cmp/cmp"
)

func TestEncrypterFromType(t *testing.T) {
	key := []byte("1234567890123456")
	plaintext := []byte("plaintext message")
	testCases := map[string]struct {
		typ        k.CommonKeyCryptType
		enc        encrypt.EncryptFunc
		ciphertext string
	}{
		"AESGCM": {
			typ:        k.CommonKeyCryptType_AESGCM,
			enc:        zaes.EncryptGCM,
			ciphertext: "66697865642076616c756520ad6e5f0773705d3964ddc5173946355bc680d98834bef3351f2b4eb25d36fae4f2",
		},
		"AESCBC": {
			typ:        k.CommonKeyCryptType_AESCBC,
			enc:        zaes.EncryptCBC,
			ciphertext: "66697865642076616c75652077696c6c152fece4213897d9a86bca729e6f6be4cce941d02efc6e0efa677c6ec2ba4b6a",
		},
		"AESCFB": {
			typ:        k.CommonKeyCryptType_AESCFB,
			enc:        zaes.EncryptCFB,
			ciphertext: "66697865642076616c75652077696c6c4c6e78d7a206d2db4dde7b34618dfae869",
		},
		"AESCTR": {
			typ:        k.CommonKeyCryptType_AESCTR,
			enc:        zaes.EncryptCTR,
			ciphertext: "66697865642076616c75652077696c6c4c6e78d7a206d2db4dde7b34618dfae85e",
		},
		"AESOFB": {
			typ:        k.CommonKeyCryptType_AESOFB,
			enc:        zaes.EncryptOFB,
			ciphertext: "66697865642076616c75652077696c6c4c6e78d7a206d2db4dde7b34618dfae87e",
		},
		"INVALID": {
			typ: 99999999,
			enc: nil,
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			enc := encrypt.EncrypterFromType(tc.typ)
			testutil.Diff(t, tc.enc, enc, cmp.Comparer(testutil.ComparePointer[encrypt.EncryptFunc]))
			if enc == nil {
				return
			}

			tmp := rand.Reader
			rand.Reader = bytes.NewReader([]byte("fixed value will be returned for reproducibility"))
			defer func() {
				rand.Reader = tmp
			}()
			ciphertext, err := enc(key, plaintext)
			testutil.Diff(t, nil, err)
			testutil.Diff(t, tc.ciphertext, hex.EncodeToString(ciphertext))
		})
	}
}

func TestDecrypterFromType(t *testing.T) {
	key := []byte("1234567890123456")
	plaintext := "plaintext message"
	testCases := map[string]struct {
		typ        k.CommonKeyCryptType
		dec        encrypt.DecryptFunc
		ciphertext string
	}{
		"AESGCM": {
			typ:        k.CommonKeyCryptType_AESGCM,
			dec:        zaes.DecryptGCM,
			ciphertext: "66697865642076616c756520ad6e5f0773705d3964ddc5173946355bc680d98834bef3351f2b4eb25d36fae4f2",
		},
		"AESCBC": {
			typ:        k.CommonKeyCryptType_AESCBC,
			dec:        zaes.DecryptCBC,
			ciphertext: "66697865642076616c75652077696c6c152fece4213897d9a86bca729e6f6be4cce941d02efc6e0efa677c6ec2ba4b6a",
		},
		"AESCFB": {
			typ:        k.CommonKeyCryptType_AESCFB,
			dec:        zaes.DecryptCFB,
			ciphertext: "66697865642076616c75652077696c6c4c6e78d7a206d2db4dde7b34618dfae869",
		},
		"AESCTR": {
			typ:        k.CommonKeyCryptType_AESCTR,
			dec:        zaes.DecryptCTR,
			ciphertext: "66697865642076616c75652077696c6c4c6e78d7a206d2db4dde7b34618dfae85e",
		},
		"AESOFB": {
			typ:        k.CommonKeyCryptType_AESOFB,
			dec:        zaes.DecryptOFB,
			ciphertext: "66697865642076616c75652077696c6c4c6e78d7a206d2db4dde7b34618dfae87e",
		},
		"INVALID": {
			typ: 99999999,
			dec: nil,
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			dec := encrypt.DecrypterFromType(tc.typ)
			testutil.Diff(t, tc.dec, dec, cmp.Comparer(testutil.ComparePointer[encrypt.DecryptFunc]))
			if dec == nil {
				return
			}
			ciphertext, _ := hex.DecodeString(tc.ciphertext)
			pt, err := dec(key, ciphertext)
			testutil.Diff(t, nil, err)
			testutil.Diff(t, plaintext, string(pt))
		})
	}
}
