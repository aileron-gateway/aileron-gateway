package common

import (
	"cmp"
	"context"

	"github.com/aileron-gateway/aileron-gateway/apis/kernel"
	"github.com/aileron-gateway/aileron-gateway/kernel/api"
	app "github.com/aileron-gateway/aileron-gateway/util/register"
	core "github.com/aileron-gateway/aileron-gateway/util/register"
)

func NewAPI() api.API[*api.Request, *api.Response] {
	f := api.NewFactoryAPI()
	core.RegisterAll(f)
	app.RegisterAll(f)

	server := api.NewDefaultServeMux()
	_ = server.Handle("core/", f)
	_ = server.Handle("app/", f)
	_ = server.Handle("container/", api.NewContainerAPI())

	return server
}

func PostTestResource(server api.API[*api.Request, *api.Response], ref *kernel.Reference, res any) {
	ref.Name = cmp.Or(ref.Name, "default")
	ref.Namespace = cmp.Or(ref.Namespace, "default")
	req := &api.Request{
		Method:  api.MethodPost,
		Key:     ref.APIVersion + "/" + ref.Kind + "/" + ref.Namespace + "/" + ref.Name,
		Content: res,
	}
	if _, err := server.Serve(context.Background(), req); err != nil {
		panic(err)
	}
}
