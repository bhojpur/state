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
	// it is ok to use math/rand here: we do not need a cryptographically secure random
	// number generator here and we can run the tests a bit faster
	"math/rand"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestKeyPath(t *testing.T) {
	var path KeyPath
	keys := make([][]byte, 10)
	alphanum := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	for d := 0; d < 1e4; d++ {
		path = nil

		for i := range keys {
			enc := keyEncoding(rand.Intn(int(KeyEncodingMax)))
			keys[i] = make([]byte, rand.Uint32()%20)
			switch enc {
			case KeyEncodingURL:
				for j := range keys[i] {
					keys[i][j] = alphanum[rand.Intn(len(alphanum))]
				}
			case KeyEncodingHex:
				rand.Read(keys[i])
			default:
				require.Fail(t, "Unexpected encoding")
			}
			path = path.AppendKey(keys[i], enc)
		}

		res, err := KeyPathToKeys(path.String())
		require.NoError(t, err)
		require.Equal(t, len(keys), len(res))

		for i, key := range keys {
			require.Equal(t, key, res[i])
		}
	}
}
