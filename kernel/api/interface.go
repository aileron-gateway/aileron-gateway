// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package api

import (
	"context"
)

// API is the interface of API component.
type API[Q, S any] interface {
	Serve(context.Context, Q) (S, error)
}

// ServeMux is the multiplexer interface for APIs.
type ServeMux[K comparable, Q, S any] interface {
	API[Q, S]
	Handle(K, API[Q, S]) error
}
