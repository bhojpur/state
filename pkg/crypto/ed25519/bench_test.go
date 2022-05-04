package ed25519

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
	"fmt"
	"io"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bhojpur/state/pkg/crypto"
	"github.com/bhojpur/state/pkg/crypto/internal/benchmarking"
)

func BenchmarkKeyGeneration(b *testing.B) {
	benchmarkKeygenWrapper := func(reader io.Reader) crypto.PrivKey {
		return genPrivKey(reader)
	}
	benchmarking.BenchmarkKeyGeneration(b, benchmarkKeygenWrapper)
}

func BenchmarkSigning(b *testing.B) {
	priv := GenPrivKey()
	benchmarking.BenchmarkSigning(b, priv)
}

func BenchmarkVerification(b *testing.B) {
	priv := GenPrivKey()
	benchmarking.BenchmarkVerification(b, priv)
}

func BenchmarkVerifyBatch(b *testing.B) {
	msg := []byte("BatchVerifyTest")

	for _, sigsCount := range []int{1, 8, 64, 1024} {
		sigsCount := sigsCount
		b.Run(fmt.Sprintf("sig-count-%d", sigsCount), func(b *testing.B) {
			// Pre-generate all of the keys, and signatures, but do not
			// benchmark key-generation and signing.
			pubs := make([]crypto.PubKey, 0, sigsCount)
			sigs := make([][]byte, 0, sigsCount)
			for i := 0; i < sigsCount; i++ {
				priv := GenPrivKey()
				sig, _ := priv.Sign(msg)
				pubs = append(pubs, priv.PubKey().(PubKey))
				sigs = append(sigs, sig)
			}
			b.ResetTimer()

			b.ReportAllocs()
			// NOTE: dividing by n so that metrics are per-signature
			for i := 0; i < b.N/sigsCount; i++ {
				// The benchmark could just benchmark the Verify()
				// routine, but there is non-trivial overhead associated
				// with BatchVerifier.Add(), which should be included
				// in the benchmark.
				v := NewBatchVerifier()
				for i := 0; i < sigsCount; i++ {
					err := v.Add(pubs[i], msg, sigs[i])
					require.NoError(b, err)
				}

				if ok, _ := v.Verify(); !ok {
					b.Fatal("signature set failed batch verification")
				}
			}
		})
	}
}
