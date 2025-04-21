// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package api_test

import (
	"context"
	"io"
	"testing"

	k "github.com/aileron-gateway/aileron-gateway/apis/kernel"
	"github.com/aileron-gateway/aileron-gateway/kernel/api"
	"github.com/aileron-gateway/aileron-gateway/kernel/encoder"
	"github.com/aileron-gateway/aileron-gateway/kernel/er"
	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
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

	cndFormatJSON := "json"
	cndFormatYAML := "yaml"
	cndFormatUnsupported := "unsupported"
	actCheckResult := "check result"
	actCheckNoError := "no error"
	actCheckError := "non-nil error"

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	tb.Condition(cndFormatJSON, "unmarshal json")
	tb.Condition(cndFormatYAML, "unmarshal yaml")
	tb.Condition(cndFormatUnsupported, "input unsupported format")
	tb.Action(actCheckResult, "check the un-marshalled values")
	tb.Action(actCheckNoError, "check that there is no error")
	tb.Action(actCheckError, "check that a non-nil error was returned")
	table := tb.Build()

	type testStruct struct {
		Foo string `json:"foo" yaml:"foo"`
		Bar int    `json:"bar" yaml:"bar"`
	}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"json",
			[]string{cndFormatJSON},
			[]string{actCheckResult, actCheckNoError},
			&condition{
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
			"yaml",
			[]string{cndFormatYAML},
			[]string{actCheckResult, actCheckNoError},
			&condition{
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
			"unsupported",
			[]string{},
			[]string{actCheckError},
			&condition{
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

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			err := tt.C().format.Unmarshal(tt.C().in, tt.C().into)
			testutil.Diff(t, tt.A().err, err, cmpopts.EquateErrors())
			testutil.Diff(t, tt.A().result, tt.C().into)
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

	cndValueExists := "value exists"
	cndNilContext := "nil api"
	cndNilAPI := "nil api"
	actCheckAPIs := "check apis"

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	tb.Condition(cndValueExists, "at least 1 API was already exist")
	tb.Condition(cndNilContext, "register multiple APIs")
	tb.Condition(cndNilAPI, "input nil API")
	tb.Action(actCheckAPIs, "check the APIs saved in the context")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"new",
			[]string{},
			[]string{actCheckAPIs},
			&condition{
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
			"append",
			[]string{cndValueExists},
			[]string{actCheckAPIs},
			&condition{
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
			"append nil",
			[]string{cndValueExists},
			[]string{actCheckAPIs},
			&condition{
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
			"nil context",
			[]string{cndNilContext},
			[]string{actCheckAPIs},
			&condition{
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

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			ctx := api.ContextWithRoute(tt.C().ctx, tt.C().a)
			testutil.Diff(t, tt.A().apis, ctx.Value(api.APIRouteContextKey))
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

	cndNoAPI := "no API in the context"
	cndSingleAPI := "only 1 API in the context"
	cndMultipleAPIs := "multiple APIs in the context"
	cndNilContext := "input nil context"
	actCheckAPI := "check returned API"
	actCheckNil := "check nil was returned"

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	tb.Condition(cndNoAPI, "no APIs were registered")
	tb.Condition(cndSingleAPI, "1 API was registered")
	tb.Condition(cndMultipleAPIs, "multiple APIs were registered")
	tb.Condition(cndNilContext, "input nil as context")
	tb.Action(actCheckAPI, "check the non-nil returned API was the one expected")
	tb.Action(actCheckNil, "check that nil was returned")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"api not registered",
			[]string{cndNoAPI},
			[]string{actCheckNil},
			&condition{
				ctx: context.WithValue(context.Background(), api.APIRouteContextKey,
					[]api.API[*api.Request, *api.Response]{}),
			},
			&action{
				a: nil,
			},
		),
		gen(
			"1 api registered",
			[]string{cndSingleAPI},
			[]string{actCheckAPI},
			&condition{
				ctx: context.WithValue(context.Background(), api.APIRouteContextKey,
					[]api.API[*api.Request, *api.Response]{stringAPI("test")}),
			},
			&action{
				a: stringAPI("test"),
			},
		),
		gen(
			"multiple apis registered",
			[]string{cndMultipleAPIs},
			[]string{actCheckAPI},
			&condition{
				ctx: context.WithValue(context.Background(), api.APIRouteContextKey,
					[]api.API[*api.Request, *api.Response]{stringAPI("test1"), stringAPI("test2")}),
			},
			&action{
				a: stringAPI("test1"),
			},
		),
		gen(
			"nil context",
			[]string{cndNilContext},
			[]string{actCheckNil},
			&condition{
				ctx: nil,
			},
			&action{
				a: nil,
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			a := api.RootAPIFromContext(tt.C().ctx)
			testutil.Diff(t, tt.A().a, a)
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

	cndFormatJSON := "json"
	cndFormatYAML := "yaml"
	cndFormatProtoMessage := "protoMessage"
	cndFormatProtoReference := "protoReference"
	cndFormatUnsupported := "unsupported format"
	cndMerge := "merge with default values"
	cndInvalidContentType := "invalid content type"
	actCheckReturnedValues := "check the returned values"
	actCheckNoError := "no error"
	actCheckError := "non-nil error"

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	tb.Condition(cndFormatJSON, "JSON format")
	tb.Condition(cndFormatYAML, "YAML format")
	tb.Condition(cndFormatProtoMessage, "ProtoMessage format")
	tb.Condition(cndFormatProtoReference, "ProtoReference format")
	tb.Condition(cndFormatUnsupported, "Unsupported format")
	tb.Condition(cndMerge, "merge with default values")
	tb.Condition(cndInvalidContentType, "the type of the content is invalid")
	tb.Action(actCheckReturnedValues, "check values of the returned proto message")
	tb.Action(actCheckNoError, "check that there is no error")
	tb.Action(actCheckError, "check that a non-nil error was returned")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"json/success",
			[]string{cndFormatJSON},
			[]string{actCheckReturnedValues, actCheckNoError},
			&condition{
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
			"json/merge",
			[]string{cndFormatJSON, cndMerge},
			[]string{actCheckReturnedValues, actCheckNoError},
			&condition{
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
			"json/error",
			[]string{cndFormatJSON, cndInvalidContentType},
			[]string{actCheckReturnedValues, actCheckError},
			&condition{
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
			"yaml/success",
			[]string{cndFormatYAML},
			[]string{actCheckReturnedValues, actCheckNoError},
			&condition{
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
			"yaml/merge",
			[]string{cndFormatYAML, cndMerge},
			[]string{actCheckReturnedValues, actCheckNoError},
			&condition{
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
			"yaml/error",
			[]string{cndFormatYAML, cndInvalidContentType},
			[]string{actCheckReturnedValues, actCheckError},
			&condition{
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
			"protoMessage/success",
			[]string{cndFormatProtoMessage},
			[]string{actCheckReturnedValues, actCheckNoError},
			&condition{
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
			"protoMessage/merge",
			[]string{cndFormatProtoMessage, cndMerge},
			[]string{actCheckReturnedValues, actCheckNoError},
			&condition{
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
			"protoMessage/error",
			[]string{cndFormatProtoMessage, cndInvalidContentType},
			[]string{actCheckReturnedValues, actCheckError},
			&condition{
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
			"protoReference/success",
			[]string{cndFormatProtoReference},
			[]string{actCheckReturnedValues, actCheckNoError},
			&condition{
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
			"protoReference/merge",
			[]string{cndFormatProtoReference, cndMerge},
			[]string{actCheckReturnedValues, actCheckNoError},
			&condition{
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
			"protoReference/error",
			[]string{cndFormatProtoReference, cndInvalidContentType},
			[]string{actCheckReturnedValues, actCheckError},
			&condition{
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
			"unsupported format",
			[]string{cndFormatUnsupported},
			[]string{actCheckReturnedValues, actCheckNoError},
			&condition{
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

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			msg, err := api.ProtoMessage(tt.C().format, tt.C().content, tt.C().msg, tt.C().opt)

			testutil.Diff(t, tt.A().err, err, cmpopts.EquateErrors())
			testutil.Diff(t, tt.A().msg, msg, cmpopts.IgnoreUnexported(k.Reference{}))
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

	cndTemplate := "from kernel.template"
	cndReference := "from kernel.Reference"
	cndInvalidMessage := "invalid message"
	actCheckID := "returned ID"
	actCheckNoError := "no error"
	actCheckError := "non-nil error"

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	tb.Condition(cndTemplate, "merge with default values")
	tb.Condition(cndReference, "merge with default values")
	tb.Condition(cndInvalidMessage, "the type of the content is invalid")
	tb.Action(actCheckID, "check values of the returned proto message")
	tb.Action(actCheckNoError, "check that there is no error")
	tb.Action(actCheckError, "check that a non-nil error was returned")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"parse from kernel.Template",
			[]string{cndTemplate},
			[]string{actCheckID, actCheckNoError},
			&condition{
				msg: &k.Template{
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
			"parse from kernel.Reference",
			[]string{cndReference},
			[]string{actCheckID, actCheckNoError},
			&condition{
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
			"nil message",
			[]string{cndInvalidMessage},
			[]string{actCheckID, actCheckError},
			&condition{
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
			"nil pointer message",
			[]string{cndInvalidMessage},
			[]string{actCheckID, actCheckError},
			&condition{
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

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			id, err := api.ParseID(tt.C().msg)
			testutil.Diff(t, tt.A().err, err, cmpopts.EquateErrors())
			testutil.Diff(t, tt.A().id, id)
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

	cndObjectExists := "object exists"
	cndErrorAPI := "API error"
	cndNilReference := "input nil reference"
	actCheckObject := "returned object"
	actCheckNoError := "no error"
	actCheckError := "non-nil error"

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	tb.Condition(cndObjectExists, "referred object exists in the API")
	tb.Condition(cndErrorAPI, "API returns an error")
	tb.Condition(cndNilReference, "input nil reference")
	tb.Action(actCheckObject, "check the returned object")
	tb.Action(actCheckNoError, "check that there is no error")
	tb.Action(actCheckError, "check that a non-nil error was returned")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"object found",
			[]string{cndObjectExists},
			[]string{actCheckObject, actCheckNoError},
			&condition{
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
			"error API",
			[]string{cndErrorAPI},
			[]string{actCheckObject, actCheckError},
			&condition{
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
			"nil reference",
			[]string{cndNilReference},
			[]string{actCheckObject, actCheckError},
			&condition{
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

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			obj, err := api.ReferObject(tt.C().a, tt.C().ref)
			testutil.Diff(t, tt.A().err, err, cmpopts.EquateErrors())
			testutil.Diff(t, tt.A().obj, obj)
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

	cndObjectExists := "object exists"
	cndErrorAPI := "API error"
	cndNilReference := "input nil reference"
	cndWrongType := "wrong type"
	actCheckObject := "returned object"
	actCheckNoError := "no error"
	actCheckError := "non-nil error"

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	tb.Condition(cndObjectExists, "referred object exists in the API")
	tb.Condition(cndErrorAPI, "API returns an error")
	tb.Condition(cndNilReference, "input nil reference")
	tb.Condition(cndWrongType, "expected wrong type")
	tb.Action(actCheckObject, "check the returned object")
	tb.Action(actCheckNoError, "check that there is no error")
	tb.Action(actCheckError, "check that a non-nil error was returned")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"object found",
			[]string{cndObjectExists},
			[]string{actCheckObject, actCheckNoError},
			&condition{
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
			"wrong value type",
			[]string{cndObjectExists, cndWrongType},
			[]string{actCheckObject, actCheckError},
			&condition{
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
			"error API",
			[]string{cndErrorAPI},
			[]string{actCheckObject, actCheckError},
			&condition{
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
			"nil reference",
			[]string{cndNilReference},
			[]string{actCheckObject, actCheckError},
			&condition{
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

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			obj, err := api.ReferTypedObject[string](tt.C().a, tt.C().ref)
			testutil.Diff(t, tt.A().err, err, cmpopts.EquateErrors())
			testutil.Diff(t, tt.A().obj, obj)
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

	cndObjectExists := "object exists"
	cndErrorAPI := "API error"
	cndWrongType := "wrong type"
	actCheckObject := "returned object"
	actCheckNoError := "no error"
	actCheckError := "non-nil error"

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	tb.Condition(cndObjectExists, "referred object exists in the API")
	tb.Condition(cndErrorAPI, "API returns an error")
	tb.Condition(cndWrongType, "expected wrong type")
	tb.Action(actCheckObject, "check the returned object")
	tb.Action(actCheckNoError, "check that there is no error")
	tb.Action(actCheckError, "check that a non-nil error was returned")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"object found",
			[]string{cndObjectExists},
			[]string{actCheckObject, actCheckNoError},
			&condition{
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
			"object found",
			[]string{cndObjectExists},
			[]string{actCheckObject, actCheckNoError},
			&condition{
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
			"wrong value type",
			[]string{cndObjectExists, cndWrongType},
			[]string{actCheckObject, actCheckError},
			&condition{
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
			"error API",
			[]string{cndErrorAPI},
			[]string{actCheckObject, actCheckError},
			&condition{
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

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			obj, err := api.ReferTypedObjects[string](tt.C().a, tt.C().refs...)
			testutil.Diff(t, tt.A().err, err, cmpopts.EquateErrors())
			testutil.Diff(t, tt.A().obj, obj)
		})
	}
}
