package mock_test

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
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/bhojpur/state/pkg/libs/bytes"
	"github.com/bhojpur/state/pkg/rpc/client/mock"
	"github.com/bhojpur/state/pkg/rpc/coretypes"
)

func TestStatus(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	m := &mock.StatusMock{
		Call: mock.Call{
			Response: &coretypes.ResultStatus{
				SyncInfo: coretypes.SyncInfo{
					LatestBlockHash:     bytes.HexBytes("block"),
					LatestAppHash:       bytes.HexBytes("app"),
					LatestBlockHeight:   10,
					MaxPeerBlockHeight:  20,
					TotalSyncedTime:     time.Second,
					RemainingTime:       time.Minute,
					TotalSnapshots:      10,
					ChunkProcessAvgTime: time.Duration(10),
					SnapshotHeight:      10,
					SnapshotChunksCount: 9,
					SnapshotChunksTotal: 10,
					BackFilledBlocks:    9,
					BackFillBlocksTotal: 10,
				},
			}},
	}

	r := mock.NewStatusRecorder(m)
	require.Equal(t, 0, len(r.Calls))

	// make sure response works proper
	status, err := r.Status(ctx)
	require.NoError(t, err)
	assert.EqualValues(t, "block", status.SyncInfo.LatestBlockHash)
	assert.EqualValues(t, 10, status.SyncInfo.LatestBlockHeight)
	assert.EqualValues(t, 20, status.SyncInfo.MaxPeerBlockHeight)
	assert.EqualValues(t, time.Second, status.SyncInfo.TotalSyncedTime)
	assert.EqualValues(t, time.Minute, status.SyncInfo.RemainingTime)

	// make sure recorder works properly
	require.Equal(t, 1, len(r.Calls))
	rs := r.Calls[0]
	assert.Equal(t, "status", rs.Name)
	assert.Nil(t, rs.Args)
	assert.Nil(t, rs.Error)
	require.NotNil(t, rs.Response)
	st, ok := rs.Response.(*coretypes.ResultStatus)
	require.True(t, ok)
	assert.EqualValues(t, "block", st.SyncInfo.LatestBlockHash)
	assert.EqualValues(t, 10, st.SyncInfo.LatestBlockHeight)
	assert.EqualValues(t, 20, st.SyncInfo.MaxPeerBlockHeight)
	assert.EqualValues(t, time.Second, status.SyncInfo.TotalSyncedTime)
	assert.EqualValues(t, time.Minute, status.SyncInfo.RemainingTime)

	assert.EqualValues(t, 10, st.SyncInfo.TotalSnapshots)
	assert.EqualValues(t, time.Duration(10), st.SyncInfo.ChunkProcessAvgTime)
	assert.EqualValues(t, 10, st.SyncInfo.SnapshotHeight)
	assert.EqualValues(t, 9, status.SyncInfo.SnapshotChunksCount)
	assert.EqualValues(t, 10, status.SyncInfo.SnapshotChunksTotal)
	assert.EqualValues(t, 9, status.SyncInfo.BackFilledBlocks)
	assert.EqualValues(t, 10, status.SyncInfo.BackFillBlocksTotal)
}
