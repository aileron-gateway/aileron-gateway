// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.4
// 	protoc        v5.29.0
// source: app/v1/meter/prommeter.proto

package v1

import (
	_ "buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go/buf/validate"
	v1 "github.com/aileron-gateway/aileron-gateway/apis/core/v1"
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

// PrometheusMeter is the definition of the PrometheusMeter object.
// PrometheusMeter implements interface of the HTTP handler, middleware.
type PrometheusMeter struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// [REQUIRED]
	// APIVersion is the defined version of the handler.
	// This value must be "app/v1".
	APIVersion string `protobuf:"bytes,1,opt,name=APIVersion,json=apiVersion,proto3" json:"APIVersion,omitempty"`
	// [REQUIRED]
	// Kind is the kind of this object.
	// This value must be "PrometheusHandler".
	Kind string `protobuf:"bytes,2,opt,name=Kind,json=kind,proto3" json:"Kind,omitempty"`
	// [OPTIONAL]
	// Metadata is the metadata of the http logger object.
	// If not set, both name and namespace in the metadata
	// are treated as "default".
	Metadata *kernel.Metadata `protobuf:"bytes,3,opt,name=Metadata,json=metadata,proto3" json:"Metadata,omitempty"`
	// [OPTIONAL]
	// Spec is the specification of the meter.
	// Default values are used when nothing is set.
	Spec          *PrometheusMeterSpec `protobuf:"bytes,4,opt,name=Spec,json=spec,proto3" json:"Spec,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *PrometheusMeter) Reset() {
	*x = PrometheusMeter{}
	mi := &file_app_v1_meter_prommeter_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *PrometheusMeter) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PrometheusMeter) ProtoMessage() {}

func (x *PrometheusMeter) ProtoReflect() protoreflect.Message {
	mi := &file_app_v1_meter_prommeter_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PrometheusMeter.ProtoReflect.Descriptor instead.
func (*PrometheusMeter) Descriptor() ([]byte, []int) {
	return file_app_v1_meter_prommeter_proto_rawDescGZIP(), []int{0}
}

func (x *PrometheusMeter) GetAPIVersion() string {
	if x != nil {
		return x.APIVersion
	}
	return ""
}

func (x *PrometheusMeter) GetKind() string {
	if x != nil {
		return x.Kind
	}
	return ""
}

func (x *PrometheusMeter) GetMetadata() *kernel.Metadata {
	if x != nil {
		return x.Metadata
	}
	return nil
}

func (x *PrometheusMeter) GetSpec() *PrometheusMeterSpec {
	if x != nil {
		return x.Spec
	}
	return nil
}

// PrometheusMeterSpec is the specifications for the PrometheusHandler object.
type PrometheusMeterSpec struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// [OPTIONAL]
	// Patterns is path patterns that this handler
	// is registered to a server.
	// Default is not set, or empty string ["/"].
	Patterns []string `protobuf:"bytes,10,rep,name=Patterns,json=patterns,proto3" json:"Patterns,omitempty"`
	// [OPTIONAL]
	// Methods is the list of HTTP method this handler can handle.
	// Note that it depends on the multiplexer, or HTTP router,
	// that the server uses if this field is used or not.
	// Default multiplexer does not use this field.
	// Default is not set.
	Methods       []v1.HTTPMethod `protobuf:"varint,11,rep,packed,name=Methods,json=methods,proto3,enum=core.v1.HTTPMethod" json:"Methods,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *PrometheusMeterSpec) Reset() {
	*x = PrometheusMeterSpec{}
	mi := &file_app_v1_meter_prommeter_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *PrometheusMeterSpec) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PrometheusMeterSpec) ProtoMessage() {}

func (x *PrometheusMeterSpec) ProtoReflect() protoreflect.Message {
	mi := &file_app_v1_meter_prommeter_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PrometheusMeterSpec.ProtoReflect.Descriptor instead.
func (*PrometheusMeterSpec) Descriptor() ([]byte, []int) {
	return file_app_v1_meter_prommeter_proto_rawDescGZIP(), []int{1}
}

func (x *PrometheusMeterSpec) GetPatterns() []string {
	if x != nil {
		return x.Patterns
	}
	return nil
}

func (x *PrometheusMeterSpec) GetMethods() []v1.HTTPMethod {
	if x != nil {
		return x.Methods
	}
	return nil
}

var File_app_v1_meter_prommeter_proto protoreflect.FileDescriptor

var file_app_v1_meter_prommeter_proto_rawDesc = string([]byte{
	0x0a, 0x1c, 0x61, 0x70, 0x70, 0x2f, 0x76, 0x31, 0x2f, 0x6d, 0x65, 0x74, 0x65, 0x72, 0x2f, 0x70,
	0x72, 0x6f, 0x6d, 0x6d, 0x65, 0x74, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x06,
	0x61, 0x70, 0x70, 0x2e, 0x76, 0x31, 0x1a, 0x1b, 0x62, 0x75, 0x66, 0x2f, 0x76, 0x61, 0x6c, 0x69,
	0x64, 0x61, 0x74, 0x65, 0x2f, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x1a, 0x12, 0x63, 0x6f, 0x72, 0x65, 0x2f, 0x76, 0x31, 0x2f, 0x68, 0x74, 0x74,
	0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x15, 0x6b, 0x65, 0x72, 0x6e, 0x65, 0x6c, 0x2f,
	0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xcb,
	0x01, 0x0a, 0x0f, 0x50, 0x72, 0x6f, 0x6d, 0x65, 0x74, 0x68, 0x65, 0x75, 0x73, 0x4d, 0x65, 0x74,
	0x65, 0x72, 0x12, 0x2d, 0x0a, 0x0a, 0x41, 0x50, 0x49, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x0d, 0xba, 0x48, 0x0a, 0x72, 0x08, 0x0a, 0x06, 0x61,
	0x70, 0x70, 0x2f, 0x76, 0x31, 0x52, 0x0a, 0x61, 0x70, 0x69, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f,
	0x6e, 0x12, 0x2a, 0x0a, 0x04, 0x4b, 0x69, 0x6e, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x42,
	0x16, 0xba, 0x48, 0x13, 0x72, 0x11, 0x0a, 0x0f, 0x50, 0x72, 0x6f, 0x6d, 0x65, 0x74, 0x68, 0x65,
	0x75, 0x73, 0x4d, 0x65, 0x74, 0x65, 0x72, 0x52, 0x04, 0x6b, 0x69, 0x6e, 0x64, 0x12, 0x2c, 0x0a,
	0x08, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x10, 0x2e, 0x6b, 0x65, 0x72, 0x6e, 0x65, 0x6c, 0x2e, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74,
	0x61, 0x52, 0x08, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x12, 0x2f, 0x0a, 0x04, 0x53,
	0x70, 0x65, 0x63, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1b, 0x2e, 0x61, 0x70, 0x70, 0x2e,
	0x76, 0x31, 0x2e, 0x50, 0x72, 0x6f, 0x6d, 0x65, 0x74, 0x68, 0x65, 0x75, 0x73, 0x4d, 0x65, 0x74,
	0x65, 0x72, 0x53, 0x70, 0x65, 0x63, 0x52, 0x04, 0x73, 0x70, 0x65, 0x63, 0x22, 0x60, 0x0a, 0x13,
	0x50, 0x72, 0x6f, 0x6d, 0x65, 0x74, 0x68, 0x65, 0x75, 0x73, 0x4d, 0x65, 0x74, 0x65, 0x72, 0x53,
	0x70, 0x65, 0x63, 0x12, 0x1a, 0x0a, 0x08, 0x50, 0x61, 0x74, 0x74, 0x65, 0x72, 0x6e, 0x73, 0x18,
	0x0a, 0x20, 0x03, 0x28, 0x09, 0x52, 0x08, 0x70, 0x61, 0x74, 0x74, 0x65, 0x72, 0x6e, 0x73, 0x12,
	0x2d, 0x0a, 0x07, 0x4d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x73, 0x18, 0x0b, 0x20, 0x03, 0x28, 0x0e,
	0x32, 0x13, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x48, 0x54, 0x54, 0x50, 0x4d,
	0x65, 0x74, 0x68, 0x6f, 0x64, 0x52, 0x07, 0x6d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x73, 0x42, 0x38,
	0x5a, 0x36, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x61, 0x69, 0x6c,
	0x65, 0x72, 0x6f, 0x6e, 0x2d, 0x67, 0x61, 0x74, 0x65, 0x77, 0x61, 0x79, 0x2f, 0x61, 0x69, 0x6c,
	0x65, 0x72, 0x6f, 0x6e, 0x2d, 0x67, 0x61, 0x74, 0x65, 0x77, 0x61, 0x79, 0x2f, 0x61, 0x70, 0x69,
	0x73, 0x2f, 0x61, 0x70, 0x70, 0x2f, 0x76, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
})

var (
	file_app_v1_meter_prommeter_proto_rawDescOnce sync.Once
	file_app_v1_meter_prommeter_proto_rawDescData []byte
)

func file_app_v1_meter_prommeter_proto_rawDescGZIP() []byte {
	file_app_v1_meter_prommeter_proto_rawDescOnce.Do(func() {
		file_app_v1_meter_prommeter_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_app_v1_meter_prommeter_proto_rawDesc), len(file_app_v1_meter_prommeter_proto_rawDesc)))
	})
	return file_app_v1_meter_prommeter_proto_rawDescData
}

var file_app_v1_meter_prommeter_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_app_v1_meter_prommeter_proto_goTypes = []any{
	(*PrometheusMeter)(nil),     // 0: app.v1.PrometheusMeter
	(*PrometheusMeterSpec)(nil), // 1: app.v1.PrometheusMeterSpec
	(*kernel.Metadata)(nil),     // 2: kernel.Metadata
	(v1.HTTPMethod)(0),          // 3: core.v1.HTTPMethod
}
var file_app_v1_meter_prommeter_proto_depIdxs = []int32{
	2, // 0: app.v1.PrometheusMeter.Metadata:type_name -> kernel.Metadata
	1, // 1: app.v1.PrometheusMeter.Spec:type_name -> app.v1.PrometheusMeterSpec
	3, // 2: app.v1.PrometheusMeterSpec.Methods:type_name -> core.v1.HTTPMethod
	3, // [3:3] is the sub-list for method output_type
	3, // [3:3] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_app_v1_meter_prommeter_proto_init() }
func file_app_v1_meter_prommeter_proto_init() {
	if File_app_v1_meter_prommeter_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_app_v1_meter_prommeter_proto_rawDesc), len(file_app_v1_meter_prommeter_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_app_v1_meter_prommeter_proto_goTypes,
		DependencyIndexes: file_app_v1_meter_prommeter_proto_depIdxs,
		MessageInfos:      file_app_v1_meter_prommeter_proto_msgTypes,
	}.Build()
	File_app_v1_meter_prommeter_proto = out.File
	file_app_v1_meter_prommeter_proto_goTypes = nil
	file_app_v1_meter_prommeter_proto_depIdxs = nil
}
