package template

import (
	"io"
	"maps"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
	"github.com/aileron-gateway/aileron-gateway/kernel/txtutil"
	utilhttp "github.com/aileron-gateway/aileron-gateway/util/http"
)

func TestHandler_ServeHTTP(t *testing.T) {
	type condition struct {
		h      *templateHandler
		header http.Header
		query  string
	}

	type action struct {
		body   string
		status int
		header http.Header
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	tpl0, _ := txtutil.NewTemplate(txtutil.TplGoText, "", "")
	tpl1, _ := txtutil.NewTemplate(txtutil.TplGoText, "{{.proto}} {{.host}} {{.method}} {{.path}}", "")
	tpl2, _ := txtutil.NewTemplate(txtutil.TplGoText, "{{.header.Foo}}", "")
	tpl3, _ := txtutil.NewTemplate(txtutil.TplGoText, "{{.query.foo}}", "")

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"no content",
			[]string{},
			[]string{},
			&condition{
				h: &templateHandler{
					eh: utilhttp.GlobalErrorHandler(utilhttp.DefaultErrorHandlerName),
				},
			},
			&action{
				status: http.StatusNotAcceptable,
				body:   `{"status":406,"statusText":"Not Acceptable"}`,
			},
		),
		gen(
			"add response header",
			[]string{},
			[]string{},
			&condition{
				h: &templateHandler{
					eh: utilhttp.GlobalErrorHandler(utilhttp.DefaultErrorHandlerName),
					contents: []*utilhttp.MIMEContent{
						{
							Template:   tpl0,
							StatusCode: http.StatusOK,
							Header:     http.Header{"Foo": []string{"bar", "baz"}},
							MIMEType:   "text/plain",
						},
					},
				},
				header: http.Header{
					"Accept":       []string{"text/plain"},
					"Content-Type": []string{"text/plain"},
				},
			},
			&action{
				status: http.StatusOK,
				header: http.Header{"Foo": []string{"bar", "baz"}},
				body:   "",
			},
		),
		gen(
			"basic info",
			[]string{},
			[]string{},
			&condition{
				h: &templateHandler{
					eh: utilhttp.GlobalErrorHandler(utilhttp.DefaultErrorHandlerName),
					contents: []*utilhttp.MIMEContent{
						{
							Template:   tpl1,
							StatusCode: http.StatusOK,
							MIMEType:   "text/plain",
						},
					},
				},
				header: http.Header{
					"Accept":       []string{"text/plain"},
					"Content-Type": []string{"text/plain"},
				},
			},
			&action{
				status: http.StatusOK,
				body:   "HTTP/1.1 test.com POST /unit-test",
			},
		),
		gen(
			"template header",
			[]string{},
			[]string{},
			&condition{
				h: &templateHandler{
					eh: utilhttp.GlobalErrorHandler(utilhttp.DefaultErrorHandlerName),
					contents: []*utilhttp.MIMEContent{
						{
							Template:   tpl2,
							StatusCode: http.StatusOK,
							MIMEType:   "application/json",
						},
					},
				},
				header: http.Header{
					"Accept":       []string{"application/json"},
					"Content-Type": []string{"application/json"},
					"Foo":          []string{"bar"},
				},
			},
			&action{
				status: http.StatusOK,
				body:   "[bar]",
			},
		),
		gen(
			"template header",
			[]string{},
			[]string{},
			&condition{
				h: &templateHandler{
					eh: utilhttp.GlobalErrorHandler(utilhttp.DefaultErrorHandlerName),
					contents: []*utilhttp.MIMEContent{
						{
							Template:   tpl3,
							StatusCode: http.StatusOK,
							MIMEType:   "application/json",
						},
					},
				},
				header: http.Header{
					"Accept":       []string{"application/json"},
					"Content-Type": []string{"application/json"},
				},
				query: "foo=bar",
			},
			&action{
				status: http.StatusOK,
				body:   "[bar]",
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			w := httptest.NewRecorder()
			r, _ := http.NewRequest(http.MethodPost, "http://test.com/unit-test?"+tt.C().query, nil)
			maps.Copy(r.Header, tt.C().header)

			tt.C().h.ServeHTTP(w, r)

			res := w.Result()
			b, _ := io.ReadAll(res.Body)
			testutil.Diff(t, tt.A().status, res.StatusCode)
			for k, v := range tt.A().header {
				testutil.Diff(t, v, res.Header[k])
			}
			testutil.Diff(t, tt.A().body, string(b))
		})
	}
}
