// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package testutil

type Case[C, A any] struct {
	C    C
	A    A
	Name string
	Cnd  []string
	Act  []string
}

func NewCase[C, A any](name string, cnd []string, act []string, c C, a A) *Case[C, A] {
	return &Case[C, A]{
		Name: name,
		Cnd:  cnd,
		Act:  act,
		C:    c,
		A:    a,
	}
}

func Register[C, A any](table *Table[C, A], cases ...*Case[C, A]) {
	builder := table.EntryBuilder()
	for _, c := range cases {
		e := builder.Reset().Name(c.Name)
		for _, cnd := range c.Cnd {
			e.CTrue(cnd)
		}
		for _, act := range c.Act {
			e.ATrue(act)
		}
		e.Condition(c.C)
		e.Action(c.A)
		e.Build().Register(table)
	}
}
