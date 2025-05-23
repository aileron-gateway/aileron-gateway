// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package main_test

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/aileron-gateway/aileron-gateway/apis/kernel"
	"github.com/aileron-gateway/aileron-gateway/cmd/aileron/app"
	"github.com/aileron-gateway/aileron-gateway/core"
	"github.com/aileron-gateway/aileron-gateway/kernel/api"
	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
	"github.com/aileron-gateway/aileron-gateway/util/register"
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
	entrypoint := getEntrypointRunner(t, "./config.yaml")

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
		timer.Stop() // Stop the timer
		cancel()     // and immediately stop the server.
	}()

	if err := entrypoint.Run(ctx); err != nil {
		t.Error(err)
	}
	testutil.Diff(t, nil, err)
	testutil.Diff(t, "Hello!12345", msg)
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
