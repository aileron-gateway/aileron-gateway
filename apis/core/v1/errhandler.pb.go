// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.3
// 	protoc        v5.29.0
// source: core/v1/errhandler.proto

package v1

import (
	_ "buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go/buf/validate"
	kernel "github.com/aileron-gateway/aileron-gateway/apis/kernel"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// ErrorHandler is the definition of the ErrorHandler object.
// ErrorHandler implements interface of error handler.
type ErrorHandler struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// [REQUIRED]
	// APIVersion is the defined version of the error handler.
	// This value must be "core/v1".
	APIVersion string `protobuf:"bytes,1,opt,name=APIVersion,json=apiVersion,proto3" json:"APIVersion,omitempty"`
	// [REQUIRED]
	// Kind is the kind of this object.
	// This value must be "ErrorHandler".
	Kind string `protobuf:"bytes,2,opt,name=Kind,json=kind,proto3" json:"Kind,omitempty"`
	// [OPTIONAL]
	// Metadata is the metadata of the error handler object.
	// If not set, both name and namespace in the metadata
	// are treated as "default".
	Metadata *kernel.Metadata `protobuf:"bytes,3,opt,name=Metadata,json=metadata,proto3" json:"Metadata,omitempty"`
	// [OPTIONAL]
	// Spec is the specification of the error handler .
	// Default values are used when nothing is set.
	Spec          *ErrorHandlerSpec `protobuf:"bytes,4,opt,name=Spec,json=spec,proto3" json:"Spec,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ErrorHandler) Reset() {
	*x = ErrorHandler{}
	mi := &file_core_v1_errhandler_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ErrorHandler) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ErrorHandler) ProtoMessage() {}

func (x *ErrorHandler) ProtoReflect() protoreflect.Message {
	mi := &file_core_v1_errhandler_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ErrorHandler.ProtoReflect.Descriptor instead.
func (*ErrorHandler) Descriptor() ([]byte, []int) {
	return file_core_v1_errhandler_proto_rawDescGZIP(), []int{0}
}

func (x *ErrorHandler) GetAPIVersion() string {
	if x != nil {
		return x.APIVersion
	}
	return ""
}

func (x *ErrorHandler) GetKind() string {
	if x != nil {
		return x.Kind
	}
	return ""
}

func (x *ErrorHandler) GetMetadata() *kernel.Metadata {
	if x != nil {
		return x.Metadata
	}
	return nil
}

func (x *ErrorHandler) GetSpec() *ErrorHandlerSpec {
	if x != nil {
		return x.Spec
	}
	return nil
}

// ErrorHandlerSpec is the specifications for the ErrorHandler object.
type ErrorHandlerSpec struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// [OPTIONAL]
	// StackAlways is the flag to output stacktrace to the logger.
	// If set to true, this error handler output stacktrace
	// even when handling client side error, or 400-499 status error.
	// Default is [false].
	StackAlways bool `protobuf:"varint,1,opt,name=StackAlways,json=stackAlways,proto3" json:"StackAlways,omitempty"`
	// [OPTIONAL]
	// ErrorMessages is the list of custom error messages to overwrite.
	// Default is not set.
	ErrorMessages []*ErrorMessageSpec `protobuf:"bytes,2,rep,name=ErrorMessages,json=errorMessages,proto3" json:"ErrorMessages,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ErrorHandlerSpec) Reset() {
	*x = ErrorHandlerSpec{}
	mi := &file_core_v1_errhandler_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ErrorHandlerSpec) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ErrorHandlerSpec) ProtoMessage() {}

func (x *ErrorHandlerSpec) ProtoReflect() protoreflect.Message {
	mi := &file_core_v1_errhandler_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ErrorHandlerSpec.ProtoReflect.Descriptor instead.
func (*ErrorHandlerSpec) Descriptor() ([]byte, []int) {
	return file_core_v1_errhandler_proto_rawDescGZIP(), []int{1}
}

func (x *ErrorHandlerSpec) GetStackAlways() bool {
	if x != nil {
		return x.StackAlways
	}
	return false
}

func (x *ErrorHandlerSpec) GetErrorMessages() []*ErrorMessageSpec {
	if x != nil {
		return x.ErrorMessages
	}
	return nil
}

// ErrorMessageSpec is the specification of HTTP error response.
type ErrorMessageSpec struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// [OPTIONAL]
	// Codes are the list of error code patterns to match this message.
	// String is evaluated by the path match.
	// For example, "E123*" or "E12*".
	// See https://pkg.go.dev/path#Match for for available expressions.
	// If an error matched to one of the Codes, Kinds or Message patterns
	// then the response is overwritten by the MIMEContents.
	// Default is not set.
	Codes []string `protobuf:"bytes,2,rep,name=Codes,json=codes,proto3" json:"Codes,omitempty"`
	// [OPTIONAL]
	// Kinds are the list of error kind patterns to match this message.
	// String is evaluated by the path match.
	// See https://pkg.go.dev/path#Match for for available expressions.
	// If an error matched to one of the Codes, Kinds or Message patterns
	// then the response is overwritten by the MIMEContents.
	// Default is not set.
	Kinds []string `protobuf:"bytes,3,rep,name=Kinds,json=kinds,proto3" json:"Kinds,omitempty"`
	// [OPTIONAL]
	// Messages are the list of error message pattarns to match this message.
	// String is evaluated by the regular expression.
	// See https://pkg.go.dev/regexp and https://github.com/google/re2/wiki/Syntax
	// for available expressions.
	// If an error matched to one of the Codes, Kinds or Message patterns
	// then the response is overwritten by the MIMEContents.
	// Default is not set.
	Messages []string `protobuf:"bytes,4,rep,name=Messages,json=messages,proto3" json:"Messages,omitempty"`
	// [OPTIONAL]
	// HeaderTemplate is the key-value pairs of HTTP headers
	// to add to the error response.
	// Unlike headers that can set in MIMEContents field,
	// values can be written in template.
	// HTTP status code `{{status}}` and status text `{{statusText}}`,
	// error code and kind `{{code}}`, `{{kind}}` can be used in the value.
	// Header names cannot be in template format.
	// This field is mainly intended to set error redirecting headers.
	// Default is not set.
	HeaderTemplate map[string]string `protobuf:"bytes,5,rep,name=HeaderTemplate,json=headerTemplate,proto3" json:"HeaderTemplate,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
	// [OPTIONAL]
	// MIMEContents is the list of mime content to be used for overwriting the error response.
	// If an error matched to one of the Codes, Kinds or Message patterns
	// then the response is overwritten by one of the MIMEContents.
	// Responses are not overwritten if this field has no content.
	// The first one is used when the Accept header did not matched to any content.
	// Default is not set.
	MIMEContents  []*MIMEContentSpec `protobuf:"bytes,6,rep,name=MIMEContents,json=mimeContents,proto3" json:"MIMEContents,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ErrorMessageSpec) Reset() {
	*x = ErrorMessageSpec{}
	mi := &file_core_v1_errhandler_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ErrorMessageSpec) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ErrorMessageSpec) ProtoMessage() {}

func (x *ErrorMessageSpec) ProtoReflect() protoreflect.Message {
	mi := &file_core_v1_errhandler_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ErrorMessageSpec.ProtoReflect.Descriptor instead.
func (*ErrorMessageSpec) Descriptor() ([]byte, []int) {
	return file_core_v1_errhandler_proto_rawDescGZIP(), []int{2}
}

func (x *ErrorMessageSpec) GetCodes() []string {
	if x != nil {
		return x.Codes
	}
	return nil
}

func (x *ErrorMessageSpec) GetKinds() []string {
	if x != nil {
		return x.Kinds
	}
	return nil
}

func (x *ErrorMessageSpec) GetMessages() []string {
	if x != nil {
		return x.Messages
	}
	return nil
}

func (x *ErrorMessageSpec) GetHeaderTemplate() map[string]string {
	if x != nil {
		return x.HeaderTemplate
	}
	return nil
}

func (x *ErrorMessageSpec) GetMIMEContents() []*MIMEContentSpec {
	if x != nil {
		return x.MIMEContents
	}
	return nil
}

var File_core_v1_errhandler_proto protoreflect.FileDescriptor

var file_core_v1_errhandler_proto_rawDesc = []byte{
	0x0a, 0x18, 0x63, 0x6f, 0x72, 0x65, 0x2f, 0x76, 0x31, 0x2f, 0x65, 0x72, 0x72, 0x68, 0x61, 0x6e,
	0x64, 0x6c, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x07, 0x63, 0x6f, 0x72, 0x65,
	0x2e, 0x76, 0x31, 0x1a, 0x1b, 0x62, 0x75, 0x66, 0x2f, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74,
	0x65, 0x2f, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x1a, 0x16, 0x63, 0x6f, 0x72, 0x65, 0x2f, 0x76, 0x31, 0x2f, 0x74, 0x65, 0x6d, 0x70, 0x6c, 0x61,
	0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x15, 0x6b, 0x65, 0x72, 0x6e, 0x65, 0x6c,
	0x2f, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22,
	0xc4, 0x01, 0x0a, 0x0c, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x48, 0x61, 0x6e, 0x64, 0x6c, 0x65, 0x72,
	0x12, 0x2e, 0x0a, 0x0a, 0x41, 0x50, 0x49, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x42, 0x0e, 0xba, 0x48, 0x0b, 0x72, 0x09, 0x0a, 0x07, 0x63, 0x6f, 0x72,
	0x65, 0x2f, 0x76, 0x31, 0x52, 0x0a, 0x61, 0x70, 0x69, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e,
	0x12, 0x27, 0x0a, 0x04, 0x4b, 0x69, 0x6e, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x42, 0x13,
	0xba, 0x48, 0x10, 0x72, 0x0e, 0x0a, 0x0c, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x48, 0x61, 0x6e, 0x64,
	0x6c, 0x65, 0x72, 0x52, 0x04, 0x6b, 0x69, 0x6e, 0x64, 0x12, 0x2c, 0x0a, 0x08, 0x4d, 0x65, 0x74,
	0x61, 0x64, 0x61, 0x74, 0x61, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x6b, 0x65,
	0x72, 0x6e, 0x65, 0x6c, 0x2e, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x52, 0x08, 0x6d,
	0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x12, 0x2d, 0x0a, 0x04, 0x53, 0x70, 0x65, 0x63, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x76, 0x31, 0x2e,
	0x45, 0x72, 0x72, 0x6f, 0x72, 0x48, 0x61, 0x6e, 0x64, 0x6c, 0x65, 0x72, 0x53, 0x70, 0x65, 0x63,
	0x52, 0x04, 0x73, 0x70, 0x65, 0x63, 0x22, 0x75, 0x0a, 0x10, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x48,
	0x61, 0x6e, 0x64, 0x6c, 0x65, 0x72, 0x53, 0x70, 0x65, 0x63, 0x12, 0x20, 0x0a, 0x0b, 0x53, 0x74,
	0x61, 0x63, 0x6b, 0x41, 0x6c, 0x77, 0x61, 0x79, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52,
	0x0b, 0x73, 0x74, 0x61, 0x63, 0x6b, 0x41, 0x6c, 0x77, 0x61, 0x79, 0x73, 0x12, 0x3f, 0x0a, 0x0d,
	0x45, 0x72, 0x72, 0x6f, 0x72, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x18, 0x02, 0x20,
	0x03, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x45, 0x72,
	0x72, 0x6f, 0x72, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x53, 0x70, 0x65, 0x63, 0x52, 0x0d,
	0x65, 0x72, 0x72, 0x6f, 0x72, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x22, 0xe2, 0x02,
	0x0a, 0x10, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x53, 0x70,
	0x65, 0x63, 0x12, 0x24, 0x0a, 0x05, 0x43, 0x6f, 0x64, 0x65, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28,
	0x09, 0x42, 0x0e, 0xba, 0x48, 0x0b, 0x92, 0x01, 0x08, 0x18, 0x01, 0x22, 0x04, 0x72, 0x02, 0x10,
	0x01, 0x52, 0x05, 0x63, 0x6f, 0x64, 0x65, 0x73, 0x12, 0x24, 0x0a, 0x05, 0x4b, 0x69, 0x6e, 0x64,
	0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x09, 0x42, 0x0e, 0xba, 0x48, 0x0b, 0x92, 0x01, 0x08, 0x18,
	0x01, 0x22, 0x04, 0x72, 0x02, 0x10, 0x01, 0x52, 0x05, 0x6b, 0x69, 0x6e, 0x64, 0x73, 0x12, 0x2a,
	0x0a, 0x08, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x18, 0x04, 0x20, 0x03, 0x28, 0x09,
	0x42, 0x0e, 0xba, 0x48, 0x0b, 0x92, 0x01, 0x08, 0x18, 0x01, 0x22, 0x04, 0x72, 0x02, 0x10, 0x01,
	0x52, 0x08, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x12, 0x55, 0x0a, 0x0e, 0x48, 0x65,
	0x61, 0x64, 0x65, 0x72, 0x54, 0x65, 0x6d, 0x70, 0x6c, 0x61, 0x74, 0x65, 0x18, 0x05, 0x20, 0x03,
	0x28, 0x0b, 0x32, 0x2d, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x45, 0x72, 0x72,
	0x6f, 0x72, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x53, 0x70, 0x65, 0x63, 0x2e, 0x48, 0x65,
	0x61, 0x64, 0x65, 0x72, 0x54, 0x65, 0x6d, 0x70, 0x6c, 0x61, 0x74, 0x65, 0x45, 0x6e, 0x74, 0x72,
	0x79, 0x52, 0x0e, 0x68, 0x65, 0x61, 0x64, 0x65, 0x72, 0x54, 0x65, 0x6d, 0x70, 0x6c, 0x61, 0x74,
	0x65, 0x12, 0x3c, 0x0a, 0x0c, 0x4d, 0x49, 0x4d, 0x45, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74,
	0x73, 0x18, 0x06, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x18, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x76,
	0x31, 0x2e, 0x4d, 0x49, 0x4d, 0x45, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x53, 0x70, 0x65,
	0x63, 0x52, 0x0c, 0x6d, 0x69, 0x6d, 0x65, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x73, 0x1a,
	0x41, 0x0a, 0x13, 0x48, 0x65, 0x61, 0x64, 0x65, 0x72, 0x54, 0x65, 0x6d, 0x70, 0x6c, 0x61, 0x74,
	0x65, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75,
	0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02,
	0x38, 0x01, 0x42, 0x39, 0x5a, 0x37, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d,
	0x2f, 0x61, 0x69, 0x6c, 0x65, 0x72, 0x6f, 0x6e, 0x2d, 0x67, 0x61, 0x74, 0x65, 0x77, 0x61, 0x79,
	0x2f, 0x61, 0x69, 0x6c, 0x65, 0x72, 0x6f, 0x6e, 0x2d, 0x67, 0x61, 0x74, 0x65, 0x77, 0x61, 0x79,
	0x2f, 0x61, 0x70, 0x69, 0x73, 0x2f, 0x63, 0x6f, 0x72, 0x65, 0x2f, 0x76, 0x31, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_core_v1_errhandler_proto_rawDescOnce sync.Once
	file_core_v1_errhandler_proto_rawDescData = file_core_v1_errhandler_proto_rawDesc
)

func file_core_v1_errhandler_proto_rawDescGZIP() []byte {
	file_core_v1_errhandler_proto_rawDescOnce.Do(func() {
		file_core_v1_errhandler_proto_rawDescData = protoimpl.X.CompressGZIP(file_core_v1_errhandler_proto_rawDescData)
	})
	return file_core_v1_errhandler_proto_rawDescData
}

var file_core_v1_errhandler_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_core_v1_errhandler_proto_goTypes = []any{
	(*ErrorHandler)(nil),     // 0: core.v1.ErrorHandler
	(*ErrorHandlerSpec)(nil), // 1: core.v1.ErrorHandlerSpec
	(*ErrorMessageSpec)(nil), // 2: core.v1.ErrorMessageSpec
	nil,                      // 3: core.v1.ErrorMessageSpec.HeaderTemplateEntry
	(*kernel.Metadata)(nil),  // 4: kernel.Metadata
	(*MIMEContentSpec)(nil),  // 5: core.v1.MIMEContentSpec
}
var file_core_v1_errhandler_proto_depIdxs = []int32{
	4, // 0: core.v1.ErrorHandler.Metadata:type_name -> kernel.Metadata
	1, // 1: core.v1.ErrorHandler.Spec:type_name -> core.v1.ErrorHandlerSpec
	2, // 2: core.v1.ErrorHandlerSpec.ErrorMessages:type_name -> core.v1.ErrorMessageSpec
	3, // 3: core.v1.ErrorMessageSpec.HeaderTemplate:type_name -> core.v1.ErrorMessageSpec.HeaderTemplateEntry
	5, // 4: core.v1.ErrorMessageSpec.MIMEContents:type_name -> core.v1.MIMEContentSpec
	5, // [5:5] is the sub-list for method output_type
	5, // [5:5] is the sub-list for method input_type
	5, // [5:5] is the sub-list for extension type_name
	5, // [5:5] is the sub-list for extension extendee
	0, // [0:5] is the sub-list for field type_name
}

func init() { file_core_v1_errhandler_proto_init() }
func file_core_v1_errhandler_proto_init() {
	if File_core_v1_errhandler_proto != nil {
		return
	}
	file_core_v1_template_proto_init()
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_core_v1_errhandler_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_core_v1_errhandler_proto_goTypes,
		DependencyIndexes: file_core_v1_errhandler_proto_depIdxs,
		MessageInfos:      file_core_v1_errhandler_proto_msgTypes,
	}.Build()
	File_core_v1_errhandler_proto = out.File
	file_core_v1_errhandler_proto_rawDesc = nil
	file_core_v1_errhandler_proto_goTypes = nil
	file_core_v1_errhandler_proto_depIdxs = nil
}
