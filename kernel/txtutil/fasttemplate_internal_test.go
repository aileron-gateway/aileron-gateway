package txtutil

import (
	"testing"

	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
)

func TestMapValue(t *testing.T) {
	type condition struct {
		tag string
		m   map[string]any
	}

	type action struct {
		expect []byte
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"nil",
			[]string{},
			[]string{},
			&condition{
				tag: "test",
				m:   nil,
			},
			&action{
				expect: []byte(nil),
			},
		),
		gen(
			"string",
			[]string{},
			[]string{},
			&condition{
				tag: "test",
				m:   map[string]any{"test": "foo"},
			},
			&action{
				expect: []byte("foo"),
			},
		),
		gen(
			"string/empty tag",
			[]string{},
			[]string{},
			&condition{
				tag: "",
				m:   map[string]any{"test": "foo"},
			},
			&action{
				expect: []byte(nil),
			},
		),
		gen(
			"[]byte",
			[]string{},
			[]string{},
			&condition{
				tag: "test",
				m:   map[string]any{"test": []byte("foo")},
			},
			&action{
				expect: []byte("foo"),
			},
		),
		gen(
			"int",
			[]string{},
			[]string{},
			&condition{
				tag: "test",
				m:   map[string]any{"test": int(123)},
			},
			&action{
				expect: []byte("123"),
			},
		),
		gen(
			"int32",
			[]string{},
			[]string{},
			&condition{
				tag: "test",
				m:   map[string]any{"test": int32(123)},
			},
			&action{
				expect: []byte("123"),
			},
		),
		gen(
			"int64",
			[]string{},
			[]string{},
			&condition{
				tag: "test",
				m:   map[string]any{"test": int64(123)},
			},
			&action{
				expect: []byte("123"),
			},
		),
		gen(
			"float32",
			[]string{},
			[]string{},
			&condition{
				tag: "test",
				m:   map[string]any{"test": float32(123.456)},
			},
			&action{
				expect: []byte("123.456"),
			},
		),
		gen(
			"float64",
			[]string{},
			[]string{},
			&condition{
				tag: "test",
				m:   map[string]any{"test": float64(123.456)},
			},
			&action{
				expect: []byte("123.456"),
			},
		),
		gen(
			"unsupported complex",
			[]string{},
			[]string{},
			&condition{
				tag: "test",
				m:   map[string]any{"test": complex64(1)},
			},
			&action{
				expect: []byte("(1+0i)"),
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			mv := mapVal(tt.C().m)
			result := mv.Value(tt.C().tag)
			testutil.Diff(t, tt.A().expect, result)
		})
	}
}
