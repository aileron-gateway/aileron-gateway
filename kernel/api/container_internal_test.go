package api

import (
	"testing"

	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
	"github.com/google/go-cmp/cmp"
)

func TestNewContainerAPI(t *testing.T) {
	type condition struct {
	}

	type action struct {
		a *ContainerAPI
	}

	cndNewDefault := "new default"
	actCheckInitialized := "check initialized "

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	tb.Condition(cndNewDefault, "create a new instance")
	tb.Action(actCheckInitialized, "check that the returned instance is initialized with expected values")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"new instance",
			[]string{cndNewDefault},
			[]string{actCheckInitialized},
			&condition{},
			&action{
				a: &ContainerAPI{
					objStore: map[string]any{},
				},
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			a := NewContainerAPI()
			testutil.Diff(t, tt.A().a, a, cmp.AllowUnexported(ContainerAPI{}))
		})
	}
}
