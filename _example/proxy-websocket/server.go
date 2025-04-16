// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"golang.org/x/net/websocket"
)

func main() {

	addr := "0.0.0.0:9999"
	log.Println("WebSocket server listens at", addr)

	dir := http.Dir("./")
	http.Handle("/", http.FileServer(dir))
	http.HandleFunc("/ws", ws())

	if err := http.ListenAndServe(addr, nil); err != nil && err != http.ErrServerClosed {
		panic(err)
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
