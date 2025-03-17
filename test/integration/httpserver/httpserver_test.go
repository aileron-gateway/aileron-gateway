//go:build integration

package httpserver_test

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/aileron-gateway/aileron-gateway/apis/kernel"
	"github.com/aileron-gateway/aileron-gateway/cmd/aileron/app"
	"github.com/aileron-gateway/aileron-gateway/core"
	"github.com/aileron-gateway/aileron-gateway/kernel/api"
	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
	"github.com/aileron-gateway/aileron-gateway/test/integration/common"
	"github.com/quic-go/quic-go/http3"
	"golang.org/x/net/http2"
)

type testHandler struct {
	called int
	sleep  int
}

func (h *testHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.called += 1
	time.Sleep(time.Millisecond * time.Duration(h.sleep))
	w.Header().Set("test", "ok")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("test"))
}

type testMiddleware struct {
	called int
}

func (m *testMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m.called += 1
		next.ServeHTTP(w, r)
	})
}

func TestAddr(t *testing.T) {

	configs := []string{
		testDataDir + "config-addr.yaml",
	}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "core/v1",
		Kind:       "HTTPServer",
	}
	runner, err := api.ReferTypedObject[core.Runner](server, ref)
	testutil.DiffError(t, nil, nil, err)

	ctx, cancel := context.WithCancel(context.Background())
	waitStop := make(chan struct{})
	defer func() {
		cancel()
		<-waitStop
	}()
	go func() {
		err := runner.Run(ctx)
		testutil.DiffError(t, nil, nil, err)
		close(waitStop)
	}()
	time.Sleep(time.Second) // Wait a little until the server starts.

	h1t := &http.Transport{
		ForceAttemptHTTP2: false,
	}

	req := httptest.NewRequest(http.MethodGet, "http://localhost:12345/test", nil)
	res, err := h1t.RoundTrip(req)
	testutil.DiffError(t, nil, nil, err)
	b, _ := io.ReadAll(res.Body)
	testutil.Diff(t, http.StatusNotFound, res.StatusCode)
	testutil.Diff(t, `{"status":404,"statusText":"Not Found"}`, string(b))

}

func TestDebugExpvar(t *testing.T) {

	configs := []string{
		testDataDir + "config-debug-expvar.yaml",
	}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "core/v1",
		Kind:       "HTTPServer",
	}
	runner, err := api.ReferTypedObject[core.Runner](server, ref)
	testutil.DiffError(t, nil, nil, err)

	ctx, cancel := context.WithCancel(context.Background())
	waitStop := make(chan struct{})
	defer func() {
		cancel()
		<-waitStop
	}()
	go func() {
		err := runner.Run(ctx)
		testutil.DiffError(t, nil, nil, err)
		close(waitStop)
	}()
	time.Sleep(time.Second) // Wait a little until the server starts.

	h1t := &http.Transport{
		ForceAttemptHTTP2: false,
	}

	r1 := httptest.NewRequest(http.MethodGet, "http://127.0.0.1:12345/debug/test", nil)
	w1, err := h1t.RoundTrip(r1)
	testutil.DiffError(t, nil, nil, err)
	testutil.Diff(t, http.StatusOK, w1.StatusCode)

}

func TestDebugProfile(t *testing.T) {

	configs := []string{
		testDataDir + "config-debug-profile.yaml",
	}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "core/v1",
		Kind:       "HTTPServer",
	}
	runner, err := api.ReferTypedObject[core.Runner](server, ref)
	testutil.DiffError(t, nil, nil, err)

	ctx, cancel := context.WithCancel(context.Background())
	waitStop := make(chan struct{})
	defer func() {
		cancel()
		<-waitStop
	}()
	go func() {
		err := runner.Run(ctx)
		testutil.DiffError(t, nil, nil, err)
		close(waitStop)
	}()
	time.Sleep(time.Second) // Wait a little until the server starts.

	h1t := &http.Transport{
		ForceAttemptHTTP2: false,
	}

	// Index
	r1 := httptest.NewRequest(http.MethodGet, "http://127.0.0.1:12345/debug/test/", nil)
	w1, err := h1t.RoundTrip(r1)
	testutil.DiffError(t, nil, nil, err)
	testutil.Diff(t, http.StatusOK, w1.StatusCode)

	// Cmdline
	r2 := httptest.NewRequest(http.MethodGet, "http://127.0.0.1:12345/debug/test/cmdline", nil)
	w2, err := h1t.RoundTrip(r2)
	testutil.DiffError(t, nil, nil, err)
	testutil.Diff(t, http.StatusOK, w2.StatusCode)

	// Profile
	r3 := httptest.NewRequest(http.MethodGet, "http://127.0.0.1:12345/debug/test/heap", nil)
	w3, err := h1t.RoundTrip(r3)
	testutil.DiffError(t, nil, nil, err)
	testutil.Diff(t, http.StatusOK, w3.StatusCode)

	// Symbol
	r4 := httptest.NewRequest(http.MethodGet, "http://127.0.0.1:12345/debug/test/symbol", nil)
	w4, err := h1t.RoundTrip(r4)
	testutil.DiffError(t, nil, nil, err)
	testutil.Diff(t, http.StatusOK, w4.StatusCode)

	// Trace
	r5 := httptest.NewRequest(http.MethodGet, "http://127.0.0.1:12345/debug/test/trace", nil)
	w5, err := h1t.RoundTrip(r5)
	testutil.DiffError(t, nil, nil, err)
	testutil.Diff(t, http.StatusOK, w5.StatusCode)

}

func TestHTTP(t *testing.T) {

	configs := []string{
		testDataDir + "config-http.yaml",
	}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "core/v1",
		Kind:       "HTTPServer",
	}
	runner, err := api.ReferTypedObject[core.Runner](server, ref)
	testutil.DiffError(t, nil, nil, err)

	ctx, cancel := context.WithCancel(context.Background())
	waitStop := make(chan struct{})
	defer func() {
		cancel()
		<-waitStop
	}()
	go func() {
		err := runner.Run(ctx)
		testutil.DiffError(t, nil, nil, err)
		close(waitStop)
	}()
	time.Sleep(time.Second) // Wait a little until the server starts.

	h1t := &http.Transport{
		ForceAttemptHTTP2: false,
	}

	// HTTP1 request.
	r1 := httptest.NewRequest(http.MethodGet, "http://127.0.0.1:12345/test", nil)
	w1, err := h1t.RoundTrip(r1)
	testutil.DiffError(t, nil, nil, err)
	b1, _ := io.ReadAll(w1.Body)
	testutil.Diff(t, http.StatusNotFound, w1.StatusCode)
	testutil.Diff(t, `{"status":404,"statusText":"Not Found"}`, string(b1))

}

func TestHTTPNetpol(t *testing.T) {

	configs := []string{
		testDataDir + "config-http-netpol.yaml",
	}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "core/v1",
		Kind:       "HTTPServer",
	}
	runner, err := api.ReferTypedObject[core.Runner](server, ref)
	testutil.DiffError(t, nil, nil, err)

	ctx, cancel := context.WithCancel(context.Background())
	waitStop := make(chan struct{})
	defer func() {
		cancel()
		<-waitStop
	}()
	go func() {
		err := runner.Run(ctx)
		testutil.DiffError(t, nil, nil, err)
		close(waitStop)
	}()
	time.Sleep(time.Second) // Wait a little until the server starts.

	h1t := &http.Transport{
		ForceAttemptHTTP2: false,
	}

	// HTTP1 request.
	// "127.0.0.1" is not allowed by security policy.
	r1 := httptest.NewRequest(http.MethodGet, "http://127.0.0.1:12345/test", nil)
	w1, err := h1t.RoundTrip(r1)
	testutil.Diff(t, true, err != nil) // EOF error or ReadTCP error
	testutil.Diff(t, (*http.Response)(nil), w1)

}

func TestHTTPTLS(t *testing.T) {

	configs := []string{
		testDataDir + "config-http-tls.yaml",
	}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "core/v1",
		Kind:       "HTTPServer",
	}
	runner, err := api.ReferTypedObject[core.Runner](server, ref)
	testutil.DiffError(t, nil, nil, err)

	ctx, cancel := context.WithCancel(context.Background())
	waitStop := make(chan struct{})
	defer func() {
		cancel()
		<-waitStop
	}()
	go func() {
		err := runner.Run(ctx)
		testutil.DiffError(t, nil, nil, err)
		close(waitStop)
	}()
	time.Sleep(time.Second) // Wait a little until the server starts.

	pem, _ := os.ReadFile(testDataDir + "testdata/cert.pem")
	pool := x509.NewCertPool()
	pool.AppendCertsFromPEM(pem)
	h1t := &http.Transport{
		ForceAttemptHTTP2: false,
		TLSClientConfig: &tls.Config{
			RootCAs: pool,
		},
	}

	// HTTP1 request with TLS.
	r1, _ := http.NewRequest(http.MethodGet, "https://127.0.0.1:12345/test", nil)
	w1, err := h1t.RoundTrip(r1)
	testutil.DiffError(t, nil, nil, err)
	b1, _ := io.ReadAll(w1.Body)
	testutil.Diff(t, http.StatusNotFound, w1.StatusCode)
	testutil.Diff(t, `{"status":404,"statusText":"Not Found"}`, string(b1))

}

func TestHTTPHandler(t *testing.T) {

	configs := []string{
		testDataDir + "config-http-handler.yaml",
	}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	testMiddlewareRef1 := &kernel.Reference{
		APIVersion: "container/v1",
		Kind:       "TestMiddleware",
		Name:       "test1",
		Namespace:  "",
	}
	m1 := &testMiddleware{}
	common.PostTestResource(server, testMiddlewareRef1, m1)
	testMiddlewareRef2 := &kernel.Reference{
		APIVersion: "container/v1",
		Kind:       "TestMiddleware",
		Name:       "test2",
		Namespace:  "",
	}
	m2 := &testMiddleware{}
	common.PostTestResource(server, testMiddlewareRef2, m2)
	testHandlerRef := &kernel.Reference{
		APIVersion: "container/v1",
		Kind:       "TestHandler",
		Name:       "test",
		Namespace:  "",
	}
	h := &testHandler{}
	common.PostTestResource(server, testHandlerRef, h)

	ref := &kernel.Reference{
		APIVersion: "core/v1",
		Kind:       "HTTPServer",
	}
	runner, err := api.ReferTypedObject[core.Runner](server, ref)
	testutil.DiffError(t, nil, nil, err)

	ctx, cancel := context.WithCancel(context.Background())
	waitStop := make(chan struct{})
	defer func() {
		cancel()
		<-waitStop
	}()
	go func() {
		err := runner.Run(ctx)
		testutil.DiffError(t, nil, nil, err)
		close(waitStop)
	}()
	time.Sleep(time.Second) // Wait a little until the server starts.

	h1t := &http.Transport{
		ForceAttemptHTTP2: false,
	}

	// HTTP1 request.
	r1, _ := http.NewRequest(http.MethodGet, "http://127.0.0.1:12345/prefix", nil)
	w1, err := h1t.RoundTrip(r1)
	testutil.DiffError(t, nil, nil, err)
	b1, _ := io.ReadAll(w1.Body)
	testutil.Diff(t, http.StatusOK, w1.StatusCode)
	testutil.Diff(t, "ok", w1.Header.Get("test"))
	testutil.Diff(t, "test", string(b1))
	testutil.Diff(t, 1, m1.called)
	testutil.Diff(t, 1, m2.called)
	testutil.Diff(t, 1, h.called)

	// "/test" does not match to the pattern.
	r2, _ := http.NewRequest(http.MethodGet, "http://127.0.0.1:12345/test", nil)
	w2, err := h1t.RoundTrip(r2)
	testutil.DiffError(t, nil, nil, err)
	b2, _ := io.ReadAll(w2.Body)
	testutil.Diff(t, http.StatusNotFound, w2.StatusCode)
	testutil.Diff(t, `{"status":404,"statusText":"Not Found"}`, string(b2))

	// POST is not allowed.
	r3, _ := http.NewRequest(http.MethodPost, "http://127.0.0.1:12345/prefix", nil)
	w3, err := h1t.RoundTrip(r3)
	testutil.DiffError(t, nil, nil, err)
	b3, _ := io.ReadAll(w3.Body)
	testutil.Diff(t, http.StatusNotFound, w3.StatusCode)
	testutil.Diff(t, `{"status":404,"statusText":"Not Found"}`, string(b3))

	// "localhost" is not allowed.
	r4, _ := http.NewRequest(http.MethodGet, "http://localhost:12345/prefix", nil)
	w4, err := h1t.RoundTrip(r4)
	testutil.DiffError(t, nil, nil, err)
	b4, _ := io.ReadAll(w4.Body)
	testutil.Diff(t, http.StatusNotFound, w4.StatusCode)
	testutil.Diff(t, `{"status":404,"statusText":"Not Found"}`, string(b4))

}

func TestHTTP2(t *testing.T) {

	configs := []string{
		testDataDir + "config-http2.yaml",
	}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "core/v1",
		Kind:       "HTTPServer",
	}
	runner, err := api.ReferTypedObject[core.Runner](server, ref)
	testutil.DiffError(t, nil, nil, err)

	ctx, cancel := context.WithCancel(context.Background())
	waitStop := make(chan struct{})
	defer func() {
		cancel()
		<-waitStop
	}()
	go func() {
		err := runner.Run(ctx)
		testutil.DiffError(t, nil, nil, err)
		close(waitStop)
	}()
	time.Sleep(time.Second) // Wait a little until the server starts.

	h1t := &http.Transport{
		ForceAttemptHTTP2: false,
	}
	h2t := &http2.Transport{
		AllowHTTP: true,
		DialTLSContext: func(_ context.Context, network, addr string, _ *tls.Config) (net.Conn, error) {
			return net.Dial(network, addr)
		},
	}

	// HTTP2 (h2c) request.
	// This will fails because the H2C is not enabled by the server.
	r1 := httptest.NewRequest(http.MethodGet, "http://127.0.0.1:12345/test", nil)
	w1, err := h2t.RoundTrip(r1)
	testutil.Diff(t, true, err != nil) // Connection closed error or Read TCP error
	testutil.Diff(t, (*http.Response)(nil), w1)

	// HTTP1 request.
	// HTTP1 requests are also succeeds.
	r2 := httptest.NewRequest(http.MethodGet, "http://127.0.0.1:12345/test", nil)
	w2, err := h1t.RoundTrip(r2)
	testutil.DiffError(t, nil, nil, err)
	b2, _ := io.ReadAll(w2.Body)
	testutil.Diff(t, http.StatusNotFound, w2.StatusCode)
	testutil.Diff(t, `{"status":404,"statusText":"Not Found"}`, string(b2))

}

func TestHTTP2H2C(t *testing.T) {

	configs := []string{
		testDataDir + "config-http2-h2c.yaml",
	}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "core/v1",
		Kind:       "HTTPServer",
	}
	runner, err := api.ReferTypedObject[core.Runner](server, ref)
	testutil.DiffError(t, nil, nil, err)

	ctx, cancel := context.WithCancel(context.Background())
	waitStop := make(chan struct{})
	defer func() {
		cancel()
		<-waitStop
	}()
	go func() {
		err := runner.Run(ctx)
		testutil.DiffError(t, nil, nil, err)
		close(waitStop)
	}()
	time.Sleep(time.Second) // Wait a little until the server starts.

	h1t := &http.Transport{
		ForceAttemptHTTP2: false,
	}
	h2t := &http2.Transport{
		AllowHTTP: true,
		DialTLSContext: func(_ context.Context, network, addr string, _ *tls.Config) (net.Conn, error) {
			return net.Dial(network, addr)
		},
	}

	// HTTP2 (h2c) request.
	r1 := httptest.NewRequest(http.MethodGet, "http://127.0.0.1:12345/test", nil)
	w1, err := h2t.RoundTrip(r1)
	testutil.DiffError(t, nil, nil, err)
	b1, _ := io.ReadAll(w1.Body)
	testutil.Diff(t, http.StatusNotFound, w1.StatusCode)
	testutil.Diff(t, `{"status":404,"statusText":"Not Found"}`, string(b1))

	// HTTP1 request.
	// HTTP1 requests are also succeeds.
	r2 := httptest.NewRequest(http.MethodGet, "http://127.0.0.1:12345/test", nil)
	w2, err := h1t.RoundTrip(r2)
	testutil.DiffError(t, nil, nil, err)
	b2, _ := io.ReadAll(w2.Body)
	testutil.Diff(t, http.StatusNotFound, w2.StatusCode)
	testutil.Diff(t, `{"status":404,"statusText":"Not Found"}`, string(b2))

}

func TestHTTP2TLS(t *testing.T) {

	configs := []string{
		testDataDir + "config-http2-tls.yaml",
	}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "core/v1",
		Kind:       "HTTPServer",
	}
	runner, err := api.ReferTypedObject[core.Runner](server, ref)
	testutil.DiffError(t, nil, nil, err)

	ctx, cancel := context.WithCancel(context.Background())
	waitStop := make(chan struct{})
	defer func() {
		cancel()
		<-waitStop
	}()
	go func() {
		err := runner.Run(ctx)
		testutil.DiffError(t, nil, nil, err)
		close(waitStop)
	}()
	time.Sleep(time.Second) // Wait a little until the server starts.

	pem, _ := os.ReadFile(testDataDir + "testdata/cert.pem")
	pool := x509.NewCertPool()
	pool.AppendCertsFromPEM(pem)
	h1t := &http.Transport{
		ForceAttemptHTTP2: false,
		TLSClientConfig: &tls.Config{
			RootCAs: pool,
		},
	}
	h2t := &http2.Transport{
		TLSClientConfig: &tls.Config{
			RootCAs: pool,
		},
	}

	// HTTP1 request with TLS.
	r1, _ := http.NewRequest(http.MethodGet, "https://127.0.0.1:12345/test", nil)
	w1, err := h1t.RoundTrip(r1)
	testutil.DiffError(t, nil, nil, err)
	b1, _ := io.ReadAll(w1.Body)
	testutil.Diff(t, http.StatusNotFound, w1.StatusCode)
	testutil.Diff(t, `{"status":404,"statusText":"Not Found"}`, string(b1))

	// HTTP2 request with TLS.
	r2, _ := http.NewRequest(http.MethodGet, "https://127.0.0.1:12345/test", nil)
	w2, err := h2t.RoundTrip(r2)
	testutil.DiffError(t, nil, nil, err)
	b2, _ := io.ReadAll(w2.Body)
	testutil.Diff(t, http.StatusNotFound, w2.StatusCode)
	testutil.Diff(t, `{"status":404,"statusText":"Not Found"}`, string(b2))

}

func TestHTTP2Netpol(t *testing.T) {

	configs := []string{
		testDataDir + "config-http2-netpol.yaml",
	}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "core/v1",
		Kind:       "HTTPServer",
	}
	runner, err := api.ReferTypedObject[core.Runner](server, ref)
	testutil.DiffError(t, nil, nil, err)

	ctx, cancel := context.WithCancel(context.Background())
	waitStop := make(chan struct{})
	defer func() {
		cancel()
		<-waitStop
	}()
	go func() {
		err := runner.Run(ctx)
		testutil.DiffError(t, nil, nil, err)
		close(waitStop)
	}()
	time.Sleep(time.Second) // Wait a little until the server starts.

	pem, _ := os.ReadFile(testDataDir + "testdata/cert.pem")
	pool := x509.NewCertPool()
	pool.AppendCertsFromPEM(pem)
	h1t := &http.Transport{
		ForceAttemptHTTP2: false,
		TLSClientConfig: &tls.Config{
			RootCAs: pool,
		},
	}
	h2t := &http2.Transport{
		TLSClientConfig: &tls.Config{
			RootCAs: pool,
		},
	}

	// HTTP1 request with TLS.
	// "127.0.0.1" is not allowed by security policy.
	r1, _ := http.NewRequest(http.MethodGet, "https://127.0.0.1:12345/test", nil)
	w1, err := h1t.RoundTrip(r1)
	testutil.Diff(t, true, err != nil) // EOF error or ReadTCP error
	testutil.Diff(t, (*http.Response)(nil), w1)

	// HTTP2 request with TLS.
	// "127.0.0.1" is not allowed by security policy.
	r2, _ := http.NewRequest(http.MethodGet, "https://127.0.0.1:12345/test", nil)
	w2, err := h2t.RoundTrip(r2)
	testutil.Diff(t, true, err != nil) // EOF error or ReadTCP error
	testutil.Diff(t, (*http.Response)(nil), w2)

}

func TestHTTP2Handler(t *testing.T) {

	configs := []string{
		testDataDir + "config-http2-handler.yaml",
	}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	testMiddlewareRef1 := &kernel.Reference{
		APIVersion: "container/v1",
		Kind:       "TestMiddleware",
		Name:       "test1",
		Namespace:  "",
	}
	m1 := &testMiddleware{}
	common.PostTestResource(server, testMiddlewareRef1, m1)
	testMiddlewareRef2 := &kernel.Reference{
		APIVersion: "container/v1",
		Kind:       "TestMiddleware",
		Name:       "test2",
		Namespace:  "",
	}
	m2 := &testMiddleware{}
	common.PostTestResource(server, testMiddlewareRef2, m2)
	testHandlerRef := &kernel.Reference{
		APIVersion: "container/v1",
		Kind:       "TestHandler",
		Name:       "test",
		Namespace:  "",
	}
	h := &testHandler{}
	common.PostTestResource(server, testHandlerRef, h)

	ref := &kernel.Reference{
		APIVersion: "core/v1",
		Kind:       "HTTPServer",
	}
	runner, err := api.ReferTypedObject[core.Runner](server, ref)
	testutil.DiffError(t, nil, nil, err)

	ctx, cancel := context.WithCancel(context.Background())
	waitStop := make(chan struct{})
	defer func() {
		cancel()
		<-waitStop
	}()
	go func() {
		err := runner.Run(ctx)
		testutil.DiffError(t, nil, nil, err)
		close(waitStop)
	}()
	time.Sleep(time.Second) // Wait a little until the server starts.

	pem, _ := os.ReadFile(testDataDir + "testdata/cert.pem")
	pool := x509.NewCertPool()
	pool.AppendCertsFromPEM(pem)
	h1t := &http.Transport{
		ForceAttemptHTTP2: false,
		TLSClientConfig: &tls.Config{
			RootCAs: pool,
		},
	}
	h2t := &http2.Transport{
		TLSClientConfig: &tls.Config{
			RootCAs: pool,
		},
	}

	// HTTP1 request with TLS.
	m1.called = 0
	m2.called = 0
	h.called = 0
	r1, _ := http.NewRequest(http.MethodGet, "https://127.0.0.1:12345/prefix", nil)
	w1, err := h1t.RoundTrip(r1)
	testutil.DiffError(t, nil, nil, err)
	b1, _ := io.ReadAll(w1.Body)
	testutil.Diff(t, http.StatusOK, w1.StatusCode)
	testutil.Diff(t, "ok", w1.Header.Get("test"))
	testutil.Diff(t, "test", string(b1))
	testutil.Diff(t, 1, m1.called)
	testutil.Diff(t, 1, m2.called)
	testutil.Diff(t, 1, h.called)

	// HTTP1 request with TLS.
	// "/test" does not match to the pattern.
	r2, _ := http.NewRequest(http.MethodGet, "https://127.0.0.1:12345/test", nil)
	w2, err := h1t.RoundTrip(r2)
	testutil.DiffError(t, nil, nil, err)
	b2, _ := io.ReadAll(w2.Body)
	testutil.Diff(t, http.StatusNotFound, w2.StatusCode)
	testutil.Diff(t, `{"status":404,"statusText":"Not Found"}`, string(b2))

	// HTTP1 request with TLS.
	// POST is not allowed.
	r3, _ := http.NewRequest(http.MethodPost, "https://127.0.0.1:12345/prefix", nil)
	w3, err := h1t.RoundTrip(r3)
	testutil.DiffError(t, nil, nil, err)
	b3, _ := io.ReadAll(w3.Body)
	testutil.Diff(t, http.StatusNotFound, w3.StatusCode)
	testutil.Diff(t, `{"status":404,"statusText":"Not Found"}`, string(b3))

	// HTTP1 request with TLS.
	// "localhost" is not allowed.
	r4, _ := http.NewRequest(http.MethodGet, "https://localhost:12345/prefix", nil)
	w4, err := h1t.RoundTrip(r4)
	testutil.DiffError(t, nil, nil, err)
	b4, _ := io.ReadAll(w4.Body)
	testutil.Diff(t, http.StatusNotFound, w4.StatusCode)
	testutil.Diff(t, `{"status":404,"statusText":"Not Found"}`, string(b4))

	// HTTP2 request with TLS.
	m1.called = 0
	m2.called = 0
	h.called = 0
	r5, _ := http.NewRequest(http.MethodGet, "https://127.0.0.1:12345/prefix", nil)
	w5, err := h2t.RoundTrip(r5)
	testutil.DiffError(t, nil, nil, err)
	b5, _ := io.ReadAll(w5.Body)
	testutil.Diff(t, http.StatusOK, w5.StatusCode)
	testutil.Diff(t, "ok", w1.Header.Get("test"))
	testutil.Diff(t, "test", string(b5))
	testutil.Diff(t, 1, m1.called)
	testutil.Diff(t, 1, m2.called)
	testutil.Diff(t, 1, h.called)

	// HTTP2 request with TLS.
	// "/test" does not match to the pattern.
	r6, _ := http.NewRequest(http.MethodGet, "https://127.0.0.1:12345/test", nil)
	w6, err := h2t.RoundTrip(r6)
	testutil.DiffError(t, nil, nil, err)
	b6, _ := io.ReadAll(w6.Body)
	testutil.Diff(t, http.StatusNotFound, w6.StatusCode)
	testutil.Diff(t, `{"status":404,"statusText":"Not Found"}`, string(b6))

	// HTTP2 request with TLS.
	// POST is not allowed.
	r7, _ := http.NewRequest(http.MethodPost, "https://127.0.0.1:12345/prefix", nil)
	w7, err := h2t.RoundTrip(r7)
	testutil.DiffError(t, nil, nil, err)
	b7, _ := io.ReadAll(w7.Body)
	testutil.Diff(t, http.StatusNotFound, w7.StatusCode)
	testutil.Diff(t, `{"status":404,"statusText":"Not Found"}`, string(b7))

	// HTTP2 request with TLS.
	// "localhost" is not allowed.
	r8, _ := http.NewRequest(http.MethodGet, "https://localhost:12345/prefix", nil)
	w8, err := h2t.RoundTrip(r8)
	testutil.DiffError(t, nil, nil, err)
	b8, _ := io.ReadAll(w8.Body)
	testutil.Diff(t, http.StatusNotFound, w8.StatusCode)
	testutil.Diff(t, `{"status":404,"statusText":"Not Found"}`, string(b8))

}

func TestHTTP3(t *testing.T) {

	configs := []string{
		testDataDir + "config-http3.yaml",
	}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "core/v1",
		Kind:       "HTTPServer",
	}
	runner, err := api.ReferTypedObject[core.Runner](server, ref)
	testutil.DiffError(t, nil, nil, err)

	ctx, cancel := context.WithCancel(context.Background())
	waitStop := make(chan struct{})
	defer func() {
		cancel()
		<-waitStop
	}()
	go func() {
		err := runner.Run(ctx)
		testutil.DiffError(t, nil, nil, err)
		close(waitStop)
	}()
	time.Sleep(time.Second) // Wait a little until the server starts.

	pem, _ := os.ReadFile(testDataDir + "testdata/cert.pem")
	pool := x509.NewCertPool()
	pool.AppendCertsFromPEM(pem)
	h3t := &http3.Transport{
		TLSClientConfig: &tls.Config{
			RootCAs: pool,
		},
	}

	// HTTP3 request.
	r1 := httptest.NewRequest(http.MethodGet, "https://127.0.0.1:12345/test", nil)
	w1, err := h3t.RoundTrip(r1)
	testutil.DiffError(t, nil, nil, err)
	b1, _ := io.ReadAll(w1.Body)
	testutil.Diff(t, http.StatusNotFound, w1.StatusCode)
	testutil.Diff(t, `{"status":404,"statusText":"Not Found"}`, string(b1))

}

func TestHTTP3Handler(t *testing.T) {

	configs := []string{
		testDataDir + "config-http3-handler.yaml",
	}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	testMiddlewareRef1 := &kernel.Reference{
		APIVersion: "container/v1",
		Kind:       "TestMiddleware",
		Name:       "test1",
		Namespace:  "",
	}
	m1 := &testMiddleware{}
	common.PostTestResource(server, testMiddlewareRef1, m1)
	testMiddlewareRef2 := &kernel.Reference{
		APIVersion: "container/v1",
		Kind:       "TestMiddleware",
		Name:       "test2",
		Namespace:  "",
	}
	m2 := &testMiddleware{}
	common.PostTestResource(server, testMiddlewareRef2, m2)
	testHandlerRef := &kernel.Reference{
		APIVersion: "container/v1",
		Kind:       "TestHandler",
		Name:       "test",
		Namespace:  "",
	}
	h := &testHandler{}
	common.PostTestResource(server, testHandlerRef, h)

	ref := &kernel.Reference{
		APIVersion: "core/v1",
		Kind:       "HTTPServer",
	}
	runner, err := api.ReferTypedObject[core.Runner](server, ref)
	testutil.DiffError(t, nil, nil, err)

	ctx, cancel := context.WithCancel(context.Background())
	waitStop := make(chan struct{})
	defer func() {
		cancel()
		<-waitStop
	}()
	go func() {
		err := runner.Run(ctx)
		testutil.DiffError(t, nil, nil, err)
		close(waitStop)
	}()
	time.Sleep(time.Second) // Wait a little until the server starts.

	pem, _ := os.ReadFile(testDataDir + "testdata/cert.pem")
	pool := x509.NewCertPool()
	pool.AppendCertsFromPEM(pem)
	h3t := &http3.Transport{
		TLSClientConfig: &tls.Config{
			RootCAs: pool,
		},
	}

	// HTTP3 request.
	r1, _ := http.NewRequest(http.MethodGet, "https://127.0.0.1:12345/prefix", nil)
	w1, err := h3t.RoundTrip(r1)
	testutil.DiffError(t, nil, nil, err)
	b1, _ := io.ReadAll(w1.Body)
	testutil.Diff(t, http.StatusOK, w1.StatusCode)
	testutil.Diff(t, "ok", w1.Header.Get("test"))
	testutil.Diff(t, "test", string(b1))
	testutil.Diff(t, 1, m1.called)
	testutil.Diff(t, 1, m2.called)
	testutil.Diff(t, 1, h.called)

	// "/test" does not match to the pattern.
	r2, _ := http.NewRequest(http.MethodGet, "https://127.0.0.1:12345/test", nil)
	w2, err := h3t.RoundTrip(r2)
	testutil.DiffError(t, nil, nil, err)
	b2, _ := io.ReadAll(w2.Body)
	testutil.Diff(t, http.StatusNotFound, w2.StatusCode)
	testutil.Diff(t, `{"status":404,"statusText":"Not Found"}`, string(b2))

	// POST is not allowed.
	r3, _ := http.NewRequest(http.MethodPost, "https://127.0.0.1:12345/prefix", nil)
	w3, err := h3t.RoundTrip(r3)
	testutil.DiffError(t, nil, nil, err)
	b3, _ := io.ReadAll(w3.Body)
	testutil.Diff(t, http.StatusNotFound, w3.StatusCode)
	testutil.Diff(t, `{"status":404,"statusText":"Not Found"}`, string(b3))

	// "localhost" is not allowed.
	r4, _ := http.NewRequest(http.MethodGet, "https://localhost:12345/prefix", nil)
	w4, err := h3t.RoundTrip(r4)
	testutil.DiffError(t, nil, nil, err)
	b4, _ := io.ReadAll(w4.Body)
	testutil.Diff(t, http.StatusNotFound, w4.StatusCode)
	testutil.Diff(t, `{"status":404,"statusText":"Not Found"}`, string(b4))

}
