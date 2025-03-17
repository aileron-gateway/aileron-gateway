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
	"github.com/quic-go/quic-go/http3"
)

const (
	proxyHTTP3CertFilePath = "./_example/proxy-http3/pki/cert.pem"
	proxyHTTP3KeyFilePath  = "./_example/proxy-http3/pki/key.pem"
)

func runHTTP3(t *testing.T, ctx context.Context) {
	svr := &http3.Server{
		Addr: ":10003",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, "OK")
		}),
	}

	log.Println("HTTP 3 server listens on", svr.Addr)
	go func() {
		if err := svr.ListenAndServeTLS(proxyHTTP3CertFilePath, proxyHTTP3KeyFilePath); err != nil && err != http.ErrServerClosed {
			t.Error(err)
		}
	}()

	<-ctx.Done()
	if err := svr.Shutdown(context.Background()); err != nil {
		t.Error(err)
	}
}

func TestProxyHTTP3(t *testing.T) {

	wd, _ := os.Getwd()
	defer changeDirectory(t, wd)
	changeDirectory(t, "./../../")

	env := []string{}
	config := []string{"./_example/proxy-http3/config-http3.yaml"}
	entrypoint := getEntrypointRunner(t, env, config)

	pem, err := os.ReadFile(proxyHTTP3CertFilePath)
	if err != nil {
		t.Error(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go runHTTP3(t, ctx)
	time.Sleep(1 * time.Second) // Wait the server to start up.

	pool := x509.NewCertPool()
	pool.AppendCertsFromPEM(pem)

	transport := &http3.Transport{
		TLSClientConfig: &tls.Config{
			RootCAs: pool,
		},
	}

	var resp *http.Response
	go func() {
		req, _ := http.NewRequest(http.MethodGet, "https://localhost:8443/test", nil)
		resp, err = transport.RoundTrip(req)
		cancel()
	}()

	if err := entrypoint.Run(ctx); err != nil {
		t.Error(err)
	}

	testutil.Diff(t, nil, err)
	testutil.Diff(t, http.StatusOK, resp.StatusCode)
}
