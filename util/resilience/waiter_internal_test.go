package resilience

import (
	"testing"
	"time"

	v1 "github.com/aileron-gateway/aileron-gateway/apis/core/v1"
	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestDefaultWaiterSpec(t *testing.T) {
	type condition struct {
		spec *v1.WaiterSpec
	}

	type action struct {
		waiter any
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"nil",
			[]string{},
			[]string{},
			&condition{
				spec: nil,
			},
			&action{
				waiter: nil,
			},
		),
		gen(
			"nil fixed backoff",
			[]string{},
			[]string{},
			&condition{
				spec: &v1.WaiterSpec{
					Waiter: &v1.WaiterSpec_FixedBackoff{},
				},
			},
			&action{
				waiter: &v1.WaiterSpec_FixedBackoff{
					FixedBackoff: &v1.FixedBackoffWaiterSpec{
						Base: 5_000,
					},
				},
			},
		),
		gen(
			"nil linear backoff",
			[]string{},
			[]string{},
			&condition{
				spec: &v1.WaiterSpec{
					Waiter: &v1.WaiterSpec_LinearBackoff{},
				},
			},
			&action{
				waiter: &v1.WaiterSpec_LinearBackoff{
					LinearBackoff: &v1.LinearBackoffWaiterSpec{
						Base: 5_000,
						Max:  1 << 21,
					},
				},
			},
		),
		gen(
			"nil polynomial backoff",
			[]string{},
			[]string{},
			&condition{
				spec: &v1.WaiterSpec{
					Waiter: &v1.WaiterSpec_PolynomialBackoff{},
				},
			},
			&action{
				waiter: &v1.WaiterSpec_PolynomialBackoff{
					PolynomialBackoff: &v1.PolynomialBackoffWaiterSpec{
						Base:     2_000,
						Exponent: 2,
						Max:      1 << 21,
					},
				},
			},
		),
		gen(
			"nil exponential backoff",
			[]string{},
			[]string{},
			&condition{
				spec: &v1.WaiterSpec{
					Waiter: &v1.WaiterSpec_ExponentialBackoff{},
				},
			},
			&action{
				waiter: &v1.WaiterSpec_ExponentialBackoff{
					ExponentialBackoff: &v1.ExponentialBackoffWaiterSpec{
						Base: 2_000,
						Max:  1 << 21,
					},
				},
			},
		),
		gen(
			"nil exponential full jitter backoff",
			[]string{},
			[]string{},
			&condition{
				spec: &v1.WaiterSpec{
					Waiter: &v1.WaiterSpec_ExponentialBackoffFullJitter{},
				},
			},
			&action{
				waiter: &v1.WaiterSpec_ExponentialBackoffFullJitter{
					ExponentialBackoffFullJitter: &v1.ExponentialBackoffFullJitterWaiterSpec{
						Base: 2_000,
						Max:  1 << 21,
					},
				},
			},
		),
		gen(
			"nil exponential equal jitter backoff",
			[]string{},
			[]string{},
			&condition{
				spec: &v1.WaiterSpec{
					Waiter: &v1.WaiterSpec_ExponentialBackoffEqualJitter{},
				},
			},
			&action{
				waiter: &v1.WaiterSpec_ExponentialBackoffEqualJitter{
					ExponentialBackoffEqualJitter: &v1.ExponentialBackoffEqualJitterWaiterSpec{
						Base: 2_000,
						Max:  1 << 21,
					},
				},
			},
		),
		gen(
			"fixed backoff",
			[]string{},
			[]string{},
			&condition{
				spec: &v1.WaiterSpec{
					Waiter: &v1.WaiterSpec_FixedBackoff{
						FixedBackoff: &v1.FixedBackoffWaiterSpec{
							Base: 1_000, // not default value.
						},
					},
				},
			},
			&action{
				waiter: &v1.WaiterSpec_FixedBackoff{
					FixedBackoff: &v1.FixedBackoffWaiterSpec{
						Base: 1_000,
					},
				},
			},
		),
		gen(
			"linear backoff",
			[]string{},
			[]string{},
			&condition{
				spec: &v1.WaiterSpec{
					Waiter: &v1.WaiterSpec_LinearBackoff{
						LinearBackoff: &v1.LinearBackoffWaiterSpec{
							Base: 1_000, // not default value.
							Max:  1 << 21,
						},
					},
				},
			},
			&action{
				waiter: &v1.WaiterSpec_LinearBackoff{
					LinearBackoff: &v1.LinearBackoffWaiterSpec{
						Base: 1_000,
						Max:  1 << 21,
					},
				},
			},
		),
		gen(
			"polynomial backoff",
			[]string{},
			[]string{},
			&condition{
				spec: &v1.WaiterSpec{
					Waiter: &v1.WaiterSpec_PolynomialBackoff{
						PolynomialBackoff: &v1.PolynomialBackoffWaiterSpec{
							Base:     1_000, // not default value.
							Exponent: 3,     // not default value.
							Max:      1 << 21,
						},
					},
				},
			},
			&action{
				waiter: &v1.WaiterSpec_PolynomialBackoff{
					PolynomialBackoff: &v1.PolynomialBackoffWaiterSpec{
						Base:     1_000,
						Exponent: 3,
						Max:      1 << 21,
					},
				},
			},
		),
		gen(
			"exponential backoff",
			[]string{},
			[]string{},
			&condition{
				spec: &v1.WaiterSpec{
					Waiter: &v1.WaiterSpec_ExponentialBackoff{
						ExponentialBackoff: &v1.ExponentialBackoffWaiterSpec{
							Base: 1_000, // not default value.
							Max:  1 << 21,
						},
					},
				},
			},
			&action{
				waiter: &v1.WaiterSpec_ExponentialBackoff{
					ExponentialBackoff: &v1.ExponentialBackoffWaiterSpec{
						Base: 1_000,
						Max:  1 << 21,
					},
				},
			},
		),
		gen(
			"exponential full jitter backoff",
			[]string{},
			[]string{},
			&condition{
				spec: &v1.WaiterSpec{
					Waiter: &v1.WaiterSpec_ExponentialBackoffFullJitter{
						ExponentialBackoffFullJitter: &v1.ExponentialBackoffFullJitterWaiterSpec{
							Base: 1_000, // not default value.
							Max:  1 << 21,
						},
					},
				},
			},
			&action{
				waiter: &v1.WaiterSpec_ExponentialBackoffFullJitter{
					ExponentialBackoffFullJitter: &v1.ExponentialBackoffFullJitterWaiterSpec{
						Base: 1_000,
						Max:  1 << 21,
					},
				},
			},
		),
		gen(
			"exponential equal jitter backoff",
			[]string{},
			[]string{},
			&condition{
				spec: &v1.WaiterSpec{
					Waiter: &v1.WaiterSpec_ExponentialBackoffEqualJitter{
						ExponentialBackoffEqualJitter: &v1.ExponentialBackoffEqualJitterWaiterSpec{
							Base: 1_000, // not default value.
							Max:  1 << 21,
						},
					},
				},
			},
			&action{
				waiter: &v1.WaiterSpec_ExponentialBackoffEqualJitter{
					ExponentialBackoffEqualJitter: &v1.ExponentialBackoffEqualJitterWaiterSpec{
						Base: 1_000,
						Max:  1 << 21,
					},
				},
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			DefaultWaiterSpec(tt.C().spec)

			if tt.A().waiter == nil {
				testutil.Diff(t, (*v1.WaiterSpec)(nil), tt.C().spec)
				return
			}

			opts := []cmp.Option{
				cmpopts.IgnoreUnexported(v1.FixedBackoffWaiterSpec{}),
				cmpopts.IgnoreUnexported(v1.LinearBackoffWaiterSpec{}),
				cmpopts.IgnoreUnexported(v1.PolynomialBackoffWaiterSpec{}),
				cmpopts.IgnoreUnexported(v1.ExponentialBackoffWaiterSpec{}),
				cmpopts.IgnoreUnexported(v1.ExponentialBackoffFullJitterWaiterSpec{}),
				cmpopts.IgnoreUnexported(v1.ExponentialBackoffEqualJitterWaiterSpec{}),
			}
			testutil.Diff(t, tt.A().waiter, tt.C().spec.Waiter, opts...)
		})
	}
}

func TestNewWaiter(t *testing.T) {
	type condition struct {
		spec *v1.WaiterSpec
	}

	type action struct {
		waiter Waiter
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"fixed backoff",
			[]string{},
			[]string{},
			&condition{
				spec: &v1.WaiterSpec{
					Waiter: &v1.WaiterSpec_FixedBackoff{
						FixedBackoff: &v1.FixedBackoffWaiterSpec{
							Base: 1_000,
						},
					},
				},
			},
			&action{
				waiter: &fixedBackoffWaiter{
					base: 1_000,
				},
			},
		),
		gen(
			"linear backoff",
			[]string{},
			[]string{},
			&condition{
				spec: &v1.WaiterSpec{
					Waiter: &v1.WaiterSpec_LinearBackoff{
						LinearBackoff: &v1.LinearBackoffWaiterSpec{
							Base: 1_000,
							Min:  100,
							Max:  1 << 21,
						},
					},
				},
			},
			&action{
				waiter: &linearBackoffWaiter{
					base: 1000,
					min:  100,
					max:  1 << 21,
				},
			},
		),
		gen(
			"polynomial backoff",
			[]string{},
			[]string{},
			&condition{
				spec: &v1.WaiterSpec{
					Waiter: &v1.WaiterSpec_PolynomialBackoff{
						PolynomialBackoff: &v1.PolynomialBackoffWaiterSpec{
							Base:     1_000,
							Exponent: 3,
							Min:      100,
							Max:      1 << 21,
						},
					},
				},
			},
			&action{
				waiter: &polynomialBackoffWaiter{
					base:     1_000,
					min:      100,
					max:      1 << 21,
					exponent: 3,
				},
			},
		),
		gen(
			"exponential backoff",
			[]string{},
			[]string{},
			&condition{
				spec: &v1.WaiterSpec{
					Waiter: &v1.WaiterSpec_ExponentialBackoff{
						ExponentialBackoff: &v1.ExponentialBackoffWaiterSpec{
							Base: 1_000,
							Min:  100,
							Max:  1 << 21,
						},
					},
				},
			},
			&action{
				waiter: &exponentialBackoffWaiter{
					base: 1_000,
					min:  100,
					max:  1 << 21,
				},
			},
		),
		gen(
			"exponential full jitter backoff",
			[]string{},
			[]string{},
			&condition{
				spec: &v1.WaiterSpec{
					Waiter: &v1.WaiterSpec_ExponentialBackoffFullJitter{
						ExponentialBackoffFullJitter: &v1.ExponentialBackoffFullJitterWaiterSpec{
							Base: 1_000,
							Min:  100,
							Max:  1 << 21,
						},
					},
				},
			},
			&action{
				waiter: &exponentialBackoffFullJitterWaiter{
					base: 1_000,
					min:  100,
					max:  1 << 21,
				},
			},
		),
		gen(
			"exponential equal jitter backoff",
			[]string{},
			[]string{},
			&condition{
				spec: &v1.WaiterSpec{
					Waiter: &v1.WaiterSpec_ExponentialBackoffEqualJitter{
						ExponentialBackoffEqualJitter: &v1.ExponentialBackoffEqualJitterWaiterSpec{
							Base: 1_000,
							Min:  100,
							Max:  1 << 21,
						},
					},
				},
			},
			&action{
				waiter: &exponentialBackoffEqualJitterWaiter{
					base: 1_000,
					min:  100,
					max:  1 << 21,
				},
			},
		),
		gen(
			"default exponential full jitter backoff",
			[]string{},
			[]string{},
			&condition{
				spec: &v1.WaiterSpec{
					Waiter: nil,
				},
			},
			&action{
				waiter: &exponentialBackoffFullJitterWaiter{
					base: 2_000,
					min:  0,
					max:  1 << 21,
				},
			},
		),
		gen(
			"nil exponential full jitter backoff",
			[]string{},
			[]string{},
			&condition{
				spec: nil,
			},
			&action{
				waiter: &exponentialBackoffFullJitterWaiter{
					base: 2_000,
					min:  0,
					max:  1 << 21,
				},
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			w := NewWaiter(tt.C().spec)
			opts := []cmp.Option{
				cmp.AllowUnexported(fixedBackoffWaiter{}),
				cmp.AllowUnexported(linearBackoffWaiter{}),
				cmp.AllowUnexported(polynomialBackoffWaiter{}),
				cmp.AllowUnexported(exponentialBackoffWaiter{}),
				cmp.AllowUnexported(exponentialBackoffFullJitterWaiter{}),
				cmp.AllowUnexported(exponentialBackoffEqualJitterWaiter{}),
			}
			testutil.Diff(t, tt.A().waiter, w, opts...)
		})
	}
}

func TestFixedBackoffWaiter_Wait(t *testing.T) {
	type condition struct {
		w Waiter
		n int
	}

	type action struct {
		d time.Duration
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"-1",
			[]string{},
			[]string{},
			&condition{
				w: &fixedBackoffWaiter{
					base: 10,
				},
				n: -1,
			},
			&action{
				d: time.Microsecond * 10,
			},
		),
		gen(
			"0",
			[]string{},
			[]string{},
			&condition{
				w: &fixedBackoffWaiter{
					base: 10,
				},
				n: 0,
			},
			&action{
				d: time.Microsecond * 10,
			},
		),
		gen(
			"1",
			[]string{},
			[]string{},
			&condition{
				w: &fixedBackoffWaiter{
					base: 10,
				},
				n: 1,
			},
			&action{
				d: time.Microsecond * 10,
			},
		),
		gen(
			"2",
			[]string{},
			[]string{},
			&condition{
				w: &fixedBackoffWaiter{
					base: 10,
				},
				n: 2,
			},
			&action{
				d: time.Microsecond * 10,
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			d := tt.C().w.Wait(tt.C().n)
			testutil.Diff(t, tt.A().d+time.Nanosecond, d)
		})
	}
}

func TestLinearBackoffWaiter_Wait(t *testing.T) {
	type condition struct {
		w Waiter
		n int
	}

	type action struct {
		d time.Duration
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"-1",
			[]string{},
			[]string{},
			&condition{
				w: &linearBackoffWaiter{
					base: 10,
					min:  30,
					max:  100,
				},
				n: -1,
			},
			&action{
				d: time.Microsecond * 30,
			},
		),
		gen(
			"0",
			[]string{},
			[]string{},
			&condition{
				w: &linearBackoffWaiter{
					base: 10,
					min:  30,
					max:  100,
				},
				n: 0,
			},
			&action{
				d: time.Microsecond * 30,
			},
		),
		gen(
			"1",
			[]string{},
			[]string{},
			&condition{
				w: &linearBackoffWaiter{
					base: 10,
					min:  30,
					max:  100,
				},
				n: 1,
			},
			&action{
				d: time.Microsecond * 30,
			},
		),
		gen(
			"2",
			[]string{},
			[]string{},
			&condition{
				w: &linearBackoffWaiter{
					base: 10,
					min:  30,
					max:  100,
				},
				n: 2,
			},
			&action{
				d: time.Microsecond * 30,
			},
		),
		gen(
			"20",
			[]string{},
			[]string{},
			&condition{
				w: &linearBackoffWaiter{
					base: 10,
					min:  30,
					max:  100,
				},
				n: 20,
			},
			&action{
				d: time.Microsecond * 100,
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			d := tt.C().w.Wait(tt.C().n)
			testutil.Diff(t, tt.A().d+time.Nanosecond, d)
		})
	}
}

func TestPolynomialBackoffWaiter_Wait(t *testing.T) {
	type condition struct {
		w Waiter
		n int
	}

	type action struct {
		d time.Duration
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"-1",
			[]string{},
			[]string{},
			&condition{
				w: &polynomialBackoffWaiter{
					exponent: 5,
					base:     10,
					min:      20,
					max:      100,
				},
				n: -1,
			},
			&action{
				d: time.Microsecond * 20,
			},
		),
		gen(
			"0",
			[]string{},
			[]string{},
			&condition{
				w: &polynomialBackoffWaiter{
					exponent: 5,
					base:     10,
					min:      20,
					max:      100,
				},
				n: 0,
			},
			&action{
				d: time.Microsecond * 20,
			},
		),
		gen(
			"1",
			[]string{},
			[]string{},
			&condition{
				w: &polynomialBackoffWaiter{
					exponent: 5,
					base:     10,
					min:      20,
					max:      100,
				},
				n: 1,
			},
			&action{
				d: time.Microsecond * 20,
			},
		),
		gen(
			"2",
			[]string{},
			[]string{},
			&condition{
				w: &polynomialBackoffWaiter{
					exponent: 5,
					base:     10,
					min:      20,
					max:      100,
				},
				n: 2,
			},
			&action{
				d: time.Microsecond * 100,
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			d := tt.C().w.Wait(tt.C().n)
			testutil.Diff(t, tt.A().d+time.Nanosecond, d)
		})
	}
}

func TestExponentialBackoffWaiter_Wait(t *testing.T) {
	type condition struct {
		w Waiter
		n int
	}

	type action struct {
		d time.Duration
	}
	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"-1",
			[]string{},
			[]string{},
			&condition{
				w: &exponentialBackoffWaiter{
					base: 10,
					min:  30,
					max:  100,
				},
				n: -1,
			},
			&action{
				d: time.Microsecond * 30,
			},
		),
		gen(
			"0",
			[]string{},
			[]string{},
			&condition{
				w: &exponentialBackoffWaiter{
					base: 10,
					min:  30,
					max:  100,
				},
				n: 0,
			},
			&action{
				d: time.Microsecond * 30,
			},
		),
		gen(
			"1",
			[]string{},
			[]string{},
			&condition{
				w: &exponentialBackoffWaiter{
					base: 10,
					min:  30,
					max:  100,
				},
				n: 1,
			},
			&action{
				d: time.Microsecond * 30,
			},
		),
		gen(
			"2",
			[]string{},
			[]string{},
			&condition{
				w: &exponentialBackoffWaiter{
					base: 10,
					min:  30,
					max:  100,
				},
				n: 2,
			},
			&action{
				d: time.Microsecond * 40,
			},
		),
		gen(
			"10",
			[]string{},
			[]string{},
			&condition{
				w: &exponentialBackoffWaiter{
					base: 10,
					min:  30,
					max:  100,
				},
				n: 10,
			},
			&action{
				d: time.Microsecond * 100,
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			d := tt.C().w.Wait(tt.C().n)
			testutil.Diff(t, tt.A().d+time.Nanosecond, d)
		})
	}
}

func TestExponentialBackoffFullJitterWaiter_Wait(t *testing.T) {
	type condition struct {
		w   Waiter
		min int
		max int
	}

	type action struct {
		min time.Duration
		max time.Duration
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"random",
			[]string{},
			[]string{},
			&condition{
				w: &exponentialBackoffFullJitterWaiter{
					base: 10,
					min:  20,
					max:  100,
				},
				min: -10,
				max: 50,
			},
			&action{
				min: time.Microsecond * 20,
				max: time.Microsecond * 100,
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			for i := tt.C().min; i < tt.C().max; i++ {
				d := tt.C().w.Wait(i)
				testutil.Diff(t, true, d >= tt.A().min+time.Nanosecond)
				testutil.Diff(t, true, d <= tt.A().max+time.Nanosecond)
			}
		})
	}
}

func TestExponentialBackoffEqualJitterWaiter_Wait(t *testing.T) {
	type condition struct {
		w   Waiter
		min int
		max int
	}

	type action struct {
		min time.Duration
		max time.Duration
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"random",
			[]string{},
			[]string{},
			&condition{
				w: &exponentialBackoffEqualJitterWaiter{
					base: 10,
					min:  20,
					max:  100,
				},
				min: -10,
				max: 50,
			},
			&action{
				min: time.Microsecond * 20,
				max: time.Microsecond * 100,
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			for i := tt.C().min; i < tt.C().max; i++ {
				d := tt.C().w.Wait(i)
				testutil.Diff(t, true, d >= tt.A().min+time.Nanosecond)
				testutil.Diff(t, true, d <= tt.A().max+time.Nanosecond)
			}
		})
	}
}
