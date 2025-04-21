// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package resilience

import (
	"math"
	"math/rand/v2"
	"time"

	v1 "github.com/aileron-gateway/aileron-gateway/apis/core/v1"
	"google.golang.org/protobuf/proto"
)

// DefaultWaiterSpec marges default value of waiters into given WaiterSpec.
// This function do nothing if nil spec was given.
func DefaultWaiterSpec(spec *v1.WaiterSpec) {
	if spec == nil || spec.Waiter == nil {
		return
	}

	switch w := spec.Waiter.(type) {
	case *v1.WaiterSpec_FixedBackoff:
		v := &v1.FixedBackoffWaiterSpec{
			Base: 5_000, // 5ms
		}
		proto.Merge(v, w.FixedBackoff)
		spec.Waiter = &v1.WaiterSpec_FixedBackoff{
			FixedBackoff: v,
		}
	case *v1.WaiterSpec_LinearBackoff:
		v := &v1.LinearBackoffWaiterSpec{
			Base: 5_000,   // 5ms
			Max:  1 << 21, // 2,097ms
		}
		proto.Merge(v, w.LinearBackoff)
		spec.Waiter = &v1.WaiterSpec_LinearBackoff{
			LinearBackoff: v,
		}
	case *v1.WaiterSpec_PolynomialBackoff:
		v := &v1.PolynomialBackoffWaiterSpec{
			Base:     2_000, // 2ms
			Exponent: 2,
			Max:      1 << 21, // 2,097ms
		}
		proto.Merge(v, w.PolynomialBackoff)
		spec.Waiter = &v1.WaiterSpec_PolynomialBackoff{
			PolynomialBackoff: v,
		}
	case *v1.WaiterSpec_ExponentialBackoff:
		v := &v1.ExponentialBackoffWaiterSpec{
			Base: 2_000,   // 2ms
			Max:  1 << 21, // 2,097ms
		}
		proto.Merge(v, w.ExponentialBackoff)
		spec.Waiter = &v1.WaiterSpec_ExponentialBackoff{
			ExponentialBackoff: v,
		}
	case *v1.WaiterSpec_ExponentialBackoffFullJitter:
		v := &v1.ExponentialBackoffFullJitterWaiterSpec{
			Base: 2_000,   // 2ms
			Max:  1 << 21, // 2,097ms
		}
		proto.Merge(v, w.ExponentialBackoffFullJitter)
		spec.Waiter = &v1.WaiterSpec_ExponentialBackoffFullJitter{
			ExponentialBackoffFullJitter: v,
		}
	case *v1.WaiterSpec_ExponentialBackoffEqualJitter:
		v := &v1.ExponentialBackoffEqualJitterWaiterSpec{
			Base: 2_000,   // 2ms
			Max:  1 << 21, // 2,097ms
		}
		proto.Merge(v, w.ExponentialBackoffEqualJitter)
		spec.Waiter = &v1.WaiterSpec_ExponentialBackoffEqualJitter{
			ExponentialBackoffEqualJitter: v,
		}
	}
}

// NewWaiter returns a new waiter.
// Values are not validated in this function.
// ExponentialBackoffFullJitter with default values are used  if nil was given.
func NewWaiter(spec *v1.WaiterSpec) Waiter {
	if spec == nil {
		spec = &v1.WaiterSpec{
			Waiter: &v1.WaiterSpec_ExponentialBackoffFullJitter{
				ExponentialBackoffFullJitter: &v1.ExponentialBackoffFullJitterWaiterSpec{
					Base: 2_000,
					Min:  0,
					Max:  1 << 21,
				},
			},
		}
	}

	switch spec.Waiter.(type) {
	case *v1.WaiterSpec_FixedBackoff:
		v := spec.GetFixedBackoff()
		return &fixedBackoffWaiter{
			base: int64(v.Base),
		}
	case *v1.WaiterSpec_LinearBackoff:
		v := spec.GetLinearBackoff()
		return &linearBackoffWaiter{
			base: int64(v.Base),
			min:  int64(v.Min),
			max:  int64(v.Max),
		}
	case *v1.WaiterSpec_PolynomialBackoff:
		v := spec.GetPolynomialBackoff()
		return &polynomialBackoffWaiter{
			exponent: int64(v.Exponent),
			base:     int64(v.Base),
			min:      int64(v.Min),
			max:      int64(v.Max),
		}
	case *v1.WaiterSpec_ExponentialBackoff:
		v := spec.GetExponentialBackoff()
		return &exponentialBackoffWaiter{
			base: int64(v.Base),
			min:  int64(v.Min),
			max:  int64(v.Max),
		}
	case *v1.WaiterSpec_ExponentialBackoffFullJitter:
		v := spec.GetExponentialBackoffFullJitter()
		return &exponentialBackoffFullJitterWaiter{
			base: int64(v.Base),
			min:  int64(v.Min),
			max:  int64(v.Max),
		}
	case *v1.WaiterSpec_ExponentialBackoffEqualJitter:
		v := spec.GetExponentialBackoffEqualJitter()
		return &exponentialBackoffEqualJitterWaiter{
			base: int64(v.Base),
			min:  int64(v.Min),
			max:  int64(v.Max),
		}
	default:
		return &exponentialBackoffFullJitterWaiter{
			base: 2_000,
			min:  0,
			max:  1 << 21,
		}
	}
}

// Waiter return a duration to wait for.
// The returned value must be grater than 0 to avoid ticker's error.
// Calculation strategy depends on the implementer.
// Wait must not block the process but return the duration immediately.
type Waiter interface {
	// Wait returns the time to wait for with time units.
	Wait(int) time.Duration
}

// fixedBackoffWaiter is a waiter with fixed backoff strategy.
// This implements http.Waiter interface.
type fixedBackoffWaiter struct {
	base int64
}

func (w *fixedBackoffWaiter) Wait(_ int) time.Duration {
	// Return with extra 1 nanosecond to avoid ticker's error.
	return time.Nanosecond + time.Microsecond*time.Duration(w.base)
}

// linearBackoffWaiter is a waiter with linear backoff strategy.
// It returns values of 'base * i' microsecond.
// Additionally, the returned value is restricted between min and max.
// This implements http.Waiter interface.
type linearBackoffWaiter struct {
	base int64
	min  int64
	max  int64
}

func (w *linearBackoffWaiter) Wait(i int) time.Duration {
	if i <= 0 {
		return time.Nanosecond + time.Microsecond*time.Duration(w.min)
	}

	t := w.base * int64(i)
	if t > w.max {
		t = w.max
	} else if t < w.min {
		t = w.min
	}

	// Return with extra 1 nanosecond to avoid ticker's error.
	return time.Nanosecond + time.Microsecond*time.Duration(t)
}

// polynomialBackoffWaiter is a waiter with polynomial backoff strategy.
// If the exponent is less than 0, it is set to 2.
// It returns values of 'base * i^exponent' microsecond.
// Additionally, the returned value is restricted between min and max.
// This implements http.Waiter interface.
type polynomialBackoffWaiter struct {
	exponent int64
	base     int64
	min      int64
	max      int64
}

func (w *polynomialBackoffWaiter) Wait(i int) time.Duration {
	if i <= 0 {
		return time.Nanosecond + time.Microsecond*time.Duration(w.min)
	}

	t := w.base * int64(math.Pow(float64(i), float64(w.exponent)))
	if t > w.max {
		t = w.max
	} else if t < w.min {
		t = w.min
	}

	// Return with extra 1 nanosecond to avoid ticker's error.
	return time.Nanosecond + time.Microsecond*time.Duration(t)
}

// exponentialBackoffWaiter is a waiter with exponential backoff strategy.
// It returns values of 'base * 2^i' microsecond.
// Additionally, the returned value is restricted between min and max.
// This implements http.Waiter interface.
type exponentialBackoffWaiter struct {
	base int64
	min  int64
	max  int64
}

func (w *exponentialBackoffWaiter) Wait(i int) time.Duration {
	if i <= 0 {
		return time.Nanosecond + time.Microsecond*time.Duration(w.min)
	}

	t := w.base * (1 << i)
	if t > w.max {
		t = w.max
	} else if t < w.min {
		t = w.min
	}
	// Return with extra 1 nanosecond to avoid ticker's error.
	return time.Nanosecond + time.Microsecond*time.Duration(t)
}

// exponentialBackoffFullJitterWaiter is a waiter with exponential backoff with full jitter strategy.
// It returns values of 'random(0, base * 2^i)' microsecond.
// Additionally, the returned value is restricted between min and max.
// This implements http.Waiter interface.
type exponentialBackoffFullJitterWaiter struct {
	base int64
	min  int64
	max  int64
}

func (w *exponentialBackoffFullJitterWaiter) Wait(i int) time.Duration {
	if i <= 0 {
		return time.Nanosecond + time.Microsecond*time.Duration(w.min)
	}

	t := w.base * (1 << i)
	if t > w.max {
		t = w.max
	}
	t = rand.Int64N(t + 1) // Add 1 to avoid 0.
	if t < w.min {
		t = w.min
	}

	// Return with extra 1 nanosecond to avoid ticker's error.
	return time.Nanosecond + time.Microsecond*time.Duration(t)
}

// exponentialBackoffEqualJitterWaiter is a waiter with exponential backoff with equal jitter strategy.
// It returns values of '(base * 2^i)/2 + random(0, (base * 2^i)/2)' microsecond.
// Additionally, the returned value is restricted between min and max.
// This implements http.Waiter interface.
type exponentialBackoffEqualJitterWaiter struct {
	base int64
	min  int64
	max  int64
}

func (w *exponentialBackoffEqualJitterWaiter) Wait(i int) time.Duration {
	if i <= 0 {
		return time.Nanosecond + time.Microsecond*time.Duration(w.min)
	}

	t := w.base * (1 << i)
	if t > w.max {
		t = w.max
	}
	t = t/2 + rand.Int64N(t/2+1) // Add 1 to avoid 0.
	if t < w.min {
		t = w.min
	}

	// Return with extra 1 nanosecond to avoid ticker's error.
	return time.Nanosecond + time.Microsecond*time.Duration(t)
}
