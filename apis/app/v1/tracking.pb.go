// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.4
// 	protoc        v5.29.0
// source: app/v1/middleware/tracking.proto

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

// TrackingMiddleware is the definition of the TrackingMiddleware object.
// TrackingMiddleware implements interface of the middleware.
type TrackingMiddleware struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// [REQUIRED]
	// APIVersion is the defined version of the midleware.
	// This value must be "app/v1".
	APIVersion string `protobuf:"bytes,1,opt,name=APIVersion,json=apiVersion,proto3" json:"APIVersion,omitempty"`
	// [REQUIRED]
	// Kind is the kind of this object.
	// This value must be "TrackingMiddleware".
	Kind string `protobuf:"bytes,2,opt,name=Kind,json=kind,proto3" json:"Kind,omitempty"`
	// [OPTIONAL]
	// Metadata is the metadata of the http logger object.
	// If not set, both name and namespace in the metadata
	// are treated as "default".
	Metadata *kernel.Metadata `protobuf:"bytes,3,opt,name=Metadata,json=metadata,proto3" json:"Metadata,omitempty"`
	// [OPTIONAL]
	// Spec is the specification of the middleware.
	// Default values are used when nothing is set.
	Spec          *TrackingMiddlewareSpec `protobuf:"bytes,4,opt,name=Spec,json=spec,proto3" json:"Spec,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *TrackingMiddleware) Reset() {
	*x = TrackingMiddleware{}
	mi := &file_app_v1_middleware_tracking_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *TrackingMiddleware) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TrackingMiddleware) ProtoMessage() {}

func (x *TrackingMiddleware) ProtoReflect() protoreflect.Message {
	mi := &file_app_v1_middleware_tracking_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TrackingMiddleware.ProtoReflect.Descriptor instead.
func (*TrackingMiddleware) Descriptor() ([]byte, []int) {
	return file_app_v1_middleware_tracking_proto_rawDescGZIP(), []int{0}
}

func (x *TrackingMiddleware) GetAPIVersion() string {
	if x != nil {
		return x.APIVersion
	}
	return ""
}

func (x *TrackingMiddleware) GetKind() string {
	if x != nil {
		return x.Kind
	}
	return ""
}

func (x *TrackingMiddleware) GetMetadata() *kernel.Metadata {
	if x != nil {
		return x.Metadata
	}
	return nil
}

func (x *TrackingMiddleware) GetSpec() *TrackingMiddlewareSpec {
	if x != nil {
		return x.Spec
	}
	return nil
}

// TrackingMiddlewareSpec is the specifications of the TrackingMiddleware object.
type TrackingMiddlewareSpec struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// [OPTIONAL]
	// ErrorHandler is the reference to a ErrorHandler object.
	// Referred object must implement ErrorHandler interface.
	// Default error handler is used when not set.
	ErrorHandler *kernel.Reference `protobuf:"bytes,1,opt,name=ErrorHandler,json=errorHandler,proto3" json:"ErrorHandler,omitempty"`
	// [OPTIONAL]
	// Encoding is the type of encoding used to generate IDs.
	// Default is [Base32HexEscaped].
	Encoding kernel.EncodingType `protobuf:"varint,2,opt,name=Encoding,json=encoding,proto3,enum=kernel.EncodingType" json:"Encoding,omitempty"`
	// [OPTIONAL]
	// RequestIDProxyName is the HTTP header name to proxy request ID.
	// If not set, request ID is not proxied.
	// Default is not set.
	RequestIDProxyName string `protobuf:"bytes,3,opt,name=RequestIDProxyName,json=requestIDProxyName,proto3" json:"RequestIDProxyName,omitempty"`
	// [OPTIONAL]
	// TraceIDProxyName is the HTTP header name to proxy trace ID.
	// If not set, trace ID is not proxied.
	// Default is not set.
	TraceIDProxyName string `protobuf:"bytes,4,opt,name=TraceIDProxyName,json=traceIDProxyName,proto3" json:"TraceIDProxyName,omitempty"`
	// [OPTIONAL]
	// TraceIDExtractName is the HTTP header name to extrace
	// a trace ID from the request headers.
	// If not set, a newly generated trace ID is always used.
	// Default is not set.
	TraceIDExtractName string `protobuf:"bytes,5,opt,name=TraceIDExtractName,json=traceIDExtractName,proto3" json:"TraceIDExtractName,omitempty"`
	unknownFields      protoimpl.UnknownFields
	sizeCache          protoimpl.SizeCache
}

func (x *TrackingMiddlewareSpec) Reset() {
	*x = TrackingMiddlewareSpec{}
	mi := &file_app_v1_middleware_tracking_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *TrackingMiddlewareSpec) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TrackingMiddlewareSpec) ProtoMessage() {}

func (x *TrackingMiddlewareSpec) ProtoReflect() protoreflect.Message {
	mi := &file_app_v1_middleware_tracking_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TrackingMiddlewareSpec.ProtoReflect.Descriptor instead.
func (*TrackingMiddlewareSpec) Descriptor() ([]byte, []int) {
	return file_app_v1_middleware_tracking_proto_rawDescGZIP(), []int{1}
}

func (x *TrackingMiddlewareSpec) GetErrorHandler() *kernel.Reference {
	if x != nil {
		return x.ErrorHandler
	}
	return nil
}

func (x *TrackingMiddlewareSpec) GetEncoding() kernel.EncodingType {
	if x != nil {
		return x.Encoding
	}
	return kernel.EncodingType(0)
}

func (x *TrackingMiddlewareSpec) GetRequestIDProxyName() string {
	if x != nil {
		return x.RequestIDProxyName
	}
	return ""
}

func (x *TrackingMiddlewareSpec) GetTraceIDProxyName() string {
	if x != nil {
		return x.TraceIDProxyName
	}
	return ""
}

func (x *TrackingMiddlewareSpec) GetTraceIDExtractName() string {
	if x != nil {
		return x.TraceIDExtractName
	}
	return ""
}

var File_app_v1_middleware_tracking_proto protoreflect.FileDescriptor

var file_app_v1_middleware_tracking_proto_rawDesc = string([]byte{
	0x0a, 0x20, 0x61, 0x70, 0x70, 0x2f, 0x76, 0x31, 0x2f, 0x6d, 0x69, 0x64, 0x64, 0x6c, 0x65, 0x77,
	0x61, 0x72, 0x65, 0x2f, 0x74, 0x72, 0x61, 0x63, 0x6b, 0x69, 0x6e, 0x67, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x06, 0x61, 0x70, 0x70, 0x2e, 0x76, 0x31, 0x1a, 0x1b, 0x62, 0x75, 0x66, 0x2f,
	0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2f, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x15, 0x6b, 0x65, 0x72, 0x6e, 0x65, 0x6c, 0x2f,
	0x65, 0x6e, 0x63, 0x6f, 0x64, 0x69, 0x6e, 0x67, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x15,
	0x6b, 0x65, 0x72, 0x6e, 0x65, 0x6c, 0x2f, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xd4, 0x01, 0x0a, 0x12, 0x54, 0x72, 0x61, 0x63, 0x6b, 0x69,
	0x6e, 0x67, 0x4d, 0x69, 0x64, 0x64, 0x6c, 0x65, 0x77, 0x61, 0x72, 0x65, 0x12, 0x2d, 0x0a, 0x0a,
	0x41, 0x50, 0x49, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x42, 0x0d, 0xba, 0x48, 0x0a, 0x72, 0x08, 0x0a, 0x06, 0x61, 0x70, 0x70, 0x2f, 0x76, 0x31, 0x52,
	0x0a, 0x61, 0x70, 0x69, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x2d, 0x0a, 0x04, 0x4b,
	0x69, 0x6e, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x42, 0x19, 0xba, 0x48, 0x16, 0x72, 0x14,
	0x0a, 0x12, 0x54, 0x72, 0x61, 0x63, 0x6b, 0x69, 0x6e, 0x67, 0x4d, 0x69, 0x64, 0x64, 0x6c, 0x65,
	0x77, 0x61, 0x72, 0x65, 0x52, 0x04, 0x6b, 0x69, 0x6e, 0x64, 0x12, 0x2c, 0x0a, 0x08, 0x4d, 0x65,
	0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x6b,
	0x65, 0x72, 0x6e, 0x65, 0x6c, 0x2e, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x52, 0x08,
	0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x12, 0x32, 0x0a, 0x04, 0x53, 0x70, 0x65, 0x63,
	0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x61, 0x70, 0x70, 0x2e, 0x76, 0x31, 0x2e,
	0x54, 0x72, 0x61, 0x63, 0x6b, 0x69, 0x6e, 0x67, 0x4d, 0x69, 0x64, 0x64, 0x6c, 0x65, 0x77, 0x61,
	0x72, 0x65, 0x53, 0x70, 0x65, 0x63, 0x52, 0x04, 0x73, 0x70, 0x65, 0x63, 0x22, 0xd8, 0x02, 0x0a,
	0x16, 0x54, 0x72, 0x61, 0x63, 0x6b, 0x69, 0x6e, 0x67, 0x4d, 0x69, 0x64, 0x64, 0x6c, 0x65, 0x77,
	0x61, 0x72, 0x65, 0x53, 0x70, 0x65, 0x63, 0x12, 0x35, 0x0a, 0x0c, 0x45, 0x72, 0x72, 0x6f, 0x72,
	0x48, 0x61, 0x6e, 0x64, 0x6c, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x11, 0x2e,
	0x6b, 0x65, 0x72, 0x6e, 0x65, 0x6c, 0x2e, 0x52, 0x65, 0x66, 0x65, 0x72, 0x65, 0x6e, 0x63, 0x65,
	0x52, 0x0c, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x48, 0x61, 0x6e, 0x64, 0x6c, 0x65, 0x72, 0x12, 0x30,
	0x0a, 0x08, 0x45, 0x6e, 0x63, 0x6f, 0x64, 0x69, 0x6e, 0x67, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0e,
	0x32, 0x14, 0x2e, 0x6b, 0x65, 0x72, 0x6e, 0x65, 0x6c, 0x2e, 0x45, 0x6e, 0x63, 0x6f, 0x64, 0x69,
	0x6e, 0x67, 0x54, 0x79, 0x70, 0x65, 0x52, 0x08, 0x65, 0x6e, 0x63, 0x6f, 0x64, 0x69, 0x6e, 0x67,
	0x12, 0x47, 0x0a, 0x12, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x49, 0x44, 0x50, 0x72, 0x6f,
	0x78, 0x79, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x42, 0x17, 0xba, 0x48,
	0x14, 0x72, 0x12, 0x32, 0x10, 0x5e, 0x5b, 0x30, 0x2d, 0x39, 0x61, 0x2d, 0x7a, 0x41, 0x2d, 0x5a,
	0x2d, 0x5f, 0x5d, 0x2a, 0x24, 0x52, 0x12, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x49, 0x44,
	0x50, 0x72, 0x6f, 0x78, 0x79, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x43, 0x0a, 0x10, 0x54, 0x72, 0x61,
	0x63, 0x65, 0x49, 0x44, 0x50, 0x72, 0x6f, 0x78, 0x79, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x09, 0x42, 0x17, 0xba, 0x48, 0x14, 0x72, 0x12, 0x32, 0x10, 0x5e, 0x5b, 0x30, 0x2d,
	0x39, 0x61, 0x2d, 0x7a, 0x41, 0x2d, 0x5a, 0x2d, 0x5f, 0x5d, 0x2a, 0x24, 0x52, 0x10, 0x74, 0x72,
	0x61, 0x63, 0x65, 0x49, 0x44, 0x50, 0x72, 0x6f, 0x78, 0x79, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x47,
	0x0a, 0x12, 0x54, 0x72, 0x61, 0x63, 0x65, 0x49, 0x44, 0x45, 0x78, 0x74, 0x72, 0x61, 0x63, 0x74,
	0x4e, 0x61, 0x6d, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x42, 0x17, 0xba, 0x48, 0x14, 0x72,
	0x12, 0x32, 0x10, 0x5e, 0x5b, 0x30, 0x2d, 0x39, 0x61, 0x2d, 0x7a, 0x41, 0x2d, 0x5a, 0x2d, 0x5f,
	0x5d, 0x2a, 0x24, 0x52, 0x12, 0x74, 0x72, 0x61, 0x63, 0x65, 0x49, 0x44, 0x45, 0x78, 0x74, 0x72,
	0x61, 0x63, 0x74, 0x4e, 0x61, 0x6d, 0x65, 0x42, 0x38, 0x5a, 0x36, 0x67, 0x69, 0x74, 0x68, 0x75,
	0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x61, 0x69, 0x6c, 0x65, 0x72, 0x6f, 0x6e, 0x2d, 0x67, 0x61,
	0x74, 0x65, 0x77, 0x61, 0x79, 0x2f, 0x61, 0x69, 0x6c, 0x65, 0x72, 0x6f, 0x6e, 0x2d, 0x67, 0x61,
	0x74, 0x65, 0x77, 0x61, 0x79, 0x2f, 0x61, 0x70, 0x69, 0x73, 0x2f, 0x61, 0x70, 0x70, 0x2f, 0x76,
	0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
})

var (
	file_app_v1_middleware_tracking_proto_rawDescOnce sync.Once
	file_app_v1_middleware_tracking_proto_rawDescData []byte
)

func file_app_v1_middleware_tracking_proto_rawDescGZIP() []byte {
	file_app_v1_middleware_tracking_proto_rawDescOnce.Do(func() {
		file_app_v1_middleware_tracking_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_app_v1_middleware_tracking_proto_rawDesc), len(file_app_v1_middleware_tracking_proto_rawDesc)))
	})
	return file_app_v1_middleware_tracking_proto_rawDescData
}

var file_app_v1_middleware_tracking_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_app_v1_middleware_tracking_proto_goTypes = []any{
	(*TrackingMiddleware)(nil),     // 0: app.v1.TrackingMiddleware
	(*TrackingMiddlewareSpec)(nil), // 1: app.v1.TrackingMiddlewareSpec
	(*kernel.Metadata)(nil),        // 2: kernel.Metadata
	(*kernel.Reference)(nil),       // 3: kernel.Reference
	(kernel.EncodingType)(0),       // 4: kernel.EncodingType
}
var file_app_v1_middleware_tracking_proto_depIdxs = []int32{
	2, // 0: app.v1.TrackingMiddleware.Metadata:type_name -> kernel.Metadata
	1, // 1: app.v1.TrackingMiddleware.Spec:type_name -> app.v1.TrackingMiddlewareSpec
	3, // 2: app.v1.TrackingMiddlewareSpec.ErrorHandler:type_name -> kernel.Reference
	4, // 3: app.v1.TrackingMiddlewareSpec.Encoding:type_name -> kernel.EncodingType
	4, // [4:4] is the sub-list for method output_type
	4, // [4:4] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_app_v1_middleware_tracking_proto_init() }
func file_app_v1_middleware_tracking_proto_init() {
	if File_app_v1_middleware_tracking_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_app_v1_middleware_tracking_proto_rawDesc), len(file_app_v1_middleware_tracking_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_app_v1_middleware_tracking_proto_goTypes,
		DependencyIndexes: file_app_v1_middleware_tracking_proto_depIdxs,
		MessageInfos:      file_app_v1_middleware_tracking_proto_msgTypes,
	}.Build()
	File_app_v1_middleware_tracking_proto = out.File
	file_app_v1_middleware_tracking_proto_goTypes = nil
	file_app_v1_middleware_tracking_proto_depIdxs = nil
}
