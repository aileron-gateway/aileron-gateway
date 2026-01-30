// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package network

import (
	"crypto/tls"
	"crypto/x509"
	"os"

	k "github.com/aileron-gateway/aileron-gateway/apis/kernel"
	"github.com/aileron-gateway/aileron-gateway/kernel/errorutil"
)

// tlsCipher is the list of TLS cipher suites.
var tlsCipher = map[k.TLSCipher]uint16{
	k.TLSCipher_TLS_RSA_WITH_RC4_128_SHA:                      tls.TLS_RSA_WITH_RC4_128_SHA,
	k.TLSCipher_TLS_RSA_WITH_3DES_EDE_CBC_SHA:                 tls.TLS_RSA_WITH_3DES_EDE_CBC_SHA,
	k.TLSCipher_TLS_RSA_WITH_AES_128_CBC_SHA:                  tls.TLS_RSA_WITH_AES_128_CBC_SHA,
	k.TLSCipher_TLS_RSA_WITH_AES_256_CBC_SHA:                  tls.TLS_RSA_WITH_AES_256_CBC_SHA,
	k.TLSCipher_TLS_RSA_WITH_AES_128_CBC_SHA256:               tls.TLS_RSA_WITH_AES_128_CBC_SHA256,
	k.TLSCipher_TLS_RSA_WITH_AES_128_GCM_SHA256:               tls.TLS_RSA_WITH_AES_128_GCM_SHA256,
	k.TLSCipher_TLS_RSA_WITH_AES_256_GCM_SHA384:               tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
	k.TLSCipher_TLS_ECDHE_ECDSA_WITH_RC4_128_SHA:              tls.TLS_ECDHE_ECDSA_WITH_RC4_128_SHA,
	k.TLSCipher_TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA:          tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA,
	k.TLSCipher_TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA:          tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA,
	k.TLSCipher_TLS_ECDHE_RSA_WITH_RC4_128_SHA:                tls.TLS_ECDHE_RSA_WITH_RC4_128_SHA,
	k.TLSCipher_TLS_ECDHE_RSA_WITH_3DES_EDE_CBC_SHA:           tls.TLS_ECDHE_RSA_WITH_3DES_EDE_CBC_SHA,
	k.TLSCipher_TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA:            tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
	k.TLSCipher_TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA:            tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
	k.TLSCipher_TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256:       tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256,
	k.TLSCipher_TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256:         tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256,
	k.TLSCipher_TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256:         tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
	k.TLSCipher_TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256:       tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
	k.TLSCipher_TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384:         tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
	k.TLSCipher_TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384:       tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
	k.TLSCipher_TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305_SHA256:   tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305_SHA256,
	k.TLSCipher_TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305_SHA256: tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305_SHA256,
	k.TLSCipher_TLS_AES_128_GCM_SHA256:                        tls.TLS_AES_128_GCM_SHA256,
	k.TLSCipher_TLS_AES_256_GCM_SHA384:                        tls.TLS_AES_256_GCM_SHA384,
	k.TLSCipher_TLS_CHACHA20_POLY1305_SHA256:                  tls.TLS_CHACHA20_POLY1305_SHA256,
	k.TLSCipher_TLS_FALLBACK_SCSV:                             tls.TLS_FALLBACK_SCSV,
	k.TLSCipher_TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305:          tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
	k.TLSCipher_TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305:        tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
}

// TLSConfig returns a new *tls.Config from the given spec.
// This function returns nil config and nil error when
// the given spec was nil.
func TLSConfig(spec *k.TLSConfig) (*tls.Config, error) {
	if spec == nil {
		return nil, nil
	}

	rootCAs, err := certPool(spec.RootCAs)
	if err != nil {
		return nil, errorutil.NewSimple(err, "internal/network: failed to load root CAs", "")
	}

	clientCAs, err := certPool(spec.ClientCAs)
	if err != nil {
		return nil, errorutil.NewSimple(err, "internal/network: failed to load root CAs", "")
	}

	if spec.ClientAuth > 4 {
		return nil, errorutil.NewSimple(nil, "internal/network: ClientAuthType must be 0 to 4.", "given=%s", spec.ClientAuth.String())
	}

	if spec.Renegotiation > 2 {
		return nil, errorutil.NewSimple(nil, "internal/network: RenegotiationSupport must be 0 to 2.", "given=%s", spec.Renegotiation.String())
	}
	certs := make([]tls.Certificate, 0, len(spec.CertKeyPairs))
	for _, pair := range spec.CertKeyPairs {
		cert, err := tls.LoadX509KeyPair(pair.CertFile, pair.KeyFile)
		if err != nil {
			return nil, errorutil.NewSimple(err, "internal/network: failed to load cert file.", "")
		}
		certs = append(certs, cert)
	}

	return &tls.Config{
		Rand:                        nil, // Use default crypto/rand
		Time:                        nil, // Use default time.Now
		Certificates:                certs,
		RootCAs:                     rootCAs,
		NextProtos:                  spec.NextProtos,
		ServerName:                  spec.ServerName,
		ClientAuth:                  tls.ClientAuthType(spec.ClientAuth),
		ClientCAs:                   clientCAs,
		InsecureSkipVerify:          spec.InsecureSkipVerify, //nolint:gosec // G402: TLS InsecureSkipVerify may be true.
		CipherSuites:                tlsCiphers(spec.TLSCiphers),
		SessionTicketsDisabled:      spec.SessionTicketsDisabled,
		MinVersion:                  uint16(spec.MinVersion), //nolint:gosec // G115: integer overflow conversion uint32 -> uint16
		MaxVersion:                  uint16(spec.MaxVersion), //nolint:gosec // G115: integer overflow conversion uint32 -> uint16
		CurvePreferences:            curveIDs(spec.CurvePreferences),
		DynamicRecordSizingDisabled: spec.DynamicRecordSizingDisabled,
		Renegotiation:               tls.RenegotiationSupport(spec.Renegotiation),
	}, nil
}

// tlsCiphers return a new slice of tls ciphers.
// Invalid ciphers will be ignored.
func tlsCiphers(cs []k.TLSCipher) []uint16 {
	if len(cs) == 0 {
		return nil
	}
	suites := make([]uint16, 0, len(cs))
	for _, c := range cs {
		if v, ok := tlsCipher[c]; ok {
			suites = append(suites, v)
		}
	}
	return suites
}

// curveIDs returns a new slice of tls.CurveID.
// Invalid curve ids are ignored.
func curveIDs(ids []k.CurveID) []tls.CurveID {
	if len(ids) == 0 {
		return nil
	}
	curves := make([]tls.CurveID, 0, len(ids))
	for _, id := range ids {
		switch id {
		case k.CurveID_CurveP256:
			curves = append(curves, tls.CurveP256)
		case k.CurveID_CurveP384:
			curves = append(curves, tls.CurveP384)
		case k.CurveID_CurveP521:
			curves = append(curves, tls.CurveP521)
		case k.CurveID_X25519:
			curves = append(curves, tls.X25519)
		}
	}
	return curves
}

// certPool returns new cert pool with given certifications
// reading from files.
func certPool(files []string) (*x509.CertPool, error) {
	pool, err := x509.SystemCertPool()
	if err != nil {
		pool = x509.NewCertPool()
	}
	for _, file := range files {
		pem, err := os.ReadFile(file)
		if err != nil {
			return nil, errorutil.NewSimple(err, "internal/network: TLS certification error", "")
		}
		if !pool.AppendCertsFromPEM(pem) {
			return nil, errorutil.NewSimple(nil, "internal/network: invalid pem file", "file=%s", file)
		}
	}
	return pool, nil
}
