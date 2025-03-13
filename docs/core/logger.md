# Package `core/logger`

## Summary

This is the design document of `core/logger` package that provides Logger resources.
Logger is the resource that output logs to standard output, standard error or files.

## Motivation

Logging is required to work as an API gateway.

### Goals

- Provides built-in loggers.
- Make loggers available through interface.

### Non-Goals

- Provides various types of loggers.
- Make many options for built-in loggers such as formatting.

## Technical Design

### NoopLogger

NoopLogger is the logger that do nothing.
This logger does not output anything.

NoopLogger implements `logger.Logger` interface.

```go
type Logger interface {
  Dispatch(LogLevel, log.Creator) log.Creator
  Enabled(LogLevel) bool
  Debug(msg string, keyValues ...any)
  Info(msg string, keyValues ...any)
  Warn(msg string, keyValues ...any)
  Error(msg string, keyValues ...any)
}
```

### SLogger

SLogger is the logger that leverages [log/slog](https://pkg.go.dev/log/slog) package.
It provides both structured logging and unstructured logging.

Unstructured logging will be key-value format like below.

```txt
key1=val1 key2=val2 key3=val3 ...
```

Structured logging will be JSON format like below.

```json
{"key1":val1,"key2":val2,"key3":val3, ...}
```

NoopLogger implements `logger.Logger` interface.

```go
type Logger interface {
  Dispatch(LogLevel, log.Creator) log.Creator
  Enabled(LogLevel) bool
  Debug(msg string, keyValues ...any)
  Info(msg string, keyValues ...any)
  Warn(msg string, keyValues ...any)
  Error(msg string, keyValues ...any)
}
```

## Test Plan

### Unit Tests

Unit tests are implemented and passed.

- All functions and methods are covered.
- Coverage objective 98%.

### Integration Tests

Integration tests are implemented with these aspects.

- NoopLogger works as a Logger.
- NoopLogger works with input configuration.
- SLogger works as a Logger.
- SLogger works with input configuration.
- SLogger provides both structured and unstructured logging.

### e2e Tests

e2e tests are implemented with these aspects.

- NoopLogger works as a Logger.
- NoopLogger works with input configuration.
- SLogger works as a Logger.
- SLogger works with input configuration.
- SLogger provides both structured and unstructured logging.

### Fuzz Tests

Not planned.

### Benchmark Tests

Not planned.

### Chaos Tests

Not planned.

## Future works

- Add logger which can format logs.

## References

None.
