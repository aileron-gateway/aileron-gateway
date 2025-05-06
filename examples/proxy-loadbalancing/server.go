// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

// main runs multiple servers.
func main() {
	var wg sync.WaitGroup
	for _, addr := range []string{":8001", ":8002", ":8003", ":8004", ":8005"} {
		wg.Add(1)
		go func() {
			defer wg.Done()
			svr := &http.Server{
				Addr: addr,
				Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					fmt.Fprintf(w, "Server %s\n", addr)
				}),
				ReadHeaderTimeout: 10 * time.Second,
			}
			log.Println("Server listens at", addr)
			if err := svr.ListenAndServe(); err != nil {
				panic(err)
			}
		}()
	}
	wg.Wait()
}
