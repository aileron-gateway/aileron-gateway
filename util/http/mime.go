package http

import (
	"mime"
	"net/http"
	"os"

	v1 "github.com/aileron-gateway/aileron-gateway/apis/core/v1"
	"github.com/aileron-gateway/aileron-gateway/apis/kernel"
	"github.com/aileron-gateway/aileron-gateway/kernel/er"
	"github.com/aileron-gateway/aileron-gateway/kernel/txtutil"
)

// MIMEContent provides HTTP content corresponding to a MIMEType.
// MIMEContent utilize txtutil.Template in it to provides
// a HTTP body content.
type MIMEContent struct {
	txtutil.Template
	MIMEType   string
	StatusCode int
	Header     http.Header
}

func NewMIMEContent(spec *v1.MIMEContentSpec) (*MIMEContent, error) {
	mt, _, err := mime.ParseMediaType(spec.MIMEType)
	if err != nil {
		return nil, (&er.Error{
			Package:     ErrPkg,
			Type:        ErrTypeMime,
			Description: ErrDscParseMime,
		}).Wrap(err)
	}

	body := spec.Template
	if spec.TemplateFile != "" {
		b, err := os.ReadFile(spec.TemplateFile)
		if err != nil {
			return nil, (&er.Error{
				Package:     ErrPkg,
				Type:        ErrTypeMime,
				Description: ErrDscIO,
			}).Wrap(err)
		}
		body = string(b)
	}

	typ := map[kernel.TemplateType]txtutil.TemplateType{
		kernel.TemplateType_Text:   txtutil.TplText,
		kernel.TemplateType_GoText: txtutil.TplGoText,
		kernel.TemplateType_GoHTML: txtutil.TplGoHTML,
	}[spec.TemplateType]
	tpl, err := txtutil.NewTemplate(typ, body, spec.FallbackText)
	if err != nil {
		return nil, err // Return err as-is.
	}

	header := http.Header{}
	for k, v := range spec.Header {
		header.Set(k, v)
	}

	return &MIMEContent{
		Template:   tpl,
		MIMEType:   mt,
		StatusCode: int(spec.StatusCode),
		Header:     header,
	}, nil
}
