package compression

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
	"sync"

	"github.com/andybalholm/brotli"
)

const (
	// gzipEncoding is the value for Content-Encoding header
	// when compress bodies with Gzip.
	// 	- https://en.wikipedia.org/wiki/HTTP_compression
	// 	- https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Content-Encoding
	// 	- https://datatracker.ietf.org/doc/rfc1952/
	gzipEncoding = "gzip"

	// brotliEncoding is the value for Content-Encoding header
	// when compress bodies with Brotli.
	// 	-	https://en.wikipedia.org/wiki/HTTP_compression
	// 	- https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Content-Encoding
	// 	- https://datatracker.ietf.org/doc/rfc7932/
	brotliEncoding = "br"
)

// compression is a gzip and brotli compression middleware.
// References:
//
//   - https://datatracker.ietf.org/doc/rfc7231/ [Hypertext Transfer Protocol (HTTP/1.1): Semantics and Content]
//   - https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Accept-Encoding
//   - https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Content-Encoding
//   - https://developer.fastly.com/learning/concepts/compression/
//   - https://developers.cloudflare.com/speed/optimization/content/brotli/content-compression/
//   - https://docs.aws.amazon.com/AmazonCloudFront/latest/DeveloperGuide/ServingCompressedFiles.html
type compression struct {
	// mimes are all mime types to compress.
	// Wildcards such as "*" CANNOT be used.
	mimes []string

	// minimumSize is the byte size of
	// response body to be compressed.
	// Note that compressing small object (~1kB)
	// could be larger than the original one.
	minimumSize int64

	// gwPool is the gzip writer pool.
	gwPool sync.Pool
	// bwPool is the brotli writer pool.
	bwPool sync.Pool
}

func (c *compression) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// See the RFC7231 for possible formats of the Accept-Encoding headers.
		// 	- https://datatracker.ietf.org/doc/rfc7231/
		encoding := r.Header.Get("Accept-Encoding")

		switch {
		case strings.Contains(encoding, brotliEncoding): // Brotli compression.
			bw := c.bwPool.Get().(*brotli.Writer)
			defer func() {
				bw.Close()
				bw.Reset(io.Discard)
				c.bwPool.Put(bw)
			}()
			w = &compressionWriter{
				ResponseWriter: w,
				writer:         bw,
				mimes:          c.mimes,
				encoding:       brotliEncoding,
				minimumSize:    c.minimumSize,
			}
		case strings.Contains(encoding, gzipEncoding): // Gzip compression.
			gw := c.gwPool.Get().(*gzip.Writer)
			defer func() {
				gw.Close()
				gw.Reset(io.Discard)
				c.gwPool.Put(gw)
			}()
			w = &compressionWriter{
				ResponseWriter: w,
				writer:         gw,
				mimes:          c.mimes,
				encoding:       gzipEncoding,
				minimumSize:    c.minimumSize,
			}
		}

		next.ServeHTTP(w, r)
	})
}
