// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package encoder

import (
	"bytes"
	"encoding/json"

	"github.com/aileron-projects/go/zerrors"
)

// MarshalJSON marshal struct into byte array of json format.
// If nil value is given, nil byte and nil error are returned.
func MarshalJSON(in any) ([]byte, error) {
	if in == nil {
		return nil, nil
	}
	var b bytes.Buffer
	enc := json.NewEncoder(&b)
	enc.SetIndent("", "  ")
	if err := enc.Encode(in); err != nil {
		return nil, zerrors.NewErr(err, "internal/encoder: marshaling from any to json failed.", "")
	}
	return b.Bytes(), nil
}

// UnmarshalJSON unmarshal byte array of json into the given struct.
// If nil byte is given, this function do nothing and return nil.
func UnmarshalJSON(in []byte, into any) error {
	if into == nil {
		return nil
	}
	err := json.Unmarshal(in, into)
	if err != nil {
		return zerrors.NewErr(err, "internal/encoder: unmarshaling json failed.", "%s", string(addLineNumber(in)))
	}
	return nil
}
