// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package encrypt

const (
	ErrPkg = "encrypt"

	ErrTypeBCrypt   = "BCRypt"
	ErrTypeSCrypt   = "SCrypt"
	ErrTypePBKDF2   = "PBKDF2"
	ErrTypeArgon2i  = "Argon2i"
	ErrTypeArgon2id = "Argon2id"

	ErrTypeStream  = "stream"
	ErrTypeBlock   = "block"
	ErrTypePadding = "padding"

	ErrTypeAES = "aes"
	ErrTypeDES = "des"
	ErrTypeRC4 = "rc4"

	// ErrDscHash is a error description.
	// This description indicates the failure of
	// hash calculation.
	ErrDscHash = "hash calculation failed."
	// ErrDscHashValid is a error description.
	// This description indicates the failure of
	// validations of hashers.
	ErrDscHashValid = "hash validation failed."

	// ErrDscPadding is a error description.
	// This description indicates failure of
	// padding data.
	ErrDscPadding = "padding failed."
	// ErrDscUnpadding is a error description.
	// This description indicates failure of
	// unpadding data.
	ErrDscUnpadding = "unpadding failed."

	// ErrDscEncrypt is a error description.
	// This description indicates failure of
	// plaintext encryption.
	ErrDscEncrypt = "failed to encrypt plain text."
	// ErrDscDecrypt is a error description.
	// This description indicates failure of
	// ciphertext decryption.
	ErrDscDecrypt = "failed to decrypt cipher text."
)
