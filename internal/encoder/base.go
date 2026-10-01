// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package encoder

import (
	"encoding/base32"
	"encoding/base64"
	"encoding/hex"

	k "github.com/aileron-gateway/aileron-gateway/apis/kernel"
)

type EncodeType int

const (
	EncodeTypeUnknown EncodeType = iota
	EncodeTypeBase16
	EncodeTypeBase32
	EncodeTypeBase32Hex
	EncodeTypeBase32Escaped
	EncodeTypeBase32HexEscaped
	EncodeTypeBase64
	EncodeTypeBase64Raw
	EncodeTypeBase64URL
	EncodeTypeBase64RawURL
)

var EncodeTypes = map[k.EncodingType]EncodeType{
	k.EncodingType_Base16:           EncodeTypeBase16,
	k.EncodingType_Base32:           EncodeTypeBase32,
	k.EncodingType_Base32Hex:        EncodeTypeBase32Hex,
	k.EncodingType_Base32Escaped:    EncodeTypeBase32Escaped,
	k.EncodingType_Base32HexEscaped: EncodeTypeBase32HexEscaped,
	k.EncodingType_Base64:           EncodeTypeBase64,
	k.EncodingType_Base64Raw:        EncodeTypeBase64Raw,
	k.EncodingType_Base64URL:        EncodeTypeBase64URL,
	k.EncodingType_Base64RawURL:     EncodeTypeBase64RawURL,
}

// Base32StdEscapedEncoding is the base32 based encoding.
// It differs from original base32 encoding in that this
// encoding does not contain vowels "A",  "E", "I", "O".
// This technique prevent the encoded string from containing un-willing words
// like "ERROR", "FATAL" or "WARN" and so on.
//   - "0" added instead of "A"
//   - "1" added instead of "E"
//   - "8" added instead of "I"
//   - "9" added instead of "O"
var Base32StdEscapedEncoding = base32.NewEncoding("BCDFGHJKLMNPQRSTUVWXYZ0123456789")

// Base32HexEscapedEncoding is the base32 hex based encoding.
// It differs from original base32 hex encoding in that this
// encoding does not contain vowels "A",  "E", "I", "O".
// This technique prevent the encoded string from containing un-willing words
// like "ERROR", "FATAL" or "WARN" and so on.
//   - "W" added instead of "A"
//   - "X" added instead of "E"
//   - "Y" added instead of "I"
//   - "Z" added instead of "O"
var Base32HexEscapedEncoding = base32.NewEncoding("0123456789BCDFGHJKLMNPQRSTUVWXYZ")

// EncodeToStringFunc is the function that encode bytes array into string.
type EncodeToStringFunc func(data []byte) (encoded string)

// DecodeStringFunc is the function that decode string into bytes array.
type DecodeStringFunc func(data string) (decoded []byte, err error)

// EncoderDecoder returns the pair of encoder and decoder function.
// If invalid argument is given, nil pointers are returned.
// Default Base32HexEscapedEncoding will be used if the unknown
// encoding type was given.
func EncoderDecoder(e k.EncodingType) (EncodeToStringFunc, DecodeStringFunc) {
	switch e {
	case k.EncodingType_Base16:
		return hex.EncodeToString, hex.DecodeString
	case k.EncodingType_Base32:
		return base32.StdEncoding.EncodeToString, base32.StdEncoding.DecodeString
	case k.EncodingType_Base32Hex:
		return base32.HexEncoding.EncodeToString, base32.HexEncoding.DecodeString
	case k.EncodingType_Base32Escaped:
		return Base32StdEscapedEncoding.EncodeToString, Base32StdEscapedEncoding.DecodeString
	case k.EncodingType_Base32HexEscaped:
		return Base32HexEscapedEncoding.EncodeToString, Base32HexEscapedEncoding.DecodeString
	case k.EncodingType_Base64:
		return base64.StdEncoding.EncodeToString, base64.StdEncoding.DecodeString
	case k.EncodingType_Base64Raw:
		return base64.RawStdEncoding.EncodeToString, base64.RawStdEncoding.DecodeString
	case k.EncodingType_Base64URL:
		return base64.URLEncoding.EncodeToString, base64.URLEncoding.DecodeString
	case k.EncodingType_Base64RawURL:
		return base64.RawURLEncoding.EncodeToString, base64.RawURLEncoding.DecodeString
	default:
		return Base32HexEscapedEncoding.EncodeToString, Base32HexEscapedEncoding.DecodeString
	}
}
