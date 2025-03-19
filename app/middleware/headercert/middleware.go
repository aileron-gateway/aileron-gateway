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

	// ローカルからルート証明書を読み込み
	for _, c := range h.rootCAs {
		rootCertPEM, err := os.ReadFile(c)
		if err != nil {
			return nil, err
		}
		// ルート証明書をCertPoolに追加
		if !roots.AppendCertsFromPEM(rootCertPEM) {
			return nil, fmt.Errorf("failed to add root certificate to CertPool")
		}
	}

	return roots, nil
}

func (h *HeaderCertAuth) convertCert(certHeader string) (*x509.Certificate, error) {
	// base64でエンコードされている文字列のクライアント証明書をデコードし、バイト列に変換
	decodedCertHeader, err := base64.StdEncoding.DecodeString(certHeader)
	if err != nil {
		return nil, fmt.Errorf("failed to decode certificate")
	}

	// バイト列に変換したクライアント証明書をPEMブロックに変換
	certBlock, _ := pem.Decode(decodedCertHeader)
	if certBlock == nil {
		return nil, fmt.Errorf("failed to decode certificate")
	}

	// PEMブロックを解析
	cert, err := x509.ParseCertificate(certBlock.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse certificate")
	}

	return cert, nil
}

func (h *HeaderCertAuth) isFingerprintMatched(cert *x509.Certificate, fingerprintHeader string) bool {
	fingerprint := sha256.Sum256(cert.Raw)
	fmt.Println(":::fingerprint:::")
	fmt.Println(hex.EncodeToString(fingerprint[:]))
	return hex.EncodeToString(fingerprint[:]) == fingerprintHeader
}

func isCertExpired(cert *x509.Certificate) bool {
	fmt.Println(":::expiration date:::")
	fmt.Println(cert.NotAfter)
	return time.Now().After(cert.NotAfter)
}

func (h *HeaderCertAuth) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.Header)

		// 証明書の抽出
		certHeader := r.Header.Get("X-SSL-Client-Cert")
		fmt.Printf("certificate: %s\n", certHeader)
		if certHeader == "" {
			http.Error(w, "Client certificate not found", http.StatusBadRequest)
			return
		}

		// フィンガープリントの抽出
		fingerprintHeader := r.Header.Get("X-SSL-Client-Fingerprint")
		fmt.Printf("fingerprint: %s\n", fingerprintHeader)
		if certHeader == "" {
			http.Error(w, "Fingerprint is not found", http.StatusBadRequest)
			return
		}

		cert, err := h.convertCert(certHeader)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		fmt.Println(":::certificate:::")
		fmt.Println(cert)

		roots, err := h.loadRootCert()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		opts := x509.VerifyOptions{
			Roots: roots,
		}

		// クライアント証明書の検証
		if _, err := cert.Verify(opts); err != nil {
			http.Error(w, "Fail to verify certificate", http.StatusUnauthorized)
			return
		}

		// フィンガープリントの検証
		if !h.isFingerprintMatched(cert, fingerprintHeader) {
			http.Error(w, "Fail to verify fingerprint", http.StatusUnauthorized)
			return
		}

		// 有効期限の確認
		if isCertExpired(cert) {
			http.Error(w, "Certificate has expired", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
