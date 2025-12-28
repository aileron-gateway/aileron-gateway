// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

//go:build integration

package httpclient_test

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"regexp"
	"testing"

	"github.com/aileron-gateway/aileron-gateway/apis/kernel"
	"github.com/aileron-gateway/aileron-gateway/cmd/aileron/app"
	"github.com/aileron-gateway/aileron-gateway/core"
	"github.com/aileron-gateway/aileron-gateway/internal/testutil"
	"github.com/aileron-gateway/aileron-gateway/kernel/api"
	"github.com/aileron-gateway/aileron-gateway/test/integration/common"
	"github.com/google/go-cmp/cmp/cmpopts"
)

// testDataDir is the path to the test data.
var testDataDir = "../../../test/integration/httpclient/"

func init() {
	// TEST_DATA_DIR is used in config file.
	os.Setenv("TEST_DATA_DIR", testDataDir+"testdata/")
}

func testRoundTripper(t *testing.T, rt http.RoundTripper) {

	t.Helper()

	svr := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("test", "ok")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("test"))
		},
	))
	defer svr.Close()

	getReq := httptest.NewRequest(http.MethodGet, svr.URL+"/test", nil)
	getRes, err := rt.RoundTrip(getReq)
	getBody, _ := io.ReadAll(getRes.Body)
	testutil.DiffError(t, nil, nil, err)
	testutil.Diff(t, http.StatusOK, getRes.StatusCode)
	testutil.Diff(t, "ok", getRes.Header.Get("test"))
	testutil.Diff(t, "test", string(getBody))

	postReq := httptest.NewRequest(http.MethodPost, svr.URL+"/test", bytes.NewReader([]byte("test")))
	postRes, err := rt.RoundTrip(postReq)
	postBody, _ := io.ReadAll(postRes.Body)
	testutil.DiffError(t, nil, nil, err)
	testutil.Diff(t, http.StatusOK, postRes.StatusCode)
	testutil.Diff(t, "ok", postRes.Header.Get("test"))
	testutil.Diff(t, "test", string(postBody))

	errSvr := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			conn, _, _ := w.(http.Hijacker).Hijack()
			conn.Close()
		},
	))
	defer svr.Close()

	req := httptest.NewRequest(http.MethodGet, errSvr.URL+"/test", nil)
	res, err := rt.RoundTrip(req)
	testutil.DiffError(t, io.EOF, nil, err, cmpopts.EquateErrors())
	testutil.Diff(t, (*http.Response)(nil), res)

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
		Kind:       "HTTPClient",
		Name:       "default",
		Namespace:  "",
	}
	rt, err := api.ReferTypedObject[http.RoundTripper](server, ref)
	testutil.DiffError(t, nil, nil, err)

	testRoundTripper(t, rt)

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
		Kind:       "HTTPClient",
		Name:       "testName",
		Namespace:  "testNamespace",
	}
	rt, err := api.ReferTypedObject[http.RoundTripper](server, ref)
	testutil.DiffError(t, nil, nil, err)

	testRoundTripper(t, rt)

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
		Kind:       "HTTPClient",
		Name:       "default",
		Namespace:  "testNamespace",
	}
	rt, err := api.ReferTypedObject[http.RoundTripper](server, ref)
	testutil.DiffError(t, nil, nil, err)

	testRoundTripper(t, rt)

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
		Kind:       "HTTPClient",
		Name:       "testName",
		Namespace:  "",
	}
	rt, err := api.ReferTypedObject[http.RoundTripper](server, ref)
	testutil.DiffError(t, nil, nil, err)

	testRoundTripper(t, rt)

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
		Kind:       "HTTPClient",
		Name:       "default",
		Namespace:  "",
	}
	_, err = api.ReferTypedObject[http.RoundTripper](server, ref)
	testutil.DiffError(t, nil, nil, err)

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
		Kind:       "HTTPClient",
		Name:       "default",
		Namespace:  "",
	}
	rt, err := api.ReferTypedObject[http.RoundTripper](server, ref)
	testutil.DiffError(t, nil, nil, err)

	testRoundTripper(t, rt)

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
