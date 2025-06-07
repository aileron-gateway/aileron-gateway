// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package throttle

import (
	"regexp"
	"testing"
	"time"

	v1 "github.com/aileron-gateway/aileron-gateway/apis/app/v1"
	corev1 "github.com/aileron-gateway/aileron-gateway/apis/core/v1"
	k "github.com/aileron-gateway/aileron-gateway/apis/kernel"
	"github.com/aileron-gateway/aileron-gateway/core"
	"github.com/aileron-gateway/aileron-gateway/kernel/api"
	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
	utilhttp "github.com/aileron-gateway/aileron-gateway/util/http"
	"github.com/aileron-gateway/aileron-gateway/util/resilience"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/testing/protocmp"
)

func TestMutate(t *testing.T) {
	type condition struct {
		manifest protoreflect.ProtoMessage
	}

	type action struct {
		manifest protoreflect.ProtoMessage
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"apply default values",
			[]string{}, []string{},
			&condition{
				manifest: Resource.Default(),
			},
			&action{
				manifest: &v1.ThrottleMiddleware{
					APIVersion: apiVersion,
					Kind:       kind,
					Metadata: &k.Metadata{
						Namespace: "default",
						Name:      "default",
					},
					Spec: &v1.ThrottleMiddlewareSpec{},
				},
			},
		),
		gen(
			"mutate MaxConnections",
			[]string{}, []string{},
			&condition{
				manifest: &v1.ThrottleMiddleware{
					Spec: &v1.ThrottleMiddlewareSpec{
						APIThrottlers: []*v1.APIThrottlerSpec{
							{
								Throttlers: &v1.APIThrottlerSpec_MaxConnections{
									MaxConnections: &v1.MaxConnectionsSpec{},
								},
								Waiter: &corev1.WaiterSpec{
									Waiter: &corev1.WaiterSpec_ExponentialBackoffFullJitter{
										ExponentialBackoffFullJitter: &corev1.ExponentialBackoffFullJitterWaiterSpec{},
									},
								},
							},
						},
					},
				},
			},
			&action{
				manifest: &v1.ThrottleMiddleware{
					Spec: &v1.ThrottleMiddlewareSpec{
						APIThrottlers: []*v1.APIThrottlerSpec{
							{
								Throttlers: &v1.APIThrottlerSpec_MaxConnections{
									MaxConnections: &v1.MaxConnectionsSpec{
										MaxConns: 128,
									},
								},
								Waiter: &corev1.WaiterSpec{
									Waiter: &corev1.WaiterSpec_ExponentialBackoffFullJitter{
										ExponentialBackoffFullJitter: &corev1.ExponentialBackoffFullJitterWaiterSpec{
											Base: 2000,
											Min:  0,
											Max:  1 << 21,
										},
									},
								},
							},
						},
					},
				},
			},
		),
		gen(
			"mutate FixedWindow",
			[]string{}, []string{},
			&condition{
				manifest: &v1.ThrottleMiddleware{
					Spec: &v1.ThrottleMiddlewareSpec{
						APIThrottlers: []*v1.APIThrottlerSpec{
							{
								Throttlers: &v1.APIThrottlerSpec_FixedWindow{
									FixedWindow: &v1.FixedWindowSpec{},
								},
								Waiter: &corev1.WaiterSpec{
									Waiter: &corev1.WaiterSpec_ExponentialBackoffFullJitter{
										ExponentialBackoffFullJitter: &corev1.ExponentialBackoffFullJitterWaiterSpec{},
									},
								},
							},
						},
					},
				},
			},
			&action{
				manifest: &v1.ThrottleMiddleware{
					Spec: &v1.ThrottleMiddlewareSpec{
						APIThrottlers: []*v1.APIThrottlerSpec{
							{
								Throttlers: &v1.APIThrottlerSpec_FixedWindow{
									FixedWindow: &v1.FixedWindowSpec{
										WindowSize: 1000,
										Limit:      1000,
									},
								},
								Waiter: &corev1.WaiterSpec{
									Waiter: &corev1.WaiterSpec_ExponentialBackoffFullJitter{
										ExponentialBackoffFullJitter: &corev1.ExponentialBackoffFullJitterWaiterSpec{
											Base: 2000,
											Min:  0,
											Max:  1 << 21,
										},
									},
								},
							},
						},
					},
				},
			},
		),
		gen(
			"mutate TokenBucket",
			[]string{}, []string{},
			&condition{
				manifest: &v1.ThrottleMiddleware{
					Spec: &v1.ThrottleMiddlewareSpec{
						APIThrottlers: []*v1.APIThrottlerSpec{
							{
								Throttlers: &v1.APIThrottlerSpec_TokenBucket{
									TokenBucket: &v1.TokenBucketSpec{},
								},
								Waiter: &corev1.WaiterSpec{
									Waiter: &corev1.WaiterSpec_ExponentialBackoffFullJitter{
										ExponentialBackoffFullJitter: &corev1.ExponentialBackoffFullJitterWaiterSpec{},
									},
								},
							},
						},
					},
				},
			},
			&action{
				manifest: &v1.ThrottleMiddleware{
					Spec: &v1.ThrottleMiddlewareSpec{
						APIThrottlers: []*v1.APIThrottlerSpec{
							{
								Throttlers: &v1.APIThrottlerSpec_TokenBucket{
									TokenBucket: &v1.TokenBucketSpec{
										BucketSize:   1000,
										FillInterval: 1000,
										FillRate:     1000,
									},
								},
								Waiter: &corev1.WaiterSpec{
									Waiter: &corev1.WaiterSpec_ExponentialBackoffFullJitter{
										ExponentialBackoffFullJitter: &corev1.ExponentialBackoffFullJitterWaiterSpec{
											Base: 2000,
											Min:  0,
											Max:  1 << 21,
										},
									},
								},
							},
						},
					},
				},
			},
		),
		gen(
			"mutate LeakyBucket",
			[]string{}, []string{},
			&condition{
				manifest: &v1.ThrottleMiddleware{
					Spec: &v1.ThrottleMiddlewareSpec{
						APIThrottlers: []*v1.APIThrottlerSpec{
							{
								Throttlers: &v1.APIThrottlerSpec_LeakyBucket{
									LeakyBucket: &v1.LeakyBucketSpec{},
								},
								Waiter: &corev1.WaiterSpec{
									Waiter: &corev1.WaiterSpec_ExponentialBackoffFullJitter{
										ExponentialBackoffFullJitter: &corev1.ExponentialBackoffFullJitterWaiterSpec{},
									},
								},
							},
						},
					},
				},
			},
			&action{
				manifest: &v1.ThrottleMiddleware{
					Spec: &v1.ThrottleMiddlewareSpec{
						APIThrottlers: []*v1.APIThrottlerSpec{
							{
								Throttlers: &v1.APIThrottlerSpec_LeakyBucket{
									LeakyBucket: &v1.LeakyBucketSpec{
										BucketSize:   1000,
										LeakInterval: 1000,
										LeakRate:     200,
									},
								},
								Waiter: &corev1.WaiterSpec{
									Waiter: &corev1.WaiterSpec_ExponentialBackoffFullJitter{
										ExponentialBackoffFullJitter: &corev1.ExponentialBackoffFullJitterWaiterSpec{
											Base: 2000,
											Min:  0,
											Max:  1 << 21,
										},
									},
								},
							},
						},
					},
				},
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			msg := Resource.Mutate(tt.C().manifest)

			opts := []cmp.Option{
				protocmp.Transform(),
				cmpopts.IgnoreUnexported(v1.ThrottleMiddleware{}, v1.ThrottleMiddlewareSpec{}, v1.APIThrottlerSpec{}),
				cmpopts.IgnoreUnexported(k.Metadata{}, k.Status{}),
			}
			testutil.Diff(t, tt.A().manifest, msg, opts...)
		})
	}
}

func TestCreate(t *testing.T) {
	type condition struct {
		manifest protoreflect.ProtoMessage
	}

	type action struct {
		err        any // error or errorutil.Kind
		errPattern *regexp.Regexp
		expect     any
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"create with default manifest",
			[]string{}, []string{},
			&condition{
				manifest: Resource.Default(),
			},
			&action{
				err: nil,
				expect: &throttle{
					eh:         utilhttp.GlobalErrorHandler(utilhttp.DefaultErrorHandlerName),
					throttlers: []*apiThrottler{},
				},
			},
		),
		gen(
			"nil matcher",
			[]string{}, []string{},
			&condition{
				manifest: &v1.ThrottleMiddleware{
					Metadata: &k.Metadata{},
					Spec: &v1.ThrottleMiddlewareSpec{
						APIThrottlers: []*v1.APIThrottlerSpec{
							{Matcher: nil},
						},
					},
				},
			},
			&action{
				err: nil,
				expect: &throttle{
					eh:         utilhttp.GlobalErrorHandler(utilhttp.DefaultErrorHandlerName),
					throttlers: []*apiThrottler{},
				},
			},
		),
		gen(
			"fail to create APIThrottlers",
			[]string{}, []string{},
			&condition{
				manifest: &v1.ThrottleMiddleware{
					Metadata: &k.Metadata{},
					Spec: &v1.ThrottleMiddlewareSpec{
						APIThrottlers: []*v1.APIThrottlerSpec{
							{
								Matcher: &k.MatcherSpec{
									Patterns: []string{
										"[0-9",
									},
									MatchType: k.MatchType_Regex,
								},
							},
						},
					},
				},
			},
			&action{
				err:        core.ErrCoreGenCreateObject,
				errPattern: regexp.MustCompile(core.ErrPrefix + `failed to create ThrottleMiddleware`),
			},
		),
		gen(
			"create with MaxConnections",
			[]string{}, []string{},
			&condition{
				manifest: &v1.ThrottleMiddleware{
					Metadata: &k.Metadata{},
					Spec: &v1.ThrottleMiddlewareSpec{
						APIThrottlers: []*v1.APIThrottlerSpec{
							{
								Throttlers: &v1.APIThrottlerSpec_MaxConnections{
									MaxConnections: &v1.MaxConnectionsSpec{
										MaxConns: 128,
									},
								},
								Waiter: &corev1.WaiterSpec{
									Waiter: &corev1.WaiterSpec_ExponentialBackoffFullJitter{
										ExponentialBackoffFullJitter: &corev1.ExponentialBackoffFullJitterWaiterSpec{},
									},
								},
								Matcher: &k.MatcherSpec{
									Patterns: []string{"/example"},
								},
							},
						},
					},
				},
			},
			&action{
				err: nil,
				expect: &throttle{
					eh: utilhttp.GlobalErrorHandler(utilhttp.DefaultErrorHandlerName),
					throttlers: []*apiThrottler{
						{
							paths: nil, // This fields will be un checked.
							throttler: &maxConnections{
								sem: make(chan struct{}, 128),
							},
						},
					},
				},
			},
		),
		gen(
			"create with FixedWindow",
			[]string{}, []string{},
			&condition{
				manifest: &v1.ThrottleMiddleware{
					Metadata: &k.Metadata{},
					Spec: &v1.ThrottleMiddlewareSpec{
						APIThrottlers: []*v1.APIThrottlerSpec{
							{
								Throttlers: &v1.APIThrottlerSpec_FixedWindow{
									FixedWindow: &v1.FixedWindowSpec{
										WindowSize: 1000,
										Limit:      1000,
									},
								},
								Waiter: &corev1.WaiterSpec{
									Waiter: &corev1.WaiterSpec_ExponentialBackoffFullJitter{
										ExponentialBackoffFullJitter: &corev1.ExponentialBackoffFullJitterWaiterSpec{},
									},
								},
								Matcher: &k.MatcherSpec{
									Patterns: []string{"/example"},
								},
							},
						},
					},
				},
			},
			&action{
				err: nil,
				expect: &throttle{
					eh: utilhttp.GlobalErrorHandler(utilhttp.DefaultErrorHandlerName),
					throttlers: []*apiThrottler{
						{
							paths: nil, // This fields will be un checked.
							throttler: &fixedWindow{
								bucket: make(chan struct{}, 1000),
								window: 1 * time.Second,
							},
						},
					},
				},
			},
		),
		gen(
			"create with TokenBucket",
			[]string{}, []string{},
			&condition{
				manifest: &v1.ThrottleMiddleware{
					Metadata: &k.Metadata{},
					Spec: &v1.ThrottleMiddlewareSpec{
						APIThrottlers: []*v1.APIThrottlerSpec{
							{
								Throttlers: &v1.APIThrottlerSpec_TokenBucket{
									TokenBucket: &v1.TokenBucketSpec{
										BucketSize:   1000,
										FillInterval: 1000,
										FillRate:     1000,
									},
								},
								Waiter: &corev1.WaiterSpec{
									Waiter: &corev1.WaiterSpec_ExponentialBackoffFullJitter{
										ExponentialBackoffFullJitter: &corev1.ExponentialBackoffFullJitterWaiterSpec{},
									},
								},
								Matcher: &k.MatcherSpec{
									Patterns: []string{"/example"},
								},
							},
						},
					},
				},
			},
			&action{
				err: nil,
				expect: &throttle{
					eh: utilhttp.GlobalErrorHandler(utilhttp.DefaultErrorHandlerName),
					throttlers: []*apiThrottler{
						{
							paths: nil, // This fields will be un checked.
							throttler: &tokenBucket{
								bucket:   make(chan struct{}, 1000),
								rate:     1,
								interval: 1000 * time.Second,
							},
						},
					},
				},
			},
		),
		gen(
			"create with LeakyBucket",
			[]string{}, []string{},
			&condition{
				manifest: &v1.ThrottleMiddleware{
					Metadata: &k.Metadata{},
					Spec: &v1.ThrottleMiddlewareSpec{
						APIThrottlers: []*v1.APIThrottlerSpec{
							{
								Throttlers: &v1.APIThrottlerSpec_LeakyBucket{
									LeakyBucket: &v1.LeakyBucketSpec{
										BucketSize:   1000,
										LeakInterval: 1000,
										LeakRate:     200,
									},
								},
								Waiter: &corev1.WaiterSpec{
									Waiter: &corev1.WaiterSpec_ExponentialBackoffFullJitter{
										ExponentialBackoffFullJitter: &corev1.ExponentialBackoffFullJitterWaiterSpec{},
									},
								},
								Matcher: &k.MatcherSpec{
									Patterns: []string{"/example"},
								},
							},
						},
					},
				},
			},
			&action{
				err: nil,
				expect: &throttle{
					eh: utilhttp.GlobalErrorHandler(utilhttp.DefaultErrorHandlerName),
					throttlers: []*apiThrottler{
						{
							paths: nil, // This fields will be un checked.
							throttler: &leakyBucket{
								bucket:   make(chan chan struct{}, 1000),
								rate:     200,
								interval: 1 * time.Second,
							},
						},
					},
				},
			},
		),
		gen(
			"create with retryThrottler",
			[]string{}, []string{},
			&condition{
				manifest: &v1.ThrottleMiddleware{
					Metadata: &k.Metadata{},
					Spec: &v1.ThrottleMiddlewareSpec{
						APIThrottlers: []*v1.APIThrottlerSpec{
							{
								Throttlers: &v1.APIThrottlerSpec_LeakyBucket{
									LeakyBucket: &v1.LeakyBucketSpec{
										BucketSize:   1000,
										LeakInterval: 1000,
										LeakRate:     200,
									},
								},
								Waiter: &corev1.WaiterSpec{
									Waiter: &corev1.WaiterSpec_ExponentialBackoffFullJitter{
										ExponentialBackoffFullJitter: &corev1.ExponentialBackoffFullJitterWaiterSpec{},
									},
								},
								Matcher: &k.MatcherSpec{
									Patterns: []string{"/example"},
								},
								MaxRetry: 3,
							},
						},
					},
				},
			},
			&action{
				err: nil,
				expect: &throttle{
					eh: utilhttp.GlobalErrorHandler(utilhttp.DefaultErrorHandlerName),
					throttlers: []*apiThrottler{
						{
							paths: nil, // This fields will be un checked.
							throttler: &retryThrottler{
								throttler: &leakyBucket{
									bucket:   make(chan chan struct{}, 1000),
									rate:     200,
									interval: 1 * time.Second,
								},
								maxRetry: 3,
								waiter: resilience.NewWaiter(&corev1.WaiterSpec{
									Waiter: &corev1.WaiterSpec_ExponentialBackoffFullJitter{
										ExponentialBackoffFullJitter: &corev1.ExponentialBackoffFullJitterWaiterSpec{},
									},
								}),
							},
						},
					},
				},
			},
		),
	}
	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			server := api.NewContainerAPI()
			a := &API{}
			got, err := a.Create(server, tt.C().manifest)

			opts := []cmp.Option{
				protocmp.Transform(),
				cmp.AllowUnexported(throttle{}),
				cmp.AllowUnexported(maxConnections{}),
				// Compare channel capacity for errors with different memory addresses
				cmp.Comparer(func(x, y maxConnections) bool {
					return cap(x.sem) == cap(y.sem)
				}),
				cmp.AllowUnexported(fixedWindow{}),
				//Compare channel capacity for errors with different memory addresses
				cmp.Comparer(func(x, y fixedWindow) bool {
					if cap(x.bucket) != cap(y.bucket) {
						return false
					}
					if x.window != y.window {
						return false
					}
					return true
				}),
				cmp.AllowUnexported(tokenBucket{}),
				//Compare channel capacity for errors with different memory addresses
				cmp.Comparer(func(x, y tokenBucket) bool {
					return cap(x.bucket) == cap(y.bucket)
				}),
				cmp.AllowUnexported(leakyBucket{}),
				//Compare channel capacity for errors with different memory addresses
				cmp.Comparer(func(x, y leakyBucket) bool {
					if cap(x.bucket) != cap(y.bucket) {
						return false
					}
					if x.interval != y.interval {
						return false
					}
					if x.rate != y.rate {
						return false
					}
					return true
				}),
				cmp.AllowUnexported(retryThrottler{}),
				cmp.Comparer(func(x, y retryThrottler) bool {
					return x.maxRetry == y.maxRetry
				}),
				cmp.AllowUnexported(apiThrottler{}),
				cmpopts.IgnoreFields(apiThrottler{}, "paths"),
				cmp.Comparer(testutil.ComparePointer[core.ErrorHandler]),
			}

			testutil.DiffError(t, tt.A().err, tt.A().errPattern, err)
			testutil.Diff(t, tt.A().expect, got, opts...)
		})
	}
}
