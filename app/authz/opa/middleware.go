package opa

import (
	"io"
	"net/http"
	"net/url"

	"github.com/aileron-gateway/aileron-gateway/app"
	"github.com/aileron-gateway/aileron-gateway/core"
	"github.com/aileron-gateway/aileron-gateway/kernel/log"
	utilhttp "github.com/aileron-gateway/aileron-gateway/util/http"
	"github.com/open-policy-agent/opa/rego"
	"github.com/open-policy-agent/opa/topdown"
)

// authzClaims holds request information
// that is be used for authorization.
type authzClaims struct {
	Auth   any            `json:"auth" yaml:"auth" msgpack:"auth"`
	Host   string         `json:"host" yaml:"host" msgpack:"host"`
	Remote string         `json:"remote" yaml:"remote" msgpack:"remote"`
	Method string         `json:"method" yaml:"method" msgpack:"method"`
	API    string         `json:"api" yaml:"api" msgpack:"api"`
	Query  url.Values     `json:"query" yaml:"query" msgpack:"query"`
	Header http.Header    `json:"header" yaml:"header" msgpack:"header"`
	Env    map[string]any `json:"env" yaml:"env" msgpack:"env"`
}

// handler do authorization with Open Policy Agent Rego language.
// This implements app.AuthorizationHandler interface.
//
// References:
//   - https://www.openpolicyagent.org/docs/latest/
//   - https://www.openpolicyagent.org/docs/latest/policy-language/
//   - https://play.openpolicyagent.org/
type authz struct {
	lg log.Logger
	w  io.Writer

	eh core.ErrorHandler

	// query run authorization with rego.
	queries []*rego.PreparedEvalQuery

	// key is the context key to extract auth information.
	key string
	// envData is the host environmental data.
	// This includes environmental variables, pid, uid and so on.
	envData map[string]any

	// Trace enables tracing logs.
	// This should be false in production environment
	// because of the performance consideration.
	trace bool
}

func (m *authz) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims := &authzClaims{
			Auth:   r.Context().Value(m.key),
			Host:   r.Host,
			Remote: r.RemoteAddr,
			Method: r.Method,
			API:    r.URL.Path,
			Query:  r.URL.Query(),
			Header: r.Header,
			Env:    m.envData,
		}

		opts := []rego.EvalOption{rego.EvalInput(claims)}
		if m.trace {
			buf := topdown.NewBufferTracer()
			defer topdown.PrettyTraceWithLocation(m.w, *buf)
			opts = append(opts, rego.EvalEarlyExit(false), rego.EvalQueryTracer(buf))
		}

		for _, query := range m.queries {
			result, err := query.Eval(r.Context(), opts...)

			// Unauthorized. Try next policy.
			if err != nil {
				if m.lg.Enabled(log.LvDebug) {
					msg := "forbidden by OPA for method=" + r.Method + " path=" + r.URL.Path
					err := app.ErrAppAuthzAuthorization.WithoutStack(err, nil)
					m.lg.Debug(r.Context(), msg, err.Name(), err.Map())
				}
				continue
			}

			// Unauthorized. Try next policy.
			if !result.Allowed() {
				if m.lg.Enabled(log.LvDebug) {
					msg := "forbidden by OPA for method=" + r.Method + " path=" + r.URL.Path
					err := app.ErrAppAuthzAuthorization.WithoutStack(nil, nil)
					m.lg.Debug(r.Context(), msg, err.Name(), err.Map())
				}
				continue
			}

			// Authorized.
			next.ServeHTTP(w, r)
			return
		}

		// Forbidden.
		m.eh.ServeHTTPError(w, r, utilhttp.ErrForbidden)
	})
}
