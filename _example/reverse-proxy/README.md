# Reverse Proxy Example

## About this example

This example shows how to configure reverse proxy handler that proxy requests to upstream services.

Files that required to run this example is shown below.
Make sure to build or download AILERON Gateway binary, `aileron`.

```txt
./
├── aileron
└── _example/
    └── reverse-proxy/
        └── config.yaml
```

## Run

Run the example with this command.
Reverse proxy server will listen on  [http://localhost:8080/](http://localhost:8080/).

```bash
./aileron -f _example/reverse-proxy/
```

## Test

Send a HTTP request like below.
The default upstream server [http://httpbin.org/](http://httpbin.org/) will return response.

```bash
$ curl http://localhost:8080/get

{
  "args": {},
  "headers": {
    "Accept": "*/*",
    "Accept-Encoding": "gzip",
    "Forwarded": "for=\"127.0.0.1\";host=\"localhost:8080\";proto=http",
    "Host": "httpbin.org",
    "User-Agent": "curl/7.68.0",
    "X-Amzn-Trace-Id": "Root=1-669d7745-1efcc4a64551d8785c2408a2",
    "X-Forwarded-Host": "localhost:8080"
  },
  "origin": "127.0.0.1, 106.73.5.65",
  "url": "http://localhost:8080/get"
}
```

POST request is configured to be rejected.

```bash
$ curl -X POST http://localhost:8080/post

{"status":404,"statusText":"Not Found"}
```
