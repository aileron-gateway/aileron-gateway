// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package redis

import (
	"context"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/aileron-gateway/aileron-gateway/app"
	"github.com/aileron-gateway/aileron-gateway/internal/testutil"
	"github.com/aileron-gateway/aileron-gateway/kernel/kvs"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/redis/go-redis/v9"
)

// errorType is the types of error for Get operation against redis.
type errorType int

const (
	noError   errorType = iota // no error
	redisNil                   // redis.Nil
	someError                  // an error
)

// mocUniversalClient is the moc redis client.
// This implements redis.UniversalClient.
type mocUniversalClient struct {
	redis.UniversalClient
	closeError bool      // flag to return an error for Close operation.
	getError   errorType // the returned error type for Get operation.
	redisError bool      // flag to return an error for Set,SetEx,Del,Exists,Expire operation.
}

func (m *mocUniversalClient) Close() error {
	if m.closeError {
		return errors.New("test error")
	}
	return nil
}

func (m *mocUniversalClient) Get(ctx context.Context, key string) *redis.StringCmd {
	switch m.getError {
	case redisNil:
		return redis.NewStringResult("redis nil", redis.Nil)
	case someError:
		return redis.NewStringResult("test error", errors.New("test error"))
	default:
		return redis.NewStringResult("test", nil)
	}
}

func (m *mocUniversalClient) Set(ctx context.Context, key string, value any, expiration time.Duration) *redis.StatusCmd {
	if m.redisError {
		return redis.NewStatusResult("test error", errors.New("test error"))
	}
	return redis.NewStatusResult("test", nil)
}

func (m *mocUniversalClient) SetEx(ctx context.Context, key string, value any, expiration time.Duration) *redis.StatusCmd {
	if m.redisError {
		return redis.NewStatusResult("test error", errors.New("test error"))
	}
	return redis.NewStatusResult("test", nil)
}

func (m *mocUniversalClient) Del(ctx context.Context, keys ...string) *redis.IntCmd {
	if m.redisError {
		return redis.NewIntResult(0, errors.New("test error"))
	}
	return redis.NewIntResult(0, nil)
}

func (m *mocUniversalClient) Exists(ctx context.Context, keys ...string) *redis.IntCmd {
	if m.redisError {
		return redis.NewIntResult(0, redis.Nil)
	}
	return redis.NewIntResult(1, nil)
}

func (m *mocUniversalClient) Expire(ctx context.Context, key string, exp time.Duration) *redis.BoolCmd {
	if m.redisError {
		return redis.NewBoolResult(false, errors.New("test error"))
	}
	return redis.NewBoolResult(true, nil)
}

func TestOpen(t *testing.T) {
	type condition struct {
		client client
	}

	type action struct {
		err any // error or errorutil.Kind
	}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"set client",
			&condition{
				client: client{
					timeout:    5 * time.Second,
					expiration: 10 * time.Second,
				},
			},
			&action{
				err: nil,
			},
		),
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			err := tt.C.client.Open(context.Background())
			testutil.Diff(t, tt.A.err, err)
		})
	}
}

func TestClose(t *testing.T) {
	type condition struct {
		mock redis.UniversalClient
	}

	type action struct {
		err        any // error or errorutil.Kind
		errPattern *regexp.Regexp
	}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"Close func returns nil",
			&condition{
				mock: &mocUniversalClient{
					closeError: false,
				},
			},
			&action{
				err:        nil,
				errPattern: nil,
			},
		),
		gen(
			"Close func returns non nil",
			&condition{
				mock: &mocUniversalClient{
					closeError: true,
				},
			},
			&action{
				err:        app.ErrAppStorageKVS,
				errPattern: regexp.MustCompile("error occurred while operating the key-value store"),
			},
		),
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			client := client{
				UniversalClient: tt.C.mock,
				timeout:         5 * time.Second,
				expiration:      10 * time.Second,
			}
			err := client.Close(context.Background())
			testutil.DiffError(t, tt.A.err, tt.A.errPattern, err)
		})
	}
}

func TestGet(t *testing.T) {
	type condition struct {
		mock    redis.UniversalClient
		context context.Context
	}

	type action struct {
		value      []byte
		err        any // error or errorutil.Kind
		errPattern *regexp.Regexp
	}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"cmd.Err() returns redis.nil",
			&condition{
				mock: &mocUniversalClient{
					getError: redisNil,
				},
				context: context.Background(),
			},
			&action{
				value:      nil,
				err:        kvs.Nil,
				errPattern: nil,
			},
		),
		gen(
			"cmd.Err() returns non nil and non redis.nil",
			&condition{
				mock: &mocUniversalClient{
					getError: someError,
				},
				context: context.Background(),
			},
			&action{
				value:      nil,
				err:        app.ErrAppStorageKVS,
				errPattern: regexp.MustCompile("error occurred while operating the key-value store"),
			},
		),
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			client := client{
				UniversalClient: tt.C.mock,
				timeout:         5 * time.Second,
				expiration:      10 * time.Second,
			}

			value, err := client.Get(tt.C.context, "")
			testutil.DiffError(t, tt.A.err, tt.A.errPattern, err, cmpopts.EquateErrors())
			testutil.Diff(t, tt.A.value, value)
		})
	}
}

func TestSet(t *testing.T) {
	type condition struct {
		mock    redis.UniversalClient
		context context.Context
	}

	type action struct {
		err        any // error or errorutil.Kind
		errPattern *regexp.Regexp
	}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"cmd.Err() returns some error",
			&condition{
				mock: &mocUniversalClient{
					redisError: true,
				},
				context: context.Background(),
			},
			&action{
				err:        app.ErrAppStorageKVS,
				errPattern: regexp.MustCompile("error occurred while operating the key-value store"),
			},
		),
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			client := client{
				UniversalClient: tt.C.mock,
				timeout:         5 * time.Second,
				expiration:      10 * time.Second,
			}
			err := client.Set(tt.C.context, "", []byte{})
			testutil.DiffError(t, tt.A.err, tt.A.errPattern, err)
		})
	}
}

func TestDelete(t *testing.T) {
	type condition struct {
		mock    redis.UniversalClient
		context context.Context
	}

	type action struct {
		err        any // error or errorutil.Kind
		errPattern *regexp.Regexp
	}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"cmd.Err() returns some error",
			&condition{
				mock: &mocUniversalClient{
					redisError: true,
				},
				context: context.Background(),
			},
			&action{
				err:        app.ErrAppStorageKVS,
				errPattern: regexp.MustCompile("error occurred while operating the key-value store"),
			},
		),
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			client := client{
				UniversalClient: tt.C.mock,
				timeout:         5 * time.Second,
				expiration:      10 * time.Second,
			}
			err := client.Delete(tt.C.context, "")
			testutil.DiffError(t, tt.A.err, tt.A.errPattern, err)
		})
	}
}

func TestExists(t *testing.T) {
	type condition struct {
		mock    redis.UniversalClient
		context context.Context
	}

	type action struct {
		exists bool
	}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"cmd.Err() returns redis.Nil",
			&condition{
				mock: &mocUniversalClient{
					redisError: true,
				},
				context: context.Background(),
			},
			&action{
				exists: false,
			},
		),
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			client := client{
				UniversalClient: tt.C.mock,
				timeout:         5 * time.Second,
				expiration:      10 * time.Second,
			}
			exists := client.Exists(tt.C.context, "")
			testutil.Diff(t, tt.A.exists, exists)
		})
	}
}
