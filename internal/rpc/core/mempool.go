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
	"math/rand"
	"time"

	"github.com/bhojpur/state/internal/mempool"
	"github.com/bhojpur/state/internal/state/indexer"
	abci "github.com/bhojpur/state/pkg/abci/types"
	libmath "github.com/bhojpur/state/pkg/libs/math"
	"github.com/bhojpur/state/pkg/rpc/coretypes"
)

// NOTE: tx should be signed, but this is only checked at the app level (not by Bhojpur State!)

// BroadcastTxAsync returns right away, with no response. Does not wait for
// CheckTx nor DeliverTx results.
func (env *Environment) BroadcastTxAsync(ctx context.Context, req *coretypes.RequestBroadcastTx) (*coretypes.ResultBroadcastTx, error) {
	err := env.Mempool.CheckTx(ctx, req.Tx, nil, mempool.TxInfo{})
	if err != nil {
		return nil, err
	}

	return &coretypes.ResultBroadcastTx{Hash: req.Tx.Hash()}, nil
}

// BroadcastTxSync returns with the response from CheckTx. Does not wait for
// DeliverTx result.
func (env *Environment) BroadcastTxSync(ctx context.Context, req *coretypes.RequestBroadcastTx) (*coretypes.ResultBroadcastTx, error) {
	resCh := make(chan *abci.ResponseCheckTx, 1)
	err := env.Mempool.CheckTx(
		ctx,
		req.Tx,
		func(res *abci.ResponseCheckTx) {
			select {
			case <-ctx.Done():
			case resCh <- res:
			}
		},
		mempool.TxInfo{},
	)
	if err != nil {
		return nil, err
	}

	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("broadcast confirmation not received: %w", ctx.Err())
	case r := <-resCh:
		return &coretypes.ResultBroadcastTx{
			Code:         r.Code,
			Data:         r.Data,
			Log:          r.Log,
			Codespace:    r.Codespace,
			MempoolError: r.MempoolError,
			Hash:         req.Tx.Hash(),
		}, nil
	}
}

// BroadcastTxCommit returns with the responses from CheckTx and DeliverTx.
func (env *Environment) BroadcastTxCommit(ctx context.Context, req *coretypes.RequestBroadcastTx) (*coretypes.ResultBroadcastTxCommit, error) {
	resCh := make(chan *abci.ResponseCheckTx, 1)
	err := env.Mempool.CheckTx(
		ctx,
		req.Tx,
		func(res *abci.ResponseCheckTx) {
			select {
			case <-ctx.Done():
			case resCh <- res:
			}
		},
		mempool.TxInfo{},
	)
	if err != nil {
		return nil, err
	}

	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("broadcast confirmation not received: %w", ctx.Err())
	case r := <-resCh:
		if r.Code != abci.CodeTypeOK {
			return &coretypes.ResultBroadcastTxCommit{
				CheckTx: *r,
				Hash:    req.Tx.Hash(),
			}, fmt.Errorf("transaction encountered error (%s)", r.MempoolError)
		}

		if !indexer.KVSinkEnabled(env.EventSinks) {
			return &coretypes.ResultBroadcastTxCommit{
					CheckTx: *r,
					Hash:    req.Tx.Hash(),
				},
				errors.New("cannot confirm transaction because kvEventSink is not enabled")
		}

		startAt := time.Now()
		timer := time.NewTimer(0)
		defer timer.Stop()

		count := 0
		for {
			count++
			select {
			case <-ctx.Done():
				env.Logger.Error("error on broadcastTxCommit",
					"duration", time.Since(startAt),
					"err", err)
				return &coretypes.ResultBroadcastTxCommit{
						CheckTx: *r,
						Hash:    req.Tx.Hash(),
					}, fmt.Errorf("timeout waiting for commit of tx %s (%s)",
						req.Tx.Hash(), time.Since(startAt))
			case <-timer.C:
				txres, err := env.Tx(ctx, &coretypes.RequestTx{
					Hash:  req.Tx.Hash(),
					Prove: false,
				})
				if err != nil {
					jitter := 100*time.Millisecond + time.Duration(rand.Int63n(int64(time.Second))) // nolint: gosec
					backoff := 100 * time.Duration(count) * time.Millisecond
					timer.Reset(jitter + backoff)
					continue
				}

				return &coretypes.ResultBroadcastTxCommit{
					CheckTx:  *r,
					TxResult: txres.TxResult,
					Hash:     req.Tx.Hash(),
					Height:   txres.Height,
				}, nil
			}
		}
	}
}

// UnconfirmedTxs gets unconfirmed transactions from the mempool in order of priority
func (env *Environment) UnconfirmedTxs(ctx context.Context, req *coretypes.RequestUnconfirmedTxs) (*coretypes.ResultUnconfirmedTxs, error) {
	totalCount := env.Mempool.Size()
	perPage := env.validatePerPage(req.PerPage.IntPtr())
	page, err := validatePage(req.Page.IntPtr(), perPage, totalCount)
	if err != nil {
		return nil, err
	}

	skipCount := validateSkipCount(page, perPage)

	txs := env.Mempool.ReapMaxTxs(skipCount + libmath.MinInt(perPage, totalCount-skipCount))
	result := txs[skipCount:]

	return &coretypes.ResultUnconfirmedTxs{
		Count:      len(result),
		Total:      totalCount,
		TotalBytes: env.Mempool.SizeBytes(),
		Txs:        result,
	}, nil
}

// NumUnconfirmedTxs gets number of unconfirmed transactions.
func (env *Environment) NumUnconfirmedTxs(ctx context.Context) (*coretypes.ResultUnconfirmedTxs, error) {
	return &coretypes.ResultUnconfirmedTxs{
		Count:      env.Mempool.Size(),
		Total:      env.Mempool.Size(),
		TotalBytes: env.Mempool.SizeBytes()}, nil
}

// CheckTx checks the transaction without executing it. The transaction won't
// be added to the mempool either.
func (env *Environment) CheckTx(ctx context.Context, req *coretypes.RequestCheckTx) (*coretypes.ResultCheckTx, error) {
	res, err := env.ProxyApp.CheckTx(ctx, &abci.RequestCheckTx{Tx: req.Tx})
	if err != nil {
		return nil, err
	}
	return &coretypes.ResultCheckTx{ResponseCheckTx: *res}, nil
}

func (env *Environment) RemoveTx(ctx context.Context, req *coretypes.RequestRemoveTx) error {
	return env.Mempool.RemoveTxByKey(req.TxKey)
}
