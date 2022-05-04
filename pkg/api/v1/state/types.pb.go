// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.19.4
// source: pkg/api/v1/state/types.proto

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

package state

import (
	types "github.com/bhojpur/state/pkg/abci/types"
	types1 "github.com/bhojpur/state/pkg/api/v1/types"
	version "github.com/bhojpur/state/pkg/api/v1/version"
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

// ABCIResponses retains the responses
// of the various ABCI calls during block processing.
// It is persisted to disk for each height before calling Commit.
type ABCIResponses struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	FinalizeBlock *types.ResponseFinalizeBlock `protobuf:"bytes,2,opt,name=finalize_block,json=finalizeBlock,proto3" json:"finalize_block,omitempty"`
}

func (x *ABCIResponses) Reset() {
	*x = ABCIResponses{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_api_v1_state_types_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ABCIResponses) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ABCIResponses) ProtoMessage() {}

func (x *ABCIResponses) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_api_v1_state_types_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ABCIResponses.ProtoReflect.Descriptor instead.
func (*ABCIResponses) Descriptor() ([]byte, []int) {
	return file_pkg_api_v1_state_types_proto_rawDescGZIP(), []int{0}
}

func (x *ABCIResponses) GetFinalizeBlock() *types.ResponseFinalizeBlock {
	if x != nil {
		return x.FinalizeBlock
	}
	return nil
}

// ValidatorsInfo represents the latest validator set, or the last height it changed
type ValidatorsInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ValidatorSet      *types1.ValidatorSet `protobuf:"bytes,1,opt,name=validator_set,json=validatorSet,proto3" json:"validator_set,omitempty"`
	LastHeightChanged int64                `protobuf:"varint,2,opt,name=last_height_changed,json=lastHeightChanged,proto3" json:"last_height_changed,omitempty"`
}

func (x *ValidatorsInfo) Reset() {
	*x = ValidatorsInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_api_v1_state_types_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ValidatorsInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ValidatorsInfo) ProtoMessage() {}

func (x *ValidatorsInfo) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_api_v1_state_types_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ValidatorsInfo.ProtoReflect.Descriptor instead.
func (*ValidatorsInfo) Descriptor() ([]byte, []int) {
	return file_pkg_api_v1_state_types_proto_rawDescGZIP(), []int{1}
}

func (x *ValidatorsInfo) GetValidatorSet() *types1.ValidatorSet {
	if x != nil {
		return x.ValidatorSet
	}
	return nil
}

func (x *ValidatorsInfo) GetLastHeightChanged() int64 {
	if x != nil {
		return x.LastHeightChanged
	}
	return 0
}

// ConsensusParamsInfo represents the latest consensus params, or the last height it changed
type ConsensusParamsInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ConsensusParams   *types1.ConsensusParams `protobuf:"bytes,1,opt,name=consensus_params,json=consensusParams,proto3" json:"consensus_params,omitempty"`
	LastHeightChanged int64                   `protobuf:"varint,2,opt,name=last_height_changed,json=lastHeightChanged,proto3" json:"last_height_changed,omitempty"`
}

func (x *ConsensusParamsInfo) Reset() {
	*x = ConsensusParamsInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_api_v1_state_types_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ConsensusParamsInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ConsensusParamsInfo) ProtoMessage() {}

func (x *ConsensusParamsInfo) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_api_v1_state_types_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ConsensusParamsInfo.ProtoReflect.Descriptor instead.
func (*ConsensusParamsInfo) Descriptor() ([]byte, []int) {
	return file_pkg_api_v1_state_types_proto_rawDescGZIP(), []int{2}
}

func (x *ConsensusParamsInfo) GetConsensusParams() *types1.ConsensusParams {
	if x != nil {
		return x.ConsensusParams
	}
	return nil
}

func (x *ConsensusParamsInfo) GetLastHeightChanged() int64 {
	if x != nil {
		return x.LastHeightChanged
	}
	return 0
}

type Version struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Consensus *version.Consensus `protobuf:"bytes,1,opt,name=consensus,proto3" json:"consensus,omitempty"`
	Software  string             `protobuf:"bytes,2,opt,name=software,proto3" json:"software,omitempty"`
}

func (x *Version) Reset() {
	*x = Version{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_api_v1_state_types_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Version) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Version) ProtoMessage() {}

func (x *Version) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_api_v1_state_types_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Version.ProtoReflect.Descriptor instead.
func (*Version) Descriptor() ([]byte, []int) {
	return file_pkg_api_v1_state_types_proto_rawDescGZIP(), []int{3}
}

func (x *Version) GetConsensus() *version.Consensus {
	if x != nil {
		return x.Consensus
	}
	return nil
}

func (x *Version) GetSoftware() string {
	if x != nil {
		return x.Software
	}
	return ""
}

type State struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Version *Version `protobuf:"bytes,1,opt,name=version,proto3" json:"version,omitempty"`
	// immutable
	ChainId       string `protobuf:"bytes,2,opt,name=chain_id,json=chainId,proto3" json:"chain_id,omitempty"`
	InitialHeight int64  `protobuf:"varint,14,opt,name=initial_height,json=initialHeight,proto3" json:"initial_height,omitempty"`
	// LastBlockHeight=0 at genesis (ie. block(H=0) does not exist)
	LastBlockHeight int64                  `protobuf:"varint,3,opt,name=last_block_height,json=lastBlockHeight,proto3" json:"last_block_height,omitempty"`
	LastBlockId     *types1.BlockID        `protobuf:"bytes,4,opt,name=last_block_id,json=lastBlockId,proto3" json:"last_block_id,omitempty"`
	LastBlockTime   *timestamppb.Timestamp `protobuf:"bytes,5,opt,name=last_block_time,json=lastBlockTime,proto3" json:"last_block_time,omitempty"`
	// LastValidators is used to validate block.LastCommit.
	// Validators are persisted to the database separately every time they change,
	// so we can query for historical validator sets.
	// Note that if s.LastBlockHeight causes a valset change,
	// we set s.LastHeightValidatorsChanged = s.LastBlockHeight + 1 + 1
	// Extra +1 due to nextValSet delay.
	NextValidators              *types1.ValidatorSet `protobuf:"bytes,6,opt,name=next_validators,json=nextValidators,proto3" json:"next_validators,omitempty"`
	Validators                  *types1.ValidatorSet `protobuf:"bytes,7,opt,name=validators,proto3" json:"validators,omitempty"`
	LastValidators              *types1.ValidatorSet `protobuf:"bytes,8,opt,name=last_validators,json=lastValidators,proto3" json:"last_validators,omitempty"`
	LastHeightValidatorsChanged int64                `protobuf:"varint,9,opt,name=last_height_validators_changed,json=lastHeightValidatorsChanged,proto3" json:"last_height_validators_changed,omitempty"`
	// Consensus parameters used for validating blocks.
	// Changes returned by EndBlock and updated after Commit.
	ConsensusParams                  *types1.ConsensusParams `protobuf:"bytes,10,opt,name=consensus_params,json=consensusParams,proto3" json:"consensus_params,omitempty"`
	LastHeightConsensusParamsChanged int64                   `protobuf:"varint,11,opt,name=last_height_consensus_params_changed,json=lastHeightConsensusParamsChanged,proto3" json:"last_height_consensus_params_changed,omitempty"`
	// Merkle root of the results from executing prev block
	LastResultsHash []byte `protobuf:"bytes,12,opt,name=last_results_hash,json=lastResultsHash,proto3" json:"last_results_hash,omitempty"`
	// the latest AppHash we've received from calling abci.Commit()
	AppHash []byte `protobuf:"bytes,13,opt,name=app_hash,json=appHash,proto3" json:"app_hash,omitempty"`
}

func (x *State) Reset() {
	*x = State{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_api_v1_state_types_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *State) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*State) ProtoMessage() {}

func (x *State) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_api_v1_state_types_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use State.ProtoReflect.Descriptor instead.
func (*State) Descriptor() ([]byte, []int) {
	return file_pkg_api_v1_state_types_proto_rawDescGZIP(), []int{4}
}

func (x *State) GetVersion() *Version {
	if x != nil {
		return x.Version
	}
	return nil
}

func (x *State) GetChainId() string {
	if x != nil {
		return x.ChainId
	}
	return ""
}

func (x *State) GetInitialHeight() int64 {
	if x != nil {
		return x.InitialHeight
	}
	return 0
}

func (x *State) GetLastBlockHeight() int64 {
	if x != nil {
		return x.LastBlockHeight
	}
	return 0
}

func (x *State) GetLastBlockId() *types1.BlockID {
	if x != nil {
		return x.LastBlockId
	}
	return nil
}

func (x *State) GetLastBlockTime() *timestamppb.Timestamp {
	if x != nil {
		return x.LastBlockTime
	}
	return nil
}

func (x *State) GetNextValidators() *types1.ValidatorSet {
	if x != nil {
		return x.NextValidators
	}
	return nil
}

func (x *State) GetValidators() *types1.ValidatorSet {
	if x != nil {
		return x.Validators
	}
	return nil
}

func (x *State) GetLastValidators() *types1.ValidatorSet {
	if x != nil {
		return x.LastValidators
	}
	return nil
}

func (x *State) GetLastHeightValidatorsChanged() int64 {
	if x != nil {
		return x.LastHeightValidatorsChanged
	}
	return 0
}

func (x *State) GetConsensusParams() *types1.ConsensusParams {
	if x != nil {
		return x.ConsensusParams
	}
	return nil
}

func (x *State) GetLastHeightConsensusParamsChanged() int64 {
	if x != nil {
		return x.LastHeightConsensusParamsChanged
	}
	return 0
}

func (x *State) GetLastResultsHash() []byte {
	if x != nil {
		return x.LastResultsHash
	}
	return nil
}

func (x *State) GetAppHash() []byte {
	if x != nil {
		return x.AppHash
	}
	return nil
}

var File_pkg_api_v1_state_types_proto protoreflect.FileDescriptor

var file_pkg_api_v1_state_types_proto_rawDesc = []byte{
	0x0a, 0x1c, 0x70, 0x6b, 0x67, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x31, 0x2f, 0x73, 0x74, 0x61,
	0x74, 0x65, 0x2f, 0x74, 0x79, 0x70, 0x65, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x08,
	0x76, 0x31, 0x2e, 0x73, 0x74, 0x61, 0x74, 0x65, 0x1a, 0x14, 0x67, 0x6f, 0x67, 0x6f, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x2f, 0x67, 0x6f, 0x67, 0x6f, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1b,
	0x70, 0x6b, 0x67, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x31, 0x2f, 0x61, 0x62, 0x63, 0x69, 0x2f,
	0x74, 0x79, 0x70, 0x65, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1c, 0x70, 0x6b, 0x67,
	0x2f, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x31, 0x2f, 0x74, 0x79, 0x70, 0x65, 0x73, 0x2f, 0x74, 0x79,
	0x70, 0x65, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x20, 0x70, 0x6b, 0x67, 0x2f, 0x61,
	0x70, 0x69, 0x2f, 0x76, 0x31, 0x2f, 0x74, 0x79, 0x70, 0x65, 0x73, 0x2f, 0x76, 0x61, 0x6c, 0x69,
	0x64, 0x61, 0x74, 0x6f, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1d, 0x70, 0x6b, 0x67,
	0x2f, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x31, 0x2f, 0x74, 0x79, 0x70, 0x65, 0x73, 0x2f, 0x70, 0x61,
	0x72, 0x61, 0x6d, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1e, 0x70, 0x6b, 0x67, 0x2f,
	0x61, 0x70, 0x69, 0x2f, 0x76, 0x31, 0x2f, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x2f, 0x74,
	0x79, 0x70, 0x65, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65,
	0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x56, 0x0a, 0x0d, 0x41,
	0x42, 0x43, 0x49, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x73, 0x12, 0x45, 0x0a, 0x0e,
	0x66, 0x69, 0x6e, 0x61, 0x6c, 0x69, 0x7a, 0x65, 0x5f, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x76, 0x31, 0x2e, 0x61, 0x62, 0x63, 0x69, 0x2e, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x46, 0x69, 0x6e, 0x61, 0x6c, 0x69, 0x7a, 0x65, 0x42,
	0x6c, 0x6f, 0x63, 0x6b, 0x52, 0x0d, 0x66, 0x69, 0x6e, 0x61, 0x6c, 0x69, 0x7a, 0x65, 0x42, 0x6c,
	0x6f, 0x63, 0x6b, 0x22, 0x7d, 0x0a, 0x0e, 0x56, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x6f, 0x72,
	0x73, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x3b, 0x0a, 0x0d, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74,
	0x6f, 0x72, 0x5f, 0x73, 0x65, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x16, 0x2e, 0x76,
	0x31, 0x2e, 0x74, 0x79, 0x70, 0x65, 0x73, 0x2e, 0x56, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x6f,
	0x72, 0x53, 0x65, 0x74, 0x52, 0x0c, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x6f, 0x72, 0x53,
	0x65, 0x74, 0x12, 0x2e, 0x0a, 0x13, 0x6c, 0x61, 0x73, 0x74, 0x5f, 0x68, 0x65, 0x69, 0x67, 0x68,
	0x74, 0x5f, 0x63, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52,
	0x11, 0x6c, 0x61, 0x73, 0x74, 0x48, 0x65, 0x69, 0x67, 0x68, 0x74, 0x43, 0x68, 0x61, 0x6e, 0x67,
	0x65, 0x64, 0x22, 0x91, 0x01, 0x0a, 0x13, 0x43, 0x6f, 0x6e, 0x73, 0x65, 0x6e, 0x73, 0x75, 0x73,
	0x50, 0x61, 0x72, 0x61, 0x6d, 0x73, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x4a, 0x0a, 0x10, 0x63, 0x6f,
	0x6e, 0x73, 0x65, 0x6e, 0x73, 0x75, 0x73, 0x5f, 0x70, 0x61, 0x72, 0x61, 0x6d, 0x73, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x76, 0x31, 0x2e, 0x74, 0x79, 0x70, 0x65, 0x73, 0x2e,
	0x43, 0x6f, 0x6e, 0x73, 0x65, 0x6e, 0x73, 0x75, 0x73, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x73, 0x42,
	0x04, 0xc8, 0xde, 0x1f, 0x00, 0x52, 0x0f, 0x63, 0x6f, 0x6e, 0x73, 0x65, 0x6e, 0x73, 0x75, 0x73,
	0x50, 0x61, 0x72, 0x61, 0x6d, 0x73, 0x12, 0x2e, 0x0a, 0x13, 0x6c, 0x61, 0x73, 0x74, 0x5f, 0x68,
	0x65, 0x69, 0x67, 0x68, 0x74, 0x5f, 0x63, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x64, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x03, 0x52, 0x11, 0x6c, 0x61, 0x73, 0x74, 0x48, 0x65, 0x69, 0x67, 0x68, 0x74, 0x43,
	0x68, 0x61, 0x6e, 0x67, 0x65, 0x64, 0x22, 0x60, 0x0a, 0x07, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f,
	0x6e, 0x12, 0x39, 0x0a, 0x09, 0x63, 0x6f, 0x6e, 0x73, 0x65, 0x6e, 0x73, 0x75, 0x73, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x15, 0x2e, 0x76, 0x31, 0x2e, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f,
	0x6e, 0x2e, 0x43, 0x6f, 0x6e, 0x73, 0x65, 0x6e, 0x73, 0x75, 0x73, 0x42, 0x04, 0xc8, 0xde, 0x1f,
	0x00, 0x52, 0x09, 0x63, 0x6f, 0x6e, 0x73, 0x65, 0x6e, 0x73, 0x75, 0x73, 0x12, 0x1a, 0x0a, 0x08,
	0x73, 0x6f, 0x66, 0x74, 0x77, 0x61, 0x72, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08,
	0x73, 0x6f, 0x66, 0x74, 0x77, 0x61, 0x72, 0x65, 0x22, 0xb1, 0x06, 0x0a, 0x05, 0x53, 0x74, 0x61,
	0x74, 0x65, 0x12, 0x31, 0x0a, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x11, 0x2e, 0x76, 0x31, 0x2e, 0x73, 0x74, 0x61, 0x74, 0x65, 0x2e, 0x56,
	0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x42, 0x04, 0xc8, 0xde, 0x1f, 0x00, 0x52, 0x07, 0x76, 0x65,
	0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x26, 0x0a, 0x08, 0x63, 0x68, 0x61, 0x69, 0x6e, 0x5f, 0x69,
	0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x42, 0x0b, 0xe2, 0xde, 0x1f, 0x07, 0x43, 0x68, 0x61,
	0x69, 0x6e, 0x49, 0x44, 0x52, 0x07, 0x63, 0x68, 0x61, 0x69, 0x6e, 0x49, 0x64, 0x12, 0x25, 0x0a,
	0x0e, 0x69, 0x6e, 0x69, 0x74, 0x69, 0x61, 0x6c, 0x5f, 0x68, 0x65, 0x69, 0x67, 0x68, 0x74, 0x18,
	0x0e, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0d, 0x69, 0x6e, 0x69, 0x74, 0x69, 0x61, 0x6c, 0x48, 0x65,
	0x69, 0x67, 0x68, 0x74, 0x12, 0x2a, 0x0a, 0x11, 0x6c, 0x61, 0x73, 0x74, 0x5f, 0x62, 0x6c, 0x6f,
	0x63, 0x6b, 0x5f, 0x68, 0x65, 0x69, 0x67, 0x68, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52,
	0x0f, 0x6c, 0x61, 0x73, 0x74, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x48, 0x65, 0x69, 0x67, 0x68, 0x74,
	0x12, 0x4a, 0x0a, 0x0d, 0x6c, 0x61, 0x73, 0x74, 0x5f, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x5f, 0x69,
	0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x11, 0x2e, 0x76, 0x31, 0x2e, 0x74, 0x79, 0x70,
	0x65, 0x73, 0x2e, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x49, 0x44, 0x42, 0x13, 0xc8, 0xde, 0x1f, 0x00,
	0xe2, 0xde, 0x1f, 0x0b, 0x4c, 0x61, 0x73, 0x74, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x49, 0x44, 0x52,
	0x0b, 0x6c, 0x61, 0x73, 0x74, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x49, 0x64, 0x12, 0x4c, 0x0a, 0x0f,
	0x6c, 0x61, 0x73, 0x74, 0x5f, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18,
	0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d,
	0x70, 0x42, 0x08, 0xc8, 0xde, 0x1f, 0x00, 0x90, 0xdf, 0x1f, 0x01, 0x52, 0x0d, 0x6c, 0x61, 0x73,
	0x74, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x3f, 0x0a, 0x0f, 0x6e, 0x65,
	0x78, 0x74, 0x5f, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x6f, 0x72, 0x73, 0x18, 0x06, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x16, 0x2e, 0x76, 0x31, 0x2e, 0x74, 0x79, 0x70, 0x65, 0x73, 0x2e, 0x56,
	0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x6f, 0x72, 0x53, 0x65, 0x74, 0x52, 0x0e, 0x6e, 0x65, 0x78,
	0x74, 0x56, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x6f, 0x72, 0x73, 0x12, 0x36, 0x0a, 0x0a, 0x76,
	0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x6f, 0x72, 0x73, 0x18, 0x07, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x16, 0x2e, 0x76, 0x31, 0x2e, 0x74, 0x79, 0x70, 0x65, 0x73, 0x2e, 0x56, 0x61, 0x6c, 0x69, 0x64,
	0x61, 0x74, 0x6f, 0x72, 0x53, 0x65, 0x74, 0x52, 0x0a, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74,
	0x6f, 0x72, 0x73, 0x12, 0x3f, 0x0a, 0x0f, 0x6c, 0x61, 0x73, 0x74, 0x5f, 0x76, 0x61, 0x6c, 0x69,
	0x64, 0x61, 0x74, 0x6f, 0x72, 0x73, 0x18, 0x08, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x16, 0x2e, 0x76,
	0x31, 0x2e, 0x74, 0x79, 0x70, 0x65, 0x73, 0x2e, 0x56, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x6f,
	0x72, 0x53, 0x65, 0x74, 0x52, 0x0e, 0x6c, 0x61, 0x73, 0x74, 0x56, 0x61, 0x6c, 0x69, 0x64, 0x61,
	0x74, 0x6f, 0x72, 0x73, 0x12, 0x43, 0x0a, 0x1e, 0x6c, 0x61, 0x73, 0x74, 0x5f, 0x68, 0x65, 0x69,
	0x67, 0x68, 0x74, 0x5f, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x6f, 0x72, 0x73, 0x5f, 0x63,
	0x68, 0x61, 0x6e, 0x67, 0x65, 0x64, 0x18, 0x09, 0x20, 0x01, 0x28, 0x03, 0x52, 0x1b, 0x6c, 0x61,
	0x73, 0x74, 0x48, 0x65, 0x69, 0x67, 0x68, 0x74, 0x56, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x6f,
	0x72, 0x73, 0x43, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x64, 0x12, 0x4a, 0x0a, 0x10, 0x63, 0x6f, 0x6e,
	0x73, 0x65, 0x6e, 0x73, 0x75, 0x73, 0x5f, 0x70, 0x61, 0x72, 0x61, 0x6d, 0x73, 0x18, 0x0a, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x76, 0x31, 0x2e, 0x74, 0x79, 0x70, 0x65, 0x73, 0x2e, 0x43,
	0x6f, 0x6e, 0x73, 0x65, 0x6e, 0x73, 0x75, 0x73, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x73, 0x42, 0x04,
	0xc8, 0xde, 0x1f, 0x00, 0x52, 0x0f, 0x63, 0x6f, 0x6e, 0x73, 0x65, 0x6e, 0x73, 0x75, 0x73, 0x50,
	0x61, 0x72, 0x61, 0x6d, 0x73, 0x12, 0x4e, 0x0a, 0x24, 0x6c, 0x61, 0x73, 0x74, 0x5f, 0x68, 0x65,
	0x69, 0x67, 0x68, 0x74, 0x5f, 0x63, 0x6f, 0x6e, 0x73, 0x65, 0x6e, 0x73, 0x75, 0x73, 0x5f, 0x70,
	0x61, 0x72, 0x61, 0x6d, 0x73, 0x5f, 0x63, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x64, 0x18, 0x0b, 0x20,
	0x01, 0x28, 0x03, 0x52, 0x20, 0x6c, 0x61, 0x73, 0x74, 0x48, 0x65, 0x69, 0x67, 0x68, 0x74, 0x43,
	0x6f, 0x6e, 0x73, 0x65, 0x6e, 0x73, 0x75, 0x73, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x73, 0x43, 0x68,
	0x61, 0x6e, 0x67, 0x65, 0x64, 0x12, 0x2a, 0x0a, 0x11, 0x6c, 0x61, 0x73, 0x74, 0x5f, 0x72, 0x65,
	0x73, 0x75, 0x6c, 0x74, 0x73, 0x5f, 0x68, 0x61, 0x73, 0x68, 0x18, 0x0c, 0x20, 0x01, 0x28, 0x0c,
	0x52, 0x0f, 0x6c, 0x61, 0x73, 0x74, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x73, 0x48, 0x61, 0x73,
	0x68, 0x12, 0x19, 0x0a, 0x08, 0x61, 0x70, 0x70, 0x5f, 0x68, 0x61, 0x73, 0x68, 0x18, 0x0d, 0x20,
	0x01, 0x28, 0x0c, 0x52, 0x07, 0x61, 0x70, 0x70, 0x48, 0x61, 0x73, 0x68, 0x42, 0x31, 0x5a, 0x2f,
	0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x62, 0x68, 0x6f, 0x6a, 0x70,
	0x75, 0x72, 0x2f, 0x73, 0x74, 0x61, 0x74, 0x65, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x61, 0x70, 0x69,
	0x2f, 0x76, 0x31, 0x2f, 0x73, 0x74, 0x61, 0x74, 0x65, 0x3b, 0x73, 0x74, 0x61, 0x74, 0x65, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_pkg_api_v1_state_types_proto_rawDescOnce sync.Once
	file_pkg_api_v1_state_types_proto_rawDescData = file_pkg_api_v1_state_types_proto_rawDesc
)

func file_pkg_api_v1_state_types_proto_rawDescGZIP() []byte {
	file_pkg_api_v1_state_types_proto_rawDescOnce.Do(func() {
		file_pkg_api_v1_state_types_proto_rawDescData = protoimpl.X.CompressGZIP(file_pkg_api_v1_state_types_proto_rawDescData)
	})
	return file_pkg_api_v1_state_types_proto_rawDescData
}

var file_pkg_api_v1_state_types_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_pkg_api_v1_state_types_proto_goTypes = []interface{}{
	(*ABCIResponses)(nil),               // 0: v1.state.ABCIResponses
	(*ValidatorsInfo)(nil),              // 1: v1.state.ValidatorsInfo
	(*ConsensusParamsInfo)(nil),         // 2: v1.state.ConsensusParamsInfo
	(*Version)(nil),                     // 3: v1.state.Version
	(*State)(nil),                       // 4: v1.state.State
	(*types.ResponseFinalizeBlock)(nil), // 5: v1.abci.ResponseFinalizeBlock
	(*types1.ValidatorSet)(nil),         // 6: v1.types.ValidatorSet
	(*types1.ConsensusParams)(nil),      // 7: v1.types.ConsensusParams
	(*version.Consensus)(nil),           // 8: v1.version.Consensus
	(*types1.BlockID)(nil),              // 9: v1.types.BlockID
	(*timestamppb.Timestamp)(nil),       // 10: google.protobuf.Timestamp
}
var file_pkg_api_v1_state_types_proto_depIdxs = []int32{
	5,  // 0: v1.state.ABCIResponses.finalize_block:type_name -> v1.abci.ResponseFinalizeBlock
	6,  // 1: v1.state.ValidatorsInfo.validator_set:type_name -> v1.types.ValidatorSet
	7,  // 2: v1.state.ConsensusParamsInfo.consensus_params:type_name -> v1.types.ConsensusParams
	8,  // 3: v1.state.Version.consensus:type_name -> v1.version.Consensus
	3,  // 4: v1.state.State.version:type_name -> v1.state.Version
	9,  // 5: v1.state.State.last_block_id:type_name -> v1.types.BlockID
	10, // 6: v1.state.State.last_block_time:type_name -> google.protobuf.Timestamp
	6,  // 7: v1.state.State.next_validators:type_name -> v1.types.ValidatorSet
	6,  // 8: v1.state.State.validators:type_name -> v1.types.ValidatorSet
	6,  // 9: v1.state.State.last_validators:type_name -> v1.types.ValidatorSet
	7,  // 10: v1.state.State.consensus_params:type_name -> v1.types.ConsensusParams
	11, // [11:11] is the sub-list for method output_type
	11, // [11:11] is the sub-list for method input_type
	11, // [11:11] is the sub-list for extension type_name
	11, // [11:11] is the sub-list for extension extendee
	0,  // [0:11] is the sub-list for field type_name
}

func init() { file_pkg_api_v1_state_types_proto_init() }
func file_pkg_api_v1_state_types_proto_init() {
	if File_pkg_api_v1_state_types_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_pkg_api_v1_state_types_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ABCIResponses); i {
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
		file_pkg_api_v1_state_types_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ValidatorsInfo); i {
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
		file_pkg_api_v1_state_types_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ConsensusParamsInfo); i {
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
		file_pkg_api_v1_state_types_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Version); i {
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
		file_pkg_api_v1_state_types_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*State); i {
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
			RawDescriptor: file_pkg_api_v1_state_types_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_pkg_api_v1_state_types_proto_goTypes,
		DependencyIndexes: file_pkg_api_v1_state_types_proto_depIdxs,
		MessageInfos:      file_pkg_api_v1_state_types_proto_msgTypes,
	}.Build()
	File_pkg_api_v1_state_types_proto = out.File
	file_pkg_api_v1_state_types_proto_rawDesc = nil
	file_pkg_api_v1_state_types_proto_goTypes = nil
	file_pkg_api_v1_state_types_proto_depIdxs = nil
}