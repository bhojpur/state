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
	"bytes"
	"context"
	"errors"
	"fmt"

	v1 "github.com/bhojpur/state/pkg/api/v1/types"
	"github.com/bhojpur/state/pkg/crypto"
	"github.com/bhojpur/state/pkg/crypto/ed25519"
)

// PrivValidatorType defines the implemtation types.
type PrivValidatorType uint8

const (
	MockSignerClient      = PrivValidatorType(0x00) // mock signer
	FileSignerClient      = PrivValidatorType(0x01) // signer client via file
	RetrySignerClient     = PrivValidatorType(0x02) // signer client with retry via socket
	SignerSocketClient    = PrivValidatorType(0x03) // signer client via socket
	ErrorMockSignerClient = PrivValidatorType(0x04) // error mock signer
	SignerGRPCClient      = PrivValidatorType(0x05) // signer client via gRPC
)

// PrivValidator defines the functionality of a local Bhojpur State validator
// that signs votes and proposals, and never double signs.
type PrivValidator interface {
	GetPubKey(context.Context) (crypto.PubKey, error)

	SignVote(ctx context.Context, chainID string, vote *v1.Vote) error
	SignProposal(ctx context.Context, chainID string, proposal *v1.Proposal) error
}

type PrivValidatorsByAddress []PrivValidator

func (pvs PrivValidatorsByAddress) Len() int {
	return len(pvs)
}

func (pvs PrivValidatorsByAddress) Less(i, j int) bool {
	pvi, err := pvs[i].GetPubKey(context.TODO())
	if err != nil {
		panic(err)
	}
	pvj, err := pvs[j].GetPubKey(context.TODO())
	if err != nil {
		panic(err)
	}

	return bytes.Compare(pvi.Address(), pvj.Address()) == -1
}

func (pvs PrivValidatorsByAddress) Swap(i, j int) {
	pvs[i], pvs[j] = pvs[j], pvs[i]
}

// MockPV

// MockPV implements PrivValidator without any safety or persistence.
// Only use it for testing.
type MockPV struct {
	PrivKey              crypto.PrivKey
	breakProposalSigning bool
	breakVoteSigning     bool
}

func NewMockPV() MockPV {
	return MockPV{ed25519.GenPrivKey(), false, false}
}

// NewMockPVWithParams allows one to create a MockPV instance, but with finer
// grained control over the operation of the mock validator. This is useful for
// mocking test failures.
func NewMockPVWithParams(privKey crypto.PrivKey, breakProposalSigning, breakVoteSigning bool) MockPV {
	return MockPV{privKey, breakProposalSigning, breakVoteSigning}
}

// Implements PrivValidator.
func (pv MockPV) GetPubKey(ctx context.Context) (crypto.PubKey, error) {
	return pv.PrivKey.PubKey(), nil
}

// Implements PrivValidator.
func (pv MockPV) SignVote(ctx context.Context, chainID string, vote *v1.Vote) error {
	useChainID := chainID
	if pv.breakVoteSigning {
		useChainID = "incorrect-chain-id"
	}

	signBytes := VoteSignBytes(useChainID, vote)
	extSignBytes := VoteExtensionSignBytes(useChainID, vote)
	sig, err := pv.PrivKey.Sign(signBytes)
	if err != nil {
		return err
	}
	vote.Signature = sig

	var extSig []byte
	// We only sign vote extensions for precommits
	if vote.Type == v1.PrecommitType {
		extSig, err = pv.PrivKey.Sign(extSignBytes)
		if err != nil {
			return err
		}
	} else if len(vote.Extension) > 0 {
		return errors.New("unexpected vote extension - vote extensions are only allowed in precommits")
	}
	vote.ExtensionSignature = extSig
	return nil
}

// Implements PrivValidator.
func (pv MockPV) SignProposal(ctx context.Context, chainID string, proposal *v1.Proposal) error {
	useChainID := chainID
	if pv.breakProposalSigning {
		useChainID = "incorrect-chain-id"
	}

	signBytes := ProposalSignBytes(useChainID, proposal)
	sig, err := pv.PrivKey.Sign(signBytes)
	if err != nil {
		return err
	}
	proposal.Signature = sig
	return nil
}

func (pv MockPV) ExtractIntoValidator(ctx context.Context, votingPower int64) *Validator {
	pubKey, _ := pv.GetPubKey(ctx)
	return &Validator{
		Address:     pubKey.Address(),
		PubKey:      pubKey,
		VotingPower: votingPower,
	}
}

// String returns a string representation of the MockPV.
func (pv MockPV) String() string {
	mpv, _ := pv.GetPubKey(context.TODO()) // mockPV will never return an error, ignored here
	return fmt.Sprintf("MockPV{%v}", mpv.Address())
}

// XXX: Implement.
func (pv MockPV) DisableChecks() {
	// Currently this does nothing,
	// as MockPV has no safety checks at all.
}

type ErroringMockPV struct {
	MockPV
}

var ErroringMockPVErr = errors.New("erroringMockPV always returns an error")

// Implements PrivValidator.
func (pv *ErroringMockPV) GetPubKey(ctx context.Context) (crypto.PubKey, error) {
	return nil, ErroringMockPVErr
}

// Implements PrivValidator.
func (pv *ErroringMockPV) SignVote(ctx context.Context, chainID string, vote *v1.Vote) error {
	return ErroringMockPVErr
}

// Implements PrivValidator.
func (pv *ErroringMockPV) SignProposal(ctx context.Context, chainID string, proposal *v1.Proposal) error {
	return ErroringMockPVErr
}

// NewErroringMockPV returns a MockPV that fails on each signing request. Again, for testing only.

func NewErroringMockPV() *ErroringMockPV {
	return &ErroringMockPV{MockPV{ed25519.GenPrivKey(), false, false}}
}
