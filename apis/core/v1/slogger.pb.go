// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.2
// 	protoc        v5.29.0
// source: core/v1/slogger.proto

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

// OutputTarget is the output destination.
type OutputTarget int32

const (
	OutputTarget_Stdout  OutputTarget = 0 // Output to standard output.
	OutputTarget_Stderr  OutputTarget = 1 // Output to standard error output.
	OutputTarget_Discard OutputTarget = 2 // Discard outputs.
	OutputTarget_File    OutputTarget = 3 // Output to files.
)

// Enum value maps for OutputTarget.
var (
	OutputTarget_name = map[int32]string{
		0: "Stdout",
		1: "Stderr",
		2: "Discard",
		3: "File",
	}
	OutputTarget_value = map[string]int32{
		"Stdout":  0,
		"Stderr":  1,
		"Discard": 2,
		"File":    3,
	}
)

func (x OutputTarget) Enum() *OutputTarget {
	p := new(OutputTarget)
	*p = x
	return p
}

func (x OutputTarget) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (OutputTarget) Descriptor() protoreflect.EnumDescriptor {
	return file_core_v1_slogger_proto_enumTypes[0].Descriptor()
}

func (OutputTarget) Type() protoreflect.EnumType {
	return &file_core_v1_slogger_proto_enumTypes[0]
}

func (x OutputTarget) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use OutputTarget.Descriptor instead.
func (OutputTarget) EnumDescriptor() ([]byte, []int) {
	return file_core_v1_slogger_proto_rawDescGZIP(), []int{0}
}

// SLogger is the definition of the SLogger object.
// SLogger implements interface of the logger.
type SLogger struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// [REQUIRED]
	// APIVersion is the defined version of the logger.
	// This value must be "core/v1".
	APIVersion string `protobuf:"bytes,1,opt,name=APIVersion,json=apiVersion,proto3" json:"APIVersion,omitempty"`
	// [REQUIRED]
	// Kind is the kind of this object.
	// This value must be "SLogger".
	Kind string `protobuf:"bytes,2,opt,name=Kind,json=kind,proto3" json:"Kind,omitempty"`
	// [OPTIONAL]
	// Metadata is the metadata of the middleware object.
	// If not set, both name and namespace in the metadata
	// are treated as "default".
	Metadata *kernel.Metadata `protobuf:"bytes,3,opt,name=Metadata,json=metadata,proto3" json:"Metadata,omitempty"`
	// [OPTIONAL]
	// Spec is the specification of the logger.
	// Default values are used when nothing is set.
	Spec          *SLoggerSpec `protobuf:"bytes,4,opt,name=Spec,json=spec,proto3" json:"Spec,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *SLogger) Reset() {
	*x = SLogger{}
	mi := &file_core_v1_slogger_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SLogger) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SLogger) ProtoMessage() {}

func (x *SLogger) ProtoReflect() protoreflect.Message {
	mi := &file_core_v1_slogger_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SLogger.ProtoReflect.Descriptor instead.
func (*SLogger) Descriptor() ([]byte, []int) {
	return file_core_v1_slogger_proto_rawDescGZIP(), []int{0}
}

func (x *SLogger) GetAPIVersion() string {
	if x != nil {
		return x.APIVersion
	}
	return ""
}

func (x *SLogger) GetKind() string {
	if x != nil {
		return x.Kind
	}
	return ""
}

func (x *SLogger) GetMetadata() *kernel.Metadata {
	if x != nil {
		return x.Metadata
	}
	return nil
}

func (x *SLogger) GetSpec() *SLoggerSpec {
	if x != nil {
		return x.Spec
	}
	return nil
}

// SLoggerSpec is the specifications for the SLogger object.
type SLoggerSpec struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// [OPTIONAL]
	// LogLevel is the log output level.
	// Default is [STDOUT].
	Level kernel.LogLevel `protobuf:"varint,1,opt,name=Level,json=level,proto3,enum=kernel.LogLevel" json:"Level,omitempty"`
	// [OPTIONAL]
	// LogOutput is the specifications of log output.
	// Default values are used if not set.
	LogOutput *LogOutputSpec `protobuf:"bytes,2,opt,name=LogOutput,json=logOutput,proto3" json:"LogOutput,omitempty"`
	// [OPTIONAL]
	// Unstructured is the flag to use text log, or non-json log.
	// Default is [false].
	Unstructured bool `protobuf:"varint,3,opt,name=Unstructured,json=unstructured,proto3" json:"Unstructured,omitempty"`
	// [OPTIONAL]
	// OutputTimeFormat is the timestamp format of the
	// "time" field which located on the top level in logs.
	// This time indicates the log 'output' time which is
	// different from log 'creation' time.
	// Use DateFormat and TimeFormat field to specify the
	// timestamp format of log 'creation' time.
	// Timezone is inherited from the LogOutput.TimeZone.
	// Check the following url for the available syntax.
	// https://pkg.go.dev/time#pkg-constants
	// Default is ["2006-01-02 15:04:05"].
	OutputTimeFormat string `protobuf:"bytes,4,opt,name=OutputTimeFormat,json=outputTimeFormat,proto3" json:"OutputTimeFormat,omitempty"`
	// [OPTIONAL]
	// DateFormat is the timestamp format of the date
	// "datetime.date" field in logs.
	// Timestamp of the "datetime" indicates the log 'creation' time
	// which is different from log 'output' time.
	// Timezone is inherited from the LogOutput.TimeZone.
	// Check the following url for the available syntax.
	// https://pkg.go.dev/time#pkg-constants
	// Default is ["2006-01-02"].
	DateFormat string `protobuf:"bytes,5,opt,name=DateFormat,json=dateFormat,proto3" json:"DateFormat,omitempty"`
	// [OPTIONAL]
	// TimeFormat is the timestamp format of the time
	// "datetime.time" field in logs.
	// Timestamp of the "datetime" indicates the log 'creation' time
	// which is different from log 'output' time.
	// Timezone is inherited from the LogOutput.TimeZone.
	// Check the following url for the available syntax.
	// https://pkg.go.dev/time#pkg-constants
	// Default is ["15:04:05.000"].
	TimeFormat string `protobuf:"bytes,6,opt,name=TimeFormat,json=timeFormat,proto3" json:"TimeFormat,omitempty"`
	// [OPTIONAL]
	// NoLocation is the flag to suppress "location" log field.
	// It contains file name, line number and function name.
	// Default is [false].
	NoLocation bool `protobuf:"varint,7,opt,name=NoLocation,json=noLocation,proto3" json:"NoLocation,omitempty"`
	// [OPTIONAL]
	// NoDatetime is the flag to suppress "datetime" log field.
	// It contains date, time and time zone.
	// Default is [false].
	NoDatetime bool `protobuf:"varint,8,opt,name=NoDatetime,json=noDatetime,proto3" json:"NoDatetime,omitempty"`
	// [OPTIONAL]
	// FieldReplaces is the list of field replace configuration.
	// This can be used for log masking.
	// Default is not set.
	FieldReplacers []*FieldReplacerSpec `protobuf:"bytes,9,rep,name=FieldReplacers,json=fieldReplacers,proto3" json:"FieldReplacers,omitempty"`
	unknownFields  protoimpl.UnknownFields
	sizeCache      protoimpl.SizeCache
}

func (x *SLoggerSpec) Reset() {
	*x = SLoggerSpec{}
	mi := &file_core_v1_slogger_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SLoggerSpec) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SLoggerSpec) ProtoMessage() {}

func (x *SLoggerSpec) ProtoReflect() protoreflect.Message {
	mi := &file_core_v1_slogger_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SLoggerSpec.ProtoReflect.Descriptor instead.
func (*SLoggerSpec) Descriptor() ([]byte, []int) {
	return file_core_v1_slogger_proto_rawDescGZIP(), []int{1}
}

func (x *SLoggerSpec) GetLevel() kernel.LogLevel {
	if x != nil {
		return x.Level
	}
	return kernel.LogLevel(0)
}

func (x *SLoggerSpec) GetLogOutput() *LogOutputSpec {
	if x != nil {
		return x.LogOutput
	}
	return nil
}

func (x *SLoggerSpec) GetUnstructured() bool {
	if x != nil {
		return x.Unstructured
	}
	return false
}

func (x *SLoggerSpec) GetOutputTimeFormat() string {
	if x != nil {
		return x.OutputTimeFormat
	}
	return ""
}

func (x *SLoggerSpec) GetDateFormat() string {
	if x != nil {
		return x.DateFormat
	}
	return ""
}

func (x *SLoggerSpec) GetTimeFormat() string {
	if x != nil {
		return x.TimeFormat
	}
	return ""
}

func (x *SLoggerSpec) GetNoLocation() bool {
	if x != nil {
		return x.NoLocation
	}
	return false
}

func (x *SLoggerSpec) GetNoDatetime() bool {
	if x != nil {
		return x.NoDatetime
	}
	return false
}

func (x *SLoggerSpec) GetFieldReplacers() []*FieldReplacerSpec {
	if x != nil {
		return x.FieldReplacers
	}
	return nil
}

type FieldReplacerSpec struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// [REQUIRED]
	// Field is the field name to replace value.
	// Inner fields, if having map structure,
	// can be accessed by paath expression
	// line "foo.bar.baz".
	// Default is not set.
	Field string `protobuf:"bytes,1,opt,name=Field,json=field,proto3" json:"Field,omitempty"`
	// [Optional]
	// Replacers is the value replace configurations.
	// If not set, target field is removed from log outpur.
	// Default is not set, which means remove field.
	Replacer      *kernel.ReplacerSpec `protobuf:"bytes,2,opt,name=Replacer,json=replacer,proto3" json:"Replacer,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *FieldReplacerSpec) Reset() {
	*x = FieldReplacerSpec{}
	mi := &file_core_v1_slogger_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *FieldReplacerSpec) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FieldReplacerSpec) ProtoMessage() {}

func (x *FieldReplacerSpec) ProtoReflect() protoreflect.Message {
	mi := &file_core_v1_slogger_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FieldReplacerSpec.ProtoReflect.Descriptor instead.
func (*FieldReplacerSpec) Descriptor() ([]byte, []int) {
	return file_core_v1_slogger_proto_rawDescGZIP(), []int{2}
}

func (x *FieldReplacerSpec) GetField() string {
	if x != nil {
		return x.Field
	}
	return ""
}

func (x *FieldReplacerSpec) GetReplacer() *kernel.ReplacerSpec {
	if x != nil {
		return x.Replacer
	}
	return nil
}

// LogOutputSpec is the specification for log output.
type LogOutputSpec struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// [OPTIONAL]
	// OutputTarget is the destination of log outsput.
	// Default is [Stdout].
	OutputTarget OutputTarget `protobuf:"varint,1,opt,name=OutputTarget,json=outputTarget,proto3,enum=core.v1.OutputTarget" json:"OutputTarget,omitempty"`
	// [OPTIONAL]
	// LogDir is the log output directory path.
	// This field is used only for "File" output.
	// Default is the working directory.
	LogDir string `protobuf:"bytes,2,opt,name=LogDir,json=logDir,proto3" json:"LogDir,omitempty"`
	// [OPTIONAL]
	// LogDir is the log output directory path.
	// This field is used only for "File" output.
	// Default is the same as LogDir.
	BackupDir string `protobuf:"bytes,3,opt,name=BackupDir,json=backupDir,proto3" json:"BackupDir,omitempty"`
	// [OPTIONAL]
	// LogFileName is the base filename of logs.
	// This field is used only for "File" output.
	// Default is ["application.log"].
	LogFileName string `protobuf:"bytes,4,opt,name=LogFileName,json=logFileName,proto3" json:"LogFileName,omitempty"`
	// [OPTIONAL]
	// Cron is the cron expression for time based log rotation.
	// If not set, time based rotation will be disabled.
	// Format should be "second minute hour day month week"
	// or "minute hour day month week".
	// TZ must be a valid timezone name.
	// Value ranges are `0-59` for second, `0-59` for minute,
	// `0-23` for hour, `1-31` for day of month,
	// `1-12 or JAN-DEC` for month, `0-6 or SUN-SAT` for day of week.
	// Special caharacters `* / , -` are allowed for all fields.
	// Timezone can be specified like "TZ=UTC * * * * *".
	// For example, "0 * * * *" means hourly rotation,
	// "0 0 * * *" means daily rotation.
	// Multiple jobs won't be run at the same time.
	// Default is not set.
	Cron string `protobuf:"bytes,5,opt,name=Cron,json=cron,proto3" json:"Cron,omitempty"`
	// [OPTIONAL]
	// RotateSize is the log file size in MiB to be rotated.
	// This field is used only for "File" output.
	// Default is [1024], or 1GiB.
	RotateSize uint32 `protobuf:"varint,6,opt,name=RotateSize,json=rotateSize,proto3" json:"RotateSize,omitempty"`
	// [OPTIONAL]
	// TimeLayout is the timestamp layout in the backup-ed log files.
	// This field is used only for "File" output.
	// Go's time format can be used.
	// See https://pkg.go.dev/time#pkg-constants for more details.
	// Default is ["2006-01-02_15-04-05"].
	TimeLayout string `protobuf:"bytes,7,opt,name=TimeLayout,json=timeLayout,proto3" json:"TimeLayout,omitempty"`
	// [OPTIONAL]
	// TimeZone is the timezone of the timestamp in the archived log files.
	// For example, "UTC", "Local", "Asia/Tokyo".
	// See https://pkg.go.dev/time#LoadLocation for more details.
	// Default is ["Local"].
	TimeZone string `protobuf:"bytes,8,opt,name=TimeZone,json=timeZone,proto3" json:"TimeZone,omitempty"`
	// [OPTIONAL]
	// CompressLevel is the gzip compression level.
	// If 0, no compression applied.
	// This field is ignored when the output target is not file.
	// 0 for no compression, 1 for best speed, 9 for best compression,
	// -2 for huffman only.
	// See https://pkg.go.dev/compress/gzip#pkg-constants for more detail.
	// This field is used only for "File" output.
	// Default is [0], or no compression.
	CompressLevel int32 `protobuf:"varint,9,opt,name=CompressLevel,json=compressLevel,proto3" json:"CompressLevel,omitempty"`
	// [OPTIONAL]
	// MaxAge is maximum age of backup logs in seconds.
	// Backups older than this age will be removed.
	// This field is used only for "File" output.
	// This field will be ignored when not set or set to negative.
	// Default is not set.
	MaxAge int32 `protobuf:"varint,10,opt,name=MaxAge,json=maxAge,proto3" json:"MaxAge,omitempty"`
	// [OPTIONAL]
	// MaxBackup is the maximum number of log backups.
	// Backups will be removed from older one if the number of backups exceeded this value.
	// This field is used only for "File" output.
	// This field will be ignored when not set or set to negative.
	// Default is not set.
	MaxBackup int32 `protobuf:"varint,11,opt,name=MaxBackup,json=maxBackup,proto3" json:"MaxBackup,omitempty"`
	// [OPTIONAL]
	// MaxTotalSize is the maximum total size of backups in MiB.
	// Backups will be removed from older one if the total file size
	// exceeded this value.
	// This field is used only for "File" output.
	// This field will be ignored when not set.
	// Default is not set.
	MaxTotalSize  uint32 `protobuf:"varint,12,opt,name=MaxTotalSize,json=maxTotalSize,proto3" json:"MaxTotalSize,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *LogOutputSpec) Reset() {
	*x = LogOutputSpec{}
	mi := &file_core_v1_slogger_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *LogOutputSpec) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LogOutputSpec) ProtoMessage() {}

func (x *LogOutputSpec) ProtoReflect() protoreflect.Message {
	mi := &file_core_v1_slogger_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LogOutputSpec.ProtoReflect.Descriptor instead.
func (*LogOutputSpec) Descriptor() ([]byte, []int) {
	return file_core_v1_slogger_proto_rawDescGZIP(), []int{3}
}

func (x *LogOutputSpec) GetOutputTarget() OutputTarget {
	if x != nil {
		return x.OutputTarget
	}
	return OutputTarget_Stdout
}

func (x *LogOutputSpec) GetLogDir() string {
	if x != nil {
		return x.LogDir
	}
	return ""
}

func (x *LogOutputSpec) GetBackupDir() string {
	if x != nil {
		return x.BackupDir
	}
	return ""
}

func (x *LogOutputSpec) GetLogFileName() string {
	if x != nil {
		return x.LogFileName
	}
	return ""
}

func (x *LogOutputSpec) GetCron() string {
	if x != nil {
		return x.Cron
	}
	return ""
}

func (x *LogOutputSpec) GetRotateSize() uint32 {
	if x != nil {
		return x.RotateSize
	}
	return 0
}

func (x *LogOutputSpec) GetTimeLayout() string {
	if x != nil {
		return x.TimeLayout
	}
	return ""
}

func (x *LogOutputSpec) GetTimeZone() string {
	if x != nil {
		return x.TimeZone
	}
	return ""
}

func (x *LogOutputSpec) GetCompressLevel() int32 {
	if x != nil {
		return x.CompressLevel
	}
	return 0
}

func (x *LogOutputSpec) GetMaxAge() int32 {
	if x != nil {
		return x.MaxAge
	}
	return 0
}

func (x *LogOutputSpec) GetMaxBackup() int32 {
	if x != nil {
		return x.MaxBackup
	}
	return 0
}

func (x *LogOutputSpec) GetMaxTotalSize() uint32 {
	if x != nil {
		return x.MaxTotalSize
	}
	return 0
}

var File_core_v1_slogger_proto protoreflect.FileDescriptor

var file_core_v1_slogger_proto_rawDesc = []byte{
	0x0a, 0x15, 0x63, 0x6f, 0x72, 0x65, 0x2f, 0x76, 0x31, 0x2f, 0x73, 0x6c, 0x6f, 0x67, 0x67, 0x65,
	0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x07, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x76, 0x31,
	0x1a, 0x1b, 0x62, 0x75, 0x66, 0x2f, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2f, 0x76,
	0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x14, 0x6b,
	0x65, 0x72, 0x6e, 0x65, 0x6c, 0x2f, 0x6c, 0x6f, 0x67, 0x67, 0x69, 0x6e, 0x67, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x1a, 0x15, 0x6b, 0x65, 0x72, 0x6e, 0x65, 0x6c, 0x2f, 0x72, 0x65, 0x73, 0x6f,
	0x75, 0x72, 0x63, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x14, 0x6b, 0x65, 0x72, 0x6e,
	0x65, 0x6c, 0x2f, 0x74, 0x78, 0x74, 0x75, 0x74, 0x69, 0x6c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x22, 0xb5, 0x01, 0x0a, 0x07, 0x53, 0x4c, 0x6f, 0x67, 0x67, 0x65, 0x72, 0x12, 0x2e, 0x0a, 0x0a,
	0x41, 0x50, 0x49, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x42, 0x0e, 0xba, 0x48, 0x0b, 0x72, 0x09, 0x0a, 0x07, 0x63, 0x6f, 0x72, 0x65, 0x2f, 0x76, 0x31,
	0x52, 0x0a, 0x61, 0x70, 0x69, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x22, 0x0a, 0x04,
	0x4b, 0x69, 0x6e, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x42, 0x0e, 0xba, 0x48, 0x0b, 0x72,
	0x09, 0x0a, 0x07, 0x53, 0x4c, 0x6f, 0x67, 0x67, 0x65, 0x72, 0x52, 0x04, 0x6b, 0x69, 0x6e, 0x64,
	0x12, 0x2c, 0x0a, 0x08, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x10, 0x2e, 0x6b, 0x65, 0x72, 0x6e, 0x65, 0x6c, 0x2e, 0x4d, 0x65, 0x74, 0x61,
	0x64, 0x61, 0x74, 0x61, 0x52, 0x08, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x12, 0x28,
	0x0a, 0x04, 0x53, 0x70, 0x65, 0x63, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x63,
	0x6f, 0x72, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x4c, 0x6f, 0x67, 0x67, 0x65, 0x72, 0x53, 0x70,
	0x65, 0x63, 0x52, 0x04, 0x73, 0x70, 0x65, 0x63, 0x22, 0xff, 0x02, 0x0a, 0x0b, 0x53, 0x4c, 0x6f,
	0x67, 0x67, 0x65, 0x72, 0x53, 0x70, 0x65, 0x63, 0x12, 0x26, 0x0a, 0x05, 0x4c, 0x65, 0x76, 0x65,
	0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x10, 0x2e, 0x6b, 0x65, 0x72, 0x6e, 0x65, 0x6c,
	0x2e, 0x4c, 0x6f, 0x67, 0x4c, 0x65, 0x76, 0x65, 0x6c, 0x52, 0x05, 0x6c, 0x65, 0x76, 0x65, 0x6c,
	0x12, 0x34, 0x0a, 0x09, 0x4c, 0x6f, 0x67, 0x4f, 0x75, 0x74, 0x70, 0x75, 0x74, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x16, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x4c, 0x6f,
	0x67, 0x4f, 0x75, 0x74, 0x70, 0x75, 0x74, 0x53, 0x70, 0x65, 0x63, 0x52, 0x09, 0x6c, 0x6f, 0x67,
	0x4f, 0x75, 0x74, 0x70, 0x75, 0x74, 0x12, 0x22, 0x0a, 0x0c, 0x55, 0x6e, 0x73, 0x74, 0x72, 0x75,
	0x63, 0x74, 0x75, 0x72, 0x65, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0c, 0x75, 0x6e,
	0x73, 0x74, 0x72, 0x75, 0x63, 0x74, 0x75, 0x72, 0x65, 0x64, 0x12, 0x2a, 0x0a, 0x10, 0x4f, 0x75,
	0x74, 0x70, 0x75, 0x74, 0x54, 0x69, 0x6d, 0x65, 0x46, 0x6f, 0x72, 0x6d, 0x61, 0x74, 0x18, 0x04,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x10, 0x6f, 0x75, 0x74, 0x70, 0x75, 0x74, 0x54, 0x69, 0x6d, 0x65,
	0x46, 0x6f, 0x72, 0x6d, 0x61, 0x74, 0x12, 0x1e, 0x0a, 0x0a, 0x44, 0x61, 0x74, 0x65, 0x46, 0x6f,
	0x72, 0x6d, 0x61, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x64, 0x61, 0x74, 0x65,
	0x46, 0x6f, 0x72, 0x6d, 0x61, 0x74, 0x12, 0x1e, 0x0a, 0x0a, 0x54, 0x69, 0x6d, 0x65, 0x46, 0x6f,
	0x72, 0x6d, 0x61, 0x74, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x74, 0x69, 0x6d, 0x65,
	0x46, 0x6f, 0x72, 0x6d, 0x61, 0x74, 0x12, 0x1e, 0x0a, 0x0a, 0x4e, 0x6f, 0x4c, 0x6f, 0x63, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x18, 0x07, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0a, 0x6e, 0x6f, 0x4c, 0x6f,
	0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x1e, 0x0a, 0x0a, 0x4e, 0x6f, 0x44, 0x61, 0x74, 0x65,
	0x74, 0x69, 0x6d, 0x65, 0x18, 0x08, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0a, 0x6e, 0x6f, 0x44, 0x61,
	0x74, 0x65, 0x74, 0x69, 0x6d, 0x65, 0x12, 0x42, 0x0a, 0x0e, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x52,
	0x65, 0x70, 0x6c, 0x61, 0x63, 0x65, 0x72, 0x73, 0x18, 0x09, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1a,
	0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x52, 0x65,
	0x70, 0x6c, 0x61, 0x63, 0x65, 0x72, 0x53, 0x70, 0x65, 0x63, 0x52, 0x0e, 0x66, 0x69, 0x65, 0x6c,
	0x64, 0x52, 0x65, 0x70, 0x6c, 0x61, 0x63, 0x65, 0x72, 0x73, 0x22, 0x5b, 0x0a, 0x11, 0x46, 0x69,
	0x65, 0x6c, 0x64, 0x52, 0x65, 0x70, 0x6c, 0x61, 0x63, 0x65, 0x72, 0x53, 0x70, 0x65, 0x63, 0x12,
	0x14, 0x0a, 0x05, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05,
	0x66, 0x69, 0x65, 0x6c, 0x64, 0x12, 0x30, 0x0a, 0x08, 0x52, 0x65, 0x70, 0x6c, 0x61, 0x63, 0x65,
	0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x6b, 0x65, 0x72, 0x6e, 0x65, 0x6c,
	0x2e, 0x52, 0x65, 0x70, 0x6c, 0x61, 0x63, 0x65, 0x72, 0x53, 0x70, 0x65, 0x63, 0x52, 0x08, 0x72,
	0x65, 0x70, 0x6c, 0x61, 0x63, 0x65, 0x72, 0x22, 0xa6, 0x03, 0x0a, 0x0d, 0x4c, 0x6f, 0x67, 0x4f,
	0x75, 0x74, 0x70, 0x75, 0x74, 0x53, 0x70, 0x65, 0x63, 0x12, 0x39, 0x0a, 0x0c, 0x4f, 0x75, 0x74,
	0x70, 0x75, 0x74, 0x54, 0x61, 0x72, 0x67, 0x65, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32,
	0x15, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x4f, 0x75, 0x74, 0x70, 0x75, 0x74,
	0x54, 0x61, 0x72, 0x67, 0x65, 0x74, 0x52, 0x0c, 0x6f, 0x75, 0x74, 0x70, 0x75, 0x74, 0x54, 0x61,
	0x72, 0x67, 0x65, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x4c, 0x6f, 0x67, 0x44, 0x69, 0x72, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x6c, 0x6f, 0x67, 0x44, 0x69, 0x72, 0x12, 0x1c, 0x0a, 0x09,
	0x42, 0x61, 0x63, 0x6b, 0x75, 0x70, 0x44, 0x69, 0x72, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x09, 0x62, 0x61, 0x63, 0x6b, 0x75, 0x70, 0x44, 0x69, 0x72, 0x12, 0x20, 0x0a, 0x0b, 0x4c, 0x6f,
	0x67, 0x46, 0x69, 0x6c, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x0b, 0x6c, 0x6f, 0x67, 0x46, 0x69, 0x6c, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x12, 0x0a, 0x04,
	0x43, 0x72, 0x6f, 0x6e, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x63, 0x72, 0x6f, 0x6e,
	0x12, 0x1e, 0x0a, 0x0a, 0x52, 0x6f, 0x74, 0x61, 0x74, 0x65, 0x53, 0x69, 0x7a, 0x65, 0x18, 0x06,
	0x20, 0x01, 0x28, 0x0d, 0x52, 0x0a, 0x72, 0x6f, 0x74, 0x61, 0x74, 0x65, 0x53, 0x69, 0x7a, 0x65,
	0x12, 0x1e, 0x0a, 0x0a, 0x54, 0x69, 0x6d, 0x65, 0x4c, 0x61, 0x79, 0x6f, 0x75, 0x74, 0x18, 0x07,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x74, 0x69, 0x6d, 0x65, 0x4c, 0x61, 0x79, 0x6f, 0x75, 0x74,
	0x12, 0x1a, 0x0a, 0x08, 0x54, 0x69, 0x6d, 0x65, 0x5a, 0x6f, 0x6e, 0x65, 0x18, 0x08, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x08, 0x74, 0x69, 0x6d, 0x65, 0x5a, 0x6f, 0x6e, 0x65, 0x12, 0x38, 0x0a, 0x0d,
	0x43, 0x6f, 0x6d, 0x70, 0x72, 0x65, 0x73, 0x73, 0x4c, 0x65, 0x76, 0x65, 0x6c, 0x18, 0x09, 0x20,
	0x01, 0x28, 0x05, 0x42, 0x12, 0xba, 0x48, 0x0f, 0x1a, 0x0d, 0x18, 0x09, 0x28, 0xfe, 0xff, 0xff,
	0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x01, 0x52, 0x0d, 0x63, 0x6f, 0x6d, 0x70, 0x72, 0x65, 0x73,
	0x73, 0x4c, 0x65, 0x76, 0x65, 0x6c, 0x12, 0x16, 0x0a, 0x06, 0x4d, 0x61, 0x78, 0x41, 0x67, 0x65,
	0x18, 0x0a, 0x20, 0x01, 0x28, 0x05, 0x52, 0x06, 0x6d, 0x61, 0x78, 0x41, 0x67, 0x65, 0x12, 0x1c,
	0x0a, 0x09, 0x4d, 0x61, 0x78, 0x42, 0x61, 0x63, 0x6b, 0x75, 0x70, 0x18, 0x0b, 0x20, 0x01, 0x28,
	0x05, 0x52, 0x09, 0x6d, 0x61, 0x78, 0x42, 0x61, 0x63, 0x6b, 0x75, 0x70, 0x12, 0x22, 0x0a, 0x0c,
	0x4d, 0x61, 0x78, 0x54, 0x6f, 0x74, 0x61, 0x6c, 0x53, 0x69, 0x7a, 0x65, 0x18, 0x0c, 0x20, 0x01,
	0x28, 0x0d, 0x52, 0x0c, 0x6d, 0x61, 0x78, 0x54, 0x6f, 0x74, 0x61, 0x6c, 0x53, 0x69, 0x7a, 0x65,
	0x2a, 0x3d, 0x0a, 0x0c, 0x4f, 0x75, 0x74, 0x70, 0x75, 0x74, 0x54, 0x61, 0x72, 0x67, 0x65, 0x74,
	0x12, 0x0a, 0x0a, 0x06, 0x53, 0x74, 0x64, 0x6f, 0x75, 0x74, 0x10, 0x00, 0x12, 0x0a, 0x0a, 0x06,
	0x53, 0x74, 0x64, 0x65, 0x72, 0x72, 0x10, 0x01, 0x12, 0x0b, 0x0a, 0x07, 0x44, 0x69, 0x73, 0x63,
	0x61, 0x72, 0x64, 0x10, 0x02, 0x12, 0x08, 0x0a, 0x04, 0x46, 0x69, 0x6c, 0x65, 0x10, 0x03, 0x42,
	0x39, 0x5a, 0x37, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x61, 0x69,
	0x6c, 0x65, 0x72, 0x6f, 0x6e, 0x2d, 0x67, 0x61, 0x74, 0x65, 0x77, 0x61, 0x79, 0x2f, 0x61, 0x69,
	0x6c, 0x65, 0x72, 0x6f, 0x6e, 0x2d, 0x67, 0x61, 0x74, 0x65, 0x77, 0x61, 0x79, 0x2f, 0x61, 0x70,
	0x69, 0x73, 0x2f, 0x63, 0x6f, 0x72, 0x65, 0x2f, 0x76, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_core_v1_slogger_proto_rawDescOnce sync.Once
	file_core_v1_slogger_proto_rawDescData = file_core_v1_slogger_proto_rawDesc
)

func file_core_v1_slogger_proto_rawDescGZIP() []byte {
	file_core_v1_slogger_proto_rawDescOnce.Do(func() {
		file_core_v1_slogger_proto_rawDescData = protoimpl.X.CompressGZIP(file_core_v1_slogger_proto_rawDescData)
	})
	return file_core_v1_slogger_proto_rawDescData
}

var file_core_v1_slogger_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_core_v1_slogger_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_core_v1_slogger_proto_goTypes = []any{
	(OutputTarget)(0),           // 0: core.v1.OutputTarget
	(*SLogger)(nil),             // 1: core.v1.SLogger
	(*SLoggerSpec)(nil),         // 2: core.v1.SLoggerSpec
	(*FieldReplacerSpec)(nil),   // 3: core.v1.FieldReplacerSpec
	(*LogOutputSpec)(nil),       // 4: core.v1.LogOutputSpec
	(*kernel.Metadata)(nil),     // 5: kernel.Metadata
	(kernel.LogLevel)(0),        // 6: kernel.LogLevel
	(*kernel.ReplacerSpec)(nil), // 7: kernel.ReplacerSpec
}
var file_core_v1_slogger_proto_depIdxs = []int32{
	5, // 0: core.v1.SLogger.Metadata:type_name -> kernel.Metadata
	2, // 1: core.v1.SLogger.Spec:type_name -> core.v1.SLoggerSpec
	6, // 2: core.v1.SLoggerSpec.Level:type_name -> kernel.LogLevel
	4, // 3: core.v1.SLoggerSpec.LogOutput:type_name -> core.v1.LogOutputSpec
	3, // 4: core.v1.SLoggerSpec.FieldReplacers:type_name -> core.v1.FieldReplacerSpec
	7, // 5: core.v1.FieldReplacerSpec.Replacer:type_name -> kernel.ReplacerSpec
	0, // 6: core.v1.LogOutputSpec.OutputTarget:type_name -> core.v1.OutputTarget
	7, // [7:7] is the sub-list for method output_type
	7, // [7:7] is the sub-list for method input_type
	7, // [7:7] is the sub-list for extension type_name
	7, // [7:7] is the sub-list for extension extendee
	0, // [0:7] is the sub-list for field type_name
}

func init() { file_core_v1_slogger_proto_init() }
func file_core_v1_slogger_proto_init() {
	if File_core_v1_slogger_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_core_v1_slogger_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_core_v1_slogger_proto_goTypes,
		DependencyIndexes: file_core_v1_slogger_proto_depIdxs,
		EnumInfos:         file_core_v1_slogger_proto_enumTypes,
		MessageInfos:      file_core_v1_slogger_proto_msgTypes,
	}.Build()
	File_core_v1_slogger_proto = out.File
	file_core_v1_slogger_proto_rawDesc = nil
	file_core_v1_slogger_proto_goTypes = nil
	file_core_v1_slogger_proto_depIdxs = nil
}
