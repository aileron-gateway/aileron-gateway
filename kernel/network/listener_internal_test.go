package network

import (
	stdcmp "cmp"
	"crypto/tls"
	"errors"
	"io"
	"math"
	"net"
	"net/netip"
	"os"
	"testing"
	"time"

	k "github.com/aileron-gateway/aileron-gateway/apis/kernel"
	"github.com/aileron-gateway/aileron-gateway/kernel/er"
	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/pion/dtls/v3"
)

// testDir is the path to the test data.
// This path can be changed by the environmental variable.
var testDir = stdcmp.Or(os.Getenv("TEST_DIR"), "../../test/")

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
					SockOption: &SockOption{
						SO:   &SockSOOption{ReceiveBuffer: 5000},
						IP:   &SockIPOption{TTL: 20},
						IPV6: &SockIPV6Option{},
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
					SockOption: &SockOption{
						SO:   &SockSOOption{ReceiveBuffer: 5000},
						IP:   &SockIPOption{TTL: 20},
						IPV6: &SockIPV6Option{},
						UDP:  &SockUDPOption{GRO: true},
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
					SockOption: &SockOption{
						SO:   &SockSOOption{ReceiveBuffer: 5000},
						IP:   &SockIPOption{TTL: 20},
						IPV6: &SockIPV6Option{},
						UDP:  &SockUDPOption{GRO: true},
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
					SockOption: &SockOption{},
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
		gen(
			"control func error",
			[]string{},
			[]string{},
			&condition{
				c: &PacketConnConfig{
					Network: "udp",
					Address: testUDPAddr,
					SockOption: &SockOption{
						SO:  &SockSOOption{KeepAlive: true},
						IP:  &SockIPOption{LocalPortRangeLower: -100},
						UDP: &SockUDPOption{Segment: math.MaxInt},
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

func mustNetContainers(addresses []string) []Container {
	cs, err := netContainers(addresses)
	if err != nil {
		panic(err)
	}
	return cs
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
			"zero spec",
			[]string{},
			[]string{},
			&condition{
				spec: &k.ListenConfig{},
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
			"valid tls config",
			[]string{},
			[]string{},
			&condition{
				spec: &k.ListenConfig{
					Network:         "tcp",
					Addr:            testAddr,
					ConnectionLimit: 10,
					Networks:        []string{"127.0.0.1/32"},
					Blacklist:       true,
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
					Network: "tcp",
					Address: testAddr,
					SockOption: &SockOption{
						SO:   &SockSOOption{ReceiveBuffer: 5000},
						IP:   &SockIPOption{TTL: 20},
						IPV6: &SockIPV6Option{},
						TCP:  &SockTCPOption{NoDelay: true},
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
					Network: "tcp",
					Address: testAddr,
					SockOption: &SockOption{
						SO:   &SockSOOption{ReceiveBuffer: 5000},
						IP:   &SockIPOption{TTL: 20},
						IPV6: &SockIPV6Option{},
						TCP:  &SockTCPOption{NoDelay: true},
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
					Network:    "tcp",
					Address:    "127.0.0.0.0.1",
					SockOption: &SockOption{},
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
					Network:    "tcp",
					Address:    "127.0.0.0.0.1",
					SockOption: &SockOption{},
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
			"valid udp with TLS",
			[]string{},
			[]string{},
			&condition{
				c: &ListenConfig{
					Network: "udp",
					Address: testAddr,
					SockOption: &SockOption{
						SO:   &SockSOOption{ReceiveBuffer: 5000},
						IP:   &SockIPOption{TTL: 20},
						IPV6: &SockIPV6Option{},
						TCP:  &SockTCPOption{NoDelay: true},
					},
					DTLSConfig: &dtls.Config{
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
			"invalid udp with TLS",
			[]string{},
			[]string{},
			&condition{
				c: &ListenConfig{
					Network: "udp",
					Address: "127.0.0.0.0.1",
					SockOption: &SockOption{
						SO:   &SockSOOption{ReceiveBuffer: 5000},
						IP:   &SockIPOption{TTL: 20},
						IPV6: &SockIPV6Option{},
						TCP:  &SockTCPOption{NoDelay: true},
					},
					DTLSConfig: &dtls.Config{
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
					Network:    "unix",
					Address:    "@test",
					SockOption: &SockOption{},
				},
			},
			&action{
				address: "@test",
				err:     nil,
			},
		),
		gen(
			"valid unix with TLS",
			[]string{},
			[]string{},
			&condition{
				c: &ListenConfig{
					Network:    "unix",
					Address:    "@test",
					SockOption: &SockOption{},
					TLSConfig: &tls.Config{
						ServerName: "test",
					},
				},
			},
			&action{
				address: "@test",
				err:     nil,
			},
		),
		gen(
			"invalid container",
			[]string{},
			[]string{},
			&condition{
				c: &ListenConfig{
					Network:  "tcp",
					Address:  testAddr,
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
					Network: "invalid",
					Address: "127.0.0.1:12358",
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
					Network:   "invalid",
					Address:   "127.0.0.1:12358",
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
			"invalid network with DTLS",
			[]string{},
			[]string{},
			&condition{
				c: &ListenConfig{
					Network:    "invalid",
					Address:    "127.0.0.1:12358",
					DTLSConfig: &dtls.Config{},
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
					Network: "unix",
					Address: "@test",
					SockOption: &SockOption{
						SO:  &SockSOOption{KeepAlive: true},
						TCP: &SockTCPOption{NoDelay: true, FastOpenConnect: true},
					},
					DTLSConfig: &dtls.Config{
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

func TestListenerWithSecure(t *testing.T) {
	type condition struct {
		inner     net.Listener
		list      []Container
		blacklist bool
	}

	type action struct {
		listener net.Listener
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	cndNonNilListener := tb.Condition("listener", "input non nil listener as inner")
	cndWithContainers := tb.Condition("with containers", "input non zero containers")
	actCheckInner := tb.Action("inner", "check that the inner listener was returned")
	actCheckSecure := tb.Action("secure lister", "check that the secure listener was returned")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"input nil listener",
			[]string{},
			[]string{actCheckInner},
			&condition{
				inner: nil,
			},
			&action{
				listener: nil,
			},
		),
		gen(
			"zero containers",
			[]string{cndNonNilListener},
			[]string{actCheckInner},
			&condition{
				inner: &testListener{
					clientAddr: "127.0.0.1:1234",
				},
				list: []Container{},
			},
			&action{
				listener: &testListener{
					clientAddr: "127.0.0.1:1234",
				},
			},
		),
		gen(
			"whitelist listener",
			[]string{cndWithContainers, cndNonNilListener},
			[]string{actCheckSecure},
			&condition{
				inner: &testListener{
					clientAddr: "127.0.0.1:1234",
				},
				list:      mustNetContainers([]string{"127.0.0.1/32:1234"}),
				blacklist: false,
			},
			&action{
				listener: &secureListener{
					Listener: &testListener{
						clientAddr: "127.0.0.1:1234",
					},
					list: mustNetContainers([]string{"127.0.0.1/32:1234"}),
				},
			},
		),
		gen(
			"blacklist listener",
			[]string{cndWithContainers, cndNonNilListener},
			[]string{actCheckSecure},
			&condition{
				inner: &testListener{
					clientAddr: "127.0.0.1:1234",
				},
				list:      mustNetContainers([]string{"127.0.0.1/32:1234"}),
				blacklist: true,
			},
			&action{
				listener: &secureListener{
					Listener: &testListener{
						clientAddr: "127.0.0.1:1234",
					},
					list:      mustNetContainers([]string{"127.0.0.1/32:1234"}),
					blacklist: true,
				},
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			l := ListenerWithSecure(tt.C().inner, tt.C().list, tt.C().blacklist)

			opts := []cmp.Option{
				cmp.AllowUnexported(testListener{}),
				cmp.AllowUnexported(secureListener{}),
			}
			testutil.Diff(t, tt.A().listener, l, opts...)
		})
	}
}

func TestSecureListener(t *testing.T) {
	type condition struct {
		listener *secureListener
	}

	type action struct {
		allowed   bool
		acceptErr error
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	cndWhitelist := tb.Condition("whitelist", "use whitelist")
	cndBlacklist := tb.Condition("blacklist", "use blacklist")
	cndContained := tb.Condition("contained", "client address is contained in the list")
	actCheckAllowed := tb.Action("allowed", "check that the connection was not allowed")
	actCheckNotAllowed := tb.Action("dis-allowed", "check that the connection was allowed")
	table := tb.Build()

	testError := errors.New("test error")

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"whitelist/nil-list",
			[]string{cndWhitelist},
			[]string{actCheckNotAllowed},
			&condition{
				listener: &secureListener{
					Listener: &testListener{
						clientAddr: "127.0.0.1:1234",
					},
					list: nil,
				},
			},
			&action{
				allowed: false,
			},
		),
		gen(
			"whitelist/not-contained",
			[]string{cndWhitelist},
			[]string{actCheckNotAllowed},
			&condition{
				listener: &secureListener{
					Listener: &testListener{
						clientAddr: "127.0.0.1:1234",
					},
					list: mustNetContainers([]string{"123.4.5.6/16"}),
				},
			},
			&action{
				allowed: false,
			},
		),
		gen(
			"whitelist/ip-allowed/port-allowed",
			[]string{cndWhitelist, cndContained},
			[]string{actCheckAllowed},
			&condition{
				listener: &secureListener{
					Listener: &testListener{
						clientAddr: "127.0.0.1:1234",
					},
					list: mustNetContainers([]string{"127.0.0.1/32:1234"}),
				},
			},
			&action{
				allowed: true,
			},
		),
		gen(
			"whitelist/ip-not-allowed/port-allowed",
			[]string{cndWhitelist},
			[]string{actCheckNotAllowed},
			&condition{
				listener: &secureListener{
					Listener: &testListener{
						clientAddr: "127.0.0.1:1234",
					},
					list: mustNetContainers([]string{"127.0.0.2/32:1234"}),
				},
			},
			&action{
				allowed: false,
			},
		),
		gen(
			"whitelist/ip-allowed/port-not-allowed",
			[]string{cndWhitelist},
			[]string{actCheckNotAllowed},
			&condition{
				listener: &secureListener{
					Listener: &testListener{
						clientAddr: "127.0.0.1:1234",
					},
					list: mustNetContainers([]string{"127.0.0.1/32:0-1000"}),
				},
			},
			&action{
				allowed: false,
			},
		),
		gen(
			"blacklist/nil-list",
			[]string{cndBlacklist},
			[]string{actCheckAllowed},
			&condition{
				listener: &secureListener{
					Listener: &testListener{
						clientAddr: "127.0.0.1:1234",
					},
					list:      nil,
					blacklist: true,
				},
			},
			&action{
				allowed: true,
			},
		),
		gen(
			"blacklist/not-contained",
			[]string{cndBlacklist},
			[]string{actCheckAllowed},
			&condition{
				listener: &secureListener{
					Listener: &testListener{
						clientAddr: "127.0.0.1:1234",
					},
					list:      mustNetContainers([]string{"123.4.5.6/16"}),
					blacklist: true,
				},
			},
			&action{
				allowed: true,
			},
		),
		gen(
			"blacklist/ip-not-allowed/port-not-allowed",
			[]string{cndBlacklist, cndContained},
			[]string{actCheckNotAllowed},
			&condition{
				listener: &secureListener{
					Listener: &testListener{
						clientAddr: "127.0.0.1:1234",
					},
					list:      mustNetContainers([]string{"127.0.0.1/32:1234"}),
					blacklist: true,
				},
			},
			&action{
				allowed: false,
			},
		),
		gen(
			"blacklist/ip-allowed/port-not-allowed",
			[]string{cndBlacklist},
			[]string{actCheckAllowed},
			&condition{
				listener: &secureListener{
					Listener: &testListener{
						clientAddr: "127.0.0.1:1234",
					},
					list:      mustNetContainers([]string{"127.0.0.2/32:1234"}),
					blacklist: true,
				},
			},
			&action{
				allowed: true,
			},
		),
		gen(
			"blacklist/ip-not-allowed/port-allowed",
			[]string{cndBlacklist},
			[]string{actCheckAllowed},
			&condition{
				listener: &secureListener{
					Listener: &testListener{
						clientAddr: "127.0.0.1:1234",
					},
					list:      mustNetContainers([]string{"127.0.0.1/32:0-1000"}),
					blacklist: true,
				},
			},
			&action{
				allowed: true,
			},
		),
		gen(
			"invalid client address",
			[]string{},
			[]string{actCheckNotAllowed},
			&condition{
				listener: &secureListener{
					Listener: &testListener{
						clientAddr: "127.0.0.0.1",
					},
				},
			},
			&action{
				allowed:   false,
				acceptErr: nil,
			},
		),
		gen(
			"error on accept",
			[]string{},
			[]string{actCheckNotAllowed},
			&condition{
				listener: &secureListener{
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
				testutil.Diff(t, tt.A().allowed, !cc.closed)
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

func TestListenerWithLimit(t *testing.T) {
	type condition struct {
		inner net.Listener
		limit int
	}

	type action struct {
		listener net.Listener
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	cndNonNilListener := tb.Condition("listener", "input non nil listener as inner")
	cndLimitNegative := tb.Condition("negative limit", "input negative value as limit")
	cndLimitZero := tb.Condition("zero limit", "input zero as limit")
	cndLimitPositive := tb.Condition("positive limit", "input positive value as limit")
	actCheckInner := tb.Action("inner", "check that the inner listener was returned")
	actCheckLimit := tb.Action("limit lister", "check that the limit listener was returned")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"input nil listener",
			[]string{cndLimitPositive},
			[]string{actCheckInner},
			&condition{
				inner: nil,
				limit: 10,
			},
			&action{
				listener: nil,
			},
		),
		gen(
			"limit is -1",
			[]string{cndLimitNegative, cndNonNilListener},
			[]string{actCheckInner},
			&condition{
				inner: &testListener{},
				limit: -1,
			},
			&action{
				listener: &testListener{},
			},
		),
		gen(
			"limit is 0",
			[]string{cndLimitZero, cndNonNilListener},
			[]string{actCheckInner},
			&condition{
				inner: &testListener{},
				limit: 0,
			},
			&action{
				listener: &testListener{},
			},
		),
		gen(
			"limit is 10",
			[]string{cndLimitPositive, cndNonNilListener},
			[]string{actCheckLimit},
			&condition{
				inner: &testListener{
					clientAddr: "127.0.0.1:1234",
				},
				limit: 10,
			},
			&action{
				listener: &limitListener{
					Listener: &testListener{
						clientAddr: "127.0.0.1:1234",
					},
					sem: make(chan struct{}, 10),
				},
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			l := ListenerWithLimit(tt.C().inner, tt.C().limit)

			opts := []cmp.Option{
				cmp.AllowUnexported(testListener{}),
				cmp.AllowUnexported(limitListener{}),
				cmpopts.IgnoreFields(limitListener{}, "sem"),
			}
			testutil.Diff(t, tt.A().listener, l, opts...)
			if ll, ok := l.(*limitListener); ok {
				testutil.Diff(t, tt.C().limit, cap(ll.sem))
			}
		})
	}
}

func TestLimitListener(t *testing.T) {
	type condition struct {
		listener     *limitListener
		successUntil int
	}

	type action struct {
		acceptErr error
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	cndConnect1 := tb.Condition("accept 1", "try to accept 1 connection without block")
	cndConnect5 := tb.Condition("accept 5", "try to accept 5 connection without block")
	cndAcceptError := tb.Condition("accept error", "return an error on accept called")
	actCheckAccepted := tb.Action("accepted", "check that the expected number of connection was accepted without blocking")
	actCheckWaited := tb.Action("waited", "check that the accept was waited when the connection exceeded limits")
	actCheckError := tb.Action("error", "check that there was an error")
	actCheckNoError := tb.Action("no error", "check that there was no error")
	table := tb.Build()

	testError := errors.New("test error")

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"1 connection accepted",
			[]string{cndConnect1},
			[]string{actCheckAccepted, actCheckWaited, actCheckNoError},
			&condition{
				listener: &limitListener{
					Listener: &testListener{
						clientAddr: "127.0.0.1:1234",
						acceptErr:  nil,
					},
					sem: make(chan struct{}, 1),
				},
				successUntil: 1,
			},
			&action{
				acceptErr: nil,
			},
		),
		gen(
			"5 connection accepted",
			[]string{cndConnect5},
			[]string{actCheckAccepted, actCheckWaited, actCheckNoError},
			&condition{
				listener: &limitListener{
					Listener: &testListener{
						clientAddr: "127.0.0.1:1234",
						acceptErr:  nil,
					},
					sem: make(chan struct{}, 5),
				},
				successUntil: 5,
			},
			&action{
				acceptErr: nil,
			},
		),
		gen(
			"error on accept",
			[]string{cndConnect1, cndAcceptError},
			[]string{actCheckAccepted, actCheckWaited, actCheckError},
			&condition{
				listener: &limitListener{
					Listener: &testListener{
						clientAddr:     "127.0.0.1:1234",
						acceptErr:      testError,
						acceptErrAfter: 1,
					},
					sem: make(chan struct{}, 1),
				},
				successUntil: 1,
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
			conns := []net.Conn{}
			for i := 0; i < tt.C().successUntil; i++ {
				before := time.Now()
				c, err := tt.C().listener.Accept() // Should successfully be accepted immediately.
				duration := time.Now().Sub(before)
				testutil.Diff(t, nil, err)
				testutil.Diff(t, true, duration < 5*time.Millisecond) // Should be accepted immediately.
				conns = append(conns, c)
			}

			// Close connection after 50 msec.
			time.AfterFunc(50*time.Millisecond, func() {
				for _, c := range conns {
					c.Close()
				}
			})

			before := time.Now()
			_, err := tt.C().listener.Accept()
			duration := time.Now().Sub(before)
			testutil.Diff(t, true, duration > 50*time.Millisecond)
			testutil.Diff(t, tt.A().acceptErr, err, cmpopts.EquateErrors())
		})
	}
}
