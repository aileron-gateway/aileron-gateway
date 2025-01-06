package entrypoint

import (
	"cmp"
	"strings"

	v1 "github.com/aileron-gateway/aileron-gateway/apis/core/v1"
	"github.com/aileron-gateway/aileron-gateway/apis/kernel"
	"github.com/aileron-gateway/aileron-gateway/core"
	"github.com/aileron-gateway/aileron-gateway/kernel/api"
	"github.com/aileron-gateway/aileron-gateway/kernel/log"
	utilhttp "github.com/aileron-gateway/aileron-gateway/util/http"
	"google.golang.org/protobuf/reflect/protoreflect"
)

const (
	apiVersion = "core/v1"
	kind       = "Entrypoint"
	Key        = apiVersion + "/" + kind
)

var Resource api.Resource = &API{
	BaseResource: &api.BaseResource{
		DefaultProto: &v1.Entrypoint{
			APIVersion: apiVersion,
			Kind:       kind,
			Metadata: &kernel.Metadata{
				Namespace: ".entrypoint", // This should not be overwritten.
				Name:      ".entrypoint", // This should not be overwritten.
			},
			Spec: &v1.EntrypointSpec{},
		},
	},
}

type API struct {
	*api.BaseResource
}

func (*API) Mutate(msg protoreflect.ProtoMessage) protoreflect.ProtoMessage {
	c := msg.(*v1.Entrypoint)
	c.Metadata.Namespace = ".entrypoint" // Prevent users from overwriting this value.
	c.Metadata.Name = ".entrypoint"      // Prevent users from overwriting this value.
	return c
}

func (*API) Create(a api.API[*api.Request, *api.Response], msg protoreflect.ProtoMessage) (any, error) {
	c := msg.(*v1.Entrypoint)

	var initializers []core.Initializer
	var finalizers []core.Finalizer

	for _, ref := range c.Spec.Loggers {
		ref.Namespace = cmp.Or(ref.Namespace, "default")
		ref.Name = cmp.Or(ref.Name, "default")
		lg, err := api.ReferTypedObject[log.Logger](a, ref)
		if err != nil {
			return nil, core.ErrCoreGenCreateObject.WithStack(err, map[string]any{"kind": kind})
		}
		// Specified logger are registered with <apiVersion>/<kind>/<namespace>/<name>
		ref.Namespace = cmp.Or(ref.Namespace, "default")
		ref.Name = cmp.Or(ref.Name, "default")
		name := strings.Join([]string{ref.APIVersion, ref.Kind, ref.Namespace, ref.Name}, "/")
		log.SetGlobalLogger(name, lg)

		initializers = appendInitializer(initializers, lg)
		finalizers = appendFinalizer(finalizers, lg)
	}

	if c.Spec.DefaultLogger != nil {
		ref := c.Spec.DefaultLogger
		lg, err := api.ReferTypedObject[log.Logger](a, ref)
		if err != nil {
			return nil, core.ErrCoreGenCreateObject.WithStack(err, map[string]any{"kind": kind})
		}
		log.SetGlobalLogger(log.DefaultLoggerName, lg)

		// Specified logger are registered with <apiVersion>/<kind>/<namespace>/<name>
		ref.Namespace = cmp.Or(ref.Namespace, "default")
		ref.Name = cmp.Or(ref.Name, "default")
		name := strings.Join([]string{ref.APIVersion, ref.Kind, ref.Namespace, ref.Name}, "/")
		log.SetGlobalLogger(name, lg)

		initializers = appendInitializer(initializers, lg)
		finalizers = appendFinalizer(finalizers, log.GlobalLogger(log.DefaultLoggerName))
	}

	if c.Spec.DefaultErrorHandler != nil {
		eh, err := api.ReferTypedObject[core.ErrorHandler](a, c.Spec.DefaultErrorHandler)
		if err != nil {
			return nil, core.ErrCoreGenCreateObject.WithStack(err, map[string]any{"kind": kind})
		}
		utilhttp.SetGlobalErrorHandler(utilhttp.DefaultErrorHandlerName, eh)
		initializers = appendInitializer(initializers, eh)
		finalizers = appendFinalizer(finalizers, eh)
	}

	lg := log.DefaultOr(c.Metadata.Logger)

	izs, err := api.ReferTypedObjects[core.Initializer](a, c.Spec.Initializers...)
	if err != nil {
		return nil, core.ErrCoreGenCreateObject.WithStack(err, map[string]any{"kind": kind})
	}
	initializers = append(izs, initializers...)

	fzs, err := api.ReferTypedObjects[core.Finalizer](a, c.Spec.Finalizers...)
	if err != nil {
		return nil, core.ErrCoreGenCreateObject.WithStack(err, map[string]any{"kind": kind})
	}
	finalizers = append(fzs, finalizers...)

	runners, err := api.ReferTypedObjects[core.Runner](a, c.Spec.Runners...)
	if err != nil {
		return nil, core.ErrCoreGenCreateObject.WithStack(err, map[string]any{"kind": kind})
	}

	if c.Spec.WaitAll {
		return &waitGroup{
			lg:           lg,
			runners:      runners,
			initializers: initializers,
			finalizers:   finalizers,
		}, nil
	} else {
		return &channelGroup{
			lg:           lg,
			runners:      runners,
			initializers: initializers,
			finalizers:   finalizers,
		}, nil
	}
}

// appendInitializer append the target to the given initializer slice
// if the target implements core.Initialize interface.
// If target does not implement core.Initializer interface,
// this function do nothing and returns the given slice.
func appendInitializer(arr []core.Initializer, target any) []core.Initializer {
	if iz, ok := target.(core.Initializer); ok {
		arr = append(arr, iz)
	}
	return arr
}

// appendFinalizer append the target to the given finalizer slice
// if the target implements core.Finalizer interface.
// If target does not implement core.Finalizer interface,
// this function do nothing and returns the given slice.
func appendFinalizer(arr []core.Finalizer, target any) []core.Finalizer {
	if fz, ok := target.(core.Finalizer); ok {
		arr = append(arr, fz)
	}
	return arr
}
