package core

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
	"errors"
	"fmt"
	"sort"

	tmquery "github.com/bhojpur/state/internal/pubsub/query"
	"github.com/bhojpur/state/internal/state/indexer"
	libmath "github.com/bhojpur/state/pkg/libs/math"
	"github.com/bhojpur/state/pkg/rpc/coretypes"
	"github.com/bhojpur/state/pkg/types"
)

// Tx allows you to query the transaction results. `nil` could mean the
// transaction is in the mempool, invalidated, or was not sent in the first
// place.
func (env *Environment) Tx(ctx context.Context, req *coretypes.RequestTx) (*coretypes.ResultTx, error) {
	// if index is disabled, return error
	if !indexer.KVSinkEnabled(env.EventSinks) {
		return nil, errors.New("transaction querying is disabled due to no kvEventSink")
	}

	for _, sink := range env.EventSinks {
		if sink.Type() == indexer.KV {
			r, err := sink.GetTxByHash(req.Hash)
			if r == nil {
				return nil, fmt.Errorf("tx (%X) not found, err: %w", req.Hash, err)
			}

			var proof types.TxProof
			if req.Prove {
				block := env.BlockStore.LoadBlock(r.Height)
				proof = block.Data.Txs.Proof(int(r.Index))
			}

			return &coretypes.ResultTx{
				Hash:     req.Hash,
				Height:   r.Height,
				Index:    r.Index,
				TxResult: r.Result,
				Tx:       r.Tx,
				Proof:    proof,
			}, nil
		}
	}

	return nil, fmt.Errorf("transaction querying is disabled on this node due to the KV event sink being disabled")
}

// TxSearch allows you to query for multiple transactions results. It returns a
// list of transactions (maximum ?per_page entries) and the total count.
func (env *Environment) TxSearch(ctx context.Context, req *coretypes.RequestTxSearch) (*coretypes.ResultTxSearch, error) {
	if !indexer.KVSinkEnabled(env.EventSinks) {
		return nil, fmt.Errorf("transaction searching is disabled due to no kvEventSink")
	} else if len(req.Query) > maxQueryLength {
		return nil, errors.New("maximum query length exceeded")
	}

	q, err := tmquery.New(req.Query)
	if err != nil {
		return nil, err
	}

	for _, sink := range env.EventSinks {
		if sink.Type() == indexer.KV {
			results, err := sink.SearchTxEvents(ctx, q)
			if err != nil {
				return nil, err
			}

			// sort results (must be done before pagination)
			switch req.OrderBy {
			case "desc", "":
				sort.Slice(results, func(i, j int) bool {
					if results[i].Height == results[j].Height {
						return results[i].Index > results[j].Index
					}
					return results[i].Height > results[j].Height
				})
			case "asc":
				sort.Slice(results, func(i, j int) bool {
					if results[i].Height == results[j].Height {
						return results[i].Index < results[j].Index
					}
					return results[i].Height < results[j].Height
				})
			default:
				return nil, fmt.Errorf("expected order_by to be either `asc` or `desc` or empty: %w", coretypes.ErrInvalidRequest)
			}

			// paginate results
			totalCount := len(results)
			perPage := env.validatePerPage(req.PerPage.IntPtr())

			page, err := validatePage(req.Page.IntPtr(), perPage, totalCount)
			if err != nil {
				return nil, err
			}

			skipCount := validateSkipCount(page, perPage)
			pageSize := libmath.MinInt(perPage, totalCount-skipCount)

			apiResults := make([]*coretypes.ResultTx, 0, pageSize)
			for i := skipCount; i < skipCount+pageSize; i++ {
				r := results[i]

				var proof types.TxProof
				if req.Prove {
					block := env.BlockStore.LoadBlock(r.Height)
					proof = block.Data.Txs.Proof(int(r.Index))
				}

				apiResults = append(apiResults, &coretypes.ResultTx{
					Hash:     types.Tx(r.Tx).Hash(),
					Height:   r.Height,
					Index:    r.Index,
					TxResult: r.Result,
					Tx:       r.Tx,
					Proof:    proof,
				})
			}

			return &coretypes.ResultTxSearch{Txs: apiResults, TotalCount: totalCount}, nil
		}
	}

	return nil, fmt.Errorf("transaction searching is disabled on this node due to the KV event sink being disabled")
}
