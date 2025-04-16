// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package errorutil

import (
	"bytes"
	"runtime"
	"sync"
)

// pool is the pool for *bytes.Buffer.
// Use like below.
//
//	buf := pool.Get().(*bytes.Buffer)
//	defer pool.Put(buf)
//	buf.Reset()  // Make sure to reset before use.
var pool = &sync.Pool{
	New: func() any {
		var buf bytes.Buffer
		return &buf
	},
}

// NewKind creates a new error kind.
// If the error message template is invalid, it returns nil.
// Creating a new error kind is not recommended to use many times on runtime
// because parsing a template could affects the latency.
// If the given template is not valid, this function panics.
func NewKind(code, kind, tpl string) *Kind {
	return &Kind{
		code: code,
		kind: kind,
		tpl:  newTemplate(tpl),
	}
}

// Kind is an error kind.
// Error kind is not an actual error it self.
type Kind struct {
	code string
	kind string
	tpl  *template
}

// Code returns an error code of this kind.
func (k *Kind) Code() string {
	return k.code
}

// Kind returns an error kind of this kind.
func (k *Kind) Kind() string {
	return k.kind
}

// WithoutStack returns a new error without stack traces.
// An internal error should be given at the first argument if any.
// The map args given by the second argument is the input
// for the error message template. Args can also be nil.
// WithoutStack removes all cumulative stack traces
// which the given err had.
func (k *Kind) WithoutStack(err error, args map[string]any) Attributes {
	var stack []byte
	if e, ok := err.(*ErrorAttrs); ok {
		stack = e.stack
	}
	return k.newError(err, args, stack)
}

// WithStack returns a new error with stack traces.
// An internal error should be given at the first argument if any.
// The map args given by the second argument is the input
// for the error message template. Args can also be nil.
func (k *Kind) WithStack(err error, args map[string]any) Attributes {
	top := []byte(k.code + "." + k.kind + ":\n")
	stack := make([]byte, 1<<12) // Read max 4kiB stack traces. May be enough.
	copy(stack, top)
	n := runtime.Stack(stack[len(top):], false)
	stack = stack[:len(top)+n]
	if e, ok := err.(*ErrorAttrs); ok {
		stack = append(stack, '\n')
		stack = append(stack, e.stack...)
	}
	return k.newError(err, args, stack)
}

// newError creates a new *errorutil.Error with the same code/kind of this kind.
// The second argument args is passed to the message template.
// The internal err can be nil.
func (k *Kind) newError(err error, args map[string]any, stack []byte) *ErrorAttrs {
	buf := pool.Get().(*bytes.Buffer)
	defer pool.Put(buf)
	buf.Reset()

	buf.Write([]byte(k.code + "." + k.kind + " "))
	k.tpl.execute(buf, args)
	if err != nil {
		buf.Write([]byte(" ["))
		buf.Write([]byte(err.Error()))
		buf.Write([]byte("]"))
	}

	return &ErrorAttrs{
		code:  k.code,
		kind:  k.kind,
		stack: stack,
		msg:   buf.String(),
		name:  keyError,
		err:   err,
	}
}

// Is checks if the given error is the same as this kind.
// Is tries to unwrap the given errors if possible.
// "same" means the errors have the same error code.
// So, given error should implement errorutil.Coder and errorutil.Unwrapper interfaces.
func (k *Kind) Is(err error) bool {
	target := err
	for target != nil {
		c, ok := target.(ErrorKind)
		if ok && k.code == c.Code() && k.kind == c.Kind() {
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
