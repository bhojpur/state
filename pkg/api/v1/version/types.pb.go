// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.19.4
// source: pkg/api/v1/version/types.proto

// Copyright (c) 2018 Bhojpur Consulting Private Limited, India. All rights reserved.

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package version

import (
	_ "github.com/gogo/protobuf/gogoproto"
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

// Consensus captures the consensus rules for processing a block in the
// blockchain, including all blockchain data structures and the rules of the
// application's state transition machine.
type Consensus struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Block uint64 `protobuf:"varint,1,opt,name=block,proto3" json:"block,omitempty"`
	App   uint64 `protobuf:"varint,2,opt,name=app,proto3" json:"app,omitempty"`
}

func (x *Consensus) Reset() {
	*x = Consensus{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_api_v1_version_types_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Consensus) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Consensus) ProtoMessage() {}

func (x *Consensus) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_api_v1_version_types_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Consensus.ProtoReflect.Descriptor instead.
func (*Consensus) Descriptor() ([]byte, []int) {
	return file_pkg_api_v1_version_types_proto_rawDescGZIP(), []int{0}
}

func (x *Consensus) GetBlock() uint64 {
	if x != nil {
		return x.Block
	}
	return 0
}

func (x *Consensus) GetApp() uint64 {
	if x != nil {
		return x.App
	}
	return 0
}

var File_pkg_api_v1_version_types_proto protoreflect.FileDescriptor

var file_pkg_api_v1_version_types_proto_rawDesc = []byte{
	0x0a, 0x1e, 0x70, 0x6b, 0x67, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x31, 0x2f, 0x76, 0x65, 0x72,
	0x73, 0x69, 0x6f, 0x6e, 0x2f, 0x74, 0x79, 0x70, 0x65, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x12, 0x0a, 0x76, 0x31, 0x2e, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x1a, 0x14, 0x67, 0x6f,
	0x67, 0x6f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x67, 0x6f, 0x67, 0x6f, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x22, 0x39, 0x0a, 0x09, 0x43, 0x6f, 0x6e, 0x73, 0x65, 0x6e, 0x73, 0x75, 0x73, 0x12,
	0x14, 0x0a, 0x05, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x05,
	0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x12, 0x10, 0x0a, 0x03, 0x61, 0x70, 0x70, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x04, 0x52, 0x03, 0x61, 0x70, 0x70, 0x3a, 0x04, 0xe8, 0xa0, 0x1f, 0x01, 0x42, 0x35, 0x5a,
	0x33, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x62, 0x68, 0x6f, 0x6a,
	0x70, 0x75, 0x72, 0x2f, 0x73, 0x74, 0x61, 0x74, 0x65, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x61, 0x70,
	0x69, 0x2f, 0x76, 0x31, 0x2f, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x3b, 0x76, 0x65, 0x72,
	0x73, 0x69, 0x6f, 0x6e, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_pkg_api_v1_version_types_proto_rawDescOnce sync.Once
	file_pkg_api_v1_version_types_proto_rawDescData = file_pkg_api_v1_version_types_proto_rawDesc
)

func file_pkg_api_v1_version_types_proto_rawDescGZIP() []byte {
	file_pkg_api_v1_version_types_proto_rawDescOnce.Do(func() {
		file_pkg_api_v1_version_types_proto_rawDescData = protoimpl.X.CompressGZIP(file_pkg_api_v1_version_types_proto_rawDescData)
	})
	return file_pkg_api_v1_version_types_proto_rawDescData
}

var file_pkg_api_v1_version_types_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_pkg_api_v1_version_types_proto_goTypes = []interface{}{
	(*Consensus)(nil), // 0: v1.version.Consensus
}
var file_pkg_api_v1_version_types_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_pkg_api_v1_version_types_proto_init() }
func file_pkg_api_v1_version_types_proto_init() {
	if File_pkg_api_v1_version_types_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_pkg_api_v1_version_types_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Consensus); i {
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
			RawDescriptor: file_pkg_api_v1_version_types_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_pkg_api_v1_version_types_proto_goTypes,
		DependencyIndexes: file_pkg_api_v1_version_types_proto_depIdxs,
		MessageInfos:      file_pkg_api_v1_version_types_proto_msgTypes,
	}.Build()
	File_pkg_api_v1_version_types_proto = out.File
	file_pkg_api_v1_version_types_proto_rawDesc = nil
	file_pkg_api_v1_version_types_proto_goTypes = nil
	file_pkg_api_v1_version_types_proto_depIdxs = nil
}
