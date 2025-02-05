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

	"golang.org/x/net/http2"
	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
)

const (
	certFilePath = "./_example/proxy-http1/pki/cert.pem"
	keyFilePath  = "./_example/proxy-http1/pki/key.pem"
)

func runHTTP2(t *testing.T, ctx context.Context) {
	svr := &http.Server{
		Addr:         ":10002",
		Handler:      http.HandlerFunc(handler),
	}

	log.Println("HTTP 2 server listens on", svr.Addr)

	go func() {
		if err := svr.ListenAndServeTLS(certFilePath, keyFilePath); err != nil {
			t.Error(err)
		}
	}()

	<-ctx.Done()

	time.Sleep(1 * time.Second)

	if err := svr.Shutdown(context.Background()); err != nil {
		t.Error(err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Method : %s\n", r.Method)
	fmt.Fprintf(w, "Path : %s\n", r.URL.Path)
	fmt.Fprintf(w, "HTTP : %d.%d\n", r.ProtoMajor, r.ProtoMinor)
	fmt.Fprintf(w, "Header:\n")
	for k, v := range r.Header {
		fmt.Fprintf(w, "  %s: %+v\n", k, v)
	}
}

func TestProxyHttp2(t *testing.T) {

	targetDir := "./../.."
	changeDirectory(t, targetDir)

	env := []string{}
	config := []string{"./_example/proxy-http2/config-http2.yaml"}
	entrypoint := getEntrypointRunner(t, env, config)

	ctx, cancel := context.WithCancel(context.Background())
	timer := time.AfterFunc(5*time.Second, cancel)

	pem, err := os.ReadFile(certFilePath)
	if err != nil {
		t.Error(err)
	}

	go runHTTP2(t, ctx)

	time.Sleep(1 * time.Second)

	pool := x509.NewCertPool()
	pool.AppendCertsFromPEM(pem)

	transport := &http2.Transport{
		TLSClientConfig: &tls.Config{
			RootCAs: pool,
		},
	}
	
	var resp *http.Response
	
	go func() {
		req, _ := http.NewRequest(http.MethodGet, "https://localhost:8443/test", nil)
		resp, err = transport.RoundTrip(req)
		timer.Stop()
		cancel()
	}()

	if err := entrypoint.Run(ctx); err != nil {
		t.Error(err)
	}

	testutil.Diff(t, nil, err)
	testutil.Diff(t, http.StatusOK, resp.StatusCode)
}
