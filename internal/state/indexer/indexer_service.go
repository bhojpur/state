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
	"time"

	"github.com/bhojpur/state/internal/eventbus"
	"github.com/bhojpur/state/internal/pubsub"
	"github.com/bhojpur/state/pkg/libs/log"
	"github.com/bhojpur/state/pkg/libs/service"
	"github.com/bhojpur/state/pkg/types"
)

// Service connects event bus, transaction and block indexers together in
// order to index transactions and blocks coming from the event bus.
type Service struct {
	service.BaseService
	logger log.Logger

	eventSinks []EventSink
	eventBus   *eventbus.EventBus
	metrics    *Metrics

	currentBlock struct {
		header types.EventDataNewBlockHeader
		height int64
		batch  *Batch
	}
}

// NewService constructs a new indexer service from the given arguments.
func NewService(args ServiceArgs) *Service {
	is := &Service{
		logger:     args.Logger,
		eventSinks: args.Sinks,
		eventBus:   args.EventBus,
		metrics:    args.Metrics,
	}
	if is.metrics == nil {
		is.metrics = NopMetrics()
	}
	is.BaseService = *service.NewBaseService(args.Logger, "IndexerService", is)
	return is
}

// publish publishes a pubsub message to the service. The service blocks until
// the message has been fully processed.
func (is *Service) publish(msg pubsub.Message) error {
	// Indexing has three states. Initially, no block is in progress (WAIT) and
	// we expect a block header. Upon seeing a header, we are waiting for zero
	// or more transactions (GATHER). Once all the expected transactions have
	// been delivered (in some order), we are ready to index. After indexing a
	// block, we revert to the WAIT state for the next block.

	if is.currentBlock.batch == nil {
		// WAIT: Start a new block.
		hdr := msg.Data().(types.EventDataNewBlockHeader)
		is.currentBlock.header = hdr
		is.currentBlock.height = hdr.Header.Height
		is.currentBlock.batch = NewBatch(hdr.NumTxs)

		if hdr.NumTxs != 0 {
			return nil
		}
		// If the block does not expect any transactions, fall through and index
		// it immediately.  This shouldn't happen, but this check ensures we do
		// not get stuck if it does.
	}

	curr := is.currentBlock.batch
	if curr.Pending != 0 {
		// GATHER: Accumulate a transaction into the current block's batch.
		txResult := msg.Data().(types.EventDataTx).TxResult
		if err := curr.Add(&txResult); err != nil {
			is.logger.Error("failed to add tx to batch",
				"height", is.currentBlock.height, "index", txResult.Index, "err", err)
		}

		// This may have been the last transaction in the batch, so fall through
		// to check whether it is time to index.
	}

	if curr.Pending == 0 {
		// INDEX: We have all the transactions we expect for the current block.
		for _, sink := range is.eventSinks {
			start := time.Now()
			if err := sink.IndexBlockEvents(is.currentBlock.header); err != nil {
				is.logger.Error("failed to index block header",
					"height", is.currentBlock.height, "err", err)
			} else {
				is.metrics.BlockEventsSeconds.Observe(time.Since(start).Seconds())
				is.metrics.BlocksIndexed.Add(1)
				is.logger.Debug("indexed block",
					"height", is.currentBlock.height, "sink", sink.Type())
			}

			if curr.Size() != 0 {
				start := time.Now()
				err := sink.IndexTxEvents(curr.Ops)
				if err != nil {
					is.logger.Error("failed to index block txs",
						"height", is.currentBlock.height, "err", err)
				} else {
					is.metrics.TxEventsSeconds.Observe(time.Since(start).Seconds())
					is.metrics.TransactionsIndexed.Add(float64(curr.Size()))
					is.logger.Debug("indexed txs",
						"height", is.currentBlock.height, "sink", sink.Type())
				}
			}
		}
		is.currentBlock.batch = nil // return to the WAIT state for the next block
	}

	return nil
}

// OnStart implements part of service.Service. It registers an observer for the
// indexer if the underlying event sinks support indexing.
//
// TODO: Can we get rid of the "enabled" check?
func (is *Service) OnStart(ctx context.Context) error {
	// If the event sinks support indexing, register an observer to capture
	// block header data for the indexer.
	if IndexingEnabled(is.eventSinks) {
		err := is.eventBus.Observe(ctx, is.publish,
			types.EventQueryNewBlockHeader, types.EventQueryTx)
		if err != nil {
			return err
		}
	}
	return nil
}

// OnStop implements service.Service by closing the event sinks.
func (is *Service) OnStop() {
	for _, sink := range is.eventSinks {
		if err := sink.Stop(); err != nil {
			is.logger.Error("failed to close eventsink", "eventsink", sink.Type(), "err", err)
		}
	}
}

// ServiceArgs are arguments for constructing a new indexer service.
type ServiceArgs struct {
	Sinks    []EventSink
	EventBus *eventbus.EventBus
	Metrics  *Metrics
	Logger   log.Logger
}

// KVSinkEnabled returns the given eventSinks is containing KVEventSink.
func KVSinkEnabled(sinks []EventSink) bool {
	for _, sink := range sinks {
		if sink.Type() == KV {
			return true
		}
	}

	return false
}

// IndexingEnabled returns the given eventSinks is supporting the indexing services.
func IndexingEnabled(sinks []EventSink) bool {
	for _, sink := range sinks {
		if sink.Type() == KV || sink.Type() == PSQL {
			return true
		}
	}

	return false
}
