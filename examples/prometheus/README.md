# Prometheus Example

## About this example

This example shows how to use PrometheusMeter that leverages [Prometheus](https://prometheus.io/).
PrometheusMeter works as a prometheus exporter.

Files that required to run this example is shown below.
Make sure to build or download AILERON Gateway binary, `aileron`.

```txt
./
├── aileron
└── _example/
    └── prometheus/
        └── config.yaml
```

## Run

Run the example with this command.
A reverse proxy server will listen on [http://localhost:8080/](http://localhost:8080/).

```bash
./aileron -f _example/prometheus/
```

## Test

Send HTTP requests to `/metrics` to access to the metrics endpoint.
Note that the `/metrics` path can be changed from the config.

```bash
$ curl http://localhost:8080/metrics

# HELP go_gc_duration_seconds A summary of the pause duration of garbage collection cycles.
# TYPE go_gc_duration_seconds summary
go_gc_duration_seconds{quantile="0"} 2.11e-05
go_gc_duration_seconds{quantile="0.25"} 7e-05
...
```

In this example, PrometheusMeter is inserted as middleware to export the number of API calls.

```yaml
middleware:
- apiVersion: app/v1
    kind: PrometheusMeter
    namespace: example
    name: default
```
