package headercert

import (
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/aileron-gateway/aileron-gateway/apis/kernel"
	"github.com/aileron-gateway/aileron-gateway/kernel/log"
	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
	utilhttp "github.com/aileron-gateway/aileron-gateway/util/http"
)

const (
	certPath        = "../../../_example/header-cert/pki/client.crt"
	fpPath          = "../../../_example/header-cert/pki/fingerprint.txt"
	failCertPath    = "../../../_example/header-cert/pki/fail-client.crt"
	failFpPath      = "../../../_example/header-cert/pki/fail-fingerprint.txt"
	expiredCertPath = "../../../_example/header-cert/pki/expired-client.crt"
	expiredFpPath   = "../../../_example/header-cert/pki/expired-fingerprint.txt"
	incompleteCertPath = "../../../_example/header-cert/pki/incomplete-client.crt"
	rootCAPath      = "../../../_example/header-cert/pki/rootCA.crt"
)

func TestMiddleware(t *testing.T) {

	type condition struct {
		method  string
		headers map[string]string
		tls     *kernel.TLSConfig
	}

	type action struct {
		status int
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	cert, err := os.ReadFile(certPath)
	if err != nil {
		t.Fatalf("fail to read client certificate: %v", err)
	}
	encodedCert := base64.StdEncoding.EncodeToString(cert)

	fp, err := os.ReadFile(fpPath)
	if err != nil {
		t.Fatalf("fail to read client fingerprint: %v", err)
	}
	fingerprint := string(fp)

	failCert, err := os.ReadFile(failCertPath)
	if err != nil {
		t.Fatalf("fail to read client certificate: %v", err)
	}
	encodedFailCert := base64.StdEncoding.EncodeToString(failCert)

	failFp, err := os.ReadFile(failFpPath)
	if err != nil {
		t.Fatalf("fail to read client fingerprint: %v", err)
	}
	failFingerprint := string(failFp)

	expiredCert, err := os.ReadFile(expiredCertPath)
	if err != nil {
		t.Fatalf("fail to read client certificate: %v", err)
	}
	encodedExpiredCert := base64.StdEncoding.EncodeToString(expiredCert)

	expiredFp, err := os.ReadFile(expiredFpPath)
	if err != nil {
		t.Fatalf("fail to read client fingerprint: %v", err)
	}
	expiredFingerprint := string(expiredFp)

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"valid client certificate and fingerprint",
			[]string{},
			[]string{},
			&condition{
				method: http.MethodGet,
				headers: map[string]string{
					"X-SSL-Client-Cert":        encodedCert,
					"X-SSL-Client-Fingerprint": fingerprint,
				},
				tls: &kernel.TLSConfig{
					RootCAs: []string{rootCAPath},
				},
			},
			&action{
				status: http.StatusOK,
			},
		),
		gen(
			"no client certificate",
			[]string{},
			[]string{},
			&condition{
				method: http.MethodGet,
				headers: map[string]string{
					"X-SSL-Client-Cert":        "", // no client certificate
					"X-SSL-Client-Fingerprint": fingerprint,
				},
				tls: &kernel.TLSConfig{
					RootCAs: []string{rootCAPath},
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
					"X-SSL-Client-Cert":        encodedCert,
					"X-SSL-Client-Fingerprint": "", // no fingerprint
				},
				tls: &kernel.TLSConfig{
					RootCAs: []string{rootCAPath},
				},
			},
			&action{
				status: http.StatusBadRequest,
			},
		),
		gen(
			"invalid client certificate",
			[]string{},
			[]string{},
			&condition{
				method: http.MethodGet,
				headers: map[string]string{
					"X-SSL-Client-Cert":        "invalid certificate",
					"X-SSL-Client-Fingerprint": "654ed0fbb21c25ce32e6cd64846af842e6e821eae3ed4b32a16a164afaf10226",
				},
				tls: &kernel.TLSConfig{
					RootCAs: []string{rootCAPath},
				},
			},
			&action{
				status: http.StatusBadRequest,
			},
		),
		gen(
			"fail to verify the client certificate",
			[]string{},
			[]string{},
			&condition{
				method: http.MethodGet,
				headers: map[string]string{
					"X-SSL-Client-Cert":        encodedFailCert, // This certificate was created by another rootCA.
					"X-SSL-Client-Fingerprint": failFingerprint,
				},
				tls: &kernel.TLSConfig{
					RootCAs: []string{rootCAPath},
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
					"X-SSL-Client-Cert":        encodedCert,
					"X-SSL-Client-Fingerprint": "invalid fingerprint",
				},
				tls: &kernel.TLSConfig{
					RootCAs: []string{rootCAPath},
				},
			},
			&action{
				status: http.StatusUnauthorized,
			},
		),
		gen(
			"expired client certificate",
			[]string{},
			[]string{},
			&condition{
				method: http.MethodGet,
				headers: map[string]string{
					"X-SSL-Client-Cert":        encodedExpiredCert, // This certificate has expired.
					"X-SSL-Client-Fingerprint": expiredFingerprint,
				},
				tls: &kernel.TLSConfig{
					RootCAs: []string{rootCAPath},
				},
			},
			&action{
				status: http.StatusUnauthorized,
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {

			// Prepare the headercert middleware
			headerCertMiddleware := &headerCert{
				lg:      log.GlobalLogger(log.DefaultLoggerName),
				eh:      utilhttp.GlobalErrorHandler(utilhttp.DefaultErrorHandlerName),
				rootCAs: tt.C().tls.RootCAs,
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

	m := &headerCert{}

	t.Run("invalid pem", func(t *testing.T) {
		invalidPEM := base64.StdEncoding.EncodeToString([]byte("-----BEGIN CERTIFICATE-----\nInvalid certificate content"))
		_, err := m.convertCert(invalidPEM)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	
	invalidCert, err := os.ReadFile(incompleteCertPath)
	if err != nil {
		t.Fatalf("fail to read client certificate: %v", err)
	}
	t.Run("invalid x509 certificate", func(t *testing.T) {

		invalidX509 := base64.StdEncoding.EncodeToString([]byte(invalidCert)) // This certificate has a negative serial number that causes an error in x509.ParseCertificate().
		_, err := m.convertCert(invalidX509)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}
