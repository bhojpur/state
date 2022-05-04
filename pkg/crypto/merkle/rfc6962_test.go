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
	"encoding/hex"
	"testing"

	"github.com/bhojpur/state/pkg/crypto"
)

func TestRFC6962Hasher(t *testing.T) {
	_, leafHashTrail := trailsFromByteSlices([][]byte{[]byte("L123456")})
	leafHash := leafHashTrail.Hash
	_, leafHashTrail = trailsFromByteSlices([][]byte{{}})
	emptyLeafHash := leafHashTrail.Hash
	_, emptyHashTrail := trailsFromByteSlices([][]byte{})
	emptyTreeHash := emptyHashTrail.Hash
	for _, tc := range []struct {
		desc string
		got  []byte
		want string
	}{
		// Check that empty trees return the hash of an empty string.
		// echo -n '' | sha256sum
		{
			desc: "RFC6962 Empty Tree",
			want: "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"[:crypto.HashSize*2],
			got:  emptyTreeHash,
		},

		// Check that the empty hash is not the same as the hash of an empty leaf.
		// echo -n 00 | xxd -r -p | sha256sum
		{
			desc: "RFC6962 Empty Leaf",
			want: "6e340b9cffb37a989ca544e6bb780a2c78901d3fb33738768511a30617afa01d"[:crypto.HashSize*2],
			got:  emptyLeafHash,
		},
		// echo -n 004C313233343536 | xxd -r -p | sha256sum
		{
			desc: "RFC6962 Leaf",
			want: "395aa064aa4c29f7010acfe3f25db9485bbd4b91897b6ad7ad547639252b4d56"[:crypto.HashSize*2],
			got:  leafHash,
		},
		// echo -n 014E3132334E343536 | xxd -r -p | sha256sum
		{
			desc: "RFC6962 Node",
			want: "aa217fe888e47007fa15edab33c2b492a722cb106c64667fc2b044444de66bbb"[:crypto.HashSize*2],
			got:  innerHash([]byte("N123"), []byte("N456")),
		},
	} {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			wantBytes, err := hex.DecodeString(tc.want)
			if err != nil {
				t.Fatalf("hex.DecodeString(%x): %v", tc.want, err)
			}
			if got, want := tc.got, wantBytes; !bytes.Equal(got, want) {
				t.Errorf("got %x, want %x", got, want)
			}
		})
	}
}

func TestRFC6962HasherCollisions(t *testing.T) {
	// Check that different leaves have different hashes.
	leaf1, leaf2 := []byte("Hello"), []byte("World")
	_, leafHashTrail := trailsFromByteSlices([][]byte{leaf1})
	hash1 := leafHashTrail.Hash
	_, leafHashTrail = trailsFromByteSlices([][]byte{leaf2})
	hash2 := leafHashTrail.Hash
	if bytes.Equal(hash1, hash2) {
		t.Errorf("leaf hashes should differ, but both are %x", hash1)
	}
	// Compute an intermediate subtree hash.
	_, subHash1Trail := trailsFromByteSlices([][]byte{hash1, hash2})
	subHash1 := subHash1Trail.Hash
	// Check that this is not the same as a leaf hash of their concatenation.
	preimage := append(hash1, hash2...)
	_, forgedHashTrail := trailsFromByteSlices([][]byte{preimage})
	forgedHash := forgedHashTrail.Hash
	if bytes.Equal(subHash1, forgedHash) {
		t.Errorf("hasher is not second-preimage resistant")
	}
	// Swap the order of nodes and check that the hash is different.
	_, subHash2Trail := trailsFromByteSlices([][]byte{hash2, hash1})
	subHash2 := subHash2Trail.Hash
	if bytes.Equal(subHash1, subHash2) {
		t.Errorf("subtree hash does not depend on the order of leaves")
	}
}
