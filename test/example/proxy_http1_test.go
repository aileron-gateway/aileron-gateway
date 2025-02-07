//go:build example
// +build example

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
)

const (
	proxyHTTP1CertFilePath = "./_example/proxy-http1/pki/cert.pem"
	proxyHTTP1KeyFilePath  = "./_example/proxy-http1/pki/key.pem"
)

func runHTTP1(t *testing.T, ctx context.Context) {
	svr := &http.Server{
		Addr: ":10001",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, "OK")
		}),
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler)),
	}

	log.Println("HTTP 1 server listens on", svr.Addr)
	go func() {
		if err := svr.ListenAndServeTLS(proxyHTTP1CertFilePath, proxyHTTP1KeyFilePath); err != nil && err != http.ErrServerClosed {
			t.Error(err)
		}
	}()

	<-ctx.Done()
	if err := svr.Shutdown(context.Background()); err != nil {
		t.Error(err)
	}
}

func TestProxyHttp1(t *testing.T) {

	wd, _ := os.Getwd()
	defer changeDirectory(t, wd)
	changeDirectory(t, "./../../")

	env := []string{}
	config := []string{"_example/proxy-http1/config-http1.yaml"}
	entrypoint := getEntrypointRunner(t, env, config)

	pem, err := os.ReadFile(proxyHTTP1CertFilePath)
	if err != nil {
		t.Error(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go runHTTP1(t, ctx)
	time.Sleep(1 * time.Second) // Wait the server to start up.

	pool := x509.NewCertPool()
	pool.AppendCertsFromPEM(pem)

	transport := &http.Transport{
		ForceAttemptHTTP2: false,
		TLSClientConfig: &tls.Config{
			RootCAs: pool,
		},
		ResponseHeaderTimeout: 3 * time.Second,
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
