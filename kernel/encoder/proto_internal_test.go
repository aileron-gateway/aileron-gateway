package encoder

import (
	"testing"

	k "github.com/aileron-gateway/aileron-gateway/apis/kernel"
	"github.com/aileron-gateway/aileron-gateway/kernel/er"
	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
	"github.com/google/go-cmp/cmp/cmpopts"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func TestMarshalProto(t *testing.T) {
	type condition struct {
		in protoreflect.ProtoMessage
	}

	type action struct {
		out string
		err error
	}

	cndNil := "input nil"
	cndInvalidVal := "input invalid value"
	actCheckExpected := "expected value returned"
	actCheckNoError := "no error"
	actCheckError := "expected error returned"

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	tb.Condition(cndNil, "give a valid encoded string")
	tb.Condition(cndInvalidVal, "give an invalid value which will result in an error")
	tb.Action(actCheckExpected, "Check that an expected value returned")
	tb.Action(actCheckNoError, "Check that returned error is nil")
	tb.Action(actCheckError, "Check that an expected error was returned")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"encode struct",
			[]string{},
			[]string{actCheckExpected, actCheckNoError},
			&condition{
				in: &k.Metadata{
					Name: "John Doe",
				},
			},
			&action{
				out: "\n\bJohn Doe",
				err: nil,
			},
		),
		gen(
			"nil",
			[]string{cndNil},
			[]string{actCheckExpected, actCheckNoError},
			&condition{
				in: nil,
			},
			&action{
				out: "",
				err: nil,
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			out, err := MarshalProto(tt.C().in, nil)
			testutil.Diff(t, tt.A().err, err, cmpopts.EquateErrors())
			testutil.Diff(t, tt.A().out, string(out))
		})
	}
}

func TestUnmarshalProto(t *testing.T) {
	type condition struct {
		in   string
		into protoreflect.ProtoMessage
	}

	type action struct {
		result protoreflect.ProtoMessage
		err    error
	}

	cndNil := "input nil"
	cndInvalidVal := "input invalid value"
	actCheckExpected := "expected value returned"
	actCheckNoError := "no error"
	actCheckError := "expected error returned"

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	tb.Condition(cndNil, "give a nil value as an input")
	tb.Condition(cndInvalidVal, "give an invalid value which will result in an error")
	tb.Action(actCheckExpected, "Check that an expected value returned")
	tb.Action(actCheckNoError, "Check that returned error is nil")
	tb.Action(actCheckError, "Check that an expected error was returned")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"decode proto string",
			[]string{},
			[]string{actCheckExpected, actCheckNoError},
			&condition{
				in:   "\n\bJohn Doe",
				into: &k.Metadata{},
			},
			&action{
				result: &k.Metadata{
					Name: "John Doe",
				},
				err: nil,
			},
		),
		gen(
			"nil",
			[]string{cndNil},
			[]string{actCheckExpected, actCheckNoError},
			&condition{
				in:   "",
				into: nil,
			},
			&action{
				result: nil,
				err:    nil,
			},
		),
		gen(
			"failed to marshal",
			[]string{cndInvalidVal},
			[]string{actCheckError},
			&condition{
				in:   "Invalid Proto",
				into: &k.Metadata{},
			},
			&action{
				result: &k.Metadata{},
				err: &er.Error{
					Package:     ErrPkg,
					Type:        ErrTypeProto,
					Description: ErrDscUnmarshal,
				},
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			err := UnmarshalProto([]byte(tt.C().in), tt.C().into, nil)
			testutil.Diff(t, tt.A().err, err, cmpopts.EquateErrors())
			testutil.Diff(t, tt.A().result, tt.C().into, cmpopts.IgnoreUnexported(k.Metadata{}))
		})
	}
}

func TestUnmarshalProtoFromJSON(t *testing.T) {
	type condition struct {
		in   string
		into protoreflect.ProtoMessage
	}

	type action struct {
		result protoreflect.ProtoMessage
		err    error
	}

	cndNil := "input nil"
	cndInvalidVal := "input invalid value"
	actCheckExpected := "expected value returned"
	actCheckNoError := "no error"
	actCheckError := "expected error returned"

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	tb.Condition(cndNil, "give a nil value as an input")
	tb.Condition(cndInvalidVal, "give an invalid value which will result in an error")
	tb.Action(actCheckExpected, "Check that an expected value returned")
	tb.Action(actCheckNoError, "Check that returned error is nil")
	tb.Action(actCheckError, "Check that an expected error was returned")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"decode json string",
			[]string{},
			[]string{actCheckExpected, actCheckNoError},
			&condition{
				in:   `{"name":"John Doe"}`,
				into: &k.Metadata{Name: "John Doe"},
			},
			&action{
				result: &k.Metadata{
					Name: "John Doe",
				},
				err: nil,
			},
		),
		gen(
			"nil",
			[]string{cndNil},
			[]string{actCheckExpected, actCheckNoError},
			&condition{
				in:   "",
				into: nil,
			},
			&action{
				result: nil,
				err:    nil,
			},
		),
		gen(
			"failed to marshal",
			[]string{cndInvalidVal},
			[]string{actCheckError},
			&condition{
				in:   `{Invalid JSON}`,
				into: &k.Metadata{},
			},
			&action{
				result: &k.Metadata{},
				err: &er.Error{
					Package:     ErrPkg,
					Type:        ErrTypeProto,
					Description: ErrDscUnmarshal,
				},
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			err := UnmarshalProtoFromJSON([]byte(tt.C().in), tt.C().into, nil)
			testutil.Diff(t, tt.A().err, err, cmpopts.EquateErrors())
			testutil.Diff(t, tt.A().result, tt.C().into, cmpopts.IgnoreUnexported(k.Metadata{}))
		})
	}
}

func TestMarshalProtoToJSON(t *testing.T) {
	type condition struct {
		in  protoreflect.ProtoMessage
		opt *protojson.MarshalOptions
	}

	type action struct {
		out string
		err error
	}

	cndNil := "input nil"
	cndInvalidVal := "input invalid value"
	actCheckExpected := "expected value returned"
	actCheckNoError := "no error"
	actCheckError := "expected error returned"

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	tb.Condition(cndNil, "give a valid encoded string")
	tb.Condition(cndInvalidVal, "give an invalid value which will result in an error")
	tb.Action(actCheckExpected, "Check that an expected value returned")
	tb.Action(actCheckNoError, "Check that returned error is nil")
	tb.Action(actCheckError, "Check that an expected error was returned")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"encode struct",
			[]string{},
			[]string{actCheckExpected, actCheckNoError},
			&condition{
				in: &k.Metadata{
					Name: "John Doe",
				},
			},
			&action{
				out: `{"name":"John Doe"}`,
				err: nil,
			},
		),
		gen(
			"nil",
			[]string{cndNil},
			[]string{actCheckExpected, actCheckNoError},
			&condition{
				in: nil,
			},
			&action{
				out: "",
				err: nil,
			},
		),
		gen(
			"invalid option",
			[]string{},
			[]string{actCheckExpected, actCheckError},
			&condition{
				in:  &k.Metadata{},
				opt: &protojson.MarshalOptions{Multiline: true, Indent: "\n"},
			},
			&action{
				out: "",
				err: &er.Error{
					Package:     ErrPkg,
					Type:        ErrTypeProto,
					Description: ErrDscMarshal,
				},
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			out, err := MarshalProtoToJSON(tt.C().in, tt.C().opt)
			testutil.Diff(t, tt.A().err, err, cmpopts.EquateErrors())
			testutil.Diff(t, tt.A().out, string(out))
		})
	}
}

func TestUnmarshalProtoFromYAML(t *testing.T) {
	type condition struct {
		in   string
		into protoreflect.ProtoMessage
	}

	type action struct {
		result protoreflect.ProtoMessage
		err    error
	}

	cndNil := "input nil"
	cndInvalidVal := "input invalid value"
	actCheckExpected := "expected value returned"
	actCheckNoError := "no error"
	actCheckError := "expected error returned"

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	tb.Condition(cndNil, "give a nil value as an input")
	tb.Condition(cndInvalidVal, "give an invalid value which will result in an error")
	tb.Action(actCheckExpected, "Check that an expected value returned")
	tb.Action(actCheckNoError, "Check that returned error is nil")
	tb.Action(actCheckError, "Check that an expected error was returned")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"decode yaml string",
			[]string{},
			[]string{actCheckExpected, actCheckNoError},
			&condition{
				in:   "name: John Doe\n",
				into: &k.Metadata{},
			},
			&action{
				result: &k.Metadata{
					Name: "John Doe",
				},
				err: nil,
			},
		),
		gen(
			"nil",
			[]string{cndNil},
			[]string{actCheckExpected, actCheckNoError},
			&condition{
				in:   "",
				into: nil,
			},
			&action{
				result: nil,
				err:    nil,
			},
		),
		gen(
			"failed to unmarshal",
			[]string{cndInvalidVal},
			[]string{actCheckError},
			&condition{
				in:   "Invalid:YAML",
				into: &k.Metadata{},
			},
			&action{
				result: &k.Metadata{},
				err: &er.Error{
					Package:     ErrPkg,
					Type:        ErrTypeYaml,
					Description: ErrDscUnmarshal,
				},
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			err := UnmarshalProtoFromYAML([]byte(tt.C().in), tt.C().into, nil)
			testutil.Diff(t, tt.A().err, err, cmpopts.EquateErrors())
			testutil.Diff(t, tt.A().result, tt.C().into, cmpopts.IgnoreUnexported(k.Metadata{}))
		})
	}
}

func TestMarshalProtoToYAML(t *testing.T) {
	type condition struct {
		in  protoreflect.ProtoMessage
		opt *protojson.MarshalOptions
	}

	type action struct {
		out string
		err error
	}

	cndNil := "input nil"
	cndInvalidVal := "input invalid value"
	actCheckExpected := "expected value returned"
	actCheckNoError := "no error"
	actCheckError := "expected error returned"

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	tb.Condition(cndNil, "give a valid encoded string")
	tb.Condition(cndInvalidVal, "give an invalid value which will result in an error")
	tb.Action(actCheckExpected, "Check that an expected value returned")
	tb.Action(actCheckNoError, "Check that returned error is nil")
	tb.Action(actCheckError, "Check that an expected error was returned")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"encode struct",
			[]string{},
			[]string{actCheckExpected, actCheckNoError},
			&condition{
				in: &k.Metadata{
					Name: "John Doe",
				},
			},
			&action{
				out: "name: John Doe\n",
				err: nil,
			},
		),
		gen(
			"nil",
			[]string{cndNil},
			[]string{actCheckExpected, actCheckNoError},
			&condition{
				in: nil,
			},
			&action{
				out: "",
				err: nil,
			},
		),
		gen(
			"invalid option",
			[]string{},
			[]string{actCheckExpected, actCheckError},
			&condition{
				in:  &k.Metadata{},
				opt: &protojson.MarshalOptions{Multiline: true, Indent: "\n"},
			},
			&action{
				out: "",
				err: &er.Error{
					Package:     ErrPkg,
					Type:        ErrTypeProto,
					Description: ErrDscMarshal,
				},
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			out, err := MarshalProtoToYAML(tt.C().in, tt.C().opt)
			testutil.Diff(t, tt.A().err, err, cmpopts.EquateErrors())
			testutil.Diff(t, tt.A().out, string(out))
		})
	}
}
