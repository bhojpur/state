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
	"crypto/sha256"
	"fmt"

	v1 "github.com/bhojpur/state/pkg/api/v1/crypto"
)

const ProofOpValue = "simple:v"

// ValueOp takes a key and a single value as argument and
// produces the root hash. The corresponding tree structure is
// the SimpleMap tree. SimpleMap takes a Hasher, and currently
// Bhojpur State uses statehash. SimpleValueOp should support
// the hash function as used in tmhash.  TODO support
// additional hash functions here as options/args to this
// operator.
//
// If the produced root hash matches the expected hash, the
// proof is good.
type ValueOp struct {
	// Encoded in ProofOp.Key.
	key []byte

	// To encode in ProofOp.Data
	Proof *Proof `json:"proof"`
}

var _ ProofOperator = ValueOp{}

func NewValueOp(key []byte, proof *Proof) ValueOp {
	return ValueOp{
		key:   key,
		Proof: proof,
	}
}

func ValueOpDecoder(pop v1.ProofOp) (ProofOperator, error) {
	if pop.Type != ProofOpValue {
		return nil, fmt.Errorf("unexpected ProofOp.Type; got %v, want %v", pop.Type, ProofOpValue)
	}
	var pbop v1.ValueOp // a bit strange as we'll discard this, but it works.
	err := pbop.Unmarshal(pop.Data)
	if err != nil {
		return nil, fmt.Errorf("decoding ProofOp.Data into ValueOp: %w", err)
	}

	sp, err := ProofFromProto(pbop.Proof)
	if err != nil {
		return nil, err
	}
	return NewValueOp(pop.Key, sp), nil
}

func (op ValueOp) ProofOp() v1.ProofOp {
	pbval := v1.ValueOp{
		Key:   op.key,
		Proof: op.Proof.ToProto(),
	}
	bz, err := pbval.Marshal()
	if err != nil {
		panic(err)
	}
	return v1.ProofOp{
		Type: ProofOpValue,
		Key:  op.key,
		Data: bz,
	}
}

func (op ValueOp) String() string {
	return fmt.Sprintf("ValueOp{%v}", op.GetKey())
}

func (op ValueOp) Run(args [][]byte) ([][]byte, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("expected 1 arg, got %v", len(args))
	}
	value := args[0]

	vhash := sha256.Sum256(value)

	bz := new(bytes.Buffer)
	// Wrap <op.Key, vhash> to hash the KVPair.
	encodeByteSlice(bz, op.key)   //nolint: errcheck // does not error
	encodeByteSlice(bz, vhash[:]) //nolint: errcheck // does not error
	kvhash := leafHash(bz.Bytes())

	if !bytes.Equal(kvhash, op.Proof.LeafHash) {
		return nil, fmt.Errorf("leaf hash mismatch: want %X got %X", op.Proof.LeafHash, kvhash)
	}

	return [][]byte{
		op.Proof.ComputeRootHash(),
	}, nil
}

func (op ValueOp) GetKey() []byte {
	return op.key
}
