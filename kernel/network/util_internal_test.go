// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package network

import (
	"math"
	"net"
	"testing"

	"github.com/aileron-gateway/aileron-gateway/kernel/er"
	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestNetContainer_Contains(t *testing.T) {
	type condition struct {
		c     *netContainer
		check []string
		port  int
	}

	type action struct {
		contained bool
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	mustParseIPNet := func(address string) *net.IPNet {
		_, in, err := net.ParseCIDR(address)
		if err != nil {
			panic(err)
		}
		return in
	}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"ipv4 contained/port contained",
			[]string{},
			[]string{},
			&condition{
				c: &netContainer{
					IPNet:    mustParseIPNet("192.168.0.0/16"),
					PortFrom: 8080,
					PortTo:   8080,
				},
				check: []string{"192.168.0.0", "192.168.0.1", "192.168.255.255"},
				port:  8080,
			},
			&action{
				contained: true,
			},
		),
		gen(
			"ipv4 contained/port not contained",
			[]string{},
			[]string{},
			&condition{
				c: &netContainer{
					IPNet:    mustParseIPNet("192.168.0.0/16"),
					PortFrom: 8081,
					PortTo:   math.MaxInt,
				},
				check: []string{"192.168.0.0", "192.168.0.1", "192.168.255.255"},
				port:  8080,
			},
			&action{
				contained: false,
			},
		),
		gen(
			"ipv4 not contained/port contained",
			[]string{},
			[]string{},
			&condition{
				c: &netContainer{
					IPNet:    mustParseIPNet("192.168.0.0/16"),
					PortFrom: 0,
					PortTo:   math.MaxInt,
				},
				check: []string{"0.0.0.0", "255.255.255.255", "10.0.0.1", "10.255.255.254", "172.16.0.1", "172.31.255.254"},
				port:  8080,
			},
			&action{
				contained: false,
			},
		),
		gen(
			"ipv4 allow all",
			[]string{},
			[]string{},
			&condition{
				c: &netContainer{
					IPNet:    mustParseIPNet("0.0.0.0/0"),
					PortFrom: 0,
					PortTo:   math.MaxInt,
				},
				check: []string{"0.0.0.0", "255.255.255.255", "10.0.0.1", "10.255.255.254", "172.16.0.1", "172.31.255.254",
					"::ffff:0.0.0.0", "::ffff:255.255.255.255"}, // IPv4 mapped address.
				port: 8080,
			},
			&action{
				contained: true,
			},
		),
		gen(
			"ipv6 contained/port contained",
			[]string{},
			[]string{},
			&condition{
				c: &netContainer{
					IPNet:    mustParseIPNet("fd00::/8"),
					PortFrom: 8080,
					PortTo:   8080,
				},
				check: []string{"fd00::", "fd00::1", "fd00::ffff"},
				port:  8080,
			},
			&action{
				contained: true,
			},
		),
		gen(
			"ipv6 contained/port not contained",
			[]string{},
			[]string{},
			&condition{
				c: &netContainer{
					IPNet:    mustParseIPNet("fd00::/8"),
					PortFrom: 8081,
					PortTo:   math.MaxInt,
				},
				check: []string{"fd00::", "fd00::1", "fd00::ffff"},
				port:  8080,
			},
			&action{
				contained: false,
			},
		),
		gen(
			"ipv6 not contained/port contained",
			[]string{},
			[]string{},
			&condition{
				c: &netContainer{
					IPNet:    mustParseIPNet("fd00::/8"),
					PortFrom: 0,
					PortTo:   math.MaxInt,
				},
				check: []string{"::", "::1", "::ffff:0.0.0.0", "::ffff:255.255.255.255",
					"::ffff:0:0.0.0.0", "::ffff:0:255.255.255.255", "64:ff9b::0.0.0.0", "64:ff9b::255.255.255.255",
					"64:ff9b:1::", "64:ff9b:1:ffff:ffff:ffff:ffff:ffff", "100::", "100::ffff:ffff:ffff:ffff",
					"2001::", "2001::ffff:ffff:ffff:ffff:ffff:ffff", "2001:20::", "2001:2f:ffff:ffff:ffff:ffff:ffff:ffff",
					"2001:db8::", "2001:db8:ffff:ffff:ffff:ffff:ffff:ffff", "2002::", "2002:ffff:ffff:ffff:ffff:ffff:ffff:ffff",
					"fc00::", "fe80::", "fe80::ffff:ffff:ffff:ffff",
					"ff00::", "ffff:ffff:ffff:ffff:ffff:ffff:ffff:ffff"},
				port: 8080,
			},
			&action{
				contained: false,
			},
		),
		gen(
			"allow all ipv6",
			[]string{},
			[]string{},
			&condition{
				c: &netContainer{
					IPNet:    mustParseIPNet("::/0"),
					PortFrom: 0,
					PortTo:   math.MaxInt,
				},
				check: []string{"::", "::1",
					"::ffff:0:0.0.0.0", "::ffff:0:255.255.255.255", "64:ff9b::0.0.0.0", "64:ff9b::255.255.255.255",
					"64:ff9b:1::", "64:ff9b:1:ffff:ffff:ffff:ffff:ffff", "100::", "100::ffff:ffff:ffff:ffff",
					"2001::", "2001::ffff:ffff:ffff:ffff:ffff:ffff", "2001:20::", "2001:2f:ffff:ffff:ffff:ffff:ffff:ffff",
					"2001:db8::", "2001:db8:ffff:ffff:ffff:ffff:ffff:ffff", "2002::", "2002:ffff:ffff:ffff:ffff:ffff:ffff:ffff",
					"fc00::", "fe80::", "fe80::ffff:ffff:ffff:ffff",
					"ff00::", "ffff:ffff:ffff:ffff:ffff:ffff:ffff:ffff"},
				port: 8080,
			},
			&action{
				contained: true,
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			for _, ip := range tt.C().check {
				contained := tt.C().c.Contains(net.ParseIP(ip), tt.C().port)
				t.Log(ip, contained)
				testutil.Diff(t, tt.A().contained, contained)
			}
		})
	}
}

func TestNetContainers(t *testing.T) {
	type condition struct {
		addresses []string
	}

	type action struct {
		cs  []Container
		err error
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	mustParseIPNet := func(address string) *net.IPNet {
		_, in, err := net.ParseCIDR(address)
		if err != nil {
			panic(err)
		}
		return in
	}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"input nil",
			[]string{},
			[]string{},
			&condition{
				addresses: nil,
			},
			&action{
				cs:  []Container{},
				err: nil,
			},
		),
		gen(
			"ipv4 wo port",
			[]string{},
			[]string{},
			&condition{
				addresses: []string{"192.168.0.0/16"},
			},
			&action{
				cs: []Container{
					&netContainer{
						IPNet:    mustParseIPNet("192.168.0.0/16"),
						PortFrom: 0,
						PortTo:   math.MaxInt,
					},
				},
				err: nil,
			},
		),
		gen(
			"ipv4 w port",
			[]string{},
			[]string{},
			&condition{
				addresses: []string{"192.168.0.0/16:8080"},
			},
			&action{
				cs: []Container{
					&netContainer{
						IPNet:    mustParseIPNet("192.168.0.0/16"),
						PortFrom: 8080,
						PortTo:   8080,
					},
				},
				err: nil,
			},
		),
		gen(
			"ipv4 w port range",
			[]string{},
			[]string{},
			&condition{
				addresses: []string{"192.168.0.0/16:8080-9090"},
			},
			&action{
				cs: []Container{
					&netContainer{
						IPNet:    mustParseIPNet("192.168.0.0/16"),
						PortFrom: 8080,
						PortTo:   9090,
					},
				},
				err: nil,
			},
		),
		gen(
			"ipv4 only",
			[]string{},
			[]string{},
			&condition{
				addresses: []string{"192.168.0.1"},
			},
			&action{
				cs: nil,
				err: &er.Error{
					Package:     ErrPkg,
					Type:        ErrTypeContainer,
					Description: ErrDscContainer,
				},
			},
		),
		gen(
			"ipv4 w invalid port",
			[]string{},
			[]string{},
			&condition{
				addresses: []string{"192.168.0.0/16:xxxx-9090"},
			},
			&action{
				cs: nil,
				err: &er.Error{
					Package:     ErrPkg,
					Type:        ErrTypeContainer,
					Description: ErrDscContainer,
				},
			},
		),
		gen(
			"ipv4 w invalid port",
			[]string{},
			[]string{},
			&condition{
				addresses: []string{"192.168.0.0/16:8080-xxxx"},
			},
			&action{
				cs: nil,
				err: &er.Error{
					Package:     ErrPkg,
					Type:        ErrTypeContainer,
					Description: ErrDscContainer,
				},
			},
		),
		gen(
			"ipv6 wo port",
			[]string{},
			[]string{},
			&condition{
				addresses: []string{"fc00::/7"},
			},
			&action{
				cs: []Container{
					&netContainer{
						IPNet:    mustParseIPNet("fc00::/7"),
						PortFrom: 0,
						PortTo:   math.MaxInt,
					},
				},
				err: nil,
			},
		),
		gen(
			"ipv6 w port",
			[]string{},
			[]string{},
			&condition{
				addresses: []string{"[fc00::/7]:8080"},
			},
			&action{
				cs: []Container{
					&netContainer{
						IPNet:    mustParseIPNet("fc00::/7"),
						PortFrom: 8080,
						PortTo:   8080,
					},
				},
				err: nil,
			},
		),
		gen(
			"ipv6 w port range",
			[]string{},
			[]string{},
			&condition{
				addresses: []string{"[fc00::/7]:8080-9090"},
			},
			&action{
				cs: []Container{
					&netContainer{
						IPNet:    mustParseIPNet("fc00::/7"),
						PortFrom: 8080,
						PortTo:   9090,
					},
				},
				err: nil,
			},
		),
		gen(
			"ipv6 only",
			[]string{},
			[]string{},
			&condition{
				addresses: []string{"fc00::1"},
			},
			&action{
				cs: nil,
				err: &er.Error{
					Package:     ErrPkg,
					Type:        ErrTypeContainer,
					Description: ErrDscContainer,
				},
			},
		),
		gen(
			"ipv6 w invalid port",
			[]string{},
			[]string{},
			&condition{
				addresses: []string{"[fc00::/7]:xxxx-9090"},
			},
			&action{
				cs: nil,
				err: &er.Error{
					Package:     ErrPkg,
					Type:        ErrTypeContainer,
					Description: ErrDscContainer,
				},
			},
		),
		gen(
			"ipv6 w invalid port",
			[]string{},
			[]string{},
			&condition{
				addresses: []string{"[fc00::/7]:8080-xxxx"},
			},
			&action{
				cs: nil,
				err: &er.Error{
					Package:     ErrPkg,
					Type:        ErrTypeContainer,
					Description: ErrDscContainer,
				},
			},
		),
		gen(
			"multiple networks",
			[]string{},
			[]string{},
			&condition{
				addresses: []string{"192.168.0.0/16:8080", "[fc00::/7]:9090"},
			},
			&action{
				cs: []Container{
					&netContainer{
						IPNet:    mustParseIPNet("192.168.0.0/16"),
						PortFrom: 8080,
						PortTo:   8080,
					},
					&netContainer{
						IPNet:    mustParseIPNet("fc00::/7"),
						PortFrom: 9090,
						PortTo:   9090,
					},
				},
				err: nil,
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			cs, err := netContainers(tt.C().addresses)
			testutil.Diff(t, tt.A().err, err, cmpopts.EquateErrors())
			testutil.Diff(t, tt.A().cs, cs)
		})
	}
}

func TestSplitHostPort(t *testing.T) {
	type condition struct {
		address string
	}

	type action struct {
		ip   net.IP
		port int
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"empty address",
			[]string{},
			[]string{},
			&condition{
				address: "",
			},
			&action{
				ip:   nil,
				port: 0,
			},
		),
		gen(
			"ipv4 wo port",
			[]string{},
			[]string{},
			&condition{
				address: "192.168.0.0",
			},
			&action{
				ip:   net.ParseIP("192.168.0.0"),
				port: 0,
			},
		),
		gen(
			"ipv4 w port",
			[]string{},
			[]string{},
			&condition{
				address: "192.168.0.0:8080",
			},
			&action{
				ip:   net.ParseIP("192.168.0.0"),
				port: 8080,
			},
		),
		gen(
			"invalid ipv4 wo port",
			[]string{},
			[]string{},
			&condition{
				address: "192.168.0.0.1",
			},
			&action{
				ip:   nil,
				port: 0,
			},
		),
		gen(
			"invalid ipv4 w port",
			[]string{},
			[]string{},
			&condition{
				address: "192.168.0.0.1:8080",
			},
			&action{
				ip:   nil,
				port: 8080, // port will be returned.
			},
		),
		gen(
			"ipv6 wo port",
			[]string{},
			[]string{},
			&condition{
				address: "fdff::1",
			},
			&action{
				ip:   net.ParseIP("fdff::1"),
				port: 0,
			},
		),
		gen(
			"ipv6 w port",
			[]string{},
			[]string{},
			&condition{
				address: "[fdff::1]:8080",
			},
			&action{
				ip:   net.ParseIP("fdff::1"),
				port: 8080,
			},
		),
		gen(
			"invalid ipv6 wo port",
			[]string{},
			[]string{},
			&condition{
				address: "fdff::ffff::1",
			},
			&action{
				ip:   nil,
				port: 0,
			},
		),
		gen(
			"invalid ipv6 w port",
			[]string{},
			[]string{},
			&condition{
				address: "[fdff::ffff::1]:8080",
			},
			&action{
				ip:   nil,
				port: 8080, // port will be returned.
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			ip, port := splitHostPort(tt.C().address)
			testutil.Diff(t, tt.A().ip, ip)
			testutil.Diff(t, tt.A().port, port)
		})
	}
}
