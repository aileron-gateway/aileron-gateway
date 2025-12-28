// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package http

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	v1 "github.com/aileron-gateway/aileron-gateway/apis/core/v1"
	"github.com/aileron-gateway/aileron-gateway/internal/testutil"
	"github.com/google/go-cmp/cmp/cmpopts"
)

type testCookieCreator struct {
}

func (c *testCookieCreator) NewCookie() *http.Cookie {
	return &http.Cookie{}
}

func TestGetCookie(t *testing.T) {
	type condition struct {
		cookies map[string]string
		key     string
	}

	type action struct {
		value string
	}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"no cookie", &condition{
				cookies: map[string]string{},
				key:     "test",
			},
			&action{
				value: "",
			},
		),
		gen(
			"dummy cookies with key", &condition{
				cookies: map[string]string{
					"test": "dummy",
					"foo":  "dummy",
					"bar":  "dummy",
				},
				key: "test",
			},
			&action{
				value: "",
			},
		),
		gen(
			"dummy cookies with empty key", &condition{
				cookies: map[string]string{
					"test": "dummy",
					"foo":  "dummy",
					"bar":  "dummy",
				},
				key: "",
			},
			&action{
				value: "",
			},
		),
		gen(
			"1 cookie", &condition{
				cookies: map[string]string{
					"test0": "xxx",
					"test":  "dummy",
					"foo":   "dummy",
				},
				key: "test",
			},
			&action{
				value: "xxx",
			},
		),
		gen(
			"2 cookie, successive index", &condition{
				cookies: map[string]string{
					"test0": "xxx",
					"test1": "yyy",
					"test":  "dummy",
					"foo":   "dummy",
				},
				key: "test",
			},
			&action{
				value: "xxxyyy",
			},
		),
		gen(
			"2 cookie, jumped index", &condition{
				cookies: map[string]string{
					"test0": "xxx",
					"test3": "yyy",
					"test":  "dummy",
					"foo":   "dummy",
				},
				key: "test",
			},
			&action{
				value: "xxxyyy",
			},
		),
		gen(
			"2 cookie, too big index", &condition{
				cookies: map[string]string{
					"test0":  "xxx",
					"test99": "dummy", // Index must be [0,len(cookie)-1].
					"test":   "dummy",
					"foo":    "dummy",
				},
				key: "test",
			},
			&action{
				value: "xxx",
			},
		),
		gen(
			"5 cookies", &condition{
				cookies: map[string]string{
					"test0":  "xxx",
					"test1":  "yyy",
					"test3":  "zzz",
					"test":   "dummy",
					"test99": "dummy",
					"foo":    "dummy",
					"bar":    "dummy",
				},
				key: "test",
			},
			&action{
				value: "xxxyyyzzz",
			},
		),
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			var cks []*http.Cookie
			for k, v := range tt.C.cookies {
				cks = append(cks, &http.Cookie{
					Name:  k,
					Value: v,
				})
			}
			value := GetCookie(tt.C.key, cks)
			testutil.Diff(t, tt.A.value, value)
		})
	}
}

func TestSetCookie(t *testing.T) {
	type condition struct {
		names  []string
		prefix string
		length int
	}

	type action struct {
		keys    []string
		deletes []string
	}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"empty prefix", &condition{
				prefix: "",
				length: 0,
			},
			&action{
				keys:    []string{},
				deletes: []string{},
			},
		),
		gen(
			"1 cookie", &condition{
				prefix: "test",
				length: 100,
			},
			&action{
				keys:    []string{"test0"},
				deletes: []string{},
			},
		),
		gen(
			"1 cookie with max size", &condition{
				prefix: "test",
				length: (1<<12 - 1<<7) / 2,
			},
			&action{
				keys:    []string{"test0"},
				deletes: []string{},
			},
		),
		gen(
			"2 cookies", &condition{
				prefix: "test",
				length: (1<<12 - 1<<7),
			},
			&action{
				keys:    []string{"test0", "test1"},
				deletes: []string{},
			},
		),
		gen(
			"4 cookies", &condition{
				prefix: "test",
				length: (1<<12 - 1<<7) * 2,
			},
			&action{
				keys:    []string{"test0", "test1", "test2", "test3"},
				deletes: []string{},
			},
		),
		gen(
			"remove unnecessary cookies", &condition{
				names:  []string{"test0", "test1", "test2", "test3"},
				prefix: "test",
				length: (1<<12 - 1<<7),
			},
			&action{
				keys:    []string{"test0", "test1"},
				deletes: []string{"test2", "test3"},
			},
		),
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			val := make([]byte, tt.C.length)
			rand.Read(val)
			value := hex.EncodeToString(val)

			w := httptest.NewRecorder()
			cc := &testCookieCreator{}
			SetCookie(w, tt.C.names, cc, tt.C.prefix, value)

			values := make([]string, len(tt.A.keys))
			deletes := []string{}
			for _, v := range w.Result().Header["Set-Cookie"] {
				for i, k := range tt.A.keys {
					if strings.HasPrefix(v, k+"=") {
						values[i] = strings.TrimPrefix(v, k+"=")
					}
				}
				if strings.HasSuffix(v, "=; Max-Age=0") {
					deletes = append(deletes, strings.TrimSuffix(v, "=; Max-Age=0"))
				}
			}

			result := strings.Join(values, "")
			testutil.Diff(t, result, value)
			testutil.Diff(t, tt.A.deletes, deletes, cmpopts.SortSlices(func(a, b string) bool { return a < b }))
		})
	}
}

func TestDeleteCookie(t *testing.T) {
	type condition struct {
		names  []string
		prefix string
	}

	type action struct {
		keys []string
	}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"empty prefix", &condition{
				prefix: "",
				names:  []string{},
			},
			&action{
				keys: []string{},
			},
		),
		gen(
			"no cookie", &condition{
				prefix: "test",
				names:  []string{},
			},
			&action{
				keys: []string{},
			},
		),
		gen(
			"irrelevant cookie", &condition{
				prefix: "test",
				names:  []string{"test", "foo"},
			},
			&action{
				keys: []string{},
			},
		),
		gen(
			"delete cookie", &condition{
				prefix: "test",
				names:  []string{"test", "test0", "test1", "test2", "foo"},
			},
			&action{
				keys: []string{"test0", "test1", "test2"},
			},
		),
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			w := httptest.NewRecorder()
			DeleteCookie(w, tt.C.names, tt.C.prefix)

			deletes := []string{}
			for _, v := range w.Result().Header["Set-Cookie"] {
				if strings.HasSuffix(v, "=; Max-Age=0") {
					deletes = append(deletes, strings.TrimSuffix(v, "=; Max-Age=0"))
				}
			}
			testutil.Diff(t, tt.A.keys, deletes, cmpopts.SortSlices(func(a, b string) bool { return a < b }))
		})
	}
}

func TestCookieNames(t *testing.T) {
	type condition struct {
		cookies map[string]string
	}

	type action struct {
		names []string
	}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"no cookie", &condition{
				cookies: map[string]string{},
			},
			&action{
				names: []string{},
			},
		),
		gen(
			"cookies", &condition{
				cookies: map[string]string{
					"foo": "value",
					"bar": "value",
				},
			},
			&action{
				names: []string{"foo", "bar"},
			},
		),
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			var cks []*http.Cookie
			for k, v := range tt.C.cookies {
				cks = append(cks, &http.Cookie{
					Name:  k,
					Value: v,
				})
			}
			names := CookieNames(cks)
			testutil.Diff(t, tt.A.names, names, cmpopts.SortSlices(func(a, b string) bool { return a < b }))
		})
	}
}

func TestMinInt(t *testing.T) {
	type condition struct {
		x int
		y int
	}

	type action struct {
		result int
	}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"same", &condition{
				x: 1,
				y: 1,
			},
			&action{
				result: 1,
			},
		),
		gen(
			"x smaller than y", &condition{
				x: -1,
				y: 1,
			},
			&action{
				result: -1,
			},
		),
		gen(
			"y smaller than x", &condition{
				x: 1,
				y: -1,
			},
			&action{
				result: -1,
			},
		),
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			result := minInt(tt.C.x, tt.C.y)
			testutil.Diff(t, tt.A.result, result)
		})
	}
}

func TestDefaultCookieCreator(t *testing.T) {
	type condition struct {
	}

	type action struct {
		ck *http.Cookie
	}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"default", &condition{},
			&action{
				ck: &http.Cookie{
					Path:     "/",
					Secure:   true,
					HttpOnly: true,
					SameSite: http.SameSiteDefaultMode,
				},
			},
		),
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			ck := DefaultCookieCreator()
			testutil.Diff(t, tt.A.ck, ck.NewCookie())
		})
	}
}

func TestNewCookieCreator(t *testing.T) {
	type condition struct {
		spec *v1.CookieSpec
	}

	type action struct {
		ck *http.Cookie
	}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"nil spec", &condition{
				spec: nil,
			},
			&action{
				ck: &http.Cookie{
					Path:     "/",
					Secure:   true,
					HttpOnly: true,
					SameSite: http.SameSiteDefaultMode,
				},
			},
		),
		gen(
			"empty spec", &condition{
				spec: &v1.CookieSpec{},
			},
			&action{
				ck: &http.Cookie{},
			},
		),
		gen(
			"full spec", &condition{
				spec: &v1.CookieSpec{
					Name:      "test",
					Value:     "value",
					Path:      "/",
					Domain:    "test.com",
					ExpiresIn: 10,
					MaxAge:    10,
					Secure:    true,
					HTTPOnly:  true,
					SameSite:  v1.SameSite_Default,
				},
			},
			&action{
				ck: &http.Cookie{
					Name:     "test",
					Value:    "value",
					Path:     "/",
					Domain:   "test.com",
					Expires:  time.Now().Add(10 * time.Second),
					MaxAge:   10,
					Secure:   true,
					HttpOnly: true,
					SameSite: http.SameSiteDefaultMode,
				},
			},
		),
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			ck := NewCookieCreator(tt.C.spec)
			testutil.Diff(t, tt.A.ck, ck.NewCookie(), cmpopts.EquateApproxTime(time.Second))
		})
	}
}

func TestSameSite(t *testing.T) {
	type condition struct {
		in v1.SameSite
	}

	type action struct {
		val http.SameSite
	}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"undefined", &condition{
				in: 999,
			},
			&action{
				val: http.SameSite(0),
			},
		),
		gen(
			"default", &condition{
				in: v1.SameSite_Default,
			},
			&action{
				val: http.SameSiteDefaultMode,
			},
		),
		gen(
			"lax", &condition{
				in: v1.SameSite_Lax,
			},
			&action{
				val: http.SameSiteLaxMode,
			},
		),
		gen(
			"strict", &condition{
				in: v1.SameSite_Strict,
			},
			&action{
				val: http.SameSiteStrictMode,
			},
		),
		gen(
			"none", &condition{
				in: v1.SameSite_None,
			},
			&action{
				val: http.SameSiteNoneMode,
			},
		),
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			val := sameSite(tt.C.in)
			testutil.Diff(t, tt.A.val, val)
		})
	}
}
