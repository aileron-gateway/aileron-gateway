// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package main_test

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/aileron-gateway/aileron-gateway/apis/kernel"
	"github.com/aileron-gateway/aileron-gateway/cmd/aileron/app"
	"github.com/aileron-gateway/aileron-gateway/core"
	"github.com/aileron-gateway/aileron-gateway/kernel/api"
	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
	"github.com/aileron-gateway/aileron-gateway/util/register"
)

const (
	certFile = "./pki/cert.pem"
	keyFile  = "./pki/key.pem"
)

func runHTTP1(t *testing.T) *http.Server {
	svr := &http.Server{
		Addr: ":10001",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, "OK")
		}),
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler)),
	}
	log.Println("HTTP 1 server listens on", svr.Addr)
	go func() {
		if err := svr.ListenAndServeTLS(certFile, keyFile); err != nil && err != http.ErrServerClosed {
			t.Error(err)
		}
	}()
	return svr
}

func TestProxyHTTP1(t *testing.T) {
	entrypoint := getEntrypointRunner(t, "config-http1.yaml")

	pem, err := os.ReadFile(certFile)
	if err != nil {
		t.Error(err)
	}
	pool := x509.NewCertPool()
	pool.AppendCertsFromPEM(pem)
	transport := &http.Transport{
		ForceAttemptHTTP2: false,
		TLSClientConfig: &tls.Config{
			RootCAs: pool,
		},
		ResponseHeaderTimeout: 3 * time.Second,
	}

	svr := runHTTP1(t)
	defer svr.Close()
	time.Sleep(1 * time.Second) // Wait the server to start up.

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
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

func getEntrypointRunner(t *testing.T, config ...string) core.Runner {
	t.Helper()
	svr := api.NewDefaultServeMux()
	f := api.NewFactoryAPI()
	register.RegisterAll(f)
	svr.Handle("core/", f)
	if err := app.LoadConfigFiles(svr, config); err != nil {
		t.Error(err)
	}
	req := &api.Request{
		Method: api.MethodGet,
		Key:    "core/v1/Entrypoint",
		Format: api.FormatProtoReference,
		Content: &kernel.Reference{
			APIVersion: "core/v1",
			Kind:       "Entrypoint",
			Namespace:  ".entrypoint",
			Name:       ".entrypoint",
		},
	}
	res, err := svr.Serve(context.Background(), req)
	if err != nil {
		t.Error(err)
	}
	entrypoint, ok := res.Content.(core.Runner)
	if !ok {
		t.Error(errors.New("failed to assert entrypoint to Runner interface"))
	}
	return entrypoint
}
