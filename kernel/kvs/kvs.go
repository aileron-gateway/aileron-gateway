package kvs

import (
	"context"
	"errors"
)

// Nil is the sentinel error returned by the operation methods against the key-value store.
// This error should not be wrapped by other error
// without implementing `Unwrap() error` interface.
//
//nolint:staticcheck // ST1012: error var Nil should have name of the form ErrFoo
//lint:ignore ST1012 // error var Nil should have name of the form ErrFoo
var Nil = errors.New("nil")

// OpenCloser opens and closes a key-value store.
// Open method should usually be called when connecting a key-value store.
// Close method should usually be called when disconnecting a key-value store.
type OpenCloser interface {
	// Open connects to the target key-value store.
	// It's implementers' responsible to define the behavior of this method.
	// In most case, connecting to the target key-value store and
	// initializing the client should be implemented.
	Open(context.Context) error

	// Close closes the connected key-value store.
	// It's implementers' responsible to define the behavior of this method.
	// It's callers' responsible to call this method
	// when disconnecting the key-value store.
	// In most case, dis-connecting to the target key-value store and
	// finalizing the client should be implemented.
	Close(context.Context) error
}

// Commander provides basic operations for the key-value Store.
type Commander[K comparable, V any] interface {
	// Get returns a value with the given key.
	// Get operation must return kvs.Nil if key was not exist in the store.
	// The sentinel error Nil must be returned when a value
	// associated to he given key was not found.
	Get(context.Context, K) (V, error)

	// Set sets value by the given key in the key-value store.
	Set(context.Context, K, V) error

	// Delete deletes values by the given key if exists.
	Delete(context.Context, K) error
}

// Client is the basic interface of a key-value Store.
type Client[K comparable, V any] interface {
	OpenCloser
	Commander[K, V]
}
