// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package header

import (
	"regexp"
	"testing"

	v1 "github.com/aileron-gateway/aileron-gateway/apis/app/v1"
	"github.com/aileron-gateway/aileron-gateway/apis/kernel"
	k "github.com/aileron-gateway/aileron-gateway/apis/kernel"
	"github.com/aileron-gateway/aileron-gateway/core"
	"github.com/aileron-gateway/aileron-gateway/internal/testutil"
	"github.com/aileron-gateway/aileron-gateway/kernel/api"
	"github.com/aileron-gateway/aileron-gateway/kernel/er"
	"github.com/aileron-gateway/aileron-gateway/kernel/log"
	"github.com/aileron-gateway/aileron-gateway/kernel/txtutil"
	httputil "github.com/aileron-gateway/aileron-gateway/util/http"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
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

	testRepl, _ := txtutil.NewStringReplacer(&kernel.ReplacerSpec{
		Replacers: &kernel.ReplacerSpec_Fixed{
			Fixed: &kernel.FixedReplacer{Value: "***"},
		},
	})

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"default manifest",
			&condition{
				manifest: Resource.Default(),
			},
			&action{
				err: nil,
				expect: &headerPolicy{
					eh: httputil.GlobalErrorHandler(httputil.DefaultErrorHandlerName),
				},
			},
		),
		gen(
			"with request policy",
			&condition{
				manifest: &v1.HeaderPolicyMiddleware{
					Metadata: &k.Metadata{},
					Spec: &v1.HeaderPolicyMiddlewareSpec{
						RequestPolicy: &v1.HeaderPolicySpec{
							Allows:  []string{"test-allows"},
							Removes: []string{"test-removes"},
							Add:     map[string]string{"test-add": "add-value"},
							Set:     map[string]string{"test-set": "set-value"},
							Rewrites: []*v1.HeaderRewriteSpec{
								{
									Name: "test-rewrites",
									Replacer: &kernel.ReplacerSpec{
										Replacers: &kernel.ReplacerSpec_Fixed{
											Fixed: &kernel.FixedReplacer{Value: "***"},
										},
									},
								},
							},
						},
					},
				},
			},
			&action{
				err: nil,
				expect: &headerPolicy{
					eh: httputil.GlobalErrorHandler(httputil.DefaultErrorHandlerName),
					reqPolicy: &policy{
						allows:  []string{"Test-Allows"},
						removes: []string{"Test-Removes"},
						add:     map[string]string{"Test-Add": "add-value"},
						set:     map[string]string{"Test-Set": "set-value"},
						repls:   map[string]txtutil.ReplaceFunc[string]{"Test-Rewrites": testRepl.Replace},
					},
				},
			},
		),
		gen(
			"with response policy",
			&condition{
				manifest: &v1.HeaderPolicyMiddleware{
					Metadata: &k.Metadata{},
					Spec: &v1.HeaderPolicyMiddlewareSpec{
						ResponsePolicy: &v1.HeaderPolicySpec{
							Allows:  []string{"test-allows"},
							Removes: []string{"test-removes"},
							Add:     map[string]string{"test-add": "add-value"},
							Set:     map[string]string{"test-set": "set-value"},
							Rewrites: []*v1.HeaderRewriteSpec{
								{
									Name: "test-rewrites",
									Replacer: &kernel.ReplacerSpec{
										Replacers: &kernel.ReplacerSpec_Fixed{
											Fixed: &kernel.FixedReplacer{Value: "***"},
										},
									},
								},
							},
						},
					},
				},
			},
			&action{
				err: nil,
				expect: &headerPolicy{
					eh: httputil.GlobalErrorHandler(httputil.DefaultErrorHandlerName),
					resPolicy: &policy{
						allows:  []string{"Test-Allows"},
						removes: []string{"Test-Removes"},
						add:     map[string]string{"Test-Add": "add-value"},
						set:     map[string]string{"Test-Set": "set-value"},
						repls:   map[string]txtutil.ReplaceFunc[string]{"Test-Rewrites": testRepl.Replace},
					},
				},
			},
		),
		gen(
			"invalid request policy",
			&condition{
				manifest: &v1.HeaderPolicyMiddleware{
					Metadata: &k.Metadata{},
					Spec: &v1.HeaderPolicyMiddlewareSpec{
						RequestPolicy: &v1.HeaderPolicySpec{
							Rewrites: []*v1.HeaderRewriteSpec{
								{
									Name: "test-rewrites",
									Replacer: &kernel.ReplacerSpec{
										Replacers: &kernel.ReplacerSpec_Regexp{
											Regexp: &kernel.RegexpReplacer{Pattern: "[0-9a-"},
										},
									},
								},
							},
						},
					},
				},
			},
			&action{
				err:        core.ErrCoreGenCreateObject,
				errPattern: regexp.MustCompile(core.ErrPrefix + `failed to create HeaderPolicyMiddleware`),
			},
		),
		gen(
			"invalid response policy",
			&condition{
				manifest: &v1.HeaderPolicyMiddleware{
					Metadata: &k.Metadata{},
					Spec: &v1.HeaderPolicyMiddlewareSpec{
						ResponsePolicy: &v1.HeaderPolicySpec{
							Rewrites: []*v1.HeaderRewriteSpec{
								{
									Name: "test-rewrites",
									Replacer: &kernel.ReplacerSpec{
										Replacers: &kernel.ReplacerSpec_Regexp{
											Regexp: &kernel.RegexpReplacer{Pattern: "[0-9a-"},
										},
									},
								},
							},
						},
					},
				},
			},
			&action{
				err:        core.ErrCoreGenCreateObject,
				errPattern: regexp.MustCompile(core.ErrPrefix + `failed to create HeaderPolicyMiddleware`),
			},
		),
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			server := api.NewContainerAPI()

			got, err := Resource.Create(server, tt.C.manifest)
			testutil.DiffError(t, tt.A.err, tt.A.errPattern, err)

			opts := []cmp.Option{
				cmp.AllowUnexported(headerPolicy{}, policy{}),
				cmp.Comparer(testutil.ComparePointer[log.Logger]),
				cmp.Comparer(testutil.ComparePointer[core.ErrorHandler]),
				cmp.Comparer(testutil.ComparePointer[txtutil.ReplaceFunc[string]]),
			}
			testutil.Diff(t, tt.A.expect, got, opts...)
		})
	}
}

func TestNewRewriters(t *testing.T) {
	type condition struct {
		specs []*v1.HeaderRewriteSpec
	}

	type action struct {
		repls map[string]txtutil.ReplaceFunc[string]
		err   error
	}

	testRepl, _ := txtutil.NewStringReplacer(&kernel.ReplacerSpec{
		Replacers: &kernel.ReplacerSpec_Fixed{
			Fixed: &kernel.FixedReplacer{Value: "***"},
		},
	})

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"nil map",
			&condition{
				specs: nil,
			},
			&action{
				repls: map[string]txtutil.ReplaceFunc[string]{},
			},
		),
		gen(
			"1 spec",
			&condition{
				specs: []*v1.HeaderRewriteSpec{
					{
						Name: "Test-Header",
						Replacer: &kernel.ReplacerSpec{
							Replacers: &kernel.ReplacerSpec_Fixed{
								Fixed: &kernel.FixedReplacer{Value: "***"},
							},
						},
					},
				},
			},
			&action{
				repls: map[string]txtutil.ReplaceFunc[string]{
					"Test-Header": testRepl.Replace,
				},
			},
		),
		gen(
			"2 specs",
			&condition{
				specs: []*v1.HeaderRewriteSpec{
					{
						Name: "Test-Foo",
						Replacer: &kernel.ReplacerSpec{
							Replacers: &kernel.ReplacerSpec_Fixed{
								Fixed: &kernel.FixedReplacer{Value: "***"},
							},
						},
					},
					{
						Name: "Test-Bar",
						Replacer: &kernel.ReplacerSpec{
							Replacers: &kernel.ReplacerSpec_Fixed{
								Fixed: &kernel.FixedReplacer{Value: "***"},
							},
						},
					},
				},
			},
			&action{
				repls: map[string]txtutil.ReplaceFunc[string]{
					"Test-Foo": testRepl.Replace,
					"Test-Bar": testRepl.Replace,
				},
			},
		),
		gen(
			"empty header name",
			&condition{
				specs: []*v1.HeaderRewriteSpec{
					{
						Name: "",
						Replacer: &kernel.ReplacerSpec{
							Replacers: &kernel.ReplacerSpec_Fixed{
								Fixed: &kernel.FixedReplacer{Value: "***"},
							},
						},
					},
				},
			},
			&action{
				repls: map[string]txtutil.ReplaceFunc[string]{},
			},
		),
		gen(
			"contain nil",
			&condition{
				specs: []*v1.HeaderRewriteSpec{
					nil,
				},
			},
			&action{
				repls: map[string]txtutil.ReplaceFunc[string]{},
			},
		),
		gen(
			"error spec",
			&condition{
				specs: []*v1.HeaderRewriteSpec{
					{
						Name: "Test-Header",
						Replacer: &kernel.ReplacerSpec{
							Replacers: &kernel.ReplacerSpec_Regexp{
								Regexp: &kernel.RegexpReplacer{Pattern: `[0-9a-`},
							},
						},
					},
				},
			},
			&action{
				err: &er.Error{
					Package:     txtutil.ErrPkg,
					Type:        txtutil.ErrTypeReplacer,
					Description: txtutil.ErrDscPattern,
				},
			},
		),
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			repls, err := newRewriters(tt.C.specs)
			testutil.Diff(t, tt.A.err, err, cmpopts.EquateErrors())

			opts := []cmp.Option{
				cmp.Comparer(testutil.ComparePointer[txtutil.ReplaceFunc[string]]),
			}
			testutil.Diff(t, tt.A.repls, repls, opts...)
		})
	}
}

func TestCanonicalSlice(t *testing.T) {
	type condition struct {
		headers []string
	}

	type action struct {
		headers []string
	}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"nil slice",
			&condition{
				headers: nil,
			},
			&action{
				headers: []string{},
			},
		),
		gen(
			"1 value",
			&condition{
				headers: []string{"foo"},
			},
			&action{
				headers: []string{"Foo"},
			},
		),
		gen(
			"2 values",
			&condition{
				headers: []string{"foo", "bar"},
			},
			&action{
				headers: []string{"Foo", "Bar"},
			},
		),
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			result := canonicalSlice(tt.C.headers)
			testutil.Diff(t, tt.A.headers, result)
		})
	}
}

func TestCanonicalMapKey(t *testing.T) {
	type condition struct {
		headers map[string]string
	}

	type action struct {
		headers map[string]string
	}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"nil map",
			&condition{
				headers: nil,
			},
			&action{
				headers: map[string]string{},
			},
		),
		gen(
			"1 value",
			&condition{
				headers: map[string]string{
					"foo": "foo",
				},
			},
			&action{
				headers: map[string]string{
					"Foo": "foo",
				},
			},
		),
		gen(
			"2 values",
			&condition{
				headers: map[string]string{
					"foo": "foo",
					"bar": "bar",
				},
			},
			&action{
				headers: map[string]string{
					"Foo": "foo",
					"Bar": "bar",
				},
			},
		),
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			result := canonicalMapKey(tt.C.headers)
			testutil.Diff(t, tt.A.headers, result)
		})
	}
}
