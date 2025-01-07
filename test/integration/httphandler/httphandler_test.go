//go:build integration
// +build integration

package httphandler_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aileron-gateway/aileron-gateway/apis/kernel"
	"github.com/aileron-gateway/aileron-gateway/cmd/aileron/app"
	"github.com/aileron-gateway/aileron-gateway/kernel/api"
	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
	"github.com/aileron-gateway/aileron-gateway/test/integration/common"
)

type testMiddleware struct {
	called int
}

func (m *testMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m.called += 1
		next.ServeHTTP(w, r)
	})
}

func TestMiddleware(t *testing.T) {

	configs := []string{
		testDataDir + "config-middleware.yaml",
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
	middle1 := &kernel.Reference{
		APIVersion: "container/v1",
		Kind:       "TestMiddleware",
		Name:       "test1",
		Namespace:  "",
	}
	m1 := &testMiddleware{}
	common.PostTestResource(server, middle1, m1)
	middle2 := &kernel.Reference{
		APIVersion: "container/v1",
		Kind:       "TestMiddleware",
		Name:       "test2",
		Namespace:  "",
	}
	m2 := &testMiddleware{}
	common.PostTestResource(server, middle2, m2)

	ref := &kernel.Reference{
		APIVersion: "core/v1",
		Kind:       "HTTPHandler",
		Name:       "default",
		Namespace:  "",
	}
	h, err := api.ReferTypedObject[http.Handler](server, ref)
	testutil.DiffError(t, nil, nil, err)

	r := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	b, _ := io.ReadAll(w.Body)
	testutil.Diff(t, http.StatusOK, w.Result().StatusCode)
	testutil.Diff(t, "ok", w.Header().Get("test"))
	testutil.Diff(t, []byte("test"), b)
	testutil.Diff(t, 1, m1.called)
	testutil.Diff(t, 1, m2.called)

}
