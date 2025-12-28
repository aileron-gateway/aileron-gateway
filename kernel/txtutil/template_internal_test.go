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
		typ  TemplateType
		tpl  string
		info map[string]any
	}

	type action struct {
		result string
		err    error
	}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"text", &condition{
				typ:  TplText,
				tpl:  "test {{.tag}}",
				info: map[string]any{"tag": "template"},
			},
			&action{
				result: "test {{.tag}}",
			},
		),
		gen(
			"text template", &condition{
				typ:  TplGoText,
				tpl:  "test {{.tag}}",
				info: map[string]any{"tag": "template"},
			},
			&action{
				result: "test template",
			},
		),
		gen(
			"html template", &condition{
				typ:  TplGoHTML,
				tpl:  "test {{.tag}}",
				info: map[string]any{"tag": "template"},
			},
			&action{
				result: "test template",
			},
		),
		gen(
			"invalid text template", &condition{
				typ:  TplGoText,
				tpl:  "test {{.tag}",
				info: map[string]any{"tag": "template"},
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
			"invalid html template", &condition{
				typ:  TplGoHTML,
				tpl:  "test {{.tag}",
				info: map[string]any{"tag": "template"},
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
			"unsupported type", &condition{
				typ:  TemplateType(999),
				tpl:  "test {{.tag}}",
				info: map[string]any{"tag": "template"},
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

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			tpl, err := NewTemplate(tt.C.typ, tt.C.tpl)
			testutil.Diff(t, tt.A.err, err, cmpopts.EquateErrors())
			if err != nil {
				testutil.Diff(t, nil, tpl)
				return
			}

			b := tpl.Content(tt.C.info)
			testutil.Diff(t, tt.A.result, string(b))
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
			"text", &condition{
				tpl: &textTemplate{
					tpl: []byte("test {{.tag}}"),
				},
			},
			&action{
				result: "test {{.tag}}",
			},
		),
		gen(
			"text template", &condition{
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
			"text template with nil map", &condition{
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
			"html template", &condition{
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
			"html template with nil map", &condition{
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
			"go text template error", &condition{
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
			"go html template error", &condition{
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

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			b := tt.C.tpl.Content(tt.C.info)
			testutil.Diff(t, tt.A.result, string(b))
		})
	}
}
