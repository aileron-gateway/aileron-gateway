package redis

import (
	"context"
	"time"

	"github.com/aileron-gateway/aileron-gateway/app"
	"github.com/aileron-gateway/aileron-gateway/kernel/kvs"
	"github.com/redis/go-redis/v9"
)

// client is the session storage client that wraps redis.UniversalClient.
type client struct {
	redis.UniversalClient

	// timeout is the operation timeout.
	timeout time.Duration

	// expiration is the expiration duration in milliseconds that the data should be
	// stored in the redis.
	// -1 means KeepTTL.
	expiration time.Duration
}

// Open opens the key-value store.
// Do nothing for this redis storage.
func (c *client) Open(_ context.Context) error {
	return nil
}

// Close closes the key-value store.
// Close redis client for this client.
func (c *client) Close(_ context.Context) error {
	if err := c.UniversalClient.Close(); err != nil {
		return app.ErrAppStorageKVS.WithStack(err, nil)
	}
	return nil
}

// Get returns a value for the given key if exists in the key-value store.
// kvs.Nil is returned when the given key was not exist in hte key-value store.
// The context of the fir st argument must not be nil.
func (c *client) Get(ctx context.Context, key string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	cmd := c.UniversalClient.Get(ctx, key)
	if err := cmd.Err(); err != nil && err != redis.Nil {
		return nil, app.ErrAppStorageKVS.WithStack(err, nil)
	}

	b, err := cmd.Bytes()
	if err == redis.Nil {
		// kvs.Nil must be returned to represent that there is no value in the storage.
		return nil, kvs.Nil
	}

	return b, nil
}

// Set sets the value in the key-value store with the given key with the given expiration.
// The context of the fir st argument must not be nil.
func (c *client) Set(ctx context.Context, key string, value []byte) error {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	cmd := c.UniversalClient.Set(ctx, key, value, c.expiration)
	if err := cmd.Err(); err != nil {
		return app.ErrAppStorageKVS.WithStack(err, nil)
	}

	return nil
}

// SetEx sets the value in the key-value store with the given key with the given expiration.
// The context of the fir st argument must not be nil.
func (c *client) SetWithTTL(ctx context.Context, key string, value []byte, exp time.Duration) error {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	cmd := c.UniversalClient.SetEx(ctx, key, value, exp)
	if err := cmd.Err(); err != nil {
		return app.ErrAppStorageKVS.WithStack(err, nil)
	}

	return nil
}

// Delete deletes the data from key-value store with the given key.
// The context of the fir st argument must not be nil.
func (c *client) Delete(ctx context.Context, key string) error {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	cmd := c.UniversalClient.Del(ctx, key)
	if err := cmd.Err(); err != nil {
		return app.ErrAppStorageKVS.WithStack(err, nil)
	}

	return nil
}

// Exists returns if the given key is exists in the key-value store.
// The context of the fir st argument must not be nil.
func (c *client) Exists(ctx context.Context, key string) bool {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	// n is the number of keys exists in the redis.
	// See https://redis.io/commands/exists/
	n, err := c.UniversalClient.Exists(ctx, key).Result()
	if err == redis.Nil {
		return false
	}

	return n > 0
}
