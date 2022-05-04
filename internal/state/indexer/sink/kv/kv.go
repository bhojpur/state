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

	dbm "github.com/bhojpur/state/pkg/database"

	"github.com/bhojpur/state/internal/pubsub/query"
	"github.com/bhojpur/state/internal/state/indexer"
	kvb "github.com/bhojpur/state/internal/state/indexer/block/kv"
	kvt "github.com/bhojpur/state/internal/state/indexer/tx/kv"
	abci "github.com/bhojpur/state/pkg/abci/types"
	"github.com/bhojpur/state/pkg/types"
)

var _ indexer.EventSink = (*EventSink)(nil)

// The EventSink is an aggregator for redirecting the call path of the tx/block kvIndexer.
// For the implementation details please see the kv.go in the indexer/block and indexer/tx folder.
type EventSink struct {
	txi   *kvt.TxIndex
	bi    *kvb.BlockerIndexer
	store dbm.DB
}

func NewEventSink(store dbm.DB) indexer.EventSink {
	return &EventSink{
		txi:   kvt.NewTxIndex(store),
		bi:    kvb.New(store),
		store: store,
	}
}

func (kves *EventSink) Type() indexer.EventSinkType {
	return indexer.KV
}

func (kves *EventSink) IndexBlockEvents(bh types.EventDataNewBlockHeader) error {
	return kves.bi.Index(bh)
}

func (kves *EventSink) IndexTxEvents(results []*abci.TxResult) error {
	return kves.txi.Index(results)
}

func (kves *EventSink) SearchBlockEvents(ctx context.Context, q *query.Query) ([]int64, error) {
	return kves.bi.Search(ctx, q)
}

func (kves *EventSink) SearchTxEvents(ctx context.Context, q *query.Query) ([]*abci.TxResult, error) {
	return kves.txi.Search(ctx, q)
}

func (kves *EventSink) GetTxByHash(hash []byte) (*abci.TxResult, error) {
	return kves.txi.Get(hash)
}

func (kves *EventSink) HasBlock(h int64) (bool, error) {
	return kves.bi.Has(h)
}

func (kves *EventSink) Stop() error {
	return kves.store.Close()
}
