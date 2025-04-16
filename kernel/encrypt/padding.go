// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package encrypt

import (
	"crypto/rand"
	"io"

	"github.com/aileron-gateway/aileron-gateway/kernel/er"
)

// PKCS7Pad appends padding bytes to given bytes slice and return the padded slice.
// This function can be used as PKCS#5 because PKCS#5 is the subset of PKCS#7.
// PKCS#7 adds padding with the range of 0x01-0x10 (1-15 in decimal).
// PKCS#5 adds padding with the range of 0x01-0x08 (1-7 in decimal).
// If a negative value or zero is given as blockSize, an error is returned.
// For example, when blockSize is 6 and []byte("abc") = {0x61, 0x61, 0x63} is given as the data,
// then padded data will be {0x61, 0x61, 0x63, 0x03, 0x03, 0x03}.
func PKCS7Pad(blockSize int, data []byte) ([]byte, error) {
	// Restrict block size from 1 to 255 because the range of byte is 0-255.
	if blockSize < 1 || blockSize > 255 {
		return nil, &er.Error{
			Package:     ErrPkg,
			Type:        ErrTypePadding,
			Description: ErrDscPadding,
			Detail:      "PKCS7 invalid block size",
		}
	}

	n := blockSize - len(data)%blockSize
	pad := make([]byte, n)
	for i := 0; i < n; i++ {
		pad[i] = byte(n)
	}

	data = append(data, pad...)
	return data, nil
}

// PKCS7UnPad removes padding using PKCS7 padding.
func PKCS7UnPad(blockSize int, data []byte) ([]byte, error) {
	// Restrict block size from 1 to 255 because the range of byte is 0-255.
	if blockSize < 1 || blockSize > 255 {
		return nil, &er.Error{
			Package:     ErrPkg,
			Type:        ErrTypePadding,
			Description: ErrDscUnpadding,
			Detail:      "PKCS7 invalid block size",
		}
	}

	// Check that the given data has valid length as padded data.
	n := len(data)
	if n < blockSize || n%blockSize != 0 {
		return nil, &er.Error{
			Package:     ErrPkg,
			Type:        ErrTypePadding,
			Description: ErrDscUnpadding,
			Detail:      "PKCS7 invalid data size",
		}
	}

	pad := int(data[n-1])
	if n < pad {
		return nil, &er.Error{
			Package:     ErrPkg,
			Type:        ErrTypePadding,
			Description: ErrDscUnpadding,
			Detail:      "PKCS7 invalid padding size",
		}
	}

	return data[:n-pad], nil
}

// ISO7816Pad appends padding bytes to given bytes slice and return the padded slice.
// If a negative value or zero is given as blockSize, an error is returned.
// For example, when blockSize is 6 and []byte("abc") = {0x61, 0x61, 0x63} is given as data,
// then padded data will be {0x61, 0x61, 0x63, 0x80, 0x00, 0x00}.
// 0x80 is the marker of the end of the actual data.
func ISO7816Pad(blockSize int, data []byte) ([]byte, error) {
	// Restrict block size from 1 to 255 because the range of byte is 0-255.
	if blockSize < 1 || blockSize > 255 {
		return nil, &er.Error{
			Package:     ErrPkg,
			Type:        ErrTypePadding,
			Description: ErrDscPadding,
			Detail:      "ISO7816 invalid block size",
		}
	}

	n := blockSize - len(data)%blockSize
	pad := make([]byte, n)
	pad[0] = 0x80
	for i := 1; i < n; i++ {
		pad[i] = 0x00
	}

	data = append(data, pad...)
	return data, nil
}

// ISO7816UnPad removes padding using ISO7816 padding.
func ISO7816UnPad(blockSize int, data []byte) ([]byte, error) {
	// Restrict block size from 1 to 255 because the range of byte is 0-255.
	if blockSize < 1 || blockSize > 255 {
		return nil, &er.Error{
			Package:     ErrPkg,
			Type:        ErrTypePadding,
			Description: ErrDscUnpadding,
			Detail:      "ISO7816 invalid block size",
		}
	}

	// Check that the given data has valid length as padded data.
	n := len(data)
	if n < blockSize || n%blockSize != 0 {
		return nil, &er.Error{
			Package:     ErrPkg,
			Type:        ErrTypePadding,
			Description: ErrDscUnpadding,
			Detail:      "ISO7816 invalid data length",
		}
	}

	pad := 0
	for i := 1; i <= blockSize; i++ {
		if data[n-i] == 0x80 {
			pad = i
			break
		} else if data[n-i] != 0x00 {
			break
		}
	}

	if pad == 0 {
		return nil, &er.Error{
			Package:     ErrPkg,
			Type:        ErrTypePadding,
			Description: ErrDscUnpadding,
			Detail:      "ISO7816 invalid padding size",
		}
	}

	return data[:n-pad], nil
}

// ISO10126Pad appends padding bytes to given bytes slice and return the padded slice.
// If a negative value or zero is given as blockSize, an error is returned.
// For example, when blockSize is 6 and []byte("abc") = {0x61, 0x61, 0x63} is given as data,
// then padded data is {0x61, 0x61, 0x63, 0xaa, 0xbb, 0x03}.
// The last byte of 0x03 represents the length of padding
// and the rest of padding bytes are filled with random value (Shown as 0xand 0xbb above).
func ISO10126Pad(blockSize int, data []byte) ([]byte, error) {
	// Restrict block size from 1 to 255 because the range of byte is 0-255.
	if blockSize < 1 || blockSize > 255 {
		return nil, &er.Error{
			Package:     ErrPkg,
			Type:        ErrTypePadding,
			Description: ErrDscPadding,
			Detail:      "ISO10126 invalid block size",
		}
	}

	n := blockSize - len(data)%blockSize
	pad := make([]byte, n)
	if _, err := io.ReadFull(rand.Reader, pad); err != nil {
		return nil, (&er.Error{
			Package:     ErrPkg,
			Type:        ErrTypePadding,
			Description: ErrDscPadding,
			Detail:      "ISO10126",
		}).Wrap(err)
	}
	pad[n-1] = byte(n)

	data = append(data, pad...)
	return data, nil
}

// ISO10126UnPad removes padding using ISO10126 padding.
func ISO10126UnPad(blockSize int, data []byte) ([]byte, error) {
	// Restrict block size from 1 to 255 because the range of byte is 0-255.
	if blockSize < 1 || blockSize > 255 {
		return nil, &er.Error{
			Package:     ErrPkg,
			Type:        ErrTypePadding,
			Description: ErrDscUnpadding,
			Detail:      "ISO10126 invalid block size",
		}
	}

	// Check that the given data has valid length as padded data.
	n := len(data)
	if n < blockSize || n%blockSize != 0 {
		return nil, &er.Error{
			Package:     ErrPkg,
			Type:        ErrTypePadding,
			Description: ErrDscUnpadding,
			Detail:      "ISO10126 invalid data length",
		}
	}

	pad := int(data[n-1])
	if n < pad {
		return nil, &er.Error{
			Package:     ErrPkg,
			Type:        ErrTypePadding,
			Description: ErrDscUnpadding,
			Detail:      "ISO10126 invalid padding size",
		}
	}

	return data[:n-pad], nil
}
