// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.3
// 	protoc        v5.29.0
// source: core/v1/entrypoint.proto

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

// Entrypoint is the definition of the Entrypoint object.
// Entrypoint implements interface of the service.
type Entrypoint struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// [REQUIRED]
	// APIVersion is the defined version of the service.
	// This value must be "core/v1".
	APIVersion string `protobuf:"bytes,1,opt,name=APIVersion,json=apiVersion,proto3" json:"APIVersion,omitempty"`
	// [REQUIRED]
	// Kind is the kind of this object.
	// This value must be "Entrypoint".
	Kind string `protobuf:"bytes,2,opt,name=Kind,json=kind,proto3" json:"Kind,omitempty"`
	// [OPTIONAL]
	// Metadata is the metadata of the entrypoint object.
	// If not set, both name and namespace in the metadata
	// are treated as "default".
	Metadata *kernel.Metadata `protobuf:"bytes,3,opt,name=Metadata,json=metadata,proto3" json:"Metadata,omitempty"`
	// [OPTIONAL]
	// Spec is the specification of the service.
	// Default values are used when nothing is set.
	Spec          *EntrypointSpec `protobuf:"bytes,4,opt,name=Spec,json=spec,proto3" json:"Spec,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Entrypoint) Reset() {
	*x = Entrypoint{}
	mi := &file_core_v1_entrypoint_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Entrypoint) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Entrypoint) ProtoMessage() {}

func (x *Entrypoint) ProtoReflect() protoreflect.Message {
	mi := &file_core_v1_entrypoint_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Entrypoint.ProtoReflect.Descriptor instead.
func (*Entrypoint) Descriptor() ([]byte, []int) {
	return file_core_v1_entrypoint_proto_rawDescGZIP(), []int{0}
}

func (x *Entrypoint) GetAPIVersion() string {
	if x != nil {
		return x.APIVersion
	}
	return ""
}

func (x *Entrypoint) GetKind() string {
	if x != nil {
		return x.Kind
	}
	return ""
}

func (x *Entrypoint) GetMetadata() *kernel.Metadata {
	if x != nil {
		return x.Metadata
	}
	return nil
}

func (x *Entrypoint) GetSpec() *EntrypointSpec {
	if x != nil {
		return x.Spec
	}
	return nil
}

// EntrypointSpec is the specifications for the Entrypoint object.
type EntrypointSpec struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// [OPTIONAL]
	// DefaultLogger is the reference to a Logger object
	// that will be used from other resources by default.
	// Object must implement Logger interface.
	// Default Logger is used when not set.
	DefaultLogger *kernel.Reference `protobuf:"bytes,1,opt,name=DefaultLogger,json=defaultLogger,proto3" json:"DefaultLogger,omitempty"`
	// [OPTIONAL]
	// Loggers is the list of references to Logger objects.
	// The specified loggers can be referred from other resources
	// by the name "<group>/<version>/<kind>/<namespace>/<name>"
	// from their metadata field
	// for example "core/v1/SLogger/myNamespace/myLogger".
	// Objects must implement Logger interface.
	// Default is not set.
	Loggers []*kernel.Reference `protobuf:"bytes,2,rep,name=Loggers,json=loggers,proto3" json:"Loggers,omitempty"`
	// [OPTIONAL]
	// DefaultErrorHandler is the reference to a ErrorHandler object
	// that will be used from other resources by default.
	// Referred object must implement ErrorHandler interface.
	// Default ErrorHandler is used when not set.
	DefaultErrorHandler *kernel.Reference `protobuf:"bytes,3,opt,name=DefaultErrorHandler,json=defaultErrorHandler,proto3" json:"DefaultErrorHandler,omitempty"`
	// [OPTIONAL]
	// Runners is the list of reference to runner resources.
	// Referred objects must implement the Runner interface.
	// HTTPServer is a typical example of the resource that can be set to this field.
	// The gateway will exit without doing anything
	// when no runner was specified.
	// The gateway will exit with failure when one of or all of the
	// runners exit with an error.
	Runners []*kernel.Reference `protobuf:"bytes,4,rep,name=Runners,json=runners,proto3" json:"Runners,omitempty"`
	// [OPTIONAL]
	// WaitAll is the flag to wait all runners to exit.
	// If false, the gateway will exit when at least 1 runner exited.
	// If true, the gateway will exit when all of the runners exited.
	// Default is [false].
	WaitAll bool `protobuf:"varint,5,opt,name=WaitAll,json=waitAll,proto3" json:"WaitAll,omitempty"`
	// [OPTIONAL]
	// Initializers is the reference to the resources
	// that should be initialized before creating runners.
	// Referred objects must implement the Initializer interface.
	// Default is not set.
	Initializers []*kernel.Reference `protobuf:"bytes,6,rep,name=Initializers,json=initializers,proto3" json:"Initializers,omitempty"`
	// [OPTIONAL]
	// Finalizers is the reference to the resources
	// that should be finalized on exit of the gateway.
	// Referred objects must implement the Finalizer interface.
	// Default is not set.
	Finalizers    []*kernel.Reference `protobuf:"bytes,7,rep,name=Finalizers,json=finalizers,proto3" json:"Finalizers,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *EntrypointSpec) Reset() {
	*x = EntrypointSpec{}
	mi := &file_core_v1_entrypoint_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *EntrypointSpec) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EntrypointSpec) ProtoMessage() {}

func (x *EntrypointSpec) ProtoReflect() protoreflect.Message {
	mi := &file_core_v1_entrypoint_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EntrypointSpec.ProtoReflect.Descriptor instead.
func (*EntrypointSpec) Descriptor() ([]byte, []int) {
	return file_core_v1_entrypoint_proto_rawDescGZIP(), []int{1}
}

func (x *EntrypointSpec) GetDefaultLogger() *kernel.Reference {
	if x != nil {
		return x.DefaultLogger
	}
	return nil
}

func (x *EntrypointSpec) GetLoggers() []*kernel.Reference {
	if x != nil {
		return x.Loggers
	}
	return nil
}

func (x *EntrypointSpec) GetDefaultErrorHandler() *kernel.Reference {
	if x != nil {
		return x.DefaultErrorHandler
	}
	return nil
}

func (x *EntrypointSpec) GetRunners() []*kernel.Reference {
	if x != nil {
		return x.Runners
	}
	return nil
}

func (x *EntrypointSpec) GetWaitAll() bool {
	if x != nil {
		return x.WaitAll
	}
	return false
}

func (x *EntrypointSpec) GetInitializers() []*kernel.Reference {
	if x != nil {
		return x.Initializers
	}
	return nil
}

func (x *EntrypointSpec) GetFinalizers() []*kernel.Reference {
	if x != nil {
		return x.Finalizers
	}
	return nil
}

var File_core_v1_entrypoint_proto protoreflect.FileDescriptor

var file_core_v1_entrypoint_proto_rawDesc = []byte{
	0x0a, 0x18, 0x63, 0x6f, 0x72, 0x65, 0x2f, 0x76, 0x31, 0x2f, 0x65, 0x6e, 0x74, 0x72, 0x79, 0x70,
	0x6f, 0x69, 0x6e, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x07, 0x63, 0x6f, 0x72, 0x65,
	0x2e, 0x76, 0x31, 0x1a, 0x1b, 0x62, 0x75, 0x66, 0x2f, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74,
	0x65, 0x2f, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x1a, 0x15, 0x6b, 0x65, 0x72, 0x6e, 0x65, 0x6c, 0x2f, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xbe, 0x01, 0x0a, 0x0a, 0x45, 0x6e, 0x74, 0x72,
	0x79, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x12, 0x2e, 0x0a, 0x0a, 0x41, 0x50, 0x49, 0x56, 0x65, 0x72,
	0x73, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x0e, 0xba, 0x48, 0x0b, 0x72,
	0x09, 0x0a, 0x07, 0x63, 0x6f, 0x72, 0x65, 0x2f, 0x76, 0x31, 0x52, 0x0a, 0x61, 0x70, 0x69, 0x56,
	0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x25, 0x0a, 0x04, 0x4b, 0x69, 0x6e, 0x64, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x42, 0x11, 0xba, 0x48, 0x0e, 0x72, 0x0c, 0x0a, 0x0a, 0x45, 0x6e, 0x74,
	0x72, 0x79, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x52, 0x04, 0x6b, 0x69, 0x6e, 0x64, 0x12, 0x2c, 0x0a,
	0x08, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x10, 0x2e, 0x6b, 0x65, 0x72, 0x6e, 0x65, 0x6c, 0x2e, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74,
	0x61, 0x52, 0x08, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x12, 0x2b, 0x0a, 0x04, 0x53,
	0x70, 0x65, 0x63, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x63, 0x6f, 0x72, 0x65,
	0x2e, 0x76, 0x31, 0x2e, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x53, 0x70,
	0x65, 0x63, 0x52, 0x04, 0x73, 0x70, 0x65, 0x63, 0x22, 0xec, 0x02, 0x0a, 0x0e, 0x45, 0x6e, 0x74,
	0x72, 0x79, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x53, 0x70, 0x65, 0x63, 0x12, 0x37, 0x0a, 0x0d, 0x44,
	0x65, 0x66, 0x61, 0x75, 0x6c, 0x74, 0x4c, 0x6f, 0x67, 0x67, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x11, 0x2e, 0x6b, 0x65, 0x72, 0x6e, 0x65, 0x6c, 0x2e, 0x52, 0x65, 0x66, 0x65,
	0x72, 0x65, 0x6e, 0x63, 0x65, 0x52, 0x0d, 0x64, 0x65, 0x66, 0x61, 0x75, 0x6c, 0x74, 0x4c, 0x6f,
	0x67, 0x67, 0x65, 0x72, 0x12, 0x2b, 0x0a, 0x07, 0x4c, 0x6f, 0x67, 0x67, 0x65, 0x72, 0x73, 0x18,
	0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x11, 0x2e, 0x6b, 0x65, 0x72, 0x6e, 0x65, 0x6c, 0x2e, 0x52,
	0x65, 0x66, 0x65, 0x72, 0x65, 0x6e, 0x63, 0x65, 0x52, 0x07, 0x6c, 0x6f, 0x67, 0x67, 0x65, 0x72,
	0x73, 0x12, 0x43, 0x0a, 0x13, 0x44, 0x65, 0x66, 0x61, 0x75, 0x6c, 0x74, 0x45, 0x72, 0x72, 0x6f,
	0x72, 0x48, 0x61, 0x6e, 0x64, 0x6c, 0x65, 0x72, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x11,
	0x2e, 0x6b, 0x65, 0x72, 0x6e, 0x65, 0x6c, 0x2e, 0x52, 0x65, 0x66, 0x65, 0x72, 0x65, 0x6e, 0x63,
	0x65, 0x52, 0x13, 0x64, 0x65, 0x66, 0x61, 0x75, 0x6c, 0x74, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x48,
	0x61, 0x6e, 0x64, 0x6c, 0x65, 0x72, 0x12, 0x2b, 0x0a, 0x07, 0x52, 0x75, 0x6e, 0x6e, 0x65, 0x72,
	0x73, 0x18, 0x04, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x11, 0x2e, 0x6b, 0x65, 0x72, 0x6e, 0x65, 0x6c,
	0x2e, 0x52, 0x65, 0x66, 0x65, 0x72, 0x65, 0x6e, 0x63, 0x65, 0x52, 0x07, 0x72, 0x75, 0x6e, 0x6e,
	0x65, 0x72, 0x73, 0x12, 0x18, 0x0a, 0x07, 0x57, 0x61, 0x69, 0x74, 0x41, 0x6c, 0x6c, 0x18, 0x05,
	0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x77, 0x61, 0x69, 0x74, 0x41, 0x6c, 0x6c, 0x12, 0x35, 0x0a,
	0x0c, 0x49, 0x6e, 0x69, 0x74, 0x69, 0x61, 0x6c, 0x69, 0x7a, 0x65, 0x72, 0x73, 0x18, 0x06, 0x20,
	0x03, 0x28, 0x0b, 0x32, 0x11, 0x2e, 0x6b, 0x65, 0x72, 0x6e, 0x65, 0x6c, 0x2e, 0x52, 0x65, 0x66,
	0x65, 0x72, 0x65, 0x6e, 0x63, 0x65, 0x52, 0x0c, 0x69, 0x6e, 0x69, 0x74, 0x69, 0x61, 0x6c, 0x69,
	0x7a, 0x65, 0x72, 0x73, 0x12, 0x31, 0x0a, 0x0a, 0x46, 0x69, 0x6e, 0x61, 0x6c, 0x69, 0x7a, 0x65,
	0x72, 0x73, 0x18, 0x07, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x11, 0x2e, 0x6b, 0x65, 0x72, 0x6e, 0x65,
	0x6c, 0x2e, 0x52, 0x65, 0x66, 0x65, 0x72, 0x65, 0x6e, 0x63, 0x65, 0x52, 0x0a, 0x66, 0x69, 0x6e,
	0x61, 0x6c, 0x69, 0x7a, 0x65, 0x72, 0x73, 0x42, 0x39, 0x5a, 0x37, 0x67, 0x69, 0x74, 0x68, 0x75,
	0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x61, 0x69, 0x6c, 0x65, 0x72, 0x6f, 0x6e, 0x2d, 0x67, 0x61,
	0x74, 0x65, 0x77, 0x61, 0x79, 0x2f, 0x61, 0x69, 0x6c, 0x65, 0x72, 0x6f, 0x6e, 0x2d, 0x67, 0x61,
	0x74, 0x65, 0x77, 0x61, 0x79, 0x2f, 0x61, 0x70, 0x69, 0x73, 0x2f, 0x63, 0x6f, 0x72, 0x65, 0x2f,
	0x76, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_core_v1_entrypoint_proto_rawDescOnce sync.Once
	file_core_v1_entrypoint_proto_rawDescData = file_core_v1_entrypoint_proto_rawDesc
)

func file_core_v1_entrypoint_proto_rawDescGZIP() []byte {
	file_core_v1_entrypoint_proto_rawDescOnce.Do(func() {
		file_core_v1_entrypoint_proto_rawDescData = protoimpl.X.CompressGZIP(file_core_v1_entrypoint_proto_rawDescData)
	})
	return file_core_v1_entrypoint_proto_rawDescData
}

var file_core_v1_entrypoint_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_core_v1_entrypoint_proto_goTypes = []any{
	(*Entrypoint)(nil),       // 0: core.v1.Entrypoint
	(*EntrypointSpec)(nil),   // 1: core.v1.EntrypointSpec
	(*kernel.Metadata)(nil),  // 2: kernel.Metadata
	(*kernel.Reference)(nil), // 3: kernel.Reference
}
var file_core_v1_entrypoint_proto_depIdxs = []int32{
	2, // 0: core.v1.Entrypoint.Metadata:type_name -> kernel.Metadata
	1, // 1: core.v1.Entrypoint.Spec:type_name -> core.v1.EntrypointSpec
	3, // 2: core.v1.EntrypointSpec.DefaultLogger:type_name -> kernel.Reference
	3, // 3: core.v1.EntrypointSpec.Loggers:type_name -> kernel.Reference
	3, // 4: core.v1.EntrypointSpec.DefaultErrorHandler:type_name -> kernel.Reference
	3, // 5: core.v1.EntrypointSpec.Runners:type_name -> kernel.Reference
	3, // 6: core.v1.EntrypointSpec.Initializers:type_name -> kernel.Reference
	3, // 7: core.v1.EntrypointSpec.Finalizers:type_name -> kernel.Reference
	8, // [8:8] is the sub-list for method output_type
	8, // [8:8] is the sub-list for method input_type
	8, // [8:8] is the sub-list for extension type_name
	8, // [8:8] is the sub-list for extension extendee
	0, // [0:8] is the sub-list for field type_name
}

func init() { file_core_v1_entrypoint_proto_init() }
func file_core_v1_entrypoint_proto_init() {
	if File_core_v1_entrypoint_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_core_v1_entrypoint_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_core_v1_entrypoint_proto_goTypes,
		DependencyIndexes: file_core_v1_entrypoint_proto_depIdxs,
		MessageInfos:      file_core_v1_entrypoint_proto_msgTypes,
	}.Build()
	File_core_v1_entrypoint_proto = out.File
	file_core_v1_entrypoint_proto_rawDesc = nil
	file_core_v1_entrypoint_proto_goTypes = nil
	file_core_v1_entrypoint_proto_depIdxs = nil
}
