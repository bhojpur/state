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
	"fmt"
	"testing"

	dbm "github.com/bhojpur/state/pkg/database"
	"github.com/gogo/protobuf/proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/bhojpur/state/internal/pubsub/query"
	"github.com/bhojpur/state/internal/state/indexer"
	abcipb "github.com/bhojpur/state/pkg/abci/types"
	librand "github.com/bhojpur/state/pkg/libs/rand"
	"github.com/bhojpur/state/pkg/types"
)

func TestTxIndex(t *testing.T) {
	txIndexer := NewTxIndex(dbm.NewMemDB())

	tx := types.Tx("HELLO WORLD")
	txResult := &abcipb.TxResult{
		Height: 1,
		Index:  0,
		Tx:     tx,
		Result: abcipb.ExecTxResult{
			Data: []byte{0},
			Code: abcipb.CodeTypeOK, Log: "", Events: nil,
		},
	}
	hash := tx.Hash()

	batch := indexer.NewBatch(1)
	if err := batch.Add(txResult); err != nil {
		t.Error(err)
	}
	err := txIndexer.Index(batch.Ops)
	require.NoError(t, err)

	loadedTxResult, err := txIndexer.Get(hash)
	require.NoError(t, err)
	assert.True(t, proto.Equal(txResult, loadedTxResult))

	tx2 := types.Tx("BYE BYE WORLD")
	txResult2 := &abcipb.TxResult{
		Height: 1,
		Index:  0,
		Tx:     tx2,
		Result: abcipb.ExecTxResult{
			Data: []byte{0},
			Code: abcipb.CodeTypeOK, Log: "", Events: nil,
		},
	}
	hash2 := tx2.Hash()

	err = txIndexer.Index([]*abcipb.TxResult{txResult2})
	require.NoError(t, err)

	loadedTxResult2, err := txIndexer.Get(hash2)
	require.NoError(t, err)
	assert.True(t, proto.Equal(txResult2, loadedTxResult2))
}

func TestTxSearch(t *testing.T) {
	indexer := NewTxIndex(dbm.NewMemDB())

	txResult := txResultWithEvents([]abcipb.Event{
		{Type: "account", Attributes: []abcipb.EventAttribute{{Key: "number", Value: "1", Index: true}}},
		{Type: "account", Attributes: []pb.EventAttribute{{Key: "owner", Value: "Ivan", Index: true}}},
		{Type: "", Attributes: []abcipb.EventAttribute{{Key: "not_allowed", Value: "Vlad", Index: true}}},
	})
	hash := types.Tx(txResult.Tx).Hash()

	err := indexer.Index([]*abcipb.TxResult{txResult})
	require.NoError(t, err)

	testCases := []struct {
		q             string
		resultsLength int
	}{
		// search by hash
		{fmt.Sprintf("tx.hash = '%X'", hash), 1},
		// search by exact match (one key)
		{"account.number = 1", 1},
		// search by exact match (two keys)
		{"account.number = 1 AND account.owner = 'Ivan'", 1},
		// search by exact match (two keys)
		{"account.number = 1 AND account.owner = 'Vlad'", 0},
		{"account.owner = 'Vlad' AND account.number = 1", 0},
		{"account.number >= 1 AND account.owner = 'Vlad'", 0},
		{"account.owner = 'Vlad' AND account.number >= 1", 0},
		{"account.number <= 0", 0},
		{"account.number <= 0 AND account.owner = 'Ivan'", 0},
		// search using a prefix of the stored value
		{"account.owner = 'Iv'", 0},
		// search by range
		{"account.number >= 1 AND account.number <= 5", 1},
		// search by range (lower bound)
		{"account.number >= 1", 1},
		// search by range (upper bound)
		{"account.number <= 5", 1},
		// search using not allowed key
		{"not_allowed = 'boom'", 0},
		// search for not existing tx result
		{"account.number >= 2 AND account.number <= 5", 0},
		// search using not existing key
		{"account.date >= TIME 2013-05-03T14:45:00Z", 0},
		// search using CONTAINS
		{"account.owner CONTAINS 'an'", 1},
		// search for non existing value using CONTAINS
		{"account.owner CONTAINS 'Vlad'", 0},
		// search using the wrong key (of numeric type) using CONTAINS
		{"account.number CONTAINS 'Iv'", 0},
		// search using EXISTS
		{"account.number EXISTS", 1},
		// search using EXISTS for non existing key
		{"account.date EXISTS", 0},
		// search using height
		{"account.number = 1 AND tx.height = 1", 1},
		// search using incorrect height
		{"account.number = 1 AND tx.height = 3", 0},
		// search using height only
		{"tx.height = 1", 1},
	}

	ctx := context.Background()

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.q, func(t *testing.T) {
			results, err := indexer.Search(ctx, query.MustCompile(tc.q))
			assert.NoError(t, err)

			assert.Len(t, results, tc.resultsLength)
			if tc.resultsLength > 0 {
				for _, txr := range results {
					assert.True(t, proto.Equal(txResult, txr))
				}
			}
		})
	}
}

func TestTxSearchWithCancelation(t *testing.T) {
	indexer := NewTxIndex(dbm.NewMemDB())

	txResult := txResultWithEvents([]abcipb.Event{
		{Type: "account", Attributes: []abcipb.EventAttribute{{Key: "number", Value: "1", Index: true}}},
		{Type: "account", Attributes: []abcipb.EventAttribute{{Key: "owner", Value: "Ivan", Index: true}}},
		{Type: "", Attributes: []abcipb.EventAttribute{{Key: "not_allowed", Value: "Vlad", Index: true}}},
	})
	err := indexer.Index([]*abcipb.TxResult{txResult})
	require.NoError(t, err)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	results, err := indexer.Search(ctx, query.MustCompile(`account.number = 1`))
	assert.NoError(t, err)
	assert.Empty(t, results)
}

func TestTxSearchDeprecatedIndexing(t *testing.T) {
	indexer := NewTxIndex(dbm.NewMemDB())

	// index tx using events indexing (composite key)
	txResult1 := txResultWithEvents([]abcipb.Event{
		{Type: "account", Attributes: []abcipb.EventAttribute{{Key: "number", Value: "1", Index: true}}},
	})
	hash1 := types.Tx(txResult1.Tx).Hash()

	err := indexer.Index([]*abcipb.TxResult{txResult1})
	require.NoError(t, err)

	// index tx also using deprecated indexing (event as key)
	txResult2 := txResultWithEvents(nil)
	txResult2.Tx = types.Tx("HELLO WORLD 2")

	hash2 := types.Tx(txResult2.Tx).Hash()
	b := indexer.store.NewBatch()

	rawBytes, err := proto.Marshal(txResult2)
	require.NoError(t, err)

	depKey := []byte(fmt.Sprintf("%s/%s/%d/%d",
		"sender",
		"addr1",
		txResult2.Height,
		txResult2.Index,
	))

	err = b.Set(depKey, hash2)
	require.NoError(t, err)
	err = b.Set(KeyFromHeight(txResult2), hash2)
	require.NoError(t, err)
	err = b.Set(hash2, rawBytes)
	require.NoError(t, err)
	err = b.Write()
	require.NoError(t, err)

	testCases := []struct {
		q       string
		results []*abcipb.TxResult
	}{
		// search by hash
		{fmt.Sprintf("tx.hash = '%X'", hash1), []*abcipb.TxResult{txResult1}},
		// search by hash
		{fmt.Sprintf("tx.hash = '%X'", hash2), []*abcipb.TxResult{txResult2}},
		// search by exact match (one key)
		{"account.number = 1", []*abcipb.TxResult{txResult1}},
		{"account.number >= 1 AND account.number <= 5", []*abcipb.TxResult{txResult1}},
		// search by range (lower bound)
		{"account.number >= 1", []*abcipb.TxResult{txResult1}},
		// search by range (upper bound)
		{"account.number <= 5", []*abcipb.TxResult{txResult1}},
		// search using not allowed key
		{"not_allowed = 'boom'", []*abcipb.TxResult{}},
		// search for not existing tx result
		{"account.number >= 2 AND account.number <= 5", []*abcipb.TxResult{}},
		// search using not existing key
		{"account.date >= TIME 2013-05-03T14:45:00Z", []*abcipb.TxResult{}},
		// search by deprecated key
		{"sender = 'addr1'", []*abcipb.TxResult{txResult2}},
	}

	ctx := context.Background()

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.q, func(t *testing.T) {
			results, err := indexer.Search(ctx, query.MustCompile(tc.q))
			require.NoError(t, err)
			for _, txr := range results {
				for _, tr := range tc.results {
					assert.True(t, proto.Equal(tr, txr))
				}
			}
		})
	}
}

func TestTxSearchOneTxWithMultipleSameTagsButDifferentValues(t *testing.T) {
	indexer := NewTxIndex(dbm.NewMemDB())

	txResult := txResultWithEvents([]abcipb.Event{
		{Type: "account", Attributes: []abcipb.EventAttribute{{Key: "number", Value: "1", Index: true}}},
		{Type: "account", Attributes: []abcipb.EventAttribute{{Key: "number", Value: "2", Index: true}}},
	})

	err := indexer.Index([]*abcipb.TxResult{txResult})
	require.NoError(t, err)

	ctx := context.Background()

	results, err := indexer.Search(ctx, query.MustCompile(`account.number >= 1`))
	assert.NoError(t, err)

	assert.Len(t, results, 1)
	for _, txr := range results {
		assert.True(t, proto.Equal(txResult, txr))
	}
}

func TestTxSearchMultipleTxs(t *testing.T) {
	indexer := NewTxIndex(dbm.NewMemDB())

	// indexed first, but bigger height (to test the order of transactions)
	txResult := txResultWithEvents([]abcipb.Event{
		{Type: "account", Attributes: []abcipb.EventAttribute{{Key: "number", Value: "1", Index: true}}},
	})

	txResult.Tx = types.Tx("Bob's account")
	txResult.Height = 2
	txResult.Index = 1
	err := indexer.Index([]*abcipb.TxResult{txResult})
	require.NoError(t, err)

	// indexed second, but smaller height (to test the order of transactions)
	txResult2 := txResultWithEvents([]abcipb.Event{
		{Type: "account", Attributes: []abcipb.EventAttribute{{Key: "number", Value: "2", Index: true}}},
	})
	txResult2.Tx = types.Tx("Alice's account")
	txResult2.Height = 1
	txResult2.Index = 2

	err = indexer.Index([]*abcipb.TxResult{txResult2})
	require.NoError(t, err)

	// indexed third (to test the order of transactions)
	txResult3 := txResultWithEvents([]abcipb.Event{
		{Type: "account", Attributes: []abcipb.EventAttribute{{Key: "number", Value: "3", Index: true}}},
	})
	txResult3.Tx = types.Tx("Jack's account")
	txResult3.Height = 1
	txResult3.Index = 1
	err = indexer.Index([]*abcipb.TxResult{txResult3})
	require.NoError(t, err)

	// indexed fourth (to test we don't include txs with similar events)
	txResult4 := txResultWithEvents([]abcipb.Event{
		{Type: "account", Attributes: []abcipb.EventAttribute{{Key: "number.id", Value: "1", Index: true}}},
	})
	txResult4.Tx = types.Tx("Mike's account")
	txResult4.Height = 2
	txResult4.Index = 2
	err = indexer.Index([]*abcipb.TxResult{txResult4})
	require.NoError(t, err)

	ctx := context.Background()

	results, err := indexer.Search(ctx, query.MustCompile(`account.number >= 1`))
	assert.NoError(t, err)

	require.Len(t, results, 3)
}

func txResultWithEvents(events []abcipb.Event) *abcipb.TxResult {
	tx := types.Tx("HELLO WORLD")
	return &abcipb.TxResult{
		Height: 1,
		Index:  0,
		Tx:     tx,
		Result: abcipb.ExecTxResult{
			Data:   []byte{0},
			Code:   abcipb.CodeTypeOK,
			Log:    "",
			Events: events,
		},
	}
}

func benchmarkTxIndex(txsCount int64, b *testing.B) {
	dir := b.TempDir()

	store, err := dbm.NewDB("tx_index", "goleveldb", dir)
	require.NoError(b, err)
	txIndexer := NewTxIndex(store)

	batch := indexer.NewBatch(txsCount)
	txIndex := uint32(0)
	for i := int64(0); i < txsCount; i++ {
		tx := librand.Bytes(250)
		txResult := &abcipb.TxResult{
			Height: 1,
			Index:  txIndex,
			Tx:     tx,
			Result: abcipb.ExecTxResult{
				Data:   []byte{0},
				Code:   abcipb.CodeTypeOK,
				Log:    "",
				Events: []abcipb.Event{},
			},
		}
		if err := batch.Add(txResult); err != nil {
			b.Fatal(err)
		}
		txIndex++
	}

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		err = txIndexer.Index(batch.Ops)
	}
	if err != nil {
		b.Fatal(err)
	}
}

func BenchmarkTxIndex1(b *testing.B)     { benchmarkTxIndex(1, b) }
func BenchmarkTxIndex500(b *testing.B)   { benchmarkTxIndex(500, b) }
func BenchmarkTxIndex1000(b *testing.B)  { benchmarkTxIndex(1000, b) }
func BenchmarkTxIndex2000(b *testing.B)  { benchmarkTxIndex(2000, b) }
func BenchmarkTxIndex10000(b *testing.B) { benchmarkTxIndex(10000, b) }
