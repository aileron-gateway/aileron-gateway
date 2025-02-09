// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.4
// 	protoc        v5.29.0
// source: app/v1/skipper.proto

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

// Skipper is the definition of the Skipper object.
// Skipper implements interface of the middleware.
type Skipper struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// [REQUIRED]
	// APIVersion is the defined version of the resource.
	// This value must be "app/v1".
	APIVersion string `protobuf:"bytes,1,opt,name=APIVersion,json=apiVersion,proto3" json:"APIVersion,omitempty"`
	// [REQUIRED]
	// Kind is the kind of this object.
	// This value must be "Skipper".
	Kind string `protobuf:"bytes,2,opt,name=Kind,json=kind,proto3" json:"Kind,omitempty"`
	// [OPTIONAL]
	// Metadata is the metadata of the http logger object.
	// If not set, both name and namespace in the metadata
	// are treated as "default".
	Metadata *kernel.Metadata `protobuf:"bytes,3,opt,name=Metadata,json=metadata,proto3" json:"Metadata,omitempty"`
	// [OPTIONAL]
	// Spec is the specification of the middleware.
	// Default values are used when nothing is set.
	Spec          *SkipperSpec `protobuf:"bytes,4,opt,name=Spec,json=spec,proto3" json:"Spec,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Skipper) Reset() {
	*x = Skipper{}
	mi := &file_app_v1_skipper_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Skipper) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Skipper) ProtoMessage() {}

func (x *Skipper) ProtoReflect() protoreflect.Message {
	mi := &file_app_v1_skipper_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Skipper.ProtoReflect.Descriptor instead.
func (*Skipper) Descriptor() ([]byte, []int) {
	return file_app_v1_skipper_proto_rawDescGZIP(), []int{0}
}

func (x *Skipper) GetAPIVersion() string {
	if x != nil {
		return x.APIVersion
	}
	return ""
}

func (x *Skipper) GetKind() string {
	if x != nil {
		return x.Kind
	}
	return ""
}

func (x *Skipper) GetMetadata() *kernel.Metadata {
	if x != nil {
		return x.Metadata
	}
	return nil
}

func (x *Skipper) GetSpec() *SkipperSpec {
	if x != nil {
		return x.Spec
	}
	return nil
}

// SkipperSpec is the specifications for the Skipper object.
type SkipperSpec struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// [OPTIONAL]
	// SkipConditions is the list of conditions
	// to skip the configured middleware or tripperware.
	// Default is not set.
	SkipConditions []*SkipConditionSpec `protobuf:"bytes,1,rep,name=SkipConditions,json=skipConditions,proto3" json:"SkipConditions,omitempty"`
	// [OPTIONAL]
	// Middleware is the list of references to middleware.
	// Configured middleware is skipped when the requests
	// matched to one of the skip conditions.
	// Default is not set.
	Middleware []*kernel.Reference `protobuf:"bytes,2,rep,name=Middleware,json=middleware,proto3" json:"Middleware,omitempty"`
	// [OPTIONAL]
	// Tripperware is the list of references to tripperware.
	// Configured tripperware is skipped when the requests
	// matched to one of the skip conditions.
	// Default is not set.
	Tripperware   []*kernel.Reference `protobuf:"bytes,3,rep,name=Tripperware,json=tripperware,proto3" json:"Tripperware,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *SkipperSpec) Reset() {
	*x = SkipperSpec{}
	mi := &file_app_v1_skipper_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SkipperSpec) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SkipperSpec) ProtoMessage() {}

func (x *SkipperSpec) ProtoReflect() protoreflect.Message {
	mi := &file_app_v1_skipper_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SkipperSpec.ProtoReflect.Descriptor instead.
func (*SkipperSpec) Descriptor() ([]byte, []int) {
	return file_app_v1_skipper_proto_rawDescGZIP(), []int{1}
}

func (x *SkipperSpec) GetSkipConditions() []*SkipConditionSpec {
	if x != nil {
		return x.SkipConditions
	}
	return nil
}

func (x *SkipperSpec) GetMiddleware() []*kernel.Reference {
	if x != nil {
		return x.Middleware
	}
	return nil
}

func (x *SkipperSpec) GetTripperware() []*kernel.Reference {
	if x != nil {
		return x.Tripperware
	}
	return nil
}

// SkipConditionSpec is the configuration spec for the matching conditions.
type SkipConditionSpec struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// [OPTIONAL]
	// Matcher is a matcher which matches to defined patterns.
	// Default is not set.
	Matcher *kernel.MatcherSpec `protobuf:"bytes,1,opt,name=Matcher,json=matcher,proto3" json:"Matcher,omitempty"`
	// [OPTIONAL]
	// Methods is the list of HTTP methods to be matched.
	// Default is not set.
	Methods       []v1.HTTPMethod `protobuf:"varint,2,rep,packed,name=Methods,json=methods,proto3,enum=core.v1.HTTPMethod" json:"Methods,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *SkipConditionSpec) Reset() {
	*x = SkipConditionSpec{}
	mi := &file_app_v1_skipper_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SkipConditionSpec) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SkipConditionSpec) ProtoMessage() {}

func (x *SkipConditionSpec) ProtoReflect() protoreflect.Message {
	mi := &file_app_v1_skipper_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SkipConditionSpec.ProtoReflect.Descriptor instead.
func (*SkipConditionSpec) Descriptor() ([]byte, []int) {
	return file_app_v1_skipper_proto_rawDescGZIP(), []int{2}
}

func (x *SkipConditionSpec) GetMatcher() *kernel.MatcherSpec {
	if x != nil {
		return x.Matcher
	}
	return nil
}

func (x *SkipConditionSpec) GetMethods() []v1.HTTPMethod {
	if x != nil {
		return x.Methods
	}
	return nil
}

var File_app_v1_skipper_proto protoreflect.FileDescriptor

var file_app_v1_skipper_proto_rawDesc = string([]byte{
	0x0a, 0x14, 0x61, 0x70, 0x70, 0x2f, 0x76, 0x31, 0x2f, 0x73, 0x6b, 0x69, 0x70, 0x70, 0x65, 0x72,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x06, 0x61, 0x70, 0x70, 0x2e, 0x76, 0x31, 0x1a, 0x1b,
	0x62, 0x75, 0x66, 0x2f, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2f, 0x76, 0x61, 0x6c,
	0x69, 0x64, 0x61, 0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x15, 0x6b, 0x65, 0x72,
	0x6e, 0x65, 0x6c, 0x2f, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x1a, 0x12, 0x63, 0x6f, 0x72, 0x65, 0x2f, 0x76, 0x31, 0x2f, 0x68, 0x74, 0x74, 0x70,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x14, 0x6b, 0x65, 0x72, 0x6e, 0x65, 0x6c, 0x2f, 0x74,
	0x78, 0x74, 0x75, 0x74, 0x69, 0x6c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xb3, 0x01, 0x0a,
	0x07, 0x53, 0x6b, 0x69, 0x70, 0x70, 0x65, 0x72, 0x12, 0x2d, 0x0a, 0x0a, 0x41, 0x50, 0x49, 0x56,
	0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x0d, 0xba, 0x48,
	0x0a, 0x72, 0x08, 0x0a, 0x06, 0x61, 0x70, 0x70, 0x2f, 0x76, 0x31, 0x52, 0x0a, 0x61, 0x70, 0x69,
	0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x22, 0x0a, 0x04, 0x4b, 0x69, 0x6e, 0x64, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x42, 0x0e, 0xba, 0x48, 0x0b, 0x72, 0x09, 0x0a, 0x07, 0x53, 0x6b,
	0x69, 0x70, 0x70, 0x65, 0x72, 0x52, 0x04, 0x6b, 0x69, 0x6e, 0x64, 0x12, 0x2c, 0x0a, 0x08, 0x4d,
	0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x10, 0x2e,
	0x6b, 0x65, 0x72, 0x6e, 0x65, 0x6c, 0x2e, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x52,
	0x08, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x12, 0x27, 0x0a, 0x04, 0x53, 0x70, 0x65,
	0x63, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x13, 0x2e, 0x61, 0x70, 0x70, 0x2e, 0x76, 0x31,
	0x2e, 0x53, 0x6b, 0x69, 0x70, 0x70, 0x65, 0x72, 0x53, 0x70, 0x65, 0x63, 0x52, 0x04, 0x73, 0x70,
	0x65, 0x63, 0x22, 0xb8, 0x01, 0x0a, 0x0b, 0x53, 0x6b, 0x69, 0x70, 0x70, 0x65, 0x72, 0x53, 0x70,
	0x65, 0x63, 0x12, 0x41, 0x0a, 0x0e, 0x53, 0x6b, 0x69, 0x70, 0x43, 0x6f, 0x6e, 0x64, 0x69, 0x74,
	0x69, 0x6f, 0x6e, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x61, 0x70, 0x70,
	0x2e, 0x76, 0x31, 0x2e, 0x53, 0x6b, 0x69, 0x70, 0x43, 0x6f, 0x6e, 0x64, 0x69, 0x74, 0x69, 0x6f,
	0x6e, 0x53, 0x70, 0x65, 0x63, 0x52, 0x0e, 0x73, 0x6b, 0x69, 0x70, 0x43, 0x6f, 0x6e, 0x64, 0x69,
	0x74, 0x69, 0x6f, 0x6e, 0x73, 0x12, 0x31, 0x0a, 0x0a, 0x4d, 0x69, 0x64, 0x64, 0x6c, 0x65, 0x77,
	0x61, 0x72, 0x65, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x11, 0x2e, 0x6b, 0x65, 0x72, 0x6e,
	0x65, 0x6c, 0x2e, 0x52, 0x65, 0x66, 0x65, 0x72, 0x65, 0x6e, 0x63, 0x65, 0x52, 0x0a, 0x6d, 0x69,
	0x64, 0x64, 0x6c, 0x65, 0x77, 0x61, 0x72, 0x65, 0x12, 0x33, 0x0a, 0x0b, 0x54, 0x72, 0x69, 0x70,
	0x70, 0x65, 0x72, 0x77, 0x61, 0x72, 0x65, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x11, 0x2e,
	0x6b, 0x65, 0x72, 0x6e, 0x65, 0x6c, 0x2e, 0x52, 0x65, 0x66, 0x65, 0x72, 0x65, 0x6e, 0x63, 0x65,
	0x52, 0x0b, 0x74, 0x72, 0x69, 0x70, 0x70, 0x65, 0x72, 0x77, 0x61, 0x72, 0x65, 0x22, 0x7b, 0x0a,
	0x11, 0x53, 0x6b, 0x69, 0x70, 0x43, 0x6f, 0x6e, 0x64, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x53, 0x70,
	0x65, 0x63, 0x12, 0x2d, 0x0a, 0x07, 0x4d, 0x61, 0x74, 0x63, 0x68, 0x65, 0x72, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x13, 0x2e, 0x6b, 0x65, 0x72, 0x6e, 0x65, 0x6c, 0x2e, 0x4d, 0x61, 0x74,
	0x63, 0x68, 0x65, 0x72, 0x53, 0x70, 0x65, 0x63, 0x52, 0x07, 0x6d, 0x61, 0x74, 0x63, 0x68, 0x65,
	0x72, 0x12, 0x37, 0x0a, 0x07, 0x4d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x73, 0x18, 0x02, 0x20, 0x03,
	0x28, 0x0e, 0x32, 0x13, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x48, 0x54, 0x54,
	0x50, 0x4d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x42, 0x08, 0xba, 0x48, 0x05, 0x92, 0x01, 0x02, 0x18,
	0x01, 0x52, 0x07, 0x6d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x73, 0x42, 0x38, 0x5a, 0x36, 0x67, 0x69,
	0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x61, 0x69, 0x6c, 0x65, 0x72, 0x6f, 0x6e,
	0x2d, 0x67, 0x61, 0x74, 0x65, 0x77, 0x61, 0x79, 0x2f, 0x61, 0x69, 0x6c, 0x65, 0x72, 0x6f, 0x6e,
	0x2d, 0x67, 0x61, 0x74, 0x65, 0x77, 0x61, 0x79, 0x2f, 0x61, 0x70, 0x69, 0x73, 0x2f, 0x61, 0x70,
	0x70, 0x2f, 0x76, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
})

var (
	file_app_v1_skipper_proto_rawDescOnce sync.Once
	file_app_v1_skipper_proto_rawDescData []byte
)

func file_app_v1_skipper_proto_rawDescGZIP() []byte {
	file_app_v1_skipper_proto_rawDescOnce.Do(func() {
		file_app_v1_skipper_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_app_v1_skipper_proto_rawDesc), len(file_app_v1_skipper_proto_rawDesc)))
	})
	return file_app_v1_skipper_proto_rawDescData
}

var file_app_v1_skipper_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_app_v1_skipper_proto_goTypes = []any{
	(*Skipper)(nil),            // 0: app.v1.Skipper
	(*SkipperSpec)(nil),        // 1: app.v1.SkipperSpec
	(*SkipConditionSpec)(nil),  // 2: app.v1.SkipConditionSpec
	(*kernel.Metadata)(nil),    // 3: kernel.Metadata
	(*kernel.Reference)(nil),   // 4: kernel.Reference
	(*kernel.MatcherSpec)(nil), // 5: kernel.MatcherSpec
	(v1.HTTPMethod)(0),         // 6: core.v1.HTTPMethod
}
var file_app_v1_skipper_proto_depIdxs = []int32{
	3, // 0: app.v1.Skipper.Metadata:type_name -> kernel.Metadata
	1, // 1: app.v1.Skipper.Spec:type_name -> app.v1.SkipperSpec
	2, // 2: app.v1.SkipperSpec.SkipConditions:type_name -> app.v1.SkipConditionSpec
	4, // 3: app.v1.SkipperSpec.Middleware:type_name -> kernel.Reference
	4, // 4: app.v1.SkipperSpec.Tripperware:type_name -> kernel.Reference
	5, // 5: app.v1.SkipConditionSpec.Matcher:type_name -> kernel.MatcherSpec
	6, // 6: app.v1.SkipConditionSpec.Methods:type_name -> core.v1.HTTPMethod
	7, // [7:7] is the sub-list for method output_type
	7, // [7:7] is the sub-list for method input_type
	7, // [7:7] is the sub-list for extension type_name
	7, // [7:7] is the sub-list for extension extendee
	0, // [0:7] is the sub-list for field type_name
}

func init() { file_app_v1_skipper_proto_init() }
func file_app_v1_skipper_proto_init() {
	if File_app_v1_skipper_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_app_v1_skipper_proto_rawDesc), len(file_app_v1_skipper_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_app_v1_skipper_proto_goTypes,
		DependencyIndexes: file_app_v1_skipper_proto_depIdxs,
		MessageInfos:      file_app_v1_skipper_proto_msgTypes,
	}.Build()
	File_app_v1_skipper_proto = out.File
	file_app_v1_skipper_proto_goTypes = nil
	file_app_v1_skipper_proto_depIdxs = nil
}
