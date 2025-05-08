// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package soaprest

import (
	"bytes"
	"io"
	"mime"
	"net/http"

	"github.com/aileron-gateway/aileron-gateway/app"
	"github.com/aileron-gateway/aileron-gateway/app/middleware/soaprest/zxml"
	"github.com/aileron-gateway/aileron-gateway/core"
	utilhttp "github.com/aileron-gateway/aileron-gateway/util/http"
)

type soapREST struct {
	eh        core.ErrorHandler
	converter *zxml.JSONConverter
}

func (s *soapREST) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var ct string     // Content-Type for response.
		var action string // Action for X-SOAP-Action header.
		mt, params, _ := mime.ParseMediaType(r.Header.Get("Content-Type"))
		switch mt {
		case "application/soap+xml":
			ct = "application/soap+xml; charset=utf-8" // SOAP1.2
			action = params["action"]                  // 6.5 SOAP Action Feature (https://www.w3.org/TR/soap12-part2/)
		case "text/xml":
			if _, ok := r.Header["Soapaction"]; ok {
				ct = "text/xml; charset=utf-8"      // SOAP1.1
				action = r.Header.Get("Soapaction") // 6.1.1 The SOAPAction HTTP Header Field (http://www.w3.org/TR/SOAP)
				break
			}
			fallthrough
		default:
			// Neither SOAP1.1 nor 1.2
			err := app.ErrAppMiddleSOAPRESTVersionMismatch.WithoutStack(nil, nil)
			s.eh.ServeHTTPError(w, r, utilhttp.NewHTTPError(err, http.StatusForbidden))
			return
		}

		xmlBody, err := io.ReadAll(r.Body)
		if err != nil {
			err = app.ErrAppMiddleSOAPRESTReadXMLBody.WithoutStack(err, nil)
			s.eh.ServeHTTPError(w, r, utilhttp.NewHTTPError(err, http.StatusBadRequest))
			return
		}

		jsonBody, err := s.converter.XMLtoJSON(xmlBody) // XML to JSON.
		if err != nil {
			err = app.ErrAppMiddleSOAPRESTConvertXMLtoJSON.WithoutStack(err, nil)
			s.eh.ServeHTTPError(w, r, utilhttp.NewHTTPError(err, http.StatusBadRequest))
			return
		}

		newReq := r.Clone(r.Context())
		// Set ContentLength to -1 to indicate that the body size is unknown as per Go documentation:
		// "The value -1 indicates that the length is unknown."
		// Reference: https://pkg.go.dev/net/http#Request
		newReq.ContentLength = -1
		newReq.Body = io.NopCloser(bytes.NewReader(jsonBody))
		newReq.Header.Set("Accept", "application/json")
		newReq.Header.Set("Content-Type", "application/json")
		newReq.Header.Set("X-SOAP-Action", action)
		newReq.Header.Set("X-Content-Type", mt)

		ww := &wrappedWriter{
			ResponseWriter: w,
			body:           &bytes.Buffer{},
		}
		next.ServeHTTP(ww, newReq)
		w.Header().Del("Content-Length") // Delete Content-Length because the body will be modified.

		// Ensure the upstream handler actually returned JSON before we convert it back to XML.
		// If the Content-Type isn’t application/json, abort with an InvalidContentType error.
		if mt, _, _ := mime.ParseMediaType(w.Header().Get("Content-Type")); mt != "application/json" {
			err := app.ErrAppMiddleSOAPRESTInvalidContentType.WithoutStack(nil, map[string]any{"type": mt})
			s.eh.ServeHTTPError(w, r, utilhttp.NewHTTPError(err, http.StatusInternalServerError))
			return
		}

		respBody, err := s.converter.JSONtoXML(ww.body.Bytes()) // JSON to XML.
		if err != nil {
			err = app.ErrAppMiddleSOAPRESTConvertJSONtoXML.WithoutStack(err, nil)
			s.eh.ServeHTTPError(w, r, utilhttp.NewHTTPError(err, http.StatusInternalServerError))
			return
		}

		w.Header().Set("Content-Type", ct)
		w.WriteHeader(ww.StatusCode())
		_, err = w.Write(respBody) // Write XML content.
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

func (w *wrappedWriter) StatusCode() int {
	if w.written && w.code == 0 {
		return http.StatusOK
	}
	return w.code
}

func (w *wrappedWriter) Flush() {
	// No-op: prevent premature header/status commit.
	// If we propagated Flush to the inner ResponseWriter,
	// net/http would immediately send headers with status=200
	// (if WriteHeader wasn’t called), and any later conversion error
	// could not override that status.
}
