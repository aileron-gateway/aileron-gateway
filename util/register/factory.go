package register

import (
	"github.com/aileron-gateway/aileron-gateway/core/entrypoint"
	"github.com/aileron-gateway/aileron-gateway/core/errhandler"
	"github.com/aileron-gateway/aileron-gateway/core/goplugin"
	"github.com/aileron-gateway/aileron-gateway/core/httpclient"
	"github.com/aileron-gateway/aileron-gateway/core/httphandler"
	"github.com/aileron-gateway/aileron-gateway/core/httplogger"
	"github.com/aileron-gateway/aileron-gateway/core/httpproxy"
	"github.com/aileron-gateway/aileron-gateway/core/httpserver"
	"github.com/aileron-gateway/aileron-gateway/core/slogger"
	"github.com/aileron-gateway/aileron-gateway/core/static"
	"github.com/aileron-gateway/aileron-gateway/core/template"
	"github.com/aileron-gateway/aileron-gateway/kernel/api"
)

type Registerer interface {
	Register(string, api.Resource) error
}

func RegisterAll(r Registerer) {
	_ = r.Register(entrypoint.Key, entrypoint.Resource)
	_ = r.Register(errhandler.Key, errhandler.Resource)
	_ = r.Register(goplugin.Key, goplugin.Resource)
	_ = r.Register(httpclient.Key, httpclient.Resource)
	_ = r.Register(httphandler.Key, httphandler.Resource)
	_ = r.Register(httplogger.Key, httplogger.Resource)
	_ = r.Register(httpproxy.Key, httpproxy.Resource)
	_ = r.Register(httpserver.Key, httpserver.Resource)
	_ = r.Register(slogger.Key, slogger.Resource)
	_ = r.Register(static.Key, static.Resource)
	_ = r.Register(template.Key, template.Resource)
}
