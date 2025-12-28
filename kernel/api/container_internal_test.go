// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

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

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"new instance",
			&condition{},
			&action{
				a: &ContainerAPI{
					objStore: map[string]any{},
				},
			},
		),
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			a := NewContainerAPI()
			testutil.Diff(t, tt.A.a, a, cmp.AllowUnexported(ContainerAPI{}))
		})
	}
}
