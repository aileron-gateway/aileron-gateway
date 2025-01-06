package static

import (
	stdcmp "cmp"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"regexp"
	"testing"

	v1 "github.com/aileron-gateway/aileron-gateway/apis/core/v1"
	k "github.com/aileron-gateway/aileron-gateway/apis/kernel"
	"github.com/aileron-gateway/aileron-gateway/core"
	"github.com/aileron-gateway/aileron-gateway/kernel/api"
	"github.com/aileron-gateway/aileron-gateway/kernel/log"
	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
	utilhttp "github.com/aileron-gateway/aileron-gateway/util/http"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// testDir is the path to the test data.
// This path can be changed by the environmental variable.
var testDir = stdcmp.Or(os.Getenv("TEST_DIR"), "../../test/")

func TestCreate(t *testing.T) {
	type condition struct {
		manifest protoreflect.ProtoMessage
		server   api.API[*api.Request, *api.Response]
		path     string
	}

	type action struct {
		expect     any
		err        any // error or errorutil.Kind
		errPattern *regexp.Regexp
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	cndDefaultManifest := tb.Condition("input default manifest", "input default manifest")
	cndErrorReference := tb.Condition("input error reference", "input an error reference to an object")
	actCheckError := tb.Action("check the returned error", "check that the returned error is the one expected")
	actCheckNoError := tb.Action("check no error", "check that there is no error returned")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"create with default manifest",
			[]string{cndDefaultManifest},
			[]string{actCheckNoError},
			&condition{
				manifest: Resource.Default(),
				server:   api.NewContainerAPI(),
			},
			&action{
				expect: &handler{
					HandlerBase: &utilhttp.HandlerBase{},
					Handler:     http.FileServer(http.Dir(".")),
					eh:          utilhttp.GlobalErrorHandler(utilhttp.DefaultErrorHandlerName),
				},
				err: nil,
			},
		),
		gen(
			"create with root dir",
			[]string{},
			[]string{actCheckNoError},
			&condition{
				manifest: &v1.StaticFileHandler{
					Metadata: &k.Metadata{},
					Spec: &v1.StaticFileHandlerSpec{
						RootDir: testDir + "ut/core/static/",
					},
				},
				server: api.NewContainerAPI(),
				path:   "/testdir/test.txt",
			},
			&action{
				expect: &handler{
					HandlerBase: &utilhttp.HandlerBase{},
					eh:          utilhttp.GlobalErrorHandler(utilhttp.DefaultErrorHandlerName),
				},
				err: nil,
			},
		),
		gen(
			"create with strip prefix",
			[]string{},
			[]string{actCheckNoError},
			&condition{
				manifest: &v1.StaticFileHandler{
					Metadata: &k.Metadata{},
					Spec: &v1.StaticFileHandlerSpec{
						StripPrefix: "/prefix",
						RootDir:     testDir + "ut/core/static/",
					},
				},
				server: api.NewContainerAPI(),
				path:   "/prefix/testdir/test.txt",
			},
			&action{
				expect: &handler{
					HandlerBase: &utilhttp.HandlerBase{},
					eh:          utilhttp.GlobalErrorHandler(utilhttp.DefaultErrorHandlerName),
				},
				err: nil,
			},
		),
		gen(
			"fail to get error handler",
			[]string{cndErrorReference},
			[]string{actCheckError},
			&condition{
				manifest: &v1.StaticFileHandler{
					Metadata: &k.Metadata{},
					Spec: &v1.StaticFileHandlerSpec{
						ErrorHandler: &k.Reference{
							APIVersion: "wrong",
						},
					},
				},
				server: api.NewContainerAPI(),
			},
			&action{
				expect:     nil,
				err:        core.ErrCoreGenCreateObject,
				errPattern: regexp.MustCompile(core.ErrPrefix + `failed to create StaticFileHandler`),
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
				cmp.Comparer(testutil.ComparePointer[log.Logger]),
				cmp.AllowUnexported(utilhttp.DefaultErrorHandler{}),
				cmp.AllowUnexported(handler{}),
				cmpopts.IgnoreFields(handler{}, "Handler"),
			}
			testutil.Diff(t, tt.A().expect, got, opts...)

			path := tt.C().path
			if path != "" {
				r, _ := http.NewRequest(http.MethodGet, "http://test.com"+path, nil)
				w := httptest.NewRecorder()

				h := got.(http.Handler)
				h.ServeHTTP(w, r)
				resp := w.Result()

				b, _ := io.ReadAll(resp.Body)
				testutil.Diff(t, "test", string(b))
			}
		})
	}
}
