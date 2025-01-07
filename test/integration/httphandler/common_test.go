//go:build integration
// +build integration

package httphandler_test

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
	"github.com/aileron-gateway/aileron-gateway/cmd/aileron/app"
	"github.com/aileron-gateway/aileron-gateway/core"
	"github.com/aileron-gateway/aileron-gateway/kernel/api"
	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
	"github.com/aileron-gateway/aileron-gateway/test/integration/common"
)

// testDataDir is the path to the test data.
// This path can be changed by the environmental variable.
var testDataDir = cmp.Or(os.Getenv("TEST_DIR"), "../../../test/") + "integration/httphandler/"

type mocHandler struct {
}

func (h *mocHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("test", "ok")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("test"))
}

func testHandler(t *testing.T, h http.Handler) {

	t.Helper()

	r1 := httptest.NewRequest(http.MethodGet, "/test", nil)
	w1 := httptest.NewRecorder()
	h.ServeHTTP(w1, r1)
	body1, _ := io.ReadAll(w1.Body)
	testutil.Diff(t, http.StatusOK, w1.Result().StatusCode)
	testutil.Diff(t, "ok", w1.Header().Get("test"))
	testutil.Diff(t, []byte("test"), body1)

	r2 := httptest.NewRequest(http.MethodPost, "/test", bytes.NewReader([]byte("test")))
	w2 := httptest.NewRecorder()
	h.ServeHTTP(w2, r2)
	body2, _ := io.ReadAll(w2.Body)
	testutil.Diff(t, http.StatusOK, w2.Result().StatusCode)
	testutil.Diff(t, "ok", w2.Header().Get("test"))
	testutil.Diff(t, []byte("test"), body2)

}

func TestMinimalWithoutMetadata(t *testing.T) {

	configs := []string{
		testDataDir + "config-minimal-without-metadata.yaml",
	}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	mocRef := &kernel.Reference{
		APIVersion: "container/v1",
		Kind:       "TestHandler",
		Name:       "test",
		Namespace:  "",
	}
	common.PostTestResource(server, mocRef, &mocHandler{})

	ref := &kernel.Reference{
		APIVersion: "core/v1",
		Kind:       "HTTPHandler",
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

	mocRef := &kernel.Reference{
		APIVersion: "container/v1",
		Kind:       "TestHandler",
		Name:       "test",
		Namespace:  "",
	}
	common.PostTestResource(server, mocRef, &mocHandler{})

	ref := &kernel.Reference{
		APIVersion: "core/v1",
		Kind:       "HTTPHandler",
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

	mocRef := &kernel.Reference{
		APIVersion: "container/v1",
		Kind:       "TestHandler",
		Name:       "test",
		Namespace:  "",
	}
	common.PostTestResource(server, mocRef, &mocHandler{})

	ref := &kernel.Reference{
		APIVersion: "core/v1",
		Kind:       "HTTPHandler",
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

	mocRef := &kernel.Reference{
		APIVersion: "container/v1",
		Kind:       "TestHandler",
		Name:       "test",
		Namespace:  "",
	}
	common.PostTestResource(server, mocRef, &mocHandler{})

	ref := &kernel.Reference{
		APIVersion: "core/v1",
		Kind:       "HTTPHandler",
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

	mocRef := &kernel.Reference{
		APIVersion: "container/v1",
		Kind:       "TestHandler",
		Name:       "test",
		Namespace:  "",
	}
	common.PostTestResource(server, mocRef, &mocHandler{})

	ref := &kernel.Reference{
		APIVersion: "core/v1",
		Kind:       "HTTPHandler",
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

	mocRef := &kernel.Reference{
		APIVersion: "container/v1",
		Kind:       "TestHandler",
		Name:       "test",
		Namespace:  "",
	}
	common.PostTestResource(server, mocRef, &mocHandler{})

	ref := &kernel.Reference{
		APIVersion: "core/v1",
		Kind:       "HTTPHandler",
		Name:       "default",
		Namespace:  "",
	}
	_, err = api.ReferTypedObject[http.Handler](server, ref)
	errPattern := regexp.MustCompile(core.ErrPrefix + `failed to create`)
	testutil.DiffError(t, core.ErrCoreGenCreateObject, errPattern, err)

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
