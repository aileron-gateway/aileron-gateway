// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package static

import (
	"bytes"
	"io/fs"
	"net/http"
	"strings"

	"github.com/aileron-gateway/aileron-gateway/core"
	utilhttp "github.com/aileron-gateway/aileron-gateway/util/http"
)

type handler struct {
	// HandlerBase is the base struct for
	// http.Handler type resource.
	// This provides Patterns() and Methods() methods
	// to fulfill the core.Handler interface.
	*utilhttp.HandlerBase

	http.Handler

	eh core.ErrorHandler

	// header is the list of HTTP response headers.
	// Users can add response headers like "MaxAge"
	// or any other headers by this field.
	header map[string]string
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	upath := r.URL.Path
	if !strings.HasPrefix(upath, "/") {
		upath = "/" + upath
	}
	r.URL.Path = upath

	w.Header().Set("X-Content-Type-Options", "nosniff")
	for k, v := range h.header {
		w.Header().Set(k, v)
	}

	dw := &discardWriter{
		ResponseWriter: w,
	}

	defer func() {
		if dw.status < 400 {
			return
		}
		// Overwrite response if the response status is grater than or equal to 400.
		// The response had been discarded by the discardWriter. Write error response we prepared.
		err := core.ErrCoreStaticServer.WithoutStack(nil, map[string]any{"body": dw.body.String()})
		h.eh.ServeHTTPError(w, r, utilhttp.NewHTTPError(err, dw.status))
	}()

	h.Handler.ServeHTTP(dw, r)
}

// fileOnlyDir is a file system only allows access
// to files and not directories.
// This implements http.FileSystem interface.
type fileOnlyDir struct {
	fs http.FileSystem
}

func (d *fileOnlyDir) Open(name string) (http.File, error) {
	f, err := d.fs.Open(name)
	if err != nil {
		return nil, err
	}

	info, err := f.Stat()
	if err != nil {
		return nil, err
	}

	// Do not accept accessing to directory.
	if info.IsDir() {
		return nil, fs.ErrNotExist
	}

	return f, nil
}

// discardWriter discard http status and response body
// when the status is grater than or equal to 400.
type discardWriter struct {
	http.ResponseWriter
	status int
	body   bytes.Buffer
}

func (w *discardWriter) Write(b []byte) (int, error) {
	if w.status >= 400 {
		return w.body.Write(b)
	}
	return w.ResponseWriter.Write(b)
}

func (w *discardWriter) WriteHeader(statusCode int) {
	w.status = statusCode
	if w.status >= 400 {
		// Do not write status into internal writer.
		return
	}
	w.ResponseWriter.WriteHeader(statusCode)
}
