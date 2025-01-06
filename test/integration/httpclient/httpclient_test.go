//go:build integration
// +build integration

package httpclient_test

import (
	"bytes"
	"crypto/tls"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/aileron-gateway/aileron-gateway/apis/kernel"
	"github.com/aileron-gateway/aileron-gateway/app"
	"github.com/aileron-gateway/aileron-gateway/core"
	"github.com/aileron-gateway/aileron-gateway/core/httpclient"
	"github.com/aileron-gateway/aileron-gateway/kernel/api"
	"github.com/aileron-gateway/aileron-gateway/kernel/er"
	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
	"github.com/aileron-gateway/aileron-gateway/test/integration/common"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/quic-go/quic-go/http3"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

type testTripperware struct {
	called int
}

func (t *testTripperware) Tripperware(next http.RoundTripper) http.RoundTripper {
	return core.RoundTripperFunc(func(r *http.Request) (*http.Response, error) {
		t.called += 1
		return next.RoundTrip(r)
	})
}

func TestRetry1(t *testing.T) {

	configs := []string{
		testDataDir + "config-retry-1.yaml",
	}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.Diff(t, nil, err)

	ref := &kernel.Reference{
		APIVersion: "core/v1",
		Kind:       "HTTPClient",
	}
	rt, err := api.ReferTypedObject[http.RoundTripper](server, ref)
	testutil.Diff(t, nil, err)

	var count int
	svr := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			count += 1
			conn, _, _ := w.(http.Hijacker).Hijack()
			conn.Close()
		},
	))
	defer svr.Close()

	count = 0 // Reset count.
	r1 := httptest.NewRequest(http.MethodGet, svr.URL+"/test", nil)
	w1, err := rt.RoundTrip(r1)
	retryErr := &er.Error{Package: httpclient.ErrPkg, Type: httpclient.ErrTypeRetry, Description: httpclient.ErrDescRetryFail}
	testutil.Diff(t, retryErr, err, cmpopts.EquateErrors())
	testutil.Diff(t, (*http.Response)(nil), w1)
	testutil.Diff(t, 2, count)

	// Requests larger than 10 bytes body won't be retried.
	count = 0 // Reset count.
	r2 := httptest.NewRequest(http.MethodPost, svr.URL+"/test", bytes.NewReader([]byte("12345678901")))
	w2, err := rt.RoundTrip(r2)
	testutil.DiffError(t, io.EOF, nil, err, cmpopts.EquateErrors())
	testutil.Diff(t, (*http.Response)(nil), w2)
	testutil.Diff(t, 1, count)

}

func TestRetry3(t *testing.T) {

	configs := []string{
		testDataDir + "config-retry-3.yaml",
	}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.Diff(t, nil, err)

	ref := &kernel.Reference{
		APIVersion: "core/v1",
		Kind:       "HTTPClient",
	}
	rt, err := api.ReferTypedObject[http.RoundTripper](server, ref)
	testutil.Diff(t, nil, err)

	var count int
	svr := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			count += 1
			conn, _, _ := w.(http.Hijacker).Hijack()
			conn.Close()
		},
	))
	defer svr.Close()

	count = 0 // Reset count.
	r1 := httptest.NewRequest(http.MethodGet, svr.URL+"/test", nil)
	w1, err := rt.RoundTrip(r1)
	retryErr := &er.Error{Package: httpclient.ErrPkg, Type: httpclient.ErrTypeRetry, Description: httpclient.ErrDescRetryFail}
	testutil.Diff(t, retryErr, err, cmpopts.EquateErrors())
	testutil.Diff(t, (*http.Response)(nil), w1)
	testutil.Diff(t, 4, count)

	// Requests larger than 10 bytes body won't be retried.
	count = 0 // Reset count.
	r2 := httptest.NewRequest(http.MethodPost, svr.URL+"/test", bytes.NewReader([]byte("12345678901")))
	w2, err := rt.RoundTrip(r2)
	testutil.DiffError(t, io.EOF, nil, err, cmpopts.EquateErrors())
	testutil.Diff(t, (*http.Response)(nil), w2)
	testutil.Diff(t, 1, count)

}

func TestRetryStatus(t *testing.T) {

	configs := []string{
		testDataDir + "config-retry-status.yaml",
	}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.Diff(t, nil, err)

	ref := &kernel.Reference{
		APIVersion: "core/v1",
		Kind:       "HTTPClient",
	}
	rt, err := api.ReferTypedObject[http.RoundTripper](server, ref)
	testutil.Diff(t, nil, err)

	var count int
	var status int
	svr := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			count += 1
			w.WriteHeader(status)
		},
	))
	defer svr.Close()

	retryErr := &er.Error{Package: httpclient.ErrPkg, Type: httpclient.ErrTypeRetry, Description: httpclient.ErrDescRetryFail}

	count = 0                               // Reset count.
	status = http.StatusInternalServerError // Use 500 error.
	r1 := httptest.NewRequest(http.MethodGet, svr.URL+"/test", nil)
	w1, err := rt.RoundTrip(r1)
	testutil.Diff(t, retryErr, err, cmpopts.EquateErrors())
	testutil.Diff(t, (*http.Response)(nil), w1)
	testutil.Diff(t, 3, count)

	count = 0                              // Reset count.
	status = http.StatusServiceUnavailable // Use 503 error.
	r2 := httptest.NewRequest(http.MethodGet, svr.URL+"/test", nil)
	w2, err := rt.RoundTrip(r2)
	testutil.Diff(t, retryErr, err, cmpopts.EquateErrors())
	testutil.Diff(t, (*http.Response)(nil), w2)
	testutil.Diff(t, 3, count)

	// Requests larger than 10 bytes body won't be retried.
	count = 0                               // Reset count.
	status = http.StatusInternalServerError // Use 500 error.
	r3 := httptest.NewRequest(http.MethodPost, svr.URL+"/test", bytes.NewReader([]byte("12345678901")))
	w3, err := rt.RoundTrip(r3)
	testutil.Diff(t, nil, err)
	testutil.Diff(t, http.StatusInternalServerError, w3.StatusCode)
	testutil.Diff(t, 1, count)

}

func TestHTTP(t *testing.T) {

	configs := []string{
		testDataDir + "config-http.yaml",
	}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.Diff(t, nil, err)

	testRef := &kernel.Reference{
		APIVersion: "container/v1",
		Kind:       "TestTripperware",
		Name:       "test",
		Namespace:  "",
	}
	testTripperware := &testTripperware{}
	common.PostTestResource(server, testRef, testTripperware)

	ref := &kernel.Reference{
		APIVersion: "core/v1",
		Kind:       "HTTPClient",
	}
	rt, err := api.ReferTypedObject[http.RoundTripper](server, ref)
	testutil.Diff(t, nil, err)

	var sleep int
	var extraHeader string
	svr := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(time.Millisecond * time.Duration(sleep))
			w.Header().Set("test", "ok")
			w.Header().Set("extra", extraHeader)
			w.WriteHeader(http.StatusFound)
			w.Write([]byte("test"))
		},
	))
	defer svr.Close()

	sleep = 0
	extraHeader = ""
	testTripperware.called = 0
	r1 := httptest.NewRequest(http.MethodGet, svr.URL+"/test", nil)
	w1, err := rt.RoundTrip(r1)
	testutil.Diff(t, nil, err)
	b1, _ := io.ReadAll(w1.Body)
	testutil.Diff(t, "ok", w1.Header.Get("test"))
	testutil.Diff(t, "test", string(b1))
	testutil.Diff(t, 1, testTripperware.called)

	// Too large response header.
	sleep = 0
	extraHeader = "1234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890"
	testTripperware.called = 0
	r2 := httptest.NewRequest(http.MethodGet, svr.URL+"/test", nil)
	w2, err := rt.RoundTrip(r2)
	testutil.Diff(t, true, strings.Contains(err.Error(), "server response headers exceeded 200 bytes"))
	testutil.Diff(t, (*http.Response)(nil), w2)
	testutil.Diff(t, 1, testTripperware.called)

	// Too late response header.
	sleep = 100
	extraHeader = ""
	r3 := httptest.NewRequest(http.MethodGet, svr.URL+"/test", nil)
	w3, err := rt.RoundTrip(r3)
	testutil.Diff(t, true, strings.Contains(err.Error(), "timeout awaiting response headers"))
	testutil.Diff(t, (*http.Response)(nil), w3)

}

func TestHTTPTLS(t *testing.T) {

	configs := []string{
		testDataDir + "config-http-tls.yaml",
	}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.Diff(t, nil, err)

	ref := &kernel.Reference{
		APIVersion: "core/v1",
		Kind:       "HTTPClient",
	}
	rt, err := api.ReferTypedObject[http.RoundTripper](server, ref)
	testutil.Diff(t, nil, err)

	var sleep int
	var extraHeader string
	svr := httptest.NewUnstartedServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(time.Millisecond * time.Duration(sleep))
			w.Header().Set("test", "ok")
			w.Header().Set("extra", extraHeader)
			w.WriteHeader(http.StatusFound)
			w.Write([]byte("test"))
		},
	))
	certs, _ := tls.LoadX509KeyPair(testDataDir+"testdata/cert.pem", testDataDir+"testdata/key.pem")
	svr.TLS = &tls.Config{
		Certificates: []tls.Certificate{certs},
	}
	svr.StartTLS()
	defer svr.Close()

	sleep = 0
	extraHeader = ""
	r1 := httptest.NewRequest(http.MethodGet, svr.URL+"/test", nil)
	w1, err := rt.RoundTrip(r1)
	testutil.Diff(t, nil, err)
	b1, _ := io.ReadAll(w1.Body)
	testutil.Diff(t, "ok", w1.Header.Get("test"))
	testutil.Diff(t, "test", string(b1))

	// Too large response header.
	sleep = 0
	extraHeader = "1234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890"
	r2 := httptest.NewRequest(http.MethodGet, svr.URL+"/test", nil)
	w2, err := rt.RoundTrip(r2)
	testutil.Diff(t, true, strings.Contains(err.Error(), "server response headers exceeded 200 bytes"))
	testutil.Diff(t, (*http.Response)(nil), w2)

	// Too late response header.
	sleep = 100
	extraHeader = ""
	r3 := httptest.NewRequest(http.MethodGet, svr.URL+"/test", nil)
	w3, err := rt.RoundTrip(r3)
	testutil.Diff(t, true, strings.Contains(err.Error(), "timeout awaiting response headers"))
	testutil.Diff(t, (*http.Response)(nil), w3)

}

func TestHTTPDialer(t *testing.T) {

	configs := []string{
		testDataDir + "config-http-dialer.yaml",
	}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.Diff(t, nil, err)

	ref := &kernel.Reference{
		APIVersion: "core/v1",
		Kind:       "HTTPClient",
		Name:       "default",
		Namespace:  "",
	}
	rt, err := api.ReferTypedObject[http.RoundTripper](server, ref)
	testutil.Diff(t, nil, err)

	var sleep int
	var extraHeader string
	svr := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(time.Millisecond * time.Duration(sleep))
			w.Header().Set("test", "ok")
			w.Header().Set("extra", extraHeader)
			w.WriteHeader(http.StatusFound)
			w.Write([]byte("test"))
		},
	))
	defer svr.Close()

	sleep = 0
	extraHeader = ""
	r1 := httptest.NewRequest(http.MethodGet, svr.URL+"/test", nil)
	w1, err := rt.RoundTrip(r1)
	testutil.Diff(t, nil, err)
	b1, _ := io.ReadAll(w1.Body)
	testutil.Diff(t, "ok", w1.Header.Get("test"))
	testutil.Diff(t, "test", string(b1))

	// Too large response header.
	sleep = 0
	extraHeader = "1234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890"
	r2 := httptest.NewRequest(http.MethodGet, svr.URL+"/test", nil)
	w2, err := rt.RoundTrip(r2)
	testutil.Diff(t, true, strings.Contains(err.Error(), "server response headers exceeded 200 bytes"))
	testutil.Diff(t, (*http.Response)(nil), w2)

	// Too late response header.
	sleep = 100
	extraHeader = ""
	r3 := httptest.NewRequest(http.MethodGet, svr.URL+"/test", nil)
	w3, err := rt.RoundTrip(r3)
	testutil.Diff(t, true, strings.Contains(err.Error(), "timeout awaiting response headers"))
	testutil.Diff(t, (*http.Response)(nil), w3)

}

func TestHTTPTLSDialer(t *testing.T) {

	configs := []string{
		testDataDir + "config-http-tls-dialer.yaml",
	}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.Diff(t, nil, err)

	ref := &kernel.Reference{
		APIVersion: "core/v1",
		Kind:       "HTTPClient",
	}
	rt, err := api.ReferTypedObject[http.RoundTripper](server, ref)
	testutil.Diff(t, nil, err)

	var sleep int
	var extraHeader string
	svr := httptest.NewUnstartedServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(time.Millisecond * time.Duration(sleep))
			w.Header().Set("test", "ok")
			w.Header().Set("extra", extraHeader)
			w.WriteHeader(http.StatusFound)
			w.Write([]byte("test"))
		},
	))
	certs, _ := tls.LoadX509KeyPair(testDataDir+"testdata/cert.pem", testDataDir+"testdata/key.pem")
	svr.TLS = &tls.Config{
		Certificates: []tls.Certificate{certs},
	}
	svr.StartTLS()
	defer svr.Close()

	sleep = 0
	extraHeader = ""
	r1 := httptest.NewRequest(http.MethodGet, svr.URL+"/test", nil)
	w1, err := rt.RoundTrip(r1)
	testutil.Diff(t, nil, err)
	b1, _ := io.ReadAll(w1.Body)
	testutil.Diff(t, "ok", w1.Header.Get("test"))
	testutil.Diff(t, "test", string(b1))

	// Too large response header.
	sleep = 0
	extraHeader = "1234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890"
	r2 := httptest.NewRequest(http.MethodGet, svr.URL+"/test", nil)
	w2, err := rt.RoundTrip(r2)
	testutil.Diff(t, true, strings.Contains(err.Error(), "server response headers exceeded 200 bytes"))
	testutil.Diff(t, (*http.Response)(nil), w2)

	// Too late response header.
	sleep = 100
	extraHeader = ""
	r3 := httptest.NewRequest(http.MethodGet, svr.URL+"/test", nil)
	w3, err := rt.RoundTrip(r3)
	testutil.Diff(t, true, strings.Contains(err.Error(), "timeout awaiting response headers"))
	testutil.Diff(t, (*http.Response)(nil), w3)

}

func TestHTTP2(t *testing.T) {

	configs := []string{
		testDataDir + "config-http2.yaml",
	}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.Diff(t, nil, err)

	testRef := &kernel.Reference{
		APIVersion: "container/v1",
		Kind:       "TestTripperware",
		Name:       "test",
		Namespace:  "",
	}
	testTripperware := &testTripperware{}
	common.PostTestResource(server, testRef, testTripperware)

	ref := &kernel.Reference{
		APIVersion: "core/v1",
		Kind:       "HTTPClient",
	}
	rt, err := api.ReferTypedObject[http.RoundTripper](server, ref)
	testutil.Diff(t, nil, err)

	available, _ := net.Listen("tcp", "127.0.0.1:0")
	available.Close()
	testAddr := available.Addr().String()

	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("test", "ok")
		w.WriteHeader(http.StatusFound)
		w.Write([]byte("test"))
	})
	svr := &http.Server{
		Addr:    testAddr,
		Handler: h2c.NewHandler(h, &http2.Server{}),
	}
	go func() { svr.ListenAndServe() }()
	time.Sleep(time.Second) // Wait a little until server starts.
	defer svr.Close()

	testTripperware.called = 0
	r1, _ := http.NewRequest(http.MethodGet, "http://"+testAddr+"/test", nil)
	w1, err := rt.RoundTrip(r1)
	testutil.Diff(t, nil, err)
	b1, _ := io.ReadAll(w1.Body)
	testutil.Diff(t, "ok", w1.Header.Get("test"))
	testutil.Diff(t, "test", string(b1))
	testutil.Diff(t, 1, testTripperware.called)

}

func TestHTTP2MultiIP(t *testing.T) {

	configs := []string{
		testDataDir + "config-http2-multiIP.yaml",
	}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.Diff(t, nil, err)

	testRef := &kernel.Reference{
		APIVersion: "container/v1",
		Kind:       "TestTripperware",
		Name:       "test",
		Namespace:  "",
	}
	testTripperware := &testTripperware{}
	common.PostTestResource(server, testRef, testTripperware)

	ref := &kernel.Reference{
		APIVersion: "core/v1",
		Kind:       "HTTPClient",
	}
	rt, err := api.ReferTypedObject[http.RoundTripper](server, ref)
	testutil.Diff(t, nil, err)

	available, _ := net.Listen("tcp", "127.0.0.1:0")
	available.Close()
	testAddr := available.Addr().String()

	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("test", "ok")
		w.WriteHeader(http.StatusFound)
		w.Write([]byte("test"))
	})
	svr := &http.Server{
		Addr:    testAddr,
		Handler: h2c.NewHandler(h, &http2.Server{}),
	}
	go func() { svr.ListenAndServe() }()
	time.Sleep(time.Second) // Wait a little until server starts.
	defer svr.Close()

	testTripperware.called = 0
	r1, _ := http.NewRequest(http.MethodGet, "http://"+testAddr+"/test", nil)
	w1, err := rt.RoundTrip(r1)
	testutil.Diff(t, nil, err)
	b1, _ := io.ReadAll(w1.Body)
	testutil.Diff(t, "ok", w1.Header.Get("test"))
	testutil.Diff(t, "test", string(b1))
	testutil.Diff(t, 1, testTripperware.called)

}

func TestHTTP2TLS(t *testing.T) {

	configs := []string{
		testDataDir + "config-http2-tls.yaml",
	}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.Diff(t, nil, err)

	ref := &kernel.Reference{
		APIVersion: "core/v1",
		Kind:       "HTTPClient",
	}
	rt, err := api.ReferTypedObject[http.RoundTripper](server, ref)
	testutil.Diff(t, nil, err)

	svr := httptest.NewUnstartedServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("test", "ok")
			w.WriteHeader(http.StatusFound)
			w.Write([]byte("test"))
		},
	))
	certs, _ := tls.LoadX509KeyPair(testDataDir+"testdata/cert.pem", testDataDir+"testdata/key.pem")
	svr.TLS = &tls.Config{
		Certificates: []tls.Certificate{certs},
	}
	svr.EnableHTTP2 = true
	svr.StartTLS()
	defer svr.Close()

	r1, _ := http.NewRequest(http.MethodGet, svr.URL+"/test", nil)
	w1, err := rt.RoundTrip(r1)
	testutil.Diff(t, nil, err)
	b1, _ := io.ReadAll(w1.Body)
	testutil.Diff(t, "ok", w1.Header.Get("test"))
	testutil.Diff(t, "test", string(b1))

}

func TestHTTP2TLSDialer(t *testing.T) {

	configs := []string{
		testDataDir + "config-http2-tls-dialer.yaml",
	}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.Diff(t, nil, err)

	ref := &kernel.Reference{
		APIVersion: "core/v1",
		Kind:       "HTTPClient",
	}
	rt, err := api.ReferTypedObject[http.RoundTripper](server, ref)
	testutil.Diff(t, nil, err)

	svr := httptest.NewUnstartedServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("test", "ok")
			w.WriteHeader(http.StatusFound)
			w.Write([]byte("test"))
		},
	))
	certs, _ := tls.LoadX509KeyPair(testDataDir+"testdata/cert.pem", testDataDir+"testdata/key.pem")
	svr.TLS = &tls.Config{
		Certificates: []tls.Certificate{certs},
	}
	svr.EnableHTTP2 = true
	svr.StartTLS()
	defer svr.Close()

	r1, _ := http.NewRequest(http.MethodGet, svr.URL+"/test", nil)
	w1, err := rt.RoundTrip(r1)
	testutil.Diff(t, nil, err)
	b1, _ := io.ReadAll(w1.Body)
	testutil.Diff(t, "ok", w1.Header.Get("test"))
	testutil.Diff(t, "test", string(b1))

}

func TestHTTP2Dialer(t *testing.T) {

	configs := []string{
		testDataDir + "config-http2-dialer.yaml",
	}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.Diff(t, nil, err)

	ref := &kernel.Reference{
		APIVersion: "core/v1",
		Kind:       "HTTPClient",
	}
	rt, err := api.ReferTypedObject[http.RoundTripper](server, ref)
	testutil.Diff(t, nil, err)

	available, _ := net.Listen("tcp", "127.0.0.1:0")
	available.Close()
	testAddr := available.Addr().String()

	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("test", "ok")
		w.WriteHeader(http.StatusFound)
		w.Write([]byte("test"))
	})
	svr := &http.Server{
		Addr:    testAddr,
		Handler: h2c.NewHandler(h, &http2.Server{}),
	}
	go func() { svr.ListenAndServe() }()
	time.Sleep(time.Second) // Wait a little until server starts.
	defer svr.Close()

	r1, _ := http.NewRequest(http.MethodGet, "http://"+testAddr+"/test", nil)
	w1, err := rt.RoundTrip(r1)
	testutil.Diff(t, nil, err)
	b1, _ := io.ReadAll(w1.Body)
	testutil.Diff(t, "ok", w1.Header.Get("test"))
	testutil.Diff(t, "test", string(b1))

}

func TestHTTP3TLS(t *testing.T) {

	configs := []string{
		testDataDir + "config-http3-tls.yaml",
	}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.Diff(t, nil, err)

	testRef := &kernel.Reference{
		APIVersion: "container/v1",
		Kind:       "TestTripperware",
		Name:       "test",
		Namespace:  "",
	}
	testTripperware := &testTripperware{}
	common.PostTestResource(server, testRef, testTripperware)

	ref := &kernel.Reference{
		APIVersion: "core/v1",
		Kind:       "HTTPClient",
	}
	rt, err := api.ReferTypedObject[http.RoundTripper](server, ref)
	testutil.Diff(t, nil, err)

	available, _ := net.Listen("tcp", "127.0.0.1:0")
	available.Close()
	testAddr := available.Addr().String()

	var extraHeader string
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("test", "ok")
		w.Header().Set("extra", extraHeader)
		w.WriteHeader(http.StatusFound)
		w.Write([]byte("test"))
	})
	certs, _ := tls.LoadX509KeyPair(testDataDir+"testdata/cert.pem", testDataDir+"testdata/key.pem")
	svr := &http3.Server{
		Addr:    testAddr,
		Handler: h,
		TLSConfig: &tls.Config{
			Certificates: []tls.Certificate{certs},
		},
	}
	go func() { svr.ListenAndServe() }()
	time.Sleep(time.Second) // Wait a little until server starts.
	defer svr.Close()

	extraHeader = ""
	testTripperware.called = 0
	r1, _ := http.NewRequest(http.MethodGet, "https://"+testAddr+"/test", nil)
	w1, err := rt.RoundTrip(r1)
	testutil.Diff(t, nil, err)
	b1, _ := io.ReadAll(w1.Body)
	testutil.Diff(t, "ok", w1.Header.Get("test"))
	testutil.Diff(t, "test", string(b1))
	testutil.Diff(t, 1, testTripperware.called)

	// Requests with not allowed "[::1]".
	extraHeader = "1234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890"
	testTripperware.called = 0
	r2, _ := http.NewRequest(http.MethodGet, "https://"+testAddr+"/test", nil)
	w2, err := rt.RoundTrip(r2)
	testutil.Diff(t, true, strings.Contains(err.Error(), "HEADERS frame too large"))
	testutil.Diff(t, (*http.Response)(nil), w2)

}
