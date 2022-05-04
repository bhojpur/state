package types

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

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	v1 "github.com/bhojpur/state/pkg/api/v1/types"
	libmath "github.com/bhojpur/state/pkg/libs/math"
)

// Check VerifyCommit, VerifyCommitLight and VerifyCommitLightTrusting basic
// verification.
func TestValidatorSet_VerifyCommit_All(t *testing.T) {
	var (
		round  = int32(0)
		height = int64(100)

		blockID    = makeBlockID([]byte("blockhash"), 1000, []byte("partshash"))
		chainID    = "Lalande21185"
		trustLevel = libmath.Fraction{Numerator: 2, Denominator: 3}
	)

	testCases := []struct {
		description string
		// vote chainID
		chainID string
		// vote blockID
		blockID BlockID
		valSize int

		// height of the commit
		height int64

		// votes
		blockVotes  int
		nilVotes    int
		absentVotes int

		expErr bool
	}{
		{"good (batch verification)", chainID, blockID, 3, height, 3, 0, 0, false},
		{"good (single verification)", chainID, blockID, 1, height, 1, 0, 0, false},

		{"wrong signature (#0)", "EpsilonEridani", blockID, 2, height, 2, 0, 0, true},
		{"wrong block ID", chainID, makeBlockIDRandom(), 2, height, 2, 0, 0, true},
		{"wrong height", chainID, blockID, 1, height - 1, 1, 0, 0, true},

		{"wrong set size: 4 vs 3", chainID, blockID, 4, height, 3, 0, 0, true},
		{"wrong set size: 1 vs 2", chainID, blockID, 1, height, 2, 0, 0, true},

		{"insufficient voting power: got 30, needed more than 66", chainID, blockID, 10, height, 3, 2, 5, true},
		{"insufficient voting power: got 0, needed more than 6", chainID, blockID, 1, height, 0, 0, 1, true},
		{"insufficient voting power: got 60, needed more than 60", chainID, blockID, 9, height, 6, 3, 0, true},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.description, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			_, valSet, vals := randVoteSet(ctx, t, tc.height, round, v1.PrecommitType, tc.valSize, 10)

			totalVotes := tc.blockVotes + tc.absentVotes + tc.nilVotes
			sigs := make([]CommitSig, totalVotes)
			vi := 0
			// add absent sigs first
			for i := 0; i < tc.absentVotes; i++ {
				sigs[vi] = NewCommitSigAbsent()
				vi++
			}
			for i := 0; i < tc.blockVotes+tc.nilVotes; i++ {

				pubKey, err := vals[vi%len(vals)].GetPubKey(ctx)
				require.NoError(t, err)
				vote := &Vote{
					ValidatorAddress: pubKey.Address(),
					ValidatorIndex:   int32(vi),
					Height:           tc.height,
					Round:            round,
					Type:             v1.PrecommitType,
					BlockID:          tc.blockID,
					Timestamp:        time.Now(),
				}
				if i >= tc.blockVotes {
					vote.BlockID = BlockID{}
				}

				v := vote.ToProto()

				require.NoError(t, vals[vi%len(vals)].SignVote(ctx, tc.chainID, v))
				vote.Signature = v.Signature

				sigs[vi] = vote.CommitSig()

				vi++
			}
			commit := NewCommit(tc.height, round, tc.blockID, sigs)

			err := valSet.VerifyCommit(chainID, blockID, height, commit)
			if tc.expErr {
				if assert.Error(t, err, "VerifyCommit") {
					assert.Contains(t, err.Error(), tc.description, "VerifyCommit")
				}
			} else {
				assert.NoError(t, err, "VerifyCommit")
			}

			err = valSet.VerifyCommitLight(chainID, blockID, height, commit)
			if tc.expErr {
				if assert.Error(t, err, "VerifyCommitLight") {
					assert.Contains(t, err.Error(), tc.description, "VerifyCommitLight")
				}
			} else {
				assert.NoError(t, err, "VerifyCommitLight")
			}

			// only a subsection of the tests apply to VerifyCommitLightTrusting
			if totalVotes != tc.valSize || !tc.blockID.Equals(blockID) || tc.height != height {
				tc.expErr = false
			}
			err = valSet.VerifyCommitLightTrusting(chainID, commit, trustLevel)
			if tc.expErr {
				if assert.Error(t, err, "VerifyCommitLightTrusting") {
					assert.Contains(t, err.Error(), tc.description, "VerifyCommitLightTrusting")
				}
			} else {
				assert.NoError(t, err, "VerifyCommitLightTrusting")
			}
		})
	}
}

func TestValidatorSet_VerifyCommit_CheckAllSignatures(t *testing.T) {
	var (
		chainID = "test_chain_id"
		h       = int64(3)
		blockID = makeBlockIDRandom()
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	voteSet, valSet, vals := randVoteSet(ctx, t, h, 0, v1.PrecommitType, 4, 10)
	commit, err := makeCommit(ctx, blockID, h, 0, voteSet, vals, time.Now())

	require.NoError(t, err)
	require.NoError(t, valSet.VerifyCommit(chainID, blockID, h, commit))

	// malleate 4th signature
	vote := voteSet.GetByIndex(3)
	v := vote.ToProto()
	err = vals[3].SignVote(ctx, "CentaurusA", v)
	require.NoError(t, err)
	vote.Signature = v.Signature
	commit.Signatures[3] = vote.CommitSig()

	err = valSet.VerifyCommit(chainID, blockID, h, commit)
	if assert.Error(t, err) {
		assert.Contains(t, err.Error(), "wrong signature (#3)")
	}
}

func TestValidatorSet_VerifyCommitLight_ReturnsAsSoonAsMajorityOfVotingPowerSigned(t *testing.T) {
	var (
		chainID = "test_chain_id"
		h       = int64(3)
		blockID = makeBlockIDRandom()
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	voteSet, valSet, vals := randVoteSet(ctx, t, h, 0, v1.PrecommitType, 4, 10)
	commit, err := makeCommit(ctx, blockID, h, 0, voteSet, vals, time.Now())

	require.NoError(t, err)
	require.NoError(t, valSet.VerifyCommit(chainID, blockID, h, commit))

	// malleate 4th signature (3 signatures are enough for 2/3+)
	vote := voteSet.GetByIndex(3)
	v := vote.ToProto()
	err = vals[3].SignVote(ctx, "CentaurusA", v)
	require.NoError(t, err)
	vote.Signature = v.Signature
	commit.Signatures[3] = vote.CommitSig()

	err = valSet.VerifyCommitLight(chainID, blockID, h, commit)
	assert.NoError(t, err)
}

func TestValidatorSet_VerifyCommitLightTrusting_ReturnsAsSoonAsTrustLevelOfVotingPowerSigned(t *testing.T) {
	var (
		chainID = "test_chain_id"
		h       = int64(3)
		blockID = makeBlockIDRandom()
	)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	voteSet, valSet, vals := randVoteSet(ctx, t, h, 0, v1.PrecommitType, 4, 10)
	commit, err := makeCommit(ctx, blockID, h, 0, voteSet, vals, time.Now())

	require.NoError(t, err)
	require.NoError(t, valSet.VerifyCommit(chainID, blockID, h, commit))

	// malleate 3rd signature (2 signatures are enough for 1/3+ trust level)
	vote := voteSet.GetByIndex(2)
	v := vote.ToProto()
	err = vals[2].SignVote(ctx, "CentaurusA", v)
	require.NoError(t, err)
	vote.Signature = v.Signature
	commit.Signatures[2] = vote.CommitSig()

	err = valSet.VerifyCommitLightTrusting(chainID, commit, libmath.Fraction{Numerator: 1, Denominator: 3})
	assert.NoError(t, err)
}

func TestValidatorSet_VerifyCommitLightTrusting(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var (
		blockID                       = makeBlockIDRandom()
		voteSet, originalValset, vals = randVoteSet(ctx, t, 1, 1, v1.PrecommitType, 6, 1)
		commit, err                   = makeCommit(ctx, blockID, 1, 1, voteSet, vals, time.Now())
		newValSet, _                  = randValidatorPrivValSet(ctx, t, 2, 1)
	)
	require.NoError(t, err)

	testCases := []struct {
		valSet *ValidatorSet
		err    bool
	}{
		// good
		0: {
			valSet: originalValset,
			err:    false,
		},
		// bad - no overlap between validator sets
		1: {
			valSet: newValSet,
			err:    true,
		},
		// good - first two are different but the rest of the same -> >1/3
		2: {
			valSet: NewValidatorSet(append(newValSet.Validators, originalValset.Validators...)),
			err:    false,
		},
	}

	for _, tc := range testCases {
		err = tc.valSet.VerifyCommitLightTrusting("test_chain_id", commit,
			libmath.Fraction{Numerator: 1, Denominator: 3})
		if tc.err {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}
	}
}

func TestValidatorSet_VerifyCommitLightTrustingErrorsOnOverflow(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var (
		blockID               = makeBlockIDRandom()
		voteSet, valSet, vals = randVoteSet(ctx, t, 1, 1, v1.PrecommitType, 1, MaxTotalVotingPower)
		commit, err           = makeCommit(ctx, blockID, 1, 1, voteSet, vals, time.Now())
	)
	require.NoError(t, err)

	err = valSet.VerifyCommitLightTrusting("test_chain_id", commit,
		libmath.Fraction{Numerator: 25, Denominator: 55})
	if assert.Error(t, err) {
		assert.Contains(t, err.Error(), "int64 overflow")
	}
}
