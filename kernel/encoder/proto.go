package encoder

import (
	"github.com/aileron-gateway/aileron-gateway/kernel/er"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// MarshalProto marshal proto message into byte array.
// If nil is given as input, this function do nothing and return nil byte and nil error.
// When marshaling proto struct into a byte array, this function use the following options by default.
// That option can be replace by the second argument.
//
//	opt := proto.MarshalOptions{
//		AllowPartial: true,
//	}
func MarshalProto(in protoreflect.ProtoMessage, opt *proto.MarshalOptions) ([]byte, error) {
	if in == nil {
		return nil, nil
	}
	if opt == nil {
		opt = &proto.MarshalOptions{
			AllowPartial: true,
		}
	}
	b, err := opt.Marshal(in)
	if err != nil {
		return nil, (&er.Error{
			Package:     ErrPkg,
			Type:        ErrTypeProto,
			Description: ErrDscMarshal,
			Detail:      "from ProtoMessage to proto",
		}).Wrap(err)
	}
	return b, nil
}

// UnmarshalProto unmarshal proto byte array into a proto struct.
// If nil is given as input byte or target proto struct, this function do nothing and return nil error.
// This function use the following unmarshal options by default.
// It can be replaced by the second argument.
//
//	opt = &proto.UnmarshalOptions{
//		Merge:          true,
//		AllowPartial:   true,
//		DiscardUnknown: false,
//	}
func UnmarshalProto(in []byte, into protoreflect.ProtoMessage, opt *proto.UnmarshalOptions) error {
	if into == nil {
		return nil
	}
	if opt == nil {
		opt = &proto.UnmarshalOptions{
			Merge:          true,
			AllowPartial:   true,
			DiscardUnknown: false,
		}
	}
	if err := opt.Unmarshal(in, into); err != nil {
		return (&er.Error{
			Package:     ErrPkg,
			Type:        ErrTypeProto,
			Description: ErrDscUnmarshal,
			Detail:      "from proto to ProtoMessage",
		}).Wrap(err)
	}
	return nil
}

// UnmarshalProtoFromJSON unmarshal json byte array into proto struct.
// If nil is given as byte array or as target proto struct,
// this function do nothing and return nil error.
// If unmarshal option is nil, the following options are used by default.
//
//	opt = &protojson.UnmarshalOptions{
//		AllowPartial:   true,
//		DiscardUnknown: false,
//	}
func UnmarshalProtoFromJSON(in []byte, into protoreflect.ProtoMessage, opt *protojson.UnmarshalOptions) error {
	if into == nil {
		return nil
	}
	if opt == nil {
		opt = &protojson.UnmarshalOptions{
			AllowPartial:   true,
			DiscardUnknown: false,
		}
	}
	if err := opt.Unmarshal(in, into); err != nil {
		return (&er.Error{
			Package:     ErrPkg,
			Type:        ErrTypeProto,
			Description: ErrDscUnmarshal,
			Detail:      string(addLineNumber(in)),
		}).Wrap(err)
	}
	return nil
}

// MarshalProtoToJSON marshal proto struct into json byte array.
// If nil is given as proto struct, this function do nothing and return nil byte and nil error.
// If marshal option is nil, the following options are used by default.
//
//	opt = &protojson.MarshalOptions{
//		Multiline:       false,
//		Indent:          "",
//		AllowPartial:    true,
//		UseProtoNames:   false,
//		UseEnumNumbers:  false,
//		EmitUnpopulated: false,
//	}
func MarshalProtoToJSON(in protoreflect.ProtoMessage, opt *protojson.MarshalOptions) ([]byte, error) {
	if in == nil {
		return nil, nil
	}
	if opt == nil {
		opt = &protojson.MarshalOptions{
			Multiline:       false,
			Indent:          "",
			AllowPartial:    true,
			UseProtoNames:   false,
			UseEnumNumbers:  false,
			EmitUnpopulated: false,
		}
	}
	b, err := opt.Marshal(in)
	if err != nil {
		return nil, (&er.Error{
			Package:     ErrPkg,
			Type:        ErrTypeProto,
			Description: ErrDscMarshal,
		}).Wrap(err)
	}
	return b, nil
}

// UnmarshalProtoFromYAML unmarshal yaml byte array into proto struct.
// If nil is given as byte array or as target proto struct,
// this function do nothing and return nil error.
// If unmarshal option is nil, the following options are used by default.
//
//	opt = &protojson.UnmarshalOptions{
//		AllowPartial:   true,
//		DiscardUnknown: false,
//	}
func UnmarshalProtoFromYAML(in []byte, into protoreflect.ProtoMessage, opt *protojson.UnmarshalOptions) error {
	if into == nil {
		return nil
	}

	// Temporarily convert yaml to json.
	tmp := &map[string]any{}
	if err := UnmarshalYAML(in, tmp); err != nil {
		return err
	}

	b, err := MarshalJSON(tmp)
	if err != nil {
		return err
	}

	return UnmarshalProtoFromJSON(b, into, opt)
}

// MarshalProtoToYAML marshal proto struct into yaml byte array.
// If nil is given as proto struct, this function do nothing and return nil byte and nil error.
// If marshal option is nil, the following options are used by default.
//
//	opt = &protojson.MarshalOptions{
//		Multiline:       false,
//		Indent:          "",
//		AllowPartial:    true,
//		UseProtoNames:   false,
//		UseEnumNumbers:  false,
//		EmitUnpopulated: false,
//	}
func MarshalProtoToYAML(in protoreflect.ProtoMessage, opt *protojson.MarshalOptions) ([]byte, error) {
	if in == nil {
		return nil, nil
	}
	// temporarily convert proto to json.
	b, err := MarshalProtoToJSON(in, opt)
	if err != nil {
		return nil, err // Return err as-is.
	}

	tmp := &map[string]any{}
	if err := UnmarshalJSON(b, tmp); err != nil {
		return nil, err // Return err as-is.
	}

	return MarshalYAML(tmp)
}
