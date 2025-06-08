// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package network

import (
	"context"
	"crypto/tls"
	"errors"
	"net"
	"net/http"
	"regexp"
	"testing"
	"time"

	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func TestDnsResolver_resolve(t *testing.T) {
	type condition struct {
		r *dnsResolver
	}

	type action struct {
		length int
		ips    []string
		next   []string
		err    error
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"127.0.0.1",
			[]string{},
			[]string{},
			&condition{
				r: &dnsResolver{
					host: "127.0.0.1",
				},
			},
			&action{
				length: 1,
				ips:    []string{"127.0.0.1"},
				next:   []string{"127.0.0.1", "127.0.0.1", "127.0.0.1"}, // Check 3 times.
			},
		),
		gen(
			"invalid host",
			[]string{},
			[]string{},
			&condition{
				r: &dnsResolver{
					host: "host.not.exist",
				},
			},
			&action{
				length: 0,
				ips:    nil,
				next:   []string{"", "", ""}, // Check 3 times.
				err: &net.DNSError{
					Err:        "no such host",
					Name:       "host.not.exist",
					IsNotFound: true,
				},
			},
		),
		gen(
			"invalid host/ips no overwrite",
			[]string{},
			[]string{},
			&condition{
				r: &dnsResolver{
					host: "host.not.exist",
					ips:  []string{"127.0.0.1", "127.0.0.2"},
				},
			},
			&action{
				length: 2,
				ips:    []string{"127.0.0.1", "127.0.0.2"},
				next:   []string{"127.0.0.2", "127.0.0.1", "127.0.0.2", "127.0.0.1"}, // Check 4 times.
				err: &net.DNSError{
					Err:        "no such host",
					Name:       "host.not.exist",
					IsNotFound: true,
				},
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			VerboseLogs = true

			resolver := tt.C().r
			err := resolver.resolve(context.Background())
			testutil.Diff(t, tt.A().err, err, cmpopts.IgnoreFields(net.DNSError{}, "Server"))
			testutil.Diff(t, tt.A().length, resolver.length())
			testutil.Diff(t, tt.A().ips, resolver.ips)
			for _, ip := range tt.A().next {
				testutil.Diff(t, resolver.next(), ip)
			}
		})
	}
}

func TestDnsResolver_resolveEveryInterval(t *testing.T) {
	type condition struct {
		r          *dnsResolver
		interval   time.Duration
		checkAfter time.Duration
	}

	type action struct {
		ips []string
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"resolved",
			[]string{},
			[]string{},
			&condition{
				r: &dnsResolver{
					host: "127.0.0.1",
				},
				interval:   100 * time.Millisecond,
				checkAfter: 500 * time.Millisecond,
			},
			&action{
				ips: []string{"127.0.0.1"},
			},
		),
		gen(
			"not yet resolved",
			[]string{},
			[]string{},
			&condition{
				r: &dnsResolver{
					host: "127.0.0.1",
				},
				interval:   500 * time.Millisecond,
				checkAfter: 100 * time.Millisecond,
			},
			&action{
				ips: nil,
			},
		),
		gen(
			"invalid host",
			[]string{},
			[]string{},
			&condition{
				r: &dnsResolver{
					host: "host.not.exist",
				},
				interval:   100 * time.Millisecond,
				checkAfter: 500 * time.Millisecond,
			},
			&action{
				ips: nil,
			},
		),
		gen(
			"invalid host / ips no overwrite",
			[]string{},
			[]string{},
			&condition{
				r: &dnsResolver{
					host: "host.not.exist",
					ips:  []string{"127.0.0.1", "127.0.0.2"},
				},
				interval:   100 * time.Millisecond,
				checkAfter: 500 * time.Millisecond,
			},
			&action{
				ips: []string{"127.0.0.1", "127.0.0.2"},
			},
		),
		gen(
			"0 interval",
			[]string{},
			[]string{},
			&condition{
				r: &dnsResolver{
					host: "127.0.0.1",
				},
				interval:   0,
				checkAfter: 100 * time.Millisecond,
			},
			&action{
				ips: nil,
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			VerboseLogs = true

			resolver := tt.C().r

			wait := make(chan struct{})
			go func() {
				close(wait)
				resolver.resolveEveryInterval(context.Background(), tt.C().interval)
			}()
			<-wait
			time.Sleep(tt.C().checkAfter)
			resolver.stopResolveLoop()

			testutil.Diff(t, tt.A().ips, resolver.ips)
		})
	}
}

func TestHostConns_markDead(t *testing.T) {
	type condition struct {
		hc   *hostConns
		conn *http2.ClientConn
	}

	type action struct {
		active   int
		conns    map[string][]*http2.ClientConn
		connToIP map[*http2.ClientConn]string
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	conn1 := &http2.ClientConn{}
	conn2 := &http2.ClientConn{}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"dead 2 of 1",
			[]string{},
			[]string{},
			&condition{
				hc: &hostConns{
					conns: map[string][]*http2.ClientConn{
						"127.0.0.1": {conn1, conn2},
					},
					connToIP: map[*http2.ClientConn]string{
						conn1: "127.0.0.1",
						conn2: "127.0.0.1",
					},
				},
				conn: conn1,
			},
			&action{
				active: 1,
				conns: map[string][]*http2.ClientConn{
					"127.0.0.1": {conn2},
				},
				connToIP: map[*http2.ClientConn]string{
					conn2: "127.0.0.1",
				},
			},
		),
		gen(
			"dead 1 of 1",
			[]string{},
			[]string{},
			&condition{
				hc: &hostConns{
					conns:    map[string][]*http2.ClientConn{"127.0.0.1": {conn1}},
					connToIP: map[*http2.ClientConn]string{conn1: "127.0.0.1"},
				},
				conn: conn1,
			},
			&action{
				active:   0,
				conns:    map[string][]*http2.ClientConn{},
				connToIP: map[*http2.ClientConn]string{},
			},
		),
		gen(
			"unknown connection",
			[]string{},
			[]string{},
			&condition{
				hc: &hostConns{
					conns:    map[string][]*http2.ClientConn{"127.0.0.1": {conn1}},
					connToIP: map[*http2.ClientConn]string{conn1: "127.0.0.1"},
				},
				conn: conn2,
			},
			&action{
				active:   1,
				conns:    map[string][]*http2.ClientConn{"127.0.0.1": {conn1}},
				connToIP: map[*http2.ClientConn]string{conn1: "127.0.0.1"},
			},
		),
		gen(
			"conns not found", // This case won't happen.
			[]string{},
			[]string{},
			&condition{
				hc: &hostConns{
					conns:    map[string][]*http2.ClientConn{"127.0.0.2": {conn1}},
					connToIP: map[*http2.ClientConn]string{conn1: "127.0.0.1"},
				},
				conn: conn1,
			},
			&action{
				active:   0,
				conns:    map[string][]*http2.ClientConn{"127.0.0.2": {conn2}},
				connToIP: map[*http2.ClientConn]string{},
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			opts := []cmp.Option{
				cmp.Comparer(testutil.ComparePointer[*http2.ClientConn]),
			}

			active := tt.C().hc.markDead(tt.C().conn)
			testutil.Diff(t, tt.A().active, active)
			testutil.Diff(t, tt.A().conns, tt.A().conns, opts...)
			testutil.Diff(t, tt.A().connToIP, tt.A().connToIP, opts...)
		})
	}
}

func TestHostConns_getClientConn(t *testing.T) {
	type condition struct {
		hc *hostConns
	}

	type action struct {
		conn          *http2.ClientConn
		conns         map[string][]*http2.ClientConn
		connToIP      map[*http2.ClientConn]string
		newConnection bool
		errPattern    *regexp.Regexp
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	testSvr := &http.Server{
		Addr:    ":12321",
		Handler: h2c.NewHandler(&http.ServeMux{}, &http2.Server{}),
	}
	go testSvr.ListenAndServe()
	time.Sleep(time.Second)
	defer testSvr.Close()

	c1, _ := (&net.Dialer{}).Dial("tcp", "127.0.0.1:12321")
	conn1, _ := (&http2.Transport{}).NewClientConn(c1)
	defer conn1.Shutdown(context.Background())

	c2, _ := (&net.Dialer{}).Dial("tcp", "127.0.0.1:12321")
	conn2, _ := (&http2.Transport{}).NewClientConn(c2)
	conn2.SetDoNotReuse()
	conn2.ReserveNewRequest() // conn2 is no longer available
	defer conn2.Shutdown(context.Background())

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"Create new connection",
			[]string{},
			[]string{},
			&condition{
				hc: &hostConns{
					port:        "12321",
					dnsResolver: &dnsResolver{ips: []string{"127.0.0.1"}},
					conns:       map[string][]*http2.ClientConn{},
					connToIP:    map[*http2.ClientConn]string{},
					t: &http2.Transport{
						DialTLSContext: func(ctx context.Context, network, addr string, cfg *tls.Config) (net.Conn, error) {
							return (&net.Dialer{}).Dial(network, addr)
						},
					},
				},
			},
			&action{
				conns:         map[string][]*http2.ClientConn{},
				connToIP:      map[*http2.ClientConn]string{},
				newConnection: true,
			},
		),
		gen(
			"Connection available",
			[]string{},
			[]string{},
			&condition{
				hc: &hostConns{
					dnsResolver: &dnsResolver{ips: []string{"127.0.0.1"}},
					conns:       map[string][]*http2.ClientConn{"127.0.0.1": {conn1}},
					connToIP:    map[*http2.ClientConn]string{},
					t: &http2.Transport{
						DialTLSContext: func(ctx context.Context, network, addr string, cfg *tls.Config) (net.Conn, error) {
							return (&net.Dialer{}).DialContext(ctx, network, addr)
						},
					},
				},
			},
			&action{
				conn:     conn1,
				conns:    map[string][]*http2.ClientConn{"127.0.0.1": {conn1}},
				connToIP: map[*http2.ClientConn]string{conn1: "127.0.0.1"},
			},
		),
		gen(
			"Connection available/skip unavailable",
			[]string{},
			[]string{},
			&condition{
				hc: &hostConns{
					dnsResolver: &dnsResolver{ips: []string{"127.0.0.1"}},
					conns:       map[string][]*http2.ClientConn{"127.0.0.1": {conn2, conn1}},
					connToIP:    map[*http2.ClientConn]string{},
					t: &http2.Transport{
						DialTLSContext: func(ctx context.Context, network, addr string, cfg *tls.Config) (net.Conn, error) {
							return (&net.Dialer{}).DialContext(ctx, network, addr)
						},
					},
				},
			},
			&action{
				conn:     conn1,
				conns:    map[string][]*http2.ClientConn{"127.0.0.1": {conn2, conn1}},
				connToIP: map[*http2.ClientConn]string{conn1: "127.0.0.1"},
			},
		),
		gen(
			"no ips",
			[]string{},
			[]string{},
			&condition{
				hc: &hostConns{
					dnsResolver: &dnsResolver{},
					conns:       map[string][]*http2.ClientConn{},
					connToIP:    map[*http2.ClientConn]string{},
				},
			},
			&action{
				conns:      map[string][]*http2.ClientConn{},
				connToIP:   map[*http2.ClientConn]string{},
				errPattern: regexp.MustCompile(`http2: no cached connection was available`),
			},
		),
		gen(
			"unreachable IP",
			[]string{},
			[]string{},
			&condition{
				hc: &hostConns{
					dnsResolver: &dnsResolver{
						ips: []string{"127.0.0.2"},
					},
					conns:    map[string][]*http2.ClientConn{},
					connToIP: map[*http2.ClientConn]string{},
					t: &http2.Transport{
						DialTLSContext: func(ctx context.Context, network, addr string, cfg *tls.Config) (net.Conn, error) {
							return nil, errors.New("unreachable")
						},
					},
				},
			},
			&action{
				conns:      map[string][]*http2.ClientConn{},
				connToIP:   map[*http2.ClientConn]string{},
				errPattern: regexp.MustCompile(`http2: no cached connection was available`),
			},
		),
		gen(
			"unreachable IPs",
			[]string{},
			[]string{},
			&condition{
				hc: &hostConns{
					dnsResolver: &dnsResolver{
						ips: []string{"127.0.0.2", "127.0.0.3"},
					},
					conns:    map[string][]*http2.ClientConn{},
					connToIP: map[*http2.ClientConn]string{},
					t: &http2.Transport{
						DialTLSContext: func(ctx context.Context, network, addr string, cfg *tls.Config) (net.Conn, error) {
							return nil, errors.New("unreachable")
						},
					},
				},
			},
			&action{
				conns:      map[string][]*http2.ClientConn{},
				connToIP:   map[*http2.ClientConn]string{},
				errPattern: regexp.MustCompile(`http2: no cached connection was available`),
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			opts := []cmp.Option{
				cmp.Comparer(testutil.ComparePointer[*http2.ClientConn]),
			}

			hc := tt.C().hc
			conn, err := hc.getClientConn(context.Background())
			if tt.A().errPattern != nil {
				t.Log(err.Error())
				testutil.Diff(t, true, tt.A().errPattern.MatchString(err.Error()))
			} else {
				testutil.Diff(t, nil, err)
			}

			if tt.A().newConnection {
				tt.A().conns["127.0.0.1"] = []*http2.ClientConn{conn}
				tt.A().connToIP[conn] = "127.0.0.1"
			} else {
				testutil.Diff(t, conn, tt.A().conn, opts...)
			}
			testutil.Diff(t, tt.A().conns, hc.conns, opts...)
			testutil.Diff(t, tt.A().connToIP, hc.connToIP, opts...)
		})
	}
}

func TestHTTP2ConnPool_GetClientConn(t *testing.T) {
	type condition struct {
		cp   *http2ConnPool
		req  *http.Request
		addr string
	}

	type action struct {
		doNotReuse    bool
		newConnection bool
		conn          *http2.ClientConn
		conns         map[string]*hostConns
		connMap       map[*http2.ClientConn]*hostConns
		errPattern    *regexp.Regexp
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	testSvr := &http.Server{
		Addr:    ":12321",
		Handler: h2c.NewHandler(&http.ServeMux{}, &http2.Server{}),
	}
	go testSvr.ListenAndServe()
	time.Sleep(time.Second)
	defer testSvr.Close()

	c1, _ := (&net.Dialer{}).Dial("tcp", "127.0.0.1:12321")
	conn1, _ := (&http2.Transport{}).NewClientConn(c1)
	defer conn1.Shutdown(context.Background())
	c2, _ := (&net.Dialer{}).Dial("tcp", "127.0.0.1:12321")
	conn2, _ := (&http2.Transport{}).NewClientConn(c2)
	defer conn2.Shutdown(context.Background())

	hc1 := &hostConns{
		dnsResolver: &dnsResolver{ips: []string{"127.0.0.1"}},
		conns:       map[string][]*http2.ClientConn{"127.0.0.1": {conn1}},
		connToIP:    make(map[*http2.ClientConn]string),
	}
	hc2 := &hostConns{
		dnsResolver: &dnsResolver{ips: []string{"127.0.0.1"}},
		conns:       map[string][]*http2.ClientConn{"127.0.0.1": {conn2}},
		connToIP:    make(map[*http2.ClientConn]string),
	}
	emptyHc := &hostConns{
		dnsResolver: &dnsResolver{},
	}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"connection available",
			[]string{},
			[]string{},
			&condition{
				cp: &http2ConnPool{
					conns:   map[string]*hostConns{"localhost:12321": hc1},
					connMap: make(map[*http2.ClientConn]*hostConns),
				},
				req:  &http.Request{},
				addr: "localhost:12321",
			},
			&action{
				conn:    conn1,
				conns:   map[string]*hostConns{"localhost:12321": hc1},
				connMap: map[*http2.ClientConn]*hostConns{conn1: hc1},
			},
		),
		gen(
			"do not reuse",
			[]string{},
			[]string{},
			&condition{
				cp: &http2ConnPool{
					conns:   map[string]*hostConns{"localhost:12321": hc2},
					connMap: make(map[*http2.ClientConn]*hostConns),
				},
				req:  &http.Request{Close: true},
				addr: "localhost:12321",
			},
			&action{
				doNotReuse: true,
				conn:       conn2,
				conns:      map[string]*hostConns{"localhost:12321": hc2},
				connMap:    map[*http2.ClientConn]*hostConns{conn2: hc2},
			},
		),
		gen(
			"connection not available",
			[]string{},
			[]string{},
			&condition{
				cp: &http2ConnPool{
					conns:   map[string]*hostConns{"localhost:12321": emptyHc},
					connMap: make(map[*http2.ClientConn]*hostConns),
				},
				req:  &http.Request{},
				addr: "localhost:12321",
			},
			&action{
				conns:      map[string]*hostConns{"localhost:12321": emptyHc},
				connMap:    map[*http2.ClientConn]*hostConns{},
				errPattern: regexp.MustCompile(`http2: no cached connection was available`),
			},
		),
		gen(
			"new connection failed / resolve error",
			[]string{},
			[]string{},
			&condition{
				cp: &http2ConnPool{
					conns:   map[string]*hostConns{},
					connMap: make(map[*http2.ClientConn]*hostConns),
					t: &http2.Transport{
						TLSClientConfig: &tls.Config{},
						DialTLSContext: func(ctx context.Context, network, addr string, cfg *tls.Config) (net.Conn, error) {
							return nil, errors.New("unreachable")
						},
					},
				},
				req:  &http.Request{},
				addr: "host.not.Found:0000",
			},
			&action{
				conns:      map[string]*hostConns{},
				connMap:    map[*http2.ClientConn]*hostConns{},
				errPattern: regexp.MustCompile(`no such host`),
			},
		),
		gen(
			"new connection failed / conn create error",
			[]string{},
			[]string{},
			&condition{
				cp: &http2ConnPool{
					conns:   map[string]*hostConns{},
					connMap: make(map[*http2.ClientConn]*hostConns),
					t: &http2.Transport{
						TLSClientConfig: &tls.Config{},
						DialTLSContext: func(ctx context.Context, network, addr string, cfg *tls.Config) (net.Conn, error) {
							return nil, errors.New("unreachable")
						},
					},
				},
				req:  &http.Request{},
				addr: "localhost:0000",
			},
			&action{
				conns:      map[string]*hostConns{},
				connMap:    map[*http2.ClientConn]*hostConns{},
				errPattern: regexp.MustCompile(`http2: no cached connection was available`),
			},
		),
		gen(
			"new connection created",
			[]string{},
			[]string{},
			&condition{
				cp: &http2ConnPool{
					conns:   map[string]*hostConns{},
					connMap: make(map[*http2.ClientConn]*hostConns),
					t: &http2.Transport{
						TLSClientConfig: &tls.Config{},
						DialTLSContext: func(ctx context.Context, network, addr string, cfg *tls.Config) (net.Conn, error) {
							return (&net.Dialer{}).Dial(network, addr)
						},
					},
				},
				req:  &http.Request{},
				addr: "localhost:12321",
			},
			&action{
				newConnection: true,
				conns:         map[string]*hostConns{},
				connMap:       map[*http2.ClientConn]*hostConns{},
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			opts := []cmp.Option{
				cmp.Comparer(testutil.ComparePointer[*http2.ClientConn]),
				cmp.Comparer(testutil.ComparePointer[*hostConns]),
			}

			cp := tt.C().cp
			conn, err := cp.GetClientConn(tt.C().req, tt.C().addr)
			if tt.A().errPattern != nil {
				t.Log(err.Error())
				testutil.Diff(t, true, tt.A().errPattern.MatchString(err.Error()))
			} else {
				testutil.Diff(t, nil, err)
			}
			if tt.A().newConnection {
				tt.A().conns[tt.C().addr] = cp.conns[tt.C().addr]
				tt.A().connMap[conn] = cp.conns[tt.C().addr]
			} else {
				testutil.Diff(t, conn, tt.A().conn, opts...)
			}
			testutil.Diff(t, tt.A().conns, cp.conns, opts...)
			testutil.Diff(t, tt.A().connMap, cp.connMap, opts...)
			if tt.A().doNotReuse {
				testutil.Diff(t, false, conn.CanTakeNewRequest())
			}
		})
	}
}

func TestHTTP2ConnPool_MarkDead(t *testing.T) {
	type condition struct {
		cp   *http2ConnPool
		conn *http2.ClientConn
	}

	type action struct {
		conns      map[string]*hostConns
		connMap    map[*http2.ClientConn]*hostConns
		errPattern *regexp.Regexp
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	conn1 := &http2.ClientConn{}
	conn2 := &http2.ClientConn{}

	hc1 := &hostConns{
		addr:        "localhost:12321",
		dnsResolver: &dnsResolver{ips: []string{"127.0.0.1"}},
		conns:       map[string][]*http2.ClientConn{"127.0.0.1": {conn1}, "127.0.0.2": {conn2}},
		connToIP:    map[*http2.ClientConn]string{conn1: "127.0.0.1", conn2: "127.0.0.2"},
	}
	hc2 := &hostConns{
		addr:        "localhost:12321",
		dnsResolver: &dnsResolver{ips: []string{"127.0.0.1"}},
		conns:       map[string][]*http2.ClientConn{"127.0.0.1": {conn1}},
		connToIP:    map[*http2.ClientConn]string{conn1: "127.0.0.1"},
	}
	// emptyHc := &hostConns{
	// 	dnsResolver: &dnsResolver{},
	// }

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"connection not found",
			[]string{},
			[]string{},
			&condition{
				cp: &http2ConnPool{
					conns:   map[string]*hostConns{},
					connMap: make(map[*http2.ClientConn]*hostConns),
				},
				conn: conn1,
			},
			&action{
				conns:   map[string]*hostConns{},
				connMap: map[*http2.ClientConn]*hostConns{},
			},
		),
		gen(
			"connection removed/active conn remains",
			[]string{},
			[]string{},
			&condition{
				cp: &http2ConnPool{
					conns:   map[string]*hostConns{"localhost:12321": hc1},
					connMap: map[*http2.ClientConn]*hostConns{conn1: hc1},
				},
				conn: conn1,
			},
			&action{
				conns:   map[string]*hostConns{"localhost:12321": hc1},
				connMap: map[*http2.ClientConn]*hostConns{},
			},
		),
		gen(
			"connection removed/no active conns",
			[]string{},
			[]string{},
			&condition{
				cp: &http2ConnPool{
					conns:   map[string]*hostConns{"localhost:12321": hc2},
					connMap: map[*http2.ClientConn]*hostConns{conn1: hc2},
				},
				conn: conn1,
			},
			&action{
				conns:   map[string]*hostConns{},
				connMap: map[*http2.ClientConn]*hostConns{},
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			VerboseLogs = true

			opts := []cmp.Option{
				cmp.Comparer(testutil.ComparePointer[*http2.ClientConn]),
				cmp.Comparer(testutil.ComparePointer[*hostConns]),
			}

			cp := tt.C().cp
			cp.MarkDead(tt.C().conn)
			testutil.Diff(t, tt.A().conns, cp.conns, opts...)
			testutil.Diff(t, tt.A().connMap, cp.connMap, opts...)
		})
	}
}
