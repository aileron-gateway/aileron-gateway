package http

import (
	"context"
	"net/http"
	"reflect"
	"testing"

	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestContextWithProxyHeader(t *testing.T) {
	type condition struct {
		ctx context.Context
		h   http.Header
	}

	type action struct {
		h http.Header
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"nil context",
			[]string{},
			[]string{},
			&condition{
				ctx: nil,
				h:   http.Header{"foo": []string{"bar"}},
			},
			&action{
				h: http.Header{"foo": []string{"bar"}},
			},
		),
		gen(
			"nil header",
			[]string{},
			[]string{},
			&condition{
				ctx: context.Background(),
				h:   nil,
			},
			&action{
				h: nil,
			},
		),
		gen(
			"empty header",
			[]string{},
			[]string{},
			&condition{
				ctx: context.Background(),
				h:   http.Header{},
			},
			&action{
				h: http.Header{},
			},
		),
		gen(
			"non empty header",
			[]string{},
			[]string{},
			&condition{
				ctx: context.Background(),
				h:   http.Header{"foo": []string{"bar"}},
			},
			&action{
				h: http.Header{"foo": []string{"bar"}},
			},
		),
		gen(
			"context with old header",
			[]string{},
			[]string{},
			&condition{
				ctx: context.WithValue(context.Background(), headerContextKey, http.Header{"test": []string{"value"}}),
				h:   http.Header{"foo": []string{"bar"}},
			},
			&action{
				h: http.Header{"test": []string{"value"}, "foo": []string{"bar"}},
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			ctx := ContextWithProxyHeader(tt.C().ctx, tt.C().h)
			h, _ := ctx.Value(headerContextKey).(http.Header)
			testutil.Diff(t, tt.A().h, h, cmpopts.SortMaps(func(a, b string) bool { return a < b }))
		})
	}
}

func TestProxyHeaderFromContext(t *testing.T) {
	type condition struct {
		ctx context.Context
	}

	type action struct {
		h http.Header
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"nil context",
			[]string{},
			[]string{},
			&condition{
				ctx: nil,
			},
			&action{
				h: nil,
			},
		),
		gen(
			"empty context",
			[]string{},
			[]string{},
			&condition{
				ctx: context.Background(),
			},
			&action{
				h: nil,
			},
		),
		gen(
			"context with header",
			[]string{},
			[]string{},
			&condition{
				ctx: context.WithValue(context.Background(), headerContextKey, http.Header{"foo": []string{"bar"}}),
			},
			&action{
				h: http.Header{"foo": []string{"bar"}},
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			h := ProxyHeaderFromContext(tt.C().ctx)
			testutil.Diff(t, tt.A().h, h, cmpopts.SortMaps(
				func(a, b string) bool { return a < b },
			))
		})
	}
}

func TestContextWithPreProxyHook(t *testing.T) {
	type condition struct {
		ctx context.Context
		h   func(r *http.Request) error
	}

	type action struct {
		hs []func(r *http.Request) error
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	h1 := func(r *http.Request) error { return nil }
	h2 := func(r *http.Request) error { return nil }

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"nil context",
			[]string{},
			[]string{},
			&condition{
				ctx: nil,
				h:   h1,
			},
			&action{
				hs: []func(r *http.Request) error{h1},
			},
		),
		gen(
			"nil func",
			[]string{},
			[]string{},
			&condition{
				ctx: context.Background(),
				h:   nil,
			},
			&action{
				hs: nil,
			},
		),
		gen(
			"non nil func",
			[]string{},
			[]string{},
			&condition{
				ctx: context.Background(),
				h:   h1,
			},
			&action{
				hs: []func(r *http.Request) error{h1},
			},
		),
		gen(
			"context with old header",
			[]string{},
			[]string{},
			&condition{
				ctx: context.WithValue(context.Background(), preProxyHookContextKey, &[]func(r *http.Request) error{h2}),
				h:   h1,
			},
			&action{
				hs: []func(r *http.Request) error{h1, h2},
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			ctx := ContextWithPreProxyHook(tt.C().ctx, tt.C().h)
			hs, _ := ctx.Value(preProxyHookContextKey).(*[]func(r *http.Request) error)
			if tt.A().hs == nil {
				testutil.Diff(t, (*[]func(*http.Request) error)(nil), hs)
			} else {
				opts := []cmp.Option{
					cmp.Comparer(testutil.ComparePointer[func(r *http.Request) error]),
					cmpopts.SortSlices(func(a, b func(r *http.Request) error) bool {
						return reflect.ValueOf(a).Pointer() > reflect.ValueOf(b).Pointer()
					}),
				}
				testutil.Diff(t, tt.A().hs, *hs, opts...)
			}
		})
	}
}

func TestPreProxyHookFromContext(t *testing.T) {
	type condition struct {
		ctx context.Context
	}

	type action struct {
		hs []func(r *http.Request) error
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	h1 := func(r *http.Request) error { return nil }
	h2 := func(r *http.Request) error { return nil }

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"nil context",
			[]string{},
			[]string{},
			&condition{
				ctx: nil,
			},
			&action{
				hs: nil,
			},
		),
		gen(
			"empty context",
			[]string{},
			[]string{},
			&condition{
				ctx: context.Background(),
			},
			&action{
				hs: nil,
			},
		),
		gen(
			"context with old header",
			[]string{},
			[]string{},
			&condition{
				ctx: context.WithValue(context.Background(), preProxyHookContextKey, &[]func(r *http.Request) error{h1, h2}),
			},
			&action{
				hs: []func(r *http.Request) error{h1, h2},
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			hs := PreProxyHookFromContext(tt.C().ctx)
			opts := []cmp.Option{
				cmp.Comparer(testutil.ComparePointer[func(r *http.Request) error]),
				cmpopts.SortSlices(func(a, b func(r *http.Request) error) bool {
					return reflect.ValueOf(a).Pointer() > reflect.ValueOf(b).Pointer()
				}),
			}
			testutil.Diff(t, tt.A().hs, hs, opts...)
		})
	}
}
