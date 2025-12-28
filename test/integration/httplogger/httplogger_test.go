// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

//go:build integration

package httplogger_test

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/aileron-gateway/aileron-gateway/apis/kernel"
	"github.com/aileron-gateway/aileron-gateway/cmd/aileron/app"
	"github.com/aileron-gateway/aileron-gateway/core"
	"github.com/aileron-gateway/aileron-gateway/internal/testutil"
	"github.com/aileron-gateway/aileron-gateway/kernel/api"
	"github.com/aileron-gateway/aileron-gateway/test/integration/common"
)

func TestMaskQuery(t *testing.T) {

	configs := []string{
		testDataDir + "config-mask-query.yaml",
	}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	tmp := os.Stdout
	defer func() { os.Stdout = tmp }()
	r, w, _ := os.Pipe()
	os.Stdout = w

	ref := &kernel.Reference{
		APIVersion: "core/v1",
		Kind:       "HTTPLogger",
		Name:       "default",
		Namespace:  "",
	}
	m, err := api.ReferTypedObject[core.Middleware](server, ref)
	testutil.DiffError(t, nil, nil, err)
	testHTTPLoggerMiddleware(t, m)

	h := m.Middleware(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("test"))
		}),
	)

	req := httptest.NewRequest(http.MethodGet, "/test?foo=bar&alice=bob", nil)
	res := httptest.NewRecorder()
	h.ServeHTTP(res, req)
	b, _ := io.ReadAll(res.Body)
	testutil.Diff(t, http.StatusOK, res.Result().StatusCode)
	testutil.Diff(t, []byte("test"), b)

	w.Close()
	var buf bytes.Buffer
	buf.ReadFrom(r)
	str := buf.String()
	t.Log(str, "\n")
	testutil.Diff(t, true, strings.Contains(str, `foo=##MASKED##`))
	testutil.Diff(t, true, strings.Contains(str, `alice=%%MASKED%%`))

}

func TestMaskHeaders_Middleware(t *testing.T) {

	configs := []string{
		testDataDir + "config-mask-headers.yaml",
	}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	tmp := os.Stdout
	defer func() { os.Stdout = tmp }()
	r, w, _ := os.Pipe()
	os.Stdout = w

	ref := &kernel.Reference{
		APIVersion: "core/v1",
		Kind:       "HTTPLogger",
		Name:       "default",
		Namespace:  "",
	}
	m, err := api.ReferTypedObject[core.Middleware](server, ref)
	testutil.DiffError(t, nil, nil, err)
	testHTTPLoggerMiddleware(t, m)

	h := m.Middleware(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("test", "ok")
			w.Header()["Foo"] = []string{"value1", "value2"}
			w.Header()["Bar"] = []string{"value1", "value2"}
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("test"))
		}),
	)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer TEST-DUMMY-TOKEN")
	req.Header.Set("Alice", "Bob")
	res := httptest.NewRecorder()
	h.ServeHTTP(res, req)
	b, _ := io.ReadAll(res.Body)
	testutil.Diff(t, http.StatusOK, res.Result().StatusCode)
	testutil.Diff(t, "ok", res.Header().Get("test"))
	testutil.Diff(t, []byte("test"), b)

	w.Close()
	var buf bytes.Buffer
	buf.ReadFrom(r)
	str := buf.String()
	t.Log(str, "\n")
	testutil.Diff(t, true, strings.Contains(str, `"Authorization":"##MASKED##"`))
	testutil.Diff(t, true, strings.Contains(str, `"Alice":"Bob"`))
	testutil.Diff(t, true, strings.Contains(str, `"Bar":"%%MASKED%%"`))
	testutil.Diff(t, true, strings.Contains(str, `"Foo":"$MASKED$"`))

}

func TestMaskHeaders_Tripperware(t *testing.T) {

	configs := []string{
		testDataDir + "config-mask-headers.yaml",
	}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	tmp := os.Stdout
	defer func() { os.Stdout = tmp }()
	r, w, _ := os.Pipe()
	os.Stdout = w

	ref := &kernel.Reference{
		APIVersion: "core/v1",
		Kind:       "HTTPLogger",
		Name:       "default",
		Namespace:  "",
	}
	tr, err := api.ReferTypedObject[core.Tripperware](server, ref)
	testutil.DiffError(t, nil, nil, err)
	testHTTPLoggerTripperware(t, tr)

	rt := tr.Tripperware(core.RoundTripperFunc(
		func(r *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Header: http.Header{
					"Test": {"ok"},
					"Foo":  {"value1", "value2"},
					"Bar":  {"value1", "value2"},
				},
				Body: io.NopCloser(bytes.NewReader([]byte("test"))),
			}, nil
		},
	))

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer TEST-DUMMY-TOKEN")
	req.Header.Set("Alice", "Bob")
	res, _ := rt.RoundTrip(req)
	b, _ := io.ReadAll(res.Body)
	testutil.Diff(t, http.StatusOK, res.StatusCode)
	testutil.Diff(t, "ok", res.Header.Get("test"))
	testutil.Diff(t, []byte("test"), b)

	w.Close()
	var buf bytes.Buffer
	buf.ReadFrom(r)
	str := buf.String()
	t.Log(str, "\n")
	testutil.Diff(t, true, strings.Contains(str, `"Authorization":"##MASKED##"`))
	testutil.Diff(t, true, strings.Contains(str, `"Alice":"Bob"`))
	testutil.Diff(t, true, strings.Contains(str, `"Bar":"%%MASKED%%"`))
	testutil.Diff(t, true, strings.Contains(str, `"Foo":"$MASKED$"`))

}

func TestMaskBodies_Middleware(t *testing.T) {

	configs := []string{
		testDataDir + "config-mask-bodies.yaml",
	}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	tmp := os.Stdout
	defer func() { os.Stdout = tmp }()
	r, w, _ := os.Pipe()
	os.Stdout = w

	ref := &kernel.Reference{
		APIVersion: "core/v1",
		Kind:       "HTTPLogger",
		Name:       "default",
		Namespace:  "",
	}
	m, err := api.ReferTypedObject[core.Middleware](server, ref)
	testutil.DiffError(t, nil, nil, err)
	testHTTPLoggerMiddleware(t, m)

	h := m.Middleware(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			body := []byte(`{"bar":"BAR","test":"response"}`)
			w.Header().Set("test", "ok")
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("Content-Length", strconv.Itoa(len(body)))
			w.WriteHeader(http.StatusOK)
			w.Write(body)
		}),
	)

	body := bytes.NewReader([]byte(`{"foo":"FOO","test":"request"}`))
	req := httptest.NewRequest(http.MethodGet, "/test", body)
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()
	h.ServeHTTP(res, req)
	b, _ := io.ReadAll(res.Body)
	testutil.Diff(t, http.StatusOK, res.Result().StatusCode)
	testutil.Diff(t, "ok", res.Header().Get("test"))
	testutil.Diff(t, `{"bar":"BAR","test":"response"}`, string(b))

	w.Close()
	var buf bytes.Buffer
	buf.ReadFrom(r)
	str := buf.String()
	t.Log(str, "\n")
	testutil.Diff(t, true, strings.Contains(str, `\"foo\":\"##MASKED##\"`))
	testutil.Diff(t, true, strings.Contains(str, `\"bar\":\"##MASKED##\"`))
	testutil.Diff(t, true, strings.Contains(str, `\"test\":\"request\"`))
	testutil.Diff(t, true, strings.Contains(str, `\"test\":\"response\"`))

}

func TestMaskBodies_Tripperware(t *testing.T) {

	configs := []string{
		testDataDir + "config-mask-bodies.yaml",
	}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	tmp := os.Stdout
	defer func() { os.Stdout = tmp }()
	r, w, _ := os.Pipe()
	os.Stdout = w

	ref := &kernel.Reference{
		APIVersion: "core/v1",
		Kind:       "HTTPLogger",
		Name:       "default",
		Namespace:  "",
	}
	tr, err := api.ReferTypedObject[core.Tripperware](server, ref)
	testutil.DiffError(t, nil, nil, err)
	testHTTPLoggerTripperware(t, tr)

	rt := tr.Tripperware(core.RoundTripperFunc(
		func(r *http.Request) (*http.Response, error) {
			body := []byte(`{"bar":"BAR","test":"response"}`)
			return &http.Response{
				ContentLength: 31,
				StatusCode:    http.StatusOK,
				Header: http.Header{
					"Test":         {"ok"},
					"Content-Type": {"application/json"},
					"Bar":          {"value1", "value2"},
					"Content-Size": {strconv.Itoa(len(body))},
				},
				Body: io.NopCloser(bytes.NewReader(body)),
			}, nil
		},
	))

	body := bytes.NewReader([]byte(`{"foo":"FOO","test":"request"}`))
	req := httptest.NewRequest(http.MethodGet, "/test", body)
	req.Header.Set("Content-Type", "application/json")
	res, _ := rt.RoundTrip(req)
	b, _ := io.ReadAll(res.Body)
	testutil.Diff(t, http.StatusOK, res.StatusCode)
	testutil.Diff(t, "ok", res.Header.Get("test"))
	testutil.Diff(t, `{"bar":"BAR","test":"response"}`, string(b))

	w.Close()
	var buf bytes.Buffer
	buf.ReadFrom(r)
	str := buf.String()
	t.Log(str, "\n")
	testutil.Diff(t, true, strings.Contains(str, `\"foo\":\"##MASKED##\"`))
	testutil.Diff(t, true, strings.Contains(str, `\"bar\":\"##MASKED##\"`))
	testutil.Diff(t, true, strings.Contains(str, `\"test\":\"request\"`))
	testutil.Diff(t, true, strings.Contains(str, `\"test\":\"response\"`))

}

func TestMaskBodiesJSON_Middleware(t *testing.T) {

	configs := []string{
		testDataDir + "config-mask-bodies-json.yaml",
	}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	tmp := os.Stdout
	defer func() { os.Stdout = tmp }()
	r, w, _ := os.Pipe()
	os.Stdout = w

	ref := &kernel.Reference{
		APIVersion: "core/v1",
		Kind:       "HTTPLogger",
		Name:       "default",
		Namespace:  "",
	}
	m, err := api.ReferTypedObject[core.Middleware](server, ref)
	testutil.DiffError(t, nil, nil, err)
	testHTTPLoggerMiddleware(t, m)

	h := m.Middleware(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			body := []byte(`{"bar":"BAR", "bob":{"bar":"BAR","test":"response"}}`)
			w.Header().Set("test", "ok")
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("Content-Length", strconv.Itoa(len(body)))
			w.WriteHeader(http.StatusOK)
			w.Write(body)
		}),
	)

	body := bytes.NewReader([]byte(`{"foo":"FOO", "alice":{"foo":"FOO","test":"request"}}`))
	req := httptest.NewRequest(http.MethodGet, "/test", body)
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()
	h.ServeHTTP(res, req)
	b, _ := io.ReadAll(res.Body)
	testutil.Diff(t, http.StatusOK, res.Result().StatusCode)
	testutil.Diff(t, "ok", res.Header().Get("test"))
	testutil.Diff(t, `{"bar":"BAR", "bob":{"bar":"BAR","test":"response"}}`, string(b))

	w.Close()
	var buf bytes.Buffer
	buf.ReadFrom(r)
	str := buf.String()
	t.Log(str, "\n")
	testutil.Diff(t, true, strings.Contains(str, `\"foo\":\"FOO\"`))
	testutil.Diff(t, true, strings.Contains(str, `\"bar\":\"BAR\"`))
	testutil.Diff(t, true, strings.Contains(str, `\"foo\":\"##MASKED##\"`))
	testutil.Diff(t, true, strings.Contains(str, `\"bar\":\"##MASKED##\"`))
	testutil.Diff(t, true, strings.Contains(str, `\"test\":\"request\"`))
	testutil.Diff(t, true, strings.Contains(str, `\"test\":\"response\"`))

}

func TestMaskBodiesJSON_Tripperware(t *testing.T) {

	configs := []string{
		testDataDir + "config-mask-bodies-json.yaml",
	}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	tmp := os.Stdout
	defer func() { os.Stdout = tmp }()
	r, w, _ := os.Pipe()
	os.Stdout = w

	ref := &kernel.Reference{
		APIVersion: "core/v1",
		Kind:       "HTTPLogger",
		Name:       "default",
		Namespace:  "",
	}
	tr, err := api.ReferTypedObject[core.Tripperware](server, ref)
	testutil.DiffError(t, nil, nil, err)
	testHTTPLoggerTripperware(t, tr)

	rt := tr.Tripperware(core.RoundTripperFunc(
		func(r *http.Request) (*http.Response, error) {
			body := []byte(`{"bar":"BAR", "bob":{"bar":"BAR","test":"response"}}`)
			return &http.Response{
				ContentLength: 52,
				StatusCode:    http.StatusOK,
				Header: http.Header{
					"Test":           {"ok"},
					"Content-Type":   {"application/json"},
					"Bar":            {"value1", "value2"},
					"Content-Length": {strconv.Itoa(len(body))},
				},
				Body: io.NopCloser(bytes.NewReader(body)),
			}, nil
		},
	))

	body := bytes.NewReader([]byte(`{"foo":"FOO", "alice":{"foo":"FOO","test":"request"}}`))
	req := httptest.NewRequest(http.MethodGet, "/test", body)
	req.Header.Set("Content-Type", "application/json")
	res, _ := rt.RoundTrip(req)
	b, _ := io.ReadAll(res.Body)
	testutil.Diff(t, http.StatusOK, res.StatusCode)
	testutil.Diff(t, "ok", res.Header.Get("test"))
	testutil.Diff(t, `{"bar":"BAR", "bob":{"bar":"BAR","test":"response"}}`, string(b))

	w.Close()
	var buf bytes.Buffer
	buf.ReadFrom(r)
	str := buf.String()
	t.Log(str, "\n")
	testutil.Diff(t, true, strings.Contains(str, `\"foo\":\"FOO\"`))
	testutil.Diff(t, true, strings.Contains(str, `\"bar\":\"BAR\"`))
	testutil.Diff(t, true, strings.Contains(str, `\"foo\":\"##MASKED##\"`))
	testutil.Diff(t, true, strings.Contains(str, `\"bar\":\"##MASKED##\"`))
	testutil.Diff(t, true, strings.Contains(str, `\"test\":\"request\"`))
	testutil.Diff(t, true, strings.Contains(str, `\"test\":\"response\"`))

}

func TestContentLengthNegativeOne_Tripperware(t *testing.T) {
	configs := []string{
		testDataDir + "config-mask-bodies.yaml",
	}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	tmp := os.Stdout
	defer func() { os.Stdout = tmp }()
	r, w, _ := os.Pipe()
	os.Stdout = w

	ref := &kernel.Reference{
		APIVersion: "core/v1",
		Kind:       "HTTPLogger",
		Name:       "default",
		Namespace:  "",
	}
	tr, err := api.ReferTypedObject[core.Tripperware](server, ref)
	testutil.DiffError(t, nil, nil, err)
	testHTTPLoggerTripperware(t, tr)

	rt := tr.Tripperware(core.RoundTripperFunc(
		func(r *http.Request) (*http.Response, error) {
			body := []byte(`{"bar":"BAR","test":"response"}`)
			return &http.Response{
				StatusCode:       http.StatusOK,
				ContentLength:    -1, // Content-Length: -1
				TransferEncoding: []string{"chunked"},
				Header: http.Header{
					"Test":         {"ok"},
					"Content-Type": {"application/json"},
				},
				Body: io.NopCloser(bytes.NewReader(body)),
			}, nil
		},
	))

	body := bytes.NewReader([]byte(`{"foo":"FOO","test":"request"}`))
	req := httptest.NewRequest(http.MethodGet, "/test", body)
	req.Header.Set("Content-Type", "application/json")
	res, _ := rt.RoundTrip(req)
	b, _ := io.ReadAll(res.Body)
	testutil.Diff(t, http.StatusOK, res.StatusCode)
	testutil.Diff(t, "ok", res.Header.Get("test"))
	testutil.Diff(t, `{"bar":"BAR","test":"response"}`, string(b))

	w.Close()
	var buf bytes.Buffer
	buf.ReadFrom(r)
	str := buf.String()
	t.Log(str, "\n")
	testutil.Diff(t, true, strings.Contains(str, `\"foo\":\"##MASKED##\"`))
	testutil.Diff(t, true, strings.Contains(str, `\"bar\":\"##MASKED##\"`))
	testutil.Diff(t, true, strings.Contains(str, `\"test\":\"request\"`))
	testutil.Diff(t, true, strings.Contains(str, `\"test\":\"response\"`))
}

func TestMaskBodies_Middleware_ContentLengthMinusOne(t *testing.T) {

	configs := []string{
		testDataDir + "config-mask-bodies.yaml",
	}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	tmp := os.Stdout
	defer func() { os.Stdout = tmp }()
	r, w, _ := os.Pipe()
	os.Stdout = w

	ref := &kernel.Reference{
		APIVersion: "core/v1",
		Kind:       "HTTPLogger",
		Name:       "default",
		Namespace:  "",
	}
	m, err := api.ReferTypedObject[core.Middleware](server, ref)
	testutil.DiffError(t, nil, nil, err)
	testHTTPLoggerMiddleware(t, m)

	h := m.Middleware(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			body := []byte(`{"bar":"BAR","test":"response"}`)
			w.Header().Set("test", "ok")
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("Transfer-Encoding", "chunked")
			w.Header().Set("Content-Length", "-1") // Content-Length: -1
			w.WriteHeader(http.StatusOK)
			w.Write(body)
		}),
	)

	body := bytes.NewReader([]byte(`{"foo":"FOO","test":"request"}`))
	req := httptest.NewRequest(http.MethodGet, "/test", body)
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()
	h.ServeHTTP(res, req)
	b, _ := io.ReadAll(res.Body)
	testutil.Diff(t, http.StatusOK, res.Result().StatusCode)
	testutil.Diff(t, "ok", res.Header().Get("test"))
	testutil.Diff(t, `{"bar":"BAR","test":"response"}`, string(b))

	w.Close()
	var buf bytes.Buffer
	buf.ReadFrom(r)
	str := buf.String()
	t.Log(str, "\n")
	testutil.Diff(t, true, strings.Contains(str, `\"foo\":\"##MASKED##\"`))
	testutil.Diff(t, true, strings.Contains(str, `\"bar\":\"##MASKED##\"`))
	testutil.Diff(t, true, strings.Contains(str, `\"test\":\"request\"`))
	testutil.Diff(t, true, strings.Contains(str, `\"test\":\"response\"`))

}

func TestGzipCompression_Middleware(t *testing.T) {
	configs := []string{
		testDataDir + "config-mask-bodies.yaml",
	}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "core/v1",
		Kind:       "HTTPLogger",
		Name:       "default",
		Namespace:  "",
	}
	m, err := api.ReferTypedObject[core.Middleware](server, ref)
	testutil.DiffError(t, nil, nil, err)
	testHTTPLoggerMiddleware(t, m)

	// Create the middleware handler
	h := m.Middleware(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			body := []byte(`{"bar":"BAR","test":"response"}`)
			w.Header().Set("test", "ok")
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("Transfer-Encoding", "chunked")
			w.Header().Set("Content-Encoding", "gzip") // Gzip compression
			// Send the response with Gzip compression
			gz := gzip.NewWriter(w)
			defer gz.Close()
			w.WriteHeader(http.StatusOK)
			gz.Write(body)
		}),
	)

	// Specify Gzip compression in the request
	body := bytes.NewReader([]byte(`{"foo":"FOO","test":"request"}`))
	req := httptest.NewRequest(http.MethodGet, "/test", body)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Encoding", "gzip") // Gzip compressed request
	res := httptest.NewRecorder()
	h.ServeHTTP(res, req)

	// Check the response status code and header
	b, _ := io.ReadAll(res.Body)

	// Check that the Content-Encoding header in the response is gzip
	testutil.Diff(t, http.StatusOK, res.Result().StatusCode)
	testutil.Diff(t, "ok", res.Header().Get("test"))
	testutil.Diff(t, "gzip", res.Header().Get("Content-Encoding")) // Verify it's compressed

	// Decompress the compressed response
	gzReader, err := gzip.NewReader(bytes.NewReader(b))
	testutil.Diff(t, nil, err)

	var uncompressedBody []byte
	uncompressedBody, err = io.ReadAll(gzReader)
	testutil.Diff(t, nil, err)

	// Check the decompressed response body
	testutil.Diff(t, `{"bar":"BAR","test":"response"}`, string(uncompressedBody))
}

func TestGzipCompression_Tripperware(t *testing.T) {
	configs := []string{
		testDataDir + "config-mask-bodies.yaml",
	}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "core/v1",
		Kind:       "HTTPLogger",
		Name:       "default",
		Namespace:  "",
	}
	tr, err := api.ReferTypedObject[core.Tripperware](server, ref)
	testutil.DiffError(t, nil, nil, err)
	testHTTPLoggerTripperware(t, tr)

	// Create the Tripperware handler
	rt := tr.Tripperware(core.RoundTripperFunc(
		func(r *http.Request) (*http.Response, error) {
			body := []byte(`{"bar":"BAR","test":"response"}`)
			// Return the Gzip compressed response
			var buf bytes.Buffer
			gz := gzip.NewWriter(&buf)
			gz.Write(body)
			gz.Close()

			return &http.Response{
				StatusCode:       http.StatusOK,
				ContentLength:    int64(buf.Len()), // Set the correct Content-Length
				TransferEncoding: []string{"chunked"},
				Header: http.Header{
					"Test":             {"ok"},
					"Content-Type":     {"application/json"},
					"Content-Encoding": {"gzip"}, // Specify Gzip compression
				},
				Body: io.NopCloser(&buf), // Gzip compressed response body
			}, nil
		},
	))

	// Specify Gzip compression in the request
	body := bytes.NewReader([]byte(`{"foo":"FOO","test":"request"}`))
	req := httptest.NewRequest(http.MethodGet, "/test", body)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Encoding", "gzip") // Gzip compressed request

	// Send the request via Tripperware
	res, _ := rt.RoundTrip(req)
	b, _ := io.ReadAll(res.Body)

	// Check the response status code and header
	testutil.Diff(t, http.StatusOK, res.StatusCode)
	testutil.Diff(t, "ok", res.Header.Get("test"))
	testutil.Diff(t, "gzip", res.Header.Get("Content-Encoding")) // Verify it's compressed

	// Decompress the compressed response
	gzReader, err := gzip.NewReader(bytes.NewReader(b))
	testutil.Diff(t, nil, err)

	var uncompressedBody []byte
	uncompressedBody, err = io.ReadAll(gzReader)
	testutil.Diff(t, nil, err)

	// Check the decompressed response body
	testutil.Diff(t, `{"bar":"BAR","test":"response"}`, string(uncompressedBody))
}
