// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package txtutil

import (
	"bytes"
	htpl "html/template"
	"strconv"
	ttpl "text/template"

	"github.com/aileron-gateway/aileron-gateway/kernel/er"
)

// TemplateType is type of document template.
// Following three types of templates are supported.
//
//   - TplText : plain text template. Static plain text document. No information embed-able.
//   - TplGoText : "html/template" template. HTML template document. External information embed-able.
//   - TplGoHTML : "text/template" template. Text template document. External information embed-able.
type TemplateType int

const (
	TplText   TemplateType = iota // TplText represents the template type of plain text.
	TplGoText                     // TplGoText represents the template of "text/template".
	TplGoHTML                     // TplGoHTML represents the template of "html/template".
)

// Template returns content from the pre-set template.
type Template interface {
	Content(map[string]any) []byte
}

// NewTemplate returns a new template.
// tplFallback is the static content that will be returned for fallback
// when embedding external information to the template failed.
// tplFallback will be ignored when the template type is TplText.
func NewTemplate(typ TemplateType, tpl string) (Template, error) {
	switch typ {
	case TplText:
		return &textTemplate{
			tpl: []byte(tpl),
		}, nil

	case TplGoText:
		template, err := ttpl.New("").Parse(tpl)
		if err != nil {
			return nil, (&er.Error{Package: ErrPkg, Type: ErrTypeTemplate, Description: ErrDscTemplate, Detail: "GoText `" + tpl + "`"}).Wrap(err)
		}
		return &goTextTemplate{
			tpl:      template,
			fallback: nil,
		}, nil

	case TplGoHTML:
		template, err := htpl.New("").Parse(tpl)
		if err != nil {
			return nil, (&er.Error{Package: ErrPkg, Type: ErrTypeTemplate, Description: ErrDscTemplate, Detail: "GoHTML `" + tpl + "`"}).Wrap(err)
		}
		return &goHTMLTemplate{
			tpl:      template,
			fallback: nil,
		}, nil

	default:
		return nil, (&er.Error{Package: ErrPkg, Type: ErrTypeTemplate, Description: ErrDscUnsupported, Detail: strconv.Itoa(int(typ))})
	}
}

// textTemplate is the plain text template.
// textTemplate implements Template interface.
type textTemplate struct {
	tpl []byte
}

func (c *textTemplate) Content(_ map[string]any) []byte {
	return c.tpl
}

// goTextTemplate is the template using "text/template".
// goTextTemplate implements Template interface.
type goTextTemplate struct {
	tpl      *ttpl.Template
	fallback []byte
}

func (c *goTextTemplate) Content(info map[string]any) []byte {
	var buf bytes.Buffer
	if err := c.tpl.Execute(&buf, info); err != nil {
		return c.fallback // Fallback.
	}
	return buf.Bytes()
}

// goHTMLTemplate is the template using "html/template".
// goHTMLTemplate implements Template interface.
type goHTMLTemplate struct {
	tpl      *htpl.Template
	fallback []byte
}

func (c *goHTMLTemplate) Content(info map[string]any) []byte {
	var buf bytes.Buffer
	if err := c.tpl.Execute(&buf, info); err != nil {
		return c.fallback // Fallback.
	}
	return buf.Bytes()
}
