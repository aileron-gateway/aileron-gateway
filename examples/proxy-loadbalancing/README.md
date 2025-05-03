# Reverse Proxy Example

## About this example

This example shows how to configure reverse proxy handler that proxy requests to upstream services.

Files that required to run this example is shown below.
Make sure to build or download AILERON Gateway binary, `aileron`.

```txt
./
├── aileron
├── server  !!!!! This binary is built from the server.go in the Test section. 
└── _example/
    └── proxy-loadbalancing/
        ├── server.go
        ├── config-direct-hash.yaml
        ├── config-maglev.yaml
        ├── config-random.yaml
        ├── config-ring-hash.yaml
        └── config-round-robin.yaml
```

## Run

Run the example with this command.
Reverse proxy server will listen on  [http://localhost:8080/](http://localhost:8080/).

```bash
# Use DirectHash load balancer.
./aileron -f _example/proxy-loadbalancing/config-direct-hash.yaml
```

```bash
# Use Maglev load balancer.
./aileron -f _example/proxy-loadbalancing/config-maglev.yaml
```

```bash
# Use Random load balancer.
./aileron -f _example/proxy-loadbalancing/config-random.yaml
```

```bash
# Use RingHash load balancer.
./aileron -f _example/proxy-loadbalancing/config-ring-hash.yaml
```

```bash
# Use RoundRobin load balancer.
./aileron -f _example/proxy-loadbalancing/config-round-robin.yaml
```

## Test

Before testing the load balancers, we need to build a upstream server.

So, run the following command to build a server.
A binary `server` will be created.

```bash
go build _example/proxy-loadbalancing/server.go
```

Then, run the server.

It runs 5 HTTP servers by default.

- <http://localhost:8001>
- <http://localhost:8002>
- <http://localhost:8003>
- <http://localhost:8004>
- <http://localhost:8005>

```bash
$ ./server

2024/08/25 23:07:21 Server listens at :8001
2024/08/25 23:07:21 Server listens at :8002
2024/08/25 23:07:21 Server listens at :8003
2024/08/25 23:07:21 Server listens at :8004
2024/08/25 23:07:21 Server listens at :8005
```

Send a HTTP request like below.

This request will be proxied by client ip and port
because there is no header named `Hash-Key`.

```bash
# RingHash proxy is used in this example.
$ curl http://localhost:8080/

Hello! 127.0.0.1:45758 from :8004  // 1st request response example.
Hello! 127.0.0.1:48032 from :8002  // 2nd
Hello! 127.0.0.1:54380 from :8005  // 3rd
Hello! 127.0.0.1:45758 from :8004  // 4th
```

This request will be proxied by the header value of `Hash-Key`.
So the request reaches to the same upstream every time.

```bash
# RingHash proxy is used in this example.
$ curl -H "Hash-Key: foo" http://localhost:8080/

Hello! 127.0.0.1:60998 from :8005  // 1st request response example.
Hello! 127.0.0.1:60998 from :8005  // 2nd
Hello! 127.0.0.1:60998 from :8005  // 3rd
Hello! 127.0.0.1:60998 from :8005  // 4th
```
