// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package redis

import (
	"crypto/tls"
	"time"

	v1 "github.com/aileron-gateway/aileron-gateway/apis/app/v1"
	"github.com/aileron-gateway/aileron-gateway/apis/kernel"
	"github.com/aileron-gateway/aileron-gateway/core"
	"github.com/aileron-gateway/aileron-gateway/kernel/api"
	"github.com/aileron-gateway/aileron-gateway/kernel/network"
	"github.com/redis/go-redis/v9"
	"google.golang.org/protobuf/reflect/protoreflect"
)

const (
	apiVersion = "app/v1"
	kind       = "RedisClient"
	Key        = apiVersion + "/" + kind
)

var Resource api.Resource = &API{
	BaseResource: &api.BaseResource{
		DefaultProto: &v1.RedisClient{
			APIVersion: apiVersion,
			Kind:       kind,
			Metadata: &kernel.Metadata{
				Namespace: "default",
				Name:      "default",
			},
			Spec: &v1.RedisClientSpec{
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
}

type API struct {
	*api.BaseResource
}

// Mutate changes configured values.
// The values of the msg which is given as the argument is the merged message of default values and user defined values.
// Changes for the fields of msg in this function make the final values which will be the input for validate and create function.
// Default values for "repeated" or "oneof" fields can also be applied in this function if necessary.
// Please check msg!=nil and asserting the mgs does not panic even they won't from the view of overall architecture of the gateway.
func (*API) Mutate(msg protoreflect.ProtoMessage) protoreflect.ProtoMessage {
	c := msg.(*v1.RedisClient)

	if len(c.Spec.Addrs) == 0 {
		c.Spec.Addrs = append(c.Spec.Addrs, "localhost:6379")
	}

	return c
}

func (*API) Create(a api.API[*api.Request, *api.Response], msg protoreflect.ProtoMessage) (any, error) {
	c := msg.(*v1.RedisClient)

	var tls *tls.Config
	if c.Spec.TLSConfig != nil {
		tlsConfig, err := network.TLSConfig(c.Spec.TLSConfig)
		if err != nil {
			return nil, core.ErrCoreGenCreateObject.WithStack(err, map[string]any{"kind": kind})
		}
		tls = tlsConfig
	}

	opts := &redis.UniversalOptions{
		Addrs: c.Spec.Addrs,
		DB:    int(c.Spec.DB),

		Username:         c.Spec.Username,
		Password:         c.Spec.Password,
		SentinelUsername: c.Spec.SentinelUsername,
		SentinelPassword: c.Spec.SentinelPassword,

		MaxRetries:      int(c.Spec.MaxRetries),
		MinRetryBackoff: time.Millisecond * time.Duration(c.Spec.MinRetryBackoff),
		MaxRetryBackoff: time.Millisecond * time.Duration(c.Spec.MaxRetryBackoff),

		DialTimeout:           time.Millisecond * time.Duration(c.Spec.DialTimeout),
		ReadTimeout:           time.Millisecond * time.Duration(c.Spec.ReadTimeout),
		WriteTimeout:          time.Millisecond * time.Duration(c.Spec.WriteTimeout),
		ContextTimeoutEnabled: c.Spec.ContextTimeoutEnabled,

		// PoolFIFO uses FIFO mode for each node connection pool GET/PUT (default LIFO).
		PoolFIFO: c.Spec.PoolFIFO,

		PoolSize:        int(c.Spec.PoolSize),
		PoolTimeout:     time.Millisecond * time.Duration(c.Spec.PoolTimeout),
		MinIdleConns:    int(c.Spec.MinIdleConns),
		MaxIdleConns:    int(c.Spec.MaxIdleConns),
		ConnMaxIdleTime: time.Millisecond * time.Duration(c.Spec.ConnMaxIdleTime),
		ConnMaxLifetime: time.Millisecond * time.Duration(c.Spec.ConnMaxLifetime),

		TLSConfig: tls,

		MaxRedirects:   int(c.Spec.MaxRedirects),
		ReadOnly:       c.Spec.ReadOnly,
		RouteByLatency: c.Spec.RouteByLatency,
		RouteRandomly:  c.Spec.RouteRandomly,

		MasterName: c.Spec.MasterName,
	}

	return &client{
		UniversalClient: redis.NewUniversalClient(opts),
		timeout:         time.Millisecond * time.Duration(c.Spec.Timeout),
		expiration:      time.Millisecond * time.Duration(c.Spec.Expiration),
	}, nil
}
