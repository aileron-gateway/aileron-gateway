// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package zxml

import (
	"errors"
	"maps"
	"slices"
	"testing"

	"github.com/davecgh/go-spew/spew"
)

// AssertEqual checks if the given two values are the same.
// See also https://go.dev/wiki/TestComments
func AssertEqual[T comparable](t *testing.T, errReason string, want, got T) {
	t.Helper()
	if want == got {
		return
	}
	errReason += "\n"
	errReason += "-want: " + spew.Sdump(want)
	errReason += "+got: " + spew.Sdump(got)
	t.Error(errReason)
}

// AssertNotEqual checks if the given two values are not same.
// See also https://go.dev/wiki/TestComments
func AssertNotEqual[T comparable](t *testing.T, errReason string, check, got T) {
	t.Helper()
	if check != got {
		return
	}
	errReason += "\n"
	errReason += "-check: " + spew.Sdump(check)
	errReason += "+got: " + spew.Sdump(got)
	t.Error(errReason)
}

// AssertEqualSlice checks if the given two slices are the same using [slices.Equal].
// See also https://go.dev/wiki/TestComments
func AssertEqualSlice[S ~[]E, E comparable](t *testing.T, errReason string, want, got S) {
	t.Helper()
	if slices.Equal(want, got) {
		return
	}
	errReason += "\n"
	errReason += "-want: " + spew.Sdump(want)
	errReason += "+got: " + spew.Sdump(got)
	t.Error(errReason)
}

// AssertEqualMap checks if the given two maps are the same using [maps.Equal].
// errReason can be a expression for fmt.Sprintf(errReason, want, got).
// See also https://go.dev/wiki/TestComments
func AssertEqualMap[M1, M2 ~map[K]V, K, V comparable](t *testing.T, errReason string, want M1, got M2) {
	t.Helper()
	if maps.Equal(want, got) {
		return
	}
	errReason += "\n"
	errReason += "-want: " + spew.Sdump(want)
	errReason += "+got: " + spew.Sdump(got)
	t.Error(errReason)
}

// AssertEqualErr checks if the given two errors are the same.
// Errors are checked by following order and considered the same
// when one of them returned true.
//
//   - Compare pointer: want == got
//   - Compare error: errors.Is(got, want)
//   - Compare message: want.Error() == got.Error()
//
// See also https://go.dev/wiki/TestComments
func AssertEqualErr(t *testing.T, errReason string, want, got error) {
	t.Helper()
	if want == got {
		return // nil == nil is also here.
	}
	if errors.Is(got, want) {
		return
	}
	if want != nil && got != nil {
		if want.Error() == got.Error() {
			return
		}
	}
	errReason += "\n"
	errReason += "-want: " + spew.Sdump(want)
	errReason += "+got: " + spew.Sdump(got)
	t.Error(errReason)
}
