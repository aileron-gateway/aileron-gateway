// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package txtutil

import (
	"crypto/md5"
	htpl "html/template"
	"testing"
	ttpl "text/template"

	"github.com/aileron-gateway/aileron-gateway/kernel/er"
	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestNewTemplate(t *testing.T) {
	type condition struct {
		typ      TemplateType
		tpl      string
		fallback string
		info     map[string]any
	}

	type action struct {
		result string
		err    error
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	cndText := tb.Condition("text", "use normal text template")
	cndGoText := tb.Condition("go text", "go text template")
	cndGoHTML := tb.Condition("go html", "go html template")
	cndValid := tb.Condition("valid template", "input valid template string")
	actCheckResult := tb.Action("check result", "check that the returned result is the expected one")
	actCheckNoError := tb.Action("no error", "check that there is no error")
	actCheckError := tb.Action("error", "check that an error was returned")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"text",
			[]string{cndText, cndValid},
			[]string{actCheckResult, actCheckNoError},
			&condition{
				typ:      TplText,
				tpl:      "test {{.tag}}",
				fallback: "fallback",
				info:     map[string]any{"tag": "template"},
			},
			&action{
				result: "test {{.tag}}",
			},
		),
		gen(
			"text template",
			[]string{cndGoText, cndValid},
			[]string{actCheckResult, actCheckNoError},
			&condition{
				typ:      TplGoText,
				tpl:      "test {{.tag}}",
				fallback: "fallback",
				info:     map[string]any{"tag": "template"},
			},
			&action{
				result: "test template",
			},
		),
		gen(
			"html template",
			[]string{cndGoHTML, cndValid},
			[]string{actCheckResult, actCheckNoError},
			&condition{
				typ:      TplGoHTML,
				tpl:      "test {{.tag}}",
				fallback: "fallback",
				info:     map[string]any{"tag": "template"},
			},
			&action{
				result: "test template",
			},
		),
		gen(
			"invalid text template",
			[]string{},
			[]string{actCheckError},
			&condition{
				typ:      TplGoText,
				tpl:      "test {{.tag}",
				fallback: "fallback",
				info:     map[string]any{"tag": "template"},
			},
			&action{
				err: &er.Error{
					Package:     ErrPkg,
					Type:        ErrTypeTemplate,
					Description: ErrDscTemplate,
				},
			},
		),
		gen(
			"invalid html template",
			[]string{},
			[]string{actCheckError},
			&condition{
				typ:      TplGoHTML,
				tpl:      "test {{.tag}",
				fallback: "fallback",
				info:     map[string]any{"tag": "template"},
			},
			&action{
				err: &er.Error{
					Package:     ErrPkg,
					Type:        ErrTypeTemplate,
					Description: ErrDscTemplate,
				},
			},
		),
		gen(
			"unsupported type",
			[]string{},
			[]string{actCheckError},
			&condition{
				typ:      TemplateType(999),
				tpl:      "test {{.tag}}",
				fallback: "fallback",
				info:     map[string]any{"tag": "template"},
			},
			&action{
				err: &er.Error{
					Package:     ErrPkg,
					Type:        ErrTypeTemplate,
					Description: ErrDscUnsupported,
				},
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			tpl, err := NewTemplate(tt.C().typ, tt.C().tpl, tt.C().fallback)
			testutil.Diff(t, tt.A().err, err, cmpopts.EquateErrors())
			if err != nil {
				testutil.Diff(t, nil, tpl)
				return
			}

			b := tpl.Content(tt.C().info)
			testutil.Diff(t, tt.A().result, string(b))
		})
	}
}

func TestTemplate_Content(t *testing.T) {
	type condition struct {
		tpl  Template
		info map[string]any
	}

	type action struct {
		result string
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	cndText := tb.Condition("text", "use normal text template")
	cndGoText := tb.Condition("go text", "go text template")
	cndGoHTML := tb.Condition("go html", "go html template")
	cndFallback := tb.Condition("fallback", "make error for fallback")
	cndNilMap := tb.Condition("nil map", "input nil as map")
	actCheckResult := tb.Action("check result", "check that the returned result is the expected one")
	actCheckFallback := tb.Action("fallback", "check that fallback occurred by an error")
	table := tb.Build()

	mustTextTpl := func(tpl string) *ttpl.Template {
		name := md5.Sum([]byte(tpl))
		t, err := ttpl.New(string(name[:])).Parse(tpl)
		if err != nil {
			panic(err)
		}
		return t
	}
	mustHTMLTpl := func(tpl string) *htpl.Template {
		name := md5.Sum([]byte(tpl))
		t, err := htpl.New(string(name[:])).Parse(tpl)
		if err != nil {
			panic(err)
		}
		return t
	}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"text",
			[]string{cndText},
			[]string{actCheckResult},
			&condition{
				tpl: &textTemplate{
					tpl: []byte("test {{.tag}}"),
				},
			},
			&action{
				result: "test {{.tag}}",
			},
		),
		gen(
			"text template",
			[]string{cndGoText},
			[]string{actCheckResult},
			&condition{
				tpl: &goTextTemplate{
					tpl: mustTextTpl("test {{.tag}}"),
				},
				info: map[string]any{"tag": "template"},
			},
			&action{
				result: "test template",
			},
		),
		gen(
			"text template with nil map",
			[]string{cndGoText, cndNilMap},
			[]string{actCheckResult},
			&condition{
				tpl: &goTextTemplate{
					tpl: mustTextTpl("test {{.tag}}"),
				},
				info: nil,
			},
			&action{
				result: "test <no value>", // "<no value>" will be input when value not found.
			},
		),
		gen(
			"html template",
			[]string{cndGoHTML},
			[]string{actCheckResult},
			&condition{
				tpl: &goHTMLTemplate{
					tpl: mustHTMLTpl("test {{.tag}}"),
				},
				info: map[string]any{"tag": "template"},
			},
			&action{
				result: "test template",
			},
		),
		gen(
			"html template with nil map",
			[]string{cndGoHTML, cndNilMap},
			[]string{actCheckResult},
			&condition{
				tpl: &goHTMLTemplate{
					tpl: mustHTMLTpl("test {{.tag}}"),
				},
				info: nil,
			},
			&action{
				result: "test ",
			},
		),
		gen(
			"go text template error",
			[]string{cndGoText, cndFallback},
			[]string{actCheckResult, actCheckFallback},
			&condition{
				tpl: &goTextTemplate{
					tpl:      ttpl.New(""),
					fallback: []byte("fallback"),
				},
			},
			&action{
				result: "fallback",
			},
		),
		gen(
			"go html template error",
			[]string{cndGoHTML, cndFallback},
			[]string{actCheckResult, actCheckFallback},
			&condition{
				tpl: &goHTMLTemplate{
					tpl:      htpl.New(""),
					fallback: []byte("fallback"),
				},
			},
			&action{
				result: "fallback",
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			b := tt.C().tpl.Content(tt.C().info)
			testutil.Diff(t, tt.A().result, string(b))
		})
	}
}
