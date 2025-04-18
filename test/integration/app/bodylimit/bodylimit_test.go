// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

//go:build integration

package bodylimit_test

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aileron-gateway/aileron-gateway/apis/kernel"
	"github.com/aileron-gateway/aileron-gateway/cmd/aileron/app"
	"github.com/aileron-gateway/aileron-gateway/core"
	"github.com/aileron-gateway/aileron-gateway/kernel/api"
	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
	"github.com/aileron-gateway/aileron-gateway/test/integration/common"
)

func TestMaxSize(t *testing.T) {

	configs := []string{"./config-max-size.yaml"}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "app/v1",
		Kind:       "BodyLimitMiddleware",
	}
	m, err := api.ReferTypedObject[core.Middleware](server, ref)
	testutil.DiffError(t, nil, nil, err)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				if rec != http.ErrAbortHandler {
					panic(rec)
				}
			}
		}()
		_, _ = io.ReadAll(r.Body)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("test"))
	})
	h := m.Middleware(handler)

	r1 := httptest.NewRequest(http.MethodPost, "http://test.com/", bytes.NewReader([]byte("12345")))
	w1 := httptest.NewRecorder()
	h.ServeHTTP(w1, r1)
	b1, _ := io.ReadAll(w1.Result().Body)
	testutil.Diff(t, http.StatusOK, w1.Result().StatusCode)
	testutil.Diff(t, "test", string(b1))

	r2 := httptest.NewRequest(http.MethodPost, "http://test.com/", bytes.NewReader([]byte("12345678901")))
	w2 := httptest.NewRecorder()
	h.ServeHTTP(w2, r2)
	b2, _ := io.ReadAll(w2.Result().Body)
	testutil.Diff(t, http.StatusRequestEntityTooLarge, w2.Result().StatusCode)
	testutil.Diff(t, `{"status":413,"statusText":"Request Entity Too Large"}`, string(b2))

	r3 := httptest.NewRequest(http.MethodPost, "http://test.com/", nil)
	r3.ContentLength = -1
	w3 := httptest.NewRecorder()
	h.ServeHTTP(w3, r3)
	b3, _ := io.ReadAll(w3.Result().Body)
	testutil.Diff(t, http.StatusOK, w3.Result().StatusCode)
	testutil.Diff(t, "test", string(b3))

	r4 := httptest.NewRequest(http.MethodPost, "http://test.com/", bytes.NewReader([]byte("12345678901")))
	r4.ContentLength = -1
	w4 := httptest.NewRecorder()
	h.ServeHTTP(w4, r4)
	b4, _ := io.ReadAll(w4.Result().Body)
	testutil.Diff(t, http.StatusRequestEntityTooLarge, w4.Result().StatusCode)
	testutil.Diff(t, `{"status":413,"statusText":"Request Entity Too Large"}`, string(b4))
}

func TestMaxSizeWithLimit(t *testing.T) {

	configs := []string{"./config-max-size-with-limit.yaml"}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "app/v1",
		Kind:       "BodyLimitMiddleware",
	}
	m, err := api.ReferTypedObject[core.Middleware](server, ref)
	testutil.DiffError(t, nil, nil, err)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				if rec != http.ErrAbortHandler {
					panic(rec)
				}
			}
		}()
		_, _ = io.ReadAll(r.Body)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("test"))
	})
	h := m.Middleware(handler)

	r1 := httptest.NewRequest(http.MethodPost, "http://test.com/", bytes.NewReader([]byte("1234567")))
	w1 := httptest.NewRecorder()
	h.ServeHTTP(w1, r1)
	b1, _ := io.ReadAll(w1.Result().Body)
	testutil.Diff(t, http.StatusOK, w1.Result().StatusCode)
	testutil.Diff(t, "test", string(b1))

	r2 := httptest.NewRequest(http.MethodPost, "http://test.com/", bytes.NewReader([]byte("12345678901")))
	w2 := httptest.NewRecorder()
	h.ServeHTTP(w2, r2)
	b2, _ := io.ReadAll(w2.Result().Body)
	testutil.Diff(t, http.StatusRequestEntityTooLarge, w2.Result().StatusCode)
	testutil.Diff(t, `{"status":413,"statusText":"Request Entity Too Large"}`, string(b2))

	r3 := httptest.NewRequest(http.MethodPost, "http://test.com/", nil)
	r3.ContentLength = -1
	w3 := httptest.NewRecorder()
	h.ServeHTTP(w3, r3)
	b3, _ := io.ReadAll(w3.Result().Body)
	testutil.Diff(t, http.StatusOK, w3.Result().StatusCode)
	testutil.Diff(t, "test", string(b3))

	r4 := httptest.NewRequest(http.MethodPost, "http://test.com/", bytes.NewReader([]byte("12345678901")))
	r4.ContentLength = -1
	w4 := httptest.NewRecorder()
	h.ServeHTTP(w4, r4)
	b4, _ := io.ReadAll(w4.Result().Body)
	testutil.Diff(t, http.StatusRequestEntityTooLarge, w4.Result().StatusCode)
	testutil.Diff(t, `{"status":413,"statusText":"Request Entity Too Large"}`, string(b4))
}
