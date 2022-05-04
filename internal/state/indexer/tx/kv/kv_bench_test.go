package kv

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
	"crypto/rand"
	"fmt"
	"testing"

	dbm "github.com/bhojpur/state/pkg/database"

	"github.com/bhojpur/state/internal/pubsub/query"
	abcipb "github.com/bhojpur/state/pkg/abci/types"
	"github.com/bhojpur/state/pkg/types"
)

func BenchmarkTxSearch(b *testing.B) {
	dbDir := b.TempDir()

	db, err := dbm.NewGoLevelDB("benchmark_tx_search_test", dbDir)
	if err != nil {
		b.Errorf("failed to create database: %s", err)
	}

	indexer := NewTxIndex(db)

	for i := 0; i < 35000; i++ {
		events := []abcipb.Event{
			{
				Type: "transfer",
				Attributes: []abcipb.EventAttribute{
					{Key: "address", Value: fmt.Sprintf("address_%d", i%100), Index: true},
					{Key: "amount", Value: "50", Index: true},
				},
			},
		}

		txBz := make([]byte, 8)
		if _, err := rand.Read(txBz); err != nil {
			b.Errorf("failed produce random bytes: %s", err)
		}

		txResult := &abcipb.TxResult{
			Height: int64(i),
			Index:  0,
			Tx:     types.Tx(string(txBz)),
			Result: abcipb.ExecTxResult{
				Data:   []byte{0},
				Code:   abcipb.CodeTypeOK,
				Log:    "",
				Events: events,
			},
		}

		if err := indexer.Index([]*abcipb.TxResult{txResult}); err != nil {
			b.Errorf("failed to index tx: %s", err)
		}
	}

	txQuery := query.MustCompile(`transfer.address = 'address_43' AND transfer.amount = 50`)

	b.ResetTimer()

	ctx := context.Background()

	for i := 0; i < b.N; i++ {
		if _, err := indexer.Search(ctx, txQuery); err != nil {
			b.Errorf("failed to query for txs: %s", err)
		}
	}
}
