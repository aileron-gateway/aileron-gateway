package httpproxy

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"io"
	"maps"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"regexp"
	"testing"

	"github.com/aileron-gateway/aileron-gateway/core"
	"github.com/aileron-gateway/aileron-gateway/kernel/log"
	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
	utilhttp "github.com/aileron-gateway/aileron-gateway/util/http"
	"github.com/aileron-gateway/aileron-gateway/util/resilience"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/quic-go/quic-go/http3"
)

type testErrorHandler struct {
	called bool
	err    error
}

func (eh *testErrorHandler) ServeHTTPError(w http.ResponseWriter, r *http.Request, e error) {
	eh.called = true
	eh.err = e
	if c, ok := e.(interface{ StatusCode() int }); ok {
		w.WriteHeader(c.StatusCode())
	}
}

type testRoundTripper struct {
	http.RoundTripper

	status       int
	header       http.Header
	trailer      http.Header
	extraTrailer http.Header
	body         io.ReadCloser

	response *http.Response
	err      error
}

func (rt *testRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	w := &http.Response{
		StatusCode: rt.status,
		Header:     rt.header,
		Trailer:    rt.trailer,
		Request:    r,
	}
	if rt.body == nil {
		// According to the *http.Response, response body is always non-nil.
		rt.body = io.NopCloser(bytes.NewReader(nil))
	}
	if rt.extraTrailer != nil {
		rt.body = &testBody{
			ReadCloser: rt.body,
			afterEOF: func() {
				maps.Copy(w.Trailer, rt.extraTrailer)
			},
		}
	}
	w.Body = rt.body
	rt.response = w
	return w, rt.err
}

type testBody struct {
	io.ReadCloser
	readErr  error
	closed   bool
	afterEOF func()
}

func (b *testBody) Read(p []byte) (int, error) {
	if b.readErr != nil {
		return 0, b.readErr
	}
	n, err := b.ReadCloser.Read(p)
	if err == io.EOF && b.afterEOF != nil {
		b.afterEOF()
	}
	return n, err
}

func (b *testBody) Close() error {
	b.closed = true
	return b.ReadCloser.Close()
}

func TestReverseProxy_ServeHTTP(t *testing.T) {
	type condition struct {
		proxy     *reverseProxy
		reqHeader http.Header
		reqBody   *testBody
	}

	type action struct {
		status     int
		written    string
		header     map[string][]string
		trailer    map[string][]string
		err        any // error or errorutil.Kind
		errPattern *regexp.Regexp
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	testErr := errors.New("test error")
	testOpErr := &net.OpError{Op: "write", Err: testErr}

	ups := &noopUpstream{
		weight:    1,
		rawURL:    "http://upstream.com/proxy",
		parsedURL: &url.URL{Scheme: "http", Host: "upstream.com", Path: "/proxy"},
	}
	rlb := &resilience.RoundRobinLB[upstream]{}
	rlb.Add(ups)

	inactive := &lbUpstream{
		circuitBreaker: &testCircuitBreaker{
			activeStatus: false,
		},
	}
	inactiveRlb := &resilience.RoundRobinLB[upstream]{}
	inactiveRlb.Add(inactive)

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"proxy success",
			[]string{},
			[]string{},
			&condition{
				proxy: &reverseProxy{
					lg: log.GlobalLogger(log.DefaultLoggerName),
					eh: &testErrorHandler{},
					lbs: []loadBalancer{
						&nonHashLB{
							lbMatcher: &lbMatcher{
								pathMatchers: []matcherFunc{func(string) (string, bool) { return "", true }},
							},
							LoadBalancer: rlb,
						},
					},
					rt: &testRoundTripper{
						status: http.StatusOK,
						body:   io.NopCloser(bytes.NewReader([]byte("test"))),
					},
				},
			},
			&action{
				status:  http.StatusOK,
				written: "test",
				err:     nil,
			},
		),
		gen(
			"proxy success with request body",
			[]string{},
			[]string{},
			&condition{
				proxy: &reverseProxy{
					lg: log.GlobalLogger(log.DefaultLoggerName),
					eh: &testErrorHandler{},
					lbs: []loadBalancer{
						&nonHashLB{
							lbMatcher: &lbMatcher{
								pathMatchers: []matcherFunc{func(string) (string, bool) { return "", true }},
							},
							LoadBalancer: rlb,
						},
					},
					rt: &testRoundTripper{
						status: http.StatusOK,
						header: http.Header{},
						body:   io.NopCloser(bytes.NewReader([]byte("test"))),
					},
				},
				reqBody: &testBody{
					ReadCloser: io.NopCloser(bytes.NewReader([]byte("test request"))),
				},
			},
			&action{
				status:  http.StatusOK,
				written: "test",
				err:     nil,
			},
		),
		gen(
			"proxy success with response header",
			[]string{},
			[]string{},
			&condition{
				proxy: &reverseProxy{
					lg: log.GlobalLogger(log.DefaultLoggerName),
					eh: &testErrorHandler{},
					lbs: []loadBalancer{
						&nonHashLB{
							lbMatcher: &lbMatcher{
								pathMatchers: []matcherFunc{func(string) (string, bool) { return "", true }},
							},
							LoadBalancer: rlb,
						},
					},
					rt: &testRoundTripper{
						status: http.StatusOK,
						header: http.Header{
							"Test-Header-Foo": {"foo1", "foo2"},
							"Test-Header-Bar": {"bar1", "bar2"},
						},
						body: io.NopCloser(bytes.NewReader([]byte("test"))),
					},
				},
			},
			&action{
				status:  http.StatusOK,
				written: "test",
				header: map[string][]string{
					"Test-Header-Foo": {"foo1", "foo2"},
					"Test-Header-Bar": {"bar1", "bar2"},
				},
				err: nil,
			},
		),
		gen(
			"proxy success with announced trailer",
			[]string{},
			[]string{},
			&condition{
				proxy: &reverseProxy{
					lg: log.GlobalLogger(log.DefaultLoggerName),
					eh: &testErrorHandler{},
					lbs: []loadBalancer{
						&nonHashLB{
							lbMatcher: &lbMatcher{
								pathMatchers: []matcherFunc{func(string) (string, bool) { return "", true }},
							},
							LoadBalancer: rlb,
						},
					},
					rt: &testRoundTripper{
						status: http.StatusOK,
						header: http.Header{},
						trailer: http.Header{
							"Trailer-Foo": {"foo1", "foo2"},
						},
						body: io.NopCloser(bytes.NewReader([]byte("test"))),
					},
				},
			},
			&action{
				status:  http.StatusOK,
				written: "test",
				header: map[string][]string{
					"Trailer": {"Trailer-Foo"},
				},
				trailer: map[string][]string{
					"Trailer-Foo": {"foo1", "foo2"},
				},
				err: nil,
			},
		),
		gen(
			"proxy success with not announced trailers",
			[]string{},
			[]string{},
			&condition{
				proxy: &reverseProxy{
					lg: log.GlobalLogger(log.DefaultLoggerName),
					eh: &testErrorHandler{},
					lbs: []loadBalancer{
						&nonHashLB{
							lbMatcher: &lbMatcher{
								pathMatchers: []matcherFunc{func(string) (string, bool) { return "", true }},
							},
							LoadBalancer: rlb,
						},
					},
					rt: &testRoundTripper{
						status:  http.StatusOK,
						header:  http.Header{},
						trailer: http.Header{},
						extraTrailer: http.Header{
							"Trailer-Foo": {"foo1", "foo2"},
							"Trailer-Bar": {"bar1", "bar2"},
						},
						body: io.NopCloser(bytes.NewReader([]byte("test"))),
					},
				},
			},
			&action{
				status:  http.StatusOK,
				written: "test",
				header:  map[string][]string{},
				trailer: map[string][]string{
					"Trailer-Foo": {"foo1", "foo2"},
					"Trailer-Bar": {"bar1", "bar2"},
				},
				err: nil,
			},
		),
		gen(
			"upstream not found",
			[]string{},
			[]string{},
			&condition{
				proxy: &reverseProxy{
					lg: log.GlobalLogger(log.DefaultLoggerName),
					eh: &testErrorHandler{},
					lbs: []loadBalancer{
						&nonHashLB{
							lbMatcher: &lbMatcher{
								pathMatchers: []matcherFunc{func(string) (string, bool) { return "", false }},
							},
							LoadBalancer: rlb,
						},
					},
					rt: &testRoundTripper{},
				},
			},
			&action{
				status:     http.StatusNotFound,
				written:    "",
				err:        core.ErrCoreProxyNoUpstream,
				errPattern: regexp.MustCompile(``),
			},
		),
		gen(
			"upstream not available",
			[]string{},
			[]string{},
			&condition{
				proxy: &reverseProxy{
					lg: log.GlobalLogger(log.DefaultLoggerName),
					eh: &testErrorHandler{},
					lbs: []loadBalancer{
						&nonHashLB{
							lbMatcher: &lbMatcher{
								pathMatchers: []matcherFunc{func(string) (string, bool) { return "", true }},
							},
							LoadBalancer: inactiveRlb,
						},
					},
					rt: &testRoundTripper{},
				},
			},
			&action{
				status:     http.StatusBadGateway,
				written:    "",
				err:        core.ErrCoreProxyUnavailable,
				errPattern: regexp.MustCompile(``),
			},
		),
		gen(
			"round trip error",
			[]string{},
			[]string{},
			&condition{
				proxy: &reverseProxy{
					lg: log.GlobalLogger(log.DefaultLoggerName),
					eh: &testErrorHandler{},
					lbs: []loadBalancer{
						&nonHashLB{
							lbMatcher: &lbMatcher{
								pathMatchers: []matcherFunc{func(string) (string, bool) { return "", true }},
							},
							LoadBalancer: rlb,
						},
					},
					rt: &testRoundTripper{
						status: 0,
						err:    testErr,
					},
				},
			},
			&action{
				status:     http.StatusInternalServerError,
				written:    "",
				err:        testErr,
				errPattern: regexp.MustCompile(`test error`),
			},
		),
		gen(
			"protocol upgrade error",
			[]string{},
			[]string{},
			&condition{
				proxy: &reverseProxy{
					lg: log.GlobalLogger(log.DefaultLoggerName),
					eh: &testErrorHandler{},
					lbs: []loadBalancer{
						&nonHashLB{
							lbMatcher: &lbMatcher{
								pathMatchers: []matcherFunc{func(string) (string, bool) { return "", true }},
							},
							LoadBalancer: rlb,
						},
					},
					rt: &testRoundTripper{
						status: http.StatusSwitchingProtocols,
						header: http.Header{},
					},
				},
			},
			&action{
				status:     http.StatusInternalServerError,
				written:    "",
				err:        core.ErrCoreProxyProtocolSwitch,
				errPattern: regexp.MustCompile(`failed to upgrade protocol. internal error`),
			},
		),
		gen(
			"copy response body/context canceled",
			[]string{},
			[]string{},
			&condition{
				proxy: &reverseProxy{
					lg: log.GlobalLogger(log.DefaultLoggerName),
					eh: &testErrorHandler{},
					lbs: []loadBalancer{
						&nonHashLB{
							lbMatcher: &lbMatcher{
								pathMatchers: []matcherFunc{func(string) (string, bool) { return "", true }},
							},
							LoadBalancer: rlb,
						},
					},
					rt: &testRoundTripper{
						status: http.StatusOK,
						body: &testBody{
							ReadCloser: io.NopCloser(bytes.NewReader(nil)),
							readErr:    context.Canceled,
						},
					},
				},
			},
			&action{
				status:     http.StatusOK,
				written:    "",
				err:        context.Canceled,
				errPattern: regexp.MustCompile(`context canceled`),
			},
		),
		gen(
			"copy response body/write error",
			[]string{},
			[]string{},
			&condition{
				proxy: &reverseProxy{
					lg: log.GlobalLogger(log.DefaultLoggerName),
					eh: &testErrorHandler{},
					lbs: []loadBalancer{
						&nonHashLB{
							lbMatcher: &lbMatcher{
								pathMatchers: []matcherFunc{func(string) (string, bool) { return "", true }},
							},
							LoadBalancer: rlb,
						},
					},
					rt: &testRoundTripper{
						status: http.StatusOK,
						body: &testBody{
							ReadCloser: io.NopCloser(bytes.NewReader(nil)),
							readErr:    testOpErr,
						},
					},
				},
			},
			&action{
				status:     http.StatusOK,
				written:    "",
				err:        testOpErr,
				errPattern: regexp.MustCompile(`write: test error`),
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			r, _ := http.NewRequest(http.MethodPost, "http://test.com/foo", nil)
			r.Header = tt.C().reqHeader
			if tt.C().reqBody != nil {
				r.Body = tt.C().reqBody
			}

			w := httptest.NewRecorder()
			tt.C().proxy.ServeHTTP(w, r)

			if tt.C().reqBody != nil {
				// Requests body should always be closed if not nil.
				testutil.Diff(t, true, tt.C().reqBody.closed)
			}

			// Check written body.
			resp := w.Result().Body
			b, _ := io.ReadAll(resp)
			testutil.Diff(t, tt.A().written, string(b))

			// Check response status code.
			testutil.Diff(t, tt.A().status, w.Result().StatusCode)

			// Check headers and trailers.
			wh := w.Result().Header
			wt := w.Result().Trailer
			for k, v := range tt.A().header {
				testutil.Diff(t, v, wh[k], cmpopts.SortSlices(func(x, y string) bool { return x > y }))
			}
			for k, v := range tt.A().trailer {
				testutil.Diff(t, v, wt[k], cmpopts.SortSlices(func(x, y string) bool { return x > y }))
			}

			eh := tt.C().proxy.eh.(*testErrorHandler)
			if e, ok := eh.err.(*utilhttp.HTTPError); ok {
				testutil.DiffError(t, tt.A().err, tt.A().errPattern, e.Unwrap(), cmpopts.EquateErrors())
			} else {
				testutil.DiffError(t, tt.A().err, tt.A().errPattern, eh.err)
			}
		})
	}
}

func TestProxyErrorResponse(t *testing.T) {
	type condition struct {
		err error
	}

	type action struct {
		er core.HTTPError
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	cndContextCanceled := tb.Condition("context canceled", "input context.Canceled error")
	cndRequestCanceled := tb.Condition("request canceled", "input http3 error of request canceled")
	cndDeadlineExceeded := tb.Condition("deadline exceeded", "input context.DeadlineExceeded error")
	actCheckLogOnly := tb.Action("log only", "check the value written to a writer")
	actCheckServerError := tb.Action("server error", "check that the expected non-nil error returned")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"context.Canceled",
			[]string{cndContextCanceled},
			[]string{actCheckLogOnly},
			&condition{
				err: context.Canceled,
			},
			&action{
				er: utilhttp.NewHTTPError(context.Canceled, -1),
			},
		),
		gen(
			"H3_REQUEST_CANCELLED",
			[]string{cndRequestCanceled},
			[]string{actCheckLogOnly},
			&condition{
				err: &http3.Error{
					ErrorCode: http3.ErrCodeRequestCanceled,
				},
			},
			&action{
				er: utilhttp.NewHTTPError(context.Canceled, -1),
			},
		),
		gen(
			"context.DeadlineExceeded",
			[]string{cndDeadlineExceeded},
			[]string{actCheckLogOnly},
			&condition{
				err: context.DeadlineExceeded,
			},
			&action{
				er: utilhttp.NewHTTPError(context.DeadlineExceeded, http.StatusGatewayTimeout),
			},
		),
		gen(
			"other errors",
			[]string{},
			[]string{actCheckServerError},
			&condition{
				err: io.ErrClosedPipe,
			},
			&action{
				er: utilhttp.NewHTTPError(io.ErrClosedPipe, http.StatusInternalServerError),
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			er := proxyErrorResponse(tt.C().err)

			opts := []cmp.Option{
				cmp.AllowUnexported(utilhttp.HTTPError{}),
				cmpopts.IgnoreFields(utilhttp.HTTPError{}, "inner"),
			}
			testutil.Diff(t, tt.A().er, er, opts...)

			err, ok := er.(*utilhttp.HTTPError)
			testutil.Diff(t, true, ok)
			testutil.Diff(t, tt.C().err, err.Unwrap(), cmpopts.EquateErrors())
		})
	}
}

type testConn struct {
	net.Conn
	reader   io.Reader
	writer   io.Writer
	closeErr error
}

func (c *testConn) Read(b []byte) (n int, err error) {
	return c.reader.Read(b)
}

func (c *testConn) Write(p []byte) (n int, err error) {
	return c.writer.Write(p)
}

func (c *testConn) Close() error {
	return c.closeErr
}

type testHijacker struct {
	http.ResponseWriter
	conn net.Conn
	rw   *bufio.ReadWriter
	err  error
}

func (h *testHijacker) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return h.conn, h.rw, h.err
}

func TestHandleUpgradeResponse(t *testing.T) {
	type condition struct {
		rw  http.ResponseWriter
		req *http.Request
		res *http.Response
	}

	type action struct {
		err        any // error or errorutil.Kind
		errPattern *regexp.Regexp
		httpErr    core.HTTPError
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"upgrade type mismatch",
			[]string{},
			[]string{},
			&condition{
				rw: &testResponse{},
				req: &http.Request{
					Header: http.Header{
						"Connection": {"Upgrade"},
						"Upgrade":    {"foo"},
					},
				},
				res: &http.Response{
					Header: http.Header{
						"Connection": {"Upgrade"},
						"Upgrade":    {"bar"},
					},
				},
			},
			&action{
				err:        core.ErrCoreProxyProtocolSwitch,
				errPattern: regexp.MustCompile(core.ErrPrefix + `failed to upgrade protocol. backend tried to switch protocol`),
				httpErr:    utilhttp.NewHTTPError(nil, http.StatusInternalServerError),
			},
		),
		gen(
			"response body not ReadWriteCloser",
			[]string{},
			[]string{},
			&condition{
				rw:  &testResponse{},
				req: &http.Request{},
				res: &http.Response{
					Body: io.NopCloser(bytes.NewReader([]byte("test"))),
				},
			},
			&action{
				err:        core.ErrCoreProxyProtocolSwitch,
				errPattern: regexp.MustCompile(core.ErrPrefix + `failed to upgrade protocol. internal error`),
				httpErr:    utilhttp.NewHTTPError(nil, http.StatusInternalServerError),
			},
		),
		gen(
			"non-Hijacker",
			[]string{},
			[]string{},
			&condition{
				rw:  &testResponse{},
				req: &http.Request{},
				res: &http.Response{
					Body: &testConn{},
				},
			},
			&action{
				err:        core.ErrCoreProxyProtocolSwitch,
				errPattern: regexp.MustCompile(core.ErrPrefix + `failed to upgrade protocol. Hijack failed from`),
				httpErr:    utilhttp.NewHTTPError(nil, http.StatusInternalServerError),
			},
		),
		gen(
			"Hijacker",
			[]string{},
			[]string{},
			&condition{
				rw: &testHijacker{
					ResponseWriter: &testResponse{},
					conn: &testConn{
						reader: bytes.NewReader([]byte("test")),
						writer: bytes.NewBuffer(make([]byte, 100)),
					},
					rw: bufio.NewReadWriter(bufio.NewReader(bytes.NewReader([]byte("test"))), bufio.NewWriter(bytes.NewBuffer(make([]byte, 100)))),
				},
				req: &http.Request{},
				res: &http.Response{
					Body: &testConn{
						reader: bytes.NewReader([]byte("test")),
						writer: bytes.NewBuffer(make([]byte, 100)),
					},
				},
			},
			&action{
				err:     nil,
				httpErr: nil,
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			httpErr := handleUpgradeResponse(tt.C().rw, tt.C().req, tt.C().res)

			opts := []cmp.Option{
				cmp.AllowUnexported(utilhttp.HTTPError{}),
				cmpopts.IgnoreFields(utilhttp.HTTPError{}, "inner"),
			}
			testutil.Diff(t, tt.A().httpErr, httpErr, opts...)
			if tt.A().httpErr != nil {
				err, ok := httpErr.(*utilhttp.HTTPError)
				testutil.Diff(t, true, ok)
				testutil.DiffError(t, tt.A().err, tt.A().errPattern, err.Unwrap())
			}
		})
	}
}
