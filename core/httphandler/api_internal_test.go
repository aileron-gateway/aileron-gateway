package httphandler

import (
	"context"
	"net/http"
	"regexp"
	"testing"

	v1 "github.com/aileron-gateway/aileron-gateway/apis/core/v1"
	k "github.com/aileron-gateway/aileron-gateway/apis/kernel"
	"github.com/aileron-gateway/aileron-gateway/core"
	"github.com/aileron-gateway/aileron-gateway/kernel/api"
	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
	utilhttp "github.com/aileron-gateway/aileron-gateway/util/http"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type testHandler struct {
	http.Handler
	patterns []string
	methods  []string
}

func (h *testHandler) Patterns() []string {
	return h.patterns
}
func (h *testHandler) Methods() []string {
	return h.methods
}

func TestCreate(t *testing.T) {
	type condition struct {
		manifest protoreflect.ProtoMessage
		server   api.API[*api.Request, *api.Response]
	}

	type action struct {
		expect     any
		err        any // error or errorutil.Kind
		errPattern *regexp.Regexp
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	cndDefaultManifest := tb.Condition("default manifest", "input default manifest")
	cndErrorReference := tb.Condition("error reference", "input an error reference to an object")
	cndWrongType := tb.Condition("wrong type", "input reference for wrong interface or type")
	actCheckError := tb.Action("check the returned error", "check that the returned error is the one expected")
	actCheckNoError := tb.Action("check no error", "check that there is no error returned")
	table := tb.Build()

	testServer := api.NewContainerAPI()
	postTestResource(testServer, "wrongHandler", "This is string, not http.handler")
	postTestResource(testServer, "handler",
		&testHandler{
			patterns: []string{"/test"},
			methods:  []string{http.MethodGet},
		},
	)

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"create with default manifest",
			[]string{cndDefaultManifest},
			[]string{actCheckNoError},
			&condition{
				manifest: Resource.Default(),
				server:   testServer,
			},
			&action{
				expect:     nil,
				err:        core.ErrCoreGenCreateObject,
				errPattern: regexp.MustCompile(core.ErrPrefix + `failed to create HTTPHandler`),
			},
		),
		gen(
			"create successful",
			[]string{},
			[]string{actCheckNoError},
			&condition{
				manifest: &v1.HTTPHandler{
					Metadata: &k.Metadata{},
					Spec: &v1.HTTPHandlerSpec{
						Handler: testResourceRef("handler"),
					},
				},
				server: testServer,
			},
			&action{
				expect: &handler{
					HandlerBase: &utilhttp.HandlerBase{
						AcceptPatterns: []string{"/test"},
						AcceptMethods:  []string{http.MethodGet},
					},
					Handler: &testHandler{
						patterns: []string{"/test"},
						methods:  []string{http.MethodGet},
					},
				},
				err: nil,
			},
		),
		gen(
			"create successful by joining pattern",
			[]string{},
			[]string{actCheckNoError},
			&condition{
				manifest: &v1.HTTPHandler{
					Metadata: &k.Metadata{},
					Spec: &v1.HTTPHandlerSpec{
						Handler: testResourceRef("handler"),
						Pattern: "/prefix",
					},
				},
				server: testServer,
			},
			&action{
				expect: &handler{
					HandlerBase: &utilhttp.HandlerBase{
						AcceptPatterns: []string{"/prefix/test"},
						AcceptMethods:  []string{http.MethodGet},
					},
					Handler: &testHandler{
						patterns: []string{"/test"},
						methods:  []string{http.MethodGet},
					},
				},
				err: nil,
			},
		),
		gen(
			"fail to get middleware",
			[]string{cndErrorReference},
			[]string{actCheckError},
			&condition{
				manifest: &v1.HTTPHandler{
					Metadata: &k.Metadata{},
					Spec: &v1.HTTPHandlerSpec{
						Middleware: []*k.Reference{
							{APIVersion: "wrong"},
						},
					},
				},
				server: testServer,
			},
			&action{
				expect:     nil,
				err:        core.ErrCoreGenCreateObject,
				errPattern: regexp.MustCompile(core.ErrPrefix + `failed to create HTTPHandler`),
			},
		),
		gen(
			"fail to type assert handler",
			[]string{cndWrongType},
			[]string{actCheckError},
			&condition{
				manifest: &v1.HTTPHandler{
					Metadata: &k.Metadata{},
					Spec: &v1.HTTPHandlerSpec{
						Handler: testResourceRef("wrongHandler"),
					},
				},
				server: testServer,
			},
			&action{
				expect:     nil,
				err:        core.ErrCoreGenCreateObject,
				errPattern: regexp.MustCompile(core.ErrPrefix + `failed to create HTTPHandler`),
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			got, err := Resource.Create(tt.C().server, tt.C().manifest)
			testutil.DiffError(t, tt.A().err, tt.A().errPattern, err)

			opts := []cmp.Option{
				cmp.AllowUnexported(handler{}, testHandler{}),
				cmpopts.SortSlices(func(a, b string) bool { return a < b }),
			}

			testutil.Diff(t, tt.A().expect, got, opts...)
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
