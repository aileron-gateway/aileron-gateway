//go:build integration

package opa_test

import (
	"context"
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
	"github.com/golang-jwt/jwt/v5"
)

func TestAuthzAPI(t *testing.T) {

	configs := []string{"./config-authz-api.yaml"}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "app/v1",
		Kind:       "OPAAuthzMiddleware",
	}
	m, err := api.ReferTypedObject[core.Middleware](server, ref)
	testutil.DiffError(t, nil, nil, err)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"message":"allowed"}`))
	})

	h := m.Middleware(handler)

	// authorization allowed
	r1 := httptest.NewRequest(http.MethodGet, "http://opa.com/allowed", nil)
	w1 := httptest.NewRecorder()
	h.ServeHTTP(w1, r1)
	b1, _ := io.ReadAll(w1.Result().Body)
	testutil.Diff(t, http.StatusOK, w1.Result().StatusCode)
	testutil.Diff(t, `{"message":"allowed"}`, string(b1))

	// authorization denied
	r2 := httptest.NewRequest(http.MethodGet, "http://opa.com/denied", nil)
	w2 := httptest.NewRecorder()
	h.ServeHTTP(w2, r2)
	b2, _ := io.ReadAll(w2.Result().Body)
	testutil.Diff(t, http.StatusForbidden, w2.Result().StatusCode)
	testutil.Diff(t, `{"status":403,"statusText":"Forbidden"}`, string(b2))

}

func TestAuthzAuth(t *testing.T) {

	configs := []string{"./config-authz-auth.yaml"}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "app/v1",
		Kind:       "OPAAuthzMiddleware",
	}
	m, err := api.ReferTypedObject[core.Middleware](server, ref)
	testutil.DiffError(t, nil, nil, err)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"message":"allowed"}`))
	})

	h := m.Middleware(handler)

	// Allowed
	// "admin" is allowed (input jwt.MapClaims)
	r1 := httptest.NewRequest(http.MethodGet, "http://opa.com/", nil).
		WithContext(context.WithValue(context.Background(), "AuthClaims", jwt.MapClaims{"role": []string{"admin", "user"}}))
	w1 := httptest.NewRecorder()
	h.ServeHTTP(w1, r1)
	b1, _ := io.ReadAll(w1.Result().Body)
	testutil.Diff(t, http.StatusOK, w1.Result().StatusCode)
	testutil.Diff(t, `{"message":"allowed"}`, string(b1))

	// Allowed
	// "admin" is allowed (input map[string]any)
	r2 := httptest.NewRequest(http.MethodGet, "http://opa.com/", nil).
		WithContext(context.WithValue(context.Background(), "AuthClaims", map[string]any{"role": []string{"admin", "user"}}))
	w2 := httptest.NewRecorder()
	h.ServeHTTP(w2, r2)
	b2, _ := io.ReadAll(w2.Result().Body)
	testutil.Diff(t, http.StatusOK, w2.Result().StatusCode)
	testutil.Diff(t, `{"message":"allowed"}`, string(b2))

	// Denied
	// "user" is not allowed
	r3 := httptest.NewRequest(http.MethodGet, "http://opa.com/", nil).
		WithContext(context.WithValue(context.Background(), "AuthClaims", map[string]any{"role": []string{"user"}}))
	w3 := httptest.NewRecorder()
	h.ServeHTTP(w3, r3)
	b3, _ := io.ReadAll(w3.Result().Body)
	testutil.Diff(t, http.StatusForbidden, w3.Result().StatusCode)
	testutil.Diff(t, `{"status":403,"statusText":"Forbidden"}`, string(b3))

}

func TestAuthzHeader(t *testing.T) {

	configs := []string{"./config-authz-header.yaml"}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "app/v1",
		Kind:       "OPAAuthzMiddleware",
	}
	m, err := api.ReferTypedObject[core.Middleware](server, ref)
	testutil.DiffError(t, nil, nil, err)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"message":"allowed"}`))
	})

	h := m.Middleware(handler)

	// Allowed
	// Allowed because "allowed" contained in the header.
	r1 := httptest.NewRequest(http.MethodGet, "http://opa.com/", nil)
	r1.Header["Test"] = []string{"foo", "bar", "allowed"}
	w1 := httptest.NewRecorder()
	h.ServeHTTP(w1, r1)
	b1, _ := io.ReadAll(w1.Result().Body)
	testutil.Diff(t, http.StatusOK, w1.Result().StatusCode)
	testutil.Diff(t, `{"message":"allowed"}`, string(b1))

	// Denied
	// Denied because "allowed" not contained in the header.
	r2 := httptest.NewRequest(http.MethodGet, "http://opa.com/", nil)
	r2.Header["Test"] = []string{"foo", "bar", "denied"}
	w2 := httptest.NewRecorder()
	h.ServeHTTP(w2, r2)
	b2, _ := io.ReadAll(w2.Result().Body)
	testutil.Diff(t, http.StatusForbidden, w2.Result().StatusCode)
	testutil.Diff(t, `{"status":403,"statusText":"Forbidden"}`, string(b2))

}

func TestAuthzQuery(t *testing.T) {

	configs := []string{"./config-authz-query.yaml"}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "app/v1",
		Kind:       "OPAAuthzMiddleware",
	}
	m, err := api.ReferTypedObject[core.Middleware](server, ref)
	testutil.DiffError(t, nil, nil, err)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"message":"allowed"}`))
	})

	h := m.Middleware(handler)

	// Allowed
	// Allowed because "allowed" contained in the query.
	r1 := httptest.NewRequest(http.MethodGet, "http://opa.com/?test=foo&test-Bar&test=allowed", nil)
	w1 := httptest.NewRecorder()
	h.ServeHTTP(w1, r1)
	b1, _ := io.ReadAll(w1.Result().Body)
	testutil.Diff(t, http.StatusOK, w1.Result().StatusCode)
	testutil.Diff(t, `{"message":"allowed"}`, string(b1))

	// Denied
	// Denied because "allowed" not contained in the query.
	r2 := httptest.NewRequest(http.MethodGet, "http://opa.com/?test=foo&test-Bar&test=denied", nil)
	w2 := httptest.NewRecorder()
	h.ServeHTTP(w2, r2)
	b2, _ := io.ReadAll(w2.Result().Body)
	testutil.Diff(t, http.StatusForbidden, w2.Result().StatusCode)
	testutil.Diff(t, `{"status":403,"statusText":"Forbidden"}`, string(b2))

}

func TestAuthzMultipleFile(t *testing.T) {

	configs := []string{"./config-authz-multiple-file.yaml"}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "app/v1",
		Kind:       "OPAAuthzMiddleware",
	}
	m, err := api.ReferTypedObject[core.Middleware](server, ref)
	testutil.DiffError(t, nil, nil, err)

	// authorization allowed
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"message":"allowed"}`))
	})

	h := m.Middleware(handler)

	r1 := httptest.NewRequest(http.MethodGet, "http://opa.com/allowed", nil)
	r1.Header = http.Header{"Authorization": {"Bearer allowed-token"}}
	w1 := httptest.NewRecorder()
	h.ServeHTTP(w1, r1)
	b1, _ := io.ReadAll(w1.Result().Body)
	testutil.Diff(t, http.StatusOK, w1.Result().StatusCode)
	testutil.Diff(t, `{"message":"allowed"}`, string(b1))

	// authorization denied
	r2 := httptest.NewRequest(http.MethodGet, "http://opa.com/denied", nil)
	r2.Header = http.Header{"Authorization": {"Bearer denied-token"}}
	w2 := httptest.NewRecorder()
	h.ServeHTTP(w2, r2)
	b2, _ := io.ReadAll(w2.Result().Body)
	testutil.Diff(t, http.StatusForbidden, w2.Result().StatusCode)
	testutil.Diff(t, `{"status":403,"statusText":"Forbidden"}`, string(b2))
}

func TestPolicyFileError(t *testing.T) {

	configs := []string{"./config-policy-file-error.yaml"}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "app/v1",
		Kind:       "OPAAuthzMiddleware",
	}
	_, err = api.ReferTypedObject[core.Middleware](server, ref)
	testutil.DiffError(t, core.ErrCoreGenCreateObject, regexp.MustCompile("failed to create OPAAuthzMiddleware"), err)
}

func TestPolicyNotFound(t *testing.T) {

	configs := []string{"./config-policy-notfound.yaml"}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "app/v1",
		Kind:       "OPAAuthzMiddleware",
	}
	_, err = api.ReferTypedObject[core.Middleware](server, ref)
	testutil.DiffError(t, core.ErrCoreGenCreateObject, regexp.MustCompile("failed to create OPAAuthzMiddleware"), err)
}
