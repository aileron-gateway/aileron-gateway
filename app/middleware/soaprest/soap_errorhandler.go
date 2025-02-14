package soaprest

import (
	"bytes"
	"encoding/xml"
	"net/http"
	"strconv"

	"github.com/aileron-gateway/aileron-gateway/app"
	"github.com/aileron-gateway/aileron-gateway/core"
	"github.com/aileron-gateway/aileron-gateway/kernel/log"
	utilhttp "github.com/aileron-gateway/aileron-gateway/util/http"
)

// "MustUnderstand" is not handled in SOAP/REST conversion.
const (
	faultCodeClient          = "Client"
	faultCodeServer          = "Server"
	faultCodeVersionMismatch = "VersionMismatch"
)

type soap11Fault struct {
	XMLName     xml.Name           `xml:"Fault"`
	Faultcode   string             `xml:"faultcode"`
	Faultstring string             `xml:"faultstring"`
	Faultactor  string             `xml:"faultactor,omitempty"`
	Detail      *soap11FaultDetail `xml:"detail,omitempty"`
}

type soap11FaultDetail struct {
	XMLName    xml.Name `xml:"detail"`
	Message    string   `xml:"message,omitempty"`
	StatusCode int      `xml:"statusCode,omitempty"`
}

type soapFaultEnvelope struct {
	XMLName xml.Name       `xml:"http://schemas.xmlsoap.org/soap/envelope/ Envelope"`
	Body    *soapFaultBody `xml:"Body"`
}

type soapFaultBody struct {
	Fault *soap11Fault `xml:"Fault"`
}

// soapErrorHandler is a SOAP error handler that
// returns a SOAP Response containing a SOAPFault to the client.
// This implements core.
type soapErrorHandler struct {
	// LG is the logger to output logs.
	lg log.Logger
	// StackAlways is the flag to output
	// stacktraces even the client-side error.
	// If set to true, this error handler always
	// show stacktrace with debug level.
	stackAlways bool
}

func (h *soapErrorHandler) ServeHTTPError(w http.ResponseWriter, r *http.Request, err error) {
	statusCode := http.StatusInternalServerError
	faultCode := faultCodeServer
	detailMessage := "An error has occurred on the upstream server."

	if c, ok := err.(core.HTTPError); ok {
		statusCode = c.StatusCode()
		if statusCode < 500 {
			if app.ErrAppMiddleSOAPRESTVersionMismatch.Is(err) {
				// If a request that is not SOAP1.1 is received, return a VersionMismatch.
				faultCode = faultCodeVersionMismatch
				detailMessage = "Expected a SOAP 1.1 request, but received a request in a different format."
			} else {
				// Client-side error handling.
				faultCode = faultCodeClient
				detailMessage = "An error has occurred while processing the request from the client."
			}
		}
		err = c.Unwrap() // HTTPError is no more needed.
	}

	fault := &soap11Fault{
		Faultcode:   faultCode,
		Faultstring: http.StatusText(statusCode),
		// <TODO>
		// By default, r.Host() is specified, but it may be subject to change due to its lack of clarity.
		Faultactor: r.Host,
		Detail: &soap11FaultDetail{
			Message:    detailMessage,
			StatusCode: statusCode,
		},
	}

	// Log output.
	if statusCode >= 500 {
		name, value := utilhttp.LogAttr(err, true)
		h.lg.Error(r.Context(), "serve http error. status="+strconv.Itoa(statusCode), name, value)
	} else if h.lg.Enabled(log.LvDebug) {
		name, value := utilhttp.LogAttr(err, h.stackAlways)
		h.lg.Debug(r.Context(), "serve http error. status="+strconv.Itoa(statusCode), name, value)
	}

	// Status code less than 100 is defined as logging only error.
	// So, return here without writing response.
	if statusCode < 100 {
		return
	}

	envelope := &soapFaultEnvelope{
		Body: &soapFaultBody{
			Fault: fault,
		},
	}

	var buf bytes.Buffer
	buf.WriteString(xml.Header)

	encoder := xml.NewEncoder(&buf)
	encoder.Indent("", "  ")

	// While the contents of the `Fault` can change due to errors,
	// there will be no failure in the encoding process, so error handling is not implemented.
	encoder.Encode(envelope)

	header := w.Header()
	header.Set("Content-Type", "text/xml; charset=utf-8")
	header.Set("X-Content-Type-Options", "nosniff")
	header.Add("Vary", "Accept")
	header.Set("Content-Length", strconv.Itoa(buf.Len()))

	w.WriteHeader(statusCode)

	if _, err := buf.WriteTo(w); err != nil {
		name, value := utilhttp.LogAttr(err, true)
		h.lg.Error(r.Context(), "failed to write response body to client", name, value)
		return
	}
}
