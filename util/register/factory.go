// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package register

import (
	"github.com/aileron-gateway/aileron-gateway/app/authn/authn"
	"github.com/aileron-gateway/aileron-gateway/app/authn/basic"
	"github.com/aileron-gateway/aileron-gateway/app/authn/digest"
	"github.com/aileron-gateway/aileron-gateway/app/authn/idkey"
	"github.com/aileron-gateway/aileron-gateway/app/authn/key"
	"github.com/aileron-gateway/aileron-gateway/app/authn/oauth"
	"github.com/aileron-gateway/aileron-gateway/app/authz/casbin"
	"github.com/aileron-gateway/aileron-gateway/app/authz/opa"
	"github.com/aileron-gateway/aileron-gateway/app/handler/echo"
	"github.com/aileron-gateway/aileron-gateway/app/handler/healthcheck"
	"github.com/aileron-gateway/aileron-gateway/app/meter/otelmeter"
	"github.com/aileron-gateway/aileron-gateway/app/meter/prommeter"
	"github.com/aileron-gateway/aileron-gateway/app/middleware/bodylimit"
	"github.com/aileron-gateway/aileron-gateway/app/middleware/compression"
	"github.com/aileron-gateway/aileron-gateway/app/middleware/cors"
	"github.com/aileron-gateway/aileron-gateway/app/middleware/csrf"
	"github.com/aileron-gateway/aileron-gateway/app/middleware/header"
	"github.com/aileron-gateway/aileron-gateway/app/middleware/headercert"
	"github.com/aileron-gateway/aileron-gateway/app/middleware/session"
	"github.com/aileron-gateway/aileron-gateway/app/middleware/soaprest"
	"github.com/aileron-gateway/aileron-gateway/app/middleware/throttle"
	"github.com/aileron-gateway/aileron-gateway/app/middleware/timeout"
	"github.com/aileron-gateway/aileron-gateway/app/middleware/tracking"
	"github.com/aileron-gateway/aileron-gateway/app/skipper"
	"github.com/aileron-gateway/aileron-gateway/app/storage/redis"
	"github.com/aileron-gateway/aileron-gateway/app/tracer/jaegertracer"
	"github.com/aileron-gateway/aileron-gateway/app/tracer/oteltracer"
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

	_ = r.Register(authn.Key, authn.Resource)
	_ = r.Register(basic.Key, basic.Resource)
	_ = r.Register(bodylimit.Key, bodylimit.Resource)
	_ = r.Register(casbin.Key, casbin.Resource)
	_ = r.Register(compression.Key, compression.Resource)
	_ = r.Register(cors.Key, cors.Resource)
	_ = r.Register(csrf.Key, csrf.Resource)
	_ = r.Register(digest.Key, digest.Resource)
	_ = r.Register(echo.Key, echo.Resource)
	_ = r.Register(header.Key, header.Resource)
	_ = r.Register(headercert.Key, headercert.Resource)
	_ = r.Register(healthcheck.Key, healthcheck.Resource)
	_ = r.Register(idkey.Key, idkey.Resource)
	_ = r.Register(jaegertracer.Key, jaegertracer.Resource)
	_ = r.Register(key.Key, key.Resource)
	_ = r.Register(oauth.Key, oauth.Resource)
	_ = r.Register(opa.Key, opa.Resource)
	_ = r.Register(otelmeter.Key, otelmeter.Resource)
	_ = r.Register(oteltracer.Key, oteltracer.Resource)
	_ = r.Register(prommeter.Key, prommeter.Resource)
	_ = r.Register(redis.Key, redis.Resource)
	_ = r.Register(session.Key, session.Resource)
	_ = r.Register(skipper.Key, skipper.Resource)
	_ = r.Register(soaprest.Key, soaprest.Resource)
	_ = r.Register(throttle.Key, throttle.Resource)
	_ = r.Register(timeout.Key, timeout.Resource)
	_ = r.Register(tracking.Key, tracking.Resource)
}
