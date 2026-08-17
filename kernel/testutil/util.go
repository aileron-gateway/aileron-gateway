// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package testutil

import (
	"errors"
	"io"
	"reflect"
	"regexp"
	"runtime/debug"
	"testing"

	"github.com/google/go-cmp/cmp"
)

// Case keeps conditions and actions.
//
// Deprecated: Do not use this.
// This feature will be removed in the future.
type Case[C, A any] struct {
	Name string // Name is the case name.
	C    C      // C is the conditions.
	A    A      // A is the actions.
}

// NewCase returns new test case.
//
// Deprecated: Do not use this.
// This feature will be removed in the future.
func NewCase[C, A any](name string, c C, a A) *Case[C, A] {
	return &Case[C, A]{
		Name: name,
		C:    c,
		A:    a,
	}
}

// Diff compares two value using go-comp.
//
// Deprecated: Do not use this.
// This feature will be removed in the future.
// Use https://github.com/aileron-projects/go-tester.
func Diff(t *testing.T, want, got any, opts ...cmp.Option) {
	t.Helper()
	// opts = append(opts, cmpopts.EquateEmpty())
	if diff := cmp.Diff(want, got, opts...); diff != "" {
		t.Errorf("mismatch (-want +got):\n%s", diff)
		t.Errorf("StackTrace:\n%s", debug.Stack())
	}
}

// DiffError compares two error.
//
// Deprecated: Do not use this.
// This feature will be removed in the future.
// Use https://github.com/aileron-projects/go-tester.
func DiffError(t *testing.T, want any, pattern *regexp.Regexp, got error, opts ...cmp.Option) {
	t.Helper()
	if want == nil || got == nil {
		Diff(t, want, got, opts...)
		return
	}
	if pattern != nil && !pattern.MatchString(got.Error()) {
		t.Errorf("error message mismatch :\n")
		t.Errorf(" (-want) %s\n", pattern.String())
		t.Errorf(" (+got)  %s\n", got.Error())
	}
	iser, ok := want.(interface{ Is(error) bool })
	if !ok {
		Diff(t, want, got, opts...)
		return
	}
	if !iser.Is(got) {
		t.Errorf("error mismatch :\n")
		t.Errorf(" (-want) %#v\n", want)
		t.Errorf(" (+got)  %#v\n", got)
	}
}

// ComparePointer is the function to compare two value pointers.
// This function is mainly intended to be used as an option of go-cmp.
// For example, use this option
//
//	cmp.Comparer(testutil.ComparePointer[foo.Bar])
//
// Deprecated: Do not use this.
// This feature will be removed in the future.
// Use https://github.com/aileron-projects/go-tester.
func ComparePointer[T any](x, y T) bool {
	return reflect.ValueOf(x).Pointer() == reflect.ValueOf(y).Pointer()
}

// ErrorReader is an io.Reader which returns an error.
// This implements io.Reader interface.
//
// Deprecated: Do not use this.
// This feature will be removed in the future.
// Use https://github.com/aileron-projects/go-tester.
type ErrorReader struct {
	io.Reader
}

func (r *ErrorReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("rand read error")
}
