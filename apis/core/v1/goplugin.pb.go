// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.1
// 	protoc        v5.29.0
// source: core/v1/goplugin.proto

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

// GoPlugin is the definition of the GoPlugin object.
type GoPlugin struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// [REQUIRED]
	// APIVersion is the defined version of the logger.
	// This value must be "core/v1".
	APIVersion string `protobuf:"bytes,1,opt,name=APIVersion,json=apiVersion,proto3" json:"APIVersion,omitempty"`
	// [REQUIRED]
	// Kind is the kind of this object.
	// This value must be "GoPlugin".
	Kind string `protobuf:"bytes,2,opt,name=Kind,json=kind,proto3" json:"Kind,omitempty"`
	// [OPTIONAL]
	// Metadata is the metadata of the middleware logger.
	// If not set, both name and namespace in the metadata
	// are treated as "default".
	Metadata *kernel.Metadata `protobuf:"bytes,3,opt,name=Metadata,json=metadata,proto3" json:"Metadata,omitempty"`
	// [OPTIONAL]
	// Spec is the specification of the GoPlugin.
	// Default values are used when nothing is set.
	Spec          *GoPluginSpec `protobuf:"bytes,4,opt,name=Spec,json=spec,proto3" json:"Spec,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GoPlugin) Reset() {
	*x = GoPlugin{}
	mi := &file_core_v1_goplugin_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GoPlugin) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GoPlugin) ProtoMessage() {}

func (x *GoPlugin) ProtoReflect() protoreflect.Message {
	mi := &file_core_v1_goplugin_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GoPlugin.ProtoReflect.Descriptor instead.
func (*GoPlugin) Descriptor() ([]byte, []int) {
	return file_core_v1_goplugin_proto_rawDescGZIP(), []int{0}
}

func (x *GoPlugin) GetAPIVersion() string {
	if x != nil {
		return x.APIVersion
	}
	return ""
}

func (x *GoPlugin) GetKind() string {
	if x != nil {
		return x.Kind
	}
	return ""
}

func (x *GoPlugin) GetMetadata() *kernel.Metadata {
	if x != nil {
		return x.Metadata
	}
	return nil
}

func (x *GoPlugin) GetSpec() *GoPluginSpec {
	if x != nil {
		return x.Spec
	}
	return nil
}

// GoPluginSpec is the specifications of the GoPlugin object.
type GoPluginSpec struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// [OPTIONAL]
	// ErrorHandler is the reference to a ErrorHandler object.
	// Referred object must implement ErrorHandler interface.
	// Default error handler is used when not set.
	ErrorHandler *kernel.Reference `protobuf:"bytes,1,opt,name=ErrorHandler,json=errorHandler,proto3" json:"ErrorHandler,omitempty"`
	// [REQUIRED]
	// PluginPath is path to the shared object of the GoPlugin.
	// The path can be absolute or relative.
	// Default is not set.
	PluginPath string `protobuf:"bytes,2,opt,name=PluginPath,json=pluginPath,proto3" json:"PluginPath,omitempty"`
	// [OPTIONAL]
	// SymbolName is synbol name to lookup.
	// Default is ["Plugin"].
	SymbolName    string `protobuf:"bytes,3,opt,name=SymbolName,json=symbolName,proto3" json:"SymbolName,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GoPluginSpec) Reset() {
	*x = GoPluginSpec{}
	mi := &file_core_v1_goplugin_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GoPluginSpec) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GoPluginSpec) ProtoMessage() {}

func (x *GoPluginSpec) ProtoReflect() protoreflect.Message {
	mi := &file_core_v1_goplugin_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GoPluginSpec.ProtoReflect.Descriptor instead.
func (*GoPluginSpec) Descriptor() ([]byte, []int) {
	return file_core_v1_goplugin_proto_rawDescGZIP(), []int{1}
}

func (x *GoPluginSpec) GetErrorHandler() *kernel.Reference {
	if x != nil {
		return x.ErrorHandler
	}
	return nil
}

func (x *GoPluginSpec) GetPluginPath() string {
	if x != nil {
		return x.PluginPath
	}
	return ""
}

func (x *GoPluginSpec) GetSymbolName() string {
	if x != nil {
		return x.SymbolName
	}
	return ""
}

var File_core_v1_goplugin_proto protoreflect.FileDescriptor

var file_core_v1_goplugin_proto_rawDesc = []byte{
	0x0a, 0x16, 0x63, 0x6f, 0x72, 0x65, 0x2f, 0x76, 0x31, 0x2f, 0x67, 0x6f, 0x70, 0x6c, 0x75, 0x67,
	0x69, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x07, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x76,
	0x31, 0x1a, 0x1b, 0x62, 0x75, 0x66, 0x2f, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2f,
	0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x15,
	0x6b, 0x65, 0x72, 0x6e, 0x65, 0x6c, 0x2f, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xb8, 0x01, 0x0a, 0x08, 0x47, 0x6f, 0x50, 0x6c, 0x75, 0x67,
	0x69, 0x6e, 0x12, 0x2e, 0x0a, 0x0a, 0x41, 0x50, 0x49, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x0e, 0xba, 0x48, 0x0b, 0x72, 0x09, 0x0a, 0x07, 0x63,
	0x6f, 0x72, 0x65, 0x2f, 0x76, 0x31, 0x52, 0x0a, 0x61, 0x70, 0x69, 0x56, 0x65, 0x72, 0x73, 0x69,
	0x6f, 0x6e, 0x12, 0x23, 0x0a, 0x04, 0x4b, 0x69, 0x6e, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x42, 0x0f, 0xba, 0x48, 0x0c, 0x72, 0x0a, 0x0a, 0x08, 0x47, 0x6f, 0x50, 0x6c, 0x75, 0x67, 0x69,
	0x6e, 0x52, 0x04, 0x6b, 0x69, 0x6e, 0x64, 0x12, 0x2c, 0x0a, 0x08, 0x4d, 0x65, 0x74, 0x61, 0x64,
	0x61, 0x74, 0x61, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x6b, 0x65, 0x72, 0x6e,
	0x65, 0x6c, 0x2e, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x52, 0x08, 0x6d, 0x65, 0x74,
	0x61, 0x64, 0x61, 0x74, 0x61, 0x12, 0x29, 0x0a, 0x04, 0x53, 0x70, 0x65, 0x63, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x15, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x47, 0x6f,
	0x50, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x53, 0x70, 0x65, 0x63, 0x52, 0x04, 0x73, 0x70, 0x65, 0x63,
	0x22, 0x85, 0x01, 0x0a, 0x0c, 0x47, 0x6f, 0x50, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x53, 0x70, 0x65,
	0x63, 0x12, 0x35, 0x0a, 0x0c, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x48, 0x61, 0x6e, 0x64, 0x6c, 0x65,
	0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x11, 0x2e, 0x6b, 0x65, 0x72, 0x6e, 0x65, 0x6c,
	0x2e, 0x52, 0x65, 0x66, 0x65, 0x72, 0x65, 0x6e, 0x63, 0x65, 0x52, 0x0c, 0x65, 0x72, 0x72, 0x6f,
	0x72, 0x48, 0x61, 0x6e, 0x64, 0x6c, 0x65, 0x72, 0x12, 0x1e, 0x0a, 0x0a, 0x50, 0x6c, 0x75, 0x67,
	0x69, 0x6e, 0x50, 0x61, 0x74, 0x68, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x70, 0x6c,
	0x75, 0x67, 0x69, 0x6e, 0x50, 0x61, 0x74, 0x68, 0x12, 0x1e, 0x0a, 0x0a, 0x53, 0x79, 0x6d, 0x62,
	0x6f, 0x6c, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x73, 0x79,
	0x6d, 0x62, 0x6f, 0x6c, 0x4e, 0x61, 0x6d, 0x65, 0x42, 0x39, 0x5a, 0x37, 0x67, 0x69, 0x74, 0x68,
	0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x61, 0x69, 0x6c, 0x65, 0x72, 0x6f, 0x6e, 0x2d, 0x67,
	0x61, 0x74, 0x65, 0x77, 0x61, 0x79, 0x2f, 0x61, 0x69, 0x6c, 0x65, 0x72, 0x6f, 0x6e, 0x2d, 0x67,
	0x61, 0x74, 0x65, 0x77, 0x61, 0x79, 0x2f, 0x61, 0x70, 0x69, 0x73, 0x2f, 0x63, 0x6f, 0x72, 0x65,
	0x2f, 0x76, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_core_v1_goplugin_proto_rawDescOnce sync.Once
	file_core_v1_goplugin_proto_rawDescData = file_core_v1_goplugin_proto_rawDesc
)

func file_core_v1_goplugin_proto_rawDescGZIP() []byte {
	file_core_v1_goplugin_proto_rawDescOnce.Do(func() {
		file_core_v1_goplugin_proto_rawDescData = protoimpl.X.CompressGZIP(file_core_v1_goplugin_proto_rawDescData)
	})
	return file_core_v1_goplugin_proto_rawDescData
}

var file_core_v1_goplugin_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_core_v1_goplugin_proto_goTypes = []any{
	(*GoPlugin)(nil),         // 0: core.v1.GoPlugin
	(*GoPluginSpec)(nil),     // 1: core.v1.GoPluginSpec
	(*kernel.Metadata)(nil),  // 2: kernel.Metadata
	(*kernel.Reference)(nil), // 3: kernel.Reference
}
var file_core_v1_goplugin_proto_depIdxs = []int32{
	2, // 0: core.v1.GoPlugin.Metadata:type_name -> kernel.Metadata
	1, // 1: core.v1.GoPlugin.Spec:type_name -> core.v1.GoPluginSpec
	3, // 2: core.v1.GoPluginSpec.ErrorHandler:type_name -> kernel.Reference
	3, // [3:3] is the sub-list for method output_type
	3, // [3:3] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_core_v1_goplugin_proto_init() }
func file_core_v1_goplugin_proto_init() {
	if File_core_v1_goplugin_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_core_v1_goplugin_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_core_v1_goplugin_proto_goTypes,
		DependencyIndexes: file_core_v1_goplugin_proto_depIdxs,
		MessageInfos:      file_core_v1_goplugin_proto_msgTypes,
	}.Build()
	File_core_v1_goplugin_proto = out.File
	file_core_v1_goplugin_proto_rawDesc = nil
	file_core_v1_goplugin_proto_goTypes = nil
	file_core_v1_goplugin_proto_depIdxs = nil
}
