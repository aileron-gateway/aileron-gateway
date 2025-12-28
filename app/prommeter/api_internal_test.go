// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package prommeter

import (
	"regexp"
	"sync"
	"sync/atomic"
	"testing"

	v1 "github.com/aileron-gateway/aileron-gateway/apis/app/v1"
	k "github.com/aileron-gateway/aileron-gateway/apis/kernel"
	"github.com/aileron-gateway/aileron-gateway/internal/testutil"
	"github.com/aileron-gateway/aileron-gateway/kernel/api"
	utilhttp "github.com/aileron-gateway/aileron-gateway/util/http"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func TestMutate(t *testing.T) {
	type condition struct {
		manifest protoreflect.ProtoMessage
	}

	type action struct {
		manifest protoreflect.ProtoMessage
	}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"default manifest",
			&condition{
				manifest: Resource.Default(),
			},
			&action{
				manifest: &v1.PrometheusMeter{
					APIVersion: apiVersion,
					Kind:       kind,
					Metadata: &k.Metadata{
						Namespace: "default",
						Name:      "default",
					},
					Spec: &v1.PrometheusMeterSpec{
						Patterns: []string{"/metrics"},
					},
				},
			},
		),
		gen(
			"not mutated",
			&condition{
				manifest: &v1.PrometheusMeter{
					APIVersion: apiVersion,
					Kind:       kind,
					Metadata: &k.Metadata{
						Name: "default",
					},
					Spec: &v1.PrometheusMeterSpec{
						Patterns: []string{"/test"},
					},
				},
			},
			&action{
				manifest: &v1.PrometheusMeter{
					APIVersion: apiVersion,
					Kind:       kind,
					Metadata: &k.Metadata{
						Name: "default",
					},
					Spec: &v1.PrometheusMeterSpec{
						Patterns: []string{"/test"},
					},
				},
			},
		),
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			msg := Resource.Mutate(tt.C.manifest)

			opts := []cmp.Option{
				cmpopts.IgnoreUnexported(v1.PrometheusMeter{}, v1.PrometheusMeterSpec{}),
				cmpopts.IgnoreUnexported(k.Metadata{}, k.Status{}),
			}
			testutil.Diff(t, tt.A.manifest, msg, opts...)
		})
	}
}

func TestCreate(t *testing.T) {
	type condition struct {
		manifest protoreflect.ProtoMessage
	}

	type action struct {
		metrics    *metrics
		err        any // error or errorutil.Kind
		errPattern *regexp.Regexp
	}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"default",
			&condition{
				manifest: Resource.Default(),
			},
			&action{
				metrics: &metrics{
					HandlerBase: &utilhttp.HandlerBase{},
					reg:         prometheus.NewRegistry(),
					mAPICalls: prometheus.NewCounterVec(
						prometheus.CounterOpts{
							Name: "http_requests_total",
							Help: "Total number of received http requests",
						},
						[]string{"host", "path", "code", "method"},
					),
					tAPICalls: prometheus.NewCounterVec(
						prometheus.CounterOpts{
							Name: "http_client_requests_total",
							Help: "Total number of sent http requests",
						},
						[]string{"host", "path", "code", "method"},
					),
				},
				err: nil,
			},
		),
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			server := api.NewContainerAPI()
			a := &API{}
			resp, err := a.Create(server, tt.C.manifest)
			testutil.DiffError(t, tt.A.err, tt.A.errPattern, err)
			if err != nil {
				return
			}

			opts := []cmp.Option{
				cmp.AllowUnexported(utilhttp.DefaultErrorHandler{}),
				cmp.AllowUnexported(sync.RWMutex{}, atomic.Bool{}),
				cmpopts.IgnoreUnexported(prometheus.Registry{}, prometheus.MetricVec{}),
				cmp.AllowUnexported(metrics{}),
				cmpopts.IgnoreFields(metrics{}, "Handler"),
			}
			testutil.Diff(t, tt.A.metrics, resp.(*metrics), opts...)
		})
	}
}
