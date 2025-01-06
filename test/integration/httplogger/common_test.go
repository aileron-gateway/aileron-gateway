//go:build integration
// +build integration

package httplogger_test

import (
	"bytes"
	"cmp"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"regexp"
	"testing"

	"github.com/aileron-gateway/aileron-gateway/apis/kernel"
	"github.com/aileron-gateway/aileron-gateway/app"
	"github.com/aileron-gateway/aileron-gateway/core"
	"github.com/aileron-gateway/aileron-gateway/kernel/api"
	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
	"github.com/aileron-gateway/aileron-gateway/test/integration/common"
)

// testDataDir is the path to the test data.
// This path can be changed by the environmental variable.
var testDataDir = cmp.Or(os.Getenv("TEST_DIR"), "../../../test/") + "integration/httplogger/"

func testHTTPLoggerMiddleware(t *testing.T, lg core.Middleware) {

	t.Helper()

	h := lg.Middleware(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("test", "ok")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("test"))
		}),
	)

	r1 := httptest.NewRequest(http.MethodGet, "/test", nil)
	w1 := httptest.NewRecorder()
	h.ServeHTTP(w1, r1)
	body1, _ := io.ReadAll(w1.Body)
	testutil.Diff(t, http.StatusOK, w1.Result().StatusCode)
	testutil.Diff(t, "ok", w1.Header().Get("test"))
	testutil.Diff(t, []byte("test"), body1)

	r2 := httptest.NewRequest(http.MethodPost, "/test", nil)
	w2 := httptest.NewRecorder()
	h.ServeHTTP(w2, r2)
	body2, _ := io.ReadAll(w2.Body)
	testutil.Diff(t, http.StatusOK, w2.Result().StatusCode)
	testutil.Diff(t, []byte("test"), body2)

}

func testHTTPLoggerTripperware(t *testing.T, lg core.Tripperware) {

	t.Helper()

	rt := lg.Tripperware(core.RoundTripperFunc(
		func(r *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Header:     http.Header{"Test": {"ok"}},
				Body:       io.NopCloser(bytes.NewReader([]byte("test"))),
			}, nil
		},
	))

	r1 := httptest.NewRequest(http.MethodGet, "/test", nil)
	w1, _ := rt.RoundTrip(r1)
	body1, _ := io.ReadAll(w1.Body)
	testutil.Diff(t, http.StatusOK, w1.StatusCode)
	testutil.Diff(t, "ok", w1.Header.Get("test"))
	testutil.Diff(t, []byte("test"), body1)

	r2 := httptest.NewRequest(http.MethodPost, "/test", bytes.NewReader([]byte("test")))
	w2, _ := rt.RoundTrip(r1)
	rt.RoundTrip(r2)
	body2, _ := io.ReadAll(w2.Body)
	testutil.Diff(t, http.StatusOK, w2.StatusCode)
	testutil.Diff(t, "ok", w2.Header.Get("test"))
	testutil.Diff(t, []byte("test"), body2)

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
		Kind:       "HTTPLogger",
		Name:       "default",
		Namespace:  "",
	}

	m, err := api.ReferTypedObject[core.Middleware](server, ref)
	testutil.DiffError(t, nil, nil, err)
	testHTTPLoggerMiddleware(t, m)

	tr, err := api.ReferTypedObject[core.Tripperware](server, ref)
	testutil.DiffError(t, nil, nil, err)
	testHTTPLoggerTripperware(t, tr)

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
		Kind:       "HTTPLogger",
		Name:       "testName",
		Namespace:  "testNamespace",
	}

	m, err := api.ReferTypedObject[core.Middleware](server, ref)
	testutil.DiffError(t, nil, nil, err)
	testHTTPLoggerMiddleware(t, m)

	tr, err := api.ReferTypedObject[core.Tripperware](server, ref)
	testutil.DiffError(t, nil, nil, err)
	testHTTPLoggerTripperware(t, tr)

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
		Kind:       "HTTPLogger",
		Name:       "default",
		Namespace:  "testNamespace",
	}

	m, err := api.ReferTypedObject[core.Middleware](server, ref)
	testutil.DiffError(t, nil, nil, err)
	testHTTPLoggerMiddleware(t, m)

	tr, err := api.ReferTypedObject[core.Tripperware](server, ref)
	testutil.DiffError(t, nil, nil, err)
	testHTTPLoggerTripperware(t, tr)

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
		Kind:       "HTTPLogger",
		Name:       "testName",
		Namespace:  "",
	}

	m, err := api.ReferTypedObject[core.Middleware](server, ref)
	testutil.DiffError(t, nil, nil, err)
	testHTTPLoggerMiddleware(t, m)

	tr, err := api.ReferTypedObject[core.Tripperware](server, ref)
	testutil.DiffError(t, nil, nil, err)
	testHTTPLoggerTripperware(t, tr)

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
		Kind:       "HTTPLogger",
		Name:       "default",
		Namespace:  "",
	}

	m, err := api.ReferTypedObject[core.Middleware](server, ref)
	testutil.DiffError(t, nil, nil, err)
	testHTTPLoggerMiddleware(t, m)

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
		Kind:       "HTTPLogger",
		Name:       "default",
		Namespace:  "",
	}

	m, err := api.ReferTypedObject[core.Middleware](server, ref)
	testutil.DiffError(t, nil, nil, err)
	testHTTPLoggerMiddleware(t, m)

	tr, err := api.ReferTypedObject[core.Tripperware](server, ref)
	testutil.DiffError(t, nil, nil, err)
	testHTTPLoggerTripperware(t, tr)

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
