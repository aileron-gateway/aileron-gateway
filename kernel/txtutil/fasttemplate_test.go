package txtutil_test

import (
	"fmt"
	"testing"

	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
	"github.com/aileron-gateway/aileron-gateway/kernel/txtutil"
)

func ExampleFastTemplate_Execute() {
	val := map[string]any{
		"foo": "alice",
		"bar": "bob",
	}

	tpl := txtutil.NewFastTemplate(`Hello {{foo}} and {{bar}}!`, "{{", "}}")

	fmt.Println(string(tpl.Execute(val)))
	// Output:
	// 	Hello alice and bob!
}

func TestFastTemplate_Execute(t *testing.T) {
	type condition struct {
		tpl   string
		start string
		end   string
		m     map[string]any
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
			"empty template",
			[]string{},
			[]string{},
			&condition{
				tpl:   "",
				start: "{{",
				end:   "}}",
				m:     map[string]any{"test": "foo"},
			},
			&action{
				expect: []byte(nil),
			},
		),
		gen(
			"1 tag only",
			[]string{},
			[]string{},
			&condition{
				tpl:   "{{test}}",
				start: "{{",
				end:   "}}",
				m:     map[string]any{"test": "foo"},
			},
			&action{
				expect: []byte("foo"),
			},
		),
		gen(
			"1 tag",
			[]string{},
			[]string{},
			&condition{
				tpl:   "123 {{test}} 456",
				start: "{{",
				end:   "}}",
				m:     map[string]any{"test": "foo"},
			},
			&action{
				expect: []byte("123 foo 456"),
			},
		),
		gen(
			"2 tag",
			[]string{},
			[]string{},
			&condition{
				tpl:   "123 {{test1}} {{test2}} 456",
				start: "{{",
				end:   "}}",
				m:     map[string]any{"test1": "foo", "test2": "bar"},
			},
			&action{
				expect: []byte("123 foo bar 456"),
			},
		),
		gen(
			"value nil",
			[]string{},
			[]string{},
			&condition{
				tpl:   "123 {{test}} 456",
				start: "{{",
				end:   "}}",
				m:     map[string]any{"test": nil},
			},
			&action{
				expect: []byte("123 <nil> 456"),
			},
		),
		gen(
			"value not found",
			[]string{},
			[]string{},
			&condition{
				tpl:   "123 {{test}} 456",
				start: "{{",
				end:   "}}",
				m:     map[string]any{"not me": "foo"},
			},
			&action{
				expect: []byte("123  456"),
			},
		),
		gen(
			"bracket []",
			[]string{},
			[]string{},
			&condition{
				tpl:   "123[test]456",
				start: "[",
				end:   "]",
				m:     map[string]any{"test": "foo"},
			},
			&action{
				expect: []byte("123foo456"),
			},
		),
		gen(
			"bracket {}",
			[]string{},
			[]string{},
			&condition{
				tpl:   "123{test}456",
				start: "{",
				end:   "}",
				m:     map[string]any{"test": "foo"},
			},
			&action{
				expect: []byte("123foo456"),
			},
		),
		gen(
			"bracket %%",
			[]string{},
			[]string{},
			&condition{
				tpl:   "123%test%456",
				start: "%",
				end:   "%",
				m:     map[string]any{"test": "foo"},
			},
			&action{
				expect: []byte("123foo456"),
			},
		),
		gen(
			"bracket #%",
			[]string{},
			[]string{},
			&condition{
				tpl:   "123#test%456",
				start: "#",
				end:   "%",
				m:     map[string]any{"test": "foo"},
			},
			&action{
				expect: []byte("123foo456"),
			},
		),
		gen(
			"empty tag",
			[]string{},
			[]string{},
			&condition{
				tpl:   "123{{test}}456",
				start: "",
				end:   "",
				m:     map[string]any{"test": "foo"},
			},
			&action{
				expect: []byte("{{foo}}"),
			},
		),
		gen(
			"tag with spaces",
			[]string{},
			[]string{},
			&condition{
				tpl:   "123{ test }456",
				start: "{",
				end:   "}",
				m:     map[string]any{"test": "foo"},
			},
			&action{
				expect: []byte("123foo456"),
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			tpl := txtutil.NewFastTemplate(tt.C().tpl, tt.C().start, tt.C().end)
			result := tpl.Execute(tt.C().m)
			testutil.Diff(t, string(tt.A().expect), string(result))
		})
	}
}

func TestFastTemplate_ExecuteFunc(t *testing.T) {
	type condition struct {
		tpl   string
		start string
		end   string
		m     map[string]any
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
			"empty template",
			[]string{},
			[]string{},
			&condition{
				tpl:   "",
				start: "{{",
				end:   "}}",
				m:     map[string]any{"test": "foo"},
			},
			&action{
				expect: []byte(nil),
			},
		),
		gen(
			"1 tag only",
			[]string{},
			[]string{},
			&condition{
				tpl:   "{{test}}",
				start: "{{",
				end:   "}}",
				m:     map[string]any{"test": "foo"},
			},
			&action{
				expect: []byte("foo"),
			},
		),
		gen(
			"1 tag",
			[]string{},
			[]string{},
			&condition{
				tpl:   "123 {{test}} 456",
				start: "{{",
				end:   "}}",
				m:     map[string]any{"test": "foo"},
			},
			&action{
				expect: []byte("123 foo 456"),
			},
		),
		gen(
			"2 tag",
			[]string{},
			[]string{},
			&condition{
				tpl:   "123 {{test1}} {{test2}} 456",
				start: "{{",
				end:   "}}",
				m:     map[string]any{"test1": "foo", "test2": "bar"},
			},
			&action{
				expect: []byte("123 foo bar 456"),
			},
		),
		gen(
			"value not found",
			[]string{},
			[]string{},
			&condition{
				tpl:   "123 {{test}} 456",
				start: "{{",
				end:   "}}",
				m:     map[string]any{"not me": "foo"},
			},
			&action{
				expect: []byte("123 <nil> 456"),
			},
		),
		gen(
			"bracket []",
			[]string{},
			[]string{},
			&condition{
				tpl:   "123[test]456",
				start: "[",
				end:   "]",
				m:     map[string]any{"test": "foo"},
			},
			&action{
				expect: []byte("123foo456"),
			},
		),
		gen(
			"bracket {}",
			[]string{},
			[]string{},
			&condition{
				tpl:   "123{test}456",
				start: "{",
				end:   "}",
				m:     map[string]any{"test": "foo"},
			},
			&action{
				expect: []byte("123foo456"),
			},
		),
		gen(
			"bracket #%",
			[]string{},
			[]string{},
			&condition{
				tpl:   "123#test%456",
				start: "#",
				end:   "%",
				m:     map[string]any{"test": "foo"},
			},
			&action{
				expect: []byte("123foo456"),
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			tpl := txtutil.NewFastTemplate(tt.C().tpl, tt.C().start, tt.C().end)
			result := tpl.ExecuteFunc(func(tag string) []byte {
				return []byte(fmt.Sprintf("%+v", tt.C().m[tag]))
			})
			testutil.Diff(t, string(tt.A().expect), string(result))
		})
	}
}
