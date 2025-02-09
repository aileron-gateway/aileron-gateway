# OpenTelemetry Example

This folder is the example of using OpenTelemetry Tracing and Metrics.  
This folder contains multiple config files, each of which can be used for the following purposes.

## metrics-config.yaml

In the metrics-config.yaml file, the collected metrics are sent to `localhost:4318` using the HTTPExporter.  
These metrics can be viewed in Prometheus and Grafana.

## tracer-config.yaml

In the tracer-config.yaml file, trace data generated when passing through the TrackingMiddleware and ReverseProxyHandler is sent to `localhost:4317` using the gRPCExporter.  
This trace data can be viewed in Jaeger.
