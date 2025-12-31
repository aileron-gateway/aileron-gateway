// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package session

import (
	"encoding/json"
	"errors"
	"net"
	"reflect"
	"testing"

	"github.com/aileron-gateway/aileron-gateway/internal/testutil"
	"github.com/google/go-cmp/cmp"
)

func TestNewDefaultSession(t *testing.T) {
	type condition struct {
	}

	type action struct {
		df *DefaultSession
	}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"json", &condition{},
			&action{
				df: &DefaultSession{
					flags:     New,
					attrs:     map[string]any{},
					data:      map[string][]byte{},
					marshal:   json.Marshal,
					unmarshal: json.Unmarshal,
				},
			},
		),
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			df := NewDefaultSession()

			opts := []cmp.Option{
				cmp.AllowUnexported(DefaultSession{}),
				cmp.Comparer(testutil.ComparePointer[func(any) ([]byte, error)]),
				cmp.Comparer(testutil.ComparePointer[func([]byte, any) error]),
			}
			testutil.Diff(t, tt.A.df, df, opts...)
		})
	}
}

func TestDefaultSession_SetFlag(t *testing.T) {
	type condition struct {
		df   *DefaultSession
		flag uint
	}

	type action struct {
		flagTrue  uint
		flagFalse uint
	}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"set zero", &condition{
				df: &DefaultSession{
					flags: 0b0001,
				},
				flag: 0b0000,
			},
			&action{
				flagTrue:  0b0001,
				flagFalse: 0b0010,
			},
		),
		gen(
			"set non- zero", &condition{
				df: &DefaultSession{
					flags: 0b0001,
				},
				flag: 0b0010,
			},
			&action{
				flagTrue:  0b0011,
				flagFalse: 0b0100,
			},
		),
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			result := tt.C.df.SetFlag(tt.C.flag)
			testutil.Diff(t, true, tt.C.df.flags&tt.A.flagTrue > 0)
			testutil.Diff(t, true, tt.C.df.flags&tt.A.flagFalse == 0)
			testutil.Diff(t, true, result&tt.A.flagTrue > 0)
			testutil.Diff(t, true, result&tt.A.flagFalse == 0)
		})
	}
}

func TestDefaultSession_Attributes(t *testing.T) {
	type condition struct {
		df *DefaultSession
	}

	type action struct {
		attrs map[string]any
	}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"non nil", &condition{
				df: &DefaultSession{
					attrs: map[string]any{"foo": "bar"},
				},
			},
			&action{
				attrs: map[string]any{"foo": "bar"},
			},
		),
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			attrs := tt.C.df.Attributes()
			testutil.Diff(t, tt.A.attrs, attrs)
		})
	}
}

func TestDefaultSession_Delete(t *testing.T) {
	type condition struct {
		df  *DefaultSession
		key string
	}

	type action struct {
		data map[string][]byte
	}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"delete existing key", &condition{
				df: &DefaultSession{
					data: map[string][]byte{
						"foo":   []byte("bar"),
						"alice": []byte("bob"),
					},
				},
				key: "alice",
			},
			&action{
				data: map[string][]byte{
					"foo": []byte("bar"),
				},
			},
		),
		gen(
			"delete non-existing key", &condition{
				df: &DefaultSession{
					data: map[string][]byte{
						"foo":   []byte("bar"),
						"alice": []byte("bob"),
					},
				},
				key: "baz",
			},
			&action{
				data: map[string][]byte{
					"foo":   []byte("bar"),
					"alice": []byte("bob"),
				},
			},
		),
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			tt.C.df.Delete(tt.C.key)
			testutil.Diff(t, tt.A.data, tt.C.df.data)
		})
	}
}

type testMarshalUnmarshaler string

func (t *testMarshalUnmarshaler) UnmarshalBinary(b []byte) error {
	val := testMarshalUnmarshaler(string(b))
	t = &val
	return nil
}

func (t *testMarshalUnmarshaler) MarshalBinary() ([]byte, error) {
	return []byte(string(*t)), nil
}

func TestDefaultSession_Persist(t *testing.T) {
	type condition struct {
		df  *DefaultSession
		key string
		val any
	}

	type action struct {
		data map[string][]byte
		err  error
	}

	testStr := testMarshalUnmarshaler("bar")

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"json/persist nil", &condition{
				df: &DefaultSession{
					data:    map[string][]byte{},
					marshal: json.Marshal,
				},
				key: "foo",
				val: nil,
			},
			&action{
				data: map[string][]byte{
					"foo": []byte("null"),
				},
			},
		),
		gen(
			"json/persist value", &condition{
				df: &DefaultSession{
					data:    map[string][]byte{},
					marshal: json.Marshal,
				},
				key: "foo",
				val: "bar",
			},
			&action{
				data: map[string][]byte{
					"foo": []byte(`"bar"`),
				},
			},
		),
		gen(
			"binary marshaler", &condition{
				df: &DefaultSession{
					data: map[string][]byte{},
				},
				key: "foo",
				val: &testStr,
			},
			&action{
				data: map[string][]byte{
					"foo": []byte("bar"),
				},
			},
		),
		gen(
			"marshal error", &condition{
				df: &DefaultSession{
					data:    map[string][]byte{},
					marshal: json.Marshal,
				},
				key: "foo",
				val: complex(1, 2),
			},
			&action{
				data: map[string][]byte{},
				err:  errors.New("json: unsupported type: complex128"),
			},
		),
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			err := tt.C.df.Persist(tt.C.key, tt.C.val)
			testutil.Diff(t, tt.A.data, tt.C.df.data)

			if tt.A.err != nil {
				testutil.Diff(t, tt.A.err.Error(), err.Error())
				testutil.Diff(t, false, tt.C.df.flags&Updated > 0)
			} else {
				testutil.Diff(t, nil, err)
				testutil.Diff(t, true, tt.C.df.flags&Updated > 0)
			}
		})
	}
}

func TestDefaultSession_Extract(t *testing.T) {
	type condition struct {
		df  *DefaultSession
		key string
		val any
	}

	type action struct {
		val  any
		data map[string][]byte
		err  error
	}

	testStrEmpty := testMarshalUnmarshaler("")
	testStr := testMarshalUnmarshaler("bar")
	strPtr := func(s string) *string { return &s }

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"json/extract nil", &condition{
				df: &DefaultSession{
					data: map[string][]byte{
						"foo": []byte("null"),
					},
					unmarshal: json.Unmarshal,
				},
				key: "foo",
				val: strPtr(""),
			},
			&action{
				val: strPtr(""),
				data: map[string][]byte{
					"foo": []byte("null"),
				},
			},
		),
		gen(
			"json/extract value", &condition{
				df: &DefaultSession{
					data: map[string][]byte{
						"foo": []byte(`"bar"`),
					},
					unmarshal: json.Unmarshal,
				},
				key: "foo",
				val: strPtr(""),
			},
			&action{
				val: strPtr("bar"),
				data: map[string][]byte{
					"foo": []byte(`"bar"`),
				},
			},
		),
		gen(
			"binary unmarshaler", &condition{
				df: &DefaultSession{
					data: map[string][]byte{
						"foo": []byte("bar"),
					},
				},
				key: "foo",
				val: &testStrEmpty,
			},
			&action{
				val: &testStr,
				data: map[string][]byte{
					"foo": []byte("bar"),
				},
			},
		),
		gen(
			"not found", &condition{
				df: &DefaultSession{
					data: map[string][]byte{
						"foo": []byte("bar"),
					},
				},
				key: "baz",
				val: strPtr(""),
			},
			&action{
				val: strPtr(""),
				data: map[string][]byte{
					"foo": []byte("bar"),
				},
				err: NoValue,
			},
		),
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			err := tt.C.df.Extract(tt.C.key, tt.A.val)
			testutil.Diff(t, &tt.A.data, &tt.C.df.data)

			if tt.A.err != nil {
				testutil.Diff(t, tt.A.err.Error(), err.Error())
			} else {
				testutil.Diff(t, nil, err)
			}
		})
	}
}

func TestDefaultSession_UnmarshalBinary(t *testing.T) {
	type condition struct {
		df *DefaultSession
		b  []byte
	}

	type action struct {
		m   map[string][]byte
		err error
	}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"json/unmarshal nil", &condition{
				df: &DefaultSession{
					data:      map[string][]byte{},
					unmarshal: json.Unmarshal,
				},
				b: []byte("null"),
			},
			&action{
				m: nil,
			},
		),
		gen(
			"json/unmarshal empty", &condition{
				df: &DefaultSession{
					data:      map[string][]byte{},
					unmarshal: json.Unmarshal,
				},
				b: []byte("{}"),
			},
			&action{
				m: map[string][]byte{},
			},
		),
		gen(
			"json/unmarshal", &condition{
				df: &DefaultSession{
					data:      map[string][]byte{},
					unmarshal: json.Unmarshal,
				},
				b: []byte(`{"foo":"dGVzdA=="}`),
			},
			&action{
				m: map[string][]byte{"foo": []byte("test")},
			},
		),
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			err := tt.C.df.UnmarshalBinary(tt.C.b)
			testutil.Diff(t, tt.A.err, err)
			testutil.Diff(t, tt.A.m, tt.C.df.data)
			testutil.Diff(t, true, tt.C.df.flags&Restored > 0)
		})
	}
}

func TestDefaultSession_MarshalBinary(t *testing.T) {
	type condition struct {
		df *DefaultSession
	}

	type action struct {
		b   []byte
		err error
	}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"not updated", &condition{
				df: &DefaultSession{
					raw:  []byte("raw"),
					data: nil,
				},
			},
			&action{
				b: []byte("raw"),
			},
		),
		gen(
			"json/marshal nil", &condition{
				df: &DefaultSession{
					flags:   Updated,
					data:    nil,
					marshal: json.Marshal,
				},
			},
			&action{
				b: []byte("null"),
			},
		),
		gen(
			"json/marshal empty", &condition{
				df: &DefaultSession{
					flags:   Updated,
					data:    map[string][]byte{},
					marshal: json.Marshal,
				},
			},
			&action{
				b: []byte("{}"),
			},
		),
		gen(
			"json/marshal", &condition{
				df: &DefaultSession{
					flags:   Updated,
					data:    map[string][]byte{"foo": []byte("test")},
					marshal: json.Marshal,
				},
			},
			&action{
				b: []byte(`{"foo":"dGVzdA=="}`),
			},
		),
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			b, err := tt.C.df.MarshalBinary()
			testutil.Diff(t, tt.A.err, err)
			t.Log(string(b))
			testutil.Diff(t, tt.A.b, b)
		})
	}
}

func TestMustPersist(t *testing.T) {
	type condition struct {
		key   string
		value any
	}

	type action struct {
		value any
		err   error
	}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"valid value", &condition{
				key:   "foo",
				value: &net.IPAddr{Zone: "test"}, // No meaning using IPAddr.
			},
			&action{
				value: &net.IPAddr{},
			},
		),
		gen(
			"invalid value", &condition{
				key:   "foo",
				value: complex(1, 2),
			},
			&action{
				err: &json.UnsupportedTypeError{
					Type: reflect.TypeOf(complex(1, 2)),
				},
			},
		),
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			defer func() {
				err := recover()
				testutil.Diff(t, tt.A.err,
					err, cmp.Comparer(func(x, y reflect.Type) bool {
						return x.String() == y.String()
					}),
				)
			}()

			ss := NewDefaultSession()
			MustPersist(ss, tt.C.key, tt.C.value)
			if tt.A.err != nil {
				ss.Extract(tt.C.key, tt.A.value)
				testutil.Diff(t, tt.C.value, tt.A.value)
			}
		})
	}
}
