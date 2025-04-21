// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

//go:build example

package example_test

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
)

func runServers(t *testing.T, ctx context.Context) {
	addrs := []string{
		":8001",
		":8002",
		":8003",
		":8004",
		":8005",
	}

	var wg sync.WaitGroup
	for _, addr := range addrs {
		wg.Add(1)
		mux := &http.ServeMux{}
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "Hello! %s from %s\n", r.RemoteAddr, addr)
		})
		svr := &http.Server{
			Addr:    addr,
			Handler: mux,
		}
		go func() {
			log.Println("Server listens at", addr)
			if err := svr.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				t.Error(err)
			}
		}()
		go func() {
			defer wg.Done()
			<-ctx.Done()
			if err := svr.Shutdown(context.Background()); err != nil {
				t.Error(err)
			}
		}()
	}
	wg.Wait()
}

func TestProxyLoadBalancing(t *testing.T) {
	wd, _ := os.Getwd()
	defer changeDirectory(t, wd)
	changeDirectory(t, "./../../")

	env := []string{}
	config := []string{"./_example/proxy-loadbalancing/config-direct-hash.yaml"}
	entrypoint := getEntrypointRunner(t, env, config)

	ctx, cancel := context.WithCancel(context.Background())
	timer := time.AfterFunc(5*time.Second, cancel)
	go runServers(t, ctx)
	time.Sleep(1 * time.Second) // Wait the server to start up.

	var resp *http.Response
	var err error
	go func() {
		req, _ := http.NewRequest(http.MethodGet, "http://localhost:8080/", nil)
		resp, err = http.DefaultTransport.RoundTrip(req)
		timer.Stop()
		cancel()
	}()

	if err := entrypoint.Run(ctx); err != nil {
		t.Error(err)
	}

	testutil.Diff(t, nil, err)
	testutil.Diff(t, http.StatusOK, resp.StatusCode)
}
