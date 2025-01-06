# gRPC Proxy Example

## About this example

This example shows how to configure reverse proxy for HTTP 3 client.

Files that required to run this example is shown below.
Make sure to build or download AILERON Gateway binary, `aileron`.

```txt
./
├── aileron
└── _example/
    └── proxy-http3/
        ├── pki/
        │   ├── cert.pem
        │   └── key.pem
        ├── server.go
        ├── client.go
        ├── config-http1.yaml
        ├── config-http2.yaml
        └── config-http3.yaml
```

Each configuration is supposed the following protocols.

| Config File          | Client Protocol | Upstream Protocol |
| -------------------- | --------------- | ----------------- |
| `config-http1.yaml`  | HTTP 3 (TLS)    | HTTP 1 (TLS)      |
| `config-http2.yaml`  | HTTP 3 (TLS)    | HTTP 2 (TLS)      |
| `config-http3.yaml`  | HTTP 3 (TLS)    | HTTP 3 (TLS)      |

This figure shows the overview of the this proxy example.
`client.go` can be used for the Client, `server.go` can be used for the Upstream.

```text
┌──────────┐          ┌─────────┐          ┌──────────┐
│          │  HTTP/3  │         │  HTTP/1  │          │
│  Client  │─────────►│  Proxy  │─────────►│ Upstream │
│          │  (TLS)   │         │  (TLS)   │          │
└──────────┘          └─────────┘          └──────────┘
┌──────────┐          ┌─────────┐          ┌──────────┐
│          │  HTTP/3  │         │  HTTP/2  │          │
│  Client  │─────────►│  Proxy  │─────────►│ Upstream │
│          │  (TLS)   │         │  (TLS)   │          │
└──────────┘          └─────────┘          └──────────┘
┌──────────┐          ┌─────────┐          ┌──────────┐
│          │  HTTP/3  │         │  HTTP/3  │          │
│  Client  │─────────►│  Proxy  │─────────►│ Upstream │
│          │  (TLS)   │         │  (TLS)   │          │
└──────────┘          └─────────┘          └──────────┘
```

## Run

Run the one of following commands to start a reverse proxy server.
A reverse proxy server starts and listens on [https://localhost:8443](https://localhost:8443).

```bash
./aileron -f _example/proxy-http3/config-http1.yaml
```

```bash
./aileron -f _example/proxy-http3/config-http2.yaml
```

```bash
./aileron -f _example/proxy-http3/config-http3.yaml
```

## Test

First, run the upstream server using `server.go`.
Run this command.

`server.go` will runs 3 servers.

- [https://localhost:10001](https://localhost:10001): HTTP 1 server
- [https://localhost:10002](https://localhost:10002): HTTP 2 server
- [https://localhost:10003](https://localhost:10003): HTTP 3 server

```bash
$ go run ./_example/proxy-http3/server.go

2024/08/26 03:32:34 HTTP 1 server listens on :10001
2024/08/26 03:32:34 HTTP 2 server listens on :10002
2024/08/26 03:32:34 HTTP 3 server listens on :10003
```

Then, send HTTP requests using `client.go` by running the command

```bash
go run ./_example/proxy-http3/client.go
```

For `config-http1.yaml`,

```bash
$ go run ./_example/proxy-http3/client.go

2024/08/26 03:33:30 Send HTTP 3 request : https://localhost:8443/test
2024/08/26 03:33:30 OK
2024/08/26 03:33:30 Method : GET
Path : /test
HTTP : 1.1
Header:
  X-Forwarded-Proto: [https]
  User-Agent: [quic-go HTTP/3]
  Accept-Encoding: [gzip]
  X-Forwarded-For: [127.0.0.1]
  X-Forwarded-Host: [localhost:8443]
  X-Forwarded-Port: [35269]
```

For `config-http2.yaml`,

```bash
$ go run ./_example/proxy-http3/client.go

2024/08/26 03:33:53 Send HTTP 3 request : https://localhost:8443/test
2024/08/26 03:33:53 OK
2024/08/26 03:33:53 Method : GET
Path : /test
HTTP : 2.0
Header:
  X-Forwarded-Host: [localhost:8443]
  X-Forwarded-Proto: [https]
  Accept-Encoding: [gzip]
  User-Agent: [quic-go HTTP/3]
  X-Forwarded-For: [127.0.0.1]
  X-Forwarded-Port: [56889]
```

For `config-http3.yaml`,

```bash
$ go run ./_example/proxy-http3/client.go

2024/08/26 03:34:10 Send HTTP 3 request : https://localhost:8443/test
2024/08/26 03:34:10 OK
2024/08/26 03:34:10 Method : GET
Path : /test
HTTP : 3.0
Header:
  X-Forwarded-For: [127.0.0.1]
  X-Forwarded-Proto: [https]
  X-Forwarded-Host: [localhost:8443]
  Accept-Encoding: [gzip]
  User-Agent: [quic-go HTTP/3]
  X-Forwarded-Port: [58061]
```
