// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package encrypt

import (
	k "github.com/aileron-gateway/aileron-gateway/apis/kernel"
	"github.com/aileron-projects/go/zcrypto/zaes"
)

// EncryptFunc is the type of function that encrypt the given plaintext.
type EncryptFunc func(key []byte, plaintext []byte) (ciphertext []byte, err error)

// DecryptFunc is the type of function that decrypt the given ciphertext.
type DecryptFunc func(key []byte, ciphertext []byte) (plaintext []byte, err error)

// EncrypterFromType returns an encryption function that corresponds to the given crypt type.
// This function returns nil when encrypter was not found.
func EncrypterFromType(t k.CommonKeyCryptType) EncryptFunc {
	typeToEnc := map[k.CommonKeyCryptType]EncryptFunc{
		k.CommonKeyCryptType_AESGCM: zaes.EncryptGCM,
		k.CommonKeyCryptType_AESCBC: zaes.EncryptCBC,
		k.CommonKeyCryptType_AESCFB: zaes.EncryptCFB,
		k.CommonKeyCryptType_AESCTR: zaes.EncryptCTR,
		k.CommonKeyCryptType_AESOFB: zaes.EncryptOFB,
	}
	return typeToEnc[t]
}

// DecrypterFromType returns a decryption function that corresponds to the given crypt type.
// This function returns nil when decrypter was not found.
func DecrypterFromType(t k.CommonKeyCryptType) DecryptFunc {
	typeToDec := map[k.CommonKeyCryptType]DecryptFunc{
		k.CommonKeyCryptType_AESGCM: zaes.DecryptGCM,
		k.CommonKeyCryptType_AESCBC: zaes.DecryptCBC,
		k.CommonKeyCryptType_AESCFB: zaes.DecryptCFB,
		k.CommonKeyCryptType_AESCTR: zaes.DecryptCTR,
		k.CommonKeyCryptType_AESOFB: zaes.DecryptOFB,
	}
	return typeToDec[t]
}
