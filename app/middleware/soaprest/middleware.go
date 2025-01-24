// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package soaprest

import (
	"bytes"
	"io"
	"mime"
	"net/http"

	"github.com/aileron-gateway/aileron-gateway/app"
	"github.com/aileron-gateway/aileron-gateway/app/middleware/soaprest/xmlconv"
	"github.com/aileron-gateway/aileron-gateway/core"
	"github.com/aileron-gateway/aileron-gateway/kernel/txtutil"
	utilhttp "github.com/aileron-gateway/aileron-gateway/util/http"
)

type soapREST struct {
	eh core.ErrorHandler

	// paths is the path matcher to apply SOAP/REST conversion.
	// paths must not be nil.
	paths txtutil.Matcher[string]

	converter *xmlconv.Converter
}

func (s *soapREST) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// If the request path does not match the configured value,
		// the conversion process will not be executed, and the request will be passed to the next handler.
		if !s.paths.Match(r.URL.Path) {
			next.ServeHTTP(w, r)
			return
		}

		// If Content-Type is exactly "application/soap+xml", treat as SOAP1.2 request.
		// If Content-Type is exactly "text/xml" and SOAPAction header exists, treat as SOAP1.1 request.
		// If neither condition is met, return a VersionMismatch error with HTTP 403 Forbidden status.
		var ct string
		reqmt, params, _ := mime.ParseMediaType(r.Header.Get("Content-Type"))
		switch reqmt {
		case "application/soap+xml":
			// SOAP1.2 Request
			ct = "application/soap+xml"
		case "text/xml":
			_, soapActionExists := r.Header["Soapaction"]
			if soapActionExists {
				// SOAP1.1 Request
				ct = "text/xml"
				break
			}
			fallthrough
		default:
			// Neither SOAP1.1 nor 1.2
			err := app.ErrAppMiddleSOAPRESTVersionMismatch.WithoutStack(nil, nil)
			s.eh.ServeHTTPError(w, r, utilhttp.NewHTTPError(err, http.StatusForbidden))
			return
		}

		// Read the request body.
		xmlBody, err := io.ReadAll(r.Body)
		if err != nil {
			err = app.ErrAppMiddleSOAPRESTReadXMLBody.WithoutStack(err, nil)
			s.eh.ServeHTTPError(w, r, utilhttp.NewHTTPError(err, http.StatusBadRequest))
			return
		}
		r.Body.Close()

		// Convert XML to JSON.
		jsonBody, err := s.converter.XMLtoJSON(xmlBody)
		if err != nil {
			err = app.ErrAppMiddleSOAPRESTConvertXMLtoJSON.WithoutStack(err, nil)
			s.eh.ServeHTTPError(w, r, utilhttp.NewHTTPError(err, http.StatusBadRequest))
			return
		}

		// Create new request with JSON body.
		newReq := r.Clone(r.Context())

		// Set ContentLength to -1 to indicate that the body size is unknown as per Go documentation:
		// "The value -1 indicates that the length is unknown."
		// Reference: https://pkg.go.dev/net/http#Request
		newReq.ContentLength = -1
		newReq.Body = io.NopCloser(bytes.NewReader(jsonBody))
		newReq.Header.Set("Content-Type", "application/json")

		// Set X-SOAP-Action header consistently for both SOAP 1.1 and 1.2.
		switch reqmt {
		case "application/soap+xml":
			// SOAP 1.2: Get action from Content-Type parameter
			if act := params["action"]; act != "" {
				newReq.Header.Set("X-SOAP-Action", act)
			}
		case "text/xml":
			// SOAP 1.1: Get action from SOAPAction header
			soapAction := r.Header.Get("Soapaction")
			newReq.Header.Set("X-SOAP-Action", soapAction)
		}

		ww := &wrappedWriter{
			ResponseWriter: w,
			body:           &bytes.Buffer{},
		}

		next.ServeHTTP(ww, newReq)

		// Delete Content-Length because the body will be modified.
		ww.ResponseWriter.Header().Del("Content-Length")

		// Ensure the upstream handler actually returned JSON before we convert it back to XML.
		// If the Content-Type isn’t application/json, abort with an InvalidContentType error.
		respmt, _, _ := mime.ParseMediaType(ww.Header().Get("Content-Type"))
		if respmt != "application/json" {
			err := app.ErrAppMiddleSOAPRESTInvalidContentType.WithoutStack(nil, map[string]any{"type": respmt})
			s.eh.ServeHTTPError(w, r, utilhttp.NewHTTPError(err, http.StatusInternalServerError))
			return
		}

		// Convert JSON to XML.
		respBody, err := s.converter.JSONtoXML(ww.body.Bytes())
		if err != nil {
			err = app.ErrAppMiddleSOAPRESTConvertJSONtoXML.WithoutStack(err, nil)
			s.eh.ServeHTTPError(w, r, utilhttp.NewHTTPError(err, http.StatusInternalServerError))
			return
		}

		// The middleware supports only UTF‑8,
		// therefore we always add `charset=utf-8` to the Content‑Type.
		ct += "; charset=utf-8"
		ww.ResponseWriter.Header().Set("Content-Type", ct)

		w.WriteHeader(ww.StatusCode())

		// Modify response with XML body.
		_, err = ww.ResponseWriter.Write(respBody)
		if err != nil {
			err = app.ErrAppMiddleSOAPRESTWriteResponseBody.WithoutStack(err, nil)
			s.eh.ServeHTTPError(w, r, utilhttp.NewHTTPError(err, http.StatusInternalServerError))
		}
	})
}

// wrappedWriter wraps http.ResponseWriter.
// This implements io.Writer interface and utilhttp.Writer interface.
type wrappedWriter struct {
	http.ResponseWriter
	code    int
	written bool
	length  int64
	body    *bytes.Buffer
}

func (w *wrappedWriter) Unwrap() http.ResponseWriter {
	return w.ResponseWriter
}

func (w *wrappedWriter) WriteHeader(statusCode int) {
	if w.written {
		return
	}
	w.code = statusCode
	w.written = true
}

func (w *wrappedWriter) Write(b []byte) (int, error) {
	w.written = true
	w.length += int64(len(b))
	return w.body.Write(b)
}

func (w *wrappedWriter) Written() bool {
	return w.written
}

func (w *wrappedWriter) StatusCode() int {
	if w.written && w.code == 0 {
		return http.StatusOK
	}
	return w.code
}

func (w *wrappedWriter) ContentLength() int64 {
	return w.length
}

func (w *wrappedWriter) Flush() {
	// No-op: prevent premature header/status commit.
	// If we propagated Flush to the real ResponseWriter,
	// net/http would immediately send headers with status=200
	// (if WriteHeader wasn’t called), and any later conversion error
	// could not override that status.
}
