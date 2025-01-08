// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.1
// 	protoc        v5.27.2
// source: app/v1/authn/basic.proto

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

// BasicAuthnMiddleware is the definition of the BasicAuthnMiddleware object.
// BasicAuthnMiddleware implements interface of the middleware.
type BasicAuthnMiddleware struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// [REQUIRED]
	// APIVersion is the defined version of the midleware.
	// This value must be "app/v1".
	APIVersion string `protobuf:"bytes,1,opt,name=APIVersion,json=apiVersion,proto3" json:"APIVersion,omitempty"`
	// [REQUIRED]
	// Kind is the kind of this object.
	// This value must be "BasicAuthnMiddleware".
	Kind string `protobuf:"bytes,2,opt,name=Kind,json=kind,proto3" json:"Kind,omitempty"`
	// [OPTIONAL]
	// Metadata is the metadata of the http logger object.
	// If not set, both name and namespace in the metadata
	// are treated as "default".
	Metadata *kernel.Metadata `protobuf:"bytes,3,opt,name=Metadata,json=metadata,proto3" json:"Metadata,omitempty"`
	// [OPTIONAL]
	// Spec is the specification of the middleware.
	// Default values are used when nothing is set.
	Spec          *BasicAuthnMiddlewareSpec `protobuf:"bytes,4,opt,name=Spec,json=spec,proto3" json:"Spec,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *BasicAuthnMiddleware) Reset() {
	*x = BasicAuthnMiddleware{}
	mi := &file_app_v1_authn_basic_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *BasicAuthnMiddleware) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BasicAuthnMiddleware) ProtoMessage() {}

func (x *BasicAuthnMiddleware) ProtoReflect() protoreflect.Message {
	mi := &file_app_v1_authn_basic_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BasicAuthnMiddleware.ProtoReflect.Descriptor instead.
func (*BasicAuthnMiddleware) Descriptor() ([]byte, []int) {
	return file_app_v1_authn_basic_proto_rawDescGZIP(), []int{0}
}

func (x *BasicAuthnMiddleware) GetAPIVersion() string {
	if x != nil {
		return x.APIVersion
	}
	return ""
}

func (x *BasicAuthnMiddleware) GetKind() string {
	if x != nil {
		return x.Kind
	}
	return ""
}

func (x *BasicAuthnMiddleware) GetMetadata() *kernel.Metadata {
	if x != nil {
		return x.Metadata
	}
	return nil
}

func (x *BasicAuthnMiddleware) GetSpec() *BasicAuthnMiddlewareSpec {
	if x != nil {
		return x.Spec
	}
	return nil
}

// BasicAuthnMiddlewareSpec is the specifications for the BasicAuthnMiddleware object.
type BasicAuthnMiddlewareSpec struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// [OPTIONAL]
	// Logger is the reference to a Logger object.
	// Referred object must implement Logger interface.
	// Default Logger is used when not set.
	Logger *kernel.Reference `protobuf:"bytes,1,opt,name=Logger,json=logger,proto3" json:"Logger,omitempty"`
	// [OPTIONAL]
	// ErrorHandler is the reference to a ErrorHandler object.
	// Referred object must implement ErrorHandler interface.
	// Default error handler is used when not set.
	ErrorHandler *kernel.Reference `protobuf:"bytes,2,opt,name=ErrorHandler,json=errorHandler,proto3" json:"ErrorHandler,omitempty"`
	// [OPTIONAL]
	// ClaimsKey is the key to set user attibutes in the context.
	// Claims can be used for authorization if necessary.
	// If not set, default value is used.
	// Default is ["AuthnClaims"].
	ClaimsKey string `protobuf:"bytes,3,opt,name=ClaimsKey,json=claimsKey,proto3" json:"ClaimsKey,omitempty"`
	// [OPTIONAL]
	// KeepCredentials is the flag to keep credentials in the header.
	// That means Authorization header is not removed in the middleware.
	// If true, Authorizatoin header will be proxied upstream services.
	// Default is [false].
	KeepCredentials bool `protobuf:"varint,4,opt,name=KeepCredentials,json=keepCredentials,proto3" json:"KeepCredentials,omitempty"`
	// [OPTIONAL]
	// PasswordCrypt is the password encryption, or password hashing config.
	// If not set, the passwords are used as is.
	// Default is not set.
	PasswordCrypt *kernel.PasswordCryptSpec `protobuf:"bytes,5,opt,name=PasswordCrypt,json=passwordCrypt,proto3" json:"PasswordCrypt,omitempty"`
	// [OPTIONAL]
	// CommonKeyCryptType is the common key encryption algorithm
	// that is used for encrypting passwords of hashed passwords.
	// Common key encryption will be enabled when CryptSecret is not empty.
	// If CryptSecret is not empty, CommonKeyCryptType should also be set
	// to specify the encryption algorithm.
	// PasswordCrypt and CommonKeyCryptType can be conbined.
	// If so, the password should be CommonKeyCrypt(PasswordCrypt(<Password>))
	// with base64 or hex encoding.
	// Default is not set.
	CommonKeyCryptType kernel.CommonKeyCryptType `protobuf:"varint,6,opt,name=CommonKeyCryptType,json=commonKeyCryptType,proto3,enum=kernel.CommonKeyCryptType" json:"CommonKeyCryptType,omitempty"`
	// [OPTIONAL]
	// CryptSecret is the Base64 encoded encryption key.
	// Base64 standard encoded with padding is ecpected to be used.
	// Common key encryption will be enabled when CryptSecret is not empty.
	// If CryptSecret is not empty, CommonKeyCryptType should also be set
	// to specify the encryption algorithm.
	// PasswordCrypt and CommonKeyCryptType can be conbined.
	// If so, the password should be CommonKeyCrypt(PasswordCrypt(<Password>))
	// with base64 or hex encoding.
	// Default is not set.
	CryptSecret string `protobuf:"bytes,7,opt,name=CryptSecret,json=cryptSecret,proto3" json:"CryptSecret,omitempty"`
	// [OPTIONAL]
	// Realm is the realm name of authentication.
	// If not set, an empty string will be used.
	// Default is not set, or empty string [""].
	Realm string `protobuf:"bytes,8,opt,name=Realm,json=realm,proto3" json:"Realm,omitempty"`
	// [OPTIONAL]
	// Providers is the credentials provider to use.
	// If not set, EnvProvider with default values are used.
	//
	// Types that are valid to be assigned to Providers:
	//
	//	*BasicAuthnMiddlewareSpec_EnvProvider
	//	*BasicAuthnMiddlewareSpec_FileProvider
	Providers     isBasicAuthnMiddlewareSpec_Providers `protobuf_oneof:"Providers"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *BasicAuthnMiddlewareSpec) Reset() {
	*x = BasicAuthnMiddlewareSpec{}
	mi := &file_app_v1_authn_basic_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *BasicAuthnMiddlewareSpec) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BasicAuthnMiddlewareSpec) ProtoMessage() {}

func (x *BasicAuthnMiddlewareSpec) ProtoReflect() protoreflect.Message {
	mi := &file_app_v1_authn_basic_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BasicAuthnMiddlewareSpec.ProtoReflect.Descriptor instead.
func (*BasicAuthnMiddlewareSpec) Descriptor() ([]byte, []int) {
	return file_app_v1_authn_basic_proto_rawDescGZIP(), []int{1}
}

func (x *BasicAuthnMiddlewareSpec) GetLogger() *kernel.Reference {
	if x != nil {
		return x.Logger
	}
	return nil
}

func (x *BasicAuthnMiddlewareSpec) GetErrorHandler() *kernel.Reference {
	if x != nil {
		return x.ErrorHandler
	}
	return nil
}

func (x *BasicAuthnMiddlewareSpec) GetClaimsKey() string {
	if x != nil {
		return x.ClaimsKey
	}
	return ""
}

func (x *BasicAuthnMiddlewareSpec) GetKeepCredentials() bool {
	if x != nil {
		return x.KeepCredentials
	}
	return false
}

func (x *BasicAuthnMiddlewareSpec) GetPasswordCrypt() *kernel.PasswordCryptSpec {
	if x != nil {
		return x.PasswordCrypt
	}
	return nil
}

func (x *BasicAuthnMiddlewareSpec) GetCommonKeyCryptType() kernel.CommonKeyCryptType {
	if x != nil {
		return x.CommonKeyCryptType
	}
	return kernel.CommonKeyCryptType(0)
}

func (x *BasicAuthnMiddlewareSpec) GetCryptSecret() string {
	if x != nil {
		return x.CryptSecret
	}
	return ""
}

func (x *BasicAuthnMiddlewareSpec) GetRealm() string {
	if x != nil {
		return x.Realm
	}
	return ""
}

func (x *BasicAuthnMiddlewareSpec) GetProviders() isBasicAuthnMiddlewareSpec_Providers {
	if x != nil {
		return x.Providers
	}
	return nil
}

func (x *BasicAuthnMiddlewareSpec) GetEnvProvider() *BasicAuthnEnvProvider {
	if x != nil {
		if x, ok := x.Providers.(*BasicAuthnMiddlewareSpec_EnvProvider); ok {
			return x.EnvProvider
		}
	}
	return nil
}

func (x *BasicAuthnMiddlewareSpec) GetFileProvider() *BasicAuthnFileProvider {
	if x != nil {
		if x, ok := x.Providers.(*BasicAuthnMiddlewareSpec_FileProvider); ok {
			return x.FileProvider
		}
	}
	return nil
}

type isBasicAuthnMiddlewareSpec_Providers interface {
	isBasicAuthnMiddlewareSpec_Providers()
}

type BasicAuthnMiddlewareSpec_EnvProvider struct {
	EnvProvider *BasicAuthnEnvProvider `protobuf:"bytes,10,opt,name=EnvProvider,json=envProvider,proto3,oneof"`
}

type BasicAuthnMiddlewareSpec_FileProvider struct {
	FileProvider *BasicAuthnFileProvider `protobuf:"bytes,11,opt,name=FileProvider,json=fileProvider,proto3,oneof"`
}

func (*BasicAuthnMiddlewareSpec_EnvProvider) isBasicAuthnMiddlewareSpec_Providers() {}

func (*BasicAuthnMiddlewareSpec_FileProvider) isBasicAuthnMiddlewareSpec_Providers() {}

type BasicAuthnEnvProvider struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// [OPTIONAL]
	// UsernamePrefix is the prefix of environmental variable
	// that provides username.
	// If the prefix is "USERNAME_", then usernames should be
	// set like "USERNAME_1=foo", "USERNAME_999=bar".
	// Note that the numbers can be zero padded which means
	// "USERNAME_1" and "USERNAME_001" are the same.
	// Both username and password must be set for each single users.
	// If empty string was set, default value is used.
	// Defailt is [GATEWAY_BASIC_USERNAME_].
	UsernamePrefix string `protobuf:"bytes,1,opt,name=UsernamePrefix,json=usernamePrefix,proto3" json:"UsernamePrefix,omitempty"`
	// [OPTIONAL]
	// PasswordPrefix is the prefix of environmental variable
	// that provides passwords.
	// If the prefix is "PASSWORDS_", then passwords should be
	// set like "PASSWORDS_1=foo", "PASSWORDS_999=bar".
	// Note that the numbers can be zero padded which means
	// "PASSWORDS_1" and "PASSWORDS_001" are the same.
	// Both username and password must be set for each single users.
	// If empty string was set, default value is used.
	// Defailt is [GATEWAY_BASIC_PASSWORD_]
	PasswordPrefix string `protobuf:"bytes,2,opt,name=PasswordPrefix,json=passwordPrefix,proto3" json:"PasswordPrefix,omitempty"`
	// [OPTIONAL]
	// Encoding is the encoding algorithm used to decode passwords.
	// If set, all password strings are decoded with configured encoding.
	// Gateway will fail to start when failed to decoding.
	// Default is [false].
	Encoding      kernel.EncodingType `protobuf:"varint,3,opt,name=Encoding,json=encoding,proto3,enum=kernel.EncodingType" json:"Encoding,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *BasicAuthnEnvProvider) Reset() {
	*x = BasicAuthnEnvProvider{}
	mi := &file_app_v1_authn_basic_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *BasicAuthnEnvProvider) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BasicAuthnEnvProvider) ProtoMessage() {}

func (x *BasicAuthnEnvProvider) ProtoReflect() protoreflect.Message {
	mi := &file_app_v1_authn_basic_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BasicAuthnEnvProvider.ProtoReflect.Descriptor instead.
func (*BasicAuthnEnvProvider) Descriptor() ([]byte, []int) {
	return file_app_v1_authn_basic_proto_rawDescGZIP(), []int{2}
}

func (x *BasicAuthnEnvProvider) GetUsernamePrefix() string {
	if x != nil {
		return x.UsernamePrefix
	}
	return ""
}

func (x *BasicAuthnEnvProvider) GetPasswordPrefix() string {
	if x != nil {
		return x.PasswordPrefix
	}
	return ""
}

func (x *BasicAuthnEnvProvider) GetEncoding() kernel.EncodingType {
	if x != nil {
		return x.Encoding
	}
	return kernel.EncodingType(0)
}

type BasicAuthnFileProvider struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// [OPTIONAL]
	// Paths are file paths that contains use information.
	// If nothing set, all authentication challenge will fail.
	// Default is not set.
	Paths []string `protobuf:"bytes,1,rep,name=Paths,json=paths,proto3" json:"Paths,omitempty"`
	// [OPTIONAL]
	// Encoding is the encoding algorithm used to decode passwords.
	// If set, all password strings are decoded with configured encoding.
	// Gateway will fail to start when failed to decoding.
	// Default is [false].
	Encoding      kernel.EncodingType `protobuf:"varint,2,opt,name=Encoding,json=encoding,proto3,enum=kernel.EncodingType" json:"Encoding,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *BasicAuthnFileProvider) Reset() {
	*x = BasicAuthnFileProvider{}
	mi := &file_app_v1_authn_basic_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *BasicAuthnFileProvider) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BasicAuthnFileProvider) ProtoMessage() {}

func (x *BasicAuthnFileProvider) ProtoReflect() protoreflect.Message {
	mi := &file_app_v1_authn_basic_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BasicAuthnFileProvider.ProtoReflect.Descriptor instead.
func (*BasicAuthnFileProvider) Descriptor() ([]byte, []int) {
	return file_app_v1_authn_basic_proto_rawDescGZIP(), []int{3}
}

func (x *BasicAuthnFileProvider) GetPaths() []string {
	if x != nil {
		return x.Paths
	}
	return nil
}

func (x *BasicAuthnFileProvider) GetEncoding() kernel.EncodingType {
	if x != nil {
		return x.Encoding
	}
	return kernel.EncodingType(0)
}

var File_app_v1_authn_basic_proto protoreflect.FileDescriptor

var file_app_v1_authn_basic_proto_rawDesc = []byte{
	0x0a, 0x18, 0x61, 0x70, 0x70, 0x2f, 0x76, 0x31, 0x2f, 0x61, 0x75, 0x74, 0x68, 0x6e, 0x2f, 0x62,
	0x61, 0x73, 0x69, 0x63, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x06, 0x61, 0x70, 0x70, 0x2e,
	0x76, 0x31, 0x1a, 0x1b, 0x62, 0x75, 0x66, 0x2f, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65,
	0x2f, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a,
	0x12, 0x6b, 0x65, 0x72, 0x6e, 0x65, 0x6c, 0x2f, 0x63, 0x72, 0x79, 0x70, 0x74, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x1a, 0x15, 0x6b, 0x65, 0x72, 0x6e, 0x65, 0x6c, 0x2f, 0x72, 0x65, 0x73, 0x6f,
	0x75, 0x72, 0x63, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x15, 0x6b, 0x65, 0x72, 0x6e,
	0x65, 0x6c, 0x2f, 0x65, 0x6e, 0x63, 0x6f, 0x64, 0x69, 0x6e, 0x67, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x22, 0xda, 0x01, 0x0a, 0x14, 0x42, 0x61, 0x73, 0x69, 0x63, 0x41, 0x75, 0x74, 0x68, 0x6e,
	0x4d, 0x69, 0x64, 0x64, 0x6c, 0x65, 0x77, 0x61, 0x72, 0x65, 0x12, 0x2d, 0x0a, 0x0a, 0x41, 0x50,
	0x49, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x0d,
	0xba, 0x48, 0x0a, 0x72, 0x08, 0x0a, 0x06, 0x61, 0x70, 0x70, 0x2f, 0x76, 0x31, 0x52, 0x0a, 0x61,
	0x70, 0x69, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x2f, 0x0a, 0x04, 0x4b, 0x69, 0x6e,
	0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x42, 0x1b, 0xba, 0x48, 0x18, 0x72, 0x16, 0x0a, 0x14,
	0x42, 0x61, 0x73, 0x69, 0x63, 0x41, 0x75, 0x74, 0x68, 0x6e, 0x4d, 0x69, 0x64, 0x64, 0x6c, 0x65,
	0x77, 0x61, 0x72, 0x65, 0x52, 0x04, 0x6b, 0x69, 0x6e, 0x64, 0x12, 0x2c, 0x0a, 0x08, 0x4d, 0x65,
	0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x6b,
	0x65, 0x72, 0x6e, 0x65, 0x6c, 0x2e, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x52, 0x08,
	0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x12, 0x34, 0x0a, 0x04, 0x53, 0x70, 0x65, 0x63,
	0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x20, 0x2e, 0x61, 0x70, 0x70, 0x2e, 0x76, 0x31, 0x2e,
	0x42, 0x61, 0x73, 0x69, 0x63, 0x41, 0x75, 0x74, 0x68, 0x6e, 0x4d, 0x69, 0x64, 0x64, 0x6c, 0x65,
	0x77, 0x61, 0x72, 0x65, 0x53, 0x70, 0x65, 0x63, 0x52, 0x04, 0x73, 0x70, 0x65, 0x63, 0x22, 0x9f,
	0x04, 0x0a, 0x18, 0x42, 0x61, 0x73, 0x69, 0x63, 0x41, 0x75, 0x74, 0x68, 0x6e, 0x4d, 0x69, 0x64,
	0x64, 0x6c, 0x65, 0x77, 0x61, 0x72, 0x65, 0x53, 0x70, 0x65, 0x63, 0x12, 0x29, 0x0a, 0x06, 0x4c,
	0x6f, 0x67, 0x67, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x11, 0x2e, 0x6b, 0x65,
	0x72, 0x6e, 0x65, 0x6c, 0x2e, 0x52, 0x65, 0x66, 0x65, 0x72, 0x65, 0x6e, 0x63, 0x65, 0x52, 0x06,
	0x6c, 0x6f, 0x67, 0x67, 0x65, 0x72, 0x12, 0x35, 0x0a, 0x0c, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x48,
	0x61, 0x6e, 0x64, 0x6c, 0x65, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x11, 0x2e, 0x6b,
	0x65, 0x72, 0x6e, 0x65, 0x6c, 0x2e, 0x52, 0x65, 0x66, 0x65, 0x72, 0x65, 0x6e, 0x63, 0x65, 0x52,
	0x0c, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x48, 0x61, 0x6e, 0x64, 0x6c, 0x65, 0x72, 0x12, 0x1c, 0x0a,
	0x09, 0x43, 0x6c, 0x61, 0x69, 0x6d, 0x73, 0x4b, 0x65, 0x79, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x09, 0x63, 0x6c, 0x61, 0x69, 0x6d, 0x73, 0x4b, 0x65, 0x79, 0x12, 0x28, 0x0a, 0x0f, 0x4b,
	0x65, 0x65, 0x70, 0x43, 0x72, 0x65, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x61, 0x6c, 0x73, 0x18, 0x04,
	0x20, 0x01, 0x28, 0x08, 0x52, 0x0f, 0x6b, 0x65, 0x65, 0x70, 0x43, 0x72, 0x65, 0x64, 0x65, 0x6e,
	0x74, 0x69, 0x61, 0x6c, 0x73, 0x12, 0x3f, 0x0a, 0x0d, 0x50, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72,
	0x64, 0x43, 0x72, 0x79, 0x70, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x6b,
	0x65, 0x72, 0x6e, 0x65, 0x6c, 0x2e, 0x50, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x43, 0x72,
	0x79, 0x70, 0x74, 0x53, 0x70, 0x65, 0x63, 0x52, 0x0d, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72,
	0x64, 0x43, 0x72, 0x79, 0x70, 0x74, 0x12, 0x4a, 0x0a, 0x12, 0x43, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e,
	0x4b, 0x65, 0x79, 0x43, 0x72, 0x79, 0x70, 0x74, 0x54, 0x79, 0x70, 0x65, 0x18, 0x06, 0x20, 0x01,
	0x28, 0x0e, 0x32, 0x1a, 0x2e, 0x6b, 0x65, 0x72, 0x6e, 0x65, 0x6c, 0x2e, 0x43, 0x6f, 0x6d, 0x6d,
	0x6f, 0x6e, 0x4b, 0x65, 0x79, 0x43, 0x72, 0x79, 0x70, 0x74, 0x54, 0x79, 0x70, 0x65, 0x52, 0x12,
	0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x4b, 0x65, 0x79, 0x43, 0x72, 0x79, 0x70, 0x74, 0x54, 0x79,
	0x70, 0x65, 0x12, 0x20, 0x0a, 0x0b, 0x43, 0x72, 0x79, 0x70, 0x74, 0x53, 0x65, 0x63, 0x72, 0x65,
	0x74, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x63, 0x72, 0x79, 0x70, 0x74, 0x53, 0x65,
	0x63, 0x72, 0x65, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x52, 0x65, 0x61, 0x6c, 0x6d, 0x18, 0x08, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x05, 0x72, 0x65, 0x61, 0x6c, 0x6d, 0x12, 0x41, 0x0a, 0x0b, 0x45, 0x6e,
	0x76, 0x50, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x1d, 0x2e, 0x61, 0x70, 0x70, 0x2e, 0x76, 0x31, 0x2e, 0x42, 0x61, 0x73, 0x69, 0x63, 0x41, 0x75,
	0x74, 0x68, 0x6e, 0x45, 0x6e, 0x76, 0x50, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72, 0x48, 0x00,
	0x52, 0x0b, 0x65, 0x6e, 0x76, 0x50, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72, 0x12, 0x44, 0x0a,
	0x0c, 0x46, 0x69, 0x6c, 0x65, 0x50, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72, 0x18, 0x0b, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x61, 0x70, 0x70, 0x2e, 0x76, 0x31, 0x2e, 0x42, 0x61, 0x73,
	0x69, 0x63, 0x41, 0x75, 0x74, 0x68, 0x6e, 0x46, 0x69, 0x6c, 0x65, 0x50, 0x72, 0x6f, 0x76, 0x69,
	0x64, 0x65, 0x72, 0x48, 0x00, 0x52, 0x0c, 0x66, 0x69, 0x6c, 0x65, 0x50, 0x72, 0x6f, 0x76, 0x69,
	0x64, 0x65, 0x72, 0x42, 0x0b, 0x0a, 0x09, 0x50, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72, 0x73,
	0x22, 0x99, 0x01, 0x0a, 0x15, 0x42, 0x61, 0x73, 0x69, 0x63, 0x41, 0x75, 0x74, 0x68, 0x6e, 0x45,
	0x6e, 0x76, 0x50, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72, 0x12, 0x26, 0x0a, 0x0e, 0x55, 0x73,
	0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x50, 0x72, 0x65, 0x66, 0x69, 0x78, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0e, 0x75, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x50, 0x72, 0x65, 0x66,
	0x69, 0x78, 0x12, 0x26, 0x0a, 0x0e, 0x50, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x50, 0x72,
	0x65, 0x66, 0x69, 0x78, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0e, 0x70, 0x61, 0x73, 0x73,
	0x77, 0x6f, 0x72, 0x64, 0x50, 0x72, 0x65, 0x66, 0x69, 0x78, 0x12, 0x30, 0x0a, 0x08, 0x45, 0x6e,
	0x63, 0x6f, 0x64, 0x69, 0x6e, 0x67, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x14, 0x2e, 0x6b,
	0x65, 0x72, 0x6e, 0x65, 0x6c, 0x2e, 0x45, 0x6e, 0x63, 0x6f, 0x64, 0x69, 0x6e, 0x67, 0x54, 0x79,
	0x70, 0x65, 0x52, 0x08, 0x65, 0x6e, 0x63, 0x6f, 0x64, 0x69, 0x6e, 0x67, 0x22, 0x60, 0x0a, 0x16,
	0x42, 0x61, 0x73, 0x69, 0x63, 0x41, 0x75, 0x74, 0x68, 0x6e, 0x46, 0x69, 0x6c, 0x65, 0x50, 0x72,
	0x6f, 0x76, 0x69, 0x64, 0x65, 0x72, 0x12, 0x14, 0x0a, 0x05, 0x50, 0x61, 0x74, 0x68, 0x73, 0x18,
	0x01, 0x20, 0x03, 0x28, 0x09, 0x52, 0x05, 0x70, 0x61, 0x74, 0x68, 0x73, 0x12, 0x30, 0x0a, 0x08,
	0x45, 0x6e, 0x63, 0x6f, 0x64, 0x69, 0x6e, 0x67, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x14,
	0x2e, 0x6b, 0x65, 0x72, 0x6e, 0x65, 0x6c, 0x2e, 0x45, 0x6e, 0x63, 0x6f, 0x64, 0x69, 0x6e, 0x67,
	0x54, 0x79, 0x70, 0x65, 0x52, 0x08, 0x65, 0x6e, 0x63, 0x6f, 0x64, 0x69, 0x6e, 0x67, 0x42, 0x38,
	0x5a, 0x36, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x61, 0x69, 0x6c,
	0x65, 0x72, 0x6f, 0x6e, 0x2d, 0x67, 0x61, 0x74, 0x65, 0x77, 0x61, 0x79, 0x2f, 0x61, 0x69, 0x6c,
	0x65, 0x72, 0x6f, 0x6e, 0x2d, 0x67, 0x61, 0x74, 0x65, 0x77, 0x61, 0x79, 0x2f, 0x61, 0x70, 0x69,
	0x73, 0x2f, 0x61, 0x70, 0x70, 0x2f, 0x76, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_app_v1_authn_basic_proto_rawDescOnce sync.Once
	file_app_v1_authn_basic_proto_rawDescData = file_app_v1_authn_basic_proto_rawDesc
)

func file_app_v1_authn_basic_proto_rawDescGZIP() []byte {
	file_app_v1_authn_basic_proto_rawDescOnce.Do(func() {
		file_app_v1_authn_basic_proto_rawDescData = protoimpl.X.CompressGZIP(file_app_v1_authn_basic_proto_rawDescData)
	})
	return file_app_v1_authn_basic_proto_rawDescData
}

var file_app_v1_authn_basic_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_app_v1_authn_basic_proto_goTypes = []any{
	(*BasicAuthnMiddleware)(nil),     // 0: app.v1.BasicAuthnMiddleware
	(*BasicAuthnMiddlewareSpec)(nil), // 1: app.v1.BasicAuthnMiddlewareSpec
	(*BasicAuthnEnvProvider)(nil),    // 2: app.v1.BasicAuthnEnvProvider
	(*BasicAuthnFileProvider)(nil),   // 3: app.v1.BasicAuthnFileProvider
	(*kernel.Metadata)(nil),          // 4: kernel.Metadata
	(*kernel.Reference)(nil),         // 5: kernel.Reference
	(*kernel.PasswordCryptSpec)(nil), // 6: kernel.PasswordCryptSpec
	(kernel.CommonKeyCryptType)(0),   // 7: kernel.CommonKeyCryptType
	(kernel.EncodingType)(0),         // 8: kernel.EncodingType
}
var file_app_v1_authn_basic_proto_depIdxs = []int32{
	4,  // 0: app.v1.BasicAuthnMiddleware.Metadata:type_name -> kernel.Metadata
	1,  // 1: app.v1.BasicAuthnMiddleware.Spec:type_name -> app.v1.BasicAuthnMiddlewareSpec
	5,  // 2: app.v1.BasicAuthnMiddlewareSpec.Logger:type_name -> kernel.Reference
	5,  // 3: app.v1.BasicAuthnMiddlewareSpec.ErrorHandler:type_name -> kernel.Reference
	6,  // 4: app.v1.BasicAuthnMiddlewareSpec.PasswordCrypt:type_name -> kernel.PasswordCryptSpec
	7,  // 5: app.v1.BasicAuthnMiddlewareSpec.CommonKeyCryptType:type_name -> kernel.CommonKeyCryptType
	2,  // 6: app.v1.BasicAuthnMiddlewareSpec.EnvProvider:type_name -> app.v1.BasicAuthnEnvProvider
	3,  // 7: app.v1.BasicAuthnMiddlewareSpec.FileProvider:type_name -> app.v1.BasicAuthnFileProvider
	8,  // 8: app.v1.BasicAuthnEnvProvider.Encoding:type_name -> kernel.EncodingType
	8,  // 9: app.v1.BasicAuthnFileProvider.Encoding:type_name -> kernel.EncodingType
	10, // [10:10] is the sub-list for method output_type
	10, // [10:10] is the sub-list for method input_type
	10, // [10:10] is the sub-list for extension type_name
	10, // [10:10] is the sub-list for extension extendee
	0,  // [0:10] is the sub-list for field type_name
}

func init() { file_app_v1_authn_basic_proto_init() }
func file_app_v1_authn_basic_proto_init() {
	if File_app_v1_authn_basic_proto != nil {
		return
	}
	file_app_v1_authn_basic_proto_msgTypes[1].OneofWrappers = []any{
		(*BasicAuthnMiddlewareSpec_EnvProvider)(nil),
		(*BasicAuthnMiddlewareSpec_FileProvider)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_app_v1_authn_basic_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_app_v1_authn_basic_proto_goTypes,
		DependencyIndexes: file_app_v1_authn_basic_proto_depIdxs,
		MessageInfos:      file_app_v1_authn_basic_proto_msgTypes,
	}.Build()
	File_app_v1_authn_basic_proto = out.File
	file_app_v1_authn_basic_proto_rawDesc = nil
	file_app_v1_authn_basic_proto_goTypes = nil
	file_app_v1_authn_basic_proto_depIdxs = nil
}
