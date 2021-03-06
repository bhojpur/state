// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package privval

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// PrivValidatorAPIClient is the client API for PrivValidatorAPI service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type PrivValidatorAPIClient interface {
	GetPubKey(ctx context.Context, in *PubKeyRequest, opts ...grpc.CallOption) (*PubKeyResponse, error)
	SignVote(ctx context.Context, in *SignVoteRequest, opts ...grpc.CallOption) (*SignedVoteResponse, error)
	SignProposal(ctx context.Context, in *SignProposalRequest, opts ...grpc.CallOption) (*SignedProposalResponse, error)
}

type privValidatorAPIClient struct {
	cc grpc.ClientConnInterface
}

func NewPrivValidatorAPIClient(cc grpc.ClientConnInterface) PrivValidatorAPIClient {
	return &privValidatorAPIClient{cc}
}

func (c *privValidatorAPIClient) GetPubKey(ctx context.Context, in *PubKeyRequest, opts ...grpc.CallOption) (*PubKeyResponse, error) {
	out := new(PubKeyResponse)
	err := c.cc.Invoke(ctx, "/v1.privval.PrivValidatorAPI/GetPubKey", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *privValidatorAPIClient) SignVote(ctx context.Context, in *SignVoteRequest, opts ...grpc.CallOption) (*SignedVoteResponse, error) {
	out := new(SignedVoteResponse)
	err := c.cc.Invoke(ctx, "/v1.privval.PrivValidatorAPI/SignVote", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *privValidatorAPIClient) SignProposal(ctx context.Context, in *SignProposalRequest, opts ...grpc.CallOption) (*SignedProposalResponse, error) {
	out := new(SignedProposalResponse)
	err := c.cc.Invoke(ctx, "/v1.privval.PrivValidatorAPI/SignProposal", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PrivValidatorAPIServer is the server API for PrivValidatorAPI service.
// All implementations must embed UnimplementedPrivValidatorAPIServer
// for forward compatibility
type PrivValidatorAPIServer interface {
	GetPubKey(context.Context, *PubKeyRequest) (*PubKeyResponse, error)
	SignVote(context.Context, *SignVoteRequest) (*SignedVoteResponse, error)
	SignProposal(context.Context, *SignProposalRequest) (*SignedProposalResponse, error)
	mustEmbedUnimplementedPrivValidatorAPIServer()
}

// UnimplementedPrivValidatorAPIServer must be embedded to have forward compatible implementations.
type UnimplementedPrivValidatorAPIServer struct {
}

func (UnimplementedPrivValidatorAPIServer) GetPubKey(context.Context, *PubKeyRequest) (*PubKeyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPubKey not implemented")
}
func (UnimplementedPrivValidatorAPIServer) SignVote(context.Context, *SignVoteRequest) (*SignedVoteResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SignVote not implemented")
}
func (UnimplementedPrivValidatorAPIServer) SignProposal(context.Context, *SignProposalRequest) (*SignedProposalResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SignProposal not implemented")
}
func (UnimplementedPrivValidatorAPIServer) mustEmbedUnimplementedPrivValidatorAPIServer() {}

// UnsafePrivValidatorAPIServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to PrivValidatorAPIServer will
// result in compilation errors.
type UnsafePrivValidatorAPIServer interface {
	mustEmbedUnimplementedPrivValidatorAPIServer()
}

func RegisterPrivValidatorAPIServer(s grpc.ServiceRegistrar, srv PrivValidatorAPIServer) {
	s.RegisterService(&PrivValidatorAPI_ServiceDesc, srv)
}

func _PrivValidatorAPI_GetPubKey_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PubKeyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PrivValidatorAPIServer).GetPubKey(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/v1.privval.PrivValidatorAPI/GetPubKey",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PrivValidatorAPIServer).GetPubKey(ctx, req.(*PubKeyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PrivValidatorAPI_SignVote_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SignVoteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PrivValidatorAPIServer).SignVote(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/v1.privval.PrivValidatorAPI/SignVote",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PrivValidatorAPIServer).SignVote(ctx, req.(*SignVoteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PrivValidatorAPI_SignProposal_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SignProposalRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PrivValidatorAPIServer).SignProposal(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/v1.privval.PrivValidatorAPI/SignProposal",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PrivValidatorAPIServer).SignProposal(ctx, req.(*SignProposalRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// PrivValidatorAPI_ServiceDesc is the grpc.ServiceDesc for PrivValidatorAPI service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var PrivValidatorAPI_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "v1.privval.PrivValidatorAPI",
	HandlerType: (*PrivValidatorAPIServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetPubKey",
			Handler:    _PrivValidatorAPI_GetPubKey_Handler,
		},
		{
			MethodName: "SignVote",
			Handler:    _PrivValidatorAPI_SignVote_Handler,
		},
		{
			MethodName: "SignProposal",
			Handler:    _PrivValidatorAPI_SignProposal_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "pkg/api/v1/privval/service.proto",
}
