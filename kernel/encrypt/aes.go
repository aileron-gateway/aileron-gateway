// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package encrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"

	"github.com/aileron-gateway/aileron-gateway/kernel/er"
)

// newAES returns a new AES cipher block and initial vector.
// The key argument should be the AES key, either 16, 24, or 32 bytes
// to select AES-128, AES-192, or AES-256.
func newAES(key []byte) (cipher.Block, []byte, error) {
	c, err := aes.NewCipher(key)
	if err != nil {
		return nil, nil, err
	}

	iv := make([]byte, c.BlockSize())
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, nil, err
	}

	return c, iv, nil
}

// EncryptAESGCM encrypts plaintext with AES GCM cipher algorithm.
// The key argument should be the AES key, either 16, 24, or 32 bytes to select AES-128, AES-192, or AES-256.
//
// See more details at:
//   - https://pkg.go.dev/crypto/aes
//   - https://pkg.go.dev/crypto/cipher
func EncryptAESGCM(key []byte, plaintext []byte) ([]byte, error) {
	c, nonce, err := newAES(key)
	if err != nil {
		return nil, (&er.Error{
			Package:     ErrPkg,
			Type:        ErrTypeAES,
			Description: ErrDscEncrypt,
			Detail:      "aes gcm failed.",
		}).Wrap(err)
	}

	// NewGCM creates a new chipher with 12 bytes nonce.
	// err should always be nil.
	aead, err := cipher.NewGCM(c)
	if err != nil {
		return nil, (&er.Error{
			Package:     ErrPkg,
			Type:        ErrTypeAES,
			Description: ErrDscEncrypt,
			Detail:      "aes gcm failed.",
		}).Wrap(err)
	}

	// nonce should be an unique value.
	// There is no need to be a random value.
	// NewGCM uses standard nonce size of 12 bytes.
	// The nonce created by the newAES is always longer than 16 bytes.
	nonce = nonce[:aead.NonceSize()]

	// To append encrypted text to nonce slice,
	// set nonce variable as dst argument.
	// That means ciphertext = append(nonce, encrypt(plaintext)).
	ciphertext := aead.Seal(nonce, nonce, plaintext, nil)

	return ciphertext, nil
}

// DecryptAESGCM decrypts ciphertext with AES GCM cipher algorithm.
// The key argument should be the AES key, either 16, 24, or 32 bytes to select AES-128, AES-192, or AES-256.
//
// See more details at:
//   - https://pkg.go.dev/crypto/aes
//   - https://pkg.go.dev/crypto/cipher
func DecryptAESGCM(key []byte, ciphertext []byte) ([]byte, error) {
	c, err := aes.NewCipher(key)
	if err != nil {
		return nil, (&er.Error{
			Package:     ErrPkg,
			Type:        ErrTypeAES,
			Description: ErrDscDecrypt,
			Detail:      "aes gcm failed.",
		}).Wrap(err)
	}

	// NewGCM creates a new chipher with 12 bytes nonce.
	// err should always be nil.
	aead, err := cipher.NewGCM(c)
	if err != nil {
		return nil, (&er.Error{
			Package:     ErrPkg,
			Type:        ErrTypeAES,
			Description: ErrDscDecrypt,
			Detail:      "aes gcm failed.",
		}).Wrap(err)
	}

	n := len(ciphertext)
	if n < aead.NonceSize() {
		// ciphertext is shorter than length of nonce.
		return nil, &er.Error{
			Package:     ErrPkg,
			Type:        ErrTypeAES,
			Description: ErrDscDecrypt,
			Detail:      "aes gcm invalid ciphertext length.",
		}
	}

	nonce, ciphertext := ciphertext[:aead.NonceSize()], ciphertext[aead.NonceSize():]

	plaintext, err := aead.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, (&er.Error{
			Package:     ErrPkg,
			Type:        ErrTypeAES,
			Description: ErrDscDecrypt,
			Detail:      "aes gcm failed.",
		}).Wrap(err)
	}

	return plaintext, nil
}

// EncryptAESCBC encrypts plaintext with AES CBC cipher algorithm.
// The key argument should be the AES key, either 16, 24, or 32 bytes to select AES-128, AES-192, or AES-256.
//
// See more details at:
//   - https://pkg.go.dev/crypto/aes
//   - https://pkg.go.dev/crypto/cipher
func EncryptAESCBC(key []byte, plaintext []byte) ([]byte, error) {
	c, iv, err := newAES(key)
	if err != nil {
		return nil, (&er.Error{
			Package:     ErrPkg,
			Type:        ErrTypeAES,
			Description: ErrDscEncrypt,
			Detail:      "aes cbc failed.",
		}).Wrap(err)
	}

	ciphertext, err := blockEncrypt(aes.BlockSize, cipher.NewCBCEncrypter(c, iv), plaintext)
	if err != nil {
		return nil, (&er.Error{
			Package:     ErrPkg,
			Type:        ErrTypeAES,
			Description: ErrDscEncrypt,
			Detail:      "aes cbc failed.",
		}).Wrap(err)
	}
	ciphertext = append(iv, ciphertext...)

	return ciphertext, nil
}

// DecryptAESCBC decrypts ciphertext with AES CBC cipher algorithm.
// The key argument should be the AES key, either 16, 24, or 32 bytes to select AES-128, AES-192, or AES-256.
//
// See more details at:
//   - https://pkg.go.dev/crypto/aes
//   - https://pkg.go.dev/crypto/cipher
func DecryptAESCBC(key []byte, ciphertext []byte) ([]byte, error) {
	n := len(ciphertext)
	if n < 2*aes.BlockSize {
		// ciphertext is shorter than length of (padding + iv).
		return nil, &er.Error{
			Package:     ErrPkg,
			Type:        ErrTypeAES,
			Description: ErrDscDecrypt,
			Detail:      "aes cbc invalid ciphertext length.",
		}
	}

	iv, ciphertext := ciphertext[:aes.BlockSize], ciphertext[aes.BlockSize:]

	c, err := aes.NewCipher(key)
	if err != nil {
		return nil, (&er.Error{
			Package:     ErrPkg,
			Type:        ErrTypeAES,
			Description: ErrDscDecrypt,
			Detail:      "aes cbc failed.",
		}).Wrap(err)
	}

	plaintext, err := blockDecrypt(aes.BlockSize, cipher.NewCBCDecrypter(c, iv), ciphertext)
	if err != nil {
		return nil, (&er.Error{
			Package:     ErrPkg,
			Type:        ErrTypeAES,
			Description: ErrDscDecrypt,
			Detail:      "aes cbc failed.",
		}).Wrap(err)
	}

	return plaintext, nil
}

// EncryptAESCFB encrypts plaintext with AES CFB cipher algorithm.
// The key argument should be the AES key, either 16, 24, or 32 bytes to select AES-128, AES-192, or AES-256.
//
// See more details at:
//   - https://pkg.go.dev/crypto/aes
//   - https://pkg.go.dev/crypto/cipher
func EncryptAESCFB(key []byte, plaintext []byte) ([]byte, error) {
	c, iv, err := newAES(key)
	if err != nil {
		return nil, (&er.Error{
			Package:     ErrPkg,
			Type:        ErrTypeAES,
			Description: ErrDscEncrypt,
			Detail:      "aes cfb failed.",
		}).Wrap(err)
	}

	ciphertext := streamEncrypt(cipher.NewCFBEncrypter(c, iv), plaintext)
	ciphertext = append(iv, ciphertext...)

	return ciphertext, nil
}

// DecryptAESCFB decrypts ciphertext with AES CFB cipher algorithm.
// The key argument should be the AES key, either 16, 24, or 32 bytes to select AES-128, AES-192, or AES-256.
//
// See more details at:
//   - https://pkg.go.dev/crypto/aes
//   - https://pkg.go.dev/crypto/cipher
func DecryptAESCFB(key []byte, ciphertext []byte) ([]byte, error) {
	n := len(ciphertext)
	if n < aes.BlockSize {
		// ciphertext is shorter than length of (padding + iv).
		return nil, &er.Error{
			Package:     ErrPkg,
			Type:        ErrTypeAES,
			Description: ErrDscDecrypt,
			Detail:      "aes cfb invalid ciphertext length.",
		}
	}

	iv, ciphertext := ciphertext[:aes.BlockSize], ciphertext[aes.BlockSize:]

	c, err := aes.NewCipher(key)
	if err != nil {
		return nil, (&er.Error{
			Package:     ErrPkg,
			Type:        ErrTypeAES,
			Description: ErrDscDecrypt,
			Detail:      "aes cfb failed.",
		}).Wrap(err)
	}

	plaintext := streamDecrypt(cipher.NewCFBDecrypter(c, iv), ciphertext)

	return plaintext, nil
}

// EncryptAESCTR encrypts plaintext with AES CTR cipher algorithm.
// The key argument should be the AES key, either 16, 24, or 32 bytes to select AES-128, AES-192, or AES-256.
//
// See more details at:
//   - https://pkg.go.dev/crypto/aes
//   - https://pkg.go.dev/crypto/cipher
func EncryptAESCTR(key []byte, plaintext []byte) ([]byte, error) {
	c, iv, err := newAES(key)
	if err != nil {
		return nil, (&er.Error{
			Package:     ErrPkg,
			Type:        ErrTypeAES,
			Description: ErrDscEncrypt,
			Detail:      "aes ctr failed.",
		}).Wrap(err)
	}

	ciphertext := streamEncrypt(cipher.NewCTR(c, iv), plaintext)
	ciphertext = append(iv, ciphertext...)

	return ciphertext, nil
}

// DecryptAESCTR decrypts ciphertext with AES CTR cipher algorithm.
// The key argument should be the AES key, either 16, 24, or 32 bytes to select AES-128, AES-192, or AES-256.
//
// See more details at:
//   - https://pkg.go.dev/crypto/aes
//   - https://pkg.go.dev/crypto/cipher
func DecryptAESCTR(key []byte, ciphertext []byte) ([]byte, error) {
	n := len(ciphertext)
	if n < aes.BlockSize {
		// ciphertext is shorter than length of (padding + iv).
		return nil, &er.Error{
			Package:     ErrPkg,
			Type:        ErrTypeAES,
			Description: ErrDscDecrypt,
			Detail:      "aes ctr invalid ciphertext length.",
		}
	}

	iv, ciphertext := ciphertext[:aes.BlockSize], ciphertext[aes.BlockSize:]

	c, err := aes.NewCipher(key)
	if err != nil {
		return nil, (&er.Error{
			Package:     ErrPkg,
			Type:        ErrTypeAES,
			Description: ErrDscDecrypt,
			Detail:      "aes ctr failed.",
		}).Wrap(err)
	}

	plaintext := streamDecrypt(cipher.NewCTR(c, iv), ciphertext)

	return plaintext, nil
}

// EncryptAESOFB encrypts plaintext with AES OFB cipher algorithm.
// The key argument should be the AES key, either 16, 24, or 32 bytes to select AES-128, AES-192, or AES-256.
//
// See more details at:
//   - https://pkg.go.dev/crypto/aes
//   - https://pkg.go.dev/crypto/cipher
func EncryptAESOFB(key []byte, plaintext []byte) ([]byte, error) {
	c, iv, err := newAES(key)
	if err != nil {
		return nil, (&er.Error{
			Package:     ErrPkg,
			Type:        ErrTypeAES,
			Description: ErrDscEncrypt,
			Detail:      "aes ofb failed.",
		}).Wrap(err)
	}

	ciphertext := streamEncrypt(cipher.NewOFB(c, iv), plaintext)
	ciphertext = append(iv, ciphertext...)

	return ciphertext, nil
}

// DecryptAESOFB decrypts ciphertext with AES OFB cipher algorithm.
// The key argument should be the AES key, either 16, 24, or 32 bytes to select AES-128, AES-192, or AES-256.
//
// See more details at:
//   - https://pkg.go.dev/crypto/aes
//   - https://pkg.go.dev/crypto/cipher
func DecryptAESOFB(key []byte, ciphertext []byte) ([]byte, error) {
	n := len(ciphertext)
	if n < aes.BlockSize {
		// ciphertext is shorter than length of (padding + iv).
		return nil, &er.Error{
			Package:     ErrPkg,
			Type:        ErrTypeAES,
			Description: ErrDscDecrypt,
			Detail:      "aes ofb invalid ciphertext length.",
		}
	}

	iv, ciphertext := ciphertext[:aes.BlockSize], ciphertext[aes.BlockSize:]

	c, err := aes.NewCipher(key)
	if err != nil {
		return nil, (&er.Error{
			Package:     ErrPkg,
			Type:        ErrTypeAES,
			Description: ErrDscDecrypt,
			Detail:      "aes ofb failed.",
		}).Wrap(err)
	}

	plaintext := streamDecrypt(cipher.NewOFB(c, iv), ciphertext)

	return plaintext, nil
}
