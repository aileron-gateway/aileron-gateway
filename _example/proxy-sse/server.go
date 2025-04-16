// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {

	addr := "0.0.0.0:9999"
	log.Println("SSE server listens at", addr)

	http.HandleFunc("/", sse)

	if err := http.ListenAndServe(addr, nil); err != nil && err != http.ErrServerClosed {
		panic(err)
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

		fmt.Fprintf(w, "Goodbye !!\n")
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
