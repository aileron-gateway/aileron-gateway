// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package log

import (
	"context"
)

// Logger is an interface that all loggers should implement.
// Defined log levels are Debug, Info, Warn, Error.
type Logger interface {
	// Enable returns if the given log level is enabled for this logger.
	Enabled(LogLevel) bool

	// Debug prints a DEBUG level log.
	// This level should show the information that can be useful for debugging.
	// The second argument must be an array of key and value pairs.
	Debug(ctx context.Context, msg string, keyValues ...any)

	// Info prints a INFO level log.
	// This level should show the information that is not suit for DEBUG, WARN, ERROR levels.
	// The second argument must be an array of key and value pairs.
	Info(ctx context.Context, msg string, keyValues ...any)

	// Warn prints a WARN level log.
	// This level should show the information that should be treated
	// not immediately but may be in the future.
	// The second argument must be an array of key and value pairs.
	Warn(ctx context.Context, msg string, keyValues ...any)

	// Error prints a ERROR level log.
	// This level should show the information that should be treated immediately.
	// The second argument must be an array of key and value pairs.
	Error(ctx context.Context, msg string, keyValues ...any)
}

// Attributes is the interface of
// log attributes.
type Attributes interface {
	Name() string
	Map() map[string]any
	KeyValues() []any
}
