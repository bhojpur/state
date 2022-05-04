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
	"crypto/rand"
	"fmt"

	"github.com/oasisprotocol/curve25519-voi/primitives/sr25519"

	"github.com/bhojpur/state/pkg/crypto"
)

var _ crypto.BatchVerifier = &BatchVerifier{}

// BatchVerifier implements batch verification for sr25519.
type BatchVerifier struct {
	*sr25519.BatchVerifier
}

func NewBatchVerifier() crypto.BatchVerifier {
	return &BatchVerifier{sr25519.NewBatchVerifier()}
}

func (b *BatchVerifier) Add(key crypto.PubKey, msg, signature []byte) error {
	pk, ok := key.(PubKey)
	if !ok {
		return fmt.Errorf("sr25519: pubkey is not sr25519")
	}

	var srpk sr25519.PublicKey
	if err := srpk.UnmarshalBinary(pk); err != nil {
		return fmt.Errorf("sr25519: invalid public key: %w", err)
	}

	var sig sr25519.Signature
	if err := sig.UnmarshalBinary(signature); err != nil {
		return fmt.Errorf("sr25519: unable to decode signature: %w", err)
	}

	st := signingCtx.NewTranscriptBytes(msg)
	b.BatchVerifier.Add(&srpk, st, &sig)

	return nil
}

func (b *BatchVerifier) Verify() (bool, []bool) {
	return b.BatchVerifier.Verify(rand.Reader)
}
