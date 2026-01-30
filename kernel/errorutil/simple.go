package errorutil

import (
	"fmt"
	"strings"
)

func NewSimple(cause error, message, detail string, a ...any) error {
	return &SimpleError{
		Cause:   cause,
		Message: message,
		Detail:  fmt.Sprintf(detail, a...),
	}
}

// SimpleError is the simple general error object.
type SimpleError struct {
	// Cause is the cause of this error.
	Cause error
	// Message is the fixed error message.
	// Message is compared in the [SimpleError.Is].
	Message string
	// Detail is the error detail.
	// Detail is NOT compared in the [SimpleError.Is].
	Detail string
}

func (e *SimpleError) Unwrap() error {
	return e.Cause
}

func (e *SimpleError) Is(err error) bool {
	if err == nil {
		return e == nil
	}
	for {
		ee, ok := err.(*SimpleError)
		if ok {
			return e.Message == ee.Message
		}
		uw, ok := err.(interface{ Unwrap() error })
		if !ok {
			return false
		}
		err = uw.Unwrap()
	}
}

func (e *SimpleError) Error() string {
	var b strings.Builder
	b.Grow(len(e.Message) + len(e.Detail) + 1)
	_, _ = b.WriteString(e.Message)
	if e.Detail != "" {
		_, _ = b.WriteString(" ")
		_, _ = b.WriteString(e.Detail)
	}
	if e.Cause != nil {
		_, _ = b.WriteString(" [ ")
		_, _ = b.WriteString(e.Cause.Error())
		_, _ = b.WriteString(" ]")
	}
	return b.String()
}
