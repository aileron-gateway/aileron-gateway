// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

//go:build integration

package httpserver_test

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"os"
	"regexp"
	"testing"

	"github.com/aileron-gateway/aileron-gateway/apis/kernel"
	"github.com/aileron-gateway/aileron-gateway/cmd/aileron/app"
	"github.com/aileron-gateway/aileron-gateway/core"
	"github.com/aileron-gateway/aileron-gateway/internal/testutil"
	"github.com/aileron-gateway/aileron-gateway/kernel/api"
	"github.com/aileron-gateway/aileron-gateway/test/integration/common"
)

// testDataDir is the path to the test data.
var testDataDir = "../../../test/integration/httpserver/"

func init() {
	// TEST_DATA_DIR is used in config file.
	os.Setenv("TEST_DATA_DIR", testDataDir+"testdata/")
}

func testServer(t *testing.T, svr core.Runner) {

	t.Helper()

	ctx, cancel := context.WithCancel(context.Background())

	closed := make(chan struct{})
	go func() {
		err := svr.Run(ctx)
		testutil.DiffError(t, nil, nil, err)
		close(closed)
	}()

	r1, _ := http.NewRequest(http.MethodGet, "http://localhost:8080/test", nil)
	w1, err := http.DefaultClient.Do(r1)
	testutil.DiffError(t, nil, nil, err)
	body1, _ := io.ReadAll(w1.Body)
	testutil.Diff(t, http.StatusNotFound, w1.StatusCode)
	testutil.Diff(t, []byte(`{"status":404,"statusText":"Not Found"}`), body1)

	r2, _ := http.NewRequest(http.MethodPost, "http://localhost:8080/test", bytes.NewReader([]byte("test")))
	w2, err := http.DefaultClient.Do(r2)
	testutil.DiffError(t, nil, nil, err)
	body2, _ := io.ReadAll(w2.Body)
	testutil.Diff(t, http.StatusNotFound, w2.StatusCode)
	testutil.Diff(t, []byte(`{"status":404,"statusText":"Not Found"}`), body2)

	cancel()
	<-closed

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
		Kind:       "HTTPServer",
		Name:       "default",
		Namespace:  "",
	}
	svr, err := api.ReferTypedObject[core.Runner](server, ref)
	testutil.DiffError(t, nil, nil, err)

	testServer(t, svr)

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
		Kind:       "HTTPServer",
		Name:       "testName",
		Namespace:  "testNamespace",
	}
	svr, err := api.ReferTypedObject[core.Runner](server, ref)
	testutil.DiffError(t, nil, nil, err)

	testServer(t, svr)

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
		Kind:       "HTTPServer",
		Name:       "default",
		Namespace:  "testNamespace",
	}
	svr, err := api.ReferTypedObject[core.Runner](server, ref)
	testutil.DiffError(t, nil, nil, err)

	testServer(t, svr)

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
		Kind:       "HTTPServer",
		Name:       "testName",
		Namespace:  "",
	}
	svr, err := api.ReferTypedObject[core.Runner](server, ref)
	testutil.DiffError(t, nil, nil, err)

	testServer(t, svr)

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
		Kind:       "HTTPServer",
		Name:       "default",
		Namespace:  "",
	}
	svr, err := api.ReferTypedObject[core.Runner](server, ref)
	testutil.DiffError(t, nil, nil, err)

	testServer(t, svr)

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
		Kind:       "HTTPServer",
		Name:       "default",
		Namespace:  "",
	}
	svr, err := api.ReferTypedObject[core.Runner](server, ref)
	testutil.DiffError(t, nil, nil, err)

	testServer(t, svr)

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
