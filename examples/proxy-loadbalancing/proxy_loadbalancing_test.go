// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package main_test

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/aileron-gateway/aileron-gateway/apis/kernel"
	"github.com/aileron-gateway/aileron-gateway/cmd/aileron/app"
	"github.com/aileron-gateway/aileron-gateway/core"
	"github.com/aileron-gateway/aileron-gateway/kernel/api"
	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
	"github.com/aileron-gateway/aileron-gateway/util/register"
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

func TestConfigRoundRobin(t *testing.T) {
	entrypoint := getEntrypointRunner(t, "./config-round-robin.yaml")

	ctx, cancel := context.WithCancel(context.Background())
	go runServers(t, ctx)
	time.Sleep(1 * time.Second) // Wait the server to start up.

	var resp *http.Response
	var err error
	go func() {
		req, _ := http.NewRequest(http.MethodGet, "http://localhost:8080/", nil)
		resp, err = http.DefaultTransport.RoundTrip(req)
		cancel() // and immediately stop the server.
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
