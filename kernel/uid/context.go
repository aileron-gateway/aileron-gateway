// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package uid

import (
	"context"
)

// idContext is the key type to save an ID in a context.
type idContext struct{}

// idContextKey is the key to save an ID in a context.
var idContextKey = idContext{}

// ContextWithID save the given ID into the given context.
// context.Background() will be used when nil context was given.
// This function accepts the empty string as the id argument.
//
// example:
//
//	id := "EXAMPLE_REQUEST_ID"
//	ctx = uid.ContextWithID(r.Context(), id) // r is *http.Request
//	r = r.WithContext(ctx)
func ContextWithID(ctx context.Context, id string) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	return context.WithValue(ctx, idContextKey, id)
}

// IDFromContext returns an ID extracted from the given context.
// Empty string will be returned if no ID found in the context
// or the given context were nil.
func IDFromContext(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	v := ctx.Value(idContextKey)
	if v == nil {
		return ""
	}
	return v.(string)
}
