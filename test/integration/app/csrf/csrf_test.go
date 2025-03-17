//go:build integration

package csrf_test

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aileron-gateway/aileron-gateway/apis/kernel"
	"github.com/aileron-gateway/aileron-gateway/cmd/aileron/app"
	"github.com/aileron-gateway/aileron-gateway/core"
	"github.com/aileron-gateway/aileron-gateway/kernel/api"
	"github.com/aileron-gateway/aileron-gateway/kernel/mac"
	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
	"github.com/aileron-gateway/aileron-gateway/test/integration/common"
	"github.com/aileron-gateway/aileron-gateway/util/session"
)

func generateValidToken(t *testing.T, secret string, seedSize, hashSize int, hmacFunc func([]byte, []byte) []byte) string {
	// create seeds
	seed := make([]byte, seedSize)
	_, err := rand.Read(seed)
	testutil.DiffError(t, nil, nil, err)

	// Digest generation with HMAC
	digest := hmacFunc(seed, []byte(secret))
	testutil.Diff(t, len(digest), hashSize)

	// Seeds and digests are combined to create tokens
	tokenBytes := append(seed, digest...)

	// Return tokens encoded as hexadecimal strings
	return hex.EncodeToString(tokenBytes)
}

func TestCustomHeader(t *testing.T) {
	configs := []string{"./config-custom-header.yaml"}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "app/v1",
		Kind:       "CSRFMiddleware",
		Name:       "",
		Namespace:  "",
	}
	m, err := api.ReferTypedObject[core.Middleware](server, ref)
	testutil.DiffError(t, nil, nil, err)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})
	h := m.Middleware(handler)

	// allowed custom header value
	r1 := httptest.NewRequest(http.MethodPost, "http://csrf.com/custom-header", nil)
	r1.Header.Set("X-CSRF-With", "ValidToken")
	w1 := httptest.NewRecorder()
	h.ServeHTTP(w1, r1)
	testutil.Diff(t, http.StatusOK, w1.Result().StatusCode)
	testutil.Diff(t, "ok", w1.Body.String())

	// disallowed custom header value
	r2 := httptest.NewRequest(http.MethodPost, "http://csrf.com/custom-header", nil)
	r2.Header.Set("X-CSRF-With", "InvalidToken")
	w2 := httptest.NewRecorder()
	h.ServeHTTP(w2, r2)
	testutil.Diff(t, http.StatusForbidden, w2.Result().StatusCode)
	b2, _ := io.ReadAll(w2.Result().Body)
	testutil.Diff(t, `{"status":403,"statusText":"Forbidden"}`, string(b2))
}

func TestDoubleSubmitHeader(t *testing.T) {
	configs := []string{"./config-double-submit-header.yaml"}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "app/v1",
		Kind:       "CSRFMiddleware",
	}
	m, err := api.ReferTypedObject[core.Middleware](server, ref)
	testutil.DiffError(t, nil, nil, err)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})
	h := m.Middleware(handler)

	// Test valid CSRF token in cookie and header
	secret := "bXktc2VjdXJlLXNlY3JldC1rZXktc3RyaW5nLWZvci1jc3Jm" // Set the secret defined in the config file
	decodedSecret, _ := base64.StdEncoding.DecodeString(secret)
	seedSize := 32
	hashSize := 32
	hmacFunc := mac.FromHashAlg(kernel.HashAlg_SHA256)
	validToken := generateValidToken(t, string(decodedSecret), seedSize, hashSize, hmacFunc)

	// Valid token in cookie and header.
	r1 := httptest.NewRequest(http.MethodPost, "http://csrf.com/double-submit", nil)
	r1.AddCookie(&http.Cookie{Name: "__csrfToken", Value: validToken})
	r1.Header.Set("X-CSRF-Token", validToken)
	w1 := httptest.NewRecorder()
	h.ServeHTTP(w1, r1)
	testutil.Diff(t, http.StatusOK, w1.Result().StatusCode)
	testutil.Diff(t, "ok", w1.Body.String())

	// Valid token in cookie and invalid token in header.
	r2 := httptest.NewRequest(http.MethodPost, "http://csrf.com/double-submit", nil)
	r2.AddCookie(&http.Cookie{Name: "__csrfToken", Value: validToken})
	r2.Header.Set("X-CSRF-Token", "invalidToken") // Invalid token.
	w2 := httptest.NewRecorder()
	h.ServeHTTP(w2, r2)
	testutil.Diff(t, http.StatusForbidden, w2.Result().StatusCode)
	testutil.Diff(t, `{"status":403,"statusText":"Forbidden"}`, w2.Body.String())

	// CSRF token missing
	r3 := httptest.NewRequest(http.MethodPost, "http://csrf.com/double-submit", nil)
	w3 := httptest.NewRecorder()
	h.ServeHTTP(w3, r3)
	testutil.Diff(t, http.StatusForbidden, w3.Result().StatusCode)
	testutil.Diff(t, `{"status":403,"statusText":"Forbidden"}`, w3.Body.String())
}

func TestDoubleSubmitForm(t *testing.T) {
	configs := []string{"./config-double-submit-form.yaml"}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "app/v1",
		Kind:       "CSRFMiddleware",
	}
	m, err := api.ReferTypedObject[core.Middleware](server, ref)
	testutil.DiffError(t, nil, nil, err)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})
	h := m.Middleware(handler)

	// Test valid CSRF token in cookie and header
	secret := "bXktc2VjdXJlLXNlY3JldC1rZXktc3RyaW5nLWZvci1jc3Jm" // Set the secret defined in the config file
	decodedSecret, _ := base64.StdEncoding.DecodeString(secret)
	seedSize := 32
	hashSize := 32
	hmacFunc := mac.FromHashAlg(kernel.HashAlg_SHA256)
	validToken := generateValidToken(t, string(decodedSecret), seedSize, hashSize, hmacFunc)

	// Valid token in cookie and header.
	b1 := bytes.NewReader([]byte("X-CSRF-Token=" + validToken))
	r1 := httptest.NewRequest(http.MethodPost, "http://csrf.com/double-submit", b1)
	r1.AddCookie(&http.Cookie{Name: "__csrfToken", Value: validToken})
	r1.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w1 := httptest.NewRecorder()
	h.ServeHTTP(w1, r1)
	testutil.Diff(t, http.StatusOK, w1.Result().StatusCode)
	testutil.Diff(t, "ok", w1.Body.String())

	// Valid token in cookie and invalid token in header.
	b2 := bytes.NewReader([]byte("X-CSRF-Token=" + "invalidToken"))
	r2 := httptest.NewRequest(http.MethodPost, "http://csrf.com/double-submit", b2)
	r2.AddCookie(&http.Cookie{Name: "__csrfToken", Value: validToken})
	r2.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w2 := httptest.NewRecorder()
	h.ServeHTTP(w2, r2)
	testutil.Diff(t, http.StatusForbidden, w2.Result().StatusCode)
	testutil.Diff(t, `{"status":403,"statusText":"Forbidden"}`, w2.Body.String())

	// CSRF token missing
	r3 := httptest.NewRequest(http.MethodPost, "http://csrf.com/double-submit", nil)
	w3 := httptest.NewRecorder()
	h.ServeHTTP(w3, r3)
	testutil.Diff(t, http.StatusForbidden, w3.Result().StatusCode)
	testutil.Diff(t, `{"status":403,"statusText":"Forbidden"}`, w3.Body.String())
}

func TestDoubleSubmitJSON(t *testing.T) {
	configs := []string{"./config-double-submit-json.yaml"}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "app/v1",
		Kind:       "CSRFMiddleware",
	}
	m, err := api.ReferTypedObject[core.Middleware](server, ref)
	testutil.DiffError(t, nil, nil, err)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})
	h := m.Middleware(handler)

	// Test valid CSRF token in cookie and header
	secret := "bXktc2VjdXJlLXNlY3JldC1rZXktc3RyaW5nLWZvci1jc3Jm" // Set the secret defined in the config file
	decodedSecret, _ := base64.StdEncoding.DecodeString(secret)
	seedSize := 32
	hashSize := 32
	hmacFunc := mac.FromHashAlg(kernel.HashAlg_SHA256)
	validToken := generateValidToken(t, string(decodedSecret), seedSize, hashSize, hmacFunc)

	// Valid token in cookie and header.
	b1 := bytes.NewReader([]byte(fmt.Sprintf(`{"X-CSRF-Token":"%s"`, validToken)))
	r1 := httptest.NewRequest(http.MethodPost, "http://csrf.com/double-submit", b1)
	r1.AddCookie(&http.Cookie{Name: "__csrfToken", Value: validToken})
	r1.Header.Set("Content-Type", "application/json")
	w1 := httptest.NewRecorder()
	h.ServeHTTP(w1, r1)
	testutil.Diff(t, http.StatusOK, w1.Result().StatusCode)
	testutil.Diff(t, "ok", w1.Body.String())

	// Valid token in cookie and invalid token in header.
	b2 := bytes.NewReader([]byte(fmt.Sprintf(`{"X-CSRF-Token":"%s"`, "invalidToken")))
	r2 := httptest.NewRequest(http.MethodPost, "http://csrf.com/double-submit", b2)
	r2.AddCookie(&http.Cookie{Name: "__csrfToken", Value: validToken})
	r2.Header.Set("Content-Type", "application/json")
	w2 := httptest.NewRecorder()
	h.ServeHTTP(w2, r2)
	testutil.Diff(t, http.StatusForbidden, w2.Result().StatusCode)
	testutil.Diff(t, `{"status":403,"statusText":"Forbidden"}`, w2.Body.String())

	// CSRF token missing
	r3 := httptest.NewRequest(http.MethodPost, "http://csrf.com/double-submit", nil)
	w3 := httptest.NewRecorder()
	h.ServeHTTP(w3, r3)
	testutil.Diff(t, http.StatusForbidden, w3.Result().StatusCode)
	testutil.Diff(t, `{"status":403,"statusText":"Forbidden"}`, w3.Body.String())
}

func withSession(ctx context.Context, token string) context.Context {
	ss := session.NewDefaultSession(session.SerializeJSON)
	_ = ss.Persist("__csrf_token__", []byte(token))
	return session.ContextWithSession(ctx, ss)
}

func TestSynchronizerToken(t *testing.T) {
	configs := []string{"./config-synchronizer-token.yaml"}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "app/v1",
		Kind:       "CSRFMiddleware",
	}
	m, err := api.ReferTypedObject[core.Middleware](server, ref)
	testutil.DiffError(t, nil, nil, err)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})
	h := m.Middleware(handler)

	// request valid CSRF token
	secret := "bXktc2VjdXJlLXNlY3JldC1rZXktc3RyaW5nLWZvci1jc3Jm" // Set the secret defined in the config file
	decodedSecret, _ := base64.StdEncoding.DecodeString(secret)
	seedSize := 32
	hashSize := 32
	hmacFunc := mac.FromHashAlg(kernel.HashAlg_SHA256)
	validToken := generateValidToken(t, string(decodedSecret), seedSize, hashSize, hmacFunc)

	// request valid CSRF token
	r1 := httptest.NewRequest(http.MethodPost, "http://csrf.com/synchronizer-token", nil)
	r1 = r1.WithContext(withSession(r1.Context(), validToken))
	r1.Header.Set("X-CSRF-Token", validToken)
	w1 := httptest.NewRecorder()
	h.ServeHTTP(w1, r1)
	testutil.Diff(t, http.StatusOK, w1.Result().StatusCode)
	testutil.Diff(t, "ok", w1.Body.String())

	// request invalid CSRF token
	r2 := httptest.NewRequest(http.MethodPost, "http://csrf.com/synchronizer-token", nil)
	r2 = r2.WithContext(withSession(r2.Context(), validToken))
	r2.Header.Set("X-CSRF-Token", "InvalidToken")
	w2 := httptest.NewRecorder()
	h.ServeHTTP(w2, r2)
	testutil.Diff(t, http.StatusForbidden, w2.Result().StatusCode)
	b2, _ := io.ReadAll(w2.Result().Body)
	testutil.Diff(t, `{"status":403,"statusText":"Forbidden"}`, string(b2))
}
