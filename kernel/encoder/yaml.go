// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package encoder

import (
	"bytes"
	"fmt"
	"os"

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

// UnmarshalYAMLFile reads yaml file and unmarshal it into a struct.
// If given file path is an empty string or target struct is nil,
// This method do nothing and return nil.
func UnmarshalYAMLFile(path string, into any) error {
	if into == nil {
		return nil
	}
	b, err := os.ReadFile(path)
	if err != nil {
		return (&er.Error{
			Package:     ErrPkg,
			Type:        ErrTypeYaml,
			Description: ErrDscUnmarshal,
			Detail:      path,
		}).Wrap(err)
	}
	if err := UnmarshalYAML(b, into); err != nil {
		return err // Return err as-is.
	}
	return nil
}
