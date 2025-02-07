// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.5
// 	protoc        v5.29.0
// source: core/v1/static.proto

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

// StaticFileHandler is the definition of the StaticFileHandler object.
// StaticFileHandler implements interface of the http handler.
type StaticFileHandler struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// [REQUIRED]
	// APIVersion is the defined version of the handler.
	// This value must be "core/v1".
	APIVersion string `protobuf:"bytes,1,opt,name=APIVersion,json=apiVersion,proto3" json:"APIVersion,omitempty"`
	// [REQUIRED]
	// Kind is the kind of this object.
	// This value must be "StaticFileHandler".
	Kind string `protobuf:"bytes,2,opt,name=Kind,json=kind,proto3" json:"Kind,omitempty"`
	// [OPTIONAL]
	// Metadata is the metadata of the handler object.
	// If not set, both name and namespace in the metadata
	// are treated as "default".
	Metadata *kernel.Metadata `protobuf:"bytes,3,opt,name=Metadata,json=metadata,proto3" json:"Metadata,omitempty"`
	// [OPTIONAL]
	// Spec is the specification of the handler.
	// Default values are used when nothing is set.
	Spec          *StaticFileHandlerSpec `protobuf:"bytes,4,opt,name=Spec,json=spec,proto3" json:"Spec,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *StaticFileHandler) Reset() {
	*x = StaticFileHandler{}
	mi := &file_core_v1_static_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *StaticFileHandler) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StaticFileHandler) ProtoMessage() {}

func (x *StaticFileHandler) ProtoReflect() protoreflect.Message {
	mi := &file_core_v1_static_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StaticFileHandler.ProtoReflect.Descriptor instead.
func (*StaticFileHandler) Descriptor() ([]byte, []int) {
	return file_core_v1_static_proto_rawDescGZIP(), []int{0}
}

func (x *StaticFileHandler) GetAPIVersion() string {
	if x != nil {
		return x.APIVersion
	}
	return ""
}

func (x *StaticFileHandler) GetKind() string {
	if x != nil {
		return x.Kind
	}
	return ""
}

func (x *StaticFileHandler) GetMetadata() *kernel.Metadata {
	if x != nil {
		return x.Metadata
	}
	return nil
}

func (x *StaticFileHandler) GetSpec() *StaticFileHandlerSpec {
	if x != nil {
		return x.Spec
	}
	return nil
}

// StaticFileHandlerSpec is the specifications for the StaticFileHandler object.
type StaticFileHandlerSpec struct {
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
	// Note that it depends on the multiplexer, or HTTP router,
	// that the server uses if this field is used or not.
	// Default is not set.
	Methods []HTTPMethod `protobuf:"varint,3,rep,packed,name=Methods,json=methods,proto3,enum=core.v1.HTTPMethod" json:"Methods,omitempty"`
	// [OPTIONAL]
	// RootDir is the root directry path
	// that is serverd by this static file server.
	// If not set, the current working directory "./" will be used.
	// Default is not set.
	RootDir string `protobuf:"bytes,4,opt,name=RootDir,json=rootDir,proto3" json:"RootDir,omitempty"`
	// [OPTIONAL]
	// StripPrefix is the prefix string to strip from the requested path.
	// For example, set "/foo/bar" to get "content.json" with the path "/foo/bar/content.json".
	// Default is not set.
	StripPrefix string `protobuf:"bytes,5,opt,name=StripPrefix,json=stripPrefix,proto3" json:"StripPrefix,omitempty"`
	// [OPTIONAL]
	// EnableListing is flag to enable directory listing under RootDir.
	// Setting this true can make the gateway vulnerable to directory listing attack.
	// Do not set this unless you know what you are doing.
	// Default is [false].
	EnableListing bool `protobuf:"varint,6,opt,name=EnableListing,json=enableListing,proto3" json:"EnableListing,omitempty"`
	// [OPTIONAL]
	// Header is the key-value pairs of HTTP headers
	// which are added to the all responses.
	// For example, headers for cache controls should be considered.
	// Content-Type header is recommended to be set when serving the same type contents
	// to avoid content detection in the gateway from the stand point view of performance.
	// Default is not set.
	Header        map[string]string `protobuf:"bytes,7,rep,name=Header,json=header,proto3" json:"Header,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *StaticFileHandlerSpec) Reset() {
	*x = StaticFileHandlerSpec{}
	mi := &file_core_v1_static_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *StaticFileHandlerSpec) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StaticFileHandlerSpec) ProtoMessage() {}

func (x *StaticFileHandlerSpec) ProtoReflect() protoreflect.Message {
	mi := &file_core_v1_static_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StaticFileHandlerSpec.ProtoReflect.Descriptor instead.
func (*StaticFileHandlerSpec) Descriptor() ([]byte, []int) {
	return file_core_v1_static_proto_rawDescGZIP(), []int{1}
}

func (x *StaticFileHandlerSpec) GetErrorHandler() *kernel.Reference {
	if x != nil {
		return x.ErrorHandler
	}
	return nil
}

func (x *StaticFileHandlerSpec) GetPatterns() []string {
	if x != nil {
		return x.Patterns
	}
	return nil
}

func (x *StaticFileHandlerSpec) GetMethods() []HTTPMethod {
	if x != nil {
		return x.Methods
	}
	return nil
}

func (x *StaticFileHandlerSpec) GetRootDir() string {
	if x != nil {
		return x.RootDir
	}
	return ""
}

func (x *StaticFileHandlerSpec) GetStripPrefix() string {
	if x != nil {
		return x.StripPrefix
	}
	return ""
}

func (x *StaticFileHandlerSpec) GetEnableListing() bool {
	if x != nil {
		return x.EnableListing
	}
	return false
}

func (x *StaticFileHandlerSpec) GetHeader() map[string]string {
	if x != nil {
		return x.Header
	}
	return nil
}

var File_core_v1_static_proto protoreflect.FileDescriptor

var file_core_v1_static_proto_rawDesc = string([]byte{
	0x0a, 0x14, 0x63, 0x6f, 0x72, 0x65, 0x2f, 0x76, 0x31, 0x2f, 0x73, 0x74, 0x61, 0x74, 0x69, 0x63,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x07, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x76, 0x31, 0x1a,
	0x1b, 0x62, 0x75, 0x66, 0x2f, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2f, 0x76, 0x61,
	0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x12, 0x63, 0x6f,
	0x72, 0x65, 0x2f, 0x76, 0x31, 0x2f, 0x68, 0x74, 0x74, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x1a, 0x15, 0x6b, 0x65, 0x72, 0x6e, 0x65, 0x6c, 0x2f, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xd3, 0x01, 0x0a, 0x11, 0x53, 0x74, 0x61, 0x74,
	0x69, 0x63, 0x46, 0x69, 0x6c, 0x65, 0x48, 0x61, 0x6e, 0x64, 0x6c, 0x65, 0x72, 0x12, 0x2e, 0x0a,
	0x0a, 0x41, 0x50, 0x49, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x42, 0x0e, 0xba, 0x48, 0x0b, 0x72, 0x09, 0x0a, 0x07, 0x63, 0x6f, 0x72, 0x65, 0x2f, 0x76,
	0x31, 0x52, 0x0a, 0x61, 0x70, 0x69, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x2c, 0x0a,
	0x04, 0x4b, 0x69, 0x6e, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x42, 0x18, 0xba, 0x48, 0x15,
	0x72, 0x13, 0x0a, 0x11, 0x53, 0x74, 0x61, 0x74, 0x69, 0x63, 0x46, 0x69, 0x6c, 0x65, 0x48, 0x61,
	0x6e, 0x64, 0x6c, 0x65, 0x72, 0x52, 0x04, 0x6b, 0x69, 0x6e, 0x64, 0x12, 0x2c, 0x0a, 0x08, 0x4d,
	0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x10, 0x2e,
	0x6b, 0x65, 0x72, 0x6e, 0x65, 0x6c, 0x2e, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x52,
	0x08, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x12, 0x32, 0x0a, 0x04, 0x53, 0x70, 0x65,
	0x63, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x76,
	0x31, 0x2e, 0x53, 0x74, 0x61, 0x74, 0x69, 0x63, 0x46, 0x69, 0x6c, 0x65, 0x48, 0x61, 0x6e, 0x64,
	0x6c, 0x65, 0x72, 0x53, 0x70, 0x65, 0x63, 0x52, 0x04, 0x73, 0x70, 0x65, 0x63, 0x22, 0x8e, 0x03,
	0x0a, 0x15, 0x53, 0x74, 0x61, 0x74, 0x69, 0x63, 0x46, 0x69, 0x6c, 0x65, 0x48, 0x61, 0x6e, 0x64,
	0x6c, 0x65, 0x72, 0x53, 0x70, 0x65, 0x63, 0x12, 0x35, 0x0a, 0x0c, 0x45, 0x72, 0x72, 0x6f, 0x72,
	0x48, 0x61, 0x6e, 0x64, 0x6c, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x11, 0x2e,
	0x6b, 0x65, 0x72, 0x6e, 0x65, 0x6c, 0x2e, 0x52, 0x65, 0x66, 0x65, 0x72, 0x65, 0x6e, 0x63, 0x65,
	0x52, 0x0c, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x48, 0x61, 0x6e, 0x64, 0x6c, 0x65, 0x72, 0x12, 0x24,
	0x0a, 0x08, 0x50, 0x61, 0x74, 0x74, 0x65, 0x72, 0x6e, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x09,
	0x42, 0x08, 0xba, 0x48, 0x05, 0x92, 0x01, 0x02, 0x18, 0x01, 0x52, 0x08, 0x70, 0x61, 0x74, 0x74,
	0x65, 0x72, 0x6e, 0x73, 0x12, 0x37, 0x0a, 0x07, 0x4d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x73, 0x18,
	0x03, 0x20, 0x03, 0x28, 0x0e, 0x32, 0x13, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x76, 0x31, 0x2e,
	0x48, 0x54, 0x54, 0x50, 0x4d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x42, 0x08, 0xba, 0x48, 0x05, 0x92,
	0x01, 0x02, 0x18, 0x01, 0x52, 0x07, 0x6d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x73, 0x12, 0x18, 0x0a,
	0x07, 0x52, 0x6f, 0x6f, 0x74, 0x44, 0x69, 0x72, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07,
	0x72, 0x6f, 0x6f, 0x74, 0x44, 0x69, 0x72, 0x12, 0x20, 0x0a, 0x0b, 0x53, 0x74, 0x72, 0x69, 0x70,
	0x50, 0x72, 0x65, 0x66, 0x69, 0x78, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x73, 0x74,
	0x72, 0x69, 0x70, 0x50, 0x72, 0x65, 0x66, 0x69, 0x78, 0x12, 0x24, 0x0a, 0x0d, 0x45, 0x6e, 0x61,
	0x62, 0x6c, 0x65, 0x4c, 0x69, 0x73, 0x74, 0x69, 0x6e, 0x67, 0x18, 0x06, 0x20, 0x01, 0x28, 0x08,
	0x52, 0x0d, 0x65, 0x6e, 0x61, 0x62, 0x6c, 0x65, 0x4c, 0x69, 0x73, 0x74, 0x69, 0x6e, 0x67, 0x12,
	0x42, 0x0a, 0x06, 0x48, 0x65, 0x61, 0x64, 0x65, 0x72, 0x18, 0x07, 0x20, 0x03, 0x28, 0x0b, 0x32,
	0x2a, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x74, 0x61, 0x74, 0x69, 0x63,
	0x46, 0x69, 0x6c, 0x65, 0x48, 0x61, 0x6e, 0x64, 0x6c, 0x65, 0x72, 0x53, 0x70, 0x65, 0x63, 0x2e,
	0x48, 0x65, 0x61, 0x64, 0x65, 0x72, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x06, 0x68, 0x65, 0x61,
	0x64, 0x65, 0x72, 0x1a, 0x39, 0x0a, 0x0b, 0x48, 0x65, 0x61, 0x64, 0x65, 0x72, 0x45, 0x6e, 0x74,
	0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x42, 0x39,
	0x5a, 0x37, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x61, 0x69, 0x6c,
	0x65, 0x72, 0x6f, 0x6e, 0x2d, 0x67, 0x61, 0x74, 0x65, 0x77, 0x61, 0x79, 0x2f, 0x61, 0x69, 0x6c,
	0x65, 0x72, 0x6f, 0x6e, 0x2d, 0x67, 0x61, 0x74, 0x65, 0x77, 0x61, 0x79, 0x2f, 0x61, 0x70, 0x69,
	0x73, 0x2f, 0x63, 0x6f, 0x72, 0x65, 0x2f, 0x76, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
})

var (
	file_core_v1_static_proto_rawDescOnce sync.Once
	file_core_v1_static_proto_rawDescData []byte
)

func file_core_v1_static_proto_rawDescGZIP() []byte {
	file_core_v1_static_proto_rawDescOnce.Do(func() {
		file_core_v1_static_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_core_v1_static_proto_rawDesc), len(file_core_v1_static_proto_rawDesc)))
	})
	return file_core_v1_static_proto_rawDescData
}

var file_core_v1_static_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_core_v1_static_proto_goTypes = []any{
	(*StaticFileHandler)(nil),     // 0: core.v1.StaticFileHandler
	(*StaticFileHandlerSpec)(nil), // 1: core.v1.StaticFileHandlerSpec
	nil,                           // 2: core.v1.StaticFileHandlerSpec.HeaderEntry
	(*kernel.Metadata)(nil),       // 3: kernel.Metadata
	(*kernel.Reference)(nil),      // 4: kernel.Reference
	(HTTPMethod)(0),               // 5: core.v1.HTTPMethod
}
var file_core_v1_static_proto_depIdxs = []int32{
	3, // 0: core.v1.StaticFileHandler.Metadata:type_name -> kernel.Metadata
	1, // 1: core.v1.StaticFileHandler.Spec:type_name -> core.v1.StaticFileHandlerSpec
	4, // 2: core.v1.StaticFileHandlerSpec.ErrorHandler:type_name -> kernel.Reference
	5, // 3: core.v1.StaticFileHandlerSpec.Methods:type_name -> core.v1.HTTPMethod
	2, // 4: core.v1.StaticFileHandlerSpec.Header:type_name -> core.v1.StaticFileHandlerSpec.HeaderEntry
	5, // [5:5] is the sub-list for method output_type
	5, // [5:5] is the sub-list for method input_type
	5, // [5:5] is the sub-list for extension type_name
	5, // [5:5] is the sub-list for extension extendee
	0, // [0:5] is the sub-list for field type_name
}

func init() { file_core_v1_static_proto_init() }
func file_core_v1_static_proto_init() {
	if File_core_v1_static_proto != nil {
		return
	}
	file_core_v1_http_proto_init()
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_core_v1_static_proto_rawDesc), len(file_core_v1_static_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_core_v1_static_proto_goTypes,
		DependencyIndexes: file_core_v1_static_proto_depIdxs,
		MessageInfos:      file_core_v1_static_proto_msgTypes,
	}.Build()
	File_core_v1_static_proto = out.File
	file_core_v1_static_proto_goTypes = nil
	file_core_v1_static_proto_depIdxs = nil
}
