package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/quic-go/quic-go/http3"
)

const (
	certFilePath = "./_example/proxy-http3/pki/cert.pem"
	keyFilePath  = "./_example/proxy-http3/pki/key.pem"
)

func main() {
	go runHTTP1()
	go runHTTP2()
	go runHTTP3()

	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	<-ctx.Done()
}

func runHTTP1() {
	svr := &http.Server{
		Addr:         ":10001",
		Handler:      http.HandlerFunc(handler),
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler)),
	}
	log.Println("HTTP 1 server listens on", svr.Addr)
	if err := svr.ListenAndServeTLS(certFilePath, keyFilePath); err != nil {
		panic(err)
	}
}

func runHTTP2() {
	svr := &http.Server{
		Addr:    ":10002",
		Handler: http.HandlerFunc(handler),
	}
	log.Println("HTTP 2 server listens on", svr.Addr)
	if err := svr.ListenAndServeTLS(certFilePath, keyFilePath); err != nil {
		panic(err)
	}
}

func runHTTP3() {
	svr := &http3.Server{
		Addr:    ":10003",
		Handler: http.HandlerFunc(handler),
	}
	log.Println("HTTP 3 server listens on", svr.Addr)
	if err := svr.ListenAndServeTLS(certFilePath, keyFilePath); err != nil {
		panic(err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Method : %s\n", r.Method)
	fmt.Fprintf(w, "Path : %s\n", r.URL.Path)
	fmt.Fprintf(w, "HTTP : %d.%d\n", r.ProtoMajor, r.ProtoMinor)
	fmt.Fprintf(w, "Header:\n")
	for k, v := range r.Header {
		fmt.Fprintf(w, "  %s: %+v\n", k, v)
	}
}
