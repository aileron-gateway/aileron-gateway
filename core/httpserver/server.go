// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package httpserver

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/aileron-gateway/aileron-gateway/core"
	"github.com/aileron-gateway/aileron-gateway/kernel/log"
	"github.com/quic-go/quic-go/http3"
)

// server is an interface of a HTTP server.
type server interface {
	Addr() string
	Shutdown(context.Context) error
	Serve() error
}

// runner is a HTTP/HTTPS server runner.
// This implements core.Runner interfaces.
type runner struct {
	svr     server
	lg      log.Logger
	timeout time.Duration
}

// Run starts this server.
func (s *runner) Run(sigCtx context.Context) error {
	// ctx, cancel := context.WithCancel(sigCtx)
	// defer cancel()

	serverClosed := make(chan struct{})

	go func(ctx context.Context) {
		// Wait until getting a signal.
		<-ctx.Done()

		// Apply graceful shutdown timeout.
		shutdownCtx, cancel := context.WithTimeout(ctx, s.timeout)
		defer cancel()

		// Graceful shutdown.
		msg := fmt.Sprintf("server shutting down %s with graceful period %.0f seconds.", s.svr.Addr(), s.timeout.Seconds())
		s.lg.Info(ctx, msg)
		if err := s.svr.Shutdown(shutdownCtx); err != nil && err != http.ErrServerClosed {
			msg := fmt.Sprintf("server shut down failed. [%v]", err)
			s.lg.Info(ctx, msg) // May be shutdown timeout. We do not treat this as ERROR.
		}
		if close, ok := s.svr.(interface{ Close() error }); ok {
			close.Close() // Just in case.
		}
		serverClosed <- struct{}{}
	}(sigCtx)

	// Start the server.
	s.lg.Info(sigCtx, "server started. listening on "+s.svr.Addr())
	if err := s.svr.Serve(); err != nil && err != http.ErrServerClosed {
		err := core.ErrCoreServer.WithStack(err, nil)
		s.lg.Error(sigCtx, "error serving.", err.Name(), err.Map())
		return err
	}
	<-serverClosed // Wait the server fully closed.

	return nil
}

type http2Server struct {
	svr *http.Server
	// listener is the network listener used by the server.
	// listener cannot be reused once closed.
	// This listener will be closed after the http.Server existed.
	listener net.Listener
}

func (s *http2Server) Addr() string {
	return s.listener.Addr().String()
}

func (s *http2Server) Shutdown(ctx context.Context) error {
	// No need to care the type of returned error
	// even it is http.ErrServerClosed.
	// The returned error will be handled by the runner.
	err := s.svr.Shutdown(ctx)
	s.svr.Close()
	return err
}

func (s *http2Server) Serve() error {
	return s.svr.Serve(s.listener)
}

// http3Server is the http3 server.
// This implements server interface.
type http3Server struct {
	svr *http3.Server
}

func (s *http3Server) Addr() string {
	return s.svr.Addr
}

func (s *http3Server) Shutdown(ctx context.Context) error {
	// No need to care the type of returned error
	// even if it is http.ErrServerClosed.
	// The returned error will be handled by the runner.
	// TODO: Change to graceful shutdown when it is implemented in the quic package.
	err := s.svr.Close()
	s.svr.Close()
	return err
}

func (s *http3Server) Serve() error {
	return s.svr.ListenAndServe()
}

// altSrvMiddleware is a middleware that append
// alt-svc header in the responses.
// See the document below for Alt-Svc header and h3 values.
//   - https://datatracker.ietf.org/doc/rfc7838/
//   - https://datatracker.ietf.org/doc/draft-duke-httpbis-quic-version-alt-svc/
type altSvcMiddleware string

func (a altSvcMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h := w.Header()
		h["Alt-Svc"] = append(h["Alt-Svc"], string(a))
		next.ServeHTTP(w, r)
	})
}
