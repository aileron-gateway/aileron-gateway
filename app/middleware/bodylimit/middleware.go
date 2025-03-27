package bodylimit

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync/atomic"
	"time"

	"github.com/aileron-gateway/aileron-gateway/app"
	"github.com/aileron-gateway/aileron-gateway/core"
	kio "github.com/aileron-gateway/aileron-gateway/kernel/io"
	utilhttp "github.com/aileron-gateway/aileron-gateway/util/http"
)

// limitReadCloser limits the data length
// that is read from this reader.
type limitReadCloser struct {
	io.ReadCloser
	maxSize    int64
	read       int64
	exceedFunc func()
}

// Unwrap returns inner ReadCloser.
func (l *limitReadCloser) Unwrap() io.ReadCloser {
	return l.ReadCloser
}

// Read read data and copy to the given p.
// Read panics http.ErrAbortHandler when the read size exceeded the limit.
func (l *limitReadCloser) Read(p []byte) (n int, err error) {
	n, err = l.ReadCloser.Read(p)
	l.read += int64(n)
	if l.read > l.maxSize {
		l.exceedFunc()
		panic(http.ErrAbortHandler)
	}
	return n, err
}

// counter is the num counter that is added to the temp file name.
// This counter is reset to zero when the
// value exceeded the maximum number of uint64 (18446744073709551615).
var counter atomic.Uint64

// bodyLimit is the middleware that limit the actual body size.
// This implements core.Middleware interface.
type bodyLimit struct {
	eh core.ErrorHandler
	// maxSize is the maximum size of request body.
	// Size check is skipped if the maxSize is zero or negative.
	maxSize int64
	// tempPath is the temporary directory path to save request body
	// that exceeds the memLimit.
	tempPath string
	// memLimit is the maximum memory size to keep request bodies
	// on the memory when checking the actual body size.
	// If request body size exceeds this limit, the body is save
	// in the temporary file under the tempPath.
	// If memLimit is zero or negative value,
	memLimit int64
}

func (m *bodyLimit) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// We skip the body size check
		// if the maxSize is zero or negative.
		if m.maxSize <= 0 {
			next.ServeHTTP(w, r)
			return
		}

		// First, we trust and check the Content-Length.
		// Serve error if it exceeds the maxSize.
		if r.ContentLength > m.maxSize {
			_, _ = io.Copy(io.Discard, r.Body)
			err := app.ErrAppMiddleBodyTooLarge.WithoutStack(nil, nil)
			m.eh.ServeHTTPError(w, r, utilhttp.NewHTTPError(err, http.StatusRequestEntityTooLarge))
			return
		}

		// This case, do not load the body in this middleware.
		// Wrap the body with limitReadCloser and panic(http.ErrAbortHandler)
		// when total read bytes exceeded the maxSize.
		if r.ContentLength < 0 || m.memLimit <= 0 {
			ww := utilhttp.WrapWriter(w)
			w = ww
			r.Body = &limitReadCloser{
				ReadCloser: r.Body,
				maxSize:    m.maxSize,
				exceedFunc: func() {
					if !ww.Written() {
						err := app.ErrAppMiddleBodyTooLarge.WithoutStack(nil, nil)
						m.eh.ServeHTTPError(w, r, utilhttp.NewHTTPError(err, http.StatusRequestEntityTooLarge))
					}
				},
			}
			next.ServeHTTP(w, r)
			return
		}

		if r.ContentLength <= m.memLimit {
			// This case, load the body content on the memory
			// up to r.ContentLength length.
			body := make([]byte, r.ContentLength)
			n, err := io.ReadFull(r.Body, body)
			_, _ = io.Copy(io.Discard, r.Body) // Just in case.
			if err != nil && err != io.EOF {
				err = app.ErrAppMiddleInvalidLength.WithoutStack(err, nil)
				m.eh.ServeHTTPError(w, r, utilhttp.NewHTTPError(err, http.StatusBadRequest))
				return
			}
			if int64(n) != r.ContentLength {
				err := app.ErrAppMiddleInvalidLength.WithoutStack(nil, nil)
				m.eh.ServeHTTPError(w, r, utilhttp.NewHTTPError(err, http.StatusBadRequest))
				return
			}
			r.Body = io.NopCloser(bytes.NewReader(body))
		} else {
			// This case, load the body content on the temp file up to r.ContentLength
			filePath := m.tempPath + "body-" + time.Now().Format("20060102150405.000000-") + fmt.Sprintf("%020d", counter.Add(1))
			f, err := os.OpenFile(filePath, os.O_CREATE|os.O_TRUNC|os.O_RDWR, os.ModePerm)
			if err != nil {
				err := app.ErrAppMiddleBodyLimit.WithoutStack(err, nil)
				m.eh.ServeHTTPError(w, r, utilhttp.NewHTTPError(err, http.StatusInternalServerError))
				return
			}
			defer func() {
				f.Close()
				os.Remove(filePath)
			}()
			n, err := kio.CopyBuffer(f, io.LimitReader(r.Body, m.maxSize))
			if err != nil {
				err := app.ErrAppMiddleBodyLimit.WithoutStack(err, nil)
				m.eh.ServeHTTPError(w, r, utilhttp.NewHTTPError(err, http.StatusInternalServerError))
				return
			}
			if n != r.ContentLength {
				// ContentLength and the actual body size are different.
				// It should be a bad request in this case.
				err := app.ErrAppMiddleInvalidLength.WithoutStack(nil, nil)
				m.eh.ServeHTTPError(w, r, utilhttp.NewHTTPError(err, http.StatusBadRequest))
				return
			}
			f.Seek(0, 0)
			r.Body = f
		}

		next.ServeHTTP(w, r)
	})
}
