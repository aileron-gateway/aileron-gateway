package headercert

import (
	v1 "github.com/aileron-gateway/aileron-gateway/apis/app/v1"
	"github.com/aileron-gateway/aileron-gateway/apis/kernel"
	"github.com/aileron-gateway/aileron-gateway/core"
	"github.com/aileron-gateway/aileron-gateway/kernel/api"
	"github.com/aileron-gateway/aileron-gateway/kernel/log"
	utilhttp "github.com/aileron-gateway/aileron-gateway/util/http"
	"google.golang.org/protobuf/reflect/protoreflect"
)

const (
	apiVersion = "app/v1"
	kind       = "HeaderCertMiddleware"
	Key        = apiVersion + "/" + kind
)

var (
	Resource api.Resource = &API{
		BaseResource: &api.BaseResource{
			DefaultProto: &v1.HeaderCertMiddleware{
				APIVersion: apiVersion,
				Kind:       kind,
				Metadata: &kernel.Metadata{
					Namespace: "default",
					Name:      "default",
				},
				Spec: &v1.HeaderCertMiddlewareSpec{
					TLSConfig: &kernel.TLSConfig{},
				},
			},
		},
	}
)

type API struct {
	*api.BaseResource
}

func (*API) Create(a api.API[*api.Request, *api.Response], msg protoreflect.ProtoMessage) (any, error) {
	c := msg.(*v1.HeaderCertMiddleware)

	eh, err := utilhttp.ErrorHandler(a, c.Spec.ErrorHandler)
	if err != nil {
		return nil, core.ErrCoreGenCreateObject.WithoutStack(err, map[string]any{"kind": kind})
	}

	return &HeaderCert{
		lg:      log.GlobalLogger(log.DefaultLoggerName),
		eh:      eh,
		rootCAs: c.Spec.TLSConfig.RootCAs,
	}, nil
}
