// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package http

import (
	"context"
	"net/http"
)

type headerContext struct{}

var headerContextKey = headerContext{}

// ContextWithProxyHeader register http header that should be proxied.
// A new http.Header will be created if nil header was given
// by the second argument h.
// context.Background() will be used if nil ctx was given.
// Update the caller's context when non nil context was returned.
//
//	ctx := r.Context()
//	ctx := ContextWithProxyHeader(ctx, http.Header{"foo":[]string{"bar"}})
//	r = r.WithContext(newCtx)
func ContextWithProxyHeader(ctx context.Context, h http.Header) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	if h == nil {
		return ctx
	}
	if v := ctx.Value(headerContextKey); v != nil {
		hh := v.(http.Header)
		for key, values := range h {
			hh[key] = append(hh[key], values...)
		}
		return ctx
	}
	return context.WithValue(ctx, headerContextKey, h)
}

// ProxyHeaderFromContext returns http headers to proxy.
// nil will be returned if no header was found in the context
// or nil context was given.
func ProxyHeaderFromContext(ctx context.Context) http.Header {
	if ctx == nil {
		return nil
	}
	if v := ctx.Value(headerContextKey); v != nil {
		return v.(http.Header)
	}
	return nil
}

type preProxyHook struct{}

var preProxyHookContextKey = preProxyHook{}

// ContextWithPreProxyHook register pre proxy hook function
// to the context.
// context.Background() will be used if nil context was given.
//
//	ctx := r.Context()
//	ctx = ContextWithPreProxyHook(ctx, func(r *http.Request) error {....})
//	r = r.WithContext(newCtx)
func ContextWithPreProxyHook(ctx context.Context, h func(r *http.Request) error) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	if h == nil {
		return ctx
	}
	if v := ctx.Value(preProxyHookContextKey); v != nil {
		fs := v.(*[]func(r *http.Request) error)
		*fs = append(*fs, h)
		return ctx
	}
	fs := &[]func(r *http.Request) error{h}
	return context.WithValue(ctx, preProxyHookContextKey, fs)
}

// PreProxyHookFromContext returns pre proxy hook functions.
// Empty slice will be returned if no function was found in the context
// or nil context was given.
func PreProxyHookFromContext(ctx context.Context) []func(r *http.Request) error {
	if ctx == nil {
		return nil
	}
	if v := ctx.Value(preProxyHookContextKey); v != nil {
		return *(v.(*[]func(r *http.Request) error))
	}
	return nil
}
