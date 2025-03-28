package httplogger

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"testing"

	v1 "github.com/aileron-gateway/aileron-gateway/apis/core/v1"
	k "github.com/aileron-gateway/aileron-gateway/apis/kernel"
	"github.com/aileron-gateway/aileron-gateway/kernel/er"
	kio "github.com/aileron-gateway/aileron-gateway/kernel/io"
	"github.com/aileron-gateway/aileron-gateway/kernel/log"
	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
	"github.com/aileron-gateway/aileron-gateway/kernel/txtutil"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestNewBaseLogger(t *testing.T) {
	type condition struct {
		spec *v1.LoggingSpec
		lg   log.Logger
	}

	type action struct {
		bl  *baseLogger
		err error
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	testStrRepl, _ := txtutil.NewStringReplacer(&k.ReplacerSpec{
		Replacers: &k.ReplacerSpec_Regexp{
			Regexp: &k.RegexpReplacer{
				Pattern: "[0-9]{5}",
				Replace: "*****",
			},
		},
	})

	testByteRepl, _ := txtutil.NewBytesReplacer(&k.ReplacerSpec{
		Replacers: &k.ReplacerSpec_Regexp{
			Regexp: &k.RegexpReplacer{
				Pattern: "[0-9]{5}",
				Replace: "*****",
			},
		},
	})

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"zero spec",
			[]string{},
			[]string{},
			&condition{
				spec: &v1.LoggingSpec{},
				lg:   log.GlobalLogger(log.DefaultLoggerName),
			},
			&action{
				bl: &baseLogger{
					lg:         log.GlobalLogger(log.DefaultLoggerName),
					w:          os.Stderr,
					queries:    []stringReplFunc{},
					headers:    map[string][]stringReplFunc{},
					headerKeys: []string{},
					bodies:     map[string][]bytesReplFunc{},
				},
			},
		),
		gen(
			"valid header spec",
			[]string{},
			[]string{},
			&condition{
				spec: &v1.LoggingSpec{
					Headers: []*v1.LogHeaderSpec{
						{
							Name: "Foo",
							Replacers: []*k.ReplacerSpec{
								{
									Replacers: &k.ReplacerSpec_Regexp{
										Regexp: &k.RegexpReplacer{Pattern: "[0-9]{5}", Replace: "*****"},
									},
								},
							},
						},
					},
				},
				lg: log.GlobalLogger(log.DefaultLoggerName),
			},
			&action{
				bl: &baseLogger{
					lg:      log.GlobalLogger(log.DefaultLoggerName),
					w:       os.Stderr,
					queries: []stringReplFunc{},
					headers: map[string][]stringReplFunc{
						"Foo": {testStrRepl.Replace},
					},
					headerKeys: []string{"Foo"},
					bodies:     map[string][]bytesReplFunc{},
				},
			},
		),
		gen(
			"invalid header spec",
			[]string{},
			[]string{},
			&condition{
				spec: &v1.LoggingSpec{
					Headers: []*v1.LogHeaderSpec{
						{
							Name: "Foo",
							Replacers: []*k.ReplacerSpec{
								{
									Replacers: &k.ReplacerSpec_Regexp{
										Regexp: &k.RegexpReplacer{Pattern: "[0-9", Replace: "*****"},
									},
								},
							},
						},
					},
				},
				lg: log.GlobalLogger(log.DefaultLoggerName),
			},
			&action{
				bl: nil,
				err: &er.Error{
					Package:     txtutil.ErrPkg,
					Type:        txtutil.ErrTypeReplacer,
					Description: txtutil.ErrDscPattern,
				},
			},
		),
		gen(
			"valid body spec",
			[]string{},
			[]string{},
			&condition{
				spec: &v1.LoggingSpec{
					Bodies: []*v1.LogBodySpec{
						{
							Mime: "application/json",
							Replacers: []*k.ReplacerSpec{
								{
									Replacers: &k.ReplacerSpec_Regexp{
										Regexp: &k.RegexpReplacer{Pattern: "[0-9]{5}", Replace: "*****"},
									},
								},
							},
						},
					},
				},
				lg: log.GlobalLogger(log.DefaultLoggerName),
			},
			&action{
				bl: &baseLogger{
					lg:         log.GlobalLogger(log.DefaultLoggerName),
					w:          os.Stderr,
					queries:    []stringReplFunc{},
					headers:    map[string][]stringReplFunc{},
					headerKeys: []string{},
					bodies: map[string][]bytesReplFunc{
						"application/json": {testByteRepl.Replace},
					},
				},
			},
		),
		gen(
			"invalid body spec",
			[]string{},
			[]string{},
			&condition{
				spec: &v1.LoggingSpec{
					Bodies: []*v1.LogBodySpec{
						{
							Mime: "application/json",
							Replacers: []*k.ReplacerSpec{
								{
									Replacers: &k.ReplacerSpec_Regexp{
										Regexp: &k.RegexpReplacer{Pattern: "[0-9", Replace: "*****"},
									},
								},
							},
						},
					},
				},
				lg: log.GlobalLogger(log.DefaultLoggerName),
			},
			&action{
				bl: nil,
				err: &er.Error{
					Package:     txtutil.ErrPkg,
					Type:        txtutil.ErrTypeReplacer,
					Description: txtutil.ErrDscPattern,
				},
			},
		),
		gen(
			"template with writer logger",
			[]string{},
			[]string{},
			&condition{
				spec: &v1.LoggingSpec{
					LogFormat: "%foo%",
				},
				lg: log.GlobalLogger(log.DefaultLoggerName),
			},
			&action{
				bl: &baseLogger{
					lg:         log.GlobalLogger(log.DefaultLoggerName),
					w:          log.GlobalLogger(log.DefaultLoggerName).(io.Writer),
					tpl:        &txtutil.FastTemplate{},
					queries:    []stringReplFunc{},
					headers:    map[string][]stringReplFunc{},
					headerKeys: []string{},
					bodies:     map[string][]bytesReplFunc{},
				},
			},
		),
		gen(
			"template with non writer logger",
			[]string{},
			[]string{},
			&condition{
				spec: &v1.LoggingSpec{
					LogFormat: "%foo%",
				},
				lg: struct{ log.Logger }{&log.Noop{}},
			},
			&action{
				bl: nil,
				err: &er.Error{
					Package:     "httplogger",
					Type:        "base logger",
					Description: "formatted log requires logger with io.Writer interface",
				},
			},
		),
		gen(
			"valid body output path",
			[]string{},
			[]string{},
			&condition{
				spec: &v1.LoggingSpec{
					BodyOutputPath: "./tmp/",
				},
				lg: log.GlobalLogger(log.DefaultLoggerName),
			},
			&action{
				bl: &baseLogger{
					lg:         log.GlobalLogger(log.DefaultLoggerName),
					w:          os.Stderr,
					queries:    []stringReplFunc{},
					headers:    map[string][]stringReplFunc{},
					headerKeys: []string{},
					bodies:     map[string][]bytesReplFunc{},
					bodyPath:   "tmp/",
				},
			},
		),
		gen(
			"invalid body output path",
			[]string{},
			[]string{},
			&condition{
				spec: &v1.LoggingSpec{
					BodyOutputPath: "\n\r\t\x00",
				},
				lg: log.GlobalLogger(log.DefaultLoggerName),
			},
			&action{
				err: &er.Error{
					Package:     kio.ErrPkg,
					Type:        kio.ErrTypeFile,
					Description: kio.ErrDscFileSys,
				},
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			if tt.C().spec.BodyOutputPath != "" {
				defer os.Remove(tt.C().spec.BodyOutputPath)
			}

			bl, err := newBaseLogger(tt.C().spec, tt.C().lg)
			testutil.Diff(t, tt.A().err, err, cmpopts.EquateErrors())

			opts := []cmp.Option{
				cmp.AllowUnexported(baseLogger{}),
				cmpopts.IgnoreTypes(txtutil.FastTemplate{}),
				cmp.Comparer(testutil.ComparePointer[io.Writer]),
				cmp.Comparer(testutil.ComparePointer[stringReplFunc]),
				cmp.Comparer(testutil.ComparePointer[bytesReplFunc]),
			}
			testutil.Diff(t, tt.A().bl, bl, opts...)
		})
	}
}

type noopReplacer[T any] struct{}

func (r *noopReplacer[T]) Replace(in T) T {
	return in
}

func TestBaseLogger_logOutput(t *testing.T) {
	type condition struct {
		bl      *baseLogger
		level   slog.Level
		attrs   []any
		tagFunc func(string) []byte
	}

	type action struct {
		check []string
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	testAttrs := &requestAttrs{
		typ:    "test-type",
		id:     "test-id",
		time:   "test-time",
		host:   "test-host",
		method: "test-method",
		path:   "test-path",
		query:  "test-query",
	}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"output by logger",
			[]string{},
			[]string{},
			&condition{
				bl:      &baseLogger{},
				level:   slog.LevelInfo,
				attrs:   testAttrs.accessKeyValues(),
				tagFunc: testAttrs.TagFunc,
			},
			&action{
				check: []string{
					`"id":"test-id"`,
					`"time":"test-time"`,
				},
			},
		),
		gen(
			"output by writer",
			[]string{},
			[]string{},
			&condition{
				bl: &baseLogger{
					tpl: txtutil.NewFastTemplate("id=%id% time=%time%", "%", "%"),
				},
				level:   slog.LevelInfo,
				attrs:   testAttrs.accessKeyValues(),
				tagFunc: testAttrs.TagFunc,
			},
			&action{
				check: []string{
					`id=test-id`,
					`time=test-time`,
				},
			},
		),
		gen(
			"not output",
			[]string{},
			[]string{},
			&condition{
				bl:      &baseLogger{},
				level:   slog.LevelWarn,
				attrs:   testAttrs.accessKeyValues(),
				tagFunc: testAttrs.TagFunc,
			},
			&action{
				check: []string{},
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			var buf bytes.Buffer
			tt.C().bl.lg = log.NewJSONSLogger(&buf, &slog.HandlerOptions{
				Level: tt.C().level,
			})
			tt.C().bl.w = &buf

			ctx := context.Background()
			tt.C().bl.logOutput(ctx, "test-msg", tt.C().attrs, tt.C().tagFunc)
			for _, s := range tt.A().check {
				testutil.Diff(t, true, strings.Contains(buf.String(), s))
			}
		})
	}
}

func TestBaseLogger_logQuery(t *testing.T) {
	type condition struct {
		bl *baseLogger
		q  string
	}

	type action struct {
		q string
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	// testReplacer is the replacer for testing.
	// "test-foobar" will be "***"
	s := &k.ReplacerSpec{Replacers: &k.ReplacerSpec_Regexp{
		Regexp: &k.RegexpReplacer{Pattern: `test-[^&]*`, Replace: `***`},
	}}
	testReplacer, _ := txtutil.NewStringReplacer(s)

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"zero logger",
			[]string{},
			[]string{},
			&condition{
				bl: &baseLogger{},
				q:  "foo=test-foo&bar=test-bar",
			},
			&action{
				q: "foo=test-foo&bar=test-bar",
			},
		),
		gen(
			"specify one query",
			[]string{},
			[]string{},
			&condition{
				bl: &baseLogger{
					queries: []stringReplFunc{
						(&noopReplacer[string]{}).Replace,
					},
				},
				q: "foo=test-foo&bar=test-bar",
			},
			&action{
				q: "foo=test-foo&bar=test-bar",
			},
		),
		gen(
			"replace",
			[]string{},
			[]string{},
			&condition{
				bl: &baseLogger{
					queries: []stringReplFunc{
						testReplacer.Replace,
					},
				},
				q: "foo=test-foo&bar=test-bar",
			},
			&action{
				q: "foo=***&bar=***",
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			q := tt.C().bl.logQuery(tt.C().q)
			testutil.Diff(t, tt.A().q, q)
		})
	}
}

func TestBaseLogger_logHeaders(t *testing.T) {
	type condition struct {
		h         http.Header
		all       bool
		replacers map[string][]stringReplFunc
	}

	type action struct {
		h map[string]string
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	s := &k.ReplacerSpec{Replacers: &k.ReplacerSpec_Regexp{
		Regexp: &k.RegexpReplacer{Pattern: `test-.*`, Replace: `***`},
	}}
	testReplacer, _ := txtutil.NewStringReplacer(s)

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"no replacer",
			[]string{},
			[]string{},
			&condition{
				h: http.Header{
					"Foo": []string{"test-foo"},
					"Bar": []string{"test-bar1", "test-bar2"},
				},
				all:       false,
				replacers: nil,
			},
			&action{
				h: map[string]string{},
			},
		),
		gen(
			"noop replace for single value",
			[]string{},
			[]string{},
			&condition{
				h: http.Header{
					"Foo": []string{"test-foo"},
					"Bar": []string{"test-bar1", "test-bar2"},
				},
				all: false,
				replacers: map[string][]stringReplFunc{
					"Foo": {(&noopReplacer[string]{}).Replace},
				},
			},
			&action{
				h: map[string]string{
					"Foo": "test-foo",
				},
			},
		),
		gen(
			"noop replace for multiple values",
			[]string{},
			[]string{},
			&condition{
				h: http.Header{
					"Foo": []string{"test-foo"},
					"Bar": []string{"test-bar1", "test-bar2"},
				},
				all: false,
				replacers: map[string][]stringReplFunc{
					"Bar": {(&noopReplacer[string]{}).Replace},
				},
			},
			&action{
				h: map[string]string{
					"Bar": "test-bar1,test-bar2",
				},
			},
		),
		gen(
			"noop replace for non-exist value",
			[]string{},
			[]string{},
			&condition{
				h: http.Header{
					"Foo": []string{"test-foo"},
					"Bar": []string{"test-bar1", "test-bar2"},
				},
				all: false,
				replacers: map[string][]stringReplFunc{
					"Baz": {(&noopReplacer[string]{}).Replace},
				},
			},
			&action{
				h: map[string]string{},
			},
		),
		gen(
			"replace for single value",
			[]string{},
			[]string{},
			&condition{
				h: http.Header{
					"Foo": []string{"test-foo"},
					"Bar": []string{"test-bar1", "test-bar2"},
				},
				all: false,
				replacers: map[string][]stringReplFunc{
					"Foo": {testReplacer.Replace},
				},
			},
			&action{
				h: map[string]string{
					"Foo": "***",
				},
			},
		),
		gen(
			"replace for multiple values",
			[]string{},
			[]string{},
			&condition{
				h: http.Header{
					"Foo": []string{"test-foo"},
					"Bar": []string{"test-bar1", "test-bar2"},
				},
				all: false,
				replacers: map[string][]stringReplFunc{
					"Bar": {testReplacer.Replace},
				},
			},
			&action{
				h: map[string]string{
					"Bar": "***",
				},
			},
		),
		gen(
			"replace for non-exist value",
			[]string{},
			[]string{},
			&condition{
				h: http.Header{
					"Foo": []string{"test-foo"},
					"Bar": []string{"test-bar1", "test-bar2"},
				},
				all: false,
				replacers: map[string][]stringReplFunc{
					"Baz": {testReplacer.Replace},
				},
			},
			&action{
				h: map[string]string{},
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			keys := []string{}
			for k := range tt.C().replacers {
				keys = append(keys, k)
			}
			bl := &baseLogger{
				allHeaders: tt.C().all,
				headers:    tt.C().replacers,
				headerKeys: keys,
			}
			h := bl.logHeaders(tt.C().h)
			testutil.Diff(t, tt.A().h, h)
		})
	}
}

func TestBaseLogger_logBody(t *testing.T) {
	type condition struct {
		bl       *baseLogger
		mimeType string
		body     string
	}

	type action struct {
		body string
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	s := &k.ReplacerSpec{Replacers: &k.ReplacerSpec_Regexp{
		Regexp: &k.RegexpReplacer{Pattern: `test-[^"]*`, Replace: `***`},
	}}
	testReplacer, _ := txtutil.NewBytesReplacer(s)

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"zero logger",
			[]string{},
			[]string{},
			&condition{
				bl:       &baseLogger{},
				mimeType: "",
				body:     "",
			},
			&action{
				body: "",
			},
		),
		gen(
			"body without mask",
			[]string{},
			[]string{},
			&condition{
				bl:       &baseLogger{},
				mimeType: "application/json",
				body:     `{"foo":"bar"}`,
			},
			&action{
				body: `{"foo":"bar"}`,
			},
		),
		gen(
			"body with mask",
			[]string{},
			[]string{},
			&condition{
				bl: &baseLogger{
					bodies: map[string][]bytesReplFunc{
						"application/json": {
							(&jsonFieldReplacer{
								fields:    []string{"foo"},
								replacers: []bytesReplFunc{testReplacer.Replace},
							}).Replace,
						},
					},
				},
				mimeType: "application/json",
				body:     `{"foo":"test-bar"}`,
			},
			&action{
				body: `{"foo":"***"}`,
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			body := tt.C().bl.logBody(tt.C().mimeType, []byte(tt.C().body))
			testutil.Diff(t, tt.A().body, string(body))
		})
	}
}

type errorReadCloser struct{}

func (e *errorReadCloser) Read(p []byte) (n int, err error) {
	return 0, errors.New("read error")
}

func (e *errorReadCloser) Close() error {
	return nil
}

func TestBaseLogger_bodyReadCloser(t *testing.T) {
	type condition struct {
		bl           *baseLogger
		fileName     string
		mimeType     string
		length       int64
		body         io.ReadCloser
		isCompressed bool
	}

	type action struct {
		b         string
		read      string
		nonNilErr bool
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"non target mime",
			[]string{},
			[]string{},
			&condition{
				bl: &baseLogger{
					mimes: []string{"application/json"},
				},
				mimeType: "text/plain",
			},
			&action{},
		),
		gen(
			"zero length body",
			[]string{},
			[]string{},
			&condition{
				bl: &baseLogger{
					mimes: []string{"application/json"},
				},
				mimeType: "application/json",
				length:   0,
			},
			&action{},
		),
		gen(
			"buffer reader wo b64",
			[]string{},
			[]string{},
			&condition{
				bl: &baseLogger{
					mimes:   []string{"application/json"},
					maxBody: 100,
				},
				mimeType: "application/json",
				length:   9,
				body:     io.NopCloser(bytes.NewBuffer([]byte("test-body"))),
			},
			&action{
				b:    "test-body",
				read: "test-body",
			},
		),
		gen(
			"buffer reader w b64",
			[]string{},
			[]string{},
			&condition{
				bl: &baseLogger{
					mimes:   []string{"application/json"},
					maxBody: 100,
					base64:  true,
				},
				mimeType: "application/json",
				length:   9,
				body:     io.NopCloser(bytes.NewBuffer([]byte("test-body"))),
			},
			&action{
				b:    base64.StdEncoding.EncodeToString([]byte("test-body")),
				read: "test-body",
			},
		),
		gen(
			"file reader",
			[]string{},
			[]string{},
			&condition{
				bl: &baseLogger{
					mimes:    []string{"application/json"},
					bodyPath: "./",
				},
				mimeType: "application/json",
				fileName: "test.txt",
				length:   9,
				body:     io.NopCloser(bytes.NewBuffer([]byte("test-body"))),
			},
			&action{
				b:    "body-test.txt",
				read: "test-body",
			},
		),
		gen(
			"file create error",
			[]string{},
			[]string{},
			&condition{
				bl: &baseLogger{
					mimes:    []string{"application/json"},
					bodyPath: "./not-exists/",
				},
				mimeType: "application/json",
				fileName: "test.txt",
				length:   9,
				body:     io.NopCloser(bytes.NewBuffer([]byte("test-body"))),
			},
			&action{
				b:         "",
				read:      "test-body",
				nonNilErr: true,
			},
		),
		gen(
			"no logging",
			[]string{},
			[]string{},
			&condition{
				bl: &baseLogger{
					mimes:   []string{"application/json"},
					maxBody: 1,
				},
				mimeType: "application/json",
				length:   9,
				body:     io.NopCloser(bytes.NewBuffer([]byte("test-body"))),
			},
			&action{
				b:    "",
				read: "test-body",
			},
		),
		gen(
			"streaming body(length = -1) without base64",
			[]string{},
			[]string{},
			&condition{
				bl: &baseLogger{
					mimes:   []string{"application/json"},
					maxBody: 100,
				},
				mimeType: "application/json",
				length:   -1, // Content-Length unknown
				body:     io.NopCloser(bytes.NewBuffer([]byte("streaming-body"))),
			},
			&action{
				b:    "streaming-body",
				read: "streaming-body",
			},
		),
		gen(
			"streaming body(length = -1) with base64",
			[]string{},
			[]string{},
			&condition{
				bl: &baseLogger{
					mimes:   []string{"application/json"},
					maxBody: 100,
					base64:  true,
				},
				mimeType: "application/json",
				length:   -1, // Content-Length unknown
				body:     io.NopCloser(bytes.NewBuffer([]byte("streaming-body"))),
			},
			&action{
				b:    base64.StdEncoding.EncodeToString([]byte("streaming-body")),
				read: "streaming-body",
			},
		),
		gen(
			"Compressed request with streaming body(length = -1)",
			[]string{},
			[]string{},
			&condition{
				bl: &baseLogger{
					mimes:   []string{"application/json"},
					maxBody: 100,
					base64:  true,
				},
				mimeType:     "application/json",
				length:       -1, // Content-Length unknown
				body:         io.NopCloser(bytes.NewBuffer([]byte("streaming-body"))),
				isCompressed: true,
			},
			&action{
				b:    base64.StdEncoding.EncodeToString([]byte("streaming-body")),
				read: "streaming-body",
			},
		),
		gen(
			"streaming body(length = -1) with read error",
			[]string{},
			[]string{},
			&condition{
				bl: &baseLogger{
					mimes:   []string{"application/json"},
					maxBody: 100,
				},
				mimeType: "application/json",
				length:   -1, // Content-Length unknown
				body:     &errorReadCloser{},
			},
			&action{
				b:         "",
				read:      "",
				nonNilErr: true,
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			if tt.C().fileName != "" {
				defer os.Remove("./body-" + tt.C().fileName)
			}

			fn := tt.C().fileName
			mt := tt.C().mimeType
			size := tt.C().length
			b, rc, err := tt.C().bl.bodyReadCloser(fn, mt, size, tt.C().body, tt.C().isCompressed)
			testutil.Diff(t, tt.A().b, string(b))
			testutil.Diff(t, tt.A().nonNilErr, err != nil)
			if tt.A().read != "" {
				read, err := io.ReadAll(rc)
				testutil.Diff(t, nil, err)
				testutil.Diff(t, tt.A().read, string(read))
				rc.Close()
			}
		})
	}
}

func TestBaseLogger_bodyWriter(t *testing.T) {
	type condition struct {
		bl           *baseLogger
		fileName     string
		mimeType     string
		length       int64
		isCompressed bool
	}

	type action struct {
		write     string
		read      string
		nonNilErr bool
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"non target mime",
			[]string{},
			[]string{},
			&condition{
				bl: &baseLogger{
					mimes: []string{"application/json"},
				},
				mimeType: "text/plain",
			},
			&action{},
		),
		gen(
			"zero length body",
			[]string{},
			[]string{},
			&condition{
				bl: &baseLogger{
					mimes: []string{"application/json"},
				},
				mimeType: "application/json",
				length:   0,
			},
			&action{},
		),
		gen(
			"buffer reader wo b64",
			[]string{},
			[]string{},
			&condition{
				bl: &baseLogger{
					mimes:   []string{"application/json"},
					maxBody: 100,
				},
				mimeType: "application/json",
				length:   9,
			},
			&action{
				write: "test-body",
				read:  "test-body",
			},
		),
		gen(
			"buffer reader w b64",
			[]string{},
			[]string{},
			&condition{
				bl: &baseLogger{
					mimes:   []string{"application/json"},
					maxBody: 100,
					base64:  true,
				},
				mimeType: "application/json",
				length:   9,
			},
			&action{
				write: "test-body",
				read:  base64.StdEncoding.EncodeToString([]byte("test-body")),
			},
		),
		gen(
			"file reader",
			[]string{},
			[]string{},
			&condition{
				bl: &baseLogger{
					mimes:    []string{"application/json"},
					bodyPath: "./",
				},
				mimeType: "application/json",
				fileName: "test.txt",
				length:   9,
			},
			&action{
				write: "test-body",
				read:  "body-test.txt",
			},
		),
		gen(
			"file create error",
			[]string{},
			[]string{},
			&condition{
				bl: &baseLogger{
					mimes:    []string{"application/json"},
					bodyPath: "./not-exists/",
				},
				mimeType: "application/json",
				fileName: "test.txt",
				length:   9,
			},
			&action{
				read:      "test-body",
				nonNilErr: true,
			},
		),
		gen(
			"no logging",
			[]string{},
			[]string{},
			&condition{
				bl: &baseLogger{
					mimes:   []string{"application/json"},
					maxBody: 1,
				},
				mimeType: "application/json",
				length:   9,
			},
			&action{
				read: "",
			},
		),
		gen(
			"unknown length body(length = -1)",
			[]string{},
			[]string{},
			&condition{
				bl: &baseLogger{
					mimes:   []string{"application/json"},
					maxBody: 100,
				},
				mimeType: "application/json",
				length:   -1, // Content-Length unknown
			},
			&action{
				write: "test-body",
				read:  "test-body",
			},
		),
		gen(
			"unknown length body(length = -1) with base64",
			[]string{},
			[]string{},
			&condition{
				bl: &baseLogger{
					mimes:   []string{"application/json"},
					maxBody: 100,
					base64:  true,
				},
				mimeType: "application/json",
				length:   -1, // Content-Length unknown
			},
			&action{
				write: "test-body",
				read:  base64.StdEncoding.EncodeToString([]byte("test-body")),
			},
		),
		gen(
			"Compressed request with unknown body length (-1)",
			[]string{},
			[]string{},
			&condition{
				bl: &baseLogger{
					mimes:   []string{"application/json"},
					maxBody: 100,
				},
				mimeType:     "application/json",
				length:       -1, // Content-Length unknown
				isCompressed: true,
			},
			&action{
				write: "test-body",
				read:  base64.StdEncoding.EncodeToString([]byte("test-body")),
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			if tt.C().fileName != "" {
				defer os.Remove("./body-" + tt.C().fileName)
			}

			fn := tt.C().fileName
			mt := tt.C().mimeType
			size := tt.C().length
			bf, w, err := tt.C().bl.bodyWriter(fn, mt, size, tt.C().isCompressed)
			if tt.A().nonNilErr {
				testutil.Diff(t, true, err != nil)
				testutil.Diff(t, (func() []byte)(nil), bf)
				testutil.Diff(t, nil, w)
				return
			}
			if tt.A().read != "" {
				w.Write([]byte(tt.A().write))
				read := bf()
				testutil.Diff(t, tt.A().read, string(read))
			} else {
				testutil.Diff(t, (func() []byte)(nil), bf)
				testutil.Diff(t, nil, w)
			}
		})
	}
}
