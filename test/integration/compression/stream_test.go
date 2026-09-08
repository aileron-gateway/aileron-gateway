// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

//go:build integration

package compression_test

import (
	"bytes"
	"compress/gzip"
	"crypto/md5"
	"encoding/hex"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/aileron-gateway/aileron-gateway/apis/kernel"
	"github.com/aileron-gateway/aileron-gateway/cmd/aileron/app"
	"github.com/aileron-gateway/aileron-gateway/core"
	"github.com/aileron-gateway/aileron-gateway/kernel/api"
	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
	"github.com/aileron-gateway/aileron-gateway/test/integration/common"
)

// testDataDir is the path to the test data.
var testDataDir = "../../../test/integration/compression/"

func TestStream_SSE(t *testing.T) {

	configs := []string{
		testDataDir + "config-stream.yaml",
	}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	proxyRef := &kernel.Reference{
		APIVersion: "core/v1",
		Kind:       "ReverseProxyHandler",
	}
	proxy, err := api.ReferTypedObject[http.Handler](server, proxyRef)
	testutil.DiffError(t, nil, nil, err)
	middleRef := &kernel.Reference{
		APIVersion: "app/v1",
		Kind:       "CompressionMiddleware",
	}
	middle, err := api.ReferTypedObject[core.Middleware](server, middleRef)
	testutil.DiffError(t, nil, nil, err)
	handler := middle.Middleware(proxy)

	svr := &http.Server{
		Addr: ":12301",
		Handler: &returnChunkedHandler{
			header: http.Header{
				"Content-Type":      []string{"text/event-stream"},
				"Transfer-Encoding": []string{"identity"},
			},
			interval: 10 * time.Millisecond,
			bodies:   []string{"1", "2", "3", "4", "5"},
		},
	}
	go func() { svr.ListenAndServe() }()
	time.Sleep(time.Second) // Wait a little until server starts.
	defer svr.Close()

	r := httptest.NewRequest(http.MethodGet, "http://test.com/test", nil)
	r.Header.Set("Accept-Encoding", "gzip")
	w := httptest.NewRecorder()
	ww := &wrappedResponseWriter{ResponseRecorder: w}
	handler.ServeHTTP(ww, r)
	testutil.Diff(t, http.StatusOK, w.Result().StatusCode)
	rr, _ := gzip.NewReader(bytes.NewReader([]byte(strings.Join(ww.bodies, ""))))
	body, _ := io.ReadAll(rr)
	testutil.Diff(t, "12345", string(body))
}

func TestStream_ChunkedResponse(t *testing.T) {
	server := common.NewAPI()
	err := app.LoadConfigFiles(server, []string{testDataDir + "config-stream.yaml"})
	testutil.DiffError(t, nil, nil, err)

	proxyRef := &kernel.Reference{
		APIVersion: "core/v1",
		Kind:       "ReverseProxyHandler",
	}
	proxy, err := api.ReferTypedObject[http.Handler](server, proxyRef)
	testutil.DiffError(t, nil, nil, err)
	middleRef := &kernel.Reference{
		APIVersion: "app/v1",
		Kind:       "CompressionMiddleware",
	}
	middle, err := api.ReferTypedObject[core.Middleware](server, middleRef)
	testutil.DiffError(t, nil, nil, err)
	handler := middle.Middleware(proxy)

	svr := &http.Server{
		Addr: ":12301",
		Handler: &returnChunkedHandler{
			header: http.Header{
				"Content-Type":      []string{"text/plain"},
				"Transfer-Encoding": []string{"chunked"},
			},
			interval: 10 * time.Millisecond,
			bodies:   []string{"1", "2", "3", "4", "5"},
		},
	}
	go func() { svr.ListenAndServe() }()
	time.Sleep(time.Second) // Wait a little until server starts.
	defer svr.Close()

	r := httptest.NewRequest(http.MethodGet, "http://test.com/test", nil)
	r.Header.Set("Accept-Encoding", "gzip")
	w := httptest.NewRecorder()
	ww := &wrappedResponseWriter{ResponseRecorder: w}
	handler.ServeHTTP(ww, r)
	testutil.Diff(t, http.StatusOK, w.Result().StatusCode)
	rr, _ := gzip.NewReader(bytes.NewReader([]byte(strings.Join(ww.bodies, ""))))
	body, _ := io.ReadAll(rr)
	testutil.Diff(t, "12345", string(body))
}

func TestStream_ReceiveOctetStream(t *testing.T) {
	server := common.NewAPI()
	err := app.LoadConfigFiles(server, []string{testDataDir + "config-stream.yaml"})
	testutil.DiffError(t, nil, nil, err)

	proxyRef := &kernel.Reference{
		APIVersion: "core/v1",
		Kind:       "ReverseProxyHandler",
	}
	proxy, err := api.ReferTypedObject[http.Handler](server, proxyRef)
	testutil.DiffError(t, nil, nil, err)
	middleRef := &kernel.Reference{
		APIVersion: "app/v1",
		Kind:       "CompressionMiddleware",
	}
	middle, err := api.ReferTypedObject[core.Middleware](server, middleRef)
	testutil.DiffError(t, nil, nil, err)
	handler := middle.Middleware(proxy)

	bh := &returnBinaryHandler{}
	svr := &http.Server{
		Addr:    ":12301",
		Handler: bh,
	}
	go func() { svr.ListenAndServe() }()
	time.Sleep(time.Second) // Wait a little until server starts.
	defer svr.Close()

	r := httptest.NewRequest(http.MethodGet, "http://test.com/test", nil)
	r.Header.Set("Accept-Encoding", "gzip")
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)
	rr, _ := gzip.NewReader(w.Body)

	hash := md5.New()
	size, _ := io.Copy(hash, rr)
	md5Hash := hash.Sum(nil)

	testutil.Diff(t, http.StatusOK, w.Result().StatusCode)
	testutil.Diff(t, size, bh.size)
	testutil.Diff(t, hex.EncodeToString(md5Hash), hex.EncodeToString(bh.md5Hash))
}
