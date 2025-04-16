// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package api

import (
	"bytes"
	"context"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/aileron-gateway/aileron-gateway/kernel/er"
	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestPrintDebug(t *testing.T) {
	type condition struct {
		level    int
		args     []any
		debugLv1 bool
		debugLv2 bool
		debugLv3 bool
	}

	type action struct {
		contains    string
		notContains string
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"print lv3",
			[]string{},
			[]string{},
			&condition{
				level:    debugLv3,
				args:     []any{"foo", "bar"},
				debugLv3: true,
				debugLv2: true,
				debugLv1: true,
			},
			&action{
				contains:    "foo bar",
				notContains: "-",
			},
		),
		gen(
			"not print lv3",
			[]string{},
			[]string{},
			&condition{
				level:    debugLv3,
				args:     []any{"foo", "bar"},
				debugLv3: false,
				debugLv2: true,
				debugLv1: true,
			},
			&action{
				notContains: "foo bar",
			},
		),
		gen(
			"print lv2",
			[]string{},
			[]string{},
			&condition{
				level:    debugLv2,
				args:     []any{"foo", "bar"},
				debugLv3: false,
				debugLv2: true,
				debugLv1: true,
			},
			&action{
				contains:    "foo bar",
				notContains: "-",
			},
		),
		gen(
			"not print lv2",
			[]string{},
			[]string{},
			&condition{
				level:    debugLv2,
				args:     []any{"foo", "bar"},
				debugLv3: false,
				debugLv2: false,
				debugLv1: true,
			},
			&action{
				notContains: "foo bar",
			},
		),
		gen(
			"print lv1",
			[]string{},
			[]string{},
			&condition{
				level:    debugLv1,
				args:     []any{"foo", "bar"},
				debugLv3: false,
				debugLv2: false,
				debugLv1: true,
			},
			&action{
				contains:    "foo bar",
				notContains: "-",
			},
		),
		gen(
			"not print lv1",
			[]string{},
			[]string{},
			&condition{
				level:    debugLv1,
				args:     []any{"foo", "bar"},
				debugLv3: false,
				debugLv2: false,
				debugLv1: false,
			},
			&action{
				notContains: "foo bar",
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			DebugLv3 = tt.C().debugLv3
			DebugLv2 = tt.C().debugLv2
			DebugLv1 = tt.C().debugLv1

			var buf bytes.Buffer
			log.SetOutput(&buf)
			defer func() { log.SetOutput(os.Stdout) }()

			printDebug(tt.C().level, tt.C().args...)
			t.Log(buf.String())
			testutil.Diff(t, true, strings.Contains(buf.String(), tt.A().contains))
			testutil.Diff(t, false, strings.Contains(buf.String(), tt.A().notContains))
		})
	}
}

func TestNewDefaultServeMux(t *testing.T) {
	type condition struct {
	}

	type action struct {
		mux *DefaultServeMux
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
			"new serve mux",
			[]string{cndNewDefault},
			[]string{actCheckInitialized},
			&condition{},
			&action{
				mux: &DefaultServeMux{
					keys: nil,
					apis: map[string]API[*Request, *Response]{},
				},
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			mux := NewDefaultServeMux()
			testutil.Diff(t, tt.A().mux, mux, cmp.AllowUnexported(DefaultServeMux{}))
		})
	}
}

type stringResponderAPI struct {
	response string
	err      error
}

func (a *stringResponderAPI) Serve(_ context.Context, req *Request) (*Response, error) {
	return &Response{Content: a.response}, a.err
}

func TestDefaultServeMux_Serve(t *testing.T) {
	type condition struct {
		keys []string
		apis []API[*Request, *Response]
		req  *Request
	}

	type action struct {
		res *Response
		err error
	}

	cndNoExactKey := "no exact match key"
	cndNoRoute := "no API route"
	actCheckResponse := "check keys are sorted"
	actCheckNoError := "check no error"
	actCheckError := "check error"

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	tb.Condition(cndNoExactKey, "there is no keys which exactly match the request")
	tb.Condition(cndNoRoute, "there is no APIs to route to")
	tb.Action(actCheckResponse, "check the responded content is the one expected")
	tb.Action(actCheckNoError, "check that no error occurred")
	tb.Action(actCheckError, "check that the expected error was returned")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"exact match",
			[]string{},
			[]string{actCheckResponse, actCheckNoError},
			&condition{
				keys: []string{"test1", "test2"},
				apis: []API[*Request, *Response]{
					&stringResponderAPI{response: "test1"},
					&stringResponderAPI{response: "test2"},
				},
				req: &Request{Key: "test1"},
			},
			&action{
				res: &Response{Content: "test1"},
			},
		),
		gen(
			"exact match",
			[]string{},
			[]string{actCheckResponse, actCheckNoError},
			&condition{
				keys: []string{"test1", "test2"},
				apis: []API[*Request, *Response]{
					&stringResponderAPI{response: "test1"},
					&stringResponderAPI{response: "test2"},
				},
				req: &Request{Key: "test2"},
			},
			&action{
				res: &Response{Content: "test2"},
			},
		),
		gen(
			"longest match",
			[]string{cndNoExactKey},
			[]string{actCheckResponse, actCheckNoError},
			&condition{
				keys: []string{"foo", "foobar"},
				apis: []API[*Request, *Response]{
					&stringResponderAPI{response: "foo"},
					&stringResponderAPI{response: "foobar"},
				},
				req: &Request{Key: "foobarbaz"}, // Should be routed to "foobar"
			},
			&action{
				res: &Response{Content: "foobar"}, // Longest match.
			},
		),
		gen(
			"nil request",
			[]string{cndNoExactKey, cndNoRoute},
			[]string{actCheckResponse, actCheckError},
			&condition{
				keys: []string{},
				apis: []API[*Request, *Response]{},
				req:  nil,
			},
			&action{
				res: nil,
				err: &er.Error{
					Package:     ErrPkg,
					Type:        ErrTypeUtil,
					Description: ErrDscNil,
				},
			},
		),
		gen(
			"no route",
			[]string{cndNoExactKey, cndNoRoute},
			[]string{actCheckResponse, actCheckError},
			&condition{
				keys: []string{},
				apis: []API[*Request, *Response]{},
				req:  &Request{Key: "no route"},
			},
			&action{
				res: nil,
				err: &er.Error{
					Package:     ErrPkg,
					Type:        ErrTypeUtil,
					Description: ErrDscNoAPI,
				},
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			a := NewDefaultServeMux()

			for i := range tt.C().keys {
				a.Handle(tt.C().keys[i], tt.C().apis[i]) // Register APIs.
			}
			for _, k := range tt.C().keys {
				res, _ := a.Serve(context.Background(), &Request{Key: k}) // Check registered APIs.
				testutil.Diff(t, k, res.Content)
			}

			ctx := context.Background()
			res, err := a.Serve(ctx, tt.C().req)
			testutil.Diff(t, tt.A().res, res, cmp.AllowUnexported(Response{}))
			testutil.Diff(t, tt.A().err, err, cmpopts.EquateErrors())
		})
	}
}

func TestDefaultServeMux_Handle(t *testing.T) {
	type condition struct {
		keys []string
		apis []API[*Request, *Response]
	}

	type action struct {
		keys []string
		err  error
	}

	cndMultipleAPI := "multiple api"
	cndNilAPI := "nil api"
	cndKeyDuplicate := "key duplicate"
	actCheckAPIRegistered := "check API registered"
	actCheckSorted := "check keys are sorted"
	actCheckNoError := "check no error"
	actCheckError := "check error"

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	tb.Condition(cndMultipleAPI, "register multiple APIs")
	tb.Condition(cndNilAPI, "input nil API")
	tb.Condition(cndKeyDuplicate, "register with the same key")
	tb.Action(actCheckAPIRegistered, "check that APIs are successfully registered")
	tb.Action(actCheckSorted, "check that the registered keys are sorted")
	tb.Action(actCheckNoError, "check that no error occurred")
	tb.Action(actCheckError, "check that the expected error was returned")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"1 api",
			[]string{},
			[]string{actCheckAPIRegistered, actCheckNoError},
			&condition{
				keys: []string{"test"},
				apis: []API[*Request, *Response]{
					&stringResponderAPI{response: "test"},
				},
			},
			&action{
				keys: []string{"test"},
			},
		),
		gen(
			"2 api",
			[]string{cndMultipleAPI},
			[]string{actCheckSorted, actCheckAPIRegistered, actCheckNoError},
			&condition{
				keys: []string{"test1", "test2"},
				apis: []API[*Request, *Response]{
					&stringResponderAPI{response: "test1"},
					&stringResponderAPI{response: "test2"},
				},
			},
			&action{
				keys: []string{"test2", "test1"}, // keys should be sorted.
			},
		),
		gen(
			"3 api",
			[]string{cndMultipleAPI},
			[]string{actCheckSorted, actCheckAPIRegistered, actCheckNoError},
			&condition{
				keys: []string{"foo", "foobar", "foobarbaz"},
				apis: []API[*Request, *Response]{
					&stringResponderAPI{response: "foo"},
					&stringResponderAPI{response: "foobar"},
					&stringResponderAPI{response: "foobarbaz"},
				},
			},
			&action{
				keys: []string{"foobarbaz", "foobar", "foo"}, // keys should be sorted.
			},
		),
		gen(
			"nil api",
			[]string{cndNilAPI},
			[]string{actCheckAPIRegistered, actCheckNoError},
			&condition{
				keys: []string{"test"},
				apis: []API[*Request, *Response]{nil},
			},
			&action{
				keys: nil,
			},
		),
		gen(
			"duplicate key",
			[]string{cndKeyDuplicate},
			[]string{actCheckAPIRegistered, actCheckError},
			&condition{
				keys: []string{"test", "test"},
				apis: []API[*Request, *Response]{
					&stringResponderAPI{response: "test"},
					&stringResponderAPI{response: "test"},
				},
			},
			&action{
				keys: []string{"test"},
				err: &er.Error{
					Package:     ErrPkg,
					Type:        ErrTypeUtil,
					Description: ErrDscDuplicateKey,
				},
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			a := NewDefaultServeMux()

			var err error
			for i := range tt.C().keys {
				if err = a.Handle(tt.C().keys[i], tt.C().apis[i]); err != nil {
					break
				}
			}
			testutil.Diff(t, tt.A().err, err, cmpopts.EquateErrors())

			testutil.Diff(t, tt.A().keys, a.keys)
			for _, k := range tt.A().keys {
				res, _ := a.apis[k].Serve(context.Background(), nil)
				testutil.Diff(t, k, res.Content)
			}
		})
	}
}
