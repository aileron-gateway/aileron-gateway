// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package httpserver

import (
	stdcmp "cmp"
	"context"
	"crypto/tls"
	"errors"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"regexp"
	"testing"
	"time"

	"github.com/aileron-gateway/aileron-gateway/core"
	"github.com/aileron-gateway/aileron-gateway/kernel/log"
	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/quic-go/quic-go/http3"
)

// testDir is the path to the test data.
// This path can be changed by the environmental variable.
var testDir = stdcmp.Or(os.Getenv("TEST_DIR"), "../../test/")

type testServer struct {
	serveTime    int
	shutdownTime int
	serveErr     error
	shutdownErr  error
}

func (s *testServer) Addr() string {
	return ""
}

func (s *testServer) Shutdown(ctx context.Context) error {
	ticker := time.NewTimer(time.Duration(s.shutdownTime) * time.Millisecond)
	defer ticker.Stop()

	select {
	case <-ctx.Done():
		return s.shutdownErr
	case <-ticker.C:
		return s.shutdownErr
	}
}

func (s *testServer) Serve() error {
	time.Sleep(time.Millisecond * time.Duration(s.serveTime))
	return s.serveErr
}

func TestRunner_Run(t *testing.T) {
	type condition struct {
		s         *runner
		doneAfter time.Duration
	}

	type action struct {
		err        any // error or errorutil.Kind
		errPattern *regexp.Regexp
	}
	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	cndServe := tb.Condition("serve", "run server")
	cndServeErr := tb.Condition("serve error", "return error on serve")
	cndShutdownErr := tb.Condition("shutdown error", "return error on shutdown")
	actCheckNoError := tb.Action("no error", "check that there is no error returned")
	actCheckError := tb.Action("error", "check the returned error")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"success to start server",
			[]string{cndServe},
			[]string{actCheckNoError},
			&condition{
				s: &runner{
					lg: log.GlobalLogger(log.DefaultLoggerName),
					svr: &testServer{
						serveTime:    1,
						shutdownTime: 1,
					},
				},
				doneAfter: 10 * time.Millisecond,
			},
			&action{
				err: nil,
			},
		),
		gen(
			"failed to start server",
			[]string{cndServe, cndServeErr},
			[]string{actCheckError},
			&condition{
				s: &runner{
					lg: log.GlobalLogger(log.DefaultLoggerName),
					svr: &testServer{
						serveTime:    1,
						shutdownTime: 1,
						serveErr:     errors.New("test error"),
					},
				},
				doneAfter: 10 * time.Millisecond,
			},
			&action{
				err:        core.ErrCoreServer,
				errPattern: regexp.MustCompile(core.ErrPrefix + `error was returned from server`),
			},
		),
		gen(
			"success to shutdown server",
			[]string{cndServe},
			[]string{actCheckNoError},
			&condition{
				s: &runner{
					lg: log.GlobalLogger(log.DefaultLoggerName),
					svr: &testServer{
						serveTime:    50,
						shutdownTime: 1,
					},
				},
				doneAfter: 10 * time.Millisecond,
			},
			&action{
				err: nil,
			},
		),
		gen(
			"failed to shutdown server",
			[]string{cndServe, cndShutdownErr},
			[]string{actCheckNoError},
			&condition{
				s: &runner{
					lg: log.GlobalLogger(log.DefaultLoggerName),
					svr: &testServer{
						serveTime:    10,
						shutdownTime: 1,
						shutdownErr:  errors.New("test error"),
					},
				},
				doneAfter: 1 * time.Millisecond,
			},
			&action{
				err: nil,
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			ctx := context.Background()
			if tt.C().doneAfter > 0 {
				newCtx, cancel := context.WithCancel(ctx)
				ctx = newCtx
				time.AfterFunc(tt.C().doneAfter, func() {
					cancel()
				})
			}

			err := tt.C().s.Run(ctx)
			testutil.DiffError(t, tt.A().err, tt.A().errPattern, err)
		})
	}
}

func TestHttp2Server(t *testing.T) {
	type condition struct {
		s             *http2Server
		shutdownAfter time.Duration
	}

	type action struct {
		addr       string
		err        any // error or errorutil.Kind
		errPattern *regexp.Regexp
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	cndHTTPServer := tb.Condition("HTTP", "start and shutdown HTTP server")
	cndHTTPSServer := tb.Condition("HTTPS", "start and shutdown HTTPS server")
	actCheckNoError := tb.Action("no error", "check that there is no error for start and shutdown except for ErrServerClosed")
	table := tb.Build()

	certFile := testDir + "ut/core/server/server.crt"
	keyFile := testDir + "ut/core/server/server.key"
	cert, _ := tls.LoadX509KeyPair(certFile, keyFile)

	// Get available address for testing.
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		panic(err)
	}
	ln.Close()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"success to start/shutdown HTTP server",
			[]string{cndHTTPServer},
			[]string{actCheckNoError},
			&condition{
				s: &http2Server{
					svr: &http.Server{
						Addr: ln.Addr().String(),
					},
				},
				shutdownAfter: 10 * time.Millisecond,
			},
			&action{
				addr: ln.Addr().String(),
				err:  http.ErrServerClosed,
			},
		),
		gen(
			"success to start/shutdown HTTPS server",
			[]string{cndHTTPSServer},
			[]string{actCheckNoError},
			&condition{
				s: &http2Server{
					svr: &http.Server{
						Addr: ln.Addr().String(),
						TLSConfig: &tls.Config{
							Certificates: []tls.Certificate{cert},
						},
					},
				},
				shutdownAfter: 10 * time.Millisecond,
			},
			&action{
				addr: ln.Addr().String(),
				err:  http.ErrServerClosed,
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			listener, err := net.Listen("tcp4", tt.C().s.svr.Addr)
			if err != nil {
				panic(err) // Should not panic for this test.
			}
			time.AfterFunc(tt.C().shutdownAfter, func() {
				err := tt.C().s.Shutdown(context.Background())
				testutil.Diff(t, nil, err, cmpopts.EquateErrors())
			})
			tt.C().s.listener = listener

			err = tt.C().s.Serve()
			testutil.Diff(t, tt.A().addr, tt.C().s.Addr())
			testutil.DiffError(t, tt.A().err, tt.A().errPattern, err, cmpopts.EquateErrors())
		})
	}
}

func TestHttp3Server(t *testing.T) {
	type condition struct {
		s             *http3Server
		shutdownAfter time.Duration
	}

	type action struct {
		addr        string
		err         error
		shutdownErr error
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	cndHTTP3Server := tb.Condition("http3 server", "start and shutdown http3 server")
	actCheckNoError := tb.Action("no error", "check that there is no error")
	actCheckError := tb.Action("error", "check that an error was returned")
	table := tb.Build()

	certFile := testDir + "ut/core/server/server.crt"
	keyFile := testDir + "ut/core/server/server.key"
	cert, _ := tls.LoadX509KeyPair(certFile, keyFile)

	// Get available port for testing.
	ln4, err := net.ListenPacket("udp4", "127.0.0.1:0")
	if err != nil {
		panic(err)
	}
	ln4.Close()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"success to start/shutdown http3 server",
			[]string{cndHTTP3Server},
			[]string{actCheckNoError},
			&condition{
				s: &http3Server{
					svr: &http3.Server{
						Addr: ln4.LocalAddr().String(),
						TLSConfig: &tls.Config{
							Certificates: []tls.Certificate{cert},
						},
					},
				},
				shutdownAfter: 10 * time.Millisecond,
			},
			&action{
				addr: ln4.LocalAddr().String(),
				err:  http.ErrServerClosed,
			},
		),
		gen(
			"serve error",
			[]string{cndHTTP3Server},
			[]string{actCheckError},
			&condition{
				s: &http3Server{
					svr: &http3.Server{
						Addr:      ln4.LocalAddr().String(),
						TLSConfig: nil,
					},
				},
				shutdownAfter: 0,
			},
			&action{
				addr: ln4.LocalAddr().String(),
				err:  errors.New("use of http3.Server without TLSConfig"),
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			time.AfterFunc(tt.C().shutdownAfter, func() {
				err := tt.C().s.Shutdown(context.Background())
				if tt.A().shutdownErr != nil {
					testutil.Diff(t, tt.A().shutdownErr.Error(), err.Error())
				}
			})

			err := tt.C().s.Serve()
			testutil.Diff(t, tt.A().addr, tt.C().s.Addr())
			if tt.A().err != nil {
				testutil.Diff(t, tt.A().err.Error(), err.Error())
			}
		})
	}
}

func TestAltSvcMiddleware(t *testing.T) {
	type condition struct {
		value string
	}

	type action struct {
		value string
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"valid address",
			[]string{},
			[]string{},
			&condition{
				value: `h3=":8443"; ma=2592000`,
			},
			&action{
				value: `h3=":8443"; ma=2592000`,
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			w := httptest.NewRecorder()
			r, _ := http.NewRequest(http.MethodGet, "http://test.com", nil)

			m := altSvcMiddleware(tt.C().value)
			h := m.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Write(nil)
			}))

			h.ServeHTTP(w, r)

			resp := w.Result()
			altSvc := resp.Header.Get("Alt-Svc")
			testutil.Diff(t, tt.A().value, altSvc)
		})
	}
}
