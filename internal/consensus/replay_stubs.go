package consensus

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

	"github.com/bhojpur/state/internal/libs/clist"
	"github.com/bhojpur/state/internal/mempool"
	"github.com/bhojpur/state/internal/proxy"
	abciclient "github.com/bhojpur/state/pkg/abci/client"
	abci "github.com/bhojpur/state/pkg/abci/types"
	v1 "github.com/bhojpur/state/pkg/api/v1/state"
	"github.com/bhojpur/state/pkg/libs/log"
	"github.com/bhojpur/state/pkg/types"
)

type emptyMempool struct{}

var _ mempool.Mempool = emptyMempool{}

func (emptyMempool) Lock()     {}
func (emptyMempool) Unlock()   {}
func (emptyMempool) Size() int { return 0 }
func (emptyMempool) CheckTx(context.Context, types.Tx, func(*abci.ResponseCheckTx), mempool.TxInfo) error {
	return nil
}
func (emptyMempool) RemoveTxByKey(txKey types.TxKey) error   { return nil }
func (emptyMempool) ReapMaxBytesMaxGas(_, _ int64) types.Txs { return types.Txs{} }
func (emptyMempool) ReapMaxTxs(n int) types.Txs              { return types.Txs{} }
func (emptyMempool) Update(
	_ context.Context,
	_ int64,
	_ types.Txs,
	_ []*abci.ExecTxResult,
	_ mempool.PreCheckFunc,
	_ mempool.PostCheckFunc,
) error {
	return nil
}
func (emptyMempool) Flush()                                 {}
func (emptyMempool) FlushAppConn(ctx context.Context) error { return nil }
func (emptyMempool) TxsAvailable() <-chan struct{}          { return make(chan struct{}) }
func (emptyMempool) EnableTxsAvailable()                    {}
func (emptyMempool) SizeBytes() int64                       { return 0 }

func (emptyMempool) TxsFront() *clist.CElement    { return nil }
func (emptyMempool) TxsWaitChan() <-chan struct{} { return nil }

func (emptyMempool) InitWAL() error { return nil }
func (emptyMempool) CloseWAL()      {}

// mockProxyApp uses ABCIResponses to give the right results.
//
// Useful because we don't want to call Commit() twice for the same block on
// the real app.

func newMockProxyApp(
	logger log.Logger,
	appHash []byte,
	abciResponses *v1.ABCIResponses,
) (abciclient.Client, error) {
	return proxy.New(abciclient.NewLocalClient(logger, &mockProxyApp{
		appHash:       appHash,
		abciResponses: abciResponses,
	}), logger, proxy.NopMetrics()), nil
}

type mockProxyApp struct {
	abci.BaseApplication

	appHash       []byte
	txCount       int
	abciResponses *v1.ABCIResponses
}

func (mock *mockProxyApp) FinalizeBlock(_ context.Context, req *abci.RequestFinalizeBlock) (*abci.ResponseFinalizeBlock, error) {
	r := mock.abciResponses.FinalizeBlock
	mock.txCount++
	if r == nil {
		return &abci.ResponseFinalizeBlock{}, nil
	}
	return r, nil
}

func (mock *mockProxyApp) Commit(context.Context) (*abci.ResponseCommit, error) {
	return &abci.ResponseCommit{Data: mock.appHash}, nil
}
