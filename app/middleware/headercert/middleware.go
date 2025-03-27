package headercert

import (
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"errors"
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
			err := app.ErrAppMiddleHeaderPolicy.WithoutStack(nil, map[string]any{"reason": "client certificate is not found"}) // TODO:エラーも自分で定義する（errors.goに）
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
			err := app.ErrAppMiddleHeaderPolicy.WithoutStack(nil, map[string]any{"reason": "fail to verify certificate"}) // TODO:エラーも自分で定義する（errors.goに）
			fmt.Println("fail to verify certificate")
			m.eh.ServeHTTPError(w, r, utilhttp.NewHTTPError(err, http.StatusUnauthorized))
			return
		}

		fh := r.Header.Get("X-SSL-Client-Fingerprint")

		// Verify the fingerprint
		if fh != "" {
			f := sha256.Sum256(cert.Raw)
			if hex.EncodeToString(f[:]) != fh {
				err := app.ErrAppMiddleHeaderPolicy.WithoutStack(nil, map[string]any{"reason": "fail to verify fingerprint"}) // TODO:エラーも自分で定義する（errors.goに）
				fmt.Println("fail to verify fingerprint")
				m.eh.ServeHTTPError(w, r, utilhttp.NewHTTPError(err, http.StatusUnauthorized))
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}

func parseCert(ch string) (*x509.Certificate, error) { 
	// Decode a Base64-encoded client certificate and convert it to a byte array
	decoded, err := base64.URLEncoding.DecodeString(ch) // TODO:StdEncodingで問題ないか確認する
	if err != nil {
		fmt.Println("fail base64")
		return nil, app.ErrAppGenCreateRequest.WithoutStack(err, map[string]any{}) // TODO:第二引数はエラーのパラメータのマップを与える
	}

	// Convert the client certificate into a PEM block
	block, _ := pem.Decode(decoded)
	if block == nil {
		fmt.Println("fail PEM")
		return nil, errors.New("failed to decode certificate")
	}

	// Analyze the PEM block
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		fmt.Println("fail x509")
		return nil, errors.New("failed to parse certificate")
	}

	return cert, nil
}
