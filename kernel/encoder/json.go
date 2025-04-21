// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package encoder

import (
	"bytes"
	"encoding/json"
	"os"

	"github.com/aileron-gateway/aileron-gateway/kernel/er"
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
		return nil, (&er.Error{
			Package:     ErrPkg,
			Type:        ErrTypeJSON,
			Description: ErrDscMarshal,
			Detail:      "from any to json",
		}).Wrap(err)
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
		return (&er.Error{
			Package:     ErrPkg,
			Type:        ErrTypeJSON,
			Description: ErrDscUnmarshal,
			Detail:      string(addLineNumber(in)),
		}).Wrap(err)
	}
	return nil
}

// UnmarshalJSONFile reads json file and unmarshal it into struct.
// If the target struct is nil, this method do nothing and return nil.
func UnmarshalJSONFile(path string, into any) error {
	if into == nil {
		return nil
	}
	b, err := os.ReadFile(path)
	if err != nil {
		return (&er.Error{
			Package:     ErrPkg,
			Type:        ErrTypeJSON,
			Description: ErrDscUnmarshal,
			Detail:      path,
		}).Wrap(err)
	}
	if err := UnmarshalJSON(b, into); err != nil {
		return err
	}
	return nil
}
