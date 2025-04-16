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

const (
	certFilePath = "./_example/proxy-http2/pki/cert.pem"
	target       = "https://localhost:8443/test"
)

func main() {

	pem, err := os.ReadFile(certFilePath)
	if err != nil {
		log.Fatalln(err.Error())
	}
	pool := x509.NewCertPool()
	pool.AppendCertsFromPEM(pem)

	r, err := http.NewRequest(http.MethodGet, target, nil)
	if err != nil {
		log.Fatalln(err.Error())
	}

	// HTTP 2 transport.
	t := &http2.Transport{
		TLSClientConfig: &tls.Config{
			RootCAs: pool,
		},
	}

	log.Println("Send HTTP 1 request :", target)

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
