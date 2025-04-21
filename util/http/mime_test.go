// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package http_test

import (
	stdcmp "cmp"
	"net/http"
	"os"
	"testing"

	v1 "github.com/aileron-gateway/aileron-gateway/apis/core/v1"
	k "github.com/aileron-gateway/aileron-gateway/apis/kernel"
	"github.com/aileron-gateway/aileron-gateway/kernel/er"
	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
	"github.com/aileron-gateway/aileron-gateway/kernel/txtutil"
	utilhttp "github.com/aileron-gateway/aileron-gateway/util/http"
	"github.com/google/go-cmp/cmp/cmpopts"
)

// testDir is the path to the test data.
// This path can be changed by the environmental variable.
var testDir = stdcmp.Or(os.Getenv("TEST_DIR"), "../../test/")

func TestNewTemplate(t *testing.T) {
	type condition struct {
		spec *v1.MIMEContentSpec
		info map[string]any
	}

	type action struct {
		status int
		header http.Header
		mime   string
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
				spec: &v1.MIMEContentSpec{
					MIMEType:     "text/plain",
					StatusCode:   http.StatusOK,
					TemplateType: k.TemplateType_Text,
					Template:     "test {{.tag}}",
					FallbackText: "fallback",
				},
				info: map[string]any{"tag": "template"},
			},
			&action{
				status: http.StatusOK,
				header: http.Header{},
				mime:   "text/plain",
				result: "test {{.tag}}",
			},
		),
		gen(
			"text with status 0",
			[]string{cndText, cndValid},
			[]string{actCheckResult, actCheckNoError},
			&condition{
				spec: &v1.MIMEContentSpec{
					MIMEType:     "text/plain",
					StatusCode:   0,
					TemplateType: k.TemplateType_Text,
					Template:     "test {{.tag}}",
					FallbackText: "fallback",
				},
				info: map[string]any{"tag": "template"},
			},
			&action{
				status: 0,
				header: http.Header{},
				mime:   "text/plain",
				result: "test {{.tag}}",
			},
		),
		gen(
			"text with header",
			[]string{cndText, cndValid},
			[]string{actCheckResult, actCheckNoError},
			&condition{
				spec: &v1.MIMEContentSpec{
					MIMEType:     "text/plain",
					StatusCode:   http.StatusOK,
					Header:       map[string]string{"foo": "bar", "alice": "bob"},
					TemplateType: k.TemplateType_Text,
					Template:     "test {{.tag}}",
					FallbackText: "fallback",
				},
				info: map[string]any{"tag": "template"},
			},
			&action{
				status: http.StatusOK,
				header: http.Header{"Foo": []string{"bar"}, "Alice": []string{"bob"}},
				mime:   "text/plain",
				result: "test {{.tag}}",
			},
		),
		gen(
			"text template",
			[]string{cndGoText, cndValid},
			[]string{actCheckResult, actCheckNoError},
			&condition{
				spec: &v1.MIMEContentSpec{
					MIMEType:     "text/plain",
					StatusCode:   http.StatusOK,
					TemplateType: k.TemplateType_GoText,
					Template:     "test {{.tag}}",
					FallbackText: "fallback",
				},
				info: map[string]any{"tag": "template"},
			},
			&action{
				status: http.StatusOK,
				header: http.Header{},
				mime:   "text/plain",
				result: "test template",
			},
		),
		gen(
			"html template",
			[]string{cndGoHTML, cndValid},
			[]string{actCheckResult, actCheckNoError},
			&condition{
				spec: &v1.MIMEContentSpec{
					MIMEType:     "text/plain",
					StatusCode:   http.StatusOK,
					TemplateType: k.TemplateType_GoHTML,
					Template:     "test {{.tag}}",
					FallbackText: "fallback",
				},
				info: map[string]any{"tag": "template"},
			},
			&action{
				status: http.StatusOK,
				header: http.Header{},
				mime:   "text/plain",
				result: "test template",
			},
		),
		gen(
			"invalid text template",
			[]string{cndGoText},
			[]string{actCheckError},
			&condition{
				spec: &v1.MIMEContentSpec{
					MIMEType:     "text/plain",
					StatusCode:   http.StatusOK,
					TemplateType: k.TemplateType_GoText,
					Template:     "test {{.tag}",
					FallbackText: "fallback",
				},
				info: map[string]any{"tag": "template"},
			},
			&action{
				err: &er.Error{
					Package:     txtutil.ErrPkg,
					Type:        txtutil.ErrTypeTemplate,
					Description: txtutil.ErrDscTemplate,
				},
			},
		),
		gen(
			"invalid html template",
			[]string{cndGoHTML},
			[]string{actCheckError},
			&condition{
				spec: &v1.MIMEContentSpec{
					MIMEType:     "text/plain",
					StatusCode:   http.StatusOK,
					TemplateType: k.TemplateType_GoHTML,
					Template:     "test {{.tag}",
					FallbackText: "fallback",
				},
				info: map[string]any{"tag": "template"},
			},
			&action{
				err: &er.Error{
					Package:     txtutil.ErrPkg,
					Type:        txtutil.ErrTypeTemplate,
					Description: txtutil.ErrDscTemplate,
				},
			},
		),
		gen(
			"empty mime type",
			[]string{cndGoText},
			[]string{actCheckError},
			&condition{
				spec: &v1.MIMEContentSpec{
					MIMEType:     "",
					StatusCode:   http.StatusOK,
					TemplateType: k.TemplateType_GoText,
					Template:     "test",
				},
				info: map[string]any{"tag": "template"},
			},
			&action{
				err: &er.Error{
					Package:     utilhttp.ErrPkg,
					Type:        utilhttp.ErrTypeMime,
					Description: utilhttp.ErrDscParseMime,
				},
			},
		),
		gen(
			"invalid mime type",
			[]string{cndGoText},
			[]string{actCheckError},
			&condition{
				spec: &v1.MIMEContentSpec{
					MIMEType:     "invalid/text/plain",
					StatusCode:   http.StatusOK,
					TemplateType: k.TemplateType_GoText,
					Template:     "test {{.tag}}",
					FallbackText: "fallback",
				},
				info: map[string]any{"tag": "template"},
			},
			&action{
				err: &er.Error{
					Package:     utilhttp.ErrPkg,
					Type:        utilhttp.ErrTypeMime,
					Description: utilhttp.ErrDscParseMime,
				},
			},
		),
		gen(
			"text from file",
			[]string{cndText, cndValid},
			[]string{actCheckResult, actCheckNoError},
			&condition{
				spec: &v1.MIMEContentSpec{
					MIMEType:     "text/plain",
					StatusCode:   http.StatusOK,
					TemplateType: k.TemplateType_Text,
					TemplateFile: testDir + "ut/core/utilhttp/template.txt",
					FallbackText: "fallback",
				},
				info: map[string]any{"tag": "template"},
			},
			&action{
				status: http.StatusOK,
				header: http.Header{},
				mime:   "text/plain",
				result: "test {{.tag}}",
			},
		),
		gen(
			"file read error",
			[]string{cndText, cndValid},
			[]string{actCheckResult, actCheckError},
			&condition{
				spec: &v1.MIMEContentSpec{
					MIMEType:     "text/plain",
					StatusCode:   http.StatusOK,
					TemplateType: k.TemplateType_Text,
					TemplateFile: testDir + "ut/core/utilhttp/not-exist.txt",
					FallbackText: "fallback",
				},
				info: map[string]any{"tag": "template"},
			},
			&action{
				err: &er.Error{
					Package:     utilhttp.ErrPkg,
					Type:        utilhttp.ErrTypeMime,
					Description: utilhttp.ErrDscIO,
				},
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			c, err := utilhttp.NewMIMEContent(tt.C().spec)
			testutil.Diff(t, tt.A().err, err, cmpopts.EquateErrors())
			if err != nil {
				testutil.Diff(t, (*utilhttp.MIMEContent)(nil), c)
				return
			}

			b := c.Content(tt.C().info)
			testutil.Diff(t, tt.A().mime, c.MIMEType)
			testutil.Diff(t, tt.A().status, c.StatusCode)
			testutil.Diff(t, tt.A().result, string(b))
			testutil.Diff(t, int(tt.C().spec.StatusCode), c.StatusCode)
			testutil.Diff(t, tt.A().header, c.Header, cmpopts.SortMaps(func(x, y string) bool { return x > y }))
		})
	}
}
