# gRPC Proxy Example

## About this example

This example shows how to configure reverse proxy for gRPC.

Files that required to run this example is shown below.
Make sure to build or download AILERON Gateway binary, `aileron`.

This example uses gRPC server and client provided by

- <https://grpc.io/docs/languages/go/basics/>
- <https://github.com/grpc/grpc-go/tree/master/examples/route_guide>

```txt
./
├── aileron
└── _example/
    └── proxy-grpc/
        ├── pki/
        │   ├── cert.pem
        │   └── key.pem
        ├── config-http-http.yaml
        ├── config-https-http.yaml
        ├── config-http-https.yaml
        └── config-https-https.yaml

https://github.com/grpc/grpc-go/
└── examples/
    └── route_guide/
        ├── client/
        │   └── client.go
        └── server/
            └── server.go
```

Because HTTP2 allow non-TLS communications,
here we show 4 variation of configs with/without TLS.

| Config File               | gRPC Client - Reverse Proxy | Reverse Proxy - gRPC Server |
| ------------------------- | --------------------------- | --------------------------- |
| `config-http-http.yaml`   | non-TLS                     | non-TLS                     |
| `config-https-http.yaml`  | TLS                         | non-TLS                     |
| `config-http-https.yaml`  | non-TLS                     | TLS                         |
| `config-https-https.yaml` | TLS                         | TLS                         |

## Run

Use appropriate config file to run a reverse proxy.

See the successive [Test](#test) section.

The command is like this for all configs.

```bash
./aileron -f _example/proxy-grpc/<Use appropriate config>
```

## Test

First, testing for the `config-http-http.yaml` is shown.

See the later descriptions to use TLS server or TLS client.

Run the reverse proxy server with this command.
Reverse proxy server will listen on  [http://localhost:50000/](http://localhost:50000/).

```bash
./aileron -f _example/proxy-grpc/config-http-http.yaml
```

Then start the upstream gRPC server without TLS.

```bash
# @ https://github.com/grpc/grpc-go/tree/master/examples/route_guide/server
$ go build
$ ./server -port 50051
```

Send gRPC requests with the client like below.
Responses will be shown in the terminal.

All 4 types of gRPC communications are done within the command.

- [Simple RPC](https://grpc.io/docs/languages/go/basics/#simple-rpc)
- [Server-side streaming RPC](https://grpc.io/docs/languages/go/basics/#server-side-streaming-rpc)
- [Client-side streaming RPC](https://grpc.io/docs/languages/go/basics/#client-side-streaming-rpc)
- [Bidirectional streaming RPC](https://grpc.io/docs/languages/go/basics/#bidirectional-streaming-rpc)

```bash
# @ https://github.com/grpc/grpc-go/tree/master/examples/route_guide/client
$ go build
$ ./client -addr localhost:50000

2024/08/24 02:42:36 Getting feature for point (409146138, -746188906)
2024/08/24 02:42:36 name:"Berkshire Valley Management Area Trail, Jefferson, NJ, USA" location:{latitude:409146138 longitude:-746188906}
2024/08/24 02:42:36 Getting feature for point (0, 0)
2024/08/24 02:42:36 location:{}
2024/08/24 02:42:36 Looking for features within lo:{latitude:400000000 longitude:-750000000} hi:{latitude:420000000 longitude:-730000000}
... output omitted
```

**gRPC server and client commands for TLS**.

`config-https-http.yaml` and `config-https-https.yaml` require TLS for the gRPC client.
So, run the gRPC client like below.

```bash
# Copy ./pki/cert.pem to the directory 
# where the gRPC client binary exists.
$ ./client -addr localhost:50000 -tls -ca_file cert.pem -server_host_override localhost
```

`config-http-https.yaml` and `config-https-https.yaml` require TLS for the gRPC server.
So, run the gRPC server like below.

```bash
# Copy ./pki/cert.pem and ./pki/key.pem to the directory
# where the gRPC server binary exists.
$ ./server -port 50051 -tls -key_file key.pem -cert_file cert.pem
```


# memo

```
go run -mod=mod google.golang.org/grpc/examples/route_guide/server
go run -mod=mod google.golang.org/grpc/examples/route_guide/client
```

```
go run -mod=mod google.golang.org/grpc/examples/helloworld/greeter_server
go run -mod=mod google.golang.org/grpc/examples/helloworld/greeter_client
```
