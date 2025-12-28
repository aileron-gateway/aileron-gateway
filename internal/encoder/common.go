// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package encoder

import (
	"fmt"
	"regexp"
)

// MarshalFunc is the function type that
// marshals given object into byte slice.
type MarshalFunc func(in any) (b []byte, err error)

// UnmarshalFunc is the function type that
// un-marshals byte slice into given any object.
type UnmarshalFunc func(in []byte, into any) error

// UnmarshalAnotherFunc is the function type that
// un-marshals given object into another object.
type UnmarshalAnotherFunc func(in any, into any) error

// addLineNumber adds line number to the input content.
func addLineNumber(in []byte) []byte {
	re := regexp.MustCompile("(?m)^")
	row := 0
	b := re.ReplaceAllFunc(in, func(b []byte) []byte { row += 1; return []byte(fmt.Sprintf("%04d|", row)) })
	line := []byte("-----1---------1---------1---------1---------")
	pre := append(append([]byte("\n\n"), line...), []byte("\n")...)
	post := append(append([]byte("\n"), line...), []byte("\n\n")...)
	return append(append(pre, b...), post...)
}
