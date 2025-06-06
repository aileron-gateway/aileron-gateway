// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.4
// 	protoc        v5.29.0
// source: core/v1/httplogger.proto

package v1

import (
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

// HTTPLogger resource definition.
// apiVersion="core/v1", kind="HTTPLogger".
type HTTPLogger struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	APIVersion    string                 `protobuf:"bytes,1,opt,name=APIVersion,json=apiVersion,proto3" json:"APIVersion,omitempty"`
	Kind          string                 `protobuf:"bytes,2,opt,name=Kind,json=kind,proto3" json:"Kind,omitempty"`
	Metadata      *kernel.Metadata       `protobuf:"bytes,3,opt,name=Metadata,json=metadata,proto3" json:"Metadata,omitempty"`
	Spec          *HTTPLoggerSpec        `protobuf:"bytes,4,opt,name=Spec,json=spec,proto3" json:"Spec,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *HTTPLogger) Reset() {
	*x = HTTPLogger{}
	mi := &file_core_v1_httplogger_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *HTTPLogger) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HTTPLogger) ProtoMessage() {}

func (x *HTTPLogger) ProtoReflect() protoreflect.Message {
	mi := &file_core_v1_httplogger_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HTTPLogger.ProtoReflect.Descriptor instead.
func (*HTTPLogger) Descriptor() ([]byte, []int) {
	return file_core_v1_httplogger_proto_rawDescGZIP(), []int{0}
}

func (x *HTTPLogger) GetAPIVersion() string {
	if x != nil {
		return x.APIVersion
	}
	return ""
}

func (x *HTTPLogger) GetKind() string {
	if x != nil {
		return x.Kind
	}
	return ""
}

func (x *HTTPLogger) GetMetadata() *kernel.Metadata {
	if x != nil {
		return x.Metadata
	}
	return nil
}

func (x *HTTPLogger) GetSpec() *HTTPLoggerSpec {
	if x != nil {
		return x.Spec
	}
	return nil
}

// HTTPLoggerSpec is the specifications of the HTTPLogger object.
type HTTPLoggerSpec struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// [OPTIONAL]
	// Logger is the reference to a Logger object.
	// Referred object must implement Logger interface.
	// Default Logger is used if not set.
	Logger *kernel.Reference `protobuf:"bytes,1,opt,name=Logger,json=logger,proto3" json:"Logger,omitempty"`
	// [OPTIONAL]
	// ErrorHandler is the reference to a ErrorHandler object.
	// Referred object must implement ErrorHandler interface.
	// Default error handler is used when not set.
	ErrorHandler *kernel.Reference `protobuf:"bytes,2,opt,name=ErrorHandler,json=errorHandler,proto3" json:"ErrorHandler,omitempty"`
	// [OPTIONAL]
	// Journal is the flag to log request and response bodies.
	// Note that not all bodies are logged by default.
	// Configure the target requests and responses to log
	// in the Request and Response field.
	// Default is [false].
	Journal bool `protobuf:"varint,3,opt,name=Journal,json=journal,proto3" json:"Journal,omitempty"`
	// [OPTIONAL]
	// TimeZone is the timezone of the access time timestamp.
	// For example, "UTC", "Local", "Asia/Tokyo".
	// See https://pkg.go.dev/time#LoadLocation for more details.
	// Default is ["Local"].
	Timezone string `protobuf:"bytes,4,opt,name=Timezone,json=timezone,proto3" json:"Timezone,omitempty"`
	// [OPTIONAL]
	// TimeFormat is the format of the access time timestamp.
	// Check the following url for examples.
	// https://pkg.go.dev/time#pkg-constants
	// Default is ["2006-01-02 15:04:05.000"].
	TimeFormat string `protobuf:"bytes,5,opt,name=TimeFormat,json=timeFormat,proto3" json:"TimeFormat,omitempty"`
	// [OPTIONAL]
	// Request is the logging configuration for requests.
	// Default values are used if not set.
	Request *LoggingSpec `protobuf:"bytes,6,opt,name=Request,json=request,proto3" json:"Request,omitempty"`
	// [OPTIONAL]
	// Response is the logging configuration for responses.
	// Default values are used if not set.
	Response      *LoggingSpec `protobuf:"bytes,7,opt,name=Response,json=response,proto3" json:"Response,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *HTTPLoggerSpec) Reset() {
	*x = HTTPLoggerSpec{}
	mi := &file_core_v1_httplogger_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *HTTPLoggerSpec) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HTTPLoggerSpec) ProtoMessage() {}

func (x *HTTPLoggerSpec) ProtoReflect() protoreflect.Message {
	mi := &file_core_v1_httplogger_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HTTPLoggerSpec.ProtoReflect.Descriptor instead.
func (*HTTPLoggerSpec) Descriptor() ([]byte, []int) {
	return file_core_v1_httplogger_proto_rawDescGZIP(), []int{1}
}

func (x *HTTPLoggerSpec) GetLogger() *kernel.Reference {
	if x != nil {
		return x.Logger
	}
	return nil
}

func (x *HTTPLoggerSpec) GetErrorHandler() *kernel.Reference {
	if x != nil {
		return x.ErrorHandler
	}
	return nil
}

func (x *HTTPLoggerSpec) GetJournal() bool {
	if x != nil {
		return x.Journal
	}
	return false
}

func (x *HTTPLoggerSpec) GetTimezone() string {
	if x != nil {
		return x.Timezone
	}
	return ""
}

func (x *HTTPLoggerSpec) GetTimeFormat() string {
	if x != nil {
		return x.TimeFormat
	}
	return ""
}

func (x *HTTPLoggerSpec) GetRequest() *LoggingSpec {
	if x != nil {
		return x.Request
	}
	return nil
}

func (x *HTTPLoggerSpec) GetResponse() *LoggingSpec {
	if x != nil {
		return x.Response
	}
	return nil
}

type LoggingSpec struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// [OPTIONAL]
	// Headers is the list of headers replacers.
	// Also, masking the header values can be configured in ths field.
	// "*" is the special name that represents all headers.
	// Set "*" to output all header values to log.
	// Note that the masking configuration for "*" is ignored.
	// Nothing is set by default.
	Headers []*LogHeaderSpec `protobuf:"bytes,1,rep,name=Headers,json=headers,proto3" json:"Headers,omitempty"`
	// [OPTIONAL]
	// Bodies is the list of body replacer.
	// Replacers can be used for masking, adding or removing content.
	// This field does not work when logging bodies to files.
	// Default is not set.
	Bodies []*LogBodySpec `protobuf:"bytes,2,rep,name=Bodies,json=bodies,proto3" json:"Bodies,omitempty"`
	// [Optional]
	// Queries is the list of query value replacers.
	// Replacers can be used for masking or removing query values.
	// All replacers are applied to the entire query string
	// with specified order.
	// If not set, query string is output as-is.
	// This field works only for request logs and not for response logs.
	// Default is not set.
	Queries []*kernel.ReplacerSpec `protobuf:"bytes,3,rep,name=Queries,json=queries,proto3" json:"Queries,omitempty"`
	// [OPTIONAL]
	// LogFormat is the customized log output format.
	// if not set, default formats determined by the logger is used.
	// Default is not set.
	LogFormat string `protobuf:"bytes,4,opt,name=LogFormat,json=logFormat,proto3" json:"LogFormat,omitempty"`
	// [OPTIONAL]
	// BodyOutputPath is the body output directory path.
	// If set, bodies that exceeds MaxContentLength and bodies with unknown size
	// are logged to files in the specidied path.
	// This feature of logging bodies to files is debugging use only
	// because of the not enough implementation of error handling.
	// To output all bodies to files, set MaxContentLength a negative value.
	// Note that body replacers do not work for bodies output to files.
	// This field is ignored when Journal mode is not enabled.
	// Default is not set.
	BodyOutputPath string `protobuf:"bytes,5,opt,name=BodyOutputPath,json=bodyOutputPath,proto3" json:"BodyOutputPath,omitempty"`
	// [OPTIONAL]
	// Base64 if enabled, encode body with base64 standard encoding.
	// This field is ignored when output mode is not enabled or BodyOutputPath is set.
	// Default is [false].
	Base64 bool `protobuf:"varint,6,opt,name=Base64,json=base64,proto3" json:"Base64,omitempty"`
	// [OPTIONAL]
	// MaxContentLength is the maximum content length in bytes to
	// allow logging request and response bodies.
	// Request and response bodies which exceeds this length are not logged.
	// Requests and response bodies with unknown sizes are ignored and not logged.
	// Streaming or HTTP2 requests and responses can have unknown sized body.
	// Note that when bodies are logged, the entire body is temporarily load on memory.
	// So do not set this value so large that can result in too much memory consumption.
	// Use BodyOutputPath for logging large or streaming bodies.
	// This field is ignored when Journal mode is not enabled.
	// Default is [4096] or 4kiB.
	MaxContentLength int64 `protobuf:"varint,7,opt,name=MaxContentLength,json=maxContentLength,proto3" json:"MaxContentLength,omitempty"`
	// [OPTIONAL]
	// MIMEs is the list of mime types to log request and response bodies.
	// When Journal mode is enabled, only requests and response bodies with
	// listed mime types are logged.
	// Mime types are evaluated with exact matching.
	// So, list all mime types to log bodies.
	// See MIME types at https://www.iana.org/assignments/media-types/media-types.xhtml.
	// This field is ignored when Journal mode is not enabled.
	// When not set, default value are used.
	// Default values are
	// ["application/json", "application/x-www-form-urlencoded", "application/xml",
	// "application/soap+xml", "application/graphql+json",
	// "text/plain", "text/html", "text/xml"].
	MIMEs         []string `protobuf:"bytes,8,rep,name=MIMEs,json=mimes,proto3" json:"MIMEs,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *LoggingSpec) Reset() {
	*x = LoggingSpec{}
	mi := &file_core_v1_httplogger_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *LoggingSpec) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LoggingSpec) ProtoMessage() {}

func (x *LoggingSpec) ProtoReflect() protoreflect.Message {
	mi := &file_core_v1_httplogger_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LoggingSpec.ProtoReflect.Descriptor instead.
func (*LoggingSpec) Descriptor() ([]byte, []int) {
	return file_core_v1_httplogger_proto_rawDescGZIP(), []int{2}
}

func (x *LoggingSpec) GetHeaders() []*LogHeaderSpec {
	if x != nil {
		return x.Headers
	}
	return nil
}

func (x *LoggingSpec) GetBodies() []*LogBodySpec {
	if x != nil {
		return x.Bodies
	}
	return nil
}

func (x *LoggingSpec) GetQueries() []*kernel.ReplacerSpec {
	if x != nil {
		return x.Queries
	}
	return nil
}

func (x *LoggingSpec) GetLogFormat() string {
	if x != nil {
		return x.LogFormat
	}
	return ""
}

func (x *LoggingSpec) GetBodyOutputPath() string {
	if x != nil {
		return x.BodyOutputPath
	}
	return ""
}

func (x *LoggingSpec) GetBase64() bool {
	if x != nil {
		return x.Base64
	}
	return false
}

func (x *LoggingSpec) GetMaxContentLength() int64 {
	if x != nil {
		return x.MaxContentLength
	}
	return 0
}

func (x *LoggingSpec) GetMIMEs() []string {
	if x != nil {
		return x.MIMEs
	}
	return nil
}

// LogValueSpec is the status of the LoggingMiddleware object.
// Values are managed by the application and therefore should not be set by users.
type LogHeaderSpec struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// [REQUIRED]
	// Name is the HTTP header name.
	// "*" is the special character to represent all header names.
	Name string `protobuf:"bytes,1,opt,name=Name,json=name,proto3" json:"Name,omitempty"`
	// [Optional]
	// Replacers is the list of replace configurations.
	// If not set, header values are output as is.
	// Default is not set.
	Replacers     []*kernel.ReplacerSpec `protobuf:"bytes,2,rep,name=Replacers,json=replacers,proto3" json:"Replacers,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *LogHeaderSpec) Reset() {
	*x = LogHeaderSpec{}
	mi := &file_core_v1_httplogger_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *LogHeaderSpec) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LogHeaderSpec) ProtoMessage() {}

func (x *LogHeaderSpec) ProtoReflect() protoreflect.Message {
	mi := &file_core_v1_httplogger_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LogHeaderSpec.ProtoReflect.Descriptor instead.
func (*LogHeaderSpec) Descriptor() ([]byte, []int) {
	return file_core_v1_httplogger_proto_rawDescGZIP(), []int{3}
}

func (x *LogHeaderSpec) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *LogHeaderSpec) GetReplacers() []*kernel.ReplacerSpec {
	if x != nil {
		return x.Replacers
	}
	return nil
}

// LogBodySpec is the HTTP body logging configuration.
type LogBodySpec struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// [REQUIRED]
	// Mime is the mime type string such as "application/json"
	// that this configuration targets to.
	// Default is not set.
	Mime string `protobuf:"bytes,1,opt,name=Mime,json=mime,proto3" json:"Mime,omitempty"`
	// [Optional]
	// Replacers is the list of replace configurations.
	// If not set, body is output as is.
	// Default is not set.
	Replacers []*kernel.ReplacerSpec `protobuf:"bytes,2,rep,name=Replacers,json=replacers,proto3" json:"Replacers,omitempty"`
	// [OPTIONAL]
	// JSONFields is the list of json key names to be replaced.
	// If set, replacers are applied to the obtained content.
	// Default is not set.
	JSONFields    []string `protobuf:"bytes,3,rep,name=JSONFields,json=jsonFields,proto3" json:"JSONFields,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *LogBodySpec) Reset() {
	*x = LogBodySpec{}
	mi := &file_core_v1_httplogger_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *LogBodySpec) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LogBodySpec) ProtoMessage() {}

func (x *LogBodySpec) ProtoReflect() protoreflect.Message {
	mi := &file_core_v1_httplogger_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LogBodySpec.ProtoReflect.Descriptor instead.
func (*LogBodySpec) Descriptor() ([]byte, []int) {
	return file_core_v1_httplogger_proto_rawDescGZIP(), []int{4}
}

func (x *LogBodySpec) GetMime() string {
	if x != nil {
		return x.Mime
	}
	return ""
}

func (x *LogBodySpec) GetReplacers() []*kernel.ReplacerSpec {
	if x != nil {
		return x.Replacers
	}
	return nil
}

func (x *LogBodySpec) GetJSONFields() []string {
	if x != nil {
		return x.JSONFields
	}
	return nil
}

var File_core_v1_httplogger_proto protoreflect.FileDescriptor

var file_core_v1_httplogger_proto_rawDesc = string([]byte{
	0x0a, 0x18, 0x63, 0x6f, 0x72, 0x65, 0x2f, 0x76, 0x31, 0x2f, 0x68, 0x74, 0x74, 0x70, 0x6c, 0x6f,
	0x67, 0x67, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x07, 0x63, 0x6f, 0x72, 0x65,
	0x2e, 0x76, 0x31, 0x1a, 0x15, 0x6b, 0x65, 0x72, 0x6e, 0x65, 0x6c, 0x2f, 0x72, 0x65, 0x70, 0x6c,
	0x61, 0x63, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x15, 0x6b, 0x65, 0x72, 0x6e,
	0x65, 0x6c, 0x2f, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x22, 0x9b, 0x01, 0x0a, 0x0a, 0x48, 0x54, 0x54, 0x50, 0x4c, 0x6f, 0x67, 0x67, 0x65, 0x72,
	0x12, 0x1e, 0x0a, 0x0a, 0x41, 0x50, 0x49, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x61, 0x70, 0x69, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e,
	0x12, 0x12, 0x0a, 0x04, 0x4b, 0x69, 0x6e, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04,
	0x6b, 0x69, 0x6e, 0x64, 0x12, 0x2c, 0x0a, 0x08, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x6b, 0x65, 0x72, 0x6e, 0x65, 0x6c, 0x2e,
	0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x52, 0x08, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61,
	0x74, 0x61, 0x12, 0x2b, 0x0a, 0x04, 0x53, 0x70, 0x65, 0x63, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x17, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x48, 0x54, 0x54, 0x50, 0x4c,
	0x6f, 0x67, 0x67, 0x65, 0x72, 0x53, 0x70, 0x65, 0x63, 0x52, 0x04, 0x73, 0x70, 0x65, 0x63, 0x22,
	0xaa, 0x02, 0x0a, 0x0e, 0x48, 0x54, 0x54, 0x50, 0x4c, 0x6f, 0x67, 0x67, 0x65, 0x72, 0x53, 0x70,
	0x65, 0x63, 0x12, 0x29, 0x0a, 0x06, 0x4c, 0x6f, 0x67, 0x67, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x11, 0x2e, 0x6b, 0x65, 0x72, 0x6e, 0x65, 0x6c, 0x2e, 0x52, 0x65, 0x66, 0x65,
	0x72, 0x65, 0x6e, 0x63, 0x65, 0x52, 0x06, 0x6c, 0x6f, 0x67, 0x67, 0x65, 0x72, 0x12, 0x35, 0x0a,
	0x0c, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x48, 0x61, 0x6e, 0x64, 0x6c, 0x65, 0x72, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x11, 0x2e, 0x6b, 0x65, 0x72, 0x6e, 0x65, 0x6c, 0x2e, 0x52, 0x65, 0x66,
	0x65, 0x72, 0x65, 0x6e, 0x63, 0x65, 0x52, 0x0c, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x48, 0x61, 0x6e,
	0x64, 0x6c, 0x65, 0x72, 0x12, 0x18, 0x0a, 0x07, 0x4a, 0x6f, 0x75, 0x72, 0x6e, 0x61, 0x6c, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x6a, 0x6f, 0x75, 0x72, 0x6e, 0x61, 0x6c, 0x12, 0x1a,
	0x0a, 0x08, 0x54, 0x69, 0x6d, 0x65, 0x7a, 0x6f, 0x6e, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x08, 0x74, 0x69, 0x6d, 0x65, 0x7a, 0x6f, 0x6e, 0x65, 0x12, 0x1e, 0x0a, 0x0a, 0x54, 0x69,
	0x6d, 0x65, 0x46, 0x6f, 0x72, 0x6d, 0x61, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a,
	0x74, 0x69, 0x6d, 0x65, 0x46, 0x6f, 0x72, 0x6d, 0x61, 0x74, 0x12, 0x2e, 0x0a, 0x07, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x63, 0x6f,
	0x72, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x4c, 0x6f, 0x67, 0x67, 0x69, 0x6e, 0x67, 0x53, 0x70, 0x65,
	0x63, 0x52, 0x07, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x30, 0x0a, 0x08, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x18, 0x07, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x63,
	0x6f, 0x72, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x4c, 0x6f, 0x67, 0x67, 0x69, 0x6e, 0x67, 0x53, 0x70,
	0x65, 0x63, 0x52, 0x08, 0x72, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0xbd, 0x02, 0x0a,
	0x0b, 0x4c, 0x6f, 0x67, 0x67, 0x69, 0x6e, 0x67, 0x53, 0x70, 0x65, 0x63, 0x12, 0x30, 0x0a, 0x07,
	0x48, 0x65, 0x61, 0x64, 0x65, 0x72, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x16, 0x2e,
	0x63, 0x6f, 0x72, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x4c, 0x6f, 0x67, 0x48, 0x65, 0x61, 0x64, 0x65,
	0x72, 0x53, 0x70, 0x65, 0x63, 0x52, 0x07, 0x68, 0x65, 0x61, 0x64, 0x65, 0x72, 0x73, 0x12, 0x2c,
	0x0a, 0x06, 0x42, 0x6f, 0x64, 0x69, 0x65, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x14,
	0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x4c, 0x6f, 0x67, 0x42, 0x6f, 0x64, 0x79,
	0x53, 0x70, 0x65, 0x63, 0x52, 0x06, 0x62, 0x6f, 0x64, 0x69, 0x65, 0x73, 0x12, 0x2e, 0x0a, 0x07,
	0x51, 0x75, 0x65, 0x72, 0x69, 0x65, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x14, 0x2e,
	0x6b, 0x65, 0x72, 0x6e, 0x65, 0x6c, 0x2e, 0x52, 0x65, 0x70, 0x6c, 0x61, 0x63, 0x65, 0x72, 0x53,
	0x70, 0x65, 0x63, 0x52, 0x07, 0x71, 0x75, 0x65, 0x72, 0x69, 0x65, 0x73, 0x12, 0x1c, 0x0a, 0x09,
	0x4c, 0x6f, 0x67, 0x46, 0x6f, 0x72, 0x6d, 0x61, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x09, 0x6c, 0x6f, 0x67, 0x46, 0x6f, 0x72, 0x6d, 0x61, 0x74, 0x12, 0x26, 0x0a, 0x0e, 0x42, 0x6f,
	0x64, 0x79, 0x4f, 0x75, 0x74, 0x70, 0x75, 0x74, 0x50, 0x61, 0x74, 0x68, 0x18, 0x05, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0e, 0x62, 0x6f, 0x64, 0x79, 0x4f, 0x75, 0x74, 0x70, 0x75, 0x74, 0x50, 0x61,
	0x74, 0x68, 0x12, 0x16, 0x0a, 0x06, 0x42, 0x61, 0x73, 0x65, 0x36, 0x34, 0x18, 0x06, 0x20, 0x01,
	0x28, 0x08, 0x52, 0x06, 0x62, 0x61, 0x73, 0x65, 0x36, 0x34, 0x12, 0x2a, 0x0a, 0x10, 0x4d, 0x61,
	0x78, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x4c, 0x65, 0x6e, 0x67, 0x74, 0x68, 0x18, 0x07,
	0x20, 0x01, 0x28, 0x03, 0x52, 0x10, 0x6d, 0x61, 0x78, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74,
	0x4c, 0x65, 0x6e, 0x67, 0x74, 0x68, 0x12, 0x14, 0x0a, 0x05, 0x4d, 0x49, 0x4d, 0x45, 0x73, 0x18,
	0x08, 0x20, 0x03, 0x28, 0x09, 0x52, 0x05, 0x6d, 0x69, 0x6d, 0x65, 0x73, 0x22, 0x57, 0x0a, 0x0d,
	0x4c, 0x6f, 0x67, 0x48, 0x65, 0x61, 0x64, 0x65, 0x72, 0x53, 0x70, 0x65, 0x63, 0x12, 0x12, 0x0a,
	0x04, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d,
	0x65, 0x12, 0x32, 0x0a, 0x09, 0x52, 0x65, 0x70, 0x6c, 0x61, 0x63, 0x65, 0x72, 0x73, 0x18, 0x02,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x6b, 0x65, 0x72, 0x6e, 0x65, 0x6c, 0x2e, 0x52, 0x65,
	0x70, 0x6c, 0x61, 0x63, 0x65, 0x72, 0x53, 0x70, 0x65, 0x63, 0x52, 0x09, 0x72, 0x65, 0x70, 0x6c,
	0x61, 0x63, 0x65, 0x72, 0x73, 0x22, 0x75, 0x0a, 0x0b, 0x4c, 0x6f, 0x67, 0x42, 0x6f, 0x64, 0x79,
	0x53, 0x70, 0x65, 0x63, 0x12, 0x12, 0x0a, 0x04, 0x4d, 0x69, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x04, 0x6d, 0x69, 0x6d, 0x65, 0x12, 0x32, 0x0a, 0x09, 0x52, 0x65, 0x70, 0x6c,
	0x61, 0x63, 0x65, 0x72, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x6b, 0x65,
	0x72, 0x6e, 0x65, 0x6c, 0x2e, 0x52, 0x65, 0x70, 0x6c, 0x61, 0x63, 0x65, 0x72, 0x53, 0x70, 0x65,
	0x63, 0x52, 0x09, 0x72, 0x65, 0x70, 0x6c, 0x61, 0x63, 0x65, 0x72, 0x73, 0x12, 0x1e, 0x0a, 0x0a,
	0x4a, 0x53, 0x4f, 0x4e, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x09,
	0x52, 0x0a, 0x6a, 0x73, 0x6f, 0x6e, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x73, 0x42, 0x39, 0x5a, 0x37,
	0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x61, 0x69, 0x6c, 0x65, 0x72,
	0x6f, 0x6e, 0x2d, 0x67, 0x61, 0x74, 0x65, 0x77, 0x61, 0x79, 0x2f, 0x61, 0x69, 0x6c, 0x65, 0x72,
	0x6f, 0x6e, 0x2d, 0x67, 0x61, 0x74, 0x65, 0x77, 0x61, 0x79, 0x2f, 0x61, 0x70, 0x69, 0x73, 0x2f,
	0x63, 0x6f, 0x72, 0x65, 0x2f, 0x76, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
})

var (
	file_core_v1_httplogger_proto_rawDescOnce sync.Once
	file_core_v1_httplogger_proto_rawDescData []byte
)

func file_core_v1_httplogger_proto_rawDescGZIP() []byte {
	file_core_v1_httplogger_proto_rawDescOnce.Do(func() {
		file_core_v1_httplogger_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_core_v1_httplogger_proto_rawDesc), len(file_core_v1_httplogger_proto_rawDesc)))
	})
	return file_core_v1_httplogger_proto_rawDescData
}

var file_core_v1_httplogger_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_core_v1_httplogger_proto_goTypes = []any{
	(*HTTPLogger)(nil),          // 0: core.v1.HTTPLogger
	(*HTTPLoggerSpec)(nil),      // 1: core.v1.HTTPLoggerSpec
	(*LoggingSpec)(nil),         // 2: core.v1.LoggingSpec
	(*LogHeaderSpec)(nil),       // 3: core.v1.LogHeaderSpec
	(*LogBodySpec)(nil),         // 4: core.v1.LogBodySpec
	(*kernel.Metadata)(nil),     // 5: kernel.Metadata
	(*kernel.Reference)(nil),    // 6: kernel.Reference
	(*kernel.ReplacerSpec)(nil), // 7: kernel.ReplacerSpec
}
var file_core_v1_httplogger_proto_depIdxs = []int32{
	5,  // 0: core.v1.HTTPLogger.Metadata:type_name -> kernel.Metadata
	1,  // 1: core.v1.HTTPLogger.Spec:type_name -> core.v1.HTTPLoggerSpec
	6,  // 2: core.v1.HTTPLoggerSpec.Logger:type_name -> kernel.Reference
	6,  // 3: core.v1.HTTPLoggerSpec.ErrorHandler:type_name -> kernel.Reference
	2,  // 4: core.v1.HTTPLoggerSpec.Request:type_name -> core.v1.LoggingSpec
	2,  // 5: core.v1.HTTPLoggerSpec.Response:type_name -> core.v1.LoggingSpec
	3,  // 6: core.v1.LoggingSpec.Headers:type_name -> core.v1.LogHeaderSpec
	4,  // 7: core.v1.LoggingSpec.Bodies:type_name -> core.v1.LogBodySpec
	7,  // 8: core.v1.LoggingSpec.Queries:type_name -> kernel.ReplacerSpec
	7,  // 9: core.v1.LogHeaderSpec.Replacers:type_name -> kernel.ReplacerSpec
	7,  // 10: core.v1.LogBodySpec.Replacers:type_name -> kernel.ReplacerSpec
	11, // [11:11] is the sub-list for method output_type
	11, // [11:11] is the sub-list for method input_type
	11, // [11:11] is the sub-list for extension type_name
	11, // [11:11] is the sub-list for extension extendee
	0,  // [0:11] is the sub-list for field type_name
}

func init() { file_core_v1_httplogger_proto_init() }
func file_core_v1_httplogger_proto_init() {
	if File_core_v1_httplogger_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_core_v1_httplogger_proto_rawDesc), len(file_core_v1_httplogger_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_core_v1_httplogger_proto_goTypes,
		DependencyIndexes: file_core_v1_httplogger_proto_depIdxs,
		MessageInfos:      file_core_v1_httplogger_proto_msgTypes,
	}.Build()
	File_core_v1_httplogger_proto = out.File
	file_core_v1_httplogger_proto_goTypes = nil
	file_core_v1_httplogger_proto_depIdxs = nil
}
