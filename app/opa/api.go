// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package opa

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
	"github.com/aileron-gateway/aileron-gateway/internal/network"
	"github.com/aileron-gateway/aileron-gateway/kernel/api"
	"github.com/aileron-gateway/aileron-gateway/kernel/er"
	"github.com/aileron-gateway/aileron-gateway/kernel/log"
	utilhttp "github.com/aileron-gateway/aileron-gateway/util/http"
	"github.com/open-policy-agent/opa/ast"
	"github.com/open-policy-agent/opa/bundle"
	"github.com/open-policy-agent/opa/loader"
	"github.com/open-policy-agent/opa/rego"
	"github.com/open-policy-agent/opa/topdown"
	"google.golang.org/protobuf/reflect/protoreflect"
)

const (
	apiVersion = "app/v1"
	kind       = "OPAAuthzMiddleware"
	Key        = apiVersion + "/" + kind
)

var Resource api.Resource = &API{
	BaseResource: &api.BaseResource{
		DefaultProto: &v1.OPAAuthzMiddleware{
			APIVersion: apiVersion,
			Kind:       kind,
			Metadata: &kernel.Metadata{
				Namespace: "default",
				Name:      "default",
			},
			Spec: &v1.OPAAuthzMiddlewareSpec{
				ClaimsKey: "AuthnClaims",
			},
		},
	},
}

type API struct {
	*api.BaseResource
}

func (*API) Create(a api.API[*api.Request, *api.Response], msg protoreflect.ProtoMessage) (any, error) {
	c := msg.(*v1.OPAAuthzMiddleware)

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

	var w io.Writer
	w, _ = authLg.(io.Writer)
	w = cmp.Or(w, io.Writer(os.Stdout))

	queries := make([]*rego.PreparedEvalQuery, 0, len(c.Spec.Regos))
	for _, rg := range c.Spec.Regos {
		query, err := regoQueries(a, rg, rego.PrintHook(topdown.NewPrintHook(w)))
		if err != nil {
			return nil, core.ErrCoreGenCreateObject.WithStack(err, map[string]any{"kind": kind})
		}
		queries = append(queries, query)
	}

	return &authz{
		lg: authLg,
		eh: eh,

		queries: queries,
		key:     c.Spec.ClaimsKey,
		envData: envData(c.Spec.EnvData),
		trace:   c.Spec.EnableTrace,
	}, nil
}

// regoQueries parse rego queries from files.
func regoQueries(a api.API[*api.Request, *api.Response], spec *v1.RegoSpec, opts ...func(*rego.Rego)) (*rego.PreparedEvalQuery, error) {
	var roundTripper http.RoundTripper = network.DefaultHTTPTransport
	if spec.RoundTripper != nil {
		rt, err := api.ReferTypedObject[http.RoundTripper](a, spec.RoundTripper)
		if err != nil {
			return nil, err
		}
		roundTripper = rt
	}
	roundTripper = addHeader(spec.Header).Tripperware(roundTripper)

	opts = append(opts, []func(*rego.Rego){
		rego.Query(spec.QueryParameter),
		rego.EnablePrintStatements(spec.EnablePrintStatements),
		rego.ShallowInlining(spec.ShallowInlining),
		rego.Strict(spec.Strict),
		rego.StrictBuiltinErrors(spec.StrictBuiltinErrors),
		rego.Instrument(false),
	}...)

	// Load .rego files.
	for _, f := range spec.PolicyFiles {
		r, err := loader.RegoWithOpts(f, ast.ParserOptions{})
		if err != nil {
			return nil, err
		}
		opts = append(opts, rego.ParsedModule(r.Parsed))
	}

	// Load bundles.
	for _, p := range spec.BundlePaths {
		vc, err := verificationConfig(spec.BundleVerification)
		if err != nil {
			return nil, err
		}
		baseLoader := loader.NewFileLoader().
			WithFollowSymlinks(true).
			WithProcessAnnotation(true).
			WithSkipBundleVerification(spec.SkipBundleVerification).
			WithBundleVerificationConfig(vc)

		b, err := loadBundle(p, roundTripper, baseLoader)
		if err != nil {
			return nil, err
		}
		opts = append(opts, rego.ParsedBundle(p, b))
	}

	switch s := spec.Stores.(type) {
	case *v1.RegoSpec_FileStore:
		store, _ := newFileStore(s.FileStore)
		opts = append(opts, rego.Store(store))
	case *v1.RegoSpec_HTTPStore:
		store, _ := newHTTPStore(s.HTTPStore, roundTripper)
		opts = append(opts, rego.Store(store))
	}

	query, err := rego.New(opts...).PrepareForEval(context.TODO())
	if err != nil {
		return nil, err
	}
	return &query, nil
}

func envData(spec *v1.EnvDataSpec) map[string]any {
	if spec == nil {
		return nil
	}

	env := map[string]any{}

	if len(spec.Vars) > 0 {
		vars := make(map[string]string, len(spec.Vars))
		env["vars"] = make(map[string]string, len(spec.Vars))
		for _, k := range spec.Vars {
			vars[k] = os.Getenv(k)
		}
	}

	if spec.PID {
		env["pid"] = os.Getpid()
	}
	if spec.PPID {
		env["ppid"] = os.Getppid()
	}
	if spec.GID {
		env["gid"] = os.Getgid()
	}
	if spec.UID {
		env["uid"] = os.Getuid()
	}
	return env
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

func loadPolicy(path string, rt http.RoundTripper, loader loader.FileLoader) (*bundle.Bundle, error) {
	if strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://") {
		req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, path, nil)
		if err != nil {
			return nil, (&er.Error{
				Package:     "authz/opa",
				Type:        "load bundle",
				Description: "failed create request.",
				Detail:      path,
			}).Wrap(err)
		}

		resp, err := rt.RoundTrip(req)
		if err != nil {
			return nil, (&er.Error{
				Package:     "authz/opa",
				Type:        "load bundle",
				Description: "failed to get bundle from endpoint.",
				Detail:      path,
			}).Wrap(err)
		}
		defer func() {
			io.Copy(io.Discard, resp.Body)
			resp.Body.Close()
		}()

		if resp.StatusCode != http.StatusOK {
			return nil, (&er.Error{
				Package:     "authz/opa",
				Type:        "load bundle",
				Description: "bundle endpoint returned non 200 OK status.",
				Detail:      "got " + strconv.Itoa(resp.StatusCode),
			}).Wrap(err)
		}

		f, err := os.CreateTemp(os.TempDir(), "*.tar.gz")
		if err != nil {
			return nil, (&er.Error{
				Package:     "authz/opa",
				Type:        "load bundle",
				Description: "failed to read bundle from response body.",
				Detail:      path,
			}).Wrap(err)
		}

		// Remove temp tarball.
		defer os.Remove(f.Name())
		path = f.Name()

		_, err = f.ReadFrom(resp.Body)
		if err != nil {
			f.Close()
			return nil, (&er.Error{
				Package:     "authz/opa",
				Type:        "load bundle",
				Description: "failed to read bundle from response body.",
				Detail:      path,
			}).Wrap(err)
		}
		f.Close()
	}

	b, err := loader.AsBundle(path)
	if err != nil {
		return nil, (&er.Error{
			Package:     "authz/opa",
			Type:        "load bundle",
			Description: "failed to load bundle.",
			Detail:      path,
		}).Wrap(err)
	}

	return b, nil
}
