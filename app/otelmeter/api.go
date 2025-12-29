// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package otelmeter

import (
	"cmp"
	"context"
	"io"
	"reflect"
	"time"

	v1 "github.com/aileron-gateway/aileron-gateway/apis/app/v1"
	"github.com/aileron-gateway/aileron-gateway/apis/kernel"
	"github.com/aileron-gateway/aileron-gateway/core"
	"github.com/aileron-gateway/aileron-gateway/internal/network"
	"github.com/aileron-gateway/aileron-gateway/kernel/api"
	"go.opentelemetry.io/contrib/instrumentation/runtime"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"google.golang.org/grpc/credentials"
	"google.golang.org/protobuf/reflect/protoreflect"
)

const (
	apiVersion = "app/v1"
	kind       = "OpenTelemetryMeter"
	Key        = apiVersion + "/" + kind
)

var (
	TestWriter io.Writer = nil
)

var Resource api.Resource = &API{
	BaseResource: &api.BaseResource{},
}

type API struct {
	*api.BaseResource
}

// Default returns a new configuration with default values.
// Returned configuration does not contain all default values.
// For example, default values for "repeated" or "oneof" fields in proto are set in mutate function.
// Use "--template" command option to show configuration with all default values.
func (o *API) Default() protoreflect.ProtoMessage {
	return &v1.OpenTelemetryMeter{
		APIVersion: apiVersion,
		Kind:       kind,
		Metadata: &kernel.Metadata{
			Namespace: "default",
			Name:      "default",
		},
		Spec: &v1.OpenTelemetryMeterSpec{
			ServiceName: "gateway",
			LibraryName: reflect.TypeOf(*o).PkgPath(),
			Exporters: &v1.OpenTelemetryMeterSpec_GRPCExporterSpec{
				GRPCExporterSpec: &v1.GRPCMetricsExporterSpec{},
			},
			PeriodicReader: &v1.PeriodicReaderSpec{
				Interval: 5,
				Timeout:  30,
			},
		},
	}
}

func (*API) Create(a api.API[*api.Request, *api.Response], msg protoreflect.ProtoMessage) (any, error) {
	c := msg.(*v1.OpenTelemetryMeter)

	var exporter sdkmetric.Exporter
	var err error
	switch c.Spec.Exporters.(type) {
	case *v1.OpenTelemetryMeterSpec_HTTPExporterSpec:
		exporter, err = newHTTPExporter(c.Spec.GetHTTPExporterSpec())
	case *v1.OpenTelemetryMeterSpec_GRPCExporterSpec:
		exporter, err = newGRPCExporter(c.Spec.GetGRPCExporterSpec())
	case *v1.OpenTelemetryMeterSpec_StdoutExporterSpec:
		exporter = newStdoutExporter(c.Spec.GetStdoutExporterSpec())
	}
	if err != nil {
		return nil, core.ErrCoreGenCreateObject.WithStack(err, map[string]any{"kind": kind})
	}

	// Create meterprovider using resource and reader options.
	mp := sdkmetric.NewMeterProvider(
		sdkmetric.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(c.Spec.ServiceName),
		)),
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(
			exporter,
			sdkmetric.WithInterval(time.Duration(c.Spec.PeriodicReader.Interval)*time.Second),
			sdkmetric.WithTimeout(time.Duration(c.Spec.PeriodicReader.Timeout)*time.Second),
		)),
	)

	// Create custom meter using meterprovider.
	meter := mp.Meter(c.Spec.LibraryName)

	// Start collecting runtime metrics.
	runtime.Start(
		runtime.WithMeterProvider(mp),
		runtime.WithMinimumReadMemStatsInterval(time.Second),
	)

	mAPICall, _ := meter.Int64Counter(
		"http_requests_total",
		metric.WithDescription("Total number of received http requests"),
	)

	tAPICall, _ := meter.Int64Counter(
		"http_client_requests_total",
		metric.WithDescription("Total number of sent http requests"),
	)

	return &otelMeter{
		mp:        mp,
		mAPICalls: mAPICall,
		tAPICalls: tAPICall,
	}, nil
}

func newHTTPExporter(spec *v1.HTTPMetricsExporterSpec) (*otlpmetrichttp.Exporter, error) {
	var opts []otlpmetrichttp.Option
	opts = appendOption(spec.EndpointURL != "", opts, otlpmetrichttp.WithEndpointURL(spec.EndpointURL))
	opts = append(opts, otlpmetrichttp.WithHeaders(spec.Headers))
	opts = appendOption(spec.Compress, opts, otlpmetrichttp.WithCompression(otlpmetrichttp.GzipCompression))
	opts = appendOption(spec.Insecure, opts, otlpmetrichttp.WithInsecure())
	opts = appendOption(spec.Timeout > 0, opts, otlpmetrichttp.WithTimeout(time.Duration(spec.Timeout)*time.Second))

	tlsConfig, err := network.TLSConfig(spec.TLSConfig)
	if err != nil {
		return nil, err
	}
	opts = appendOption(tlsConfig != nil, opts, otlpmetrichttp.WithTLSClientConfig(tlsConfig))

	retry := cmp.Or(spec.OTLPRetry, &v1.OTLPMetricsRetrySpec{})
	retryConfig := otlpmetrichttp.RetryConfig{
		Enabled:         retry.Enabled,
		InitialInterval: time.Duration(retry.InitialInterval) * time.Second,
		MaxInterval:     time.Duration(retry.MaxInterval) * time.Second,
		MaxElapsedTime:  time.Duration(retry.MaxElapsedTime) * time.Second,
	}
	opts = appendOption(spec.OTLPRetry != nil, opts, otlpmetrichttp.WithRetry(retryConfig))
	return otlpmetrichttp.New(context.Background(), opts...)
}

func newGRPCExporter(spec *v1.GRPCMetricsExporterSpec) (*otlpmetricgrpc.Exporter, error) {
	var opts []otlpmetricgrpc.Option
	opts = appendOption(spec.EndpointURL != "", opts, otlpmetricgrpc.WithEndpointURL(spec.EndpointURL))
	opts = append(opts, otlpmetricgrpc.WithHeaders(spec.Headers))
	opts = appendOption(spec.Compress, opts, otlpmetricgrpc.WithCompressor("gzip"))
	opts = appendOption(spec.Timeout > 0, opts, otlpmetricgrpc.WithTimeout(time.Duration(spec.Timeout)*time.Second))
	opts = appendOption(spec.ReconnectionPeriod > 0, opts, otlpmetricgrpc.WithReconnectionPeriod(time.Duration(spec.ReconnectionPeriod)*time.Second))
	opts = appendOption(spec.ServiceConfig != "", opts, otlpmetricgrpc.WithServiceConfig(spec.ServiceConfig))
	opts = appendOption(spec.Insecure, opts, otlpmetricgrpc.WithInsecure())

	tlsConfig, err := network.TLSConfig(spec.TLSConfig)
	if err != nil {
		return nil, err
	}
	opts = appendOption(tlsConfig != nil, opts, otlpmetricgrpc.WithTLSCredentials(credentials.NewTLS(tlsConfig)))

	retry := cmp.Or(spec.OTLPRetry, &v1.OTLPMetricsRetrySpec{})
	retryConfig := otlpmetricgrpc.RetryConfig{
		Enabled:         retry.Enabled,
		InitialInterval: time.Duration(retry.InitialInterval) * time.Second,
		MaxInterval:     time.Duration(retry.MaxInterval) * time.Second,
		MaxElapsedTime:  time.Duration(retry.MaxElapsedTime) * time.Second,
	}
	opts = appendOption(spec.OTLPRetry != nil, opts, otlpmetricgrpc.WithRetry(retryConfig))
	return otlpmetricgrpc.New(context.Background(), opts...)
}

func newStdoutExporter(spec *v1.StdoutMetricsExporterSpec) sdkmetric.Exporter {
	var opts []stdoutmetric.Option
	opts = appendOption(spec.PrettyPrint, opts, stdoutmetric.WithPrettyPrint())
	opts = appendOption(spec.WithoutTimestamps, opts, stdoutmetric.WithoutTimestamps())
	opts = appendOption(TestWriter != nil, opts, stdoutmetric.WithWriter(TestWriter))
	exporter, _ := stdoutmetric.New(opts...)
	return exporter
}

func appendOption[T any](shouldAppend bool, opts []T, opt T) []T {
	if shouldAppend {
		return append(opts, opt)
	}
	return opts
}
