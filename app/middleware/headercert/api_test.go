// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package headercert

import (
	"crypto/x509"
	"errors"
	"io/fs"
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

	defaultRoots, _ := loadRootCert([]string{})
	defaultOpts := x509.VerifyOptions{
		Roots: defaultRoots,
	}

	roots, _ := loadRootCert([]string{rootCAPath})
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
					eh:         utilhttp.GlobalErrorHandler(utilhttp.DefaultErrorHandlerName),
					opts:       defaultOpts,
					certHeader: "X-SSL-Client-Cert",
					fpHeader:   "",
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
						Namespace: "default",
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
						RootCAs:            []string{rootCAPath},
						CertHeader:         "X-SSL-Client-Cert",
						FingerprintpHeader: "X-SSL-Client-Fingerprint",
					},
				},
			},
			&action{
				err: nil,
				expect: &headerCert{
					eh:         utilhttp.GlobalErrorHandler(utilhttp.DefaultErrorHandlerName),
					opts:       opts,
					certHeader: "X-SSL-Client-Cert",
					fpHeader:   "X-SSL-Client-Fingerprint",
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
						RootCAs:            []string{"wrong"},
						CertHeader:         "X-SSL-Client-Cert",
						FingerprintpHeader: "X-SSL-Client-Fingerprint",
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
				cmp.AllowUnexported(headerCert{}, x509.VerifyOptions{}),
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
		pool, err := loadRootCert([]string{"wrong"})
		if pool != nil {
			t.Errorf("expected nil pool, got %v", pool)
		}
		if !errors.Is(err, fs.ErrNotExist) {
			t.Errorf("expected error %v, got %v", fs.ErrNotExist, err)
		}
	})
	t.Run("invalid root cert", func(t *testing.T) {
		pool, err := loadRootCert([]string{incompleteCertPath})
		if pool != nil {
			t.Errorf("expected nil pool, got %v", pool)
		}
		if !errors.Is(err, ErrAddCert) {
			t.Errorf("expected error %v, got %v", ErrAddCert, err)
		}
	})
}
