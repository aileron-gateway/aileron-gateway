package session

import (
	"context"
)

// sessionContext is the key type to saves session in the context.Context.
type sessionContext struct{}

// sessionContextKey is the key which saves session in the context.Context.
var sessionContextKey = &sessionContext{}

// SessionFromContext get session from the given context.
// nil will be returned when a session was not found in the context.
// This function panics when nil context was given as argument.
func SessionFromContext(c context.Context) Session {
	if v := c.Value(sessionContextKey); v != nil {
		return v.(Session)
	}
	return nil
}

// ContextWithSession returns a new context with given session.
// This function panics when nil context was given as argument.
// This function do nothing when the given session s is nil.
func ContextWithSession(ctx context.Context, s Session) context.Context {
	if s == nil {
		return ctx
	}
	return context.WithValue(ctx, sessionContextKey, s)
}
