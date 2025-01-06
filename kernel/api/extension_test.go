package api_test

import (
	"context"
	"fmt"
	"io"
	"testing"

	k "github.com/aileron-gateway/aileron-gateway/apis/kernel"
	"github.com/aileron-gateway/aileron-gateway/kernel/api"
	"github.com/aileron-gateway/aileron-gateway/kernel/encoder"
	"github.com/aileron-gateway/aileron-gateway/kernel/er"
	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
	"github.com/google/go-cmp/cmp/cmpopts"
)

type MyCreator struct{}

func (c *MyCreator) Create(a api.API[*api.Request, *api.Response], f api.Format, manifest any) (any, error) {
	// Suppose that yaml manifest with the following structure is given
	// with the manifest argument in this example.
	// Note that the manifest is not always YAML but it depends on the usage.
	m := &struct {
		APIVersion string `yaml:"apiVersion"`
		Kind       string `yaml:"kind"`
		Metadata   *struct {
			Name      string `yaml:"name"`
			Namespace string `yaml:"namespace"`
		} `yaml:"metadata"`
		Spec *struct { // We use only this field in this example. So, other fields can be omitted.
			Value string `yaml:"value"`
		} `yaml:"spec"`
	}{}

	if err := encoder.UnmarshalYAML(manifest.([]byte), m); err != nil {
		return nil, err
	}
	return m.Spec.Value, nil
}

func ExampleExtensionAPI() {
	// Create a new Extension API.
	e := api.NewExtensionAPI()

	// Register a creator.
	// The key is supposed to be "APIGroup/APIVersion/Kind".
	e.Register("extension/v1/MyCreator", &MyCreator{})

	postReq := &api.Request{
		Method: api.MethodPost,
		Key:    "extension/v1/MyCreator", // Must be the same as registered key.
		Format: api.FormatYAML,           // Must be specified correctly.
		Content: []byte(`
apiVersion: extension/v1
kind: MyCreator
metadata:
    namespace: example
    name: test
spec:
    value: Hello from ExtensionAPI`),
	}
	if _, err := e.Serve(context.Background(), postReq); err != nil { // Store the object.
		panic("handle error here")
	}

	getReq := &api.Request{
		Method: api.MethodGet,
		Key:    "extension/v1/MyCreator", // Must be the same as registered key.
		Format: api.FormatProtoReference, // Must be specified correctly.
		Content: &k.Reference{
			APIVersion: "extension/v1",
			Kind:       "MyCreator",
			Namespace:  "example",
			Name:       "test",
		},
	}
	getResp, err := e.Serve(context.Background(), getReq) // Get the stored object.
	if err != nil {
		panic("handle error here")
	}
	fmt.Println(getResp.Content)
	// Output:
	// 	Hello from ExtensionAPI

	deleteReq := &api.Request{
		Method: api.MethodDelete,
		Key:    "extension/v1/MyCreator",
		Format: api.FormatProtoReference,
		Content: &k.Reference{
			APIVersion: "extension/v1",
			Kind:       "MyCreator",
			Namespace:  "example",
			Name:       "test",
		},
	}
	if _, err = e.Serve(context.Background(), deleteReq); err != nil {
		panic("handle error here")
	}
}

type stringCreator string

func (c stringCreator) Create(a api.API[*api.Request, *api.Response], f api.Format, manifest any) (any, error) {
	return string(c), nil
}

type errorCreator struct {
	err error
}

func (c errorCreator) Create(a api.API[*api.Request, *api.Response], f api.Format, manifest any) (any, error) {
	return nil, c.err
}

func TestExtensionAPI_Serve(t *testing.T) {
	type condition struct {
		creators map[string]api.Creator
		reqs     []*api.Request
	}

	type action struct {
		res *api.Response // Result for the final request in the reqs.
		err error
	}

	cndPost := "register"
	cndDelete := "delete"
	cndGet := "get"
	cndInstanceExists := "instance exists"
	cndErrorCreate := "error on create"
	cndWrongType := "wrong content"
	cndNilRequest := "nil request"
	cndDuplicateKey := "duplicate key"
	cndNoCreator := "creator not registered"
	cndUnsupportedMethod := "unsupported method"
	actCheckResponse := "check response"
	actCheckNoError := "no error"
	actCheckError := "non-nil error"

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	tb.Condition(cndPost, "send Post request")
	tb.Condition(cndDelete, "send Delete request")
	tb.Condition(cndGet, "send Get request")
	tb.Condition(cndInstanceExists, "an instance has already been created")
	tb.Condition(cndErrorCreate, "create method returns an error")
	tb.Condition(cndWrongType, "content in the request is invalid")
	tb.Condition(cndNilRequest, "input nil request")
	tb.Condition(cndDuplicateKey, "try to register with a duplicate key")
	tb.Condition(cndNoCreator, "creator is not registered in the API")
	tb.Condition(cndUnsupportedMethod, "request with an unsupported method")
	tb.Action(actCheckResponse, "check the returned response")
	tb.Action(actCheckNoError, "check that there is no error")
	tb.Action(actCheckError, "check that a non-nil error was returned")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"Post an object",
			[]string{cndPost, cndGet},
			[]string{actCheckResponse, actCheckNoError},
			&condition{
				creators: map[string]api.Creator{
					"test1/test2": stringCreator("hello from test creator"),
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
					Content: "hello from test creator",
				},
			},
		),
		gen(
			"Delete an object",
			[]string{cndPost, cndDelete, cndGet},
			[]string{actCheckResponse, actCheckError},
			&condition{
				creators: map[string]api.Creator{
					"test1/test2": stringCreator("hello from test creator"),
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
					Type:        api.ErrTypeExt,
					Description: api.ErrDscNoManifest,
				},
			},
		),
		gen(
			"Get already created object",
			[]string{cndPost, cndGet, cndInstanceExists},
			[]string{actCheckResponse, actCheckNoError},
			&condition{
				creators: map[string]api.Creator{
					"test1/test2": stringCreator("hello from test creator"),
				},
				reqs: []*api.Request{
					{
						Method:  api.MethodPost, // Register a manifest for instantiation.
						Key:     "test1/test2",
						Format:  api.FormatJSON,
						Content: []byte(`{"apiVersion":"test1", "kind":"test2", "metadata": {"namespace":"test3", "name":"test4"}}`),
					},
					{
						Method:  api.MethodGet, // A new instance will be created with this request.
						Key:     "test1/test2",
						Format:  api.FormatJSON,
						Content: []byte(`{"apiVersion":"test1", "kind":"test2", "metadata": {"namespace":"test3", "name":"test4"}}`),
					},
					{
						Method:  api.MethodGet, // A new instance won't be created with this request. The existing one will be returned.
						Key:     "test1/test2",
						Format:  api.FormatJSON,
						Content: []byte(`{"apiVersion":"test1", "kind":"test2", "metadata": {"namespace":"test3", "name":"test4"}}`),
					},
				},
			},
			&action{
				res: &api.Response{
					Content: "hello from test creator",
				},
			},
		),
		gen(
			"Error on creation",
			[]string{cndPost, cndGet, cndErrorCreate},
			[]string{actCheckResponse, actCheckError},
			&condition{
				creators: map[string]api.Creator{
					"test1/test2": errorCreator{err: io.ErrUnexpectedEOF}, // Use dummy error
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
				res: nil,
				err: io.ErrUnexpectedEOF,
			},
		),
		gen(
			"nil request",
			[]string{cndNilRequest},
			[]string{actCheckResponse, actCheckError},
			&condition{
				creators: map[string]api.Creator{},
				reqs: []*api.Request{
					nil,
				},
			},
			&action{
				res: nil,
				err: &er.Error{
					Package:     api.ErrPkg,
					Type:        api.ErrTypeExt,
					Description: api.ErrDscNil,
				},
			},
		),
		gen(
			"invalid content",
			[]string{cndWrongType},
			[]string{actCheckResponse, actCheckError},
			&condition{
				creators: map[string]api.Creator{
					"test1/test2": stringCreator("hello from test creator"),
				},
				reqs: []*api.Request{
					{
						Method:  api.MethodPost, // Register a manifest for instantiation.
						Key:     "test1/test2",
						Format:  api.FormatJSON,
						Content: 999, // FormatJSON require []byte JSON content. But set int here to make an error.
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
			"register with duplicate key",
			[]string{cndDuplicateKey},
			[]string{actCheckResponse, actCheckError},
			&condition{
				creators: map[string]api.Creator{
					"test1/test2": stringCreator("hello from test creator"),
				},
				reqs: []*api.Request{
					{
						Method:  api.MethodPost, // Register a manifest for instantiation.
						Key:     "test1/test2",
						Format:  api.FormatJSON,
						Content: []byte(`{"apiVersion":"test1", "kind":"test2", "metadata": {"namespace":"test3", "name":"test4"}}`),
					},
					{
						Method:  api.MethodPost, // Register a manifest with the same key "test1/test2/test3/test4".
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
					Type:        api.ErrTypeExt,
					Description: api.ErrDscDuplicateKey,
				},
			},
		),
		gen(
			"no creator",
			[]string{cndNoCreator},
			[]string{actCheckResponse, actCheckError},
			&condition{
				creators: map[string]api.Creator{},
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
					Type:        api.ErrTypeExt,
					Description: api.ErrDscNoAPI,
				},
			},
		),
		gen(
			"unsupported method",
			[]string{cndUnsupportedMethod},
			[]string{actCheckResponse, actCheckError},
			&condition{
				creators: map[string]api.Creator{
					"test1/test2": stringCreator("hello from test creator"),
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
					Type:        api.ErrTypeExt,
					Description: api.ErrDscNoMethod,
				},
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			a := api.NewExtensionAPI()
			for k, v := range tt.C().creators {
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
