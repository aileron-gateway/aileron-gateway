package headercert

import (
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/aileron-gateway/aileron-gateway/core"
	"github.com/aileron-gateway/aileron-gateway/kernel/log"
)

type HeaderCertAuth struct {
	lg      log.Logger
	eh      core.ErrorHandler
	rootCAs []string
}

func (h *HeaderCertAuth) loadRootCert() (*x509.CertPool, error) {

	roots := x509.NewCertPool()

	// Read the root certificate specified in the local file
	for _, c := range h.rootCAs {
		rootCertPEM, err := os.ReadFile(c)
		if err != nil {
			return nil, err
		}
		// Add the root certificate to CertPool
		if !roots.AppendCertsFromPEM(rootCertPEM) {
			return nil, fmt.Errorf("failed to add root certificate to CertPool")
		}
	}

	return roots, nil
}

func (h *HeaderCertAuth) convertCert(certHeader string) (*x509.Certificate, error) {
	// Decode a Base64-encoded client certificate and convert it to a byte array
	decodedCertHeader, err := base64.StdEncoding.DecodeString(certHeader)
	if err != nil {
		return nil, fmt.Errorf("failed to decode certificate")
	}

	// Convert the client certificate into a PEM block
	certBlock, _ := pem.Decode(decodedCertHeader)
	if certBlock == nil {
		return nil, fmt.Errorf("failed to decode certificate")
	}

	// Analyze the PEM block
	cert, err := x509.ParseCertificate(certBlock.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse certificate")
	}

	return cert, nil
}

func (h *HeaderCertAuth) isFingerprintMatched(cert *x509.Certificate, fingerprintHeader string) bool {
	fingerprint := sha256.Sum256(cert.Raw)
	return hex.EncodeToString(fingerprint[:]) == fingerprintHeader
}

func isCertExpired(cert *x509.Certificate) bool {
	return time.Now().After(cert.NotAfter)
}

func (h *HeaderCertAuth) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.Header)

		certHeader := r.Header.Get("X-SSL-Client-Cert")
		if certHeader == "" {
			http.Error(w, "Client certificate not found", http.StatusBadRequest)
			return
		}

		fingerprintHeader := r.Header.Get("X-SSL-Client-Fingerprint")
		if certHeader == "" {
			http.Error(w, "Fingerprint is not found", http.StatusBadRequest)
			return
		}

		cert, err := h.convertCert(certHeader)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		roots, err := h.loadRootCert()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		opts := x509.VerifyOptions{
			Roots: roots,
		}

		// Verify the client certificate
		if _, err := cert.Verify(opts); err != nil {
			http.Error(w, "Fail to verify certificate", http.StatusUnauthorized)
			return
		}

		// Verify the fingerprint
		if !h.isFingerprintMatched(cert, fingerprintHeader) {
			http.Error(w, "Fail to verify fingerprint", http.StatusUnauthorized)
			return
		}

		// Check the expiration date
		if isCertExpired(cert) {
			http.Error(w, "Certificate has expired", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
