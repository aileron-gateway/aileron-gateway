package jaegertracer

import (
	"context"
	"crypto/tls"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aileron-gateway/aileron-gateway/core"
	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/mocktracer"
)

func TestMiddleware(t *testing.T) {
	type condition struct {
		mNames     map[int]string
		mCtxKey    any
		headers    []string
		parentSpan bool
		childSpan  bool
		httpsFlag  bool
	}

	type action struct {
		statusCode    int
		tags          map[string]any
		body          string
		operationName string
	}

	CndSetNoJaegerTracer := "set no jaegerTracer"
	CndSetmCtxKey := "set mCtxKey"
	CndSetmNames := "set mNames"
	CndSetStartChildSpan := "start ChildSpan"
	CndSetParentSpan := "set ParentSpan"
	CndSetHeaders := "set Headers"
	CndSetHTTPSSchema := "set HTTPS Schema"
	ActCheckExpected := "expected value returned"

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	tb.Condition(CndSetNoJaegerTracer, "set no jaegerTracer")
	tb.Condition(CndSetmCtxKey, "set mCtxKey")
	tb.Condition(CndSetmNames, "set mNames")
	tb.Condition(CndSetStartChildSpan, "start ChildSpan")
	tb.Condition(CndSetParentSpan, "set ParentSpan")
	tb.Condition(CndSetHeaders, "set Headers")
	tb.Condition(CndSetHTTPSSchema, "set HTTPS Schema")
	tb.Action(ActCheckExpected, "check that an expected value returned")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"nil tracer",
			[]string{CndSetNoJaegerTracer},
			[]string{ActCheckExpected},
			&condition{},
			&action{
				statusCode: http.StatusOK,
				tags: map[string]any{
					"http.id":          "",
					"http.method":      "GET",
					"http.path":        "/",
					"http.query":       "",
					"net.addr":         "192.0.2.1:1234",
					"net.host":         "example.com",
					"caller.file":      "jaegertracer/tracer_test.go",
					"caller.func":      "net/http.HandlerFunc.ServeHTTP",
					"http.status_code": 0,
					"http.schema":      "http",
				},
				operationName: "0:middleware",
			},
		),
		gen(
			"set mCtxKey",
			[]string{CndSetmCtxKey},
			[]string{ActCheckExpected},
			&condition{
				mCtxKey: 1,
			},
			&action{
				statusCode: http.StatusOK,
				tags: map[string]any{
					"caller.file": "jaegertracer/tracer_test.go",
					"caller.func": "net/http.HandlerFunc.ServeHTTP",
				},
				operationName: "2:middleware",
			},
		),
		gen(
			"start ChildSpan",
			[]string{CndSetStartChildSpan},
			[]string{ActCheckExpected},
			&condition{
				childSpan: true,
			},
			&action{
				statusCode: http.StatusOK,
				tags: map[string]any{
					"http.id":          "",
					"http.method":      "GET",
					"http.path":        "/",
					"http.query":       "",
					"net.addr":         "192.0.2.1:1234",
					"net.host":         "example.com",
					"caller.file":      "jaegertracer/tracer_test.go",
					"caller.func":      "net/http.HandlerFunc.ServeHTTP",
					"http.status_code": 0,
					"http.schema":      "http",
				},
				operationName: "0:middleware",
			},
		),
		gen(
			"set ParentSpan",
			[]string{CndSetParentSpan},
			[]string{ActCheckExpected},
			&condition{
				parentSpan: true,
			},
			&action{
				statusCode: http.StatusOK,
				tags: map[string]any{
					"http.id":          "",
					"http.method":      "GET",
					"http.path":        "/",
					"http.query":       "",
					"net.addr":         "192.0.2.1:1234",
					"net.host":         "example.com",
					"caller.file":      "jaegertracer/tracer_test.go",
					"caller.func":      "net/http.HandlerFunc.ServeHTTP",
					"http.status_code": 0,
					"http.schema":      "http",
				},
				operationName: "0:middleware",
			},
		),
		gen(
			"set mNames",
			[]string{CndSetmNames},
			[]string{ActCheckExpected},
			&condition{
				mNames: map[int]string{
					0: "testName",
				},
			},
			&action{
				statusCode: http.StatusOK,
				tags: map[string]any{
					"http.id":          "",
					"http.method":      "GET",
					"http.path":        "/",
					"http.query":       "",
					"net.addr":         "192.0.2.1:1234",
					"net.host":         "example.com",
					"caller.file":      "jaegertracer/tracer_test.go",
					"caller.func":      "net/http.HandlerFunc.ServeHTTP",
					"http.status_code": 0,
					"http.schema":      "http",
				},
				operationName: "testName",
			},
		),
		gen(
			"set Headers",
			[]string{CndSetHeaders},
			[]string{ActCheckExpected},
			&condition{
				headers: []string{"testHeader"},
			},
			&action{
				statusCode: http.StatusOK,
				tags: map[string]any{
					"http.id":                "",
					"http.method":            "GET",
					"http.path":              "/",
					"http.query":             "",
					"net.addr":               "192.0.2.1:1234",
					"net.host":               "example.com",
					"caller.file":            "jaegertracer/tracer_test.go",
					"caller.func":            "net/http.HandlerFunc.ServeHTTP",
					"http.status_code":       0,
					"http.schema":            "http",
					"http.header.testheader": []string(nil),
				},
				operationName: "0:middleware",
			},
		),
		gen(
			"set HTTPS schema",
			[]string{CndSetHTTPSSchema},
			[]string{ActCheckExpected},
			&condition{
				httpsFlag: true,
			},
			&action{
				statusCode: http.StatusOK,
				tags: map[string]any{
					"http.id":          "",
					"http.method":      "GET",
					"http.path":        "/",
					"http.query":       "",
					"net.addr":         "192.0.2.1:1234",
					"net.host":         "example.com",
					"caller.file":      "jaegertracer/tracer_test.go",
					"caller.func":      "net/http.HandlerFunc.ServeHTTP",
					"http.status_code": 0,
					"http.schema":      "https",
				},
				operationName: "0:middleware",
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			mockTracer := mocktracer.New()
			jt := &jaegerTracer{
				tracer:  mockTracer,
				mNames:  tt.C().mNames,
				headers: tt.C().headers,
			}

			ctx := context.Background()

			if tt.C().parentSpan {
				parentSpan := jt.tracer.StartSpan("parent")
				defer parentSpan.Finish()

				ctx = opentracing.ContextWithSpan(ctx, parentSpan)
			} else {
				ctx = context.WithValue(ctx, mCtxKey, tt.C().mCtxKey)
			}

			req := httptest.NewRequest(http.MethodGet, "/", nil)

			if tt.C().childSpan {
				childSpan := mockTracer.StartSpan("testSpan")
				defer childSpan.Finish()

				jt.tracer.Inject(
					childSpan.Context(),
					opentracing.HTTPHeaders,
					opentracing.HTTPHeadersCarrier(req.Header),
				)
			}

			if tt.C().httpsFlag {
				req.TLS = &tls.ConnectionState{}
			}

			req = req.WithContext(ctx)
			resp := httptest.NewRecorder()
			h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
			jt.Middleware(h).ServeHTTP(resp, req)

			span := mockTracer.FinishedSpans()[0]

			testutil.Diff(t, tt.A().tags, span.Tags())
			testutil.Diff(t, tt.A().operationName, span.OperationName)
			testutil.Diff(t, tt.A().statusCode, resp.Code)
			testutil.Diff(t, tt.A().body, resp.Body.String())
		})
	}
}

func TestTripperware(t *testing.T) {
	type condition struct {
		tNames       map[int]string
		tCtxKey      any
		headers      []string
		parentSpan   bool
		roundTripErr bool
	}

	type action struct {
		err           any // error or errorutil.Kind
		statusCode    int
		tags          map[string]any
		operationName string
	}

	CndSetNoJaegerTracer := "set no jaegerTracer"
	CndSettCtxKey := "set tCtxKey"
	CndSettNames := "set tNames"
	CndSetHeaders := "set Headers"
	CndSetRoundTripErr := "cause RoundTrip Err"
	ActCheckExpected := "expected value returned"

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	tb.Condition(CndSetNoJaegerTracer, "set no jaegerTracer")
	tb.Condition(CndSettCtxKey, "set tCtxKey")
	tb.Condition(CndSettNames, "set tNames")
	tb.Condition(CndSetHeaders, "set Headers")
	tb.Condition(CndSetRoundTripErr, "cause RoundTrip Err")
	tb.Action(ActCheckExpected, "check that an expected value returned")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"nil tracer",
			[]string{CndSetNoJaegerTracer},
			[]string{ActCheckExpected},
			&condition{},
			&action{
				statusCode: http.StatusOK,
				tags: map[string]any{
					"http.id":          "",
					"http.method":      "GET",
					"http.path":        "/",
					"http.query":       "",
					"caller.file":      "jaegertracer/tracer_test.go",
					"caller.func":      "github.com/aileron-gateway/aileron-gateway/core.RoundTripperFunc.RoundTrip",
					"peer.host":        "",
					"http.status_code": 200,
					"http.schema":      "http",
				},
				operationName: "0:tripperware",
			},
		),
		gen(
			"set tCtxKey",
			[]string{CndSettCtxKey},
			[]string{ActCheckExpected},
			&condition{
				tCtxKey: 1,
			},
			&action{
				statusCode: http.StatusOK,
				tags: map[string]any{
					"caller.file": "jaegertracer/tracer_test.go",
					"caller.func": "github.com/aileron-gateway/aileron-gateway/core.RoundTripperFunc.RoundTrip",
				},
				operationName: "2:tripperware",
			},
		),
		gen(
			"set tNames",
			[]string{CndSettNames},
			[]string{ActCheckExpected},
			&condition{
				tNames: map[int]string{
					0: "testName",
				},
			},
			&action{
				statusCode: http.StatusOK,
				tags: map[string]any{
					"http.id":          "",
					"http.method":      "GET",
					"http.path":        "/",
					"http.query":       "",
					"caller.file":      "jaegertracer/tracer_test.go",
					"caller.func":      "github.com/aileron-gateway/aileron-gateway/core.RoundTripperFunc.RoundTrip",
					"peer.host":        "",
					"http.status_code": 200,
					"http.schema":      "http",
				},
				operationName: "testName",
			},
		),
		gen(
			"set Headers",
			[]string{CndSetHeaders},
			[]string{ActCheckExpected},
			&condition{
				headers: []string{"testHeader"},
			},
			&action{
				statusCode: http.StatusOK,
				tags: map[string]any{
					"http.id":                "",
					"http.method":            "GET",
					"http.path":              "/",
					"http.query":             "",
					"caller.file":            "jaegertracer/tracer_test.go",
					"caller.func":            "github.com/aileron-gateway/aileron-gateway/core.RoundTripperFunc.RoundTrip",
					"peer.host":              "",
					"http.status_code":       200,
					"http.schema":            "http",
					"http.header.testheader": []string(nil),
				},
				operationName: "0:tripperware",
			},
		),
		gen(
			"cause RoundTripError",
			[]string{CndSetRoundTripErr},
			[]string{ActCheckExpected},
			&condition{
				roundTripErr: true,
			},
			&action{
				statusCode: http.StatusOK,
				tags: map[string]any{
					"http.id":          "",
					"http.method":      "GET",
					"http.path":        "/",
					"http.query":       "",
					"caller.file":      "jaegertracer/tracer_test.go",
					"caller.func":      "github.com/aileron-gateway/aileron-gateway/core.RoundTripperFunc.RoundTrip",
					"peer.host":        "",
					"http.status_code": 0,
					"http.schema":      "http",
				},
				operationName: "0:tripperware",
				err:           opentracing.ErrSpanContextNotFound,
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			mockTracer := mocktracer.New()
			jt := &jaegerTracer{
				tracer:  mockTracer,
				tNames:  tt.C().tNames,
				headers: tt.C().headers,
			}

			ctx := context.Background()

			if tt.C().parentSpan {
				parentSpan := jt.tracer.StartSpan("parent")
				defer parentSpan.Finish()

				ctx = opentracing.ContextWithSpan(ctx, parentSpan)
			} else {
				ctx = context.WithValue(ctx, tCtxKey, tt.C().tCtxKey)
			}

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req = req.WithContext(ctx)

			r := core.RoundTripperFunc(func(r *http.Request) (*http.Response, error) {
				resp := &http.Response{
					StatusCode: http.StatusOK,
					Status:     "StatusOK",
					Proto:      "HTTP/1.0",
					ProtoMajor: 1,
					ProtoMinor: 0,
					Body:       nil,
				}
				if tt.C().roundTripErr {
					return resp, opentracing.ErrSpanContextNotFound
				}
				return resp, nil
			})

			opts := []cmp.Option{
				cmpopts.EquateErrors(),
			}

			resp, err := jt.Tripperware(r).RoundTrip(req)
			testutil.Diff(t, tt.A().err, err, opts...)

			span := mockTracer.FinishedSpans()[0]

			testutil.Diff(t, tt.A().tags, span.Tags())
			testutil.Diff(t, tt.A().operationName, span.OperationName)
			testutil.Diff(t, tt.A().statusCode, resp.StatusCode)
		})
	}
}

func TestTrace(t *testing.T) {
	type condition struct {
		name       string
		tags       map[string]string
		parentSpan bool
	}

	type action struct {
		name string
		tags map[string]any
	}

	CndSetEmptyNameTags := "set empty name and tags"
	CndSetName := "set name"
	CndSetSingleTag := "set single tag"
	CndSetMultipleTags := "set multiple tags"
	CndSetParentSpan := "set parent span"
	ActCheckExpected := "expected value returned"

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	tb.Condition(CndSetEmptyNameTags, "set empty name and tags")
	tb.Condition(CndSetName, "set name")
	tb.Condition(CndSetSingleTag, "set single tag")
	tb.Condition(CndSetMultipleTags, "set multiple tags")
	tb.Condition(CndSetParentSpan, "set parent span")
	tb.Action(ActCheckExpected, "check that an expected value returned")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"empty name and tags",
			[]string{CndSetEmptyNameTags},
			[]string{ActCheckExpected},
			&condition{},
			&action{
				name: "",
				tags: map[string]any{},
			},
		),
		gen(
			"set single tag",
			[]string{CndSetSingleTag},
			[]string{ActCheckExpected},
			&condition{
				name: "test",
				tags: map[string]string{
					"testKey": "testValue",
				},
			},
			&action{
				name: "test",
				tags: map[string]any{
					"testKey": "testValue",
				},
			},
		),
		gen(
			"set multiple tags",
			[]string{CndSetMultipleTags},
			[]string{ActCheckExpected},
			&condition{
				name: "test",
				tags: map[string]string{
					"testFirstKey":  "testFirstValue",
					"testSecondKey": "testSecondValue",
				},
			},
			&action{
				name: "test",
				tags: map[string]any{
					"testFirstKey":  "testFirstValue",
					"testSecondKey": "testSecondValue",
				},
			},
		),
		gen(
			"set parent span",
			[]string{CndSetParentSpan},
			[]string{ActCheckExpected},
			&condition{
				name:       "test",
				parentSpan: true,
			},
			&action{
				name: "test",
				tags: map[string]any{},
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			mockTracer := mocktracer.New()
			jt := &jaegerTracer{
				tracer: mockTracer,
			}

			ctx := context.Background()

			if tt.C().parentSpan {
				parentSpan := jt.tracer.StartSpan("parent")
				defer parentSpan.Finish()

				ctx = opentracing.ContextWithSpan(ctx, parentSpan)
			}

			_, finish := jt.Trace(ctx, tt.C().name, tt.C().tags)
			finish()

			span := mockTracer.FinishedSpans()[0]

			tags := span.Tags()
			testutil.Diff(t, tt.A().tags, tags)

			name := span.OperationName
			testutil.Diff(t, tt.A().name, name)
		})
	}
}

type mockCloser struct {
	isClosed bool
}

func (m *mockCloser) Close() error {
	m.isClosed = true
	return nil
}

func TestFinalize(t *testing.T) {
	type condition struct {
		mockCloser mockCloser
	}

	type action struct {
		isClosed bool
	}

	CndInputMockCloser := "input mockCloser"
	ActCheckNoError := "check no error was returned"

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	tb.Condition(CndInputMockCloser, "input mockCloser")
	tb.Action(ActCheckNoError, "check that no error was returned")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"input mockCloser",
			[]string{CndInputMockCloser},
			[]string{ActCheckNoError},
			&condition{
				mockCloser: mockCloser{},
			},
			&action{
				isClosed: true,
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			mockCloser := tt.C().mockCloser
			jt := &jaegerTracer{
				closer: &mockCloser,
			}

			jt.Finalize()
			testutil.Diff(t, tt.A().isClosed, mockCloser.isClosed)
		})
	}
}
