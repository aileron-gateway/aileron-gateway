package prommeter

import (
	v1 "github.com/aileron-gateway/aileron-gateway/apis/app/v1"
	"github.com/aileron-gateway/aileron-gateway/apis/kernel"
	"github.com/aileron-gateway/aileron-gateway/kernel/api"
	utilhttp "github.com/aileron-gateway/aileron-gateway/util/http"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/protobuf/reflect/protoreflect"
)

const (
	apiVersion = "app/v1"
	kind       = "PrometheusMeter"
	Key        = apiVersion + "/" + kind
)

var Resource api.Resource = &API{
	BaseResource: &api.BaseResource{
		DefaultProto: &v1.PrometheusMeter{
			APIVersion: apiVersion,
			Kind:       kind,
			Metadata: &kernel.Metadata{
				Namespace: "default",
				Name:      "default",
			},
			Spec: &v1.PrometheusMeterSpec{},
		},
	},
}

type API struct {
	*api.BaseResource
}

// Mutate changes configured values.
// The values of the msg which is given as the argument is the merged message of default values and user defined values.
// Changes for the fields of msg in this function make the final values which will be the input for validate and create function.
// Default values for "repeated" or "oneof" fields can also be applied in this function if necessary.
// Please check msg!=nil and asserting the mgs does not panic even they won't from the view of overall architecture of the gateway.
func (*API) Mutate(msg protoreflect.ProtoMessage) protoreflect.ProtoMessage {
	c := msg.(*v1.PrometheusMeter)

	if len(c.Spec.Patterns) == 0 {
		c.Spec.Patterns = append(c.Spec.Patterns, "/metrics")
	}

	return c
}

func (*API) Create(a api.API[*api.Request, *api.Response], msg protoreflect.ProtoMessage) (any, error) {
	c := msg.(*v1.PrometheusMeter)

	reg := prometheus.NewRegistry()
	collectors.NewProcessCollector(collectors.ProcessCollectorOpts{})
	reg.MustRegister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))
	reg.MustRegister(collectors.NewGoCollector())

	opts := promhttp.HandlerOpts{}
	handler := promhttp.InstrumentMetricHandler(reg, promhttp.HandlerFor(reg, opts))

	mAPICalls := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of received http requests",
		},
		[]string{"host", "path", "code", "method"},
	)
	reg.MustRegister(mAPICalls)

	tAPICalls := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_client_requests_total",
			Help: "Total number of sent http requests",
		},
		[]string{"host", "path", "code", "method"},
	)
	reg.MustRegister(tAPICalls)

	return &metrics{
		HandlerBase: &utilhttp.HandlerBase{
			AcceptPatterns: c.Spec.Patterns,
			AcceptMethods:  utilhttp.Methods(c.Spec.Methods),
		},
		Handler: handler,

		reg: reg,

		mAPICalls: mAPICalls,
		tAPICalls: tAPICalls,
	}, nil
}
