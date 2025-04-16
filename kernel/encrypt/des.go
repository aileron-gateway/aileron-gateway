// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package encrypt

import (
	"crypto/cipher"
	"crypto/des"
	"crypto/rand"
	"io"

	"github.com/aileron-gateway/aileron-gateway/kernel/er"
)

type desFunc func(key []byte) (cipher.Block, []byte, error)

func newDES(key []byte) (cipher.Block, []byte, error) {
	c, err := des.NewCipher(key)
	if err != nil {
		return nil, nil, err
	}

	iv := make([]byte, c.BlockSize())
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, nil, err
	}

	return c, iv[:], nil
}

func newTripleDES(key []byte) (cipher.Block, []byte, error) {
	c, err := des.NewTripleDESCipher(key)
	if err != nil {
		return nil, nil, err
	}

	iv := make([]byte, c.BlockSize())
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, nil, err
	}

	return c, iv[:], nil
}

// EncryptDESCBC encrypts plaintext with DES CBC cipher algorithm.
// The key must be the DES key, 8 bytes.
//
// See more details at:
//   - https://pkg.go.dev/crypto/des
//   - https://pkg.go.dev/crypto/cipher
func EncryptDESCBC(key []byte, plaintext []byte) ([]byte, error) {
	return encryptDESCBC(key, plaintext, newDES)
}

// DecryptDESCBC decrypts ciphertext with DES CBC cipher algorithm.
// The key must be the DES key, 8 bytes.
//
// See more details at:
//   - https://pkg.go.dev/crypto/des
//   - https://pkg.go.dev/crypto/cipher
func DecryptDESCBC(key []byte, plaintext []byte) ([]byte, error) {
	return decryptDESCBC(key, plaintext, newDES)
}

// EncryptTripleDESCBC encrypts plaintext with 3DES CBC cipher algorithm.
// The key must be the 3DES key, 24 bytes.
//
// See more details at:
//   - https://pkg.go.dev/crypto/des
//   - https://pkg.go.dev/crypto/cipher
func EncryptTripleDESCBC(key []byte, plaintext []byte) ([]byte, error) {
	return encryptDESCBC(key, plaintext, newTripleDES)
}

// DecryptTripleDESCBC decrypts ciphertext with 3DES CBC cipher algorithm.
// The key must be the 3DES key, 24 bytes.
//
// See more details at:
//   - https://pkg.go.dev/crypto/des
//   - https://pkg.go.dev/crypto/cipher
func DecryptTripleDESCBC(key []byte, plaintext []byte) ([]byte, error) {
	return decryptDESCBC(key, plaintext, newTripleDES)
}

func encryptDESCBC(key []byte, plaintext []byte, f desFunc) ([]byte, error) {
	c, iv, err := f(key)
	if err != nil {
		return nil, (&er.Error{
			Package:     ErrPkg,
			Type:        ErrTypeDES,
			Description: ErrDscEncrypt,
			Detail:      "des cbc failed.",
		}).Wrap(err)
	}

	ciphertext, err := blockEncrypt(des.BlockSize, cipher.NewCBCEncrypter(c, iv), plaintext)
	if err != nil {
		return nil, (&er.Error{
			Package:     ErrPkg,
			Type:        ErrTypeDES,
			Description: ErrDscEncrypt,
			Detail:      "des cbc failed.",
		}).Wrap(err)
	}
	ciphertext = append(iv, ciphertext...)

	return ciphertext, nil
}

func decryptDESCBC(key []byte, ciphertext []byte, f desFunc) ([]byte, error) {
	n := len(ciphertext)
	if n < 2*des.BlockSize {
		// ciphertext is shorter than length of (padding + iv).
		return nil, &er.Error{
			Package:     ErrPkg,
			Type:        ErrTypeDES,
			Description: ErrDscDecrypt,
			Detail:      "des cbc invalid ciphertext length.",
		}
	}

	iv, ciphertext := ciphertext[:des.BlockSize], ciphertext[des.BlockSize:]

	c, _, err := f(key)
	if err != nil {
		return nil, (&er.Error{
			Package:     ErrPkg,
			Type:        ErrTypeDES,
			Description: ErrDscDecrypt,
			Detail:      "des cbc failed.",
		}).Wrap(err)
	}

	plaintext, err := blockDecrypt(des.BlockSize, cipher.NewCBCDecrypter(c, iv), ciphertext)
	if err != nil {
		return nil, (&er.Error{
			Package:     ErrPkg,
			Type:        ErrTypeDES,
			Description: ErrDscDecrypt,
			Detail:      "des cbc failed.",
		}).Wrap(err)
	}

	return plaintext, nil
}

// EncryptDESCFB encrypts plaintext with DES CFB cipher algorithm.
// The key must be the DES key, 8 bytes.
//
// See more details at:
//   - https://pkg.go.dev/crypto/des
//   - https://pkg.go.dev/crypto/cipher
func EncryptDESCFB(key []byte, plaintext []byte) ([]byte, error) {
	return encryptDESCFB(key, plaintext, newDES)
}

// DecryptDESCFB decrypts ciphertext with DES CFC cipher algorithm.
// The key must be the DES key, 8 bytes.
//
// See more details at:
//   - https://pkg.go.dev/crypto/des
//   - https://pkg.go.dev/crypto/cipher
func DecryptDESCFB(key []byte, plaintext []byte) ([]byte, error) {
	return decryptDESCFB(key, plaintext, newDES)
}

// EncryptTripleDESCFB encrypts plaintext with 3DES CFB cipher algorithm.
// The key must be the 3DES key, 24 bytes.
//
// See more details at:
//   - https://pkg.go.dev/crypto/des
//   - https://pkg.go.dev/crypto/cipher
func EncryptTripleDESCFB(key []byte, plaintext []byte) ([]byte, error) {
	return encryptDESCFB(key, plaintext, newTripleDES)
}

// DecryptTripleDESCFB decrypts ciphertext with 3DES CFC cipher algorithm.
// The key must be the 3DES key, 24 bytes.
//
// See more details at:
//   - https://pkg.go.dev/crypto/des
//   - https://pkg.go.dev/crypto/cipher
func DecryptTripleDESCFB(key []byte, plaintext []byte) ([]byte, error) {
	return decryptDESCFB(key, plaintext, newTripleDES)
}

func encryptDESCFB(key []byte, plaintext []byte, f desFunc) ([]byte, error) {
	c, iv, err := f(key)
	if err != nil {
		return nil, (&er.Error{
			Package:     ErrPkg,
			Type:        ErrTypeDES,
			Description: ErrDscEncrypt,
			Detail:      "des cfb failed.",
		}).Wrap(err)
	}

	ciphertext := streamEncrypt(cipher.NewCFBEncrypter(c, iv), plaintext)
	ciphertext = append(iv, ciphertext...)

	return ciphertext, nil
}

func decryptDESCFB(key []byte, ciphertext []byte, f desFunc) ([]byte, error) {
	n := len(ciphertext)
	if n < des.BlockSize {
		// ciphertext is shorter than length of (padding + iv).
		return nil, &er.Error{
			Package:     ErrPkg,
			Type:        ErrTypeDES,
			Description: ErrDscDecrypt,
			Detail:      "des cfb invalid ciphertext length.",
		}
	}

	iv, ciphertext := ciphertext[:des.BlockSize], ciphertext[des.BlockSize:]

	c, _, err := f(key)
	if err != nil {
		return nil, (&er.Error{
			Package:     ErrPkg,
			Type:        ErrTypeDES,
			Description: ErrDscDecrypt,
			Detail:      "des cfb failed.",
		}).Wrap(err)
	}

	plaintext := streamDecrypt(cipher.NewCFBDecrypter(c, iv), ciphertext)

	return plaintext, nil
}

// EncryptDESCTR encrypts plaintext with DES CTR cipher algorithm.
// The key must be the DES key, 8 bytes.
//
// See more details at:
//   - https://pkg.go.dev/crypto/des
//   - https://pkg.go.dev/crypto/cipher
func EncryptDESCTR(key []byte, plaintext []byte) ([]byte, error) {
	return encryptDESCTR(key, plaintext, newDES)
}

// DecryptDESCTR decrypts ciphertext with DES CTR cipher algorithm.
// The key must be the DES key, 8 bytes.
//
// See more details at:
//   - https://pkg.go.dev/crypto/des
//   - https://pkg.go.dev/crypto/cipher
func DecryptDESCTR(key []byte, plaintext []byte) ([]byte, error) {
	return decryptDESCTR(key, plaintext, newDES)
}

// EncryptTripleDESCTR encrypts plaintext with 3DES CTR cipher algorithm.
// The key must be the 3DES key, 24 bytes.
//
// See more details at:
//   - https://pkg.go.dev/crypto/des
//   - https://pkg.go.dev/crypto/cipher
func EncryptTripleDESCTR(key []byte, plaintext []byte) ([]byte, error) {
	return encryptDESCTR(key, plaintext, newTripleDES)
}

// DecryptTripleDESCTR decrypts ciphertext with 3DES CTR cipher algorithm.
// The key must be the 3DES key, 24 bytes.
//
// See more details at:
//   - https://pkg.go.dev/crypto/des
//   - https://pkg.go.dev/crypto/cipher
func DecryptTripleDESCTR(key []byte, plaintext []byte) ([]byte, error) {
	return decryptDESCTR(key, plaintext, newTripleDES)
}

func encryptDESCTR(key []byte, plaintext []byte, f desFunc) ([]byte, error) {
	c, iv, err := f(key)
	if err != nil {
		return nil, (&er.Error{
			Package:     ErrPkg,
			Type:        ErrTypeDES,
			Description: ErrDscEncrypt,
			Detail:      "des ctr failed.",
		}).Wrap(err)
	}

	ciphertext := streamEncrypt(cipher.NewCTR(c, iv), plaintext)
	ciphertext = append(iv, ciphertext...)

	return ciphertext, nil
}

func decryptDESCTR(key []byte, ciphertext []byte, f desFunc) ([]byte, error) {
	n := len(ciphertext)
	if n < des.BlockSize {
		// ciphertext is shorter than length of (padding + iv).
		return nil, &er.Error{
			Package:     ErrPkg,
			Type:        ErrTypeDES,
			Description: ErrDscDecrypt,
			Detail:      "des ofb invalid ciphertext length.",
		}
	}

	iv, ciphertext := ciphertext[:des.BlockSize], ciphertext[des.BlockSize:]

	c, _, err := f(key)
	if err != nil {
		return nil, (&er.Error{
			Package:     ErrPkg,
			Type:        ErrTypeDES,
			Description: ErrDscDecrypt,
			Detail:      "des ctr failed.",
		}).Wrap(err)
	}

	plaintext := streamDecrypt(cipher.NewCTR(c, iv), ciphertext)

	return plaintext, nil
}

// EncryptDESOFB encrypts plaintext with DES OFB cipher algorithm.
// The key must be the DES key, 8 bytes.
//
// See more details at:
//   - https://pkg.go.dev/crypto/des
//   - https://pkg.go.dev/crypto/cipher
func EncryptDESOFB(key []byte, plaintext []byte) ([]byte, error) {
	return encryptDESOFB(key, plaintext, newDES)
}

// DecryptDESOFB decrypts ciphertext with DES OFB cipher algorithm.
// The key must be the DES key, 8 bytes.
//
// See more details at:
//   - https://pkg.go.dev/crypto/des
//   - https://pkg.go.dev/crypto/cipher
func DecryptDESOFB(key []byte, plaintext []byte) ([]byte, error) {
	return decryptDESOFB(key, plaintext, newDES)
}

// EncryptTripleDESOFB encrypts plaintext with 3DES OFB cipher algorithm.
// The key must be the 3DES key, 24 bytes.
//
// See more details at:
//   - https://pkg.go.dev/crypto/des
//   - https://pkg.go.dev/crypto/cipher
func EncryptTripleDESOFB(key []byte, plaintext []byte) ([]byte, error) {
	return encryptDESOFB(key, plaintext, newTripleDES)
}

// DecryptTripleDESOFB decrypts ciphertext with 3DES OFB cipher algorithm.
// The key must be the 3DES key, 24 bytes.
//
// See more details at:
//   - https://pkg.go.dev/crypto/des
//   - https://pkg.go.dev/crypto/cipher
func DecryptTripleDESOFB(key []byte, plaintext []byte) ([]byte, error) {
	return decryptDESOFB(key, plaintext, newTripleDES)
}

func encryptDESOFB(key []byte, plaintext []byte, f desFunc) ([]byte, error) {
	c, iv, err := f(key)
	if err != nil {
		return nil, (&er.Error{
			Package:     ErrPkg,
			Type:        ErrTypeDES,
			Description: ErrDscEncrypt,
			Detail:      "des ofb failed.",
		}).Wrap(err)
	}

	ciphertext := streamEncrypt(cipher.NewOFB(c, iv), plaintext)
	ciphertext = append(iv, ciphertext...)

	return ciphertext, nil
}

func decryptDESOFB(key []byte, ciphertext []byte, f desFunc) ([]byte, error) {
	n := len(ciphertext)
	if n < des.BlockSize {
		// ciphertext is shorter than length of (padding + iv).
		return nil, &er.Error{
			Package:     ErrPkg,
			Type:        ErrTypeDES,
			Description: ErrDscDecrypt,
			Detail:      "des ofb invalid ciphertext length.",
		}
	}

	iv, ciphertext := ciphertext[:des.BlockSize], ciphertext[des.BlockSize:]

	c, _, err := f(key)
	if err != nil {
		return nil, (&er.Error{
			Package:     ErrPkg,
			Type:        ErrTypeDES,
			Description: ErrDscDecrypt,
			Detail:      "des ofb failed.",
		}).Wrap(err)
	}

	plaintext := streamDecrypt(cipher.NewOFB(c, iv), ciphertext)

	return plaintext, nil
}
