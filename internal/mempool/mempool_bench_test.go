package mempool

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
	"context"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/bhojpur/state/example/kvstore"
	abciclient "github.com/bhojpur/state/pkg/abci/client"
	"github.com/bhojpur/state/pkg/libs/log"
)

func BenchmarkTxMempool_CheckTx(b *testing.B) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	client := abciclient.NewLocalClient(log.NewNopLogger(), kvstore.NewApplication())
	if err := client.Start(ctx); err != nil {
		b.Fatal(err)
	}

	// setup the cache and the mempool number for hitting GetEvictableTxs during the
	// benchmark. 5000 is the current default mempool size in the TM config.
	txmp := setup(b, client, 10000)
	txmp.config.Size = 5000

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	const peerID = 1

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		b.StopTimer()
		prefix := make([]byte, 20)
		_, err := rng.Read(prefix)
		require.NoError(b, err)

		priority := int64(rng.Intn(9999-1000) + 1000)
		tx := []byte(fmt.Sprintf("sender-%d-%d=%X=%d", n, peerID, prefix, priority))
		txInfo := TxInfo{SenderID: uint16(peerID)}

		b.StartTimer()

		require.NoError(b, txmp.CheckTx(ctx, tx, nil, txInfo))
	}
}
