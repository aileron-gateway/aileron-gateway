package headercert

import (
	"crypto/x509"
	"regexp"
	"testing"

	v1 "github.com/aileron-gateway/aileron-gateway/apis/app/v1"
	"github.com/aileron-gateway/aileron-gateway/apis/kernel"
	"github.com/aileron-gateway/aileron-gateway/core"
	"github.com/aileron-gateway/aileron-gateway/kernel/api"
	"github.com/aileron-gateway/aileron-gateway/kernel/log"
	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
	utilhttp "github.com/aileron-gateway/aileron-gateway/util/http"
	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/reflect/protoreflect"
)

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

	defaultRoots, err := loadRootCert([]string{})
	if err != nil {
		t.Errorf("fail to load default RootCA: %v", err)
	}
	defaultOpts := x509.VerifyOptions{
		Roots: defaultRoots,
	}

	roots, err := loadRootCert([]string{rootCAPath})
	if err != nil {
		t.Errorf("fail to load RootCA: %v", err)
	}
	opts := x509.VerifyOptions{
		Roots: roots,
	}

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
				expect: &headerCert{
					eh:   utilhttp.GlobalErrorHandler(utilhttp.DefaultErrorHandlerName),
					opts: defaultOpts,
				},
			},
		),
		gen(
			"fail to get errorhandler",
			[]string{},
			[]string{},
			&condition{
				manifest: &v1.HeaderCertMiddleware{
					APIVersion: apiVersion,
					Kind:       kind,
					Metadata: &kernel.Metadata{
						Namespace: "defalut",
						Name:      "default",
					},
					Spec: &v1.HeaderCertMiddlewareSpec{
						ErrorHandler: &kernel.Reference{
							APIVersion: "wrong",
						},
					},
				},
			},
			&action{
				err:        core.ErrCoreGenCreateObject,
				errPattern: regexp.MustCompile(core.ErrPrefix + `failed to create HeaderCertMiddleware`),
			},
		),
		gen(
			"valid root cert path",
			[]string{},
			[]string{},
			&condition{
				manifest: &v1.HeaderCertMiddleware{
					APIVersion: apiVersion,
					Kind:       kind,
					Metadata: &kernel.Metadata{
						Namespace: "default",
						Name:      "default",
					},
					Spec: &v1.HeaderCertMiddlewareSpec{
						RootCAs: []string{rootCAPath},
					},
				},
			},
			&action{
				err: nil,
				expect: &headerCert{
					eh:   utilhttp.GlobalErrorHandler(utilhttp.DefaultErrorHandlerName),
					opts: opts,
				},
			},
		),
		gen(
			"invalid root cert path",
			[]string{},
			[]string{},
			&condition{
				manifest: &v1.HeaderCertMiddleware{
					APIVersion: apiVersion,
					Kind:       kind,
					Metadata: &kernel.Metadata{
						Namespace: "default",
						Name:      "default",
					},
					Spec: &v1.HeaderCertMiddlewareSpec{
						RootCAs: []string{"wrong"},
					},
				},
			},
			&action{
				err:        core.ErrCoreGenCreateObject,
				errPattern: regexp.MustCompile(core.ErrPrefix + `failed to create HeaderCertMiddleware`),
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
				cmp.AllowUnexported(headerCert{}),
				cmp.Comparer(testutil.ComparePointer[core.ErrorHandler]),
				cmp.Comparer(testutil.ComparePointer[log.Logger]),
			}
			testutil.DiffError(t, tt.A().err, tt.A().errPattern, err)
			testutil.Diff(t, tt.A().expect, got, opts...)
		})

	}
}

func TestLoadRootCert(t *testing.T) {

	t.Run("no root cert", func(t *testing.T) {
		_, err := loadRootCert([]string{"wrong"})
		if err == nil {
			t.Errorf("expected error, got nil")
		}
	})
	t.Run("invalid root cert", func(t *testing.T) {
		_, err := loadRootCert([]string{incompleteCertPath})
		if err == nil {
			t.Errorf("expected error, got nil")
		}
	})
}
