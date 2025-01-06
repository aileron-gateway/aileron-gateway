package errorutil

import (
	"bytes"
	"testing"

	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
)

func TestFastTemplate_Execute(t *testing.T) {
	type condition struct {
		tpl string
		m   map[string]any
	}

	type action struct {
		expect string
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"empty template",
			[]string{},
			[]string{},
			&condition{
				tpl: "",
				m:   map[string]any{"test": "foo"},
			},
			&action{
				expect: "",
			},
		),
		gen(
			"1 tag only",
			[]string{},
			[]string{},
			&condition{
				tpl: "{{test}}",
				m:   map[string]any{"test": "foo"},
			},
			&action{
				expect: "foo",
			},
		),
		gen(
			"1 tag",
			[]string{},
			[]string{},
			&condition{
				tpl: "123 {{test}} 456",
				m:   map[string]any{"test": "foo"},
			},
			&action{
				expect: "123 foo 456",
			},
		),
		gen(
			"2 tag",
			[]string{},
			[]string{},
			&condition{
				tpl: "123 {{test1}} {{test2}} 456",
				m:   map[string]any{"test1": "foo", "test2": "bar"},
			},
			&action{
				expect: "123 foo bar 456",
			},
		),
		gen(
			"value is nil",
			[]string{},
			[]string{},
			&condition{
				tpl: "123 {{test}} 456",
				m:   map[string]any{"test": nil},
			},
			&action{
				expect: "123 <nil> 456",
			},
		),
		gen(
			"value not found",
			[]string{},
			[]string{},
			&condition{
				tpl: "123 {{test}} 456",
				m:   map[string]any{"not me": "foo"},
			},
			&action{
				expect: "123  456",
			},
		),
		gen(
			"tag with space",
			[]string{},
			[]string{},
			&condition{
				tpl: "{{ test }}",
				m:   map[string]any{"test": "foo"},
			},
			&action{
				expect: "foo",
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			var buf bytes.Buffer

			tpl := newTemplate(tt.C().tpl)
			tpl.execute(&buf, tt.C().m)
			testutil.Diff(t, tt.A().expect, buf.String())
		})
	}
}

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
			"nil/empty tag",
			[]string{},
			[]string{},
			&condition{
				tag: "",
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
			result := mv.tagValue(tt.C().tag)
			testutil.Diff(t, tt.A().expect, result)
		})
	}
}
