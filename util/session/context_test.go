// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package session

import (
	"context"
	"testing"

	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
	"github.com/google/go-cmp/cmp"
)

func TestSessionFromContext(t *testing.T) {
	type condition struct {
		ctx context.Context
	}

	type action struct {
		ss Session
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	testSession := NewDefaultSession(SerializeJSON)
	testSession.Persist("foo", nil)

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"session exists",
			[]string{},
			[]string{},
			&condition{
				ctx: context.WithValue(context.Background(), sessionContextKey, testSession),
			},
			&action{
				ss: testSession,
			},
		),
		gen(
			"session not exists",
			[]string{},
			[]string{},
			&condition{
				ctx: context.Background(),
			},
			&action{
				ss: nil,
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			ss := SessionFromContext(tt.C().ctx)

			opts := []cmp.Option{
				cmp.AllowUnexported(DefaultSession{}),
				cmp.Comparer(testutil.ComparePointer[func(any) ([]byte, error)]),
				cmp.Comparer(testutil.ComparePointer[func([]byte, any) error]),
			}
			testutil.Diff(t, tt.A().ss, ss, opts...)
		})
	}
}

func TestContextWithSession(t *testing.T) {
	type condition struct {
		ctx context.Context
		ss  Session
	}

	type action struct {
		ss Session
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	testSession1 := NewDefaultSession(SerializeJSON)
	testSession1.Persist("foo", nil)
	testSession2 := NewDefaultSession(SerializeJSON)
	testSession2.Persist("bar", nil)

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"session saved",
			[]string{},
			[]string{},
			&condition{
				ctx: context.Background(),
				ss:  testSession1,
			},
			&action{
				ss: testSession1,
			},
		),
		gen(
			"session override",
			[]string{},
			[]string{},
			&condition{
				ctx: context.WithValue(context.Background(), sessionContextKey, testSession1),
				ss:  testSession2,
			},
			&action{
				ss: testSession2,
			},
		),
		gen(
			"nil, session",
			[]string{},
			[]string{},
			&condition{
				ctx: context.Background(),
				ss:  nil,
			},
			&action{
				ss: nil,
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			ctx := ContextWithSession(tt.C().ctx, tt.C().ss)
			ss := SessionFromContext(ctx)

			opts := []cmp.Option{
				cmp.AllowUnexported(DefaultSession{}),
				cmp.Comparer(testutil.ComparePointer[func(any) ([]byte, error)]),
				cmp.Comparer(testutil.ComparePointer[func([]byte, any) error]),
			}
			testutil.Diff(t, tt.A().ss, ss, opts...)
		})
	}
}
