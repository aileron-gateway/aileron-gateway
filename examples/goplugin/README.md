# GoPlugin Example

## About this example

This example shows how to extend AILERON Gateway using GoPlugin.

**Note: GoPlugin can be used on Linux,

Files that required to run this example is shown below.
Make sure to build or download AILERON Gateway binary, `aileron`.

```txt
./
├── aileron
└── _example/
    └── goplugin/
        ├── plugin.go
        └── config.yaml
```

## Run

Before running aileron binary, we have to get dynamically build aileron binary.
So, build a binary with `CGO_ENABLED=1`.

```bash
export CGO_ENABLED=1
make build
```

Next, generate a shared object from `goplugin.go` by the command.
**Because of the limitation of GoPlugin, options such as `-trimpath` or `-tag` MUST be
the same as the options aileron binary built with.**

```bash
cd <This directory>
go build -buildmode=plugin -trimpath -tags="netgo,osusergo"
```

Run the example with this command.
Reverse proxy server will listen on  [http://localhost:8080/](http://localhost:8080/).

```bash
./aileron -f _example/goplugin/
```

## Test

Send a HTTP request like below.
The default upstream server [http://httpbin.org/](http://httpbin.org/) will return response.

We can see that the `Alice: bob` is added to the request header and `Foo: bar` is added to the response header.
Both are added in the plugin.

```bash
$ curl -v http://localhost:8080/get

> GET /get HTTP/1.1
> Host: localhost:8080
> User-Agent: curl/7.68.0
> Accept: */*

< HTTP/1.1 200 OK
< Access-Control-Allow-Credentials: true
< Access-Control-Allow-Origin: *
< Content-Length: 330
< Content-Type: application/json
< Date: Sun, 08 Sep 2024 08:14:15 GMT
< Foo: bar
< Server: gunicorn/19.9.0

{
  "args": {},
  "headers": {
    "Accept": "*/*",
    "Alice": "bob",
    "Host": "httpbin.org",
    "User-Agent": "curl/7.68.0",
    "X-Amzn-Trace-Id": "Root=1-66dd5cd7-6498359e3515af1162a50300",
    "X-Forwarded-Host": "localhost:8080"
  },
  "origin": "127.0.0.1, 106.73.5.65",
  "url": "http://localhost:8080/get"
}
```

Because the example plugin is also registered as a HTTP handler,
GET request to the `/goplugin` returns response from the plugin.

```bash
$ curl http://localhost:8080/goplugin

GoPlugin Handler !!
```
