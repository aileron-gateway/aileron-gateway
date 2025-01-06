package slogger

import (
	"fmt"
	"log/slog"
	"testing"

	v1 "github.com/aileron-gateway/aileron-gateway/apis/core/v1"
	"github.com/aileron-gateway/aileron-gateway/apis/kernel"
	"github.com/aileron-gateway/aileron-gateway/kernel/er"
	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
	"github.com/aileron-gateway/aileron-gateway/kernel/txtutil"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func Example_replaceAttr() {
	a := slog.AnyValue(
		map[string]any{
			"foo":  "bar",
			"john": "doe",
			"alice": map[string]any{
				"one": "1",
				"two": "2",
				"three": map[string]any{
					"japan": "tokyo",
				},
			},
		},
	)

	replaceFunc := func(s string) string {
		return "xxx"
	}

	f := replaceAttrFunc([]string{"alice", "three", "japan"}, replaceFunc)
	attr, _ := f(slog.Attr{Value: a})
	fmt.Printf("%+v", attr.Value.Any())
	// Output:
	// map[alice:map[one:1 three:map[japan:xxx] two:2] foo:bar john:doe]
}

func TestNewReplacer(t *testing.T) {
	type condition struct {
		specs []*v1.FieldReplacerSpec
		in    slog.Attr
	}

	type action struct {
		result any
		err    error
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"no spec",
			[]string{},
			[]string{},
			&condition{
				specs: []*v1.FieldReplacerSpec{},
				in:    slog.Any("foo", slog.AnyValue("bar")),
			},
			&action{
				result: "bar",
			},
		),
		gen(
			"delete string",
			[]string{},
			[]string{},
			&condition{
				specs: []*v1.FieldReplacerSpec{
					{Field: "foo"},
				},
				in: slog.Any("foo", slog.AnyValue("bar")),
			},
			&action{
				result: nil,
			},
		),
		gen(
			"replace string",
			[]string{},
			[]string{},
			&condition{
				specs: []*v1.FieldReplacerSpec{
					{
						Field: "foo",
						Replacer: &kernel.ReplacerSpec{
							Replacers: &kernel.ReplacerSpec_Fixed{
								Fixed: &kernel.FixedReplacer{
									Value: "xxx",
								},
							},
						},
					},
				},
				in: slog.Any("foo", slog.AnyValue("bar")),
			},
			&action{
				result: "xxx",
			},
		),
		gen(
			"empty field",
			[]string{},
			[]string{},
			&condition{
				specs: []*v1.FieldReplacerSpec{
					{Field: ""},
					{Field: ""},
					{Field: ""},
					{
						Field: "foo",
						Replacer: &kernel.ReplacerSpec{
							Replacers: &kernel.ReplacerSpec_Fixed{
								Fixed: &kernel.FixedReplacer{
									Value: "xxx",
								},
							},
						},
					},
				},
				in: slog.Any("", slog.AnyValue("bar")),
			},
			&action{
				result: "bar",
			},
		),
		gen(
			"invalid replacer",
			[]string{},
			[]string{},
			&condition{
				specs: []*v1.FieldReplacerSpec{
					{
						Field: "foo",
						Replacer: &kernel.ReplacerSpec{
							Replacers: &kernel.ReplacerSpec_Regexp{
								Regexp: &kernel.RegexpReplacer{
									Pattern: "[0-9",
								},
							},
						},
					},
				},
			},
			&action{
				err: &er.Error{
					Package:     txtutil.ErrPkg,
					Type:        txtutil.ErrTypeReplacer,
					Description: txtutil.ErrDscPattern,
				},
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			repl, err := newReplaceFunc(tt.C().specs)
			testutil.Diff(t, tt.A().err, err, cmpopts.EquateErrors())
			if tt.A().err != nil {
				return
			}

			a := repl.replaceAttr(nil, tt.C().in)
			testutil.Diff(t, tt.A().result, a.Value.Any())
		})
	}
}

func TestReplacer_replaceAttr(t *testing.T) {
	type condition struct {
		specs []*v1.FieldReplacerSpec
		in    slog.Attr
	}

	type action struct {
		result any
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"delete string",
			[]string{},
			[]string{},
			&condition{
				specs: []*v1.FieldReplacerSpec{
					{Field: "foo"},
				},
				in: slog.Any("foo", slog.AnyValue("bar")),
			},
			&action{
				result: nil,
			},
		),
		gen(
			"delete map string",
			[]string{},
			[]string{},
			&condition{
				specs: []*v1.FieldReplacerSpec{
					{Field: "foo.bar"},
				},
				in: slog.Any("foo",
					slog.AnyValue(map[string]any{"bar": "xxx", "baz": "yyy"}),
				),
			},
			&action{
				result: map[string]any{"baz": "yyy"},
			},
		),
		gen(
			"delete map inner string",
			[]string{},
			[]string{},
			&condition{
				specs: []*v1.FieldReplacerSpec{
					{Field: "foo.bar.baz"},
				},
				in: slog.Any("foo",
					slog.AnyValue(map[string]any{"bar": map[string]any{"baz": "yyy", "alice": "bob"}}),
				),
			},
			&action{
				result: map[string]any{"bar": map[string]any{"alice": "bob"}},
			},
		),
		gen(
			"delete int",
			[]string{},
			[]string{},
			&condition{
				specs: []*v1.FieldReplacerSpec{
					{Field: "foo"},
				},
				in: slog.Any("foo", slog.AnyValue(123)),
			},
			&action{
				result: nil,
			},
		),
		gen(
			"delete map int",
			[]string{},
			[]string{},
			&condition{
				specs: []*v1.FieldReplacerSpec{
					{Field: "foo.bar"},
				},
				in: slog.Any("foo",
					slog.AnyValue(map[string]any{"bar": 123, "baz": "yyy"}),
				),
			},
			&action{
				result: map[string]any{"baz": "yyy"},
			},
		),
		gen(
			"delete map inner int",
			[]string{},
			[]string{},
			&condition{
				specs: []*v1.FieldReplacerSpec{
					{Field: "foo.bar.baz"},
				},
				in: slog.Any("foo",
					slog.AnyValue(map[string]any{"bar": map[string]any{"baz": 123, "alice": "bob"}}),
				),
			},
			&action{
				result: map[string]any{"bar": map[string]any{"alice": "bob"}},
			},
		),
		gen(
			"string not match",
			[]string{},
			[]string{},
			&condition{
				specs: []*v1.FieldReplacerSpec{
					{Field: "piyo"},
				},
				in: slog.Any("foo", slog.AnyValue("bar")),
			},
			&action{
				result: "bar",
			},
		),
		gen(
			"map string not match",
			[]string{},
			[]string{},
			&condition{
				specs: []*v1.FieldReplacerSpec{
					{Field: "foo.piyo"},
				},
				in: slog.Any("foo",
					slog.AnyValue(map[string]any{"bar": "xxx", "baz": "yyy"}),
				),
			},
			&action{
				result: map[string]any{"bar": "xxx", "baz": "yyy"},
			},
		),
		gen(
			"map inner string not match",
			[]string{},
			[]string{},
			&condition{
				specs: []*v1.FieldReplacerSpec{
					{Field: "foo.bar.piyo"},
				},
				in: slog.Any("foo",
					slog.AnyValue(map[string]any{"bar": map[string]any{"baz": "yyy", "alice": "bob"}}),
				),
			},
			&action{
				result: map[string]any{"bar": map[string]any{"baz": "yyy", "alice": "bob"}},
			},
		),
		gen(
			"replace string",
			[]string{},
			[]string{},
			&condition{
				specs: []*v1.FieldReplacerSpec{
					{
						Field: "foo",
						Replacer: &kernel.ReplacerSpec{
							Replacers: &kernel.ReplacerSpec_Fixed{
								Fixed: &kernel.FixedReplacer{Value: "xxx"},
							},
						},
					},
				},
				in: slog.Any("foo", slog.AnyValue("bar")),
			},
			&action{
				result: "xxx",
			},
		),
		gen(
			"replace map string",
			[]string{},
			[]string{},
			&condition{
				specs: []*v1.FieldReplacerSpec{
					{
						Field: "foo.bar",
						Replacer: &kernel.ReplacerSpec{
							Replacers: &kernel.ReplacerSpec_Fixed{
								Fixed: &kernel.FixedReplacer{Value: "xxx"},
							},
						},
					},
				},
				in: slog.Any("foo",
					slog.AnyValue(map[string]any{"bar": "xxx", "baz": "yyy"}),
				),
			},
			&action{
				result: map[string]any{"bar": "xxx", "baz": "yyy"},
			},
		),
		gen(
			"replace map inner string",
			[]string{},
			[]string{},
			&condition{
				specs: []*v1.FieldReplacerSpec{
					{
						Field: "foo.bar.baz",
						Replacer: &kernel.ReplacerSpec{
							Replacers: &kernel.ReplacerSpec_Fixed{
								Fixed: &kernel.FixedReplacer{Value: "xxx"},
							},
						},
					},
				},
				in: slog.Any("foo",
					slog.AnyValue(map[string]any{"bar": map[string]any{"baz": "yyy", "alice": "bob"}}),
				),
			},
			&action{
				result: map[string]any{"bar": map[string]any{"baz": "xxx", "alice": "bob"}},
			},
		),
		gen(
			"replace int",
			[]string{},
			[]string{},
			&condition{
				specs: []*v1.FieldReplacerSpec{
					{
						Field: "foo",
						Replacer: &kernel.ReplacerSpec{
							Replacers: &kernel.ReplacerSpec_Fixed{
								Fixed: &kernel.FixedReplacer{Value: "xxx"},
							},
						},
					},
				},
				in: slog.Any("foo", slog.AnyValue(123)),
			},
			&action{
				result: int64(123),
			},
		),
		gen(
			"replace map int",
			[]string{},
			[]string{},
			&condition{
				specs: []*v1.FieldReplacerSpec{
					{
						Field: "foo.bar",
						Replacer: &kernel.ReplacerSpec{
							Replacers: &kernel.ReplacerSpec_Fixed{
								Fixed: &kernel.FixedReplacer{Value: "xxx"},
							},
						},
					},
				},
				in: slog.Any("foo",
					slog.AnyValue(map[string]any{"bar": 123, "baz": "yyy"}),
				),
			},
			&action{
				result: map[string]any{"bar": 123, "baz": "yyy"},
			},
		),
		gen(
			"replace map inner int",
			[]string{},
			[]string{},
			&condition{
				specs: []*v1.FieldReplacerSpec{
					{
						Field: "foo.bar.baz",
						Replacer: &kernel.ReplacerSpec{
							Replacers: &kernel.ReplacerSpec_Fixed{
								Fixed: &kernel.FixedReplacer{Value: "xxx"},
							},
						},
					},
				},
				in: slog.Any("foo",
					slog.AnyValue(map[string]any{"bar": map[string]any{"baz": 123, "alice": "bob"}}),
				),
			},
			&action{
				result: map[string]any{"bar": map[string]any{"baz": 123, "alice": "bob"}},
			},
		),
		gen(
			"cannot reach final value",
			[]string{},
			[]string{},
			&condition{
				specs: []*v1.FieldReplacerSpec{
					{
						Field: "foo.bar.baz",
						Replacer: &kernel.ReplacerSpec{
							Replacers: &kernel.ReplacerSpec_Fixed{
								Fixed: &kernel.FixedReplacer{Value: "xxx"},
							},
						},
					},
				},
				in: slog.Any("foo",
					slog.AnyValue(map[string]any{"bar": 123}),
				),
			},
			&action{
				result: map[string]any{"bar": 123},
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			repl, err := newReplaceFunc(tt.C().specs)
			testutil.Diff(t, nil, err)

			a := repl.replaceAttr(nil, tt.C().in)
			testutil.Diff(t, tt.A().result, a.Value.Any())
		})
	}
}
