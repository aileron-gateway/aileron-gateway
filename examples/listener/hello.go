package main

import (
	"fmt"
	"net/http"
)

// main just runs a simple HTTP server.
func main() {
	addr := ":9090"
	println("server listening at", addr)
	err := http.ListenAndServe(addr, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, AILERON Gateway!!")
	}))
	if err != nil {
		panic(err)
	}
}
