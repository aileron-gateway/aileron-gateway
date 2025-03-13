# Package `kernel/kvs`

## Summary

This is the design document of kernel/kvs package.
kernel/kvs package provides basic interface of key-value stores.

## Motivation

API Gateways uses key-value stores in some of their core features such as

- Session management
- Idempotent management
- CSRF token management

Some of the key-values that may be used are as follows.

- Memory (e.g. Go's map)
- HTTP header including Cookie
- [Redis](https://redis.io/)
- [Valkey](https://valkey.io/)
- [Memcached](https://memcached.org/)
- [Dragonfly](https://www.dragonflydb.io/)
- [SQLite](https://sqlite.org/)

So, it is reasonable to define common interface for key-value stores.  

### Goals

- Provide stable and common interface for key-value store.

### Non-Goals

- Implement client of a specific key-value store.

## Technical Design

### Connecting to stores

Connecting and dis-connecting to key-value stores are required to use remote stores over networks.
`OpenCloser` interface is defined for that case.

```go
type OpenCloser interface {
  // Open connects to the target key-value store and initializes the client.
  Open(context.Context) error
  // Close dis-connects the target key-value store and finalize the client.
  Close(context.Context) error
}
```

### Operation sets

`Commander` interface provides minimum operation sets for key-value stores.
Batch operations are not defined.

```go
type Commander[K comparable, V any] interface {
  // Get gets the value with the given key from the key-value store.
  Get(context.Context, K) (V, error)
  // Set sets the value with the given key to the key-value store.
  Set(context.Context, K, V) error
  // Delete deletes the value associated to the given key from the key-value store.
  Delete(context.Context, K) error
}
```

To unify the error that returned by the Commander's methods when no values are found,
a sentinel error named `Nil` is defined.

```go
var (
  Nil = errors.New("kernel/kvs: NIL")
)
```

The behavior of the basic operations when the given key was or was not found should follow the table.

| Operation | Key exists | Key not-exists     |
| --------- | ---------- | ------------------ |
| Get       | Success    | Nil error returned |
| Set       | Success    | Success            |
| Delete    | Success    | Success            |

### Client

The `Client` interface provides minimum operation set for key-value store clients.
It embeds `OpenCloser` and `Commander` interfaces.

```go
type Client[K comparable, V any] interface {
  OpenCloser
  Commander[K, V]
}
```

## Test Plan

### Unit Tests

Not planned.
This package only have interfaces.

### Integration Tests

Not planned.

### e2e Tests

Not planned.

### Fuzz Tests

Not planned.

### Benchmark Tests

Not planned.

### Chaos Tests

Not planned.

## Future works

None.

## References

None.
