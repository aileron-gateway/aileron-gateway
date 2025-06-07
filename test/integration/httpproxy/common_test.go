// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

//go:build integration

package httpproxy_test

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/aileron-gateway/aileron-gateway/apis/kernel"
	"github.com/aileron-gateway/aileron-gateway/cmd/aileron/app"
	"github.com/aileron-gateway/aileron-gateway/core"
	"github.com/aileron-gateway/aileron-gateway/kernel/api"
	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
	"github.com/aileron-gateway/aileron-gateway/test/integration/common"
)

// testDataDir is the path to the test data.
var testDataDir = "../../../test/integration/httpproxy/"

func testHandler(t *testing.T, h http.Handler) {

	t.Helper()

	r1 := httptest.NewRequest(http.MethodGet, "http://localhost/test", nil)
	w1 := httptest.NewRecorder()
	h.ServeHTTP(w1, r1)
	body1, _ := io.ReadAll(w1.Body)
	testutil.Diff(t, http.StatusNotFound, w1.Result().StatusCode)
	testutil.Diff(t, []byte(`{"status":404,"statusText":"Not Found"}`), body1)

	r2 := httptest.NewRequest(http.MethodPost, "http://localhost/test", bytes.NewReader([]byte("test")))
	w2 := httptest.NewRecorder()
	h.ServeHTTP(w2, r2)
	body2, _ := io.ReadAll(w2.Body)
	testutil.Diff(t, http.StatusNotFound, w2.Result().StatusCode)
	testutil.Diff(t, []byte(`{"status":404,"statusText":"Not Found"}`), body2)

}

func TestMinimalWithoutMetadata(t *testing.T) {

	configs := []string{
		testDataDir + "config-minimal-without-metadata.yaml",
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
	h, err := api.ReferTypedObject[http.Handler](server, ref)
	testutil.DiffError(t, nil, nil, err)

	testHandler(t, h)

}

func TestMinimalWithMetadata(t *testing.T) {

	configs := []string{
		testDataDir + "config-minimal-with-metadata.yaml",
	}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "core/v1",
		Kind:       "ReverseProxyHandler",
		Name:       "testName",
		Namespace:  "testNamespace",
	}
	h, err := api.ReferTypedObject[http.Handler](server, ref)
	testutil.DiffError(t, nil, nil, err)

	testHandler(t, h)

}

func TestEmptyName(t *testing.T) {

	configs := []string{
		testDataDir + "config-empty-name.yaml",
	}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "core/v1",
		Kind:       "ReverseProxyHandler",
		Name:       "default",
		Namespace:  "testNamespace",
	}
	h, err := api.ReferTypedObject[http.Handler](server, ref)
	testutil.DiffError(t, nil, nil, err)

	testHandler(t, h)

}

func TestEmptyNamespace(t *testing.T) {

	configs := []string{
		testDataDir + "config-empty-namespace.yaml",
	}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "core/v1",
		Kind:       "ReverseProxyHandler",
		Name:       "testName",
		Namespace:  "",
	}
	h, err := api.ReferTypedObject[http.Handler](server, ref)
	testutil.DiffError(t, nil, nil, err)

	testHandler(t, h)

}

func TestEmptyNameNamespace(t *testing.T) {

	configs := []string{
		testDataDir + "config-empty-name-namespace.yaml",
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
	h, err := api.ReferTypedObject[http.Handler](server, ref)
	testutil.DiffError(t, nil, nil, err)

	testHandler(t, h)

}

func TestEmptySpec(t *testing.T) {

	configs := []string{
		testDataDir + "config-empty-spec.yaml",
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
	h, err := api.ReferTypedObject[http.Handler](server, ref)
	testutil.DiffError(t, nil, nil, err)

	testHandler(t, h)

}

func TestInvalidSpec(t *testing.T) {

	configs := []string{
		testDataDir + "config-invalid-spec.yaml",
	}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	errPattern := regexp.MustCompile(core.ErrPrefix + `failed to load configs.`)
	testutil.DiffError(t, app.ErrAppMainLoadConfigs, errPattern, err)

}
