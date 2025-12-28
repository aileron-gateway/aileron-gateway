// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package httplogger

import (
	"context"
	"io"
	"os"
	"regexp"
	"testing"
	"time"

	v1 "github.com/aileron-gateway/aileron-gateway/apis/core/v1"
	"github.com/aileron-gateway/aileron-gateway/apis/kernel"
	k "github.com/aileron-gateway/aileron-gateway/apis/kernel"
	"github.com/aileron-gateway/aileron-gateway/core"
	"github.com/aileron-gateway/aileron-gateway/kernel/api"
	"github.com/aileron-gateway/aileron-gateway/kernel/log"
	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
	utilhttp "github.com/aileron-gateway/aileron-gateway/util/http"
	"github.com/aileron-projects/go/ztext"
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

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"apply default values",
			&condition{
				manifest: Resource.Default(),
			},
			&action{
				manifest: &v1.HTTPLogger{
					APIVersion: apiVersion,
					Kind:       kind,
					Metadata: &k.Metadata{
						Namespace: "default",
						Name:      "default",
					},
					Spec: &v1.HTTPLoggerSpec{
						Timezone:   "Local",
						TimeFormat: "2006-01-02 15:04:05.000",
						Request: &v1.LoggingSpec{
							MaxContentLength: 1 << 12,
							MIMEs: []string{
								"application/json",
								"application/x-www-form-urlencoded",
								"application/xml",
								"application/soap+xml",
								"application/graphql+json",
								"text/plain",
								"text/html",
								"text/xml",
							},
						},
						Response: &v1.LoggingSpec{
							MaxContentLength: 1 << 12,
							MIMEs: []string{
								"application/json",
								"application/x-www-form-urlencoded",
								"application/xml",
								"application/soap+xml",
								"application/graphql+json",
								"text/plain",
								"text/html",
								"text/xml",
							},
						},
					},
				},
			},
		),
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			msg := Resource.Mutate(tt.C.manifest)

			opts := []cmp.Option{
				cmpopts.IgnoreUnexported(v1.HTTPLogger{}, v1.HTTPLoggerSpec{}, v1.LoggingSpec{}),
				cmpopts.IgnoreUnexported(k.Metadata{}, k.Status{}),
				cmpopts.IgnoreUnexported(v1.LogHeaderSpec{}, v1.LogBodySpec{}),
			}
			testutil.Diff(t, tt.A.manifest, msg, opts...)
		})
	}
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

	testServer := api.NewContainerAPI()
	postTestResource(testServer, "logger", log.GlobalLogger(log.DefaultLoggerName))

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"default manifest",
			&condition{
				manifest: Resource.Default(),
				server:   testServer,
			},
			&action{
				expect: &httpLogger{
					zone:    time.Local,
					timeFmt: "2006-01-02 15:04:05.000",
					req: &baseLogger{
						lg:         log.GlobalLogger(log.DefaultLoggerName),
						w:          os.Stderr,
						queries:    []stringReplFunc{},
						headers:    map[string][]stringReplFunc{},
						headerKeys: []string{},
						bodies:     map[string][]bytesReplFunc{},
						maxBody:    1 << 12,
					},
					res: &baseLogger{
						lg:         log.GlobalLogger(log.DefaultLoggerName),
						w:          os.Stderr,
						queries:    []stringReplFunc{},
						headers:    map[string][]stringReplFunc{},
						headerKeys: []string{},
						bodies:     map[string][]bytesReplFunc{},
						maxBody:    1 << 12,
					},
				},
				err: nil,
			},
		),
		gen(
			"create journal",
			&condition{
				manifest: &v1.HTTPLogger{
					Metadata: &kernel.Metadata{},
					Spec: &v1.HTTPLoggerSpec{
						Journal:  true,
						Request:  &v1.LoggingSpec{},
						Response: &v1.LoggingSpec{},
					},
				},
				server: testServer,
			},
			&action{
				expect: &journalLogger{
					zone: time.UTC,
					lg:   log.GlobalLogger(log.DefaultLoggerName),
					eh:   utilhttp.GlobalErrorHandler(utilhttp.DefaultErrorHandlerName),
					req: &baseLogger{
						lg:         log.GlobalLogger(log.DefaultLoggerName),
						w:          os.Stderr,
						queries:    []stringReplFunc{},
						headers:    map[string][]stringReplFunc{},
						headerKeys: []string{},
						bodies:     map[string][]bytesReplFunc{},
					},
					res: &baseLogger{
						lg:         log.GlobalLogger(log.DefaultLoggerName),
						w:          os.Stderr,
						queries:    []stringReplFunc{},
						headers:    map[string][]stringReplFunc{},
						headerKeys: []string{},
						bodies:     map[string][]bytesReplFunc{},
					},
				},
				err: nil,
			},
		),
		gen(
			"use logger",
			&condition{
				manifest: &v1.HTTPLogger{
					Metadata: &kernel.Metadata{},
					Spec: &v1.HTTPLoggerSpec{
						Logger:   testResourceRef("logger"),
						Request:  &v1.LoggingSpec{},
						Response: &v1.LoggingSpec{},
					},
				},
				server: testServer,
			},
			&action{
				expect: &httpLogger{
					zone: time.UTC,
					req: &baseLogger{
						lg:         log.GlobalLogger(log.DefaultLoggerName),
						w:          os.Stderr,
						queries:    []stringReplFunc{},
						headers:    map[string][]stringReplFunc{},
						headerKeys: []string{},
						bodies:     map[string][]bytesReplFunc{},
					},
					res: &baseLogger{
						lg:         log.GlobalLogger(log.DefaultLoggerName),
						w:          os.Stderr,
						queries:    []stringReplFunc{},
						headers:    map[string][]stringReplFunc{},
						headerKeys: []string{},
						bodies:     map[string][]bytesReplFunc{},
					},
				},
				err: nil,
			},
		),
		gen(
			"fail to get logger",
			&condition{
				manifest: &v1.HTTPLogger{
					Metadata: &kernel.Metadata{},
					Spec: &v1.HTTPLoggerSpec{
						Logger:   testResourceRef("none"),
						Request:  &v1.LoggingSpec{},
						Response: &v1.LoggingSpec{},
					},
				},
				server: testServer,
			},
			&action{
				expect:     nil,
				err:        core.ErrCoreGenCreateObject,
				errPattern: regexp.MustCompile(core.ErrPrefix + `failed to create HTTPLogger`),
			},
		),
		gen(
			"fail to parse timezone",
			&condition{
				manifest: &v1.HTTPLogger{
					Metadata: &kernel.Metadata{},
					Spec: &v1.HTTPLoggerSpec{
						Timezone: "Unknown",
						Request:  &v1.LoggingSpec{},
						Response: &v1.LoggingSpec{},
					},
				},
				server: testServer,
			},
			&action{
				expect:     nil,
				err:        core.ErrCoreGenCreateObject,
				errPattern: regexp.MustCompile(core.ErrPrefix + `failed to create HTTPLogger`),
			},
		),
		gen(
			"fail to get error handler",
			&condition{
				manifest: &v1.HTTPLogger{
					Metadata: &kernel.Metadata{},
					Spec: &v1.HTTPLoggerSpec{
						ErrorHandler: testResourceRef("none"),
						Request:      &v1.LoggingSpec{},
						Response:     &v1.LoggingSpec{},
					},
				},
				server: testServer,
			},
			&action{
				expect:     nil,
				err:        core.ErrCoreGenCreateObject,
				errPattern: regexp.MustCompile(core.ErrPrefix + `failed to create HTTPLogger`),
			},
		),
		gen(
			"fail to create req",
			&condition{
				manifest: &v1.HTTPLogger{
					Metadata: &kernel.Metadata{},
					Spec: &v1.HTTPLoggerSpec{
						Request: &v1.LoggingSpec{
							Queries: []*k.ReplacerSpec{
								{
									Replacers: &k.ReplacerSpec_Regexp{
										Regexp: &k.RegexpReplacer{Pattern: "[0-9"},
									},
								},
							},
						},
						Response: &v1.LoggingSpec{},
					},
				},
				server: testServer,
			},
			&action{
				expect:     nil,
				err:        core.ErrCoreGenCreateObject,
				errPattern: regexp.MustCompile(core.ErrPrefix + `failed to create HTTPLogger`),
			},
		),
		gen(
			"fail to create req",
			&condition{
				manifest: &v1.HTTPLogger{
					Metadata: &kernel.Metadata{},
					Spec: &v1.HTTPLoggerSpec{
						Request: &v1.LoggingSpec{},
						Response: &v1.LoggingSpec{
							Queries: []*k.ReplacerSpec{
								{
									Replacers: &k.ReplacerSpec_Regexp{
										Regexp: &k.RegexpReplacer{Pattern: "[0-9"},
									},
								},
							},
						},
					},
				},
				server: testServer,
			},
			&action{
				expect:     nil,
				err:        core.ErrCoreGenCreateObject,
				errPattern: regexp.MustCompile(core.ErrPrefix + `failed to create HTTPLogger`),
			},
		),
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			got, err := Resource.Create(tt.C.server, tt.C.manifest)
			testutil.DiffError(t, tt.A.err, tt.A.errPattern, err)

			opts := []cmp.Option{
				cmp.Comparer(testutil.ComparePointer[log.Logger]),
				cmp.Comparer(testutil.ComparePointer[core.ErrorHandler]),
				cmp.Comparer(testutil.ComparePointer[io.Writer]),
				cmp.Comparer(testutil.ComparePointer[*time.Location]),
				cmp.AllowUnexported(baseLogger{}, httpLogger{}, journalLogger{}),
				cmpopts.IgnoreUnexported(ztext.Template{}),
				cmpopts.SortSlices(func(a, b string) bool { return a < b }),
			}

			testutil.Diff(t, tt.A.expect, got, opts...)
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
