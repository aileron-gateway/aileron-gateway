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

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"new instance",
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

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			a := NewFactoryAPI()
			testutil.Diff(t, tt.A.a, a, cmp.AllowUnexported(FactoryAPI{}))
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

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"register 1 resource",
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

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			a := NewFactoryAPI()

			var err error
			for i := range tt.C.keys {
				err = a.Register(tt.C.keys[i], tt.C.resources[i])
			}
			testutil.Diff(t, tt.A.err, err, cmpopts.EquateErrors())
			testutil.Diff(t, tt.A.resources, a.resources)
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

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"Delete nothing",
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

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			a := tt.C.a
			err := a.delete(context.Background(), tt.C.req, tt.C.resource)

			testutil.Diff(t, tt.A.err, err, cmpopts.EquateErrors())
			testutil.Diff(t, tt.A.objStore, a.objStore)
			testutil.Diff(t, tt.A.protoStore, a.protoStore)
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

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"Post manifest",
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

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			a := tt.C.a
			err := a.post(context.Background(), tt.C.req, tt.C.resource)

			testutil.Diff(t, tt.A.err, err, cmpopts.EquateErrors())
			testutil.Diff(t, tt.A.objStore, a.objStore)
			testutil.Diff(t, tt.A.protoStore, a.protoStore, cmpopts.IgnoreUnexported(k.Resource{}, k.Metadata{}))
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

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"Get new instance",
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

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			a := tt.C.a
			obj, err := a.get(context.Background(), tt.C.req, tt.C.resource)

			testutil.Diff(t, tt.A.err, err, cmpopts.EquateErrors())
			testutil.Diff(t, tt.A.objStore, a.objStore)
			testutil.Diff(t, tt.A.protoStore, a.protoStore, cmpopts.IgnoreUnexported(k.Resource{}, k.Metadata{}))

			// Because JSON does not ensure the order of objects init,
			// except the check for JSON returned by the get method.
			// Note that the protoreflect package intentionally use single space " " and double space "  " randomly
			// when marshalling ProtoMessage to JSON.
			if tt.C.req.Params != nil && tt.C.req.Params[KeyAccept] != string(FormatJSON) {
				testutil.Diff(t, tt.A.obj, obj, cmpopts.IgnoreUnexported(k.Resource{}, k.Metadata{}))
			}
		})
	}
}
