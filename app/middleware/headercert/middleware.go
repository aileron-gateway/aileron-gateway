package headercert

import (
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"net/http"

	"github.com/aileron-gateway/aileron-gateway/app"
	"github.com/aileron-gateway/aileron-gateway/core"
	utilhttp "github.com/aileron-gateway/aileron-gateway/util/http"
)

type headerCert struct {
	eh   core.ErrorHandler
	opts x509.VerifyOptions
}

func (m *headerCert) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.Header)

		ch := r.Header.Get("X-SSL-Client-Cert")
		if ch == "" {
			err := app.ErrAppMiddleInvalidCert.WithoutStack(nil, map[string]any{"reason": "cert not found"})
			m.eh.ServeHTTPError(w, r, utilhttp.NewHTTPError(err, http.StatusBadRequest))
			return
		}

		cert, err := parseCert(ch)
		if err != nil {
			m.eh.ServeHTTPError(w, r, utilhttp.NewHTTPError(err, http.StatusBadRequest))
			return
		}

		// Verify the client certificate
		if _, err := cert.Verify(m.opts); err != nil {
			m.eh.ServeHTTPError(w, r, utilhttp.NewHTTPError(err, http.StatusUnauthorized))
			return
		}

		fh := r.Header.Get("X-SSL-Client-Fingerprint")

		// Verify the fingerprint
		if fh != "" {
			f := sha256.Sum256(cert.Raw)
			if hex.EncodeToString(f[:]) != fh {
				err := app.ErrAppMiddleInvalidFingerprint.WithoutStack(nil, nil)
				m.eh.ServeHTTPError(w, r, utilhttp.NewHTTPError(err, http.StatusUnauthorized))
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}

func parseCert(ch string) (*x509.Certificate, error) {
	// Decode a Base64-encoded client certificate and convert it to a byte array
	decoded, err := base64.URLEncoding.DecodeString(ch)
	if err != nil {
		return nil, app.ErrAppMiddleInvalidCert.WithoutStack(err, map[string]any{"reason": "fail base64 decode"})
	}

	// Convert the client certificate into a PEM block
	block, _ := pem.Decode(decoded)
	if block == nil {
		return nil, app.ErrAppMiddleInvalidCert.WithoutStack(nil, map[string]any{"reason": "fail PEM decode"})
	}

	// Analyze the PEM block
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, app.ErrAppMiddleInvalidCert.WithoutStack(err, map[string]any{"reason": "fail x509 parse"})
	}

	return cert, nil
}
