//go:build example
// +build example

package example_test

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"net/http"
	"os"
	"os/exec"
	"testing"
	"time"
	
	"golang.org/x/net/http2"
	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
)

func MoveToRootDirectory(t *testing.T, targetDir string) {
	err := os.Chdir(targetDir)
	if err != nil {
		t.Fatal(err)
		return
	}
}

func TestProxyHttp2(t *testing.T) {

	targetDir := "./../.."
	MoveToRootDirectory(t, targetDir)

	env := []string{}
	config := []string{"./_example/proxy-http2/config-http2.yaml"}
	entrypoint := getEntrypointRunner(t, env, config)
	certFilePath := "./_example/proxy-http2/pki/cert.pem"

	ctx, cancel := context.WithCancel(context.Background())
	timer := time.AfterFunc(5*time.Second, cancel)

	// Read the certificate
	pem, err := os.ReadFile(certFilePath)
	if err != nil {
		t.Fatal(err)
	}

	// Create and start a http server
	cmd := exec.Command("go", "run", "./_example/proxy-http2/server.go")
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}

	// Wait until the server finishes after all tests are done
	defer func() {
		cmd.Process.Kill()
		cmd.Wait()
	}()

	// Wait until the server starts completely
	time.Sleep(3 * time.Second)

	// Create a new certificate pool and add the certificate to the pool
	pool := x509.NewCertPool()
	pool.AppendCertsFromPEM(pem)

	var resp *http.Response
	transport := &http2.Transport{
		TLSClientConfig: &tls.Config{
			RootCAs: pool,
		},
	}

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
