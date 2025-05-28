// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package opa

// func TestLoadBundle(t *testing.T) {
// 	type condition struct {
// 		path   string
// 		rt     http.RoundTripper
// 		loader loader.FileLoader
// 	}

// 	type action struct {
// 		err error
// 	}

// 	tb := testutil.NewTableBuilder[*condition, *action]()
// 	tb.Name(t.Name())
// 	table := tb.Build()

// 	testData := "../../../test/ut/app/opa/"

// 	gen := testutil.NewCase[*condition, *action]
// 	testCases := []*testutil.Case[*condition, *action]{
// 		gen(
// 			"http error",
// 			[]string{},
// 			[]string{},
// 			&condition{
// 				path:   "http://tset.com/bundle.tar.gz\n",
// 				rt:     nil,
// 				loader: loader.NewFileLoader().WithSkipBundleVerification(true),
// 			},
// 			&action{
// 				err: nil,
// 			},
// 		),
// 		gen(
// 			"skip verify",
// 			[]string{},
// 			[]string{},
// 			&condition{
// 				path:   testData + "bundle.tar.gz",
// 				rt:     nil,
// 				loader: loader.NewFileLoader().WithSkipBundleVerification(true),
// 			},
// 			&action{
// 				err: nil,
// 			},
// 		),
// 		gen(
// 			"load failed",
// 			[]string{},
// 			[]string{},
// 			&condition{
// 				path:   testData + "bundle.tar.gz",
// 				rt:     nil,
// 				loader: loader.NewFileLoader(),
// 			},
// 			&action{
// 				err: &er.Error{
// 					Package:     "authz/opa",
// 					Type:        "load bundle",
// 					Description: "failed to load bundle.",
// 				},
// 			},
// 		),
// 	}

// 	testutil.Register(table, testCases...)

// 	for _, tt := range table.Entries() {
// 		tt := tt
// 		t.Run(tt.Name(), func(t *testing.T) {

// 			b, err := loadBundle(tt.C().path, tt.C().rt, tt.C().loader)
// 			testutil.Diff(t, tt.A().err, err, cmpopts.EquateErrors())
// 			if err != nil {
// 				testutil.Diff(t, (*bundle.Bundle)(nil), b)
// 				return
// 			}

// 			opts := []func(*rego.Rego){
// 				rego.Query("data.authz.allow"),
// 				rego.ParsedBundle("test", b),
// 			}
// 			query, err := rego.New(opts...).PrepareForEval(context.TODO())
// 			testutil.Diff(t, nil, err)
// 			result, err := query.Eval(context.TODO(), rego.EvalInput(map[string]any{"user": "bob"}))
// 			testutil.Diff(t, nil, err)
// 			testutil.Diff(t, true, result.Allowed())

// 		})
// 	}
// }
