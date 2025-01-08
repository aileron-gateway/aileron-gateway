// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.1
// 	protoc        v5.27.2
// source: core/v1/httphandler.proto

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

// HTTPHandler is the definition of the HTTPHandler object.
// HTTPHandler implements interface of the http handler.
type HTTPHandler struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// [REQUIRED]
	// APIVersion is the defined version of the handler.
	// This value must be "core/v1".
	APIVersion string `protobuf:"bytes,1,opt,name=APIVersion,json=apiVersion,proto3" json:"APIVersion,omitempty"`
	// [REQUIRED]
	// Kind is the kind of this object.
	// This value must be "Handler".
	Kind string `protobuf:"bytes,2,opt,name=Kind,json=kind,proto3" json:"Kind,omitempty"`
	// [OPTIONAL]
	// Metadata is the metadata of the HTTP handler object.
	// If not set, both name and namespace in the metadata
	// are treated as "default".
	Metadata *kernel.Metadata `protobuf:"bytes,3,opt,name=Metadata,json=metadata,proto3" json:"Metadata,omitempty"`
	// [OPTIONAL]
	// Spec is the specification of the handler.
	// Default values are used when nothing is set.
	Spec          *HTTPHandlerSpec `protobuf:"bytes,4,opt,name=Spec,json=spec,proto3" json:"Spec,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *HTTPHandler) Reset() {
	*x = HTTPHandler{}
	mi := &file_core_v1_httphandler_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *HTTPHandler) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HTTPHandler) ProtoMessage() {}

func (x *HTTPHandler) ProtoReflect() protoreflect.Message {
	mi := &file_core_v1_httphandler_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HTTPHandler.ProtoReflect.Descriptor instead.
func (*HTTPHandler) Descriptor() ([]byte, []int) {
	return file_core_v1_httphandler_proto_rawDescGZIP(), []int{0}
}

func (x *HTTPHandler) GetAPIVersion() string {
	if x != nil {
		return x.APIVersion
	}
	return ""
}

func (x *HTTPHandler) GetKind() string {
	if x != nil {
		return x.Kind
	}
	return ""
}

func (x *HTTPHandler) GetMetadata() *kernel.Metadata {
	if x != nil {
		return x.Metadata
	}
	return nil
}

func (x *HTTPHandler) GetSpec() *HTTPHandlerSpec {
	if x != nil {
		return x.Spec
	}
	return nil
}

// HTTPHandlerSpec is the specifications for the Handler object.
type HTTPHandlerSpec struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// [OPTIONAL]
	// Pattern is path pattern that this handler
	// is registered to servers.
	// The pattern will be joined with the internal handler
	// that is specified with Handler.
	// Default is not set.
	Pattern string `protobuf:"bytes,1,opt,name=Pattern,json=pattern,proto3" json:"Pattern,omitempty"`
	// [OPTIONAL]
	// Middleware is the list of middleware applied for all handlers.
	// Default is not set.
	Middleware []*kernel.Reference `protobuf:"bytes,2,rep,name=Middleware,json=middleware,proto3" json:"Middleware,omitempty"`
	// [REQUIRED]
	// Handler is the reference to a handler to use.
	// Default is not set.
	Handler       *kernel.Reference `protobuf:"bytes,3,opt,name=Handler,json=handler,proto3" json:"Handler,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *HTTPHandlerSpec) Reset() {
	*x = HTTPHandlerSpec{}
	mi := &file_core_v1_httphandler_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *HTTPHandlerSpec) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HTTPHandlerSpec) ProtoMessage() {}

func (x *HTTPHandlerSpec) ProtoReflect() protoreflect.Message {
	mi := &file_core_v1_httphandler_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HTTPHandlerSpec.ProtoReflect.Descriptor instead.
func (*HTTPHandlerSpec) Descriptor() ([]byte, []int) {
	return file_core_v1_httphandler_proto_rawDescGZIP(), []int{1}
}

func (x *HTTPHandlerSpec) GetPattern() string {
	if x != nil {
		return x.Pattern
	}
	return ""
}

func (x *HTTPHandlerSpec) GetMiddleware() []*kernel.Reference {
	if x != nil {
		return x.Middleware
	}
	return nil
}

func (x *HTTPHandlerSpec) GetHandler() *kernel.Reference {
	if x != nil {
		return x.Handler
	}
	return nil
}

var File_core_v1_httphandler_proto protoreflect.FileDescriptor

var file_core_v1_httphandler_proto_rawDesc = []byte{
	0x0a, 0x19, 0x63, 0x6f, 0x72, 0x65, 0x2f, 0x76, 0x31, 0x2f, 0x68, 0x74, 0x74, 0x70, 0x68, 0x61,
	0x6e, 0x64, 0x6c, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x07, 0x63, 0x6f, 0x72,
	0x65, 0x2e, 0x76, 0x31, 0x1a, 0x1b, 0x62, 0x75, 0x66, 0x2f, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61,
	0x74, 0x65, 0x2f, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x1a, 0x15, 0x6b, 0x65, 0x72, 0x6e, 0x65, 0x6c, 0x2f, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72,
	0x63, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xc1, 0x01, 0x0a, 0x0b, 0x48, 0x54, 0x54,
	0x50, 0x48, 0x61, 0x6e, 0x64, 0x6c, 0x65, 0x72, 0x12, 0x2e, 0x0a, 0x0a, 0x41, 0x50, 0x49, 0x56,
	0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x0e, 0xba, 0x48,
	0x0b, 0x72, 0x09, 0x0a, 0x07, 0x63, 0x6f, 0x72, 0x65, 0x2f, 0x76, 0x31, 0x52, 0x0a, 0x61, 0x70,
	0x69, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x26, 0x0a, 0x04, 0x4b, 0x69, 0x6e, 0x64,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x42, 0x12, 0xba, 0x48, 0x0f, 0x72, 0x0d, 0x0a, 0x0b, 0x48,
	0x54, 0x54, 0x50, 0x48, 0x61, 0x6e, 0x64, 0x6c, 0x65, 0x72, 0x52, 0x04, 0x6b, 0x69, 0x6e, 0x64,
	0x12, 0x2c, 0x0a, 0x08, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x10, 0x2e, 0x6b, 0x65, 0x72, 0x6e, 0x65, 0x6c, 0x2e, 0x4d, 0x65, 0x74, 0x61,
	0x64, 0x61, 0x74, 0x61, 0x52, 0x08, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x12, 0x2c,
	0x0a, 0x04, 0x53, 0x70, 0x65, 0x63, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x18, 0x2e, 0x63,
	0x6f, 0x72, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x48, 0x54, 0x54, 0x50, 0x48, 0x61, 0x6e, 0x64, 0x6c,
	0x65, 0x72, 0x53, 0x70, 0x65, 0x63, 0x52, 0x04, 0x73, 0x70, 0x65, 0x63, 0x22, 0x8b, 0x01, 0x0a,
	0x0f, 0x48, 0x54, 0x54, 0x50, 0x48, 0x61, 0x6e, 0x64, 0x6c, 0x65, 0x72, 0x53, 0x70, 0x65, 0x63,
	0x12, 0x18, 0x0a, 0x07, 0x50, 0x61, 0x74, 0x74, 0x65, 0x72, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x07, 0x70, 0x61, 0x74, 0x74, 0x65, 0x72, 0x6e, 0x12, 0x31, 0x0a, 0x0a, 0x4d, 0x69,
	0x64, 0x64, 0x6c, 0x65, 0x77, 0x61, 0x72, 0x65, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x11,
	0x2e, 0x6b, 0x65, 0x72, 0x6e, 0x65, 0x6c, 0x2e, 0x52, 0x65, 0x66, 0x65, 0x72, 0x65, 0x6e, 0x63,
	0x65, 0x52, 0x0a, 0x6d, 0x69, 0x64, 0x64, 0x6c, 0x65, 0x77, 0x61, 0x72, 0x65, 0x12, 0x2b, 0x0a,
	0x07, 0x48, 0x61, 0x6e, 0x64, 0x6c, 0x65, 0x72, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x11,
	0x2e, 0x6b, 0x65, 0x72, 0x6e, 0x65, 0x6c, 0x2e, 0x52, 0x65, 0x66, 0x65, 0x72, 0x65, 0x6e, 0x63,
	0x65, 0x52, 0x07, 0x68, 0x61, 0x6e, 0x64, 0x6c, 0x65, 0x72, 0x42, 0x39, 0x5a, 0x37, 0x67, 0x69,
	0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x61, 0x69, 0x6c, 0x65, 0x72, 0x6f, 0x6e,
	0x2d, 0x67, 0x61, 0x74, 0x65, 0x77, 0x61, 0x79, 0x2f, 0x61, 0x69, 0x6c, 0x65, 0x72, 0x6f, 0x6e,
	0x2d, 0x67, 0x61, 0x74, 0x65, 0x77, 0x61, 0x79, 0x2f, 0x61, 0x70, 0x69, 0x73, 0x2f, 0x63, 0x6f,
	0x72, 0x65, 0x2f, 0x76, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_core_v1_httphandler_proto_rawDescOnce sync.Once
	file_core_v1_httphandler_proto_rawDescData = file_core_v1_httphandler_proto_rawDesc
)

func file_core_v1_httphandler_proto_rawDescGZIP() []byte {
	file_core_v1_httphandler_proto_rawDescOnce.Do(func() {
		file_core_v1_httphandler_proto_rawDescData = protoimpl.X.CompressGZIP(file_core_v1_httphandler_proto_rawDescData)
	})
	return file_core_v1_httphandler_proto_rawDescData
}

var file_core_v1_httphandler_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_core_v1_httphandler_proto_goTypes = []any{
	(*HTTPHandler)(nil),      // 0: core.v1.HTTPHandler
	(*HTTPHandlerSpec)(nil),  // 1: core.v1.HTTPHandlerSpec
	(*kernel.Metadata)(nil),  // 2: kernel.Metadata
	(*kernel.Reference)(nil), // 3: kernel.Reference
}
var file_core_v1_httphandler_proto_depIdxs = []int32{
	2, // 0: core.v1.HTTPHandler.Metadata:type_name -> kernel.Metadata
	1, // 1: core.v1.HTTPHandler.Spec:type_name -> core.v1.HTTPHandlerSpec
	3, // 2: core.v1.HTTPHandlerSpec.Middleware:type_name -> kernel.Reference
	3, // 3: core.v1.HTTPHandlerSpec.Handler:type_name -> kernel.Reference
	4, // [4:4] is the sub-list for method output_type
	4, // [4:4] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_core_v1_httphandler_proto_init() }
func file_core_v1_httphandler_proto_init() {
	if File_core_v1_httphandler_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_core_v1_httphandler_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_core_v1_httphandler_proto_goTypes,
		DependencyIndexes: file_core_v1_httphandler_proto_depIdxs,
		MessageInfos:      file_core_v1_httphandler_proto_msgTypes,
	}.Build()
	File_core_v1_httphandler_proto = out.File
	file_core_v1_httphandler_proto_rawDesc = nil
	file_core_v1_httphandler_proto_goTypes = nil
	file_core_v1_httphandler_proto_depIdxs = nil
}
