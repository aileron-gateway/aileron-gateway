package oteltracer

import (
	"cmp"
	"context"
	"reflect"
	"time"

	v1 "github.com/aileron-gateway/aileron-gateway/apis/app/v1"
	"github.com/aileron-gateway/aileron-gateway/apis/kernel"
	"github.com/aileron-gateway/aileron-gateway/core"
	"github.com/aileron-gateway/aileron-gateway/kernel/api"
	"github.com/aileron-gateway/aileron-gateway/kernel/network"
	"go.opentelemetry.io/contrib/propagators/autoprop"
	"go.opentelemetry.io/contrib/propagators/aws/xray"
	"go.opentelemetry.io/contrib/propagators/b3"
	"go.opentelemetry.io/contrib/propagators/jaeger"
	"go.opentelemetry.io/contrib/propagators/opencensus"
	"go.opentelemetry.io/contrib/propagators/ot"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"google.golang.org/grpc/credentials"
	"google.golang.org/protobuf/reflect/protoreflect"
)

const (
	apiVersion = "app/v1"
	kind       = "OpenTelemetryTracer"
	Key        = apiVersion + "/" + kind
)

var Resource api.Resource = &API{
	BaseResource: &api.BaseResource{},
}

type API struct {
	*api.BaseResource
}

func (o *API) Default() protoreflect.ProtoMessage {
	return &v1.OpenTelemetryTracer{
		APIVersion: apiVersion,
		Kind:       kind,
		Metadata: &kernel.Metadata{
			Namespace: "default",
			Name:      "default",
		},
		Spec: &v1.OpenTelemetryTracerSpec{
			ServiceName: "gateway",
			LibraryName: reflect.TypeOf(*o).PkgPath(),
			Exporters: &v1.OpenTelemetryTracerSpec_GRPCExporterSpec{
				GRPCExporterSpec: &v1.GRPCTraceExporterSpec{},
			},
			TracerProviderBatch: &v1.TracerProviderBatchSpec{
				MaxQueueSize:       2048,
				BatchTimeout:       5,
				ExportTimeout:      30,
				MaxExportBatchSize: 512,
				Blocking:           false,
			},
			TracerProviderLimit: &v1.TracerProviderLimitSpec{
				AttributeValueLengthLimit:   -1,
				AttributeCountLimit:         -1,
				EventCountLimit:             -1,
				LinkCountLimit:              -1,
				AttributePerEventCountLimit: -1,
				AttributePerLinkCountLimit:  -1,
			},
			TraceIDRatioBased: 1.0,
		},
	}
}

func (*API) Create(a api.API[*api.Request, *api.Response], msg protoreflect.ProtoMessage) (any, error) {
	c := msg.(*v1.OpenTelemetryTracer)

	var exporter sdktrace.SpanExporter
	var err error
	switch c.Spec.Exporters.(type) {
	case *v1.OpenTelemetryTracerSpec_HTTPExporterSpec:
		exporter, err = newHTTPExporter(c.Spec.GetHTTPExporterSpec())
	case *v1.OpenTelemetryTracerSpec_GRPCExporterSpec:
		exporter, err = newGRPCExporter(c.Spec.GetGRPCExporterSpec())
	case *v1.OpenTelemetryTracerSpec_StdoutExporterSpec:
		exporter = newStdoutExporter(c.Spec.GetStdoutExporterSpec())
	case *v1.OpenTelemetryTracerSpec_ZipkinExporterSpec:
		exporter, err = newZipkinExporter(c.Spec.GetZipkinExporterSpec())
	}
	if err != nil {
		return nil, core.ErrCoreGenCreateObject.WithStack(err, map[string]any{"kind": kind})
	}

	tpl := c.Spec.TracerProviderLimit
	topts := []sdktrace.TracerProviderOption{
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(c.Spec.ServiceName),
		)),
		sdktrace.WithRawSpanLimits(
			sdktrace.SpanLimits{
				AttributeValueLengthLimit:   int(tpl.AttributeValueLengthLimit),
				AttributeCountLimit:         int(tpl.AttributeCountLimit),
				EventCountLimit:             int(tpl.EventCountLimit),
				LinkCountLimit:              int(tpl.LinkCountLimit),
				AttributePerEventCountLimit: int(tpl.AttributePerEventCountLimit),
				AttributePerLinkCountLimit:  int(tpl.AttributePerLinkCountLimit),
			},
		),
		sdktrace.WithSampler(sdktrace.TraceIDRatioBased(float64(c.Spec.TraceIDRatioBased))),
	}

	tpb := c.Spec.TracerProviderBatch
	if tpb.Blocking {
		topts = append(topts, sdktrace.WithBatcher(
			exporter,
			sdktrace.WithMaxQueueSize(int(tpb.MaxQueueSize)),
			sdktrace.WithBatchTimeout(time.Duration(tpb.BatchTimeout)*time.Second),
			sdktrace.WithExportTimeout(time.Duration(tpb.ExportTimeout)*time.Second),
			sdktrace.WithMaxExportBatchSize(int(tpb.MaxExportBatchSize)),
			sdktrace.WithBlocking(),
		))
	} else {
		topts = append(topts, sdktrace.WithBatcher(
			exporter,
			sdktrace.WithMaxQueueSize(int(tpb.MaxQueueSize)),
			sdktrace.WithBatchTimeout(time.Duration(tpb.BatchTimeout)*time.Second),
			sdktrace.WithExportTimeout(time.Duration(tpb.ExportTimeout)*time.Second),
			sdktrace.WithMaxExportBatchSize(int(tpb.MaxExportBatchSize)),
		))
	}

	tracerProvider := sdktrace.NewTracerProvider(topts...)
	tracer := tracerProvider.Tracer(c.Spec.LibraryName)

	return &otelTracer{
		tracer:  tracer,
		tp:      tracerProvider,
		pg:      autoprop.NewTextMapPropagator(propagators(c.Spec.PropagationTypes)...),
		headers: c.Spec.Headers,
	}, nil
}

func propagators(types []v1.PropagationType) []propagation.TextMapPropagator {
	if len(types) == 0 {
		return []propagation.TextMapPropagator{propagation.TraceContext{}, propagation.Baggage{}}
	}
	props := make([]propagation.TextMapPropagator, 0, len(types))
	for _, pt := range types {
		switch pt {
		case v1.PropagationType_W3CTraceContext:
			props = append(props, propagation.TraceContext{})
		case v1.PropagationType_W3CBaggage:
			props = append(props, propagation.Baggage{})
		case v1.PropagationType_B3:
			props = append(props, b3.New())
		case v1.PropagationType_Jaeger:
			props = append(props, jaeger.Jaeger{})
		case v1.PropagationType_XRay:
			props = append(props, xray.Propagator{})
		case v1.PropagationType_OpenCensus:
			props = append(props, opencensus.Binary{})
		case v1.PropagationType_OpenTracing:
			props = append(props, ot.OT{})
		}
	}
	return props
}

func newHTTPExporter(spec *v1.HTTPTraceExporterSpec) (*otlptrace.Exporter, error) {
	var opts []otlptracehttp.Option
	opts = appendOption(spec.EndpointURL != "", opts, otlptracehttp.WithEndpointURL(spec.EndpointURL))
	opts = append(opts, otlptracehttp.WithHeaders(spec.Headers))
	opts = appendOption(spec.Compress, opts, otlptracehttp.WithCompression(otlptracehttp.GzipCompression))
	opts = appendOption(spec.Insecure, opts, otlptracehttp.WithInsecure())
	opts = appendOption(spec.Timeout > 0, opts, otlptracehttp.WithTimeout(time.Duration(spec.Timeout)*time.Second))

	tlsConfig, err := network.TLSConfig(spec.TLSConfig)
	if err != nil {
		return nil, err
	}
	opts = appendOption(tlsConfig != nil, opts, otlptracehttp.WithTLSClientConfig(tlsConfig))

	retry := cmp.Or(spec.OTLPRetry, &v1.OTLPTraceRetrySpec{})
	retryConfig := otlptracehttp.RetryConfig{
		Enabled:         retry.Enabled,
		InitialInterval: time.Duration(retry.InitialInterval) * time.Second,
		MaxInterval:     time.Duration(retry.MaxInterval) * time.Second,
		MaxElapsedTime:  time.Duration(retry.MaxElapsedTime) * time.Second,
	}
	opts = appendOption(spec.OTLPRetry != nil, opts, otlptracehttp.WithRetry(retryConfig))

	// otlptracehttp.New doesn't return errors.
	// Reference: https://github.com/open-telemetry/opentelemetry-go/blob/v1.27.0/exporters/otlp/otlptrace/otlptracehttp/client.go#L101-L110
	return otlptracehttp.New(context.Background(), opts...)
}

func newGRPCExporter(spec *v1.GRPCTraceExporterSpec) (*otlptrace.Exporter, error) {
	var opts []otlptracegrpc.Option
	opts = appendOption(spec.Compress, opts, otlptracegrpc.WithCompressor("gzip"))
	opts = append(opts, otlptracegrpc.WithHeaders(spec.Headers))
	opts = appendOption(spec.EndpointURL != "", opts, otlptracegrpc.WithEndpointURL(spec.EndpointURL))
	opts = appendOption(spec.Timeout > 0, opts, otlptracegrpc.WithTimeout(time.Duration(spec.Timeout)*time.Second))
	opts = appendOption(spec.ReconnectionPeriod > 0, opts, otlptracegrpc.WithReconnectionPeriod(time.Duration(spec.ReconnectionPeriod)*time.Second))
	opts = appendOption(spec.ServiceConfig != "", opts, otlptracegrpc.WithServiceConfig(spec.ServiceConfig))

	if spec.Insecure {
		opts = append(opts, otlptracegrpc.WithInsecure())
	}

	tlsConfig, err := network.TLSConfig(spec.TLSConfig)
	if err != nil {
		return nil, err
	}
	opts = appendOption(tlsConfig != nil, opts, otlptracegrpc.WithTLSCredentials(credentials.NewTLS(tlsConfig)))

	retry := cmp.Or(spec.OTLPRetry, &v1.OTLPTraceRetrySpec{})
	retryConfig := otlptracegrpc.RetryConfig{
		Enabled:         retry.Enabled,
		InitialInterval: time.Duration(retry.InitialInterval) * time.Second,
		MaxInterval:     time.Duration(retry.MaxInterval) * time.Second,
		MaxElapsedTime:  time.Duration(retry.MaxElapsedTime) * time.Second,
	}
	opts = appendOption(spec.OTLPRetry != nil, opts, otlptracegrpc.WithRetry(retryConfig))

	return otlptracegrpc.New(context.Background(), opts...)
}

func newStdoutExporter(spec *v1.StdoutTraceExporterSpec) *stdouttrace.Exporter {
	var opts []stdouttrace.Option
	opts = appendOption(spec.PrettyPrint, opts, stdouttrace.WithPrettyPrint())
	opts = appendOption(spec.WithoutTimestamps, opts, stdouttrace.WithoutTimestamps())
	exporter, _ := stdouttrace.New(opts...) // No error returned.
	return exporter
}

func newZipkinExporter(spec *v1.ZipkinTraceExporterSpec) (*zipkin.Exporter, error) {
	var opts []zipkin.Option
	opts = append(opts, zipkin.WithHeaders(spec.Headers))
	// If EndpointURL is empty, default "http://localhost:9411/api/v2/spans" is used.
	return zipkin.New(spec.EndpointURL, opts...)
}

func appendOption[T any](shouldAppend bool, opts []T, opt T) []T {
	if shouldAppend {
		return append(opts, opt)
	}
	return opts
}
