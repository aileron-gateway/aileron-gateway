# gRPC Proxy Example

## About this example

This example shows how to configure reverse proxy for HTTP 2 client.

Files that required to run this example is shown below.
Make sure to build or download AILERON Gateway binary, `aileron`.

```txt
./
├── aileron
└── _example/
    └── proxy-http2/
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
| `config-http1.yaml`  | HTTP 2 (TLS)    | HTTP 1 (TLS)      |
| `config-http2.yaml`  | HTTP 2 (TLS)    | HTTP 2 (TLS)      |
| `config-http3.yaml`  | HTTP 2 (TLS)    | HTTP 3 (TLS)      |

This figure shows the overview of the this proxy example.
`client.go` can be used for the Client, `server.go` can be used for the Upstream.

```text
┌──────────┐          ┌─────────┐          ┌──────────┐
│          │  HTTP/2  │         │  HTTP/1  │          │
│  Client  │─────────►│  Proxy  │─────────►│ Upstream │
│          │  (TLS)   │         │  (TLS)   │          │
└──────────┘          └─────────┘          └──────────┘
┌──────────┐          ┌─────────┐          ┌──────────┐
│          │  HTTP/2  │         │  HTTP/2  │          │
│  Client  │─────────►│  Proxy  │─────────►│ Upstream │
│          │  (TLS)   │         │  (TLS)   │          │
└──────────┘          └─────────┘          └──────────┘
┌──────────┐          ┌─────────┐          ┌──────────┐
│          │  HTTP/2  │         │  HTTP/3  │          │
│  Client  │─────────►│  Proxy  │─────────►│ Upstream │
│          │  (TLS)   │         │  (TLS)   │          │
└──────────┘          └─────────┘          └──────────┘
```

## Run

Run the one of following commands to start a reverse proxy server.
A reverse proxy server starts and listens on [https://localhost:8443](https://localhost:8443).

```bash
./aileron -f _example/proxy-http2/config-http1.yaml
```

```bash
./aileron -f _example/proxy-http2/config-http2.yaml
```

```bash
./aileron -f _example/proxy-http2/config-http3.yaml
```

## Test

First, run the upstream server using `server.go`.
Run this command.

`server.go` will runs 3 servers.

- [https://localhost:10001](https://localhost:10001): HTTP 1 server
- [https://localhost:10002](https://localhost:10002): HTTP 2 server
- [https://localhost:10003](https://localhost:10003): HTTP 3 server

```bash
$ go run ./_example/proxy-http2/server.go

2024/08/26 03:26:59 HTTP 1 server listens on :10001
2024/08/26 03:26:59 HTTP 2 server listens on :10002
2024/08/26 03:26:59 HTTP 3 server listens on :10003
```

Then, send HTTP requests using `client.go` by running the command

```bash
go run ./_example/proxy-http2/client.go
```

For `config-http1.yaml`,

```bash
$ go run ./_example/proxy-http2/client.go

2024/08/26 03:27:54 Send HTTP 1 request : https://localhost:8443/test
2024/08/26 03:27:54 OK
2024/08/26 03:27:54 Method : GET
Path : /test
HTTP : 1.1
Header:
  Accept-Encoding: [gzip]
  X-Forwarded-For: [127.0.0.1]
  X-Forwarded-Host: [localhost:8443]
  X-Forwarded-Port: [41040]
  X-Forwarded-Proto: [https]
  User-Agent: [Go-http-client/2.0]
```

For `config-http2.yaml`,

```bash
$ go run ./_example/proxy-http2/client.go

2024/08/26 03:28:31 Send HTTP 1 request : https://localhost:8443/test
2024/08/26 03:28:31 OK
2024/08/26 03:28:31 Method : GET
Path : /test
HTTP : 2.0
Header:
  X-Forwarded-Port: [48934]
  X-Forwarded-Host: [localhost:8443]
  X-Forwarded-Proto: [https]
  Accept-Encoding: [gzip]
  User-Agent: [Go-http-client/2.0]
  X-Forwarded-For: [127.0.0.1]
```

For `config-http3.yaml`,

```bash
$ go run ./_example/proxy-http2/client.go

2024/08/26 03:28:45 Send HTTP 1 request : https://localhost:8443/test
2024/08/26 03:28:46 OK
2024/08/26 03:28:46 Method : GET
Path : /test
HTTP : 3.0
Header:
  X-Forwarded-Port: [33132]
  X-Forwarded-Host: [localhost:8443]
  X-Forwarded-Proto: [https]
  Accept-Encoding: [gzip]
  User-Agent: [Go-http-client/2.0]
  X-Forwarded-For: [127.0.0.1]
```
