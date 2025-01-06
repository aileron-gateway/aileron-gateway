package log

import (
	"context"
)

// attrsContext is the key type to
// saves attributes in the context.Context.
type attrsContext struct{}

// attrsContextKey is the key which
// saves log attributes in the context.Context.
var attrsContextKey = &attrsContext{}

// ContextWithAttrs saves the given log attributes
// into the given context.
// This function panics when a nil context was given.
func ContextWithAttrs(ctx context.Context, attr ...Attributes) context.Context {
	var attrs []Attributes
	if v := ctx.Value(attrsContextKey); v != nil {
		attrs = v.([]Attributes) // When v is not []Attributes, then a panic occurs.
	} else {
		attrs = make([]Attributes, len(attrs))
	}
	attrs = append(attrs, attr...)
	return context.WithValue(ctx, attrsContextKey, attrs)
}

// AttrsFromContext returns log attributes
// bounded to the given context.
// Nil slice will be returned when no attributes were found.
// This function panics when a nil context was given.
func AttrsFromContext(ctx context.Context) []Attributes {
	v := ctx.Value(attrsContextKey)
	if v == nil {
		return nil
	}
	return v.([]Attributes)
}
