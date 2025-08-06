// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package network

import (
	"crypto/tls"
	"io"
	"net"
	"os"
	"testing"

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
	net.Listener // For testing RemoveSocketListener
	addr         net.Addr
	closeErr     error
	closed       bool
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
