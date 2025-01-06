package api

const (
	ErrPkg = "api"

	ErrTypeUtil      = "util"
	ErrTypeContainer = "container"
	ErrTypeExt       = "extension"
	ErrTypeFactory   = "factory"

	// ErrDscNil is a error description.
	// This description indicates unexpected
	// nil input.
	ErrDscNil = "nil value was given."

	// ErrDscNoAPI is a error description.
	// This description indicates there
	// is no requested API registered.
	ErrDscNoAPI = "api is not registered."

	// ErrDscDuplicateKey is a error description.
	// This description indicates there is
	// the same named APIs registered to the API handler.
	ErrDscDuplicateKey = "key duplication error."

	// ErrDscNoMethod is a error description.
	// This description indicates there is
	// no requested method implemented.
	ErrDscNoMethod = "method not implemented."

	// ErrDscNoManifest is a error description.
	// This description indicates no manifest
	// found which is requested.
	ErrDscNoManifest = "manifest not found."

	// ErrDscFormatSupport is a error description.
	// This description indicates the specified
	// format is not supported.
	ErrDscFormatSupport = "unsupported format."

	// ErrDscAssert is a error description.
	// This description indicates the failure of
	// type assertion.
	ErrDscAssert = "type assertion failed."

	// ErrDscProtoValidate is a error description.
	// This description indicates the failure of
	// validating proto message by protovalidate.
	ErrDscProtoValidate = "validating proto message failed."
)
