package merkle

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
	"errors"
	"fmt"

	v1 "github.com/bhojpur/state/pkg/api/v1/crypto"
)

// ProofOp gets converted to an instance of ProofOperator:

// ProofOperator is a layer for calculating intermediate Merkle roots
// when a series of Merkle trees are chained together.
// Run() takes leaf values from a tree and returns the Merkle
// root for the corresponding tree. It takes and returns a list of bytes
// to allow multiple leaves to be part of a single proof, for instance in a range proof.
// ProofOp() encodes the ProofOperator in a generic way so it can later be
// decoded with OpDecoder.
type ProofOperator interface {
	Run([][]byte) ([][]byte, error)
	GetKey() []byte
	ProofOp() v1.ProofOp
}

// Operations on a list of ProofOperators

// ProofOperators is a slice of ProofOperator(s).
// Each operator will be applied to the input value sequentially
// and the last Merkle root will be verified with already known data
type ProofOperators []ProofOperator

func (poz ProofOperators) VerifyValue(root []byte, keypath string, value []byte) (err error) {
	return poz.Verify(root, keypath, [][]byte{value})
}

func (poz ProofOperators) Verify(root []byte, keypath string, args [][]byte) (err error) {
	keys, err := KeyPathToKeys(keypath)
	if err != nil {
		return
	}

	for i, op := range poz {
		key := op.GetKey()
		if len(key) != 0 {
			if len(keys) == 0 {
				return fmt.Errorf("key path has insufficient # of parts: expected no more keys but got %+v", string(key))
			}
			lastKey := keys[len(keys)-1]
			if !bytes.Equal(lastKey, key) {
				return fmt.Errorf("key mismatch on operation #%d: expected %+v but got %+v", i, string(lastKey), string(key))
			}
			keys = keys[:len(keys)-1]
		}
		args, err = op.Run(args)
		if err != nil {
			return
		}
	}
	if !bytes.Equal(root, args[0]) {
		return fmt.Errorf("calculated root hash is invalid: expected %X but got %X", root, args[0])
	}
	if len(keys) != 0 {
		return errors.New("keypath not consumed all")
	}
	return nil
}

// ProofRuntime - main entrypoint

type OpDecoder func(v1.ProofOp) (ProofOperator, error)

type ProofRuntime struct {
	decoders map[string]OpDecoder
}

func NewProofRuntime() *ProofRuntime {
	return &ProofRuntime{
		decoders: make(map[string]OpDecoder),
	}
}

func (prt *ProofRuntime) RegisterOpDecoder(typ string, dec OpDecoder) {
	_, ok := prt.decoders[typ]
	if ok {
		panic("already registered for type " + typ)
	}
	prt.decoders[typ] = dec
}

func (prt *ProofRuntime) Decode(pop v1.ProofOp) (ProofOperator, error) {
	decoder := prt.decoders[pop.Type]
	if decoder == nil {
		return nil, fmt.Errorf("unrecognized proof type %v", pop.Type)
	}
	return decoder(pop)
}

func (prt *ProofRuntime) DecodeProof(proof *v1.ProofOps) (ProofOperators, error) {
	poz := make(ProofOperators, 0, len(proof.Ops))
	for _, pop := range proof.Ops {
		operator, err := prt.Decode(pop)
		if err != nil {
			return nil, fmt.Errorf("decoding a proof operator: %w", err)
		}
		poz = append(poz, operator)
	}
	return poz, nil
}

func (prt *ProofRuntime) VerifyValue(proof *v1.ProofOps, root []byte, keypath string, value []byte) (err error) {
	return prt.Verify(proof, root, keypath, [][]byte{value})
}

// TODO In the long run we'll need a method of classifcation of ops,
// whether existence or absence or perhaps a third?
func (prt *ProofRuntime) VerifyAbsence(proof *v1.ProofOps, root []byte, keypath string) (err error) {
	return prt.Verify(proof, root, keypath, nil)
}

func (prt *ProofRuntime) Verify(proof *v1.ProofOps, root []byte, keypath string, args [][]byte) (err error) {
	poz, err := prt.DecodeProof(proof)
	if err != nil {
		return fmt.Errorf("decoding proof: %w", err)
	}
	return poz.Verify(root, keypath, args)
}

// DefaultProofRuntime only knows about value proofs.
// To use e.g. IAVL proofs, register op-decoders as
// defined in the IAVL package.
func DefaultProofRuntime() (prt *ProofRuntime) {
	prt = NewProofRuntime()
	prt.RegisterOpDecoder(ProofOpValue, ValueOpDecoder)
	return
}
