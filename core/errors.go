package core

import (
	"github.com/aileron-gateway/aileron-gateway/kernel/errorutil"
)

var ErrPrefix = `^E[0-9]{4}.[a-zA-Z]+ `

var (
	// ErrPrimitive wraps any primitive error type.
	ErrPrimitive = errorutil.NewKind("E0000", "Undefined", "{{msg}}")

	// general: E2000 - E2049
	// TODO: remove ErrCoreGen*
	ErrCoreGenCreateComponent = errorutil.NewKind("E2000", "CoreGenCreateComponent", "failed to create component. {{reason}}")
	ErrCoreGenCreateObject    = errorutil.NewKind("E2004", "CoreGenCreateObject", "failed to create {{kind}}")

	// core/httplogger: E2050 - E2059
	ErrCoreLogger = errorutil.NewKind("E2050", "CoreLogger", "http logging error.")

	// core/entrypoint: E2060 - E2069
	ErrCoreEntrypointRun = errorutil.NewKind("E2060", "CoreEntrypointRun", "error on running entrypoint")

	// core/httpproxy: E2120 - E2129
	ErrCoreProxyUnavailable      = errorutil.NewKind("E2121", "CoreProxyUpstreamUnavailable", "upstream unavailable for {{path}}")
	ErrCoreProxyNoUpstream       = errorutil.NewKind("E2122", "CoreProxyNoUpstream", "cannot find upstream for {{path}}")
	ErrCoreProxyTimeout          = errorutil.NewKind("E2123", "CoreProxyTimeout", "request timeout")
	ErrCoreProxyRoundtrip        = errorutil.NewKind("E2124", "CoreProxyRoundtrip", "error occurred in the round trippers")
	ErrCoreProxyProtocolSwitch   = errorutil.NewKind("E2125", "CoreProxyProtocolSwitch", "failed to upgrade protocol. {{reason}}")
	ErrCoreProxyBidirectionalCom = errorutil.NewKind("E2126", "CoreProxyBidirectionalCom", "failed to bidirectional communication")
	ErrCoreProxyNoRecovery       = errorutil.NewKind("E2127", "CoreProxyNoRecovery", "un-recoverable error occurred in proxy. this message is logging only")

	// core/httpserver: E2130 - E2139
	ErrCoreServer         = errorutil.NewKind("E2130", "CoreServer", "error was returned from server")
	ErrCoreServerNotFound = errorutil.NewKind("E2131", "CoreServerNotFound", "handler not found for {{pattern}}")
	ErrCoreServerRecover  = errorutil.NewKind("E2132", "CoreServerRecover", "panic recovered")

	// core/static: E2140 - E2149
	ErrCoreStaticServer = errorutil.NewKind("E2140", "CoreStaticServer", "failed to serve static file. {{body}}")
)
