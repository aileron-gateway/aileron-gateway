// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package headercert

import (
	"crypto/x509"
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/aileron-gateway/aileron-gateway/app"
	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
	utilhttp "github.com/aileron-gateway/aileron-gateway/util/http"
)

const (
	testDataDir        = "../../../test/ut/app/headercert/"
	certPath           = testDataDir + "client.crt"
	fpPath             = testDataDir + "fingerprint.txt"
	failCertPath       = testDataDir + "fail-client.crt"
	failFpPath         = testDataDir + "fail-fingerprint.txt"
	expiredCertPath    = testDataDir + "expired-client.crt"
	expiredFpPath      = testDataDir + "expired-fingerprint.txt"
	incompleteCertPath = testDataDir + "incomplete-client.crt"
	rootCAPath         = testDataDir + "rootCA.crt"
)

func TestMiddleware(t *testing.T) {
	type condition struct {
		method     string
		certHeader string
		fpHeader   string
		headers    map[string]string
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
				method:     http.MethodGet,
				certHeader: "X-SSL-Client-Cert",
				fpHeader:   "X-SSL-Client-Fingerprint",
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
				method:   http.MethodGet,
				fpHeader: "X-SSL-Client-Fingerprint",
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
			"no necessary fingerprint",
			[]string{},
			[]string{},
			&condition{
				method:     http.MethodGet,
				certHeader: "X-SSL-Client-Cert",
				fpHeader:   "X-SSL-Client-Fingerprint",
				headers: map[string]string{
					"X-SSL-Client-Cert":        base64.URLEncoding.EncodeToString(cert),
					"X-SSL-Client-Fingerprint": "", // no fingerprint
				},
			},
			&action{
				status: http.StatusUnauthorized,
			},
		),
		gen(
			"no unnecessary fingerprint",
			[]string{},
			[]string{},
			&condition{
				method:     http.MethodGet,
				certHeader: "X-SSL-Client-Cert",
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
				method:     http.MethodGet,
				certHeader: "X-SSL-Client-Cert",
				fpHeader:   "X-SSL-Client-Fingerprint",
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
				method:     http.MethodGet,
				certHeader: "X-SSL-Client-Cert",
				fpHeader:   "X-SSL-Client-Fingerprint",
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
				method:     http.MethodGet,
				certHeader: "X-SSL-Client-Cert",
				fpHeader:   "X-SSL-Client-Fingerprint",
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
				method:     http.MethodGet,
				certHeader: "X-SSL-Client-Cert",
				fpHeader:   "X-SSL-Client-Fingerprint",
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
				eh:         utilhttp.GlobalErrorHandler(utilhttp.DefaultErrorHandlerName),
				opts:       opts,
				certHeader: tt.C().certHeader,
				fpHeader:   tt.C().fpHeader,
			}

			// Create a test request
			req := httptest.NewRequest(tt.C().method, "http://test.com", nil)
			for k, v := range tt.C().headers {
				req.Header.Set(k, v)
			}

			// Create a test response recorder.
			resp := httptest.NewRecorder()

			// Call the middleware
			headerCertMiddleware.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})).ServeHTTP(resp, req)

			// Verify the status code
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
		if !app.ErrAppMiddleInvalidCert.Is(err) {
			t.Errorf("expected ErrAppMiddleInvalidCert, got %v", err)
		}
	})
	t.Run("invalid x509 cert", func(t *testing.T) {
		invalidCert, _ := os.ReadFile(incompleteCertPath)
		invalidX509 := base64.URLEncoding.EncodeToString([]byte(invalidCert)) // This cert has a negative serial number that causes an error in [x509.ParseCertificate].
		cert, err := parseCert(invalidX509)
		if cert != nil {
			t.Errorf("expected nil cert, got %v", cert)
		}
		if !app.ErrAppMiddleInvalidCert.Is(err) {
			t.Errorf("expected ErrAppMiddleInvalidCert, got %v", err)
		}
	})
}
