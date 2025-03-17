//go:build example

package example

import (
	"context"
	"errors"
	"testing"

	"github.com/aileron-gateway/aileron-gateway/apis/kernel"
	"github.com/aileron-gateway/aileron-gateway/cmd/aileron/app"
	"github.com/aileron-gateway/aileron-gateway/kernel/api"
	apps "github.com/aileron-gateway/aileron-gateway/util/register"
	core "github.com/aileron-gateway/aileron-gateway/util/register"
)

type Runner interface {
	Run(context.Context) error
}

func getEntrypointRunner(t *testing.T, env, config []string) Runner {

	t.Helper()

	svr := api.NewDefaultServeMux()
	f := api.NewFactoryAPI()
	core.RegisterAll(f)
	apps.RegisterAll(f)
	svr.Handle("core/", f)
	svr.Handle("app/", f)

	if err := app.LoadEnvFiles(env); err != nil {
		t.Error(err)
	}
	if err := app.LoadConfigFiles(svr, config); err != nil {
		t.Error(err)
	}

	req := &api.Request{
		Method: api.MethodGet,
		Key:    "core/v1/Entrypoint",
		Format: api.FormatProtoReference,
		Content: &kernel.Reference{
			APIVersion: "core/v1",
			Kind:       "Entrypoint",
			Namespace:  ".entrypoint",
			Name:       ".entrypoint",
		},
	}

	res, err := svr.Serve(context.Background(), req)
	if err != nil {
		t.Error(err)
	}
	entrypoint, ok := res.Content.(Runner)
	if !ok {
		t.Error(errors.New("failed to assert entrypoint to Runner interface"))
	}

	return entrypoint

}
