// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package uid

import (
	"context"
	"testing"

	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
)

func TestContextWithID(t *testing.T) {
	type condition struct {
		ctx context.Context
		id  string
	}

	type action struct {
		id string
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	cndInputEmpty := tb.Condition("empty string", "input empty string as the id")
	cndInputNonEmpty := tb.Condition("non-empty string", "input non-empty string as the id")
	cndAlreadyExists := tb.Condition("id already exists", "id is already exists in the given context")
	actCheckStoredID := tb.Action("check stored ID", "check the stored id is the same as expected")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"nil context",
			[]string{},
			[]string{actCheckStoredID},
			&condition{
				ctx: nil,
				id:  "",
			},
			&action{
				id: "",
			},
		),
		gen(
			"empty string",
			[]string{cndInputEmpty},
			[]string{actCheckStoredID},
			&condition{
				ctx: context.Background(),
				id:  "",
			},
			&action{
				id: "",
			},
		),
		gen(
			"non empty string",
			[]string{cndInputNonEmpty},
			[]string{actCheckStoredID},
			&condition{
				ctx: context.Background(),
				id:  "test",
			},
			&action{
				id: "test",
			},
		),
		gen(
			"override",
			[]string{cndAlreadyExists},
			[]string{actCheckStoredID},
			&condition{
				ctx: context.WithValue(context.Background(), idContextKey, "test1"),
				id:  "test2",
			},
			&action{
				id: "test2",
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			ctx := ContextWithID(tt.C().ctx, tt.C().id)
			id := ctx.Value(idContextKey).(string)

			testutil.Diff(t, tt.A().id, id)
		})
	}
}

func TestIDFromContext(t *testing.T) {
	type condition struct {
		ctx context.Context
	}

	type action struct {
		id string
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	cndExists := tb.Condition("id exists", "id exists in the given context")
	actCheckID := tb.Action("check returned ID", "check the returned id is the same as expected")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"nil context",
			[]string{},
			[]string{actCheckID},
			&condition{
				ctx: nil,
			},
			&action{
				id: "",
			},
		),
		gen(
			"id not exists",
			[]string{},
			[]string{actCheckID},
			&condition{
				ctx: context.Background(),
			},
			&action{
				id: "",
			},
		),
		gen(
			"id exists",
			[]string{cndExists},
			[]string{actCheckID},
			&condition{
				ctx: context.WithValue(context.Background(), idContextKey, "test"),
			},
			&action{
				id: "test",
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			id := IDFromContext(tt.C().ctx)
			testutil.Diff(t, tt.A().id, id)
		})
	}
}
