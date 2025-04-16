// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

//go:build integration

package httpproxy_test

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/aileron-gateway/aileron-gateway/apis/kernel"
	"github.com/aileron-gateway/aileron-gateway/cmd/aileron/app"
	"github.com/aileron-gateway/aileron-gateway/core"
	"github.com/aileron-gateway/aileron-gateway/kernel/api"
	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
	"github.com/aileron-gateway/aileron-gateway/test/integration/common"
	"github.com/google/go-cmp/cmp/cmpopts"
)

type testRoundTripper struct {
	called int
	body   []byte
}

func (rt *testRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	rt.called += 1
	return &http.Response{
		Header: http.Header{
			"Test": []string{"ok"},
		},
		StatusCode: http.StatusFound,
		Body:       io.NopCloser(bytes.NewReader(rt.body)),
	}, nil
}

type testTripperware struct {
	called int
}

func (t *testTripperware) Tripperware(next http.RoundTripper) http.RoundTripper {
	return core.RoundTripperFunc(func(r *http.Request) (*http.Response, error) {
		t.called += 1
		return next.RoundTrip(r)
	})
}

func Test0LoadBalancer(t *testing.T) {

	configs := []string{
		testDataDir + "config-0-loadbalancer.yaml",
	}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "core/v1",
		Kind:       "ReverseProxyHandler",
		Name:       "default",
		Namespace:  "",
	}
	handler, err := api.ReferTypedObject[http.Handler](server, ref)
	testutil.DiffError(t, nil, nil, err)

	r := httptest.NewRequest(http.MethodGet, "http://test.com/test", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)
	testutil.DiffError(t, nil, nil, err)
	b, _ := io.ReadAll(w.Body)
	testutil.Diff(t, http.StatusNotFound, w.Result().StatusCode)
	testutil.Diff(t, `{"status":404,"statusText":"Not Found"}`, string(b))

}

func Test0Upstream(t *testing.T) {

	configs := []string{
		testDataDir + "config-0-upstream.yaml",
	}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "core/v1",
		Kind:       "ReverseProxyHandler",
		Name:       "default",
		Namespace:  "",
	}
	handler, err := api.ReferTypedObject[http.Handler](server, ref)
	testutil.DiffError(t, nil, nil, err)

	svr := &http.Server{
		Addr: ":10001",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("test", "ok")
			w.WriteHeader(http.StatusFound)
			w.Write([]byte("test"))
		}),
	}
	go func() { svr.ListenAndServe() }()
	time.Sleep(time.Second) // Wait a little until server starts.
	defer svr.Close()

	r1 := httptest.NewRequest(http.MethodGet, "http://test.com/test", nil)
	w1 := httptest.NewRecorder()
	handler.ServeHTTP(w1, r1)
	b1, _ := io.ReadAll(w1.Body)
	testutil.Diff(t, http.StatusBadGateway, w1.Result().StatusCode)
	testutil.Diff(t, "", w1.Result().Header.Get("test"))
	testutil.Diff(t, `{"status":502,"statusText":"Bad Gateway"}`, string(b1))

	r2 := httptest.NewRequest(http.MethodGet, "http://test.com/none", nil)
	w2 := httptest.NewRecorder()
	handler.ServeHTTP(w2, r2)
	b2, _ := io.ReadAll(w2.Body)
	testutil.Diff(t, http.StatusNotFound, w2.Result().StatusCode)
	testutil.Diff(t, "", w2.Result().Header.Get("test"))
	testutil.Diff(t, `{"status":404,"statusText":"Not Found"}`, string(b2))

}

func Test1Upstream(t *testing.T) {

	configs := []string{
		testDataDir + "config-1-upstream.yaml",
	}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "core/v1",
		Kind:       "ReverseProxyHandler",
		Name:       "default",
		Namespace:  "",
	}
	handler, err := api.ReferTypedObject[http.Handler](server, ref)
	testutil.DiffError(t, nil, nil, err)

	svr := &http.Server{
		Addr: ":10001",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("test", "ok")
			w.WriteHeader(http.StatusFound)
			w.Write([]byte("test"))
		}),
	}
	go func() { svr.ListenAndServe() }()
	time.Sleep(time.Second) // Wait a little until server starts.
	defer svr.Close()

	r1 := httptest.NewRequest(http.MethodGet, "http://test.com/test", nil)
	w1 := httptest.NewRecorder()
	handler.ServeHTTP(w1, r1)
	b1, _ := io.ReadAll(w1.Body)
	testutil.Diff(t, http.StatusFound, w1.Result().StatusCode)
	testutil.Diff(t, "ok", w1.Result().Header.Get("test"))
	testutil.Diff(t, "test", string(b1))

	r2 := httptest.NewRequest(http.MethodGet, "http://test.com/none", nil)
	w2 := httptest.NewRecorder()
	handler.ServeHTTP(w2, r2)
	b2, _ := io.ReadAll(w2.Body)
	testutil.Diff(t, http.StatusNotFound, w2.Result().StatusCode)
	testutil.Diff(t, "", w2.Result().Header.Get("test"))
	testutil.Diff(t, `{"status":404,"statusText":"Not Found"}`, string(b2))

}

func Test2Upstream(t *testing.T) {

	configs := []string{
		testDataDir + "config-2-upstream.yaml",
	}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "core/v1",
		Kind:       "ReverseProxyHandler",
		Name:       "default",
		Namespace:  "",
	}
	handler, err := api.ReferTypedObject[http.Handler](server, ref)
	testutil.DiffError(t, nil, nil, err)

	svr1 := &http.Server{
		Addr: ":10001",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("test", "h1")
			w.WriteHeader(http.StatusFound)
			w.Write([]byte(r.URL.Path))
		}),
	}
	svr2 := &http.Server{
		Addr: ":10002",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("test", "h2")
			w.WriteHeader(http.StatusFound)
			w.Write([]byte(r.URL.Path))
		}),
	}
	go func() { svr1.ListenAndServe() }()
	go func() { svr2.ListenAndServe() }()
	time.Sleep(time.Second) // Wait a little until server starts.
	defer svr1.Close()
	defer svr2.Close()

	r1 := httptest.NewRequest(http.MethodGet, "http://test.com/test/foo", nil)
	w1 := httptest.NewRecorder()
	handler.ServeHTTP(w1, r1)
	b1, _ := io.ReadAll(w1.Body)
	testutil.Diff(t, http.StatusFound, w1.Result().StatusCode)
	testutil.Diff(t, "h1", w1.Result().Header.Get("test"))
	testutil.Diff(t, "/test/foo", string(b1))

	r2 := httptest.NewRequest(http.MethodGet, "http://test.com/test/bar", nil)
	w2 := httptest.NewRecorder()
	handler.ServeHTTP(w2, r2)
	b2, _ := io.ReadAll(w2.Body)
	testutil.Diff(t, http.StatusFound, w2.Result().StatusCode)
	testutil.Diff(t, "h2", w2.Result().Header.Get("test"))
	testutil.Diff(t, "/test/bar", string(b2))

	r3 := httptest.NewRequest(http.MethodGet, "http://test.com/none", nil)
	w3 := httptest.NewRecorder()
	handler.ServeHTTP(w3, r3)
	b3, _ := io.ReadAll(w3.Body)
	testutil.Diff(t, http.StatusNotFound, w3.Result().StatusCode)
	testutil.Diff(t, `{"status":404,"statusText":"Not Found"}`, string(b3))

}

func Test2LoadBalancer(t *testing.T) {

	configs := []string{
		testDataDir + "config-2-upstream.yaml",
	}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "core/v1",
		Kind:       "ReverseProxyHandler",
		Name:       "default",
		Namespace:  "",
	}
	handler, err := api.ReferTypedObject[http.Handler](server, ref)
	testutil.DiffError(t, nil, nil, err)

	svr1 := &http.Server{
		Addr: ":10001",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("test", "h1")
			w.WriteHeader(http.StatusFound)
			w.Write([]byte(r.URL.Path))
		}),
	}
	svr2 := &http.Server{
		Addr: ":10002",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("test", "h2")
			w.WriteHeader(http.StatusFound)
			w.Write([]byte(r.URL.Path))
		}),
	}
	go func() { svr1.ListenAndServe() }()
	go func() { svr2.ListenAndServe() }()
	time.Sleep(time.Second) // Wait a little until server starts.
	defer svr1.Close()
	defer svr2.Close()

	r1 := httptest.NewRequest(http.MethodGet, "http://test.com/test1/foo", nil)
	w1 := httptest.NewRecorder()
	handler.ServeHTTP(w1, r1)
	b1, _ := io.ReadAll(w1.Body)
	testutil.Diff(t, http.StatusFound, w1.Result().StatusCode)
	testutil.Diff(t, "h1", w1.Result().Header.Get("test"))
	testutil.Diff(t, "/test1/foo", string(b1))

	r2 := httptest.NewRequest(http.MethodGet, "http://test.com/test2/bar", nil)
	w2 := httptest.NewRecorder()
	handler.ServeHTTP(w2, r2)
	b2, _ := io.ReadAll(w2.Body)
	testutil.Diff(t, http.StatusFound, w2.Result().StatusCode)
	testutil.Diff(t, "h2", w2.Result().Header.Get("test"))
	testutil.Diff(t, "/test2/bar", string(b2))

	r3 := httptest.NewRequest(http.MethodGet, "http://test.com/none", nil)
	w3 := httptest.NewRecorder()
	handler.ServeHTTP(w3, r3)
	b3, _ := io.ReadAll(w3.Body)
	testutil.Diff(t, http.StatusNotFound, w3.Result().StatusCode)
	testutil.Diff(t, `{"status":404,"statusText":"Not Found"}`, string(b3))

}

func TestMatchTypes(t *testing.T) {

	configs := []string{
		testDataDir + "config-match-types.yaml",
	}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "core/v1",
		Kind:       "ReverseProxyHandler",
		Name:       "default",
		Namespace:  "",
	}
	handler, err := api.ReferTypedObject[http.Handler](server, ref)
	testutil.DiffError(t, nil, nil, err)

	svr1 := &http.Server{
		Addr: ":10001",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("test", "h1")
			w.WriteHeader(http.StatusFound)
		}),
	}
	svr2 := &http.Server{
		Addr: ":10002",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("test", "h2")
			w.WriteHeader(http.StatusFound)
		}),
	}
	svr3 := &http.Server{
		Addr: ":10003",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("test", "h3")
			w.WriteHeader(http.StatusFound)
		}),
	}
	svr4 := &http.Server{
		Addr: ":10004",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("test", "h4")
			w.WriteHeader(http.StatusFound)
		}),
	}
	svr5 := &http.Server{
		Addr: ":10005",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("test", "h5")
			w.WriteHeader(http.StatusFound)
		}),
	}
	svr6 := &http.Server{
		Addr: ":10006",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("test", "h6")
			w.WriteHeader(http.StatusFound)
		}),
	}
	svr7 := &http.Server{
		Addr: ":10007",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("test", "h7")
			w.WriteHeader(http.StatusFound)
		}),
	}
	svr8 := &http.Server{
		Addr: ":10008",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("test", "h8")
			w.WriteHeader(http.StatusFound)
		}),
	}
	go func() { svr1.ListenAndServe() }()
	go func() { svr2.ListenAndServe() }()
	go func() { svr3.ListenAndServe() }()
	go func() { svr4.ListenAndServe() }()
	go func() { svr5.ListenAndServe() }()
	go func() { svr6.ListenAndServe() }()
	go func() { svr7.ListenAndServe() }()
	go func() { svr8.ListenAndServe() }()
	time.Sleep(time.Second) // Wait a little until server starts.
	defer svr1.Close()
	defer svr2.Close()
	defer svr3.Close()
	defer svr4.Close()
	defer svr5.Close()
	defer svr6.Close()
	defer svr7.Close()
	defer svr8.Close()

	// Match to Exact
	r1 := httptest.NewRequest(http.MethodGet, "http://test.com/test1", nil)
	w1 := httptest.NewRecorder()
	handler.ServeHTTP(w1, r1)
	testutil.Diff(t, http.StatusFound, w1.Result().StatusCode)
	testutil.Diff(t, "h1", w1.Result().Header.Get("test"))

	// Match to Prefix
	r2 := httptest.NewRequest(http.MethodGet, "http://test.com/test2/foo/bar", nil)
	w2 := httptest.NewRecorder()
	handler.ServeHTTP(w2, r2)
	testutil.Diff(t, http.StatusFound, w2.Result().StatusCode)
	testutil.Diff(t, "h2", w2.Result().Header.Get("test"))

	// Match to Suffix
	r3 := httptest.NewRequest(http.MethodGet, "http://test.com/foo/bar/test3", nil)
	w3 := httptest.NewRecorder()
	handler.ServeHTTP(w3, r3)
	testutil.Diff(t, http.StatusFound, w3.Result().StatusCode)
	testutil.Diff(t, "h3", w3.Result().Header.Get("test"))

	// Match to Contain
	r4 := httptest.NewRequest(http.MethodGet, "http://test.com/foo/test4/bar", nil)
	w4 := httptest.NewRecorder()
	handler.ServeHTTP(w4, r4)
	testutil.Diff(t, http.StatusFound, w3.Result().StatusCode)
	testutil.Diff(t, "h4", w4.Result().Header.Get("test"))

	// Match to Path
	r5 := httptest.NewRequest(http.MethodGet, "http://test.com/test5/foo", nil)
	w5 := httptest.NewRecorder()
	handler.ServeHTTP(w5, r5)
	testutil.Diff(t, http.StatusFound, w5.Result().StatusCode)
	testutil.Diff(t, "h5", w5.Result().Header.Get("test"))

	// Match to FilePath
	r6 := httptest.NewRequest(http.MethodGet, "http://test.com/test6/foo", nil)
	w6 := httptest.NewRecorder()
	handler.ServeHTTP(w6, r6)
	testutil.Diff(t, http.StatusFound, w6.Result().StatusCode)
	testutil.Diff(t, "h6", w6.Result().Header.Get("test"))

	// Match to Regex
	r7 := httptest.NewRequest(http.MethodGet, "http://test.com/test/123/foo/bar", nil)
	w7 := httptest.NewRecorder()
	handler.ServeHTTP(w7, r7)
	testutil.Diff(t, http.StatusFound, w7.Result().StatusCode)
	testutil.Diff(t, "h7", w7.Result().Header.Get("test"))

	// Match to RegexPOSIX
	r8 := httptest.NewRequest(http.MethodGet, "http://test.com/test/abc/foo/bar", nil)
	w8 := httptest.NewRecorder()
	handler.ServeHTTP(w8, r8)
	testutil.Diff(t, http.StatusFound, w8.Result().StatusCode)
	testutil.Diff(t, "h8", w8.Result().Header.Get("test"))

	// Match to none
	r9 := httptest.NewRequest(http.MethodGet, "http://test.com/none", nil)
	w9 := httptest.NewRecorder()
	handler.ServeHTTP(w9, r9)
	b9, _ := io.ReadAll(w9.Body)
	testutil.Diff(t, http.StatusNotFound, w9.Result().StatusCode)
	testutil.Diff(t, `{"status":404,"statusText":"Not Found"}`, string(b9))

}

func TestTrimPrefix(t *testing.T) {

	configs := []string{
		testDataDir + "config-trim-prefix.yaml",
	}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "core/v1",
		Kind:       "ReverseProxyHandler",
		Name:       "default",
		Namespace:  "",
	}
	handler, err := api.ReferTypedObject[http.Handler](server, ref)
	testutil.DiffError(t, nil, nil, err)

	svr := &http.Server{
		Addr: ":10001",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("test", "ok")
			w.WriteHeader(http.StatusFound)
			w.Write([]byte("test"))
		}),
	}
	go func() { svr.ListenAndServe() }()
	time.Sleep(time.Second) // Wait a little until server starts.
	defer svr.Close()

	r1 := httptest.NewRequest(http.MethodGet, "http://test.com/pre/test", nil)
	w1 := httptest.NewRecorder()
	handler.ServeHTTP(w1, r1)
	b1, _ := io.ReadAll(w1.Body)
	testutil.Diff(t, http.StatusFound, w1.Result().StatusCode)
	testutil.Diff(t, "ok", w1.Result().Header.Get("test"))
	testutil.Diff(t, "test", string(b1))

	r2 := httptest.NewRequest(http.MethodGet, "http://test.com/none", nil)
	w2 := httptest.NewRecorder()
	handler.ServeHTTP(w2, r2)
	b2, _ := io.ReadAll(w2.Body)
	testutil.Diff(t, http.StatusNotFound, w2.Result().StatusCode)
	testutil.Diff(t, "", w2.Result().Header.Get("test"))
	testutil.Diff(t, `{"status":404,"statusText":"Not Found"}`, string(b2))

}

func TestRoundTripper(t *testing.T) {

	configs := []string{
		testDataDir + "config-roundtripper.yaml",
	}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	testRef := &kernel.Reference{
		APIVersion: "container/v1",
		Kind:       "TestRoundTripper",
		Name:       "test",
		Namespace:  "",
	}
	testRT := &testRoundTripper{}
	common.PostTestResource(server, testRef, testRT)

	ref := &kernel.Reference{
		APIVersion: "core/v1",
		Kind:       "ReverseProxyHandler",
		Name:       "default",
		Namespace:  "",
	}
	handler, err := api.ReferTypedObject[http.Handler](server, ref)
	testutil.DiffError(t, nil, nil, err)

	r1 := httptest.NewRequest(http.MethodGet, "http://test.com/test", nil)
	w1 := httptest.NewRecorder()
	handler.ServeHTTP(w1, r1)
	testutil.Diff(t, http.StatusFound, w1.Result().StatusCode)
	testutil.Diff(t, "ok", w1.Result().Header.Get("test"))
	testutil.Diff(t, 1, testRT.called)

	r2 := httptest.NewRequest(http.MethodGet, "http://test.com/none", nil)
	w2 := httptest.NewRecorder()
	handler.ServeHTTP(w2, r2)
	b2, _ := io.ReadAll(w2.Body)
	testutil.Diff(t, http.StatusNotFound, w2.Result().StatusCode)
	testutil.Diff(t, "", w2.Result().Header.Get("test"))
	testutil.Diff(t, `{"status":404,"statusText":"Not Found"}`, string(b2))

}

func TestTripperware(t *testing.T) {

	configs := []string{
		testDataDir + "config-tripperware.yaml",
	}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	testRef := &kernel.Reference{
		APIVersion: "container/v1",
		Kind:       "TestRoundTripper",
		Name:       "test",
		Namespace:  "",
	}
	testRT := &testRoundTripper{}
	common.PostTestResource(server, testRef, testRT)
	testTWRef1 := &kernel.Reference{
		APIVersion: "container/v1",
		Kind:       "TestTripperware",
		Name:       "test1",
		Namespace:  "",
	}
	testTW1 := &testTripperware{}
	common.PostTestResource(server, testTWRef1, testTW1)
	testTWRef2 := &kernel.Reference{
		APIVersion: "container/v1",
		Kind:       "TestTripperware",
		Name:       "test2",
		Namespace:  "",
	}
	testTW2 := &testTripperware{}
	common.PostTestResource(server, testTWRef2, testTW2)

	ref := &kernel.Reference{
		APIVersion: "core/v1",
		Kind:       "ReverseProxyHandler",
		Name:       "default",
		Namespace:  "",
	}
	handler, err := api.ReferTypedObject[http.Handler](server, ref)
	testutil.DiffError(t, nil, nil, err)

	r1 := httptest.NewRequest(http.MethodGet, "http://test.com/test", nil)
	w1 := httptest.NewRecorder()
	handler.ServeHTTP(w1, r1)
	testutil.Diff(t, http.StatusFound, w1.Result().StatusCode)
	testutil.Diff(t, "ok", w1.Result().Header.Get("test"))
	testutil.Diff(t, 1, testRT.called)
	testutil.Diff(t, 1, testTW1.called)
	testutil.Diff(t, 1, testTW2.called)

	r2 := httptest.NewRequest(http.MethodGet, "http://test.com/none", nil)
	w2 := httptest.NewRecorder()
	handler.ServeHTTP(w2, r2)

	b2, _ := io.ReadAll(w2.Body)
	testutil.Diff(t, http.StatusNotFound, w2.Result().StatusCode)
	testutil.Diff(t, "", w2.Result().Header.Get("test"))
	testutil.Diff(t, `{"status":404,"statusText":"Not Found"}`, string(b2))

}

func TestLBMaglev(t *testing.T) {

	configs := []string{
		testDataDir + "config-lb-maglev.yaml",
	}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "core/v1",
		Kind:       "ReverseProxyHandler",
		Name:       "default",
		Namespace:  "",
	}
	handler, err := api.ReferTypedObject[http.Handler](server, ref)
	testutil.DiffError(t, nil, nil, err)

	var count1, count2, count3 int
	svr1 := &http.Server{
		Addr: ":10001",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			count1 += 1
			w.WriteHeader(http.StatusOK)
		}),
	}
	svr2 := &http.Server{
		Addr: ":10002",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			count2 += 1
			w.WriteHeader(http.StatusOK)
		}),
	}
	svr3 := &http.Server{
		Addr: ":10003",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			count3 += 1
			w.WriteHeader(http.StatusOK)
		}),
	}
	go func() { svr1.ListenAndServe() }()
	go func() { svr2.ListenAndServe() }()
	go func() { svr3.ListenAndServe() }()
	time.Sleep(time.Second) // Wait a little until server starts.
	defer svr1.Close()
	defer svr2.Close()
	defer svr3.Close()

	// ?proxy=baz reached svr2.
	r1 := httptest.NewRequest(http.MethodGet, "http://test.com/test?proxy=baz", nil)
	for i := 0; i < 10; i++ {
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, r1)
		testutil.Diff(t, http.StatusOK, w.Result().StatusCode)
		testutil.Diff(t, i+1, count2)
	}

	// ?proxy=bar reached svr3.
	r2 := httptest.NewRequest(http.MethodGet, "http://test.com/test?proxy=bar", nil)
	for i := 0; i < 10; i++ {
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, r2)
		testutil.Diff(t, http.StatusOK, w.Result().StatusCode)
		testutil.Diff(t, i+1, count3)
	}

	// ?proxy=FooBar reached svr1.
	r3 := httptest.NewRequest(http.MethodGet, "http://test.com/test?proxy=proxy", nil)
	for i := 0; i < 10; i++ {
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, r3)
		testutil.Diff(t, http.StatusOK, w.Result().StatusCode)
		testutil.Diff(t, i+1, count1)
	}

	count1, count2, count3 = 0, 0, 0
	n := 1200
	for i := 0; i < n; i++ {
		r := httptest.NewRequest(http.MethodGet, "http://test.com/test?proxy="+strconv.Itoa(i), nil)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, r)
		testutil.Diff(t, http.StatusOK, w.Result().StatusCode)
	}
	testutil.Diff(t, float64(1*n/(1+2+3)), float64(count1), cmpopts.EquateApprox(0, 10))
	testutil.Diff(t, float64(2*n/(1+2+3)), float64(count2), cmpopts.EquateApprox(0, 20))
	testutil.Diff(t, float64(3*n/(1+2+3)), float64(count3), cmpopts.EquateApprox(0, 30))

}

func TestLBRingHash(t *testing.T) {

	configs := []string{
		testDataDir + "config-lb-ringhash.yaml",
	}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "core/v1",
		Kind:       "ReverseProxyHandler",
		Name:       "default",
		Namespace:  "",
	}
	handler, err := api.ReferTypedObject[http.Handler](server, ref)
	testutil.DiffError(t, nil, nil, err)

	var count1, count2, count3 int
	svr1 := &http.Server{
		Addr: ":10001",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			count1 += 1
			w.WriteHeader(http.StatusOK)
		}),
	}
	svr2 := &http.Server{
		Addr: ":10002",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			count2 += 1
			w.WriteHeader(http.StatusOK)
		}),
	}
	svr3 := &http.Server{
		Addr: ":10003",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			count3 += 1
			w.WriteHeader(http.StatusOK)
		}),
	}
	go func() { svr1.ListenAndServe() }()
	go func() { svr2.ListenAndServe() }()
	go func() { svr3.ListenAndServe() }()
	time.Sleep(time.Second) // Wait a little until server starts.
	defer svr1.Close()
	defer svr2.Close()
	defer svr3.Close()

	// ?proxy=alice reached svr2.
	r1 := httptest.NewRequest(http.MethodGet, "http://test.com/test?proxy=alice", nil)
	for i := 0; i < 10; i++ {
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, r1)
		testutil.Diff(t, http.StatusOK, w.Result().StatusCode)
		testutil.Diff(t, i+1, count2)
	}

	// ?proxy=bob reached svr3.
	r2 := httptest.NewRequest(http.MethodGet, "http://test.com/test?proxy=bob", nil)
	for i := 0; i < 10; i++ {
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, r2)
		testutil.Diff(t, http.StatusOK, w.Result().StatusCode)
		testutil.Diff(t, i+1, count3)
	}

	// ?proxy=FooBar reached svr2.
	r3 := httptest.NewRequest(http.MethodGet, "http://test.com/test?proxy=xyz", nil)
	for i := 0; i < 10; i++ {
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, r3)
		testutil.Diff(t, http.StatusOK, w.Result().StatusCode)
		testutil.Diff(t, i+1, count1)
	}

	count1, count2, count3 = 0, 0, 0
	n := 1200
	for i := 0; i < n; i++ {
		r := httptest.NewRequest(http.MethodGet, "http://test.com/test?proxy="+strconv.Itoa(i), nil)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, r)
		testutil.Diff(t, http.StatusOK, w.Result().StatusCode)
	}
	testutil.Diff(t, float64(1*n/(1+2+3)), float64(count1), cmpopts.EquateApprox(0, 50))
	testutil.Diff(t, float64(2*n/(1+2+3)), float64(count2), cmpopts.EquateApprox(0, 50))
	testutil.Diff(t, float64(3*n/(1+2+3)), float64(count3), cmpopts.EquateApprox(0, 50))

}

func TestLBDirectHash(t *testing.T) {

	configs := []string{
		testDataDir + "config-lb-directhash.yaml",
	}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "core/v1",
		Kind:       "ReverseProxyHandler",
		Name:       "default",
		Namespace:  "",
	}
	handler, err := api.ReferTypedObject[http.Handler](server, ref)
	testutil.DiffError(t, nil, nil, err)

	var count1, count2, count3 int
	svr1 := &http.Server{
		Addr: ":10001",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			count1 += 1
			w.WriteHeader(http.StatusOK)
		}),
	}
	svr2 := &http.Server{
		Addr: ":10002",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			count2 += 1
			w.WriteHeader(http.StatusOK)
		}),
	}
	svr3 := &http.Server{
		Addr: ":10003",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			count3 += 1
			w.WriteHeader(http.StatusOK)
		}),
	}
	go func() { svr1.ListenAndServe() }()
	go func() { svr2.ListenAndServe() }()
	go func() { svr3.ListenAndServe() }()
	time.Sleep(time.Second) // Wait a little until server starts.
	defer svr1.Close()
	defer svr2.Close()
	defer svr3.Close()

	// ?proxy=foo reached svr3.
	r1 := httptest.NewRequest(http.MethodGet, "http://test.com/test?proxy=foo", nil)
	for i := 0; i < 10; i++ {
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, r1)
		testutil.Diff(t, http.StatusOK, w.Result().StatusCode)
		testutil.Diff(t, i+1, count3)
	}

	// ?proxy=hoge reached svr2.
	r2 := httptest.NewRequest(http.MethodGet, "http://test.com/test?proxy=bob", nil)
	for i := 0; i < 10; i++ {
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, r2)
		testutil.Diff(t, http.StatusOK, w.Result().StatusCode)
		testutil.Diff(t, i+1, count2)
	}

	// ?proxy=abcdef reached svr1.
	r3 := httptest.NewRequest(http.MethodGet, "http://test.com/test?proxy=baz", nil)
	for i := 0; i < 10; i++ {
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, r3)
		testutil.Diff(t, http.StatusOK, w.Result().StatusCode)
		testutil.Diff(t, i+1, count1)
	}

	count1, count2, count3 = 0, 0, 0
	n := 1200
	for i := 0; i < n; i++ {
		r := httptest.NewRequest(http.MethodGet, "http://test.com/test?proxy="+strconv.Itoa(i), nil)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, r)
		testutil.Diff(t, http.StatusOK, w.Result().StatusCode)
	}
	testutil.Diff(t, float64(1*n/(1+2+3)), float64(count1), cmpopts.EquateApprox(0, 10))
	testutil.Diff(t, float64(2*n/(1+2+3)), float64(count2), cmpopts.EquateApprox(0, 20))
	testutil.Diff(t, float64(3*n/(1+2+3)), float64(count3), cmpopts.EquateApprox(0, 30))

}

func TestLBRoundRobin(t *testing.T) {

	configs := []string{
		testDataDir + "config-lb-roundrobin.yaml",
	}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "core/v1",
		Kind:       "ReverseProxyHandler",
		Name:       "default",
		Namespace:  "",
	}
	handler, err := api.ReferTypedObject[http.Handler](server, ref)
	testutil.DiffError(t, nil, nil, err)

	var count1, count2, count3 int
	svr1 := &http.Server{
		Addr: ":10001",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			count1 += 1
			w.WriteHeader(http.StatusOK)
		}),
	}
	svr2 := &http.Server{
		Addr: ":10002",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			count2 += 1
			w.WriteHeader(http.StatusOK)
		}),
	}
	svr3 := &http.Server{
		Addr: ":10003",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			count3 += 1
			w.WriteHeader(http.StatusOK)
		}),
	}
	go func() { svr1.ListenAndServe() }()
	go func() { svr2.ListenAndServe() }()
	go func() { svr3.ListenAndServe() }()
	time.Sleep(time.Second) // Wait a little until server starts.
	defer svr1.Close()
	defer svr2.Close()
	defer svr3.Close()

	r := httptest.NewRequest(http.MethodGet, "http://test.com/test", nil)
	for i := 0; i < 10; i++ {
		w1 := httptest.NewRecorder()
		handler.ServeHTTP(w1, r)
		testutil.Diff(t, http.StatusOK, w1.Result().StatusCode)
		testutil.Diff(t, i+1, count1)

		w2 := httptest.NewRecorder()
		handler.ServeHTTP(w2, r) // Weight 2
		handler.ServeHTTP(w2, r) // Weight 2
		testutil.Diff(t, http.StatusOK, w2.Result().StatusCode)
		testutil.Diff(t, 2*(i+1), count2)

		w3 := httptest.NewRecorder()
		handler.ServeHTTP(w3, r) // Weight 3
		handler.ServeHTTP(w3, r) // Weight 3
		handler.ServeHTTP(w3, r) // Weight 3
		testutil.Diff(t, http.StatusOK, w3.Result().StatusCode)
		testutil.Diff(t, 3*(i+1), count3)
	}

}

func TestLBRandom(t *testing.T) {

	configs := []string{
		testDataDir + "config-lb-random.yaml",
	}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "core/v1",
		Kind:       "ReverseProxyHandler",
		Name:       "default",
		Namespace:  "",
	}
	handler, err := api.ReferTypedObject[http.Handler](server, ref)
	testutil.DiffError(t, nil, nil, err)

	var count1, count2, count3 int
	svr1 := &http.Server{
		Addr: ":10001",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			count1 += 1
			w.WriteHeader(http.StatusOK)
		}),
	}
	svr2 := &http.Server{
		Addr: ":10002",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			count2 += 1
			w.WriteHeader(http.StatusOK)
		}),
	}
	svr3 := &http.Server{
		Addr: ":10003",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			count3 += 1
			w.WriteHeader(http.StatusOK)
		}),
	}
	go func() { svr1.ListenAndServe() }()
	go func() { svr2.ListenAndServe() }()
	go func() { svr3.ListenAndServe() }()
	time.Sleep(time.Second) // Wait a little until server starts.
	defer svr1.Close()
	defer svr2.Close()
	defer svr3.Close()

	r := httptest.NewRequest(http.MethodGet, "http://test.com/test", nil)

	n := 1800
	for i := 0; i < n; i++ {
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, r)
		testutil.Diff(t, http.StatusOK, w.Result().StatusCode)
	}
	testutil.Diff(t, float64(1*n/(1+2+3)), float64(count1), cmpopts.EquateApprox(0, 50))
	testutil.Diff(t, float64(2*n/(1+2+3)), float64(count2), cmpopts.EquateApprox(0, 50))
	testutil.Diff(t, float64(3*n/(1+2+3)), float64(count3), cmpopts.EquateApprox(0, 50))

}

func TestLBMatchHost(t *testing.T) {

	configs := []string{
		testDataDir + "config-lb-match-host.yaml",
	}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "core/v1",
		Kind:       "ReverseProxyHandler",
		Name:       "default",
		Namespace:  "",
	}
	handler, err := api.ReferTypedObject[http.Handler](server, ref)
	testutil.DiffError(t, nil, nil, err)

	svr := &http.Server{
		Addr: ":10001",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("test", "ok")
			w.WriteHeader(http.StatusFound)
			w.Write([]byte("test"))
		}),
	}
	go func() { svr.ListenAndServe() }()
	time.Sleep(time.Second) // Wait a little until server starts.
	defer svr.Close()

	// test.com matches.
	r1 := httptest.NewRequest(http.MethodGet, "http://test.com/", nil)
	w1 := httptest.NewRecorder()
	handler.ServeHTTP(w1, r1)
	b1, _ := io.ReadAll(w1.Body)
	testutil.Diff(t, http.StatusFound, w1.Result().StatusCode)
	testutil.Diff(t, "ok", w1.Result().Header.Get("test"))
	testutil.Diff(t, "test", string(b1))

	// wrong.com not matches.
	r2 := httptest.NewRequest(http.MethodGet, "http://wrong.com/", nil)
	w2 := httptest.NewRecorder()
	handler.ServeHTTP(w2, r2)
	b2, _ := io.ReadAll(w2.Body)
	testutil.Diff(t, http.StatusNotFound, w2.Result().StatusCode)
	testutil.Diff(t, "", w2.Result().Header.Get("test"))
	testutil.Diff(t, `{"status":404,"statusText":"Not Found"}`, string(b2))

}

func TestLBMatchMethod(t *testing.T) {

	configs := []string{
		testDataDir + "config-lb-match-method.yaml",
	}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "core/v1",
		Kind:       "ReverseProxyHandler",
		Name:       "default",
		Namespace:  "",
	}
	handler, err := api.ReferTypedObject[http.Handler](server, ref)
	testutil.DiffError(t, nil, nil, err)

	svr := &http.Server{
		Addr: ":10001",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("test", "ok")
			w.WriteHeader(http.StatusFound)
			w.Write([]byte("test"))
		}),
	}
	go func() { svr.ListenAndServe() }()
	time.Sleep(time.Second) // Wait a little until server starts.
	defer svr.Close()

	// Get matches.
	r1 := httptest.NewRequest(http.MethodGet, "http://test.com/", nil)
	w1 := httptest.NewRecorder()
	handler.ServeHTTP(w1, r1)
	b1, _ := io.ReadAll(w1.Body)
	testutil.Diff(t, http.StatusFound, w1.Result().StatusCode)
	testutil.Diff(t, "ok", w1.Result().Header.Get("test"))
	testutil.Diff(t, "test", string(b1))

	// Head matches.
	r2 := httptest.NewRequest(http.MethodPut, "http://test.com/", nil)
	w2 := httptest.NewRecorder()
	handler.ServeHTTP(w2, r2)
	b2, _ := io.ReadAll(w2.Body)
	testutil.Diff(t, http.StatusFound, w2.Result().StatusCode)
	testutil.Diff(t, "ok", w2.Result().Header.Get("test"))
	testutil.Diff(t, "test", string(b2))

	// Post not matches.
	r3 := httptest.NewRequest(http.MethodPost, "http://test.com/", nil)
	w3 := httptest.NewRecorder()
	handler.ServeHTTP(w3, r3)
	b3, _ := io.ReadAll(w3.Body)
	testutil.Diff(t, http.StatusNotFound, w3.Result().StatusCode)
	testutil.Diff(t, "", w3.Result().Header.Get("test"))
	testutil.Diff(t, `{"status":404,"statusText":"Not Found"}`, string(b3))

}

func TestLBMatchPathParam(t *testing.T) {

	configs := []string{
		testDataDir + "config-lb-match-path-param.yaml",
	}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "core/v1",
		Kind:       "ReverseProxyHandler",
		Name:       "default",
		Namespace:  "",
	}
	handler, err := api.ReferTypedObject[http.Handler](server, ref)
	testutil.DiffError(t, nil, nil, err)
	mux := &http.ServeMux{}
	mux.Handle("/test/{param}", handler)

	svr := &http.Server{
		Addr: ":10001",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("test", "ok")
			w.WriteHeader(http.StatusFound)
			w.Write([]byte("test"))
		}),
	}
	go func() { svr.ListenAndServe() }()
	time.Sleep(time.Second) // Wait a little until server starts.
	defer svr.Close()

	// Path parameter "foo" matches.
	r1 := httptest.NewRequest(http.MethodGet, "http://test.com/test/foo", nil)
	w1 := httptest.NewRecorder()
	mux.ServeHTTP(w1, r1)
	b1, _ := io.ReadAll(w1.Body)
	testutil.Diff(t, http.StatusFound, w1.Result().StatusCode)
	testutil.Diff(t, "ok", w1.Result().Header.Get("test"))
	testutil.Diff(t, "test", string(b1))

	// Path parameter "bar" matches.
	r2 := httptest.NewRequest(http.MethodGet, "http://test.com/test/foo", nil)
	w2 := httptest.NewRecorder()
	mux.ServeHTTP(w2, r2)
	b2, _ := io.ReadAll(w2.Body)
	testutil.Diff(t, http.StatusFound, w2.Result().StatusCode)
	testutil.Diff(t, "ok", w2.Result().Header.Get("test"))
	testutil.Diff(t, "test", string(b2))

	// Path parameter "baz" not matches.
	r3 := httptest.NewRequest(http.MethodPost, "http://test.com/test/baz", nil)
	w3 := httptest.NewRecorder()
	mux.ServeHTTP(w3, r3)
	b3, _ := io.ReadAll(w3.Body)
	testutil.Diff(t, http.StatusNotFound, w3.Result().StatusCode)
	testutil.Diff(t, "", w3.Result().Header.Get("test"))
	testutil.Diff(t, `{"status":404,"statusText":"Not Found"}`, string(b3))

}

func TestLBMatchHeader(t *testing.T) {

	configs := []string{
		testDataDir + "config-lb-match-header.yaml",
	}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "core/v1",
		Kind:       "ReverseProxyHandler",
		Name:       "default",
		Namespace:  "",
	}
	handler, err := api.ReferTypedObject[http.Handler](server, ref)
	testutil.DiffError(t, nil, nil, err)

	svr := &http.Server{
		Addr: ":10001",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("test", "ok")
			w.WriteHeader(http.StatusFound)
			w.Write([]byte("test"))
		}),
	}
	go func() { svr.ListenAndServe() }()
	time.Sleep(time.Second) // Wait a little until server starts.
	defer svr.Close()

	// Header "foo" matches.
	r1 := httptest.NewRequest(http.MethodGet, "http://test.com/", nil)
	r1.Header.Add("param", "foo")
	w1 := httptest.NewRecorder()
	handler.ServeHTTP(w1, r1)
	b1, _ := io.ReadAll(w1.Body)
	testutil.Diff(t, http.StatusFound, w1.Result().StatusCode)
	testutil.Diff(t, "ok", w1.Result().Header.Get("test"))
	testutil.Diff(t, "test", string(b1))

	// Header "bar" matches.
	r2 := httptest.NewRequest(http.MethodGet, "http://test.com/", nil)
	r2.Header.Add("param", "bar")
	w2 := httptest.NewRecorder()
	handler.ServeHTTP(w2, r2)
	b2, _ := io.ReadAll(w2.Body)
	testutil.Diff(t, http.StatusFound, w2.Result().StatusCode)
	testutil.Diff(t, "ok", w2.Result().Header.Get("test"))
	testutil.Diff(t, "test", string(b2))

	// Header "baz" not matches.
	r3 := httptest.NewRequest(http.MethodPost, "http://test.com/", nil)
	r3.Header.Add("param", "baz")
	w3 := httptest.NewRecorder()
	handler.ServeHTTP(w3, r3)
	b3, _ := io.ReadAll(w3.Body)
	testutil.Diff(t, http.StatusNotFound, w3.Result().StatusCode)
	testutil.Diff(t, "", w3.Result().Header.Get("test"))
	testutil.Diff(t, `{"status":404,"statusText":"Not Found"}`, string(b3))

}

func TestLBMatchPaths(t *testing.T) {

	configs := []string{
		testDataDir + "config-lb-match-paths.yaml",
	}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "core/v1",
		Kind:       "ReverseProxyHandler",
	}
	handler, err := api.ReferTypedObject[http.Handler](server, ref)
	testutil.DiffError(t, nil, nil, err)

	svr := &http.Server{
		Addr: ":10001",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("test", "ok")
			w.WriteHeader(http.StatusFound)
			w.Write([]byte("test"))
		}),
	}
	go func() { svr.ListenAndServe() }()
	time.Sleep(time.Second) // Wait a little until server starts.
	defer svr.Close()

	// "/foo" matches.
	r1 := httptest.NewRequest(http.MethodGet, "http://test.com/foo", nil)
	w1 := httptest.NewRecorder()
	handler.ServeHTTP(w1, r1)
	b1, _ := io.ReadAll(w1.Body)
	testutil.Diff(t, http.StatusFound, w1.Result().StatusCode)
	testutil.Diff(t, "ok", w1.Result().Header.Get("test"))
	testutil.Diff(t, "test", string(b1))

	// "/bar" matches.
	r2 := httptest.NewRequest(http.MethodGet, "http://test.com/bar", nil)
	w2 := httptest.NewRecorder()
	handler.ServeHTTP(w2, r2)
	b2, _ := io.ReadAll(w2.Body)
	testutil.Diff(t, http.StatusFound, w2.Result().StatusCode)
	testutil.Diff(t, "ok", w2.Result().Header.Get("test"))
	testutil.Diff(t, "test", string(b2))

	// "/baz" not matches.
	r3 := httptest.NewRequest(http.MethodPost, "http://test.com/baz", nil)
	w3 := httptest.NewRecorder()
	handler.ServeHTTP(w3, r3)
	b3, _ := io.ReadAll(w3.Body)
	testutil.Diff(t, http.StatusNotFound, w3.Result().StatusCode)
	testutil.Diff(t, "", w3.Result().Header.Get("test"))
	testutil.Diff(t, `{"status":404,"statusText":"Not Found"}`, string(b3))
}

func TestLBMatchQuery(t *testing.T) {

	configs := []string{
		testDataDir + "config-lb-match-query.yaml",
	}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "core/v1",
		Kind:       "ReverseProxyHandler",
	}
	handler, err := api.ReferTypedObject[http.Handler](server, ref)
	testutil.DiffError(t, nil, nil, err)

	svr := &http.Server{
		Addr: ":10001",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("test", "ok")
			w.WriteHeader(http.StatusFound)
			w.Write([]byte("test"))
		}),
	}
	go func() { svr.ListenAndServe() }()
	time.Sleep(time.Second) // Wait a little until server starts.
	defer svr.Close()

	// Query parameter "foo" matches.
	r1 := httptest.NewRequest(http.MethodGet, "http://test.com/?param=foo", nil)
	w1 := httptest.NewRecorder()
	handler.ServeHTTP(w1, r1)
	b1, _ := io.ReadAll(w1.Body)
	testutil.Diff(t, http.StatusFound, w1.Result().StatusCode)
	testutil.Diff(t, "ok", w1.Result().Header.Get("test"))
	testutil.Diff(t, "test", string(b1))

	// Query parameter "bar" matches.
	r2 := httptest.NewRequest(http.MethodGet, "http://test.com/?param=bar", nil)
	w2 := httptest.NewRecorder()
	handler.ServeHTTP(w2, r2)
	b2, _ := io.ReadAll(w2.Body)
	testutil.Diff(t, http.StatusFound, w2.Result().StatusCode)
	testutil.Diff(t, "ok", w2.Result().Header.Get("test"))
	testutil.Diff(t, "test", string(b2))

	// Query parameter "baz" not matches.
	r3 := httptest.NewRequest(http.MethodPost, "http://test.com/?param=baz", nil)
	w3 := httptest.NewRecorder()
	handler.ServeHTTP(w3, r3)
	b3, _ := io.ReadAll(w3.Body)
	testutil.Diff(t, http.StatusNotFound, w3.Result().StatusCode)
	testutil.Diff(t, "", w3.Result().Header.Get("test"))
	testutil.Diff(t, `{"status":404,"statusText":"Not Found"}`, string(b3))

}
