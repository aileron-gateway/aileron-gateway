//go:build examle
// +build examle

package example_test

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
)

func runServer(t *testing.T, testCtx context.Context) {
	addr := "0.0.0.0:9999"
	log.Println("SSE server listens at", addr)

	wg := sync.WaitGroup{}

	svr := &http.Server{
		Addr:    addr,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { sse(w, r) }),
	}

	wg.Add(1)

	go func() {
		defer wg.Done()
		if err := svr.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			t.Fatal(err)
		}
	}()

	wg.Wait()

	<-testCtx.Done()

	err := svr.Shutdown(context.Background())
	if err != nil {
		t.Fatal(err)
	}

}

func sse(w http.ResponseWriter, r *http.Request) {
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
		fmt.Fprintf(w, "Hello !!\n")

		for i := 0; i < 30; i++ {
			select {
			case <-r.Context().Done():
				return
			case <-time.After(time.Second):
			}

			n, err := fmt.Fprintf(w, "It's %s\n", time.Now().Format(http.TimeFormat))
			if n > 0 {
				flusher.Flush()
			}
			if err != nil {
				panic(err)
			}
		}

		fmt.Fprintf(w, "Goodbye!!\n")
		flusher.Flush()
		close(done)
	}()

	select {
	case <-r.Context().Done():
		log.Println("Client closed connection.")
	case <-done:
		log.Println("Done!!")
	}
}

func TestProxySEE(t *testing.T) {

	targetDir := "./../.."
	changeDirectory(t, targetDir)

	env := []string{}
	config := []string{"./_example/proxy-sse/config.yaml"}
	entrypoint := getEntrypointRunner(t, env, config)

	ctx, cancel := context.WithCancel(context.Background())
	timer := time.AfterFunc(5*time.Second, cancel)

	go runServer(t, ctx)

	time.Sleep(1 * time.Second)

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

}
