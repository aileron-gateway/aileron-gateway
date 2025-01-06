//go:build !linux && !windows

package network

import (
	"syscall"
	"testing"

	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
)

func TestSockSOOption_Controllers(t *testing.T) {

	type condition struct {
		opt *SockSOOption
	}

	type action struct {
		values map[int]any
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	cndEnabled := tb.Condition("enabled", "enabled")
	actCheckEnabled := tb.Action("check enabled", "check that the option was enabled")
	actCheckNil := tb.Action("check nil", "check that nil was returned")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"no options",
			[]string{},
			[]string{actCheckNil},
			&condition{
				opt: &SockSOOption{},
			},
			&action{
				values: map[int]any{},
			},
		),
		gen(
			"all",
			[]string{cndEnabled},
			[]string{actCheckEnabled},
			&condition{
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
			&action{
				values: map[int]any{},
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {

			tmp1 := setsockoptInt
			tmp2 := setsockoptLinger
			tmp3 := setsockoptTimeval
			defer func() {
				setsockoptInt = tmp1
				setsockoptLinger = tmp2
				setsockoptTimeval = tmp3
			}()

			called := false
			setsockoptInt = func(fd, level, opt, value int) (err error) {
				called = true
				testutil.Diff(t, syscall.SOL_SOCKET, level)
				testutil.Diff(t, tt.A().values[opt], value)
				return nil
			}
			setsockoptLinger = func(fd, level, opt int, value *syscall.Linger) (err error) {
				called = true
				testutil.Diff(t, syscall.SOL_SOCKET, level)
				testutil.Diff(t, tt.A().values[opt], value)
				return nil
			}
			setsockoptTimeval = func(fd, level, opt int, value *syscall.Timeval) (err error) {
				called = true
				testutil.Diff(t, syscall.SOL_SOCKET, level)
				testutil.Diff(t, tt.A().values[opt], value)
				return nil
			}

			for _, f := range tt.C().opt.Controllers() {
				called = false
				f(0)
				testutil.Diff(t, true, called)
			}

		})
	}

}

func TestSockIPOption_Controllers(t *testing.T) {

	type condition struct {
		opt *SockIPOption
	}

	type action struct {
		values map[int]any
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	cndEnabled := tb.Condition("enabled", "enabled")
	actCheckEnabled := tb.Action("check enabled", "check that the option was enabled")
	actCheckNil := tb.Action("check nil", "check that nil was returned")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"no options",
			[]string{},
			[]string{actCheckNil},
			&condition{
				opt: &SockIPOption{},
			},
			&action{
				values: map[int]any{},
			},
		),
		gen(
			"all",
			[]string{cndEnabled},
			[]string{actCheckEnabled},
			&condition{
				opt: &SockIPOption{
					BindAddressNoPort:   true,
					FreeBind:            true,
					LocalPortRangeUpper: 10,
					LocalPortRangeLower: 11,
					Transparent:         true,
					TTL:                 12,
				},
			},
			&action{
				values: map[int]any{},
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {

			tmp1 := setsockoptInt
			tmp2 := setsockoptLinger
			tmp3 := setsockoptTimeval
			defer func() {
				setsockoptInt = tmp1
				setsockoptLinger = tmp2
				setsockoptTimeval = tmp3
			}()

			called := false
			setsockoptInt = func(fd, level, opt, value int) (err error) {
				called = true
				testutil.Diff(t, syscall.IPPROTO_IP, level)
				testutil.Diff(t, tt.A().values[opt], value)
				return nil
			}
			setsockoptLinger = func(fd, level, opt int, value *syscall.Linger) (err error) {
				called = true
				testutil.Diff(t, syscall.IPPROTO_IP, level)
				testutil.Diff(t, tt.A().values[opt], value)
				return nil
			}
			setsockoptTimeval = func(fd, level, opt int, value *syscall.Timeval) (err error) {
				called = true
				testutil.Diff(t, syscall.IPPROTO_IP, level)
				testutil.Diff(t, tt.A().values[opt], value)
				return nil
			}

			for _, f := range tt.C().opt.Controllers() {
				called = false
				f(0)
				testutil.Diff(t, true, called)
			}

		})
	}

}

func TestSockIPV6Option_Controllers(t *testing.T) {

	type condition struct {
		opt *SockIPV6Option
	}

	type action struct {
		values map[int]any
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	cndEnabled := tb.Condition("enabled", "enabled")
	actCheckEnabled := tb.Action("check enabled", "check that the option was enabled")
	actCheckNil := tb.Action("check nil", "check that nil was returned")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"no options",
			[]string{},
			[]string{actCheckNil},
			&condition{
				opt: &SockIPV6Option{},
			},
			&action{
				values: map[int]any{},
			},
		),
		gen(
			"all",
			[]string{cndEnabled},
			[]string{actCheckEnabled},
			&condition{
				opt: &SockIPV6Option{},
			},
			&action{
				values: map[int]any{},
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {

			tmp1 := setsockoptInt
			tmp2 := setsockoptLinger
			tmp3 := setsockoptTimeval
			defer func() {
				setsockoptInt = tmp1
				setsockoptLinger = tmp2
				setsockoptTimeval = tmp3
			}()

			called := false
			setsockoptInt = func(fd, level, opt, value int) (err error) {
				called = true
				testutil.Diff(t, syscall.IPPROTO_IPV6, level)
				testutil.Diff(t, tt.A().values[opt], value)
				return nil
			}
			setsockoptLinger = func(fd, level, opt int, value *syscall.Linger) (err error) {
				called = true
				testutil.Diff(t, syscall.IPPROTO_IPV6, level)
				testutil.Diff(t, tt.A().values[opt], value)
				return nil
			}
			setsockoptTimeval = func(fd, level, opt int, value *syscall.Timeval) (err error) {
				called = true
				testutil.Diff(t, syscall.IPPROTO_IPV6, level)
				testutil.Diff(t, tt.A().values[opt], value)
				return nil
			}

			for _, f := range tt.C().opt.Controllers() {
				called = false
				f(0)
				testutil.Diff(t, true, called)
			}

		})
	}

}

func TestSockTCPOption_Controllers(t *testing.T) {

	type condition struct {
		opt *SockTCPOption
	}

	type action struct {
		values map[int]any
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	cndEnabled := tb.Condition("enabled", "enabled")
	actCheckEnabled := tb.Action("check enabled", "check that the option was enabled")
	actCheckNil := tb.Action("check nil", "check that nil was returned")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"no options",
			[]string{},
			[]string{actCheckNil},
			&condition{
				opt: &SockTCPOption{},
			},
			&action{
				values: map[int]any{},
			},
		),
		gen(
			"all",
			[]string{cndEnabled},
			[]string{actCheckEnabled},
			&condition{
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
			&action{
				values: map[int]any{},
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {

			tmp1 := setsockoptInt
			tmp2 := setsockoptLinger
			tmp3 := setsockoptTimeval
			defer func() {
				setsockoptInt = tmp1
				setsockoptLinger = tmp2
				setsockoptTimeval = tmp3
			}()

			called := false
			setsockoptInt = func(fd, level, opt, value int) (err error) {
				called = true
				testutil.Diff(t, syscall.IPPROTO_TCP, level)
				testutil.Diff(t, tt.A().values[opt], value)
				return nil
			}
			setsockoptLinger = func(fd, level, opt int, value *syscall.Linger) (err error) {
				called = true
				testutil.Diff(t, syscall.IPPROTO_TCP, level)
				testutil.Diff(t, tt.A().values[opt], value)
				return nil
			}
			setsockoptTimeval = func(fd, level, opt int, value *syscall.Timeval) (err error) {
				called = true
				testutil.Diff(t, syscall.IPPROTO_TCP, level)
				testutil.Diff(t, tt.A().values[opt], value)
				return nil
			}

			for _, f := range tt.C().opt.Controllers() {
				called = false
				f(0)
				testutil.Diff(t, true, called)
			}

		})
	}

}

func TestSockUDPOption_Controllers(t *testing.T) {

	type condition struct {
		opt *SockUDPOption
	}

	type action struct {
		values map[int]any
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	cndEnabled := tb.Condition("enabled", "enabled")
	actCheckEnabled := tb.Action("check enabled", "check that the option was enabled")
	actCheckNil := tb.Action("check nil", "check that nil was returned")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"no options",
			[]string{},
			[]string{actCheckNil},
			&condition{
				opt: &SockUDPOption{},
			},
			&action{
				values: map[int]any{},
			},
		),
		gen(
			"all",
			[]string{cndEnabled},
			[]string{actCheckEnabled},
			&condition{
				opt: &SockUDPOption{
					CORK:    true,
					Segment: 10,
					GRO:     true,
				},
			},
			&action{
				values: map[int]any{},
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {

			tmp1 := setsockoptInt
			tmp2 := setsockoptLinger
			tmp3 := setsockoptTimeval
			defer func() {
				setsockoptInt = tmp1
				setsockoptLinger = tmp2
				setsockoptTimeval = tmp3
			}()

			called := false
			setsockoptInt = func(fd, level, opt, value int) (err error) {
				called = true
				testutil.Diff(t, syscall.IPPROTO_UDP, level)
				testutil.Diff(t, tt.A().values[opt], value)
				return nil
			}
			setsockoptLinger = func(fd, level, opt int, value *syscall.Linger) (err error) {
				called = true
				testutil.Diff(t, syscall.IPPROTO_UDP, level)
				testutil.Diff(t, tt.A().values[opt], value)
				return nil
			}
			setsockoptTimeval = func(fd, level, opt int, value *syscall.Timeval) (err error) {
				called = true
				testutil.Diff(t, syscall.IPPROTO_UDP, level)
				testutil.Diff(t, tt.A().values[opt], value)
				return nil
			}

			for _, f := range tt.C().opt.Controllers() {
				called = false
				f(0)
				testutil.Diff(t, true, called)
			}

		})
	}

}
