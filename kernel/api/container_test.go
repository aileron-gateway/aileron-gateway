// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

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
)

func ExampleContainerAPI_one() {
	// Create a new Container API.
	c := api.NewContainerAPI()

	// Any types of objects can be stored.
	// Use string object so that this example test can be verified.
	obj := "Hello from ContainerAPI"

	postReq := &api.Request{
		Method: api.MethodPost,
		// Key to the object to store typically in the format of "APIGroup/APIVersion/Kind/namespace/Name".
		Key:     "example/v1/Container/test/object",
		Content: obj,
	}
	if _, err := c.Serve(context.Background(), postReq); err != nil { // Store the object.
		panic("handle error here")
	}

	getReq := &api.Request{
		Method: api.MethodGet,
		Key:    "example/v1/Container/test/object",
	}
	getResp, err := c.Serve(context.Background(), getReq) // Get the stored object.
	if err != nil {
		panic("handle error here")
	}
	fmt.Println(getResp.Content)
	// Output:
	// 	Hello from ContainerAPI

	deleteReq := &api.Request{
		Method: api.MethodDelete,
		Key:    "example/v1/Container/test/object",
	}
	if _, err = c.Serve(context.Background(), deleteReq); err != nil {
		panic("handle error here")
	}
}

func ExampleContainerAPI_two() {
	// Create a new Container API.
	c := api.NewContainerAPI()

	// Any types of objects can be stored.
	// Use string object so that this example test can be verified.
	obj := "Hello from ContainerAPI"

	postReq := &api.Request{
		Method: api.MethodPost,
		// Key to the object to store typically in the format of "APIGroup/APIVersion/Kind/namespace/Name".
		Key:     "example/v1/Container/test/object",
		Content: obj,
	}
	if _, err := c.Serve(context.Background(), postReq); err != nil { // Store the object.
		panic("handle error here")
	}

	getReq := &api.Request{
		Method: api.MethodGet,
		Key:    "",                       // Key is built from the content in the format of "APIVersion/Kind/Namespace/Name".
		Format: api.FormatProtoReference, // The format of the content must be specified correctly.
		Content: &k.Reference{
			APIVersion: "example/v1",
			Kind:       "Container",
			Namespace:  "test",
			Name:       "object",
		},
	}
	getResp, err := c.Serve(context.Background(), getReq)
	if err != nil {
		panic("handle error here")
	}
	fmt.Println(getResp.Content)
	// Output:
	// 	Hello from ContainerAPI

	deleteReq := &api.Request{
		Method: api.MethodDelete,
		Key:    "",
		Format: api.FormatProtoReference,
		Content: &k.Reference{
			APIVersion: "example/v1",
			Kind:       "Container",
			Namespace:  "test",
			Name:       "object",
		},
	}
	if _, err = c.Serve(context.Background(), deleteReq); err != nil {
		panic("handle error here")
	}
}

func TestContainerAPI_Serve(t *testing.T) {
	type condition struct {
		a    api.API[*api.Request, *api.Response]
		reqs []*api.Request
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
				a: api.NewContainerAPI(),
				reqs: []*api.Request{
					{
						Method:  api.MethodPost, // Register an object.
						Key:     "test key string",
						Content: "test object",
					},
					{
						Method: api.MethodGet, // Get registered object.
						Key:    "test key string",
					},
				},
			},
			&action{
				res: &api.Response{
					Content: "test object",
				},
			},
		),
		gen(
			"Delete an object",
			[]string{cndPost, cndDelete, cndGet},
			[]string{actCheckResponse, actCheckNoError},
			&condition{
				a: api.NewContainerAPI(),
				reqs: []*api.Request{
					{
						Method:  api.MethodPost, // Register an object.
						Key:     "test key string",
						Content: "test object",
					},
					{
						Method: api.MethodDelete, // Delete the registered object.
						Key:    "test key string",
					},
					{
						Method: api.MethodGet, // Get registered object. But the object should be deleted.
						Key:    "test key string",
					},
				},
			},
			&action{
				res: &api.Response{
					Content: nil, // Nil because the object was deleted.
				},
			},
		),
		gen(
			"Get with reference",
			[]string{cndPost, cndGet},
			[]string{actCheckResponse, actCheckNoError},
			&condition{
				a: api.NewContainerAPI(),
				reqs: []*api.Request{
					{
						Method:  api.MethodPost, // Register an object.
						Key:     "test1/test2/test3/test4",
						Content: "test object",
					},
					{
						Method: api.MethodGet, // Get registered object.
						Format: api.FormatProtoReference,
						Content: &k.Reference{
							APIVersion: "test1",
							Kind:       "test2",
							Namespace:  "test3",
							Name:       "test4",
						},
					},
				},
			},
			&action{
				res: &api.Response{
					Content: "test object",
				},
			},
		),
		gen(
			"nil request",
			[]string{cndNilRequest},
			[]string{actCheckResponse, actCheckError},
			&condition{
				a: api.NewContainerAPI(),
				reqs: []*api.Request{
					nil,
				},
			},
			&action{
				res: nil,
				err: &er.Error{
					Package:     api.ErrPkg,
					Type:        api.ErrTypeContainer,
					Description: api.ErrDscNil,
				},
			},
		),
		gen(
			"invalid content",
			[]string{cndWrongType},
			[]string{actCheckResponse, actCheckError},
			&condition{
				a: api.NewContainerAPI(),
				reqs: []*api.Request{
					{
						Method:  api.MethodPost,
						Format:  api.FormatJSON,
						Key:     "test key string",
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
				a: api.NewContainerAPI(),
				reqs: []*api.Request{
					{
						Method:  api.MethodPost, // Register an object.
						Key:     "test key string",
						Content: "test1 object",
					},
					{
						Method:  api.MethodPost,
						Key:     "test key string", // Register with the same key.
						Content: "test2 object",
					},
				},
			},
			&action{
				res: nil,
				err: &er.Error{
					Package:     api.ErrPkg,
					Type:        api.ErrTypeContainer,
					Description: api.ErrDscDuplicateKey,
				},
			},
		),
		gen(
			"unsupported method",
			[]string{cndUnsupportedMethod},
			[]string{actCheckResponse, actCheckError},
			&condition{
				a: api.NewContainerAPI(),
				reqs: []*api.Request{
					{
						Method:  api.Method("UNSUPPORTED METHOD"),
						Key:     "test key string",
						Content: "test object",
					},
				},
			},
			&action{
				res: nil,
				err: &er.Error{
					Package:     api.ErrPkg,
					Type:        api.ErrTypeContainer,
					Description: api.ErrDscNoMethod,
				},
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			var res *api.Response
			var err error

			ctx := context.Background()
			for _, r := range tt.C().reqs {
				res, err = tt.C().a.Serve(ctx, r)
			}

			// Check the response for the final request.
			testutil.Diff(t, tt.A().err, err, cmpopts.EquateErrors())
			testutil.Diff(t, tt.A().res, res)
		})
	}
}
