// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package api_test

import (
	"context"
	"io"
	"testing"

	k "github.com/aileron-gateway/aileron-gateway/apis/kernel"
	"github.com/aileron-gateway/aileron-gateway/internal/encoder"
	"github.com/aileron-gateway/aileron-gateway/internal/testutil"
	"github.com/aileron-gateway/aileron-gateway/kernel/api"
	"github.com/aileron-gateway/aileron-gateway/kernel/er"
	"github.com/google/go-cmp/cmp/cmpopts"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func TestFormat_Unmarshal(t *testing.T) {
	type condition struct {
		format api.Format
		in     any
		into   any
	}

	type action struct {
		result any
		err    error
	}

	type testStruct struct {
		Foo string `json:"foo" yaml:"foo"`
		Bar int    `json:"bar" yaml:"bar"`
	}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"json", &condition{
				format: api.FormatJSON,
				in:     []byte(`{"foo":"test", "bar":999}`),
				into:   &testStruct{},
			},
			&action{
				result: &testStruct{
					Foo: "test",
					Bar: 999,
				},
			},
		),
		gen(
			"yaml", &condition{
				format: api.FormatYAML,
				in:     []byte("foo: test \nbar: 999"),
				into:   &testStruct{},
			},
			&action{
				result: &testStruct{
					Foo: "test",
					Bar: 999,
				},
			},
		),
		gen(
			"unsupported", &condition{
				format: api.Format("UNSUPPORTED"),
				in:     []byte(`{"foo":"test", "bar":999}`),
				into:   &testStruct{},
			},
			&action{
				result: &testStruct{},
				err: &er.Error{
					Package:     api.ErrPkg,
					Type:        api.ErrTypeUtil,
					Description: api.ErrDscFormatSupport,
				},
			},
		),
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			err := tt.C.format.Unmarshal(tt.C.in, tt.C.into)
			testutil.Diff(t, tt.A.err, err, cmpopts.EquateErrors())
			testutil.Diff(t, tt.A.result, tt.C.into)
		})
	}
}

type stringAPI string

func (s stringAPI) Serve(_ context.Context, req *api.Request) (*api.Response, error) {
	return &api.Response{Content: s}, nil
}

func TestContextWithRoute(t *testing.T) {
	type condition struct {
		ctx context.Context
		a   api.API[*api.Request, *api.Response]
	}

	type action struct {
		apis []api.API[*api.Request, *api.Response]
	}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"new", &condition{
				ctx: context.Background(),
				a:   stringAPI("test"),
			},
			&action{
				apis: []api.API[*api.Request, *api.Response]{
					stringAPI("test"),
				},
			},
		),
		gen(
			"append", &condition{
				ctx: context.WithValue(context.Background(), api.APIRouteContextKey, []api.API[*api.Request, *api.Response]{stringAPI("test1")}),
				a:   stringAPI("test2"),
			},
			&action{
				apis: []api.API[*api.Request, *api.Response]{
					stringAPI("test1"),
					stringAPI("test2"),
				},
			},
		),
		gen(
			"append nil", &condition{
				ctx: context.WithValue(context.Background(), api.APIRouteContextKey, []api.API[*api.Request, *api.Response]{stringAPI("test")}),
				a:   nil,
			},
			&action{
				apis: []api.API[*api.Request, *api.Response]{
					stringAPI("test"),
				},
			},
		),
		gen(
			"nil context", &condition{
				ctx: nil,
				a:   stringAPI("test"),
			},
			&action{
				apis: []api.API[*api.Request, *api.Response]{
					stringAPI("test"),
				},
			},
		),
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			ctx := api.ContextWithRoute(tt.C.ctx, tt.C.a)
			testutil.Diff(t, tt.A.apis, ctx.Value(api.APIRouteContextKey))
		})
	}
}

func TestRootAPIFromContext(t *testing.T) {
	type condition struct {
		ctx context.Context
	}

	type action struct {
		a api.API[*api.Request, *api.Response]
	}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"api not registered", &condition{
				ctx: context.WithValue(context.Background(), api.APIRouteContextKey,
					[]api.API[*api.Request, *api.Response]{}),
			},
			&action{
				a: nil,
			},
		),
		gen(
			"1 api registered", &condition{
				ctx: context.WithValue(context.Background(), api.APIRouteContextKey,
					[]api.API[*api.Request, *api.Response]{stringAPI("test")}),
			},
			&action{
				a: stringAPI("test"),
			},
		),
		gen(
			"multiple apis registered", &condition{
				ctx: context.WithValue(context.Background(), api.APIRouteContextKey,
					[]api.API[*api.Request, *api.Response]{stringAPI("test1"), stringAPI("test2")}),
			},
			&action{
				a: stringAPI("test1"),
			},
		),
		gen(
			"nil context", &condition{
				ctx: nil,
			},
			&action{
				a: nil,
			},
		),
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			a := api.RootAPIFromContext(tt.C.ctx)
			testutil.Diff(t, tt.A.a, a)
		})
	}
}

func TestProtoMessage(t *testing.T) {
	type condition struct {
		format  api.Format
		content any
		msg     protoreflect.ProtoMessage
		opt     *protojson.UnmarshalOptions
	}

	type action struct {
		msg protoreflect.ProtoMessage
		err error
	}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"json/success", &condition{
				format:  api.FormatJSON,
				content: []byte(`{"apiVersion":"test1", "kind":"test2", "namespace":"test3", "name":"test4"}`),
				msg:     &k.Reference{},
				opt:     nil,
			},
			&action{
				msg: &k.Reference{
					APIVersion: "test1",
					Kind:       "test2",
					Namespace:  "test3",
					Name:       "test4",
				},
			},
		),
		gen(
			"json/merge", &condition{
				format:  api.FormatJSON,
				content: []byte(`{"apiVersion":"test1", "kind":"test2", "namespace":"test3"}`),
				msg: &k.Reference{
					Namespace: "default",
					Name:      "default",
				},
				opt: nil,
			},
			&action{
				msg: &k.Reference{
					APIVersion: "test1",
					Kind:       "test2",
					Namespace:  "test3",   // Value in the content is prior to the default.
					Name:       "default", // Default value.
				},
			},
		),
		gen(
			"json/error", &condition{
				format:  api.FormatJSON,
				content: nil,
				msg:     nil,
				opt:     nil,
			},
			&action{
				msg: nil,
				err: &er.Error{
					Package:     api.ErrPkg,
					Type:        api.ErrTypeUtil,
					Description: api.ErrDscAssert,
				},
			},
		),
		gen(
			"yaml/success", &condition{
				format:  api.FormatYAML,
				content: []byte("apiVersion: test1 \nkind: test2 \nnamespace: test3 \nname: test4"),
				msg:     &k.Reference{},
				opt:     nil,
			},
			&action{
				msg: &k.Reference{
					APIVersion: "test1",
					Kind:       "test2",
					Namespace:  "test3",
					Name:       "test4",
				},
			},
		),
		gen(
			"yaml/merge", &condition{
				format:  api.FormatYAML,
				content: []byte("apiVersion: test1 \nkind: test2 \nnamespace: test3"),
				msg: &k.Reference{
					Namespace: "default",
					Name:      "default",
				},
				opt: nil,
			},
			&action{
				msg: &k.Reference{
					APIVersion: "test1",
					Kind:       "test2",
					Namespace:  "test3",   // Value in the content is prior to the default.
					Name:       "default", // Default value.
				},
			},
		),
		gen(
			"yaml/error", &condition{
				format:  api.FormatYAML,
				content: nil,
				msg:     nil,
				opt:     nil,
			},
			&action{
				msg: nil,
				err: &er.Error{
					Package:     api.ErrPkg,
					Type:        api.ErrTypeUtil,
					Description: api.ErrDscAssert,
				},
			},
		),
		gen(
			"protoMessage/success", &condition{
				format: api.FormatProtoMessage,
				content: &k.Reference{
					APIVersion: "test1",
					Kind:       "test2",
					Namespace:  "test3",
					Name:       "test4",
				},
				msg: &k.Reference{},
				opt: nil,
			},
			&action{
				msg: &k.Reference{
					APIVersion: "test1",
					Kind:       "test2",
					Namespace:  "test3",
					Name:       "test4",
				},
			},
		),
		gen(
			"protoMessage/merge", &condition{
				format: api.FormatProtoMessage,
				content: &k.Reference{
					APIVersion: "test1",
					Kind:       "test2",
					Namespace:  "test3",
				},
				msg: &k.Reference{
					Namespace: "default",
					Name:      "default",
				},
				opt: nil,
			},
			&action{
				msg: &k.Reference{
					APIVersion: "test1",
					Kind:       "test2",
					Namespace:  "test3",   // Value in the content is prior to the default.
					Name:       "default", // Default value.
				},
			},
		),
		gen(
			"protoMessage/error", &condition{
				format:  api.FormatProtoMessage,
				content: nil,
				msg:     nil,
				opt:     nil,
			},
			&action{
				msg: nil,
				err: &er.Error{
					Package:     api.ErrPkg,
					Type:        api.ErrTypeUtil,
					Description: api.ErrDscAssert,
				},
			},
		),
		gen(
			"protoReference/success", &condition{
				format: api.FormatProtoReference,
				content: &k.Reference{
					APIVersion: "test1",
					Kind:       "test2",
					Namespace:  "test3",
					Name:       "test4",
				},
				msg: &k.Reference{},
				opt: nil,
			},
			&action{
				msg: &k.Reference{
					APIVersion: "test1",
					Kind:       "test2",
					Namespace:  "test3",
					Name:       "test4",
				},
			},
		),
		gen(
			"protoReference/merge", &condition{
				format: api.FormatProtoReference,
				content: &k.Reference{
					APIVersion: "test1",
					Kind:       "test2",
					Namespace:  "test3",
				},
				msg: &k.Reference{
					Namespace: "default",
					Name:      "default",
				},
				opt: nil,
			},
			&action{
				msg: &k.Reference{
					APIVersion: "test1",
					Kind:       "test2",
					Namespace:  "test3",
					Name:       "", //Default message is not used for api.FormatProtoReference type.
				},
			},
		),
		gen(
			"protoReference/error", &condition{
				format:  api.FormatProtoReference,
				content: nil,
				msg:     nil,
				opt:     nil,
			},
			&action{
				msg: nil,
				err: &er.Error{
					Package:     api.ErrPkg,
					Type:        api.ErrTypeUtil,
					Description: api.ErrDscAssert,
				},
			},
		),
		gen(
			"unsupported format", &condition{
				format:  api.Format("unsupported"),
				content: nil,
				msg:     nil,
				opt:     nil,
			},
			&action{
				msg: nil, // nil for unsupported format.
				err: nil, // nil for unsupported format.
			},
		),
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			msg, err := api.ProtoMessage(tt.C.format, tt.C.content, tt.C.msg, tt.C.opt)

			testutil.Diff(t, tt.A.err, err, cmpopts.EquateErrors())
			testutil.Diff(t, tt.A.msg, msg, cmpopts.IgnoreUnexported(k.Reference{}))
		})
	}
}

func TestParseID(t *testing.T) {
	type condition struct {
		msg protoreflect.ProtoMessage
	}

	type action struct {
		id  string
		err error
	}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"parse from kernel.Template", &condition{
				msg: &k.Resource{
					APIVersion: "test1",
					Kind:       "test2",
					Metadata: &k.Metadata{
						Namespace: "test3",
						Name:      "test4",
					},
				},
			},
			&action{
				id: "test1/test2/test3/test4",
			},
		),
		gen(
			"parse from kernel.Reference", &condition{
				msg: &k.Reference{
					APIVersion: "test1",
					Kind:       "test2",
					Namespace:  "test3",
					Name:       "test4",
				},
			},
			&action{
				id: "test1/test2/test3/test4",
			},
		),
		gen(
			"nil message", &condition{
				msg: nil,
			},
			&action{
				id: "",
				err: &er.Error{
					Package:     encoder.ErrPkg,
					Type:        encoder.ErrTypeJSON,
					Description: encoder.ErrDscUnmarshal,
				},
			},
		),
		gen(
			"nil pointer message", &condition{
				msg: *new(protoreflect.ProtoMessage),
			},
			&action{
				id: "",
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
			id, err := api.ParseID(tt.C.msg)
			testutil.Diff(t, tt.A.err, err, cmpopts.EquateErrors())
			testutil.Diff(t, tt.A.id, id)
		})
	}
}

type testContainer struct {
	objStore map[string]any
	err      error
}

func (c *testContainer) Serve(ctx context.Context, req *api.Request) (*api.Response, error) {
	return &api.Response{Content: c.objStore[req.Key]}, c.err
}

func TestReferObject(t *testing.T) {
	type condition struct {
		a   api.API[*api.Request, *api.Response]
		ref *k.Reference
	}

	type action struct {
		obj any
		err error
	}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"object found", &condition{
				a: &testContainer{
					objStore: map[string]any{"foo/bar": "test object"},
				},
				ref: &k.Reference{
					APIVersion: "foo",
					Kind:       "bar",
				},
			},
			&action{
				obj: "test object",
			},
		),
		gen(
			"error API", &condition{
				a: &testContainer{
					objStore: map[string]any{},
					err:      io.ErrUnexpectedEOF, // Dummy error
				},
				ref: &k.Reference{
					APIVersion: "foo",
					Kind:       "bar",
				},
			},
			&action{
				err: io.ErrUnexpectedEOF,
			},
		),
		gen(
			"nil reference", &condition{
				ref: nil,
			},
			&action{
				err: &er.Error{
					Package:     api.ErrPkg,
					Type:        api.ErrTypeUtil,
					Description: api.ErrDscNil,
				},
			},
		),
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			obj, err := api.ReferObject(tt.C.a, tt.C.ref)
			testutil.Diff(t, tt.A.err, err, cmpopts.EquateErrors())
			testutil.Diff(t, tt.A.obj, obj)
		})
	}
}

func TestReferTypedObject(t *testing.T) {
	type condition struct {
		a   api.API[*api.Request, *api.Response]
		ref *k.Reference
	}

	type action struct {
		obj any
		err error
	}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"object found", &condition{
				a: &testContainer{
					objStore: map[string]any{"foo/bar": "test object"},
				},
				ref: &k.Reference{
					APIVersion: "foo",
					Kind:       "bar",
				},
			},
			&action{
				obj: "test object",
			},
		),
		gen(
			"wrong value type", &condition{
				a: &testContainer{
					objStore: map[string]any{"foo/bar": 999}, // Expect string but the actual value is int.
				},
				ref: &k.Reference{
					APIVersion: "foo",
					Kind:       "bar",
				},
			},
			&action{
				obj: "", // Zero value of string is returned when error.
				err: &er.Error{
					Package:     api.ErrPkg,
					Type:        api.ErrTypeUtil,
					Description: api.ErrDscAssert,
				},
			},
		),
		gen(
			"error API", &condition{
				a: &testContainer{
					objStore: map[string]any{},
					err:      io.ErrUnexpectedEOF, // Dummy error
				},
				ref: &k.Reference{
					APIVersion: "foo",
					Kind:       "bar",
				},
			},
			&action{
				obj: "", // Zero value of string is returned when error.
				err: io.ErrUnexpectedEOF,
			},
		),
		gen(
			"nil reference", &condition{
				ref: nil,
			},
			&action{
				obj: "", // Zero value of string is returned when error.
				err: &er.Error{
					Package:     api.ErrPkg,
					Type:        api.ErrTypeUtil,
					Description: api.ErrDscNil,
				},
			},
		),
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			obj, err := api.ReferTypedObject[string](tt.C.a, tt.C.ref)
			testutil.Diff(t, tt.A.err, err, cmpopts.EquateErrors())
			testutil.Diff(t, tt.A.obj, obj)
		})
	}
}

func TestReferTypedObjects(t *testing.T) {
	type condition struct {
		a    api.API[*api.Request, *api.Response]
		refs []*k.Reference
	}

	type action struct {
		obj any
		err error
	}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"object found", &condition{
				a: &testContainer{
					objStore: map[string]any{"foo/bar": "test object"},
				},
				refs: []*k.Reference{
					{
						APIVersion: "foo",
						Kind:       "bar",
					},
				},
			},
			&action{
				obj: []string{"test object"},
			},
		),
		gen(
			"object found", &condition{
				a: &testContainer{
					objStore: map[string]any{"foo/bar": "test1 object", "alice/bob": "test2 object"},
				},
				refs: []*k.Reference{
					{
						APIVersion: "foo",
						Kind:       "bar",
					},
					{
						APIVersion: "alice",
						Kind:       "bob",
					},
				},
			},
			&action{
				obj: []string{"test1 object", "test2 object"},
			},
		),
		gen(
			"wrong value type", &condition{
				a: &testContainer{
					objStore: map[string]any{"foo/bar": 999}, // Expect string but the actual value is int.
				},
				refs: []*k.Reference{
					{
						APIVersion: "foo",
						Kind:       "bar",
					},
				},
			},
			&action{
				obj: []string(nil), // Typed nil slice is returned when error.
				err: &er.Error{
					Package:     api.ErrPkg,
					Type:        api.ErrTypeUtil,
					Description: api.ErrDscAssert,
				},
			},
		),
		gen(
			"error API", &condition{
				a: &testContainer{
					objStore: map[string]any{},
					err:      io.ErrUnexpectedEOF, // Dummy error
				},
				refs: []*k.Reference{
					{
						APIVersion: "foo",
						Kind:       "bar",
					},
				},
			},
			&action{
				obj: []string(nil), // Typed nil slice is returned when error.
				err: io.ErrUnexpectedEOF,
			},
		),
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			obj, err := api.ReferTypedObjects[string](tt.C.a, tt.C.refs...)
			testutil.Diff(t, tt.A.err, err, cmpopts.EquateErrors())
			testutil.Diff(t, tt.A.obj, obj)
		})
	}
}
