// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

//go:build integration

package template_test

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aileron-gateway/aileron-gateway/apis/kernel"
	"github.com/aileron-gateway/aileron-gateway/cmd/aileron/app"
	"github.com/aileron-gateway/aileron-gateway/internal/testutil"
	"github.com/aileron-gateway/aileron-gateway/kernel/api"
	"github.com/aileron-gateway/aileron-gateway/test/integration/common"
)

func TestContentsNone(t *testing.T) {

	configs := []string{
		testDataDir + "config-contents-none.yaml",
	}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "core/v1",
		Kind:       "TemplateHandler",
		Name:       "default",
		Namespace:  "",
	}
	h, err := api.ReferTypedObject[http.Handler](server, ref)
	testutil.DiffError(t, nil, nil, err)

	r1 := httptest.NewRequest(http.MethodGet, "/test", nil)
	r1.Header.Set("Accept", "application/json")
	w1 := httptest.NewRecorder()
	h.ServeHTTP(w1, r1)
	b1, _ := io.ReadAll(w1.Body)
	testutil.Diff(t, http.StatusNotAcceptable, w1.Result().StatusCode)
	testutil.Diff(t, `{"status":406,"statusText":"Not Acceptable"}`, string(b1))

	r2 := httptest.NewRequest(http.MethodPost, "/test", bytes.NewReader([]byte("test")))
	r2.Header.Set("Accept", "text/plain")
	w2 := httptest.NewRecorder()
	h.ServeHTTP(w2, r2)
	b2, _ := io.ReadAll(w2.Body)
	testutil.Diff(t, http.StatusNotAcceptable, w2.Result().StatusCode)
	testutil.Diff(t, "status: 406\nstatusText: Not Acceptable\n", string(b2))

	// First content application/json will be used.
	r3 := httptest.NewRequest(http.MethodGet, "/test", nil)
	r3.Header.Set("Accept", "invalid/content")
	w3 := httptest.NewRecorder()
	h.ServeHTTP(w3, r3)
	b3, _ := io.ReadAll(w3.Body)
	testutil.Diff(t, http.StatusNotAcceptable, w3.Result().StatusCode)
	testutil.Diff(t, `{"status":406,"statusText":"Not Acceptable"}`, string(b3))

}

func TestContentsText(t *testing.T) {

	configs := []string{
		testDataDir + "config-contents-text.yaml",
	}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "core/v1",
		Kind:       "TemplateHandler",
		Name:       "default",
		Namespace:  "",
	}
	h, err := api.ReferTypedObject[http.Handler](server, ref)
	testutil.DiffError(t, nil, nil, err)

	r1 := httptest.NewRequest(http.MethodGet, "/test", nil)
	r1.Header.Set("Accept", "application/json")
	w1 := httptest.NewRecorder()
	h.ServeHTTP(w1, r1)
	b1, _ := io.ReadAll(w1.Body)
	testutil.Diff(t, http.StatusFound, w1.Result().StatusCode)
	testutil.Diff(t, "bob", w1.Result().Header.Get("alice"))
	testutil.Diff(t, `{"foo":"bar"}`, string(b1))

	r2 := httptest.NewRequest(http.MethodGet, "/test?key=value", nil)
	r2.Header.Set("Accept", "text/plain")
	w2 := httptest.NewRecorder()
	h.ServeHTTP(w2, r2)
	b2, _ := io.ReadAll(w2.Body)
	testutil.Diff(t, http.StatusForbidden, w2.Result().StatusCode)
	testutil.Diff(t, "bob", w2.Result().Header.Get("alice"))
	testutil.Diff(t, `foo=bar`, string(b2))

	r3 := httptest.NewRequest(http.MethodGet, "/test", nil)
	r3.Header.Set("Accept", "invalid/content")
	w3 := httptest.NewRecorder()
	h.ServeHTTP(w3, r3)
	b3, _ := io.ReadAll(w3.Body)
	testutil.Diff(t, http.StatusNotAcceptable, w3.Result().StatusCode)
	testutil.Diff(t, `{"status":406,"statusText":"Not Acceptable"}`, string(b3))

}

func TestContentsGoText(t *testing.T) {

	configs := []string{
		testDataDir + "config-contents-go-text.yaml",
	}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "core/v1",
		Kind:       "TemplateHandler",
		Name:       "default",
		Namespace:  "",
	}
	h, err := api.ReferTypedObject[http.Handler](server, ref)
	testutil.DiffError(t, nil, nil, err)

	r1 := httptest.NewRequest(http.MethodGet, "/test?key=value", nil)
	r1.Header.Set("Accept", "application/json")
	w1 := httptest.NewRecorder()
	h.ServeHTTP(w1, r1)
	b1, _ := io.ReadAll(w1.Body)
	testutil.Diff(t, http.StatusFound, w1.Result().StatusCode)
	testutil.Diff(t, "bob", w1.Result().Header.Get("alice"))
	testutil.Diff(t, `{"method":"GET"}`, string(b1))

	r2 := httptest.NewRequest(http.MethodGet, "/test?key=value", nil)
	r2.Header.Set("Accept", "text/plain")
	w2 := httptest.NewRecorder()
	h.ServeHTTP(w2, r2)
	b2, _ := io.ReadAll(w2.Body)
	testutil.Diff(t, http.StatusForbidden, w2.Result().StatusCode)
	testutil.Diff(t, "bob", w2.Result().Header.Get("alice"))
	testutil.Diff(t, `method=GET`, string(b2))

	r3 := httptest.NewRequest(http.MethodGet, "/test?key=value", nil)
	r3.Header.Set("Accept", "invalid/content")
	w3 := httptest.NewRecorder()
	h.ServeHTTP(w3, r3)
	b3, _ := io.ReadAll(w3.Body)
	testutil.Diff(t, http.StatusNotAcceptable, w3.Result().StatusCode)
	testutil.Diff(t, `{"status":406,"statusText":"Not Acceptable"}`, string(b3))

}

func TestContentsGoHTML(t *testing.T) {

	configs := []string{
		testDataDir + "config-contents-go-html.yaml",
	}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "core/v1",
		Kind:       "TemplateHandler",
		Name:       "default",
		Namespace:  "",
	}
	h, err := api.ReferTypedObject[http.Handler](server, ref)
	testutil.DiffError(t, nil, nil, err)

	r1 := httptest.NewRequest(http.MethodGet, "/test?key=value", nil)
	r1.Header.Set("Accept", "application/json")
	w1 := httptest.NewRecorder()
	h.ServeHTTP(w1, r1)
	b1, _ := io.ReadAll(w1.Body)
	testutil.Diff(t, http.StatusNotAcceptable, w1.Result().StatusCode)
	testutil.Diff(t, `{"status":406,"statusText":"Not Acceptable"}`, string(b1))

	r2 := httptest.NewRequest(http.MethodGet, "/test?key=value", nil)
	r2.Header.Set("Accept", "text/html")
	w2 := httptest.NewRecorder()
	h.ServeHTTP(w2, r2)
	b2, _ := io.ReadAll(w2.Body)
	testutil.Diff(t, http.StatusFound, w2.Result().StatusCode)
	testutil.Diff(t, "bob", w2.Result().Header.Get("alice"))
	testutil.Diff(t, `<a href="/test">/test</a>`, string(b2))

}
