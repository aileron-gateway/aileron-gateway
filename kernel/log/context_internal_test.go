// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package log

import (
	"context"
	"testing"

	"github.com/aileron-gateway/aileron-gateway/internal/testutil"
)

type testAttribute struct {
	M  map[string]any
	KV []any
	ID string // To make the instance comparative in the tests.
}

func (a *testAttribute) String() string {
	return a.ID
}

func (a *testAttribute) Name() string {
	return a.ID
}

func (a *testAttribute) Map() map[string]any {
	return a.M
}

func (a *testAttribute) KeyValues() []any {
	return a.KV
}

func TestContextWithAttrs(t *testing.T) {
	type condition struct {
		ctx   context.Context
		attrs []Attributes
	}

	type action struct {
		attrs []Attributes
	}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"input nil",
			&condition{
				ctx:   context.Background(),
				attrs: nil,
			},
			&action{
				attrs: []Attributes{},
			},
		),
		gen(
			"add one",
			&condition{
				ctx: context.Background(),
				attrs: []Attributes{
					&testAttribute{ID: "test"},
				},
			},
			&action{
				attrs: []Attributes{
					&testAttribute{ID: "test"},
				},
			},
		),
		gen(
			"add multiple",
			&condition{
				ctx: context.Background(),
				attrs: []Attributes{
					&testAttribute{ID: "test1"},
					&testAttribute{ID: "test2"},
				},
			},
			&action{
				attrs: []Attributes{
					&testAttribute{ID: "test1"},
					&testAttribute{ID: "test2"},
				}},
		),
		gen(
			"append",
			&condition{
				ctx: context.WithValue(context.Background(), attrsContextKey, []Attributes{&testAttribute{ID: "test1"}}),
				attrs: []Attributes{
					&testAttribute{ID: "test2"},
					&testAttribute{ID: "test3"},
				},
			},
			&action{
				attrs: []Attributes{
					&testAttribute{ID: "test1"},
					&testAttribute{ID: "test2"},
					&testAttribute{ID: "test3"},
				},
			},
		),
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			ctx := ContextWithAttrs(tt.C.ctx, tt.C.attrs...)
			attrs := ctx.Value(attrsContextKey).([]Attributes)

			testutil.Diff(t, tt.A.attrs, attrs)
		})
	}
}

func TestAttrsFromContext(t *testing.T) {
	type condition struct {
		ctx context.Context
	}

	type action struct {
		attrs []Attributes
	}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"no attributes",
			&condition{
				ctx: context.Background(),
			},
			&action{
				attrs: nil,
			},
		),
		gen(
			"append",
			&condition{
				ctx: context.WithValue(context.Background(), attrsContextKey,
					[]Attributes{
						&testAttribute{ID: "test1"},
						&testAttribute{ID: "test2"},
					},
				),
			},
			&action{
				attrs: []Attributes{
					&testAttribute{ID: "test1"},
					&testAttribute{ID: "test2"},
				},
			},
		),
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			attrs := AttrsFromContext(tt.C.ctx)
			testutil.Diff(t, tt.A.attrs, attrs)
		})
	}
}
