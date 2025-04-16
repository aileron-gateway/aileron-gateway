// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package er

import (
	"errors"
	"strings"
)

// Error is a general error object.
// Error instance is considered to be same
// when they have the same Package, Type and Description.
type Error struct {
	// inner is the inner error error.
	// inner is ignored when comparing the Error.
	inner error

	// Package is the package name
	// that this error belongs to.
	// This field should not be empty.
	// Package is compared in Is method.
	Package string

	// Type is the error type of this error.
	// This field should not be empty.
	// Type is compared in Is method.
	Type string

	// Description is the description type of this error.
	// This field should not be empty.
	// Description is compared in Is method.
	Description string

	// Detail is the detail about this error.
	// Detail is ignored when comparing the Error.
	Detail string
}

func (e *Error) Wrap(err error) error {
	if e.inner == nil {
		e.inner = err
	} else {
		e.inner = errors.Join(err, e.inner)
	}
	return e
}

func (e *Error) Unwrap() error {
	return e.inner
}

func (e *Error) Is(err error) bool {
	if err == nil {
		return e == nil
	}
	for {
		ee, ok := err.(*Error)
		if ok {
			return e.Package == ee.Package && e.Type == ee.Type && e.Description == ee.Description
		}
		uw, ok := err.(interface{ Unwrap() error })
		if !ok {
			return false
		}
		err = uw.Unwrap()
	}
}

func (e *Error) Error() string {
	var b strings.Builder
	b.Grow(100)
	b.WriteString(e.Package)
	b.WriteString(": ")
	b.WriteString(e.Type)
	b.WriteString(": ")
	b.WriteString(e.Description)
	if e.Detail != "" {
		b.WriteString(" ")
		b.WriteString(e.Detail)
	}
	if e.inner != nil {
		b.WriteString(" [ ")
		b.WriteString(e.inner.Error())
		b.WriteString(" ]")
	}
	return b.String()
}
