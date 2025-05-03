# Access Logging Example

## About this example

This example shows how to get access logs.

Files that required to run this example is shown below.
Make sure to build or download AILERON Gateway binary, `aileron`.

```txt
./
├── aileron
└── _example/
    └── access-logging/
        └── config.yaml
```

## Run

A reverse proxy server with access logger will be run with the command.
The server listens on [http://localhost:8080/](http://localhost:8080/).

```bash
./aileron -f _example/access-logging/
```

## Test

Send a HTTP request like below.
The default upstream [http://httpbin.org/](http://httpbin.org/) will return response.

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
    "X-Amzn-Trace-Id": "Root=1-669bf9d7-570ac10959cfd49c16a68f3c",
    "X-Forwarded-Host": "localhost:8080"
  },
  "origin": "127.0.0.1, 106.73.5.65",
  "url": "http://localhost:8080/get"
}
```
