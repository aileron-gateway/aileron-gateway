//go:build integration
// +build integration

package soaprest_test

import (
	"io"
	"net/http"
	"net/http/httptest"
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
	respJSON = `{"soap_Envelope":{"namespaceKey":{"soap":"http://schemas.xmlsoap.org/soap/envelope/"},"soap_Body":{"Value":"Hello"},"soap_Header":{}}}`

	reqModifiedKeyXML = `<?xml version="1.0" encoding="UTF-8"?>
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/">
  <SOAP-ENV:Header></SOAP-ENV:Header>
  <SOAP-ENV:Body>
    <ElementNode testAttribute="example">testTextNode</ElementNode>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`
	respModifiedKeyJSON = `{"SOAP-ENV_Envelope":{"SOAP-ENV_Body":{"ElementNode":{"#text":"testTextNode","@attr":{"testAttribute":"example"}}},"SOAP-ENV_Header":{},"_nsKey":{"SOAP-ENV":"http://schemas.xmlsoap.org/soap/envelope/"}}}`

	reqExtractStringElementXML = `<?xml version="1.0" encoding="UTF-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Header></soap:Header>
  <soap:Body>
    <StringElementNode>"StringWithDoubleQuotations"</StringElementNode>
  </soap:Body>
</soap:Envelope>`

	respExtractStringElementXML = `<?xml version="1.0" encoding="UTF-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Header></soap:Header>
  <soap:Body>
    <StringElementNode>StringWithDoubleQuotations</StringElementNode>
  </soap:Body>
</soap:Envelope>`

	respExtractStringElementJSON = `{"soap_Envelope":{"namespaceKey":{"soap":"http://schemas.xmlsoap.org/soap/envelope/"},"soap_Body":{"StringElementNode":"StringWithDoubleQuotations"},"soap_Header":{}}}`

	reqExtractBooleanElementXML = `<?xml version="1.0" encoding="UTF-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Header></soap:Header>
  <soap:Body>
    <BooleanElementNode>true</BooleanElementNode>
    <BooleanElementNode>false</BooleanElementNode>
  </soap:Body>
</soap:Envelope>`

	respExtractBooleanElementJSON = `{"soap_Envelope":{"namespaceKey":{"soap":"http://schemas.xmlsoap.org/soap/envelope/"},"soap_Body":{"BooleanElementNode":[true,false]},"soap_Header":{}}}`

	reqExtractIntegerElementXML = `<?xml version="1.0" encoding="UTF-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Header></soap:Header>
  <soap:Body>
    <IntegerElementNode>0</IntegerElementNode>
    <IntegerElementNode>100</IntegerElementNode>
  </soap:Body>
</soap:Envelope>`

	respExtractIntegerElementJSON = `{"soap_Envelope":{"namespaceKey":{"soap":"http://schemas.xmlsoap.org/soap/envelope/"},"soap_Body":{"IntegerElementNode":[0,100]},"soap_Header":{}}}`

	reqExtractFloatElementXML = `<?xml version="1.0" encoding="UTF-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Header></soap:Header>
  <soap:Body>
    <FloatElementNode>0.1</FloatElementNode>
    <FloatElementNode>3.141592653589793238462643</FloatElementNode>
  </soap:Body>
</soap:Envelope>`

	respExtractFloatElementXML = `<?xml version="1.0" encoding="UTF-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Header></soap:Header>
  <soap:Body>
    <FloatElementNode>0.1</FloatElementNode>
    <FloatElementNode>3.141592653589793</FloatElementNode>
  </soap:Body>
</soap:Envelope>`

	respExtractFloatElementJSON = `{"soap_Envelope":{"namespaceKey":{"soap":"http://schemas.xmlsoap.org/soap/envelope/"},"soap_Body":{"FloatElementNode":[0.1,3.141592653589793]},"soap_Header":{}}}`
)

func TestCorrectPath(t *testing.T) {
	configs := []string{"./config-path.yaml"}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "app/v1",
		Kind:       "SOAPRESTMiddleware",
		Name:       "test",
		Namespace:  "testNamespace",
	}

	m, err := api.ReferTypedObject[core.Middleware](server, ref)
	testutil.DiffError(t, nil, nil, err)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "failed to read body", http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()
		testutil.Diff(t, respJSON, string(body))

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Write([]byte(respJSON))
	})

	h := m.Middleware(handler)

	r := httptest.NewRequest(http.MethodPost, "http://soaprest.com/soap", strings.NewReader(reqXML))
	r.Header.Set("Content-Type", "text/xml")
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	b, _ := io.ReadAll(w.Result().Body)
	testutil.Diff(t, http.StatusOK, w.Result().StatusCode)
	testutil.Diff(t, reqXML, string(b))
}

func TestWrongPath(t *testing.T) {
	configs := []string{"./config-path.yaml"}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "app/v1",
		Kind:       "SOAPRESTMiddleware",
		Name:       "test",
		Namespace:  "testNamespace",
	}

	m, err := api.ReferTypedObject[core.Middleware](server, ref)
	testutil.DiffError(t, nil, nil, err)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "failed to read body", http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()

		// since requests are being sent to different paths,
		// the conversion from SOAP/XML to REST/JSON will not take place
		testutil.Diff(t, reqXML, string(body))

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Write([]byte(respJSON))
	})

	h := m.Middleware(handler)

	r := httptest.NewRequest(http.MethodPost, "http://soaprest.com/wrong", strings.NewReader(reqXML))
	r.Header.Set("Content-Type", "text/xml")
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	b, _ := io.ReadAll(w.Result().Body)
	testutil.Diff(t, http.StatusOK, w.Result().StatusCode)

	// similarly to the request,
	// the response will not undergo conversion from REST/JSON to SOAP/XML
	testutil.Diff(t, respJSON, string(b))
}

func TestModifiedKey(t *testing.T) {
	configs := []string{"./config-modified-key.yaml"}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "app/v1",
		Kind:       "SOAPRESTMiddleware",
		Name:       "test",
		Namespace:  "testNamespace",
	}

	m, err := api.ReferTypedObject[core.Middleware](server, ref)
	testutil.DiffError(t, nil, nil, err)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "failed to read body", http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()

		testutil.Diff(t, respModifiedKeyJSON, string(body))

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Write([]byte(respModifiedKeyJSON))
	})

	h := m.Middleware(handler)

	r := httptest.NewRequest(http.MethodPost, "http://soaprest.com/", strings.NewReader(reqModifiedKeyXML))
	r.Header.Set("Content-Type", "text/xml")
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	b, _ := io.ReadAll(w.Result().Body)
	testutil.Diff(t, http.StatusOK, w.Result().StatusCode)

	testutil.Diff(t, reqModifiedKeyXML, string(b))
}

func TestEnabledExtractStringElement(t *testing.T) {
	configs := []string{"./config-enabled-extract-string-element.yaml"}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "app/v1",
		Kind:       "SOAPRESTMiddleware",
		Name:       "test",
		Namespace:  "testNamespace",
	}

	m, err := api.ReferTypedObject[core.Middleware](server, ref)
	testutil.DiffError(t, nil, nil, err)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "failed to read body", http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()

		testutil.Diff(t, respExtractStringElementJSON, string(body))

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Write([]byte(respExtractStringElementJSON))
	})

	h := m.Middleware(handler)

	r := httptest.NewRequest(http.MethodPost, "http://soaprest.com/", strings.NewReader(reqExtractStringElementXML))
	r.Header.Set("Content-Type", "text/xml")
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	b, _ := io.ReadAll(w.Result().Body)
	testutil.Diff(t, http.StatusOK, w.Result().StatusCode)

	// In the conversion to XML, double quotes are not added,
	// so a comparison will be made with XML that does not include double quotes
	testutil.Diff(t, respExtractStringElementXML, string(b))
}

func TestEnabledExtractBooleanElement(t *testing.T) {
	configs := []string{"./config-enabled-extract-boolean-element.yaml"}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "app/v1",
		Kind:       "SOAPRESTMiddleware",
		Name:       "test",
		Namespace:  "testNamespace",
	}

	m, err := api.ReferTypedObject[core.Middleware](server, ref)
	testutil.DiffError(t, nil, nil, err)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "failed to read body", http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()

		testutil.Diff(t, respExtractBooleanElementJSON, string(body))

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Write([]byte(respExtractBooleanElementJSON))
	})

	h := m.Middleware(handler)

	r := httptest.NewRequest(http.MethodPost, "http://soaprest.com/", strings.NewReader(reqExtractBooleanElementXML))
	r.Header.Set("Content-Type", "text/xml")
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	b, _ := io.ReadAll(w.Result().Body)
	testutil.Diff(t, http.StatusOK, w.Result().StatusCode)
	testutil.Diff(t, reqExtractBooleanElementXML, string(b))
}

func TestEnabledExtractIntegerElement(t *testing.T) {
	configs := []string{"./config-enabled-extract-integer-element.yaml"}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "app/v1",
		Kind:       "SOAPRESTMiddleware",
		Name:       "test",
		Namespace:  "testNamespace",
	}

	m, err := api.ReferTypedObject[core.Middleware](server, ref)
	testutil.DiffError(t, nil, nil, err)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "failed to read body", http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()

		testutil.Diff(t, respExtractIntegerElementJSON, string(body))

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Write([]byte(respExtractIntegerElementJSON))
	})

	h := m.Middleware(handler)

	r := httptest.NewRequest(http.MethodPost, "http://soaprest.com/", strings.NewReader(reqExtractIntegerElementXML))
	r.Header.Set("Content-Type", "text/xml")
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	b, _ := io.ReadAll(w.Result().Body)
	testutil.Diff(t, http.StatusOK, w.Result().StatusCode)
	testutil.Diff(t, reqExtractIntegerElementXML, string(b))
}

func TestEnabledExtractFloatElement(t *testing.T) {
	configs := []string{"./config-enabled-extract-float-element.yaml"}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "app/v1",
		Kind:       "SOAPRESTMiddleware",
		Name:       "test",
		Namespace:  "testNamespace",
	}

	m, err := api.ReferTypedObject[core.Middleware](server, ref)
	testutil.DiffError(t, nil, nil, err)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "failed to read body", http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()

		testutil.Diff(t, respExtractFloatElementJSON, string(body))

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Write([]byte(respExtractFloatElementJSON))
	})

	h := m.Middleware(handler)

	r := httptest.NewRequest(http.MethodPost, "http://soaprest.com/", strings.NewReader(reqExtractFloatElementXML))
	r.Header.Set("Content-Type", "text/xml")
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	b, _ := io.ReadAll(w.Result().Body)
	testutil.Diff(t, http.StatusOK, w.Result().StatusCode)

	// In the conversion to XML, if `extractFloatElement` is enabled, precision degradation will occur;
	// therefore, a comparison will be made with the XML containing elements that have experienced precision loss
	testutil.Diff(t, respExtractFloatElementXML, string(b))
}
