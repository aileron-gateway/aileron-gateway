// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.1
// 	protoc        v5.27.2
// source: kernel/logging.proto

package kernel

import (
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

// LogLevel is the defined log output level.
type LogLevel int32

const (
	LogLevel_LogLevelUnknown LogLevel = 0 // Log level Unknown.
	LogLevel_Trace           LogLevel = 1 // Log level Trace.
	LogLevel_Debug           LogLevel = 2 // Log level Debug.
	LogLevel_Info            LogLevel = 3 // Log level Info.
	LogLevel_Warn            LogLevel = 4 // Log level Warn.
	LogLevel_Error           LogLevel = 5 // Log level Error.
	LogLevel_Fatal           LogLevel = 6 // Log level Fatal.
)

// Enum value maps for LogLevel.
var (
	LogLevel_name = map[int32]string{
		0: "LogLevelUnknown",
		1: "Trace",
		2: "Debug",
		3: "Info",
		4: "Warn",
		5: "Error",
		6: "Fatal",
	}
	LogLevel_value = map[string]int32{
		"LogLevelUnknown": 0,
		"Trace":           1,
		"Debug":           2,
		"Info":            3,
		"Warn":            4,
		"Error":           5,
		"Fatal":           6,
	}
)

func (x LogLevel) Enum() *LogLevel {
	p := new(LogLevel)
	*p = x
	return p
}

func (x LogLevel) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (LogLevel) Descriptor() protoreflect.EnumDescriptor {
	return file_kernel_logging_proto_enumTypes[0].Descriptor()
}

func (LogLevel) Type() protoreflect.EnumType {
	return &file_kernel_logging_proto_enumTypes[0]
}

func (x LogLevel) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use LogLevel.Descriptor instead.
func (LogLevel) EnumDescriptor() ([]byte, []int) {
	return file_kernel_logging_proto_rawDescGZIP(), []int{0}
}

var File_kernel_logging_proto protoreflect.FileDescriptor

var file_kernel_logging_proto_rawDesc = []byte{
	0x0a, 0x14, 0x6b, 0x65, 0x72, 0x6e, 0x65, 0x6c, 0x2f, 0x6c, 0x6f, 0x67, 0x67, 0x69, 0x6e, 0x67,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x06, 0x6b, 0x65, 0x72, 0x6e, 0x65, 0x6c, 0x2a, 0x5f,
	0x0a, 0x08, 0x4c, 0x6f, 0x67, 0x4c, 0x65, 0x76, 0x65, 0x6c, 0x12, 0x13, 0x0a, 0x0f, 0x4c, 0x6f,
	0x67, 0x4c, 0x65, 0x76, 0x65, 0x6c, 0x55, 0x6e, 0x6b, 0x6e, 0x6f, 0x77, 0x6e, 0x10, 0x00, 0x12,
	0x09, 0x0a, 0x05, 0x54, 0x72, 0x61, 0x63, 0x65, 0x10, 0x01, 0x12, 0x09, 0x0a, 0x05, 0x44, 0x65,
	0x62, 0x75, 0x67, 0x10, 0x02, 0x12, 0x08, 0x0a, 0x04, 0x49, 0x6e, 0x66, 0x6f, 0x10, 0x03, 0x12,
	0x08, 0x0a, 0x04, 0x57, 0x61, 0x72, 0x6e, 0x10, 0x04, 0x12, 0x09, 0x0a, 0x05, 0x45, 0x72, 0x72,
	0x6f, 0x72, 0x10, 0x05, 0x12, 0x09, 0x0a, 0x05, 0x46, 0x61, 0x74, 0x61, 0x6c, 0x10, 0x06, 0x42,
	0x38, 0x5a, 0x36, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x61, 0x69,
	0x6c, 0x65, 0x72, 0x6f, 0x6e, 0x2d, 0x67, 0x61, 0x74, 0x65, 0x77, 0x61, 0x79, 0x2f, 0x61, 0x69,
	0x6c, 0x65, 0x72, 0x6f, 0x6e, 0x2d, 0x67, 0x61, 0x74, 0x65, 0x77, 0x61, 0x79, 0x2f, 0x61, 0x70,
	0x69, 0x73, 0x2f, 0x6b, 0x65, 0x72, 0x6e, 0x65, 0x6c, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var (
	file_kernel_logging_proto_rawDescOnce sync.Once
	file_kernel_logging_proto_rawDescData = file_kernel_logging_proto_rawDesc
)

func file_kernel_logging_proto_rawDescGZIP() []byte {
	file_kernel_logging_proto_rawDescOnce.Do(func() {
		file_kernel_logging_proto_rawDescData = protoimpl.X.CompressGZIP(file_kernel_logging_proto_rawDescData)
	})
	return file_kernel_logging_proto_rawDescData
}

var file_kernel_logging_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_kernel_logging_proto_goTypes = []any{
	(LogLevel)(0), // 0: kernel.LogLevel
}
var file_kernel_logging_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_kernel_logging_proto_init() }
func file_kernel_logging_proto_init() {
	if File_kernel_logging_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_kernel_logging_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_kernel_logging_proto_goTypes,
		DependencyIndexes: file_kernel_logging_proto_depIdxs,
		EnumInfos:         file_kernel_logging_proto_enumTypes,
	}.Build()
	File_kernel_logging_proto = out.File
	file_kernel_logging_proto_rawDesc = nil
	file_kernel_logging_proto_goTypes = nil
	file_kernel_logging_proto_depIdxs = nil
}
