# Package `core/httpserver` for `HTTPServer`

## Summary

This is the design document of `core/httpserver` package which provides `HTTPServer` resource.
`HTTPServer` can run HTTP 1/2/3 servers.

## Motivation

HTTP server is required to work as an API gateway.

### Goals

- HTTPServer can run HTTP1 server.
    - TLS is configurable.
- HTTPServer can run HTTP2 server.
- HTTPServer can run HTTP3 server.

### Non-Goals

## Technical Design

### HTTP server

HTTPServer runs a HTTP server.
HTTPServer can run HTTP1, HTTP2 and HTTP3 server.

HTTPServer implements `core.Runner` interface.
HTTPServer is intended to be run by registering to the Entrypoint resource.

```go
type Runner interface {
  Run(context.Context) error
}
```

Following 3 types of servers can be run by configuring with the configuration.

- HTTP1 server
    - HTTPServer can run HTTP1 server by leveraging [net/http](https://pkg.go.dev/net/http#hdr-HTTP_2)
    - HTTP2 server is enabled by default and [can be disabled](https://pkg.go.dev/net/http#hdr-HTTP_2).
    - [http.Server](https://pkg.go.dev/net/http#Server) is used.
    - TLS is configurable.
- HTTP2 server
    - HTTPServer can run HTTP2 server by leveraging [net/http](https://pkg.go.dev/net/http#hdr-HTTP_2) and [golang.org/x/net/http2](https://pkg.go.dev/golang.org/x/net/http2).
    - [http2.Server](https://pkg.go.dev/golang.org/x/net/http2#Server) is used.
    - TLS is required.
- HTTP3 server
    - HTTPServer can run HTTP3 server by leveraging [github.com/quic-go/quic-go/http3](github.com/quic-go/quic-go/http3).
    - [http3.Server](https://pkg.go.dev/github.com/quic-go/quic-go/http3#Server) is used.

## Test Plan

### Unit Tests

Unit tests are implemented and passed.

- All functions and methods are covered.
- Coverage objective 98%.

### Integration Tests

Integration tests are implemented with these aspects.

- HTTPServer works as a runner.
- HTTPServer works with input configuration.
- HTTPServer can run HTTP1 server with/without TLS.
- HTTPServer can run HTTP2 server with TLS.
- HTTPServer can run HTTP3 server.

### e2e Tests

e2e tests are implemented with these aspects.

- HTTPServer works as a runner.
- HTTPServer works with input configuration.
- HTTPServer can run HTTP1 server with/without TLS.
- HTTPServer can run HTTP2 server with TLS.
- HTTPServer can run HTTP3 server.

### Fuzz Tests

Not planned.

### Benchmark Tests

Not planned.

### Chaos Tests

Not planned.

## Future works

None.

## References

- [net/http - pkg.go](https://pkg.go.dev/net/http#Server)
- [golang.org/x/net/http2 - pkg.go](https://pkg.go.dev/golang.org/x/net/http2#Server)
- [github.com/quic-go/quic-go/http3 - pkg.go](https://pkg.go.dev/github.com/quic-go/quic-go/http3#Server)
