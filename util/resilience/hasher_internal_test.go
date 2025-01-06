package resilience

import (
	"net/http"
	"net/http/httptest"
	"testing"

	v1 "github.com/aileron-gateway/aileron-gateway/apis/core/v1"
	"github.com/aileron-gateway/aileron-gateway/apis/kernel"
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
					&clientAddrHasher{hashFunc: hash.FNV1a_32},
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
						HashAlg:    kernel.HashAlg_FNV1_32,
					},
				},
			},
			&action{
				hs: []HTTPHasher{
					&headerHasher{name: "Test", hashFunc: hash.FNV1_32},
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
						HashAlg:    kernel.HashAlg_FNV1_32,
					},
					{
						HasherType: v1.HTTPHasherType_Query,
						Key:        "Test",
						HashAlg:    kernel.HashAlg_FNV1_32,
					},
				},
			},
			&action{
				hs: []HTTPHasher{
					&headerHasher{name: "Test", hashFunc: hash.FNV1_32},
					&queryHasher{name: "Test", hashFunc: hash.FNV1_32},
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
						HashAlg:    kernel.HashAlg_FNV1_32,
					},
					nil, nil,
					{
						HasherType: v1.HTTPHasherType_Query,
						Key:        "Test",
						HashAlg:    kernel.HashAlg_FNV1_32,
					},
				},
			},
			&action{
				hs: []HTTPHasher{
					&headerHasher{name: "Test", hashFunc: hash.FNV1_32},
					&queryHasher{name: "Test", hashFunc: hash.FNV1_32},
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
				h: &clientAddrHasher{
					hashFunc: hash.FNV1a_32,
				},
			},
		),
		gen(
			"unknown hasher",
			[]string{},
			[]string{},
			&condition{
				spec: &v1.HTTPHasherSpec{
					HasherType: v1.HTTPHasherType(9999),
				},
			},
			&action{
				h: &clientAddrHasher{
					hashFunc: hash.FNV1a_32,
				},
			},
		),
		gen(
			"unknown hasher/set hash func",
			[]string{},
			[]string{},
			&condition{
				spec: &v1.HTTPHasherSpec{
					HasherType: v1.HTTPHasherType(9999),
					HashAlg:    kernel.HashAlg_FNV1_32,
				},
			},
			&action{
				h: &clientAddrHasher{
					hashFunc: hash.FNV1_32,
				},
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
					HashAlg:    kernel.HashAlg_FNV1_32,
				},
			},
			&action{
				h: &headerHasher{
					name:     "Test",
					hashFunc: hash.FNV1_32,
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
					HashAlg:    kernel.HashAlg_FNV1_32,
				},
			},
			&action{
				h: &multiHeaderHasher{
					names:    []string{"Test1", "Test2"},
					hashFunc: hash.FNV1_32,
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
					HashAlg:    kernel.HashAlg_FNV1_32,
				},
			},
			&action{
				h: &cookieHasher{
					name:     "Test",
					hashFunc: hash.FNV1_32,
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
					HashAlg:    kernel.HashAlg_FNV1_32,
				},
			},
			&action{
				h: &queryHasher{
					name:     "Test",
					hashFunc: hash.FNV1_32,
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
					HashAlg:    kernel.HashAlg_FNV1_32,
				},
			},
			&action{
				h: &pathParamHasher{
					name:     "Test",
					hashFunc: hash.FNV1_32,
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
		hf   hash.HashFunc
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
				hf:   hash.FNV1_32,
				addr: "",
			},
			&action{
				val: 1083068130,
			},
		),
		gen(
			"FNV1_32",
			[]string{},
			[]string{},
			&condition{
				hf:   hash.FNV1_32,
				addr: "192.168.0.1",
			},
			&action{
				val: 156524783,
			},
		),
		gen(
			"FNV1a_32",
			[]string{},
			[]string{},
			&condition{
				hf:   hash.FNV1a_32,
				addr: "192.168.0.1",
			},
			&action{
				val: 2076768497,
			},
		),
		gen(
			"FNV1_64",
			[]string{},
			[]string{},
			&condition{
				hf:   hash.FNV1_64,
				addr: "192.168.0.1",
			},
			&action{
				val: 1237457740,
			},
		),
		gen(
			"FNV1a_64",
			[]string{},
			[]string{},
			&condition{
				hf:   hash.FNV1a_64,
				addr: "192.168.0.1",
			},
			&action{
				val: 390611308,
			},
		),
		gen(
			"FNV1_128",
			[]string{},
			[]string{},
			&condition{
				hf:   hash.FNV1_128,
				addr: "192.168.0.1",
			},
			&action{
				val: 1997240431,
			},
		),
		gen(
			"FNV1a_128",
			[]string{},
			[]string{},
			&condition{
				hf:   hash.FNV1a_128,
				addr: "192.168.0.1",
			},
			&action{
				val: 1695602422,
			},
		),
		gen(
			"SHA512",
			[]string{},
			[]string{},
			&condition{
				hf:   hash.SHA512,
				addr: "192.168.0.1",
			},
			&action{
				val: 1032341250,
			},
		),
		gen(
			"BLAKE2b_512",
			[]string{},
			[]string{},
			&condition{
				hf:   hash.BLAKE2b_512,
				addr: "192.168.0.1",
			},
			&action{
				val: 132499035,
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			h := &clientAddrHasher{
				hashFunc: tt.C().hf,
			}

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
		hf     hash.HashFunc
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
				hf:     hash.FNV1_32,
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
				hf:     hash.FNV1_32,
				header: http.Header{"Test": []string{}},
			},
			&action{
				val: -1,
				ok:  false,
			},
		),
		gen(
			"FNV1_32",
			[]string{},
			[]string{},
			&condition{
				name:   "Test",
				hf:     hash.FNV1_32,
				header: http.Header{"Test": []string{"foo"}},
			},
			&action{
				val: 541568777,
				ok:  true,
			},
		),
		gen(
			"FNV1a_32",
			[]string{},
			[]string{},
			&condition{
				name:   "Test",
				hf:     hash.FNV1a_32,
				header: http.Header{"Test": []string{"foo"}},
			},
			&action{
				val: 1425653611,
				ok:  true,
			},
		),
		gen(
			"FNV1_64",
			[]string{},
			[]string{},
			&condition{
				name:   "Test",
				hf:     hash.FNV1_64,
				header: http.Header{"Test": []string{"foo"}},
			},
			&action{
				val: 1818616716,
				ok:  true,
			},
		),
		gen(
			"FNV1a_64",
			[]string{},
			[]string{},
			&condition{
				name:   "Test",
				hf:     hash.FNV1a_64,
				header: http.Header{"Test": []string{"foo"}},
			},
			&action{
				val: 1851341452,
				ok:  true,
			},
		),
		gen(
			"FNV1_128",
			[]string{},
			[]string{},
			&condition{
				name:   "Test",
				hf:     hash.FNV1_128,
				header: http.Header{"Test": []string{"foo"}},
			},
			&action{
				val: 1397086540,
				ok:  true,
			},
		),
		gen(
			"FNV1a_128",
			[]string{},
			[]string{},
			&condition{
				name:   "Test",
				hf:     hash.FNV1a_128,
				header: http.Header{"Test": []string{"foo"}},
			},
			&action{
				val: 1397141352,
				ok:  true,
			},
		),
		gen(
			"SHA512",
			[]string{},
			[]string{},
			&condition{
				name:   "Test",
				hf:     hash.SHA512,
				header: http.Header{"Test": []string{"foo"}},
			},
			&action{
				val: 2080234807,
				ok:  true,
			},
		),
		gen(
			"BLAKE2b_512",
			[]string{},
			[]string{},
			&condition{
				name:   "Test",
				hf:     hash.BLAKE2b_512,
				header: http.Header{"Test": []string{"foo"}},
			},
			&action{
				val: 1694503320,
				ok:  true,
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			h := &headerHasher{
				name:     tt.C().name,
				hashFunc: tt.C().hf,
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
		hf     hash.HashFunc
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
				hf:     hash.FNV1_32,
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
				hf:     hash.FNV1_32,
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
				hf:     hash.FNV1_32,
				header: http.Header{"Test1": []string{"foo"}, "Test2": []string{"bar"}},
			},
			&action{
				val: 541568777,
				ok:  true,
			},
		),
		gen(
			"2 names",
			[]string{},
			[]string{},
			&condition{
				names:  []string{"Test1", "Test2"},
				hf:     hash.FNV1_32,
				header: http.Header{"Test1": []string{"foo"}, "Test2": []string{"bar"}},
			},
			&action{
				val: 418928945,
				ok:  true,
			},
		),
		gen(
			"3 names",
			[]string{},
			[]string{},
			&condition{
				names:  []string{"Test1", "Test2", "Test3"},
				hf:     hash.FNV1_32,
				header: http.Header{"Test1": []string{"foo"}, "Test2": []string{"bar"}},
			},
			&action{
				val: 418928945,
				ok:  true,
			},
		),
		gen(
			"FNV1_32",
			[]string{},
			[]string{},
			&condition{
				names:  []string{"Test1", "Test2"},
				hf:     hash.FNV1_32,
				header: http.Header{"Test1": []string{"foo"}, "Test2": []string{"bar"}},
			},
			&action{
				val: 418928945,
				ok:  true,
			},
		),
		gen(
			"FNV1a_32",
			[]string{},
			[]string{},
			&condition{
				names:  []string{"Test1", "Test2"},
				hf:     hash.FNV1a_32,
				header: http.Header{"Test1": []string{"foo"}, "Test2": []string{"bar"}},
			},
			&action{
				val: 1607367860,
				ok:  true,
			},
		),
		gen(
			"FNV1_64",
			[]string{},
			[]string{},
			&condition{
				names:  []string{"Test1", "Test2"},
				hf:     hash.FNV1_64,
				header: http.Header{"Test1": []string{"foo"}, "Test2": []string{"bar"}},
			},
			&action{
				val: 436650930,
				ok:  true,
			},
		),
		gen(
			"FNV1a_64",
			[]string{},
			[]string{},
			&condition{
				names:  []string{"Test1", "Test2"},
				hf:     hash.FNV1a_64,
				header: http.Header{"Test1": []string{"foo"}, "Test2": []string{"bar"}},
			},
			&action{
				val: 1120542904,
				ok:  true,
			},
		),
		gen(
			"SHA512",
			[]string{},
			[]string{},
			&condition{
				names:  []string{"Test1", "Test2"},
				hf:     hash.SHA512,
				header: http.Header{"Test1": []string{"foo"}, "Test2": []string{"bar"}},
			},
			&action{
				val: 86512399,
				ok:  true,
			},
		),
		gen(
			"BLAKE2b_512",
			[]string{},
			[]string{},
			&condition{
				names:  []string{"Test1", "Test2"},
				hf:     hash.BLAKE2b_512,
				header: http.Header{"Test1": []string{"foo"}, "Test2": []string{"bar"}},
			},
			&action{
				val: 1190760368,
				ok:  true,
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			h := &multiHeaderHasher{
				names:    tt.C().names,
				hashFunc: tt.C().hf,
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
		hf     hash.HashFunc
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
				hf:     hash.FNV1_32,
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
				hf:     hash.FNV1_32,
				header: http.Header{"Cookie": []string{"Dummy1=dum1;Test=;Dummy2=dum2;"}},
			},
			&action{
				val: -1,
				ok:  false,
			},
		),
		gen(
			"FNV1_32",
			[]string{},
			[]string{},
			&condition{
				name:   "Test",
				hf:     hash.FNV1_32,
				header: http.Header{"Cookie": []string{"Dummy1=dum1;Test=foo;Dummy2=dum2;"}},
			},
			&action{
				val: 541568777,
				ok:  true,
			},
		),
		gen(
			"FNV1a_32",
			[]string{},
			[]string{},
			&condition{
				name:   "Test",
				hf:     hash.FNV1a_32,
				header: http.Header{"Cookie": []string{"Dummy1=dum1;Test=foo;Dummy2=dum2;"}},
			},
			&action{
				val: 1425653611,
				ok:  true,
			},
		),
		gen(
			"FNV1_64",
			[]string{},
			[]string{},
			&condition{
				name:   "Test",
				hf:     hash.FNV1_64,
				header: http.Header{"Cookie": []string{"Dummy1=dum1;Test=foo;Dummy2=dum2;"}},
			},
			&action{
				val: 1818616716,
				ok:  true,
			},
		),
		gen(
			"FNV1a_64",
			[]string{},
			[]string{},
			&condition{
				name:   "Test",
				hf:     hash.FNV1a_64,
				header: http.Header{"Cookie": []string{"Dummy1=dum1;Test=foo;Dummy2=dum2;"}},
			},
			&action{
				val: 1851341452,
				ok:  true,
			},
		),
		gen(
			"FNV1_128",
			[]string{},
			[]string{},
			&condition{
				name:   "Test",
				hf:     hash.FNV1_128,
				header: http.Header{"Cookie": []string{"Dummy1=dum1;Test=foo;Dummy2=dum2;"}},
			},
			&action{
				val: 1397086540,
				ok:  true,
			},
		),
		gen(
			"FNV1a_128",
			[]string{},
			[]string{},
			&condition{
				name:   "Test",
				hf:     hash.FNV1a_128,
				header: http.Header{"Cookie": []string{"Dummy1=dum1;Test=foo;Dummy2=dum2;"}},
			},
			&action{
				val: 1397141352,
				ok:  true,
			},
		),
		gen(
			"SHA512",
			[]string{},
			[]string{},
			&condition{
				name:   "Test",
				hf:     hash.SHA512,
				header: http.Header{"Cookie": []string{"Dummy1=dum1;Test=foo;Dummy2=dum2;"}},
			},
			&action{
				val: 2080234807,
				ok:  true,
			},
		),
		gen(
			"BLAKE2b_512",
			[]string{},
			[]string{},
			&condition{
				name:   "Test",
				hf:     hash.BLAKE2b_512,
				header: http.Header{"Cookie": []string{"Dummy1=dum1;Test=foo;Dummy2=dum2;"}},
			},
			&action{
				val: 1694503320,
				ok:  true,
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			h := &cookieHasher{
				name:     tt.C().name,
				hashFunc: tt.C().hf,
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
		hf   hash.HashFunc
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
				hf:   hash.FNV1_32,
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
				hf:   hash.FNV1_32,
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
				hf:   hash.FNV1_32,
				url:  "/path?dummy=dum&test=foo",
			},
			&action{
				val: -1,
				ok:  false,
			},
		),
		gen(
			"FNV1_32",
			[]string{},
			[]string{},
			&condition{
				name: "test",
				hf:   hash.FNV1_32,
				url:  "/path?dummy=dum&test=foo",
			},
			&action{
				val: 541568777,
				ok:  true,
			},
		),
		gen(
			"FNV1a_32",
			[]string{},
			[]string{},
			&condition{
				name: "test",
				hf:   hash.FNV1a_32,
				url:  "/path?dummy=dum&test=foo",
			},
			&action{
				val: 1425653611,
				ok:  true,
			},
		),
		gen(
			"FNV1_64",
			[]string{},
			[]string{},
			&condition{
				name: "test",
				hf:   hash.FNV1_64,
				url:  "/path?dummy=dum&test=foo",
			},
			&action{
				val: 1818616716,
				ok:  true,
			},
		),
		gen(
			"FNV1a_64",
			[]string{},
			[]string{},
			&condition{
				name: "test",
				hf:   hash.FNV1a_64,
				url:  "/path?dummy=dum&test=foo",
			},
			&action{
				val: 1851341452,
				ok:  true,
			},
		),
		gen(
			"FNV1_128",
			[]string{},
			[]string{},
			&condition{
				name: "test",
				hf:   hash.FNV1_128,
				url:  "/path?dummy=dum&test=foo",
			},
			&action{
				val: 1397086540,
				ok:  true,
			},
		),
		gen(
			"FNV1a_128",
			[]string{},
			[]string{},
			&condition{
				name: "test",
				hf:   hash.FNV1a_128,
				url:  "/path?dummy=dum&test=foo",
			},
			&action{
				val: 1397141352,
				ok:  true,
			},
		),
		gen(
			"SHA512",
			[]string{},
			[]string{},
			&condition{
				name: "test",
				hf:   hash.SHA512,
				url:  "/path?dummy=dum&test=foo",
			},
			&action{
				val: 2080234807,
				ok:  true,
			},
		),
		gen(
			"BLAKE2b_512",
			[]string{},
			[]string{},
			&condition{
				name: "test",
				hf:   hash.BLAKE2b_512,
				url:  "/path?dummy=dum&test=foo",
			},
			&action{
				val: 1694503320,
				ok:  true,
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			h := &queryHasher{
				name:     tt.C().name,
				hashFunc: tt.C().hf,
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
		hf   hash.HashFunc
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
				hf:   hash.FNV1_32,
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
				hf:   hash.FNV1_32,
				url:  "/test/foo",
			},
			&action{
				val: -1,
				ok:  false,
			},
		),
		gen(
			"FNV1_32",
			[]string{},
			[]string{},
			&condition{
				name: "test",
				hf:   hash.FNV1_32,
				url:  "/test/foo",
			},
			&action{
				val: 541568777,
				ok:  true,
			},
		),
		gen(
			"FNV1a_32",
			[]string{},
			[]string{},
			&condition{
				name: "test",
				hf:   hash.FNV1a_32,
				url:  "/test/foo",
			},
			&action{
				val: 1425653611,
				ok:  true,
			},
		),
		gen(
			"FNV1_64",
			[]string{},
			[]string{},
			&condition{
				name: "test",
				hf:   hash.FNV1_64,
				url:  "/test/foo",
			},
			&action{
				val: 1818616716,
				ok:  true,
			},
		),
		gen(
			"FNV1a_64",
			[]string{},
			[]string{},
			&condition{
				name: "test",
				hf:   hash.FNV1a_64,
				url:  "/test/foo",
			},
			&action{
				val: 1851341452,
				ok:  true,
			},
		),
		gen(
			"FNV1_128",
			[]string{},
			[]string{},
			&condition{
				name: "test",
				hf:   hash.FNV1_128,
				url:  "/test/foo",
			},
			&action{
				val: 1397086540,
				ok:  true,
			},
		),
		gen(
			"FNV1a_128",
			[]string{},
			[]string{},
			&condition{
				name: "test",
				hf:   hash.FNV1a_128,
				url:  "/test/foo",
			},
			&action{
				val: 1397141352,
				ok:  true,
			},
		),
		gen(
			"SHA512",
			[]string{},
			[]string{},
			&condition{
				name: "test",
				hf:   hash.SHA512,
				url:  "/test/foo",
			},
			&action{
				val: 2080234807,
				ok:  true,
			},
		),
		gen(
			"BLAKE2b_512",
			[]string{},
			[]string{},
			&condition{
				name: "test",
				hf:   hash.BLAKE2b_512,
				url:  "/test/foo",
			},
			&action{
				val: 1694503320,
				ok:  true,
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			h := &pathParamHasher{
				name:     tt.C().name,
				hashFunc: tt.C().hf,
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
