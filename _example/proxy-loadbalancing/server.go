// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
)

func main() {

	addrs := []string{
		":8001",
		":8002",
		":8003",
		":8004",
		":8005",
	}

	wg := sync.WaitGroup{}

	for _, addr := range addrs {
		wg.Add(1)
		a := addr
		go func() {
			defer wg.Done()
			mux := &http.ServeMux{}
			mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprintf(w, "Hello! %s from %s\n", r.RemoteAddr, a)
			})
			svr := &http.Server{
				Addr:    a,
				Handler: mux,
			}
			log.Println("Server listens at", a)
			if err := svr.ListenAndServe(); err != nil {
				panic(err)
			}
		}()
	}

	wg.Wait()

}
