//go:build integration
// +build integration

package errorhandler_test

import (
	"cmp"
	"errors"
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
	utilhttp "github.com/aileron-gateway/aileron-gateway/util/http"
)

// testDataDir is the path to the test data.
// This path can be changed by the environmental variable.
var testDataDir = cmp.Or(os.Getenv("TEST_DIR"), "../../../test/") + "integration/errorhandler/"

func testErrorHandler(t *testing.T, eh core.ErrorHandler) {

	t.Helper()

	r := httptest.NewRequest(http.MethodGet, "/test", nil)

	w1 := httptest.NewRecorder()
	eh.ServeHTTPError(w1, r, nil)
	testutil.Diff(t, http.StatusInternalServerError, w1.Result().StatusCode)

	w2 := httptest.NewRecorder()
	eh.ServeHTTPError(w2, r, errors.New("test error"))
	testutil.Diff(t, http.StatusInternalServerError, w2.Result().StatusCode)

	w3 := httptest.NewRecorder()
	eh.ServeHTTPError(w3, r, utilhttp.NewHTTPError(nil, http.StatusUnauthorized))
	testutil.Diff(t, http.StatusUnauthorized, w3.Result().StatusCode)

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
		Kind:       "ErrorHandler",
		Name:       "default",
		Namespace:  "",
	}
	eh, err := api.ReferTypedObject[core.ErrorHandler](server, ref)
	testutil.DiffError(t, nil, nil, err)

	testErrorHandler(t, eh)

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
		Kind:       "ErrorHandler",
		Name:       "testName",
		Namespace:  "testNamespace",
	}
	eh, err := api.ReferTypedObject[core.ErrorHandler](server, ref)
	testutil.DiffError(t, nil, nil, err)

	testErrorHandler(t, eh)

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
		Kind:       "ErrorHandler",
		Name:       "default",
		Namespace:  "testNamespace",
	}
	eh, err := api.ReferTypedObject[core.ErrorHandler](server, ref)
	testutil.DiffError(t, nil, nil, err)

	testErrorHandler(t, eh)

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
		Kind:       "ErrorHandler",
		Name:       "testName",
		Namespace:  "",
	}
	eh, err := api.ReferTypedObject[core.ErrorHandler](server, ref)
	testutil.DiffError(t, nil, nil, err)

	testErrorHandler(t, eh)

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
		Kind:       "ErrorHandler",
		Name:       "default",
		Namespace:  "",
	}
	eh, err := api.ReferTypedObject[core.ErrorHandler](server, ref)
	testutil.DiffError(t, nil, nil, err)

	testErrorHandler(t, eh)

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
		Kind:       "ErrorHandler",
		Name:       "default",
		Namespace:  "",
	}
	eh, err := api.ReferTypedObject[core.ErrorHandler](server, ref)
	testutil.DiffError(t, nil, nil, err)

	testErrorHandler(t, eh)

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
