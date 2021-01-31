// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.13.0
// source: proto/server_stream.proto

package proto

import (
	proto "github.com/golang/protobuf/proto"
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

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

// 定义发送请求信息
type SimpleReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Para string `protobuf:"bytes,1,opt,name=para,proto3" json:"para,omitempty"`
}

func (x *SimpleReq) Reset() {
	*x = SimpleReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_server_stream_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SimpleReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SimpleReq) ProtoMessage() {}

func (x *SimpleReq) ProtoReflect() protoreflect.Message {
	mi := &file_proto_server_stream_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SimpleReq.ProtoReflect.Descriptor instead.
func (*SimpleReq) Descriptor() ([]byte, []int) {
	return file_proto_server_stream_proto_rawDescGZIP(), []int{0}
}

func (x *SimpleReq) GetPara() string {
	if x != nil {
		return x.Para
	}
	return ""
}

// 定义流式响应信息
type StreamResp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// 流式响应数据
	Val string `protobuf:"bytes,1,opt,name=val,proto3" json:"val,omitempty"`
}

func (x *StreamResp) Reset() {
	*x = StreamResp{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_server_stream_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StreamResp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StreamResp) ProtoMessage() {}

func (x *StreamResp) ProtoReflect() protoreflect.Message {
	mi := &file_proto_server_stream_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StreamResp.ProtoReflect.Descriptor instead.
func (*StreamResp) Descriptor() ([]byte, []int) {
	return file_proto_server_stream_proto_rawDescGZIP(), []int{1}
}

func (x *StreamResp) GetVal() string {
	if x != nil {
		return x.Val
	}
	return ""
}

var File_proto_server_stream_proto protoreflect.FileDescriptor

var file_proto_server_stream_proto_rawDesc = []byte{
	0x0a, 0x19, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x5f, 0x73,
	0x74, 0x72, 0x65, 0x61, 0x6d, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x22, 0x1f, 0x0a, 0x09, 0x53, 0x69, 0x6d, 0x70, 0x6c, 0x65, 0x52, 0x65, 0x71, 0x12,
	0x12, 0x0a, 0x04, 0x70, 0x61, 0x72, 0x61, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x70,
	0x61, 0x72, 0x61, 0x22, 0x1e, 0x0a, 0x0a, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x52, 0x65, 0x73,
	0x70, 0x12, 0x10, 0x0a, 0x03, 0x76, 0x61, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03,
	0x76, 0x61, 0x6c, 0x32, 0x44, 0x0a, 0x0c, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x53, 0x65, 0x72,
	0x76, 0x65, 0x72, 0x12, 0x34, 0x0a, 0x09, 0x4c, 0x69, 0x73, 0x74, 0x56, 0x61, 0x6c, 0x75, 0x65,
	0x12, 0x10, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x53, 0x69, 0x6d, 0x70, 0x6c, 0x65, 0x52,
	0x65, 0x71, 0x1a, 0x11, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x53, 0x74, 0x72, 0x65, 0x61,
	0x6d, 0x52, 0x65, 0x73, 0x70, 0x22, 0x00, 0x30, 0x01, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var (
	file_proto_server_stream_proto_rawDescOnce sync.Once
	file_proto_server_stream_proto_rawDescData = file_proto_server_stream_proto_rawDesc
)

func file_proto_server_stream_proto_rawDescGZIP() []byte {
	file_proto_server_stream_proto_rawDescOnce.Do(func() {
		file_proto_server_stream_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_server_stream_proto_rawDescData)
	})
	return file_proto_server_stream_proto_rawDescData
}

var file_proto_server_stream_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_proto_server_stream_proto_goTypes = []interface{}{
	(*SimpleReq)(nil),  // 0: proto.SimpleReq
	(*StreamResp)(nil), // 1: proto.StreamResp
}
var file_proto_server_stream_proto_depIdxs = []int32{
	0, // 0: proto.StreamServer.ListValue:input_type -> proto.SimpleReq
	1, // 1: proto.StreamServer.ListValue:output_type -> proto.StreamResp
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_proto_server_stream_proto_init() }
func file_proto_server_stream_proto_init() {
	if File_proto_server_stream_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_server_stream_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SimpleReq); i {
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
		file_proto_server_stream_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StreamResp); i {
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
			RawDescriptor: file_proto_server_stream_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_server_stream_proto_goTypes,
		DependencyIndexes: file_proto_server_stream_proto_depIdxs,
		MessageInfos:      file_proto_server_stream_proto_msgTypes,
	}.Build()
	File_proto_server_stream_proto = out.File
	file_proto_server_stream_proto_rawDesc = nil
	file_proto_server_stream_proto_goTypes = nil
	file_proto_server_stream_proto_depIdxs = nil
}
