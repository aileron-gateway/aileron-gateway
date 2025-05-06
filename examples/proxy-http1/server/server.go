// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/quic-go/quic-go/http3"
)

const (
	certFile = "./pki/cert.pem"
	keyFile  = "./pki/key.pem"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(3)
	go runHTTP1(&wg)
	go runHTTP2(&wg)
	go runHTTP3(&wg)
	wg.Wait()
}

func runHTTP1(wg *sync.WaitGroup) {
	defer wg.Done()
	svr := &http.Server{
		Addr: ":10001",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "Method : %s\n", r.Method)
			fmt.Fprintf(w, "Path : %s\n", r.URL.Path)
			fmt.Fprintf(w, "HTTP : %d.%d\n", r.ProtoMajor, r.ProtoMinor)
		}),
		TLSNextProto:      make(map[string]func(*http.Server, *tls.Conn, http.Handler)),
		ReadHeaderTimeout: 10 * time.Second,
	}
	log.Println("HTTP 1 server listens on", svr.Addr)
	if err := svr.ListenAndServeTLS(certFile, keyFile); err != nil {
		panic(err)
	}
}

func runHTTP2(wg *sync.WaitGroup) {
	defer wg.Done()
	svr := &http.Server{
		Addr: ":10002",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "Method : %s\n", r.Method)
			fmt.Fprintf(w, "Path : %s\n", r.URL.Path)
			fmt.Fprintf(w, "HTTP : %d.%d\n", r.ProtoMajor, r.ProtoMinor)
		}),
		ReadHeaderTimeout: 10 * time.Second,
	}
	log.Println("HTTP 2 server listens on", svr.Addr)
	if err := svr.ListenAndServeTLS(certFile, keyFile); err != nil {
		panic(err)
	}
}

func runHTTP3(wg *sync.WaitGroup) {
	defer wg.Done()
	svr := &http3.Server{
		Addr: ":10003",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "Method : %s\n", r.Method)
			fmt.Fprintf(w, "Path : %s\n", r.URL.Path)
			fmt.Fprintf(w, "HTTP : %d.%d\n", r.ProtoMajor, r.ProtoMinor)
		}),
	}
	log.Println("HTTP 3 server listens on", svr.Addr)
	if err := svr.ListenAndServeTLS(certFile, keyFile); err != nil {
		panic(err)
	}
}
