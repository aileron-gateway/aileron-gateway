// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package main_test

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"testing"

	"github.com/aileron-gateway/aileron-gateway/apis/kernel"
	"github.com/aileron-gateway/aileron-gateway/cmd/aileron/app"
	"github.com/aileron-gateway/aileron-gateway/core"
	"github.com/aileron-gateway/aileron-gateway/kernel/api"
	"github.com/aileron-gateway/aileron-gateway/util/register"
)

func runSSEServer(t *testing.T, ctx context.Context) {
	addr := "0.0.0.0:9999"
	log.Println("SSE server listens at", addr)

	svr := &http.Server{
		Addr:    addr,
		Handler: http.HandlerFunc(sseHandlerFunc),
	}

	go func() {
		if err := svr.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			t.Error(err)
		}
	}()

	<-ctx.Done()
	err := svr.Shutdown(context.Background())
	if err != nil {
		t.Error(err)
	}
}

func sseHandlerFunc(w http.ResponseWriter, r *http.Request) {
	flusher, _ := w.(http.Flusher)
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Transfer-Encoding", "identity")
	flusher.Flush() // Flush before sending the body.
	done := make(chan struct{})
	go func() {
		for i := 1; i <= 5; i++ {
			fmt.Fprintf(w, strconv.Itoa(i))
			flusher.Flush()
		}
		close(done)
	}()
	select {
	case <-r.Context().Done():
		log.Println("Client closed connection.")
	case <-done:
		log.Println("Done!!")
	}
}

// func TestProxySSE(t *testing.T) {
// 	entrypoint := getEntrypointRunner(t, "./config.yaml")

// 	ctx, cancel := context.WithCancel(context.Background())
// 	go runSSEServer(t, ctx)
// 	time.Sleep(1 * time.Second) // Wait for the server start up.

// 	var resp *http.Response
// 	var err error
// 	go func() {
// 		req, _ := http.NewRequest(http.MethodGet, "http://localhost:8080/test", nil)
// 		resp, err = http.DefaultTransport.RoundTrip(req)
// 		cancel()     // and immediately stop the server.
// 	}()

// 	if err := entrypoint.Run(ctx); err != nil {
// 		t.Error(err)
// 	}
// 	testutil.Diff(t, nil, err)
// 	testutil.Diff(t, http.StatusOK, resp.StatusCode)
// 	body, _ := io.ReadAll(resp.Body)
// 	testutil.Diff(t, "12345", string(body))
// }

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
