// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package resilience

import (
	"net/http"
	"net/http/httptest"
	"testing"

	v1 "github.com/aileron-gateway/aileron-gateway/apis/core/v1"
	"github.com/aileron-gateway/aileron-gateway/kernel/hash"
	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
	"github.com/google/go-cmp/cmp"
)

func TestNewHashers(t *testing.T) {
	type condition struct {
		specs []*v1.HTTPHasherSpec
	}

	type action struct {
		hs []HTTPHasher
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"nil spec",
			[]string{},
			[]string{},
			&condition{
				specs: nil,
			},
			&action{
				hs: []HTTPHasher{},
			},
		),
		gen(
			"empty spec",
			[]string{},
			[]string{},
			&condition{
				specs: []*v1.HTTPHasherSpec{
					{},
				},
			},
			&action{
				hs: []HTTPHasher{
					&clientAddrHasher{},
				},
			},
		),
		gen(
			"single hasher",
			[]string{},
			[]string{},
			&condition{
				specs: []*v1.HTTPHasherSpec{
					{
						HasherType: v1.HTTPHasherType_Header,
						Key:        "Test",
					},
				},
			},
			&action{
				hs: []HTTPHasher{
					&headerHasher{name: "Test"},
				},
			},
		),
		gen(
			"multiple hasher",
			[]string{},
			[]string{},
			&condition{
				specs: []*v1.HTTPHasherSpec{
					{
						HasherType: v1.HTTPHasherType_Header,
						Key:        "Test",
					},
					{
						HasherType: v1.HTTPHasherType_Query,
						Key:        "Test",
					},
				},
			},
			&action{
				hs: []HTTPHasher{
					&headerHasher{name: "Test"},
					&queryHasher{name: "Test"},
				},
			},
		),
		gen(
			"contains nil",
			[]string{},
			[]string{},
			&condition{
				specs: []*v1.HTTPHasherSpec{
					{
						HasherType: v1.HTTPHasherType_Header,
						Key:        "Test",
					},
					nil, nil,
					{
						HasherType: v1.HTTPHasherType_Query,
						Key:        "Test",
					},
				},
			},
			&action{
				hs: []HTTPHasher{
					&headerHasher{name: "Test"},
					&queryHasher{name: "Test"},
				},
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			hs := NewHTTPHashers(tt.C().specs)

			opts := []cmp.Option{
				cmp.AllowUnexported(clientAddrHasher{}),
				cmp.AllowUnexported(headerHasher{}, multiHeaderHasher{}),
				cmp.AllowUnexported(cookieHasher{}),
				cmp.AllowUnexported(queryHasher{}),
				cmp.AllowUnexported(pathParamHasher{}),
				cmp.Comparer(testutil.ComparePointer[hash.HashFunc]),
			}
			testutil.Diff(t, tt.A().hs, hs, opts...)
		})
	}
}

func TestNewHasher(t *testing.T) {
	type condition struct {
		spec *v1.HTTPHasherSpec
	}

	type action struct {
		h HTTPHasher
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"nil spec",
			[]string{},
			[]string{},
			&condition{
				spec: nil,
			},
			&action{
				h: nil,
			},
		),
		gen(
			"empty spec",
			[]string{},
			[]string{},
			&condition{
				spec: &v1.HTTPHasherSpec{},
			},
			&action{
				h: &clientAddrHasher{},
			},
		),
		gen(
			"header hasher",
			[]string{},
			[]string{},
			&condition{
				spec: &v1.HTTPHasherSpec{
					HasherType: v1.HTTPHasherType_Header,
					Key:        "Test",
				},
			},
			&action{
				h: &headerHasher{
					name: "Test",
				},
			},
		),
		gen(
			"multiHeader hasher",
			[]string{},
			[]string{},
			&condition{
				spec: &v1.HTTPHasherSpec{
					HasherType: v1.HTTPHasherType_MultiHeader,
					Keys:       []string{"Test1", "Test2"},
				},
			},
			&action{
				h: &multiHeaderHasher{
					names: []string{"Test1", "Test2"},
				},
			},
		),
		gen(
			"cookie hasher",
			[]string{},
			[]string{},
			&condition{
				spec: &v1.HTTPHasherSpec{
					HasherType: v1.HTTPHasherType_Cookie,
					Key:        "Test",
				},
			},
			&action{
				h: &cookieHasher{
					name: "Test",
				},
			},
		),
		gen(
			"query hasher",
			[]string{},
			[]string{},
			&condition{
				spec: &v1.HTTPHasherSpec{
					HasherType: v1.HTTPHasherType_Query,
					Key:        "Test",
				},
			},
			&action{
				h: &queryHasher{
					name: "Test",
				},
			},
		),
		gen(
			"pathParam hasher",
			[]string{},
			[]string{},
			&condition{
				spec: &v1.HTTPHasherSpec{
					HasherType: v1.HTTPHasherType_PathParam,
					Key:        "Test",
				},
			},
			&action{
				h: &pathParamHasher{
					name: "Test",
				},
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			h := NewHTTPHasher(tt.C().spec)

			opts := []cmp.Option{
				cmp.AllowUnexported(clientAddrHasher{}),
				cmp.AllowUnexported(headerHasher{}, multiHeaderHasher{}),
				cmp.AllowUnexported(cookieHasher{}),
				cmp.AllowUnexported(queryHasher{}),
				cmp.AllowUnexported(pathParamHasher{}),
				cmp.Comparer(testutil.ComparePointer[hash.HashFunc]),
			}
			testutil.Diff(t, tt.A().h, h, opts...)
		})
	}
}

func TestClientAddrHasher(t *testing.T) {
	type condition struct {
		addr string
	}

	type action struct {
		val int
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"empty value",
			[]string{},
			[]string{},
			&condition{
				addr: "",
			},
			&action{
				val: 1083068130,
			},
		),
		gen(
			"FNV1a_32",
			[]string{},
			[]string{},
			&condition{
				addr: "192.168.0.1",
			},
			&action{
				val: 2076768497,
			},
		),
	}
	testutil.Register(table, testCases...)
	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			h := &clientAddrHasher{}

			r, _ := http.NewRequest(http.MethodGet, "http://test.com/test", nil)
			r.RemoteAddr = tt.C().addr

			val, ok := h.Hash(r)
			testutil.Diff(t, tt.A().val, val)
			testutil.Diff(t, true, ok)
		})
	}
}

func TestHeaderHasher(t *testing.T) {
	type condition struct {
		name   string
		header http.Header
	}

	type action struct {
		val int
		ok  bool
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"empty name",
			[]string{},
			[]string{},
			&condition{
				name:   "",
				header: http.Header{"Test": []string{"foo"}},
			},
			&action{
				val: -1,
				ok:  false,
			},
		),
		gen(
			"empty value",
			[]string{},
			[]string{},
			&condition{
				name:   "Test",
				header: http.Header{"Test": []string{}},
			},
			&action{
				val: -1,
				ok:  false,
			},
		),
		gen(
			"FNV1a_32",
			[]string{},
			[]string{},
			&condition{
				name:   "Test",
				header: http.Header{"Test": []string{"foo"}},
			},
			&action{
				val: 1425653611,
				ok:  true,
			},
		),
	}
	testutil.Register(table, testCases...)
	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			h := &headerHasher{
				name: tt.C().name,
			}
			r, _ := http.NewRequest(http.MethodGet, "http://test.com/test", nil)
			r.Header = tt.C().header
			val, ok := h.Hash(r)
			testutil.Diff(t, tt.A().val, val)
			testutil.Diff(t, tt.A().ok, ok)
		})
	}
}

func TestMultiHeaderHasher(t *testing.T) {
	type condition struct {
		names  []string
		header http.Header
	}

	type action struct {
		val int
		ok  bool
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"empty name",
			[]string{},
			[]string{},
			&condition{
				names:  []string{""},
				header: http.Header{"Test1": []string{"foo"}, "Test2": []string{"bar"}},
			},
			&action{
				val: -1,
				ok:  false,
			},
		),
		gen(
			"empty value",
			[]string{},
			[]string{},
			&condition{
				names:  []string{"Test1", "Test2"},
				header: http.Header{"Test1": []string{}, "Test2": []string{}},
			},
			&action{
				val: -1,
				ok:  false,
			},
		),
		gen(
			"1 name",
			[]string{},
			[]string{},
			&condition{
				names:  []string{"Test1"},
				header: http.Header{"Test1": []string{"foo"}, "Test2": []string{"bar"}},
			},
			&action{
				val: 1425653611,
				ok:  true,
			},
		),
		gen(
			"2 names",
			[]string{},
			[]string{},
			&condition{
				names:  []string{"Test1", "Test2"},
				header: http.Header{"Test1": []string{"foo"}, "Test2": []string{"bar"}},
			},
			&action{
				val: 1607367860,
				ok:  true,
			},
		),
		gen(
			"3 names",
			[]string{},
			[]string{},
			&condition{
				names:  []string{"Test1", "Test2", "Test3"},
				header: http.Header{"Test1": []string{"foo"}, "Test2": []string{"bar"}},
			},
			&action{
				val: 1607367860,
				ok:  true,
			},
		),
		gen(
			"FNV1a_32",
			[]string{},
			[]string{},
			&condition{
				names:  []string{"Test1", "Test2"},
				header: http.Header{"Test1": []string{"foo"}, "Test2": []string{"bar"}},
			},
			&action{
				val: 1607367860,
				ok:  true,
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			h := &multiHeaderHasher{
				names: tt.C().names,
			}
			r, _ := http.NewRequest(http.MethodGet, "http://test.com/test", nil)
			r.Header = tt.C().header
			val, ok := h.Hash(r)
			testutil.Diff(t, tt.A().val, val)
			testutil.Diff(t, tt.A().ok, ok)
		})
	}
}

func TestCookieHasher(t *testing.T) {
	type condition struct {
		name   string
		header http.Header
	}

	type action struct {
		val int
		ok  bool
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"empty name",
			[]string{},
			[]string{},
			&condition{
				name:   "",
				header: http.Header{"Cookie": []string{"Dummy1=dum1;Test=foo;Dummy2=dum2;"}},
			},
			&action{
				val: -1,
				ok:  false,
			},
		),
		gen(
			"empty value",
			[]string{},
			[]string{},
			&condition{
				name:   "Test",
				header: http.Header{"Cookie": []string{"Dummy1=dum1;Test=;Dummy2=dum2;"}},
			},
			&action{
				val: -1,
				ok:  false,
			},
		),
		gen(
			"FNV1a_32",
			[]string{},
			[]string{},
			&condition{
				name:   "Test",
				header: http.Header{"Cookie": []string{"Dummy1=dum1;Test=foo;Dummy2=dum2;"}},
			},
			&action{
				val: 1425653611,
				ok:  true,
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			h := &cookieHasher{
				name: tt.C().name,
			}
			r, _ := http.NewRequest(http.MethodGet, "http://test.com/test", nil)
			r.Header = tt.C().header
			val, ok := h.Hash(r)
			testutil.Diff(t, tt.A().val, val)
			testutil.Diff(t, tt.A().ok, ok)
		})
	}
}

func TestQueryHasher(t *testing.T) {
	type condition struct {
		name string
		url  string
	}

	type action struct {
		val int
		ok  bool
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"empty name",
			[]string{},
			[]string{},
			&condition{
				name: "",
				url:  "/path?dummy=dum&test=foo",
			},
			&action{
				val: -1,
				ok:  false,
			},
		),
		gen(
			"empty value",
			[]string{},
			[]string{},
			&condition{
				name: "test",
				url:  "/path?dummy=dum&test=",
			},
			&action{
				val: -1,
				ok:  false,
			},
		),
		gen(
			"not found",
			[]string{},
			[]string{},
			&condition{
				name: "wrong",
				url:  "/path?dummy=dum&test=foo",
			},
			&action{
				val: -1,
				ok:  false,
			},
		),
		gen(
			"FNV1a_32",
			[]string{},
			[]string{},
			&condition{
				name: "test",
				url:  "/path?dummy=dum&test=foo",
			},
			&action{
				val: 1425653611,
				ok:  true,
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			h := &queryHasher{
				name: tt.C().name,
			}
			r, _ := http.NewRequest(http.MethodGet, "http://test.com"+tt.C().url, nil)
			val, ok := h.Hash(r)
			testutil.Diff(t, tt.A().val, val)
			testutil.Diff(t, tt.A().ok, ok)
		})
	}
}

func TestPathParamHasher(t *testing.T) {
	type condition struct {
		name string
		url  string
	}

	type action struct {
		val int
		ok  bool
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"empty name",
			[]string{},
			[]string{},
			&condition{
				name: "",
				url:  "/test/foo",
			},
			&action{
				val: -1,
				ok:  false,
			},
		),
		gen(
			"not found",
			[]string{},
			[]string{},
			&condition{
				name: "wrong",
				url:  "/test/foo",
			},
			&action{
				val: -1,
				ok:  false,
			},
		),
		gen(
			"FNV1a_32",
			[]string{},
			[]string{},
			&condition{
				name: "test",
				url:  "/test/foo",
			},
			&action{
				val: 1425653611,
				ok:  true,
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			h := &pathParamHasher{
				name: tt.C().name,
			}
			r, _ := http.NewRequest(http.MethodGet, "http://test.com"+tt.C().url, nil)
			w := httptest.NewRecorder()
			var rr *http.Request
			mux := &http.ServeMux{}
			mux.Handle("/test/{test}", http.HandlerFunc(
				func(w http.ResponseWriter, r *http.Request) {
					rr = r
				},
			))
			mux.ServeHTTP(w, r)

			val, ok := h.Hash(rr)
			testutil.Diff(t, tt.A().val, val)
			testutil.Diff(t, tt.A().ok, ok)
		})
	}
}
