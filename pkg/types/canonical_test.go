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
	"reflect"
	"testing"

	v1 "github.com/bhojpur/state/pkg/api/v1/types"
	"github.com/bhojpur/state/pkg/crypto"
	librand "github.com/bhojpur/state/pkg/libs/rand"
)

func TestCanonicalizeBlockID(t *testing.T) {
	randhash := librand.Bytes(crypto.HashSize)
	block1 := v1.BlockID{Hash: randhash,
		PartSetHeader: v1.PartSetHeader{Total: 5, Hash: randhash}}
	block2 := v1.BlockID{Hash: randhash,
		PartSetHeader: v1.PartSetHeader{Total: 10, Hash: randhash}}
	cblock1 := v1.CanonicalBlockID{Hash: randhash,
		PartSetHeader: v1.CanonicalPartSetHeader{Total: 5, Hash: randhash}}
	cblock2 := v1.CanonicalBlockID{Hash: randhash,
		PartSetHeader: v1.CanonicalPartSetHeader{Total: 10, Hash: randhash}}

	tests := []struct {
		name string
		args v1.BlockID
		want *v1.CanonicalBlockID
	}{
		{"first", block1, &cblock1},
		{"second", block2, &cblock2},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if got := CanonicalizeBlockID(tt.args); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CanonicalizeBlockID() = %v, want %v", got, tt.want)
			}
		})
	}
}
