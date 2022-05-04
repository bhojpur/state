package null

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

	"github.com/bhojpur/state/internal/pubsub/query"
	"github.com/bhojpur/state/internal/state/indexer"
	abci "github.com/bhojpur/state/pkg/abci/types"
	"github.com/bhojpur/state/pkg/types"
)

var _ indexer.EventSink = (*EventSink)(nil)

// EventSink implements a no-op indexer.
type EventSink struct{}

func NewEventSink() indexer.EventSink {
	return &EventSink{}
}

func (nes *EventSink) Type() indexer.EventSinkType {
	return indexer.NULL
}

func (nes *EventSink) IndexBlockEvents(bh types.EventDataNewBlockHeader) error {
	return nil
}

func (nes *EventSink) IndexTxEvents(results []*abci.TxResult) error {
	return nil
}

func (nes *EventSink) SearchBlockEvents(ctx context.Context, q *query.Query) ([]int64, error) {
	return nil, nil
}

func (nes *EventSink) SearchTxEvents(ctx context.Context, q *query.Query) ([]*abci.TxResult, error) {
	return nil, nil
}

func (nes *EventSink) GetTxByHash(hash []byte) (*abci.TxResult, error) {
	return nil, nil
}

func (nes *EventSink) HasBlock(h int64) (bool, error) {
	return false, nil
}

func (nes *EventSink) Stop() error {
	return nil
}
