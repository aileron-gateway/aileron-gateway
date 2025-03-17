//go:build integration

package otelmeter_test

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/aileron-gateway/aileron-gateway/apis/kernel"
	"github.com/aileron-gateway/aileron-gateway/app/meter/otelmeter"
	"github.com/aileron-gateway/aileron-gateway/cmd/aileron/app"
	"github.com/aileron-gateway/aileron-gateway/core"
	"github.com/aileron-gateway/aileron-gateway/kernel/api"
	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
	"github.com/aileron-gateway/aileron-gateway/test/integration/common"
)

type meter interface {
	core.Middleware
	core.Tripperware
	core.Finalizer
}

// mockRoundTripper is a mock implementation of http.RoundTripper
type mockRoundTripper struct {
	err error
}

func (m *mockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	if m.err != nil {
		return nil, m.err
	}
	return &http.Response{
		StatusCode: http.StatusOK,
		Request:    req,
	}, nil
}

func check(t *testing.T, m meter) {
	t.Helper()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})
	h := m.Middleware(handler)

	r1 := httptest.NewRequest(http.MethodGet, "http://metrics.com/middleware", nil)
	w1 := httptest.NewRecorder()
	h.ServeHTTP(w1, r1)
	b1, _ := io.ReadAll(w1.Result().Body)
	testutil.Diff(t, http.StatusOK, w1.Result().StatusCode)
	testutil.Diff(t, "ok", string(b1))

	mockTransport := &mockRoundTripper{}
	rt := m.Tripperware(mockTransport)
	r2 := httptest.NewRequest(http.MethodGet, "http://metrics.com/tripperware", nil)
	_, err := rt.RoundTrip(r2)
	testutil.DiffError(t, nil, nil, err)

	time.Sleep(2 * time.Second)
}

func TestEmptyNameNamespace(t *testing.T) {
	configs := []string{"./config-empty-name-namespace.yaml"}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	var buf bytes.Buffer
	otelmeter.TestWriter = &buf
	defer func() {
		otelmeter.TestWriter = nil
	}()

	ref := &kernel.Reference{
		APIVersion: "app/v1",
		Kind:       "OpenTelemetryMeter",
		Name:       "",
		Namespace:  "",
	}
	m, err := api.ReferTypedObject[meter](server, ref)
	testutil.DiffError(t, nil, nil, err)
	defer m.Finalize()

	check(t, m)
	metrics := buf.String()
	t.Log(metrics)
	testutil.Diff(t, true, strings.Contains(metrics, `{"Type":"STRING","Value":"/middleware"}`))
	testutil.Diff(t, true, strings.Contains(metrics, `{"Type":"STRING","Value":"/tripperware"}`))
}

func TestEmptyName(t *testing.T) {
	configs := []string{"./config-empty-name.yaml"}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	var buf bytes.Buffer
	otelmeter.TestWriter = &buf
	defer func() {
		otelmeter.TestWriter = nil
	}()

	ref := &kernel.Reference{
		APIVersion: "app/v1",
		Kind:       "OpenTelemetryMeter",
		Name:       "",
		Namespace:  "testNamespace",
	}
	m, err := api.ReferTypedObject[meter](server, ref)
	testutil.DiffError(t, nil, nil, err)
	defer m.Finalize()

	check(t, m)
	metrics := buf.String()
	t.Log(metrics)
	testutil.Diff(t, true, strings.Contains(metrics, `{"Type":"STRING","Value":"/middleware"}`))
	testutil.Diff(t, true, strings.Contains(metrics, `{"Type":"STRING","Value":"/tripperware"}`))
}

func TestEmptyNamespace(t *testing.T) {
	configs := []string{"./config-empty-namespace.yaml"}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	var buf bytes.Buffer
	otelmeter.TestWriter = &buf
	defer func() {
		otelmeter.TestWriter = nil
	}()

	ref := &kernel.Reference{
		APIVersion: "app/v1",
		Kind:       "OpenTelemetryMeter",
		Name:       "testName",
		Namespace:  "",
	}
	m, err := api.ReferTypedObject[meter](server, ref)
	testutil.DiffError(t, nil, nil, err)
	defer m.Finalize()

	check(t, m)
	metrics := buf.String()
	t.Log(metrics)
	testutil.Diff(t, true, strings.Contains(metrics, `{"Type":"STRING","Value":"/middleware"}`))
	testutil.Diff(t, true, strings.Contains(metrics, `{"Type":"STRING","Value":"/tripperware"}`))
}

func TestEmptySpec(t *testing.T) {
	configs := []string{"./config-empty-spec.yaml"}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	var buf bytes.Buffer
	otelmeter.TestWriter = &buf
	defer func() {
		otelmeter.TestWriter = nil
	}()

	ref := &kernel.Reference{
		APIVersion: "app/v1",
		Kind:       "OpenTelemetryMeter",
		Name:       "testName",
		Namespace:  "testNamespace",
	}
	m, err := api.ReferTypedObject[meter](server, ref)
	testutil.DiffError(t, nil, nil, err)
	defer m.Finalize()

	check(t, m)
	metrics := buf.String()
	t.Log(metrics)
	testutil.Diff(t, true, strings.Contains(metrics, `{"Type":"STRING","Value":"/middleware"}`))
	testutil.Diff(t, true, strings.Contains(metrics, `{"Type":"STRING","Value":"/tripperware"}`))
}

func TestMinimalWithoutMetadata(t *testing.T) {
	configs := []string{"./config-minimal-without-metadata.yaml"}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	var buf bytes.Buffer
	otelmeter.TestWriter = &buf
	defer func() {
		otelmeter.TestWriter = nil
	}()

	ref := &kernel.Reference{
		APIVersion: "app/v1",
		Kind:       "OpenTelemetryMeter",
	}
	m, err := api.ReferTypedObject[meter](server, ref)
	testutil.DiffError(t, nil, nil, err)
	defer m.Finalize()

	check(t, m)
	metrics := buf.String()
	t.Log(metrics)
	testutil.Diff(t, true, strings.Contains(metrics, `{"Type":"STRING","Value":"/middleware"}`))
	testutil.Diff(t, true, strings.Contains(metrics, `{"Type":"STRING","Value":"/tripperware"}`))
}

func TestMinimalWithMetadata(t *testing.T) {
	configs := []string{"./config-minimal-with-metadata.yaml"}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	var buf bytes.Buffer
	otelmeter.TestWriter = &buf
	defer func() {
		otelmeter.TestWriter = nil
	}()

	ref := &kernel.Reference{
		APIVersion: "app/v1",
		Kind:       "OpenTelemetryMeter",
		Name:       "testName",
		Namespace:  "testNamespace",
	}
	m, err := api.ReferTypedObject[meter](server, ref)
	testutil.DiffError(t, nil, nil, err)
	defer m.Finalize()

	check(t, m)
	metrics := buf.String()
	t.Log(metrics)
	testutil.Diff(t, true, strings.Contains(metrics, `{"Type":"STRING","Value":"/middleware"}`))
	testutil.Diff(t, true, strings.Contains(metrics, `{"Type":"STRING","Value":"/tripperware"}`))
}

func TestInvalidSpec(t *testing.T) {
	configs := []string{"./config-invalid-spec.yaml"}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	errPattern := regexp.MustCompile(core.ErrPrefix + `failed to load configs. config-invalid-spec.yaml`)
	testutil.DiffError(t, app.ErrAppMainLoadConfigs, errPattern, err)
}
