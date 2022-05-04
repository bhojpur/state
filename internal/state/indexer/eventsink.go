package indexer

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
	v1 "github.com/bhojpur/state/pkg/abci/types"
	"github.com/bhojpur/state/pkg/types"
)

type EventSinkType string

const (
	NULL EventSinkType = "null"
	KV   EventSinkType = "kv"
	PSQL EventSinkType = "psql"
)

//go:generate ../../../scripts/mockery_generate.sh EventSink

// EventSink interface is defined the APIs for the IndexerService to interact with the data store,
// including the block/transaction indexing and the search functions.
//
// The IndexerService will accept a list of one or more EventSink types. During the OnStart method
// it will call the appropriate APIs on each EventSink to index both block and transaction events.
type EventSink interface {

	// IndexBlockEvents indexes the blockheader.
	IndexBlockEvents(types.EventDataNewBlockHeader) error

	// IndexTxEvents indexes the given result of transactions. To call it with multi transactions,
	// must guarantee the index of given transactions are in order.
	IndexTxEvents([]*v1.TxResult) error

	// SearchBlockEvents provides the block search by given query conditions. This function only
	// supported by the kvEventSink.
	SearchBlockEvents(context.Context, *query.Query) ([]int64, error)

	// SearchTxEvents provides the transaction search by given query conditions. This function only
	// supported by the kvEventSink.
	SearchTxEvents(context.Context, *query.Query) ([]*v1.TxResult, error)

	// GetTxByHash provides the transaction search by given transaction hash. This function only
	// supported by the kvEventSink.
	GetTxByHash([]byte) (*v1.TxResult, error)

	// HasBlock provides the transaction search by given transaction hash. This function only
	// supported by the kvEventSink.
	HasBlock(int64) (bool, error)

	// Type checks the eventsink structure type.
	Type() EventSinkType

	// Stop will close the data store connection, if the eventsink supports it.
	Stop() error
}
