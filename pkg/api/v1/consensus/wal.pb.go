// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.19.4
// source: pkg/api/v1/consensus/wal.proto

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

package consensus

import (
	types "github.com/bhojpur/state/pkg/api/v1/types"
	_ "github.com/gogo/protobuf/gogoproto"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	durationpb "google.golang.org/protobuf/types/known/durationpb"
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

// MsgInfo are msgs from the reactor which may update the state
type MsgInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Msg    *Message `protobuf:"bytes,1,opt,name=msg,proto3" json:"msg,omitempty"`
	PeerId string   `protobuf:"bytes,2,opt,name=peer_id,json=peerId,proto3" json:"peer_id,omitempty"`
}

func (x *MsgInfo) Reset() {
	*x = MsgInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_api_v1_consensus_wal_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MsgInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MsgInfo) ProtoMessage() {}

func (x *MsgInfo) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_api_v1_consensus_wal_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MsgInfo.ProtoReflect.Descriptor instead.
func (*MsgInfo) Descriptor() ([]byte, []int) {
	return file_pkg_api_v1_consensus_wal_proto_rawDescGZIP(), []int{0}
}

func (x *MsgInfo) GetMsg() *Message {
	if x != nil {
		return x.Msg
	}
	return nil
}

func (x *MsgInfo) GetPeerId() string {
	if x != nil {
		return x.PeerId
	}
	return ""
}

// TimeoutInfo internally generated messages which may update the state
type TimeoutInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Duration *durationpb.Duration `protobuf:"bytes,1,opt,name=duration,proto3" json:"duration,omitempty"`
	Height   int64                `protobuf:"varint,2,opt,name=height,proto3" json:"height,omitempty"`
	Round    int32                `protobuf:"varint,3,opt,name=round,proto3" json:"round,omitempty"`
	Step     uint32               `protobuf:"varint,4,opt,name=step,proto3" json:"step,omitempty"`
}

func (x *TimeoutInfo) Reset() {
	*x = TimeoutInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_api_v1_consensus_wal_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TimeoutInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TimeoutInfo) ProtoMessage() {}

func (x *TimeoutInfo) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_api_v1_consensus_wal_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TimeoutInfo.ProtoReflect.Descriptor instead.
func (*TimeoutInfo) Descriptor() ([]byte, []int) {
	return file_pkg_api_v1_consensus_wal_proto_rawDescGZIP(), []int{1}
}

func (x *TimeoutInfo) GetDuration() *durationpb.Duration {
	if x != nil {
		return x.Duration
	}
	return nil
}

func (x *TimeoutInfo) GetHeight() int64 {
	if x != nil {
		return x.Height
	}
	return 0
}

func (x *TimeoutInfo) GetRound() int32 {
	if x != nil {
		return x.Round
	}
	return 0
}

func (x *TimeoutInfo) GetStep() uint32 {
	if x != nil {
		return x.Step
	}
	return 0
}

// EndHeight marks the end of the given height inside WAL.
// @internal used by scripts/wal2json util.
type EndHeight struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Height int64 `protobuf:"varint,1,opt,name=height,proto3" json:"height,omitempty"`
}

func (x *EndHeight) Reset() {
	*x = EndHeight{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_api_v1_consensus_wal_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EndHeight) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EndHeight) ProtoMessage() {}

func (x *EndHeight) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_api_v1_consensus_wal_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EndHeight.ProtoReflect.Descriptor instead.
func (*EndHeight) Descriptor() ([]byte, []int) {
	return file_pkg_api_v1_consensus_wal_proto_rawDescGZIP(), []int{2}
}

func (x *EndHeight) GetHeight() int64 {
	if x != nil {
		return x.Height
	}
	return 0
}

type WALMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to Sum:
	//	*WALMessage_EventDataRoundState
	//	*WALMessage_MsgInfo
	//	*WALMessage_TimeoutInfo
	//	*WALMessage_EndHeight
	Sum isWALMessage_Sum `protobuf_oneof:"sum"`
}

func (x *WALMessage) Reset() {
	*x = WALMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_api_v1_consensus_wal_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *WALMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*WALMessage) ProtoMessage() {}

func (x *WALMessage) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_api_v1_consensus_wal_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use WALMessage.ProtoReflect.Descriptor instead.
func (*WALMessage) Descriptor() ([]byte, []int) {
	return file_pkg_api_v1_consensus_wal_proto_rawDescGZIP(), []int{3}
}

func (m *WALMessage) GetSum() isWALMessage_Sum {
	if m != nil {
		return m.Sum
	}
	return nil
}

func (x *WALMessage) GetEventDataRoundState() *types.EventDataRoundState {
	if x, ok := x.GetSum().(*WALMessage_EventDataRoundState); ok {
		return x.EventDataRoundState
	}
	return nil
}

func (x *WALMessage) GetMsgInfo() *MsgInfo {
	if x, ok := x.GetSum().(*WALMessage_MsgInfo); ok {
		return x.MsgInfo
	}
	return nil
}

func (x *WALMessage) GetTimeoutInfo() *TimeoutInfo {
	if x, ok := x.GetSum().(*WALMessage_TimeoutInfo); ok {
		return x.TimeoutInfo
	}
	return nil
}

func (x *WALMessage) GetEndHeight() *EndHeight {
	if x, ok := x.GetSum().(*WALMessage_EndHeight); ok {
		return x.EndHeight
	}
	return nil
}

type isWALMessage_Sum interface {
	isWALMessage_Sum()
}

type WALMessage_EventDataRoundState struct {
	EventDataRoundState *types.EventDataRoundState `protobuf:"bytes,1,opt,name=event_data_round_state,json=eventDataRoundState,proto3,oneof"`
}

type WALMessage_MsgInfo struct {
	MsgInfo *MsgInfo `protobuf:"bytes,2,opt,name=msg_info,json=msgInfo,proto3,oneof"`
}

type WALMessage_TimeoutInfo struct {
	TimeoutInfo *TimeoutInfo `protobuf:"bytes,3,opt,name=timeout_info,json=timeoutInfo,proto3,oneof"`
}

type WALMessage_EndHeight struct {
	EndHeight *EndHeight `protobuf:"bytes,4,opt,name=end_height,json=endHeight,proto3,oneof"`
}

func (*WALMessage_EventDataRoundState) isWALMessage_Sum() {}

func (*WALMessage_MsgInfo) isWALMessage_Sum() {}

func (*WALMessage_TimeoutInfo) isWALMessage_Sum() {}

func (*WALMessage_EndHeight) isWALMessage_Sum() {}

// TimedWALMessage wraps WALMessage and adds Time for debugging purposes.
type TimedWALMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Time *timestamppb.Timestamp `protobuf:"bytes,1,opt,name=time,proto3" json:"time,omitempty"`
	Msg  *WALMessage            `protobuf:"bytes,2,opt,name=msg,proto3" json:"msg,omitempty"`
}

func (x *TimedWALMessage) Reset() {
	*x = TimedWALMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_api_v1_consensus_wal_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TimedWALMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TimedWALMessage) ProtoMessage() {}

func (x *TimedWALMessage) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_api_v1_consensus_wal_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TimedWALMessage.ProtoReflect.Descriptor instead.
func (*TimedWALMessage) Descriptor() ([]byte, []int) {
	return file_pkg_api_v1_consensus_wal_proto_rawDescGZIP(), []int{4}
}

func (x *TimedWALMessage) GetTime() *timestamppb.Timestamp {
	if x != nil {
		return x.Time
	}
	return nil
}

func (x *TimedWALMessage) GetMsg() *WALMessage {
	if x != nil {
		return x.Msg
	}
	return nil
}

var File_pkg_api_v1_consensus_wal_proto protoreflect.FileDescriptor

var file_pkg_api_v1_consensus_wal_proto_rawDesc = []byte{
	0x0a, 0x1e, 0x70, 0x6b, 0x67, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x31, 0x2f, 0x63, 0x6f, 0x6e,
	0x73, 0x65, 0x6e, 0x73, 0x75, 0x73, 0x2f, 0x77, 0x61, 0x6c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x12, 0x0c, 0x76, 0x31, 0x2e, 0x63, 0x6f, 0x6e, 0x73, 0x65, 0x6e, 0x73, 0x75, 0x73, 0x1a, 0x14,
	0x67, 0x6f, 0x67, 0x6f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x67, 0x6f, 0x67, 0x6f, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x20, 0x70, 0x6b, 0x67, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x31,
	0x2f, 0x63, 0x6f, 0x6e, 0x73, 0x65, 0x6e, 0x73, 0x75, 0x73, 0x2f, 0x74, 0x79, 0x70, 0x65, 0x73,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1d, 0x70, 0x6b, 0x67, 0x2f, 0x61, 0x70, 0x69, 0x2f,
	0x76, 0x31, 0x2f, 0x74, 0x79, 0x70, 0x65, 0x73, 0x2f, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x64, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x5d, 0x0a, 0x07, 0x4d, 0x73, 0x67, 0x49, 0x6e, 0x66,
	0x6f, 0x12, 0x2d, 0x0a, 0x03, 0x6d, 0x73, 0x67, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x15,
	0x2e, 0x76, 0x31, 0x2e, 0x63, 0x6f, 0x6e, 0x73, 0x65, 0x6e, 0x73, 0x75, 0x73, 0x2e, 0x4d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x42, 0x04, 0xc8, 0xde, 0x1f, 0x00, 0x52, 0x03, 0x6d, 0x73, 0x67,
	0x12, 0x23, 0x0a, 0x07, 0x70, 0x65, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x42, 0x0a, 0xe2, 0xde, 0x1f, 0x06, 0x50, 0x65, 0x65, 0x72, 0x49, 0x44, 0x52, 0x06, 0x70,
	0x65, 0x65, 0x72, 0x49, 0x64, 0x22, 0x90, 0x01, 0x0a, 0x0b, 0x54, 0x69, 0x6d, 0x65, 0x6f, 0x75,
	0x74, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x3f, 0x0a, 0x08, 0x64, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x44, 0x75, 0x72, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x42, 0x08, 0xc8, 0xde, 0x1f, 0x00, 0x98, 0xdf, 0x1f, 0x01, 0x52, 0x08, 0x64, 0x75,
	0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x16, 0x0a, 0x06, 0x68, 0x65, 0x69, 0x67, 0x68, 0x74,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x06, 0x68, 0x65, 0x69, 0x67, 0x68, 0x74, 0x12, 0x14,
	0x0a, 0x05, 0x72, 0x6f, 0x75, 0x6e, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x05, 0x72,
	0x6f, 0x75, 0x6e, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x73, 0x74, 0x65, 0x70, 0x18, 0x04, 0x20, 0x01,
	0x28, 0x0d, 0x52, 0x04, 0x73, 0x74, 0x65, 0x70, 0x22, 0x23, 0x0a, 0x09, 0x45, 0x6e, 0x64, 0x48,
	0x65, 0x69, 0x67, 0x68, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x68, 0x65, 0x69, 0x67, 0x68, 0x74, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x06, 0x68, 0x65, 0x69, 0x67, 0x68, 0x74, 0x22, 0x97, 0x02,
	0x0a, 0x0a, 0x57, 0x41, 0x4c, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x54, 0x0a, 0x16,
	0x65, 0x76, 0x65, 0x6e, 0x74, 0x5f, 0x64, 0x61, 0x74, 0x61, 0x5f, 0x72, 0x6f, 0x75, 0x6e, 0x64,
	0x5f, 0x73, 0x74, 0x61, 0x74, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1d, 0x2e, 0x76,
	0x31, 0x2e, 0x74, 0x79, 0x70, 0x65, 0x73, 0x2e, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x44, 0x61, 0x74,
	0x61, 0x52, 0x6f, 0x75, 0x6e, 0x64, 0x53, 0x74, 0x61, 0x74, 0x65, 0x48, 0x00, 0x52, 0x13, 0x65,
	0x76, 0x65, 0x6e, 0x74, 0x44, 0x61, 0x74, 0x61, 0x52, 0x6f, 0x75, 0x6e, 0x64, 0x53, 0x74, 0x61,
	0x74, 0x65, 0x12, 0x32, 0x0a, 0x08, 0x6d, 0x73, 0x67, 0x5f, 0x69, 0x6e, 0x66, 0x6f, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x15, 0x2e, 0x76, 0x31, 0x2e, 0x63, 0x6f, 0x6e, 0x73, 0x65, 0x6e,
	0x73, 0x75, 0x73, 0x2e, 0x4d, 0x73, 0x67, 0x49, 0x6e, 0x66, 0x6f, 0x48, 0x00, 0x52, 0x07, 0x6d,
	0x73, 0x67, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x3e, 0x0a, 0x0c, 0x74, 0x69, 0x6d, 0x65, 0x6f, 0x75,
	0x74, 0x5f, 0x69, 0x6e, 0x66, 0x6f, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x76,
	0x31, 0x2e, 0x63, 0x6f, 0x6e, 0x73, 0x65, 0x6e, 0x73, 0x75, 0x73, 0x2e, 0x54, 0x69, 0x6d, 0x65,
	0x6f, 0x75, 0x74, 0x49, 0x6e, 0x66, 0x6f, 0x48, 0x00, 0x52, 0x0b, 0x74, 0x69, 0x6d, 0x65, 0x6f,
	0x75, 0x74, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x38, 0x0a, 0x0a, 0x65, 0x6e, 0x64, 0x5f, 0x68, 0x65,
	0x69, 0x67, 0x68, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x76, 0x31, 0x2e,
	0x63, 0x6f, 0x6e, 0x73, 0x65, 0x6e, 0x73, 0x75, 0x73, 0x2e, 0x45, 0x6e, 0x64, 0x48, 0x65, 0x69,
	0x67, 0x68, 0x74, 0x48, 0x00, 0x52, 0x09, 0x65, 0x6e, 0x64, 0x48, 0x65, 0x69, 0x67, 0x68, 0x74,
	0x42, 0x05, 0x0a, 0x03, 0x73, 0x75, 0x6d, 0x22, 0x77, 0x0a, 0x0f, 0x54, 0x69, 0x6d, 0x65, 0x64,
	0x57, 0x41, 0x4c, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x38, 0x0a, 0x04, 0x74, 0x69,
	0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73,
	0x74, 0x61, 0x6d, 0x70, 0x42, 0x08, 0xc8, 0xde, 0x1f, 0x00, 0x90, 0xdf, 0x1f, 0x01, 0x52, 0x04,
	0x74, 0x69, 0x6d, 0x65, 0x12, 0x2a, 0x0a, 0x03, 0x6d, 0x73, 0x67, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x18, 0x2e, 0x76, 0x31, 0x2e, 0x63, 0x6f, 0x6e, 0x73, 0x65, 0x6e, 0x73, 0x75, 0x73,
	0x2e, 0x57, 0x41, 0x4c, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x52, 0x03, 0x6d, 0x73, 0x67,
	0x42, 0x39, 0x5a, 0x37, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x62,
	0x68, 0x6f, 0x6a, 0x70, 0x75, 0x72, 0x2f, 0x73, 0x74, 0x61, 0x74, 0x65, 0x2f, 0x70, 0x6b, 0x67,
	0x2f, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x31, 0x2f, 0x63, 0x6f, 0x6e, 0x73, 0x65, 0x6e, 0x73, 0x75,
	0x73, 0x3b, 0x63, 0x6f, 0x6e, 0x73, 0x65, 0x6e, 0x73, 0x75, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
}

var (
	file_pkg_api_v1_consensus_wal_proto_rawDescOnce sync.Once
	file_pkg_api_v1_consensus_wal_proto_rawDescData = file_pkg_api_v1_consensus_wal_proto_rawDesc
)

func file_pkg_api_v1_consensus_wal_proto_rawDescGZIP() []byte {
	file_pkg_api_v1_consensus_wal_proto_rawDescOnce.Do(func() {
		file_pkg_api_v1_consensus_wal_proto_rawDescData = protoimpl.X.CompressGZIP(file_pkg_api_v1_consensus_wal_proto_rawDescData)
	})
	return file_pkg_api_v1_consensus_wal_proto_rawDescData
}

var file_pkg_api_v1_consensus_wal_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_pkg_api_v1_consensus_wal_proto_goTypes = []interface{}{
	(*MsgInfo)(nil),                   // 0: v1.consensus.MsgInfo
	(*TimeoutInfo)(nil),               // 1: v1.consensus.TimeoutInfo
	(*EndHeight)(nil),                 // 2: v1.consensus.EndHeight
	(*WALMessage)(nil),                // 3: v1.consensus.WALMessage
	(*TimedWALMessage)(nil),           // 4: v1.consensus.TimedWALMessage
	(*Message)(nil),                   // 5: v1.consensus.Message
	(*durationpb.Duration)(nil),       // 6: google.protobuf.Duration
	(*types.EventDataRoundState)(nil), // 7: v1.types.EventDataRoundState
	(*timestamppb.Timestamp)(nil),     // 8: google.protobuf.Timestamp
}
var file_pkg_api_v1_consensus_wal_proto_depIdxs = []int32{
	5, // 0: v1.consensus.MsgInfo.msg:type_name -> v1.consensus.Message
	6, // 1: v1.consensus.TimeoutInfo.duration:type_name -> google.protobuf.Duration
	7, // 2: v1.consensus.WALMessage.event_data_round_state:type_name -> v1.types.EventDataRoundState
	0, // 3: v1.consensus.WALMessage.msg_info:type_name -> v1.consensus.MsgInfo
	1, // 4: v1.consensus.WALMessage.timeout_info:type_name -> v1.consensus.TimeoutInfo
	2, // 5: v1.consensus.WALMessage.end_height:type_name -> v1.consensus.EndHeight
	8, // 6: v1.consensus.TimedWALMessage.time:type_name -> google.protobuf.Timestamp
	3, // 7: v1.consensus.TimedWALMessage.msg:type_name -> v1.consensus.WALMessage
	8, // [8:8] is the sub-list for method output_type
	8, // [8:8] is the sub-list for method input_type
	8, // [8:8] is the sub-list for extension type_name
	8, // [8:8] is the sub-list for extension extendee
	0, // [0:8] is the sub-list for field type_name
}

func init() { file_pkg_api_v1_consensus_wal_proto_init() }
func file_pkg_api_v1_consensus_wal_proto_init() {
	if File_pkg_api_v1_consensus_wal_proto != nil {
		return
	}
	file_pkg_api_v1_consensus_types_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_pkg_api_v1_consensus_wal_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MsgInfo); i {
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
		file_pkg_api_v1_consensus_wal_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TimeoutInfo); i {
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
		file_pkg_api_v1_consensus_wal_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*EndHeight); i {
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
		file_pkg_api_v1_consensus_wal_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*WALMessage); i {
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
		file_pkg_api_v1_consensus_wal_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TimedWALMessage); i {
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
	file_pkg_api_v1_consensus_wal_proto_msgTypes[3].OneofWrappers = []interface{}{
		(*WALMessage_EventDataRoundState)(nil),
		(*WALMessage_MsgInfo)(nil),
		(*WALMessage_TimeoutInfo)(nil),
		(*WALMessage_EndHeight)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_pkg_api_v1_consensus_wal_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_pkg_api_v1_consensus_wal_proto_goTypes,
		DependencyIndexes: file_pkg_api_v1_consensus_wal_proto_depIdxs,
		MessageInfos:      file_pkg_api_v1_consensus_wal_proto_msgTypes,
	}.Build()
	File_pkg_api_v1_consensus_wal_proto = out.File
	file_pkg_api_v1_consensus_wal_proto_rawDesc = nil
	file_pkg_api_v1_consensus_wal_proto_goTypes = nil
	file_pkg_api_v1_consensus_wal_proto_depIdxs = nil
}