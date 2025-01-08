package jaegertracer

import (
	"io"
	"regexp"
	"testing"

	v1 "github.com/aileron-gateway/aileron-gateway/apis/app/v1"
	k "github.com/aileron-gateway/aileron-gateway/apis/kernel"
	"github.com/aileron-gateway/aileron-gateway/core"
	"github.com/aileron-gateway/aileron-gateway/kernel/api"
	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go/config"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type mockOption struct {
	options []config.Option
}

func (m *mockOption) JaegerOptions() []config.Option {
	return m.options
}

func TestCreate(t *testing.T) {
	type condition struct {
		manifest protoreflect.ProtoMessage
	}

	type action struct {
		tracer     any
		err        any // error or errorutil.Kind
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
				tracer: &jaegerTracer{
					mNames: map[int]string{},
					tNames: map[int]string{},
				},
				err: nil,
			},
		),
		gen(
			"create with options",
			[]string{},
			[]string{},
			&condition{
				manifest: &v1.JaegerTracer{
					Metadata: &k.Metadata{},
					Spec: &v1.JaegerTracerSpec{
						ServiceName:         "jaeger",
						Sampler:             &v1.JaegerSamplerSpec{},
						Reporter:            &v1.JaegerReporterSpec{},
						Headers:             &v1.JaegerHeadersSpec{},
						BaggageRestrictions: &v1.JaegerBaggageRestrictionsSpec{},
						Throttler:           &v1.JaegerThrottlerSpec{},
					},
				},
			},
			&action{
				tracer: &jaegerTracer{
					mNames: map[int]string{},
					tNames: map[int]string{},
				},
				err: nil,
			},
		),
		gen(
			"create with span names",
			[]string{},
			[]string{},
			&condition{
				manifest: &v1.JaegerTracer{
					Metadata: &k.Metadata{},
					Spec: &v1.JaegerTracerSpec{
						ServiceName:          "jaeger",
						MiddlewareSpanNames:  map[int32]string{1: "foo"},
						TripperwareSpanNames: map[int32]string{2: "bar"},
					},
				},
			},
			&action{
				tracer: &jaegerTracer{
					mNames: map[int]string{1: "foo"},
					tNames: map[int]string{2: "bar"},
				},
				err: nil,
			},
		),
		gen(
			"input tags",
			[]string{},
			[]string{},
			&condition{
				manifest: &v1.JaegerTracer{
					Metadata: &k.Metadata{},
					Spec: &v1.JaegerTracerSpec{
						ServiceName: "jaeger",
						Tags: map[string]string{
							"testKey": "testValue",
						},
					},
				},
			},
			&action{
				tracer: &jaegerTracer{
					mNames: map[int]string{},
					tNames: map[int]string{},
				},
				err: nil,
			},
		),
		gen(
			"fail to process NewTracer",
			[]string{},
			[]string{},
			&condition{
				manifest: &v1.JaegerTracer{
					Metadata: &k.Metadata{},
					Spec: &v1.JaegerTracerSpec{
						ServiceName: "",
					},
				},
			},
			&action{
				err:        core.ErrCoreGenCreateObject,
				errPattern: regexp.MustCompile(core.ErrPrefix + `failed to create JaegerTracer`),
			},
		),
		gen(
			"input k8s attributes tags",
			[]string{},
			[]string{},
			&condition{
				manifest: &v1.JaegerTracer{
					Metadata: &k.Metadata{},
					Spec: &v1.JaegerTracerSpec{
						ServiceName: "jaeger",
						K8SAttributes: &v1.K8SAttributesSpec{
							ContainerName: "testContainer",
						},
					},
				},
			},
			&action{
				tracer: &jaegerTracer{
					mNames: map[int]string{},
					tNames: map[int]string{},
				},
				err: nil,
			},
		),
		gen(
			"input container attributes tags",
			[]string{},
			[]string{},
			&condition{
				manifest: &v1.JaegerTracer{
					Metadata: &k.Metadata{},
					Spec: &v1.JaegerTracerSpec{
						ServiceName: "jaeger",
						ContainerAttributes: &v1.ContainerAttributesSpec{
							ImageName: "testImage",
						},
					},
				},
			},
			&action{
				tracer: &jaegerTracer{
					mNames: map[int]string{},
					tNames: map[int]string{},
				},
				err: nil,
			},
		),
		gen(
			"input host attributes tags",
			[]string{},
			[]string{},
			&condition{
				manifest: &v1.JaegerTracer{
					Metadata: &k.Metadata{},
					Spec: &v1.JaegerTracerSpec{
						ServiceName: "jaeger",
						HostAttributes: &v1.HostAttributesSpec{
							ImageName: "testImage",
						},
					},
				},
			},
			&action{
				tracer: &jaegerTracer{
					mNames: map[int]string{},
					tNames: map[int]string{},
				},
				err: nil,
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {

			server := api.NewContainerAPI()
			tracer, err := Resource.Create(server, tt.C().manifest)
			testutil.DiffError(t, tt.A().err, tt.A().errPattern, err)

			opts := []cmp.Option{
				cmp.AllowUnexported(jaegerTracer{}),
				cmpopts.IgnoreInterfaces(struct{ opentracing.Tracer }{}),
				cmpopts.IgnoreInterfaces(struct{ io.Closer }{}),
			}
			testutil.Diff(t, tt.A().tracer, tracer, opts...)

		})
	}
}
