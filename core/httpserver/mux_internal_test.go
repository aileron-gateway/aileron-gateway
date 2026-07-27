// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package httpserver

import (
	"net/http"
	"net/http/httptest"
	"regexp"
	"slices"
	"testing"

	v1 "github.com/aileron-gateway/aileron-gateway/apis/core/v1"
	k "github.com/aileron-gateway/aileron-gateway/apis/kernel"
	"github.com/aileron-gateway/aileron-gateway/core"
	"github.com/aileron-gateway/aileron-gateway/internal/testutil"
	"github.com/aileron-gateway/aileron-gateway/kernel/api"
	utilhttp "github.com/aileron-gateway/aileron-gateway/util/http"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestNotFoundHandler(t *testing.T) {
	type condition struct {
		url string
	}

	type action struct {
		body string
	}
	// Create API for test.
	testAPI := api.NewContainerAPI()
	postTestResource(testAPI, "handler", &testHandler{
		headers: map[string]string{"test": "handler"},
		body:    "test",
	})
	postTestResource(testAPI, "middleware", &testMiddleware{
		headers: map[string]string{"test": "middleware"},
	})

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"not found",
			&condition{
				url: "http://test.com/test",
			},
			&action{
				body: `{"status":404,"statusText":"Not Found"}`,
			},
		),
		gen(
			"found",
			&condition{
				url: "http://test.com/foo",
			},
			&action{
				body: "404 page not found\n",
			},
		),
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			handler := &http.ServeMux{}
			handler.Handle("/foo", http.DefaultServeMux) // Register dummy handler.
			handler.Handle("/bar", http.DefaultServeMux) // Register dummy handler.

			eh := utilhttp.GlobalErrorHandler(utilhttp.DefaultErrorHandlerName)
			handler.Handle("/", notFoundHandler(eh)) // Register NotFound handler.

			w := httptest.NewRecorder()
			r, _ := http.NewRequest(http.MethodGet, tt.C.url, nil)

			handler.ServeHTTP(w, r)
			resp := w.Result()
			defer resp.Body.Close()

			testutil.Diff(t, tt.A.body, w.Body.String())
		})
	}
}

func TestRegisterHandlers(t *testing.T) {
	type condition struct {
		specs []*v1.VirtualHostSpec
	}

	type action struct {
		handlers   map[string]http.Handler
		err        any // error or errorutil.Kind
		errPattern *regexp.Regexp
	}

	// Create API for test.
	testAPI := api.NewContainerAPI()
	h1 := &testHandler{
		id:      "handler1",
		headers: map[string]string{"handler": "h1"},
	}
	h2 := &testHandler{
		id:       "handler2",
		patterns: []string{"/test1", "/test2/"}, // One with tailing slash, one without.
		headers:  map[string]string{"handler": "h2"},
	}
	h3 := &testHandler{
		id:      "handler3",
		methods: []string{http.MethodGet, http.MethodHead},
		headers: map[string]string{"handler": "h3"},
	}
	h4 := &testHandler{
		id:       "handler3",
		patterns: []string{"/test1", "/test2/"}, // One with tailing slash, one without.
		methods:  []string{http.MethodGet, http.MethodHead},
		headers:  map[string]string{"handler": "h4"},
	}
	m := &testMiddleware{
		id:      "middleware",
		headers: map[string]string{"middleware": "m"},
	}
	postTestResource(testAPI, "handler1", h1)
	postTestResource(testAPI, "handler2", h2)
	postTestResource(testAPI, "handler3", h3)
	postTestResource(testAPI, "handler4", h4)
	postTestResource(testAPI, "middleware", m)

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"no handler",
			&condition{
				specs: []*v1.VirtualHostSpec{
					{
						Handlers: nil, // no handler
					},
				},
			},
			&action{
				handlers: map[string]http.Handler{},
				err:      nil,
			},
		),
		gen(
			"wo/path wo/methods",
			&condition{
				specs: []*v1.VirtualHostSpec{
					{
						Handlers: []*v1.HTTPHandlerSpec{
							{Handler: testResourceRef("handler1")},
						},
					},
				},
			},
			&action{
				handlers: map[string]http.Handler{
					"GET /": h1, "HEAD /": h1, "POST /": h1,
					"PUT /": h1, "PATCH /": h1, "DELETE /": h1,
					"CONNECT /": h1, "OPTIONS /": h1, "TRACE /": h1,
				},
				err: nil,
			},
		),
		gen(
			"w/path wo/methods",
			&condition{
				specs: []*v1.VirtualHostSpec{
					{
						Handlers: []*v1.HTTPHandlerSpec{
							{Handler: testResourceRef("handler2")},
						},
					},
				},
			},
			&action{
				handlers: map[string]http.Handler{
					"GET /test1": h2, "HEAD /test1": h2, "POST /test1": h2,
					"PUT /test1": h2, "PATCH /test1": h2, "DELETE /test1": h2,
					"CONNECT /test1": h2, "OPTIONS /test1": h2, "TRACE /test1": h2,
					"GET /test2/": h2, "HEAD /test2/": h2, "POST /test2/": h2,
					"PUT /test2/": h2, "PATCH /test2/": h2, "DELETE /test2/": h2,
					"CONNECT /test2/": h2, "OPTIONS /test2/": h2, "TRACE /test2/": h2,
				},
				err: nil,
			},
		),
		gen(
			"wo/path w/methods",
			&condition{
				specs: []*v1.VirtualHostSpec{
					{
						Handlers: []*v1.HTTPHandlerSpec{
							{Handler: testResourceRef("handler3")},
						},
					},
				},
			},
			&action{
				handlers: map[string]http.Handler{"GET /": h3, "HEAD /": h3},
				err:      nil,
			},
		),
		gen(
			"w/path w/methods",
			&condition{
				specs: []*v1.VirtualHostSpec{
					{
						Handlers: []*v1.HTTPHandlerSpec{
							{Handler: testResourceRef("handler4")},
						},
					},
				},
			},
			&action{
				handlers: map[string]http.Handler{
					"GET /test1": h4, "HEAD /test1": h4,
					"GET /test2/": h4, "HEAD /test2/": h4,
				},
				err: nil,
			},
		),
		gen(
			"wo/path wo/methods and host path",
			&condition{
				specs: []*v1.VirtualHostSpec{
					{
						Hosts:   []string{"example.com"},
						Pattern: "/pattern",
						Handlers: []*v1.HTTPHandlerSpec{
							{Handler: testResourceRef("handler1")},
						},
					},
				},
			},
			&action{
				handlers: map[string]http.Handler{
					"GET example.com/pattern": h1, "HEAD example.com/pattern": h1, "POST example.com/pattern": h1,
					"PUT example.com/pattern": h1, "PATCH example.com/pattern": h1, "DELETE example.com/pattern": h1,
					"CONNECT example.com/pattern": h1, "OPTIONS example.com/pattern": h1, "TRACE example.com/pattern": h1,
				},
				err: nil,
			},
		),
		gen(
			"w/path wo/methods and host path",
			&condition{
				specs: []*v1.VirtualHostSpec{
					{
						Hosts:   []string{"example.com"},
						Pattern: "/pattern",
						Handlers: []*v1.HTTPHandlerSpec{
							{Handler: testResourceRef("handler2")},
						},
					},
				},
			},
			&action{
				handlers: map[string]http.Handler{
					"GET example.com/pattern/test1": h2, "HEAD example.com/pattern/test1": h2, "POST example.com/pattern/test1": h2,
					"PUT example.com/pattern/test1": h2, "PATCH example.com/pattern/test1": h2, "DELETE example.com/pattern/test1": h2,
					"CONNECT example.com/pattern/test1": h2, "OPTIONS example.com/pattern/test1": h2, "TRACE example.com/pattern/test1": h2,
					"GET example.com/pattern/test2/": h2, "HEAD example.com/pattern/test2/": h2, "POST example.com/pattern/test2/": h2,
					"PUT example.com/pattern/test2/": h2, "PATCH example.com/pattern/test2/": h2, "DELETE example.com/pattern/test2/": h2,
					"CONNECT example.com/pattern/test2/": h2, "OPTIONS example.com/pattern/test2/": h2, "TRACE example.com/pattern/test2/": h2,
				},
				err: nil,
			},
		),
		gen(
			"wo/path w/methods and host path",
			&condition{
				specs: []*v1.VirtualHostSpec{
					{
						Hosts:   []string{"example.com"},
						Pattern: "/pattern",
						Handlers: []*v1.HTTPHandlerSpec{
							{Handler: testResourceRef("handler3")},
						},
					},
				},
			},
			&action{
				handlers: map[string]http.Handler{
					"GET example.com/pattern": h3, "HEAD example.com/pattern": h3,
				},
				err: nil,
			},
		),
		gen(
			"w/path w/methods and host path",
			&condition{
				specs: []*v1.VirtualHostSpec{
					{
						Hosts:   []string{"example.com"},
						Pattern: "/pattern",
						Handlers: []*v1.HTTPHandlerSpec{
							{Handler: testResourceRef("handler4")},
						},
					},
				},
			},
			&action{
				handlers: map[string]http.Handler{
					"GET example.com/pattern/test1": h4, "HEAD example.com/pattern/test1": h4,
					"GET example.com/pattern/test2/": h4, "HEAD example.com/pattern/test2/": h4,
				},
				err: nil,
			},
		),
		gen(
			"wo/path wo/methods and host methods",
			&condition{
				specs: []*v1.VirtualHostSpec{
					{
						Methods: []v1.HTTPMethod{v1.HTTPMethod_GET, v1.HTTPMethod_DELETE},
						Handlers: []*v1.HTTPHandlerSpec{
							{Handler: testResourceRef("handler1")},
						},
					},
				},
			},
			&action{
				handlers: map[string]http.Handler{
					"GET /": h1, "DELETE /": h1,
				},
				err: nil,
			},
		),
		gen(
			"w/path wo/methods and host methods",
			&condition{
				specs: []*v1.VirtualHostSpec{
					{
						Methods: []v1.HTTPMethod{v1.HTTPMethod_GET, v1.HTTPMethod_DELETE},
						Handlers: []*v1.HTTPHandlerSpec{
							{Handler: testResourceRef("handler2")},
						},
					},
				},
			},
			&action{
				handlers: map[string]http.Handler{
					"GET /test1": h2, "DELETE /test1": h2,
					"GET /test2/": h2, "DELETE /test2/": h2,
				},
				err: nil,
			},
		),
		gen(
			"wo/path w/methods and host methods",
			&condition{
				specs: []*v1.VirtualHostSpec{
					{
						Methods: []v1.HTTPMethod{v1.HTTPMethod_GET, v1.HTTPMethod_DELETE},
						Handlers: []*v1.HTTPHandlerSpec{
							{Handler: testResourceRef("handler3")},
						},
					},
				},
			},
			&action{
				handlers: map[string]http.Handler{
					"GET /": h3,
				},
				err: nil,
			},
		),
		gen(
			"w/path w/methods and host methods",
			&condition{
				specs: []*v1.VirtualHostSpec{
					{
						Methods: []v1.HTTPMethod{v1.HTTPMethod_GET, v1.HTTPMethod_DELETE},
						Handlers: []*v1.HTTPHandlerSpec{
							{Handler: testResourceRef("handler4")},
						},
					},
				},
			},
			&action{
				handlers: map[string]http.Handler{
					"GET /test1":  h4,
					"GET /test2/": h4,
				},
				err: nil,
			},
		),
		gen(
			"wo/path wo/methods and host path and host methods",
			&condition{
				specs: []*v1.VirtualHostSpec{
					{
						Hosts:   []string{"example.com"},
						Pattern: "/pattern",
						Methods: []v1.HTTPMethod{v1.HTTPMethod_GET, v1.HTTPMethod_DELETE},
						Handlers: []*v1.HTTPHandlerSpec{
							{Handler: testResourceRef("handler1")},
						},
					},
				},
			},
			&action{
				handlers: map[string]http.Handler{
					"GET example.com/pattern": h1, "DELETE example.com/pattern": h1,
				},
				err: nil,
			},
		),
		gen(
			"w/path wo/methods and host path and host methods",
			&condition{
				specs: []*v1.VirtualHostSpec{
					{
						Hosts:   []string{"example.com"},
						Pattern: "/pattern",
						Methods: []v1.HTTPMethod{v1.HTTPMethod_GET, v1.HTTPMethod_DELETE},
						Handlers: []*v1.HTTPHandlerSpec{
							{Handler: testResourceRef("handler2")},
						},
					},
				},
			},
			&action{
				handlers: map[string]http.Handler{
					"GET example.com/pattern/test1": h2, "DELETE example.com/pattern/test1": h2,
					"GET example.com/pattern/test2/": h2, "DELETE example.com/pattern/test2/": h2,
				},
				err: nil,
			},
		),
		gen(
			"wo/path w/methods and host path and host methods",
			&condition{
				specs: []*v1.VirtualHostSpec{
					{
						Hosts:   []string{"example.com"},
						Pattern: "/pattern",
						Methods: []v1.HTTPMethod{v1.HTTPMethod_GET, v1.HTTPMethod_DELETE},
						Handlers: []*v1.HTTPHandlerSpec{
							{Handler: testResourceRef("handler3")},
						},
					},
				},
			},
			&action{
				handlers: map[string]http.Handler{
					"GET example.com/pattern": h3,
				},
				err: nil,
			},
		),
		gen(
			"w/path w/methods and host path and host methods",
			&condition{
				specs: []*v1.VirtualHostSpec{
					{
						Hosts:   []string{"example.com"},
						Pattern: "/pattern",
						Methods: []v1.HTTPMethod{v1.HTTPMethod_GET, v1.HTTPMethod_DELETE},
						Handlers: []*v1.HTTPHandlerSpec{
							{Handler: testResourceRef("handler4")},
						},
					},
				},
			},
			&action{
				handlers: map[string]http.Handler{
					"GET example.com/pattern/test1":  h4,
					"GET example.com/pattern/test2/": h4,
				},
				err: nil,
			},
		),
		gen(
			"use middleware",
			&condition{
				specs: []*v1.VirtualHostSpec{
					{
						Hosts:   []string{"example.com"},
						Pattern: "/pattern",
						Methods: []v1.HTTPMethod{v1.HTTPMethod_GET, v1.HTTPMethod_DELETE},
						Handlers: []*v1.HTTPHandlerSpec{
							{Handler: testResourceRef("handler4")},
						},
						Middleware: []*k.Reference{
							testResourceRef("middleware"),
						},
					},
				},
			},
			&action{
				handlers: map[string]http.Handler{
					"GET example.com/pattern/test1":  m.Middleware(h4),
					"GET example.com/pattern/test2/": m.Middleware(h4),
				},
				err: nil,
			},
		),
		gen(
			"invalid path ptatern",
			&condition{
				specs: []*v1.VirtualHostSpec{
					{
						Pattern: "{bar..}",
						Hosts:   []string{""},
						Handlers: []*v1.HTTPHandlerSpec{
							{Handler: testResourceRef("handler1")},
						},
					},
				},
			},
			&action{
				err:        core.ErrCoreGenCreateComponent,
				errPattern: regexp.MustCompile(core.ErrPrefix + `failed to create component. failed to register handler to HTTP router`),
			},
		),
		gen(
			"reference to a invalid middleware",
			&condition{
				specs: []*v1.VirtualHostSpec{
					{
						Middleware: []*k.Reference{
							testResourceRef("not exist middleware"),
						},
					},
				},
			},
			&action{
				err:        core.ErrCoreGenCreateComponent,
				errPattern: regexp.MustCompile(core.ErrPrefix + `failed to create component. failed to get middleware`),
			},
		),
		gen(
			"reference to a invalid handler",
			&condition{
				specs: []*v1.VirtualHostSpec{
					{
						Handlers: []*v1.HTTPHandlerSpec{
							{
								Handler: testResourceRef("not exist handler"),
							},
						},
					},
				},
			},
			&action{
				err:        core.ErrCoreGenCreateComponent,
				errPattern: regexp.MustCompile(core.ErrPrefix + `failed to create component. failed to create handler`),
			},
		),
	}

	for _, tt := range testCases {
		t.Run(tt.Name, func(t *testing.T) {
			mux := &testMux{
				Mux: &http.ServeMux{},
				hs:  make(map[string]http.Handler),
			}
			err := registerHandlers(testAPI, mux, tt.C.specs)
			testutil.DiffError(t, tt.A.err, tt.A.errPattern, err)
			if err != nil {
				return
			}

			opts := []cmp.Option{
				cmp.AllowUnexported(testHandler{}, testMiddleware{}),
				cmpopts.EquateEmpty(),
				cmpopts.SortMaps(func(x, y string) bool { return x > y }),
				cmp.Comparer(func(x, y http.HandlerFunc) bool {
					r := httptest.NewRequest(http.MethodGet, "/test", nil)
					wx := httptest.NewRecorder()
					wy := httptest.NewRecorder()
					x.ServeHTTP(wx, r)
					y.ServeHTTP(wy, r)
					xh := wx.Result().Header
					yh := wy.Result().Header
					return slices.Equal(xh["Handler"], yh["Handler"]) && slices.Equal(xh["Middleware"], yh["Middleware"])
				}),
			}
			testutil.Diff(t, tt.A.handlers, mux.hs, opts...)
			testutil.Diff(t, len(tt.A.handlers), len(mux.hs))
		})
	}
}

func TestIntersectionString(t *testing.T) {
	type condition struct {
		set1 []string
		set2 []string
	}

	type action struct {
		set []string
	}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"mutually exclusive",
			&condition{
				set1: []string{"test1", "test2"},
				set2: []string{"test3", "test4"},
			},
			&action{
				set: []string{},
			},
		),
		gen(
			"extract duplicates",
			&condition{
				set1: []string{"test1", "test2"},
				set2: []string{"test2", "test3"},
			},
			&action{
				set: []string{"test2"},
			},
		),
		gen(
			"set1 nil",
			&condition{
				set1: nil,
				set2: []string{"test1", "test2"},
			},
			&action{
				set: nil,
			},
		),
		gen(
			"set2 nil",
			&condition{
				set1: []string{"test1", "test2"},
				set2: nil,
			},
			&action{
				set: nil,
			},
		),
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			set := intersectionString(tt.C.set1, tt.C.set2)
			testutil.Diff(t, tt.A.set, set)
		})
	}
}

func TestGeneratePatterns(t *testing.T) {
	type condition struct {
		methods []string
		hosts   []string
		paths   []string
	}

	type action struct {
		patterns []string
	}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"0 method, 0 hosts, 0 paths",
			&condition{
				methods: []string{},
				hosts:   []string{},
				paths:   []string{},
			},
			&action{
				patterns: []string{"GET /", "HEAD /", "POST /", "PUT /", "PATCH /",
					"DELETE /", "CONNECT /", "OPTIONS /", "TRACE /"},
			},
		),
		gen(
			"0 method, 0 hosts, 1 paths",
			&condition{
				methods: []string{},
				hosts:   []string{},
				paths:   []string{"/foo"},
			},
			&action{
				patterns: []string{"GET /foo", "HEAD /foo", "POST /foo", "PUT /foo", "PATCH /foo",
					"DELETE /foo", "CONNECT /foo", "OPTIONS /foo", "TRACE /foo"},
			},
		),
		gen(
			"0 method, 0 hosts, 2 paths",
			&condition{
				methods: []string{},
				hosts:   []string{},
				paths:   []string{"/foo", "bar"},
			},
			&action{
				patterns: []string{"GET /foo", "HEAD /foo", "POST /foo", "PUT /foo", "PATCH /foo",
					"DELETE /foo", "CONNECT /foo", "OPTIONS /foo", "TRACE /foo",
					"GET /bar", "HEAD /bar", "POST /bar", "PUT /bar", "PATCH /bar",
					"DELETE /bar", "CONNECT /bar", "OPTIONS /bar", "TRACE /bar"},
			},
		),
		gen(
			"0 method, 1 hosts, 0 paths",
			&condition{
				methods: []string{},
				hosts:   []string{"foo.com"},
				paths:   []string{},
			},
			&action{
				patterns: []string{"GET foo.com/", "HEAD foo.com/", "POST foo.com/", "PUT foo.com/", "PATCH foo.com/",
					"DELETE foo.com/", "CONNECT foo.com/", "OPTIONS foo.com/", "TRACE foo.com/"},
			},
		),
		gen(
			"0 method, 1 hosts, 1 paths",
			&condition{
				methods: []string{},
				hosts:   []string{"foo.com"},
				paths:   []string{"/foo"},
			},
			&action{
				patterns: []string{"GET foo.com/foo", "HEAD foo.com/foo", "POST foo.com/foo", "PUT foo.com/foo", "PATCH foo.com/foo",
					"DELETE foo.com/foo", "CONNECT foo.com/foo", "OPTIONS foo.com/foo", "TRACE foo.com/foo"},
			},
		),
		gen(
			"0 method, 1 hosts, 2 paths",
			&condition{
				methods: []string{},
				hosts:   []string{"foo.com"},
				paths:   []string{"/foo", "bar"},
			},
			&action{
				patterns: []string{"GET foo.com/foo", "HEAD foo.com/foo", "POST foo.com/foo", "PUT foo.com/foo", "PATCH foo.com/foo",
					"DELETE foo.com/foo", "CONNECT foo.com/foo", "OPTIONS foo.com/foo", "TRACE foo.com/foo",
					"GET foo.com/bar", "HEAD foo.com/bar", "POST foo.com/bar", "PUT foo.com/bar", "PATCH foo.com/bar",
					"DELETE foo.com/bar", "CONNECT foo.com/bar", "OPTIONS foo.com/bar", "TRACE foo.com/bar"},
			},
		),
		gen(
			"0 method, 2 hosts, 0 paths",
			&condition{
				methods: []string{},
				hosts:   []string{"foo.com", "bar.com"},
				paths:   []string{},
			},
			&action{
				patterns: []string{"GET foo.com/", "HEAD foo.com/", "POST foo.com/", "PUT foo.com/", "PATCH foo.com/",
					"DELETE foo.com/", "CONNECT foo.com/", "OPTIONS foo.com/", "TRACE foo.com/",
					"GET bar.com/", "HEAD bar.com/", "POST bar.com/", "PUT bar.com/", "PATCH bar.com/",
					"DELETE bar.com/", "CONNECT bar.com/", "OPTIONS bar.com/", "TRACE bar.com/"},
			},
		),
		gen(
			"0 method, 2 hosts, 1 paths",
			&condition{
				methods: []string{},
				hosts:   []string{"foo.com", "bar.com"},
				paths:   []string{"/foo"},
			},
			&action{
				patterns: []string{"GET foo.com/foo", "HEAD foo.com/foo", "POST foo.com/foo", "PUT foo.com/foo", "PATCH foo.com/foo",
					"DELETE foo.com/foo", "CONNECT foo.com/foo", "OPTIONS foo.com/foo", "TRACE foo.com/foo",
					"GET bar.com/foo", "HEAD bar.com/foo", "POST bar.com/foo", "PUT bar.com/foo", "PATCH bar.com/foo",
					"DELETE bar.com/foo", "CONNECT bar.com/foo", "OPTIONS bar.com/foo", "TRACE bar.com/foo"},
			},
		),
		// gen(
		// 	"0 method, 2 hosts, 2 paths",
		// 	&condition{
		// 		methods: []string{},
		// 		hosts:   []string{"foo.com", "bar.com"},
		// 		paths:   []string{"/foo", "bar"},
		// 	},
		// 	&action{
		// 		patterns: []string{"foo.com/foo", "foo.com/bar", "bar.com/foo", "bar.com/bar"},
		// 	},
		// ),
		gen(
			"1 method, 0 hosts, 0 paths",
			&condition{
				methods: []string{http.MethodGet},
				hosts:   []string{},
				paths:   []string{},
			},
			&action{
				patterns: []string{"GET /"},
			},
		),
		gen(
			"1 method, 0 hosts, 1 paths",
			&condition{
				methods: []string{http.MethodGet},
				hosts:   []string{},
				paths:   []string{"/foo"},
			},
			&action{
				patterns: []string{"GET /foo"},
			},
		),
		gen(
			"1 method, 0 hosts, 2 paths",
			&condition{
				methods: []string{http.MethodGet},
				hosts:   []string{},
				paths:   []string{"/foo", "bar"},
			},
			&action{
				patterns: []string{"GET /foo", "GET /bar"},
			},
		),
		gen(
			"1 method, 1 hosts, 0 paths",
			&condition{
				methods: []string{http.MethodGet},
				hosts:   []string{"foo.com"},
				paths:   []string{},
			},
			&action{
				patterns: []string{"GET foo.com/"},
			},
		),
		gen(
			"1 method, 1 hosts, 1 paths",
			&condition{
				methods: []string{http.MethodGet},
				hosts:   []string{"foo.com"},
				paths:   []string{"/foo"},
			},
			&action{
				patterns: []string{"GET foo.com/foo"},
			},
		),
		gen(
			"1 method, 1 hosts, 2 paths",
			&condition{
				methods: []string{http.MethodGet},
				hosts:   []string{"foo.com"},
				paths:   []string{"/foo", "bar"},
			},
			&action{
				patterns: []string{"GET foo.com/foo", "GET foo.com/bar"},
			},
		),
		gen(
			"1 method, 2 hosts, 0 paths",
			&condition{
				methods: []string{http.MethodGet},
				hosts:   []string{"foo.com", "bar.com"},
				paths:   []string{},
			},
			&action{
				patterns: []string{"GET foo.com/", "GET bar.com/"},
			},
		),
		gen(
			"1 method, 2 hosts, 1 paths",
			&condition{
				methods: []string{http.MethodGet},
				hosts:   []string{"foo.com", "bar.com"},
				paths:   []string{"/foo"},
			},
			&action{
				patterns: []string{"GET foo.com/foo", "GET bar.com/foo"},
			},
		),
		gen(
			"1 method, 2 hosts, 2 paths",
			&condition{
				methods: []string{http.MethodGet},
				hosts:   []string{"foo.com", "bar.com"},
				paths:   []string{"/foo", "bar"},
			},
			&action{
				patterns: []string{"GET foo.com/foo", "GET foo.com/bar", "GET bar.com/foo", "GET bar.com/bar"},
			},
		),
		gen(
			"2 method, 0 hosts, 0 paths",
			&condition{
				methods: []string{http.MethodGet, http.MethodPost},
				hosts:   []string{},
				paths:   []string{},
			},
			&action{
				patterns: []string{"GET /", "POST /"},
			},
		),
		gen(
			"2 method, 0 hosts, 1 paths",
			&condition{
				methods: []string{http.MethodGet, http.MethodPost},
				hosts:   []string{},
				paths:   []string{"/foo"},
			},
			&action{
				patterns: []string{"GET /foo", "POST /foo"},
			},
		),
		gen(
			"2 method, 0 hosts, 2 paths",
			&condition{
				methods: []string{http.MethodGet, http.MethodPost},
				hosts:   []string{},
				paths:   []string{"/foo", "bar"},
			},
			&action{
				patterns: []string{"GET /foo", "GET /bar", "POST /foo", "POST /bar"},
			},
		),
		gen(
			"2 method, 1 hosts, 0 paths",
			&condition{
				methods: []string{http.MethodGet, http.MethodPost},
				hosts:   []string{"foo.com"},
				paths:   []string{},
			},
			&action{
				patterns: []string{"GET foo.com/", "POST foo.com/"},
			},
		),
		gen(
			"2 method, 1 hosts, 1 paths",
			&condition{
				methods: []string{http.MethodGet, http.MethodPost},
				hosts:   []string{"foo.com"},
				paths:   []string{"/foo"},
			},
			&action{
				patterns: []string{"GET foo.com/foo", "POST foo.com/foo"},
			},
		),
		gen(
			"2 method, 1 hosts, 2 paths",
			&condition{
				methods: []string{http.MethodGet, http.MethodPost},
				hosts:   []string{"foo.com"},
				paths:   []string{"/foo", "bar"},
			},
			&action{
				patterns: []string{"GET foo.com/foo", "GET foo.com/bar", "POST foo.com/foo", "POST foo.com/bar"},
			},
		),
		gen(
			"2 method, 2 hosts, 0 paths",
			&condition{
				methods: []string{http.MethodGet, http.MethodPost},
				hosts:   []string{"foo.com", "bar.com"},
				paths:   []string{},
			},
			&action{
				patterns: []string{"GET foo.com/", "GET bar.com/", "POST foo.com/", "POST bar.com/"},
			},
		),
		gen(
			"2 method, 2 hosts, 1 paths",
			&condition{
				methods: []string{http.MethodGet, http.MethodPost},
				hosts:   []string{"foo.com", "bar.com"},
				paths:   []string{"/foo"},
			},
			&action{
				patterns: []string{"GET foo.com/foo", "GET bar.com/foo", "POST foo.com/foo", "POST bar.com/foo"},
			},
		),
		gen(
			"2 method, 2 hosts, 2 paths",
			&condition{
				methods: []string{http.MethodGet, http.MethodPost},
				hosts:   []string{"foo.com", "bar.com"},
				paths:   []string{"/foo", "bar"},
			},
			&action{
				patterns: []string{
					"GET foo.com/foo", "GET foo.com/bar", "GET bar.com/foo", "GET bar.com/bar",
					"POST foo.com/foo", "POST foo.com/bar", "POST bar.com/foo", "POST bar.com/bar",
				},
			},
		),
		gen(
			"invalid pattern path",
			&condition{
				methods: []string{},
				hosts:   []string{},
				paths:   []string{"foo.com/foo"},
			},
			&action{
				// Unexpected pattern.
				patterns: []string{"GET /foo.com/foo", "HEAD /foo.com/foo", "POST /foo.com/foo", "PUT /foo.com/foo", "PATCH /foo.com/foo",
					"DELETE /foo.com/foo", "CONNECT /foo.com/foo", "OPTIONS /foo.com/foo", "TRACE /foo.com/foo"},
			},
		),
		gen(
			"invalid pattern method+path",
			&condition{
				methods: []string{http.MethodGet},
				hosts:   []string{},
				paths:   []string{"foo.com/foo"},
			},
			&action{
				patterns: []string{"GET /foo.com/foo"}, // Unexpected pattern.
			},
		),
		gen(
			"valid pattern host",
			&condition{
				methods: []string{},
				hosts:   []string{"foo.com/foo"},
				paths:   []string{},
			},
			&action{
				// Irregular setting, but valid pattern.
				patterns: []string{"GET foo.com/foo/", "HEAD foo.com/foo/", "POST foo.com/foo/", "PUT foo.com/foo/", "PATCH foo.com/foo/",
					"DELETE foo.com/foo/", "CONNECT foo.com/foo/", "OPTIONS foo.com/foo/", "TRACE foo.com/foo/"},
			},
		),
		gen(
			"valid pattern method+host",
			&condition{
				methods: []string{http.MethodGet},
				hosts:   []string{"foo.com/foo"},
				paths:   []string{},
			},
			&action{
				patterns: []string{"GET foo.com/foo/"}, // Irregular setting, but valid pattern.
			},
		),
		gen(
			"valid pattern method+host+path",
			&condition{
				methods: []string{http.MethodGet},
				hosts:   []string{"foo.com/foo"},
				paths:   []string{"bar"},
			},
			&action{
				patterns: []string{"GET foo.com/foo/bar"}, // Irregular setting, but valid pattern.
			},
		),
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			patterns := generatePatterns(tt.C.methods, tt.C.hosts, tt.C.paths)

			opts := []cmp.Option{
				cmpopts.SortSlices(func(x, y string) bool { return x > y }),
			}
			testutil.Diff(t, tt.A.patterns, patterns, opts...)
		})
	}
}
