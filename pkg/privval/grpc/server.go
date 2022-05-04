package grpc

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

import (
	context "context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	privvalproto "github.com/bhojpur/state/pkg/api/v1/privval"
	"github.com/bhojpur/state/pkg/crypto"
	"github.com/bhojpur/state/pkg/crypto/encoding"
	"github.com/bhojpur/state/pkg/libs/log"
	"github.com/bhojpur/state/pkg/types"
)

// SignerServer implements PrivValidatorAPIServer 9generated via protobuf services)
// Handles remote validator connections that provide signing services
type SignerServer struct {
	privvalproto.UnimplementedPrivValidatorAPIServer
	logger  log.Logger
	chainID string
	privVal types.PrivValidator
}

func NewSignerServer(logger log.Logger, chainID string, privVal types.PrivValidator) *SignerServer {
	return &SignerServer{
		logger:  logger,
		chainID: chainID,
		privVal: privVal,
	}
}

var _ privvalproto.PrivValidatorAPIServer = (*SignerServer)(nil)

// PubKey receives a request for the pubkey
// returns the pubkey on success and error on failure
func (ss *SignerServer) GetPubKey(ctx context.Context, req *privvalproto.PubKeyRequest) (
	*privvalproto.PubKeyResponse, error) {
	var pubKey crypto.PubKey

	pubKey, err := ss.privVal.GetPubKey(ctx)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "error getting pubkey: %v", err)
	}

	pk, err := encoding.PubKeyToProto(pubKey)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error transitioning pubkey to proto: %v", err)
	}

	ss.logger.Info("SignerServer: GetPubKey Success")

	return &privvalproto.PubKeyResponse{PubKey: &pk}, nil
}

// SignVote receives a vote sign requests, attempts to sign it
// returns SignedVoteResponse on success and error on failure
func (ss *SignerServer) SignVote(ctx context.Context, req *privvalproto.SignVoteRequest) (*privvalproto.SignedVoteResponse, error) {
	vote := req.Vote

	err := ss.privVal.SignVote(ctx, req.ChainId, vote)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "error signing vote: %v", err)
	}

	ss.logger.Info("SignerServer: SignVote Success", "height", req.Vote.Height)

	return &privvalproto.SignedVoteResponse{Vote: vote}, nil
}

// SignProposal receives a proposal sign requests, attempts to sign it
// returns SignedProposalResponse on success and error on failure
func (ss *SignerServer) SignProposal(ctx context.Context, req *privvalproto.SignProposalRequest) (*privvalproto.SignedProposalResponse, error) {
	proposal := req.Proposal

	err := ss.privVal.SignProposal(ctx, req.ChainId, proposal)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "error signing proposal: %v", err)
	}

	ss.logger.Info("SignerServer: SignProposal Success", "height", req.Proposal.Height)

	return &privvalproto.SignedProposalResponse{Proposal: proposal}, nil
}
