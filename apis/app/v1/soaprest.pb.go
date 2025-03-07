// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.4
// 	protoc        v5.29.0
// source: app/v1/middleware/soaprest.proto

package v1

import (
	_ "buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go/buf/validate"
	kernel "github.com/aileron-gateway/aileron-gateway/apis/kernel"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// SOAPRESTMiddleware is the definition of the SOAPRESTMiddleware object.
// SOAPRESTMiddleware implements interface of the middleware.
type SOAPRESTMiddleware struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// [REQUIRED]
	// APIVersion is the defined version of the middleware.
	// This value must be "app/v1".
	APIVersion string `protobuf:"bytes,1,opt,name=APIVersion,json=apiVersion,proto3" json:"APIVersion,omitempty"`
	// [REQUIRED]
	// Kind is the kind of this object.
	// This value must be "SOAPRESTMiddleware".
	Kind string `protobuf:"bytes,2,opt,name=Kind,json=kind,proto3" json:"Kind,omitempty"`
	// [OPTIONAL]
	// Metadata is the metadata of the http logger object.
	// If not set, both name and namespace in the metadata
	// are treated as "default".
	Metadata *kernel.Metadata `protobuf:"bytes,3,opt,name=Metadata,json=metadata,proto3" json:"Metadata,omitempty"`
	// [OPTIONAL]
	// Spec is the specification of the middleware.
	// Default values are used when nothing is set.
	Spec          *SOAPRESTMiddlewareSpec `protobuf:"bytes,4,opt,name=Spec,json=spec,proto3" json:"Spec,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *SOAPRESTMiddleware) Reset() {
	*x = SOAPRESTMiddleware{}
	mi := &file_app_v1_middleware_soaprest_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SOAPRESTMiddleware) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SOAPRESTMiddleware) ProtoMessage() {}

func (x *SOAPRESTMiddleware) ProtoReflect() protoreflect.Message {
	mi := &file_app_v1_middleware_soaprest_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SOAPRESTMiddleware.ProtoReflect.Descriptor instead.
func (*SOAPRESTMiddleware) Descriptor() ([]byte, []int) {
	return file_app_v1_middleware_soaprest_proto_rawDescGZIP(), []int{0}
}

func (x *SOAPRESTMiddleware) GetAPIVersion() string {
	if x != nil {
		return x.APIVersion
	}
	return ""
}

func (x *SOAPRESTMiddleware) GetKind() string {
	if x != nil {
		return x.Kind
	}
	return ""
}

func (x *SOAPRESTMiddleware) GetMetadata() *kernel.Metadata {
	if x != nil {
		return x.Metadata
	}
	return nil
}

func (x *SOAPRESTMiddleware) GetSpec() *SOAPRESTMiddlewareSpec {
	if x != nil {
		return x.Spec
	}
	return nil
}

// SOAPRESTMiddlewareSpec is the status of the SOAPRESTMiddleware object.
// Values are managed by the application and therefore should not be set by users.
type SOAPRESTMiddlewareSpec struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// [REQUIRED]
	// Matcher is a matcher which matches to defined patterns.
	// Default is not set.
	Matcher *kernel.MatcherSpec `protobuf:"bytes,1,opt,name=Matcher,json=matcher,proto3" json:"Matcher,omitempty"`
	// [OPTIONAL]
	// ErrorHandler is the reference to a ErrorHandler object.
	// Referred object must implement ErrorHandler interface.
	// Default error handler is used when not set.
	ErrorHandler *kernel.Reference `protobuf:"bytes,2,opt,name=ErrorHandler,json=errorHandler,proto3" json:"ErrorHandler,omitempty"`
	// [OPTIONAL]
	// AttributeKey is the key used in JSON to hold attributes found in XML.
	// AttributeKey is also used during the conversion from JSON to XML.
	// Default value is "@attribute".
	AttributeKey string `protobuf:"bytes,3,opt,name=AttributeKey,json=attributeKey,proto3" json:"AttributeKey,omitempty"`
	// [OPTIONAL]
	// TextKey is the key used in JSON to hold text nodes found in XML.
	// TextKey is also used during the conversion from JSON to XML.
	// Default value is "#text".
	TextKey string `protobuf:"bytes,4,opt,name=TextKey,json=textKey,proto3" json:"TextKey,omitempty"`
	// [OPTIONAL]
	// NamespaceKey is the key used in JSON to hold namespaces found in XML.
	// NamespaceKey is also used during the conversion from JSON to XML.
	// Default value is "_namespace".
	NamespaceKey string `protobuf:"bytes,5,opt,name=NamespaceKey,json=namespaceKey,proto3" json:"NamespaceKey,omitempty"`
	// [OPTIONAL]
	// ArrayKey is the key used to represent array elements during the conversion from JSON to XML.
	// Default value is "item".
	ArrayKey string `protobuf:"bytes,6,opt,name=ArrayKey,json=arrayKey,proto3" json:"ArrayKey,omitempty"`
	// [OPTIONAL]
	// SeparatorChar is the symbol used to distinguish between namespaces and keys.
	// Default value is ":".
	SeparatorChar string `protobuf:"bytes,7,opt,name=SeparatorChar,json=separatorChar,proto3" json:"SeparatorChar,omitempty"`
	// [OPTIONAL]
	// ExtractStringElement is a boolean value used to extract elements enclosed in double quotes within XML.
	// When this option is set to false, double quotes will be escaped.
	// Specifically, the processing occurs as follows:
	//
	//	""     → "\"\""
	//	"100"  → "\"100\""
	//	"true" → "\"true\""
	//
	// When its set to true, the escaping process will not be performed, and elements are extracted as follows:
	//
	//	""     → ""
	//	"100"  → "100"
	//	"true" → "true"
	//
	// Default value is "false".
	ExtractStringElement bool `protobuf:"varint,8,opt,name=ExtractStringElement,json=extractStringElement,proto3" json:"ExtractStringElement,omitempty"`
	// [OPTIONAL]
	// ExtractBooleanElement is a boolean value used to extract elements enclosed in boolean tags within XML.
	// When this option is set to false, do not extract boolean elements.
	//
	//	"true"  → "true"
	//	"false" → "false"
	//
	// When its set to true, the extract process will be performed, and elements are extracted as follows:
	//
	//	"true"  → true
	//	"false" → false
	//
	// Default value is "false".
	ExtractBooleanElement bool `protobuf:"varint,9,opt,name=ExtractBooleanElement,json=extractBooleanElement,proto3" json:"ExtractBooleanElement,omitempty"`
	// [OPTIONAL]
	// ExtractIntegerElement is a boolean value used to extract integer elements within XML.
	// When this option is set to false, do not extract integer elements.
	//
	//	"0"   → "0"
	//	"100" → "100"
	//
	// When its set to true, the extract process will be performed, and elements are extracted as follows:
	//
	//	"0"   → 0
	//	"100" → 100
	//
	// The precision is based on the 64-bit integer type.
	// Default value is "false".
	ExtractIntegerElement bool `protobuf:"varint,10,opt,name=ExtractIntegerElement,json=extractIntegerElement,proto3" json:"ExtractIntegerElement,omitempty"`
	// [OPTIONAL]
	// ExtractFloatElement is a value used to extract float elements within XML.
	// When this option is set to false, do not extract float elements.
	//
	//	"0.1"   → "0.1"
	//	"3.14"  → "3.14"
	//
	// When its set to true, the extract process will be performed, and elements are extracted as follows:
	//
	//	"0.1"  → 0.1
	//	"3.14" → 3.14
	//
	// The precision is based on the 64-bit float type.
	// Default value is "false".
	ExtractFloatElement bool `protobuf:"varint,11,opt,name=ExtractFloatElement,json=extractFloatElement,proto3" json:"ExtractFloatElement,omitempty"`
	unknownFields       protoimpl.UnknownFields
	sizeCache           protoimpl.SizeCache
}

func (x *SOAPRESTMiddlewareSpec) Reset() {
	*x = SOAPRESTMiddlewareSpec{}
	mi := &file_app_v1_middleware_soaprest_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SOAPRESTMiddlewareSpec) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SOAPRESTMiddlewareSpec) ProtoMessage() {}

func (x *SOAPRESTMiddlewareSpec) ProtoReflect() protoreflect.Message {
	mi := &file_app_v1_middleware_soaprest_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SOAPRESTMiddlewareSpec.ProtoReflect.Descriptor instead.
func (*SOAPRESTMiddlewareSpec) Descriptor() ([]byte, []int) {
	return file_app_v1_middleware_soaprest_proto_rawDescGZIP(), []int{1}
}

func (x *SOAPRESTMiddlewareSpec) GetMatcher() *kernel.MatcherSpec {
	if x != nil {
		return x.Matcher
	}
	return nil
}

func (x *SOAPRESTMiddlewareSpec) GetErrorHandler() *kernel.Reference {
	if x != nil {
		return x.ErrorHandler
	}
	return nil
}

func (x *SOAPRESTMiddlewareSpec) GetAttributeKey() string {
	if x != nil {
		return x.AttributeKey
	}
	return ""
}

func (x *SOAPRESTMiddlewareSpec) GetTextKey() string {
	if x != nil {
		return x.TextKey
	}
	return ""
}

func (x *SOAPRESTMiddlewareSpec) GetNamespaceKey() string {
	if x != nil {
		return x.NamespaceKey
	}
	return ""
}

func (x *SOAPRESTMiddlewareSpec) GetArrayKey() string {
	if x != nil {
		return x.ArrayKey
	}
	return ""
}

func (x *SOAPRESTMiddlewareSpec) GetSeparatorChar() string {
	if x != nil {
		return x.SeparatorChar
	}
	return ""
}

func (x *SOAPRESTMiddlewareSpec) GetExtractStringElement() bool {
	if x != nil {
		return x.ExtractStringElement
	}
	return false
}

func (x *SOAPRESTMiddlewareSpec) GetExtractBooleanElement() bool {
	if x != nil {
		return x.ExtractBooleanElement
	}
	return false
}

func (x *SOAPRESTMiddlewareSpec) GetExtractIntegerElement() bool {
	if x != nil {
		return x.ExtractIntegerElement
	}
	return false
}

func (x *SOAPRESTMiddlewareSpec) GetExtractFloatElement() bool {
	if x != nil {
		return x.ExtractFloatElement
	}
	return false
}

var File_app_v1_middleware_soaprest_proto protoreflect.FileDescriptor

var file_app_v1_middleware_soaprest_proto_rawDesc = string([]byte{
	0x0a, 0x20, 0x61, 0x70, 0x70, 0x2f, 0x76, 0x31, 0x2f, 0x6d, 0x69, 0x64, 0x64, 0x6c, 0x65, 0x77,
	0x61, 0x72, 0x65, 0x2f, 0x73, 0x6f, 0x61, 0x70, 0x72, 0x65, 0x73, 0x74, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x06, 0x61, 0x70, 0x70, 0x2e, 0x76, 0x31, 0x1a, 0x1b, 0x62, 0x75, 0x66, 0x2f,
	0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2f, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x15, 0x6b, 0x65, 0x72, 0x6e, 0x65, 0x6c, 0x2f,
	0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x14,
	0x6b, 0x65, 0x72, 0x6e, 0x65, 0x6c, 0x2f, 0x74, 0x78, 0x74, 0x75, 0x74, 0x69, 0x6c, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x22, 0xd4, 0x01, 0x0a, 0x12, 0x53, 0x4f, 0x41, 0x50, 0x52, 0x45, 0x53,
	0x54, 0x4d, 0x69, 0x64, 0x64, 0x6c, 0x65, 0x77, 0x61, 0x72, 0x65, 0x12, 0x2d, 0x0a, 0x0a, 0x41,
	0x50, 0x49, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42,
	0x0d, 0xba, 0x48, 0x0a, 0x72, 0x08, 0x0a, 0x06, 0x61, 0x70, 0x70, 0x2f, 0x76, 0x31, 0x52, 0x0a,
	0x61, 0x70, 0x69, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x2d, 0x0a, 0x04, 0x4b, 0x69,
	0x6e, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x42, 0x19, 0xba, 0x48, 0x16, 0x72, 0x14, 0x0a,
	0x12, 0x53, 0x4f, 0x41, 0x50, 0x52, 0x45, 0x53, 0x54, 0x4d, 0x69, 0x64, 0x64, 0x6c, 0x65, 0x77,
	0x61, 0x72, 0x65, 0x52, 0x04, 0x6b, 0x69, 0x6e, 0x64, 0x12, 0x2c, 0x0a, 0x08, 0x4d, 0x65, 0x74,
	0x61, 0x64, 0x61, 0x74, 0x61, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x6b, 0x65,
	0x72, 0x6e, 0x65, 0x6c, 0x2e, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x52, 0x08, 0x6d,
	0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x12, 0x32, 0x0a, 0x04, 0x53, 0x70, 0x65, 0x63, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x61, 0x70, 0x70, 0x2e, 0x76, 0x31, 0x2e, 0x53,
	0x4f, 0x41, 0x50, 0x52, 0x45, 0x53, 0x54, 0x4d, 0x69, 0x64, 0x64, 0x6c, 0x65, 0x77, 0x61, 0x72,
	0x65, 0x53, 0x70, 0x65, 0x63, 0x52, 0x04, 0x73, 0x70, 0x65, 0x63, 0x22, 0xa4, 0x05, 0x0a, 0x16,
	0x53, 0x4f, 0x41, 0x50, 0x52, 0x45, 0x53, 0x54, 0x4d, 0x69, 0x64, 0x64, 0x6c, 0x65, 0x77, 0x61,
	0x72, 0x65, 0x53, 0x70, 0x65, 0x63, 0x12, 0x35, 0x0a, 0x07, 0x4d, 0x61, 0x74, 0x63, 0x68, 0x65,
	0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x13, 0x2e, 0x6b, 0x65, 0x72, 0x6e, 0x65, 0x6c,
	0x2e, 0x4d, 0x61, 0x74, 0x63, 0x68, 0x65, 0x72, 0x53, 0x70, 0x65, 0x63, 0x42, 0x06, 0xba, 0x48,
	0x03, 0xc8, 0x01, 0x01, 0x52, 0x07, 0x6d, 0x61, 0x74, 0x63, 0x68, 0x65, 0x72, 0x12, 0x35, 0x0a,
	0x0c, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x48, 0x61, 0x6e, 0x64, 0x6c, 0x65, 0x72, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x11, 0x2e, 0x6b, 0x65, 0x72, 0x6e, 0x65, 0x6c, 0x2e, 0x52, 0x65, 0x66,
	0x65, 0x72, 0x65, 0x6e, 0x63, 0x65, 0x52, 0x0c, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x48, 0x61, 0x6e,
	0x64, 0x6c, 0x65, 0x72, 0x12, 0x4c, 0x0a, 0x0c, 0x41, 0x74, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74,
	0x65, 0x4b, 0x65, 0x79, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x42, 0x28, 0xba, 0x48, 0x25, 0x72,
	0x23, 0x32, 0x21, 0x5e, 0x5b, 0x61, 0x2d, 0x7a, 0x41, 0x2d, 0x5a, 0x3a, 0x2e, 0x5c, 0x5f, 0x7c,
	0x2f, 0x23, 0x40, 0x5d, 0x5b, 0x61, 0x2d, 0x7a, 0x41, 0x2d, 0x5a, 0x30, 0x2d, 0x39, 0x5f, 0x2e,
	0x2d, 0x5d, 0x2a, 0x24, 0x52, 0x0c, 0x61, 0x74, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x65, 0x4b,
	0x65, 0x79, 0x12, 0x42, 0x0a, 0x07, 0x54, 0x65, 0x78, 0x74, 0x4b, 0x65, 0x79, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x09, 0x42, 0x28, 0xba, 0x48, 0x25, 0x72, 0x23, 0x32, 0x21, 0x5e, 0x5b, 0x61, 0x2d,
	0x7a, 0x41, 0x2d, 0x5a, 0x3a, 0x2e, 0x5c, 0x5f, 0x7c, 0x2f, 0x23, 0x40, 0x5d, 0x5b, 0x61, 0x2d,
	0x7a, 0x41, 0x2d, 0x5a, 0x30, 0x2d, 0x39, 0x5f, 0x2e, 0x2d, 0x5d, 0x2a, 0x24, 0x52, 0x07, 0x74,
	0x65, 0x78, 0x74, 0x4b, 0x65, 0x79, 0x12, 0x4c, 0x0a, 0x0c, 0x4e, 0x61, 0x6d, 0x65, 0x73, 0x70,
	0x61, 0x63, 0x65, 0x4b, 0x65, 0x79, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x42, 0x28, 0xba, 0x48,
	0x25, 0x72, 0x23, 0x32, 0x21, 0x5e, 0x5b, 0x61, 0x2d, 0x7a, 0x41, 0x2d, 0x5a, 0x3a, 0x2e, 0x5c,
	0x5f, 0x7c, 0x2f, 0x23, 0x40, 0x5d, 0x5b, 0x61, 0x2d, 0x7a, 0x41, 0x2d, 0x5a, 0x30, 0x2d, 0x39,
	0x5f, 0x2e, 0x2d, 0x5d, 0x2a, 0x24, 0x52, 0x0c, 0x6e, 0x61, 0x6d, 0x65, 0x73, 0x70, 0x61, 0x63,
	0x65, 0x4b, 0x65, 0x79, 0x12, 0x44, 0x0a, 0x08, 0x41, 0x72, 0x72, 0x61, 0x79, 0x4b, 0x65, 0x79,
	0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x42, 0x28, 0xba, 0x48, 0x25, 0x72, 0x23, 0x32, 0x21, 0x5e,
	0x5b, 0x61, 0x2d, 0x7a, 0x41, 0x2d, 0x5a, 0x3a, 0x2e, 0x5c, 0x5f, 0x7c, 0x2f, 0x23, 0x40, 0x5d,
	0x5b, 0x61, 0x2d, 0x7a, 0x41, 0x2d, 0x5a, 0x30, 0x2d, 0x39, 0x5f, 0x2e, 0x2d, 0x5d, 0x2a, 0x24,
	0x52, 0x08, 0x61, 0x72, 0x72, 0x61, 0x79, 0x4b, 0x65, 0x79, 0x12, 0x24, 0x0a, 0x0d, 0x53, 0x65,
	0x70, 0x61, 0x72, 0x61, 0x74, 0x6f, 0x72, 0x43, 0x68, 0x61, 0x72, 0x18, 0x07, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0d, 0x73, 0x65, 0x70, 0x61, 0x72, 0x61, 0x74, 0x6f, 0x72, 0x43, 0x68, 0x61, 0x72,
	0x12, 0x32, 0x0a, 0x14, 0x45, 0x78, 0x74, 0x72, 0x61, 0x63, 0x74, 0x53, 0x74, 0x72, 0x69, 0x6e,
	0x67, 0x45, 0x6c, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x18, 0x08, 0x20, 0x01, 0x28, 0x08, 0x52, 0x14,
	0x65, 0x78, 0x74, 0x72, 0x61, 0x63, 0x74, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x45, 0x6c, 0x65,
	0x6d, 0x65, 0x6e, 0x74, 0x12, 0x34, 0x0a, 0x15, 0x45, 0x78, 0x74, 0x72, 0x61, 0x63, 0x74, 0x42,
	0x6f, 0x6f, 0x6c, 0x65, 0x61, 0x6e, 0x45, 0x6c, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x18, 0x09, 0x20,
	0x01, 0x28, 0x08, 0x52, 0x15, 0x65, 0x78, 0x74, 0x72, 0x61, 0x63, 0x74, 0x42, 0x6f, 0x6f, 0x6c,
	0x65, 0x61, 0x6e, 0x45, 0x6c, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x12, 0x34, 0x0a, 0x15, 0x45, 0x78,
	0x74, 0x72, 0x61, 0x63, 0x74, 0x49, 0x6e, 0x74, 0x65, 0x67, 0x65, 0x72, 0x45, 0x6c, 0x65, 0x6d,
	0x65, 0x6e, 0x74, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x08, 0x52, 0x15, 0x65, 0x78, 0x74, 0x72, 0x61,
	0x63, 0x74, 0x49, 0x6e, 0x74, 0x65, 0x67, 0x65, 0x72, 0x45, 0x6c, 0x65, 0x6d, 0x65, 0x6e, 0x74,
	0x12, 0x30, 0x0a, 0x13, 0x45, 0x78, 0x74, 0x72, 0x61, 0x63, 0x74, 0x46, 0x6c, 0x6f, 0x61, 0x74,
	0x45, 0x6c, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x08, 0x52, 0x13, 0x65,
	0x78, 0x74, 0x72, 0x61, 0x63, 0x74, 0x46, 0x6c, 0x6f, 0x61, 0x74, 0x45, 0x6c, 0x65, 0x6d, 0x65,
	0x6e, 0x74, 0x42, 0x38, 0x5a, 0x36, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d,
	0x2f, 0x61, 0x69, 0x6c, 0x65, 0x72, 0x6f, 0x6e, 0x2d, 0x67, 0x61, 0x74, 0x65, 0x77, 0x61, 0x79,
	0x2f, 0x61, 0x69, 0x6c, 0x65, 0x72, 0x6f, 0x6e, 0x2d, 0x67, 0x61, 0x74, 0x65, 0x77, 0x61, 0x79,
	0x2f, 0x61, 0x70, 0x69, 0x73, 0x2f, 0x61, 0x70, 0x70, 0x2f, 0x76, 0x31, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
})

var (
	file_app_v1_middleware_soaprest_proto_rawDescOnce sync.Once
	file_app_v1_middleware_soaprest_proto_rawDescData []byte
)

func file_app_v1_middleware_soaprest_proto_rawDescGZIP() []byte {
	file_app_v1_middleware_soaprest_proto_rawDescOnce.Do(func() {
		file_app_v1_middleware_soaprest_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_app_v1_middleware_soaprest_proto_rawDesc), len(file_app_v1_middleware_soaprest_proto_rawDesc)))
	})
	return file_app_v1_middleware_soaprest_proto_rawDescData
}

var file_app_v1_middleware_soaprest_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_app_v1_middleware_soaprest_proto_goTypes = []any{
	(*SOAPRESTMiddleware)(nil),     // 0: app.v1.SOAPRESTMiddleware
	(*SOAPRESTMiddlewareSpec)(nil), // 1: app.v1.SOAPRESTMiddlewareSpec
	(*kernel.Metadata)(nil),        // 2: kernel.Metadata
	(*kernel.MatcherSpec)(nil),     // 3: kernel.MatcherSpec
	(*kernel.Reference)(nil),       // 4: kernel.Reference
}
var file_app_v1_middleware_soaprest_proto_depIdxs = []int32{
	2, // 0: app.v1.SOAPRESTMiddleware.Metadata:type_name -> kernel.Metadata
	1, // 1: app.v1.SOAPRESTMiddleware.Spec:type_name -> app.v1.SOAPRESTMiddlewareSpec
	3, // 2: app.v1.SOAPRESTMiddlewareSpec.Matcher:type_name -> kernel.MatcherSpec
	4, // 3: app.v1.SOAPRESTMiddlewareSpec.ErrorHandler:type_name -> kernel.Reference
	4, // [4:4] is the sub-list for method output_type
	4, // [4:4] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_app_v1_middleware_soaprest_proto_init() }
func file_app_v1_middleware_soaprest_proto_init() {
	if File_app_v1_middleware_soaprest_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_app_v1_middleware_soaprest_proto_rawDesc), len(file_app_v1_middleware_soaprest_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_app_v1_middleware_soaprest_proto_goTypes,
		DependencyIndexes: file_app_v1_middleware_soaprest_proto_depIdxs,
		MessageInfos:      file_app_v1_middleware_soaprest_proto_msgTypes,
	}.Build()
	File_app_v1_middleware_soaprest_proto = out.File
	file_app_v1_middleware_soaprest_proto_goTypes = nil
	file_app_v1_middleware_soaprest_proto_depIdxs = nil
}
