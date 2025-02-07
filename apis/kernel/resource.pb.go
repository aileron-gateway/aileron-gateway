// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.5
// 	protoc        v5.29.0
// source: kernel/resource.proto

package kernel

import (
	_ "buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go/buf/validate"
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

// Metadata is the metadata of resources.
type Metadata struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// [OPTIONAL]
	// Name is the unique name for the resource.
	// The name must be unique in a namespace.
	// If not set, the name "default" is used.
	// Default is not set.
	Name string `protobuf:"bytes,1,opt,name=Name,json=name,proto3" json:"Name,omitempty"`
	// [OPTIONAL]
	// Namespace is the name of a logical group of resources.
	// If not set, the name "default" is used.
	// Default is not set.
	Namespace string `protobuf:"bytes,2,opt,name=Namespace,json=namespace,proto3" json:"Namespace,omitempty"`
	// [OPTIONAL]
	// Logger is the logger name to use.
	// The format should be "<group>/<version>/<kind>/<namespace>/<name>"
	// for example "core/v1/SLogger/myNamespace/myLogger".
	// The named logger must be registered to the entrypoint resource.
	// "__default__" is the special name to refer to the default logger.
	// The default logger, which is globally unique, is used
	// if logger was not found by the specified name.
	// If not set, the globally unique default logger is used.
	// Logger for access logs, authentication or authorization logs, etc
	// should be configured under spec.
	// Default is not set.
	Logger        string `protobuf:"bytes,3,opt,name=Logger,json=logger,proto3" json:"Logger,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Metadata) Reset() {
	*x = Metadata{}
	mi := &file_kernel_resource_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Metadata) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Metadata) ProtoMessage() {}

func (x *Metadata) ProtoReflect() protoreflect.Message {
	mi := &file_kernel_resource_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Metadata.ProtoReflect.Descriptor instead.
func (*Metadata) Descriptor() ([]byte, []int) {
	return file_kernel_resource_proto_rawDescGZIP(), []int{0}
}

func (x *Metadata) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Metadata) GetNamespace() string {
	if x != nil {
		return x.Namespace
	}
	return ""
}

func (x *Metadata) GetLogger() string {
	if x != nil {
		return x.Logger
	}
	return ""
}

// Reference is the reference to a resource instance.
type Reference struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// [REQUIRED]
	// APIVersion is the api version of the resource instance to fetch.
	// This field must not be empty.
	APIVersion string `protobuf:"bytes,1,opt,name=APIVersion,json=apiVersion,proto3" json:"APIVersion,omitempty"`
	// [REQUIRED]
	// Kind is the api kind of the instance.
	// This field must not be empty.
	Kind string `protobuf:"bytes,2,opt,name=Kind,json=kind,proto3" json:"Kind,omitempty"`
	// [OPTIONAL]
	// Name is the name of the object to fetch.
	// If not set, "default" is used.
	Name string `protobuf:"bytes,3,opt,name=Name,json=name,proto3" json:"Name,omitempty"`
	// [OPTIONAL]
	// Namespace is the namespace name where to fetch the instance from.
	// If not set, "default" is used.
	Namespace     string `protobuf:"bytes,4,opt,name=Namespace,json=namespace,proto3" json:"Namespace,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Reference) Reset() {
	*x = Reference{}
	mi := &file_kernel_resource_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Reference) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Reference) ProtoMessage() {}

func (x *Reference) ProtoReflect() protoreflect.Message {
	mi := &file_kernel_resource_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Reference.ProtoReflect.Descriptor instead.
func (*Reference) Descriptor() ([]byte, []int) {
	return file_kernel_resource_proto_rawDescGZIP(), []int{1}
}

func (x *Reference) GetAPIVersion() string {
	if x != nil {
		return x.APIVersion
	}
	return ""
}

func (x *Reference) GetKind() string {
	if x != nil {
		return x.Kind
	}
	return ""
}

func (x *Reference) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Reference) GetNamespace() string {
	if x != nil {
		return x.Namespace
	}
	return ""
}

// Status is the status of resources.
// This message is defined for future
// use and have no effests currently.
type Status struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Status) Reset() {
	*x = Status{}
	mi := &file_kernel_resource_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Status) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Status) ProtoMessage() {}

func (x *Status) ProtoReflect() protoreflect.Message {
	mi := &file_kernel_resource_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Status.ProtoReflect.Descriptor instead.
func (*Status) Descriptor() ([]byte, []int) {
	return file_kernel_resource_proto_rawDescGZIP(), []int{2}
}

// Template is the template manifest of resources.
// This messages is used internally in the gateway.
type Template struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// [REQUIRED]
	// APIVersion is the version string of the API.
	APIVersion string `protobuf:"bytes,1,opt,name=APIVersion,json=apiVersion,proto3" json:"APIVersion,omitempty"`
	// [REQUIRED]
	// Kind is the API kind.
	Kind string `protobuf:"bytes,2,opt,name=Kind,json=kind,proto3" json:"Kind,omitempty"`
	// [OPTIONAL]
	// Metadata is the metadata of the resource.
	Metadata *Metadata `protobuf:"bytes,3,opt,name=Metadata,json=metadata,proto3" json:"Metadata,omitempty"`
	// [OPTIONAL]
	// Spec is the specificationa of the resource.
	Spec          *TemplateSpec `protobuf:"bytes,4,opt,name=Spec,json=spec,proto3" json:"Spec,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Template) Reset() {
	*x = Template{}
	mi := &file_kernel_resource_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Template) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Template) ProtoMessage() {}

func (x *Template) ProtoReflect() protoreflect.Message {
	mi := &file_kernel_resource_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Template.ProtoReflect.Descriptor instead.
func (*Template) Descriptor() ([]byte, []int) {
	return file_kernel_resource_proto_rawDescGZIP(), []int{3}
}

func (x *Template) GetAPIVersion() string {
	if x != nil {
		return x.APIVersion
	}
	return ""
}

func (x *Template) GetKind() string {
	if x != nil {
		return x.Kind
	}
	return ""
}

func (x *Template) GetMetadata() *Metadata {
	if x != nil {
		return x.Metadata
	}
	return nil
}

func (x *Template) GetSpec() *TemplateSpec {
	if x != nil {
		return x.Spec
	}
	return nil
}

// TemplateSpec is the specifications for the Template.
// This messages is used internally in the gateway.
type TemplateSpec struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *TemplateSpec) Reset() {
	*x = TemplateSpec{}
	mi := &file_kernel_resource_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *TemplateSpec) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TemplateSpec) ProtoMessage() {}

func (x *TemplateSpec) ProtoReflect() protoreflect.Message {
	mi := &file_kernel_resource_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TemplateSpec.ProtoReflect.Descriptor instead.
func (*TemplateSpec) Descriptor() ([]byte, []int) {
	return file_kernel_resource_proto_rawDescGZIP(), []int{4}
}

var File_kernel_resource_proto protoreflect.FileDescriptor

var file_kernel_resource_proto_rawDesc = string([]byte{
	0x0a, 0x15, 0x6b, 0x65, 0x72, 0x6e, 0x65, 0x6c, 0x2f, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x06, 0x6b, 0x65, 0x72, 0x6e, 0x65, 0x6c, 0x1a,
	0x1b, 0x62, 0x75, 0x66, 0x2f, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2f, 0x76, 0x61,
	0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x86, 0x01, 0x0a,
	0x08, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x12, 0x2b, 0x0a, 0x04, 0x4e, 0x61, 0x6d,
	0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x17, 0xba, 0x48, 0x14, 0x72, 0x12, 0x32, 0x10,
	0x5e, 0x5b, 0x30, 0x2d, 0x39, 0x41, 0x2d, 0x5a, 0x61, 0x2d, 0x7a, 0x2d, 0x2e, 0x5d, 0x2a, 0x24,
	0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x35, 0x0a, 0x09, 0x4e, 0x61, 0x6d, 0x65, 0x73, 0x70,
	0x61, 0x63, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x42, 0x17, 0xba, 0x48, 0x14, 0x72, 0x12,
	0x32, 0x10, 0x5e, 0x5b, 0x30, 0x2d, 0x39, 0x41, 0x2d, 0x5a, 0x61, 0x2d, 0x7a, 0x2d, 0x2e, 0x5d,
	0x2a, 0x24, 0x52, 0x09, 0x6e, 0x61, 0x6d, 0x65, 0x73, 0x70, 0x61, 0x63, 0x65, 0x12, 0x16, 0x0a,
	0x06, 0x4c, 0x6f, 0x67, 0x67, 0x65, 0x72, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x6c,
	0x6f, 0x67, 0x67, 0x65, 0x72, 0x22, 0xde, 0x01, 0x0a, 0x09, 0x52, 0x65, 0x66, 0x65, 0x72, 0x65,
	0x6e, 0x63, 0x65, 0x12, 0x42, 0x0a, 0x0a, 0x41, 0x50, 0x49, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f,
	0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x22, 0xba, 0x48, 0x1f, 0x72, 0x1d, 0x32, 0x1b,
	0x5e, 0x5b, 0x30, 0x2d, 0x39, 0x41, 0x2d, 0x5a, 0x61, 0x2d, 0x7a, 0x5d, 0x2b, 0x2f, 0x5b, 0x30,
	0x2d, 0x39, 0x41, 0x2d, 0x5a, 0x61, 0x2d, 0x7a, 0x5d, 0x2b, 0x24, 0x52, 0x0a, 0x61, 0x70, 0x69,
	0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x29, 0x0a, 0x04, 0x4b, 0x69, 0x6e, 0x64, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x42, 0x15, 0xba, 0x48, 0x12, 0x72, 0x10, 0x32, 0x0e, 0x5e, 0x5b,
	0x30, 0x2d, 0x39, 0x41, 0x2d, 0x5a, 0x61, 0x2d, 0x7a, 0x5d, 0x2b, 0x24, 0x52, 0x04, 0x6b, 0x69,
	0x6e, 0x64, 0x12, 0x2b, 0x0a, 0x04, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09,
	0x42, 0x17, 0xba, 0x48, 0x14, 0x72, 0x12, 0x32, 0x10, 0x5e, 0x5b, 0x30, 0x2d, 0x39, 0x41, 0x2d,
	0x5a, 0x61, 0x2d, 0x7a, 0x2d, 0x2e, 0x5d, 0x2a, 0x24, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12,
	0x35, 0x0a, 0x09, 0x4e, 0x61, 0x6d, 0x65, 0x73, 0x70, 0x61, 0x63, 0x65, 0x18, 0x04, 0x20, 0x01,
	0x28, 0x09, 0x42, 0x17, 0xba, 0x48, 0x14, 0x72, 0x12, 0x32, 0x10, 0x5e, 0x5b, 0x30, 0x2d, 0x39,
	0x41, 0x2d, 0x5a, 0x61, 0x2d, 0x7a, 0x2d, 0x2e, 0x5d, 0x2a, 0x24, 0x52, 0x09, 0x6e, 0x61, 0x6d,
	0x65, 0x73, 0x70, 0x61, 0x63, 0x65, 0x22, 0x08, 0x0a, 0x06, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73,
	0x22, 0x96, 0x01, 0x0a, 0x08, 0x54, 0x65, 0x6d, 0x70, 0x6c, 0x61, 0x74, 0x65, 0x12, 0x1e, 0x0a,
	0x0a, 0x41, 0x50, 0x49, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0a, 0x61, 0x70, 0x69, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x12, 0x0a,
	0x04, 0x4b, 0x69, 0x6e, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6b, 0x69, 0x6e,
	0x64, 0x12, 0x2c, 0x0a, 0x08, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x6b, 0x65, 0x72, 0x6e, 0x65, 0x6c, 0x2e, 0x4d, 0x65, 0x74,
	0x61, 0x64, 0x61, 0x74, 0x61, 0x52, 0x08, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x12,
	0x28, 0x0a, 0x04, 0x53, 0x70, 0x65, 0x63, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e,
	0x6b, 0x65, 0x72, 0x6e, 0x65, 0x6c, 0x2e, 0x54, 0x65, 0x6d, 0x70, 0x6c, 0x61, 0x74, 0x65, 0x53,
	0x70, 0x65, 0x63, 0x52, 0x04, 0x73, 0x70, 0x65, 0x63, 0x22, 0x0e, 0x0a, 0x0c, 0x54, 0x65, 0x6d,
	0x70, 0x6c, 0x61, 0x74, 0x65, 0x53, 0x70, 0x65, 0x63, 0x42, 0x38, 0x5a, 0x36, 0x67, 0x69, 0x74,
	0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x61, 0x69, 0x6c, 0x65, 0x72, 0x6f, 0x6e, 0x2d,
	0x67, 0x61, 0x74, 0x65, 0x77, 0x61, 0x79, 0x2f, 0x61, 0x69, 0x6c, 0x65, 0x72, 0x6f, 0x6e, 0x2d,
	0x67, 0x61, 0x74, 0x65, 0x77, 0x61, 0x79, 0x2f, 0x61, 0x70, 0x69, 0x73, 0x2f, 0x6b, 0x65, 0x72,
	0x6e, 0x65, 0x6c, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
})

var (
	file_kernel_resource_proto_rawDescOnce sync.Once
	file_kernel_resource_proto_rawDescData []byte
)

func file_kernel_resource_proto_rawDescGZIP() []byte {
	file_kernel_resource_proto_rawDescOnce.Do(func() {
		file_kernel_resource_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_kernel_resource_proto_rawDesc), len(file_kernel_resource_proto_rawDesc)))
	})
	return file_kernel_resource_proto_rawDescData
}

var file_kernel_resource_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_kernel_resource_proto_goTypes = []any{
	(*Metadata)(nil),     // 0: kernel.Metadata
	(*Reference)(nil),    // 1: kernel.Reference
	(*Status)(nil),       // 2: kernel.Status
	(*Template)(nil),     // 3: kernel.Template
	(*TemplateSpec)(nil), // 4: kernel.TemplateSpec
}
var file_kernel_resource_proto_depIdxs = []int32{
	0, // 0: kernel.Template.Metadata:type_name -> kernel.Metadata
	4, // 1: kernel.Template.Spec:type_name -> kernel.TemplateSpec
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_kernel_resource_proto_init() }
func file_kernel_resource_proto_init() {
	if File_kernel_resource_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_kernel_resource_proto_rawDesc), len(file_kernel_resource_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_kernel_resource_proto_goTypes,
		DependencyIndexes: file_kernel_resource_proto_depIdxs,
		MessageInfos:      file_kernel_resource_proto_msgTypes,
	}.Build()
	File_kernel_resource_proto = out.File
	file_kernel_resource_proto_goTypes = nil
	file_kernel_resource_proto_depIdxs = nil
}
