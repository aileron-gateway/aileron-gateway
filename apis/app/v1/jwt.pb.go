// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.4
// 	protoc        v5.29.0
// source: app/v1/jwt.proto

package v1

import (
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

// SigningKeyAlgorithm is algorithm to sign JWTs.
type SigningKeyAlgorithm int32

const (
	SigningKeyAlgorithm_ALGORITHM_UNKNOWN SigningKeyAlgorithm = 0  // Unknown signing algorithm
	SigningKeyAlgorithm_NONE              SigningKeyAlgorithm = 1  // No digital signature or MAC
	SigningKeyAlgorithm_ES256             SigningKeyAlgorithm = 2  // ECDSA using P-256 and SHA-256
	SigningKeyAlgorithm_ES384             SigningKeyAlgorithm = 3  // ECDSA using P-384 and SHA-384
	SigningKeyAlgorithm_ES512             SigningKeyAlgorithm = 4  // ECDSA using P-521 and SHA-512
	SigningKeyAlgorithm_EdDSA             SigningKeyAlgorithm = 5  // EdDSA using Ed25519
	SigningKeyAlgorithm_HS256             SigningKeyAlgorithm = 6  // HMAC using SHA-256
	SigningKeyAlgorithm_HS384             SigningKeyAlgorithm = 7  // HMAC using SHA-384
	SigningKeyAlgorithm_HS512             SigningKeyAlgorithm = 8  // HMAC using SHA-512
	SigningKeyAlgorithm_RS256             SigningKeyAlgorithm = 9  // RSASSA-PKCS1-v1_5 using SHA-256
	SigningKeyAlgorithm_RS384             SigningKeyAlgorithm = 10 // RSASSA-PKCS1-v1_5 using SHA-384
	SigningKeyAlgorithm_RS512             SigningKeyAlgorithm = 11 // RSASSA-PKCS1-v1_5 using SHA-512
	SigningKeyAlgorithm_PS256             SigningKeyAlgorithm = 12 // RSASSA-PSS using SHA-256 and MGF1 with SHA-256
	SigningKeyAlgorithm_PS384             SigningKeyAlgorithm = 13 // RSASSA-PSS using SHA-384 and MGF1 with SHA-384
	SigningKeyAlgorithm_PS512             SigningKeyAlgorithm = 14 // RSASSA-PSS using SHA-512 and MGF1 with SHA-512
)

// Enum value maps for SigningKeyAlgorithm.
var (
	SigningKeyAlgorithm_name = map[int32]string{
		0:  "ALGORITHM_UNKNOWN",
		1:  "NONE",
		2:  "ES256",
		3:  "ES384",
		4:  "ES512",
		5:  "EdDSA",
		6:  "HS256",
		7:  "HS384",
		8:  "HS512",
		9:  "RS256",
		10: "RS384",
		11: "RS512",
		12: "PS256",
		13: "PS384",
		14: "PS512",
	}
	SigningKeyAlgorithm_value = map[string]int32{
		"ALGORITHM_UNKNOWN": 0,
		"NONE":              1,
		"ES256":             2,
		"ES384":             3,
		"ES512":             4,
		"EdDSA":             5,
		"HS256":             6,
		"HS384":             7,
		"HS512":             8,
		"RS256":             9,
		"RS384":             10,
		"RS512":             11,
		"PS256":             12,
		"PS384":             13,
		"PS512":             14,
	}
)

func (x SigningKeyAlgorithm) Enum() *SigningKeyAlgorithm {
	p := new(SigningKeyAlgorithm)
	*p = x
	return p
}

func (x SigningKeyAlgorithm) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (SigningKeyAlgorithm) Descriptor() protoreflect.EnumDescriptor {
	return file_app_v1_jwt_proto_enumTypes[0].Descriptor()
}

func (SigningKeyAlgorithm) Type() protoreflect.EnumType {
	return &file_app_v1_jwt_proto_enumTypes[0]
}

func (x SigningKeyAlgorithm) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use SigningKeyAlgorithm.Descriptor instead.
func (SigningKeyAlgorithm) EnumDescriptor() ([]byte, []int) {
	return file_app_v1_jwt_proto_rawDescGZIP(), []int{0}
}

// SigningKeyType is the type of key for encryption.
type SigningKeyType int32

const (
	SigningKeyType_KEY_TYPE_UNKNOWN SigningKeyType = 0 // Unknown encryption key type.
	SigningKeyType_COMMON           SigningKeyType = 1 // Common key for common key encryption.
	SigningKeyType_PUBLIC           SigningKeyType = 2 // Public keys for public key encryption.
	SigningKeyType_PRIVATE          SigningKeyType = 3 // Private keys for public key encryption.
)

// Enum value maps for SigningKeyType.
var (
	SigningKeyType_name = map[int32]string{
		0: "KEY_TYPE_UNKNOWN",
		1: "COMMON",
		2: "PUBLIC",
		3: "PRIVATE",
	}
	SigningKeyType_value = map[string]int32{
		"KEY_TYPE_UNKNOWN": 0,
		"COMMON":           1,
		"PUBLIC":           2,
		"PRIVATE":          3,
	}
)

func (x SigningKeyType) Enum() *SigningKeyType {
	p := new(SigningKeyType)
	*p = x
	return p
}

func (x SigningKeyType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (SigningKeyType) Descriptor() protoreflect.EnumDescriptor {
	return file_app_v1_jwt_proto_enumTypes[1].Descriptor()
}

func (SigningKeyType) Type() protoreflect.EnumType {
	return &file_app_v1_jwt_proto_enumTypes[1]
}

func (x SigningKeyType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use SigningKeyType.Descriptor instead.
func (SigningKeyType) EnumDescriptor() ([]byte, []int) {
	return file_app_v1_jwt_proto_rawDescGZIP(), []int{1}
}

// SigningKeySpec is the definition of the JWT signing key object.
type SigningKeySpec struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// [OPTIONAL]
	// KeyID is the key identifier.
	// This value is set in the "kid" attribute of JWT headers.
	// This value is automatically set if not set.
	KeyID string `protobuf:"bytes,1,opt,name=KeyID,json=keyID,proto3" json:"KeyID,omitempty"`
	// [REQUIRED]
	// Algorithm is the signing algorith when generating JWTs.
	Algorithm SigningKeyAlgorithm `protobuf:"varint,2,opt,name=Algorithm,json=algorithm,proto3,enum=app.v1.SigningKeyAlgorithm" json:"Algorithm,omitempty"`
	// [REQUIRED]
	// KeyType is the type of key.
	KeyType SigningKeyType `protobuf:"varint,3,opt,name=KeyType,json=keyType,proto3,enum=app.v1.SigningKeyType" json:"KeyType,omitempty"`
	// [OPTIONAL]
	// KeyFilePath is the file path to a common key or a pem key.
	// KeyFilePath is used when both keyFilePath and keyString are set.
	KeyFilePath string `protobuf:"bytes,4,opt,name=KeyFilePath,json=keyFilePath,proto3" json:"KeyFilePath,omitempty"`
	// [OPTIONAL]
	// KeyString is the base64 encoded string of a common key or a pem key.
	// KeyFilePath is used when both keyFilePath and keyString are set.
	KeyString string `protobuf:"bytes,5,opt,name=KeyString,json=keyString,proto3" json:"KeyString,omitempty"`
	// [OPTIONAL]
	// JWTHeader is the user defined header values in the JWT's header.
	JWTHeader     map[string]string `protobuf:"bytes,6,rep,name=JWTHeader,json=jwtHeader,proto3" json:"JWTHeader,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *SigningKeySpec) Reset() {
	*x = SigningKeySpec{}
	mi := &file_app_v1_jwt_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SigningKeySpec) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SigningKeySpec) ProtoMessage() {}

func (x *SigningKeySpec) ProtoReflect() protoreflect.Message {
	mi := &file_app_v1_jwt_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SigningKeySpec.ProtoReflect.Descriptor instead.
func (*SigningKeySpec) Descriptor() ([]byte, []int) {
	return file_app_v1_jwt_proto_rawDescGZIP(), []int{0}
}

func (x *SigningKeySpec) GetKeyID() string {
	if x != nil {
		return x.KeyID
	}
	return ""
}

func (x *SigningKeySpec) GetAlgorithm() SigningKeyAlgorithm {
	if x != nil {
		return x.Algorithm
	}
	return SigningKeyAlgorithm_ALGORITHM_UNKNOWN
}

func (x *SigningKeySpec) GetKeyType() SigningKeyType {
	if x != nil {
		return x.KeyType
	}
	return SigningKeyType_KEY_TYPE_UNKNOWN
}

func (x *SigningKeySpec) GetKeyFilePath() string {
	if x != nil {
		return x.KeyFilePath
	}
	return ""
}

func (x *SigningKeySpec) GetKeyString() string {
	if x != nil {
		return x.KeyString
	}
	return ""
}

func (x *SigningKeySpec) GetJWTHeader() map[string]string {
	if x != nil {
		return x.JWTHeader
	}
	return nil
}

// JWTHandlerSpec is the specification of JWTHandler object.
type JWTHandlerSpec struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// [OPTIONAL]
	// PrivateKeys is list of private key objects for signing JWTs.
	// Default is empty.
	PrivateKeys []*SigningKeySpec `protobuf:"bytes,1,rep,name=PrivateKeys,json=privateKeys,proto3" json:"PrivateKeys,omitempty"`
	// [OPTIONAL]
	// PublicKeys is list of public key objects for validating JWTs.
	// Default is empty.
	PublicKeys []*SigningKeySpec `protobuf:"bytes,2,rep,name=PublicKeys,json=publicKeys,proto3" json:"PublicKeys,omitempty"`
	// [OPTIONAL]
	// JWKs is the pairs of issuer and JWKs URLs.
	// The keys are issuer and the values are JWKs URLs corresponding to the issuer.
	// JWKs URLs are used when a validating key was not found for JWTs.
	JWKs map[string]string `protobuf:"bytes,3,rep,name=JWKs,proto3" json:"JWKs,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
	// [OPTIONAL]
	// UseJKU is the flag to use JWKs endpoint set in "jku" header field.
	// JWKs endpoint in "jku" header is used if there is no key cached
	// in the handler for validating a token.
	// Default is [false].
	UseJKU        bool `protobuf:"varint,4,opt,name=UseJKU,json=useJKU,proto3" json:"UseJKU,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *JWTHandlerSpec) Reset() {
	*x = JWTHandlerSpec{}
	mi := &file_app_v1_jwt_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *JWTHandlerSpec) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*JWTHandlerSpec) ProtoMessage() {}

func (x *JWTHandlerSpec) ProtoReflect() protoreflect.Message {
	mi := &file_app_v1_jwt_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use JWTHandlerSpec.ProtoReflect.Descriptor instead.
func (*JWTHandlerSpec) Descriptor() ([]byte, []int) {
	return file_app_v1_jwt_proto_rawDescGZIP(), []int{1}
}

func (x *JWTHandlerSpec) GetPrivateKeys() []*SigningKeySpec {
	if x != nil {
		return x.PrivateKeys
	}
	return nil
}

func (x *JWTHandlerSpec) GetPublicKeys() []*SigningKeySpec {
	if x != nil {
		return x.PublicKeys
	}
	return nil
}

func (x *JWTHandlerSpec) GetJWKs() map[string]string {
	if x != nil {
		return x.JWKs
	}
	return nil
}

func (x *JWTHandlerSpec) GetUseJKU() bool {
	if x != nil {
		return x.UseJKU
	}
	return false
}

var File_app_v1_jwt_proto protoreflect.FileDescriptor

var file_app_v1_jwt_proto_rawDesc = string([]byte{
	0x0a, 0x10, 0x61, 0x70, 0x70, 0x2f, 0x76, 0x31, 0x2f, 0x6a, 0x77, 0x74, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x06, 0x61, 0x70, 0x70, 0x2e, 0x76, 0x31, 0x22, 0xd6, 0x02, 0x0a, 0x0e, 0x53,
	0x69, 0x67, 0x6e, 0x69, 0x6e, 0x67, 0x4b, 0x65, 0x79, 0x53, 0x70, 0x65, 0x63, 0x12, 0x14, 0x0a,
	0x05, 0x4b, 0x65, 0x79, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x6b, 0x65,
	0x79, 0x49, 0x44, 0x12, 0x39, 0x0a, 0x09, 0x41, 0x6c, 0x67, 0x6f, 0x72, 0x69, 0x74, 0x68, 0x6d,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x1b, 0x2e, 0x61, 0x70, 0x70, 0x2e, 0x76, 0x31, 0x2e,
	0x53, 0x69, 0x67, 0x6e, 0x69, 0x6e, 0x67, 0x4b, 0x65, 0x79, 0x41, 0x6c, 0x67, 0x6f, 0x72, 0x69,
	0x74, 0x68, 0x6d, 0x52, 0x09, 0x61, 0x6c, 0x67, 0x6f, 0x72, 0x69, 0x74, 0x68, 0x6d, 0x12, 0x30,
	0x0a, 0x07, 0x4b, 0x65, 0x79, 0x54, 0x79, 0x70, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0e, 0x32,
	0x16, 0x2e, 0x61, 0x70, 0x70, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x69, 0x67, 0x6e, 0x69, 0x6e, 0x67,
	0x4b, 0x65, 0x79, 0x54, 0x79, 0x70, 0x65, 0x52, 0x07, 0x6b, 0x65, 0x79, 0x54, 0x79, 0x70, 0x65,
	0x12, 0x20, 0x0a, 0x0b, 0x4b, 0x65, 0x79, 0x46, 0x69, 0x6c, 0x65, 0x50, 0x61, 0x74, 0x68, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x6b, 0x65, 0x79, 0x46, 0x69, 0x6c, 0x65, 0x50, 0x61,
	0x74, 0x68, 0x12, 0x1c, 0x0a, 0x09, 0x4b, 0x65, 0x79, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x18,
	0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x6b, 0x65, 0x79, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67,
	0x12, 0x43, 0x0a, 0x09, 0x4a, 0x57, 0x54, 0x48, 0x65, 0x61, 0x64, 0x65, 0x72, 0x18, 0x06, 0x20,
	0x03, 0x28, 0x0b, 0x32, 0x25, 0x2e, 0x61, 0x70, 0x70, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x69, 0x67,
	0x6e, 0x69, 0x6e, 0x67, 0x4b, 0x65, 0x79, 0x53, 0x70, 0x65, 0x63, 0x2e, 0x4a, 0x57, 0x54, 0x48,
	0x65, 0x61, 0x64, 0x65, 0x72, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x09, 0x6a, 0x77, 0x74, 0x48,
	0x65, 0x61, 0x64, 0x65, 0x72, 0x1a, 0x3c, 0x0a, 0x0e, 0x4a, 0x57, 0x54, 0x48, 0x65, 0x61, 0x64,
	0x65, 0x72, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c,
	0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a,
	0x02, 0x38, 0x01, 0x22, 0x89, 0x02, 0x0a, 0x0e, 0x4a, 0x57, 0x54, 0x48, 0x61, 0x6e, 0x64, 0x6c,
	0x65, 0x72, 0x53, 0x70, 0x65, 0x63, 0x12, 0x38, 0x0a, 0x0b, 0x50, 0x72, 0x69, 0x76, 0x61, 0x74,
	0x65, 0x4b, 0x65, 0x79, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x16, 0x2e, 0x61, 0x70,
	0x70, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x69, 0x67, 0x6e, 0x69, 0x6e, 0x67, 0x4b, 0x65, 0x79, 0x53,
	0x70, 0x65, 0x63, 0x52, 0x0b, 0x70, 0x72, 0x69, 0x76, 0x61, 0x74, 0x65, 0x4b, 0x65, 0x79, 0x73,
	0x12, 0x36, 0x0a, 0x0a, 0x50, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x4b, 0x65, 0x79, 0x73, 0x18, 0x02,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x16, 0x2e, 0x61, 0x70, 0x70, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x69,
	0x67, 0x6e, 0x69, 0x6e, 0x67, 0x4b, 0x65, 0x79, 0x53, 0x70, 0x65, 0x63, 0x52, 0x0a, 0x70, 0x75,
	0x62, 0x6c, 0x69, 0x63, 0x4b, 0x65, 0x79, 0x73, 0x12, 0x34, 0x0a, 0x04, 0x4a, 0x57, 0x4b, 0x73,
	0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x20, 0x2e, 0x61, 0x70, 0x70, 0x2e, 0x76, 0x31, 0x2e,
	0x4a, 0x57, 0x54, 0x48, 0x61, 0x6e, 0x64, 0x6c, 0x65, 0x72, 0x53, 0x70, 0x65, 0x63, 0x2e, 0x4a,
	0x57, 0x4b, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x04, 0x4a, 0x57, 0x4b, 0x73, 0x12, 0x16,
	0x0a, 0x06, 0x55, 0x73, 0x65, 0x4a, 0x4b, 0x55, 0x18, 0x04, 0x20, 0x01, 0x28, 0x08, 0x52, 0x06,
	0x75, 0x73, 0x65, 0x4a, 0x4b, 0x55, 0x1a, 0x37, 0x0a, 0x09, 0x4a, 0x57, 0x4b, 0x73, 0x45, 0x6e,
	0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x2a,
	0xc5, 0x01, 0x0a, 0x13, 0x53, 0x69, 0x67, 0x6e, 0x69, 0x6e, 0x67, 0x4b, 0x65, 0x79, 0x41, 0x6c,
	0x67, 0x6f, 0x72, 0x69, 0x74, 0x68, 0x6d, 0x12, 0x15, 0x0a, 0x11, 0x41, 0x4c, 0x47, 0x4f, 0x52,
	0x49, 0x54, 0x48, 0x4d, 0x5f, 0x55, 0x4e, 0x4b, 0x4e, 0x4f, 0x57, 0x4e, 0x10, 0x00, 0x12, 0x08,
	0x0a, 0x04, 0x4e, 0x4f, 0x4e, 0x45, 0x10, 0x01, 0x12, 0x09, 0x0a, 0x05, 0x45, 0x53, 0x32, 0x35,
	0x36, 0x10, 0x02, 0x12, 0x09, 0x0a, 0x05, 0x45, 0x53, 0x33, 0x38, 0x34, 0x10, 0x03, 0x12, 0x09,
	0x0a, 0x05, 0x45, 0x53, 0x35, 0x31, 0x32, 0x10, 0x04, 0x12, 0x09, 0x0a, 0x05, 0x45, 0x64, 0x44,
	0x53, 0x41, 0x10, 0x05, 0x12, 0x09, 0x0a, 0x05, 0x48, 0x53, 0x32, 0x35, 0x36, 0x10, 0x06, 0x12,
	0x09, 0x0a, 0x05, 0x48, 0x53, 0x33, 0x38, 0x34, 0x10, 0x07, 0x12, 0x09, 0x0a, 0x05, 0x48, 0x53,
	0x35, 0x31, 0x32, 0x10, 0x08, 0x12, 0x09, 0x0a, 0x05, 0x52, 0x53, 0x32, 0x35, 0x36, 0x10, 0x09,
	0x12, 0x09, 0x0a, 0x05, 0x52, 0x53, 0x33, 0x38, 0x34, 0x10, 0x0a, 0x12, 0x09, 0x0a, 0x05, 0x52,
	0x53, 0x35, 0x31, 0x32, 0x10, 0x0b, 0x12, 0x09, 0x0a, 0x05, 0x50, 0x53, 0x32, 0x35, 0x36, 0x10,
	0x0c, 0x12, 0x09, 0x0a, 0x05, 0x50, 0x53, 0x33, 0x38, 0x34, 0x10, 0x0d, 0x12, 0x09, 0x0a, 0x05,
	0x50, 0x53, 0x35, 0x31, 0x32, 0x10, 0x0e, 0x2a, 0x4b, 0x0a, 0x0e, 0x53, 0x69, 0x67, 0x6e, 0x69,
	0x6e, 0x67, 0x4b, 0x65, 0x79, 0x54, 0x79, 0x70, 0x65, 0x12, 0x14, 0x0a, 0x10, 0x4b, 0x45, 0x59,
	0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x55, 0x4e, 0x4b, 0x4e, 0x4f, 0x57, 0x4e, 0x10, 0x00, 0x12,
	0x0a, 0x0a, 0x06, 0x43, 0x4f, 0x4d, 0x4d, 0x4f, 0x4e, 0x10, 0x01, 0x12, 0x0a, 0x0a, 0x06, 0x50,
	0x55, 0x42, 0x4c, 0x49, 0x43, 0x10, 0x02, 0x12, 0x0b, 0x0a, 0x07, 0x50, 0x52, 0x49, 0x56, 0x41,
	0x54, 0x45, 0x10, 0x03, 0x42, 0x38, 0x5a, 0x36, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63,
	0x6f, 0x6d, 0x2f, 0x61, 0x69, 0x6c, 0x65, 0x72, 0x6f, 0x6e, 0x2d, 0x67, 0x61, 0x74, 0x65, 0x77,
	0x61, 0x79, 0x2f, 0x61, 0x69, 0x6c, 0x65, 0x72, 0x6f, 0x6e, 0x2d, 0x67, 0x61, 0x74, 0x65, 0x77,
	0x61, 0x79, 0x2f, 0x61, 0x70, 0x69, 0x73, 0x2f, 0x61, 0x70, 0x70, 0x2f, 0x76, 0x31, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
})

var (
	file_app_v1_jwt_proto_rawDescOnce sync.Once
	file_app_v1_jwt_proto_rawDescData []byte
)

func file_app_v1_jwt_proto_rawDescGZIP() []byte {
	file_app_v1_jwt_proto_rawDescOnce.Do(func() {
		file_app_v1_jwt_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_app_v1_jwt_proto_rawDesc), len(file_app_v1_jwt_proto_rawDesc)))
	})
	return file_app_v1_jwt_proto_rawDescData
}

var file_app_v1_jwt_proto_enumTypes = make([]protoimpl.EnumInfo, 2)
var file_app_v1_jwt_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_app_v1_jwt_proto_goTypes = []any{
	(SigningKeyAlgorithm)(0), // 0: app.v1.SigningKeyAlgorithm
	(SigningKeyType)(0),      // 1: app.v1.SigningKeyType
	(*SigningKeySpec)(nil),   // 2: app.v1.SigningKeySpec
	(*JWTHandlerSpec)(nil),   // 3: app.v1.JWTHandlerSpec
	nil,                      // 4: app.v1.SigningKeySpec.JWTHeaderEntry
	nil,                      // 5: app.v1.JWTHandlerSpec.JWKsEntry
}
var file_app_v1_jwt_proto_depIdxs = []int32{
	0, // 0: app.v1.SigningKeySpec.Algorithm:type_name -> app.v1.SigningKeyAlgorithm
	1, // 1: app.v1.SigningKeySpec.KeyType:type_name -> app.v1.SigningKeyType
	4, // 2: app.v1.SigningKeySpec.JWTHeader:type_name -> app.v1.SigningKeySpec.JWTHeaderEntry
	2, // 3: app.v1.JWTHandlerSpec.PrivateKeys:type_name -> app.v1.SigningKeySpec
	2, // 4: app.v1.JWTHandlerSpec.PublicKeys:type_name -> app.v1.SigningKeySpec
	5, // 5: app.v1.JWTHandlerSpec.JWKs:type_name -> app.v1.JWTHandlerSpec.JWKsEntry
	6, // [6:6] is the sub-list for method output_type
	6, // [6:6] is the sub-list for method input_type
	6, // [6:6] is the sub-list for extension type_name
	6, // [6:6] is the sub-list for extension extendee
	0, // [0:6] is the sub-list for field type_name
}

func init() { file_app_v1_jwt_proto_init() }
func file_app_v1_jwt_proto_init() {
	if File_app_v1_jwt_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_app_v1_jwt_proto_rawDesc), len(file_app_v1_jwt_proto_rawDesc)),
			NumEnums:      2,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_app_v1_jwt_proto_goTypes,
		DependencyIndexes: file_app_v1_jwt_proto_depIdxs,
		EnumInfos:         file_app_v1_jwt_proto_enumTypes,
		MessageInfos:      file_app_v1_jwt_proto_msgTypes,
	}.Build()
	File_app_v1_jwt_proto = out.File
	file_app_v1_jwt_proto_goTypes = nil
	file_app_v1_jwt_proto_depIdxs = nil
}
