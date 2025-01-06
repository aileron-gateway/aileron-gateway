package httphandler

import (
	"fmt"
	"net/http"
	"path"
	"slices"

	v1 "github.com/aileron-gateway/aileron-gateway/apis/core/v1"
	k "github.com/aileron-gateway/aileron-gateway/apis/kernel"
	"github.com/aileron-gateway/aileron-gateway/core"
	"github.com/aileron-gateway/aileron-gateway/kernel/api"
	utilhttp "github.com/aileron-gateway/aileron-gateway/util/http"
	"google.golang.org/protobuf/reflect/protoreflect"
)

const (
	apiVersion = "core/v1"
	kind       = "HTTPHandler"
	Key        = apiVersion + "/" + kind
)

var Resource api.Resource = &API{
	BaseResource: &api.BaseResource{
		DefaultProto: &v1.HTTPHandler{
			APIVersion: apiVersion,
			Kind:       kind,
			Metadata: &k.Metadata{
				Namespace: "default",
				Name:      "default",
			},
			Spec: &v1.HTTPHandlerSpec{},
		},
	},
}

type API struct {
	*api.BaseResource
}

func (*API) Create(a api.API[*api.Request, *api.Response], msg protoreflect.ProtoMessage) (any, error) {
	c := msg.(*v1.HTTPHandler)

	ms, err := api.ReferTypedObjects[core.Middleware](a, c.Spec.Middleware...)
	if err != nil {
		return nil, core.ErrCoreGenCreateObject.WithStack(err, map[string]any{"kind": kind})
	}

	obj, err := api.ReferObject(a, c.Spec.Handler)
	if err != nil {
		return nil, core.ErrCoreGenCreateObject.WithStack(err, map[string]any{"kind": kind})
	}

	h, ok := obj.(http.Handler)
	if !ok {
		err := fmt.Errorf("fail to convert type from %T to http.Handler", obj)
		return nil, core.ErrCoreGenCreateObject.WithStack(err, map[string]any{"kind": kind})
	}

	handler := &handler{
		HandlerBase: &utilhttp.HandlerBase{},
		Handler:     utilhttp.MiddlewareChain(ms, h),
	}

	// Join the pattern and paths given from child handler
	// nad use it as patterns of this handler.
	if p, ok := obj.(interface{ Patterns() []string }); ok {
		paths := make([]string, len(p.Patterns()))
		for i, p := range p.Patterns() {
			paths[i] = path.Clean("/" + c.Spec.Pattern + p)
		}
		handler.HandlerBase.AcceptPatterns = paths
	}

	// Inherit methods from child handler.
	if m, ok := obj.(interface{ Methods() []string }); ok {
		ms := m.Methods()
		slices.Sort(ms)                                                     // slices.Compact requires sorted slice.
		handler.HandlerBase.AcceptMethods = slices.Clip(slices.Compact(ms)) // Remove duplicates.
	}

	return handler, nil
}

type handler struct {
	http.Handler
	// HandlerBase is the base struct for
	// http.Handler type resource.
	// This provides Patterns() and Methods() methods
	// to fulfill the core.Handler interface.
	*utilhttp.HandlerBase
}
