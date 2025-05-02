// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package http

import (
	"cmp"
	"math"
	"mime"
	"net/http"
	"net/textproto"
	"path"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"sync"

	v1 "github.com/aileron-gateway/aileron-gateway/apis/core/v1"
	k "github.com/aileron-gateway/aileron-gateway/apis/kernel"
	"github.com/aileron-gateway/aileron-gateway/core"
	"github.com/aileron-gateway/aileron-gateway/kernel/api"
	"github.com/aileron-gateway/aileron-gateway/kernel/er"
	"github.com/aileron-gateway/aileron-gateway/kernel/errorutil"
	"github.com/aileron-gateway/aileron-gateway/kernel/log"
	"github.com/aileron-gateway/aileron-gateway/kernel/txtutil"
)

// Following error statuses do not contain inner error.
// These error status are shared across the entire application.
// So, their state should not be modified which means
// Header or AddError methods are not allowed to be used.
// Errors with headers or ErrorElems should be created by NewHTTPError.

var (
	ErrLoggingOnly                 error = NewHTTPError(nil, -1)                                     // -1  Logging only status for default error handler.
	ErrBadRequest                  error = NewHTTPError(nil, http.StatusBadRequest)                  // 400 RFC 9110, 15.5.1
	ErrUnauthorized                error = NewHTTPError(nil, http.StatusUnauthorized)                // 401 RFC 9110, 15.5.2
	ErrPaymentRequired             error = NewHTTPError(nil, http.StatusPaymentRequired)             // 402 RFC 9110, 15.5.3
	ErrForbidden                   error = NewHTTPError(nil, http.StatusForbidden)                   // 403 RFC 9110, 15.5.4
	ErrNotFound                    error = NewHTTPError(nil, http.StatusNotFound)                    // 404 RFC 9110, 15.5.5
	ErrMethodNotAllowed            error = NewHTTPError(nil, http.StatusMethodNotAllowed)            // 405 RFC 9110, 15.5.6
	ErrNotAcceptable               error = NewHTTPError(nil, http.StatusNotAcceptable)               // 406 RFC 9110, 15.5.7
	ErrProxyAuthRequired           error = NewHTTPError(nil, http.StatusProxyAuthRequired)           // 407 RFC 9110, 15.5.8
	ErrRequestTimeout              error = NewHTTPError(nil, http.StatusRequestTimeout)              // 408 RFC 9110, 15.5.9
	ErrConflict                    error = NewHTTPError(nil, http.StatusConflict)                    // 409 RFC 9110, 15.5.10
	ErrGone                        error = NewHTTPError(nil, http.StatusGone)                        // 410 RFC 9110, 15.5.11
	ErrLengthRequired              error = NewHTTPError(nil, http.StatusLengthRequired)              // 411 RFC 9110, 15.5.12
	ErrPreconditionFailed          error = NewHTTPError(nil, http.StatusPreconditionFailed)          // 412 RFC 9110, 15.5.13
	ErrRequestEntityTooLarge       error = NewHTTPError(nil, http.StatusRequestEntityTooLarge)       // 413 RFC 9110, 15.5.14
	ErrRequestURITooLong           error = NewHTTPError(nil, http.StatusRequestURITooLong)           // 414 RFC 9110, 15.5.15
	ErrUnsupportedMediaType        error = NewHTTPError(nil, http.StatusUnsupportedMediaType)        // 415 RFC 9110, 15.5.16
	ErrTooManyRequests             error = NewHTTPError(nil, http.StatusTooManyRequests)             // 429 RFC 6585, 4
	ErrRequestHeaderFieldsTooLarge error = NewHTTPError(nil, http.StatusRequestHeaderFieldsTooLarge) // 431 RFC 6585, 5
	ErrInternalServerError         error = NewHTTPError(nil, http.StatusInternalServerError)         // 500 RFC 9110, 15.6.1
	ErrNotImplemented              error = NewHTTPError(nil, http.StatusNotImplemented)              // 501 RFC 9110, 15.6.2
	ErrBadGateway                  error = NewHTTPError(nil, http.StatusBadGateway)                  // 502 RFC 9110, 15.6.3
	ErrServiceUnavailable          error = NewHTTPError(nil, http.StatusServiceUnavailable)          // 503 RFC 9110, 15.6.4
	ErrGatewayTimeout              error = NewHTTPError(nil, http.StatusGatewayTimeout)              // 504 RFC 9110, 15.6.5
)

// DefaultErrorHandlerName is the name of default error handler.
// Use SetGlobalErrorHandler method to replace the default error handler.
// Even the default error handler can be replaced, it cannot be set to nil.
// To get the default error handler use this like below.
//
//	eh := http.GlobalErrorHandler(http.DefaultErrorHandlerName)
const DefaultErrorHandlerName = "__default__"

var (
	// mu protects handlers.
	mu = sync.RWMutex{}
	// handlers is the global error handler set.
	handlers = map[string]core.ErrorHandler{
		DefaultErrorHandlerName: &DefaultErrorHandler{
			LG: log.GlobalLogger(log.DefaultLoggerName),
		},
	}
)

// GlobalErrorHandler returns a error handler which is stored in the global error handler holder by name.
// A default registered error handler can be obtained by the name http.DefaultErrorHandlerName.
// If there is no handler, GlobalErrorHandler returns nil.
//
// When getting the default error handler, use like below.
// The error handler gotten by the name DefaultErrorHandlerName is always non nil.
//
//	eh := http.GlobalErrorHandler(http.DefaultErrorHandlerName)
//
// If it's not the default logger, nil check should be taken.
//
//	eh := http.GlobalErrorHandler("yourErrorHandlerName")
//	if eh == nil {
//		// use other error handler.
//	}
func GlobalErrorHandler(name string) core.ErrorHandler {
	mu.RLock()
	defer mu.RUnlock()
	eh, ok := handlers[name]
	if !ok {
		return nil
	}
	return eh
}

// SetGlobalErrorHandler stores the given error handler in the global error handler holder.
// This replaces the existing error handler if there have already been the same named handler.
// To delete the handler, set nil as the second argument.
// The error handler named http.DefaultErrorHandlerName can be replaced but cannot be deleted.
//
// To delete error handler:
//
//	http.SetGlobalErrorHandler("errorHandlerName", nil)
//
// To replace default error handler:
//
//	var eh core.ErrorHandler
//	eh = <error handler you want to use>
//	http.SetGlobalErrorHandler(http.DefaultErrorHandlerName, eh)
func SetGlobalErrorHandler(name string, handler core.ErrorHandler) {
	mu.RLock()
	defer mu.RUnlock()
	if handler == nil {
		if name != DefaultErrorHandlerName {
			delete(handlers, name)
		}
		return
	}
	handlers[name] = handler
}

// ErrorHandler returns a error handler by getting it from the given api.
// The default error handler will be returned when a nil reference was given by the
// second argument ref.
// This function can panic if the first argument api is nil.
func ErrorHandler(a api.API[*api.Request, *api.Response], ref *k.Reference) (core.ErrorHandler, error) {
	if ref == nil {
		return GlobalErrorHandler(DefaultErrorHandlerName), nil
	}
	eh, err := api.ReferTypedObject[core.ErrorHandler](a, ref)
	if err != nil {
		return nil, err // Return err as-is
	}
	return eh, nil
}

func NewErrorMessage(spec *v1.ErrorMessageSpec) (*ErrorMessage, error) {
	if spec == nil || len(spec.MIMEContents) == 0 {
		return nil, nil
	}

	m := &ErrorMessage{
		codes:     spec.Codes,
		kinds:     spec.Kinds,
		headerTpl: make(map[string]*txtutil.FastTemplate, len(spec.HeaderTemplate)),
	}

	for k, tpl := range spec.HeaderTemplate {
		m.headerTpl[textproto.CanonicalMIMEHeaderKey(k)] = txtutil.NewFastTemplate(tpl, "{{", "}}")
	}

	for _, msg := range spec.Messages {
		tpl, err := regexp.Compile(msg)
		if err != nil {
			return nil, (&er.Error{
				Package:     ErrPkg,
				Type:        ErrTypeErrHandler,
				Description: ErrDscRegexp,
			}).Wrap(err)
		}
		m.messages = append(m.messages, tpl)
	}

	for _, cs := range spec.MIMEContents {
		c, err := NewMIMEContent(cs)
		if err != nil {
			return nil, err // Return err as-is.
		}
		m.contents = append(m.contents, c)
	}

	return m, nil
}

type ErrorMessage struct {
	codes     []string
	kinds     []string
	messages  []*regexp.Regexp
	headerTpl map[string]*txtutil.FastTemplate
	contents  []*MIMEContent
}

// Match returns if the given error matched to this message.
// Codes, kinds and messages are evaluated by AND condition.
func (m *ErrorMessage) Match(code, kind string, msg []byte) bool {
	for _, c := range m.codes {
		if matched, _ := path.Match(c, code); matched {
			return true
		}
	}
	for _, k := range m.kinds {
		if matched, _ := path.Match(k, kind); matched {
			return true
		}
	}
	for _, m := range m.messages {
		if m.Match(msg) {
			return true
		}
	}
	return false
}

func (m *ErrorMessage) Content(accept string) *MIMEContent {
	if len(m.contents) == 0 {
		return nil
	}
	accepts := strings.Split(accept, ",")
	for i := range accepts {
		mimeType, _, _ := mime.ParseMediaType(accepts[i])
		for j := 0; j < len(m.contents); j++ {
			matched, _ := path.Match(mimeType, m.contents[j].MIMEType)
			if matched {
				return m.contents[j]
			}
		}
	}
	// Use the first MIMEContent if not match to content type.
	return m.contents[0]
}

// DefaultErrorHandler is a HTTP error handler that
// returns HTTP error response to the clients.
// This implements core.
type DefaultErrorHandler struct {
	// LG is the logger to output logs.
	LG log.Logger
	// StackAlways is the flag to output
	// stacktraces even the client-side error.
	// If set to true, this error handler always
	// show stacktrace with debug level.
	StackAlways bool
	// Msgs is the list of error messages
	// to overwrite the default.
	Msgs []*ErrorMessage
}

func (h *DefaultErrorHandler) ServeHTTPError(w http.ResponseWriter, r *http.Request, err error) {
	statusCode := math.MaxInt // Define initial state.
	var mimeType string
	var body []byte
	var header1, header2 http.Header

	if h, ok := err.(interface{ Header() http.Header }); ok {
		header1 = h.Header()
	}
	if c, ok := err.(core.HTTPError); ok {
		statusCode = c.StatusCode()
		mimeType, body = c.Content(r.Header.Get("Accept"))
		err = c.Unwrap() // HTTPError is no more needed.
	}

	// Overwrite response if configured.
	// Do not overwrite logging only error (statusCode<100).
	if statusCode >= 100 && err != nil {
		var errCode, errKind string
		errMsg := []byte(err.Error())
		if ek, ok := err.(errorutil.ErrorKind); ok {
			errCode, errKind = ek.Code(), ek.Kind()
		}
		for _, m := range h.Msgs {
			if m.Match(errCode, errKind, errMsg) {
				mc := m.Content(r.Header.Get("Accept"))
				if mc == nil {
					continue
				}
				statusCode = cmp.Or(mc.StatusCode, statusCode) // Use current status if zero.
				mimeType = mc.MIMEType
				header2 = mc.Header
				info := map[string]any{"status": statusCode, "statusText": http.StatusText(statusCode), "code": errCode, "kind": errKind}
				body = mc.Content(info)
				for k, tpl := range m.headerTpl {
					header2.Add(k, string(tpl.Execute(info)))
				}
				break
			}
		}
	}

	if statusCode == math.MaxInt { // Use default status 500.
		e := NewHTTPError(err, http.StatusInternalServerError)
		statusCode = http.StatusInternalServerError
		mimeType, body = e.Content("application/json")
	}

	// Log output.
	if statusCode >= 500 {
		name, value := logAttr(err, true)
		h.LG.Error(r.Context(), "serve http error. status="+strconv.Itoa(statusCode), name, value)
	} else if h.LG.Enabled(log.LvDebug) {
		name, value := logAttr(err, h.StackAlways)
		h.LG.Debug(r.Context(), "serve http error. status="+strconv.Itoa(statusCode), name, value)
	}

	// Status code less than 100 is defined as logging only error.
	// So, return here without writing response.
	if statusCode < 100 {
		return
	}

	header := w.Header()
	copyHeader(header, header1)
	copyHeader(header, header2)
	header.Set("Content-Type", mimeType+"; charset=utf-8")
	header.Set("X-Content-Type-Options", "nosniff")
	header.Add("Vary", "Accept")
	w.WriteHeader(statusCode)
	_, _ = w.Write(body)
}

// logAttr returns key value to output logs.
// This remove or add stack trace depending on the given stack.
// First argument err must not be nil.
// If the stack is true, returned map always contains stack trace.
// if the stack is false, returned map does not always contain stack trace.
func logAttr(err error, stack bool) (string, map[string]any) {
	// First, assert the error to log attribute.
	attr, ok := err.(log.Attributes)
	if !ok {
		attr = core.ErrPrimitive.WithoutStack(err, nil)
	}
	attrMap := attr.Map()

	// Depending on the kernel/errorutil package,
	// the map always contains "stack" filed.
	// Delete that stack trace when the stack is false.
	if !stack {
		attrMap["stack"] = ""
		return attr.Name(), attrMap
	}

	// If the map does not have any stack trace,
	// add stack traces here.
	s, ok := attrMap["stack"].(string)
	if !ok || s == "" {
		b := make([]byte, 3*1<<10)
		n := runtime.Stack(b, false)
		attrMap["stack"] = string(b[:n])
	}

	// Here, attrMap always contains non-empty stack trace.
	return attr.Name(), attrMap
}

// copyHeader copies headers from src to dst.
// The destination header, dst, must not be nil.
func copyHeader(dst, src http.Header) {
	for k, vv := range src {
		dst[k] = append(dst[k], vv...)
	}
}
