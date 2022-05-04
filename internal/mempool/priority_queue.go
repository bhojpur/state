package mempool

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
	"container/heap"
	"sort"
	"sync"
)

var _ heap.Interface = (*TxPriorityQueue)(nil)

// TxPriorityQueue defines a thread-safe priority queue for valid transactions.
type TxPriorityQueue struct {
	mtx sync.RWMutex
	txs []*WrappedTx
}

func NewTxPriorityQueue() *TxPriorityQueue {
	pq := &TxPriorityQueue{
		txs: make([]*WrappedTx, 0),
	}

	heap.Init(pq)

	return pq
}

// GetEvictableTxs attempts to find and return a list of *WrappedTx than can be
// evicted to make room for another *WrappedTx with higher priority. If no such
// list of *WrappedTx exists, nil will be returned. The returned list of *WrappedTx
// indicate that these transactions can be removed due to them being of lower
// priority and that their total sum in size allows room for the incoming
// transaction according to the mempool's configured limits.
func (pq *TxPriorityQueue) GetEvictableTxs(priority, txSize, totalSize, cap int64) []*WrappedTx {
	pq.mtx.RLock()
	defer pq.mtx.RUnlock()

	txs := make([]*WrappedTx, len(pq.txs))
	copy(txs, pq.txs)

	sort.Slice(txs, func(i, j int) bool {
		return txs[i].priority < txs[j].priority
	})

	var (
		toEvict []*WrappedTx
		i       int
	)

	currSize := totalSize

	// Loop over all transactions in ascending priority order evaluating those
	// that are only of less priority than the provided argument. We continue
	// evaluating transactions until there is sufficient capacity for the new
	// transaction (size) as defined by txSize.
	for i < len(txs) && txs[i].priority < priority {
		toEvict = append(toEvict, txs[i])
		currSize -= int64(txs[i].Size())

		if currSize+txSize <= cap {
			return toEvict
		}

		i++
	}

	return nil
}

// NumTxs returns the number of transactions in the priority queue. It is
// thread safe.
func (pq *TxPriorityQueue) NumTxs() int {
	pq.mtx.RLock()
	defer pq.mtx.RUnlock()

	return len(pq.txs)
}

// RemoveTx removes a specific transaction from the priority queue.
func (pq *TxPriorityQueue) RemoveTx(tx *WrappedTx) {
	pq.mtx.Lock()
	defer pq.mtx.Unlock()

	if tx.heapIndex < len(pq.txs) {
		heap.Remove(pq, tx.heapIndex)
	}
}

// PushTx adds a valid transaction to the priority queue. It is thread safe.
func (pq *TxPriorityQueue) PushTx(tx *WrappedTx) {
	pq.mtx.Lock()
	defer pq.mtx.Unlock()

	heap.Push(pq, tx)
}

// PopTx removes the top priority transaction from the queue. It is thread safe.
func (pq *TxPriorityQueue) PopTx() *WrappedTx {
	pq.mtx.Lock()
	defer pq.mtx.Unlock()

	x := heap.Pop(pq)
	if x != nil {
		return x.(*WrappedTx)
	}

	return nil
}

// Push implements the Heap interface.
//
// NOTE: A caller should never call Push. Use PushTx instead.
func (pq *TxPriorityQueue) Push(x interface{}) {
	n := len(pq.txs)
	item := x.(*WrappedTx)
	item.heapIndex = n
	pq.txs = append(pq.txs, item)
}

// Pop implements the Heap interface.
//
// NOTE: A caller should never call Pop. Use PopTx instead.
func (pq *TxPriorityQueue) Pop() interface{} {
	old := pq.txs
	n := len(old)
	item := old[n-1]
	old[n-1] = nil      // avoid memory leak
	item.heapIndex = -1 // for safety
	pq.txs = old[0 : n-1]
	return item
}

// Len implements the Heap interface.
//
// NOTE: A caller should never call Len. Use NumTxs instead.
func (pq *TxPriorityQueue) Len() int {
	return len(pq.txs)
}

// Less implements the Heap interface. It returns true if the transaction at
// position i in the queue is of less priority than the transaction at position j.
func (pq *TxPriorityQueue) Less(i, j int) bool {
	// If there exists two transactions with the same priority, consider the one
	// that we saw the earliest as the higher priority transaction.
	if pq.txs[i].priority == pq.txs[j].priority {
		return pq.txs[i].timestamp.Before(pq.txs[j].timestamp)
	}

	// We want Pop to give us the highest, not lowest, priority so we use greater
	// than here.
	return pq.txs[i].priority > pq.txs[j].priority
}

// Swap implements the Heap interface. It swaps two transactions in the queue.
func (pq *TxPriorityQueue) Swap(i, j int) {
	pq.txs[i], pq.txs[j] = pq.txs[j], pq.txs[i]
	pq.txs[i].heapIndex = i
	pq.txs[j].heapIndex = j
}
