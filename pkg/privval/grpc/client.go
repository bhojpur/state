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
	"context"

	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/status"

	privvalproto "github.com/bhojpur/state/pkg/api/v1/privval"
	v1 "github.com/bhojpur/state/pkg/api/v1/types"
	"github.com/bhojpur/state/pkg/crypto"
	"github.com/bhojpur/state/pkg/crypto/encoding"
	"github.com/bhojpur/state/pkg/libs/log"
	"github.com/bhojpur/state/pkg/types"
)

// SignerClient implements PrivValidator.
// Handles remote validator connections that provide signing services
type SignerClient struct {
	logger log.Logger

	client  privvalproto.PrivValidatorAPIClient
	conn    *grpc.ClientConn
	chainID string
}

var _ types.PrivValidator = (*SignerClient)(nil)

// NewSignerClient returns an instance of SignerClient.
// it will start the endpoint (if not already started)
func NewSignerClient(conn *grpc.ClientConn,
	chainID string, log log.Logger) (*SignerClient, error) {

	sc := &SignerClient{
		logger:  log,
		chainID: chainID,
		client:  privvalproto.NewPrivValidatorAPIClient(conn), // Create the Private Validator Client
	}

	return sc, nil
}

// Close closes the underlying connection
func (sc *SignerClient) Close() error {
	sc.logger.Info("Stopping service")
	if sc.conn != nil {
		return sc.conn.Close()
	}
	return nil
}

// Implement PrivValidator

// GetPubKey retrieves a public key from a remote signer
// returns an error if client is not able to provide the key
func (sc *SignerClient) GetPubKey(ctx context.Context) (crypto.PubKey, error) {
	resp, err := sc.client.GetPubKey(ctx, &privvalproto.PubKeyRequest{ChainId: sc.chainID})
	if err != nil {
		errStatus, _ := status.FromError(err)
		sc.logger.Error("SignerClient::GetPubKey", "err", errStatus.Message())
		return nil, errStatus.Err()
	}

	pk, err := encoding.PubKeyFromProto(resp.PubKey)
	if err != nil {
		return nil, err
	}

	return pk, nil
}

// SignVote requests a remote signer to sign a vote
func (sc *SignerClient) SignVote(ctx context.Context, chainID string, vote *v1.Vote) error {
	resp, err := sc.client.SignVote(ctx, &privvalproto.SignVoteRequest{ChainId: sc.chainID, Vote: vote})
	if err != nil {
		errStatus, _ := status.FromError(err)
		sc.logger.Error("Client SignVote", "err", errStatus.Message())
		return errStatus.Err()
	}

	*vote = resp.Vote

	return nil
}

// SignProposal requests a remote signer to sign a proposal
func (sc *SignerClient) SignProposal(ctx context.Context, chainID string, proposal *v1.Proposal) error {
	resp, err := sc.client.SignProposal(
		ctx, &privvalproto.SignProposalRequest{ChainId: chainID, Proposal: proposal})

	if err != nil {
		errStatus, _ := status.FromError(err)
		sc.logger.Error("SignerClient::SignProposal", "err", errStatus.Message())
		return errStatus.Err()
	}

	*proposal = resp.Proposal

	return nil
}
