// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package entrypoint

import (
	"context"
	"errors"
	"regexp"
	"sync"
	"sync/atomic"
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

func TestMutate(t *testing.T) {
	type condition struct {
		manifest protoreflect.ProtoMessage
	}

	type action struct {
		manifest protoreflect.ProtoMessage
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	cndDefault := tb.Condition("default", "input default config")
	cndWrongMetadata := tb.Condition("wrong metadata", "name and name space is wrong")
	actCheckMutated := tb.Action("check mutated", "check that the intended fields are mutated")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"mutate default",
			[]string{cndDefault},
			[]string{},
			&condition{
				manifest: Resource.Default(),
			},
			&action{
				manifest: &v1.Entrypoint{
					APIVersion: apiVersion,
					Kind:       kind,
					Metadata: &k.Metadata{
						Namespace: ".entrypoint",
						Name:      ".entrypoint",
					},
					Spec: &v1.EntrypointSpec{},
				},
			},
		),
		gen(
			"mutate metadata",
			[]string{cndWrongMetadata},
			[]string{actCheckMutated},
			&condition{
				manifest: &v1.Entrypoint{
					APIVersion: apiVersion,
					Kind:       kind,
					Metadata: &k.Metadata{
						Namespace: "test",
						Name:      "test",
					},
					Spec: &v1.EntrypointSpec{},
				},
			},
			&action{
				manifest: &v1.Entrypoint{
					APIVersion: apiVersion,
					Kind:       kind,
					Metadata: &k.Metadata{
						Namespace: ".entrypoint",
						Name:      ".entrypoint",
					},
					Spec: &v1.EntrypointSpec{},
				},
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			manifest := Resource.Mutate(tt.C().manifest)

			opts := []cmp.Option{
				cmpopts.IgnoreUnexported(v1.Entrypoint{}, v1.EntrypointSpec{}),
				cmpopts.IgnoreUnexported(k.Metadata{}, k.Status{}, k.Reference{}),
			}
			testutil.Diff(t, tt.A().manifest, manifest, opts...)
		})
	}
}

type testInitializer struct {
	name   string
	called bool
	err    error
}

func (f *testInitializer) Initialize() error {
	f.called = true
	return f.err
}

type testFinalizer struct {
	name   string
	called bool
	err    error
}

func (f *testFinalizer) Finalize() error {
	f.called = true
	return f.err
}

func TestCreate(t *testing.T) {
	type condition struct {
		manifest protoreflect.ProtoMessage
	}

	type action struct {
		expect     any
		err        any // error or errorutil.Kind
		errPattern *regexp.Regexp
	}

	noopLogger := &struct{ log.Logger }{}
	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	cndDefaultManifest := tb.Condition("default manifest", "input default manifest")
	cndErrorReference := tb.Condition("error reference", "input an error reference to an object")
	actCheckExpected := tb.Action("check returned wait group", "check that an expected object was returned")
	actCheckError := tb.Action("check the returned error", "check that the returned error is the one expected")
	actCheckNoError := tb.Action("check no error", "check that there is no error returned")
	table := tb.Build()

	testServer := api.NewContainerAPI()
	postTestResource(testServer, "nilRunner", &testRunner{})
	postTestResource(testServer, "errRunner", &testRunner{err: errors.New("test")})
	postTestResource(testServer, "nilFinalizer", &testFinalizer{name: "nil"})
	postTestResource(testServer, "errFinalizer", &testFinalizer{err: errors.New("test")})
	postTestResource(testServer, "noopLogger", noopLogger)
	testEH := &utilhttp.DefaultErrorHandler{}
	postTestResource(testServer, "errorHandler", testEH)

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"create with default manifest",
			[]string{cndDefaultManifest},
			[]string{actCheckExpected, actCheckNoError},
			&condition{
				manifest: Resource.Default(),
			},
			&action{
				expect: &channelGroup{
					lg:           log.GlobalLogger(log.DefaultLoggerName),
					runners:      []core.Runner{},
					initializers: []core.Initializer{},
					finalizers:   []core.Finalizer{},
				},
				err: nil,
			},
		),
		gen(
			"create with logger",
			[]string{cndDefaultManifest},
			[]string{actCheckExpected, actCheckNoError},
			&condition{
				manifest: &v1.Entrypoint{
					Metadata: &k.Metadata{},
					Spec: &v1.EntrypointSpec{
						Loggers: []*k.Reference{testResourceRef("noopLogger")},
					},
				},
			},
			&action{
				expect: &channelGroup{
					lg:           log.GlobalLogger(log.DefaultLoggerName),
					runners:      []core.Runner{},
					initializers: []core.Initializer{},
					finalizers:   []core.Finalizer{},
				},
				err: nil,
			},
		),
		gen(
			"create with default logger",
			[]string{cndDefaultManifest},
			[]string{actCheckExpected, actCheckNoError},
			&condition{
				manifest: &v1.Entrypoint{
					Metadata: &k.Metadata{},
					Spec: &v1.EntrypointSpec{
						DefaultLogger: testResourceRef("noopLogger"),
					},
				},
			},
			&action{
				expect: &channelGroup{
					lg:           noopLogger,
					runners:      []core.Runner{},
					initializers: []core.Initializer{},
					finalizers:   []core.Finalizer{},
				},
				err: nil,
			},
		),
		gen(
			"create with default error handler",
			[]string{cndDefaultManifest},
			[]string{actCheckExpected, actCheckNoError},
			&condition{
				manifest: &v1.Entrypoint{
					Metadata: &k.Metadata{},
					Spec: &v1.EntrypointSpec{
						DefaultErrorHandler: testResourceRef("errorHandler"),
					},
				},
			},
			&action{
				expect: &channelGroup{
					lg:           log.GlobalLogger(log.DefaultLoggerName),
					runners:      []core.Runner{},
					initializers: []core.Initializer{},
					finalizers:   []core.Finalizer{},
				},
				err: nil,
			},
		),
		gen(
			"use metadata logger",
			[]string{cndDefaultManifest},
			[]string{actCheckExpected, actCheckNoError},
			&condition{
				manifest: &v1.Entrypoint{
					Metadata: &k.Metadata{
						Logger: "core/v1/Container/test/noopLogger",
					},
					Spec: &v1.EntrypointSpec{
						Loggers: []*k.Reference{testResourceRef("noopLogger")},
					},
				},
			},
			&action{
				expect: &channelGroup{
					lg:           noopLogger,
					runners:      []core.Runner{},
					initializers: []core.Initializer{},
					finalizers:   []core.Finalizer{},
				},
				err: nil,
			},
		),
		gen(
			"create with normal runner/channelgroup",
			[]string{},
			[]string{actCheckExpected, actCheckNoError},
			&condition{
				manifest: &v1.Entrypoint{
					Metadata: &k.Metadata{},
					Spec: &v1.EntrypointSpec{
						Runners: []*k.Reference{
							testResourceRef("nilRunner"),
						},
					},
				},
			},
			&action{
				expect: &channelGroup{
					lg: log.GlobalLogger(log.DefaultLoggerName),
					runners: []core.Runner{
						&testRunner{},
					},
					initializers: []core.Initializer{},
					finalizers:   []core.Finalizer{},
				},
				err: nil,
			},
		),
		gen(
			"create with finalizers",
			[]string{},
			[]string{actCheckExpected, actCheckNoError},
			&condition{
				manifest: &v1.Entrypoint{
					Metadata: &k.Metadata{},
					Spec: &v1.EntrypointSpec{
						Runners:    []*k.Reference{testResourceRef("nilRunner")},
						Finalizers: []*k.Reference{testResourceRef("nilFinalizer")},
					},
				},
			},
			&action{
				expect: &channelGroup{
					lg:           log.GlobalLogger(log.DefaultLoggerName),
					runners:      []core.Runner{&testRunner{}},
					initializers: []core.Initializer{},
					finalizers:   []core.Finalizer{&testFinalizer{name: "nil"}},
				},
				err: nil,
			},
		),
		gen(
			"fail to get logger",
			[]string{cndErrorReference},
			[]string{actCheckError},
			&condition{
				manifest: &v1.Entrypoint{
					Metadata: &k.Metadata{},
					Spec: &v1.EntrypointSpec{
						Loggers: []*k.Reference{
							{APIVersion: "wrong"},
						},
					},
				},
			},
			&action{
				expect:     nil,
				err:        core.ErrCoreGenCreateObject,
				errPattern: regexp.MustCompile(core.ErrPrefix + `failed to create Entrypoint`),
			},
		),
		gen(
			"fail to get default logger",
			[]string{cndErrorReference},
			[]string{actCheckError},
			&condition{
				manifest: &v1.Entrypoint{
					Metadata: &k.Metadata{},
					Spec: &v1.EntrypointSpec{
						DefaultLogger: &k.Reference{APIVersion: "wrong"},
					},
				},
			},
			&action{
				expect:     nil,
				err:        core.ErrCoreGenCreateObject,
				errPattern: regexp.MustCompile(core.ErrPrefix + `failed to create Entrypoint`),
			},
		),
		gen(
			"fail to get error handler",
			[]string{cndErrorReference},
			[]string{actCheckError},
			&condition{
				manifest: &v1.Entrypoint{
					Metadata: &k.Metadata{},
					Spec: &v1.EntrypointSpec{
						DefaultErrorHandler: &k.Reference{APIVersion: "wrong"},
					},
				},
			},
			&action{
				expect:     nil,
				err:        core.ErrCoreGenCreateObject,
				errPattern: regexp.MustCompile(core.ErrPrefix + `failed to create Entrypoint`),
			},
		),
		gen(
			"fail to get runner",
			[]string{cndErrorReference},
			[]string{actCheckError},
			&condition{
				manifest: &v1.Entrypoint{
					Metadata: &k.Metadata{},
					Spec: &v1.EntrypointSpec{
						Runners: []*k.Reference{
							{APIVersion: "wrong"},
						},
					},
				},
			},
			&action{
				expect:     nil,
				err:        core.ErrCoreGenCreateObject,
				errPattern: regexp.MustCompile(core.ErrPrefix + `failed to create Entrypoint`),
			},
		),
		gen(
			"fail to get finalizer",
			[]string{cndErrorReference},
			[]string{actCheckError},
			&condition{
				manifest: &v1.Entrypoint{
					Metadata: &k.Metadata{},
					Spec: &v1.EntrypointSpec{
						Finalizers: []*k.Reference{nil},
					},
				},
			},
			&action{
				expect:     nil,
				err:        core.ErrCoreGenCreateObject,
				errPattern: regexp.MustCompile(core.ErrPrefix + `failed to create Entrypoint`),
			},
		),
	}

	testutil.Register(table, testCases...)

	lgTmp := log.GlobalLogger(log.DefaultLoggerName)
	ehTmp := utilhttp.GlobalErrorHandler(utilhttp.DefaultErrorHandlerName)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			defer func() {
				log.SetGlobalLogger("core/v1/Container/test/noopLogger", nil) // Remove
				log.SetGlobalLogger(log.DefaultLoggerName, lgTmp)
				log.SetGlobalLogger("core/v1/Container/test/errorHandler", nil) // Remove
				utilhttp.SetGlobalErrorHandler(utilhttp.DefaultErrorHandlerName, ehTmp)
			}()

			got, err := Resource.Create(testServer, tt.C().manifest)
			testutil.DiffError(t, tt.A().err, tt.A().errPattern, err)

			opts := []cmp.Option{
				cmp.Comparer(testutil.ComparePointer[log.Logger]),
				cmp.AllowUnexported(utilhttp.DefaultErrorHandler{}),
				cmp.AllowUnexported(channelGroup{}),
				cmp.AllowUnexported(testRunner{}, testFinalizer{}),
				cmp.AllowUnexported(sync.WaitGroup{}, atomic.Uint64{}),
				cmpopts.EquateErrors(),
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

func TestAppendInitializer(t *testing.T) {
	type condition struct {
		initializers []core.Initializer
		append       core.Initializer
	}

	type action struct {
		initializers []core.Initializer
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	cndAppendNil := tb.Condition("append nil", "append nil")
	cndAppendInitializer := tb.Condition("append initializer", "append initializer")
	actCheckSkipped := tb.Action("skipped", "check that initializer was not appended")
	actCheckAppended := tb.Action("appended", "check that initializer was appended")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"append nil",
			[]string{cndAppendNil},
			[]string{actCheckSkipped},
			&condition{
				initializers: []core.Initializer{},
				append:       nil,
			},
			&action{
				initializers: []core.Initializer{},
			},
		),
		gen(
			"append initializers",
			[]string{cndAppendInitializer},
			[]string{actCheckAppended},
			&condition{
				initializers: []core.Initializer{},
				append:       &testInitializer{name: "test"},
			},
			&action{
				initializers: []core.Initializer{
					&testInitializer{name: "test"},
				},
			},
		),
		gen(
			"append initializer",
			[]string{cndAppendInitializer},
			[]string{actCheckAppended},
			&condition{
				initializers: []core.Initializer{
					&testInitializer{name: "test1"},
				},
				append: &testInitializer{name: "test2"},
			},
			&action{
				initializers: []core.Initializer{
					&testInitializer{name: "test1"},
					&testInitializer{name: "test2"},
				},
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt

		t.Run(tt.Name(), func(t *testing.T) {
			result := appendInitializer(tt.C().initializers, tt.C().append)
			testutil.Diff(t, tt.A().initializers, result, cmp.AllowUnexported(testInitializer{}))
		})
	}
}

func TestAppendFinalizer(t *testing.T) {
	type condition struct {
		finalizers []core.Finalizer
		append     core.Finalizer
	}

	type action struct {
		finalizers []core.Finalizer
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	cndAppendNil := tb.Condition("append nil", "append nil")
	cndAppendFinalizer := tb.Condition("append finalizer", "append finalizer")
	actCheckSkipped := tb.Action("skipped", "check that finalizer was not appended")
	actCheckAppended := tb.Action("appended", "check that finalizer was appended")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"append nil",
			[]string{cndAppendNil},
			[]string{actCheckSkipped},
			&condition{
				finalizers: []core.Finalizer{},
				append:     nil,
			},
			&action{
				finalizers: []core.Finalizer{},
			},
		),
		gen(
			"append finalizer",
			[]string{cndAppendFinalizer},
			[]string{actCheckAppended},
			&condition{
				finalizers: []core.Finalizer{},
				append:     &testFinalizer{name: "test"},
			},
			&action{
				finalizers: []core.Finalizer{
					&testFinalizer{name: "test"},
				},
			},
		),
		gen(
			"append finalizer",
			[]string{cndAppendFinalizer},
			[]string{actCheckAppended},
			&condition{
				finalizers: []core.Finalizer{
					&testFinalizer{name: "test1"},
				},
				append: &testFinalizer{name: "test2"},
			},
			&action{
				finalizers: []core.Finalizer{
					&testFinalizer{name: "test1"},
					&testFinalizer{name: "test2"},
				},
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt

		t.Run(tt.Name(), func(t *testing.T) {
			result := appendFinalizer(tt.C().finalizers, tt.C().append)
			testutil.Diff(t, tt.A().finalizers, result, cmp.AllowUnexported(testFinalizer{}))
		})
	}
}
