//go:build example
//+build example

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
	
	"github.com/quic-go/quic-go/http3"
	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
)

func MoveToRootDirectory(t *testing.T, targetDir string) {
	err := os.Chdir(targetDir)
	if err != nil {
		t.Fatal(err)
		return
	}
}

func TestProxyHttp3(t *testing.T) {

	targetDir := "./../.."
	MoveToRootDirectory(t, targetDir)

	env := []string{}
	config := []string{"./_example/proxy-http3/config-http3.yaml"}
	entrypoint := getEntrypointRunner(t, env, config)
	certFilePath := "./_example/proxy-http3/pki/cert.pem"

	ctx, cancel := context.WithCancel(context.Background())
	timer := time.AfterFunc(5*time.Second, cancel)

	pem, err := os.ReadFile(certFilePath)
	if err != nil {
		t.Fatal(err)
	}

	cmd := exec.Command("go", "run", "./_example/proxy-http3/server.go")
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}

	defer func() {
		cmd.Process.Kill()
		cmd.Wait()
	}()

	time.Sleep(3*time.Second)

	pool := x509.NewCertPool()
	pool.AppendCertsFromPEM(pem)

	var resp *http.Response
	transport := &http3.Transport{
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