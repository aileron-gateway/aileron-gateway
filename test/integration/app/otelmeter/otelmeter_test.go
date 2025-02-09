//go:build integration
// +build integration

package otelmeter_test

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
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
	v1 "go.opentelemetry.io/proto/otlp/metrics/v1"
	"google.golang.org/protobuf/proto"
)

func TestExportStdout_middleware(t *testing.T) {
	configs := []string{"./config-export-stdout.yaml"}

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
	m, err := api.ReferTypedObject[core.Middleware](server, ref)
	testutil.DiffError(t, nil, nil, err)

	// Use OpenTelemetryMeter Middleware on an actual handler
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})
	h := m.Middleware(handler)

	// Perform GET requests.
	for i := 0; i < 5; i++ {
		r := httptest.NewRequest(http.MethodGet, "http://metrics.com/get", nil)
		w := httptest.NewRecorder()
		h.ServeHTTP(w, r)
		b, _ := io.ReadAll(w.Result().Body)
		testutil.Diff(t, http.StatusOK, w.Result().StatusCode)
		testutil.Diff(t, "ok", string(b))
	}
	// Perform POST requests.
	for i := 0; i < 5; i++ {
		r := httptest.NewRequest(http.MethodPost, "http://metrics.com/post", bytes.NewBuffer([]byte("test body")))
		w := httptest.NewRecorder()
		h.ServeHTTP(w, r)
		b, _ := io.ReadAll(w.Result().Body)
		testutil.Diff(t, http.StatusOK, w.Result().StatusCode)
		testutil.Diff(t, "ok", string(b))
	}

	time.Sleep(2 * time.Second)
	metrics := buf.String()
	t.Log(metrics)
	testutil.Diff(t, true, strings.Contains(metrics, `{"Key":"code","Value":{"Type":"INT64","Value":200}},{"Key":"host","Value":{"Type":"STRING","Value":"metrics.com"}},{"Key":"method","Value":{"Type":"STRING","Value":"GET"}},{"Key":"path","Value":{"Type":"STRING","Value":"/get"}}`))
	testutil.Diff(t, true, strings.Contains(metrics, `{"Key":"code","Value":{"Type":"INT64","Value":200}},{"Key":"host","Value":{"Type":"STRING","Value":"metrics.com"}},{"Key":"method","Value":{"Type":"STRING","Value":"POST"}},{"Key":"path","Value":{"Type":"STRING","Value":"/post"}}`))
}

func TestExportStdout_tripperware(t *testing.T) {
	configs := []string{"./config-export-stdout.yaml"}

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
	m, err := api.ReferTypedObject[core.Tripperware](server, ref)
	testutil.DiffError(t, nil, nil, err)

	mockTransport := &mockRoundTripper{}
	rt := m.Tripperware(mockTransport)

	// Perform GET requests.
	for i := 0; i < 5; i++ {
		r := httptest.NewRequest(http.MethodGet, "http://metrics.com/get", nil)
		_, err := rt.RoundTrip(r)
		testutil.DiffError(t, nil, nil, err)
	}
	// Perform POST requests.
	for i := 0; i < 5; i++ {
		r := httptest.NewRequest(http.MethodPost, "http://metrics.com/post", bytes.NewBuffer([]byte("test body")))
		_, err := rt.RoundTrip(r)
		testutil.DiffError(t, nil, nil, err)
	}

	time.Sleep(2 * time.Second)
	metrics := buf.String()
	t.Log(metrics)
	testutil.Diff(t, true, strings.Contains(metrics, `{"Key":"code","Value":{"Type":"INT64","Value":200}},{"Key":"host","Value":{"Type":"STRING","Value":"metrics.com"}},{"Key":"method","Value":{"Type":"STRING","Value":"GET"}},{"Key":"path","Value":{"Type":"STRING","Value":"/get"}}`))
	testutil.Diff(t, true, strings.Contains(metrics, `{"Key":"code","Value":{"Type":"INT64","Value":200}},{"Key":"host","Value":{"Type":"STRING","Value":"metrics.com"}},{"Key":"method","Value":{"Type":"STRING","Value":"POST"}},{"Key":"path","Value":{"Type":"STRING","Value":"/post"}}`))
}

func TestExportHTTP(t *testing.T) {
	configs := []string{"./config-export-http.yaml"}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	msg := &v1.MetricsData{}
	called := false
	svr := &http.Server{
		Addr: "0.0.0.0:4318",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if called {
				return
			}
			called = true
			b, _ := io.ReadAll(r.Body)
			proto.Unmarshal(b, msg)
		}),
	}
	go svr.ListenAndServe()
	defer svr.Close()
	time.Sleep(time.Second) // Wait the server start.

	ref := &kernel.Reference{
		APIVersion: "app/v1",
		Kind:       "OpenTelemetryMeter",
	}
	m, err := api.ReferTypedObject[core.Tripperware](server, ref)
	testutil.DiffError(t, nil, nil, err)

	mockTransport := &mockRoundTripper{}
	rt := m.Tripperware(mockTransport)

	// Perform GET requests.
	for i := 0; i < 5; i++ {
		r := httptest.NewRequest(http.MethodGet, "http://metrics.com/get", nil)
		_, err := rt.RoundTrip(r)
		testutil.DiffError(t, nil, nil, err)
	}
	// Perform POST requests.
	for i := 0; i < 5; i++ {
		r := httptest.NewRequest(http.MethodPost, "http://metrics.com/post", bytes.NewBuffer([]byte("test body")))
		_, err := rt.RoundTrip(r)
		testutil.DiffError(t, nil, nil, err)
	}

	time.Sleep(2 * time.Second)
	t.Log(msg)

	metrics := ""
	for _, m := range msg.ResourceMetrics[0].ScopeMetrics[1].Metrics {
		if m.Name == "http_client_requests_total" {
			for _, p := range m.Data.(*v1.Metric_Sum).Sum.DataPoints {
				metrics += p.String() + "\n"
			}
		}
	}

	t.Log(metrics)
	testutil.Diff(t, true, strings.Contains(metrics, `value:{string_value:"GET"}`))
	testutil.Diff(t, true, strings.Contains(metrics, `value:{string_value:"POST"}`))
	testutil.Diff(t, true, strings.Contains(metrics, `value:{string_value:"/get"}`))
	testutil.Diff(t, true, strings.Contains(metrics, `value:{string_value:"/post"}`))
}
