//go:build integration
// +build integration

package healthcheck_test

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/aileron-gateway/aileron-gateway/apis/kernel"
	"github.com/aileron-gateway/aileron-gateway/cmd/aileron/app"
	"github.com/aileron-gateway/aileron-gateway/kernel/api"
	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
	"github.com/aileron-gateway/aileron-gateway/test/integration/common"
)

// healthChecker is a mock implementation that simulates a HealthChecker.
type healthChecker struct {
	timeout time.Duration
	status  bool
}

func (m *healthChecker) HealthCheck(ctx context.Context) (context.Context, bool) {
	if m.timeout > 0 {
		time.Sleep(m.timeout)
	}
	return ctx, m.status
}

func TestNoProbes(t *testing.T) {
	configs := []string{"./config-no-probes.yaml"}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "app/v1",
		Kind:       "HealthCheckHandler",
	}
	h, err := api.ReferTypedObject[http.Handler](server, ref)
	testutil.DiffError(t, nil, nil, err)

	r := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)

	resp := w.Result()
	body, err := io.ReadAll(resp.Body)
	testutil.Diff(t, nil, err)
	testutil.Diff(t, http.StatusOK, resp.StatusCode)
	testutil.Diff(t, []byte("{}"), body)
}

func TestSingleProbe(t *testing.T) {
	configs := []string{"./config-single-probe.yaml"}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	// Register testProbe as a healthChecker
	probeRef := &kernel.Reference{
		APIVersion: "container/v1",
		Kind:       "Container",
		Name:       "testProbe",
		Namespace:  "test",
	}
	probe := &healthChecker{timeout: time.Millisecond, status: true}
	common.PostTestResource(server, probeRef, probe)

	ref := &kernel.Reference{
		APIVersion: "app/v1",
		Kind:       "HealthCheckHandler",
	}
	h, err := api.ReferTypedObject[http.Handler](server, ref)
	testutil.DiffError(t, nil, nil, err)

	r := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)

	resp := w.Result()
	body, err := io.ReadAll(resp.Body)
	testutil.Diff(t, nil, err)
	testutil.Diff(t, http.StatusOK, resp.StatusCode)
	testutil.Diff(t, []byte("{}"), body)
}

func TestMultipleProbes(t *testing.T) {
	configs := []string{"./config-multiple-probes.yaml"}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	// Register testProbe1 and testProbe2 as healthCheckers
	probeRef1 := &kernel.Reference{
		APIVersion: "container/v1",
		Kind:       "Container",
		Name:       "testProbe1",
		Namespace:  "test",
	}
	probe1 := &healthChecker{timeout: time.Millisecond, status: true}
	common.PostTestResource(server, probeRef1, probe1)

	probeRef2 := &kernel.Reference{
		APIVersion: "container/v1",
		Kind:       "Container",
		Name:       "testProbe2",
		Namespace:  "test",
	}
	probe2 := &healthChecker{timeout: time.Millisecond, status: true}
	common.PostTestResource(server, probeRef2, probe2)

	ref := &kernel.Reference{
		APIVersion: "app/v1",
		Kind:       "HealthCheckHandler",
	}
	h, err := api.ReferTypedObject[http.Handler](server, ref)
	testutil.DiffError(t, nil, nil, err)

	r := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)

	resp := w.Result()
	body, err := io.ReadAll(resp.Body)
	testutil.Diff(t, nil, err)
	testutil.Diff(t, http.StatusOK, resp.StatusCode)
	testutil.Diff(t, []byte("{}"), body)
}

func TestErrorProbe(t *testing.T) {
	configs := []string{"./config-error-probe.yaml"}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	probeRef := &kernel.Reference{
		APIVersion: "container/v1",
		Kind:       "Container",
		Name:       "testProbe",
		Namespace:  "test",
	}
	probe := &healthChecker{timeout: time.Millisecond, status: false}
	common.PostTestResource(server, probeRef, probe)

	ref := &kernel.Reference{
		APIVersion: "app/v1",
		Kind:       "HealthCheckHandler",
	}
	h, err := api.ReferTypedObject[http.Handler](server, ref)
	testutil.DiffError(t, nil, nil, err)

	r := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)

	resp := w.Result()
	body, err := io.ReadAll(resp.Body)
	testutil.Diff(t, nil, err)
	testutil.Diff(t, http.StatusInternalServerError, resp.StatusCode)
	testutil.Diff(t, []byte(`{"status":500,"statusText":"Internal Server Error"}`), body)
}

func TestTimeout(t *testing.T) {
	configs := []string{"./config-timeout.yaml"}
	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	probeRef := &kernel.Reference{
		APIVersion: "container/v1",
		Kind:       "Container",
		Name:       "testProbe",
		Namespace:  "test",
	}
	probe := &healthChecker{timeout: time.Second * 2, status: true}
	common.PostTestResource(server, probeRef, probe)

	ref := &kernel.Reference{
		APIVersion: "app/v1",
		Kind:       "HealthCheckHandler",
	}
	h, err := api.ReferTypedObject[http.Handler](server, ref)
	testutil.DiffError(t, nil, nil, err)

	r := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	w := httptest.NewRecorder()

	ctx, cancel := context.WithTimeout(r.Context(), time.Millisecond*500)
	defer cancel()
	h.ServeHTTP(w, r.WithContext(ctx))

	resp := w.Result()
	body, err := io.ReadAll(resp.Body)
	testutil.Diff(t, nil, err)
	testutil.Diff(t, http.StatusGatewayTimeout, resp.StatusCode)
	testutil.Diff(t, []byte(`{"status":504,"statusText":"Gateway Timeout"}`), body)
}
