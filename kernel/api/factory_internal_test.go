// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package api

import (
	"context"
	"testing"

	k "github.com/aileron-gateway/aileron-gateway/apis/kernel"
	"github.com/aileron-gateway/aileron-gateway/kernel/encoder"
	"github.com/aileron-gateway/aileron-gateway/kernel/er"
	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func TestNewFactoryAPI(t *testing.T) {
	type condition struct {
	}

	type action struct {
		a *FactoryAPI
	}

	cndNewDefault := "new default"
	actCheckInitialized := "check initialized "

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	tb.Condition(cndNewDefault, "create a new instance")
	tb.Action(actCheckInitialized, "check that the returned instance is initialized with expected values")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"new instance",
			[]string{cndNewDefault},
			[]string{actCheckInitialized},
			&condition{},
			&action{
				a: &FactoryAPI{
					protoStore: map[string]protoreflect.ProtoMessage{},
					objStore:   map[string]any{},
					resources:  map[string]Resource{},
				},
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			a := NewFactoryAPI()
			testutil.Diff(t, tt.A().a, a, cmp.AllowUnexported(FactoryAPI{}))
		})
	}
}

type noopResource struct {
	*BaseResource
	ID string // To make this struct comparative in the test.
}

func (r *noopResource) Default() protoreflect.ProtoMessage {
	return nil
}

func (r *noopResource) Create(a API[*Request, *Response], msg protoreflect.ProtoMessage) (any, error) {
	return nil, nil
}

func TestFactoryAPI_Register(t *testing.T) {
	type condition struct {
		keys      []string
		resources []Resource
	}

	type action struct {
		resources map[string]Resource
		err       error
	}

	cndRegisterOne := "1 resource"
	cndRegisterMultiple := "multiple resources"
	cndRegisterNil := "register nil"
	cndRegisterDuplicateKey := "duplicate key"
	actCheckRegistered := "check registered resources"
	actCheckNoError := "no error"
	actCheckError := "non-nil error"

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	tb.Condition(cndRegisterOne, "register 1 non-nil resource")
	tb.Condition(cndRegisterMultiple, "register multiple non-nil resource with different keys")
	tb.Condition(cndRegisterNil, "try to register nil resource")
	tb.Condition(cndRegisterDuplicateKey, "try to register resources with the same key")
	tb.Action(actCheckRegistered, "check that the registered resources are the same as expected")
	tb.Action(actCheckNoError, "check that there is no error")
	tb.Action(actCheckError, "check that a non-nil error was returned")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"register 1 resource",
			[]string{cndRegisterOne},
			[]string{actCheckRegistered, actCheckNoError},
			&condition{
				keys:      []string{"test"},
				resources: []Resource{&noopResource{ID: "foo"}},
			},
			&action{
				resources: map[string]Resource{
					"test": &noopResource{ID: "foo"},
				},
			},
		),
		gen(
			"register multiple resources",
			[]string{cndRegisterMultiple},
			[]string{actCheckRegistered, actCheckNoError},
			&condition{
				keys:      []string{"test1", "test2"},
				resources: []Resource{&noopResource{ID: "foo"}, &noopResource{ID: "bar"}},
			},
			&action{
				resources: map[string]Resource{
					"test1": &noopResource{ID: "foo"},
					"test2": &noopResource{ID: "bar"},
				},
			},
		),
		gen(
			"register nil",
			[]string{cndRegisterNil},
			[]string{actCheckRegistered, actCheckNoError},
			&condition{
				keys:      []string{"test"},
				resources: []Resource{nil},
			},
			&action{
				resources: map[string]Resource{},
			},
		),
		gen(
			"duplicate key",
			[]string{cndRegisterDuplicateKey},
			[]string{actCheckRegistered, actCheckError},
			&condition{
				keys:      []string{"test", "test"},
				resources: []Resource{&noopResource{ID: "foo"}, &noopResource{ID: "bar"}},
			},
			&action{
				resources: map[string]Resource{
					"test": &noopResource{ID: "foo"},
				},
				err: &er.Error{
					Package:     ErrPkg,
					Type:        ErrTypeFactory,
					Description: ErrDscDuplicateKey,
				},
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			a := NewFactoryAPI()

			var err error
			for i := range tt.C().keys {
				err = a.Register(tt.C().keys[i], tt.C().resources[i])
			}
			testutil.Diff(t, tt.A().err, err, cmpopts.EquateErrors())
			testutil.Diff(t, tt.A().resources, a.resources)
		})
	}
}

type testResource struct {
	err error
}

func (r *testResource) Default() protoreflect.ProtoMessage {
	return &k.Resource{}
}

func (r *testResource) Create(a API[*Request, *Response], msg protoreflect.ProtoMessage) (any, error) {
	if msg == nil {
		return nil, r.err
	}
	c := msg.(*k.Resource)
	return c.Metadata.Namespace + " " + c.Metadata.Name, r.err
}

func (r *testResource) Validate(msg protoreflect.ProtoMessage) error {
	return r.err
}

func (r *testResource) Mutate(msg protoreflect.ProtoMessage) protoreflect.ProtoMessage {
	return msg
}

func (r *testResource) Delete(a API[*Request, *Response], msg protoreflect.ProtoMessage, obj any) error {
	return r.err
}

func TestFactoryAPI_delete(t *testing.T) {
	type condition struct {
		a        *FactoryAPI
		resource Resource
		req      *Request
	}

	type action struct {
		protoStore map[string]protoreflect.ProtoMessage
		objStore   map[string]any
		err        error
	}

	cndErrorDelete := "delete error"
	cndWrongType := "wrong content"
	cndUnsupportedFormat := "unsupported format"
	actCheckNoError := "no error"
	actCheckError := "non-nil error"

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	tb.Condition(cndErrorDelete, "error occurred in delete method")
	tb.Condition(cndWrongType, "content in the request is invalid")
	tb.Condition(cndUnsupportedFormat, "specify unsupported format")
	tb.Action(actCheckNoError, "check that there is no error")
	tb.Action(actCheckError, "check that a non-nil error was returned")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"Delete nothing",
			[]string{},
			[]string{actCheckNoError},
			&condition{
				a:        NewFactoryAPI(),
				resource: &testResource{},
				req: &Request{
					Method:  MethodDelete,
					Key:     "test1/test2",
					Format:  FormatJSON,
					Content: []byte(`{"apiVersion":"test1", "kind":"test2", "metadata": {"namespace":"test3", "name":"test4"}}`),
				},
			},
			&action{
				protoStore: map[string]protoreflect.ProtoMessage{},
				objStore:   map[string]any{},
			},
		),
		gen(
			"Delete object",
			[]string{},
			[]string{actCheckNoError},
			&condition{
				a: &FactoryAPI{
					protoStore: map[string]protoreflect.ProtoMessage{"test1/test2/test3/test4": &k.Reference{}},
					objStore:   map[string]any{"test1/test2/test3/test4": "test3 test4"},
				},
				resource: &testResource{},
				req: &Request{
					Method:  MethodDelete,
					Key:     "test1/test2",
					Format:  FormatJSON,
					Content: []byte(`{"apiVersion":"test1", "kind":"test2", "metadata": {"namespace":"test3", "name":"test4"}}`),
				},
			},
			&action{
				protoStore: map[string]protoreflect.ProtoMessage{},
				objStore:   map[string]any{},
			},
		),
		gen(
			"Delete fails",
			[]string{cndErrorDelete},
			[]string{actCheckError},
			&condition{
				a:        NewFactoryAPI(),
				resource: &testResource{err: &er.Error{}}, // Use APIError for dummy.
				req: &Request{
					Method:  MethodDelete,
					Key:     "test1/test2",
					Format:  FormatJSON,
					Content: []byte(`{"apiVersion":"test1", "kind":"test2", "metadata": {"namespace":"test3", "name":"test4"}}`),
				},
			},
			&action{
				protoStore: map[string]protoreflect.ProtoMessage{},
				objStore:   map[string]any{},
				err:        &er.Error{},
			},
		),
		gen(
			"Invalid content type",
			[]string{cndWrongType},
			[]string{actCheckError},
			&condition{
				a:        NewFactoryAPI(),
				resource: &testResource{},
				req: &Request{
					Method:  MethodDelete,
					Key:     "test1/test2",
					Format:  FormatJSON,
					Content: 999, // FormatJSON require []byte JSON but set int to make an error.
				},
			},
			&action{
				protoStore: map[string]protoreflect.ProtoMessage{},
				objStore:   map[string]any{},
				err: &er.Error{
					Package:     ErrPkg,
					Type:        ErrTypeUtil,
					Description: ErrDscAssert,
				},
			},
		),
		gen(
			"Invalid format",
			[]string{cndUnsupportedFormat},
			[]string{actCheckError},
			&condition{
				a:        NewFactoryAPI(),
				resource: &testResource{},
				req: &Request{
					Method: MethodDelete,
					Key:    "test1/test2",
					Format: Format("UNSUPPORTED"),
				},
			},
			&action{
				protoStore: map[string]protoreflect.ProtoMessage{},
				objStore:   map[string]any{},
				err: &er.Error{
					Package:     encoder.ErrPkg,
					Type:        encoder.ErrTypeJSON,
					Description: encoder.ErrDscUnmarshal,
				},
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			a := tt.C().a
			err := a.delete(context.Background(), tt.C().req, tt.C().resource)

			testutil.Diff(t, tt.A().err, err, cmpopts.EquateErrors())
			testutil.Diff(t, tt.A().objStore, a.objStore)
			testutil.Diff(t, tt.A().protoStore, a.protoStore)
		})
	}
}

func TestFactoryAPI_post(t *testing.T) {
	type condition struct {
		a        *FactoryAPI
		resource Resource
		req      *Request
	}

	type action struct {
		protoStore map[string]protoreflect.ProtoMessage
		objStore   map[string]any
		err        error
	}

	cndErrorPost := "post error"
	cndDuplicateKey := "duplicate key"
	cndWrongType := "wrong content"
	cndUnsupportedFormat := "unsupported format"
	actCheckNoError := "no error"
	actCheckError := "non-nil error"

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	tb.Condition(cndErrorPost, "error occurred in post method")
	tb.Condition(cndDuplicateKey, "post manifest with the same key")
	tb.Condition(cndWrongType, "content in the request is invalid")
	tb.Condition(cndUnsupportedFormat, "specify unsupported format")
	tb.Action(actCheckNoError, "check that there is no error")
	tb.Action(actCheckError, "check that a non-nil error was returned")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"Post manifest",
			[]string{},
			[]string{actCheckNoError},
			&condition{
				a:        NewFactoryAPI(),
				resource: &testResource{},
				req: &Request{
					Method:  MethodPost,
					Key:     "test1/test2",
					Format:  FormatJSON,
					Content: []byte(`{"apiVersion":"test1", "kind":"test2", "metadata": {"namespace":"test3", "name":"test4"}}`),
				},
			},
			&action{
				protoStore: map[string]protoreflect.ProtoMessage{
					"test1/test2/test3/test4": &k.Resource{
						APIVersion: "test1",
						Kind:       "test2",
						Metadata: &k.Metadata{
							Namespace: "test3",
							Name:      "test4",
						},
					},
				},
				objStore: map[string]any{},
			},
		),
		gen(
			"duplicate key",
			[]string{cndDuplicateKey},
			[]string{actCheckNoError},
			&condition{
				a: &FactoryAPI{
					protoStore: map[string]protoreflect.ProtoMessage{"test1/test2/test3/test4": nil},
				},
				resource: &testResource{},
				req: &Request{
					Method:  MethodPost,
					Key:     "test1/test2",
					Format:  FormatJSON,
					Content: []byte(`{"apiVersion":"test1", "kind":"test2", "metadata": {"namespace":"test3", "name":"test4"}}`),
				},
			},
			&action{
				protoStore: map[string]protoreflect.ProtoMessage{"test1/test2/test3/test4": nil},
				err: &er.Error{
					Package:     ErrPkg,
					Type:        ErrTypeFactory,
					Description: ErrDscDuplicateKey,
				},
			},
		),
		gen(
			"Post fails",
			[]string{cndErrorPost},
			[]string{actCheckError},
			&condition{
				a:        NewFactoryAPI(),
				resource: &testResource{err: &er.Error{}}, // Use APIError for dummy.
				req: &Request{
					Method:  MethodPost,
					Key:     "test1/test2",
					Format:  FormatJSON,
					Content: []byte(`{"apiVersion":"test1", "kind":"test2", "metadata": {"namespace":"test3", "name":"test4"}}`),
				},
			},
			&action{
				protoStore: map[string]protoreflect.ProtoMessage{},
				objStore:   map[string]any{},
				err:        &er.Error{},
			},
		),
		gen(
			"Invalid content type",
			[]string{cndWrongType},
			[]string{actCheckError},
			&condition{
				a:        NewFactoryAPI(),
				resource: &testResource{},
				req: &Request{
					Method:  MethodPost,
					Key:     "test1/test2",
					Format:  FormatJSON,
					Content: 999, // FormatJSON require []byte JSON but set int to make an error.
				},
			},
			&action{
				protoStore: map[string]protoreflect.ProtoMessage{},
				objStore:   map[string]any{},
				err: &er.Error{
					Package:     ErrPkg,
					Type:        ErrTypeUtil,
					Description: ErrDscAssert,
				},
			},
		),
		gen(
			"Unsupported format",
			[]string{cndUnsupportedFormat},
			[]string{actCheckError},
			&condition{
				a:        NewFactoryAPI(),
				resource: &testResource{},
				req: &Request{
					Method: MethodPost,
					Key:    "test1/test2",
					Format: Format("UNSUPPORTED"),
				},
			},
			&action{
				protoStore: map[string]protoreflect.ProtoMessage{},
				objStore:   map[string]any{},
				err: &er.Error{
					Package:     encoder.ErrPkg,
					Type:        encoder.ErrTypeJSON,
					Description: encoder.ErrDscUnmarshal,
				},
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			a := tt.C().a
			err := a.post(context.Background(), tt.C().req, tt.C().resource)

			testutil.Diff(t, tt.A().err, err, cmpopts.EquateErrors())
			testutil.Diff(t, tt.A().objStore, a.objStore)
			testutil.Diff(t, tt.A().protoStore, a.protoStore, cmpopts.IgnoreUnexported(k.Resource{}, k.Metadata{}))
		})
	}
}

func TestFactoryAPI_get(t *testing.T) {
	type condition struct {
		a        *FactoryAPI
		resource Resource
		req      *Request
	}

	type action struct {
		obj        any
		protoStore map[string]protoreflect.ProtoMessage
		objStore   map[string]any
		err        error
	}

	cndAcceptJSON := "accept JSON"
	cndAcceptYAML := "accept YAML"
	cndAcceptProtoMessage := "accept ProtoMessage"
	cndDefault := "use default"
	cndErrorGet := "get error"
	cndWrongType := "wrong content"
	cndUnsupportedFormat := "unsupported format"
	actCheckObject := "check returned object"
	actCheckNoError := "no error"
	actCheckError := "non-nil error"

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	tb.Condition(cndAcceptJSON, "specify accept parameter to get JSON response")
	tb.Condition(cndAcceptYAML, "specify accept parameter to get YAML response")
	tb.Condition(cndAcceptProtoMessage, "specify accept parameter to get ProtoMessage response")
	tb.Condition(cndDefault, "specify namespace and name to use default ProtoMessage")
	tb.Condition(cndErrorGet, "error occurred in get method")
	tb.Condition(cndWrongType, "content in the request is invalid")
	tb.Condition(cndUnsupportedFormat, "specify unsupported format")
	tb.Action(actCheckObject, "check that the returned object is the same as expected")
	tb.Action(actCheckNoError, "check that there is no error")
	tb.Action(actCheckError, "check that a non-nil error was returned")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"Get new instance",
			[]string{},
			[]string{actCheckObject, actCheckNoError},
			&condition{
				a: &FactoryAPI{
					protoStore: map[string]protoreflect.ProtoMessage{
						"test1/test2/test3/test4": &k.Resource{
							APIVersion: "test1",
							Kind:       "test2",
							Metadata: &k.Metadata{
								Namespace: "test3",
								Name:      "test4",
							},
						},
					},
					objStore: map[string]any{},
				},
				resource: &testResource{},
				req: &Request{
					Method:  MethodGet,
					Key:     "test1/test2",
					Format:  FormatJSON,
					Content: []byte(`{"apiVersion":"test1", "kind":"test2", "metadata": {"namespace":"test3", "name":"test4"}}`),
				},
			},
			&action{
				obj: "test3 test4",
				protoStore: map[string]protoreflect.ProtoMessage{
					"test1/test2/test3/test4": &k.Resource{
						APIVersion: "test1",
						Kind:       "test2",
						Metadata: &k.Metadata{
							Namespace: "test3",
							Name:      "test4",
						},
					},
				},
				objStore: map[string]any{"test1/test2/test3/test4": "test3 test4"},
			},
		),
		gen(
			"Get existing instance",
			[]string{},
			[]string{actCheckObject, actCheckNoError},
			&condition{
				a: &FactoryAPI{
					protoStore: map[string]protoreflect.ProtoMessage{
						"test1/test2/test3/test4": &k.Resource{
							APIVersion: "test1",
							Kind:       "test2",
							Metadata: &k.Metadata{
								Namespace: "test3",
								Name:      "test4",
							},
						},
					},
					objStore: map[string]any{"test1/test2/test3/test4": "test3 test4"},
				},
				resource: &testResource{},
				req: &Request{
					Method:  MethodGet,
					Key:     "test1/test2",
					Format:  FormatJSON,
					Content: []byte(`{"apiVersion":"test1", "kind":"test2", "metadata": {"namespace":"test3", "name":"test4"}}`),
				},
			},
			&action{
				obj: "test3 test4",
				protoStore: map[string]protoreflect.ProtoMessage{
					"test1/test2/test3/test4": &k.Resource{
						APIVersion: "test1",
						Kind:       "test2",
						Metadata: &k.Metadata{
							Namespace: "test3",
							Name:      "test4",
						},
					},
				},
				objStore: map[string]any{"test1/test2/test3/test4": "test3 test4"},
			},
		),
		gen(
			"accept as json",
			[]string{cndAcceptJSON},
			[]string{actCheckNoError},
			&condition{
				a: &FactoryAPI{
					protoStore: map[string]protoreflect.ProtoMessage{
						"test1/test2/test3/test4": &k.Resource{
							APIVersion: "test1",
							Kind:       "test2",
							Metadata: &k.Metadata{
								Namespace: "test3",
								Name:      "test4",
							},
						},
					},
				},
				resource: &testResource{},
				req: &Request{
					Method:  MethodGet,
					Key:     "test1/test2",
					Params:  map[string]string{KeyAccept: string(FormatJSON)},
					Format:  FormatJSON,
					Content: []byte(`{"apiVersion":"test1", "kind":"test2", "metadata": {"namespace":"test3", "name":"test4"}}`),
				},
			},
			&action{
				obj: "<<< Do not check this value because single space and double spaces are randomly used for marshalling ProtoMessage to JSON >>>",
				protoStore: map[string]protoreflect.ProtoMessage{
					"test1/test2/test3/test4": &k.Resource{
						APIVersion: "test1",
						Kind:       "test2",
						Metadata: &k.Metadata{
							Namespace: "test3",
							Name:      "test4",
						},
					},
				},
			},
		),
		gen(
			"accept as yaml",
			[]string{cndAcceptYAML},
			[]string{actCheckObject, actCheckNoError},
			&condition{
				a: &FactoryAPI{
					protoStore: map[string]protoreflect.ProtoMessage{
						"test1/test2/test3/test4": &k.Resource{
							APIVersion: "test1",
							Kind:       "test2",
							Metadata: &k.Metadata{
								Namespace: "test3",
								Name:      "test4",
							},
						},
					},
				},
				resource: &testResource{},
				req: &Request{
					Method:  MethodGet,
					Key:     "test1/test2",
					Params:  map[string]string{KeyAccept: string(FormatYAML)},
					Format:  FormatJSON,
					Content: []byte(`{"apiVersion":"test1", "kind":"test2", "metadata": {"namespace":"test3", "name":"test4"}}`),
				},
			},
			&action{
				obj: []byte("apiVersion: test1\nkind: test2\nmetadata:\n  errorHandler: \"\"\n  logger: \"\"\n  name: test4\n  namespace: test3\nspec: null\n"),
				protoStore: map[string]protoreflect.ProtoMessage{
					"test1/test2/test3/test4": &k.Resource{
						APIVersion: "test1",
						Kind:       "test2",
						Metadata: &k.Metadata{
							Namespace: "test3",
							Name:      "test4",
						},
						Spec: nil,
					},
				},
			},
		),
		gen(
			"accept as proto message",
			[]string{cndAcceptProtoMessage},
			[]string{actCheckObject, actCheckNoError},
			&condition{
				a: &FactoryAPI{
					protoStore: map[string]protoreflect.ProtoMessage{
						"test1/test2/test3/test4": &k.Resource{
							APIVersion: "test1",
							Kind:       "test2",
							Metadata: &k.Metadata{
								Namespace: "test3",
								Name:      "test4",
							},
						},
					},
				},
				resource: &testResource{},
				req: &Request{
					Method:  MethodGet,
					Key:     "test1/test2",
					Params:  map[string]string{KeyAccept: string(FormatProtoMessage)},
					Format:  FormatJSON,
					Content: []byte(`{"apiVersion":"test1", "kind":"test2", "metadata": {"namespace":"test3", "name":"test4"}}`),
				},
			},
			&action{
				obj: &k.Resource{
					APIVersion: "test1",
					Kind:       "test2",
					Metadata: &k.Metadata{
						Namespace: "test3",
						Name:      "test4",
					},
				},
				protoStore: map[string]protoreflect.ProtoMessage{
					"test1/test2/test3/test4": &k.Resource{
						APIVersion: "test1",
						Kind:       "test2",
						Metadata: &k.Metadata{
							Namespace: "test3",
							Name:      "test4",
						},
					},
				},
			},
		),
		gen(
			"use template",
			[]string{cndDefault, cndAcceptProtoMessage},
			[]string{actCheckObject, actCheckNoError},
			&condition{
				a:        &FactoryAPI{},
				resource: &testResource{},
				req: &Request{
					Method:  MethodGet,
					Key:     "test1/test2",
					Params:  map[string]string{KeyAccept: string(FormatProtoMessage)},
					Format:  FormatJSON,
					Content: []byte(`{"apiVersion":"test1", "kind":"test2", "metadata": {"namespace":"template", "name":"template"}}`),
				},
			},
			&action{
				obj: &k.Resource{},
			},
		),
		gen(
			"Get fails",
			[]string{cndErrorGet},
			[]string{actCheckObject, actCheckError},
			&condition{
				a: &FactoryAPI{
					protoStore: map[string]protoreflect.ProtoMessage{"test1/test2/test3/test4": nil},
				},
				resource: &testResource{err: &er.Error{}}, // Use APIError for dummy.
				req: &Request{
					Method:  MethodGet,
					Key:     "test1/test2",
					Format:  FormatJSON,
					Content: []byte(`{"apiVersion":"test1", "kind":"test2", "metadata": {"namespace":"test3", "name":"test4"}}`),
				},
			},
			&action{
				objStore:   nil,
				protoStore: map[string]protoreflect.ProtoMessage{"test1/test2/test3/test4": nil},
				err:        &er.Error{},
			},
		),
		gen(
			"Invalid content type",
			[]string{cndWrongType},
			[]string{actCheckObject, actCheckError},
			&condition{
				a:        NewFactoryAPI(),
				resource: &testResource{},
				req: &Request{
					Method:  MethodGet,
					Key:     "test1/test2",
					Format:  FormatJSON,
					Content: 999, // FormatJSON require []byte JSON but set int to make an error.
				},
			},
			&action{
				objStore:   map[string]any{},
				protoStore: map[string]protoreflect.ProtoMessage{},
				err: &er.Error{
					Package:     ErrPkg,
					Type:        ErrTypeUtil,
					Description: ErrDscAssert,
				},
			},
		),
		gen(
			"Unsupported format",
			[]string{cndUnsupportedFormat},
			[]string{actCheckObject, actCheckError},
			&condition{
				a:        NewFactoryAPI(),
				resource: &testResource{},
				req: &Request{
					Method: MethodGet,
					Key:    "test1/test2",
					Format: Format("UNSUPPORTED"),
				},
			},
			&action{
				objStore:   map[string]any{},
				protoStore: map[string]protoreflect.ProtoMessage{},
				err: &er.Error{
					Package:     encoder.ErrPkg,
					Type:        encoder.ErrTypeJSON,
					Description: encoder.ErrDscUnmarshal,
				},
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			a := tt.C().a
			obj, err := a.get(context.Background(), tt.C().req, tt.C().resource)

			testutil.Diff(t, tt.A().err, err, cmpopts.EquateErrors())
			testutil.Diff(t, tt.A().objStore, a.objStore)
			testutil.Diff(t, tt.A().protoStore, a.protoStore, cmpopts.IgnoreUnexported(k.Resource{}, k.Metadata{}))

			// Because JSON does not ensure the order of objects init,
			// except the check for JSON returned by the get method.
			// Note that the protoreflect package intentionally use single space " " and double space "  " randomly
			// when marshalling ProtoMessage to JSON.
			if tt.C().req.Params != nil && tt.C().req.Params[KeyAccept] != string(FormatJSON) {
				testutil.Diff(t, tt.A().obj, obj, cmpopts.IgnoreUnexported(k.Resource{}, k.Metadata{}))
			}
		})
	}
}
