package p2p

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
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bhojpur/state/pkg/libs/log"
	"github.com/bhojpur/state/pkg/types"
)

func TestRouter_ConstructQueueFactory(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	t.Run("ValidateOptionsPopulatesDefaultQueue", func(t *testing.T) {
		opts := RouterOptions{}
		require.NoError(t, opts.Validate())
		require.Equal(t, "fifo", opts.QueueType)
	})
	t.Run("Default", func(t *testing.T) {
		require.Zero(t, os.Getenv("TM_P2P_QUEUE"))
		opts := RouterOptions{}
		r, err := NewRouter(log.NewNopLogger(), nil, nil, nil, func() *types.NodeInfo { return &types.NodeInfo{} }, nil, nil, opts)
		require.NoError(t, err)
		require.NoError(t, r.setupQueueFactory(ctx))

		_, ok := r.queueFactory(1).(*fifoQueue)
		require.True(t, ok)
	})
	t.Run("Fifo", func(t *testing.T) {
		opts := RouterOptions{QueueType: queueTypeFifo}
		r, err := NewRouter(log.NewNopLogger(), nil, nil, nil, func() *types.NodeInfo { return &types.NodeInfo{} }, nil, nil, opts)
		require.NoError(t, err)
		require.NoError(t, r.setupQueueFactory(ctx))

		_, ok := r.queueFactory(1).(*fifoQueue)
		require.True(t, ok)
	})
	t.Run("Priority", func(t *testing.T) {
		opts := RouterOptions{QueueType: queueTypePriority}
		r, err := NewRouter(log.NewNopLogger(), nil, nil, nil, func() *types.NodeInfo { return &types.NodeInfo{} }, nil, nil, opts)
		require.NoError(t, err)
		require.NoError(t, r.setupQueueFactory(ctx))

		q, ok := r.queueFactory(1).(*pqScheduler)
		require.True(t, ok)
		defer q.close()
	})
	t.Run("NonExistant", func(t *testing.T) {
		opts := RouterOptions{QueueType: "fast"}
		_, err := NewRouter(log.NewNopLogger(), nil, nil, nil, func() *types.NodeInfo { return &types.NodeInfo{} }, nil, nil, opts)
		require.Error(t, err)
		require.Contains(t, err.Error(), "fast")
	})
	t.Run("InternalsSafeWhenUnspecified", func(t *testing.T) {
		r := &Router{}
		require.Zero(t, r.options.QueueType)

		fn, err := r.createQueueFactory(ctx)
		require.Error(t, err)
		require.Nil(t, fn)
	})
}
