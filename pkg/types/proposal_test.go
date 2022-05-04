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
	"math"
	"testing"
	"time"

	"github.com/gogo/protobuf/proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/bhojpur/state/internal/libs/protoio"
	v1 "github.com/bhojpur/state/pkg/api/v1/types"
	"github.com/bhojpur/state/pkg/crypto"
	librand "github.com/bhojpur/state/pkg/libs/rand"
	libtime "github.com/bhojpur/state/pkg/libs/time"
)

func getTestProposal(t testing.TB) *Proposal {
	t.Helper()

	stamp, err := time.Parse(TimeFormat, "2018-02-11T07:09:22.765Z")
	require.NoError(t, err)

	return &Proposal{
		Height: 12345,
		Round:  23456,
		BlockID: BlockID{Hash: []byte("--June_15_2020_amino_was_removed"),
			PartSetHeader: PartSetHeader{Total: 111, Hash: []byte("--June_15_2020_amino_was_removed")}},
		POLRound:  -1,
		Timestamp: stamp,
	}
}

func TestProposalSignable(t *testing.T) {
	chainID := "test_chain_id"
	signBytes := ProposalSignBytes(chainID, getTestProposal(t).ToProto())
	pb := CanonicalizeProposal(chainID, getTestProposal(t).ToProto())

	expected, err := protoio.MarshalDelimited(&pb)
	require.NoError(t, err)
	require.Equal(t, expected, signBytes, "Got unexpected sign bytes for Proposal")
}

func TestProposalString(t *testing.T) {
	str := getTestProposal(t).String()
	expected := `Proposal{12345/23456 (2D2D4A756E655F31355F323032305F616D696E6F5F7761735F72656D6F766564:111:2D2D4A756E65, -1) 000000000000 @ 2018-02-11T07:09:22.765Z}`
	if str != expected {
		t.Errorf("got unexpected string for Proposal. Expected:\n%v\nGot:\n%v", expected, str)
	}
}

func TestProposalVerifySignature(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	privVal := NewMockPV()
	pubKey, err := privVal.GetPubKey(ctx)
	require.NoError(t, err)

	prop := NewProposal(
		4, 2, 2,
		BlockID{librand.Bytes(crypto.HashSize), PartSetHeader{777, librand.Bytes(crypto.HashSize)}}, libtime.Now())
	p := prop.ToProto()
	signBytes := ProposalSignBytes("test_chain_id", p)

	// sign it
	err = privVal.SignProposal(ctx, "test_chain_id", p)
	require.NoError(t, err)
	prop.Signature = p.Signature

	// verify the same proposal
	valid := pubKey.VerifySignature(signBytes, prop.Signature)
	require.True(t, valid)

	// serialize, deserialize and verify again....
	newProp := new(v1.Proposal)
	pb := prop.ToProto()

	bs, err := proto.Marshal(pb)
	require.NoError(t, err)

	err = proto.Unmarshal(bs, newProp)
	require.NoError(t, err)

	np, err := ProposalFromProto(newProp)
	require.NoError(t, err)

	// verify the transmitted proposal
	newSignBytes := ProposalSignBytes("test_chain_id", pb)
	require.Equal(t, string(signBytes), string(newSignBytes))
	valid = pubKey.VerifySignature(newSignBytes, np.Signature)
	require.True(t, valid)
}

func BenchmarkProposalWriteSignBytes(b *testing.B) {
	pbp := getTestProposal(b).ToProto()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ProposalSignBytes("test_chain_id", pbp)
	}
}

func BenchmarkProposalSign(b *testing.B) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	privVal := NewMockPV()

	pbp := getTestProposal(b).ToProto()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		err := privVal.SignProposal(ctx, "test_chain_id", pbp)
		if err != nil {
			b.Error(err)
		}
	}
}

func BenchmarkProposalVerifySignature(b *testing.B) {
	testProposal := getTestProposal(b)
	pbp := testProposal.ToProto()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	privVal := NewMockPV()
	err := privVal.SignProposal(ctx, "test_chain_id", pbp)
	require.NoError(b, err)
	pubKey, err := privVal.GetPubKey(ctx)
	require.NoError(b, err)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		pubKey.VerifySignature(ProposalSignBytes("test_chain_id", pbp), testProposal.Signature)
	}
}

func TestProposalValidateBasic(t *testing.T) {

	privVal := NewMockPV()
	testCases := []struct {
		testName         string
		malleateProposal func(*Proposal)
		expectErr        bool
	}{
		{"Good Proposal", func(p *Proposal) {}, false},
		{"Invalid Type", func(p *Proposal) { p.Type = v1.PrecommitType }, true},
		{"Invalid Height", func(p *Proposal) { p.Height = -1 }, true},
		{"Invalid Round", func(p *Proposal) { p.Round = -1 }, true},
		{"Invalid POLRound", func(p *Proposal) { p.POLRound = -2 }, true},
		{"Invalid BlockId", func(p *Proposal) {
			p.BlockID = BlockID{[]byte{1, 2, 3}, PartSetHeader{111, []byte("blockparts")}}
		}, true},
		{"Invalid Signature", func(p *Proposal) {
			p.Signature = make([]byte, 0)
		}, true},
		{"Too big Signature", func(p *Proposal) {
			p.Signature = make([]byte, MaxSignatureSize+1)
		}, true},
	}
	blockID := makeBlockID(crypto.Checksum([]byte("blockhash")), math.MaxInt32, crypto.Checksum([]byte("partshash")))

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.testName, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			prop := NewProposal(
				4, 2, 2,
				blockID, libtime.Now())
			p := prop.ToProto()
			err := privVal.SignProposal(ctx, "test_chain_id", p)
			prop.Signature = p.Signature
			require.NoError(t, err)
			tc.malleateProposal(prop)
			assert.Equal(t, tc.expectErr, prop.ValidateBasic() != nil, "Validate Basic had an unexpected result")
		})
	}
}

func TestProposalProtoBuf(t *testing.T) {
	proposal := NewProposal(1, 2, 3, makeBlockID([]byte("hash"), 2, []byte("part_set_hash")), libtime.Now())
	proposal.Signature = []byte("sig")
	proposal2 := NewProposal(1, 2, 3, BlockID{}, libtime.Now())

	testCases := []struct {
		msg     string
		p1      *Proposal
		expPass bool
	}{
		{"success", proposal, true},
		{"success", proposal2, false}, // blcokID cannot be empty
		{"empty proposal failure validatebasic", &Proposal{}, false},
		{"nil proposal", nil, false},
	}
	for _, tc := range testCases {
		protoProposal := tc.p1.ToProto()

		p, err := ProposalFromProto(protoProposal)
		if tc.expPass {
			require.NoError(t, err)
			require.Equal(t, tc.p1, p, tc.msg)
		} else {
			require.Error(t, err)
		}
	}
}

func TestIsTimely(t *testing.T) {
	genesisTime, err := time.Parse(time.RFC3339, "2019-03-13T23:00:00Z")
	require.NoError(t, err)
	testCases := []struct {
		name         string
		proposalTime time.Time
		recvTime     time.Time
		precision    time.Duration
		msgDelay     time.Duration
		expectTimely bool
		round        int32
	}{
		// proposalTime - precision <= localTime <= proposalTime + msgDelay + precision
		{
			// Checking that the following inequality evaluates to true:
			// 0 - 2 <= 1 <= 0 + 1 + 2
			name:         "basic timely",
			proposalTime: genesisTime,
			recvTime:     genesisTime.Add(1 * time.Nanosecond),
			precision:    time.Nanosecond * 2,
			msgDelay:     time.Nanosecond,
			expectTimely: true,
		},
		{
			// Checking that the following inequality evaluates to false:
			// 0 - 2 <= 4 <= 0 + 1 + 2
			name:         "local time too large",
			proposalTime: genesisTime,
			recvTime:     genesisTime.Add(4 * time.Nanosecond),
			precision:    time.Nanosecond * 2,
			msgDelay:     time.Nanosecond,
			expectTimely: false,
		},
		{
			// Checking that the following inequality evaluates to false:
			// 4 - 2 <= 0 <= 4 + 2 + 1
			name:         "proposal time too large",
			proposalTime: genesisTime.Add(4 * time.Nanosecond),
			recvTime:     genesisTime,
			precision:    time.Nanosecond * 2,
			msgDelay:     time.Nanosecond,
			expectTimely: false,
		},
		{
			// Checking that the following inequality evaluates to true:
			// 0 - (2 * 2)  <= 4 <= 0 + (1 * 2) + 2
			name:         "message delay adapts after 10 rounds",
			proposalTime: genesisTime,
			recvTime:     genesisTime.Add(4 * time.Nanosecond),
			precision:    time.Nanosecond * 2,
			msgDelay:     time.Nanosecond,
			expectTimely: true,
			round:        10,
		},
		{
			// check that values that overflow time.Duration still correctly register
			// as timely when round relaxation applied.
			name:         "message delay fixed to not overflow time.Duration",
			proposalTime: genesisTime,
			recvTime:     genesisTime.Add(4 * time.Nanosecond),
			precision:    time.Nanosecond * 2,
			msgDelay:     time.Nanosecond,
			expectTimely: true,
			round:        5000,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			p := Proposal{
				Timestamp: testCase.proposalTime,
			}

			sp := SynchronyParams{
				Precision:    testCase.precision,
				MessageDelay: testCase.msgDelay,
			}

			ti := p.IsTimely(testCase.recvTime, sp, testCase.round)
			assert.Equal(t, testCase.expectTimely, ti)
		})
	}
}
