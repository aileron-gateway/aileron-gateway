// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package encoder

import (
	"encoding/base32"
	"encoding/base64"
	"encoding/hex"

	k "github.com/aileron-gateway/aileron-gateway/apis/kernel"
	"github.com/aileron-gateway/aileron-gateway/kernel/er"
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

// EncodeFunc is the function that encode the given data.
// type EncodeFunc func(data []byte) (encoded []byte)

// EncodeToStringFunc is the function that encode bytes array into string.
type EncodeToStringFunc func(data []byte) (encoded string)

// DecodeFunc is the function that decode the given data.
// type DecodeFunc func(data []byte) (decoded []byte, err error)

// DecodeStringFunc is the function that decode string into bytes array.
type DecodeStringFunc func(data string) (decoded []byte, err error)

// Base16Encode encode byte array into string with Base16 (hex) encoding.
func Base16Encode(src []byte) string {
	return hex.EncodeToString(src)
}

// Base16Decode decode string into byte array with Base16 (hex) decoding.
func Base16Decode(src string) ([]byte, error) {
	b, err := hex.DecodeString(src)
	if err != nil {
		return nil, (&er.Error{
			Package:     ErrPkg,
			Type:        ErrTypeBase16,
			Description: ErrDscDecode,
		}).Wrap(err)
	}
	return b, nil
}

// Base32Encode encode byte array into string with Base32 standard encoding.
func Base32Encode(src []byte) string {
	return base32.StdEncoding.EncodeToString(src)
}

// Base32Decode decode string into byte array with Base32 decoding.
func Base32Decode(src string) ([]byte, error) {
	b, err := base32.StdEncoding.DecodeString(src)
	if err != nil {
		return nil, (&er.Error{
			Package:     ErrPkg,
			Type:        ErrTypeBase32,
			Description: ErrDscDecode,
			Detail:      "base32 standard encoding.",
		}).Wrap(err)
	}
	return b, nil
}

// Base32EscapedEncode encode byte array into string with Base32 standard encoding with escape.
func Base32EscapedEncode(src []byte) string {
	return Base32StdEscapedEncoding.EncodeToString(src)
}

// Base32EscapedDecode decode string into byte array with Base32 Escaped decoding.
func Base32EscapedDecode(src string) ([]byte, error) {
	b, err := Base32StdEscapedEncoding.DecodeString(src)
	if err != nil {
		return nil, (&er.Error{
			Package:     ErrPkg,
			Type:        ErrTypeBase32,
			Description: ErrDscDecode,
			Detail:      "base32 standard escaped encoding.",
		}).Wrap(err)
	}
	return b, nil
}

// Base32HexEncode encode byte array into string with Base32 Hex encoding.
func Base32HexEncode(src []byte) string {
	return base32.HexEncoding.EncodeToString(src)
}

// Base32HexDecode decode string into byte array with Base32 Hex decoding.
func Base32HexDecode(src string) ([]byte, error) {
	b, err := base32.HexEncoding.DecodeString(src)
	if err != nil {
		return nil, (&er.Error{
			Package:     ErrPkg,
			Type:        ErrTypeBase32,
			Description: ErrDscDecode,
			Detail:      "base32 hex encoding.",
		}).Wrap(err)
	}
	return b, nil
}

// Base32HexEscapedEncode encode byte array into string with Base32 Hex encoding with escaping.
func Base32HexEscapedEncode(src []byte) string {
	return Base32HexEscapedEncoding.EncodeToString(src)
}

// Base32HexEscapedDecode decode string into byte array with Base32 Hex Escaped decoding.
func Base32HexEscapedDecode(src string) ([]byte, error) {
	b, err := Base32HexEscapedEncoding.DecodeString(src)
	if err != nil {
		return nil, (&er.Error{
			Package:     ErrPkg,
			Type:        ErrTypeBase32,
			Description: ErrDscDecode,
			Detail:      "base32 hex escaped encoding.",
		}).Wrap(err)
	}
	return b, nil
}

// Base64Encode encode byte array into string with Base64 encoding.
func Base64Encode(src []byte) string {
	return base64.StdEncoding.EncodeToString(src)
}

// Base64Decode decode string into byte array with Base64 decoding.
func Base64Decode(src string) ([]byte, error) {
	b, err := base64.StdEncoding.DecodeString(src)
	if err != nil {
		return nil, (&er.Error{
			Package:     ErrPkg,
			Type:        ErrTypeBase64,
			Description: ErrDscDecode,
			Detail:      "base64 standard encoding.",
		}).Wrap(err)
	}
	return b, nil
}

// Base64RawEncode encode byte array into string with Base64 Raw encoding.
func Base64RawEncode(src []byte) string {
	return base64.RawStdEncoding.EncodeToString(src)
}

// Base64RawDecode decode string into byte array with Base64 Raw decoding.
func Base64RawDecode(src string) ([]byte, error) {
	b, err := base64.RawStdEncoding.DecodeString(src)
	if err != nil {
		return nil, (&er.Error{
			Package:     ErrPkg,
			Type:        ErrTypeBase64,
			Description: ErrDscDecode,
			Detail:      "base64 raw encoding.",
		}).Wrap(err)
	}
	return b, nil
}

// Base64URLEncode encode byte array into string with Base64 URL encoding.
func Base64URLEncode(src []byte) string {
	return base64.URLEncoding.EncodeToString(src)
}

// Base64URLDecode decode string into byte array with Base64 URL decoding.
func Base64URLDecode(src string) ([]byte, error) {
	b, err := base64.URLEncoding.DecodeString(src)
	if err != nil {
		return nil, (&er.Error{
			Package:     ErrPkg,
			Type:        ErrTypeBase64,
			Description: ErrDscDecode,
			Detail:      "base64 url encoding.",
		}).Wrap(err)
	}
	return b, nil
}

// Base64RawURLEncode encode byte array into string with Base64 Raw URL encoding.
func Base64RawURLEncode(src []byte) string {
	return base64.RawURLEncoding.EncodeToString(src)
}

// Base64RawURLDecode decode string into byte array with Base64 Raw URL decoding.
func Base64RawURLDecode(src string) ([]byte, error) {
	b, err := base64.RawURLEncoding.DecodeString(src)
	if err != nil {
		return nil, (&er.Error{
			Package:     ErrPkg,
			Type:        ErrTypeBase64,
			Description: ErrDscDecode,
			Detail:      "base64 raw url encoding.",
		}).Wrap(err)
	}
	return b, nil
}

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
