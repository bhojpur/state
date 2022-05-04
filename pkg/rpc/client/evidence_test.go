package client_test

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
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	typespb "github.com/bhojpur/state/pkg/api/v1/types"
	"github.com/bhojpur/state/pkg/crypto"
	librand "github.com/bhojpur/state/pkg/libs/rand"
	"github.com/bhojpur/state/pkg/privval"
	"github.com/bhojpur/state/pkg/rpc/client"
	"github.com/bhojpur/state/pkg/types"
)

func newEvidence(t *testing.T, val *privval.FilePV,
	vote *types.Vote, vote2 *types.Vote,
	chainID string,
	timestamp time.Time,
) *types.DuplicateVoteEvidence {
	t.Helper()
	var err error

	v := vote.ToProto()
	v2 := vote2.ToProto()

	vote.Signature, err = val.Key.PrivKey.Sign(types.VoteSignBytes(chainID, v))
	require.NoError(t, err)

	vote2.Signature, err = val.Key.PrivKey.Sign(types.VoteSignBytes(chainID, v2))
	require.NoError(t, err)

	validator := types.NewValidator(val.Key.PubKey, 10)
	valSet := types.NewValidatorSet([]*types.Validator{validator})

	ev, err := types.NewDuplicateVoteEvidence(vote, vote2, timestamp, valSet)
	require.NoError(t, err)
	return ev
}

func makeEvidences(
	t *testing.T,
	val *privval.FilePV,
	chainID string,
	timestamp time.Time,
) (correct *types.DuplicateVoteEvidence, fakes []*types.DuplicateVoteEvidence) {
	vote := types.Vote{
		ValidatorAddress: val.Key.Address,
		ValidatorIndex:   0,
		Height:           1,
		Round:            0,
		Type:             typespb.PrevoteType,
		Timestamp:        timestamp,
		BlockID: types.BlockID{
			Hash: crypto.Checksum(librand.Bytes(crypto.HashSize)),
			PartSetHeader: types.PartSetHeader{
				Total: 1000,
				Hash:  crypto.Checksum([]byte("partset")),
			},
		},
	}

	vote2 := vote
	vote2.BlockID.Hash = crypto.Checksum([]byte("blockhash2"))
	correct = newEvidence(t, val, &vote, &vote2, chainID, timestamp)

	fakes = make([]*types.DuplicateVoteEvidence, 0)

	// different address
	{
		v := vote2
		v.ValidatorAddress = []byte("some_address")
		fakes = append(fakes, newEvidence(t, val, &vote, &v, chainID, timestamp))
	}

	// different height
	{
		v := vote2
		v.Height = vote.Height + 1
		fakes = append(fakes, newEvidence(t, val, &vote, &v, chainID, timestamp))
	}

	// different round
	{
		v := vote2
		v.Round = vote.Round + 1
		fakes = append(fakes, newEvidence(t, val, &vote, &v, chainID, timestamp))
	}

	// different type
	{
		v := vote2
		v.Type = typespb.PrecommitType
		fakes = append(fakes, newEvidence(t, val, &vote, &v, chainID, timestamp))
	}

	// exactly same vote
	{
		v := vote
		fakes = append(fakes, newEvidence(t, val, &vote, &v, chainID, timestamp))
	}

	return correct, fakes
}

func waitForBlock(ctx context.Context, t *testing.T, c client.Client, height int64) {
	timer := time.NewTimer(0 * time.Millisecond)
	defer timer.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-timer.C:
			status, err := c.Status(ctx)
			require.NoError(t, err)
			if status.SyncInfo.LatestBlockHeight >= height {
				return
			}
			timer.Reset(200 * time.Millisecond)
		}
	}
}
