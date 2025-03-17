# Package `core/log`

## Summary

This is the design document of core/log package.
core/log package includes sub packages that provides LogCreator resources.
LogCreator resources creates log lines.

## Motivation

Logging is required to work as an API gateway.

### Goals

- Provides built-in log creators.
- Make loggers available through interface.

### Non-Goals

- Provides various types of log creators.
- Make many options for built-in log creators.

## Technical Design

### LogCreators

core/log package includes sub packages that provide LogCreator resources.
LogCreator is defined as resources that have the following interface.

```go
// Creator is defined in the kernel/log package.
type Creator interface {
   New(string, ...Attributes) Log
}
```

### MessageLogCreator

MessageLogCreator is one of the LogCreator that is defined in the kernel/log package.
See `kernel/log` package documents.

## Test Plan

### Unit Tests

Unit tests are implemented and passed.

- All functions and methods are covered.
- Coverage objective 98%.

### Integration Tests

Integration tests are implemented with these aspects.

- MessageLogCreator works as a LogCreator.
- MessageLogCreator works with input configuration.
- MessageLogCreator creates intended log lines.

### e2e Tests

e2e tests are implemented with these aspects.

- MessageLogCreator works as a LogCreator.
- MessageLogCreator works with input configuration.
- MessageLogCreator creates intended log lines.

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
