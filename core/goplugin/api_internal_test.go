// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

//go:build (!linux && !darwin && !freebsd) || !cgo

package goplugin

import (
	"regexp"
	"testing"

	"github.com/aileron-gateway/aileron-gateway/core"
	"github.com/aileron-gateway/aileron-gateway/kernel/api"
	"github.com/aileron-gateway/aileron-gateway/kernel/log"
	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
	utilhttp "github.com/aileron-gateway/aileron-gateway/util/http"
	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func TestCreate(t *testing.T) {

	type condition struct {
		manifest protoreflect.ProtoMessage
	}

	type action struct {
		expect     any
		err        any // error or errorutil.Kind
		errPattern *regexp.Regexp
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	testServer := api.NewContainerAPI()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"default manifest",
			[]string{}, []string{},
			&condition{
				manifest: Resource.Default(),
			},
			&action{
				expect:     nil,
				err:        core.ErrCoreGenCreateObject,
				errPattern: regexp.MustCompile(core.ErrPrefix + `failed to create GoPlugin`),
			},
		),
	}
	testutil.Register(table, testCases...)
	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			got, err := Resource.Create(testServer, tt.C().manifest)
			testutil.DiffError(t, tt.A().err, tt.A().errPattern, err)
			opts := []cmp.Option{
				cmp.Comparer(testutil.ComparePointer[log.Logger]),
				cmp.AllowUnexported(utilhttp.DefaultErrorHandler{}),
			}
			testutil.Diff(t, tt.A().expect, got, opts...)
		})
	}

}
