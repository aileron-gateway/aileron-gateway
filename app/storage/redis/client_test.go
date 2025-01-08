package redis

import (
	"context"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/aileron-gateway/aileron-gateway/app"
	"github.com/aileron-gateway/aileron-gateway/kernel/kvs"
	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
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

	CndClient := "set client"
	ActCheckNil := "expected nil returned"

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	tb.Condition(CndClient, "set client")
	tb.Action(ActCheckNil, "check that expected nil returned")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"set client",
			[]string{CndClient},
			[]string{ActCheckNil},
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

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			err := tt.C().client.Open(context.Background())
			testutil.Diff(t, tt.A().err, err)
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

	CndCloseReturnsNil := "set Close method returns nil"
	CndCloseReturnsNonNil := "set Close method returns non nil"
	ActCheckNil := "expected nil returned"
	ActCheckError := "expected error returned"

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	tb.Condition(CndCloseReturnsNil, "set Close method returns nil")
	tb.Condition(CndCloseReturnsNonNil, "set Close method returns non nil")
	tb.Action(ActCheckNil, "check that expected nil returned")
	tb.Action(ActCheckError, "check that expected error returned")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"Close func returns nil",
			[]string{CndCloseReturnsNil},
			[]string{ActCheckNil},
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
			[]string{CndCloseReturnsNonNil},
			[]string{ActCheckError},
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

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			client := client{
				UniversalClient: tt.C().mock,
				timeout:         5 * time.Second,
				expiration:      10 * time.Second,
			}
			err := client.Close(context.Background())
			testutil.DiffError(t, tt.A().err, tt.A().errPattern, err)
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

	CndRedisNil := "set cmd.Err() returns redis.nil"
	CndNonNilNonRedisNilError := "set cmd.Err() returns non nil and non redis.nil"
	ActCheckNil := "expected nil returned"
	ActCheckError := "expected error returned"

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	tb.Condition(CndRedisNil, "set cmd.Err() returns redis.nil")
	tb.Condition(CndNonNilNonRedisNilError, "set cmd.Err() returns non nil and non redis.nil")
	tb.Action(ActCheckNil, "check that expected nil returned")
	tb.Action(ActCheckError, "check that expected error returned")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"cmd.Err() returns redis.nil",
			[]string{CndRedisNil},
			[]string{ActCheckError},
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
			[]string{CndNonNilNonRedisNilError},
			[]string{ActCheckError},
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

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			client := client{
				UniversalClient: tt.C().mock,
				timeout:         5 * time.Second,
				expiration:      10 * time.Second,
			}

			value, err := client.Get(tt.C().context, "")
			testutil.DiffError(t, tt.A().err, tt.A().errPattern, err, cmpopts.EquateErrors())
			testutil.Diff(t, tt.A().value, value)
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

	CndNonNilError := "set cmd.Err() returns non nil"
	ActCheckNil := "expected nil returned"
	ActCheckError := "expected error returned"

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	tb.Condition(CndNonNilError, "set cmd.Err() returns non nil")
	tb.Action(ActCheckNil, "check that expected nil returned")
	tb.Action(ActCheckError, "check that expected error returned")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"cmd.Err() returns some error",
			[]string{CndNonNilError},
			[]string{ActCheckError},
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

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			client := client{
				UniversalClient: tt.C().mock,
				timeout:         5 * time.Second,
				expiration:      10 * time.Second,
			}
			err := client.Set(tt.C().context, "", []byte{})
			testutil.DiffError(t, tt.A().err, tt.A().errPattern, err)
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

	CndNonNilError := "set cmd.Err() returns non nil"
	ActCheckNil := "expected nil returned"
	ActCheckError := "expected error returned"

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	tb.Condition(CndNonNilError, "set cmd.Err() returns non nil")
	tb.Action(ActCheckNil, "check that expected nil returned")
	tb.Action(ActCheckError, "check that expected error returned")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"cmd.Err() returns some error",
			[]string{CndNonNilError},
			[]string{ActCheckError},
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

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			client := client{
				UniversalClient: tt.C().mock,
				timeout:         5 * time.Second,
				expiration:      10 * time.Second,
			}
			err := client.Delete(tt.C().context, "")
			testutil.DiffError(t, tt.A().err, tt.A().errPattern, err)
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

	CndRedisNil := "set cmd.Err() returns redis.Nil"
	ActCheckNil := "expected nil returned"
	ActCheckError := "expected error returned"

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	tb.Condition(CndRedisNil, "set cmd.Err() returns redis.Nil")
	tb.Action(ActCheckNil, "check that expected nil returned")
	tb.Action(ActCheckError, "check that expected error returned")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"cmd.Err() returns redis.Nil",
			[]string{CndRedisNil},
			[]string{ActCheckError},
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

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			client := client{
				UniversalClient: tt.C().mock,
				timeout:         5 * time.Second,
				expiration:      10 * time.Second,
			}
			exists := client.Exists(tt.C().context, "")
			testutil.Diff(t, tt.A().exists, exists)
		})
	}
}
