package jaegertracer

import (
	"cmp"
	"time"

	v1 "github.com/aileron-gateway/aileron-gateway/apis/app/v1"
	"github.com/aileron-gateway/aileron-gateway/apis/kernel"
	"github.com/aileron-gateway/aileron-gateway/core"
	"github.com/aileron-gateway/aileron-gateway/kernel/api"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
	"google.golang.org/protobuf/reflect/protoreflect"
)

const (
	apiVersion = "app/v1"
	kind       = "JaegerTracer"
	Key        = apiVersion + "/" + kind
)

var Resource api.Resource = &API{
	BaseResource: &api.BaseResource{
		DefaultProto: &v1.JaegerTracer{
			APIVersion: apiVersion,
			Kind:       kind,
			Metadata: &kernel.Metadata{
				Namespace: "default",
				Name:      "default",
			},
			Spec: &v1.JaegerTracerSpec{
				ServiceName: "gateway",
			},
		},
	},
}

// TestOptions is the list of jaeger tracer options.
// NOTE: This is testing purpose only.
var TestOptions []config.Option

type API struct {
	*api.BaseResource
}

func (*API) Create(a api.API[*api.Request, *api.Response], msg protoreflect.ProtoMessage) (any, error) {
	c := msg.(*v1.JaegerTracer)

	tags := make([]opentracing.Tag, 0, len(c.Spec.Tags))
	tags = append(tags, k8sAttributes(c.Spec.K8SAttributes)...)
	tags = append(tags, containerAttributes(c.Spec.ContainerAttributes)...)
	tags = append(tags, hostAttributes(c.Spec.HostAttributes)...)
	for k, v := range c.Spec.Tags {
		tags = append(tags, opentracing.Tag{
			Key:   k,
			Value: v,
		})
	}

	cfg := config.Configuration{
		ServiceName:         c.Spec.ServiceName,
		Disabled:            c.Spec.Disabled,
		Gen128Bit:           c.Spec.Gen128Bit,
		RPCMetrics:          c.Spec.RPCMetrics,
		Tags:                tags,
		Sampler:             newSamplerConfig(c.Spec.Sampler),
		Reporter:            newReporterConfig(c.Spec.Reporter),
		Headers:             newHeaderConfig(c.Spec.Headers),
		BaggageRestrictions: newBaggageRestrictionsConfig(c.Spec.BaggageRestrictions),
		Throttler:           newThrottlerConfig(c.Spec.Throttler),
	}

	tracer, closer, err := cfg.NewTracer(TestOptions...)
	if err != nil {
		return nil, core.ErrCoreGenCreateObject.WithStack(err, map[string]any{"kind": kind})
	}

	return &jaegerTracer{
		tracer: tracer,
		closer: closer,

		headers: c.Spec.HeaderNames,

		mNames: spanNames(c.Spec.MiddlewareSpanNames),
		tNames: spanNames(c.Spec.TripperwareSpanNames),
	}, nil
}

func spanNames(names map[int32]string) map[int]string {
	sn := make(map[int]string, len(names))
	for k, v := range names {
		sn[int(k)] = v
	}
	return sn
}

func newSamplerConfig(c *v1.JaegerSamplerSpec) *config.SamplerConfig {
	if c == nil {
		c = &v1.JaegerSamplerSpec{
			Type:  v1.JaegerSamplerSpec_Const,
			Param: 1.0,
		}
	}
	typ := map[v1.JaegerSamplerSpec_JaegerSamplerType]string{
		v1.JaegerSamplerSpec_Const:         jaeger.SamplerTypeConst,
		v1.JaegerSamplerSpec_Remote:        jaeger.SamplerTypeRemote,
		v1.JaegerSamplerSpec_Probabilistic: jaeger.SamplerTypeProbabilistic,
		v1.JaegerSamplerSpec_RateLimiting:  jaeger.SamplerTypeRateLimiting,
		v1.JaegerSamplerSpec_LowerBound:    jaeger.SamplerTypeLowerBound,
	}[c.Type]
	return &config.SamplerConfig{
		Type:                     cmp.Or(typ, jaeger.SamplerTypeConst),
		Param:                    c.Param,
		SamplingServerURL:        c.SamplingServerURL,
		SamplingRefreshInterval:  time.Millisecond * time.Duration(c.SamplingRefreshInterval),
		MaxOperations:            int(c.MaxOperations),
		OperationNameLateBinding: false,
	}
}

func newReporterConfig(c *v1.JaegerReporterSpec) *config.ReporterConfig {
	if c == nil {
		return nil
	}
	return &config.ReporterConfig{
		QueueSize:                  int(c.QueueSize),
		BufferFlushInterval:        time.Millisecond * time.Duration(c.BufferFlushInterval),
		LogSpans:                   c.LogSpans,
		LocalAgentHostPort:         c.LocalAgentHostPort,
		DisableAttemptReconnecting: c.DisableAttemptReconnecting,
		AttemptReconnectInterval:   time.Millisecond * time.Duration(c.AttemptReconnectInterval),
		CollectorEndpoint:          c.CollectorEndpoint,
		User:                       c.User,
		Password:                   c.Password,
		HTTPHeaders:                c.HTTPHeaders,
	}
}

func newHeaderConfig(c *v1.JaegerHeadersSpec) *jaeger.HeadersConfig {
	if c == nil {
		return nil
	}
	return &jaeger.HeadersConfig{
		JaegerDebugHeader:        c.JaegerDebugHeader,
		JaegerBaggageHeader:      c.JaegerBaggageHeader,
		TraceContextHeaderName:   c.TraceContextHeaderName,
		TraceBaggageHeaderPrefix: c.TraceBaggageHeaderPrefix,
	}
}

func newBaggageRestrictionsConfig(c *v1.JaegerBaggageRestrictionsSpec) *config.BaggageRestrictionsConfig {
	if c == nil {
		return nil
	}
	return &config.BaggageRestrictionsConfig{
		DenyBaggageOnInitializationFailure: c.DenyBaggageOnInitializationFailure,
		HostPort:                           c.HostPort,
		RefreshInterval:                    time.Millisecond * time.Duration(c.RefreshInterval),
	}
}

func newThrottlerConfig(c *v1.JaegerThrottlerSpec) *config.ThrottlerConfig {
	if c == nil {
		return nil
	}
	return &config.ThrottlerConfig{
		HostPort:                  c.HostPort,
		RefreshInterval:           time.Millisecond * time.Duration(c.RefreshInterval),
		SynchronousInitialization: c.SynchronousInitialization,
	}
}

func k8sAttributes(c *v1.K8SAttributesSpec) []opentracing.Tag {
	if c == nil {
		return nil
	}

	vals := []struct {
		key string
		val string
	}{
		{key: "k8s.cluster.name", val: c.ClusterName},
		{key: "k8s.container.name", val: c.ContainerName},
		{key: "k8s.container.restart_count", val: c.ContainerRestartCount},
		{key: "k8s.cronjob.name", val: c.CronJobName},
		{key: "k8s.cronjob.uid", val: c.CronJobUID},
		{key: "k8s.daemonset.name", val: c.DaemonSetName},
		{key: "k8s.daemonset.uid", val: c.DaemonSetUID},
		{key: "k8s.deployment.name", val: c.DeploymentName},
		{key: "k8s.deployment.uid", val: c.DeploymentUID},
		{key: "k8s.job.name", val: c.JobName},
		{key: "k8s.job.uid", val: c.JobUID},
		{key: "k8s.namespace.name", val: c.NamespaceName},
		{key: "k8s.node.name", val: c.NodeName},
		{key: "k8s.node.uid", val: c.NodeUID},
		{key: "k8s.pod.name", val: c.PodName},
		{key: "k8s.pod.uid", val: c.PodUID},
		{key: "k8s.replicaset.name", val: c.ReplicaSetName},
		{key: "k8s.replicaset.uid", val: c.ReplicaSetUID},
		{key: "k8s.statefulset.name", val: c.StatefulSetName},
		{key: "k8s.statefulset.uid", val: c.StatefulSetUID},
	}

	var kvs []opentracing.Tag
	for _, v := range vals {
		if v.val != "" {
			kvs = append(kvs, opentracing.Tag{
				Key:   v.key,
				Value: v.val,
			})
		}
	}

	return kvs
}

func containerAttributes(c *v1.ContainerAttributesSpec) []opentracing.Tag {
	if c == nil {
		return nil
	}

	vals := []struct {
		key string
		val string
	}{
		{key: "container.id", val: c.ID},
		{key: "container.image.name", val: c.ImageName},
		{key: "container.image.tag", val: c.ImageTag},
		{key: "container.name", val: c.Name},
		{key: "container.runtime", val: c.Runtime},
	}

	var kvs []opentracing.Tag
	for _, v := range vals {
		if v.val != "" {
			kvs = append(kvs, opentracing.Tag{
				Key:   v.key,
				Value: v.val,
			})
		}
	}

	return kvs
}

func hostAttributes(c *v1.HostAttributesSpec) []opentracing.Tag {
	if c == nil {
		return nil
	}

	vals := []struct {
		key string
		val string
	}{
		{key: "host.id", val: c.ID},
		{key: "host.image.id", val: c.ImageID},
		{key: "host.image.name", val: c.ImageName},
		{key: "host.image.version", val: c.ImageVersion},
		{key: "host.name", val: c.Name},
		{key: "host.type", val: c.Type},
	}

	var kvs []opentracing.Tag
	for _, v := range vals {
		if v.val != "" {
			kvs = append(kvs, opentracing.Tag{
				Key:   v.key,
				Value: v.val,
			})
		}
	}

	return kvs
}
