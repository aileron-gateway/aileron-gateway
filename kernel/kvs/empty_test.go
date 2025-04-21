// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package kvs_test

import (
	"context"
	"testing"
	"time"

	"github.com/aileron-gateway/aileron-gateway/kernel/kvs"
	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
	"github.com/google/go-cmp/cmp/cmpopts"
)

var (
	_   = kvs.Client[string, string](&kvs.MapKVS[string, string]{})
	ctx = context.Background()
)

func TestMapKVS(t *testing.T) {
	type condition struct {
		s     *kvs.MapKVS[string, string]
		key   string
		value string
	}

	type action struct {
		value string
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"non empty key",
			[]string{},
			[]string{},
			&condition{
				s:     &kvs.MapKVS[string, string]{},
				key:   "foo",
				value: "bar",
			},
			&action{
				value: "bar",
			},
		),
		gen(
			"empty key",
			[]string{},
			[]string{},
			&condition{
				s:     &kvs.MapKVS[string, string]{},
				key:   "",
				value: "bar",
			},
			&action{
				value: "bar",
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			store := tt.C().s

			var err error
			err = store.Open(ctx)
			testutil.Diff(t, nil, err)

			var v string
			_, err = store.Get(ctx, tt.C().key)
			testutil.Diff(t, kvs.Nil, err, cmpopts.EquateErrors())

			err = store.Set(ctx, tt.C().key, tt.C().value)
			testutil.Diff(t, nil, err)

			ok := store.Exists(ctx, tt.C().key)
			testutil.Diff(t, true, ok)
			v, err = store.Get(ctx, tt.C().key)
			testutil.Diff(t, nil, err)
			testutil.Diff(t, tt.A().value, v)

			err = store.Delete(ctx, tt.C().key)
			testutil.Diff(t, nil, err)
			_, err = store.Get(ctx, tt.C().key)
			testutil.Diff(t, kvs.Nil, err, cmpopts.EquateErrors())

			err = store.Close(ctx)
			testutil.Diff(t, nil, err)
		})
	}
}

func TestMapKVS_SetWithTTL(t *testing.T) {
	type condition struct {
		s     *kvs.MapKVS[string, string]
		key   string
		value string
		ttl   time.Duration
	}

	type action struct {
		deleted bool
		value   string
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"no ttl",
			[]string{},
			[]string{},
			&condition{
				s:     &kvs.MapKVS[string, string]{},
				key:   "foo",
				value: "bar",
				ttl:   0,
			},
			&action{
				deleted: false,
				value:   "bar",
			},
		),
		gen(
			"with ttl",
			[]string{},
			[]string{},
			&condition{
				s:     &kvs.MapKVS[string, string]{},
				key:   "",
				value: "bar",
				ttl:   100 * time.Millisecond,
			},
			&action{
				deleted: true,
				value:   "bar",
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			store := tt.C().s

			var err error
			err = store.Open(ctx)
			testutil.Diff(t, nil, err)

			err = store.SetWithTTL(ctx, tt.C().key, tt.C().value, tt.C().ttl)
			testutil.Diff(t, nil, err)

			ok := store.Exists(ctx, tt.C().key)
			testutil.Diff(t, true, ok)

			// Key exists after ttl reset.
			time.Sleep(tt.C().ttl / 2)
			store.SetWithTTL(ctx, tt.C().key, tt.C().value, tt.C().ttl)
			time.Sleep(tt.C().ttl / 2)
			testutil.Diff(t, true, store.Exists(ctx, tt.C().key))

			// After ttl.
			time.Sleep(2 * tt.C().ttl)
			ok = store.Exists(ctx, tt.C().key)
			testutil.Diff(t, tt.A().deleted, !ok)

			err = store.Delete(ctx, tt.C().key)
			testutil.Diff(t, nil, err)
			err = store.Close(ctx)
			testutil.Diff(t, nil, err)
		})
	}
}
