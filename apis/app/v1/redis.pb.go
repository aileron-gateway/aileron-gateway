// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.3
// 	protoc        v5.29.0
// source: app/v1/storage/redis.proto

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

type RedisClient struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// [REQUIRED]
	// APIVersion is the defined version of the midleware.
	// This value must be "app/v1".
	APIVersion string `protobuf:"bytes,1,opt,name=APIVersion,json=apiVersion,proto3" json:"APIVersion,omitempty"`
	// [REQUIRED]
	// Kind is the kind of this object.
	// This value must be "RedisClient".
	Kind string `protobuf:"bytes,2,opt,name=Kind,json=kind,proto3" json:"Kind,omitempty"`
	// [OPTIONAL]
	// Metadata is the metadata of the http logger object.
	// If not set, both name and namespace in the metadata
	// are treated as "default".
	Metadata *kernel.Metadata `protobuf:"bytes,3,opt,name=Metadata,json=metadata,proto3" json:"Metadata,omitempty"`
	// [OPTIONAL]
	// Spec is the specification of the storage.
	// Default values are used when nothing is set.
	Spec          *RedisClientSpec `protobuf:"bytes,4,opt,name=Spec,json=spec,proto3" json:"Spec,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *RedisClient) Reset() {
	*x = RedisClient{}
	mi := &file_app_v1_storage_redis_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *RedisClient) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RedisClient) ProtoMessage() {}

func (x *RedisClient) ProtoReflect() protoreflect.Message {
	mi := &file_app_v1_storage_redis_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RedisClient.ProtoReflect.Descriptor instead.
func (*RedisClient) Descriptor() ([]byte, []int) {
	return file_app_v1_storage_redis_proto_rawDescGZIP(), []int{0}
}

func (x *RedisClient) GetAPIVersion() string {
	if x != nil {
		return x.APIVersion
	}
	return ""
}

func (x *RedisClient) GetKind() string {
	if x != nil {
		return x.Kind
	}
	return ""
}

func (x *RedisClient) GetMetadata() *kernel.Metadata {
	if x != nil {
		return x.Metadata
	}
	return nil
}

func (x *RedisClient) GetSpec() *RedisClientSpec {
	if x != nil {
		return x.Spec
	}
	return nil
}

// RedisClientSpec is the spec for redis universal client.
// See https://pkg.go.dev/github.com/go-redis/redis#UniversalOptions for details.
type RedisClientSpec struct {
	state                 protoimpl.MessageState `protogen:"open.v1"`
	Addrs                 []string               `protobuf:"bytes,1,rep,name=Addrs,json=addrs,proto3" json:"Addrs,omitempty"`
	Name                  string                 `protobuf:"bytes,2,opt,name=Name,json=name,proto3" json:"Name,omitempty"`
	DB                    int32                  `protobuf:"varint,3,opt,name=DB,json=db,proto3" json:"DB,omitempty"`
	Username              string                 `protobuf:"bytes,4,opt,name=Username,json=username,proto3" json:"Username,omitempty"`
	Password              string                 `protobuf:"bytes,5,opt,name=Password,json=password,proto3" json:"Password,omitempty"`
	SentinelUsername      string                 `protobuf:"bytes,6,opt,name=SentinelUsername,json=sentinelUsername,proto3" json:"SentinelUsername,omitempty"`
	SentinelPassword      string                 `protobuf:"bytes,7,opt,name=SentinelPassword,json=sentinelPassword,proto3" json:"SentinelPassword,omitempty"`
	MaxRetries            int32                  `protobuf:"varint,8,opt,name=MaxRetries,json=maxRetries,proto3" json:"MaxRetries,omitempty"`
	MinRetryBackoff       int32                  `protobuf:"varint,9,opt,name=MinRetryBackoff,json=minRetryBackoff,proto3" json:"MinRetryBackoff,omitempty"`
	MaxRetryBackoff       int32                  `protobuf:"varint,10,opt,name=MaxRetryBackoff,json=maxRetryBackoff,proto3" json:"MaxRetryBackoff,omitempty"`
	DialTimeout           int32                  `protobuf:"varint,11,opt,name=DialTimeout,json=dialTimeout,proto3" json:"DialTimeout,omitempty"`
	ReadTimeout           int32                  `protobuf:"varint,12,opt,name=ReadTimeout,json=readTimeout,proto3" json:"ReadTimeout,omitempty"`
	WriteTimeout          int32                  `protobuf:"varint,13,opt,name=WriteTimeout,json=writeTimeout,proto3" json:"WriteTimeout,omitempty"`
	ContextTimeoutEnabled bool                   `protobuf:"varint,14,opt,name=ContextTimeoutEnabled,json=contextTimeoutEnabled,proto3" json:"ContextTimeoutEnabled,omitempty"`
	PoolFIFO              bool                   `protobuf:"varint,15,opt,name=PoolFIFO,json=poolFIFO,proto3" json:"PoolFIFO,omitempty"`
	PoolSize              int32                  `protobuf:"varint,16,opt,name=PoolSize,json=poolSize,proto3" json:"PoolSize,omitempty"`
	PoolTimeout           int32                  `protobuf:"varint,17,opt,name=PoolTimeout,json=poolTimeout,proto3" json:"PoolTimeout,omitempty"`
	MinIdleConns          int32                  `protobuf:"varint,18,opt,name=MinIdleConns,json=minIdleConns,proto3" json:"MinIdleConns,omitempty"`
	MaxIdleConns          int32                  `protobuf:"varint,19,opt,name=MaxIdleConns,json=maxIdleConns,proto3" json:"MaxIdleConns,omitempty"`
	ConnMaxIdleTime       int32                  `protobuf:"varint,20,opt,name=ConnMaxIdleTime,json=connMaxIdleTime,proto3" json:"ConnMaxIdleTime,omitempty"`
	ConnMaxLifetime       int32                  `protobuf:"varint,21,opt,name=ConnMaxLifetime,json=connMaxLifetime,proto3" json:"ConnMaxLifetime,omitempty"`
	TLSConfig             *kernel.TLSConfig      `protobuf:"bytes,22,opt,name=TLSConfig,json=tlsConfig,proto3" json:"TLSConfig,omitempty"`
	MaxRedirects          int32                  `protobuf:"varint,23,opt,name=MaxRedirects,json=maxRedirects,proto3" json:"MaxRedirects,omitempty"`
	ReadOnly              bool                   `protobuf:"varint,24,opt,name=ReadOnly,json=readOnly,proto3" json:"ReadOnly,omitempty"`
	RouteByLatency        bool                   `protobuf:"varint,25,opt,name=RouteByLatency,json=routeByLatency,proto3" json:"RouteByLatency,omitempty"`
	RouteRandomly         bool                   `protobuf:"varint,26,opt,name=RouteRandomly,json=routeRandomly,proto3" json:"RouteRandomly,omitempty"`
	MasterName            string                 `protobuf:"bytes,27,opt,name=MasterName,json=masterName,proto3" json:"MasterName,omitempty"`
	Timeout               int64                  `protobuf:"varint,28,opt,name=Timeout,json=timeout,proto3" json:"Timeout,omitempty"`
	Expiration            int64                  `protobuf:"varint,29,opt,name=Expiration,json=expiration,proto3" json:"Expiration,omitempty"`
	unknownFields         protoimpl.UnknownFields
	sizeCache             protoimpl.SizeCache
}

func (x *RedisClientSpec) Reset() {
	*x = RedisClientSpec{}
	mi := &file_app_v1_storage_redis_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *RedisClientSpec) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RedisClientSpec) ProtoMessage() {}

func (x *RedisClientSpec) ProtoReflect() protoreflect.Message {
	mi := &file_app_v1_storage_redis_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RedisClientSpec.ProtoReflect.Descriptor instead.
func (*RedisClientSpec) Descriptor() ([]byte, []int) {
	return file_app_v1_storage_redis_proto_rawDescGZIP(), []int{1}
}

func (x *RedisClientSpec) GetAddrs() []string {
	if x != nil {
		return x.Addrs
	}
	return nil
}

func (x *RedisClientSpec) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *RedisClientSpec) GetDB() int32 {
	if x != nil {
		return x.DB
	}
	return 0
}

func (x *RedisClientSpec) GetUsername() string {
	if x != nil {
		return x.Username
	}
	return ""
}

func (x *RedisClientSpec) GetPassword() string {
	if x != nil {
		return x.Password
	}
	return ""
}

func (x *RedisClientSpec) GetSentinelUsername() string {
	if x != nil {
		return x.SentinelUsername
	}
	return ""
}

func (x *RedisClientSpec) GetSentinelPassword() string {
	if x != nil {
		return x.SentinelPassword
	}
	return ""
}

func (x *RedisClientSpec) GetMaxRetries() int32 {
	if x != nil {
		return x.MaxRetries
	}
	return 0
}

func (x *RedisClientSpec) GetMinRetryBackoff() int32 {
	if x != nil {
		return x.MinRetryBackoff
	}
	return 0
}

func (x *RedisClientSpec) GetMaxRetryBackoff() int32 {
	if x != nil {
		return x.MaxRetryBackoff
	}
	return 0
}

func (x *RedisClientSpec) GetDialTimeout() int32 {
	if x != nil {
		return x.DialTimeout
	}
	return 0
}

func (x *RedisClientSpec) GetReadTimeout() int32 {
	if x != nil {
		return x.ReadTimeout
	}
	return 0
}

func (x *RedisClientSpec) GetWriteTimeout() int32 {
	if x != nil {
		return x.WriteTimeout
	}
	return 0
}

func (x *RedisClientSpec) GetContextTimeoutEnabled() bool {
	if x != nil {
		return x.ContextTimeoutEnabled
	}
	return false
}

func (x *RedisClientSpec) GetPoolFIFO() bool {
	if x != nil {
		return x.PoolFIFO
	}
	return false
}

func (x *RedisClientSpec) GetPoolSize() int32 {
	if x != nil {
		return x.PoolSize
	}
	return 0
}

func (x *RedisClientSpec) GetPoolTimeout() int32 {
	if x != nil {
		return x.PoolTimeout
	}
	return 0
}

func (x *RedisClientSpec) GetMinIdleConns() int32 {
	if x != nil {
		return x.MinIdleConns
	}
	return 0
}

func (x *RedisClientSpec) GetMaxIdleConns() int32 {
	if x != nil {
		return x.MaxIdleConns
	}
	return 0
}

func (x *RedisClientSpec) GetConnMaxIdleTime() int32 {
	if x != nil {
		return x.ConnMaxIdleTime
	}
	return 0
}

func (x *RedisClientSpec) GetConnMaxLifetime() int32 {
	if x != nil {
		return x.ConnMaxLifetime
	}
	return 0
}

func (x *RedisClientSpec) GetTLSConfig() *kernel.TLSConfig {
	if x != nil {
		return x.TLSConfig
	}
	return nil
}

func (x *RedisClientSpec) GetMaxRedirects() int32 {
	if x != nil {
		return x.MaxRedirects
	}
	return 0
}

func (x *RedisClientSpec) GetReadOnly() bool {
	if x != nil {
		return x.ReadOnly
	}
	return false
}

func (x *RedisClientSpec) GetRouteByLatency() bool {
	if x != nil {
		return x.RouteByLatency
	}
	return false
}

func (x *RedisClientSpec) GetRouteRandomly() bool {
	if x != nil {
		return x.RouteRandomly
	}
	return false
}

func (x *RedisClientSpec) GetMasterName() string {
	if x != nil {
		return x.MasterName
	}
	return ""
}

func (x *RedisClientSpec) GetTimeout() int64 {
	if x != nil {
		return x.Timeout
	}
	return 0
}

func (x *RedisClientSpec) GetExpiration() int64 {
	if x != nil {
		return x.Expiration
	}
	return 0
}

var File_app_v1_storage_redis_proto protoreflect.FileDescriptor

var file_app_v1_storage_redis_proto_rawDesc = []byte{
	0x0a, 0x1a, 0x61, 0x70, 0x70, 0x2f, 0x76, 0x31, 0x2f, 0x73, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65,
	0x2f, 0x72, 0x65, 0x64, 0x69, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x06, 0x61, 0x70,
	0x70, 0x2e, 0x76, 0x31, 0x1a, 0x1b, 0x62, 0x75, 0x66, 0x2f, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61,
	0x74, 0x65, 0x2f, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x1a, 0x14, 0x6b, 0x65, 0x72, 0x6e, 0x65, 0x6c, 0x2f, 0x6e, 0x65, 0x74, 0x77, 0x6f, 0x72,
	0x6b, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x15, 0x6b, 0x65, 0x72, 0x6e, 0x65, 0x6c, 0x2f,
	0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xbf,
	0x01, 0x0a, 0x0b, 0x52, 0x65, 0x64, 0x69, 0x73, 0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x12, 0x2d,
	0x0a, 0x0a, 0x41, 0x50, 0x49, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x42, 0x0d, 0xba, 0x48, 0x0a, 0x72, 0x08, 0x0a, 0x06, 0x61, 0x70, 0x70, 0x2f, 0x76,
	0x31, 0x52, 0x0a, 0x61, 0x70, 0x69, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x26, 0x0a,
	0x04, 0x4b, 0x69, 0x6e, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x42, 0x12, 0xba, 0x48, 0x0f,
	0x72, 0x0d, 0x0a, 0x0b, 0x52, 0x65, 0x64, 0x69, 0x73, 0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x52,
	0x04, 0x6b, 0x69, 0x6e, 0x64, 0x12, 0x2c, 0x0a, 0x08, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74,
	0x61, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x6b, 0x65, 0x72, 0x6e, 0x65, 0x6c,
	0x2e, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x52, 0x08, 0x6d, 0x65, 0x74, 0x61, 0x64,
	0x61, 0x74, 0x61, 0x12, 0x2b, 0x0a, 0x04, 0x53, 0x70, 0x65, 0x63, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x17, 0x2e, 0x61, 0x70, 0x70, 0x2e, 0x76, 0x31, 0x2e, 0x52, 0x65, 0x64, 0x69, 0x73,
	0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x53, 0x70, 0x65, 0x63, 0x52, 0x04, 0x73, 0x70, 0x65, 0x63,
	0x22, 0xfc, 0x07, 0x0a, 0x0f, 0x52, 0x65, 0x64, 0x69, 0x73, 0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74,
	0x53, 0x70, 0x65, 0x63, 0x12, 0x14, 0x0a, 0x05, 0x41, 0x64, 0x64, 0x72, 0x73, 0x18, 0x01, 0x20,
	0x03, 0x28, 0x09, 0x52, 0x05, 0x61, 0x64, 0x64, 0x72, 0x73, 0x12, 0x12, 0x0a, 0x04, 0x4e, 0x61,
	0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x0e,
	0x0a, 0x02, 0x44, 0x42, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x02, 0x64, 0x62, 0x12, 0x1a,
	0x0a, 0x08, 0x55, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x08, 0x75, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x50, 0x61,
	0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x70, 0x61,
	0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x12, 0x2a, 0x0a, 0x10, 0x53, 0x65, 0x6e, 0x74, 0x69, 0x6e,
	0x65, 0x6c, 0x55, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x10, 0x73, 0x65, 0x6e, 0x74, 0x69, 0x6e, 0x65, 0x6c, 0x55, 0x73, 0x65, 0x72, 0x6e, 0x61,
	0x6d, 0x65, 0x12, 0x2a, 0x0a, 0x10, 0x53, 0x65, 0x6e, 0x74, 0x69, 0x6e, 0x65, 0x6c, 0x50, 0x61,
	0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52, 0x10, 0x73, 0x65,
	0x6e, 0x74, 0x69, 0x6e, 0x65, 0x6c, 0x50, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x12, 0x1e,
	0x0a, 0x0a, 0x4d, 0x61, 0x78, 0x52, 0x65, 0x74, 0x72, 0x69, 0x65, 0x73, 0x18, 0x08, 0x20, 0x01,
	0x28, 0x05, 0x52, 0x0a, 0x6d, 0x61, 0x78, 0x52, 0x65, 0x74, 0x72, 0x69, 0x65, 0x73, 0x12, 0x28,
	0x0a, 0x0f, 0x4d, 0x69, 0x6e, 0x52, 0x65, 0x74, 0x72, 0x79, 0x42, 0x61, 0x63, 0x6b, 0x6f, 0x66,
	0x66, 0x18, 0x09, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0f, 0x6d, 0x69, 0x6e, 0x52, 0x65, 0x74, 0x72,
	0x79, 0x42, 0x61, 0x63, 0x6b, 0x6f, 0x66, 0x66, 0x12, 0x28, 0x0a, 0x0f, 0x4d, 0x61, 0x78, 0x52,
	0x65, 0x74, 0x72, 0x79, 0x42, 0x61, 0x63, 0x6b, 0x6f, 0x66, 0x66, 0x18, 0x0a, 0x20, 0x01, 0x28,
	0x05, 0x52, 0x0f, 0x6d, 0x61, 0x78, 0x52, 0x65, 0x74, 0x72, 0x79, 0x42, 0x61, 0x63, 0x6b, 0x6f,
	0x66, 0x66, 0x12, 0x20, 0x0a, 0x0b, 0x44, 0x69, 0x61, 0x6c, 0x54, 0x69, 0x6d, 0x65, 0x6f, 0x75,
	0x74, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0b, 0x64, 0x69, 0x61, 0x6c, 0x54, 0x69, 0x6d,
	0x65, 0x6f, 0x75, 0x74, 0x12, 0x20, 0x0a, 0x0b, 0x52, 0x65, 0x61, 0x64, 0x54, 0x69, 0x6d, 0x65,
	0x6f, 0x75, 0x74, 0x18, 0x0c, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0b, 0x72, 0x65, 0x61, 0x64, 0x54,
	0x69, 0x6d, 0x65, 0x6f, 0x75, 0x74, 0x12, 0x22, 0x0a, 0x0c, 0x57, 0x72, 0x69, 0x74, 0x65, 0x54,
	0x69, 0x6d, 0x65, 0x6f, 0x75, 0x74, 0x18, 0x0d, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0c, 0x77, 0x72,
	0x69, 0x74, 0x65, 0x54, 0x69, 0x6d, 0x65, 0x6f, 0x75, 0x74, 0x12, 0x34, 0x0a, 0x15, 0x43, 0x6f,
	0x6e, 0x74, 0x65, 0x78, 0x74, 0x54, 0x69, 0x6d, 0x65, 0x6f, 0x75, 0x74, 0x45, 0x6e, 0x61, 0x62,
	0x6c, 0x65, 0x64, 0x18, 0x0e, 0x20, 0x01, 0x28, 0x08, 0x52, 0x15, 0x63, 0x6f, 0x6e, 0x74, 0x65,
	0x78, 0x74, 0x54, 0x69, 0x6d, 0x65, 0x6f, 0x75, 0x74, 0x45, 0x6e, 0x61, 0x62, 0x6c, 0x65, 0x64,
	0x12, 0x1a, 0x0a, 0x08, 0x50, 0x6f, 0x6f, 0x6c, 0x46, 0x49, 0x46, 0x4f, 0x18, 0x0f, 0x20, 0x01,
	0x28, 0x08, 0x52, 0x08, 0x70, 0x6f, 0x6f, 0x6c, 0x46, 0x49, 0x46, 0x4f, 0x12, 0x1a, 0x0a, 0x08,
	0x50, 0x6f, 0x6f, 0x6c, 0x53, 0x69, 0x7a, 0x65, 0x18, 0x10, 0x20, 0x01, 0x28, 0x05, 0x52, 0x08,
	0x70, 0x6f, 0x6f, 0x6c, 0x53, 0x69, 0x7a, 0x65, 0x12, 0x20, 0x0a, 0x0b, 0x50, 0x6f, 0x6f, 0x6c,
	0x54, 0x69, 0x6d, 0x65, 0x6f, 0x75, 0x74, 0x18, 0x11, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0b, 0x70,
	0x6f, 0x6f, 0x6c, 0x54, 0x69, 0x6d, 0x65, 0x6f, 0x75, 0x74, 0x12, 0x22, 0x0a, 0x0c, 0x4d, 0x69,
	0x6e, 0x49, 0x64, 0x6c, 0x65, 0x43, 0x6f, 0x6e, 0x6e, 0x73, 0x18, 0x12, 0x20, 0x01, 0x28, 0x05,
	0x52, 0x0c, 0x6d, 0x69, 0x6e, 0x49, 0x64, 0x6c, 0x65, 0x43, 0x6f, 0x6e, 0x6e, 0x73, 0x12, 0x22,
	0x0a, 0x0c, 0x4d, 0x61, 0x78, 0x49, 0x64, 0x6c, 0x65, 0x43, 0x6f, 0x6e, 0x6e, 0x73, 0x18, 0x13,
	0x20, 0x01, 0x28, 0x05, 0x52, 0x0c, 0x6d, 0x61, 0x78, 0x49, 0x64, 0x6c, 0x65, 0x43, 0x6f, 0x6e,
	0x6e, 0x73, 0x12, 0x28, 0x0a, 0x0f, 0x43, 0x6f, 0x6e, 0x6e, 0x4d, 0x61, 0x78, 0x49, 0x64, 0x6c,
	0x65, 0x54, 0x69, 0x6d, 0x65, 0x18, 0x14, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0f, 0x63, 0x6f, 0x6e,
	0x6e, 0x4d, 0x61, 0x78, 0x49, 0x64, 0x6c, 0x65, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x28, 0x0a, 0x0f,
	0x43, 0x6f, 0x6e, 0x6e, 0x4d, 0x61, 0x78, 0x4c, 0x69, 0x66, 0x65, 0x74, 0x69, 0x6d, 0x65, 0x18,
	0x15, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0f, 0x63, 0x6f, 0x6e, 0x6e, 0x4d, 0x61, 0x78, 0x4c, 0x69,
	0x66, 0x65, 0x74, 0x69, 0x6d, 0x65, 0x12, 0x2f, 0x0a, 0x09, 0x54, 0x4c, 0x53, 0x43, 0x6f, 0x6e,
	0x66, 0x69, 0x67, 0x18, 0x16, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x11, 0x2e, 0x6b, 0x65, 0x72, 0x6e,
	0x65, 0x6c, 0x2e, 0x54, 0x4c, 0x53, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x52, 0x09, 0x74, 0x6c,
	0x73, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x12, 0x22, 0x0a, 0x0c, 0x4d, 0x61, 0x78, 0x52, 0x65,
	0x64, 0x69, 0x72, 0x65, 0x63, 0x74, 0x73, 0x18, 0x17, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0c, 0x6d,
	0x61, 0x78, 0x52, 0x65, 0x64, 0x69, 0x72, 0x65, 0x63, 0x74, 0x73, 0x12, 0x1a, 0x0a, 0x08, 0x52,
	0x65, 0x61, 0x64, 0x4f, 0x6e, 0x6c, 0x79, 0x18, 0x18, 0x20, 0x01, 0x28, 0x08, 0x52, 0x08, 0x72,
	0x65, 0x61, 0x64, 0x4f, 0x6e, 0x6c, 0x79, 0x12, 0x26, 0x0a, 0x0e, 0x52, 0x6f, 0x75, 0x74, 0x65,
	0x42, 0x79, 0x4c, 0x61, 0x74, 0x65, 0x6e, 0x63, 0x79, 0x18, 0x19, 0x20, 0x01, 0x28, 0x08, 0x52,
	0x0e, 0x72, 0x6f, 0x75, 0x74, 0x65, 0x42, 0x79, 0x4c, 0x61, 0x74, 0x65, 0x6e, 0x63, 0x79, 0x12,
	0x24, 0x0a, 0x0d, 0x52, 0x6f, 0x75, 0x74, 0x65, 0x52, 0x61, 0x6e, 0x64, 0x6f, 0x6d, 0x6c, 0x79,
	0x18, 0x1a, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0d, 0x72, 0x6f, 0x75, 0x74, 0x65, 0x52, 0x61, 0x6e,
	0x64, 0x6f, 0x6d, 0x6c, 0x79, 0x12, 0x1e, 0x0a, 0x0a, 0x4d, 0x61, 0x73, 0x74, 0x65, 0x72, 0x4e,
	0x61, 0x6d, 0x65, 0x18, 0x1b, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x6d, 0x61, 0x73, 0x74, 0x65,
	0x72, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x54, 0x69, 0x6d, 0x65, 0x6f, 0x75, 0x74,
	0x18, 0x1c, 0x20, 0x01, 0x28, 0x03, 0x52, 0x07, 0x74, 0x69, 0x6d, 0x65, 0x6f, 0x75, 0x74, 0x12,
	0x1e, 0x0a, 0x0a, 0x45, 0x78, 0x70, 0x69, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x1d, 0x20,
	0x01, 0x28, 0x03, 0x52, 0x0a, 0x65, 0x78, 0x70, 0x69, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x42,
	0x38, 0x5a, 0x36, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x61, 0x69,
	0x6c, 0x65, 0x72, 0x6f, 0x6e, 0x2d, 0x67, 0x61, 0x74, 0x65, 0x77, 0x61, 0x79, 0x2f, 0x61, 0x69,
	0x6c, 0x65, 0x72, 0x6f, 0x6e, 0x2d, 0x67, 0x61, 0x74, 0x65, 0x77, 0x61, 0x79, 0x2f, 0x61, 0x70,
	0x69, 0x73, 0x2f, 0x61, 0x70, 0x70, 0x2f, 0x76, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var (
	file_app_v1_storage_redis_proto_rawDescOnce sync.Once
	file_app_v1_storage_redis_proto_rawDescData = file_app_v1_storage_redis_proto_rawDesc
)

func file_app_v1_storage_redis_proto_rawDescGZIP() []byte {
	file_app_v1_storage_redis_proto_rawDescOnce.Do(func() {
		file_app_v1_storage_redis_proto_rawDescData = protoimpl.X.CompressGZIP(file_app_v1_storage_redis_proto_rawDescData)
	})
	return file_app_v1_storage_redis_proto_rawDescData
}

var file_app_v1_storage_redis_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_app_v1_storage_redis_proto_goTypes = []any{
	(*RedisClient)(nil),      // 0: app.v1.RedisClient
	(*RedisClientSpec)(nil),  // 1: app.v1.RedisClientSpec
	(*kernel.Metadata)(nil),  // 2: kernel.Metadata
	(*kernel.TLSConfig)(nil), // 3: kernel.TLSConfig
}
var file_app_v1_storage_redis_proto_depIdxs = []int32{
	2, // 0: app.v1.RedisClient.Metadata:type_name -> kernel.Metadata
	1, // 1: app.v1.RedisClient.Spec:type_name -> app.v1.RedisClientSpec
	3, // 2: app.v1.RedisClientSpec.TLSConfig:type_name -> kernel.TLSConfig
	3, // [3:3] is the sub-list for method output_type
	3, // [3:3] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_app_v1_storage_redis_proto_init() }
func file_app_v1_storage_redis_proto_init() {
	if File_app_v1_storage_redis_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_app_v1_storage_redis_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_app_v1_storage_redis_proto_goTypes,
		DependencyIndexes: file_app_v1_storage_redis_proto_depIdxs,
		MessageInfos:      file_app_v1_storage_redis_proto_msgTypes,
	}.Build()
	File_app_v1_storage_redis_proto = out.File
	file_app_v1_storage_redis_proto_rawDesc = nil
	file_app_v1_storage_redis_proto_goTypes = nil
	file_app_v1_storage_redis_proto_depIdxs = nil
}
