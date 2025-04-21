// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package casbin

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/aileron-gateway/aileron-gateway/kernel/log"
	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
	utilhttp "github.com/aileron-gateway/aileron-gateway/util/http"
	"github.com/casbin/casbin/v3"
)

func TestMiddleware(t *testing.T) {
	rootDir := "../../../test/ut/app/casbin/middleware/"

	// Enforcer.
	enf1, _ := casbin.NewEnforcer(rootDir+"abac_model.conf", rootDir+"abac_policy1.csv")
	enf1.AddFunction("authValue", mapValue)

	enf2, _ := casbin.NewEnforcer(rootDir+"abac_model.conf", rootDir+"abac_policy2.csv")
	enf2.AddFunction("authValue", mapValue)
	enf2.AddFunction("containsString", contains[string])

	enf3, _ := casbin.NewEnforcer(rootDir+"abac_model.conf", rootDir+"abac_policy3.csv")
	enf3.AddFunction("authValue", mapValue)
	enf3.AddFunction("containsInt", containsNumber[int])

	enf4, _ := casbin.NewEnforcer(rootDir+"abac_model.conf", rootDir+"abac_policy4.csv")
	enf4.AddFunction("authValue", mapValue)
	enf4.AddFunction("containsInt", containsNumber[int])
	enf4.AddFunction("asIntSlice", asSliceNumber[int])

	enf4sub, _ := casbin.NewEnforcer(rootDir+"abac_model_4sub.conf", rootDir+"abac_policy_sub4.csv")

	enf5sub, _ := casbin.NewEnforcer(rootDir+"abac_model_5sub.conf", rootDir+"abac_policy_sub5.csv")
	enf5sub.AddFunction("mapValue", mapValue)

	// RBAC Enforcer
	enf6, _ := casbin.NewEnforcer(rootDir+"rbac_model.conf", rootDir+"rbac_policy1.csv")
	enf6.AddFunction("mapValue", mapValue)

	// ACL Enforcer
	enf7, _ := casbin.NewEnforcer(rootDir+"acl_model.conf", "")
	enf7.AddFunction("mapValue", mapValue)

	type condition struct {
		r       *http.Request
		cv      any
		enf     []casbin.IEnforcer
		keys    []string
		explain bool
	}

	type action struct {
		authorized bool
		status     int
		body       any
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	condZeroQueries := tb.Condition("zero queries", "zero queries")
	condOneQuery := tb.Condition("one query", "one query")
	condTwoQueries := tb.Condition("two queries", "two queries")
	condExtraKey := tb.Condition("extra key", "set extra keys")
	condSuccessFirstQuery := tb.Condition("success of the first query", "success of the first query")
	condSuccessSecondQuery := tb.Condition("success of the second query", "success of the second query")
	actOk := tb.Action("authorization success", "authorization success")
	actNg := tb.Action("authorization failure", "authorization failure")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"failure/0 enforcer",
			[]string{condZeroQueries},
			[]string{actNg},
			&condition{
				r: httptest.NewRequest(http.MethodGet, "https://test.com/", nil),
				cv: map[string]any{
					"age":  20,
					"role": []string{"role1", "role2"},
				},
				enf: []casbin.IEnforcer{},
			},
			&action{
				authorized: false,
				status:     http.StatusForbidden,
				body:       `{"status":403,"statusText":"Forbidden"}`,
			},
		),
		gen(
			"success/1 enforcer",
			[]string{condOneQuery, condSuccessFirstQuery},
			[]string{actOk},
			&condition{
				r: httptest.NewRequest(http.MethodGet, "https://test.com/", nil),
				cv: map[string]any{
					"age":  20,
					"role": []string{"role1", "role2"},
				},
				enf: []casbin.IEnforcer{enf1},
			},
			&action{
				authorized: true,
				status:     http.StatusOK,
				body:       "ok",
			},
		),
		gen(
			"failure/1 enforcer",
			[]string{condOneQuery},
			[]string{actNg},
			&condition{
				r: httptest.NewRequest(http.MethodGet, "https://test.com/", nil),
				cv: map[string]any{
					"age":  10,
					"role": []string{"role1", "role2"},
				},
				enf: []casbin.IEnforcer{enf1},
			},
			&action{
				authorized: false,
				status:     http.StatusForbidden,
				body:       `{"status":403,"statusText":"Forbidden"}`,
			},
		),
		gen(
			"failure/use extra key for r=sub,obj,act definition", // Model must have 4 keys.
			[]string{condOneQuery, condSuccessFirstQuery},
			[]string{actOk},
			&condition{
				r: httptest.NewRequest(http.MethodGet, "https://test.com/", nil),
				cv: map[string]any{
					"age":  20,
					"role": []string{"role1", "role2"},
				},
				enf:  []casbin.IEnforcer{enf1},
				keys: []string{"ext_string"},
			},
			&action{
				authorized: false,
				status:     http.StatusForbidden,
				body:       `{"status":403,"statusText":"Forbidden"}`,
			},
		),
		gen(
			"success/1 enforcer with extra string key",
			[]string{condOneQuery, condSuccessFirstQuery, condExtraKey},
			[]string{actOk},
			&condition{
				r:    httptest.NewRequest(http.MethodPost, "https://test.com/", nil),
				cv:   map[string]any{},
				enf:  []casbin.IEnforcer{enf4sub},
				keys: []string{"ext_string"},
			},
			&action{
				authorized: true,
				status:     http.StatusOK,
				body:       "ok",
			},
		),
		gen(
			"failure/1 enforcer with extra string key",
			[]string{condOneQuery, condExtraKey},
			[]string{actNg},
			&condition{
				r:    httptest.NewRequest(http.MethodGet, "https://test.com/", nil),
				cv:   map[string]any{},
				enf:  []casbin.IEnforcer{enf4sub},
				keys: []string{"ext_string"},
			},
			&action{
				authorized: false,
				status:     http.StatusForbidden,
				body:       `{"status":403,"statusText":"Forbidden"}`,
			},
		),
		gen(
			"success/1 enforcer with extra map key",
			[]string{condOneQuery, condSuccessFirstQuery, condExtraKey},
			[]string{actOk},
			&condition{
				r:    httptest.NewRequest(http.MethodGet, "https://test.com/", nil),
				cv:   map[string]any{},
				enf:  []casbin.IEnforcer{enf5sub},
				keys: []string{"ext_string", "ext_map1"},
			},
			&action{

				authorized: true,
				status:     http.StatusOK,
				body:       "ok",
			},
		),
		gen(
			"success/1 enforcer with extra map key",
			[]string{condOneQuery, condSuccessFirstQuery, condExtraKey},
			[]string{actOk},
			&condition{
				r:    httptest.NewRequest(http.MethodGet, "https://test.com/", nil),
				cv:   map[string]any{},
				enf:  []casbin.IEnforcer{enf5sub},
				keys: []string{"ext_unused", "ext_map1"},
			},
			&action{
				authorized: true, // Model and policy do not used 4th input. It can be nil.
				status:     http.StatusOK,
				body:       "ok",
			},
		),
		gen(
			"failure/1 enforcer with extra map key",
			[]string{condOneQuery, condSuccessFirstQuery, condExtraKey},
			[]string{actNg},
			&condition{
				r:    httptest.NewRequest(http.MethodGet, "https://test.com/", nil),
				cv:   map[string]any{},
				enf:  []casbin.IEnforcer{enf5sub},
				keys: []string{"ext_map1", "ext_map2"},
			},
			&action{
				authorized: false,
				status:     http.StatusForbidden,
				body:       `{"status":403,"statusText":"Forbidden"}`,
			},
		),
		gen(
			"success/1st of 2 enforcer",
			[]string{condTwoQueries, condSuccessFirstQuery},
			[]string{actOk},
			&condition{
				r: httptest.NewRequest(http.MethodGet, "https://test.com/", nil),
				cv: map[string]any{
					"age":  20,
					"role": []string{"role2"},
				},
				enf: []casbin.IEnforcer{enf1, enf2},
			},
			&action{
				authorized: true,
				status:     http.StatusOK,
				body:       "ok",
			},
		),
		gen(
			"success/2nd of 2 enforcers",
			[]string{condTwoQueries, condSuccessSecondQuery},
			[]string{actOk},
			&condition{
				r: httptest.NewRequest(http.MethodGet, "https://test.com/", nil),
				cv: map[string]any{
					"age":  10,
					"role": []string{"role1", "role2"},
				},
				enf: []casbin.IEnforcer{enf1, enf2},
			},
			&action{
				authorized: true,
				status:     http.StatusOK,
				body:       "ok",
			},
		),
		gen(
			"failure/2 enforcers",
			[]string{condTwoQueries},
			[]string{actNg},
			&condition{
				r: httptest.NewRequest(http.MethodGet, "https://test.com/", nil),
				cv: map[string]any{
					"age":  10,
					"role": []string{"role2"},
				},
				enf: []casbin.IEnforcer{enf1, enf2},
			},
			&action{
				authorized: false,
				status:     http.StatusForbidden,
				body:       `{"status":403,"statusText":"Forbidden"}`,
			},
		),
		gen(
			"authorization success with explain enforcer",
			[]string{},
			[]string{actOk},
			&condition{
				r: httptest.NewRequest(http.MethodGet, "https://test.com/", nil),
				cv: map[string]any{
					"age": 20,
				},
				enf: []casbin.IEnforcer{enf1},
			},
			&action{
				authorized: true,
				status:     http.StatusOK,
				body:       "ok",
			},
		),
		gen(
			"invalid auth type",
			[]string{},
			[]string{actNg},
			&condition{
				r:   httptest.NewRequest(http.MethodGet, "https://test.com/", nil),
				enf: []casbin.IEnforcer{enf1},
				cv:  "invalid", // map[string]any expected.
			},
			&action{
				authorized: false,
				status:     http.StatusForbidden,
				body:       `{"status":403,"statusText":"Forbidden"}`,
			},
		),
		gen(
			"success/1 enforcer/with int slice",
			[]string{condOneQuery, condSuccessFirstQuery},
			[]string{actOk},
			&condition{
				r: httptest.NewRequest(http.MethodGet, "https://test.com/", nil),
				cv: map[string]any{
					"age":  25,
					"role": []int{123, 456},
				},
				enf: []casbin.IEnforcer{enf3},
			},
			&action{
				authorized: true,
				status:     http.StatusOK,
				body:       "ok",
			},
		),
		gen(
			"success/1 enforcer/with int",
			[]string{condOneQuery, condSuccessFirstQuery},
			[]string{actOk},
			&condition{
				r: httptest.NewRequest(http.MethodGet, "https://test.com/", nil),
				cv: map[string]any{
					"age":   25,
					"role1": 123,
					"role2": 456,
				},
				enf: []casbin.IEnforcer{enf4},
			},
			&action{
				authorized: true,
				status:     http.StatusOK,
				body:       "ok",
			},
		),
		gen(
			"success/rbac by name",
			[]string{condOneQuery},
			[]string{actOk},
			&condition{
				r:       httptest.NewRequest(http.MethodGet, "https://test.com/alice/data", nil),
				cv:      map[string]any{"name": "alice"}, // alice can GET /alice/*
				enf:     []casbin.IEnforcer{enf6},
				explain: true,
			},
			&action{
				authorized: true,
				status:     http.StatusOK,
				body:       "ok",
			},
		),
		gen(
			"failure/rbac by name",
			[]string{condOneQuery},
			[]string{actNg},
			&condition{
				r:       httptest.NewRequest(http.MethodPost, "https://test.com/alice/data", nil),
				cv:      map[string]any{"name": "alice"}, // alice cannot POST /alice/*
				enf:     []casbin.IEnforcer{enf6},
				explain: true,
			},
			&action{
				authorized: false,
				status:     http.StatusForbidden,
				body:       `{"status":403,"statusText":"Forbidden"}`,
			},
		),
		gen(
			"success/rbac by role",
			[]string{condOneQuery},
			[]string{actOk},
			&condition{
				r:       httptest.NewRequest(http.MethodPost, "https://test.com/foo/data", nil),
				cv:      map[string]any{"name": "alice"}, // admin alice can POST /foo/*
				enf:     []casbin.IEnforcer{enf6},
				explain: true,
			},
			&action{
				authorized: true,
				status:     http.StatusOK,
				body:       "ok",
			},
		),
		gen(
			"failure/rbac by role",
			[]string{condOneQuery},
			[]string{actNg},
			&condition{
				r:       httptest.NewRequest(http.MethodPost, "https://test.com/foo/data", nil),
				cv:      map[string]any{"name": "bob"}, // non admin bob cannot POST /foo/*
				enf:     []casbin.IEnforcer{enf6},
				explain: true,
			},
			&action{
				authorized: false,
				status:     http.StatusForbidden,
				body:       `{"status":403,"statusText":"Forbidden"}`,
			},
		),
		gen(
			"success/acl",
			[]string{condOneQuery},
			[]string{actOk},
			&condition{
				r:   httptest.NewRequest(http.MethodPost, "https://test.com/", nil),
				cv:  map[string]any{"name": "alice", "boss": "alice"}, // name==boss is allowed.
				enf: []casbin.IEnforcer{enf7},
			},
			&action{
				authorized: true,
				status:     http.StatusOK,
				body:       "ok",
			},
		),
		gen(
			"failure/acl",
			[]string{condOneQuery},
			[]string{actOk},
			&condition{
				r:   httptest.NewRequest(http.MethodPost, "https://test.com/", nil),
				cv:  map[string]any{"name": "alice", "boss": "bob"}, // name!=boss is not allowed.
				enf: []casbin.IEnforcer{enf7},
			},
			&action{
				authorized: false,
				status:     http.StatusForbidden,
				body:       `{"status":403,"statusText":"Forbidden"}`,
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			lg := log.NewJSONSLogger(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})
			m := &authz{
				lg:        lg,
				w:         lg,
				eh:        utilhttp.GlobalErrorHandler(utilhttp.DefaultErrorHandlerName),
				key:       "AuthnClaims",
				extraKeys: tt.C().keys,
				enforcers: tt.C().enf,
				explain:   tt.C().explain,
			}

			authorized := false
			h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				authorized = true
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("ok"))
			})

			r := tt.C().r
			if tt.C().cv != nil {
				ctx := context.WithValue(r.Context(), m.key, tt.C().cv)
				ctx = context.WithValue(ctx, "ext_string", "foo")
				ctx = context.WithValue(ctx, "ext_map1", map[string]any{"foo": "bar"})
				ctx = context.WithValue(ctx, "ext_map2", map[string]any{"foo": "123"})
				r = r.WithContext(ctx)
			}
			w := httptest.NewRecorder()
			m.Middleware(h).ServeHTTP(w, r)

			body, _ := io.ReadAll(w.Result().Body)
			testutil.Diff(t, tt.A().authorized, authorized)
			testutil.Diff(t, tt.A().status, w.Result().StatusCode)
			testutil.Diff(t, tt.A().body, string(body))
		})
	}
}
