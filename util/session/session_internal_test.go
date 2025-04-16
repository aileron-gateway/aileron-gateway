// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package session

import (
	"encoding/json"
	"errors"
	"net"
	"reflect"
	"testing"

	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
	"github.com/google/go-cmp/cmp"
	"github.com/vmihailenco/msgpack/v5"
)

func TestNewDefaultSession(t *testing.T) {
	type condition struct {
		sm SerializeMethod
	}

	type action struct {
		df *DefaultSession
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"default",
			[]string{},
			[]string{},
			&condition{},
			&action{
				df: &DefaultSession{
					flags:     New,
					attrs:     map[string]any{},
					data:      map[string][]byte{},
					marshal:   msgpack.Marshal,
					unmarshal: msgpack.Unmarshal,
				},
			},
		),
		gen(
			"msgpack",
			[]string{},
			[]string{},
			&condition{
				sm: SerializeMsgPack,
			},
			&action{
				df: &DefaultSession{
					flags:     New,
					attrs:     map[string]any{},
					data:      map[string][]byte{},
					marshal:   msgpack.Marshal,
					unmarshal: msgpack.Unmarshal,
				},
			},
		),
		gen(
			"json",
			[]string{},
			[]string{},
			&condition{
				sm: SerializeJSON,
			},
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

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			df := NewDefaultSession(tt.C().sm)

			opts := []cmp.Option{
				cmp.AllowUnexported(DefaultSession{}),
				cmp.Comparer(testutil.ComparePointer[func(any) ([]byte, error)]),
				cmp.Comparer(testutil.ComparePointer[func([]byte, any) error]),
			}
			testutil.Diff(t, tt.A().df, df, opts...)
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

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"set zero",
			[]string{},
			[]string{},
			&condition{
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
			"set non- zero",
			[]string{},
			[]string{},
			&condition{
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

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			result := tt.C().df.SetFlag(tt.C().flag)
			testutil.Diff(t, true, tt.C().df.flags&tt.A().flagTrue > 0)
			testutil.Diff(t, true, tt.C().df.flags&tt.A().flagFalse == 0)
			testutil.Diff(t, true, result&tt.A().flagTrue > 0)
			testutil.Diff(t, true, result&tt.A().flagFalse == 0)
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

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"non nil",
			[]string{},
			[]string{},
			&condition{
				df: &DefaultSession{
					attrs: map[string]any{"foo": "bar"},
				},
			},
			&action{
				attrs: map[string]any{"foo": "bar"},
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			attrs := tt.C().df.Attributes()
			testutil.Diff(t, tt.A().attrs, attrs)
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

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"delete existing key",
			[]string{},
			[]string{},
			&condition{
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
			"delete non-existing key",
			[]string{},
			[]string{},
			&condition{
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

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			tt.C().df.Delete(tt.C().key)
			testutil.Diff(t, tt.A().data, tt.C().df.data)
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

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	testStr := testMarshalUnmarshaler("bar")

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"msgpack/persist nil",
			[]string{},
			[]string{},
			&condition{
				df: &DefaultSession{
					data:    map[string][]byte{},
					marshal: msgpack.Marshal,
				},
				key: "foo",
				val: nil,
			},
			&action{
				data: map[string][]byte{
					"foo": {0xc0},
				},
			},
		),
		gen(
			"msgpack/persist value",
			[]string{},
			[]string{},
			&condition{
				df: &DefaultSession{
					data:    map[string][]byte{},
					marshal: msgpack.Marshal,
				},
				key: "foo",
				val: "bar",
			},
			&action{
				data: map[string][]byte{
					"foo": {0xa3, 0x62, 0x61, 0x72},
				},
			},
		),
		gen(
			"json/persist nil",
			[]string{},
			[]string{},
			&condition{
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
			"json/persist value",
			[]string{},
			[]string{},
			&condition{
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
			"binary marshaler",
			[]string{},
			[]string{},
			&condition{
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
			"marshal error",
			[]string{},
			[]string{},
			&condition{
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

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			err := tt.C().df.Persist(tt.C().key, tt.C().val)
			testutil.Diff(t, tt.A().data, tt.C().df.data)

			if tt.A().err != nil {
				testutil.Diff(t, tt.A().err.Error(), err.Error())
				testutil.Diff(t, false, tt.C().df.flags&Updated > 0)
			} else {
				testutil.Diff(t, nil, err)
				testutil.Diff(t, true, tt.C().df.flags&Updated > 0)
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

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	testStrEmpty := testMarshalUnmarshaler("")
	testStr := testMarshalUnmarshaler("bar")
	strPtr := func(s string) *string { return &s }

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"msgpack/extract nil",
			[]string{},
			[]string{},
			&condition{
				df: &DefaultSession{
					data: map[string][]byte{
						"foo": {0xc0},
					},
					unmarshal: msgpack.Unmarshal,
				},
				key: "foo",
				val: strPtr(""),
			},
			&action{
				val: strPtr(""),
				data: map[string][]byte{
					"foo": {0xc0},
				},
			},
		),
		gen(
			"msgpack/extract value",
			[]string{},
			[]string{},
			&condition{
				df: &DefaultSession{
					data: map[string][]byte{
						"foo": {0xa3, 0x62, 0x61, 0x72},
					},
					unmarshal: msgpack.Unmarshal,
				},
				key: "foo",
				val: strPtr("bar"),
			},
			&action{
				val: strPtr("bar"),
				data: map[string][]byte{
					"foo": {0xa3, 0x62, 0x61, 0x72},
				},
			},
		),
		gen(
			"json/extract nil",
			[]string{},
			[]string{},
			&condition{
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
			"json/extract value",
			[]string{},
			[]string{},
			&condition{
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
			"binary unmarshaler",
			[]string{},
			[]string{},
			&condition{
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
			"not found",
			[]string{},
			[]string{},
			&condition{
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

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			err := tt.C().df.Extract(tt.C().key, tt.A().val)
			testutil.Diff(t, &tt.A().data, &tt.C().df.data)

			if tt.A().err != nil {
				testutil.Diff(t, tt.A().err.Error(), err.Error())
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

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"msgpack/unmarshal nil",
			[]string{},
			[]string{},
			&condition{
				df: &DefaultSession{
					data:      map[string][]byte{},
					unmarshal: msgpack.Unmarshal,
				},
				b: []byte{0xc0},
			},
			&action{
				m: nil,
			},
		),
		gen(
			"msgpack/unmarshal empty",
			[]string{},
			[]string{},
			&condition{
				df: &DefaultSession{
					data:      map[string][]byte{},
					unmarshal: msgpack.Unmarshal,
				},
				b: []byte{0x80},
			},
			&action{
				m: map[string][]byte{},
			},
		),
		gen(
			"msgpack/unmarshal",
			[]string{},
			[]string{},
			&condition{
				df: &DefaultSession{
					data:      map[string][]byte{},
					unmarshal: msgpack.Unmarshal,
				},
				b: []byte{0x81, 0xa3, 0x66, 0x6f, 0x6f, 0xc4, 0x04, 0x74, 0x65, 0x73, 0x74},
			},
			&action{
				m: map[string][]byte{"foo": []byte("test")},
			},
		),
		gen(
			"json/unmarshal nil",
			[]string{},
			[]string{},
			&condition{
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
			"json/unmarshal empty",
			[]string{},
			[]string{},
			&condition{
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
			"json/unmarshal",
			[]string{},
			[]string{},
			&condition{
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

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			err := tt.C().df.UnmarshalBinary(tt.C().b)
			testutil.Diff(t, tt.A().err, err)
			testutil.Diff(t, tt.A().m, tt.C().df.data)
			testutil.Diff(t, true, tt.C().df.flags&Restored > 0)
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

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"not updated",
			[]string{},
			[]string{},
			&condition{
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
			"msgpack/marshal nil",
			[]string{},
			[]string{},
			&condition{
				df: &DefaultSession{
					flags:   Updated,
					data:    nil,
					marshal: msgpack.Marshal,
				},
			},
			&action{
				b: []byte{0xc0},
			},
		),
		gen(
			"msgpack/marshal empty",
			[]string{},
			[]string{},
			&condition{
				df: &DefaultSession{
					flags:   Updated,
					data:    map[string][]byte{},
					marshal: msgpack.Marshal,
				},
			},
			&action{
				b: []byte{0x80},
			},
		),
		gen(
			"msgpack/marshal",
			[]string{},
			[]string{},
			&condition{
				df: &DefaultSession{
					flags:   Updated,
					data:    map[string][]byte{"foo": []byte("test")},
					marshal: msgpack.Marshal,
				},
			},
			&action{
				b: []byte{0x81, 0xa3, 0x66, 0x6f, 0x6f, 0xc4, 0x04, 0x74, 0x65, 0x73, 0x74},
			},
		),
		gen(
			"json/marshal nil",
			[]string{},
			[]string{},
			&condition{
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
			"json/marshal empty",
			[]string{},
			[]string{},
			&condition{
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
			"json/marshal",
			[]string{},
			[]string{},
			&condition{
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

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			b, err := tt.C().df.MarshalBinary()
			testutil.Diff(t, tt.A().err, err)
			t.Log(string(b))
			testutil.Diff(t, tt.A().b, b)
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

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"valid value",
			[]string{},
			[]string{},
			&condition{
				key:   "foo",
				value: &net.IPAddr{Zone: "test"}, // No meaning using IPAddr.
			},
			&action{
				value: &net.IPAddr{},
			},
		),
		gen(
			"invalid value",
			[]string{},
			[]string{},
			&condition{
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

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			defer func() {
				err := recover()
				testutil.Diff(t, tt.A().err,
					err, cmp.Comparer(func(x, y reflect.Type) bool {
						return x.String() == y.String()
					}),
				)
			}()

			ss := NewDefaultSession(SerializeJSON)
			MustPersist(ss, tt.C().key, tt.C().value)
			if tt.A().err != nil {
				ss.Extract(tt.C().key, tt.A().value)
				testutil.Diff(t, tt.C().value, tt.A().value)
			}
		})
	}
}
