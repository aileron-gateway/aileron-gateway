//go:build example

// + build example

package example_test

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"testing"
	"time"

	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
	"golang.org/x/net/websocket"
)

func runServer(t *testing.T, ctx context.Context) {
	addr := "0.0.0.0:9999"
	log.Println("WebSocket server listens at", addr)

	dir := http.Dir("./")
	http.Handle("/", http.FileServer(dir))
	http.HandleFunc("/ws", ws())

	svr := &http.Server{
		Addr: addr,
	}

	go func() {
		if err := svr.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			t.Error(err)
		}
	}()

	time.Sleep(time.Second * 1)

	<-ctx.Done()

	err := svr.Shutdown(context.Background())
	if err != nil {
		t.Error(err)
	}
}

func ws() http.HandlerFunc {
	return websocket.Handler(func(ws *websocket.Conn) {
		defer ws.Close()

		err := websocket.Message.Send(ws, "Hello!! This is a WebSocket server!!")
		if err != nil {
			panic(err)
		}

		done := make(chan struct{})

		// Receive message from client.
		go func() {
			for {
				msg := ""
				err = websocket.Message.Receive(ws, &msg)
				if err != nil {
					log.Println(err)
					close(done)
					return
				}
				err := websocket.Message.Send(ws, fmt.Sprintf("Your message arrived: %s", msg))
				if err != nil {
					log.Println(err)
					close(done)
					return
				}
			}
		}()

		// Send message to client.
		for {
			select {
			case <-ws.Request().Context().Done():
				return
			case <-done:
				return
			case <-time.After(time.Second):
			}
			err := websocket.Message.Send(ws, fmt.Sprintf("It's %s\n", time.Now().Format(http.TimeFormat)))
			if err != nil {
				log.Println(err)
				close(done)
				break
			}
		}

	}).ServeHTTP
}

func TestProxyWebsocket(t *testing.T) {

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
		req, _ := http.NewRequest(http.MethodGet, "http://localhost:8080/_example/proxy-websocket/", nil)
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
