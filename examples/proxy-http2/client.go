// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package main

import (
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
			RootCAs: pool,
		},
	}

	log.Println("Send HTTP 1 request :", target)
	r, err := http.NewRequest(http.MethodGet, target, nil)
	if err != nil {
		log.Fatalln(err.Error())
	}
	w, err := t.RoundTrip(r)
	if err != nil {
		log.Fatalln(err.Error())
	}
	b, err := io.ReadAll(w.Body)
	if err != nil {
		log.Fatalln(err.Error())
	}
	log.Println(http.StatusText(w.StatusCode))
	log.Println(string(b))
}
