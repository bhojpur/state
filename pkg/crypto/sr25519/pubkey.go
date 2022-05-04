package sr25519

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
	"fmt"

	"github.com/oasisprotocol/curve25519-voi/primitives/sr25519"

	"github.com/bhojpur/state/pkg/crypto"
)

var _ crypto.PubKey = PubKey{}

const (
	// PubKeySize is the size of a sr25519 public key in bytes.
	PubKeySize = sr25519.PublicKeySize

	// SignatureSize is the size of a sr25519 signature in bytes.
	SignatureSize = sr25519.SignatureSize
)

// PubKey implements crypto.PubKey.
type PubKey []byte

// TypeTag satisfies the jsontypes.Tagged interface.
func (PubKey) TypeTag() string { return PubKeyName }

// Address is the SHA256-20 of the raw pubkey bytes.
func (pubKey PubKey) Address() crypto.Address {
	if len(pubKey) != PubKeySize {
		panic("pubkey is incorrect size")
	}
	return crypto.AddressHash(pubKey)
}

// Bytes returns the PubKey byte format.
func (pubKey PubKey) Bytes() []byte {
	return []byte(pubKey)
}

func (pubKey PubKey) Equals(other crypto.PubKey) bool {
	if otherSr, ok := other.(PubKey); ok {
		return bytes.Equal(pubKey[:], otherSr[:])
	}

	return false
}

func (pubKey PubKey) VerifySignature(msg []byte, sigBytes []byte) bool {
	var srpk sr25519.PublicKey
	if err := srpk.UnmarshalBinary(pubKey); err != nil {
		return false
	}

	var sig sr25519.Signature
	if err := sig.UnmarshalBinary(sigBytes); err != nil {
		return false
	}

	st := signingCtx.NewTranscriptBytes(msg)
	return srpk.Verify(st, &sig)
}

func (pubKey PubKey) Type() string {
	return KeyType
}

func (pubKey PubKey) String() string {
	return fmt.Sprintf("PubKeySr25519{%X}", []byte(pubKey))
}
