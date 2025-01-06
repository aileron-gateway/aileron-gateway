# Static File Server Example

## About this example

This example shows how to configure static file server that returns static contents in a directory.

Files that required to run this example is shown below.
Make sure to build or download AILERON Gateway binary, `aileron`.

Static contents in the root/ directory will be served.

```txt
./
├── aileron
└── _example/
    └── static-server/
        ├── config.yaml
        └── root/
            ├── index.html
            └── content.json
```

## Run

```bash
./aileron -f _example/static-server/
```

## Test

Send a HTTP request that returns `_example/static-server/root/index.html`

```bash
$ curl http://localhost:8080/

<!doctype html>
<html>
  <head>
    <title>AILERON Gateway</title>
  </head>
  <body>
    <h1>AILERON Gateway</h1>
  </body>
</html>
```

Another example serving content.json file.

```bash
curl -v http://localhost:8080/content.json

> GET /content.json HTTP/1.1
> Host: localhost:8080
> User-Agent: curl/7.68.0
> Accept: */*

< HTTP/1.1 200 OK
< Accept-Ranges: bytes
< Cache-Control: no-cache
< Content-Length: 34
< Content-Type: application/json
< Last-Modified: Wed, 31 Jan 2024 02:24:06 GMT
< X-Content-Type-Options: nosniff
< Date: Sun, 21 Jul 2024 20:39:47 GMT

{
    "Hello": "AILERON Gateway"
}
```
