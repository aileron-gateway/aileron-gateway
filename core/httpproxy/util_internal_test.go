// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package httpproxy

import (
	"crypto/tls"
	"net/http"
	"net/url"
	"testing"

	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
)

func TestUpgradeType(t *testing.T) {
	type condition struct {
		h http.Header
	}

	type action struct {
		typ string
	}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"header is empty",
			&condition{
				h: http.Header{},
			},
			&action{
				typ: "",
			},
		),
		gen(
			"value is not Upgrade",
			&condition{
				h: http.Header{
					"Connection": []string{"value"},
				},
			},
			&action{
				typ: "",
			},
		),
		gen(
			"No Upgrade header",
			&condition{
				h: http.Header{
					"Connection": []string{"Upgrade"},
				},
			},
			&action{
				typ: "",
			},
		),
		gen(
			"only upgrade in connection header",
			&condition{
				h: http.Header{
					"Connection": []string{"upgrade"},
					"Upgrade":    []string{"test"},
				},
			},
			&action{
				typ: "test",
			},
		),
		gen(
			"upgrade in connection header",
			&condition{
				h: http.Header{
					"Connection": []string{"foo,upgrade,bar"},
					"Upgrade":    []string{"test"},
				},
			},
			&action{
				typ: "test",
			},
		),
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			typ := upgradeType(tt.C.h)
			testutil.Diff(t, tt.A.typ, typ)
		})
	}
}

func TestCopyHeader(t *testing.T) {
	type condition struct {
		src http.Header
		dst http.Header
	}

	type action struct {
		dst http.Header
	}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"copy header values",
			&condition{
				src: http.Header{
					"test":  []string{"src"},
					"test1": []string{"src"},
				},
				dst: http.Header{
					"test":  []string{"dst"},
					"test2": []string{"dst"},
				},
			},
			&action{
				dst: http.Header{
					"Test":  []string{"src"},
					"Test1": []string{"src"},
					"test":  []string{"dst"},
					"test2": []string{"dst"},
				},
			},
		),
		gen(
			"src/dst is empty",
			&condition{
				src: http.Header{},
				dst: http.Header{},
			},
			&action{
				dst: http.Header{},
			},
		),
		gen(
			"src is empty",
			&condition{
				src: http.Header{},
				dst: http.Header{
					"test": []string{"dst"},
				},
			},
			&action{
				dst: http.Header{
					"test": []string{"dst"},
				},
			},
		),
		gen(
			"dst is empty",
			&condition{
				src: http.Header{
					"test": []string{"src"},
				},
				dst: http.Header{},
			},
			&action{
				dst: http.Header{
					"Test": []string{"src"},
				},
			},
		),
		gen(
			"src/dst is nil",
			&condition{
				src: nil,
				dst: nil,
			},
			&action{
				dst: nil,
			},
		),
		gen(
			"src is nil",
			&condition{
				src: nil,
				dst: http.Header{
					"test": []string{"dst"},
				},
			},
			&action{
				dst: http.Header{
					"test": []string{"dst"},
				},
			},
		),
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			copyHeader(tt.C.dst, tt.C.src)
			testutil.Diff(t, tt.A.dst, tt.C.dst)
		})
	}
}

func TestCopyTrailer(t *testing.T) {
	type condition struct {
		src http.Header
		dst http.Header
	}

	type action struct {
		dst http.Header
	}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"copy header values",
			&condition{
				src: http.Header{
					"test":  []string{"src"},
					"test1": []string{"src"},
				},
				dst: http.Header{
					"test":  []string{"dst"},
					"test2": []string{"dst"},
				},
			},
			&action{
				dst: http.Header{
					http.TrailerPrefix + "test":  []string{"src"},
					http.TrailerPrefix + "test1": []string{"src"},
					"test":                       []string{"dst"},
					"test2":                      []string{"dst"},
				},
			},
		),
		gen(
			"src/dst is empty",
			&condition{
				src: http.Header{},
				dst: http.Header{},
			},
			&action{
				dst: http.Header{},
			},
		),
		gen(
			"src is empty",
			&condition{
				src: http.Header{},
				dst: http.Header{
					"test": []string{"dst"},
				},
			},
			&action{
				dst: http.Header{
					"test": []string{"dst"},
				},
			},
		),
		gen(
			"dst is empty",
			&condition{
				src: http.Header{
					"test": []string{"src"},
				},
				dst: http.Header{},
			},
			&action{
				dst: http.Header{
					http.TrailerPrefix + "test": []string{"src"},
				},
			},
		),
		gen(
			"src/dst is nil",
			&condition{
				src: nil,
				dst: nil,
			},
			&action{
				dst: nil,
			},
		),
		gen(
			"src is nil",
			&condition{
				src: nil,
				dst: http.Header{
					"test": []string{"dst"},
				},
			},
			&action{
				dst: http.Header{
					"test": []string{"dst"},
				},
			},
		),
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			copyTrailer(tt.C.dst, tt.C.src)
			testutil.Diff(t, tt.A.dst, tt.C.dst)
		})
	}
}

func TestHandleHopByHopHeaders(t *testing.T) {
	type condition struct {
		in  http.Header
		out http.Header
	}

	type action struct {
		outHeader http.Header
	}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"UA not exists",
			&condition{
				in:  http.Header{},
				out: http.Header{},
			},
			&action{
				outHeader: http.Header{
					"User-Agent": {""},
				},
			},
		),
		gen(
			"UA exists",
			&condition{
				in: http.Header{},
				out: http.Header{
					"User-Agent": {"test"},
				},
			},
			&action{
				outHeader: http.Header{
					"User-Agent": {"test"},
				},
			},
		),
		gen(
			"TE trailers exists",
			&condition{
				in: http.Header{
					"Te": []string{"trailers"},
				},
				out: http.Header{},
			},
			&action{
				outHeader: http.Header{
					"User-Agent": []string{""},
					"Te":         []string{"trailers"},
				},
			},
		),
		gen(
			"TE trailers exists",
			&condition{
				in: http.Header{
					"Te": []string{"trailers, deflate;q=0.5"},
				},
				out: http.Header{},
			},
			&action{
				outHeader: http.Header{
					"User-Agent": []string{""},
					"Te":         []string{"trailers"},
				},
			},
		),
		gen(
			"TE trailers not exists",
			&condition{
				in: http.Header{
					"Te": []string{"compress"},
				},
				out: http.Header{},
			},
			&action{
				outHeader: http.Header{
					"User-Agent": []string{""},
				},
			},
		),
		gen(
			"Connection upgrade exists",
			&condition{
				in: http.Header{
					"Connection": []string{"upgrade"},
					"Upgrade":    []string{"websocket"},
				},
				out: http.Header{},
			},
			&action{
				outHeader: http.Header{
					"User-Agent": []string{""},
					"Connection": []string{"Upgrade"},
					"Upgrade":    []string{"websocket"},
				},
			},
		),
		gen(
			"Connection upgrade exists",
			&condition{
				in: http.Header{
					"Connection": []string{"keep-alive, upgrade"},
					"Upgrade":    []string{"websocket"},
				},
				out: http.Header{},
			},
			&action{
				outHeader: http.Header{
					"User-Agent": []string{""},
					"Connection": []string{"Upgrade"},
					"Upgrade":    []string{"websocket"},
				},
			},
		),
		gen(
			"Connection upgrade not exists",
			&condition{
				in: http.Header{
					"Connection": []string{"keep-alive"},
					"Upgrade":    []string{"websocket"},
				},
				out: http.Header{},
			},
			&action{
				outHeader: http.Header{
					"User-Agent": []string{""},
				},
			},
		),
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			handleHopByHopHeaders(tt.C.in, tt.C.out)
			testutil.Diff(t, tt.A.outHeader, tt.C.out)
		})
	}
}

func TestRemoveHopByHopHeaders(t *testing.T) {
	type condition struct {
		header http.Header
		key    string
	}

	type action struct {
	}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"Connection",
			&condition{
				header: http.Header{
					"Connection": []string{"test"},
				},
				key: "Connection",
			},
			&action{},
		),
		gen(
			"Keep-Alive",
			&condition{
				header: http.Header{
					"Keep-Alive": []string{"test"},
				},
				key: "Keep-Alive",
			},
			&action{},
		),
		gen(
			"Proxy-Connection",
			&condition{
				header: http.Header{
					"Proxy-Connection": []string{"test"},
				},
				key: "Proxy-Connection",
			},
			&action{},
		),
		gen(
			"Proxy-Authenticate",
			&condition{
				header: http.Header{
					"Proxy-Authenticate": []string{"test"},
				},
				key: "Proxy-Authenticate",
			},
			&action{},
		),
		gen(
			"Proxy-Authorization",
			&condition{
				header: http.Header{
					"Proxy-Authorization": []string{"test"},
				},
				key: "Proxy-Authorization",
			},
			&action{},
		),
		gen(
			"Te",
			&condition{
				header: http.Header{
					"Te": []string{"test"},
				},
				key: "Te",
			},
			&action{},
		),
		gen(
			"Trailer",
			&condition{
				header: http.Header{
					"Trailer": []string{"test"},
				},
				key: "Trailer",
			},
			&action{},
		),
		gen(
			"Transfer-Encoding",
			&condition{
				header: http.Header{
					"Transfer-Encoding": []string{"test"},
				},
				key: "Transfer-Encoding",
			},
			&action{},
		),
		gen(
			"Upgrade",
			&condition{
				header: http.Header{
					"Upgrade": []string{"test"},
				},
				key: "Upgrade",
			},
			&action{},
		),
		gen(
			"Connection Foo",
			&condition{
				header: http.Header{
					"Connection": []string{"foo", "bar"},
					"Foo":        []string{"foo-val"},
					"Bar":        []string{"bar-val"},
				},
				key: "Foo",
			},
			&action{},
		),
		gen(
			"Connection Bar",
			&condition{
				header: http.Header{
					"Connection": []string{"foo", "bar"},
					"Foo":        []string{"foo-val"},
					"Bar":        []string{"bar-val"},
				},
				key: "Bar",
			},
			&action{},
		),
		gen(
			"Connection FooBar",
			&condition{
				header: http.Header{
					"Connection": []string{"foo, bar"},
					"Foo":        []string{"foo-val"},
					"Bar":        []string{"bar-val"},
				},
				key: "Foo",
			},
			&action{},
		),
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			removeHopByHopHeaders(tt.C.header)

			v, ok := tt.C.header[tt.C.key]

			testutil.Diff(t, false, ok)
			testutil.Diff(t, []string(nil), v)
		})
	}
}

func TestRewriteRequestURL(t *testing.T) {
	type condition struct {
		req    string
		target string
	}

	type action struct {
		result string
	}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"host rewritten",
			&condition{
				req:    "http://req.com",
				target: "http://target.com",
			},
			&action{
				result: "http://target.com",
			},
		),
		gen(
			"schema rewritten",
			&condition{
				req:    "http://req.com",
				target: "https://target.com",
			},
			&action{
				result: "https://target.com",
			},
		),
		gen(
			"port rewritten",
			&condition{
				req:    "http://req.com:81",
				target: "http://target.com:82",
			},
			&action{
				result: "http://target.com:82",
			},
		),
		gen(
			"path rewritten",
			&condition{
				req:    "http://req.com/foo",
				target: "https://target.com/bar",
			},
			&action{
				result: "https://target.com/bar",
			},
		),
		gen(
			"query rewritten",
			&condition{
				req:    "http://req.com/?foo=xyz",
				target: "https://target.com/?bar=abc",
			},
			&action{
				result: "https://target.com/?bar=abc&foo=xyz",
			},
		),
		gen(
			"query rewritten",
			&condition{
				req:    "http://req.com/",
				target: "https://target.com/?bar=abc",
			},
			&action{
				result: "https://target.com/?bar=abc",
			},
		),
		gen(
			"query rewritten",
			&condition{
				req:    "http://req.com/?foo=xyz",
				target: "https://target.com/",
			},
			&action{
				result: "https://target.com/?foo=xyz",
			},
		),
		gen(
			"fragment rewritten",
			&condition{
				req:    "http://req.com/#foo",
				target: "https://target.com/#bar",
			},
			&action{
				result: "https://target.com/#bar",
			},
		),
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			reqURL, err := url.Parse(tt.C.req)
			testutil.Diff(t, nil, err)

			targetURL, err := url.Parse(tt.C.target)
			testutil.Diff(t, nil, err)

			resultURL, err := url.Parse(tt.A.result)
			testutil.Diff(t, nil, err)

			rewriteRequestURL(reqURL, targetURL)

			testutil.Diff(t, resultURL, reqURL)
		})
	}
}

func TestSetXForwardedHeaders(t *testing.T) {
	type condition struct {
		req    *http.Request
		header http.Header
	}

	type action struct {
		header http.Header
	}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"non TLS",
			&condition{
				req: &http.Request{
					RemoteAddr: "127.0.0.1:80",
					Host:       "test.com",
					TLS:        nil,
					Header:     http.Header{},
				},
				header: http.Header{},
			},
			&action{
				http.Header{
					"X-Forwarded-For":   []string{"127.0.0.1"},
					"X-Forwarded-Port":  []string{"80"},
					"X-Forwarded-Host":  []string{"test.com"},
					"X-Forwarded-Proto": []string{"http"},
				},
			},
		),
		gen(
			"ipv4/https",
			&condition{
				req: &http.Request{
					RemoteAddr: "127.0.0.1:80",
					Host:       "test.com",
					TLS:        &tls.ConnectionState{},
					Header:     http.Header{},
				},
				header: http.Header{},
			},
			&action{
				http.Header{
					"X-Forwarded-For":   []string{"127.0.0.1"},
					"X-Forwarded-Port":  []string{"80"},
					"X-Forwarded-Host":  []string{"test.com"},
					"X-Forwarded-Proto": []string{"https"},
				},
			},
		),
		gen(
			"ipv6/https",
			&condition{
				req: &http.Request{
					RemoteAddr: "[::1]:80",
					Host:       "test.com",
					TLS:        &tls.ConnectionState{},
					Header:     http.Header{},
				},
				header: http.Header{},
			},
			&action{
				http.Header{
					"X-Forwarded-For":   []string{"::1"},
					"X-Forwarded-Port":  []string{"80"},
					"X-Forwarded-Host":  []string{"test.com"},
					"X-Forwarded-Proto": []string{"https"},
				},
			},
		),
		gen(
			"prior X-Forwarded-For",
			&condition{
				req: &http.Request{
					RemoteAddr: "127.0.0.1:80",
					Host:       "test.com",
					TLS:        nil,
					Header: http.Header{
						"X-Forwarded-For": []string{"::1"},
					},
				},
				header: http.Header{},
			},
			&action{
				http.Header{
					"X-Forwarded-For":   []string{"::1, 127.0.0.1"},
					"X-Forwarded-Port":  []string{"80"},
					"X-Forwarded-Host":  []string{"test.com"},
					"X-Forwarded-Proto": []string{"http"},
				},
			},
		),
		gen(
			"prior X-Forwarded-Port",
			&condition{
				req: &http.Request{
					RemoteAddr: "127.0.0.1:80",
					Host:       "test.com",
					TLS:        nil,
					Header: http.Header{
						"X-Forwarded-Port": []string{"12345"},
					},
				},
				header: http.Header{},
			},
			&action{
				http.Header{
					"X-Forwarded-For":   []string{"127.0.0.1"},
					"X-Forwarded-Port":  []string{"80"},
					"X-Forwarded-Host":  []string{"test.com"},
					"X-Forwarded-Proto": []string{"http"},
				},
			},
		),
		gen(
			"prior X-Forwarded-Host",
			&condition{
				req: &http.Request{
					RemoteAddr: "127.0.0.1:80",
					Host:       "test.com",
					TLS:        nil,
					Header: http.Header{
						"X-Forwarded-Host": []string{"prior.com"},
					},
				},
				header: http.Header{},
			},
			&action{
				http.Header{
					"X-Forwarded-For":   []string{"127.0.0.1"},
					"X-Forwarded-Port":  []string{"80"},
					"X-Forwarded-Host":  []string{"test.com"}, // Won't be added.
					"X-Forwarded-Proto": []string{"http"},
				},
			},
		),
		gen(
			"prior X-Forwarded-Proto",
			&condition{
				req: &http.Request{
					RemoteAddr: "127.0.0.1:80",
					Host:       "test.com",
					TLS:        nil,
					Header: http.Header{
						"X-Forwarded-Proto": []string{"https"},
					},
				},
				header: http.Header{},
			},
			&action{
				http.Header{
					"X-Forwarded-For":   []string{"127.0.0.1"},
					"X-Forwarded-Port":  []string{"80"},
					"X-Forwarded-Host":  []string{"test.com"},
					"X-Forwarded-Proto": []string{"http"}, // Won't be added.
				},
			},
		),
		gen(
			"invalid remote address",
			&condition{
				req: &http.Request{
					RemoteAddr: "invalid address",
					Host:       "test.com",
					TLS:        nil,
					Header:     http.Header{},
				},
				header: http.Header{
					"X-Forwarded-For": []string{"::1"},
				},
			},
			&action{
				http.Header{
					"X-Forwarded-Host":  []string{"test.com"},
					"X-Forwarded-Proto": []string{"http"},
				},
			},
		),
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			setXForwardedHeaders(tt.C.req, tt.C.header)
			testutil.Diff(t, tt.A.header, tt.C.header)
		})
	}
}
