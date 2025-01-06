package log

import (
	"context"
	"testing"

	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
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

	cndInputNil := "input nil"
	cndAddOne := "add one attribute"
	cndAddMultiple := "add multiple attributes"
	cndAppend := "append attributes"
	actCheckNil := "check nil"
	actCheckAttributes := "check attributes"

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	tb.Condition(cndInputNil, "input nil slice as an attributes")
	tb.Condition(cndAddOne, "input 1 attribute with a fresh context")
	tb.Condition(cndAddMultiple, "input multiple attributes with a fresh context")
	tb.Condition(cndAppend, "input attributes with a context which already has some attributes")
	tb.Action(actCheckNil, "check that nil was returned")
	tb.Action(actCheckAttributes, "check that the attributes was properly saved in the context")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"input nil",
			[]string{cndInputNil},
			[]string{actCheckNil},
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
			[]string{cndAddOne},
			[]string{actCheckAttributes},
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
			[]string{cndAddMultiple},
			[]string{actCheckAttributes},
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
			[]string{cndAppend},
			[]string{actCheckAttributes},
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

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			ctx := ContextWithAttrs(tt.C().ctx, tt.C().attrs...)
			attrs := ctx.Value(attrsContextKey).([]Attributes)

			testutil.Diff(t, tt.A().attrs, attrs)
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

	cndNoAttributes := "no attributes"
	cndSomeAttributes := "some attributes"
	actCheckNil := "check nil"
	actCheckAttributes := "check line"

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	tb.Condition(cndNoAttributes, "no attributes in the input context")
	tb.Condition(cndSomeAttributes, "some attributes in the input context")
	tb.Action(actCheckNil, "check that nil was returned")
	tb.Action(actCheckAttributes, "check that the attributes was properly obtained from the context")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"no attributes",
			[]string{cndNoAttributes},
			[]string{actCheckNil},
			&condition{
				ctx: context.Background(),
			},
			&action{
				attrs: nil,
			},
		),
		gen(
			"append",
			[]string{cndSomeAttributes},
			[]string{actCheckAttributes},
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

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			attrs := AttrsFromContext(tt.C().ctx)
			testutil.Diff(t, tt.A().attrs, attrs)
		})
	}
}
