package casbin

import (
	"cmp"
	"context"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	v1 "github.com/aileron-gateway/aileron-gateway/apis/app/v1"
	"github.com/aileron-gateway/aileron-gateway/apis/kernel"
	"github.com/aileron-gateway/aileron-gateway/core"
	"github.com/aileron-gateway/aileron-gateway/kernel/api"
	"github.com/aileron-gateway/aileron-gateway/kernel/cron"
	"github.com/aileron-gateway/aileron-gateway/kernel/er"
	"github.com/aileron-gateway/aileron-gateway/kernel/log"
	"github.com/aileron-gateway/aileron-gateway/kernel/network"
	utilhttp "github.com/aileron-gateway/aileron-gateway/util/http"
	"github.com/casbin/casbin/v3"
	"github.com/casbin/casbin/v3/model"
	"github.com/casbin/casbin/v3/persist"
	"google.golang.org/protobuf/reflect/protoreflect"
)

const (
	apiVersion = "app/v1"
	kind       = "CasbinAuthzMiddleware"
	Key        = apiVersion + "/" + kind
)

var Resource api.Resource = &API{
	BaseResource: &api.BaseResource{
		DefaultProto: &v1.CasbinAuthzMiddleware{
			APIVersion: apiVersion,
			Kind:       kind,
			Metadata: &kernel.Metadata{
				Namespace: "default",
				Name:      "default",
			},
			Spec: &v1.CasbinAuthzMiddlewareSpec{
				ClaimsKey: "AuthnClaims",
			},
		},
	},
}

type API struct {
	*api.BaseResource
}

func (*API) Create(a api.API[*api.Request, *api.Response], msg protoreflect.ProtoMessage) (any, error) {
	c := msg.(*v1.CasbinAuthzMiddleware)

	lg := log.DefaultOr(c.Metadata.Logger)

	authLg := lg
	if c.Spec.Logger != nil {
		alg, err := api.ReferTypedObject[log.Logger](a, c.Spec.Logger)
		if err != nil {
			return nil, core.ErrCoreGenCreateObject.WithStack(err, map[string]any{"kind": kind})
		}
		authLg = alg
	}

	// Obtain an error handler.
	// Default error handler is returned when not configured.
	eh, err := utilhttp.ErrorHandler(a, c.Spec.ErrorHandler)
	if err != nil {
		return nil, core.ErrCoreGenCreateObject.WithStack(err, map[string]any{"kind": kind})
	}

	enforcers, err := newEnforcers(a, lg, c.Spec.Enforcers)
	if err != nil {
		return nil, core.ErrCoreGenCreateObject.WithStack(err, map[string]any{"kind": kind})
	}

	var w io.Writer
	w, _ = authLg.(io.Writer)
	w = cmp.Or(w, io.Writer(os.Stdout))

	return &authz{
		lg: authLg,
		w:  w,
		eh: eh,

		enforcers: enforcers,
		key:       c.Spec.ClaimsKey,
		extraKeys: c.Spec.ExtraKeys,

		explain: c.Spec.Explain,
	}, nil
}

func newEnforcers(a api.API[*api.Request, *api.Response], lg log.Logger, specs []*v1.EnforcerSpec) ([]casbin.IEnforcer, error) {
	// Create enforcer from models and policies.
	// See https://casbin.org/docs/supported-models for about models and policies.
	enforcers := make([]casbin.IEnforcer, 0, len(specs))
	for _, spec := range specs {
		var roundTripper http.RoundTripper = network.DefaultHTTPTransport
		if spec.RoundTripper != nil {
			rt, err := api.ReferTypedObject[http.RoundTripper](a, spec.RoundTripper)
			if err != nil {
				return nil, err
			}
			roundTripper = rt
		}
		roundTripper = addHeader(spec.Header).Tripperware(roundTripper)

		var adapter any
		var err error
		switch s := spec.Policies.(type) { // Policy is optional.
		case *v1.EnforcerSpec_PolicyPath:
			adapter, err = newFileAdapter(s.PolicyPath)
		case *v1.EnforcerSpec_PolicyURL:
			adapter = &httpAdapter{
				Adapter:  &noopAdapter{},
				endpoint: s.PolicyURL,
				rt:       roundTripper,
			}
		case *v1.EnforcerSpec_ExternalAdapter:
			adapter, err = api.ReferTypedObject[persist.Adapter](a, s.ExternalAdapter)
		}
		if err != nil {
			return nil, err
		}

		// Load model from file or HTTP endpoint.
		mdl, err := loadModel(spec.ModelPath, roundTripper)
		if err != nil {
			return nil, err
		}

		enf, err := casbin.NewEnforcer(mdl, adapter)
		if err != nil {
			return nil, err
		}

		enf.AddFunction("authValue", mapValue)
		enf.AddFunction("mapValue", mapValue)
		enf.AddFunction("containsString", contains[string])
		enf.AddFunction("containsInt", containsNumber[int])
		enf.AddFunction("asStringSlice", asSlice[string])
		enf.AddFunction("asIntSlice", asSliceNumber[int])
		enf.AddFunction("queryValue", queryValue)
		enf.AddFunction("queryValues", queryValues)
		enf.AddFunction("headerValue", headerValue)
		enf.AddFunction("headerValues", headerValues)

		if spec.Cron != "" {
			ct, err := cron.Parse(spec.Cron)
			if err != nil {
				return nil, err
			}
			job := ct.NewJob(func() {
				// This implementation is based on (casbin.Enforcer).LoadPolicy
				// The autoBuildRoleLinks is considered to be false.
				model := enf.GetModel()
				adapter := enf.GetAdapter()
				lg.Info(context.Background(), "casbin policy reloading")
				if err := adapter.LoadPolicy(model); err != nil {
					lg.Error(context.Background(), "casbin policy reload error: "+err.Error())
				}
				model.PrintPolicy()
			})
			go job.Start() // This job never stops.
		}

		enforcers = append(enforcers, enf)
	}

	return enforcers, nil
}

// addHeader is a tripperware that
// adds http headers to the requests.
type addHeader map[string]string

func (t addHeader) Tripperware(next http.RoundTripper) http.RoundTripper {
	if len(t) == 0 {
		return next
	}
	return core.RoundTripperFunc(func(r *http.Request) (w *http.Response, err error) {
		for k, v := range t {
			r.Header.Add(k, v)
		}
		return next.RoundTrip(r)
	})
}

func loadModel(path string, rt http.RoundTripper) (any, error) {
	if !strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://") {
		return path, nil
	}

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, path, nil)
	if err != nil {
		return nil, (&er.Error{
			Package:     "authz/casbin",
			Type:        "load model",
			Description: "failed create request.",
			Detail:      path,
		}).Wrap(err)
	}

	resp, err := rt.RoundTrip(req)
	if err != nil {
		return nil, (&er.Error{
			Package:     "authz/casbin",
			Type:        "load model",
			Description: "failed to get model from endpoint.",
			Detail:      path,
		}).Wrap(err)
	}
	defer func() {
		io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, (&er.Error{
			Package:     "authz/casbin",
			Type:        "load model",
			Description: "model endpoint returned non 200 OK status.",
			Detail:      "got " + strconv.Itoa(resp.StatusCode),
		}).Wrap(err)
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, (&er.Error{
			Package:     "authz/casbin",
			Type:        "load model",
			Description: "failed to read model from response body.",
			Detail:      path,
		}).Wrap(err)
	}

	m, err := model.NewModelFromString(string(b))
	if err != nil {
		return nil, (&er.Error{
			Package:     "authz/casbin",
			Type:        "load model",
			Description: "failed to load model.",
			Detail:      path,
		}).Wrap(err)
	}

	return m, nil
}
