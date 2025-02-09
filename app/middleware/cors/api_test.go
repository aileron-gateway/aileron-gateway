package cors

import (
	"regexp"
	"testing"

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

func TestMutate(t *testing.T) {
	type condition struct {
		manifest protoreflect.ProtoMessage
	}

	type action struct {
		manifest protoreflect.ProtoMessage
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"apply default values",
			[]string{},
			[]string{},
			&condition{
				manifest: Resource.Default(),
			},
			&action{
				manifest: &v1.CORSMiddleware{
					APIVersion: apiVersion,
					Kind:       kind,
					Metadata: &k.Metadata{
						Namespace: "default",
						Name:      "default",
					},
					Spec: &v1.CORSMiddlewareSpec{
						CORSPolicy: &v1.CORSPolicySpec{
							AllowedOrigins: []string{"*"},
							AllowedMethods: []corev1.HTTPMethod{corev1.HTTPMethod_POST, corev1.HTTPMethod_GET, corev1.HTTPMethod_OPTIONS},
							AllowedHeaders: []string{"Content-Type", "X-Requested-With"},
						},
					},
				},
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			msg := Resource.Mutate(tt.C().manifest)

			opts := []cmp.Option{
				cmpopts.IgnoreUnexported(v1.CORSMiddleware{}, v1.CORSMiddlewareSpec{}),
				cmpopts.IgnoreUnexported(v1.CORSPolicySpec{}),
				cmpopts.IgnoreUnexported(k.Metadata{}, k.Status{}),
			}
			testutil.Diff(t, tt.A().manifest, msg, opts...)
		})
	}
}
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
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"create with default manifest",
			[]string{},
			[]string{},
			&condition{
				manifest: Resource.Default(),
			},
			&action{
				err: nil,
				expect: &cors{
					eh: utilhttp.GlobalErrorHandler(utilhttp.DefaultErrorHandlerName),
					policy: &corsPolicy{
						maxAge:                "0",
						allowPrivateNetwork:   false,
						disableWildCardOrigin: false,
					},
				},
			},
		),
		gen(
			"open policy unsafe-none",
			[]string{},
			[]string{},
			&condition{
				manifest: &v1.CORSMiddleware{
					Metadata: &k.Metadata{},
					Spec: &v1.CORSMiddlewareSpec{
						CORSPolicy: &v1.CORSPolicySpec{
							CORSOpenerPolicy: v1.CORSOpenerPolicy_OpenerUnsafeNone,
						},
					},
				},
			},
			&action{
				err: nil,
				expect: &cors{
					eh: utilhttp.GlobalErrorHandler(utilhttp.DefaultErrorHandlerName),
					policy: &corsPolicy{
						openerPolicy:          "unsafe-none",
						maxAge:                "0",
						allowPrivateNetwork:   false,
						disableWildCardOrigin: false,
					},
				},
			},
		),
		gen(
			"open policy same-origin-allow-popups",
			[]string{},
			[]string{},
			&condition{
				manifest: &v1.CORSMiddleware{
					Metadata: &k.Metadata{},
					Spec: &v1.CORSMiddlewareSpec{
						CORSPolicy: &v1.CORSPolicySpec{
							CORSOpenerPolicy: v1.CORSOpenerPolicy_OpenerSameOriginAllowPopups,
						},
					},
				},
			},
			&action{
				err: nil,
				expect: &cors{
					eh: utilhttp.GlobalErrorHandler(utilhttp.DefaultErrorHandlerName),
					policy: &corsPolicy{
						openerPolicy:          "same-origin-allow-popups",
						maxAge:                "0",
						allowPrivateNetwork:   false,
						disableWildCardOrigin: false,
					},
				},
			},
		),
		gen(
			"open policy same-origin",
			[]string{},
			[]string{},
			&condition{
				manifest: &v1.CORSMiddleware{
					Metadata: &k.Metadata{},
					Spec: &v1.CORSMiddlewareSpec{
						CORSPolicy: &v1.CORSPolicySpec{
							CORSOpenerPolicy: v1.CORSOpenerPolicy_OpenerSameOrigin,
						},
					},
				},
			},
			&action{
				err: nil,
				expect: &cors{
					eh: utilhttp.GlobalErrorHandler(utilhttp.DefaultErrorHandlerName),
					policy: &corsPolicy{
						openerPolicy:          "same-origin",
						maxAge:                "0",
						allowPrivateNetwork:   false,
						disableWildCardOrigin: false,
					},
				},
			},
		),
		gen(
			"resource policy same-saite",
			[]string{},
			[]string{},
			&condition{
				manifest: &v1.CORSMiddleware{
					Metadata: &k.Metadata{},
					Spec: &v1.CORSMiddlewareSpec{
						CORSPolicy: &v1.CORSPolicySpec{
							CORSResourcePolicy: v1.CORSResourcePolicy_ResourceSameSite,
						},
					},
				},
			},
			&action{
				err: nil,
				expect: &cors{
					eh: utilhttp.GlobalErrorHandler(utilhttp.DefaultErrorHandlerName),
					policy: &corsPolicy{
						resourcePolicy:        "same-site",
						maxAge:                "0",
						allowPrivateNetwork:   false,
						disableWildCardOrigin: false,
					},
				},
			},
		),
		gen(
			"resource policy same-origin",
			[]string{},
			[]string{},
			&condition{
				manifest: &v1.CORSMiddleware{
					Metadata: &k.Metadata{},
					Spec: &v1.CORSMiddlewareSpec{
						CORSPolicy: &v1.CORSPolicySpec{
							CORSResourcePolicy: v1.CORSResourcePolicy_ResourceSameOrigin,
						},
					},
				},
			},
			&action{
				err: nil,
				expect: &cors{
					eh: utilhttp.GlobalErrorHandler(utilhttp.DefaultErrorHandlerName),
					policy: &corsPolicy{
						resourcePolicy:        "same-origin",
						maxAge:                "0",
						allowPrivateNetwork:   false,
						disableWildCardOrigin: false,
					},
				},
			},
		),
		gen(
			"resource policy cross-origin",
			[]string{},
			[]string{},
			&condition{
				manifest: &v1.CORSMiddleware{
					Metadata: &k.Metadata{},
					Spec: &v1.CORSMiddlewareSpec{
						CORSPolicy: &v1.CORSPolicySpec{
							CORSResourcePolicy: v1.CORSResourcePolicy_ResourceCrossOrigin,
						},
					},
				},
			},
			&action{
				err: nil,
				expect: &cors{
					eh: utilhttp.GlobalErrorHandler(utilhttp.DefaultErrorHandlerName),
					policy: &corsPolicy{
						resourcePolicy:        "cross-origin",
						maxAge:                "0",
						allowPrivateNetwork:   false,
						disableWildCardOrigin: false,
					},
				},
			},
		),
		gen(
			"embedder policy unsafe-none",
			[]string{},
			[]string{},
			&condition{
				manifest: &v1.CORSMiddleware{
					Metadata: &k.Metadata{},
					Spec: &v1.CORSMiddlewareSpec{
						CORSPolicy: &v1.CORSPolicySpec{
							CORSEmbedderPolicy: v1.CORSEmbedderPolicy_EmbedderUnsafeNone,
						},
					},
				},
			},
			&action{
				err: nil,
				expect: &cors{
					eh: utilhttp.GlobalErrorHandler(utilhttp.DefaultErrorHandlerName),
					policy: &corsPolicy{
						embedderPolicy:        "unsafe-none",
						maxAge:                "0",
						allowPrivateNetwork:   false,
						disableWildCardOrigin: false,
					},
				},
			},
		),
		gen(
			"embedder policy require-corp",
			[]string{},
			[]string{},
			&condition{
				manifest: &v1.CORSMiddleware{
					Metadata: &k.Metadata{},
					Spec: &v1.CORSMiddlewareSpec{
						CORSPolicy: &v1.CORSPolicySpec{
							CORSEmbedderPolicy: v1.CORSEmbedderPolicy_EmbedderRequireCorp,
						},
					},
				},
			},
			&action{
				err: nil,
				expect: &cors{
					eh: utilhttp.GlobalErrorHandler(utilhttp.DefaultErrorHandlerName),
					policy: &corsPolicy{
						embedderPolicy:        "require-corp",
						maxAge:                "0",
						allowPrivateNetwork:   false,
						disableWildCardOrigin: false,
					},
				},
			},
		),
		gen(
			"embedder policy require-corp",
			[]string{},
			[]string{},
			&condition{
				manifest: &v1.CORSMiddleware{
					Metadata: &k.Metadata{},
					Spec: &v1.CORSMiddlewareSpec{
						CORSPolicy: &v1.CORSPolicySpec{
							CORSEmbedderPolicy: v1.CORSEmbedderPolicy_EmbedderCredentialless,
						},
					},
				},
			},
			&action{
				err: nil,
				expect: &cors{
					eh: utilhttp.GlobalErrorHandler(utilhttp.DefaultErrorHandlerName),
					policy: &corsPolicy{
						embedderPolicy:        "credentialless",
						maxAge:                "0",
						allowPrivateNetwork:   false,
						disableWildCardOrigin: false,
					},
				},
			},
		),
		gen(
			"fail to get errorhandler",
			[]string{},
			[]string{},
			&condition{
				manifest: &v1.CORSMiddleware{
					Metadata: &k.Metadata{},
					Spec: &v1.CORSMiddlewareSpec{
						ErrorHandler: &k.Reference{
							APIVersion: "wrong",
						},
					},
				},
			},
			&action{
				err:        core.ErrCoreGenCreateObject,
				errPattern: regexp.MustCompile(core.ErrPrefix + `failed to create CORSMiddleware`),
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
				cmp.AllowUnexported(cors{}),
				cmp.AllowUnexported(corsPolicy{}),
				cmp.Comparer(testutil.ComparePointer[core.ErrorHandler]),
			}
			testutil.DiffError(t, tt.A().err, tt.A().errPattern, err)
			testutil.Diff(t, tt.A().expect, got, opts...)
		})
	}
}
