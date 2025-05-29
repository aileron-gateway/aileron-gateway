// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package main

import (
	"fmt"
	"net/http"
	"time"
)

// main just runs a simple HTTP server.
func main() {
	addr := ":9090"
	println("server listening at", addr)
	svr := &http.Server{
		Addr: addr,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, "Hello, AILERON Gateway!!")
		}),
		ReadHeaderTimeout: 10 * time.Second,
	}
	if err := svr.ListenAndServe(); err != nil {
		panic(err)
	}
}
