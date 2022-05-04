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
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bhojpur/state/pkg/crypto"
	librand "github.com/bhojpur/state/pkg/libs/rand"
)

func TestBlockMeta_ToProto(t *testing.T) {
	h := MakeRandHeader()
	bi := BlockID{Hash: h.Hash(), PartSetHeader: PartSetHeader{Total: 123, Hash: librand.Bytes(crypto.HashSize)}}

	bm := &BlockMeta{
		BlockID:   bi,
		BlockSize: 200,
		Header:    h,
		NumTxs:    0,
	}

	tests := []struct {
		testName string
		bm       *BlockMeta
		expErr   bool
	}{
		{"success", bm, false},
		{"failure nil", nil, true},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.testName, func(t *testing.T) {
			pb := tt.bm.ToProto()

			bm, err := BlockMetaFromProto(pb)

			if !tt.expErr {
				require.NoError(t, err, tt.testName)
				require.Equal(t, tt.bm, bm, tt.testName)
			} else {
				require.Error(t, err, tt.testName)
			}
		})
	}
}

func TestBlockMeta_ValidateBasic(t *testing.T) {
	h := MakeRandHeader()
	bi := BlockID{Hash: h.Hash(), PartSetHeader: PartSetHeader{Total: 123, Hash: librand.Bytes(crypto.HashSize)}}
	bi2 := BlockID{Hash: librand.Bytes(crypto.HashSize),
		PartSetHeader: PartSetHeader{Total: 123, Hash: librand.Bytes(crypto.HashSize)}}
	bi3 := BlockID{Hash: []byte("incorrect hash"),
		PartSetHeader: PartSetHeader{Total: 123, Hash: []byte("incorrect hash")}}

	bm := &BlockMeta{
		BlockID:   bi,
		BlockSize: 200,
		Header:    h,
		NumTxs:    0,
	}

	bm2 := &BlockMeta{
		BlockID:   bi2,
		BlockSize: 200,
		Header:    h,
		NumTxs:    0,
	}

	bm3 := &BlockMeta{
		BlockID:   bi3,
		BlockSize: 200,
		Header:    h,
		NumTxs:    0,
	}

	tests := []struct {
		name    string
		bm      *BlockMeta
		wantErr bool
	}{
		{"success", bm, false},
		{"failure wrong blockID hash", bm2, true},
		{"failure wrong length blockID hash", bm3, true},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.bm.ValidateBasic(); (err != nil) != tt.wantErr {
				t.Errorf("BlockMeta.ValidateBasic() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
