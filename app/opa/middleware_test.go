// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package opa

import (
	"context"
	"io"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/aileron-gateway/aileron-gateway/kernel/log"
	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
	utilhttp "github.com/aileron-gateway/aileron-gateway/util/http"
	"github.com/open-policy-agent/opa/rego"
)

type testClaims struct {
	User string `json:"user"`
	Age  int    `json:"age"`
}

type transaction struct{}

func (t *transaction) ID() uint64 {
	return rand.Uint64()
}

// func TestFile(t *testing.T) {

// 	module1 := `package example.authz
// 	allow {
// 		print(data.auth)
// 		print("foo****************")
// 		input.auth.user == "alice"
// 	}
// 	`

// 	store, err := newFileStore(&v1.FileStore{
// 		Path:      map[string]string{"/": "data.json"},
// 		Directory: "./",
// 	})
// 	if err != nil {
// 		panic(err)
// 	}

// 	q, err := rego.New(
// 		rego.EnablePrintStatements(true),
// 		rego.PrintHook(topdown.NewPrintHook(os.Stdout)),
// 		rego.Query("data.example.authz.allow"),
// 		rego.Module("module", module1),
// 		rego.Store(store),
// 	).PrepareForEval(context.Background())
// 	if err != nil {
// 		panic(err)
// 	}

// 	r, _ := q.Eval(context.Background(), rego.EvalInput(map[string]any{"host": "foo"}))
// 	fmt.Println(r)
// 	t.Error()
// }

// func TestServe(t *testing.T) {

// 	module1 := `package example.authz
// 	foo {false}
// 	allow {
// 		print(data.baz)
// 		print("foo****************")
// 		input.auth.user == "alice"
// 	}
// 	`
// 	_ = module1

// 	store := inmem.NewFromObject(map[string]any{"foo": "bar"})
// 	txn, err := store.NewTransaction(context.Background(), storage.WriteParams)

// 	b1, err := loader.NewFileLoader().AsBundle("../../../test/ut/app/opa/bundle2/b1/")
// 	b2, err := loader.NewFileLoader().AsBundle("../../../test/ut/app/opa/bundle2/b2/")
// 	_ = b2

// 	q, err := rego.New(
// 		rego.EnablePrintStatements(true),
// 		rego.PrintHook(topdown.NewPrintHook(os.Stdout)),
// 		rego.Query("data.example.authn.allow"),
// 		// rego.Query("data.example.authz.bar"),
// 		// rego.Query("data.example.authz.allow = allow"),
// 		// rego.Module("example1.rego", module1),
// 		rego.ParsedBundle("bbb", b1),
// 		// rego.LoadBundle("../../../test/ut/app/opa/bundle"),
// 		// rego.LoadBundle("../../../test/ut/app/opa/bundle2"),
// 		rego.Transaction(txn),
// 		rego.Store(store),
// 	).PrepareForEval(context.Background())
// 	if err != nil {
// 		panic(err)
// 	}
// 	_ = q
// 	for i := 0; i < 20; i++ {
// 		r, _ := q.Eval(context.Background(), rego.EvalInput(map[string]any{"host": "foo"}))
// 		fmt.Println(i, r)
// 		// time.Sleep(time.Second)
// 		// store.Write(context.Background(), txn, storage.AddOp, storage.MustParsePath("/"), map[string]any{"baz": "hoge"})
// 		// store.Commit(context.Background(), txn)
// 		if i == 10 {
// 			txn, _ := store.NewTransaction(context.Background(), storage.WriteParams)
// 			store.Truncate(context.Background(), txn, storage.WriteParams, nil)
// 			store.Commit(context.Background(), txn)
// 			// p1 := unsafe.Pointer(b1)
// 			// p2 := unsafe.Pointer(b2)
// 			// atomic.SwapPointer(&p1, p2)
// 		}
// 	}

// 	t.Error(q.Eval(context.Background(), rego.EvalInput(map[string]any{"host": "foo"})))
// }

func TestServeAuthz(t *testing.T) {
	module1 := `
	package example.authz
	allow {
		input.auth.user == "alice"
	}
	`
	module2 := `
	package example.authz
	allow {
		input.auth.age > 10
	}
	`
	q1, _ := rego.New(
		rego.Query("data.example.authz.allow"),
		rego.Module("example1.rego", module1),
	).PrepareForEval(context.Background())

	q2, _ := rego.New(
		rego.Dump(os.Stdout),
		rego.Query("data.example.authz.allow"),
		rego.Module("example2.rego", module2),
	).PrepareForEval(context.Background())

	eq, _ := rego.New(
		rego.Query("data.example.authz.allow"),
		rego.Module("example.rego", module1),
		rego.DisableInlining([]string{"error"}),
	).PrepareForEval(context.Background())

	type condition struct {
		r        *http.Request
		queries  []*rego.PreparedEvalQuery
		ctxValue any
	}

	type action struct {
		authorized bool
		status     int
		body       any
	}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"authorization failure with zero queries",
			&condition{
				r:       httptest.NewRequest(http.MethodGet, "https://test.com/", nil),
				queries: []*rego.PreparedEvalQuery{},
				ctxValue: testClaims{
					User: "alice",
					Age:  20,
				},
			},
			&action{
				authorized: false,
				status:     http.StatusForbidden,
				body:       `{"status":403,"statusText":"Forbidden"}`,
			},
		),
		gen(
			"authorization success with one query",
			&condition{
				r:       httptest.NewRequest(http.MethodGet, "https://test.com/", nil),
				queries: []*rego.PreparedEvalQuery{&q1},
				ctxValue: testClaims{
					User: "alice",
					Age:  20,
				},
			},
			&action{
				authorized: true,
				status:     http.StatusOK,
				body:       "",
			},
		),
		gen(
			"authorization failure with one query",
			&condition{
				r:       httptest.NewRequest(http.MethodGet, "https://test.com/", nil),
				queries: []*rego.PreparedEvalQuery{&q1},
				ctxValue: testClaims{
					User: "jon",
					Age:  20,
				},
			},
			&action{
				authorized: false,
				status:     http.StatusForbidden,
				body:       `{"status":403,"statusText":"Forbidden"}`,
			},
		),
		gen(
			"the first query authorization success with two queries",
			&condition{
				r:       httptest.NewRequest(http.MethodGet, "https://test.com/", nil),
				queries: []*rego.PreparedEvalQuery{&q1, &q2},
				ctxValue: testClaims{
					User: "alice",
					Age:  10,
				},
			},
			&action{
				authorized: true,
				status:     http.StatusOK,
				body:       "",
			},
		),
		gen(
			"the second query authorization success with two queries",
			&condition{
				r:       httptest.NewRequest(http.MethodGet, "https://test.com/", nil),
				queries: []*rego.PreparedEvalQuery{&q1, &q2},
				ctxValue: testClaims{
					User: "jon",
					Age:  20,
				},
			},
			&action{
				authorized: true,
				status:     http.StatusOK,
				body:       "",
			},
		),
		gen(
			"authorization failure with two queries",
			&condition{
				r:       httptest.NewRequest(http.MethodGet, "https://test.com/", nil),
				queries: []*rego.PreparedEvalQuery{&q1, &q2},
				ctxValue: testClaims{
					User: "jon",
					Age:  10,
				},
			},
			&action{
				authorized: false,
				status:     http.StatusForbidden,
				body:       `{"status":403,"statusText":"Forbidden"}`,
			},
		),
		gen(
			"failed to evaluation",
			&condition{
				r:        httptest.NewRequest(http.MethodGet, "https://test.com/", nil),
				queries:  []*rego.PreparedEvalQuery{&eq},
				ctxValue: testClaims{},
			},
			&action{
				authorized: false,
				status:     http.StatusForbidden,
				body:       `{"status":403,"statusText":"Forbidden"}`,
			},
		),
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			m := &authz{
				lg:      log.GlobalLogger(log.DefaultLoggerName),
				eh:      utilhttp.GlobalErrorHandler(utilhttp.DefaultErrorHandlerName),
				key:     "AuthnClaims",
				queries: tt.C.queries,
			}

			w := httptest.NewRecorder()

			r := tt.C.r
			if tt.C.ctxValue != nil {
				ctx := context.WithValue(r.Context(), m.key, tt.C.ctxValue)
				r = r.WithContext(ctx)
			}

			authorized := false
			h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				authorized = true
				w.WriteHeader(http.StatusOK)
			})

			m.Middleware(h).ServeHTTP(w, r)

			body, _ := io.ReadAll(w.Result().Body)
			testutil.Diff(t, tt.A.authorized, authorized)
			testutil.Diff(t, tt.A.status, w.Result().StatusCode)
			testutil.Diff(t, tt.A.body, string(body))
		})
	}
}
