package timeout

import (
	"net/http"
	"regexp"
	"testing"
	"time"

	v1 "github.com/aileron-gateway/aileron-gateway/apis/app/v1"
	corev1 "github.com/aileron-gateway/aileron-gateway/apis/core/v1"
	k "github.com/aileron-gateway/aileron-gateway/apis/kernel"
	"github.com/aileron-gateway/aileron-gateway/core"
	"github.com/aileron-gateway/aileron-gateway/kernel/api"
	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
	utilhttp "github.com/aileron-gateway/aileron-gateway/util/http"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func TestCreate(t *testing.T) {
	type condition struct {
		manifest protoreflect.ProtoMessage
	}

	type action struct {
		err        any // error or errorutil.Kind
		errPattern *regexp.Regexp
		expect     any
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	CndDefaultManifest := tb.Condition("input default manifest", "input default manifest")
	CndErrorErrorHandlerSet := tb.Condition("input error reference to errorhandler", "input error reference to errorhandler")
	CndSetAPITimeout := tb.Condition("input apiTimeout", "input apiTimeout")
	CndErrorPattern := tb.Condition("input pattern", "input pattern")
	ActCheckNoError := tb.Action("check no error was returned", "check no error was returned")
	ActCheckErrorMsg := tb.Action("check error message", "check the error messages that was returned")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"create with default manifest",
			[]string{CndDefaultManifest},
			[]string{ActCheckNoError},
			&condition{
				manifest: Resource.Default(),
			},
			&action{
				err: nil,
				expect: &timeout{
					eh:             utilhttp.GlobalErrorHandler(utilhttp.DefaultErrorHandlerName),
					defaultTimeout: 0,
					apiTimeouts:    []*apiTimeout{},
				},
			},
		),
		gen(
			"default timeout only",
			[]string{CndSetAPITimeout},
			[]string{ActCheckNoError},
			&condition{
				manifest: &v1.TimeoutMiddleware{
					Metadata: &k.Metadata{},
					Spec: &v1.TimeoutMiddlewareSpec{
						DefaultTimeout: 1000,
					},
				},
			},
			&action{
				err: nil,
				expect: &timeout{
					eh:             utilhttp.GlobalErrorHandler(utilhttp.DefaultErrorHandlerName),
					defaultTimeout: 1 * time.Second,
					apiTimeouts:    []*apiTimeout{},
				},
			},
		),
		gen(
			"default timeout negative",
			[]string{CndSetAPITimeout},
			[]string{ActCheckNoError},
			&condition{
				manifest: &v1.TimeoutMiddleware{
					Metadata: &k.Metadata{},
					Spec: &v1.TimeoutMiddlewareSpec{
						DefaultTimeout: -1,
					},
				},
			},
			&action{
				err: nil,
				expect: &timeout{
					eh:             utilhttp.GlobalErrorHandler(utilhttp.DefaultErrorHandlerName),
					defaultTimeout: -1 * time.Millisecond,
					apiTimeouts:    []*apiTimeout{},
				},
			},
		),
		gen(
			"all options",
			[]string{CndSetAPITimeout},
			[]string{ActCheckNoError},
			&condition{
				manifest: &v1.TimeoutMiddleware{
					Metadata: &k.Metadata{},
					Spec: &v1.TimeoutMiddlewareSpec{
						APITimeouts: []*v1.APITimeoutSpec{
							{
								Matcher: &k.MatcherSpec{
									MatchType: k.MatchType_Prefix,
									Patterns:  []string{"/test"},
								},
								Methods: []corev1.HTTPMethod{corev1.HTTPMethod_GET},
								Timeout: 1,
							},
						},
						DefaultTimeout: 1000,
					},
				},
			},
			&action{
				err: nil,
				expect: &timeout{
					eh:             utilhttp.GlobalErrorHandler(utilhttp.DefaultErrorHandlerName),
					defaultTimeout: 1 * time.Second,
					apiTimeouts: []*apiTimeout{
						{
							paths:   nil, // Not checked
							methods: []string{http.MethodGet},
							timeout: 1 * time.Millisecond,
						},
					},
				},
			},
		),
		gen(
			"api timeout negative",
			[]string{CndSetAPITimeout},
			[]string{ActCheckNoError},
			&condition{
				manifest: &v1.TimeoutMiddleware{
					Metadata: &k.Metadata{},
					Spec: &v1.TimeoutMiddlewareSpec{
						APITimeouts: []*v1.APITimeoutSpec{
							{
								Matcher: &k.MatcherSpec{
									MatchType: k.MatchType_Prefix,
									Patterns:  []string{"/test"},
								},
								Methods: []corev1.HTTPMethod{corev1.HTTPMethod_GET},
								Timeout: -1,
							},
						},
						DefaultTimeout: 1000,
					},
				},
			},
			&action{
				err: nil,
				expect: &timeout{
					eh:             utilhttp.GlobalErrorHandler(utilhttp.DefaultErrorHandlerName),
					defaultTimeout: 1 * time.Second,
					apiTimeouts:    []*apiTimeout{},
				},
			},
		),
		gen(
			"fail to get errorhandler",
			[]string{CndErrorErrorHandlerSet},
			[]string{ActCheckErrorMsg},
			&condition{
				manifest: &v1.TimeoutMiddleware{
					Metadata: &k.Metadata{},
					Spec: &v1.TimeoutMiddlewareSpec{
						ErrorHandler: &k.Reference{
							APIVersion: "wrong",
						},
					},
				},
			},
			&action{
				err:        core.ErrCoreGenCreateObject,
				errPattern: regexp.MustCompile(core.ErrPrefix + `failed to create TimeoutMiddleware`),
			},
		),
		gen(
			"default timeout negative",
			[]string{CndSetAPITimeout},
			[]string{ActCheckNoError},
			&condition{
				manifest: &v1.TimeoutMiddleware{
					Metadata: &k.Metadata{},
					Spec: &v1.TimeoutMiddlewareSpec{
						APITimeouts: []*v1.APITimeoutSpec{
							{
								Matcher: &k.MatcherSpec{
									MatchType: k.MatchType_Prefix,
									Patterns:  []string{"/test"},
								},
								Methods: []corev1.HTTPMethod{corev1.HTTPMethod_GET},
								Timeout: 1,
							},
						},
						DefaultTimeout: 1000,
					},
				},
			},
			&action{
				err: nil,
				expect: &timeout{
					eh:             utilhttp.GlobalErrorHandler(utilhttp.DefaultErrorHandlerName),
					defaultTimeout: 1 * time.Second,
					apiTimeouts: []*apiTimeout{
						{
							paths:   nil, // Not checked
							methods: []string{http.MethodGet},
							timeout: 1 * time.Millisecond,
						},
					},
				},
			},
		),
		gen(
			"input options to pattern error",
			[]string{CndErrorPattern},
			[]string{ActCheckErrorMsg},
			&condition{
				manifest: &v1.TimeoutMiddleware{
					Metadata: &k.Metadata{},
					Spec: &v1.TimeoutMiddlewareSpec{
						APITimeouts: []*v1.APITimeoutSpec{
							{
								Matcher: &k.MatcherSpec{
									MatchType: k.MatchType_Regex,
									Patterns:  []string{"[0-9"},
								},
								Methods: []corev1.HTTPMethod{corev1.HTTPMethod_GET},
								Timeout: 1,
							},
						},
					},
				},
			},
			&action{
				err:        core.ErrCoreGenCreateObject,
				errPattern: regexp.MustCompile(core.ErrPrefix + `failed to create TimeoutMiddleware`),
			},
		),
	}
	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			server := api.NewContainerAPI()
			a := &API{}
			got, err := a.Create(server, tt.C().manifest)
			opts := []cmp.Option{
				cmp.AllowUnexported(timeout{}),
				cmp.AllowUnexported(apiTimeout{}),
				cmp.Comparer(testutil.ComparePointer[core.ErrorHandler]),
				cmpopts.IgnoreFields(apiTimeout{}, "paths"),
			}
			testutil.DiffError(t, tt.A().err, tt.A().errPattern, err)
			testutil.Diff(t, tt.A().expect, got, opts...)
		})
	}
}
