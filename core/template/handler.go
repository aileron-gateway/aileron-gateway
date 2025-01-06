package template

import (
	"mime"
	"net/http"
	"path"
	"strings"

	"github.com/aileron-gateway/aileron-gateway/core"
	utilhttp "github.com/aileron-gateway/aileron-gateway/util/http"
)

type templateHandler struct {
	// HandlerBase is the base struct for
	// http.Handler type resource.
	// This provides Patterns() and Methods() methods
	// to fulfill the core.Handler interface.
	*utilhttp.HandlerBase

	// eh is the http error handler to
	// serve HTTP errors.
	eh core.ErrorHandler

	// contents is the list of mime contents to be
	// responded to clients.
	// NotAcceptable error will be returned
	// when a content was not found clients can accept.
	contents []*utilhttp.MIMEContent
}

func (h *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := h.findContent(r.Header.Get("Accept"))
	if c == nil {
		h.eh.ServeHTTPError(w, r, utilhttp.ErrNotAcceptable)
		return
	}

	// info is the information that is
	// given to the template.
	info := map[string]any{
		"proto":  r.Proto,
		"host":   r.URL.Host,
		"method": r.Method,
		"path":   r.URL.Path,
		"remote": r.RemoteAddr,
		"header": r.Header,
		"query":  r.URL.Query(),
	}

	header := w.Header()
	for k, v := range c.Header {
		header[k] = append(header[k], v...)
	}
	header.Set("Content-Type", c.MIMEType+"; charset=utf-8")
	header.Set("X-Content-Type-Options", "nosniff")
	header.Add("Vary", "Accept")

	w.WriteHeader(c.StatusCode)
	_, _ = w.Write(c.Content(info))
}

// findContent returns content depending on the Accept header.
// findContent returns nil when an appropriate type of content
// was not found.
func (h *templateHandler) findContent(accept string) *utilhttp.MIMEContent {
	accepts := strings.Split(accept, ",")

	for i := 0; i < len(accepts); i++ {
		mimeType, _, _ := mime.ParseMediaType(accepts[i])
		for j := 0; j < len(h.contents); j++ {
			matched, _ := path.Match(mimeType, h.contents[j].MIMEType)
			if matched {
				return h.contents[j]
			}
		}
	}

	return nil
}
