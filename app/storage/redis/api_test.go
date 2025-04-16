// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package redis

import (
	"testing"

	v1 "github.com/aileron-gateway/aileron-gateway/apis/app/v1"
	k "github.com/aileron-gateway/aileron-gateway/apis/kernel"
	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"google.golang.org/protobuf/reflect/protoreflect"
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
			[]string{},
			[]string{},
			&condition{
				manifest: Resource.Default(),
			},
			&action{
				manifest: &v1.RedisClient{
					APIVersion: apiVersion,
					Kind:       kind,
					Metadata: &k.Metadata{
						Namespace: "default",
						Name:      "default",
					},
					Spec: &v1.RedisClientSpec{
						Addrs:           []string{"localhost:6379"},
						MaxRetries:      2,
						MaxRetryBackoff: 3,
						ReadTimeout:     5000,
						WriteTimeout:    5000,
						PoolSize:        1000,
						MaxIdleConns:    100,
						MinIdleConns:    10,
						ConnMaxIdleTime: 10_000,
						ConnMaxLifetime: 600_000,
						Timeout:         10_000,
						Expiration:      0,
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
				cmpopts.IgnoreUnexported(v1.RedisClient{}, v1.RedisClientSpec{}),
				cmpopts.IgnoreUnexported(k.Metadata{}, k.Status{}),
			}
			testutil.Diff(t, tt.A().manifest, msg, opts...)
		})
	}
}
