// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"io"
	"log"
	"net/http"
	"os"

	"golang.org/x/net/http2"
)

func main() {
	const certFile = "./pki/cert.pem"
	const target = "https://localhost:8443/test"

	pem, err := os.ReadFile(certFile)
	if err != nil {
		log.Fatalln(err.Error())
	}
	pool := x509.NewCertPool()
	pool.AppendCertsFromPEM(pem)

	t := &http2.Transport{ // HTTP 2 transport.
		TLSClientConfig: &tls.Config{
			RootCAs:    pool,
			MinVersion: tls.VersionTLS12,
		},
	}

	log.Println("Send HTTP 2 request :", target)
	r, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, target, nil)
	w, err := t.RoundTrip(r)
	if err != nil {
		log.Fatalln(err.Error())
	}
	defer w.Body.Close()
	b, err := io.ReadAll(w.Body)
	if err != nil {
		log.Fatalln(err.Error())
	}
	log.Println(http.StatusText(w.StatusCode))
	log.Println(string(b))
}
