package http

const (
	ErrPkg = "utilhttp"

	ErrTypeMime       = "mime"
	ErrTypeErrHandler = "error handler"
	ErrTypeChain      = "chain"

	// ErrDscIO is a error description.
	// This description indicates io error.
	// Especially for file io error
	ErrDscIO = "io error."

	// ErrDscParseMime is a error description.
	// This description indicates the failure of
	// parsing media type string.
	ErrDscParseMime = "failed to parse media type."

	// ErrDscRegexp is a error description.
	// This description indicates the failure of
	// parsing regular expression.
	ErrDscRegexp = "failed to parse regular expression."

	// ErrDscAssert is a error description.
	// This description indicates the failure of
	// type assertion.
	ErrDscAssert = "type assertion failed."
)
