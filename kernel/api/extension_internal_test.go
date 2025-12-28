// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package api

import (
	"testing"

	"github.com/aileron-gateway/aileron-gateway/internal/testutil"
	"github.com/aileron-gateway/aileron-gateway/kernel/er"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestNewExtensionAPI(t *testing.T) {
	type condition struct {
	}

	type action struct {
		a *ExtensionAPI
	}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"new instance",
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

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			a := NewExtensionAPI()
			testutil.Diff(t, tt.A.a, a, cmp.AllowUnexported(ExtensionAPI{}))
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

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"register 1 creator",
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

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			a := NewExtensionAPI()

			var err error
			for i := range tt.C.keys {
				err = a.Register(tt.C.keys[i], tt.C.creators[i])
			}

			testutil.Diff(t, tt.A.err, err, cmpopts.EquateErrors())
			testutil.Diff(t, tt.A.creators, a.creators)
		})
	}
}
