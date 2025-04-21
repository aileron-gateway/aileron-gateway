// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package http

import (
	"encoding/json"
	"encoding/xml"
	"mime"
	"net/http"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

// ErrorElems represents aggregated errors.
type ErrorElems struct {
	// XMLName is the name of this xml space.
	// This field is only for xml marshaling.
	XMLName xml.Name `xml:"result" json:"-" yaml:"-"`
	// Status is a HTTP status code.
	// For example, 500.
	Status int `json:"status" yaml:"status" xml:"status"`
	// StatusText is a HTTP status text.
	// For example, "InternalServerError".
	StatusText string `json:"statusText" yaml:"statusText" xml:"statusText"`
	// Errors is the list of errors.
	Errors []*ErrorElem `json:"errors,omitempty" yaml:"errors,omitempty" xml:"errors>error,omitempty"`
}

// ErrorElem represents a single error.
type ErrorElem struct {
	// XMLName is the name of this xml space.
	// This field is only for xml marshaling.
	XMLName xml.Name `xml:"error" json:"-" yaml:"-"`
	// Code is an error code.
	// For example "E1234".
	Code string `json:"code" yaml:"code" xml:"code"`
	// Message is an error message.
	// For example, "invalid request header".
	Message string `json:"message" yaml:"message" xml:"message"`
	// Detail is more details about this error.
	Detail string `json:"detail,omitempty" yaml:"detail,omitempty" xml:"detail,omitempty"`
}

// NewHTTPError returns a new HTTP error instance.
func NewHTTPError(err error, status int) *HTTPError {
	return &HTTPError{
		inner:  err,
		status: status,
		header: make(http.Header, 0),
	}
}

// HTTPError is the error response information.
// This struct holds and provides information about HTTP response error.
// Use NewHTTPError function to create a new instance.
// This implements core.ErrorResponse interface.
type HTTPError struct {
	inner  error        // error is the internal error.
	errs   []*ErrorElem // Internal errors.
	status int          // Status is the response status code.
	header http.Header  // Header is the additional response header.
}

func (e *HTTPError) Error() string {
	if e.inner != nil {
		return e.inner.Error()
	}
	return "http status " + strconv.Itoa(e.status) + " " + http.StatusText(e.status)
}

// AddError adds error element to this error object.
func (e *HTTPError) AddError(err *ErrorElem) {
	if err != nil {
		e.errs = append(e.errs, err)
	}
}

// Unwrap returns the internal error if any.
func (e *HTTPError) Unwrap() error {
	return e.inner
}

// StatusCode returns http status code.
func (e *HTTPError) StatusCode() int {
	return e.status
}

// Header returns http header.
// Note that modifying the returned header changes
// the state of this HTTPError object.
func (e *HTTPError) Header() http.Header {
	return e.header
}

// Content returns content type and body content.
// Supported content types are as follows.
// Default "application/json" will be returned when
// unsupported content type was given.
//
//   - "application/json", "text/json"
//   - "application/xml", "text/xml"
//   - "application/yaml", "text/yaml"
//   - "text/plain"
func (e *HTTPError) Content(accept string) (string, []byte) {
	elems := &ErrorElems{
		Status:     e.status,
		StatusText: http.StatusText(e.status),
		Errors:     e.errs,
	}
	accepts := strings.Split(accept, ",")
	for i := range accepts {
		mediatype, _, _ := mime.ParseMediaType(accepts[i])
		if mediatype == "" {
			continue
		}
		switch mediatype {
		case "application/json", "text/json":
			body, _ := json.Marshal(elems)
			return mediatype, body
		case "application/xml", "text/xml":
			body, _ := xml.Marshal(elems)
			return mediatype, body
		case "application/yaml", "text/yaml", "text/plain":
			body, _ := yaml.Marshal(elems)
			return mediatype, body
		default:
			body, _ := json.Marshal(elems)
			return "application/json", body
		}
	}
	// Fallback to json.
	body, _ := json.Marshal(elems)
	return "application/json", body
}
