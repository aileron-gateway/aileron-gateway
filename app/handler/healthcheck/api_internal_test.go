package healthcheck

import (
	"net/http"
	"regexp"
	"testing"
	"time"

	v1 "github.com/aileron-gateway/aileron-gateway/apis/app/v1"
	k "github.com/aileron-gateway/aileron-gateway/apis/kernel"
	"github.com/aileron-gateway/aileron-gateway/app"
	"github.com/aileron-gateway/aileron-gateway/core"
	"github.com/aileron-gateway/aileron-gateway/kernel/api"
	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
	utilhttp "github.com/aileron-gateway/aileron-gateway/util/http"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/testing/protocmp"
)

func TestMutate(t *testing.T) {
	type condition struct {
		manifest protoreflect.ProtoMessage
	}

	type action struct {
		err        any
		errPattern *regexp.Regexp
		expect     protoreflect.ProtoMessage
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())

	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"default patterns applied",
			[]string{},
			[]string{},
			&condition{
				manifest: &v1.HealthCheckHandler{
					Metadata: &k.Metadata{},
					Spec:     &v1.HealthCheckHandlerSpec{},
				},
			},
			&action{
				err: nil,
				expect: &v1.HealthCheckHandler{
					Metadata: &k.Metadata{},
					Spec: &v1.HealthCheckHandlerSpec{
						Patterns: []string{"/healthz"},
					},
				},
			},
		),
		gen(
			"custom patterns, no changes applied",
			[]string{},
			[]string{},
			&condition{
				manifest: &v1.HealthCheckHandler{
					Metadata: &k.Metadata{},
					Spec: &v1.HealthCheckHandlerSpec{
						Patterns: []string{"/custom"},
					},
				},
			},
			&action{
				err: nil,
				expect: &v1.HealthCheckHandler{
					Metadata: &k.Metadata{},
					Spec: &v1.HealthCheckHandlerSpec{
						Patterns: []string{"/custom"},
					},
				},
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			a := &API{}
			got := a.Mutate(tt.C().manifest)

			opts := []cmp.Option{
				protocmp.Transform(),
				cmpopts.IgnoreUnexported(k.Metadata{}, k.Status{}),
			}

			testutil.DiffError(t, tt.A().err, tt.A().errPattern, nil)
			testutil.Diff(t, tt.A().expect, got, opts...)
		})
	}
}

func TestCreate(t *testing.T) {
	type condition struct {
		manifest protoreflect.ProtoMessage
	}

	type action struct {
		err        any
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
				expect: &healthCheck{
					HandlerBase: &utilhttp.HandlerBase{
						AcceptPatterns: nil,
						AcceptMethods:  nil,
					},
					eh:       utilhttp.GlobalErrorHandler(utilhttp.DefaultErrorHandlerName),
					timeout:  30 * time.Second,
					checkers: []app.HealthChecker{},
				},
			},
		),
		gen(
			"fail to get errorhandler",
			[]string{},
			[]string{},
			&condition{
				manifest: &v1.HealthCheckHandler{
					Metadata: &k.Metadata{},
					Spec: &v1.HealthCheckHandlerSpec{
						ErrorHandler: &k.Reference{
							APIVersion: "wrong",
						},
					},
				},
			},
			&action{
				err:        core.ErrCoreGenCreateObject,
				errPattern: regexp.MustCompile(core.ErrPrefix + `failed to create HealthCheckHandler`),
			},
		),
		gen(
			"input custom timeout",
			[]string{},
			[]string{},
			&condition{
				manifest: &v1.HealthCheckHandler{
					Metadata: &k.Metadata{},
					Spec: &v1.HealthCheckHandlerSpec{
						Timeout: 10,
					},
				},
			},
			&action{
				err: nil,
				expect: &healthCheck{
					HandlerBase: &utilhttp.HandlerBase{
						AcceptPatterns: nil,
						AcceptMethods:  nil,
					},
					eh:       utilhttp.GlobalErrorHandler(utilhttp.DefaultErrorHandlerName),
					timeout:  10 * time.Second,
					checkers: []app.HealthChecker{},
				},
			},
		),
		gen(
			"fail to get external probes",
			[]string{},
			[]string{},
			&condition{
				manifest: &v1.HealthCheckHandler{
					Metadata: &k.Metadata{},
					Spec: &v1.HealthCheckHandlerSpec{
						ExternalProbes: []*k.Reference{
							{
								APIVersion: "wrong",
							},
						},
					},
				},
			},
			&action{
				err:        core.ErrCoreGenCreateObject,
				errPattern: regexp.MustCompile(core.ErrPrefix + `failed to create HealthCheckHandler`),
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
				cmp.AllowUnexported(healthCheck{}),
				cmpopts.IgnoreInterfaces(struct{ core.Matcher[*http.Request] }{}),
				cmp.Comparer(testutil.ComparePointer[core.ErrorHandler]),
			}
			testutil.DiffError(t, tt.A().err, tt.A().errPattern, err)
			testutil.Diff(t, tt.A().expect, got, opts...)
		})
	}
}
