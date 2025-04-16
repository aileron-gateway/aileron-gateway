// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package httplogger

import (
	"encoding/base64"
	"encoding/binary"
	"testing"
	"time"

	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
)

func TestRequestAttrs_accessKeyValues(t *testing.T) {
	type condition struct {
		attr *requestAttrs
	}

	type action struct {
		val []any
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"full value",
			[]string{},
			[]string{},
			&condition{
				attr: &requestAttrs{
					typ:    "test-type",
					id:     "test-id",
					time:   "test-time",
					host:   "test-host",
					method: "test-method",
					path:   "test-path",
					query:  "test-query",
					remote: "test-remote",
					proto:  "test-proto",
					size:   10,
					header: map[string]string{"key": "value"},
					body:   "test-body",
				},
			},
			&action{
				val: []any{
					"request",
					map[string]any{
						keyID:     "test-id",
						keyTime:   "test-time",
						keyHost:   "test-host",
						keyMethod: "test-method",
						keyPath:   "test-path",
						keyQuery:  "test-query",
						keyRemote: "test-remote",
						keyProto:  "test-proto",
						keySize:   int64(10),
						keyHeader: map[string]string{"key": "value"},
					},
				},
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			val := tt.C().attr.accessKeyValues()
			testutil.Diff(t, tt.A().val, val)
		})
	}
}

func TestRequestAttrs_journalKeyValues(t *testing.T) {
	type condition struct {
		attr *requestAttrs
	}

	type action struct {
		val []any
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"full value",
			[]string{},
			[]string{},
			&condition{
				attr: &requestAttrs{
					typ:    "test-type",
					id:     "test-id",
					time:   "test-time",
					host:   "test-host",
					method: "test-method",
					path:   "test-path",
					query:  "test-query",
					remote: "test-remote",
					proto:  "test-proto",
					size:   10,
					header: map[string]string{"key": "value"},
					body:   "test-body",
				},
			},
			&action{
				val: []any{
					"request",
					map[string]any{
						keyID:     "test-id",
						keyTime:   "test-time",
						keyHost:   "test-host",
						keyMethod: "test-method",
						keyPath:   "test-path",
						keyQuery:  "test-query",
						keyRemote: "test-remote",
						keyProto:  "test-proto",
						keySize:   int64(10),
						keyHeader: map[string]string{"key": "value"},
						keyBody:   "test-body",
					},
				},
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			val := tt.C().attr.journalKeyValues()
			testutil.Diff(t, tt.A().val, val)
		})
	}
}

func TestRequestAttrs_TagFunc(t *testing.T) {
	type condition struct {
		attr *requestAttrs
		tag  string
	}

	type action struct {
		val string
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	testAttr := &requestAttrs{
		typ:    "test-type",
		id:     "test-id",
		time:   "test-time",
		host:   "test-host",
		method: "test-method",
		path:   "test-path",
		query:  "test-query",
		remote: "test-remote",
		proto:  "test-proto",
		size:   10,
		header: map[string]string{"Key": "value"},
		body:   "test-body",
	}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"id",
			[]string{},
			[]string{},
			&condition{
				attr: testAttr,
				tag:  "id",
			},
			&action{
				val: "test-id",
			},
		),
		gen(
			"time",
			[]string{},
			[]string{},
			&condition{
				attr: testAttr,
				tag:  "time",
			},
			&action{
				val: "test-time",
			},
		),
		gen(
			"host",
			[]string{},
			[]string{},
			&condition{
				attr: testAttr,
				tag:  "host",
			},
			&action{
				val: "test-host",
			},
		),
		gen(
			"method",
			[]string{},
			[]string{},
			&condition{
				attr: testAttr,
				tag:  "method",
			},
			&action{
				val: "test-method",
			},
		),
		gen(
			"path",
			[]string{},
			[]string{},
			&condition{
				attr: testAttr,
				tag:  "path",
			},
			&action{
				val: "test-path",
			},
		),
		gen(
			"query",
			[]string{},
			[]string{},
			&condition{
				attr: testAttr,
				tag:  "query",
			},
			&action{
				val: "test-query",
			},
		),
		gen(
			"remote",
			[]string{},
			[]string{},
			&condition{
				attr: testAttr,
				tag:  "remote",
			},
			&action{
				val: "test-remote",
			},
		),
		gen(
			"proto",
			[]string{},
			[]string{},
			&condition{
				attr: testAttr,
				tag:  "proto",
			},
			&action{
				val: "test-proto",
			},
		),
		gen(
			"size",
			[]string{},
			[]string{},
			&condition{
				attr: testAttr,
				tag:  "size",
			},
			&action{
				val: "10",
			},
		),
		gen(
			"header",
			[]string{},
			[]string{},
			&condition{
				attr: testAttr,
				tag:  "header",
			},
			&action{
				val: "map[Key:value]",
			},
		),
		gen(
			"r.body",
			[]string{},
			[]string{},
			&condition{
				attr: testAttr,
				tag:  "body",
			},
			&action{
				val: "test-body",
			},
		),
		gen(
			"header.key",
			[]string{},
			[]string{},
			&condition{
				attr: testAttr,
				tag:  "header.key",
			},
			&action{
				val: "value",
			},
		),
		gen(
			"type",
			[]string{},
			[]string{},
			&condition{
				attr: testAttr,
				tag:  "type",
			},
			&action{
				val: "test-type",
			},
		),
		gen(
			"not.exist",
			[]string{},
			[]string{},
			&condition{
				attr: testAttr,
				tag:  "not.exist",
			},
			&action{
				val: "<undefined:not.exist>",
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			val := tt.C().attr.TagFunc(tt.C().tag)
			testutil.Diff(t, tt.A().val, string(val))
		})
	}
}

func TestResponseAttrs_accessKeyValues(t *testing.T) {
	type condition struct {
		attr *responseAttrs
	}

	type action struct {
		val []any
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"full value",
			[]string{},
			[]string{},
			&condition{
				attr: &responseAttrs{
					typ:      "test-type",
					id:       "test-id",
					time:     "test-time",
					duration: 1,
					status:   200,
					size:     10,
					header:   map[string]string{"key": "value"},
					body:     "test-body",
				},
			},
			&action{
				val: []any{
					"response",
					map[string]any{
						keyID:       "test-id",
						keyTime:     "test-time",
						keyDuration: int64(1),
						keyStatus:   200,
						keySize:     int64(10),
						keyHeader:   map[string]string{"key": "value"},
					},
				},
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			val := tt.C().attr.accessKeyValues()
			testutil.Diff(t, tt.A().val, val)
		})
	}
}

func TestResponseAttrs_journalKeyValues(t *testing.T) {
	type condition struct {
		attr *responseAttrs
	}

	type action struct {
		val []any
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"full value",
			[]string{},
			[]string{},
			&condition{
				attr: &responseAttrs{
					typ:      "test-type",
					id:       "test-id",
					time:     "test-time",
					duration: 1,
					status:   200,
					size:     10,
					header:   map[string]string{"key": "value"},
					body:     "test-body",
				},
			},
			&action{
				val: []any{
					"response",
					map[string]any{
						keyID:       "test-id",
						keyTime:     "test-time",
						keyDuration: int64(1),
						keyStatus:   200,
						keySize:     int64(10),
						keyHeader:   map[string]string{"key": "value"},
						keyBody:     "test-body",
					},
				},
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			val := tt.C().attr.journalKeyValues()
			testutil.Diff(t, tt.A().val, val)
		})
	}
}

func TestResponseAttrs_TagFunc(t *testing.T) {
	type condition struct {
		attr *responseAttrs
		tag  string
	}

	type action struct {
		val string
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	testAttr := &responseAttrs{
		typ:      "test-type",
		id:       "test-id",
		time:     "test-time",
		duration: 1,
		status:   200,
		size:     10,
		header:   map[string]string{"Key": "value"},
		body:     "test-body",
	}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"id",
			[]string{},
			[]string{},
			&condition{
				attr: testAttr,
				tag:  "id",
			},
			&action{
				val: "test-id",
			},
		),
		gen(
			"time",
			[]string{},
			[]string{},
			&condition{
				attr: testAttr,
				tag:  "time",
			},
			&action{
				val: "test-time",
			},
		),
		gen(
			"duration",
			[]string{},
			[]string{},
			&condition{
				attr: testAttr,
				tag:  "duration",
			},
			&action{
				val: "1",
			},
		),
		gen(
			"status",
			[]string{},
			[]string{},
			&condition{
				attr: testAttr,
				tag:  "status",
			},
			&action{
				val: "200",
			},
		),
		gen(
			"size",
			[]string{},
			[]string{},
			&condition{
				attr: testAttr,
				tag:  "size",
			},
			&action{
				val: "10",
			},
		),
		gen(
			"header",
			[]string{},
			[]string{},
			&condition{
				attr: testAttr,
				tag:  "header",
			},
			&action{
				val: "map[Key:value]",
			},
		),
		gen(
			"r.body",
			[]string{},
			[]string{},
			&condition{
				attr: testAttr,
				tag:  "body",
			},
			&action{
				val: "test-body",
			},
		),
		gen(
			"header.key",
			[]string{},
			[]string{},
			&condition{
				attr: testAttr,
				tag:  "header.key",
			},
			&action{
				val: "value",
			},
		),
		gen(
			"type",
			[]string{},
			[]string{},
			&condition{
				attr: testAttr,
				tag:  "type",
			},
			&action{
				val: "test-type",
			},
		),
		gen(
			"not.exist",
			[]string{},
			[]string{},
			&condition{
				attr: testAttr,
				tag:  "not.exist",
			},
			&action{
				val: "<undefined:not.exist>",
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			val := tt.C().attr.TagFunc(tt.C().tag)
			testutil.Diff(t, tt.A().val, string(val))
		})
	}
}

func TestNewLogID(t *testing.T) {
	type condition struct {
		count uint64
	}

	type action struct {
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"counter 1",
			[]string{},
			[]string{},
			&condition{
				count: 1,
			},
			&action{},
		),
		gen(
			"counter 99999",
			[]string{},
			[]string{},
			&condition{
				count: 99999,
			},
			&action{},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			counter.Store(tt.C().count)
			now := time.Now().Unix()

			id := newLogID()
			b, err := base64.URLEncoding.DecodeString(id)
			testutil.Diff(t, nil, err)

			expectSec := now & 0x00000000_ffffffff
			testutil.Diff(t, expectSec, int64(binary.BigEndian.Uint32(b[4:8])))

			expectCount := counter.Load() & 0x00ffffff_ffffffff
			testutil.Diff(t, expectCount, binary.BigEndian.Uint64(append([]byte{0}, b[8:]...)))
		})
	}
}
