package errorutil

// ErrorKind is the interface which provides error code and kind
// of this error.
type ErrorKind interface {
	Code() string
	Kind() string
}

// Error is the interface of an error with stack traces.
type Error interface {
	error
	// Unwrap returns the internal error if any..
	Unwrap() error
	// StackTrace return the string of stacktrace.
	// This method return an empty string when no stacktrace available.
	StackTrace() string
}

// Attributes is the interface of error attributes
// which also satisfy log attributes.
type Attributes interface {
	Error
	Name() string
	Map() map[string]any
	KeyValues() []any
}

// Creator creates an error with/without stack traces.
type Creator interface {
	WithoutStack(error, map[string]any) Attributes
	WithStack(error, map[string]any) Attributes
}
