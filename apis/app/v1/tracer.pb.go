// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.3
// 	protoc        v5.29.0
// source: app/v1/tracer/tracer.proto

package v1

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

type K8SAttributesSpec struct {
	state                 protoimpl.MessageState `protogen:"open.v1"`
	ClusterName           string                 `protobuf:"bytes,1,opt,name=ClusterName,json=clusterName,proto3" json:"ClusterName,omitempty"`
	ContainerName         string                 `protobuf:"bytes,2,opt,name=ContainerName,json=containerName,proto3" json:"ContainerName,omitempty"`
	ContainerRestartCount string                 `protobuf:"bytes,3,opt,name=ContainerRestartCount,json=containerRestartCount,proto3" json:"ContainerRestartCount,omitempty"`
	CronJobName           string                 `protobuf:"bytes,4,opt,name=CronJobName,json=cronJobName,proto3" json:"CronJobName,omitempty"`
	CronJobUID            string                 `protobuf:"bytes,5,opt,name=CronJobUID,json=cronJobUID,proto3" json:"CronJobUID,omitempty"`
	DaemonSetName         string                 `protobuf:"bytes,6,opt,name=DaemonSetName,json=daemonSetName,proto3" json:"DaemonSetName,omitempty"`
	DaemonSetUID          string                 `protobuf:"bytes,7,opt,name=DaemonSetUID,json=daemonSetUID,proto3" json:"DaemonSetUID,omitempty"`
	DeploymentName        string                 `protobuf:"bytes,8,opt,name=DeploymentName,json=deploymentName,proto3" json:"DeploymentName,omitempty"`
	DeploymentUID         string                 `protobuf:"bytes,9,opt,name=DeploymentUID,json=deploymentUID,proto3" json:"DeploymentUID,omitempty"`
	JobName               string                 `protobuf:"bytes,10,opt,name=JobName,json=jobName,proto3" json:"JobName,omitempty"`
	JobUID                string                 `protobuf:"bytes,11,opt,name=JobUID,json=jobUID,proto3" json:"JobUID,omitempty"`
	NamespaceName         string                 `protobuf:"bytes,12,opt,name=NamespaceName,json=namespaceName,proto3" json:"NamespaceName,omitempty"`
	NodeName              string                 `protobuf:"bytes,13,opt,name=NodeName,json=nodeName,proto3" json:"NodeName,omitempty"`
	NodeUID               string                 `protobuf:"bytes,14,opt,name=NodeUID,json=nodeUID,proto3" json:"NodeUID,omitempty"`
	PodName               string                 `protobuf:"bytes,15,opt,name=PodName,json=podName,proto3" json:"PodName,omitempty"`
	PodUID                string                 `protobuf:"bytes,16,opt,name=PodUID,json=podUID,proto3" json:"PodUID,omitempty"`
	ReplicaSetName        string                 `protobuf:"bytes,17,opt,name=ReplicaSetName,json=replicaSetName,proto3" json:"ReplicaSetName,omitempty"`
	ReplicaSetUID         string                 `protobuf:"bytes,18,opt,name=ReplicaSetUID,json=replicaSetUID,proto3" json:"ReplicaSetUID,omitempty"`
	StatefulSetName       string                 `protobuf:"bytes,19,opt,name=StatefulSetName,json=statefulSetName,proto3" json:"StatefulSetName,omitempty"`
	StatefulSetUID        string                 `protobuf:"bytes,20,opt,name=StatefulSetUID,json=statefulSetUID,proto3" json:"StatefulSetUID,omitempty"`
	unknownFields         protoimpl.UnknownFields
	sizeCache             protoimpl.SizeCache
}

func (x *K8SAttributesSpec) Reset() {
	*x = K8SAttributesSpec{}
	mi := &file_app_v1_tracer_tracer_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *K8SAttributesSpec) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*K8SAttributesSpec) ProtoMessage() {}

func (x *K8SAttributesSpec) ProtoReflect() protoreflect.Message {
	mi := &file_app_v1_tracer_tracer_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use K8SAttributesSpec.ProtoReflect.Descriptor instead.
func (*K8SAttributesSpec) Descriptor() ([]byte, []int) {
	return file_app_v1_tracer_tracer_proto_rawDescGZIP(), []int{0}
}

func (x *K8SAttributesSpec) GetClusterName() string {
	if x != nil {
		return x.ClusterName
	}
	return ""
}

func (x *K8SAttributesSpec) GetContainerName() string {
	if x != nil {
		return x.ContainerName
	}
	return ""
}

func (x *K8SAttributesSpec) GetContainerRestartCount() string {
	if x != nil {
		return x.ContainerRestartCount
	}
	return ""
}

func (x *K8SAttributesSpec) GetCronJobName() string {
	if x != nil {
		return x.CronJobName
	}
	return ""
}

func (x *K8SAttributesSpec) GetCronJobUID() string {
	if x != nil {
		return x.CronJobUID
	}
	return ""
}

func (x *K8SAttributesSpec) GetDaemonSetName() string {
	if x != nil {
		return x.DaemonSetName
	}
	return ""
}

func (x *K8SAttributesSpec) GetDaemonSetUID() string {
	if x != nil {
		return x.DaemonSetUID
	}
	return ""
}

func (x *K8SAttributesSpec) GetDeploymentName() string {
	if x != nil {
		return x.DeploymentName
	}
	return ""
}

func (x *K8SAttributesSpec) GetDeploymentUID() string {
	if x != nil {
		return x.DeploymentUID
	}
	return ""
}

func (x *K8SAttributesSpec) GetJobName() string {
	if x != nil {
		return x.JobName
	}
	return ""
}

func (x *K8SAttributesSpec) GetJobUID() string {
	if x != nil {
		return x.JobUID
	}
	return ""
}

func (x *K8SAttributesSpec) GetNamespaceName() string {
	if x != nil {
		return x.NamespaceName
	}
	return ""
}

func (x *K8SAttributesSpec) GetNodeName() string {
	if x != nil {
		return x.NodeName
	}
	return ""
}

func (x *K8SAttributesSpec) GetNodeUID() string {
	if x != nil {
		return x.NodeUID
	}
	return ""
}

func (x *K8SAttributesSpec) GetPodName() string {
	if x != nil {
		return x.PodName
	}
	return ""
}

func (x *K8SAttributesSpec) GetPodUID() string {
	if x != nil {
		return x.PodUID
	}
	return ""
}

func (x *K8SAttributesSpec) GetReplicaSetName() string {
	if x != nil {
		return x.ReplicaSetName
	}
	return ""
}

func (x *K8SAttributesSpec) GetReplicaSetUID() string {
	if x != nil {
		return x.ReplicaSetUID
	}
	return ""
}

func (x *K8SAttributesSpec) GetStatefulSetName() string {
	if x != nil {
		return x.StatefulSetName
	}
	return ""
}

func (x *K8SAttributesSpec) GetStatefulSetUID() string {
	if x != nil {
		return x.StatefulSetUID
	}
	return ""
}

type ContainerAttributesSpec struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	ID            string                 `protobuf:"bytes,1,opt,name=ID,json=id,proto3" json:"ID,omitempty"`
	ImageName     string                 `protobuf:"bytes,2,opt,name=ImageName,json=imageName,proto3" json:"ImageName,omitempty"`
	ImageTag      string                 `protobuf:"bytes,3,opt,name=ImageTag,json=imageTag,proto3" json:"ImageTag,omitempty"`
	Name          string                 `protobuf:"bytes,4,opt,name=Name,json=name,proto3" json:"Name,omitempty"`
	Runtime       string                 `protobuf:"bytes,5,opt,name=Runtime,json=runtime,proto3" json:"Runtime,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ContainerAttributesSpec) Reset() {
	*x = ContainerAttributesSpec{}
	mi := &file_app_v1_tracer_tracer_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ContainerAttributesSpec) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ContainerAttributesSpec) ProtoMessage() {}

func (x *ContainerAttributesSpec) ProtoReflect() protoreflect.Message {
	mi := &file_app_v1_tracer_tracer_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ContainerAttributesSpec.ProtoReflect.Descriptor instead.
func (*ContainerAttributesSpec) Descriptor() ([]byte, []int) {
	return file_app_v1_tracer_tracer_proto_rawDescGZIP(), []int{1}
}

func (x *ContainerAttributesSpec) GetID() string {
	if x != nil {
		return x.ID
	}
	return ""
}

func (x *ContainerAttributesSpec) GetImageName() string {
	if x != nil {
		return x.ImageName
	}
	return ""
}

func (x *ContainerAttributesSpec) GetImageTag() string {
	if x != nil {
		return x.ImageTag
	}
	return ""
}

func (x *ContainerAttributesSpec) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *ContainerAttributesSpec) GetRuntime() string {
	if x != nil {
		return x.Runtime
	}
	return ""
}

type HostAttributesSpec struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	ID            string                 `protobuf:"bytes,1,opt,name=ID,json=id,proto3" json:"ID,omitempty"`
	ImageID       string                 `protobuf:"bytes,2,opt,name=ImageID,json=imageID,proto3" json:"ImageID,omitempty"`
	ImageName     string                 `protobuf:"bytes,3,opt,name=ImageName,json=imageName,proto3" json:"ImageName,omitempty"`
	ImageVersion  string                 `protobuf:"bytes,4,opt,name=ImageVersion,json=imageVersion,proto3" json:"ImageVersion,omitempty"`
	Name          string                 `protobuf:"bytes,5,opt,name=Name,json=name,proto3" json:"Name,omitempty"`
	Type          string                 `protobuf:"bytes,6,opt,name=Type,json=type,proto3" json:"Type,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *HostAttributesSpec) Reset() {
	*x = HostAttributesSpec{}
	mi := &file_app_v1_tracer_tracer_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *HostAttributesSpec) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HostAttributesSpec) ProtoMessage() {}

func (x *HostAttributesSpec) ProtoReflect() protoreflect.Message {
	mi := &file_app_v1_tracer_tracer_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HostAttributesSpec.ProtoReflect.Descriptor instead.
func (*HostAttributesSpec) Descriptor() ([]byte, []int) {
	return file_app_v1_tracer_tracer_proto_rawDescGZIP(), []int{2}
}

func (x *HostAttributesSpec) GetID() string {
	if x != nil {
		return x.ID
	}
	return ""
}

func (x *HostAttributesSpec) GetImageID() string {
	if x != nil {
		return x.ImageID
	}
	return ""
}

func (x *HostAttributesSpec) GetImageName() string {
	if x != nil {
		return x.ImageName
	}
	return ""
}

func (x *HostAttributesSpec) GetImageVersion() string {
	if x != nil {
		return x.ImageVersion
	}
	return ""
}

func (x *HostAttributesSpec) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *HostAttributesSpec) GetType() string {
	if x != nil {
		return x.Type
	}
	return ""
}

var File_app_v1_tracer_tracer_proto protoreflect.FileDescriptor

var file_app_v1_tracer_tracer_proto_rawDesc = []byte{
	0x0a, 0x1a, 0x61, 0x70, 0x70, 0x2f, 0x76, 0x31, 0x2f, 0x74, 0x72, 0x61, 0x63, 0x65, 0x72, 0x2f,
	0x74, 0x72, 0x61, 0x63, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x06, 0x61, 0x70,
	0x70, 0x2e, 0x76, 0x31, 0x22, 0xcb, 0x05, 0x0a, 0x11, 0x4b, 0x38, 0x73, 0x41, 0x74, 0x74, 0x72,
	0x69, 0x62, 0x75, 0x74, 0x65, 0x73, 0x53, 0x70, 0x65, 0x63, 0x12, 0x20, 0x0a, 0x0b, 0x43, 0x6c,
	0x75, 0x73, 0x74, 0x65, 0x72, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x0b, 0x63, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x24, 0x0a, 0x0d,
	0x43, 0x6f, 0x6e, 0x74, 0x61, 0x69, 0x6e, 0x65, 0x72, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0d, 0x63, 0x6f, 0x6e, 0x74, 0x61, 0x69, 0x6e, 0x65, 0x72, 0x4e, 0x61,
	0x6d, 0x65, 0x12, 0x34, 0x0a, 0x15, 0x43, 0x6f, 0x6e, 0x74, 0x61, 0x69, 0x6e, 0x65, 0x72, 0x52,
	0x65, 0x73, 0x74, 0x61, 0x72, 0x74, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x15, 0x63, 0x6f, 0x6e, 0x74, 0x61, 0x69, 0x6e, 0x65, 0x72, 0x52, 0x65, 0x73, 0x74,
	0x61, 0x72, 0x74, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x20, 0x0a, 0x0b, 0x43, 0x72, 0x6f, 0x6e,
	0x4a, 0x6f, 0x62, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x63,
	0x72, 0x6f, 0x6e, 0x4a, 0x6f, 0x62, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x1e, 0x0a, 0x0a, 0x43, 0x72,
	0x6f, 0x6e, 0x4a, 0x6f, 0x62, 0x55, 0x49, 0x44, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a,
	0x63, 0x72, 0x6f, 0x6e, 0x4a, 0x6f, 0x62, 0x55, 0x49, 0x44, 0x12, 0x24, 0x0a, 0x0d, 0x44, 0x61,
	0x65, 0x6d, 0x6f, 0x6e, 0x53, 0x65, 0x74, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0d, 0x64, 0x61, 0x65, 0x6d, 0x6f, 0x6e, 0x53, 0x65, 0x74, 0x4e, 0x61, 0x6d, 0x65,
	0x12, 0x22, 0x0a, 0x0c, 0x44, 0x61, 0x65, 0x6d, 0x6f, 0x6e, 0x53, 0x65, 0x74, 0x55, 0x49, 0x44,
	0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x64, 0x61, 0x65, 0x6d, 0x6f, 0x6e, 0x53, 0x65,
	0x74, 0x55, 0x49, 0x44, 0x12, 0x26, 0x0a, 0x0e, 0x44, 0x65, 0x70, 0x6c, 0x6f, 0x79, 0x6d, 0x65,
	0x6e, 0x74, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x08, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0e, 0x64, 0x65,
	0x70, 0x6c, 0x6f, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x24, 0x0a, 0x0d,
	0x44, 0x65, 0x70, 0x6c, 0x6f, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x55, 0x49, 0x44, 0x18, 0x09, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0d, 0x64, 0x65, 0x70, 0x6c, 0x6f, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x55,
	0x49, 0x44, 0x12, 0x18, 0x0a, 0x07, 0x4a, 0x6f, 0x62, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x0a, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x07, 0x6a, 0x6f, 0x62, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x16, 0x0a, 0x06,
	0x4a, 0x6f, 0x62, 0x55, 0x49, 0x44, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x6a, 0x6f,
	0x62, 0x55, 0x49, 0x44, 0x12, 0x24, 0x0a, 0x0d, 0x4e, 0x61, 0x6d, 0x65, 0x73, 0x70, 0x61, 0x63,
	0x65, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x0c, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x6e, 0x61, 0x6d,
	0x65, 0x73, 0x70, 0x61, 0x63, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x4e, 0x6f,
	0x64, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x0d, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x6e, 0x6f,
	0x64, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x4e, 0x6f, 0x64, 0x65, 0x55, 0x49,
	0x44, 0x18, 0x0e, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6e, 0x6f, 0x64, 0x65, 0x55, 0x49, 0x44,
	0x12, 0x18, 0x0a, 0x07, 0x50, 0x6f, 0x64, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x0f, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x07, 0x70, 0x6f, 0x64, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x50, 0x6f,
	0x64, 0x55, 0x49, 0x44, 0x18, 0x10, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x70, 0x6f, 0x64, 0x55,
	0x49, 0x44, 0x12, 0x26, 0x0a, 0x0e, 0x52, 0x65, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x53, 0x65, 0x74,
	0x4e, 0x61, 0x6d, 0x65, 0x18, 0x11, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0e, 0x72, 0x65, 0x70, 0x6c,
	0x69, 0x63, 0x61, 0x53, 0x65, 0x74, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x24, 0x0a, 0x0d, 0x52, 0x65,
	0x70, 0x6c, 0x69, 0x63, 0x61, 0x53, 0x65, 0x74, 0x55, 0x49, 0x44, 0x18, 0x12, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0d, 0x72, 0x65, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x53, 0x65, 0x74, 0x55, 0x49, 0x44,
	0x12, 0x28, 0x0a, 0x0f, 0x53, 0x74, 0x61, 0x74, 0x65, 0x66, 0x75, 0x6c, 0x53, 0x65, 0x74, 0x4e,
	0x61, 0x6d, 0x65, 0x18, 0x13, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0f, 0x73, 0x74, 0x61, 0x74, 0x65,
	0x66, 0x75, 0x6c, 0x53, 0x65, 0x74, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x26, 0x0a, 0x0e, 0x53, 0x74,
	0x61, 0x74, 0x65, 0x66, 0x75, 0x6c, 0x53, 0x65, 0x74, 0x55, 0x49, 0x44, 0x18, 0x14, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0e, 0x73, 0x74, 0x61, 0x74, 0x65, 0x66, 0x75, 0x6c, 0x53, 0x65, 0x74, 0x55,
	0x49, 0x44, 0x22, 0x91, 0x01, 0x0a, 0x17, 0x43, 0x6f, 0x6e, 0x74, 0x61, 0x69, 0x6e, 0x65, 0x72,
	0x41, 0x74, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x65, 0x73, 0x53, 0x70, 0x65, 0x63, 0x12, 0x0e,
	0x0a, 0x02, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x1c,
	0x0a, 0x09, 0x49, 0x6d, 0x61, 0x67, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x09, 0x69, 0x6d, 0x61, 0x67, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x1a, 0x0a, 0x08,
	0x49, 0x6d, 0x61, 0x67, 0x65, 0x54, 0x61, 0x67, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08,
	0x69, 0x6d, 0x61, 0x67, 0x65, 0x54, 0x61, 0x67, 0x12, 0x12, 0x0a, 0x04, 0x4e, 0x61, 0x6d, 0x65,
	0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x18, 0x0a, 0x07,
	0x52, 0x75, 0x6e, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x72,
	0x75, 0x6e, 0x74, 0x69, 0x6d, 0x65, 0x22, 0xa8, 0x01, 0x0a, 0x12, 0x48, 0x6f, 0x73, 0x74, 0x41,
	0x74, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x65, 0x73, 0x53, 0x70, 0x65, 0x63, 0x12, 0x0e, 0x0a,
	0x02, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x18, 0x0a,
	0x07, 0x49, 0x6d, 0x61, 0x67, 0x65, 0x49, 0x44, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07,
	0x69, 0x6d, 0x61, 0x67, 0x65, 0x49, 0x44, 0x12, 0x1c, 0x0a, 0x09, 0x49, 0x6d, 0x61, 0x67, 0x65,
	0x4e, 0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x69, 0x6d, 0x61, 0x67,
	0x65, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x22, 0x0a, 0x0c, 0x49, 0x6d, 0x61, 0x67, 0x65, 0x56, 0x65,
	0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x69, 0x6d, 0x61,
	0x67, 0x65, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x12, 0x0a, 0x04, 0x4e, 0x61, 0x6d,
	0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x12, 0x0a,
	0x04, 0x54, 0x79, 0x70, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x74, 0x79, 0x70,
	0x65, 0x42, 0x38, 0x5a, 0x36, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f,
	0x61, 0x69, 0x6c, 0x65, 0x72, 0x6f, 0x6e, 0x2d, 0x67, 0x61, 0x74, 0x65, 0x77, 0x61, 0x79, 0x2f,
	0x61, 0x69, 0x6c, 0x65, 0x72, 0x6f, 0x6e, 0x2d, 0x67, 0x61, 0x74, 0x65, 0x77, 0x61, 0x79, 0x2f,
	0x61, 0x70, 0x69, 0x73, 0x2f, 0x61, 0x70, 0x70, 0x2f, 0x76, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
}

var (
	file_app_v1_tracer_tracer_proto_rawDescOnce sync.Once
	file_app_v1_tracer_tracer_proto_rawDescData = file_app_v1_tracer_tracer_proto_rawDesc
)

func file_app_v1_tracer_tracer_proto_rawDescGZIP() []byte {
	file_app_v1_tracer_tracer_proto_rawDescOnce.Do(func() {
		file_app_v1_tracer_tracer_proto_rawDescData = protoimpl.X.CompressGZIP(file_app_v1_tracer_tracer_proto_rawDescData)
	})
	return file_app_v1_tracer_tracer_proto_rawDescData
}

var file_app_v1_tracer_tracer_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_app_v1_tracer_tracer_proto_goTypes = []any{
	(*K8SAttributesSpec)(nil),       // 0: app.v1.K8sAttributesSpec
	(*ContainerAttributesSpec)(nil), // 1: app.v1.ContainerAttributesSpec
	(*HostAttributesSpec)(nil),      // 2: app.v1.HostAttributesSpec
}
var file_app_v1_tracer_tracer_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_app_v1_tracer_tracer_proto_init() }
func file_app_v1_tracer_tracer_proto_init() {
	if File_app_v1_tracer_tracer_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_app_v1_tracer_tracer_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_app_v1_tracer_tracer_proto_goTypes,
		DependencyIndexes: file_app_v1_tracer_tracer_proto_depIdxs,
		MessageInfos:      file_app_v1_tracer_tracer_proto_msgTypes,
	}.Build()
	File_app_v1_tracer_tracer_proto = out.File
	file_app_v1_tracer_tracer_proto_rawDesc = nil
	file_app_v1_tracer_tracer_proto_goTypes = nil
	file_app_v1_tracer_tracer_proto_depIdxs = nil
}
