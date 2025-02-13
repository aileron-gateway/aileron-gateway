//go:build integration
// +build integration

package soaprest_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"

	"github.com/aileron-gateway/aileron-gateway/apis/kernel"
	"github.com/aileron-gateway/aileron-gateway/cmd/aileron/app"
	"github.com/aileron-gateway/aileron-gateway/core"
	"github.com/aileron-gateway/aileron-gateway/kernel/api"
	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
	"github.com/aileron-gateway/aileron-gateway/test/integration/common"
)

const (
	reqXML = `<?xml version="1.0" encoding="UTF-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Header></soap:Header>
  <soap:Body>
    <Value>Hello</Value>
  </soap:Body>
</soap:Envelope>`
	respJSON = `{"soap:Envelope":{"_namespace":{"soap":"http://schemas.xmlsoap.org/soap/envelope/"},"soap:Body":{"Value":"Hello"}}}`
)

func testSOAPRESTMiddleware(t *testing.T, m core.Middleware) {
	t.Helper()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Write([]byte(respJSON))
	})

	h := m.Middleware(handler)

	r1 := httptest.NewRequest(http.MethodGet, "http://SOAPREST.com/test", strings.NewReader(reqXML))
	r1.Header.Set("Content-Type", "text/xml")
	w1 := httptest.NewRecorder()
	h.ServeHTTP(w1, r1)
	b1, _ := io.ReadAll(w1.Result().Body)
	testutil.Diff(t, http.StatusOK, w1.Result().StatusCode)
	testutil.Diff(t, reqXML, string(b1))
}

func TestMinimalWithoutMetadata(t *testing.T) {
	configs := []string{"./config-minimal-without-metadata.yaml"}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "app/v1",
		Kind:       "SOAPRESTMiddleware",
		Name:       "default",
		Namespace:  "",
	}
	eh, err := api.ReferTypedObject[core.Middleware](server, ref)
	testutil.DiffError(t, nil, nil, err)

	testSOAPRESTMiddleware(t, eh)
}

func TestMinimalWithMetadata(t *testing.T) {
	configs := []string{"./config-minimal-with-metadata.yaml"}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "app/v1",
		Kind:       "SOAPRESTMiddleware",
		Name:       "testName",
		Namespace:  "testNamespace",
	}
	eh, err := api.ReferTypedObject[core.Middleware](server, ref)
	testutil.DiffError(t, nil, nil, err)

	testSOAPRESTMiddleware(t, eh)
}

func TestEmptyName(t *testing.T) {
	configs := []string{"./config-empty-name.yaml"}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "app/v1",
		Kind:       "SOAPRESTMiddleware",
		Name:       "default",
		Namespace:  "testNamespace",
	}
	eh, err := api.ReferTypedObject[core.Middleware](server, ref)
	testutil.DiffError(t, nil, nil, err)

	testSOAPRESTMiddleware(t, eh)
}

func TestEmptyNamespace(t *testing.T) {
	configs := []string{"./config-empty-namespace.yaml"}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "app/v1",
		Kind:       "SOAPRESTMiddleware",
		Name:       "testName",
		Namespace:  "",
	}
	eh, err := api.ReferTypedObject[core.Middleware](server, ref)
	testutil.DiffError(t, nil, nil, err)

	testSOAPRESTMiddleware(t, eh)
}

func TestEmptyNameNamespace(t *testing.T) {
	configs := []string{"./config-empty-name-namespace.yaml"}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "app/v1",
		Kind:       "SOAPRESTMiddleware",
		Name:       "default",
		Namespace:  "",
	}
	eh, err := api.ReferTypedObject[core.Middleware](server, ref)
	testutil.DiffError(t, nil, nil, err)

	testSOAPRESTMiddleware(t, eh)
}

func TestEmptySpec(t *testing.T) {
	configs := []string{"./config-empty-spec.yaml"}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "app/v1",
		Kind:       "SOAPRESTMiddleware",
		Name:       "default",
		Namespace:  "",
	}
	eh, err := api.ReferTypedObject[core.Middleware](server, ref)
	testutil.DiffError(t, nil, nil, err)

	testSOAPRESTMiddleware(t, eh)
}

func TestInvalidSpec(t *testing.T) {
	configs := []string{"./config-invalid-spec.yaml"}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	errPattern := regexp.MustCompile(core.ErrPrefix + `failed to load configs. config-invalid-spec.yaml`)
	testutil.DiffError(t, app.ErrAppMainLoadConfigs, errPattern, err)
}
