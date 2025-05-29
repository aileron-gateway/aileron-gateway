// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package main

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

func main() {
	addr := ":9090"
	println("server listening at", addr)
	svr := &http.Server{
		Addr: addr,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", r.Header.Get("Content-Type"))
			b, _ := io.ReadAll(r.Body)
			_, _ = w.Write(b)
			fmt.Println("--------------------------------------------------")
			fmt.Println(string(b))
		}),
		ReadHeaderTimeout: 10 * time.Second,
	}
	if err := svr.ListenAndServe(); err != nil {
		panic(err)
	}
}
