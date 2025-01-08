package app

import (
	"context"
	"net/http"
	"time"

	"github.com/aileron-gateway/aileron-gateway/kernel/kvs"
)

// AuthStatus is the status of authentication.
type AuthStatus uint

const (
	AuthNone    AuthStatus = 1 << iota // AuthNone has no meaning and effect
	AuthSuccess                        // AuthSuccess indicates successfully authenticated
	AuthFail                           // AuthFail indicates authn challenged but failed
	AuthSkip                           // AuthNone indicates authentication should be skipped
	AuthReturn                         // AuthReturn indicates the request should immediately returned to client
)

type AuthResult bool

const (
	// AuthFailed is the AuthResult that represent the requirement of next actions.
	AuthContinue AuthResult = false
	// AuthFailed is the AuthResult that represent failure authn/authz.
	AuthFailed AuthResult = false
	// AuthSucceeded is the AuthResult that represent successful authn/authz.
	AuthSucceeded AuthResult = true
)

// AuthenticationHandler is the interface of an authentication (AuthN) handler.
// Implement AuthorizationHandler for authorization (AuthZ).
type AuthenticationHandler interface {
	// ServeAuthn serves authentication.
	// The returned values are newRequest, authenticated, shouldReturn, error.
	// Caller must replace the original request with returned newRequest because it may contain new values in the request context.
	// Given request was successfully authenticated only when authenticated is true.
	// Caller must not call the successive middleware or handler but return immediately when shouldReturn is true.
	// Caller should pass the returned error to the error handler if it is not nil.
	// It is implementer's responsible to return an http error as error if any.
	ServeAuthn(http.ResponseWriter, *http.Request) (*http.Request, AuthResult, bool, error)
}

// SessionStore is the interface of a session store.
type SessionStore interface {
	kvs.Client[string, []byte]
	// SetEX set value with expiration.
	// This method can panic when nil context was given as the value for first argument.
	SetEx(context.Context, string, []byte, time.Duration) error
	// Expire sets expiration for the key.
	// This method can panic when nil context was given as the value for first argument.
	Expire(context.Context, string, time.Duration) (bool, error)
}

// HealthChecker is the interface for checking health status of resources.
type HealthChecker interface {
	// HealthCheck returns if this resource is healthy or not.
	// "true" must be returned as the second response when this resource is healthy.
	// The returned context must not be nil.
	// Return the given context as it is if there is no new context to return.
	HealthCheck(ctx context.Context) (context.Context, bool)
}

// Tracer is an interface of a tracer.
type Tracer interface {
	// Trace starts a new span.
	// The context must not be nil to take over the parent span.
	// Additional span attributes are given through the map of the second argument.
	// The returned function must be called by caller to end up the span.
	Trace(ctx context.Context, name string, attributes map[string]string) (context.Context, func())
}
