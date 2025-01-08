package echo

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"

	utilhttp "github.com/aileron-gateway/aileron-gateway/util/http"
)

type echo struct {
	*utilhttp.HandlerBase
}

func (h *echo) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var body []byte
	if r.ContentLength > 1<<24 {
		body = []byte("REQUEST BODY TOO LARGE\n")
		_, _ = io.Copy(io.Discard, r.Body)
	} else {
		body, _ = io.ReadAll(r.Body)
	}

	unescapedPath, err := url.PathUnescape(r.URL.Path)
	if err != nil {
		w.Header().Set("Content-Type", "text/plain;charset=utf-8")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	unescapedQuery, err := url.QueryUnescape(r.URL.RawQuery)
	if err != nil {
		w.Header().Set("Content-Type", "text/plain;charset=utf-8")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	var b bytes.Buffer
	enc := json.NewEncoder(&b)
	enc.SetIndent("", "  ")
	_ = enc.Encode(r.Header)
	header := b.Bytes()

	w.Header().Set("Content-Type", "text/plain;charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(http.StatusOK)

	var buf bytes.Buffer
	buf.Write([]byte("---------- Request ----------\n"))
	buf.Write([]byte("Proto   : " + r.Proto + "\n"))
	buf.Write([]byte("Host   : " + r.Host + "\n"))
	buf.Write([]byte("Method : " + r.Method + "\n"))
	buf.Write([]byte("URI    : " + r.RequestURI + "\n"))
	buf.Write([]byte("Path   : " + unescapedPath + "\n"))
	buf.Write([]byte("Query  : " + unescapedQuery + "\n"))
	buf.Write([]byte("Remote : " + r.RemoteAddr + "\n"))
	buf.Write([]byte("---------- Header ----------\n"))
	buf.Write(header)
	buf.Write([]byte("---------- Body ----------\n"))
	buf.Write(body)
	buf.Write([]byte("\n--------------------------\n"))

	// Write to the response
	_, _ = w.Write(buf.Bytes())
}
