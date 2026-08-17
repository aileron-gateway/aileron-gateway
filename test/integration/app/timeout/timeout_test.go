// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

//go:build integration

package timeout_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/aileron-gateway/aileron-gateway/apis/kernel"
	"github.com/aileron-gateway/aileron-gateway/cmd/aileron/app"
	"github.com/aileron-gateway/aileron-gateway/core"
	"github.com/aileron-gateway/aileron-gateway/kernel/api"
	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
	"github.com/aileron-gateway/aileron-gateway/test/integration/common"
)

func TestDefaultTimeout(t *testing.T) {
	configs := []string{"./config-default-timeout.yaml"}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "app/v1",
		Kind:       "TimeoutMiddleware",
		Name:       "test",
		Namespace:  "testNamespace",
	}
	m, err := api.ReferTypedObject[core.Middleware](server, ref)
	testutil.DiffError(t, nil, nil, err)

	var sleep int
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		select {
		case <-time.After(time.Millisecond * time.Duration(sleep)):
			w.Write([]byte("timeout not occurred"))
		case <-r.Context().Done():
			return
		}
	})

	h := m.Middleware(handler)

	// default timeout occurred
	sleep = 100
	r1 := httptest.NewRequest(http.MethodGet, "http://timeout.com/default", nil)
	w1 := httptest.NewRecorder()
	h.ServeHTTP(w1, r1)
	b1, _ := io.ReadAll(w1.Result().Body)
	testutil.Diff(t, http.StatusGatewayTimeout, w1.Result().StatusCode)
	testutil.Diff(t, `{"status":504,"statusText":"Gateway Timeout"}`, string(b1))

	// default timeout not occurred
	sleep = 0
	r2 := httptest.NewRequest(http.MethodGet, "http://timeout.com/default", nil)
	w2 := httptest.NewRecorder()
	h.ServeHTTP(w2, r2)
	b2, _ := io.ReadAll(w2.Result().Body)
	testutil.Diff(t, http.StatusOK, w2.Result().StatusCode)
	testutil.Diff(t, "timeout not occurred", string(b2))

}

func TestAppliedApiTimeout(t *testing.T) {
	configs := []string{"./config-api-timeout.yaml"}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "app/v1",
		Kind:       "TimeoutMiddleware",
		Name:       "test",
		Namespace:  "testNamespace",
	}
	m, err := api.ReferTypedObject[core.Middleware](server, ref)
	testutil.DiffError(t, nil, nil, err)

	var sleep int
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		select {
		case <-time.After(time.Millisecond * time.Duration(sleep)):
			w.Write([]byte("timeout not occurred"))
		case <-r.Context().Done():
			return
		}
	})

	h := m.Middleware(handler)

	// API timeout occurred
	sleep = 150
	r1 := httptest.NewRequest(http.MethodGet, "http://timeout.com/applied-api-timeout", nil)
	w1 := httptest.NewRecorder()
	h.ServeHTTP(w1, r1)
	b1, _ := io.ReadAll(w1.Result().Body)
	testutil.Diff(t, http.StatusGatewayTimeout, w1.Result().StatusCode)
	testutil.Diff(t, `{"status":504,"statusText":"Gateway Timeout"}`, string(b1))

	// API timeout not occurred
	sleep = 100
	r2 := httptest.NewRequest(http.MethodGet, "http://timeout.com/applied-api-timeout", nil)
	w2 := httptest.NewRecorder()
	h.ServeHTTP(w2, r2)
	b2, _ := io.ReadAll(w2.Result().Body)
	testutil.Diff(t, http.StatusOK, w2.Result().StatusCode)
	testutil.Diff(t, "timeout not occurred", string(b2))

	// API timeout not occurred but default timeout occurred
	sleep = 100
	r3 := httptest.NewRequest(http.MethodGet, "http://timeout.com/not-applied-api-timeout", nil)
	w3 := httptest.NewRecorder()
	h.ServeHTTP(w3, r3)
	b3, _ := io.ReadAll(w3.Result().Body)
	testutil.Diff(t, http.StatusGatewayTimeout, w3.Result().StatusCode)
	testutil.Diff(t, `{"status":504,"statusText":"Gateway Timeout"}`, string(b3))

	// API timeout and default timeout not occurred
	sleep = 0
	r4 := httptest.NewRequest(http.MethodGet, "http://timeout.com/not-applied-api-timeout", nil)
	w4 := httptest.NewRecorder()
	h.ServeHTTP(w4, r4)
	b4, _ := io.ReadAll(w4.Result().Body)
	testutil.Diff(t, http.StatusOK, w4.Result().StatusCode)
	testutil.Diff(t, "timeout not occurred", string(b4))
}
