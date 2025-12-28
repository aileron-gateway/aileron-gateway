// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package encoder

const (
	ErrPkg = "encoder"

	ErrTypeBase16 = "json"
	ErrTypeBase32 = "json"
	ErrTypeBase64 = "json"
	ErrTypeJSON   = "json"
	ErrTypeYaml   = "yaml"
	ErrTypeProto  = "proto"

	// ErrDscHash is a error description.
	// This description indicates the failure of
	// hash calculation.
	ErrDscDecode = "hash calculation failed."

	// ErrDscMarshal is a error description.
	// This description indicates the failure of
	// marshaling data.
	ErrDscMarshal = "marshaling failed"
	// ErrDscUnmarshal is a error description.
	// This description indicates the failure of
	// unmarshaling data.
	ErrDscUnmarshal = "unmarshaling failed"

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
