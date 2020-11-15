// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0-devel
// 	protoc        v3.6.1
// source: transporter/v1/transporter.proto

package transporterpb

import (
	v1 "github.com/ihtkas/farm/account/v1"
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

type Profile struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	User *v1.User `protobuf:"bytes,1,opt,name=user,proto3" json:"user,omitempty"`
}

func (x *Profile) Reset() {
	*x = Profile{}
	if protoimpl.UnsafeEnabled {
		mi := &file_transporter_v1_transporter_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Profile) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Profile) ProtoMessage() {}

func (x *Profile) ProtoReflect() protoreflect.Message {
	mi := &file_transporter_v1_transporter_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Profile.ProtoReflect.Descriptor instead.
func (*Profile) Descriptor() ([]byte, []int) {
	return file_transporter_v1_transporter_proto_rawDescGZIP(), []int{0}
}

func (x *Profile) GetUser() *v1.User {
	if x != nil {
		return x.User
	}
	return nil
}

var File_transporter_v1_transporter_proto protoreflect.FileDescriptor

var file_transporter_v1_transporter_proto_rawDesc = []byte{
	0x0a, 0x20, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x70, 0x6f, 0x72, 0x74, 0x65, 0x72, 0x2f, 0x76, 0x31,
	0x2f, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x70, 0x6f, 0x72, 0x74, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x0e, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x70, 0x6f, 0x72, 0x74, 0x65, 0x72, 0x2e,
	0x76, 0x31, 0x1a, 0x15, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x2f, 0x76, 0x31, 0x2f, 0x75,
	0x73, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x2f, 0x0a, 0x07, 0x50, 0x72, 0x6f,
	0x66, 0x69, 0x6c, 0x65, 0x12, 0x24, 0x0a, 0x04, 0x75, 0x73, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x10, 0x2e, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x2e, 0x76, 0x31, 0x2e,
	0x55, 0x73, 0x65, 0x72, 0x52, 0x04, 0x75, 0x73, 0x65, 0x72, 0x42, 0x35, 0x5a, 0x33, 0x67, 0x69,
	0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x69, 0x68, 0x74, 0x6b, 0x61, 0x73, 0x2f,
	0x66, 0x61, 0x72, 0x6d, 0x2f, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x70, 0x6f, 0x72, 0x74, 0x65, 0x72,
	0x2f, 0x76, 0x31, 0x3b, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x70, 0x6f, 0x72, 0x74, 0x65, 0x72, 0x70,
	0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_transporter_v1_transporter_proto_rawDescOnce sync.Once
	file_transporter_v1_transporter_proto_rawDescData = file_transporter_v1_transporter_proto_rawDesc
)

func file_transporter_v1_transporter_proto_rawDescGZIP() []byte {
	file_transporter_v1_transporter_proto_rawDescOnce.Do(func() {
		file_transporter_v1_transporter_proto_rawDescData = protoimpl.X.CompressGZIP(file_transporter_v1_transporter_proto_rawDescData)
	})
	return file_transporter_v1_transporter_proto_rawDescData
}

var file_transporter_v1_transporter_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_transporter_v1_transporter_proto_goTypes = []interface{}{
	(*Profile)(nil), // 0: transporter.v1.Profile
	(*v1.User)(nil), // 1: account.v1.User
}
var file_transporter_v1_transporter_proto_depIdxs = []int32{
	1, // 0: transporter.v1.Profile.user:type_name -> account.v1.User
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_transporter_v1_transporter_proto_init() }
func file_transporter_v1_transporter_proto_init() {
	if File_transporter_v1_transporter_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_transporter_v1_transporter_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Profile); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_transporter_v1_transporter_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_transporter_v1_transporter_proto_goTypes,
		DependencyIndexes: file_transporter_v1_transporter_proto_depIdxs,
		MessageInfos:      file_transporter_v1_transporter_proto_msgTypes,
	}.Build()
	File_transporter_v1_transporter_proto = out.File
	file_transporter_v1_transporter_proto_rawDesc = nil
	file_transporter_v1_transporter_proto_goTypes = nil
	file_transporter_v1_transporter_proto_depIdxs = nil
}
