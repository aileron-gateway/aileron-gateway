//go:build integration

package casbin_test

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

func testCasbinAuthzMiddleware(t *testing.T, m core.Middleware) {

	t.Helper()

	handler := m.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	}))

	h := m.Middleware(handler)

	r1 := httptest.NewRequest(http.MethodGet, "http://casbin.com/denied", nil)
	r1.Header = http.Header{"Authorization": {"Bearer denied-token"}}
	w1 := httptest.NewRecorder()
	h.ServeHTTP(w1, r1)
	b1, _ := io.ReadAll(w1.Result().Body)
	testutil.Diff(t, http.StatusForbidden, w1.Result().StatusCode)
	testutil.Diff(t, `{"status":403,"statusText":"Forbidden"}`, string(b1))

}

func TestMinimalWithoutMetadata(t *testing.T) {

	configs := []string{"./config-minimal-without-metadata.yaml"}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "app/v1",
		Kind:       "CasbinAuthzMiddleware",
		Name:       "default",
		Namespace:  "",
	}
	eh, err := api.ReferTypedObject[core.Middleware](server, ref)
	testutil.DiffError(t, nil, nil, err)

	testCasbinAuthzMiddleware(t, eh)

}

func TestMinimalWithMetadata(t *testing.T) {

	configs := []string{"./config-minimal-with-metadata.yaml"}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "app/v1",
		Kind:       "CasbinAuthzMiddleware",
		Name:       "testName",
		Namespace:  "testNamespace",
	}
	eh, err := api.ReferTypedObject[core.Middleware](server, ref)
	testutil.DiffError(t, nil, nil, err)

	testCasbinAuthzMiddleware(t, eh)

}

func TestEmptyName(t *testing.T) {

	configs := []string{"./config-empty-name.yaml"}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "app/v1",
		Kind:       "CasbinAuthzMiddleware",
		Name:       "default",
		Namespace:  "testNamespace",
	}
	eh, err := api.ReferTypedObject[core.Middleware](server, ref)
	testutil.DiffError(t, nil, nil, err)

	testCasbinAuthzMiddleware(t, eh)

}

func TestEmptyNamespace(t *testing.T) {

	configs := []string{"./config-empty-namespace.yaml"}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "app/v1",
		Kind:       "CasbinAuthzMiddleware",
		Name:       "testName",
		Namespace:  "",
	}
	eh, err := api.ReferTypedObject[core.Middleware](server, ref)
	testutil.DiffError(t, nil, nil, err)

	testCasbinAuthzMiddleware(t, eh)

}

func TestEmptyNameNamespace(t *testing.T) {

	configs := []string{"./config-empty-name-namespace.yaml"}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "app/v1",
		Kind:       "CasbinAuthzMiddleware",
		Name:       "default",
		Namespace:  "",
	}
	eh, err := api.ReferTypedObject[core.Middleware](server, ref)
	testutil.DiffError(t, nil, nil, err)

	testCasbinAuthzMiddleware(t, eh)

}

func TestEmptySpec(t *testing.T) {

	configs := []string{"./config-empty-spec.yaml"}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "app/v1",
		Kind:       "CasbinAuthzMiddleware",
		Name:       "default",
		Namespace:  "",
	}
	eh, err := api.ReferTypedObject[core.Middleware](server, ref)
	testutil.DiffError(t, nil, nil, err)

	testCasbinAuthzMiddleware(t, eh)

}

func TestInvalidSpec(t *testing.T) {

	configs := []string{"./config-invalid-spec.yaml"}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	errPattern := regexp.MustCompile(core.ErrPrefix + `failed to load configs. config-invalid-spec.yaml`)
	testutil.DiffError(t, app.ErrAppMainLoadConfigs, errPattern, err)

}
