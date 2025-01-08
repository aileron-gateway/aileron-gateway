package session

import (
	"bytes"
	"io"
	"net/http"
	"net/url"
)

const (
	// requestSessionKey is the key to save HTTP request in the session.
	// This name must not be used elsewhere.
	requestSessionKey = "__sec_request__"
)

// httpRequestInfo is the request data to be saved in the session.
type httpRequestInfo struct {
	Method string      `json:"m" msgpack:"m"`
	URL    *url.URL    `json:"u" msgpack:"u"`
	Header http.Header `json:"h" msgpack:"h"`
	Body   []byte      `json:"b" msgpack:"b"`
}

// PersistRequest persists request header, URL, method and body in the session.
// Cookie in the header won't be saved for security reason and for use of cookie storage.
// This method will consume request body and won't rewind it.
// ss must not be nil.
func PersistRequest(ss Session, r *http.Request) error {
	var body []byte
	if r.Body != nil && r.Body != http.NoBody {
		var err error
		body, err = io.ReadAll(r.Body)
		if err != nil {
			return err
		}
	}

	h := r.Header.Clone()
	delete(h, "Cookie")
	delete(h, "Authorization")

	u := &url.URL{
		Scheme:      r.URL.Scheme,
		Opaque:      r.URL.Opaque,
		User:        nil,
		Host:        r.URL.Host,
		Path:        r.URL.Path,
		RawPath:     r.URL.RawPath,
		OmitHost:    r.URL.OmitHost,
		ForceQuery:  r.URL.ForceQuery,
		RawQuery:    r.URL.RawQuery,
		Fragment:    r.URL.Fragment,
		RawFragment: r.URL.RawFragment,
	}

	req := &httpRequestInfo{
		Header: h,
		URL:    u,
		Method: r.Method,
		Body:   body,
	}
	return ss.Persist(requestSessionKey, req)
}

// ExtractRequest extracts request header, URL, method and body
// from the session and return the request with  extracted request values.
// ss must not be nil.
func ExtractRequest(ss Session, r *http.Request) (*http.Request, error) {
	req := &httpRequestInfo{
		URL:    &url.URL{},
		Header: http.Header{},
	}

	if err := ss.Extract(requestSessionKey, req); err != nil {
		return nil, err
	}

	newReq, _ := http.NewRequestWithContext(r.Context(), req.Method, "", nil)
	newReq.Method = req.Method
	newReq.URL = req.URL
	newReq.Header = req.Header
	newReq.Body = io.NopCloser(bytes.NewReader(req.Body))

	// Keep URL.User, Cookie and Authorization headers
	// because they contains authentication or sensitive data.
	newReq.URL.User = r.URL.User
	if cks := r.Header["Cookie"]; len(cks) > 0 {
		newReq.Header["Cookie"] = cks
	}
	if ahs := r.Header["Authorization"]; len(ahs) > 0 {
		newReq.Header["Authorization"] = ahs
	}

	ss.Delete(requestSessionKey)
	return newReq, nil
}
