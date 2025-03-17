//go:build example

package example_test

import (
	"context"
	"log"
	"net/http"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
	"golang.org/x/net/websocket"
)

func runWebSocketServer(t *testing.T, ctx context.Context) {
	addr := "0.0.0.0:9999"
	log.Println("WebSocket server listens at", addr)

	mux := http.NewServeMux()
	mux.Handle("/", ws())
	svr := &http.Server{
		Addr:    addr,
		Handler: mux,
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

func ws() http.HandlerFunc {
	return websocket.Handler(func(ws *websocket.Conn) {
		defer ws.Close()

		err := websocket.Message.Send(ws, "Hello!")
		if err != nil {
			panic(err)
		}

		done := make(chan struct{})
		// Send message to client.
		for i := 1; i <= 5; i++ {
			err := websocket.Message.Send(ws, strconv.Itoa(i))
			if err != nil {
				log.Println(err)
				close(done)
				break
			}
		}
		select {
		case <-ws.Request().Context().Done():
			return
		case <-done:
			return
		}
	}).ServeHTTP
}

func TestProxyWebsocket(t *testing.T) {
	wd, _ := os.Getwd()
	defer changeDirectory(t, wd)
	changeDirectory(t, "./../../")

	env := []string{}
	config := []string{"./_example/proxy-sse/config.yaml"}
	entrypoint := getEntrypointRunner(t, env, config)

	ctx, cancel := context.WithCancel(context.Background())
	timer := time.AfterFunc(5*time.Second, cancel)
	go runWebSocketServer(t, ctx)
	time.Sleep(1 * time.Second) // Wait for the server start up.

	var err error
	var msg string
	go func() {
		time.Sleep(time.Second)
		ws, dialErr := websocket.Dial("ws://localhost:8080", "", "http://localhost:8080")
		if dialErr != nil {
			panic(dialErr)
			// t.Error(dialErr)
		}
		defer ws.Close()
		for {
			tmp := ""
			if err := websocket.Message.Receive(ws, &tmp); err != nil {
				t.Error(err)
				break
			}
			msg += tmp
		}
		timer.Stop()
		cancel()
	}()

	if err := entrypoint.Run(ctx); err != nil {
		t.Error(err)
	}

	testutil.Diff(t, nil, err)
	testutil.Diff(t, "Hello!12345", msg)
}
