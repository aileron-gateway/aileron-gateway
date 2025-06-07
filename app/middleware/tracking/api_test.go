// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package tracking

import (
	"regexp"
	"testing"

	v1 "github.com/aileron-gateway/aileron-gateway/apis/app/v1"
	k "github.com/aileron-gateway/aileron-gateway/apis/kernel"
	"github.com/aileron-gateway/aileron-gateway/kernel/api"
	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func TestCreate(t *testing.T) {
	type condition struct {
		manifest protoreflect.ProtoMessage
	}

	type action struct {
		err        any // error or errorutil.Kind
		errPattern *regexp.Regexp
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"create with default manifest",
			[]string{},
			[]string{},
			&condition{
				manifest: &v1.TrackingMiddleware{
					APIVersion: apiVersion,
					Kind:       kind,
					Metadata: &k.Metadata{
						Namespace: "",
						Name:      "default",
					},
					Spec: &v1.TrackingMiddlewareSpec{},
				},
			},
			&action{
				err: nil,
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			server := api.NewContainerAPI()

			a := &API{}
			_, err := a.Create(server, tt.C().manifest)
			testutil.DiffError(t, tt.A().err, tt.A().errPattern, err)
		})
	}
}
