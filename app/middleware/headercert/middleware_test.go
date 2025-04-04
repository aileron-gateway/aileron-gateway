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
	failCertPath       = "../../../test/ut/app/headercert/cert/fail-client.crt"
	failFpPath         = "../../../test/ut/app/headercert/fingerprint/fail-fingerprint.txt"
	expiredCertPath    = "../../../test/ut/app/headercert/cert/expired-client.crt"
	expiredFpPath      = "../../../test/ut/app/headercert/fingerprint/expired-fingerprint.txt"
	incompleteCertPath = "../../../test/ut/app/headercert/cert/incomplete-client.crt"
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

	cert, _ := os.ReadFile(certPath)
	fp, _ := os.ReadFile(fpPath)
	failCert, _ := os.ReadFile(failCertPath)
	failFp, _ := os.ReadFile(failFpPath)
	expiredCert, _ := os.ReadFile(expiredCertPath)
	expiredFp, _ := os.ReadFile(expiredFpPath)

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"valid client cert and fingerprint",
			[]string{},
			[]string{},
			&condition{
				method: http.MethodGet,
				headers: map[string]string{
					"X-SSL-Client-Cert":        base64.URLEncoding.EncodeToString(cert),
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
					"X-SSL-Client-Cert":        base64.URLEncoding.EncodeToString(cert),
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
					"X-SSL-Client-Cert":        base64.URLEncoding.EncodeToString(failCert), // created by another rootCA.
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
					"X-SSL-Client-Cert":        base64.URLEncoding.EncodeToString(cert),
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
					"X-SSL-Client-Cert":        base64.URLEncoding.EncodeToString(expiredCert), // expired client cert
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
		rootCertPEM, _ := os.ReadFile(c)
		// Add the root cert to CertPool
		roots.AppendCertsFromPEM(rootCertPEM)
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

func TestParseCert(t *testing.T) {
	t.Run("invalid pem", func(t *testing.T) {
		invalidPEM := base64.URLEncoding.EncodeToString([]byte("-----BEGIN cert-----\nInvalid cert content"))
		cert, err := parseCert(invalidPEM)
		if cert != nil {
			t.Errorf("expected nil cert, got %v", cert)
		}
		want := "E3214.AppMiddleInvalidCert client certificate invalid or not found. pem not found"
		if want != err.Error() {
			t.Errorf("expected %v, got %v", want, err)
		}
	})
	
	t.Run("invalid x509 cert", func(t *testing.T) {
		invalidCert, _ := os.ReadFile(incompleteCertPath)
		invalidX509 := base64.URLEncoding.EncodeToString([]byte(invalidCert)) // This cert has a negative serial number that causes an error in [x509.ParseCertificate].
		cert, err := parseCert(invalidX509)
		if cert != nil {
			t.Errorf("expected nil cert, got %v", cert)
		}
		want := "E3214.AppMiddleInvalidCert client certificate invalid or not found. x509 parse failed [x509: negative serial number]"
		if want != err.Error() {
			t.Errorf("expected %v, got %v", want, err)
		}
	})
}
