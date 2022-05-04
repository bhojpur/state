package crypto_test

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
)

// the purpose of this test is primarily to ensure that the randomness
// generation won't error.
func TestRandomConsistency(t *testing.T) {
	x1 := crypto.CRandBytes(256)
	x2 := crypto.CRandBytes(256)
	x3 := crypto.CRandBytes(256)
	x4 := crypto.CRandBytes(256)
	x5 := crypto.CRandBytes(256)
	require.NotEqual(t, x1, x2)
	require.NotEqual(t, x3, x4)
	require.NotEqual(t, x4, x5)
	require.NotEqual(t, x1, x5)
}
