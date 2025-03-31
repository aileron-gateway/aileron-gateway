package httplogger

import (
	"bytes"
	"context"
	"encoding/base64"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"

	v1 "github.com/aileron-gateway/aileron-gateway/apis/core/v1"
	"github.com/aileron-gateway/aileron-gateway/kernel/er"
	kio "github.com/aileron-gateway/aileron-gateway/kernel/io"
	"github.com/aileron-gateway/aileron-gateway/kernel/log"
	"github.com/aileron-gateway/aileron-gateway/kernel/txtutil"
)

// stringReplFunc is the value replace function for string.
// This functions is defined for header value replacing.
type stringReplFunc txtutil.ReplaceFunc[string]

// bytesReplFunc is the value replace function for []byte.
// This functions is defined for body value replacing.
type bytesReplFunc txtutil.ReplaceFunc[[]byte]

func stringReplacerToFunc(rs []txtutil.Replacer[string]) []stringReplFunc {
	fs := make([]stringReplFunc, 0, len(rs))
	for _, r := range rs {
		fs = append(fs, r.Replace)
	}
	return fs
}

func bytesReplacerToFunc(rs []txtutil.Replacer[[]byte]) []bytesReplFunc {
	fs := make([]bytesReplFunc, 0, len(rs))
	for _, r := range rs {
		fs = append(fs, r.Replace)
	}
	return fs
}

// pool is the pool of *bytes.Buffer.
// Make sure to reset the buffer before use it.
// Do not use the obtained []bytes by buf.Bytes()
// after putting the buffer back.
//
//	buf := pool.Get().(*bytes.Buffer)
//	defer pool.Put(buf)
//	buf.Reset()
var pool = &sync.Pool{
	New: func() any {
		var buf bytes.Buffer
		return &buf
	},
}

func newBaseLogger(spec *v1.LoggingSpec, lg log.Logger) (*baseLogger, error) {
	bl := &baseLogger{
		lg:      lg,
		w:       os.Stderr, // Will be replaced.
		mimes:   spec.MIMEs,
		maxBody: spec.MaxContentLength,
		base64:  spec.Base64,
	}

	var err error
	bl.headers, bl.allHeaders, err = headerReplacers(spec.Headers)
	if err != nil {
		return nil, err
	}
	bl.headerKeys = make([]string, 0, len(bl.headers))
	for k := range bl.headers {
		bl.headerKeys = append(bl.headerKeys, k)
	}

	bl.bodies, err = bodyReplacers(spec.Bodies)
	if err != nil {
		return nil, err
	}

	fs, err := txtutil.NewStringReplacers(spec.Queries...)
	if err != nil {
		return nil, err
	}
	bl.queries = stringReplacerToFunc(fs)

	if spec.LogFormat != "" {
		bl.tpl = txtutil.NewFastTemplate(spec.LogFormat+"\n", "%", "%")
		writer, ok := lg.(io.Writer) // Format logging requires raw io.Writer.
		if !ok {
			return nil, &er.Error{
				Package:     "httplogger",
				Type:        "base logger",
				Description: "formatted log requires logger with io.Writer interface",
			}
		}
		bl.w = writer
	}

	if spec.BodyOutputPath != "" {
		bl.bodyPath = filepath.Clean(spec.BodyOutputPath) + "/"
		_ = os.MkdirAll(bl.bodyPath, os.ModePerm)
		if err := kio.ReadWriteTest(bl.bodyPath); err != nil {
			return nil, err
		}
	}

	return bl, nil
}

type baseLogger struct {
	lg  log.Logger
	w   io.Writer // Used for formatted logger.
	tpl *txtutil.FastTemplate

	// queries is the URL query value replacers.
	queries []stringReplFunc

	// allHeaders is the flag to output all
	// header values to the log.
	// If false, only headers set in headers are output.
	allHeaders bool
	// headers is the header names to output.
	// Header values are masked by replaceFunc if specified.
	// Map keys are formatted by textproto.CanonicalMIMEHeaderKey.
	headers map[string][]stringReplFunc
	// headerKeys is the list of header names
	// that should be output or replaced.
	// headerKeys is the same with all key strings
	// os headers field.
	headerKeys []string

	// reqBodies is the functions that mask
	// the request body.
	// Map keys must be a valid MIME type such as
	// "application/json" or "application/x-www-form-urlencoded".
	bodies map[string][]bytesReplFunc

	// mimes is the list of media types
	// to log request/response bodies.
	mimes []string
	// maxBody is the maximum body size
	// that can load on memory for logging bodies.
	// maxBody is ignored when bodyPath is set.
	maxBody int64
	// bodyPath is the file directory
	// to output request and response bodies.
	// If bodyPath is set, maxBody and base64
	// and body replacers are ignored.
	bodyPath string
	// base64 if true, output logs
	// in base64 encoded format.
	// base64 is ignored when bodyPath is set.
	base64 bool
}

func (l *baseLogger) logOutput(ctx context.Context, msg string, attrs []any, tagFunc func(tag string) []byte) {
	if !l.lg.Enabled(log.LvInfo) {
		return
	}
	if l.tpl != nil {
		buf := pool.Get().(*bytes.Buffer)
		defer pool.Put(buf)
		buf.Reset()
		l.tpl.ExecuteFuncWriter(buf, tagFunc)
		_, _ = l.w.Write(buf.Bytes())
	} else {
		l.lg.Info(ctx, msg, attrs...)
	}
}

// logQuery returns query string that should be output to logs.
// logQuery applies query replacers to the given query string.
func (l *baseLogger) logQuery(q string) string {
	for i := range l.queries {
		q = l.queries[i](q)
	}
	return q
}

// logHeaders returns the header key values
// that should be output to logs.
// logHeaders applies header replacers to the given header.
func (l *baseLogger) logHeaders(h http.Header) map[string]string {
	var keys []string
	if l.allHeaders {
		keys = append(keys, make([]string, 0, len(h))...)
		for k := range h {
			keys = append(keys, k)
		}
	} else {
		keys = l.headerKeys
	}

	hs := make(map[string]string, len(keys))

	for _, key := range keys {
		vs := h[key]
		if len(vs) == 0 {
			continue
		}
		var val string
		if len(vs) == 1 {
			val = vs[0]
		} else {
			val = strings.Join(vs, ",")
		}

		repls := l.headers[key]
		for i := 0; i < len(repls); i++ {
			val = repls[i](val)
		}
		hs[key] = val
	}

	return hs
}

// logBody returns the body that should be output to logs.
// logBody applies body replacers to the given body.
func (l *baseLogger) logBody(mimeType string, body []byte) []byte {
	if len(body) == 0 {
		return nil
	}
	replacers := l.bodies[mimeType]
	for i := 0; i < len(replacers); i++ {
		body = replacers[i](body)
	}
	return body
}

func (l *baseLogger) bodyReadCloser(fileName, mimeType string, length int64, body io.ReadCloser, isCompressed bool) ([]byte, io.ReadCloser, error) {
	if !slices.Contains(l.mimes, mimeType) {
		return nil, body, nil
	}

	if length == 0 {
		return nil, body, nil
	}

	if length > 0 && length < l.maxBody {
		b, err := io.ReadAll(body)
		logB := l.logBody(mimeType, b)
		if l.base64 {
			dst := make([]byte, 0, 8*len(logB)/6+4)
			dst = base64.StdEncoding.AppendEncode(dst, logB)
			return dst, io.NopCloser(bytes.NewReader(b)), err
		}
		return logB, io.NopCloser(bytes.NewReader(b)), err
	}

	// When output to files,
	// Content-Lengths is not checked to allow logging streaming body.
	if l.bodyPath != "" {
		filePath := l.bodyPath + "body-" + fileName
		f, err := os.OpenFile(filePath, os.O_CREATE|os.O_TRUNC|os.O_RDWR, os.ModePerm)
		if err != nil {
			return nil, body, err
		}
		return []byte("body-" + fileName), &teeReadCloser{r: body, w: f}, nil
	}

	// Handle length == -1 (unknown content length) case
	// We treat this as a streaming body and log it directly into a buffer.
	if length == -1 {
		var buf bytes.Buffer
		tr := io.TeeReader(body, &buf)
		b, err := io.ReadAll(tr)
		if err != nil {
			return nil, body, err
		}
		logB := l.logBody(mimeType, b)
		if l.base64 || isCompressed {
			dst := make([]byte, 0, 8*len(logB)/6+4)
			dst = base64.StdEncoding.AppendEncode(dst, logB)
			return dst, io.NopCloser(bytes.NewReader(b)), err
		}
		// Return the buffer and the function for logging the body
		return logB, io.NopCloser(bytes.NewReader(b)), nil
	}
	return nil, body, nil // No logging.
}

// bodyWriter returns writer that should be used for writing log bodies.
func (l *baseLogger) bodyWriter(fileName, mimeType string, length int64, isCompressed bool) (func() []byte, io.Writer, error) {
	if !slices.Contains(l.mimes, mimeType) {
		return nil, nil, nil
	}

	if length == 0 {
		return nil, nil, nil
	}

	if length > 0 && length < l.maxBody {
		var buf bytes.Buffer
		buf.Grow(int(length))
		bf := func() []byte {
			logB := l.logBody(mimeType, buf.Bytes())
			if l.base64 {
				dst := make([]byte, 0, 8*len(logB)/6+4)
				return base64.StdEncoding.AppendEncode(dst, logB)
			}
			return logB
		}
		return bf, &buf, nil
	}

	// When output to files,
	// Content-Lengths is not checked to allow logging streaming body.
	if l.bodyPath != "" {
		filePath := l.bodyPath + "body-" + fileName
		f, err := os.OpenFile(filePath, os.O_CREATE|os.O_TRUNC|os.O_RDWR, os.ModePerm)
		if err != nil {
			return nil, nil, err
		}
		bf := func() []byte {
			f.Close()
			return []byte("body-" + fileName)
		}
		return bf, f, nil
	}

	// Handle length == -1 (unknown content length) case
	// We treat this as a streaming body and log it directly into a buffer.
	if length == -1 {
		var buf bytes.Buffer
		bf := func() []byte {
			logBody := l.logBody(mimeType, buf.Bytes())
			if l.base64 || isCompressed {
				dst := make([]byte, 0, 8*len(logBody)/6+4)
				return base64.StdEncoding.AppendEncode(dst, logBody)
			}
			return logBody
		}

		// Return the buffer and the function for logging the body
		return bf, &buf, nil
	}
	return nil, nil, nil // No logging.
}

// teeReadCloser reads from r and write
// the read bytes into w.
// teeReadCloser call Closer of both r and w
// when Close method was called.
type teeReadCloser struct {
	r io.ReadCloser
	w io.WriteCloser
}

func (t *teeReadCloser) Close() error {
	_ = t.w.Close()
	return t.r.Close()
}

func (t *teeReadCloser) Read(p []byte) (n int, err error) {
	n, err = t.r.Read(p)
	if n > 0 {
		_, _ = t.w.Write(p[:n])
	}
	if err == io.EOF {
		_ = t.w.Close()
	}
	return n, err
}
