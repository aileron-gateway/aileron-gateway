// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package compression

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"sync"
	"testing"

	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
	"github.com/andybalholm/brotli"
)

func TestCompressionMiddleware(t *testing.T) {
	type condition struct {
		acceptEncoding string
		encoding       string
		contentType    string
		body           []byte
	}

	type action struct {
		encoding string
		body     string
	}

	var gzipBody, brBody bytes.Buffer
	gw := gzip.NewWriter(&gzipBody)
	gw.Write([]byte("test response body"))
	gw.Close()
	bw := brotli.NewWriter(&brBody)
	bw.Write([]byte("test response body"))
	bw.Close()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"accept header not exist",
			&condition{
				acceptEncoding: "",
				encoding:       "",
				contentType:    "text/plain",
				body:           []byte("test response body"),
			},
			&action{
				encoding: "",
				body:     "test response body",
			},
		),
		gen(
			"accept gzip",
			&condition{
				acceptEncoding: "gzip",
				encoding:       "",
				contentType:    "text/plain",
				body:           []byte("test response body"),
			},
			&action{
				encoding: "gzip",
				body:     "test response body",
			},
		),
		gen(
			"accept brotli",
			&condition{
				acceptEncoding: "br",
				encoding:       "",
				contentType:    "text/plain",
				body:           []byte("test response body"),
			},
			&action{
				encoding: "br",
				body:     "test response body",
			},
		),
		gen(
			"accept deflate",
			&condition{
				acceptEncoding: "deflate",
				encoding:       "",
				contentType:    "text/plain",
				body:           []byte("test response body"),
			},
			&action{
				encoding: "",
				body:     "test response body",
			},
		),
		gen(
			"accept gzip and brotli",
			&condition{
				acceptEncoding: "deflate, gzip, br",
				encoding:       "",
				contentType:    "text/plain",
				body:           []byte("test response body"),
			},
			&action{
				encoding: "br",
				body:     "test response body",
			},
		),
		gen(
			"short body",
			&condition{
				acceptEncoding: "gzip, br",
				encoding:       "",
				contentType:    "text/plain",
				body:           []byte("short body"),
			},
			&action{
				encoding: "",
				body:     "short body",
			},
		),
		gen(
			"already compressed with gzip",
			&condition{
				acceptEncoding: "gzip, br",
				encoding:       "gzip",
				contentType:    "text/plain",
				body:           gzipBody.Bytes(),
			},
			&action{
				encoding: "gzip",
				body:     "test response body",
			},
		),
		gen(
			"already compressed with brotli",
			&condition{
				acceptEncoding: "gzip, br",
				encoding:       "br",
				contentType:    "text/plain",
				body:           brBody.Bytes(),
			},
			&action{
				encoding: "br",
				body:     "test response body",
			},
		),
		gen(
			"no content-length gzip",
			&condition{
				acceptEncoding: "gzip, br",
				encoding:       "",
				contentType:    "text/plain",
				body:           []byte(nil),
			},
			&action{
				encoding: "",
				body:     "",
			},
		),
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			comp := &compression{
				mimes:       []string{"text/plain", "application/json"},
				minimumSize: 12,
				gwPool: sync.Pool{
					New: func() any { return gzip.NewWriter(io.Discard) },
				},
				bwPool: sync.Pool{
					New: func() any { return brotli.NewWriter(io.Discard) },
				},
			}

			h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", tt.C.contentType)
				if len(tt.C.body) > 0 {
					w.Header().Set("Content-Length", strconv.Itoa(len(tt.C.body)))
				}
				if tt.C.encoding != "" {
					w.Header().Set("Content-Encoding", tt.C.encoding)
				}
				w.WriteHeader(http.StatusOK)
				w.Write(tt.C.body)
			})

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			if tt.C.acceptEncoding != "" {
				req.Header.Set("Accept-Encoding", tt.C.acceptEncoding)
			}
			resp := httptest.NewRecorder()
			comp.Middleware(h).ServeHTTP(resp, req)

			testutil.Diff(t, http.StatusOK, resp.Code)
			testutil.Diff(t, tt.A.encoding, resp.Result().Header.Get("Content-Encoding"))

			b, _ := io.ReadAll(resp.Body)
			if strings.Contains(tt.A.encoding, "gzip") {
				r, _ := gzip.NewReader(bytes.NewReader(b))
				b, _ = io.ReadAll(r)
			}
			if strings.Contains(tt.A.encoding, "br") {
				r := brotli.NewReader(bytes.NewReader(b))
				b, _ = io.ReadAll(r)
			}
			testutil.Diff(t, tt.A.body, string(b))
		})
	}
}
