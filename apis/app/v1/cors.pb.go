// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.5
// 	protoc        v5.29.0
// source: app/v1/middleware/cors.proto

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

// CORSEmbedderPolicy is the collection of cross origin embedder policy types.
// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Cross-Origin-Embedder-Policy
type CORSEmbedderPolicy int32

const (
	CORSEmbedderPolicy_EmbedderPolicy_Unknown CORSEmbedderPolicy = 0 // ""
	CORSEmbedderPolicy_EmbedderUnsafeNone     CORSEmbedderPolicy = 1 // "unsafe-none"
	CORSEmbedderPolicy_EmbedderRequireCorp    CORSEmbedderPolicy = 2 // "require-corp"
	CORSEmbedderPolicy_EmbedderCredentialless CORSEmbedderPolicy = 3 // "credentialless"
)

// Enum value maps for CORSEmbedderPolicy.
var (
	CORSEmbedderPolicy_name = map[int32]string{
		0: "EmbedderPolicy_Unknown",
		1: "EmbedderUnsafeNone",
		2: "EmbedderRequireCorp",
		3: "EmbedderCredentialless",
	}
	CORSEmbedderPolicy_value = map[string]int32{
		"EmbedderPolicy_Unknown": 0,
		"EmbedderUnsafeNone":     1,
		"EmbedderRequireCorp":    2,
		"EmbedderCredentialless": 3,
	}
)

func (x CORSEmbedderPolicy) Enum() *CORSEmbedderPolicy {
	p := new(CORSEmbedderPolicy)
	*p = x
	return p
}

func (x CORSEmbedderPolicy) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (CORSEmbedderPolicy) Descriptor() protoreflect.EnumDescriptor {
	return file_app_v1_middleware_cors_proto_enumTypes[0].Descriptor()
}

func (CORSEmbedderPolicy) Type() protoreflect.EnumType {
	return &file_app_v1_middleware_cors_proto_enumTypes[0]
}

func (x CORSEmbedderPolicy) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use CORSEmbedderPolicy.Descriptor instead.
func (CORSEmbedderPolicy) EnumDescriptor() ([]byte, []int) {
	return file_app_v1_middleware_cors_proto_rawDescGZIP(), []int{0}
}

// CORSOpenerPolicy is the collection of cross origin opener policy types.
// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Cross-Origin-Opener-Policy
type CORSOpenerPolicy int32

const (
	CORSOpenerPolicy_OpenerPolicy_Unknown        CORSOpenerPolicy = 0 // ""
	CORSOpenerPolicy_OpenerUnsafeNone            CORSOpenerPolicy = 1 // "unsafe-none"
	CORSOpenerPolicy_OpenerSameOriginAllowPopups CORSOpenerPolicy = 2 // "same-origin-allow-popups"
	CORSOpenerPolicy_OpenerSameOrigin            CORSOpenerPolicy = 3 // "same-origin"
)

// Enum value maps for CORSOpenerPolicy.
var (
	CORSOpenerPolicy_name = map[int32]string{
		0: "OpenerPolicy_Unknown",
		1: "OpenerUnsafeNone",
		2: "OpenerSameOriginAllowPopups",
		3: "OpenerSameOrigin",
	}
	CORSOpenerPolicy_value = map[string]int32{
		"OpenerPolicy_Unknown":        0,
		"OpenerUnsafeNone":            1,
		"OpenerSameOriginAllowPopups": 2,
		"OpenerSameOrigin":            3,
	}
)

func (x CORSOpenerPolicy) Enum() *CORSOpenerPolicy {
	p := new(CORSOpenerPolicy)
	*p = x
	return p
}

func (x CORSOpenerPolicy) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (CORSOpenerPolicy) Descriptor() protoreflect.EnumDescriptor {
	return file_app_v1_middleware_cors_proto_enumTypes[1].Descriptor()
}

func (CORSOpenerPolicy) Type() protoreflect.EnumType {
	return &file_app_v1_middleware_cors_proto_enumTypes[1]
}

func (x CORSOpenerPolicy) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use CORSOpenerPolicy.Descriptor instead.
func (CORSOpenerPolicy) EnumDescriptor() ([]byte, []int) {
	return file_app_v1_middleware_cors_proto_rawDescGZIP(), []int{1}
}

// CORSResourcePolicy is the collection of cross origin resource policy types.
// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Cross-Origin-Resource-Policy
type CORSResourcePolicy int32

const (
	CORSResourcePolicy_ResourcePolicy_Unknown CORSResourcePolicy = 0 // ""
	CORSResourcePolicy_ResourceSameSite       CORSResourcePolicy = 1 // "same-site"
	CORSResourcePolicy_ResourceSameOrigin     CORSResourcePolicy = 2 // "same-origin"
	CORSResourcePolicy_ResourceCrossOrigin    CORSResourcePolicy = 3 // "cross-origin"
)

// Enum value maps for CORSResourcePolicy.
var (
	CORSResourcePolicy_name = map[int32]string{
		0: "ResourcePolicy_Unknown",
		1: "ResourceSameSite",
		2: "ResourceSameOrigin",
		3: "ResourceCrossOrigin",
	}
	CORSResourcePolicy_value = map[string]int32{
		"ResourcePolicy_Unknown": 0,
		"ResourceSameSite":       1,
		"ResourceSameOrigin":     2,
		"ResourceCrossOrigin":    3,
	}
)

func (x CORSResourcePolicy) Enum() *CORSResourcePolicy {
	p := new(CORSResourcePolicy)
	*p = x
	return p
}

func (x CORSResourcePolicy) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (CORSResourcePolicy) Descriptor() protoreflect.EnumDescriptor {
	return file_app_v1_middleware_cors_proto_enumTypes[2].Descriptor()
}

func (CORSResourcePolicy) Type() protoreflect.EnumType {
	return &file_app_v1_middleware_cors_proto_enumTypes[2]
}

func (x CORSResourcePolicy) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use CORSResourcePolicy.Descriptor instead.
func (CORSResourcePolicy) EnumDescriptor() ([]byte, []int) {
	return file_app_v1_middleware_cors_proto_rawDescGZIP(), []int{2}
}

// CORSMiddleware is the definition of the CORSMiddleware object.
// CORSMiddleware implements interface of the middleware.
type CORSMiddleware struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// [REQUIRED]
	// APIVersion is the defined version of the midleware.
	// This value must be "app/v1".
	APIVersion string `protobuf:"bytes,1,opt,name=APIVersion,json=apiVersion,proto3" json:"APIVersion,omitempty"`
	// [REQUIRED]
	// Kind is the kind of this object.
	// This value must be "CORSMiddleware".
	Kind string `protobuf:"bytes,2,opt,name=Kind,json=kind,proto3" json:"Kind,omitempty"`
	// [OPTIONAL]
	// Metadata is the metadata of the http logger object.
	// If not set, both name and namespace in the metadata
	// are treated as "default".
	Metadata *kernel.Metadata `protobuf:"bytes,3,opt,name=Metadata,json=metadata,proto3" json:"Metadata,omitempty"`
	// [OPTIONAL]
	// Spec is the specification of the middleware.
	// Default values are used when nothing is set.
	Spec          *CORSMiddlewareSpec `protobuf:"bytes,4,opt,name=Spec,json=spec,proto3" json:"Spec,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CORSMiddleware) Reset() {
	*x = CORSMiddleware{}
	mi := &file_app_v1_middleware_cors_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CORSMiddleware) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CORSMiddleware) ProtoMessage() {}

func (x *CORSMiddleware) ProtoReflect() protoreflect.Message {
	mi := &file_app_v1_middleware_cors_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CORSMiddleware.ProtoReflect.Descriptor instead.
func (*CORSMiddleware) Descriptor() ([]byte, []int) {
	return file_app_v1_middleware_cors_proto_rawDescGZIP(), []int{0}
}

func (x *CORSMiddleware) GetAPIVersion() string {
	if x != nil {
		return x.APIVersion
	}
	return ""
}

func (x *CORSMiddleware) GetKind() string {
	if x != nil {
		return x.Kind
	}
	return ""
}

func (x *CORSMiddleware) GetMetadata() *kernel.Metadata {
	if x != nil {
		return x.Metadata
	}
	return nil
}

func (x *CORSMiddleware) GetSpec() *CORSMiddlewareSpec {
	if x != nil {
		return x.Spec
	}
	return nil
}

// CORSMiddlewareSpec is the specifications for the CORSMiddleware object.
type CORSMiddlewareSpec struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// [OPTIONAL]
	// ErrorHandler is the reference to a ErrorHandler object.
	// Referred object must implement ErrorHandler interface.
	// Default error handler is used when not set.
	ErrorHandler *kernel.Reference `protobuf:"bytes,1,opt,name=ErrorHandler,json=errorHandler,proto3" json:"ErrorHandler,omitempty"`
	// [REQUIRED]
	// CORSPolicy is the configuration of CORS policy object.
	CORSPolicy    *CORSPolicySpec `protobuf:"bytes,2,opt,name=CORSPolicy,json=corsPolicy,proto3" json:"CORSPolicy,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CORSMiddlewareSpec) Reset() {
	*x = CORSMiddlewareSpec{}
	mi := &file_app_v1_middleware_cors_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CORSMiddlewareSpec) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CORSMiddlewareSpec) ProtoMessage() {}

func (x *CORSMiddlewareSpec) ProtoReflect() protoreflect.Message {
	mi := &file_app_v1_middleware_cors_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CORSMiddlewareSpec.ProtoReflect.Descriptor instead.
func (*CORSMiddlewareSpec) Descriptor() ([]byte, []int) {
	return file_app_v1_middleware_cors_proto_rawDescGZIP(), []int{1}
}

func (x *CORSMiddlewareSpec) GetErrorHandler() *kernel.Reference {
	if x != nil {
		return x.ErrorHandler
	}
	return nil
}

func (x *CORSMiddlewareSpec) GetCORSPolicy() *CORSPolicySpec {
	if x != nil {
		return x.CORSPolicy
	}
	return nil
}

// CORSPolicySpec is the specifications for the CORSPolicy object.
type CORSPolicySpec struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// [OPTIONAL]
	// AllowedOrigins is the list of origins to accept.
	// Only one origin is retuned as the value of "Access-Control-Allow-Origin" header if matched.
	// The wildcard origin "*" can be used here.
	// If nothing set, the wildcard origin "*" is used.
	// See https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Allow-Origin.
	// Default is ["*"].
	AllowedOrigins []string `protobuf:"bytes,1,rep,name=AllowedOrigins,json=allowedOrigins,proto3" json:"AllowedOrigins,omitempty"`
	// [OPTIONAL]
	// AllowedMethods is the list of HTTP methods to accept.
	// All methods should be listed. "ALL" is not allowed here for security reason.
	// See https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Allow-Methods.
	// Default is [POST, GET, OPTIONS].
	AllowedMethods []v1.HTTPMethod `protobuf:"varint,2,rep,packed,name=AllowedMethods,json=allowedMethods,proto3,enum=core.v1.HTTPMethod" json:"AllowedMethods,omitempty"`
	// [OPTIONAL]
	// AllowedHeaders is the list of HTTP header names to acccept.
	// The listed headers are responded to the clients as "Access-Control-Allow-Headers" header.
	// See https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Allow-Headers.
	// Set an empty string "" to disable all headers.
	// Default is ["Content-Type", "X-Requested-With"].
	AllowedHeaders []string `protobuf:"bytes,3,rep,name=AllowedHeaders,json=allowedHeaders,proto3" json:"AllowedHeaders,omitempty"`
	// [OPTIONAL]
	// ExposedHeaders are the list of HTTP header names to expose.
	// The listed headers are responded to the clients as "Access-Control-Expose-Headers" header.
	// See https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Expose-Headers.
	// Default is not set.
	ExposedHeaders []string `protobuf:"bytes,4,rep,name=ExposedHeaders,json=exposedHeaders,proto3" json:"ExposedHeaders,omitempty"`
	// [OPTIONAL]
	// AllowCredentials is the flag to allow credentials.
	// When this field is set to true, "Access-Control-Allow-Credentials: true" header is returned.
	// See https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Allow-Credentials.
	// Default is [false].
	AllowCredentials bool `protobuf:"varint,5,opt,name=AllowCredentials,json=allowCredentials,proto3" json:"AllowCredentials,omitempty"`
	// [OPTIONAL]
	// MaxAge is the duration that the user-agents can cache the preflight requests.
	// Mx age is returned to the client as "Access-Control-Max-Age" header when this fieled is set.
	// See https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Max-Age.
	// Default is [0] and the header is not returned.
	MaxAge int32 `protobuf:"varint,6,opt,name=MaxAge,json=maxAge,proto3" json:"MaxAge,omitempty"`
	// [OPTIONAL]
	// CORSEmbedderPolicy is the cross origin embedder policy to respond to the clients.
	// This header is returned in both preflight and actual requests.
	// See https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Cross-Origin-Embedder-Policy.
	// Default is not set.
	CORSEmbedderPolicy CORSEmbedderPolicy `protobuf:"varint,7,opt,name=CORSEmbedderPolicy,json=corsEmbedderPolicy,proto3,enum=app.v1.CORSEmbedderPolicy" json:"CORSEmbedderPolicy,omitempty"`
	// [OPTIONAL]
	// CORSOpenerPolicy is the cross origin opener policy to respond to the clients.
	// This header is returned in both preflight and actual requests.
	// See https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Cross-Origin-Opener-Policy.
	// Default is not set.
	CORSOpenerPolicy CORSOpenerPolicy `protobuf:"varint,8,opt,name=CORSOpenerPolicy,json=corsOpenerPolicy,proto3,enum=app.v1.CORSOpenerPolicy" json:"CORSOpenerPolicy,omitempty"`
	// [OPTIONAL]
	// CORSResourcePolicy is the cross origin resource policy to respond to the clients.
	// This header is returned in both preflight and actual requests.
	// See https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Cross-Origin-Resource-Policy.
	// Default is not set.
	CORSResourcePolicy CORSResourcePolicy `protobuf:"varint,9,opt,name=CORSResourcePolicy,json=corsResourcePolicy,proto3,enum=app.v1.CORSResourcePolicy" json:"CORSResourcePolicy,omitempty"`
	// [OPTIONAL]
	// AllowPrivateNetwork is the flag to allow shareing resources with external networks..
	// When this field is set to true, "Access-Control-Allow-Private-Network: true" header is returned.
	// See https://wicg.github.io/private-network-access/.
	// Default is [false].
	AllowPrivateNetwork bool `protobuf:"varint,10,opt,name=AllowPrivateNetwork,json=allowPrivateNetwork,proto3" json:"AllowPrivateNetwork,omitempty"`
	// [OPTIONAL]
	// DisableWildCardOrigin if true, set the requested origin to the
	// "Access-Control-Allow-Origin" header rather than the wildcard origin "*".
	// This is, in most cases, insecure than the wildcard origin "*".
	// This options is used when the AllowedOrigins contains wildcard origin "*".
	// Default is [false].
	DisableWildCardOrigin bool `protobuf:"varint,11,opt,name=DisableWildCardOrigin,json=disableWildCardOrigin,proto3" json:"DisableWildCardOrigin,omitempty"`
	unknownFields         protoimpl.UnknownFields
	sizeCache             protoimpl.SizeCache
}

func (x *CORSPolicySpec) Reset() {
	*x = CORSPolicySpec{}
	mi := &file_app_v1_middleware_cors_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CORSPolicySpec) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CORSPolicySpec) ProtoMessage() {}

func (x *CORSPolicySpec) ProtoReflect() protoreflect.Message {
	mi := &file_app_v1_middleware_cors_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CORSPolicySpec.ProtoReflect.Descriptor instead.
func (*CORSPolicySpec) Descriptor() ([]byte, []int) {
	return file_app_v1_middleware_cors_proto_rawDescGZIP(), []int{2}
}

func (x *CORSPolicySpec) GetAllowedOrigins() []string {
	if x != nil {
		return x.AllowedOrigins
	}
	return nil
}

func (x *CORSPolicySpec) GetAllowedMethods() []v1.HTTPMethod {
	if x != nil {
		return x.AllowedMethods
	}
	return nil
}

func (x *CORSPolicySpec) GetAllowedHeaders() []string {
	if x != nil {
		return x.AllowedHeaders
	}
	return nil
}

func (x *CORSPolicySpec) GetExposedHeaders() []string {
	if x != nil {
		return x.ExposedHeaders
	}
	return nil
}

func (x *CORSPolicySpec) GetAllowCredentials() bool {
	if x != nil {
		return x.AllowCredentials
	}
	return false
}

func (x *CORSPolicySpec) GetMaxAge() int32 {
	if x != nil {
		return x.MaxAge
	}
	return 0
}

func (x *CORSPolicySpec) GetCORSEmbedderPolicy() CORSEmbedderPolicy {
	if x != nil {
		return x.CORSEmbedderPolicy
	}
	return CORSEmbedderPolicy_EmbedderPolicy_Unknown
}

func (x *CORSPolicySpec) GetCORSOpenerPolicy() CORSOpenerPolicy {
	if x != nil {
		return x.CORSOpenerPolicy
	}
	return CORSOpenerPolicy_OpenerPolicy_Unknown
}

func (x *CORSPolicySpec) GetCORSResourcePolicy() CORSResourcePolicy {
	if x != nil {
		return x.CORSResourcePolicy
	}
	return CORSResourcePolicy_ResourcePolicy_Unknown
}

func (x *CORSPolicySpec) GetAllowPrivateNetwork() bool {
	if x != nil {
		return x.AllowPrivateNetwork
	}
	return false
}

func (x *CORSPolicySpec) GetDisableWildCardOrigin() bool {
	if x != nil {
		return x.DisableWildCardOrigin
	}
	return false
}

var File_app_v1_middleware_cors_proto protoreflect.FileDescriptor

var file_app_v1_middleware_cors_proto_rawDesc = string([]byte{
	0x0a, 0x1c, 0x61, 0x70, 0x70, 0x2f, 0x76, 0x31, 0x2f, 0x6d, 0x69, 0x64, 0x64, 0x6c, 0x65, 0x77,
	0x61, 0x72, 0x65, 0x2f, 0x63, 0x6f, 0x72, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x06,
	0x61, 0x70, 0x70, 0x2e, 0x76, 0x31, 0x1a, 0x1b, 0x62, 0x75, 0x66, 0x2f, 0x76, 0x61, 0x6c, 0x69,
	0x64, 0x61, 0x74, 0x65, 0x2f, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x1a, 0x12, 0x63, 0x6f, 0x72, 0x65, 0x2f, 0x76, 0x31, 0x2f, 0x68, 0x74, 0x74,
	0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x15, 0x6b, 0x65, 0x72, 0x6e, 0x65, 0x6c, 0x2f,
	0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xc8,
	0x01, 0x0a, 0x0e, 0x43, 0x4f, 0x52, 0x53, 0x4d, 0x69, 0x64, 0x64, 0x6c, 0x65, 0x77, 0x61, 0x72,
	0x65, 0x12, 0x2d, 0x0a, 0x0a, 0x41, 0x50, 0x49, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x0d, 0xba, 0x48, 0x0a, 0x72, 0x08, 0x0a, 0x06, 0x61, 0x70,
	0x70, 0x2f, 0x76, 0x31, 0x52, 0x0a, 0x61, 0x70, 0x69, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e,
	0x12, 0x29, 0x0a, 0x04, 0x4b, 0x69, 0x6e, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x42, 0x15,
	0xba, 0x48, 0x12, 0x72, 0x10, 0x0a, 0x0e, 0x43, 0x4f, 0x52, 0x53, 0x4d, 0x69, 0x64, 0x64, 0x6c,
	0x65, 0x77, 0x61, 0x72, 0x65, 0x52, 0x04, 0x6b, 0x69, 0x6e, 0x64, 0x12, 0x2c, 0x0a, 0x08, 0x4d,
	0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x10, 0x2e,
	0x6b, 0x65, 0x72, 0x6e, 0x65, 0x6c, 0x2e, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x52,
	0x08, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x12, 0x2e, 0x0a, 0x04, 0x53, 0x70, 0x65,
	0x63, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x61, 0x70, 0x70, 0x2e, 0x76, 0x31,
	0x2e, 0x43, 0x4f, 0x52, 0x53, 0x4d, 0x69, 0x64, 0x64, 0x6c, 0x65, 0x77, 0x61, 0x72, 0x65, 0x53,
	0x70, 0x65, 0x63, 0x52, 0x04, 0x73, 0x70, 0x65, 0x63, 0x22, 0x83, 0x01, 0x0a, 0x12, 0x43, 0x4f,
	0x52, 0x53, 0x4d, 0x69, 0x64, 0x64, 0x6c, 0x65, 0x77, 0x61, 0x72, 0x65, 0x53, 0x70, 0x65, 0x63,
	0x12, 0x35, 0x0a, 0x0c, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x48, 0x61, 0x6e, 0x64, 0x6c, 0x65, 0x72,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x11, 0x2e, 0x6b, 0x65, 0x72, 0x6e, 0x65, 0x6c, 0x2e,
	0x52, 0x65, 0x66, 0x65, 0x72, 0x65, 0x6e, 0x63, 0x65, 0x52, 0x0c, 0x65, 0x72, 0x72, 0x6f, 0x72,
	0x48, 0x61, 0x6e, 0x64, 0x6c, 0x65, 0x72, 0x12, 0x36, 0x0a, 0x0a, 0x43, 0x4f, 0x52, 0x53, 0x50,
	0x6f, 0x6c, 0x69, 0x63, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x16, 0x2e, 0x61, 0x70,
	0x70, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x4f, 0x52, 0x53, 0x50, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x53,
	0x70, 0x65, 0x63, 0x52, 0x0a, 0x63, 0x6f, 0x72, 0x73, 0x50, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x22,
	0xcf, 0x04, 0x0a, 0x0e, 0x43, 0x4f, 0x52, 0x53, 0x50, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x53, 0x70,
	0x65, 0x63, 0x12, 0x26, 0x0a, 0x0e, 0x41, 0x6c, 0x6c, 0x6f, 0x77, 0x65, 0x64, 0x4f, 0x72, 0x69,
	0x67, 0x69, 0x6e, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x09, 0x52, 0x0e, 0x61, 0x6c, 0x6c, 0x6f,
	0x77, 0x65, 0x64, 0x4f, 0x72, 0x69, 0x67, 0x69, 0x6e, 0x73, 0x12, 0x3b, 0x0a, 0x0e, 0x41, 0x6c,
	0x6c, 0x6f, 0x77, 0x65, 0x64, 0x4d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x73, 0x18, 0x02, 0x20, 0x03,
	0x28, 0x0e, 0x32, 0x13, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x48, 0x54, 0x54,
	0x50, 0x4d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x52, 0x0e, 0x61, 0x6c, 0x6c, 0x6f, 0x77, 0x65, 0x64,
	0x4d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x73, 0x12, 0x26, 0x0a, 0x0e, 0x41, 0x6c, 0x6c, 0x6f, 0x77,
	0x65, 0x64, 0x48, 0x65, 0x61, 0x64, 0x65, 0x72, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x09, 0x52,
	0x0e, 0x61, 0x6c, 0x6c, 0x6f, 0x77, 0x65, 0x64, 0x48, 0x65, 0x61, 0x64, 0x65, 0x72, 0x73, 0x12,
	0x26, 0x0a, 0x0e, 0x45, 0x78, 0x70, 0x6f, 0x73, 0x65, 0x64, 0x48, 0x65, 0x61, 0x64, 0x65, 0x72,
	0x73, 0x18, 0x04, 0x20, 0x03, 0x28, 0x09, 0x52, 0x0e, 0x65, 0x78, 0x70, 0x6f, 0x73, 0x65, 0x64,
	0x48, 0x65, 0x61, 0x64, 0x65, 0x72, 0x73, 0x12, 0x2a, 0x0a, 0x10, 0x41, 0x6c, 0x6c, 0x6f, 0x77,
	0x43, 0x72, 0x65, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x61, 0x6c, 0x73, 0x18, 0x05, 0x20, 0x01, 0x28,
	0x08, 0x52, 0x10, 0x61, 0x6c, 0x6c, 0x6f, 0x77, 0x43, 0x72, 0x65, 0x64, 0x65, 0x6e, 0x74, 0x69,
	0x61, 0x6c, 0x73, 0x12, 0x16, 0x0a, 0x06, 0x4d, 0x61, 0x78, 0x41, 0x67, 0x65, 0x18, 0x06, 0x20,
	0x01, 0x28, 0x05, 0x52, 0x06, 0x6d, 0x61, 0x78, 0x41, 0x67, 0x65, 0x12, 0x4a, 0x0a, 0x12, 0x43,
	0x4f, 0x52, 0x53, 0x45, 0x6d, 0x62, 0x65, 0x64, 0x64, 0x65, 0x72, 0x50, 0x6f, 0x6c, 0x69, 0x63,
	0x79, 0x18, 0x07, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x1a, 0x2e, 0x61, 0x70, 0x70, 0x2e, 0x76, 0x31,
	0x2e, 0x43, 0x4f, 0x52, 0x53, 0x45, 0x6d, 0x62, 0x65, 0x64, 0x64, 0x65, 0x72, 0x50, 0x6f, 0x6c,
	0x69, 0x63, 0x79, 0x52, 0x12, 0x63, 0x6f, 0x72, 0x73, 0x45, 0x6d, 0x62, 0x65, 0x64, 0x64, 0x65,
	0x72, 0x50, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x12, 0x44, 0x0a, 0x10, 0x43, 0x4f, 0x52, 0x53, 0x4f,
	0x70, 0x65, 0x6e, 0x65, 0x72, 0x50, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x18, 0x08, 0x20, 0x01, 0x28,
	0x0e, 0x32, 0x18, 0x2e, 0x61, 0x70, 0x70, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x4f, 0x52, 0x53, 0x4f,
	0x70, 0x65, 0x6e, 0x65, 0x72, 0x50, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x52, 0x10, 0x63, 0x6f, 0x72,
	0x73, 0x4f, 0x70, 0x65, 0x6e, 0x65, 0x72, 0x50, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x12, 0x4a, 0x0a,
	0x12, 0x43, 0x4f, 0x52, 0x53, 0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x50, 0x6f, 0x6c,
	0x69, 0x63, 0x79, 0x18, 0x09, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x1a, 0x2e, 0x61, 0x70, 0x70, 0x2e,
	0x76, 0x31, 0x2e, 0x43, 0x4f, 0x52, 0x53, 0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x50,
	0x6f, 0x6c, 0x69, 0x63, 0x79, 0x52, 0x12, 0x63, 0x6f, 0x72, 0x73, 0x52, 0x65, 0x73, 0x6f, 0x75,
	0x72, 0x63, 0x65, 0x50, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x12, 0x30, 0x0a, 0x13, 0x41, 0x6c, 0x6c,
	0x6f, 0x77, 0x50, 0x72, 0x69, 0x76, 0x61, 0x74, 0x65, 0x4e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b,
	0x18, 0x0a, 0x20, 0x01, 0x28, 0x08, 0x52, 0x13, 0x61, 0x6c, 0x6c, 0x6f, 0x77, 0x50, 0x72, 0x69,
	0x76, 0x61, 0x74, 0x65, 0x4e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x12, 0x34, 0x0a, 0x15, 0x44,
	0x69, 0x73, 0x61, 0x62, 0x6c, 0x65, 0x57, 0x69, 0x6c, 0x64, 0x43, 0x61, 0x72, 0x64, 0x4f, 0x72,
	0x69, 0x67, 0x69, 0x6e, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x08, 0x52, 0x15, 0x64, 0x69, 0x73, 0x61,
	0x62, 0x6c, 0x65, 0x57, 0x69, 0x6c, 0x64, 0x43, 0x61, 0x72, 0x64, 0x4f, 0x72, 0x69, 0x67, 0x69,
	0x6e, 0x2a, 0x7d, 0x0a, 0x12, 0x43, 0x4f, 0x52, 0x53, 0x45, 0x6d, 0x62, 0x65, 0x64, 0x64, 0x65,
	0x72, 0x50, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x12, 0x1a, 0x0a, 0x16, 0x45, 0x6d, 0x62, 0x65, 0x64,
	0x64, 0x65, 0x72, 0x50, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x5f, 0x55, 0x6e, 0x6b, 0x6e, 0x6f, 0x77,
	0x6e, 0x10, 0x00, 0x12, 0x16, 0x0a, 0x12, 0x45, 0x6d, 0x62, 0x65, 0x64, 0x64, 0x65, 0x72, 0x55,
	0x6e, 0x73, 0x61, 0x66, 0x65, 0x4e, 0x6f, 0x6e, 0x65, 0x10, 0x01, 0x12, 0x17, 0x0a, 0x13, 0x45,
	0x6d, 0x62, 0x65, 0x64, 0x64, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x69, 0x72, 0x65, 0x43, 0x6f,
	0x72, 0x70, 0x10, 0x02, 0x12, 0x1a, 0x0a, 0x16, 0x45, 0x6d, 0x62, 0x65, 0x64, 0x64, 0x65, 0x72,
	0x43, 0x72, 0x65, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x61, 0x6c, 0x6c, 0x65, 0x73, 0x73, 0x10, 0x03,
	0x2a, 0x79, 0x0a, 0x10, 0x43, 0x4f, 0x52, 0x53, 0x4f, 0x70, 0x65, 0x6e, 0x65, 0x72, 0x50, 0x6f,
	0x6c, 0x69, 0x63, 0x79, 0x12, 0x18, 0x0a, 0x14, 0x4f, 0x70, 0x65, 0x6e, 0x65, 0x72, 0x50, 0x6f,
	0x6c, 0x69, 0x63, 0x79, 0x5f, 0x55, 0x6e, 0x6b, 0x6e, 0x6f, 0x77, 0x6e, 0x10, 0x00, 0x12, 0x14,
	0x0a, 0x10, 0x4f, 0x70, 0x65, 0x6e, 0x65, 0x72, 0x55, 0x6e, 0x73, 0x61, 0x66, 0x65, 0x4e, 0x6f,
	0x6e, 0x65, 0x10, 0x01, 0x12, 0x1f, 0x0a, 0x1b, 0x4f, 0x70, 0x65, 0x6e, 0x65, 0x72, 0x53, 0x61,
	0x6d, 0x65, 0x4f, 0x72, 0x69, 0x67, 0x69, 0x6e, 0x41, 0x6c, 0x6c, 0x6f, 0x77, 0x50, 0x6f, 0x70,
	0x75, 0x70, 0x73, 0x10, 0x02, 0x12, 0x14, 0x0a, 0x10, 0x4f, 0x70, 0x65, 0x6e, 0x65, 0x72, 0x53,
	0x61, 0x6d, 0x65, 0x4f, 0x72, 0x69, 0x67, 0x69, 0x6e, 0x10, 0x03, 0x2a, 0x77, 0x0a, 0x12, 0x43,
	0x4f, 0x52, 0x53, 0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x50, 0x6f, 0x6c, 0x69, 0x63,
	0x79, 0x12, 0x1a, 0x0a, 0x16, 0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x50, 0x6f, 0x6c,
	0x69, 0x63, 0x79, 0x5f, 0x55, 0x6e, 0x6b, 0x6e, 0x6f, 0x77, 0x6e, 0x10, 0x00, 0x12, 0x14, 0x0a,
	0x10, 0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x53, 0x61, 0x6d, 0x65, 0x53, 0x69, 0x74,
	0x65, 0x10, 0x01, 0x12, 0x16, 0x0a, 0x12, 0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x53,
	0x61, 0x6d, 0x65, 0x4f, 0x72, 0x69, 0x67, 0x69, 0x6e, 0x10, 0x02, 0x12, 0x17, 0x0a, 0x13, 0x52,
	0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x43, 0x72, 0x6f, 0x73, 0x73, 0x4f, 0x72, 0x69, 0x67,
	0x69, 0x6e, 0x10, 0x03, 0x42, 0x38, 0x5a, 0x36, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63,
	0x6f, 0x6d, 0x2f, 0x61, 0x69, 0x6c, 0x65, 0x72, 0x6f, 0x6e, 0x2d, 0x67, 0x61, 0x74, 0x65, 0x77,
	0x61, 0x79, 0x2f, 0x61, 0x69, 0x6c, 0x65, 0x72, 0x6f, 0x6e, 0x2d, 0x67, 0x61, 0x74, 0x65, 0x77,
	0x61, 0x79, 0x2f, 0x61, 0x70, 0x69, 0x73, 0x2f, 0x61, 0x70, 0x70, 0x2f, 0x76, 0x31, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
})

var (
	file_app_v1_middleware_cors_proto_rawDescOnce sync.Once
	file_app_v1_middleware_cors_proto_rawDescData []byte
)

func file_app_v1_middleware_cors_proto_rawDescGZIP() []byte {
	file_app_v1_middleware_cors_proto_rawDescOnce.Do(func() {
		file_app_v1_middleware_cors_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_app_v1_middleware_cors_proto_rawDesc), len(file_app_v1_middleware_cors_proto_rawDesc)))
	})
	return file_app_v1_middleware_cors_proto_rawDescData
}

var file_app_v1_middleware_cors_proto_enumTypes = make([]protoimpl.EnumInfo, 3)
var file_app_v1_middleware_cors_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_app_v1_middleware_cors_proto_goTypes = []any{
	(CORSEmbedderPolicy)(0),    // 0: app.v1.CORSEmbedderPolicy
	(CORSOpenerPolicy)(0),      // 1: app.v1.CORSOpenerPolicy
	(CORSResourcePolicy)(0),    // 2: app.v1.CORSResourcePolicy
	(*CORSMiddleware)(nil),     // 3: app.v1.CORSMiddleware
	(*CORSMiddlewareSpec)(nil), // 4: app.v1.CORSMiddlewareSpec
	(*CORSPolicySpec)(nil),     // 5: app.v1.CORSPolicySpec
	(*kernel.Metadata)(nil),    // 6: kernel.Metadata
	(*kernel.Reference)(nil),   // 7: kernel.Reference
	(v1.HTTPMethod)(0),         // 8: core.v1.HTTPMethod
}
var file_app_v1_middleware_cors_proto_depIdxs = []int32{
	6, // 0: app.v1.CORSMiddleware.Metadata:type_name -> kernel.Metadata
	4, // 1: app.v1.CORSMiddleware.Spec:type_name -> app.v1.CORSMiddlewareSpec
	7, // 2: app.v1.CORSMiddlewareSpec.ErrorHandler:type_name -> kernel.Reference
	5, // 3: app.v1.CORSMiddlewareSpec.CORSPolicy:type_name -> app.v1.CORSPolicySpec
	8, // 4: app.v1.CORSPolicySpec.AllowedMethods:type_name -> core.v1.HTTPMethod
	0, // 5: app.v1.CORSPolicySpec.CORSEmbedderPolicy:type_name -> app.v1.CORSEmbedderPolicy
	1, // 6: app.v1.CORSPolicySpec.CORSOpenerPolicy:type_name -> app.v1.CORSOpenerPolicy
	2, // 7: app.v1.CORSPolicySpec.CORSResourcePolicy:type_name -> app.v1.CORSResourcePolicy
	8, // [8:8] is the sub-list for method output_type
	8, // [8:8] is the sub-list for method input_type
	8, // [8:8] is the sub-list for extension type_name
	8, // [8:8] is the sub-list for extension extendee
	0, // [0:8] is the sub-list for field type_name
}

func init() { file_app_v1_middleware_cors_proto_init() }
func file_app_v1_middleware_cors_proto_init() {
	if File_app_v1_middleware_cors_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_app_v1_middleware_cors_proto_rawDesc), len(file_app_v1_middleware_cors_proto_rawDesc)),
			NumEnums:      3,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_app_v1_middleware_cors_proto_goTypes,
		DependencyIndexes: file_app_v1_middleware_cors_proto_depIdxs,
		EnumInfos:         file_app_v1_middleware_cors_proto_enumTypes,
		MessageInfos:      file_app_v1_middleware_cors_proto_msgTypes,
	}.Build()
	File_app_v1_middleware_cors_proto = out.File
	file_app_v1_middleware_cors_proto_goTypes = nil
	file_app_v1_middleware_cors_proto_depIdxs = nil
}
