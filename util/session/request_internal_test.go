// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package session

import (
	"bytes"
	"io"
	"net/http"
	"net/url"
	"testing"

	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

type testErrorReader struct {
	io.Reader
	err error
}

func (r *testErrorReader) Close() error {
	return nil
}

func (r *testErrorReader) Read(p []byte) (n int, err error) {
	if r.err != nil {
		return 0, r.err
	}
	return r.Read(p)
}

func TestPersistRequest(t *testing.T) {
	type condition struct {
		req *http.Request
	}

	type action struct {
		req *httpRequestInfo
		err error
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"empty request",
			[]string{},
			[]string{},
			&condition{
				req: &http.Request{
					Method: http.MethodGet,
					URL:    &url.URL{},
					Header: http.Header{},
					Body:   io.NopCloser(bytes.NewBuffer(nil)),
				},
			},
			&action{
				req: &httpRequestInfo{
					Method: http.MethodGet,
					URL:    &url.URL{},
					Header: http.Header{},
					Body:   []byte{},
				},
				err: nil,
			},
		),
		gen(
			"persist url",
			[]string{},
			[]string{},
			&condition{
				req: &http.Request{
					Method: http.MethodGet,
					URL: &url.URL{
						Scheme:      "Scheme",
						Opaque:      "Opaque",
						User:        url.UserPassword("foo", "bar"),
						Host:        "Host",
						Path:        "Path",
						RawPath:     "RawPath",
						OmitHost:    true,
						ForceQuery:  true,
						RawQuery:    "RawQuery",
						Fragment:    "Fragment",
						RawFragment: "RawFragment",
					},
					Header: http.Header{},
					Body:   io.NopCloser(bytes.NewBuffer(nil)),
				},
			},
			&action{
				req: &httpRequestInfo{
					Method: http.MethodGet,
					URL: &url.URL{
						Scheme:      "Scheme",
						Opaque:      "Opaque",
						User:        nil, // User is not persisted.
						Host:        "Host",
						Path:        "Path",
						RawPath:     "RawPath",
						OmitHost:    true,
						ForceQuery:  true,
						RawQuery:    "RawQuery",
						Fragment:    "Fragment",
						RawFragment: "RawFragment",
					},
					Header: http.Header{},
					Body:   []byte{},
				},
				err: nil,
			},
		),
		gen(
			"persist header",
			[]string{},
			[]string{},
			&condition{
				req: &http.Request{
					Method: http.MethodGet,
					URL:    &url.URL{},
					Header: http.Header{
						"Foo":           []string{"bar", "baz"},
						"Cookie":        []string{"alice=1", "bob=2"},
						"Authorization": []string{"Basic foo:bar"},
					},
					Body: io.NopCloser(bytes.NewBuffer(nil)),
				},
			},
			&action{
				req: &httpRequestInfo{
					Method: http.MethodGet,
					URL:    &url.URL{},
					Header: http.Header{
						"Foo": []string{"bar", "baz"},
					},
					Body: []byte{},
				},
				err: nil,
			},
		),
		gen(
			"persist body",
			[]string{},
			[]string{},
			&condition{
				req: &http.Request{
					Method: http.MethodGet,
					URL:    &url.URL{},
					Header: http.Header{},
					Body:   io.NopCloser(bytes.NewBuffer([]byte("foobar"))),
				},
			},
			&action{
				req: &httpRequestInfo{
					Method: http.MethodGet,
					URL:    &url.URL{},
					Header: http.Header{},
					Body:   []byte("foobar"),
				},
				err: nil,
			},
		),
		gen(
			"body read error",
			[]string{},
			[]string{},
			&condition{
				req: &http.Request{
					Method: http.MethodGet,
					URL:    &url.URL{},
					Header: http.Header{},
					Body:   &testErrorReader{err: io.ErrUnexpectedEOF},
				},
			},
			&action{
				req: &httpRequestInfo{},
				err: io.ErrUnexpectedEOF,
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			ss := NewDefaultSession(SerializeJSON)
			err := PersistRequest(ss, tt.C().req)
			testutil.Diff(t, tt.A().err, err, cmpopts.EquateErrors())

			req := &httpRequestInfo{}
			ss.Extract(requestSessionKey, req)
			testutil.Diff(t, tt.A().req, req)
		})
	}
}

func TestExtractRequest(t *testing.T) {
	type condition struct {
		info *httpRequestInfo
		req  *http.Request
	}

	type action struct {
		req  *http.Request
		body string
		err  error
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"no data",
			[]string{},
			[]string{},
			&condition{
				info: nil,
				req:  &http.Request{},
			},
			&action{
				req: nil,
				err: NoValue,
			},
		),
		gen(
			"empty info",
			[]string{},
			[]string{},
			&condition{
				info: &httpRequestInfo{
					URL:    &url.URL{},
					Header: http.Header{},
				},
				req: &http.Request{
					Method: http.MethodGet,
					URL:    &url.URL{Scheme: "Schema", User: url.UserPassword("foo", "bar")},
					Header: http.Header{
						"Foo":           []string{"bar", "baz"},
						"Cookie":        []string{"alice=1", "bob=2"},
						"Authorization": []string{"Basic foo:bar"},
					},
					Body: io.NopCloser(bytes.NewReader([]byte("body"))),
				},
			},
			&action{
				req: &http.Request{
					URL: &url.URL{
						User: url.UserPassword("foo", "bar"),
					},
					Header: http.Header{
						"Cookie":        []string{"alice=1", "bob=2"},
						"Authorization": []string{"Basic foo:bar"},
					},
				},
				err: nil,
			},
		),
		gen(
			"non empty request",
			[]string{},
			[]string{},
			&condition{
				info: &httpRequestInfo{
					Method: http.MethodPost,
					URL: &url.URL{
						Scheme:      "Scheme",
						Opaque:      "Opaque",
						User:        url.UserPassword("alice", "bob"),
						Host:        "Host",
						Path:        "Path",
						RawPath:     "RawPath",
						OmitHost:    true,
						ForceQuery:  true,
						RawQuery:    "RawQuery",
						Fragment:    "Fragment",
						RawFragment: "RawFragment",
					},
					Header: http.Header{
						"Test":          []string{"john", "doe"},
						"Cookie":        []string{"key=value"},
						"Authorization": []string{"Basic credentials"},
					},
				},
				req: &http.Request{
					Method: http.MethodGet,
					URL:    &url.URL{Scheme: "Schema", User: url.UserPassword("foo", "bar")},
					Header: http.Header{
						"Foo":           []string{"bar", "baz"},
						"Cookie":        []string{"alice=1", "bob=2"},
						"Authorization": []string{"Basic foo:bar"},
					},
					Body: io.NopCloser(bytes.NewReader([]byte("body"))),
				},
			},
			&action{
				req: &http.Request{
					Method: http.MethodPost,
					URL: &url.URL{
						Scheme:      "Scheme",
						Opaque:      "Opaque",
						User:        url.UserPassword("foo", "bar"),
						Host:        "Host",
						Path:        "Path",
						RawPath:     "RawPath",
						OmitHost:    true,
						ForceQuery:  true,
						RawQuery:    "RawQuery",
						Fragment:    "Fragment",
						RawFragment: "RawFragment",
					},
					Header: http.Header{
						"Test":          []string{"john", "doe"},
						"Cookie":        []string{"alice=1", "bob=2"},
						"Authorization": []string{"Basic foo:bar"},
					},
				},
				err: nil,
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			ss := NewDefaultSession(SerializeJSON)
			if tt.C().info != nil {
				err := ss.Persist(requestSessionKey, tt.C().info)
				testutil.Diff(t, nil, err)
			}

			req, err := ExtractRequest(ss, tt.C().req)
			testutil.Diff(t, tt.A().err, err, cmpopts.EquateErrors())
			if err != nil {
				return
			}
			_, err = ExtractRequest(ss, &http.Request{})
			testutil.Diff(t, NoValue, err, cmpopts.EquateErrors())

			body, _ := io.ReadAll(req.Body)
			testutil.Diff(t, tt.A().req.Method, req.Method)
			testutil.Diff(t, tt.A().req.URL, req.URL, cmp.AllowUnexported(url.Userinfo{}))
			testutil.Diff(t, tt.A().req.Header, req.Header)
			testutil.Diff(t, tt.A().body, string(body))
		})
	}
}
