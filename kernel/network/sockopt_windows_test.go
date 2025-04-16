// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

//go:build windows

package network

import (
	"io"
	"syscall"
	"testing"

	"github.com/aileron-gateway/aileron-gateway/kernel/er"
	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
	"github.com/google/go-cmp/cmp/cmpopts"
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
				values: map[int]any{
					SO_DEBUG:             1,
					syscall.SO_KEEPALIVE: 1,
					syscall.SO_LINGER: &syscall.Linger{
						Onoff:  1,
						Linger: 10,
					},
					syscall.SO_RCVBUF:    12,
					syscall.SO_REUSEADDR: 1,
					syscall.SO_SNDBUF:    16,
				},
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
			setsockoptInt = func(fd syscall.Handle, level, opt, value int) (err error) {
				called = true
				testutil.Diff(t, syscall.SOL_SOCKET, level)
				testutil.Diff(t, tt.A().values[opt], value)
				return nil
			}
			setsockoptLinger = func(fd syscall.Handle, level, opt int, value *syscall.Linger) (err error) {
				called = true
				testutil.Diff(t, syscall.SOL_SOCKET, level)
				testutil.Diff(t, tt.A().values[opt], value)
				return nil
			}
			setsockoptTimeval = func(fd syscall.Handle, level, opt int, value *syscall.Timeval) (err error) {
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

func TestSODebug(t *testing.T) {

	type condition struct {
		enabled bool
		setErr  error
	}

	type action struct {
		shouldNil bool
		level     int
		opt       int
		value     int
		err       error
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
			"disabled",
			[]string{},
			[]string{actCheckNil},
			&condition{
				enabled: false,
			},
			&action{
				shouldNil: true,
			},
		),
		gen(
			"enabled",
			[]string{cndEnabled},
			[]string{actCheckEnabled},
			&condition{
				enabled: true,
			},
			&action{
				shouldNil: false,
				level:     syscall.SOL_SOCKET,
				opt:       SO_DEBUG,
				value:     1,
			},
		),
		gen(
			"error",
			[]string{cndEnabled},
			[]string{actCheckEnabled},
			&condition{
				enabled: true,
				setErr:  io.EOF, // Use dummy error.
			},
			&action{
				shouldNil: false,
				level:     syscall.SOL_SOCKET,
				opt:       SO_DEBUG,
				value:     1,
				err: &er.Error{
					Package:     ErrPkg,
					Type:        ErrTypeSocket,
					Description: ErrDscSockOpts,
				},
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {

			tmp := setsockoptInt
			defer func() {
				setsockoptInt = tmp
			}()

			called := false
			setsockoptInt = func(fd syscall.Handle, level, opt, value int) (err error) {
				called = true
				testutil.Diff(t, level, tt.A().level)
				testutil.Diff(t, opt, tt.A().opt)
				testutil.Diff(t, value, tt.A().value)
				return tt.C().setErr
			}

			cf := soDebug(tt.C().enabled)

			if tt.A().shouldNil {
				testutil.Diff(t, Controller(nil), cf)
				return
			}
			err := cf(0)
			testutil.Diff(t, tt.A().err, err, cmpopts.EquateErrors())
			testutil.Diff(t, true, called)

		})
	}

}

func TestSOKeepAlive(t *testing.T) {

	type condition struct {
		enabled bool
		setErr  error
	}

	type action struct {
		shouldNil bool
		level     int
		opt       int
		value     int
		err       error
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
			"disabled",
			[]string{},
			[]string{actCheckNil},
			&condition{
				enabled: false,
			},
			&action{
				shouldNil: true,
			},
		),
		gen(
			"enabled",
			[]string{cndEnabled},
			[]string{actCheckEnabled},
			&condition{
				enabled: true,
			},
			&action{
				shouldNil: false,
				level:     syscall.SOL_SOCKET,
				opt:       syscall.SO_KEEPALIVE,
				value:     1,
			},
		),
		gen(
			"error",
			[]string{cndEnabled},
			[]string{actCheckEnabled},
			&condition{
				enabled: true,
				setErr:  io.EOF, // Use dummy error.
			},
			&action{
				shouldNil: false,
				level:     syscall.SOL_SOCKET,
				opt:       syscall.SO_KEEPALIVE,
				value:     1,
				err: &er.Error{
					Package:     ErrPkg,
					Type:        ErrTypeSocket,
					Description: ErrDscSockOpts,
				},
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {

			tmp := setsockoptInt
			defer func() {
				setsockoptInt = tmp
			}()

			called := false
			setsockoptInt = func(fd syscall.Handle, level, opt, value int) (err error) {
				called = true
				testutil.Diff(t, level, tt.A().level)
				testutil.Diff(t, opt, tt.A().opt)
				testutil.Diff(t, value, tt.A().value)
				return tt.C().setErr
			}

			cf := soKeepAlive(tt.C().enabled)

			if tt.A().shouldNil {
				testutil.Diff(t, Controller(nil), cf)
				return
			}
			err := cf(0)
			testutil.Diff(t, tt.A().err, err, cmpopts.EquateErrors())
			testutil.Diff(t, true, called)

		})
	}

}

func TestSOLinger(t *testing.T) {

	type condition struct {
		value  int32
		setErr error
	}

	type action struct {
		shouldNil bool
		level     int
		opt       int
		value     *syscall.Linger
		err       error
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	cndNegative := tb.Condition("negative", "input ngative value <0")
	cndZero := tb.Condition("zero", "input zero")
	cndPositive := tb.Condition("positive", "input positive value >0")
	actCheckEnabled := tb.Action("check enabled", "check that the option was enabled")
	actCheckNil := tb.Action("check nil", "check that nil was returned")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"neative",
			[]string{cndNegative},
			[]string{actCheckNil},
			&condition{
				value: -1,
			},
			&action{
				shouldNil: true,
			},
		),
		gen(
			"zero",
			[]string{cndZero},
			[]string{actCheckNil},
			&condition{
				value: 0,
			},
			&action{
				shouldNil: true,
			},
		),
		gen(
			"positive",
			[]string{cndPositive},
			[]string{actCheckEnabled},
			&condition{
				value: 1,
			},
			&action{
				shouldNil: false,
				level:     syscall.SOL_SOCKET,
				opt:       syscall.SO_LINGER,
				value: &syscall.Linger{
					Onoff:  1,
					Linger: 1,
				},
			},
		),
		gen(
			"error",
			[]string{cndPositive},
			[]string{actCheckEnabled},
			&condition{
				value:  1,
				setErr: io.EOF, // Use dummy error.
			},
			&action{
				shouldNil: false,
				level:     syscall.SOL_SOCKET,
				opt:       syscall.SO_LINGER,
				value: &syscall.Linger{
					Onoff:  1,
					Linger: 1,
				},
				err: &er.Error{
					Package:     ErrPkg,
					Type:        ErrTypeSocket,
					Description: ErrDscSockOpts,
				},
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {

			tmp := setsockoptLinger
			defer func() {
				setsockoptLinger = tmp
			}()

			called := false
			setsockoptLinger = func(fd syscall.Handle, level, opt int, value *syscall.Linger) (err error) {
				called = true
				testutil.Diff(t, level, tt.A().level)
				testutil.Diff(t, opt, tt.A().opt)
				testutil.Diff(t, value, tt.A().value)
				return tt.C().setErr
			}

			cf := soLinger(tt.C().value)

			if tt.A().shouldNil {
				testutil.Diff(t, Controller(nil), cf)
				return
			}
			err := cf(0)
			testutil.Diff(t, tt.A().err, err, cmpopts.EquateErrors())
			testutil.Diff(t, true, called)

		})
	}

}

func TestSORcvbuf(t *testing.T) {

	type condition struct {
		value  int
		setErr error
	}

	type action struct {
		shouldNil bool
		level     int
		opt       int
		value     int
		err       error
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	cndNegative := tb.Condition("negative", "input ngative value <0")
	cndZero := tb.Condition("zero", "input zero")
	cndPositive := tb.Condition("positive", "input positive value >0")
	actCheckEnabled := tb.Action("check enabled", "check that the option was enabled")
	actCheckNil := tb.Action("check nil", "check that nil was returned")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"neative",
			[]string{cndNegative},
			[]string{actCheckNil},
			&condition{
				value: -1,
			},
			&action{
				shouldNil: true,
			},
		),
		gen(
			"zero",
			[]string{cndZero},
			[]string{actCheckNil},
			&condition{
				value: 0,
			},
			&action{
				shouldNil: true,
			},
		),
		gen(
			"positive",
			[]string{cndPositive},
			[]string{actCheckEnabled},
			&condition{
				value: 1,
			},
			&action{
				shouldNil: false,
				level:     syscall.SOL_SOCKET,
				opt:       syscall.SO_RCVBUF,
				value:     1,
			},
		),
		gen(
			"error",
			[]string{cndPositive},
			[]string{actCheckEnabled},
			&condition{
				value:  1,
				setErr: io.EOF, // Use dummy error.
			},
			&action{
				shouldNil: false,
				level:     syscall.SOL_SOCKET,
				opt:       syscall.SO_RCVBUF,
				value:     1,
				err: &er.Error{
					Package:     ErrPkg,
					Type:        ErrTypeSocket,
					Description: ErrDscSockOpts,
				},
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {

			tmp := setsockoptInt
			defer func() {
				setsockoptInt = tmp
			}()

			called := false
			setsockoptInt = func(fd syscall.Handle, level, opt, value int) (err error) {
				called = true
				testutil.Diff(t, level, tt.A().level)
				testutil.Diff(t, opt, tt.A().opt)
				testutil.Diff(t, value, tt.A().value)
				return tt.C().setErr
			}

			cf := soRcvbuf(tt.C().value)

			if tt.A().shouldNil {
				testutil.Diff(t, Controller(nil), cf)
				return
			}
			err := cf(0)
			testutil.Diff(t, tt.A().err, err, cmpopts.EquateErrors())
			testutil.Diff(t, true, called)

		})
	}

}

func TestSOSndbuf(t *testing.T) {

	type condition struct {
		value  int
		setErr error
	}

	type action struct {
		shouldNil bool
		level     int
		opt       int
		value     int
		err       error
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	cndNegative := tb.Condition("negative", "input ngative value <0")
	cndZero := tb.Condition("zero", "input zero")
	cndPositive := tb.Condition("positive", "input positive value >0")
	actCheckEnabled := tb.Action("check enabled", "check that the option was enabled")
	actCheckNil := tb.Action("check nil", "check that nil was returned")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"neative",
			[]string{cndNegative},
			[]string{actCheckNil},
			&condition{
				value: -1,
			},
			&action{
				shouldNil: true,
			},
		),
		gen(
			"zero",
			[]string{cndZero},
			[]string{actCheckNil},
			&condition{
				value: 0,
			},
			&action{
				shouldNil: true,
			},
		),
		gen(
			"positive",
			[]string{cndPositive},
			[]string{actCheckEnabled},
			&condition{
				value: 1,
			},
			&action{
				shouldNil: false,
				level:     syscall.SOL_SOCKET,
				opt:       syscall.SO_SNDBUF,
				value:     1,
			},
		),
		gen(
			"error",
			[]string{cndPositive},
			[]string{actCheckEnabled},
			&condition{
				value:  1,
				setErr: io.EOF, // Use dummy error.
			},
			&action{
				shouldNil: false,
				level:     syscall.SOL_SOCKET,
				opt:       syscall.SO_SNDBUF,
				value:     1,
				err: &er.Error{
					Package:     ErrPkg,
					Type:        ErrTypeSocket,
					Description: ErrDscSockOpts,
				},
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {

			tmp := setsockoptInt
			defer func() {
				setsockoptInt = tmp
			}()

			called := false
			setsockoptInt = func(fd syscall.Handle, level, opt, value int) (err error) {
				called = true
				testutil.Diff(t, level, tt.A().level)
				testutil.Diff(t, opt, tt.A().opt)
				testutil.Diff(t, value, tt.A().value)
				return tt.C().setErr
			}

			cf := soSndbuf(tt.C().value)

			if tt.A().shouldNil {
				testutil.Diff(t, Controller(nil), cf)
				return
			}
			err := cf(0)
			testutil.Diff(t, tt.A().err, err, cmpopts.EquateErrors())
			testutil.Diff(t, true, called)

		})
	}

}

func TestSOReuseaddr(t *testing.T) {

	type condition struct {
		enabled bool
		setErr  error
	}

	type action struct {
		shouldNil bool
		level     int
		opt       int
		value     int
		err       error
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
			"disabled",
			[]string{},
			[]string{actCheckNil},
			&condition{
				enabled: false,
			},
			&action{
				shouldNil: true,
			},
		),
		gen(
			"enabled",
			[]string{cndEnabled},
			[]string{actCheckEnabled},
			&condition{
				enabled: true,
			},
			&action{
				shouldNil: false,
				level:     syscall.SOL_SOCKET,
				opt:       syscall.SO_REUSEADDR,
				value:     1,
			},
		),
		gen(
			"error",
			[]string{cndEnabled},
			[]string{actCheckEnabled},
			&condition{
				enabled: true,
				setErr:  io.EOF, // Use dummy error.
			},
			&action{
				shouldNil: false,
				level:     syscall.SOL_SOCKET,
				opt:       syscall.SO_REUSEADDR,
				value:     1,
				err: &er.Error{
					Package:     ErrPkg,
					Type:        ErrTypeSocket,
					Description: ErrDscSockOpts,
				},
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {

			tmp := setsockoptInt
			defer func() {
				setsockoptInt = tmp
			}()

			called := false
			setsockoptInt = func(fd syscall.Handle, level, opt, value int) (err error) {
				called = true
				testutil.Diff(t, level, tt.A().level)
				testutil.Diff(t, opt, tt.A().opt)
				testutil.Diff(t, value, tt.A().value)
				return tt.C().setErr
			}

			cf := soReuseaddr(tt.C().enabled)

			if tt.A().shouldNil {
				testutil.Diff(t, Controller(nil), cf)
				return
			}
			err := cf(0)
			testutil.Diff(t, tt.A().err, err, cmpopts.EquateErrors())
			testutil.Diff(t, true, called)

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
				values: map[int]any{
					syscall.IP_TTL: 12,
				},
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
			setsockoptInt = func(fd syscall.Handle, level, opt, value int) (err error) {
				called = true
				testutil.Diff(t, syscall.IPPROTO_IP, level)
				testutil.Diff(t, tt.A().values[opt], value)
				return nil
			}
			setsockoptLinger = func(fd syscall.Handle, level, opt int, value *syscall.Linger) (err error) {
				called = true
				testutil.Diff(t, syscall.IPPROTO_IP, level)
				testutil.Diff(t, tt.A().values[opt], value)
				return nil
			}
			setsockoptTimeval = func(fd syscall.Handle, level, opt int, value *syscall.Timeval) (err error) {
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

func TestIPTTL(t *testing.T) {

	type condition struct {
		value  int
		setErr error
	}

	type action struct {
		shouldNil bool
		level     int
		opt       int
		value     int
		err       error
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	cndNegative := tb.Condition("negative", "input ngative value <0")
	cndZero := tb.Condition("zero", "input zero")
	cndPositive := tb.Condition("positive", "input positive value >0")
	actCheckEnabled := tb.Action("check enabled", "check that the option was enabled")
	actCheckNil := tb.Action("check nil", "check that nil was returned")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"neative",
			[]string{cndNegative},
			[]string{actCheckNil},
			&condition{
				value: -1,
			},
			&action{
				shouldNil: true,
			},
		),
		gen(
			"zero",
			[]string{cndZero},
			[]string{actCheckNil},
			&condition{
				value: 0,
			},
			&action{
				shouldNil: true,
			},
		),
		gen(
			"positive",
			[]string{cndPositive},
			[]string{actCheckEnabled},
			&condition{
				value: 1,
			},
			&action{
				shouldNil: false,
				level:     syscall.IPPROTO_IP,
				opt:       syscall.IP_TTL,
				value:     1,
			},
		),
		gen(
			"error",
			[]string{cndPositive},
			[]string{actCheckEnabled},
			&condition{
				value:  1,
				setErr: io.EOF, // Use dummy error.
			},
			&action{
				shouldNil: false,
				level:     syscall.IPPROTO_IP,
				opt:       syscall.IP_TTL,
				value:     1,
				err: &er.Error{
					Package:     ErrPkg,
					Type:        ErrTypeSocket,
					Description: ErrDscSockOpts,
				},
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {

			tmp := setsockoptInt
			defer func() {
				setsockoptInt = tmp
			}()

			called := false
			setsockoptInt = func(fd syscall.Handle, level, opt, value int) (err error) {
				called = true
				testutil.Diff(t, level, tt.A().level)
				testutil.Diff(t, opt, tt.A().opt)
				testutil.Diff(t, value, tt.A().value)
				return tt.C().setErr
			}

			cf := ipTTL(tt.C().value)

			if tt.A().shouldNil {
				testutil.Diff(t, Controller(nil), cf)
				return
			}
			err := cf(0)
			testutil.Diff(t, tt.A().err, err, cmpopts.EquateErrors())
			testutil.Diff(t, true, called)

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
			setsockoptInt = func(fd syscall.Handle, level, opt, value int) (err error) {
				called = true
				testutil.Diff(t, syscall.IPPROTO_IPV6, level)
				testutil.Diff(t, tt.A().values[opt], value)
				return nil
			}
			setsockoptLinger = func(fd syscall.Handle, level, opt int, value *syscall.Linger) (err error) {
				called = true
				testutil.Diff(t, syscall.IPPROTO_IPV6, level)
				testutil.Diff(t, tt.A().values[opt], value)
				return nil
			}
			setsockoptTimeval = func(fd syscall.Handle, level, opt int, value *syscall.Timeval) (err error) {
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
				values: map[int]any{
					syscall.TCP_NODELAY: 1,
				},
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
			setsockoptInt = func(fd syscall.Handle, level, opt, value int) (err error) {
				called = true
				testutil.Diff(t, syscall.IPPROTO_TCP, level)
				testutil.Diff(t, tt.A().values[opt], value)
				return nil
			}
			setsockoptLinger = func(fd syscall.Handle, level, opt int, value *syscall.Linger) (err error) {
				called = true
				testutil.Diff(t, syscall.IPPROTO_TCP, level)
				testutil.Diff(t, tt.A().values[opt], value)
				return nil
			}
			setsockoptTimeval = func(fd syscall.Handle, level, opt int, value *syscall.Timeval) (err error) {
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

func TestTCPNoDelay(t *testing.T) {

	type condition struct {
		enabled bool
		setErr  error
	}

	type action struct {
		shouldNil bool
		level     int
		opt       int
		value     int
		err       error
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
			"disabled",
			[]string{},
			[]string{actCheckNil},
			&condition{
				enabled: false,
			},
			&action{
				shouldNil: true,
			},
		),
		gen(
			"enabled",
			[]string{cndEnabled},
			[]string{actCheckEnabled},
			&condition{
				enabled: true,
			},
			&action{
				shouldNil: false,
				level:     syscall.IPPROTO_TCP,
				opt:       syscall.TCP_NODELAY,
				value:     1,
			},
		),
		gen(
			"error",
			[]string{cndEnabled},
			[]string{actCheckEnabled},
			&condition{
				enabled: true,
				setErr:  io.EOF, // Use dummy error.
			},
			&action{
				shouldNil: false,
				level:     syscall.IPPROTO_TCP,
				opt:       syscall.TCP_NODELAY,
				value:     1,
				err: &er.Error{
					Package:     ErrPkg,
					Type:        ErrTypeSocket,
					Description: ErrDscSockOpts,
				},
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {

			tmp := setsockoptInt
			defer func() {
				setsockoptInt = tmp
			}()

			called := false
			setsockoptInt = func(fd syscall.Handle, level, opt, value int) (err error) {
				called = true
				testutil.Diff(t, level, tt.A().level)
				testutil.Diff(t, opt, tt.A().opt)
				testutil.Diff(t, value, tt.A().value)
				return tt.C().setErr
			}

			cf := tcpNoDelay(tt.C().enabled)

			if tt.A().shouldNil {
				testutil.Diff(t, Controller(nil), cf)
				return
			}
			err := cf(0)
			testutil.Diff(t, tt.A().err, err, cmpopts.EquateErrors())
			testutil.Diff(t, true, called)

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
			setsockoptInt = func(fd syscall.Handle, level, opt, value int) (err error) {
				called = true
				testutil.Diff(t, syscall.IPPROTO_UDP, level)
				testutil.Diff(t, tt.A().values[opt], value)
				return nil
			}
			setsockoptLinger = func(fd syscall.Handle, level, opt int, value *syscall.Linger) (err error) {
				called = true
				testutil.Diff(t, syscall.IPPROTO_UDP, level)
				testutil.Diff(t, tt.A().values[opt], value)
				return nil
			}
			setsockoptTimeval = func(fd syscall.Handle, level, opt int, value *syscall.Timeval) (err error) {
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
