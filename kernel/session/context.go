// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package session

import (
	"context"
)

type ctxKey struct{ string }

var sessionCtxKey = &ctxKey{"Context value from kernel/session."}

// SessionFromContext get session from the given context.
// nil will be returned when a session was not found in the context.
// This function panics when nil context was given as argument.
func SessionFromContext(c context.Context) Session {
	if v := c.Value(sessionCtxKey); v != nil {
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
	return context.WithValue(ctx, sessionCtxKey, s)
}
