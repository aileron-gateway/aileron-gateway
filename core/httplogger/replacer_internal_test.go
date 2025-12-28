// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package httplogger

import (
	"net/http"
	"testing"

	v1 "github.com/aileron-gateway/aileron-gateway/apis/core/v1"
	k "github.com/aileron-gateway/aileron-gateway/apis/kernel"
	"github.com/aileron-gateway/aileron-gateway/internal/testutil"
	"github.com/aileron-gateway/aileron-gateway/kernel/er"
	"github.com/aileron-gateway/aileron-gateway/kernel/txtutil"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestHeaderReplacers(t *testing.T) {
	type condition struct {
		specs []*v1.LogHeaderSpec
		data  http.Header
	}

	type action struct {
		data       map[string]string
		allHeaders bool
		err        error
	}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"nil spec",
			&condition{
				specs: nil,
				data:  http.Header{"Alice": {"bob"}, "Foo": {"bar"}},
			},
			&action{
				data: map[string]string{"Alice": "bob", "Foo": "bar"},
				err:  nil,
			},
		),
		gen(
			"empty spec",
			&condition{
				specs: []*v1.LogHeaderSpec{},
				data:  http.Header{"Alice": {"bob"}, "Foo": {"bar"}},
			},
			&action{
				data: map[string]string{"Alice": "bob", "Foo": "bar"},
				err:  nil,
			},
		),
		gen(
			"all headers",
			&condition{
				specs: []*v1.LogHeaderSpec{
					{Name: "*"},
				},
				data: http.Header{"Alice": {"bob"}, "Foo": {"bar"}},
			},
			&action{
				data:       map[string]string{"Alice": "bob", "Foo": "bar"},
				allHeaders: true,
				err:        nil,
			},
		),
		gen(
			"one replacer",
			&condition{
				specs: []*v1.LogHeaderSpec{
					{
						Name: "alice",
						Replacers: []*k.ReplacerSpec{
							{Replacers: &k.ReplacerSpec_Fixed{Fixed: &k.FixedReplacer{Value: "***"}}},
						},
					},
				},
				data: http.Header{"Alice": {"bob"}, "Foo": {"bar"}},
			},
			&action{
				data: map[string]string{
					"Alice": "***",
					"Foo":   "bar",
				},
				err: nil,
			},
		),
		gen(
			"one replacer multi value",
			&condition{
				specs: []*v1.LogHeaderSpec{
					{
						Name: "alice",
						Replacers: []*k.ReplacerSpec{
							{Replacers: &k.ReplacerSpec_Fixed{Fixed: &k.FixedReplacer{Value: "***"}}},
						},
					},
				},
				data: http.Header{"Alice": {"bob", "charlie"}, "Foo": {"bar"}},
			},
			&action{
				data: map[string]string{
					"Alice": "***",
					"Foo":   "bar",
				},
				err: nil,
			},
		),
		gen(
			"two replacer",
			&condition{
				specs: []*v1.LogHeaderSpec{
					{
						Name: "alice",
						Replacers: []*k.ReplacerSpec{
							{Replacers: &k.ReplacerSpec_Fixed{Fixed: &k.FixedReplacer{Value: "***"}}},
						},
					},
					{
						Name: "foo",
						Replacers: []*k.ReplacerSpec{
							{Replacers: &k.ReplacerSpec_Fixed{Fixed: &k.FixedReplacer{Value: "+++"}}},
						},
					},
				},
				data: http.Header{"Alice": {"bob"}, "Foo": {"bar"}},
			},
			&action{
				data: map[string]string{
					"Alice": "***",
					"Foo":   "+++",
				},
				err: nil,
			},
		),
		gen(
			"duplicate names",
			&condition{
				specs: []*v1.LogHeaderSpec{
					{
						Name: "alice",
						Replacers: []*k.ReplacerSpec{
							{Replacers: &k.ReplacerSpec_Fixed{Fixed: &k.FixedReplacer{Value: "***"}}},
						},
					},
					{
						Name: "alice",
						Replacers: []*k.ReplacerSpec{
							{Replacers: &k.ReplacerSpec_Fixed{Fixed: &k.FixedReplacer{Value: "+++"}}},
						},
					},
				},
				data: http.Header{"Alice": {"bob"}, "Foo": {"bar"}},
			},
			&action{
				data: map[string]string{
					"Alice": "+++", // Replace by *** and replace by +++. Finally it becomes +++
					"Foo":   "bar",
				},
				err: nil,
			},
		),
		gen(
			"nil spec",
			&condition{
				specs: []*v1.LogHeaderSpec{
					nil, nil, nil, // nil spec
					{
						Name: "alice",
						Replacers: []*k.ReplacerSpec{
							{Replacers: &k.ReplacerSpec_Fixed{Fixed: &k.FixedReplacer{Value: "***"}}},
						},
					},
				},
				data: http.Header{"Alice": {"bob"}, "Foo": {"bar"}},
			},
			&action{
				data: map[string]string{
					"Alice": "***",
					"Foo":   "bar",
				},
				err: nil,
			},
		),
		gen(
			"invalid spec",
			&condition{
				specs: []*v1.LogHeaderSpec{
					{
						Name: "alice",
						Replacers: []*k.ReplacerSpec{
							{Replacers: &k.ReplacerSpec_Regexp{Regexp: &k.RegexpReplacer{Pattern: "[0-9a-"}}},
						},
					},
				},
				data: http.Header{"Alice": {"bob"}, "Foo": {"bar"}},
			},
			&action{
				data: map[string]string{"Alice": "bob", "Foo": "bar"},
				err: &er.Error{
					Package:     txtutil.ErrPkg,
					Type:        txtutil.ErrTypeReplacer,
					Description: txtutil.ErrDscPattern,
				},
			},
		),
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			replacers, all, err := headerReplacers(tt.C.specs)

			testutil.Diff(t, tt.A.err, err, cmpopts.EquateErrors())
			testutil.Diff(t, tt.A.allHeaders, all)

			bl := &baseLogger{
				allHeaders: true,
				headers:    replacers,
			}
			result := bl.logHeaders(tt.C.data)
			testutil.Diff(t, tt.A.data, result)
		})
	}
}

func TestBodyReplacers(t *testing.T) {
	type condition struct {
		specs []*v1.LogBodySpec
		mime  string
		data  string
	}

	type action struct {
		data string
		err  error
	}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"nil spec",
			&condition{
				specs: nil,
				data:  `{"foo":"bar"}`,
			},
			&action{
				data: `{"foo":"bar"}`,
				err:  nil,
			},
		),
		gen(
			"empty spec",
			&condition{
				specs: []*v1.LogBodySpec{},
				data:  `{"foo":"bar"}`,
			},
			&action{
				data: `{"foo":"bar"}`,
				err:  nil,
			},
		),
		gen(
			"one replacer",
			&condition{
				specs: []*v1.LogBodySpec{
					{
						Mime: "application/json",
						Replacers: []*k.ReplacerSpec{
							{
								Replacers: &k.ReplacerSpec_Value{Value: &k.ValueReplacer{FromTo: map[string]string{"bar": "baz"}}},
							},
						},
					},
				},
				mime: "application/json",
				data: `{"foo":"bar", "alice":"bob"}`,
			},
			&action{
				data: `{"foo":"baz", "alice":"bob"}`,
				err:  nil,
			},
		),
		gen(
			"one replacer/wrong mime",
			&condition{
				specs: []*v1.LogBodySpec{
					{
						Mime: "application/json",
						Replacers: []*k.ReplacerSpec{
							{
								Replacers: &k.ReplacerSpec_Value{Value: &k.ValueReplacer{FromTo: map[string]string{"bar": "baz"}}},
							},
						},
					},
				},
				mime: "text/plain",
				data: `"foo"="bar"`,
			},
			&action{
				data: `"foo"="bar"`,
				err:  nil,
			},
		),
		gen(
			"two replacers",
			&condition{
				specs: []*v1.LogBodySpec{
					{
						Mime: "application/json",
						Replacers: []*k.ReplacerSpec{
							{
								Replacers: &k.ReplacerSpec_Value{
									Value: &k.ValueReplacer{FromTo: map[string]string{"bar": "baz"}},
								},
							},
							{
								Replacers: &k.ReplacerSpec_Value{
									Value: &k.ValueReplacer{FromTo: map[string]string{"bob": "charlie"}},
								},
							},
						},
					},
				},
				mime: "application/json",
				data: `{"foo":"bar", "alice":"bob"}`,
			},
			&action{
				data: `{"foo":"baz", "alice":"charlie"}`,
				err:  nil,
			},
		),
		gen(
			"two replacers different mime",
			&condition{
				specs: []*v1.LogBodySpec{
					{
						Mime: "application/json",
						Replacers: []*k.ReplacerSpec{
							{
								Replacers: &k.ReplacerSpec_Value{
									Value: &k.ValueReplacer{FromTo: map[string]string{"bar": "baz"}},
								},
							},
						},
					},
					{
						Mime: "text/plain",
						Replacers: []*k.ReplacerSpec{
							{
								Replacers: &k.ReplacerSpec_Value{
									Value: &k.ValueReplacer{FromTo: map[string]string{"bob": "charlie"}},
								},
							},
						},
					},
				},
				mime: "text/plain",
				data: `"foo"="bar", "alice"="bob"`,
			},
			&action{
				data: `"foo"="bar", "alice"="charlie"`,
				err:  nil,
			},
		),
		gen(
			"json field replace string",
			&condition{
				specs: []*v1.LogBodySpec{
					{
						Mime: "application/json",
						Replacers: []*k.ReplacerSpec{
							{
								Replacers: &k.ReplacerSpec_Regexp{
									Regexp: &k.RegexpReplacer{Pattern: `\"([^"]*\")`, Replace: `"***"`},
								},
							},
						},
						JSONFields: []string{"foo"},
					},
				},
				mime: "application/json",
				data: `{"foo":"bar", "alice":{"bob", "charlie"}}`,
			},
			&action{
				data: `{"foo":"***", "alice":{"bob", "charlie"}}`,
				err:  nil,
			},
		),
		gen(
			"json field replace object",
			&condition{
				specs: []*v1.LogBodySpec{
					{
						Mime: "application/json",
						Replacers: []*k.ReplacerSpec{
							{
								Replacers: &k.ReplacerSpec_Regexp{
									Regexp: &k.RegexpReplacer{Pattern: `\"([^"]*\")`, Replace: `"***"`},
								},
							},
						},
						JSONFields: []string{"alice"},
					},
				},
				mime: "application/json",
				data: `{"foo":"bar", "alice":{"bob", "charlie"}}`,
			},
			&action{
				data: `{"foo":"bar", "alice":{"***", "***"}}`,
				err:  nil,
			},
		),
		gen(
			"json field not found",
			&condition{
				specs: []*v1.LogBodySpec{
					{
						Mime: "application/json",
						Replacers: []*k.ReplacerSpec{
							{
								Replacers: &k.ReplacerSpec_Regexp{
									Regexp: &k.RegexpReplacer{Pattern: `\"([^"]*\")`, Replace: `"***"`},
								},
							},
						},
						JSONFields: []string{"alice.john"},
					},
				},
				mime: "application/json",
				data: `{"foo":"bar", "alice":{"bob", "charlie"}}`,
			},
			&action{
				data: `{"foo":"bar", "alice":{"bob", "charlie"}}`,
				err:  nil,
			},
		),
		gen(
			"replacer contains nil",
			&condition{
				specs: []*v1.LogBodySpec{
					nil, nil, nil, // nil spec.
					{
						Mime: "application/json",
						Replacers: []*k.ReplacerSpec{
							{
								Replacers: &k.ReplacerSpec_Value{
									Value: &k.ValueReplacer{FromTo: map[string]string{"bar": "baz"}},
								},
							},
						},
					},
				},
				mime: "application/json",
				data: `{"foo":"bar", "alice":"bob"}`,
			},
			&action{
				data: `{"foo":"baz", "alice":"bob"}`,
				err:  nil,
			},
		),
		gen(
			"empty mime",
			&condition{
				specs: []*v1.LogBodySpec{
					{
						Mime: "",
						Replacers: []*k.ReplacerSpec{
							{
								Replacers: &k.ReplacerSpec_Value{
									Value: &k.ValueReplacer{FromTo: map[string]string{"bar": "baz"}},
								},
							},
						},
					},
				},
				mime: "application/json",
				data: `{"foo":"bar", "alice":"bob"}`,
			},
			&action{
				data: `{"foo":"bar", "alice":"bob"}`,
				err:  nil,
			},
		),
		gen(
			"invalid replacer",
			&condition{
				specs: []*v1.LogBodySpec{
					{
						Mime: "application/json",
						Replacers: []*k.ReplacerSpec{
							{
								Replacers: &k.ReplacerSpec_Regexp{
									Regexp: &k.RegexpReplacer{Pattern: `[0-9a-`},
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

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			replacers, err := bodyReplacers(tt.C.specs)
			testutil.Diff(t, tt.A.err, err, cmpopts.EquateErrors())

			bl := &baseLogger{
				bodies: replacers,
			}
			result := bl.logBody(tt.C.mime, []byte(tt.C.data))
			testutil.Diff(t, tt.A.data, string(result))
		})
	}
}
