package api_test

import (
	"context"
	"fmt"
	"testing"

	k "github.com/aileron-gateway/aileron-gateway/apis/kernel"
	"github.com/aileron-gateway/aileron-gateway/kernel/api"
	"github.com/aileron-gateway/aileron-gateway/kernel/er"
	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
	"github.com/google/go-cmp/cmp/cmpopts"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type MyResource struct {
	*api.BaseResource
}

func (r *MyResource) Default() protoreflect.ProtoMessage {
	// Use template message for this example.
	return &k.Template{
		APIVersion: "factory/v1",
		Kind:       "MyResource",
		Metadata: &k.Metadata{
			Namespace: "default",
			Name:      "default",
		},
	}
}

func (r *MyResource) Create(a api.API[*api.Request, *api.Response], msg protoreflect.ProtoMessage) (any, error) {
	c := msg.(*k.Template)
	// Just return the namespace and name values in the manifest
	// because the kernel.Template message does not contain any meaningful fields
	// that can be used in this example test.
	return c.Metadata.Namespace + " " + c.Metadata.Name, nil
}

func ExampleFactoryAPI() {
	// Create a new Factory API.
	f := api.NewFactoryAPI()

	// Register a resource.
	// The key is supposed to be "APIGroup/APIVersion/Kind".
	f.Register("factory/v1/MyResource", &MyResource{&api.BaseResource{}})

	postReq := &api.Request{
		Method: api.MethodPost,
		Key:    "factory/v1/MyResource", // Must be the same as registered key.
		Format: api.FormatYAML,          // Must be specified correctly.
		Content: []byte(`
apiVersion: factory/v1
kind: MyResource
metadata:
    namespace: Hello
    name: FactoryAPI
spec: {}`),
	}
	if _, err := f.Serve(context.Background(), postReq); err != nil { // Store the object.
		panic("handle error here")
	}

	getReq := &api.Request{
		Method: api.MethodGet,
		Key:    "factory/v1/MyResource",  // Must be the same as registered key.
		Format: api.FormatProtoReference, // Must be specified correctly.
		Content: &k.Reference{
			APIVersion: "factory/v1",
			Kind:       "MyResource",
			Namespace:  "Hello",
			Name:       "FactoryAPI",
		},
	}
	getResp, err := f.Serve(context.Background(), getReq) // Get the stored object.
	if err != nil {
		panic("handle error here")
	}
	fmt.Println(getResp.Content)
	// Output:
	// 	Hello FactoryAPI

	deleteReq := &api.Request{
		Method: api.MethodDelete,
		Key:    "factory/v1/MyResource",
		Format: api.FormatProtoReference,
		Content: &k.Reference{
			APIVersion: "factory/v1",
			Kind:       "MyResource",
			Namespace:  "Hello",
			Name:       "FactoryAPI",
		},
	}
	if _, err = f.Serve(context.Background(), deleteReq); err != nil {
		panic("handle error here")
	}
}

func TestBaseResource(t *testing.T) {
	type condition struct {
		r *api.BaseResource
	}

	type action struct {
		err error // validation error
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"nil manifest",
			[]string{},
			[]string{},
			&condition{
				r: &api.BaseResource{
					DefaultProto: nil,
				},
			},
			&action{
				err: nil,
			},
		),
		gen(
			"valid manifest",
			[]string{},
			[]string{},
			&condition{
				r: &api.BaseResource{
					DefaultProto: &k.Reference{
						APIVersion: "test/v1",
						Kind:       "Test",
						Name:       "foo",
						Namespace:  "bar",
					},
				},
			},
			&action{
				err: nil,
			},
		),
		gen(
			"invalid manifest",
			[]string{},
			[]string{},
			&condition{
				r: &api.BaseResource{
					DefaultProto: &k.Reference{
						APIVersion: "test/v1",
						Kind:       "Invalid Kind",
						Name:       "foo",
						Namespace:  "bar",
					},
				},
			},
			&action{
				err: &er.Error{
					Package:     api.ErrPkg,
					Type:        api.ErrTypeFactory,
					Description: api.ErrDscProtoValidate,
				},
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			d := tt.C().r.Default()
			err := tt.C().r.Validate(d)
			testutil.Diff(t, tt.A().err, err, cmpopts.EquateErrors())
		})
	}
}

func TestFactoryAPI_Serve(t *testing.T) {
	type condition struct {
		resources map[string]api.Resource
		reqs      []*api.Request
	}

	type action struct {
		res *api.Response // Result for the final request in the reqs.
		err error
	}

	cndPost := "register"
	cndDelete := "delete"
	cndGet := "get"
	cndWrongType := "wrong content"
	cndNilRequest := "nil request"
	cndDuplicateKey := "duplicate key"
	cndNoResource := "resource is not registered"
	cndUnsupportedMethod := "unsupported method"
	actCheckResponse := "check response"
	actCheckNoError := "no error"
	actCheckError := "non-nil error"

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	tb.Condition(cndPost, "send Post request")
	tb.Condition(cndDelete, "send Delete request")
	tb.Condition(cndGet, "send Get request")
	tb.Condition(cndWrongType, "content in the request is invalid")
	tb.Condition(cndNilRequest, "input nil request")
	tb.Condition(cndDuplicateKey, "try to register with a duplicate key")
	tb.Condition(cndNoResource, "resource is not registered in the API")
	tb.Condition(cndUnsupportedMethod, "request with an unsupported method")
	tb.Action(actCheckResponse, "check the returned response")
	tb.Action(actCheckNoError, "check that there is no error")
	tb.Action(actCheckError, "check that a non-nil error was returned")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"Post successful",
			[]string{cndPost, cndGet},
			[]string{actCheckResponse, actCheckNoError},
			&condition{
				resources: map[string]api.Resource{
					"test1/test2": &MyResource{},
				},
				reqs: []*api.Request{
					{
						Method:  api.MethodPost, // Register a manifest for instantiation.
						Key:     "test1/test2",
						Format:  api.FormatJSON,
						Content: []byte(`{"apiVersion":"test1", "kind":"test2", "metadata": {"namespace":"test3", "name":"test4"}}`),
					},
					{
						Method:  api.MethodGet,
						Key:     "test1/test2",
						Format:  api.FormatJSON,
						Content: []byte(`{"apiVersion":"test1", "kind":"test2", "metadata": {"namespace":"test3", "name":"test4"}}`),
					},
				},
			},
			&action{
				res: &api.Response{
					Params:  map[string]string{},
					Content: "test3 test4", // MyResource creates string "<namespace> <name>".
				},
			},
		),
		gen(
			"Delete successful",
			[]string{cndPost, cndDelete, cndGet},
			[]string{actCheckResponse, actCheckError},
			&condition{
				resources: map[string]api.Resource{
					"test1/test2": &MyResource{},
				},
				reqs: []*api.Request{
					{
						Method:  api.MethodPost, // Register a manifest for instantiation.
						Key:     "test1/test2",
						Format:  api.FormatJSON,
						Content: []byte(`{"apiVersion":"test1", "kind":"test2", "metadata": {"namespace":"test3", "name":"test4"}}`),
					},
					{
						Method:  api.MethodDelete,
						Key:     "test1/test2",
						Format:  api.FormatJSON,
						Content: []byte(`{"apiVersion":"test1", "kind":"test2", "metadata": {"namespace":"test3", "name":"test4"}}`),
					},
					{
						Method:  api.MethodGet,
						Key:     "test1/test2",
						Format:  api.FormatJSON,
						Content: []byte(`{"apiVersion":"test1", "kind":"test2", "metadata": {"namespace":"test3", "name":"test4"}}`),
					},
				},
			},
			&action{
				res: nil,
				err: &er.Error{
					Package:     api.ErrPkg,
					Type:        api.ErrTypeFactory,
					Description: api.ErrDscNoManifest,
				},
			},
		),
		gen(
			"Delete fails",
			[]string{cndDelete, cndWrongType},
			[]string{actCheckResponse, actCheckError},
			&condition{
				resources: map[string]api.Resource{
					"test1/test2": &MyResource{},
				},
				reqs: []*api.Request{
					{
						Method:  api.MethodDelete,
						Key:     "test1/test2",
						Format:  api.FormatJSON,
						Content: 999, // FormatJSON require []byte JSON but set int to make an error.
					},
				},
			},
			&action{
				res: nil,
				err: &er.Error{
					Package:     api.ErrPkg,
					Type:        api.ErrTypeUtil,
					Description: api.ErrDscAssert,
				},
			},
		),
		gen(
			"Get fails",
			[]string{cndGet, cndWrongType},
			[]string{actCheckResponse, actCheckError},
			&condition{
				resources: map[string]api.Resource{
					"test1/test2": &MyResource{},
				},
				reqs: []*api.Request{
					{
						Method:  api.MethodGet,
						Key:     "test1/test2",
						Format:  api.FormatJSON,
						Content: 999, // FormatJSON require []byte JSON but set int to make an error.
					},
				},
			},
			&action{
				res: nil,
				err: &er.Error{
					Package:     api.ErrPkg,
					Type:        api.ErrTypeUtil,
					Description: api.ErrDscAssert,
				},
			},
		),
		gen(
			"Post fails",
			[]string{cndPost, cndWrongType},
			[]string{actCheckResponse, actCheckError},
			&condition{
				resources: map[string]api.Resource{
					"test1/test2": &MyResource{},
				},
				reqs: []*api.Request{
					{
						Method:  api.MethodPost,
						Key:     "test1/test2",
						Format:  api.FormatJSON,
						Content: 999, // FormatJSON require []byte JSON but set int to make an error.
					},
				},
			},
			&action{
				res: nil,
				err: &er.Error{
					Package:     api.ErrPkg,
					Type:        api.ErrTypeUtil,
					Description: api.ErrDscAssert,
				},
			},
		),
		gen(
			"nil request",
			[]string{cndNilRequest},
			[]string{actCheckResponse, actCheckError},
			&condition{
				resources: map[string]api.Resource{},
				reqs: []*api.Request{
					nil,
				},
			},
			&action{
				res: nil,
				err: &er.Error{
					Package:     api.ErrPkg,
					Type:        api.ErrTypeFactory,
					Description: api.ErrDscNil,
				},
			},
		),
		gen(
			"no resource",
			[]string{cndNoResource},
			[]string{actCheckResponse, actCheckError},
			&condition{
				resources: map[string]api.Resource{},
				reqs: []*api.Request{
					{
						Method: api.MethodGet,
						Key:    "test1/test2",
					},
				},
			},
			&action{
				res: nil,
				err: &er.Error{
					Package:     api.ErrPkg,
					Type:        api.ErrTypeFactory,
					Description: api.ErrDscNoAPI,
				},
			},
		),
		gen(
			"unsupported method",
			[]string{cndUnsupportedMethod},
			[]string{actCheckResponse, actCheckError},
			&condition{
				resources: map[string]api.Resource{
					"test1/test2": &MyResource{},
				},
				reqs: []*api.Request{
					{
						Method: api.Method("UNSUPPORTED METHOD"),
						Key:    "test1/test2",
					},
				},
			},
			&action{
				res: nil,
				err: &er.Error{
					Package:     api.ErrPkg,
					Type:        api.ErrTypeFactory,
					Description: api.ErrDscNoMethod,
				},
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			a := api.NewFactoryAPI()
			for k, v := range tt.C().resources {
				a.Register(k, v)
			}

			var res *api.Response
			var err error

			ctx := context.Background()
			for _, r := range tt.C().reqs {
				res, err = a.Serve(ctx, r)
			}

			// Check the response for the final request.
			testutil.Diff(t, tt.A().err, err, cmpopts.EquateErrors())
			testutil.Diff(t, tt.A().res, res)
		})
	}
}
