# Proxy SSE Example

## About this example

This example shows reverse proxy handler can proxy SSE (Server sent event).

Files that required to run this example is shown below.
Make sure to build or download AILERON Gateway binary, `aileron`.

```txt
./
├── aileron
├── server  !!!!! This binary is built from the server.go in the Test section. 
└── _example/
    └── proxy-sse/
        ├── server.go
        └── config.yaml
```

## Run

Run the example with this command.
Reverse proxy server will listen on  [http://localhost:8080/](http://localhost:8080/).

```bash
./aileron -f _example/proxy-sse/
```

## Test

Before testing the SSE proxy, we need to build a SSE server.

So, run the following command to build a SSE server.
A binary `server` will be created.

```bash
go build _example/proxy-sse/server.go
```

Then run the SSE server.

```bash
$ ./server

2024/08/24 20:41:17 SSE server listens at 0.0.0.0:9999
```

Send a HTTP request like below.
An streaming response will be obtained.

Note that you can access to the URL from your browser.

```bash
$ curl localhost:8080

Hello !!
It's Sat, 24 Aug 2024 20:42:12 GMT
It's Sat, 24 Aug 2024 20:42:13 GMT
It's Sat, 24 Aug 2024 20:42:14 GMT
It's Sat, 24 Aug 2024 20:42:15 GMT

..... output omitted
```
