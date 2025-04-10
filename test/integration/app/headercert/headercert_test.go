//go:build integration

package headercert_test

import (
	"encoding/base64"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/aileron-gateway/aileron-gateway/apis/kernel"
	"github.com/aileron-gateway/aileron-gateway/cmd/aileron/app"
	"github.com/aileron-gateway/aileron-gateway/core"
	"github.com/aileron-gateway/aileron-gateway/kernel/api"
	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
	"github.com/aileron-gateway/aileron-gateway/test/integration/common"
)

var (
	cert, _ = os.ReadFile("./client.crt")
	fp, _   = os.ReadFile("./fingerprint.txt")
)

func TestRootCAs(t *testing.T) {
	configs := []string{"./config-rootcas.yaml"}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "app/v1",
		Kind:       "HeaderCertMiddleware",
		Name:       "",
		Namespace:  "",
	}
	m, err := api.ReferTypedObject[core.Middleware](server, ref)
	testutil.DiffError(t, nil, nil, err)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})
	h := m.Middleware(handler)

	r1 := httptest.NewRequest(http.MethodGet, "http://headercert.com/rootcas", nil)
	r1.Header.Set("X-SSL-Client-Cert", base64.URLEncoding.EncodeToString(cert))
	w1 := httptest.NewRecorder()
	h.ServeHTTP(w1, r1)
	b1, _ := io.ReadAll(w1.Result().Body)
	testutil.Diff(t, http.StatusOK, w1.Result().StatusCode)
	testutil.Diff(t, "ok", string(b1))
	return
}

func TestCertHeader(t *testing.T) {
	configs := []string{"./config-certheader.yaml"}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "app/v1",
		Kind:       "HeaderCertMiddleware",
		Name:       "",
		Namespace:  "",
	}
	m, err := api.ReferTypedObject[core.Middleware](server, ref)
	testutil.DiffError(t, nil, nil, err)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})
	h := m.Middleware(handler)

	r1 := httptest.NewRequest(http.MethodGet, "http://headercert.com/certHeader", nil)
	r1.Header.Set("X-SSL-Client-Cert-Test", base64.URLEncoding.EncodeToString(cert))
	w1 := httptest.NewRecorder()
	h.ServeHTTP(w1, r1)
	b1, _ := io.ReadAll(w1.Result().Body)
	testutil.Diff(t, http.StatusOK, w1.Result().StatusCode)
	testutil.Diff(t, "ok", string(b1))
	return
}

func TestFingerprintHeader(t *testing.T) {
	configs := []string{"./config-fingerprintheader.yaml"}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "app/v1",
		Kind:       "HeaderCertMiddleware",
		Name:       "",
		Namespace:  "",
	}
	m, err := api.ReferTypedObject[core.Middleware](server, ref)
	testutil.DiffError(t, nil, nil, err)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})
	h := m.Middleware(handler)

	// test valid fingerprint header name
	r1 := httptest.NewRequest(http.MethodGet, "http://headercert.com/certHeader", nil)
	r1.Header.Set("X-SSL-Client-Cert-Test", base64.URLEncoding.EncodeToString(cert))
	r1.Header.Set("X-SSL-Client-Fingerprint-Test", string(fp))
	w1 := httptest.NewRecorder()
	h.ServeHTTP(w1, r1)
	b1, _ := io.ReadAll(w1.Result().Body)
	testutil.Diff(t, http.StatusOK, w1.Result().StatusCode)
	testutil.Diff(t, "ok", string(b1))
	
	// test invalid fingerprint header name
	r2 := httptest.NewRequest(http.MethodGet, "http://headercert.com/certHeader", nil)
	r2.Header.Set("X-SSL-Client-Cert-Test", base64.URLEncoding.EncodeToString(cert))
	r2.Header.Set("X-SSL-Client-Fingerprint", string(fp))
	w2 := httptest.NewRecorder()
	h.ServeHTTP(w2, r2)
	testutil.Diff(t, http.StatusUnauthorized, w2.Result().StatusCode)
	testutil.Diff(t, `{"status":401,"statusText":"Unauthorized"}`, w2.Body.String())
	return
}
