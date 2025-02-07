// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.5
// 	protoc        v5.29.0
// source: kernel/crypt.proto

package kernel

import (
	_ "buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go/buf/validate"
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

// CommonKeyCryptType is the algorithms of common key encryption.
type CommonKeyCryptType int32

const (
	CommonKeyCryptType_CommonKeyCryptTypeUnknown CommonKeyCryptType = 0  // Unknown crypt type.
	CommonKeyCryptType_AESGCM                    CommonKeyCryptType = 1  // AES-GCM chipher. Key length must be 16, 24, 32 bytes for AES-128, AES-192, AES-256.
	CommonKeyCryptType_AESCBC                    CommonKeyCryptType = 2  // AES-CBC chipher. Key length must be 16, 24, 32 bytes for AES-128, AES-192, AES-256.
	CommonKeyCryptType_AESCFB                    CommonKeyCryptType = 3  // AES-CFB chipher. Key length must be 16, 24, 32 bytes for AES-128, AES-192, AES-256.
	CommonKeyCryptType_AESCTR                    CommonKeyCryptType = 4  // AES-CTR chipher. Key length must be 16, 24, 32 bytes for AES-128, AES-192, AES-256.
	CommonKeyCryptType_AESOFB                    CommonKeyCryptType = 5  // AES-OFB chipher. Key length must be 16, 24, 32 bytes for AES-128, AES-192, AES-256.
	CommonKeyCryptType_DESCBC                    CommonKeyCryptType = 6  // DES-CBC chipher. Key length must be 8 bytes.
	CommonKeyCryptType_DESCFB                    CommonKeyCryptType = 7  // DES-CFB chipher. Key length must be 8 bytes.
	CommonKeyCryptType_DESCTR                    CommonKeyCryptType = 8  // DES-CTR chipher. Key length must be 8 bytes.
	CommonKeyCryptType_DESOFB                    CommonKeyCryptType = 9  // DES-OFB chipher. Key length must be 8 bytes.
	CommonKeyCryptType_TripleDESCBC              CommonKeyCryptType = 10 // 3DES-CBC chipher. Key length must be 24 bytes.
	CommonKeyCryptType_TripleDESCFB              CommonKeyCryptType = 11 // 3DES-CFB chipher. Key length must be 24 bytes.
	CommonKeyCryptType_TripleDESCTR              CommonKeyCryptType = 12 // 3DES-CTR chipher. Key length must be 24 bytes.
	CommonKeyCryptType_TripleDESOFB              CommonKeyCryptType = 13 // 3DES-OFB chipher. Key length must be 24 bytes.
	CommonKeyCryptType_RC4                       CommonKeyCryptType = 14 // RC4 chipher. Key length must be 5 to 256 bytes.
)

// Enum value maps for CommonKeyCryptType.
var (
	CommonKeyCryptType_name = map[int32]string{
		0:  "CommonKeyCryptTypeUnknown",
		1:  "AESGCM",
		2:  "AESCBC",
		3:  "AESCFB",
		4:  "AESCTR",
		5:  "AESOFB",
		6:  "DESCBC",
		7:  "DESCFB",
		8:  "DESCTR",
		9:  "DESOFB",
		10: "TripleDESCBC",
		11: "TripleDESCFB",
		12: "TripleDESCTR",
		13: "TripleDESOFB",
		14: "RC4",
	}
	CommonKeyCryptType_value = map[string]int32{
		"CommonKeyCryptTypeUnknown": 0,
		"AESGCM":                    1,
		"AESCBC":                    2,
		"AESCFB":                    3,
		"AESCTR":                    4,
		"AESOFB":                    5,
		"DESCBC":                    6,
		"DESCFB":                    7,
		"DESCTR":                    8,
		"DESOFB":                    9,
		"TripleDESCBC":              10,
		"TripleDESCFB":              11,
		"TripleDESCTR":              12,
		"TripleDESOFB":              13,
		"RC4":                       14,
	}
)

func (x CommonKeyCryptType) Enum() *CommonKeyCryptType {
	p := new(CommonKeyCryptType)
	*p = x
	return p
}

func (x CommonKeyCryptType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (CommonKeyCryptType) Descriptor() protoreflect.EnumDescriptor {
	return file_kernel_crypt_proto_enumTypes[0].Descriptor()
}

func (CommonKeyCryptType) Type() protoreflect.EnumType {
	return &file_kernel_crypt_proto_enumTypes[0]
}

func (x CommonKeyCryptType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use CommonKeyCryptType.Descriptor instead.
func (CommonKeyCryptType) EnumDescriptor() ([]byte, []int) {
	return file_kernel_crypt_proto_rawDescGZIP(), []int{0}
}

type PasswordCryptSpec struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Types that are valid to be assigned to PasswordCrypts:
	//
	//	*PasswordCryptSpec_BCrypt
	//	*PasswordCryptSpec_SCrypt
	//	*PasswordCryptSpec_PBKDF2
	//	*PasswordCryptSpec_Argon2I
	//	*PasswordCryptSpec_Argon2Id
	PasswordCrypts isPasswordCryptSpec_PasswordCrypts `protobuf_oneof:"PasswordCrypts"`
	unknownFields  protoimpl.UnknownFields
	sizeCache      protoimpl.SizeCache
}

func (x *PasswordCryptSpec) Reset() {
	*x = PasswordCryptSpec{}
	mi := &file_kernel_crypt_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *PasswordCryptSpec) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PasswordCryptSpec) ProtoMessage() {}

func (x *PasswordCryptSpec) ProtoReflect() protoreflect.Message {
	mi := &file_kernel_crypt_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PasswordCryptSpec.ProtoReflect.Descriptor instead.
func (*PasswordCryptSpec) Descriptor() ([]byte, []int) {
	return file_kernel_crypt_proto_rawDescGZIP(), []int{0}
}

func (x *PasswordCryptSpec) GetPasswordCrypts() isPasswordCryptSpec_PasswordCrypts {
	if x != nil {
		return x.PasswordCrypts
	}
	return nil
}

func (x *PasswordCryptSpec) GetBCrypt() *BCryptSpec {
	if x != nil {
		if x, ok := x.PasswordCrypts.(*PasswordCryptSpec_BCrypt); ok {
			return x.BCrypt
		}
	}
	return nil
}

func (x *PasswordCryptSpec) GetSCrypt() *SCryptSpec {
	if x != nil {
		if x, ok := x.PasswordCrypts.(*PasswordCryptSpec_SCrypt); ok {
			return x.SCrypt
		}
	}
	return nil
}

func (x *PasswordCryptSpec) GetPBKDF2() *PBKDF2Spec {
	if x != nil {
		if x, ok := x.PasswordCrypts.(*PasswordCryptSpec_PBKDF2); ok {
			return x.PBKDF2
		}
	}
	return nil
}

func (x *PasswordCryptSpec) GetArgon2I() *Argon2Spec {
	if x != nil {
		if x, ok := x.PasswordCrypts.(*PasswordCryptSpec_Argon2I); ok {
			return x.Argon2I
		}
	}
	return nil
}

func (x *PasswordCryptSpec) GetArgon2Id() *Argon2Spec {
	if x != nil {
		if x, ok := x.PasswordCrypts.(*PasswordCryptSpec_Argon2Id); ok {
			return x.Argon2Id
		}
	}
	return nil
}

type isPasswordCryptSpec_PasswordCrypts interface {
	isPasswordCryptSpec_PasswordCrypts()
}

type PasswordCryptSpec_BCrypt struct {
	BCrypt *BCryptSpec `protobuf:"bytes,1,opt,name=BCrypt,json=bcrypt,proto3,oneof"`
}

type PasswordCryptSpec_SCrypt struct {
	SCrypt *SCryptSpec `protobuf:"bytes,2,opt,name=SCrypt,json=scrypt,proto3,oneof"`
}

type PasswordCryptSpec_PBKDF2 struct {
	PBKDF2 *PBKDF2Spec `protobuf:"bytes,3,opt,name=PBKDF2,json=pbkdf2,proto3,oneof"`
}

type PasswordCryptSpec_Argon2I struct {
	Argon2I *Argon2Spec `protobuf:"bytes,4,opt,name=Argon2i,json=argon2i,proto3,oneof"`
}

type PasswordCryptSpec_Argon2Id struct {
	Argon2Id *Argon2Spec `protobuf:"bytes,5,opt,name=Argon2id,json=argon2id,proto3,oneof"`
}

func (*PasswordCryptSpec_BCrypt) isPasswordCryptSpec_PasswordCrypts() {}

func (*PasswordCryptSpec_SCrypt) isPasswordCryptSpec_PasswordCrypts() {}

func (*PasswordCryptSpec_PBKDF2) isPasswordCryptSpec_PasswordCrypts() {}

func (*PasswordCryptSpec_Argon2I) isPasswordCryptSpec_PasswordCrypts() {}

func (*PasswordCryptSpec_Argon2Id) isPasswordCryptSpec_PasswordCrypts() {}

type BCryptSpec struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// [OPTIONAL]
	// Cost is the cost of BCrypt hashing.
	// Cost is automatically limited to 4<=cost<=32.
	// Default value is used if not set or set to zero.
	// Default is [10].
	Cost          int32 `protobuf:"varint,1,opt,name=Cost,json=cost,proto3" json:"Cost,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *BCryptSpec) Reset() {
	*x = BCryptSpec{}
	mi := &file_kernel_crypt_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *BCryptSpec) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BCryptSpec) ProtoMessage() {}

func (x *BCryptSpec) ProtoReflect() protoreflect.Message {
	mi := &file_kernel_crypt_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BCryptSpec.ProtoReflect.Descriptor instead.
func (*BCryptSpec) Descriptor() ([]byte, []int) {
	return file_kernel_crypt_proto_rawDescGZIP(), []int{1}
}

func (x *BCryptSpec) GetCost() int32 {
	if x != nil {
		return x.Cost
	}
	return 0
}

type SCryptSpec struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// [OPTIONAL]
	// SaltLen is the salt length in bytes.
	// Random bytes read by a random reader of crypt/rand is
	// used for generating a specified length of random salt.
	// Salts are appended to the resulting hash value.
	// Default value is used if not set or set to zeto.
	// Default is [32].
	SaltLen int32 `protobuf:"varint,1,opt,name=SaltLen,json=saltLen,proto3" json:"SaltLen,omitempty"`
	// [OPTIONAL]
	// N is the "N" parameter for SCrypt algorith.
	// Default value is used if not set or set to zeto.
	// Default is [32768].
	N int32 `protobuf:"varint,2,opt,name=N,json=n,proto3" json:"N,omitempty"`
	// [OPTIONAL]
	// R is the "r" parameter for SCrypt algorith.
	// Default value is used if not set or set to zeto.
	// Default is [8].
	R int32 `protobuf:"varint,3,opt,name=R,json=r,proto3" json:"R,omitempty"`
	// [OPTIONAL]
	// P is the "p" parameter for SCrypt algorith.
	// Default value is used if not set or set to zeto.
	// Default is [1].
	P int32 `protobuf:"varint,4,opt,name=P,json=p,proto3" json:"P,omitempty"`
	// [OPTIONAL]
	// KeyLen is the hashed key length.
	// Default value is used if not set or set to zeto.
	// Default is [32].
	KeyLen        int32 `protobuf:"varint,5,opt,name=KeyLen,json=keyLen,proto3" json:"KeyLen,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *SCryptSpec) Reset() {
	*x = SCryptSpec{}
	mi := &file_kernel_crypt_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SCryptSpec) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SCryptSpec) ProtoMessage() {}

func (x *SCryptSpec) ProtoReflect() protoreflect.Message {
	mi := &file_kernel_crypt_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SCryptSpec.ProtoReflect.Descriptor instead.
func (*SCryptSpec) Descriptor() ([]byte, []int) {
	return file_kernel_crypt_proto_rawDescGZIP(), []int{2}
}

func (x *SCryptSpec) GetSaltLen() int32 {
	if x != nil {
		return x.SaltLen
	}
	return 0
}

func (x *SCryptSpec) GetN() int32 {
	if x != nil {
		return x.N
	}
	return 0
}

func (x *SCryptSpec) GetR() int32 {
	if x != nil {
		return x.R
	}
	return 0
}

func (x *SCryptSpec) GetP() int32 {
	if x != nil {
		return x.P
	}
	return 0
}

func (x *SCryptSpec) GetKeyLen() int32 {
	if x != nil {
		return x.KeyLen
	}
	return 0
}

type PBKDF2Spec struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// [OPTIONAL]
	// SaltLen is the salt length in bytes.
	// Random bytes read by a random reader of crypt/rand is
	// used for generating a specified length of random salt.
	// Salts are appended to the resulting hash value.
	// Default value is used if not set or set to zeto.
	// Default is [32].
	SaltLen int32 `protobuf:"varint,1,opt,name=SaltLen,json=saltLen,proto3" json:"SaltLen,omitempty"`
	// [OPTIONAL]
	// Iter is the iteration count parameter for PBKDF2.
	// Default value is used if not set or set to zeto.
	// Default is [4096].
	Iter int32 `protobuf:"varint,2,opt,name=Iter,json=iter,proto3" json:"Iter,omitempty"`
	// [OPTIONAL]
	// KeyLen is the hashed key length.
	// Default value is used if not set or set to zeto.
	// Default is [32].
	KeyLen int32 `protobuf:"varint,3,opt,name=KeyLen,json=keyLen,proto3" json:"KeyLen,omitempty"`
	// Currently following algorithms are available.
	// SHA1, SHA224, SHA256, SHA384, SHA512, SHA512_224,
	// SHA512_256, SHA3_224, SHA3_256, SHA3_384, SHA3_512, MD5.
	HashAlg       HashAlg `protobuf:"varint,4,opt,name=HashAlg,json=hashAlg,proto3,enum=kernel.HashAlg" json:"HashAlg,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *PBKDF2Spec) Reset() {
	*x = PBKDF2Spec{}
	mi := &file_kernel_crypt_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *PBKDF2Spec) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PBKDF2Spec) ProtoMessage() {}

func (x *PBKDF2Spec) ProtoReflect() protoreflect.Message {
	mi := &file_kernel_crypt_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PBKDF2Spec.ProtoReflect.Descriptor instead.
func (*PBKDF2Spec) Descriptor() ([]byte, []int) {
	return file_kernel_crypt_proto_rawDescGZIP(), []int{3}
}

func (x *PBKDF2Spec) GetSaltLen() int32 {
	if x != nil {
		return x.SaltLen
	}
	return 0
}

func (x *PBKDF2Spec) GetIter() int32 {
	if x != nil {
		return x.Iter
	}
	return 0
}

func (x *PBKDF2Spec) GetKeyLen() int32 {
	if x != nil {
		return x.KeyLen
	}
	return 0
}

func (x *PBKDF2Spec) GetHashAlg() HashAlg {
	if x != nil {
		return x.HashAlg
	}
	return HashAlg_HashAlgUnknown
}

type Argon2Spec struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// [OPTIONAL]
	// SaltLen is the salt length in bytes.
	// Random bytes read by a random reader of crypt/rand is
	// used for generating a specified length of random salt.
	// Salts are appended to the resulting hash value.
	// Default value is used if not set or set to zeto.
	// Default is [32].
	SaltLen uint32 `protobuf:"varint,1,opt,name=SaltLen,json=saltLen,proto3" json:"SaltLen,omitempty"`
	// [OPTIONAL]
	// Time is the time parameter for Argon2i and Argon2id.
	// Default is [3] for Argon2i and [1] for Argon2id.
	Time uint32 `protobuf:"varint,2,opt,name=Time,json=time,proto3" json:"Time,omitempty"`
	// [OPTIONAL]
	// Memory is the memory parameter for Argon2i and Argon2id.
	// Default value is used if not set or set to zero.
	// Default is [32*1024] for Argon2i and [64*1024] for Argon2id.
	Memory uint32 `protobuf:"varint,3,opt,name=Memory,json=memory,proto3" json:"Memory,omitempty"`
	// [OPTIONAL]
	// Threads is the thread number to use for hash calculation.
	// Default value is used if not set or set to zero.
	// Default is [4].
	Threads uint32 `protobuf:"varint,4,opt,name=Threads,json=threads,proto3" json:"Threads,omitempty"`
	// [OPTIONAL]
	// KeyLen is the hashed key length.
	// Default value is used if not set or set to zeto.
	// Default is [32].
	KeyLen        uint32 `protobuf:"varint,5,opt,name=KeyLen,json=keyLen,proto3" json:"KeyLen,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Argon2Spec) Reset() {
	*x = Argon2Spec{}
	mi := &file_kernel_crypt_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Argon2Spec) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Argon2Spec) ProtoMessage() {}

func (x *Argon2Spec) ProtoReflect() protoreflect.Message {
	mi := &file_kernel_crypt_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Argon2Spec.ProtoReflect.Descriptor instead.
func (*Argon2Spec) Descriptor() ([]byte, []int) {
	return file_kernel_crypt_proto_rawDescGZIP(), []int{4}
}

func (x *Argon2Spec) GetSaltLen() uint32 {
	if x != nil {
		return x.SaltLen
	}
	return 0
}

func (x *Argon2Spec) GetTime() uint32 {
	if x != nil {
		return x.Time
	}
	return 0
}

func (x *Argon2Spec) GetMemory() uint32 {
	if x != nil {
		return x.Memory
	}
	return 0
}

func (x *Argon2Spec) GetThreads() uint32 {
	if x != nil {
		return x.Threads
	}
	return 0
}

func (x *Argon2Spec) GetKeyLen() uint32 {
	if x != nil {
		return x.KeyLen
	}
	return 0
}

var File_kernel_crypt_proto protoreflect.FileDescriptor

var file_kernel_crypt_proto_rawDesc = string([]byte{
	0x0a, 0x12, 0x6b, 0x65, 0x72, 0x6e, 0x65, 0x6c, 0x2f, 0x63, 0x72, 0x79, 0x70, 0x74, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x12, 0x06, 0x6b, 0x65, 0x72, 0x6e, 0x65, 0x6c, 0x1a, 0x1b, 0x62, 0x75,
	0x66, 0x2f, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2f, 0x76, 0x61, 0x6c, 0x69, 0x64,
	0x61, 0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x11, 0x6b, 0x65, 0x72, 0x6e, 0x65,
	0x6c, 0x2f, 0x68, 0x61, 0x73, 0x68, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x91, 0x02, 0x0a,
	0x11, 0x50, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x43, 0x72, 0x79, 0x70, 0x74, 0x53, 0x70,
	0x65, 0x63, 0x12, 0x2c, 0x0a, 0x06, 0x42, 0x43, 0x72, 0x79, 0x70, 0x74, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x12, 0x2e, 0x6b, 0x65, 0x72, 0x6e, 0x65, 0x6c, 0x2e, 0x42, 0x43, 0x72, 0x79,
	0x70, 0x74, 0x53, 0x70, 0x65, 0x63, 0x48, 0x00, 0x52, 0x06, 0x62, 0x63, 0x72, 0x79, 0x70, 0x74,
	0x12, 0x2c, 0x0a, 0x06, 0x53, 0x43, 0x72, 0x79, 0x70, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x12, 0x2e, 0x6b, 0x65, 0x72, 0x6e, 0x65, 0x6c, 0x2e, 0x53, 0x43, 0x72, 0x79, 0x70, 0x74,
	0x53, 0x70, 0x65, 0x63, 0x48, 0x00, 0x52, 0x06, 0x73, 0x63, 0x72, 0x79, 0x70, 0x74, 0x12, 0x2c,
	0x0a, 0x06, 0x50, 0x42, 0x4b, 0x44, 0x46, 0x32, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x12,
	0x2e, 0x6b, 0x65, 0x72, 0x6e, 0x65, 0x6c, 0x2e, 0x50, 0x42, 0x4b, 0x44, 0x46, 0x32, 0x53, 0x70,
	0x65, 0x63, 0x48, 0x00, 0x52, 0x06, 0x70, 0x62, 0x6b, 0x64, 0x66, 0x32, 0x12, 0x2e, 0x0a, 0x07,
	0x41, 0x72, 0x67, 0x6f, 0x6e, 0x32, 0x69, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x12, 0x2e,
	0x6b, 0x65, 0x72, 0x6e, 0x65, 0x6c, 0x2e, 0x41, 0x72, 0x67, 0x6f, 0x6e, 0x32, 0x53, 0x70, 0x65,
	0x63, 0x48, 0x00, 0x52, 0x07, 0x61, 0x72, 0x67, 0x6f, 0x6e, 0x32, 0x69, 0x12, 0x30, 0x0a, 0x08,
	0x41, 0x72, 0x67, 0x6f, 0x6e, 0x32, 0x69, 0x64, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x12,
	0x2e, 0x6b, 0x65, 0x72, 0x6e, 0x65, 0x6c, 0x2e, 0x41, 0x72, 0x67, 0x6f, 0x6e, 0x32, 0x53, 0x70,
	0x65, 0x63, 0x48, 0x00, 0x52, 0x08, 0x61, 0x72, 0x67, 0x6f, 0x6e, 0x32, 0x69, 0x64, 0x42, 0x10,
	0x0a, 0x0e, 0x50, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x43, 0x72, 0x79, 0x70, 0x74, 0x73,
	0x22, 0x20, 0x0a, 0x0a, 0x42, 0x43, 0x72, 0x79, 0x70, 0x74, 0x53, 0x70, 0x65, 0x63, 0x12, 0x12,
	0x0a, 0x04, 0x43, 0x6f, 0x73, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x63, 0x6f,
	0x73, 0x74, 0x22, 0x68, 0x0a, 0x0a, 0x53, 0x43, 0x72, 0x79, 0x70, 0x74, 0x53, 0x70, 0x65, 0x63,
	0x12, 0x18, 0x0a, 0x07, 0x53, 0x61, 0x6c, 0x74, 0x4c, 0x65, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x05, 0x52, 0x07, 0x73, 0x61, 0x6c, 0x74, 0x4c, 0x65, 0x6e, 0x12, 0x0c, 0x0a, 0x01, 0x4e, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x01, 0x6e, 0x12, 0x0c, 0x0a, 0x01, 0x52, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x05, 0x52, 0x01, 0x72, 0x12, 0x0c, 0x0a, 0x01, 0x50, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x05, 0x52, 0x01, 0x70, 0x12, 0x16, 0x0a, 0x06, 0x4b, 0x65, 0x79, 0x4c, 0x65, 0x6e, 0x18, 0x05,
	0x20, 0x01, 0x28, 0x05, 0x52, 0x06, 0x6b, 0x65, 0x79, 0x4c, 0x65, 0x6e, 0x22, 0x7d, 0x0a, 0x0a,
	0x50, 0x42, 0x4b, 0x44, 0x46, 0x32, 0x53, 0x70, 0x65, 0x63, 0x12, 0x18, 0x0a, 0x07, 0x53, 0x61,
	0x6c, 0x74, 0x4c, 0x65, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x07, 0x73, 0x61, 0x6c,
	0x74, 0x4c, 0x65, 0x6e, 0x12, 0x12, 0x0a, 0x04, 0x49, 0x74, 0x65, 0x72, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x05, 0x52, 0x04, 0x69, 0x74, 0x65, 0x72, 0x12, 0x16, 0x0a, 0x06, 0x4b, 0x65, 0x79, 0x4c,
	0x65, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x06, 0x6b, 0x65, 0x79, 0x4c, 0x65, 0x6e,
	0x12, 0x29, 0x0a, 0x07, 0x48, 0x61, 0x73, 0x68, 0x41, 0x6c, 0x67, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x0e, 0x32, 0x0f, 0x2e, 0x6b, 0x65, 0x72, 0x6e, 0x65, 0x6c, 0x2e, 0x48, 0x61, 0x73, 0x68, 0x41,
	0x6c, 0x67, 0x52, 0x07, 0x68, 0x61, 0x73, 0x68, 0x41, 0x6c, 0x67, 0x22, 0x8e, 0x01, 0x0a, 0x0a,
	0x41, 0x72, 0x67, 0x6f, 0x6e, 0x32, 0x53, 0x70, 0x65, 0x63, 0x12, 0x18, 0x0a, 0x07, 0x53, 0x61,
	0x6c, 0x74, 0x4c, 0x65, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x07, 0x73, 0x61, 0x6c,
	0x74, 0x4c, 0x65, 0x6e, 0x12, 0x12, 0x0a, 0x04, 0x54, 0x69, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x0d, 0x52, 0x04, 0x74, 0x69, 0x6d, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x4d, 0x65, 0x6d, 0x6f,
	0x72, 0x79, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x06, 0x6d, 0x65, 0x6d, 0x6f, 0x72, 0x79,
	0x12, 0x22, 0x0a, 0x07, 0x54, 0x68, 0x72, 0x65, 0x61, 0x64, 0x73, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x0d, 0x42, 0x08, 0xba, 0x48, 0x05, 0x2a, 0x03, 0x18, 0xff, 0x01, 0x52, 0x07, 0x74, 0x68, 0x72,
	0x65, 0x61, 0x64, 0x73, 0x12, 0x16, 0x0a, 0x06, 0x4b, 0x65, 0x79, 0x4c, 0x65, 0x6e, 0x18, 0x05,
	0x20, 0x01, 0x28, 0x0d, 0x52, 0x06, 0x6b, 0x65, 0x79, 0x4c, 0x65, 0x6e, 0x2a, 0xf0, 0x01, 0x0a,
	0x12, 0x43, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x4b, 0x65, 0x79, 0x43, 0x72, 0x79, 0x70, 0x74, 0x54,
	0x79, 0x70, 0x65, 0x12, 0x1d, 0x0a, 0x19, 0x43, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x4b, 0x65, 0x79,
	0x43, 0x72, 0x79, 0x70, 0x74, 0x54, 0x79, 0x70, 0x65, 0x55, 0x6e, 0x6b, 0x6e, 0x6f, 0x77, 0x6e,
	0x10, 0x00, 0x12, 0x0a, 0x0a, 0x06, 0x41, 0x45, 0x53, 0x47, 0x43, 0x4d, 0x10, 0x01, 0x12, 0x0a,
	0x0a, 0x06, 0x41, 0x45, 0x53, 0x43, 0x42, 0x43, 0x10, 0x02, 0x12, 0x0a, 0x0a, 0x06, 0x41, 0x45,
	0x53, 0x43, 0x46, 0x42, 0x10, 0x03, 0x12, 0x0a, 0x0a, 0x06, 0x41, 0x45, 0x53, 0x43, 0x54, 0x52,
	0x10, 0x04, 0x12, 0x0a, 0x0a, 0x06, 0x41, 0x45, 0x53, 0x4f, 0x46, 0x42, 0x10, 0x05, 0x12, 0x0a,
	0x0a, 0x06, 0x44, 0x45, 0x53, 0x43, 0x42, 0x43, 0x10, 0x06, 0x12, 0x0a, 0x0a, 0x06, 0x44, 0x45,
	0x53, 0x43, 0x46, 0x42, 0x10, 0x07, 0x12, 0x0a, 0x0a, 0x06, 0x44, 0x45, 0x53, 0x43, 0x54, 0x52,
	0x10, 0x08, 0x12, 0x0a, 0x0a, 0x06, 0x44, 0x45, 0x53, 0x4f, 0x46, 0x42, 0x10, 0x09, 0x12, 0x10,
	0x0a, 0x0c, 0x54, 0x72, 0x69, 0x70, 0x6c, 0x65, 0x44, 0x45, 0x53, 0x43, 0x42, 0x43, 0x10, 0x0a,
	0x12, 0x10, 0x0a, 0x0c, 0x54, 0x72, 0x69, 0x70, 0x6c, 0x65, 0x44, 0x45, 0x53, 0x43, 0x46, 0x42,
	0x10, 0x0b, 0x12, 0x10, 0x0a, 0x0c, 0x54, 0x72, 0x69, 0x70, 0x6c, 0x65, 0x44, 0x45, 0x53, 0x43,
	0x54, 0x52, 0x10, 0x0c, 0x12, 0x10, 0x0a, 0x0c, 0x54, 0x72, 0x69, 0x70, 0x6c, 0x65, 0x44, 0x45,
	0x53, 0x4f, 0x46, 0x42, 0x10, 0x0d, 0x12, 0x07, 0x0a, 0x03, 0x52, 0x43, 0x34, 0x10, 0x0e, 0x42,
	0x38, 0x5a, 0x36, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x61, 0x69,
	0x6c, 0x65, 0x72, 0x6f, 0x6e, 0x2d, 0x67, 0x61, 0x74, 0x65, 0x77, 0x61, 0x79, 0x2f, 0x61, 0x69,
	0x6c, 0x65, 0x72, 0x6f, 0x6e, 0x2d, 0x67, 0x61, 0x74, 0x65, 0x77, 0x61, 0x79, 0x2f, 0x61, 0x70,
	0x69, 0x73, 0x2f, 0x6b, 0x65, 0x72, 0x6e, 0x65, 0x6c, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
})

var (
	file_kernel_crypt_proto_rawDescOnce sync.Once
	file_kernel_crypt_proto_rawDescData []byte
)

func file_kernel_crypt_proto_rawDescGZIP() []byte {
	file_kernel_crypt_proto_rawDescOnce.Do(func() {
		file_kernel_crypt_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_kernel_crypt_proto_rawDesc), len(file_kernel_crypt_proto_rawDesc)))
	})
	return file_kernel_crypt_proto_rawDescData
}

var file_kernel_crypt_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_kernel_crypt_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_kernel_crypt_proto_goTypes = []any{
	(CommonKeyCryptType)(0),   // 0: kernel.CommonKeyCryptType
	(*PasswordCryptSpec)(nil), // 1: kernel.PasswordCryptSpec
	(*BCryptSpec)(nil),        // 2: kernel.BCryptSpec
	(*SCryptSpec)(nil),        // 3: kernel.SCryptSpec
	(*PBKDF2Spec)(nil),        // 4: kernel.PBKDF2Spec
	(*Argon2Spec)(nil),        // 5: kernel.Argon2Spec
	(HashAlg)(0),              // 6: kernel.HashAlg
}
var file_kernel_crypt_proto_depIdxs = []int32{
	2, // 0: kernel.PasswordCryptSpec.BCrypt:type_name -> kernel.BCryptSpec
	3, // 1: kernel.PasswordCryptSpec.SCrypt:type_name -> kernel.SCryptSpec
	4, // 2: kernel.PasswordCryptSpec.PBKDF2:type_name -> kernel.PBKDF2Spec
	5, // 3: kernel.PasswordCryptSpec.Argon2i:type_name -> kernel.Argon2Spec
	5, // 4: kernel.PasswordCryptSpec.Argon2id:type_name -> kernel.Argon2Spec
	6, // 5: kernel.PBKDF2Spec.HashAlg:type_name -> kernel.HashAlg
	6, // [6:6] is the sub-list for method output_type
	6, // [6:6] is the sub-list for method input_type
	6, // [6:6] is the sub-list for extension type_name
	6, // [6:6] is the sub-list for extension extendee
	0, // [0:6] is the sub-list for field type_name
}

func init() { file_kernel_crypt_proto_init() }
func file_kernel_crypt_proto_init() {
	if File_kernel_crypt_proto != nil {
		return
	}
	file_kernel_hash_proto_init()
	file_kernel_crypt_proto_msgTypes[0].OneofWrappers = []any{
		(*PasswordCryptSpec_BCrypt)(nil),
		(*PasswordCryptSpec_SCrypt)(nil),
		(*PasswordCryptSpec_PBKDF2)(nil),
		(*PasswordCryptSpec_Argon2I)(nil),
		(*PasswordCryptSpec_Argon2Id)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_kernel_crypt_proto_rawDesc), len(file_kernel_crypt_proto_rawDesc)),
			NumEnums:      1,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_kernel_crypt_proto_goTypes,
		DependencyIndexes: file_kernel_crypt_proto_depIdxs,
		EnumInfos:         file_kernel_crypt_proto_enumTypes,
		MessageInfos:      file_kernel_crypt_proto_msgTypes,
	}.Build()
	File_kernel_crypt_proto = out.File
	file_kernel_crypt_proto_goTypes = nil
	file_kernel_crypt_proto_depIdxs = nil
}
