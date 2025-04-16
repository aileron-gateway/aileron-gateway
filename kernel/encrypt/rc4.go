// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package encrypt

import (
	"crypto/rand"
	"crypto/rc4"
	"io"

	"github.com/aileron-gateway/aileron-gateway/kernel/er"
)

// EncryptRC4 encrypt plain text to cipher text using RC4 encryption cipher.
// RC4 is a common key stream encryption.
// The key argument should be the RC4 key, at least 1 byte and at most 256 bytes.
// See more details at https://pkg.go.dev/crypto/rc4
func EncryptRC4(key []byte, plaintext []byte) ([]byte, error) {
	// key size must be 1 - 256 bytes.
	c, err := rc4.NewCipher(key)
	if err != nil {
		return nil, (&er.Error{
			Package:     ErrPkg,
			Type:        ErrTypeRC4,
			Description: ErrDscEncrypt,
		}).Wrap(err)
	}

	iv := make([]byte, 24)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, (&er.Error{
			Package:     ErrPkg,
			Type:        ErrTypeRC4,
			Description: ErrDscEncrypt,
		}).Wrap(err)
	}

	return streamEncrypt(c, append(iv, plaintext...)), nil
}

// DecryptRC4 decrypt cipher text to plain text using RC4 cipher.
// RC4 is a common key stream encryption.
// The key argument should be the RC4 key,  at least 1 byte and at most 256 bytes.
// See more details at https://pkg.go.dev/crypto/rc4
func DecryptRC4(key []byte, ciphertext []byte) ([]byte, error) {
	// ciphertext should have at least 24 bytes which correspond to the iv length.
	if len(ciphertext) < 24 {
		return nil, &er.Error{
			Package:     ErrPkg,
			Type:        ErrTypeRC4,
			Description: ErrDscDecrypt,
			Detail:      "ciphertext is shorter than 24 bytes.",
		}
	}

	// key size must be 1 - 256 bytes.
	c, err := rc4.NewCipher(key)
	if err != nil {
		return nil, (&er.Error{
			Package:     ErrPkg,
			Type:        ErrTypeRC4,
			Description: ErrDscDecrypt,
		}).Wrap(err)
	}

	plaintext := streamDecrypt(c, ciphertext)

	return plaintext[24:], nil
}
