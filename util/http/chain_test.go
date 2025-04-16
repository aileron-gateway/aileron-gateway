// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package http_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	v1 "github.com/aileron-gateway/aileron-gateway/apis/core/v1"
	k "github.com/aileron-gateway/aileron-gateway/apis/kernel"
	"github.com/aileron-gateway/aileron-gateway/core"
	"github.com/aileron-gateway/aileron-gateway/kernel/api"
	"github.com/aileron-gateway/aileron-gateway/kernel/er"
	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
	utilhttp "github.com/aileron-gateway/aileron-gateway/util/http"
	"github.com/google/go-cmp/cmp/cmpopts"
)

type appendHeaderHandler string

func (h appendHeaderHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	v := w.Header().Get("Test-Key") + string(h)
	w.Header().Set("Test-Key", v)
}

func (h appendHeaderHandler) Patterns() []string {
	return []string{"/foo", "/bar/"} // One with tailing slash, one without.
}
func (h appendHeaderHandler) Methods() []string {
	return []string{http.MethodGet, http.MethodPost}
}

type appendHeaderMiddleware string

func (m appendHeaderMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		v := w.Header().Get("Test-Key") + string(m)
		w.Header().Set("Test-Key", v)
		next.ServeHTTP(w, r)
	})
}

type appendHeaderRoundTripper string

func (t appendHeaderRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	v := r.Header.Get("Test-Key") + string(t)
	r.Header.Set("Test-Key", v)
	return nil, nil
}

type appendHeaderTripperware string

func (t appendHeaderTripperware) Tripperware(next http.RoundTripper) http.RoundTripper {
	return core.RoundTripperFunc(func(r *http.Request) (*http.Response, error) {
		v := r.Header.Get("Test-Key") + string(t)
		r.Header.Set("Test-Key", v)
		return next.RoundTrip(r)
	})
}

func TestHandler(t *testing.T) {
	type condition struct {
		a    api.API[*api.Request, *api.Response]
		spec *v1.HTTPHandlerSpec
	}

	type action struct {
		methods []string
		paths   []string
		result  string
		err     error
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	cndZero := tb.Condition("zero", "input no middleware")
	cndOne := tb.Condition("one", "input one middleware")
	cndMultiple := tb.Condition("multiple", "input multiple middleware")
	cndMidNil := tb.Condition("middleware not exist", "input nil reference for middleware")
	cndHandNil := tb.Condition("handler not exist", "input nil reference for handler")
	cndMidNotExist := tb.Condition("middleware not exist", "referred middleware does not exist")
	cndHandNotExist := tb.Condition("handler not exist", "referred handler does not exist")
	actCheckOrder := tb.Action("check order", "check the middleware order")
	actCheckError := tb.Action("check error", "check that there is an error")
	actCheckNoError := tb.Action("check no error", "check that there is no error")
	table := tb.Build()

	testAPI := api.NewContainerAPI()
	postTestResource(testAPI, "mid0", appendHeaderMiddleware("0"))
	postTestResource(testAPI, "mid1", appendHeaderMiddleware("1"))
	postTestResource(testAPI, "mid2", appendHeaderMiddleware("2"))
	postTestResource(testAPI, "hand", appendHeaderHandler("H"))

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"no middleware",
			[]string{cndZero},
			[]string{actCheckOrder, actCheckNoError},
			&condition{
				a: testAPI,
				spec: &v1.HTTPHandlerSpec{
					Middleware: []*k.Reference{},
					Handler:    testResourceRef("hand"),
				},
			},
			&action{
				result:  "H",
				paths:   []string{"/foo", "/bar/"},
				methods: []string{http.MethodGet, http.MethodPost},
			},
		),
		gen(
			"no middleware, join path",
			[]string{cndZero},
			[]string{actCheckOrder, actCheckNoError},
			&condition{
				a: testAPI,
				spec: &v1.HTTPHandlerSpec{
					Middleware: []*k.Reference{},
					Handler:    testResourceRef("hand"),
					Pattern:    "/test",
				},
			},
			&action{
				result:  "H",
				paths:   []string{"/test/foo", "/test/bar/"},
				methods: []string{http.MethodGet, http.MethodPost},
			},
		),
		gen(
			"one middleware",
			[]string{cndOne},
			[]string{actCheckOrder, actCheckNoError},
			&condition{
				a: testAPI,
				spec: &v1.HTTPHandlerSpec{
					Middleware: []*k.Reference{
						testResourceRef("mid0"),
					},
					Handler: testResourceRef("hand"),
				},
			},
			&action{
				result:  "0H",
				paths:   []string{"/foo", "/bar/"},
				methods: []string{http.MethodGet, http.MethodPost},
			},
		),
		gen(
			"multiple middleware",
			[]string{cndMultiple},
			[]string{actCheckOrder, actCheckNoError},
			&condition{
				a: testAPI,
				spec: &v1.HTTPHandlerSpec{
					Middleware: []*k.Reference{
						testResourceRef("mid0"),
						testResourceRef("mid1"),
						testResourceRef("mid2"),
					},
					Handler: testResourceRef("hand"),
				},
			},
			&action{
				result:  "012H",
				paths:   []string{"/foo", "/bar/"},
				methods: []string{http.MethodGet, http.MethodPost},
			},
		),
		gen(
			"middleware not exist",
			[]string{cndMidNotExist},
			[]string{actCheckError},
			&condition{
				a: testAPI,
				spec: &v1.HTTPHandlerSpec{
					Middleware: []*k.Reference{
						testResourceRef("not-exist"),
					},
					Handler: testResourceRef("hand"),
				},
			},
			&action{
				result: "",
				err: &er.Error{
					Package:     api.ErrPkg,
					Type:        api.ErrTypeUtil,
					Description: api.ErrDscAssert,
				},
			},
		),
		gen(
			"middleware nil reference",
			[]string{cndMidNil},
			[]string{actCheckError},
			&condition{
				a: testAPI,
				spec: &v1.HTTPHandlerSpec{
					Middleware: []*k.Reference{
						nil,
					},
					Handler: testResourceRef("hand"),
				},
			},
			&action{
				result: "",
				err: &er.Error{
					Package:     api.ErrPkg,
					Type:        api.ErrTypeUtil,
					Description: api.ErrDscNil,
				},
			},
		),
		gen(
			"handler not exist",
			[]string{cndHandNotExist},
			[]string{actCheckError},
			&condition{
				a: testAPI,
				spec: &v1.HTTPHandlerSpec{
					Middleware: []*k.Reference{},
					Handler:    testResourceRef("not-exist"),
				},
			},
			&action{
				result: "",
				err: &er.Error{
					Package:     utilhttp.ErrPkg,
					Type:        utilhttp.ErrTypeChain,
					Description: utilhttp.ErrDscAssert,
				},
			},
		),
		gen(
			"handler nil reference",
			[]string{cndHandNil},
			[]string{actCheckError},
			&condition{
				a: testAPI,
				spec: &v1.HTTPHandlerSpec{
					Middleware: []*k.Reference{},
					Handler:    nil,
				},
			},
			&action{
				result: "",
				err: &er.Error{
					Package:     api.ErrPkg,
					Type:        api.ErrTypeUtil,
					Description: api.ErrDscNil,
				},
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			methods, paths, h, err := utilhttp.Handler(tt.C().a, tt.C().spec)

			testutil.Diff(t, tt.A().paths, paths, cmpopts.SortSlices(func(s1, s2 string) bool { return s1 > s2 }))
			testutil.Diff(t, tt.A().methods, methods, cmpopts.SortSlices(func(s1, s2 string) bool { return s1 > s2 }))
			testutil.Diff(t, tt.A().err, err, cmpopts.EquateErrors())
			if err != nil {
				return
			}

			w := httptest.NewRecorder()
			r, _ := http.NewRequest(http.MethodGet, "http://test.com", nil)
			h.ServeHTTP(w, r)

			resp := w.Result()
			testutil.Diff(t, tt.A().result, resp.Header.Get("Test-Key"))
		})
	}
}

func TestMiddlewareChain(t *testing.T) {
	type condition struct {
		ms []core.Middleware
		h  http.Handler
	}

	type action struct {
		result string
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	cndZero := tb.Condition("zero", "input no middleware")
	cndOne := tb.Condition("one", "input one middleware")
	cndMultiple := tb.Condition("multiple", "input multiple middleware")
	cndNil := tb.Condition("multiple", "middleware contains nil")
	actCheckOrder := tb.Action("check order", "check the middleware order")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"no middleware",
			[]string{cndZero},
			[]string{actCheckOrder},
			&condition{
				ms: []core.Middleware{},
				h:  appendHeaderHandler("H"),
			},
			&action{
				result: "H",
			},
		),
		gen(
			"one middleware",
			[]string{cndOne},
			[]string{actCheckOrder},
			&condition{
				ms: []core.Middleware{
					appendHeaderMiddleware("0"),
				},
				h: appendHeaderHandler("H"),
			},
			&action{
				result: "0H",
			},
		),
		gen(
			"multiple middleware",
			[]string{cndMultiple},
			[]string{actCheckOrder},
			&condition{
				ms: []core.Middleware{
					appendHeaderMiddleware("0"),
					appendHeaderMiddleware("1"),
					appendHeaderMiddleware("2"),
				},
				h: appendHeaderHandler("H"),
			},
			&action{
				result: "012H",
			},
		),
		gen(
			"contains nil",
			[]string{cndMultiple, cndNil},
			[]string{actCheckOrder},
			&condition{
				ms: []core.Middleware{
					appendHeaderMiddleware("0"),
					appendHeaderMiddleware("1"),
					nil,
					appendHeaderMiddleware("2"),
				},
				h: appendHeaderHandler("H"),
			},
			&action{
				result: "012H",
			},
		),
		gen(
			"contains nil only",
			[]string{cndNil},
			[]string{actCheckOrder},
			&condition{
				ms: []core.Middleware{
					nil,
				},
				h: appendHeaderHandler("H"),
			},
			&action{
				result: "H",
			},
		),
		gen(
			"contains multiple nil",
			[]string{cndNil},
			[]string{actCheckOrder},
			&condition{
				ms: []core.Middleware{
					nil,
					nil,
					nil,
				},
				h: appendHeaderHandler("H"),
			},
			&action{
				result: "H",
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			w := httptest.NewRecorder()
			r, _ := http.NewRequest(http.MethodGet, "http://test.com", nil)

			h := utilhttp.MiddlewareChain(tt.C().ms, tt.C().h)
			h.ServeHTTP(w, r)

			resp := w.Result()
			testutil.Diff(t, tt.A().result, resp.Header.Get("Test-Key"))
		})
	}
}

func TestTripperwareChain(t *testing.T) {
	type condition struct {
		ts []core.Tripperware
		t  http.RoundTripper
	}

	type action struct {
		result string
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	cndZero := tb.Condition("zero", "input no tripperware")
	cndOne := tb.Condition("one", "input one tripperware")
	cndMultiple := tb.Condition("multiple", "input multiple tripperware")
	cndNil := tb.Condition("multiple", "tripperware contains nil")
	actCheckOrder := tb.Action("check order", "check the tripperware order")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"no tripperware",
			[]string{cndZero},
			[]string{actCheckOrder},
			&condition{
				ts: []core.Tripperware{},
				t:  appendHeaderRoundTripper("T"),
			},
			&action{
				result: "T",
			},
		),
		gen(
			"one tripperware",
			[]string{cndOne},
			[]string{actCheckOrder},
			&condition{
				ts: []core.Tripperware{
					appendHeaderTripperware("0"),
				},
				t: appendHeaderRoundTripper("T"),
			},
			&action{
				result: "0T",
			},
		),
		gen(
			"multiple tripperware",
			[]string{cndMultiple},
			[]string{actCheckOrder},
			&condition{
				ts: []core.Tripperware{
					appendHeaderTripperware("0"),
					appendHeaderTripperware("1"),
					appendHeaderTripperware("2"),
				},
				t: appendHeaderRoundTripper("T"),
			},
			&action{
				result: "012T",
			},
		),
		gen(
			"contains nil",
			[]string{cndMultiple, cndNil},
			[]string{actCheckOrder},
			&condition{
				ts: []core.Tripperware{
					appendHeaderTripperware("0"),
					appendHeaderTripperware("1"),
					nil,
					appendHeaderTripperware("2"),
				},
				t: appendHeaderRoundTripper("T"),
			},
			&action{
				result: "012T",
			},
		),
		gen(
			"contains nil only",
			[]string{cndNil},
			[]string{actCheckOrder},
			&condition{
				ts: []core.Tripperware{
					nil,
				},
				t: appendHeaderRoundTripper("T"),
			},
			&action{
				result: "T",
			},
		),
		gen(
			"contains multiple nil",
			[]string{cndNil},
			[]string{actCheckOrder},
			&condition{
				ts: []core.Tripperware{
					nil,
					nil,
					nil,
				},
				t: appendHeaderRoundTripper("T"),
			},
			&action{
				result: "T",
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			r, _ := http.NewRequest(http.MethodGet, "http://test.com", nil)

			h := utilhttp.TripperwareChain(tt.C().ts, tt.C().t)
			h.RoundTrip(r)

			testutil.Diff(t, tt.A().result, r.Header.Get("Test-Key"))
		})
	}
}

func postTestResource(server api.API[*api.Request, *api.Response], name string, res any) {
	ref := testResourceRef(name)
	req := &api.Request{
		Method:  api.MethodPost,
		Key:     ref.APIVersion + "/" + ref.Kind + "/" + ref.Namespace + "/" + ref.Name,
		Content: res,
	}
	if _, err := server.Serve(context.Background(), req); err != nil {
		panic(err)
	}
}

func testResourceRef(name string) *k.Reference {
	return &k.Reference{
		APIVersion: "core/v1",
		Kind:       "Container",
		Namespace:  "test",
		Name:       name,
	}
}
