//go:build integration

package casbin_test

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"regexp"
	"testing"
	"time"

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
		Kind:       "CasbinAuthzMiddleware",
	}
	m, err := api.ReferTypedObject[core.Middleware](server, ref)
	testutil.DiffError(t, nil, nil, err)

	handler := m.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"message":"allowed"}`))
	}))

	h := m.Middleware(handler)

	// Path "/allowed" is allowed
	r1 := httptest.NewRequest(http.MethodGet, "http://casbin.com/allowed", nil)
	w1 := httptest.NewRecorder()
	h.ServeHTTP(w1, r1)
	b1, _ := io.ReadAll(w1.Result().Body)
	testutil.Diff(t, http.StatusOK, w1.Result().StatusCode)
	testutil.Diff(t, `{"message":"allowed"}`, string(b1))

	// POST is denied.
	r2 := httptest.NewRequest(http.MethodPost, "http://casbin.com/allowed", nil)
	w2 := httptest.NewRecorder()
	h.ServeHTTP(w2, r2)
	b2, _ := io.ReadAll(w2.Result().Body)
	testutil.Diff(t, http.StatusForbidden, w2.Result().StatusCode)
	testutil.Diff(t, `{"status":403,"statusText":"Forbidden"}`, string(b2))

	// Path "/denied" is denied.
	r3 := httptest.NewRequest(http.MethodGet, "http://casbin.com/denied", nil)
	w3 := httptest.NewRecorder()
	h.ServeHTTP(w3, r3)
	b3, _ := io.ReadAll(w3.Result().Body)
	testutil.Diff(t, http.StatusForbidden, w3.Result().StatusCode)
	testutil.Diff(t, `{"status":403,"statusText":"Forbidden"}`, string(b3))

}

func TestAuthzAuthCsv(t *testing.T) {

	configs := []string{"./config-authz-auth-csv.yaml"}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "app/v1",
		Kind:       "CasbinAuthzMiddleware",
	}
	m, err := api.ReferTypedObject[core.Middleware](server, ref)
	testutil.DiffError(t, nil, nil, err)

	handler := m.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"message":"allowed"}`))
	}))

	h := m.Middleware(handler)

	// Age>15 is allowed (input jwt.MapClaims)
	r1 := httptest.NewRequest(http.MethodGet, "http://casbin.com/test", nil).
		WithContext(context.WithValue(context.Background(), "AuthClaims", jwt.MapClaims{"age": 20}))
	w1 := httptest.NewRecorder()
	h.ServeHTTP(w1, r1)
	b1, _ := io.ReadAll(w1.Result().Body)
	testutil.Diff(t, http.StatusOK, w1.Result().StatusCode)
	testutil.Diff(t, `{"message":"allowed"}`, string(b1))

	// Age>15 is allowed (input map[string]any)
	r2 := httptest.NewRequest(http.MethodGet, "http://casbin.com/test", nil).
		WithContext(context.WithValue(context.Background(), "AuthClaims", map[string]any{"age": 20}))
	w2 := httptest.NewRecorder()
	h.ServeHTTP(w2, r2)
	b2, _ := io.ReadAll(w2.Result().Body)
	testutil.Diff(t, http.StatusOK, w2.Result().StatusCode)
	testutil.Diff(t, `{"message":"allowed"}`, string(b2))

	// POST is denied.
	r3 := httptest.NewRequest(http.MethodPost, "http://casbin.com/test", nil).
		WithContext(context.WithValue(context.Background(), "AuthClaims", map[string]any{"age": 20}))
	w3 := httptest.NewRecorder()
	h.ServeHTTP(w3, r3)
	b3, _ := io.ReadAll(w3.Result().Body)
	testutil.Diff(t, http.StatusForbidden, w3.Result().StatusCode)
	testutil.Diff(t, `{"status":403,"statusText":"Forbidden"}`, string(b3))

	// Age<=15 is denied
	r4 := httptest.NewRequest(http.MethodGet, "http://casbin.com/test", nil).
		WithContext(context.WithValue(context.Background(), "AuthClaims", map[string]any{"age": 15}))
	w4 := httptest.NewRecorder()
	h.ServeHTTP(w4, r4)
	b4, _ := io.ReadAll(w4.Result().Body)
	testutil.Diff(t, http.StatusForbidden, w4.Result().StatusCode)
	testutil.Diff(t, `{"status":403,"statusText":"Forbidden"}`, string(b4))

}

func TestAuthzAuthJSON(t *testing.T) {

	configs := []string{"./config-authz-auth-json.yaml"}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "app/v1",
		Kind:       "CasbinAuthzMiddleware",
	}
	m, err := api.ReferTypedObject[core.Middleware](server, ref)
	testutil.DiffError(t, nil, nil, err)

	handler := m.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"message":"allowed"}`))
	}))

	h := m.Middleware(handler)

	// Age>15 is allowed (input jwt.MapClaims)
	r1 := httptest.NewRequest(http.MethodGet, "http://casbin.com/test", nil).
		WithContext(context.WithValue(context.Background(), "AuthClaims", jwt.MapClaims{"age": 20}))
	w1 := httptest.NewRecorder()
	h.ServeHTTP(w1, r1)
	b1, _ := io.ReadAll(w1.Result().Body)
	testutil.Diff(t, http.StatusOK, w1.Result().StatusCode)
	testutil.Diff(t, `{"message":"allowed"}`, string(b1))

	// Age>15 is allowed (input map[string]any)
	r2 := httptest.NewRequest(http.MethodGet, "http://casbin.com/test", nil).
		WithContext(context.WithValue(context.Background(), "AuthClaims", map[string]any{"age": 20}))
	w2 := httptest.NewRecorder()
	h.ServeHTTP(w2, r2)
	b2, _ := io.ReadAll(w2.Result().Body)
	testutil.Diff(t, http.StatusOK, w2.Result().StatusCode)
	testutil.Diff(t, `{"message":"allowed"}`, string(b2))

	// POST is denied.
	r3 := httptest.NewRequest(http.MethodPost, "http://casbin.com/test", nil).
		WithContext(context.WithValue(context.Background(), "AuthClaims", map[string]any{"age": 20}))
	w3 := httptest.NewRecorder()
	h.ServeHTTP(w3, r3)
	b3, _ := io.ReadAll(w3.Result().Body)
	testutil.Diff(t, http.StatusForbidden, w3.Result().StatusCode)
	testutil.Diff(t, `{"status":403,"statusText":"Forbidden"}`, string(b3))

	// Age<=15 is denied
	r4 := httptest.NewRequest(http.MethodGet, "http://casbin.com/test", nil).
		WithContext(context.WithValue(context.Background(), "AuthClaims", map[string]any{"age": 15}))
	w4 := httptest.NewRecorder()
	h.ServeHTTP(w4, r4)
	b4, _ := io.ReadAll(w4.Result().Body)
	testutil.Diff(t, http.StatusForbidden, w4.Result().StatusCode)
	testutil.Diff(t, `{"status":403,"statusText":"Forbidden"}`, string(b4))

}

func TestAuthzAuthHTTP(t *testing.T) {

	configs := []string{"./config-authz-auth-http.yaml"}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	svr := &http.Server{ // Policy server.
		Addr: ":12121",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/csv")
			b, _ := os.ReadFile("./enforcer-policy-auth.csv")
			w.Write(b)
		}),
	}
	go func() { svr.ListenAndServe() }()
	defer svr.Close()
	time.Sleep(time.Second)

	ref := &kernel.Reference{
		APIVersion: "app/v1",
		Kind:       "CasbinAuthzMiddleware",
	}
	m, err := api.ReferTypedObject[core.Middleware](server, ref)
	testutil.DiffError(t, nil, nil, err)

	handler := m.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"message":"allowed"}`))
	}))

	h := m.Middleware(handler)

	// Age>15 is allowed (input jwt.MapClaims)
	r1 := httptest.NewRequest(http.MethodGet, "http://casbin.com/test", nil).
		WithContext(context.WithValue(context.Background(), "AuthClaims", jwt.MapClaims{"age": 20}))
	w1 := httptest.NewRecorder()
	h.ServeHTTP(w1, r1)
	b1, _ := io.ReadAll(w1.Result().Body)
	testutil.Diff(t, http.StatusOK, w1.Result().StatusCode)
	testutil.Diff(t, `{"message":"allowed"}`, string(b1))

	// Age>15 is allowed (input map[string]any)
	r2 := httptest.NewRequest(http.MethodGet, "http://casbin.com/test", nil).
		WithContext(context.WithValue(context.Background(), "AuthClaims", map[string]any{"age": 20}))
	w2 := httptest.NewRecorder()
	h.ServeHTTP(w2, r2)
	b2, _ := io.ReadAll(w2.Result().Body)
	testutil.Diff(t, http.StatusOK, w2.Result().StatusCode)
	testutil.Diff(t, `{"message":"allowed"}`, string(b2))

	// POST is denied.
	r3 := httptest.NewRequest(http.MethodPost, "http://casbin.com/test", nil).
		WithContext(context.WithValue(context.Background(), "AuthClaims", map[string]any{"age": 20}))
	w3 := httptest.NewRecorder()
	h.ServeHTTP(w3, r3)
	b3, _ := io.ReadAll(w3.Result().Body)
	testutil.Diff(t, http.StatusForbidden, w3.Result().StatusCode)
	testutil.Diff(t, `{"status":403,"statusText":"Forbidden"}`, string(b3))

	// Age<=15 is denied
	r4 := httptest.NewRequest(http.MethodGet, "http://casbin.com/test", nil).
		WithContext(context.WithValue(context.Background(), "AuthClaims", map[string]any{"age": 15}))
	w4 := httptest.NewRecorder()
	h.ServeHTTP(w4, r4)
	b4, _ := io.ReadAll(w4.Result().Body)
	testutil.Diff(t, http.StatusForbidden, w4.Result().StatusCode)
	testutil.Diff(t, `{"status":403,"statusText":"Forbidden"}`, string(b4))

}

func TestAuthzHeader(t *testing.T) {

	configs := []string{"./config-authz-header.yaml"}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "app/v1",
		Kind:       "CasbinAuthzMiddleware",
	}
	m, err := api.ReferTypedObject[core.Middleware](server, ref)
	testutil.DiffError(t, nil, nil, err)

	handler := m.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"message":"allowed"}`))
	}))

	h := m.Middleware(handler)

	// Header value ".*allowed.*" is allowed
	r1 := httptest.NewRequest(http.MethodGet, "http://casbin.com/test", nil)
	r1.Header.Set("Test", "this is allowed")
	w1 := httptest.NewRecorder()
	h.ServeHTTP(w1, r1)
	b1, _ := io.ReadAll(w1.Result().Body)
	testutil.Diff(t, http.StatusOK, w1.Result().StatusCode)
	testutil.Diff(t, `{"message":"allowed"}`, string(b1))

	// POST is denied.
	r2 := httptest.NewRequest(http.MethodPost, "http://casbin.com/test", nil)
	r2.Header.Set("Test", "this is allowed")
	w2 := httptest.NewRecorder()
	h.ServeHTTP(w2, r2)
	b2, _ := io.ReadAll(w2.Result().Body)
	testutil.Diff(t, http.StatusForbidden, w2.Result().StatusCode)
	testutil.Diff(t, `{"status":403,"statusText":"Forbidden"}`, string(b2))

	// Header value without "allowed" is denied.
	r3 := httptest.NewRequest(http.MethodGet, "http://casbin.com/test", nil)
	r3.Header.Set("Test", "this is denied")
	w3 := httptest.NewRecorder()
	h.ServeHTTP(w3, r3)
	b3, _ := io.ReadAll(w3.Result().Body)
	testutil.Diff(t, http.StatusForbidden, w3.Result().StatusCode)
	testutil.Diff(t, `{"status":403,"statusText":"Forbidden"}`, string(b3))

}

func TestAuthzHeaderMultiple(t *testing.T) {

	configs := []string{"./config-authz-header-multiple.yaml"}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "app/v1",
		Kind:       "CasbinAuthzMiddleware",
	}
	m, err := api.ReferTypedObject[core.Middleware](server, ref)
	testutil.DiffError(t, nil, nil, err)

	handler := m.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"message":"allowed"}`))
	}))

	h := m.Middleware(handler)

	// Header value ".*allowed.*" is allowed
	r1 := httptest.NewRequest(http.MethodGet, "http://casbin.com/test", nil)
	r1.Header["Test"] = []string{"foo", "bar", "allowed"}
	w1 := httptest.NewRecorder()
	h.ServeHTTP(w1, r1)
	b1, _ := io.ReadAll(w1.Result().Body)
	testutil.Diff(t, http.StatusOK, w1.Result().StatusCode)
	testutil.Diff(t, `{"message":"allowed"}`, string(b1))

	// POST is denied.
	r2 := httptest.NewRequest(http.MethodPost, "http://casbin.com/test", nil)
	r2.Header["Test"] = []string{"foo", "bar", "allowed"}
	w2 := httptest.NewRecorder()
	h.ServeHTTP(w2, r2)
	b2, _ := io.ReadAll(w2.Result().Body)
	testutil.Diff(t, http.StatusForbidden, w2.Result().StatusCode)
	testutil.Diff(t, `{"status":403,"statusText":"Forbidden"}`, string(b2))

	// Header value without "allowed" is denied.
	r3 := httptest.NewRequest(http.MethodGet, "http://casbin.com/test", nil)
	r3.Header["Test"] = []string{"foo", "bar"}
	w3 := httptest.NewRecorder()
	h.ServeHTTP(w3, r3)
	b3, _ := io.ReadAll(w3.Result().Body)
	testutil.Diff(t, http.StatusForbidden, w3.Result().StatusCode)
	testutil.Diff(t, `{"status":403,"statusText":"Forbidden"}`, string(b3))

}

func TestAuthzHost(t *testing.T) {

	configs := []string{"./config-authz-query.yaml"}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "app/v1",
		Kind:       "CasbinAuthzMiddleware",
	}
	m, err := api.ReferTypedObject[core.Middleware](server, ref)
	testutil.DiffError(t, nil, nil, err)

	handler := m.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"message":"allowed"}`))
	}))

	h := m.Middleware(handler)

	// Query value "this-is-allowed" is allowed
	r1 := httptest.NewRequest(http.MethodGet, "http://casbin.com/?test=this-is-allowed", nil)
	w1 := httptest.NewRecorder()
	h.ServeHTTP(w1, r1)
	b1, _ := io.ReadAll(w1.Result().Body)
	testutil.Diff(t, http.StatusOK, w1.Result().StatusCode)
	testutil.Diff(t, `{"message":"allowed"}`, string(b1))

	// POST is denied.
	r2 := httptest.NewRequest(http.MethodPost, "http://casbin.com/?test=this-is-allowed", nil)
	w2 := httptest.NewRecorder()
	h.ServeHTTP(w2, r2)
	b2, _ := io.ReadAll(w2.Result().Body)
	testutil.Diff(t, http.StatusForbidden, w2.Result().StatusCode)
	testutil.Diff(t, `{"status":403,"statusText":"Forbidden"}`, string(b2))

	// Query value "test=this-is-denied" is denied.
	r3 := httptest.NewRequest(http.MethodGet, "http://casbin.com/?test=this-is-denied", nil)
	w3 := httptest.NewRecorder()
	h.ServeHTTP(w3, r3)
	b3, _ := io.ReadAll(w3.Result().Body)
	testutil.Diff(t, http.StatusForbidden, w3.Result().StatusCode)
	testutil.Diff(t, `{"status":403,"statusText":"Forbidden"}`, string(b3))

}

func TestAuthzPolicyNotfound(t *testing.T) {

	configs := []string{"./config-authz-policy-notfound.yaml"}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "app/v1",
		Kind:       "CasbinAuthzMiddleware",
	}
	_, err = api.ReferTypedObject[core.Middleware](server, ref)
	errPattern := regexp.MustCompile(core.ErrPrefix + `failed to create CasbinAuthzMiddleware`)
	testutil.DiffError(t, core.ErrCoreGenCreateObject, errPattern, err)

}
