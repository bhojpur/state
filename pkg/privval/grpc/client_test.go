package grpc_test

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
	"net"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	privvalproto "github.com/bhojpur/state/pkg/api/v1/privval"
	v1 "github.com/bhojpur/state/pkg/api/v1/types"
	"github.com/bhojpur/state/pkg/crypto"
	"github.com/bhojpur/state/pkg/libs/log"
	librand "github.com/bhojpur/state/pkg/libs/rand"
	privrpc "github.com/bhojpur/state/pkg/privval/grpc"
	"github.com/bhojpur/state/pkg/types"
)

const chainID = "chain-id"

func dialer(t *testing.T, pv types.PrivValidator, logger log.Logger) (*grpc.Server, func(context.Context, string) (net.Conn, error)) {
	listener := bufconn.Listen(1024 * 1024)

	server := grpc.NewServer()

	s := privrpc.NewSignerServer(logger, chainID, pv)

	privvalproto.RegisterPrivValidatorAPIServer(server, s)

	go func() { require.NoError(t, server.Serve(listener)) }()

	return server, func(context.Context, string) (net.Conn, error) {
		return listener.Dial()
	}
}

func TestSignerClient_GetPubKey(t *testing.T) {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mockPV := types.NewMockPV()
	logger := log.NewTestingLogger(t)
	srv, dialer := dialer(t, mockPV, logger)
	defer srv.Stop()

	conn, err := grpc.DialContext(ctx, "",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithContextDialer(dialer),
	)
	require.NoError(t, err)
	defer conn.Close()

	client, err := privrpc.NewSignerClient(conn, chainID, logger)
	require.NoError(t, err)

	pk, err := client.GetPubKey(ctx)
	require.NoError(t, err)
	assert.Equal(t, mockPV.PrivKey.PubKey(), pk)
}

func TestSignerClient_SignVote(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mockPV := types.NewMockPV()
	logger := log.NewTestingLogger(t)
	srv, dialer := dialer(t, mockPV, logger)
	defer srv.Stop()

	conn, err := grpc.DialContext(ctx, "",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithContextDialer(dialer),
	)
	require.NoError(t, err)
	defer conn.Close()

	client, err := privrpc.NewSignerClient(conn, chainID, logger)
	require.NoError(t, err)

	ts := time.Now()
	hash := librand.Bytes(crypto.HashSize)
	valAddr := librand.Bytes(crypto.AddressSize)

	want := &types.Vote{
		Type:             v1.PrecommitType,
		Height:           1,
		Round:            2,
		BlockID:          types.BlockID{Hash: hash, PartSetHeader: types.PartSetHeader{Hash: hash, Total: 2}},
		Timestamp:        ts,
		ValidatorAddress: valAddr,
		ValidatorIndex:   1,
	}

	have := &types.Vote{
		Type:             v1.PrecommitType,
		Height:           1,
		Round:            2,
		BlockID:          types.BlockID{Hash: hash, PartSetHeader: types.PartSetHeader{Hash: hash, Total: 2}},
		Timestamp:        ts,
		ValidatorAddress: valAddr,
		ValidatorIndex:   1,
	}

	pbHave := have.ToProto()

	err = client.SignVote(ctx, chainID, pbHave)
	require.NoError(t, err)

	pbWant := want.ToProto()

	require.NoError(t, mockPV.SignVote(ctx, chainID, pbWant))

	assert.Equal(t, pbWant.Signature, pbHave.Signature)
}

func TestSignerClient_SignProposal(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mockPV := types.NewMockPV()
	logger := log.NewTestingLogger(t)
	srv, dialer := dialer(t, mockPV, logger)
	defer srv.Stop()

	conn, err := grpc.DialContext(ctx, "",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithContextDialer(dialer),
	)
	require.NoError(t, err)
	defer conn.Close()

	client, err := privrpc.NewSignerClient(conn, chainID, logger)
	require.NoError(t, err)

	ts := time.Now()
	hash := librand.Bytes(crypto.HashSize)

	have := &types.Proposal{
		Type:      v1.ProposalType,
		Height:    1,
		Round:     2,
		POLRound:  2,
		BlockID:   types.BlockID{Hash: hash, PartSetHeader: types.PartSetHeader{Hash: hash, Total: 2}},
		Timestamp: ts,
	}
	want := &types.Proposal{
		Type:      v1.ProposalType,
		Height:    1,
		Round:     2,
		POLRound:  2,
		BlockID:   types.BlockID{Hash: hash, PartSetHeader: types.PartSetHeader{Hash: hash, Total: 2}},
		Timestamp: ts,
	}

	pbHave := have.ToProto()

	err = client.SignProposal(ctx, chainID, pbHave)
	require.NoError(t, err)

	pbWant := want.ToProto()

	require.NoError(t, mockPV.SignProposal(ctx, chainID, pbWant))

	assert.Equal(t, pbWant.Signature, pbHave.Signature)
}
