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
	"errors"

	"github.com/aileron-gateway/aileron-gateway/app"
	"github.com/aileron-gateway/aileron-gateway/core"
	"github.com/aileron-gateway/aileron-gateway/kernel/log"
	utilhttp "github.com/aileron-gateway/aileron-gateway/util/http"
)

type HeaderCert struct {
	lg      log.Logger
	eh      core.ErrorHandler
	rootCAs []string
}

func (m *HeaderCert) loadRootCert() (*x509.CertPool, error) {

	roots := x509.NewCertPool()

	// Read the root certificate specified in the local file
	for _, c := range m.rootCAs {
		rootCertPEM, err := os.ReadFile(c)
		if err != nil {
			return nil, err
		}
		// Add the root certificate to CertPool
		if !roots.AppendCertsFromPEM(rootCertPEM) {
			return nil, errors.New("failed to add root certificate to CertPool")
		}
	}

	return roots, nil
}

func (m *HeaderCert) convertCert(certHeader string) (*x509.Certificate, error) {
	// Decode a Base64-encoded client certificate and convert it to a byte array
	decodedCertHeader, err := base64.StdEncoding.DecodeString(certHeader)
	if err != nil {
		return nil, errors.New("failed to decode certificate")
	}

	// Convert the client certificate into a PEM block
	certBlock, _ := pem.Decode(decodedCertHeader)
	if certBlock == nil {
		return nil, errors.New("failed to decode certificate")
	}

	// Analyze the PEM block
	cert, err := x509.ParseCertificate(certBlock.Bytes)
	if err != nil {
		return nil, errors.New("failed to parse certificate")
	}

	return cert, nil
}

func (m *HeaderCert) isFingerprintMatched(cert *x509.Certificate, fingerprintHeader string) bool {
	fingerprint := sha256.Sum256(cert.Raw)
	return hex.EncodeToString(fingerprint[:]) == fingerprintHeader
}

func isCertExpired(cert *x509.Certificate) bool {
	return time.Now().After(cert.NotAfter)
}

func (m *HeaderCert) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.Header)

		certHeader := r.Header.Get("X-SSL-Client-Cert")
		if certHeader == "" {
			err := app.ErrAppMiddleHeaderPolicy.WithoutStack(nil, map[string]any{"reason": "client certificate is not found"})
			m.eh.ServeHTTPError(w, r, utilhttp.NewHTTPError(err, http.StatusBadRequest))
			return
		}
		
		fingerprintHeader := r.Header.Get("X-SSL-Client-Fingerprint")
		if certHeader == "" {
			err := app.ErrAppMiddleHeaderPolicy.WithoutStack(nil, map[string]any{"reason": "fingerprint is not found"})
			m.eh.ServeHTTPError(w, r, utilhttp.NewHTTPError(err, http.StatusBadRequest))
			return
		}
		
		cert, err := m.convertCert(certHeader)
		if err != nil {
			err := app.ErrAppMiddleHeaderPolicy.WithoutStack(nil, map[string]any{"reason": err.Error()})
			m.eh.ServeHTTPError(w, r, utilhttp.NewHTTPError(err, http.StatusBadRequest))
			return
		}
		
		roots, err := m.loadRootCert()
		if err != nil {
			err := app.ErrAppMiddleHeaderPolicy.WithoutStack(nil, map[string]any{"reason": err.Error()})
			m.eh.ServeHTTPError(w, r, utilhttp.NewHTTPError(err, http.StatusInternalServerError))
			return
		}
		
		opts := x509.VerifyOptions{
			Roots: roots,
		}
		
		// Verify the client certificate
		if _, err := cert.Verify(opts); err != nil {
			err := app.ErrAppMiddleHeaderPolicy.WithoutStack(nil, map[string]any{"reason": "fail to verify certificate"})
			m.eh.ServeHTTPError(w, r, utilhttp.NewHTTPError(err, http.StatusUnauthorized))
			return
		}
		
		// Verify the fingerprint
		if !m.isFingerprintMatched(cert, fingerprintHeader) {
			err := app.ErrAppMiddleHeaderPolicy.WithoutStack(nil, map[string]any{"reason": "fail to verify fingerprint"})
			m.eh.ServeHTTPError(w, r, utilhttp.NewHTTPError(err, http.StatusUnauthorized))
			return
		}
		
		// Check the expiration date
		if isCertExpired(cert) {
			err := app.ErrAppMiddleHeaderPolicy.WithoutStack(nil, map[string]any{"reason": "certificate has expired"})
			m.eh.ServeHTTPError(w, r, utilhttp.NewHTTPError(err, http.StatusUnauthorized))
			return
		}

		next.ServeHTTP(w, r)
	})
}
