// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.19.4
// source: pkg/api/v1/types/canonical.proto

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

package types

import (
	_ "github.com/gogo/protobuf/gogoproto"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type CanonicalBlockID struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Hash          []byte                  `protobuf:"bytes,1,opt,name=hash,proto3" json:"hash,omitempty"`
	PartSetHeader *CanonicalPartSetHeader `protobuf:"bytes,2,opt,name=part_set_header,json=partSetHeader,proto3" json:"part_set_header,omitempty"`
}

func (x *CanonicalBlockID) Reset() {
	*x = CanonicalBlockID{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_api_v1_types_canonical_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CanonicalBlockID) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CanonicalBlockID) ProtoMessage() {}

func (x *CanonicalBlockID) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_api_v1_types_canonical_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CanonicalBlockID.ProtoReflect.Descriptor instead.
func (*CanonicalBlockID) Descriptor() ([]byte, []int) {
	return file_pkg_api_v1_types_canonical_proto_rawDescGZIP(), []int{0}
}

func (x *CanonicalBlockID) GetHash() []byte {
	if x != nil {
		return x.Hash
	}
	return nil
}

func (x *CanonicalBlockID) GetPartSetHeader() *CanonicalPartSetHeader {
	if x != nil {
		return x.PartSetHeader
	}
	return nil
}

type CanonicalPartSetHeader struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Total uint32 `protobuf:"varint,1,opt,name=total,proto3" json:"total,omitempty"`
	Hash  []byte `protobuf:"bytes,2,opt,name=hash,proto3" json:"hash,omitempty"`
}

func (x *CanonicalPartSetHeader) Reset() {
	*x = CanonicalPartSetHeader{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_api_v1_types_canonical_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CanonicalPartSetHeader) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CanonicalPartSetHeader) ProtoMessage() {}

func (x *CanonicalPartSetHeader) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_api_v1_types_canonical_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CanonicalPartSetHeader.ProtoReflect.Descriptor instead.
func (*CanonicalPartSetHeader) Descriptor() ([]byte, []int) {
	return file_pkg_api_v1_types_canonical_proto_rawDescGZIP(), []int{1}
}

func (x *CanonicalPartSetHeader) GetTotal() uint32 {
	if x != nil {
		return x.Total
	}
	return 0
}

func (x *CanonicalPartSetHeader) GetHash() []byte {
	if x != nil {
		return x.Hash
	}
	return nil
}

type CanonicalProposal struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Type      SignedMsgType          `protobuf:"varint,1,opt,name=type,proto3,enum=v1.types.SignedMsgType" json:"type,omitempty"` // type alias for byte
	Height    int64                  `protobuf:"fixed64,2,opt,name=height,proto3" json:"height,omitempty"`                        // canonicalization requires fixed size encoding here
	Round     int64                  `protobuf:"fixed64,3,opt,name=round,proto3" json:"round,omitempty"`                          // canonicalization requires fixed size encoding here
	PolRound  int64                  `protobuf:"varint,4,opt,name=pol_round,json=polRound,proto3" json:"pol_round,omitempty"`
	BlockId   *CanonicalBlockID      `protobuf:"bytes,5,opt,name=block_id,json=blockId,proto3" json:"block_id,omitempty"`
	Timestamp *timestamppb.Timestamp `protobuf:"bytes,6,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	ChainId   string                 `protobuf:"bytes,7,opt,name=chain_id,json=chainId,proto3" json:"chain_id,omitempty"`
}

func (x *CanonicalProposal) Reset() {
	*x = CanonicalProposal{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_api_v1_types_canonical_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CanonicalProposal) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CanonicalProposal) ProtoMessage() {}

func (x *CanonicalProposal) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_api_v1_types_canonical_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CanonicalProposal.ProtoReflect.Descriptor instead.
func (*CanonicalProposal) Descriptor() ([]byte, []int) {
	return file_pkg_api_v1_types_canonical_proto_rawDescGZIP(), []int{2}
}

func (x *CanonicalProposal) GetType() SignedMsgType {
	if x != nil {
		return x.Type
	}
	return SignedMsgType_SIGNED_MSG_TYPE_UNKNOWN
}

func (x *CanonicalProposal) GetHeight() int64 {
	if x != nil {
		return x.Height
	}
	return 0
}

func (x *CanonicalProposal) GetRound() int64 {
	if x != nil {
		return x.Round
	}
	return 0
}

func (x *CanonicalProposal) GetPolRound() int64 {
	if x != nil {
		return x.PolRound
	}
	return 0
}

func (x *CanonicalProposal) GetBlockId() *CanonicalBlockID {
	if x != nil {
		return x.BlockId
	}
	return nil
}

func (x *CanonicalProposal) GetTimestamp() *timestamppb.Timestamp {
	if x != nil {
		return x.Timestamp
	}
	return nil
}

func (x *CanonicalProposal) GetChainId() string {
	if x != nil {
		return x.ChainId
	}
	return ""
}

type CanonicalVote struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Type      SignedMsgType          `protobuf:"varint,1,opt,name=type,proto3,enum=v1.types.SignedMsgType" json:"type,omitempty"` // type alias for byte
	Height    int64                  `protobuf:"fixed64,2,opt,name=height,proto3" json:"height,omitempty"`                        // canonicalization requires fixed size encoding here
	Round     int64                  `protobuf:"fixed64,3,opt,name=round,proto3" json:"round,omitempty"`                          // canonicalization requires fixed size encoding here
	BlockId   *CanonicalBlockID      `protobuf:"bytes,4,opt,name=block_id,json=blockId,proto3" json:"block_id,omitempty"`
	Timestamp *timestamppb.Timestamp `protobuf:"bytes,5,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	ChainId   string                 `protobuf:"bytes,6,opt,name=chain_id,json=chainId,proto3" json:"chain_id,omitempty"`
}

func (x *CanonicalVote) Reset() {
	*x = CanonicalVote{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_api_v1_types_canonical_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CanonicalVote) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CanonicalVote) ProtoMessage() {}

func (x *CanonicalVote) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_api_v1_types_canonical_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CanonicalVote.ProtoReflect.Descriptor instead.
func (*CanonicalVote) Descriptor() ([]byte, []int) {
	return file_pkg_api_v1_types_canonical_proto_rawDescGZIP(), []int{3}
}

func (x *CanonicalVote) GetType() SignedMsgType {
	if x != nil {
		return x.Type
	}
	return SignedMsgType_SIGNED_MSG_TYPE_UNKNOWN
}

func (x *CanonicalVote) GetHeight() int64 {
	if x != nil {
		return x.Height
	}
	return 0
}

func (x *CanonicalVote) GetRound() int64 {
	if x != nil {
		return x.Round
	}
	return 0
}

func (x *CanonicalVote) GetBlockId() *CanonicalBlockID {
	if x != nil {
		return x.BlockId
	}
	return nil
}

func (x *CanonicalVote) GetTimestamp() *timestamppb.Timestamp {
	if x != nil {
		return x.Timestamp
	}
	return nil
}

func (x *CanonicalVote) GetChainId() string {
	if x != nil {
		return x.ChainId
	}
	return ""
}

// CanonicalVoteExtension provides us a way to serialize a vote extension from
// a particular validator such that we can sign over those serialized bytes.
type CanonicalVoteExtension struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Extension []byte `protobuf:"bytes,1,opt,name=extension,proto3" json:"extension,omitempty"`
	Height    int64  `protobuf:"fixed64,2,opt,name=height,proto3" json:"height,omitempty"`
	Round     int64  `protobuf:"fixed64,3,opt,name=round,proto3" json:"round,omitempty"`
	ChainId   string `protobuf:"bytes,4,opt,name=chain_id,json=chainId,proto3" json:"chain_id,omitempty"`
}

func (x *CanonicalVoteExtension) Reset() {
	*x = CanonicalVoteExtension{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_api_v1_types_canonical_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CanonicalVoteExtension) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CanonicalVoteExtension) ProtoMessage() {}

func (x *CanonicalVoteExtension) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_api_v1_types_canonical_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CanonicalVoteExtension.ProtoReflect.Descriptor instead.
func (*CanonicalVoteExtension) Descriptor() ([]byte, []int) {
	return file_pkg_api_v1_types_canonical_proto_rawDescGZIP(), []int{4}
}

func (x *CanonicalVoteExtension) GetExtension() []byte {
	if x != nil {
		return x.Extension
	}
	return nil
}

func (x *CanonicalVoteExtension) GetHeight() int64 {
	if x != nil {
		return x.Height
	}
	return 0
}

func (x *CanonicalVoteExtension) GetRound() int64 {
	if x != nil {
		return x.Round
	}
	return 0
}

func (x *CanonicalVoteExtension) GetChainId() string {
	if x != nil {
		return x.ChainId
	}
	return ""
}

var File_pkg_api_v1_types_canonical_proto protoreflect.FileDescriptor

var file_pkg_api_v1_types_canonical_proto_rawDesc = []byte{
	0x0a, 0x20, 0x70, 0x6b, 0x67, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x31, 0x2f, 0x74, 0x79, 0x70,
	0x65, 0x73, 0x2f, 0x63, 0x61, 0x6e, 0x6f, 0x6e, 0x69, 0x63, 0x61, 0x6c, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x08, 0x76, 0x31, 0x2e, 0x74, 0x79, 0x70, 0x65, 0x73, 0x1a, 0x14, 0x67, 0x6f,
	0x67, 0x6f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x67, 0x6f, 0x67, 0x6f, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x1a, 0x1c, 0x70, 0x6b, 0x67, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x31, 0x2f, 0x74,
	0x79, 0x70, 0x65, 0x73, 0x2f, 0x74, 0x79, 0x70, 0x65, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x22, 0x76, 0x0a, 0x10, 0x43, 0x61, 0x6e, 0x6f, 0x6e, 0x69, 0x63, 0x61, 0x6c, 0x42, 0x6c,
	0x6f, 0x63, 0x6b, 0x49, 0x44, 0x12, 0x12, 0x0a, 0x04, 0x68, 0x61, 0x73, 0x68, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0c, 0x52, 0x04, 0x68, 0x61, 0x73, 0x68, 0x12, 0x4e, 0x0a, 0x0f, 0x70, 0x61, 0x72,
	0x74, 0x5f, 0x73, 0x65, 0x74, 0x5f, 0x68, 0x65, 0x61, 0x64, 0x65, 0x72, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x20, 0x2e, 0x76, 0x31, 0x2e, 0x74, 0x79, 0x70, 0x65, 0x73, 0x2e, 0x43, 0x61,
	0x6e, 0x6f, 0x6e, 0x69, 0x63, 0x61, 0x6c, 0x50, 0x61, 0x72, 0x74, 0x53, 0x65, 0x74, 0x48, 0x65,
	0x61, 0x64, 0x65, 0x72, 0x42, 0x04, 0xc8, 0xde, 0x1f, 0x00, 0x52, 0x0d, 0x70, 0x61, 0x72, 0x74,
	0x53, 0x65, 0x74, 0x48, 0x65, 0x61, 0x64, 0x65, 0x72, 0x22, 0x42, 0x0a, 0x16, 0x43, 0x61, 0x6e,
	0x6f, 0x6e, 0x69, 0x63, 0x61, 0x6c, 0x50, 0x61, 0x72, 0x74, 0x53, 0x65, 0x74, 0x48, 0x65, 0x61,
	0x64, 0x65, 0x72, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0d, 0x52, 0x05, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x12, 0x12, 0x0a, 0x04, 0x68, 0x61, 0x73,
	0x68, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x68, 0x61, 0x73, 0x68, 0x22, 0xc9, 0x02,
	0x0a, 0x11, 0x43, 0x61, 0x6e, 0x6f, 0x6e, 0x69, 0x63, 0x61, 0x6c, 0x50, 0x72, 0x6f, 0x70, 0x6f,
	0x73, 0x61, 0x6c, 0x12, 0x2b, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0e, 0x32, 0x17, 0x2e, 0x76, 0x31, 0x2e, 0x74, 0x79, 0x70, 0x65, 0x73, 0x2e, 0x53, 0x69, 0x67,
	0x6e, 0x65, 0x64, 0x4d, 0x73, 0x67, 0x54, 0x79, 0x70, 0x65, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65,
	0x12, 0x16, 0x0a, 0x06, 0x68, 0x65, 0x69, 0x67, 0x68, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x10,
	0x52, 0x06, 0x68, 0x65, 0x69, 0x67, 0x68, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x72, 0x6f, 0x75, 0x6e,
	0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x10, 0x52, 0x05, 0x72, 0x6f, 0x75, 0x6e, 0x64, 0x12, 0x29,
	0x0a, 0x09, 0x70, 0x6f, 0x6c, 0x5f, 0x72, 0x6f, 0x75, 0x6e, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x03, 0x42, 0x0c, 0xe2, 0xde, 0x1f, 0x08, 0x50, 0x4f, 0x4c, 0x52, 0x6f, 0x75, 0x6e, 0x64, 0x52,
	0x08, 0x70, 0x6f, 0x6c, 0x52, 0x6f, 0x75, 0x6e, 0x64, 0x12, 0x42, 0x0a, 0x08, 0x62, 0x6c, 0x6f,
	0x63, 0x6b, 0x5f, 0x69, 0x64, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x76, 0x31,
	0x2e, 0x74, 0x79, 0x70, 0x65, 0x73, 0x2e, 0x43, 0x61, 0x6e, 0x6f, 0x6e, 0x69, 0x63, 0x61, 0x6c,
	0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x49, 0x44, 0x42, 0x0b, 0xe2, 0xde, 0x1f, 0x07, 0x42, 0x6c, 0x6f,
	0x63, 0x6b, 0x49, 0x44, 0x52, 0x07, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x49, 0x64, 0x12, 0x42, 0x0a,
	0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x42, 0x08, 0xc8, 0xde,
	0x1f, 0x00, 0x90, 0xdf, 0x1f, 0x01, 0x52, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d,
	0x70, 0x12, 0x26, 0x0a, 0x08, 0x63, 0x68, 0x61, 0x69, 0x6e, 0x5f, 0x69, 0x64, 0x18, 0x07, 0x20,
	0x01, 0x28, 0x09, 0x42, 0x0b, 0xe2, 0xde, 0x1f, 0x07, 0x43, 0x68, 0x61, 0x69, 0x6e, 0x49, 0x44,
	0x52, 0x07, 0x63, 0x68, 0x61, 0x69, 0x6e, 0x49, 0x64, 0x22, 0x9a, 0x02, 0x0a, 0x0d, 0x43, 0x61,
	0x6e, 0x6f, 0x6e, 0x69, 0x63, 0x61, 0x6c, 0x56, 0x6f, 0x74, 0x65, 0x12, 0x2b, 0x0a, 0x04, 0x74,
	0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x17, 0x2e, 0x76, 0x31, 0x2e, 0x74,
	0x79, 0x70, 0x65, 0x73, 0x2e, 0x53, 0x69, 0x67, 0x6e, 0x65, 0x64, 0x4d, 0x73, 0x67, 0x54, 0x79,
	0x70, 0x65, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x68, 0x65, 0x69, 0x67,
	0x68, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x10, 0x52, 0x06, 0x68, 0x65, 0x69, 0x67, 0x68, 0x74,
	0x12, 0x14, 0x0a, 0x05, 0x72, 0x6f, 0x75, 0x6e, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x10, 0x52,
	0x05, 0x72, 0x6f, 0x75, 0x6e, 0x64, 0x12, 0x42, 0x0a, 0x08, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x5f,
	0x69, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x76, 0x31, 0x2e, 0x74, 0x79,
	0x70, 0x65, 0x73, 0x2e, 0x43, 0x61, 0x6e, 0x6f, 0x6e, 0x69, 0x63, 0x61, 0x6c, 0x42, 0x6c, 0x6f,
	0x63, 0x6b, 0x49, 0x44, 0x42, 0x0b, 0xe2, 0xde, 0x1f, 0x07, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x49,
	0x44, 0x52, 0x07, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x49, 0x64, 0x12, 0x42, 0x0a, 0x09, 0x74, 0x69,
	0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e,
	0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x42, 0x08, 0xc8, 0xde, 0x1f, 0x00, 0x90,
	0xdf, 0x1f, 0x01, 0x52, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x12, 0x26,
	0x0a, 0x08, 0x63, 0x68, 0x61, 0x69, 0x6e, 0x5f, 0x69, 0x64, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09,
	0x42, 0x0b, 0xe2, 0xde, 0x1f, 0x07, 0x43, 0x68, 0x61, 0x69, 0x6e, 0x49, 0x44, 0x52, 0x07, 0x63,
	0x68, 0x61, 0x69, 0x6e, 0x49, 0x64, 0x22, 0x7f, 0x0a, 0x16, 0x43, 0x61, 0x6e, 0x6f, 0x6e, 0x69,
	0x63, 0x61, 0x6c, 0x56, 0x6f, 0x74, 0x65, 0x45, 0x78, 0x74, 0x65, 0x6e, 0x73, 0x69, 0x6f, 0x6e,
	0x12, 0x1c, 0x0a, 0x09, 0x65, 0x78, 0x74, 0x65, 0x6e, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0c, 0x52, 0x09, 0x65, 0x78, 0x74, 0x65, 0x6e, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x16,
	0x0a, 0x06, 0x68, 0x65, 0x69, 0x67, 0x68, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x10, 0x52, 0x06,
	0x68, 0x65, 0x69, 0x67, 0x68, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x72, 0x6f, 0x75, 0x6e, 0x64, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x10, 0x52, 0x05, 0x72, 0x6f, 0x75, 0x6e, 0x64, 0x12, 0x19, 0x0a, 0x08,
	0x63, 0x68, 0x61, 0x69, 0x6e, 0x5f, 0x69, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07,
	0x63, 0x68, 0x61, 0x69, 0x6e, 0x49, 0x64, 0x42, 0x31, 0x5a, 0x2f, 0x67, 0x69, 0x74, 0x68, 0x75,
	0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x62, 0x68, 0x6f, 0x6a, 0x70, 0x75, 0x72, 0x2f, 0x73, 0x74,
	0x61, 0x74, 0x65, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x31, 0x2f, 0x74,
	0x79, 0x70, 0x65, 0x73, 0x3b, 0x74, 0x79, 0x70, 0x65, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_pkg_api_v1_types_canonical_proto_rawDescOnce sync.Once
	file_pkg_api_v1_types_canonical_proto_rawDescData = file_pkg_api_v1_types_canonical_proto_rawDesc
)

func file_pkg_api_v1_types_canonical_proto_rawDescGZIP() []byte {
	file_pkg_api_v1_types_canonical_proto_rawDescOnce.Do(func() {
		file_pkg_api_v1_types_canonical_proto_rawDescData = protoimpl.X.CompressGZIP(file_pkg_api_v1_types_canonical_proto_rawDescData)
	})
	return file_pkg_api_v1_types_canonical_proto_rawDescData
}

var file_pkg_api_v1_types_canonical_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_pkg_api_v1_types_canonical_proto_goTypes = []interface{}{
	(*CanonicalBlockID)(nil),       // 0: v1.types.CanonicalBlockID
	(*CanonicalPartSetHeader)(nil), // 1: v1.types.CanonicalPartSetHeader
	(*CanonicalProposal)(nil),      // 2: v1.types.CanonicalProposal
	(*CanonicalVote)(nil),          // 3: v1.types.CanonicalVote
	(*CanonicalVoteExtension)(nil), // 4: v1.types.CanonicalVoteExtension
	(SignedMsgType)(0),             // 5: v1.types.SignedMsgType
	(*timestamppb.Timestamp)(nil),  // 6: google.protobuf.Timestamp
}
var file_pkg_api_v1_types_canonical_proto_depIdxs = []int32{
	1, // 0: v1.types.CanonicalBlockID.part_set_header:type_name -> v1.types.CanonicalPartSetHeader
	5, // 1: v1.types.CanonicalProposal.type:type_name -> v1.types.SignedMsgType
	0, // 2: v1.types.CanonicalProposal.block_id:type_name -> v1.types.CanonicalBlockID
	6, // 3: v1.types.CanonicalProposal.timestamp:type_name -> google.protobuf.Timestamp
	5, // 4: v1.types.CanonicalVote.type:type_name -> v1.types.SignedMsgType
	0, // 5: v1.types.CanonicalVote.block_id:type_name -> v1.types.CanonicalBlockID
	6, // 6: v1.types.CanonicalVote.timestamp:type_name -> google.protobuf.Timestamp
	7, // [7:7] is the sub-list for method output_type
	7, // [7:7] is the sub-list for method input_type
	7, // [7:7] is the sub-list for extension type_name
	7, // [7:7] is the sub-list for extension extendee
	0, // [0:7] is the sub-list for field type_name
}

func init() { file_pkg_api_v1_types_canonical_proto_init() }
func file_pkg_api_v1_types_canonical_proto_init() {
	if File_pkg_api_v1_types_canonical_proto != nil {
		return
	}
	file_pkg_api_v1_types_types_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_pkg_api_v1_types_canonical_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CanonicalBlockID); i {
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
		file_pkg_api_v1_types_canonical_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CanonicalPartSetHeader); i {
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
		file_pkg_api_v1_types_canonical_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CanonicalProposal); i {
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
		file_pkg_api_v1_types_canonical_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CanonicalVote); i {
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
		file_pkg_api_v1_types_canonical_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CanonicalVoteExtension); i {
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
			RawDescriptor: file_pkg_api_v1_types_canonical_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_pkg_api_v1_types_canonical_proto_goTypes,
		DependencyIndexes: file_pkg_api_v1_types_canonical_proto_depIdxs,
		MessageInfos:      file_pkg_api_v1_types_canonical_proto_msgTypes,
	}.Build()
	File_pkg_api_v1_types_canonical_proto = out.File
	file_pkg_api_v1_types_canonical_proto_rawDesc = nil
	file_pkg_api_v1_types_canonical_proto_goTypes = nil
	file_pkg_api_v1_types_canonical_proto_depIdxs = nil
}
