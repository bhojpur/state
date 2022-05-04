package types_test

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

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	abci "github.com/bhojpur/state/pkg/abci/types"
	"github.com/bhojpur/state/pkg/crypto/merkle"
)

func TestHashAndProveResults(t *testing.T) {
	trs := []*abci.ExecTxResult{
		// Note, these tests rely on the first two entries being in this order.
		{Code: 0, Data: nil},
		{Code: 0, Data: []byte{}},

		{Code: 0, Data: []byte("one")},
		{Code: 14, Data: nil},
		{Code: 14, Data: []byte("foo")},
		{Code: 14, Data: []byte("bar")},
	}

	// Nil and []byte{} should produce the same bytes
	bz0, err := trs[0].Marshal()
	require.NoError(t, err)
	bz1, err := trs[1].Marshal()
	require.NoError(t, err)
	require.Equal(t, bz0, bz1)

	// Make sure that we can get a root hash from results and verify proofs.
	rs, err := abci.MarshalTxResults(trs)
	require.NoError(t, err)
	root := merkle.HashFromByteSlices(rs)
	assert.NotEmpty(t, root)

	_, proofs := merkle.ProofsFromByteSlices(rs)
	for i, tr := range trs {
		bz, err := tr.Marshal()
		require.NoError(t, err)

		valid := proofs[i].Verify(root, bz)
		assert.NoError(t, valid, "%d", i)
	}
}

func TestHashDeterministicFieldsOnly(t *testing.T) {
	tr1 := abci.ExecTxResult{
		Code:      1,
		Data:      []byte("transaction"),
		Log:       "nondeterministic data: abc",
		Info:      "nondeterministic data: abc",
		GasWanted: 1000,
		GasUsed:   1000,
		Events:    []abci.Event{},
		Codespace: "nondeterministic.data.abc",
	}
	tr2 := abci.ExecTxResult{
		Code:      1,
		Data:      []byte("transaction"),
		Log:       "nondeterministic data: def",
		Info:      "nondeterministic data: def",
		GasWanted: 1000,
		GasUsed:   1000,
		Events:    []abci.Event{},
		Codespace: "nondeterministic.data.def",
	}
	r1, err := abci.MarshalTxResults([]*abci.ExecTxResult{&tr1})
	require.NoError(t, err)
	r2, err := abci.MarshalTxResults([]*abci.ExecTxResult{&tr2})
	require.NoError(t, err)
	require.Equal(t, merkle.HashFromByteSlices(r1), merkle.HashFromByteSlices(r2))
}
