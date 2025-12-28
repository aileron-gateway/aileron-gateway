// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package testutil

type Case[C, A any] struct {
	Name string
	C    C
	A    A
}

// NewCase returns new test case.
//
// Deprecated: Do not use this.
func NewCase[C, A any](name string, c C, a A) *Case[C, A] {
	return &Case[C, A]{
		Name: name,
		C:    c,
		A:    a,
	}
}
