package api

import (
	"testing"

	"github.com/aileron-gateway/aileron-gateway/kernel/er"
	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestNewExtensionAPI(t *testing.T) {
	type condition struct {
	}

	type action struct {
		a *ExtensionAPI
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
				a: &ExtensionAPI{
					manifestStore: map[string]any{},
					objStore:      map[string]any{},
					creators:      map[string]Creator{},
					formatStore:   map[string]Format{},
				},
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			a := NewExtensionAPI()
			testutil.Diff(t, tt.A().a, a, cmp.AllowUnexported(ExtensionAPI{}))
		})
	}
}

type stringCreator string

func (c stringCreator) Create(a API[*Request, *Response], f Format, manifest any) (any, error) {
	return c, nil
}

func TestExtensionAPI_Register(t *testing.T) {
	type condition struct {
		keys     []string
		creators []Creator
	}

	type action struct {
		creators map[string]Creator
		err      error
	}

	cndRegisterOne := "1 creator"
	cndRegisterMultiple := "multiple creators"
	cndRegisterNil := "register nil"
	cndRegisterDuplicateKey := "duplicate key"
	actCheckRegistered := "check registered creators"
	actCheckNoError := "no error"
	actCheckError := "non-nil error"

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	tb.Condition(cndRegisterOne, "register 1 non-nil creator")
	tb.Condition(cndRegisterMultiple, "register multiple non-nil creator with different keys")
	tb.Condition(cndRegisterNil, "try to register nil creator")
	tb.Condition(cndRegisterDuplicateKey, "try to register creators with the same key")
	tb.Action(actCheckRegistered, "check that the registered creators are the same as expected")
	tb.Action(actCheckNoError, "check that there is no error")
	tb.Action(actCheckError, "check that a non-nil error was returned")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"register 1 creator",
			[]string{cndRegisterOne},
			[]string{actCheckRegistered, actCheckNoError},
			&condition{
				keys:     []string{"test"},
				creators: []Creator{stringCreator("foo")},
			},
			&action{
				creators: map[string]Creator{
					"test": stringCreator("foo"),
				},
			},
		),
		gen(
			"register multiple creators",
			[]string{cndRegisterMultiple},
			[]string{actCheckRegistered, actCheckNoError},
			&condition{
				keys:     []string{"test1", "test2"},
				creators: []Creator{stringCreator("foo"), stringCreator("bar")},
			},
			&action{
				creators: map[string]Creator{
					"test1": stringCreator("foo"),
					"test2": stringCreator("bar"),
				},
			},
		),
		gen(
			"register nil",
			[]string{cndRegisterNil},
			[]string{actCheckRegistered, actCheckNoError},
			&condition{
				keys:     []string{"test"},
				creators: []Creator{nil},
			},
			&action{
				creators: map[string]Creator{},
			},
		),
		gen(
			"duplicate key",
			[]string{cndRegisterDuplicateKey},
			[]string{actCheckRegistered, actCheckError},
			&condition{
				keys:     []string{"test", "test"},
				creators: []Creator{stringCreator("foo"), stringCreator("bar")},
			},
			&action{
				creators: map[string]Creator{
					"test": stringCreator("foo"),
				},
				err: &er.Error{
					Package:     ErrPkg,
					Type:        ErrTypeExt,
					Description: ErrDscDuplicateKey,
				},
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			a := NewExtensionAPI()

			var err error
			for i := range tt.C().keys {
				err = a.Register(tt.C().keys[i], tt.C().creators[i])
			}

			testutil.Diff(t, tt.A().err, err, cmpopts.EquateErrors())
			testutil.Diff(t, tt.A().creators, a.creators)
		})
	}
}
