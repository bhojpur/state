package state

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
	"sync"
	"time"

	"github.com/bhojpur/state/internal/mempool"
	abci "github.com/bhojpur/state/pkg/abci/types"
	"github.com/bhojpur/state/pkg/types"
)

func cachingStateFetcher(store Store) func() (State, error) {
	const ttl = time.Second

	var (
		last  time.Time
		mutex = &sync.Mutex{}
		cache State
		err   error
	)

	return func() (State, error) {
		mutex.Lock()
		defer mutex.Unlock()

		if time.Since(last) < ttl && cache.ChainID != "" {
			return cache, nil
		}

		cache, err = store.Load()
		if err != nil {
			return State{}, err
		}
		last = time.Now()

		return cache, nil
	}

}

// TxPreCheckFromStore returns a function to filter transactions before processing.
// The function limits the size of a transaction to the block's maximum data size.
func TxPreCheckFromStore(store Store) mempool.PreCheckFunc {
	fetch := cachingStateFetcher(store)

	return func(tx types.Tx) error {
		state, err := fetch()
		if err != nil {
			return err
		}

		return TxPreCheckForState(state)(tx)
	}
}

func TxPreCheckForState(state State) mempool.PreCheckFunc {
	return func(tx types.Tx) error {
		maxDataBytes := types.MaxDataBytesNoEvidence(
			state.ConsensusParams.Block.MaxBytes,
			state.Validators.Size(),
		)
		return mempool.PreCheckMaxBytes(maxDataBytes)(tx)
	}

}

// TxPostCheckFromStore returns a function to filter transactions after processing.
// The function limits the gas wanted by a transaction to the block's maximum total gas.
func TxPostCheckFromStore(store Store) mempool.PostCheckFunc {
	fetch := cachingStateFetcher(store)

	return func(tx types.Tx, resp *abci.ResponseCheckTx) error {
		state, err := fetch()
		if err != nil {
			return err
		}
		return mempool.PostCheckMaxGas(state.ConsensusParams.Block.MaxGas)(tx, resp)
	}
}

func TxPostCheckForState(state State) mempool.PostCheckFunc {
	return func(tx types.Tx, resp *abci.ResponseCheckTx) error {
		return mempool.PostCheckMaxGas(state.ConsensusParams.Block.MaxGas)(tx, resp)
	}
}
