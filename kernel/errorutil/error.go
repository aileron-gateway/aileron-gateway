package errorutil

const (
	keyError = "error"
	keyCode  = "code"
	keyKind  = "kind"
	keyMsg   = "msg"
	keyStack = "stack"
)

// New creates a new error attribute.
func New(code, kind, msg string, stack []byte, err error) *ErrorAttrs {
	return &ErrorAttrs{
		code:  code,
		kind:  kind,
		stack: stack,
		name:  keyError,
		msg:   msg,
		err:   err,
	}
}

// ErrorAttrs is an error which can have identifiers code,kind and can hold stack traces.
// Use errorutil.New to create a new error.
// This implements Attributes interface.
type ErrorAttrs struct {
	code  string
	kind  string
	stack []byte
	name  string
	msg   string
	err   error
}

// Error returns an error string.
// This method implements error interface.
func (e *ErrorAttrs) Error() string {
	return e.msg
}

// StackTrace returns the string of stacktraces.
// An empty string will be returned when
// there is no available stacktraces.
func (e *ErrorAttrs) StackTrace() string {
	return string(e.stack)
}

// Code returns the error code.
func (e *ErrorAttrs) Code() string {
	return e.code
}

// Kind returns the error kind.
func (e *ErrorAttrs) Kind() string {
	return e.kind
}

func (e *ErrorAttrs) Name() string {
	return e.name
}

func (e *ErrorAttrs) Map() map[string]any {
	return map[string]any{
		keyCode:  e.code,
		keyKind:  e.kind,
		keyMsg:   e.msg,
		keyStack: e.StackTrace(),
	}
}

func (e *ErrorAttrs) KeyValues() []any {
	return []any{
		keyCode, e.code,
		keyKind, e.kind,
		keyMsg, e.msg,
		keyStack, e.StackTrace(),
	}
}

// Unwrap returns an internal error.
func (e *ErrorAttrs) Unwrap() error {
	return e.err
}

// Is checks if the given error is the same as this error.
// "same" means the errors have the same code and kind.
// This checks wrapped errors for the given error
// through interface{ Unwrap() error }.
func (e *ErrorAttrs) Is(err error) bool {
	target := err
	for target != nil {
		c, ok := target.(ErrorKind)
		if ok && e.code == c.Code() && e.kind == c.Kind() {
			return true
		}

		uw, ok := target.(interface{ Unwrap() error })
		if !ok {
			return false
		}

		target = uw.Unwrap()
	}

	return false
}
