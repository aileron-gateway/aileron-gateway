// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.5
// 	protoc        v5.29.0
// source: core/v1/template.proto

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

// TemplateHandler is the definition of the TemplateHandler object.
// TemplateHandler implements interface of the http handler.
type TemplateHandler struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// [REQUIRED]
	// APIVersion is the defined version of the handler.
	// This value must be "core/v1".
	APIVersion string `protobuf:"bytes,1,opt,name=APIVersion,json=apiVersion,proto3" json:"APIVersion,omitempty"`
	// [REQUIRED]
	// Kind is the kind of this object.
	// This value must be "TemplateHandler".
	Kind string `protobuf:"bytes,2,opt,name=Kind,json=kind,proto3" json:"Kind,omitempty"`
	// [OPTIONAL]
	// Metadata is the metadata of the handler object.
	// If not set, both name and namespace in the metadata
	// are treated as "default".
	Metadata *kernel.Metadata `protobuf:"bytes,3,opt,name=Metadata,json=metadata,proto3" json:"Metadata,omitempty"`
	// [OPTIONAL]
	// Spec is the specification of the handler.
	// Default values are used when nothing is set.
	Spec          *TemplateHandlerSpec `protobuf:"bytes,4,opt,name=Spec,json=spec,proto3" json:"Spec,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *TemplateHandler) Reset() {
	*x = TemplateHandler{}
	mi := &file_core_v1_template_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *TemplateHandler) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TemplateHandler) ProtoMessage() {}

func (x *TemplateHandler) ProtoReflect() protoreflect.Message {
	mi := &file_core_v1_template_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TemplateHandler.ProtoReflect.Descriptor instead.
func (*TemplateHandler) Descriptor() ([]byte, []int) {
	return file_core_v1_template_proto_rawDescGZIP(), []int{0}
}

func (x *TemplateHandler) GetAPIVersion() string {
	if x != nil {
		return x.APIVersion
	}
	return ""
}

func (x *TemplateHandler) GetKind() string {
	if x != nil {
		return x.Kind
	}
	return ""
}

func (x *TemplateHandler) GetMetadata() *kernel.Metadata {
	if x != nil {
		return x.Metadata
	}
	return nil
}

func (x *TemplateHandler) GetSpec() *TemplateHandlerSpec {
	if x != nil {
		return x.Spec
	}
	return nil
}

// TemplateHandlerSpec is the specifications for the TemplateHandler object.
type TemplateHandlerSpec struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// [OPTIONAL]
	// ErrorHandler is the reference to a ErrorHandler object.
	// Referred object must implement ErrorHandler interface.
	// Default error handler is used when not set.
	ErrorHandler *kernel.Reference `protobuf:"bytes,1,opt,name=ErrorHandler,json=errorHandler,proto3" json:"ErrorHandler,omitempty"`
	// [OPTIONAL]
	// Patterns is path patterns that this handler
	// is registered to a server.
	// Default is not set.
	Patterns []string `protobuf:"bytes,2,rep,name=Patterns,json=patterns,proto3" json:"Patterns,omitempty"`
	// [OPTIONA]
	// Methods is the list of HTTP method this handler can handle.
	// Note that it depends on the multiplexer, or HTTP router
	// if this field can be used.
	// If not set, all methods are accepted.
	// Default is not set.
	Methods []HTTPMethod `protobuf:"varint,3,rep,packed,name=Methods,json=methods,proto3,enum=core.v1.HTTPMethod" json:"Methods,omitempty"`
	// [REQUIRED]
	// MIMEContents is the list of content that will be returned by this handler.
	// When no appropriate content were found,
	// not acceptable error will be returned to clients.
	// Default is not set.
	MIMEContents  []*MIMEContentSpec `protobuf:"bytes,4,rep,name=MIMEContents,json=mimeContents,proto3" json:"MIMEContents,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *TemplateHandlerSpec) Reset() {
	*x = TemplateHandlerSpec{}
	mi := &file_core_v1_template_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *TemplateHandlerSpec) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TemplateHandlerSpec) ProtoMessage() {}

func (x *TemplateHandlerSpec) ProtoReflect() protoreflect.Message {
	mi := &file_core_v1_template_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TemplateHandlerSpec.ProtoReflect.Descriptor instead.
func (*TemplateHandlerSpec) Descriptor() ([]byte, []int) {
	return file_core_v1_template_proto_rawDescGZIP(), []int{1}
}

func (x *TemplateHandlerSpec) GetErrorHandler() *kernel.Reference {
	if x != nil {
		return x.ErrorHandler
	}
	return nil
}

func (x *TemplateHandlerSpec) GetPatterns() []string {
	if x != nil {
		return x.Patterns
	}
	return nil
}

func (x *TemplateHandlerSpec) GetMethods() []HTTPMethod {
	if x != nil {
		return x.Methods
	}
	return nil
}

func (x *TemplateHandlerSpec) GetMIMEContents() []*MIMEContentSpec {
	if x != nil {
		return x.MIMEContents
	}
	return nil
}

// MIMEContentSpec is the specification for the MIMEContent.
type MIMEContentSpec struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// [REQUIRED]
	// MIMEType is the mediatype of this content.
	// See the following documents for available type.
	// https://www.iana.org/assignments/media-types/media-types.xhtml
	// Default is not set.
	MIMEType string `protobuf:"bytes,1,opt,name=MIMEType,json=mimeType,proto3" json:"MIMEType,omitempty"`
	// [OPTIONAL]
	// StatusCode is the http status code used to respond this content.
	// Default is [200].
	StatusCode int32 `protobuf:"varint,2,opt,name=StatusCode,json=statusCode,proto3" json:"StatusCode,omitempty"`
	// [OPTIONAL]
	// Header is the key-value pairs of HTTP headers
	// to add to the response.
	// Keys must be a valid http header name.
	// Default is not set.
	Header map[string]string `protobuf:"bytes,3,rep,name=Header,json=header,proto3" json:"Header,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
	// [OPTIONAL]
	// TemplateType is the template engine type of this content.
	// Default is [Text].
	TemplateType kernel.TemplateType `protobuf:"varint,4,opt,name=TemplateType,json=templateType,proto3,enum=kernel.TemplateType" json:"TemplateType,omitempty"`
	// [OPTIONAL]
	// Template is the template text to generate response body.
	// TemplateFile is prior to Template if both parameters are set.
	// Default is not set.
	Template string `protobuf:"bytes,5,opt,name=Template,json=template,proto3" json:"Template,omitempty"`
	// [OPTIONAL]
	// TemplateFile is the fail path to read template from.
	// TemplateFile is prior to Template if both parameters are set.
	// It does not matter wheather the path is relative or absolute.
	// Default is not set.
	TemplateFile string `protobuf:"bytes,6,opt,name=TemplateFile,json=templateFile,proto3" json:"TemplateFile,omitempty"`
	// [OPTIONAL]
	// FallbackText is the response body that will be used for fallback.
	// This text will be used when generating response body
	// with GoText, GoHTML failed.
	// The value is not used when the TemplateType is Text.
	// Default is not set.
	FallbackText  string `protobuf:"bytes,7,opt,name=FallbackText,json=fallbackText,proto3" json:"FallbackText,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *MIMEContentSpec) Reset() {
	*x = MIMEContentSpec{}
	mi := &file_core_v1_template_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *MIMEContentSpec) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MIMEContentSpec) ProtoMessage() {}

func (x *MIMEContentSpec) ProtoReflect() protoreflect.Message {
	mi := &file_core_v1_template_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MIMEContentSpec.ProtoReflect.Descriptor instead.
func (*MIMEContentSpec) Descriptor() ([]byte, []int) {
	return file_core_v1_template_proto_rawDescGZIP(), []int{2}
}

func (x *MIMEContentSpec) GetMIMEType() string {
	if x != nil {
		return x.MIMEType
	}
	return ""
}

func (x *MIMEContentSpec) GetStatusCode() int32 {
	if x != nil {
		return x.StatusCode
	}
	return 0
}

func (x *MIMEContentSpec) GetHeader() map[string]string {
	if x != nil {
		return x.Header
	}
	return nil
}

func (x *MIMEContentSpec) GetTemplateType() kernel.TemplateType {
	if x != nil {
		return x.TemplateType
	}
	return kernel.TemplateType(0)
}

func (x *MIMEContentSpec) GetTemplate() string {
	if x != nil {
		return x.Template
	}
	return ""
}

func (x *MIMEContentSpec) GetTemplateFile() string {
	if x != nil {
		return x.TemplateFile
	}
	return ""
}

func (x *MIMEContentSpec) GetFallbackText() string {
	if x != nil {
		return x.FallbackText
	}
	return ""
}

var File_core_v1_template_proto protoreflect.FileDescriptor

var file_core_v1_template_proto_rawDesc = string([]byte{
	0x0a, 0x16, 0x63, 0x6f, 0x72, 0x65, 0x2f, 0x76, 0x31, 0x2f, 0x74, 0x65, 0x6d, 0x70, 0x6c, 0x61,
	0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x07, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x76,
	0x31, 0x1a, 0x1b, 0x62, 0x75, 0x66, 0x2f, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2f,
	0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x12,
	0x63, 0x6f, 0x72, 0x65, 0x2f, 0x76, 0x31, 0x2f, 0x68, 0x74, 0x74, 0x70, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x1a, 0x15, 0x6b, 0x65, 0x72, 0x6e, 0x65, 0x6c, 0x2f, 0x72, 0x65, 0x73, 0x6f, 0x75,
	0x72, 0x63, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x14, 0x6b, 0x65, 0x72, 0x6e, 0x65,
	0x6c, 0x2f, 0x74, 0x78, 0x74, 0x75, 0x74, 0x69, 0x6c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22,
	0xcd, 0x01, 0x0a, 0x0f, 0x54, 0x65, 0x6d, 0x70, 0x6c, 0x61, 0x74, 0x65, 0x48, 0x61, 0x6e, 0x64,
	0x6c, 0x65, 0x72, 0x12, 0x2e, 0x0a, 0x0a, 0x41, 0x50, 0x49, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f,
	0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x0e, 0xba, 0x48, 0x0b, 0x72, 0x09, 0x0a, 0x07,
	0x63, 0x6f, 0x72, 0x65, 0x2f, 0x76, 0x31, 0x52, 0x0a, 0x61, 0x70, 0x69, 0x56, 0x65, 0x72, 0x73,
	0x69, 0x6f, 0x6e, 0x12, 0x2a, 0x0a, 0x04, 0x4b, 0x69, 0x6e, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x42, 0x16, 0xba, 0x48, 0x13, 0x72, 0x11, 0x0a, 0x0f, 0x54, 0x65, 0x6d, 0x70, 0x6c, 0x61,
	0x74, 0x65, 0x48, 0x61, 0x6e, 0x64, 0x6c, 0x65, 0x72, 0x52, 0x04, 0x6b, 0x69, 0x6e, 0x64, 0x12,
	0x2c, 0x0a, 0x08, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x10, 0x2e, 0x6b, 0x65, 0x72, 0x6e, 0x65, 0x6c, 0x2e, 0x4d, 0x65, 0x74, 0x61, 0x64,
	0x61, 0x74, 0x61, 0x52, 0x08, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x12, 0x30, 0x0a,
	0x04, 0x53, 0x70, 0x65, 0x63, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x63, 0x6f,
	0x72, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x54, 0x65, 0x6d, 0x70, 0x6c, 0x61, 0x74, 0x65, 0x48, 0x61,
	0x6e, 0x64, 0x6c, 0x65, 0x72, 0x53, 0x70, 0x65, 0x63, 0x52, 0x04, 0x73, 0x70, 0x65, 0x63, 0x22,
	0xe9, 0x01, 0x0a, 0x13, 0x54, 0x65, 0x6d, 0x70, 0x6c, 0x61, 0x74, 0x65, 0x48, 0x61, 0x6e, 0x64,
	0x6c, 0x65, 0x72, 0x53, 0x70, 0x65, 0x63, 0x12, 0x35, 0x0a, 0x0c, 0x45, 0x72, 0x72, 0x6f, 0x72,
	0x48, 0x61, 0x6e, 0x64, 0x6c, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x11, 0x2e,
	0x6b, 0x65, 0x72, 0x6e, 0x65, 0x6c, 0x2e, 0x52, 0x65, 0x66, 0x65, 0x72, 0x65, 0x6e, 0x63, 0x65,
	0x52, 0x0c, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x48, 0x61, 0x6e, 0x64, 0x6c, 0x65, 0x72, 0x12, 0x24,
	0x0a, 0x08, 0x50, 0x61, 0x74, 0x74, 0x65, 0x72, 0x6e, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x09,
	0x42, 0x08, 0xba, 0x48, 0x05, 0x92, 0x01, 0x02, 0x18, 0x01, 0x52, 0x08, 0x70, 0x61, 0x74, 0x74,
	0x65, 0x72, 0x6e, 0x73, 0x12, 0x37, 0x0a, 0x07, 0x4d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x73, 0x18,
	0x03, 0x20, 0x03, 0x28, 0x0e, 0x32, 0x13, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x76, 0x31, 0x2e,
	0x48, 0x54, 0x54, 0x50, 0x4d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x42, 0x08, 0xba, 0x48, 0x05, 0x92,
	0x01, 0x02, 0x18, 0x01, 0x52, 0x07, 0x6d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x73, 0x12, 0x3c, 0x0a,
	0x0c, 0x4d, 0x49, 0x4d, 0x45, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x73, 0x18, 0x04, 0x20,
	0x03, 0x28, 0x0b, 0x32, 0x18, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x4d, 0x49,
	0x4d, 0x45, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x53, 0x70, 0x65, 0x63, 0x52, 0x0c, 0x6d,
	0x69, 0x6d, 0x65, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x73, 0x22, 0x91, 0x03, 0x0a, 0x0f,
	0x4d, 0x49, 0x4d, 0x45, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x53, 0x70, 0x65, 0x63, 0x12,
	0x3b, 0x0a, 0x08, 0x4d, 0x49, 0x4d, 0x45, 0x54, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x42, 0x1f, 0xba, 0x48, 0x1c, 0x72, 0x1a, 0x32, 0x18, 0x5e, 0x5b, 0x61, 0x2d, 0x7a, 0x5d,
	0x2b, 0x2f, 0x5b, 0x30, 0x2d, 0x39, 0x61, 0x2d, 0x7a, 0x41, 0x2d, 0x5a, 0x2e, 0x2b, 0x2d, 0x5d,
	0x2b, 0x24, 0x52, 0x08, 0x6d, 0x69, 0x6d, 0x65, 0x54, 0x79, 0x70, 0x65, 0x12, 0x2a, 0x0a, 0x0a,
	0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x43, 0x6f, 0x64, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05,
	0x42, 0x0a, 0xba, 0x48, 0x07, 0x1a, 0x05, 0x18, 0xe7, 0x07, 0x28, 0x00, 0x52, 0x0a, 0x73, 0x74,
	0x61, 0x74, 0x75, 0x73, 0x43, 0x6f, 0x64, 0x65, 0x12, 0x3c, 0x0a, 0x06, 0x48, 0x65, 0x61, 0x64,
	0x65, 0x72, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x24, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e,
	0x76, 0x31, 0x2e, 0x4d, 0x49, 0x4d, 0x45, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x53, 0x70,
	0x65, 0x63, 0x2e, 0x48, 0x65, 0x61, 0x64, 0x65, 0x72, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x06,
	0x68, 0x65, 0x61, 0x64, 0x65, 0x72, 0x12, 0x38, 0x0a, 0x0c, 0x54, 0x65, 0x6d, 0x70, 0x6c, 0x61,
	0x74, 0x65, 0x54, 0x79, 0x70, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x14, 0x2e, 0x6b,
	0x65, 0x72, 0x6e, 0x65, 0x6c, 0x2e, 0x54, 0x65, 0x6d, 0x70, 0x6c, 0x61, 0x74, 0x65, 0x54, 0x79,
	0x70, 0x65, 0x52, 0x0c, 0x74, 0x65, 0x6d, 0x70, 0x6c, 0x61, 0x74, 0x65, 0x54, 0x79, 0x70, 0x65,
	0x12, 0x1a, 0x0a, 0x08, 0x54, 0x65, 0x6d, 0x70, 0x6c, 0x61, 0x74, 0x65, 0x18, 0x05, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x08, 0x74, 0x65, 0x6d, 0x70, 0x6c, 0x61, 0x74, 0x65, 0x12, 0x22, 0x0a, 0x0c,
	0x54, 0x65, 0x6d, 0x70, 0x6c, 0x61, 0x74, 0x65, 0x46, 0x69, 0x6c, 0x65, 0x18, 0x06, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0c, 0x74, 0x65, 0x6d, 0x70, 0x6c, 0x61, 0x74, 0x65, 0x46, 0x69, 0x6c, 0x65,
	0x12, 0x22, 0x0a, 0x0c, 0x46, 0x61, 0x6c, 0x6c, 0x62, 0x61, 0x63, 0x6b, 0x54, 0x65, 0x78, 0x74,
	0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x66, 0x61, 0x6c, 0x6c, 0x62, 0x61, 0x63, 0x6b,
	0x54, 0x65, 0x78, 0x74, 0x1a, 0x39, 0x0a, 0x0b, 0x48, 0x65, 0x61, 0x64, 0x65, 0x72, 0x45, 0x6e,
	0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x42,
	0x39, 0x5a, 0x37, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x61, 0x69,
	0x6c, 0x65, 0x72, 0x6f, 0x6e, 0x2d, 0x67, 0x61, 0x74, 0x65, 0x77, 0x61, 0x79, 0x2f, 0x61, 0x69,
	0x6c, 0x65, 0x72, 0x6f, 0x6e, 0x2d, 0x67, 0x61, 0x74, 0x65, 0x77, 0x61, 0x79, 0x2f, 0x61, 0x70,
	0x69, 0x73, 0x2f, 0x63, 0x6f, 0x72, 0x65, 0x2f, 0x76, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
})

var (
	file_core_v1_template_proto_rawDescOnce sync.Once
	file_core_v1_template_proto_rawDescData []byte
)

func file_core_v1_template_proto_rawDescGZIP() []byte {
	file_core_v1_template_proto_rawDescOnce.Do(func() {
		file_core_v1_template_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_core_v1_template_proto_rawDesc), len(file_core_v1_template_proto_rawDesc)))
	})
	return file_core_v1_template_proto_rawDescData
}

var file_core_v1_template_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_core_v1_template_proto_goTypes = []any{
	(*TemplateHandler)(nil),     // 0: core.v1.TemplateHandler
	(*TemplateHandlerSpec)(nil), // 1: core.v1.TemplateHandlerSpec
	(*MIMEContentSpec)(nil),     // 2: core.v1.MIMEContentSpec
	nil,                         // 3: core.v1.MIMEContentSpec.HeaderEntry
	(*kernel.Metadata)(nil),     // 4: kernel.Metadata
	(*kernel.Reference)(nil),    // 5: kernel.Reference
	(HTTPMethod)(0),             // 6: core.v1.HTTPMethod
	(kernel.TemplateType)(0),    // 7: kernel.TemplateType
}
var file_core_v1_template_proto_depIdxs = []int32{
	4, // 0: core.v1.TemplateHandler.Metadata:type_name -> kernel.Metadata
	1, // 1: core.v1.TemplateHandler.Spec:type_name -> core.v1.TemplateHandlerSpec
	5, // 2: core.v1.TemplateHandlerSpec.ErrorHandler:type_name -> kernel.Reference
	6, // 3: core.v1.TemplateHandlerSpec.Methods:type_name -> core.v1.HTTPMethod
	2, // 4: core.v1.TemplateHandlerSpec.MIMEContents:type_name -> core.v1.MIMEContentSpec
	3, // 5: core.v1.MIMEContentSpec.Header:type_name -> core.v1.MIMEContentSpec.HeaderEntry
	7, // 6: core.v1.MIMEContentSpec.TemplateType:type_name -> kernel.TemplateType
	7, // [7:7] is the sub-list for method output_type
	7, // [7:7] is the sub-list for method input_type
	7, // [7:7] is the sub-list for extension type_name
	7, // [7:7] is the sub-list for extension extendee
	0, // [0:7] is the sub-list for field type_name
}

func init() { file_core_v1_template_proto_init() }
func file_core_v1_template_proto_init() {
	if File_core_v1_template_proto != nil {
		return
	}
	file_core_v1_http_proto_init()
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_core_v1_template_proto_rawDesc), len(file_core_v1_template_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_core_v1_template_proto_goTypes,
		DependencyIndexes: file_core_v1_template_proto_depIdxs,
		MessageInfos:      file_core_v1_template_proto_msgTypes,
	}.Build()
	File_core_v1_template_proto = out.File
	file_core_v1_template_proto_goTypes = nil
	file_core_v1_template_proto_depIdxs = nil
}
