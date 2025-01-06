//go:build linux

package network

import (
	"encoding/binary"
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
					syscall.SO_BINDTODEVICE: "eth0",
					syscall.SO_DEBUG:        1,
					SO_INCOMING_CPU:         1,
					syscall.SO_KEEPALIVE:    1,
					syscall.SO_LINGER: &syscall.Linger{
						Onoff:  1,
						Linger: 10,
					},
					syscall.SO_MARK:        11,
					syscall.SO_RCVBUF:      12,
					syscall.SO_RCVBUFFORCE: 13,
					syscall.SO_RCVTIMEO:    &syscall.Timeval{Sec: 14},
					syscall.SO_SNDTIMEO:    &syscall.Timeval{Sec: 15},
					syscall.SO_REUSEADDR:   1,
					SO_REUSEPORT:           1,
					syscall.SO_SNDBUF:      16,
					syscall.SO_SNDBUFFORCE: 17,
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
			tmp4 := setsockoptString
			defer func() {
				setsockoptInt = tmp1
				setsockoptLinger = tmp2
				setsockoptTimeval = tmp3
				setsockoptString = tmp4
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
			setsockoptString = func(fd, level, opt int, value string) (err error) {
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

func TestSOBindToDevice(t *testing.T) {
	type condition struct {
		value  string
		setErr error
	}

	type action struct {
		shouldNil bool
		level     int
		opt       int
		value     string
		err       error
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	cndNonEmpty := tb.Condition("non empty", "input non empty string")
	actCheckEnabled := tb.Action("check enabled", "check that the option was enabled")
	actCheckNil := tb.Action("check nil", "check that nil was returned")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"empty",
			[]string{},
			[]string{actCheckNil},
			&condition{
				value: "",
			},
			&action{
				shouldNil: true,
			},
		),
		gen(
			"non empty",
			[]string{cndNonEmpty},
			[]string{actCheckEnabled},
			&condition{
				value: "eth0",
			},
			&action{
				shouldNil: false,
				level:     syscall.SOL_SOCKET,
				opt:       syscall.SO_BINDTODEVICE,
				value:     "eth0",
			},
		),
		gen(
			"error",
			[]string{cndNonEmpty},
			[]string{actCheckEnabled},
			&condition{
				value:  "eth0",
				setErr: io.EOF, // Use dummy error.
			},
			&action{
				shouldNil: false,
				level:     syscall.SOL_SOCKET,
				opt:       syscall.SO_BINDTODEVICE,
				value:     "eth0",
				err: &er.Error{
					Package:     ErrPkg,
					Type:        ErrTypeSocket,
					Description: ErrDscSockOpts,
					Detail:      "SOL_SOCKET.SO_BINDTODEVICE",
				},
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			tmp := setsockoptString
			defer func() {
				setsockoptString = tmp
			}()

			called := false
			setsockoptString = func(fd, level, opt int, value string) (err error) {
				called = true
				testutil.Diff(t, level, tt.A().level)
				testutil.Diff(t, opt, tt.A().opt)
				testutil.Diff(t, value, tt.A().value)
				return tt.C().setErr
			}

			cf := soBindToDevice(tt.C().value)

			if tt.A().shouldNil {
				testutil.Diff(t, Controller(nil), cf)
				return
			}
			err := cf(0)
			testutil.Diff(t, true, called)
			testutil.Diff(t, tt.A().err, err, cmpopts.EquateErrors())
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
				opt:       syscall.SO_DEBUG,
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
				opt:       syscall.SO_DEBUG,
				value:     1,
				err: &er.Error{
					Package:     ErrPkg,
					Type:        ErrTypeSocket,
					Description: ErrDscSockOpts,
					Detail:      "SOL_SOCKET.SO_DEBUG",
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
			setsockoptInt = func(fd, level, opt, value int) (err error) {
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
			testutil.Diff(t, true, called)
			testutil.Diff(t, tt.A().err, err, cmpopts.EquateErrors())
		})
	}
}

func TestSOIncomingCPU(t *testing.T) {
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
				opt:       SO_INCOMING_CPU,
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
				opt:       SO_INCOMING_CPU,
				value:     1,
				err: &er.Error{
					Package:     ErrPkg,
					Type:        ErrTypeSocket,
					Description: ErrDscSockOpts,
					Detail:      "SOL_SOCKET.SO_INCOMING_CPU",
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
			setsockoptInt = func(fd, level, opt, value int) (err error) {
				called = true
				testutil.Diff(t, level, tt.A().level)
				testutil.Diff(t, opt, tt.A().opt)
				testutil.Diff(t, value, tt.A().value)
				return tt.C().setErr
			}

			cf := soIncomingCPU(tt.C().enabled)

			if tt.A().shouldNil {
				testutil.Diff(t, Controller(nil), cf)
				return
			}
			err := cf(0)
			testutil.Diff(t, true, called)
			testutil.Diff(t, tt.A().err, err, cmpopts.EquateErrors())
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
					Detail:      "SOL_SOCKET.SO_KEEPALIVE",
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
			setsockoptInt = func(fd, level, opt, value int) (err error) {
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
			testutil.Diff(t, true, called)
			testutil.Diff(t, tt.A().err, err, cmpopts.EquateErrors())
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
					Detail:      "SOL_SOCKET.SO_LINGER",
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
			setsockoptLinger = func(fd, level, opt int, value *syscall.Linger) (err error) {
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
			testutil.Diff(t, true, called)
			testutil.Diff(t, tt.A().err, err, cmpopts.EquateErrors())
		})
	}
}

func TestSOMark(t *testing.T) {
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
				opt:       syscall.SO_MARK,
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
				opt:       syscall.SO_MARK,
				value:     1,
				err: &er.Error{
					Package:     ErrPkg,
					Type:        ErrTypeSocket,
					Description: ErrDscSockOpts,
					Detail:      "SOL_SOCKET.SO_MARK",
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
			setsockoptInt = func(fd, level, opt, value int) (err error) {
				called = true
				testutil.Diff(t, level, tt.A().level)
				testutil.Diff(t, opt, tt.A().opt)
				testutil.Diff(t, value, tt.A().value)
				return tt.C().setErr
			}

			cf := soMark(tt.C().value)

			if tt.A().shouldNil {
				testutil.Diff(t, Controller(nil), cf)
				return
			}
			err := cf(0)
			testutil.Diff(t, true, called)
			testutil.Diff(t, tt.A().err, err, cmpopts.EquateErrors())
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
					Detail:      "SOL_SOCKET.SO_RCVBUF",
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
			setsockoptInt = func(fd, level, opt, value int) (err error) {
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
			testutil.Diff(t, true, called)
			testutil.Diff(t, tt.A().err, err, cmpopts.EquateErrors())
		})
	}
}

func TestSORcvbufForce(t *testing.T) {
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
				opt:       syscall.SO_RCVBUFFORCE,
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
				opt:       syscall.SO_RCVBUFFORCE,
				value:     1,
				err: &er.Error{
					Package:     ErrPkg,
					Type:        ErrTypeSocket,
					Description: ErrDscSockOpts,
					Detail:      "SOL_SOCKET.SO_RCVBUFFORCE",
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
			setsockoptInt = func(fd, level, opt, value int) (err error) {
				called = true
				testutil.Diff(t, level, tt.A().level)
				testutil.Diff(t, opt, tt.A().opt)
				testutil.Diff(t, value, tt.A().value)
				return tt.C().setErr
			}

			cf := soRcvbufForce(tt.C().value)

			if tt.A().shouldNil {
				testutil.Diff(t, Controller(nil), cf)
				return
			}
			err := cf(0)
			testutil.Diff(t, true, called)
			testutil.Diff(t, tt.A().err, err, cmpopts.EquateErrors())
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
					Detail:      "SOL_SOCKET.SO_SNDBUF",
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
			setsockoptInt = func(fd, level, opt, value int) (err error) {
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
			testutil.Diff(t, true, called)
			testutil.Diff(t, tt.A().err, err, cmpopts.EquateErrors())
		})
	}
}

func TestSOSndbufForce(t *testing.T) {
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
				opt:       syscall.SO_SNDBUFFORCE,
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
				opt:       syscall.SO_SNDBUFFORCE,
				value:     1,
				err: &er.Error{
					Package:     ErrPkg,
					Type:        ErrTypeSocket,
					Description: ErrDscSockOpts,
					Detail:      "SOL_SOCKET.SO_SNDBUFFORCE",
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
			setsockoptInt = func(fd, level, opt, value int) (err error) {
				called = true
				testutil.Diff(t, level, tt.A().level)
				testutil.Diff(t, opt, tt.A().opt)
				testutil.Diff(t, value, tt.A().value)
				return tt.C().setErr
			}

			cf := soSndbufForce(tt.C().value)

			if tt.A().shouldNil {
				testutil.Diff(t, Controller(nil), cf)
				return
			}
			err := cf(0)
			testutil.Diff(t, true, called)
			testutil.Diff(t, tt.A().err, err, cmpopts.EquateErrors())
		})
	}
}

func TestSOSndtimeo(t *testing.T) {
	type condition struct {
		value  float64
		setErr error
	}

	type action struct {
		shouldNil bool
		level     int
		opt       int
		value     *syscall.Timeval
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
				opt:       syscall.SO_SNDTIMEO,
				value: &syscall.Timeval{
					Sec:  1,
					Usec: 0,
				},
			},
		),
		gen(
			"positive",
			[]string{cndPositive},
			[]string{actCheckEnabled},
			&condition{
				value: 9.123456789,
			},
			&action{
				shouldNil: false,
				level:     syscall.SOL_SOCKET,
				opt:       syscall.SO_SNDTIMEO,
				value: &syscall.Timeval{
					Sec:  9,
					Usec: 123457,
				},
			},
		),
		gen(
			"error",
			[]string{cndPositive},
			[]string{actCheckEnabled},
			&condition{
				value:  9.123456789,
				setErr: io.EOF, // Use dummy error.
			},
			&action{
				shouldNil: false,
				level:     syscall.SOL_SOCKET,
				opt:       syscall.SO_SNDTIMEO,
				value: &syscall.Timeval{
					Sec:  9,
					Usec: 123457,
				},
				err: &er.Error{
					Package:     ErrPkg,
					Type:        ErrTypeSocket,
					Description: ErrDscSockOpts,
					Detail:      "SOL_SOCKET.SO_SNDTIMEO",
				},
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			tmp := setsockoptTimeval
			defer func() {
				setsockoptTimeval = tmp
			}()

			called := false
			setsockoptTimeval = func(fd, level, opt int, value *syscall.Timeval) (err error) {
				called = true
				testutil.Diff(t, tt.A().level, level)
				testutil.Diff(t, tt.A().opt, opt)
				testutil.Diff(t, tt.A().value, value)
				return tt.C().setErr
			}

			cf := soSndtimeo(tt.C().value)

			if tt.A().shouldNil {
				testutil.Diff(t, Controller(nil), cf)
				return
			}
			err := cf(0)
			testutil.Diff(t, true, called)
			testutil.Diff(t, tt.A().err, err, cmpopts.EquateErrors())
		})
	}
}

func TestSORcvtimeo(t *testing.T) {
	type condition struct {
		value  float64
		setErr error
	}

	type action struct {
		shouldNil bool
		level     int
		opt       int
		value     *syscall.Timeval
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
				opt:       syscall.SO_RCVTIMEO,
				value: &syscall.Timeval{
					Sec:  1,
					Usec: 0,
				},
			},
		),
		gen(
			"positive",
			[]string{cndPositive},
			[]string{actCheckEnabled},
			&condition{
				value: 9.123456789,
			},
			&action{
				shouldNil: false,
				level:     syscall.SOL_SOCKET,
				opt:       syscall.SO_RCVTIMEO,
				value: &syscall.Timeval{
					Sec:  9,
					Usec: 123457,
				},
			},
		),
		gen(
			"error",
			[]string{cndPositive},
			[]string{actCheckEnabled},
			&condition{
				value:  9.123456789,
				setErr: io.EOF, // Use dummy error.
			},
			&action{
				shouldNil: false,
				level:     syscall.SOL_SOCKET,
				opt:       syscall.SO_RCVTIMEO,
				value: &syscall.Timeval{
					Sec:  9,
					Usec: 123457,
				},
				err: &er.Error{
					Package:     ErrPkg,
					Type:        ErrTypeSocket,
					Description: ErrDscSockOpts,
					Detail:      "SOL_SOCKET.SO_RCVTIMEO",
				},
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			tmp := setsockoptTimeval
			defer func() {
				setsockoptTimeval = tmp
			}()

			called := false
			setsockoptTimeval = func(fd, level, opt int, value *syscall.Timeval) (err error) {
				called = true
				testutil.Diff(t, tt.A().level, level)
				testutil.Diff(t, tt.A().opt, opt)
				testutil.Diff(t, tt.A().value, value)
				return tt.C().setErr
			}

			cf := soRcvtimeo(tt.C().value)

			if tt.A().shouldNil {
				testutil.Diff(t, Controller(nil), cf)
				return
			}
			err := cf(0)
			testutil.Diff(t, true, called)
			testutil.Diff(t, tt.A().err, err, cmpopts.EquateErrors())
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
					Detail:      "SOL_SOCKET.SO_REUSEADDR",
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
			setsockoptInt = func(fd, level, opt, value int) (err error) {
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
			testutil.Diff(t, true, called)
			testutil.Diff(t, tt.A().err, err, cmpopts.EquateErrors())
		})
	}
}

func TestSOReuseport(t *testing.T) {
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
				opt:       SO_REUSEPORT,
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
				opt:       SO_REUSEPORT,
				value:     1,
				err: &er.Error{
					Package:     ErrPkg,
					Type:        ErrTypeSocket,
					Description: ErrDscSockOpts,
					Detail:      "SOL_SOCKET.SO_REUSEPORT",
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
			setsockoptInt = func(fd, level, opt, value int) (err error) {
				called = true
				testutil.Diff(t, level, tt.A().level)
				testutil.Diff(t, opt, tt.A().opt)
				testutil.Diff(t, value, tt.A().value)
				return tt.C().setErr
			}

			cf := soReuseport(tt.C().enabled)

			if tt.A().shouldNil {
				testutil.Diff(t, Controller(nil), cf)
				return
			}
			err := cf(0)
			testutil.Diff(t, true, called)
			testutil.Diff(t, tt.A().err, err, cmpopts.EquateErrors())
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
					IP_BIND_ADDRESS_NO_PORT: 1,
					syscall.IP_FREEBIND:     1,
					IP_LOCAL_PORT_RANGE:     int(binary.BigEndian.Uint32(binary.NativeEndian.AppendUint16(binary.NativeEndian.AppendUint16(nil, 10), 11))),
					syscall.IP_TRANSPARENT:  1,
					syscall.IP_TTL:          12,
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

func TestIPBindAddressNoPort(t *testing.T) {
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
				level:     syscall.IPPROTO_IP,
				opt:       IP_BIND_ADDRESS_NO_PORT,
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
				level:     syscall.IPPROTO_IP,
				opt:       IP_BIND_ADDRESS_NO_PORT,
				value:     1,
				err: &er.Error{
					Package:     ErrPkg,
					Type:        ErrTypeSocket,
					Description: ErrDscSockOpts,
					Detail:      "IPPROTO_IP.IP_BIND_ADDRESS_NO_PORT",
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
			setsockoptInt = func(fd, level, opt, value int) (err error) {
				called = true
				testutil.Diff(t, level, tt.A().level)
				testutil.Diff(t, opt, tt.A().opt)
				testutil.Diff(t, value, tt.A().value)
				return tt.C().setErr
			}

			cf := ipBindAddressNoPort(tt.C().enabled)

			if tt.A().shouldNil {
				testutil.Diff(t, Controller(nil), cf)
				return
			}
			err := cf(0)
			testutil.Diff(t, true, called)
			testutil.Diff(t, tt.A().err, err, cmpopts.EquateErrors())
		})
	}
}

func TestIPFreeBind(t *testing.T) {
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
				level:     syscall.IPPROTO_IP,
				opt:       syscall.IP_FREEBIND,
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
				level:     syscall.IPPROTO_IP,
				opt:       syscall.IP_FREEBIND,
				value:     1,
				err: &er.Error{
					Package:     ErrPkg,
					Type:        ErrTypeSocket,
					Description: ErrDscSockOpts,
					Detail:      "IPPROTO_IP.IP_FREEBIND",
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
			setsockoptInt = func(fd, level, opt, value int) (err error) {
				called = true
				testutil.Diff(t, level, tt.A().level)
				testutil.Diff(t, opt, tt.A().opt)
				testutil.Diff(t, value, tt.A().value)
				return tt.C().setErr
			}

			cf := ipFreeBind(tt.C().enabled)

			if tt.A().shouldNil {
				testutil.Diff(t, Controller(nil), cf)
				return
			}
			err := cf(0)
			testutil.Diff(t, true, called)
			testutil.Diff(t, tt.A().err, err, cmpopts.EquateErrors())
		})
	}
}

func TestIPLocalPortRange(t *testing.T) {
	type condition struct {
		upper  int
		lower  int
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
	cndEnabled := tb.Condition("enabled", "enabled")
	actCheckEnabled := tb.Action("check enabled", "check that the option was enabled")
	actCheckNil := tb.Action("check nil", "check that nil was returned")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"negative,negative",
			[]string{},
			[]string{actCheckNil},
			&condition{
				upper: -1,
				lower: -1,
			},
			&action{
				shouldNil: true,
			},
		),
		gen(
			"negative,zero",
			[]string{},
			[]string{actCheckNil},
			&condition{
				upper: -1,
				lower: 0,
			},
			&action{
				shouldNil: true,
			},
		),
		gen(
			"zero,neative",
			[]string{},
			[]string{actCheckNil},
			&condition{
				upper: 0,
				lower: -1,
			},
			&action{
				shouldNil: true,
			},
		),
		gen(
			"zero,zero",
			[]string{},
			[]string{actCheckNil},
			&condition{
				upper: 0,
				lower: 0,
			},
			&action{
				shouldNil: true,
			},
		),
		gen(
			"positiv,positive",
			[]string{cndEnabled},
			[]string{actCheckEnabled},
			&condition{
				upper: 100,
				lower: 50,
			},
			&action{
				shouldNil: false,
				level:     syscall.IPPROTO_IP,
				opt:       IP_LOCAL_PORT_RANGE,
				value:     int(binary.BigEndian.Uint32(binary.NativeEndian.AppendUint16(binary.NativeEndian.AppendUint16(nil, 100), 50))),
			},
		),
		gen(
			"error",
			[]string{cndEnabled},
			[]string{actCheckEnabled},
			&condition{
				upper:  100,
				lower:  50,
				setErr: io.EOF, // Use dummy error.
			},
			&action{
				shouldNil: false,
				level:     syscall.IPPROTO_IP,
				opt:       IP_LOCAL_PORT_RANGE,
				value:     int(binary.BigEndian.Uint32(binary.NativeEndian.AppendUint16(binary.NativeEndian.AppendUint16(nil, 100), 50))),
				err: &er.Error{
					Package:     ErrPkg,
					Type:        ErrTypeSocket,
					Description: ErrDscSockOpts,
					Detail:      "IPPROTO_IP.IP_LOCAL_PORT_RANGE",
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
			setsockoptInt = func(fd, level, opt, value int) (err error) {
				called = true
				testutil.Diff(t, level, tt.A().level)
				testutil.Diff(t, opt, tt.A().opt)
				testutil.Diff(t, value, tt.A().value)
				return tt.C().setErr
			}

			cf := ipLocalPortRange(tt.C().upper, tt.C().lower)

			if tt.A().shouldNil {
				testutil.Diff(t, Controller(nil), cf)
				return
			}
			err := cf(0)
			testutil.Diff(t, true, called)
			testutil.Diff(t, tt.A().err, err, cmpopts.EquateErrors())
		})
	}
}

func TestIPTransparent(t *testing.T) {
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
				level:     syscall.IPPROTO_IP,
				opt:       syscall.IP_TRANSPARENT,
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
				level:     syscall.IPPROTO_IP,
				opt:       syscall.IP_TRANSPARENT,
				value:     1,
				err: &er.Error{
					Package:     ErrPkg,
					Type:        ErrTypeSocket,
					Description: ErrDscSockOpts,
					Detail:      "IPPROTO_IP.IP_TRANSPARENT",
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
			setsockoptInt = func(fd, level, opt, value int) (err error) {
				called = true
				testutil.Diff(t, level, tt.A().level)
				testutil.Diff(t, opt, tt.A().opt)
				testutil.Diff(t, value, tt.A().value)
				return tt.C().setErr
			}

			cf := ipTransparent(tt.C().enabled)

			if tt.A().shouldNil {
				testutil.Diff(t, Controller(nil), cf)
				return
			}
			err := cf(0)
			testutil.Diff(t, true, called)
			testutil.Diff(t, tt.A().err, err, cmpopts.EquateErrors())
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
					Detail:      "IPPROTO_IP.IP_TTL",
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
			setsockoptInt = func(fd, level, opt, value int) (err error) {
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
			testutil.Diff(t, true, called)
			testutil.Diff(t, tt.A().err, err, cmpopts.EquateErrors())
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
				values: map[int]any{
					syscall.TCP_CORK:         1,
					syscall.TCP_DEFER_ACCEPT: 10,
					syscall.TCP_KEEPCNT:      11,
					syscall.TCP_KEEPIDLE:     12,
					syscall.TCP_KEEPINTVL:    13,
					syscall.TCP_LINGER2: &syscall.Linger{
						Onoff:  1,
						Linger: 14,
					},
					syscall.TCP_MAXSEG:       15,
					syscall.TCP_NODELAY:      1,
					TCP_QUICKACK:             1,
					TCP_SYNCNT:               16,
					TCP_USER_TIMEOUT:         17,
					syscall.TCP_WINDOW_CLAMP: 18,
					TCP_FASTOPEN:             1,
					TCP_FASTOPEN_CONNECT:     1,
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

func TestTCPCORK(t *testing.T) {
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
				opt:       syscall.TCP_CORK,
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
				opt:       syscall.TCP_CORK,
				value:     1,
				err: &er.Error{
					Package:     ErrPkg,
					Type:        ErrTypeSocket,
					Description: ErrDscSockOpts,
					Detail:      "IPPROTO_TCP.TCP_CORK",
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
			setsockoptInt = func(fd, level, opt, value int) (err error) {
				called = true
				testutil.Diff(t, level, tt.A().level)
				testutil.Diff(t, opt, tt.A().opt)
				testutil.Diff(t, value, tt.A().value)
				return tt.C().setErr
			}

			cf := tcpCORK(tt.C().enabled)

			if tt.A().shouldNil {
				testutil.Diff(t, Controller(nil), cf)
				return
			}
			err := cf(0)
			testutil.Diff(t, true, called)
			testutil.Diff(t, tt.A().err, err, cmpopts.EquateErrors())
		})
	}
}

func TestTCPDeferAccept(t *testing.T) {
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
				level:     syscall.IPPROTO_TCP,
				opt:       syscall.TCP_DEFER_ACCEPT,
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
				level:     syscall.IPPROTO_TCP,
				opt:       syscall.TCP_DEFER_ACCEPT,
				value:     1,
				err: &er.Error{
					Package:     ErrPkg,
					Type:        ErrTypeSocket,
					Description: ErrDscSockOpts,
					Detail:      "IPPROTO_TCP.TCP_DEFER_ACCEPT",
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
			setsockoptInt = func(fd, level, opt, value int) (err error) {
				called = true
				testutil.Diff(t, level, tt.A().level)
				testutil.Diff(t, opt, tt.A().opt)
				testutil.Diff(t, value, tt.A().value)
				return tt.C().setErr
			}

			cf := tcpDeferAccept(tt.C().value)

			if tt.A().shouldNil {
				testutil.Diff(t, Controller(nil), cf)
				return
			}
			err := cf(0)
			testutil.Diff(t, true, called)
			testutil.Diff(t, tt.A().err, err, cmpopts.EquateErrors())
		})
	}
}

func TestTCPKeepCount(t *testing.T) {
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
				level:     syscall.IPPROTO_TCP,
				opt:       syscall.TCP_KEEPCNT,
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
				level:     syscall.IPPROTO_TCP,
				opt:       syscall.TCP_KEEPCNT,
				value:     1,
				err: &er.Error{
					Package:     ErrPkg,
					Type:        ErrTypeSocket,
					Description: ErrDscSockOpts,
					Detail:      "IPPROTO_TCP.TCP_KEEPCNT",
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
			setsockoptInt = func(fd, level, opt, value int) (err error) {
				called = true
				testutil.Diff(t, level, tt.A().level)
				testutil.Diff(t, opt, tt.A().opt)
				testutil.Diff(t, value, tt.A().value)
				return tt.C().setErr
			}

			cf := tcpKeepCount(tt.C().value)

			if tt.A().shouldNil {
				testutil.Diff(t, Controller(nil), cf)
				return
			}
			err := cf(0)
			testutil.Diff(t, true, called)
			testutil.Diff(t, tt.A().err, err, cmpopts.EquateErrors())
		})
	}
}

func TestTCPKeepIdle(t *testing.T) {
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
				level:     syscall.IPPROTO_TCP,
				opt:       syscall.TCP_KEEPIDLE,
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
				level:     syscall.IPPROTO_TCP,
				opt:       syscall.TCP_KEEPIDLE,
				value:     1,
				err: &er.Error{
					Package:     ErrPkg,
					Type:        ErrTypeSocket,
					Description: ErrDscSockOpts,
					Detail:      "IPPROTO_TCP.TCP_KEEPIDLE",
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
			setsockoptInt = func(fd, level, opt, value int) (err error) {
				called = true
				testutil.Diff(t, level, tt.A().level)
				testutil.Diff(t, opt, tt.A().opt)
				testutil.Diff(t, value, tt.A().value)
				return tt.C().setErr
			}

			cf := tcpKeepIdle(tt.C().value)

			if tt.A().shouldNil {
				testutil.Diff(t, Controller(nil), cf)
				return
			}
			err := cf(0)
			testutil.Diff(t, true, called)
			testutil.Diff(t, tt.A().err, err, cmpopts.EquateErrors())
		})
	}
}

func TestTCPKeepInterval(t *testing.T) {
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
				level:     syscall.IPPROTO_TCP,
				opt:       syscall.TCP_KEEPINTVL,
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
				level:     syscall.IPPROTO_TCP,
				opt:       syscall.TCP_KEEPINTVL,
				value:     1,
				err: &er.Error{
					Package:     ErrPkg,
					Type:        ErrTypeSocket,
					Description: ErrDscSockOpts,
					Detail:      "IPPROTO_TCP.TCP_KEEPINTVL",
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
			setsockoptInt = func(fd, level, opt, value int) (err error) {
				called = true
				testutil.Diff(t, level, tt.A().level)
				testutil.Diff(t, opt, tt.A().opt)
				testutil.Diff(t, value, tt.A().value)
				return tt.C().setErr
			}

			cf := tcpKeepInterval(tt.C().value)

			if tt.A().shouldNil {
				testutil.Diff(t, Controller(nil), cf)
				return
			}
			err := cf(0)
			testutil.Diff(t, true, called)
			testutil.Diff(t, tt.A().err, err, cmpopts.EquateErrors())
		})
	}
}

func TestTCPLinger2(t *testing.T) {
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
				level:     syscall.IPPROTO_TCP,
				opt:       syscall.TCP_LINGER2,
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
				level:     syscall.IPPROTO_TCP,
				opt:       syscall.TCP_LINGER2,
				value: &syscall.Linger{
					Onoff:  1,
					Linger: 1,
				},
				err: &er.Error{
					Package:     ErrPkg,
					Type:        ErrTypeSocket,
					Description: ErrDscSockOpts,
					Detail:      "IPPROTO_TCP.TCP_LINGER2",
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
			setsockoptLinger = func(fd, level, opt int, value *syscall.Linger) (err error) {
				called = true
				testutil.Diff(t, level, tt.A().level)
				testutil.Diff(t, opt, tt.A().opt)
				testutil.Diff(t, value, tt.A().value)
				return tt.C().setErr
			}

			cf := tcpLinger2(tt.C().value)

			if tt.A().shouldNil {
				testutil.Diff(t, Controller(nil), cf)
				return
			}
			err := cf(0)
			testutil.Diff(t, true, called)
			testutil.Diff(t, tt.A().err, err, cmpopts.EquateErrors())
		})
	}
}

func TestTCPMaxSegment(t *testing.T) {
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
				level:     syscall.IPPROTO_TCP,
				opt:       syscall.TCP_MAXSEG,
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
				level:     syscall.IPPROTO_TCP,
				opt:       syscall.TCP_MAXSEG,
				value:     1,
				err: &er.Error{
					Package:     ErrPkg,
					Type:        ErrTypeSocket,
					Description: ErrDscSockOpts,
					Detail:      "IPPROTO_TCP.TCP_MAXSEG",
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
			setsockoptInt = func(fd, level, opt, value int) (err error) {
				called = true
				testutil.Diff(t, level, tt.A().level)
				testutil.Diff(t, opt, tt.A().opt)
				testutil.Diff(t, value, tt.A().value)
				return tt.C().setErr
			}

			cf := tcpMaxSegment(tt.C().value)

			if tt.A().shouldNil {
				testutil.Diff(t, Controller(nil), cf)
				return
			}
			err := cf(0)
			testutil.Diff(t, true, called)
			testutil.Diff(t, tt.A().err, err, cmpopts.EquateErrors())
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
					Detail:      "IPPROTO_TCP.TCP_NODELAY",
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
			setsockoptInt = func(fd, level, opt, value int) (err error) {
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
			testutil.Diff(t, true, called)
			testutil.Diff(t, tt.A().err, err, cmpopts.EquateErrors())
		})
	}
}

func TestTCPQuickAck(t *testing.T) {
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
				opt:       TCP_QUICKACK,
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
				opt:       TCP_QUICKACK,
				value:     1,
				err: &er.Error{
					Package:     ErrPkg,
					Type:        ErrTypeSocket,
					Description: ErrDscSockOpts,
					Detail:      "IPPROTO_TCP.TCP_QUICKACK",
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
			setsockoptInt = func(fd, level, opt, value int) (err error) {
				called = true
				testutil.Diff(t, level, tt.A().level)
				testutil.Diff(t, opt, tt.A().opt)
				testutil.Diff(t, value, tt.A().value)
				return tt.C().setErr
			}

			cf := tcpQuickAck(tt.C().enabled)

			if tt.A().shouldNil {
				testutil.Diff(t, Controller(nil), cf)
				return
			}
			err := cf(0)
			testutil.Diff(t, true, called)
			testutil.Diff(t, tt.A().err, err, cmpopts.EquateErrors())
		})
	}
}

func TestTCPSynCount(t *testing.T) {
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
				level:     syscall.IPPROTO_TCP,
				opt:       TCP_SYNCNT,
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
				level:     syscall.IPPROTO_TCP,
				opt:       TCP_SYNCNT,
				value:     1,
				err: &er.Error{
					Package:     ErrPkg,
					Type:        ErrTypeSocket,
					Description: ErrDscSockOpts,
					Detail:      "IPPROTO_TCP.TCP_SYNCNT",
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
			setsockoptInt = func(fd, level, opt, value int) (err error) {
				called = true
				testutil.Diff(t, level, tt.A().level)
				testutil.Diff(t, opt, tt.A().opt)
				testutil.Diff(t, value, tt.A().value)
				return tt.C().setErr
			}

			cf := tcpSynCount(tt.C().value)

			if tt.A().shouldNil {
				testutil.Diff(t, Controller(nil), cf)
				return
			}
			err := cf(0)
			testutil.Diff(t, true, called)
			testutil.Diff(t, tt.A().err, err, cmpopts.EquateErrors())
		})
	}
}

func TestTCPUserTimeout(t *testing.T) {
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
				level:     syscall.IPPROTO_TCP,
				opt:       TCP_USER_TIMEOUT,
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
				level:     syscall.IPPROTO_TCP,
				opt:       TCP_USER_TIMEOUT,
				value:     1,
				err: &er.Error{
					Package:     ErrPkg,
					Type:        ErrTypeSocket,
					Description: ErrDscSockOpts,
					Detail:      "IPPROTO_TCP.TCP_USER_TIMEOUT",
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
			setsockoptInt = func(fd, level, opt, value int) (err error) {
				called = true
				testutil.Diff(t, level, tt.A().level)
				testutil.Diff(t, opt, tt.A().opt)
				testutil.Diff(t, value, tt.A().value)
				return tt.C().setErr
			}

			cf := tcpUserTimeout(tt.C().value)

			if tt.A().shouldNil {
				testutil.Diff(t, Controller(nil), cf)
				return
			}
			err := cf(0)
			testutil.Diff(t, true, called)
			testutil.Diff(t, tt.A().err, err, cmpopts.EquateErrors())
		})
	}
}

func TestTCPWindowClamp(t *testing.T) {
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
				level:     syscall.IPPROTO_TCP,
				opt:       syscall.TCP_WINDOW_CLAMP,
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
				level:     syscall.IPPROTO_TCP,
				opt:       syscall.TCP_WINDOW_CLAMP,
				value:     1,
				err: &er.Error{
					Package:     ErrPkg,
					Type:        ErrTypeSocket,
					Description: ErrDscSockOpts,
					Detail:      "IPPROTO_TCP.TCP_WINDOW_CLAMP",
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
			setsockoptInt = func(fd, level, opt, value int) (err error) {
				called = true
				testutil.Diff(t, level, tt.A().level)
				testutil.Diff(t, opt, tt.A().opt)
				testutil.Diff(t, value, tt.A().value)
				return tt.C().setErr
			}

			cf := tcpWindowClamp(tt.C().value)

			if tt.A().shouldNil {
				testutil.Diff(t, Controller(nil), cf)
				return
			}
			err := cf(0)
			testutil.Diff(t, true, called)
			testutil.Diff(t, tt.A().err, err, cmpopts.EquateErrors())
		})
	}
}

func TestTCPFastOpen(t *testing.T) {
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
				opt:       TCP_FASTOPEN,
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
				opt:       TCP_FASTOPEN,
				value:     1,
				err: &er.Error{
					Package:     ErrPkg,
					Type:        ErrTypeSocket,
					Description: ErrDscSockOpts,
					Detail:      "IPPROTO_TCP.TCP_FASTOPEN",
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
			setsockoptInt = func(fd, level, opt, value int) (err error) {
				called = true
				testutil.Diff(t, level, tt.A().level)
				testutil.Diff(t, opt, tt.A().opt)
				testutil.Diff(t, value, tt.A().value)
				return tt.C().setErr
			}

			cf := tcpFastOpen(tt.C().enabled)

			if tt.A().shouldNil {
				testutil.Diff(t, Controller(nil), cf)
				return
			}
			err := cf(0)
			testutil.Diff(t, true, called)
			testutil.Diff(t, tt.A().err, err, cmpopts.EquateErrors())
		})
	}
}

func TestTCPFastOpenConnect(t *testing.T) {
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
				opt:       TCP_FASTOPEN_CONNECT,
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
				opt:       TCP_FASTOPEN_CONNECT,
				value:     1,
				err: &er.Error{
					Package:     ErrPkg,
					Type:        ErrTypeSocket,
					Description: ErrDscSockOpts,
					Detail:      "IPPROTO_TCP.TCP_FASTOPEN_CONNECT",
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
			setsockoptInt = func(fd, level, opt, value int) (err error) {
				called = true
				testutil.Diff(t, level, tt.A().level)
				testutil.Diff(t, opt, tt.A().opt)
				testutil.Diff(t, value, tt.A().value)
				return tt.C().setErr
			}

			cf := tcpFastOpenConnect(tt.C().enabled)

			if tt.A().shouldNil {
				testutil.Diff(t, Controller(nil), cf)
				return
			}
			err := cf(0)
			testutil.Diff(t, true, called)
			testutil.Diff(t, tt.A().err, err, cmpopts.EquateErrors())
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
				values: map[int]any{
					UDP_CORK:    1,
					UDP_SEGMENT: 10,
					UDP_GRO:     1,
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

func TestUDPCORK(t *testing.T) {
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
				level:     syscall.IPPROTO_UDP,
				opt:       UDP_CORK,
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
				level:     syscall.IPPROTO_UDP,
				opt:       UDP_CORK,
				value:     1,
				err: &er.Error{
					Package:     ErrPkg,
					Type:        ErrTypeSocket,
					Description: ErrDscSockOpts,
					Detail:      "IPPROTO_UDP.UDP_CORK",
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
			setsockoptInt = func(fd, level, opt, value int) (err error) {
				called = true
				testutil.Diff(t, level, tt.A().level)
				testutil.Diff(t, opt, tt.A().opt)
				testutil.Diff(t, value, tt.A().value)
				return tt.C().setErr
			}

			cf := udpCORK(tt.C().enabled)

			if tt.A().shouldNil {
				testutil.Diff(t, Controller(nil), cf)
				return
			}
			err := cf(0)
			testutil.Diff(t, true, called)
			testutil.Diff(t, tt.A().err, err, cmpopts.EquateErrors())
		})
	}
}

func TestUDPSegment(t *testing.T) {
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
				level:     syscall.IPPROTO_UDP,
				opt:       UDP_SEGMENT,
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
				level:     syscall.IPPROTO_UDP,
				opt:       UDP_SEGMENT,
				value:     1,
				err: &er.Error{
					Package:     ErrPkg,
					Type:        ErrTypeSocket,
					Description: ErrDscSockOpts,
					Detail:      "IPPROTO_UDP.UDP_SEGMENT",
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
			setsockoptInt = func(fd, level, opt, value int) (err error) {
				called = true
				testutil.Diff(t, level, tt.A().level)
				testutil.Diff(t, opt, tt.A().opt)
				testutil.Diff(t, value, tt.A().value)
				return tt.C().setErr
			}

			cf := udpSegment(tt.C().value)

			if tt.A().shouldNil {
				testutil.Diff(t, Controller(nil), cf)
				return
			}
			err := cf(0)
			testutil.Diff(t, true, called)
			testutil.Diff(t, tt.A().err, err, cmpopts.EquateErrors())
		})
	}
}

func TestUDPGRO(t *testing.T) {
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
				level:     syscall.IPPROTO_UDP,
				opt:       UDP_GRO,
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
				level:     syscall.IPPROTO_UDP,
				opt:       UDP_GRO,
				value:     1,
				err: &er.Error{
					Package:     ErrPkg,
					Type:        ErrTypeSocket,
					Description: ErrDscSockOpts,
					Detail:      "IPPROTO_UDP.UDP_GRO",
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
			setsockoptInt = func(fd, level, opt, value int) (err error) {
				called = true
				testutil.Diff(t, level, tt.A().level)
				testutil.Diff(t, opt, tt.A().opt)
				testutil.Diff(t, value, tt.A().value)
				return tt.C().setErr
			}

			cf := udpGRO(tt.C().enabled)

			if tt.A().shouldNil {
				testutil.Diff(t, Controller(nil), cf)
				return
			}
			err := cf(0)
			testutil.Diff(t, true, called)
			testutil.Diff(t, tt.A().err, err, cmpopts.EquateErrors())
		})
	}
}
