//go:build example
// +build example

package example_test

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
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
	// Set response headers.
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Transfer-Encoding", "identity")
	// Flush before sending the body.
	flusher.Flush()
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

func TestProxySSE(t *testing.T) {
	wd, _ := os.Getwd()
	defer changeDirectory(t, wd)
	changeDirectory(t, "./../../")

	env := []string{}
	config := []string{"./_example/proxy-sse/config.yaml"}
	entrypoint := getEntrypointRunner(t, env, config)

	ctx, cancel := context.WithCancel(context.Background())
	timer := time.AfterFunc(5*time.Second, cancel)
	go runSSEServer(t, ctx)
	time.Sleep(1 * time.Second) // Wait for the server start up.

	var resp *http.Response
	var err error
	go func() {
		req, _ := http.NewRequest(http.MethodGet, "http://localhost:8080", nil)
		resp, err = http.DefaultTransport.RoundTrip(req)
		timer.Stop()
		cancel()
	}()

	if err := entrypoint.Run(ctx); err != nil {
		t.Error(err)
	}

	testutil.Diff(t, nil, err)
	testutil.Diff(t, http.StatusOK, resp.StatusCode)
	body, _ := io.ReadAll(resp.Body)
	testutil.Diff(t, "12345", string(body))
}
