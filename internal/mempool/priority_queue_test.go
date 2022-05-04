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
	"math/rand"
	"sort"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestTxPriorityQueue(t *testing.T) {
	pq := NewTxPriorityQueue()
	numTxs := 1000

	priorities := make([]int, numTxs)

	var wg sync.WaitGroup
	for i := 1; i <= numTxs; i++ {
		priorities[i-1] = i
		wg.Add(1)

		go func(i int) {
			pq.PushTx(&WrappedTx{
				priority:  int64(i),
				timestamp: time.Now(),
			})

			wg.Done()
		}(i)
	}

	sort.Sort(sort.Reverse(sort.IntSlice(priorities)))

	wg.Wait()
	require.Equal(t, numTxs, pq.NumTxs())

	// Wait a second and push a tx with a duplicate priority
	time.Sleep(time.Second)
	now := time.Now()
	pq.PushTx(&WrappedTx{
		priority:  1000,
		timestamp: now,
	})
	require.Equal(t, 1001, pq.NumTxs())

	tx := pq.PopTx()
	require.Equal(t, 1000, pq.NumTxs())
	require.Equal(t, int64(1000), tx.priority)
	require.NotEqual(t, now, tx.timestamp)

	gotPriorities := make([]int, 0)
	for pq.NumTxs() > 0 {
		gotPriorities = append(gotPriorities, int(pq.PopTx().priority))
	}

	require.Equal(t, priorities, gotPriorities)
}

func TestTxPriorityQueue_GetEvictableTxs(t *testing.T) {
	pq := NewTxPriorityQueue()
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	values := make([]int, 1000)

	for i := 0; i < 1000; i++ {
		tx := make([]byte, 5) // each tx is 5 bytes
		_, err := rng.Read(tx)
		require.NoError(t, err)

		x := rng.Intn(100000)
		pq.PushTx(&WrappedTx{
			tx:       tx,
			priority: int64(x),
		})

		values[i] = x
	}

	sort.Ints(values)

	max := values[len(values)-1]
	min := values[0]
	totalSize := int64(len(values) * 5)

	testCases := []struct {
		name                             string
		priority, txSize, totalSize, cap int64
		expectedLen                      int
	}{
		{
			name:        "larest priority; single tx",
			priority:    int64(max + 1),
			txSize:      5,
			totalSize:   totalSize,
			cap:         totalSize,
			expectedLen: 1,
		},
		{
			name:        "larest priority; multi tx",
			priority:    int64(max + 1),
			txSize:      17,
			totalSize:   totalSize,
			cap:         totalSize,
			expectedLen: 4,
		},
		{
			name:        "larest priority; out of capacity",
			priority:    int64(max + 1),
			txSize:      totalSize + 1,
			totalSize:   totalSize,
			cap:         totalSize,
			expectedLen: 0,
		},
		{
			name:        "smallest priority; no tx",
			priority:    int64(min - 1),
			txSize:      5,
			totalSize:   totalSize,
			cap:         totalSize,
			expectedLen: 0,
		},
		{
			name:        "small priority; no tx",
			priority:    int64(min),
			txSize:      5,
			totalSize:   totalSize,
			cap:         totalSize,
			expectedLen: 0,
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			evictTxs := pq.GetEvictableTxs(tc.priority, tc.txSize, tc.totalSize, tc.cap)
			require.Len(t, evictTxs, tc.expectedLen)
		})
	}
}

func TestTxPriorityQueue_RemoveTx(t *testing.T) {
	pq := NewTxPriorityQueue()
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	numTxs := 1000

	values := make([]int, numTxs)

	for i := 0; i < numTxs; i++ {
		x := rng.Intn(100000)
		pq.PushTx(&WrappedTx{
			priority: int64(x),
		})

		values[i] = x
	}

	require.Equal(t, numTxs, pq.NumTxs())

	sort.Ints(values)
	max := values[len(values)-1]

	wtx := pq.txs[pq.NumTxs()/2]
	pq.RemoveTx(wtx)
	require.Equal(t, numTxs-1, pq.NumTxs())
	require.Equal(t, int64(max), pq.PopTx().priority)
	require.Equal(t, numTxs-2, pq.NumTxs())

	require.NotPanics(t, func() {
		pq.RemoveTx(&WrappedTx{heapIndex: numTxs})
		pq.RemoveTx(&WrappedTx{heapIndex: numTxs + 1})
	})
	require.Equal(t, numTxs-2, pq.NumTxs())
}
