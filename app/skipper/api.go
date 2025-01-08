package skipper

import (
	"strings"

	v1 "github.com/aileron-gateway/aileron-gateway/apis/app/v1"
	"github.com/aileron-gateway/aileron-gateway/apis/kernel"
	"github.com/aileron-gateway/aileron-gateway/core"
	"github.com/aileron-gateway/aileron-gateway/kernel/api"
	"github.com/aileron-gateway/aileron-gateway/kernel/log"
	"github.com/aileron-gateway/aileron-gateway/kernel/txtutil"
	utilhttp "github.com/aileron-gateway/aileron-gateway/util/http"
	"google.golang.org/protobuf/reflect/protoreflect"
)

const (
	apiVersion = "app/v1"
	kind       = "Skipper"
	Key        = apiVersion + "/" + kind
)

var Resource api.Resource = &API{
	BaseResource: &api.BaseResource{
		DefaultProto: &v1.Skipper{
			APIVersion: apiVersion,
			Kind:       kind,
			Metadata: &kernel.Metadata{
				Namespace: "default",
				Name:      "default",
			},
			Spec: &v1.SkipperSpec{},
		},
	},
}

type API struct {
	*api.BaseResource
}

func (*API) Create(a api.API[*api.Request, *api.Response], msg protoreflect.ProtoMessage) (any, error) {
	c := msg.(*v1.Skipper)

	ms, err := api.ReferTypedObjects[core.Middleware](a, c.Spec.Middleware...)
	if err != nil {
		return nil, core.ErrCoreGenCreateObject.WithStack(err, map[string]any{"kind": kind})
	}

	ts, err := api.ReferTypedObjects[core.Tripperware](a, c.Spec.Tripperware...)
	if err != nil {
		return nil, core.ErrCoreGenCreateObject.WithStack(err, map[string]any{"kind": kind})
	}

	skippers := make([]*skipper, 0, len(c.Spec.SkipConditions))
	for _, sc := range c.Spec.SkipConditions {
		if sc == nil || sc.Matcher == nil {
			continue
		}
		m, err := txtutil.NewStringMatcher(txtutil.MatchTypes[sc.Matcher.MatchType], sc.Matcher.Patterns...)
		if err != nil {
			return nil, err // Return err as-is.
		}
		skippers = append(skippers, &skipper{
			methods: utilhttp.Methods(sc.Methods),
			paths:   m,
		})
	}

	return &skippable{
		skippers: skippers,
		ms:       ms,
		ts:       ts,
		lg:       log.DefaultOr(c.Metadata.Logger),
		name:     strings.Join([]string{apiVersion, c.Kind, c.Metadata.Namespace, c.Metadata.Name}, "/"),
	}, nil
}
