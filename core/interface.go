// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package core

import (
	"context"
	"net/http"
)

// Runner is the interface to serve some application.
type Runner interface {
	// Run starts the runner.
	// Run can block the process.
	Run(context.Context) error
}

// Middleware is the interface of http middleware.
// Following code is the example of middleware that
// adds a http response header.
//
// ------------------------------
//
//	type foo struct {
//		value string
//	}
//
//	func (f *foo) Middleware(next http.Handler) http.Handler {
//		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//
//			w.Header().Add("foo", f.value)
//			next.ServeHTTP(w, r)
//
//		})
//	}
//
// ------------------------------
type Middleware interface {
	Middleware(http.Handler) http.Handler
}

// Tripperware is the interface of client middleware,
// or round tripper middleware.
// Following code is the example of tripperware that
// adds a http request header.
//
// ------------------------------
//
//	type foo struct {
//		value string
//	}
//
//	func (f *foo) Tripperware(next http.RoundTripper) http.RoundTripper {
//		return RoundTripperFunc(func(r *http.Request) (w *http.Response, err error) {
//
//			r.Header.Add("foo", f.value)
//			return next.RoundTrip(r)
//
//		})
//	}
//
// ------------------------------
type Tripperware interface {
	Tripperware(http.RoundTripper) http.RoundTripper
}

// RoundTripperFunc wraps a function with the signature of HTTP round tripper.
// This implements the http.RoundTripper interface to the function.
//
// ------------------------------
//
//	func addHeader(r *http.Request) (*http.Response, error) {
//		r.Header.Add("foo", "bar")
//		return http.DefaultTransport.RoundTrip(r)
//	}
//
//	roundTripper := RoundTripperFunc(addHeader)
//
// ------------------------------
type RoundTripperFunc func(*http.Request) (*http.Response, error)

func (f RoundTripperFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}

// Matcher is an interface of value matcher.
type Matcher[T any] interface {
	// Match returns true when the given value
	// was matched to this object.
	// It's implementer's responsible to provides
	// the matching algorithm.
	Match(T) bool
}

// HTTPError represents the HTTP error response.
type HTTPError interface {
	error
	// Unwrap returns the internal error if any.
	Unwrap() error
	// StatusCode returns the response status code of t.
	StatusCode() int
	// Content returns content type and body.
	Content(accept string) (string, []byte)
}

// ErrorHandler is an interface of a HTTP error handler
// that handle and returns HTTP error response to the clients.
type ErrorHandler interface {
	// ServeHTTPError serves HTTP error.
	// It's callers responsibility to make sure
	// the all given arguments are non nil.
	ServeHTTPError(http.ResponseWriter, *http.Request, error)
}

// Initializer initialize the resource.
type Initializer interface {
	Initialize() error
}

// Finalizer finalize the resource.
type Finalizer interface {
	Finalize() error
}

// CookieCreator creates a new instance of a HTTP cookie.
type CookieCreator interface {
	// NewCookie returns a new instance of a cookie.
	// Default attributes depends on the implementers.
	// Callers should change at least name and value fields.
	NewCookie() *http.Cookie
}
