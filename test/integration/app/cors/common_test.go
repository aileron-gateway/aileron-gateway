//go:build integration
// +build integration

package cors_test

import (
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

func check(t *testing.T, m core.Middleware) {
	t.Helper()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})
	h := m.Middleware(handler)

	// Preflight request without Origin.
	r1 := httptest.NewRequest(http.MethodOptions, "http://cors-test.com/test", nil)
	r1.Header.Set("Access-Control-Request-Method", http.MethodGet)
	r1.Header.Set("Access-Control-Request-Headers", "Foo,Bar")
	w1 := httptest.NewRecorder()
	h.ServeHTTP(w1, r1)
	b1, _ := io.ReadAll(w1.Result().Body)
	testutil.Diff(t, "Content-Type,X-Requested-With", w1.Header().Get("Access-Control-Allow-Headers"))
	testutil.Diff(t, "GET,OPTIONS,POST", w1.Header().Get("Access-Control-Allow-Methods"))
	testutil.Diff(t, http.StatusForbidden, w1.Result().StatusCode)
	testutil.Diff(t, "", string(b1))

	// Preflight request with Origin.
	r2 := httptest.NewRequest(http.MethodOptions, "http://cors-test.com/test", nil)
	r2.Header.Set("Access-Control-Request-Method", http.MethodGet)
	r2.Header.Set("Access-Control-Request-Headers", "Foo,Bar")
	r2.Header.Set("Origin", "http://cors-test.com")
	w2 := httptest.NewRecorder()
	h.ServeHTTP(w2, r2)
	b2, _ := io.ReadAll(w2.Result().Body)
	testutil.Diff(t, "Content-Type,X-Requested-With", w2.Header().Get("Access-Control-Allow-Headers"))
	testutil.Diff(t, "GET,OPTIONS,POST", w2.Header().Get("Access-Control-Allow-Methods"))
	testutil.Diff(t, "*", w2.Header().Get("Access-Control-Allow-Origin"))
	testutil.Diff(t, http.StatusOK, w2.Result().StatusCode)
	testutil.Diff(t, "", string(b2))

	// Actual request without Origin.
	r3 := httptest.NewRequest(http.MethodGet, "http://cors-test.com/test", nil)
	w3 := httptest.NewRecorder()
	h.ServeHTTP(w3, r3)
	b3, _ := io.ReadAll(w3.Result().Body)
	testutil.Diff(t, "Content-Type,X-Requested-With", w3.Header().Get("Access-Control-Allow-Headers"))
	testutil.Diff(t, "GET,OPTIONS,POST", w3.Header().Get("Access-Control-Allow-Methods"))
	testutil.Diff(t, http.StatusForbidden, w3.Result().StatusCode)
	testutil.Diff(t, `{"status":403,"statusText":"Forbidden"}`, string(b3))

	// Actual request with Origin.
	r4 := httptest.NewRequest(http.MethodGet, "http://cors-test.com/test", nil)
	r4.Header.Set("Origin", "http://cors-test.com")
	w4 := httptest.NewRecorder()
	h.ServeHTTP(w4, r4)
	b4, _ := io.ReadAll(w4.Result().Body)
	testutil.Diff(t, "Content-Type,X-Requested-With", w4.Header().Get("Access-Control-Allow-Headers"))
	testutil.Diff(t, "GET,OPTIONS,POST", w4.Header().Get("Access-Control-Allow-Methods"))
	testutil.Diff(t, "*", w4.Header().Get("Access-Control-Allow-Origin"))
	testutil.Diff(t, http.StatusOK, w4.Result().StatusCode)
	testutil.Diff(t, "ok", string(b4))

	// Disallowed method.
	r5 := httptest.NewRequest(http.MethodDelete, "http://cors-test.com/test", nil)
	r5.Header.Set("Origin", "http://cors-test.com")
	w5 := httptest.NewRecorder()
	h.ServeHTTP(w5, r5)
	b5, _ := io.ReadAll(w5.Result().Body)
	testutil.Diff(t, "Content-Type,X-Requested-With", w5.Header().Get("Access-Control-Allow-Headers"))
	testutil.Diff(t, "GET,OPTIONS,POST", w5.Header().Get("Access-Control-Allow-Methods"))
	testutil.Diff(t, "*", w5.Header().Get("Access-Control-Allow-Origin"))
	testutil.Diff(t, http.StatusForbidden, w5.Result().StatusCode)
	testutil.Diff(t, `{"status":403,"statusText":"Forbidden"}`, string(b5))
}

func TestMinimalWithoutMetadata(t *testing.T) {
	configs := []string{"./config-minimal-without-metadata.yaml"}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "app/v1",
		Kind:       "CORSMiddleware",
	}
	m, err := api.ReferTypedObject[core.Middleware](server, ref)
	testutil.DiffError(t, nil, nil, err)

	check(t, m)
}

func TestMinimalWithMetadata(t *testing.T) {
	configs := []string{"./config-minimal-with-metadata.yaml"}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "app/v1",
		Kind:       "CORSMiddleware",
		Name:       "testName",
		Namespace:  "testNamespace",
	}
	m, err := api.ReferTypedObject[core.Middleware](server, ref)
	testutil.DiffError(t, nil, nil, err)

	check(t, m)
}

func TestEmptyName(t *testing.T) {
	configs := []string{"./config-empty-name.yaml"}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "app/v1",
		Kind:       "CORSMiddleware",
		Name:       "",
		Namespace:  "testNamespace",
	}
	m, err := api.ReferTypedObject[core.Middleware](server, ref)
	testutil.DiffError(t, nil, nil, err)

	check(t, m)
}

func TestEmptyNamespace(t *testing.T) {
	configs := []string{"./config-empty-namespace.yaml"}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "app/v1",
		Kind:       "CORSMiddleware",
		Name:       "testName",
		Namespace:  "",
	}
	m, err := api.ReferTypedObject[core.Middleware](server, ref)
	testutil.DiffError(t, nil, nil, err)

	check(t, m)
}

func TestEmptyNameNamespace(t *testing.T) {
	configs := []string{"./config-empty-name-namespace.yaml"}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "app/v1",
		Kind:       "CORSMiddleware",
		Name:       "",
		Namespace:  "",
	}
	m, err := api.ReferTypedObject[core.Middleware](server, ref)
	testutil.DiffError(t, nil, nil, err)

	check(t, m)
}

func TestEmptySpec(t *testing.T) {
	configs := []string{"./config-empty-spec.yaml"}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "app/v1",
		Kind:       "CORSMiddleware",
	}
	m, err := api.ReferTypedObject[core.Middleware](server, ref)
	testutil.DiffError(t, nil, nil, err)

	check(t, m)
}

func TestInvalidSpec(t *testing.T) {
	configs := []string{"./config-invalid-spec.yaml"}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	errPattern := regexp.MustCompile(core.ErrPrefix + `failed to load configs. config-invalid-spec.yaml`)
	testutil.DiffError(t, app.ErrAppMainLoadConfigs, errPattern, err)
}
