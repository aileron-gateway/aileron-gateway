// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package encoder

import (
	"bytes"
	"fmt"

	"github.com/aileron-gateway/aileron-gateway/kernel/er"
	"gopkg.in/yaml.v3"
)

// MarshalYAML marshal struct into byte array of yaml format.
// If nil value is given, nil byte and nil error are returned.
func MarshalYAML(in any) (b []byte, err error) {
	if in == nil {
		return nil, nil
	}
	// Recover panic of enc.Encode if any.
	defer func() {
		if e := recover(); e != nil {
			err = (&er.Error{
				Package:     ErrPkg,
				Type:        ErrTypeYaml,
				Description: ErrDscMarshal,
				Detail:      "from any to yaml",
			}).Wrap(fmt.Errorf("%v", e))
		}
	}()
	var buf bytes.Buffer
	enc := yaml.NewEncoder(&buf)
	enc.SetIndent(2)
	err = enc.Encode(in) // Not return an error but panics.
	enc.Close()
	return buf.Bytes(), err
}

// UnmarshalYAML unmarshal byte array of yaml into the given struct.
// If nil byte is given, this function do nothing and return nil.
func UnmarshalYAML(in []byte, into any) error {
	if into == nil {
		return nil
	}
	err := yaml.Unmarshal(in, into)
	if err != nil {
		return (&er.Error{
			Package:     ErrPkg,
			Type:        ErrTypeYaml,
			Description: ErrDscUnmarshal,
			Detail:      string(addLineNumber(in)),
		}).Wrap(err)
	}
	return nil
}
