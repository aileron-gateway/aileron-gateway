package headercert

import (
	"crypto/x509"
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
	utilhttp "github.com/aileron-gateway/aileron-gateway/util/http"
)

const (
	certPath           = "../../../_example/header-cert/pki/client.crt"
	fpPath             = "../../../_example/header-cert/pki/fingerprint.txt"
	failCertPath       = "../../../_example/header-cert/pki/fail-client.crt"
	failFpPath         = "../../../_example/header-cert/pki/fail-fingerprint.txt"
	expiredCertPath    = "../../../_example/header-cert/pki/expired-client.crt"
	expiredFpPath      = "../../../_example/header-cert/pki/expired-fingerprint.txt"
	incompleteCertPath = "../../../_example/header-cert/pki/incomplete-client.crt"
	rootCAPath         = "../../../_example/header-cert/pki/rootCA.crt"
)

func TestMiddleware(t *testing.T) {

	type condition struct {
		method  string
		headers map[string]string
	}

	type action struct {
		status int
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	cert, _:= os.ReadFile(certPath)
	if err != nil {
		t.Errorf("fail to read client cert: %v", err)
	}

	fp, err := os.ReadFile(fpPath)
	if err != nil {
		t.Errorf("fail to read client fingerprint: %v", err)
	}

	failCert, err := os.ReadFile(failCertPath)
	if err != nil {
		t.Errorf("fail to read client cert: %v", err)
	}

	failFp, err := os.ReadFile(failFpPath)
	if err != nil {
		t.Errorf("fail to read client fingerprint: %v", err)
	}

	expiredCert, err := os.ReadFile(expiredCertPath)
	if err != nil {
		t.Errorf("fail to read client cert: %v", err)
	}

	expiredFp, err := os.ReadFile(expiredFpPath)
	if err != nil {
		t.Errorf("fail to read client fingerprint: %v", err)
	}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"valid client cert and fingerprint",
			[]string{},
			[]string{},
			&condition{
				method: http.MethodGet,
				headers: map[string]string{
					"X-SSL-Client-Cert":        base64.StdEncoding.EncodeToString(cert),
					"X-SSL-Client-Fingerprint": string(fp),
				},
			},
			&action{
				status: http.StatusOK,
			},
		),
		gen(
			"no client cert",
			[]string{},
			[]string{},
			&condition{
				method: http.MethodGet,
				headers: map[string]string{
					"X-SSL-Client-Cert":        "", // no client cert
					"X-SSL-Client-Fingerprint": string(fp),
				},
			},
			&action{
				status: http.StatusBadRequest,
			},
		),
		gen(
			"no fingerprint",
			[]string{},
			[]string{},
			&condition{
				method: http.MethodGet,
				headers: map[string]string{
					"X-SSL-Client-Cert":        base64.StdEncoding.EncodeToString(cert),
					"X-SSL-Client-Fingerprint": "", // no fingerprint
				},
			},
			&action{
				status: http.StatusOK,
			},
		),
		gen(
			"invalid client cert",
			[]string{},
			[]string{},
			&condition{
				method: http.MethodGet,
				headers: map[string]string{
					"X-SSL-Client-Cert":        "invalid cert", // invalid client cert
					"X-SSL-Client-Fingerprint": "654ed0fbb21c25ce32e6cd64846af842e6e821eae3ed4b32a16a164afaf10226",
				},
			},
			&action{
				status: http.StatusBadRequest,
			},
		),
		gen(
			"fail to verify the client cert",
			[]string{},
			[]string{},
			&condition{
				method: http.MethodGet,
				headers: map[string]string{
					"X-SSL-Client-Cert":        base64.StdEncoding.EncodeToString(failCert), // created by another rootCA.
					"X-SSL-Client-Fingerprint": string(failFp),
				},
			},
			&action{
				status: http.StatusUnauthorized,
			},
		),
		gen(
			"fingerprint not matched",
			[]string{},
			[]string{},
			&condition{
				method: http.MethodGet,
				headers: map[string]string{
					"X-SSL-Client-Cert":        base64.StdEncoding.EncodeToString(cert),
					"X-SSL-Client-Fingerprint": "invalid fingerprint",
				},
			},
			&action{
				status: http.StatusUnauthorized,
			},
		),
		gen(
			"expired client cert",
			[]string{},
			[]string{},
			&condition{
				method: http.MethodGet,
				headers: map[string]string{
					"X-SSL-Client-Cert":        base64.StdEncoding.EncodeToString(expiredCert), // expired client cert
					"X-SSL-Client-Fingerprint": string(expiredFp),
				},
			},
			&action{
				status: http.StatusUnauthorized,
			},
		),
	}

	testutil.Register(table, testCases...)

	rootCAs := []string{rootCAPath}
	roots := x509.NewCertPool()

	// Read the root cert specified in the local file
	for _, c := range rootCAs {
		rootCertPEM, err := os.ReadFile(c)
		if err != nil {
			t.Errorf("fail to load rootCA: %v", err)
		}
		// Add the root cert to CertPool
		if !roots.AppendCertsFromPEM(rootCertPEM) {
			t.Errorf("fail to load rootCA: %v", err)
		}
	}

	opts := x509.VerifyOptions{
		Roots: roots,
	}

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {

			// Prepare the headercert middleware
			headerCertMiddleware := &headerCert{
				eh:   utilhttp.GlobalErrorHandler(utilhttp.DefaultErrorHandlerName),
				opts: opts,
			}

			// Create a test request
			req := httptest.NewRequest(tt.C().method, "http://test.com", nil)
			for k, v := range tt.C().headers {
				req.Header.Set(k, v)
			}

			// Create a test respose recoder
			resp := httptest.NewRecorder()

			// Call the middleware
			headerCertMiddleware.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})).ServeHTTP(resp, req)

			// verify the status code
			testutil.Diff(t, tt.A().status, resp.Code)
		})
	}
}

func TestConvertCert(t *testing.T) {

	t.Run("invalid pem", func(t *testing.T) {
		invalidPEM := base64.StdEncoding.EncodeToString([]byte("-----BEGIN cert-----\nInvalid cert content"))
		_, err := parseCert(invalidPEM)
		if err == nil {
			t.Errorf("expected error, got nil")
		}
	})

	invalidCert, err := os.ReadFile(incompleteCertPath)
	if err != nil {
		t.Errorf("fail to read client cert: %v", err)
	}
	t.Run("invalid x509 cert", func(t *testing.T) {

		invalidX509 := base64.StdEncoding.EncodeToString([]byte(invalidCert)) // This cert has a negative serial number that causes an error in x509.Parsecert().
		_, err := parseCert(invalidX509)
		if err == nil {
			t.Errorf("expected error, got nil")
		}
	})
}
