// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package network

import (
	"testing"

	k "github.com/aileron-gateway/aileron-gateway/apis/kernel"
	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
)

func TestSockOptionFromSpec(t *testing.T) {
	type condition struct {
		opt *k.SockOption
	}

	type action struct {
		opt *SockOption
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	cndNil := tb.Condition("nil option", "input nil as an option")
	cndZero := tb.Condition("zero option", "input zero value as an option")
	actCheckNil := tb.Action("check nil", "check that nil was returned")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"nil",
			[]string{cndNil},
			[]string{actCheckNil},
			&condition{
				opt: nil,
			},
			&action{
				opt: nil,
			},
		),
		gen(
			"zero",
			[]string{cndZero},
			[]string{},
			&condition{
				opt: &k.SockOption{},
			},
			&action{
				opt: &SockOption{},
			},
		),
		gen(
			"with so",
			[]string{},
			[]string{},
			&condition{
				opt: &k.SockOption{
					SOOption: &k.SockSOOption{},
				},
			},
			&action{
				opt: &SockOption{
					SO: &SockSOOption{},
				},
			},
		),
		gen(
			"with ip",
			[]string{},
			[]string{},
			&condition{
				opt: &k.SockOption{
					IPOption: &k.SockIPOption{},
				},
			},
			&action{
				opt: &SockOption{
					IP: &SockIPOption{},
				},
			},
		),
		gen(
			"with ipv6",
			[]string{},
			[]string{},
			&condition{
				opt: &k.SockOption{
					IPV6Option: &k.SockIPV6Option{},
				},
			},
			&action{
				opt: &SockOption{
					IPV6: &SockIPV6Option{},
				},
			},
		),
		gen(
			"with tcp",
			[]string{},
			[]string{},
			&condition{
				opt: &k.SockOption{
					TCPOption: &k.SockTCPOption{},
				},
			},
			&action{
				opt: &SockOption{
					TCP: &SockTCPOption{},
				},
			},
		),
		gen(
			"with udp",
			[]string{},
			[]string{},
			&condition{
				opt: &k.SockOption{
					UDPOption: &k.SockUDPOption{},
				},
			},
			&action{
				opt: &SockOption{
					UDP: &SockUDPOption{},
				},
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			opt := SockOptionFromSpec(tt.C().opt)
			testutil.Diff(t, tt.A().opt, opt)
		})
	}
}

func TestSockSOOptionFromSpec(t *testing.T) {
	type condition struct {
		opt *k.SockSOOption
	}

	type action struct {
		opt *SockSOOption
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	cndNil := tb.Condition("nil option", "input nil as an option")
	cndZero := tb.Condition("zero option", "input zero value as an option")
	actCheckNil := tb.Action("check nil", "check that nil was returned")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"nil",
			[]string{cndNil},
			[]string{actCheckNil},
			&condition{
				opt: nil,
			},
			&action{
				opt: nil,
			},
		),
		gen(
			"zero",
			[]string{cndZero},
			[]string{},
			&condition{
				opt: &k.SockSOOption{},
			},
			&action{
				opt: &SockSOOption{},
			},
		),
		gen(
			"all",
			[]string{},
			[]string{},
			&condition{
				opt: &k.SockSOOption{
					BindToDevice:       "eth0",
					Debug:              true,
					IncomingCPU:        true,
					KeepAlive:          true,
					Linger:             10,
					Mark:               11,
					ReceiveBuffer:      12,
					ReceiveBufferForce: 13,
					ReceiveTimeout:     14,
					SendTimeout:        15,
					ReuseAddr:          true,
					ReusePort:          true,
					SendBuffer:         16,
					SendBufferForce:    17,
				},
			},
			&action{
				opt: &SockSOOption{
					BindToDevice:       "eth0",
					Debug:              true,
					IncomingCPU:        true,
					KeepAlive:          true,
					Linger:             10,
					Mark:               11,
					ReceiveBuffer:      12,
					ReceiveBufferForce: 13,
					ReceiveTimeout:     14,
					SendTimeout:        15,
					ReuseAddr:          true,
					ReusePort:          true,
					SendBuffer:         16,
					SendBufferForce:    17,
				},
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			opt := SockSOOptionFromSpec(tt.C().opt)
			testutil.Diff(t, tt.A().opt, opt)
		})
	}
}

func TestSockIPOptionFromSpec(t *testing.T) {
	type condition struct {
		opt *k.SockIPOption
	}

	type action struct {
		opt *SockIPOption
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	cndNil := tb.Condition("nil option", "input nil as an option")
	cndZero := tb.Condition("zero option", "input zero value as an option")
	actCheckNil := tb.Action("check nil", "check that nil was returned")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"nil",
			[]string{cndNil},
			[]string{actCheckNil},
			&condition{
				opt: nil,
			},
			&action{
				opt: nil,
			},
		),
		gen(
			"zero",
			[]string{cndZero},
			[]string{},
			&condition{
				opt: &k.SockIPOption{},
			},
			&action{
				opt: &SockIPOption{},
			},
		),
		gen(
			"all",
			[]string{},
			[]string{},
			&condition{
				opt: &k.SockIPOption{
					BindAddressNoPort:   true,
					FreeBind:            true,
					LocalPortRangeUpper: 10,
					LocalPortRangeLower: 11,
					Transparent:         true,
					TTL:                 12,
				},
			},
			&action{
				opt: &SockIPOption{
					BindAddressNoPort:   true,
					FreeBind:            true,
					LocalPortRangeUpper: 10,
					LocalPortRangeLower: 11,
					Transparent:         true,
					TTL:                 12,
				},
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			opt := SockIPOptionFromSpec(tt.C().opt)
			testutil.Diff(t, tt.A().opt, opt)
		})
	}
}

func TestSockIPV6OptionFromSpec(t *testing.T) {
	type condition struct {
		opt *k.SockIPV6Option
	}

	type action struct {
		opt *SockIPV6Option
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	cndNil := tb.Condition("nil option", "input nil as an option")
	cndZero := tb.Condition("zero option", "input zero value as an option")
	actCheckNil := tb.Action("check nil", "check that nil was returned")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"nil",
			[]string{cndNil},
			[]string{actCheckNil},
			&condition{
				opt: nil,
			},
			&action{
				opt: nil,
			},
		),
		gen(
			"zero",
			[]string{cndZero},
			[]string{},
			&condition{
				opt: &k.SockIPV6Option{},
			},
			&action{
				opt: &SockIPV6Option{},
			},
		),
		gen(
			"all",
			[]string{},
			[]string{},
			&condition{
				opt: &k.SockIPV6Option{},
			},
			&action{
				opt: &SockIPV6Option{},
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			opt := SockIPV6OptionFromSpec(tt.C().opt)
			testutil.Diff(t, tt.A().opt, opt)
		})
	}
}

func TestSockTCPOptionFromSpec(t *testing.T) {
	type condition struct {
		opt *k.SockTCPOption
	}

	type action struct {
		opt *SockTCPOption
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	cndNil := tb.Condition("nil option", "input nil as an option")
	cndZero := tb.Condition("zero option", "input zero value as an option")
	actCheckNil := tb.Action("check nil", "check that nil was returned")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"nil",
			[]string{cndNil},
			[]string{actCheckNil},
			&condition{
				opt: nil,
			},
			&action{
				opt: nil,
			},
		),
		gen(
			"zero",
			[]string{cndZero},
			[]string{},
			&condition{
				opt: &k.SockTCPOption{},
			},
			&action{
				opt: &SockTCPOption{},
			},
		),
		gen(
			"all",
			[]string{},
			[]string{},
			&condition{
				opt: &k.SockTCPOption{
					CORK:            true,
					DeferAccept:     10,
					KeepCount:       11,
					KeepIdle:        12,
					KeepInterval:    13,
					Linger2:         14,
					MaxSegment:      15,
					NoDelay:         true,
					QuickAck:        true,
					SynCount:        16,
					UserTimeout:     17,
					WindowClamp:     18,
					FastOpen:        true,
					FastOpenConnect: true,
				},
			},
			&action{
				opt: &SockTCPOption{
					CORK:            true,
					DeferAccept:     10,
					KeepCount:       11,
					KeepIdle:        12,
					KeepInterval:    13,
					Linger2:         14,
					MaxSegment:      15,
					NoDelay:         true,
					QuickAck:        true,
					SynCount:        16,
					UserTimeout:     17,
					WindowClamp:     18,
					FastOpen:        true,
					FastOpenConnect: true,
				},
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			opt := SockTCPOptionFromSpec(tt.C().opt)
			testutil.Diff(t, tt.A().opt, opt)
		})
	}
}

func TestSockUDPOptionFromSpec(t *testing.T) {
	type condition struct {
		opt *k.SockUDPOption
	}

	type action struct {
		opt *SockUDPOption
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	cndNil := tb.Condition("nil option", "input nil as an option")
	cndZero := tb.Condition("zero option", "input zero value as an option")
	actCheckNil := tb.Action("check nil", "check that nil was returned")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"nil",
			[]string{cndNil},
			[]string{actCheckNil},
			&condition{
				opt: nil,
			},
			&action{
				opt: nil,
			},
		),
		gen(
			"zero",
			[]string{cndZero},
			[]string{},
			&condition{
				opt: &k.SockUDPOption{},
			},
			&action{
				opt: &SockUDPOption{},
			},
		),
		gen(
			"all",
			[]string{},
			[]string{},
			&condition{
				opt: &k.SockUDPOption{
					CORK:    true,
					Segment: 10,
					GRO:     true,
				},
			},
			&action{
				opt: &SockUDPOption{
					CORK:    true,
					Segment: 10,
					GRO:     true,
				},
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			opt := SockUDPOptionFromSpec(tt.C().opt)
			testutil.Diff(t, tt.A().opt, opt)
		})
	}
}
