// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package encrypt

import (
	"crypto/cipher"

	k "github.com/aileron-gateway/aileron-gateway/apis/kernel"
	"github.com/aileron-gateway/aileron-gateway/kernel/er"
)

// EncryptFunc is the type of function that encrypt the given plaintext.
type EncryptFunc func(key []byte, plaintext []byte) (ciphertext []byte, err error)

// DecryptFunc is the type of function that decrypt the given ciphertext.
type DecryptFunc func(key []byte, ciphertext []byte) (plaintext []byte, err error)

// EncrypterFromType returns an encryption function that corresponds to the given crypt type.
// This function returns nil when encrypter was not found.
func EncrypterFromType(t k.CommonKeyCryptType) EncryptFunc {
	typeToEnc := map[k.CommonKeyCryptType]EncryptFunc{
		k.CommonKeyCryptType_AESGCM: EncryptAESGCM,
		k.CommonKeyCryptType_AESCBC: EncryptAESCBC,
		k.CommonKeyCryptType_AESCFB: EncryptAESCFB,
		k.CommonKeyCryptType_AESCTR: EncryptAESCTR,
		k.CommonKeyCryptType_AESOFB: EncryptAESOFB,
	}
	return typeToEnc[t]
}

// DecrypterFromType returns a decryption function that corresponds to the given crypt type.
// This function returns nil when decrypter was not found.
func DecrypterFromType(t k.CommonKeyCryptType) DecryptFunc {
	typeToDec := map[k.CommonKeyCryptType]DecryptFunc{
		k.CommonKeyCryptType_AESGCM: DecryptAESGCM,
		k.CommonKeyCryptType_AESCBC: DecryptAESCBC,
		k.CommonKeyCryptType_AESCFB: DecryptAESCFB,
		k.CommonKeyCryptType_AESCTR: DecryptAESCTR,
		k.CommonKeyCryptType_AESOFB: DecryptAESOFB,
	}
	return typeToDec[t]
}

// blockEncrypt encrypt plaintext to ciphertext with given cipher and block mode.
// This function panics when the second argument c is nil.
func blockEncrypt(blockSize int, c cipher.BlockMode, plaintext []byte) ([]byte, error) {
	data, err := PKCS7Pad(blockSize, plaintext)
	if err != nil {
		return nil, (&er.Error{
			Package:     ErrPkg,
			Type:        ErrTypeBlock,
			Description: ErrDscEncrypt,
		}).Wrap(err)
	}

	dst := make([]byte, len(data))
	c.CryptBlocks(dst, data)

	return dst, nil
}

// blockDecrypt decrypts ciphertext to plaintext with given cipher and block mode.
// This function panics when the second argument c is nil.
func blockDecrypt(blockSize int, c cipher.BlockMode, ciphertext []byte) ([]byte, error) {
	n := len(ciphertext)
	if n < blockSize || n%blockSize != 0 {
		return nil, &er.Error{
			Package:     ErrPkg,
			Type:        ErrTypeBlock,
			Description: ErrDscDecrypt,
			Detail:      "invalid ciphertext length",
		}
	}

	plaintext := make([]byte, n)
	c.CryptBlocks(plaintext, ciphertext)

	unpadded, err := PKCS7UnPad(blockSize, plaintext)
	if err != nil {
		return nil, (&er.Error{
			Package:     ErrPkg,
			Type:        ErrTypeBlock,
			Description: ErrDscDecrypt,
		}).Wrap(err)
	}

	return unpadded, nil
}

// streamEncrypt encrypts plaintext with stream encryption with given cipher.
// This function panics when the first argument c is nil.
func streamEncrypt(c cipher.Stream, plaintext []byte) []byte {
	ciphertext := make([]byte, len(plaintext))
	c.XORKeyStream(ciphertext, plaintext)
	return ciphertext
}

// streamDecrypt decrypts ciphertext with stream encryption with given cipher.
// This function panics when the first argument c is nil.
func streamDecrypt(c cipher.Stream, ciphertext []byte) []byte {
	n := len(ciphertext)
	plaintext := make([]byte, n)
	c.XORKeyStream(plaintext, ciphertext)
	return plaintext
}
