// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package uid

import (
	"crypto/rand"
	"encoding/binary"
	"hash/fnv"
	"io"
	"os"
	"time"

	"github.com/aileron-gateway/aileron-gateway/kernel/er"
)

// hostnameFNV1 is the FNV1 hashed hostname with 8 bytes.
// This value is used for generating a unique ID.
var hostnameFNV1 = func() []byte {
	hostname, _ := os.Hostname()
	h := fnv.New64()
	h.Write([]byte(hostname))
	return h.Sum(nil)
}()

// NewID creates a new 30 bytes ID which is intended to be used as session IDs.
// Check the best practice for session ID generation at
// https://cheatsheetseries.owasp.org/cheatsheets/Session_Management_Cheat_Sheet.html.
//
// Generated IDs consists of
//   - 8 bytes unix time in microsecond. Valid until 9,223,372,036,854 unix seconds, or January 10th, 294247.
//   - 22 bytes random fetched from crypt/rand.
//
// ID will have the following format when encoded with the encoder:
//
//   - Base16           : ^[0-9a-fA-F]{60}$
//   - Base32           : ^[2-7A-Z]{48}$
//   - Base32Hex        : ^[0-9A-V]{48}$
//   - Base32Escaped    : ^[0-9B-DF-HJ-NP-TV-Z]{48}$
//   - Base32HexEscaped : ^[0-9B-DF-HJ-NP-TV-Z]{48}$
//   - Base64           : ^[0-9a-zA-Z+/]{40}$
//   - Base64Raw        : ^[0-9a-zA-Z+/]{40}$
//   - Base64URL        : ^[0-9a-zA-Z-_]{40}$
//   - Base64RawURL     : ^[0-9a-zA-Z-_]{40}$
func NewID() ([]byte, error) {
	// x is a 30 bytes unique ID.
	x := [30]byte{}

	// Initial 8 bytes are timestamp.
	//nolint:gosec // G115: integer overflow conversion int64 -> uint64
	binary.BigEndian.PutUint64(x[0:], uint64(time.Now().UnixMicro()))

	// Fill the next 22 bytes with random values.
	if _, err := io.ReadFull(rand.Reader, x[8:]); err != nil {
		return nil, (&er.Error{Package: ErrPkg, Type: "id", Description: ErrDscNew}).Wrap(err)
	}

	// Return encoded ID.
	return x[:], nil
}

// NewHostedID creates a new 30 bytes ID which is intended to be used
// as request IDs and trace IDs.
//
// Generated IDs consist of
//   - 8 bytes unix time in microsecond. Valid until 9,223,372,036,854 unix seconds, or January 10th, 294247.
//   - 8 bytes FNV1 hash of the hostname.
//   - 14 bytes random value read from crypt/rand.
//
// ID will have the following format when encoded with the encoder:
//
//   - Base16           : ^[0-9a-fA-F]{60}$
//   - Base32           : ^[2-7A-Z]{48}$
//   - Base32Hex        : ^[0-9A-V]{48}$
//   - Base32Escaped    : ^[0-9B-DF-HJ-NP-TV-Z]{48}$
//   - Base32HexEscaped : ^[0-9B-DF-HJ-NP-TV-Z]{48}$
//   - Base64           : ^[0-9a-zA-Z+/]{40}$
//   - Base64Raw        : ^[0-9a-zA-Z+/]{40}$
//   - Base64URL        : ^[0-9a-zA-Z-_]{40}$
//   - Base64RawURL     : ^[0-9a-zA-Z-_]{40}$
func NewHostedID() ([]byte, error) {
	// x is a 30 bytes unique ID.
	x := [30]byte{}

	// Initial 8 bytes are timestamp.
	//nolint:gosec // G115: integer overflow conversion int64 -> uint64
	binary.BigEndian.PutUint64(x[0:], uint64(time.Now().UnixMicro()))

	// Next 8 bytes are hash of the hostname.
	copy(x[8:], hostnameFNV1)

	// Last 14 bytes are random.
	if _, err := io.ReadFull(rand.Reader, x[16:]); err != nil {
		return nil, (&er.Error{Package: ErrPkg, Type: "hosted id", Description: ErrDscNew}).Wrap(err)
	}

	// Return encoded ID with given encoder.
	return x[:], nil
}
