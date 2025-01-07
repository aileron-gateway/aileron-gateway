//go:build integration
// +build integration

package static_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/aileron-gateway/aileron-gateway/apis/kernel"
	"github.com/aileron-gateway/aileron-gateway/cmd/aileron/app"
	"github.com/aileron-gateway/aileron-gateway/kernel/api"
	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
	"github.com/aileron-gateway/aileron-gateway/test/integration/common"
)

func TestEnableListing(t *testing.T) {

	configs := []string{
		testDataDir + "config-enable-listing.yaml",
	}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "core/v1",
		Kind:       "StaticFileHandler",
		Name:       "default",
		Namespace:  "",
	}
	h, err := api.ReferTypedObject[http.Handler](server, ref)
	testutil.DiffError(t, nil, nil, err)

	r := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	b, _ := io.ReadAll(w.Body)
	testutil.Diff(t, http.StatusOK, w.Result().StatusCode)
	testutil.Diff(t, "bar", w.Result().Header.Get("foo"))
	testutil.Diff(t, true, strings.Contains(string(b), "foo.json"))
	testutil.Diff(t, true, strings.Contains(string(b), "foo.txt"))

	rr := httptest.NewRequest(http.MethodGet, "/not-found", nil)
	ww := httptest.NewRecorder()
	h.ServeHTTP(ww, rr)
	testutil.Diff(t, http.StatusNotFound, ww.Result().StatusCode)
	testutil.Diff(t, "bar", ww.Result().Header.Get("foo"))

}

func TestStripPrefix(t *testing.T) {

	configs := []string{
		testDataDir + "config-strip-prefix.yaml",
	}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "core/v1",
		Kind:       "StaticFileHandler",
		Name:       "default",
		Namespace:  "",
	}
	h, err := api.ReferTypedObject[http.Handler](server, ref)
	testutil.DiffError(t, nil, nil, err)

	r := httptest.NewRequest(http.MethodGet, "/prefix/foo.txt", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	b, _ := io.ReadAll(w.Body)
	testutil.Diff(t, http.StatusOK, w.Result().StatusCode)
	testutil.Diff(t, "bar", w.Result().Header.Get("foo"))
	testutil.Diff(t, `foo=bar`, string(b))

	rr := httptest.NewRequest(http.MethodGet, "/not-found", nil)
	ww := httptest.NewRecorder()
	h.ServeHTTP(ww, rr)
	testutil.Diff(t, http.StatusNotFound, ww.Result().StatusCode)
	testutil.Diff(t, "bar", ww.Result().Header.Get("foo"))

}
