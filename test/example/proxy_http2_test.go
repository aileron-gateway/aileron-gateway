// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

//go:build example

package example_test

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
	"golang.org/x/net/http2"
)

const (
	proxyHTTP2CertFilePath = "./_example/proxy-http2/pki/cert.pem"
	proxyHTTP2KeyFilePath  = "./_example/proxy-http2/pki/key.pem"
)

func runHTTP2(t *testing.T, ctx context.Context) {
	svr := &http.Server{
		Addr: ":10002",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, "OK")
		}),
	}

	log.Println("HTTP 2 server listens on", svr.Addr)
	go func() {
		if err := svr.ListenAndServeTLS(proxyHTTP2CertFilePath, proxyHTTP2KeyFilePath); err != nil && err != http.ErrServerClosed {
			t.Error(err)
		}
	}()

	<-ctx.Done()
	if err := svr.Shutdown(context.Background()); err != nil {
		t.Error(err)
	}
}

func TestProxyHttp2(t *testing.T) {

	wd, _ := os.Getwd()
	defer changeDirectory(t, wd)
	changeDirectory(t, "./../../")

	env := []string{}
	config := []string{"_example/proxy-http2/config-http2.yaml"}
	entrypoint := getEntrypointRunner(t, env, config)

	pem, err := os.ReadFile(proxyHTTP2CertFilePath)
	if err != nil {
		t.Error(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go runHTTP2(t, ctx)
	time.Sleep(1 * time.Second) // Wait the server to start up.

	pool := x509.NewCertPool()
	pool.AppendCertsFromPEM(pem)

	transport := &http2.Transport{
		TLSClientConfig: &tls.Config{
			RootCAs: pool,
		},
		ReadIdleTimeout: 3 * time.Second,
	}

	var resp *http.Response
	go func() {
		req, _ := http.NewRequest(http.MethodGet, "https://localhost:8443/test", nil)
		resp, err = transport.RoundTrip(req)
		cancel() // Stop the server and the proxy.
	}()

	if err := entrypoint.Run(ctx); err != nil {
		t.Error(err)
	}

	testutil.Diff(t, nil, err)
	testutil.Diff(t, http.StatusOK, resp.StatusCode)
}
