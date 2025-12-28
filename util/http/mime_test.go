// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package http_test

import (
	"net/http"
	"testing"

	v1 "github.com/aileron-gateway/aileron-gateway/apis/core/v1"
	"github.com/aileron-gateway/aileron-gateway/internal/testutil"
	"github.com/aileron-gateway/aileron-gateway/kernel/er"
	"github.com/aileron-gateway/aileron-gateway/kernel/txtutil"
	utilhttp "github.com/aileron-gateway/aileron-gateway/util/http"
	"github.com/google/go-cmp/cmp/cmpopts"
)

// testDir is the path to the test data.
var testDir = "../../test/"

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

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"text", &condition{
				spec: &v1.MIMEContentSpec{
					MIMEType:     "text/plain",
					StatusCode:   http.StatusOK,
					TemplateType: v1.TemplateType_Text,
					Template:     "test {{.tag}}",
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
			"text with status 0", &condition{
				spec: &v1.MIMEContentSpec{
					MIMEType:     "text/plain",
					StatusCode:   0,
					TemplateType: v1.TemplateType_Text,
					Template:     "test {{.tag}}",
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
			"text with header", &condition{
				spec: &v1.MIMEContentSpec{
					MIMEType:     "text/plain",
					StatusCode:   http.StatusOK,
					Header:       map[string]string{"foo": "bar", "alice": "bob"},
					TemplateType: v1.TemplateType_Text,
					Template:     "test {{.tag}}",
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
			"text template", &condition{
				spec: &v1.MIMEContentSpec{
					MIMEType:     "text/plain",
					StatusCode:   http.StatusOK,
					TemplateType: v1.TemplateType_GoText,
					Template:     "test {{.tag}}",
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
			"html template", &condition{
				spec: &v1.MIMEContentSpec{
					MIMEType:     "text/plain",
					StatusCode:   http.StatusOK,
					TemplateType: v1.TemplateType_GoHTML,
					Template:     "test {{.tag}}",
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
			"invalid text template", &condition{
				spec: &v1.MIMEContentSpec{
					MIMEType:     "text/plain",
					StatusCode:   http.StatusOK,
					TemplateType: v1.TemplateType_GoText,
					Template:     "test {{.tag}",
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
			"invalid html template", &condition{
				spec: &v1.MIMEContentSpec{
					MIMEType:     "text/plain",
					StatusCode:   http.StatusOK,
					TemplateType: v1.TemplateType_GoHTML,
					Template:     "test {{.tag}",
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
			"empty mime type", &condition{
				spec: &v1.MIMEContentSpec{
					MIMEType:     "",
					StatusCode:   http.StatusOK,
					TemplateType: v1.TemplateType_GoText,
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
			"invalid mime type", &condition{
				spec: &v1.MIMEContentSpec{
					MIMEType:     "invalid/text/plain",
					StatusCode:   http.StatusOK,
					TemplateType: v1.TemplateType_GoText,
					Template:     "test {{.tag}}",
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
			"text from file", &condition{
				spec: &v1.MIMEContentSpec{
					MIMEType:     "text/plain",
					StatusCode:   http.StatusOK,
					TemplateType: v1.TemplateType_Text,
					TemplateFile: testDir + "ut/core/utilhttp/template.txt",
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
			"file read error", &condition{
				spec: &v1.MIMEContentSpec{
					MIMEType:     "text/plain",
					StatusCode:   http.StatusOK,
					TemplateType: v1.TemplateType_Text,
					TemplateFile: testDir + "ut/core/utilhttp/not-exist.txt",
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

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			c, err := utilhttp.NewMIMEContent(tt.C.spec)
			testutil.Diff(t, tt.A.err, err, cmpopts.EquateErrors())
			if err != nil {
				testutil.Diff(t, (*utilhttp.MIMEContent)(nil), c)
				return
			}

			b := c.Content(tt.C.info)
			testutil.Diff(t, tt.A.mime, c.MIMEType)
			testutil.Diff(t, tt.A.status, c.StatusCode)
			testutil.Diff(t, tt.A.result, string(b))
			testutil.Diff(t, int(tt.C.spec.StatusCode), c.StatusCode)
			testutil.Diff(t, tt.A.header, c.Header, cmpopts.SortMaps(func(x, y string) bool { return x > y }))
		})
	}
}
