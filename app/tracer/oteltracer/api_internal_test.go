package oteltracer

import (
	"context"
	"regexp"
	"testing"

	v1 "github.com/aileron-gateway/aileron-gateway/apis/app/v1"
	k "github.com/aileron-gateway/aileron-gateway/apis/kernel"
	"github.com/aileron-gateway/aileron-gateway/core"
	"github.com/aileron-gateway/aileron-gateway/kernel/api"
	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"go.opentelemetry.io/contrib/propagators/aws/xray"
	"go.opentelemetry.io/contrib/propagators/b3"
	"go.opentelemetry.io/contrib/propagators/jaeger"
	"go.opentelemetry.io/contrib/propagators/opencensus"
	"go.opentelemetry.io/contrib/propagators/ot"
	"go.opentelemetry.io/otel/propagation"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func TestCreate(t *testing.T) {
	type condition struct {
		manifest protoreflect.ProtoMessage
	}

	type action struct {
		tracer     any
		err        any
		errPattern *regexp.Regexp
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"create with default manifest",
			[]string{},
			[]string{},
			&condition{
				manifest: Resource.Default(),
			},
			&action{
				tracer: &otelTracer{},
				err:    nil,
			},
		),
		gen(
			"create with service and library names",
			[]string{},
			[]string{},
			&condition{
				manifest: &v1.OpenTelemetryTracer{
					Spec: &v1.OpenTelemetryTracerSpec{
						ServiceName: "test",
						LibraryName: "test",
						Exporters: &v1.OpenTelemetryTracerSpec_HTTPExporterSpec{
							HTTPExporterSpec: &v1.HTTPTraceExporterSpec{},
						},
						TracerProviderBatch: &v1.TracerProviderBatchSpec{},
						TracerProviderLimit: &v1.TracerProviderLimitSpec{},
					},
				},
			},
			&action{
				tracer: &otelTracer{},
				err:    nil,
			},
		),
		gen(
			"create with headers",
			[]string{},
			[]string{},
			&condition{
				manifest: &v1.OpenTelemetryTracer{
					Spec: &v1.OpenTelemetryTracerSpec{
						Headers: []string{"foo", "bar"},
						Exporters: &v1.OpenTelemetryTracerSpec_HTTPExporterSpec{
							HTTPExporterSpec: &v1.HTTPTraceExporterSpec{},
						},
						TracerProviderBatch: &v1.TracerProviderBatchSpec{},
						TracerProviderLimit: &v1.TracerProviderLimitSpec{},
					},
				},
			},
			&action{
				tracer: &otelTracer{
					headers: []string{"foo", "bar"},
				},
				err: nil,
			},
		),
		gen(
			"create with limit options",
			[]string{},
			[]string{},
			&condition{
				manifest: &v1.OpenTelemetryTracer{
					Spec: &v1.OpenTelemetryTracerSpec{
						Exporters: &v1.OpenTelemetryTracerSpec_HTTPExporterSpec{
							HTTPExporterSpec: &v1.HTTPTraceExporterSpec{},
						},
						TracerProviderBatch: &v1.TracerProviderBatchSpec{},
						TracerProviderLimit: &v1.TracerProviderLimitSpec{
							AttributeValueLengthLimit:   64,
							AttributeCountLimit:         64,
							EventCountLimit:             64,
							LinkCountLimit:              64,
							AttributePerEventCountLimit: 64,
							AttributePerLinkCountLimit:  64,
						},
					},
				},
			},
			&action{
				tracer: &otelTracer{},
				err:    nil,
			},
		),
		gen(
			"create with batch options",
			[]string{},
			[]string{},
			&condition{
				manifest: &v1.OpenTelemetryTracer{
					Spec: &v1.OpenTelemetryTracerSpec{
						Exporters: &v1.OpenTelemetryTracerSpec_HTTPExporterSpec{
							HTTPExporterSpec: &v1.HTTPTraceExporterSpec{},
						},
						TracerProviderBatch: &v1.TracerProviderBatchSpec{
							MaxQueueSize:       1024,
							BatchTimeout:       1,
							ExportTimeout:      2,
							MaxExportBatchSize: 128,
							Blocking:           true,
						},
						TracerProviderLimit: &v1.TracerProviderLimitSpec{},
					},
				},
			},
			&action{
				tracer: &otelTracer{},
				err:    nil,
			},
		),
		gen(
			"HTTPExporter",
			[]string{},
			[]string{},
			&condition{
				manifest: &v1.OpenTelemetryTracer{
					Spec: &v1.OpenTelemetryTracerSpec{
						Exporters: &v1.OpenTelemetryTracerSpec_HTTPExporterSpec{
							HTTPExporterSpec: &v1.HTTPTraceExporterSpec{},
						},
						TracerProviderBatch: &v1.TracerProviderBatchSpec{},
						TracerProviderLimit: &v1.TracerProviderLimitSpec{},
					},
				},
			},
			&action{
				tracer: &otelTracer{},
				err:    nil,
			},
		),
		gen(
			"HTTPExporter with invalid TLS",
			[]string{},
			[]string{},
			&condition{
				manifest: &v1.OpenTelemetryTracer{
					Spec: &v1.OpenTelemetryTracerSpec{
						Exporters: &v1.OpenTelemetryTracerSpec_HTTPExporterSpec{
							HTTPExporterSpec: &v1.HTTPTraceExporterSpec{
								TLSConfig: &k.TLSConfig{ClientAuth: k.ClientAuthType(999)},
							},
						},
						TracerProviderBatch: &v1.TracerProviderBatchSpec{},
						TracerProviderLimit: &v1.TracerProviderLimitSpec{},
					},
				},
			},
			&action{
				tracer:     nil,
				err:        core.ErrCoreGenCreateObject,
				errPattern: regexp.MustCompile(core.ErrPrefix + `failed to create OpenTelemetryTracer`),
			},
		),
		gen(
			"gRPCExporter",
			[]string{},
			[]string{},
			&condition{
				manifest: &v1.OpenTelemetryTracer{
					Spec: &v1.OpenTelemetryTracerSpec{
						Exporters: &v1.OpenTelemetryTracerSpec_GRPCExporterSpec{
							GRPCExporterSpec: &v1.GRPCTraceExporterSpec{},
						},
						TracerProviderBatch: &v1.TracerProviderBatchSpec{},
						TracerProviderLimit: &v1.TracerProviderLimitSpec{},
					},
				},
			},
			&action{
				tracer: &otelTracer{},
				err:    nil,
			},
		),
		gen(
			"gRPCExporter with Insecure option",
			[]string{},
			[]string{},
			&condition{
				manifest: &v1.OpenTelemetryTracer{
					Spec: &v1.OpenTelemetryTracerSpec{
						Exporters: &v1.OpenTelemetryTracerSpec_GRPCExporterSpec{
							GRPCExporterSpec: &v1.GRPCTraceExporterSpec{
								Insecure: true,
							},
						},
						TracerProviderBatch: &v1.TracerProviderBatchSpec{},
						TracerProviderLimit: &v1.TracerProviderLimitSpec{},
					},
				},
			},
			&action{
				tracer: &otelTracer{},
				err:    nil,
			},
		),
		gen(
			"gRPCExporter with invalid TLS",
			[]string{},
			[]string{},
			&condition{
				manifest: &v1.OpenTelemetryTracer{
					Spec: &v1.OpenTelemetryTracerSpec{
						Exporters: &v1.OpenTelemetryTracerSpec_GRPCExporterSpec{
							GRPCExporterSpec: &v1.GRPCTraceExporterSpec{
								TLSConfig: &k.TLSConfig{ClientAuth: k.ClientAuthType(999)},
							},
						},
						TracerProviderBatch: &v1.TracerProviderBatchSpec{},
						TracerProviderLimit: &v1.TracerProviderLimitSpec{},
					},
				},
			},
			&action{
				tracer:     nil,
				err:        core.ErrCoreGenCreateObject,
				errPattern: regexp.MustCompile(core.ErrPrefix + `failed to create OpenTelemetryTracer`),
			},
		),
		gen(
			"stdoutExporter",
			[]string{},
			[]string{},
			&condition{
				manifest: &v1.OpenTelemetryTracer{
					Spec: &v1.OpenTelemetryTracerSpec{
						Exporters: &v1.OpenTelemetryTracerSpec_StdoutExporterSpec{
							StdoutExporterSpec: &v1.StdoutTraceExporterSpec{
								PrettyPrint:       true,
								WithoutTimestamps: true,
							},
						},
						TracerProviderBatch: &v1.TracerProviderBatchSpec{},
						TracerProviderLimit: &v1.TracerProviderLimitSpec{},
					},
				},
			},
			&action{
				tracer: &otelTracer{},
				err:    nil,
			},
		),
		gen(
			"zipkinExporter",
			[]string{},
			[]string{},
			&condition{
				manifest: &v1.OpenTelemetryTracer{
					Spec: &v1.OpenTelemetryTracerSpec{
						Exporters: &v1.OpenTelemetryTracerSpec_ZipkinExporterSpec{
							ZipkinExporterSpec: &v1.ZipkinTraceExporterSpec{
								EndpointURL: "http://localhost:9411/api/v2/spans",
							},
						},
						TracerProviderBatch: &v1.TracerProviderBatchSpec{},
						TracerProviderLimit: &v1.TracerProviderLimitSpec{},
					},
				},
			},
			&action{
				tracer: &otelTracer{},
				err:    nil,
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			server := api.NewContainerAPI()
			postTestResource(server, nil)

			tracer, err := Resource.Create(server, tt.C().manifest)
			testutil.DiffError(t, tt.A().err, tt.A().errPattern, err)

			opts := []cmp.Option{
				cmp.AllowUnexported(otelTracer{}),
				cmpopts.IgnoreFields(otelTracer{}, "tracer", "tp", "pg"),
			}
			testutil.Diff(t, tt.A().tracer, tracer, opts...)
		})
	}
}

func postTestResource(server api.API[*api.Request, *api.Response], res any) {
	ref := &k.Reference{
		APIVersion: "container/v1",
		Kind:       "Container",
		Namespace:  "externalOptions",
		Name:       "externalOptions",
	}
	req := &api.Request{
		Method:  api.MethodPost,
		Key:     ref.APIVersion + "/" + ref.Kind + "/" + ref.Namespace + "/" + ref.Name,
		Content: res,
	}
	if _, err := server.Serve(context.Background(), req); err != nil {
		panic(err)
	}
}

func TestPropagators(t *testing.T) {
	type condition struct {
		types []v1.PropagationType
	}

	type action struct {
		props []propagation.TextMapPropagator
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"nil",
			[]string{},
			[]string{},
			&condition{
				types: nil,
			},
			&action{
				props: []propagation.TextMapPropagator{
					propagation.TraceContext{},
					propagation.Baggage{},
				},
			},
		),
		gen(
			"TraceContext",
			[]string{},
			[]string{},
			&condition{
				types: []v1.PropagationType{
					v1.PropagationType_W3CTraceContext,
				},
			},
			&action{
				props: []propagation.TextMapPropagator{
					propagation.TraceContext{},
				},
			},
		),
		gen(
			"Baggage",
			[]string{},
			[]string{},
			&condition{
				types: []v1.PropagationType{
					v1.PropagationType_W3CBaggage,
				},
			},
			&action{
				props: []propagation.TextMapPropagator{
					propagation.Baggage{},
				},
			},
		),
		gen(
			"B3",
			[]string{},
			[]string{},
			&condition{
				types: []v1.PropagationType{
					v1.PropagationType_B3,
				},
			},
			&action{
				props: []propagation.TextMapPropagator{
					b3.New(),
				},
			},
		),
		gen(
			"Jaeger",
			[]string{},
			[]string{},
			&condition{
				types: []v1.PropagationType{
					v1.PropagationType_Jaeger,
				},
			},
			&action{
				props: []propagation.TextMapPropagator{
					jaeger.Jaeger{},
				},
			},
		),
		gen(
			"XRay",
			[]string{},
			[]string{},
			&condition{
				types: []v1.PropagationType{
					v1.PropagationType_XRay,
				},
			},
			&action{
				props: []propagation.TextMapPropagator{
					xray.Propagator{},
				},
			},
		),
		gen(
			"OpenCensus",
			[]string{},
			[]string{},
			&condition{
				types: []v1.PropagationType{
					v1.PropagationType_OpenCensus,
				},
			},
			&action{
				props: []propagation.TextMapPropagator{
					opencensus.Binary{},
				},
			},
		),
		gen(
			"OpenTracing",
			[]string{},
			[]string{},
			&condition{
				types: []v1.PropagationType{
					v1.PropagationType_OpenTracing,
				},
			},
			&action{
				props: []propagation.TextMapPropagator{
					ot.OT{},
				},
			},
		),
		gen(
			"all",
			[]string{},
			[]string{},
			&condition{
				types: []v1.PropagationType{
					v1.PropagationType_W3CTraceContext,
					v1.PropagationType_W3CBaggage,
					v1.PropagationType_B3,
					v1.PropagationType_Jaeger,
					v1.PropagationType_XRay,
					v1.PropagationType_OpenCensus,
					v1.PropagationType_OpenTracing,
				},
			},
			&action{
				props: []propagation.TextMapPropagator{
					propagation.TraceContext{},
					propagation.Baggage{},
					b3.New(),
					jaeger.Jaeger{},
					xray.Propagator{},
					opencensus.Binary{},
					ot.OT{},
				},
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			props := propagators(tt.C().types)
			opts := []cmp.Option{
				cmp.AllowUnexported(b3.New()),
			}
			testutil.Diff(t, tt.A().props, props, opts...)
		})
	}
}
