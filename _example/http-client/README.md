# HTTP Client Example

## About this example

This example shows how to configure http client.
Configured http client is used in the reverse proxy.

Files that required to run this example is shown below.
Make sure to build or download AILERON Gateway binary, `aileron`.

```txt
./
├── aileron
└── _example/
    └── http-client/
        └── config.yaml
```

## Run

Run the example with this command.
Reverse proxy server will listen on  [http://localhost:8080/](http://localhost:8080/).

```bash
./aileron -f _example/http-client/
```

## Test

Send a HTTP request like below.
The default upstream server of the reverse proxy is [http://httpbin.org/](http://httpbin.org/)

This request return 200 OK status.
No error logs in the gateway log.

```bash
$ curl -v http://localhost:8080/status/200

> GET /status/200 HTTP/1.1
> Host: localhost:8080
> User-Agent: curl/7.68.0
> Accept: */*

< HTTP/1.1 200 OK
< Access-Control-Allow-Credentials: true
< Access-Control-Allow-Origin: *
< Content-Length: 0
< Content-Type: text/html; charset=utf-8
< Date: Tue, 30 Jul 2024 13:30:33 GMT
< Server: gunicorn/19.9.0
```

Next, send the following request.
A http response with status code 500 will be returned.
An error log will be output in the gateway log.
The error log is telling that the request was retried some times because the http client is configured to retry when error responses with status code 500 were returned.

```bash
$ curl -v http://localhost:8080/status/500

> GET /status/500 HTTP/1.1
> Host: localhost:8080
> User-Agent: curl/7.68.0
> Accept: */*

< HTTP/1.1 500 Internal Server Error
< Content-Type: application/json; charset=utf-8
< Vary: Accept
< X-Content-Type-Options: nosniff
< Date: Tue, 30 Jul 2024 13:32:21 GMT
< Content-Length: 51

{"status":500,"statusText":"Internal Server Error"}
```

This is a part of the error log that tells the request was retried by the configured http client.

```txt
0th request failed at unix millis 1722346339104 by core/client: returned status code is 500; 
1th request failed at unix millis 1722346339974 by core/client: returned status code is 500;
2th request failed at unix millis 1722346341267 by core/client: returned status code is 500;
3th request failed at unix millis 1722346341652 by core/client: returned status code is 500
```
