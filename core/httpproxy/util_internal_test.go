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

	cndConnectionExists := "Connection header exists"
	cndUpgradeExists := "Upgrade value exists"
	actCheckProtocol := "check non-empty protocol"
	actCheckEmpty := "check empty string"

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	tb.Condition(cndConnectionExists, "connection header is exist in the header")
	tb.Condition(cndUpgradeExists, "upgrade value exists in the header")
	tb.Action(actCheckProtocol, "check the returned string which is not empty")
	tb.Action(actCheckEmpty, "check that the returned string is empty")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"header is empty",
			[]string{},
			[]string{actCheckEmpty},
			&condition{
				h: http.Header{},
			},
			&action{
				typ: "",
			},
		),
		gen(
			"value is not Upgrade",
			[]string{cndConnectionExists},
			[]string{actCheckEmpty},
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
			[]string{cndConnectionExists},
			[]string{actCheckEmpty},
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
			[]string{cndConnectionExists, cndUpgradeExists},
			[]string{actCheckProtocol},
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
			[]string{cndConnectionExists, cndUpgradeExists},
			[]string{actCheckProtocol},
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

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			typ := upgradeType(tt.C().h)
			testutil.Diff(t, tt.A().typ, typ)
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

	cndSrcEmpty := "empty source header"
	cndDstEmpty := "empty destination header"
	cndSrcNil := "source header is nil"
	actCheckCopied := "check copied header values"

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	tb.Condition(cndSrcEmpty, "source header has no value")
	tb.Condition(cndDstEmpty, "destination header has no value")
	tb.Condition(cndSrcNil, "source header is nil")
	tb.Action(actCheckCopied, "check that the returned string is empty")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"copy header values",
			[]string{},
			[]string{actCheckCopied},
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
			[]string{cndSrcEmpty, cndDstEmpty},
			[]string{actCheckCopied},
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
			[]string{cndSrcEmpty},
			[]string{actCheckCopied},
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
			[]string{cndDstEmpty},
			[]string{actCheckCopied},
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
			[]string{cndSrcNil},
			[]string{actCheckCopied},
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
			[]string{cndSrcNil},
			[]string{actCheckCopied},
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

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			copyHeader(tt.C().dst, tt.C().src)
			testutil.Diff(t, tt.A().dst, tt.C().dst)
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

	cndSrcEmpty := "empty source header"
	cndDstEmpty := "empty destination header"
	cndSrcNil := "source header is nil"
	actCheckCopied := "check copied header values"

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	tb.Condition(cndSrcEmpty, "source header has no value")
	tb.Condition(cndDstEmpty, "destination header has no value")
	tb.Condition(cndSrcNil, "source header is nil")
	tb.Action(actCheckCopied, "check that the returned string is empty")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"copy header values",
			[]string{},
			[]string{actCheckCopied},
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
			[]string{cndSrcEmpty, cndDstEmpty},
			[]string{actCheckCopied},
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
			[]string{cndSrcEmpty},
			[]string{actCheckCopied},
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
			[]string{cndDstEmpty},
			[]string{actCheckCopied},
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
			[]string{cndSrcNil},
			[]string{actCheckCopied},
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
			[]string{cndSrcNil},
			[]string{actCheckCopied},
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

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			copyTrailer(tt.C().dst, tt.C().src)
			testutil.Diff(t, tt.A().dst, tt.C().dst)
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

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"UA not exists",
			[]string{},
			[]string{},
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
			[]string{},
			[]string{},
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
			[]string{},
			[]string{},
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
			[]string{},
			[]string{},
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
			[]string{},
			[]string{},
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
			[]string{},
			[]string{},
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
			[]string{},
			[]string{},
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
			[]string{},
			[]string{},
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

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			handleHopByHopHeaders(tt.C().in, tt.C().out)
			testutil.Diff(t, tt.A().outHeader, tt.C().out)
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

	actCheckRemoved := "check copied header values"

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	tb.Action(actCheckRemoved, "check that the header value was removed")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"Connection",
			[]string{},
			[]string{actCheckRemoved},
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
			[]string{},
			[]string{actCheckRemoved},
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
			[]string{},
			[]string{actCheckRemoved},
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
			[]string{},
			[]string{actCheckRemoved},
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
			[]string{},
			[]string{actCheckRemoved},
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
			[]string{},
			[]string{actCheckRemoved},
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
			[]string{},
			[]string{actCheckRemoved},
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
			[]string{},
			[]string{actCheckRemoved},
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
			[]string{},
			[]string{actCheckRemoved},
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
			[]string{},
			[]string{actCheckRemoved},
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
			[]string{},
			[]string{actCheckRemoved},
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
			[]string{},
			[]string{actCheckRemoved},
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

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			removeHopByHopHeaders(tt.C().header)

			v, ok := tt.C().header[tt.C().key]

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

	cndWithPort := "url with port"
	cndWithPath := "url with path"
	cndWithReqQuery := "request has url query"
	cndWithTargetQuery := "target has url query"
	cndWithFragment := "target has fragment"
	actCheckResult := "check the rewritten url"

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	tb.Condition(cndWithPort, "request/target url have port number")
	tb.Condition(cndWithPath, "request/target url have path")
	tb.Condition(cndWithReqQuery, "request url has queries")
	tb.Condition(cndWithTargetQuery, "target url has queries")
	tb.Condition(cndWithFragment, "target url has fragment")
	tb.Action(actCheckResult, "check result of rewritten url")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"host rewritten",
			[]string{},
			[]string{actCheckResult},
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
			[]string{},
			[]string{actCheckResult},
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
			[]string{cndWithPort},
			[]string{actCheckResult},
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
			[]string{cndWithPath},
			[]string{actCheckResult},
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
			[]string{cndWithReqQuery, cndWithTargetQuery},
			[]string{actCheckResult},
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
			[]string{cndWithTargetQuery},
			[]string{actCheckResult},
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
			[]string{cndWithReqQuery},
			[]string{actCheckResult},
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
			[]string{cndWithFragment},
			[]string{actCheckResult},
			&condition{
				req:    "http://req.com/#foo",
				target: "https://target.com/#bar",
			},
			&action{
				result: "https://target.com/#bar",
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			reqURL, err := url.Parse(tt.C().req)
			testutil.Diff(t, nil, err)

			targetURL, err := url.Parse(tt.C().target)
			testutil.Diff(t, nil, err)

			resultURL, err := url.Parse(tt.A().result)
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

	cndTLS := "TLS"
	cndInvalidAddress := "invalid address"
	cndPriorXForwardedFor := "prior X-Forwarded-For"
	cndPriorXForwardedPort := "prior X-Forwarded-Port"
	cndPriorXForwardedHost := "prior X-Forwarded-Host"
	cndPriorXForwardedProto := "prior X-Forwarded-Proto"
	actCheckHeaders := "check forward headers"

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	tb.Condition(cndTLS, "TLS request")
	tb.Condition(cndPriorXForwardedFor, "X-Forwarded-For header is in the request")
	tb.Condition(cndPriorXForwardedPort, "X-Forwarded-Port header is in the request")
	tb.Condition(cndPriorXForwardedHost, "X-Forwarded-Host header is in the request")
	tb.Condition(cndPriorXForwardedProto, "X-Forwarded-Proto header is in the request")
	tb.Condition(cndInvalidAddress, "remote address is invalid")
	tb.Action(actCheckHeaders, "check result of rewritten url")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"non TLS",
			[]string{},
			[]string{actCheckHeaders},
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
			[]string{cndTLS},
			[]string{actCheckHeaders},
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
			[]string{cndTLS},
			[]string{actCheckHeaders},
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
			[]string{cndTLS, cndPriorXForwardedFor},
			[]string{actCheckHeaders},
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
			[]string{cndTLS, cndPriorXForwardedPort},
			[]string{actCheckHeaders},
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
			[]string{cndTLS, cndPriorXForwardedHost},
			[]string{actCheckHeaders},
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
			[]string{cndTLS, cndPriorXForwardedProto},
			[]string{actCheckHeaders},
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
			[]string{cndTLS, cndInvalidAddress},
			[]string{actCheckHeaders},
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

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			setXForwardedHeaders(tt.C().req, tt.C().header)
			testutil.Diff(t, tt.A().header, tt.C().header)
		})
	}
}
