package txtutil

const (
	ErrPkg = "txtutil"

	ErrTypeTemplate = "template"
	ErrTypeMatcher  = "matcher"
	ErrTypeReplacer = "replacer"

	// ErrDscNil is a error description.
	// This description indicates a nil spec was given.
	ErrDscNil = "nil spec was given."

	// ErrDscUnsupported is a error description.
	// This description indicates the some unsupported
	// config was given.
	ErrDscUnsupported = "unsupported"

	// ErrDscPattern is a error description.
	// This description indicates the invalid patterns
	// of regular expression or any other string pattern.
	ErrDscPattern = "invalid pattern."

	// ErrDscTemplate is a error description.
	// This description indicates the invalid template.
	ErrDscTemplate = "invalid template."
)
