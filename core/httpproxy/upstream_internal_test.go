package httpproxy

import (
	"errors"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"testing"
	"time"

	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
	"github.com/aileron-gateway/aileron-gateway/util/resilience"
)

var (
	_ resilience.Entry = &noopUpstream{}
	_ resilience.Entry = &lbUpstream{}
)

type testCircuitBreaker struct {
	activeStatus bool
	successCount int
	failureCount int
}

func (cb *testCircuitBreaker) Active() bool  { return cb.activeStatus }
func (cb *testCircuitBreaker) countSuccess() { cb.successCount += 1 }
func (cb *testCircuitBreaker) countFailure() { cb.failureCount += 1 }

func TestLBUpstream_url(t *testing.T) {
	type condition struct {
		upstream *lbUpstream
	}

	type action struct {
		rawURL    string
		parsedURL *url.URL
	}

	actCheckURL := "check url string"

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	tb.Action(actCheckURL, "check the returned string")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"single upstream",
			[]string{},
			[]string{actCheckURL},
			&condition{
				upstream: &lbUpstream{
					rawURL: "http://test.com",
					parsedURL: &url.URL{
						Scheme: "http",
						Host:   "test.com",
					},
				},
			},
			&action{
				rawURL: "http://test.com",
				parsedURL: &url.URL{
					Scheme: "http",
					Host:   "test.com",
				},
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			rawURL := tt.C().upstream.ID()
			testutil.Diff(t, tt.A().rawURL, rawURL)

			parsedURL := tt.C().upstream.url()
			testutil.Diff(t, tt.A().parsedURL, parsedURL)
		})
	}
}

func TestLBUpstream_notify(t *testing.T) {
	type condition struct {
		upstream *lbUpstream
		status   int
		err      error
	}

	type action struct {
		success int
		failure int
	}

	cndPassiveDisabled := "passive disabled"
	cndCountSuccess := "input success condition"
	cndCountFailure := "input failure condition"
	actCheckSuccessCount := "check the count of success"
	actCheckFailureCount := "check the count of failure"

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	tb.Condition(cndPassiveDisabled, "disable passive health check")
	tb.Condition(cndCountSuccess, "input success condition")
	tb.Condition(cndCountFailure, "input failure condition")
	tb.Action(actCheckSuccessCount, "check the count of the success")
	tb.Action(actCheckFailureCount, "check the count of failure")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"passive disabled",
			[]string{cndPassiveDisabled},
			[]string{actCheckSuccessCount, actCheckFailureCount},
			&condition{
				upstream: &lbUpstream{
					circuitBreaker: &testCircuitBreaker{},
					passiveEnabled: false,
				},
			},
			&action{
				success: 0,
				failure: 0,
			},
		),
		gen(
			"count success with status 0",
			[]string{cndCountSuccess},
			[]string{actCheckSuccessCount, actCheckFailureCount},
			&condition{
				upstream: &lbUpstream{
					circuitBreaker: &testCircuitBreaker{},
					passiveEnabled: true,
				},
				status: 0,
			},
			&action{
				success: 1,
				failure: 0,
			},
		),
		gen(
			"count success with status 200",
			[]string{cndCountSuccess},
			[]string{actCheckSuccessCount, actCheckFailureCount},
			&condition{
				upstream: &lbUpstream{
					circuitBreaker: &testCircuitBreaker{},
					passiveEnabled: true,
				},
				status: 200,
			},
			&action{
				success: 1,
				failure: 0,
			},
		),
		gen(
			"count success with status 499",
			[]string{cndCountSuccess},
			[]string{actCheckSuccessCount, actCheckFailureCount},
			&condition{
				upstream: &lbUpstream{
					circuitBreaker: &testCircuitBreaker{},
					passiveEnabled: true,
				},
				status: 499,
			},
			&action{
				success: 1,
				failure: 0,
			},
		),
		gen(
			"count failure with status 500",
			[]string{cndCountFailure},
			[]string{actCheckSuccessCount, actCheckFailureCount},
			&condition{
				upstream: &lbUpstream{
					circuitBreaker: &testCircuitBreaker{},
					passiveEnabled: true,
				},
				status: 500,
			},
			&action{
				success: 0,
				failure: 1,
			},
		),
		gen(
			"count failure with an error",
			[]string{cndCountFailure},
			[]string{actCheckSuccessCount, actCheckFailureCount},
			&condition{
				upstream: &lbUpstream{
					circuitBreaker: &testCircuitBreaker{},
					passiveEnabled: true,
				},
				err: errors.New("test"),
			},
			&action{
				success: 0,
				failure: 1,
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			tt.C().upstream.notify(tt.C().status, tt.C().err)

			cb := tt.C().upstream.circuitBreaker.(*testCircuitBreaker)
			testutil.Diff(t, tt.A().success, cb.successCount)
			testutil.Diff(t, tt.A().failure, cb.failureCount)
		})
	}
}

func TestLBUpstream_activeCheckHTTP(t *testing.T) {
	type condition struct {
		upstream *lbUpstream
		rt       http.RoundTripper
	}

	type action struct {
		success int
		failure int
	}

	cndCountSuccess := "input success condition"
	cndCountFailure := "input failure condition"
	actCheckSuccessCount := "check the count of success"
	actCheckFailureCount := "check the count of failure"

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	tb.Condition(cndCountSuccess, "input success condition")
	tb.Condition(cndCountFailure, "input failure condition")
	tb.Action(actCheckSuccessCount, "check the count of the success")
	tb.Action(actCheckFailureCount, "check the count of failure")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"200 OK response",
			[]string{cndCountSuccess},
			[]string{actCheckSuccessCount, actCheckFailureCount},
			&condition{
				upstream: &lbUpstream{
					circuitBreaker: &testCircuitBreaker{},
					interval:       time.Millisecond,
				},
				rt: &testRoundTripper{
					status: http.StatusOK,
				},
			},
			&action{
				success: 1,
				failure: 0,
			},
		),
		gen(
			"200 OK response with 0 interval",
			[]string{cndCountSuccess},
			[]string{actCheckSuccessCount, actCheckFailureCount},
			&condition{
				upstream: &lbUpstream{
					circuitBreaker: &testCircuitBreaker{},
					interval:       0,
				},
				rt: &testRoundTripper{
					status: http.StatusOK,
				},
			},
			&action{
				success: 1,
				failure: 0,
			},
		),
		gen(
			"server error",
			[]string{cndCountFailure},
			[]string{actCheckSuccessCount, actCheckFailureCount},
			&condition{
				upstream: &lbUpstream{
					circuitBreaker: &testCircuitBreaker{},
					interval:       time.Millisecond,
				},
				rt: &testRoundTripper{
					status: http.StatusInternalServerError,
				},
			},
			&action{
				success: 0,
				failure: 1,
			},
		),
		gen(
			"error on round trip",
			[]string{cndCountFailure},
			[]string{actCheckSuccessCount, actCheckFailureCount},
			&condition{
				upstream: &lbUpstream{
					circuitBreaker: &testCircuitBreaker{},
					interval:       time.Millisecond,
				},
				rt: &testRoundTripper{
					err: http.ErrHandlerTimeout,
				},
			},
			&action{
				success: 0,
				failure: 1,
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			go func() {
				time.Sleep(50 * time.Millisecond)
				tt.C().upstream.close()
			}()

			req, _ := http.NewRequest(http.MethodGet, "http://test.com/example", nil)
			tt.C().upstream.activeCheckHTTP(tt.C().rt, req)

			cb := tt.C().upstream.circuitBreaker.(*testCircuitBreaker)
			testutil.Diff(t, true, tt.A().success <= cb.successCount)
			testutil.Diff(t, true, tt.A().failure <= cb.failureCount)
		})
	}
}

func TestLBUpstream_activeCheck(t *testing.T) {
	type condition struct {
		upstream *lbUpstream
		nw       string
		addr     string
	}

	type action struct {
		success int
		failure int
	}

	cndTCP := "tcp,tcp4,tcp6"
	cndUDP := "udp,udp4,udp6"
	cndIP := "ip,ip4,ip6"
	cndCountSuccess := "input success condition"
	cndCountFailure := "input failure condition"
	actCheckSuccessCount := "check the count of success"
	actCheckFailureCount := "check the count of failure"

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	tb.Condition(cndTCP, "specify tcp,tcp4,tcp6 as network type")
	tb.Condition(cndUDP, "specify udp,udp4,udp6 as network type")
	tb.Condition(cndIP, "specify ip,ip4,ip6 as network type")
	tb.Condition(cndCountSuccess, "input success condition")
	tb.Condition(cndCountFailure, "input failure condition")
	tb.Action(actCheckSuccessCount, "check the count of the success")
	tb.Action(actCheckFailureCount, "check the count of failure")
	table := tb.Build()

	// Listen with IPv4 address.
	ln4, err := net.Listen("tcp4", "127.0.0.1:0")
	if err != nil {
		t.Errorf("%#v\n", err)
		return
	}
	defer ln4.Close()
	testPort4 := strconv.Itoa(ln4.Addr().(*net.TCPAddr).Port)

	// Listen with IPv6 address.
	ln6, err := net.Listen("tcp6", "[::1]:0")
	if err != nil {
		t.Errorf("%#v\n", err)
		return
	}
	defer ln6.Close()
	testPort6 := strconv.Itoa(ln6.Addr().(*net.TCPAddr).Port)

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"tcp success with 0 interval",
			[]string{cndTCP, cndCountSuccess},
			[]string{actCheckSuccessCount, actCheckFailureCount},
			&condition{
				upstream: &lbUpstream{
					circuitBreaker: &testCircuitBreaker{},
					interval:       0,
				},
				nw:   "tcp",
				addr: "127.0.0.1:" + testPort4,
			},
			&action{
				success: 1,
				failure: 0,
			},
		),
		gen(
			"tcp success as localhost",
			[]string{cndTCP, cndCountSuccess},
			[]string{actCheckSuccessCount, actCheckFailureCount},
			&condition{
				upstream: &lbUpstream{
					circuitBreaker: &testCircuitBreaker{},
					interval:       time.Millisecond,
				},
				nw:   "tcp",
				addr: "localhost:" + testPort4,
			},
			&action{
				success: 1,
				failure: 0,
			},
		),
		gen(
			"tcp",
			[]string{cndTCP, cndCountSuccess},
			[]string{actCheckSuccessCount, actCheckFailureCount},
			&condition{
				upstream: &lbUpstream{
					circuitBreaker: &testCircuitBreaker{},
					interval:       time.Millisecond,
				},
				nw:   "tcp",
				addr: "127.0.0.1:" + testPort4,
			},
			&action{
				success: 1,
				failure: 0,
			},
		),
		// gen(
		// 	"tcp4",
		// 	[]string{cndTCP, cndCountSuccess},
		// 	[]string{actCheckSuccessCount, actCheckFailureCount},
		// 	&condition{
		// 		upstream: &lbUpstream{
		// 			circuitBreaker: &testCircuitBreaker{},
		// 			interval:       time.Millisecond,
		// 		},
		// 		nw:   "tcp4",
		// 		addr: "127.0.0.1:" + testPort4,
		// 	},
		// 	&action{
		// 		success: 1,
		// 		failure: 0,
		// 	},
		// ),
		gen(
			"tcp6",
			[]string{cndTCP, cndCountSuccess},
			[]string{actCheckSuccessCount, actCheckFailureCount},
			&condition{
				upstream: &lbUpstream{
					circuitBreaker: &testCircuitBreaker{},
					interval:       time.Millisecond,
				},
				nw:   "tcp6",
				addr: "[::1]:" + testPort6,
			},
			&action{
				success: 1,
				failure: 0,
			},
		),
		gen(
			"udp success",
			[]string{cndUDP, cndCountSuccess},
			[]string{actCheckSuccessCount, actCheckFailureCount},
			&condition{
				upstream: &lbUpstream{
					circuitBreaker: &testCircuitBreaker{},
					interval:       time.Millisecond,
				},
				nw:   "udp",
				addr: "127.0.0.1:" + testPort4,
			},
			&action{
				success: 1,
				failure: 0,
			},
		),
		gen(
			"udp4 success",
			[]string{cndUDP, cndCountSuccess},
			[]string{actCheckSuccessCount, actCheckFailureCount},
			&condition{
				upstream: &lbUpstream{
					circuitBreaker: &testCircuitBreaker{},
					interval:       time.Millisecond,
				},
				nw:   "udp4",
				addr: "127.0.0.1:" + testPort4,
			},
			&action{
				success: 1,
				failure: 0,
			},
		),
		gen(
			"udp6 success",
			[]string{cndUDP, cndCountSuccess},
			[]string{actCheckSuccessCount, actCheckFailureCount},
			&condition{
				upstream: &lbUpstream{
					circuitBreaker: &testCircuitBreaker{},
					interval:       time.Millisecond,
				},
				nw:   "udp6",
				addr: "[::1]:" + testPort6,
			},
			&action{
				success: 1,
				failure: 0,
			},
		),
		// Following tests fail in some environment.
		// gen(
		// 	"ip success",
		// 	[]string{cndIP, cndCountSuccess},
		// 	[]string{actCheckSuccessCount, actCheckFailureCount},
		// 	&condition{
		// 		upstream: &lbUpstream{
		// 			circuitBreaker: &testCircuitBreaker{},
		// 			interval:       time.Millisecond,
		// 		},
		// 		nw:   "ip:1",
		// 		addr: "127.0.0.1",
		// 	},
		// 	&action{
		// 		success: 1,
		// 		failure: 0,
		// 	},
		// ),
		// gen(
		// 	"ip4 success",
		// 	[]string{cndIP, cndCountSuccess},
		// 	[]string{actCheckSuccessCount, actCheckFailureCount},
		// 	&condition{
		// 		upstream: &lbUpstream{
		// 			circuitBreaker: &testCircuitBreaker{},
		// 			interval:       time.Millisecond,
		// 		},
		// 		nw:   "ip4:icmp",
		// 		addr: "127.0.0.1",
		// 	},
		// 	&action{
		// 		success: 1,
		// 		failure: 0,
		// 	},
		// ),
		// gen(
		// 	"ip6 success",
		// 	[]string{cndIP, cndCountSuccess},
		// 	[]string{actCheckSuccessCount, actCheckFailureCount},
		// 	&condition{
		// 		upstream: &lbUpstream{
		// 			circuitBreaker: &testCircuitBreaker{},
		// 			interval:       time.Millisecond,
		// 		},
		// 		nw:   "ip6:icmp",
		// 		addr: "[::1]",
		// 	},
		// 	&action{
		// 		success: 1,
		// 		failure: 0,
		// 	},
		// ),
		// gen(
		// 	"ip:icmp success",
		// 	[]string{cndIP, cndCountSuccess},
		// 	[]string{actCheckSuccessCount, actCheckFailureCount},
		// 	&condition{
		// 		upstream: &lbUpstream{
		// 			circuitBreaker: &testCircuitBreaker{},
		// 			interval:       time.Millisecond,
		// 		},
		// 		nw:   "ip:icmp",
		// 		addr: "127.0.0.1",
		// 	},
		// 	&action{
		// 		success: 1,
		// 		failure: 0,
		// 	},
		// ),
		// gen(
		// 	"ip4:icmp success",
		// 	[]string{cndIP, cndCountSuccess},
		// 	[]string{actCheckSuccessCount, actCheckFailureCount},
		// 	&condition{
		// 		upstream: &lbUpstream{
		// 			circuitBreaker: &testCircuitBreaker{},
		// 			interval:       time.Millisecond,
		// 		},
		// 		nw:   "ip4:icmp",
		// 		addr: "127.0.0.1",
		// 	},
		// 	&action{
		// 		success: 1,
		// 		failure: 0,
		// 	},
		// ),
		// gen(
		// 	"ip6:icmp success",
		// 	[]string{cndIP, cndCountSuccess},
		// 	[]string{actCheckSuccessCount, actCheckFailureCount},
		// 	&condition{
		// 		upstream: &lbUpstream{
		// 			circuitBreaker: &testCircuitBreaker{},
		// 			interval:       time.Millisecond,
		// 		},
		// 		nw:   "ip6:icmp",
		// 		addr: "[::1]",
		// 	},
		// 	&action{
		// 		success: 1,
		// 		failure: 0,
		// 	},
		// ),
		gen(
			"invalid network ipv4",
			[]string{cndCountFailure},
			[]string{actCheckSuccessCount, actCheckFailureCount},
			&condition{
				upstream: &lbUpstream{
					circuitBreaker: &testCircuitBreaker{},
					interval:       time.Millisecond,
				},
				nw:   "INVALID",
				addr: "127.0.0.1:" + testPort4,
			},
			&action{
				success: 0,
				failure: 1,
			},
		),
		gen(
			"invalid network ipv6",
			[]string{cndCountFailure},
			[]string{actCheckSuccessCount, actCheckFailureCount},
			&condition{
				upstream: &lbUpstream{
					circuitBreaker: &testCircuitBreaker{},
					interval:       time.Millisecond,
				},
				nw:   "INVALID",
				addr: "[::1]:" + testPort6,
			},
			&action{
				success: 0,
				failure: 1,
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			go func() {
				time.Sleep(50 * time.Millisecond)
				tt.C().upstream.close()
			}()

			tt.C().upstream.activeCheck(tt.C().nw, tt.C().addr)

			cb := tt.C().upstream.circuitBreaker.(*testCircuitBreaker)
			testutil.Diff(t, true, tt.A().success <= cb.successCount)
			testutil.Diff(t, true, tt.A().failure <= cb.failureCount)
		})
	}
}
