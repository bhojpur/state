package kv_test

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
	"fmt"
	"testing"

	dbm "github.com/bhojpur/state/pkg/database"
	"github.com/stretchr/testify/require"

	"github.com/bhojpur/state/internal/pubsub/query"
	blockidxkv "github.com/bhojpur/state/internal/state/indexer/block/kv"
	abci "github.com/bhojpur/state/pkg/abci/types"
	"github.com/bhojpur/state/pkg/types"
)

func TestBlockIndexer(t *testing.T) {
	store := dbm.NewPrefixDB(dbm.NewMemDB(), []byte("block_events"))
	indexer := blockidxkv.New(store)

	require.NoError(t, indexer.Index(types.EventDataNewBlockHeader{
		Header: types.Header{Height: 1},
		ResultFinalizeBlock: abci.ResponseFinalizeBlock{
			Events: []abci.Event{
				{
					Type: "finalize_event1",
					Attributes: []abci.EventAttribute{
						{
							Key:   "proposer",
							Value: "FCAA001",
							Index: true,
						},
					},
				},
				{
					Type: "finalize_event2",
					Attributes: []abci.EventAttribute{
						{
							Key:   "foo",
							Value: "100",
							Index: true,
						},
					},
				},
			},
		},
	}))

	for i := 2; i < 12; i++ {
		var index bool
		if i%2 == 0 {
			index = true
		}
		require.NoError(t, indexer.Index(types.EventDataNewBlockHeader{
			Header: types.Header{Height: int64(i)},
			ResultFinalizeBlock: abci.ResponseFinalizeBlock{
				Events: []abci.Event{
					{
						Type: "finalize_event1",
						Attributes: []abci.EventAttribute{
							{
								Key:   "proposer",
								Value: "FCAA001",
								Index: true,
							},
						},
					},
					{
						Type: "finalize_event2",
						Attributes: []abci.EventAttribute{
							{
								Key:   "foo",
								Value: fmt.Sprintf("%d", i),
								Index: index,
							},
						},
					},
				},
			},
		}))
	}

	testCases := map[string]struct {
		q       *query.Query
		results []int64
	}{
		"block.height = 100": {
			q:       query.MustCompile(`block.height = 100`),
			results: []int64{},
		},
		"block.height = 5": {
			q:       query.MustCompile(`block.height = 5`),
			results: []int64{5},
		},
		"finalize_event.key1 = 'value1'": {
			q:       query.MustCompile(`finalize_event1.key1 = 'value1'`),
			results: []int64{},
		},
		"finalize_event.proposer = 'FCAA001'": {
			q:       query.MustCompile(`finalize_event1.proposer = 'FCAA001'`),
			results: []int64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11},
		},
		"finalize_event.foo <= 5": {
			q:       query.MustCompile(`finalize_event2.foo <= 5`),
			results: []int64{2, 4},
		},
		"finalize_event.foo >= 100": {
			q:       query.MustCompile(`finalize_event2.foo >= 100`),
			results: []int64{1},
		},
		"block.height > 2 AND finalize_event2.foo <= 8": {
			q:       query.MustCompile(`block.height > 2 AND finalize_event2.foo <= 8`),
			results: []int64{4, 6, 8},
		},
		"finalize_event.proposer CONTAINS 'FFFFFFF'": {
			q:       query.MustCompile(`finalize_event1.proposer CONTAINS 'FFFFFFF'`),
			results: []int64{},
		},
		"finalize_event.proposer CONTAINS 'FCAA001'": {
			q:       query.MustCompile(`finalize_event1.proposer CONTAINS 'FCAA001'`),
			results: []int64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11},
		},
	}

	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			results, err := indexer.Search(ctx, tc.q)
			require.NoError(t, err)
			require.Equal(t, tc.results, results)
		})
	}
}
