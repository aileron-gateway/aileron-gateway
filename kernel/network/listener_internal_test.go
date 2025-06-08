// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package network

import (
	"crypto/tls"
	"errors"
	"io"
	"net"
	"net/netip"
	"os"
	"testing"
	"time"

	k "github.com/aileron-gateway/aileron-gateway/apis/kernel"
	"github.com/aileron-gateway/aileron-gateway/kernel/er"
	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
	"github.com/aileron-projects/go/zsyscall"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

// testDir is the path to the test data.
var testDir = "../../test/"

type removeSocketTest struct {
	net.Listener   // For testing RemoveSocketListener
	net.PacketConn // For testing RemoveSocketConn
	addr           net.Addr
	closeErr       error
	closed         bool
}

func (t *removeSocketTest) LocalAddr() net.Addr {
	return t.addr
}

func (t *removeSocketTest) Addr() net.Addr {
	return t.addr
}

func (t *removeSocketTest) Close() error {
	t.closed = true
	return t.closeErr
}

func TestRemoveSocketConn(t *testing.T) {
	type condition struct {
		conn net.PacketConn
		// socket file path.
		// Created before test.
		socket    string
		testClose bool
	}

	type action struct {
		expect   net.PacketConn
		closeErr error
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"nil conn",
			[]string{},
			[]string{},
			&condition{
				conn:      nil,
				testClose: false,
			},
			&action{
				expect:   nil,
				closeErr: nil,
			},
		),
		gen(
			"nil addr",
			[]string{},
			[]string{},
			&condition{
				conn: &removeSocketTest{
					addr:     nil,
					closeErr: io.ErrUnexpectedEOF, // Return dummy error.
				},
				testClose: true,
			},
			&action{
				expect: &removeSocketTest{
					closeErr: io.ErrUnexpectedEOF,
				},
				closeErr: io.ErrUnexpectedEOF,
			},
		),
		gen(
			"non unix",
			[]string{},
			[]string{},
			&condition{
				conn:      (*net.UDPConn)(nil),
				testClose: false,
			},
			&action{
				expect:   (*net.UDPConn)(nil),
				closeErr: nil,
			},
		),
		gen(
			"abstract socket",
			[]string{},
			[]string{},
			&condition{
				conn: &removeSocketTest{
					addr:     &net.UnixAddr{Name: "@test", Net: "unix"},
					closeErr: io.ErrUnexpectedEOF, // Return dummy error.
				},
				testClose: true,
			},
			&action{
				expect: &removeSocketTest{
					addr:     &net.UnixAddr{Name: "@test", Net: "unix"},
					closeErr: io.ErrUnexpectedEOF,
				},
				closeErr: io.ErrUnexpectedEOF,
			},
		),
		gen(
			"path name socket",
			[]string{},
			[]string{},
			&condition{
				conn: &removeSocketTest{
					addr:     &net.UnixAddr{Name: testDir + "ut/kernel/network/remove-socket-test.sock", Net: "unix"},
					closeErr: io.ErrUnexpectedEOF,
				},
				testClose: true,
				socket:    testDir + "ut/kernel/network/remove-socket-test.sock",
			},
			&action{
				expect: &removeSocketConn{
					PacketConn: &removeSocketTest{
						addr:     &net.UnixAddr{Name: testDir + "ut/kernel/network/remove-socket-test.sock", Net: "unix"},
						closeErr: io.ErrUnexpectedEOF,
					},
					socket: testDir + "ut/kernel/network/remove-socket-test.sock",
				},
				closeErr: io.ErrUnexpectedEOF,
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			opts := []cmp.Option{
				cmp.AllowUnexported(removeSocketConn{}),
				cmp.AllowUnexported(removeSocketTest{}),
				cmpopts.IgnoreUnexported(net.UDPConn{}),
				cmpopts.EquateErrors(),
			}

			conn := RemoveSocketConn(tt.C().conn)
			testutil.Diff(t, tt.A().expect, conn, opts...)
			if !tt.C().testClose {
				return
			}

			if tt.C().socket != "" {
				f, err := os.Create(tt.C().socket)
				f.Close()
				testutil.Diff(t, nil, err)
			}

			err := conn.Close()
			testutil.Diff(t, tt.A().closeErr, err, cmpopts.EquateErrors())
			testutil.Diff(t, true, tt.C().conn.(*removeSocketTest).closed)

			if tt.C().socket != "" {
				info, err := os.Stat(tt.C().socket)
				testutil.Diff(t, true, os.IsNotExist(err))
				testutil.Diff(t, nil, info)
			}
		})
	}
}

func TestRemoveSocketListener(t *testing.T) {
	type condition struct {
		ln net.Listener
		// socket file path.
		// Created before test.
		socket    string
		testClose bool
	}

	type action struct {
		expect   net.Listener
		closeErr error
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"nil listener",
			[]string{},
			[]string{},
			&condition{
				ln:        nil,
				testClose: false,
			},
			&action{
				expect:   nil,
				closeErr: nil,
			},
		),
		gen(
			"nil addr",
			[]string{},
			[]string{},
			&condition{
				ln: &removeSocketTest{
					addr:     nil,
					closeErr: io.ErrUnexpectedEOF, // Return dummy error.
				},
				testClose: true,
			},
			&action{
				expect: &removeSocketTest{
					closeErr: io.ErrUnexpectedEOF,
				},
				closeErr: io.ErrUnexpectedEOF,
			},
		),
		gen(
			"non unix",
			[]string{},
			[]string{},
			&condition{
				ln:        (*net.TCPListener)(nil),
				testClose: false,
			},
			&action{
				expect:   (*net.TCPListener)(nil),
				closeErr: nil,
			},
		),
		gen(
			"abstract socket",
			[]string{},
			[]string{},
			&condition{
				ln: &removeSocketTest{
					addr:     &net.UnixAddr{Name: "@test", Net: "unix"},
					closeErr: io.ErrUnexpectedEOF, // Return dummy error.
				},
				testClose: true,
			},
			&action{
				expect: &removeSocketTest{
					addr:     &net.UnixAddr{Name: "@test", Net: "unix"},
					closeErr: io.ErrUnexpectedEOF,
				},
				closeErr: io.ErrUnexpectedEOF,
			},
		),
		gen(
			"path name socket",
			[]string{},
			[]string{},
			&condition{
				ln: &removeSocketTest{
					addr:     &net.UnixAddr{Name: testDir + "ut/kernel/network/remove-socket-test.sock", Net: "unix"},
					closeErr: io.ErrUnexpectedEOF,
				},
				testClose: true,
				socket:    testDir + "ut/kernel/network/remove-socket-test.sock",
			},
			&action{
				expect: &removeSocketListener{
					Listener: &removeSocketTest{
						addr:     &net.UnixAddr{Name: testDir + "ut/kernel/network/remove-socket-test.sock", Net: "unix"},
						closeErr: io.ErrUnexpectedEOF,
					},
					socket: testDir + "ut/kernel/network/remove-socket-test.sock",
				},
				closeErr: io.ErrUnexpectedEOF,
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			opts := []cmp.Option{
				cmp.AllowUnexported(removeSocketListener{}),
				cmp.AllowUnexported(removeSocketTest{}),
				cmpopts.IgnoreUnexported(net.TCPListener{}),
				cmpopts.EquateErrors(),
			}

			ln := RemoveSocketListener(tt.C().ln)
			testutil.Diff(t, tt.A().expect, ln, opts...)
			if !tt.C().testClose {
				return
			}

			if tt.C().socket != "" {
				f, err := os.Create(tt.C().socket)
				f.Close()
				testutil.Diff(t, nil, err)
			}

			err := ln.Close()
			testutil.Diff(t, tt.A().closeErr, err, cmpopts.EquateErrors())
			testutil.Diff(t, true, tt.C().ln.(*removeSocketTest).closed)

			if tt.C().socket != "" {
				info, err := os.Stat(tt.C().socket)
				testutil.Diff(t, true, os.IsNotExist(err))
				testutil.Diff(t, nil, info)
			}
		})
	}
}

type testConn struct {
	net.Conn
	addr     string
	closeErr error
	closed   bool

	readDead  time.Time
	writeDead time.Time

	readBuffer  int
	writeBuffer int
}

func (c *testConn) RemoteAddr() net.Addr {
	addrPort, err := netip.ParseAddrPort(c.addr)
	if err != nil {
		var addr *net.TCPAddr
		return addr
	}
	addr := net.TCPAddrFromAddrPort(addrPort)
	return addr
}

func (c *testConn) Close() error {
	c.closed = true
	return c.closeErr
}

func (c *testConn) SetReadDeadline(t time.Time) error {
	c.readDead = t
	return nil
}

func (c *testConn) SetWriteDeadline(t time.Time) error {
	c.writeDead = t
	return nil
}

func (c *testConn) SetReadBuffer(bytes int) error {
	c.readBuffer = bytes
	return nil
}

func (c *testConn) SetWriteBuffer(bytes int) error {
	c.writeBuffer = bytes
	return nil
}

func TestNewPacketConn(t *testing.T) {
	type condition struct {
		c *PacketConnConfig
	}

	type action struct {
		address string
		err     error
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	udpChecker, _ := net.ListenPacket("udp", "127.0.0.1:0")
	udpChecker.Close()
	testUDPAddr := udpChecker.LocalAddr().String()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"nil config",
			[]string{},
			[]string{},
			&condition{
				c: nil,
			},
			&action{
				address: "",
				err: &er.Error{
					Package:     ErrPkg,
					Type:        ErrTypePackConn,
					Description: ErrDscPackConn,
				},
			},
		),
		gen(
			"invalid ip",
			[]string{},
			[]string{},
			&condition{
				c: &PacketConnConfig{
					Network: "ip:1",
					Address: "127.0.0.0.1",
					SockOption: &zsyscall.SockOption{
						SO:   &zsyscall.SockSOOption{ReceiveBuffer: 5000},
						IP:   &zsyscall.SockIPOption{TTL: 20},
						IPV6: &zsyscall.SockIPV6Option{},
					},
				},
			},
			&action{
				address: "",
				err: &er.Error{
					Package:     ErrPkg,
					Type:        ErrTypePackConn,
					Description: ErrDscPackConn,
				},
			},
		),
		gen(
			"valid udp",
			[]string{},
			[]string{},
			&condition{
				c: &PacketConnConfig{
					Network: "udp",
					Address: testUDPAddr,
					SockOption: &zsyscall.SockOption{
						SO:   &zsyscall.SockSOOption{ReceiveBuffer: 5000},
						IP:   &zsyscall.SockIPOption{TTL: 20},
						IPV6: &zsyscall.SockIPV6Option{},
						UDP:  &zsyscall.SockUDPOption{GRO: true},
					},
				},
			},
			&action{
				address: testUDPAddr,
				err:     nil,
			},
		),
		gen(
			"invalid udp",
			[]string{},
			[]string{},
			&condition{
				c: &PacketConnConfig{
					Network: "udp",
					Address: "127.0.0.0.1",
					SockOption: &zsyscall.SockOption{
						SO:   &zsyscall.SockSOOption{ReceiveBuffer: 5000},
						IP:   &zsyscall.SockIPOption{TTL: 20},
						IPV6: &zsyscall.SockIPV6Option{},
						UDP:  &zsyscall.SockUDPOption{GRO: true},
					},
				},
			},
			&action{
				address: "",
				err: &er.Error{
					Package:     ErrPkg,
					Type:        ErrTypePackConn,
					Description: ErrDscPackConn,
				},
			},
		),
		gen(
			"invalid unixgram",
			[]string{},
			[]string{},
			&condition{
				c: &PacketConnConfig{
					Network:    "unixgram",
					Address:    ` !"#$%&\'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[]^_abcdefghijklmnopqrstuvwxyz{|}~`,
					SockOption: &zsyscall.SockOption{},
				},
			},
			&action{
				address: "",
				err: &er.Error{
					Package:     ErrPkg,
					Type:        ErrTypePackConn,
					Description: ErrDscPackConn,
				},
			},
		),
		gen(
			"invalid network",
			[]string{},
			[]string{},
			&condition{
				c: &PacketConnConfig{
					Network: "invalid",
					Address: "127.0.0.1:12358",
				},
			},
			&action{
				address: "",
				err: &er.Error{
					Package:     ErrPkg,
					Type:        ErrTypePackConn,
					Description: ErrDscPackConn,
				},
			},
		),
		// This case success in macos.
		// gen(
		// 	"control func error",
		// 	[]string{},
		// 	[]string{},
		// 	&condition{
		// 		c: &PacketConnConfig{
		// 			Network: "udp",
		// 			Address: testUDPAddr,
		// 			SockOption: &zsyscall.SockOption{
		// 				SO:  &zsyscall.SockSOOption{KeepAlive: true},
		// 				IP:  &zsyscall.SockIPOption{LocalPortRangeLower: -100},
		// 				UDP: &zsyscall.SockUDPOption{Segment: math.MaxInt},
		// 			},
		// 		},
		// 	},
		// 	&action{
		// 		address: "",
		// 		err: &er.Error{
		// 			Package:     ErrPkg,
		// 			Type:        ErrTypePackConn,
		// 			Description: ErrDscPackConn,
		// 		},
		// 	},
		// ),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			conn, err := NewPacketConn(tt.C().c)
			testutil.Diff(t, tt.A().err, err, cmpopts.EquateErrors())

			if tt.A().address != "" {
				testutil.Diff(t, tt.A().address, conn.LocalAddr().String())
				conn.Close()
			} else {
				testutil.Diff(t, nil, conn)
			}
		})
	}
}

type testListener struct {
	addr  net.Addr
	count int

	acceptErrAfter int
	acceptErr      error

	closeErr error

	clientAddr string

	closed bool // Closed will be set to true when the Close was called.
}

func (l *testListener) Accept() (net.Conn, error) {
	l.count += 1
	if l.acceptErrAfter < l.count && l.acceptErr != nil {
		return nil, l.acceptErr
	}
	return &testConn{
		addr: l.clientAddr,
	}, nil
}

func (l *testListener) Close() error {
	l.closed = true
	return l.closeErr
}

func (l *testListener) Addr() net.Addr {
	return l.addr
}

func TestNewListenerFromSpec(t *testing.T) {
	type condition struct {
		spec *k.ListenConfig
	}

	type action struct {
		address  string
		listener net.Listener
		err      error
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	checker, _ := net.Listen("tcp", "127.0.0.1:0")
	checker.Close()
	testAddr := checker.Addr().String()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"nil spec",
			[]string{},
			[]string{},
			&condition{
				spec: nil,
			},
			&action{
				listener: nil,
				err:      nil,
			},
		),
		gen(
			"valid tls config",
			[]string{},
			[]string{},
			&condition{
				spec: &k.ListenConfig{
					Addr:            "tcp://" + testAddr,
					ConnectionLimit: 10,
					Networks:        []string{"127.0.0.1/32"},
					ReadDeadline:    10,
					WriteDeadline:   20,
				},
			},
			&action{
				address: testAddr,
				err:     nil,
			},
		),
		gen(
			"invalid tls config",
			[]string{},
			[]string{},
			&condition{
				spec: &k.ListenConfig{
					TLSConfig: &k.TLSConfig{
						ClientAuth: k.ClientAuthType(9999), // Invalid value.
					},
				},
			},
			&action{
				err: &er.Error{
					Package:     ErrPkg,
					Type:        ErrTypeListener,
					Description: ErrDscListener,
				},
			},
		),
		gen(
			"with keep-alive disabled",
			[]string{},
			[]string{},
			&condition{
				spec: &k.ListenConfig{
					Addr: "tcp://" + testAddr,
					KeepAliveConfig: &k.KeepAliveConfig{
						Disable: true,
					},
				},
			},
			&action{
				address: testAddr,
				err:     nil,
			},
		),
		gen(
			"with keep-alive enabled",
			[]string{},
			[]string{},
			&condition{
				spec: &k.ListenConfig{
					Addr: "tcp://" + testAddr,
					KeepAliveConfig: &k.KeepAliveConfig{
						Disable:  false,
						Idle:     1,
						Interval: 2,
						Count:    3,
					},
				},
			},
			&action{
				address: testAddr,
				err:     nil,
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			ln, err := NewListenerFromSpec(tt.C().spec)
			testutil.Diff(t, tt.A().err, err, cmpopts.EquateErrors())

			// TODO: better way to check the result.
			// opts := []cmp.Option{
			// 	cmp.AllowUnexported(testListener{}),
			// 	cmp.AllowUnexported(secureListener{}, deadlineListener{}, keepAliveListener{}, limitListener{}),
			// }
			// testutil.Diff(t, tt.A().listener, ln, opts...)

			if tt.A().address != "" {
				testutil.Diff(t, tt.A().address, ln.Addr().String())
				ln.Close()
			} else {
				testutil.Diff(t, nil, ln)
			}
		})
	}
}

func TestNewListener(t *testing.T) {
	type condition struct {
		c *ListenConfig
	}

	type action struct {
		address string
		err     error
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	checker, _ := net.Listen("tcp", "127.0.0.1:0")
	checker.Close()
	testAddr := checker.Addr().String()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"nil config",
			[]string{},
			[]string{},
			&condition{
				c: nil,
			},
			&action{
				err: &er.Error{
					Package:     ErrPkg,
					Type:        ErrTypeListener,
					Description: ErrDscListener,
				},
			},
		),
		gen(
			"valid tcp",
			[]string{},
			[]string{},
			&condition{
				c: &ListenConfig{
					Address: "tcp://" + testAddr,
					SockOption: &zsyscall.SockOption{
						SO:   &zsyscall.SockSOOption{ReceiveBuffer: 5000},
						IP:   &zsyscall.SockIPOption{TTL: 20},
						IPV6: &zsyscall.SockIPV6Option{},
						TCP:  &zsyscall.SockTCPOption{NoDelay: true},
					},
				},
			},
			&action{
				address: testAddr,
				err:     nil,
			},
		),
		gen(
			"valid tcp with TLS",
			[]string{},
			[]string{},
			&condition{
				c: &ListenConfig{
					Address: "tcp://" + testAddr,
					SockOption: &zsyscall.SockOption{
						SO:   &zsyscall.SockSOOption{ReceiveBuffer: 5000},
						IP:   &zsyscall.SockIPOption{TTL: 20},
						IPV6: &zsyscall.SockIPV6Option{},
						TCP:  &zsyscall.SockTCPOption{NoDelay: true},
					},
					TLSConfig: &tls.Config{
						ServerName: "test",
					},
				},
			},
			&action{
				address: testAddr,
				err:     nil,
			},
		),
		gen(
			"invalid tcp",
			[]string{},
			[]string{},
			&condition{
				c: &ListenConfig{
					Address:    "tcp://" + "127.0.0.0.0.1",
					SockOption: &zsyscall.SockOption{},
				},
			},
			&action{
				address: "",
				err: &er.Error{
					Package:     ErrPkg,
					Type:        ErrTypeListener,
					Description: ErrDscListener,
				},
			},
		),
		gen(
			"invalid tcp with TLS",
			[]string{},
			[]string{},
			&condition{
				c: &ListenConfig{
					Address:    "tcp://" + "127.0.0.0.0.1",
					SockOption: &zsyscall.SockOption{},
					TLSConfig: &tls.Config{
						ServerName: "test",
					},
				},
			},
			&action{
				address: "",
				err: &er.Error{
					Package:     ErrPkg,
					Type:        ErrTypeListener,
					Description: ErrDscListener,
				},
			},
		),
		gen(
			"valid unix",
			[]string{},
			[]string{},
			&condition{
				c: &ListenConfig{
					Address:    "unix://" + "@test",
					SockOption: &zsyscall.SockOption{},
				},
			},
			&action{
				address: "@test",
				err:     nil,
			},
		),
		// This case fails on macos.
		// gen(
		// 	"valid unix with TLS",
		// 	[]string{},
		// 	[]string{},
		// 	&condition{
		// 		c: &ListenConfig{
		// 			Network:    "unix",
		// 			Address:    "@test",
		// 			SockOption: &zsyscall.SockOption{},
		// 			TLSConfig: &tls.Config{
		// 				ServerName: "test",
		// 			},
		// 		},
		// 	},
		// 	&action{
		// 		address: "@test",
		// 		err:     nil,
		// 	},
		// ),
		gen(
			"invalid container",
			[]string{},
			[]string{},
			&condition{
				c: &ListenConfig{
					Address:  "tcp://" + testAddr,
					Networks: []string{"127.0.0.0.1"},
				},
			},
			&action{
				address: "",
				err: &er.Error{
					Package:     ErrPkg,
					Type:        ErrTypeListener,
					Description: ErrDscListener,
				},
			},
		),
		gen(
			"invalid network",
			[]string{},
			[]string{},
			&condition{
				c: &ListenConfig{
					Address: "invalid://" + "127.0.0.1:12358",
				},
			},
			&action{
				address: "",
				err: &er.Error{
					Package:     ErrPkg,
					Type:        ErrTypeListener,
					Description: ErrDscListener,
				},
			},
		),
		gen(
			"invalid network with TLS",
			[]string{},
			[]string{},
			&condition{
				c: &ListenConfig{
					Address:   "invalid://" + "127.0.0.1:12358",
					TLSConfig: &tls.Config{},
				},
			},
			&action{
				address: "",
				err: &er.Error{
					Package:     ErrPkg,
					Type:        ErrTypeListener,
					Description: ErrDscListener,
				},
			},
		),
		gen(
			"control error",
			[]string{},
			[]string{},
			&condition{
				c: &ListenConfig{
					Address: "unix://" + "@test",
					SockOption: &zsyscall.SockOption{
						SO:  &zsyscall.SockSOOption{KeepAlive: true},
						TCP: &zsyscall.SockTCPOption{NoDelay: true, FastOpenConnect: true},
					},
				},
			},
			&action{
				address: "",
				err: &er.Error{
					Package:     ErrPkg,
					Type:        ErrTypeListener,
					Description: ErrDscListener,
				},
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			ln, err := NewListener(tt.C().c)
			testutil.Diff(t, tt.A().err, err, cmpopts.EquateErrors())
			// TODO: better way to check the result.
			// opts := []cmp.Option{
			// 	cmp.AllowUnexported(testListener{}),
			// 	cmp.AllowUnexported(secureListener{}, readDeadlineListener{}, writeDeadlineListener{}, limitListener{}),
			// 	cmpopts.IgnoreUnexported(net.TCPListener{}),
			// }
			// testutil.Diff(t, tt.A().listener, ln, opts...)
			if tt.A().address != "" {
				testutil.Diff(t, tt.A().address, ln.Addr().String())
				ln.Close()
			} else {
				testutil.Diff(t, nil, ln)
			}
		})
	}
}

func TestListenerWithReadDeadline(t *testing.T) {
	type condition struct {
		inner    net.Listener
		deadline time.Duration
	}

	type action struct {
		listener net.Listener
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	cndNonNilListener := tb.Condition("listener", "input non nil listener as inner")
	cndNegative := tb.Condition("negative", "input negative value as deadline")
	cndZero := tb.Condition("zero", "input zero as deadline")
	cndPositive := tb.Condition("positive", "input positive value as deadline")
	actCheckInner := tb.Action("inner", "check that the inner listener was returned")
	actCheckDeadline := tb.Action("deadline lister", "check that the deadline listener was returned")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"input nil listener",
			[]string{cndPositive},
			[]string{actCheckInner},
			&condition{
				inner:    nil,
				deadline: 10,
			},
			&action{
				listener: nil,
			},
		),
		gen(
			"negative deadline",
			[]string{cndNegative, cndNonNilListener},
			[]string{actCheckInner},
			&condition{
				inner: &testListener{
					clientAddr: "127.0.0.1:1234",
				},
				deadline: -1,
			},
			&action{
				listener: &readDeadlineListener{
					Listener: &testListener{
						clientAddr: "127.0.0.1:1234",
					},
					read: -1,
				},
			},
		),
		gen(
			"deadline 0",
			[]string{cndZero, cndNonNilListener},
			[]string{actCheckInner},
			&condition{
				inner: &testListener{
					clientAddr: "127.0.0.1:1234",
				},
				deadline: 0,
			},
			&action{
				listener: &testListener{
					clientAddr: "127.0.0.1:1234",
				},
			},
		),
		gen(
			"positive deadline",
			[]string{cndPositive, cndNonNilListener},
			[]string{actCheckDeadline},
			&condition{
				inner: &testListener{
					clientAddr: "127.0.0.1:1234",
				},
				deadline: 10,
			},
			&action{
				listener: &readDeadlineListener{
					Listener: &testListener{
						clientAddr: "127.0.0.1:1234",
					},
					read: 10,
				},
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			l := ListenerWithReadDeadline(tt.C().inner, tt.C().deadline)

			opts := []cmp.Option{
				cmp.AllowUnexported(testListener{}),
				cmp.AllowUnexported(readDeadlineListener{}),
			}
			testutil.Diff(t, tt.A().listener, l, opts...)
		})
	}
}
func TestListenerWithWriteDeadline(t *testing.T) {
	type condition struct {
		inner    net.Listener
		deadline time.Duration
	}

	type action struct {
		listener net.Listener
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	cndNonNilListener := tb.Condition("listener", "input non nil listener as inner")
	cndNegative := tb.Condition("negative", "input negative value as deadline")
	cndZero := tb.Condition("zero", "input zero as deadline")
	cndPositive := tb.Condition("positive", "input positive value as deadline")
	actCheckInner := tb.Action("inner", "check that the inner listener was returned")
	actCheckDeadline := tb.Action("deadline lister", "check that the deadline listener was returned")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"input nil listener",
			[]string{cndPositive},
			[]string{actCheckInner},
			&condition{
				inner:    nil,
				deadline: 10,
			},
			&action{
				listener: nil,
			},
		),
		gen(
			"negative deadline",
			[]string{cndNegative, cndNonNilListener},
			[]string{actCheckInner},
			&condition{
				inner: &testListener{
					clientAddr: "127.0.0.1:1234",
				},
				deadline: -1,
			},
			&action{
				listener: &writeDeadlineListener{
					Listener: &testListener{
						clientAddr: "127.0.0.1:1234",
					},
					write: -1,
				},
			},
		),
		gen(
			"deadline 0",
			[]string{cndZero, cndNonNilListener},
			[]string{actCheckInner},
			&condition{
				inner: &testListener{
					clientAddr: "127.0.0.1:1234",
				},
				deadline: 0,
			},
			&action{
				listener: &testListener{
					clientAddr: "127.0.0.1:1234",
				},
			},
		),
		gen(
			"positive deadline",
			[]string{cndPositive, cndNonNilListener},
			[]string{actCheckDeadline},
			&condition{
				inner: &testListener{
					clientAddr: "127.0.0.1:1234",
				},
				deadline: 10,
			},
			&action{
				listener: &writeDeadlineListener{
					Listener: &testListener{
						clientAddr: "127.0.0.1:1234",
					},
					write: 10,
				},
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			l := ListenerWithWriteDeadline(tt.C().inner, tt.C().deadline)

			opts := []cmp.Option{
				cmp.AllowUnexported(testListener{}),
				cmp.AllowUnexported(writeDeadlineListener{}),
			}
			testutil.Diff(t, tt.A().listener, l, opts...)
		})
	}
}

func TestReadDeadlineListener(t *testing.T) {
	type condition struct {
		listener *readDeadlineListener
	}

	type action struct {
		deadline  time.Time
		acceptErr error
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	cndPositive := tb.Condition("positive", "set positive deadline")
	cndZero := tb.Condition("zero", "set zero deadline")
	cndNegative := tb.Condition("negative", "set negative deadline")
	actCheckDead := tb.Action("check deadline", "check that the deadline was set")
	actCheckError := tb.Action("check error", "check that non nil error was returned")
	table := tb.Build()

	testError := errors.New("test error")

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"positive deadline",
			[]string{cndPositive},
			[]string{actCheckDead},
			&condition{
				listener: &readDeadlineListener{
					Listener: &testListener{
						clientAddr: "127.0.0.1:1234",
					},
					read: 100 * time.Second,
				},
			},
			&action{
				deadline: time.Now().Add(100 * time.Second), // Allow +-5 seconds.
			},
		),
		gen(
			"zero deadline",
			[]string{cndZero},
			[]string{actCheckDead},
			&condition{
				listener: &readDeadlineListener{
					Listener: &testListener{
						clientAddr: "127.0.0.1:1234",
					},
					read: 0,
				},
			},
			&action{
				deadline: time.Time{},
			},
		),
		gen(
			"negative deadline",
			[]string{cndNegative},
			[]string{actCheckDead},
			&condition{
				listener: &readDeadlineListener{
					Listener: &testListener{
						clientAddr: "127.0.0.1:1234",
					},
					read: -1,
				},
			},
			&action{
				deadline: time.Time{},
			},
		),
		gen(
			"error on accept",
			[]string{},
			[]string{actCheckError},
			&condition{
				listener: &readDeadlineListener{
					Listener: &testListener{
						clientAddr: "127.0.0.1:1234",
						acceptErr:  testError,
					},
				},
			},
			&action{
				acceptErr: testError,
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			c, err := tt.C().listener.Accept()
			testutil.Diff(t, tt.A().acceptErr, err, cmpopts.EquateErrors())
			if err == nil {
				cc := c.(*testConn)
				testutil.Diff(t, tt.A().deadline, cc.readDead, cmpopts.EquateApproxTime(5*time.Second))
			}
		})
	}
}

func TestWriteDeadlineListener(t *testing.T) {
	type condition struct {
		listener *writeDeadlineListener
	}

	type action struct {
		deadline  time.Time
		acceptErr error
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	cndPositive := tb.Condition("positive", "set positive deadline")
	cndZero := tb.Condition("zero", "set zero deadline")
	cndNegative := tb.Condition("negative", "set negative deadline")
	actCheckDead := tb.Action("check deadline", "check that the deadline was set")
	actCheckError := tb.Action("check error", "check that non nil error was returned")
	table := tb.Build()

	testError := errors.New("test error")

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"positive deadline",
			[]string{cndPositive},
			[]string{actCheckDead},
			&condition{
				listener: &writeDeadlineListener{
					Listener: &testListener{
						clientAddr: "127.0.0.1:1234",
					},
					write: 100 * time.Second,
				},
			},
			&action{
				deadline: time.Now().Add(100 * time.Second), // Allow +-5 seconds.
			},
		),
		gen(
			"zero deadline",
			[]string{cndZero},
			[]string{actCheckDead},
			&condition{
				listener: &writeDeadlineListener{
					Listener: &testListener{
						clientAddr: "127.0.0.1:1234",
					},
					write: 0,
				},
			},
			&action{
				deadline: time.Time{},
			},
		),
		gen(
			"negative deadline",
			[]string{cndNegative},
			[]string{actCheckDead},
			&condition{
				listener: &writeDeadlineListener{
					Listener: &testListener{
						clientAddr: "127.0.0.1:1234",
					},
					write: -1,
				},
			},
			&action{
				deadline: time.Time{},
			},
		),
		gen(
			"error on accept",
			[]string{},
			[]string{actCheckError},
			&condition{
				listener: &writeDeadlineListener{
					Listener: &testListener{
						clientAddr: "127.0.0.1:1234",
						acceptErr:  testError,
					},
				},
			},
			&action{
				acceptErr: testError,
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			c, err := tt.C().listener.Accept()
			testutil.Diff(t, tt.A().acceptErr, err, cmpopts.EquateErrors())
			if err == nil {
				cc := c.(*testConn)
				testutil.Diff(t, tt.A().deadline, cc.writeDead, cmpopts.EquateApproxTime(5*time.Second))
			}
		})
	}
}
